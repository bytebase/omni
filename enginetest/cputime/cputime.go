// Package cputime holds a test's time budget to the CPU time its work uses,
// not to the clock.
//
// CI runs on a shared box. Under load a single-threaded parse keeps its CPU
// time while its clock time grows: at eight times oversubscription a parse
// that takes 390ms still used 420ms of CPU time and 3.3s of clock time. A
// budget read off the clock therefore failed tests that had not regressed.
// What a budget guards against, work that turns quadratic or worse, multiplies
// the CPU time as much as the clock time, so the CPU time still catches it.
//
// The work runs on its own goroutine, so a regression that never returns
// fails the test at a clock deadline instead of hanging it. That deadline is
// Hang(budget), far past the budget, so load alone does not reach it.
package cputime

import (
	"fmt"
	"syscall"
	"testing"
	"time"
)

// Hang is how long Run and Within wait on the clock for work with the given
// budget before they give up on it: ten budgets, and never under a minute.
func Hang(budget time.Duration) time.Duration {
	return max(10*budget, time.Minute)
}

// Now reports the CPU time, user plus system, that this process has used.
// The process's other goroutines count too, the collector's among them; a
// test that times its work does not run beside other tests.
func Now() time.Duration {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		panic(fmt.Sprintf("getrusage: %v", err))
	}
	// Timeval's fields are int64 on linux/amd64; Usec is int32 on darwin.
	elapsed := func(tv syscall.Timeval) time.Duration {
		return time.Duration(tv.Sec)*time.Second + time.Duration(tv.Usec)*time.Microsecond
	}
	return elapsed(ru.Utime) + elapsed(ru.Stime)
}

// Run runs f on its own goroutine and returns the CPU time the process used
// until f returned. It reports done false, with the CPU time used so far, if
// f has not returned after hang of clock time; f then keeps running.
func Run(hang time.Duration, f func()) (used time.Duration, done bool) {
	start := Now()
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		f()
	}()
	select {
	case <-finished:
		return Now() - start, true
	case <-time.After(hang):
		return Now() - start, false
	}
}

// Within runs f and fails t if f uses more than budget of CPU time, or has
// not returned after Hang(budget) of clock time. what names the work in the
// failure. f runs on another goroutine, so it must not call t.Fatal; it
// records what it found for the caller to check after Within returns.
func Within(t testing.TB, budget time.Duration, what string, f func()) {
	t.Helper()
	used, done := Run(Hang(budget), f)
	if !done {
		t.Fatalf("%s had not returned after %v (%v of CPU time; budget %v)", what, Hang(budget), used, budget)
		return
	}
	t.Logf("%s: %v of CPU time, budget %v", what, used, budget)
	if used > budget {
		t.Fatalf("%s used %v of CPU time, over its %v budget", what, used, budget)
	}
}
