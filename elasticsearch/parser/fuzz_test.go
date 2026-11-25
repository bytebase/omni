package parser

import (
	"path/filepath"
	"testing"
	"unicode/utf8"

	"github.com/bytebase/omni/enginetest/fuzz"
)

// FuzzParse holds Parse to the engine contract on any input (package fuzz).
// The console parser recovers from errors, so it reports each one in
// ParseResult.Errors and fails only on text it cannot decode as UTF-8. A
// request carries its byte range in StartOffset and EndOffset, not a Loc.
// `go test` replays compat.tsv and testdata/fuzz/FuzzParse.
func FuzzParse(f *testing.F) {
	fuzz.Seed(f, filepath.Join("testdata", "compat.tsv"))
	f.Fuzz(func(t *testing.T, text string) {
		res, err := Parse(text)
		if err != nil {
			if utf8.ValidString(text) {
				t.Fatalf("Parse(%q) failed outside ParseResult.Errors: %v", text, err)
			}
			return
		}
		for _, e := range res.Errors {
			if e.ByteOffset < 0 || e.ByteOffset > len(text) {
				t.Fatalf("SyntaxError at byte %d, outside a %d-byte input: %s", e.ByteOffset, len(text), e.Message)
			}
		}
		for _, r := range res.Requests {
			fuzz.Range(t, text, "ParsedRequest", r.StartOffset, r.EndOffset)
		}
	})
}
