package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// AT / BEFORE (time travel)
// ---------------------------------------------------------------------------

func TestTimeTravel_AtTimestamp(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM tt1 AT(TIMESTAMP => '2024-06-05 15:29:00'::TIMESTAMP)")
	ref := firstTableRef(t, sel)
	tt := ref.TimeTravel
	if tt == nil {
		t.Fatal("expected TimeTravel clause")
	}
	if tt.Kind != ast.TimeTravelAt {
		t.Errorf("kind = %v, want AT", tt.Kind)
	}
	if tt.Anchor != ast.TimeTravelTimestamp {
		t.Errorf("anchor = %v, want TIMESTAMP", tt.Anchor)
	}
	if tt.Expr == nil {
		t.Error("expected anchor expr")
	}
}

func TestTimeTravel_AtLowercase(t *testing.T) {
	// corpus example_04: lowercase `at`.
	sel := firstSelect(t, "SELECT * FROM tt1 at(TIMESTAMP => '2024-06-05 15:29:00'::TIMESTAMP_LTZ)")
	if firstTableRef(t, sel).TimeTravel == nil {
		t.Fatal("expected TimeTravel for lowercase at()")
	}
}

func TestTimeTravel_AtOffsetWithAlias(t *testing.T) {
	// corpus example_09: AT(OFFSET => -60*5) AS T WHERE T.flag = 'valid'.
	sel := firstSelect(t, "SELECT * FROM my_table AT(OFFSET => -60*5) AS T WHERE T.flag = 'valid'")
	ref := firstTableRef(t, sel)
	if ref.TimeTravel == nil || ref.TimeTravel.Anchor != ast.TimeTravelOffset {
		t.Fatalf("expected AT OFFSET, got %+v", ref.TimeTravel)
	}
	if ref.Alias.Name != "T" {
		t.Errorf("alias = %q, want T", ref.Alias.Name)
	}
	if sel.Where == nil {
		t.Error("expected WHERE clause")
	}
}

func TestTimeTravel_AtStatementFunc(t *testing.T) {
	// corpus example_08: AT(TIMESTAMP => TO_TIMESTAMP(...)).
	sel := firstSelect(t, "SELECT * FROM my_table AT(TIMESTAMP => TO_TIMESTAMP(1432669154242, 3))")
	ref := firstTableRef(t, sel)
	if ref.TimeTravel == nil {
		t.Fatal("expected TimeTravel")
	}
	if _, ok := ref.TimeTravel.Expr.(*ast.FuncCallExpr); !ok {
		t.Errorf("anchor expr = %T, want FuncCallExpr", ref.TimeTravel.Expr)
	}
}

func TestTimeTravel_BeforeStatement(t *testing.T) {
	// corpus example_10: BEFORE(STATEMENT => '...').
	sel := firstSelect(t, "SELECT * FROM my_table BEFORE(STATEMENT => '8e5d0ca9-005e-44e6-b858-a8f5b37c5726')")
	ref := firstTableRef(t, sel)
	if ref.TimeTravel == nil {
		t.Fatal("expected TimeTravel")
	}
	if ref.TimeTravel.Kind != ast.TimeTravelBefore {
		t.Errorf("kind = %v, want BEFORE", ref.TimeTravel.Kind)
	}
	if ref.TimeTravel.Anchor != ast.TimeTravelStatement {
		t.Errorf("anchor = %v, want STATEMENT", ref.TimeTravel.Anchor)
	}
}

func TestTimeTravel_JoinBothSides(t *testing.T) {
	// corpus example_11: BEFORE(...) AS oldt FULL OUTER JOIN ... AT(...) AS newt.
	in := `SELECT oldt.*, newt.*
  FROM my_table BEFORE(STATEMENT => '8e5d0ca9-005e-44e6-b858-a8f5b37c5726') AS oldt
    FULL OUTER JOIN my_table AT(STATEMENT => '8e5d0ca9-005e-44e6-b858-a8f5b37c5726') AS newt
    ON oldt.id = newt.id
  WHERE oldt.id IS NULL OR newt.id IS NULL`
	sel := firstSelect(t, in)
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("FROM[0] = %T, want JoinExpr", sel.From[0])
	}
	left, ok := join.Left.(*ast.TableRef)
	if !ok || left.TimeTravel == nil || left.TimeTravel.Kind != ast.TimeTravelBefore {
		t.Fatalf("join.Left = %+v, want BEFORE time travel", join.Left)
	}
	if left.Alias.Name != "oldt" {
		t.Errorf("left alias = %q, want oldt", left.Alias.Name)
	}
	right, ok := join.Right.(*ast.TableRef)
	if !ok || right.TimeTravel == nil || right.TimeTravel.Kind != ast.TimeTravelAt {
		t.Fatalf("join.Right = %+v, want AT time travel", join.Right)
	}
}

