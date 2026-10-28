package review

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkSchema is the synced shop database the unit tests walk on, as
// Bytebase's sync since #21665 lists it: a sequence s, and a
// system-versioned table vt, marked in its create options.
func walkSchema() *metadata.DatabaseSchemaMetadata {
	cols := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	return &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name:        "t",
					Columns:     cols("id", "a"),
					Indexes:     []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}, {Name: "i"}},
					ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk1", ReferencedTable: "r"}},
					Triggers:    []*metadata.TriggerMetadata{{Name: "tr"}},
				},
				{Name: "r", Columns: cols("id"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}}},
				{Name: "nokey", Columns: cols("a")},
				{Name: "vt", Columns: cols("id", "a"), CreateOptions: "WITH SYSTEM VERSIONING", Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}}},
			},
			Views:      []*metadata.ViewMetadata{{Name: "v"}},
			Sequences:  []*metadata.SequenceMetadata{{Name: "s"}},
			Procedures: []*metadata.ProcedureMetadata{{Name: "p"}},
			Functions:  []*metadata.FunctionMetadata{{Name: "f"}},
			Events:     []*metadata.EventMetadata{{Name: "e"}},
		}},
	}
}

func walkFindings(t *testing.T, sql string, targets ...review.Target) []review.Finding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.WalkThrough}}, targets)
	if err != nil {
		t.Fatal(err)
	}
	return result.Findings
}

// oneFinding asserts the change's one WALK_THROUGH finding.
func oneFinding(t *testing.T, sql string, statement int, message string, target review.Target) {
	t.Helper()
	got := walkFindings(t, sql, target)
	if len(got) != 1 || got[0].Statement != statement || got[0].Message != message {
		t.Errorf("%s\nfindings = %+v, want statement %d: %q", sql, got, statement, message)
	}
}

func noFinding(t *testing.T, sql string, target review.Target) {
	t.Helper()
	if got := walkFindings(t, sql, target); len(got) != 0 {
		t.Errorf("%s\nfindings = %+v, want none", sql, got)
	}
}

func TestWalkThroughMariaDB(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table t already exists", shop)
	oneFinding(t, "ALTER TABLE v ADD COLUMN b int;", 0, "v is not a base table", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN A int;", 0, "Column A already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column zz does not exist in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX I (a);", 0, "Index I already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "Table t already has a primary key", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);", 0, "Foreign key FK1 already exists in the database", shop)
	oneFinding(t, "DROP VIEW t;", 0, "t is not a view", shop)
	oneFinding(t, "DROP TABLE v;", 0, "v is a view", shop)
	oneFinding(t, "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", 0, "Trigger tr already exists", shop)
	oneFinding(t, "DROP PROCEDURE nosuch;", 0, "PROCEDURE nosuch does not exist", shop)
	oneFinding(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 0, "Event e already exists", shop)
	// ALTER TABLE judges drops against the table as it was.
	noFinding(t, "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN b int, ADD COLUMN b int;", 0, "Column b is added twice to table t", shop)
	// RENAME TABLE runs its pairs in order.
	noFinding(t, "RENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN zz int;\nALTER TABLE t ADD COLUMN zz int;", shop)
	oneFinding(t, "RENAME TABLE t TO u;\nCREATE TABLE u (a int);", 1, "Table u already exists", shop)
	// MariaDB has no invisible primary key.
	oneFinding(t, "ALTER TABLE nokey DROP PRIMARY KEY;", 0, "Table nokey has no primary key", shop)
	oneFinding(t, "CREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;", 1, "Table g has no primary key", shop)
}

