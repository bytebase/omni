package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableComment(t *testing.T) {
	type commentFinding struct {
		Code      v1review.Code
		Name      string
		Statement int
		Text      string
	}
	missing := func(name string, statement int, text string) commentFinding {
		return commentFinding{v1review.TableCommentMissing, name, statement, text}
	}
	tooLong := func(name string, statement int, text string) commentFinding {
		return commentFinding{v1review.TableCommentTooLong, name, statement, text}
	}
	required := &v1review.TableCommentConfig{Required: true, MaxLength: 5}
	for _, tt := range []struct {
		name   string
		sql    string
		config *v1review.TableCommentConfig
		want   []commentFinding
	}{
		{name: "commented", sql: "CREATE TABLE t (a int);\nCOMMENT ON TABLE t IS 'ok';", config: required},
		{
			name:   "no comment",
			sql:    "CREATE TABLE t (a int);",
			config: required,
			want:   []commentFinding{missing("t", 0, "CREATE TABLE t (a int)")},
		},
		{name: "comment before the create counts", sql: "COMMENT ON TABLE t IS 'ok';\nCREATE TABLE t (a int);", config: required},
		{
			name:   "empty and null comments are missing",
			sql:    "CREATE TABLE t (a int);\nCOMMENT ON TABLE t IS '';\nCREATE TABLE u (a int);\nCOMMENT ON TABLE u IS NULL;",
			config: required,
			want:   []commentFinding{missing("t", 0, "CREATE TABLE t (a int)"), missing("u", 2, "CREATE TABLE u (a int)")},
		},
		{
			name:   "last comment wins",
			sql:    "CREATE TABLE t (a int);\nCOMMENT ON TABLE t IS 'ok';\nCOMMENT ON TABLE t IS '';",
			config: required,
			want:   []commentFinding{missing("t", 0, "CREATE TABLE t (a int)")},
		},
		{
			name:   "too long in bytes anchors the last comment",
			sql:    "CREATE TABLE t (a int);\nCOMMENT ON TABLE t IS 'x';\nCOMMENT ON TABLE t IS 'ééé';",
			config: required,
			want:   []commentFinding{tooLong("t", 2, "COMMENT ON TABLE t IS 'ééé'")},
		},
		{name: "zero length is no limit", sql: "CREATE TABLE t (a int);\nCOMMENT ON TABLE t IS 'a long comment';", config: &v1review.TableCommentConfig{Required: true}},
		{name: "negative length is no limit", sql: "CREATE TABLE t (a int);\nCOMMENT ON TABLE t IS 'a long comment';", config: &v1review.TableCommentConfig{MaxLength: -1}},
		{name: "not required", sql: "CREATE TABLE t (a int);", config: &v1review.TableCommentConfig{MaxLength: 5}},
		{
			name:   "schema must match: public is the default, others are printed",
			sql:    "CREATE TABLE public.t (a int);\nCOMMENT ON TABLE t IS 'ok';\nCREATE TABLE s.u (a int);\nCOMMENT ON TABLE u IS 'ok';",
			config: required,
			want:   []commentFinding{missing("s.u", 2, "CREATE TABLE s.u (a int)")},
		},
		{
			name:   "names compare as parsed",
			sql:    "CREATE TABLE \"T\" (a int);\nCOMMENT ON TABLE T IS 'ok';",
			config: required,
			want:   []commentFinding{missing("T", 0, "CREATE TABLE \"T\" (a int)")},
		},
		{
			name:   "three-part comment is ignored",
			sql:    "CREATE TABLE s.t (a int);\nCOMMENT ON TABLE db.s.t IS 'ok';",
			config: required,
			want:   []commentFinding{missing("s.t", 0, "CREATE TABLE s.t (a int)")},
		},
		{
			name:   "a re-created table anchors its last create, ordered by first create",
			sql:    "CREATE TABLE b (a int);\nCREATE TABLE a (a int);\nCREATE TABLE b (x int);",
			config: required,
			want:   []commentFinding{missing("b", 2, "CREATE TABLE b (x int)"), missing("a", 1, "CREATE TABLE a (a int)")},
		},
		{name: "column comment not a table comment", sql: "CREATE TABLE t (a int);\nCOMMENT ON COLUMN t.a IS 'ok';", config: &v1review.TableCommentConfig{MaxLength: 1}},
		{name: "create table as not judged", sql: "CREATE TABLE t AS SELECT 1 AS a;", config: required},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void AS $$ CREATE TABLE t (a int) $$ LANGUAGE sql;", config: required},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []commentFinding
			for _, f := range TableComment(list.Items, tt.config) {
				got = append(got, commentFinding{f.Code, f.Name, f.Statement, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableCommentNilConfigPanicsOnlyWhenJudging(t *testing.T) {
	list, err := parser.Parse("COMMENT ON TABLE t IS 'x';\nALTER TABLE t ADD COLUMN b int;")
	if err != nil {
		t.Fatal(err)
	}
	if got := TableComment(list.Items, nil); got != nil {
		t.Errorf("got %v, want none", got)
	}
	list, err = parser.Parse("CREATE TABLE t (a int);")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recover() == nil {
			t.Error("want a panic")
		}
	}()
	TableComment(list.Items, nil)
}
