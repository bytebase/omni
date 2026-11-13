package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

// indexFinding is an index rule's finding as the tests compare it, with the
// text its range anchors in the statement it lies in.
type indexFinding struct {
	Code               v1review.Code
	Table, Name, Other string
	Type               string
	Text               string
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
		findings = append(findings, indexFinding{f.Code, f.Table, f.Name, f.Other, f.Value, sqls[f.Statement][f.Start:f.End]})
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
				exceeded("t", "", v1review.KeyPrimary, "PRIMARY KEY (a, b, c)"),
				exceeded("t", "uk", v1review.KeyUnique, "UNIQUE KEY uk (a, b, c)"),
				exceeded("t", "idx", v1review.KeyIndex, "KEY idx (a, b, c)"),
				exceeded("t", "fk", v1review.KeyForeign, "CONSTRAINT fk FOREIGN KEY fk_idx (a, b, c) REFERENCES p (a, b, c)"),
			},
		},
		{
			name:  "a foreign key without a constraint name takes its index name",
			sql:   "CREATE TABLE t (a INT, b INT, FOREIGN KEY fk_idx (a, b) REFERENCES p (a, b))",
			limit: 1,
			want:  []indexFinding{exceeded("t", "fk_idx", v1review.KeyForeign, "FOREIGN KEY fk_idx (a, b) REFERENCES p (a, b)")},
		},
		{
			name:  "a name after PRIMARY KEY is not kept",
			sql:   "CREATE TABLE t (a INT, b INT, PRIMARY KEY pk_ab (a, b))",
			limit: 1,
			want:  []indexFinding{exceeded("t", "", v1review.KeyPrimary, "PRIMARY KEY pk_ab (a, b)")},
		},
		{
			name:  "a constraint name on a primary key is kept",
			sql:   "CREATE TABLE t (a INT, b INT, CONSTRAINT pk_c PRIMARY KEY (a, b))",
			limit: 1,
			want:  []indexFinding{exceeded("t", "pk_c", v1review.KeyPrimary, "CONSTRAINT pk_c PRIMARY KEY (a, b)")},
		},
		{
			name:  "alter table anchors the statement, once per key",
			sql:   "ALTER TABLE t ADD INDEX idx (a, b), ADD CONSTRAINT uk UNIQUE (a, b), ADD COLUMN c INT",
			limit: 1,
			want: []indexFinding{
				exceeded("t", "idx", v1review.KeyIndex, "ALTER TABLE t ADD INDEX idx (a, b), ADD CONSTRAINT uk UNIQUE (a, b), ADD COLUMN c INT"),
				exceeded("t", "uk", v1review.KeyUnique, "ALTER TABLE t ADD INDEX idx (a, b), ADD CONSTRAINT uk UNIQUE (a, b), ADD COLUMN c INT"),
			},
		},
		{
			name:  "create index",
			sql:   "CREATE UNIQUE INDEX idx ON db.t (a, b(10), c DESC)",
			limit: 2,
			want:  []indexFinding{exceeded("t", "idx", v1review.KeyUniqueIndex, "CREATE UNIQUE INDEX idx ON db.t (a, b(10), c DESC)")},
		},
		{name: "expression key parts are not counted", sql: "CREATE TABLE t (a INT, b INT, INDEX idx (a, (b + 1), (a * 2)))", limit: 1},
		{
			name:  "create index counts a parenthesized column",
			sql:   "CREATE INDEX idx ON t (a, (b), (a + b))",
			limit: 1,
			want:  []indexFinding{exceeded("t", "idx", v1review.KeyIndex, "CREATE INDEX idx ON t (a, (b), (a + b))")},
		},
		{name: "a parenthesized column of a definition is not counted", sql: "CREATE TABLE t (a INT, b INT, INDEX idx (a, (b), (a)))", limit: 1},
		{
			name:  "local and global indexes",
			sql:   "ALTER TABLE t ADD INDEX l (a, b) LOCAL, ADD UNIQUE INDEX g (a, b) GLOBAL",
			limit: 1,
			want: []indexFinding{
				exceeded("t", "l", v1review.KeyIndex, "ALTER TABLE t ADD INDEX l (a, b) LOCAL, ADD UNIQUE INDEX g (a, b) GLOBAL"),
				exceeded("t", "g", v1review.KeyUnique, "ALTER TABLE t ADD INDEX l (a, b) LOCAL, ADD UNIQUE INDEX g (a, b) GLOBAL"),
			},
		},
		{
			name:  "create vector index is judged as an index",
			sql:   "CREATE VECTOR INDEX v ON t (a, b) WITH (distance=l2, type=hnsw)",
			limit: 1,
			want:  []indexFinding{exceeded("t", "v", v1review.KeyIndex, "CREATE VECTOR INDEX v ON t (a, b) WITH (distance=l2, type=hnsw)")},
		},
		{name: "a vector key of a create table is not judged", sql: "CREATE TABLE t (a VECTOR(3), b INT, VECTOR INDEX v (a, b))", limit: 1},
		{name: "fulltext and spatial are not judged", sql: "CREATE TABLE t (a TEXT, b TEXT, FULLTEXT KEY ft (a, b), SPATIAL KEY sp (a, b))", limit: 1},
		{name: "create fulltext index is not judged", sql: "CREATE FULLTEXT INDEX ft ON t (a, b)", limit: 1},
		{name: "a column's own key is not judged", sql: "CREATE TABLE t (a INT PRIMARY KEY, b INT UNIQUE)", limit: 0},
		{name: "limit 0 reports nothing", sql: "CREATE INDEX idx ON t (a, b, c)", limit: 0},
		{name: "negative limit reports nothing", sql: "CREATE INDEX idx ON t (a, b, c)", limit: -1},
		{
			name:  "a routine body is not entered",
			sql:   "CREATE PROCEDURE p() BEGIN CREATE INDEX idx ON t (a, b); END",
			limit: 1,
		},
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
