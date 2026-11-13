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

	// SIZE follows the subpartition clause, where the parser reads it, so SHOW CREATE
	// TABLE of a subpartitioned table reads back.
	sub := build("CREATE TABLE t (a INT PRIMARY KEY) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) SUBPARTITIONS 2 SIZE ('1G') (PARTITION p0 VALUES LESS THAN (100))")
	readBack := New()
	mustExec(t, readBack, "CREATE DATABASE test")
	readBack.SetCurrentDatabase("test")
	if _, err := readBack.Exec(sub.ShowCreateTable("test", "t"), nil); err != nil {
		t.Errorf("SHOW CREATE TABLE of a subpartitioned table with SIZE does not read back: %v\n%s", err, sub.ShowCreateTable("test", "t"))
	}

	// A migration creates such a table with its partitioning in the CREATE TABLE, and
	// adds no ALTER TABLE ... PARTITION BY, which OceanBase refuses for SIZE.
	empty := New()
	mustExec(t, empty, "CREATE DATABASE test")
	plan := GenerateMigration(empty, one, Diff(empty, one))
	creates, alters := 0, 0
	for _, op := range plan.Ops {
		switch op.Type {
		case OpCreateTable:
			creates++
			if !strings.Contains(op.SQL, "SIZE ('1G')") {
				t.Errorf("the CREATE TABLE lost SIZE ('1G'):\n%s", op.SQL)
			}
		case OpAlterTable:
			alters++
		}
	}
	if creates != 1 || alters != 0 {
		t.Errorf("plan has %d CREATE TABLE and %d ALTER TABLE, want 1 and 0:\n%s", creates, alters, plan.SQL())
	}
	// Changing it on a table that exists is an ALTER OceanBase refuses, and the plan says so.
	plan = GenerateMigration(one, two, Diff(one, two))
	warned := false
	for _, op := range plan.Ops {
		warned = warned || (op.Type == OpAlterTable && op.Warning != "")
	}
	if !warned {
		t.Errorf("a SIZE change carries no warning:\n%s", plan.SQL())
	}
}
