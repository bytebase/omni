package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnRequired(t *testing.T) {
	missing := func(table, name, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnMissing, Table: table, Name: name, Text: text}
	}
	required := []string{"id", "created_ts", "updated_ts"}
	for _, tt := range []struct {
		name     string
		sql      string
		required []string
		want     []ruleFinding
	}{
		{name: "all present", sql: "CREATE TABLE t (id int, created_ts int, updated_ts int, x int);", required: required},
		{
			name: "missing in name order, anchored at the statement", sql: "CREATE TABLE s.t (created_ts int);", required: required,
			want: []ruleFinding{missing("t", "id", "CREATE TABLE s.t (created_ts int)"), missing("t", "updated_ts", "CREATE TABLE s.t (created_ts int)")},
		},
		{
			name: "empty and duplicate required names", sql: "CREATE TABLE t (id int);", required: []string{"id", "id", ""},
			want: []ruleFinding{missing("t", "", "CREATE TABLE t (id int)")},
		},
		{
			name: "like lacks every column, create table as not judged", sql: "CREATE TABLE t AS SELECT 1 AS id; CREATE TABLE u (LIKE v);", required: []string{"id"},
			want: []ruleFinding{missing("u", "id", "CREATE TABLE u (LIKE v)")},
		},
		{
			name: "drop column anchored at its action, once per action", sql: "ALTER TABLE t DROP COLUMN id, DROP COLUMN x, DROP id;", required: required,
			want: []ruleFinding{missing("t", "id", "DROP COLUMN id"), missing("t", "id", "DROP id")},
		},
		{
			name: "rename away from a required name", sql: "ALTER TABLE t RENAME COLUMN id TO pk; ALTER TABLE t RENAME id TO id; ALTER TABLE t RENAME COLUMN x TO id;", required: required,
			want: []ruleFinding{missing("t", "id", "ALTER TABLE t RENAME COLUMN id TO pk")},
		},
		{
			name: "rename on a view", sql: "ALTER VIEW v RENAME COLUMN id TO pk;", required: required,
			want: []ruleFinding{missing("v", "id", "ALTER VIEW v RENAME COLUMN id TO pk")},
		},
		{name: "statements judged on their own", sql: "CREATE TABLE t (id int, created_ts int, updated_ts int); ALTER TABLE t ADD COLUMN x int;", required: required},
		{name: "quoted names compared as stored", sql: `CREATE TABLE t ("ID" int, created_ts int, updated_ts int);`, required: []string{"ID"}},
		{name: "no required columns", sql: "CREATE TABLE t (a int); ALTER TABLE t DROP COLUMN a;"},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC CREATE TABLE t (a int); END;", required: required},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return ColumnRequired(stmts, tt.required) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnRequired (-want +got):\n%s", diff)
			}
		})
	}
}
