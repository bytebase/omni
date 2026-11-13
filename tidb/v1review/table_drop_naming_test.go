package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
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
		{name: "drop ok", sql: "DROP TABLE IF EXISTS foo_delete;", config: pattern},
		{
			name:   "every table judged, anchoring the statement",
			sql:    "DROP TABLE IF EXISTS foo,\n  bar_delete,\n  baz;",
			config: pattern,
			want: []located{
				mismatch("foo", "DROP TABLE IF EXISTS foo,\n  bar_delete,\n  baz"),
				mismatch("baz", "DROP TABLE IF EXISTS foo,\n  bar_delete,\n  baz"),
			},
		},
		{
			name:   "qualifier not judged",
			sql:    "DROP TABLE shop_delete.foo, Shop.bar_delete;",
			config: v1review.NamingConfig{Pattern: "^[a-z_]+_delete$"},
			want:   []located{mismatch("foo", "DROP TABLE shop_delete.foo, Shop.bar_delete")},
		},
		{
			name:   "temporary tables judged",
			sql:    "DROP TEMPORARY TABLE tmp;\nDROP GLOBAL TEMPORARY TABLE gtmp;",
			config: pattern,
			want: []located{
				mismatch("tmp", "DROP TEMPORARY TABLE tmp"),
				mismatch("gtmp", "DROP GLOBAL TEMPORARY TABLE gtmp"),
			},
		},
		{
			name:   "quoted name judged unquoted",
			sql:    "DROP TABLE `foo_delete`, `foo_delete `;",
			config: pattern,
			want:   []located{mismatch("foo_delete ", "DROP TABLE `foo_delete`, `foo_delete `")},
		},
		{name: "drop view not judged", sql: "DROP VIEW foo;", config: pattern},
		{name: "other statements not judged", sql: "TRUNCATE TABLE foo;\nALTER TABLE foo DROP COLUMN a;\nDROP DATABASE foo;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN DROP TABLE foo; END;", config: pattern},
		{
			name:   "length not judged",
			sql:    "DROP TABLE foo_to_delete;",
			config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 3},
		},
		{name: "length alone requires nothing", sql: "DROP TABLE foo;", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE foo;"},
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

// The parser builds no nil entries, so a hand-built tree pins that the walk
// skips them rather than dereferencing them.
func TestTableDropNamingNilEntries(t *testing.T) {
	stmts := []ast.Node{nil, &ast.DropTableStmt{Tables: []*ast.TableRef{nil, {Name: "foo"}}}}
	got, err := TableDropNaming(stmts, v1review.NamingConfig{Pattern: "_delete$"})
	if err != nil {
		t.Fatalf("TableDropNaming: %v", err)
	}
	if diff := cmp.Diff([]v1review.Finding{{Code: v1review.NameMismatch, Name: "foo"}}, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
