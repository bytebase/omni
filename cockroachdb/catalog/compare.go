package catalog

import (
	"fmt"
	"go/constant"
	"go/token"
	"math/big"
	"strings"

	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	plpgsqlparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"

	metadatapb "github.com/bytebase/omni/metadata"
)

// The keys below compare two definitions of one object, one synced and one
// written by hand, as CockroachDB would store them. They are for comparison
// only; what a migration writes is the target's own text.

// exprKey is a check's or a predicate's expression without what CockroachDB
// adds to the literals it stores ('a'::STRING compared with a column,
// 's.seq'::REGCLASS), and without parentheses.
func exprKey(expr string, columns map[string]string) string {
	if strings.TrimSpace(expr) == "" {
		return ""
	}
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return expr
	}
	parsed, _ = tree.WalkExpr(literalCastStripper{columns}, parsed)
	return groupedKey(parsed)
}

// generationKey is exprKey for a computed column's expression, which
// CockroachDB stores with each literal annotated with the type it resolved
// it to (a * 2:::INT8, price * 100:::DECIMAL). An annotation that holds the
// literal exactly and that the literal resolves to bare goes
// (impliedAnnotation); one that picks another overload ('é':::BYTES) or may
// round (:::FLOAT8) stays.
func generationKey(expr string) string {
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return expr
	}
	parsed, err = tree.SimpleVisit(parsed, func(e tree.Expr) (bool, tree.Expr, error) {
		if a, ok := e.(*tree.AnnotateTypeExpr); ok && impliedAnnotation(unparen(a.Expr), a.Type) {
			return false, unparen(a.Expr), nil
		}
		return true, e, nil
	})
	if err != nil {
		return expr
	}
	return groupedKey(parsed)
}

// impliedAnnotation reports whether an annotation names a type its literal
// resolves to unannotated and holds it exactly: INT8 or DECIMAL on an
// integer, DECIMAL on a fraction, STRING on a string.
func impliedAnnotation(literal tree.Expr, ref tree.ResolvableTypeReference) bool {
	t, ok := tree.GetStaticallyKnownType(ref)
	if !ok {
		return false
	}
	switch l := literal.(type) {
	case *tree.NumVal:
		// An integer as INT or DECIMAL, or a number with a fraction as
		// DECIMAL, holds its value exactly; FLOAT may round it.
		switch t.Family() {
		case types.IntFamily:
			return !strings.ContainsAny(l.OrigString(), ".eE")
		case types.DecimalFamily:
			return true
		}
	case *tree.StrVal:
		return t.Identical(types.String)
	}
	return false
}

// columnTypes maps a table's columns to their types, for exprKey.
func columnTypes(columns []*metadatapb.ColumnMetadata) map[string]string {
	typed := make(map[string]string, len(columns))
	for _, column := range columns {
		typed[column.Name] = column.Type
	}
	return typed
}

// defaultKey is exprKey for a default of type typ. CockroachDB folds a cast
// of a literal to the column's own type ('x'::STRING to 'x'), adds COLLATE to
// a collated column's literal and takes ARRAY[]::STRING[] down to ARRAY[], so
// those go; a cast to another type converts the value first ('2020-01-01
// 12:34:56'::DATE in a TIMESTAMP is midnight) and stays.
func defaultKey(expr, typ string) string {
	if strings.TrimSpace(expr) == "" {
		return ""
	}
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return expr
	}
	if collate, ok := unparen(parsed).(*tree.CollateExpr); ok && isLiteral(collate.Expr) {
		parsed = collate.Expr
	}
	parsed = stripDefaultCast(parsed, typ)
	parsed, _ = tree.WalkExpr(literalCastStripper{}, parsed)
	return groupedKey(parsed)
}

// stripDefaultCast drops a default's cast or annotation on a literal where it
// names the default's own type, or the literal's own.
func stripDefaultCast(expr tree.Expr, typ string) tree.Expr {
	if cast, ok := unparen(expr).(*tree.CastExpr); ok && emptyArray(cast) {
		return unparen(cast.Expr)
	}
	return stripTypedCast(expr, typ)
}

