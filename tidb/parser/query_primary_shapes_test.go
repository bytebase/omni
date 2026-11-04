package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
)

// The shape a TABLE or VALUES query primary and a locking-read clause parse
// into. The corpus in this package records only whether a statement parses, so
// it cannot see a TABLE that became a statement of its own instead of a query
// primary, a VALUES whose rows were flattened into one list, or a WAIT that
// landed nowhere. These assertions are the other half, and they run without a
// container.

// parseOneSelect parses sql and returns its single statement as a SelectStmt.
func parseOneSelect(t *testing.T, sql string) *nodes.SelectStmt {
	t.Helper()
	list := ParseAndCheck(t, sql)
	if len(list.Items) != 1 {
		t.Fatalf("%s: got %d statements, want 1", sql, len(list.Items))
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("%s: statement type = %T, want *ast.SelectStmt", sql, list.Items[0])
	}
	return sel
}

// TestTableValuesPrimaryIsAQueryPrimary: a statement that begins with TABLE or
// VALUES is the query expression it begins, not a statement of its own, so the
// set operation around it is a real set operation with the primary as an
// operand.
func TestTableValuesPrimaryIsAQueryPrimary(t *testing.T) {
	t.Run("TABLE t1 UNION SELECT c2 FROM t2", func(t *testing.T) {
		sel := parseOneSelect(t, "TABLE t1 UNION SELECT c2 FROM t2")
		if sel.SetOp != nodes.SetOpUnion {
			t.Fatalf("SetOp = %v, want SetOpUnion", sel.SetOp)
		}
		if sel.TableSource != nil {
			t.Errorf("the set-operation node carries a TableSource of its own")
		}
		if sel.Left == nil || sel.Left.TableSource == nil {
			t.Fatalf("Left = %+v, want a SelectStmt over a TABLE primary", sel.Left)
		}
		if got := sel.Left.TableSource.Table.Name; got != "t1" {
			t.Errorf("Left.TableSource.Table.Name = %q, want %q", got, "t1")
		}
		if sel.Right == nil || sel.Right.TableSource != nil || len(sel.Right.From) != 1 {
			t.Errorf("Right = %+v, want the plain SELECT leaf", sel.Right)
		}
	})

	t.Run("TABLE t1 UNION TABLE t2 ORDER BY c1 LIMIT 2", func(t *testing.T) {
		sel := parseOneSelect(t, "TABLE t1 UNION TABLE t2 ORDER BY c1 LIMIT 2")
		if sel.SetOp != nodes.SetOpUnion {
			t.Fatalf("SetOp = %v, want SetOpUnion", sel.SetOp)
		}
		// The clauses after the last operand apply to the whole operation, so
		// they are the set-operation node's, not the operand's.
		if len(sel.OrderBy) != 1 || sel.Limit == nil {
			t.Errorf("the set operation holds order=%d limit=%v, want one ORDER BY and a LIMIT",
				len(sel.OrderBy), sel.Limit)
		}
		if sel.Right == nil || len(sel.Right.OrderBy) != 0 || sel.Right.Limit != nil {
			t.Errorf("the last operand kept clauses the set operation holds")
		}
		for i, operand := range []*nodes.SelectStmt{sel.Left, sel.Right} {
			if operand == nil || operand.TableSource == nil {
				t.Fatalf("operand %d = %+v, want a SelectStmt over a TABLE primary", i, operand)
			}
		}
	})

	t.Run("VALUES ROW(1,2) UNION TABLE t", func(t *testing.T) {
		sel := parseOneSelect(t, "VALUES ROW(1,2) UNION TABLE t")
		if sel.SetOp != nodes.SetOpUnion {
			t.Fatalf("SetOp = %v, want SetOpUnion", sel.SetOp)
		}
		if sel.Left == nil || sel.Left.ValuesSource == nil {
			t.Fatalf("Left = %+v, want a SelectStmt over a VALUES primary", sel.Left)
		}
		if sel.Right == nil || sel.Right.TableSource == nil {
			t.Fatalf("Right = %+v, want a SelectStmt over a TABLE primary", sel.Right)
		}
	})

	t.Run("SELECT * FROM t1 WHERE a IN (TABLE t2)", func(t *testing.T) {
		sel := parseOneSelect(t, "SELECT * FROM t1 WHERE a IN (TABLE t2)")
		in, ok := sel.Where.(*nodes.InExpr)
		if !ok {
			t.Fatalf("Where type = %T, want *ast.InExpr", sel.Where)
		}
		if in.Select == nil || in.Select.TableSource == nil {
			t.Fatalf("IN operand = %+v, want a SelectStmt over a TABLE primary", in.Select)
		}
		if got := in.Select.TableSource.Table.Name; got != "t2" {
			t.Errorf("IN operand table = %q, want %q", got, "t2")
		}
	})
}

