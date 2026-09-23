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

// explainAnalyzes reports whether the EXPLAIN carries ANALYZE, with a
// value that is on: ANALYZE alone, or true/on/yes/1.
//
// pg: src/backend/commands/define.c — defGetBoolean
func explainAnalyzes(e *ast.ExplainStmt) bool {
	if e.Options == nil {
		return false
	}
	for _, item := range e.Options.Items {
		opt, ok := item.(*ast.DefElem)
		if !ok || opt.Defname != "analyze" {
			continue
		}
		switch v := opt.Arg.(type) {
		case nil:
			return true
		case *ast.Integer:
			return v.Ival != 0
		case *ast.Boolean:
			return v.Boolval
		case *ast.String:
			switch strings.ToLower(v.Str) {
			case "true", "on", "yes", "1":
				return true
			}
			return false
		}
	}
	return false
}
