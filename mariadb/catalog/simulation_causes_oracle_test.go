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
// a REFERENCES clause with no column list of its own, an ENGINE the pinned
// server does not ship, a view or CTAS body naming a column no relation has, a
// SPATIAL key over a column whose type it cannot index, and the table options
// MariaDB stores differently from the text it was given: ROW_FORMAT DEFAULT,
// TABLESPACE, SERIAL, CHAR BYTE, CHAR COLLATE BINARY, and MAX_ROWS.
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
		// absent names a table both sides must leave absent, for a statement a
		// missing table takes silently.
		absent string
		// contains and omits are text both SHOW CREATE TABLE texts must hold
		// and must not hold.
		contains []string
		omits    []string
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
			// The columns a key and its reference name count the same.
			name:  "ADD COLUMN with REFERENCES whose columns do not count the same",
			setup: []string{oneTable, "CREATE TABLE z (q INT UNIQUE)"},
			stmt:  "ALTER TABLE a ADD q INT REFERENCES z(q, q)",
			code:  1239,
		},
		{
			// A backing index named after the constraint keeps a name the
			// table's indexes already hold from being taken twice.
			name:  "ADD COLUMN with a named REFERENCES whose name an index holds",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY k (c))", "CREATE TABLE p (id INT PRIMARY KEY)"},
			stmt:  "ALTER TABLE a ADD q INT CONSTRAINT k REFERENCES p(id)",
			code:  1061,
		},
		{
			// The parser a FULLTEXT key names is looked up before its columns
			// are judged.
			name:  "CREATE FULLTEXT INDEX over an INT column with a parser",
			setup: []string{intColumnTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a(a) WITH PARSER p",
			code:  1128,
		},
		{
			name: "CREATE TABLE with a FULLTEXT key over an INT column and a parser",
			stmt: "CREATE TABLE x (a INT, FULLTEXT KEY x (a) WITH PARSER p)",
			code: 1128,
		},
		{
			// A partition definition's ENGINE is judged with the table's.
			name:  "ADD PARTITION with an unknown engine",
			setup: []string{"CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2"},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p ENGINE=zzz)",
			code:  1286,
		},
		{
			name: "CREATE TABLE with an unknown engine in a partition definition",
			stmt: "CREATE TABLE x (a INT) PARTITION BY HASH(a) (PARTITION p ENGINE=zzz)",
			code: 1286,
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
			name: "CREATE VIEW over a missing column",
			stmt: "CREATE VIEW a AS SELECT a",
			code: 1054,
		},
		{
			name:  "CREATE VIEW over a table with a missing column",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW a AS SELECT a",
			code:  1054,
		},
		{
			name:  "CREATE VIEW over a table with a resolved body",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW a AS SELECT b FROM a",
			code:  1050,
		},
		{
			name: "CREATE VIEW over a missing column in a function",
			stmt: "CREATE VIEW a AS SELECT ABS(a)",
			code: 1054,
		},
		{
			name: "CREATE VIEW over a missing column in a set operation",
			stmt: "CREATE VIEW a AS SELECT 1 UNION SELECT a",
			code: 1054,
		},
		{
			name: "CREATE VIEW over a missing column and an alias-bearing clause",
			stmt: "CREATE VIEW a AS SELECT a HAVING 1",
			code: 1054,
		},
		{
			name: "CREATE VIEW over a missing column in WHERE",
			stmt: "CREATE VIEW a AS SELECT 1 WHERE a",
			code: 1054,
		},
		{
			name: "CREATE VIEW over a nested select alias",
			stmt: "CREATE VIEW a AS SELECT (SELECT 1 AS x HAVING x)",
			code: 0,
		},
		{
			name: "CREATE TABLE ... SELECT over a missing column",
			stmt: "CREATE TABLE a SELECT a",
			code: 1054,
		},
		{
			name:  "CREATE TABLE ... SELECT over a table with a missing column",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a SELECT a",
			code:  1050,
		},
		{
			name:  "ADD SPATIAL over an INT column",
			setup: []string{"CREATE TABLE a (g INT)"},
			stmt:  "ALTER TABLE a ADD SPATIAL(g)",
			code:  1210,
		},
		{
			name:  "ADD SPATIAL over a NOT NULL INT column",
			setup: []string{"CREATE TABLE a (g INT NOT NULL)"},
			stmt:  "ALTER TABLE a ADD SPATIAL(g)",
			code:  1210,
		},
		{
			name:  "ADD SPATIAL over a nullable POINT column",
			setup: []string{"CREATE TABLE a (g POINT)"},
			stmt:  "ALTER TABLE a ADD SPATIAL(g)",
			code:  1252,
		},
		{
			name:  "ADD SPATIAL over a NOT NULL POINT column",
			setup: []string{"CREATE TABLE a (g POINT NOT NULL)"},
			stmt:  "ALTER TABLE a ADD SPATIAL(g)",
			code:  0,
		},
		{
			name:  "ALTER TABLE ROW_FORMAT DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ROW_FORMAT DEFAULT",
			code:  0,
			omits: []string{"ROW_FORMAT="},
		},
		{
			name:  "ALTER TABLE TABLESPACE",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TABLESPACE a",
			code:  0,
			omits: []string{"TABLESPACE"},
		},
		{
			name:     "ALTER TABLE ADD SERIAL",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD a SERIAL",
			code:     0,
			contains: []string{"`a` bigint(20) unsigned NOT NULL AUTO_INCREMENT", "UNIQUE KEY `a` (`a`)"},
		},
		{
			name:     "ALTER TABLE ADD CHAR BYTE",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD a CHAR BYTE",
			code:     0,
			contains: []string{"`a` binary(1)"},
		},
		{
			name:     "ALTER TABLE ADD CHAR COLLATE BINARY",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD a CHAR COLLATE BINARY",
			code:     0,
			contains: []string{"`a` binary(1)"},
		},
		{
			name:     "ALTER TABLE MAX_ROWS above 32 bits",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a MAX_ROWS 9223372036854775808",
			code:     0,
			contains: []string{"MAX_ROWS=4294967295"},
		},
		{
			name:     "ALTER TABLE MAX_ROWS in exponent notation above one",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a MAX_ROWS 2e3",
			code:     0,
			contains: []string{"MAX_ROWS=2"},
		},
		{
			name:     "ALTER TABLE MIN_ROWS in exponent notation above one",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a MIN_ROWS 3e2",
			code:     0,
			contains: []string{"MIN_ROWS=3"},
		},
		{
			name:     "ALTER TABLE MAX_ROWS with a fraction",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a MAX_ROWS 1.5",
			code:     0,
			contains: []string{"MAX_ROWS=1"},
		},
		{
			name:     "ALTER TABLE MAX_ROWS in exponent notation",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a MAX_ROWS 1e0",
			code:     0,
			contains: []string{"MAX_ROWS=1"},
		},
		{
			name:     "ALTER TABLE MODIFY SERIAL",
			setup:    []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT)"},
			stmt:     "ALTER TABLE a MODIFY c SERIAL",
			code:     0,
			contains: []string{"`c` bigint(20) unsigned NOT NULL AUTO_INCREMENT", "UNIQUE KEY `c` (`c`)"},
		},
		{
			name:     "ALTER TABLE CHANGE SERIAL",
			setup:    []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT)"},
			stmt:     "ALTER TABLE a CHANGE c y SERIAL",
			code:     0,
			contains: []string{"`y` bigint(20) unsigned NOT NULL AUTO_INCREMENT", "UNIQUE KEY `y` (`y`)"},
		},
		{
			name:     "ALTER TABLE MIN_ROWS with a fraction",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a MIN_ROWS 1.5",
			code:     0,
			contains: []string{"MIN_ROWS=1"},
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
		// An expression that reads the column it defines is 4029, after the
		// table is looked up (1146) and after a column of the name exists (1060).
		{
			name:  "ADD a generated column reading itself",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a INT AS(a)",
			code:  4029,
		},
		{
			name:  "ADD a generated column reading itself as STORED",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a INT AS(a) STORED",
			code:  4029,
		},
		{
			name:  "ADD a column whose default reads it",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a INT DEFAULT a",
			code:  4029,
		},
		{
			name:  "ADD a column whose parenthesized default reads it",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a INT DEFAULT(a)",
			code:  4029,
		},
		{
			name:  "SET DEFAULT on a column that reads itself",
			setup: []string{intColumnTable},
			stmt:  "ALTER TABLE a ALTER a SET DEFAULT a",
			code:  4029,
		},
		{
			name:  "ADD a column reading a generated column defined before it",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT AS(b + 1) VIRTUAL)"},
			stmt:  "ALTER TABLE a ADD d INT AS(c)",
			code:  0,
		},
		{
			name:  "MODIFY a generated column to read itself",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT AS(b + 1) VIRTUAL)"},
			stmt:  "ALTER TABLE a MODIFY c INT AS(c)",
			code:  4029,
		},
		{
			name:  "CHANGE a generated column to read itself",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT AS(b + 1) VIRTUAL)"},
			stmt:  "ALTER TABLE a CHANGE c c INT AS(c)",
			code:  4029,
		},
		{
			name: "CREATE TABLE with a generated column reading itself",
			stmt: "CREATE TABLE x (a INT AS(a))",
			code: 4029,
		},
		{
			name: "CREATE TABLE with a default reading its column",
			stmt: "CREATE TABLE x (a INT DEFAULT a)",
			code: 4029,
		},
		{
			name: "CREATE TABLE with a generated column reading a later generated column",
			stmt: "CREATE TABLE x (a INT AS(b), b INT AS(a))",
			code: 4029,
		},
		{
			name: "CREATE TABLE with a generated column reading a later plain column",
			stmt: "CREATE TABLE x (a INT AS(b), b INT)",
			code: 0,
		},
		// A collation of another character set than the one the type settled is
		// 1253, answered while the statement is parsed: before the table is looked
		// up (1146), before a column of the name exists (1060), and before a
		// trailing token the parser would otherwise refuse.
		{
			name:  "ADD a column whose collation is of another character set",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin",
			code:  1253,
		},
		{
			name:  "ADD a column whose collation is of another character set, with a trailing token",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz",
			code:  1253,
		},
		{
			name:  "ADD a column of BINARY's collation to an ASCII column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY",
			code:  1253,
		},
		{
			name:  "ADD a BYTE column with a collation of another character set",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a CHAR BYTE COLLATE utf8mb4_bin",
			code:  1253,
		},
		{
			name:  "ADD a NCHAR column with a collation of another character set",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a NCHAR COLLATE utf8mb4_bin",
			code:  1253,
		},
		{
			name: "ADD a column over no table with a collation of another character set",
			stmt: "ALTER TABLE x ADD a CHAR ASCII COLLATE utf8mb4_bin",
			code: 1253,
		},
		{
			name:  "ADD a column whose collation belongs to its character set",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a CHAR CHARACTER SET ascii COLLATE ascii_general_ci",
			code:  0,
		},
		{
			name:  "CREATE TABLE over an existing name with a collation of another character set",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a (x CHAR ASCII COLLATE utf8mb4_bin)",
			code:  1253,
		},
		// A FULLTEXT key's parser is looked up while the statement is parsed: 1128
		// before the table is looked up (1146) and before a trailing token.
		{
			name:   "ADD FULLTEXT WITH PARSER over no table",
			stmt:   "ALTER TABLE x ADD FULLTEXT(a) WITH PARSER a",
			code:   1128,
			absent: "x",
		},
		{
			name:  "ADD FULLTEXT WITH PARSER with a trailing token",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz",
			code:  1128,
		},
		// A ROW_FORMAT FIXED or PAGE the statement writes on an engine that does
		// not take it is 1478; MyISAM and Aria take both.
		{
			name:  "ROW_FORMAT FIXED on an InnoDB table",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ROW_FORMAT FIXED",
			code:  1478,
		},
		{
			name:  "ROW_FORMAT PAGE on an InnoDB table",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ROW_FORMAT PAGE",
			code:  1478,
		},
		{
			name:  "ROW_FORMAT FIXED on a table the statement makes MyISAM",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ROW_FORMAT FIXED ENGINE=MyISAM",
			code:  0,
		},
		{
			name:  "ROW_FORMAT PAGE on a MyISAM table",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d)) ENGINE=MyISAM"},
			stmt:  "ALTER TABLE a ROW_FORMAT PAGE",
			code:  0,
		},
		{
			name:  "ROW_FORMAT FIXED on an Aria table",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d)) ENGINE=Aria"},
			stmt:  "ALTER TABLE a ROW_FORMAT FIXED",
			code:  0,
		},
		// A partition maintenance command takes a table that is not there
		// silently, where REBUILD PARTITION and REMOVE PARTITIONING are 1146.
		{
			name:   "OPTIMIZE PARTITION of a missing table",
			stmt:   "ALTER TABLE x OPTIMIZE PARTITION ALL",
			code:   0,
			absent: "x",
		},
		{
			name:   "ANALYZE PARTITION of a missing table",
			stmt:   "ALTER TABLE x ANALYZE PARTITION p1",
			code:   0,
			absent: "x",
		},
		{
			name:   "CHECK PARTITION of a missing table",
			stmt:   "ALTER TABLE x CHECK PARTITION ALL",
			code:   0,
			absent: "x",
		},
		{
			name:   "REPAIR PARTITION of a missing table",
			stmt:   "ALTER TABLE x REPAIR PARTITION ALL",
			code:   0,
			absent: "x",
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
			if tc.absent != "" {
				if got := catalogTableState(c, "d", tc.absent); got != "absent" {
					t.Fatalf("catalog keeps %s after %s:\n%s", tc.absent, tc.stmt, got)
				}
				db := simulationDatabase
				_ = db
				var count int
				if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tc.absent).Scan(&count); err != nil {
					t.Fatalf("count %s: %v", tc.absent, err)
				}
				if count != 0 {
					t.Fatalf("mariadb keeps %s after %s", tc.absent, tc.stmt)
				}
			}
			if len(tc.contains) > 0 || len(tc.omits) > 0 {
				var name, engineCreate string
				if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&name, &engineCreate); err != nil {
					t.Fatalf("SHOW CREATE TABLE a: %v", err)
				}
				catalogCreate := c.ShowCreateTable("d", "a")
				for _, want := range tc.contains {
					if !strings.Contains(engineCreate, want) {
						t.Fatalf("mariadb stored\n%s\nwithout %s", engineCreate, want)
					}
					if !strings.Contains(catalogCreate, want) {
						t.Fatalf("catalog stored\n%s\nwithout %s", catalogCreate, want)
					}
				}
				for _, no := range tc.omits {
					if strings.Contains(engineCreate, no) {
						t.Fatalf("mariadb stored\n%s\nwith %s", engineCreate, no)
					}
					if strings.Contains(catalogCreate, no) {
						t.Fatalf("catalog stored\n%s\nwith %s", catalogCreate, no)
					}
				}
			}
		})
	}
}