// TestValuesPrimaryKeepsItsRows: the rows of a VALUES primary stay one list
// per ROW(...). Flattening them into a single list parses just as cleanly and
// reports the wrong arity to every consumer.
func TestValuesPrimaryKeepsItsRows(t *testing.T) {
	sel := parseOneSelect(t, "VALUES ROW(1,2), ROW(3,4)")
	if sel.ValuesSource == nil {
		t.Fatalf("ValuesSource = nil, want the VALUES primary")
	}
	rows := sel.ValuesSource.Rows
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2: %#v", len(rows), rows)
	}
	for i, row := range rows {
		if len(row) != 2 {
			t.Errorf("row %d holds %d values, want 2", i, len(row))
		}
	}

	// An empty ROW() is one row of no values, not no row.
	empty := parseOneSelect(t, "VALUES ROW()")
	if empty.ValuesSource == nil || len(empty.ValuesSource.Rows) != 1 || len(empty.ValuesSource.Rows[0]) != 0 {
		t.Errorf("VALUES ROW() rows = %#v, want one row of no values", empty.ValuesSource)
	}
}

// TestQueryPrimaryClausesHaveOneOwner: the primary holds only its table or
// rows; the ORDER BY, LIMIT, locking clause and INTO after it are the
// SelectStmt's. The TableStmt and ValuesStmt types carry no clause fields, so
// the compiler enforces half of this; what is asserted here is that the
// SelectStmt really did take them.
func TestQueryPrimaryClausesHaveOneOwner(t *testing.T) {
	sel := parseOneSelect(t, "TABLE t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/f'")
	if sel.TableSource == nil || sel.TableSource.Table.Name != "t" {
		t.Fatalf("TableSource = %+v, want the TABLE t primary", sel.TableSource)
	}
	if len(sel.OrderBy) != 1 {
		t.Errorf("OrderBy = %d items, want 1", len(sel.OrderBy))
	}
	if sel.Limit == nil {
		t.Errorf("Limit = nil, want LIMIT 1")
	}
	if sel.ForUpdate == nil || sel.ForUpdate.Share {
		t.Errorf("ForUpdate = %+v, want FOR UPDATE", sel.ForUpdate)
	}
	if sel.Into == nil || sel.Into.Outfile != "/tmp/f" {
		t.Errorf("Into = %+v, want INTO OUTFILE '/tmp/f'", sel.Into)
	}
	// The primary's Loc covers only "TABLE t", so a caller reporting a
	// position inside it does not point at the clauses.
	if got := "TABLE t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/f'"[sel.TableSource.Loc.Start:sel.TableSource.Loc.End]; got != "TABLE t" {
		t.Errorf("TableSource.Loc covers %q, want %q", got, "TABLE t")
	}
}

