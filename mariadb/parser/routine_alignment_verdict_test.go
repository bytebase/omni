package parser

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"testing"
)

// verdictRecorder stands in for the test an oracle verdict helper is handed, so a
// test can see whether the helper fails it.
type verdictRecorder struct {
	testing.TB
	failed string
}

func (r *verdictRecorder) Helper() {}

func (r *verdictRecorder) Fatalf(format string, args ...any) {
	r.failed = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

// recordVerdict runs f against a recorder and returns the failure f reported, if
// any. f runs on its own goroutine because Fatalf ends the goroutine it runs on.
func recordVerdict(t *testing.T, f func(testing.TB)) string {
	rec := &verdictRecorder{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		f(rec)
	}()
	<-done
	return rec.failed
}

// TestRoutineEngineCodeFailsWithoutVerdict pins that a closed oracle connection is
// no verdict. TestRoutineAlignment used to count any error as MySQL refusing the
// routine, so a dead connection read as a rejection of every probe. This needs no
// container, so it guards the helper wherever it runs.
func TestRoutineEngineCodeFailsWithoutVerdict(t *testing.T) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:1)/")
	if err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	o := &parserOracle{db: db, ctx: context.Background()}

	var code uint16
	if msg := recordVerdict(t, func(tb testing.TB) { code = routineEngineCode(tb, o, "SELECT 1") }); msg == "" {
		t.Errorf("a closed connection came back as a verdict (code %d), not a failure", code)
	}
	if msg := recordVerdict(t, func(tb testing.TB) { mustExec(tb, o, "SELECT 1") }); msg == "" {
		t.Errorf("a seed on a closed connection did not fail the test")
	}
}
