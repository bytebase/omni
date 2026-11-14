package catalog

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

// columnPartitionClauseCases hold what the column, key and partition clauses of
// MariaDB's grammar that the parser reads since the grammar's rows held it to them
// do to the tables: a column's engine-defined options, CONSTRAINT before its
// REFERENCES, NOT NULL ENABLE, a qualified column name, the PARTITION keyword a
// partition definition may leave off, and CONVERT PARTITION and CONVERT TABLE.
// Each case runs its setup and then sql on mariadb:11.8.8 and on the catalog, which
// must answer the same code and print each of tables the same.
var columnPartitionClauseCases = []struct {
	name    string
	setup   []string
	sql     string
	tables  []string
	errCode int
}{
	{
		name:   "a column's reference takes a constraint name",
		setup:  []string{"CREATE TABLE p (id INT PRIMARY KEY)"},
		sql:    "CREATE TABLE c (a INT CONSTRAINT fk1 REFERENCES p (id), b INT CONSTRAINT REFERENCES p (id), d INT CHECK (d > 0) CONSTRAINT fk2 REFERENCES p (id) ON DELETE CASCADE)",
		tables: []string{"c"},
	},
	{
		name:   "NOT NULL ENABLE is NOT NULL",
		sql:    "CREATE TABLE n (a INT NOT NULL ENABLE DEFAULT 1, b INT NOT NULL ENABLE)",
		tables: []string{"n"},
	},
	{
		name:   "INVISIBLE comes before COMMENT",
		sql:    "CREATE TABLE i (a INT COMMENT 'c' INVISIBLE, b INT, c INT AS (b) VIRTUAL COMMENT 'g' INVISIBLE, d INT COMMENT 'd' WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING",
		tables: []string{"i"},
	},
	{
		name:   "a ROW START column takes a key",
		sql:    "CREATE TABLE v (a INT, s TIMESTAMP(6) AS ROW START KEY, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING",
		tables: []string{"v"},
	},
	{
		name:   "a ROW START column takes asrow_attribute",
		sql:    "CREATE TABLE v (a INT, s TIMESTAMP(6) AS ROW START NOT NULL ENABLE UNIQUE COMMENT 'r', e TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING",
		tables: []string{"v"},
	},
	{
		name:    "a ROW START column ALTER adds to a table not versioned",
		setup:   []string{"CREATE TABLE v (a INT PRIMARY KEY)"},
		sql:     "ALTER TABLE v ADD a TIMESTAMP(6) AS ROW START KEY",
		errCode: 4124,
	},
	{
		name:   "ROW START and ROW END ALTER adds with SYSTEM VERSIONING",
		setup:  []string{"CREATE TABLE v (a INT, KEY k (a))"},
		sql:    "ALTER TABLE v ADD s TIMESTAMP(6) AS ROW START, ADD e TIMESTAMP(6) AS ROW END, ADD PERIOD FOR SYSTEM_TIME (s, e), ADD SYSTEM VERSIONING",
		tables: []string{"v"},
	},
	{
		name:    "an engine-defined option InnoDB does not know",
		sql:     "CREATE TABLE o (a INT x = 1)",
		errCode: 1911,
	},
	{
		name:    "an engine-defined option ALTER adds",
		setup:   []string{"CREATE TABLE o (a INT)"},
		sql:     "ALTER TABLE o ADD b INT x = 's'",
		errCode: 1911,
	},
	{
		name:   "engine-defined options under IGNORE_BAD_TABLE_OPTIONS",
		setup:  []string{"SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'"},
		sql:    "CREATE TABLE o (a INT NOT NULL DEFAULT 1 COMMENT 'c' x = 1 INVISIBLE CHECK (a > 0), b INT y = 2 z := 'it''s' Y = w, c INT q = 1 q = DEFAULT, d INT `select` = 0x41 READ_ONLY = 7)",
		tables: []string{"o"},
	},
	{
		name:   "MODIFY gives a column the options it writes",
		setup:  []string{"SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'", "CREATE TABLE o (a INT x = 1, b INT y = 2)"},
		sql:    "ALTER TABLE o MODIFY a INT, CHANGE b b INT k = 3",
		tables: []string{"o"},
	},
	{
		name:   "options outside IGNORE_BAD_TABLE_OPTIONS print in a comment",
		setup:  []string{"SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'", "CREATE TABLE o (a INT x = 1, b INT y = 2 z = 3)", "SET sql_mode = DEFAULT"},
		sql:    "ALTER TABLE o ADD c INT",
		tables: []string{"o"},
	},
	{
		name:   "a column named after its table",
		setup:  []string{"CREATE TABLE q (a INT, b INT, c INT DEFAULT 1)"},
		sql:    "ALTER TABLE q DROP q.b, ALTER q.c DROP DEFAULT, CHANGE .a z INT, ADD q.d INT",
		tables: []string{"q"},
	},
	{
		name:   "partitions without the PARTITION keyword",
		sql:    "CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))",
		tables: []string{"pr"},
	},
	{
		name:   "ADD PARTITION without the keyword",
		setup:  []string{"CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))"},
		sql:    "ALTER TABLE pr ADD PARTITION (p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN (30))",
		tables: []string{"pr"},
	},
	{
		name:   "REORGANIZE PARTITION without the keyword",
		setup:  []string{"CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))"},
		sql:    "ALTER TABLE pr REORGANIZE PARTITION p1 INTO (p1a VALUES LESS THAN (15), p1b VALUES LESS THAN (20))",
		tables: []string{"pr"},
	},
	{
		name:   "CONVERT PARTITION makes a table",
		setup:  []string{"CREATE TABLE pt (a INT PRIMARY KEY, b VARCHAR(10) COMMENT 'bb') PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN (30))"},
		sql:    "ALTER TABLE pt CONVERT PARTITION p1 TO TABLE ct",
		tables: []string{"pt", "ct"},
	},
	{
		name:   "CONVERT TABLE makes a partition",
		setup:  []string{"CREATE TABLE pt (a INT PRIMARY KEY, b VARCHAR(10) COMMENT 'bb') PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", "CREATE TABLE ct (a INT PRIMARY KEY, b VARCHAR(10) COMMENT 'bb')"},
		sql:    "ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20) WITHOUT VALIDATION",
		tables: []string{"pt", "ct"},
	},
	{
		name:    "CONVERT TABLE of another definition",
		setup:   []string{"CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", "CREATE TABLE ct (a INT)"},
		sql:     "ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20)",
		errCode: 1736,
	},
	{
		name:    "CONVERT TABLE of another definition, the clause in a comment",
		setup:   []string{"CREATE TABLE pt (a INT COMMENT '\\n PARTITION BY x', b INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", "CREATE TABLE ct (a INT COMMENT '\\n PARTITION BY x', b VARCHAR(20))"},
		sql:     "ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20)",
		errCode: 1736,
	},
	{
		name:    "CONVERT TABLE of no table",
		setup:   []string{"CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))"},
		sql:     "ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20)",
		errCode: 1146,
	},
	{
		name:    "CONVERT PARTITION of no partition",
		setup:   []string{"CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))"},
		sql:     "ALTER TABLE pt CONVERT PARTITION p9 TO TABLE ct",
		errCode: 1507,
	},
	{
		name:    "CONVERT PARTITION to a table there",
		setup:   []string{"CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))", "CREATE TABLE ct (a INT)"},
		sql:     "ALTER TABLE pt CONVERT PARTITION p1 TO TABLE ct",
		errCode: 1050,
	},
	{
		name:    "CONVERT PARTITION of HASH partitions",
		setup:   []string{"CREATE TABLE pt (a INT) PARTITION BY HASH (a) PARTITIONS 2"},
		sql:     "ALTER TABLE pt CONVERT PARTITION p0 TO TABLE ct",
		errCode: 1512,
	},
	{
		name:    "CONVERT PARTITION of a table with none",
		setup:   []string{"CREATE TABLE pt (a INT)"},
		sql:     "ALTER TABLE pt CONVERT PARTITION p0 TO TABLE ct",
		errCode: 1505,
	},
	{
		name:    "CONVERT TABLE into a table with none",
		setup:   []string{"CREATE TABLE pt (a INT)", "CREATE TABLE ct (a INT)"},
		sql:     "ALTER TABLE pt CONVERT TABLE ct TO PARTITION p0",
		errCode: 1505,
	},
}

