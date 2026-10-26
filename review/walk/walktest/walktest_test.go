package walktest

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/review/walk"
)

// fakeEngine runs statements of the form "name" (accepted) or
// "name!code" (rejected with code), a rejected one ending in "+" changing
// the database. Its walk reports the statement named in report, with the
// message message.
type fakeEngine struct {
	digest  int
	ran     []string
	report  string
	message string
}

func (e *fakeEngine) Statements(change string) []string { return strings.Split(change, ";") }

func (e *fakeEngine) Exec(_ testing.TB, stmt string) Outcome {
	e.ran = append(e.ran, stmt)
	name, code, ok := strings.Cut(stmt, "!")
	if !ok {
		e.digest++
		return Outcome{}
	}
	if strings.HasSuffix(code, "+") {
		e.digest++
		code = strings.TrimSuffix(code, "+")
	}
	return Outcome{Code: code, Message: name}
}

func (e *fakeEngine) Report(_ testing.TB, change string) *Finding {
	stmts := strings.Split(change, "\n")
	for i, s := range stmts {
		if s == e.report {
			return &Finding{Statement: i, Message: e.message}
		}
	}
	return nil
}

func (e *fakeEngine) Violation(testing.TB, string) *walk.Violation { return &walk.Violation{} }

func (e *fakeEngine) Codes(message string) []string { return []string{message} }

func (e *fakeEngine) Holds(testing.TB, *walk.Violation, string) (bool, string) {
	return false, "the fake holds no claim"
}

func (e *fakeEngine) Digest(testing.TB) string { return fmt.Sprint(e.digest) }

type partialEngine struct{ fakeEngine }

func (*partialEngine) RejectsPartially() {}

// recorder is a testing.TB that records failures; Fatalf ends the
// goroutine, as testing.T's does.
type recorder struct {
	testing.TB
	errors []string
	fatal  string
}

func (r *recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *recorder) Fatalf(format string, args ...any) {
	r.fatal = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

func judge(t *testing.T, e Engine, change string) (r *recorder, found, rejected int) {
	r = &recorder{TB: t}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		found, rejected = Judge(r, e, change, true)
	}()
	wg.Wait()
	return r, found, rejected
}

func TestJudgeStopsAfterAPartialRejection(t *testing.T) {
	e := &partialEngine{fakeEngine{report: "b!1051+", message: "1051"}}
	r, found, rejected := judge(t, e, "a;b!1051+;c!1146")
	if r.fatal != "" || len(r.errors) != 0 || found != 1 || rejected != 1 {
		t.Errorf("fatal %q, errors %q, found %d of %d; want 1 of 1 and no failure", r.fatal, r.errors, found, rejected)
	}
	if got := strings.Join(e.ran, ";"); got != "a;b!1051+" {
		t.Errorf("ran %q, want the change to stop after the partial rejection", got)
	}
}

func TestJudgeHoldsAPartialRejectionToItsCode(t *testing.T) {
	e := &partialEngine{fakeEngine{report: "b!1051+", message: "1146"}}
	r, _, _ := judge(t, e, "a;b!1051+")
	if len(r.errors) != 1 || !strings.Contains(r.errors[0], "cannot be checked") {
		t.Errorf("errors %q, want one that the claim cannot be checked", r.errors)
	}
}

func TestJudgeFailsAPartialRejectionWithoutOptIn(t *testing.T) {
	e := &fakeEngine{report: "b!1051+", message: "1051"}
	r, _, _ := judge(t, e, "a;b!1051+;c")
	if !strings.Contains(r.fatal, "changed the database") {
		t.Errorf("fatal %q, want the change of the database reported", r.fatal)
	}
}

func TestJudgeGoesOnAfterAnAtomicRejection(t *testing.T) {
	e := &partialEngine{fakeEngine{report: "b!1051", message: "1051"}}
	r, found, rejected := judge(t, e, "a;b!1051;c!1146")
	// c's rejection is missed, which complete reports.
	if r.fatal != "" || len(r.errors) != 1 || found != 1 || rejected != 2 {
		t.Errorf("fatal %q, errors %q, found %d of %d; want 1 of 2 and the miss of c", r.fatal, r.errors, found, rejected)
	}
}
