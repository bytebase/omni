package review

import (
	"syscall"
	"testing"
	"time"
)

// cpuTime reports the CPU time, user plus system, that this process has used.
//
// The performance budgets in this package measure their work with it rather
// than with the clock. The work is single-threaded, so its CPU time is the
// same on an idle machine and on one running eight other jobs, while its
// wall-clock time grows with the load: at eight times oversubscription a parse
// that takes 390ms still took 420ms of CPU time and 3.3s of clock time, which
// is how these budgets failed on the shared box CI runs on. The figures the
// budgets cite are clock time measured idle, where the two differ only by the
// collector's parallel work.
//
// Every engine keeps its own copy, as no engine imports another (AGENTS.md).
func cpuTime(t *testing.T) time.Duration {
	t.Helper()
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		t.Fatalf("Getrusage: %v", err)
	}
	// Timeval's fields are int64 on linux/amd64; Usec is int32 on darwin.
	elapsed := func(tv syscall.Timeval) time.Duration {
		return time.Duration(tv.Sec)*time.Second + time.Duration(tv.Usec)*time.Microsecond
	}
	return elapsed(ru.Utime) + elapsed(ru.Stime)
}

// cpuSince returns the CPU time used since start, a reading from cpuTime.
func cpuSince(t *testing.T, start time.Duration) time.Duration {
	t.Helper()
	return cpuTime(t) - start
}
