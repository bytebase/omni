package analysis

import (
	"errors"
	"fmt"
	"strings"
)

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns rows:
// (readOnly, returnsData, error).
//
// Decision is lexical: every non-empty statement in a script must be a
// read-only Databricks SQL query (SELECT, WITH ... SELECT/VALUES/TABLE, VALUES,
// TABLE, SHOW, DESCRIBE/DESC, or non-executing EXPLAIN over a read-only
// statement). EXPLAIN ANALYZE, EXPLAIN over writes, CTEs containing DML or
// feeding non-query statements, DML (INSERT, UPDATE, DELETE, MERGE, COPY INTO),
// DDL (CREATE, ALTER, DROP, TRUNCATE, OPTIMIZE, VACUUM, RESTORE, REFRESH,
// MSCK), DCL (GRANT, REVOKE, DENY), SET, USE, CALL, and multi-statement batches
// containing any non-read-only statement are rejected. Empty, whitespace-only,
// comment-only, or semicolon-only input carries no statements to refuse and
// returns (true, true, nil). Lexical errors (unterminated comments, string
// literals, or quoted identifiers) are returned so callers fail closed.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	stmts, err := tokenizeStatements(statement)
	if err != nil {
		return false, false, err
	}
	for _, stmt := range stmts {
		if !isReadOnlyStatement(stmt) {
			return false, false, nil
		}
	}
	return true, true, nil
}

type tokenKind uint8

const (
	tokWord tokenKind = iota + 1
	tokIdent
	tokLiteral
	tokLParen
	tokRParen
	tokComma
	tokPipe
	tokOther
)

type token struct {
	kind tokenKind
	val  string // uppercased text for tokWord; empty otherwise
}