// stripTypedCast drops a cast or annotation on a literal, or on each of a
// tuple's elements, where the literal takes the type it names anyway: the
// type typ of what it is assigned or compared to, or with typ unknown the
// literal's own type (a string's STRING), and unless the cast can change the
// value. A number's trailing fractional zeros go too, as CockroachDB states a
// FLOAT's default or comparand as a float (0 as 0.0).
func stripTypedCast(expr tree.Expr, typ string) tree.Expr {
	switch e := unparen(expr).(type) {
	case *tree.CastExpr:
		literal := unparen(e.Expr)
		if isLiteral(literal) && !castRounds(e) && (sameType(e.Type, typ) || typ == "" && ownType(literal, e.Type)) {
			return trimZeros(literal)
		}
	case *tree.AnnotateTypeExpr:
		literal := unparen(e.Expr)
		if isLiteral(literal) && (sameType(e.Type, typ) || ownType(literal, e.Type)) {
			return trimZeros(literal)
		}
		// An annotation converts nothing, so one naming the type the
		// expression is assigned to, as CockroachDB writes on a default
		// it stores (now():::TIMESTAMPTZ), says nothing.
		if _, call := literal.(*tree.FuncExpr); call && sameType(e.Type, typ) {
			return literal
		}
	case *tree.NumVal:
		return trimZeros(e)
	case *tree.Tuple:
		tuple := *e
		tuple.Exprs = make(tree.Exprs, len(e.Exprs))
		for i, element := range e.Exprs {
			tuple.Exprs[i] = stripTypedCast(element, typ)
		}
		return &tuple
	}
	return expr
}

// sameType reports whether a cast's type is typ, a type as the sync states
// it.
func sameType(ref tree.ResolvableTypeReference, typ string) bool {
	if typ == "" {
		return false
	}
	t, ok := tree.GetStaticallyKnownType(ref)
	if !ok {
		// A user-defined type, which CockroachDB writes on a default it
		// stores ('a':::public.mood), named by its schema as the sync names
		// a column's type.
		return schemaQualified(ref).SQLString() == typ
	}
	other, err := crparser.GetTypeFromValidSQLSyntax(typ)
	if err != nil {
		return false
	}
	o, ok := tree.GetStaticallyKnownType(other)
	return ok && t.Identical(o)
}

// emptyArray reports whether a cast is of an empty ARRAY[], whose element
// type the cast only names.
func emptyArray(cast *tree.CastExpr) bool {
	array, ok := unparen(cast.Expr).(*tree.Array)
	return ok && len(array.Exprs) == 0
}

// ownType reports whether an annotation names the type a literal takes on
// its own: STRING for a string, INT8 for an integer, DECIMAL for a number
// with a fraction, BOOL for a boolean.
func ownType(literal tree.Expr, ref tree.ResolvableTypeReference) bool {
	t, ok := tree.GetStaticallyKnownType(ref)
	if !ok {
		return false
	}
	switch l := literal.(type) {
	case *tree.StrVal:
		return t.Identical(types.String)
	case *tree.NumVal:
		if strings.ContainsAny(l.OrigString(), ".eE") {
			return t.Identical(types.Decimal)
		}
		return t.Identical(types.Int)
	case *tree.DBool:
		return t.Identical(types.Bool)
	}
	return false
}

func unparen(expr tree.Expr) tree.Expr {
	for {
		paren, ok := expr.(*tree.ParenExpr)
		if !ok {
			return expr
		}
		expr = paren.Expr
	}
}

// literalCastStripper knows the types of the columns an expression reads.
type literalCastStripper struct{ columns map[string]string }

// VisitPre drops a cast or annotation on a literal only where the context
// fixes the literal's type anyway, which is where CockroachDB adds them to
// what it stores: a side of a comparison whose other side is not a literal
// (s > 'a'::STRING, an IN list's elements), where the cast names the type of
// the column on the other side, and the sequence's REGCLASS in
// nextval('s'::REGCLASS), currval and setval. Anywhere else a cast chooses a type, of an
// overload, an operator, a COALESCE or a comparison of literals, and stays.
func (v literalCastStripper) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	switch e := expr.(type) {
	case *tree.ComparisonExpr:
		left, right := literalOperand(e.Left), literalOperand(e.Right)
		if left == right {
			return true, expr
		}
		cmp := *e
		if left {
			cmp.Left = stripTypedCast(cmp.Left, v.columnType(cmp.Right))
		} else {
			cmp.Right = stripTypedCast(cmp.Right, v.columnType(cmp.Left))
		}
		return true, &cmp
	case *tree.FuncExpr:
		// nextval, currval and setval take their sequence as a REGCLASS,
		// which CockroachDB writes on the name it stores.
		name, ok := e.Func.FunctionReference.(*tree.UnresolvedName)
		if !ok || len(e.Exprs) == 0 {
			break
		}
		switch strings.ToLower(name.Parts[0]) {
		case "nextval", "currval", "setval":
		default:
			return true, expr
		}
		arg := unparen(e.Exprs[0])
		if cast, ok := arg.(*tree.CastExpr); ok && isLiteral(unparen(cast.Expr)) && strings.EqualFold(cast.Type.SQLString(), "REGCLASS") {
			arg = unparen(cast.Expr)
		}
		literal, ok := arg.(*tree.StrVal)
		if !ok {
			break
		}
		// CockroachDB stores the name it resolved, in its schema:
		// nextval('s') as nextval('public.s'::REGCLASS).
		if sequence, err := crparser.ParseQualifiedTableName(literal.RawString()); err == nil {
			literal = tree.NewStrVal(schemaOf(sequence) + "." + string(sequence.ObjectName))
		}
		call := *e
		call.Exprs = append(tree.Exprs{literal}, e.Exprs[1:]...)
		return true, &call
	}
	return true, expr
}

