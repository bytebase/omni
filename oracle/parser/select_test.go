package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseSelectSimple tests a basic SELECT statement.
func TestParseSelectSimple(t *testing.T) {
	result := ParseAndCheck(t, "SELECT 1 FROM dual")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	sel, ok := raw.Stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", raw.Stmt)
	}
	if sel.TargetList.Len() != 1 {
		t.Errorf("expected 1 target, got %d", sel.TargetList.Len())
	}
}

// TestParseSelectColumns tests SELECT with multiple columns.
func TestParseSelectColumns(t *testing.T) {
	result := ParseAndCheck(t, "SELECT a, b, c FROM t")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.TargetList.Len() != 3 {
		t.Errorf("expected 3 targets, got %d", sel.TargetList.Len())
	}
}

// TestParseSelectStar tests SELECT *.
func TestParseSelectStar(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.TargetList.Len() != 1 {
		t.Errorf("expected 1 target, got %d", sel.TargetList.Len())
	}
}

// TestParseSelectAlias tests SELECT with column aliases.
func TestParseSelectAlias(t *testing.T) {
	result := ParseAndCheck(t, "SELECT salary AS sal, name n FROM employees e")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)

	rt := sel.TargetList.Items[0].(*ast.ResTarget)
	if rt.Name != "SAL" {
		t.Errorf("expected alias SAL, got %q", rt.Name)
	}
	rt2 := sel.TargetList.Items[1].(*ast.ResTarget)
	if rt2.Name != "N" {
		t.Errorf("expected alias N, got %q", rt2.Name)
	}
}

// TestParseSelectDistinct tests SELECT DISTINCT.
func TestParseSelectDistinct(t *testing.T) {
	result := ParseAndCheck(t, "SELECT DISTINCT status FROM orders")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if !sel.Distinct {
		t.Error("expected Distinct=true")
	}
}

// TestParseSelectWhere tests SELECT with WHERE clause.
func TestParseSelectWhere(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees WHERE salary > 50000")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.WhereClause == nil {
		t.Error("expected non-nil WhereClause")
	}
}

func TestParseSelectWhereRowValueNotInSubquery(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM t WHERE (a, b, typ) NOT IN (SELECT x, y, typ FROM u)")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	in, ok := sel.WhereClause.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected InExpr, got %T", sel.WhereClause)
	}
	if !in.Not {
		t.Fatal("expected NOT IN")
	}
	tuple, ok := in.Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected row-value tuple, got %T", in.Expr)
	}
	if tuple.Args.Len() != 3 {
		t.Fatalf("expected 3 tuple args, got %d", tuple.Args.Len())
	}
}

// TestParseSelectOrderBy tests SELECT with ORDER BY.
func TestParseSelectOrderBy(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees ORDER BY salary DESC, name ASC")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.OrderBy.Len() != 2 {
		t.Fatalf("expected 2 ORDER BY items, got %d", sel.OrderBy.Len())
	}
	sb := sel.OrderBy.Items[0].(*ast.SortBy)
	if sb.Dir != ast.SORTBY_DESC {
		t.Errorf("expected DESC, got %d", sb.Dir)
	}
}

// TestParseSelectGroupBy tests SELECT with GROUP BY and HAVING.
func TestParseSelectGroupBy(t *testing.T) {
	result := ParseAndCheck(t, "SELECT department_id, COUNT(*) FROM employees GROUP BY department_id HAVING COUNT(*) > 5")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.GroupClause.Len() != 1 {
		t.Errorf("expected 1 GROUP BY item, got %d", sel.GroupClause.Len())
	}
	if sel.HavingClause == nil {
		t.Error("expected non-nil HavingClause")
	}
}

func TestParseSelectHavingBeforeGroupBy(t *testing.T) {
	result := ParseAndCheck(t, "SELECT department_id, COUNT(*) FROM employees HAVING COUNT(*) > 5 GROUP BY department_id")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.HavingClause == nil {
		t.Fatal("expected non-nil HavingClause")
	}
	if sel.GroupClause.Len() != 1 {
		t.Fatalf("expected 1 GROUP BY item, got %d", sel.GroupClause.Len())
	}
}

