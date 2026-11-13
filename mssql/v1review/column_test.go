package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

// columnFinding is a column finding as the tests compare it, with the text
// its range anchors in place of the range.
type columnFinding struct {
	Code  v1review.Code
	Name  string
	Table string
	Value string
	Text  string
}

func runColumn(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []columnFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []columnFinding
	for _, f := range rule(list.Items) {
		got = append(got, columnFinding{f.Code, f.Name, f.Table, f.Value, sql[f.Start:f.End]})
	}
	return got
}

func TestColumnMaximumVarcharLength(t *testing.T) {
	long := func(text string) columnFinding {
		return columnFinding{Code: v1review.ColumnVarcharTooLong, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "within the limit", sql: "CREATE TABLE t(a VARCHAR(100), b NCHAR(100));"},
		{
			name: "every character type, MAX counted as the largest",
			sql:  "CREATE TABLE t(a VARCHAR(101), b NVARCHAR(MAX), c char(200), d NChar(101), e VARBINARY(MAX), f TEXT);",
			want: []columnFinding{long("VARCHAR(101)"), long("NVARCHAR(MAX)"), long("char(200)"), long("NChar(101)")},
		},
		{name: "no length counts as 0", sql: "CREATE TABLE t(a VARCHAR);"},
		{
			name: "nested data types are judged",
			sql:  "SELECT CAST(x AS VARCHAR(500)) FROM t;\nDECLARE @v NVARCHAR(500);",
			want: []columnFinding{long("VARCHAR(500)"), long("NVARCHAR(500)")},
		},
		{
			name: "alter column and add column",
			sql:  "ALTER TABLE t ALTER COLUMN a VARCHAR(300);\nALTER TABLE t ADD b CHAR(300);",
			want: []columnFinding{long("VARCHAR(300)"), long("CHAR(300)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumn(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return ColumnMaximumVarcharLength(stmts, 100) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestColumnNoNull(t *testing.T) {
	nullable := func(name, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnNullable, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{
			name: "not null, inline primary key, table primary key in any case",
			sql:  "CREATE TABLE t(Id INT, CONSTRAINT pk PRIMARY KEY (ID), a INT NOT NULL, b INT PRIMARY KEY, c INT CONSTRAINT nn NOT NULL);",
		},
		{
			name: "nullable columns, name lower-cased",
			sql:  "CREATE TABLE t(\n  Name INT,\n  b INT NULL UNIQUE\n);",
			want: []columnFinding{nullable("name", "Name INT"), nullable("b", "b INT NULL UNIQUE")},
		},
		{
			name: "alter column without NOT NULL and add column",
			sql:  "ALTER TABLE t ALTER COLUMN a INT;\nALTER TABLE t ADD b INT, c INT NOT NULL;",
			want: []columnFinding{nullable("a", "ALTER COLUMN a INT"), nullable("b", "b INT")},
		},
		{name: "a table constraint added by ALTER is not judged", sql: "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a);"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p AS BEGIN CREATE TABLE t(a INT) END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumn(t, tt.sql, ColumnNoNull)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestColumnRequired(t *testing.T) {
	missing := func(name, table, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnMissing, Name: name, Table: table, Text: text}
	}
	required := []string{"id", "Created_TS", "id"}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "create with every column in any case", sql: "CREATE TABLE t(ID INT, created_ts INT);"},
		{
			name: "create lacks columns, named as given and sorted",
			sql:  "CREATE TABLE db.dbo.t(a INT);",
			want: []columnFinding{
				missing("Created_TS", "db.dbo.t", "CREATE TABLE db.dbo.t(a INT)"),
				missing("id", "db.dbo.t", "CREATE TABLE db.dbo.t(a INT)"),
			},
		},
		{
			name: "drop column, lower-cased, matched exactly",
			sql:  "ALTER TABLE dbo.t DROP COLUMN ID, Created_TS, id;",
			want: []columnFinding{
				missing("id", "dbo.t", "ALTER TABLE dbo.t DROP COLUMN ID, Created_TS, id"),
				missing("id", "dbo.t", "ALTER TABLE dbo.t DROP COLUMN ID, Created_TS, id"),
			},
		},
		{name: "add and alter are not judged", sql: "ALTER TABLE t ADD x INT;\nALTER TABLE t ALTER COLUMN id BIGINT;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumn(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return ColumnRequired(stmts, required) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestColumnTypeDisallowList(t *testing.T) {
	disallowed := func(value, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnTypeDisallowed, Value: value, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "allowed types", sql: "CREATE TABLE t(a INT, b NVARCHAR(10));"},
		{
			name: "any case, bare name only",
			sql:  "CREATE TABLE t(a json, b Xml, c VARCHAR(10));",
			want: []columnFinding{disallowed("JSON", "json"), disallowed("XML", "Xml")},
		},
		{
			name: "nested data types are judged",
			sql:  "SELECT CAST(x AS XML) FROM t;",
			want: []columnFinding{disallowed("XML", "XML")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumn(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnTypeDisallowList(stmts, []string{"JSON", "xml", "VARCHAR(10)"})
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestColumnUnparsedForms pins branches V1 kept for tree shapes the parser
// does not build: a NOT NULL constraint node, and a DROP COLUMN action
// listing its columns in Names.
func TestColumnUnparsedForms(t *testing.T) {
	notNull := &ast.ColumnDef{Name: "a", Constraints: &ast.List{Items: []ast.Node{&ast.ConstraintDef{Type: ast.ConstraintNotNull}}}}
	create := &ast.CreateTableStmt{Name: &ast.TableRef{Object: "t"}, Columns: &ast.List{Items: []ast.Node{notNull}}}
	if got := ColumnNoNull([]ast.Node{create}); len(got) != 0 {
		t.Errorf("ColumnNoNull with a NOT NULL constraint: got %v", got)
	}

	drop := &ast.AlterTableStmt{Name: &ast.TableRef{Object: "t"}, Actions: &ast.List{Items: []ast.Node{&ast.AlterTableAction{
		Type:  ast.ATDropColumn,
		Names: &ast.List{Items: []ast.Node{&ast.TableRef{Object: "ID"}, &ast.TableRef{Object: "x"}}},
	}}}}
	got := ColumnRequired([]ast.Node{drop}, []string{"id"})
	if len(got) != 1 || got[0].Name != "id" || got[0].Table != "t" {
		t.Errorf("ColumnRequired with a Names drop list: got %v", got)
	}
}
