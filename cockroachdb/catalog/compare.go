package catalog

import (
	"go/constant"
	"go/token"
	"regexp"
	"strings"

	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	plpgsqlparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
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
	parsed, _ = tree.WalkExpr(literalCastStripper{}, parsed)
	return sdlExpr(parsed)
}

type literalCastStripper struct{}

func (literalCastStripper) VisitPre(expr tree.Expr) (bool, tree.Expr) { return true, expr }

func (literalCastStripper) VisitPost(expr tree.Expr) tree.Expr {
	switch e := expr.(type) {
	case *tree.CastExpr:
		if isLiteral(e.Expr) {
			return e.Expr
		}
	case *tree.AnnotateTypeExpr:
		if isLiteral(e.Expr) {
			return e.Expr
		}
	case *tree.CollateExpr:
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

// indexKey is an index's statement with its predicate as exprKey states it.
func indexKey(schema, table string, index *metadataIndex) string {
	statement := indexStatement(schema, table, index)
	parsed, ok := parseIndex(statement)
	if !ok || parsed.Predicate == nil {
		return statement
	}
	parsed.Predicate, _ = tree.WalkExpr(literalCastStripper{}, parsed.Predicate)
	return tree.AsString(parsed)
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
	return tree.AsString(routine)
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
		// CockroachDB stores a body's literals annotated (1:::INT8).
		return typeAnnotation.ReplaceAllString(strings.Join(parts, "; "), "")
	case "plpgsql":
		block, err := plpgsqlparser.Parse(body)
		if err != nil {
			return body
		}
		return tree.AsString(block.AST)
	}
	return body
}

// typeAnnotation matches the :::type CockroachDB adds to a stored literal.
var typeAnnotation = regexp.MustCompile(`:::[A-Za-z0-9_]+(\[\])*`)

// identifiers returns the words of a query, a superset of the columns it
// reads, and whether it reads every column through a star.
func identifiers(query string) (map[string]bool, bool) {
	tokens, ok := crparser.Tokens(query)
	if !ok {
		return nil, true
	}
	words := make(map[string]bool)
	star := false
	for _, token := range tokens {
		words[token.Str] = true
		if token.Str == "*" {
			star = true
		}
	}
	return words, star
}
