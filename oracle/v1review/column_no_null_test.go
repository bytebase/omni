package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnNoNull(t *testing.T) {
	null := func(name string, statement int, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnNullable, Name: name, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{
			name: "nullable columns in key order",
			sql:  "CREATE TABLE t (b int, a int NOT NULL, c int PRIMARY KEY, d int CONSTRAINT nn NOT NULL, \"A\" int);",
			want: []columnFinding{null("A", 0, `"A" int`), null("B", 0, "b int")},
		},
		{name: "primary key constraint forgets", sql: "CREATE TABLE t (a int, b int, PRIMARY KEY (a, b));"},
		{
			name: "dotted name prints its tail",
			sql:  `CREATE TABLE t ("x.y" int);`,
			want: []columnFinding{null("y", 0, `"x.y" int`)},
		},
		{
			name: "later definition anchors and statement order kept",
			sql:  "ALTER TABLE t ADD a int;\nALTER TABLE t MODIFY a varchar2(10);",
			want: []columnFinding{null("A", 1, "a varchar2(10)")},
		},
		{name: "later not null forgets", sql: "ALTER TABLE t ADD a int;\nALTER TABLE t MODIFY a NOT NULL;"},
		{
			name: "drop does not forget",
			sql:  "ALTER TABLE t ADD a int;\nALTER TABLE t DROP COLUMN a;\nDROP TABLE t;",
			want: []columnFinding{null("A", 0, "a int")},
		},
		{
			name: "current database qualifies keys",
			sql:  "ALTER TABLE t ADD a int;\nALTER TABLE db.t ADD a int NOT NULL;\nALTER TABLE s.t ADD a int;",
			want: []columnFinding{null("A", 2, "a int")},
		},
		{
			name: "key order not statement order",
			sql:  "CREATE TABLE z (a int);\nCREATE TABLE b (a int);",
			want: []columnFinding{null("A", 1, "a int"), null("A", 0, "a int")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnNoNull(parseStmts(t, tt.sql), "DB"))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// The parser records NOT NULL as ColumnDef.NotNull, never as a constraint;
// a NOT NULL constraint node still makes a column non-nullable.
func TestColumnNoNullConstraintNode(t *testing.T) {
	col := &ast.ColumnDef{Name: "A", Constraints: &ast.List{Items: []ast.Node{&ast.ColumnConstraint{Type: ast.CONSTRAINT_NOT_NULL}}}}
	stmt := &ast.CreateTableStmt{Name: &ast.ObjectName{Name: "T"}, Columns: &ast.List{Items: []ast.Node{col}}}
	if got := ColumnNoNull([]ast.Node{stmt}, "DB"); len(got) != 0 {
		t.Errorf("got %v, want no findings", got)
	}
}
