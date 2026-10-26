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

// walkSchema is the synced shop database the unit tests walk on.
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
					Columns:     cols("id", "a", "b"),
					Indexes:     []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}, {Name: "i"}},
					ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk1", ReferencedTable: "r"}},
				},
				{Name: "r", Columns: cols("id"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}}},
				{Name: "nokey", Columns: cols("a")},
			},
			Views:     []*metadata.ViewMetadata{{Name: "v"}},
			Sequences: []*metadata.SequenceMetadata{{Name: "sq"}},
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

func TestWalkThroughTiDB(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table t already exists", shop)
	oneFinding(t, "ALTER TABLE nosuch ADD COLUMN c int;", 0, "Table nosuch does not exist", shop)
	oneFinding(t, "ALTER TABLE v ADD COLUMN c int;", 0, "v is not a base table", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN A int;", 0, "Column A already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column zz does not exist in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX I (a);", 0, "Index I already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "Table t already has a primary key", shop)
	oneFinding(t, "DROP TABLE t, nosuch;", 0, "Table nosuch does not exist", shop)
	oneFinding(t, "DROP VIEW t;", 0, "t is not a view", shop)
	noFinding(t, "RENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN zz int;\nALTER TABLE t ADD COLUMN zz int;", shop)
	oneFinding(t, "RENAME TABLE t TO u;\nCREATE TABLE u (a int);", 1, "Table u already exists", shop)
	// A quoted `db.t` is one name, of the current database.
	oneFinding(t, "ALTER TABLE t RENAME TO `q.t`;\nCREATE TABLE `q.t` (a int);", 1, "Table `q.t` already exists", shop)
}

// TestWalkThroughTiDBAlterTable pins spec §3.3: every subcommand is judged
// against the table as it was, two subcommands on one column or index fail
// from the text, and RENAME TO takes no other subcommand.
func TestWalkThroughTiDBAlterTable(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// MySQL runs the drops first; TiDB judges each against the table as it
	// was (1060, 1061, 1054).
	oneFinding(t, "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;", 0, "Column a already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP INDEX i, ADD INDEX i (b);", 0, "Index i already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t RENAME COLUMN a TO c, ADD COLUMN a int;", 0, "Column a already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t RENAME COLUMN a TO b, RENAME COLUMN b TO a;", 0, "Column b already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN y1 int, ADD COLUMN y2 int AFTER y1;", 0, "Column y1 does not exist in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN x int, MODIFY x bigint;", 0, "Column x does not exist in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (a);", 0, "Table t already has a primary key", shop)
	noFinding(t, "ALTER TABLE t DROP INDEX i, ADD INDEX j (b);\nALTER TABLE t DROP INDEX j;\nALTER TABLE t ADD INDEX i (a);", shop)
	// Repeated subcommands (8200).
	oneFinding(t, "ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;", 0, "Column c is added twice to table t", shop)
	oneFinding(t, "ALTER TABLE t DROP INDEX i, DROP INDEX I;", 0, "Index I is dropped twice from table t", shop)
	oneFinding(t, "ALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;", 0, "Index `PRIMARY` is dropped twice from table t", shop)
	oneFinding(t, "ALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);", 0, "Multiple primary keys for table nokey", shop)
	// TiDB drops fk1 and then fails (1091), so the oracle cannot run it.
	oneFinding(t, "ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;", 0, "Foreign key fk1 is dropped twice from table t", shop)
	// TiDB loops on this one's rollback, so the oracle cannot run it.
	oneFinding(t, "ALTER TABLE t ADD CONSTRAINT f2 FOREIGN KEY (a) REFERENCES r (id), ADD CONSTRAINT F2 FOREIGN KEY (b) REFERENCES r (id);", 0, "Foreign key F2 is added twice to table t", shop)
	// RENAME TO beside another subcommand (8200), but ALGORITHM and LOCK.
	oneFinding(t, "ALTER TABLE t COMMENT 'x', RENAME TO t2;", 0, "RENAME TO cannot be combined with other changes to table t", shop)
	oneFinding(t, "ALTER TABLE t RENAME TO t2, RENAME TO t3;", 0, "RENAME TO cannot be combined with other changes to table t", shop)
	noFinding(t, "ALTER TABLE t COMMENT 'x', RENAME TO t;\nALTER TABLE t ADD COLUMN zz int, RENAME TO T;", shop)
	noFinding(t, "ALTER TABLE t ALGORITHM = INSTANT, LOCK = NONE, RENAME TO t2;\nCREATE TABLE t (a int);\nALTER TABLE t2 ADD COLUMN zz int;", shop)
	// DROP COLUMN IF EXISTS of a missing column is a note, and the
	// statement goes on without it.
	noFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS zz;\nALTER TABLE t DROP COLUMN IF EXISTS zz, DROP COLUMN IF EXISTS zz;", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN q int, DROP COLUMN IF EXISTS q;\nALTER TABLE t ADD COLUMN q int;", 1, "Column q already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN IF EXISTS a;\nALTER TABLE t DROP COLUMN a;", 1, "Column a does not exist in table t", shop)
	// TiDB checks no length in CHANGE or RENAME COLUMN.
	long := strings.Repeat("c", 65)
	noFinding(t, "ALTER TABLE t CHANGE a "+long+" int;\nALTER TABLE t RENAME COLUMN b TO "+long+"2;", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN "+long+" int;", 0, "Identifier name "+long+" is too long", shop)
	// A local temporary table ends the walk (8200).
	noFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nALTER TABLE tt ADD COLUMN b int;\nCREATE TABLE t (a int);", shop)
}

