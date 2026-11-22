package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/SqlBaseParser.g4 (Apache Spark's official ANTLR 4
// SQL grammar, which Databricks SQL is built on) and its near-misses, with
// omni's verdict and the grammar's. The judge is a recognizer over
// SqlBaseParser.g4 and SqlBaseLexer.g4, so every run asks it again and needs
// no server or network.
//
// testdata/README.md describes what the grammar covers, what the recognizer
// cannot decide, and how to upgrade to a live Databricks SQL warehouse judge
// when DATABRICKS_* credentials are provided.
func TestCompatibility(t *testing.T) {
	grammar := filepath.Join("..", "testdata", "SqlBaseParser.g4")
	lexer := filepath.Join("..", "testdata", "SqlBaseLexer.g4")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: grammar,
		Lexer:   lexer,
		Vendor:  []string{grammar, lexer},
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceGrammar,
	})
}

// omniVerdict accepts a row that parses cleanly as one Databricks SQL statement.
func omniVerdict(sql string) string {
	res, err := Parse(sql)
	if err != nil || res == nil || len(res.Statements) != 1 {
		return compat.Reject
	}
	return compat.Accept
}