func TestTimeTravel_QualifiedTableJoin(t *testing.T) {
	// corpus example_12: db1.public.htt1 AT(...) h JOIN db1.public.tt1 AT(...) t ON ...
	in := `SELECT * FROM db1.public.htt1
    AT(TIMESTAMP => '2024-06-05 17:50:00'::TIMESTAMP_LTZ) h
    JOIN db1.public.tt1
    AT(TIMESTAMP => '2024-06-05 17:50:00'::TIMESTAMP_LTZ) t
    ON h.c1=t.c1`
	sel := firstSelect(t, in)
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("FROM[0] = %T, want JoinExpr", sel.From[0])
	}
	left := join.Left.(*ast.TableRef)
	if left.TimeTravel == nil {
		t.Error("left table should carry AT")
	}
	if left.Alias.Name != "h" {
		t.Errorf("left alias = %q, want h", left.Alias.Name)
	}
}

func TestTimeTravel_OnCTE(t *testing.T) {
	// corpus example_13: WITH mycte AS (...) SELECT * FROM mycte AT(...).
	in := `WITH mycte AS (SELECT mytable.* FROM mytable) SELECT * FROM mycte AT(TIMESTAMP => '2024-03-13 13:56:09.553 +0100'::TIMESTAMP_TZ)`
	sel := firstSelect(t, in)
	if firstTableRef(t, sel).TimeTravel == nil {
		t.Fatal("expected TimeTravel on CTE reference")
	}
}

// ---------------------------------------------------------------------------
// CHANGES
// ---------------------------------------------------------------------------

func TestChanges_DefaultAt(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM t CHANGES(INFORMATION => DEFAULT) AT(TIMESTAMP => '2024-06-05 15:29:00'::TIMESTAMP)")
	ref := firstTableRef(t, sel)
	ch := ref.Changes
	if ch == nil {
		t.Fatal("expected Changes clause")
	}
	if ch.Info != ast.ChangesDefault {
		t.Errorf("info = %v, want DEFAULT", ch.Info)
	}
	if ch.Start == nil || ch.Start.Kind != ast.TimeTravelAt {
		t.Errorf("changes anchor = %+v, want AT", ch.Start)
	}
}

func TestChanges_AppendOnlyBefore(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM t CHANGES(INFORMATION => APPEND_ONLY) BEFORE(STATEMENT => 'abc')")
	ch := firstTableRef(t, sel).Changes
	if ch.Info != ast.ChangesAppendOnly {
		t.Errorf("info = %v, want APPEND_ONLY", ch.Info)
	}
	if ch.Start.Kind != ast.TimeTravelBefore {
		t.Errorf("anchor kind = %v, want BEFORE", ch.Start.Kind)
	}
}

func TestChanges_WithEnd(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM t CHANGES(INFORMATION => DEFAULT) AT(OFFSET => -60) END(OFFSET => -10)")
	ch := firstTableRef(t, sel).Changes
	if ch.End == nil {
		t.Fatal("expected END bound")
	}
	if ch.End.Anchor != ast.TimeTravelOffset {
		t.Errorf("end anchor = %v, want OFFSET", ch.End.Anchor)
	}
}

// ---------------------------------------------------------------------------
// CONNECT BY / START WITH
// ---------------------------------------------------------------------------

func TestConnectBy_StartWith(t *testing.T) {
	// corpus example_03.
	in := `SELECT employee_ID, manager_ID, title
  FROM employees
    START WITH title = 'President'
    CONNECT BY manager_ID = PRIOR employee_id
  ORDER BY employee_ID`
	sel := firstSelect(t, in)
	if sel.StartWith == nil {
		t.Fatal("expected START WITH condition")
	}
	if len(sel.ConnectBy) != 1 {
		t.Fatalf("connect by = %d, want 1", len(sel.ConnectBy))
	}
	// The CONNECT BY condition is `manager_ID = PRIOR employee_id`; the RHS of
	// the comparison must be a UnaryExpr(PRIOR, employee_id).
	bin, ok := sel.ConnectBy[0].(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("connect by[0] = %T, want BinaryExpr", sel.ConnectBy[0])
	}
	un, ok := bin.Right.(*ast.UnaryExpr)
	if !ok || un.Op != ast.UnaryPrior {
		t.Fatalf("RHS = %T/%v, want UnaryExpr(PRIOR)", bin.Right, un)
	}
	// ORDER BY still applies to the SELECT.
	if len(sel.OrderBy) != 1 {
		t.Errorf("order by = %d, want 1", len(sel.OrderBy))
	}
}

