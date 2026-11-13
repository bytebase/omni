package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
	"github.com/bytebase/omni/v1review"
)

// indexFinding is an index rule's finding as the tests compare it, with the
// text its range anchors in the script: empty when it has no range.
type indexFinding struct {
	Code               v1review.Code
	Table, Name, Other string
	Type               string
	Statement          int
	Text               string
}

func parseScript(t *testing.T, sql string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse %q: %v", sql, err)
	}
	return list.Items
}

func indexFindings(sql string, got []v1review.Finding) []indexFinding {
	var findings []indexFinding
	for _, f := range got {
		findings = append(findings, indexFinding{f.Code, f.Table, f.Name, f.Other, f.Value, f.Statement, sql[f.Start:f.End]})
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
		{name: "within the limit", sql: "CREATE TABLE t (a INT, b INT, INDEX idx_ab (a, b))", limit: 2},
		{
			name:  "create table anchors each definition",
			sql:   "CREATE TABLE db.t (a INT, b INT, c INT,\n  PRIMARY KEY (a, b, c),\n  UNIQUE KEY uk (a, b, c),\n  KEY idx (a, b, c),\n  CONSTRAINT fk FOREIGN KEY fk_idx (a, b, c) REFERENCES p (a, b, c))",
			limit: 2,
			want: []indexFinding{
				exceeded("t", "PRIMARY", v1review.KeyPrimary, "PRIMARY KEY (a, b, c)"),
				exceeded("t", "uk", v1review.KeyUnique, "UNIQUE KEY uk (a, b, c)"),
				exceeded("t", "idx", v1review.KeyIndex, "KEY idx (a, b, c)"),
				exceeded("t", "fk", v1review.KeyForeign, "CONSTRAINT fk FOREIGN KEY fk_idx (a, b, c) REFERENCES p (a, b, c)"),
			},
		},
		{
			name:  "a primary key is PRIMARY unless CONSTRAINT names it",
			sql:   "CREATE TABLE t (a INT, b INT, PRIMARY KEY pk_ab (a, b), CONSTRAINT pk_c PRIMARY KEY (a, b))",
			limit: 1,
			want: []indexFinding{
				exceeded("t", "PRIMARY", v1review.KeyPrimary, "PRIMARY KEY pk_ab (a, b)"),
				exceeded("t", "pk_c", v1review.KeyPrimary, "CONSTRAINT pk_c PRIMARY KEY (a, b)"),
			},
		},
		{
			name:  "a foreign key without a constraint name takes its index name",
			sql:   "CREATE TABLE t (a INT, b INT, FOREIGN KEY fk_idx (a, b) REFERENCES p (a, b), FOREIGN KEY (a, b) REFERENCES p (a, b))",
			limit: 1,
			want: []indexFinding{
				exceeded("t", "fk_idx", v1review.KeyForeign, "FOREIGN KEY fk_idx (a, b) REFERENCES p (a, b)"),
				exceeded("t", "", v1review.KeyForeign, "FOREIGN KEY (a, b) REFERENCES p (a, b)"),
			},
		},
		{
			name:  "expression key parts count",
			sql:   "CREATE TABLE t (a INT, b INT, INDEX idx ((a + 1), (lower(b))))",
			limit: 1,
			want:  []indexFinding{exceeded("t", "idx", v1review.KeyIndex, "INDEX idx ((a + 1), (lower(b)))")},
		},
		{
			name:  "column attributes and FULLTEXT and VECTOR definitions are not judged",
			sql:   "CREATE TABLE t (a TEXT PRIMARY KEY, b TEXT UNIQUE, v VECTOR(3), FULLTEXT KEY ft (a, b), VECTOR INDEX vi ((VEC_COSINE_DISTANCE(v))))",
			limit: 0,
		},
		{
			name:  "create index of any kind counts every key part and has no range",
			sql:   "CREATE INDEX idx ON db.t (a, (a + 1));\nCREATE UNIQUE INDEX uk ON t (a, b);\nCREATE FULLTEXT INDEX ft ON t (a, b)",
			limit: 1,
			want: []indexFinding{
				exceeded("t", "idx", v1review.KeyIndex, ""),
				exceeded("t", "uk", v1review.KeyIndex, ""),
				exceeded("t", "ft", v1review.KeyIndex, ""),
			},
		},
		{
			name:  "alter table additions have no range; the parenthesized ADD form is not judged",
			sql:   "ALTER TABLE db.t ADD INDEX idx (a, b), ADD CONSTRAINT uk UNIQUE (a, b), ADD (c INT, KEY k (a, b, c)), ADD COLUMN d INT PRIMARY KEY",
			limit: 1,
			want: []indexFinding{
				exceeded("t", "idx", v1review.KeyIndex, ""),
				exceeded("t", "uk", v1review.KeyUnique, ""),
			},
		},
		{name: "a limit of zero reports nothing", sql: "CREATE INDEX idx ON t (a, b)", limit: 0},
		{name: "a negative limit reports nothing", sql: "CREATE INDEX idx ON t (a, b)", limit: -1},
		{name: "routine bodies are not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE INDEX idx ON t (a, b); END", limit: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := IndexKeyNumberLimit(parseScript(t, tt.sql), tt.limit)
			if diff := cmp.Diff(tt.want, indexFindings(tt.sql, got)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
