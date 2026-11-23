package analysis

import "github.com/bytebase/omni/databricks/parser"

// ExplainStatement returns the Databricks EXPLAIN statement for statement.
// Databricks prints Spark's FORMATTED plan regardless of the requested format;
// a statement that already asks for a plan is returned unchanged.
func ExplainStatement(statement, _ string) string {
	if StartsWithExplain(statement) {
		return statement
	}
	return "EXPLAIN FORMATTED " + statement
}

// StartsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading whitespace and comments (including nested block comments).
// An unclosed block comment returns false so a write hidden behind a malformed
// comment is never mistaken for an EXPLAIN.
func StartsWithExplain(statement string) bool {
	l := parser.NewLexer(statement)
	tok := l.NextToken()
	return l.Err() == nil && tok.Kind == parser.KwExplain
}

// PlanFormat reports the plan format ("TEXT") of a statement that already
// starts with EXPLAIN. Databricks prints Spark's text plan for all EXPLAIN
// forms.
func PlanFormat(statement string) (string, bool) {
	if !StartsWithExplain(statement) {
		return "", false
	}
	return "TEXT", true
}

// PlanExecuted reports whether statement is an EXPLAIN ANALYZE, whose output
// is a plan of an execution that happened. `EXPLAIN ANALYZE TABLE ...` and
// `EXPLAIN ANALYZE TABLES ...` plan an ANALYZE command without executing it
// and return false.
func PlanExecuted(statement string) bool {
	_, executed, ok := DescribeExplain(statement)
	return ok && executed
}

// DescribeExplain returns the plan format ("TEXT") and whether statement
// executes what it plans (EXPLAIN ANALYZE). ok is false when statement does
// not start with EXPLAIN.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	l := parser.NewLexer(statement)
	first := l.NextToken()
	if l.Err() != nil || first.Kind != parser.KwExplain {
		return "", false, false
	}
	second := l.NextToken()
	if l.Err() != nil {
		return "TEXT", false, true
	}
	if second.Kind == parser.KwAnalyze {
		third := l.NextToken()
		if l.Err() == nil && third.Kind != parser.KwTable && third.Kind != parser.KwTables {
			return "TEXT", true, true
		}
	}
	return "TEXT", false, true
}

// explainAnalyzeBody returns the statement substring after `EXPLAIN ANALYZE`
// and its byte offset in statement when statement is an EXPLAIN ANALYZE.
func explainAnalyzeBody(statement string) (body string, offset int, ok bool) {
	l := parser.NewLexer(statement)
	first := l.NextToken()
	if l.Err() != nil || first.Kind != parser.KwExplain {
		return "", 0, false
	}
	second := l.NextToken()
	if l.Err() != nil || second.Kind != parser.KwAnalyze {
		return "", 0, false
	}
	third := l.NextToken()
	if l.Err() != nil || third.Kind == parser.TokEOF || third.Kind == parser.KwTable || third.Kind == parser.KwTables {
		return "", 0, false
	}
	return statement[third.Start:], third.Start, true
}