// TestLockingClauseShape: the OF table list and the wait modifier both land on
// the ForUpdate, and at most one modifier is ever set.
func TestLockingClauseShape(t *testing.T) {
	cases := []struct {
		sql         string
		share       bool
		lockInShare bool
		tables      []string
		noWait      bool
		wait        string // "" for no WAIT
		skipLocked  bool
	}{
		{sql: "SELECT * FROM t FOR UPDATE"},
		{sql: "SELECT * FROM t FOR SHARE", share: true},
		{sql: "SELECT * FROM t LOCK IN SHARE MODE", share: true, lockInShare: true},
		{sql: "SELECT * FROM t FOR UPDATE NOWAIT", noWait: true},
		{sql: "SELECT * FROM t FOR UPDATE SKIP LOCKED", skipLocked: true},
		{sql: "SELECT * FROM t FOR UPDATE WAIT 5", wait: "5"},
		// WAIT 0 is a WAIT, not the absence of one.
		{sql: "SELECT * FROM t FOR UPDATE WAIT 0", wait: "0"},
		{sql: "SELECT * FROM t FOR UPDATE OF t", tables: []string{"t"}},
		{sql: "SELECT * FROM t FOR UPDATE OF t, t1", tables: []string{"t", "t1"}},
		{sql: "SELECT * FROM t FOR UPDATE OF t WAIT 5", tables: []string{"t"}, wait: "5"},
		{sql: "SELECT * FROM t FOR UPDATE OF t, t1 WAIT 11", tables: []string{"t", "t1"}, wait: "11"},
		{sql: "SELECT * FROM t FOR UPDATE OF t NOWAIT", tables: []string{"t"}, noWait: true},
		{sql: "SELECT * FROM t FOR UPDATE OF t SKIP LOCKED", tables: []string{"t"}, skipLocked: true},
		{sql: "SELECT * FROM t FOR SHARE OF t NOWAIT", share: true, tables: []string{"t"}, noWait: true},
		{sql: "TABLE t FOR UPDATE OF t WAIT 5", tables: []string{"t"}, wait: "5"},
		{sql: "VALUES ROW(1) FOR UPDATE WAIT 5", wait: "5"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			sel := parseOneSelect(t, tc.sql)
			fu := sel.ForUpdate
			if fu == nil {
				t.Fatalf("ForUpdate = nil, want a locking clause")
			}
			if fu.Share != tc.share {
				t.Errorf("Share = %v, want %v", fu.Share, tc.share)
			}
			if fu.LockInShareMode != tc.lockInShare {
				t.Errorf("LockInShareMode = %v, want %v", fu.LockInShareMode, tc.lockInShare)
			}
			var tables []string
			for _, ref := range fu.Tables {
				tables = append(tables, ref.Name)
			}
			if strings.Join(tables, ",") != strings.Join(tc.tables, ",") {
				t.Errorf("Tables = %v, want %v", tables, tc.tables)
			}
			if fu.NoWait != tc.noWait {
				t.Errorf("NoWait = %v, want %v", fu.NoWait, tc.noWait)
			}
			if fu.SkipLocked != tc.skipLocked {
				t.Errorf("SkipLocked = %v, want %v", fu.SkipLocked, tc.skipLocked)
			}
			if tc.wait == "" {
				if fu.HasWait {
					t.Errorf("HasWait = true, want no WAIT (WaitSeconds=%q)", fu.WaitSeconds)
				}
			} else {
				if !fu.HasWait {
					t.Fatalf("HasWait = false, want WAIT %s", tc.wait)
				}
				if fu.WaitSeconds != tc.wait {
					t.Errorf("WaitSeconds = %q, want %q", fu.WaitSeconds, tc.wait)
				}
			}
			// At most one wait modifier, whatever the form.
			set := 0
			for _, on := range []bool{fu.NoWait, fu.HasWait, fu.SkipLocked} {
				if on {
					set++
				}
			}
			if set > 1 {
				t.Errorf("%d wait modifiers set at once: %+v", set, fu)
			}
		})
	}
}

// TestQueryPrimaryDeparse: the primary and the clauses it carries round-trip
// through deparse, which is what a view definition is built from. Before the
// primary existed on SelectStmt, deparse wrote "select " for one of these and
// lost the TABLE or the rows entirely.
func TestQueryPrimaryDeparse(t *testing.T) {
	cases := []struct{ sql, want string }{
		{"TABLE t1", "table `t1`"},
		{"TABLE t1 ORDER BY a LIMIT 2", "table `t1` order by `a` limit 2"},
		{"TABLE t1 FOR UPDATE WAIT 5", "table `t1` for update wait 5"},
		{"TABLE t1 FOR SHARE", "table `t1` for share"},
		{"TABLE t1 LOCK IN SHARE MODE", "table `t1` lock in share mode"},
		{"TABLE t1 FOR UPDATE OF t1 SKIP LOCKED", "table `t1` for update of `t1` skip locked"},
		{"VALUES ROW(1) FOR UPDATE", "values row(1) for update"},
		{"VALUES ROW(1) ORDER BY 1 LIMIT 1 FOR SHARE", "values row(1) order by 1 limit 1 for share"},
		{"VALUES ROW(1,2), ROW(3,4)", "values row(1,2),row(3,4)"},
		{"VALUES ROW(1) LIMIT 1", "values row(1) limit 1"},
		{"TABLE t1 UNION TABLE t2", "table `t1` union table `t2`"},
		{"SELECT a FROM t FOR UPDATE OF t WAIT 5", "select `a` AS `a` from `t` for update of `t` wait 5"},
		{"SELECT a FROM t FOR UPDATE SKIP LOCKED", "select `a` AS `a` from `t` for update skip locked"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			sel := parseOneSelect(t, tc.sql)
			if got := deparse.DeparseSelect(sel); got != tc.want {
				t.Errorf("DeparseSelect = %q, want %q", got, tc.want)
			}
		})
	}
}
