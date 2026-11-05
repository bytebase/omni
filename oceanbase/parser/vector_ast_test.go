package parser

import (
	"testing"

	ast "github.com/bytebase/omni/oceanbase/ast"
)

func oneNode(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

func TestVectorColumnTypeAST(t *testing.T) {
	sql := "CREATE TABLE vector (vector VECTOR(3) NOT NULL, b INT)"
	stmt := oneNode(t, sql).(*ast.CreateTableStmt)
	dt := stmt.Columns[0].TypeName
	if stmt.Columns[0].Name != "vector" || dt.Name != "VECTOR" || dt.Length != 3 {
		t.Errorf("column: %+v, type %+v", stmt.Columns[0], dt)
	}
	if got := sql[dt.Loc.Start:dt.Loc.End]; got != "VECTOR(3)" {
		t.Errorf("Loc spans %q", got)
	}
	if len(stmt.Columns) != 2 || stmt.Columns[1].TypeName.Name != "INT" {
		t.Errorf("the column after the vector type was lost: %+v", stmt.Columns)
	}
}

// TestCreateVectorIndexAST: the statement says it is a vector index and WITH keeps its
// list, in order, with each value's kind.
func TestCreateVectorIndexAST(t *testing.T) {
	sql := "CREATE VECTOR INDEX IF NOT EXISTS idx USING BTREE ON d.t (v) WITH (Distance=l2, m=16, ef_search=1.5, a=TRUE) COMMENT 'c'"
	stmt := oneNode(t, sql).(*ast.CreateIndexStmt)
	if !stmt.Vector || stmt.Unique || stmt.Fulltext || stmt.Spatial {
		t.Errorf("Vector=%v Unique=%v Fulltext=%v Spatial=%v", stmt.Vector, stmt.Unique, stmt.Fulltext, stmt.Spatial)
	}
	if !stmt.IfNotExists || stmt.IndexName != "idx" || stmt.IndexType != "BTREE" {
		t.Errorf("IfNotExists=%v IndexName=%q IndexType=%q", stmt.IfNotExists, stmt.IndexName, stmt.IndexType)
	}
	if stmt.Table.Schema != "d" || stmt.Table.Name != "t" || len(stmt.Columns) != 1 {
		t.Errorf("Table=%s.%s Columns=%+v", stmt.Table.Schema, stmt.Table.Name, stmt.Columns)
	}
	if len(stmt.Options) != 2 || stmt.Options[0].Name != "WITH" || stmt.Options[1].Name != "COMMENT" {
		t.Fatalf("Options=%+v", stmt.Options)
	}
	props := stmt.Options[0].Properties
	if len(props) != 4 {
		t.Fatalf("Properties=%+v", props)
	}
	if v, ok := props[0].Value.(*ast.StringLit); props[0].Name != "distance" || !ok || v.Value != "l2" {
		t.Errorf("distance: %+v", props[0])
	}
	if v, ok := props[1].Value.(*ast.IntLit); props[1].Name != "m" || !ok || v.Value != 16 {
		t.Errorf("m: %+v", props[1])
	}
	if v, ok := props[2].Value.(*ast.StringLit); props[2].Name != "ef_search" || !ok || v.Value != "1.5" {
		t.Errorf("ef_search: %+v", props[2])
	}
	if v, ok := props[3].Value.(*ast.StringLit); props[3].Name != "a" || !ok || v.Value != "TRUE" {
		t.Errorf("a: %+v", props[3])
	}
	if got := sql[stmt.Options[0].Loc.Start:stmt.Options[0].Loc.End]; got != "WITH (Distance=l2, m=16, ef_search=1.5, a=TRUE)" {
		t.Errorf("WITH Loc spans %q", got)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(sql) {
		t.Errorf("Loc=%+v, want the whole statement of %d bytes", stmt.Loc, len(sql))
	}
}

func TestVectorKeyAST(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX k (v) WITH (distance=l2))",
		"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY k (v) WITH (distance=l2))",
	} {
		stmt := oneNode(t, sql).(*ast.CreateTableStmt)
		if len(stmt.Columns) != 1 || len(stmt.Constraints) != 1 {
			t.Fatalf("%q: Columns=%+v Constraints=%+v", sql, stmt.Columns, stmt.Constraints)
		}
		c := stmt.Constraints[0]
		if c.Type != ast.ConstrVector || c.Name != "k" || len(c.Columns) != 1 || c.Columns[0] != "v" {
			t.Errorf("%q: %+v", sql, c)
		}
		if len(c.IndexOptions) != 1 || c.IndexOptions[0].Name != "WITH" || len(c.IndexOptions[0].Properties) != 1 {
			t.Errorf("%q: IndexOptions=%+v", sql, c.IndexOptions)
		}
	}
	// A bare ADD reads vector as a column's name, which it may be here.
	alter := oneNode(t, "ALTER TABLE t ADD vector VECTOR(3)").(*ast.AlterTableStmt)
	if cmd := alter.Commands[0]; cmd.Type != ast.ATAddColumn {
		t.Errorf("command: %+v", cmd)
	}
}
