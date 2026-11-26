package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/fuzz"
)

// FuzzParse holds Parse to the engine contract on any input (package fuzz).
// `go test` replays compat.tsv and testdata/fuzz/FuzzParse.
func FuzzParse(f *testing.F) {
	fuzz.Seed(f, filepath.Join("testdata", "compat.tsv"))
	f.Fuzz(func(t *testing.T, sql string) {
		tree, err := Parse(sql)
		fuzz.Error[*ParseError](t, err, sql)
		fuzz.Tree(t, tree, sql)
	})
}
