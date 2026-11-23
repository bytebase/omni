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
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE t (a int);", 1, "Table `t` already exists.", shop)
	oneFinding(t, "ALTER TABLE nosuch ADD COLUMN b int;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE v ADD COLUMN b int;", 0, "`v` is a view, not a table.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN A int;", 0, "Column `A` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column `zz` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD INDEX I (a);", 0, "Index `I` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (a);", 0, "`t` already has a primary key.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);", 0, "Foreign key `FK1` already exists.", shop)
	oneFinding(t, "DROP TABLE t, nosuch;", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP VIEW t;", 0, "`t` is a table, not a view.", shop)
	oneFinding(t, "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", 0, "Trigger `tr` already exists.", shop)
	oneFinding(t, "DROP PROCEDURE nosuch;", 0, "Procedure `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;", 0, "Event `e` already exists.", shop)
	// ALTER TABLE judges drops against the table as it was.
	noFinding(t, "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN b int, ADD COLUMN b int;", 0, "This statement adds column `b` twice.", shop)
	// RENAME TABLE runs its pairs in order.
	noFinding(t, "RENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN zz int;\nALTER TABLE t ADD COLUMN zz int;", shop)
	oneFinding(t, "RENAME TABLE t TO u;\nCREATE TABLE u (a int);", 1, "Table `u` already exists.", shop)
}

func TestWalkThroughMySQLCase(t *testing.T) {
	sensitive := review.Target{Schema: walkSchema()}
	folded := review.Target{Schema: walkSchema(), LowerCaseTableNames: 1}
	noFinding(t, "CREATE TABLE T (a int);", sensitive)
	oneFinding(t, "CREATE TABLE T (a int);", 0, "Table `T` already exists.", folded)
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
	oneFinding(t, "DROP TRIGGER tr;\nDROP TRIGGER tr;", 1, "Trigger `tr` doesn't exist.", sensitive)
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
	oneFinding(t, "CREATE TABLE x (a int) COMMENT '中文';\nCREATE TABLE x (a int);", 1, "Table `x` already exists.", shop)
	oneFinding(t, "INSERT INTO nokey VALUES ('é');\nDROP TEMPORARY TABLE nosuch;", 1, "Temporary table `nosuch` doesn't exist.", shop)
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
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column `zz` doesn't exist in `t`.", target)
	oneFinding(t, "DROP EVENT nosuch;", 0, "Event `nosuch` doesn't exist.", target)
}

func TestWalkThroughMySQLSession(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Another database is unknown, until the change makes it.
	noFinding(t, "USE other;\nCREATE TABLE t (a int);", shop)
	oneFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);", 2, "Table `q.x` already exists.", shop)
	oneFinding(t, "DROP DATABASE shop;\nCREATE TABLE x (a int);", 1, "No database is selected for `x`.", shop)
	// A temporary table shadows the base table.
	noFinding(t, "CREATE TEMPORARY TABLE nosuch (a int);\nALTER TABLE t ADD COLUMN zz int;\nINSERT INTO nosuch VALUES (1);", shop)
	oneFinding(t, "DROP TEMPORARY TABLE t;", 0, "Temporary table `t` doesn't exist.", shop)
	// A trigger may create a temporary table of any name.
	noFinding(t, "INSERT INTO t VALUES (1, 1);\nINSERT INTO nosuch VALUES (1);", shop)
	oneFinding(t, "INSERT INTO nokey VALUES (1);\nINSERT INTO nosuch VALUES (1);", 1, "`nosuch` doesn't exist.", shop)
	// A foreign key's target is checked only while the checks are on.
	noFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", shop)
	oneFinding(t, "SET foreign_key_checks = 1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", 1, "Referenced table `nosuch` doesn't exist.", shop)
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	oneFinding(t, "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", 0, "Referenced table `nosuch` doesn't exist.", on)
	noFinding(t, "SET FOREIGN_KEY_CHECKS = 0;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", on)
	// Each assignment has a scope of its own: an unprefixed one the
	// statement's, past a prefixed one, and an @@session or bare @@ one the
	// session's. A global one leaves the session's checks as they are.
	const fk = "CREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));"
	const missing = "Referenced table `nosuch` doesn't exist."
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
	// A later assignment writes its own scope, and the most recent scope reaches one that
	// writes none.
	oneFinding(t, "SET @a = 1, GLOBAL foreign_key_checks = 0;\n"+fk, 1, missing, on)
	oneFinding(t, "SET @a = 1, PERSIST foreign_key_checks = 0;\n"+fk, 1, missing, on)
	noFinding(t, "SET GLOBAL sql_mode = DEFAULT, SESSION foreign_key_checks = 0;\n"+fk, on)
	oneFinding(t, "SET SESSION sql_mode = DEFAULT, GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, 1, missing, on)
	noFinding(t, "SET GLOBAL sql_mode = DEFAULT, SESSION sql_mode = DEFAULT, foreign_key_checks = 0;\n"+fk, on)
	// With another spec of the plan, only the change's own facts count.
	other := review.Target{Schema: walkSchema(), OtherSpecs: true}
	noFinding(t, "CREATE TABLE t (a int);", other)
	oneFinding(t, "CREATE TABLE x (a int);\nCREATE TABLE x (a int);", 1, "Table `x` already exists.", other)
}

