package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkRequireWhere reports an UPDATE or DELETE without a WHERE clause,
// wherever it sits: top level, a data-modifying CTE, or a SQL-standard
// function body. A plain EXPLAIN only plans its query, so the DML under
// it runs nothing and is skipped; EXPLAIN ANALYZE runs it and is not.
// MERGE is not in scope.
func checkRequireWhere(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt:
			return explainAnalyzes(v)
		case *ast.UpdateStmt:
			if v.WhereClause == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "UPDATE "+relation(v.Relation)+" has no WHERE clause")
			}
		case *ast.DeleteStmt:
			if v.WhereClause == nil {
				r.report(review.RequireWhere, s.index, rangeOf(v.Loc), "DELETE FROM "+relation(v.Relation)+" has no WHERE clause")
			}
		}
		return true
	})
}

// explainAnalyzes reports whether the EXPLAIN carries ANALYZE with a value
// that is on: ANALYZE alone, or a value the server reads as true. The
// options apply in order, so a repeated ANALYZE takes its last value. A
// value the server would reject reads as off.
//
// pg: src/backend/commands/explain.c — ExplainQuery; define.c — defGetBoolean
func explainAnalyzes(e *ast.ExplainStmt) bool {
	on := false
	if e.Options == nil {
		return on
	}
	for _, item := range e.Options.Items {
		opt, ok := item.(*ast.DefElem)
		if !ok || opt.Defname != "analyze" {
			continue
		}
		switch v := opt.Arg.(type) {
		case nil:
			on = true
		case *ast.Integer:
			on = v.Ival != 0
		case *ast.Boolean:
			on = v.Boolval
		case *ast.String:
			on, _ = parseBool(v.Str)
		default:
			on = false
		}
	}
	return on
}

// parseBool reads a boolean the way the server's option parser does: any
// prefix of true, false, yes, or no, a prefix of at least two letters of
// on or off, or 1 or 0, ignoring case. ok is false for anything else.
//
// pg: src/backend/utils/adt/bool.c — parse_bool_with_len
func parseBool(s string) (value, ok bool) {
	s = strings.ToLower(s)
	if s == "" {
		return false, false
	}
	for _, word := range []string{"true", "yes", "on"} {
		if strings.HasPrefix(word, s) && (word != "on" || len(s) >= 2) {
			return true, true
		}
	}
	for _, word := range []string{"false", "no", "off"} {
		if strings.HasPrefix(word, s) && (word != "off" || len(s) >= 2) {
			return false, true
		}
	}
	switch s {
	case "1":
		return true, true
	case "0":
		return false, true
	}
	return false, false
}
