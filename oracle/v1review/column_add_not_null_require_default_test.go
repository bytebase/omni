package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnAddNotNullRequireDefault(t *testing.T) {
	nn := func(name, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnNotNullNoDefault, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "add not null", sql: "ALTER TABLE t ADD a int NOT NULL;", want: []columnFinding{nn("A", "a int NOT NULL")}},
		{name: "default satisfies", sql: "ALTER TABLE t ADD a int DEFAULT 0 NOT NULL;"},
		{name: "nullable not judged", sql: "ALTER TABLE t ADD a int;"},
		{
			name: "named not null constraint",
			sql:  "ALTER TABLE t ADD a int CONSTRAINT nn NOT NULL;",
			want: []columnFinding{nn("A", "a int CONSTRAINT nn NOT NULL")},
		},
		{name: "primary key alone not judged", sql: "ALTER TABLE t ADD a int PRIMARY KEY;"},
		{
			name: "column list in order",
			sql:  "ALTER TABLE t ADD (\"b\" int NOT NULL, c int, a int NOT NULL);",
			want: []columnFinding{nn("b", "\"b\" int NOT NULL"), nn("A", "a int NOT NULL")},
		},
		{name: "modify not judged", sql: "ALTER TABLE t MODIFY a int NOT NULL;"},
		{name: "create table not judged", sql: "CREATE TABLE t (a int NOT NULL);"},
		{
			name: "statements in order",
			sql:  "ALTER TABLE t ADD a int NOT NULL;\nALTER TABLE u ADD b int NOT NULL;",
			want: []columnFinding{nn("A", "a int NOT NULL"), nn("B", "b int NOT NULL")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnAddNotNullRequireDefault(parseStmts(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// The parser records NOT NULL as ColumnDef.NotNull, never as a constraint;
// a NOT NULL constraint node is still judged.
func TestColumnAddNotNullRequireDefaultConstraintNode(t *testing.T) {
	col := &ast.ColumnDef{Name: "A", Constraints: &ast.List{Items: []ast.Node{&ast.ColumnConstraint{Type: ast.CONSTRAINT_NOT_NULL}}}, Loc: ast.Loc{Start: 3, End: 5}}
	stmt := &ast.AlterTableStmt{Actions: &ast.List{Items: []ast.Node{&ast.AlterTableCmd{Action: ast.AT_ADD_COLUMN, ColumnDef: col}}}}
	want := []v1review.Finding{{Code: v1review.ColumnNotNullNoDefault, Name: "A", Start: 3, End: 5}}
	if diff := cmp.Diff(want, ColumnAddNotNullRequireDefault([]ast.Node{stmt})); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
