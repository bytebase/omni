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
	for _, f := range result.Findings {
		if err := review.CheckCodes(f); err != nil {
			t.Error(err)
		}
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
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table `t` already exists.", shop)
	oneFinding(t, "ALTER TABLE v ADD COLUMN b int;", 0, "`v` is a view, not a table.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN A int;", 0, "Column `A` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column `zz` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX I (a);", 0, "Index `I` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "`t` already has a primary key.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);", 0, "Foreign key `FK1` already exists.", shop)
	oneFinding(t, "DROP VIEW t;", 0, "`t` is a table, not a view.", shop)
	oneFinding(t, "DROP TABLE v;", 0, "`v` is a view, not a table.", shop)
	oneFinding(t, "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", 0, "Trigger `tr` already exists.", shop)
	oneFinding(t, "DROP PROCEDURE nosuch;", 0, "Procedure `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 0, "Event `e` already exists.", shop)
	// ALTER TABLE judges drops against the table as it was.
	noFinding(t, "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN b int, ADD COLUMN b int;", 0, "This statement adds column `b` twice.", shop)
	// RENAME TABLE runs its pairs in order.
	noFinding(t, "RENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN zz int;\nALTER TABLE t ADD COLUMN zz int;", shop)
	oneFinding(t, "RENAME TABLE t TO u;\nCREATE TABLE u (a int);", 1, "Table `u` already exists.", shop)
	// MariaDB has no invisible primary key.
	oneFinding(t, "ALTER TABLE nokey DROP PRIMARY KEY;", 0, "`nokey` has no primary key.", shop)
	oneFinding(t, "CREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;", 1, "`g` has no primary key.", shop)
}

// TestWalkThroughMariaDBAuthority pins what the sync is trusted for: the
// relations, sequences and system-versioned tables among them, and the
// triggers it does not list are missing; the columns of a
// system-versioned table are not all listed.
func TestWalkThroughMariaDBAuthority(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE nosuch ADD COLUMN b int;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE t, nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP VIEW nosuch;", 0, "View `nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP SEQUENCE nosuch;", 0, "Sequence `nosuch` doesn't exist.", shop)
	oneFinding(t, "TRUNCATE TABLE nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "INSERT INTO nosuch VALUES (1);", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP TRIGGER nosuch;", 0, "Trigger `nosuch` doesn't exist.", shop)
	noFinding(t, "DROP SEQUENCE s;\nDROP TABLE vt;", shop)
	// What the change drops is followed.
	oneFinding(t, "DROP TABLE nokey;\nDROP TABLE nokey;", 1, "Table `nokey` doesn't exist.", shop)
	oneFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER tr;", 1, "Trigger `tr` doesn't exist.", shop)
	// A system-versioned table hides row_start and row_end, and an ALTER
	// TABLE of one ends the walk.
	noFinding(t, "ALTER TABLE vt ADD COLUMN x int AFTER row_start;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE TABLE w (a int) WITH SYSTEM VERSIONING;\nALTER TABLE w ADD COLUMN x int AFTER row_end;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE TABLE w (a int WITH SYSTEM VERSIONING, b int);\nALTER TABLE w DROP COLUMN row_start;", shop)
	noFinding(t, "CREATE TABLE w LIKE vt;\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "RENAME TABLE vt TO w;\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nCREATE TABLE t (a int);", shop)
	// Its indexes are listed as any table's.
	oneFinding(t, "DROP INDEX nosuch ON vt;", 0, "Index `nosuch` doesn't exist in `vt`.", shop)
	// A table replaced without versioning is followed again.
	oneFinding(t, "CREATE OR REPLACE TABLE vt (a int);\nALTER TABLE vt DROP COLUMN row_start;", 1, "Column `row_start` doesn't exist in `vt`.", shop)
}

func TestWalkThroughMariaDBOrReplace(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "CREATE OR REPLACE TABLE t (id int PRIMARY KEY, CONSTRAINT fk1 FOREIGN KEY (id) REFERENCES r (id));", shop)
	noFinding(t, "CREATE OR REPLACE TABLE s (a int);\nALTER TABLE s ADD COLUMN b int;", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE v (a int);", 0, "`v` is a view, not a table.", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE nokey (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));", 0, "Foreign key `fk1` already exists.", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE t (b int);\nALTER TABLE t DROP COLUMN a;", 1, "Column `a` doesn't exist in `t`.", shop)
	noFinding(t, "CREATE OR REPLACE INDEX i ON t (a);", shop)
	noFinding(t, "CREATE OR REPLACE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;", shop)
	noFinding(t, "CREATE OR REPLACE PROCEDURE p() SELECT 1;\nCREATE OR REPLACE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;", shop)
	noFinding(t, "CREATE OR REPLACE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", shop)
	// OR REPLACE of the current database leaves none selected.
	oneFinding(t, "CREATE OR REPLACE DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database is selected for `x`.", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE OR REPLACE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 4, "Table `q.x` already exists.", shop)
}

func TestWalkThroughMariaDBIfExists(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS zz, DROP INDEX IF EXISTS zz, DROP FOREIGN KEY IF EXISTS zz, DROP CONSTRAINT IF EXISTS zz;", shop)
	noFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS a, DROP COLUMN IF EXISTS a;", shop)
	noFinding(t, "ALTER TABLE t DROP INDEX IF EXISTS `PRIMARY`;\nALTER TABLE t ADD PRIMARY KEY (a);", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS a;\nALTER TABLE t DROP COLUMN a;", 1, "Column `a` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t DROP INDEX IF EXISTS i;\nDROP INDEX i ON t;", 1, "Index `i` doesn't exist in `t`.", shop)
	// The standalone DROP INDEX reads IF EXISTS too, and it speaks for the
	// index, not for the table.
	noFinding(t, "DROP INDEX IF EXISTS nosuch ON t;", shop)
	noFinding(t, "DROP INDEX i ON t;\nDROP INDEX IF EXISTS i ON t;", shop)
	noFinding(t, "DROP INDEX IF EXISTS i ON t;\nCREATE INDEX i ON t (a);", shop)
	oneFinding(t, "DROP INDEX IF EXISTS i ON nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	// Another table's foreign key of the name stays.
	oneFinding(t, "ALTER TABLE nokey DROP FOREIGN KEY IF EXISTS fk1;\nALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);", 1, "Foreign key `fk1` already exists.", shop)
	noFinding(t, "ALTER TABLE t DROP FOREIGN KEY IF EXISTS fk1;\nALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);", shop)
}

func TestWalkThroughMariaDBCreateTable(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// A temporary table needs no database.
	noFinding(t, "CREATE DATABASE q;\nDROP DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nDROP DATABASE q;\nCREATE TABLE q.x (a int);", 2, "Database `q` doesn't exist.", shop)
	// IF NOT EXISTS of a table that exists reads no LIKE source; of a
	// temporary table it does.
	noFinding(t, "DROP TABLE r;\nCREATE TABLE IF NOT EXISTS t LIKE r;", shop)
	oneFinding(t, "DROP TABLE r;\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp LIKE r;", 2, "Table `r` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE r;\nCREATE TABLE x LIKE r;", 1, "Table `r` doesn't exist.", shop)
	// LIKE of a sequence makes a sequence, and of a view fails.
	noFinding(t, "CREATE TABLE x LIKE s;\nDROP SEQUENCE x;", shop)
	oneFinding(t, "CREATE TABLE x LIKE v;", 0, "`v` is a view, not a table.", shop)
}

func TestWalkThroughMariaDBSequences(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "CREATE SEQUENCE t;", 0, "The name `t` is already taken.", shop)
	oneFinding(t, "CREATE TABLE s (a int);", 0, "The name `s` is already taken.", shop)
	oneFinding(t, "CREATE VIEW s AS SELECT 1 AS a;", 0, "The name `s` is already taken.", shop)
	oneFinding(t, "CREATE OR REPLACE VIEW s AS SELECT 1 AS a;", 0, "`s` is a sequence, not a view.", shop)
	oneFinding(t, "CREATE SEQUENCE x;\nCREATE SEQUENCE x;", 1, "Sequence `x` already exists.", shop)
	noFinding(t, "CREATE SEQUENCE IF NOT EXISTS t;\nCREATE OR REPLACE SEQUENCE t;\nDROP SEQUENCE t;", shop)
	oneFinding(t, "CREATE OR REPLACE SEQUENCE v;", 0, "`v` is a view, not a table or sequence.", shop)
	oneFinding(t, "DROP SEQUENCE t;", 0, "`t` is a table, not a sequence.", shop)
	oneFinding(t, "DROP SEQUENCE v;", 0, "`v` is a view, not a sequence.", shop)
	oneFinding(t, "DROP SEQUENCE s;\nDROP SEQUENCE s;", 1, "Sequence `s` doesn't exist.", shop)
	oneFinding(t, "DROP VIEW s;", 0, "`s` is a sequence, not a view.", shop)
	noFinding(t, "DROP SEQUENCE IF EXISTS t, v, nosuch;", shop)
	noFinding(t, "DROP TABLE s;\nCREATE TABLE s (a int);", shop)
	noFinding(t, "RENAME TABLE s TO s2;\nDROP SEQUENCE s2;", shop)
	// Sequences: valid ALTER TABLE options keep the walk going; invalid
	// sequence operations report findings directly.
	oneFinding(t, "ALTER TABLE s COMMENT 'x';\nCREATE TABLE t (a int);", 1, "Table `t` already exists.", shop)
	oneFinding(t, "ALTER TABLE s ADD COLUMN z int;", 0, "Sequence `s` doesn't support this `ALTER TABLE` operation.", shop)
	oneFinding(t, "TRUNCATE TABLE s;\nCREATE TABLE t (a int);", 0, "`s` is a sequence, not a table.", shop)
	oneFinding(t, "CREATE INDEX k ON s (cache_size);\nCREATE TABLE t (a int);", 0, "`s` is a sequence, not a table.", shop)
	oneFinding(t, "DROP INDEX k ON s;", 0, "Index `k` doesn't exist in `s`.", shop)
	noFinding(t, "DROP INDEX IF EXISTS k ON s;", shop)
	oneFinding(t, "CREATE TRIGGER trs BEFORE INSERT ON s FOR EACH ROW SET @x = 1;\nCREATE TABLE t (a int);", 0, "`s` is a sequence, not a table.", shop)
	// ALTER SEQUENCE checks target existence, kind, and signed 64-bit increment bounds.
	noFinding(t, "ALTER SEQUENCE s NOMINVALUE RESTART 1;", shop)
	noFinding(t, "ALTER SEQUENCE IF EXISTS nosuch NOMINVALUE;", shop)
	noFinding(t, "ALTER SEQUENCE IF EXISTS t NOMINVALUE;", shop)
	oneFinding(t, "ALTER SEQUENCE nosuch NOMINVALUE;", 0, "Sequence `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER SEQUENCE t NOMINVALUE;", 0, "`t` is a table, not a sequence.", shop)
	oneFinding(t, "ALTER SEQUENCE v NOMINVALUE;", 0, "`v` is a view, not a sequence.", shop)
	oneFinding(t, "ALTER SEQUENCE s INCREMENT -9223372036854775808;", 0, "Sequence `s` has an out-of-range increment.", shop)
	// Temporary sequences: DROP SEQUENCE drops one before a base one,
	// passes over a temporary table, and DROP TEMPORARY SEQUENCE finds
	// only a temporary sequence.
	noFinding(t, "CREATE TEMPORARY SEQUENCE s;\nDROP SEQUENCE s;\nDROP SEQUENCE s;", shop)
	oneFinding(t, "CREATE TEMPORARY SEQUENCE s;\nDROP SEQUENCE s;\nDROP SEQUENCE s;\nDROP SEQUENCE s;", 3, "Sequence `s` doesn't exist.", shop)
	noFinding(t, "CREATE TEMPORARY TABLE s (a int);\nDROP SEQUENCE s;\nDROP TEMPORARY TABLE s;", shop)
	oneFinding(t, "DROP TEMPORARY SEQUENCE s;", 0, "Temporary sequence `s` doesn't exist.", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nDROP TEMPORARY SEQUENCE tt;", 1, "`tt` is a temporary table, not a sequence.", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY SEQUENCE tt;", 1, "The name `tt` is already taken.", shop)
	noFinding(t, "CREATE TEMPORARY SEQUENCE q.ts;", shop)
}

// TestWalkThroughMariaDBMySQLForms pins that the walk ends at the MySQL
// ALTER TABLE forms omni's MariaDB parser takes and MariaDB rejects.
func TestWalkThroughMariaDBMySQLForms(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range []string{
		"ALTER TABLE t ALTER CHECK c NOT ENFORCED;\nCREATE TABLE t (a int);",
		"ALTER TABLE t SECONDARY_LOAD;\nCREATE TABLE t (a int);",
	} {
		noFinding(t, sql, shop)
	}
	// MariaDB's ALTER INDEX ... [NOT] IGNORED checks index existence unless IF EXISTS.
	oneFinding(t, "ALTER TABLE t ALTER INDEX nosuch IGNORED;\nCREATE TABLE t (a int);", 0, "Index `nosuch` doesn't exist in `t`.", shop)
	noFinding(t, "ALTER TABLE t ALTER INDEX i IGNORED;\nALTER TABLE t ALTER INDEX i NOT IGNORED;", shop)
	noFinding(t, "ALTER TABLE t ALTER INDEX IF EXISTS nosuch IGNORED;", shop)
}

func TestWalkThroughMariaDBCase(t *testing.T) {
	sensitive := review.Target{Schema: walkSchema()}
	folded := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1}
	noFinding(t, "CREATE TABLE T (a int);", sensitive)
	oneFinding(t, "CREATE TABLE T (a int);", 0, "Table `T` already exists.", folded)
	// Trigger names compare with case: TR beside tr is a second trigger.
	noFinding(t, "CREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", sensitive)
	noFinding(t, "CREATE TRIGGER TR BEFORE UPDATE ON nokey FOR EACH ROW SET NEW.a = 1;\nDROP TRIGGER TR;\nDROP TRIGGER tr;", sensitive)
	oneFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER tr;", 1, "Trigger `tr` doesn't exist.", sensitive)
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
	oneFinding(t, "CREATE TABLE x (a int) COMMENT '中文';\nCREATE TABLE x (a int);", 1, "Table `x` already exists.", shop)
	oneFinding(t, "INSERT INTO nokey VALUES ('é');\nDROP TEMPORARY TABLE nosuch;", 1, "Temporary table `nosuch` doesn't exist.", shop)
	mixed := walkSchema()
	schema := mixed.Schemas[0]
	schema.Tables[2].Columns = append(schema.Tables[2].Columns, &metadata.ColumnMetadata{Name: "ſ"})
	schema.Procedures = append(schema.Procedures, &metadata.ProcedureMetadata{Name: "pé"})
	target := review.Target{Schema: mixed}
	noFinding(t, "ALTER TABLE nokey DROP COLUMN s;\nALTER TABLE nokey DROP COLUMN zz;\nDROP PROCEDURE pe;\nDROP PROCEDURE nosuch;", target)
	// A space that lists none still knows its absent names.
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column `zz` doesn't exist in `t`.", target)
	oneFinding(t, "DROP EVENT nosuch;", 0, "Event `nosuch` doesn't exist.", target)
}

