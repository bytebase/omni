package review

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// finding is a finding as the tests compare it: the statement, the text
// its range anchors, the message, and the targets.
type finding struct {
	Statement int
	Text      string
	Message   string
	Targets   []int
}

func run(t *testing.T, rule review.Rule, sql string, change review.Change, targets []review.Target) []finding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{rule}, Change: change}, targets)
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	if len(result.Failures) != 0 {
		t.Fatalf("Failures = %v", result.Failures)
	}
	var got []finding
	for _, f := range result.Findings {
		if f.Rule == review.Syntax && rule != review.Syntax {
			t.Fatalf("syntax finding: %s", f.Message)
		}
		if f.Rule != rule {
			continue
		}
		got = append(got, finding{f.Statement, sql[f.Range.Start:f.Range.End], f.Message, f.Targets})
	}
	return got
}

type ruleCase struct {
	name    string
	sql     string
	change  review.Change
	targets []review.Target
	want    []finding
}

func runCases(t *testing.T, rule review.Rule, cases []ruleCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			targets := tt.targets
			if targets == nil {
				targets = []review.Target{{Schema: shop(), BackupDatabaseExists: true}}
			}
			got := run(t, rule, tt.sql, tt.change, targets)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// shop is a small database: dbo.t with a primary key, a unique
// constraint, a check, and a foreign key to dbo.p; dbo.nokey without a
// key; dbo.v, a view; and s.t, a second t in another schema.
func shop() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{
						Name:    "p",
						Columns: columns("id", "note"),
						Indexes: []*metadata.IndexMetadata{{Name: "pk_p", Primary: true, Unique: true}},
					},
					{
						Name:    "t",
						Columns: columns("id", "code", "p_id", "n"),
						Indexes: []*metadata.IndexMetadata{
							{Name: "pk_t", Primary: true, Unique: true},
							{Name: "uk_code", Unique: true, IsConstraint: true},
							{Name: "ix_n"},
							{Name: "ux_n", Unique: true},
						},
						ForeignKeys:      []*metadata.ForeignKeyMetadata{{Name: "fk_p", ReferencedTable: "p"}},
						CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck_n"}},
					},
					{Name: "nokey", Columns: columns("a")},
				},
				Views: []*metadata.ViewMetadata{{Name: "v"}},
			},
			{
				Name:   "s",
				Tables: []*metadata.TableMetadata{{Name: "t2", Columns: columns("id")}},
			},
		},
	}
}

func columns(names ...string) []*metadata.ColumnMetadata {
	out := make([]*metadata.ColumnMetadata, len(names))
	for i, n := range names {
		out[i] = &metadata.ColumnMetadata{Name: n}
	}
	return out
}

func TestSyntax(t *testing.T) {
	runCases(t, review.Syntax, []ruleCase{
		{name: "valid", sql: "SELECT 1; UPDATE t SET n = 1 WHERE id = 1"},
		{
			name: "one finding for the first error",
			sql:  "SELECT 1;\nSELEC 2;\nSELEC 3;",
			want: []finding{{1, "", "syntax error at or near \"SELEC\"", []int{0}}},
		},
	})
}

func TestRequireWhere(t *testing.T) {
	runCases(t, review.RequireWhere, []ruleCase{
		{
			name: "UPDATE and DELETE without WHERE",
			sql:  "UPDATE t SET n = 1;\nDELETE FROM dbo.t;\nUPDATE t SET n = 1 WHERE id = 1;\nDELETE x FROM t AS x JOIN p ON x.p_id = p.id;",
			want: []finding{
				{0, "UPDATE t SET n = 1", "UPDATE t has no WHERE clause", []int{0}},
				{1, "DELETE FROM dbo.t", "DELETE FROM dbo.t has no WHERE clause", []int{0}},
				{3, "DELETE x FROM t AS x JOIN p ON x.p_id = p.id", "DELETE FROM x has no WHERE clause", []int{0}},
			},
		},
		{
			name: "inside IF and a procedure",
			sql:  "IF 1 = 1 DELETE FROM t;\nGO\nCREATE PROCEDURE pr AS BEGIN UPDATE t SET n = 1 END",
			want: []finding{
				{0, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0}},
				{2, "UPDATE t SET n = 1", "UPDATE t has no WHERE clause", []int{0}},
			},
		},
	})
}

