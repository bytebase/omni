package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/sql.g4, Microsoft's grammar of the Cosmos DB
// for NoSQL query language as the .NET SDK publishes it, and its
// near-misses, with omni's verdict and the grammar's. Cosmos DB has no
// server CI can start, so the grammar judges: a statement it derives is
// accepted, any other rejected. The judge is a recognizer, so every run
// asks it again and needs no network. Where Microsoft's query reference
// takes a statement the grammar does not derive, a row of
// ../testdata/grammar-overrides.tsv cites the page and the docs judge it.
//
// testdata/README.md says what the grammar judge cannot decide.
func TestCompatibility(t *testing.T) {
	grammar := filepath.Join("..", "testdata", "sql.g4")
	compat.Check(t, compat.Suite{
		Path:      filepath.Join("testdata", "compat.tsv"),
		Pins:      filepath.Join("..", "testdata", "pins.tsv"),
		Grammar:   grammar,
		Vendor:    []string{grammar},
		Tokens:    filepath.Join("..", "testdata", "tokens.tsv"),
		Overrides: filepath.Join("..", "testdata", "grammar-overrides.tsv"),
		Omni:      omniVerdict,
		Source:    compat.SourceGrammar,
	})
}

// omniVerdict accepts a row that parses as exactly one query.
func omniVerdict(sql string) string {
	list, err := Parse(sql)
	if err != nil || list == nil || list.Len() != 1 {
		return compat.Reject
	}
	return compat.Accept
}
