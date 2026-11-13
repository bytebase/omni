package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func runTableDropNaming(t *testing.T, sql string, cfg v1review.NamingConfig) []v1review.Finding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := TableDropNaming(list.Items, cfg)
	if err != nil {
		t.Fatalf("TableDropNaming: %v", err)
	}
	return got
}

func TestTableDropNaming(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "_delete$"}
	// dropFinding is a finding as the tests compare it: its code, name,
	// table and pattern, and the text its range anchors.
	type dropFinding struct {
		Code    v1review.Code
		Name    string
		Table   string
		Pattern string
		Text    string
	}
	mismatch := func(name, text string) dropFinding {
		return dropFinding{Code: v1review.NameMismatch, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []dropFinding
	}{
		{name: "match", sql: "DROP TABLE IF EXISTS [Foo_delete];", config: pattern},
		{
			name:   "mismatch is lower-cased and anchors the statement",
			sql:    "DROP TABLE Foo;",
			config: pattern,
			want:   []dropFinding{mismatch("foo", "DROP TABLE Foo")},
		},
		{name: "name lower-cased before matching", sql: "DROP TABLE Foo_DELETE;", config: pattern},
		{
			name:   "upper-case pattern does not match the lower-cased name",
			sql:    "DROP TABLE Foo_DELETE;",
			config: v1review.NamingConfig{Pattern: "_DELETE$"},
			want:   []dropFinding{mismatch("foo_delete", "DROP TABLE Foo_DELETE")},
		},
		{
			name:   "every listed table judged in order",
			sql:    "DROP TABLE [Foo], [Bar_delete], Baz;",
			config: pattern,
			want: []dropFinding{
				mismatch("foo", "DROP TABLE [Foo], [Bar_delete], Baz"),
				mismatch("baz", "DROP TABLE [Foo], [Bar_delete], Baz"),
			},
		},
		{
			name:   "qualifier not judged",
			sql:    "DROP TABLE x_delete.y_delete.Foo, srv.db.dbo.Bar, db..Baz_delete;",
			config: pattern,
			want: []dropFinding{
				mismatch("foo", "DROP TABLE x_delete.y_delete.Foo, srv.db.dbo.Bar, db..Baz_delete"),
				mismatch("bar", "DROP TABLE x_delete.y_delete.Foo, srv.db.dbo.Bar, db..Baz_delete"),
			},
		},
		{
			name:   "bracketed and quoted names judged unquoted",
			sql:    `DROP TABLE [Bad]]Name], "Q uoted";`,
			config: pattern,
			want: []dropFinding{
				mismatch("bad]name", `DROP TABLE [Bad]]Name], "Q uoted"`),
				mismatch("q uoted", `DROP TABLE [Bad]]Name], "Q uoted"`),
			},
		},
		{
			name:   "temporary table judged with its prefix",
			sql:    "DROP TABLE #Tmp, ##Global_delete;",
			config: v1review.NamingConfig{Pattern: "^[a-z_]+$"},
			want: []dropFinding{
				mismatch("#tmp", "DROP TABLE #Tmp, ##Global_delete"),
				mismatch("##global_delete", "DROP TABLE #Tmp, ##Global_delete"),
			},
		},
		{
			name:   "non-ASCII name lower-cased",
			sql:    "DROP TABLE [ÀÉ数据表];",
			config: pattern,
			want:   []dropFinding{mismatch("àé数据表", "DROP TABLE [ÀÉ数据表]")},
		},
		{
			name:   "unanchored pattern matches anywhere",
			sql:    "DROP TABLE a_delete_b, c;",
			config: v1review.NamingConfig{Pattern: "_delete"},
			want:   []dropFinding{mismatch("c", "DROP TABLE a_delete_b, c")},
		},
		{
			name:   "length not judged",
			sql:    "DROP TABLE a_very_long_table_name_delete;",
			config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 4},
		},
		{name: "length alone requires nothing", sql: "DROP TABLE Foo;", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE Foo;"},
		{
			name:   "statement on a later line",
			sql:    "DROP TABLE a_delete;\n\n-- note\nDROP TABLE\n  Foo;",
			config: pattern,
			want:   []dropFinding{mismatch("foo", "DROP TABLE\n  Foo")},
		},
		{name: "drop view not judged", sql: "DROP VIEW Foo;", config: pattern},
		{name: "drop index not judged", sql: "DROP INDEX ix ON Foo;", config: pattern},
		{name: "drop external table not judged", sql: "DROP EXTERNAL TABLE Foo;", config: pattern},
		{name: "truncate not judged", sql: "TRUNCATE TABLE Foo;", config: pattern},
		{name: "alter table drop column not judged", sql: "ALTER TABLE Foo DROP COLUMN c;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr AS DROP TABLE Foo;", config: pattern},
		{name: "if statement not judged", sql: "IF 1 = 1 DROP TABLE Foo;", config: pattern},
		{name: "begin block not judged", sql: "BEGIN DROP TABLE Foo; END", config: pattern},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got []dropFinding
			for _, f := range runTableDropNaming(t, tt.sql, tt.config) {
				got = append(got, dropFinding{f.Code, f.Name, f.Table, f.Pattern, tt.sql[f.Start:f.End]})
			}
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
