package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/bytebase/omni/v1review"
)

// anchored is a finding as the keyword and case tests compare it: its
// code, name, table and pattern, and the text its range anchors, empty for
// a finding without a range.
type anchored struct {
	Code    v1review.Code
	Name    string
	Table   string
	Pattern string
	Text    string
}

func anchor(sql string, findings []v1review.Finding) []anchored {
	var got []anchored
	for _, f := range findings {
		got = append(got, anchored{f.Code, f.Name, f.Table, f.Pattern, sql[f.Start:f.End]})
	}
	return got
}

func parse(t *testing.T, sql string) []ast.Node {
	t.Helper()
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return file.Stmts
}

func TestNamingIdentifierNoKeyword(t *testing.T) {
	keyword := func(name, text string) anchored {
		return anchored{Code: v1review.NameIsKeyword, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []anchored
	}{
		{name: "no keyword", sql: "CREATE TABLE t (name VARCHAR, age INT);"},
		{
			name: "a quoted keyword column, anchored at the first column",
			sql:  "CREATE TABLE t(\n  name VARCHAR,\n  \"CONNECT\" VARCHAR\n);",
			want: []anchored{keyword(`"CONNECT"`, "name VARCHAR")},
		},
		{
			name: "a quoted name matches whatever its case",
			sql:  `CREATE TABLE t (id INT, "connect" INT, "Select" INT);`,
			want: []anchored{keyword(`"connect"`, "id INT"), keyword(`"Select"`, "id INT")},
		},
		{
			name: "an unquoted name is reported as written",
			sql:  "CREATE TABLE t (id INT, issue INT, Trigger INT, GSCLUSTER INT);",
			want: []anchored{keyword("issue", "id INT"), keyword("Trigger", "id INT"), keyword("GSCLUSTER", "id INT")},
		},
		{
			name: "near misses",
			sql:  `CREATE TABLE t ("CONNECTS" INT, "SELEC" INT, "SELECT " INT, "TRY CAST" INT, "" INT, index INT, "LIMIT" INT);`,
		},
		{
			name: "create table as select",
			sql:  "CREATE TABLE c (\n  name INT,\n  \"connect\" INT\n) AS SELECT a, b FROM t;",
			want: []anchored{keyword(`"connect"`, "name INT")},
		},
		{
			name: "the anchor is a constraint that comes first",
			sql:  "CREATE TABLE t (\n  PRIMARY KEY (id),\n  id INT,\n  \"from\" INT\n);",
			want: []anchored{keyword(`"from"`, "PRIMARY KEY (id)")},
		},
		{
			name: "the anchor is an inline index that comes first",
			sql:  "CREATE HYBRID TABLE t (\n  INDEX idx (id),\n  id INT PRIMARY KEY,\n  \"from\" INT\n);",
			want: []anchored{keyword(`"from"`, "INDEX idx (id)")},
		},
		{
			name: "rename column judges only the new name, anchored at it",
			sql:  "ALTER TABLE t RENAME COLUMN \"select\" TO b;\nALTER TABLE t RENAME COLUMN a\n  TO \"connect\";",
			want: []anchored{keyword(`"connect"`, `"connect"`)},
		},
		{
			name: "each rename column of one statement is judged",
			sql:  "ALTER TABLE t RENAME COLUMN a TO \"select\",\n  RENAME TO \"WHERE\",\n  RENAME COLUMN b TO c,\n  RENAME COLUMN d TO \"of\";",
			want: []anchored{keyword(`"select"`, `"select"`), keyword(`"of"`, `"of"`)},
		},
		{
			name: "table, constraint and index names are not judged",
			sql:  "CREATE TABLE \"SELECT\" (id INT, CONSTRAINT \"from\" PRIMARY KEY (id));\nCREATE HYBRID TABLE t (id INT PRIMARY KEY, INDEX \"where\" (id));",
		},
		{
			name: "other alter table actions are not judged",
			sql:  "ALTER TABLE t ADD COLUMN \"select\" INT;\nALTER TABLE t DROP COLUMN \"select\";\nALTER TABLE t RENAME TO \"SELECT\";\nALTER TABLE t RENAME CONSTRAINT a TO \"select\";",
		},
		{
			name: "other statements are not judged",
			sql:  "CREATE VIEW v (\"select\") AS SELECT 1;\nSELECT 1 AS \"select\";\nALTER VIEW v RENAME TO \"SELECT\";",
		},
		{name: "a statement under EXPLAIN is not judged", sql: `EXPLAIN CREATE TABLE t ("select" INT);`},
		{name: "a scripting block is not judged", sql: "BEGIN\n  CREATE TABLE t (\"select\" INT);\nEND"},
		{
			name: "a task body is not judged",
			sql:  `CREATE TASK tk WAREHOUSE = wh SCHEDULE = '1 MINUTE' AS CREATE TABLE t ("select" INT);`,
		},
		{
			name: "findings in statement order",
			sql:  "CREATE TABLE a (\"to\" INT);\nALTER TABLE a RENAME COLUMN x TO \"of\";\nCREATE TABLE b (\"on\" INT, \"or\" INT);",
			want: []anchored{keyword(`"to"`, `"to" INT`), keyword(`"of"`, `"of"`), keyword(`"on"`, `"on" INT`), keyword(`"or"`, `"on" INT`)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := anchor(tt.sql, NamingIdentifierNoKeyword(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
