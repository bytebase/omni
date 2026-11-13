package catalog

import (
	"runtime"
	"testing"
)

// TestOversizedDefaultStaysCheap holds integerInRange to a constant cost. A
// default of a few bytes names a number with any exponent, and reading it against
// the column's range once allocated a mantissa as wide as the exponent: about
// 40 MiB for '1e100000000', from an input of thirteen. The engine's answers for
// these are in TestOracle_OceanBaseColumnDefinition.
func TestOversizedDefaultStaysCheap(t *testing.T) {
	cost := func(sql string) uint64 {
		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)
		c, err := LoadSQL("CREATE DATABASE d; USE d")
		if err != nil {
			t.Fatal(err)
		}
		if code := omniCode(c, sql); code != 1067 {
			t.Fatalf("%s: got %d, want 1067", sql, code)
		}
		runtime.ReadMemStats(&after)
		return after.TotalAlloc - before.TotalAlloc
	}
	small := cost("CREATE TABLE t (a TINYINT DEFAULT '1e100')")
	large := cost("CREATE TABLE t (a TINYINT DEFAULT '1e100000000')")
	// The two differ only in the exponent, so the larger must not cost by it. A
	// megabyte is far above the statement's own cost and far below the old 40 MiB.
	if large > small+(1<<20) {
		t.Errorf("a bigger exponent cost %d bytes more than a small one (%d vs %d)", large-small, large, small)
	}
}