// TestWalkThroughMySQLForeignKeyNames pins the names a foreign key takes:
// the CONSTRAINT symbol, or without one a generated name, never the name
// after FOREIGN KEY, which names only its index (verified on 8.0.46). A
// long name there fails only without a symbol.
func TestWalkThroughMySQLForeignKeyNames(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	long := strings.Repeat("f", 65)
	noFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY cx;\nALTER TABLE x DROP INDEX cx;", shop)
	oneFinding(t, "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY fx;", 1, "Foreign key `fx` doesn't exist in `x`.", shop)
	noFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nALTER TABLE y DROP FOREIGN KEY y_ibfk_1;", shop)
	noFinding(t, "CREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nCREATE TABLE z (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY f (a) REFERENCES r (id), ADD CONSTRAINT nk2 FOREIGN KEY f (a) REFERENCES r (id);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY nk (a) REFERENCES r (id);", 0, "Foreign key `fk1` already exists.", shop)
	noFinding(t, "CREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id));", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT cw FOREIGN KEY "+long+" (a) REFERENCES r (id);", shop)
	oneFinding(t, "CREATE TABLE w (a int, FOREIGN KEY "+long+" (a) REFERENCES r (id));", 0, "The name `"+long+"` is longer than MySQL's 64-character limit.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY "+long+" (a) REFERENCES r (id);", 0, "The name `"+long+"` is longer than MySQL's 64-character limit.", shop)
}

// TestWalkThroughMySQLNoDatabase pins what fails once the change drops
// the current database: an unqualified name (1046), but for DROP FUNCTION,
// which then looks for a loadable function alone (verified on 8.0.46).
func TestWalkThroughMySQLNoDatabase(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS nosuch;\nDROP FUNCTION nosuch;\nDROP FUNCTION IF EXISTS shop.f;", shop)
	for stmt, name := range map[string]string{"DROP PROCEDURE IF EXISTS p;": "p", "DROP TRIGGER IF EXISTS tr;": "tr", "DROP VIEW IF EXISTS v;": "v", "DROP EVENT IF EXISTS e;": "e", "DROP TABLE IF EXISTS nokey;": "nokey", "CREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;": "f2", "ALTER FUNCTION f COMMENT 'x';": "f"} {
		oneFinding(t, "DROP DATABASE shop;\n"+stmt, 1, "No database is selected for `"+name+"`.", shop)
	}
}

// TestWalkThroughMySQLParentBesideRename pins that an ALTER TABLE that
// also renames its table checks no foreign key's parent: MySQL resolves
// it in the database the table moves to (verified on 8.0.46).
func TestWalkThroughMySQLParentBesideRename(t *testing.T) {
	on := review.Target{Schema: walkSchema(), ForeignKeyChecks: review.SwitchOn}
	noFinding(t, "CREATE DATABASE q;\nCREATE TABLE q.p2 (id int PRIMARY KEY);\nALTER TABLE nokey RENAME TO q.nokey, ADD FOREIGN KEY (a) REFERENCES p2 (id);", on)
	noFinding(t, "ALTER TABLE nokey RENAME TO nk, ADD FOREIGN KEY (a) REFERENCES nosuch (id);", on)
	oneFinding(t, "ALTER TABLE nokey ADD FOREIGN KEY (a) REFERENCES nosuch (id);", 0, "Referenced table `nosuch` doesn't exist.", on)
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
	oneFinding(t, "SELECT * FROM nokey;\n"+drop, 1, "Temporary table `made` doesn't exist.", shop)
	oneFinding(t, "SELECT mk();\n"+drop, 1, "Temporary table `made` doesn't exist.", shop)
	oneFinding(t, "CREATE DATABASE q;\nUSE q;\nSELECT mk();\nUSE shop;\n"+drop, 4, "Temporary table `made` doesn't exist.", shop)
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

// TestWalkThroughMySQLPartitionMaintenance pins that ALTER TABLE with only
// ANALYZE, CHECK, OPTIMIZE, or REPAIR PARTITION reports a status row instead
// of a SQL error on a missing table or view (measured on mysql:8.0.46), so
// WALK_THROUGH must stay silent on them.
func TestWalkThroughMySQLPartitionMaintenance(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, stmt := range []string{
		"ALTER TABLE nosuch OPTIMIZE PARTITION ALL;",
		"ALTER TABLE nosuch ANALYZE PARTITION ALL;",
		"ALTER TABLE nosuch CHECK PARTITION ALL;",
		"ALTER TABLE nosuch REPAIR PARTITION ALL;",
		"ALTER TABLE v CHECK PARTITION ALL;",
	} {
		noFinding(t, stmt, shop)
	}
	oneFinding(t, "ALTER TABLE nosuch REBUILD PARTITION ALL;", 0, "Table `nosuch` doesn't exist.", shop)
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
	walktest.Coverage(t, walkForms, formLowering(walkSchema()))
}

func TestWalkThroughMySQLFKGeneratedIndex(t *testing.T) {
	schema := &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name: "parent",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}},
					},
				},
				{
					Name: "child",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id"},
						{Name: "pid"},
						{Name: "other"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}},
						{Name: "fk_pid", Type: "BTREE", Expressions: []string{"pid"}},
						{Name: "idx_other", Type: "BTREE", Expressions: []string{"other"}},
						{Name: "uq_pid", Type: "BTREE", Unique: true, Expressions: []string{"pid"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{
						{Name: "fk_pid", Columns: []string{"pid"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}},
					},
				},
			},
		}},
	}
	target := review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn}

	// Synced index fk_pid matches the FK's columns and could be HA_GENERATED_KEY,
	// so CREATE/ADD INDEX reusing its name must not false-positive.
	noFinding(t, "CREATE INDEX fk_pid ON child (pid);", target)
	noFinding(t, "CREATE UNIQUE INDEX fk_pid ON child (pid);", target)
	noFinding(t, "ALTER TABLE child ADD KEY fk_pid (pid);", target)

	// Once created explicitly, a second CREATE INDEX of the same name fails (1061).
	oneFinding(t, "CREATE INDEX fk_pid ON child (pid);\nCREATE INDEX fk_pid ON child (pid);", 1, "Index `fk_pid` already exists in `child`.", target)

	// Non-matching or UNIQUE indexes cannot be HA_GENERATED_KEY and still report 1061.
	oneFinding(t, "CREATE INDEX idx_other ON child (other);", 0, "Index `idx_other` already exists in `child`.", target)
	oneFinding(t, "CREATE INDEX uq_pid ON child (pid);", 0, "Index `uq_pid` already exists in `child`.", target)

	// Adding another index may drop the generated fk_pid, so a later DROP INDEX fk_pid
	// must stay silent; without an intervening ADD, dropping twice still catches the second.
	noFinding(t, "ALTER TABLE child ADD KEY idx_new (pid);\nDROP INDEX fk_pid ON child;", target)
	oneFinding(t, "DROP INDEX fk_pid ON child;\nDROP INDEX fk_pid ON child;", 1, "Index `fk_pid` doesn't exist in `child`.", target)

	// RENAME TABLE carries the FK-generated index candidate to the new table name.
	noFinding(t, "RENAME TABLE child TO child2;\nCREATE INDEX fk_pid ON child2 (pid);", target)

	// Explicitly adding fk_pid in a multi-command ALTER preserves its Present fact
	// regardless of whether an unrelated ADD INDEX comes before or after it.
	oneFinding(t, "ALTER TABLE child ADD INDEX fk_pid (pid), ADD INDEX idx_new (other);\nCREATE INDEX fk_pid ON child (pid);", 1, "Index `fk_pid` already exists in `child`.", target)
	oneFinding(t, "ALTER TABLE child ADD INDEX idx_new (other), ADD INDEX fk_pid (pid);\nCREATE INDEX fk_pid ON child (pid);", 1, "Index `fk_pid` already exists in `child`.", target)

	// Dropping or renaming fk_pid alongside another ADD INDEX in the same ALTER
	// preserves its Absent fact for subsequent statements.
	oneFinding(t, "ALTER TABLE child DROP INDEX fk_pid, ADD INDEX idx_new (other);\nDROP INDEX fk_pid ON child;", 1, "Index `fk_pid` doesn't exist in `child`.", target)
	oneFinding(t, "ALTER TABLE child RENAME INDEX fk_pid TO fk_pid_explicit, ADD INDEX idx_new (other);\nDROP INDEX fk_pid ON child;", 1, "Index `fk_pid` doesn't exist in `child`.", target)

	// Dropping the backing FK retires fk_pid as a generated candidate for later
	// statements (where creating fk_pid again fails with 1061), while an ADD KEY
	// in the same ALTER still succeeds without a false positive.
	oneFinding(t, "ALTER TABLE child DROP FOREIGN KEY fk_pid;\nCREATE INDEX fk_pid ON child (pid);", 1, "Index `fk_pid` already exists in `child`.", target)
	noFinding(t, "ALTER TABLE child DROP FOREIGN KEY fk_pid, ADD KEY fk_pid (pid);", target)

	// RENAME TABLE rewrites generated <tbl>_ibfk_<n> FK owner names when moving
	// index candidates so dropping the renamed FK retires the candidate.
	schemaIbfk := &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				schema.Schemas[0].Tables[0],
				{
					Name: "child",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id"},
						{Name: "pid"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}},
						{Name: "pid", Type: "BTREE", Expressions: []string{"pid"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{
						{Name: "child_ibfk_1", Columns: []string{"pid"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}},
					},
				},
			},
		}},
	}
	targetIbfk := review.Target{Schema: schemaIbfk, ForeignKeyChecks: review.SwitchOn}
	oneFinding(t, "RENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nCREATE INDEX pid ON kid (pid);", 2, "Index `pid` already exists in `kid`.", targetIbfk)
}

