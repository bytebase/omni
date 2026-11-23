// Package completion provides parser-native C3-style SQL completion for MSSQL.
package completion

import (
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/mssql/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType = analysis.CandidateType

const (
	CandidateKeyword   = analysis.CandidateKeyword   // SQL keyword
	CandidateSchema    = analysis.CandidateSchema    // schema name
	CandidateTable     = analysis.CandidateTable     // table name
	CandidateView      = analysis.CandidateView      // view name
	CandidateColumn    = analysis.CandidateColumn    // column name
	CandidateFunction  = analysis.CandidateFunction  // function name
	CandidateProcedure = analysis.CandidateProcedure // stored procedure name
	CandidateIndex     = analysis.CandidateIndex     // index name
	CandidateTrigger   = analysis.CandidateTrigger   // trigger name
	CandidateSequence  = analysis.CandidateSequence  // sequence name
	CandidateType_     = analysis.CandidateType_     // SQL type name
	CandidateLogin     = analysis.CandidateLogin     // login name
	CandidateUser      = analysis.CandidateUser      // user name
	CandidateRole      = analysis.CandidateRole      // role name
	CandidateDatabase  = analysis.CandidateDatabase  // database name
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
	if collectOffset > 0 && collectOffset <= len(sql) {
		if delim := caretQuoteDelimiter(sql, cursorOffset); delim != 0 && sql[collectOffset-1] == delim {
			collectOffset--
		}
	}

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
	startByte, endByte := parser.StatementBounds(sql, cursorOffset)
	return sql[startByte:endByte], cursorOffset - startByte
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
