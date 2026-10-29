package review

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk/walktest"
)

// oceanBaseSchema is walkSchema as Bytebase syncs it from OceanBase, with
// an external table and a materialized view, whose storage table the sync
// lists as a table.
func oceanBaseSchema() *metadata.DatabaseSchemaMetadata {
	db := walkSchema()
	s := db.Schemas[0]
	s.Tables = append(s.Tables, &metadata.TableMetadata{Name: "__mv_container_500001", Columns: []*metadata.ColumnMetadata{{Name: "a"}}})
	s.ExternalTables = []*metadata.ExternalTableMetadata{{Name: "ext", Columns: []*metadata.ColumnMetadata{{Name: "a"}}}}
	s.MaterializedViews = []*metadata.MaterializedViewMetadata{{Name: "mv", SkipDump: true}}
	return db
}

func oneFindingOB(t *testing.T, sql string, statement int, message string, target review.Target) {
	t.Helper()
	oneFindingOf(t, ReviewOceanBase, sql, statement, message, target)
}

func noFindingOB(t *testing.T, sql string, target review.Target) {
	t.Helper()
	noFindingOf(t, ReviewOceanBase, sql, target)
}

// TestWalkThroughOceanBase pins each delta of the OceanBase spec
// (docs/specs/2026-10-27-walk-through-oceanbase.md), against what MySQL's
// walk says of the same change where they differ.
func TestWalkThroughOceanBase(t *testing.T) {
	shop := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1}
	mysql := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1}
	oneFindingOB(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table t already exists", shop)

	t.Run("external tables and materialized views", func(t *testing.T) {
		// They share the tables' namespace, of no kind the walk claims.
		oneFindingOB(t, "CREATE TABLE ext (a int);", 0, "Table ext already exists", shop)
		oneFindingOB(t, "CREATE VIEW mv AS SELECT 1 AS a;", 0, "Table mv already exists", shop)
		oneFindingOB(t, "RENAME TABLE t TO mv;", 0, "Table mv already exists", shop)
		noFindingOB(t, "DROP VIEW ext;\nDROP TABLE mv;\nTRUNCATE TABLE ext;\nCREATE INDEX k ON mv (a);\nCREATE OR REPLACE VIEW ext AS SELECT 1 AS a;\nCREATE TABLE x LIKE mv;\nINSERT INTO ext VALUES (1);", shop)
		// What they hold is unknown, and ALTER TABLE of one ends the walk.
		noFindingOB(t, "ALTER TABLE ext DROP COLUMN zz;\nALTER TABLE t DROP COLUMN zz;", shop)
		noFindingOB(t, "RENAME TABLE mv TO mv2;\nALTER TABLE mv2 ADD COLUMN a int;", shop)
		oneFindingOB(t, "RENAME TABLE ext TO e2;\nDROP TABLE ext;", 1, "Table ext does not exist", shop)
	})

	t.Run("materialized views the sync could not read", func(t *testing.T) {
		unread := oceanBaseSchema()
		unread.Schemas[0].MaterializedViews = nil
		target := review.Target{Schema: unread, LowerCaseTableNames: 1}
		noFindingOB(t, "DROP TABLE nosuch;\nALTER TABLE nosuch2 ADD COLUMN a int;\nRENAME TABLE nosuch3 TO x;", target)
		oneFindingOB(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column zz does not exist in table t", target)
		oneFindingOB(t, "DROP TABLE nosuch;", 0, "Table nosuch does not exist", shop)
	})

	t.Run("ALTER TABLE of a view", func(t *testing.T) {
		oneFinding(t, "ALTER TABLE v COMMENT 'x';", 0, "v is not a base table", mysql)
		noFindingOB(t, "ALTER TABLE v COMMENT 'x';\nCREATE TABLE t (a int);", shop)
		oneFindingOB(t, "ALTER TABLE nosuch COMMENT 'x';", 0, "Table nosuch does not exist", shop)
	})

	t.Run("the primary key", func(t *testing.T) {
		// No invisible primary key: a table listed without PRIMARY has none.
		noFinding(t, "ALTER TABLE nokey DROP PRIMARY KEY;", mysql)
		oneFindingOB(t, "ALTER TABLE nokey DROP PRIMARY KEY;", 0, "Table nokey has no primary key", shop)
		oneFindingOB(t, "CREATE TABLE x (a int);\nALTER TABLE x DROP PRIMARY KEY;", 1, "Table x has no primary key", shop)
		oneFindingOB(t, "ALTER TABLE nokey ADD COLUMN x int AFTER my_row_id;", 0, "Column my_row_id does not exist in table nokey", shop)
		// A listed PRIMARY may be a secondary index of the name.
		oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "Table t already has a primary key", mysql)
		noFindingOB(t, "ALTER TABLE t ADD PRIMARY KEY (a);", shop)
		// The hidden column of a keyless table is never listed.
		noFindingOB(t, "ALTER TABLE nokey ADD COLUMN y int AFTER __pk_increment;\nALTER TABLE nokey ADD COLUMN __pk_increment int;\nALTER TABLE t ALTER COLUMN __pk_increment SET DEFAULT 1;", shop)
		noFindingOB(t, "CREATE TABLE x (id int PRIMARY KEY);\nALTER TABLE x ADD COLUMN y int AFTER __pk_increment;", shop)
		noFindingOB(t, "CREATE TABLE x (id int PRIMARY KEY, y int);\nALTER TABLE x ADD COLUMN __pk_increment int;\nALTER TABLE x DROP COLUMN __pk_increment;\nALTER TABLE x DROP PRIMARY KEY;\nALTER TABLE x ADD COLUMN z int AFTER __pk_increment;", shop)
	})

	t.Run("indexes named PRIMARY", func(t *testing.T) {
		oneFinding(t, "CREATE INDEX `PRIMARY` ON t (a);", 0, "Index name PRIMARY is reserved for the primary key", mysql)
		noFindingOB(t, "CREATE INDEX `PRIMARY` ON t (a);\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER TABLE t ADD INDEX `PRIMARY` (a);\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "DROP INDEX `PRIMARY` ON nokey;\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER TABLE nokey DROP INDEX `PRIMARY`;\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER TABLE t RENAME INDEX `PRIMARY` TO k;\nCREATE TABLE t (a int);", shop)
		oneFindingOB(t, "CREATE TABLE x (a int, KEY `PRIMARY` (a));", 0, "Index name PRIMARY is reserved for the primary key", shop)
	})

	t.Run("foreign keys", func(t *testing.T) {
		// Every table keeps its foreign keys, whatever its engine.
		noFinding(t, "ALTER TABLE t ENGINE=MyISAM;\nALTER TABLE t DROP FOREIGN KEY nosuch;", mysql)
		oneFindingOB(t, "ALTER TABLE t ENGINE=MyISAM;\nALTER TABLE t DROP FOREIGN KEY nosuch;", 1, "Foreign key nosuch does not exist in table t", shop)
		oneFindingOB(t, "CREATE TABLE my (id int, a int) ENGINE=MyISAM;\nALTER TABLE my DROP FOREIGN KEY nosuch;", 1, "Foreign key nosuch does not exist in table my", shop)
		// A foreign key dropped twice in one statement is dropped.
		oneFinding(t, "ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;", 0, "Foreign key fk1 is dropped twice from table t", mysql)
		oneFindingOB(t, "ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;", 1, "Foreign key fk1 does not exist in table t", shop)
		// Generated names are <table>_OBFK_<16 digits>, and a rename keeps
		// them.
		const g = "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\n"
		noFindingOB(t, g+"ALTER TABLE g DROP FOREIGN KEY g_obfk_1700000000000000;", shop)
		oneFindingOB(t, g+"ALTER TABLE g DROP FOREIGN KEY g_ibfk_1;", 1, "Foreign key g_ibfk_1 does not exist in table g", shop)
		noFindingOB(t, g+"RENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY g_obfk_1700000000000000;", shop)
		oneFindingOB(t, g+"RENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY h_obfk_1700000000000000;", 2, "Foreign key h_obfk_1700000000000000 does not exist in table h", shop)
		// An unqualified parent is in the current database.
		on := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}
		const child = "CREATE DATABASE q;\nCREATE TABLE q.parent (id int PRIMARY KEY);\nCREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES "
		oneFinding(t, child+"r (id));", 2, "Foreign key target table r does not exist", review.Target{Schema: walkSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn})
		noFindingOB(t, child+"r (id));", on)
		oneFindingOB(t, child+"parent (id));", 2, "Foreign key target table parent does not exist", on)
		noFindingOB(t, "DROP DATABASE shop;\nCREATE TABLE q.c (id int, FOREIGN KEY (id) REFERENCES nosuch (id));", on)
	})

	t.Run("name lengths", func(t *testing.T) {
		long := func(n int) string { return strings.Repeat("a", n) }
		oneFinding(t, "CREATE DATABASE "+long(100)+";", 0, "Identifier name "+long(100)+" is too long", mysql)
		noFindingOB(t, "CREATE DATABASE "+long(100)+";\nALTER TABLE t ADD COLUMN "+long(100)+" int;\nCREATE TABLE x (a int, CONSTRAINT "+long(100)+" FOREIGN KEY (a) REFERENCES r (id));", shop)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN "+long(129)+" int;", 0, "Identifier name "+long(129)+" is too long", shop)
		oneFindingOB(t, "CREATE TABLE "+long(65)+" (a int);", 0, "Identifier name "+long(65)+" is too long", shop)
		oneFindingOB(t, "CREATE INDEX "+long(65)+" ON t (a);", 0, "Identifier name "+long(65)+" is too long", shop)
		// The name after FOREIGN KEY names nothing: the foreign key takes
		// its symbol or a generated name, and adds no index.
		for _, stmt := range []string{"CREATE TABLE w (a int, FOREIGN KEY " + long(129) + " (a) REFERENCES r (id));", "ALTER TABLE nokey ADD FOREIGN KEY " + long(129) + " (a) REFERENCES r (id);"} {
			oneFinding(t, stmt, 0, "Identifier name "+long(129)+" is too long", mysql)
			noFindingOB(t, stmt, shop)
		}
		oneFindingOB(t, "CREATE TABLE w (a int, FOREIGN KEY fw (a) REFERENCES r (id));\nALTER TABLE w DROP FOREIGN KEY fw;", 1, "Foreign key fw does not exist in table w", shop)
	})

	t.Run("SET foreign_key_checks", func(t *testing.T) {
		on := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}
		onMySQL := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}
		const fk = "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));"
		const missing = "Foreign key target table nosuch does not exist"
		// A scope keyword reaches only the first assignment.
		const scoped = "SET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\n" + fk
		oneFinding(t, scoped, 1, missing, onMySQL)
		noFindingOB(t, scoped, on)
		oneFindingOB(t, "SET GLOBAL foreign_key_checks = 0;\n"+fk, 1, missing, on)
		oneFindingOB(t, "SET GLOBAL foreign_key_checks = 0, wait_timeout = 28800;\n"+fk, 1, missing, on)
		// An assignment after one to DEFAULT may not take its value.
		const afterDefault = "SET foreign_key_checks = 0;\nSET sql_mode = DEFAULT, foreign_key_checks = 1;\n" + fk
		oneFinding(t, afterDefault, 2, missing, onMySQL)
		noFindingOB(t, afterDefault, on)
		oneFindingOB(t, "SET foreign_key_checks = 0;\nSET foreign_key_checks = 1, sql_mode = DEFAULT;\n"+fk, 2, missing, on)
		oneFindingOB(t, "SET foreign_key_checks = 0;\nSET sql_mode = 'STRICT_TRANS_TABLES', foreign_key_checks = 1;\n"+fk, 2, missing, on)
	})

	t.Run("DROP DATABASE drops its temporary tables", func(t *testing.T) {
		const change = "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;"
		noFinding(t, change, mysql)
		oneFindingOB(t, change, 4, "Temporary table q.tmp does not exist", shop)
	})

	t.Run("views", func(t *testing.T) {
		oneFindingOB(t, "DROP VIEW IF EXISTS t;", 0, "t is not a view", shop)
		noFinding(t, "DROP VIEW IF EXISTS t;", mysql)
		noFindingOB(t, "CREATE TEMPORARY TABLE t (a int);\nDROP VIEW IF EXISTS t;", shop)
		const move = "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;"
		oneFinding(t, move, 1, "View v cannot move to another database", mysql)
		oneFindingOB(t, move, 2, "View v does not exist", shop)
	})

	t.Run("routines and events", func(t *testing.T) {
		// No loadable function: an unqualified DROP FUNCTION is the stored
		// one's, but with no database selected.
		noFinding(t, "DROP FUNCTION nosuch;", mysql)
		oneFindingOB(t, "DROP FUNCTION nosuch;", 0, "FUNCTION nosuch does not exist", shop)
		noFindingOB(t, "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS f;\nDROP FUNCTION f;", shop)
		oneFindingOB(t, "DROP DATABASE shop;\nDROP PROCEDURE IF EXISTS p;", 1, "No database selected", shop)
		// A listed event may be a dropped database's.
		oneFinding(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 0, "Event e already exists", mysql)
		noFindingOB(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", shop)
		noFindingOB(t, "CREATE EVENT e9 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nALTER EVENT e9 RENAME TO e;", shop)
		oneFindingOB(t, "DROP EVENT nosuch;", 0, "Event nosuch does not exist", shop)
		oneFindingOB(t, "CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE EVENT E2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 1, "Event E2 already exists", shop)
		// OceanBase splits an event name at its period.
		noFindingOB(t, "CREATE EVENT `shop.e3` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER EVENT e RENAME TO `q.e`;\nCREATE TABLE t (a int);", shop)
	})
}

// oceanBaseForms are the statement forms OceanBase's deltas add to
// walkForms, on oceanBaseSchema.
var oceanBaseForms = []walktest.Form{
	{SQL: "ALTER TABLE v COMMENT 'x';", Last: true},
	{SQL: "ALTER TABLE ext ADD COLUMN b int;", Last: true},
	{SQL: "ALTER TABLE mv COMMENT 'x';", Last: true},
	{SQL: "DROP TABLE ext;"},
	{SQL: "RENAME TABLE mv TO mv2;"},
	{SQL: "CREATE INDEX k ON mv (a);"},
	{SQL: "DROP VIEW IF EXISTS t;"},
	{SQL: "DROP FUNCTION f;"},
	{SQL: "ALTER TABLE t ENGINE = MyISAM;"},
}

func TestWalkThroughOceanBaseLowersEveryForm(t *testing.T) {
	walktest.Coverage(t, append(append([]walktest.Form{}, walkForms...), oceanBaseForms...), formLowering(oceanBaseSchema(), dialect{oceanBase: true}))
}
