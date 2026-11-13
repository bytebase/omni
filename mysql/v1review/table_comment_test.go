package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableComment(t *testing.T) {
	missing := func(name, text string) finding { return finding{v1review.TableCommentMissing, name, text} }
	tooLong := func(name, text string) finding { return finding{v1review.TableCommentTooLong, name, text} }
	required := &v1review.TableCommentConfig{Required: true, MaxLength: 10}
	for _, tt := range []struct {
		name   string
		sql    string
		config *v1review.TableCommentConfig
		want   []finding
	}{
		{name: "comment within limit", sql: "CREATE TABLE t (a INT) COMMENT 'short';", config: required},
		{
			name:   "no comment",
			sql:    "CREATE TABLE t (a INT);",
			config: required,
			want:   []finding{missing("t", "CREATE TABLE t (a INT)")},
		},
		{name: "empty comment is a comment", sql: "CREATE TABLE t (a INT) COMMENT '';", config: required},
		{name: "lower-case option is a comment", sql: "CREATE TABLE t (a INT) comment = 'x';", config: required},
		{
			name:   "too long, counted in bytes",
			sql:    "CREATE TABLE t (a INT) COMMENT 'ééééé ';",
			config: required,
			want:   []finding{tooLong("t", "CREATE TABLE t (a INT) COMMENT 'ééééé '")},
		},
		{name: "not required", sql: "CREATE TABLE t (a INT);", config: &v1review.TableCommentConfig{MaxLength: 10}},
		{name: "negative length is no limit", sql: "CREATE TABLE t (a INT) COMMENT 'a long comment indeed';", config: &v1review.TableCommentConfig{MaxLength: -1}},
		{
			name:   "zero length allows only an empty comment",
			sql:    "CREATE TABLE t (a INT) COMMENT 'x';",
			config: &v1review.TableCommentConfig{},
			want:   []finding{tooLong("t", "CREATE TABLE t (a INT) COMMENT 'x'")},
		},
		{
			name:   "like and as-select judged, name without database",
			sql:    "CREATE TABLE db.t2 LIKE t;\nCREATE TABLE t3 AS SELECT 1;",
			config: required,
			want:   []finding{missing("t2", "CREATE TABLE db.t2 LIKE t"), missing("t3", "CREATE TABLE t3 AS SELECT 1")},
		},
		{name: "alter table not judged", sql: "ALTER TABLE t COMMENT 'a very long comment';", config: required},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT); END;", config: required},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []finding
			for _, f := range TableComment(list.Items, tt.config) {
				got = append(got, finding{f.Code, f.Name, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableCommentNilConfigPanicsOnlyWhenJudging(t *testing.T) {
	list, err := parser.Parse("ALTER TABLE t COMMENT 'x';")
	if err != nil {
		t.Fatal(err)
	}
	if got := TableComment(list.Items, nil); got != nil {
		t.Errorf("got %v, want none", got)
	}
	list, err = parser.Parse("CREATE TABLE t (a INT);")
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
