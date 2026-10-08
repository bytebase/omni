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

// explainAnalyzes reports whether the EXPLAIN runs its query: ANALYZE is
// on and every option is one the server accepts, since an option it
// rejects fails the statement before anything runs. The options apply in
// order, so a repeated one takes its last value. Boolean options read as
// defGetBoolean does; FORMAT and SERIALIZE take their listed values;
// ANALYZE and GENERIC_PLAN together are refused.
//
// pg: src/backend/commands/explain.c — ExplainQuery
func explainAnalyzes(e *ast.ExplainStmt) bool {
	if e.Options == nil {
		return false
	}
	analyze, genericPlan := false, false
	for _, item := range e.Options.Items {
		opt, ok := item.(*ast.DefElem)
		if !ok {
			return false
		}
		switch opt.Defname {
		case "analyze":
			value, ok := optionBool(opt.Arg)
			if !ok {
				return false
			}
			analyze = value
		case "generic_plan":
			value, ok := optionBool(opt.Arg)
			if !ok {
				return false
			}
			genericPlan = value
		case "verbose", "costs", "settings", "buffers", "wal", "timing", "summary", "memory":
			if _, ok := optionBool(opt.Arg); !ok {
				return false
			}
		case "format":
			if !optionIn(opt.Arg, "text", "xml", "json", "yaml") {
				return false
			}
		case "serialize":
			if opt.Arg != nil && !optionIn(opt.Arg, "off", "none", "text", "binary") {
				return false
			}
		default:
			return false
		}
	}
	return analyze && !genericPlan
}

// optionBool reads an option value as defGetBoolean does: no value is
// true, an integer must be 0 or 1, a string is parsed as a boolean.
//
// pg: src/backend/commands/define.c — defGetBoolean
func optionBool(arg ast.Node) (value, ok bool) {
	switch v := arg.(type) {
	case nil:
		return true, true
	case *ast.Integer:
		return v.Ival == 1, v.Ival == 0 || v.Ival == 1
	case *ast.Boolean:
		return v.Boolval, true
	case *ast.String:
		return parseBool(v.Str)
	}
	return false, false
}

// optionIn reports whether the option value is one of the listed words,
// compared exactly as the server does.
func optionIn(arg ast.Node, words ...string) bool {
	s, ok := arg.(*ast.String)
	if !ok {
		return false
	}
	for _, w := range words {
		if s.Str == w {
			return true
		}
	}
	return false
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
