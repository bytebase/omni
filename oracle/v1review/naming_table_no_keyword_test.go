package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestNamingTableNoKeyword(t *testing.T) {
	keyword := func(name, text string) finding { return finding{v1review.NameIsKeyword, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{
			name: "quoted keyword",
			sql:  `CREATE TABLE "table"(id int, name varchar(255));`,
			want: []finding{keyword("table", `CREATE TABLE "table"(id int, name varchar(255))`)},
		},
		{
			name: "unquoted keyword folds to upper case",
			sql:  "CREATE TABLE avg(id int);",
			want: []finding{keyword("AVG", "CREATE TABLE avg(id int)")},
		},
		{name: "not a keyword", sql: "CREATE TABLE analyze(id int);"},
		{name: "column keyword not judged", sql: "CREATE TABLE t (name NUMBER, \"DATE\" DATE);"},
		{name: "qualifier not judged", sql: `CREATE TABLE "TABLE".order_item (id NUMBER);`},
		{
			name: "qualified name judged",
			sql:  `CREATE TABLE shop."Size" (id NUMBER);`,
			want: []finding{keyword("Size", `CREATE TABLE shop."Size" (id NUMBER)`)},
		},
		{
			name: "create as select",
			sql:  `CREATE TABLE "DATE" AS SELECT 1 AS id FROM dual;`,
			want: []finding{keyword("DATE", `CREATE TABLE "DATE" AS SELECT 1 AS id FROM dual`)},
		},
		{
			name: "global temporary table judged",
			sql:  `CREATE GLOBAL TEMPORARY TABLE "LEVEL" (id NUMBER);`,
			want: []finding{keyword("LEVEL", `CREATE GLOBAL TEMPORARY TABLE "LEVEL" (id NUMBER)`)},
		},
		{
			name: "alter rename to",
			sql:  "ALTER TABLE tech_book RENAME TO avg;",
			want: []finding{keyword("AVG", "RENAME TO avg")},
		},
		{name: "alter rename to a non-keyword", sql: "ALTER TABLE tech_book RENAME TO analyze;"},
		{name: "alter rename column not judged", sql: `ALTER TABLE t RENAME COLUMN a TO "LEVEL";`},
		{name: "rename statement not judged", sql: `RENAME t TO "LEVEL";`},
		{name: "view not judged", sql: `CREATE VIEW "LEVEL" AS SELECT 1 AS id FROM dual;`},
		{name: "explain plan not judged", sql: `EXPLAIN PLAN FOR CREATE TABLE "LEVEL" AS SELECT 1 AS id FROM dual`},
		{
			name: "create schema nested table not judged",
			sql:  `CREATE SCHEMA AUTHORIZATION shop CREATE TABLE "LEVEL" (id NUMBER)`,
		},
		{
			name: "findings in statement order",
			sql:  "CREATE TABLE t (id NUMBER);\nCREATE TABLE\n  \"SIZE\" (id NUMBER);\nALTER TABLE t\n  RENAME TO \"type\";",
			want: []finding{
				keyword("SIZE", "CREATE TABLE\n  \"SIZE\" (id NUMBER)"),
				keyword("type", `RENAME TO "type"`),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := anchored(t, tt.sql, NamingTableNoKeyword(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingTableNoKeywordBareStatement(t *testing.T) {
	sql := `CREATE TABLE "table" (id NUMBER)`
	stmts := parse(t, sql)
	raw, ok := stmts[0].(*ast.RawStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.RawStmt", stmts[0])
	}
	want := []finding{{v1review.NameIsKeyword, "table", sql}}
	if diff := cmp.Diff(want, anchored(t, sql, NamingTableNoKeyword([]ast.Node{raw.Stmt}))); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
