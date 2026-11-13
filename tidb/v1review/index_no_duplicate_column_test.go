package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestIndexNoDuplicateColumn(t *testing.T) {
	duplicated := func(kind, table, name, column, text string) indexFinding {
		return indexFinding{Code: v1review.IndexColumnDuplicated, Table: table, Name: name, Other: column, Type: kind, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []indexFinding
	}{
		{name: "no duplicate", sql: "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b), KEY k (b, a))"},
		{
			name: "create table anchors each definition, first repeat only",
			sql:  "CREATE TABLE db.t (a INT, b INT,\n  PRIMARY KEY pk (a, b, a),\n  UNIQUE uk (b, a, b, a),\n  KEY (a, a),\n  CONSTRAINT fk FOREIGN KEY (b, b) REFERENCES p (a, b))",
			want: []indexFinding{
				duplicated(v1review.KeyPrimary, "t", "PRIMARY", "a", "PRIMARY KEY pk (a, b, a)"),
				duplicated(v1review.KeyUnique, "t", "uk", "b", "UNIQUE uk (b, a, b, a)"),
				duplicated(v1review.KeyIndex, "t", "", "a", "KEY (a, a)"),
				duplicated(v1review.KeyForeign, "t", "fk", "b", "CONSTRAINT fk FOREIGN KEY (b, b) REFERENCES p (a, b)"),
			},
		},
		{
			name: "columns compare as written; expressions never repeat, a parenthesized column does",
			sql:  "CREATE TABLE t (a INT, KEY k1 (a, A), KEY k2 ((a + 1), (a + 1)), KEY k3 ((a), (a)))",
			want: []indexFinding{duplicated(v1review.KeyIndex, "t", "k3", "a", "KEY k3 ((a), (a))")},
		},
		{
			name: "FULLTEXT definitions and column attributes are not judged",
			sql:  "CREATE TABLE t (a TEXT PRIMARY KEY, FULLTEXT KEY ft (a, a))",
		},
		{
			name: "create index of any kind is INDEX and has no range",
			sql:  "CREATE UNIQUE INDEX uk ON db.t (a, b, a);\nCREATE FULLTEXT INDEX ft ON t (a, a)",
			want: []indexFinding{
				duplicated(v1review.KeyIndex, "t", "uk", "a", ""),
				duplicated(v1review.KeyIndex, "t", "ft", "a", ""),
			},
		},
		{
			name: "alter table additions have no range; the parenthesized ADD form is not judged",
			sql:  "ALTER TABLE t ADD PRIMARY KEY (a, a), ADD FOREIGN KEY fk_idx (b, b) REFERENCES p (a, b), ADD (c INT, KEY k (c, c))",
			want: []indexFinding{
				duplicated(v1review.KeyPrimary, "t", "PRIMARY", "a", ""),
				duplicated(v1review.KeyForeign, "t", "fk_idx", "b", ""),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := IndexNoDuplicateColumn(parseScript(t, tt.sql))
			if diff := cmp.Diff(tt.want, indexFindings(tt.sql, got)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