// TestWalkThroughTiDBForeignKeys pins spec §3.4: names are the table's,
// generated ones are fk_<n> and survive RENAME TABLE, and ENGINE is
// ignored.
func TestWalkThroughTiDBForeignKeys(t *testing.T) {
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE x (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));", on)
	oneFinding(t, "ALTER TABLE t ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);", 0, "Foreign key FK1 already exists in table t", on)
	oneFinding(t, "ALTER TABLE nokey DROP FOREIGN KEY fk1;", 0, "Foreign key fk1 does not exist in table nokey", on)
	oneFinding(t, "ALTER TABLE t DROP FOREIGN KEY fk1, ADD CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id);", 0, "Foreign key fk1 already exists in table t", on)
	oneFinding(t, "CREATE TABLE x (a int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));", 0, "Foreign key F is specified twice in table x", on)
	noFinding(t, "CREATE TABLE g (id int, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY fk_1;", on)
	oneFinding(t, "CREATE TABLE g (id int, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY g_ibfk_1;", 2, "Foreign key g_ibfk_1 does not exist in table h", on)
	oneFinding(t, "RENAME TABLE t TO u;\nALTER TABLE u DROP FOREIGN KEY fk1;\nALTER TABLE u DROP FOREIGN KEY fk1;", 2, "Foreign key fk1 does not exist in table u", on)
	oneFinding(t, "CREATE TABLE my (a int, FOREIGN KEY (a) REFERENCES nosuch (id)) ENGINE=MyISAM;", 0, "Foreign key target table nosuch does not exist", on)
	oneFinding(t, "CREATE TABLE my (a int) ENGINE=MyISAM;\nALTER TABLE my DROP FOREIGN KEY nosuch;", 1, "Foreign key nosuch does not exist in table my", on)
	oneFinding(t, "ALTER TABLE nokey ENGINE=MyISAM;\nALTER TABLE nokey DROP FOREIGN KEY nosuch;", 1, "Foreign key nosuch does not exist in table nokey", on)
	oneFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES sq (id));", 0, "Foreign key target sq is not a base table", on)
	oneFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES v (id));", 0, "Foreign key target v is not a base table", on)
}