// tokenizeStatements lexes Databricks SQL into per-statement token slices split
// on top-level semicolons. It handles `--` line comments, nested `/* /* */ */`
// block comments, single/double-quoted strings with backslash and doubled-quote
// escapes, raw strings (`r'...'`, `R'...'`, `r"..."`, `R"..."`), hex binary
// literals (`x'...'`, `X'...'`), and backtick-quoted identifiers.
func tokenizeStatements(script string) ([][]token, error) {
	var stmts [][]token
	var cur []token
	i := 0
	for i < len(script) {
		c := script[i]
		switch {
		case isSpace(c):
			i++
		case c == '-' && i+1 < len(script) && script[i+1] == '-':
			i += 2
			for i < len(script) && script[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < len(script) && script[i+1] == '*':
			next, err := skipNestedBlockComment(script, i)
			if err != nil {
				return nil, err
			}
			i = next
		case c == '\'' || c == '"':
			next, err := skipQuotedString(script, i, false)
			if err != nil {
				return nil, err
			}
			cur = append(cur, token{kind: tokLiteral})
			i = next
		case c == '`':
			next, err := skipBacktickIdent(script, i)
			if err != nil {
				return nil, err
			}
			cur = append(cur, token{kind: tokIdent})
			i = next
		case c == ';':
			if len(cur) > 0 {
				stmts = append(stmts, cur)
				cur = nil
			}
			i++
		case c == '(':
			cur = append(cur, token{kind: tokLParen})
			i++
		case c == ')':
			cur = append(cur, token{kind: tokRParen})
			i++
		case c == ',':
			cur = append(cur, token{kind: tokComma})
			i++
		case c == '|' && i+1 < len(script) && script[i+1] == '>':
			cur = append(cur, token{kind: tokPipe})
			i += 2
		case isWordStart(c):
			start := i
			i++
			for i < len(script) && isWordPart(script[i]) {
				i++
			}
			word := script[start:i]
			if len(word) == 1 && i < len(script) {
				b := word[0]
				if ((b == 'r' || b == 'R') && (script[i] == '\'' || script[i] == '"')) ||
					((b == 'x' || b == 'X') && script[i] == '\'') {
					next, err := skipQuotedString(script, i, true)
					if err != nil {
						return nil, err
					}
					cur = append(cur, token{kind: tokLiteral})
					i = next
					continue
				}
			}
			cur = append(cur, token{kind: tokWord, val: strings.ToUpper(word)})
		default:
			cur = append(cur, token{kind: tokOther})
			i++
		}
	}
	if len(cur) > 0 {
		stmts = append(stmts, cur)
	}
	return stmts, nil
}

func skipNestedBlockComment(s string, start int) (int, error) {
	depth := 1
	i := start + 2
	for i < len(s) && depth > 0 {
		if s[i] == '/' && i+1 < len(s) && s[i+1] == '*' {
			depth++
			i += 2
			continue
		}
		if s[i] == '*' && i+1 < len(s) && s[i+1] == '/' {
			depth--
			i += 2
			continue
		}
		i++
	}
	if depth > 0 {
		return 0, errors.New("unterminated block comment")
	}
	return i, nil
}

func skipQuotedString(s string, start int, raw bool) (int, error) {
	quote := s[start]
	i := start + 1
	for i < len(s) {
		c := s[i]
		if !raw && c == '\\' {
			if i+1 >= len(s) {
				return 0, fmt.Errorf("unterminated string literal: missing closing %c", quote)
			}
			i += 2
			continue
		}
		if c == quote {
			if !raw && i+1 < len(s) && s[i+1] == quote {
				i += 2
				continue
			}
			return i + 1, nil
		}
		i++
	}
	return 0, fmt.Errorf("unterminated string literal: missing closing %c", quote)
}

func skipBacktickIdent(s string, start int) (int, error) {
	i := start + 1
	for i < len(s) {
		if s[i] == '`' {
			if i+1 < len(s) && s[i+1] == '`' {
				i += 2
				continue
			}
			return i + 1, nil
		}
		i++
	}
	return 0, errors.New("unterminated quoted identifier: missing closing `")
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

func isWordStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isWordPart(c byte) bool {
	return isWordStart(c) || (c >= '0' && c <= '9') || c == '$'
}

func isReadOnlyStatement(tokens []token) bool {
	if len(tokens) == 0 || !hasBalancedParens(tokens) {
		return false
	}
	switch tokens[0].kind {
	case tokLParen:
		return isReadOnlyQuery(tokens)
	case tokWord:
		switch tokens[0].val {
		case "SELECT", "VALUES", "TABLE", "SHOW":
			return len(tokens) >= 2
		case "DESCRIBE", "DESC":
			if len(tokens) < 2 {
				return false
			}
			if tokens[1].kind == tokWord {
				if tokens[1].val == "QUERY" {
					return isReadOnlyQuery(tokens[2:])
				}
				if isNonReadOnlyLeadKeyword(tokens[1].val) {
					return false
				}
			}
			return true
		case "WITH":
			return isReadOnlyWith(tokens)
		case "EXPLAIN":
			return isReadOnlyExplain(tokens)
		default:
			return false
		}
	default:
		return false
	}
}

func isReadOnlyExplain(tokens []token) bool {
	if len(tokens) < 2 {
		return false
	}
	idx := 1
	if tokens[idx].kind == tokWord {
		switch tokens[idx].val {
		case "ANALYZE":
			return false
		case "EXTENDED", "CODEGEN", "COST", "FORMATTED":
			idx++
		}
	}
	if idx >= len(tokens) {
		return false
	}
	if tokens[idx].kind == tokWord && tokens[idx].val == "EXPLAIN" {
		return false
	}
	return isReadOnlyStatement(tokens[idx:])
}

func isReadOnlyQuery(tokens []token) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokens[0].kind == tokLParen {
		closeIdx := matchingRParen(tokens, 0)
		if closeIdx < 0 || !isReadOnlyQuery(tokens[1:closeIdx]) {
			return false
		}
		rest := tokens[closeIdx+1:]
		if len(rest) == 0 {
			return true
		}
		if rest[0].kind == tokPipe {
			return true
		}
		if rest[0].kind == tokWord {
			switch rest[0].val {
			case "UNION", "INTERSECT", "EXCEPT", "MINUS":
				next := 1
				if next < len(rest) && rest[next].kind == tokWord && (rest[next].val == "ALL" || rest[next].val == "DISTINCT") {
					next++
				}
				if next+1 < len(rest) && rest[next].kind == tokWord && rest[next].val == "BY" &&
					rest[next+1].kind == tokWord && rest[next+1].val == "NAME" {
					next += 2
				}
				return isReadOnlyQuery(rest[next:])
			case "ORDER", "SORT", "CLUSTER", "DISTRIBUTE", "LIMIT", "OFFSET", "FETCH", "WINDOW", "QUALIFY":
				return true
			}
		}
		return false
	}
	if tokens[0].kind == tokWord {
		switch tokens[0].val {
		case "SELECT", "VALUES", "TABLE":
			return len(tokens) >= 2
		case "WITH":
			return isReadOnlyWith(tokens)
		}
	}
	return false
}

// isReadOnlyWith validates a CTE query:
// WITH [RECURSIVE] name [(cols)] AS ([NOT] MATERIALIZED) (query) [, ...] main_query
// Both each CTE body and the trailing main query must be read-only queries.
func isReadOnlyWith(tokens []token) bool {
	i := 1
	if i < len(tokens) && tokens[i].kind == tokWord && tokens[i].val == "RECURSIVE" {
		i++
	}
	for {
		if i >= len(tokens) || (tokens[i].kind != tokWord && tokens[i].kind != tokIdent) {
			return false
		}
		i++
		// Optional CTE column list: (c1, c2, ...)
		if i < len(tokens) && tokens[i].kind == tokLParen {
			closeCols := matchingRParen(tokens, i)
			if closeCols < 0 {
				return false
			}
			i = closeCols + 1
		}
		if i >= len(tokens) || tokens[i].kind != tokWord || tokens[i].val != "AS" {
			return false
		}
		i++
		if i+1 < len(tokens) && tokens[i].kind == tokWord && tokens[i].val == "NOT" &&
			tokens[i+1].kind == tokWord && tokens[i+1].val == "MATERIALIZED" {
			i += 2
		} else if i < len(tokens) && tokens[i].kind == tokWord && tokens[i].val == "MATERIALIZED" {
			i++
		}
		if i >= len(tokens) || tokens[i].kind != tokLParen {
			return false
		}
		closeBody := matchingRParen(tokens, i)
		if closeBody < 0 {
			return false
		}
		if !isReadOnlyQuery(tokens[i+1 : closeBody]) {
			return false
		}
		i = closeBody + 1
		if i < len(tokens) && tokens[i].kind == tokComma {
			i++
			continue
		}
		break
	}
	return isReadOnlyQuery(tokens[i:])
}

func isNonReadOnlyLeadKeyword(word string) bool {
	switch word {
	case "INSERT", "UPDATE", "DELETE", "MERGE", "COPY",
		"CREATE", "ALTER", "DROP", "TRUNCATE", "REPLACE",
		"OPTIMIZE", "VACUUM", "RESTORE", "REFRESH", "MSCK",
		"GRANT", "REVOKE", "DENY",
		"SET", "USE", "CALL", "EXECUTE":
		return true
	default:
		return false
	}
}

func hasBalancedParens(tokens []token) bool {
	depth := 0
	for _, tok := range tokens {
		switch tok.kind {
		case tokLParen:
			depth++
		case tokRParen:
			depth--
			if depth < 0 {
				return false
			}
		}
	}
	return depth == 0
}

func matchingRParen(tokens []token, openIdx int) int {
	if openIdx >= len(tokens) || tokens[openIdx].kind != tokLParen {
		return -1
	}
	depth := 0
	for i := openIdx; i < len(tokens); i++ {
		switch tokens[i].kind {
		case tokLParen:
			depth++
		case tokRParen:
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
