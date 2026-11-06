package validate

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

// TestEngineExecAcceptsFailsWithoutVerdict pins that a closed oracle connection is no
// verdict. engineExecAccepts (once mysqlExecAccepts) used to read any error that was
// not a server error as a rejection, so a dead connection passed every probe that
// expects the engine to refuse.
// This needs no container, so it guards the helper wherever it runs.
func TestEngineExecAcceptsFailsWithoutVerdict(t *testing.T) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:1)/")
	if err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	o := &routineValidationOracle{db: db, ctx: context.Background()}

	var accepted bool
	if msg := recordVerdict(t, func(tb testing.TB) { accepted, _ = engineExecAccepts(tb, o, "SELECT 1") }); msg == "" {
		t.Errorf("a closed connection came back as a verdict (accepted=%v), not a failure", accepted)
	}
}