func (v literalCastStripper) VisitPost(expr tree.Expr) tree.Expr {
	switch e := expr.(type) {
	case *tree.ParenExpr:
		return e.Expr
	}
	return expr
}

// columnType is the type of an expression that is a column the stripper
// knows, "" otherwise.
func (v literalCastStripper) columnType(expr tree.Expr) string {
	switch e := unparen(expr).(type) {
	case *tree.UnresolvedName:
		if !e.Star {
			return v.columns[e.Parts[0]]
		}
	case *tree.ColumnItem:
		return v.columns[string(e.ColumnName)]
	}
	return ""
}

// castRounds reports whether a cast can change its literal's value: a cast
// to a type with a width, length, precision or scale (VARCHAR(1),
// DECIMAL(3,1), TIMESTAMP(0)), to FLOAT4, to FLOAT8 of a number it cannot
// hold exactly, or of a number with a fraction to an integer. CockroachDB adds no such cast to what it stores, so it is kept.
// trimZeros drops a number's trailing fractional zeros.
func trimZeros(expr tree.Expr) tree.Expr {
	n, ok := expr.(*tree.NumVal)
	if !ok {
		return expr
	}
	orig := n.OrigString()
	if !strings.Contains(orig, ".") || strings.ContainsAny(orig, "eE") {
		return expr
	}
	trimmed := strings.TrimSuffix(strings.TrimRight(orig, "0"), ".")
	if trimmed == orig {
		return expr
	}
	return tree.NewNumVal(constant.MakeFromLiteral(trimmed, token.FLOAT, 0), trimmed, false)
}

func castRounds(cast *tree.CastExpr) bool {
	t, ok := tree.GetStaticallyKnownType(cast.Type)
	if !ok {
		return true
	}
	if strings.Contains(t.SQLString(), "(") {
		return true
	}
	switch t.Family() {
	case types.FloatFamily:
		return t.Width() != 64 || !exactFloat(cast.Expr)
	case types.IntFamily:
		fraction := false
		tree.WalkExprConst(&numVisitor{visit: func(n *tree.NumVal) {
			fraction = fraction || strings.ContainsAny(n.OrigString(), ".eE")
		}}, cast.Expr)
		return fraction
	}
	return false
}

// exactFloat reports whether every number in a literal is a float64 exactly.
func exactFloat(expr tree.Expr) bool {
	exact := true
	tree.WalkExprConst(&numVisitor{visit: func(n *tree.NumVal) {
		r, ok := new(big.Rat).SetString(n.OrigString())
		if !ok {
			exact = false
			return
		}
		f, _ := r.Float64()
		exact = exact && new(big.Rat).SetFloat64(f) != nil && new(big.Rat).SetFloat64(f).Cmp(r) == 0
	}}, expr)
	return exact
}

type numVisitor struct{ visit func(*tree.NumVal) }

func (v *numVisitor) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	if n, ok := expr.(*tree.NumVal); ok {
		v.visit(n)
	}
	return true, expr
}

func (*numVisitor) VisitPost(expr tree.Expr) tree.Expr { return expr }

// groupedKey formats a node with every expression in parentheses. The keys
// drop the parentheses an expression was written with, and the tree keeps
// its grouping, so (1 + 2) * 3 and 1 + 2 * 3 stay apart.
func groupedKey(node tree.NodeFormatter) string {
	return tree.AsStringWithFlags(node, tree.FmtAlwaysGroupExprs)
}

