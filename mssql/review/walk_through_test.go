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

// walkSchema is the synced shop database the unit tests walk on: dbo with
// tables t (keyed, with an index, constraints, and a default), r, and
// nokey, a view, a sequence, a procedure, and a function; and schema s.
func walkSchema() *metadata.DatabaseSchemaMetadata {
	cols := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	t := &metadata.TableMetadata{
		Name:    "t",
		Columns: cols("id", "a", "b"),
		Indexes: []*metadata.IndexMetadata{
			{Name: "pk_t", Primary: true, Unique: true, IsConstraint: true},
			{Name: "ix_a"},
		},
		ForeignKeys:      []*metadata.ForeignKeyMetadata{{Name: "fk_r", ReferencedTable: "r"}},
		CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck_b"}},
	}
	t.Columns[1].DefaultConstraintName = "df_a"
	return &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					t,
					{Name: "r", Columns: cols("id"), Indexes: []*metadata.IndexMetadata{{Name: "pk_r", Primary: true, Unique: true, IsConstraint: true}}},
					{Name: "nokey", Columns: cols("a")},
				},
				Views:      []*metadata.ViewMetadata{{Name: "v"}},
				Sequences:  []*metadata.SequenceMetadata{{Name: "sq"}},
				Synonyms:   []*metadata.SynonymMetadata{{Name: "syn", TargetSchema: "dbo", TargetName: "t"}},
				Procedures: []*metadata.ProcedureMetadata{{Name: "pr"}},
				Functions:  []*metadata.FunctionMetadata{{Name: "fn"}},
			},
			{Name: "s", Tables: []*metadata.TableMetadata{{Name: "st", Columns: cols("a")}}},
		},
	}
}

func walkTarget() review.Target {
	return review.Target{Schema: walkSchema(), DefaultSchema: "dbo"}
}

func walkFindings(t *testing.T, sql string, targets ...review.Target) []review.Finding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.WalkThrough}}, targets)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range result.Failures {
		t.Errorf("target %d could not be evaluated: %v", f.Target, f.Err)
	}
	for _, f := range result.Findings {
		if err := review.CheckCodes(f); err != nil {
			t.Error(err)
		}
		if f.Rule == review.Syntax {
			t.Fatalf("syntax: %s\n%s", f.Message, sql)
		}
	}
	return result.Findings
}

