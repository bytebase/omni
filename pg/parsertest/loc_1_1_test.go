package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// ---------------------------------------------------------------------------
// Section 1.1: FROM/Join nodes
// ---------------------------------------------------------------------------

func TestLocJoinExpr(t *testing.T) {
	sql := "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	j := sel.FromClause.Items[0].(*nodes.JoinExpr)
	got := sql[j.Loc.Start:j.Loc.End]
	want := "t1 JOIN t2 ON t1.id = t2.id"
	if got != want {
		t.Errorf("JoinExpr text = %q, want %q", got, want)
	}
}

func TestLocJoinExprLeft(t *testing.T) {
	sql := "SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	j := sel.FromClause.Items[0].(*nodes.JoinExpr)
	got := sql[j.Loc.Start:j.Loc.End]
	want := "t1 LEFT JOIN t2 ON t1.id = t2.id"
	if got != want {
		t.Errorf("JoinExpr LEFT text = %q, want %q", got, want)
	}
}

func TestLocJoinExprCross(t *testing.T) {
	sql := "SELECT * FROM t1 CROSS JOIN t2"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	j := sel.FromClause.Items[0].(*nodes.JoinExpr)
	got := sql[j.Loc.Start:j.Loc.End]
	want := "t1 CROSS JOIN t2"
	if got != want {
		t.Errorf("JoinExpr CROSS text = %q, want %q", got, want)
	}
}

func TestLocJoinExprNatural(t *testing.T) {
	sql := "SELECT * FROM t1 NATURAL JOIN t2"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	j := sel.FromClause.Items[0].(*nodes.JoinExpr)
	got := sql[j.Loc.Start:j.Loc.End]
	want := "t1 NATURAL JOIN t2"
	if got != want {
		t.Errorf("JoinExpr NATURAL text = %q, want %q", got, want)
	}
}

func TestLocRangeSubselect(t *testing.T) {
	sql := "SELECT * FROM (SELECT 1) AS sub"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rs := sel.FromClause.Items[0].(*nodes.RangeSubselect)
	got := sql[rs.Loc.Start:rs.Loc.End]
	want := "(SELECT 1) AS sub"
	if got != want {
		t.Errorf("RangeSubselect text = %q, want %q", got, want)
	}
}

func TestLocRangeSubselectLateral(t *testing.T) {
	sql := "SELECT * FROM t1, LATERAL (SELECT * FROM t2) AS sub"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rs := sel.FromClause.Items[1].(*nodes.RangeSubselect)
	got := sql[rs.Loc.Start:rs.Loc.End]
	want := "LATERAL (SELECT * FROM t2) AS sub"
	if got != want {
		t.Errorf("RangeSubselect LATERAL text = %q, want %q", got, want)
	}
}

func TestLocRangeFunction(t *testing.T) {
	sql := "SELECT * FROM generate_series(1, 10)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rf := sel.FromClause.Items[0].(*nodes.RangeFunction)
	got := sql[rf.Loc.Start:rf.Loc.End]
	want := "generate_series(1, 10)"
	if got != want {
		t.Errorf("RangeFunction text = %q, want %q", got, want)
	}
}

func TestLocRangeFunctionWithOrdinality(t *testing.T) {
	sql := "SELECT * FROM generate_series(1, 10) WITH ORDINALITY"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rf := sel.FromClause.Items[0].(*nodes.RangeFunction)
	got := sql[rf.Loc.Start:rf.Loc.End]
	want := "generate_series(1, 10) WITH ORDINALITY"
	if got != want {
		t.Errorf("RangeFunction WITH ORDINALITY text = %q, want %q", got, want)
	}
}

func TestLocRangeFunctionRowsFrom(t *testing.T) {
	sql := "SELECT * FROM ROWS FROM(generate_series(1,3), generate_series(1,4))"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rf := sel.FromClause.Items[0].(*nodes.RangeFunction)
	got := sql[rf.Loc.Start:rf.Loc.End]
	want := "ROWS FROM(generate_series(1,3), generate_series(1,4))"
	if got != want {
		t.Errorf("RangeFunction ROWS FROM text = %q, want %q", got, want)
	}
}

func TestLocCurrentOfExpr(t *testing.T) {
	sql := "DELETE FROM t1 WHERE CURRENT OF cursor1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	del := raw.Stmt.(*nodes.DeleteStmt)
	co := del.WhereClause.(*nodes.CurrentOfExpr)
	got := sql[co.Loc.Start:co.Loc.End]
	want := "CURRENT OF cursor1"
	if got != want {
		t.Errorf("CurrentOfExpr text = %q, want %q", got, want)
	}
}

func TestLocLockingClause(t *testing.T) {
	sql := "SELECT * FROM t1 FOR UPDATE"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if sel.LockingClause == nil || len(sel.LockingClause.Items) == 0 {
		t.Fatal("expected LockingClause")
	}
	lc := sel.LockingClause.Items[0].(*nodes.LockingClause)
	got := sql[lc.Loc.Start:lc.Loc.End]
	want := "FOR UPDATE"
	if got != want {
		t.Errorf("LockingClause text = %q, want %q", got, want)
	}
}

func TestLocIntoClause(t *testing.T) {
	sql := "SELECT * INTO newtable FROM t1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if sel.IntoClause == nil {
		t.Fatal("expected IntoClause")
	}
	got := sql[sel.IntoClause.Loc.Start:sel.IntoClause.Loc.End]
	want := "INTO newtable"
	if got != want {
		t.Errorf("IntoClause text = %q, want %q", got, want)
	}
}
