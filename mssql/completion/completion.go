// Package completion provides parser-native C3-style SQL completion for MSSQL.
package completion

import (
	"strings"

	"github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/mssql/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType int

const (
	CandidateKeyword   CandidateType = iota // SQL keyword
	CandidateSchema                         // schema name
	CandidateTable                          // table name
	CandidateView                           // view name
	CandidateColumn                         // column name
	CandidateFunction                       // function name
	CandidateProcedure                      // stored procedure name
	CandidateIndex                          // index name
	CandidateTrigger                        // trigger name
	CandidateSequence                       // sequence name
	CandidateType_                          // SQL type name
	CandidateLogin                          // login name
	CandidateUser                           // user name
	CandidateRole                           // role name
	CandidateDatabase                       // database name
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition/signature
	Comment    string        // optional doc comment
}

// Complete returns completion candidates for the given SQL at the cursor offset.
// cat may be nil if no catalog context is available; in that case only keyword
// candidates are returned.
func Complete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	sql, cursorOffset = isolateStatementAtCursor(sql, cursorOffset)
	prefix := extractPrefix(sql, cursorOffset)

	// When the cursor is mid-token, back up to the start of the token
	// so the parser sees the position before the partial text.
	collectOffset := cursorOffset - len(prefix)

	result := standardComplete(sql, collectOffset, cat)
	if len(result) == 0 {
		result = trickyComplete(sql, collectOffset, cat)
	}

	return filterByPrefix(result, prefix)
}

// isolateStatementAtCursor slices sql to the statement containing cursorOffset
// between surrounding top-level semicolons or standalone GO batch commands, so
// syntax errors or table references in earlier or following statements do not
// leak into the active completion scope.
func isolateStatementAtCursor(sql string, cursorOffset int) (string, int) {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	tokens := parser.Tokenize(sql)
	depth := 0
	startByte := 0
	endByte := len(sql)
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if sepEnd, ok := statementSeparatorEnd(sql, tokens, i, depth); ok {
			if sepEnd <= cursorOffset {
				startByte = sepEnd
				continue
			}
			if tok.Loc >= cursorOffset {
				endByte = tok.Loc
				break
			}
		}
		if tok.Loc >= cursorOffset {
			continue
		}
		switch tok.Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
	}
	if endByte < cursorOffset {
		endByte = cursorOffset
	}
	if startByte > cursorOffset {
		startByte = 0
	}
	if startByte == 0 && endByte == len(sql) {
		return sql, cursorOffset
	}
	return sql[startByte:endByte], cursorOffset - startByte
}

func statementSeparatorEnd(sql string, tokens []parser.Token, i int, depth int) (int, bool) {
	tok := tokens[i]
	if tok.Type == ';' {
		if depth == 0 {
			return tok.End, true
		}
		return 0, false
	}
	if depth != 0 || tok.Type != parser.GO {
		return 0, false
	}
	return standaloneGoSeparatorEnd(sql, tokens, i)
}

func standaloneGoSeparatorEnd(sql string, tokens []parser.Token, i int) (int, bool) {
	tok := tokens[i]
	if tok.Loc < 0 || tok.Loc >= len(sql) || sql[tok.Loc] == '[' || sql[tok.Loc] == '"' {
		return 0, false
	}
	if i > 0 && tokens[i-1].Type != ';' && !hasLineBreakBetween(sql, tokens[i-1].End, tok.Loc) {
		return 0, false
	}
	sepEnd := tok.End
	nextIdx := i + 1
	if nextIdx < len(tokens) && tokens[nextIdx].Type == parser.ICONST && !hasLineBreakBetween(sql, sepEnd, tokens[nextIdx].Loc) {
		sepEnd = tokens[nextIdx].End
		nextIdx++
	}
	if nextIdx < len(tokens) && tokens[nextIdx].Type != ';' && !hasLineBreakBetween(sql, sepEnd, tokens[nextIdx].Loc) {
		return 0, false
	}
	return sepEnd, true
}

func hasLineBreakBetween(sql string, start, end int) bool {
	if start < 0 {
		start = 0
	}
	if end > len(sql) {
		end = len(sql)
	}
	if start >= end {
		return false
	}
	return strings.ContainsAny(sql[start:end], "\r\n")
}

// extractPrefix returns the partial token the user is typing at cursorOffset.
func extractPrefix(sql string, cursorOffset int) string {
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	i := cursorOffset
	for i > 0 {
		c := sql[i-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' {
			i--
		} else {
			break
		}
	}
	return sql[i:cursorOffset]
}

// filterByPrefix filters candidates whose text starts with prefix (case-insensitive).
func filterByPrefix(candidates []Candidate, prefix string) []Candidate {
	if prefix == "" {
		return candidates
	}
	upper := strings.ToUpper(prefix)
	var result []Candidate
	for _, c := range candidates {
		if strings.HasPrefix(strings.ToUpper(c.Text), upper) {
			result = append(result, c)
		}
	}
	return result
}

// standardComplete collects parser-level candidates using Collect, then
// resolves them into Candidate values.
func standardComplete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	cs := parser.Collect(sql, cursorOffset)
	return resolve(cs, cat, sql, cursorOffset)
}

// trickyComplete handles edge cases that the standard C3 approach cannot
// resolve (e.g., trailing spaces after FROM, truncated keywords, commas, operators).
// It inserts placeholder tokens to make incomplete SQL parseable, then
// collects candidates from the patched SQL.
func trickyComplete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	prefix := sql[:cursorOffset]
	suffix := ""
	if cursorOffset < len(sql) {
		suffix = sql[cursorOffset:]
	}

	// Try several placeholder strategies to make the SQL parseable.
	strategies := []string{
		prefix + " __placeholder__" + suffix,  // add identifier
		prefix + " __placeholder__ " + suffix, // add identifier with space
		prefix + " 1" + suffix,                // numeric placeholder
	}

	for _, patched := range strategies {
		cs := parser.Collect(patched, cursorOffset)
		if cs != nil && (len(cs.Tokens) > 0 || len(cs.Rules) > 0) {
			return resolve(cs, cat, sql, cursorOffset)
		}
	}
	return nil
}

// resolve, resolveRule, dedup, and related helpers are in resolve.go.
