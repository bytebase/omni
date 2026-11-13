package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

func runTableDropNaming(t *testing.T, sql string, cfg v1review.NamingConfig) []finding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := TableDropNaming(list.Items, cfg)
	if err != nil {
		t.Fatalf("TableDropNaming: %v", err)
	}
	var findings []finding
	for _, f := range got {
		if f.Table != "" || f.Pattern != "" || f.Detail != "" {
			t.Errorf("finding %+v: Table, Pattern, and Detail should be empty", f)
		}
		findings = append(findings, finding{f.Code, f.Name, sql[f.Start:f.End]})
	}
	return findings
}

func TestTableDropNaming(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "_delete$"}
	mismatch := func(name, text string) finding { return finding{v1review.NameMismatch, name, text} }
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []finding
	}{
		{name: "drop ok", sql: "DROP TABLE IF EXISTS foo_delete;", config: pattern},
		{
			name:   "drop mismatch",
			sql:    "DROP TABLE IF EXISTS foo;",
			config: pattern,
			want:   []finding{mismatch("foo", "DROP TABLE IF EXISTS foo")},
		},
		{
			name:   "every listed table judged in order",
			sql:    `DROP TABLE foo, bar_delete, "Baz" CASCADE;`,
			config: pattern,
			want: []finding{
				mismatch("foo", `DROP TABLE foo, bar_delete, "Baz" CASCADE`),
				mismatch("Baz", `DROP TABLE foo, bar_delete, "Baz" CASCADE`),
			},
		},
		{
			name:   "schema not judged",
			sql:    `DROP TABLE old_delete.foo, s.bar_delete;`,
			config: pattern,
			want:   []finding{mismatch("foo", `DROP TABLE old_delete.foo, s.bar_delete`)},
		},
		{name: "three-part name skipped", sql: `DROP TABLE db.s.foo;`, config: pattern},
		{
			name:   "unquoted name folds to lower",
			sql:    `DROP TABLE Foo_DELETE; DROP TABLE "Foo_DELETE";`,
			config: pattern,
			want:   []finding{mismatch("Foo_DELETE", `DROP TABLE "Foo_DELETE"`)},
		},
		{
			name:   "name truncated to 63 bytes before judging",
			sql:    `DROP TABLE ` + strings.Repeat("a", 60) + `_delete;`,
			config: pattern,
			want:   []finding{mismatch(strings.Repeat("a", 60)+"_de", `DROP TABLE `+strings.Repeat("a", 60)+`_delete`)},
		},
		{
			name:   "other drops not judged",
			sql:    `DROP VIEW foo; DROP MATERIALIZED VIEW foo; DROP FOREIGN TABLE foo; DROP INDEX foo; DROP SEQUENCE foo; DROP SCHEMA foo; TRUNCATE TABLE foo;`,
			config: pattern,
		},
		{
			name:   "routine body not judged",
			sql:    `CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC DROP TABLE foo; END;`,
			config: pattern,
		},
		{name: "length not consulted", sql: "DROP TABLE a_very_long_table_name_delete;", config: v1review.NamingConfig{Pattern: "_delete$", MaxLength: 3}},
		{name: "length alone requires nothing", sql: "DROP TABLE foo;", config: v1review.NamingConfig{MaxLength: 1}},
		{name: "empty config requires nothing", sql: "DROP TABLE foo;"},
		{
			name:   "findings in statement order",
			sql:    "DROP TABLE a;\n-- c\nDROP TABLE b_delete;\nDROP TABLE c;",
			config: pattern,
			want:   []finding{mismatch("a", "DROP TABLE a"), mismatch("c", "DROP TABLE c")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runTableDropNaming(t, tt.sql, tt.config)); diff != "" {
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
