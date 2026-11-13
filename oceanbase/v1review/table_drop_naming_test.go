package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableDropNaming(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "_delete$"}
	mismatch := func(name, text string) judged {
		return judged{Code: v1review.NameMismatch, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []judged
	}{
		{name: "drop ok", sql: "DROP TABLE foo_delete;", config: pattern},
		{
			name:   "drop mismatch",
			sql:    "DROP TABLE IF EXISTS foo;",
			config: pattern,
			want:   []judged{mismatch("foo", "DROP TABLE IF EXISTS foo")},
		},
		{
			name:   "each table in order, anchoring the statement",
			sql:    "DROP TABLE foo,\n  bar_delete,\n  baz;",
			config: pattern,
			want: []judged{
				mismatch("foo", "DROP TABLE foo,\n  bar_delete,\n  baz"),
				mismatch("baz", "DROP TABLE foo,\n  bar_delete,\n  baz"),
			},
		},
		{name: "qualifier not judged", sql: "DROP TABLE Shop.foo_delete;", config: v1review.NamingConfig{Pattern: "^[a-z_]+$"}},
		{
			name:   "qualified name judged without its database",
			sql:    "DROP TABLE shop_delete.foo;",
			config: pattern,
			want:   []judged{mismatch("foo", "DROP TABLE shop_delete.foo")},
		},
		{
			name:   "quoted name judged without its quotes",
			sql:    "DROP TABLE `foo _delete`, `_delete``x`;",
			config: pattern,
			want:   []judged{mismatch("_delete`x", "DROP TABLE `foo _delete`, `_delete``x`")},
		},
		{
			name:   "temporary table judged",
			sql:    "DROP TEMPORARY TABLE foo;",
			config: pattern,
			want:   []judged{mismatch("foo", "DROP TEMPORARY TABLE foo")},
		},
		{name: "drop view not judged", sql: "DROP VIEW foo;", config: pattern},
		{name: "truncate not judged", sql: "TRUNCATE TABLE foo;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN DROP TABLE foo; END;", config: pattern},
		{name: "length not judged", sql: "DROP TABLE foo_delete;", config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 3}},
		{name: "length alone requires nothing", sql: "DROP TABLE foo;", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE foo;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(t, tt.sql, func(stmts []ast.Node) ([]v1review.Finding, error) { return TableDropNaming(stmts, tt.config) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableDropNamingInvalidPattern(t *testing.T) {
	list, err := parser.Parse("DROP TABLE t;")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := TableDropNaming(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
