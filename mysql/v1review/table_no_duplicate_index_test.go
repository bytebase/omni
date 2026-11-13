package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

// duplicateFinding is a TableNoDuplicateIndex finding as the tests compare
// it: the duplicate's name, the reporting table, the anchored text, and
// the anchored statement's index.
type duplicateFinding struct {
	Name, Table, Text string
	Stmt              int
}

func TestTableNoDuplicateIndex(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  []string
		want []duplicateFinding
	}{
		{
			name: "duplicate in a create table anchors its constraint",
			sql:  []string{"CREATE TABLE t (a INT, KEY i1 (a), KEY i2 (a))"},
			want: []duplicateFinding{{"i2", "t", "KEY i2 (a)", 0}},
		},
		{name: "different columns", sql: []string{"CREATE TABLE t (a INT, b INT, KEY i1 (a), KEY i2 (a, b))"}},
		{name: "different kinds", sql: []string{"CREATE TABLE t (a INT, KEY i1 (a), UNIQUE KEY u1 (a), FULLTEXT KEY f1 (a))"}},
		{name: "different index types", sql: []string{"CREATE TABLE t (a INT, KEY i1 USING HASH (a), KEY i2 (a))"}},
		{
			name: "index type after the columns not read",
			sql:  []string{"CREATE TABLE t (a INT, KEY i1 (a) USING HASH, KEY i2 (a))"},
			want: []duplicateFinding{{"i2", "t", "KEY i2 (a)", 0}},
		},
		{
			name: "primary key and plain index are duplicates",
			sql:  []string{"CREATE TABLE t (a INT, PRIMARY KEY (a), KEY i1 (a))"},
			want: []duplicateFinding{{"i1", "t", "KEY i1 (a)", 0}},
		},
		{name: "column-level keys not judged", sql: []string{"CREATE TABLE t (a INT PRIMARY KEY, b INT UNIQUE, KEY i1 (a), KEY i2 (b))"}},
		{name: "foreign keys not judged", sql: []string{"CREATE TABLE t (a INT, FOREIGN KEY f1 (a) REFERENCES u (a), FOREIGN KEY f2 (a) REFERENCES u (a))"}},
		{name: "table names compared case-sensitively", sql: []string{"CREATE TABLE t (a INT, KEY i1 (a))", "CREATE TABLE T (a INT, KEY i2 (a))"}},
		{name: "same columns on different tables", sql: []string{"CREATE TABLE t (a INT, KEY i1 (a))", "CREATE TABLE u (a INT, KEY i1 (a))"}},
		{
			name: "alter table anchors its statement",
			sql:  []string{"CREATE TABLE t (a INT, KEY i1 (a))", "ALTER TABLE t ADD INDEX i2 (a)"},
			want: []duplicateFinding{{"i2", "t", "ALTER TABLE t ADD INDEX i2 (a)", 1}},
		},
		{
			name: "create index anchors its statement, expression parts left out",
			sql:  []string{"CREATE TABLE t (a INT, KEY i1 (a))", "CREATE INDEX i2 ON db.t (a, (a + 1))"},
			want: []duplicateFinding{{"i2", "t", "CREATE INDEX i2 ON db.t (a, (a + 1))", 1}},
		},
		{
			name: "every later judged statement repeats the first duplicate, naming its own table",
			sql: []string{
				"CREATE TABLE t (a INT, KEY i1 (a), KEY i2 (a), KEY i3 (a))",
				"INSERT INTO t VALUES (1)",
				"DROP INDEX i2 ON t",
				"CREATE TABLE u (b INT)",
				"ALTER TABLE u ADD COLUMN c INT",
				"CREATE INDEX j ON u (b)",
			},
			want: []duplicateFinding{
				{"i2", "t", "KEY i2 (a)", 0},
				{"i2", "u", "KEY i2 (a)", 0},
				{"i2", "u", "KEY i2 (a)", 0},
				{"i2", "u", "KEY i2 (a)", 0},
			},
		},
		{name: "routine body not judged", sql: []string{"CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT, KEY i1 (a), KEY i2 (a)); END"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Parse each statement on its own, as Bytebase does, so ranges
			// index each statement's own text.
			var stmts []ast.Node
			for _, sql := range tt.sql {
				list, err := parser.Parse(sql)
				if err != nil {
					t.Fatalf("Parse %q: %v", sql, err)
				}
				stmts = append(stmts, list.Items...)
			}
			var got []duplicateFinding
			for _, f := range TableNoDuplicateIndex(stmts) {
				if f.Code != v1review.TableDuplicateIndex {
					t.Errorf("code %q", f.Code)
				}
				got = append(got, duplicateFinding{f.Name, f.Table, tt.sql[f.Statement][f.Start:f.End], f.Statement})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
