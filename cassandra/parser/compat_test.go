package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/Cql.g4 and its near-misses, with omni's
// verdict and the grammar's. Cql.g4 is Apache Cassandra's own CQL grammar,
// ANTLR 3, converted to ANTLR 4 by cmd/antlr3to4 (TestCqlG4IsTheConversion
// holds it to that). Cassandra has no server this suite starts, so the
// grammar judges: a statement it derives is accepted, any other rejected.
// The judge is a recognizer, so every run asks it again and needs no
// network.
//
// testdata/README.md says what the grammar judge cannot decide.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:      filepath.Join("testdata", "compat.tsv"),
		Pins:      filepath.Join(td, "pins.tsv"),
		Grammar:   filepath.Join(td, "Cql.g4"),
		Vendor:    []string{filepath.Join(td, "Cql.g"), filepath.Join(td, "Parser.g"), filepath.Join(td, "Lexer.g")},
		Tokens:    filepath.Join(td, "tokens.tsv"),
		Overrides: filepath.Join(td, "grammar-overrides.tsv"),
		Omni:      omniVerdict,
		Source:    compat.SourceGrammar,
	})
}

// omniVerdict accepts a row that parses as exactly one statement.
func omniVerdict(sql string) string {
	list, err := Parse(sql)
	if err != nil || list == nil || len(list.Items) != 1 {
		return compat.Reject
	}
	return compat.Accept
}
