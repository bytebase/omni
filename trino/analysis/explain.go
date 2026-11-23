package analysis

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// ErrExplainAnalyzeExecutes is ExplainStatement's error for EXPLAIN ANALYZE,
// which executes the statement it plans.
var ErrExplainAnalyzeExecutes = analysis.ErrExplainAnalyzeExecutes

// ExplainStatement returns the statement that plans statement:
// EXPLAIN (FORMAT JSON) when json is set, otherwise EXPLAIN, the engine's
// default text plan. An EXPLAIN is rebuilt around the statement it plans;
// one whose inner statement cannot be located is returned as written.
// EXPLAIN ANALYZE returns ErrExplainAnalyzeExecutes, and a statement that
// does not parse returns the parser's first error.
func ExplainStatement(statement string, json bool) (string, error) {
	node, planned, err := singleStatement(statement)
	if err != nil {
		return "", err
	}
	if explain, ok := node.(*parser.ExplainStmt); ok {
		if explain.Analyze {
			return "", ErrExplainAnalyzeExecutes
		}
		if inner, ok := explainedText(explain, statement); ok {
			planned = inner
		} else {
			// Nothing to rebuild around: run what the user wrote rather than
			// plan a plan.
			return planned, nil
		}
	}
	if json {
		return "EXPLAIN (FORMAT JSON) " + planned, nil
	}
	return "EXPLAIN " + planned, nil
}

// DescribeExplain returns the plan format statement asks for ("TEXT", "JSON", or
// "GRAPHVIZ") and whether producing it executes the statement it plans
// (EXPLAIN ANALYZE). ok is false when statement is not a single EXPLAIN or is
// EXPLAIN (TYPE VALIDATE), which returns a boolean validation result rather
// than a query plan.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return "", false, false
	}
	explain, ok := file.Stmts[0].(*parser.ExplainStmt)
	if !ok {
		return "", false, false
	}
	if explain.Type == parser.ExplainTypeValidate {
		return "", false, false
	}
	format = "TEXT"
	switch explain.Format {
	case parser.ExplainFormatJSON:
		format = "JSON"
	case parser.ExplainFormatGraphviz:
		format = "GRAPHVIZ"
	case parser.ExplainFormatUnset:
		// The IO plan is a JSON document unless a FORMAT option says otherwise.
		if explain.Type == parser.ExplainTypeIO {
			format = "JSON"
		}
	}
	return format, explain.Analyze, true
}

// singleStatement returns the one statement in text and its isolated text
// without surrounding empty semicolon-delimited segments. An explain request
// plans one statement at a time: both the drivers and the API gate split a
// multi-statement request first.
func singleStatement(text string) (ast.Node, string, error) {
	file, err := parser.Parse(text)
	if err != nil {
		return nil, "", parser.FirstError(err)
	}
	if file == nil || len(file.Stmts) == 0 || file.Stmts[0] == nil {
		return nil, "", errors.New("EXPLAIN needs a statement to plan")
	}
	if len(file.Stmts) > 1 {
		return nil, "", errors.New("EXPLAIN plans a single statement")
	}
	planned := text
	for _, seg := range parser.Split(text) {
		if seg.Empty() {
			continue
		}
		// Keep a single trailing semicolon when present, while dropping any
		// leading empty segments (e.g. "; SELECT 1") so prefixing EXPLAIN
		// never produces "EXPLAIN ; SELECT 1".
		end := seg.ByteEnd
		if end < len(text) && text[end] == ';' {
			end++
		}
		planned = text[seg.ByteStart:end]
		if seg.ByteStart > 0 {
			planned = strings.TrimLeft(planned, " \t\r\n")
		}
		break
	}
	return file.Stmts[0], planned, nil
}

// explainedText returns the text within text of the statement explain plans,
// and whether it could be located.
func explainedText(explain *parser.ExplainStmt, text string) (string, bool) {
	inner, ok := nodeSpan(explain.Statement)
	if !ok {
		return "", false
	}
	// The plan ends where the EXPLAIN does, which is past clauses that the
	// planned statement's own location can leave out.
	if inner.Start < 0 || inner.Start >= explain.Loc.End || explain.Loc.End > len(text) {
		return "", false
	}
	return text[inner.Start:explain.Loc.End], true
}
