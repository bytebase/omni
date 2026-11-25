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
)

// The keys below compare two definitions of one object, one synced and one
// written by hand, as CockroachDB would store them. They are for comparison
// only; what a migration writes is the target's own text.

// exprKey is an expression without what CockroachDB adds to the literals it
// stores (”::STRING, 's.seq'::REGCLASS, a collated column's COLLATE) or
// takes off them (ARRAY[]::STRING[]), and without parentheses.
func exprKey(expr string) string {
	if strings.TrimSpace(expr) == "" {
		return ""
	}
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return expr
	}
	// CockroachDB adds COLLATE to a collated column's literal default as a
	// whole; inside an expression a collation is the expression's own.
	if collate, ok := unparen(parsed).(*tree.CollateExpr); ok && isLiteral(collate.Expr) {
		parsed = collate.Expr
	}
	// The whole expression is a column's default, whose type the column
	// fixes.
	parsed = stripLiteralCast(parsed)
	parsed, _ = tree.WalkExpr(literalCastStripper{}, parsed)
	return groupedKey(parsed)
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

type literalCastStripper struct{}

// VisitPre drops a cast or annotation on a literal only where the context
// fixes the literal's type anyway, which is where CockroachDB adds them to
// what it stores: a side of a comparison whose other side is not a literal
// (s > 'a'::STRING, an IN list's elements), and the sequence's REGCLASS in
// nextval('s'::REGCLASS), currval and setval. Anywhere else a cast chooses a type, of an
// overload, an operator, a COALESCE or a comparison of literals, and stays.
func (literalCastStripper) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	switch e := expr.(type) {
	case *tree.ComparisonExpr:
		left, right := literalOperand(e.Left), literalOperand(e.Right)
		if left == right {
			return true, expr
		}
		cmp := *e
		if left {
			cmp.Left = stripLiteralCast(cmp.Left)
		} else {
			cmp.Right = stripLiteralCast(cmp.Right)
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
		cast, ok := unparen(e.Exprs[0]).(*tree.CastExpr)
		if !ok || !isLiteral(unparen(cast.Expr)) || !strings.EqualFold(cast.Type.SQLString(), "REGCLASS") {
			break
		}
		call := *e
		call.Exprs = append(tree.Exprs{unparen(cast.Expr)}, e.Exprs[1:]...)
		return true, &call
	}
	return true, expr
}

func (literalCastStripper) VisitPost(expr tree.Expr) tree.Expr {
	switch e := expr.(type) {
	case *tree.ParenExpr:
		return e.Expr
	}
	return expr
}

// stripLiteralCast drops a cast or annotation on a literal, or on each of a
// tuple's elements, unless the cast can change the literal's value. In the
// same contexts a number's trailing fractional zeros go: CockroachDB states a
// FLOAT column's default or comparand as a float (0 as 0.0).
func stripLiteralCast(expr tree.Expr) tree.Expr {
	switch e := unparen(expr).(type) {
	case *tree.CastExpr:
		if isLiteral(unparen(e.Expr)) && !castRounds(e) {
			return trimZeros(unparen(e.Expr))
		}
	case *tree.AnnotateTypeExpr:
		if isLiteral(unparen(e.Expr)) {
			return trimZeros(unparen(e.Expr))
		}
	case *tree.NumVal:
		return trimZeros(e)
	case *tree.Tuple:
		tuple := *e
		tuple.Exprs = make(tree.Exprs, len(e.Exprs))
		for i, element := range e.Exprs {
			tuple.Exprs[i] = stripLiteralCast(element)
		}
		return &tuple
	}
	return expr
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
			routine.Params[i].DefaultVal, _ = tree.WalkExpr(literalCastStripper{}, stripLiteralCast(routine.Params[i].DefaultVal))
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
