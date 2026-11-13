package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestTableNoForeignKey(t *testing.T) {
	fk := func(table, text string) located {
		return located{Code: v1review.TableForeignKey, Table: table, Text: text}
	}
	const alter = "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u(b), ADD FOREIGN KEY (c) REFERENCES u(d)"
	for _, tt := range []struct {
		name string
		sql  string
		want []located
	}{
		{name: "no foreign key", sql: "CREATE TABLE t (a INT, PRIMARY KEY (a), UNIQUE KEY uk (a));"},
		{
			name: "create anchors each constraint, name without database",
			sql:  "CREATE TABLE db.t (\n  a INT,\n  CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES u(b),\n  FOREIGN KEY (a) REFERENCES u(c)\n);",
			want: []located{
				fk("t", "CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES u(b)"),
				fk("t", "FOREIGN KEY (a) REFERENCES u(c)"),
			},
		},
		{
			name: "alter anchors the statement, once per foreign key",
			sql:  alter + ";",
			want: []located{fk("t", alter), fk("t", alter)},
		},
		{name: "column references not judged", sql: "CREATE TABLE t (a INT REFERENCES u(b));"},
		{name: "drop foreign key not judged", sql: "ALTER TABLE t DROP FOREIGN KEY fk;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := locate(tt.sql, TableNoForeignKey(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableNoForeignKeyNilEntries(t *testing.T) {
	stmts := []ast.Node{
		&ast.CreateTableStmt{Table: &ast.TableRef{Name: "t"}, Constraints: []*ast.Constraint{nil, {Type: ast.ConstrForeignKey}}},
		&ast.CreateTableStmt{Constraints: []*ast.Constraint{{Type: ast.ConstrForeignKey}}},
		&ast.AlterTableStmt{Table: &ast.TableRef{Name: "u"}, Commands: []*ast.AlterTableCmd{nil, {Type: ast.ATAddConstraint}}},
	}
	want := []v1review.Finding{{Code: v1review.TableForeignKey, Table: "t"}}
	if diff := cmp.Diff(want, TableNoForeignKey(stmts)); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
