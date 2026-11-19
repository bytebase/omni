package catalog

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestSimulationCausesOracle holds the refusals testdata/simulate.tsv records to
// pingcap/tidb:v8.5.5 cause by cause: a partition action on a table that has no
// partitions, a view body's star that has no table to expand, a CREATE TABLE
// that lists no column, a NOW()-family DEFAULT the column's type cannot hold,
// and the second catalog round's causes — a partition clause over a column the
// table does not have, a generated column's expression, a TTL config, a storage
// engine or charset name, an unknown placement policy, ADD of AUTO_RANDOM or a
// STORED generated column, a table-level foreign key with no referenced
// column, a column reference inside SET DEFAULT, DROP VIEW over a table, a
// sequence reference to one, an index name the table does not have, a
// duplicated ENUM member, an ON UPDATE a non-temporal column cannot carry, and
// a DEFAULT that calls a function TiDB disallows. Each case is answered by the
// pinned server and by the catalog, and both must agree with the code the case
// states. A case whose statement TiDB takes also compares the schema the two
// store, so CACHE and a repartition are held to the text TiDB prints.
func TestSimulationCausesOracle(t *testing.T) {
	conn, cleanup := startContainer(t)
	defer cleanup()
	ctx := context.Background()
	// A table named a is what most cases leave, and DROP VIEW fails on it, so
	// the reset drops tables; the view cases make v, which no table shares.
	// A cached table cannot be dropped at all (8242, read on v8.5.5), so the
	// CACHE case's table is uncached first, on failure too.
	reset := func(tb testing.TB) {
		for _, stmt := range []string{"ALTER TABLE a NOCACHE", "DROP TABLE IF EXISTS a", "DROP TABLE IF EXISTS s", "DROP TABLE IF EXISTS z", "DROP VIEW IF EXISTS v", "DROP TABLE IF EXISTS `dual`"} {
			if stmt == "ALTER TABLE a NOCACHE" {
				// Only a cached table has anything to uncache, and the table may
				// not exist at all.
				_, _ = conn.db.ExecContext(ctx, stmt)
				continue
			}
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
		// sameSchema compares the catalog's SHOW CREATE TABLE of table (a when
		// empty) with the server's, for a statement TiDB takes that changes the
		// stored definition.
		sameSchema bool
		table      string
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
			name:  "CREATE VIEW whose qualified star names its relation in another schema",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT other.a.* FROM a",
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
		// A partition clause resolves its function's columns before anything
		// else, whichever form carries it: 1054 for a RANGE or HASH expression
		// that names a column the table does not have, and 1488 for one a
		// COLUMNS or KEY column list names.
		{
			name:  "PARTITION BY HASH over a column the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY HASH(a)",
			code:  1054,
		},
		{
			name:  "PARTITION BY RANGE over a column the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))",
			code:  1054,
		},
		{
			name:  "PARTITION BY KEY naming a column the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY(a)",
			code:  1488,
		},
		{
			// The partition function is resolved first, so the interval's own
			// FIRST/LAST requirement (8200) never answers this statement, and
			// the bound a column names is 1054 in 'expression'.
			name:  "INTERVAL partition whose FIRST bound names a column",
			setup: []string{oneTable, "ALTER TABLE a ADD c2 INT"},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(c2) INTERVAL(1) FIRST PARTITION LESS THAN (c) LAST PARTITION LESS THAN (20)",
			code:  1054,
		},
		{
			name:  "INTERVAL partition over a column the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(a) INTERVAL(a) FIRST PARTITION LESS THAN (10) LAST PARTITION LESS THAN (20)",
			code:  1054,
		},
		// A generated column's expression is resolved against the table, and a
		// STORED one cannot be added at all.
		{
			name:  "ADD a generated column whose expression names no column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT AS (a)",
			code:  1054,
		},
		{
			name:  "ADD a generated column whose expression reads a column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT AS (c + 1)",
			code:  0,
		},
		{
			name:  "ADD a STORED generated column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT AS (c) STORED",
			code:  3106,
		},
		// A TTL config resolves its column and holds its type.
		{
			name:  "TTL over a column the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TTL a + INTERVAL 1 DAY",
			code:  1054,
		},
		{
			name:  "TTL over a column of a type it cannot use",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TTL c + INTERVAL 1 DAY",
			code:  8148,
		},
		{
			name:  "TTL over a DATETIME column",
			setup: []string{oneTable, "ALTER TABLE a ADD dt DATETIME"},
			stmt:  "ALTER TABLE a TTL dt + INTERVAL 1 DAY",
			code:  0,
		},
		// DROP VIEW names a view; a table of the name is 1347, and IF EXISTS
		// covers only a name nothing holds.
		{
			name:  "DROP VIEW over a table",
			setup: []string{oneTable},
			stmt:  "DROP VIEW a",
			code:  1347,
		},
		{
			name:  "DROP VIEW IF EXISTS over a table",
			setup: []string{oneTable},
			stmt:  "DROP VIEW IF EXISTS a",
			code:  1347,
		},
		{
			name:  "DROP VIEW IF EXISTS over no object",
			setup: []string{oneTable},
			stmt:  "DROP VIEW IF EXISTS nope",
			code:  0,
		},
		// A sequence a DEFAULT names is resolved before the duplicate column
		// check: a table of the name is 1347, nothing of it 1146.
		{
			name:  "ADD a column whose DEFAULT names a table as a sequence",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT DEFAULT NEXT VALUE FOR a",
			code:  1347,
		},
		{
			name:  "ADD a duplicate column whose DEFAULT names a table as a sequence",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD c INT DEFAULT NEXT VALUE FOR a",
			code:  1347,
		},
		// An index name the table has no index under is 1176 for RENAME and for
		// ALTER INDEX, not the 1061 a duplicate name answers.
		{
			name:  "RENAME KEY of an index the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a RENAME KEY c TO q",
			code:  1176,
		},
		{
			name:  "ALTER INDEX VISIBLE of an index the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER INDEX c VISIBLE",
			code:  1176,
		},
		{
			name:  "ALTER INDEX INVISIBLE of an index the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER INDEX c INVISIBLE",
			code:  1176,
		},
		// AUTO_RANDOM is a CREATE TABLE attribute: no ALTER adds one, and
		// AUTO_RANDOM_BASE needs one to move.
		{
			name:  "ADD a column carrying AUTO_RANDOM",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TINYINT AUTO_RANDOM",
			code:  8216,
		},
		{
			name:  "ADD a duplicate column carrying AUTO_RANDOM",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD c TINYINT AUTO_RANDOM",
			code:  8216,
		},
		{
			name:  "AUTO_RANDOM_BASE on a table without one",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a AUTO_RANDOM_BASE 1",
			code:  8216,
		},
		{
			name:  "AUTO_RANDOM_BASE on a table with an auto-random column",
			setup: []string{"CREATE TABLE a (b BIGINT PRIMARY KEY AUTO_RANDOM)"},
			stmt:  "ALTER TABLE a AUTO_RANDOM_BASE 2",
			code:  0,
		},
		// A table-level FOREIGN KEY names the referenced columns it matches; a
		// column-level REFERENCES is another form and takes the referenced
		// table's key.
		{
			name:  "ADD a FOREIGN KEY with no referenced column list",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES a",
			code:  1239,
		},
		{
			name:  "ADD a FOREIGN KEY with a referenced column list",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES a(b)",
			code:  0,
		},
		// A charset a column definition names is resolved while the statement
		// is prepared: UNICODE is ucs2, which TiDB does not have, even with no
		// table to alter.
		{
			name:  "ADD a UNICODE text column over a table",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TINYTEXT UNICODE",
			code:  1115,
		},
		{
			name: "ADD a UNICODE text column over no table",
			stmt: "ALTER TABLE nope ADD e TINYTEXT UNICODE",
			code: 1115,
		},
		// A placement policy a statement names must exist.
		{
			name:  "PLACEMENT POLICY the catalog does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PLACEMENT POLICY 'binary'",
			code:  8239,
		},
		// Column options the ALTER path used to skip.
		{
			name:  "ADD a column with ON UPDATE CURRENT_TIMESTAMP on a TINYINT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TINYINT ON UPDATE CURRENT_TIMESTAMP",
			code:  1294,
		},
		{
			name:  "ADD a TIMESTAMP column with ON UPDATE CURRENT_TIMESTAMP",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP",
			code:  0,
		},
		{
			name:  "ADD an ENUM column with a duplicated member",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e ENUM('binary', 'binary')",
			code:  1291,
		},
		{
			name:  "ADD an ENUM column with distinct members",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e ENUM('x', 'y')",
			code:  0,
		},
		{
			name:  "ADD a column whose DEFAULT calls a disallowed function",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e TINYINT DEFAULT a()",
			code:  3770,
		},
		{
			name:  "ADD a column whose DEFAULT is NOW() on a DATETIME",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e DATETIME DEFAULT NOW()",
			code:  0,
		},
		// CACHE and a repartition both reach SHOW CREATE TABLE, so the schema
		// the catalog keeps is compared with the server's.
		{
			name:       "ALTER TABLE CACHE",
			setup:      []string{oneTable},
			stmt:       "ALTER TABLE a CACHE",
			code:       0,
			sameSchema: true,
		},
		{
			name:       "ALTER TABLE PARTITION BY KEY with no column",
			setup:      []string{oneTable},
			stmt:       "ALTER TABLE a PARTITION BY KEY()",
			code:       0,
			sameSchema: true,
		},
		{
			name:       "ALTER TABLE PARTITION BY LINEAR KEY with no column",
			setup:      []string{oneTable},
			stmt:       "ALTER TABLE a PARTITION BY LINEAR KEY()",
			code:       0,
			sameSchema: true,
		},
		{
			name:       "ALTER TABLE PARTITION BY HASH",
			setup:      []string{oneTable},
			stmt:       "ALTER TABLE a PARTITION BY HASH(b)",
			code:       0,
			sameSchema: true,
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
			if tc.sameSchema {
				table := tc.table
				if table == "" {
					table = "a"
				}
				engineText := engineShowCreateTable(t, conn.db, ctx, table)
				if mine := c.ShowCreateTable("d", table); mine != engineText {
					t.Fatalf("%s: the catalog keeps\n%s\ntidb keeps\n%s", tc.stmt, mine, engineText)
				}
			}
		})
	}
}

// engineShowCreateTable reads the server's SHOW CREATE TABLE text, which the
// sameSchema cases compare with what the catalog keeps.
func engineShowCreateTable(t *testing.T, db *sql.DB, ctx context.Context, table string) string {
	t.Helper()
	rows, err := db.QueryContext(ctx, "SHOW CREATE TABLE "+quoteIdent(table))
	if err != nil {
		t.Fatalf("SHOW CREATE TABLE %s: %v", table, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("SHOW CREATE TABLE %s: %v", table, err)
	}
	values := make([]sql.NullString, len(columns))
	into := make([]any, len(columns))
	for i := range values {
		into[i] = &values[i]
	}
	if !rows.Next() || len(values) < 2 {
		t.Fatalf("SHOW CREATE TABLE %s: no text", table)
	}
	if err := rows.Scan(into...); err != nil {
		t.Fatalf("SHOW CREATE TABLE %s: %v", table, err)
	}
	return values[1].String
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