// TestWalkThroughMariaDBSequenceEngineTables pins that the walk judges no
// name the SEQUENCE engine answers for (sequenceEngineTable): such a table
// is in every database though the sync lists none, so DROP TABLE and
// CREATE TRIGGER take it, and DML on it fails for another reason (1031).
// A name that only resembles one is still judged.
func TestWalkThroughMariaDBSequenceEngineTables(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range []string{
		"DROP TABLE seq_1_to_4;",
		"DROP TABLE seq_1_to_4;\nDROP TABLE seq_1_to_4;",
		"DROP TABLE seq_0_to_0;",
		"CREATE TRIGGER trq BEFORE INSERT ON seq_1_to_3 FOR EACH ROW SET @x = 1;",
		"INSERT INTO seq_1_to_3 VALUES (1);",
		"UPDATE seq_1_to_10_step_2 SET seq = 1;",
		"ALTER TABLE seq_1_to_100 ENGINE = BLACKHOLE;",
	} {
		noFinding(t, sql, shop)
	}
	oneFinding(t, "DROP TABLE seq_1_to_3x;", 0, "Table `seq_1_to_3x` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE seq_1_to_;", 0, "Table `seq_1_to_` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE SEQ_1_TO_3;", 0, "Table `SEQ_1_TO_3` doesn't exist.", shop)
	noFinding(t, "DROP TABLE SEQ_1_TO_3;", review.Target{Schema: walkSchema(), LowerCaseTableNames: 1})
	// A table of such a name the sync lists keeps its triggers.
	synced := walkSchema()
	synced.Schemas[0].Tables = append(synced.Schemas[0].Tables, &metadata.TableMetadata{Name: "seq_1_to_3", Triggers: []*metadata.TriggerMetadata{{Name: "trq"}}})
	noFinding(t, "DROP TRIGGER trq;", review.Target{Schema: synced})
	oneFinding(t, "DROP TRIGGER trq;\nDROP TRIGGER trq;", 1, "Trigger `trq` doesn't exist.", review.Target{Schema: synced})
}

