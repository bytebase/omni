package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
	"github.com/bytebase/omni/v1review"
)

// finding is a finding as the tests compare it: its code, name, table and
// pattern, and the text its range anchors.
type finding struct {
	Code    v1review.Code
	Name    string
	Table   string
	Pattern string
	Text    string
}

func TestTableDropNaming(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "_delete$"}
	mismatch := func(name, text string) finding {
		return finding{Code: v1review.NameMismatch, Name: name, Text: text}
	}
	long := strings.Repeat("a", 60) + "_delete"
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []finding
	}{
		{name: "match", sql: "DROP TABLE IF EXISTS foo_delete;", config: pattern},
		{
			name:   "mismatch anchors the statement",
			sql:    "DROP TABLE active;",
			config: pattern,
			want:   []finding{mismatch("active", "DROP TABLE active")},
		},
		{name: "unquoted name folded before matching", sql: "DROP TABLE Foo_DELETE;", config: pattern},
		{
			name:   "quoted name keeps its case",
			sql:    `DROP TABLE "Foo_DELETE";`,
			config: pattern,
			want:   []finding{mismatch("Foo_DELETE", `DROP TABLE "Foo_DELETE"`)},
		},
		{
			name:   "every listed table judged in order",
			sql:    "DROP TABLE Foo, bar_delete, Baz CASCADE;",
			config: pattern,
			want: []finding{
				mismatch("foo", "DROP TABLE Foo, bar_delete, Baz CASCADE"),
				mismatch("baz", "DROP TABLE Foo, bar_delete, Baz CASCADE"),
			},
		},
		{
			name:   "last part judged whatever the qualifier",
			sql:    "DROP TABLE x_delete.Foo, x_delete.y_delete.Bar, a.b.c.Qux, a.b.c.d_delete;",
			config: pattern,
			want: []finding{
				mismatch("foo", "DROP TABLE x_delete.Foo, x_delete.y_delete.Bar, a.b.c.Qux, a.b.c.d_delete"),
				mismatch("bar", "DROP TABLE x_delete.Foo, x_delete.y_delete.Bar, a.b.c.Qux, a.b.c.d_delete"),
				mismatch("qux", "DROP TABLE x_delete.Foo, x_delete.y_delete.Bar, a.b.c.Qux, a.b.c.d_delete"),
			},
		},
		{
			name:   "doubled quote read as one",
			sql:    `DROP TABLE "Foo""Bar";`,
			config: pattern,
			want:   []finding{mismatch(`Foo"Bar`, `DROP TABLE "Foo""Bar"`)},
		},
		{
			name:   "unquoted name truncated to 63 bytes",
			sql:    "DROP TABLE " + long + ";",
			config: pattern,
			want:   []finding{mismatch(long[:63], "DROP TABLE "+long)},
		},
		{name: "quoted name not truncated", sql: `DROP TABLE "` + long + `";`, config: pattern},
		{
			name:   "non-ASCII names",
			sql:    `DROP TABLE "数据表", ÀÉ;`,
			config: pattern,
			want: []finding{
				mismatch("数据表", `DROP TABLE "数据表", ÀÉ`),
				mismatch("ÀÉ", `DROP TABLE "数据表", ÀÉ`),
			},
		},
		{
			name:   "unanchored pattern matches anywhere",
			sql:    "DROP TABLE a_delete_b, c;",
			config: v1review.NamingConfig{Pattern: "_delete"},
			want:   []finding{mismatch("c", "DROP TABLE a_delete_b, c")},
		},
		{name: "length not judged", sql: "DROP TABLE a_very_long_table_name_delete;", config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 4}},
		{name: "length alone requires nothing", sql: "DROP TABLE foo;", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE foo;"},
		{
			name:   "statement on a later line",
			sql:    "DROP TABLE a_delete;\n\n-- note\nDROP TABLE\n  foo;",
			config: pattern,
			want:   []finding{mismatch("foo", "DROP TABLE\n  foo")},
		},
		{name: "drop view not judged", sql: "DROP VIEW foo;", config: pattern},
		{name: "drop materialized view not judged", sql: "DROP MATERIALIZED VIEW foo;", config: pattern},
		{name: "drop foreign table not judged", sql: "DROP FOREIGN TABLE foo;", config: pattern},
		{name: "drop schema not judged", sql: "DROP SCHEMA foo;", config: pattern},
		{name: "truncate not judged", sql: "TRUNCATE TABLE foo;", config: pattern},
		{name: "alter table drop column not judged", sql: "ALTER TABLE foo DROP COLUMN c;", config: pattern},
		{
			name:   "routine body not judged",
			sql:    "CREATE OR REPLACE PROCEDURE p() AS $$ BEGIN DROP TABLE foo; END; $$ LANGUAGE plpgsql;",
			config: pattern,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			findings, err := TableDropNaming(list.Items, tt.config)
			if err != nil {
				t.Fatalf("TableDropNaming: %v", err)
			}
			var got []finding
			for _, f := range findings {
				got = append(got, finding{f.Code, f.Name, f.Table, f.Pattern, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// The parser writes every dropped table as a list of name parts; the V1
// advisor also read a RangeVar, and skipped an object with no name.
func TestTableDropNamingObjectForms(t *testing.T) {
	drop := &ast.DropStmt{
		RemoveType: int(ast.OBJECT_TABLE),
		Objects: &ast.List{Items: []ast.Node{
			&ast.RangeVar{Schemaname: "s_delete", Relname: "ranged"},
			&ast.List{},
			&ast.List{Items: []ast.Node{&ast.String{Str: "s"}, &ast.String{Str: ""}}},
			&ast.RangeVar{},
			&ast.String{Str: "bare"},
		}},
		Loc: ast.Loc{Start: -1, End: -1},
	}
	got, err := TableDropNaming([]ast.Node{drop}, v1review.NamingConfig{Pattern: "_delete$"})
	if err != nil {
		t.Fatalf("TableDropNaming: %v", err)
	}
	want := []v1review.Finding{{Code: v1review.NameMismatch, Name: "ranged"}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
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
