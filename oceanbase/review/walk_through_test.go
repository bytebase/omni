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
					Columns:     cols("id", "a"),
					Indexes:     []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}, {Name: "i"}},
					ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk1", ReferencedTable: "r"}},
					Triggers:    []*metadata.TriggerMetadata{{Name: "tr"}},
				},
				{Name: "r", Columns: cols("id"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true}}},
				{Name: "nokey", Columns: cols("a")},
			},
			Views:      []*metadata.ViewMetadata{{Name: "v"}},
			Procedures: []*metadata.ProcedureMetadata{{Name: "p"}},
			Functions:  []*metadata.FunctionMetadata{{Name: "f"}},
			Events:     []*metadata.EventMetadata{{Name: "e"}},
		}},
	}
}

func walkFindings(t *testing.T, sql string, targets ...review.Target) []review.Finding {
	t.Helper()
	return reviewFindings(t, Review, sql, targets...)
}

func reviewFindings(t *testing.T, rev review.ReviewFunc, sql string, targets ...review.Target) []review.Finding {
	t.Helper()
	result, err := rev(context.Background(), sql, review.Options{Rules: []review.Rule{review.WalkThrough}}, targets)
	if err != nil {
		t.Fatal(err)
	}
	return result.Findings
}

// oneFinding asserts the change's one WALK_THROUGH finding.
func oneFinding(t *testing.T, sql string, statement int, message string, target review.Target) {
	t.Helper()
	oneFindingOf(t, Review, sql, statement, message, target)
}

func noFinding(t *testing.T, sql string, target review.Target) {
	t.Helper()
	noFindingOf(t, Review, sql, target)
}

func oneFindingOf(t *testing.T, rev review.ReviewFunc, sql string, statement int, message string, target review.Target) {
	t.Helper()
	got := reviewFindings(t, rev, sql, target)
	if len(got) != 1 || got[0].Statement != statement || got[0].Message != message {
		t.Errorf("%s\nfindings = %+v, want statement %d: %q", sql, got, statement, message)
	}
}

func noFindingOf(t *testing.T, rev review.ReviewFunc, sql string, target review.Target) {
	t.Helper()
	if got := reviewFindings(t, rev, sql, target); len(got) != 0 {
		t.Errorf("%s\nfindings = %+v, want none", sql, got)
	}
}

func TestWalkThroughMySQL(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table t already exists", shop)
	oneFinding(t, "ALTER TABLE nosuch ADD COLUMN b int;", 0, "Table nosuch does not exist", shop)
	oneFinding(t, "ALTER TABLE v ADD COLUMN b int;", 0, "v is not a base table", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN A int;", 0, "Column A already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column zz does not exist in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX I (a);", 0, "Index I already exists in table t", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "Table t already has a primary key", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);", 0, "Foreign key FK1 already exists in the database", shop)
	oneFinding(t, "DROP TABLE t, nosuch;", 0, "Table nosuch does not exist", shop)
	oneFinding(t, "DROP VIEW t;", 0, "t is not a view", shop)
	oneFinding(t, "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", 0, "Trigger tr already exists", shop)
	oneFinding(t, "DROP PROCEDURE nosuch;", 0, "PROCEDURE nosuch does not exist", shop)
	oneFinding(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 0, "Event e already exists", shop)
	// ALTER TABLE judges drops against the table as it was.
	noFinding(t, "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN b int, ADD COLUMN b int;", 0, "Column b is added twice to table t", shop)
	// RENAME TABLE runs its pairs in order.
	noFinding(t, "RENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN zz int;\nALTER TABLE t ADD COLUMN zz int;", shop)
	oneFinding(t, "RENAME TABLE t TO u;\nCREATE TABLE u (a int);", 1, "Table u already exists", shop)
}

func TestWalkThroughMySQLCase(t *testing.T) {
	sensitive := review.Target{Schema: walkSchema()}
	folded := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1}
	noFinding(t, "CREATE TABLE T (a int);", sensitive)
	oneFinding(t, "CREATE TABLE T (a int);", 0, "Table T already exists", folded)
	// A trigger is claimed taken only by its own spelling: 5.7 compares
	// trigger names with case.
	noFinding(t, "CREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", sensitive)
	// On 5.7, TR beside tr is a second trigger, and DROP TRIGGER TR leaves
	// tr (verified); the oracle runs no 5.7.
	noFinding(t, "CREATE TRIGGER TR BEFORE UPDATE ON nokey FOR EACH ROW SET NEW.a = 1;\nDROP TRIGGER TR;\nDROP TRIGGER tr;", sensitive)
	noFinding(t, "DROP TRIGGER TR;\nDROP TRIGGER tr;", sensitive)
	both := walkSchema()
	both.Schemas[0].Tables[0].Triggers = append(both.Schemas[0].Tables[0].Triggers, &metadata.TriggerMetadata{Name: "TR"})
	noFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER TR;", review.Target{Schema: both})
	oneFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER tr;", 1, "Trigger tr does not exist", sensitive)
}

