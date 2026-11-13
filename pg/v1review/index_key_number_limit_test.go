package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// indexFinding is an index rule's finding as the tests compare it, with the
// text its range anchors in the statement it lies in.
type indexFinding struct {
	Code                     v1review.Code
	Table, Name, Other, Type string
	Search                   string
	Text                     string
}

// parseEach parses each statement on its own, as Bytebase does, so a
// finding's range is an offset into its own statement.
func parseEach(t *testing.T, sqls []string) []ast.Node {
	t.Helper()
	var nodes []ast.Node
	for _, sql := range sqls {
		list, err := parser.Parse(sql)
		if err != nil {
			t.Fatalf("Parse %q: %v", sql, err)
		}
		if len(list.Items) != 1 {
			t.Fatalf("Parse %q: %d statements, want 1", sql, len(list.Items))
		}
		nodes = append(nodes, list.Items[0])
	}
	return nodes
}

func indexFindings(sqls []string, got []v1review.Finding) []indexFinding {
	var findings []indexFinding
	for _, f := range got {
		findings = append(findings, indexFinding{f.Code, f.Table, f.Name, f.Other, f.Value, f.Search, sqls[f.Statement][f.Start:f.End]})
	}
	return findings
}

func TestIndexKeyNumberLimit(t *testing.T) {
	exceeded := func(table, name, kind, text string) indexFinding {
		return indexFinding{Code: v1review.IndexKeyCountExceeded, Table: table, Name: name, Type: kind, Text: text}
	}
	for _, tt := range []struct {
		name  string
		sql   string
		limit int
		want  []indexFinding
	}{
		{name: "within the limit", sql: "CREATE TABLE t (a int, b int, PRIMARY KEY (a, b))", limit: 2},
		{
			name:  "table constraints",
			sql:   "CREATE TABLE s.t (a int, b int, c int, CONSTRAINT pk PRIMARY KEY (a, b, c), UNIQUE (a, b, c), FOREIGN KEY (a, b, c) REFERENCES p (x, y, z))",
			limit: 2,
			want: []indexFinding{
				exceeded("t", "pk", v1review.KeyPrimary, "CREATE TABLE s.t (a int, b int, c int, CONSTRAINT pk PRIMARY KEY (a, b, c), UNIQUE (a, b, c), FOREIGN KEY (a, b, c) REFERENCES p (x, y, z))"),
				exceeded("t", "", v1review.KeyUnique, "CREATE TABLE s.t (a int, b int, c int, CONSTRAINT pk PRIMARY KEY (a, b, c), UNIQUE (a, b, c), FOREIGN KEY (a, b, c) REFERENCES p (x, y, z))"),
				exceeded("t", "", v1review.KeyForeign, "CREATE TABLE s.t (a int, b int, c int, CONSTRAINT pk PRIMARY KEY (a, b, c), UNIQUE (a, b, c), FOREIGN KEY (a, b, c) REFERENCES p (x, y, z))"),
			},
		},
		{
			name:  "column constraints, EXCLUDE and INCLUDE are not counted",
			sql:   "CREATE TABLE t (a int PRIMARY KEY, b int UNIQUE, EXCLUDE USING gist (a WITH =, b WITH =, a WITH =), UNIQUE (a) INCLUDE (a, b, a))",
			limit: 1,
		},
		{
			name:  "alter table adds",
			sql:   "ALTER TABLE t ADD CONSTRAINT uk UNIQUE (a, b), ADD COLUMN c int PRIMARY KEY, ADD PRIMARY KEY USING INDEX i",
			limit: 1,
			want:  []indexFinding{exceeded("t", "uk", v1review.KeyUnique, "ALTER TABLE t ADD CONSTRAINT uk UNIQUE (a, b), ADD COLUMN c int PRIMARY KEY, ADD PRIMARY KEY USING INDEX i")},
		},
		{
			name:  "create index counts plain columns only",
			sql:   "CREATE UNIQUE INDEX idx ON t (a, (b + 1), lower(c), d) INCLUDE (e, f)",
			limit: 1,
			want:  []indexFinding{exceeded("t", "idx", v1review.KeyIndex, "CREATE UNIQUE INDEX idx ON t (a, (b + 1), lower(c), d) INCLUDE (e, f)")},
		},
		{name: "expressions within the limit", sql: "CREATE INDEX idx ON t (a, (b + 1), lower(c))", limit: 1},
		{name: "a limit of 0 reports nothing", sql: "CREATE INDEX idx ON t (a, b, c)", limit: 0},
		{name: "a negative limit reports nothing", sql: "CREATE INDEX idx ON t (a, b, c)", limit: -1},
		{name: "routine bodies are not judged", sql: "CREATE FUNCTION f() RETURNS void AS $$ CREATE INDEX i ON t (a, b) $$ LANGUAGE sql", limit: 1},
		{name: "create schema elements are not judged", sql: "CREATE SCHEMA s CREATE INDEX i ON t (a, b)", limit: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sqls := []string{tt.sql}
			got := indexFindings(sqls, IndexKeyNumberLimit(parseEach(t, sqls), tt.limit))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