// TestParseSelectJoin tests SELECT with JOIN.
func TestParseSelectJoin(t *testing.T) {
	result := ParseAndCheck(t, "SELECT e.name, d.name FROM employees e INNER JOIN departments d ON e.dept_id = d.id")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FromClause.Len() != 1 {
		t.Fatalf("expected 1 FROM item, got %d", sel.FromClause.Len())
	}
	jc, ok := sel.FromClause.Items[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("expected JoinClause, got %T", sel.FromClause.Items[0])
	}
	if jc.Type != ast.JOIN_INNER {
		t.Errorf("expected JOIN_INNER, got %d", jc.Type)
	}
}

func TestParseSelectParenthesizedJoin(t *testing.T) {
	result := ParseAndCheck(t, "SELECT o.id FROM (orders o INNER JOIN customers c ON o.customer_id = c.id)")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FromClause.Len() != 1 {
		t.Fatalf("expected 1 FROM item, got %d", sel.FromClause.Len())
	}
	if _, ok := sel.FromClause.Items[0].(*ast.JoinClause); !ok {
		t.Fatalf("expected JoinClause, got %T", sel.FromClause.Items[0])
	}
}

func TestParseJsonTableColumnFormatJson(t *testing.T) {
	result := ParseAndCheck(t, "SELECT jt.payload FROM src s, JSON_TABLE(s.doc, '$' COLUMNS (payload VARCHAR2(4000) FORMAT JSON WITHOUT WRAPPER PATH '$.payload')) jt")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FromClause.Len() != 2 {
		t.Fatalf("expected 2 FROM items, got %d", sel.FromClause.Len())
	}
	jt := sel.FromClause.Items[1].(*ast.JsonTableRef)
	if jt.Alias == nil || jt.Alias.Name != "JT" {
		t.Fatalf("expected JSON_TABLE alias JT, got %#v", jt.Alias)
	}
	if jt.Columns == nil || jt.Columns.Len() != 1 {
		t.Fatalf("expected 1 JSON_TABLE column, got %#v", jt.Columns)
	}
	col := jt.Columns.Items[0].(*ast.JsonTableColumn)
	if col.Name != "PAYLOAD" {
		t.Fatalf("expected PAYLOAD column, got %q", col.Name)
	}
}

func TestParseJsonTableSimplifiedPathColumns(t *testing.T) {
	result := ParseAndCheck(t, "SELECT jt.id FROM src s CROSS JOIN JSON_TABLE(s.doc.rows[*] COLUMNS (id NUMBER PATH id)) jt")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FromClause.Len() != 1 {
		t.Fatalf("expected 1 FROM item, got %d", sel.FromClause.Len())
	}
	join := sel.FromClause.Items[0].(*ast.JoinClause)
	jt := join.Right.(*ast.JsonTableRef)
	if jt.Columns == nil || jt.Columns.Len() != 1 {
		t.Fatalf("expected 1 JSON_TABLE column, got %#v", jt.Columns)
	}
	col := jt.Columns.Items[0].(*ast.JsonTableColumn)
	if col.Name != "ID" {
		t.Fatalf("expected ID column, got %q", col.Name)
	}
}

// TestParseSelectLeftJoin tests SELECT with LEFT JOIN.
func TestParseSelectLeftJoin(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM a LEFT OUTER JOIN b ON a.id = b.id")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	jc := sel.FromClause.Items[0].(*ast.JoinClause)
	if jc.Type != ast.JOIN_LEFT {
		t.Errorf("expected JOIN_LEFT, got %d", jc.Type)
	}
}

// TestParseSelectCrossJoin tests SELECT with CROSS JOIN.
func TestParseSelectCrossJoin(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM a CROSS JOIN b")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	jc := sel.FromClause.Items[0].(*ast.JoinClause)
	if jc.Type != ast.JOIN_CROSS {
		t.Errorf("expected JOIN_CROSS, got %d", jc.Type)
	}
}

// TestParseSelectConnectBy tests hierarchical query.
func TestParseSelectConnectBy(t *testing.T) {
	result := ParseAndCheck(t, "SELECT employee_id, manager_id FROM employees START WITH manager_id IS NULL CONNECT BY PRIOR employee_id = manager_id")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.Hierarchical == nil {
		t.Fatal("expected non-nil Hierarchical")
	}
	if sel.Hierarchical.ConnectBy == nil {
		t.Error("expected non-nil ConnectBy")
	}
	if sel.Hierarchical.StartWith == nil {
		t.Error("expected non-nil StartWith")
	}
}