// TestWalkThroughMySQLRefusalFamilies pins every MySQL 8.0.46 refusal family
// added for DDL and stored-program simulation recall alongside valid accepted
// controls.
func TestWalkThroughMySQLRefusalFamilies(t *testing.T) {
	schema := &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name:   "parent",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}},
					},
				},
				{
					Name:   "plain",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
						{Name: "n", Type: "int"},
						{Name: "s", Type: "varchar(64)"},
						{Name: "g", Type: "point"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}},
						{Name: "idx_n", Type: "BTREE", Expressions: []string{"n"}},
					},
					CheckConstraints: []*metadata.CheckConstraintMetadata{
						{Name: "chk_n", Expression: "(`n` > 0)"},
					},
				},
				{
					Name:   "child",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
						{Name: "pid", Type: "int"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}},
						{Name: "fk_pid", Type: "BTREE", Expressions: []string{"pid"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{
						{Name: "fk_pid", Columns: []string{"pid"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}},
					},
				},
				{
					Name:   "parted",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
						{Name: "pcol", Type: "int"},
						{Name: "txt", Type: "text"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id", "pcol"}},
					},
					Partitions: []*metadata.TablePartitionMetadata{
						{Name: "p0", Type: metadata.TablePartitionMetadata_RANGE, Expression: "`pcol`"},
						{Name: "p1", Type: metadata.TablePartitionMetadata_RANGE, Expression: "`pcol`"},
					},
				},
				{
					Name:   "my",
					Engine: "MyISAM",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
					},
				},
				{
					Name:   "self_fk",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
						{Name: "pid", Type: "int"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}},
						{Name: "fk_self", Type: "BTREE", Expressions: []string{"pid"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{
						{Name: "fk_self", Columns: []string{"pid"}, ReferencedTable: "self_fk", ReferencedColumns: []string{"id"}},
					},
				},
				{
					Name:   "other_db_parent",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}},
					},
				},
				{
					Name:   "other_db_child",
					Engine: "InnoDB",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int"},
						{Name: "pid", Type: "int"},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}},
						{Name: "fk_other", Type: "BTREE", Expressions: []string{"pid"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{
						{Name: "fk_other", Columns: []string{"pid"}, ReferencedSchema: "Shop", ReferencedTable: "other_db_parent", ReferencedColumns: []string{"id"}},
					},
				},
			},
			Views:  []*metadata.ViewMetadata{{Name: "v"}},
			Events: []*metadata.EventMetadata{{Name: "ev"}},
		}},
	}
	target := review.Target{Schema: schema, ForeignKeyChecks: review.SwitchOn}

	// Refusals: every newly covered family produces a finding.
	refusals := []struct {
		name string
		sql  string
		want string
	}{
		// Partitioning management & clauses (1505, 1066, 1492, 1500, 1479, 1517, 1508, 1509, 1511, 1506, 1488, 1652, 1503, 3855, 1054).
		{"1505 add partition on non-partitioned", "ALTER TABLE plain ADD PARTITION (PARTITION p0 VALUES LESS THAN (10));", "Table `plain` is not partitioned."},
		{"1505 remove partitioning on non-partitioned", "ALTER TABLE plain REMOVE PARTITIONING;", "Table `plain` is not partitioned."},
		{"1066 exchange partition with self", "ALTER TABLE parted EXCHANGE PARTITION p0 WITH TABLE parted;", "Table `parted` cannot exchange a partition with itself."},
		{"1492 partition by range without defs", "ALTER TABLE plain PARTITION BY RANGE (id);", "Each partition must be defined for `RANGE` or `LIST` partitioning."},
		{"1492 add partition without defs on range", "ALTER TABLE parted ADD PARTITION PARTITIONS 1;", "Each partition must be defined for `RANGE` or `LIST` partitioning."},
		{"1500 subpartition on key", "ALTER TABLE plain PARTITION BY KEY () SUBPARTITION BY HASH (id);", "`SUBPARTITION BY` is only supported with `RANGE` or `LIST` partitioning."},
		{"1479 add partition missing values on range", "ALTER TABLE parted ADD PARTITION (PARTITION p2);", "Partition `p2` requires a `VALUES` definition."},
		{"1517 duplicate partition name", "ALTER TABLE parted ADD PARTITION (PARTITION p0 VALUES LESS THAN (300));", "Partition `p0` already exists in `parted`."},
		{"1508 drop all partitions", "ALTER TABLE parted DROP PARTITION p0, p1;", "Cannot remove all partitions from `parted`."},
		{"1509 coalesce on range", "ALTER TABLE parted COALESCE PARTITION 1;", "`COALESCE PARTITION` is only valid for `HASH` or `KEY` partitioning."},
		{"1511 reorganize without params on range", "ALTER TABLE parted REORGANIZE PARTITION;", "`REORGANIZE PARTITION` without parameters is only valid for `HASH` or `KEY` partitioning."},
		{"1506 partition table with fk", "ALTER TABLE child PARTITION BY KEY ();", "Partitioned table `child` cannot have foreign keys."},
		{"1506 add fk on partitioned table", "ALTER TABLE parted ADD FOREIGN KEY (id) REFERENCES parent (id);", "Partitioned table `parted` cannot define foreign keys."},
		{"1488 partition by key missing col", "ALTER TABLE plain PARTITION BY KEY (nosuch);", "Partition column `nosuch` doesn't exist in `plain`."},
		{"1652 duplicate partition key col", "ALTER TABLE plain PARTITION BY KEY (id, id);", "This statement lists partition column `id` twice."},
		{"1503 unique key missing partition col", "ALTER TABLE plain PARTITION BY HASH (n);", "Every unique key of `plain` must include all columns in its partitioning function."},
		{"3855 drop partition col", "ALTER TABLE parted DROP COLUMN pcol;", "Column `pcol` in `parted` is used by the table's partitioning function."},
		{"3855 rename partition col to self", "ALTER TABLE parted RENAME COLUMN pcol TO pcol;", "Column `pcol` in `parted` is used by the table's partitioning function."},
		{"1054 partition bound colref", "ALTER TABLE parted ADD PARTITION (PARTITION p2 VALUES LESS THAN (badcol));", "Column `badcol` doesn't exist in the partition function."},
		{"1054 partition hash missing col", "ALTER TABLE plain PARTITION BY HASH (nosuch);", "Column `nosuch` doesn't exist in the partition function."},

		// Index & FK columns / types / backing (1072, 1060, 3762, 1214, 1283, 1687, 1089, 1239, 1553, 1828).
		{"1072 create table missing key col", "CREATE TABLE x (a INT, KEY (missing));", "Key column `missing` doesn't exist in `x`."},
		{"1072 create index missing col", "CREATE INDEX idx_m ON plain (missing);", "Key column `missing` doesn't exist in `plain`."},
		{"1072 alter add key missing col", "ALTER TABLE plain ADD KEY (missing);", "Key column `missing` doesn't exist in `plain`."},
		{"1072 alter drop col then index it", "ALTER TABLE plain DROP COLUMN n, ADD KEY k_n (n);", "Key column `n` doesn't exist in `plain`."},
		{"1060 duplicate key col", "ALTER TABLE plain ADD KEY (n, n);", "This statement names column `n` twice."},
		{"3762 functional index on bare col", "ALTER TABLE plain ADD KEY ((n));", "Functional index key part cannot be a bare column reference."},
		{"1214 fulltext on partitioned table", "ALTER TABLE parted ADD FULLTEXT (txt);", "Partitioned table `parted` doesn't support `FULLTEXT` indexes."},
		{"1283 fulltext on int col", "ALTER TABLE plain ADD FULLTEXT (n);", "Column `n` cannot be part of a `FULLTEXT` index."},
		{"1687 spatial on int col", "ALTER TABLE plain ADD SPATIAL (n);", "Column `n` cannot be part of a `SPATIAL` index."},
		{"1687 rtree on int col", "ALTER TABLE plain ADD KEY USING RTREE (n);", "Column `n` cannot be part of a `SPATIAL` index."},
		{"1089 prefix on int col", "ALTER TABLE plain ADD KEY (n(1));", "Column `n` cannot use an index prefix length."},
		{"1239 fk missing ref cols", "ALTER TABLE plain ADD FOREIGN KEY (n) REFERENCES parent;", "Foreign key in `plain` must name matching referenced columns."},
		{"1553 drop index needed by fk", "DROP INDEX fk_pid ON child;", "Index `fk_pid` in `child` is needed by a foreign key constraint."},
		{"1553 alter drop index needed by fk", "ALTER TABLE child DROP INDEX fk_pid;", "Index `fk_pid` in `child` is needed by a foreign key constraint."},
		{"1828 drop col needed by fk", "ALTER TABLE child DROP COLUMN pid;", "Column `pid` in `child` is needed by a foreign key constraint."},

		// CHECK & general constraints (3812, 3821, 3940, 3941).
		{"3812 non-boolean check", "ALTER TABLE plain ADD CHECK (n);", "Check constraint in `plain` must be a boolean expression."},
		{"3812 non-boolean col check", "ALTER TABLE plain ADD COLUMN x INT CHECK (x);", "Check constraint on `x` must be a boolean expression."},
		{"3821 drop missing check", "ALTER TABLE plain DROP CHECK nosuch;", "Check constraint `nosuch` doesn't exist in `plain`."},
		{"3821 alter missing check", "ALTER TABLE plain ALTER CHECK nosuch NOT ENFORCED;", "Check constraint `nosuch` doesn't exist in `plain`."},
		{"3940 drop missing constraint", "ALTER TABLE plain DROP CONSTRAINT idx_n;", "Constraint `idx_n` doesn't exist in `plain`."},
		{"3940 alter missing constraint", "ALTER TABLE plain ALTER CONSTRAINT nosuch ENFORCED;", "Constraint `nosuch` doesn't exist in `plain`."},
		{"3941 alter non-check constraint", "ALTER TABLE child ALTER CONSTRAINT fk_pid ENFORCED;", "Constraint `fk_pid` in `child` is not a check constraint."},

		// Columns, defaults, table options, empty table, SRS, RENAME self (1439, 1221, 1294, 1067, 3107, 3767, 1075, 4028, 3979, 3184, 1478, 4025, 3708, 1050/1146).
		{"1439 int display width", "ALTER TABLE plain ADD COLUMN x INT(256);", "Column `x` has a display width or length past MySQL's limit."},
		{"1439 bit width", "ALTER TABLE plain ADD COLUMN x BIT(65);", "Column `x` has a display width or length past MySQL's limit."},
		{"1221 srid on non-geometry", "ALTER TABLE plain ADD COLUMN x INT SRID 1;", "Column `x` cannot declare `SRID` on a non-geometry type."},
		{"1221 standalone with validation", "ALTER TABLE plain WITH VALIDATION;", "`WITH/WITHOUT VALIDATION` requires a column modification or `EXCHANGE PARTITION`."},
		{"1294 on update on int", "ALTER TABLE plain ADD COLUMN x INT ON UPDATE CURRENT_TIMESTAMP;", "Column `x` cannot use `ON UPDATE` with type `int`."},
		{"1067 bare now default on int", "ALTER TABLE plain ADD COLUMN x INT DEFAULT CURRENT_TIMESTAMP;", "Column `x` has an invalid default value."},
		{"1067 non-numeric string default on int", "ALTER TABLE plain ALTER COLUMN n SET DEFAULT 'abc';", "Column `n` has an invalid default value."},
		{"3107 generated self ref", "ALTER TABLE plain ADD COLUMN x INT AS (x);", "Generated column `x` cannot refer to itself."},
		{"3767 expr default self ref", "ALTER TABLE plain ALTER COLUMN n SET DEFAULT (n);", "Default expression of column `n` cannot refer to itself."},
		{"1075 unindexed auto_increment", "ALTER TABLE plain ADD COLUMN x INT AUTO_INCREMENT;", "Column `x` is `AUTO_INCREMENT` and must be defined as a key."},
		{"4028 empty create table", "CREATE TABLE empty_tbl;", "Table `empty_tbl` must have at least one column."},
		{"3979 start transaction option", "ALTER TABLE plain START TRANSACTION;", "MySQL doesn't accept `START TRANSACTION` as a table option."},
		{"3184 invalid encryption", "ALTER TABLE plain ENCRYPTION = 'X';", "Invalid `ENCRYPTION` option value `X`."},
		{"1478 innodb invalid compression", "ALTER TABLE plain COMPRESSION = 'bad';", "InnoDB doesn't support `COMPRESSION` value `bad`."},
		{"1478 innodb row_format fixed", "ALTER TABLE plain ROW_FORMAT = FIXED;", "InnoDB doesn't support `ROW_FORMAT = FIXED`."},
		{"4025 innodb autoextend_size out of range", "ALTER TABLE plain AUTOEXTEND_SIZE = 1;", "`AUTOEXTEND_SIZE` value `1` is out of range for InnoDB."},
		{"3708 srs missing name", "CREATE SPATIAL REFERENCE SYSTEM 1;", "Spatial reference system requires a `NAME` attribute."},
		{"3708 srs missing definition", "CREATE SPATIAL REFERENCE SYSTEM 1 NAME 's';", "Spatial reference system requires a `DEFINITION` attribute."},
		{"1050 rename table to self", "RENAME TABLE plain TO plain;", "Table `plain` already exists."},
		{"1146 rename missing table to self", "RENAME TABLE missing TO missing;", "`missing` doesn't exist."},

		// Views / CTAS / Events (1054, 1551).
		{"1054 view from-less colref", "CREATE VIEW v2 AS SELECT badcol;", "Column `badcol` doesn't exist in the query."},
		{"1054 alter view from-less colref", "ALTER VIEW v AS SELECT badcol;", "Column `badcol` doesn't exist in the query."},
		{"1054 ctas from-less colref", "CREATE TABLE t2 AS SELECT badcol;", "Column `badcol` doesn't exist in the query."},
		{"1054 create event schedule colref", "CREATE EVENT ev2 ON SCHEDULE AT badcol DO COMMIT;", "Column `badcol` doesn't exist in the event schedule."},
		{"1054 alter event schedule colref", "ALTER EVENT ev ON SCHEDULE EVERY badcol DAY;", "Column `badcol` doesn't exist in the event schedule."},
		{"1551 alter event rename to self", "ALTER EVENT ev RENAME TO ev;", "Event `ev` cannot be renamed to the same name."},
	}
	for _, tc := range refusals {
		t.Run("refuse/"+tc.name, func(t *testing.T) {
			oneFinding(t, tc.sql, 0, tc.want, target)
		})
	}

	// Valid accepted controls: none of these may produce a finding.
	accepted := []struct {
		name string
		sql  string
	}{
		{"fulltext on varchar", "ALTER TABLE plain ADD FULLTEXT (s);"},
		{"spatial on geometry", "ALTER TABLE plain ADD SPATIAL (g);"},
		{"prefix on varchar", "ALTER TABLE plain ADD KEY (s(10));"},
		{"functional index on expr", "ALTER TABLE plain ADD KEY ((n + 1));"},
		{"boolean check", "ALTER TABLE plain ADD CHECK (n > 0);"},
		{"drop existing check", "ALTER TABLE plain DROP CHECK chk_n;"},
		{"alter existing check", "ALTER TABLE plain ALTER CHECK chk_n NOT ENFORCED;"},
		{"alter constraint on existing check", "ALTER TABLE plain ALTER CONSTRAINT chk_n NOT ENFORCED;"},
		{"drop index with fk dropped in same alter", "ALTER TABLE child DROP FOREIGN KEY fk_pid, DROP INDEX fk_pid;"},
		{"drop index with replacement index in same alter", "ALTER TABLE child DROP INDEX fk_pid, ADD INDEX fk_pid2 (pid);"},
		{"drop index with inline unique in same alter", "ALTER TABLE child DROP INDEX fk_pid, MODIFY COLUMN pid INT UNIQUE;"},
		{"drop col with fk dropped in same alter", "ALTER TABLE child DROP FOREIGN KEY fk_pid, DROP COLUMN pid;"},
		{"partition by hash covering pk", "ALTER TABLE plain PARTITION BY HASH (id);"},
		{"add valid range partition", "ALTER TABLE parted ADD PARTITION (PARTITION p2 VALUES LESS THAN (300));"},
		{"drop one of two partitions", "ALTER TABLE parted DROP PARTITION p0;"},
		{"remove partitioning then drop col", "ALTER TABLE parted REMOVE PARTITIONING;\nALTER TABLE parted DROP COLUMN pcol;"},
		{"myisam row_format fixed", "ALTER TABLE my ROW_FORMAT = FIXED;"},
		{"alter to myisam with row_format fixed", "ALTER TABLE plain ENGINE = MyISAM, ROW_FORMAT = FIXED;"},
		{"innodb valid compression", "ALTER TABLE plain COMPRESSION = 'ZLIB';"},
		{"innodb valid encryption", "ALTER TABLE plain ENCRYPTION = 'Y';"},
		{"innodb valid autoextend_size", "ALTER TABLE plain AUTOEXTEND_SIZE = 4M;"},
		{"innodb zero autoextend_size", "ALTER TABLE plain AUTOEXTEND_SIZE = 0;"},
		{"numeric string default on int", "ALTER TABLE plain ALTER COLUMN n SET DEFAULT '42';"},
		{"decimal string default on int", "CREATE TABLE x (a INT DEFAULT '1.5');"},
		{"scientific string default on int", "CREATE TABLE x (a INT DEFAULT '1e2');"},
		{"drop self fk and partition in same alter", "ALTER TABLE self_fk DROP FOREIGN KEY fk_self PARTITION BY HASH (id);"},
		{"partition table whose namesake in different-case db is referenced", "ALTER TABLE other_db_parent PARTITION BY HASH (id);"},
		{"drop pk col and repartition in same alter", "CREATE TABLE x (a INT PRIMARY KEY, b INT);\nALTER TABLE x DROP COLUMN a PARTITION BY HASH (b);"},
		{"rename fk col drop old index add replacement index", "ALTER TABLE child RENAME COLUMN pid TO newpid, DROP INDEX fk_pid, ADD INDEX fk_pid2 (newpid);"},
		{"set session default_storage_engine myisam with row_format fixed", "SET SESSION default_storage_engine = MyISAM;\nCREATE TABLE x (a INT) ROW_FORMAT = FIXED;"},
		{"set session innodb_strict_mode off with row_format fixed", "SET SESSION innodb_strict_mode = 0;\nALTER TABLE plain ROW_FORMAT = FIXED;"},
		{"expr default referencing other col", "ALTER TABLE plain ADD COLUMN x INT DEFAULT (n + 1);"},
		{"generated col referencing other col", "ALTER TABLE plain ADD COLUMN x INT AS (n + 1);"},
		{"auto_increment with key in same alter", "ALTER TABLE plain ADD COLUMN x INT AUTO_INCREMENT, ADD KEY (x);"},
		{"on update on timestamp", "ALTER TABLE plain ADD COLUMN ts TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;"},
		{"srid on point", "ALTER TABLE plain ADD COLUMN pt POINT SRID 4326;"},
		{"view with literal and builtins", "CREATE VIEW v2 AS SELECT 1, CURRENT_TIMESTAMP, TIMESTAMPDIFF(DAY, '2020-01-01', '2020-01-02');"},
		{"view with subquery", "CREATE VIEW v2 AS SELECT (SELECT id FROM plain LIMIT 1) AS x;"},
		{"event with interval schedule", "CREATE EVENT ev2 ON SCHEDULE EVERY 1 DAY STARTS TIMESTAMPADD(HOUR, 1, NOW()) DO COMMIT;"},
		{"event rename to new name", "ALTER EVENT ev RENAME TO ev2;"},
		{"valid srs", "CREATE SPATIAL REFERENCE SYSTEM 4326 NAME 'WGS 84' DEFINITION 'GEOGCS[\"WGS 84\"]';"},
		{"alter table rename col then partition by it", "ALTER TABLE plain RENAME COLUMN id TO id2 PARTITION BY HASH (id2);"},
		{"rename pk col then partition in next stmt", "CREATE TABLE x (a INT PRIMARY KEY);\nALTER TABLE x RENAME COLUMN a TO b;\nALTER TABLE x PARTITION BY HASH (b);"},
		{"swap cols then partition by renamed pk", "CREATE TABLE x (a INT PRIMARY KEY, b INT);\nALTER TABLE x RENAME COLUMN a TO b, RENAME COLUMN b TO a;\nALTER TABLE x PARTITION BY HASH (b);"},
		{"rename child drop child then partition parent", "RENAME TABLE child TO child2;\nDROP TABLE child2;\nALTER TABLE parent PARTITION BY HASH (id);"},
		{"rename child drop fk then partition parent", "RENAME TABLE child TO child2;\nALTER TABLE child2 DROP FOREIGN KEY fk_pid;\nALTER TABLE parent PARTITION BY HASH (id);"},
		{"other specs ignores snapshot partitioning", "ALTER TABLE plain ADD PARTITION (PARTITION p0 VALUES LESS THAN (10));"},
		{"other specs optional create table preserves partition uncertainty", "DROP TEMPORARY TABLE IF EXISTS x;\nCREATE TABLE IF NOT EXISTS x (a INT);\nALTER TABLE x ADD PARTITION (PARTITION p1 VALUES LESS THAN (20));"},
	}
	for _, tc := range accepted {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			tgt := target
			if strings.HasPrefix(tc.name, "other specs ") {
				tgt.OtherSpecs = true
			}
			noFinding(t, tc.sql, tgt)
		})
	}
}

// formLowering walks a form on the schema, as walktest.Coverage asks.
func formLowering(schema *metadata.DatabaseSchemaMetadata) func(t *testing.T, sql string) (n, lowered int, last bool) {
	return func(t *testing.T, sql string) (n, lowered int, last bool) {
		t.Helper()
		ranges := statementRanges(sql)
		stmts, syntax := parse(sql, ranges)
		if syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, sql)
		}
		w := newWalker(sql, stmts, review.Target{Schema: schema})
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