// TestWalkThroughTiDBObjects pins spec §3.2 and §3.5: no triggers,
// functions, events, or procedures, sequences as relations, the kinds IF
// EXISTS does not cover, and DROP forms TiDB does not run atomically.
func TestWalkThroughTiDBObjects(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// CREATE and DROP PROCEDURE fail (8108) and end the walk.
	noFinding(t, "CREATE PROCEDURE p() SELECT 1;\nCREATE TABLE t (a int);", shop)
	noFinding(t, "DROP PROCEDURE IF EXISTS p;\nCREATE TABLE t (a int);", shop)
	// ALTER VIEW, which TiDB rejects (1064), ends the walk.
	noFinding(t, "ALTER VIEW v AS SELECT 1 AS a;\nCREATE TABLE t (a int);", shop)
	// DML runs no trigger: a temporary table stays known.
	oneFinding(t, "INSERT INTO t (id) VALUES (1);\nDROP TEMPORARY TABLE nosuch;", 1, "Temporary table nosuch does not exist", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nUPDATE t JOIN r ON r.id = t.a SET t.a = 1;\nCREATE TEMPORARY TABLE tt (a int);", 2, "Table tt already exists", shop)
	// Sequences are synced relations of their own kind.
	oneFinding(t, "CREATE TABLE sq (a int);", 0, "Table sq already exists", shop)
	oneFinding(t, "CREATE SEQUENCE SQ;", 0, "Table SQ already exists", shop)
	oneFinding(t, "CREATE SEQUENCE t;", 0, "Table t already exists", shop)
	noFinding(t, "CREATE SEQUENCE IF NOT EXISTS sq;\nCREATE TEMPORARY TABLE sq (a int);", shop)
	oneFinding(t, "DROP SEQUENCE nosuch;", 0, "Sequence nosuch does not exist", shop)
	oneFinding(t, "DROP SEQUENCE t;", 0, "t is not a sequence", shop)
	noFinding(t, "DROP SEQUENCE IF EXISTS t, nosuch;\nDROP TABLE t;", shop)
	oneFinding(t, "DROP SEQUENCE sq;\nDROP SEQUENCE sq;", 1, "Sequence sq does not exist", shop)
	oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nDROP SEQUENCE tt;", 1, "Sequence tt does not exist", shop)
	oneFinding(t, "CREATE SEQUENCE s2;\nRENAME TABLE s2 TO s3;\nDROP SEQUENCE s2;", 2, "Sequence s2 does not exist", shop)
	oneFinding(t, "DROP DATABASE shop;\nDROP SEQUENCE IF EXISTS sq;", 1, "No database selected", shop)
	oneFinding(t, "ALTER TABLE sq COMMENT 'x';", 0, "sq is not a base table", shop)
	oneFinding(t, "TRUNCATE TABLE sq;", 0, "Table sq does not exist", shop)
	oneFinding(t, "DROP TABLE sq;", 0, "Table sq does not exist", shop)
	oneFinding(t, "CREATE OR REPLACE VIEW sq AS SELECT 1 AS a;", 0, "sq is not a view", shop)
	oneFinding(t, "CREATE TABLE x LIKE sq;", 0, "sq is not a base table", shop)
	// DROP TABLE IF EXISTS leaves a sequence or a view with a note.
	oneFinding(t, "DROP TABLE IF EXISTS sq;\nDROP SEQUENCE sq;\nDROP SEQUENCE sq;", 2, "Sequence sq does not exist", shop)
	noFinding(t, "DROP TABLE IF EXISTS v;\nDROP VIEW v;", shop)
	// DROP VIEW IF EXISTS of a table or a sequence fails (1347).
	oneFinding(t, "DROP VIEW IF EXISTS t;", 0, "t is not a view", shop)
	oneFinding(t, "DROP VIEW IF EXISTS sq;", 0, "sq is not a view", shop)
	noFinding(t, "DROP VIEW IF EXISTS nosuch, v;", shop)
	// CREATE INDEX of an expression succeeds on a view or a sequence: the
	// walk ends there.
	noFinding(t, "CREATE INDEX k ON v ((1 + 1));\nCREATE TABLE t (a int);", shop)
	noFinding(t, "CREATE INDEX k ON sq ((1 + 1));\nDROP INDEX k ON sq;\nCREATE TABLE t (a int);", shop)
	oneFinding(t, "CREATE INDEX k ON nosuch (a);", 0, "Table nosuch does not exist", shop)
	// TiDB drops what it can of these before it fails, so the oracle cannot
	// run them; the statement fails all the same.
	oneFinding(t, "DROP TABLE t, nosuch;", 0, "Table nosuch does not exist", shop)
	oneFinding(t, "DROP TABLE t, v;", 0, "Table v does not exist", shop)
	oneFinding(t, "DROP TABLE t, sq;", 0, "Table sq does not exist", shop)
	oneFinding(t, "DROP VIEW v, nosuch;", 0, "View nosuch does not exist", shop)
	oneFinding(t, "DROP SEQUENCE sq, nosuch;", 0, "Sequence nosuch does not exist", shop)
	// The hidden row id is Unknown; the primary key's absence is trusted.
	noFinding(t, "ALTER TABLE nokey ADD COLUMN x int AFTER _tidb_rowid;\nALTER TABLE nokey DROP COLUMN _tidb_rowid;", shop)
	noFinding(t, "CREATE TABLE g (a int);\nALTER TABLE g DROP COLUMN _tidb_rowid;", shop)
	oneFinding(t, "ALTER TABLE nokey DROP PRIMARY KEY;", 0, "Table nokey has no primary key", shop)
	oneFinding(t, "CREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;", 1, "Table g has no primary key", shop)
	// TiDB checks no length in DROP TABLE.
	long := strings.Repeat("d", 65)
	noFinding(t, "DROP TABLE IF EXISTS "+long+";", shop)
	oneFinding(t, "DROP TABLE "+long+";", 0, "Table "+long+" does not exist", shop)
}

func TestWalkThroughTiDBCase(t *testing.T) {
	// lower_case_table_names is 2 on TiDB, whatever the target says.
	for _, lctn := range []int{0, 1, 2} {
		target := review.Target{Schema: walkSchema(), LowerCaseTableNames: lctn}
		oneFinding(t, "CREATE TABLE T (a int);", 0, "Table T already exists", target)
		oneFinding(t, "CREATE DATABASE SHOP;", 0, "Database SHOP already exists", target)
		oneFinding(t, "CREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY TABLE TT (a int);", 1, "Table TT already exists", target)
	}
	// A name outside ASCII is never absent.
	noFinding(t, "ALTER TABLE t DROP COLUMN `é`;", review.Target{Schema: walkSchema()})
}

func TestWalkThroughTiDBSession(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Another database is unknown, until the change makes it.
	noFinding(t, "USE other;\nCREATE TABLE t (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 2, "Table q.x already exists", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database selected", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE SEQUENCE s;", 1, "No database selected", shop)
	// A temporary table shadows the base table.
	noFinding(t, "CREATE TEMPORARY TABLE nosuch (a int);\nALTER TABLE t ADD COLUMN zz int;\nINSERT INTO nosuch VALUES (1);", shop)
	oneFinding(t, "DROP TEMPORARY TABLE t;", 0, "Temporary table t does not exist", shop)
	// A foreign key's target is checked only while the checks are on.
	noFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", shop)
	oneFinding(t, "SET foreign_key_checks = 1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", 1, "Foreign key target table nosuch does not exist", shop)
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	noFinding(t, "SET FOREIGN_KEY_CHECKS = 0;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", on)
	// With another spec of the plan, only the change's own facts count.
	other := review.Target{Schema: walkSchema(), OtherSpecs: true}
	noFinding(t, "CREATE TABLE t (a int);", other)
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE x (a int);", 1, "Table x already exists", other)
}

// TestWalkThroughTiDBTargetGroups pins OtherSpecs and ForeignKeyChecks in
// the target grouping: two targets of one schema that differ in either are
// walked apart.
func TestWalkThroughTiDBTargetGroups(t *testing.T) {
	schema := walkSchema()
	got := walkFindings(t, "CREATE TABLE t (a int);", review.Target{Schema: schema}, review.Target{Schema: schema, OtherSpecs: true})
	if len(got) != 1 || len(got[0].Targets) != 1 || got[0].Targets[0] != 0 {
		t.Errorf("OtherSpecs: findings = %+v, want one on target 0", got)
	}
	got = walkFindings(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", review.Target{Schema: schema}, review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn})
	if len(got) != 1 || len(got[0].Targets) != 1 || got[0].Targets[0] != 1 {
		t.Errorf("ForeignKeyChecks: findings = %+v, want one on target 1", got)
	}
}

// walkForms are the statement forms the walk lowers (the MySQL spec's §5
// table, with spec §3's TiDB deltas), each on walkSchema. A form the table
// lists must lower.
var walkForms = []walktest.Form{
	{SQL: "CREATE DATABASE q;"},
	{SQL: "CREATE DATABASE IF NOT EXISTS shop;"},
	{SQL: "DROP DATABASE IF EXISTS q;"},
	{SQL: "ALTER DATABASE shop CHARACTER SET utf8mb4;"},
	{SQL: "USE shop;"},
	{SQL: "CREATE TABLE x (a int PRIMARY KEY, b int UNIQUE, KEY k (b), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id), CHECK (a > 0));"},
	{SQL: "CREATE TABLE IF NOT EXISTS t (a int);"},
	{SQL: "CREATE TEMPORARY TABLE tmp (a int);"},
	{SQL: "CREATE TABLE x LIKE t;"},
	{SQL: "CREATE TABLE x AS SELECT 1 AS a;"},
	{SQL: "ALTER TABLE t ADD COLUMN c int AFTER a, ADD COLUMN (d int, e int);"},
	{SQL: "ALTER TABLE t DROP COLUMN a;"},
	{SQL: "ALTER TABLE t DROP COLUMN IF EXISTS zz;"},
	{SQL: "ALTER TABLE t MODIFY COLUMN a bigint;"},
	{SQL: "ALTER TABLE t CHANGE a c int;"},
	{SQL: "ALTER TABLE t RENAME COLUMN a TO c;"},
	{SQL: "ALTER TABLE t ALTER COLUMN a SET DEFAULT 1;"},
	{SQL: "ALTER TABLE t ADD INDEX k (a), ADD UNIQUE (b), ADD FULLTEXT INDEX ft (a);"},
	{SQL: "ALTER TABLE nokey ADD PRIMARY KEY (a);"},
	{SQL: "ALTER TABLE nokey ADD CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES r (id);"},
	{SQL: "ALTER TABLE t ADD CONSTRAINT c CHECK (a > 0);"},
	{SQL: "ALTER TABLE t DROP INDEX i, DROP KEY i2;"},
	{SQL: "ALTER TABLE t DROP PRIMARY KEY;"},
	{SQL: "ALTER TABLE t DROP FOREIGN KEY fk1;"},
	{SQL: "ALTER TABLE t DROP CHECK c;"},
	{SQL: "ALTER TABLE t DROP CONSTRAINT c;"},
	{SQL: "ALTER TABLE t RENAME INDEX i TO j;"},
	{SQL: "ALTER TABLE t ALTER INDEX i INVISIBLE;"},
	{SQL: "ALTER TABLE t RENAME TO t2;"},
	{SQL: "ALTER TABLE t ALGORITHM = INSTANT, LOCK = NONE, RENAME TO t2;"},
	{SQL: "ALTER TABLE t ENGINE = InnoDB, COMMENT 'x', ALGORITHM = INPLACE, LOCK = NONE;"},
	{SQL: "ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4;"},
	{SQL: "ALTER TABLE t FORCE;"},
	{SQL: "RENAME TABLE t TO t2, r TO r2;"},
	{SQL: "DROP TABLE IF EXISTS t, nosuch;"},
	{SQL: "DROP TEMPORARY TABLE IF EXISTS tmp;"},
	{SQL: "TRUNCATE TABLE t;"},
	{SQL: "CREATE UNIQUE INDEX k ON t (a);"},
	{SQL: "CREATE INDEX IF NOT EXISTS i ON t (a);"},
	{SQL: "DROP INDEX i ON t;"},
	{SQL: "CREATE VIEW w AS SELECT 1 AS a;"},
	{SQL: "CREATE OR REPLACE VIEW v AS SELECT 1 AS a;"},
	{SQL: "DROP VIEW IF EXISTS v, w;"},
	{SQL: "CREATE SEQUENCE s2;"},
	{SQL: "CREATE SEQUENCE IF NOT EXISTS sq;"},
	{SQL: "DROP SEQUENCE IF EXISTS sq, s2;"},
	{SQL: "INSERT INTO nokey VALUES (1);"},
	{SQL: "REPLACE INTO nokey VALUES (1);"},
	{SQL: "UPDATE nokey SET a = 1;"},
	{SQL: "DELETE FROM nokey;"},
	{SQL: "SELECT 1;"},
	{SQL: "SET foreign_key_checks = 0;\nSET @@session.foreign_key_checks = 1;\nSET NAMES utf8mb4;"},
	{SQL: "BEGIN;\nSAVEPOINT a;\nRELEASE SAVEPOINT a;\nCOMMIT;"},
	{SQL: "SHOW TABLES;\nANALYZE TABLE t;"},
	{SQL: "GRANT SELECT ON shop.t TO 'u'@'%';"},
}

func TestWalkThroughTiDBLowersEveryForm(t *testing.T) {
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
