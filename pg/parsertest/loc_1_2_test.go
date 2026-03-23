package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// ---------------------------------------------------------------------------
// Section 1.2: Expression helper nodes
// ---------------------------------------------------------------------------

func TestLocAIndicesSingle(t *testing.T) {
	sql := "SELECT a[1] FROM t"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	// a[1] is parsed as A_Indirection{Arg: ColumnRef{a}, Indirection: [A_Indices{Uidx: 1}]}
	ind := sel.TargetList.Items[0].(*nodes.ResTarget).Val.(*nodes.A_Indirection)
	ai := ind.Indirection.Items[0].(*nodes.A_Indices)
	got := sql[ai.Loc.Start:ai.Loc.End]
	want := "[1]"
	if got != want {
		t.Errorf("A_Indices single text = %q, want %q", got, want)
	}
}

func TestLocAIndicesSlice(t *testing.T) {
	sql := "SELECT a[1:3] FROM t"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	ind := sel.TargetList.Items[0].(*nodes.ResTarget).Val.(*nodes.A_Indirection)
	ai := ind.Indirection.Items[0].(*nodes.A_Indices)
	got := sql[ai.Loc.Start:ai.Loc.End]
	want := "[1:3]"
	if got != want {
		t.Errorf("A_Indices slice text = %q, want %q", got, want)
	}
}

func TestLocAIndirection(t *testing.T) {
	sql := "SELECT (rec).field FROM t"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	ind := sel.TargetList.Items[0].(*nodes.ResTarget).Val.(*nodes.A_Indirection)
	got := sql[ind.Loc.Start:ind.Loc.End]
	want := ".field"
	if got != want {
		t.Errorf("A_Indirection text = %q, want %q", got, want)
	}
}

func TestLocAStar(t *testing.T) {
	sql := "SELECT * FROM t"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	cr := sel.TargetList.Items[0].(*nodes.ResTarget).Val.(*nodes.ColumnRef)
	star := cr.Fields.Items[0].(*nodes.A_Star)
	got := sql[star.Loc.Start:star.Loc.End]
	want := "*"
	if got != want {
		t.Errorf("A_Star text = %q, want %q", got, want)
	}
}

func TestLocMergeWhenClauseMatched(t *testing.T) {
	sql := "MERGE INTO t USING s ON t.id=s.id WHEN MATCHED THEN UPDATE SET x=1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	merge := raw.Stmt.(*nodes.MergeStmt)
	wc := merge.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	got := sql[wc.Loc.Start:wc.Loc.End]
	want := "WHEN MATCHED THEN UPDATE SET x=1"
	if got != want {
		t.Errorf("MergeWhenClause MATCHED text = %q, want %q", got, want)
	}
}

func TestLocMergeWhenClauseNotMatched(t *testing.T) {
	sql := "MERGE INTO t USING s ON t.id=s.id WHEN NOT MATCHED THEN INSERT VALUES(1)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	merge := raw.Stmt.(*nodes.MergeStmt)
	wc := merge.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	got := sql[wc.Loc.Start:wc.Loc.End]
	want := "WHEN NOT MATCHED THEN INSERT VALUES(1)"
	if got != want {
		t.Errorf("MergeWhenClause NOT MATCHED text = %q, want %q", got, want)
	}
}

func TestLocMultiAssignRef(t *testing.T) {
	sql := "UPDATE t SET (a, b) = (SELECT 1, 2)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	upd := raw.Stmt.(*nodes.UpdateStmt)
	rt := upd.TargetList.Items[0].(*nodes.ResTarget)
	mar := rt.Val.(*nodes.MultiAssignRef)
	got := sql[mar.Loc.Start:mar.Loc.End]
	// The MultiAssignRef's Loc tracks the source expression
	// (SELECT 1, 2) is the source — a SubLink
	want := "(SELECT 1, 2)"
	if got != want {
		t.Errorf("MultiAssignRef text = %q, want %q", got, want)
	}
}

func TestLocTableLikeClause(t *testing.T) {
	sql := "CREATE TABLE t2 (LIKE t1 INCLUDING ALL)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	cs := raw.Stmt.(*nodes.CreateStmt)
	tlc := cs.TableElts.Items[0].(*nodes.TableLikeClause)
	got := sql[tlc.Loc.Start:tlc.Loc.End]
	want := "LIKE t1 INCLUDING ALL"
	if got != want {
		t.Errorf("TableLikeClause text = %q, want %q", got, want)
	}
}