// TestWalkThroughMySQLNamesOutsideASCII pins that the walk follows no name
// outside ASCII, which MySQL compares by rules of its own per space: a
// statement that names one ends the walk, and a synced space that lists
// one has no name known to be absent, since the Kelvin sign is k as a
// column and é is e as a routine or a trigger.
func TestWalkThroughMySQLNamesOutsideASCII(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "ALTER TABLE t DROP COLUMN `é`;", shop)
	noFinding(t, "CREATE TABLE x (id INT, café INT);\nALTER TABLE x DROP COLUMN CAFé;", shop)
	noFinding(t, "CREATE TABLE x (k INT);\nALTER TABLE x DROP COLUMN `\u212a`;\nALTER TABLE x DROP COLUMN zz;", shop)
	noFinding(t, "CREATE DATABASE q;\nCREATE PROCEDURE q.`pé`() SELECT 1;\nDROP PROCEDURE q.pe;", shop)
	// A string is no name.
	oneFinding(t, "CREATE TABLE x (a int) COMMENT '中文';\nCREATE TABLE x (a int);", 1, "Table x already exists", shop)
	oneFinding(t, "INSERT INTO nokey VALUES ('é');\nDROP TEMPORARY TABLE nosuch;", 1, "Temporary table nosuch does not exist", shop)
	mixed := walkSchema()
	schema := mixed.Schemas[0]
	schema.Tables[2].Columns = append(schema.Tables[2].Columns, &metadata.ColumnMetadata{Name: "\u212a"})
	schema.Tables[0].Triggers = append(schema.Tables[0].Triggers, &metadata.TriggerMetadata{Name: "tré"})
	schema.Tables = append(schema.Tables, &metadata.TableMetadata{Name: "t\u212a"})
	schema.Procedures = append(schema.Procedures, &metadata.ProcedureMetadata{Name: "pé"})
	target := review.Target{Schema: mixed}
	noFinding(t, "ALTER TABLE nokey DROP COLUMN k;\nALTER TABLE nokey DROP COLUMN zz;\nDROP PROCEDURE pe;\nDROP PROCEDURE nosuch;\nDROP TRIGGER tre;\nDROP TRIGGER nosuch;", target)
	noFinding(t, "DROP TABLE tk;\nDROP TABLE nosuch;", review.Target{Schema: mixed, LowerCaseTableNames: 1})
	// A space that lists none still knows its absent names.
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column zz does not exist in table t", target)
	oneFinding(t, "DROP EVENT nosuch;", 0, "Event nosuch does not exist", target)
}

