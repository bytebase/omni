package review

import (
	"context"
	"slices"
	"strings"
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
		if err := review.CheckCodes(f); err != nil {
			t.Error(err)
		}
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
		{
			name: "an empty statement before one",
			sql:  "SELECT 1;; SELEC 2;",
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
				{0, "UPDATE t SET n = 1", "`UPDATE` of `t` has no `WHERE` clause.", []int{0}},
				{1, "DELETE FROM dbo.t", "`DELETE` from `dbo.t` has no `WHERE` clause.", []int{0}},
				{3, "DELETE x FROM t AS x JOIN p ON x.p_id = p.id", "`DELETE` from `x` has no `WHERE` clause.", []int{0}},
			},
		},
		{
			name: "inside IF and a procedure",
			sql:  "IF 1 = 1 DELETE FROM t;\nGO\nCREATE PROCEDURE pr AS BEGIN UPDATE t SET n = 1 END",
			want: []finding{
				{0, "DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause.", []int{0}},
				{2, "UPDATE t SET n = 1", "`UPDATE` of `t` has no `WHERE` clause.", []int{0}},
			},
		},
		{
			name: "a table variable",
			sql:  "DECLARE @rows TABLE (n INT);\nDELETE FROM @rows;",
			want: []finding{{1, "DELETE FROM @rows", "`DELETE` from `@rows` has no `WHERE` clause.", []int{0}}},
		},
	})
}