func TestRequireIsNull(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "comparisons a predicate tests",
			sql:  "SELECT * FROM t WHERE n = NULL OR NULL <> code;\nUPDATE t SET n = 1 WHERE NOT (code != NULL);\nIF @x = NULL PRINT 1;",
			want: []finding{
				{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}},
				{0, "NULL <> code", `"NULL <> code" is never true; use IS NOT NULL`, []int{0}},
				{1, "code != NULL", `"code != NULL" is never true; use IS NOT NULL`, []int{0}},
				{2, "@x = NULL", `"@x = NULL" is never true; use IS NULL`, []int{0}},
			},
		},
		{
			name: "a value, IS NULL, and after SET ANSI_NULLS OFF",
			sql:  "SELECT CASE WHEN 1 = 1 THEN 1 END FROM t WHERE n IS NULL;\nSET ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;",
		},
		{
			name: "ANSI_NULLS set back ON",
			sql:  "SET ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_NULLS ON;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_DEFAULTS OFF;\nSELECT * FROM t WHERE n = NULL;",
			want: []finding{{3, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}}},
		},
		{
			name: "ANSI_NULLS set in a block, under a condition, or in a routine body",
			sql:  "BEGIN SET ANSI_NULLS OFF; END;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_NULLS ON;\nIF @x = 1 SET ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_NULLS ON;\nCREATE PROCEDURE pr2 AS BEGIN SET ANSI_NULLS OFF; END;\nGO\nSELECT * FROM t WHERE n = NULL;",
			want: []finding{{8, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}}},
		},
		{
			name: "IIF",
			sql:  "SELECT IIF(n = NULL, 1, 0) FROM t;",
			want: []finding{{0, "n = NULL", `"n = NULL" is never true; use IS NULL`, []int{0}}},
		},
		{
			name: "MERGE conditions",
			sql:  "MERGE t USING p ON t.id = NULL WHEN MATCHED AND p.note <> NULL THEN UPDATE SET n = 1;",
			want: []finding{
				{0, "t.id = NULL", `"t.id = NULL" is never true; use IS NULL`, []int{0}},
				{0, "p.note <> NULL", `"p.note <> NULL" is never true; use IS NOT NULL`, []int{0}},
			},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t, dbo.p;\nDROP VIEW v;\nDROP INDEX ix_n ON t;\nDROP PROCEDURE pr;\nDROP FUNCTION f;\nDROP SCHEMA s;\nDROP DATABASE d;\nDROP USER u;\nALTER TABLE t DROP COLUMN n, code;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t, dbo.p", "drops table t, dbo.p", []int{0}},
				{1, "DROP VIEW v", "drops view v", []int{0}},
				{2, "DROP INDEX ix_n ON t", "drops index ix_n", []int{0}},
				{3, "DROP PROCEDURE pr", "drops procedure pr", []int{0}},
				{4, "DROP FUNCTION f", "drops function f", []int{0}},
				{5, "DROP SCHEMA s", "drops schema s", []int{0}},
				{6, "DROP DATABASE d", "drops database d", []int{0}},
				{7, "DROP USER u", "drops user u", []int{0}},
				{8, "n", "drops column n of t", []int{0}},
				{8, "code", "drops column code of t", []int{0}},
			},
		},
		{name: "a constraint belongs to DisallowDropConstraint", sql: "ALTER TABLE t DROP CONSTRAINT ck_n;"},
		{
			name: "a name with a line break",
			sql:  "DROP TABLE [a\nb];",
			want: []finding{{0, "DROP TABLE [a\nb]", "drops table [a\\nb]", []int{0}}},
		},
		{
			name: "DROP SENSITIVITY CLASSIFICATION and DROP SIGNATURE",
			sql:  "DROP SENSITIVITY CLASSIFICATION FROM dbo.t.n;\nDROP SIGNATURE FROM pr BY CERTIFICATE c1;",
			want: []finding{
				{0, "DROP SENSITIVITY CLASSIFICATION FROM dbo.t.n", "drops sensitivity classification of dbo.t.n", []int{0}},
				{1, "DROP SIGNATURE FROM pr BY CERTIFICATE c1", "drops signature of pr", []int{0}},
			},
		},
		{
			name: "a name with a part left out",
			sql:  "DROP TABLE shop..t;",
			want: []finding{{0, "DROP TABLE shop..t", "drops table shop..t", []int{0}}},
		},
	})
}

