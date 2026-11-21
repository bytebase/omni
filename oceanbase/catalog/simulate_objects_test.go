package catalog

import "testing"

// TestAlterViewBehavior covers OceanBase 4.3.5.6's ALTER TABLE / CREATE INDEX /
// DROP INDEX semantics when the target relation is a view (objects@view).
func TestAlterViewBehavior(t *testing.T) {
	const setup = `CREATE DATABASE d DEFAULT CHARSET=utf8mb4;
USE d;
CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
CREATE VIEW a AS SELECT b, a, d FROM t;`

	cases := []struct {
		sql  string
		want int
	}{
		{"ALTER TABLE a COMMENT 's'", 0},
		{"ALTER TABLE a", 1347},
		{"ALTER TABLE a ADD a TINYINT", 1347},
		{"ALTER TABLE a ADD a TINYINT DEFAULT 's'", 1347},
		{"ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP", 1067},
		{"ALTER TABLE a ADD(a TINYINT, a TINYINT)", 1060},
		{"ALTER TABLE a ADD a INT, ADD a INT", 1060},
		{"ALTER TABLE a ADD PRIMARY KEY(a)", 1068},
		{"ALTER TABLE a ADD CONSTRAINT FOREIGN KEY(a) REFERENCES a(a)", 1215},
		{"ALTER TABLE a ADD CHECK(a)", 3812},
		{"ALTER TABLE a ALTER a SET DEFAULT(a)", 1235},
		{"ALTER TABLE a ALTER INDEX a STORAGE_CACHE_POLICY(GLOBAL = 's')", 1235},
		{"ALTER TABLE a ALTER INDEX a VISIBLE", 1347},
		{"ALTER TABLE a TABLEGROUP a", 1347},
		{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE)", 1347},
		{"ALTER TABLE a PARTITION BY KEY()", 1488},
		{"ALTER TABLE a PARTITION BY RANGE()", 1235},
		{"ALTER TABLE a PARTITION BY HASH(a)", 1503},
		{"ALTER TABLE a ADD PARTITION(PARTITION p1 VALUES LESS THAN MAXVALUE)", 1505},
		{"ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS)", 1235},
		{"ALTER TABLE a DROP CHECK a", 3821},
		{"ALTER TABLE a DROP CONSTRAINT a", 3940},
		{"ALTER TABLE a DROP SUBPARTITION", 1091},
		{"CREATE INDEX a ON a(a)", 1347},
		{"CREATE FULLTEXT INDEX a ON a(a)", 1283},
		{"CREATE SPATIAL INDEX a ON a(a)", 1687},
		{"CREATE VECTOR INDEX a ON a(a)", 7601},
		{"DROP INDEX a ON a", 1091},
	}

	for _, tc := range cases {
		c, err := LoadSQL(setup)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		if got := omniCode(c, tc.sql); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.want)
		}
	}
}