// TestWalkThroughMariaDBAuthority pins what the sync is not trusted for:
// the relations and triggers it does not list (decision N2), and the
// columns of a system-versioned table.
func TestWalkThroughMariaDBAuthority(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "ALTER TABLE nosuch ADD COLUMN b int;", shop)
	noFinding(t, "DROP TABLE t, nosuch;", shop)
	noFinding(t, "DROP VIEW nosuch;", shop)
	noFinding(t, "TRUNCATE TABLE nosuch;", shop)
	noFinding(t, "INSERT INTO nosuch VALUES (1);", shop)
	noFinding(t, "DROP TRIGGER nosuch;", shop)
	// What the change drops is followed.
	oneFinding(t, "DROP TABLE nokey;\nDROP TABLE nokey;", 1, "Table nokey does not exist", shop)
	oneFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER tr;", 1, "Trigger tr does not exist", shop)
	// A system-versioned table hides row_start and row_end, and an ALTER
	// TABLE of one ends the walk.
	noFinding(t, "ALTER TABLE vt ADD COLUMN x int AFTER row_start;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE TABLE w (a int) WITH SYSTEM VERSIONING;\nALTER TABLE w ADD COLUMN x int AFTER row_end;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE TABLE w (a int WITH SYSTEM VERSIONING, b int);\nALTER TABLE w DROP COLUMN row_start;", shop)
	noFinding(t, "CREATE TABLE w LIKE vt;\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "RENAME TABLE vt TO w;\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nCREATE TABLE t (a int);", shop)
	// Its indexes are listed as any table's.
	oneFinding(t, "DROP INDEX nosuch ON vt;", 0, "Index nosuch does not exist in table vt", shop)
	// A table replaced without versioning is followed again.
	oneFinding(t, "CREATE OR REPLACE TABLE vt (a int);\nALTER TABLE vt DROP COLUMN row_start;", 1, "Column row_start does not exist in table vt", shop)
}

func TestWalkThroughMariaDBOrReplace(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "CREATE OR REPLACE TABLE t (id int PRIMARY KEY, CONSTRAINT fk1 FOREIGN KEY (id) REFERENCES r (id));", shop)
	noFinding(t, "CREATE OR REPLACE TABLE s (a int);\nALTER TABLE s ADD COLUMN b int;", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE v (a int);", 0, "v is a view", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE nokey (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));", 0, "Foreign key fk1 already exists in the database", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE t (b int);\nALTER TABLE t DROP COLUMN a;", 1, "Column a does not exist in table t", shop)
	noFinding(t, "CREATE OR REPLACE INDEX i ON t (a);", shop)
	noFinding(t, "CREATE OR REPLACE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;", shop)
	noFinding(t, "CREATE OR REPLACE PROCEDURE p() SELECT 1;\nCREATE OR REPLACE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;", shop)
	noFinding(t, "CREATE OR REPLACE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", shop)
	// OR REPLACE of the current database leaves none selected.
	oneFinding(t, "CREATE OR REPLACE DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database selected", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE OR REPLACE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 4, "Table q.x already exists", shop)
}

