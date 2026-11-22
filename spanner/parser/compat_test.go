package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// TestCompatibility holds spanner/parser (DialectSpanner) to
// testdata/compat.tsv: a statement for each alternative of GoogleSQL's grammar
// (../../googlesql/testdata/GoogleSQL.g4, converted from googlesql.tm by
// googlesql/parser/cmd/tm2g4) and its near-misses, with omni's verdict and the
// pinned Cloud Spanner emulator's (../testdata/pins.tsv). The test reads the
// committed file and starts no server in offline replay. OMNI_RECORD=1 judges
// rows on the pinned emulator; OMNI_RECORD=verify re-judges every row read-only.
func TestCompatibility(t *testing.T) {
	gtd := filepath.Join("..", "..", "googlesql", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join(gtd, "GoogleSQL.g4"),
		Tokens:  filepath.Join(gtd, "tokens.tsv"),
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

// compatJudge asks the pinned Spanner emulator whether it accepts each
// statement, through googlesql/spanner-oracle, which runs a query, runs DML in
// a transaction it rolls back, and applies DDL to a database of its own. It
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

// spannerVerdict reads the harness's answer as a verdict. The harness rejects a
// statement on the parser's own error ("Syntax error:" for a query or DML, a
// grammar refusal, or "Error parsing Spanner DDL statement:" for DDL) and
// accepts it on any other answer of the emulator's, including Unimplemented
// (a query the emulator parsed and analyzed and cannot run, such as one
// returning a STRUCT column). Every other error is no verdict.
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