// TestAlterPartitionedBehavior covers OceanBase 4.3.5.6's DDL validations and
// generated subpartition numbering on a RANGE-partitioned table (objects@partitioned).
func TestAlterPartitionedBehavior(t *testing.T) {
	const setup = `CREATE DATABASE d DEFAULT CHARSET=utf8mb4;
USE d;
CREATE TABLE a (b INT, a INT, d INT, PRIMARY KEY (b, a), KEY idx_d (d))
  PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (100), PARTITION a VALUES LESS THAN (200));`

	cases := []struct {
		sql  string
		want int
	}{
		{"ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE)", 1235},
		{"ALTER TABLE a PARTITION BY RANGE()", 1235},
		{"ALTER TABLE a PARTITION BY RANGE COLUMNS()", 1235},
		{"ALTER TABLE a PARTITION BY RANGE() SIZE('s')", 1064},
		{"ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE)", 1517},
		{"ALTER TABLE a ADD PARTITION(PARTITION p2 VALUES IN(1))", 1149},
		{"ALTER TABLE a ADD PARTITION(PARTITION p2 VALUES LESS THAN(a))", 1054},
		{"ALTER TABLE a ADD PARTITION(PARTITION p2 VALUES LESS THAN(a, a))", 1653},
		{"ALTER TABLE a ADD PARTITION(PARTITION p2 VALUES LESS THAN(300))", 0},
		{"ALTER TABLE a DROP PARTITION a, a", 1517},
		{"ALTER TABLE a DROP PARTITION a, zzz", 1508},
		{"ALTER TABLE a DROP PARTITION missing_p", 1507},
		{"ALTER TABLE a DROP PARTITION a", 0},
		{"ALTER TABLE a TRUNCATE PARTITION a, zzz", 1507},
		{"ALTER TABLE a TRUNCATE PARTITION a", 0},
		{"ALTER TABLE a TRUNCATE PARTITION ALL", 0},
		{"ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION s0 VALUES LESS THAN MAXVALUE)", 5875},
		{"ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION s0 VALUES IN(1))", 5875},
		{"ALTER TABLE a DROP SUBPARTITION a", 5875},
		{"ALTER TABLE a TRUNCATE SUBPARTITION a", 5875},
		{"ALTER TABLE a ALTER PARTITION a STORAGE_CACHE_POLICY = 's'", 1235},
		{"ALTER TABLE a ALTER SUBPARTITION a STORAGE_CACHE_POLICY = 's'", 1235},
		{"ALTER TABLE a EXCHANGE PARTITION a WITH TABLE a WITHOUT VALIDATION", 1066},
		{"ALTER TABLE a DROP COLUMN a", 3855},
		{"ALTER TABLE a RENAME COLUMN a TO a", 3855},
	}

	for _, tc := range cases {
		c, err := LoadSQL(setup)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		if got := omniCode(c, tc.sql); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.want)
		}
	}

	// Repartitioning a 2-partition table into RANGE/LIST subpartitioning with no
	// explicit subpartitions starts generated subpartition names at P8193.
	c, err := LoadSQL(setup)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if got := omniCode(c, "ALTER TABLE a PARTITION BY KEY() PARTITIONS 1 SUBPARTITION BY RANGE(a)"); got != 0 {
		t.Fatalf("repartition failed with code %d", got)
	}
	layout := subpartitionValueLayout(c.GetDatabase("d").GetTable("a"))
	if len(layout) != 1 || layout[0] != "p0:P8193<MAXVALUE>" {
		t.Errorf("subpartition layout = %v, want [p0:P8193<MAXVALUE>]", layout)
	}
}

// TestAlterForeignKeyTableBehavior covers OceanBase 4.3.5.6's DDL validations
// on a table with a foreign key constraint (objects@fk).
func TestAlterForeignKeyTableBehavior(t *testing.T) {
	const setup = `CREATE DATABASE d DEFAULT CHARSET=utf8mb4;
USE d;
CREATE TABLE parent (b INT PRIMARY KEY);
CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d), CONSTRAINT a FOREIGN KEY (a) REFERENCES parent (b), CONSTRAINT chk_d CHECK (d > 0));`

	cases := []struct {
		sql  string
		want int
	}{
		{"ALTER TABLE a DROP COLUMN a", 5320},
		{"ALTER TABLE a MODIFY COLUMN a TINYINT", 5320},
		{"ALTER TABLE a CHANGE COLUMN a a TINYINT", 5320},
		{"ALTER TABLE a MODIFY COLUMN a INT NOT NULL", 0},
		{"ALTER TABLE a MODIFY COLUMN a INT COMMENT 'ok'", 0},
		{"ALTER TABLE a ALTER CONSTRAINT a NOT ENFORCED", 3941},
		{"ALTER TABLE a ALTER CONSTRAINT a ENFORCED", 3941},
		{"ALTER TABLE a ALTER CONSTRAINT chk_d NOT ENFORCED", 0},
		{"ALTER TABLE a ALTER CHECK a ENFORCED", 3821},
		{"ALTER TABLE a ALTER INDEX a VISIBLE", 1146},
		{"ALTER TABLE a ALTER INDEX a INVISIBLE", 1146},
		{"ALTER TABLE a ALTER INDEX idx_d INVISIBLE", 0},
		{"DROP INDEX a ON a", 1146},
		{"ALTER TABLE a DROP INDEX a", 1091},
	}

	for _, tc := range cases {
		c, err := LoadSQL(setup)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		if got := omniCode(c, tc.sql); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.want)
		}
	}
}
