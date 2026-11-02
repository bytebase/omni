package catalog

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

const (
	errPKNeedsPartitionColumns = "A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered)."
	errUKNeedsPartitionColumns = "A UNIQUE INDEX must include all columns in the table's partitioning function (prefixed columns are not considered)."
	errGIPKOnPartitionedTable  = "This version of MySQL doesn't yet support 'generating invisible primary key for the partitioned tables'"
)

// partitionCoverageCase is one statement MySQL 8.0 holds to the rule that every
// unique key of a partitioned table covers the partitioning columns. setup runs first
// and must succeed; table names the table the statement works on.
type partitionCoverageCase struct {
	name  string
	setup []string
	stmt  string
	table string
	code  int    // 0 when MySQL accepts the statement
	msg   string // the message for a 1503
}

const coverageHashTable = "CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2"

// partitionCoverageCases are read off MySQL 8.0. The CREATE TABLE path used to check
// only HASH over a bare column and called every key a UNIQUE INDEX; the ALTER TABLE
// and CREATE INDEX paths did not check at all.
var partitionCoverageCases = []partitionCoverageCase{
	{name: "create hash unique", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create hash primary", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create range", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create list", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1, 2))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create range columns", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create list columns", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create range columns pair", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(a, d) (PARTITION p0 VALUES LESS THAN (10, 10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create key", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create key over the primary key", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create key over a promoted unique key", stmt: "CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2", table: "c"},
	{name: "create hash expression", stmt: "CREATE TABLE c (d DATE, b INT, UNIQUE KEY (b)) PARTITION BY HASH(YEAR(d)) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create hash of two columns", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create linear hash", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LINEAR HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create subpartition hash", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create subpartition key", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create primary key named first", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create primary key covers", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (d), PRIMARY KEY (b, a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create column primary key", stmt: "CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create column unique", stmt: "CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create serial is promoted", stmt: "CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create first NOT NULL unique key is promoted and covers", stmt: "CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create first NOT NULL unique key is promoted and leaves a column out", stmt: "CREATE TABLE c (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create a nullable unique key is not promoted", stmt: "CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY ub (b), UNIQUE KEY ua (a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create a NOT NULL unique key is promoted past a nullable one", stmt: "CREATE TABLE c (a INT, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create with a generated invisible primary key", setup: []string{"SET SESSION sql_generate_invisible_primary_key = ON"}, stmt: "CREATE TABLE c (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))", table: "c", code: 1235, msg: errGIPKOnPartitionedTable},
	{name: "create hash with a generated invisible primary key", setup: []string{"SET SESSION sql_generate_invisible_primary_key = ON"}, stmt: "CREATE TABLE c (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1235, msg: errGIPKOnPartitionedTable},
	{name: "create with a primary key when invisible keys are generated", setup: []string{"SET SESSION sql_generate_invisible_primary_key = ON"}, stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (a)) PARTITION BY HASH(a) PARTITIONS 2", table: "c"},
	{name: "create functional unique", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY ((a + 1))) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errUKNeedsPartitionColumns},
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
	{name: "add functional unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE ((b + 1))", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "duplicate column comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b), ADD COLUMN a INT", table: "t", code: 1060},
	{name: "duplicate key name comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE u (b), ADD UNIQUE KEY u (c)", table: "t", code: 1061},
	{name: "missing drop comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b), DROP INDEX nope", table: "t", code: 1091},
	{name: "create unique index", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ui ON t (c)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create unique index covered", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ui ON t (c, a)", table: "t"},
	{name: "create unique index on a NOT NULL column is promoted", setup: []string{"CREATE TABLE t (a INT, b INT NOT NULL) PARTITION BY HASH(a) PARTITIONS 2"}, stmt: "CREATE UNIQUE INDEX ub ON t (b)", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "create functional unique index", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ue ON t ((c + 2))", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "create plain index", setup: []string{coverageHashTable}, stmt: "CREATE INDEX ic ON t (c)", table: "t"},
	{name: "add unique to range", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to list columns", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to range expression", setup: []string{"CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to range expression covered", setup: []string{"CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))"}, stmt: "ALTER TABLE t ADD UNIQUE (b, d)", table: "t"},
	{name: "add unique to key over the primary key", setup: []string{"CREATE TABLE t (a INT PRIMARY KEY, b INT) PARTITION BY KEY() PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to key over a promoted key", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to key", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add unique to subpartitioned", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))"}, stmt: "ALTER TABLE t ADD UNIQUE (a)", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add prefix unique", setup: []string{"CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (a(5))", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "add prefix unique of the whole column", setup: []string{"CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (a(20))", table: "t"},
	{name: "add unique after remove partitioning", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE t REMOVE PARTITIONING"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t"},
	{name: "repartition", setup: []string{"CREATE TABLE t (a INT, b INT, UNIQUE KEY (b))"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
	{name: "repartition over a primary key", setup: []string{"CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 1503, msg: errPKNeedsPartitionColumns},
	{name: "repartition again", setup: []string{"CREATE TABLE t (a INT, b INT, c INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE t PARTITION BY HASH(b) PARTITIONS 2", "ALTER TABLE t ADD UNIQUE (b, c)"}, stmt: "ALTER TABLE t PARTITION BY HASH(c) PARTITIONS 2", table: "t", code: 1503, msg: errUKNeedsPartitionColumns},
}

// TestPartitionKeyCoverage holds the catalog to MySQL 8.0's verdicts, and checks that
// a refused statement leaves the table as it was: a refused ADD PRIMARY KEY (b)
// leaves b nullable on MySQL, and a refused ADD UNIQUE (b), ADD COLUMN f adds no f.
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

// TestContainer_PartitionKeyCoverage reads each case's verdict off MySQL 8.0
// (container test).
func TestContainer_PartitionKeyCoverage(t *testing.T) {
	for _, tc := range partitionCoverageCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr, cleanup := startContainer(t)
			defer cleanup()
			for _, stmt := range tc.setup {
				if err := ctr.execSQL(stmt); err != nil {
					t.Fatalf("setup %q: %v", stmt, err)
				}
			}
			err := ctr.execSQL(tc.stmt)
			if tc.code == 0 {
				if err != nil {
					t.Fatalf("MySQL refused %q: %v", tc.stmt, err)
				}
				return
			}
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("MySQL gave no verdict for %q: %v", tc.stmt, err)
			}
			if int(myErr.Number) != tc.code {
				t.Fatalf("MySQL answered %d %q, want %d", myErr.Number, myErr.Message, tc.code)
			}
			if tc.msg != "" && myErr.Message != tc.msg {
				t.Errorf("MySQL message = %q, want %q", myErr.Message, tc.msg)
			}
		})
	}
}