// TestParseSelectUnion tests UNION.
func TestParseSelectUnion(t *testing.T) {
	result := ParseAndCheck(t, "SELECT 1 FROM dual UNION SELECT 2 FROM dual")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.Op != ast.SETOP_UNION {
		t.Errorf("expected SETOP_UNION, got %d", sel.Op)
	}
}

// TestParseSelectUnionAll tests UNION ALL.
func TestParseSelectUnionAll(t *testing.T) {
	result := ParseAndCheck(t, "SELECT 1 FROM dual UNION ALL SELECT 2 FROM dual")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.Op != ast.SETOP_UNION {
		t.Errorf("expected SETOP_UNION, got %d", sel.Op)
	}
	if !sel.SetAll {
		t.Error("expected SetAll=true")
	}
}

// TestParseSelectMinus tests MINUS.
func TestParseSelectMinus(t *testing.T) {
	result := ParseAndCheck(t, "SELECT 1 FROM dual MINUS SELECT 2 FROM dual")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.Op != ast.SETOP_MINUS {
		t.Errorf("expected SETOP_MINUS, got %d", sel.Op)
	}
}

// TestParseSelectForUpdate tests FOR UPDATE.
func TestParseSelectForUpdate(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees FOR UPDATE")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.ForUpdate == nil {
		t.Error("expected non-nil ForUpdate")
	}
}

// TestParseSelectForUpdateNoWait tests FOR UPDATE NOWAIT.
func TestParseSelectForUpdateNoWait(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees FOR UPDATE NOWAIT")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if !sel.ForUpdate.NoWait {
		t.Error("expected NoWait=true")
	}
}

// TestParseSelectFetchFirst tests FETCH FIRST.
func TestParseSelectFetchFirst(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees FETCH FIRST 10 ROWS ONLY")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FetchFirst == nil {
		t.Error("expected non-nil FetchFirst")
	}
}

// TestParseSelectOffset tests OFFSET FETCH.
func TestParseSelectOffset(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM employees OFFSET 5 ROWS FETCH NEXT 10 ROWS ONLY")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FetchFirst == nil {
		t.Error("expected non-nil FetchFirst")
	}
	if sel.FetchFirst.Offset == nil {
		t.Error("expected non-nil Offset")
	}
}

// TestParseSelectWithCTE tests WITH clause (CTE).
func TestParseSelectWithCTE(t *testing.T) {
	sql := "WITH dept_summary AS (SELECT department_id, COUNT(*) cnt FROM employees GROUP BY department_id) SELECT * FROM dept_summary"
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.WithClause == nil {
		t.Fatal("expected non-nil WithClause")
	}
	if sel.WithClause.CTEs.Len() != 1 {
		t.Errorf("expected 1 CTE, got %d", sel.WithClause.CTEs.Len())
	}
}

// TestParseSelectSubquery tests subquery in FROM clause.
func TestParseSelectSubquery(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM (SELECT id, name FROM employees) sub")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FromClause.Len() != 1 {
		t.Fatalf("expected 1 FROM item, got %d", sel.FromClause.Len())
	}
	_, ok := sel.FromClause.Items[0].(*ast.SubqueryRef)
	if !ok {
		t.Fatalf("expected SubqueryRef, got %T", sel.FromClause.Items[0])
	}
}

// TestParseSelectMultipleFromTables tests comma-separated FROM tables.
func TestParseSelectMultipleFromTables(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM a, b, c")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	if sel.FromClause.Len() != 3 {
		t.Errorf("expected 3 FROM items, got %d", sel.FromClause.Len())
	}
}

// TestParseSelectNullsFirstLast tests ORDER BY with NULLS FIRST/LAST.
func TestParseSelectNullsFirstLast(t *testing.T) {
	result := ParseAndCheck(t, "SELECT * FROM t ORDER BY a NULLS FIRST, b DESC NULLS LAST")
	raw := result.Items[0].(*ast.RawStmt)
	sel := raw.Stmt.(*ast.SelectStmt)
	sb0 := sel.OrderBy.Items[0].(*ast.SortBy)
	if sb0.NullOrder != ast.SORTBY_NULLS_FIRST {
		t.Errorf("expected NULLS FIRST, got %d", sb0.NullOrder)
	}
	sb1 := sel.OrderBy.Items[1].(*ast.SortBy)
	if sb1.NullOrder != ast.SORTBY_NULLS_LAST {
		t.Errorf("expected NULLS LAST, got %d", sb1.NullOrder)
	}
}
