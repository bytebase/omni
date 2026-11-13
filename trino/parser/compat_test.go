package parser

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/trino/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of Trino's SqlBase.g4 and its near-misses, with omni's verdict
// and Trino's. A rejection carries the byte it is reported at. The test
// reads the file and starts no server. OMNI_RECORD=1 judges the rows that
// have no verdict yet on the pinned Trino.
//
// testdata/README.md says what Trino's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "SqlBase.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as Trino's statement endpoint does: one
// statement, with nothing after it. Parse's error is a rejection at its
// byte. Parse reads a script, so a ';' ends the statement where Trino's
// grammar has no place for it, and a text with no statement is refused at
// its end, where Trino finds no statement.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		if errs := AllErrors(err); len(errs) > 0 && errs[0].Position >= 0 {
			return compat.RejectAt(errs[0].Position)
		}
		return compat.Reject
	}
	segs := Split(sql)
	switch {
	case len(segs) == 0:
		return compat.RejectAt(len(sql))
	case segs[0].ByteEnd < len(sql):
		return compat.RejectAt(segs[0].ByteEnd)
	}
	return compat.Accept
}

// compatJudge asks the pinned Trino whether its parser accepts each
// statement, by running it and reading the answer with compat.TrinoVerdict.
// A statement that parses may run, on the memory catalog the session
// defaults to; nothing it leaves changes whether a later one parses. It
// runs only in a recording, outside any tape.
func compatJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	o := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return compat.TrinoVerdict(sql, o.Run(ctx, sql), trinoError)
	}
}

// trinoError reads a Trino answer for compat.TrinoVerdict.
func trinoError(err error) (compat.TrinoError, bool) {
	var qe *testcontainer.QueryError
	if !errors.As(err, &qe) {
		return compat.TrinoError{}, false
	}
	return compat.TrinoError{Name: qe.Name, Line: qe.Line, Column: qe.Column}, true
}