func TestWalkThroughMariaDBSession(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Another database is unknown, until the change makes it.
	noFinding(t, "USE other;\nCREATE TABLE t (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 2, "Table `q.x` already exists.", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database is selected for `x`.", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TEMPORARY TABLE x (a int);", 1, "No database is selected for `x`.", shop)
	// A temporary table shadows the base table.
	noFinding(t, "DROP TABLE nokey;\nCREATE TEMPORARY TABLE nokey (a int);\nINSERT INTO nokey VALUES (1);", shop)
	oneFinding(t, "DROP TEMPORARY TABLE t;", 0, "Temporary table `t` doesn't exist.", shop)
	// A trigger may create a temporary table of any name.
	noFinding(t, "CREATE TEMPORARY TABLE x (a int);\nDROP TEMPORARY TABLE x;\nINSERT INTO t VALUES (1, 1);\nDROP TEMPORARY TABLE x;", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE x (a int);\nDROP TEMPORARY TABLE x;\nINSERT INTO nokey VALUES (1);\nDROP TEMPORARY TABLE x;", 3, "Temporary table `x` doesn't exist.", shop)
	// A foreign key's target is checked only while the checks are on.
	noFinding(t, "DROP TABLE r;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES r (id));", shop)
	oneFinding(t, "DROP TABLE nokey;\nSET foreign_key_checks = 1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nokey (a));", 2, "Referenced table `nokey` doesn't exist.", shop)
	oneFinding(t, "SET foreign_key_checks = 1;\nCREATE TABLE x (a bigint, FOREIGN KEY (a) REFERENCES s (next_not_cached_value));", 1, "`s` is a sequence, not a table a foreign key can reference.", shop)
	// Each assignment has a scope of its own: an unprefixed one the
	// statement's, past a prefixed one, and an @@session or bare @@ one the
	// session's. A global one leaves the session's checks as they are.
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	const fk = "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nokey (a));"
	const missing = "Referenced table `nokey` doesn't exist."
	oneFinding(t, "DROP TABLE nokey;\nSET GLOBAL foreign_key_checks = 0;\n"+fk, 2, missing, on)
	oneFinding(t, "DROP TABLE nokey;\nSET @@global.foreign_key_checks = 0;\n"+fk, 2, missing, on)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL foreign_key_checks = 1;\n"+fk, shop)
	oneFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, 2, missing, on)
	oneFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, @@session.foreign_key_checks = 1, foreign_key_checks = 0;\n"+fk, 2, missing, shop)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, @@session.foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, @@foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "DROP TABLE nokey;\nSET @@global.foreign_key_checks = 1, foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "DROP TABLE nokey;\nSET LOCAL foreign_key_checks = 0;\n"+fk, on)
	// SET STATEMENT sets foreign_key_checks for its statement alone, and puts it back after,
	// whatever the statement set it to. A variable it does not name stays as the statement set it.
	noFinding(t, "DROP TABLE nokey;\nSET STATEMENT foreign_key_checks = 0 FOR "+fk, on)
	oneFinding(t, "DROP TABLE nokey;\nSET STATEMENT foreign_key_checks = 0 FOR "+fk+"\n"+fk, 2, "Table `x` already exists.", on)
	oneFinding(t, "DROP TABLE nokey;\nSET STATEMENT foreign_key_checks = 0 FOR SET foreign_key_checks = 0;\n"+fk, 2, missing, on)
	noFinding(t, "DROP TABLE nokey;\nSET STATEMENT max_statement_time = 1 FOR SET foreign_key_checks = 0;\n"+fk, on)
	oneFinding(t, "DROP TABLE nokey;\nSET STATEMENT max_statement_time = 1 FOR "+fk, 1, missing, on)
	oneFinding(t, "DROP TABLE nokey;\nSET STATEMENT foreign_key_checks = 0, foreign_key_checks = 1 FOR "+fk, 1, missing, on)
	noFinding(t, "DROP TABLE nokey;\nSET STATEMENT sql_mode = '' FOR SET STATEMENT foreign_key_checks = 0 FOR "+fk, on)
	// SET STATEMENT refuses a stored function or a table in an option, and what it wraps does not
	// happen, so the walk follows none of it. A built-in function is no such thing.
	noFinding(t, "CREATE FUNCTION mk() RETURNS int DETERMINISTIC RETURN 1;\nSET STATEMENT max_statement_time = mk() FOR DROP TABLE nokey;\nDROP TABLE nokey;", shop)
	noFinding(t, "SET STATEMENT max_statement_time = (SELECT COUNT(*) FROM t) FOR DROP TABLE nokey;\nDROP TABLE nokey;", shop)
	oneFinding(t, "SET STATEMENT max_statement_time = ABS(1) FOR DROP TABLE nokey;\nDROP TABLE nokey;", 1, "Table `nokey` doesn't exist.", shop)
	oneFinding(t, "SET STATEMENT max_statement_time = (SELECT 1) FOR DROP TABLE nokey;\nDROP TABLE nokey;", 1, "Table `nokey` doesn't exist.", shop)
	// A later assignment writes its own scope, and the most recent scope reaches one that
	// writes none.
	oneFinding(t, "DROP TABLE nokey;\nSET @a = 1, GLOBAL foreign_key_checks = 0;\n"+fk, 2, missing, on)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, SESSION foreign_key_checks = 0;\n"+fk, on)
	oneFinding(t, "DROP TABLE nokey;\nSET SESSION sql_mode = DEFAULT, GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, 2, missing, on)
	noFinding(t, "DROP TABLE nokey;\nSET GLOBAL sql_mode = DEFAULT, SESSION sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, on)
	// With another spec of the plan, only the change's own facts count.
	other := review.Target{Schema: walkSchema(), OtherSpecs: true}
	noFinding(t, "CREATE TABLE t (a int);", other)
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE x (a int);", 1, "Table `x` already exists.", other)
}

// TestWalkThroughMariaDBForeignKeyNames pins the name a foreign key and
// its index take: the CONSTRAINT symbol, or without one the name after
// FOREIGN KEY, which counts as long only then (verified on 11.8.9).
func TestWalkThroughMariaDBForeignKeyNames(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	long := strings.Repeat("f", 65)
	noFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY cx;\nALTER TABLE x DROP INDEX cx;", shop)
	oneFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY fx;", 1, "Foreign key `fx` doesn't exist in `x`.", shop)
	noFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nALTER TABLE y DROP FOREIGN KEY fy;", shop)
	oneFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nCREATE TABLE z (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));", 1, "Foreign key `fy` already exists.", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY fk1 (a) REFERENCES r (id);", 0, "Foreign key `fk1` already exists.", shop)
	noFinding(t, "CREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id));", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id);", shop)
	oneFinding(t, "CREATE TABLE w (a int, FOREIGN KEY "+long+" (a) REFERENCES r (id));", 0, "The name `"+long+"` is longer than MariaDB's 64-character limit.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY "+long+" (a) REFERENCES r (id);", 0, "The name `"+long+"` is longer than MariaDB's 64-character limit.", shop)
}

