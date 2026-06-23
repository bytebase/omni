package catalog

import (
	"strings"
	"testing"
)

// TestSystemTimePartitioningCatalog: the catalog stores SYSTEM_TIME partitioning
// (structural, not byte-exact — MariaDB injects a non-deterministic STARTS
// timestamp and omni's partition renderer is MySQL-flavored) and rejects it on a
// non-versioned table with 4124 (grounded vs 11.8.8).
func TestSystemTimePartitioningCatalog(t *testing.T) {
	t.Run("versioned renders PARTITION BY SYSTEM_TIME", func(t *testing.T) {
		ddl := defineAndShow(t, "t",
			"CREATE TABLE t (x INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME (PARTITION p0 HISTORY, PARTITION pn CURRENT)")
		if !strings.Contains(ddl, "PARTITION BY SYSTEM_TIME") {
			t.Errorf("ShowCreateTable missing SYSTEM_TIME partitioning:\n%s", ddl)
		}
	})
	t.Run("versioned interval form stores interval", func(t *testing.T) {
		ddl := defineAndShow(t, "t",
			"CREATE TABLE t (x INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME INTERVAL 1 MONTH PARTITIONS 12")
		if !strings.Contains(ddl, "PARTITION BY SYSTEM_TIME") || !strings.Contains(ddl, "INTERVAL 1 MONTH") {
			t.Errorf("ShowCreateTable missing SYSTEM_TIME INTERVAL:\n%s", ddl)
		}
	})
	t.Run("non-versioned rejected 4124", func(t *testing.T) {
		err := execErr(t, "CREATE TABLE t (x INT) PARTITION BY SYSTEM_TIME (PARTITION p0 HISTORY, PARTITION pn CURRENT)")
		catErr, ok := err.(*Error)
		if !ok || catErr.Code != 4124 {
			t.Fatalf("want *Error 4124, got %v", err)
		}
	})

	const sv = "CREATE TABLE t (x INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME"
	reject := []struct {
		name string
		ddl  string
		code int
	}{
		{"history_under_hash", "CREATE TABLE t (x INT) PARTITION BY HASH(x) (PARTITION p0 HISTORY)", 4113},
		{"history_under_range", "CREATE TABLE t (x INT) PARTITION BY RANGE(x) (PARTITION p0 HISTORY)", 4113},
		{"no_current", sv + " (PARTITION p0 HISTORY)", 4128},
		{"multiple_current", sv + " (PARTITION p0 CURRENT, PARTITION p1 CURRENT)", 4128},
		{"current_not_last", sv + " (PARTITION pc CURRENT, PARTITION p0 HISTORY)", 4128},
		{"only_current", sv + " (PARTITION pc CURRENT)", 4128},
		{"values_under_system_time", sv + " (PARTITION p0 VALUES LESS THAN (10))", 1480},
	}
	for _, tc := range reject {
		t.Run(tc.name, func(t *testing.T) {
			err := execErr(t, tc.ddl)
			catErr, ok := err.(*Error)
			if !ok {
				t.Fatalf("expected *Error, got %v", err)
			}
			if catErr.Code != tc.code {
				t.Errorf("Code = %d, want %d (message: %q)", catErr.Code, tc.code, catErr.Message)
			}
		})
	}
}
