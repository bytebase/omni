package analysis

import (
	"strings"

	"github.com/bytebase/omni/clickhouse/parser"
)

// ExplainStatement returns the ClickHouse EXPLAIN statement for statement.
// When format is "JSON" (case-insensitive), it requests `json = 1, indexes = 1`
// so reads report part and granule index pruning. A statement that already
// begins with EXPLAIN is returned unchanged so a typed EXPLAIN keeps its kind.
func ExplainStatement(statement, format string) string {
	if StartsWithExplain(statement) {
		return statement
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN json = 1, indexes = 1 " + statement
	}
	return "EXPLAIN " + statement
}

// StartsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading whitespace and comments (including nested block comments).
// An unclosed block comment or quoted identifier returns false.
func StartsWithExplain(statement string) bool {
	return parser.NewLexer(statement).NextToken().Type == parser.KwExplain
}

// PlanFormat reports the plan format ("JSON" or "TEXT") of a statement that
// already starts with EXPLAIN.
func PlanFormat(statement string) (string, bool) {
	info := inspectExplain(statement)
	if !info.isExplain {
		return "", false
	}
	return info.format, true
}

// DescribeExplain returns the plan format ("JSON" or "TEXT") and whether
// statement executes what it plans (EXPLAIN ANALYZE). ok is false when
// statement does not start with EXPLAIN.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	info := inspectExplain(statement)
	if !info.isExplain {
		return "", false, false
	}
	return info.format, info.executed, true
}

type explainHeader struct {
	isExplain bool
	format    string
	executed  bool
}

// inspectExplain inspects a ClickHouse EXPLAIN statement via clickhouse/parser.
// When the full statement parses as an *parser.ExplainStmt, the AST directly
// determines both the plan format and execution flag. Otherwise (for example,
// partial/non-standard EXPLAIN headers such as `EXPLAIN JSON SELECT 1`), it
// scans the EXPLAIN header using clickhouse/parser's lexer so nested block
// comments, line comments (`--`, `#`), quoted identifiers, and inner query
// expressions (such as `EXPLAIN SELECT json = 1 FROM t`) are handled
// accurately.
func inspectExplain(statement string) explainHeader {
	l := parser.NewLexer(statement)
	if l.NextToken().Type != parser.KwExplain {
		return explainHeader{}
	}

	if stmts, err := parser.Parse(statement); err == nil && len(stmts) == 1 {
		if ex, ok := stmts[0].(*parser.ExplainStmt); ok {
			return explainHeaderFromAST(ex)
		}
	}

	return inspectExplainHeaderLex(l)
}

func explainHeaderFromAST(ex *parser.ExplainStmt) explainHeader {
	format := "TEXT"
	for _, s := range ex.Settings {
		if strings.EqualFold(s.Name, "JSON") {
			if settingExprIsOff(s.Value) {
				format = "TEXT"
			} else {
				format = "JSON"
			}
		}
	}
	return explainHeader{
		isExplain: true,
		format:    format,
		executed:  strings.EqualFold(ex.Kind, "ANALYZE"),
	}
}

func inspectExplainHeaderLex(l *parser.Lexer) explainHeader {
	info := explainHeader{
		isExplain: true,
		format:    "TEXT",
	}

	cur := l.NextToken()
	next := l.NextToken()
	advance := func() {
		cur = next
		next = l.NextToken()
	}

	// Optional EXPLAIN kind when not followed by `=` (which begins settings).
	if next.Type != parser.TokEq {
		switch cur.Type {
		case parser.KwAst, parser.KwSyntax, parser.KwPlan, parser.KwPipeline, parser.KwEstimate:
			advance()
		case parser.KwQuery:
			if next.Type == parser.KwTree {
				advance()
				advance()
			}
		case parser.KwTable:
			if next.Type == parser.KwOverride {
				advance()
				advance()
			}
		case parser.KwCurrent:
			if next.Type == parser.TokIdent && strings.EqualFold(next.Val, "TRANSACTION") {
				return info
			}
		case parser.TokIdent:
			switch {
			case strings.EqualFold(cur.Val, "ANALYZE"):
				info.executed = true
				advance()
			case strings.EqualFold(cur.Val, "JSON"):
				info.format = "JSON"
				return info
			}
		}
	}

	// Optional `name = value (, name = value)*` settings before the query body.
	for (cur.Type == parser.TokIdent || cur.IsKeyword()) && next.Type == parser.TokEq {
		name := cur.Val
		advance() // move cur to '='
		advance() // move cur to value (or sign)
		if cur.Type == parser.TokEOF || cur.Type == parser.TokIllegal {
			break
		}
		val := cur.Val
		if (cur.Type == parser.TokPlus || cur.Type == parser.TokMinus) &&
			(next.Type == parser.TokInt || next.Type == parser.TokFloat) {
			val = cur.Raw + next.Val
			advance()
		}
		if strings.EqualFold(name, "JSON") {
			if settingIsOff(val) {
				info.format = "TEXT"
			} else {
				info.format = "JSON"
			}
		}
		advance() // move past value
		if cur.Type != parser.TokComma {
			break
		}
		advance() // consume ','
	}
	return info
}

func settingExprIsOff(expr parser.Expr) bool {
	switch e := expr.(type) {
	case *parser.LiteralExpr:
		return settingIsOff(e.Val)
	case *parser.IdentExpr:
		if len(e.Parts) == 1 {
			return settingIsOff(e.Parts[0])
		}
	case *parser.UnaryExpr:
		if e.Op == "+" || e.Op == "-" {
			return settingExprIsOff(e.Expr)
		}
	}
	return false
}

func settingIsOff(value string) bool {
	switch strings.ToUpper(value) {
	case "0", "+0", "-0", "FALSE", "OFF", "NO":
		return true
	}
	return false
}
