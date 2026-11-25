package parser

import (
	"errors"
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
		// PostgreSQL reports a grammar action's error, such as a constraint
		// attribute its kind cannot take, with no position, and so does
		// Parse: Position -1, which pg/review places at the statement.
		var pe *ParseError
		if errors.As(err, &pe) && pe.Position == -1 {
			return
		}
		fuzz.Error[*ParseError](t, err, sql)
		fuzz.Tree(t, tree, sql)
	})
}
