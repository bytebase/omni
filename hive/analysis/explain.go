package analysis

import (
	"strings"

	"github.com/bytebase/omni/hive/parser"
)

// ExplainStatement returns the Hive EXPLAIN statement for statement:
// EXPLAIN FORMATTED when format is "JSON" (case-insensitive), which Hive
// answers as one JSON object, and EXPLAIN otherwise. A statement that already
// starts with EXPLAIN is returned unchanged.
func ExplainStatement(statement, format string) string {
	if _, ok := PlanFormat(statement); ok {
		return statement
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN FORMATTED " + statement
	}
	return "EXPLAIN " + statement
}

// StartsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading whitespace and comments (including nested block comments).
func StartsWithExplain(statement string) bool {
	_, ok := PlanFormat(statement)
	return ok
}

// PlanFormat reports the format ("JSON" or "TEXT") of a statement that is
// already an EXPLAIN: EXPLAIN FORMATTED answers with JSON, and every other
// EXPLAIN with a text tree.
func PlanFormat(statement string) (string, bool) {
	info := inspectExplain(statement)
	if !info.isExplain {
		return "", false
	}
	if info.formatted {
		return "JSON", true
	}
	return "TEXT", true
}

// PlanExecuted reports whether statement is an EXPLAIN ANALYZE, which executes
// the statement it plans rather than only compiling a plan. EXPLAIN ANALYZE
// TABLE ... plans an ANALYZE TABLE statement and does not execute it.
func PlanExecuted(statement string) bool {
	info := inspectExplain(statement)
	return info.isExplain && info.analyze
}

type explainHeader struct {
	isExplain bool
	formatted bool
	analyze   bool
}

// inspectExplain scans the EXPLAIN keyword and any Hive 4.0.1 explainOption*
// modifiers using hive/parser's lexer so nested/unclosed block comments,
// backtick identifiers, and inner statement keywords (such as DESCRIBE
// FORMATTED or ANALYZE TABLE) are handled accurately.
func inspectExplain(statement string) explainHeader {
	l := parser.NewLexer(statement)
	if l.NextToken().Type != parser.KW_EXPLAIN {
		return explainHeader{}
	}
	info := explainHeader{isExplain: true}
	cur := l.NextToken()
	if cur.Type == parser.KW_REWRITE {
		return info
	}
	for {
		switch cur.Type {
		case parser.KW_FORMATTED:
			info.formatted = true
			cur = l.NextToken()
		case parser.KW_EXTENDED, parser.KW_DEPENDENCY, parser.KW_LOGICAL,
			parser.KW_AUTHORIZATION, parser.KW_REOPTIMIZATION, parser.KW_LOCKS,
			parser.KW_AST, parser.KW_DEBUG, parser.KW_DDL:
			cur = l.NextToken()
		case parser.KW_ANALYZE:
			next := l.NextToken()
			// ANALYZE TABLE ... starts analyzeStatement rather than explainOption.
			if next.Type == parser.KW_TABLE {
				return info
			}
			info.analyze = true
			cur = next
		case parser.KW_CBO:
			cur = l.NextToken()
			if cur.Type == parser.KW_COST || cur.Type == parser.KW_JOINCOST {
				cur = l.NextToken()
			}
		case parser.KW_VECTORIZATION:
			cur = l.NextToken()
			if cur.Type == parser.KW_ONLY {
				cur = l.NextToken()
			}
			switch cur.Type {
			case parser.KW_SUMMARY, parser.KW_OPERATOR, parser.KW_EXPRESSION, parser.KW_DETAIL:
				cur = l.NextToken()
			}
		default:
			return info
		}
	}
}
