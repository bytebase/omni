package catalog

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// TestSimulationCausesOracle holds the refusals testdata/simulate.tsv records to
// mysql:8.0.46 cause by cause: the object an ALTER or DROP VIEW names, a key
// built on a column the table does not have, a FULLTEXT key's parser, a
// partition action on a table that has no partitions, a type's display width, a
// view body naming a column no relation has, and the columns MySQL gives
// another type from the text they were: SERIAL, CHAR BYTE and CHAR COLLATE
// BINARY. Each case is answered by the pinned server and by the catalog, and
// both must agree with the code the case states.
func TestSimulationCausesOracle(t *testing.T) {
	conn := testcontainer.OpenEmpty(t)
	ctx := context.Background()
	reset := func() {
		for _, stmt := range []string{"SET SESSION sql_mode = DEFAULT", "DROP VIEW IF EXISTS a", "DROP TABLE IF EXISTS a", "DROP TABLE IF EXISTS p", "DROP TABLE IF EXISTS s"} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
	}
	t.Cleanup(reset)

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	// A table with no primary key, for the one a statement adds beside a column.
	noPrimaryKey := "CREATE TABLE a (c INT, d INT, KEY idx_d (d))"
	cases := []struct {
		name  string
		setup []string
		// session runs on the connection before the statement, for the cases
		// whose answer depends on the session's sql_mode.
		session []string
		stmt    string
		code    int
		// contains and omits are text both SHOW CREATE TABLE texts must hold
		// and must not hold.
		contains []string
		omits    []string
	}{
		{
			name:  "ALTER VIEW over a table",
			setup: []string{oneTable},
			stmt:  "ALTER VIEW a AS SELECT b FROM a",
			code:  1347,
		},
		{
			name: "ALTER VIEW over no object",
			stmt: "ALTER VIEW a AS SELECT b FROM a",
			code: 1146,
		},
		{
			name:  "ALTER VIEW over a view",
			setup: []string{"CREATE VIEW a AS SELECT 1 AS b"},
			stmt:  "ALTER VIEW a AS SELECT 1 AS b",
			code:  0,
		},
		{
			name:  "DROP VIEW over a table",
			setup: []string{oneTable},
			stmt:  "DROP VIEW a",
			code:  1347,
		},
		{
			name: "DROP VIEW over no object",
			stmt: "DROP VIEW a",
			code: 1051,
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
			// A FULLTEXT key takes the character types only.
			name:  "ADD FULLTEXT over a column of another type",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(b)",
			code:  1283,
		},
		{
			name:  "ADD FULLTEXT of a name the table has over a column of another type",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT idx_d(c)",
			code:  1061,
		},
		{
			name:  "ADD FULLTEXT over a later-added column of another type",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(x), ADD x INT",
			code:  1283,
		},
		{
			name:  "ADD FULLTEXT over a later-added TEXT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(x), ADD x TEXT",
			code:  0,
		},
		{
			name:  "CREATE FULLTEXT INDEX over a column of another type",
			setup: []string{oneTable},
			stmt:  "CREATE FULLTEXT INDEX ft ON a(b)",
			code:  1283,
		},
		{
			name:  "CREATE TABLE with a FULLTEXT key over a column of another type",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT, FULLTEXT KEY ft (b))",
			code:  1283,
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
			code:  1061,
		},
		{
			name:  "ADD PRIMARY KEY over a missing column on a table with one",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PRIMARY KEY(missing)",
			code:  1068,
		},
		{
			name:  "ADD FULLTEXT WITH PARSER over a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a",
			code:  1128,
		},
		{
			name:  "ADD FULLTEXT WITH PARSER ngram",
			setup: []string{"CREATE TABLE a (b TEXT)"},
			stmt:  "ALTER TABLE a ADD FULLTEXT(b) WITH PARSER ngram",
			code:  0,
		},
		{
			name:  "CREATE INDEX over a missing column",
			setup: []string{oneTable},
			stmt:  "CREATE INDEX i ON a(missing)",
			code:  1072,
		},
		{
			name:  "CREATE INDEX of a name the table has over a missing column",
			setup: []string{oneTable},
			stmt:  "CREATE INDEX idx_d ON a(missing)",
			code:  1061,
		},
		{
			name:  "CREATE FULLTEXT INDEX WITH PARSER over a missing column",
			setup: []string{oneTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a(missing) WITH PARSER a",
			code:  1128,
		},
		{
			name:  "CREATE TABLE with a key over a missing column",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT, KEY(missing))",
			code:  1072,
		},
		{
			// MySQL judges a key against the columns the statement leaves, so a
			// key may name a column a later command adds.
			name:  "ADD KEY over a column a later command adds",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c2), ADD c2 INT",
			code:  0,
		},
		{
			name:  "ADD KEY over a column the statement drops",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c), DROP COLUMN c",
			code:  1072,
		},
		{
			name:  "ADD KEY over a column the statement drops and adds again",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c), DROP COLUMN c, ADD COLUMN c INT",
			code:  0,
		},
		{
			name:  "ADD KEY over a column a later command renames",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(y), CHANGE c y INT",
			code:  0,
		},
		{
			name:  "ADD KEY over a column an earlier command renames",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a RENAME COLUMN c TO y, ADD KEY(c)",
			code:  1072,
		},
		{
			name:  "ADD KEY over a column the statement adds beside a rename",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(x), ADD COLUMN x INT, RENAME COLUMN c TO y",
			code:  0,
		},
		{
			// A MODIFY, CHANGE or RENAME starts from a column of the table the
			// statement found, so one the statement adds is 1054.
			name:  "CHANGE of a column the statement adds",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD COLUMN x INT, CHANGE x y INT",
			code:  1054,
		},
		{
			// A key is judged against the column a later command adds, with the
			// definition that command writes: a prefix past its length is 1089.
			name:  "ADD KEY with a prefix past the length of a later-added column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(x(11)), ADD x VARCHAR(10)",
			code:  1089,
		},
		{
			// A TEXT part with no length is 1170.
			name:  "ADD KEY over a later-added TEXT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(x), ADD x TEXT",
			code:  1170,
		},
		{
			// A foreign key is judged on the columns the statement leaves: 3780
			// compares the local column a later command adds.
			name:  "ADD FOREIGN KEY whose local column a later command adds",
			setup: []string{oneTable, "CREATE TABLE s (id INT PRIMARY KEY)"},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id), ADD x2 VARCHAR(10)",
			code:  3780,
		},
		{
			name:  "ADD FOREIGN KEY SET NULL beside a later-added NOT NULL column",
			setup: []string{oneTable, "CREATE TABLE s (id INT PRIMARY KEY)"},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id) ON DELETE SET NULL, ADD x2 INT NOT NULL",
			code:  1830,
		},
		{
			name:  "ADD FOREIGN KEY whose local column a later command adds as its type",
			setup: []string{oneTable, "CREATE TABLE s (id INT PRIMARY KEY)"},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id), ADD x2 INT",
			code:  0,
		},
		{
			name:  "ADD SPATIAL KEY over a later-added nullable column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD SPATIAL KEY(x), ADD x GEOMETRY",
			code:  1252,
		},
		{
			name:  "ADD PRIMARY KEY over a later-added column",
			setup: []string{noPrimaryKey},
			stmt:  "ALTER TABLE a ADD PRIMARY KEY(x2), ADD x2 INT",
			code:  0,
		},
		{
			// The rename's source is resolved against the table the statement
			// found, so it answers before the key's columns are read.
			name:  "ADD KEY over a column the statement adds but a rename cannot take",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(x), ADD COLUMN x INT, RENAME COLUMN x TO y",
			code:  1054,
		},
		{
			name:  "RENAME COLUMN of no column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a RENAME COLUMN x TO y",
			code:  1054,
		},
		{
			// A DROP COLUMN is resolved against the table the statement found,
			// so a column it adds cannot be dropped by it.
			name:  "DROP COLUMN of a column the statement adds",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD COLUMN x INT, ADD INDEX idx_x (x), DROP COLUMN x",
			code:  1091,
		},
		{
			name:  "CREATE TABLE with a FULLTEXT key's unknown parser",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT, FULLTEXT(b) WITH PARSER a)",
			code:  1128,
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
			name: "CREATE VIEW over a missing column in a function",
			stmt: "CREATE VIEW a AS SELECT ABS(a)",
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
			name:  "CREATE VIEW over a table with a resolved body",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW a AS SELECT b FROM a",
			code:  1050,
		},
		{
			name:     "ALTER TABLE ADD SERIAL",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD a SERIAL",
			code:     0,
			contains: []string{"`a` bigint", "NOT NULL AUTO_INCREMENT", "UNIQUE KEY `a` (`a`)"},
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
			name:     "ALTER TABLE ADD SERIAL NULL",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD a SERIAL NULL",
			code:     0,
			contains: []string{"`a` bigint unsigned AUTO_INCREMENT", "UNIQUE KEY `a` (`a`)"},
			omits:    []string{"`a` bigint unsigned NOT NULL"},
		},
		{
			name:     "ALTER TABLE MODIFY SERIAL",
			setup:    []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT)"},
			stmt:     "ALTER TABLE a MODIFY c SERIAL",
			code:     0,
			contains: []string{"`c` bigint unsigned NOT NULL AUTO_INCREMENT", "UNIQUE KEY `c` (`c`)"},
		},
		{
			name:     "ALTER TABLE CHANGE SERIAL",
			setup:    []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT)"},
			stmt:     "ALTER TABLE a CHANGE c y SERIAL",
			code:     0,
			contains: []string{"`y` bigint unsigned NOT NULL AUTO_INCREMENT", "UNIQUE KEY `y` (`y`)"},
		},
		{
			name:  "REMOVE PARTITIONING of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  1505,
		},
		{
			name:  "ADD PARTITION to a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p0)",
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
			// The new definitions are judged before the table they would
			// replace partitions of, as ADD PARTITION's are.
			name:  "REORGANIZE PARTITION into a bound that reads a name as a column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (b))",
			code:  1054,
		},
		{
			name:  "REORGANIZE PARTITION into an engine the server does not ship",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 ENGINE nosuch)",
			code:  1286,
		},
		{
			name:  "REORGANIZE PARTITION into a constant bound of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10))",
			code:  1505,
		},
		{
			name:  "REORGANIZE PARTITION whose subpartition names an unknown engine",
			setup: []string{"CREATE TABLE a (b INT) PARTITION BY RANGE(b) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10))"},
			stmt:  "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10) ENGINE InnoDB (SUBPARTITION s0 ENGINE nosuch))",
			code:  1286,
		},
		{
			// A definition's own errors answer before the table is looked up.
			name: "ADD PARTITION over a name its bound reads as a column on no table",
			stmt: "ALTER TABLE a ADD PARTITION (PARTITION p0 VALUES LESS THAN (b))",
			code: 1054,
		},
		{
			name: "REORGANIZE PARTITION naming an unknown engine on no table",
			stmt: "ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 ENGINE nosuch)",
			code: 1286,
		},
		{
			name:  "REBUILD PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a REBUILD PARTITION ALL",
			code:  1505,
		},
		{
			name:  "DISCARD PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DISCARD PARTITION ALL TABLESPACE",
			code:  1505,
		},
		{
			name:  "IMPORT PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a IMPORT PARTITION ALL TABLESPACE",
			code:  1505,
		},
		{
			name:  "EXCHANGE PARTITION with another table on a table with none",
			setup: []string{oneTable, "CREATE TABLE s (b INT)"},
			stmt:  "ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s",
			code:  1505,
		},
		{
			name:  "EXCHANGE PARTITION with the table itself",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a",
			code:  1066,
		},
		{
			// MySQL's maintenance actions are taken on a table with no
			// partitions, where REBUILD PARTITION is refused (1505).
			name:  "ANALYZE PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ANALYZE PARTITION ALL",
			code:  0,
		},
		{
			name:  "CHECK PARTITION of a table with none",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CHECK PARTITION ALL",
			code:  0,
		},
		{
			name:  "REMOVE PARTITIONING of a partitioned table",
			setup: []string{"CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2"},
			stmt:  "ALTER TABLE a REMOVE PARTITIONING",
			code:  0,
		},
		{
			name:  "ADD BIT past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e BIT(65)",
			code:  1439,
		},
		{
			name:  "ADD INT past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT(256)",
			code:  1439,
		},
		{
			name:  "ADD INT past the width MySQL holds at all",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT(4294967296)",
			code:  1439,
		},
		{
			name:  "ADD CHAR past the width MySQL holds at all",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e CHAR(4294967296)",
			code:  1439,
		},
		{
			name:  "ADD DECIMAL past the width MySQL holds at all",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e DECIMAL(4294967296)",
			code:  1439,
		},
		{
			// The width is checked before the column's name.
			name:  "ADD a column the table has with a width past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD b INT(256)",
			code:  1439,
		},
		{
			// The width is checked before the table is looked up.
			name: "ADD a column with a width past its display width to no table",
			stmt: "ALTER TABLE a ADD e INT(256)",
			code: 1439,
		},
		{
			name:  "MODIFY a column past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a MODIFY b INT(256)",
			code:  1439,
		},
		{
			name: "CREATE TABLE with a column past its display width",
			stmt: "CREATE TABLE s (b INT(256))",
			code: 1439,
		},
		{
			// The bound is a constant expression, so MySQL reads a name in it
			// as a column of the partition function, even a column the table has.
			name:  "ADD PARTITION over a name its bound reads as a column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (b))",
			code:  1054,
		},
		{
			name:  "ADD PARTITION naming an engine the server does not ship",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p1 ENGINE nosuch)",
			code:  1286,
		},
		{
			name:  "ADD PARTITION whose bound column is judged before its engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (b) ENGINE nosuch)",
			code:  1054,
		},
		{
			name:  "ALTER TABLE naming an engine the server does not ship",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE = nosuch",
			code:  1286,
		},
		{
			name: "ALTER TABLE naming an engine the server does not ship on no table",
			stmt: "ALTER TABLE a ENGINE = nosuch",
			code: 1286,
		},
		{
			name:  "CREATE TABLE naming an engine the server does not ship",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) ENGINE = nosuch",
			code:  1286,
		},
		{
			// The checks run in the statement's order, both being made while it
			// is prepared.
			name:  "ADD past its display width before an engine the server does not ship",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT(256), ENGINE = nosuch",
			code:  1439,
		},
		{
			name:  "an engine the server does not ship before an ADD past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE = nosuch, ADD e INT(256)",
			code:  1286,
		},
		{
			name:  "CREATE TABLE with a column past its display width and an unknown engine",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT(256)) ENGINE = nosuch",
			code:  1439,
		},
		{
			// PERFORMANCE_SCHEMA holds the performance schema's own tables.
			name:  "CREATE TABLE with the PERFORMANCE_SCHEMA engine",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) ENGINE = PERFORMANCE_SCHEMA",
			code:  1683,
		},
		{
			name:  "ALTER TABLE with the PERFORMANCE_SCHEMA engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE = PERFORMANCE_SCHEMA",
			code:  1031,
		},
		{
			name: "ALTER TABLE with the PERFORMANCE_SCHEMA engine on no table",
			stmt: "ALTER TABLE a ENGINE = PERFORMANCE_SCHEMA",
			code: 1146,
		},
		{
			// Without NO_ENGINE_SUBSTITUTION, MySQL substitutes the default
			// engine and takes the statement with a warning.
			name:    "CREATE TABLE with an unknown engine where the session substitutes",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "CREATE TABLE s (b INT) ENGINE = nosuch",
			code:    0,
		},
		{
			name:    "ALTER TABLE with an unknown engine where the session substitutes",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "ALTER TABLE a ENGINE = nosuch",
			code:    0,
		},
		{
			// TRADITIONAL expands to a list that carries NO_ENGINE_SUBSTITUTION;
			// ANSI's does not.
			name:    "CREATE TABLE with an unknown engine under TRADITIONAL",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'TRADITIONAL'"},
			stmt:    "CREATE TABLE s (b INT) ENGINE = nosuch",
			code:    1286,
		},
		{
			name:    "CREATE TABLE with an unknown engine under ANSI",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'ANSI'"},
			stmt:    "CREATE TABLE s (b INT) ENGINE = nosuch",
			code:    0,
		},
		{
			name:    "ALTER TABLE with an unknown engine where the session refuses it",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'NO_ENGINE_SUBSTITUTION'"},
			stmt:    "ALTER TABLE a ENGINE = nosuch",
			code:    1286,
		},
		{
			// MERGE is MRG_MYISAM's other name, as HEAP is MEMORY's.
			name:  "CREATE TABLE with the MERGE engine",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) ENGINE = MERGE",
			code:  0,
		},
		{
			name:  "ALTER TABLE naming the MERGE engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE = MERGE",
			code:  0,
		},
		{
			name:  "CREATE TABLE with the HEAP engine",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) ENGINE = HEAP",
			code:  0,
		},
		{
			name:  "CREATE TABLE whose partition bound reads a name as a column",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) PARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (b))",
			code:  1054,
		},
		{
			name:  "CREATE TABLE whose subpartition names an unknown engine",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) PARTITION BY RANGE(b) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) ENGINE InnoDB (SUBPARTITION s0 ENGINE nosuch))",
			code:  1286,
		},
		{
			name:  "CREATE TABLE with the DEFAULT engine",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT) ENGINE = 'DEFAULT'",
			code:  0,
		},
		{
			name:  "ALTER TABLE with the DEFAULT engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE = 'DEFAULT'",
			code:  0,
		},
		{
			name:  "ADD FLOAT with a scale past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e FLOAT(256,2)",
			code:  1439,
		},
		{
			name:  "ADD DOUBLE with a scale past its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e DOUBLE(256,2)",
			code:  1439,
		},
		{
			name:  "ADD FLOAT with a scale at its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e FLOAT(255,2)",
			code:  0,
		},
		{
			name:  "ADD INT at its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT(255)",
			code:  0,
		},
		{
			name:  "ADD BIT at its display width",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e BIT(64)",
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
			for _, stmt := range tc.session {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			_, err := conn.ExecContext(ctx, tc.stmt)
			engineCode, engineMessage := simulationEngineRefusal(t, tc.stmt, err)
			if engineCode != tc.code {
				t.Fatalf("mysql: %s answered %d (%s), want %d", tc.stmt, engineCode, engineMessage, tc.code)
			}

			c := New()
			execAlterClauses(t, c, "CREATE DATABASE d; USE d;")
			for _, stmt := range tc.setup {
				execAlterClauses(t, c, stmt)
			}
			for _, stmt := range tc.session {
				execAlterClauses(t, c, stmt)
			}
			catalogCode, catalogMessage := simulationCatalogRefusal(t, c, tc.stmt)
			if catalogCode != engineCode {
				t.Fatalf("catalog: %s answered %d (%s), mysql answered %d (%s)",
					tc.stmt, catalogCode, catalogMessage, engineCode, engineMessage)
			}
			if tc.code == 0 && tc.stmt == "ALTER TABLE a REMOVE PARTITIONING" {
				db := c.GetDatabase("d")
				if db == nil || db.GetTable("a") == nil || db.GetTable("a").Partitioning != nil {
					t.Fatalf("catalog kept the table partitioned after %s", tc.stmt)
				}
			}
			// MySQL stores an engine's alias under its canonical name.
			if tc.code == 0 {
				tbl := c.GetDatabase("d").GetTable("s")
				switch tc.stmt {
				case "CREATE TABLE s (b INT) ENGINE = MERGE":
					if tbl == nil || tbl.Engine != "MRG_MyISAM" {
						t.Fatalf("catalog stored engine %q for MERGE, want MRG_MyISAM", tbl.Engine)
					}
				case "CREATE TABLE s (b INT) ENGINE = HEAP":
					if tbl == nil || tbl.Engine != "MEMORY" {
						t.Fatalf("catalog stored engine %q for HEAP, want MEMORY", tbl.Engine)
					}
				}
			}
			if tc.code == 0 && tc.stmt == "ALTER TABLE a ADD PRIMARY KEY(x2), ADD x2 INT" {
				tbl := c.GetDatabase("d").GetTable("a")
				if col := tbl.GetColumn("x2"); col == nil || col.Nullable {
					t.Fatalf("catalog left x2 nullable after %s", tc.stmt)
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
						t.Fatalf("mysql stored\n%s\nwithout %s", engineCreate, want)
					}
					if !strings.Contains(catalogCreate, want) {
						t.Fatalf("catalog stored\n%s\nwithout %s", catalogCreate, want)
					}
				}
				for _, no := range tc.omits {
					if strings.Contains(engineCreate, no) {
						t.Fatalf("mysql stored\n%s\nwith %s", engineCreate, no)
					}
					if strings.Contains(catalogCreate, no) {
						t.Fatalf("catalog stored\n%s\nwith %s", catalogCreate, no)
					}
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
		t.Fatalf("%s: not a MySQL error: %v", stmt, err)
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

// TestDefineTableKeepsWhatOnlyItsServerHas: a snapshot describes a server that
// may ship a storage engine or a FULLTEXT parser the pinned one does not, and
// DefineTable installs it as it stands. A CREATE TABLE statement is judged
// against the pin instead.
func TestDefineTableKeepsWhatOnlyItsServerHas(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d;")

	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE s (b TEXT, FULLTEXT KEY ft (b) WITH PARSER mecab)")); err != nil {
		t.Fatalf("DefineTable with a parser only another server has: %v", err)
	}
	if tbl := c.GetDatabase("d").GetTable("s"); tbl == nil || len(tbl.Indexes) != 1 {
		t.Fatalf("DefineTable did not install the fulltext key: %+v", tbl)
	}
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE s2 (b INT) ENGINE = MyRocks")); err != nil {
		t.Fatalf("DefineTable with an engine only another server has: %v", err)
	}

	for _, stmt := range []string{
		"CREATE TABLE s3 (b TEXT, FULLTEXT KEY ft (b) WITH PARSER mecab)",
		"CREATE TABLE s4 (b INT) ENGINE = MyRocks",
	} {
		results, err := c.Exec(stmt, nil)
		if err != nil {
			t.Fatalf("Exec(%q): %v", stmt, err)
		}
		if results[0].Error == nil {
			t.Errorf("Exec(%q) took what the pinned server does not have", stmt)
		}
	}
}
