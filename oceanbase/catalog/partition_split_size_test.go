package catalog

import (
	"strings"
	"testing"
)

// TestPartitionSplitAndAutoSize holds a RANGE partitioning's SIZE ('size') and
// PARTITION SIZE ... PARTITIONS AUTO to the catalog: SHOW CREATE TABLE prints
// each, and Diff sees one change between two tables that differ only there.
func TestPartitionSplitAndAutoSize(t *testing.T) {
	build := func(table string) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, table)
		return c
	}
	split := "CREATE TABLE t (a INT PRIMARY KEY) PARTITION BY RANGE (a) SIZE ('%s') (PARTITION p0 VALUES LESS THAN (100))"
	one, two := build(strings.Replace(split, "%s", "1G", 1)), build(strings.Replace(split, "%s", "2G", 1))
	if got := one.ShowCreateTable("test", "t"); !strings.Contains(got, "SIZE ('1G')") {
		t.Errorf("SHOW CREATE TABLE lost SIZE ('1G'):\n%s", got)
	}
	if Diff(one, two).IsEmpty() {
		t.Error("Diff saw no change between SIZE ('1G') and SIZE ('2G')")
	}
	if !Diff(one, one).IsEmpty() {
		t.Error("Diff saw a change between a table and itself")
	}

	auto := build("CREATE TABLE t (a INT PRIMARY KEY) PARTITION BY RANGE (a) PARTITION SIZE '1G' PARTITIONS AUTO")
	if got := auto.ShowCreateTable("test", "t"); !strings.Contains(got, "PARTITION SIZE '1G' PARTITIONS AUTO") {
		t.Errorf("SHOW CREATE TABLE lost the automatic partitioning:\n%s", got)
	}
}