func TestWalkThroughMySQLSession(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Another database is unknown, until the change makes it.
	noFinding(t, "USE other;\nCREATE TABLE t (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 2, "Table q.x already exists", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database selected", shop)
	// A temporary table shadows the base table.
	noFinding(t, "CREATE TEMPORARY TABLE nosuch (a int);\nALTER TABLE t ADD COLUMN zz int;\nINSERT INTO nosuch VALUES (1);", shop)
	oneFinding(t, "DROP TEMPORARY TABLE t;", 0, "Temporary table t does not exist", shop)
	// A trigger may create a temporary table of any name.
	noFinding(t, "INSERT INTO t VALUES (1, 1);\nINSERT INTO nosuch VALUES (1);", shop)
	oneFinding(t, "INSERT INTO nokey VALUES (1);\nINSERT INTO nosuch VALUES (1);", 1, "Table nosuch does not exist", shop)
	// A foreign key's target is checked only while the checks are on.
	noFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", shop)
	oneFinding(t, "SET foreign_key_checks = 1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", 1, "Foreign key target table nosuch does not exist", shop)
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	oneFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", 0, "Foreign key target table nosuch does not exist", on)
	noFinding(t, "SET FOREIGN_KEY_CHECKS = 0;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", on)
	// Each assignment has a scope of its own: an unprefixed one the
	// statement's, past a prefixed one, and an @@session or bare @@ one the
	// session's. A global one leaves the session's checks as they are.
	const fk = "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));"
	const missing = "Foreign key target table nosuch does not exist"
	oneFinding(t, "SET GLOBAL foreign_key_checks = 0;\n"+fk, 1, missing, on)
	oneFinding(t, "SET @@global.foreign_key_checks = 0;\n"+fk, 1, missing, on)
	oneFinding(t, "SET PERSIST foreign_key_checks = 0;\n"+fk, 1, missing, on)
	oneFinding(t, "SET PERSIST_ONLY foreign_key_checks = 0;\n"+fk, 1, missing, on)
	noFinding(t, "SET GLOBAL foreign_key_checks = 1;\n"+fk, shop)
	oneFinding(t, "SET GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, 1, missing, on)
	oneFinding(t, "SET GLOBAL sql_mode = DEFAULT, @@session.foreign_key_checks = 1, foreign_key_checks = 0;\n"+fk, 1, missing, shop)
	noFinding(t, "SET GLOBAL sql_mode = DEFAULT, @@session.foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "SET GLOBAL sql_mode = DEFAULT, @@foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "SET PERSIST sql_mode = DEFAULT, @@local.foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "SET @@global.foreign_key_checks = 1, foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "SET SESSION foreign_key_checks = 0;\n"+fk, on)
	noFinding(t, "SET LOCAL foreign_key_checks = 0;\n"+fk, on)
	// With another spec of the plan, only the change's own facts count.
	other := review.Target{Schema: walkSchema(), OtherSpecs: true}
	noFinding(t, "CREATE TABLE t (a int);", other)
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE x (a int);", 1, "Table x already exists", other)
}

// TestWalkThroughMySQLForeignKeyNames pins the names a foreign key takes:
// the CONSTRAINT symbol, or without one a generated name, never the name
// after FOREIGN KEY, which names only its index (verified on 8.0.46). A
// long name there fails only without a symbol.
func TestWalkThroughMySQLForeignKeyNames(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	long := strings.Repeat("f", 65)
	noFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY cx;\nALTER TABLE x DROP INDEX cx;", shop)
	oneFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY fx;", 1, "Foreign key fx does not exist in table x", shop)
	noFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nALTER TABLE y DROP FOREIGN KEY y_ibfk_1;", shop)
	noFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nCREATE TABLE z (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY f (a) REFERENCES r (id), ADD CONSTRAINT nk2 FOREIGN KEY f (a) REFERENCES r (id);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY nk (a) REFERENCES r (id);", 0, "Foreign key fk1 already exists in the database", shop)
	noFinding(t, "CREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id));", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id);", shop)
	oneFinding(t, "CREATE TABLE w (a int, FOREIGN KEY "+long+" (a) REFERENCES r (id));", 0, "Identifier name "+long+" is too long", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY "+long+" (a) REFERENCES r (id);", 0, "Identifier name "+long+" is too long", shop)
}

// TestWalkThroughMySQLNoDatabase pins what fails once the change drops
// the current database: an unqualified name (1046), but for DROP FUNCTION,
// which then looks for a loadable function alone (verified on 8.0.46).
func TestWalkThroughMySQLNoDatabase(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS nosuch;\nDROP FUNCTION nosuch;\nDROP FUNCTION IF EXISTS shop.f;", shop)
	for _, stmt := range []string{"DROP PROCEDURE IF EXISTS p;", "DROP TRIGGER IF EXISTS tr;", "DROP VIEW IF EXISTS v;", "DROP EVENT IF EXISTS e;", "DROP TABLE IF EXISTS nokey;", "CREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;", "ALTER FUNCTION f COMMENT 'x';"} {
		oneFinding(t, "DROP DATABASE shop;\n"+stmt, 1, "No database selected", shop)
	}
}

// TestWalkThroughMySQLParentBesideRename pins that an ALTER TABLE that
// also renames its table checks no foreign key's parent: MySQL resolves
// it in the database the table moves to (verified on 8.0.46).
func TestWalkThroughMySQLParentBesideRename(t *testing.T) {
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	noFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.p2 (id int PRIMARY KEY);\nALTER TABLE nokey RENAME TO q.nokey, ADD FOREIGN KEY (a) REFERENCES p2 (id);", on)
	noFinding(t, "ALTER TABLE nokey RENAME TO nk, ADD FOREIGN KEY (a) REFERENCES nosuch (id);", on)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY (a) REFERENCES nosuch (id);", 0, "Foreign key target table nosuch does not exist", on)
}

// TestWalkThroughMySQLStoredFunctions pins the calls that may run a stored
// function, which may create a temporary table (verified on 8.0.46): one a
// view's query makes, and an unqualified one after a USE of a database the
// change did not make.
func TestWalkThroughMySQLStoredFunctions(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	const drop = "DROP TEMPORARY TABLE made;"
	noFinding(t, "SELECT * FROM v;\n"+drop, shop)
	noFinding(t, "INSERT INTO nokey SELECT a FROM q.qv;\n"+drop, shop)
	noFinding(t, "SELECT 1 FROM nokey WHERE a IN (SELECT a FROM v);\n"+drop, shop)
	noFinding(t, "USE q;\nSELECT mk();\nUSE shop;\n"+drop, shop)
	noFinding(t, "SELECT f();\n"+drop, shop)
	noFinding(t, "SELECT q.mk();\n"+drop, shop)
	// A call found stays found past a base table read after it.
	noFinding(t, "SELECT f() FROM r;\n"+drop, shop)
	noFinding(t, "SELECT q.mk() FROM nokey WHERE a IN (SELECT id FROM r);\n"+drop, shop)
	// A base table runs no function, nor does a call to one the current
	// database certainly lacks.
	oneFinding(t, "SELECT * FROM nokey;\n"+drop, 1, "Temporary table made does not exist", shop)
	oneFinding(t, "SELECT mk();\n"+drop, 1, "Temporary table made does not exist", shop)
	oneFinding(t, "CREATE DATABASE q;\nUSE q;\nSELECT mk();\nUSE shop;\n"+drop, 4, "Temporary table made does not exist", shop)
}

// TestWalkThroughMySQLTargetGroups pins OtherSpecs and ForeignKeyChecks in
// the target grouping: two targets of one schema that differ in either are
// walked apart.
func TestWalkThroughMySQLTargetGroups(t *testing.T) {
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

// walkForms are the statement forms of the MySQL spec's §5 table, each on
// walkSchema. A form the table lists must lower.
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
	{SQL: "ALTER TABLE t ADD COLUMN b int AFTER a, ADD COLUMN (c int, d int);"},
	{SQL: "ALTER TABLE t DROP COLUMN a;"},
	{SQL: "ALTER TABLE t MODIFY COLUMN a bigint;"},
	{SQL: "ALTER TABLE t CHANGE a b int;"},
	{SQL: "ALTER TABLE t RENAME COLUMN a TO b;"},
	{SQL: "ALTER TABLE t ALTER COLUMN a SET DEFAULT 1;"},
	{SQL: "ALTER TABLE t ADD INDEX k (a), ADD UNIQUE (a), ADD FULLTEXT INDEX ft (a);"},
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
	{SQL: "ALTER TABLE t ENGINE = InnoDB, COMMENT 'x', ALGORITHM = INPLACE, LOCK = NONE;"},
	{SQL: "ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4;"},
	{SQL: "ALTER TABLE t FORCE;"},
	{SQL: "RENAME TABLE t TO t2, r TO r2;"},
	{SQL: "DROP TABLE IF EXISTS t, nosuch;"},
	{SQL: "DROP TEMPORARY TABLE IF EXISTS tmp;"},
	{SQL: "TRUNCATE TABLE t;"},
	{SQL: "CREATE UNIQUE INDEX k ON t (a);"},
	{SQL: "DROP INDEX i ON t;"},
	{SQL: "CREATE VIEW w AS SELECT 1 AS a;"},
	{SQL: "CREATE OR REPLACE VIEW v AS SELECT 1 AS a;"},
	{SQL: "ALTER VIEW v AS SELECT 2 AS a;"},
	{SQL: "DROP VIEW IF EXISTS v, w;"},
	{SQL: "CREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;"},
	{SQL: "DROP TRIGGER IF EXISTS tr;"},
	{SQL: "CREATE PROCEDURE p2() SELECT 1;"},
	{SQL: "CREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;"},
	{SQL: "ALTER PROCEDURE p COMMENT 'x';"},
	{SQL: "DROP FUNCTION IF EXISTS f;"},
	{SQL: "CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;"},
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

func TestWalkThroughMySQLLowersEveryForm(t *testing.T) {
	walktest.Coverage(t, walkForms, formLowering(walkSchema(), dialect{}))
}

// formLowering walks a form on the schema in the dialect, as
// walktest.Coverage asks.
func formLowering(schema *metadata.DatabaseSchemaMetadata, d dialect) func(t *testing.T, sql string) (n, lowered int, last bool) {
	return func(t *testing.T, sql string) (n, lowered int, last bool) {
		t.Helper()
		ranges := statementRanges(sql)
		stmts, syntax := parse(sql, ranges)
		if syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, sql)
		}
		w := newWalker(sql, stmts, review.Target{Schema: schema}, d)
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
	}
}
