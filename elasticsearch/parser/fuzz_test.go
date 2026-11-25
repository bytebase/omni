package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/fuzz"
)

// FuzzParse holds Parse to the engine contract on any input (package fuzz).
// The console parser recovers from errors, so it reports each one in
// ParseResult.Errors and fails only on an input it cannot read at all.
// `go test` replays compat.tsv and testdata/fuzz/FuzzParse.
func FuzzParse(f *testing.F) {
	fuzz.Seed(f, filepath.Join("testdata", "compat.tsv"))
	f.Fuzz(func(t *testing.T, text string) {
		res, err := Parse(text)
		if err != nil {
			t.Fatalf("Parse failed outside ParseResult.Errors: %v", err)
		}
		for _, e := range res.Errors {
			if e.ByteOffset < 0 || e.ByteOffset > len(text) {
				t.Fatalf("SyntaxError at byte %d, outside a %d-byte input: %s", e.ByteOffset, len(text), e.Message)
			}
		}
		fuzz.Tree(t, res, text)
	})
}
