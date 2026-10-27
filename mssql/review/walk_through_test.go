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
	for _, f := range result.Findings {
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
		{"table exists", "CREATE TABLE dbo.t (a int);", 0, "an object named dbo.t already exists"},
		{"table named like a constraint", "CREATE TABLE dbo.pk_t (a int);", 0, "an object named dbo.pk_t already exists"},
		{"case variant is unknown", "CREATE TABLE dbo.T (a int);", -1, ""},
		{"unqualified with the default schema", "CREATE TABLE t (a int);", 0, "an object named t already exists"},
		{"missing schema", "CREATE TABLE nosuch.x (a int);", 0, "schema nosuch does not exist"},
		{"other schema is free", "CREATE TABLE s.t (a int);", -1, ""},
		{"created twice", "CREATE TABLE dbo.x (a int);\nCREATE TABLE dbo.x (a int);", 1, "an object named dbo.x already exists"},
		{"column twice", "CREATE TABLE dbo.x (a int, a int);", 0, "column a is named twice in the statement"},
		{"two primary keys", "CREATE TABLE dbo.x (a int PRIMARY KEY, b int PRIMARY KEY);", 0, "the statement adds two primary keys to dbo.x"},
		{"constraint name taken", "CREATE TABLE dbo.x (a int CONSTRAINT ck_b CHECK (a > 0));", 0, "an object named ck_b already exists"},
		{"constraint named like the table", "CREATE TABLE dbo.x (a int CONSTRAINT x CHECK (a > 0));", 0, "an object named x already exists"},
		{"constraint name twice", "CREATE TABLE dbo.x (a int CONSTRAINT c1 CHECK (a > 0), b int CONSTRAINT c1 CHECK (b > 0));", 0, "the name c1 is given twice in the statement"},
		{"foreign key to a missing table", "CREATE TABLE dbo.x (a int REFERENCES dbo.nosuch (id));", 0, "foreign key target dbo.nosuch does not exist"},
		{"foreign key to a view", "CREATE TABLE dbo.x (a int REFERENCES dbo.v (id));", 0, "foreign key target dbo.v is not a table"},
		{"foreign key to itself", "CREATE TABLE dbo.x (id int PRIMARY KEY, p int REFERENCES dbo.x (id));", -1, ""},
		{"view over a table", "CREATE VIEW dbo.t AS SELECT 1 AS a;", 0, "an object named dbo.t already exists"},
		{"sequence over a procedure", "CREATE SEQUENCE dbo.pr;", 0, "an object named dbo.pr already exists"},
		{"temporary table twice", "CREATE TABLE #x (a int);\nCREATE TABLE #x (a int);", 1, "an object named #x already exists"},
		{"select into an existing table", "SELECT 1 AS a INTO dbo.t;", 0, "an object named dbo.t already exists"},
		// Drops.
		{"drop missing table", "DROP TABLE dbo.nosuch;", 0, "table dbo.nosuch does not exist"},
		{"drop a view as a table", "DROP TABLE dbo.v;", 0, "dbo.v is not a table"},
		{"drop if exists of a view", "DROP TABLE IF EXISTS dbo.v;", 0, "dbo.v is not a table"},
		{"drop if exists of a procedure", "DROP TABLE IF EXISTS dbo.pr;", 0, "dbo.pr is not a table"},
		{"drop if exists missing", "DROP TABLE IF EXISTS dbo.nosuch;", -1, ""},
		{"drop twice", "DROP TABLE dbo.t;\nDROP TABLE dbo.t;", 1, "table dbo.t does not exist"},
		{"drop twice in one", "DROP TABLE dbo.r, dbo.r;", 0, "dbo.r is dropped twice"},
		{"drop case variant twice", "DROP TABLE dbo.t;\nDROP TABLE dbo.T;", 1, "table dbo.T does not exist"},
		{"drop then create", "DROP TABLE dbo.t;\nCREATE TABLE dbo.t (a int);\nCREATE TABLE dbo.pk_t (a int);", -1, ""},
		{"drop then use", "DROP TABLE dbo.nokey;\nINSERT INTO dbo.nokey VALUES (1);", 1, "object dbo.nokey does not exist"},
		{"drop missing view", "DROP VIEW dbo.nosuch;", 0, "view dbo.nosuch does not exist"},
		{"drop a table as a view", "DROP VIEW dbo.t;", 0, "dbo.t is not a view"},
		{"drop procedure unsynced may exist", "DROP PROCEDURE dbo.nosuch;", -1, ""},
		{"drop procedure of a table", "DROP PROCEDURE dbo.t;", 0, "dbo.t is not a procedure"},
		{"drop function of a procedure", "DROP FUNCTION dbo.pr;", 0, "dbo.pr is not a function"},
		{"drop procedure twice", "DROP PROCEDURE dbo.pr;\nDROP PROCEDURE dbo.pr;", 1, "procedure dbo.pr does not exist"},
		{"drop sequence missing", "DROP SEQUENCE dbo.nosuch;", 0, "sequence dbo.nosuch does not exist"},
		{"drop schema missing", "DROP SCHEMA nosuch;", 0, "schema nosuch does not exist"},
		{"drop schema with a table", "DROP SCHEMA s;", 0, "schema s still holds objects"},
		{"drop schema emptied", "DROP TABLE s.st;\nDROP SCHEMA s;", -1, ""},
		{"create schema twice", "CREATE SCHEMA s;", 0, "schema s already exists"},
		// ALTER TABLE.
		{"alter missing table", "ALTER TABLE dbo.nosuch ADD x int;", 0, "table dbo.nosuch does not exist"},
		{"alter a view", "ALTER TABLE dbo.v ADD x int;", 0, "dbo.v is not a table"},
		{"add existing column", "ALTER TABLE dbo.t ADD a int;", 0, "column a already exists in dbo.t"},
		{"add column twice", "ALTER TABLE dbo.t ADD x int, x int;", 0, "column x is named twice in the statement"},
		{"drop missing column", "ALTER TABLE dbo.t DROP COLUMN nosuch;", 0, "column nosuch does not exist in dbo.t"},
		{"drop missing column if exists", "ALTER TABLE dbo.t DROP COLUMN IF EXISTS nosuch;", -1, ""},
		{"drop column twice", "ALTER TABLE dbo.t DROP COLUMN b, COLUMN b;", 0, "the name b is given twice in the statement"},
		{"alter missing column", "ALTER TABLE dbo.t ALTER COLUMN nosuch int;", 0, "column nosuch does not exist in dbo.t"},
		{"second primary key", "ALTER TABLE dbo.t ADD CONSTRAINT pk2 PRIMARY KEY (id);", 0, "table dbo.t already has a primary key"},
		{"primary key on a keyless table", "ALTER TABLE dbo.nokey ADD CONSTRAINT pk2 PRIMARY KEY (a);", -1, ""},
		{"constraint name taken in ALTER", "ALTER TABLE dbo.nokey ADD CONSTRAINT pk_t CHECK (a > 0);", 0, "an object named pk_t already exists"},
		{"unique named like an index", "ALTER TABLE dbo.t ADD CONSTRAINT ix_a UNIQUE (b);", 0, "index ix_a already exists on dbo.t"},
		{"drop missing constraint", "ALTER TABLE dbo.t DROP CONSTRAINT nosuch;", 0, "nosuch is not a constraint of dbo.t"},
		{"drop another table's constraint", "ALTER TABLE dbo.r DROP CONSTRAINT pk_t;", 0, "pk_t is not a constraint of dbo.r"},
		{"drop a default", "ALTER TABLE dbo.t DROP CONSTRAINT df_a;", -1, ""},
		{"drop then drop the constraint", "ALTER TABLE dbo.t DROP CONSTRAINT ck_b;\nALTER TABLE dbo.t DROP CONSTRAINT ck_b;", 1, "ck_b is not a constraint of dbo.t"},
		{"primary key dropped then added", "ALTER TABLE dbo.t DROP CONSTRAINT pk_t;\nALTER TABLE dbo.t ADD CONSTRAINT pk_t2 PRIMARY KEY (id);", -1, ""},
		{"nocheck missing constraint", "ALTER TABLE dbo.t NOCHECK CONSTRAINT nosuch;", 0, "nosuch is not a constraint of dbo.t"},
		{"column added then dropped twice", "ALTER TABLE dbo.t ADD x int;\nALTER TABLE dbo.t DROP COLUMN x;\nALTER TABLE dbo.t DROP COLUMN x;", 2, "column x does not exist in dbo.t"},
		// Indexes.
		{"index exists", "CREATE INDEX ix_a ON dbo.t (b);", 0, "index ix_a already exists on dbo.t"},
		{"index named like the key", "CREATE INDEX pk_t ON dbo.t (b);", 0, "index pk_t already exists on dbo.t"},
		{"index on another table is free", "CREATE INDEX pk_t ON dbo.r (id);", -1, ""},
		{"index on a missing table", "CREATE INDEX i ON dbo.nosuch (a);", 0, "table dbo.nosuch does not exist"},
		{"drop missing index", "DROP INDEX nosuch ON dbo.t;", 0, "index nosuch does not exist on dbo.t"},
		{"drop a key's index", "DROP INDEX pk_t ON dbo.t;", 0, "index pk_t of dbo.t enforces a constraint"},
		{"drop a key's index if exists", "DROP INDEX IF EXISTS pk_t ON dbo.t;", 0, "index pk_t of dbo.t enforces a constraint"},
		{"old drop index form", "DROP INDEX dbo.t.ix_a;\nDROP INDEX dbo.t.ix_a;", 1, "index ix_a does not exist on dbo.t"},
		{"drop_existing of a missing index", "CREATE INDEX nosuch ON dbo.t (a) WITH (DROP_EXISTING = ON);", 0, "index nosuch does not exist on dbo.t"},
		{"alter missing index", "ALTER INDEX nosuch ON dbo.t REBUILD;", 0, "index nosuch does not exist on dbo.t"},
		{"truncate a view", "TRUNCATE TABLE dbo.v;", 0, "dbo.v is not a table"},
		{"truncate missing", "TRUNCATE TABLE dbo.nosuch;", 0, "table dbo.nosuch does not exist"},
		// Modules.
		{"alter missing view", "ALTER VIEW dbo.nosuch AS SELECT 1 AS a;", 0, "view dbo.nosuch does not exist"},
		{"alter a table as a view", "ALTER VIEW dbo.t AS SELECT 1 AS a;", 0, "dbo.t is not a view"},
		{"create or alter a view over a table", "CREATE OR ALTER VIEW dbo.t AS SELECT 1 AS a;", 0, "dbo.t is not a view"},
		{"create or alter a procedure over a function", "CREATE OR ALTER PROCEDURE dbo.fn AS SELECT 1;", 0, "dbo.fn is not a procedure"},
		{"alter procedure unsynced", "ALTER PROCEDURE dbo.nosuch AS SELECT 1;", -1, ""},
		{"view not first", "CREATE TABLE dbo.x (a int);\nCREATE VIEW dbo.w AS SELECT 1 AS a;", 1, "CREATE VIEW must be the first statement in its batch"},
		{"view after GO", "CREATE TABLE dbo.x (a int);\nGO\nCREATE VIEW dbo.w AS SELECT 1 AS a;", -1, ""},
		{"procedure body swallows", "CREATE PROCEDURE dbo.p2 AS SELECT 1;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"procedure body ends at GO", "CREATE PROCEDURE dbo.p2 AS SELECT 1;\nCREATE TABLE dbo.x (a int);\nGO\nCREATE TABLE dbo.t (a int);", 3, "an object named dbo.t already exists"},
		{"trigger on a missing table", "CREATE TRIGGER dbo.tr ON dbo.nosuch AFTER INSERT AS SELECT 1;", 0, "table dbo.nosuch does not exist"},
		{"trigger in another schema", "CREATE TRIGGER s.tr ON dbo.t AFTER INSERT AS SELECT 1;", 0, "trigger s.tr must be in the schema of its table dbo.t"},
		{"alter sequence of a table", "ALTER SEQUENCE dbo.t RESTART;", 0, "dbo.t is not a sequence"},
		{"type twice", "CREATE TYPE dbo.ty FROM int;\nCREATE TYPE dbo.ty FROM int;", 1, "type dbo.ty already exists"},
		{"type dropped twice", "CREATE TYPE dbo.ty FROM int;\nDROP TYPE dbo.ty;\nDROP TYPE dbo.ty;", 2, "type dbo.ty does not exist"},
		// sp_rename.
		{"rename to a taken name", "EXEC sp_rename 'dbo.t', 'r';", 0, "an object named r already exists"},
		{"rename a dropped table", "DROP TABLE dbo.nokey;\nEXEC sp_rename 'dbo.nokey', 'x';", 1, "object dbo.nokey does not exist"},
		{"rename then use", "EXEC sp_rename 'dbo.nokey', 'n2';\nALTER TABLE dbo.n2 ADD a int;", 1, "column a already exists in dbo.n2"},
		{"rename frees the name", "EXEC sp_rename 'dbo.nokey', 'n2';\nCREATE TABLE dbo.nokey (a int);", -1, ""},
		{"rename column", "EXEC sp_rename 'dbo.t.a', 'z', 'COLUMN';\nALTER TABLE dbo.t DROP COLUMN a;", 1, "column a does not exist in dbo.t"},
		{"rename column to a taken name", "EXEC sp_rename 'dbo.t.a', 'b', 'COLUMN';", 0, "column b already exists in dbo.t"},
		{"rename missing column", "EXEC sp_rename 'dbo.t.nosuch', 'z', 'COLUMN';", 0, "column nosuch does not exist in dbo.t"},
		{"rename index", "EXEC sp_rename 'dbo.t.ix_a', 'ix_b', 'INDEX';\nDROP INDEX ix_a ON dbo.t;", 1, "index ix_a does not exist on dbo.t"},
		// Unknown by design.
		{"no default schema", "CREATE TABLE t (a int);", -1, ""},
		{"DML on unsynced name", "INSERT INTO dbo.nosuch VALUES (1);", -1, ""},
		{"exec ends the walk", "EXEC dbo.pr;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"if makes effects optional", "IF OBJECT_ID('dbo.t') IS NOT NULL DROP TABLE dbo.t;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"if keeps the walk going", "IF 1 = 1 DROP TABLE dbo.nokey;\nCREATE TABLE dbo.r (a int);", 1, "an object named dbo.r already exists"},
		{"try ends the walk", "BEGIN TRY DROP TABLE dbo.nosuch END TRY BEGIN CATCH END CATCH;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"GO forgets temporary tables", "CREATE TABLE #x (a int);\nGO\nCREATE TABLE #x (a int);", -1, ""},
		{"GO not alone ends the walk", "SELECT 1 GO\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"GO twice ends the walk", "CREATE TABLE dbo.x (a int)\nGO 2\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"USE elsewhere", "USE master;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"USE back", "USE master;\nUSE shop;\nCREATE TABLE dbo.t (a int);", 2, "an object named dbo.t already exists"},
		{"three-part target", "CREATE TABLE shop.dbo.t (a int);", 0, "an object named shop.dbo.t already exists"},
		{"three-part elsewhere", "CREATE TABLE other.dbo.t (a int);", -1, ""},
		{"noexec ends the walk", "SET NOEXEC ON;\nCREATE TABLE dbo.t (a int);", -1, ""},
		{"a schema elsewhere", "USE master;\nGO\nCREATE SCHEMA s;\nGO\nDROP SCHEMA nosuch;", -1, ""},
		{"a schema in an unknown database", "USE SHOP;\nGO\nCREATE SCHEMA s;\nGO\nUSE shop;\nCREATE TABLE s.st (a int);", -1, ""},
		{"a default schema that may not exist", "CREATE TABLE t (a int);", -1, ""},
		{"a missing schema holds no table", "DROP TABLE nosuch.t;", 0, "table nosuch.t does not exist"},
		// Codex round 1 on #514.
		{"a marked name may be the table another spelling means", "ALTER TABLE dbo.CAFE ADD a int;", -1, ""},
		{"the exact spelling means its own table", "ALTER TABLE dbo.cafe ADD a int;", 0, "column a already exists in dbo.cafe"},
		{"a dropped spelling may leave a marked table", "DROP TABLE dbo.cafe;\nDROP TABLE dbo.CAFE;", -1, ""},
		{"a created marked name", "CREATE TABLE dbo.[NOKÉY] (b int);\nALTER TABLE dbo.NOKEY ADD a int;", -1, ""},
		{"a check and an index of one name", "CREATE TABLE dbo.x (a int, CONSTRAINT ix CHECK (a > 0), INDEX ix NONCLUSTERED (a));", 0, "the name ix is given twice in the statement"},
		{"a unique constraint and an index of one name", "CREATE TABLE dbo.x (a int CONSTRAINT u1 UNIQUE, INDEX u1 (a));", 0, "the name u1 is given twice in the statement"},
		{"DROP CONSTRAINT IF EXISTS of another object's name", "ALTER TABLE dbo.t DROP CONSTRAINT IF EXISTS nokey;\nINSERT INTO dbo.nokey VALUES (1);\nCREATE TABLE dbo.nokey (a int);", 2, "an object named dbo.nokey already exists"},
		// Codex round 2 on #514.
		{"a marked schema may be the one another spelling means", "CREATE TABLE CAFE.t (a int);", -1, ""},
		{"the exact schema spelling means its own schema", "CREATE TABLE cafe.t (a int);", 0, "an object named cafe.t already exists"},
		{"sp_rename cuts a long name", "EXEC sp_rename 'dbo.nokey', '" + strings.Repeat("n", 129) + "';", -1, ""},
		{"a primary XML index takes its secondary ones", "DROP INDEX px ON dbo.t;\nCREATE INDEX sx ON dbo.t (a);", -1, ""},
		{"a secondary XML index", "CREATE INDEX sx ON dbo.t (a);", 0, "index sx already exists on dbo.t"},
		{"a self-reference of another spelling", "CREATE TABLE dbo.x (id int NOT NULL PRIMARY KEY, p int REFERENCES dbo.X (id));", -1, ""},
		{"an alias of a marked spelling", "DROP TABLE dbo.nokey;\nUPDATE a SET a = 1 FROM dbo.r AS [á];", -1, ""},
		{"a marked arrival leaves its space's absences unknown", "CREATE TABLE dbo.[é] (a int UNIQUE);\nDROP TABLE dbo.cafe;", -1, ""},
		// Codex round 3 on #514.
		{"an XML index's name is no object", "CREATE TABLE dbo.px (a int);", -1, ""},
		{"a dropped schema beside a marked one", "DROP SCHEMA cafe;\nALTER TABLE CAFE.t ADD z int;", -1, ""},
		{"an unplaced drop may reach a marked name", "DROP TABLE CAFE;\nCREATE TABLE dbo.[CAFÉ] (a int);", -1, ""},
		{"dotted names are distinct", "DROP TABLE dbo.[x.y], [dbo.x].y;", -1, ""},
		{"DROP_EXISTING keeps a key's constraint", "CREATE UNIQUE CLUSTERED INDEX pk_t ON dbo.t (id) WITH (DROP_EXISTING = ON);\nDROP INDEX pk_t ON dbo.t;", 1, "index pk_t of dbo.t enforces a constraint"},
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
			case "no default schema":
				target.DefaultSchema = ""
			case "a default schema that may not exist":
				target.DefaultSchema = "nosuch"
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
	// Sequences, synonyms, types.
	{SQL: "CREATE SEQUENCE dbo.sq2 START WITH 1;"},
	{SQL: "ALTER SEQUENCE dbo.sq RESTART;"},
	{SQL: "DROP SEQUENCE dbo.sq;"},
	{SQL: "CREATE SYNONYM dbo.syn FOR dbo.t;\nDROP SYNONYM dbo.syn;"},
	{SQL: "CREATE TYPE dbo.ty FROM int NOT NULL;\nCREATE TYPE dbo.tt AS TABLE (a int);\nDROP TYPE dbo.ty;\nDROP TYPE IF EXISTS dbo.tt;"},
	// sp_rename and extended properties.
	{SQL: "EXEC sp_rename 'dbo.nokey', 'n2';"},
	{SQL: "EXEC sp_rename @objname = N'dbo.v', @newname = N'v2', @objtype = N'OBJECT';"},
	{SQL: "EXEC sys.sp_rename 'dbo.t.a', 'z', 'COLUMN';"},
	{SQL: "EXEC sp_rename 'dbo.t.ix_a', 'ix_z', 'INDEX';"},
	{SQL: "EXEC sp_addextendedproperty N'MS_Description', N'x', N'SCHEMA', N'dbo', N'TABLE', N't';"},
	// DML: no trigger in walkSchema, so the walk goes on.
	{SQL: "INSERT INTO dbo.nokey (a) VALUES (1);"},
	{SQL: "UPDATE dbo.nokey SET a = 1;\nDELETE FROM dbo.nokey WHERE a = 1;"},
	{SQL: "UPDATE x SET a = 1 FROM dbo.nokey AS x;"},
	{SQL: "MERGE dbo.nokey AS x USING dbo.r AS y ON x.a = y.id WHEN MATCHED THEN DELETE;"},
	{SQL: "DECLARE @tv TABLE (a int);\nINSERT INTO @tv VALUES (1);"},
	{SQL: "CREATE TABLE #x (a int);\nINSERT INTO #x VALUES (1);"},
	// IF.
	{SQL: "IF OBJECT_ID('dbo.t') IS NOT NULL DROP TABLE dbo.t;"},
	{SQL: "IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'x') BEGIN CREATE TABLE dbo.x (a int); CREATE INDEX ix ON dbo.x (a) END ELSE ALTER TABLE dbo.x ADD b int;"},
	// DML on a table a trigger may fire on, or an object that may be a
	// synonym: lowered, and the walk ends after it.
	{SQL: "CREATE TRIGGER dbo.tr ON dbo.nokey AFTER INSERT AS SELECT 1;\nGO\nINSERT INTO dbo.nokey VALUES (1);", Last: true},
	{SQL: "INSERT INTO dbo.nosuch VALUES (1);", Last: true},
	{SQL: "INSERT INTO other.dbo.t VALUES (1);", Last: true},
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
