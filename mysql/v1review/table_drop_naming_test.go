package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableDropNaming(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "_delete$"}
	mismatch := func(name, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []located
	}{
		{name: "drop ok", sql: "DROP TABLE foo_delete;", config: pattern},
		{
			name:   "drop mismatch",
			sql:    "DROP TABLE IF EXISTS foo;",
			config: pattern,
			want:   []located{mismatch("foo", "DROP TABLE IF EXISTS foo")},
		},
		{
			name:   "every table judged, anchored to the statement",
			sql:    "DROP TABLE foo,\n  bar_delete,\n  baz;",
			config: pattern,
			want: []located{
				mismatch("foo", "DROP TABLE foo,\n  bar_delete,\n  baz"),
				mismatch("baz", "DROP TABLE foo,\n  bar_delete,\n  baz"),
			},
		},
		{name: "database not judged", sql: "DROP TABLE Shop.foo_delete;", config: pattern},
		{name: "database not part of the judged name", sql: "DROP TABLE Shop.tmp_foo;", config: v1review.NamingConfig{Pattern: "^tmp_"}},
		{
			name:   "name judged without its database",
			sql:    "DROP TABLE shop_delete.foo;",
			config: pattern,
			want:   []located{mismatch("foo", "DROP TABLE shop_delete.foo")},
		},
		{
			name:   "quoted name judged without its quotes",
			sql:    "DROP TABLE `foo_delete$`, `bar_delete`;",
			config: pattern,
			want:   []located{mismatch("foo_delete$", "DROP TABLE `foo_delete$`, `bar_delete`")},
		},
		{
			name:   "temporary table judged",
			sql:    "DROP TEMPORARY TABLE tmp;",
			config: pattern,
			want:   []located{mismatch("tmp", "DROP TEMPORARY TABLE tmp")},
		},
		{name: "drop view not judged", sql: "DROP VIEW foo;", config: pattern},
		{name: "truncate not judged", sql: "TRUNCATE TABLE foo;", config: pattern},
		{name: "alter drop column not judged", sql: "ALTER TABLE foo DROP COLUMN bar;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN DROP TABLE foo; END;", config: pattern},
		{name: "length not consulted", sql: "DROP TABLE a_very_long_name_delete;", config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 4}},
		{name: "length alone requires nothing", sql: "DROP TABLE a_very_long_name;", config: v1review.NamingConfig{MaxLength: 4}},
		{name: "empty config requires nothing", sql: "DROP TABLE foo;"},
		{
			name:   "findings in statement order",
			sql:    "DROP TABLE one;\nDROP TABLE two;",
			config: pattern,
			want:   []located{mismatch("one", "DROP TABLE one"), mismatch("two", "DROP TABLE two")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TableDropNaming(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("TableDropNaming: %v", err)
			}
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableDropNamingInvalidPattern(t *testing.T) {
	if _, err := TableDropNaming(parse(t, "DROP TABLE t;"), v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
