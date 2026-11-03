package parser

import (
	"database/sql/driver"
	"fmt"
	"runtime"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
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

// TestTiDBRejectedSyntaxFailsWithoutVerdict pins that a server which never answered
// is no verdict. tidbRejectedSyntax used to read any error that was not a server
// error as acceptance, so a dropped connection passed every case that expects TiDB
// to accept. This needs no container, so it guards the helper wherever it runs.
func TestTiDBRejectedSyntaxFailsWithoutVerdict(t *testing.T) {
	var rejected bool
	if msg := recordVerdict(t, func(tb testing.TB) { rejected = tidbRejectedSyntax(tb, driver.ErrBadConn) }); msg == "" {
		t.Errorf("a dropped connection came back as a verdict (rejected=%v), not a failure", rejected)
	}
	// The verdicts the helper exists for are unchanged.
	for _, c := range []struct {
		err  error
		want bool
	}{
		{nil, false},
		{&gomysql.MySQLError{Number: 1064}, true},
		{&gomysql.MySQLError{Number: 1146}, false},
	} {
		if msg := recordVerdict(t, func(tb testing.TB) { rejected = tidbRejectedSyntax(tb, c.err) }); msg != "" {
			t.Errorf("%v: failed the test: %s", c.err, msg)
		} else if rejected != c.want {
			t.Errorf("%v: rejected=%v, want %v", c.err, rejected, c.want)
		}
	}
}