func TestRequireIsNull(t *testing.T) {
	runCases(t, review.RequireIsNull, []ruleCase{
		{
			name: "comparisons a predicate tests",
			sql:  "SELECT * FROM t WHERE n = NULL OR NULL <> code;\nUPDATE t SET n = 1 WHERE NOT (code != NULL);\nIF @x = NULL PRINT 1;",
			want: []finding{
				{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{0, "NULL <> code", "`NULL <> code` is never true.\n\nUse `IS NOT NULL` instead.", []int{0}},
				{1, "code != NULL", "`code != NULL` is never true.\n\nUse `IS NOT NULL` instead.", []int{0}},
				{2, "@x = NULL", "`@x = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
			},
		},
		{
			name: "a comparison with a block comment holding dashes",
			sql:  "SELECT * FROM t WHERE n /* -- note */ = NULL;",
			want: []finding{{0, "n /* -- note */ = NULL", "`n /* -- note */ = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "a comparison with a line comment inside it",
			sql:  "SELECT * FROM t WHERE n -- note\n = NULL;",
			want: []finding{{0, "n -- note\n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "a bracketed name with an escaped bracket",
			sql:  "SELECT * FROM t WHERE [a]]  b] =\n NULL;",
			want: []finding{{0, "[a]]  b] =\n NULL", "`[a]]  b] = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "a value, IS NULL, and after SET ANSI_NULLS OFF",
			sql:  "SELECT CASE WHEN 1 = 1 THEN 1 END FROM t WHERE n IS NULL;\nSET ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;",
		},
		{
			name: "ANSI_NULLS set back ON",
			sql:  "SET ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_NULLS ON;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_DEFAULTS OFF;\nSELECT * FROM t WHERE n = NULL;",
			want: []finding{{3, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "ANSI_NULLS set in a block, under a condition, or in a routine body",
			sql:  "BEGIN SET ANSI_NULLS OFF; END;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_NULLS ON;\nIF @x = 1 SET ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_NULLS ON;\nCREATE PROCEDURE pr2 AS BEGIN SET ANSI_NULLS OFF; END;\nGO\nSELECT * FROM t WHERE n = NULL;",
			want: []finding{{8, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "ANSI_NULLS set in a list of options",
			sql:  "SET QUOTED_IDENTIFIER, ANSI_NULLS OFF;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_PADDING, ANSI_NULLS ON;\nSELECT * FROM t WHERE n = NULL;\nSET ANSI_PADDING, ANSI_DEFAULTS OFF;\nSELECT * FROM t WHERE n = NULL;",
			want: []finding{{3, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "ANSI_NULLS set back ON in a block after a condition",
			sql:  "BEGIN IF @x = 1 SET ANSI_NULLS OFF; SET ANSI_NULLS ON; END;\nSELECT * FROM t WHERE n = NULL;",
			want: []finding{{1, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "IIF",
			sql:  "SELECT IIF(n = NULL, 1, 0) FROM t;",
			want: []finding{{0, "n = NULL", "`n = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
		{
			name: "MERGE conditions",
			sql:  "MERGE t USING p ON t.id = NULL WHEN MATCHED AND p.note <> NULL THEN UPDATE SET n = 1;",
			want: []finding{
				{0, "t.id = NULL", "`t.id = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}},
				{0, "p.note <> NULL", "`p.note <> NULL` is never true.\n\nUse `IS NOT NULL` instead.", []int{0}},
			},
		},
		{
			name: "quoted text kept as written",
			sql:  "SELECT * FROM t WHERE [a  b] = NULL;",
			want: []finding{{0, "[a  b] = NULL", "`[a  b] = NULL` is never true.\n\nUse `IS NULL` instead.", []int{0}}},
		},
	})
}

func TestDisallowDropObject(t *testing.T) {
	runCases(t, review.DisallowDropObject, []ruleCase{
		{
			name: "a user, whose name is one part",
			sql:  "DROP USER [a. b];",
			want: []finding{{0, "DROP USER [a. b]", "Drops user `[a. b]`.", []int{0}}},
		},
		{
			name: "an external library, whose name is one part",
			sql:  "DROP EXTERNAL LIBRARY [a. b];",
			want: []finding{{0, "DROP EXTERNAL LIBRARY [a. b]", "Drops external library `[a. b]`.", []int{0}}},
		},
		{
			name: "Service Broker objects, whose names are one part",
			sql:  "DROP MESSAGE TYPE [//Adventure-Works.com/Expenses];\nDROP QUEUE dbo.q;",
			want: []finding{
				{0, "DROP MESSAGE TYPE [//Adventure-Works.com/Expenses]", "Drops message type `[//Adventure-Works.com/Expenses]`.", []int{0}},
				{1, "DROP QUEUE dbo.q", "Drops queue `dbo.q`.", []int{0}},
			},
		},
		{
			name: "indexes of the same name on two tables",
			sql:  "DROP INDEX ix ON dbo.t, ix ON dbo.p;\nDROP INDEX t.ix_n;",
			want: []finding{
				{0, "DROP INDEX ix ON dbo.t, ix ON dbo.p", "Drops indexes `ix` of `dbo.t` and `ix` of `dbo.p`.", []int{0}},
				{1, "DROP INDEX t.ix_n", "Drops index `t.ix_n`.", []int{0}},
			},
		},
		{
			name: "every DROP",
			sql:  "DROP TABLE IF EXISTS t, dbo.p;\nDROP VIEW v;\nDROP INDEX ix_n ON t;\nDROP PROCEDURE pr;\nDROP FUNCTION f;\nDROP SCHEMA s;\nDROP DATABASE d;\nDROP USER u;\nALTER TABLE t DROP COLUMN n, code;",
			want: []finding{
				{0, "DROP TABLE IF EXISTS t, dbo.p", "Drops tables `t` and `dbo.p`.", []int{0}},
				{1, "DROP VIEW v", "Drops view `v`.", []int{0}},
				{2, "DROP INDEX ix_n ON t", "Drops index `ix_n` of `t`.", []int{0}},
				{3, "DROP PROCEDURE pr", "Drops procedure `pr`.", []int{0}},
				{4, "DROP FUNCTION f", "Drops function `f`.", []int{0}},
				{5, "DROP SCHEMA s", "Drops schema `s`.", []int{0}},
				{6, "DROP DATABASE d", "Drops database `d`.", []int{0}},
				{7, "DROP USER u", "Drops user `u`.", []int{0}},
				{8, "n", "Drops column `n` of `t`.", []int{0}},
				{8, "code", "Drops column `code` of `t`.", []int{0}},
			},
		},
		{
			name: "a column after a constraint in one DROP",
			sql:  "ALTER TABLE t DROP CONSTRAINT uk_code, COLUMN n, code;",
			want: []finding{
				{0, "n", "Drops column `n` of `t`.", []int{0}},
				{0, "code", "Drops column `code` of `t`.", []int{0}},
			},
		},
		{name: "a constraint belongs to DisallowDropConstraint", sql: "ALTER TABLE t DROP CONSTRAINT ck_n;"},
		{
			name: "a name with a line break",
			sql:  "DROP TABLE [a\nb];",
			want: []finding{{0, "DROP TABLE [a\nb]", "Drops table `[a\\nb]`.", []int{0}}},
		},
		{
			name: "DROP SENSITIVITY CLASSIFICATION and DROP SIGNATURE",
			sql:  "DROP SENSITIVITY CLASSIFICATION FROM dbo.t.n;\nDROP SIGNATURE FROM pr BY CERTIFICATE c1;",
			want: []finding{
				{0, "DROP SENSITIVITY CLASSIFICATION FROM dbo.t.n", "Drops the sensitivity classification of `dbo.t.n`.", []int{0}},
				{1, "DROP SIGNATURE FROM pr BY CERTIFICATE c1", "Drops the signature of `pr`.", []int{0}},
			},
		},
		{
			name: "DROP COUNTER SIGNATURE and a key without a name",
			sql:  "DROP COUNTER SIGNATURE FROM pr BY CERTIFICATE c1;\nDROP MASTER KEY;",
			want: []finding{
				{0, "DROP COUNTER SIGNATURE FROM pr BY CERTIFICATE c1", "Drops the counter signature of `pr`.", []int{0}},
				{1, "DROP MASTER KEY", "Drops the master key.", []int{0}},
			},
		},
		{
			name: "a name with a part left out",
			sql:  "DROP TABLE shop..t;",
			want: []finding{{0, "DROP TABLE shop..t", "Drops table `shop..t`.", []int{0}}},
		},
		{
			name: "DROP STATISTICS names each",
			sql:  "DROP STATISTICS t.s1, dbo.t.s2;",
			want: []finding{{0, "DROP STATISTICS t.s1, dbo.t.s2", "Drops statistics `t.s1` and `dbo.t.s2`.", []int{0}}},
		},
		{
			name: "dotted and listed names of DROP",
			sql:  "DROP EXTERNAL TABLE dbo.t;\nDROP EVENT NOTIFICATION n1, n2 ON SERVER;\nDROP QUEUE dbo.q;",
			want: []finding{
				{0, "DROP EXTERNAL TABLE dbo.t", "Drops external table `dbo.t`.", []int{0}},
				{1, "DROP EVENT NOTIFICATION n1, n2 ON SERVER", "Drops event notifications `n1` and `n2`.", []int{0}},
				{2, "DROP QUEUE dbo.q", "Drops queue `dbo.q`.", []int{0}},
			},
		},
	})
}

func TestDisallowTruncate(t *testing.T) {
	runCases(t, review.DisallowTruncate, []ruleCase{
		{
			name: "TRUNCATE TABLE",
			sql:  "TRUNCATE TABLE t;\nTRUNCATE TABLE dbo.p WITH (PARTITIONS (1));",
			want: []finding{
				{0, "TRUNCATE TABLE t", "Truncates `t`.", []int{0}},
				{1, "TRUNCATE TABLE dbo.p WITH (PARTITIONS (1))", "Truncates partitions of `dbo.p`.", []int{0}},
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
				{0, "EXEC sp_rename 'dbo.t.n', 'm', 'COLUMN'", "Renames column `n` of `dbo.t` to `m`.", []int{0}},
				{1, "EXEC sp_rename @objname = N't', @newname = N't3'", "Renames table `t` to `t3`.", []int{0}},
				{2, "EXEC sys.sp_rename 'v', 'v2', 'OBJECT'", "Renames view `v` to `v2`.", []int{0}},
			},
		},
		{
			name: "sp_rename of a column whose name keeps spaces inside brackets",
			sql:  "EXEC sp_rename N'dbo.t.[ c ]', N'd', 'COLUMN';",
			want: []finding{{0, "EXEC sp_rename N'dbo.t.[ c ]', N'd', 'COLUMN'", "Renames column `[ c ]` of `dbo.t` to `d`.", []int{0}}},
		},
		{
			name: "sp_rename of a name the change renamed away and reused",
			sql:  "EXEC sp_rename 'dbo.t', 'moved', 'OBJECT';\nGO\nCREATE VIEW dbo.t AS SELECT 1 AS n;\nGO\nEXEC sp_rename 'dbo.t', 'new', 'OBJECT';",
			want: []finding{{0, "EXEC sp_rename 'dbo.t', 'moved', 'OBJECT'", "Renames table `dbo.t` to `moved`.", []int{0}}},
		},
		{
			name: "sp_rename of a name renamed away and reused in the same block",
			sql:  "BEGIN EXEC sp_rename 'dbo.v', 'moved', 'OBJECT'; CREATE TABLE dbo.v (id INT); EXEC sp_rename 'dbo.v', 'new', 'OBJECT'; END;",
			want: []finding{{0, "EXEC sp_rename 'dbo.v', 'moved', 'OBJECT'", "Renames view `dbo.v` to `moved`.", []int{0}}},
		},
		{
			name: "sp_rename of a name the change dropped and reused",
			sql:  "DROP TABLE dbo.t;\nGO\nCREATE PROCEDURE dbo.t AS SELECT 1;\nGO\nEXEC sp_rename 'dbo.t', 'new', 'OBJECT';",
		},
		{
			name: "sp_rename through another database",
			sql:  "EXEC other.sys.sp_rename 't', 'x';",
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
				{0, "pk_t", "Drops primary key `pk_t` of `t`.", []int{0}},
				{0, "fk_p", "Drops foreign key `fk_p` of `t`.", []int{0}},
				{1, "uk_code", "Drops unique constraint `uk_code` of `dbo.t`.", []int{0}},
				{2, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}},
			},
		},
		{
			name: "an item that restates CONSTRAINT or COLUMN",
			sql:  "ALTER TABLE t DROP CONSTRAINT pk_t, CONSTRAINT fk_p;\nALTER TABLE t DROP CONSTRAINT uk_code, COLUMN n;\nALTER TABLE t DROP COLUMN code, CONSTRAINT IF EXISTS ck_n, df_n;",
			want: []finding{
				{0, "pk_t", "Drops primary key `pk_t` of `t`.", []int{0}},
				{0, "fk_p", "Drops foreign key `fk_p` of `t`.", []int{0}},
				{1, "uk_code", "Drops unique constraint `uk_code` of `t`.", []int{0}},
				{2, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}},
			},
		},
		{
			name: "constraints the change defined",
			sql:  "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;\nALTER TABLE nokey ADD CONSTRAINT df_a DEFAULT 0 FOR a;\nALTER TABLE nokey DROP CONSTRAINT df_a;",
			want: []finding{{1, "ck_a", "Drops check constraint `ck_a` of `nokey`.", []int{0}}},
		},
		{
			name: "a constraint on a column the change adds, or of the database named",
			sql:  "ALTER TABLE nokey ADD extra INT CONSTRAINT ck_extra CHECK (extra >= 0);\nALTER TABLE nokey DROP CONSTRAINT ck_extra;\nALTER TABLE other.dbo.t ADD CONSTRAINT ck_y CHECK (n > 0);\nALTER TABLE dbo.t ADD CONSTRAINT ck_y CHECK (n > 0);\nALTER TABLE other.dbo.t DROP CONSTRAINT ck_y;",
			want: []finding{
				{1, "ck_extra", "Drops check constraint `ck_extra` of `nokey`.", []int{0}},
				{4, "ck_y", "Drops check constraint `ck_y` of `other.dbo.t`.", []int{0}},
			},
		},
		{
			name: "a constraint added and dropped in a block, or defined on a table of a named schema",
			sql:  "BEGIN ALTER TABLE nokey ADD CONSTRAINT ck_b CHECK (a > 0); ALTER TABLE nokey DROP CONSTRAINT ck_b; END;\nALTER TABLE s.t2 ADD CONSTRAINT ck_z CHECK (id > 0);\nALTER TABLE t2 DROP CONSTRAINT ck_z;",
			want: []finding{{0, "ck_b", "Drops check constraint `ck_b` of `nokey`.", []int{0}}},
		},
		{
			name: "a definition a rollback may undo",
			sql:  "BEGIN TRAN;\nALTER TABLE t ADD CONSTRAINT df_x CHECK (n >= 0);\nROLLBACK;\nALTER TABLE t DROP CONSTRAINT df_x;",
		},
		{
			name: "a constraint dropped, added back as another kind, and dropped again in control flow",
			sql:  "IF 1 = 1 BEGIN ALTER TABLE t DROP CONSTRAINT ck_n; ALTER TABLE t ADD CONSTRAINT ck_n DEFAULT 0 FOR n; ALTER TABLE t DROP CONSTRAINT ck_n; END;\nALTER TABLE t DROP CONSTRAINT ck_n;",
			want: []finding{{0, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}}},
		},
		{
			name: "constraints and tables the change renames",
			sql:  "CREATE TABLE dbo.a (n INT CONSTRAINT ck CHECK (n > 0));\nEXEC sp_rename 'dbo.a.ck', 'old_ck', 'OBJECT';\nCREATE TABLE dbo.b (n INT CONSTRAINT ck DEFAULT 0);\nEXEC sp_rename 'dbo.a', 'old_a', 'OBJECT';\nEXEC sp_rename 'dbo.b', 'a', 'OBJECT';\nALTER TABLE dbo.a DROP CONSTRAINT ck;",
		},
		{
			name: "a table the change made, moved away and replaced by a transfer",
			sql:  "CREATE TABLE dbo.a (n INT CONSTRAINT ck CHECK (n > 0));\nALTER SCHEMA archive TRANSFER dbo.a;\nALTER SCHEMA dbo TRANSFER staging.a;\nALTER TABLE dbo.a DROP CONSTRAINT ck;",
		},
		{
			name: "a table of the name moved in from another schema",
			sql:  "ALTER SCHEMA archive TRANSFER dbo.t;\nALTER SCHEMA dbo TRANSFER staging.t;\nALTER TABLE dbo.t DROP CONSTRAINT ck_n;",
		},
		{
			name: "a constraint the change renames to a synced constraint's name",
			sql:  "EXEC sp_rename 'dbo.t.ck_n', 'old_ck', 'OBJECT';\nEXEC sp_rename 'dbo.t.df_n', 'ck_n', 'OBJECT';\nALTER TABLE dbo.t DROP CONSTRAINT ck_n;",
		},
		{
			name: "a table recreated in control flow with a constraint of the name",
			sql:  "DROP TABLE dbo.t;\nIF 1 = 1 CREATE TABLE dbo.t (n INT CONSTRAINT ck_n DEFAULT 0);\nALTER TABLE dbo.t DROP CONSTRAINT ck_n;",
		},
		{
			name: "a constraint added back on a new column in control flow",
			sql:  "IF 1 = 1 BEGIN ALTER TABLE t DROP CONSTRAINT ck_n; ALTER TABLE t ADD n2 INT CONSTRAINT ck_n DEFAULT 0; ALTER TABLE t DROP CONSTRAINT ck_n; END;",
			want: []finding{{0, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}}},
		},
		{
			name: "a rollback to a savepoint, which keeps what came before it",
			sql:  "BEGIN TRAN;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t ADD CONSTRAINT ck_n DEFAULT 0 FOR n;\nSAVE TRAN keep;\nROLLBACK TRAN keep;\nALTER TABLE t DROP CONSTRAINT ck_n;",
			want: []finding{{1, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}}},
		},
		{
			name: "a rollback in a procedure body, which does not run",
			sql:  "ALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t ADD CONSTRAINT ck_n DEFAULT 0 FOR n;\nGO\nCREATE PROCEDURE pr AS ROLLBACK;\nGO\nALTER TABLE t DROP CONSTRAINT ck_n;",
			want: []finding{{0, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}}},
		},
		{
			name: "a definition a rollback in control flow may undo",
			sql:  "BEGIN TRAN;\nALTER TABLE t DROP CONSTRAINT df_n;\nALTER TABLE t ADD CONSTRAINT df_n CHECK (n >= 0);\nIF 1 = 1 ROLLBACK;\nALTER TABLE t DROP CONSTRAINT df_n;",
		},
		{
			name: "an unqualified name only a schema other than dbo holds",
			sql:  "ALTER TABLE t3 DROP CONSTRAINT ck_a;\nALTER TABLE audit.t3 DROP CONSTRAINT ck_a;",
			targets: []review.Target{{Schema: &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{
				{Name: "dbo"},
				{Name: "audit", Tables: []*metadata.TableMetadata{{Name: "t3", CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck_a"}}}}},
			}}}},
			want: []finding{{1, "ck_a", "Drops check constraint `ck_a` of `audit.t3`.", []int{0}}},
		},
		{
			name: "a USE in a block",
			sql:  "BEGIN USE other; END;\nALTER TABLE t DROP CONSTRAINT ck_n;",
		},
		{
			name: "a USE before a drop in the same block",
			sql:  "BEGIN USE other; ALTER TABLE t DROP CONSTRAINT ck_n; END;",
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
				{0, "fk_p", "Drops foreign key `fk_p` of `t`.", []int{0}},
				{2, "ck_n", "Drops check constraint `ck_n` of `t`.", []int{0}},
			},
		},
	})
}

func TestRequirePrimaryKey(t *testing.T) {
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name: "CREATE TABLE without a primary key",
			sql:  "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);\nCREATE TABLE c (id INT, CONSTRAINT pk_c PRIMARY KEY (id));\nCREATE TABLE #tmp (id INT);",
			want: []finding{{0, "CREATE TABLE a (id INT)", "Creates table `a` without a primary key.", []int{0}}},
		},
		{
			name: "settled by a later key, drop, or rename",
			sql:  "CREATE TABLE a (id INT NOT NULL);\nALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id);\nCREATE TABLE b (id INT);\nDROP TABLE b;\nCREATE TABLE c (id INT);\nEXEC sp_rename 'c', 'c2';\nCREATE DATABASE d;\nCREATE TABLE d.dbo.e (id INT);\nDROP DATABASE d;",
		},
		{
			name: "a column rename, which does not move a table named like it",
			sql:  "CREATE TABLE [t].[id] (x INT);\nEXEC sp_rename 't.id', 'id2', 'COLUMN';",
			want: []finding{{0, "CREATE TABLE [t].[id] (x INT)", "Creates table `t.id` without a primary key.", []int{0}}},
		},
		{
			name: "a primary key dropped and not added back",
			sql:  "ALTER TABLE t DROP CONSTRAINT pk_t;\nALTER TABLE p DROP CONSTRAINT pk_p;\nALTER TABLE p ADD CONSTRAINT pk_p2 PRIMARY KEY (id, note);",
			want: []finding{{0, "pk_t", "Removes the primary key of `t`, and the change doesn't add one back.", []int{0}}},
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
			want: []finding{{0, "CREATE TABLE s2.k (id INT)", "Creates table `s2.k` without a primary key.", []int{0}}},
		},
		{
			name: "not settled by a column rename, nor judged in a routine body",
			sql:  "CREATE TABLE f (id INT);\nRENAME OBJECT f COLUMN id TO id2;\nGO\nCREATE PROCEDURE pr AS ALTER TABLE t DROP CONSTRAINT pk_t;",
			want: []finding{{0, "CREATE TABLE f (id INT)", "Creates table `f` without a primary key.", []int{0}}},
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
			want:   []finding{{-1, "", "SQL Server doesn't support online migration.\n\nTurn off **Online migration** for this change.", []int{0}}},
		},
	})
}

func TestPriorBackup(t *testing.T) {
	backup := review.Change{PriorBackup: true, MaxBackupSize: 1 << 20}
	runCases(t, review.PriorBackup, []ruleCase{
		{
			name:   "a temporary target, whose backup Bytebase skips",
			sql:    "CREATE TABLE n2 (id INT);\nCREATE TABLE #tmp (id INT);\nUPDATE x SET id = n2.id FROM #tmp AS x JOIN n2 ON n2.id = x.id;",
			change: backup,
		},
		{
			name:   "a table the change makes in another database, or in place of a synonym it drops",
			sql:    "CREATE TABLE other.dbo.n (id INT);\nUPDATE n SET id = 2 WHERE id = 1;\nDROP SYNONYM n3;\nCREATE TABLE n3 (id INT);\nUPDATE n3 SET id = 2 WHERE id = 1;",
			change: backup,
		},
		{name: "a plain backup", sql: "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM p WHERE id = 2;", change: backup},
		{
			name:    "no bbdataarchive",
			sql:     "ALTER TABLE t ADD m INT;",
			change:  backup,
			targets: []review.Target{{Schema: shop()}},
			want:    []finding{{-1, "", "Prior backup needs the backup database `bbdataarchive`, which doesn't exist on this instance.\n\nCreate it, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "what the backup refuses",
			sql:    "DECLARE @v TABLE (id INT);\nUPDATE @v SET id = 1 WHERE id = 2;\nWITH c AS (SELECT * FROM t) DELETE FROM c WHERE id = 1;\nUPDATE TOP (3) p SET note = '' WHERE id > 0;",
			change: backup,
			want: []finding{
				{1, "UPDATE @v SET id = 1 WHERE id = 2", "Prior backup can't back up the table variable `@v`.\n\nTurn off **Prior backup** for this change.", []int{0}},
				{2, "WITH c AS (SELECT * FROM t) DELETE FROM c WHERE id = 1", "Prior backup can't back up through the `WITH` query `c`.\n\nChange the base table instead, or turn off **Prior backup**.", []int{0}},
				{3, "UPDATE TOP (3) p SET note = '' WHERE id > 0", "Prior backup can't copy this statement's `TOP`: SQL Server rejects it where the backup query puts it.\n\nUse a `WHERE` instead of `TOP`, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a change over the size limit",
			sql:    strings.Repeat("UPDATE t SET n = 1 WHERE id = 1;\n", 32),
			change: review.Change{PriorBackup: true, MaxBackupSize: 1024},
			want:   []finding{{-1, "", "The change is larger than prior backup's 1 KiB limit.\n\nSplit it into smaller changes, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "statement forms the backup refuses",
			sql:    "DELETE FROM t WHERE CURRENT OF cur;\nWITH XMLNAMESPACES ('urn:x' AS x) UPDATE t SET n = 1 WHERE id = 1;",
			change: backup,
			want: []finding{
				{0, "DELETE FROM t WHERE CURRENT OF cur", "Prior backup can't back up a statement with `WHERE CURRENT OF`.\n\nUse a `WHERE` on the table's key, or turn off **Prior backup**.", []int{0}},
				{1, "WITH XMLNAMESPACES ('urn:x' AS x) UPDATE t SET n = 1 WHERE id = 1", "Prior backup can't back up a statement with `WITH XMLNAMESPACES`.\n\nTurn off **Prior backup** for this change.", []int{0}},
			},
		},
		{
			name:   "UPDATE and DELETE on one table, and different WITH clauses",
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM dbo.t WHERE id = 2;\nWITH a AS (SELECT 1 AS id) UPDATE p SET note = '' WHERE id IN (SELECT id FROM a);\nUPDATE p SET note = 'x' WHERE id = 2;",
			change: backup,
			want: []finding{
				{1, "DELETE FROM dbo.t WHERE id = 2", "Prior backup can't back up `dbo.t`: this change both updates and deletes its rows.\n\nMove the `UPDATE` and the `DELETE` into separate changes, or turn off **Prior backup**.", []int{0}},
				{3, "UPDATE p SET note = 'x' WHERE id = 2", "Prior backup can't back up `p`: its statements have different `WITH` clauses.\n\nGive them the same `WITH` clause, or move them into separate changes.", []int{0}},
			},
		},
		{
			name:   "a name without a database, the target's",
			sql:    "DELETE FROM shop.dbo.p WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nDELETE FROM other.dbo.nokey WHERE a = 1;\nUPDATE nokey SET a = 2 WHERE a = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE p SET note = '' WHERE id = 1", "Prior backup can't back up `p`: this change both updates and deletes its rows.\n\nMove the `UPDATE` and the `DELETE` into separate changes, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a table SELECT INTO makes",
			sql:    "SELECT 1 AS id INTO n4;\nUPDATE n4 SET id = 2 WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE n4 SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `n4` doesn't exist yet.\n\nCreate `n4` in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a column sp_rename renames to, and a FROM list without the target",
			sql:    "EXEC sp_rename 'dbo.t.n', 'm', 'COLUMN';\nUPDATE dbo.t SET code = '' WHERE m = 1;\nUPDATE dbo.t SET n = p.id FROM dbo.p AS p;\nUPDATE t SET n = 1 FROM t JOIN p ON p.id = t.p_id;",
			change: backup,
			want: []finding{
				{1, "UPDATE dbo.t SET code = '' WHERE m = 1", "Prior backup runs before the change, when column `m` of `dbo.t` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
				{2, "UPDATE dbo.t SET n = p.id FROM dbo.p AS p", "Prior backup reads `dbo.t` from the statement's `FROM` list, which doesn't name it.\n\nName `dbo.t` in the `FROM` list, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a table sp_rename renames to, and FROM lists by alias",
			sql:    "EXEC sp_rename 'dbo.nokey', 'n8';\nUPDATE dbo.n8 SET a = 2 WHERE a = 1;\nUPDATE dbo.t SET n = s.n FROM audit.t AS s;\nUPDATE x SET n = 1 FROM t AS x;",
			change: backup,
			want: []finding{
				{1, "UPDATE dbo.n8 SET a = 2 WHERE a = 1", "Prior backup runs before the change, when `dbo.n8` doesn't exist yet.\n\nCreate `dbo.n8` in an earlier change, or turn off **Prior backup**.", []int{0}},
				{2, "UPDATE dbo.t SET n = s.n FROM audit.t AS s", "Prior backup reads `dbo.t` from the statement's `FROM` list, which doesn't name it.\n\nName `dbo.t` in the `FROM` list, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a table and a column RENAME OBJECT renames to",
			sql:    "RENAME OBJECT dbo.nokey TO n9;\nUPDATE dbo.n9 SET a = 2 WHERE a = 1;\nRENAME OBJECT dbo.t COLUMN n TO m2;\nUPDATE dbo.t SET code = '' WHERE m2 = 1;",
			change: backup,
			want: []finding{
				{1, "UPDATE dbo.n9 SET a = 2 WHERE a = 1", "Prior backup runs before the change, when `dbo.n9` doesn't exist yet.\n\nCreate `dbo.n9` in an earlier change, or turn off **Prior backup**.", []int{0}},
				{3, "UPDATE dbo.t SET code = '' WHERE m2 = 1", "Prior backup runs before the change, when column `m2` of `dbo.t` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "tables CTAS and CLONE make, and an alias that cannot bind a qualified target",
			sql:    "CREATE TABLE n10 WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT 1 AS id;\nUPDATE n10 SET id = 2 WHERE id = 1;\nUPDATE dbo.t SET n = t.id FROM dbo.p AS t;",
			change: backup,
			want: []finding{
				{1, "UPDATE n10 SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `n10` doesn't exist yet.\n\nCreate `n10` in an earlier change, or turn off **Prior backup**.", []int{0}},
				{2, "UPDATE dbo.t SET n = t.id FROM dbo.p AS t", "Prior backup reads `dbo.t` from the statement's `FROM` list, which doesn't name it.\n\nName `dbo.t` in the `FROM` list, or turn off **Prior backup**.", []int{0}},
			},
		},
		{
			name:   "a table a BEGIN block makes",
			sql:    "BEGIN CREATE TABLE n7 (id INT PRIMARY KEY); END;\nUPDATE n7 SET id = 2 WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE n7 SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `n7` doesn't exist yet.\n\nCreate `n7` in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a table a CREATE SCHEMA makes",
			sql:    "CREATE SCHEMA s3 CREATE TABLE s3.n6 (id INT);\nUPDATE s3.n6 SET id = 2 WHERE id = 1;",
			change: backup,
			want:   []finding{{1, "UPDATE s3.n6 SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `s3.n6` doesn't exist yet.\n\nCreate `s3.n6` in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a synonym, which the synced schema does not tell of",
			sql:    "DROP SYNONYM n5;\nCREATE SYNONYM n5 FOR dbo.t;\nUPDATE n5 SET n = 2 WHERE id = 1;",
			change: backup,
		},
		{
			name:   "statements before a USE",
			sql:    "ALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE gone = 1;\nUPDATE n3 SET id = 1 WHERE id = 2;\nUSE other;\nCREATE TABLE n3 (id INT PRIMARY KEY);\nDELETE FROM nokey WHERE gone = 2;",
			change: backup,
			want:   []finding{{1, "DELETE FROM nokey WHERE gone = 1", "Prior backup runs before the change, when column `gone` of `nokey` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}}},
		},
		{
			name:   "a table or column the change makes",
			sql:    "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;\nALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE gone = 1;\nALTER TABLE p ADD memo VARCHAR(10);\nUPDATE x SET x.n = 1 FROM t AS x JOIN p ON p.memo = x.code WHERE x.id = 3;",
			change: backup,
			want: []finding{
				{1, "UPDATE n2 SET id = 2 WHERE id = 1", "Prior backup runs before the change, when `n2` doesn't exist yet.\n\nCreate `n2` in an earlier change, or turn off **Prior backup**.", []int{0}},
				{3, "DELETE FROM nokey WHERE gone = 1", "Prior backup runs before the change, when column `gone` of `nokey` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
				{5, "UPDATE x SET x.n = 1 FROM t AS x JOIN p ON p.memo = x.code WHERE x.id = 3", "Prior backup runs before the change, when column `memo` of `p` doesn't exist yet.\n\nAdd the column in an earlier change, or turn off **Prior backup**.", []int{0}},
			},
		},
	})
}

func TestStatementRanges(t *testing.T) {
	sql := "SELECT 1;; UPDATE t SET n = 1 WHERE id = 1;\n;\nDELETE FROM t WHERE id = 2"
	result, err := Review(context.Background(), sql, review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, rng := range result.Statements {
		got = append(got, sql[rng.Start:rng.End])
	}
	want := []string{"SELECT 1", "UPDATE t SET n = 1 WHERE id = 1", "DELETE FROM t WHERE id = 2"}
	if !slices.Equal(got, want) {
		t.Errorf("statements = %q, want %q", got, want)
	}
}

// TestFindingCodes pins the code of every condition the package reports:
// a case's findings, in order, are all of its rule, and each carries
// exactly the one code listed for it.
func TestFindingCodes(t *testing.T) {
	backup := review.Change{PriorBackup: true, MaxBackupSize: 1 << 20}
	walks := []review.Target{walkTarget()}
	repeat := func(code review.FindingCode, n int) []review.FindingCode {
		return slices.Repeat([]review.FindingCode{code}, n)
	}
	tests := []struct {
		name    string
		rule    review.Rule
		sql     string
		change  review.Change
		targets []review.Target // nil: shop, with its backup database
		want    []review.FindingCode
	}{
		{name: "syntax error", rule: review.Syntax, sql: "SELECT 1;\nSELEC 2;", want: []review.FindingCode{review.SyntaxError}},

		{name: "online migration", rule: review.OnlineMigration, sql: "ALTER TABLE t ADD m INT;", change: review.Change{OnlineMigration: true}, want: []review.FindingCode{review.OnlineMigrationUnsupported}},

		{name: "comparison with NULL", rule: review.RequireIsNull, sql: "SELECT * FROM t WHERE n = NULL OR n <> NULL;", want: repeat(review.NullComparison, 2)},

		{name: "no WHERE", rule: review.RequireWhere, sql: "UPDATE t SET n = 1;\nDELETE FROM t;", want: repeat(review.MissingWhere, 2)},

		{
			name: "drops of data",
			rule: review.DisallowDropObject,
			sql:  "DROP TABLE t;\nDROP SCHEMA s;\nDROP DATABASE d;\nDROP SEQUENCE q;\nALTER TABLE t DROP COLUMN n;",
			want: repeat(review.DropsData, 5),
		},
		{
			name: "drops of definitions",
			rule: review.DisallowDropObject,
			sql:  "DROP VIEW v;\nDROP INDEX ix_n ON t;\nDROP PROCEDURE pr;\nDROP TYPE ty;\nDROP SYNONYM sy;\nDROP STATISTICS t.s1;\nDROP ASSEMBLY a;\nDROP PARTITION FUNCTION pf;\nDROP PARTITION SCHEME ps;\nDROP DEFAULT d;\nDROP RULE r;\nDROP QUEUE dbo.q;",
			want: repeat(review.DropsDefinition, 12),
		},
		{
			name: "drops of access",
			rule: review.DisallowDropObject,
			sql:  "DROP USER u;\nDROP LOGIN l;\nDROP ROLE r;\nDROP APPLICATION ROLE ar;\nDROP SYMMETRIC KEY k;\nDROP MASTER KEY;\nDROP CERTIFICATE c;\nDROP CREDENTIAL cr;\nDROP DATABASE SCOPED CREDENTIAL dcr;\nDROP SERVER AUDIT a;\nDROP DATABASE AUDIT SPECIFICATION das;\nDROP SECURITY POLICY sp;\nDROP SIGNATURE FROM pr BY CERTIFICATE c1;\nDROP COUNTER SIGNATURE FROM pr BY CERTIFICATE c1;\nDROP SENSITIVITY CLASSIFICATION FROM dbo.t.n;",
			want: repeat(review.DropsAccess, 15),
		},

		{name: "truncate", rule: review.DisallowTruncate, sql: "TRUNCATE TABLE t;\nTRUNCATE TABLE t WITH (PARTITIONS (1));", want: repeat(review.TruncatesData, 2)},

		{
			name: "drops of a constraint the change made and a synced one",
			rule: review.DisallowDropConstraint,
			sql:  "CREATE TABLE x (a INT CONSTRAINT ck_x CHECK (a > 0));\nALTER TABLE x DROP CONSTRAINT ck_x;\nALTER TABLE t DROP CONSTRAINT ck_n;",
			want: repeat(review.DropsConstraint, 2),
		},

		{
			name: "renames",
			rule: review.DisallowRename,
			sql:  "EXEC sp_rename 'dbo.t.n', 'm', 'COLUMN';\nRENAME OBJECT dbo.p TO p2;\nEXEC sp_rename 'dbo.v', 'v2';",
			want: repeat(review.RenamesObject, 3),
		},

		{name: "a table without a primary key", rule: review.RequirePrimaryKey, sql: "CREATE TABLE x (a INT);", want: []review.FindingCode{review.TableWithoutPrimaryKey}},
		{name: "a primary key removed", rule: review.RequirePrimaryKey, sql: "ALTER TABLE t DROP CONSTRAINT pk_t;", want: []review.FindingCode{review.PrimaryKeyRemoved}},

		{
			name:   "a change over the backup size",
			rule:   review.PriorBackup,
			sql:    "UPDATE t SET n = 1 WHERE id = 1;",
			change: review.Change{PriorBackup: true, MaxBackupSize: 8},
			want:   []review.FindingCode{review.ChangeTooLarge},
		},
		{
			name:    "no backup database",
			rule:    review.PriorBackup,
			sql:     "UPDATE t SET n = 1 WHERE id = 1;",
			change:  backup,
			targets: []review.Target{{Schema: shop()}},
			want:    []review.FindingCode{review.BackupLocationMissing},
		},
		{
			name:   "statements that cannot share a backup",
			rule:   review.PriorBackup,
			sql:    "UPDATE t SET n = 1 WHERE id = 1;\nDELETE FROM dbo.t WHERE id = 2;\nWITH a AS (SELECT 1 AS id) UPDATE p SET note = '' WHERE id IN (SELECT id FROM a);\nUPDATE p SET note = 'x' WHERE id = 2;",
			change: backup,
			want:   repeat(review.StatementsShareBackup, 2),
		},
		{
			name:   "a FROM list without the target",
			rule:   review.PriorBackup,
			sql:    "UPDATE dbo.t SET n = p.id FROM dbo.p AS p;",
			change: backup,
			want:   []review.FindingCode{review.BackupReadsWrongTable},
		},
		{
			name:   "targets the backup cannot copy",
			rule:   review.PriorBackup,
			sql:    "DECLARE @v TABLE (id INT);\nUPDATE @v SET id = 1 WHERE id = 2;\nWITH c AS (SELECT * FROM t) DELETE FROM c WHERE id = 1;",
			change: backup,
			want:   repeat(review.TableNotBackupable, 2),
		},
		{
			name:   "statement forms the backup cannot copy",
			rule:   review.PriorBackup,
			sql:    "UPDATE TOP (3) p SET note = '' WHERE id > 0;\nDELETE FROM t WHERE CURRENT OF cur;\nWITH XMLNAMESPACES ('urn:x' AS x) UPDATE t SET n = 1 WHERE id = 1;",
			change: backup,
			want:   repeat(review.StatementNotBackupable, 3),
		},
		{
			name:   "a table and a column the change makes",
			rule:   review.PriorBackup,
			sql:    "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;\nALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE gone = 1;",
			change: backup,
			want:   repeat(review.BackupReadsNewObject, 2),
		},

		{name: "a missing table", rule: review.WalkThrough, sql: "DROP TABLE dbo.nosuch;", targets: walks, want: []review.FindingCode{review.ObjectNotFound}},
		{name: "a missing column", rule: review.WalkThrough, sql: "ALTER TABLE dbo.t DROP COLUMN nosuch;", targets: walks, want: []review.FindingCode{review.ObjectNotFound}},
		{name: "a missing schema", rule: review.WalkThrough, sql: "CREATE TABLE nosuch.x (a INT);", targets: walks, want: []review.FindingCode{review.ObjectNotFound}},
		{name: "an index on a missing table", rule: review.WalkThrough, sql: "DROP INDEX ix ON dbo.nosuch;", targets: walks, want: []review.FindingCode{review.ObjectNotFound}},
		{name: "a missing constraint", rule: review.WalkThrough, sql: "ALTER TABLE dbo.t DROP CONSTRAINT nosuch;", targets: walks, want: []review.FindingCode{review.ObjectNotFound}},
		{name: "a taken name", rule: review.WalkThrough, sql: "CREATE TABLE dbo.t (a INT);", targets: walks, want: []review.FindingCode{review.NameTaken}},
		{name: "a taken index name", rule: review.WalkThrough, sql: "CREATE INDEX ix_a ON dbo.t (b);", targets: walks, want: []review.FindingCode{review.NameTaken}},
		{name: "a second primary key", rule: review.WalkThrough, sql: "ALTER TABLE dbo.t ADD CONSTRAINT pk2 PRIMARY KEY (id);", targets: walks, want: []review.FindingCode{review.NameTaken}},
		{name: "two primary keys in one statement", rule: review.WalkThrough, sql: "CREATE TABLE dbo.x (a INT PRIMARY KEY, b INT PRIMARY KEY);", targets: walks, want: []review.FindingCode{review.NameTaken}},
		{name: "another kind", rule: review.WalkThrough, sql: "DROP TABLE dbo.v;", targets: walks, want: []review.FindingCode{review.WrongObjectKind}},
		{name: "a column named twice", rule: review.WalkThrough, sql: "CREATE TABLE dbo.x (a INT, a INT);", targets: walks, want: []review.FindingCode{review.DuplicateName}},
		{name: "a constraint named like its table", rule: review.WalkThrough, sql: "CREATE TABLE dbo.x (a INT CONSTRAINT x CHECK (a > 0));", targets: walks, want: []review.FindingCode{review.DuplicateName}},
		{name: "a table dropped twice", rule: review.WalkThrough, sql: "DROP TABLE dbo.r, dbo.r;", targets: walks, want: []review.FindingCode{review.DuplicateName}},
		{name: "a name too long", rule: review.WalkThrough, sql: "CREATE TABLE dbo." + strings.Repeat("a", 129) + " (a INT);", targets: walks, want: []review.FindingCode{review.NameTooLong}},
		{name: "a schema that holds objects", rule: review.WalkThrough, sql: "DROP SCHEMA s;", targets: walks, want: []review.FindingCode{review.SchemaNotEmpty}},
		{name: "a view not first in its batch", rule: review.WalkThrough, sql: "CREATE TABLE dbo.x (a INT);\nCREATE VIEW dbo.w AS SELECT 1 AS a;", targets: walks, want: []review.FindingCode{review.NotFirstInBatch}},
		{name: "a trigger in another schema", rule: review.WalkThrough, sql: "CREATE TRIGGER s.tr ON dbo.t AFTER INSERT AS SELECT 1;", targets: walks, want: []review.FindingCode{review.TriggerSchemaMismatch}},
		{name: "an index that backs a constraint", rule: review.WalkThrough, sql: "DROP INDEX pk_t ON dbo.t;", targets: walks, want: []review.FindingCode{review.IndexBacksConstraint}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targets := tt.targets
			if targets == nil {
				targets = []review.Target{{Schema: shop(), BackupDatabaseExists: true}}
			}
			result, err := Review(context.Background(), tt.sql, review.Options{Rules: []review.Rule{tt.rule}, Change: tt.change}, targets)
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			var got [][]review.FindingCode
			for _, f := range result.Findings {
				if err := review.CheckCodes(f); err != nil {
					t.Error(err)
				}
				if f.Rule != tt.rule {
					t.Errorf("finding of %s, want %s: %s", f.Rule, tt.rule, f.Message)
				}
				got = append(got, f.Codes)
			}
			var want [][]review.FindingCode
			for _, c := range tt.want {
				want = append(want, []review.FindingCode{c})
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("codes mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
