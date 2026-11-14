package cputime

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func spin(d time.Duration) {
	for start := Now(); Now()-start < d; {
	}
}

func TestRunMeasuresCPUTime(t *testing.T) {
	used, done := Run(time.Minute, func() { spin(50 * time.Millisecond) })
	if !done {
		t.Fatal("Run gave up on work that returns")
	}
	if used < 50*time.Millisecond {
		t.Errorf("Run reported %v for 50ms of CPU work", used)
	}
}

// A goroutine that sleeps uses no CPU time, so a budget does not count the
// time the work waits for a CPU on a loaded box.
func TestRunDoesNotCountSleep(t *testing.T) {
	used, done := Run(time.Minute, func() { time.Sleep(200 * time.Millisecond) })
	if !done {
		t.Fatal("Run gave up on work that returns")
	}
	if used > 100*time.Millisecond {
		t.Errorf("Run reported %v of CPU time for a 200ms sleep", used)
	}
}

func TestRunGivesUpAtHang(t *testing.T) {
	block := make(chan struct{})
	defer close(block)
	if _, done := Run(50*time.Millisecond, func() { <-block }); done {
		t.Fatal("Run reported blocked work as done")
	}
}

// recorder stands in for a test so Within's verdict can be read.
type recorder struct {
	testing.TB
	failure string
}

func (r *recorder) Helper()                        {}
func (r *recorder) Logf(string, ...any)            {}
func (r *recorder) Fatalf(format string, a ...any) { r.failure = fmt.Sprintf(format, a...) }

func TestWithinFailsOverBudget(t *testing.T) {
	r := &recorder{TB: t}
	Within(r, 20*time.Millisecond, "spin", func() { spin(100 * time.Millisecond) })
	if !strings.Contains(r.failure, "spin used") || !strings.Contains(r.failure, "over its 20ms budget") {
		t.Errorf("Within over budget reported %q", r.failure)
	}
	r = &recorder{TB: t}
	Within(r, time.Second, "spin", func() { spin(10 * time.Millisecond) })
	if r.failure != "" {
		t.Errorf("Within under budget reported %q", r.failure)
	}
}

func TestHang(t *testing.T) {
	if got := Hang(2 * time.Second); got != time.Minute {
		t.Errorf("Hang(2s) = %v, want 1m", got)
	}
	if got := Hang(10 * time.Second); got != 100*time.Second {
		t.Errorf("Hang(10s) = %v, want 100s", got)
	}
}
