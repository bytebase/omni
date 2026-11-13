package catalog

import "testing"

// TestOracle_OceanBasePartitionNames holds checkPartitionNames, with its order around
// 1500 and 1503, and the 1505 of a
// partition command on a table with no partitions to the engine.
func TestOracle_OceanBasePartitionNames(t *testing.T) {
	const table = "CREATE TABLE a (b INT PRIMARY KEY, c INT)"
	runExecSteps(t, "oceanbase_catalog_partition_names", []struct {
		name  string
		steps []engineStep
	}{
		{"names in a partitioning", []engineStep{
			{table, taken},
			{"ALTER TABLE a PARTITION BY HASH(b + zz)", 1054},
			{"ALTER TABLE a PARTITION BY RANGE COLUMNS(zz)", 1054},
			{"ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(zz)", 1054},
			{"ALTER TABLE a PARTITION BY HASH(c)", 1503},
			{"ALTER TABLE a PARTITION BY HASH(c)(PARTITION p(SUBPARTITION s))", 1500},
			{"ALTER TABLE a PARTITION BY KEY()(PARTITION p(SUBPARTITION s VALUES IN(zz)))", 1500},
			{"ALTER TABLE a PARTITION BY KEY(zz)(PARTITION p(SUBPARTITION s))", 1054},
			{"ALTER TABLE a PARTITION BY LIST(zz)(PARTITION p VALUES IN(c))", 1054},
			{"ALTER TABLE a PARTITION BY LIST(b)(PARTITION p VALUES IN(c))", 1054},
			{"ALTER TABLE a PARTITION BY RANGE(b)(PARTITION p VALUES LESS THAN(c))", 1054},
			{"CREATE TABLE t2 (b INT PRIMARY KEY) PARTITION BY HASH(zz)", 1054},
			{"CREATE TABLE t2 (b INT PRIMARY KEY) PARTITION BY LIST(b)(PARTITION p VALUES IN(c))", 1054},
		}},
		{"a table with no partitions", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD PARTITION (PARTITION p VALUES LESS THAN MAXVALUE)", 1505},
			{"ALTER TABLE a DROP PARTITION p", 1505},
			{"ALTER TABLE a TRUNCATE PARTITION p", 1505},
			{"ALTER TABLE a REORGANIZE PARTITION p INTO (PARTITION q VALUES LESS THAN MAXVALUE)", 1505},
			{"ALTER TABLE a EXCHANGE PARTITION p WITH TABLE a WITHOUT VALIDATION", 1505},
		}},
	})
}