func TestWalkThroughMariaDBIfExists(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS zz, DROP INDEX IF EXISTS zz, DROP FOREIGN KEY IF EXISTS zz, DROP CONSTRAINT IF EXISTS zz;", shop)
	noFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS a, DROP COLUMN IF EXISTS a;", shop)
	noFinding(t, "ALTER TABLE t DROP INDEX IF EXISTS `PRIMARY`;\nALTER TABLE t ADD PRIMARY KEY (a);", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS a;\nALTER TABLE t DROP COLUMN a;", 1, "Column a does not exist in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP INDEX IF EXISTS i;\nDROP INDEX i ON t;", 1, "Index i does not exist in table t", shop)
	// Another table's foreign key of the name stays.
	oneFinding(t, "ALTER TABLE nokey DROP FOREIGN KEY IF EXISTS fk1;\nALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);", 1, "Foreign key fk1 already exists in the database", shop)
	noFinding(t, "ALTER TABLE t DROP FOREIGN KEY IF EXISTS fk1;\nALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);", shop)
}

func TestWalkThroughMariaDBCreateTable(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// A temporary table needs no database.
	noFinding(t, "CREATE DATABASE q;\nDROP DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nDROP DATABASE q;\nCREATE TABLE q.x (a int);", 2, "Database q does not exist", shop)
	// IF NOT EXISTS of a table that exists reads no LIKE source; of a
	// temporary table it does.
	noFinding(t, "DROP TABLE r;\nCREATE TABLE IF NOT EXISTS t LIKE r;", shop)
	oneFinding(t, "DROP TABLE r;\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp LIKE r;", 2, "Table r does not exist", shop)
	oneFinding(t, "DROP TABLE r;\nCREATE TABLE x LIKE r;", 1, "Table r does not exist", shop)
	// LIKE of a sequence makes a sequence, and of a view fails.
	noFinding(t, "CREATE TABLE x LIKE s;\nDROP SEQUENCE x;", shop)
	oneFinding(t, "CREATE TABLE x LIKE v;", 0, "v is not a base table", shop)
}

func TestWalkThroughMariaDBSequences(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "CREATE SEQUENCE t;", 0, "Table t already exists", shop)
	oneFinding(t, "CREATE TABLE s (a int);", 0, "Table s already exists", shop)
	oneFinding(t, "CREATE VIEW s AS SELECT 1 AS a;", 0, "Table s already exists", shop)
	oneFinding(t, "CREATE OR REPLACE VIEW s AS SELECT 1 AS a;", 0, "s is not a view", shop)
	oneFinding(t, "CREATE SEQUENCE x;\nCREATE SEQUENCE x;", 1, "Table x already exists", shop)
	noFinding(t, "CREATE SEQUENCE IF NOT EXISTS t;\nCREATE OR REPLACE SEQUENCE t;\nDROP SEQUENCE t;", shop)
	oneFinding(t, "CREATE OR REPLACE SEQUENCE v;", 0, "v is a view", shop)
	oneFinding(t, "DROP SEQUENCE t;", 0, "t is not a sequence", shop)
	oneFinding(t, "DROP SEQUENCE v;", 0, "v is a view", shop)
	oneFinding(t, "DROP SEQUENCE s;\nDROP SEQUENCE s;", 1, "Sequence s does not exist", shop)
	oneFinding(t, "DROP VIEW s;", 0, "s is not a view", shop)
	noFinding(t, "DROP SEQUENCE IF EXISTS t, v, nosuch;", shop)
	noFinding(t, "DROP TABLE s;\nCREATE TABLE s (a int);", shop)
	noFinding(t, "RENAME TABLE s TO s2;\nDROP SEQUENCE s2;", shop)
	// A sequence in ALTER TABLE, TRUNCATE, CREATE INDEX, or CREATE TRIGGER
	// ends the walk.
	noFinding(t, "ALTER TABLE s COMMENT 'x';\nCREATE TABLE t (a int);", shop)
	noFinding(t, "TRUNCATE TABLE s;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE INDEX k ON s (cache_size);\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE TRIGGER trs BEFORE INSERT ON s FOR EACH ROW SET @x = 1;\nCREATE TABLE t (a int);", shop)
	// Temporary sequences: DROP SEQUENCE drops one before a base one,
	// passes over a temporary table, and DROP TEMPORARY SEQUENCE finds
	// only a temporary sequence.
	noFinding(t, "CREATE TEMPORARY SEQUENCE s;\nDROP SEQUENCE s;\nDROP SEQUENCE s;", shop)
	oneFinding(t, "CREATE TEMPORARY SEQUENCE s;\nDROP SEQUENCE s;\nDROP SEQUENCE s;\nDROP SEQUENCE s;", 3, "Sequence s does not exist", shop)
	noFinding(t, "CREATE TEMPORARY TABLE s (a int);\nDROP SEQUENCE s;\nDROP TEMPORARY TABLE s;", shop)
	oneFinding(t, "DROP TEMPORARY SEQUENCE s;", 0, "Temporary sequence s does not exist", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nDROP TEMPORARY SEQUENCE tt;", 1, "Temporary table tt is not a sequence", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY SEQUENCE tt;", 1, "Table tt already exists", shop)
	noFinding(t, "CREATE TEMPORARY SEQUENCE q.ts;", shop)
}

