package catalog

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// TestSimulationCausesOracle holds the refusals testdata/simulate.tsv records to
// mysql:8.0.46 cause by cause: the object an ALTER or DROP VIEW names, a key
// built on a column the table does not have, a FULLTEXT key's parser, a
// partition action on a table that has no partitions, a type's display width, a
// view body naming a column no relation has or expanding a star no table backs,
// a CREATE TABLE with no column, a NOW()-family DEFAULT the column's type cannot
// hold, the columns MySQL gives another type from the text they were: SERIAL,
// CHAR BYTE and CHAR COLLATE BINARY, and the keyword a DROP or ALTER wrote
// before the constraint it names, DROP CHECK / ALTER CHECK beside DROP
// CONSTRAINT / ALTER CONSTRAINT, which decides the code a missing one answers:
// 3821 for CHECK and 3940 for CONSTRAINT. Each case is answered by the pinned
// server and by the catalog, and both must agree with the code the case states.
func TestSimulationCausesOracle(t *testing.T) {
	conn := testcontainer.OpenEmpty(t)
	ctx := context.Background()
	reset := func() {
		for _, stmt := range []string{"SET SESSION sql_mode = DEFAULT", "DROP VIEW IF EXISTS a", "DROP TABLE IF EXISTS a", "DROP TABLE IF EXISTS p", "DROP TABLE IF EXISTS s", "DROP TABLE IF EXISTS z", "DROP VIEW IF EXISTS v", "SET SESSION sql_generate_invisible_primary_key = OFF", "DROP DATABASE IF EXISTS cteother", "DROP TABLE IF EXISTS `dual`", "DROP TABLE IF EXISTS c2", "DROP EVENT IF EXISTS e"} {
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
		// keeps is the partitioning type the catalog must hold after a
		// statement it took, "" for one it may not.
		keeps string
		// cleanup runs on the server after the subtest finishes, for instance-wide
		// objects such as general tablespaces that DROP TABLE / reset do not clear.
		cleanup []string
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
		// The keyword a DROP or ALTER wrote before a constraint is what the
		// engine answers a missing one with: 3821 for CHECK, 3940 for
		// CONSTRAINT. It also says what the name may be: DROP CHECK takes a
		// check alone, DROP CONSTRAINT any constraint, and DROP FOREIGN KEY a
		// foreign key alone (measured on mysql:8.0.46).
		{
			name:  "DROP CHECK of a name no check constraint has",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DROP CHECK nope",
			code:  3821,
		},
		{
			name:  "DROP CONSTRAINT of a name no constraint has",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a DROP CONSTRAINT nope",
			code:  3940,
		},
		{
			name: "DROP CHECK of a foreign key's name",
			setup: []string{
				"CREATE TABLE p (id INT PRIMARY KEY)",
				"CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT fk1 FOREIGN KEY (c) REFERENCES p(id))",
			},
			stmt:     "ALTER TABLE a DROP CHECK fk1",
			code:     3821,
			contains: []string{"fk1"},
		},
		{
			name: "DROP FOREIGN KEY of a check's name",
			setup: []string{
				"CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT ck1 CHECK (c > 0))",
			},
			stmt:     "ALTER TABLE a DROP FOREIGN KEY ck1",
			code:     1091,
			contains: []string{"ck1"},
		},
		{
			name: "DROP CONSTRAINT of a check's name",
			setup: []string{
				"CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT ck1 CHECK (c > 0))",
			},
			stmt:  "ALTER TABLE a DROP CONSTRAINT ck1",
			code:  0,
			omits: []string{"ck1"},
		},
		{
			name:  "ALTER CHECK of a name no check constraint has",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER CHECK nope ENFORCED",
			code:  3821,
		},
		{
			name:  "ALTER CONSTRAINT of a name no constraint has",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER CONSTRAINT nope NOT ENFORCED",
			code:  3940,
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
		{
			// A partition function that reads a column the table does not hold
			// answers before the unique keys are judged: 1488 for KEY's list,
			// 1054 for a partitioning expression, and 1492 for a RANGE or LIST
			// partitioning that defines no partition, which answers even on a
			// table that is not there.
			name:  "ALTER TABLE PARTITION BY KEY over a column the table does not hold",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY(zzz)",
			code:  1488,
		},
		{
			name:  "ALTER TABLE PARTITION BY HASH over a column the table does not hold",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY HASH(zzz)",
			code:  1054,
		},
		{
			name:  "ALTER TABLE PARTITION BY KEY naming a column twice",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY(c, c)",
			code:  1652,
		},
		{
			// MySQL mixes subpartitioning with RANGE and LIST only, and answers
			// 1500 before it reads the column the subpartition function names.
			name:  "ALTER TABLE PARTITION BY KEY with a SUBPARTITION BY",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(zzz)",
			code:  1500,
		},
		{
			name:  "CREATE TABLE PARTITION BY HASH with a SUBPARTITION BY",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT PRIMARY KEY, c INT) PARTITION BY HASH(b) PARTITIONS 2 SUBPARTITION BY HASH(c)",
			code:  1500,
		},
		{
			name:  "ALTER TABLE PARTITION BY RANGE with no partition",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(c)",
			code:  1492,
		},
		{
			name: "ALTER TABLE PARTITION BY RANGE with no partition on no table",
			stmt: "ALTER TABLE a PARTITION BY RANGE(c)",
			code: 1492,
		},
		{
			// The forms that read a column answer 1146 first.
			name: "ALTER TABLE PARTITION BY HASH over a column on no table",
			stmt: "ALTER TABLE a PARTITION BY HASH(c)",
			code: 1146,
		},
		{
			name:  "CREATE TABLE PARTITION BY KEY over a column the table does not hold",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT PRIMARY KEY, c INT) PARTITION BY KEY(zzz)",
			code:  1488,
		},
		{
			// A partitioning the statement writes is the one the table keeps.
			name:  "ALTER TABLE PARTITION BY KEY over the primary key",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY()",
			code:  0,
			keeps: "KEY",
		},
		{
			name:  "ALTER TABLE PARTITION BY LINEAR KEY with a partition",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY LINEAR KEY() PARTITIONS 1",
			code:  0,
			keeps: "KEY",
		},
		{
			name:  "ALTER TABLE PARTITION BY RANGE over the primary key",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10))",
			code:  0,
			keeps: "RANGE",
		},
		{
			// A SPATIAL key's column is judged for its type before its
			// nullability: an INT column is 1687 whether or not it is NOT NULL,
			// and a nullable geometry column is what answers 1252.
			name:  "ADD SPATIAL over a nullable column of another type",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD SPATIAL(c)",
			code:  1687,
		},
		{
			name:  "ADD SPATIAL over a NOT NULL column of another type",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD SPATIAL(b)",
			code:  1687,
		},
		{
			name:  "ADD SPATIAL over a nullable geometry column",
			setup: []string{oneTable, "ALTER TABLE a ADD g POINT"},
			stmt:  "ALTER TABLE a ADD SPATIAL(g)",
			code:  1252,
		},
		{
			name:  "ADD SPATIAL over a NOT NULL geometry column",
			setup: []string{oneTable, "ALTER TABLE a ADD g POINT NOT NULL"},
			stmt:  "ALTER TABLE a ADD SPATIAL(g)",
			code:  0,
		},
		{
			name:  "ADD KEY USING RTREE over a column of another type",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY USING RTREE(c)",
			code:  1687,
		},
		{
			// USING RTREE makes the key a spatial one, which MySQL stores and
			// prints as SPATIAL KEY.
			name:     "ADD KEY USING RTREE over a NOT NULL geometry column",
			setup:    []string{oneTable, "ALTER TABLE a ADD g POINT NOT NULL"},
			stmt:     "ALTER TABLE a ADD KEY USING RTREE(g)",
			code:     0,
			contains: []string{"SPATIAL KEY `g` (`g`)"},
		},
		{
			name:  "CREATE TABLE with a SPATIAL key over a column of another type",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE s (b INT NOT NULL, SPATIAL KEY (b))",
			code:  1687,
		},
		{
			// A spatial reference system MySQL defines must write both the
			// attributes it requires, and answers 3708 naming the one it misses.
			name:  "CREATE SPATIAL REFERENCE SYSTEM with no attribute",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 5000",
			code:  3708,
		},
		{
			name:  "CREATE SPATIAL REFERENCE SYSTEM with no NAME",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 5000 DEFINITION 's'",
			code:  3708,
		},
		{
			name:  "CREATE SPATIAL REFERENCE SYSTEM with no DEFINITION",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 5000 NAME 's'",
			code:  3708,
		},
		{
			name:  "CREATE OR REPLACE SPATIAL REFERENCE SYSTEM with no NAME",
			setup: []string{oneTable},
			stmt:  "CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 5000 ORGANIZATION 's' IDENTIFIED BY 1",
			code:  3708,
		},
		{
			// An attribute the statement writes with an empty literal is present,
			// and MySQL judges its value instead of calling it missing.
			name:  "CREATE SPATIAL REFERENCE SYSTEM with an empty NAME",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 5000 NAME ''",
			code:  3710,
		},
		{
			name:  "CREATE SPATIAL REFERENCE SYSTEM with a NAME that starts with a space",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 5000 NAME ' s'",
			code:  3710,
		},
		{
			name:  "CREATE SPATIAL REFERENCE SYSTEM with an empty DEFINITION",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 5000 NAME 's' DEFINITION ''",
			code:  3517,
		},
		{
			name:  "CREATE SPATIAL REFERENCE SYSTEM on SRID 0",
			setup: []string{oneTable},
			stmt:  "CREATE SPATIAL REFERENCE SYSTEM 0 NAME 's' DEFINITION 's'",
			code:  3714,
		},
		// A view body's star is resolved before the view's own name: a `*`
		// with no FROM expands from nothing and is 1096, even over a table a,
		// where the name alone would be 1050.
		{
			name:  "CREATE VIEW over a table whose body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW a AS SELECT *",
			code:  1096,
		},
		{
			name:  "CREATE VIEW whose body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT *",
			code:  1096,
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
		// A table with no column at all is 4028, after the name collision and
		// after a key's own columns are resolved.
		{
			name:  "CREATE TABLE with no columns over a table",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a;",
			code:  1050,
		},
		{
			name: "CREATE TABLE with no columns",
			stmt: "CREATE TABLE z;",
			code: 4028,
		},
		{
			name: "CREATE TABLE whose only element is a key over a missing column",
			stmt: "CREATE TABLE z (KEY (b));",
			code: 1072,
		},
		{
			name: "CREATE TABLE whose only element is a partition clause",
			stmt: "CREATE TABLE z PARTITION BY KEY();",
			code: 4028,
		},
		{
			name:  "CREATE VIEW whose derived table expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT * FROM (SELECT *) d",
			code:  1096,
		},
		{
			name:  "CREATE VIEW whose CTE body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM c",
			code:  1096,
		},
		{
			// MySQL resolves a CTE only where the body reads it, so a body it
			// never names is not judged.
			name:  "CREATE VIEW whose unreferenced CTE body expands an empty star",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *) SELECT 1",
			code:  0,
		},
		{
			name:  "CREATE VIEW whose body expands a star over DUAL",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT * FROM DUAL",
			code:  1096,
		},
		{
			// A table really named dual is written quoted, which is not the
			// dummy keyword MySQL expands from nothing.
			name:  "CREATE VIEW whose body expands a star over a quoted dual table",
			setup: []string{"CREATE TABLE `dual` (x INT)"},
			stmt:  "CREATE VIEW v AS SELECT * FROM `dual`",
			code:  0,
		},
		{
			// A CTE read only by another CTE nothing reads is still unreached.
			name:  "CREATE VIEW whose unreferenced CTE reads another",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c1 AS (SELECT *), c2 AS (SELECT * FROM c1) SELECT 1",
			code:  0,
		},
		{
			// A non-recursive CTE does not see itself, so this inner c reads the
			// outer one, whose body expands an empty star.
			name:  "CREATE VIEW whose non-recursive CTE names itself",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *), d AS (WITH c AS (SELECT * FROM c) SELECT * FROM c) SELECT * FROM d",
			code:  1096,
		},
		{
			name:  "CREATE VIEW whose recursive CTE names itself",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH RECURSIVE c AS (SELECT 1 AS x) SELECT * FROM c",
			code:  0,
		},
		{
			// A CTE does not see a sibling declared after it, so the name reads
			// the real table c2 and the later CTE stays unreached.
			name:  "CREATE VIEW whose CTE reads a later sibling's name",
			setup: []string{oneTable, "CREATE TABLE c2 (x INT)"},
			stmt:  "CREATE VIEW v AS WITH c1 AS (SELECT * FROM c2), c2 AS (SELECT *) SELECT * FROM c1",
			code:  0,
		},
		{
			// The inner block's CTE shadows the outer one, so the outer CTE is
			// never read and MySQL never resolves it.
			name:  "CREATE VIEW whose nested CTE shadows an unreferenced one",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM (WITH c AS (SELECT 1 AS x) SELECT x FROM c) d",
			code:  0,
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
			// The pin runs with lower_case_table_names=0: a name and an alias are
			// compared as written.
			name:  "CREATE VIEW whose qualified star differs in case from its alias",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT T.* FROM a AS t",
			code:  1051,
		},
		{
			name:  "CREATE VIEW whose qualified star names a JSON_TABLE alias",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT jt.* FROM JSON_TABLE('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt",
			code:  0,
		},
		{
			name:  "CREATE VIEW whose qualified star names its relation's alias",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT t.* FROM a AS t",
			code:  0,
		},
		{
			// The dummy DUAL reads no CTE, whatever a CTE of that name is.
			name:  "CREATE VIEW whose quoted dual CTE the dummy DUAL does not read",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS WITH `dual` AS (SELECT *) SELECT 1 FROM DUAL",
			code:  0,
		},
		{
			name:  "CREATE VIEW whose body reads only DUAL for a column",
			setup: []string{oneTable},
			stmt:  "CREATE VIEW v AS SELECT 1 FROM DUAL",
			code:  0,
		},
		{
			// `other.c` names another schema's table, not the CTE c, so MySQL
			// never resolves the CTE's star.
			name:  "CREATE VIEW whose qualified name shares an unreferenced CTE's",
			setup: []string{"CREATE DATABASE cteother", "CREATE TABLE cteother.c (x INT)"},
			stmt:  "CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM cteother.c",
			code:  0,
		},
		{
			name:  "CREATE TABLE LIKE over a table",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE z (LIKE a);",
			code:  0,
		},
		// A bare NOW()-family DEFAULT is a current-timestamp value, which only
		// a DATETIME or TIMESTAMP column holds. MySQL judges it while it
		// prepares the statement: ahead of the table's lookup (1146) and of a
		// duplicate column (1060). A parenthesized DEFAULT (NOW()) is a general
		// expression default, which the same server takes on an INT column.
		{
			name: "ADD a NOW() default over no table",
			stmt: "ALTER TABLE nope ADD e INT DEFAULT NOW()",
			code: 1067,
		},
		{
			name:  "ADD a CURRENT_TIMESTAMP default over a duplicate column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD b INT DEFAULT CURRENT_TIMESTAMP",
			code:  1067,
		},
		{
			name:  "MODIFY a column to a NOW() default",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a MODIFY b INT DEFAULT NOW()",
			code:  1067,
		},
		{
			name:  "ADD a parenthesized NOW() default to an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e INT DEFAULT (NOW())",
			code:  0,
		},
		{
			name:  "ADD a NOW() default to a DATETIME column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD e DATETIME DEFAULT NOW()",
			code:  0,
		},
		{
			name: "CREATE TABLE with a NOW() default on an INT column",
			stmt: "CREATE TABLE z (a INT DEFAULT NOW());",
			code: 1067,
		},
		{
			name: "CREATE TABLE with a NOW() default on a DATETIME column",
			stmt: "CREATE TABLE z (a DATETIME DEFAULT NOW());",
			code: 0,
		},
		{
			// 4028 asks for a column a reader may see, not merely for a column.
			name: "CREATE TABLE whose only column is invisible",
			stmt: "CREATE TABLE z (a INT INVISIBLE);",
			code: 4028,
		},
		{
			name:  "CREATE TABLE whose only column is invisible, with a generated primary key",
			setup: []string{"SET SESSION sql_generate_invisible_primary_key = ON"},
			stmt:  "CREATE TABLE z (a INT INVISIBLE);",
			code:  4028,
		},
		{
			name: "CREATE TABLE with an invisible and a visible column",
			stmt: "CREATE TABLE z (a INT INVISIBLE, b INT);",
			code: 0,
		},
		{
			name:  "ALTER COLUMN SET INVISIBLE on the last visible column",
			setup: []string{"CREATE TABLE z (a INT, b INT INVISIBLE)"},
			stmt:  "ALTER TABLE z ALTER COLUMN a SET INVISIBLE",
			code:  4028,
		},
		{
			// A functional index's hidden system column is not a visible one.
			name:  "ALTER COLUMN SET INVISIBLE beside a functional index",
			setup: []string{"CREATE TABLE z (a VARCHAR(10), INDEX ((LOWER(a))))"},
			stmt:  "ALTER TABLE z ALTER COLUMN a SET INVISIBLE",
			code:  4028,
		},
		{
			// The visible column is judged over the table the statement leaves.
			name:  "hide a column and add another in one ALTER",
			setup: []string{"CREATE TABLE z (a INT, b INT INVISIBLE)"},
			stmt:  "ALTER TABLE z ALTER COLUMN a SET INVISIBLE, ADD COLUMN c INT",
			code:  0,
		},
		{
			name:  "hide two columns in one ALTER",
			setup: []string{"CREATE TABLE z (a INT, b INT)"},
			stmt:  "ALTER TABLE z ALTER COLUMN a SET INVISIBLE, ALTER COLUMN b SET INVISIBLE",
			code:  4028,
		},
		{
			name:  "MODIFY the last visible column INVISIBLE",
			setup: []string{"CREATE TABLE z (a INT, b INT INVISIBLE)"},
			stmt:  "ALTER TABLE z MODIFY a INT INVISIBLE",
			code:  4028,
		},
		{
			// A CTAS appends the query's own columns as visible ones.
			name: "CREATE TABLE AS SELECT beside an invisible column",
			stmt: "CREATE TABLE z (a INT INVISIBLE) AS SELECT 1 AS b;",
			code: 0,
		},
		{
			// MySQL resolves a collation name while it parses, so a name it does
			// not have is 1273 before the table is looked up: utf8mb4 is a
			// character set and no collation.
			name:  "COLLATE a character set name as a collation",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COLLATE utf8mb4",
			code:  1273,
		},
		{
			name:  "COLLATE an unknown collation on a column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD c2 INT COLLATE utf8mb4",
			code:  1273,
		},
		{
			// The lookup answers before the trailing token a statement reads past
			// the name: this text does not parse, and MySQL still answers 1273.
			name:  "an unknown collation before a syntax error",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COLLATE utf8mb4 zzz",
			code:  1273,
		},
		{
			// A collation the server does have is read, and the trailing zzz is
			// the 1064 it is.
			name:  "a known collation before a syntax error",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COLLATE utf8mb4_bin zzz",
			code:  1064,
		},
		{
			// A CTAS body is read after the table's own name: a table that is not
			// there leaves the body's column to answer 1054.
			name: "CREATE TABLE AS SELECT a column no table has",
			stmt: "CREATE TABLE z SELECT nosuchcol",
			code: 1054,
		},
		{
			// A star with no table to expand is 1096, as it is in a view body.
			name: "CREATE TABLE AS SELECT a star with no table",
			stmt: "CREATE TABLE z AS SELECT *",
			code: 1096,
		},
		{
			name:  "CREATE TABLE AS SELECT over the table it makes",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a SELECT nosuchcol",
			code:  1050,
		},
		{
			name: "CREATE TABLE AS SELECT a body it takes",
			stmt: "CREATE TABLE z AS SELECT 1 AS b",
			code: 0,
		},
		{
			// Declaring the columns does not keep the body from being read.
			name: "CREATE TABLE with columns AS SELECT a column no table has",
			stmt: "CREATE TABLE z (a INT) AS SELECT nosuchcol",
			code: 1054,
		},
		{
			name: "CREATE TABLE with columns AS SELECT a body it takes",
			stmt: "CREATE TABLE z (a INT) AS SELECT 1",
			code: 0,
		},
		{
			// An event's schedule sees no table, so a column reference names
			// nothing: 1054, before the event's own name is looked up.
			name:  "an event schedule naming a column",
			setup: []string{oneTable},
			stmt:  "CREATE EVENT e ON SCHEDULE AT nosuchcol DO COMMIT",
			code:  1054,
		},
		{
			name:  "an event schedule it takes",
			setup: []string{oneTable},
			stmt:  "CREATE DEFINER = root EVENT e ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY DO COMMIT",
			code:  0,
		},
		{
			// ALTER EVENT reads the schedule the same way, before the event.
			name:  "ALTER EVENT's schedule naming a column",
			setup: []string{oneTable},
			stmt:  "ALTER EVENT e ON SCHEDULE AT nosuchcol DO COMMIT",
			code:  1054,
		},
		{
			name:  "ALTER EVENT's schedule over an event that is not there",
			setup: []string{oneTable},
			stmt:  "ALTER EVENT nosuchevent ON SCHEDULE AT nosuchcol DO COMMIT",
			code:  1054,
		},
		{
			name:  "ALTER EVENT's schedule it takes",
			setup: []string{oneTable, "CREATE EVENT e ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY DO COMMIT"},
			stmt:  "ALTER EVENT e ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 2 DAY DO COMMIT",
			code:  0,
		},
		{
			// InnoDB refuses ROW_FORMAT FIXED and COMPRESSION, and a table that
			// names no engine is InnoDB's.
			name:  "ROW_FORMAT FIXED on InnoDB",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ROW_FORMAT FIXED",
			code:  1478,
		},
		{
			name:  "COMPRESSION on InnoDB",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COMPRESSION 's'",
			code:  1478,
		},
		{
			// An empty COMPRESSION is a value InnoDB takes: it is the option
			// written with no value, not a value the engine refuses.
			name:  "an empty COMPRESSION on InnoDB",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COMPRESSION=''",
			code:  0,
		},
		{
			// ZLIB, LZ4 and NONE are InnoDB's page compression and its off
			// switch, so the option is refused on the value and not on the name.
			name:  "COMPRESSION ZLIB on InnoDB",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COMPRESSION='ZLIB'",
			code:  0,
		},
		{
			name:  "COMPRESSION none on InnoDB",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COMPRESSION='none'",
			code:  0,
		},
		{
			// The row option is judged after the statement's other commands.
			name:  "ROW_FORMAT FIXED beside a key on a column that is not there",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ROW_FORMAT FIXED, ADD KEY(zzz)",
			code:  1072,
		},
		{
			// The engine is the one the statement leaves, so MyISAM takes it and
			// MySQL stores the table it makes.
			name:     "ROW_FORMAT FIXED on MyISAM",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ENGINE=MyISAM, ROW_FORMAT FIXED",
			code:     0,
			contains: []string{"ENGINE=MyISAM"},
		},
		// A CHECK expression whose top level is not boolean is 3812 while the
		// statement is prepared, before the column it names is resolved (a bare
		// unknown column is this, not 1054) and before a duplicate column.
		{
			name:  "ADD CHECK over a bare column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD CHECK(zzz)",
			code:  3812,
		},
		{
			name:  "ADD a column whose CHECK is the column itself",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD zzz INT CHECK(zzz)",
			code:  3812,
		},
		{
			name:  "ADD CHECK over a comparison",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD CHECK(b > 0)",
			code:  0,
		},
		// A generated column whose expression names the column being added is
		// 3107 after the duplicate-column check, not 1054.
		{
			name:  "ADD a generated column that names itself",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD zzz INT AS(zzz)",
			code:  3107,
		},
		{
			name:  "ADD a generated column over a base column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD zzz INT AS(b + 1)",
			code:  0,
		},
		// ALTER INDEX and RENAME KEY over a key the table does not hold answer
		// 1176, not the 1061 a duplicate name answers.
		{
			name:  "ALTER INDEX over a key the table does not hold",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER INDEX zzz VISIBLE",
			code:  1176,
		},
		{
			name:  "RENAME KEY over a key the table does not hold",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a RENAME KEY zzz TO yyy",
			code:  1176,
		},
		// A SRID clause belongs to a geometry column, and a bare WITH VALIDATION
		// belongs to EXCHANGE PARTITION. Both are 1221 while the statement is
		// prepared, ahead of the table lookup.
		{
			name:  "ADD a non-geometry column with SRID",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD zzz INT SRID 1",
			code:  1221,
		},
		{
			name:  "ALTER TABLE WITH VALIDATION",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a WITH VALIDATION",
			code:  1221,
		},
		{
			name:  "ALTER TABLE WITHOUT VALIDATION",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a WITHOUT VALIDATION",
			code:  1221,
		},
		// CREATE TABLE ... LIKE naming its own source is 1066 before either
		// table is looked up: over the table itself, and over one that is not
		// there.
		{
			name:  "CREATE TABLE LIKE itself over the table",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a LIKE a",
			code:  1066,
		},
		{
			name: "CREATE TABLE LIKE itself",
			stmt: "CREATE TABLE a LIKE a",
			code: 1066,
		},
		// A key that lists the same column twice is 1060, after the columns
		// exist and after a prefix's own check.
		{
			name:  "ADD KEY over the same column twice",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(b, b)",
			code:  1060,
		},
		{
			name:  "ADD FOREIGN KEY over the same column twice",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(c, c) REFERENCES a(b)",
			code:  1060,
		},
		// ANALYZE, CHECK, OPTIMIZE and REPAIR PARTITION do not open the table
		// when each is the statement's only action.
		{
			name: "OPTIMIZE PARTITION over a table that is not there",
			stmt: "ALTER TABLE a OPTIMIZE PARTITION ALL",
			code: 0,
		},
		// START TRANSACTION is not an ALTER TABLE option, and an ENCRYPTION
		// value other than Y or N is 3184.
		{
			name:  "ALTER TABLE START TRANSACTION",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a START TRANSACTION",
			code:  3979,
		},
		{
			name:  "ALTER TABLE ENCRYPTION over a value it does not take",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENCRYPTION 's'",
			code:  3184,
		},
		// 1. GEOMETRYCOLLECTION / GEOMCOLLECTION column type aliases store as
		// geomcollection in SHOW CREATE TABLE (mysql:8.0.46).
		{
			name:     "ALTER TABLE ADD GEOMETRYCOLLECTION",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD g GEOMETRYCOLLECTION",
			code:     0,
			contains: []string{"`g` geomcollection DEFAULT NULL"},
		},
		{
			name:     "ALTER TABLE ADD GEOMCOLLECTION",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD g GEOMCOLLECTION",
			code:     0,
			contains: []string{"`g` geomcollection DEFAULT NULL"},
		},
		// 2. Functional index naming (`functional_index`), bare-column refusal
		// (3762 ahead of 1061 and 1054), primary-key functional refusal (3756 ahead
		// of 1068), and functional expression column resolution (1054 / forward ADD).
		{
			name:     "ALTER TABLE ADD KEY over a constant expression",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD KEY((1))",
			code:     0,
			contains: []string{"KEY `functional_index` ((1))"},
		},
		{
			name:     "ALTER TABLE ADD UNIQUE over an expression",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD UNIQUE((b + 1))",
			code:     0,
			contains: []string{"UNIQUE KEY `functional_index` ("},
		},
		{
			name:  "ALTER TABLE ADD KEY over a bare column expression",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY((b))",
			code:  3762,
		},
		{
			name:  "ALTER TABLE ADD KEY over a nested bare column expression",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(((b)))",
			code:  3762,
		},
		{
			name:  "ALTER TABLE ADD KEY over a bare column expression with a duplicate index name",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY idx_d ((b))",
			code:  3762,
		},
		{
			name:  "ALTER TABLE ADD KEY over an unknown bare column expression",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY((nosuch))",
			code:  3762,
		},
		{
			name:  "CREATE INDEX over a bare column expression with a duplicate index name",
			setup: []string{oneTable},
			stmt:  "CREATE INDEX idx_d ON a ((b))",
			code:  3762,
		},
		{
			name: "CREATE TABLE with a functional key over an unknown bare column",
			stmt: "CREATE TABLE z (b INT, KEY ((nosuch)))",
			code: 3762,
		},
		{
			name:  "ALTER TABLE ADD KEY over an unknown column in an expression",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY((nosuch + 1))",
			code:  1054,
		},
		{
			name:     "ALTER TABLE ADD KEY over a column a later command adds in an expression",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a ADD KEY((x2 + 1)), ADD x2 INT",
			code:     0,
			contains: []string{"KEY `functional_index` (", "`x2` int DEFAULT NULL"},
		},
		{
			name:  "ALTER TABLE ADD PRIMARY KEY over a functional expression on a table with a PK",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD PRIMARY KEY((b + 1))",
			code:  3756,
		},
		{
			name: "CREATE TABLE with a functional PRIMARY KEY",
			stmt: "CREATE TABLE z (b INT, PRIMARY KEY((b + 1)))",
			code: 3756,
		},
		// 3. CONVERT TO CHARSET DEFAULT resolves to the database's charset and
		// collation, clearing column-level overrides and re-checking ENUM/SET.
		{
			name: "CONVERT TO CHARSET DEFAULT",
			setup: []string{
				"CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=latin1 COLLATE=latin1_bin",
			},
			stmt:     "ALTER TABLE a CONVERT TO CHARSET DEFAULT",
			code:     0,
			contains: []string{"DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"},
			omits:    []string{"latin1"},
		},
		{
			name: "CONVERT TO CHARSET DEFAULT with an explicit COLLATE",
			setup: []string{
				"CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=latin1",
			},
			stmt:     "ALTER TABLE a CONVERT TO CHARSET DEFAULT COLLATE utf8mb4_bin",
			code:     0,
			contains: []string{"DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"},
			omits:    []string{"latin1"},
		},
		{
			name:  "CONVERT TO CHARSET DEFAULT with a mismatched COLLATE",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CONVERT TO CHARSET DEFAULT COLLATE latin1_bin",
			code:  1253,
		},
		{
			name: "CONVERT TO CHARSET DEFAULT colliding case-distinct ENUM values",
			setup: []string{
				"CREATE TABLE a (e ENUM('a', 'A') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
			},
			stmt: "ALTER TABLE a CONVERT TO CHARSET DEFAULT",
			code: 1291,
		},
		// 4. SECONDARY_LOAD and SECONDARY_UNLOAD require a defined SECONDARY_ENGINE (3889).
		{
			name:  "SECONDARY_LOAD with no secondary engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a SECONDARY_LOAD",
			code:  3889,
		},
		{
			name:  "ALGORITHM DEFAULT, SECONDARY_LOAD with no secondary engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALGORITHM DEFAULT, SECONDARY_LOAD",
			code:  3889,
		},
		{
			name:  "SECONDARY_UNLOAD with no secondary engine",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a SECONDARY_UNLOAD",
			code:  3889,
		},
		{
			name:     "SECONDARY_LOAD with a secondary engine defined",
			setup:    []string{"CREATE TABLE a (b INT PRIMARY KEY) SECONDARY_ENGINE = rapid"},
			stmt:     "ALTER TABLE a SECONDARY_LOAD",
			code:     0,
			contains: []string{"SECONDARY_ENGINE=rapid"},
		},
		{
			name:  "SECONDARY_UNLOAD with a secondary engine defined",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY) SECONDARY_ENGINE = rapid"},
			stmt:  "ALTER TABLE a SECONDARY_UNLOAD",
			code:  0,
		},
		{
			name: "SECONDARY_LOAD after clearing SECONDARY_ENGINE to NULL",
			setup: []string{
				"CREATE TABLE a (b INT PRIMARY KEY) SECONDARY_ENGINE = rapid",
				"ALTER TABLE a SECONDARY_ENGINE = NULL",
			},
			stmt: "ALTER TABLE a SECONDARY_LOAD",
			code: 3889,
		},
		// 5. TABLESPACE validation on InnoDB (3510, 3723, 3119, 1478) and CREATE TABLESPACE.
		{
			name:  "ALTER TABLE TABLESPACE over a missing tablespace",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TABLESPACE nosuch_ts",
			code:  3510,
		},
		{
			name:  "ALTER TABLE TABLESPACE judged before a missing key column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(nosuch_col), TABLESPACE nosuch_ts",
			code:  3510,
		},
		{
			name:  "CREATE TABLE TABLESPACE judged before a duplicate table name",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE a (b INT) TABLESPACE nosuch_ts",
			code:  3510,
		},
		{
			name:  "ALTER TABLE TABLESPACE mysql reserved",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TABLESPACE mysql",
			code:  3723,
		},
		{
			name:  "ALTER TABLE TABLESPACE innodb_ prefix reserved",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TABLESPACE innodb_custom",
			code:  3119,
		},
		{
			name:  "ALTER TABLE TABLESPACE innodb_temporary on a non-temporary table",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a TABLESPACE innodb_temporary",
			code:  1478,
		},
		{
			name:     "ALTER TABLE TABLESPACE innodb_file_per_table",
			setup:    []string{oneTable},
			stmt:     "ALTER TABLE a TABLESPACE innodb_file_per_table",
			code:     0,
			contains: []string{"TABLESPACE `innodb_file_per_table`"},
		},
		{
			name:  "ALTER TABLE TABLESPACE ignored on MyISAM",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE=MyISAM, TABLESPACE nosuch_ts",
			code:  0,
		},
		{
			name: "ALTER TABLE TABLESPACE over a created general tablespace",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
				oneTable,
			},
			stmt:     "ALTER TABLE a TABLESPACE ts_oracle",
			code:     0,
			contains: []string{"TABLESPACE `ts_oracle`"},
			cleanup:  []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_oracle"},
		},
		{
			name: "ALTER TABLE TABLESPACE is case-sensitive on general tablespaces",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
				oneTable,
			},
			stmt:    "ALTER TABLE a TABLESPACE TS_ORACLE",
			code:    3510,
			cleanup: []string{"DROP TABLESPACE ts_oracle"},
		},
		// 6. AUTOEXTEND_SIZE range [4M, 4096M] (4025) and 4M alignment (4023) on InnoDB.
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE below 4M",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a AUTOEXTEND_SIZE 1",
			code:  4025,
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE above 4096M",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a AUTOEXTEND_SIZE 8192M",
			code:  4025,
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE not a multiple of 4M",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a AUTOEXTEND_SIZE 5M",
			code:  4023,
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE 0",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a AUTOEXTEND_SIZE 0",
			code:  0,
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE 4M",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a AUTOEXTEND_SIZE 4M",
			code:  0,
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE judged after a missing key column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(nosuch_col), AUTOEXTEND_SIZE 1",
			code:  1072,
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE ignored on MyISAM",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ENGINE=MyISAM, AUTOEXTEND_SIZE 1",
			code:  0,
		},
		{
			name: "CREATE TABLE with AUTOEXTEND_SIZE below 4M",
			stmt: "CREATE TABLE z (b INT) AUTOEXTEND_SIZE 1",
			code: 4025,
		},
		// 7. DISCARD / IMPORT TABLESPACE (1813, 1812, 1031, 1148, and TRUNCATE 1814).
		{
			name:  "IMPORT TABLESPACE when tablespace has not been discarded",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a IMPORT TABLESPACE",
			code:  1813,
		},
		{
			name:  "IMPORT TABLESPACE after DISCARD TABLESPACE with no .ibd file",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a IMPORT TABLESPACE",
			code:  1812,
		},
		{
			name:  "TRUNCATE TABLE after DISCARD TABLESPACE",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "TRUNCATE TABLE a",
			code:  1814,
		},
		{
			name:  "IMPORT TABLESPACE on a MyISAM table",
			setup: []string{"CREATE TABLE a (b INT) ENGINE=MyISAM"},
			stmt:  "ALTER TABLE a IMPORT TABLESPACE",
			code:  1031,
		},
		{
			name: "DISCARD TABLESPACE on a table in a general tablespace",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
				"CREATE TABLE a (b INT) TABLESPACE ts_oracle",
			},
			stmt:    "ALTER TABLE a DISCARD TABLESPACE",
			code:    1148,
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_oracle"},
		},
		// 8. FOREIGN KEY without referenced columns or with mismatched arity (1239).
		{
			name:  "ALTER TABLE ADD FOREIGN KEY with no referenced columns",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES a",
			code:  1239,
		},
		{
			name:  "ALTER TABLE ADD FOREIGN KEY with mismatched referenced column count",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(b, c) REFERENCES a(b)",
			code:  1239,
		},
		{
			name:  "CREATE TABLE with FOREIGN KEY and no referenced columns",
			setup: []string{oneTable},
			stmt:  "CREATE TABLE z (b INT, FOREIGN KEY(b) REFERENCES a)",
			code:  1239,
		},
		{
			name:    "ALTER TABLE ADD FOREIGN KEY with no referenced columns under foreign_key_checks=0",
			setup:   []string{oneTable},
			session: []string{"SET SESSION foreign_key_checks = 0"},
			stmt:    "ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES nosuch_tbl",
			code:    1239,
			cleanup: []string{"SET SESSION foreign_key_checks = ON"},
		},
		{
			name: "ALTER TABLE ADD FOREIGN KEY with no referenced columns judged before duplicate FK name",
			setup: []string{
				"CREATE TABLE p (id INT PRIMARY KEY)",
				"CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT fk1 FOREIGN KEY (c) REFERENCES p(id))",
			},
			stmt: "ALTER TABLE a ADD CONSTRAINT fk1 FOREIGN KEY(c) REFERENCES p",
			code: 1239,
		},
		// 9. Key prefix length on non-prefixable types (1089) and JSON key (3152).
		{
			name:  "ALTER TABLE ADD KEY with a prefix on an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD KEY(c(1))",
			code:  1089,
		},
		{
			name:  "ALTER TABLE ADD PRIMARY KEY with a prefix on an INT column",
			setup: []string{noPrimaryKey},
			stmt:  "ALTER TABLE a ADD PRIMARY KEY(c(1))",
			code:  1089,
		},
		{
			name:  "ALTER TABLE ADD FOREIGN KEY with a prefix on an INT column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD FOREIGN KEY(c(1)) REFERENCES a(b)",
			code:  1089,
		},
		{
			name:  "CREATE INDEX with a prefix on an INT column",
			setup: []string{oneTable},
			stmt:  "CREATE INDEX i ON a(c(1))",
			code:  1089,
		},
		{
			name: "CREATE TABLE with a prefix key on a DATETIME column",
			stmt: "CREATE TABLE z (c DATETIME, KEY(c(1)))",
			code: 1089,
		},
		{
			name: "CREATE TABLE with a prefix key on an ENUM column",
			stmt: "CREATE TABLE z (c ENUM('a','b'), KEY(c(1)))",
			code: 1089,
		},
		{
			name: "CREATE TABLE with a prefix on a SPATIAL KEY",
			stmt: "CREATE TABLE z (g POINT NOT NULL, SPATIAL KEY(g(1)))",
			code: 1089,
		},
		{
			name: "CREATE TABLE with a prefix on a BIT column",
			stmt: "CREATE TABLE z (b BIT(8), KEY(b(4)))",
			code: 1089,
		},
		{
			name: "CREATE TABLE with a prefix on a YEAR column",
			stmt: "CREATE TABLE z (b YEAR, KEY(b(2)))",
			code: 1089,
		},
		{
			name: "CREATE TABLE with a valid prefix on a VARCHAR column",
			stmt: "CREATE TABLE z (s VARCHAR(8), KEY(s(4)))",
			code: 0,
		},
		{
			name: "CREATE TABLE with a regular key on a JSON column",
			stmt: "CREATE TABLE z (j JSON, KEY(j))",
			code: 3152,
		},
		// 10. ON UPDATE CURRENT_TIMESTAMP prepare-time checks (1294) and FSP mismatch.
		{
			name: "ADD INT ON UPDATE CURRENT_TIMESTAMP over no table",
			stmt: "ALTER TABLE nope ADD x INT ON UPDATE CURRENT_TIMESTAMP",
			code: 1294,
		},
		{
			name:  "ADD INT ON UPDATE CURRENT_TIMESTAMP over a duplicate column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD b INT ON UPDATE CURRENT_TIMESTAMP",
			code:  1294,
		},
		{
			name: "ADD DATETIME(3) ON UPDATE CURRENT_TIMESTAMP with mismatched FSP over no table",
			stmt: "ALTER TABLE nope ADD x DATETIME(3) ON UPDATE CURRENT_TIMESTAMP",
			code: 1294,
		},
		{
			name: "ADD DATETIME(3) DEFAULT CURRENT_TIMESTAMP with mismatched FSP over no table",
			stmt: "ALTER TABLE nope ADD x DATETIME(3) DEFAULT CURRENT_TIMESTAMP",
			code: 1067,
		},
		{
			name:  "ADD DATETIME(3) with matching FSP on DEFAULT and ON UPDATE",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)",
			code:  0,
		},
		// 11. AUTO_INCREMENT key requirement on CREATE TABLE, ALTER TABLE and DROP INDEX (1075).
		{
			name:  "ALTER TABLE ADD unkeyed AUTO_INCREMENT column",
			setup: []string{noPrimaryKey},
			stmt:  "ALTER TABLE a ADD x INT AUTO_INCREMENT",
			code:  1075,
		},
		{
			name:  "ALTER TABLE ADD AUTO_INCREMENT column keyed by a later command",
			setup: []string{noPrimaryKey},
			stmt:  "ALTER TABLE a ADD x INT AUTO_INCREMENT, ADD KEY(x)",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD AUTO_INCREMENT column as non-leading key part on InnoDB",
			setup: []string{noPrimaryKey},
			stmt:  "ALTER TABLE a ADD x INT AUTO_INCREMENT, ADD KEY(c, x)",
			code:  1075,
		},
		{
			name: "CREATE TABLE with AUTO_INCREMENT column as secondary key part on MyISAM",
			stmt: "CREATE TABLE z (c INT NOT NULL, x INT AUTO_INCREMENT, PRIMARY KEY (c, x)) ENGINE=MyISAM",
			code: 0,
		},
		{
			name:  "DROP INDEX of the only key covering an AUTO_INCREMENT column",
			setup: []string{"CREATE TABLE a (b INT AUTO_INCREMENT, c INT, KEY k (b))"},
			stmt:  "DROP INDEX k ON a",
			code:  1075,
		},
		{
			name:  "ALTER TABLE DROP INDEX of the only key covering an AUTO_INCREMENT column",
			setup: []string{"CREATE TABLE a (b INT AUTO_INCREMENT, c INT, KEY k (b))"},
			stmt:  "ALTER TABLE a DROP INDEX k",
			code:  1075,
		},
		// 12. ENUM and SET duplicate member checks under strict sql_mode and collation (1291).
		{
			name:  "ALTER TABLE ADD ENUM with exact duplicate values",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('s', 's')",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with case-colliding values under default collation",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('s', 'S')",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with trailing-space duplicate under binary collation",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('s ', 's') COLLATE utf8mb4_bin",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with case-distinct values under binary collation",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('s', 'S') COLLATE utf8mb4_bin",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with COLLATE after NOT NULL",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('s', 'S') NOT NULL COLLATE utf8mb4_bin",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD SET with case-colliding values",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x SET('a', 'A')",
			code:  1291,
		},
		{
			name:    "ALTER TABLE ADD ENUM with duplicate values under non-strict sql_mode",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "ALTER TABLE a ADD x ENUM('s', 's')",
			code:    0,
		},
		// 13. Expression DEFAULT `(expr)` validation (3767, 3768, 1054, 3769, 3772).
		{
			name:  "ALTER TABLE ADD column with self-referencing expression DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT (x)",
			code:  3767,
		},
		{
			name:  "ALTER TABLE ALTER COLUMN SET self-referencing expression DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER COLUMN c SET DEFAULT (c)",
			code:  3767,
		},
		{
			name: "CREATE TABLE with forward expression DEFAULT referencing an expr-default column",
			stmt: "CREATE TABLE z (x INT DEFAULT (y), y INT DEFAULT (1))",
			code: 3767,
		},
		{
			name: "CREATE TABLE with forward expression DEFAULT referencing a generated column",
			stmt: "CREATE TABLE z (x INT DEFAULT (y), y INT AS (1))",
			code: 3767,
		},
		{
			name: "CREATE TABLE with forward expression DEFAULT referencing an ordinary column",
			stmt: "CREATE TABLE z (x INT DEFAULT (y), y INT)",
			code: 0,
		},
		{
			name: "CREATE TABLE with backward expression DEFAULT referencing an expr-default column",
			stmt: "CREATE TABLE z (y INT DEFAULT (1), x INT DEFAULT (y))",
			code: 0,
		},
		{
			name: "CREATE TABLE with forward expression DEFAULT referencing a bare NOW() column",
			stmt: "CREATE TABLE z (x INT DEFAULT (dt), dt DATETIME DEFAULT NOW())",
			code: 0,
		},
		{
			name: "CREATE TABLE with forward expression DEFAULT referencing a parenthesized (NOW()) column",
			stmt: "CREATE TABLE z (x INT DEFAULT (dt), dt DATETIME DEFAULT (NOW()))",
			code: 3767,
		},
		{
			name: "CREATE TABLE with expression DEFAULT referencing an AUTO_INCREMENT column",
			stmt: "CREATE TABLE z (id INT AUTO_INCREMENT PRIMARY KEY, x INT DEFAULT (id))",
			code: 3768,
		},
		{
			name:  "ALTER TABLE ADD column with expression DEFAULT referencing a missing column",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT (nosuch + 1)",
			code:  1054,
		},
		{
			name:  "ALTER TABLE ADD column with subquery in expression DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT ((SELECT 1))",
			code:  3769,
		},
		{
			name:  "ALTER TABLE ADD column with variable in expression DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT (@v)",
			code:  3772,
		},
		// 14. Literal DEFAULT validation against column type and range (1067).
		{
			name:  "ALTER TABLE ADD INT with non-numeric string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT 's'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ALTER INT COLUMN SET non-numeric string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ALTER COLUMN c SET DEFAULT 's'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD INT with numeric string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT '123'",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD INT with parenthesized string expression DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x INT DEFAULT ('s')",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD TINYINT with out-of-range DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x TINYINT DEFAULT 128",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD TINYINT UNSIGNED with negative DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x TINYINT UNSIGNED DEFAULT -1",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD DECIMAL(4,2) with overflowing DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x DECIMAL(4,2) DEFAULT 100.00",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD VARCHAR(3) with overlong DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x VARCHAR(3) DEFAULT 'abcd'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD BIT(1) with overflowing byte string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x BIT(1) DEFAULT 's'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD BIT(8) with fitting byte string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x BIT(8) DEFAULT 's'",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with non-member DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('a', 'b') DEFAULT 'c'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD YEAR with invalid string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x YEAR DEFAULT 's'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD DATE with invalid string DEFAULT",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x DATE DEFAULT 's'",
			code:  1067,
		},
		{
			name:  "ALTER TABLE ADD DATE with zero date DEFAULT under default sql_mode",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x DATE DEFAULT '0000-00-00'",
			code:  1067,
		},
		{
			name:    "ALTER TABLE ADD DATE with zero date DEFAULT under empty sql_mode",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "ALTER TABLE a ADD x DATE DEFAULT '0000-00-00'",
			code:    0,
		},
		// Partitioned InnoDB tables in shared tablespaces (1478 vs 3510),
		// repeated DISCARD TABLESPACE (0), accent-insensitive vs accent-sensitive
		// ENUM collations (1291 vs 0), ALLOW_INVALID_DATES (0 vs 1067), and
		// expression-default column drop/rename dependencies (3766).
		{
			name: "CREATE TABLE partitioned in a shared general tablespace",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
			},
			stmt:    "CREATE TABLE z (b INT) TABLESPACE ts_oracle PARTITION BY HASH (b) PARTITIONS 2",
			code:    1478,
			cleanup: []string{"DROP TABLESPACE ts_oracle"},
		},
		{
			name: "CREATE TABLE with a partition in a shared general tablespace",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
			},
			stmt:    "CREATE TABLE z (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE ts_oracle)",
			code:    1478,
			cleanup: []string{"DROP TABLESPACE ts_oracle"},
		},
		{
			name: "CREATE TABLE with a partition naming a missing tablespace",
			stmt: "CREATE TABLE z (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE nosuch_ts)",
			code: 3510,
		},
		{
			name: "CREATE TABLE with a partition in innodb_file_per_table",
			stmt: "CREATE TABLE z (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE innodb_file_per_table)",
			code: 0,
		},
		{
			name: "ALTER TABLE partitioning a table in a shared general tablespace",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
				"CREATE TABLE a (b INT) TABLESPACE ts_oracle",
			},
			stmt:    "ALTER TABLE a PARTITION BY HASH (b) PARTITIONS 2",
			code:    1478,
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_oracle"},
		},
		{
			name: "ALTER TABLE moving a partitioned table into a shared general tablespace",
			setup: []string{
				"CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB",
				"CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))",
			},
			stmt:    "ALTER TABLE a TABLESPACE ts_oracle",
			code:    1478,
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_oracle"},
		},
		{
			name:  "DISCARD TABLESPACE twice succeeds",
			setup: []string{oneTable, "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a DISCARD TABLESPACE",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with accent-colliding values under utf8mb4_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('a', 'á')",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with accent-distinct values under utf8mb4_0900_as_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('a', 'á') COLLATE utf8mb4_0900_as_ci",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with case-colliding values under utf8mb4_0900_as_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('a', 'A') COLLATE utf8mb4_0900_as_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with sharp-s and ss under utf8mb4_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ß', 'ss') COLLATE utf8mb4_0900_ai_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with sharp-s and ss under utf8mb4_general_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ß', 'ss') COLLATE utf8mb4_general_ci",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with sharp-s and s under utf8mb4_general_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ß', 's') COLLATE utf8mb4_general_ci",
			code:  1291,
		},
		{
			name:    "ALTER TABLE ADD DATE with Feb 31 DEFAULT under ALLOW_INVALID_DATES",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'"},
			stmt:    "ALTER TABLE a ADD x DATE DEFAULT '2024-02-31'",
			code:    0,
		},
		{
			name:    "ALTER TABLE ADD DATETIME with Feb 31 DEFAULT under ALLOW_INVALID_DATES",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'"},
			stmt:    "ALTER TABLE a ADD x DATETIME DEFAULT '2024-02-31 12:00:00'",
			code:    0,
		},
		{
			name:    "ALTER TABLE ADD TIMESTAMP with Feb 31 DEFAULT under ALLOW_INVALID_DATES",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'"},
			stmt:    "ALTER TABLE a ADD x TIMESTAMP DEFAULT '2024-02-31 12:00:00'",
			code:    1067,
		},
		{
			name:    "ALTER TABLE ADD DATE with day 32 DEFAULT under ALLOW_INVALID_DATES",
			setup:   []string{oneTable},
			session: []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'"},
			stmt:    "ALTER TABLE a ADD x DATE DEFAULT '2024-02-32'",
			code:    1067,
		},
		{
			name:  "DROP COLUMN referenced by an expression DEFAULT",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))"},
			stmt:  "ALTER TABLE a DROP COLUMN c",
			code:  3766,
		},
		{
			name:  "RENAME COLUMN referenced by an expression DEFAULT",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))"},
			stmt:  "ALTER TABLE a RENAME COLUMN c TO c2",
			code:  3766,
		},
		{
			name:  "CHANGE COLUMN renaming a column referenced by an expression DEFAULT",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))"},
			stmt:  "ALTER TABLE a CHANGE COLUMN c c2 INT",
			code:  3766,
		},
		{
			name:  "DROP and re-ADD a column referenced by an expression DEFAULT",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))"},
			stmt:  "ALTER TABLE a DROP COLUMN c, ADD COLUMN c INT",
			code:  3766,
		},
		{
			name:  "RENAME COLUMN while updating the dependent expression DEFAULT",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))"},
			stmt:  "ALTER TABLE a RENAME COLUMN c TO c2, ALTER COLUMN x SET DEFAULT (c2 + 1)",
			code:  0,
		},
		{
			name:  "DROP COLUMN together with the dependent expression-DEFAULT column",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))"},
			stmt:  "ALTER TABLE a DROP COLUMN c, DROP COLUMN x",
			code:  0,
		},
		// Round-2 review oracle cases: 64-bit SET masks, general tablespace engines,
		// coerced numeric string defaults, and non-InnoDB TABLESPACE occupancy.
		{
			name: "CREATE TABLE with 64-member SET and MaxUint64 numeric DEFAULT",
			stmt: "CREATE TABLE s (c SET('18446744073709551615','m1','m2','m3','m4','m5','m6','m7','m8','m9','m10','m11','m12','m13','m14','m15','m16','m17','m18','m19','m20','m21','m22','m23','m24','m25','m26','m27','m28','m29','m30','m31','m32','m33','m34','m35','m36','m37','m38','m39','m40','m41','m42','m43','m44','m45','m46','m47','m48','m49','m50','m51','m52','m53','m54','m55','m56','m57','m58','m59','m60','m61','m62','m63') DEFAULT 18446744073709551615)",
			code: 0,
		},
		{
			name: "CREATE TABLE with 64-member SET and 2^64 decimal-lexed DEFAULT matching member",
			stmt: "CREATE TABLE s (c SET('18446744073709551616','m1','m2','m3','m4','m5','m6','m7','m8','m9','m10','m11','m12','m13','m14','m15','m16','m17','m18','m19','m20','m21','m22','m23','m24','m25','m26','m27','m28','m29','m30','m31','m32','m33','m34','m35','m36','m37','m38','m39','m40','m41','m42','m43','m44','m45','m46','m47','m48','m49','m50','m51','m52','m53','m54','m55','m56','m57','m58','m59','m60','m61','m62','m63') DEFAULT 18446744073709551616)",
			code: 0,
		},
		{
			name: "CREATE TABLE with 64-member SET and MaxUint64 DEFAULT not in member names",
			stmt: "CREATE TABLE s (c SET('m0','m1','m2','m3','m4','m5','m6','m7','m8','m9','m10','m11','m12','m13','m14','m15','m16','m17','m18','m19','m20','m21','m22','m23','m24','m25','m26','m27','m28','m29','m30','m31','m32','m33','m34','m35','m36','m37','m38','m39','m40','m41','m42','m43','m44','m45','m46','m47','m48','m49','m50','m51','m52','m53','m54','m55','m56','m57','m58','m59','m60','m61','m62','m63') DEFAULT 18446744073709551615)",
			code: 1067,
		},
		{
			name: "CREATE TABLE with SET bitmask DEFAULT matching a member name",
			stmt: "CREATE TABLE a (c SET('a','3') DEFAULT 3)",
			code: 0,
		},
		{
			name: "CREATE TABLE with SET bitmask DEFAULT exceeding member mask",
			stmt: "CREATE TABLE a (c SET('a','4') DEFAULT 4)",
			code: 1067,
		},
		{
			name: "CREATE TABLESPACE with ENGINE=MyISAM",
			stmt: "CREATE TABLESPACE ts_myisam ADD DATAFILE 'ts_myisam.ibd' ENGINE=MyISAM",
			code: 1478,
		},
		{
			name: "CREATE TABLESPACE with ENGINE=MEMORY",
			stmt: "CREATE TABLESPACE ts_mem ADD DATAFILE 'ts_mem.ibd' ENGINE=MEMORY",
			code: 1478,
		},
		{
			name: "CREATE TABLESPACE with unknown ENGINE",
			stmt: "CREATE TABLESPACE ts_bogus ADD DATAFILE 'ts_bogus.ibd' ENGINE=Bogus",
			code: 1286,
		},
		{
			name:    "CREATE TABLESPACE with unknown ENGINE substituted under empty sql_mode",
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "CREATE TABLESPACE ts_sub ADD DATAFILE 'ts_sub.ibd' ENGINE=Bogus",
			code:    0,
			cleanup: []string{"DROP TABLESPACE ts_sub"},
		},
		{
			name:    "ALTER TABLESPACE with mismatched ENGINE=MyISAM",
			setup:   []string{"CREATE TABLESPACE ts_eng ADD DATAFILE 'ts_eng.ibd' ENGINE=InnoDB"},
			stmt:    "ALTER TABLESPACE ts_eng ENGINE=MyISAM",
			code:    3642,
			cleanup: []string{"DROP TABLESPACE ts_eng"},
		},
		{
			name:    "DROP TABLESPACE with mismatched ENGINE=MyISAM",
			setup:   []string{"CREATE TABLESPACE ts_eng ADD DATAFILE 'ts_eng.ibd' ENGINE=InnoDB"},
			stmt:    "DROP TABLESPACE ts_eng ENGINE=MyISAM",
			code:    3642,
			cleanup: []string{"DROP TABLESPACE ts_eng"},
		},
		{
			name: "DROP missing TABLESPACE with ENGINE=MyISAM",
			stmt: "DROP TABLESPACE ts_missing ENGINE=MyISAM",
			code: 3510,
		},
		{
			name:     "CREATE TABLE CHAR(1) with leading-zero numeric DEFAULT 01",
			stmt:     "CREATE TABLE a (c CHAR(1) DEFAULT 01)",
			code:     0,
			contains: []string{"DEFAULT '1'"},
		},
		{
			name:     "CREATE TABLE CHAR(1) with signed zero DEFAULT -0",
			stmt:     "CREATE TABLE a (c CHAR(1) DEFAULT -0)",
			code:     0,
			contains: []string{"DEFAULT '0'"},
		},
		{
			name:     "CREATE TABLE CHAR(1) with scientific DEFAULT 1e0",
			stmt:     "CREATE TABLE a (c CHAR(1) DEFAULT 1e0)",
			code:     0,
			contains: []string{"DEFAULT '1'"},
		},
		{
			name: "CREATE TABLE CHAR(3) with trailing-zero decimal DEFAULT 1.50",
			stmt: "CREATE TABLE a (c CHAR(3) DEFAULT 1.50)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(2) with leading-dot decimal DEFAULT .5",
			stmt: "CREATE TABLE a (c CHAR(2) DEFAULT .5)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE CHAR(3) with leading-dot decimal DEFAULT .5",
			stmt:     "CREATE TABLE a (c CHAR(3) DEFAULT .5)",
			code:     0,
			contains: []string{"DEFAULT '0.5'"},
		},
		{
			name: "CREATE TABLE CHAR(1) with scientific DEFAULT 1e-2",
			stmt: "CREATE TABLE a (c CHAR(1) DEFAULT 1e-2)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE CHAR(2) with scientific DEFAULT 1e-2 rounded by gcvt",
			stmt:     "CREATE TABLE a (c CHAR(2) DEFAULT 1e-2)",
			code:     0,
			contains: []string{"DEFAULT '0'"},
		},
		{
			name: "CREATE TABLE CHAR(4) with inexact scientific DEFAULT 1.23e10",
			stmt: "CREATE TABLE a (c CHAR(4) DEFAULT 1.23e10)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE CHAR(5) with scientific DEFAULT 1.23e10",
			stmt:     "CREATE TABLE a (c CHAR(5) DEFAULT 1.23e10)",
			code:     0,
			contains: []string{"DEFAULT '1e10'"},
		},
		{
			name:     "CREATE TABLE BINARY(1) with leading-zero numeric DEFAULT 01",
			stmt:     "CREATE TABLE a (c BINARY(1) DEFAULT 01)",
			code:     0,
			contains: []string{"DEFAULT '1'"},
		},
		{
			name: "CREATE TABLE BINARY(3) with trailing-zero decimal DEFAULT 1.50",
			stmt: "CREATE TABLE a (c BINARY(3) DEFAULT 1.50)",
			code: 1067,
		},
		{
			name: "ALTER TABLE to MyISAM retains syntactic TABLESPACE in SHOW CREATE TABLE",
			setup: []string{
				"CREATE TABLESPACE ts_m ADD DATAFILE 'ts_m.ibd'",
				oneTable,
			},
			stmt:     "ALTER TABLE a ENGINE=MyISAM, TABLESPACE ts_m",
			code:     0,
			contains: []string{"TABLESPACE `ts_m`", "ENGINE=MyISAM"},
			cleanup:  []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_m"},
		},
		{
			name: "DROP TABLESPACE succeeds when only a MyISAM table names it",
			setup: []string{
				"CREATE TABLESPACE ts_m ADD DATAFILE 'ts_m.ibd'",
				"CREATE TABLE a (b INT) ENGINE=InnoDB TABLESPACE ts_m",
				"ALTER TABLE a ENGINE=MyISAM",
			},
			stmt: "DROP TABLESPACE ts_m",
			code: 0,
		},
		{
			name: "ALTER TABLE converting MyISAM back to InnoDB validates retained missing TABLESPACE",
			setup: []string{
				"CREATE TABLESPACE ts_m ADD DATAFILE 'ts_m.ibd'",
				"CREATE TABLE a (b INT) ENGINE=MyISAM TABLESPACE ts_m",
				"DROP TABLESPACE ts_m",
			},
			stmt: "ALTER TABLE a ENGINE=InnoDB",
			code: 3510,
		},
		// Round-3 review oracle cases: datetime timezone offsets, trailing spaces on
		// CHAR/VARCHAR defaults, hex/bit defaults on CHAR/VARCHAR, and non-strict
		// TINYINT range behavior.
		{
			name: "CREATE TABLE DATETIME with +14:00 timezone offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+14:00')",
			code: 0,
		},
		{
			name: "CREATE TABLE DATETIME with +14:01 out-of-range timezone offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+14:01')",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATETIME with +99:00 invalid timezone offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+99:00')",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATETIME with -13:59 timezone offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 12:00:00-13:59')",
			code: 0,
		},
		{
			name: "CREATE TABLE DATETIME with -14:00 timezone offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 12:00:00-14:00')",
			code: 0,
		},
		{
			name: "CREATE TABLE DATETIME with -14:01 out-of-range timezone offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 12:00:00-14:01')",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATETIME with -00:00 disallowed negative zero offset",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00-00:00')",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATETIME with trailing Z suffix rejected",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00Z')",
			code: 1067,
		},
		{
			name: "CREATE TABLE TIMESTAMP underflowing 1970-01-01 UTC after timezone offset",
			stmt: "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 02:00:00+05:00')",
			code: 1067,
		},
		{
			name: "CREATE TABLE TIMESTAMP valid in UTC after negative timezone offset",
			stmt: "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01-05:00')",
			code: 0,
		},
		{
			name:     "CREATE TABLE CHAR(2) with excess trailing spaces in DEFAULT",
			stmt:     "CREATE TABLE a (c CHAR(2) DEFAULT 'ab   ')",
			code:     0,
			contains: []string{"DEFAULT 'ab'"},
		},
		{
			name:     "CREATE TABLE VARCHAR(2) truncating only excess trailing spaces in DEFAULT",
			stmt:     "CREATE TABLE a (c VARCHAR(2) DEFAULT 'a   ')",
			code:     0,
			contains: []string{"DEFAULT 'a '"},
		},
		{
			name: "CREATE TABLE BINARY(2) rejects excess trailing spaces in DEFAULT",
			stmt: "CREATE TABLE a (c BINARY(2) DEFAULT 'ab   ')",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) with 1-char hex literal DEFAULT 0x41",
			stmt: "CREATE TABLE a (c CHAR(1) DEFAULT 0x41)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) rejects 2-char hex literal DEFAULT 0x4142",
			stmt: "CREATE TABLE a (c CHAR(1) DEFAULT 0x4142)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) accepts 3-byte UTF-8 1-rune hex literal DEFAULT",
			stmt: "CREATE TABLE a (c CHAR(1) DEFAULT 0xE4B8AD)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf8mb4 rejects invalid UTF-8 hex literal DEFAULT 0xFF",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 0xFF)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) latin1 accepts single-byte hex literal DEFAULT 0xFF",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET latin1 DEFAULT 0xFF)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) rejects 2-byte bit literal DEFAULT",
			stmt: "CREATE TABLE a (c CHAR(1) DEFAULT b'0100000101000010')",
			code: 1067,
		},
		{
			name:    "CREATE TABLE TINYINT out-of-range DEFAULT 128 rejected even under empty sql_mode",
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "CREATE TABLE a (c TINYINT DEFAULT 128)",
			code:    1067,
		},
		{
			name:    "CREATE TABLE DECIMAL non-numeric string DEFAULT coerced to 0 under empty sql_mode",
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "CREATE TABLE a (c DECIMAL(3,1) DEFAULT 'abc')",
			code:    0,
		},
		// Round-4 review oracle cases: full UCA 9.0.0 weights (ligatures/digraphs),
		// undo tablespace lifecycle, FLOAT(M,D) precision/scale, and UNIQUE USING RTREE.
		{
			name:  "ALTER TABLE ADD ENUM with ff and U+FB00 ligature under utf8mb4_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_ai_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with ff and U+FB00 ligature under utf8mb4_0900_as_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_as_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with ff and U+FB00 ligature distinct under utf8mb4_0900_as_cs",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_as_cs",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with ff and U+FB00 ligature distinct under utf8mb4_general_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_general_ci",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with dz and U+02A3 digraph under utf8mb4_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('dz', 'ʣ') COLLATE utf8mb4_0900_ai_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with l-middot and U+0140 distinct under utf8mb4_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('l·', 'ŀ') COLLATE utf8mb4_0900_ai_ci",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with SMP bold A and ASCII a under utf8mb4_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('a', '𝐀') COLLATE utf8mb4_0900_ai_ci",
			code:  1291,
		},
		{
			name:    "CREATE duplicate UNDO TABLESPACE rejected with 1813",
			setup:   []string{"CREATE UNDO TABLESPACE u_dup ADD DATAFILE 'u_dup.ibu'"},
			stmt:    "CREATE UNDO TABLESPACE u_dup ADD DATAFILE 'u_dup2.ibu'",
			code:    1813,
			cleanup: []string{"ALTER UNDO TABLESPACE u_dup SET INACTIVE", "SELECT SLEEP(1)", "DROP UNDO TABLESPACE u_dup"},
		},
		{
			name: "CREATE UNDO TABLESPACE with non-.ibu datafile rejected with 3121",
			stmt: "CREATE UNDO TABLESPACE u_ext ADD DATAFILE 'u_ext.ibd'",
			code: 3121,
		},
		{
			name: "ALTER missing UNDO TABLESPACE rejected with 3510",
			stmt: "ALTER UNDO TABLESPACE u_nosuch SET INACTIVE",
			code: 3510,
		},
		{
			name: "DROP missing UNDO TABLESPACE rejected with 3510",
			stmt: "DROP UNDO TABLESPACE u_nosuch",
			code: 3510,
		},
		{
			name:    "DROP active UNDO TABLESPACE rejected with 1529",
			setup:   []string{"CREATE UNDO TABLESPACE u_act ADD DATAFILE 'u_act.ibu'"},
			stmt:    "DROP UNDO TABLESPACE u_act",
			code:    1529,
			cleanup: []string{"ALTER UNDO TABLESPACE u_act SET INACTIVE", "SELECT SLEEP(1)", "DROP UNDO TABLESPACE u_act"},
		},
		{
			name: "DROP inactive UNDO TABLESPACE succeeds",
			setup: []string{
				"CREATE UNDO TABLESPACE u_inact ADD DATAFILE 'u_inact.ibu'",
				"ALTER UNDO TABLESPACE u_inact SET INACTIVE",
				"SELECT SLEEP(1)",
			},
			stmt: "DROP UNDO TABLESPACE u_inact",
			code: 0,
		},
		{
			name:    "DROP TABLESPACE on an undo tablespace rejected with 3119",
			setup:   []string{"CREATE UNDO TABLESPACE u_kind ADD DATAFILE 'u_kind.ibu'"},
			stmt:    "DROP TABLESPACE u_kind",
			code:    3119,
			cleanup: []string{"ALTER UNDO TABLESPACE u_kind SET INACTIVE", "SELECT SLEEP(1)", "DROP UNDO TABLESPACE u_kind"},
		},
		{
			name: "CREATE TABLE FLOAT(4,2) out-of-range DEFAULT 100",
			stmt: "CREATE TABLE a (c FLOAT(4,2) DEFAULT 100)",
			code: 1067,
		},
		{
			name: "CREATE TABLE FLOAT(4,2) in-range DEFAULT 99.994",
			stmt: "CREATE TABLE a (c FLOAT(4,2) DEFAULT 99.994)",
			code: 0,
		},
		{
			name: "CREATE TABLE FLOAT(4,2) rounding-overflow DEFAULT 99.995",
			stmt: "CREATE TABLE a (c FLOAT(4,2) DEFAULT 99.995)",
			code: 1067,
		},
		{
			name: "CREATE TABLE DOUBLE(4,2) out-of-range DEFAULT 100",
			stmt: "CREATE TABLE a (c DOUBLE(4,2) DEFAULT 100)",
			code: 1067,
		},
		{
			name:  "CREATE UNIQUE INDEX USING RTREE on geometry column rejected with 3728",
			setup: []string{"CREATE TABLE s (pt POINT NOT NULL)"},
			stmt:  "CREATE UNIQUE INDEX u USING RTREE ON s(pt)",
			code:  3728,
		},
		{
			name:  "CREATE UNIQUE INDEX USING RTREE on non-geometry column rejected with 1687",
			setup: []string{oneTable},
			stmt:  "CREATE UNIQUE INDEX u USING RTREE ON a(b)",
			code:  1687,
		},
		{
			name: "CREATE TABLE UNIQUE KEY USING RTREE rejected with 3728",
			stmt: "CREATE TABLE s (pt POINT NOT NULL, UNIQUE KEY u (pt) USING RTREE)",
			code: 3728,
		},
		{
			name:  "ALTER TABLE ADD UNIQUE KEY USING RTREE rejected with 3728",
			setup: []string{"CREATE TABLE s (pt POINT NOT NULL)"},
			stmt:  "ALTER TABLE s ADD UNIQUE KEY u (pt) USING RTREE",
			code:  3728,
		},
		// Round-5 review oracle cases: session time_zone TIMESTAMP bounds, locale-
		// tailored collations, and stored coerced numeric defaults.
		{
			name:    "CREATE TABLE TIMESTAMP overflowing 2038 UTC under negative session time_zone",
			session: []string{"SET SESSION time_zone = '-13:59'"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07')",
			code:    1067,
		},
		{
			name:    "CREATE TABLE TIMESTAMP underflowing 1970 UTC under positive session time_zone",
			session: []string{"SET SESSION time_zone = '+14:00'"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')",
			code:    1067,
		},
		{
			name:    "CREATE TABLE TIMESTAMP valid in UTC under positive session time_zone",
			session: []string{"SET SESSION time_zone = '+14:00'"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 14:00:01')",
			code:    0,
		},
		{
			name: "SET time_zone -14:00 rejected with 1298",
			stmt: "SET SESSION time_zone = '-14:00'",
			code: 1298,
		},
		{
			name:  "ALTER TABLE ADD ENUM with ae and a-umlaut under utf8mb4_de_pb_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ae', 'ä') COLLATE utf8mb4_de_pb_0900_ai_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with ae and decomposed a-umlaut distinct under utf8mb4_de_pb_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ae', 'a\u0308') COLLATE utf8mb4_de_pb_0900_ai_ci",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD ENUM with ae and a-umlaut under utf8mb4_german2_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('ae', 'ä') COLLATE utf8mb4_german2_ci",
			code:  1291,
		},
		{
			name:  "ALTER TABLE ADD ENUM with aa and a-ring under utf8mb4_da_0900_ai_ci",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD x ENUM('aa', 'å') COLLATE utf8mb4_da_0900_ai_ci",
			code:  1291,
		},
		{
			name:     "CREATE TABLE INT DEFAULT 1.5 stores rounded DEFAULT '2'",
			stmt:     "CREATE TABLE a (c INT DEFAULT 1.5)",
			code:     0,
			contains: []string{"DEFAULT '2'"},
		},
		{
			name:     "CREATE TABLE DECIMAL(5,2) DEFAULT 1.235 stores rounded DEFAULT '1.24'",
			stmt:     "CREATE TABLE a (c DECIMAL(5,2) DEFAULT 1.235)",
			code:     0,
			contains: []string{"DEFAULT '1.24'"},
		},
		{
			name:     "CREATE TABLE FLOAT(4,2) DEFAULT 1.5 stores scaled DEFAULT '1.50'",
			stmt:     "CREATE TABLE a (c FLOAT(4,2) DEFAULT 1.5)",
			code:     0,
			contains: []string{"DEFAULT '1.50'"},
		},
		// Round-6 review oracle cases: BINARY/VARBINARY numeric default coercion
		// and partition/subpartition TABLESPACE rendering in SHOW CREATE TABLE.
		{
			name:     "CREATE TABLE VARBINARY(2) with leading-zero numeric DEFAULT 007",
			stmt:     "CREATE TABLE a (c VARBINARY(2) DEFAULT 007)",
			code:     0,
			contains: []string{"DEFAULT '7'"},
		},
		{
			name:     "CREATE TABLE VARBINARY(4) with scientific numeric DEFAULT 1e-2",
			stmt:     "CREATE TABLE a (c VARBINARY(4) DEFAULT 1e-2)",
			code:     0,
			contains: []string{"DEFAULT '0.01'"},
		},
		{
			name:     "CREATE TABLE BINARY(3) with negative-zero decimal DEFAULT -0.0",
			stmt:     "CREATE TABLE a (c BINARY(3) DEFAULT -0.0)",
			code:     0,
			contains: []string{"DEFAULT '0.0'"},
		},
		{
			name:     "CREATE TABLE with partition TABLESPACE renders TABLESPACE in SHOW CREATE TABLE",
			stmt:     "CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = innodb_file_per_table)",
			code:     0,
			contains: []string{"PARTITION p0 VALUES LESS THAN (10) TABLESPACE = `innodb_file_per_table` ENGINE = InnoDB"},
		},
		{
			name:     "CREATE TABLE with subpartition TABLESPACE renders TABLESPACE in SHOW CREATE TABLE",
			stmt:     "CREATE TABLE a (b INT) PARTITION BY RANGE (b) SUBPARTITION BY HASH (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = innodb_file_per_table (SUBPARTITION sp0, SUBPARTITION sp1))",
			code:     0,
			contains: []string{"SUBPARTITION sp0 TABLESPACE = `innodb_file_per_table` ENGINE = InnoDB", "SUBPARTITION sp1 TABLESPACE = `innodb_file_per_table` ENGINE = InnoDB"},
		},
		// Round-7 review oracle cases: user-variable time_zone assignments, padded
		// numeric temporal defaults and rounding, and UTF-16/Unicode binary defaults.
		{
			name:    "SET time_zone via user variable enforces TIMESTAMP UTC lower bound",
			session: []string{"SET @tz = '+14:00'", "SET time_zone = @tz"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')",
			code:    1067,
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name:    "SET time_zone via quoted user variable enforces TIMESTAMP UTC upper bound",
			session: []string{"SET @'tz_q' = '-13:59'", "SET time_zone = @tz_q"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07')",
			code:    1067,
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name:    "SET time_zone via CONCAT user variable accepts valid offset",
			session: []string{"SET @tz_cat = CONCAT('+', '14:00')", "SET time_zone = @tz_cat"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 14:00:01')",
			code:    0,
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name: "SET time_zone to unset user variable rejected with 1231",
			stmt: "SET time_zone = @omni_unset_tz",
			code: 1231,
		},
		{
			name:    "SET time_zone to NULL user variable rejected with 1231",
			session: []string{"SET @tz_null = NULL"},
			stmt:    "SET time_zone = @tz_null",
			code:    1231,
		},
		{
			name:    "SET time_zone to numeric user variable rejected with 1232",
			session: []string{"SET @tz_num = 1"},
			stmt:    "SET time_zone = @tz_num",
			code:    1232,
		},
		{
			name:    "SET time_zone to invalid user variable rejected with 1298",
			session: []string{"SET @tz_bad = 'not_a_tz'"},
			stmt:    "SET time_zone = @tz_bad",
			code:    1298,
		},
		{
			name:     "CREATE TABLE DATE DEFAULT 10101 pads to 2001-01-01",
			stmt:     "CREATE TABLE a (c DATE DEFAULT 10101)",
			code:     0,
			contains: []string{"DEFAULT '2001-01-01'"},
		},
		{
			name:     "CREATE TABLE DATE DEFAULT 101 pads to 2000-01-01",
			stmt:     "CREATE TABLE a (c DATE DEFAULT 101)",
			code:     0,
			contains: []string{"DEFAULT '2000-01-01'"},
		},
		{
			name: "CREATE TABLE DATE DEFAULT 100 rejected with 1067",
			stmt: "CREATE TABLE a (c DATE DEFAULT 100)",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATE DEFAULT '10101' parses as 2010-10-01",
			stmt: "CREATE TABLE a (c DATE DEFAULT '10101')",
			code: 0,
		},
		{
			name: "CREATE TABLE DATE DEFAULT 13101 pads to invalid month 31",
			stmt: "CREATE TABLE a (c DATE DEFAULT 13101)",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATE DEFAULT '13101' parses as valid 2013-10-01",
			stmt: "CREATE TABLE a (c DATE DEFAULT '13101')",
			code: 0,
		},
		{
			name: "CREATE TABLE DATE DEFAULT '101' rejected with 1067",
			stmt: "CREATE TABLE a (c DATE DEFAULT '101')",
			code: 1067,
		},
		{
			name:     "CREATE TABLE DATE DEFAULT 1010101 pads 7 digits to 0101-01-01",
			stmt:     "CREATE TABLE a (c DATE DEFAULT 1010101)",
			code:     0,
			contains: []string{"DEFAULT '0101-01-01'"},
		},
		{
			name: "CREATE TABLE TIMESTAMP DEFAULT 1010101 rejected before 1970",
			stmt: "CREATE TABLE a (c TIMESTAMP DEFAULT 1010101)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE TIMESTAMP DEFAULT 00010101 evaluates as 10101 (2001-01-01)",
			stmt:     "CREATE TABLE a (c TIMESTAMP DEFAULT 00010101)",
			code:     0,
			contains: []string{"DEFAULT '2001-01-01 00:00:00'"},
		},
		{
			name:     "CREATE TABLE DATETIME(3) DEFAULT 101010101.5 pads 9 digits with fraction",
			stmt:     "CREATE TABLE a (c DATETIME(3) DEFAULT 101010101.5)",
			code:     0,
			contains: []string{"DEFAULT '2000-01-01 01:01:01.500'"},
		},
		{
			name: "CREATE TABLE DATETIME DEFAULT 20200101.5 rejects fraction on date-only integer",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT 20200101.5)",
			code: 1067,
		},
		{
			name: "CREATE TABLE DATETIME DEFAULT '9999-12-31 23:59:59.6' overflows on rounding",
			stmt: "CREATE TABLE a (c DATETIME DEFAULT '9999-12-31 23:59:59.6')",
			code: 1067,
		},
		{
			name: "CREATE TABLE TIMESTAMP DEFAULT '2038-01-19 03:14:07.6' overflows on rounding",
			stmt: "CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07.6')",
			code: 1067,
		},
		{
			name:     "CREATE TABLE TIMESTAMP(1) DEFAULT '2038-01-19 03:14:07.6' accepted within max UTC",
			stmt:     "CREATE TABLE a (c TIMESTAMP(1) DEFAULT '2038-01-19 03:14:07.6')",
			code:     0,
			contains: []string{"DEFAULT '2038-01-19 03:14:07.6'"},
		},
		{
			name: "CREATE TABLE CHAR(1) utf16 DEFAULT 0xFFFD accepts replacement character",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xFFFD)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf16le DEFAULT 0xFDFF accepts replacement character",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16le DEFAULT 0xFDFF)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf16 DEFAULT 0xD800 rejects unpaired high surrogate",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xD800)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) utf16 DEFAULT 0xDC00 rejects unpaired low surrogate",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xDC00)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) utf16 DEFAULT 0xD83DDE00 accepts supplementary surrogate pair",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xD83DDE00)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf16 DEFAULT 0x41 left-pads single byte",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0x41)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf16 DEFAULT 0x004142 left-pads 3 bytes to 2 chars and rejects",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0x004142)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(2) utf16 DEFAULT 0x004142 left-pads 3 bytes to 2 chars and accepts",
			stmt: "CREATE TABLE a (c CHAR(2) CHARACTER SET utf16 DEFAULT 0x004142)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf32 DEFAULT 0x41 left-pads single byte",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf32 DEFAULT 0x41)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) utf8mb3 DEFAULT supplementary emoji rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb3 DEFAULT '😀')",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) ascii DEFAULT non-ASCII character rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET ascii DEFAULT 'é')",
			code: 1067,
		},
		// Round-8 review oracle cases: TIME_TRUNCATE_FRACTIONAL temporal defaults
		// and bit/hex defaults across TIME, DATE, DATETIME, TIMESTAMP, YEAR, ENUM, SET.
		{
			name:     "CREATE TABLE DATETIME DEFAULT '9999-12-31 23:59:59.9' truncates under TIME_TRUNCATE_FRACTIONAL",
			session:  []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'"},
			stmt:     "CREATE TABLE a (c DATETIME DEFAULT '9999-12-31 23:59:59.9')",
			code:     0,
			contains: []string{"DEFAULT '9999-12-31 23:59:59'"},
		},
		{
			name:     "CREATE TABLE TIMESTAMP DEFAULT '2038-01-19 03:14:07.9' truncates under TIME_TRUNCATE_FRACTIONAL",
			session:  []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'"},
			stmt:     "CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07.9')",
			code:     0,
			contains: []string{"DEFAULT '2038-01-19 03:14:07'"},
		},
		{
			name:    "CREATE TABLE TIMESTAMP DEFAULT '1970-01-01 00:00:00.9' truncates below min UTC under TIME_TRUNCATE_FRACTIONAL",
			session: []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:00.9')",
			code:    1067,
		},
		{
			name:     "CREATE TABLE TIMESTAMP DEFAULT '1970-01-01 00:00:00.9' rounds up to min UTC by default",
			stmt:     "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:00.9')",
			code:     0,
			contains: []string{"DEFAULT '1970-01-01 00:00:01'"},
		},
		{
			name:    "CREATE TABLE TIME DEFAULT '838:59:59.9' still exceeds input max under TIME_TRUNCATE_FRACTIONAL",
			session: []string{"SET SESSION sql_mode = 'TIME_TRUNCATE_FRACTIONAL'"},
			stmt:    "CREATE TABLE a (c TIME DEFAULT '838:59:59.9')",
			code:    1067,
		},
		{
			name:     "CREATE TABLE TIME DEFAULT '838:59:58.9' truncates to 838:59:58 under TIME_TRUNCATE_FRACTIONAL",
			session:  []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'"},
			stmt:     "CREATE TABLE a (c TIME DEFAULT '838:59:58.9')",
			code:     0,
			contains: []string{"DEFAULT '838:59:58'"},
		},
		{
			name:     "CREATE TABLE TIME DEFAULT '838:59:58.9' rounds to 838:59:59 by default",
			stmt:     "CREATE TABLE a (c TIME DEFAULT '838:59:58.9')",
			code:     0,
			contains: []string{"DEFAULT '838:59:59'"},
		},
		{
			name:     "CREATE TABLE TIME(2) DEFAULT 10203.459 truncates to 01:02:03.45 under TIME_TRUNCATE_FRACTIONAL",
			session:  []string{"SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'"},
			stmt:     "CREATE TABLE a (c TIME(2) DEFAULT 10203.459)",
			code:     0,
			contains: []string{"DEFAULT '01:02:03.45'"},
		},
		{
			name:     "CREATE TABLE TIME DEFAULT bit literal '12:00:00' decodes and coerces",
			stmt:     "CREATE TABLE a (c TIME DEFAULT b'0011000100110010001110100011000000110000001110100011000000110000')",
			code:     0,
			contains: []string{"DEFAULT '12:00:00'"},
		},
		{
			name:     "CREATE TABLE TIME DEFAULT hex literal '12:00:00' decodes and coerces",
			stmt:     "CREATE TABLE a (c TIME DEFAULT 0x31323A30303A3030)",
			code:     0,
			contains: []string{"DEFAULT '12:00:00'"},
		},
		{
			name:     "CREATE TABLE DATE DEFAULT bit literal '2020-01-01' decodes and coerces",
			stmt:     "CREATE TABLE a (c DATE DEFAULT b'00110010001100000011001000110000001011010011000000110001001011010011000000110001')",
			code:     0,
			contains: []string{"DEFAULT '2020-01-01'"},
		},
		{
			name:     "CREATE TABLE YEAR DEFAULT 0x07E4 evaluates as integer 2020",
			stmt:     "CREATE TABLE a (c YEAR DEFAULT 0x07E4)",
			code:     0,
			contains: []string{"DEFAULT '2020'"},
		},
		{
			name: "CREATE TABLE YEAR DEFAULT 0x32303230 rejected out of YEAR integer range",
			stmt: "CREATE TABLE a (c YEAR DEFAULT 0x32303230)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE ENUM DEFAULT b'01100001' matches 'a'",
			stmt:     "CREATE TABLE a (c ENUM('a','b') DEFAULT b'01100001')",
			code:     0,
			contains: []string{"DEFAULT 'a'"},
		},
		{
			name: "CREATE TABLE ENUM DEFAULT 0x7A rejected when not a member",
			stmt: "CREATE TABLE a (c ENUM('a','b') DEFAULT 0x7A)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE SET DEFAULT 0x612C62 matches 'a,b'",
			stmt:     "CREATE TABLE a (c SET('a','b') DEFAULT 0x612C62)",
			code:     0,
			contains: []string{"DEFAULT 'a,b'"},
		},
		{
			name: "CREATE TABLE SET DEFAULT b'1' rejected when not a member",
			stmt: "CREATE TABLE a (c SET('a','b') DEFAULT b'1')",
			code: 1067,
		},
		// 16. Offset-bearing temporal defaults normalized to session time_zone, tablespace
		//     ADD DATAFILE / FILE_BLOCK_SIZE checks, and table AUTOEXTEND_SIZE preservation.
		{
			name:     "CREATE TABLE DATETIME and TIMESTAMP defaults with explicit offsets normalize into session time_zone",
			session:  []string{"SET SESSION time_zone = '+08:00'"},
			stmt:     "CREATE TABLE a (dt DATETIME DEFAULT '2024-01-01 00:00:00+14:00', ts TIMESTAMP DEFAULT '2024-01-01 00:00:00+14:00', d DATE DEFAULT '2024-01-01 00:00:00+14:00')",
			code:     0,
			contains: []string{"`dt` datetime DEFAULT '2023-12-31 18:00:00'", "`ts` timestamp NULL DEFAULT '2023-12-31 18:00:00'", "`d` date DEFAULT '2023-12-31'"},
			omits:    []string{"+14:00"},
		},
		{
			name:     "CREATE TABLE DATETIME(3) default with offset and hex offset normalize into session time_zone",
			session:  []string{"SET SESSION time_zone = '+00:00'"},
			stmt:     "CREATE TABLE a (dt3 DATETIME(3) DEFAULT '2024-01-01 00:00:00.5678+05:30', dth DATETIME DEFAULT 0x323032342D30312D30312030303A30303A30302B31343A3030)",
			code:     0,
			contains: []string{"`dt3` datetime(3) DEFAULT '2023-12-31 18:30:00.568'", "`dth` datetime DEFAULT '2023-12-31 10:00:00'"},
		},
		{
			name:    "CREATE TABLE DATETIME default with offset overflowing year 9999 in session time_zone rejected",
			session: []string{"SET SESSION time_zone = '+05:00'"},
			stmt:    "CREATE TABLE a (dt DATETIME DEFAULT '9999-12-31 23:00:00+00:00')",
			code:    1067,
		},
		{
			name:    "CREATE TABLE DATETIME default with offset underflowing year 0 in session time_zone rejected",
			session: []string{"SET SESSION time_zone = '-05:00'"},
			stmt:    "CREATE TABLE a (dt DATETIME DEFAULT '0000-01-01 01:00:00+00:00')",
			code:    1067,
		},
		{
			name:     "ALTER TABLE ALTER COLUMN SET DEFAULT with explicit offset normalizes into session time_zone",
			setup:    []string{"CREATE TABLE a (dt DATETIME, ts TIMESTAMP NULL)"},
			session:  []string{"SET SESSION time_zone = '-05:00'"},
			stmt:     "ALTER TABLE a ALTER COLUMN dt SET DEFAULT '2024-01-01 02:00:00+00:00'",
			code:     0,
			contains: []string{"`dt` datetime DEFAULT '2023-12-31 21:00:00'"},
		},
		{
			name: "CREATE TABLESPACE ADD DATAFILE without .ibd extension rejected with 3121",
			stmt: "CREATE TABLESPACE ts_bad_ext ADD DATAFILE 'ts_bad_ext.txt' ENGINE=InnoDB",
			code: 3121,
		},
		{
			name: "CREATE TABLESPACE ADD DATAFILE empty filename rejected with 3121",
			stmt: "CREATE TABLESPACE ts_emp_df ADD DATAFILE '' ENGINE=InnoDB",
			code: 3121,
		},
		{
			name: "CREATE TABLESPACE ADD DATAFILE bare .ibd filename rejected with 3121",
			stmt: "CREATE TABLESPACE ts_dot_df ADD DATAFILE '.ibd' ENGINE=InnoDB",
			code: 3121,
		},
		{
			name: "CREATE TABLESPACE ADD DATAFILE uppercase .IBD rejected with 3121",
			stmt: "CREATE TABLESPACE ts_up_df ADD DATAFILE 'ts_up_df.IBD' ENGINE=InnoDB",
			code: 3121,
		},
		{
			name: "CREATE TABLESPACE ADD DATAFILE outside directory rejected with 3121",
			stmt: "CREATE TABLESPACE ts_dir_df ADD DATAFILE '../ts_dir_df.ibd' ENGINE=InnoDB",
			code: 3121,
		},
		{
			name:    "CREATE TABLESPACE duplicate ADD DATAFILE rejected with 3606",
			setup:   []string{"CREATE TABLESPACE ts_df1 ADD DATAFILE 'shared_df.ibd' ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLESPACE ts_df1"},
			stmt:    "CREATE TABLESPACE ts_df2 ADD DATAFILE 'shared_df.ibd' ENGINE=InnoDB",
			code:    3606,
		},
		{
			name: "CREATE TABLESPACE FILE_BLOCK_SIZE=1 rejected with 1478",
			stmt: "CREATE TABLESPACE ts_fbs_bad1 ADD DATAFILE 'ts_fbs_bad1.ibd' FILE_BLOCK_SIZE=1 ENGINE=InnoDB",
			code: 1478,
		},
		{
			name: "CREATE TABLESPACE FILE_BLOCK_SIZE=32K rejected above innodb_page_size with 1478",
			stmt: "CREATE TABLESPACE ts_fbs_32k ADD DATAFILE 'ts_fbs_32k.ibd' FILE_BLOCK_SIZE=32K ENGINE=InnoDB",
			code: 1478,
		},
		{
			name:    "CREATE TABLESPACE FILE_BLOCK_SIZE=8K accepted",
			cleanup: []string{"DROP TABLESPACE ts_fbs_8k"},
			stmt:    "CREATE TABLESPACE ts_fbs_8k ADD DATAFILE 'ts_fbs_8k.ibd' FILE_BLOCK_SIZE=8K ENGINE=InnoDB",
			code:    0,
		},
		{
			name: "CREATE TABLESPACE FILE_BLOCK_SIZE validated ahead of AUTOEXTEND_SIZE and DATAFILE",
			stmt: "CREATE TABLESPACE ts_fbs_prec ADD DATAFILE 'ts_fbs_prec.txt' FILE_BLOCK_SIZE=1 AUTOEXTEND_SIZE=1 ENGINE=InnoDB",
			code: 1478,
		},
		{
			name:     "CREATE TABLE AUTOEXTEND_SIZE=4M preserved as byte value in SHOW CREATE TABLE",
			stmt:     "CREATE TABLE a (c INT) AUTOEXTEND_SIZE=4M",
			code:     0,
			contains: []string{"/*!80023 AUTOEXTEND_SIZE=4194304 */ ENGINE=InnoDB"},
		},
		{
			name:     "ALTER TABLE AUTOEXTEND_SIZE=8M preserved in SHOW CREATE TABLE",
			setup:    []string{"CREATE TABLE a (c INT)"},
			stmt:     "ALTER TABLE a AUTOEXTEND_SIZE=8M",
			code:     0,
			contains: []string{"/*!80023 AUTOEXTEND_SIZE=8388608 */ ENGINE=InnoDB"},
		},
		{
			name:  "ALTER TABLE AUTOEXTEND_SIZE=0 clears stored AUTOEXTEND_SIZE",
			setup: []string{"CREATE TABLE a (c INT) AUTOEXTEND_SIZE=8M"},
			stmt:  "ALTER TABLE a AUTOEXTEND_SIZE=0",
			code:  0,
			omits: []string{"AUTOEXTEND_SIZE"},
		},
		{
			name:     "CREATE TABLE LIKE preserves AUTOEXTEND_SIZE",
			setup:    []string{"CREATE TABLE src_ae (c INT) AUTOEXTEND_SIZE=12M"},
			stmt:     "CREATE TABLE a LIKE src_ae",
			code:     0,
			contains: []string{"/*!80023 AUTOEXTEND_SIZE=12582912 */ ENGINE=InnoDB"},
		},
		{
			name:    "CREATE TABLE AUTOEXTEND_SIZE on general tablespace rejected with 4024",
			setup:   []string{"CREATE TABLESPACE ts_ae_gen ADD DATAFILE 'ts_ae_gen.ibd'"},
			cleanup: []string{"DROP TABLESPACE ts_ae_gen"},
			stmt:    "CREATE TABLE a (c INT) TABLESPACE ts_ae_gen AUTOEXTEND_SIZE=8M",
			code:    4024,
		},
		{
			name: "CREATE TEMPORARY TABLE non-zero AUTOEXTEND_SIZE rejected with 4026",
			stmt: "CREATE TEMPORARY TABLE a_tmp_ae (c INT) AUTOEXTEND_SIZE=8M",
			code: 4026,
		},
		{
			name: "CREATE TABLE AUTOEXTEND_SIZE on innodb_system rejected with 1478",
			stmt: "CREATE TABLE a (c INT) TABLESPACE innodb_system AUTOEXTEND_SIZE=8M",
			code: 1478,
		},
		// 17. Binary floating/decimal/integer defaults enforce 8-byte / signed-int64 limits without panic.
		{
			name: "CREATE TABLE DOUBLE DEFAULT >1024-bit hex rejected with 1067 without panic",
			stmt: "CREATE TABLE a (c DOUBLE DEFAULT 0x" + strings.Repeat("FF", 140) + ")",
			code: 1067,
		},
		{
			name: "CREATE TABLE FLOAT DEFAULT >8-byte hex rejected with 1067",
			stmt: "CREATE TABLE a (c FLOAT DEFAULT 0x000000000000000001)",
			code: 1067,
		},
		{
			name: "CREATE TABLE signed DOUBLE DEFAULT 0x8000000000000000 rejected with 1067",
			stmt: "CREATE TABLE a (c DOUBLE DEFAULT 0x8000000000000000)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE DOUBLE UNSIGNED DEFAULT 0xFFFFFFFFFFFFFFFF accepted and formatted",
			stmt:     "CREATE TABLE a (c DOUBLE UNSIGNED DEFAULT 0xFFFFFFFFFFFFFFFF)",
			code:     0,
			contains: []string{"DEFAULT '1.8446744073709552e19'"},
		},
		{
			name:     "CREATE TABLE FLOAT DEFAULT 0x7FFFFFFFFFFFFFFF accepted and formatted",
			stmt:     "CREATE TABLE a (c FLOAT DEFAULT 0x7FFFFFFFFFFFFFFF)",
			code:     0,
			contains: []string{"DEFAULT '9.22337e18'"},
		},
		{
			name: "CREATE TABLE DECIMAL(30,2) signed DEFAULT 0x8000000000000000 rejected with 1067",
			stmt: "CREATE TABLE a (c DECIMAL(30,2) DEFAULT 0x8000000000000000)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE FLOAT DEFAULT 0x41424344 rounds to 6 significant digits",
			stmt:     "CREATE TABLE a (c FLOAT DEFAULT 0x41424344)",
			code:     0,
			contains: []string{"DEFAULT '1094860000'"},
		},
		{
			name: "CREATE TABLE BIGINT DEFAULT >8-byte leading-zero hex rejected with 1067",
			stmt: "CREATE TABLE a (c BIGINT DEFAULT 0x000000000000000001)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE non-strict DOUBLE and BIGINT clamp oversized hex defaults",
			session:  []string{"SET SESSION sql_mode = ''"},
			stmt:     "CREATE TABLE a (d DOUBLE DEFAULT 0xFFFFFFFFFFFFFFFF, b BIGINT DEFAULT 0x000000000000000001)",
			code:     0,
			contains: []string{"DEFAULT '9.223372036854776e18'", "DEFAULT '9223372036854775807'"},
			cleanup:  []string{"SET SESSION sql_mode = DEFAULT"},
		},
		// 18. Legacy multibyte binary defaults, bit-literal user-variable time_zone,
		// and YEAR / ENUM / SET default canonicalization.
		{
			name: "CREATE TABLE CHAR(1) sjis DEFAULT 0x82A0 accepts single 2-byte character",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT 0x82A0)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) sjis DEFAULT 0x82 rejects truncated lead byte",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT 0x82)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) ujis DEFAULT 0x8FA2AF accepts single 3-byte character",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET ujis DEFAULT 0x8FA2AF)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) ujis DEFAULT 0x8FA2 rejects truncated 3-byte sequence",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET ujis DEFAULT 0x8FA2)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) gb18030 DEFAULT 0x81308130 accepts single 4-byte character",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET gb18030 DEFAULT 0x81308130)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) gb18030 DEFAULT 0x813081 rejects truncated 4-byte sequence",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET gb18030 DEFAULT 0x813081)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) big5 DEFAULT 0xA440 accepts single 2-byte character",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET big5 DEFAULT 0xA440)",
			code: 0,
		},
		{
			name: "CREATE TABLE CHAR(1) euckr DEFAULT 0x815B rejects invalid trail byte",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET euckr DEFAULT 0x815B)",
			code: 1067,
		},
		{
			name:    "SET time_zone via bit-literal user variable enforces TIMESTAMP UTC lower bound",
			session: []string{"SET @tz = b'001010110011000100110100001110100011000000110000'", "SET time_zone = @tz"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')",
			code:    1067,
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name:    "SET time_zone via bit-literal user variable accepts valid offset",
			session: []string{"SET @tz_b = b'001010110011000100110100001110100011000000110000'", "SET time_zone = @tz_b"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 14:00:01')",
			code:    0,
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name: "CREATE TABLE YEAR normalizes string and numeric two-digit defaults",
			stmt: "CREATE TABLE a (y0 YEAR DEFAULT 0, ys0 YEAR DEFAULT '0', ys000 YEAR DEFAULT '000', ys0000 YEAR DEFAULT '0000', y70 YEAR DEFAULT 70, ys01 YEAR DEFAULT '01')",
			code: 0,
			contains: []string{
				"`y0` year DEFAULT '0000'",
				"`ys0` year DEFAULT '2000'",
				"`ys000` year DEFAULT '2000'",
				"`ys0000` year DEFAULT '0000'",
				"`y70` year DEFAULT '1970'",
				"`ys01` year DEFAULT '2001'",
			},
		},
		{
			name: "CREATE TABLE ENUM and SET canonicalize textual and numeric defaults to member spellings",
			stmt: "CREATE TABLE a (e1 ENUM('a','b') DEFAULT 'A', e2 ENUM('2','y') DEFAULT 2, e3 ENUM('ae','b') COLLATE utf8mb4_de_pb_0900_ai_ci DEFAULT 'ä', s1 SET('a','b') DEFAULT 'B,A,b')",
			code: 0,
			contains: []string{
				"`e1` enum('a','b') DEFAULT 'a'",
				"`e2` enum('2','y') DEFAULT 'y'",
				"COLLATE utf8mb4_de_pb_0900_ai_ci DEFAULT 'ae'",
				"`s1` set('a','b') DEFAULT 'a,b'",
			},
		},
		// 19. Single-digit time_zone offsets, datadir tablespace paths, sql_mode
		// user-variable validation, and UCS-2 surrogate defaults.
		{
			name:    "SET time_zone '+5:00' accepted and enforces TIMESTAMP UTC lower bound",
			session: []string{"SET time_zone = '+5:00'"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 05:00:00')",
			code:    1067,
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name:     "SET time_zone '-5:00' accepted and shifts offset DATETIME default",
			session:  []string{"SET time_zone = '-5:00'"},
			stmt:     "CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+00:00')",
			code:     0,
			contains: []string{"DEFAULT '2023-12-31 19:00:00'"},
			cleanup:  []string{"SET SESSION time_zone = DEFAULT"},
		},
		{
			name:    "CREATE TABLESPACE in /var/lib/mysql datadir accepted",
			stmt:    "CREATE TABLESPACE ts_datadir_abs ADD DATAFILE '/var/lib/mysql/ts_datadir_abs.ibd' ENGINE=InnoDB",
			code:    0,
			cleanup: []string{"DROP TABLESPACE ts_datadir_abs"},
		},
		{
			name: "CREATE TABLESPACE with circular .. directory rejected with 3121",
			stmt: "CREATE TABLESPACE ts_circ ADD DATAFILE '/var/lib/mysql/no_such_dir/../ts_circ.ibd' ENGINE=InnoDB",
			code: 3121,
		},
		{
			name:    "SET sql_mode to NULL user variable rejected with 1231",
			session: []string{"SET @mode_null = NULL"},
			stmt:    "SET sql_mode = @mode_null",
			code:    1231,
		},
		{
			name: "SET sql_mode to unset user variable rejected with 1231",
			stmt: "SET sql_mode = @omni_unset_mode",
			code: 1231,
		},
		{
			name:    "SET sql_mode to invalid string user variable rejected with 1231",
			session: []string{"SET @mode_bad = 'NOT_A_REAL_MODE'"},
			stmt:    "SET sql_mode = @mode_bad",
			code:    1231,
		},
		{
			name:    "SET sql_mode to float user variable rejected with 1232",
			session: []string{"SET @mode_flt = 1.5"},
			stmt:    "SET sql_mode = @mode_flt",
			code:    1232,
		},
		{
			name: "CREATE TABLE CHAR(1) ucs2 DEFAULT 0xD800 rejected with 3507",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET ucs2 DEFAULT 0xD800)",
			code: 3507,
		},
		{
			name: "CREATE TABLE CHAR(1) ucs2 DEFAULT 0xD83DDE00 rejected for width first with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET ucs2 DEFAULT 0xD83DDE00)",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(2) ucs2 DEFAULT 0xD83DDE00 rejected for surrogates with 3507",
			stmt: "CREATE TABLE a (c CHAR(2) CHARACTER SET ucs2 DEFAULT 0xD83DDE00)",
			code: 3507,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			if len(tc.cleanup) > 0 {
				t.Cleanup(func() {
					for _, stmt := range tc.cleanup {
						if _, err := conn.ExecContext(ctx, stmt); err != nil {
							t.Fatalf("cleanup %s: %v", stmt, err)
						}
					}
				})
			}
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
			if tc.code == 0 && tc.keeps != "" {
				tbl := c.GetDatabase("d").GetTable("a")
				if tbl == nil || tbl.Partitioning == nil || tbl.Partitioning.Type != tc.keeps {
					t.Fatalf("catalog did not keep the %s partitioning %s wrote: %+v", tc.keeps, tc.stmt, tbl)
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
	var syntax *mysqlparser.ParseError
	if errors.As(err, &syntax) {
		// A statement the parser refuses is the code it recorded, 1064 for a
		// syntax error: see mysql/catalog/simulate_test.go's parseCode.
		code, _ := strconv.Atoi(parseCode(syntax))
		return code, syntax.Message
	}
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

func TestSimulationRound4CatalogInvariants(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d DEFAULT CHARSET=latin1 COLLATE=latin1_bin; USE d;")
	execAlterClauses(t, c, "CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4;")

	// CONVERT TO CHARSET DEFAULT inherits the database's non-utf8mb4 charset and collation.
	execAlterClauses(t, c, "ALTER TABLE a CONVERT TO CHARSET DEFAULT;")
	tbl := c.GetDatabase("d").GetTable("a")
	if tbl.Charset != "latin1" || tbl.Collation != "latin1_bin" {
		t.Fatalf("after CONVERT TO CHARSET DEFAULT got (%q, %q), want (latin1, latin1_bin)", tbl.Charset, tbl.Collation)
	}
	colS := tbl.GetColumn("s")
	if colS.Charset != "latin1" || colS.Collation != "latin1_bin" || colS.CharsetExplicit || colS.CollationExplicit {
		t.Fatalf("column s after CONVERT TO CHARSET DEFAULT: %+v", colS)
	}

	// Single-command ALTER failures leave the table completely untouched.
	before := c.ShowCreateTable("d", "a")
	for _, bad := range []struct {
		stmt string
		code int
	}{
		{"ALTER TABLE a ADD x INT AUTO_INCREMENT", 1075},
		{"ALTER TABLE a ADD x INT DEFAULT (x)", 3767},
		{"ALTER TABLE a AUTOEXTEND_SIZE 1", 4025},
		{"ALTER TABLE a ADD e ENUM('x', 'x')", 1291},
	} {
		got, msg := simulationCatalogRefusal(t, c, bad.stmt)
		if got != bad.code {
			t.Fatalf("%s: got %d (%s), want %d", bad.stmt, got, msg, bad.code)
		}
		if after := c.ShowCreateTable("d", "a"); after != before {
			t.Fatalf("failed %s mutated table:\nbefore:\n%s\nafter:\n%s", bad.stmt, before, after)
		}
	}

	// Tablespace DDL lifecycle: create, duplicate (1813), rename, non-empty drop (3120), drop, missing drop (3510).
	execAlterClauses(t, c, "CREATE TABLESPACE ts1 ADD DATAFILE 'ts1.ibd' AUTOEXTEND_SIZE = 4M;")
	if got, _ := simulationCatalogRefusal(t, c, "CREATE TABLESPACE ts1 ADD DATAFILE 'ts1b.ibd'"); got != 1813 {
		t.Fatalf("duplicate CREATE TABLESPACE got %d, want 1813", got)
	}
	execAlterClauses(t, c, "ALTER TABLE a TABLESPACE ts1;")
	execAlterClauses(t, c, "ALTER TABLESPACE ts1 RENAME TO ts2;")
	if tbl.Tablespace != "ts2" {
		t.Fatalf("table tablespace after rename = %q, want ts2", tbl.Tablespace)
	}
	if got, _ := simulationCatalogRefusal(t, c, "DROP TABLESPACE ts2"); got != 3120 {
		t.Fatalf("DROP non-empty TABLESPACE got %d, want 3120", got)
	}
	execAlterClauses(t, c, "ALTER TABLE a TABLESPACE innodb_file_per_table;")
	execAlterClauses(t, c, "DROP TABLESPACE ts2;")
	if got, _ := simulationCatalogRefusal(t, c, "DROP TABLESPACE ts2"); got != 3510 {
		t.Fatalf("DROP missing TABLESPACE got %d, want 3510", got)
	}

	// Snapshot-loaded partition/subpartition tablespace ownership blocks DROP TABLESPACE (3120) and tracks RENAME.
	execAlterClauses(t, c, "CREATE TABLESPACE ts_p ADD DATAFILE 'ts_p.ibd';")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_snap (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = ts_p)")); err != nil {
		t.Fatalf("DefineTable with partition tablespace: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, c, "DROP TABLESPACE ts_p"); got != 3120 {
		t.Fatalf("DROP partition-referenced TABLESPACE got %d, want 3120", got)
	}
	execAlterClauses(t, c, "ALTER TABLESPACE ts_p RENAME TO ts_p2;")
	if pSnap := c.GetDatabase("d").GetTable("p_snap"); pSnap.Partitioning.Partitions[0].Tablespace != "ts_p2" {
		t.Fatalf("partition tablespace after rename = %q, want ts_p2", pSnap.Partitioning.Partitions[0].Tablespace)
	}
	pSnapDDL := c.ShowCreateTable("d", "p_snap")
	if !strings.Contains(pSnapDDL, "TABLESPACE = `ts_p2`") {
		t.Fatalf("ShowCreateTable omitted partition TABLESPACE:\n%s", pSnapDDL)
	}
	if got, _ := simulationCatalogRefusal(t, c, "DROP TABLESPACE ts_p2"); got != 3120 {
		t.Fatalf("DROP renamed partition-referenced TABLESPACE got %d, want 3120", got)
	}

	// Subpartition TABLESPACE also renders in ShowCreateTable and survives reload.
	execAlterClauses(t, c, "CREATE TABLESPACE ts_sub ADD DATAFILE 'ts_sub.ibd';")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_sub_snap (b INT) PARTITION BY RANGE (b) SUBPARTITION BY HASH (b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION sp0 TABLESPACE = ts_sub, SUBPARTITION sp1))")); err != nil {
		t.Fatalf("DefineTable with subpartition tablespace: %v", err)
	}
	pSubSnapDDL := c.ShowCreateTable("d", "p_sub_snap")
	if !strings.Contains(pSubSnapDDL, "SUBPARTITION sp0 TABLESPACE = `ts_sub` ENGINE = InnoDB") {
		t.Fatalf("ShowCreateTable omitted subpartition TABLESPACE:\n%s", pSubSnapDDL)
	}
	if strings.Contains(pSubSnapDDL, "SUBPARTITION sp1 TABLESPACE") {
		t.Fatalf("ShowCreateTable unexpectedly emitted TABLESPACE on sp1:\n%s", pSubSnapDDL)
	}
	reloaded := New()
	execAlterClauses(t, reloaded, "CREATE DATABASE d; USE d; CREATE TABLESPACE ts_p2 ADD DATAFILE 'ts_p2.ibd'; CREATE TABLESPACE ts_sub ADD DATAFILE 'ts_sub.ibd';")
	if err := reloaded.DefineTable(mustParseTable(t, pSnapDDL)); err != nil {
		t.Fatalf("reload p_snap ShowCreateTable: %v", err)
	}
	if err := reloaded.DefineTable(mustParseTable(t, pSubSnapDDL)); err != nil {
		t.Fatalf("reload p_sub_snap ShowCreateTable: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, reloaded, "DROP TABLESPACE ts_p2"); got != 3120 {
		t.Fatalf("reloaded DROP partition-referenced TABLESPACE got %d, want 3120", got)
	}
	if got, _ := simulationCatalogRefusal(t, reloaded, "DROP TABLESPACE ts_sub"); got != 3120 {
		t.Fatalf("reloaded DROP subpartition-referenced TABLESPACE got %d, want 3120", got)
	}

	// BINARY and VARBINARY numeric defaults coerce and replay cleanly via ShowCreateTable.
	execAlterClauses(t, c, "CREATE TABLE bin_num_def (b1 BINARY(1) DEFAULT 01, vb2 VARBINARY(2) DEFAULT 007);")
	binNumDDL := c.ShowCreateTable("d", "bin_num_def")
	if !strings.Contains(binNumDDL, "`b1` binary(1) DEFAULT '1'") || !strings.Contains(binNumDDL, "`vb2` varbinary(2) DEFAULT '7'") {
		t.Fatalf("ShowCreateTable did not coerce BINARY/VARBINARY numeric defaults:\n%s", binNumDDL)
	}
	execAlterClauses(t, reloaded, binNumDDL)

	// Snapshot-loaded non-InnoDB partitioned table with partition TABLESPACE is
	// validated when converted to ENGINE=InnoDB (3510 when missing, 1478 when shared).
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_myisam_miss (b INT) ENGINE=MyISAM PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = ts_gone)")); err != nil {
		t.Fatalf("DefineTable non-InnoDB partition missing tablespace: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, c, "ALTER TABLE p_myisam_miss ENGINE=InnoDB"); got != 3510 {
		t.Fatalf("ALTER non-InnoDB partitioned table with missing partition tablespace to InnoDB got %d, want 3510", got)
	}
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_myisam_shared (b INT) ENGINE=MyISAM PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = ts_p2)")); err != nil {
		t.Fatalf("DefineTable non-InnoDB partition shared tablespace: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, c, "ALTER TABLE p_myisam_shared ENGINE=InnoDB"); got != 1478 {
		t.Fatalf("ALTER non-InnoDB partitioned table with shared partition tablespace to InnoDB got %d, want 1478", got)
	}
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_myisam_add (b INT) ENGINE=MyISAM PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))")); err != nil {
		t.Fatalf("DefineTable non-InnoDB partitioned table: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, c, "ALTER TABLE p_myisam_add ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) TABLESPACE = ts_gone), ENGINE=InnoDB"); got != 3510 {
		t.Fatalf("ALTER non-InnoDB partitioned table ADD PARTITION missing tablespace + ENGINE=InnoDB got %d, want 3510", got)
	}
	if got, _ := simulationCatalogRefusal(t, c, "ALTER TABLE p_myisam_add REORGANIZE PARTITION p0 INTO (PARTITION p0a VALUES LESS THAN (10) TABLESPACE = ts_p2), ENGINE=InnoDB"); got != 1478 {
		t.Fatalf("ALTER non-InnoDB partitioned table REORGANIZE PARTITION shared tablespace + ENGINE=InnoDB got %d, want 1478", got)
	}

	// SET GLOBAL time_zone updates global default without clobbering session time_zone,
	// and numeric temporal defaults round-trip cleanly via ShowCreateTable.
	execAlterClauses(t, c, "SET SESSION time_zone = '+00:00'; SET GLOBAL time_zone = '+14:00';")
	if c.session.TimeZone != "+00:00" || c.globalTimeZone != "+14:00" {
		t.Fatalf("unexpected time zones after SET GLOBAL: session=%q global=%q", c.session.TimeZone, c.globalTimeZone)
	}
	if got, _ := simulationCatalogRefusal(t, c, "SET @tz_atomic = '+05:00', time_zone = @tz_atomic"); got != 1231 {
		t.Fatalf("single-SET read-before-update @tz_atomic got %d, want 1231", got)
	}
	if _, exists := c.userVars["tz_atomic"]; exists {
		t.Fatalf("failed SET mutated @tz_atomic: %+v", c.userVars["tz_atomic"])
	}
	execAlterClauses(t, c, "SET SESSION time_zone = DEFAULT;")
	if c.session.TimeZone != "+14:00" {
		t.Fatalf("SET SESSION time_zone = DEFAULT got %q, want +14:00", c.session.TimeZone)
	}
	execAlterClauses(t, c, "SET GLOBAL time_zone = DEFAULT; SET SESSION time_zone = DEFAULT;")
	execAlterClauses(t, c, "CREATE TABLE num_temp_def (d DATE DEFAULT 10101, dt DATETIME(2) DEFAULT 101010101.25, tm TIME DEFAULT 10203);")
	numTempDDL := c.ShowCreateTable("d", "num_temp_def")
	for _, want := range []string{"`d` date DEFAULT '2001-01-01'", "`dt` datetime(2) DEFAULT '2000-01-01 01:01:01.25'", "`tm` time DEFAULT '01:02:03'"} {
		if !strings.Contains(numTempDDL, want) {
			t.Fatalf("ShowCreateTable missing %s:\n%s", want, numTempDDL)
		}
	}
	execAlterClauses(t, reloaded, numTempDDL)

	// LoadSDL accepts InnoDB tables and partitions referencing general tablespaces
	// without CREATE TABLESPACE statements in SDL, and registers them on the catalog.
	sdlCat, err := LoadSDL(`
		CREATE DATABASE d_sdl;
		USE d_sdl;
		CREATE TABLE t_ts (id INT PRIMARY KEY) ENGINE=InnoDB TABLESPACE ts_general;
		CREATE TABLE t_part_ts (id INT) ENGINE=InnoDB PARTITION BY RANGE (id) (
			PARTITION p0 VALUES LESS THAN (10) TABLESPACE = ts_part
		);
	`)
	if err != nil {
		t.Fatalf("LoadSDL with general tablespace failed: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, sdlCat, "DROP TABLESPACE ts_general"); got != 3120 {
		t.Fatalf("DROP SDL-registered general tablespace got %d, want 3120", got)
	}
	if got, _ := simulationCatalogRefusal(t, sdlCat, "DROP TABLESPACE ts_part"); got != 3120 {
		t.Fatalf("DROP SDL-registered partition tablespace got %d, want 3120", got)
	}

	// LoadSDL round-trips ShowCreateTable output with /*!80023 AUTOEXTEND_SIZE=... */
	// and offset-normalized temporal defaults.
	execAlterClauses(t, c, "SET SESSION time_zone = '+08:00'; CREATE TABLE t_ae_sdl (id INT PRIMARY KEY, dt DATETIME DEFAULT '2024-01-01 00:00:00+14:00') AUTOEXTEND_SIZE=16M;")
	aeSDL := c.ShowCreateTable("d", "t_ae_sdl")
	if !strings.Contains(aeSDL, "/*!80023 AUTOEXTEND_SIZE=16777216 */ ENGINE=InnoDB") || !strings.Contains(aeSDL, "DEFAULT '2023-12-31 18:00:00'") {
		t.Fatalf("unexpected ShowCreateTable for t_ae_sdl:\n%s", aeSDL)
	}
	sdlCat2, err := LoadSDL("CREATE DATABASE d;\nUSE d;\n" + aeSDL + ";")
	if err != nil {
		t.Fatalf("LoadSDL with AUTOEXTEND_SIZE failed: %v", err)
	}
	if gotSDL := sdlCat2.ShowCreateTable("d", "t_ae_sdl"); gotSDL != aeSDL {
		t.Fatalf("LoadSDL round-trip mismatch:\ngot:\n%s\nwant:\n%s", gotSDL, aeSDL)
	}
}