func TestWalkThroughFindings(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		statement int // -1 for no finding
		message   string
	}{
		// Creations.
		{"table exists", "CREATE TABLE dbo.t (a int);", 0, "Table `dbo.t` already exists."},
		{"table named like a constraint", "CREATE TABLE dbo.pk_t (a int);", 0, "The name `dbo.pk_t` is already taken."},
		{"case variant is unknown", "CREATE TABLE dbo.T (a int);", -1, ""},
		{"unqualified with the default schema", "CREATE TABLE t (a int);", 0, "Table `t` already exists."},
		{"missing schema", "CREATE TABLE nosuch.x (a int);", 0, "Schema `nosuch` doesn't exist."},
		{"other schema is free", "CREATE TABLE s.t (a int);", -1, ""},
		{"created twice", "CREATE TABLE dbo.x (a int);\nCREATE TABLE dbo.x (a int);", 1, "Table `dbo.x` already exists."},
		{"column twice", "CREATE TABLE dbo.x (a int, a int);", 0, "This statement declares column `a` twice."},
		{"two primary keys", "CREATE TABLE dbo.x (a int PRIMARY KEY, b int PRIMARY KEY);", 0, "This statement gives `dbo.x` two primary keys."},
		{"constraint name taken", "CREATE TABLE dbo.x (a int CONSTRAINT ck_b CHECK (a > 0));", 0, "Check constraint `ck_b` already exists."},
		{"constraint named like the table", "CREATE TABLE dbo.x (a int CONSTRAINT x CHECK (a > 0));", 0, "This statement names `x` twice."},
		{"constraint name twice", "CREATE TABLE dbo.x (a int CONSTRAINT c1 CHECK (a > 0), b int CONSTRAINT c1 CHECK (b > 0));", 0, "This statement declares constraint `c1` twice."},
		{"foreign key to a missing table", "CREATE TABLE dbo.x (a int REFERENCES dbo.nosuch (id));", 0, "Table `dbo.nosuch` doesn't exist."},
		{"foreign key to a view", "CREATE TABLE dbo.x (a int REFERENCES dbo.v (id));", 0, "`dbo.v` is a view, not a table."},
		{"foreign key to itself", "CREATE TABLE dbo.x (id int PRIMARY KEY, p int REFERENCES dbo.x (id));", -1, ""},
		{"view over a table", "CREATE VIEW dbo.t AS SELECT 1 AS a;", 0, "The name `dbo.t` is already taken."},
		{"sequence over a procedure", "CREATE SEQUENCE dbo.pr;", 0, "The name `dbo.pr` is already taken."},
		{"temporary table twice", "CREATE TABLE #x (a int);\nCREATE TABLE #x (a int);", 1, "Table `#x` already exists."},
		{"select into an existing table", "SELECT 1 AS a INTO dbo.t;", 0, "Table `dbo.t` already exists."},
		// Drops.
		{"drop missing table", "DROP TABLE dbo.nosuch;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"drop a view as a table", "DROP TABLE dbo.v;", 0, "`dbo.v` is a view, not a table."},
		{"drop if exists of a view", "DROP TABLE IF EXISTS dbo.v;", 0, "`dbo.v` is a view, not a table."},
		{"drop if exists of a procedure", "DROP TABLE IF EXISTS dbo.pr;", 0, "`dbo.pr` is a procedure, not a table."},
		{"drop if exists missing", "DROP TABLE IF EXISTS dbo.nosuch;", -1, ""},
		{"drop twice", "DROP TABLE dbo.t;\nDROP TABLE dbo.t;", 1, "Table `dbo.t` doesn't exist."},
		{"drop twice in one", "DROP TABLE dbo.r, dbo.r;", 0, "This statement drops table `dbo.r` twice."},
		{"drop case variant twice", "DROP TABLE dbo.t;\nDROP TABLE dbo.T;", 1, "Table `dbo.T` doesn't exist."},
		{"drop then create", "DROP TABLE dbo.t;\nCREATE TABLE dbo.t (a int);\nCREATE TABLE dbo.pk_t (a int);", -1, ""},
		{"drop then use", "DROP TABLE dbo.nokey;\nINSERT INTO dbo.nokey VALUES (1);", 1, "`dbo.nokey` doesn't exist."},
		// DML resolves through tables, views, sequences, and synonyms, which
		// the sync lists.
		{"DML on a missing name", "INSERT INTO dbo.nosuch VALUES (1);", 0, "`dbo.nosuch` doesn't exist."},
		{"unqualified DML on a missing name", "UPDATE nosuch SET a = 1;", 0, "`nosuch` doesn't exist."},
		{"DML through a synonym", "INSERT INTO dbo.syn (a) VALUES (1);", -1, ""},
		{"DML on a view", "INSERT INTO dbo.v VALUES (1);", -1, ""},
		{"DML after a dropped synonym", "DROP SYNONYM dbo.syn;\nDELETE FROM dbo.syn;", 1, "`dbo.syn` doesn't exist."},
		{"DML through a synonym the change makes", "CREATE SYNONYM dbo.s2 FOR dbo.t;\nINSERT INTO dbo.s2 (a) VALUES (1);", -1, ""},
		{"DML through a renamed synonym", "EXEC sp_rename 'dbo.syn', 'syn2';\nINSERT INTO dbo.syn (a) VALUES (1);", 1, "`dbo.syn` doesn't exist."},
		{"DML through a synonym ends the walk", "INSERT INTO dbo.syn (a) VALUES (1);\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"a synonym takes its name", "CREATE TABLE dbo.syn (a int);", 0, "The name `dbo.syn` is already taken."},
		{"DROP SYNONYM of a missing name", "DROP SYNONYM dbo.nosuch;", 0, "Synonym `dbo.nosuch` doesn't exist."},
		{"DROP VIEW of a synonym", "DROP VIEW dbo.syn;", 0, "`dbo.syn` is a synonym, not a view."},
		{"drop missing view", "DROP VIEW dbo.nosuch;", 0, "View `dbo.nosuch` doesn't exist."},
		{"drop a table as a view", "DROP VIEW dbo.t;", 0, "`dbo.t` is a table, not a view."},
		{"drop procedure unsynced may exist", "DROP PROCEDURE dbo.nosuch;", -1, ""},
		{"drop procedure of a table", "DROP PROCEDURE dbo.t;", 0, "`dbo.t` is a table, not a procedure."},
		{"drop function of a procedure", "DROP FUNCTION dbo.pr;", 0, "`dbo.pr` is a procedure, not a function."},
		{"drop procedure twice", "DROP PROCEDURE dbo.pr;\nDROP PROCEDURE dbo.pr;", 1, "Procedure `dbo.pr` doesn't exist."},
		{"drop sequence missing", "DROP SEQUENCE dbo.nosuch;", 0, "Sequence `dbo.nosuch` doesn't exist."},
		{"drop schema missing", "DROP SCHEMA nosuch;", 0, "Schema `nosuch` doesn't exist."},
		{"drop schema with a table", "DROP SCHEMA s;", 0, "Schema `s` still holds objects."},
		{"drop schema emptied", "DROP TABLE s.st;\nDROP SCHEMA s;", -1, ""},
		{"create schema twice", "CREATE SCHEMA s;", 0, "Schema `s` already exists."},
		// ALTER TABLE.
		{"alter missing table", "ALTER TABLE dbo.nosuch ADD x int;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"alter a view", "ALTER TABLE dbo.v ADD x int;", 0, "`dbo.v` is a view, not a table."},
		{"add existing column", "ALTER TABLE dbo.t ADD a int;", 0, "Column `a` already exists in `dbo.t`."},
		{"add column twice", "ALTER TABLE dbo.t ADD x int, x int;", 0, "This statement adds column `x` twice."},
		{"drop missing column", "ALTER TABLE dbo.t DROP COLUMN nosuch;", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"drop missing column if exists", "ALTER TABLE dbo.t DROP COLUMN IF EXISTS nosuch;", -1, ""},
		{"drop column twice", "ALTER TABLE dbo.t DROP COLUMN b, COLUMN b;", 0, "This statement drops column `b` twice."},
		{"alter missing column", "ALTER TABLE dbo.t ALTER COLUMN nosuch int;", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"second primary key", "ALTER TABLE dbo.t ADD CONSTRAINT pk2 PRIMARY KEY (id);", 0, "`dbo.t` already has a primary key."},
		{"primary key on a keyless table", "ALTER TABLE dbo.nokey ADD CONSTRAINT pk2 PRIMARY KEY (a);", -1, ""},
		{"constraint name taken in ALTER", "ALTER TABLE dbo.nokey ADD CONSTRAINT pk_t CHECK (a > 0);", 0, "The name `pk_t` is already taken."},
		{"unique named like an index", "ALTER TABLE dbo.t ADD CONSTRAINT ix_a UNIQUE (b);", 0, "Index `ix_a` already exists in `dbo.t`."},
		{"drop missing constraint", "ALTER TABLE dbo.t DROP CONSTRAINT nosuch;", 0, "Constraint `nosuch` doesn't exist in `dbo.t`."},
		{"drop another table's constraint", "ALTER TABLE dbo.r DROP CONSTRAINT pk_t;", 0, "Constraint `pk_t` doesn't exist in `dbo.r`."},
		{"drop a default", "ALTER TABLE dbo.t DROP CONSTRAINT df_a;", -1, ""},
		{"drop then drop the constraint", "ALTER TABLE dbo.t DROP CONSTRAINT ck_b;\nALTER TABLE dbo.t DROP CONSTRAINT ck_b;", 1, "Constraint `ck_b` doesn't exist in `dbo.t`."},
		{"primary key dropped then added", "ALTER TABLE dbo.t DROP CONSTRAINT pk_t;\nALTER TABLE dbo.t ADD CONSTRAINT pk_t2 PRIMARY KEY (id);", -1, ""},
		{"nocheck missing constraint", "ALTER TABLE dbo.t NOCHECK CONSTRAINT nosuch;", 0, "Constraint `nosuch` doesn't exist in `dbo.t`."},
		{"column added then dropped twice", "ALTER TABLE dbo.t ADD x int;\nALTER TABLE dbo.t DROP COLUMN x;\nALTER TABLE dbo.t DROP COLUMN x;", 2, "Column `x` doesn't exist in `dbo.t`."},
		{"add index", "ALTER TABLE dbo.t ADD INDEX ix_new (b);", -1, ""},
		{"add index named like one", "ALTER TABLE dbo.t ADD INDEX ix_a (b);", 0, "Index `ix_a` already exists in `dbo.t`."},
		{"add index keeps the walk going", "ALTER TABLE dbo.t ADD INDEX ix_new (b);\nDROP TABLE dbo.nosuch;", 1, "Table `dbo.nosuch` doesn't exist."},
		{"add index twice", "ALTER TABLE dbo.t ADD INDEX ix, INDEX ix;", 0, "This statement adds index `ix` twice."},
		{"add an index and a constraint of one name", "ALTER TABLE dbo.t ADD CONSTRAINT c1 CHECK (b > 0), INDEX c1 (b);", 0, "This statement adds `c1` twice."},
		// Indexes.
		{"index exists", "CREATE INDEX ix_a ON dbo.t (b);", 0, "Index `ix_a` already exists in `dbo.t`."},
		{"index named like the key", "CREATE INDEX pk_t ON dbo.t (b);", 0, "Index `pk_t` already exists in `dbo.t`."},
		{"index on another table is free", "CREATE INDEX pk_t ON dbo.r (id);", -1, ""},
		{"index on a missing table", "CREATE INDEX i ON dbo.nosuch (a);", 0, "Table `dbo.nosuch` doesn't exist."},
		{"drop missing index", "DROP INDEX nosuch ON dbo.t;", 0, "Index `nosuch` doesn't exist in `dbo.t`."},
		{"drop a key's index", "DROP INDEX pk_t ON dbo.t;", 0, "Index `pk_t` enforces constraint `pk_t`."},
		{"drop a key's index if exists", "DROP INDEX IF EXISTS pk_t ON dbo.t;", 0, "Index `pk_t` enforces constraint `pk_t`."},
		{"old drop index form", "DROP INDEX dbo.t.ix_a;\nDROP INDEX dbo.t.ix_a;", 1, "Index `ix_a` doesn't exist in `dbo.t`."},
		{"drop_existing of a missing index", "CREATE INDEX nosuch ON dbo.t (a) WITH (DROP_EXISTING = ON);", 0, "Index `nosuch` doesn't exist in `dbo.t`."},
		{"alter missing index", "ALTER INDEX nosuch ON dbo.t REBUILD;", 0, "Index `nosuch` doesn't exist in `dbo.t`."},
		{"truncate a view", "TRUNCATE TABLE dbo.v;", 0, "`dbo.v` is a view, not a table."},
		{"truncate missing", "TRUNCATE TABLE dbo.nosuch;", 0, "Table `dbo.nosuch` doesn't exist."},
		// Modules.
		{"alter missing view", "ALTER VIEW dbo.nosuch AS SELECT 1 AS a;", 0, "View `dbo.nosuch` doesn't exist."},
		{"alter a table as a view", "ALTER VIEW dbo.t AS SELECT 1 AS a;", 0, "`dbo.t` is a table, not a view."},
		{"create or alter a view over a table", "CREATE OR ALTER VIEW dbo.t AS SELECT 1 AS a;", 0, "`dbo.t` is a table, not a view."},
		{"create or alter a procedure over a function", "CREATE OR ALTER PROCEDURE dbo.fn AS SELECT 1;", 0, "`dbo.fn` is a function, not a procedure."},
		{"alter procedure unsynced", "ALTER PROCEDURE dbo.nosuch AS SELECT 1;", -1, ""},
		{"view not first", "CREATE TABLE dbo.x (a int);\nCREATE VIEW dbo.w AS SELECT 1 AS a;", 1, "`CREATE VIEW` must be the first statement of its batch."},
		{"view after GO", "CREATE TABLE dbo.x (a int);\nGO\nCREATE VIEW dbo.w AS SELECT 1 AS a;", -1, ""},
		{"procedure body swallows", "CREATE PROCEDURE dbo.p2 AS SELECT 1;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"procedure body ends at GO", "CREATE PROCEDURE dbo.p2 AS SELECT 1;\nCREATE TABLE dbo.x (a int);\nGO\nCREATE TABLE dbo.t (a int);", 3, "Table `dbo.t` already exists."},
		{"trigger on a missing table", "CREATE TRIGGER dbo.tr ON dbo.nosuch AFTER INSERT AS SELECT 1;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"trigger in another schema", "CREATE TRIGGER s.tr ON dbo.t AFTER INSERT AS SELECT 1;", 0, "Trigger `s.tr` must be in the schema of `dbo.t`."},
		{"alter sequence of a table", "ALTER SEQUENCE dbo.t RESTART;", 0, "`dbo.t` is a table, not a sequence."},
		{"type twice", "CREATE TYPE dbo.ty FROM int;\nCREATE TYPE dbo.ty FROM int;", 1, "Type `dbo.ty` already exists."},
		{"type dropped twice", "CREATE TYPE dbo.ty FROM int;\nDROP TYPE dbo.ty;\nDROP TYPE dbo.ty;", 2, "Type `dbo.ty` doesn't exist."},
		// sp_rename.
		{"rename to a taken name", "EXEC sp_rename 'dbo.t', 'r';", 0, "Table `r` already exists."},
		{"rename a dropped table", "DROP TABLE dbo.nokey;\nEXEC sp_rename 'dbo.nokey', 'x';", 1, "`dbo.nokey` doesn't exist."},
		{"rename then use", "EXEC sp_rename 'dbo.nokey', 'n2';\nALTER TABLE dbo.n2 ADD a int;", 1, "Column `a` already exists in `dbo.n2`."},
		{"rename frees the name", "EXEC sp_rename 'dbo.nokey', 'n2';\nCREATE TABLE dbo.nokey (a int);", -1, ""},
		{"rename column", "EXEC sp_rename 'dbo.t.a', 'z', 'COLUMN';\nALTER TABLE dbo.t DROP COLUMN a;", 1, "Column `a` doesn't exist in `dbo.t`."},
		{"rename column to a taken name", "EXEC sp_rename 'dbo.t.a', 'b', 'COLUMN';", 0, "Column `b` already exists in `dbo.t`."},
		{"rename missing column", "EXEC sp_rename 'dbo.t.nosuch', 'z', 'COLUMN';", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"rename index", "EXEC sp_rename 'dbo.t.ix_a', 'ix_b', 'INDEX';\nDROP INDEX ix_a ON dbo.t;", 1, "Index `ix_a` doesn't exist in `dbo.t`."},
		// Names too long, cut in the message.
		{"a name too long", "CREATE TABLE dbo." + strings.Repeat("a", 129) + " (a int);", 0, "The name `" + strings.Repeat("a", 119) + "…` is longer than SQL Server's 128-character limit."},
		{"a local temporary table's name too long", "CREATE TABLE #" + strings.Repeat("l", 117) + " (a int);", 0, "The name `#" + strings.Repeat("l", 117) + "` is longer than SQL Server's 116-character limit."},
		// Unknown by design.
		{"no default schema", "CREATE TABLE t (a int);", -1, ""},
		{"exec ends the walk", "EXEC dbo.pr;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"if makes effects optional", "IF OBJECT_ID('dbo.t') IS NOT NULL DROP TABLE dbo.t;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"if keeps the walk going", "IF 1 = 1 DROP TABLE dbo.nokey;\nCREATE TABLE dbo.r (a int);", 1, "Table `dbo.r` already exists."},
		{"try ends the walk", "BEGIN TRY DROP TABLE dbo.nosuch END TRY BEGIN CATCH END CATCH;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"GO forgets temporary tables", "CREATE TABLE #x (a int);\nGO\nCREATE TABLE #x (a int);", -1, ""},
		{"GO not alone ends the walk", "SELECT 1 GO\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"GO twice ends the walk", "CREATE TABLE dbo.x (a int)\nGO 2\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"USE elsewhere", "USE master;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"USE back", "USE master;\nUSE shop;\nCREATE TABLE dbo.t (a int);", 2, "Table `dbo.t` already exists."},
		{"three-part target", "CREATE TABLE shop.dbo.t (a int);", 0, "Table `shop.dbo.t` already exists."},
		{"three-part elsewhere", "CREATE TABLE other.dbo.t (a int);", -1, ""},
		{"noexec ends the walk", "SET NOEXEC ON;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"a schema elsewhere", "USE master;\nGO\nCREATE SCHEMA s;\nGO\nDROP SCHEMA nosuch;", -1, ""},
		{"a schema in an unknown database", "USE SHOP;\nGO\nCREATE SCHEMA s;\nGO\nUSE shop;\nCREATE TABLE s.st (a int);", -1, ""},
		{"a default schema that may not exist", "CREATE TABLE t (a int);", -1, ""},
		{"a missing schema holds no table", "DROP TABLE nosuch.t;", 0, "Table `nosuch.t` doesn't exist."},
		// Codex round 1 on #514.
		{"a marked name may be the table another spelling means", "ALTER TABLE dbo.CAFE ADD a int;", -1, ""},
		{"the exact spelling means its own table", "ALTER TABLE dbo.cafe ADD a int;", 0, "Column `a` already exists in `dbo.cafe`."},
		{"a dropped spelling may leave a marked table", "DROP TABLE dbo.cafe;\nDROP TABLE dbo.CAFE;", -1, ""},
		{"a created marked name", "CREATE TABLE dbo.[NOKÉY] (b int);\nALTER TABLE dbo.NOKEY ADD a int;", -1, ""},
		{"a check and an index of one name", "CREATE TABLE dbo.x (a int, CONSTRAINT ix CHECK (a > 0), INDEX ix NONCLUSTERED (a));", 0, "This statement declares `ix` twice."},
		{"a unique constraint and an index of one name", "CREATE TABLE dbo.x (a int CONSTRAINT u1 UNIQUE, INDEX u1 (a));", 0, "This statement declares `u1` twice."},
		{"DROP CONSTRAINT IF EXISTS of another object's name", "ALTER TABLE dbo.t DROP CONSTRAINT IF EXISTS nokey;\nINSERT INTO dbo.nokey VALUES (1);\nCREATE TABLE dbo.nokey (a int);", 2, "Table `dbo.nokey` already exists."},
		// Codex round 2 on #514.
		{"a marked schema may be the one another spelling means", "CREATE TABLE CAFE.t (a int);", -1, ""},
		{"the exact schema spelling means its own schema", "CREATE TABLE cafe.t (a int);", 0, "Table `cafe.t` already exists."},
		{"sp_rename cuts a long name", "EXEC sp_rename 'dbo.nokey', '" + strings.Repeat("n", 129) + "';", -1, ""},
		{"a primary XML index takes its secondary ones", "DROP INDEX px ON dbo.t;\nCREATE INDEX sx ON dbo.t (a);", -1, ""},
		{"a secondary XML index", "CREATE INDEX sx ON dbo.t (a);", 0, "Index `sx` already exists in `dbo.t`."},
		{"a self-reference of another spelling", "CREATE TABLE dbo.x (id int NOT NULL PRIMARY KEY, p int REFERENCES dbo.X (id));", -1, ""},
		{"an alias of a marked spelling", "DROP TABLE dbo.nokey;\nUPDATE a SET a = 1 FROM dbo.r AS [á];", -1, ""},
		{"a marked arrival leaves its space's absences unknown", "CREATE TABLE dbo.[é] (a int UNIQUE);\nDROP TABLE dbo.cafe;", -1, ""},
		// Codex round 3 on #514.
		{"an XML index's name is no object", "CREATE TABLE dbo.px (a int);", -1, ""},
		{"a dropped schema beside a marked one", "DROP SCHEMA cafe;\nALTER TABLE CAFE.t ADD z int;", -1, ""},
		{"an unplaced drop may reach a marked name", "DROP TABLE CAFE;\nCREATE TABLE dbo.[CAFÉ] (a int);", -1, ""},
		{"dotted names are distinct", "DROP TABLE dbo.[x.y], [dbo.x].y;", -1, ""},
		{"DROP_EXISTING keeps a key's constraint", "CREATE UNIQUE CLUSTERED INDEX pk_t ON dbo.t (id) WITH (DROP_EXISTING = ON);\nDROP INDEX pk_t ON dbo.t;", 1, "Index `pk_t` enforces constraint `pk_t`."},
		// Another spelling means the class's one table or nothing, and
		// fails either way.
		{"another spelling's primary key", "ALTER TABLE dbo.T ADD CONSTRAINT pk2 PRIMARY KEY (id);", 0, "`dbo.T` already has a primary key."},
		// Codex round 5 on #514.
		{"a database in front of a temporary table", "CREATE TABLE #x (a int, b int);\nALTER TABLE tempdb..#x DROP COLUMN a;\nALTER TABLE #x ADD a int;", -1, ""},
		{"a temporary table in any database", "CREATE TABLE #x (a int);\nALTER TABLE nosuch..#x ADD a int;", 1, "Column `a` already exists in `nosuch..#x`."},
		{"a linked server's temporary table", "CREATE TABLE #x (a int, b int, INDEX i (b));\nALTER TABLE srv.tempdb..#x DROP COLUMN a;\nALTER TABLE #x ADD a int;\nCREATE INDEX i ON #x (b);\nDROP INDEX j ON #x;", -1, ""},
		// Codex round 6 on #514; each ends with a statement the walk still
		// judges.
		{"a global temporary table's long name", "CREATE TABLE ##" + strings.Repeat("g", 126) + " (x int);\nDROP TABLE dbo.nosuch;", 1, "Table `dbo.nosuch` doesn't exist."},
		{"a marked name of another case", "ALTER TABLE dbo.[CAFÉ] DROP COLUMN a;\nGO\nALTER TABLE dbo.[Café] ADD a int;\nCREATE INDEX i9 ON dbo.[CAFÉ] (b);\nDROP INDEX i9 ON dbo.[Café];\nDROP TABLE s.nosuch;", 5, "Table `s.nosuch` doesn't exist."},
		{"an index on a marked name of another case", "CREATE INDEX i9 ON dbo.[CAFÉ] (b);\nDROP INDEX i9 ON dbo.[Café];\nDROP TABLE s.nosuch;", 2, "Table `s.nosuch` doesn't exist."},
		{"another spelling beside a marked name", "ALTER TABLE dbo.NOKEY DROP COLUMN a;\nALTER TABLE dbo.nokey ADD a int;\nDROP TABLE s.nosuch;", 2, "Table `s.nosuch` doesn't exist."},
		// Codex round 7 on #514.
		{"a constraint dropped by an accent-equal spelling", "ALTER TABLE dbo.tc DROP CONSTRAINT [pk_café];\nCREATE INDEX pk_cafe ON dbo.tc (id);\nDROP TABLE s.nosuch;", 2, "Table `s.nosuch` doesn't exist."},
		{"a dropped ledger table frees its history table and view", "DROP TABLE dbo.lt;\nCREATE TABLE s.lt_h (a int);\nGO\nCREATE VIEW s.lt_v AS SELECT 1 AS a;\nGO\nDROP TABLE s.nosuch;", 5, "Table `s.nosuch` doesn't exist."},
		{"a table of one NOT NULL bigint is no ledger table", "DROP TABLE dbo.lt;\nCREATE TABLE s.lt_h (a int);", 1, "Table `s.lt_h` already exists."},
		{"a numbered procedure joins its group", "CREATE PROCEDURE dbo.p9;1 AS SELECT 1;\nGO\nCREATE PROCEDURE dbo.p9;2 AS SELECT 2;\nGO\nDROP TABLE s.nosuch;", 4, "Table `s.nosuch` doesn't exist."},
		{"a numbered procedure of a table's name", "CREATE PROCEDURE dbo.t;2 AS SELECT 2;", 0, "`dbo.t` is a table, not a procedure."},
		{"an unqualified type in the default schema", "CREATE TYPE s.ty FROM int;\nDROP TABLE s.ty;\nDROP TYPE ty;\nCREATE TYPE s.ty FROM int;\nDROP TABLE dbo.nosuch;", 4, "Table `dbo.nosuch` doesn't exist."},
		// The review of 2026-10-28.
		{"a trigger in another schema on a missing table", "CREATE TRIGGER s.tr ON nosuch AFTER INSERT AS SELECT 1;", 0, "Table `nosuch` doesn't exist."},
		{"an INSERT fires its table's triggers alone", "INSERT INTO dbo.r VALUES (1);\nDROP TABLE dbo.nosuch;", 1, "Table `dbo.nosuch` doesn't exist."},
		{"a DELETE may cascade to a trigger the sync does not list", "DELETE FROM dbo.r;\nDROP TABLE dbo.nosuch;", -1, ""},
		{"the last DROP_EXISTING counts", "CREATE INDEX nosuch ON dbo.t (a) WITH (DROP_EXISTING = ON, DROP_EXISTING = OFF);\nCREATE INDEX ix_a ON dbo.t (a) WITH (DROP_EXISTING = ON, DROP_EXISTING = OFF);", 1, "Index `ix_a` already exists in `dbo.t`."},
		{"a dropped table's capture instance", "DROP TABLE dbo.nokey;\nGO\nCREATE TABLE cdc.dbo_nokey_CT (a int);", -1, ""},
		{"a change table without change data capture", "DROP TABLE dbo.nokey;\nGO\nCREATE TABLE cdc.dbo_nokey_CT (a int);", 2, "Table `cdc.dbo_nokey_CT` already exists."},
		// Walk simulation recall fixes.
		{"alter column add sparse on missing column", "ALTER TABLE dbo.t ALTER COLUMN nosuch ADD SPARSE;", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"alter column drop rowguidcol on existing column", "ALTER TABLE dbo.t ALTER COLUMN a DROP ROWGUIDCOL;", -1, ""},
		{"create primary xml index on missing table", "CREATE PRIMARY XML INDEX px1 ON dbo.nosuch (a);", 0, "Table `dbo.nosuch` doesn't exist."},
		{"create primary xml index on view", "CREATE PRIMARY XML INDEX px1 ON dbo.v (a);", 0, "`dbo.v` is a view, not a table."},
		{"create primary xml index name taken", "CREATE PRIMARY XML INDEX ix_a ON dbo.t (a);", 0, "Index `ix_a` already exists in `dbo.t`."},
		{"create selective xml index on view", "CREATE SELECTIVE XML INDEX sxi1 ON dbo.v (a) FOR (p = '/a');", 0, "`dbo.v` is a view, not a table."},
		{"create spatial index on missing table", "CREATE SPATIAL INDEX sp1 ON dbo.nosuch (a);", 0, "Table `dbo.nosuch` doesn't exist."},
		{"create spatial index on view", "CREATE SPATIAL INDEX sp1 ON dbo.v (a);", 0, "`dbo.v` is a view, not a table."},
		{"create spatial index name taken", "CREATE SPATIAL INDEX ix_a ON dbo.t (a);", 0, "Index `ix_a` already exists in `dbo.t`."},
		{"create spatial index then duplicate", "CREATE SPATIAL INDEX sp1 ON dbo.t (a);\nCREATE SPATIAL INDEX sp1 ON dbo.t (a);", 1, "Index `sp1` already exists in `dbo.t`."},
		{"disable trigger on missing table", "DISABLE TRIGGER ALL ON dbo.nosuch;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"disable trigger on synonym", "DISABLE TRIGGER ALL ON dbo.syn;", 0, "`dbo.syn` is a synonym, not a table or view."},
		{"enable trigger on view", "ENABLE TRIGGER ALL ON dbo.v;", -1, ""},
		{"insert into sequence", "INSERT INTO dbo.sq VALUES (1);", 0, "`dbo.sq` is a sequence, not a table, view, or synonym."},
		{"delete from sequence", "DELETE FROM dbo.sq;", 0, "`dbo.sq` is a sequence, not a table, view, or synonym."},
		{"update sequence", "UPDATE dbo.sq SET a = 1;", 0, "`dbo.sq` is a sequence, not a table, view, or synonym."},
		{"insert exec on missing table", "INSERT INTO dbo.nosuch EXEC dbo.pr;", 0, "`dbo.nosuch` doesn't exist."},
		{"insert exec ends walk after target check", "INSERT INTO dbo.nokey EXEC dbo.pr;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"alter table switch on missing table", "ALTER TABLE dbo.nosuch SWITCH TO dbo.r;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"alter table switch on view", "ALTER TABLE dbo.v SWITCH TO dbo.r;", 0, "`dbo.v` is a view, not a table."},
		{"alter table system_versioning on existing table ends walk", "ALTER TABLE dbo.t SET (SYSTEM_VERSIONING = OFF);\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"create edge table over existing name", "CREATE TABLE dbo.t AS EDGE;", 0, "Table `dbo.t` already exists."},
		{"create edge table in missing schema", "CREATE TABLE nosuch.e AS EDGE;", 0, "Schema `nosuch` doesn't exist."},
		{"create edge table ends walk", "CREATE TABLE dbo.e AS EDGE;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"create schema with elements over existing schema", "CREATE SCHEMA s CREATE VIEW s.v1 AS SELECT 1 AS a;", 0, "Schema `s` already exists."},
		{"create schema with elements ends walk", "CREATE SCHEMA s2 CREATE VIEW s2.v1 AS SELECT 1 AS a;\nCREATE TABLE dbo.t (a int);", -1, ""},
		// Statistics checks.
		{"update statistics missing table", "UPDATE STATISTICS dbo.nosuch;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"update statistics synonym", "UPDATE STATISTICS dbo.syn;", 0, "`dbo.syn` is a synonym, not a table or view."},
		{"update statistics existing table", "UPDATE STATISTICS dbo.t;", -1, ""},
		{"update statistics view allowed", "UPDATE STATISTICS dbo.v;", -1, ""},
		{"create statistics missing table", "CREATE STATISTICS st1 ON dbo.nosuch (a);", 0, "Table `dbo.nosuch` doesn't exist."},
		{"create statistics synonym", "CREATE STATISTICS st1 ON dbo.syn (a);", 0, "`dbo.syn` is a synonym, not a table or view."},
		{"create statistics duplicate column", "CREATE STATISTICS st1 ON dbo.t (a, a);", 0, "This statement lists column `a` twice."},
		{"create statistics missing column", "CREATE STATISTICS st1 ON dbo.t (nosuch);", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"create statistics collides with index", "CREATE STATISTICS ix_a ON dbo.t (a);", 0, "Index `ix_a` already exists in `dbo.t`."},
		{"create statistics valid", "CREATE STATISTICS st1 ON dbo.t (a);", -1, ""},
		// ALTER SCHEMA.
		{"alter schema missing target", "ALTER SCHEMA nosuch TRANSFER dbo.t;", 0, "Schema `nosuch` doesn't exist."},
		{"alter schema existing target", "ALTER SCHEMA s TRANSFER dbo.t;", -1, ""},
		// Legacy defaults and rules.
		{"create default over table", "CREATE DEFAULT dbo.t AS 1;", 0, "The name `dbo.t` is already taken."},
		{"create default missing schema", "CREATE DEFAULT nosuch.d1 AS 1;", 0, "Schema `nosuch` doesn't exist."},
		{"create default valid", "CREATE DEFAULT dbo.d1 AS 1;", -1, ""},
		{"drop default of table", "DROP DEFAULT dbo.t;", 0, "`dbo.t` is a table, not a default."},
		{"drop default if exists of view", "DROP DEFAULT IF EXISTS dbo.v;", 0, "`dbo.v` is a view, not a default."},
		{"drop default unsynced may exist", "DROP DEFAULT dbo.nosuch;", -1, ""},
		{"drop rule of table", "DROP RULE dbo.t;", 0, "`dbo.t` is a table, not a rule."},
		{"drop rule if exists of procedure", "DROP RULE IF EXISTS dbo.pr;", 0, "`dbo.pr` is a procedure, not a rule."},
		{"drop rule unsynced may exist", "DROP RULE dbo.nosuch;", -1, ""},
		// CREATE TYPE base type validation.
		{"create type from non-system type", "CREATE TYPE dbo.ty FROM dbo.custom_t;", 0, "Base type `dbo.custom_t` is not a system data type."},
		{"create type from unknown base type", "CREATE TYPE dbo.ty FROM nosuch_type;", 0, "Base type `nosuch_type` is not a system data type."},
		{"create type from sys system type", "CREATE TYPE dbo.ty FROM sys.int;", -1, ""},
		{"create type from national char varying", "CREATE TYPE dbo.ty FROM NATIONAL CHAR VARYING(10);", -1, ""},
		// SELECT * without FROM and function body returning data.
		{"select star without from", "SELECT * INTO dbo.x;", 0, "`SELECT *` requires a `FROM` clause."},
		{"create view select star without from", "CREATE VIEW dbo.v2 AS SELECT *;", 0, "`SELECT *` requires a `FROM` clause."},
		{"insert select star without from", "INSERT INTO dbo.nokey SELECT *;", 0, "`SELECT *` requires a `FROM` clause."},
		{"inline tvf select star without from", "CREATE FUNCTION dbo.tvf() RETURNS TABLE RETURN SELECT *;", 0, "`SELECT *` requires a `FROM` clause."},
		{"create schema inline view select star", "CREATE SCHEMA s2 CREATE VIEW s2.v2 AS SELECT *;", 0, "`SELECT *` requires a `FROM` clause."},
		{"view with TABLESAMPLE", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t TABLESAMPLE (10 PERCENT);", 0, "SQL Server doesn't accept `TABLESAMPLE` in a view."},
		{"view with NEXT VALUE FOR", "CREATE VIEW dbo.v2 AS SELECT NEXT VALUE FOR dbo.sq AS n;", 0, "SQL Server doesn't accept `NEXT VALUE FOR` in a view."},
		{"view with SELECT INTO", "CREATE VIEW dbo.v2 AS SELECT a INTO dbo.x FROM dbo.t;", 0, "SQL Server doesn't accept `SELECT INTO` in a view."},
		{"ROW_NUMBER without ORDER BY", "CREATE VIEW dbo.v2 AS SELECT ROW_NUMBER() OVER () AS rn FROM dbo.t;", 0, "SQL Server's `ROW_NUMBER` needs an `OVER` clause with `ORDER BY`."},
		{"aggregate in WHERE", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t WHERE SUM(b) > 1;", 0, "SQL Server doesn't accept an aggregate in `WHERE`."},
		{"aggregate of an outer column in a subquery's WHERE", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t GROUP BY a HAVING EXISTS (SELECT 1 FROM dbo.r WHERE SUM(t.b) > 1);", -1, ""},
		{"subquery of two columns", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t WHERE a IN (SELECT id, id FROM dbo.r);", 0, "SQL Server doesn't accept a subquery of 2 columns in an expression."},
		{"FROM exposes a name twice", "CREATE VIEW dbo.v2 AS SELECT x.a FROM dbo.t x JOIN dbo.r x ON 1 = 1;", 0, "This query exposes `x` twice in its `FROM` clause."},
		{"derived table column list", "CREATE VIEW dbo.v2 AS SELECT x.a FROM (SELECT a FROM dbo.t) x (n);", 0, "Column `a` doesn't exist in `x`."},
		{"CTE column list", "CREATE VIEW dbo.v2 AS WITH x (n) AS (SELECT a FROM dbo.t) SELECT a FROM x;", 0, "Column `a` doesn't exist in `x`."},
		{"CTE name a later CTE takes reads the table", "CREATE VIEW dbo.v2 AS WITH x AS (SELECT id FROM r), r AS (SELECT 1 AS m) SELECT id FROM x;", -1, ""},
		{"CTE name of another statement reads the table", "CREATE PROCEDURE dbo.p AS BEGIN WITH r AS (SELECT 1 AS m) SELECT m FROM r; SELECT id FROM r; END;", -1, ""},
		{"PIVOT consumes its value column", "CREATE VIEW dbo.v2 AS SELECT p.b FROM (SELECT id, a, b FROM dbo.t) AS src PIVOT (SUM(b) FOR a IN ([x])) AS p;", 0, "Column `b` doesn't exist in `p`."},
		{"derived table reads a sibling under APPLY", "CREATE VIEW dbo.v2 AS SELECT t.a, x.n FROM dbo.t CROSS APPLY (SELECT t.b AS n) x;", -1, ""},
		{"derived table reads a sibling under JOIN", "CREATE VIEW dbo.v2 AS SELECT t.a FROM dbo.t JOIN (SELECT t.b AS n) x ON 1 = 1;", 0, "`t` is not a table or alias of the query."},
		{"table called as a function", "CREATE VIEW dbo.v2 AS SELECT * FROM dbo.t(1);", 0, "`dbo.t` is a table, not a function."},
		{"table with an old-style hint", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t(NOLOCK);", -1, ""},
		{"NOEXPAND on a table", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t WITH (NOEXPAND);", 0, "`dbo.t` is a table, not a view."},
		{"LEN of two arguments", "CREATE VIEW dbo.v2 AS SELECT LEN(a, 1) AS l FROM dbo.t;", 0, "SQL Server's `LEN` takes 1 argument."},
		{"unknown datepart", "CREATE VIEW dbo.v2 AS SELECT DATEADD(nope, 1, a) AS d FROM dbo.t;", 0, "`nope` is not a datepart of `DATEADD`."},
		{"HAVING names an ungrouped column", "CREATE VIEW dbo.v2 AS SELECT a FROM dbo.t GROUP BY a HAVING b = 1;", 0, "SQL Server doesn't accept column `b` in `HAVING`, which is neither grouped nor aggregated."},
		{"TRIM of an int column", "CREATE VIEW dbo.v2 AS SELECT TRIM(i) AS x FROM dbo.typed;", 0, "SQL Server's `TRIM` doesn't accept an argument of type `int`."},
		{"COLLATE on an int column", "CREATE VIEW dbo.v2 AS SELECT i COLLATE Latin1_General_BIN2 AS x FROM dbo.typed;", 0, "SQL Server doesn't accept `COLLATE` on an expression of type `int`."},
		{"select count star without from valid", "SELECT COUNT(*) AS c INTO dbo.x;", -1, ""},
		{"scalar function returning select data", "CREATE FUNCTION dbo.f2() RETURNS int AS BEGIN SELECT 1 END;", 0, "SELECT statements inside a function cannot return data to a client."},
		{"scalar function variable select assign valid", "CREATE FUNCTION dbo.f2() RETURNS int AS BEGIN DECLARE @x int; SELECT @x = 1; RETURN @x; END;", -1, ""},
		// Batch-local @variable tracking.
		{"delete undeclared table variable", "DELETE @tv;", 0, "Table variable `@tv` is not declared."},
		{"delete declared table variable", "DECLARE @tv TABLE (a int);\nDELETE @tv;", -1, ""},
		{"delete table variable across GO", "DECLARE @tv TABLE (a int);\nGO\nDELETE @tv;", 2, "Table variable `@tv` is not declared."},
		{"update undeclared scalar variable", "UPDATE dbo.t SET @v = 1;", 0, "Variable `@v` is not declared."},
		{"update declared scalar variable", "DECLARE @v int;\nUPDATE dbo.t SET @v = 1;", -1, ""},
		{"delete where current of undeclared cursor var", "DELETE dbo.t WHERE CURRENT OF @c;", 0, "Variable `@c` is not declared."},
		{"delete where current of declared cursor var", "DECLARE @c CURSOR;\nDELETE dbo.t WHERE CURRENT OF @c;", -1, ""},
		{"declare inside if counts as declared", "IF 1 = 0 DECLARE @v int;\nUPDATE dbo.t SET @v = 1;", -1, ""},
		// Index / constraint column checks, inline index without cols, XML/spatial index checks.
		{"create index duplicate column", "CREATE INDEX ix_dup ON dbo.t (a, a);", 0, "This statement lists column `a` twice."},
		{"create index duplicate include column", "CREATE INDEX ix_dup ON dbo.t (a) INCLUDE (a);", 0, "This statement lists column `a` twice."},
		{"create index missing column", "CREATE INDEX ix_bad ON dbo.t (nosuch);", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"alter table add index missing columns", "ALTER TABLE dbo.t ADD INDEX ix_nocol;", 0, "Index `ix_nocol` must specify at least one column."},
		{"create table clustered columnstore inline index valid", "CREATE TABLE dbo.cc (a int, INDEX ix_cc CLUSTERED COLUMNSTORE);", -1, ""},
		{"alter table add unique duplicate column", "ALTER TABLE dbo.t ADD UNIQUE (a, a);", 0, "This statement lists column `a` twice."},
		{"alter table add unique missing column", "ALTER TABLE dbo.t ADD UNIQUE (nosuch);", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"alter table add unique sees new column in same statement", "ALTER TABLE dbo.nokey ADD UNIQUE (x), x int;", -1, ""},
		{"alter table add fk without columns", "ALTER TABLE dbo.nokey ADD FOREIGN KEY REFERENCES dbo.r;", 0, "Table-level foreign key constraint must specify referencing columns."},
		{"alter table add fk missing column", "ALTER TABLE dbo.nokey ADD FOREIGN KEY (nosuch) REFERENCES dbo.r (id);", 0, "Column `nosuch` doesn't exist in `dbo.nokey`."},
		{"alter table switch to itself", "ALTER TABLE dbo.t SWITCH TO dbo.t;", 0, "Cannot switch table `dbo.t` to itself."},
		{"alter table switch to missing target", "ALTER TABLE dbo.t SWITCH TO dbo.nosuch;", 0, "Table `dbo.nosuch` doesn't exist."},
		{"alter table switch to other table valid", "ALTER TABLE dbo.t SWITCH TO dbo.r;", -1, ""},
		{"alter table add period duplicate column", "ALTER TABLE dbo.t ADD PERIOD FOR SYSTEM_TIME (a, a);", 0, "This statement lists column `a` twice."},
		{"alter table add period missing column", "ALTER TABLE dbo.t ADD PERIOD FOR SYSTEM_TIME (a, nosuch);", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"alter table add period distinct columns valid", "ALTER TABLE dbo.t ADD PERIOD FOR SYSTEM_TIME (a, b);", -1, ""},
		{"alter table add connection on non-edge table", "ALTER TABLE dbo.t ADD CONNECTION (dbo.t TO dbo.r);", 0, "Table `dbo.t` is not a graph edge table."},
		{"alter index for on non-xml index", "ALTER INDEX ix_a ON dbo.t FOR (REMOVE p);", 0, "`ix_a` is an index, not an XML index."},
		{"create secondary xml index missing parent", "CREATE XML INDEX sx1 ON dbo.t (a) USING XML INDEX nosuch FOR VALUE;", 0, "Index `nosuch` doesn't exist in `dbo.t`."},
		{"create secondary xml index non-xml parent", "CREATE XML INDEX sx1 ON dbo.t (a) USING XML INDEX ix_a FOR VALUE;", 0, "`ix_a` is an index, not an XML index."},
		// Column type incompatibilities and valid controls.
		{"create spatial index on int column", "CREATE SPATIAL INDEX sp1 ON dbo.typed (i);", 0, "Column `i` of type `int` cannot be used in a spatial index."},
		{"create spatial index on geography column valid", "CREATE SPATIAL INDEX sp1 ON dbo.typed (g);", -1, ""},
		{"create primary xml index on int column", "CREATE PRIMARY XML INDEX px1 ON dbo.typed (i);", 0, "Column `i` of type `int` cannot be used in an XML index."},
		{"create primary xml index on xml column valid", "CREATE PRIMARY XML INDEX px1 ON dbo.typed (x);", -1, ""},
		{"create index on xml key column", "CREATE INDEX ix_x ON dbo.typed (x);", 0, "Column `x` of type `xml` cannot be used as an index or statistics key column."},
		{"create index include xml column valid", "CREATE INDEX ix_inc ON dbo.typed (i) INCLUDE (x);", -1, ""},
		{"create index on geography key column", "CREATE INDEX ix_g ON dbo.typed (g);", 0, "Column `g` of type `geography` cannot be used as an index or statistics key column."},
		{"alter column invalidates synced column type", "ALTER TABLE dbo.typed ALTER COLUMN x int;\nCREATE INDEX ix_x ON dbo.typed (x);", -1, ""},
		// Trigger name kind, DML function call / EXEC target / DML column checks.
		{"disable trigger named like table", "DISABLE TRIGGER dbo.t ON dbo.t;", 0, "`dbo.t` is a table, not a trigger."},
		{"disable trigger in wrong schema", "DISABLE TRIGGER s.tr1 ON dbo.t;", 0, "Trigger `s.tr1` must be in the schema of `dbo.t`."},
		{"delete function call on table", "DELETE dbo.t();", 0, "`dbo.t` is a table, not a function."},
		{"delete function call on function allowed", "DELETE dbo.fn();", -1, ""},
		{"insert empty column list on table", "INSERT dbo.t() DEFAULT VALUES;", 0, "`dbo.t` is a table, not a function."},
		{"insert missing column", "INSERT INTO dbo.t (nosuch) VALUES (1);", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"insert identity default values", "INSERT INTO dbo.t ($IDENTITY) DEFAULT VALUES;", 0, "Cannot insert `DEFAULT VALUES` into `$IDENTITY`."},
		{"insert exec table target", "INSERT INTO dbo.nokey EXEC dbo.t;", 0, "`dbo.t` is a table, not a procedure."},
		{"update missing column", "UPDATE dbo.t SET nosuch = 1;", 0, "Column `nosuch` doesn't exist in `dbo.t`."},
		{"update identity pseudocolumn", "UPDATE dbo.t SET $IDENTITY = 1;", 0, "Cannot update `$IDENTITY` column."},
		{"update with from clause skips missing column check", "UPDATE dbo.t SET nosuch = 1 FROM dbo.r;", -1, ""},
		// A view's and an inline function's query bind when the module is
		// created (mssql/catalog/testdata/queries/simulate.tsv).
		{"default for a missing column", "ALTER TABLE dbo.t ADD DEFAULT 1 FOR nope;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"default for a column the statement adds", "ALTER TABLE dbo.t ADD c int, CONSTRAINT df_c DEFAULT 1 FOR c;", -1, ""},
		{"default for a column that has one", "ALTER TABLE dbo.t ADD DEFAULT 1 FOR a;", 0, "Column `a` of `dbo.t` already has a default."},
		{"default for a column whose default the change dropped", "ALTER TABLE dbo.t DROP CONSTRAINT df_a;\nALTER TABLE dbo.t ADD DEFAULT 1 FOR a;", -1, ""},
		{"masking function SQL Server lacks", "ALTER TABLE dbo.t ALTER COLUMN a ADD MASKED WITH (FUNCTION = 'nope()');", 0, "SQL Server doesn't accept the masking function `nope()`."},
		{"masking function SQL Server has", "ALTER TABLE dbo.t ALTER COLUMN a ADD MASKED WITH (FUNCTION = 'default()');", -1, ""},
		{"ROWGUIDCOL on an int column", "ALTER TABLE dbo.typed ALTER COLUMN i ADD ROWGUIDCOL;", 0, "SQL Server doesn't accept `ADD ROWGUIDCOL` on column `i` of type `int`."},
		{"HIDDEN on an xml column", "ALTER TABLE dbo.typed ALTER COLUMN x ADD HIDDEN;", 0, "SQL Server doesn't accept `ADD HIDDEN` on column `x` of type `xml`."},
		{"PERSISTED on a column that is not computed", "ALTER TABLE dbo.typed ALTER COLUMN g ADD PERSISTED;", 0, "SQL Server doesn't accept `ADD PERSISTED` on column `g` of type `geography`."},
		{"ROWGUIDCOL on a column of unknown type", "ALTER TABLE dbo.t ALTER COLUMN a ADD ROWGUIDCOL;", -1, ""},
		{"view column without a name", "CREATE VIEW dbo.w AS SELECT 1;", 0, "SQL Server doesn't accept a view column without a name."},
		{"view column named by an alias", "CREATE VIEW dbo.w AS SELECT 1 AS x;", -1, ""},
		{"view column named by its list", "CREATE VIEW dbo.w (x) AS SELECT 1;", -1, ""},
		{"view column in parentheses keeps its name", "CREATE VIEW dbo.w AS SELECT (a) FROM dbo.t;", -1, ""},
		{"view FOR XML returns an unnamed column", "CREATE VIEW dbo.w AS SELECT id FROM dbo.t FOR XML AUTO;", 0, "SQL Server doesn't accept a view column without a name."},
		{"view column list too short", "CREATE VIEW dbo.w (x) AS SELECT id, a FROM dbo.t;", 0, "SQL Server doesn't accept a view column list of 1 name for a query of 2 columns."},
		{"view column list with a star", "CREATE VIEW dbo.w (x) AS SELECT * FROM dbo.r;", -1, ""},
		{"view column list names a column twice", "CREATE VIEW dbo.w (x, x) AS SELECT id, a FROM dbo.t;", 0, "This statement names column `x` twice."},
		{"view names a column twice", "CREATE VIEW dbo.w AS SELECT t.id, r.id FROM dbo.t JOIN dbo.r ON r.id = t.id;", 0, "This statement names column `id` twice."},
		{"view names columns of two spellings", "CREATE VIEW dbo.w AS SELECT a, A AS a2 FROM dbo.t;", -1, ""},
		{"view reads a variable", "CREATE VIEW dbo.w AS SELECT @x AS v;", 0, "Variable `@x` is not declared."},
		{"view reads a table variable", "CREATE VIEW dbo.w AS SELECT * FROM @t;", 0, "Table variable `@t` is not declared."},
		{"view reads a system variable", "CREATE VIEW dbo.w AS SELECT @@SPID AS v;", -1, ""},
		{"view reads a temporary table", "CREATE VIEW dbo.w AS SELECT a FROM #x;", 0, "SQL Server doesn't accept a temporary table in a view."},
		{"view reads a missing table", "CREATE VIEW dbo.w AS SELECT a FROM dbo.nosuch;", 0, "`dbo.nosuch` doesn't exist."},
		{"view reads a table dropped before it", "DROP TABLE dbo.nokey;\nGO\nCREATE VIEW dbo.w AS SELECT a FROM dbo.nokey;", 2, "`dbo.nokey` doesn't exist."},
		{"view reads a CTE", "CREATE VIEW dbo.w AS WITH nosuch AS (SELECT 1 AS a) SELECT a FROM nosuch;", -1, ""},
		{"view reads a compatibility view", "CREATE VIEW dbo.w AS SELECT name FROM sysobjects;", -1, ""},
		{"view reads a missing column", "CREATE VIEW dbo.w AS SELECT nope FROM dbo.t;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"view reads a column by its alias", "CREATE VIEW dbo.w AS SELECT x.nope FROM dbo.t AS x;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"view reads a column of another case", "CREATE VIEW dbo.w AS SELECT ID FROM dbo.t;", -1, ""},
		{"view reads a column of the outer query", "CREATE VIEW dbo.w AS SELECT id FROM dbo.r WHERE EXISTS (SELECT 1 FROM dbo.nokey WHERE nokey.a = id);", -1, ""},
		{"view reads a column of either source", "CREATE VIEW dbo.w AS SELECT b FROM dbo.r JOIN dbo.t ON t.id = r.id;", -1, ""},
		{"view reads a column of no source", "CREATE VIEW dbo.w AS SELECT nope FROM dbo.r JOIN dbo.t ON t.id = r.id;", 0, "Column `nope` doesn't exist in `dbo.r`."},
		{"view reads a column of a view", "CREATE VIEW dbo.w AS SELECT nope FROM dbo.v;", -1, ""},
		{"view reads a column of a derived table", "CREATE VIEW dbo.w AS SELECT nope FROM (SELECT id FROM dbo.t) AS d;", 0, "Column `nope` doesn't exist in `d`."},
		{"view reads a column of a derived table's star", "CREATE VIEW dbo.w AS SELECT nope FROM (SELECT * FROM dbo.t) AS d;", -1, ""},
		{"view reads a column with no FROM", "CREATE VIEW dbo.w AS SELECT a;", 0, "Column `a` doesn't exist."},
		{"view reads a column of a name no source has", "CREATE VIEW dbo.w AS SELECT x.a FROM dbo.t;", 0, "`x` is not a table or alias of the query."},
		{"view reads a table its alias hides", "CREATE VIEW dbo.w AS SELECT t.a FROM dbo.t AS x;", 0, "`t` is not a table or alias of the query."},
		{"view reads a property of a column", "CREATE VIEW dbo.w AS SELECT a.STX AS x FROM dbo.t AS y;", -1, ""},
		{"view reads a derived table's unnamed column", "CREATE VIEW dbo.w AS SELECT n FROM (SELECT COUNT(*) FROM dbo.t) AS d;", 0, "SQL Server doesn't accept a derived table or CTE column without a name."},
		{"view names a derived table's columns", "CREATE VIEW dbo.w AS SELECT n FROM (SELECT COUNT(*) FROM dbo.t) AS d (n);", -1, ""},
		{"view defines an unused CTE with an unnamed column", "CREATE VIEW dbo.w AS WITH c AS (SELECT 1) SELECT 1 AS x;", 0, "SQL Server doesn't accept a derived table or CTE column without a name."},
		{"view reads VALUES without names", "CREATE VIEW dbo.w AS SELECT * FROM (VALUES (1)) AS v;", 0, "SQL Server doesn't accept a derived table or CTE column without a name."},
		{"view in another schema reads its schema's table", "CREATE VIEW s.w AS SELECT a FROM st;", -1, ""},
		{"view in another schema reads its schema's missing column", "CREATE VIEW s.w AS SELECT nope FROM st;", 0, "Column `nope` doesn't exist in `st`."},
		{"view in another schema reads dbo's table", "CREATE VIEW s.w AS SELECT b FROM t;", -1, ""},
		{"procedure in another schema reads its schema's table", "CREATE PROCEDURE s.p2 AS SELECT a FROM st;", -1, ""},
		{"view orders by an alias", "CREATE VIEW dbo.w AS SELECT TOP 1 a AS x FROM dbo.t ORDER BY x;", -1, ""},
		{"view orders without TOP", "CREATE VIEW dbo.w AS SELECT a FROM dbo.t ORDER BY a;", 0, "SQL Server doesn't accept `ORDER BY` without `TOP`, `OFFSET`, or `FOR XML` in a view."},
		{"view selects * without FROM in a subquery", "CREATE VIEW dbo.w AS SELECT (SELECT *) AS x;", 0, "`SELECT *` requires a `FROM` clause."},
		{"view selects * without FROM in EXISTS", "CREATE VIEW dbo.w AS SELECT a FROM dbo.t WHERE EXISTS (SELECT *);", -1, ""},
		{"view selects * without FROM in a set operation in EXISTS", "CREATE VIEW dbo.w AS SELECT a FROM dbo.t WHERE EXISTS (SELECT * UNION SELECT *);", 0, "`SELECT *` requires a `FROM` clause."},
		{"view reads a column two tables hold", "CREATE VIEW dbo.w AS SELECT a FROM dbo.t, dbo.t AS t2;", 0, "Column `a` is in more than one table of the query."},
		{"view orders by a column two tables hold through its alias", "CREATE VIEW dbo.w AS SELECT TOP 1 t.a FROM dbo.t, dbo.t AS t2 ORDER BY a;", -1, ""},
		{"view unions queries of two widths", "CREATE VIEW dbo.w AS SELECT id, a FROM dbo.t UNION SELECT id FROM dbo.r;", 0, "SQL Server doesn't accept a `UNION`, `INTERSECT`, or `EXCEPT` of queries with different numbers of columns."},
		{"view selects an ungrouped column", "CREATE VIEW dbo.w AS SELECT a, b FROM dbo.t GROUP BY a;", 0, "SQL Server doesn't accept column `b`, which is neither grouped nor aggregated."},
		{"view selects grouped columns", "CREATE VIEW dbo.w AS SELECT a, COUNT(*) AS n FROM dbo.t GROUP BY a;", -1, ""},
		{"schema-bound view names a table by one part", "CREATE VIEW dbo.w WITH SCHEMABINDING AS SELECT a FROM t;", 0, "SQL Server doesn't accept the one-part name `t` in a schema-bound view."},
		{"schema-bound view selects *", "CREATE VIEW dbo.w WITH SCHEMABINDING AS SELECT * FROM dbo.t;", 0, "SQL Server doesn't accept `*` in a schema-bound view."},
		{"inline function column without a name", "CREATE FUNCTION dbo.f2() RETURNS TABLE AS RETURN SELECT COUNT(*) FROM dbo.t;", 0, "SQL Server doesn't accept a function column without a name."},
		{"inline function reads its parameter", "CREATE FUNCTION dbo.f2(@c int) RETURNS TABLE AS RETURN SELECT a FROM dbo.t WHERE id = @c;", -1, ""},
		{"inline function reads a missing table", "CREATE FUNCTION dbo.f2() RETURNS TABLE AS RETURN SELECT a FROM dbo.nosuch;", 0, "`dbo.nosuch` doesn't exist."},
		// A procedure's, a function's, and a trigger's statements bind when
		// the module is created, each unless a table it names may not
		// exist (mssql/catalog/testdata/programs/simulate.tsv).
		{"procedure reads a missing column", "CREATE PROCEDURE dbo.p2 AS SELECT nope FROM dbo.t;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"procedure defers a statement on a missing table", "CREATE PROCEDURE dbo.p2 AS SELECT nope FROM dbo.t JOIN dbo.nosuch ON 1 = 1;", -1, ""},
		{"procedure defers a temporary table", "CREATE PROCEDURE dbo.p2 AS SELECT nope FROM #x;", -1, ""},
		{"procedure binds each statement", "CREATE PROCEDURE dbo.p2 AS BEGIN SELECT nope FROM dbo.nosuch; SELECT nope FROM dbo.t; END", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"procedure updates a missing column", "CREATE PROCEDURE dbo.p2 AS UPDATE dbo.t SET nope = 1;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"procedure deletes by a missing column", "CREATE PROCEDURE dbo.p2 AS DELETE FROM dbo.t WHERE nope = 1;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"procedure inserts a missing column", "CREATE PROCEDURE dbo.p2 AS INSERT INTO dbo.t (nope) VALUES (1);", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"procedure merges on a missing column", "CREATE PROCEDURE dbo.p2 AS MERGE dbo.t AS x USING dbo.r AS y ON x.id = y.nope WHEN MATCHED THEN DELETE;", 0, "Column `nope` doesn't exist in `dbo.r`."},
		{"procedure that changes a table binds nothing", "CREATE PROCEDURE dbo.p2 AS BEGIN ALTER TABLE dbo.t ADD nope int; SELECT nope FROM dbo.t; END", -1, ""},
		{"procedure reads a column with no FROM", "CREATE PROCEDURE dbo.p2 AS SELECT a;", 0, "Column `a` doesn't exist."},
		{"procedure passes a named argument", "CREATE PROCEDURE dbo.p2 AS EXEC sp_executesql N'SELECT @a', N'@a int', @a = 1;", -1, ""},
		{"procedure reads an undeclared variable", "CREATE PROCEDURE dbo.p2 AS SELECT @x AS v;", 0, "Variable `@x` is not declared."},
		{"procedure declares a variable after it", "CREATE PROCEDURE dbo.p2 AS BEGIN SELECT @n AS v; DECLARE @n int; END", 0, "Variable `@n` is not declared."},
		{"procedure reads its parameter", "CREATE PROCEDURE dbo.p2 @x int AS SELECT @X AS v;", -1, ""},
		{"procedure reads a declared variable", "CREATE PROCEDURE dbo.p2 AS BEGIN DECLARE @n int = 1; SELECT @n AS v; END", -1, ""},
		{"procedure reads a table variable's missing column", "CREATE PROCEDURE dbo.p2 AS BEGIN DECLARE @t TABLE (x int); SELECT y FROM @t; END", 0, "Column `y` doesn't exist in `@t`."},
		{"procedure inserts fewer values than columns", "CREATE PROCEDURE dbo.p2 AS INSERT INTO dbo.t (a) VALUES (1, 2);", 0, "SQL Server doesn't accept an `INSERT` that lists 1 column and 2 values."},
		{"insert lists more columns than its query", "INSERT INTO dbo.t (a, b) SELECT id FROM dbo.r;", 0, "SQL Server doesn't accept an `INSERT` that lists 2 columns and 1 value."},
		{"insert with a star query", "INSERT INTO dbo.t (a, b) SELECT * FROM dbo.r;", -1, ""},
		{"union of two widths", "SELECT id, a FROM dbo.t UNION SELECT id FROM dbo.r;", 0, "SQL Server doesn't accept a `UNION`, `INTERSECT`, or `EXCEPT` of queries with different numbers of columns."},
		{"trigger reads a missing column of inserted", "CREATE TRIGGER dbo.tr ON dbo.t AFTER INSERT AS SELECT nope FROM inserted;", 0, "Column `nope` doesn't exist in `inserted`."},
		{"trigger reads a column of deleted", "CREATE TRIGGER dbo.tr ON dbo.t AFTER DELETE AS SELECT a FROM deleted;", -1, ""},
		{"trigger tests a column of its table", "CREATE TRIGGER dbo.tr ON dbo.t AFTER UPDATE AS IF UPDATE(a) RETURN;", -1, ""},
		{"trigger tests a missing column", "CREATE TRIGGER dbo.tr ON dbo.t AFTER UPDATE AS IF UPDATE(nope) RETURN;", 0, "Column `nope` doesn't exist in `dbo.t`."},
		{"trigger tests a column in a query of another table", "CREATE TRIGGER dbo.tr ON dbo.t AFTER UPDATE AS SELECT id FROM dbo.r WHERE UPDATE(a);", -1, ""},
		{"altered procedure binds in the schema it resolves to", "ALTER PROCEDURE pr AS SELECT a FROM t;", -1, ""},
		{"AFTER trigger on a view", "CREATE TRIGGER dbo.tr ON dbo.v AFTER INSERT AS SELECT 1;", 0, "`dbo.v` is a view, not a table."},
		{"INSTEAD OF trigger on a view", "CREATE TRIGGER dbo.tr ON dbo.v INSTEAD OF INSERT AS SELECT 1;", -1, ""},
		{"function changes a table", "CREATE FUNCTION dbo.f2() RETURNS int AS BEGIN INSERT INTO dbo.t (a) VALUES (1); RETURN 1 END", 0, "SQL Server doesn't accept `INSERT` of a table in a function."},
		{"function changes its table variable", "CREATE FUNCTION dbo.f2() RETURNS @r TABLE (id int) AS BEGIN INSERT INTO @r VALUES (1); RETURN; END", -1, ""},
		{"function calls NEWID", "CREATE FUNCTION dbo.f2() RETURNS uniqueidentifier AS BEGIN RETURN NEWID() END", 0, "SQL Server doesn't accept `NEWID` in a function."},
		{"function does not end with RETURN", "CREATE FUNCTION dbo.f2(@x int) RETURNS int AS BEGIN IF @x > 0 RETURN 1; END", 0, "SQL Server doesn't accept a function whose last statement isn't `RETURN`."},
		{"function ends with IF and ELSE", "CREATE FUNCTION dbo.f2(@x int) RETURNS int AS BEGIN IF @x > 0 RETURN 1 ELSE RETURN 0 END", 0, "SQL Server doesn't accept a function whose last statement isn't `RETURN`."},
		{"function ends with a block that returns", "CREATE FUNCTION dbo.f2() RETURNS int AS BEGIN BEGIN RETURN 1 END END", -1, ""},
		{"function's table declares a column twice", "CREATE FUNCTION dbo.f2() RETURNS @r TABLE (id int, id int) AS BEGIN RETURN; END", 0, "This statement declares column `id` twice."},
		{"function reads a missing column of its table", "CREATE FUNCTION dbo.f2() RETURNS @r TABLE (id int) AS BEGIN INSERT INTO @r (nope) VALUES (1); RETURN; END", 0, "Column `nope` doesn't exist in `@r`."},
		{"schema-bound function names a table by one part", "CREATE FUNCTION dbo.f2() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT COUNT(*) FROM t) END", 0, "SQL Server doesn't accept the one-part name `t` in a schema-bound function."},
		{"function reads an undeclared variable", "CREATE FUNCTION dbo.f2() RETURNS int AS BEGIN RETURN @x END", 0, "Variable `@x` is not declared."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := walkTarget()
			switch tt.name {
			case "a marked name may be the table another spelling means", "the exact spelling means its own table", "a dropped spelling may leave a marked table":
				dbo := target.Schema.Schemas[0]
				dbo.Tables = append(dbo.Tables,
					&metadata.TableMetadata{Name: "cafe", Columns: []*metadata.ColumnMetadata{{Name: "a"}}},
					&metadata.TableMetadata{Name: "CAFÉ", Columns: []*metadata.ColumnMetadata{{Name: "b"}}})
			case "a marked schema may be the one another spelling means", "the exact schema spelling means its own schema":
				target.Schema.Schemas = append(target.Schema.Schemas,
					&metadata.SchemaMetadata{Name: "cafe", Tables: []*metadata.TableMetadata{{Name: "t"}}},
					&metadata.SchemaMetadata{Name: "CAFÉ"})
			case "a dropped schema beside a marked one":
				target.Schema.Schemas = append(target.Schema.Schemas,
					&metadata.SchemaMetadata{Name: "cafe"},
					&metadata.SchemaMetadata{Name: "CAFÉ", Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "a"}}}}})
			case "an unplaced drop may reach a marked name":
				target.DefaultSchema = ""
				target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{Name: "CAFÉ"})
			case "an XML index's name is no object", "a primary XML index takes its secondary ones", "a secondary XML index":
				t := target.Schema.Schemas[0].Tables[0]
				t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: "px", Type: "XML"}, &metadata.IndexMetadata{Name: "sx", Type: "XML"})
			case "a marked name of another case", "an index on a marked name of another case":
				target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{Name: "Café", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}}})
			case "another spelling beside a marked name":
				target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{Name: "é"})
			case "a constraint dropped by an accent-equal spelling":
				target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{Name: "tc", Columns: []*metadata.ColumnMetadata{{Name: "id"}}, Indexes: []*metadata.IndexMetadata{{Name: "pk_cafe", Primary: true, Unique: true, IsConstraint: true}}})
			case "a dropped ledger table frees its history table and view", "a table of one NOT NULL bigint is no ledger table":
				cols := []*metadata.ColumnMetadata{{Name: "a", Type: "int", Nullable: true}, {Name: "tx", Type: "bigint"}, {Name: "sq", Type: "bi"}}
				if tt.name == "a table of one NOT NULL bigint is no ledger table" {
					cols[2].Nullable = true
				}
				target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{Name: "lt", Columns: cols})
				s := target.Schema.Schemas[len(target.Schema.Schemas)-1]
				s.Tables = append(s.Tables, &metadata.TableMetadata{Name: "lt_h"})
				s.Views = append(s.Views, &metadata.ViewMetadata{Name: "lt_v"})
			case "an unqualified type in the default schema":
				target.DefaultSchema = "s"
				s := target.Schema.Schemas[len(target.Schema.Schemas)-1]
				s.Tables = append(s.Tables, &metadata.TableMetadata{Name: "ty", Columns: []*metadata.ColumnMetadata{{Name: "a"}}})
			case "a dropped table's capture instance", "a change table without change data capture":
				cdc := &metadata.SchemaMetadata{Name: "cdc", Tables: []*metadata.TableMetadata{{Name: "dbo_nokey_CT"}}}
				if tt.name == "a dropped table's capture instance" {
					cdc.Tables = append(cdc.Tables, &metadata.TableMetadata{Name: "change_tables"})
				}
				target.Schema.Schemas = append(target.Schema.Schemas, cdc)
			case "no default schema":
				target.DefaultSchema = ""
			case "altered procedure binds in the schema it resolves to":
				target.DefaultSchema = "s"
				s := target.Schema.Schemas[1]
				s.Tables = append(s.Tables, &metadata.TableMetadata{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "b"}}})
			case "a default schema that may not exist":
				target.DefaultSchema = "nosuch"
			case "create spatial index on int column", "create spatial index on geography column valid",
				"create primary xml index on int column", "create primary xml index on xml column valid",
				"create index on xml key column", "create index include xml column valid",
				"create index on geography key column", "alter column invalidates synced column type",
				"ROWGUIDCOL on an int column", "HIDDEN on an xml column", "PERSISTED on a column that is not computed",
				"TRIM of an int column", "COLLATE on an int column":
				target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{
					Name: "typed",
					Columns: []*metadata.ColumnMetadata{
						{Name: "i", Type: "int"},
						{Name: "x", Type: "xml"},
						{Name: "g", Type: "geography"},
					},
				})
			}
			got := walkFindings(t, tt.sql, target)
			if tt.statement < 0 {
				if len(got) != 0 {
					t.Errorf("findings = %+v, want none", got)
				}
				return
			}
			if len(got) != 1 || got[0].Statement != tt.statement || got[0].Message != tt.message {
				t.Errorf("findings = %+v, want statement %d: %q", got, tt.statement, tt.message)
			}
			// The oracle reads the message for the errors it stands for.
			if (&walkEngine{}).Codes(tt.message) == nil {
				t.Errorf("walkCodes maps no error to %q", tt.message)
			}
		})
	}
}

// TestWalkThroughTargetGroups pins DefaultSchema and OtherSpecs in the
// target grouping: two targets of one schema that differ in either are
// walked apart.
func TestWalkThroughTargetGroups(t *testing.T) {
	schema := walkSchema()
	got := walkFindings(t, "CREATE TABLE st (a int);", review.Target{Schema: schema, DefaultSchema: "dbo"}, review.Target{Schema: schema, DefaultSchema: "s"})
	if len(got) != 1 || len(got[0].Targets) != 1 || got[0].Targets[0] != 1 {
		t.Errorf("DefaultSchema: findings = %+v, want one on target 1", got)
	}
	got = walkFindings(t, "CREATE TABLE dbo.t (a int);", review.Target{Schema: schema, DefaultSchema: "dbo"}, review.Target{Schema: schema, DefaultSchema: "dbo", OtherSpecs: true})
	if len(got) != 1 || len(got[0].Targets) != 1 || got[0].Targets[0] != 0 {
		t.Errorf("OtherSpecs: findings = %+v, want one on target 0", got)
	}
}

// TestWalkThroughAmbiguousParents pins what the snapshot keeps of a table
// or schema whose class another synced one shares: the walk keeps none of
// its children, but its triggers and the names its children hold in the
// schema count.
func TestWalkThroughAmbiguousParents(t *testing.T) {
	trigger := []*metadata.TriggerMetadata{{Name: "mk"}}
	tests := []struct {
		name      string
		edit      func(db *metadata.DatabaseSchemaMetadata)
		sql       string
		statement int // -1 for no finding
		message   string
	}{
		{"a trigger on the second table of a class", func(db *metadata.DatabaseSchemaMetadata) {
			dbo := db.Schemas[0]
			dbo.Tables = append(dbo.Tables, &metadata.TableMetadata{Name: "t1"}, &metadata.TableMetadata{Name: "t11", Triggers: trigger})
		}, "INSERT INTO dbo.t11 VALUES (1);\nDROP TABLE dbo.nosuch;", -1, ""},
		{"a trigger on a table of the second schema of a class", func(db *metadata.DatabaseSchemaMetadata) {
			db.Schemas = append(db.Schemas, &metadata.SchemaMetadata{Name: "s1"}, &metadata.SchemaMetadata{Name: "s11", Tables: []*metadata.TableMetadata{{Name: "x", Triggers: trigger}}})
		}, "INSERT INTO dbo.nokey VALUES (1);\nDROP TABLE dbo.nosuch;", -1, ""},
		{"a trigger on a view of the second schema of a class", func(db *metadata.DatabaseSchemaMetadata) {
			db.Schemas = append(db.Schemas, &metadata.SchemaMetadata{Name: "s1"}, &metadata.SchemaMetadata{Name: "s11", Views: []*metadata.ViewMetadata{{Name: "v", Triggers: trigger}}})
		}, "INSERT INTO dbo.nokey VALUES (1);\nDROP TABLE dbo.nosuch;", -1, ""},
		{"no trigger", func(db *metadata.DatabaseSchemaMetadata) {
			dbo := db.Schemas[0]
			dbo.Tables = append(dbo.Tables, &metadata.TableMetadata{Name: "t1"}, &metadata.TableMetadata{Name: "t11"})
		}, "INSERT INTO dbo.t11 VALUES (1);\nDROP TABLE dbo.nosuch;", 1, "Table `dbo.nosuch` doesn't exist."},
		{"the second table's constraint names", func(db *metadata.DatabaseSchemaMetadata) {
			dbo := db.Schemas[0]
			dbo.Tables = append(dbo.Tables, &metadata.TableMetadata{Name: "t1"}, &metadata.TableMetadata{Name: "t11", CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck_t11"}}})
		}, "CREATE TABLE dbo.ck_t11 (a int);", 0, "The name `dbo.ck_t11` is already taken."},
		{"the second table's children are unknown", func(db *metadata.DatabaseSchemaMetadata) {
			dbo := db.Schemas[0]
			dbo.Tables = append(dbo.Tables, &metadata.TableMetadata{Name: "t1"}, &metadata.TableMetadata{Name: "t11", Columns: []*metadata.ColumnMetadata{{Name: "a"}}})
		}, "ALTER TABLE dbo.t11 DROP COLUMN b;\nALTER TABLE dbo.t1 DROP COLUMN b;", -1, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := walkTarget()
			tt.edit(target.Schema)
			got := walkFindings(t, tt.sql, target)
			if tt.statement < 0 {
				if len(got) != 0 {
					t.Errorf("findings = %+v, want none", got)
				}
				return
			}
			if len(got) != 1 || got[0].Statement != tt.statement || got[0].Message != tt.message {
				t.Errorf("findings = %+v, want statement %d: %q", got, tt.statement, tt.message)
			}
		})
	}
}

// walkForms are the statement forms of the spec's §5 table, each on
// walkSchema with dbo the default schema. A form the table lists must
// lower: one the parser gives in a shape the lowering does not expect
// would end the walk unseen.
var walkForms = []walktest.Form{
	// Batches and the session.
	{SQL: "CREATE TABLE dbo.x (a int);\nGO\nCREATE TABLE dbo.y (a int);"},
	{SQL: "USE shop;"},
	{SQL: "USE master;\nGO\nUSE shop;"},
	{SQL: "SET NOCOUNT ON;\nSET XACT_ABORT ON;\nSET NOEXEC OFF;\nSET QUOTED_IDENTIFIER ON;"},
	{SQL: "DECLARE @x int = 1;\nSET @x = 2;\nPRINT 'x';"},
	{SQL: "BEGIN TRANSACTION;\nSAVE TRANSACTION sp;\nCOMMIT;"},
	{SQL: "GRANT SELECT ON dbo.t TO public;\nDENY SELECT ON dbo.t TO public;\nREVOKE SELECT ON dbo.t FROM public;"},
	{SQL: "UPDATE STATISTICS dbo.t;\nWAITFOR DELAY '00:00:00';"},
	{SQL: "RAISERROR ('x', 10, 1);\nTHROW 50000, 'x', 1;"},
	// Schemas.
	{SQL: "CREATE SCHEMA s2;"},
	{SQL: "CREATE SCHEMA s3 AUTHORIZATION dbo;"},
	{SQL: "DROP SCHEMA IF EXISTS s2;"},
	{SQL: "DROP TABLE s.st;\nDROP SCHEMA s;"},
	// Tables.
	{SQL: "CREATE TABLE dbo.x (id int NOT NULL CONSTRAINT pk_x PRIMARY KEY, a int DEFAULT 0 CONSTRAINT ck_x CHECK (a > 0), b int UNIQUE, c int REFERENCES dbo.r (id), INDEX ix_x (a));"},
	{SQL: "CREATE TABLE x (a int);"},
	{SQL: "CREATE TABLE #x (a int PRIMARY KEY);"},
	{SQL: "CREATE TABLE shop.dbo.x (a int);"},
	{SQL: "SELECT 1 AS a INTO dbo.x;"},
	{SQL: "SELECT 1 AS a INTO #x;"},
	{SQL: "SELECT 1 AS a;"},
	{SQL: "ALTER TABLE dbo.t ADD x int, y int CONSTRAINT df_y DEFAULT 0, CONSTRAINT ck_x CHECK (x > 0);"},
	{SQL: "ALTER TABLE dbo.nokey ADD CONSTRAINT pk_n PRIMARY KEY (a);"},
	{SQL: "ALTER TABLE dbo.t WITH CHECK ADD CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES dbo.r (id);"},
	{SQL: "ALTER TABLE dbo.t DROP COLUMN b;"},
	{SQL: "ALTER TABLE dbo.t DROP COLUMN IF EXISTS nosuch, CONSTRAINT IF EXISTS nosuch2;"},
	{SQL: "ALTER TABLE dbo.t DROP CONSTRAINT ck_b, fk_r;"},
	{SQL: "ALTER TABLE dbo.t DROP ck_b;"},
	{SQL: "ALTER TABLE dbo.t ALTER COLUMN b bigint NOT NULL;"},
	{SQL: "ALTER TABLE dbo.t ALTER COLUMN b ADD SPARSE;"},
	{SQL: "ALTER TABLE dbo.t SWITCH TO dbo.r;", Last: true},
	{SQL: "CREATE TABLE dbo.e AS EDGE;", Last: true},
	{SQL: "CREATE SCHEMA s4 CREATE TABLE s4.x (a int);", Last: true},
	{SQL: "CREATE PRIMARY XML INDEX px_new ON dbo.t (a);"},
	{SQL: "CREATE SELECTIVE XML INDEX sxi_new ON dbo.t (a) FOR (p = '/a');"},
	{SQL: "CREATE SPATIAL INDEX sp_new ON dbo.t (a);"},
	{SQL: "INSERT INTO dbo.nokey EXEC dbo.pr;", Last: true},
	{SQL: "ALTER TABLE dbo.t NOCHECK CONSTRAINT ck_b;\nALTER TABLE dbo.t WITH CHECK CHECK CONSTRAINT ALL;"},
	{SQL: "ALTER TABLE dbo.t DISABLE TRIGGER ALL;"},
	{SQL: "ALTER TABLE dbo.t SET (LOCK_ESCALATION = AUTO);\nALTER TABLE dbo.t REBUILD;"},
	{SQL: "ALTER TABLE dbo.t ENABLE TRIGGER ALL;"},
	{SQL: "DROP TABLE dbo.nokey;"},
	{SQL: "DROP TABLE IF EXISTS dbo.x, #x, dbo.nokey;"},
	{SQL: "TRUNCATE TABLE dbo.nokey;"},
	// Indexes.
	{SQL: "CREATE UNIQUE NONCLUSTERED INDEX ix2 ON dbo.t (b);"},
	{SQL: "CREATE INDEX ix_a ON dbo.t (a) WITH (DROP_EXISTING = ON);"},
	{SQL: "CREATE INDEX ix_v ON dbo.v (id);"},
	{SQL: "ALTER INDEX ix_a ON dbo.t REBUILD;\nALTER INDEX ALL ON dbo.t REORGANIZE;"},
	{SQL: "DROP INDEX ix_a ON dbo.t;"},
	{SQL: "DROP INDEX IF EXISTS nosuch ON dbo.t, dbo.t.ix_a;"},
	// Views, routines, triggers, each in a batch of its own.
	{SQL: "CREATE VIEW dbo.w AS SELECT 1 AS a;"},
	{SQL: "CREATE OR ALTER VIEW dbo.v AS SELECT 1 AS id;"},
	{SQL: "ALTER VIEW dbo.v AS SELECT 1 AS id;"},
	{SQL: "DROP VIEW dbo.v;"},
	{SQL: "DROP VIEW IF EXISTS dbo.w, dbo.v;"},
	{SQL: "CREATE PROCEDURE dbo.p2 AS SELECT 1;\nGO\nALTER PROCEDURE dbo.p2 AS SELECT 2;"},
	{SQL: "CREATE OR ALTER PROCEDURE dbo.pr AS BEGIN SELECT 1 END;\nSELECT 2;"},
	{SQL: "CREATE FUNCTION dbo.f2() RETURNS int AS BEGIN RETURN 1 END;"},
	{SQL: "ALTER FUNCTION dbo.fn() RETURNS int AS BEGIN RETURN 2 END;"},
	{SQL: "DROP PROCEDURE dbo.pr;\nDROP FUNCTION IF EXISTS dbo.fn;"},
	{SQL: "CREATE TRIGGER dbo.tr ON dbo.t AFTER INSERT AS SELECT 1;"},
	{SQL: "CREATE TRIGGER tr ON dbo.t AFTER INSERT AS SELECT 1;\nGO\nDROP TRIGGER dbo.tr;"},
	{SQL: "ENABLE TRIGGER ALL ON dbo.t;\nDISABLE TRIGGER ALL ON dbo.t;"},
	// Sequences, synonyms, types, statistics, defaults, schema transfer.
	{SQL: "CREATE SEQUENCE dbo.sq2 START WITH 1;"},
	{SQL: "ALTER SEQUENCE dbo.sq RESTART;"},
	{SQL: "DROP SEQUENCE dbo.sq;"},
	{SQL: "CREATE SYNONYM dbo.syn FOR dbo.t;\nDROP SYNONYM dbo.syn;"},
	{SQL: "CREATE TYPE dbo.ty FROM int NOT NULL;\nCREATE TYPE dbo.tt AS TABLE (a int);\nDROP TYPE dbo.ty;\nDROP TYPE IF EXISTS dbo.tt;"},
	{SQL: "CREATE STATISTICS st1 ON dbo.t (a);\nUPDATE STATISTICS dbo.t;"},
	{SQL: "CREATE DEFAULT dbo.d1 AS 1;\nGO\nDROP DEFAULT dbo.d1;"},
	{SQL: "DROP RULE IF EXISTS dbo.r1;"},
	{SQL: "ALTER SCHEMA s TRANSFER dbo.t;", Last: true},
	{SQL: "ALTER TABLE dbo.t SWITCH TO dbo.r;", Last: true},
	{SQL: "ALTER TABLE dbo.t ADD PERIOD FOR SYSTEM_TIME (a, b);", Last: true},
	// sp_rename and extended properties.
	{SQL: "EXEC sp_rename 'dbo.nokey', 'n2';"},
	{SQL: "EXEC sp_rename @objname = N'dbo.v', @newname = N'v2', @objtype = N'OBJECT';"},
	{SQL: "EXEC sys.sp_rename 'dbo.t.a', 'z', 'COLUMN';"},
	{SQL: "EXEC sp_rename 'dbo.t.ix_a', 'ix_z', 'INDEX';"},
	{SQL: "EXEC sp_addextendedproperty N'MS_Description', N'x', N'SCHEMA', N'dbo', N'TABLE', N't';"},
	// DML: no trigger in walkSchema, so the walk goes on after an INSERT,
	// and after an UPDATE, DELETE, or MERGE of what no foreign key
	// references.
	{SQL: "INSERT INTO dbo.nokey (a) VALUES (1);"},
	{SQL: "DECLARE @tv TABLE (a int);\nINSERT INTO @tv VALUES (1);\nUPDATE @tv SET a = 2;\nDELETE FROM @tv;"},
	{SQL: "CREATE TABLE #x (a int);\nINSERT INTO #x VALUES (1);\nUPDATE #x SET a = 2;\nDELETE FROM #x;\nMERGE #x AS x USING dbo.r AS y ON x.a = y.id WHEN MATCHED THEN DELETE;"},
	// IF.
	{SQL: "IF OBJECT_ID('dbo.t') IS NOT NULL DROP TABLE dbo.t;"},
	{SQL: "IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'x') BEGIN CREATE TABLE dbo.x (a int); CREATE INDEX ix ON dbo.x (a) END ELSE ALTER TABLE dbo.x ADD b int;"},
	// DML on a table a trigger may fire on, or through a synonym, which may
	// lead to another database: lowered, and the walk ends after it.
	{SQL: "CREATE TRIGGER dbo.tr ON dbo.nokey AFTER INSERT AS SELECT 1;\nGO\nINSERT INTO dbo.nokey VALUES (1);", Last: true},
	{SQL: "INSERT INTO dbo.syn (a) VALUES (1);", Last: true},
	{SQL: "INSERT INTO other.dbo.t VALUES (1);", Last: true},
	// An UPDATE, DELETE, or MERGE of a table may cascade to a table whose
	// triggers the sync does not list.
	{SQL: "UPDATE dbo.nokey SET a = 1;", Last: true},
	{SQL: "DELETE FROM dbo.nokey WHERE a = 1;", Last: true},
	{SQL: "UPDATE x SET a = 1 FROM dbo.nokey AS x;", Last: true},
	{SQL: "MERGE dbo.nokey AS x USING dbo.r AS y ON x.a = y.id WHEN MATCHED THEN DELETE;", Last: true},
}

func TestWalkThroughLowersEveryForm(t *testing.T) {
	walktest.Coverage(t, walkForms, func(t *testing.T, sql string) (n, lowered int, last bool) {
		t.Helper()
		_, stmts, syntax := parseReview(sql)
		if syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, sql)
		}
		w := newWalker(sql, stmts, walkTarget())
		l := walk.NewLedger(w.snap)
		for i := range stmts {
			s := w.lower(i, l)
			if s == nil {
				return len(stmts), i, false
			}
			if s.Fail != nil {
				t.Errorf("%q fails at statement %d", sql, i)
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
