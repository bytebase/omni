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
	oneFindingOB(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table `t` already exists.", shop)
	// method_opt is a statement that runs nothing; the walk follows it.
	noFindingOB(t, "FOR COLUMNS;", shop)
	oneFindingOB(t, "FOR COLUMNS;\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;", 2, "Column `c` already exists in `t`.", shop)
	// An external table is a relation of no kind, as a synced one is: its
	// CREATE claims the name and its ALTER is followed like O7's.
	noFindingOB(t, "CREATE EXTERNAL TABLE ext2 (a INT) LOCATION '/tmp/x' FORMAT(TYPE = 'CSV');", shop)
	oneFindingOB(t, "CREATE EXTERNAL TABLE ext (a INT) LOCATION '/tmp/x' FORMAT(TYPE = 'CSV');", 0, "The name `ext` is already taken.", shop)
	noFindingOB(t, "CREATE EXTERNAL TABLE ext2 (a INT) LOCATION '/tmp/x' FORMAT(TYPE = 'CSV');\nALTER EXTERNAL TABLE ext2 ADD COLUMN b INT;", shop)
	// An ALTER EXTERNAL TABLE changes nothing the walk holds, and it does not
	// end the walk: the duplicate below is still found.
	oneFindingOB(t, "ALTER EXTERNAL TABLE ext ADD PARTITION() LOCATION '/tmp/x';\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;", 2, "Column `c` already exists in `t`.", shop)
	// DROP TABLE takes an external table, the snapshot's or one the change made.
	noFindingOB(t, "DROP TABLE ext;", shop)
	noFindingOB(t, "CREATE EXTERNAL TABLE ext2 (a INT) LOCATION '/tmp/x' FORMAT(TYPE = 'CSV');\nDROP TABLE ext2;", shop)
	// A rename moves the ledger's name with it, so the new name is an external
	// table too, and the name goes when it is dropped.
	noFindingOB(t, "CREATE EXTERNAL TABLE ext2 (a INT) LOCATION '/tmp/x' FORMAT(TYPE = 'CSV');\nALTER EXTERNAL TABLE ext2 RENAME TO ext3;\nDROP TABLE ext3;", shop)
	oneFindingOB(t, "ALTER EXTERNAL TABLE ext RENAME TO t;", 0, "The name `t` is already taken.", shop)
	// ALTER EXTERNAL TABLE names an external table.
	oneFindingOB(t, "ALTER EXTERNAL TABLE nosuch ADD COLUMN b INT;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFindingOB(t, "ALTER EXTERNAL TABLE t ADD COLUMN b INT;", 0, "`t` is a table, not an external table.", shop)
	// A definition the server refuses is refused here too.
	oneFindingOB(t, "CREATE EXTERNAL TABLE ext2 (a INT, a INT) LOCATION '/tmp/x' FORMAT(TYPE = 'CSV');", 0, "This statement declares column `a` twice.", shop)

	t.Run("external tables and materialized views", func(t *testing.T) {
		// They share the tables' namespace, of no kind the walk claims.
		oneFindingOB(t, "CREATE TABLE ext (a int);", 0, "The name `ext` is already taken.", shop)
		oneFindingOB(t, "CREATE VIEW mv AS SELECT 1 AS a;", 0, "The name `mv` is already taken.", shop)
		oneFindingOB(t, "RENAME TABLE t TO mv;", 0, "The name `mv` is already taken.", shop)
		noFindingOB(t, "DROP VIEW ext;\nDROP TABLE mv;\nTRUNCATE TABLE ext;\nCREATE INDEX k ON mv (a);\nCREATE OR REPLACE VIEW ext AS SELECT 1 AS a;\nCREATE TABLE x LIKE mv;\nINSERT INTO ext VALUES (1);", shop)
		// What they hold is unknown, and ALTER TABLE of one ends the walk.
		noFindingOB(t, "ALTER TABLE ext DROP COLUMN zz;\nALTER TABLE t DROP COLUMN zz;", shop)
		noFindingOB(t, "RENAME TABLE mv TO mv2;\nALTER TABLE mv2 ADD COLUMN a int;", shop)
		oneFindingOB(t, "RENAME TABLE ext TO e2;\nDROP TABLE ext;", 1, "Table `ext` doesn't exist.", shop)
	})

	t.Run("materialized views the sync could not read", func(t *testing.T) {
		unread := oceanBaseSchema()
		unread.Schemas[0].MaterializedViews = nil
		target := review.Target{Schema: unread, LowerCaseTableNames: 1}
		noFindingOB(t, "DROP TABLE nosuch;\nALTER TABLE nosuch2 ADD COLUMN a int;\nRENAME TABLE nosuch3 TO x;", target)
		oneFindingOB(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column `zz` doesn't exist in `t`.", target)
		oneFindingOB(t, "DROP TABLE nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	})

	t.Run("ALTER TABLE of a view", func(t *testing.T) {
		oneFinding(t, "ALTER TABLE v COMMENT 'x';", 0, "`v` is a view, not a table.", mysql)
		noFindingOB(t, "ALTER TABLE v COMMENT 'x';\nCREATE TABLE t (a int);", shop)
		oneFindingOB(t, "ALTER TABLE nosuch COMMENT 'x';", 0, "Table `nosuch` doesn't exist.", shop)
	})

	t.Run("the primary key", func(t *testing.T) {
		// No invisible primary key: a table listed without PRIMARY has none.
		noFinding(t, "ALTER TABLE nokey DROP PRIMARY KEY;", mysql)
		oneFindingOB(t, "ALTER TABLE nokey DROP PRIMARY KEY;", 0, "`nokey` has no primary key.", shop)
		oneFindingOB(t, "CREATE TABLE x (a int);\nALTER TABLE x DROP PRIMARY KEY;", 1, "`x` has no primary key.", shop)
		oneFindingOB(t, "ALTER TABLE nokey ADD COLUMN x int AFTER my_row_id;", 0, "Column `my_row_id` doesn't exist in `nokey`.", shop)
		// A listed PRIMARY may be a secondary index of the name.
		oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "`t` already has a primary key.", mysql)
		noFindingOB(t, "ALTER TABLE t ADD PRIMARY KEY (a);", shop)
		// The hidden column of a keyless table is never listed.
		noFindingOB(t, "ALTER TABLE nokey ADD COLUMN y int AFTER __pk_increment;\nALTER TABLE nokey ADD COLUMN __pk_increment int;\nALTER TABLE t ALTER COLUMN __pk_increment SET DEFAULT 1;", shop)
		noFindingOB(t, "CREATE TABLE x (id int PRIMARY KEY);\nALTER TABLE x ADD COLUMN y int AFTER __pk_increment;", shop)
		noFindingOB(t, "CREATE TABLE x (id int PRIMARY KEY, y int);\nALTER TABLE x ADD COLUMN __pk_increment int;\nALTER TABLE x DROP COLUMN __pk_increment;\nALTER TABLE x DROP PRIMARY KEY;\nALTER TABLE x ADD COLUMN z int AFTER __pk_increment;", shop)
	})

	t.Run("indexes named PRIMARY", func(t *testing.T) {
		oneFinding(t, "CREATE INDEX `PRIMARY` ON t (a);", 0, "`PRIMARY` is reserved for the primary key's index.", mysql)
		noFindingOB(t, "CREATE INDEX `PRIMARY` ON t (a);\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER TABLE t ADD INDEX `PRIMARY` (a);\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "DROP INDEX `PRIMARY` ON nokey;\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER TABLE nokey DROP INDEX `PRIMARY`;\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER TABLE t RENAME INDEX `PRIMARY` TO k;\nCREATE TABLE t (a int);", shop)
		oneFindingOB(t, "CREATE TABLE x (a int, KEY `PRIMARY` (a));", 0, "`PRIMARY` is reserved for the primary key's index.", shop)
	})

	t.Run("foreign keys", func(t *testing.T) {
		// Every table keeps its foreign keys, whatever its engine.
		noFinding(t, "ALTER TABLE t ENGINE=MyISAM;\nALTER TABLE t DROP FOREIGN KEY nosuch;", mysql)
		oneFindingOB(t, "ALTER TABLE t ENGINE=MyISAM;\nALTER TABLE t DROP FOREIGN KEY nosuch;", 1, "Foreign key `nosuch` doesn't exist in `t`.", shop)
		oneFindingOB(t, "CREATE TABLE my (id int, a int) ENGINE=MyISAM;\nALTER TABLE my DROP FOREIGN KEY nosuch;", 1, "Foreign key `nosuch` doesn't exist in `my`.", shop)
		// A foreign key dropped twice in one statement is dropped.
		oneFinding(t, "ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;", 0, "This statement drops foreign key `fk1` twice.", mysql)
		oneFindingOB(t, "ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;", 1, "Foreign key `fk1` doesn't exist in `t`.", shop)
		// Generated names are <table>_OBFK_<16 digits>, and a rename keeps
		// them.
		const g = "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\n"
		noFindingOB(t, g+"ALTER TABLE g DROP FOREIGN KEY g_obfk_1700000000000000;", shop)
		oneFindingOB(t, g+"ALTER TABLE g DROP FOREIGN KEY g_ibfk_1;", 1, "Foreign key `g_ibfk_1` doesn't exist in `g`.", shop)
		noFindingOB(t, g+"RENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY g_obfk_1700000000000000;", shop)
		oneFindingOB(t, g+"RENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY h_obfk_1700000000000000;", 2, "Foreign key `h_obfk_1700000000000000` doesn't exist in `h`.", shop)
		// An unqualified parent is in the current database.
		on := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}
		const child = "CREATE DATABASE q;\nCREATE TABLE q.parent (id int PRIMARY KEY);\nCREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES "
		oneFinding(t, child+"r (id));", 2, "Referenced table `r` doesn't exist.", review.Target{Schema: walkSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn})
		noFindingOB(t, child+"r (id));", on)
		oneFindingOB(t, child+"parent (id));", 2, "Referenced table `parent` doesn't exist.", on)
		noFindingOB(t, "DROP DATABASE shop;\nCREATE TABLE q.c (id int, FOREIGN KEY (id) REFERENCES nosuch (id));", on)
	})

	t.Run("name lengths", func(t *testing.T) {
		long := func(n int) string { return strings.Repeat("a", n) }
		oneFinding(t, "CREATE DATABASE "+long(100)+";", 0, "The name `"+long(100)+"` is longer than MySQL's 64-character limit.", mysql)
		noFindingOB(t, "CREATE DATABASE "+long(100)+";\nALTER TABLE t ADD COLUMN "+long(100)+" int;\nCREATE TABLE x (a int, CONSTRAINT "+long(100)+" FOREIGN KEY (a) REFERENCES r (id));", shop)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN "+long(129)+" int;", 0, "The name "+review.Code(long(129))+" is longer than OceanBase's 128-character limit.", shop)
		oneFindingOB(t, "CREATE TABLE "+long(65)+" (a int);", 0, "The name `"+long(65)+"` is longer than OceanBase's 64-character limit.", shop)
		oneFindingOB(t, "CREATE INDEX "+long(65)+" ON t (a);", 0, "The name `"+long(65)+"` is longer than OceanBase's 64-character limit.", shop)
		// The name after FOREIGN KEY names nothing: the foreign key takes
		// its symbol or a generated name, and adds no index.
		for _, stmt := range []string{"CREATE TABLE w (a int, FOREIGN KEY " + long(129) + " (a) REFERENCES r (id));", "ALTER TABLE nokey ADD FOREIGN KEY " + long(129) + " (a) REFERENCES r (id);"} {
			oneFinding(t, stmt, 0, "The name "+review.Code(long(129))+" is longer than MySQL's 64-character limit.", mysql)
			noFindingOB(t, stmt, shop)
		}
		oneFindingOB(t, "CREATE TABLE w (a int, FOREIGN KEY fw (a) REFERENCES r (id));\nALTER TABLE w DROP FOREIGN KEY fw;", 1, "Foreign key `fw` doesn't exist in `w`.", shop)
	})

	t.Run("SET foreign_key_checks", func(t *testing.T) {
		on := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}
		onMySQL := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}
		const fk = "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));"
		const missing = "Referenced table `nosuch` doesn't exist."
		// A scope keyword reaches only the first assignment.
		const scoped = "SET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\n" + fk
		oneFinding(t, scoped, 1, missing, onMySQL)
		noFindingOB(t, scoped, on)
		// A later assignment writes its own scope, and it reaches no other.
		oneFindingOB(t, "SET @a = 1, GLOBAL foreign_key_checks = 0;\n"+fk, 1, missing, on)
		noFindingOB(t, "SET GLOBAL sql_mode = '', SESSION foreign_key_checks = 0;\n"+fk, on)
		noFindingOB(t, "SET SESSION sql_mode = '', GLOBAL sql_mode = '', foreign_key_checks = 0;\n"+fk, on)
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
		oneFindingOB(t, change, 4, "Temporary table `q.tmp` doesn't exist.", shop)
	})

	t.Run("views", func(t *testing.T) {
		oneFindingOB(t, "DROP VIEW IF EXISTS t;", 0, "`t` is a table, not a view.", shop)
		noFinding(t, "DROP VIEW IF EXISTS t;", mysql)
		noFindingOB(t, "CREATE TEMPORARY TABLE t (a int);\nDROP VIEW IF EXISTS t;", shop)
		const move = "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;"
		oneFinding(t, move, 1, "MySQL doesn't accept `RENAME TABLE` of view `v` to another database.", mysql)
		oneFindingOB(t, move, 2, "View `v` doesn't exist.", shop)
	})

	t.Run("routines and events", func(t *testing.T) {
		// No loadable function: an unqualified DROP FUNCTION is the stored
		// one's, but with no database selected.
		noFinding(t, "DROP FUNCTION nosuch;", mysql)
		oneFindingOB(t, "DROP FUNCTION nosuch;", 0, "Function `nosuch` doesn't exist.", shop)
		noFindingOB(t, "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS f;\nDROP FUNCTION f;", shop)
		oneFindingOB(t, "DROP DATABASE shop;\nDROP PROCEDURE IF EXISTS p;", 1, "No database is selected for `p`.", shop)
		// A listed event may be a dropped database's.
		oneFinding(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 0, "Event `e` already exists.", mysql)
		noFindingOB(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", shop)
		noFindingOB(t, "CREATE EVENT e9 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nALTER EVENT e9 RENAME TO e;", shop)
		oneFindingOB(t, "DROP EVENT nosuch;", 0, "Event `nosuch` doesn't exist.", shop)
		oneFindingOB(t, "CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE EVENT E2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 1, "Event `E2` already exists.", shop)
		// OceanBase splits an event name at its period.
		noFindingOB(t, "CREATE EVENT `shop.e3` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE TABLE t (a int);", shop)
		noFindingOB(t, "ALTER EVENT e RENAME TO `q.e`;\nCREATE TABLE t (a int);", shop)
	})

	t.Run("FLASHBACK stops the walk and SET CATALOG marks session lost", func(t *testing.T) {
		noFindingOB(t, "DROP TABLE t;\nFLASHBACK TABLE t TO BEFORE DROP;\nALTER TABLE t ADD COLUMN c int;", shop)
		noFindingOB(t, "SET CATALOG ext_cat;\nALTER TABLE nosuch ADD COLUMN c int;", shop)
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

func TestWalkThroughRejectedAlterActions(t *testing.T) {
	shop := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1}
	// A column is no BEFORE or AFTER target of its own move (1054).
	oneFindingOB(t, "ALTER TABLE t MODIFY a INT BEFORE a;", 0, "Column `a` doesn't exist in `t`.", shop)
	oneFindingOB(t, "ALTER TABLE t MODIFY a INT AFTER a;", 0, "Column `a` doesn't exist in `t`.", shop)
	oneFindingOB(t, "ALTER TABLE t CHANGE a b INT BEFORE a;", 0, "Column `a` doesn't exist in `t`.", shop)
	oneFindingOB(t, "ALTER TABLE t CHANGE a b INT AFTER a;", 0, "Column `a` doesn't exist in `t`.", shop)
	oneFindingOB(t, "ALTER TABLE t CHANGE a b INT BEFORE b;", 0, "Column `b` doesn't exist in `t`.", shop)
	// REFRESH, storage cache policy, LOB storage, column ID/CHUNK/ORIG_DEFAULT,
	// quoted data types, and view TABLE_ID are refused rather than mutating the ledger.
	oneFindingOB(t, "ALTER TABLE t REFRESH;", 0, "OceanBase refreshes no metadata of table `t`.", shop)
	oneFindingOB(t, "ALTER TABLE t ALTER INDEX i STORAGE_CACHE_POLICY(GLOBAL = 's');", 0, "OceanBase has no storage cache policy for an index.", shop)
	oneFindingOB(t, "ALTER TABLE t ADD COLUMN c INT ID 1;", 0, "OceanBase refuses explicit column ID on ALTER TABLE.", shop)
	oneFindingOB(t, "ALTER TABLE t ADD (c JSON) JSON(c) STORE AS (CHUNK 1);", 0, "OceanBase refuses LOB storage parameters on ALTER TABLE ADD.", shop)
	oneFindingOB(t, "ALTER TABLE t ADD COLUMN c 'int';", 0, "OceanBase refuses a quoted string literal as a column data type.", shop)
	oneFindingOB(t, "CREATE VIEW v2 TABLE_ID = 1 AS SELECT 1;", 0, "OceanBase refuses TABLE_ID on a view.", shop)
}

func TestWalkThroughOceanBaseRecallFamilies(t *testing.T) {
	// Rich schema with typed columns, index column lists, and FK column lists
	// matching synced OceanBase metadata.
	richDB := func() *metadata.DatabaseSchemaMetadata {
		return &metadata.DatabaseSchemaMetadata{
			Name: "shop",
			Schemas: []*metadata.SchemaMetadata{{
				Tables: []*metadata.TableMetadata{
					{
						Name: "parent",
						Columns: []*metadata.ColumnMetadata{
							{Name: "id", Type: "int(11)"},
						},
						Indexes: []*metadata.IndexMetadata{
							{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}},
						},
					},
					{
						Name: "t",
						Columns: []*metadata.ColumnMetadata{
							{Name: "id", Type: "int(11)"},
							{Name: "a", Type: "int(11)"},
							{Name: "name", Type: "varchar(64)"},
							{Name: "ts", Type: "datetime"},
							{Name: "geo", Type: "point"},
							{Name: "vec", Type: "vector(3)"},
							{Name: "parent_id", Type: "int(11)"},
						},
						Indexes: []*metadata.IndexMetadata{
							{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}},
							{Name: "idx_a", Expressions: []string{"a"}},
						},
						ForeignKeys: []*metadata.ForeignKeyMetadata{
							{Name: "fk_parent", Columns: []string{"parent_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}},
						},
						Triggers: []*metadata.TriggerMetadata{
							{Name: "tr_existing"},
						},
					},
				},
				Views: []*metadata.ViewMetadata{{Name: "v"}},
			}},
		}
	}
	target := review.Target{Schema: richDB(), LowerCaseTableNames: 1, ForeignKeyChecks: review.SwitchOn}

	t.Run("ALTER TABLE on view", func(t *testing.T) {
		noFindingOB(t, "ALTER TABLE v COMMENT 'ok';", target)
		oneFindingOB(t, "ALTER TABLE v ADD COLUMN c INT;", 0, "`v` is a view, not a table.", target)
		oneFindingOB(t, "ALTER TABLE v ENGINE = InnoDB;", 0, "`v` is a view, not a table.", target)
		oneFindingOB(t, "ALTER TABLE v;", 0, "`v` is a view, not a table.", target)
	})

	t.Run("table options, charsets, column groups, and TTL", func(t *testing.T) {
		oneFindingOB(t, "ALTER TABLE t FORMAT(ENCODING = 'utf8');", 0, "OceanBase refuses external table option `FORMAT` on a base table.", target)
		oneFindingOB(t, "ALTER TABLE t TABLE_ID 1;", 0, "OceanBase refuses altering `TABLE_ID`.", target)
		oneFindingOB(t, "ALTER TABLE t ORGANIZATION INDEX;", 0, "OceanBase refuses `ORGANIZATION` on ALTER TABLE.", target)
		oneFindingOB(t, "ALTER TABLE t MERGE_ENGINE PARTIAL_UPDATE;", 0, "OceanBase refuses altering `MERGE_ENGINE`.", target)
		oneFindingOB(t, "ALTER TABLE t AUTO_REFRESH IMMEDIATE;", 0, "OceanBase refuses `AUTO_REFRESH` on ALTER TABLE.", target)
		oneFindingOB(t, "ALTER TABLE t STORAGE_CACHE_POLICY(GLOBAL = 's');", 0, "OceanBase has no storage cache policy for a table.", target)
		oneFindingOB(t, "ALTER TABLE t TABLE_MODE 'bad';", 0, "OceanBase refuses `TABLE_MODE` value `bad`.", target)
		noFindingOB(t, "ALTER TABLE t TABLE_MODE 'queuing';", target)
		oneFindingOB(t, "ALTER TABLE t DUPLICATE_SCOPE 'bad';", 0, "OceanBase refuses `DUPLICATE_SCOPE` value `bad`.", target)
		noFindingOB(t, "ALTER TABLE t DUPLICATE_SCOPE 'cluster';", target)
		oneFindingOB(t, "ALTER TABLE t BLOCK_SIZE 1;", 0, "OceanBase requires `BLOCK_SIZE` between 1024 and 1048576.", target)
		noFindingOB(t, "ALTER TABLE t BLOCK_SIZE 16384;", target)
		oneFindingOB(t, "ALTER TABLE t TABLET_SIZE 1;", 0, "OceanBase requires `TABLET_SIZE` to be a multiple of 2 MiB.", target)
		noFindingOB(t, "ALTER TABLE t TABLET_SIZE 2097152;", target)
		oneFindingOB(t, "ALTER TABLE t COMPRESSION 'bad';", 0, "OceanBase refuses `COMPRESSION` value `bad`.", target)
		noFindingOB(t, "ALTER TABLE t COMPRESSION 'zstd_1.3.8';", target)
		oneFindingOB(t, "ALTER TABLE t KV_ATTRIBUTES 'not-json';", 0, "OceanBase requires valid JSON for `KV_ATTRIBUTES`.", target)
		noFindingOB(t, `ALTER TABLE t KV_ATTRIBUTES '{"ttl": 1}';`, target)
		oneFindingOB(t, "ALTER TABLE t CHARSET nosuch;", 0, "OceanBase has no character set `nosuch`.", target)
		noFindingOB(t, "ALTER TABLE t CHARSET utf8mb4;", target)
		oneFindingOB(t, "ALTER TABLE t COLLATE nosuch;", 0, "OceanBase has no collation `nosuch`.", target)
		noFindingOB(t, "ALTER TABLE t COLLATE utf8mb4_general_ci;", target)
		oneFindingOB(t, "ALTER TABLE t DYNAMIC_PARTITION_POLICY(TIME_ZONE = 'bad');", 0, "OceanBase refuses `TIME_ZONE` value `bad` in `DYNAMIC_PARTITION_POLICY`.", target)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS) DELAYED;", 0, "OceanBase refuses `DELAYED` on ADD COLUMN GROUP without EACH COLUMN.", target)
		noFindingOB(t, "ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN) DELAYED;", target)
		oneFindingOB(t, "ALTER TABLE t TTL(nosuch + INTERVAL 1 DAY);", 0, "Column `nosuch` doesn't exist in `t`.", target)
		oneFindingOB(t, "ALTER TABLE t TTL(a + INTERVAL 1 DAY);", 0, "OceanBase requires a DATETIME or TIMESTAMP column for TTL, not `a`.", target)
		noFindingOB(t, "ALTER TABLE t TTL(ts + INTERVAL 1 DAY);", target)
		noFindingOB(t, "ALTER TABLE t ADD COLUMN created_at DATETIME, TTL(created_at + INTERVAL 1 DAY);", target)
	})

	t.Run("index options and key parts", func(t *testing.T) {
		oneFindingOB(t, "ALTER TABLE t ADD KEY(a) WITH(k = 1);", 0, "OceanBase refuses `WITH` index parameters on a non-vector index.", target)
		oneFindingOB(t, "ALTER TABLE t ADD KEY(a) DATA_TABLE_ID 1;", 0, "OceanBase refuses `DATA_TABLE_ID` on an index.", target)
		oneFindingOB(t, "ALTER TABLE t ADD KEY(a) BLOCK_SIZE 1;", 0, "OceanBase requires `BLOCK_SIZE` between 1024 and 1048576.", target)
		oneFindingOB(t, "ALTER TABLE t ADD KEY(a, a);", 0, "This statement names column `a` twice.", target)
		oneFindingOB(t, "ALTER TABLE t ADD KEY(nosuch);", 0, "Column `nosuch` doesn't exist in `t`.", target)
		oneFindingOB(t, "CREATE INDEX k_missing ON t(nosuch);", 0, "Column `nosuch` doesn't exist in `t`.", target)
		oneFindingOB(t, "ALTER TABLE t ADD KEY((a));", 0, "OceanBase refuses a functional index on bare column `a`.", target)
		oneFindingOB(t, "ALTER TABLE t ADD KEY(a(1));", 0, "OceanBase refuses a prefix key on non-string column `a`.", target)
		noFindingOB(t, "ALTER TABLE t ADD KEY(name(10));", target)
		oneFindingOB(t, "ALTER TABLE t ADD FULLTEXT(a);", 0, "OceanBase refuses a FULLTEXT index on non-string column `a`.", target)
		noFindingOB(t, "ALTER TABLE t ADD FULLTEXT(name);", target)
		noFindingOB(t, "ALTER TABLE t MODIFY COLUMN a VARCHAR(32), ADD FULLTEXT(a);", target)
		oneFindingOB(t, "CREATE SPATIAL INDEX k_sp ON t(a);", 0, "OceanBase refuses a SPATIAL index on non-spatial column `a`.", target)
		noFindingOB(t, "CREATE SPATIAL INDEX k_sp ON t(geo);", target)
		oneFindingOB(t, "CREATE VECTOR INDEX k_vec ON t(a);", 0, "OceanBase refuses a VECTOR index on non-vector column `a`.", target)
		noFindingOB(t, "CREATE VECTOR INDEX k_vec ON t(vec) WITH (distance=l2, type=hnsw, lib=vsag);", target)
		oneFindingOB(t, "ALTER TABLE t RENAME INDEX idx_a TO idx_a;", 0, "Index `idx_a` already exists in `t`.", target)
	})

	t.Run("column definitions, CHECK, FK, partitioning, and table statements", func(t *testing.T) {
		oneFindingOB(t, "CREATE TABLE empty_t (PRIMARY KEY (a));", 0, "Table `empty_t` must have at least one column.", target)
		oneFindingOB(t, "CREATE TABLE ctas AS SELECT nosuch;", 0, "Column `nosuch` doesn't exist.", target)
		oneFindingOB(t, "CREATE TABLE ctas AS SELECT nosuch alias_c;", 0, "Column `nosuch` doesn't exist.", target)
		noFindingOB(t, "CREATE TABLE ctas AS SELECT 1 AS c;", target)
		oneFindingOB(t, "CREATE VIEW v_bad AS SELECT nosuch;", 0, "Column `nosuch` doesn't exist.", target)
		noFindingOB(t, "CREATE VIEW v_ok AS SELECT TIMESTAMPDIFF(DAY, '2020-01-01', '2020-01-02') AS d;", target)
		oneFindingOB(t, "ALTER VIEW v AS SELECT nosuch;", 0, "Column `nosuch` doesn't exist.", target)
		oneFindingOB(t, "DROP TABLE t, t;", 0, "This statement drops table `t` twice.", target)
		oneFindingOB(t, "RENAME TABLE t TO t, t TO t;", 0, "Table `t` already exists.", target)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN g INT AS (nosuch);", 0, "Column `nosuch` doesn't exist in `t`.", target)
		noFindingOB(t, "ALTER TABLE t ADD COLUMN (x INT, y INT AS (x + 1));", target)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN x TINYINT DEFAULT 'bad';", 0, "Invalid default value for column `x`.", target)
		noFindingOB(t, "ALTER TABLE t ADD COLUMN x TINYINT DEFAULT '12';", target)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN x INT ON UPDATE CURRENT_TIMESTAMP;", 0, "OceanBase refuses `ON UPDATE CURRENT_TIMESTAMP` on non-temporal column `x`.", target)
		oneFindingOB(t, "ALTER TABLE t ADD COLUMN x INT SRID 4326;", 0, "OceanBase refuses `SRID` on non-spatial column `x`.", target)
		oneFindingOB(t, "ALTER TABLE t ALTER a SET DEFAULT (a);", 0, "OceanBase refuses column references in a default value expression.", target)
		oneFindingOB(t, "ALTER TABLE t ALTER a SET DEFAULT 'bad';", 0, "Invalid default value for column `a`.", target)
		oneFindingOB(t, "ALTER TABLE t ADD CHECK (a);", 0, "OceanBase requires a boolean expression for a CHECK constraint.", target)
		noFindingOB(t, "ALTER TABLE t ADD CHECK (a > 0);", target)
		oneFindingOB(t, "ALTER TABLE t ALTER CONSTRAINT fk_parent NOT ENFORCED;", 0, "OceanBase refuses altering enforcement on foreign key `fk_parent`.", target)
		oneFindingOB(t, "ALTER TABLE t DROP COLUMN parent_id;", 0, "Column `parent_id` is used in a foreign key constraint on `t`.", target)
		oneFindingOB(t, "ALTER TABLE t MODIFY COLUMN parent_id TINYINT;", 0, "Column `parent_id` is used in a foreign key constraint on `t`.", target)
		noFindingOB(t, "ALTER TABLE t MODIFY COLUMN parent_id INT NOT NULL;", target)
		noFindingOB(t, "ALTER TABLE t DROP FOREIGN KEY fk_parent, DROP COLUMN parent_id;", target)
		oneFindingOB(t, "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t (name);", 0, "OceanBase requires an index on referenced columns of `t` for a foreign key.", target)
		noFindingOB(t, "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t (id);", target)
		oneFindingOB(t, "ALTER TABLE t PARTITION BY HASH(a);", 0, "OceanBase requires the primary key of `t` to include partitioning column `a`.", target)
		noFindingOB(t, "ALTER TABLE t PARTITION BY HASH(id);", target)
		oneFindingOB(t, "ALTER TABLE t PARTITION BY KEY()(PARTITION p1 (SUBPARTITION s1));", 0, "OceanBase refuses subpartitions without SUBPARTITION BY.", target)
		oneFindingOB(t, "ALTER TABLE t PARTITION BY KEY()(PARTITION p1, PARTITION p1);", 0, "This statement declares partition `p1` twice.", target)
		oneFindingOB(t, "ALTER TABLE t ALTER PARTITION p0 STORAGE_CACHE_POLICY = 's';", 0, "OceanBase has no storage cache policy for a partition.", target)
		oneFindingOB(t, "ALTER TABLE t EXCHANGE PARTITION p0 WITH TABLE t WITHOUT VALIDATION;", 0, "OceanBase refuses exchanging a partition of `t` with itself.", target)
	})

	t.Run("triggers and events", func(t *testing.T) {
		oneFindingOB(t, "CREATE TRIGGER otherdb.tr_bad BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;", 0, "Trigger `otherdb.tr_bad` and table `t` are in different databases.", target)
		oneFindingOB(t, "CREATE TRIGGER tr_ord BEFORE INSERT ON t FOR EACH ROW FOLLOWS nosuch SET NEW.a = 1;", 0, "Trigger `nosuch` doesn't exist.", target)
		noFindingOB(t, "CREATE TRIGGER tr_ord BEFORE INSERT ON t FOR EACH ROW FOLLOWS tr_existing SET NEW.a = 1;", target)
		oneFindingOB(t, "CREATE TRIGGER tr_old BEFORE INSERT ON t FOR EACH ROW SET @x = OLD.a;", 0, "OceanBase has no OLD row in an INSERT trigger.", target)
		oneFindingOB(t, "CREATE TRIGGER tr_new BEFORE DELETE ON t FOR EACH ROW SET @x = NEW.a;", 0, "OceanBase has no NEW row in a DELETE trigger.", target)
		oneFindingOB(t, "CREATE TRIGGER tr_after AFTER INSERT ON t FOR EACH ROW SET NEW.a = 1;", 0, "OceanBase refuses updating NEW row in an AFTER trigger.", target)
		noFindingOB(t, "CREATE TRIGGER tr_alias BEFORE INSERT ON t FOR EACH ROW BEGIN DECLARE x INT; SELECT old.id INTO x FROM parent AS old LIMIT 1; END;", target)
		oneFindingOB(t, "CREATE EVENT ev_bad ON SCHEDULE AT a DO COMMIT;", 0, "OceanBase refuses column reference `a` in an event schedule.", target)
		noFindingOB(t, "CREATE EVENT ev_ok ON SCHEDULE AT '2030-01-01 00:00:00' DO COMMIT;", target)
		oneFindingOB(t, "ALTER EVENT ev_ok RENAME TO ev_ok;", 0, "Event `ev_ok` cannot be renamed to itself.", target)
	})
}