// TestColumnPartitionClausesOracle holds the catalog to mariadb:11.8.8 on
// columnPartitionClauseCases.
func TestColumnPartitionClausesOracle(t *testing.T) {
	ctx := context.Background()
	conn := mariadbDatabase(t)
	reset := func() {
		for _, stmt := range []string{"SET sql_mode = DEFAULT", "DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct"} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Cleanup(reset)
	for _, tc := range columnPartitionClauseCases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			cat := New()
			execAlterClauses(t, cat, "CREATE DATABASE d; USE d;")
			for _, stmt := range tc.setup {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("mariadb setup %q: %v", stmt, err)
				}
				execAlterClauses(t, cat, stmt)
			}
			engCode, engMsg := mariadbErr(t, tc.sql, func() error { _, err := conn.ExecContext(ctx, tc.sql); return err }())
			if engCode != tc.errCode {
				t.Fatalf("mariadb answered %d (%s), the case says %d", engCode, engMsg, tc.errCode)
			}
			catCode, catMsg := 0, ""
			if err := execAlterClausesExpectError(t, cat, tc.sql); err != nil {
				var cerr *Error
				if !errors.As(err, &cerr) {
					t.Fatalf("catalog: %v", err)
				}
				catCode, catMsg = cerr.Code, cerr.Message
			}
			if catCode != engCode {
				t.Fatalf("catalog answered %d (%s), mariadb %d (%s)", catCode, catMsg, engCode, engMsg)
			}
			for _, table := range tc.tables {
				eng := showCreateOrNone(t, conn, table)
				omni := cat.ShowCreateTable("d", table)
				if normalizeWhitespace(eng) != normalizeWhitespace(omni) {
					t.Errorf("SHOW CREATE TABLE %s differs:\n--- mariadb ---\n%s\n--- catalog ---\n%s", table, eng, omni)
				}
			}
		})
	}
}

// showCreateOrNone is SHOW CREATE TABLE of a table on the server, or "" when it
// has none of the name (1146).
func showCreateOrNone(t *testing.T, conn *sql.DB, table string) string {
	t.Helper()
	var name, ddl string
	err := conn.QueryRowContext(context.Background(), "SHOW CREATE TABLE "+table).Scan(&name, &ddl)
	if code, _ := mariadbErr(t, "SHOW CREATE TABLE "+table, err); code == 1146 {
		return ""
	}
	return ddl
}
