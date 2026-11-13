package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableDropNaming(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "_delete$"}
	mismatch := func(name, text string) anchored {
		return anchored{Code: v1review.NameMismatch, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []anchored
	}{
		{name: "drop ok", sql: "DROP TABLE foo_delete;", config: pattern},
		{
			name:   "drop anchors the statement",
			sql:    "DROP TABLE IF EXISTS\n  foo;",
			config: pattern,
			want:   []anchored{mismatch("foo", "DROP TABLE IF EXISTS\n  foo")},
		},
		{
			name:   "every table judged, in written order",
			sql:    "DROP TABLE foo, bar_delete, baz CASCADE;",
			config: pattern,
			want: []anchored{
				mismatch("foo", "DROP TABLE foo, bar_delete, baz CASCADE"),
				mismatch("baz", "DROP TABLE foo, bar_delete, baz CASCADE"),
			},
		},
		{name: "qualifier not judged", sql: "DROP TABLE Shop.del_foo;", config: v1review.NamingConfig{Pattern: "^del_"}},
		{
			name:   "qualified name judged without its qualifier",
			sql:    "DROP TABLE del_shop.foo;",
			config: v1review.NamingConfig{Pattern: "^del_"},
			want:   []anchored{mismatch("foo", "DROP TABLE del_shop.foo")},
		},
		{
			name:   "temporary table judged",
			sql:    "DROP TEMPORARY TABLE foo;",
			config: pattern,
			want:   []anchored{mismatch("foo", "DROP TEMPORARY TABLE foo")},
		},
		{
			name:   "name judged as parsed, quotes removed",
			sql:    "DROP TABLE `foo``_delete`, `bar_delete`, `baz _delete$`;",
			config: pattern,
			want:   []anchored{mismatch("baz _delete$", "DROP TABLE `foo``_delete`, `bar_delete`, `baz _delete$`")},
		},
		{
			name:   "unanchored pattern matches anywhere in the name",
			sql:    "DROP TABLE foo_delete_old, foo;",
			config: v1review.NamingConfig{Pattern: "_delete"},
			want:   []anchored{mismatch("foo", "DROP TABLE foo_delete_old, foo")},
		},
		{
			name:   "length not consulted",
			sql:    "DROP TABLE foo_delete, foo;",
			config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 3},
			want:   []anchored{mismatch("foo", "DROP TABLE foo_delete, foo")},
		},
		{name: "length alone requires nothing", sql: "DROP TABLE foo;", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE foo;"},
		{name: "drop view not judged", sql: "DROP VIEW foo;", config: pattern},
		{name: "drop sequence not judged", sql: "DROP SEQUENCE foo;", config: pattern},
		{name: "drop database not judged", sql: "DROP DATABASE foo;", config: pattern},
		{name: "alter drop column and truncate not judged", sql: "ALTER TABLE foo DROP COLUMN a; TRUNCATE TABLE foo;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN DROP TABLE foo; END;", config: pattern},
		{
			name:   "findings in statement order",
			sql:    "DROP TABLE foo;\nDROP TABLE bar;",
			config: pattern,
			want: []anchored{
				mismatch("foo", "DROP TABLE foo"),
				mismatch("bar", "DROP TABLE bar"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TableDropNaming(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("TableDropNaming: %v", err)
			}
			if diff := cmp.Diff(tt.want, anchor(tt.sql, got)); diff != "" {
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
