package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnComment(t *testing.T) {
	missing := func(table, name string, stmt int, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnCommentMissing, Table: table, Name: name, Statement: stmt, Text: text}
	}
	tooLong := func(table, name string, stmt int, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnCommentTooLong, Table: table, Name: name, Statement: stmt, Text: text}
	}
	for _, tt := range []struct {
		name      string
		sql       string
		required  bool
		maxLength int
		want      []ruleFinding
	}{
		{name: "commented column ok", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN t.a IS 'note';", required: true, maxLength: 10},
		{
			name: "missing comment anchors the defining statement", sql: "CREATE TABLE t (a int, b int); COMMENT ON COLUMN t.a IS 'note';", required: true, maxLength: -1,
			want: []ruleFinding{missing("t", "b", 0, "CREATE TABLE t (a int, b int)")},
		},
		{name: "missing comment not required", sql: "CREATE TABLE t (a int);", maxLength: 1},
		{
			name: "too long anchors the comment statement", sql: "CREATE TABLE t (a int);\nCOMMENT ON COLUMN t.a IS 'long comment';", maxLength: 4,
			want: []ruleFinding{tooLong("t", "a", 1, "COMMENT ON COLUMN t.a IS 'long comment'")},
		},
		{name: "length within limit", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN t.a IS 'four';", maxLength: 4},
		{name: "length in bytes", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN t.a IS 'éé';", maxLength: 3, want: []ruleFinding{tooLong("t", "a", 1, "COMMENT ON COLUMN t.a IS 'éé'")}},
		{name: "negative max length sets no limit", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN t.a IS 'long comment';", required: true, maxLength: -1},
		{
			name: "zero max length reports every non-empty comment", sql: "CREATE TABLE t (a int, b int); COMMENT ON COLUMN t.a IS 'x';", required: true, maxLength: 0,
			want: []ruleFinding{tooLong("t", "a", 1, "COMMENT ON COLUMN t.a IS 'x'"), missing("t", "b", 0, "CREATE TABLE t (a int, b int)")},
		},
		{
			name: "last comment wins and IS NULL empties it", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN t.a IS 'note'; COMMENT ON COLUMN t.a IS NULL;", required: true, maxLength: 10,
			want: []ruleFinding{missing("t", "a", 0, "CREATE TABLE t (a int)")},
		},
		{
			name: "last comment wins over an earlier long one", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN t.a IS 'long comment'; COMMENT ON COLUMN t.a IS 'ok';", maxLength: 4,
		},
		{name: "comment before the column counts", sql: "COMMENT ON COLUMN t.a IS 'note'; CREATE TABLE t (a int);", required: true, maxLength: -1},
		{name: "schema ignored on both sides", sql: "CREATE TABLE s1.t (a int); COMMENT ON COLUMN s2.t.a IS 'note';", required: true, maxLength: -1},
		{
			name: "four-part comment name is not read", sql: "CREATE TABLE t (a int); COMMENT ON COLUMN d.s.t.a IS 'note';", required: true, maxLength: -1,
			want: []ruleFinding{missing("t", "a", 0, "CREATE TABLE t (a int)")},
		},
		{
			name: "add column judged on any relation", sql: "ALTER TABLE s.t ADD COLUMN a int, ADD b int, DROP COLUMN c; ALTER FOREIGN TABLE f ADD COLUMN x int;", required: true, maxLength: -1,
			want: []ruleFinding{
				missing("t", "a", 0, "ALTER TABLE s.t ADD COLUMN a int, ADD b int, DROP COLUMN c"),
				missing("t", "b", 0, "ALTER TABLE s.t ADD COLUMN a int, ADD b int, DROP COLUMN c"),
				missing("f", "x", 1, "ALTER FOREIGN TABLE f ADD COLUMN x int"),
			},
		},
		{
			name: "column defined twice judged twice", sql: "CREATE TABLE t (a int); CREATE TABLE t (a int);", required: true, maxLength: -1,
			want: []ruleFinding{missing("t", "a", 0, "CREATE TABLE t (a int)"), missing("t", "a", 1, "CREATE TABLE t (a int)")},
		},
		{name: "names compared as stored", sql: `CREATE TABLE "T" ("A" int, b int); COMMENT ON COLUMN "T"."A" IS 'x'; COMMENT ON COLUMN "T".B IS 'y';`, required: true, maxLength: -1},
		{name: "case-folded name does not match a quoted one", sql: `CREATE TABLE t ("A" int); COMMENT ON COLUMN t.a IS 'x';`, required: true, maxLength: -1, want: []ruleFinding{missing("t", "A", 0, `CREATE TABLE t ("A" int)`)}},
		{name: "routine body and other statements not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC CREATE TABLE t (a int); END; SELECT 1; COMMENT ON TABLE t IS 'x';", required: true, maxLength: -1},
		{name: "create table as has no columns", sql: "CREATE TABLE t AS SELECT 1 AS a;", required: true, maxLength: -1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return ColumnComment(stmts, tt.required, tt.maxLength) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnComment (-want +got):\n%s", diff)
			}
		})
	}
}
