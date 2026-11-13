package catalog

import (
	"errors"
	"testing"
)

// TestAlterSubpartitionCommands holds the subpartition actions to OceanBase's
// 1505 on a table without partitions (simulate.tsv), and checks that a drop and
// an add change the subpartitions a partition lists.
func TestAlterSubpartitionCommands(t *testing.T) {
	c := setupTestTable(t)
	for _, sql := range []string{
		"ALTER TABLE t1 DROP SUBPARTITION s0",
		"ALTER TABLE t1 TRUNCATE SUBPARTITION s0",
		"ALTER TABLE t1 ALTER SUBPARTITION s0 STORAGE_CACHE_POLICY = 'hot'",
		"ALTER TABLE t1 MODIFY PARTITION p0 ADD SUBPARTITION (SUBPARTITION s9 VALUES LESS THAN (10))",
	} {
		results, _ := c.Exec(sql, &ExecOptions{ContinueOnError: true})
		var ce *Error
		if len(results) != 1 || !errors.As(results[0].Error, &ce) || ce.Code != 1505 {
			t.Errorf("%s: got %v, want 1505", sql, results[0].Error)
		}
	}

	mustExec(t, c, "CREATE TABLE t2 (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY RANGE (b) "+
		"(PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (10)))")
	mustExec(t, c, "ALTER TABLE t2 DROP SUBPARTITION s0")
	mustExec(t, c, "ALTER TABLE t2 MODIFY PARTITION p0 ADD SUBPARTITION (SUBPARTITION s2 VALUES LESS THAN (20))")
	var names []string
	for _, sp := range c.GetDatabase("test").GetTable("t2").Partitioning.Partitions[0].SubPartitions {
		names = append(names, sp.Name)
	}
	if len(names) != 2 || names[0] != "s1" || names[1] != "s2" {
		t.Errorf("p0's subpartitions = %q, want [s1 s2]", names)
	}
}