// TestWalkThroughMariaDBMySQLForms pins that the walk ends at the MySQL
// ALTER TABLE forms omni's MariaDB parser takes and MariaDB rejects.
func TestWalkThroughMariaDBMySQLForms(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range []string{
		"ALTER TABLE t ALTER INDEX nosuch INVISIBLE;",
		"ALTER TABLE t ALTER COLUMN nosuch SET NOT NULL;",
		"ALTER TABLE t ALTER COLUMN nosuch SET INVISIBLE;",
		"ALTER TABLE t ALTER CHECK c NOT ENFORCED;\nCREATE TABLE t (a int);",
		"ALTER TABLE t SECONDARY_LOAD;\nCREATE TABLE t (a int);",
	} {
		noFinding(t, sql, shop)
	}
}

func TestWalkThroughMariaDBCase(t *testing.T) {
	sensitive := review.Target{Schema: walkSchema()}
	folded := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1}
	noFinding(t, "CREATE TABLE T (a int);", sensitive)
	oneFinding(t, "CREATE TABLE T (a int);", 0, "Table T already exists", folded)
	// Trigger names compare with case: TR beside tr is a second trigger.
	noFinding(t, "CREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", sensitive)
	noFinding(t, "CREATE TRIGGER TR BEFORE UPDATE ON nokey FOR EACH ROW SET NEW.a = 1;\nDROP TRIGGER TR;\nDROP TRIGGER tr;", sensitive)
	oneFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER tr;", 1, "Trigger tr does not exist", sensitive)
}

// TestWalkThroughMariaDBNamesOutsideASCII pins that the walk follows no
// name outside ASCII, which MariaDB compares by rules of its own per
// space: a statement that names one ends the walk, and a synced space
// that lists one has no name known to be absent, since ſ is s as a column
// and é is e as a routine.
func TestWalkThroughMariaDBNamesOutsideASCII(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "ALTER TABLE t DROP COLUMN `é`;", shop)
	noFinding(t, "CREATE TABLE x (id INT, café INT);\nALTER TABLE x DROP COLUMN CAFé;", shop)
	noFinding(t, "CREATE TABLE x (s INT);\nALTER TABLE x DROP COLUMN ſ;\nALTER TABLE x DROP COLUMN zz;", shop)
	noFinding(t, "CREATE DATABASE q;\nCREATE PROCEDURE q.`pé`() SELECT 1;\nDROP PROCEDURE q.pe;", shop)
	// A string is no name.
	oneFinding(t, "CREATE TABLE x (a int) COMMENT '中文';\nCREATE TABLE x (a int);", 1, "Table x already exists", shop)
	oneFinding(t, "INSERT INTO nokey VALUES ('é');\nDROP TEMPORARY TABLE nosuch;", 1, "Temporary table nosuch does not exist", shop)
	mixed := walkSchema()
	schema := mixed.Schemas[0]
	schema.Tables[2].Columns = append(schema.Tables[2].Columns, &metadata.ColumnMetadata{Name: "ſ"})
	schema.Procedures = append(schema.Procedures, &metadata.ProcedureMetadata{Name: "pé"})
	target := review.Target{Schema: mixed}
	noFinding(t, "ALTER TABLE nokey DROP COLUMN s;\nALTER TABLE nokey DROP COLUMN zz;\nDROP PROCEDURE pe;\nDROP PROCEDURE nosuch;", target)
	// A space that lists none still knows its absent names.
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column zz does not exist in table t", target)
	oneFinding(t, "DROP EVENT nosuch;", 0, "Event nosuch does not exist", target)
}

