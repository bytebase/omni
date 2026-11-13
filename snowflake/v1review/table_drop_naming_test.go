package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/snowflake/parser"
	"github.com/bytebase/omni/v1review"
)

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
	mismatch := func(name, text string) []dropFinding {
		return []dropFinding{{Code: v1review.NameMismatch, Name: name, Text: text}}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []dropFinding
	}{
		{name: "quoted match", sql: `DROP TABLE IF EXISTS "FOO_delete"`, config: pattern},
		{
			name:   "unquoted name upper-cased and anchored at the name",
			sql:    "DROP TABLE IF EXISTS Foo",
			config: pattern,
			want:   mismatch("FOO", "Foo"),
		},
		{
			name:   "unquoted name upper-cased before matching",
			sql:    "DROP TABLE foo_delete",
			config: pattern,
			want:   mismatch("FOO_DELETE", "foo_delete"),
		},
		{name: "upper-case pattern matches an unquoted name", sql: "DROP TABLE foo_delete", config: v1review.NamingConfig{Pattern: "_DELETE$"}},
		{
			name:   "quoted name keeps its case",
			sql:    `DROP TABLE "foo_Delete"`,
			config: pattern,
			want:   mismatch("foo_Delete", `"foo_Delete"`),
		},
		{
			name:   "qualifier not judged, last part anchored",
			sql:    `DROP TABLE "x_delete"."y_delete".Foo CASCADE`,
			config: pattern,
			want:   mismatch("FOO", "Foo"),
		},
		{
			name:   "doubled quote read as one",
			sql:    `DROP TABLE "Foo""Bar"`,
			config: pattern,
			want:   mismatch(`Foo"Bar`, `"Foo""Bar"`),
		},
		{
			name:   "IDENTIFIER name judged as written, upper-cased whole",
			sql:    "DROP TABLE IDENTIFIER('db.sch.foo_delete')",
			config: pattern,
			want:   mismatch("IDENTIFIER('DB.SCH.FOO_DELETE')", "IDENTIFIER('db.sch.foo_delete')"),
		},
		{
			name:   "empty quoted name judged",
			sql:    `DROP TABLE ""`,
			config: pattern,
			want:   mismatch("", `""`),
		},
		{
			name:   "non-ASCII quoted name",
			sql:    `DROP TABLE "数据表é"`,
			config: pattern,
			want:   mismatch("数据表é", `"数据表é"`),
		},
		{name: "unanchored pattern matches anywhere", sql: `DROP TABLE "a_delete_b"`, config: v1review.NamingConfig{Pattern: "_delete"}},
		{name: "length not judged", sql: `DROP TABLE "a_very_long_table_name_delete"`, config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 4}},
		{name: "length alone requires nothing", sql: "DROP TABLE Foo", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE Foo"},
		{
			name:   "name on a later line",
			sql:    "DROP TABLE \"a_delete\";\n\n-- note\nDROP TABLE\n  sch\n  .\n  Foo;",
			config: pattern,
			want:   mismatch("FOO", "Foo"),
		},
		{name: "drop view not judged", sql: "DROP VIEW Foo", config: pattern},
		{name: "drop dynamic table not judged", sql: "DROP DYNAMIC TABLE Foo", config: pattern},
		{name: "drop external table not judged", sql: "DROP EXTERNAL TABLE Foo", config: pattern},
		{name: "truncate not judged", sql: "TRUNCATE TABLE Foo", config: pattern},
		{name: "undrop not judged", sql: "UNDROP TABLE Foo", config: pattern},
		{name: "alter table drop column not judged", sql: "ALTER TABLE Foo DROP COLUMN c", config: pattern},
		{name: "statement under EXPLAIN not judged", sql: "EXPLAIN DROP TABLE Foo", config: pattern},
		{name: "scripting block not judged", sql: "BEGIN DROP TABLE Foo; END;", config: pattern},
		{name: "task body not judged", sql: "CREATE TASK t AS DROP TABLE Foo", config: pattern},
		{name: "execute immediate not judged", sql: "EXECUTE IMMEDIATE 'DROP TABLE Foo'", config: pattern},
	} {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			findings, err := TableDropNaming(file.Stmts, tt.config)
			if err != nil {
				t.Fatalf("TableDropNaming: %v", err)
			}
			var got []dropFinding
			for _, f := range findings {
				got = append(got, dropFinding{f.Code, f.Name, f.Table, f.Pattern, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableDropNamingInvalidPattern(t *testing.T) {
	file, err := parser.Parse("DROP TABLE t")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := TableDropNaming(file.Stmts, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
