package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
)

// explainAnalyzes reports whether an EXPLAIN runs the statement it carries,
// which only EXPLAIN ANALYZE without GENERIC_PLAN does: a plain EXPLAIN
// plans a CREATE TABLE AS or SELECT INTO without creating anything. An
// option the server rejects fails the statement before anything runs. The
// options apply in order, so a repeated one takes its last value.
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

// optionBool reads an EXPLAIN option's boolean value; a bare option is
// true. ok is false for a value the server rejects.
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
// on or off, or 1 or 0, ignoring case.
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
