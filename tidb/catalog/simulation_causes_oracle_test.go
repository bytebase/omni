package catalog

import (
	"context"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestSimulationCausesOracle holds the refusals testdata/simulate.tsv records to
// pingcap/tidb:v8.5.5 cause by cause: a partition action on a table that has no
// partitions, a view body's star that has no table to expand, a CREATE TABLE
// that lists no column, and a NOW()-family DEFAULT the column's type cannot
// hold. Each case is answered by the pinned server and by the catalog, and both
// must agree with the code the case states.
func TestSimulationCausesOracle(t *testing.T) {
	conn, cleanup := startContainer(t)
	defer cleanup()
	ctx := context.Background()
	// A table named a is what most cases leave, and DROP VIEW fails on it, so
	// the reset drops tables; the view cases make v, which no table shares.
	reset := func(tb testing.TB) {
		for _, stmt := range []string{"DROP TABLE IF EXISTS a", "DROP TABLE IF EXISTS s", "DROP TABLE IF EXISTS z", "DROP VIEW IF EXISTS v", "DROP TABLE IF EXISTS `dual`"} {
			if _, err := conn.db.ExecContext(ctx, stmt); err != nil {
				tb.Errorf("%s: %v", stmt, err)
			}
		}
	}
	t.Cleanup(func() { reset(t) })

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	partitioned := "CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))"
	// The interval partitions FIRST and LAST name need the INTERVAL clause: on
	// a table that only ranges, TiDB answers 8200 for them.
	interval := "CREATE TABLE a (b INT) PARTITION BY RANGE (b) INTERVAL (10) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)"
	cases := []struct {
		name  string
		setup []string
		stmt  string
		code  int
	}{
		{
			name:  "REMOVE PARTITIONING of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  1505,
		},
		{
			name:  "ADD PARTITION to a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (10))",
			code:  1505,
		},
		{
			name:  "ADD PARTITION PARTITIONS to a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION PARTITIONS 2",
			code:  1505,
		},
		{
			name:  "DROP PARTITION from a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DROP PARTITION p0",
			code:  1505,
		},
		{
			// IF EXISTS covers a missing partition, not a missing partitioning.
			name:  "DROP PARTITION IF EXISTS from a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DROP PARTITION IF EXISTS p0",
			code:  1505,
		},
		{
			name:  "TRUNCATE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TRUNCATE PARTITION ALL",
			code:  1505,
		},
		{
			name:  "COALESCE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COALESCE PARTITION 1",
			code:  1505,
		},
		{
			name:  "REORGANIZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1)",
			code:  1505,
		},
		{
			name:  "EXCHANGE PARTITION with another table on a table with none",
			setup: []string{oneTable, "CREATE TABLE s (b INT)"},
			stmt:  "ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s",
			code:  1505,
		},
		{
			// TiDB answers 1505 here where MySQL 8.0.46 answers 1066.
			name:  "EXCHANGE PARTITION with the table itself",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a",
			code:  1505,
		},
		{
			name:  "ANALYZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ANALYZE PARTITION p0",
			code:  1505,
		},
		{
			name:  "PARTITION ATTRIBUTES of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT",
			code:  1505,
		},
		{
			name:  "LAST PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a LAST PARTITION LESS THAN (10)",
			code:  1505,
		},
		{
			name:  "FIRST PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a FIRST PARTITION LESS THAN (20)",
			code:  1505,
		},
		{
			name:  "REMOVE PARTITIONING of a partitioned table",
			setup: []string{partitioned},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  0,
		},
		{
			name:  "ANALYZE PARTITION of a partitioned table",
			setup: []string{partitioned},
			stmt:  "ALTER TABLE a ANALYZE PARTITION p0",
			code:  0,
		},
		{
			name:  "PARTITION ATTRIBUTES of a partitioned table",
			setup: []string{partitioned},
			stmt:  "ALTER TABLE a PARTITION p0 ATTRIBUTES = DEFAULT",
			code:  0,
		},
		{
			name:  "LAST PARTITION of an interval-partitioned table",
			setup: []string{interval},
			stmt:  "ALTER TABLE a LAST PARTITION LESS THAN (1000)",
			code:  0,
		},
		// A view body's star is resolved before the view's own name: a `*`
		// with no FROM expands from nothing and is 1051, even over a table a,
		// where the name alone would be 1050.
		{
			name:  "CREATE VIEW over a table whose body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW a AS SELECT *",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT *",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose body expands an empty qualified star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT t.*",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose body expands a star over a table",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT t.* FROM a t",
			code:  0,
		},
		{
			name:  "CREATE VIEW whose body counts a star with no table",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT count(*)",
			code:  0,
		},
		// A CREATE TABLE's own definition is read before the name is looked
		// up, so a list with no column is 1113 ahead of the collision's 1050.
		{
			name:  "CREATE TABLE with no columns over a table",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a",
			code:  1113,
		},
		{
			name: "CREATE TABLE with no columns",
			stmt: "CREATE TABLE z",
			code: 1113,
		},
		{
			name: "CREATE TABLE whose only element is a key",
			stmt: "CREATE TABLE z (KEY (b))",
			code: 1113,
		},
		{
			name:  "CREATE VIEW whose derived table expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT * FROM (SELECT *) d",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose CTE body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *) SELECT 1",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose body expands a star over DUAL",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT * FROM DUAL",
			code:  1051,
		},
		{
			// A table really named dual is written quoted, which is not the
			// dummy keyword TiDB expands from nothing.
			name:  "CREATE VIEW whose body expands a star over a quoted dual table",
			setup: []string{"CREATE TABLE `dual` (x INT)"},
			stmt:  "CREATE VIEW v AS SELECT * FROM `dual`",
			code:  0,
		},
		{
			// TiDB resolves every CTE, read or not, and a nearer declaration
			// does not stop it resolving the outer one.
			name:  "CREATE VIEW whose nested CTE shadows an unreferenced one",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM (WITH c AS (SELECT 1 AS x) SELECT x FROM c) d",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose qualified star names no relation of the block",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT missing.* FROM a",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose qualified star names its relation's alias",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT t.* FROM a AS t",
			code:  0,
		},
		{
			name:  "CREATE VIEW whose body reads only DUAL for a column",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT 1 FROM DUAL",
			code:  0,
		},
		{
			name:  "CREATE TABLE LIKE over a table",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE z (LIKE a)",
			code:  0,
		},
		// A NOW()-family DEFAULT is a current-timestamp value, which only a
		// DATETIME or TIMESTAMP column holds. TiDB judges it while it prepares
		// the statement: ahead of the table's lookup (1146) and of a duplicate
		// column (1060).
		{
			name: "ADD a NOW() default over no table",
			stmt: "ALTER TABLE nope ADD e TINYINT DEFAULT NOW()",
			code: 1067,
		},
		{
			name:  "ADD a CURRENT_TIMESTAMP default over a duplicate column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD b TINYINT DEFAULT CURRENT_TIMESTAMP",
			code:  1067,
		},
		{
			name:  "ADD a parenthesized CURRENT_TIMESTAMP default",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TINYINT DEFAULT(CURRENT_TIMESTAMP)",
			code:  1067,
		},
		{
			name:  "ADD a CURRENT_TIMESTAMP default with a precision",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TINYINT DEFAULT CURRENT_TIMESTAMP(1)",
			code:  1067,
		},
		{
			name:  "MODIFY a column to a CURRENT_TIMESTAMP default",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a MODIFY b INT DEFAULT CURRENT_TIMESTAMP",
			code:  1067,
		},
		{
			name:  "ADD a NOW() default to a DATETIME column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e DATETIME DEFAULT NOW()",
			code:  0,
		},
		{
			name: "CREATE TABLE with a NOW() default on an INT column",
			stmt: "CREATE TABLE z (a INT DEFAULT NOW())",
			code: 1067,
		},
		{
			name: "CREATE TABLE with a NOW() default on a DATETIME column",
			stmt: "CREATE TABLE z (a DATETIME DEFAULT NOW())",
			code: 0,
		},
		{
			// A parenthesized NOW()-family call is the same default, and its
			// precision is held to the column's own either way.
			name: "CREATE TABLE with a parenthesized precision past the column's",
			stmt: "CREATE TABLE z (a DATETIME DEFAULT(CURRENT_TIMESTAMP(1)))",
			code: 1067,
		},
		{
			name: "CREATE TABLE with a parenthesized precision at the column's",
			stmt: "CREATE TABLE z (a DATETIME(1) DEFAULT(CURRENT_TIMESTAMP(1)))",
			code: 0,
		},
		{
			// SET DEFAULT is judged against the type the column already has,
			// and only the CURRENT_TIMESTAMP spelling is this rule's:
			// SET DEFAULT (NOW()) is a general expression TiDB refuses as 3770.
			name:  "SET DEFAULT CURRENT_TIMESTAMP on an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER b SET DEFAULT (CURRENT_TIMESTAMP)",
			code:  1067,
		},
		{
			name:  "SET DEFAULT CURRENT_TIMESTAMP on a DATETIME column",
			setup: []string{oneTable, "ALTER TABLE a ADD dt DATETIME"},
			stmt:  "ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP)",
			code:  0,
		},
		{
			// SET DEFAULT names the function, so a NOW()-family alias is a
			// disallowed function rather than the current-timestamp default.
			name:  "SET DEFAULT NOW on an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER b SET DEFAULT (NOW())",
			code:  3770,
		},
		{
			name:  "SET DEFAULT LOCALTIME on an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER b SET DEFAULT (LOCALTIME)",
			code:  3770,
		},
		{
			name:  "SET DEFAULT CURRENT_TIMESTAMP past a DATETIME(1) column",
			setup: []string{oneTable, "ALTER TABLE a ADD dt DATETIME(1)"},
			stmt:  "ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP)",
			code:  1067,
		},
		{
			name:  "SET DEFAULT CURRENT_TIMESTAMP(1) at a DATETIME(1) column",
			setup: []string{oneTable, "ALTER TABLE a ADD dt DATETIME(1)"},
			stmt:  "ALTER TABLE a ALTER dt SET DEFAULT (CURRENT_TIMESTAMP(1))",
			code:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset(t)
			for _, stmt := range tc.setup {
				if _, err := conn.db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			_, err := conn.db.ExecContext(ctx, tc.stmt)
			engineCode, engineMessage := simulationEngineRefusal(t, tc.stmt, err)
			if engineCode != tc.code {
				t.Fatalf("tidb: %s answered %d (%s), want %d", tc.stmt, engineCode, engineMessage, tc.code)
			}

			c := New()
			execAll(t, c, "CREATE DATABASE d; USE d;")
			for _, stmt := range tc.setup {
				execAll(t, c, stmt)
			}
			catalogCode, catalogMessage := simulationCatalogRefusal(t, c, tc.stmt)
			if catalogCode != engineCode {
				t.Fatalf("catalog: %s answered %d (%s), tidb answered %d (%s)",
					tc.stmt, catalogCode, catalogMessage, engineCode, engineMessage)
			}
			if tc.code == 0 && tc.stmt == "ALTER TABLE a REMOVE PARTITIONING" {
				db := c.GetDatabase("d")
				if db == nil || db.GetTable("a") == nil || db.GetTable("a").Partitioning != nil {
					t.Fatalf("catalog kept the table partitioned after %s", tc.stmt)
				}
			}
		})
	}
}

// simulationEngineRefusal reads the server's answer as the error number and
// message, 0 for a statement it took.
func simulationEngineRefusal(t *testing.T, stmt string, err error) (int, string) {
	t.Helper()
	if err == nil {
		return 0, ""
	}
	var refusal *gomysql.MySQLError
	if !errors.As(err, &refusal) {
		t.Fatalf("%s: not a TiDB error: %v", stmt, err)
	}
	return int(refusal.Number), refusal.Message
}

// simulationCatalogRefusal reads the catalog's answer as the error number and
// message, 0 for a statement it took.
func simulationCatalogRefusal(t *testing.T, c *Catalog, stmt string) (int, string) {
	t.Helper()
	results, err := c.Exec(stmt, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", stmt, err)
	}
	if len(results) != 1 {
		t.Fatalf("Exec(%q): %d results, want 1", stmt, len(results))
	}
	if results[0].Error == nil {
		return 0, ""
	}
	var refusal *Error
	if !errors.As(results[0].Error, &refusal) {
		t.Fatalf("%s: not a catalog error: %v", stmt, results[0].Error)
	}
	return refusal.Code, refusal.Message
}
