package parser

import (
	"testing"
	"time"
)

func TestCompletionPerformance(t *testing.T) {
	// Warm cache
	Collect("SELECT ", 7)

	start := time.Now()
	for i := 0; i < 100; i++ {
		Collect("SELECT ", 7)
	}
	elapsed := time.Since(start)
	avg := elapsed / 100
	t.Logf("avg completion time: %v", avg)
	if avg > 5*time.Millisecond {
		t.Errorf("completion too slow: avg %v (want < 5ms)", avg)
	}
}
