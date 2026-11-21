package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/ElasticsearchConsole.g4 — generated from Elastic's
// official elasticsearch-specification (8.19, Apache-2.0) by cmd/specgrammar —
// and its near-misses, with omni's verdict and the grammar's. The judge is a
// recognizer over ElasticsearchConsole.g4, so every run asks it again and needs
// no server or network.
//
// testdata/README.md describes what the grammar covers and what the recognizer
// cannot decide.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "ElasticsearchConsole.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceGrammar,
	})
}

// omniVerdict accepts a row that parses cleanly as exactly one REST request.
func omniVerdict(sql string) string {
	res, err := Parse(sql)
	if err != nil || res == nil || len(res.Errors) > 0 || len(res.Requests) != 1 {
		return compat.Reject
	}
	return compat.Accept
}
