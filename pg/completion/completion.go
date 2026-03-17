// Package completion provides parser-native C3-style SQL completion for PostgreSQL.
package completion

import (
	"strings"

	"github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/pg/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType int

const (
	CandidateKeyword          CandidateType = iota // SQL keyword
	CandidateSchema                                // schema name
	CandidateTable                                 // table name
	CandidateView                                  // view name
	CandidateMaterializedView                      // materialized view name
	CandidateColumn                                // column name
	CandidateFunction                              // function name
	CandidateSequence                              // sequence name
	CandidateIndex                                 // index name
	CandidateType_                                 // SQL type name
	CandidateTrigger                               // trigger name
	CandidatePolicy                                // policy name
	CandidateExtension                             // extension name
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition/signature
	Comment    string        // optional doc comment
}

// Complete returns completion candidates for the given SQL at the cursor offset.
// cat may be nil if no catalog context is available.
func Complete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
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
// resolves them against the catalog. Stub for now.
func standardComplete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	cs := parser.Collect(sql, cursorOffset)
	return resolve(cs, cat, sql, cursorOffset)
}

// trickyComplete handles edge cases that the standard C3 approach cannot
// resolve (e.g., partially typed identifiers in ambiguous positions).
// It inserts placeholder tokens to make incomplete SQL parseable, then
// collects candidates from the patched input.
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
