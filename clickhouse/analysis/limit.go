package analysis

import (
	"fmt"
	"strings"
	"unicode"
)

// StatementWithResultLimit bounds a subqueryable ClickHouse statement with a
// CTE wrapper (`WITH result AS (...) SELECT * FROM result LIMIT N;`).
// Statements that cannot stand as a subquery (SHOW, DESCRIBE, INSERT, EXPLAIN,
// etc.) are returned unchanged.
func StatementWithResultLimit(statement string, limit int) string {
	wrapped := trimStatementForWrapper(statement)
	if !isSubqueryable(wrapped) {
		return statement
	}
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", wrapped, limit)
}

// WrapResultLimit trims whitespace, trailing semicolons, and trailing comments
// from statement and unconditionally wraps it in ClickHouse's result-limit CTE.
func WrapResultLimit(statement string, limit int) string {
	wrapped := trimStatementForWrapper(statement)
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", wrapped, limit)
}

// trimStatementForWrapper removes trailing comments and the terminators before
// them, which a CTE cannot hold.
func trimStatementForWrapper(statement string) string {
	trimmed := trimStatement(statement)
	for {
		start := trailingCommentStart(trimmed)
		if start < 0 {
			return trimmed
		}
		trimmed = trimStatement(strings.TrimSpace(string([]rune(trimmed)[:start])))
	}
}

// isSubqueryable reports whether statement can stand inside a CTE.
func isSubqueryable(statement string) bool {
	body := statementBody(statement)
	if strings.HasPrefix(body, "(") {
		return true
	}
	upper := strings.ToUpper(body)
	for _, keyword := range []string{"SELECT", "WITH", "TABLE", "VALUES"} {
		if startsWithKeyword(upper, keyword) {
			return true
		}
	}
	return false
}

// trailingCommentStart is the rune index where a comment that reaches the end
// of statement starts, or -1. It skips strings, quoted identifiers, heredocs,
// and earlier comments so a `--` or `#` inside one does not count.
func trailingCommentStart(statement string) int {
	runes := []rune(statement)
	last := -1
	line, block, str := false, false, false
	var ident rune
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case line:
			if c == '\n' {
				line = false
			}
		case block:
			if c == '*' && i+1 < len(runes) && runes[i+1] == '/' {
				block = false
				i++
			}
		case str:
			if c == '\\' {
				i++
			} else if c == '\'' {
				str = false
			}
		case ident != 0:
			if c == '\\' {
				i++
			} else if c == ident {
				ident = 0
			}
		default:
			switch {
			case c == '-' && i+1 < len(runes) && runes[i+1] == '-':
				line = true
				i++
				last = i - 1
			case c == '#':
				line = true
				last = i
			case c == '$' && isHeredocStart(runes, i):
				if err := scanHeredoc(runes, &i); err != nil {
					last = -1
					i = len(runes)
				} else {
					i--
					last = -1
				}
			case c == '/' && i+1 < len(runes) && runes[i+1] == '*':
				block = true
				i++
				last = i - 1
			case unicode.IsSpace(c):
			case c == '\'':
				str = true
				last = -1
			case c == '"' || c == '`':
				ident = c
				last = -1
			default:
				last = -1
			}
		}
	}
	return last
}

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
