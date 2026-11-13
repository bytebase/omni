package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestTableNoForeignKey(t *testing.T) {
	fk := func(name string, statement int, text string) tableFinding {
		return tableFinding{v1review.TableForeignKey, name, statement, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []tableFinding
	}{
		{name: "no foreign key", sql: "CREATE TABLE t (a INT PRIMARY KEY);"},
		{
			name: "column references anchors the create",
			sql:  "CREATE TABLE t (a INT REFERENCES u (a));",
			want: []tableFinding{fk("S.T", 0, "CREATE TABLE t (a INT REFERENCES u (a))")},
		},
		{
			name: "table constraint",
			sql:  "CREATE TABLE x.t (a INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a));",
			want: []tableFinding{fk("X.T", 0, "CREATE TABLE x.t (a INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a))")},
		},
		{
			name: "one finding per table, anchored on its last foreign key, ordered by first",
			sql:  "CREATE TABLE b (a INT REFERENCES u);\nCREATE TABLE a (a INT REFERENCES u);\nALTER TABLE b ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a);",
			want: []tableFinding{
				fk("S.B", 2, "CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a)"),
				fk("S.A", 1, "CREATE TABLE a (a INT REFERENCES u)"),
			},
		},
		{name: "alter add column references not judged", sql: "ALTER TABLE t ADD b INT REFERENCES u;"},
		{name: "alter add check not judged", sql: "ALTER TABLE t ADD CONSTRAINT c CHECK (a > 0);"},
		{name: "plsql block not judged", sql: "BEGIN EXECUTE IMMEDIATE 'x'; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runTableRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return TableNoForeignKey(stmts, "S") })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
