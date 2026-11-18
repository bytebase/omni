package catalog

import (
	"context"
	"strings"
	"testing"
)

// TestSimulationCausesOracle holds the refusals testdata/simulate.tsv records to
// mariadb:11.8.8 cause by cause: the object an ALTER or DROP VIEW names, a key
// built on a column the table does not have, a partition action on a table that
// has no partitions, a FULLTEXT key over a column whose type it cannot index,
// a REFERENCES clause with no column list of its own, and an ENGINE the pinned
// server does not ship.
func TestSimulationCausesOracle(t *testing.T) {
	conn := mariadbDatabase(t)
	ctx := context.Background()
	reset := func() {
		for _, stmt := range []string{
			"DROP VIEW IF EXISTS a", "DROP TABLE IF EXISTS a",
			"DROP SEQUENCE IF EXISTS a", "DROP TABLE IF EXISTS p",
			"DROP TABLE IF EXISTS x", "DROP TABLE IF EXISTS z",
		} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
	}
	t.Cleanup(reset)

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	// The same table with a column a, whose only key leads with d.
	intColumnTable := "CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))"
	cases := []struct {
		name  string
		setup []string
		stmt  string
		code  int
		// refs is the REFERENCES clause a taken foreign key must hold, on
		// both sides: MariaDB prints the referenced columns a statement left
		// out, and so must the catalog.
		refs string
		// unpartitioned holds the case whose statement leaves table a
		// unpartitioned, which both sides must then show.
		unpartitioned bool
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
			name:  "ADD FULLTEXT WITH PARSER over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a",
			code:  1128,
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
			name:          "REMOVE PARTITIONING of a partitioned table",
			setup:         []string{"CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2"},
			stmt:          "ALTER TABLE a REMOVE PARTITIONING",
			code:          0,
			unpartitioned: true,
		},
		{
			// A FULLTEXT key takes the character types only: every other type
			// is 1283, measured on mariadb:11.8.8, BLOB and VARBINARY included.
			name:  "ADD FULLTEXT over an INT column",
			setup: []string{intColumnTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(a)",
			code:  1283,
		},
		{
			name:  "ADD FULLTEXT KEY over a BLOB column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c BLOB)"},
			stmt:  "ALTER TABLE a ADD FULLTEXT KEY(c)",
			code:  1283,
		},
		{
			name:  "CREATE FULLTEXT INDEX over a VARBINARY column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c VARBINARY(10))"},
			stmt:  "CREATE FULLTEXT INDEX i ON a(c)",
			code:  1283,
		},
		{
			name:  "ADD FULLTEXT over a TEXT column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c TEXT)"},
			stmt:  "ALTER TABLE a ADD FULLTEXT(c)",
			code:  0,
		},
		{
			// A REFERENCES clause with no referenced column list takes the
			// key's own column names: a column no key of the referenced table
			// leads with cannot form the foreign key.
			name:  "ADD COLUMN with REFERENCES over keys that do not lead with it",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a INT REFERENCES a",
			code:  1005,
		},
		{
			name:  "ADD FOREIGN KEY with REFERENCES and no column list over a column no key leads with",
			setup: []string{intColumnTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a",
			code:  1005,
		},
		{
			name:  "ADD COLUMN with REFERENCES the referenced table indexes by its name",
			setup: []string{oneTable, "CREATE TABLE z (q INT UNIQUE)"},
			stmt:  "ALTER TABLE a ADD q INT REFERENCES z",
			code:  0,
			refs:  "REFERENCES `z` (`q`)",
		},
		{
			name:  "ADD COLUMN with REFERENCES the same table's primary key",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD q INT REFERENCES a(b)",
			code:  0,
			refs:  "REFERENCES `a` (`b`)",
		},
		{
			// An ENGINE option naming an engine the pinned server does not
			// ship is 1286, before the table is looked up, before a name
			// collides and before a column-less definition.
			name: "CREATE TABLE with an unknown engine",
			stmt: "CREATE TABLE x (a INT) ENGINE zzz",
			code: 1286,
		},
		{
			name:  "CREATE TABLE with an unknown engine over an existing name",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a ENGINE zzz",
			code:  1286,
		},
		{
			name: "CREATE TABLE with an unknown engine and no columns",
			stmt: "CREATE TABLE x ENGINE zzz",
			code: 1286,
		},
		{
			name: "ALTER TABLE with an unknown engine over no table",
			stmt: "ALTER TABLE x ENGINE zzz",
			code: 1286,
		},
		{
			name:  "ALTER TABLE with an unknown engine over an existing table",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE zzz",
			code:  1286,
		},
		{
			name: "CREATE TABLE with the Aria engine",
			stmt: "CREATE TABLE x (a INT) ENGINE=Aria",
			code: 0,
		},
		{
			name:  "ALTER TABLE to MyISAM",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE=MyISAM",
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
			if tc.refs != "" {
				want := tc.refs
				var name, create string
				if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&name, &create); err != nil {
					t.Fatalf("SHOW CREATE TABLE a: %v", err)
				}
				if !strings.Contains(create, want) {
					t.Fatalf("mariadb stored\n%s\nwithout %s", create, want)
				}
				if got := c.ShowCreateTable("d", "a"); !strings.Contains(got, want) {
					t.Fatalf("catalog stored\n%s\nwithout %s", got, want)
				}
			}
			if tc.unpartitioned {
				if tbl := c.GetDatabase("d").GetTable("a"); tbl == nil || tbl.Partitioning != nil {
					t.Fatalf("catalog kept the table partitioned after %s", tc.stmt)
				}
			}
		})
	}
}