func TestWalkThroughMariaDBSession(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Another database is unknown, until the change makes it.
	noFinding(t, "USE other;\nCREATE TABLE t (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 2, "Table q.x already exists", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database selected", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TEMPORARY TABLE x (a int);", 1, "No database selected", shop)
	// A temporary table shadows the base table.
	noFinding(t, "DROP TABLE nokey;\nCREATE TEMPORARY TABLE nokey (a int);\nINSERT INTO nokey VALUES (1);", shop)
	oneFinding(t, "DROP TEMPORARY TABLE t;", 0, "Temporary table t does not exist", shop)
	// A trigger may create a temporary table of any name.
	noFinding(t, "CREATE TEMPORARY TABLE x (a int);\nDROP TEMPORARY TABLE x;\nINSERT INTO t VALUES (1, 1);\nDROP TEMPORARY TABLE x;", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE x (a int);\nDROP TEMPORARY TABLE x;\nINSERT INTO nokey VALUES (1);\nDROP TEMPORARY TABLE x;", 3, "Temporary table x does not exist", shop)
	// A foreign key's target is checked only while the checks are on.
	noFinding(t, "DROP TABLE r;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES r (id));", shop)
	oneFinding(t, "DROP TABLE nokey;\nSET foreign_key_checks = 1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nokey (a));", 2, "Foreign key target table nokey does not exist", shop)
	oneFinding(t, "SET foreign_key_checks = 1;\nCREATE TABLE x (a bigint, FOREIGN KEY (a) REFERENCES s (next_not_cached_value));", 1, "Foreign key target s is not a base table", shop)
	// Each assignment has a scope of its own: an unprefixed one the
	// statement's, past a prefixed one, and an @@session or bare @@ one the
	// session's. A global one leaves the session's checks as they are.
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	const fk = "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nokey (a));"
	const missing = "Foreign key target table nokey does not exist"
	oneFinding(t, "DROP TABLE nokey;\nSET GLOBAL foreign_key_checks = 0;\n"+fk, 2, missing, on)
	oneFinding(t, "DROP TABLE nokey;\nSET @@global.foreign_key_checks = 0;\n"+fk, 2, missing, on)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL foreign_key_checks = 1;\n"+fk, shop)
	oneFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, 2, missing, on)
	oneFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, @@session.foreign_key_checks = 1, foreign_key_checks = 0;\n"+fk, 2, missing, shop)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, @@session.foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, @@foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "DROP TABLE nokey;\nSET @@global.foreign_key_checks = 1, foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "DROP TABLE nokey;\nSET LOCAL foreign_key_checks = 0;\n"+fk, on)
	// With another spec of the plan, only the change's own facts count.
	other := review.Target{Schema: walkSchema(), OtherSpecs: true}
	noFinding(t, "CREATE TABLE t (a int);", other)
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE x (a int);", 1, "Table x already exists", other)
}

// TestWalkThroughMariaDBForeignKeyNames pins the name a foreign key and
// its index take: the CONSTRAINT symbol, or without one the name after
// FOREIGN KEY, which counts as long only then (verified on 11.8.9).
func TestWalkThroughMariaDBForeignKeyNames(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	long := strings.Repeat("f", 65)
	noFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY cx;\nALTER TABLE x DROP INDEX cx;", shop)
	oneFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY fx;", 1, "Foreign key fx does not exist in table x", shop)
	noFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nALTER TABLE y DROP FOREIGN KEY fy;", shop)
	oneFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nCREATE TABLE z (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));", 1, "Foreign key fy already exists in the database", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY fk1 (a) REFERENCES r (id);", 0, "Foreign key fk1 already exists in the database", shop)
	noFinding(t, "CREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id));", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id);", shop)
	oneFinding(t, "CREATE TABLE w (a int, FOREIGN KEY "+long+" (a) REFERENCES r (id));", 0, "Identifier name "+long+" is too long", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY "+long+" (a) REFERENCES r (id);", 0, "Identifier name "+long+" is too long", shop)
}

// TestWalkThroughMariaDBNoDatabase pins what fails once the change drops
// the current database: an unqualified name (1046), but for DROP
// FUNCTION, which then looks for a loadable function alone (verified on
// 11.8.9).
func TestWalkThroughMariaDBNoDatabase(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS nosuch;\nDROP FUNCTION nosuch;\nDROP FUNCTION IF EXISTS shop.f;", shop)
	for _, stmt := range []string{"DROP PROCEDURE IF EXISTS p;", "DROP TRIGGER IF EXISTS tr;", "DROP VIEW IF EXISTS v;", "DROP EVENT IF EXISTS e;", "DROP SEQUENCE IF EXISTS s;", "CREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;", "ALTER FUNCTION f COMMENT 'x';"} {
		oneFinding(t, "DROP DATABASE shop;\n"+stmt, 1, "No database selected", shop)
	}
}

