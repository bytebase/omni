package catalog

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/mysql/validate"
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
		for _, stmt := range []string{"SET NAMES utf8mb4", "SET SESSION sql_mode = DEFAULT", "DROP VIEW IF EXISTS a", "DROP TABLE IF EXISTS a", "DROP TABLE IF EXISTS p", "DROP TABLE IF EXISTS s", "DROP TABLE IF EXISTS z", "DROP VIEW IF EXISTS v", "SET SESSION sql_generate_invisible_primary_key = OFF", "DROP DATABASE IF EXISTS cteother", "DROP TABLE IF EXISTS `dual`", "DROP TABLE IF EXISTS c2", "DROP EVENT IF EXISTS e"} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
	}
	t.Cleanup(reset)

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	// Tables for the functional key part cases.
	fiTable := "CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))"
	fiGeometry := "CREATE TABLE s (g GEOMETRY NOT NULL SRID 0)"
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
		// 20. Numeric ENUM/SET DDL string-membership check, column-charset binary
		// ENUM/SET decoding, and @@sql_mode DEFAULT expansion.
		{
			name: "CREATE TABLE ENUM('a','b') DEFAULT 1 rejected with 1067 when '1' not a member",
			stmt: "CREATE TABLE a (c ENUM('a','b') DEFAULT 1)",
			code: 1067,
		},
		{
			name: "CREATE TABLE SET('a','b') DEFAULT 3 rejected with 1067 when '3' not a member",
			stmt: "CREATE TABLE a (c SET('a','b') DEFAULT 3)",
			code: 1067,
		},
		{
			name:     "CREATE TABLE ENUM sjis DEFAULT 0x82A0 decodes in column charset",
			stmt:     "CREATE TABLE a (c ENUM('あ') CHARACTER SET sjis DEFAULT 0x82A0)",
			code:     0,
			contains: []string{"DEFAULT 'あ'"},
		},
		{
			name:     "CREATE TABLE SET sjis DEFAULT 0x82A22C82A0 decodes and orders in column charset",
			stmt:     "CREATE TABLE a (c SET('あ','い') CHARACTER SET sjis DEFAULT 0x82A22C82A0)",
			code:     0,
			contains: []string{"DEFAULT 'あ,い'"},
		},
		{
			name:     "CREATE TABLE ENUM utf16 and latin1 hex defaults decode in column charset",
			stmt:     "CREATE TABLE a (e1 ENUM('a','b') CHARACTER SET utf16 DEFAULT 0x41, e2 ENUM('é','b') CHARACTER SET latin1 DEFAULT 0xE9)",
			code:     0,
			contains: []string{"DEFAULT 'a'", "DEFAULT 'é'"},
		},
		{
			name:    "SET sql_mode = @@sql_mode and @saved_sql_mode round-trip default mode",
			session: []string{"SET @saved_sql_mode = @@sql_mode", "SET sql_mode = ''", "SET sql_mode = @saved_sql_mode", "SET sql_mode = @@sql_mode"},
			stmt:    "CREATE TABLE a (c DATE DEFAULT '0000-00-00')",
			code:    1067,
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
		},
		// 21. Read-only innodb_directories SQL assignment refusal and FULLTEXT prefix clearing.
		{
			name: "SET GLOBAL innodb_directories rejected with 1238",
			stmt: "SET GLOBAL innodb_directories = '/tmp'",
			code: 1238,
		},
		{
			name:     "CREATE TABLE FULLTEXT key ignores oversized prefix length",
			stmt:     "CREATE TABLE a (c VARCHAR(10), FULLTEXT KEY fk1 (c(100)))",
			code:     0,
			contains: []string{"FULLTEXT KEY `fk1` (`c`)"},
			omits:    []string{"(100)"},
		},
		{
			name:     "ALTER TABLE ADD FULLTEXT index ignores oversized prefix length",
			setup:    []string{"CREATE TABLE a (c VARCHAR(10))"},
			stmt:     "ALTER TABLE a ADD FULLTEXT INDEX fk2 (c(100))",
			code:     0,
			contains: []string{"FULLTEXT KEY `fk2` (`c`)"},
			omits:    []string{"(100)"},
		},
		{
			name:     "CREATE FULLTEXT INDEX ignores oversized prefix length",
			setup:    []string{"CREATE TABLE a (c VARCHAR(10))"},
			stmt:     "CREATE FULLTEXT INDEX fk3 ON a (c(100))",
			code:     0,
			contains: []string{"FULLTEXT KEY `fk3` (`c`)"},
			omits:    []string{"(100)"},
		},
		// 22. Quoted CHAR/VARCHAR charset representability and stored/disallowed functions in DEFAULT (expr).
		{
			name: "CREATE TABLE CHAR sjis DEFAULT unrepresentable emoji rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT '😀')",
			code: 1067,
		},
		{
			name:    "CREATE TABLE CHAR sjis DEFAULT unrepresentable emoji rejected even in non-strict sql_mode",
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT '😀')",
			code:    1067,
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
		},
		{
			name: "CREATE TABLE CHAR cp1251 DEFAULT 'é' rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET cp1251 DEFAULT 'é')",
			code: 1067,
		},
		{
			name:     "CREATE TABLE CHAR cp1251 DEFAULT 'Д' accepted",
			stmt:     "CREATE TABLE a (c CHAR(1) CHARACTER SET cp1251 DEFAULT 'Д')",
			code:     0,
			contains: []string{"DEFAULT 'Д'"},
		},
		{
			name: "CREATE TABLE CHAR swe7 DEFAULT '[' rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET swe7 DEFAULT '[')",
			code: 1067,
		},
		{
			name:     "CREATE TABLE CHAR swe7 DEFAULT 'Å' accepted",
			stmt:     "CREATE TABLE a (c CHAR(1) CHARACTER SET swe7 DEFAULT 'Å')",
			code:     0,
			contains: []string{"DEFAULT 'Å'"},
		},
		{
			name:     "CREATE TABLE CHAR latin1 DEFAULT '€' accepted",
			stmt:     "CREATE TABLE a (c CHAR(1) CHARACTER SET latin1 DEFAULT '€')",
			code:     0,
			contains: []string{"DEFAULT '€'"},
		},
		{
			name:     "CREATE TABLE CHAR gb18030 DEFAULT '😀' accepted",
			stmt:     "CREATE TABLE a (c CHAR(4) CHARACTER SET gb18030 DEFAULT '😀')",
			code:     0,
			contains: []string{"`c` char(4) CHARACTER SET gb18030"},
		},
		{
			name:    "CREATE TABLE DEFAULT calling stored function rejected with 3770",
			setup:   []string{"CREATE FUNCTION f() RETURNS INT DETERMINISTIC RETURN 1"},
			cleanup: []string{"DROP FUNCTION IF EXISTS f"},
			stmt:    "CREATE TABLE a (c INT DEFAULT (f()))",
			code:    3770,
		},
		{
			name: "CREATE TABLE DEFAULT calling unknown function rejected with 3770",
			stmt: "CREATE TABLE a (c INT DEFAULT (unknown_fn()))",
			code: 3770,
		},
		{
			name: "CREATE TABLE DEFAULT calling disallowed builtin VERSION() rejected with 3770",
			stmt: "CREATE TABLE a (c VARCHAR(64) DEFAULT (VERSION()))",
			code: 3770,
		},
		{
			name:  "ALTER TABLE ALTER COLUMN SET DEFAULT calling unknown function rejected with 3770",
			setup: []string{"CREATE TABLE a (c INT)"},
			stmt:  "ALTER TABLE a ALTER COLUMN c SET DEFAULT (unknown_fn())",
			code:  3770,
		},
		{
			name: "CREATE TABLE DEFAULT calling allowed builtins accepted",
			stmt: "CREATE TABLE a (x INT DEFAULT (ABS(-1)), y VARCHAR(36) DEFAULT (UUID()), z INT DEFAULT (IF(1,2,3)), g GEOMETRY DEFAULT (POINT(0,0)))",
			code: 0,
		},
		// 23. Literal introducer/national charset defaults, discarded-table ALTER blocking,
		//     FK check on DISCARD TABLESPACE, and slash rejection in tablespace names.
		{
			name: "CREATE TABLE utf8mb4 DEFAULT N supplementary rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT N'😀')",
			code: 1067,
		},
		{
			name: "CREATE TABLE utf8mb4 DEFAULT _utf8mb3 supplementary rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _utf8mb3'😀')",
			code: 1067,
		},
		{
			name: "CREATE TABLE utf8mb4 DEFAULT _ascii non-ASCII rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _ascii'é')",
			code: 1067,
		},
		{
			name: "CREATE TABLE utf8mb4 DEFAULT _sjis invalid sequence rejected with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _sjis'😀')",
			code: 1067,
		},
		{
			name: "CREATE TABLE CHAR(1) utf8mb4 DEFAULT _latin1 multi-byte UTF-8 input overflows width with 1067",
			stmt: "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _latin1'€')",
			code: 1067,
		},
		{
			name:     "CREATE TABLE CHAR(10) utf8mb4 DEFAULT _latin1 decodes introducer bytes",
			stmt:     "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT _latin1'€')",
			code:     0,
			contains: []string{"DEFAULT 'â‚¬'"},
		},
		{
			name:  "ALTER TABLE ADD COLUMN after DISCARD TABLESPACE rejected with 1814",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD COLUMN c INT",
			code:  1814,
		},
		{
			name:     "ALTER TABLE COMMENT after DISCARD TABLESPACE accepted",
			setup:    []string{"CREATE TABLE a (id INT PRIMARY KEY) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:     "ALTER TABLE a COMMENT = 'ok'",
			code:     0,
			contains: []string{"COMMENT='ok'"},
		},
		{
			name:    "ALTER TABLE DISCARD TABLESPACE on FK parent with foreign_key_checks=1 rejected with 1451",
			setup:   []string{"CREATE TABLE p_fk (id INT PRIMARY KEY)", "CREATE TABLE c_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES p_fk(id))"},
			cleanup: []string{"DROP TABLE IF EXISTS c_fk", "DROP TABLE IF EXISTS p_fk"},
			stmt:    "ALTER TABLE p_fk DISCARD TABLESPACE",
			code:    1451,
		},
		{
			name:    "ALTER TABLE DISCARD TABLESPACE on FK child with outbound FK accepted",
			setup:   []string{"CREATE TABLE p_fk (id INT PRIMARY KEY)", "CREATE TABLE c_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES p_fk(id))"},
			cleanup: []string{"DROP TABLE IF EXISTS c_fk", "DROP TABLE IF EXISTS p_fk"},
			stmt:    "ALTER TABLE c_fk DISCARD TABLESPACE",
			code:    0,
		},
		{
			name:    "ALTER TABLE DISCARD TABLESPACE on FK parent with foreign_key_checks=0 accepted",
			setup:   []string{"CREATE TABLE p_fk (id INT PRIMARY KEY)", "CREATE TABLE c_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES p_fk(id))"},
			session: []string{"SET foreign_key_checks = 0"},
			cleanup: []string{"SET foreign_key_checks = 1", "DROP TABLE IF EXISTS c_fk", "DROP TABLE IF EXISTS p_fk"},
			stmt:    "ALTER TABLE p_fk DISCARD TABLESPACE",
			code:    0,
		},
		{
			name: "CREATE TABLESPACE with slash in name rejected with 3119",
			stmt: "CREATE TABLESPACE `a/b` ENGINE=InnoDB",
			code: 3119,
		},
		{
			name:    "ALTER TABLESPACE RENAME TO with slash in name rejected with 3119",
			setup:   []string{"CREATE TABLESPACE ts_ren_slash ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLESPACE ts_ren_slash"},
			stmt:    "ALTER TABLESPACE ts_ren_slash RENAME TO `a/b`",
			code:    3119,
		},
		// 24. Extended time_zone offset digits accepted by MySQL, arithmetic user-variable
		//     SET expressions, and empty ENCRYPTION='' rejection.
		{
			name:    "SET time_zone '+005:00' accepted and normalizes to +05:00",
			session: []string{"SET time_zone = '+005:00'"},
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 05:00:00')",
			code:    1067,
		},
		{
			name:    "SET time_zone '+5:0' accepted and normalizes to +05:00",
			session: []string{"SET time_zone = '+5:0'"},
			cleanup: []string{"SET SESSION time_zone = DEFAULT"},
			stmt:    "CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 05:00:00')",
			code:    1067,
		},
		{
			name:     "SET @mode = 2 + 2 evaluates arithmetic for SET sql_mode",
			session:  []string{"SET @mode = 2 + 2", "SET sql_mode = @mode"},
			cleanup:  []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:     "CREATE TABLE a (c VARCHAR(10) DEFAULT 'abc')",
			code:     0,
			contains: []string{"DEFAULT 'abc'"},
		},
		{
			name: "CREATE TABLESPACE with empty ENCRYPTION rejected with 3184",
			stmt: "CREATE TABLESPACE ts_emp_enc ADD DATAFILE 'ts_emp_enc.ibd' ENCRYPTION=''",
			code: 3184,
		},
		// 25. Division preserving fractional decimals, ALTER TABLESPACE empty DATAFILE
		//     clauses, and general tablespace FILE_BLOCK_SIZE compatibility.
		{
			name:    "SET @m = 3 / 2 preserves fractional decimal so SET sql_mode rejects with 1232",
			session: []string{"SET @m = 3 / 2"},
			stmt:    "SET sql_mode = @m",
			code:    1232,
		},
		{
			name:    "SET @m = 4 / 2 preserves fractional decimal so SET sql_mode rejects with 1232",
			session: []string{"SET @m = 4 / 2"},
			stmt:    "SET sql_mode = @m",
			code:    1232,
		},
		{
			name:    "ALTER TABLESPACE ADD DATAFILE empty filename rejected with 1533",
			setup:   []string{"CREATE TABLESPACE ts_alt_emp ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLESPACE ts_alt_emp"},
			stmt:    "ALTER TABLESPACE ts_alt_emp ADD DATAFILE ''",
			code:    1533,
		},
		{
			name:    "ALTER TABLESPACE DROP DATAFILE empty filename rejected with 3629",
			setup:   []string{"CREATE TABLESPACE ts_alt_drp ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLESPACE ts_alt_drp"},
			stmt:    "ALTER TABLESPACE ts_alt_drp DROP DATAFILE ''",
			code:    3629,
		},
		{
			name:    "CREATE TABLE COMPRESSED KEY_BLOCK_SIZE=4 in FILE_BLOCK_SIZE=8K tablespace rejected with 1478",
			setup:   []string{"CREATE TABLESPACE ts_fbs8 ADD DATAFILE 'ts_fbs8.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_fbs8"},
			stmt:    "CREATE TABLE a (id INT) TABLESPACE ts_fbs8 ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=4",
			code:    1478,
		},
		{
			name:    "CREATE TABLE COMPRESSED KEY_BLOCK_SIZE=8 in FILE_BLOCK_SIZE=8K tablespace accepted",
			setup:   []string{"CREATE TABLESPACE ts_fbs8_ok ADD DATAFILE 'ts_fbs8_ok.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_fbs8_ok"},
			stmt:    "CREATE TABLE a (id INT) TABLESPACE ts_fbs8_ok ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8",
			code:    0,
		},
		{
			name:    "CREATE TABLE uncompressed in FILE_BLOCK_SIZE=8K tablespace rejected with 1478",
			setup:   []string{"CREATE TABLESPACE ts_fbs8_uncomp ADD DATAFILE 'ts_fbs8_uncomp.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_fbs8_uncomp"},
			stmt:    "CREATE TABLE a (id INT) TABLESPACE ts_fbs8_uncomp",
			code:    1478,
		},
		{
			name:    "CREATE TABLE COMPRESSED in default 16K general tablespace rejected with 1478",
			setup:   []string{"CREATE TABLESPACE ts_fbs16 ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_fbs16"},
			stmt:    "CREATE TABLE a (id INT) TABLESPACE ts_fbs16 ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8",
			code:    1478,
		},
		{
			name:    "ALTER TABLE ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=4 into FILE_BLOCK_SIZE=8K tablespace rejected with 1478",
			setup:   []string{"CREATE TABLESPACE ts_fbs8_alt ADD DATAFILE 'ts_fbs8_alt.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB", "CREATE TABLE a (id INT) ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_fbs8_alt"},
			stmt:    "ALTER TABLE a TABLESPACE ts_fbs8_alt ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=4",
			code:    1478,
		},
		{
			name:    "SET NAMES ascii rejects utf8mb4 quoted literal DEFAULT with 1067",
			session: []string{"SET NAMES ascii"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:    1067,
		},
		{
			name:    "SET NAMES ascii rejects ENUM non-ascii DEFAULT with 1067",
			session: []string{"SET NAMES ascii"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (s ENUM('a') DEFAULT 'é')",
			code:    1067,
		},
		{
			name:    "SET NAMES sjis rejects 4-byte utf8mb4 quoted DEFAULT with 1067",
			session: []string{"SET NAMES sjis"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT '😀')",
			code:    1067,
		},
		{
			name:    "SET NAMES latin1 interprets UTF-8 bytes as 2 chars so CHAR(1) rejects with 1067",
			session: []string{"SET NAMES latin1"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:    1067,
		},
		{
			name:    "SET NAMES latin1 ALTER TABLE ADD COLUMN CHAR(1) DEFAULT rejects with 1067",
			setup:   []string{"CREATE TABLE a (id INT)"},
			session: []string{"SET NAMES latin1"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "ALTER TABLE a ADD COLUMN d CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é'",
			code:    1067,
		},
		{
			name:    "SET NAMES latin1 ALTER COLUMN SET DEFAULT on CHAR(1) rejects with 1067",
			setup:   []string{"CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4)"},
			session: []string{"SET NAMES latin1"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "ALTER TABLE a ALTER COLUMN c SET DEFAULT 'é'",
			code:    1067,
		},
		{
			name:    "SET NAMES binary accepts valid raw utf8mb4 bytes in utf8mb4 column DEFAULT",
			session: []string{"SET NAMES binary"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:    0,
		},
		{
			name:    "SET NAMES binary rejects ENUM DEFAULT not in list with 1067",
			session: []string{"SET NAMES binary"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (s ENUM('a') DEFAULT 'é')",
			code:    1067,
		},
		{
			name:     "SET character_set_client=ascii with utf8mb4 connection replaces non-ascii with question marks",
			session:  []string{"SET NAMES utf8mb4", "SET character_set_client = ascii"},
			cleanup:  []string{"SET NAMES utf8mb4"},
			stmt:     "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:     0,
			contains: []string{"DEFAULT '??'"},
		},
		{
			name:     "SET character_set_client=utf8mb4 with ascii connection replaces non-ascii with question mark",
			session:  []string{"SET NAMES ascii", "SET character_set_client = utf8mb4"},
			cleanup:  []string{"SET NAMES utf8mb4"},
			stmt:     "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:     0,
			contains: []string{"DEFAULT '?'"},
		},
		{
			name:    "SET CHARACTER SET latin1 with utf8mb4 connection decodes latin1 bytes for CHAR(1) length check",
			session: []string{"SET CHARACTER SET latin1"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:    1067,
		},
		{
			name: "SET NAMES unknown charset rejected with 1115",
			stmt: "SET NAMES nonexistent",
			code: 1115,
		},
		{
			name: "SET character_set_client unknown charset rejected with 1115",
			stmt: "SET character_set_client = nonexistent",
			code: 1115,
		},
		{
			name:  "ALTER TABLE DROP PRIMARY KEY after DISCARD TABLESPACE rejected with 1814",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY)", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a DROP PRIMARY KEY",
			code:  1814,
		},
		{
			name:    "non-strict sql_mode still rejects out-of-range TINYINT DEFAULT 999 with 1067",
			session: []string{"SET sql_mode=''"},
			stmt:    "CREATE TABLE a (a TINYINT DEFAULT 999)",
			code:    1067,
		},
		{
			name:     "BINARY(2) DEFAULT _latin1'é' counts encoded bytes as 2",
			stmt:     "CREATE TABLE a (b BINARY(2) DEFAULT _latin1'é')",
			code:     0,
			contains: []string{"DEFAULT 'é'"},
		},
		{
			name:     "VARBINARY(2) DEFAULT _latin1'é' counts encoded bytes as 2",
			stmt:     "CREATE TABLE a (b VARBINARY(2) DEFAULT _latin1'é')",
			code:     0,
			contains: []string{"DEFAULT 'é'"},
		},
		{
			name:    "SET NAMES latin1 BINARY(2) DEFAULT 'é' accepted",
			session: []string{"SET NAMES latin1"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (b BINARY(2) DEFAULT 'é')",
			code:    0,
		},
		{
			name: "BINARY(2) DEFAULT _latin1'aé' exceeds 2 encoded bytes rejected with 1067",
			stmt: "CREATE TABLE a (b BINARY(2) DEFAULT _latin1'aé')",
			code: 1067,
		},
		{
			name: "TINYBLOB prefix length 256 rejected with 1071",
			stmt: "CREATE TABLE a (c TINYBLOB, KEY(c(256)))",
			code: 1071,
		},
		{
			name: "TINYTEXT utf8mb4 prefix length 64 exceeds 255 bytes rejected with 1071",
			stmt: "CREATE TABLE a (c TINYTEXT CHARACTER SET utf8mb4, KEY(c(64)))",
			code: 1071,
		},
		{
			name: "TINYTEXT utf8mb4 prefix length 63 fits 255 bytes accepted",
			stmt: "CREATE TABLE a (c TINYTEXT CHARACTER SET utf8mb4, KEY(c(63)))",
			code: 0,
		},
		{
			name: "BLOB prefix length 3073 exceeds InnoDB key limit rejected with 1071",
			stmt: "CREATE TABLE a (c BLOB, KEY(c(3073)))",
			code: 1071,
		},
		{
			name: "BLOB MyISAM prefix length 1001 exceeds 1000 bytes rejected with 1071",
			stmt: "CREATE TABLE a (c BLOB, KEY(c(1001))) ENGINE=MyISAM",
			code: 1071,
		},
		{
			name:     "non-strict sql_mode clamps TINYBLOB prefix length 256 to 255",
			session:  []string{"SET sql_mode=''"},
			stmt:     "CREATE TABLE a (c TINYBLOB, KEY(c(256)))",
			code:     0,
			contains: []string{"KEY `c` (`c`(255))"},
		},
		{
			name:    "CREATE TEMPORARY TABLE on general tablespace rejected with 1478",
			setup:   []string{"CREATE TABLESPACE ts_tmp_tbl ADD DATAFILE 'ts_tmp_tbl.ibd' ENGINE=InnoDB"},
			cleanup: []string{"DROP TABLESPACE ts_tmp_tbl"},
			stmt:    "CREATE TEMPORARY TABLE a (id INT) TABLESPACE ts_tmp_tbl",
			code:    1478,
		},
		{
			name: "CREATE TEMPORARY TABLE on innodb_system rejected with 1478",
			stmt: "CREATE TEMPORARY TABLE a (id INT) TABLESPACE innodb_system",
			code: 1478,
		},
		{
			name: "CREATE TEMPORARY TABLE on innodb_file_per_table rejected with 1478",
			stmt: "CREATE TEMPORARY TABLE a (id INT) TABLESPACE innodb_file_per_table",
			code: 1478,
		},
		{
			name:    "CREATE TEMPORARY TABLE on innodb_temporary accepted",
			cleanup: []string{"DROP TEMPORARY TABLE IF EXISTS a"},
			stmt:    "CREATE TEMPORARY TABLE a (id INT) TABLESPACE innodb_temporary",
			code:    0,
		},
		{
			name: "SET user var overflowing multiplication rejected with 1690",
			stmt: "SET @m = 4294967296 * 4294967296",
			code: 1690,
		},
		{
			name: "SET user var overflowing addition rejected with 1690",
			stmt: "SET @m = 9223372036854775807 + 1",
			code: 1690,
		},
		{
			name: "SET user var overflowing subtraction rejected with 1690",
			stmt: "SET @m = -9223372036854775808 - 1",
			code: 1690,
		},
		{
			name:     "SET collation_connection controls unintroduced string literal conversion",
			session:  []string{"SET NAMES utf8mb4", "SET collation_connection = latin1_swedish_ci"},
			cleanup:  []string{"SET NAMES utf8mb4"},
			stmt:     "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT '😀')",
			code:     0,
			contains: []string{"DEFAULT '?'"},
		},
		{
			name: "SET collation_connection unknown collation rejected with 1273",
			stmt: "SET collation_connection = nonexistent_coll",
			code: 1273,
		},
		{
			name: "CREATE TABLESPACE with non-JSON ENGINE_ATTRIBUTE rejected with 3980",
			stmt: "CREATE TABLESPACE ts_ea_bad ADD DATAFILE 'ts_ea_bad.ibd' ENGINE_ATTRIBUTE='not-json'",
			code: 3980,
		},
		{
			name:    "ALTER TABLESPACE with non-JSON ENGINE_ATTRIBUTE rejected with 3980",
			setup:   []string{"CREATE TABLESPACE ts_ea_ok ADD DATAFILE 'ts_ea_ok.ibd' ENGINE_ATTRIBUTE='{\"a\":1}'"},
			cleanup: []string{"DROP TABLESPACE ts_ea_ok"},
			stmt:    "ALTER TABLESPACE ts_ea_ok ENGINE_ATTRIBUTE='not-json'",
			code:    3980,
		},
		{
			name:    "CREATE UNDO TABLESPACE accepts relative datafile name under datadir",
			cleanup: []string{"ALTER UNDO TABLESPACE u_rel_ok SET INACTIVE", "SELECT SLEEP(1)", "DROP UNDO TABLESPACE u_rel_ok"},
			stmt:    "CREATE UNDO TABLESPACE u_rel_ok ADD DATAFILE 'u_rel_ok.ibu'",
			code:    0,
		},
		{
			name: "SET collation_connection rejects shaped nonexistent collation with 1273",
			stmt: "SET collation_connection = utf8mb4_not_a_real_collation_ci",
			code: 1273,
		},
		{
			name: "SET NAMES COLLATE rejects shaped nonexistent collation with 1273",
			stmt: "SET NAMES utf8mb4 COLLATE utf8mb4_not_a_real_collation_ci",
			code: 1273,
		},
		{
			name:     "BINARY(1) preserves 1-byte latin1 connection default from utf8mb4 client",
			session:  []string{"SET character_set_client = utf8mb4", "SET collation_connection = latin1_swedish_ci"},
			cleanup:  []string{"SET NAMES utf8mb4"},
			stmt:     "CREATE TABLE a (b BINARY(1) DEFAULT 'é')",
			code:     0,
			contains: []string{"DEFAULT 0xE9"},
		},
		{
			name:    "BINARY(2) preserves 2-byte sjis connection default from utf8mb4 client",
			session: []string{"SET character_set_client = utf8mb4", "SET collation_connection = sjis_japanese_ci"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (b BINARY(2) DEFAULT 'あ')",
			code:    0,
		},
		{
			name:    "BINARY(1) rejects 2-byte sjis connection default with 1067",
			session: []string{"SET character_set_client = utf8mb4", "SET collation_connection = sjis_japanese_ci"},
			cleanup: []string{"SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (b BINARY(1) DEFAULT 'あ')",
			code:    1067,
		},
		{
			name: "SET @p = @@SESSION.innodb_page_size rejected with 1238",
			stmt: "SET @p = @@SESSION.innodb_page_size",
			code: 1238,
		},
		{
			name: "SET @p = @@LOCAL.innodb_page_size rejected with 1238",
			stmt: "SET @p = @@LOCAL.innodb_page_size",
			code: 1238,
		},
		{
			name: "SET @p = @@PERSIST.innodb_page_size rejected with 1193",
			stmt: "SET @p = @@PERSIST.innodb_page_size",
			code: 1193,
		},
		{
			name: "SET @p = @@GLOBAL.innodb_page_size accepted",
			stmt: "SET @p = @@GLOBAL.innodb_page_size",
			code: 0,
		},
		{
			name: "default sql_mode rejects spaced builtin in DEFAULT expr with 3770",
			stmt: "CREATE TABLE a (a VARCHAR(1) DEFAULT (SUBSTRING ('x',1,1)))",
			code: 3770,
		},
		{
			name:    "IGNORE_SPACE sql_mode accepts spaced builtin in DEFAULT expr",
			session: []string{"SET sql_mode = 'IGNORE_SPACE'"},
			stmt:    "CREATE TABLE a (a VARCHAR(1) DEFAULT (SUBSTRING ('x',1,1)))",
			code:    0,
		},
		{
			name:    "ANSI sql_mode accepts spaced builtin in DEFAULT expr",
			session: []string{"SET sql_mode = 'ANSI'"},
			stmt:    "CREATE TABLE a (a VARCHAR(1) DEFAULT (SUBSTRING ('x',1,1)))",
			code:    0,
		},
		{
			name:    "IGNORE_SPACE sql_mode still rejects spaced UDF in DEFAULT expr with 3770",
			session: []string{"SET sql_mode = 'IGNORE_SPACE'"},
			stmt:    "CREATE TABLE a (a VARCHAR(1) DEFAULT (my_udf ('x')))",
			code:    3770,
		},
		{
			name: "SET character_set_client impermissible ucs2 rejected with 1231",
			stmt: "SET character_set_client = ucs2",
			code: 1231,
		},
		{
			name: "SET character_set_client impermissible utf16 rejected with 1231",
			stmt: "SET character_set_client = utf16",
			code: 1231,
		},
		{
			name: "SET character_set_client impermissible utf16le rejected with 1231",
			stmt: "SET character_set_client = utf16le",
			code: 1231,
		},
		{
			name: "SET character_set_client impermissible utf32 rejected with 1231",
			stmt: "SET character_set_client = utf32",
			code: 1231,
		},
		{
			name: "SET NAMES impermissible utf16 rejected with 1231",
			stmt: "SET NAMES utf16",
			code: 1231,
		},
		{
			name: "SET CHARACTER SET impermissible utf32 rejected with 1231",
			stmt: "SET CHARACTER SET utf32",
			code: 1231,
		},
		{
			name: "SET NAMES mismatched COLLATE rejected with 1253",
			stmt: "SET NAMES utf8mb4 COLLATE latin1_swedish_ci",
			code: 1253,
		},
		{
			name: "SET character_set_client NULL rejected with 1231",
			stmt: "SET character_set_client = NULL",
			code: 1231,
		},
		{
			name:    "SET GLOBAL character_set_client does not alter session client charset",
			session: []string{"SET GLOBAL character_set_client = latin1"},
			cleanup: []string{"SET GLOBAL character_set_client = DEFAULT", "SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:    0,
		},
		{
			name:    "SET PERSIST character_set_client does not alter session client charset",
			session: []string{"SET PERSIST character_set_client = latin1"},
			cleanup: []string{"SET GLOBAL character_set_client = DEFAULT", "SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:    0,
		},
		{
			name:    "SET GLOBAL collation_connection does not alter session connection charset",
			session: []string{"SET GLOBAL collation_connection = latin1_swedish_ci"},
			cleanup: []string{"SET GLOBAL collation_connection = DEFAULT", "SET NAMES utf8mb4"},
			stmt:    "CREATE TABLE a (b BINARY(2) DEFAULT 'é')",
			code:    0,
		},
		{
			name:     "SET character_set_client DEFAULT inherits SET GLOBAL value",
			session:  []string{"SET GLOBAL character_set_client = ascii", "SET character_set_client = DEFAULT"},
			cleanup:  []string{"SET GLOBAL character_set_client = DEFAULT", "SET NAMES utf8mb4"},
			stmt:     "CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')",
			code:     0,
			contains: []string{"DEFAULT '??'"},
		},
		{
			name: "SET system_time_zone read-only rejected with 1238",
			stmt: "SET system_time_zone = 'UTC'",
			code: 1238,
		},
		{
			name: "SET PERSIST_ONLY system_time_zone non-persistent read-only rejected with 1238",
			stmt: "SET PERSIST_ONLY system_time_zone = 'UTC'",
			code: 1238,
		},
		{
			name: "SET foreign_key_checks invalid value rejected with 1231",
			stmt: "SET foreign_key_checks = 'bad'",
			code: 1231,
		},
		{
			name:    "SET GLOBAL keyword carries into subsequent assignments without changing session foreign_key_checks",
			session: []string{"SET GLOBAL sql_mode='', foreign_key_checks=0"},
			cleanup: []string{"SET GLOBAL sql_mode = DEFAULT", "SET GLOBAL foreign_key_checks = DEFAULT"},
			stmt:    "CREATE TABLE a (id INT, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES missing_parent(id))",
			code:    1824,
		},
		{
			name:    "SET GLOBAL keyword carries across user variable assignment without changing session foreign_key_checks",
			session: []string{"SET GLOBAL sql_mode='', @u=1, foreign_key_checks=0"},
			cleanup: []string{"SET GLOBAL sql_mode = DEFAULT", "SET GLOBAL foreign_key_checks = DEFAULT"},
			stmt:    "CREATE TABLE a (id INT, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES missing_parent(id))",
			code:    1824,
		},
		{
			name:    "SET @@GLOBAL qualifier does not carry into subsequent assignments",
			session: []string{"SET @@GLOBAL.sql_mode='', foreign_key_checks=0"},
			cleanup: []string{"SET GLOBAL sql_mode = DEFAULT", "SET SESSION foreign_key_checks = DEFAULT"},
			stmt:    "CREATE TABLE a (id INT, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES missing_parent(id))",
			code:    0,
		},
		{
			name:    "CREATE TABLE encrypted in unencrypted general tablespace rejected with 3825",
			setup:   []string{"CREATE TABLESPACE ts_enc_def ADD DATAFILE 'ts_enc_def.ibd'"},
			cleanup: []string{"DROP TABLESPACE ts_enc_def"},
			stmt:    "CREATE TABLE a (id INT) TABLESPACE ts_enc_def ENCRYPTION='Y'",
			code:    3825,
		},
		{
			name:    "ALTER TABLE moving encrypted table into unencrypted general tablespace rejected with 3825",
			setup:   []string{"CREATE TABLESPACE ts_enc_n ADD DATAFILE 'ts_enc_n.ibd' ENCRYPTION='N'", "CREATE TABLE a (id INT)"},
			cleanup: []string{"DROP TABLE IF EXISTS a", "DROP TABLESPACE ts_enc_n"},
			stmt:    "ALTER TABLE a ENCRYPTION='Y', TABLESPACE ts_enc_n",
			code:    3825,
		},
		{
			name: "SET time_zone with unknown zone rejected with 1298",
			stmt: "SET time_zone = 'Invalid/Zone'",
			code: 1298,
		},
		{
			name:    "CREATE TEMPORARY TABLE LIKE general tablespace source rejected with 1478",
			setup:   []string{"CREATE TABLESPACE ts_tmp_like ADD DATAFILE 'ts_tmp_like.ibd'", "CREATE TABLE s (a INT) TABLESPACE ts_tmp_like"},
			cleanup: []string{"DROP TABLE IF EXISTS s", "DROP TABLESPACE ts_tmp_like"},
			stmt:    "CREATE TEMPORARY TABLE a LIKE s",
			code:    1478,
		},
		{
			name:    "CREATE TEMPORARY TABLE LIKE innodb_system source rejected with 1478",
			setup:   []string{"CREATE TABLE s (a INT) TABLESPACE innodb_system"},
			cleanup: []string{"DROP TABLE IF EXISTS s"},
			stmt:    "CREATE TEMPORARY TABLE a LIKE s",
			code:    1478,
		},
		{
			name:    "CREATE TEMPORARY TABLE LIKE innodb_file_per_table source rejected with 1478",
			setup:   []string{"CREATE TABLE s (a INT) TABLESPACE innodb_file_per_table"},
			cleanup: []string{"DROP TABLE IF EXISTS s"},
			stmt:    "CREATE TEMPORARY TABLE a LIKE s",
			code:    1478,
		},
		{
			name:    "CREATE TEMPORARY TABLE LIKE MyISAM source with TABLESPACE clause accepted",
			setup:   []string{"CREATE TABLESPACE ts_myisam_like ADD DATAFILE 'ts_myisam_like.ibd'", "CREATE TABLE s (a INT) ENGINE=MyISAM TABLESPACE ts_myisam_like"},
			cleanup: []string{"DROP TEMPORARY TABLE IF EXISTS a", "DROP TABLE IF EXISTS s", "DROP TABLESPACE ts_myisam_like"},
			stmt:    "CREATE TEMPORARY TABLE a LIKE s",
			code:    0,
		},
		{
			name:  "Repeated DISCARD TABLESPACE accepted with warning on mysql 8.0.46",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY)", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a DISCARD TABLESPACE",
			code:  0,
		},
		{
			name:    "Repeated DISCARD TABLESPACE accepted under TRADITIONAL sql_mode",
			setup:   []string{"CREATE TABLE a (id INT PRIMARY KEY)", "ALTER TABLE a DISCARD TABLESPACE"},
			session: []string{"SET sql_mode = 'TRADITIONAL'"},
			stmt:    "ALTER TABLE a DISCARD TABLESPACE",
			code:    0,
		},
		{
			name:    "SET sql_mode from decimal arithmetic user var rejected with 1232",
			session: []string{"SET @mode = ''", "SET @mode = 1.5 + 1"},
			stmt:    "SET sql_mode = @mode",
			code:    1232,
		},
		{
			name: "Composite TEXT index exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(2000), b(2000)))",
			code: 1071,
		},
		{
			name:    "Composite TEXT index exceeding 3072 bytes rejected with 1071 in non-strict mode",
			session: []string{"SET sql_mode = ''"},
			stmt:    "CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(2000), b(2000)))",
			code:    1071,
		},
		{
			name: "Composite VARCHAR utf8mb4 index exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a VARCHAR(500) CHARACTER SET utf8mb4, b VARCHAR(500) CHARACTER SET utf8mb4, KEY(a, b))",
			code: 1071,
		},
		{
			name: "Composite TEXT index at exact 3072-byte boundary accepted",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(2000), b(1072)))",
			code: 0,
		},
		{
			name:  "ALTER TABLE ADD KEY composite exceeding 3072 bytes rejected with 1071",
			setup: []string{"CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii)"},
			stmt:  "ALTER TABLE a ADD KEY(a(2000), b(2000))",
			code:  1071,
		},
		{
			name: "Composite MyISAM TEXT index exceeding 1000 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(600), b(500))) ENGINE=MyISAM",
			code: 1071,
		},
		{
			name:    "SET SESSION sql_mode REPLACE expression accepted",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET SESSION sql_mode = REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', '')",
			code:    0,
		},
		{
			name: "SET PERSIST_ONLY innodb_page_size 32768 rejected with 1238",
			stmt: "SET PERSIST_ONLY innodb_page_size = 32768",
			code: 1238,
		},
		{
			name: "SET PERSIST_ONLY innodb_page_size invalid value rejected with 1238",
			stmt: "SET PERSIST_ONLY innodb_page_size = 12345",
			code: 1238,
		},
		{
			name: "SET PERSIST_ONLY innodb_directories rejected with 1238",
			stmt: "SET PERSIST_ONLY innodb_directories = '/var/lib/mysql'",
			code: 1238,
		},
		{
			name: "CREATE TABLE ROW_FORMAT=COMPACT key over 767 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(200))) ROW_FORMAT=COMPACT",
			code: 1071,
		},
		{
			name: "CREATE TABLE ROW_FORMAT=REDUNDANT key over 767 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(200))) ROW_FORMAT=REDUNDANT",
			code: 1071,
		},
		{
			name: "CREATE TABLE ROW_FORMAT=COMPACT key within 767 bytes accepted",
			stmt: "CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(191))) ROW_FORMAT=COMPACT",
			code: 0,
		},
		{
			name:  "ALTER TABLE ROW_FORMAT=COMPACT revalidates existing index key limit with 1071",
			setup: []string{"CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(200))) ROW_FORMAT=DYNAMIC"},
			stmt:  "ALTER TABLE a ROW_FORMAT=COMPACT",
			code:  1071,
		},
		{
			name:  "ALTER TABLE ADD KEY on ROW_FORMAT=COMPACT table over 767 bytes rejected with 1071",
			setup: []string{"CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(191))) ROW_FORMAT=COMPACT"},
			stmt:  "ALTER TABLE a ADD KEY k2 (a(200))",
			code:  1071,
		},
		{
			name: "CREATE TABLE innodb_system ROW_FORMAT=COMPRESSED rejected with 1478",
			stmt: "CREATE TABLE a (id INT) TABLESPACE innodb_system ROW_FORMAT=COMPRESSED",
			code: 1478,
		},
		{
			name: "SET sql_mode missing schema-qualified function rejected with 1305",
			stmt: "SET sql_mode = missing_fn_db.CONCAT('ANSI')",
			code: 1305,
		},
		{
			name: "SET user var missing schema-qualified function rejected with 1305",
			stmt: "SET @u_qfn = missing_fn_db.CONCAT('ANSI')",
			code: 1305,
		},
		{
			name:    "Non-strict CREATE TABLE UNIQUE TEXT prefix over limit rejected with 1071",
			session: []string{"SET sql_mode = ''"},
			stmt:    "CREATE TABLE a (a TEXT CHARACTER SET utf8mb4, UNIQUE KEY(a(1000)))",
			code:    1071,
		},
		{
			name:    "Non-strict CREATE TABLE non-unique TEXT prefix clamped to 0",
			session: []string{"SET sql_mode = ''"},
			stmt:    "CREATE TABLE a (a TEXT CHARACTER SET utf8mb4, KEY(a(1000)))",
			code:    0,
		},
		{
			name:    "Non-strict ALTER TABLE ADD UNIQUE KEY TEXT prefix over limit rejected with 1071",
			setup:   []string{"CREATE TABLE a (a TEXT CHARACTER SET utf8mb4)"},
			session: []string{"SET sql_mode = ''"},
			stmt:    "ALTER TABLE a ADD UNIQUE KEY (a(1000))",
			code:    1071,
		},
		{
			name:    "Non-strict CREATE UNIQUE INDEX TEXT prefix over limit rejected with 1071",
			setup:   []string{"CREATE TABLE a (a TEXT CHARACTER SET utf8mb4)"},
			session: []string{"SET sql_mode = ''"},
			stmt:    "CREATE UNIQUE INDEX idx_uq ON a (a(1000))",
			code:    1071,
		},
		{
			name: "Mixed composite TEXT+BIGINT exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, n BIGINT, KEY(a(3072), n))",
			code: 1071,
		},
		{
			name: "Mixed composite TEXT+BIGINT at exact 3072 bytes accepted",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, n BIGINT, KEY(a(3064), n))",
			code: 0,
		},
		{
			name: "Mixed composite TEXT+DATETIME(6) exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, d DATETIME(6), KEY(a(3065), d))",
			code: 1071,
		},
		{
			name: "Mixed composite TEXT+DECIMAL(10,2) exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, d DECIMAL(10,2), KEY(a(3068), d))",
			code: 1071,
		},
		{
			name: "Mixed composite TEXT+BIT(16) exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, b BIT(16), KEY(a(3071), b))",
			code: 1071,
		},
		{
			name: "Mixed composite TEXT+ENUM exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT CHARACTER SET ascii, e ENUM('x','y'), KEY(a(3072), e))",
			code: 1071,
		},
		{
			name: "SET sql_mode IFNULL wrong parameter count rejected with 1582",
			stmt: "SET sql_mode = IFNULL('', 'ANSI', 'TRADITIONAL')",
			code: 1582,
		},
		{
			name: "SET sql_mode empty CONCAT rejected with 1582",
			stmt: "SET sql_mode = CONCAT()",
			code: 1582,
		},
		{
			name:  "ALTER TABLE RENAME COLUMN referenced by surviving functional index rejected with 3837",
			setup: []string{"CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))"},
			stmt:  "ALTER TABLE a RENAME COLUMN a TO b",
			code:  3837,
		},
		{
			name:  "ALTER TABLE DROP COLUMN referenced by surviving functional index rejected with 3837",
			setup: []string{"CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))"},
			stmt:  "ALTER TABLE a DROP COLUMN a",
			code:  3837,
		},
		{
			name:  "ALTER TABLE CHANGE COLUMN renaming functional index column rejected with 3837",
			setup: []string{"CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))"},
			stmt:  "ALTER TABLE a CHANGE COLUMN a b INT",
			code:  3837,
		},
		{
			name:  "ALTER TABLE RENAME COLUMN accepted when same ALTER drops dependent functional index",
			setup: []string{"CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))"},
			stmt:  "ALTER TABLE a RENAME COLUMN a TO b, DROP INDEX idx",
			code:  0,
		},
		{
			name:  "ALTER TABLE MODIFY COLUMN preserving functional index column name accepted",
			setup: []string{"CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))"},
			stmt:  "ALTER TABLE a MODIFY COLUMN a VARCHAR(20)",
			code:  0,
		},
		{
			name:    "SET time_zone numeric +014:00 normalized per mysql str_to_time_zone",
			cleanup: []string{"SET time_zone = DEFAULT"},
			stmt:    "SET time_zone = '+014:00'",
			code:    0,
		},
		{
			name:    "SET time_zone numeric +1:000 normalized per mysql str_to_time_zone",
			cleanup: []string{"SET time_zone = DEFAULT"},
			stmt:    "SET time_zone = '+1:000'",
			code:    0,
		},
		{
			name: "SET time_zone ++01:00 rejected with 1298",
			stmt: "SET time_zone = '++01:00'",
			code: 1298,
		},
		{
			name: "CREATE TABLE DEFAULT (IFNULL(1)) wrong arity rejected with 1582",
			stmt: "CREATE TABLE a (a INT DEFAULT (IFNULL(1)))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (IFNULL(1, 2, 3)) wrong arity rejected with 1582",
			stmt: "CREATE TABLE a (a INT DEFAULT (IFNULL(1, 2, 3)))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (CONCAT()) wrong arity rejected with 1582",
			stmt: "CREATE TABLE a (a VARCHAR(10) DEFAULT (CONCAT()))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (ABS()) wrong arity rejected with 1582",
			stmt: "CREATE TABLE a (a INT DEFAULT (ABS()))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (ABS(1, 2)) wrong arity rejected with 1582",
			stmt: "CREATE TABLE a (a INT DEFAULT (ABS(1, 2)))",
			code: 1582,
		},
		{
			name:  "ALTER TABLE ADD COLUMN DEFAULT (IFNULL(1)) wrong arity rejected with 1582",
			setup: []string{"CREATE TABLE a (id INT)"},
			stmt:  "ALTER TABLE a ADD COLUMN b INT DEFAULT (IFNULL(1))",
			code:  1582,
		},
		{
			name:  "ALTER TABLE SET DEFAULT (IFNULL(1)) wrong arity rejected with 1582",
			setup: []string{"CREATE TABLE a (id INT)"},
			stmt:  "ALTER TABLE a ALTER COLUMN id SET DEFAULT (IFNULL(1))",
			code:  1582,
		},
		{
			name: "CREATE TABLE DEFAULT (IFNULL(1, 2)) valid arity accepted",
			stmt: "CREATE TABLE a (a INT DEFAULT (IFNULL(1, 2)))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_OBJECT('k')) odd arg count rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT('k')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_OBJECT('k', 1, 'k2')) odd arg count rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT('k', 1, 'k2')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_OBJECT()) empty accepted",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT()))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_OBJECT('k', 1)) pair accepted",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT('k', 1)))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_SET('{}', '$.a')) missing value arg rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_SET('{}', '$.a')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_SET('{}', '$.a', 1)) valid pair accepted",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_SET('{}', '$.a', 1)))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_INSERT('{}', '$.a')) missing value arg rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_INSERT('{}', '$.a')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_REPLACE('{}', '$.a')) missing value arg rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_REPLACE('{}', '$.a')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_ARRAY_APPEND('[]', '$')) missing value arg rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_ARRAY_APPEND('[]', '$')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (JSON_ARRAY_INSERT('[]', '$[0]')) missing value arg rejected with 1582",
			stmt: "CREATE TABLE a (j JSON DEFAULT (JSON_ARRAY_INSERT('[]', '$[0]')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE functional composite index exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a VARCHAR(400), b VARCHAR(400), KEY ((CONCAT(a, '')), (CONCAT(b, ''))))",
			code: 1071,
		},
		{
			name: "CREATE TABLE functional composite unique index exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a VARCHAR(400), b VARCHAR(400), UNIQUE KEY ((CONCAT(a, '')), (CONCAT(b, ''))))",
			code: 1071,
		},
		{
			name:  "ALTER TABLE ADD functional composite index exceeding 3072 bytes rejected with 1071",
			setup: []string{"CREATE TABLE a (a VARCHAR(400), b VARCHAR(400))"},
			stmt:  "ALTER TABLE a ADD KEY ((CONCAT(a, '')), (CONCAT(b, '')))",
			code:  1071,
		},
		{
			name:  "ALTER TABLE ADD functional composite unique index exceeding 3072 bytes rejected with 1071",
			setup: []string{"CREATE TABLE a (a VARCHAR(400), b VARCHAR(400))"},
			stmt:  "ALTER TABLE a ADD UNIQUE KEY ((CONCAT(a, '')), (CONCAT(b, '')))",
			code:  1071,
		},
		{
			name:  "ALTER TABLE ADD functional key before ADD columns exceeding 3072 bytes rejected with 1071",
			setup: []string{"CREATE TABLE a (id INT)"},
			stmt:  "ALTER TABLE a ADD KEY ((CONCAT(a, '')), (CONCAT(b, ''))), ADD a VARCHAR(400), ADD b VARCHAR(400)",
			code:  1071,
		},
		{
			name:  "CREATE INDEX functional composite exceeding 3072 bytes rejected with 1071",
			setup: []string{"CREATE TABLE a (a VARCHAR(400), b VARCHAR(400))"},
			stmt:  "CREATE INDEX idx_c ON a ((CONCAT(a, '')), (CONCAT(b, '')))",
			code:  1071,
		},
		{
			name: "CREATE TABLE DEFAULT (ST_X()) too few args rejected with 1582",
			stmt: "CREATE TABLE a (x DOUBLE DEFAULT (ST_X()))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (ST_X(POINT(1, 2))) accepted with 0",
			stmt: "CREATE TABLE a (x DOUBLE DEFAULT (ST_X(POINT(1, 2))))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (ST_X(POINT(1, 2), 3)) 2-arg setter accepted with 0",
			stmt: "CREATE TABLE a (x DOUBLE DEFAULT (ST_X(POINT(1, 2), 3)))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (MBRWITHIN(POINT(1, 2))) too few args rejected with 1582",
			stmt: "CREATE TABLE a (x INT DEFAULT (MBRWITHIN(POINT(1, 2))))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (ROLES_GRAPHML('x')) too many args rejected with 1582",
			stmt: "CREATE TABLE a (x VARCHAR(100) DEFAULT (ROLES_GRAPHML('x')))",
			code: 1582,
		},
		{
			name: "CREATE TABLE DEFAULT (ST_CONTAINS(POINT(1, 2))) too few args rejected with 1582",
			stmt: "CREATE TABLE a (x INT DEFAULT (ST_CONTAINS(POINT(1, 2))))",
			code: 1582,
		},
		{
			name: "CREATE TABLE functional SUBSTRING index exceeding 3072 bytes rejected with 1071",
			stmt: "CREATE TABLE a (a TEXT, b TEXT, KEY ((SUBSTRING(a, 1, 400)), (SUBSTRING(b, 1, 400))))",
			code: 1071,
		},
		{
			name: "CREATE TABLE DEFAULT (TRIM(LEADING 'x' FROM 'xxx')) accepted with 0",
			stmt: "CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM(LEADING 'x' FROM 'xxx')))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (TRIM(TRAILING 'x' FROM 'xxx')) accepted with 0",
			stmt: "CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM(TRAILING 'x' FROM 'xxx')))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (TRIM(BOTH 'x' FROM 'xxx')) accepted with 0",
			stmt: "CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM(BOTH 'x' FROM 'xxx')))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (TRIM('x' FROM 'xxx')) accepted with 0",
			stmt: "CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM('x' FROM 'xxx')))",
			code: 0,
		},
		{
			name:  "ALTER TABLE ADD KEY after DISCARD TABLESPACE accepted with 0 on mysql 8.0.46",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, v INT) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD KEY k (v)",
			code:  0,
		},
		{
			name:  "ALTER TABLE ADD KEY ALGORITHM=COPY after DISCARD TABLESPACE rejected with 1814",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, v INT) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD KEY k (v), ALGORITHM=COPY",
			code:  1814,
		},
		{
			name:  "ALTER TABLE ADD PRIMARY KEY after DISCARD TABLESPACE rejected with 1814",
			setup: []string{"CREATE TABLE a (id INT, v INT) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a ADD PRIMARY KEY (id)",
			code:  1814,
		},
		{
			name:  "ALTER TABLE FORCE after DISCARD TABLESPACE rejected with 1814",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "ALTER TABLE a FORCE",
			code:  1814,
		},
		{
			name:  "CREATE INDEX after DISCARD TABLESPACE accepted with 0 on mysql 8.0.46",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, v INT) ENGINE=InnoDB", "ALTER TABLE a DISCARD TABLESPACE"},
			stmt:  "CREATE INDEX idx_v ON a (v)",
			code:  0,
		},
		{
			name:    "SET sql_mode = IF(1, 'ANSI', '') enables ANSI mode on both server and catalog",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IF(1, 'ANSI', '')",
			code:    0,
		},
		{
			name: "CREATE TABLE DEFAULT (REPLACE('abc', 'a', 'd')) accepted with 0",
			stmt: "CREATE TABLE a (id INT, s VARCHAR(16) DEFAULT (REPLACE('abc', 'a', 'd')))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (ATAN2(1)) 1-arg form accepted with 0",
			stmt: "CREATE TABLE a (d DOUBLE DEFAULT (ATAN2(1)))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (AES_ENCRYPT 6-arg HKDF form) accepted with 0",
			stmt: "CREATE TABLE a (s VARBINARY(64) DEFAULT (AES_ENCRYPT('a', 'b', 'c', 'd', 'e', 'f')))",
			code: 0,
		},
		{
			name: "CREATE TABLE DEFAULT (AES_DECRYPT 6-arg HKDF form) accepted with 0",
			stmt: "CREATE TABLE a (s VARCHAR(64) DEFAULT (AES_DECRYPT('a', 'b', 'c', 'd', 'e', 'f')))",
			code: 0,
		},
		{
			name:  "ALTER TABLE MODIFY COLUMN widening existing indexed column past 3072 rejected with 1071",
			setup: []string{"CREATE TABLE a (v VARCHAR(20) CHARACTER SET utf8mb4, KEY k(v))"},
			stmt:  "ALTER TABLE a MODIFY COLUMN v VARCHAR(1000) CHARACTER SET utf8mb4",
			code:  1071,
		},
		{
			name:  "ALTER TABLE CONVERT TO CHARSET utf8mb4 widening indexed latin1 column past 3072 rejected with 1071",
			setup: []string{"CREATE TABLE a (v VARCHAR(800)) DEFAULT CHARSET=latin1, ENGINE=InnoDB", "ALTER TABLE a ADD KEY k(v)"},
			stmt:  "ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4",
			code:  1071,
		},
		{
			name: "CREATE TABLE with table-qualified DEFAULT and functional index accepted",
			stmt: "CREATE TABLE a (x INT, y INT DEFAULT (a.x + 1), KEY ((a.x + 1)))",
			code: 0,
		},
		{
			name:     "ALTER TABLE CONVERT TO CHARSET utf8mb4 promotes TEXT to MEDIUMTEXT",
			setup:    []string{"CREATE TABLE a (c1 TINYTEXT, c2 TEXT, c3 MEDIUMTEXT) DEFAULT CHARSET=latin1"},
			stmt:     "ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4",
			code:     0,
			contains: []string{"`c1` text", "`c2` mediumtext", "`c3` longtext"},
		},
		{
			name: "CREATE TABLE COMPACT composite key exceeding 767 total but <=767 per part accepted",
			stmt: "CREATE TABLE a (x VARBINARY(600), y VARBINARY(600), KEY k(x, y)) ROW_FORMAT=COMPACT",
			code: 0,
		},
		{
			name:  "ALTER TABLE ADD KEY with ROW_FORMAT=DYNAMIC on COMPACT table accepted",
			setup: []string{"CREATE TABLE a (v VARCHAR(200) CHARACTER SET utf8mb4) ROW_FORMAT=COMPACT"},
			stmt:  "ALTER TABLE a ADD KEY k(v), ROW_FORMAT=DYNAMIC",
			code:  0,
		},
		{
			name: "CREATE TABLE inline JSON UNIQUE rejected with 3152",
			stmt: "CREATE TABLE a (j JSON UNIQUE)",
			code: 3152,
		},
		{
			name: "CREATE TABLE inline oversized VARCHAR PRIMARY KEY rejected with 1071",
			stmt: "CREATE TABLE a (v VARCHAR(900) CHARACTER SET utf8mb4 PRIMARY KEY)",
			code: 1071,
		},
		{
			name:     "Non-strict sql_mode shortens single-column non-unique index prefix to 768",
			session:  []string{"SET sql_mode = ''"},
			cleanup:  []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:     "CREATE TABLE a (v VARCHAR(900) CHARACTER SET utf8mb4, KEY i(v))",
			code:     0,
			contains: []string{"KEY `i` (`v`(768))"},
		},
		{
			name: "CONCAT functional index uses aggregated utf8mb4 charset and rejects >3072 with 1071",
			stmt: "CREATE TABLE a (a VARCHAR(1) CHARACTER SET latin1, b VARCHAR(510) CHARACTER SET utf8mb4, c VARBINARY(1050), KEY ((CONCAT(a, b)), c))",
			code: 1071,
		},
		{
			name:     "Approximate numeric literal DEFAULT 2.5E0 uses round-half-to-even while exact 2.5 rounds away from zero",
			stmt:     "CREATE TABLE a (i INT DEFAULT 2.5E0, j INT DEFAULT 2.5, k TINYINT DEFAULT -128.5E0)",
			code:     0,
			contains: []string{"DEFAULT '2'", "DEFAULT '3'", "DEFAULT '-128'"},
		},
		{
			name:    "ENUM latin1_general_ci distinguishes e and e-acute under STRICT_ALL_TABLES",
			session: []string{"SET sql_mode = 'STRICT_ALL_TABLES'"},
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "CREATE TABLE a (c ENUM('e', 'é') COLLATE latin1_general_ci)",
			code:    0,
		},
		{
			name:    "ENUM utf8mb4_unicode_520_ci distinguishes supplementary emoji under STRICT_ALL_TABLES",
			session: []string{"SET sql_mode = 'STRICT_ALL_TABLES'"},
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "CREATE TABLE a (c ENUM('😀', '😁') COLLATE utf8mb4_unicode_520_ci)",
			code:    0,
		},
		{
			name:  "IMPORT PARTITION on non-discarded p1 rejected with 1813 when p0 is discarded",
			setup: []string{"CREATE TABLE a (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))", "ALTER TABLE a DISCARD PARTITION p0 TABLESPACE"},
			stmt:  "ALTER TABLE a IMPORT PARTITION p1 TABLESPACE",
			code:  1813,
		},
		{
			name:  "TRUNCATE PARTITION on discarded p0 rejected with 1814",
			setup: []string{"CREATE TABLE a (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))", "ALTER TABLE a DISCARD PARTITION p0 TABLESPACE"},
			stmt:  "ALTER TABLE a TRUNCATE PARTITION p0",
			code:  1814,
		},
		{
			name:  "TRUNCATE PARTITION on non-discarded p1 accepted when p0 is discarded",
			setup: []string{"CREATE TABLE a (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))", "ALTER TABLE a DISCARD PARTITION p0 TABLESPACE"},
			stmt:  "ALTER TABLE a TRUNCATE PARTITION p1",
			code:  0,
		},
		{
			name:    "SET sql_mode CASE uses session collation case-insensitivity",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = CASE 'ansi' WHEN 'ANSI' THEN 'TRADITIONAL' ELSE '' END",
			code:    0,
		},
		{
			name:    "SET mixed numeric string comparison coerces nonnumeric string to 0",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IF(0 = 'abc', 'ANSI', '')",
			code:    0,
		},
		{
			name:    "SET one-arg TRIM preserves non-space ASCII tab X'09'",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IF(TRIM(X'09') = X'', 'ANSI', '')",
			code:    0,
		},
		{
			name:    "SET CHAR_LENGTH counts bytes for binary hex literal X'C3A9'",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IF(CHAR_LENGTH(X'C3A9') = 2, 'ANSI', '')",
			code:    0,
		},
		{
			name:    "SET IF rejects missing qualified function in unselected false branch with 1305",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IF(1, 'ANSI', sim_oracle_missing.no_such())",
			code:    1305,
		},
		{
			name:    "SET sql_mode = CAST('ANSI' AS CHAR) enables ANSI mode on both server and catalog",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = CAST('ANSI' AS CHAR)",
			code:    0,
		},
		{
			name:    "SET IFNULL validates skipped second argument for missing function with 1305",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IFNULL('ANSI', sim_oracle_missing.no_such())",
			code:    1305,
		},
		{
			name:    "SET IF mixed numeric and string branch aggregates to string before comparison",
			cleanup: []string{"SET SESSION sql_mode = DEFAULT"},
			stmt:    "SET sql_mode = IF(IF(1, 1, 'x') = '01', 'ANSI', '')",
			code:    0,
		},
		{
			name:  "DISCARD TABLESPACE on self-referencing FK table with foreign_key_checks=1 accepted with 0 on mysql 8.0.46",
			setup: []string{"CREATE TABLE a (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_self FOREIGN KEY (pid) REFERENCES a(id))"},
			stmt:  "ALTER TABLE a DISCARD TABLESPACE",
			code:  0,
		},
		{
			// ATAN2 takes one argument as ATAN does: mysql:8.0.46 stores it.
			name: "DEFAULT (ATAN2(1))",
			stmt: "CREATE TABLE z (x DOUBLE DEFAULT (ATAN2(1)))",
			code: 0,
		},
		// A hex or bit literal in a functional CONCAT is a binary string of its bytes.
		{name: "functional CONCAT with a hex literal", stmt: "CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, X'00'))))", code: 0},
		{name: "functional CONCAT with a bit literal", stmt: "CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, b'0000000100000001'))))", code: 0},
		{name: "functional CONCAT of a character column and a hex literal", stmt: "CREATE TABLE z (v VARCHAR(10), KEY ((CONCAT(v, X'4142'))))", code: 0},
		{name: "functional CONCAT with a long binary REPEAT", stmt: "CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, REPEAT(X'00', 600)))))", code: 3757},
		// Built-in tablespace names are case-sensitive, and an innodb_ prefix
		// is reserved only for CREATE.
		{name: "ALTER UNDO TABLESPACE MYSQL", stmt: "ALTER UNDO TABLESPACE MYSQL SET INACTIVE", code: 3510},
		{name: "DROP UNDO TABLESPACE MYSQL", stmt: "DROP UNDO TABLESPACE MYSQL", code: 3510},
		{name: "DROP UNDO TABLESPACE INNODB_SYSTEM", stmt: "DROP UNDO TABLESPACE INNODB_SYSTEM", code: 3510},
		{name: "ALTER UNDO TABLESPACE mysql", stmt: "ALTER UNDO TABLESPACE mysql SET INACTIVE", code: 3119},
		{name: "DROP UNDO TABLESPACE mysql", stmt: "DROP UNDO TABLESPACE mysql", code: 3119},
		{name: "DROP UNDO TABLESPACE innodb_system", stmt: "DROP UNDO TABLESPACE innodb_system", code: 3119},
		{name: "DROP UNDO TABLESPACE innodb_temporary", stmt: "DROP UNDO TABLESPACE innodb_temporary", code: 3119},
		{name: "DROP UNDO TABLESPACE with an innodb_ prefix", stmt: "DROP UNDO TABLESPACE innodb_x", code: 3510},
		{name: "ALTER UNDO TABLESPACE with an innodb_ prefix", stmt: "ALTER UNDO TABLESPACE INNODB_X SET INACTIVE", code: 3510},
		{name: "temporary table with an unkeyed AUTO_INCREMENT column", stmt: "CREATE TEMPORARY TABLE z (a INT AUTO_INCREMENT)", code: 1075},
		// REPEAT keeps a binary argument binary, and a huge count is a LOB.
		{name: "functional CONCAT with a binary REPEAT", stmt: "CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, REPEAT(X'00', 128)))))", code: 0},
		{name: "functional CONCAT with a long binary REPEAT of 512", stmt: "CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, REPEAT(X'00', 512)))))", code: 3757},
		{name: "functional REPEAT of a binary column within 512 bytes", stmt: "CREATE TABLE z (b VARBINARY(10), KEY ((REPEAT(b, 51))))", code: 0},
		{name: "functional REPEAT of a binary column past 512 bytes", stmt: "CREATE TABLE z (b VARBINARY(10), KEY ((REPEAT(b, 52))))", code: 3757},
		{name: "functional REPEAT with a huge count", stmt: "CREATE TABLE z (v VARCHAR(2), KEY ((REPEAT(v, 9223372036854775807))))", code: 3757},
		// A functional index's hidden column resolves in no expression.
		{
			name:  "ALTER ADD functional KEY over another's hidden column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY i ((c + 1)))"},
			stmt:  "ALTER TABLE a ADD KEY j ((`!hidden!i!0!0` + 1))",
			code:  1054,
		},
		{
			name:  "CREATE INDEX over another's hidden column",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY i ((c + 1)))"},
			stmt:  "CREATE INDEX j ON a ((`!hidden!i!0!0` + 1))",
			code:  1054,
		},
		// A functional key part meets the 767-byte per-part limit of a COMPACT
		// or REDUNDANT table in strict mode.
		{
			name:  "CREATE INDEX functional over 767 bytes on a COMPACT table",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(300)) ROW_FORMAT=COMPACT"},
			stmt:  "CREATE INDEX i ON a ((LOWER(s)))",
			code:  1071,
		},
		{
			name:  "ALTER ADD functional KEY over 767 bytes on a REDUNDANT table",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(300)) ROW_FORMAT=REDUNDANT"},
			stmt:  "ALTER TABLE a ADD KEY ((LOWER(s)))",
			code:  1071,
		},
		{
			name:  "CREATE INDEX functional over 767 bytes on a DYNAMIC table",
			setup: []string{"CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(300)) ROW_FORMAT=DYNAMIC"},
			stmt:  "CREATE INDEX i ON a ((LOWER(s)))",
			code:  0,
		},
		{name: "CREATE TABLE functional key part of 764 bytes on a COMPACT table", stmt: "CREATE TABLE z (s VARCHAR(191), KEY ((LOWER(s)))) ROW_FORMAT=COMPACT", code: 0},
		{name: "CREATE TABLE functional key part of 768 bytes on a COMPACT table", stmt: "CREATE TABLE z (s VARCHAR(192), KEY ((LOWER(s)))) ROW_FORMAT=COMPACT", code: 1071},
		{
			name:    "CREATE TABLE functional key part over 767 bytes on a COMPACT table under sql_mode ''",
			session: []string{"SET SESSION sql_mode = ''"},
			stmt:    "CREATE TABLE z (s VARCHAR(300), KEY ((LOWER(s)))) ROW_FORMAT=COMPACT",
			code:    0,
		},
		// On a COMPACT table, a binary LPAD/RPAD keeps its byte width and a
		// comparison or predicate is an integer, within the 767-byte part limit.
		{name: "COMPACT functional key part: KEY ((RPAD(b, 512, X'00')))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((RPAD(b, 512, X'00')))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((LPAD(b, 512, X'00')))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((LPAD(b, 512, X'00')))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((RPAD(b, 513, X'00')))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((RPAD(b, 513, X'00')))) ROW_FORMAT=COMPACT", code: 3757},
		{name: "COMPACT functional key part: KEY ((a = 1))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a = 1))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((a < 1 AND a > 0))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a < 1 AND a > 0))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((NOT a))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((NOT a))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((a IS NULL))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a IS NULL))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((a IN (1, 2)))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a IN (1, 2)))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((a BETWEEN 1 AND 2))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a BETWEEN 1 AND 2))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((s LIKE 'x%'))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((s LIKE 'x%'))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part: KEY ((a <=> 1))", stmt: "CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a <=> 1))) ROW_FORMAT=COMPACT", code: 0},
		// A unary minus keeps a number's type and makes a string a double.
		{name: "COMPACT functional unary key part: a INT, KEY ((-a))", stmt: "CREATE TABLE z (a INT, KEY ((-a))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional unary key part: s VARCHAR(10), KEY ((-s))", stmt: "CREATE TABLE z (s VARCHAR(10), KEY ((-s))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional unary key part: a TINYINT, KEY ((-a))", stmt: "CREATE TABLE z (a TINYINT, KEY ((-a))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional unary key part: a INT, KEY ((~a))", stmt: "CREATE TABLE z (a INT, KEY ((~a))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional unary key part: a INT, KEY ((+a))", stmt: "CREATE TABLE z (a INT, KEY ((+a))) ROW_FORMAT=COMPACT", code: 3762},
		{name: "COMPACT functional unary key part: d DECIMAL(10,2), KEY ((-d))", stmt: "CREATE TABLE z (d DECIMAL(10,2), KEY ((-d))) ROW_FORMAT=COMPACT", code: 0},
		// Literals, REPLACE and functions the catalog cannot size, on a COMPACT
		// table: a known width meets the 767-byte part limit, an unknown one is
		// not judged on a placeholder.
		{name: "COMPACT functional key part width: a INT, KEY ((1))", stmt: "CREATE TABLE z (a INT, KEY ((1))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: a INT, KEY ((1.5))", stmt: "CREATE TABLE z (a INT, KEY ((1.5))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: a INT, KEY ((1e0))", stmt: "CREATE TABLE z (a INT, KEY ((1e0))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: a INT, KEY (('abc'))", stmt: "CREATE TABLE z (a INT, KEY (('abc'))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: a INT, KEY ((TRUE))", stmt: "CREATE TABLE z (a INT, KEY ((TRUE))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: a INT, KEY ((DATE '2024-01-01'))", stmt: "CREATE TABLE z (a INT, KEY ((DATE '2024-01-01'))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: a INT, KEY ((REPEAT('x', 200)))", stmt: "CREATE TABLE z (a INT, KEY ((REPEAT('x', 200)))) ROW_FORMAT=COMPACT", code: 1071},
		{name: "COMPACT functional key part width: s VARCHAR(100), KEY ((REPLACE(s, 'x', 'y')))", stmt: "CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, 'x', 'y')))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: s VARCHAR(100), KEY ((REPLACE(s, 'x', 'yy')))", stmt: "CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, 'x', 'yy')))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: s VARCHAR(100), KEY ((REPLACE(s, 'xx', 'yyy')))", stmt: "CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, 'xx', 'yyy')))) ROW_FORMAT=COMPACT", code: 1071},
		{name: "COMPACT functional key part width: s VARCHAR(191), KEY ((REPLACE(s, 'x', '')))", stmt: "CREATE TABLE z (s VARCHAR(191), KEY ((REPLACE(s, 'x', '')))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: s VARCHAR(192), KEY ((REPLACE(s, 'x', '')))", stmt: "CREATE TABLE z (s VARCHAR(192), KEY ((REPLACE(s, 'x', '')))) ROW_FORMAT=COMPACT", code: 1071},
		{name: "COMPACT functional key part width: s VARCHAR(100), KEY ((REPLACE(s, '', 'yyyyy')))", stmt: "CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, '', 'yyyyy')))) ROW_FORMAT=COMPACT", code: 1071},
		{name: "COMPACT functional key part width: s VARCHAR(10), KEY ((SHA1(s)))", stmt: "CREATE TABLE z (s VARCHAR(10), KEY ((SHA1(s)))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: s VARCHAR(10), KEY ((MD5(s)))", stmt: "CREATE TABLE z (s VARCHAR(10), KEY ((MD5(s)))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional key part width: s VARCHAR(10), KEY ((HEX(s)))", stmt: "CREATE TABLE z (s VARCHAR(10), KEY ((HEX(s)))) ROW_FORMAT=COMPACT", code: 0},
		// An explicit zero key-part length is refused as the grammar reduces it
		// (1391), before the table or column is looked up.
		{name: "CREATE TABLE key part s(0)", stmt: "CREATE TABLE z (s VARCHAR(10), KEY (s(0)))", code: 1391},
		{name: "CREATE TABLE key part s(0) over a missing column", stmt: "CREATE TABLE z (s VARCHAR(10), KEY (nosuch(0)))", code: 1391},
		{name: "CREATE TABLE PRIMARY KEY part s(0)", stmt: "CREATE TABLE z (s VARCHAR(10), PRIMARY KEY (s(0)))", code: 1391},
		{name: "ALTER ADD KEY part c(0)", setup: []string{oneTable}, stmt: "ALTER TABLE a ADD KEY (c(0))", code: 1391},
		{name: "CREATE INDEX part c(0)", setup: []string{oneTable}, stmt: "CREATE INDEX i ON a (c(0))", code: 1391},
		{name: "CREATE INDEX part s(0) on a missing table", stmt: "CREATE INDEX i ON nosuch (s(0))", code: 1391},
		{name: "CREATE TABLE key part longer than the column", stmt: "CREATE TABLE z (s VARCHAR(10), KEY (s(11)))", code: 1089},
		// SPACE(n) is n characters; SPACE of a non-constant is a LOB.
		{name: "functional SPACE(512)", stmt: "CREATE TABLE z (a INT, KEY ((SPACE(512))))", code: 0},
		{name: "functional SPACE(513)", stmt: "CREATE TABLE z (a INT, KEY ((SPACE(513))))", code: 3757},
		{name: "functional SPACE of a huge count", stmt: "CREATE TABLE z (a INT, KEY ((SPACE(9223372036854775807))))", code: 3757},
		{name: "functional SPACE of a column", stmt: "CREATE TABLE z (a INT, KEY ((SPACE(a))))", code: 3757},
		{name: "COMPACT functional SPACE(191)", stmt: "CREATE TABLE z (a INT, KEY ((SPACE(191)))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional SPACE(192)", stmt: "CREATE TABLE z (a INT, KEY ((SPACE(192)))) ROW_FORMAT=COMPACT", code: 1071},
		// Integer-returning functions are small key parts, not placeholders.
		{name: "functional key of 13 ASCII parts", stmt: "CREATE TABLE z (c1 VARCHAR(10), c2 VARCHAR(10), c3 VARCHAR(10), c4 VARCHAR(10), c5 VARCHAR(10), c6 VARCHAR(10), c7 VARCHAR(10), c8 VARCHAR(10), c9 VARCHAR(10), c10 VARCHAR(10), c11 VARCHAR(10), c12 VARCHAR(10), c13 VARCHAR(10), KEY ((ASCII(c1)), (ASCII(c2)), (ASCII(c3)), (ASCII(c4)), (ASCII(c5)), (ASCII(c6)), (ASCII(c7)), (ASCII(c8)), (ASCII(c9)), (ASCII(c10)), (ASCII(c11)), (ASCII(c12)), (ASCII(c13))))", code: 0},
		{name: "COMPACT functional ASCII part", stmt: "CREATE TABLE z (s VARCHAR(10), KEY ((ASCII(s)))) ROW_FORMAT=COMPACT", code: 0},
		// A binary CONCAT_WS counts a character separator at 4 bytes.
		{name: "functional CONCAT_WS of binary columns at 512 bytes", stmt: "CREATE TABLE z (b VARBINARY(254), KEY ((CONCAT_WS('x', b, b))))", code: 0},
		{name: "functional CONCAT_WS of binary columns past 512 bytes", stmt: "CREATE TABLE z (b VARBINARY(255), KEY ((CONCAT_WS('x', b, b))))", code: 3757},
		{name: "functional CONCAT_WS with a binary separator", stmt: "CREATE TABLE z (b VARBINARY(10), KEY ((CONCAT_WS(X'00', b, b, b))))", code: 0},
		// IGNORE_SPACE makes a spaced call a built-in in a functional key part.
		{name: "functional spaced SUBSTRING under IGNORE_SPACE", session: []string{"SET SESSION sql_mode = 'IGNORE_SPACE'"}, stmt: "CREATE TABLE z (a VARCHAR(10), KEY ((SUBSTRING (a,1,1))))", code: 0},
		{name: "functional spaced SUBSTRING under ANSI", session: []string{"SET SESSION sql_mode = 'ANSI'"}, stmt: "CREATE TABLE z (a VARCHAR(10), KEY ((SUBSTRING (a,1,1))))", code: 0},
		{name: "functional spaced SUBSTRING without IGNORE_SPACE", stmt: "CREATE TABLE z (a VARCHAR(10), KEY ((SUBSTRING (a,1,1))))", code: 3758},
		{name: "functional spaced unknown function under IGNORE_SPACE", session: []string{"SET SESSION sql_mode = 'IGNORE_SPACE'"}, stmt: "CREATE TABLE z (a VARCHAR(10), KEY ((nosuch (a))))", code: 3758},
		{name: "functional spaced NOW under IGNORE_SPACE", session: []string{"SET SESSION sql_mode = 'IGNORE_SPACE'"}, stmt: "CREATE TABLE z (a VARCHAR(10), KEY ((NOW ())))", code: 3758},
		{name: "CREATE INDEX spaced SUBSTRING under IGNORE_SPACE", setup: []string{oneTable}, session: []string{"SET SESSION sql_mode = 'IGNORE_SPACE'"}, stmt: "CREATE INDEX i ON a ((SUBSTRING (c,1,1)))", code: 0},
		// A literal's own charset sizes its key part: N'' is utf8mb3.
		{name: "COMPACT functional literal N of 255 characters", stmt: "CREATE TABLE z (a INT, KEY ((N'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional literal N of 256 characters", stmt: "CREATE TABLE z (a INT, KEY ((N'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'))) ROW_FORMAT=COMPACT", code: 1071},
		{name: "COMPACT functional literal _latin1 of 256 characters", stmt: "CREATE TABLE z (a INT, KEY ((_latin1'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'))) ROW_FORMAT=COMPACT", code: 0},
		{name: "COMPACT functional literal plain of 192 characters", stmt: "CREATE TABLE z (a INT, KEY (('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'))) ROW_FORMAT=COMPACT", code: 1071},
		{name: "COMPACT functional literal _utf8mb3 of 256 characters", stmt: "CREATE TABLE z (a INT, KEY ((_utf8mb3'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'))) ROW_FORMAT=COMPACT", code: 1071},
		// A literal's own charset sizes it in a binary CONCAT or REPLACE.
		{name: "functional binary CONCAT with a latin1 literal", stmt: "CREATE TABLE z (b VARBINARY(511), KEY ((CONCAT(b, _latin1'x'))))", code: 0},
		{name: "functional binary CONCAT with a utf8mb4 literal", stmt: "CREATE TABLE z (b VARBINARY(511), KEY ((CONCAT(b, 'x'))))", code: 3757},
		{name: "functional binary CONCAT with a national literal", stmt: "CREATE TABLE z (b VARBINARY(509), KEY ((CONCAT(b, N'x'))))", code: 0},
		{name: "functional binary REPLACE with a latin1 replacement", stmt: "CREATE TABLE z (b VARBINARY(100), KEY ((REPLACE(b, 'x', _latin1'yyyyy'))))", code: 0},
		{name: "functional binary REPLACE with a utf8mb4 replacement", stmt: "CREATE TABLE z (b VARBINARY(100), KEY ((REPLACE(b, 'x', 'yyyyy'))))", code: 3757},
		// A key part's byte length wraps at 32 bits on the server.
		{name: "TEXT key prefix whose bytes wrap to 4", stmt: "CREATE TABLE a (s TEXT, KEY (s(1073741825)))", code: 0, contains: []string{"KEY `s` (`s`(1))"}},
		{name: "TEXT key prefix whose bytes wrap to 0", stmt: "CREATE TABLE z (s TEXT, KEY (s(1073741824)))", code: 1170},
		{name: "latin1 TEXT key prefix of 2^30", stmt: "CREATE TABLE z (s TEXT CHARACTER SET latin1, KEY (s(1073741824)))", code: 1071},
		{name: "VARCHAR key prefix whose bytes wrap to 0", stmt: "CREATE TABLE z (s VARCHAR(10), KEY (s(1073741824)))", code: 0},
		{name: "TEXT key prefix of 2^31-1", stmt: "CREATE TABLE z (s TEXT, KEY (s(2147483647)))", code: 1071},
		{name: "TEXT key prefix past 2^31-1", stmt: "CREATE TABLE z (s TEXT, KEY (s(9223372036854775807)))", code: 1064},
		// A signed literal renders with its sign in a binary CONCAT.
		{name: "functional binary CONCAT with -1 at 513 bytes", stmt: "CREATE TABLE z (b VARBINARY(511), KEY ((CONCAT(b, -1))))", code: 3757},
		{name: "functional binary CONCAT with -1 at 512 bytes", stmt: "CREATE TABLE z (b VARBINARY(510), KEY ((CONCAT(b, -1))))", code: 0},
		{name: "functional binary CONCAT with -1.5 at 512 bytes", stmt: "CREATE TABLE z (b VARBINARY(508), KEY ((CONCAT(b, -1.5))))", code: 0},
		// Binary strings of unequal length do not combine bitwise.
		{name: "bitwise AND of unequal binary strings", stmt: "SET @x = CAST(X'0102' AS BINARY) & CAST(X'FF' AS BINARY)", code: 3513},
		// default_table_encryption sets the default for schemas and general
		// tablespaces, not tables: a table that omits ENCRYPTION is unencrypted,
		// so an unencrypted general tablespace takes it (and, on a server with a
		// keyring, an encrypted one refuses it with 3825).
		{
			name:    "default_table_encryption ON and a table in an unencrypted tablespace",
			setup:   []string{"CREATE TABLESPACE ts_dte_n ADD DATAFILE 'ts_dte_n.ibd' ENCRYPTION='N'"},
			session: []string{"SET SESSION default_table_encryption = ON"},
			stmt:    "CREATE TABLE z (id INT) TABLESPACE ts_dte_n",
			code:    0,
			cleanup: []string{"SET SESSION default_table_encryption = DEFAULT", "DROP TABLE IF EXISTS z", "DROP TABLESPACE ts_dte_n"},
		},
		// Functions in a SET expression are resolved before a branch is chosen:
		// an unknown unqualified name is a stored function of the current
		// database (1305), and aggregates (1111) and window functions (3593)
		// are refused even in the branch not taken.
		{name: "SET IF over an unknown function", stmt: "SET sql_mode = IF(1, 'ANSI', no_such())", code: 1305},
		{name: "SET user variable IF over an unknown function", stmt: "SET @x = IF(1, 1, no_such())", code: 1305},
		{name: "SET IF over an aggregate", stmt: "SET sql_mode = IF(1, 'ANSI', COUNT(1))", code: 1111},
		{name: "SET IF over a window function", stmt: "SET sql_mode = IF(1, 'ANSI', ROW_NUMBER() OVER ())", code: 3593},
		{name: "SET IF over a built-in DEFAULT refuses", stmt: "SET sql_mode = IF(1, 'ANSI', VERSION())", code: 0},
		{name: "SET IF over REPEAT", stmt: "SET sql_mode = IF(1, 'ANSI', REPEAT('a', 2))", code: 0},
		{
			name:    "SET IF over a stored function",
			setup:   []string{"CREATE FUNCTION f_set_branch() RETURNS INT DETERMINISTIC RETURN 1"},
			stmt:    "SET @x = IF(1, 1, f_set_branch())",
			code:    0,
			cleanup: []string{"DROP FUNCTION IF EXISTS f_set_branch"},
		},
		// REPEAT is a deterministic built-in DEFAULT (expr) takes; CREATE_DIGEST
		// belongs to a component the pinned server has not loaded.
		{name: "DEFAULT (REPEAT)", stmt: "CREATE TABLE z (s VARCHAR(10) DEFAULT (REPEAT('a', 3)))", code: 0},
		{name: "DEFAULT (REPEAT) with one argument", stmt: "CREATE TABLE z (s VARCHAR(10) DEFAULT (REPEAT('a')))", code: 1064},
		{name: "DEFAULT (CREATE_DIGEST)", stmt: "CREATE TABLE z (s VARCHAR(64) DEFAULT (CREATE_DIGEST('SHA256', 'a')))", code: 3770},
		{
			// A general tablespace with no ADD DATAFILE gets a file named by a
			// UUID, not <name>.ibd, so another may take that name.
			name:    "ADD DATAFILE named like an implicit tablespace file",
			setup:   []string{"CREATE TABLESPACE ts_impl1 ENGINE=InnoDB"},
			stmt:    "CREATE TABLESPACE ts_impl2 ADD DATAFILE 'ts_impl1.ibd' ENGINE=InnoDB",
			code:    0,
			cleanup: []string{"DROP TABLESPACE ts_impl2", "DROP TABLESPACE ts_impl1"},
		},
		// The type a functional key part's result takes: a BLOB, TEXT or binary
		// column keeps its type through string functions, and a string or binary
		// result longer than 512 characters is a LOB (3757).
		{
			name: "functional key part type: KEY ((REVERSE(b)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(b))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LEFT(b, 10)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(b, 10))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((SUBSTRING(b, 1, 10)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((SUBSTRING(b, 1, 10))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((CONCAT(b, 'x')))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(b, 'x'))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LOWER(b)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LOWER(b))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LEFT(b, 512)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(b, 512))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((LEFT(b, 513)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(b, 513))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((REVERSE(tb)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(tb))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((UPPER(tb)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((UPPER(tb))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((CONCAT(tb, tb, 'x')))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(tb, tb, 'x'))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((REVERSE(vb)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(vb))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((CONCAT(vb, 'x')))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(vb, 'x'))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((REVERSE(bn)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(bn))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((REVERSE(vbw)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(vbw))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((CONCAT(vbw, vbw)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(vbw, vbw))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LEFT(vbw, 512)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(vbw, 512))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((LEFT(vbw, 513)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(vbw, 513))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((REVERSE(tt)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(tt))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((REVERSE(ttl)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(ttl))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((LEFT(tx, 10)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(tx, 10))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((REVERSE(v)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(v))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LEFT(v, 512)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(v, 512))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((LEFT(v, 513)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(v, 513))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((REVERSE(l)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(l))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((CAST(i AS CHAR(512))))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CAST(i AS CHAR(512)))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((CAST(i AS CHAR(513))))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CAST(i AS CHAR(513)))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((CAST(i AS BINARY(513))))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CAST(i AS BINARY(513)))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((CONCAT(v5, v5)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(v5, v5))))",
			code: 0,
		},
		{
			name: "functional key part type: KEY ((CONCAT(v5, v5, v5)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(v5, v5, v5))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LOWER(v)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LOWER(v))))",
			code: 3757,
		},
		{
			name: "functional key part type: KEY ((LOWER(w1)), (LOWER(w2)))",
			stmt: "CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LOWER(w1)), (LOWER(w2))))",
			code: 1071,
		},
		// Functional key parts (mysql:8.0.46): the built-ins DEFAULT (expr)
		// takes but a generated column does not, the ones DEFAULT refuses too,
		// stored and unknown functions, variables judged in post-order,
		// aggregates and window functions, and FULLTEXT / SPATIAL keys, which
		// take no functional part (3759 / 3760, after the expression's checks).
		{
			name:  "functional key part: CREATE INDEX i ON a ((CONNECTION_ID()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((CONNECTION_ID()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((UTC_TIMESTAMP()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((UTC_TIMESTAMP()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((UTC_TIMESTAMP(3)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((UTC_TIMESTAMP(3)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((CURRENT_DATE))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((CURRENT_DATE))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((CURRENT_USER()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((CURRENT_USER()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((DATABASE()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((DATABASE()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((UNIX_TIMESTAMP(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((UNIX_TIMESTAMP(c)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((UUID_SHORT()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((UUID_SHORT()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((STATEMENT_DIGEST_TEXT('select 1')))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((STATEMENT_DIGEST_TEXT('select 1')))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((LAST_INSERT_ID()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((LAST_INSERT_ID()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((VERSION()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((VERSION()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((GET_LOCK('a', 1)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((GET_LOCK('a', 1)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((VALUES(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((VALUES(c)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((NAME_CONST('a', 1)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((NAME_CONST('a', 1)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((nosuchfn(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((nosuchfn(c)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((d.nosuchfn(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((d.nosuchfn(c)))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((MATCH(c) AGAINST ('x')))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((MATCH(c) AGAINST ('x')))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((@x))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((@x))",
			code:  3772,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((@x + NOW()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((@x + NOW()))",
			code:  3772,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((NOW() + @x))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((NOW() + @x))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((SLEEP(@x)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((SLEEP(@x)))",
			code:  3772,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((COUNT(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((COUNT(c)))",
			code:  1111,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((GROUP_CONCAT(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((GROUP_CONCAT(c)))",
			code:  1111,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((COUNT(c) + NOW()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((COUNT(c) + NOW()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((ROW_NUMBER() OVER ()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((ROW_NUMBER() OVER ()))",
			code:  3593,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((COUNT(c) OVER ()))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((COUNT(c) OVER ()))",
			code:  3593,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((ANY_VALUE(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((ANY_VALUE(c)))",
			code:  0,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((CHAR(65)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((CHAR(65)))",
			code:  0,
		},
		{
			name:  "functional key part: CREATE INDEX i ON a ((ELT(1, c)))",
			setup: []string{fiTable},
			stmt:  "CREATE INDEX i ON a ((ELT(1, c)))",
			code:  0,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a ((LOWER(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a ((LOWER(c)))",
			code:  3759,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a (c, (LOWER(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a (c, (LOWER(c)))",
			code:  3759,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a ((CAST(c AS JSON)))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a ((CAST(c AS JSON)))",
			code:  3759,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a ((c))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a ((c))",
			code:  3762,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a ((NOW()))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a ((NOW()))",
			code:  3758,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a ((COUNT(c)))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a ((COUNT(c)))",
			code:  1111,
		},
		{
			name:  "functional key part: CREATE FULLTEXT INDEX i ON a ((LOWER(nosuch)))",
			setup: []string{fiTable},
			stmt:  "CREATE FULLTEXT INDEX i ON a ((LOWER(nosuch)))",
			code:  1054,
		},
		{
			name:  "functional key part: CREATE SPATIAL INDEX i ON s ((ST_SRID(g)))",
			setup: []string{fiGeometry},
			stmt:  "CREATE SPATIAL INDEX i ON s ((ST_SRID(g)))",
			code:  3760,
		},
		{
			name:  "functional key part: CREATE SPATIAL INDEX i ON s ((g))",
			setup: []string{fiGeometry},
			stmt:  "CREATE SPATIAL INDEX i ON s ((g))",
			code:  3762,
		},
		{
			name:  "functional key part: ALTER TABLE a ADD FULLTEXT ((LOWER(c)))",
			setup: []string{fiTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT ((LOWER(c)))",
			code:  3759,
		},
		{
			name:  "functional key part: ALTER TABLE a ADD FULLTEXT ((c))",
			setup: []string{fiTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT ((c))",
			code:  3762,
		},
		{
			name:  "functional key part: ALTER TABLE a ADD FULLTEXT nosuch_idx (nosuch, (LOWER(c)))",
			setup: []string{fiTable},
			stmt:  "ALTER TABLE a ADD FULLTEXT nosuch_idx (nosuch, (LOWER(c)))",
			code:  1072,
		},
		{
			name:  "functional key part: ALTER TABLE a ADD SPATIAL ((LOWER(c)))",
			setup: []string{fiTable},
			stmt:  "ALTER TABLE a ADD SPATIAL ((LOWER(c)))",
			code:  3760,
		},
		{
			name: "functional key part: CREATE TABLE z (c VARCHAR(10), FULLTEXT ((LOWER(c))))",
			stmt: "CREATE TABLE z (c VARCHAR(10), FULLTEXT ((LOWER(c))))",
			code: 3759,
		},
		{
			name: "functional key part: CREATE TABLE z (c VARCHAR(10), FULLTEXT ((c)))",
			stmt: "CREATE TABLE z (c VARCHAR(10), FULLTEXT ((c)))",
			code: 3762,
		},
		{
			name: "functional key part: CREATE TABLE z (c VARCHAR(10), FULLTEXT ((NOW())))",
			stmt: "CREATE TABLE z (c VARCHAR(10), FULLTEXT ((NOW())))",
			code: 3758,
		},
		{
			name: "functional key part: CREATE TABLE z (c VARCHAR(10), FULLTEXT ((LOWER(nosuch))))",
			stmt: "CREATE TABLE z (c VARCHAR(10), FULLTEXT ((LOWER(nosuch))))",
			code: 1054,
		},
		{
			name: "functional key part: CREATE TABLE z (c VARCHAR(10), SPATIAL ((LOWER(c))))",
			stmt: "CREATE TABLE z (c VARCHAR(10), SPATIAL ((LOWER(c))))",
			code: 3760,
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

	// Chained SET NAMES latin1 -> CREATE TABLE DEFAULT -> SET NAMES utf8mb4
	// stores decoded latin1 characters on both live MySQL 8.0.46 and Catalog,
	// and @@character_set_client save/restore round-trips cleanly.
	reset()
	for _, stmt := range []string{
		"SET NAMES latin1",
		"CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')",
		"SET NAMES utf8mb4",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	var srvTblName, srvCreate string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&srvTblName, &srvCreate); err != nil {
		t.Fatalf("SHOW CREATE TABLE a: %v", err)
	}
	latinCat := New()
	execAlterClauses(t, latinCat, "CREATE DATABASE d; USE d; SET NAMES latin1; CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é'); SET NAMES utf8mb4;")
	catCreate := latinCat.ShowCreateTable("d", "a")
	if !strings.Contains(srvCreate, "DEFAULT 'Ã©'") || !strings.Contains(catCreate, "DEFAULT 'Ã©'") {
		t.Fatalf("SET NAMES latin1 decoded default mismatch:\nmysql:\n%s\ncatalog:\n%s", srvCreate, catCreate)
	}
	execAlterClauses(t, latinCat, "SET @saved_cs = @@character_set_client; SET character_set_client = ascii; SET character_set_client = @saved_cs;")
	if latinCat.effectiveClientCharset() != "utf8mb4" {
		t.Fatalf("@@character_set_client round-trip got %q, want utf8mb4", latinCat.effectiveClientCharset())
	}

	// Named time zones behave identically to live MySQL 8.0.46 when configured
	// on the catalog via SetNamedTimeZones.
	reset()
	if _, err := conn.ExecContext(ctx, "SET SESSION foreign_key_checks = 1"); err != nil {
		t.Fatalf("mysql reset foreign_key_checks: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "SET time_zone = DEFAULT")
	})
	for _, stmt := range []string{
		"SET time_zone = 'America/New_York'",
		"CREATE TABLE a (dt DATETIME DEFAULT '2026-06-01 12:00:00+00:00')",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	var tzName, tzEngineCreate string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&tzName, &tzEngineCreate); err != nil {
		t.Fatalf("mysql SHOW CREATE TABLE a (tz): %v", err)
	}
	tzOracleCat := New()
	tzOracleCat.SetNamedTimeZones("America/New_York")
	execAlterClauses(t, tzOracleCat, "CREATE DATABASE d; USE d; SET time_zone = 'America/New_York'; CREATE TABLE a (dt DATETIME DEFAULT '2026-06-01 12:00:00+00:00');")
	if !strings.Contains(tzEngineCreate, "DEFAULT '2026-06-01 08:00:00'") {
		t.Fatalf("mysql named time_zone SHOW CREATE TABLE:\n%s", tzEngineCreate)
	}
	if gotShow := tzOracleCat.ShowCreateTable("d", "a"); !strings.Contains(gotShow, "DEFAULT '2026-06-01 08:00:00'") {
		t.Fatalf("catalog named time_zone SHOW CREATE TABLE:\n%s", gotShow)
	}

	// SHOW CREATE TABLE output for a general-tablespace table loaded via
	// LoadSDL preserves tablespace encryption inference across subsequent ALTER.
	reset()
	if _, err := conn.ExecContext(ctx, "SET SESSION foreign_key_checks = 1"); err != nil {
		t.Fatalf("mysql reset foreign_key_checks: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "DROP TABLE IF EXISTS a")
		_, _ = conn.ExecContext(ctx, "DROP TABLESPACE ts_sdl_plain")
	})
	for _, stmt := range []string{
		"CREATE TABLESPACE ts_sdl_plain ADD DATAFILE 'ts_sdl_plain.ibd'",
		"CREATE TABLE a (id INT PRIMARY KEY) TABLESPACE ts_sdl_plain",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	var sdlTblName, sdlEngineCreate string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&sdlTblName, &sdlEngineCreate); err != nil {
		t.Fatalf("mysql SHOW CREATE TABLE a (sdl): %v", err)
	}
	sdlCat, err := LoadSDL("CREATE DATABASE d; USE d; " + sdlEngineCreate + ";")
	if err != nil {
		t.Fatalf("LoadSDL(%s): %v", sdlEngineCreate, err)
	}
	execAlterClauses(t, sdlCat, "ALTER TABLE a ADD COLUMN c INT;")
	if got, _ := simulationCatalogRefusal(t, sdlCat, "ALTER TABLE a ENCRYPTION='Y'"); got != 3825 {
		t.Fatalf("post-SDL conflicting ENCRYPTION='Y' got %d, want 3825", got)
	}

	// With keyring_file loaded on live MySQL 8.0.46, pin ALTER TABLESPACE
	// ENCRYPTION='Y'/'N' propagation to resident tables AND LoadSDL of an
	// encrypted general-tablespace table from live SHOW CREATE TABLE output.
	reset()
	if _, err := conn.ExecContext(ctx, "SET SESSION foreign_key_checks = 1"); err != nil {
		t.Fatalf("mysql reset foreign_key_checks: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "DROP TABLE IF EXISTS a")
		_, _ = conn.ExecContext(ctx, "DROP TABLESPACE ts_enc_chain")
		_, _ = conn.ExecContext(ctx, "UNINSTALL PLUGIN keyring_file")
	})
	if _, err := conn.ExecContext(ctx, "INSTALL PLUGIN keyring_file SONAME 'keyring_file.so'"); err != nil {
		t.Fatalf("mysql INSTALL PLUGIN keyring_file: %v", err)
	}
	encOracleCat := New()
	execAlterClauses(t, encOracleCat, "CREATE DATABASE d; USE d;")
	for _, stmt := range []string{
		"CREATE TABLESPACE ts_enc_chain ADD DATAFILE 'ts_enc_chain.ibd'",
		"CREATE TABLE a (id INT PRIMARY KEY) TABLESPACE ts_enc_chain",
		"ALTER TABLESPACE ts_enc_chain ENCRYPTION='Y'",
		"ALTER TABLE a ADD COLUMN c INT",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, encOracleCat, stmt+";")
	}
	var encTblName, encEngineCreate string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&encTblName, &encEngineCreate); err != nil {
		t.Fatalf("mysql SHOW CREATE TABLE a (encrypted): %v", err)
	}
	encCatCreate := encOracleCat.ShowCreateTable("d", "a")
	if !strings.Contains(encEngineCreate, "ENCRYPTION='Y'") || !strings.Contains(encCatCreate, "ENCRYPTION='Y'") {
		t.Fatalf("encrypted resident SHOW CREATE TABLE mismatch:\nmysql:\n%s\ncatalog:\n%s", encEngineCreate, encCatCreate)
	}
	sdlEncLiveCat, err := LoadSDL("CREATE DATABASE d; USE d; " + encEngineCreate + ";")
	if err != nil {
		t.Fatalf("LoadSDL(%s): %v", encEngineCreate, err)
	}
	execAlterClauses(t, sdlEncLiveCat, "ALTER TABLE a ADD COLUMN c_sdl INT;")
	if got, _ := simulationCatalogRefusal(t, sdlEncLiveCat, "CREATE TABLE t_plain_sdl (id INT) TABLESPACE ts_enc_chain"); got != 3825 {
		t.Fatalf("post-encrypted-SDL plain table got %d, want 3825", got)
	}
	for _, stmt := range []string{
		"ALTER TABLESPACE ts_enc_chain ENCRYPTION='N'",
		"ALTER TABLE a ADD COLUMN c2 INT",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, encOracleCat, stmt+";")
	}
	var decTblName, decEngineCreate string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&decTblName, &decEngineCreate); err != nil {
		t.Fatalf("mysql SHOW CREATE TABLE a (decrypted): %v", err)
	}
	decCatCreate := encOracleCat.ShowCreateTable("d", "a")
	if strings.Contains(decEngineCreate, "ENCRYPTION=") || strings.Contains(decCatCreate, "ENCRYPTION=") {
		t.Fatalf("decrypted resident SHOW CREATE TABLE still had ENCRYPTION:\nmysql:\n%s\ncatalog:\n%s", decEngineCreate, decCatCreate)
	}

	// Under default_table_encryption=ON with keyring_file loaded, an omitted
	// ENCRYPTION clause on CREATE TABLESPACE defaults to ENCRYPTION='Y'.
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "SET SESSION default_table_encryption = OFF")
		_, _ = conn.ExecContext(ctx, "DROP TABLE IF EXISTS a_def_enc")
		_, _ = conn.ExecContext(ctx, "DROP TABLESPACE ts_def_enc")
	})
	for _, stmt := range []string{
		"SET SESSION default_table_encryption = ON",
		"CREATE TABLESPACE ts_def_enc ADD DATAFILE 'ts_def_enc.ibd'",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, encOracleCat, stmt+";")
	}
	if got, _ := simulationEngineRefusal(t, "CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc", func() error {
		_, err := conn.ExecContext(ctx, "CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc")
		return err
	}()); got != 3825 {
		t.Fatalf("mysql plain table in default-encrypted tablespace got %d, want 3825", got)
	}
	if got, _ := simulationCatalogRefusal(t, encOracleCat, "CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc"); got != 3825 {
		t.Fatalf("catalog plain table in default-encrypted tablespace got %d, want 3825", got)
	}
	if _, err := conn.ExecContext(ctx, "CREATE TABLESPACE ts_omit_enc ADD DATAFILE 'ts_omit_enc.ibd' ENCRYPTION='Y'"); err != nil {
		t.Fatalf("mysql CREATE TABLESPACE ts_omit_enc: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "DROP TABLE IF EXISTS a_def_omit_y")
		_, _ = conn.ExecContext(ctx, "DROP TABLESPACE ts_omit_enc")
	})
	execAlterClauses(t, encOracleCat, "CREATE TABLESPACE ts_omit_enc ADD DATAFILE 'ts_omit_enc.ibd' ENCRYPTION='Y';")
	if _, err := conn.ExecContext(ctx, "CREATE TABLE a_def_omit_y (id INT) TABLESPACE ts_omit_enc"); err == nil {
		t.Fatal("mysql plain table with omitted ENCRYPTION in encrypted tablespace succeeded, want 3825")
	} else if code, _ := simulationEngineRefusal(t, "CREATE TABLE a_def_omit_y (id INT) TABLESPACE ts_omit_enc", err); code != 3825 {
		t.Fatalf("mysql plain table in encrypted ts_omit_enc got %d, want 3825", code)
	}
	if got, _ := simulationCatalogRefusal(t, encOracleCat, "CREATE TABLE a_def_omit_y (id INT) TABLESPACE ts_omit_enc"); got != 3825 {
		t.Fatalf("catalog plain table in encrypted ts_omit_enc got %d, want 3825", got)
	}
	if _, err := conn.ExecContext(ctx, "CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc ENCRYPTION='Y'"); err != nil {
		t.Fatalf("mysql encrypted table in default-encrypted tablespace: %v", err)
	}
	execAlterClauses(t, encOracleCat, "CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc ENCRYPTION='Y';")

	// SET CHARACTER SET preserves the current database's non-default collation
	// as @@collation_connection on both live MySQL 8.0.46 and Catalog.
	reset()
	var origDB string
	if err := conn.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&origDB); err != nil {
		t.Fatalf("mysql SELECT DATABASE(): %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "USE `"+origDB+"`")
		_, _ = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS db_cs_bin")
		_, _ = conn.ExecContext(ctx, "SET NAMES utf8mb4")
	})
	for _, stmt := range []string{
		"CREATE DATABASE db_cs_bin CHARACTER SET latin1 COLLATE latin1_bin",
		"USE db_cs_bin",
		"SET CHARACTER SET utf8mb4",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	var srvCSClient, srvCollConn string
	if err := conn.QueryRowContext(ctx, "SELECT @@character_set_client, @@collation_connection").Scan(&srvCSClient, &srvCollConn); err != nil {
		t.Fatalf("mysql SELECT charset/collation: %v", err)
	}
	csOracleCat := New()
	execAlterClauses(t, csOracleCat, "CREATE DATABASE db_cs_bin CHARACTER SET latin1 COLLATE latin1_bin; USE db_cs_bin; SET CHARACTER SET utf8mb4; SET @cc = @@collation_connection;")
	if srvCSClient != "utf8mb4" || srvCollConn != "latin1_bin" {
		t.Fatalf("mysql SET CHARACTER SET got client=%q coll=%q, want utf8mb4/latin1_bin", srvCSClient, srvCollConn)
	}
	if csOracleCat.effectiveClientCharset() != "utf8mb4" || csOracleCat.effectiveSessionCollationConn() != "latin1_bin" || csOracleCat.userVars["cc"].str != "latin1_bin" {
		t.Fatalf("catalog SET CHARACTER SET lost db.Collation: client=%q coll=%q var=%q", csOracleCat.effectiveClientCharset(), csOracleCat.effectiveSessionCollationConn(), csOracleCat.userVars["cc"].str)
	}

	// Pin multi-assignment SET atomicity, pre-statement RHS evaluation, and
	// REPLACE(@@sql_mode, ...) directly against live MySQL 8.0.46.
	reset()
	if _, err := conn.ExecContext(ctx, "SET SESSION foreign_key_checks = 1"); err != nil {
		t.Fatalf("mysql reset foreign_key_checks: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "SET SESSION foreign_key_checks = 1")
		_, _ = conn.ExecContext(ctx, "SET SESSION sql_mode = DEFAULT")
	})
	_, errFailSET := conn.ExecContext(ctx, "SET foreign_key_checks = 0, sql_mode = 'NO_SUCH_MODE'")
	if code, _ := simulationEngineRefusal(t, "SET foreign_key_checks = 0, sql_mode = 'NO_SUCH_MODE'", errFailSET); code != 1231 {
		t.Fatalf("mysql partial-fail SET got %d, want 1231", code)
	}
	var fkAfter int
	if err := conn.QueryRowContext(ctx, "SELECT @@foreign_key_checks").Scan(&fkAfter); err != nil {
		t.Fatalf("mysql SELECT @@foreign_key_checks: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "SET @x_seq = 1, @y_seq = @x_seq + 1"); err != nil {
		t.Fatalf("mysql SET @x_seq,@y_seq: %v", err)
	}
	var seqY *string
	if err := conn.QueryRowContext(ctx, "SELECT @y_seq").Scan(&seqY); err != nil {
		t.Fatalf("mysql SELECT @y_seq: %v", err)
	}
	atomCat := New()
	if got, _ := simulationCatalogRefusal(t, atomCat, "SET foreign_key_checks = 0, sql_mode = 'NO_SUCH_MODE'"); got != 1231 {
		t.Fatalf("catalog partial-fail SET got %d, want 1231", got)
	}
	if fkAfter != 1 || !atomCat.ForeignKeyChecks() {
		t.Fatalf("SET partial failure leaked foreign_key_checks=0: mysql=%d catalog=%v", fkAfter, atomCat.ForeignKeyChecks())
	}
	execAlterClauses(t, atomCat, "SET @x_seq = 1, @y_seq = @x_seq + 1;")
	if seqY != nil || atomCat.userVars["y_seq"].kind != userVarNull {
		t.Fatalf("SET RHS read uncommitted earlier assignment: mysql=%v catalog=%+v", seqY, atomCat.userVars["y_seq"])
	}
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', '')"); err != nil {
		t.Fatalf("mysql SET REPLACE sql_mode: %v", err)
	}
	var srvModeAfter string
	if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&srvModeAfter); err != nil {
		t.Fatalf("mysql SELECT @@sql_mode: %v", err)
	}
	execAlterClauses(t, atomCat, "SET SESSION sql_mode = REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', ''); SET @sm = @@sql_mode;")
	if atomCat.userVars["sm"].str != srvModeAfter {
		t.Fatalf("REPLACE(@@sql_mode, ...) mismatch: catalog=%q mysql=%q", atomCat.userVars["sm"].str, srvModeAfter)
	}

	// Schema-qualified SET function resolves stored functions (1305 when missing,
	// accepted without evaluating the builtin when defined on the schema).
	reset()
	var qfnOrigDB string
	if err := conn.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&qfnOrigDB); err != nil {
		t.Fatalf("mysql SELECT DATABASE() (qfn): %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "USE `"+qfnOrigDB+"`")
		_, _ = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS d_qfn")
		_, _ = conn.ExecContext(ctx, "SET SESSION sql_mode = DEFAULT")
	})
	if _, err := conn.ExecContext(ctx, "CREATE DATABASE d_qfn; USE d_qfn;"); err != nil {
		t.Fatalf("mysql CREATE DATABASE d_qfn: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "SET sql_mode = d_qfn.CONCAT('ANSI')"); err == nil {
		t.Fatal("mysql missing d_qfn.CONCAT succeeded, want 1305")
	} else if code, _ := simulationEngineRefusal(t, "SET sql_mode = d_qfn.CONCAT('ANSI')", err); code != 1305 {
		t.Fatalf("mysql missing d_qfn.CONCAT got %d, want 1305", code)
	}
	qfnCat := New()
	execAlterClauses(t, qfnCat, "CREATE DATABASE d_qfn; USE d_qfn;")
	if got, _ := simulationCatalogRefusal(t, qfnCat, "SET sql_mode = d_qfn.CONCAT('ANSI')"); got != 1305 {
		t.Fatalf("catalog missing d_qfn.CONCAT got %d, want 1305", got)
	}
	if _, err := conn.ExecContext(ctx, "CREATE FUNCTION d_qfn.CONCAT(s VARCHAR(20)) RETURNS VARCHAR(20) DETERMINISTIC RETURN ''"); err != nil {
		t.Fatalf("mysql CREATE FUNCTION d_qfn.CONCAT: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "SET sql_mode = d_qfn.CONCAT('ANSI')"); err != nil {
		t.Fatalf("mysql defined d_qfn.CONCAT: %v", err)
	}
	execAlterClauses(t, qfnCat, "CREATE FUNCTION d_qfn.CONCAT(s VARCHAR(20)) RETURNS VARCHAR(20) DETERMINISTIC RETURN ''; SET sql_mode = d_qfn.CONCAT('ANSI');")
	if strings.Contains(qfnCat.session.SQLMode, "ANSI") {
		t.Fatalf("catalog evaluated builtin CONCAT instead of stored d_qfn.CONCAT: sql_mode=%q", qfnCat.session.SQLMode)
	}

	// Schema-and-table qualified references (d.t.a) in column DEFAULT expressions
	// and functional indexes are accepted when schema and table match the owning table.
	reset()
	if _, err := conn.ExecContext(ctx, "CREATE DATABASE d_qual_own; USE d_qual_own;"); err != nil {
		t.Fatalf("mysql CREATE DATABASE d_qual_own: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "USE `"+qfnOrigDB+"`")
		_, _ = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS d_qual_own")
	})
	if _, err := conn.ExecContext(ctx, "CREATE TABLE q_own (a INT, b INT DEFAULT (d_qual_own.q_own.a + 1), KEY ((d_qual_own.q_own.a + 1)))"); err != nil {
		t.Fatalf("mysql schema+table qualified CREATE TABLE: %v", err)
	}
	qualCat := New()
	execAlterClauses(t, qualCat, "CREATE DATABASE d_qual_own; USE d_qual_own; CREATE TABLE q_own (a INT, b INT DEFAULT (d_qual_own.q_own.a + 1), KEY ((d_qual_own.q_own.a + 1)));")

	// Live MySQL 8.0.46 oracle checks for skipped-branch sysvar validation and scientific notation SET coercion.
	reset()
	t.Cleanup(func() {
		_, _ = conn.ExecContext(ctx, "SET SESSION sql_mode = DEFAULT")
	})
	setExprOracleCat := New()
	execAlterClauses(t, setExprOracleCat, "SET SESSION sql_mode = '';")
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = ''"); err != nil {
		t.Fatalf("mysql reset sql_mode: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "SET sql_mode = IF(1, 'ANSI', @@no_such_variable)"); err == nil {
		t.Fatal("mysql IF(..., @@no_such_variable) succeeded, want 1193")
	} else if code, _ := simulationEngineRefusal(t, "SET sql_mode = IF(1, 'ANSI', @@no_such_variable)", err); code != 1193 {
		t.Fatalf("mysql IF(..., @@no_such_variable) got %d, want 1193", code)
	}
	if got, _ := simulationCatalogRefusal(t, setExprOracleCat, "SET sql_mode = IF(1, 'ANSI', @@no_such_variable)"); got != 1193 {
		t.Fatalf("catalog IF(..., @@no_such_variable) got %d, want 1193", got)
	}
	if setExprOracleCat.session.SQLMode != "" {
		t.Fatalf("unknown sysvar failure mutated sql_mode: %q", setExprOracleCat.session.SQLMode)
	}
	for _, tc := range []struct {
		stmt string
		want string
	}{
		{"SET sql_mode = IF(1e2 = 100, 'ANSI', '')", "ANSI"},
		{"SET sql_mode = IF('123euro' = 123, 'STRICT_ALL_TABLES', '')", "STRICT_ALL_TABLES"},
		{"SET sql_mode = IF(1e2 + 1 = 101, 'ANSI', '')", "ANSI"},
	} {
		if _, err := conn.ExecContext(ctx, tc.stmt); err != nil {
			t.Fatalf("mysql %s: %v", tc.stmt, err)
		}
		var srvMode string
		if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&srvMode); err != nil {
			t.Fatalf("mysql SELECT @@sql_mode: %v", err)
		}
		execAlterClauses(t, setExprOracleCat, tc.stmt+"; SET @m = @@sql_mode;")
		if !strings.Contains(srvMode, tc.want) || setExprOracleCat.userVars["m"].str != srvMode {
			t.Fatalf("%s mismatch: srv=%q catalog=%q", tc.stmt, srvMode, setExprOracleCat.userVars["m"].str)
		}
	}
	if _, err := conn.ExecContext(ctx, "SET sql_mode = 1e2 + 1"); err == nil {
		t.Fatal("mysql SET sql_mode = 1e2 + 1 succeeded, want 1232")
	} else if code, _ := simulationEngineRefusal(t, "SET sql_mode = 1e2 + 1", err); code != 1232 {
		t.Fatalf("mysql SET sql_mode = 1e2 + 1 got %d, want 1232", code)
	}
	if got, _ := simulationCatalogRefusal(t, setExprOracleCat, "SET sql_mode = 1e2 + 1"); got != 1232 {
		t.Fatalf("catalog SET sql_mode = 1e2 + 1 got %d, want 1232", got)
	}

	// ELT in SET: N read as val_int reads it (a string by its integer prefix,
	// an exact number half away from zero, an approximate one half to even),
	// only strN evaluated, and a NULL answer refused by sql_mode.
	eltCat := New()
	execAlterClauses(t, eltCat, "SET SESSION sql_mode = '';")
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = ''"); err != nil {
		t.Fatalf("mysql reset sql_mode: %v", err)
	}
	for _, stmt := range []string{
		"SET sql_mode = ELT(1, 'ANSI', '')",
		"SET sql_mode = ELT('2', '', 'STRICT_ALL_TABLES')",
		"SET sql_mode = ELT('1.6', 'ANSI', 'STRICT_ALL_TABLES')",
		"SET sql_mode = ELT(1.6, 'ANSI', 'STRICT_ALL_TABLES')",
		"SET sql_mode = ELT(2.5, 'ANSI', 'TRADITIONAL', 'STRICT_ALL_TABLES')",
		"SET sql_mode = ELT(2.5e0, 'ANSI', 'TRADITIONAL', 'STRICT_ALL_TABLES')",
		"SET sql_mode = ELT(' 2', 'ANSI', '')",
		"SET sql_mode = ELT(1, 'ANSI', @@sql_mode)",
		"SET sql_mode = ELT(4294967297, 'ANSI', 'TRADITIONAL')",
		"SET sql_mode = ELT(0, 'ANSI')",
		"SET sql_mode = ELT(-1, 'ANSI')",
		"SET sql_mode = ELT(NULL, 'ANSI')",
		"SET sql_mode = ELT('x', 'ANSI')",
		"SET sql_mode = ELT(1, NULL, 'ANSI')",
		"SET sql_mode = ELT(1)",
		"SET sql_mode = ELT(1, 'ANSI', @@no_such_variable)",
	} {
		var srvCode int
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			srvCode, _ = simulationEngineRefusal(t, stmt, err)
		}
		catCode, _ := simulationCatalogRefusal(t, eltCat, stmt)
		if catCode != srvCode {
			t.Fatalf("%s: catalog %d, mysql %d", stmt, catCode, srvCode)
		}
		var srvMode string
		if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&srvMode); err != nil {
			t.Fatalf("mysql SELECT @@sql_mode: %v", err)
		}
		execAlterClauses(t, eltCat, "SET @m = @@sql_mode;")
		if got := eltCat.userVars["m"].str; got != srvMode {
			t.Fatalf("%s: catalog sql_mode %q, mysql %q", stmt, got, srvMode)
		}
	}
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = DEFAULT"); err != nil {
		t.Fatalf("mysql reset sql_mode: %v", err)
	}

	// ALTER TABLESPACE on the data dictionary tablespace `mysql`: it exists,
	// is InnoDB, keeps mysql.ibd and its name.
	for _, tc := range []struct {
		stmt string
		code int
	}{
		{"ALTER TABLESPACE mysql AUTOEXTEND_SIZE=4M", 0},
		{"ALTER TABLESPACE mysql ENGINE=InnoDB", 0},
		{"ALTER TABLESPACE mysql AUTOEXTEND_SIZE=3M", 4025},
		{"ALTER TABLESPACE mysql ENGINE=MyISAM", 3642},
		{"ALTER TABLESPACE mysql ENCRYPTION='Z' ENGINE_ATTRIBUTE='bad'", 3980},
		{"ALTER TABLESPACE mysql ADD DATAFILE 'x.ibd'", 1533},
		{"ALTER TABLESPACE mysql DROP DATAFILE 'x.ibd'", 3629},
		{"ALTER TABLESPACE mysql DROP DATAFILE 'mysql.ibd'", 1533},
		{"ALTER TABLESPACE mysql RENAME TO x", 3119},
		{"ALTER TABLESPACE mysql RENAME TO MYSQL", 3119},
		{"ALTER TABLESPACE mysql RENAME TO mysql", 1813},
		{"ALTER TABLESPACE mysql RENAME TO innodb_temporary", 1813},
		{"ALTER TABLESPACE MYSQL ENCRYPTION='N'", 3510},
	} {
		var srvCode int
		if _, err := conn.ExecContext(ctx, tc.stmt); err != nil {
			srvCode, _ = simulationEngineRefusal(t, tc.stmt, err)
		}
		if srvCode != tc.code {
			t.Fatalf("mysql %s got %d, want %d", tc.stmt, srvCode, tc.code)
		}
		if got, _ := simulationCatalogRefusal(t, New(), tc.stmt); got != tc.code {
			t.Fatalf("catalog %s got %d, want %d", tc.stmt, got, tc.code)
		}
	}
	// ENCRYPTION is the documented way to encrypt the tablespace: 'Y' and 'N'
	// both take, and the server's INNODB_TABLESPACES row follows.
	mysqlTSCat := New()
	t.Cleanup(func() { _, _ = conn.ExecContext(ctx, "ALTER TABLESPACE mysql ENCRYPTION='N'") })
	for _, enc := range []string{"Y", "N"} {
		stmt := "ALTER TABLESPACE mysql ENCRYPTION='" + enc + "'"
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		var srvEnc string
		if err := conn.QueryRowContext(ctx, "SELECT ENCRYPTION FROM information_schema.INNODB_TABLESPACES WHERE NAME = 'mysql'").Scan(&srvEnc); err != nil {
			t.Fatalf("mysql read mysql tablespace encryption: %v", err)
		}
		execAlterClauses(t, mysqlTSCat, stmt+";")
		if srvEnc != enc || mysqlTSCat.tablespaceEncrypted["mysql"] != (enc == "Y") {
			t.Fatalf("%s: mysql ENCRYPTION=%q, catalog encrypted=%v", stmt, srvEnc, mysqlTSCat.tablespaceEncrypted["mysql"])
		}
	}

	// CAST(... AS SIGNED / UNSIGNED) and SUBSTRING bounds in SET: the value a
	// user variable holds after each, and the sql_mode a cast chains into.
	castCat := New()
	for _, expr := range []string{
		"CAST(3.6 AS SIGNED)", "CAST(-3.6 AS SIGNED)", "CAST(2.5 AS SIGNED)", "CAST(-2.5 AS SIGNED)",
		"CAST(2.5e0 AS SIGNED)", "CAST(3.5e0 AS SIGNED)", "CAST(-1.5e0 AS SIGNED)",
		"CAST('3.6' AS SIGNED)", "CAST(' 3.6x' AS SIGNED)", "CAST('1e2' AS SIGNED)",
		"CAST(3.6 AS UNSIGNED)", "CAST(-1.6 AS UNSIGNED)", "CAST(-1 AS UNSIGNED)", "CAST('-1' AS UNSIGNED)",
		"CAST(99999999999999999999 AS SIGNED)", "CAST(-99999999999999999999 AS SIGNED)",
		"CAST(99999999999999999999 AS UNSIGNED)", "CAST(-99999999999999999999 AS UNSIGNED)",
		"CAST('99999999999999999999' AS SIGNED)", "CAST('10000000000000000000' AS SIGNED)",
		"CAST('-99999999999999999999' AS SIGNED)", "CAST('-99999999999999999999' AS UNSIGNED)",
		"CAST(18446744073709551615 AS SIGNED)", "CAST(9223372036854775808 AS SIGNED)",
		"CAST(9223372036854775807.5 AS SIGNED)", "CAST(1e30 AS SIGNED)", "CAST(1e30 AS UNSIGNED)",
		"CAST(1.5e19 AS UNSIGNED)", "CAST(-2.5e0 AS UNSIGNED)", "CAST(0.5 AS SIGNED)", "CAST(-0.5 AS UNSIGNED)",
		"SUBSTRING('ANSI', 2, 9223372036854775807)", "SUBSTRING('ANSI', -2, 9223372036854775807)",
		"SUBSTRING(_binary'ANSI', 2, 9223372036854775807)", "SUBSTRING('ANSI', 9223372036854775807, 2)",
		"SUBSTRING('ANSI', -9223372036854775807, 2)", "SUBSTRING('ANSI', 2, 18446744073709551615)",
		"SUBSTRING('ANSI', 18446744073709551615)", "SUBSTRING('ANSI', 2, -1)",
		"LEFT('ANSI', 9223372036854775807)", "LEFT('ANSI', 18446744073709551615)", "RIGHT('ANSI', 9223372036854775807)",
		"ABS(-4)", "ABS('-4')", "ABS(-4.5)", "ABS('abc')", "ABS('-1.5e0')", "ABS(-1.5e0)", "ABS(' -3x')",
		"ABS(4)", "ABS(-9223372036854775807)", "ABS(18446744073709551615)",
		"NOT 0", "NOT 1", "NOT 'abc'", "NOT '1x'", "!0", "NOT 0.0", "NOT 0.5",
		"LEFT('ANSI', 1.6)", "LEFT('ANSI', 1.5)", "LEFT('ANSI', '1.6')", "RIGHT('ANSI', 2.5e0)", "LEFT('ANSI', 1.4)",
		"SUBSTRING('ANSI', 1.6, 1.6)", "SUBSTRING('ANSI', '2.6')", "SUBSTRING('ANSI', 2.5, 2.5e0)",
		"LEFT('ANSI', 'x')", "LEFT('ANSI', -0.4)",
		"-99999999999999999999", "99999999999999999999 + 1", "-(99999999999999999999)",
		"ROUND(3.6)", "ROUND(-3.5)", "ROUND(2.5e0)", "ROUND(3.5e0)", "ROUND(1.2345, 2)", "ROUND(1234.5, -2)",
		"ROUND('3.6')", "ROUND(1.2, 3)", "ROUND(3.14159, 2.6)", "ROUND(-1.25, 1)", "ROUND(1.25e0, 1)",
		"ROUND(15, -1)", "ROUND(-15, -1)", "ROUND(14, -5)", "ROUND(1.5, 40)",
		"TRUNCATE(3.69, 1)", "TRUNCATE(-3.69, 1)", "TRUNCATE(1234, -2)", "TRUNCATE(1.2, 3)", "TRUNCATE('3.69', 1)", "TRUNCATE(3.69e0, 1)",
		"FLOOR(-3.5)", "FLOOR(3.5)", "CEIL(3.2)", "CEILING(-3.2)", "FLOOR(2.5e0)", "FLOOR('3.7')", "CEIL('x')", "CEIL(-0.5)",
		"FLOOR(99999999999999999999.5)", "SIGN(-2.5)", "SIGN(0)", "SIGN('x')", "SIGN('-3')", "SIGN(2e0)",
		"MOD(5, 4)", "MOD(-5, 4)", "MOD(5.5, 2)", "MOD('7', 4)", "MOD(5e0, 2)",
		"~-2", "~0", "~'3'", "~1.6", "~1.5e0", "~'x'",
		"BIT_COUNT(1)", "BIT_COUNT(-1)", "BIT_COUNT('3')", "BIT_COUNT(1.6)", "BIT_COUNT(1.5e0)", "BIT_COUNT('x')",
		"BIT_COUNT(18446744073709551615)", "BIT_COUNT(X'FF')", "BIT_COUNT(b'101')",
		"ASCII(X'01')", "ASCII('A')", "ASCII('')", "ASCII(65)",
		"ORD(X'01')", "ORD('A')", "ORD('')", "ORD('é')", "ORD('€')", "ORD('😀')", "ORD(65)", "ORD(_binary'é')",
		"BIT_LENGTH('a')", "BIT_LENGTH('é')", "BIT_LENGTH(X'0102')", "BIT_LENGTH(123)", "BIT_LENGTH('')", "BIT_LENGTH(1.5e0)",
		"LENGTH(1.5e0)", "OCTET_LENGTH('é')", "CHAR_LENGTH(1.5e0)",
		"LOCATE('a', 'a')", "LOCATE('A', 'a')", "LOCATE('b', 'abcb')", "LOCATE('b', 'abcb', 3)", "LOCATE('b', 'abcb', 0)",
		"LOCATE('', 'abc')", "LOCATE('', 'abc', 3)", "LOCATE('', 'abc', 4)", "LOCATE('', 'abc', 5)", "LOCATE('b', 'abcb', -1)",
		"INSTR('abcb', 'b')", "INSTR('a', 'A')", "LOCATE('b', 'abcb', 2.6)", "LOCATE(1, 21)",
		"CRC32('a')", "CRC32('')", "CRC32(123)", "CRC32('é')", "CRC32(X'00')", "CRC32(1.5e0)",
		"X'01' | X'02'", "X'0102' | X'0301'", "X'01' | X'0102'", "X'01' | 2", "~X'01'", "b'1' | b'10'", "X'01' << 1",
		"CAST(X'01' AS BINARY) | 2", "CAST(X'01' AS BINARY) | CAST(X'02' AS BINARY)",
		"CAST(X'0102' AS BINARY) & CAST(X'FF03' AS BINARY)", "CAST(X'01' AS BINARY) << 1", "~CAST(X'01' AS BINARY)",
		"CAST(X'8001' AS BINARY) >> 1",
		"CAST(2.5 AS DOUBLE)", "CAST(2.5 AS FLOAT)", "CAST(0.1 AS FLOAT)", "CAST('2.5x' AS DOUBLE)", "CAST(2.5 AS REAL)",
		"ABS(CAST(-9223372036854775808 AS DECIMAL(20,0)))", "ABS(CAST(-5 AS DECIMAL(3,0)))",
		"CONCAT(1.5e0)", "CONCAT(1e20)", "CONCAT(1e15)", "CONCAT(1.5e-7)", "CONCAT(123456789012345678e0)",
		"CONCAT(0.1e0)", "CONCAT(1e-5)", "CONCAT(-2.5e0)", "CONCAT(1e16)", "CONCAT(100e0)", "CONCAT(1e-16)",
		"CONCAT(1.2345678901234567e-10)", "CONCAT(1.5e-15)", "CONCAT(1e14)", "CONCAT(0.000123456789012345678e0)",
		"CONCAT(1.2345678901234567e-5)", "CONCAT(1.5e0 + 1)",
		"QUARTER('2008-01-01')", "QUARTER('2008-12-31')", "QUARTER(20080401)", "QUARTER('2008-04-01 10:00:00')",
		"QUARTER('2008-00-01')", "MONTH('2008-05-01')", "YEAR('2008-05-01')", "DAYOFMONTH('2008-05-07')", "DAY('2008-05-07')",
		"1 << 63", "(1 << 63) > 0", "-1 | 0", "1 << 64", "5 >> 1", "1.6 | 0", "'3' ^ 1", "-1 & 255",
	} {
		stmt := "SET @x = " + expr
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		// The value's bytes, so a binary result such as 0xFE reads back whole.
		var srvBytes []byte
		if err := conn.QueryRowContext(ctx, "SELECT CAST(@x AS BINARY)").Scan(&srvBytes); err != nil {
			t.Fatalf("mysql read @x after %s: %v", stmt, err)
		}
		srv := string(srvBytes)
		execAlterClauses(t, castCat, stmt+";")
		if got := castCat.userVars["x"].str; got != srv {
			t.Fatalf("%s: catalog %q, mysql %q", stmt, got, srv)
		}
	}
	// REPLACE takes its character set from its first argument only, so a
	// binary search or replacement string leaves the comparison case-insensitive.
	for _, tc := range []struct{ stmt, want string }{
		{"SET SESSION sql_mode = IF(REPLACE('ansi', X'78', '') = 'ANSI', 'ANSI', '')", "ANSI"},
		{"SET SESSION sql_mode = IF(REPLACE('ansi', 'x', X'78') = 'ANSI', 'ANSI', '')", "ANSI"},
		{"SET SESSION sql_mode = IF(REPLACE(X'616E7369', 'x', '') = 'ANSI', 'ANSI', '')", ""},
	} {
		for _, stmt := range []string{"SET SESSION sql_mode = ''", tc.stmt} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("mysql %s: %v", stmt, err)
			}
		}
		var srvMode string
		if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&srvMode); err != nil {
			t.Fatalf("mysql SELECT @@sql_mode: %v", err)
		}
		execAlterClauses(t, castCat, "SET SESSION sql_mode = ''; "+tc.stmt+"; SET @m = @@sql_mode;")
		got := castCat.userVars["m"].str
		if got != srvMode || strings.Contains(srvMode, "ANSI") != (tc.want == "ANSI") {
			t.Fatalf("%s: catalog %q, mysql %q, want ANSI=%v", tc.stmt, got, srvMode, tc.want == "ANSI")
		}
	}
	// AND and OR stop at a deciding left operand; the right is not evaluated.
	if _, err := conn.ExecContext(ctx, "SET @x_max = 9223372036854775807"); err != nil {
		t.Fatalf("mysql SET @x_max: %v", err)
	}
	execAlterClauses(t, castCat, "SET @x_max = 9223372036854775807;")
	for _, tc := range []struct {
		stmt string
		code int
	}{
		{"SET SESSION sql_mode = 0 AND (@x_max + 1)", 0},
		{"SET SESSION sql_mode = 1 OR (@x_max + 1)", 0},
		{"SET SESSION sql_mode = 1 AND (@x_max + 1)", 1690},
		{"SET SESSION sql_mode = NULL AND (@x_max + 1)", 1690},
		{"SET SESSION sql_mode = 0 OR (@x_max + 1)", 1690},
	} {
		for _, stmt := range []string{"SET SESSION sql_mode = 'ANSI'"} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("mysql %s: %v", stmt, err)
			}
			execAlterClauses(t, castCat, stmt+";")
		}
		var srvCode int
		if _, err := conn.ExecContext(ctx, tc.stmt); err != nil {
			srvCode, _ = simulationEngineRefusal(t, tc.stmt, err)
		}
		catCode, _ := simulationCatalogRefusal(t, castCat, tc.stmt)
		if srvCode != tc.code || catCode != tc.code {
			t.Fatalf("%s: mysql %d, catalog %d, want %d", tc.stmt, srvCode, catCode, tc.code)
		}
		var srvMode string
		if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&srvMode); err != nil {
			t.Fatalf("mysql SELECT @@sql_mode: %v", err)
		}
		execAlterClauses(t, castCat, "SET @m = @@sql_mode;")
		if got := castCat.userVars["m"].str; got != srvMode {
			t.Fatalf("%s: catalog sql_mode %q, mysql %q", tc.stmt, got, srvMode)
		}
	}

	// A CAST result longer than max_allowed_packet (64 MiB by default) is NULL;
	// a CHAR(n) whose widest bytes pass 2^32-1 is 1690 first.
	if got, _ := simulationCatalogRefusal(t, castCat, "SET @b_cast = CAST('a' AS CHAR(4294967295))"); got != 1690 {
		t.Fatalf("catalog CAST('a' AS CHAR(4294967295)) under utf8mb4 got %d, want 1690", got)
	}
	if _, err := conn.ExecContext(ctx, "SET @b_cast = CAST('a' AS CHAR(4294967295))"); err == nil {
		t.Fatal("mysql CAST('a' AS CHAR(4294967295)) under utf8mb4 succeeded, want 1690")
	} else if code, _ := simulationEngineRefusal(t, "SET @b_cast = CAST('a' AS CHAR(4294967295))", err); code != 1690 {
		t.Fatalf("mysql CAST('a' AS CHAR(4294967295)) got %d, want 1690", code)
	}
	for _, expr := range []string{
		"CAST('a' AS BINARY(4294967295))", "CAST('a' AS BINARY(67108865))",
		"CAST('a' AS CHAR(67108865))", "CAST('a' AS CHAR(67108864))",
		"CAST('a' AS CHAR(4294967295) CHARACTER SET latin1)",
	} {
		stmt := "SET @b_cast = " + expr
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		var srvNull bool
		if err := conn.QueryRowContext(ctx, "SELECT @b_cast IS NULL").Scan(&srvNull); err != nil {
			t.Fatalf("mysql read @b_cast after %s: %v", stmt, err)
		}
		execAlterClauses(t, castCat, stmt+";")
		if got := castCat.userVars["b_cast"].kind == userVarNull; got != srvNull {
			t.Fatalf("%s: catalog NULL=%v, mysql NULL=%v", stmt, got, srvNull)
		}
	}

	// The numeric functions' result types decide what sql_mode takes: an INT
	// is a mode mask, a DECIMAL or a double is 1232, even ROUND(3.6) = 4.
	for _, expr := range []string{
		"FLOOR(3.5)", "CEIL(3.2)", "ROUND(3.6)", "ROUND(4)", "ROUND(3.6e0)", "TRUNCATE(3.6, 0)",
		"TRUNCATE(4, 1)", "ROUND(14, -1)", "SIGN(-2.5)", "SIGN(2.5)", "FLOOR('3.7')", "ABS(-4.0)",
		"FLOOR(99999999999999999999.5)", "ROUND('4')", "MOD(5, 4)", "MOD(5.5, 2)", "~-2", "~0", "BIT_COUNT(1)",
		"ASCII(X'01')", "ORD(X'01')", "BIT_LENGTH('a')", "LOCATE('a', 'a')", "(CRC32('a') > 0)",
		"BIT_COUNT(CAST(X'01' AS BINARY) | CAST(X'02' AS BINARY))", "(CAST(2.5 AS DOUBLE) = 2.5)",
		"(ABS(CAST(-9223372036854775808 AS DECIMAL(20,0))) > 0)", "(CONCAT(1.5e0) = '1.5')", "QUARTER('2008-01-01')",
		"IFNULL(QUARTER('2008-02-30'), 1)", "IF(_binary'a' = 'A', 'ANSI', '')", "IF(_utf8mb4'a' = 'A', 'ANSI', '')", "((1 << 63) > 0)", "IF(IF(1, REPLACE('ansi', X'78', ''), '') = 'ANSI', 'ANSI', '')",
	} {
		stmt := "SET SESSION sql_mode = " + expr
		for _, pre := range []string{"SET SESSION sql_mode = ''"} {
			if _, err := conn.ExecContext(ctx, pre); err != nil {
				t.Fatalf("mysql %s: %v", pre, err)
			}
			execAlterClauses(t, castCat, pre+";")
		}
		var srvCode int
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			srvCode, _ = simulationEngineRefusal(t, stmt, err)
		}
		catCode, _ := simulationCatalogRefusal(t, castCat, stmt)
		if srvCode != catCode {
			t.Fatalf("%s: mysql %d, catalog %d", stmt, srvCode, catCode)
		}
		var srvMode string
		if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&srvMode); err != nil {
			t.Fatalf("mysql SELECT @@sql_mode: %v", err)
		}
		execAlterClauses(t, castCat, "SET @m = @@sql_mode;")
		if got := castCat.userVars["m"].str; got != srvMode {
			t.Fatalf("%s: catalog sql_mode %q, mysql %q", stmt, got, srvMode)
		}
	}
	// A user variable keeps ROUND's DECIMAL type.
	for _, stmt := range []string{"SET SESSION sql_mode = ''", "SET @m_round = ROUND(3.6)"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, castCat, stmt+";")
	}
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = @m_round"); err == nil {
		t.Fatal("mysql SET sql_mode = @m_round (ROUND(3.6)) succeeded, want 1232")
	} else if code, _ := simulationEngineRefusal(t, "SET SESSION sql_mode = @m_round", err); code != 1232 {
		t.Fatalf("mysql SET sql_mode = @m_round got %d, want 1232", code)
	}
	if got, _ := simulationCatalogRefusal(t, castCat, "SET SESSION sql_mode = @m_round"); got != 1232 {
		t.Fatalf("catalog SET sql_mode = @m_round got %d, want 1232", got)
	}

	// ABS of the smallest BIGINT is 1690, and the whole multi-assignment SET
	// fails with it, foreign_key_checks included.
	absOverflow := "SET foreign_key_checks = 0, @x_abs = ABS(-9223372036854775808)"
	if _, err := conn.ExecContext(ctx, absOverflow); err == nil {
		t.Fatalf("mysql %s succeeded, want 1690", absOverflow)
	} else if code, _ := simulationEngineRefusal(t, absOverflow, err); code != 1690 {
		t.Fatalf("mysql %s got %d, want 1690", absOverflow, code)
	}
	if got, _ := simulationCatalogRefusal(t, castCat, absOverflow); got != 1690 {
		t.Fatalf("catalog %s got %d, want 1690", absOverflow, got)
	}
	var srvFK string
	if err := conn.QueryRowContext(ctx, "SELECT CAST(@@foreign_key_checks AS CHAR)").Scan(&srvFK); err != nil {
		t.Fatalf("mysql read foreign_key_checks: %v", err)
	}
	execAlterClauses(t, castCat, "SET @fk = @@foreign_key_checks;")
	if got := castCat.userVars["fk"].str; got != srvFK || srvFK != "1" {
		t.Fatalf("foreign_key_checks after the failed SET: catalog %q, mysql %q, want 1", got, srvFK)
	}

	// A system variable the catalog does not model keeps its scope: reading a
	// global-only one at session scope, or a session-only one at global scope,
	// is 1238 and fails the whole SET.
	for _, stmt := range []string{
		"SET foreign_key_checks = 0, @x_scope = @@SESSION.max_connections",
		"SET foreign_key_checks = 0, @x_scope = @@GLOBAL.warning_count",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err == nil {
			t.Fatalf("mysql %s succeeded, want 1238", stmt)
		} else if code, _ := simulationEngineRefusal(t, stmt, err); code != 1238 {
			t.Fatalf("mysql %s got %d, want 1238", stmt, code)
		}
		if got, _ := simulationCatalogRefusal(t, castCat, stmt); got != 1238 {
			t.Fatalf("catalog %s got %d, want 1238", stmt, got)
		}
		var fk string
		if err := conn.QueryRowContext(ctx, "SELECT CAST(@@foreign_key_checks AS CHAR)").Scan(&fk); err != nil {
			t.Fatalf("mysql read foreign_key_checks: %v", err)
		}
		execAlterClauses(t, castCat, "SET @fk = @@foreign_key_checks;")
		if got := castCat.userVars["fk"].str; got != fk || fk != "1" {
			t.Fatalf("%s: foreign_key_checks catalog %q, mysql %q, want 1", stmt, got, fk)
		}
	}
	for _, stmt := range []string{"SET @x_scope = @@GLOBAL.max_connections", "SET @x_scope = @@SESSION.warning_count", "SET @x_scope = @@max_connections"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		if got, _ := simulationCatalogRefusal(t, castCat, stmt); got != 0 {
			t.Fatalf("catalog %s got %d, want 0", stmt, got)
		}
	}

	// PAD SPACE collations ignore trailing spaces; the _0900_ ones do not.
	for _, coll := range []string{"utf8mb4_0900_ai_ci", "utf8mb4_general_ci", "utf8mb4_bin", "utf8mb4_0900_bin", "latin1_swedish_ci", "utf8mb4_unicode_ci", "binary"} {
		for _, stmt := range []string{"SET collation_connection = " + coll, "SET SESSION sql_mode = ''", "SET SESSION sql_mode = IF('a' = 'a ', 'ANSI', '')"} {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("mysql %s: %v", stmt, err)
			}
			execAlterClauses(t, castCat, stmt+";")
		}
		var padMode string
		if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&padMode); err != nil {
			t.Fatalf("mysql SELECT @@sql_mode: %v", err)
		}
		execAlterClauses(t, castCat, "SET @m = @@sql_mode;")
		if got := castCat.userVars["m"].str; got != padMode {
			t.Fatalf("'a' = 'a ' under %s: catalog sql_mode %q, mysql %q", coll, got, padMode)
		}
	}
	for _, stmt := range []string{"SET NAMES utf8mb4", "SET SESSION sql_mode = ''"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, castCat, stmt+";")
	}

	// A scope keyword carries to the assignments after it: SET GLOBAL
	// foreign_key_checks = 1, sql_mode = 'ANSI' sets the global sql_mode and
	// leaves the session's alone.
	t.Cleanup(func() { _, _ = conn.ExecContext(ctx, "SET GLOBAL sql_mode = DEFAULT") })
	for _, stmt := range []string{"SET SESSION sql_mode = ''", "SET GLOBAL foreign_key_checks = 1, sql_mode = 'ANSI'"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, castCat, stmt+";")
	}
	var scopeMode string
	if err := conn.QueryRowContext(ctx, "SELECT @@SESSION.sql_mode").Scan(&scopeMode); err != nil {
		t.Fatalf("mysql SELECT @@SESSION.sql_mode: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "SET GLOBAL sql_mode = DEFAULT"); err != nil {
		t.Fatalf("mysql reset global sql_mode: %v", err)
	}
	execAlterClauses(t, castCat, "SET @m = @@SESSION.sql_mode; SET GLOBAL sql_mode = DEFAULT;")
	if got := castCat.userVars["m"].str; got != scopeMode || scopeMode != "" {
		t.Fatalf("SET GLOBAL ..., sql_mode = 'ANSI': catalog session sql_mode %q, mysql %q, want ''", got, scopeMode)
	}

	// A table's general tablespace is diffed and migrated: moving a table
	// from one tablespace to another and creating one in a tablespace both
	// run on the server and read back with the target's tablespace.
	var curDB string
	if err := conn.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&curDB); err != nil {
		t.Fatalf("mysql SELECT DATABASE(): %v", err)
	}
	t.Cleanup(func() {
		for _, stmt := range []string{"DROP TABLE IF EXISTS m_ts", "DROP TABLE IF EXISTS n_ts", "DROP TABLESPACE ts_mig1", "DROP TABLESPACE ts_mig2"} {
			_, _ = conn.ExecContext(ctx, stmt)
		}
	})
	for _, stmt := range []string{"CREATE TABLESPACE ts_mig1 ADD DATAFILE 'ts_mig1.ibd'", "CREATE TABLESPACE ts_mig2 ADD DATAFILE 'ts_mig2.ibd'", "CREATE TABLE m_ts (id INT) TABLESPACE ts_mig1"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	tsBase := "CREATE DATABASE `" + curDB + "`; USE `" + curDB + "`; CREATE TABLESPACE ts_mig1 ADD DATAFILE 'ts_mig1.ibd'; CREATE TABLESPACE ts_mig2 ADD DATAFILE 'ts_mig2.ibd';"
	fromCat, toCat := New(), New()
	execAlterClauses(t, fromCat, tsBase+" CREATE TABLE m_ts (id INT) TABLESPACE ts_mig1;")
	execAlterClauses(t, toCat, tsBase+" CREATE TABLE m_ts (id INT) TABLESPACE ts_mig2; CREATE TABLE n_ts (id INT) TABLESPACE ts_mig2;")
	plan := GenerateMigration(fromCat, toCat, Diff(fromCat, toCat))
	if len(plan.Ops) == 0 {
		t.Fatal("tablespace move and new table produced no migration")
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("mysql migration %s: %v", op.SQL, err)
		}
	}
	for _, tbl := range []string{"m_ts", "n_ts"} {
		var name, create string
		if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE "+tbl).Scan(&name, &create); err != nil {
			t.Fatalf("mysql SHOW CREATE TABLE %s: %v", tbl, err)
		}
		if !strings.Contains(create, "TABLESPACE `ts_mig2`") {
			t.Fatalf("migrated %s is not in ts_mig2:\n%s\nplan:\n%s", tbl, create, plan.SQL())
		}
	}

	// NOT 0 chains into sql_mode as mode 1, REAL_AS_FLOAT.
	for _, stmt := range []string{"SET SESSION sql_mode = ''", "SET SESSION sql_mode = NOT 0"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	var notMode string
	if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&notMode); err != nil {
		t.Fatalf("mysql SELECT @@sql_mode: %v", err)
	}
	execAlterClauses(t, castCat, "SET SESSION sql_mode = ''; SET SESSION sql_mode = NOT 0; SET @m = @@sql_mode;")
	if got := castCat.userVars["m"].str; got != notMode || notMode != "REAL_AS_FLOAT" {
		t.Fatalf("SET sql_mode = NOT 0: catalog %q, mysql %q, want REAL_AS_FLOAT", got, notMode)
	}

	// ABS(-4) chains into sql_mode as mode 4, ANSI_QUOTES.
	for _, stmt := range []string{"SET SESSION sql_mode = ''", "SET @m_abs = ABS(-4)", "SET SESSION sql_mode = @m_abs"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
	}
	var absMode string
	if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&absMode); err != nil {
		t.Fatalf("mysql SELECT @@sql_mode: %v", err)
	}
	execAlterClauses(t, castCat, "SET SESSION sql_mode = ''; SET @m_abs = ABS(-4); SET SESSION sql_mode = @m_abs; SET @m = @@sql_mode;")
	if got := castCat.userVars["m"].str; got != absMode || absMode != "ANSI_QUOTES" {
		t.Fatalf("SET sql_mode = ABS(-4): catalog %q, mysql %q, want ANSI_QUOTES", got, absMode)
	}

	// A TIMESTAMP default written with an introducer follows the session time
	// zone like any other: stored under +00:00, it shows 15:00:00 under +03:00.
	reset()
	introTSCat := New()
	execAlterClauses(t, introTSCat, "CREATE DATABASE d; USE d;")
	t.Cleanup(func() { _, _ = conn.ExecContext(ctx, "SET SESSION time_zone = DEFAULT") })
	for _, stmt := range []string{"SET SESSION time_zone = '+00:00'", "CREATE TABLE z (ts TIMESTAMP DEFAULT _utf8mb4'2026-01-01 12:00:00', tn TIMESTAMP DEFAULT N'2026-01-01 12:00:00')", "SET SESSION time_zone = '+03:00'"} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("mysql %s: %v", stmt, err)
		}
		execAlterClauses(t, introTSCat, stmt+";")
	}
	var tsName, tsCreate string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE z").Scan(&tsName, &tsCreate); err != nil {
		t.Fatalf("mysql SHOW CREATE TABLE z: %v", err)
	}
	introCreate := introTSCat.ShowCreateTable("d", "z")
	for _, want := range []string{"`ts` timestamp NULL DEFAULT '2026-01-01 15:00:00'", "`tn` timestamp NULL DEFAULT '2026-01-01 15:00:00'"} {
		if !strings.Contains(tsCreate, want) || !strings.Contains(introCreate, want) {
			t.Fatalf("introduced TIMESTAMP default under +03:00: want %s\nmysql:\n%s\ncatalog:\n%s", want, tsCreate, introCreate)
		}
	}
	if _, err := conn.ExecContext(ctx, "SET SESSION time_zone = DEFAULT"); err != nil {
		t.Fatalf("mysql reset time_zone: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = CAST(3.6 AS SIGNED)"); err != nil {
		t.Fatalf("mysql SET sql_mode = CAST(3.6 AS SIGNED): %v", err)
	}
	var castMode string
	if err := conn.QueryRowContext(ctx, "SELECT @@sql_mode").Scan(&castMode); err != nil {
		t.Fatalf("mysql SELECT @@sql_mode: %v", err)
	}
	execAlterClauses(t, castCat, "SET SESSION sql_mode = CAST(3.6 AS SIGNED); SET @m = @@sql_mode;")
	if got := castCat.userVars["m"].str; got != castMode || castMode != "ANSI_QUOTES" {
		t.Fatalf("SET sql_mode = CAST(3.6 AS SIGNED): catalog %q, mysql %q, want ANSI_QUOTES", got, castMode)
	}
	if _, err := conn.ExecContext(ctx, "SET SESSION sql_mode = DEFAULT"); err != nil {
		t.Fatalf("mysql reset sql_mode: %v", err)
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

	// A partition that records ENGINE=MyISAM is rebuilt in the ALTER's target
	// engine: its missing tablespace is judged (3510), and a conversion that
	// succeeds records the new engine on the partition.
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_myisam_rec (b INT) ENGINE=MyISAM PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) ENGINE = MyISAM TABLESPACE = ts_gone)")); err != nil {
		t.Fatalf("DefineTable partition recording MyISAM: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, c, "ALTER TABLE p_myisam_rec ENGINE=InnoDB"); got != 3510 {
		t.Fatalf("ALTER partition recording MyISAM with a missing tablespace to InnoDB got %d, want 3510", got)
	}
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE p_myisam_rec2 (b INT) ENGINE=MyISAM PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) ENGINE = MyISAM)")); err != nil {
		t.Fatalf("DefineTable partition recording MyISAM: %v", err)
	}
	execAlterClauses(t, c, "ALTER TABLE p_myisam_rec2 ENGINE=InnoDB;")
	if got := c.ShowCreateTable(c.currentDB, "p_myisam_rec2"); strings.Contains(got, "MyISAM") {
		t.Fatalf("converted partitioned table still records MyISAM:\n%s", got)
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

	// Configured system_time_zone governs session time_zone = 'SYSTEM' in both
	// TIMESTAMP UTC bounds and offset-bearing DATETIME normalization; unresolved
	// abbreviations like "EDT" are rejected rather than mapped to UTC, while an
	// explicitly preserved host Location handles abbreviations and DST.
	sysTZCat := New()
	sysTZCat.SetSystemTimeZone("+14:00")
	execAlterClauses(t, sysTZCat, "CREATE DATABASE d; USE d;")
	if got, _ := simulationCatalogRefusal(t, sysTZCat, "CREATE TABLE ts_bad (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')"); got != 1067 {
		t.Fatalf("SYSTEM=+14:00 TIMESTAMP lower bound got %d, want 1067", got)
	}
	execAlterClauses(t, sysTZCat, "CREATE TABLE ts_ok (ts TIMESTAMP DEFAULT '1970-01-01 14:00:01', dt DATETIME DEFAULT '2024-01-01 00:00:00+00:00');")
	if gotShow := sysTZCat.ShowCreateTable("d", "ts_ok"); !strings.Contains(gotShow, "`dt` datetime DEFAULT '2024-01-01 14:00:00'") {
		t.Fatalf("SYSTEM=+14:00 offset DATETIME normalization mismatch:\n%s", gotShow)
	}
	sysTZCat.SetSystemTimeZone("EDT")
	if got, _ := simulationCatalogRefusal(t, sysTZCat, "CREATE TABLE ts_unres (ts TIMESTAMP DEFAULT '1970-01-01 00:00:01')"); got != 1298 {
		t.Fatalf("unresolved SYSTEM='EDT' TIMESTAMP default got %d, want 1298", got)
	}
	if got, _ := simulationCatalogRefusal(t, sysTZCat, "CREATE TABLE dt_unres (dt DATETIME DEFAULT '2024-01-01 00:00:00+00:00')"); got != 1298 {
		t.Fatalf("unresolved SYSTEM='EDT' offset DATETIME default got %d, want 1298", got)
	}
	nyLoc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("LoadLocation(America/New_York): %v", err)
	}
	sysTZCat.SetSystemTimeZone("EST")
	sysTZCat.SetSystemTimeZoneLocation(nyLoc)
	if sysTZCat.SystemTimeZone() != "EST" {
		t.Fatalf("SystemTimeZone() = %q, want EST", sysTZCat.SystemTimeZone())
	}
	execAlterClauses(t, sysTZCat, "CREATE TABLE ts_ny (ts TIMESTAMP DEFAULT '1970-01-01 00:00:01', dt DATETIME DEFAULT '2026-06-01 12:00:00+00:00');")
	if gotShow := sysTZCat.ShowCreateTable("d", "ts_ny"); !strings.Contains(gotShow, "`ts` timestamp NULL DEFAULT '1970-01-01 00:00:01'") || !strings.Contains(gotShow, "`dt` datetime DEFAULT '2026-06-01 08:00:00'") {
		t.Fatalf("preserved host Location SYSTEM=America/New_York normalization mismatch:\n%s", gotShow)
	}

	// innodb_directories is read-only at runtime via SQL (1238), while Go
	// configuration setters remain available to configure external directories.
	dirCat := New()
	for _, sql := range []string{
		"SET innodb_directories = '/opt/innodb1'",
		"SET GLOBAL innodb_directories = '/opt/innodb1;/opt/innodb2'",
		"SET PERSIST innodb_directories = '/opt/innodb1'",
	} {
		if got, _ := simulationCatalogRefusal(t, dirCat, sql); got != 1238 {
			t.Fatalf("%s got %d, want 1238", sql, got)
		}
	}
	dirCat.AddTablespaceDir("/srv/mysql")
	dirCat.AddTablespaceDir("/opt/innodb1;/opt/innodb2")
	execAlterClauses(t, dirCat, "CREATE TABLESPACE ts_srv ADD DATAFILE '/srv/mysql/ts_srv.ibd' ENGINE=InnoDB;")
	execAlterClauses(t, dirCat, "CREATE TABLESPACE ts_opt ADD DATAFILE '/opt/innodb2/ts_opt.ibd' ENGINE=InnoDB;")
	if got, _ := simulationCatalogRefusal(t, dirCat, "CREATE TABLESPACE ts_unconf ADD DATAFILE '/unconfigured/ts_unconf.ibd' ENGINE=InnoDB"); got != 3121 {
		t.Fatalf("unconfigured tablespace directory got %d, want 3121", got)
	}

	// Replacing DataDir removes the previous automatic /var/lib/mysql allowance
	// unless /var/lib/mysql was explicitly added via AddTablespaceDir.
	replacedDirCat := New()
	replacedDirCat.SetDataDir("/srv/mysql")
	if got, _ := simulationCatalogRefusal(t, replacedDirCat, "CREATE TABLESPACE ts_old_dd ADD DATAFILE '/var/lib/mysql/ts_old_dd.ibd' ENGINE=InnoDB"); got != 3121 {
		t.Fatalf("replaced DataDir still allowed previous /var/lib/mysql datadir: got %d, want 3121", got)
	}
	execAlterClauses(t, replacedDirCat, "CREATE TABLESPACE ts_new_dd ADD DATAFILE '/srv/mysql/ts_new_dd.ibd' ENGINE=InnoDB;")
	replacedDirCat.AddTablespaceDir("/var/lib/mysql")
	replacedDirCat.SetDataDir("/srv/mysql2")
	execAlterClauses(t, replacedDirCat, "CREATE TABLESPACE ts_explicit_old ADD DATAFILE '/var/lib/mysql/ts_explicit_old.ibd' ENGINE=InnoDB;")
	if got, _ := simulationCatalogRefusal(t, replacedDirCat, "CREATE TABLESPACE ts_prev_dd2 ADD DATAFILE '/srv/mysql/ts_prev_dd2.ibd' ENGINE=InnoDB"); got != 3121 {
		t.Fatalf("second SetDataDir still allowed /srv/mysql: got %d, want 3121", got)
	}

	// Assigning an unsupported RHS to an existing user variable preserves the
	// prior value instead of deleting it, and ALTER TABLESPACE ENCRYPTION=''
	// rejects with 3184.
	execAlterClauses(t, c, "SET @keep_tz = '+08:00'; SET @keep_tz = SHA1('x'); SET time_zone = @keep_tz;")
	if c.session.TimeZone != "+08:00" {
		t.Fatalf("unsupported user-var RHS overwrote @keep_tz: session.TimeZone = %q, want +08:00", c.session.TimeZone)
	}
	execAlterClauses(t, c, "SET @mode = 2 + 2; SET sql_mode = @mode;")
	if !strings.Contains(c.session.SQLMode, "ANSI_QUOTES") {
		t.Fatalf("arithmetic @mode=2+2 did not enable ANSI_QUOTES: %q", c.session.SQLMode)
	}
	execAlterClauses(t, c, "SET SESSION sql_mode = DEFAULT; SET SESSION time_zone = DEFAULT;")
	execAlterClauses(t, c, "CREATE TABLESPACE ts_enc_empty ENGINE=InnoDB;")
	if got, _ := simulationCatalogRefusal(t, c, "ALTER TABLESPACE ts_enc_empty ENCRYPTION=''"); got != 3184 {
		t.Fatalf("ALTER TABLESPACE ENCRYPTION='' got %d, want 3184", got)
	}

	// renderWithSessionContext save/restore framing executes cleanly on a fresh
	// catalog whose session sql_mode starts at "DEFAULT".
	framedCat := New()
	execAlterClauses(t, framedCat, "CREATE DATABASE d; USE d; CREATE TABLE t (id INT);")
	framedSQL := renderWithSessionContext("CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = 1", SessionContext{SQLMode: ""}, false)
	execAlterClauses(t, framedCat, framedSQL+";")
	if !framedCat.isStrictSQLMode() {
		t.Fatalf("renderWithSessionContext did not restore default strict sql_mode: %q", framedCat.session.SQLMode)
	}
	execAlterClauses(t, framedCat, "SET @saved_cs = @@character_set_client; SET character_set_client = ascii; SET character_set_client = @saved_cs;")
	if framedCat.effectiveClientCharset() != "utf8mb4" {
		t.Fatalf("@@character_set_client round-trip got %q, want utf8mb4", framedCat.effectiveClientCharset())
	}

	// Configurable innodb_page_size allows matching FILE_BLOCK_SIZE (e.g. 32K)
	// and defaults tablespaces to that configured page size.
	pgCat := New()
	pgCat.SetInnoDBPageSize(32768)
	if pgCat.InnoDBPageSize() != 32768 {
		t.Fatalf("InnoDBPageSize() = %d, want 32768", pgCat.InnoDBPageSize())
	}
	if got, _ := simulationCatalogRefusal(t, pgCat, "SET innodb_page_size = 32768"); got != 1238 {
		t.Fatalf("SET innodb_page_size got %d, want 1238", got)
	}
	execAlterClauses(t, pgCat, "CREATE DATABASE d; USE d; CREATE TABLESPACE ts32 ADD DATAFILE 'ts32.ibd' FILE_BLOCK_SIZE=32768; CREATE TABLE t32 (id INT) TABLESPACE ts32;")
	if err := pgCat.createTablespace(&ast.CreateTablespaceStmt{
		Name:            "ts_ast_bad_ea",
		DataFile:        "ts_ast_bad_ea.ibd",
		HasDataFile:     true,
		EngineAttribute: "not-json",
	}); err == nil {
		t.Fatal("direct createTablespace with invalid ENGINE_ATTRIBUTE succeeded, want 3980")
	}

	// Configured 4 KiB innodb_page_size enforces the 768-byte InnoDB key prefix cap.
	pg4k := New()
	pg4k.SetInnoDBPageSize(4096)
	execAlterClauses(t, pg4k, "CREATE DATABASE d; USE d;")
	if got, _ := simulationCatalogRefusal(t, pg4k, "CREATE TABLE t4k_bad (c TEXT CHARACTER SET ascii, KEY(c(769)))"); got != 1071 {
		t.Fatalf("4K page size KEY(c(769)) got %d, want 1071", got)
	}
	execAlterClauses(t, pg4k, "CREATE TABLE t4k_ok (c TEXT CHARACTER SET ascii, KEY(c(768)));")

	// Global vs session sysvar separation across bool and charset variables.
	svCat := New()
	execAlterClauses(t, svCat, "SET GLOBAL foreign_key_checks = 0; SET GLOBAL character_set_client = latin1; SET @g_cs = @@GLOBAL.character_set_client; SET @s_cs = @@SESSION.character_set_client; SET @g_fk = @@GLOBAL.foreign_key_checks; SET @s_fk = @@SESSION.foreign_key_checks;")
	if !svCat.ForeignKeyChecks() || svCat.effectiveClientCharset() != "utf8mb4" {
		t.Fatalf("SET GLOBAL mutated session state: fk=%v cs=%q", svCat.ForeignKeyChecks(), svCat.effectiveClientCharset())
	}
	if svCat.userVars["g_cs"].str != "latin1" || svCat.userVars["s_cs"].str != "utf8mb4" || svCat.userVars["g_fk"].str != "0" || svCat.userVars["s_fk"].str != "1" {
		t.Fatalf("unexpected @@GLOBAL/@@SESSION readback: %+v", svCat.userVars)
	}

	// SET GLOBAL/PERSIST/SESSION scope keywords carry into subsequent
	// assignments in the same SET statement while @@GLOBAL prefixes do not.
	stickyCat := New()
	execAlterClauses(t, stickyCat, "SET GLOBAL sql_mode='', @u=1, foreign_key_checks=0; SET @gf=@@GLOBAL.foreign_key_checks, @sf=@@SESSION.foreign_key_checks;")
	if stickyCat.userVars["gf"].str != "0" || stickyCat.userVars["sf"].str != "1" || !stickyCat.ForeignKeyChecks() {
		t.Fatalf("SET GLOBAL sticky scope mismatch: %+v, fk=%v", stickyCat.userVars, stickyCat.ForeignKeyChecks())
	}
	execAlterClauses(t, stickyCat, "SET GLOBAL foreign_key_checks=1, SESSION sql_mode='ANSI', show_gipk_in_create_table_and_information_schema=1;")
	if !stickyCat.globalForeignKeyChecks || !stickyCat.session.ShowGIPK || stickyCat.globalShowGIPK {
		t.Fatalf("SET GLOBAL ..., SESSION ... scope override mismatch: globalFK=%v sessionGIPK=%v globalGIPK=%v", stickyCat.globalForeignKeyChecks, stickyCat.session.ShowGIPK, stickyCat.globalShowGIPK)
	}
	execAlterClauses(t, stickyCat, "SET @@GLOBAL.sql_mode='STRICT_ALL_TABLES', foreign_key_checks=0; SET @gf2=@@GLOBAL.foreign_key_checks, @sf2=@@SESSION.foreign_key_checks;")
	if stickyCat.userVars["gf2"].str != "1" || stickyCat.userVars["sf2"].str != "0" || stickyCat.ForeignKeyChecks() {
		t.Fatalf("SET @@GLOBAL prefix leak: %+v, fk=%v", stickyCat.userVars, stickyCat.ForeignKeyChecks())
	}

	// General tablespace ENCRYPTION state persists across ALTER/RENAME/DROP and
	// validates table encryption compatibility (3825).
	encCat := New()
	execAlterClauses(t, encCat, "CREATE DATABASE d; USE d; CREATE TABLESPACE ts_y ADD DATAFILE 'ts_y.ibd' ENCRYPTION='Y';")
	if got, _ := simulationCatalogRefusal(t, encCat, "CREATE TABLE t_plain (id INT) TABLESPACE ts_y"); got != 3825 {
		t.Fatalf("unencrypted table in encrypted tablespace got %d, want 3825", got)
	}
	if got, _ := simulationCatalogRefusal(t, encCat, "CREATE TABLE t_n (id INT) TABLESPACE ts_y ENCRYPTION='N'"); got != 3825 {
		t.Fatalf("ENCRYPTION='N' table in encrypted tablespace got %d, want 3825", got)
	}
	execAlterClauses(t, encCat, "CREATE TABLE t_ok_y (id INT) TABLESPACE ts_y ENCRYPTION='Y';")
	execAlterClauses(t, encCat, "ALTER TABLESPACE ts_y RENAME TO ts_y2;")
	if got, _ := simulationCatalogRefusal(t, encCat, "CREATE TABLE t_n2 (id INT) TABLESPACE ts_y2"); got != 3825 {
		t.Fatalf("renamed encrypted tablespace did not keep encryption: got %d, want 3825", got)
	}
	execAlterClauses(t, encCat, "CREATE TABLESPACE ts_toggle ADD DATAFILE 'ts_toggle.ibd' ENCRYPTION='Y'; ALTER TABLESPACE ts_toggle ENCRYPTION='N';")
	if got, _ := simulationCatalogRefusal(t, encCat, "CREATE TABLE t_tog_y (id INT) TABLESPACE ts_toggle ENCRYPTION='Y'"); got != 3825 {
		t.Fatalf("altered unencrypted tablespace accepted ENCRYPTION='Y': got %d, want 3825", got)
	}
	execAlterClauses(t, encCat, "DROP TABLESPACE ts_toggle;")

	// ALTER TABLESPACE ... ENCRYPTION updates resident InnoDB tables (as verified
	// against live mysql:8.0.46 with keyring_file.so) so SHOW CREATE TABLE and
	// subsequent ALTER TABLE statements stay consistent.
	execAlterClauses(t, encCat, "CREATE TABLESPACE ts_prop ADD DATAFILE 'ts_prop.ibd'; CREATE TABLE t_prop (id INT PRIMARY KEY) TABLESPACE ts_prop; ALTER TABLESPACE ts_prop ENCRYPTION='Y'; ALTER TABLE t_prop ADD COLUMN c INT;")
	if gotShow := encCat.ShowCreateTable("d", "t_prop"); !strings.Contains(gotShow, "ENCRYPTION='Y'") {
		t.Fatalf("resident table missing ENCRYPTION='Y' after ALTER TABLESPACE ENCRYPTION='Y':\n%s", gotShow)
	}
	execAlterClauses(t, encCat, "ALTER TABLESPACE ts_prop ENCRYPTION='N'; ALTER TABLE t_prop ADD COLUMN c2 INT;")
	if gotShow := encCat.ShowCreateTable("d", "t_prop"); strings.Contains(gotShow, "ENCRYPTION=") {
		t.Fatalf("resident table still showed ENCRYPTION after ALTER TABLESPACE ENCRYPTION='N':\n%s", gotShow)
	}
	if got, _ := simulationCatalogRefusal(t, encCat, "ALTER TABLE t_prop ENCRYPTION='Y'"); got != 3825 {
		t.Fatalf("ALTER TABLE t_prop ENCRYPTION='Y' on decrypted tablespace got %d, want 3825", got)
	}

	// LoadSDL infers unknown general tablespace encryption from loaded tables
	// instead of assuming false.
	sdlEncCat, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TABLE t_sdl_y (id INT PRIMARY KEY) TABLESPACE ts_sdl_y ENCRYPTION='Y';")
	if err != nil {
		t.Fatalf("LoadSDL encrypted general tablespace failed: %v", err)
	}
	execAlterClauses(t, sdlEncCat, "ALTER TABLE t_sdl_y ADD COLUMN c INT;")
	if got, _ := simulationCatalogRefusal(t, sdlEncCat, "CREATE TABLE t_sdl_plain (id INT) TABLESPACE ts_sdl_y"); got != 3825 {
		t.Fatalf("table omitting ENCRYPTION in SDL-inferred encrypted tablespace got %d, want 3825", got)
	}
	if _, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TABLE t1 (id INT) TABLESPACE ts_s ENCRYPTION='Y'; CREATE TABLE t2 (id INT) TABLESPACE ts_s;"); err == nil {
		t.Fatal("LoadSDL accepted conflicting encryption states for the same general tablespace, want 3825")
	}

	// Deeply nested DEFAULT (ABS(...)) validates in linear time without quadratic
	// AST re-traversal.
	const absDepth = 2000
	deepExpr := strings.Repeat("ABS(", absDepth) + "1" + strings.Repeat(")", absDepth)
	deepCat := New()
	execAlterClauses(t, deepCat, "CREATE DATABASE d; USE d; CREATE TABLE t_deep (a INT DEFAULT ("+deepExpr+"));")

	// Named time zones do not consult host OS tzdata unless explicitly
	// configured on the catalog via SetNamedTimeZones/SetNamedTimeZone.
	unloadedTZ := New()
	if got, _ := simulationCatalogRefusal(t, unloadedTZ, "SET time_zone = 'America/New_York'"); got != 1298 {
		t.Fatalf("unloaded named time_zone got %d, want 1298", got)
	}
	unloadedTZ.SetNamedTimeZones("America/New_York")
	if got, _ := simulationCatalogRefusal(t, unloadedTZ, "SET time_zone = 'Europe/London'"); got != 1298 {
		t.Fatalf("unconfigured Europe/London got %d, want 1298", got)
	}
	execAlterClauses(t, unloadedTZ, "CREATE DATABASE d; USE d; SET time_zone = 'America/New_York'; CREATE TABLE t_tz (dt DATETIME DEFAULT '2026-06-01 12:00:00+00:00');")
	if gotShow := unloadedTZ.ShowCreateTable("d", "t_tz"); !strings.Contains(gotShow, "DEFAULT '2026-06-01 08:00:00'") {
		t.Fatalf("configured SetNamedTimeZones normalization mismatch:\n%s", gotShow)
	}

	// Configured 32 KiB and 64 KiB innodb_page_size scale AUTOEXTEND_SIZE minimum
	// and alignment step to 8M and 16M respectively (matching live mysql:8.0.46).
	ae32k := New()
	ae32k.SetInnoDBPageSize(32768)
	execAlterClauses(t, ae32k, "CREATE DATABASE d; USE d;")
	if got, _ := simulationCatalogRefusal(t, ae32k, "CREATE TABLESPACE ts32_4 ADD DATAFILE 'ts32_4.ibd' AUTOEXTEND_SIZE=4M"); got != 4025 {
		t.Fatalf("32K tablespace AUTOEXTEND_SIZE=4M got %d, want 4025", got)
	}
	if got, _ := simulationCatalogRefusal(t, ae32k, "CREATE TABLE t32_4 (id INT) AUTOEXTEND_SIZE=4M"); got != 4025 {
		t.Fatalf("32K table AUTOEXTEND_SIZE=4M got %d, want 4025", got)
	}
	if got, _ := simulationCatalogRefusal(t, ae32k, "CREATE TABLESPACE ts32_12 ADD DATAFILE 'ts32_12.ibd' AUTOEXTEND_SIZE=12M"); got != 4023 {
		t.Fatalf("32K tablespace AUTOEXTEND_SIZE=12M got %d, want 4023", got)
	}
	if got, _ := simulationCatalogRefusal(t, ae32k, "CREATE TABLE t32_12 (id INT) AUTOEXTEND_SIZE=12M"); got != 4023 {
		t.Fatalf("32K table AUTOEXTEND_SIZE=12M got %d, want 4023", got)
	}
	execAlterClauses(t, ae32k, "CREATE TABLESPACE ts32_8 ADD DATAFILE 'ts32_8.ibd' AUTOEXTEND_SIZE=8M; CREATE TABLE t32_8 (id INT) AUTOEXTEND_SIZE=8M; ALTER TABLE t32_8 AUTOEXTEND_SIZE=16M;")

	ae64k := New()
	ae64k.SetInnoDBPageSize(65536)
	execAlterClauses(t, ae64k, "CREATE DATABASE d; USE d;")
	for _, bad := range []string{"4M", "8M", "12M"} {
		if got, _ := simulationCatalogRefusal(t, ae64k, "CREATE TABLESPACE ts64_bad ADD DATAFILE 'ts64_bad.ibd' AUTOEXTEND_SIZE="+bad); got != 4025 {
			t.Fatalf("64K tablespace AUTOEXTEND_SIZE=%s got %d, want 4025", bad, got)
		}
		if got, _ := simulationCatalogRefusal(t, ae64k, "CREATE TABLE t64_bad (id INT) AUTOEXTEND_SIZE="+bad); got != 4025 {
			t.Fatalf("64K table AUTOEXTEND_SIZE=%s got %d, want 4025", bad, got)
		}
	}
	execAlterClauses(t, ae64k, "CREATE TABLESPACE ts64_16 ADD DATAFILE 'ts64_16.ibd' AUTOEXTEND_SIZE=16M; CREATE TABLE t64_16 (id INT) AUTOEXTEND_SIZE=16M;")

	// Decimal arithmetic in SET user variable assignments evaluates exact
	// decimal strings and division-by-zero NULL rather than dropping the update.
	decCat := New()
	execAlterClauses(t, decCat, "SET @a = 1.5 + 1, @b = 1.50 * 2, @c = 10.5 - 3, @d = 1 / 2, @e = 1.5 / 0, @f = 5 % 2.5, @g = -1.25 + 3;")
	if decCat.userVars["a"].str != "2.5" || decCat.userVars["b"].str != "3.00" || decCat.userVars["c"].str != "7.5" || decCat.userVars["d"].str != "0.5000" || decCat.userVars["e"].kind != userVarNull || decCat.userVars["g"].str != "1.75" {
		t.Fatalf("decimal SET arithmetic mismatch: %+v", decCat.userVars)
	}
	if got, _ := simulationCatalogRefusal(t, decCat, "SET sql_mode = @a"); got != 1232 {
		t.Fatalf("SET sql_mode = @a (2.5) got %d, want 1232", got)
	}

	// default_table_encryption separates GLOBAL and SESSION scopes and drives
	// omitted CREATE TABLESPACE ENCRYPTION state.
	dteCat := New()
	execAlterClauses(t, dteCat, "CREATE DATABASE d; USE d; SET GLOBAL default_table_encryption = ON; SET @gd=@@GLOBAL.default_table_encryption, @sd=@@SESSION.default_table_encryption;")
	if dteCat.userVars["gd"].str != "1" || dteCat.userVars["sd"].str != "0" {
		t.Fatalf("default_table_encryption scope leak: %+v", dteCat.userVars)
	}
	execAlterClauses(t, dteCat, "SET SESSION default_table_encryption = DEFAULT; CREATE TABLESPACE ts_dte ADD DATAFILE 'ts_dte.ibd';")
	if got, _ := simulationCatalogRefusal(t, dteCat, "CREATE TABLE t_dte_plain (id INT) TABLESPACE ts_dte"); got != 3825 {
		t.Fatalf("plain table in default-encrypted tablespace got %d, want 3825", got)
	}

	// Composite index key byte total scales with configured 4K InnoDB page size.
	comp4k := New()
	comp4k.SetInnoDBPageSize(4096)
	execAlterClauses(t, comp4k, "CREATE DATABASE d; USE d;")
	if got, _ := simulationCatalogRefusal(t, comp4k, "CREATE TABLE t4k_comp_bad (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(400), b(400)))"); got != 1071 {
		t.Fatalf("4K composite KEY(a(400),b(400)) got %d, want 1071", got)
	}
	execAlterClauses(t, comp4k, "CREATE TABLE t4k_comp_ok (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(400), b(368)));")

	// Configured 32K/64K InnoDB page size rejects file-per-table COMPRESSED and
	// KEY_BLOCK_SIZE tables (1031 on CREATE, 1478 on innodb_system and ALTER).
	for _, pgSize := range []uint64{32768, 65536} {
		compPG := New()
		compPG.SetInnoDBPageSize(pgSize)
		execAlterClauses(t, compPG, "CREATE DATABASE d; USE d;")
		if got, _ := simulationCatalogRefusal(t, compPG, "CREATE TABLE t_cmp (id INT) ROW_FORMAT=COMPRESSED"); got != 1031 {
			t.Fatalf("%d CREATE COMPRESSED got %d, want 1031", pgSize, got)
		}
		if got, _ := simulationCatalogRefusal(t, compPG, "CREATE TABLE t_fpt (id INT) TABLESPACE innodb_file_per_table ROW_FORMAT=COMPRESSED"); got != 1031 {
			t.Fatalf("%d CREATE file_per_table COMPRESSED got %d, want 1031", pgSize, got)
		}
		if got, _ := simulationCatalogRefusal(t, compPG, "CREATE TABLE t_kbs (id INT) KEY_BLOCK_SIZE=8"); got != 1031 {
			t.Fatalf("%d CREATE KEY_BLOCK_SIZE=8 got %d, want 1031", pgSize, got)
		}
		if got, _ := simulationCatalogRefusal(t, compPG, "CREATE TABLE t_sys (id INT) TABLESPACE innodb_system ROW_FORMAT=COMPRESSED"); got != 1478 {
			t.Fatalf("%d CREATE innodb_system COMPRESSED got %d, want 1478", pgSize, got)
		}
		execAlterClauses(t, compPG, "CREATE TABLE t_dyn (id INT) ROW_FORMAT=DYNAMIC; CREATE TABLE t_myisam (id INT) ENGINE=MyISAM KEY_BLOCK_SIZE=8;")
		if got, _ := simulationCatalogRefusal(t, compPG, "ALTER TABLE t_dyn ROW_FORMAT=COMPRESSED"); got != 1478 {
			t.Fatalf("%d ALTER COMPRESSED got %d, want 1478", pgSize, got)
		}
	}

	// Every builtin allowed in DEFAULT expressions has an explicit arity entry,
	// and every entry in defaultBuiltinFuncArity / defaultBuiltinFuncStep is allowed.
	for fn := range allowedDefaultExprBuiltinFuncs {
		if _, ok := defaultBuiltinFuncArity[fn]; !ok {
			t.Errorf("allowedDefaultExprBuiltinFuncs[%q] missing from defaultBuiltinFuncArity", fn)
		}
	}
	for fn := range defaultBuiltinFuncArity {
		if !allowedDefaultExprBuiltinFuncs[fn] {
			t.Errorf("defaultBuiltinFuncArity[%q] missing from allowedDefaultExprBuiltinFuncs", fn)
		}
	}
	for fn := range defaultBuiltinFuncStep {
		if !allowedDefaultExprBuiltinFuncs[fn] {
			t.Errorf("defaultBuiltinFuncStep[%q] missing from allowedDefaultExprBuiltinFuncs", fn)
		}
	}

	// A synthesized functional index exceeding 3072 bytes is rejected after
	// hidden generated column synthesis. Each part stays within 512 characters:
	// a longer one is a TEXT, which 3757 refuses first (TestSimulationCausesOracle).
	fnCat := New()
	execAlterClauses(t, fnCat, "CREATE DATABASE d; USE d;")
	if got, _ := simulationCatalogRefusal(t, fnCat, "CREATE TABLE t_fn_width (a VARCHAR(500), b VARCHAR(500), KEY ((LOWER(a)), (LOWER(b))))"); got != 1071 {
		t.Fatalf("functional KEY ((LOWER(VARCHAR(500))), (LOWER(VARCHAR(500)))) got %d, want 1071", got)
	}

	// SET IF(), CASE, TRIM_*, NULLIF, comparison/logical ops, and ContinueOnError
	// state propagation.
	ifCat := New()
	execAlterClauses(t, ifCat, "CREATE DATABASE d; USE d;")
	res1, err := ifCat.Exec("SET sql_mode = IF(1 = 1, 'ANSI', '');", &ExecOptions{ContinueOnError: true})
	if err != nil || len(res1) != 1 || res1[0].Error != nil {
		t.Fatalf("ifCat ANSI SET failed: res=%+v err=%v", res1, err)
	}
	if !strings.Contains(ifCat.session.SQLMode, "ANSI") {
		t.Fatalf("expected ANSI enabled, got %q", ifCat.session.SQLMode)
	}
	res2, err := ifCat.Exec(
		"SET sql_mode = IF('abc', 'ANSI', TRIM(LEADING 'x' FROM 'xxSTRICT_ALL_TABLES,NO_ZERO_DATE')); CREATE TABLE t_strict_def (d DATE DEFAULT '0000-00-00'); SET @nz = CASE WHEN 1 > 2 THEN 'a' ELSE NULLIF('b', 'c') END;",
		&ExecOptions{ContinueOnError: true},
	)
	if err != nil || len(res2) != 3 {
		t.Fatalf("ifCat.Exec chained: res=%+v err=%v", res2, err)
	}
	if res2[0].Error != nil || res2[1].Error == nil || res2[2].Error != nil {
		t.Fatalf("ifCat chained unexpected errors: [0]=%v [1]=%v [2]=%v", res2[0].Error, res2[1].Error, res2[2].Error)
	}
	if !strings.Contains(ifCat.session.SQLMode, "STRICT_ALL_TABLES") {
		t.Fatalf("sql_mode=%q, want STRICT_ALL_TABLES", ifCat.session.SQLMode)
	}
	if ifCat.userVars["nz"].str != "b" {
		t.Fatalf("userVars[nz]=%+v, want b", ifCat.userVars["nz"])
	}

	// #2/#6: Diff, GenerateMigration, and LoadMetadata preserve AUTOEXTEND_SIZE and SECONDARY_ENGINE.
	optFrom := New()
	execAlterClauses(t, optFrom, "CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY);")
	optTo := New()
	execAlterClauses(t, optTo, "CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY) /*!80023 AUTOEXTEND_SIZE=4M */ SECONDARY_ENGINE=RAPID;")
	optDiff := Diff(optFrom, optTo)
	if len(optDiff.Tables) != 1 || optDiff.Tables[0].Action != DiffModify {
		t.Fatalf("Diff missed AUTOEXTEND_SIZE/SECONDARY_ENGINE: %+v", optDiff.Tables)
	}
	optMig := GenerateMigration(optFrom, optTo, optDiff)
	if len(optMig.Ops) != 1 || !strings.Contains(optMig.Ops[0].SQL, "AUTOEXTEND_SIZE=4194304") || !strings.Contains(optMig.Ops[0].SQL, "SECONDARY_ENGINE=RAPID") {
		t.Fatalf("GenerateMigration alter missing options: %+v", optMig.Ops)
	}
	optAdd := GenerateMigration(New(), optTo, Diff(New(), optTo))
	foundCreateOpt := false
	for _, op := range optAdd.Ops {
		if op.Type == OpCreateTable && strings.Contains(op.SQL, "AUTOEXTEND_SIZE=4194304") && strings.Contains(op.SQL, "SECONDARY_ENGINE=RAPID") {
			foundCreateOpt = true
		}
	}
	if !foundCreateOpt {
		t.Fatalf("GenerateMigration CREATE TABLE missing options: %+v", optAdd.Ops)
	}

	// #3: Failed CREATE INDEX rolls back synthesized hidden column so that column name can be added afterwards.
	ciCat := New()
	execAlterClauses(t, ciCat, "CREATE DATABASE d; USE d; CREATE TABLE t (a VARCHAR(100) CHARACTER SET utf8mb4, b VARBINARY(2850));")
	if got, _ := simulationCatalogRefusal(t, ciCat, "CREATE INDEX orphan ON t (b, (LOWER(a)))"); got != 1071 {
		t.Fatalf("oversized CREATE INDEX got %d, want 1071", got)
	}
	execAlterClauses(t, ciCat, "ALTER TABLE t ADD COLUMN `!hidden!orphan!1!0` INT;")

	// #5: TIMESTAMP literal default renders shifted to current session time_zone in SHOW CREATE TABLE.
	tzTSCat := New()
	execAlterClauses(t, tzTSCat, "CREATE DATABASE d; USE d; SET time_zone = '+00:00'; CREATE TABLE ts_tz (ts TIMESTAMP DEFAULT '2026-01-01 12:00:00+00:00'); SET time_zone = '+03:00';")
	if show := tzTSCat.ShowCreateTable("d", "ts_tz"); !strings.Contains(show, "DEFAULT '2026-01-01 15:00:00'") {
		t.Fatalf("TIMESTAMP ShowCreateTable after TZ switch got:\n%s", show)
	}

	// #7: LoadSDL checks deferred AUTO_INCREMENT key requirement across final tables.
	if _, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TABLE bad_ai (a INT AUTO_INCREMENT);"); err == nil {
		t.Fatal("LoadSDL accepted unindexed AUTO_INCREMENT, want 1075")
	}
	if _, err := LoadSDL("CREATE DATABASE d; USE d; CREATE INDEX i ON ok_ai (a); CREATE TABLE ok_ai (a INT AUTO_INCREMENT);"); err != nil {
		t.Fatalf("LoadSDL rejected indexed AUTO_INCREMENT: %v", err)
	}

	// #9: LIKE and RENAME TABLE rebind source table qualifiers on defaults and functional indexes.
	likeRenCat := New()
	execAlterClauses(t, likeRenCat, "CREATE DATABASE d; USE d; CREATE TABLE src (a INT, b INT DEFAULT (src.a + 1)); CREATE TABLE dst LIKE src; ALTER TABLE dst ADD c INT; CREATE TABLE src_fn (a INT, KEY ((src_fn.a + 1))); CREATE TABLE dst_fn LIKE src_fn;")
	if got, _ := simulationCatalogRefusal(t, likeRenCat, "ALTER TABLE dst_fn DROP COLUMN a"); got != 3837 {
		t.Fatalf("dst_fn functional index column dependency got %d, want 3837", got)
	}
	execAlterClauses(t, likeRenCat, "RENAME TABLE src TO src2, src_fn TO src_fn2; ALTER TABLE src2 ADD c INT;")
	if got, _ := simulationCatalogRefusal(t, likeRenCat, "ALTER TABLE src_fn2 DROP COLUMN a"); got != 3837 {
		t.Fatalf("src_fn2 functional index column dependency got %d, want 3837", got)
	}

	// #18: SET GLOBAL sql_mode='' followed by SET SESSION sql_mode=DEFAULT inherits empty mode.
	gmCat := New()
	execAlterClauses(t, gmCat, "CREATE DATABASE d; USE d; SET GLOBAL sql_mode = ''; SET SESSION sql_mode = DEFAULT; CREATE TABLE t_ns (c ENUM('a', 'A')); SET SESSION sql_mode = @@GLOBAL.sql_mode; CREATE TABLE t_ns2 (c ENUM('b', 'B'));")

	// Session collation comparison in SET CASE/equality: case-insensitive default vs _bin.
	collCat := New()
	execAlterClauses(t, collCat, "SET sql_mode = CASE 'ansi' WHEN 'ANSI' THEN 'TRADITIONAL' ELSE '' END;")
	if !strings.Contains(collCat.session.SQLMode, "TRADITIONAL") {
		t.Fatalf("CASE under utf8mb4_0900_ai_ci got sql_mode=%q, want TRADITIONAL", collCat.session.SQLMode)
	}
	execAlterClauses(t, collCat, "SET collation_connection = utf8mb4_bin; SET sql_mode = CASE 'ansi' WHEN 'ANSI' THEN 'TRADITIONAL' ELSE '' END;")
	if collCat.session.SQLMode != "" {
		t.Fatalf("CASE under utf8mb4_bin got sql_mode=%q, want empty", collCat.session.SQLMode)
	}

	// SET mixed numeric-string comparison, ASCII-space TRIM, bytewise CHAR_LENGTH
	// on binary literals, and unselected-branch function/arity validation.
	setEdge := New()
	execAlterClauses(t, setEdge, "CREATE DATABASE d; USE d; SET sql_mode = IF(0 = 'abc', 'ANSI', '');")
	if !strings.Contains(setEdge.session.SQLMode, "ANSI") {
		t.Fatalf("IF(0 = 'abc', 'ANSI', '') got %q, want ANSI", setEdge.session.SQLMode)
	}
	execAlterClauses(t, setEdge, "SET sql_mode = IF(TRIM(X'09') = X'', 'ANSI', '');")
	if setEdge.session.SQLMode != "" {
		t.Fatalf("TRIM(X'09') should preserve tab, got sql_mode=%q", setEdge.session.SQLMode)
	}
	execAlterClauses(t, setEdge, "SET sql_mode = IF(CHAR_LENGTH(X'C3A9') = 2, 'ANSI', '');")
	if !strings.Contains(setEdge.session.SQLMode, "ANSI") {
		t.Fatalf("CHAR_LENGTH(X'C3A9') should be 2, got sql_mode=%q", setEdge.session.SQLMode)
	}
	if got, _ := simulationCatalogRefusal(t, setEdge, "SET sql_mode = IF(1, 'TRADITIONAL', d.no_such())"); got != 1305 {
		t.Fatalf("IF unselected missing func got %d, want 1305", got)
	}
	if !strings.Contains(setEdge.session.SQLMode, "ANSI") {
		t.Fatalf("failed IF mutated sql_mode to %q", setEdge.session.SQLMode)
	}
	if got, _ := simulationCatalogRefusal(t, setEdge, "SET sql_mode = IF(1, 'TRADITIONAL', CONCAT())"); got != 1582 {
		t.Fatalf("IF unselected bad arity got %d, want 1582", got)
	}

	// Implicit compressed block size derives half of innodb_page_size (4K/8K/16K).
	for _, tc := range []struct {
		pageSize uint64
		sql      string
	}{
		{4096, "CREATE DATABASE d; USE d; CREATE TABLESPACE ts_c ADD DATAFILE 'ts_c.ibd' FILE_BLOCK_SIZE=2K ENGINE=InnoDB; CREATE TABLE t_c (id INT PRIMARY KEY) TABLESPACE ts_c ROW_FORMAT=COMPRESSED;"},
		{8192, "CREATE DATABASE d; USE d; CREATE TABLESPACE ts_c ADD DATAFILE 'ts_c.ibd' FILE_BLOCK_SIZE=4K ENGINE=InnoDB; CREATE TABLE t_c (id INT PRIMARY KEY) TABLESPACE ts_c ROW_FORMAT=COMPRESSED;"},
	} {
		compCat := New()
		compCat.SetInnoDBPageSize(tc.pageSize)
		execAlterClauses(t, compCat, tc.sql)
	}

	// CAST in SET, skipped IFNULL/COALESCE validation, mixed IF result typing,
	// failed ALTER RENAME AST preservation, and LoadMetadata pre-existing tablespaces.
	castCat := New()
	execAlterClauses(t, castCat, "CREATE DATABASE d; USE d; SET sql_mode = CAST('ANSI' AS CHAR);")
	if !strings.Contains(castCat.session.SQLMode, "ANSI") {
		t.Fatalf("CAST('ANSI' AS CHAR) got %q, want ANSI", castCat.session.SQLMode)
	}
	if got, _ := simulationCatalogRefusal(t, castCat, "SET sql_mode = IFNULL('TRADITIONAL', d.no_such())"); got != 1305 {
		t.Fatalf("IFNULL skipped missing func got %d, want 1305", got)
	}
	if got, _ := simulationCatalogRefusal(t, castCat, "SET sql_mode = COALESCE(NULL, 'TRADITIONAL', d.no_such())"); got != 1305 {
		t.Fatalf("COALESCE skipped missing func got %d, want 1305", got)
	}
	execAlterClauses(t, castCat, "SET sql_mode = IF(IF(1, 1, 'x') = '01', 'ANSI', '');")
	if castCat.session.SQLMode != "" {
		t.Fatalf("mixed IF(1,1,'x') = '01' enabled ANSI: %q", castCat.session.SQLMode)
	}

	renFailCat := New()
	execAlterClauses(t, renFailCat, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT DEFAULT (t.a + 1), KEY ((t.a + 1))); CREATE TABLE u_def (c INT, d INT DEFAULT (c + 1));")
	if got, _ := simulationCatalogRefusal(t, renFailCat, "ALTER TABLE t RENAME TO u, ADD COLUMN a INT"); got != 1060 {
		t.Fatalf("failed RENAME+ADD COLUMN got %d, want 1060", got)
	}
	if got, _ := simulationCatalogRefusal(t, renFailCat, "ALTER TABLE u_def RENAME TO u_def2, RENAME COLUMN c TO c2, ADD COLUMN c INT"); got != 3766 {
		t.Fatalf("RENAME TABLE + RENAME COLUMN + ADD COLUMN bypassing default dependency got %d, want 3766", got)
	}
	execAlterClauses(t, renFailCat, "ALTER TABLE t RENAME TO u2; ALTER TABLE u2 ADD COLUMN c INT;")

	metaCat := New()
	metaCat.tablespacesPartial = true
	execAlterClauses(t, metaCat, "CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB; ALTER TABLE t TABLESPACE pre_existing_ts; CREATE TABLE u (id INT PRIMARY KEY) TABLESPACE pre_existing_ts ENGINE=InnoDB;")
	if err := metaCat.DefineTable(mustParseTable(t, "CREATE TABLE loaded_ts_tbl (id INT PRIMARY KEY) TABLESPACE loaded_ts ENGINE=InnoDB")); err != nil {
		t.Fatalf("DefineTable loaded_ts_tbl: %v", err)
	}
	if got, _ := simulationCatalogRefusal(t, metaCat, "CREATE TABLESPACE loaded_ts ADD DATAFILE 'loaded_ts.ibd' ENGINE=InnoDB"); got != 1813 {
		t.Fatalf("CREATE TABLESPACE on metadata-loaded tablespace got %d, want 1813", got)
	}
	execAlterClauses(t, metaCat, "ALTER TABLESPACE loaded_ts ENCRYPTION='N';")
	execAlterClauses(t, metaCat, "ALTER TABLESPACE loaded_ts RENAME TO loaded_ts2;")
	if got, _ := simulationCatalogRefusal(t, metaCat, "CREATE TABLE u_old_ts (id INT PRIMARY KEY) TABLESPACE loaded_ts ENGINE=InnoDB"); got != 3510 {
		t.Fatalf("CREATE TABLE on renamed-away loaded_ts got %d, want 3510", got)
	}
	execAlterClauses(t, metaCat, "ALTER TABLESPACE loaded_ts2 RENAME TO loaded_ts; CREATE TABLE u_restored_ts (id INT PRIMARY KEY) TABLESPACE loaded_ts ENGINE=InnoDB;")
	if got, _ := simulationCatalogRefusal(t, metaCat, "CREATE TABLE u_old_ts2 (id INT PRIMARY KEY) TABLESPACE loaded_ts2 ENGINE=InnoDB"); got != 3510 {
		t.Fatalf("CREATE TABLE on renamed-away loaded_ts2 got %d, want 3510", got)
	}

	// Bidirectional invariant: every modeled sysvar in mysqlSysvarSpecs must be a known MySQL 8.0 sysvar.
	for name := range mysqlSysvarSpecs {
		if !validate.IsKnownSystemVariable(name) {
			t.Fatalf("mysqlSysvarSpecs entry %q missing from validate.IsKnownSystemVariable", name)
		}
	}
	sciCat := New()
	execAlterClauses(t, sciCat, "SET @saved_cr = @@CHARACTER_SET_RESULTS; SET sql_mode = IF(1.5e-2 * 2 = 0.03, 'ANSI', '');")
	if !strings.Contains(sciCat.session.SQLMode, "ANSI") {
		t.Fatalf("scientific arithmetic IF got sql_mode=%q, want ANSI", sciCat.session.SQLMode)
	}
	execAlterClauses(t, sciCat, "SET sql_mode = IF('0e5', 'ANSI', '');")
	if sciCat.session.SQLMode != "" {
		t.Fatalf("'0e5' truthiness got sql_mode=%q, want empty", sciCat.session.SQLMode)
	}
	execAlterClauses(t, sciCat, "SET sql_mode = IF(1e-400 = 0, 'ANSI', '');")
	if !strings.Contains(sciCat.session.SQLMode, "ANSI") {
		t.Fatalf("1e-400 underflow comparison got sql_mode=%q, want ANSI", sciCat.session.SQLMode)
	}
	if got, _ := simulationCatalogRefusal(t, sciCat, "SET sql_mode = (1e2 * 2) + 1"); got != 1232 {
		t.Fatalf("compound float arithmetic SQLMode got %d, want 1232", got)
	}

	// ALTER TABLESPACE mysql ENCRYPTION takes 'Y' and 'N' in any case and
	// refuses any other value with 3184.
	mysqlTSCat := New()
	execAlterClauses(t, mysqlTSCat, "ALTER TABLESPACE mysql ENCRYPTION='Y';")
	if !mysqlTSCat.tablespaceEncrypted["mysql"] {
		t.Fatal("ALTER TABLESPACE mysql ENCRYPTION='Y' left the tablespace unencrypted")
	}
	execAlterClauses(t, mysqlTSCat, "ALTER TABLESPACE mysql ENCRYPTION='n';")
	if mysqlTSCat.tablespaceEncrypted["mysql"] {
		t.Fatal("ALTER TABLESPACE mysql ENCRYPTION='n' left the tablespace encrypted")
	}
	for _, stmt := range []string{"ALTER TABLESPACE mysql ENCRYPTION='Z'", "ALTER TABLESPACE mysql ENCRYPTION=''"} {
		if got, _ := simulationCatalogRefusal(t, mysqlTSCat, stmt); got != 3184 {
			t.Fatalf("%s got %d, want 3184", stmt, got)
		}
	}

	// LoadSDL holds a temporary table to the AUTO_INCREMENT key rule, as
	// CREATE TEMPORARY TABLE does on mysql:8.0.46 (1075).
	if _, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TEMPORARY TABLE t (a INT AUTO_INCREMENT);"); err == nil {
		t.Fatal("LoadSDL took a temporary table with an unkeyed AUTO_INCREMENT column, want 1075")
	} else {
		var ce *Error
		if !errors.As(err, &ce) || ce.Code != 1075 {
			t.Fatalf("LoadSDL temporary unkeyed AUTO_INCREMENT: %v, want 1075", err)
		}
	}
	if _, err := LoadSDL("CREATE DATABASE d; USE d; CREATE TEMPORARY TABLE t (a INT AUTO_INCREMENT); CREATE INDEX i ON t (a);"); err != nil {
		t.Fatalf("LoadSDL temporary AUTO_INCREMENT keyed by a later index: %v", err)
	}

	selfFKCat := New()
	execAlterClauses(t, selfFKCat, "CREATE DATABASE d; USE d; CREATE TABLE s_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_s FOREIGN KEY (pid) REFERENCES s_fk(id)); ALTER TABLE s_fk DISCARD TABLESPACE;")
}
