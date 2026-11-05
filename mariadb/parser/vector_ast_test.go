package parser

import (
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
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

// TestVectorColumnTypeAST: the dimension is the integer part of the first number, which
// is what SHOW CREATE TABLE writes on 11.8.8 for all of these.
func TestVectorColumnTypeAST(t *testing.T) {
	sql := "CREATE TABLE t (a VECTOR(3), b VECTOR(3,2), c VECTOR(1.5), d `vector`(4))"
	stmt := oneNode(t, sql).(*ast.CreateTableStmt)
	want := []struct {
		text   string
		length int
	}{{"VECTOR(3)", 3}, {"VECTOR(3,2)", 3}, {"VECTOR(1.5)", 1}, {"`vector`(4)", 4}}
	for i, w := range want {
		dt := stmt.Columns[i].TypeName
		if dt.Name != "VECTOR" || dt.Length != w.length {
			t.Errorf("column %d: Name=%q Length=%d, want VECTOR and %d", i, dt.Name, dt.Length, w.length)
		}
		if got := sql[dt.Loc.Start:dt.Loc.End]; got != w.text {
			t.Errorf("column %d: Loc spans %q, want %q", i, got, w.text)
		}
	}
}

// TestCreateVectorIndexAST: the statement says it is a vector index and the
// engine-defined options keep their spelling and the kind of their value.
func TestCreateVectorIndexAST(t *testing.T) {
	sql := "CREATE OR REPLACE VECTOR INDEX idx ON d.t (v) WAIT 3 M=6 DISTANCE=cosine `s`='x' d=DEFAULT COMMENT 'c' ALGORITHM=INPLACE LOCK=NONE"
	stmt := oneNode(t, sql).(*ast.CreateIndexStmt)
	if !stmt.Vector || !stmt.OrReplace || stmt.Unique || stmt.Fulltext || stmt.Spatial {
		t.Errorf("Vector=%v OrReplace=%v Unique=%v Fulltext=%v Spatial=%v", stmt.Vector, stmt.OrReplace, stmt.Unique, stmt.Fulltext, stmt.Spatial)
	}
	if stmt.IndexName != "idx" || stmt.Table.Schema != "d" || stmt.Table.Name != "t" {
		t.Errorf("IndexName=%q Table=%s.%s", stmt.IndexName, stmt.Table.Schema, stmt.Table.Name)
	}
	if len(stmt.Columns) != 1 || stmt.Columns[0].Functional {
		t.Fatalf("Columns=%+v", stmt.Columns)
	}
	if ref, ok := stmt.Columns[0].Expr.(*ast.ColumnRef); !ok || ref.Column != "v" {
		t.Errorf("the key part is %+v, want the column v", stmt.Columns[0].Expr)
	}
	if len(stmt.Options) != 5 {
		t.Fatalf("Options=%+v, want 5", stmt.Options)
	}
	m, distance, s, d, comment := stmt.Options[0], stmt.Options[1], stmt.Options[2], stmt.Options[3], stmt.Options[4]
	if v, ok := m.Value.(*ast.IntLit); !m.Engine || m.Name != "M" || !ok || v.Value != 6 {
		t.Errorf("M: %+v", m)
	}
	if v, ok := distance.Value.(*ast.StringLit); !distance.Engine || !distance.ValueIsWord || distance.Name != "DISTANCE" || !ok || v.Value != "cosine" {
		t.Errorf("DISTANCE: %+v", distance)
	}
	if v, ok := s.Value.(*ast.StringLit); !s.Engine || s.ValueIsWord || s.Name != "s" || !ok || v.Value != "x" {
		t.Errorf("s: %+v", s)
	}
	if !d.Engine || d.Value != nil {
		t.Errorf("d: %+v, want DEFAULT as no value", d)
	}
	if comment.Engine || comment.Name != "COMMENT" {
		t.Errorf("COMMENT: %+v", comment)
	}
	if stmt.Algorithm != "INPLACE" || stmt.Lock != "NONE" {
		t.Errorf("Algorithm=%q Lock=%q", stmt.Algorithm, stmt.Lock)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(sql) {
		t.Errorf("Loc=%+v, want the whole statement of %d bytes", stmt.Loc, len(sql))
	}
}

// TestVectorKeyAST: every spelling of the table element and of ALTER's ADD comes out as
// the one constraint, and IF NOT EXISTS is kept where it may stand.
func TestVectorKeyAST(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR k (v) M=6)",
		"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY k (v) M=6)",
		"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX k (v) M=6)",
	} {
		stmt := oneNode(t, sql).(*ast.CreateTableStmt)
		if len(stmt.Constraints) != 1 {
			t.Fatalf("%q: Constraints=%+v", sql, stmt.Constraints)
		}
		c := stmt.Constraints[0]
		if c.Type != ast.ConstrVector || c.Name != "k" || c.IfNotExists {
			t.Errorf("%q: Type=%v Name=%q IfNotExists=%v", sql, c.Type, c.Name, c.IfNotExists)
		}
		if len(c.Columns) != 1 || c.Columns[0] != "v" || len(c.IndexOptions) != 1 {
			t.Errorf("%q: Columns=%v IndexOptions=%+v", sql, c.Columns, c.IndexOptions)
		}
	}
	alter := oneNode(t, "ALTER TABLE t ADD VECTOR IF NOT EXISTS k (v), ADD COLUMN b INT").(*ast.AlterTableStmt)
	if len(alter.Commands) != 2 {
		t.Fatalf("Commands=%+v", alter.Commands)
	}
	c := alter.Commands[0].Constraint
	if alter.Commands[0].Type != ast.ATAddConstraint || c == nil || c.Type != ast.ConstrVector || !c.IfNotExists || c.Name != "k" {
		t.Errorf("first command: %+v, constraint %+v", alter.Commands[0], c)
	}
	unnamed := oneNode(t, "ALTER TABLE t ADD VECTOR (v)").(*ast.AlterTableStmt).Commands[0].Constraint
	if unnamed.Name != "" || len(unnamed.Columns) != 1 {
		t.Errorf("unnamed: %+v", unnamed)
	}
}
