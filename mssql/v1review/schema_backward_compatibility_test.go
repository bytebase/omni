package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestSchemaBackwardCompatibility(t *testing.T) {
	f := func(code v1review.Code, name string, stmt int, text string) anchoredFinding {
		return anchoredFinding{Code: code, Name: name, Statement: stmt, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		db   string
		want []anchoredFinding
	}{
		{
			name: "drop table keyed with defaults, case kept",
			sql:  "DROP TABLE Foo, s.Bar, d.s.Baz, srv.d.s.Qux;",
			db:   "MyDB",
			want: []anchoredFinding{
				f(v1review.IncompatibleDropTable, "MyDB.dbo.Foo", 0, "DROP TABLE Foo, s.Bar, d.s.Baz, srv.d.s.Qux"),
				f(v1review.IncompatibleDropTable, "MyDB.s.Bar", 0, "DROP TABLE Foo, s.Bar, d.s.Baz, srv.d.s.Qux"),
				f(v1review.IncompatibleDropTable, "d.s.Baz", 0, "DROP TABLE Foo, s.Bar, d.s.Baz, srv.d.s.Qux"),
				f(v1review.IncompatibleDropTable, "d.s.Qux", 0, "DROP TABLE Foo, s.Bar, d.s.Baz, srv.d.s.Qux"),
			},
		},
		{
			name: "empty current database left out",
			sql:  "DROP TABLE Foo;",
			want: []anchoredFinding{f(v1review.IncompatibleDropTable, "dbo.Foo", 0, "DROP TABLE Foo")},
		},
		{
			name: "table created earlier dropped freely, then forgotten",
			sql:  "CREATE TABLE t(a INT);\nDROP TABLE dbo.t;\nDROP TABLE t;",
			db:   "db",
			want: []anchoredFinding{f(v1review.IncompatibleDropTable, "db.dbo.t", 2, "DROP TABLE t")},
		},
		{
			name: "table key is case sensitive",
			sql:  "CREATE TABLE t(a INT);\nDROP TABLE T;",
			db:   "db",
			want: []anchoredFinding{f(v1review.IncompatibleDropTable, "db.dbo.T", 1, "DROP TABLE T")},
		},
		{
			name: "schema and database keys lower-cased",
			sql:  "CREATE SCHEMA S1;\nDROP SCHEMA s1;\nDROP SCHEMA S2;\nCREATE DATABASE D1;\nDROP DATABASE d1;\nDROP DATABASE D2;",
			db:   "Db",
			want: []anchoredFinding{
				f(v1review.IncompatibleDropSchema, "Db.s2", 2, "DROP SCHEMA S2"),
				f(v1review.IncompatibleDropDatabase, "d2", 5, "DROP DATABASE D2"),
			},
		},
		{
			name: "schema created by authorization alone",
			sql:  "CREATE SCHEMA AUTHORIZATION Owner1;\nDROP SCHEMA owner1;",
			db:   "db",
		},
		{
			name: "alter table actions, each dropped column its own action",
			sql:  "ALTER TABLE t DROP COLUMN A, B;\nALTER TABLE t ALTER COLUMN C INT;\nALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a);",
			db:   "db",
			want: []anchoredFinding{
				f(v1review.IncompatibleDropColumn, "a", 0, "A"),
				f(v1review.IncompatibleDropColumn, "b", 0, "B"),
				f(v1review.IncompatibleAlterColumn, "c", 1, "ALTER COLUMN C INT"),
				f(v1review.IncompatibleAddPrimaryKey, "", 2, "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a)"),
			},
		},
		{
			name: "added constraints by type",
			sql:  "ALTER TABLE t ADD CONSTRAINT u UNIQUE (a), CONSTRAINT c CHECK (a > 0), CONSTRAINT fk FOREIGN KEY (a) REFERENCES r(a), CONSTRAINT d DEFAULT 1 FOR a;",
			db:   "db",
			want: []anchoredFinding{
				f(v1review.IncompatibleAddUniqueKey, "", 0, "ALTER TABLE t ADD CONSTRAINT u UNIQUE (a), CONSTRAINT c CHECK (a > 0), CONSTRAINT fk FOREIGN KEY (a) REFERENCES r(a), CONSTRAINT d DEFAULT 1 FOR a"),
				f(v1review.IncompatibleAddCheck, "", 0, "ALTER TABLE t ADD CONSTRAINT u UNIQUE (a), CONSTRAINT c CHECK (a > 0), CONSTRAINT fk FOREIGN KEY (a) REFERENCES r(a), CONSTRAINT d DEFAULT 1 FOR a"),
				f(v1review.IncompatibleAddForeignKey, "", 0, "ALTER TABLE t ADD CONSTRAINT u UNIQUE (a), CONSTRAINT c CHECK (a > 0), CONSTRAINT fk FOREIGN KEY (a) REFERENCES r(a), CONSTRAINT d DEFAULT 1 FOR a"),
			},
		},
		{
			name: "with nocheck anchors the action, with check not reported",
			sql:  "ALTER TABLE t WITH NOCHECK ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES r(a);\nALTER TABLE t WITH NOCHECK ADD CONSTRAINT c CHECK (a > 0);\nALTER TABLE t WITH CHECK ADD CONSTRAINT c2 CHECK (a > 0);\nALTER TABLE t WITH NOCHECK ADD CONSTRAINT u UNIQUE (a);",
			db:   "db",
			want: []anchoredFinding{
				f(v1review.IncompatibleAddForeignKeyNoCheck, "", 0, "CONSTRAINT fk FOREIGN KEY (a) REFERENCES r(a)"),
				f(v1review.IncompatibleAddCheckNoCheck, "", 1, "CONSTRAINT c CHECK (a > 0)"),
			},
		},
		{
			name: "alter of a table created earlier not judged",
			sql:  "CREATE TABLE db.dbo.t(a INT);\nALTER TABLE t DROP COLUMN a;\nALTER TABLE T DROP COLUMN a;",
			db:   "db",
			want: []anchoredFinding{f(v1review.IncompatibleDropColumn, "a", 2, "a")},
		},
		{
			name: "sp_rename with a string first argument",
			sql:  "EXEC sp_rename 'a.b', 'c';\nEXEC SP_RENAME 'a', 'c';\nEXEC dbo.sp_rename 'a', 'c';\nEXEC sp_rename @old, 'c';\nEXEC sp_rename;",
			db:   "db",
			want: []anchoredFinding{
				f(v1review.IncompatibleRenameTable, "", 0, "EXEC sp_rename 'a.b', 'c'"),
				f(v1review.IncompatibleRenameTable, "", 1, "EXEC SP_RENAME 'a', 'c'"),
			},
		},
		{name: "nocheck constraint by name carries no constraint", sql: "ALTER TABLE t NOCHECK CONSTRAINT fk;", db: "db"},
		{name: "drop view not judged", sql: "DROP VIEW v; DROP INDEX ix ON t;", db: "db"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p AS DROP TABLE t;", db: "db"},
		{name: "if statement not judged", sql: "IF 1 = 1 DROP TABLE t;", db: "db"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			got := anchored(t, tt.sql, SchemaBackwardCompatibility(list.Items, tt.db))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// The parser gives each dropped column its own action, so the Names list a
// DROP COLUMN action can carry is built by hand.
func TestSchemaBackwardCompatibilityDropColumnNames(t *testing.T) {
	stmt := &ast.AlterTableStmt{
		Name: &ast.TableRef{Object: "t"},
		Actions: &ast.List{Items: []ast.Node{&ast.AlterTableAction{
			Type:    ast.ATDropColumn,
			ColName: "A",
			Names:   &ast.List{Items: []ast.Node{&ast.TableRef{Object: "B"}, &ast.String{Str: "skipped"}, &ast.TableRef{Object: "C"}}},
			Loc:     ast.Loc{Start: 3, End: 5},
		}}},
		Loc: ast.Loc{Start: 0, End: 9},
	}
	got := SchemaBackwardCompatibility([]ast.Node{stmt}, "db")
	want := []v1review.Finding{{Code: v1review.IncompatibleDropColumn, Name: "a, b, c", Start: 3, End: 5}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
