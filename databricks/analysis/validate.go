package analysis

import (
	"errors"
	"fmt"
	"strings"
)

// maxNestingDepth bounds parenthesized nesting in ValidateQuery, matching
// omni's repository-wide nesting bound so deeply nested input fails closed
// without quadratic scans or stack growth.
const maxNestingDepth = 10000

var errNestingDepthExceeded = errors.New("maximum nesting depth exceeded")

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns rows:
// (readOnly, returnsData, error).
//
// Decision is lexical: every non-empty statement in a script must be a
// read-only Databricks SQL query (SELECT, WITH ... SELECT/VALUES/TABLE/FROM
// pipeline, FROM ... |> ... pipeline, VALUES, TABLE, SHOW, DESCRIBE/DESC, or
// non-executing EXPLAIN over a read-only statement). EXPLAIN ANALYZE, EXPLAIN
// over writes, CTEs containing DML or feeding non-query statements, FROM-first
// inserts, DML (INSERT, UPDATE, DELETE, MERGE, COPY INTO), DDL (CREATE, ALTER,
// DROP, TRUNCATE, OPTIMIZE, VACUUM, RESTORE, REFRESH, MSCK), DCL (GRANT,
// REVOKE, DENY), SET, USE, CALL, and multi-statement batches containing any
// non-read-only statement are rejected. Empty, whitespace-only, comment-only,
// or semicolon-only input carries no statements to refuse and returns
// (true, true, nil). Lexical errors (unterminated comments, string literals, or
// quoted identifiers) and excessive nesting depth are returned so callers fail
// closed.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	stmts, err := tokenizeStatements(statement)
	if err != nil {
		return false, false, err
	}
	for _, stmt := range stmts {
		pairs, balanced, err := matchParens(stmt)
		if err != nil {
			return false, false, err
		}
		if !balanced || !isReadOnlyStatement(stmt, pairs, 0, len(stmt)) {
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
	tokBar
	tokOther
)

type token struct {
	kind tokenKind
	val  string // uppercased text for tokWord; empty otherwise
}

// tokenizeStatements lexes Databricks SQL into per-statement token slices split
// on top-level semicolons. It handles Spark's `--` line comments (terminating
// on CR or LF and supporting `\` + LF continuation), nested `/* /* */ */` block
// comments, single/double-quoted strings with backslash and doubled-quote
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
			// Match Spark's SIMPLE_COMMENT rule: '--' ('\\\n' | ~[\r\n])* '\r'? '\n'?
			i += 2
			for i < len(script) {
				if script[i] == '\\' && i+1 < len(script) && script[i+1] == '\n' {
					i += 2
					continue
				}
				if script[i] == '\r' || script[i] == '\n' {
					if script[i] == '\r' && i+1 < len(script) && script[i+1] == '\n' {
						i += 2
					} else {
						i++
					}
					break
				}
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
		case c == '|' && i+1 < len(script) && script[i+1] == '|':
			cur = append(cur, token{kind: tokOther})
			i += 2
		case c == '|':
			cur = append(cur, token{kind: tokBar})
			i++
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

// matchParens precomputes matching closing parenthesis indices in a single O(n)
// pass and enforces maxNestingDepth. Unbalanced parentheses return balanced=false
// with a nil error; exceeding maxNestingDepth returns errNestingDepthExceeded.
func matchParens(tokens []token) ([]int, bool, error) {
	pairs := make([]int, len(tokens))
	for i := range pairs {
		pairs[i] = -1
	}
	var stack []int
	for i, tok := range tokens {
		switch tok.kind {
		case tokLParen:
			if len(stack) >= maxNestingDepth {
				return nil, false, errNestingDepthExceeded
			}
			stack = append(stack, i)
		case tokRParen:
			if len(stack) == 0 {
				return nil, false, nil
			}
			openIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pairs[openIdx] = i
		}
	}
	if len(stack) != 0 {
		return nil, false, nil
	}
	return pairs, true, nil
}

func isReadOnlyStatement(tokens []token, pairs []int, start, end int) bool {
	if start >= end {
		return false
	}
	switch tokens[start].kind {
	case tokLParen:
		return isReadOnlyQuery(tokens, pairs, start, end)
	case tokWord:
		switch tokens[start].val {
		case "SELECT", "VALUES", "TABLE", "SHOW":
			return end-start >= 2
		case "DESCRIBE", "DESC":
			if end-start < 2 {
				return false
			}
			if tokens[start+1].kind == tokWord {
				if tokens[start+1].val == "QUERY" {
					return isReadOnlyQuery(tokens, pairs, start+2, end)
				}
				if isNonReadOnlyLeadKeyword(tokens[start+1].val) {
					return false
				}
			}
			return true
		case "WITH":
			return isReadOnlyWith(tokens, pairs, start, end)
		case "FROM":
			return isReadOnlyFromPipeline(tokens, pairs, start, end)
		case "EXPLAIN":
			return isReadOnlyExplain(tokens, pairs, start, end)
		default:
			return false
		}
	default:
		return false
	}
}

func isReadOnlyExplain(tokens []token, pairs []int, start, end int) bool {
	if end-start < 2 {
		return false
	}
	idx := start + 1
	if tokens[idx].kind == tokWord {
		switch tokens[idx].val {
		case "ANALYZE":
			return false
		case "EXTENDED", "CODEGEN", "COST", "FORMATTED":
			idx++
		}
	}
	if idx >= end {
		return false
	}
	if tokens[idx].kind == tokWord && tokens[idx].val == "EXPLAIN" {
		return false
	}
	return isReadOnlyStatement(tokens, pairs, idx, end)
}

func isReadOnlyQuery(tokens []token, pairs []int, start, end int) bool {
	for {
		// Unwrap full-span outer parentheses iteratively so deeply nested
		// queries stay O(n) with constant stack space.
		for start < end && tokens[start].kind == tokLParen && pairs[start] == end-1 {
			start++
			end--
		}
		if start >= end {
			return false
		}
		if tokens[start].kind == tokLParen {
			closeIdx := pairs[start]
			if closeIdx < 0 || closeIdx >= end || !isReadOnlyQuery(tokens, pairs, start+1, closeIdx) {
				return false
			}
			restStart := closeIdx + 1
			if restStart >= end {
				return true
			}
			if tokens[restStart].kind == tokPipe || tokens[restStart].kind == tokBar {
				return isReadOnlyPipeTail(tokens, pairs, restStart, end)
			}
			if tokens[restStart].kind == tokWord {
				switch tokens[restStart].val {
				case "UNION", "INTERSECT", "EXCEPT", "MINUS":
					next := restStart + 1
					if next < end && tokens[next].kind == tokWord && (tokens[next].val == "ALL" || tokens[next].val == "DISTINCT") {
						next++
					}
					if next+1 < end && tokens[next].kind == tokWord && tokens[next].val == "BY" &&
						tokens[next+1].kind == tokWord && tokens[next+1].val == "NAME" {
						next += 2
					}
					start = next
					continue
				case "ORDER", "SORT", "CLUSTER", "DISTRIBUTE", "LIMIT", "OFFSET", "FETCH", "WINDOW", "QUALIFY":
					return true
				}
			}
			return false
		}
		if tokens[start].kind == tokWord {
			switch tokens[start].val {
			case "SELECT", "VALUES", "TABLE":
				return end-start >= 2
			case "WITH":
				return isReadOnlyWith(tokens, pairs, start, end)
			case "FROM":
				return isReadOnlyFromPipeline(tokens, pairs, start, end)
			}
		}
		return false
	}
}

// isReadOnlyFromPipeline validates Databricks SQL FROM-first pipeline queries
// using either `|>` or Databricks Runtime 18.0+ `|` separators, while rejecting
// FROM-first inserts (FROM s INSERT INTO t SELECT ...) and bare FROM clauses.
func isReadOnlyFromPipeline(tokens []token, pairs []int, start, end int) bool {
	for i := start + 1; i < end; i++ {
		switch tokens[i].kind {
		case tokLParen:
			closeIdx := pairs[i]
			if closeIdx < 0 || closeIdx >= end {
				return false
			}
			if !isReadOnlyParenFragment(tokens, pairs, i+1, closeIdx) {
				return false
			}
			i = closeIdx
		case tokPipe:
			if i == start+1 {
				return false
			}
			return isReadOnlyPipeTail(tokens, pairs, i, end)
		case tokBar:
			if i == start+1 {
				return false
			}
			if i+1 < end && tokens[i+1].kind == tokWord && isReadOnlyPipeOperator(tokens[i+1].val) {
				return isReadOnlyPipeTail(tokens, pairs, i, end)
			}
		case tokWord:
			if isInsertClause(tokens, i, end) {
				return false
			}
		}
	}
	return false
}

// isReadOnlyPipeTail validates one or more `|>` / `|` pipe stages starting at
// start. Each stage must begin with a documented read-only Databricks SQL pipe
// operator; bitwise OR `|` inside stage expressions is distinguished from stage
// separators, and write operations are rejected by syntactic position rather
// than blocking scalar function names such as replace(...) in expressions.
func isReadOnlyPipeTail(tokens []token, pairs []int, start, end int) bool {
	i := start
	for i < end {
		if tokens[i].kind != tokPipe && tokens[i].kind != tokBar {
			return false
		}
		i++
		if i >= end || tokens[i].kind != tokWord || !isReadOnlyPipeOperator(tokens[i].val) {
			return false
		}
		i++
		for i < end && tokens[i].kind != tokPipe {
			if tokens[i].kind == tokBar {
				if i+1 >= end {
					return false
				}
				if tokens[i+1].kind == tokWord {
					if isReadOnlyPipeOperator(tokens[i+1].val) {
						break
					}
					isCall := i+2 < end && tokens[i+2].kind == tokLParen
					if !isCall && (isNonReadOnlyLeadKeyword(tokens[i+1].val) || tokens[i+1].val == "EXPORT") {
						return false
					}
				}
				i++
				continue
			}
			switch tokens[i].kind {
			case tokLParen:
				closeIdx := pairs[i]
				if closeIdx < 0 || closeIdx >= end {
					return false
				}
				if !isReadOnlyParenFragment(tokens, pairs, i+1, closeIdx) {
					return false
				}
				i = closeIdx
			case tokWord:
				if isInsertClause(tokens, i, end) {
					return false
				}
			}
			i++
		}
	}
	return true
}

// isReadOnlyParenFragment checks a parenthesized fragment inside a FROM or pipe
// stage: if it starts a WITH or FROM query it must pass isReadOnlyQuery, and if
// it starts a write statement (not a scalar function call like replace(...)) it
// is rejected.
func isReadOnlyParenFragment(tokens []token, pairs []int, start, end int) bool {
	for start < end && tokens[start].kind == tokLParen && pairs[start] == end-1 {
		start++
		end--
	}
	if start >= end || tokens[start].kind != tokWord {
		return true
	}
	if start+1 < end && tokens[start+1].kind == tokLParen {
		return true
	}
	if tokens[start].val == "WITH" || tokens[start].val == "FROM" {
		return isReadOnlyQuery(tokens, pairs, start, end)
	}
	return !isNonReadOnlyLeadKeyword(tokens[start].val)
}

func isInsertClause(tokens []token, idx, end int) bool {
	if tokens[idx].kind != tokWord || tokens[idx].val != "INSERT" || idx+1 >= end || tokens[idx+1].kind != tokWord {
		return false
	}
	switch tokens[idx+1].val {
	case "INTO", "OVERWRITE", "TABLE":
		return true
	default:
		return false
	}
}

func isReadOnlyPipeOperator(word string) bool {
	switch word {
	case "SELECT", "EXTEND", "SET", "DROP", "AS",
		"WHERE", "LIMIT", "OFFSET", "AGGREGATE",
		"ORDER", "SORT", "CLUSTER", "DISTRIBUTE",
		"UNION", "INTERSECT", "EXCEPT", "MINUS",
		"JOIN", "INNER", "LEFT", "RIGHT", "FULL", "CROSS", "SEMI", "ANTI", "NATURAL", "LATERAL",
		"TABLESAMPLE", "PIVOT", "UNPIVOT", "WINDOW", "QUALIFY", "HAVING":
		return true
	default:
		return false
	}
}

// isReadOnlyWith validates a CTE query:
// WITH [RECURSIVE] name [(cols)] AS ([NOT] MATERIALIZED) (query) [, ...] main_query
// Both each CTE body and the trailing main query must be read-only queries.
func isReadOnlyWith(tokens []token, pairs []int, start, end int) bool {
	i := start + 1
	if i < end && tokens[i].kind == tokWord && tokens[i].val == "RECURSIVE" {
		i++
	}
	for {
		if i >= end || (tokens[i].kind != tokWord && tokens[i].kind != tokIdent) {
			return false
		}
		i++
		// Optional CTE column list: (c1, c2, ...)
		if i < end && tokens[i].kind == tokLParen {
			closeCols := pairs[i]
			if closeCols < 0 || closeCols >= end {
				return false
			}
			i = closeCols + 1
		}
		if i >= end || tokens[i].kind != tokWord || tokens[i].val != "AS" {
			return false
		}
		i++
		if i+1 < end && tokens[i].kind == tokWord && tokens[i].val == "NOT" &&
			tokens[i+1].kind == tokWord && tokens[i+1].val == "MATERIALIZED" {
			i += 2
		} else if i < end && tokens[i].kind == tokWord && tokens[i].val == "MATERIALIZED" {
			i++
		}
		if i >= end || tokens[i].kind != tokLParen {
			return false
		}
		closeBody := pairs[i]
		if closeBody < 0 || closeBody >= end {
			return false
		}
		if !isReadOnlyQuery(tokens, pairs, i+1, closeBody) {
			return false
		}
		i = closeBody + 1
		if i < end && tokens[i].kind == tokComma {
			i++
			continue
		}
		break
	}
	return isReadOnlyQuery(tokens, pairs, i, end)
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
