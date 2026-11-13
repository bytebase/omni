package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestIndexNoDuplicateColumn(t *testing.T) {
	duplicated := func(table, name, other, kind, search, text string) indexFinding {
		return indexFinding{Code: v1review.IndexColumnDuplicated, Table: table, Name: name, Other: other, Type: kind, Search: search, Text: text}
	}
	const create = `CREATE TABLE s.t (a int, b int, CONSTRAINT pk PRIMARY KEY (a, b, b, a), UNIQUE (b, b), FOREIGN KEY (a, a) REFERENCES p (x, y))`
	for _, tt := range []struct {
		name string
		sql  string
		want []indexFinding
	}{
		{name: "no duplicate", sql: "CREATE TABLE t (a int, b int, PRIMARY KEY (a, b))"},
		{
			name: "one finding per constraint, for the first repeat",
			sql:  create,
			want: []indexFinding{
				duplicated("t", "pk", "b", v1review.KeyPrimary, "PRIMARY KEY", create),
				duplicated("t", "", "b", v1review.KeyUnique, "UNIQUE", create),
				duplicated("t", "", "a", v1review.KeyForeign, "FOREIGN KEY", create),
			},
		},
		{
			name: "unquoted names fold, quoted ones do not",
			sql:  `CREATE TABLE t ("A" int, a int, UNIQUE ("A", a), PRIMARY KEY (A, a))`,
			want: []indexFinding{duplicated("t", "", "a", v1review.KeyPrimary, "PRIMARY KEY", `CREATE TABLE t ("A" int, a int, UNIQUE ("A", a), PRIMARY KEY (A, a))`)},
		},
		{
			name: "alter table adds",
			sql:  "ALTER TABLE t ADD CONSTRAINT uk UNIQUE (a, a), ADD CONSTRAINT ex EXCLUDE USING gist (a WITH =, a WITH =)",
			want: []indexFinding{duplicated("t", "uk", "a", v1review.KeyUnique, "UNIQUE", "ALTER TABLE t ADD CONSTRAINT uk UNIQUE (a, a), ADD CONSTRAINT ex EXCLUDE USING gist (a WITH =, a WITH =)")},
		},
		{
			name: "a unique index is an INDEX",
			sql:  "CREATE UNIQUE INDEX idx ON t (a, (a + 1), (a + 1), a)",
			want: []indexFinding{duplicated("t", "idx", "a", v1review.KeyIndex, "", "CREATE UNIQUE INDEX idx ON t (a, (a + 1), (a + 1), a)")},
		},
		{name: "expressions and INCLUDE columns are never duplicates", sql: "CREATE INDEX idx ON t (a, (a + 1), (a + 1)) INCLUDE (a, a)"},
		{name: "routine bodies are not judged", sql: "CREATE FUNCTION f() RETURNS void AS $$ CREATE INDEX i ON t (a, a) $$ LANGUAGE sql"},
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