func TestDisallowTruncate(t *testing.T) {
	runCases(t, review.DisallowTruncate, []ruleCase{
		{
			name: "TRUNCATE TABLE",
			sql:  "TRUNCATE TABLE t;\nTRUNCATE TABLE dbo.p WITH (PARTITIONS (1));",
			want: []finding{
				{0, "TRUNCATE TABLE t", "truncates t", []int{0}},
				{1, "TRUNCATE TABLE dbo.p WITH (PARTITIONS (1))", "truncates partitions of dbo.p", []int{0}},
			},
		},
	})
}

func TestDisallowRename(t *testing.T) {
	runCases(t, review.DisallowRename, []ruleCase{
		{
			name: "sp_rename of a column, a table, and a view",
			sql:  "EXEC sp_rename 'dbo.t.n', 'm', 'COLUMN';\nEXEC sp_rename @objname = N't', @newname = N't3';\nEXEC sys.sp_rename 'v', 'v2', 'OBJECT';\nEXEC sp_rename 'dbo.t.ix_n', 'ix_m', 'INDEX';\nEXEC sp_rename 'pr', 'pr2';",
			want: []finding{
				{0, "EXEC sp_rename 'dbo.t.n', 'm', 'COLUMN'", "renames column n of dbo.t to m", []int{0}},
				{1, "EXEC sp_rename @objname = N't', @newname = N't3'", "renames table t to t3", []int{0}},
				{2, "EXEC sys.sp_rename 'v', 'v2', 'OBJECT'", "renames view v to v2", []int{0}},
			},
		},
		{
			name: "an object rename after a USE",
			sql:  "USE other;\nEXEC sp_rename 't', 'x';",
		},
	})
}

func TestDisallowDropConstraint(t *testing.T) {
	runCases(t, review.DisallowDropConstraint, []ruleCase{
		{
			name: "constraints of the synced schema",
			sql:  "ALTER TABLE t DROP CONSTRAINT pk_t, fk_p;\nALTER TABLE dbo.t DROP CONSTRAINT uk_code;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT df_n;",
			want: []finding{
				{0, "pk_t", "drops primary key pk_t of t", []int{0}},
				{0, "fk_p", "drops foreign key fk_p of t", []int{0}},
				{1, "uk_code", "drops unique constraint uk_code of dbo.t", []int{0}},
				{2, "ck_n", "drops check constraint ck_n of t", []int{0}},
			},
		},
		{
			name: "constraints the change defined",
			sql:  "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;\nALTER TABLE nokey ADD CONSTRAINT df_a DEFAULT 0 FOR a;\nALTER TABLE nokey DROP CONSTRAINT df_a;",
			want: []finding{{1, "ck_a", "drops check constraint ck_a of nokey", []int{0}}},
		},
		{
			name: "a constraint defined on a table of another database",
			sql:  "ALTER TABLE dbo.t ADD CONSTRAINT ck_x CHECK (n > 0);\nALTER TABLE other.dbo.t DROP CONSTRAINT ck_x;",
		},
		{
			name: "a constraint the change defined, then a USE",
			sql:  "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nUSE other;\nALTER TABLE nokey DROP CONSTRAINT ck_a;",
		},
		{
			name: "a unique index, not a constraint",
			sql:  "ALTER TABLE t DROP CONSTRAINT ux_n;",
		},
		{
			name: "inside a routine body or IF",
			sql:  "IF 1 = 1 ALTER TABLE t DROP CONSTRAINT fk_p;\nGO\nCREATE PROCEDURE pr AS ALTER TABLE t DROP CONSTRAINT ck_n;",
			want: []finding{
				{0, "fk_p", "drops foreign key fk_p of t", []int{0}},
				{2, "ck_n", "drops check constraint ck_n of t", []int{0}},
			},
		},
	})
}

