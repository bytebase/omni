package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/PartiQLParser.g4, the grammar of the PartiQL
// reference implementation, and its near-misses, with omni's verdict and
// the grammar's. PartiQL has no server CI can start, so the grammar judges:
// a statement it derives is accepted, any other rejected. The judge is a
// recognizer, so every run asks it again and needs no network.
//
// The grammar is the reference language, and DynamoDB, which Bytebase
// parses PartiQL for, takes a part of it. testdata/README.md says which
// judge decides what, and what neither can.
func TestCompatibility(t *testing.T) {
	parser := filepath.Join("..", "testdata", "PartiQLParser.g4")
	lexer := filepath.Join("..", "testdata", "PartiQLTokens.g4")
	compat.Check(t, compat.Suite{
		Path:      filepath.Join("testdata", "compat.tsv"),
		Pins:      filepath.Join("..", "testdata", "pins.tsv"),
		Grammar:   parser,
		Lexer:     lexer,
		Vendor:    []string{parser, lexer},
		Tokens:    filepath.Join("..", "testdata", "tokens.tsv"),
		Overrides: filepath.Join("..", "testdata", "grammar-overrides.tsv"),
		Omni:      omniVerdict,
		Source:    compat.SourceGrammar,
	})
}

// omniVerdict accepts a row that parses as exactly one statement.
func omniVerdict(sql string) string {
	list, err := Parse(sql)
	if err != nil || list == nil || list.Len() != 1 {
		return compat.Reject
	}
	return compat.Accept
}