// TestWalkThroughMariaDBParentBesideRename pins that an ALTER TABLE that
// also renames its table checks no foreign key's parent: MariaDB resolves
// it in the database the table moves to (verified on 11.8.9).
func TestWalkThroughMariaDBParentBesideRename(t *testing.T) {
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	const gone = "CREATE TABLE p2 (id int PRIMARY KEY);\nDROP TABLE p2;\n"
	noFinding(t, gone+"CREATE DATABASE q;\nCREATE TABLE q.p2 (id int PRIMARY KEY);\nALTER TABLE nokey RENAME TO q.nokey, ADD FOREIGN KEY (a) REFERENCES p2 (id);", on)
	noFinding(t, gone+"ALTER TABLE nokey RENAME TO nk, ADD FOREIGN KEY (a) REFERENCES p2 (id);", on)
	oneFinding(t, gone+"ALTER TABLE nokey ADD FOREIGN KEY (a) REFERENCES p2 (id);", 2, "Foreign key target table p2 does not exist", on)
}

// TestWalkThroughMariaDBStoredFunctions pins the calls that may run a
// stored function, which may create a temporary table (verified on
// 11.8.9): one a view's query makes, and an unqualified one after a USE of
// a database the change did not make.
func TestWalkThroughMariaDBStoredFunctions(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	const drop = "DROP TEMPORARY TABLE made;"
	noFinding(t, "SELECT * FROM v;\n"+drop, shop)
	noFinding(t, "INSERT INTO nokey SELECT a FROM q.qv;\n"+drop, shop)
	noFinding(t, "USE q;\nSELECT mk();\nUSE shop;\n"+drop, shop)
	noFinding(t, "SELECT f();\n"+drop, shop)
	// A call found stays found past a base table read after it.
	noFinding(t, "SELECT f() FROM r;\n"+drop, shop)
	noFinding(t, "SELECT q.mk() FROM nokey WHERE a IN (SELECT id FROM r);\n"+drop, shop)
	// A base table or a sequence runs no function, nor does a call to one
	// the current database certainly lacks.
	oneFinding(t, "SELECT * FROM nokey;\n"+drop, 1, "Temporary table made does not exist", shop)
	oneFinding(t, "SELECT NEXTVAL(s) FROM s;\n"+drop, 1, "Temporary table made does not exist", shop)
	oneFinding(t, "SELECT mk();\n"+drop, 1, "Temporary table made does not exist", shop)
}

// TestWalkThroughMariaDBSequenceRenamed pins that a relation keeps its
// kind through ALTER TABLE ... RENAME: one made LIKE another database's
// sequence may be a sequence, and stays one (verified on 11.8.9).
func TestWalkThroughMariaDBSequenceRenamed(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "CREATE TABLE t2 LIKE q.s;\nALTER TABLE t2 RENAME TO t3;\nDROP SEQUENCE t3;", shop)
	noFinding(t, "CREATE TABLE t2 LIKE q.s;\nALTER TABLE t2 RENAME TO q.t3, COMMENT 'x';\nDROP SEQUENCE q.t3;", shop)
	oneFinding(t, "CREATE TABLE t2 (a int);\nALTER TABLE t2 RENAME TO t3;\nDROP SEQUENCE t3;", 2, "t3 is not a sequence", shop)
}

// TestWalkThroughMariaDBTargetGroups pins OtherSpecs and ForeignKeyChecks
// in the target grouping: two targets of one schema that differ in either
// are walked apart.
func TestWalkThroughMariaDBTargetGroups(t *testing.T) {
	schema := walkSchema()
	got := walkFindings(t, "CREATE TABLE t (a int);", review.Target{Schema: schema}, review.Target{Schema: schema, OtherSpecs: true})
	if len(got) != 1 || len(got[0].Targets) != 1 || got[0].Targets[0] != 0 {
		t.Errorf("OtherSpecs: findings = %+v, want one on target 0", got)
	}
	const fk = "CREATE TABLE gone (id int PRIMARY KEY);\nDROP TABLE gone;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES gone (id));"
	got = walkFindings(t, fk, review.Target{Schema: schema}, review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn})
	if len(got) != 1 || len(got[0].Targets) != 1 || got[0].Targets[0] != 1 {
		t.Errorf("ForeignKeyChecks: findings = %+v, want one on target 1", got)
	}
}

