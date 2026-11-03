package catalog

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

const errClusteredNeedsPartitionColumns = "A CLUSTERED INDEX must include all columns in the table's partitioning function"

func errGlobalIndexNeededFor(name string) string {
	return "Global Index is needed for index '" + name + "', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption"
}

// partitionCoverageCase is one statement TiDB v8.5.5 holds to the rule that every
// unique key of a partitioned table covers the partitioning columns. setup runs first
// and must succeed; table names the table the statement works on.
type partitionCoverageCase struct {
	name  string
	setup []string
	stmt  string
	table string
	code  int // 0 when TiDB accepts the statement
	msg   string
}

const coverageHashTable = "CREATE TABLE t (a INT, b INT, c INT, KEY (b)) PARTITION BY HASH(a) PARTITIONS 2"

// partitionCoverageCases are read off TiDB v8.5.5. The catalog used to check none of
// them.
var partitionCoverageCases = []partitionCoverageCase{
	{name: "create hash unique", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create integer primary key", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errClusteredNeedsPartitionColumns},
	{name: "create clustered primary key", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (b) CLUSTERED) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errClusteredNeedsPartitionColumns},
	{name: "create nonclustered primary key", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (b) NONCLUSTERED) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "create varchar primary key", stmt: "CREATE TABLE c (a INT, b VARCHAR(10), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errClusteredNeedsPartitionColumns},
	{name: "create range", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create list", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1, 2))", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create range columns", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10))", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create list columns", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create range columns pair", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY RANGE COLUMNS(a, d) (PARTITION p0 VALUES LESS THAN (10, 10))", table: "c", code: 8264, msg: errGlobalIndexNeededFor("a")},
	{name: "create key", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create key over the primary key", stmt: "CREATE TABLE c (a INT, b INT, PRIMARY KEY (a), UNIQUE KEY (b)) PARTITION BY KEY() PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create key over a unique key", stmt: "CREATE TABLE c (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2", table: "c"},
	{name: "create hash expression", stmt: "CREATE TABLE c (d DATE, b INT, UNIQUE KEY (b)) PARTITION BY HASH(YEAR(d)) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create hash of two columns", stmt: "CREATE TABLE c (a INT, b INT, d INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a + d) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("a")},
	{name: "create linear hash", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (b)) PARTITION BY LINEAR HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create subpartition is ignored", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (a)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", table: "c"},
	{name: "create row handle checked last", stmt: "CREATE TABLE c (a INT, b INT, d INT, PRIMARY KEY (b), UNIQUE KEY (d)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("d")},
	{name: "create clustered key in order", stmt: "CREATE TABLE c (a INT, b VARCHAR(10), d INT, PRIMARY KEY (b), UNIQUE KEY (d)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 1503, msg: errClusteredNeedsPartitionColumns},
	{name: "create clustered key after a unique key", stmt: "CREATE TABLE c (a INT, b VARCHAR(10), d INT, UNIQUE KEY (d), PRIMARY KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("d")},
	{name: "create nonclustered key in order", stmt: "CREATE TABLE c (a INT, b INT, d INT, PRIMARY KEY (b) NONCLUSTERED, UNIQUE KEY ud (d)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "create table keys before column keys", stmt: "CREATE TABLE c (a INT, b INT UNIQUE, d INT, UNIQUE KEY ud (d)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("ud")},
	{name: "create column primary key after table keys", stmt: "CREATE TABLE c (a INT, b VARCHAR(10) PRIMARY KEY, d INT, UNIQUE KEY ud (d)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("ud")},
	{name: "create column primary key", stmt: "CREATE TABLE c (a INT PRIMARY KEY, b INT) PARTITION BY HASH(b) PARTITIONS 2", table: "c", code: 1503, msg: errClusteredNeedsPartitionColumns},
	{name: "create column unique", stmt: "CREATE TABLE c (a INT UNIQUE, b INT) PARTITION BY HASH(b) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("a")},
	{name: "create generated name", stmt: "CREATE TABLE c (a INT, b INT, KEY b (b), UNIQUE KEY (b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b_2")},
	{name: "create serial", stmt: "CREATE TABLE c (a INT, b SERIAL) PARTITION BY HASH(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "create prefix", stmt: "CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(5), b)) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("a")},
	{name: "create prefix short by one", stmt: "CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(19))) PARTITION BY KEY(a) PARTITIONS 2", table: "c", code: 8264, msg: errGlobalIndexNeededFor("a")},
	{name: "create prefix of the whole column", stmt: "CREATE TABLE c (a VARCHAR(20), b INT, UNIQUE KEY (a(20))) PARTITION BY KEY(a) PARTITIONS 2", table: "c"},
	{name: "create covered", stmt: "CREATE TABLE c (a INT, b INT, UNIQUE KEY (A, b)) PARTITION BY HASH(a) PARTITIONS 2", table: "c"},

	{name: "add unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b_2")},
	{name: "add unique covered", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE KEY u (a, b)", table: "t"},
	{name: "add constraint unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD CONSTRAINT uc UNIQUE (c)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("uc")},
	{name: "add primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD PRIMARY KEY (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "add nonclustered primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD PRIMARY KEY (b) NONCLUSTERED", table: "t", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "add constraint primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (c)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "add primary key covered", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD PRIMARY KEY (a, b)", table: "t"},
	{name: "add unique then column", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (b), ADD COLUMN f INT", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b_2")},
	{name: "add column then unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD COLUMN f INT, ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b_2")},
	{name: "unique before primary key", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (c), ADD PRIMARY KEY (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("c")},
	{name: "primary key before unique", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD PRIMARY KEY (b), ADD UNIQUE (c)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "add two uniques", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (c), ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("c")},
	{name: "the key comes before a later duplicate column", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (c), ADD COLUMN a INT", table: "t", code: 8264, msg: errGlobalIndexNeededFor("c")},
	{name: "an earlier duplicate column comes first", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD COLUMN a INT, ADD UNIQUE (c)", table: "t", code: 1060},
	{name: "the key comes before a missing drop", setup: []string{coverageHashTable}, stmt: "ALTER TABLE t ADD UNIQUE (c), DROP INDEX nope", table: "t", code: 8264, msg: errGlobalIndexNeededFor("c")},
	{name: "create unique index", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ui ON t (c)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("ui")},
	{name: "create unique index covered", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ui ON t (c, a)", table: "t"},
	{name: "create unique index over an expression", setup: []string{coverageHashTable}, stmt: "CREATE UNIQUE INDEX ue ON t ((c + 2))", table: "t", code: 8264, msg: errGlobalIndexNeededFor("ue")},
	{name: "create plain index", setup: []string{coverageHashTable}, stmt: "CREATE INDEX ic ON t (c)", table: "t"},
	{name: "add unique to range", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "add unique to list columns", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) (PARTITION p0 VALUES IN (1))"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "add unique to range expression", setup: []string{"CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "add unique to range expression covered", setup: []string{"CREATE TABLE t (d DATE, b INT) PARTITION BY RANGE(YEAR(d)) (PARTITION p0 VALUES LESS THAN (2000))"}, stmt: "ALTER TABLE t ADD UNIQUE (b, d)", table: "t"},
	{name: "add unique to key over the primary key", setup: []string{"CREATE TABLE t (a INT PRIMARY KEY, b INT) PARTITION BY KEY() PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "add unique to key over a unique key", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY (a)) PARTITION BY KEY() PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "add unique to key", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "add unique to subpartitioned", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))"}, stmt: "ALTER TABLE t ADD UNIQUE (a)", table: "t"},
	{name: "add prefix unique", setup: []string{"CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (a(5))", table: "t", code: 8264, msg: errGlobalIndexNeededFor("a")},
	{name: "add prefix unique of the whole column", setup: []string{"CREATE TABLE t (a VARCHAR(20), b INT) PARTITION BY KEY(a) PARTITIONS 2"}, stmt: "ALTER TABLE t ADD UNIQUE (a(20))", table: "t"},
	{name: "add unique after remove partitioning", setup: []string{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE t REMOVE PARTITIONING"}, stmt: "ALTER TABLE t ADD UNIQUE (b)", table: "t"},
	{name: "repartition", setup: []string{"CREATE TABLE t (a INT, b INT, UNIQUE KEY (b))"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 8264, msg: errGlobalIndexNeededFor("b")},
	{name: "repartition over a clustered primary key", setup: []string{"CREATE TABLE t (a INT, b INT, PRIMARY KEY (b))"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 1503, msg: errClusteredNeedsPartitionColumns},
	{name: "repartition over an added primary key", setup: []string{"CREATE TABLE t (a INT, b INT NOT NULL)", "ALTER TABLE t ADD PRIMARY KEY (b)"}, stmt: "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2", table: "t", code: 8264, msg: errGlobalIndexNeededFor("PRIMARY")},
	{name: "repartition again", setup: []string{"CREATE TABLE t (a INT, b INT, c INT, UNIQUE KEY (a, b)) PARTITION BY HASH(a) PARTITIONS 2", "ALTER TABLE t PARTITION BY HASH(b) PARTITIONS 2", "ALTER TABLE t ADD UNIQUE (b, c)"}, stmt: "ALTER TABLE t PARTITION BY HASH(c) PARTITIONS 2", table: "t", code: 8264, msg: errGlobalIndexNeededFor("a")},
}

// TestPartitionKeyCoverage holds the catalog to TiDB v8.5.5's verdicts, and checks
// that a refused statement leaves the table as it was.
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

// TestContainerPartitionKeyCoverage reads each case's verdict off TiDB v8.5.5
// (container test).
func TestContainerPartitionKeyCoverage(t *testing.T) {
	tc := startTiDBForCatalog(t)
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS omni_partition_coverage",
		"CREATE DATABASE omni_partition_coverage",
		"USE omni_partition_coverage",
	} {
		mustExecTiDB(t, tc, stmt)
	}
	defer tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_partition_coverage")
	for _, cs := range partitionCoverageCases {
		t.Run(cs.name, func(t *testing.T) {
			for _, stmt := range append([]string{"DROP TABLE IF EXISTS c, t"}, cs.setup...) {
				mustExecTiDB(t, tc, stmt)
			}
			_, err := tc.db.ExecContext(tc.ctx, cs.stmt)
			if cs.code == 0 {
				if err != nil {
					t.Fatalf("TiDB refused %q: %v", cs.stmt, err)
				}
				return
			}
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("TiDB gave no verdict for %q: %v", cs.stmt, err)
			}
			if int(myErr.Number) != cs.code {
				t.Fatalf("TiDB answered %d %q, want %d", myErr.Number, myErr.Message, cs.code)
			}
			if cs.msg != "" && myErr.Message != cs.msg {
				t.Errorf("TiDB message = %q, want %q", myErr.Message, cs.msg)
			}
		})
	}
}
