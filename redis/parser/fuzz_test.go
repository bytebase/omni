package parser

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/fuzz"
)

// FuzzParse holds Parse to the engine contract on any input (package fuzz).
// A Command carries its byte range in Start and End rather than a Loc.
// `go test` replays compat.tsv and testdata/fuzz/FuzzParse.
func FuzzParse(f *testing.F) {
	fuzz.Seed(f, filepath.Join("testdata", "compat.tsv"))
	f.Fuzz(func(t *testing.T, sql string) {
		cmds, err := Parse(sql)
		fuzz.Error[*ParseError](t, err, sql)
		for _, c := range cmds {
			fuzz.Range(t, sql, "Command "+c.Name, c.Start, c.End)
		}
	})
}