// TestWalkThroughMariaDBOldSync walks a target as Bytebase synced it
// before #21665, which listed no sequence and no system-versioned table:
// such names are unknown, never absent.
func TestWalkThroughMariaDBOldSync(t *testing.T) {
	old := walkSchema()
	old.Schemas[0].Sequences = nil
	old.Schemas[0].Tables = old.Schemas[0].Tables[:3]
	target := review.Target{Schema: old}
	noFinding(t, "DROP SEQUENCE s;\nDROP TABLE vt;\nCREATE TABLE s (a int);", target)
	noFinding(t, "CREATE TABLE IF NOT EXISTS vt (a int);\nALTER TABLE vt ADD COLUMN b int AFTER row_start;", target)
}

// walkForms are the statement forms of the MySQL spec's §5 table and the
// MariaDB spec's §2.3, each on walkSchema. A form the tables list must
// lower.
var walkForms = []walktest.Form{
	{SQL: "CREATE DATABASE q;"},
	{SQL: "CREATE DATABASE IF NOT EXISTS shop;"},
	{SQL: "CREATE OR REPLACE DATABASE q;"},
	{SQL: "DROP DATABASE IF EXISTS q;"},
	{SQL: "ALTER DATABASE shop CHARACTER SET utf8mb4;"},
	{SQL: "USE shop;"},
	{SQL: "CREATE TABLE x (a int PRIMARY KEY, b int UNIQUE, KEY k (b), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id), CHECK (a > 0));"},
	{SQL: "CREATE TABLE IF NOT EXISTS t (a int);"},
	{SQL: "CREATE OR REPLACE TABLE t (a int);"},
	{SQL: "CREATE TABLE x (a int) WITH SYSTEM VERSIONING;"},
	{SQL: "CREATE TEMPORARY TABLE tmp (a int);"},
	{SQL: "CREATE OR REPLACE TEMPORARY TABLE tmp (a int);"},
	{SQL: "CREATE TABLE x LIKE t;"},
	{SQL: "CREATE TABLE x AS SELECT 1 AS a;"},
	{SQL: "ALTER TABLE t ADD COLUMN b int AFTER a, ADD COLUMN (c int, d int);"},
	{SQL: "ALTER TABLE t DROP COLUMN a;"},
	{SQL: "ALTER TABLE t DROP COLUMN IF EXISTS a;"},
	{SQL: "ALTER TABLE t MODIFY COLUMN a bigint;"},
	{SQL: "ALTER TABLE t CHANGE a b int;"},
	{SQL: "ALTER TABLE t RENAME COLUMN a TO b;"},
	{SQL: "ALTER TABLE t ALTER COLUMN a SET DEFAULT 1;"},
	{SQL: "ALTER TABLE t ADD INDEX k (a), ADD UNIQUE (a), ADD FULLTEXT INDEX ft (a);"},
	{SQL: "ALTER TABLE nokey ADD PRIMARY KEY (a);"},
	{SQL: "ALTER TABLE nokey ADD CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES r (id);"},
	{SQL: "ALTER TABLE t ADD CONSTRAINT c CHECK (a > 0);"},
	{SQL: "ALTER TABLE t DROP INDEX i, DROP KEY i2;"},
	{SQL: "ALTER TABLE t DROP INDEX IF EXISTS i, DROP KEY IF EXISTS i2;"},
	{SQL: "ALTER TABLE t DROP PRIMARY KEY;"},
	{SQL: "ALTER TABLE t DROP FOREIGN KEY fk1;"},
	{SQL: "ALTER TABLE t DROP FOREIGN KEY IF EXISTS fk1;"},
	{SQL: "ALTER TABLE t DROP CONSTRAINT c;"},
	{SQL: "ALTER TABLE t DROP CONSTRAINT IF EXISTS c;"},
	{SQL: "ALTER TABLE t RENAME INDEX i TO j;"},
	{SQL: "ALTER TABLE t RENAME TO t2;"},
	{SQL: "ALTER TABLE t ENGINE = InnoDB, COMMENT 'x', ALGORITHM = INPLACE, LOCK = NONE;"},
	{SQL: "ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4;"},
	{SQL: "ALTER TABLE t FORCE;"},
	{SQL: "RENAME TABLE t TO t2, r TO r2;"},
	{SQL: "DROP TABLE IF EXISTS t, nosuch;"},
	{SQL: "DROP TABLE s;"},
	{SQL: "DROP TEMPORARY TABLE IF EXISTS tmp;"},
	{SQL: "TRUNCATE TABLE t;"},
	{SQL: "CREATE SEQUENCE x;"},
	{SQL: "CREATE OR REPLACE SEQUENCE s;"},
	{SQL: "CREATE SEQUENCE IF NOT EXISTS s;"},
	{SQL: "CREATE TEMPORARY SEQUENCE x;"},
	{SQL: "DROP SEQUENCE s;"},
	{SQL: "DROP SEQUENCE IF EXISTS s, nosuch;"},
	{SQL: "DROP TEMPORARY SEQUENCE IF EXISTS x;"},
	{SQL: "CREATE UNIQUE INDEX k ON t (a);"},
	{SQL: "CREATE OR REPLACE INDEX i ON t (a);"},
	{SQL: "CREATE INDEX IF NOT EXISTS i ON t (a);"},
	{SQL: "DROP INDEX i ON t;"},
	{SQL: "CREATE VIEW w AS SELECT 1 AS a;"},
	{SQL: "CREATE OR REPLACE VIEW v AS SELECT 1 AS a;"},
	{SQL: "ALTER VIEW v AS SELECT 2 AS a;"},
	{SQL: "DROP VIEW IF EXISTS v, w;"},
	{SQL: "CREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;"},
	{SQL: "CREATE OR REPLACE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;"},
	{SQL: "DROP TRIGGER IF EXISTS tr;"},
	{SQL: "CREATE PROCEDURE p2() SELECT 1;"},
	{SQL: "CREATE OR REPLACE PROCEDURE p() SELECT 1;"},
	{SQL: "CREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;"},
	{SQL: "ALTER PROCEDURE p COMMENT 'x';"},
	{SQL: "DROP FUNCTION IF EXISTS f;"},
	{SQL: "CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;"},
	{SQL: "CREATE OR REPLACE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;"},
	{SQL: "ALTER EVENT e RENAME TO e2;"},
	{SQL: "DROP EVENT IF EXISTS e;"},
	{SQL: "INSERT INTO nokey VALUES (1);"},
	{SQL: "REPLACE INTO nokey VALUES (1);"},
	{SQL: "UPDATE nokey SET a = 1;"},
	{SQL: "DELETE FROM nokey;"},
	{SQL: "SELECT 1;"},
	{SQL: "SET foreign_key_checks = 0;\nSET @@session.foreign_key_checks = ON;\nSET NAMES utf8mb4;"},
	{SQL: "BEGIN;\nSAVEPOINT a;\nRELEASE SAVEPOINT a;\nCOMMIT;"},
	{SQL: "SHOW TABLES;\nANALYZE TABLE t;\nOPTIMIZE TABLE t;\nCHECK TABLE t;"},
	{SQL: "GRANT SELECT ON shop.t TO 'u'@'%';"},
}

func TestWalkThroughMariaDBLowersEveryForm(t *testing.T) {
	walktest.Coverage(t, walkForms, func(t *testing.T, sql string) (n, lowered int, last bool) {
		t.Helper()
		ranges := statementRanges(sql)
		stmts, syntax := parse(sql, ranges)
		if syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, sql)
		}
		w := newWalker(sql, stmts, review.Target{Schema: walkSchema()})
		l := walk.NewLedger(w.snap)
		for i := range stmts {
			s := w.lower(i, l)
			if s == nil {
				return len(stmts), i, false
			}
			if s.Last {
				return len(stmts), i + 1, true
			}
			if s.Apply != nil {
				s.Apply(l)
			}
		}
		return len(stmts), len(stmts), false
	})
}