// literalOperand reports whether an operand is a literal, cast or annotated
// or not, or a tuple or arithmetic of them.
func literalOperand(expr tree.Expr) bool {
	switch e := unparen(expr).(type) {
	case *tree.BinaryExpr:
		return literalOperand(e.Left) && literalOperand(e.Right)
	case *tree.UnaryExpr:
		return literalOperand(e.Expr)
	case *tree.CastExpr:
		return literalOperand(e.Expr)
	case *tree.AnnotateTypeExpr:
		return literalOperand(e.Expr)
	case *tree.Tuple:
		for _, element := range e.Exprs {
			if !literalOperand(element) {
				return false
			}
		}
		return true
	}
	return isLiteral(unparen(expr))
}

func isLiteral(expr tree.Expr) bool {
	switch e := expr.(type) {
	case *tree.StrVal, *tree.NumVal, tree.Datum:
		return true
	case *tree.Array:
		for _, element := range e.Exprs {
			if !isLiteral(element) {
				return false
			}
		}
		return true
	}
	return false
}

// queryKey is a view's query without the database's name and with each
// relation it names by its schema, an unqualified one being the one in
// public, which is what CockroachDB's default search path finds.
func queryKey(query, database string) string {
	query = relationQuery(query, database)
	stmt, err := parseOne(query)
	if err != nil {
		return query
	}
	refs := cteReferences(stmt)
	i := 0
	ctx := tree.NewFmtCtx(tree.FmtSimple, tree.FmtReformatTableNames(func(ctx *tree.FmtCtx, tn *tree.TableName) {
		name := *tn
		if !name.ExplicitSchema && (i >= len(refs) || !refs[i]) {
			name.SchemaName, name.ExplicitSchema = "public", true
		}
		i++
		ctx.WithReformatTableNames(nil, func() { ctx.FormatNode(&name) })
	}))
	defer ctx.Close()
	ctx.FormatNode(stmt)
	return ctx.String()
}

// indexKey is an index's statement with its predicate as exprKey states it.
func indexKey(schema, table string, index *metadataIndex) string {
	statement := indexStatement(schema, table, index)
	parsed, ok := parseIndex(statement)
	if !ok || parsed.Predicate == nil {
		return statement
	}
	parsed.Predicate, _ = tree.WalkExpr(literalCastStripper{}, parsed.Predicate)
	return groupedKey(parsed)
}

// routineKey is a routine definition without the options CockroachDB states
// by default (VOLATILE, NOT LEAKPROOF, CALLED ON NULL INPUT, SECURITY
// INVOKER), its options in one order, and its body as CockroachDB formats it.
func routineKey(definition string) string {
	stmt, err := parseOne(strings.TrimSuffix(strings.TrimSpace(definition), ";"))
	if err != nil {
		return definition
	}
	routine, ok := stmt.(*tree.CreateRoutine)
	if !ok {
		return definition
	}
	routine.Replace = false
	// An unqualified routine is created in public; the database's name is not
	// part of it.
	if !routine.Name.ExplicitSchema {
		routine.Name.SchemaName, routine.Name.ExplicitSchema = "public", true
	}
	routine.Name.CatalogName, routine.Name.ExplicitCatalog = "", false
	// A parameter default is stored annotated (0:::INT8).
	for i := range routine.Params {
		if routine.Params[i].DefaultVal != nil {
			// Its parameter types it, as a column types its default.
			routine.Params[i].DefaultVal, _ = tree.WalkExpr(literalCastStripper{}, stripDefaultCast(routine.Params[i].DefaultVal, routine.Params[i].Type.SQLString()))
		}
	}
	var language tree.RoutineLanguage
	var options [5]tree.RoutineOption
	for _, option := range routine.Options {
		switch o := option.(type) {
		case tree.RoutineLanguage:
			language = tree.RoutineLanguage(strings.ToLower(string(o)))
			options[0] = language
		case tree.RoutineVolatility:
			if o != tree.RoutineVolatile {
				options[1] = o
			}
		case tree.RoutineNullInputBehavior:
			if o != tree.RoutineCalledOnNullInput {
				options[2] = o
			}
		case tree.RoutineLeakproof:
			if o {
				options[3] = o
			}
		case tree.RoutineSecurity:
			if o != tree.RoutineInvoker {
				options[4] = o
			}
		}
	}
	var normalized tree.RoutineOptions
	for _, option := range options {
		if option != nil {
			normalized = append(normalized, option)
		}
	}
	for _, option := range routine.Options {
		if body, ok := option.(tree.RoutineBodyStr); ok {
			normalized = append(normalized, tree.RoutineBodyStr(bodyKey(string(body), language)))
		}
	}
	routine.Options = normalized
	return groupedKey(routine)
}