func TestConnectBy_PriorOnLeft(t *testing.T) {
	in := `SELECT * FROM components START WITH component_ID = 1 CONNECT BY parent_component_ID = PRIOR component_ID ORDER BY 1`
	sel := firstSelect(t, in)
	if len(sel.ConnectBy) != 1 {
		t.Fatalf("connect by = %d, want 1", len(sel.ConnectBy))
	}
}

func TestConnectBy_NoStartWith(t *testing.T) {
	in := `SELECT * FROM t CONNECT BY id = PRIOR parent_id`
	sel := firstSelect(t, in)
	if sel.StartWith != nil {
		t.Error("START WITH should be nil")
	}
	if len(sel.ConnectBy) != 1 {
		t.Fatalf("connect by = %d, want 1", len(sel.ConnectBy))
	}
}

func TestConnectBy_SysConnectByPath(t *testing.T) {
	// corpus example_04: SYS_CONNECT_BY_PATH(...) in the SELECT list.
	in := `SELECT SYS_CONNECT_BY_PATH(title, ' -> '), employee_ID
  FROM employees
    START WITH title = 'President'
    CONNECT BY manager_ID = PRIOR employee_id
  ORDER BY employee_ID`
	sel := firstSelect(t, in)
	if len(sel.ConnectBy) != 1 {
		t.Fatalf("connect by = %d, want 1", len(sel.ConnectBy))
	}
	fc, ok := sel.Targets[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("target[0] = %T, want FuncCallExpr", sel.Targets[0].Expr)
	}
	_ = fc
}

func TestConnectBy_ConnectByRoot(t *testing.T) {
	// corpus example_05: CONNECT_BY_ROOT title AS root_title in the SELECT list.
	in := `SELECT employee_ID, CONNECT_BY_ROOT title AS root_title
  FROM employees
    START WITH title = 'President'
    CONNECT BY manager_ID = PRIOR employee_id
  ORDER BY employee_ID`
	sel := firstSelect(t, in)
	// The 2nd target is CONNECT_BY_ROOT title.
	un, ok := sel.Targets[1].Expr.(*ast.UnaryExpr)
	if !ok || un.Op != ast.UnaryConnectByRoot {
		t.Fatalf("target[1] = %T/%v, want UnaryExpr(CONNECT_BY_ROOT)", sel.Targets[1].Expr, un)
	}
	if sel.Targets[1].Alias.Name != "root_title" {
		t.Errorf("target[1] alias = %q, want root_title", sel.Targets[1].Alias.Name)
	}
}

// PRIOR must NOT be treated as a prefix operator outside CONNECT BY — a column
// literally named "prior" must still parse as a column reference.
func TestConnectBy_PriorAsColumnOutsideConnectBy(t *testing.T) {
	sel := firstSelect(t, "SELECT prior FROM t")
	cr, ok := sel.Targets[0].Expr.(*ast.ColumnRef)
	if !ok || len(cr.Parts) != 1 || cr.Parts[0].Name != "prior" {
		t.Fatalf("target[0] = %T (%+v), want ColumnRef{prior}", sel.Targets[0].Expr, sel.Targets[0].Expr)
	}
}

func TestConnectBy_Negatives(t *testing.T) {
	cases := []string{
		"SELECT * FROM t START WITH a = 1",         // START WITH without CONNECT BY
		"SELECT * FROM t CONNECT a = PRIOR b",      // CONNECT without BY
		"SELECT * FROM t START a = 1 CONNECT BY b", // START without WITH
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) == 0 {
				t.Errorf("expected a parse error for %q, got none", in)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Time-travel negatives
// ---------------------------------------------------------------------------

func TestTimeTravel_Negatives(t *testing.T) {
	cases := []string{
		"SELECT * FROM t AT(FOO => 1)",                                  // unknown anchor
		"SELECT * FROM t AT TIMESTAMP => 1",                             // missing parens
		"SELECT * FROM t AT(TIMESTAMP 1)",                               // missing =>
		"SELECT * FROM t CHANGES(INFORMATION => DEFAULT)",               // CHANGES without AT/BEFORE
		"SELECT * FROM t CHANGES(FOO => DEFAULT) AT(OFFSET => 1)",       // wrong key
		"SELECT * FROM t CHANGES(INFORMATION => BOGUS) AT(OFFSET => 1)", // bad info value
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) == 0 {
				t.Errorf("expected a parse error for %q, got none", in)
			}
		})
	}
}