func TestRequirePrimaryKey(t *testing.T) {
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name: "CREATE TABLE without a primary key",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);\nCREATE TABLE c (id INT, CONSTRAINT pk_c PRIMARY KEY (id));\nCREATE TABLE #tmp (id INT);",
			want: []finding{{0, "CREATE TABLE a (id INT)", "creates table a without a primary key", []int{0}}},
		},
		{
			name: "settled by a later key, drop, or rename",
			sql:  "CREATE TABLE a (id INT NOT NULL);\nALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id);\nCREATE TABLE b (id INT);\nDROP TABLE b;\nCREATE TABLE c (id INT);\nEXEC sp_rename 'c', 'c2';\nCREATE DATABASE d;\nCREATE TABLE d.dbo.e (id INT);\nDROP DATABASE d;",
		},
		{
			name: "a primary key dropped and not added back",
			sql:  "ALTER TABLE t DROP CONSTRAINT pk_t;\nALTER TABLE p DROP CONSTRAINT pk_p;\nALTER TABLE p ADD CONSTRAINT pk_p2 PRIMARY KEY (id, note);",
			want: []finding{{0, "pk_t", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
		{
			name: "a FileTable, whose fixed schema has a primary key",
			sql:  "CREATE TABLE files AS FILETABLE;",
		},
		{
			name: "settled in a block, or rolled back",
			sql:  "CREATE TABLE a2 (id INT NOT NULL);\nBEGIN ALTER TABLE a2 ADD CONSTRAINT pk_a2 PRIMARY KEY (id); END;\nBEGIN TRAN;\nALTER TABLE t DROP CONSTRAINT pk_t;\nROLLBACK;",
		},
		{
			name: "in a CREATE SCHEMA",
			sql:  "CREATE SCHEMA s2 CREATE TABLE s2.k (id INT);",
			want: []finding{{0, "CREATE TABLE s2.k (id INT)", "creates table s2.k without a primary key", []int{0}}},
		},
		{
			name: "not settled by a column rename, nor judged in a routine body",
			sql:  "CREATE TABLE f (id INT);\nRENAME OBJECT f COLUMN id TO id2;\nGO\nCREATE PROCEDURE pr AS ALTER TABLE t DROP CONSTRAINT pk_t;",
			want: []finding{{0, "CREATE TABLE f (id INT)", "creates table f without a primary key", []int{0}}},
		},
	})
}

func TestOnlineMigration(t *testing.T) {
	runCases(t, review.OnlineMigration, []ruleCase{
		{name: "not asked for", sql: "ALTER TABLE t ADD m INT;"},
		{
			name:   "asked for",
			sql:    "ALTER TABLE t ADD m INT;",
			change: review.Change{OnlineMigration: true},
			want:   []finding{{-1, "", "the change requests online migration, which SQL Server does not support", []int{0}}},
		},
	})
}

