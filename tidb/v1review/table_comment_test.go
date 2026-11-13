package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableComment(t *testing.T) {
	missing := func(name, text string) located {
		return located{Code: v1review.TableCommentMissing, Name: name, Text: text}
	}
	tooLong := func(name, text string) located {
		return located{Code: v1review.TableCommentTooLong, Name: name, Text: text}
	}
	required := &v1review.TableCommentConfig{Required: true, MaxLength: 10}
	for _, tt := range []struct {
		name   string
		sql    string
		config *v1review.TableCommentConfig
		want   []located
	}{
		{name: "comment within limit", sql: "CREATE TABLE t (a INT) COMMENT 'short';", config: required},
		{
			name:   "no comment",
			sql:    "CREATE TABLE t (a INT);",
			config: required,
			want:   []located{missing("t", "CREATE TABLE t (a INT)")},
		},
		{name: "empty comment is a comment", sql: "CREATE TABLE t (a INT) COMMENT '';", config: required},
		{name: "lower-case empty comment is a comment", sql: "CREATE TABLE t (a INT) comment '';", config: required},
		{
			name:   "too long, counted in bytes",
			sql:    "CREATE TABLE t (a INT) COMMENT 'ééééé ';",
			config: required,
			want:   []located{tooLong("t", "CREATE TABLE t (a INT) COMMENT 'ééééé '")},
		},
		{
			name:   "first comment option is the one judged",
			sql:    "CREATE TABLE t (a INT) COMMENT 'a much longer one' COMMENT 'ab';",
			config: required,
			want:   []located{tooLong("t", "CREATE TABLE t (a INT) COMMENT 'a much longer one' COMMENT 'ab'")},
		},
		{name: "not required", sql: "CREATE TABLE t (a INT);", config: &v1review.TableCommentConfig{MaxLength: 10}},
		{name: "negative length is no limit", sql: "CREATE TABLE t (a INT) COMMENT 'a long comment indeed';", config: &v1review.TableCommentConfig{MaxLength: -1}},
		{
			name:   "zero length allows only an empty comment",
			sql:    "CREATE TABLE t (a INT) COMMENT 'x';",
			config: &v1review.TableCommentConfig{},
			want:   []located{tooLong("t", "CREATE TABLE t (a INT) COMMENT 'x'")},
		},
		{
			name:   "like and as-select judged, name without database",
			sql:    "CREATE TABLE db.t2 LIKE t;\nCREATE TABLE t3 AS SELECT 1;",
			config: required,
			want:   []located{missing("t2", "CREATE TABLE db.t2 LIKE t"), missing("t3", "CREATE TABLE t3 AS SELECT 1")},
		},
		{name: "column comment is not a table comment", sql: "CREATE TABLE t (a INT COMMENT 'x');", config: &v1review.TableCommentConfig{}},
		{name: "alter table not judged", sql: "ALTER TABLE t COMMENT 'a very long comment';", config: required},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := locate(tt.sql, TableComment(parse(t, tt.sql), tt.config))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableCommentNilConfigPanicsOnlyWhenJudging(t *testing.T) {
	if got := TableComment(parse(t, "ALTER TABLE t COMMENT 'x';"), nil); got != nil {
		t.Errorf("got %v, want none", got)
	}
	stmts := parse(t, "CREATE TABLE t (a INT);")
	defer func() {
		if recover() == nil {
			t.Error("want a panic")
		}
	}()
	TableComment(stmts, nil)
}
