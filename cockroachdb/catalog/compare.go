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

// VisitPre keeps a cast of a cast whole: '01'::INT8::STRING is '1', and
// stripping its inner cast would make it '01'::STRING.
func (literalCastStripper) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	if cast, ok := expr.(*tree.CastExpr); ok {
		if _, ok := unparen(cast.Expr).(*tree.CastExpr); ok {
			return false, expr
		}
	}
	return true, expr
}

func (literalCastStripper) VisitPost(expr tree.Expr) tree.Expr {
	switch e := expr.(type) {
	case *tree.CastExpr:
		if isLiteral(e.Expr) && !castRounds(e) {
			return e.Expr
		}
	case *tree.AnnotateTypeExpr:
		if isLiteral(e.Expr) {
			return e.Expr
		}
	case *tree.ParenExpr:
		return e.Expr
	case *tree.NumVal:
		// CockroachDB states a number compared with a FLOAT column as a
		// float (0 as 0.0), so a decimal point's trailing zeros go.
		orig := e.OrigString()
		if strings.Contains(orig, ".") && !strings.ContainsAny(orig, "eE") {
			trimmed := strings.TrimSuffix(strings.TrimRight(orig, "0"), ".")
			if trimmed != orig {
				return tree.NewNumVal(constant.MakeFromLiteral(trimmed, token.FLOAT, 0), trimmed, false)
			}
		}
	}
	return expr
}

// castRounds reports whether a cast can change its literal's value: a cast
// to a type with a width, length, precision or scale (VARCHAR(1),
// DECIMAL(3,1), TIMESTAMP(0)), to FLOAT4, to FLOAT8 of a number it cannot
// hold exactly, or of a number with a fraction to an integer. CockroachDB adds no such cast to what it stores, so it is kept.
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
			routine.Params[i].DefaultVal, _ = tree.WalkExpr(literalCastStripper{}, routine.Params[i].DefaultVal)
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
	for i := 0; i < len(tokens); i++ {
		if tokens[i].TokenID == annotateToken {
			// Skip the type's name, its modifiers and its array brackets.
			i++
			if i+1 < len(tokens) && tokens[i+1].Str == "(" {
				for i+1 < len(tokens) && tokens[i].Str != ")" {
					i++
				}
			}
			for i+2 < len(tokens) && tokens[i+1].Str == "[" && tokens[i+2].Str == "]" {
				i += 2
			}
			continue
		}
		fmt.Fprintf(&b, "%d:%s\x00", tokens[i].TokenID, tokens[i].Str)
	}
	return b.String()
}

// annotateToken is the lexer's token for :::.
var annotateToken = func() int32 {
	tokens, _ := crparser.Tokens("1:::INT8")
	return tokens[1].TokenID
}()

// identifiers returns the words of a query, a superset of the columns it
// reads, and whether it may read columns it does not name: through a star,
// or through a NATURAL join's shared columns.
func identifiers(query string) (map[string]bool, bool) {
	tokens, ok := crparser.Tokens(query)
	if !ok {
		return nil, true
	}
	words := make(map[string]bool)
	star := false
	for _, token := range tokens {
		words[token.Str] = true
		if token.Str == "*" || strings.EqualFold(token.Str, "natural") {
			star = true
		}
	}
	return words, star
}
