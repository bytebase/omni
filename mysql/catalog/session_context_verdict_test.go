package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// setupRecorder stands in for the test a setup helper is handed, so a test can
// see whether the helper fails it.
type setupRecorder struct {
	testing.TB
	failed string
}

func (r *setupRecorder) Helper() {}

func (r *setupRecorder) Fatalf(format string, args ...any) {
	r.failed = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

func (r *setupRecorder) Skipf(format string, args ...any) {
	r.failed = "skipped: " + fmt.Sprintf(format, args...)
	runtime.Goexit()
}

// TestExecUnderSessionFailsWithoutSeed pins that a seed the engine never ran fails
// the test. execUnderSession used to log the error and let its callers skip, so a
// dropped connection passed the session-context proofs without running them. This
// needs no container, so it guards the helper wherever it runs.
func TestExecUnderSessionFailsWithoutSeed(t *testing.T) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:1)/")
	if err != nil {
		t.Fatal(err)
	}
	_ = db.Close()

	rec := &setupRecorder{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		execUnderSession(rec, db, context.Background(), "PIPES_AS_CONCAT", "", "SELECT 1")
	}()
	<-done
	if rec.failed == "" {
		t.Fatal("a seed on a closed connection did not fail the test")
	}
	if strings.HasPrefix(rec.failed, "skipped: ") {
		t.Fatalf("a seed on a closed connection skipped the test: %s", rec.failed)
	}
}
