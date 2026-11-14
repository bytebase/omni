package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/googlesql/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of GoogleSQL's grammar (../testdata/GoogleSQL.g4, converted
// from googlesql.tm by ./cmd/tm2g4) and its near-misses, with omni's verdict
// and the pinned Spanner emulator's. The test reads the file and starts no
// server. OMNI_RECORD=1 judges the rows that have no verdict yet on the
// pinned emulator.
//
// testdata/README.md says what the emulator's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "GoogleSQL.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as one statement, as the judge sends it.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the pinned emulator whether it accepts each statement,
// through googlesql/spanner-oracle, which runs a query, runs DML in a
// transaction it rolls back, and applies DDL to a database of its own. It
// runs only in a recording, outside any tape.
func compatJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	o := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		v, err := o.Verdict(ctx, sql)
		if err != nil {
			return "", err
		}
		return spannerVerdict(v)
	}
}

// spannerVerdict reads the harness's answer as a verdict. The harness
// rejects a statement on the parser's own error, "Syntax error:" for a
// query or DML and "Error parsing Spanner DDL statement:" for DDL, and
// accepts it on any other answer of the emulator's
// (googlesql/spanner-oracle/README.md). One more answer is an acceptance
// here: Unimplemented, which the emulator gives a query it parsed and
// analyzed and cannot run, as one that returns a STRUCT column. Every other
// error is no verdict.
func spannerVerdict(v testcontainer.Verdict) (string, error) {
	switch {
	case v.Verdict == "accept":
		return compat.Accept, nil
	case v.Verdict == "reject":
		return compat.Reject, nil
	case v.Code == "Unimplemented":
		return compat.Accept, nil
	}
	return "", fmt.Errorf("no verdict: %s %s (%s): %s", v.Verdict, v.Reason, v.Code, v.Message)
}
