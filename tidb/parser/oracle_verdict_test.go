package parser

import (
	"database/sql/driver"
	"errors"
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

// TestPrepareVerdictFailsWithoutVerdict pins the same property for the corpus
// harness's reader of PREPARE errors: a server that never answered is no
// verdict, never silent agreement. The corpus has a thousand rows the engine
// accepts, so reading a dropped connection as acceptance would turn the gap
// count into a number about the network. This needs no container.
func TestPrepareVerdictFailsWithoutVerdict(t *testing.T) {
	unprobed := func() error {
		t.Errorf("the verdict asked the probe")
		return nil
	}
	if _, err := prepareVerdict(driver.ErrBadConn, unprobed); err == nil {
		t.Errorf("a dropped connection came back as a verdict, not as an error")
	}
	// go-sql-driver refuses a local-file request the test did not register, and
	// that error is not a server error either.
	if _, err := prepareVerdict(errors.New("local file '/tmp/stats.json' is not registered"), unprobed); err == nil {
		t.Errorf("a refused local-file request came back as a verdict, not as an error")
	}
	// A probe that lost the server is no verdict, and neither is a probe the server
	// took, which would mean it did not read the tail as a syntax error at all.
	missing := &gomysql.MySQLError{Number: 1146}
	if _, err := prepareVerdict(missing, func() error { return driver.ErrBadConn }); err == nil {
		t.Errorf("a probe on a dropped connection came back as a verdict")
	}
	if _, err := prepareVerdict(missing, func() error { return nil }); err == nil {
		t.Errorf("a probe the server took came back as a verdict")
	}
}

// TestPrepareVerdictReadsTheProbe pins how a verdict reads the code. A grammar action's
// code comes back with the tail after the statement, and anything raised after the
// parse gives way to the tail's 1064. This needs no container; the live control is
// TestSyntaxProbeControls.
func TestPrepareVerdictReadsTheProbe(t *testing.T) {
	for _, c := range []struct {
		name  string
		err   error
		probe uint16 // 0: the verdict must not ask
		want  bool
	}{
		{"prepared", nil, 0, true},
		{"parse error", &gomysql.MySQLError{Number: 1064}, 0, false},
		{"a name it looks up", &gomysql.MySQLError{Number: 1193}, 0, true},
		{"table missing", &gomysql.MySQLError{Number: 1146}, 1064, true},
		{"not preparable", &gomysql.MySQLError{Number: 1295}, 1064, true},
		{"grammar action", &gomysql.MySQLError{Number: 1149}, 1149, false},
		{"a code raised both ways, after the parse", &gomysql.MySQLError{Number: 1105}, 1064, true},
		{"a code raised both ways, while parsing", &gomysql.MySQLError{Number: 1105}, 1105, false},
	} {
		got, err := prepareVerdict(c.err, func() error {
			if c.probe == 0 {
				t.Errorf("%s: the verdict asked the probe", c.name)
				return nil
			}
			return &gomysql.MySQLError{Number: c.probe}
		})
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
		} else if got != c.want {
			t.Errorf("%s: accepts=%v, want %v", c.name, got, c.want)
		}
	}
}
