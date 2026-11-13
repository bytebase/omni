package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// placed is a column rule finding as the tests compare it, with Text the
// source its range anchors: empty when it has no range.
type placed struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
	Text      string
}

func place(sql string, findings []v1review.Finding) []placed {
	var out []placed
	for _, f := range findings {
		out = append(out, placed{f.Code, f.Name, f.Table, f.Value, f.Statement, sql[f.Start:f.End]})
	}
	return out
}

// checkPlaced parses sql, runs rule over it, and compares the findings.
func checkPlaced(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding, want []placed) {
	t.Helper()
	if diff := cmp.Diff(want, place(sql, rule(parse(t, sql)))); diff != "" {
		t.Errorf("%s\n(-want +got):\n%s", sql, diff)
	}
}

func TestWalkColumns(t *testing.T) {
	const sql = "CREATE TABLE t (a INT, b INT);\n" +
		"ALTER TABLE u ADD COLUMN c INT, ADD COLUMN (d INT, e INT), CHANGE f g INT, MODIFY h INT, RENAME COLUMN i TO j, DROP COLUMN k, ALTER COLUMN l SET DEFAULT 1;\n" +
		"CREATE TABLE v LIKE t;"
	type visit struct{ Table, Column, Text string }
	var got []visit
	walkColumns(parse(t, sql), func(table string, col *ast.ColumnDef, loc ast.Loc) {
		got = append(got, visit{table, col.Name, sql[loc.Start:loc.End]})
	})
	alter := "ALTER TABLE u ADD COLUMN c INT, ADD COLUMN (d INT, e INT), CHANGE f g INT, MODIFY h INT, RENAME COLUMN i TO j, DROP COLUMN k, ALTER COLUMN l SET DEFAULT 1"
	want := []visit{
		{"t", "a", "a INT"}, {"t", "b", "b INT"},
		{"u", "c", alter}, {"u", "d", alter}, {"u", "e", alter}, {"u", "g", alter}, {"u", "h", alter},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestWalkColumnsNilEntries(t *testing.T) {
	stmts := []ast.Node{
		&ast.CreateTableStmt{Columns: []*ast.ColumnDef{{Name: "a"}}},
		&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{{Type: ast.ATAddColumn, Column: &ast.ColumnDef{Name: "b"}}}},
		&ast.CreateTableStmt{Table: &ast.TableRef{Name: "t"}, Columns: []*ast.ColumnDef{nil, {Name: "c"}}},
		&ast.AlterTableStmt{Table: &ast.TableRef{Name: "u"}, Commands: []*ast.AlterTableCmd{nil, {Type: ast.ATAddColumn, Columns: []*ast.ColumnDef{nil}}, {Type: ast.ATModifyColumn}}},
	}
	var got []string
	walkColumns(stmts, func(table string, col *ast.ColumnDef, _ ast.Loc) {
		got = append(got, table+"."+col.Name)
	})
	if diff := cmp.Diff([]string{"t.c"}, got); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}