func bodyKey(body string, language tree.RoutineLanguage) string {
	switch language {
	case "sql":
		stmts, err := crparser.Parse(body)
		if err != nil {
			return body
		}
		parts := make([]string, len(stmts))
		for i, s := range stmts {
			parts[i] = tree.AsString(s.AST)
		}
		return withoutAnnotations(strings.Join(parts, "; "))
	case "plpgsql":
		block, err := plpgsqlparser.Parse(body)
		if err != nil {
			return body
		}
		return tree.AsString(block.AST)
	}
	return body
}

// withoutAnnotations returns a SQL text's tokens without the :::type
// CockroachDB adds to the literals of a stored body (1:::INT8), as a key: an
// annotation inside a string or an identifier is the token's own text and
// stays.
func withoutAnnotations(sql string) string {
	tokens, ok := crparser.Tokens(sql)
	if !ok {
		return sql
	}
	var b strings.Builder
	// arrays holds, for each open bracket, the one type its elements were
	// annotated with ("" for none yet, "?" for several or a kept one), and
	// closed is the last closed bracket's.
	var arrays []string
	closed := "?"
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.TokenID {
		case tokenID("["):
			arrays = append(arrays, "")
		case tokenID("]"):
			if len(arrays) > 0 {
				closed = arrays[len(arrays)-1]
				arrays = arrays[:len(arrays)-1]
			}
		case annotateToken:
			if i == 0 || i+1 >= len(tokens) {
				break
			}
			typ := strings.ToLower(tokens[i+1].Str)
			next := func(k int) int32 {
				if i+k < len(tokens) {
					return tokens[i+k].TokenID
				}
				return 0
			}
			// A literal's own type, the one CockroachDB annotates it with
			// in a stored body, is no part of it; another type chooses an
			// overload ('é':::BYTES) and stays.
			if defaultAnnotation(tokens[i-1].TokenID) == typ && next(2) != tokenID("(") && next(2) != tokenID("[") {
				if len(arrays) > 0 {
					if top := arrays[len(arrays)-1]; top == "" || top == typ {
						arrays[len(arrays)-1] = typ
					} else {
						arrays[len(arrays)-1] = "?"
					}
				}
				i++
				continue
			}
			// So is an array literal's, ARRAY[1:::INT8, 2:::INT8]:::INT8[].
			if tokens[i-1].TokenID == tokenID("]") && closed != "" && closed != "?" && closed == typ &&
				next(2) == tokenID("[") && next(3) == tokenID("]") && next(4) != tokenID("[") {
				i += 3
				continue
			}
			if len(arrays) > 0 {
				arrays[len(arrays)-1] = "?"
			}
		}
		fmt.Fprintf(&b, "%d:%s\x00", t.TokenID, t.Str)
	}
	return b.String()
}

// defaultAnnotation is the type a literal token takes on its own: STRING for
// a string, INT8 for an integer, DECIMAL for a number with a fraction, BOOL
// for TRUE and FALSE.
func defaultAnnotation(literal int32) string {
	switch literal {
	case tokenID("'x'"):
		return "string"
	case tokenID("1"):
		return "int8"
	case tokenID("1.5"):
		return "decimal"
	case tokenID("true"), tokenID("false"):
		return "bool"
	}
	return ""
}

// annotateToken is the lexer's token for :::.
var annotateToken = func() int32 {
	tokens, _ := crparser.Tokens("1:::INT8")
	return tokens[1].TokenID
}()

// identifiers returns the words of a query, a superset of the columns it
// reads, and whether it may read columns it does not name: through a star,
// a NATURAL join's shared columns, or a relation's column aliases (FROM t AS
// x (a, b), whose a and b name t's columns by position).
func identifiers(query string) (map[string]bool, bool) {
	tokens, ok := crparser.Tokens(query)
	if !ok {
		return nil, true
	}
	words := make(map[string]bool)
	star := false
	ident := tokenID("x")
	for i, token := range tokens {
		words[token.Str] = true
		if token.Str == "*" || strings.EqualFold(token.Str, "natural") {
			star = true
		}
		// An alias with a column list: AS x ( or t x (.
		if i+2 < len(tokens) && tokens[i+1].TokenID == ident && tokens[i+2].TokenID == tokenID("(") &&
			(token.TokenID == tokenID("as") || token.TokenID == ident) {
			star = true
		}
	}
	return words, star
}
