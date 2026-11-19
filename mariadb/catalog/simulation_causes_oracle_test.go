package catalog

import (
	"context"
	"database/sql"
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
// mariadbShowCreate reads the server's SHOW CREATE TABLE of a table in the test's
// database.
func mariadbShowCreate(t *testing.T, conn *sql.DB, name string) string {
	t.Helper()
	var table, create string
	if err := conn.QueryRow("SHOW CREATE TABLE "+name).Scan(&table, &create); err != nil {
		t.Fatalf("SHOW CREATE TABLE %s: %v", name, err)
	}
	return normalizeShowCreate(create)
}

func TestSimulationCausesOracle(t *testing.T) {
	conn := mariadbDatabase(t)
	ctx := context.Background()
	reset := func() {
		for _, stmt := range []string{
			"DROP TEMPORARY TABLE IF EXISTS a",
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
		// keepsSchema holds the case whose statement fails with the table as it
		// was, which both sides must then still show: table a, or keptTable.
		keepsSchema bool
		keptTable   string
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
		// The columns a statement adds are read before any of them is built, so a
		// generated one may not read a generated column added after it; a plain
		// column added later is fine. The statement leaves the table as it was.
		{
			name:        "ADD two generated columns that read each other",
			setup:       []string{"CREATE TABLE x (k INT PRIMARY KEY)"},
			stmt:        "ALTER TABLE x ADD q INT AS(r), ADD r INT AS(q)",
			code:        4029,
			keepsSchema: true,
			keptTable:   "x",
		},
		{
			name:  "ADD two generated columns that read each other in one list",
			setup: []string{"CREATE TABLE x (k INT PRIMARY KEY)"},
			stmt:  "ALTER TABLE x ADD (q INT AS(r), r INT AS(q))",
			code:  4029,
		},
		{
			name:  "ADD a generated column reading a plain column added after it",
			setup: []string{"CREATE TABLE x (k INT PRIMARY KEY)"},
			stmt:  "ALTER TABLE x ADD q INT AS(r), ADD r INT",
			code:  0,
		},
		{
			name:        "ROW_FORMAT FIXED with the engine the table has",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ENGINE=InnoDB ROW_FORMAT FIXED",
			code:        1478,
			keepsSchema: true,
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
		// SYSTEM_TIME partitioning is a versioned table's (4124), and a
		// PERIOD FOR SYSTEM_TIME a CREATE TABLE declares needs WITH SYSTEM
		// VERSIONING (4125), answered before the table's name is looked up. An
		// ALTER that adds a period over ordinary columns runs and stores nothing.
		{
			name:        "PARTITION BY SYSTEM_TIME of a table that is not versioned",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a PARTITION BY SYSTEM_TIME",
			code:        4124,
			keepsSchema: true,
		},
		{
			name:        "CREATE TABLE with a period and no versioning",
			setup:       []string{oneTable},
			stmt:        "CREATE TABLE a(PERIOD FOR SYSTEM_TIME(a, a))",
			code:        4125,
			keepsSchema: true,
		},
		{
			name:        "ADD PERIOD FOR SYSTEM_TIME over ordinary columns",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a)",
			code:        0,
			keepsSchema: true,
		},
		// A CHECK or DEFAULT expression names columns of the table (1054).
		{
			name:  "ADD CHECK over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD CHECK(a)",
			code:  1054,
		},
		{
			name:  "ADD CONSTRAINT CHECK over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK(a)",
			code:  1054,
		},
		{
			name:  "ADD a column whose DEFAULT names no column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a INT DEFAULT VALUE",
			code:  1054,
		},
		// A key names a column once (1060), takes a length over a string column
		// only (1089), and takes no engine-defined option (1911); USING RTREE is
		// an R-tree only through SPATIAL (1210).
		{
			name:  "ADD KEY naming a column twice",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c, c)",
			code:  1060,
		},
		{
			name:  "ADD KEY with a length over an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c(1))",
			code:  1089,
		},
		{
			name:  "ADD KEY USING RTREE over an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY USING RTREE(c)",
			code:  1210,
		},
		{
			// The key's own type is judged before its name: the table holds an
			// index called a, and the column is no geometry.
			name:  "CREATE SPATIAL INDEX of a name the table has over an INT column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY a (c))"},
			stmt:  "CREATE SPATIAL INDEX a ON a(b)",
			code:  1210,
		},
		{
			name:  "CREATE FULLTEXT INDEX of a name the table has over an INT column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY a (c))"},
			stmt:  "CREATE FULLTEXT INDEX a ON a(b)",
			code:  1283,
		},
		{
			name:  "ADD KEY with an engine-defined option",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c) x = 1",
			code:  1911,
		},
		// AUTO_INCREMENT takes the first part of a key (1075), an ENUM lists a
		// value once (1291), a WITHOUT OVERLAPS part names a period (4156), and
		// IMPORT TABLESPACE needs a DISCARDed tablespace (1813).
		{
			name:        "ADD an AUTO_INCREMENT column no key leads with",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ADD a INT AUTO_INCREMENT",
			code:        1075,
			keepsSchema: true,
		},
		{
			name:  "ADD an ENUM column listing a value twice",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a ENUM('s', 's')",
			code:  1291,
		},
		{
			name:  "ADD UNIQUE over a period the table does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD UNIQUE(c, c WITHOUT OVERLAPS)",
			code:  4156,
		},
		{
			name:        "IMPORT TABLESPACE over a table whose tablespace was not DISCARDed",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a IMPORT TABLESPACE",
			code:        1813,
			keepsSchema: true,
		},
		{
			// A DISCARDed tablespace is gone, so there is no file to import: the
			// import asks for the secondary indexes to go first, and reports the
			// tablespace missing once they have.
			name:        "IMPORT TABLESPACE after DISCARD TABLESPACE with a secondary key",
			setup:       []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:        "ALTER TABLE a IMPORT TABLESPACE",
			code:        1815,
			keepsSchema: true,
		},
		{
			name:        "IMPORT TABLESPACE after DISCARD TABLESPACE without one",
			setup:       []string{"CREATE TABLE a (b INT PRIMARY KEY)", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:        "ALTER TABLE a IMPORT TABLESPACE",
			code:        1030,
			keepsSchema: true,
		},
		{
			name:        "IMPORT TABLESPACE over a MyISAM table",
			setup:       []string{"CREATE TABLE a (b INT) ENGINE=MyISAM"},
			stmt:        "ALTER TABLE a IMPORT TABLESPACE",
			code:        1031,
			keepsSchema: true,
		},
		{
			// An AUTO_INCREMENT column is the first part of a key on InnoDB, and
			// anywhere in one on MyISAM.
			name:        "ADD a column to a MyISAM table whose auto column is not first in its key",
			setup:       []string{"CREATE TABLE a (g INT, id INT AUTO_INCREMENT, KEY k (g, id)) ENGINE=MyISAM"},
			stmt:        "ALTER TABLE a ADD c INT",
			code:        0,
			keepsSchema: false,
		},
		// ALTER ... RENAME KEY of a key the table does not have is 1176, which the
		// new name does not change.
		{
			name:        "RENAME KEY of a key the table does not have",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a RENAME KEY nosuch TO idx_d",
			code:        1176,
			keepsSchema: true,
		},
		// The table the statement leaves decides SYSTEM_TIME partitioning, so a
		// statement that also adds the versioning runs.
		{
			name:  "ADD SYSTEM VERSIONING with PARTITION BY SYSTEM_TIME",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD SYSTEM VERSIONING PARTITION BY SYSTEM_TIME",
			code:  0,
		},
		// A key part naming no column is 1072 before the repeated name is 1060,
		// on a CREATE TABLE as much as on an ALTER.
		{
			name:   "CREATE TABLE with a key naming a missing column twice",
			stmt:   "CREATE TABLE a (b INT, KEY k (missing, missing))",
			code:   1072,
			absent: "a",
		},
		{
			name:   "CREATE TABLE with a primary key naming a missing column twice",
			stmt:   "CREATE TABLE a (b INT, PRIMARY KEY (missing, missing))",
			code:   1072,
			absent: "a",
		},
		// ENUM compares its values with the column's collation, so a
		// case-sensitive one lists 's' and 'S' as two.
		{
			name:     "ENUM of two spellings under a case-sensitive collation",
			stmt:     "CREATE TABLE a (e ENUM('s', 'S') COLLATE utf8mb4_bin)",
			code:     0,
			contains: []string{"enum('s','S')"},
		},
		// The auto column's key is the engine's: MEMORY wants it first, as
		// InnoDB does, where MyISAM and Aria take it anywhere.
		{
			name: "CREATE TABLE on MEMORY with the auto column second in its key",
			stmt: "CREATE TABLE a (g INT, id INT AUTO_INCREMENT, KEY k (g, id)) ENGINE=MEMORY",
			code: 1075,
		},
		{
			name: "CREATE TABLE on MyISAM with the auto column second in its key",
			stmt: "CREATE TABLE a (g INT, id INT AUTO_INCREMENT, KEY k (g, id)) ENGINE=MyISAM",
			code: 0,
		},
		// A key's columns are looked up before its type, on CREATE INDEX too.
		{
			name:   "CREATE INDEX USING RTREE over a missing column",
			setup:  []string{"CREATE TABLE a (b INT, KEY k (b))"},
			stmt:   "CREATE INDEX i USING RTREE ON a (missing)",
			code:   1072,
			absent: "i",
		},
		{
			name:  "CREATE INDEX USING RTREE over an INT column",
			setup: []string{"CREATE TABLE a (b INT, KEY k (b))"},
			stmt:  "CREATE INDEX i USING RTREE ON a (b)",
			code:  1210,
		},
		// A quoted name is a column name all the same, and one the table lacks is
		// 1072.
		{
			name:   "CREATE TABLE with a quoted key column the table lacks",
			stmt:   "CREATE TABLE a (b INT, KEY k (`missing name`))",
			code:   1072,
			absent: "a",
		},
		// ENGINE='DEFAULT' is the server's own engine for a tablespace.
		{
			name:        "IMPORT TABLESPACE over an ENGINE='DEFAULT' table",
			setup:       []string{"CREATE TABLE a (b INT) ENGINE='DEFAULT'"},
			stmt:        "ALTER TABLE a IMPORT TABLESPACE",
			code:        1813,
			keepsSchema: true,
		},
		// CREATE OR REPLACE reads its own collation, not the database's.
		{
			name:     "CREATE OR REPLACE TABLE with a case-sensitive table collation",
			setup:    []string{"CREATE TABLE a (b INT)"},
			stmt:     "CREATE OR REPLACE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin",
			code:     0,
			contains: []string{"enum('s','S')"},
		},
		// An ENUM's trailing spaces are its collation's: `('a','a ')` lists one
		// value twice, and `('a ','b ')` is stored as enum('a','b').
		{
			name: "ENUM listing a value and its trailing-space form",
			stmt: "CREATE TABLE a (e ENUM('a', 'a '))",
			code: 1291,
		},
		{
			name:     "ENUM of two values with trailing spaces",
			stmt:     "CREATE TABLE a (e ENUM('a ', 'b '))",
			code:     0,
			contains: []string{"enum('a','b')"},
		},
		// A discarded tablespace leaves nothing to copy into: the ALTERs that
		// need a copy are 1814, the in-place ones run.
		{
			name:        "ROW_FORMAT of a table whose tablespace was DISCARDed",
			setup:       []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:        "ALTER TABLE a ROW_FORMAT=COMPACT",
			code:        1814,
			keepsSchema: true,
		},
		{
			name:  "ADD COLUMN to a table whose tablespace was DISCARDed",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD e INT",
			code:  0,
		},
		// A key IF NOT EXISTS finds is a note, and the definition beside it is
		// not read.
		{
			name:  "ADD KEY IF NOT EXISTS of a name the table holds, with an option",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY IF NOT EXISTS idx_d (d) x = 1",
			code:  0,
		},
		// The table the statement leaves decides SYSTEM_TIME partitioning: one
		// that drops the versioning is 4124.
		{
			name:        "DROP SYSTEM VERSIONING with PARTITION BY SYSTEM_TIME",
			setup:       []string{"CREATE TABLE a (b INT) WITH SYSTEM VERSIONING"},
			stmt:        "ALTER TABLE a DROP SYSTEM VERSIONING PARTITION BY SYSTEM_TIME",
			code:        4124,
			keepsSchema: true,
		},
		// Turkish is its own case mapping: I and i are different letters, and
		// i and İ are the same one.
		{
			name: "ENUM of I and i under a Turkish collation",
			stmt: "CREATE TABLE a (e ENUM('I', 'i') COLLATE utf8mb4_turkish_ci)",
			code: 0,
		},
		{
			name: "ENUM of i and İ under a Turkish collation",
			stmt: "CREATE TABLE a (e ENUM('i', 'İ') COLLATE utf8mb4_turkish_ci)",
			code: 1291,
		},
		// A FULLTEXT key's parser is looked up before its columns.
		{
			name:   "FULLTEXT key over a missing column with an unknown parser",
			stmt:   "CREATE TABLE a (b TEXT, FULLTEXT KEY f (missing) WITH PARSER unavailable)",
			code:   1128,
			absent: "a",
		},
		{
			name:   "FULLTEXT key over a missing column",
			stmt:   "CREATE TABLE a (b TEXT, FULLTEXT KEY f (missing))",
			code:   1072,
			absent: "a",
		},
		// ENGINE='DEFAULT' is InnoDB's list, and a key IF NOT EXISTS adds is a
		// new key whose definition is read.
		{
			name:   "key option on an ENGINE='DEFAULT' table",
			stmt:   "CREATE TABLE a (b INT, KEY k (b) x = 1) ENGINE='DEFAULT'",
			code:   1911,
			absent: "a",
		},
		{
			name:  "ADD KEY IF NOT EXISTS of a name the table does not hold, with an option",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY IF NOT EXISTS j (d) x = 1",
			code:  1911,
		},
		// Every copy ALTER is 1814 on a discarded tablespace: FORCE, ORDER BY, and
		// the keys that rebuild the table.
		{
			name:        "FORCE of a table whose tablespace was DISCARDed",
			setup:       []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:        "ALTER TABLE a FORCE",
			code:        1814,
			keepsSchema: true,
		},
		{
			name:  "ADD INDEX to a table whose tablespace was DISCARDed",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD INDEX i (d)",
			code:  0,
		},
		// TRUNCATE gives the tablespace back, so the import that follows reports
		// it exists again.
		{
			name:  "IMPORT TABLESPACE after TRUNCATE of a discarded table",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE", "TRUNCATE TABLE a"},
			stmt:  "ALTER TABLE a IMPORT TABLESPACE",
			code:  1813,
		},
		// A key IF NOT EXISTS finds is a note, and its WITHOUT OVERLAPS part is
		// not read.
		{
			name:  "ADD UNIQUE KEY IF NOT EXISTS of a name the table holds, over a missing period",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD UNIQUE KEY IF NOT EXISTS idx_d (d, missing_period WITHOUT OVERLAPS)",
			code:  0,
		},
		// ALGORITHM=COPY asks for the copy by name.
		{
			name:        "ALGORITHM=COPY on a table whose tablespace was DISCARDed",
			setup:       []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:        "ALTER TABLE a ALGORITHM=COPY, ADD e INT",
			code:        1814,
			keepsSchema: true,
		},
		{
			name:  "ALGORITHM=INPLACE on a table whose tablespace was DISCARDed",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ALGORITHM=INPLACE, ADD e INT",
			code:  0,
		},
		// A column's own CHARACTER SET gives it that set's default collation, not
		// the table's.
		{
			name: "ENUM under a CHARACTER SET whose default collation folds case",
			stmt: "CREATE TABLE a (e ENUM('s', 'S') CHARACTER SET latin1) COLLATE utf8mb4_bin",
			code: 1291,
		},
		// A primary key's parts are held to their types before a repeated name.
		{
			name:   "primary key of a length over an INT column, repeated",
			stmt:   "CREATE TABLE a (b INT, PRIMARY KEY (b(1), b))",
			code:   1089,
			absent: "a",
		},
		{
			name:   "primary key over a TEXT column",
			stmt:   "CREATE TABLE a (b TEXT, PRIMARY KEY (b))",
			code:   1170,
			absent: "a",
		},
		// A CHECK may name a column the same statement adds after it.
		{
			name:  "two added columns whose CHECKs read each other",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT CHECK (e > f), ADD f INT CHECK (f < e)",
			code:  0,
		},
		// A column's own UNIQUE and REFERENCES rebuild a discarded table where its
		// CHECK and DEFAULT do not, and a UNIQUE or FULLTEXT index does where an
		// ordinary one does not.
		{
			name:  "ADD a column with a CHECK to a table whose tablespace was DISCARDed",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD e INT CHECK (e > 0)",
			code:  0,
		},
		{
			name:  "CREATE INDEX over a table whose tablespace was DISCARDed",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "CREATE INDEX i ON a (c)",
			code:  0,
		},
		{
			name:  "CREATE UNIQUE INDEX over a table whose tablespace was DISCARDed, with a primary key",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "CREATE UNIQUE INDEX u ON a (c)",
			code:  0,
		},
		// A column IF NOT EXISTS finds is a note whose definition is not read, its
		// ENUM values and its DEFAULT alike.
		{
			name:  "ADD COLUMN IF NOT EXISTS of a column the table holds, with duplicated ENUM values",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD COLUMN IF NOT EXISTS c ENUM('x', 'x'), ADD e INT AS (1)",
			code:  0,
		},
		{
			name:  "ADD COLUMN IF NOT EXISTS of a column the table holds, with a DEFAULT that names nobody",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD COLUMN IF NOT EXISTS c INT DEFAULT missing",
			code:  0,
		},
		// A column the table already holds is 1060 before the key it would be
		// is read, where the same statement over a free name is 1068.
		{
			name:  "ADD an existing column SERIAL KEY KEY",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, a INT, KEY idx_d (b))"},
			stmt:  "ALTER TABLE a ADD a SERIAL KEY KEY",
			code:  1060,
		},
		// RENAME KEY over a key the table does not hold is 1176, not the 1061
		// a duplicate name answers.
		{
			name:  "RENAME KEY over a key the table does not hold",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a RENAME KEY zzz TO yyy",
			code:  1176,
		},
		// A key that lists the same column twice is 1060.
		{
			name:  "ADD KEY over the same column twice",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT)"},
			stmt:  "ALTER TABLE a ADD KEY(a, a)",
			code:  1060,
		},
		// Parenthesized ADD (KEY ...) runs the same engine-option check as ADD KEY.
		{
			name:        "parenthesized ADD KEY with an engine-defined option",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ADD (KEY k (c) x = 1)",
			code:        1911,
			keepsSchema: true,
		},
		// ALTER COLUMN ... SET DEFAULT and GENERATED ALWAYS AS read their
		// columns against the table the statement leaves (1054).
		{
			name:        "ALTER COLUMN SET DEFAULT over a missing column",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ALTER COLUMN c SET DEFAULT missing",
			code:        1054,
			keepsSchema: true,
		},
		{
			name:     "ALTER COLUMN SET DEFAULT naming a column added after it",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ALTER COLUMN c SET DEFAULT e, ADD COLUMN e INT",
			code:     0,
			contains: []string{"DEFAULT `e`"},
		},
		{
			name:   "CREATE TABLE generated column over a missing column",
			stmt:   "CREATE TABLE a (b INT, g INT AS (missing))",
			code:   1054,
			absent: "a",
		},
		{
			name:        "ALTER TABLE ADD generated column over a missing column",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ADD g INT AS (missing)",
			code:        1054,
			keepsSchema: true,
		},
		// Skipped ADD COLUMN IF NOT EXISTS clauses neither read generated
		// columns added later nor count as generated columns themselves.
		{
			name:  "skipped ADD COLUMN IF NOT EXISTS generated forward ref",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD COLUMN IF NOT EXISTS c INT AS (e), ADD e INT AS (1)",
			code:  0,
		},
		{
			name:  "generated column reading an existing column beside a skipped generated IF NOT EXISTS",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT AS (c), ADD COLUMN IF NOT EXISTS c INT AS (1)",
			code:  0,
		},
		// ENUM collation in ALTER follows the table collation the statement
		// leaves, CONVERT TO re-checks existing ENUM values, MODIFY resolves the
		// column name before 1291, and CREATE OR REPLACE drops the old table
		// before 1291.
		{
			name:     "ADD ENUM followed by table COLLATE utf8mb4_bin",
			setup:    []string{"CREATE TABLE a (b INT)"},
			stmt:     "ALTER TABLE a ADD e ENUM('s', 'S'), COLLATE utf8mb4_bin",
			code:     0,
			contains: []string{"enum('s','S')"},
			omits:    []string{"CHARACTER SET utf8mb4 COLLATE"},
		},
		{
			name:        "CONVERT TO case-insensitive charset on a binary-collated ENUM",
			setup:       []string{"CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin"},
			stmt:        "ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4",
			code:        1291,
			keepsSchema: true,
		},
		{
			name:        "MODIFY a missing column to a duplicated ENUM",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a MODIFY missing ENUM('s', 's')",
			code:        1054,
			keepsSchema: true,
		},
		{
			name:   "CREATE OR REPLACE TABLE with a duplicated ENUM drops the old table",
			setup:  []string{"CREATE TABLE a (b INT)"},
			stmt:   "CREATE OR REPLACE TABLE a (e ENUM('s', 'S'))",
			code:   1291,
			absent: "a",
		},
		// WITHOUT OVERLAPS resolves its period before missing key columns on
		// CREATE INDEX and CREATE TABLE alike.
		{
			name:        "CREATE UNIQUE INDEX over a missing column and a missing period",
			setup:       []string{oneTable},
			stmt:        "CREATE UNIQUE INDEX u ON a (missing, missing WITHOUT OVERLAPS)",
			code:        4156,
			keepsSchema: true,
		},
		{
			name:   "CREATE TABLE UNIQUE over a missing column and a missing period",
			stmt:   "CREATE TABLE a (b INT, UNIQUE (missing, missing WITHOUT OVERLAPS))",
			code:   4156,
			absent: "a",
		},
		// Table-level CHECK constraints resolve their columns against the table
		// the atomic ALTER leaves, and DISCARD/IMPORT TABLESPACE refuse
		// temporary InnoDB tables (3007) and partitioned InnoDB tables (1031).
		{
			name:     "table CHECK added before the column it reads",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD CHECK (e > 0), ADD e INT",
			code:     0,
			contains: []string{"CHECK (`e` > 0)"},
		},
		{
			name:        "DISCARD TABLESPACE on an InnoDB temporary table",
			setup:       []string{"CREATE TEMPORARY TABLE a (b INT)"},
			stmt:        "ALTER TABLE a DISCARD TABLESPACE",
			code:        3007,
			keepsSchema: true,
		},
		{
			name:        "IMPORT TABLESPACE on an InnoDB temporary table",
			setup:       []string{"CREATE TEMPORARY TABLE a (b INT)"},
			stmt:        "ALTER TABLE a IMPORT TABLESPACE",
			code:        3007,
			keepsSchema: true,
		},
		{
			name:        "DISCARD TABLESPACE on a MyISAM temporary table",
			setup:       []string{"CREATE TEMPORARY TABLE a (b INT) ENGINE=MyISAM"},
			stmt:        "ALTER TABLE a DISCARD TABLESPACE",
			code:        1031,
			keepsSchema: true,
		},
		{
			name:        "DISCARD TABLESPACE on a partitioned InnoDB table",
			setup:       []string{"CREATE TABLE a (b INT) PARTITION BY HASH (b) PARTITIONS 2"},
			stmt:        "ALTER TABLE a DISCARD TABLESPACE",
			code:        1031,
			keepsSchema: true,
		},
		// Expressions written in an ALTER bind against the final table's column
		// names rather than pre-rename names: RENAME COLUMN rewrites only stored
		// expressions from earlier statements.
		{
			name:        "added CHECK naming a column renamed later in the same ALTER",
			setup:       []string{oneTable},
			stmt:        "ALTER TABLE a ADD e INT CHECK (c > 0), RENAME COLUMN c TO z",
			code:        1054,
			keepsSchema: true,
		},
		{
			name:     "added CHECK naming the post-rename column of the same ALTER",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD e INT CHECK (z > 0), RENAME COLUMN c TO z",
			code:     0,
			contains: []string{"CHECK (`z` > 0)"},
		},
		// A charset or collation written as DEFAULT is resolved while it is read:
		// `CHARSET DEFAULT` is utf8mb4 with utf8mb4_general_ci, and `COLLATE
		// DEFAULT` the table's charset with its 11.8 default (a case-insensitive
		// one for utf8mb4, so an existing ENUM('s','S') is 1291 on CONVERT).
		{
			name:     "CONVERT TO CHARACTER SET DEFAULT",
			setup:    []string{"CREATE TABLE a (b VARCHAR(3)) COLLATE utf8mb4_bin"},
			stmt:     "ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT",
			code:     0,
			contains: []string{"COLLATE=utf8mb4_general_ci"},
		},
		{
			name:  "CONVERT TO CHARACTER SET DEFAULT over a binary-collated ENUM",
			setup: []string{"CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin"},
			stmt:  "ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT",
			code:  1291,
		},
		// The collation a CONVERT TO CHARACTER SET DEFAULT leaves is the one an
		// ENUM the same statement adds is read under, wherever in it they stand.
		{
			name:        "CONVERT TO CHARACTER SET DEFAULT with an ENUM added after it",
			setup:       []string{"CREATE TABLE a (b INT) COLLATE utf8mb4_bin"},
			stmt:        "ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT, ADD e ENUM('s', 'S')",
			code:        1291,
			keepsSchema: true,
		},
		{
			name:     "CONVERT TO CHARACTER SET DEFAULT with an ENUM added after it, listing two values",
			setup:    []string{"CREATE TABLE a (b INT) COLLATE utf8mb4_bin"},
			stmt:     "ALTER TABLE a CONVERT TO CHARACTER SET DEFAULT, ADD e ENUM('s', 't')",
			code:     0,
			contains: []string{"COLLATE=utf8mb4_general_ci"},
		},
		{
			name:     "COLLATE DEFAULT on a binary-collated table",
			setup:    []string{"CREATE TABLE a (e ENUM('s', 'S')) COLLATE utf8mb4_bin"},
			stmt:     "ALTER TABLE a COLLATE DEFAULT",
			code:     0,
			contains: []string{"COLLATE=utf8mb4_uca1400_ai_ci"},
		},
		// An engine-defined key option IGNORE_BAD_TABLE_OPTIONS tolerated is
		// commented out once the mode is gone, a VECTOR key's own M and DISTANCE
		// excepted.
		{
			name:     "VECTOR key option the engine defines beside one it does not",
			setup:    []string{"SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'", "CREATE TABLE a (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=cosine x=1)", "SET sql_mode = DEFAULT"},
			stmt:     "ALTER TABLE a ADD c INT",
			code:     0,
			contains: []string{"VECTOR KEY `i` (`v`) `M`=6 `DISTANCE`=cosine /* `x`=1 */"},
		},
		// A CREATE INDEX name the table holds is 1061 before the unknown engine
		// option is 1911, where a free name answers 1911, and a missing key column
		// answers 1072 before either.
		{
			name:        "CREATE INDEX of a name the table holds, with an engine-defined option",
			setup:       []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"},
			stmt:        "CREATE INDEX idx_d ON a (c) x = 1",
			code:        1061,
			keepsSchema: true,
		},
		{
			name:        "CREATE INDEX of a free name, with an engine-defined option",
			setup:       []string{oneTable},
			stmt:        "CREATE INDEX i ON a (c) x = 1",
			code:        1911,
			keepsSchema: true,
		},
		{
			name:        "CREATE INDEX of a name the table holds, over a missing column, with an engine-defined option",
			setup:       []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"},
			stmt:        "CREATE INDEX idx_d ON a (missing) x = 1",
			code:        1072,
			keepsSchema: true,
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
			keptTable := tc.keptTable
			if keptTable == "" {
				keptTable = "a"
			}
			beforeEngine := ""
			if tc.keepsSchema {
				beforeEngine = mariadbShowCreate(t, conn, keptTable)
			}

			c := New()
			execAlterClauses(t, c, "CREATE DATABASE d; USE d;")
			for _, stmt := range tc.setup {
				execAlterClauses(t, c, stmt)
			}
			beforeCatalog := c.ShowCreateTable("d", keptTable)
			catalogCode, _ := catalogErr(t, c, tc.stmt)
			if catalogCode != engineCode {
				t.Fatalf("catalog: %s answered %d, mariadb answered %d", tc.stmt, catalogCode, engineCode)
			}
			if tc.keepsSchema {
				if got := mariadbShowCreate(t, conn, keptTable); got != beforeEngine {
					t.Fatalf("mariadb kept\n%s\nafter %s, where it had\n%s", got, tc.stmt, beforeEngine)
				}
				if got := c.ShowCreateTable("d", keptTable); got != beforeCatalog {
					t.Fatalf("catalog kept\n%s\nafter %s, where it had\n%s", got, tc.stmt, beforeCatalog)
				}
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
