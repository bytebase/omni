package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

// tableFindings parses sql, runs rule over its statements, and returns the
// findings as the tests compare them.
func tableFindings(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []finding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []finding
	for _, f := range rule(list.Items) {
		got = append(got, finding{f.Code, f.Name, sql[f.Start:f.End]})
	}
	return got
}

func TestTableDisallowDDL(t *testing.T) {
	disallow := []string{"t", "u", "t"}
	rule := func(stmts []ast.Node) []v1review.Finding { return TableDisallowDDL(stmts, disallow) }
	f := func(name, text string) finding { return finding{v1review.TableDisallowed, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "create", sql: "CREATE TABLE t (a INT);", want: []finding{f("t", "CREATE TABLE t (a INT)")}},
		{name: "alter, database dropped", sql: "ALTER TABLE db.t ADD COLUMN b INT;", want: []finding{f("t", "ALTER TABLE db.t ADD COLUMN b INT")}},
		{name: "drop judges each table", sql: "DROP TABLE t, x, u;", want: []finding{f("t", "DROP TABLE t, x, u"), f("u", "DROP TABLE t, x, u")}},
		{name: "rename judges the old name", sql: "RENAME TABLE t TO x, x TO u;", want: []finding{f("t", "RENAME TABLE t TO x, x TO u")}},
		{name: "truncate", sql: "TRUNCATE TABLE u;", want: []finding{f("u", "TRUNCATE TABLE u")}},
		{name: "case-sensitive", sql: "CREATE TABLE T (a INT);"},
		{name: "other tables", sql: "CREATE TABLE x (a INT);"},
		{name: "dml not judged", sql: "INSERT INTO t VALUES (1);"},
		{name: "index statements not judged", sql: "CREATE INDEX i ON t (a);"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN DROP TABLE t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, rule)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
	list, err := parser.Parse("DROP TABLE t;")
	if err != nil {
		t.Fatal(err)
	}
	if got := TableDisallowDDL(list.Items, nil); got != nil {
		t.Errorf("nil list: got %v, want none", got)
	}
}
