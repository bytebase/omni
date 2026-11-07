package parser

import (
	"testing"
	"time"
)

// TestCompletionPerformance holds completion at the caret to a budget in CPU
// time, so a loaded machine does not fail it.
func TestCompletionPerformance(t *testing.T) {
	// Warm cache
	Collect("SELECT ", 7)

	start := cpuTime(t)
	for i := 0; i < 100; i++ {
		Collect("SELECT ", 7)
	}
	elapsed := cpuSince(t, start)
	avg := elapsed / 100
	t.Logf("avg completion CPU time: %v", avg)
	if avg > 5*time.Millisecond {
		t.Errorf("completion too slow: avg %v of CPU time (want < 5ms)", avg)
	}
}