// TestWalkThroughMariaDBNoDatabase pins what fails once the change drops
// the current database: an unqualified name (1046), but for DROP
// FUNCTION, which then looks for a loadable function alone (verified on
// 11.8.9).
func TestWalkThroughMariaDBNoDatabase(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS nosuch;\nDROP FUNCTION nosuch;\nDROP FUNCTION IF EXISTS shop.f;", shop)
	for stmt, name := range map[string]string{"DROP PROCEDURE IF EXISTS p;": "p", "DROP TRIGGER IF EXISTS tr;": "tr", "DROP VIEW IF EXISTS v;": "v", "DROP EVENT IF EXISTS e;": "e", "DROP SEQUENCE IF EXISTS s;": "s", "CREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;": "f2", "ALTER FUNCTION f COMMENT 'x';": "f"} {
		oneFinding(t, "DROP DATABASE shop;\n"+stmt, 1, "No database is selected for `"+name+"`.", shop)
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
	oneFinding(t, gone+"ALTER TABLE nokey ADD FOREIGN KEY (a) REFERENCES p2 (id);", 2, "Referenced table `p2` doesn't exist.", on)
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
	oneFinding(t, "SELECT * FROM nokey;\n"+drop, 1, "Temporary table `made` doesn't exist.", shop)
	oneFinding(t, "SELECT NEXTVAL(s) FROM s;\n"+drop, 1, "Temporary table `made` doesn't exist.", shop)
	oneFinding(t, "SELECT mk();\n"+drop, 1, "Temporary table `made` doesn't exist.", shop)
}

// TestWalkThroughMariaDBSequenceRenamed pins that a relation keeps its
// kind through ALTER TABLE ... RENAME: one made LIKE another database's
// sequence may be a sequence, and stays one (verified on 11.8.9).
func TestWalkThroughMariaDBSequenceRenamed(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "CREATE TABLE t2 LIKE q.s;\nALTER TABLE t2 RENAME TO t3;\nDROP SEQUENCE t3;", shop)
	noFinding(t, "CREATE TABLE t2 LIKE q.s;\nALTER TABLE t2 RENAME TO q.t3, COMMENT 'x';\nDROP SEQUENCE q.t3;", shop)
	oneFinding(t, "CREATE TABLE t2 (a int);\nALTER TABLE t2 RENAME TO t3;\nDROP SEQUENCE t3;", 2, "`t3` is a table, not a sequence.", shop)
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

// TestWalkThroughMariaDBAlterIfExistsAndPartitionMaintenance pins that
// ALTER TABLE IF EXISTS on a missing table and partition maintenance actions
// (OPTIMIZE, ANALYZE, CHECK, REPAIR PARTITION) produce no SQL error on
// mariadb:11.8.8, so WALK_THROUGH stays silent on them.
func TestWalkThroughMariaDBAlterIfExistsAndPartitionMaintenance(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, stmt := range []string{
		"ALTER TABLE IF EXISTS nosuch;",
		"ALTER TABLE IF EXISTS nosuch ADD COLUMN x int;",
		"ALTER TABLE nosuch OPTIMIZE PARTITION ALL;",
		"ALTER TABLE nosuch ANALYZE PARTITION ALL;",
		"ALTER TABLE nosuch CHECK PARTITION ALL;",
		"ALTER TABLE nosuch REPAIR PARTITION ALL;",
		"ALTER TABLE v CHECK PARTITION ALL;",
	} {
		noFinding(t, stmt, shop)
	}
	// ALTER TABLE IF EXISTS on an existing table still checks its subcommands.
	oneFinding(t, "ALTER TABLE IF EXISTS t ADD COLUMN a int;", 0, "Column `a` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE nosuch REBUILD PARTITION ALL;", 0, "Table `nosuch` doesn't exist.", shop)
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
	{SQL: "ALTER SEQUENCE s NOMINVALUE;"},
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

func TestWalkThroughMariaDBFKGeneratedIndex(t *testing.T) {
	schema := &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name:    "parent",
					Columns: []*metadata.ColumnMetadata{{Name: "b"}},
					Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Expressions: []string{"b"}, Primary: true, Unique: true}},
				},
				{
					Name:    "a",
					Columns: []*metadata.ColumnMetadata{{Name: "b"}, {Name: "a"}, {Name: "d"}},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Expressions: []string{"b"}, Primary: true, Unique: true},
						{Name: "idx_d", Expressions: []string{"d"}},
						{Name: "a", Expressions: []string{"a"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{
						{Name: "a", Columns: []string{"a"}, ReferencedTable: "parent"},
					},
				},
			},
		}},
	}
	target := review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn}

	// MariaDB silently drops an InnoDB auto-generated FK backing index when an
	// explicit BTREE/default index covers its columns, even under the same name.
	noFinding(t, "CREATE INDEX a ON a (a);", target)
	noFinding(t, "CREATE UNIQUE INDEX a ON a (a);", target)
	noFinding(t, "ALTER TABLE a ADD KEY a (a);", target)
	noFinding(t, "ALTER TABLE a ADD KEY a (a, d);", target)
	noFinding(t, "ALTER TABLE a ADD UNIQUE a TYPE BTREE (a);", target)

	// Once replaced by an explicit index, a second index of the same name fails (1061).
	oneFinding(t, "CREATE INDEX a ON a (a);\nCREATE INDEX a ON a (a);", 1, "Index `a` already exists in `a`.", target)

	// Non-covering, HASH, or unrelated existing indexes still report 1061.
	oneFinding(t, "ALTER TABLE a ADD KEY a (d);", 0, "Index `a` already exists in `a`.", target)
	oneFinding(t, "CREATE INDEX a USING HASH ON a (a);", 0, "Index `a` already exists in `a`.", target)
	oneFinding(t, "CREATE INDEX idx_d ON a (d);", 0, "Index `idx_d` already exists in `a`.", target)

	// Tracking survives column rename and table rename.
	noFinding(t, "ALTER TABLE a RENAME COLUMN a TO x;\nCREATE INDEX a ON a (x);", target)
	noFinding(t, "RENAME TABLE a TO c;\nCREATE INDEX a ON c (a);", target)

	// Adding a different covering index may drop the implicit FK index `a`, so
	// a follow-up ADD KEY a (d) is no longer certain to collide.
	noFinding(t, "ALTER TABLE a ADD KEY zzz (a);\nALTER TABLE a ADD KEY a (d);", target)
}

