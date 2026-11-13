package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestIndexNoDuplicateColumn(t *testing.T) {
	dup := func(table, name, kind, column, text string) indexFinding {
		return indexFinding{Code: v1review.IndexColumnDuplicated, Table: table, Name: name, Other: column, Type: kind, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []indexFinding
	}{
		{name: "no duplicate", sql: "CREATE TABLE t (a INT, b INT, KEY idx (a, b))"},
		{
			name: "create table, one finding per key, the first repeat",
			sql:  "CREATE TABLE db.t (a INT, b INT,\n  PRIMARY KEY (a, b, a),\n  UNIQUE KEY uk (b, a, a, b),\n  KEY idx (a, a),\n  FOREIGN KEY fk (b, b) REFERENCES p (a, b))",
			want: []indexFinding{
				dup("t", "", v1review.KeyPrimary, "a", "PRIMARY KEY (a, b, a)"),
				dup("t", "uk", v1review.KeyUnique, "a", "UNIQUE KEY uk (b, a, a, b)"),
				dup("t", "idx", v1review.KeyIndex, "a", "KEY idx (a, a)"),
				dup("t", "fk", v1review.KeyForeign, "b", "FOREIGN KEY fk (b, b) REFERENCES p (a, b)"),
			},
		},
		{name: "letter case differs", sql: "CREATE TABLE t (a INT, KEY idx (a, A))"},
		{name: "expressions are not columns", sql: "CREATE TABLE t (a INT, KEY idx ((a + 1), (a + 1)))"},
		{name: "fulltext and spatial are not judged", sql: "CREATE TABLE t (a TEXT, FULLTEXT KEY ft (a, a))"},
		{
			name: "alter table anchors the statement",
			sql:  "ALTER TABLE t ADD UNIQUE INDEX uk (a, a), ADD CONSTRAINT c PRIMARY KEY (b, b)",
			want: []indexFinding{
				dup("t", "uk", v1review.KeyUnique, "a", "ALTER TABLE t ADD UNIQUE INDEX uk (a, a), ADD CONSTRAINT c PRIMARY KEY (b, b)"),
				dup("t", "c", v1review.KeyPrimary, "b", "ALTER TABLE t ADD UNIQUE INDEX uk (a, a), ADD CONSTRAINT c PRIMARY KEY (b, b)"),
			},
		},
		{
			name: "create unique index",
			sql:  "CREATE UNIQUE INDEX uk ON t (a, (a))",
			want: []indexFinding{dup("t", "uk", v1review.KeyUniqueIndex, "a", "CREATE UNIQUE INDEX uk ON t (a, (a))")},
		},
		{
			name: "create index",
			sql:  "CREATE INDEX idx ON t (a(10), a DESC)",
			want: []indexFinding{dup("t", "idx", v1review.KeyIndex, "a", "CREATE INDEX idx ON t (a(10), a DESC)")},
		},
		{name: "create fulltext index is not judged", sql: "CREATE FULLTEXT INDEX ft ON t (a, a)"},
		{name: "a routine body is not entered", sql: "CREATE PROCEDURE p() BEGIN CREATE INDEX idx ON t (a, a); END"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sqls := []string{tt.sql}
			got := indexFindings(sqls, IndexNoDuplicateColumn(parseEach(t, sqls)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
