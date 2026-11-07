package catalog

import (
	"context"
	"errors"
	"testing"
)

const (
	errPKNeedsPartitionColumns = "A PRIMARY KEY must include all columns in the table's partitioning function"
	errUKNeedsPartitionColumns = "A UNIQUE INDEX must include all columns in the table's partitioning function"
)

// partitionCoverageCase is one statement MariaDB 11.8 holds to the rule that every
// unique key of a partitioned table covers the partitioning columns. setup runs first
// and must succeed; table names the table the statement works on.
type partitionCoverageCase struct {
	name  string
	setup []string
	stmt  string
	table string
	code  int    // 0 when MariaDB accepts the statement
	msg   string // the message for a 1503
}

const coverageHashTable = "CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2"

// partitionCoverageCases are read off MariaDB 11.8.8. The CREATE TABLE path used to check
// only HASH over a bare column and called every key a UNIQUE INDEX; the ALTER TABLE
// and CREATE INDEX paths did not check at all.
var partitionCoverageCases = []partitionCoverageCase{
	{name: "create hash unique", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create hash primary", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create range", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION `p0` VALUES LESS THAN (10), PARTITION `p1` VALUES LESS THAN MAXVALUE)", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create list", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION `p0` VALUES IN (1, 2))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create range columns", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(`a`) (PARTITION `p0` VALUES LESS THAN (10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create list columns", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(`a`) (PARTITION `p0` VALUES IN (1))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create range columns pair", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(`a`,`d`) (PARTITION `p0` VALUES LESS THAN (10, 10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create key", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create key over the primary key", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create key over a promoted unique key", stmt: "CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2", table: "c"},
	{name: "create hash expression", stmt: "CREATE TABLE c (d DATE, b INT, UNIQUE KEY (b)) PARTITION BY HASH(YEAR(d)) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create hash of two columns", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create linear hash", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LINEAR HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create subpartition hash", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create subpartition key", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create primary key named first", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create primary key covers", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b, a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create column primary key", stmt: "CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create column unique", stmt: "CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create serial is promoted", stmt: "CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create first NOT NULL unique key is promoted and covers", stmt: "CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create first NOT NULL unique key is promoted and leaves a column out", stmt: "CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create a nullable unique key is not promoted", stmt: "CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create a NOT NULL unique key is promoted past a nullable one", stmt: "CREATE TABLE c (a INT, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create prefix", stmt: "CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(5), b)) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create prefix short by one", stmt: "CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(19))) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create prefix of the whole column", stmt: "CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(20))) PARTITION BY KEY(a) PARTITIONS 2", table: "c"},
	{name: "create prefix on another column", stmt: "CREATE TABLE c (a VARCHAR(20), b VARCHAR(20), UNIQUE KEY (a, b(3))) PARTITION BY KEY(a) PARTITIONS 2", table: "c"},
	{name: "create covered", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (A, b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c"},

	{name: "add unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique covered", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE KEY u (a, b)", table: "t"},
	{name: "add constraint unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD CONSTRAINT uc UNIQUE (c)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD PRIMARY KEY (b)", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "add constraint primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (c)", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "add primary key covered", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD PRIMARY KEY (a, b)", table: "t"},
	{name: "add unique on a NOT NULL column is promoted", setup: []string{"CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "add column unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD COLUMN d INT UNIQUE", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add column primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD COLUMN d INT PRIMARY KEY", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "add unique then column", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b), ADD COLUMN f INT", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add column then unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD COLUMN f INT, ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "primary key named before unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (c), ADD PRIMARY KEY (b)", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "add two uniques", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (c), ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "duplicate column comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b), ADD COLUMN a INT", table: "t", code: 1060},
	{name: "duplicate key name comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE u (b), ADD UNIQUE KEY u (c)", table: "t", code: 1061},
	{name: "missing drop comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b), DROP INDEX nope", table: "t", code: 1091},
	{name: "create unique index", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ui ON t (c)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create unique index covered", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ui ON t (c, a)", table: "t"},
	{name: "create unique index on a NOT NULL column is promoted", setup: []string{"CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2"}, stmt: "CREATE UNIQUE INDEX ub ON t (b)", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create plain index", setup: []string{coverageHashTable}, stmt: "CREATE INDEX ic ON t (c)", table: "t"},
	{name: "add unique to range", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION `p0` VALUES LESS THAN (10), PARTITION `p1` VALUES LESS THAN MAXVALUE)"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to list columns", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(`a`) (PARTITION `p0` VALUES IN (1))"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to range expression", setup: []string{"CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION `p0` VALUES LESS THAN (2000))"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to range expression covered", setup: []string{"CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION `p0` VALUES LESS THAN (2000))"}, stmt: "ALTER TABLE t ADD UNIQUE (b, d)", table: "t"},
	{name: "add unique to key over the primary key", setup: []string{"CREATE TABLE t (a INT PRIMARY KEY, b INT) PARTITION BY KEY() PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to key over a promoted key", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to key", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to subpartitioned", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (10))"}, stmt: "ALTER TABLE t ADD UNIQUE (a)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add prefix unique", setup: []string{"CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (a(5))", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add prefix unique of the whole column", setup: []string{"CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (a(20))", table: "t"},
	{name: "add unique after remove partitioning", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE t REMOVE PARTITIONING"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t"},
	{name: "repartition", setup: []string{"CREATE TABLE t (a INT, b INT, UNIQUE KEY (b))"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "repartition over a primary key", setup: []string{"CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "repartition again", setup: []string{"CREATE TABLE t (a INT, b INT, c INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE t PARTITION BY HASH(b) PARTITIONS 2", "ALTER TABLE t ADD UNIQUE (b, c)"}, stmt: "ALTER TABLE t PARTITION BY HASH(c) PARTITIONS 2", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "rename a partitioning column", setup: []string{"CREATE TABLE t (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(b) PARTITIONS 2", "ALTER TABLE t RENAME COLUMN b TO b2"}, stmt: "ALTER TABLE t ADD UNIQUE (b2, a)", table: "t"},
	{name: "rename a partitioning column, then leave it out", setup: []string{"CREATE TABLE t (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(b) PARTITIONS 2", "ALTER TABLE t CHANGE b b2 INT"}, stmt: "ALTER TABLE t ADD UNIQUE (a)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
}

// TestPartitionKeyCoverage holds the catalog to MariaDB 11.8.8's verdicts, and checks
// that a refused statement leaves the table as it was: a refused ADD PRIMARY KEY (b)
// leaves b nullable, and a refused ADD UNIQUE (b), ADD COLUMN f adds no f. MariaDB
// renames a partitioning column with the partitioning, so the check follows it.
func TestPartitionKeyCoverage(t *testing.T) {
	for _, tc := range partitionCoverageCases {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			for _, stmt := range tc.setup {
				mustExec(t, c, stmt)
			}
			before := c.ShowCreateTable("testdb", tc.table)
			results, err := c.Exec(tc.stmt, nil)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("got %d results, want 1", len(results))
			}
			gotErr := results[0].Error
			if tc.code == 0 {
				if gotErr != nil {
					t.Fatalf("got %v, want the statement accepted", gotErr)
				}
				return
			}
			var catErr *Error
			if !errors.As(gotErr, &catErr) || catErr.Code != tc.code {
				t.Fatalf("got %v, want error %d", gotErr, tc.code)
			}
			if tc.msg != "" && catErr.Message != tc.msg {
				t.Errorf("message = %q, want %q", catErr.Message, tc.msg)
			}
			if after := c.ShowCreateTable("testdb", tc.table); after != before {
				t.Errorf("the refused statement changed the table:\n--- before ---\n%s\n--- after ---\n%s", before, after)
			}
		})
	}
}

// TestContainerPartitionKeyCoverage reads each case's verdict off MariaDB 11.8.8
// (container test).
func TestContainerPartitionKeyCoverage(t *testing.T) {
	db := mariadbDatabase(t)
	ctx := context.Background()
	for _, tc := range partitionCoverageCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, stmt := range append([]string{"DROP TABLE IF EXISTS c, t"}, tc.setup...) {
				if _, err := db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("setup %q: %v", stmt, err)
				}
			}
			_, err := db.ExecContext(ctx, tc.stmt)
			code, msg := mariadbErr(t, tc.stmt, err)
			if code != tc.code {
				t.Fatalf("MariaDB answered %d %q, want %d", code, msg, tc.code)
			}
			if tc.msg != "" && msg != tc.msg {
				t.Errorf("MariaDB message = %q, want %q", msg, tc.msg)
			}
		})
	}
}

// crossDatabaseRenameCases are ALTER TABLE statements MariaDB 11.8.8 refuses while they
// also rename the table into another database. It keeps d1.t whichever command
// comes first; the catalog used to undo only a rename within the database.
var crossDatabaseRenameCases = []struct {
	setup, stmt string
	code        int
}{
	{"CREATE TABLE d1.t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE d1.t ADD UNIQUE (b), RENAME TO d2.t", 1503},
	{"CREATE TABLE d1.t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE d1.t RENAME TO d2.u, ADD UNIQUE (b)", 1503},
	{"CREATE TABLE d1.t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE d1.t ADD UNIQUE (b), RENAME TO d1.v", 1503},
	{"CREATE TABLE d1.t (a INT, b INT)", "ALTER TABLE d1.t RENAME TO d2.u, ADD COLUMN a INT", 1060},
	{"CREATE TABLE d1.t (a INT, b INT)", "ALTER TABLE d1.t ADD COLUMN a INT, RENAME TO d2.u", 1060},
}

// TestRefusedAlterKeepsTableInPlace pins that a refused ALTER TABLE leaves the table
// in its database under its name, whatever RENAME TO the statement carries.
func TestRefusedAlterKeepsTableInPlace(t *testing.T) {
	for _, tc := range crossDatabaseRenameCases {
		t.Run(tc.stmt, func(t *testing.T) {
			c := New()
			mustExec(t, c, "CREATE DATABASE d1")
			mustExec(t, c, "CREATE DATABASE d2")
			mustExec(t, c, tc.setup)
			results, err := c.Exec(tc.stmt, nil)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			var catErr *Error
			if len(results) != 1 || !errors.As(results[0].Error, &catErr) || catErr.Code != tc.code {
				t.Fatalf("got %+v, want error %d", results, tc.code)
			}
			d1, d2 := c.GetDatabase("d1"), c.GetDatabase("d2")
			if tbl := d1.GetTable("t"); tbl == nil || tbl.Database != d1 || tbl.Name != "t" {
				t.Errorf("d1.t is gone or moved: %+v", tbl)
			}
			if len(d1.Tables) != 1 || len(d2.Tables) != 0 {
				t.Errorf("tables after the refused statement: d1 %d, d2 %d, want 1 and 0", len(d1.Tables), len(d2.Tables))
			}
		})
	}
}

// TestContainerRefusedAlterKeepsTableInPlace reads each case's verdict and the tables
// it leaves off MariaDB 11.8.8 (container test).
func TestContainerRefusedAlterKeepsTableInPlace(t *testing.T) {
	db := mariadbDatabase(t)
	ctx := context.Background()
	for _, tc := range crossDatabaseRenameCases {
		t.Run(tc.stmt, func(t *testing.T) {
			for _, stmt := range []string{"DROP DATABASE IF EXISTS d1", "DROP DATABASE IF EXISTS d2", "CREATE DATABASE d1", "CREATE DATABASE d2", tc.setup} {
				if _, err := db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("setup %q: %v", stmt, err)
				}
			}
			_, err := db.ExecContext(ctx, tc.stmt)
			if code, msg := mariadbErr(t, tc.stmt, err); code != tc.code {
				t.Fatalf("MariaDB answered %d %q, want %d", code, msg, tc.code)
			}
			var tables string
			if err := db.QueryRowContext(ctx, "SELECT GROUP_CONCAT(CONCAT(table_schema, '.', table_name)) FROM information_schema.tables WHERE table_schema IN ('d1', 'd2')").Scan(&tables); err != nil {
				t.Fatal(err)
			}
			if tables != "d1.t" {
				t.Errorf("MariaDB leaves %q, want d1.t", tables)
			}
		})
	}
}