func TestPriorBackup(t *testing.T) {
	backup := review.Change{PriorBackup: true, MaxBackupSize: 1 << 20}
	runCases(t, review.PriorBackup, []ruleCase{
		{name: "a plain backup", sql: "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM p WHERE id = 2;", change: backup},
		{
			name:    "no bbdataarchive",
			sql:     "ALTER TABLE t ADD m INT;",
			change:  backup,
			targets: []review.Target{{Schema: shop()}},
			want:    []finding{{-1, "", "prior backup needs database bbdataarchive, which does not exist", []int{0}}},
		},
		{
			name:   "what the backup refuses",
			sql:    "DECLARE @v TABLE (id INT);\nUPDATE @v SET id = 1 WHERE id = 2;\nWITH c AS (SELECT * FROM t) DELETE FROM c WHERE id = 1;\nUPDATE TOP (3) p SET note = '' WHERE id > 0;",
			change: backup,
			want: []finding{
				{1, "UPDATE @v SET id = 1 WHERE id = 2", "prior backup does not support the table variable @v as the target", []int{0}},
				{2, "WITH c AS (SELECT * FROM t) DELETE FROM c WHERE id = 1", "prior backup does not support the WITH query c as the target", []int{0}},
				{3, "UPDATE TOP (3) p SET note = '' WHERE id > 0", "prior backup copies TOP after the select list of its query, which SQL Server rejects", []int{0}},
			},
		},
		{
			name:   "UPDATE and DELETE on one table, and different WITH clauses",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM dbo.t WHERE id = 2;\nWITH a AS (SELECT 1 AS id) UPDATE p SET note = '' WHERE id IN (SELECT id FROM a);\nUPDATE p SET note = 'x' WHERE id = 2;",
			change: backup,
			want: []finding{
				{1, "DELETE FROM dbo.t WHERE id = 2", "prior backup cannot back up both UPDATE and DELETE on dbo.t", []int{0}},
				{3, "UPDATE p SET note = 'x' WHERE id = 2", "prior backup cannot back up statements on p with different WITH clauses", []int{0}},
			},
		},
		{
			name:   "a name without a database, the target's",
			sql:    "DELETE FROM shop.dbo.p WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nDELETE FROM other.dbo.nokey WHERE a = 1;\nUPDATE nokey SET a = 2 WHERE a = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE p SET note = '' WHERE id = 1", "prior backup cannot back up both UPDATE and DELETE on p", []int{0}}},
		},
		{
			name:   "a table SELECT INTO makes",
			sql:    "SELECT 1 AS id INTO n4;\nUPDATE n4 SET id = 2 WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE n4 SET id = 2 WHERE id = 1", "prior backup runs before the change, when n4 does not exist yet", []int{0}}},
		},
		{
			name:   "a table a CREATE SCHEMA makes",
			sql:    "CREATE SCHEMA s3 CREATE TABLE s3.n6 (id INT);\nUPDATE s3.n6 SET id = 2 WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE s3.n6 SET id = 2 WHERE id = 1", "prior backup runs before the change, when s3.n6 does not exist yet", []int{0}}},
		},
		{
			name:   "a synonym the change creates",
			sql:    "CREATE SYNONYM n5 FOR dbo.t;\nUPDATE n5 SET n = 2 WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE n5 SET n = 2 WHERE id = 1", "prior backup runs before the change, when n5 does not exist yet", []int{0}}},
		},
		{
			name:   "statements before a USE",
			sql:    "ALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE gone = 1;\nUPDATE n3 SET id = 1 WHERE id = 2;\nUSE other;\nCREATE TABLE n3 (id INT PRIMARY KEY);\nDELETE FROM nokey WHERE gone = 2;",
			change: backup,
			want:   []finding{{1, "DELETE FROM nokey WHERE gone = 1", "prior backup runs before the change, when column gone of nokey does not exist yet", []int{0}}},
		},
		{
			name:   "a table or column the change makes",
			sql:    "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;\nALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE gone = 1;\nALTER TABLE p ADD memo VARCHAR(10);\nUPDATE x SET x.n = 1 FROM t AS x JOIN p ON p.memo = x.code WHERE x.id = 3;",
			change: backup,
			want: []finding{
				{1, "UPDATE n2 SET id = 2 WHERE id = 1", "prior backup runs before the change, when n2 does not exist yet", []int{0}},
				{3, "DELETE FROM nokey WHERE gone = 1", "prior backup runs before the change, when column gone of nokey does not exist yet", []int{0}},
				{5, "UPDATE x SET x.n = 1 FROM t AS x JOIN p ON p.memo = x.code WHERE x.id = 3", "prior backup runs before the change, when column memo of p does not exist yet", []int{0}},
			},
		},
	})
}
