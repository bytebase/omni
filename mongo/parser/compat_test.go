package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/MongoShell.g4 (converted from mongodb-js/mongosh
// TypeScript shell-api and shell-bson definitions by cmd/mongoshgrammar) and
// its near-misses, with omni's verdict and the docs recognizer's.
//
// MongoDB's server does not parse mongosh shell syntax, so the docs recognizer
// over MongoShell.g4 judges: a statement the call-shape grammar derives is
// accepted, and a near-miss it does not derive is undocumented (passing
// whatever omni says so omni is never pushed to reject valid JavaScript /
// mongosh expressions outside the call-shape subset). The judge is a
// recognizer, so every run asks it again and needs no network or server.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		PinRole: "mongosh",
		Grammar: filepath.Join(td, "MongoShell.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceDocs,
	})
}

// omniVerdict accepts a row that parses as exactly one statement with no errors.
func omniVerdict(sql string) string {
	res := ParseBestEffort(sql)
	if len(res.Errors) != 0 || len(res.Nodes) != 1 {
		return compat.Reject
	}
	return compat.Accept
}
