// Package completion provides parser-native C3-style SQL completion for PostgreSQL.
package completion

import (
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
	return standardComplete(sql, cursorOffset, cat)
}

// standardComplete collects parser-level candidates using Collect, then
// resolves them against the catalog. Stub for now.
func standardComplete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	_ = parser.Collect(sql, cursorOffset)
	return nil
}

// trickyComplete handles edge cases that the standard C3 approach cannot
// resolve (e.g., partially typed identifiers in ambiguous positions). Stub.
func trickyComplete(_ string, _ int, _ *catalog.Catalog) []Candidate {
	return nil
}
