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
}