// TestWalkThroughMariaDBRecallFamilies covers every MariaDB 11.8.8 refusal
// family added for walk recall along with valid accepted controls.
func TestWalkThroughMariaDBRecallFamilies(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Tables = append(schema.Schemas[0].Tables, &metadata.TableMetadata{
		Name:       "pt",
		Columns:    []*metadata.ColumnMetadata{{Name: "id"}, {Name: "a"}},
		Indexes:    []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}},
		Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}, {Name: "p1"}},
	})
	shop := review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn}

	// 1. ALTER SEQUENCE on missing, non-sequence, and valid targets.
	oneFinding(t, "ALTER SEQUENCE nosuch RESTART 1;", 0, "Sequence `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER SEQUENCE t RESTART 1;", 0, "`t` is a table, not a sequence.", shop)
	oneFinding(t, "ALTER SEQUENCE v RESTART 1;", 0, "`v` is a view, not a sequence.", shop)
	oneFinding(t, "ALTER SEQUENCE s INCREMENT -9223372036854775808;", 0, "Sequence `s` has an out-of-range increment.", shop)
	noFinding(t, "ALTER SEQUENCE s INCREMENT 2 START 10 CACHE 100 CYCLE;", shop)
	noFinding(t, "CREATE TEMPORARY SEQUENCE ts;\nALTER SEQUENCE ts RESTART 5;", shop)

	// 2. Sequence operations: ALTER TABLE, CREATE/DROP INDEX, TRUNCATE, CREATE TRIGGER.
	oneFinding(t, "ALTER TABLE s ADD COLUMN x INT;", 0, "Sequence `s` doesn't support this `ALTER TABLE` operation.", shop)
	oneFinding(t, "ALTER TABLE s ADD KEY (cache_size);", 0, "Sequence `s` doesn't support this `ALTER TABLE` operation.", shop)
	oneFinding(t, "ALTER TABLE s DROP COLUMN cache_size;", 0, "Sequence `s` doesn't support this `ALTER TABLE` operation.", shop)
	oneFinding(t, "ALTER TABLE s DROP COLUMN nosuch;", 0, "Column `nosuch` doesn't exist in `s`.", shop)
	oneFinding(t, "ALTER TABLE s MODIFY nosuch INT;", 0, "Column `nosuch` doesn't exist in `s`.", shop)
	oneFinding(t, "ALTER TABLE s DROP PRIMARY KEY;", 0, "`s` has no primary key.", shop)
	oneFinding(t, "ALTER TABLE s DROP INDEX nosuch;", 0, "Index `nosuch` doesn't exist in `s`.", shop)
	oneFinding(t, "ALTER TABLE s ALTER KEY nosuch IGNORED;", 0, "Index `nosuch` doesn't exist in `s`.", shop)
	oneFinding(t, "ALTER TABLE s REMOVE PARTITIONING;", 0, "Table `s` is not partitioned.", shop)
	oneFinding(t, "ALTER TABLE s DROP SYSTEM VERSIONING;", 0, "Table `s` is not system-versioned.", shop)
	noFinding(t, "ALTER TABLE s COMMENT 'ok', FORCE;\nALTER TABLE s DROP COLUMN IF EXISTS nosuch;\nALTER TABLE s RENAME TO s2;\nDROP SEQUENCE s2;", shop)

	// 3. Partition management on non-partitioned vs partitioned tables.
	oneFinding(t, "ALTER TABLE t REMOVE PARTITIONING;", 0, "Table `t` is not partitioned.", shop)
	oneFinding(t, "ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (200));", 0, "Table `t` is not partitioned.", shop)
	oneFinding(t, "ALTER TABLE t ANALYZE PARTITION ALL;", 0, "Table `t` is not partitioned.", shop)
	noFinding(t, "ALTER TABLE pt DROP PARTITION p0;", shop)
	noFinding(t, "ALTER TABLE pt ANALYZE PARTITION ALL;", shop)
	noFinding(t, "ALTER TABLE pt REMOVE PARTITIONING;", shop)
	oneFinding(t, "ALTER TABLE pt REMOVE PARTITIONING;\nALTER TABLE pt DROP PARTITION p1;", 1, "Table `pt` is not partitioned.", shop)
	noFinding(t, "CREATE TABLE np (id INT) PARTITION BY HASH (id) PARTITIONS 2;\nALTER TABLE np COALESCE PARTITION 1;", shop)
	oneFinding(t, "CREATE TABLE plain_new (id INT);\nALTER TABLE plain_new REMOVE PARTITIONING;", 1, "Table `plain_new` is not partitioned.", shop)
	noFinding(t, "ALTER TABLE t PARTITION BY HASH (id) PARTITIONS 2;\nALTER TABLE t COALESCE PARTITION 1;", shop)

	// 4. Index/key column existence and duplicate key columns.
	oneFinding(t, "CREATE INDEX k ON t (nosuch);", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	oneFinding(t, "CREATE INDEX IF NOT EXISTS k ON t (nosuch);", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	noFinding(t, "CREATE INDEX IF NOT EXISTS i ON t (nosuch);", shop)
	oneFinding(t, "CREATE INDEX k ON t (a, a);", 0, "This statement names column `a` twice.", shop)
	oneFinding(t, "ALTER TABLE t ADD KEY k (nosuch);", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD KEY k (a, a);", 0, "This statement names column `a` twice.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY (nosuch) REFERENCES r (id);", 0, "Column `nosuch` doesn't exist in `nokey`.", shop)
	oneFinding(t, "CREATE TABLE bad_idx (a INT, KEY k (nosuch));", 0, "Column `nosuch` doesn't exist in `bad_idx`.", shop)
	oneFinding(t, "CREATE TABLE bad_dup (a INT, KEY k (a, a));", 0, "This statement names column `a` twice.", shop)
	// Valid controls: columns added/renamed in the same ALTER, functional key parts, WITHOUT OVERLAPS.
	noFinding(t, "ALTER TABLE t ADD COLUMN q INT, ADD INDEX kq (q);", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX kq (q), ADD COLUMN q INT;", shop)
	noFinding(t, "ALTER TABLE t ADD (q INT, KEY kq (q));", shop)
	noFinding(t, "ALTER TABLE t RENAME COLUMN a TO a2, ADD INDEX ka2 (a2);", shop)
	noFinding(t, "CREATE TABLE app_p (id INT, s DATE, e DATE, PERIOD FOR p (s, e), UNIQUE (id, p WITHOUT OVERLAPS));", shop)

	// 5. CREATE TABLE columnless (1113) and PERIOD FOR SYSTEM_TIME without versioning (4125).
	oneFinding(t, "CREATE TABLE empty_t;", 0, "Table `empty_t` must have at least one column.", shop)
	oneFinding(t, "CREATE OR REPLACE TABLE t;", 0, "Table `t` must have at least one column.", shop)
	oneFinding(t, "CREATE TABLE t;", 0, "Table `t` already exists.", shop)
	noFinding(t, "CREATE TABLE IF NOT EXISTS t;", shop)
	oneFinding(t, "CREATE TABLE IF NOT EXISTS empty_t;", 0, "Table `empty_t` must have at least one column.", shop)
	oneFinding(t, "CREATE TABLE bad_per (PERIOD FOR SYSTEM_TIME (a, b));", 0, "Table `bad_per` requires `WITH SYSTEM VERSIONING` for `PERIOD FOR SYSTEM_TIME`.", shop)
	noFinding(t, "CREATE TABLE ok_per (id INT, s TIMESTAMP(6) GENERATED ALWAYS AS ROW START, e TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING;", shop)
	// Column-level WITHOUT SYSTEM VERSIONING normalizes the period away into an ordinary table.
	noFinding(t, "CREATE TABLE ok_norm (x INT WITHOUT SYSTEM VERSIONING, rs INT GENERATED ALWAYS AS ROW START, re INT GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME (rs, re));", shop)

	// 6. Bare column without FROM in views/CTAS and event schedules (1054), and event self-rename (1551).
	oneFinding(t, "CREATE VIEW v_bad AS SELECT nosuch;", 0, "Column `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER VIEW v AS SELECT nosuch;", 0, "Column `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE TABLE ct_bad AS SELECT nosuch;", 0, "Column `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE EVENT ev_bad ON SCHEDULE AT nosuch DO COMMIT;", 0, "Column `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER EVENT e ON SCHEDULE EVERY nosuch DAY;", 0, "Column `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER EVENT e RENAME TO e;", 0, "Event `e` cannot be renamed to itself.", shop)
	noFinding(t, "CREATE VIEW v_ok AS SELECT 1 AS a, CURRENT_TIMESTAMP AS ts, @x AS u, @@version AS ver;", shop)
	noFinding(t, "CREATE VIEW v_ok2 AS SELECT a FROM t;", shop)
	noFinding(t, "CREATE VIEW v_ok3 AS SELECT (SELECT a FROM t LIMIT 1) AS x;", shop)
	noFinding(t, "CREATE TABLE ct_ok AS SELECT 1 AS a;", shop)
	noFinding(t, "CREATE EVENT ev_ok ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 HOUR DO COMMIT;", shop)

	// 7. System-versioning ALTER TABLE checks on non-versioned tables (4124) and unknown LIKE control.
	oneFinding(t, "ALTER TABLE t DROP SYSTEM VERSIONING;", 0, "Table `t` is not system-versioned.", shop)
	oneFinding(t, "ALTER TABLE t DROP PERIOD FOR SYSTEM_TIME;", 0, "Table `t` is not system-versioned.", shop)
	oneFinding(t, "ALTER TABLE t PARTITION BY SYSTEM_TIME;", 0, "Table `t` is not system-versioned.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN s TIMESTAMP(6) AS ROW START;", 0, "Table `t` is not system-versioned.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN x INT WITH SYSTEM VERSIONING;", 0, "Table `t` is not system-versioned.", shop)
	noFinding(t, "ALTER TABLE t ADD SYSTEM VERSIONING PARTITION BY SYSTEM_TIME;", shop)
	noFinding(t, "CREATE TABLE from_other LIKE other_db.vt;\nALTER TABLE from_other DROP SYSTEM VERSIONING;", shop)

	// 8. 3-part field DB mismatch (1102), RENAME self (1050/1146), ALTER TABLE IF EXISTS view (1347).
	oneFinding(t, "ALTER TABLE t DROP other.t.a;", 0, "Column qualifier `other` doesn't match database `shop`.", shop)
	noFinding(t, "ALTER TABLE t DROP shop.t.a;", shop)
	oneFinding(t, "RENAME TABLE t TO t, t TO t;", 0, "Table `t` already exists.", shop)
	oneFinding(t, "RENAME TABLE nosuch TO nosuch, nosuch TO nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE IF EXISTS v;", 0, "`v` is a view, not a table.", shop)

	// 9. Self-referencing column expressions (4029), CHECK/DEFAULT missing columns (1054), FK col count (1239).
	oneFinding(t, "ALTER TABLE t ALTER a SET DEFAULT a;", 0, "Expression for column `a` cannot reference itself.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN z INT DEFAULT z;", 0, "Expression for column `z` cannot reference itself.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN z INT AS (z + 1);", 0, "Expression for column `z` cannot reference itself.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN z INT DEFAULT nosuch;", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD CHECK (nosuch > 0);", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN z INT REFERENCES r (id, id);", 0, "Foreign key on `t` has a mismatched column count.", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN z INT DEFAULT (a + 1);", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN z INT DEFAULT (NEXTVAL(s));", shop)
	noFinding(t, "ALTER TABLE t ADD (z1 INT, z2 INT DEFAULT (z1 + 1));", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN z INT, ADD CHECK (z > 0);", shop)
	// Earlier additions in the same ALTER satisfy later IF NOT EXISTS guards, skipping validation.
	noFinding(t, "ALTER TABLE t ADD COLUMN c INT, ADD COLUMN IF NOT EXISTS c INT DEFAULT nosuch;", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN c INT, ADD COLUMN IF NOT EXISTS c INT DEFAULT c;", shop)
	noFinding(t, "ALTER TABLE t RENAME COLUMN a TO c, ADD COLUMN IF NOT EXISTS c INT DEFAULT nosuch;", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX k (a), ADD INDEX IF NOT EXISTS k (nosuch);", shop)
	noFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS k (a), ADD INDEX IF NOT EXISTS k (nosuch);", shop)
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

func TestWalkThroughMariaDBTypesAndQueries(t *testing.T) {
	schema := walkSchema()
	typed := func(name, typ string) *metadata.ColumnMetadata {
		return &metadata.ColumnMetadata{Name: name, Type: typ}
	}
	pk := &metadata.IndexMetadata{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}}
	fk := []*metadata.ForeignKeyMetadata{{Name: "fk_p", Columns: []string{"p"}, ReferencedTable: "r", ReferencedColumns: []string{"id"}}}
	schema.Schemas[0].Tables = append(schema.Schemas[0].Tables,
		&metadata.TableMetadata{
			Name:    "ty",
			Columns: []*metadata.ColumnMetadata{typed("id", "int(11)"), typed("n", "int(11)"), typed("s", "varchar(10)"), typed("g", "point")},
			Indexes: []*metadata.IndexMetadata{pk},
		},
		&metadata.TableMetadata{Name: "my", Engine: "MyISAM", Columns: []*metadata.ColumnMetadata{typed("n", "int(11)")}},
		&metadata.TableMetadata{
			Name:       "rp",
			Columns:    []*metadata.ColumnMetadata{typed("id", "int(11)")},
			Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_RANGE, Expression: "id"}},
		},
		&metadata.TableMetadata{
			Name:       "hp",
			Columns:    []*metadata.ColumnMetadata{typed("id", "int(11)")},
			Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Type: metadata.TablePartitionMetadata_HASH, Expression: "id"}, {Name: "p1", Type: metadata.TablePartitionMetadata_HASH, Expression: "id"}},
		},
		&metadata.TableMetadata{
			Name:        "ch",
			Columns:     []*metadata.ColumnMetadata{typed("id", "int(11)"), typed("p", "int(11)")},
			Indexes:     []*metadata.IndexMetadata{pk, {Name: "fk_p", Expressions: []string{"p"}}},
			ForeignKeys: fk,
		},
		&metadata.TableMetadata{
			Name:        "ch2",
			Columns:     []*metadata.ColumnMetadata{typed("id", "int(11)"), typed("p", "int(11)")},
			Indexes:     []*metadata.IndexMetadata{pk, {Name: "fk_p", Expressions: []string{"p"}}, {Name: "k2", Expressions: []string{"p", "id"}}},
			ForeignKeys: fk,
		},
	)
	schema.Schemas[0].Sequences = append(schema.Schemas[0].Sequences, &metadata.SequenceMetadata{Name: "sq", Start: "1", MinValue: "1", MaxValue: "100", Increment: "1"})
	shop := review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn}

	// Key parts judged by their columns' types (1283, 1210, 1089).
	oneFinding(t, "ALTER TABLE ty ADD FULLTEXT (n);", 0, "Column `n` cannot be part of a `FULLTEXT` index.", shop)
	oneFinding(t, "CREATE SPATIAL INDEX k ON ty (n);", 0, "Column `n` cannot be part of a `SPATIAL` index.", shop)
	oneFinding(t, "ALTER TABLE ty ADD KEY USING RTREE (s);", 0, "Column `s` cannot be part of a `SPATIAL` index.", shop)
	oneFinding(t, "ALTER TABLE ty ADD VECTOR (n);", 0, "Column `n` cannot be part of a `VECTOR` index.", shop)
	oneFinding(t, "ALTER TABLE ty ADD KEY (n(3));", 0, "Column `n` cannot use an index prefix length.", shop)
	oneFinding(t, "ALTER TABLE ty ADD COLUMN m DATE;\nCREATE FULLTEXT INDEX k ON ty (m);", 1, "Column `m` cannot be part of a `FULLTEXT` index.", shop)
	oneFinding(t, "ALTER TABLE ty ADD COLUMN m DATE, ADD FULLTEXT (m);", 0, "Column `m` cannot be part of a `FULLTEXT` index.", shop)
	oneFinding(t, "CREATE TABLE nt (a INT, FULLTEXT (a));", 0, "Column `a` cannot be part of a `FULLTEXT` index.", shop)
	noFinding(t, "ALTER TABLE ty ADD FULLTEXT (s);\nCREATE SPATIAL INDEX k ON ty (g);\nALTER TABLE ty ADD KEY (s(3));", shop)
	noFinding(t, "ALTER TABLE ty MODIFY n TEXT, ADD FULLTEXT (n);", shop)
	noFinding(t, "ALTER TABLE ty MODIFY n TEXT;\nALTER TABLE ty ADD FULLTEXT (n);", shop)
	noFinding(t, "ALTER TABLE ty CHANGE n n TEXT;\nCREATE FULLTEXT INDEX k ON ty (n);", shop)
	noFinding(t, "ALTER TABLE ty RENAME COLUMN s TO n2;\nALTER TABLE ty ADD SPATIAL (n2);", shop)
	noFinding(t, "ALTER TABLE ty ADD COLUMN IF NOT EXISTS n TEXT, ADD FULLTEXT (n);", shop)
	noFinding(t, "ALTER TABLE my ADD FULLTEXT (n);", shop)
	noFinding(t, "CREATE TABLE nt (a INT, FULLTEXT (a)) ENGINE=Mroonga;", shop)

	// AUTO_INCREMENT needs a key (1075).
	oneFinding(t, "ALTER TABLE ty ADD COLUMN x INT AUTO_INCREMENT;", 0, "Column `x` is `AUTO_INCREMENT` and must be defined as a key.", shop)
	oneFinding(t, "CREATE TABLE ai (a INT AUTO_INCREMENT, b INT, KEY (b));", 0, "Column `a` is `AUTO_INCREMENT` and must be defined as a key.", shop)
	noFinding(t, "ALTER TABLE ty ADD COLUMN x INT AUTO_INCREMENT, ADD KEY (x);", shop)
	noFinding(t, "ALTER TABLE ty ADD COLUMN x INT AUTO_INCREMENT UNIQUE;", shop)
	noFinding(t, "CREATE TABLE ai (a INT AUTO_INCREMENT PRIMARY KEY);\nCREATE TABLE ai2 (a INT AUTO_INCREMENT, KEY (a));", shop)

	// SEQUENCE=1 on a table of other columns (4086).
	oneFinding(t, "ALTER TABLE ty SEQUENCE=1;", 0, "Table `ty` has column `g`, which a sequence has not.", shop)
	noFinding(t, "ALTER TABLE ty SEQUENCE=0;", shop)
	noFinding(t, "ALTER TABLE ty DROP COLUMN g;\nALTER TABLE ty DROP COLUMN s;\nALTER TABLE ty DROP COLUMN n;\nALTER TABLE ty DROP COLUMN id, ADD next_not_cached_value BIGINT;\nALTER TABLE ty SEQUENCE=1;", shop)

	// ALTER SEQUENCE values against the synced ones (4085).
	oneFinding(t, "ALTER SEQUENCE sq MAXVALUE 1;", 0, "Sequence `sq` has invalid values.", shop)
	oneFinding(t, "ALTER SEQUENCE sq START -1;", 0, "Sequence `sq` has invalid values.", shop)
	oneFinding(t, "ALTER SEQUENCE sq MAXVALUE 50;\nALTER SEQUENCE sq START 60;", 1, "Sequence `sq` has invalid values.", shop)
	noFinding(t, "ALTER SEQUENCE sq MAXVALUE 50 START 40;\nALTER SEQUENCE sq NO MAXVALUE;\nALTER SEQUENCE s MAXVALUE 1;", shop)
	noFinding(t, "DROP SEQUENCE sq;\nCREATE SEQUENCE sq MINVALUE -10 START -5;\nALTER SEQUENCE sq START -1;", shop)

	// Partition maintenance of a RANGE table (1479, 1492, 1509, 1511).
	oneFinding(t, "ALTER TABLE rp COALESCE PARTITION 1;", 0, "Table `rp` is partitioned by `RANGE` or `LIST`, which `COALESCE PARTITION` does not take.", shop)
	oneFinding(t, "ALTER TABLE rp ADD PARTITION (PARTITION p9);", 0, "Partition `p9` of `rp` requires a `VALUES` clause.", shop)
	oneFinding(t, "ALTER TABLE rp ADD PARTITION PARTITIONS 2;", 0, "Table `rp` is partitioned by `RANGE` or `LIST`, and each partition it adds must be defined.", shop)
	oneFinding(t, "ALTER TABLE rp REORGANIZE PARTITION;", 0, "Table `rp` is partitioned by `RANGE` or `LIST`, which `REORGANIZE PARTITION` without partitions does not take.", shop)
	noFinding(t, "ALTER TABLE rp ADD PARTITION (PARTITION p9 VALUES LESS THAN (1000));", shop)
	noFinding(t, "ALTER TABLE rp PARTITION BY HASH (id) PARTITIONS 4;\nALTER TABLE rp COALESCE PARTITION 1;", shop)
	noFinding(t, "ALTER TABLE hp COALESCE PARTITION 1;", shop)

	// A foreign key's only index, or its column (1553).
	oneFinding(t, "ALTER TABLE ch DROP INDEX fk_p;", 0, "Index `fk_p` is needed by a foreign key constraint.", shop)
	oneFinding(t, "DROP INDEX fk_p ON ch;", 0, "Index `fk_p` is needed by a foreign key constraint.", shop)
	oneFinding(t, "ALTER TABLE ch DROP COLUMN p;", 0, "Column `p` is needed by a foreign key constraint.", shop)
	noFinding(t, "ALTER TABLE ch DROP FOREIGN KEY fk_p, DROP INDEX fk_p;", shop)
	noFinding(t, "ALTER TABLE ch2 DROP INDEX fk_p;", shop)
	noFinding(t, "ALTER TABLE ch ADD KEY k2 (p, id);\nALTER TABLE ch DROP INDEX fk_p;", shop)
	noFinding(t, "SET foreign_key_checks = 0;\nALTER TABLE ch DROP INDEX fk_p;", shop)
	noFinding(t, "ALTER TABLE ch DROP INDEX fk_p;", review.Target{Schema: schema})

	// A view's query (1054, 1060, 1096, 1146, 1210, 1241, 4009, 4017, 4089).
	oneFinding(t, "CREATE VIEW vw AS SELECT ROW_NUMBER() OVER w FROM t;", 0, "Window `w` is not defined.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT RANK() OVER (ROWS UNBOUNDED PRECEDING) FROM t;", 0, "Window function `RANK` takes no frame.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT NEXTVAL(t) FROM t;", 0, "`t` is a table, not a sequence.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT NEXT VALUE FOR nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT a FROM t WHERE a LIKE 'x' ESCAPE a;", 0, "The `ESCAPE` of `LIKE` must be a constant.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT MATCH (a) AGAINST (a) FROM t;", 0, "The `AGAINST` of `MATCH` must be a constant.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT nosuch FROM t;", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT a FROM t WHERE nosuch > 0;", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT a, a FROM t;", 0, "This statement names column `a` twice.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT a FROM nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT *;", 0, "This query reads no table for `*`.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT t.*;", 0, "Table `t` is not in this query.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT a FROM DUAL;", 0, "Column `a` doesn't exist.", shop)
	oneFinding(t, "CREATE VIEW vw AS SELECT ROW(1, 2);", 0, "A view's column cannot be a row.", shop)
	oneFinding(t, "ALTER VIEW v AS SELECT nosuch FROM t;", 0, "Column `nosuch` doesn't exist in `t`.", shop)
	noFinding(t, "CREATE VIEW vw AS SELECT ROW_NUMBER() OVER w, SUM(id) OVER (ROWS UNBOUNDED PRECEDING) FROM t WINDOW w AS (ORDER BY id);", shop)
	noFinding(t, "CREATE VIEW vw AS SELECT NEXTVAL(s), LASTVAL(s), NEXT VALUE FOR s;", shop)
	noFinding(t, "CREATE VIEW vw AS SELECT a FROM t WHERE a LIKE 'x' ESCAPE '|';", shop)
	noFinding(t, "CREATE VIEW vw AS SELECT x.a, a AS z, (SELECT nosuch FROM v) AS q FROM t AS x WHERE x.id > 0 ORDER BY z;", shop)
	noFinding(t, "CREATE VIEW vw (c1, c2) AS SELECT a, a FROM t;", shop)
	noFinding(t, "CREATE VIEW vw AS SELECT a, a + 1, a + 1 FROM t;", shop)
	noFinding(t, "CREATE VIEW vw AS WITH c AS (SELECT 1 AS x) SELECT x FROM c;", shop)
	noFinding(t, "CREATE VIEW vw AS SELECT nosuch FROM v;", shop)
	noFinding(t, "CREATE TABLE nt (q INT);\nCREATE VIEW vw AS SELECT q FROM nt;", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN q INT;\nCREATE VIEW vw AS SELECT q FROM t;", shop)
}
