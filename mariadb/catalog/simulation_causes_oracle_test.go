package catalog

import (
	"context"
	"testing"
)

// TestSimulationCausesOracle holds three refusals that simulate.tsv records to
// mariadb:11.8.8: the object an ALTER or DROP VIEW names, a key built on a
// column the table does not have, and a partition action on a table that has
// no partitions.
func TestSimulationCausesOracle(t *testing.T) {
	conn := mariadbDatabase(t)
	ctx := context.Background()
	reset := func() {
		for _, stmt := range []string{"DROP VIEW IF EXISTS a", "DROP TABLE IF EXISTS a", "DROP SEQUENCE IF EXISTS a", "DROP TABLE IF EXISTS p"} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
	}
	t.Cleanup(reset)

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	cases := []struct {
		name  string
		setup []string
		stmt  string
		code  int
	}{
		{
			name:  "ALTER VIEW over a table",
			setup: []string{oneTable},
			stmt:  "ALTER VIEW a AS SELECT a",
			code:  1347,
		},
		{
			name: "ALTER VIEW over no object",
			stmt: "ALTER VIEW a AS SELECT a",
			code: 1146,
		},
		{
			name:  "DROP VIEW over a table",
			setup: []string{oneTable},
			stmt:  "DROP VIEW a",
			code:  4092,
		},
		{
			name: "DROP VIEW over no object",
			stmt: "DROP VIEW a",
			code: 4092,
		},
		{
			name:  "ALTER VIEW over a sequence",
			setup: []string{"CREATE SEQUENCE a"},
			stmt:  "ALTER VIEW a AS SELECT a",
			code:  1347,
		},
		{
			name:  "ADD KEY over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(a)",
			code:  1072,
		},
		{
			name:  "ADD FULLTEXT over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(a)",
			code:  1072,
		},
		{
			name:  "ADD SPATIAL over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD SPATIAL(a)",
			code:  1072,
		},
		{
			name:  "ADD UNIQUE over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD UNIQUE(a)",
			code:  1072,
		},
		{
			name:  "ADD FOREIGN KEY over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a",
			code:  1072,
		},
		{
			name:  "ADD KEY of a name the table has over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY idx_d(missing)",
			code:  1072,
		},
		{
			name:  "ADD PRIMARY KEY over a missing column on a table with one",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PRIMARY KEY(missing)",
			code:  1072,
		},
		{
			name:  "ADD KEY before the column it names",
			setup: []string{"CREATE TABLE a (b INT)"},
			stmt:  "ALTER TABLE a ADD KEY k(missing), ADD COLUMN missing INT",
			code:  0,
		},
		{
			name:  "REMOVE PARTITIONING of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  1505,
		},
		{
			name:  "FORCE REMOVE PARTITIONING of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a FORCE REMOVE PARTITIONING",
			code:  1505,
		},
		{
			name:  "REBUILD PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REBUILD PARTITION ALL",
			code:  1505,
		},
		{
			name:  "OPTIMIZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a OPTIMIZE PARTITION ALL",
			code:  1505,
		},
		{
			name:  "ANALYZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ANALYZE PARTITION ALL",
			code:  1505,
		},
		{
			name:  "CHECK PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CHECK PARTITION ALL",
			code:  1505,
		},
		{
			name:  "REPAIR PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REPAIR PARTITION ALL",
			code:  1505,
		},
		{
			name:  "REMOVE PARTITIONING of a partitioned table",
			setup: []string{"CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2"},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			for _, stmt := range tc.setup {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			_, err := conn.ExecContext(ctx, tc.stmt)
			engineCode, _ := mariadbErr(t, tc.stmt, err)
			if engineCode != tc.code {
				t.Fatalf("mariadb: %s answered %d, want %d", tc.stmt, engineCode, tc.code)
			}

			c := New()
			execAlterClauses(t, c, "CREATE DATABASE d; USE d;")
			for _, stmt := range tc.setup {
				execAlterClauses(t, c, stmt)
			}
			catalogCode, _ := catalogErr(t, c, tc.stmt)
			if catalogCode != engineCode {
				t.Fatalf("catalog: %s answered %d, mariadb answered %d", tc.stmt, catalogCode, engineCode)
			}
			if tc.code == 0 {
				if tbl := c.GetDatabase("d").GetTable("a"); tbl == nil || tbl.Partitioning != nil {
					t.Fatalf("catalog kept the table partitioned after %s", tc.stmt)
				}
			}
		})
	}
}
