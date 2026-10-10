package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

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

// TestParseOffsetColumnName: OFFSET is not reserved in Oracle, and a column
// named OFFSET works wherever an operand can start. Oracle 23ai accepts each
// statement. OFFSET after a complete expression or table reference still
// starts the row_limiting_clause.
func TestParseOffsetColumnName(t *testing.T) {
	tests := []string{
		"SELECT offset, partition FROM t_offset",
		"SELECT offset + 1 AS next_offset FROM t_offset",
		"SELECT t.offset FROM t_offset t ORDER BY offset",
		"SELECT a FROM t_offset WHERE offset = 1",
		"SELECT a FROM t_offset WHERE a = 1 AND offset > 0",
		"UPDATE t_offset SET offset = 1 WHERE offset > 0",
		"DELETE FROM t_offset WHERE offset IS NULL",
		"CREATE INDEX i ON t_offset (a, offset)",
		"CREATE UNIQUE INDEX app.idx1 ON app.hub (offset, partition, txn_date) NOLOGGING  LOCAL",
		"SELECT a FROM t_offset ORDER BY offset OFFSET 1 ROWS",
		"SELECT a FROM t_offset ORDER BY a OFFSET 5 ROWS FETCH NEXT 5 ROWS ONLY",
		"SELECT a FROM t_offset OFFSET 5 ROWS",
		"SELECT a FROM t_offset WHERE a = 1 OFFSET 5 ROWS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

func TestParseOffsetColumnNameAST(t *testing.T) {
	result := ParseAndCheck(t, "SELECT offset FROM t_offset ORDER BY a OFFSET 2 ROWS")
	sel := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt)
	col, ok := sel.TargetList.Items[0].(*ast.ResTarget).Expr.(*ast.ColumnRef)
	if !ok || col.Column != "OFFSET" {
		t.Fatalf("target = %s, want column OFFSET", ast.NodeToString(sel.TargetList))
	}
	if sel.FetchFirst == nil || sel.FetchFirst.Offset == nil {
		t.Fatalf("row limiting clause missing: %s", ast.NodeToString(sel))
	}
}

// TestParseCaseBooleanSelector: CASE opens an expression when WHEN follows it
// or follows an expression after it; the selector may start with NOT, NULLS,
// or LIKEC and may be a Boolean comparison. Otherwise CASE is a column:
// case NOT IN (...), ORDER BY case NULLS FIRST, case LIKEC '...', and a column
// with an alias. Oracle 23ai accepts all of these.
func TestParseCaseBooleanSelector(t *testing.T) {
	result := ParseAndCheck(t, "SELECT CASE NOT TRUE WHEN TRUE THEN 1 ELSE 0 END FROM t")
	target := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt).TargetList.Items[0].(*ast.ResTarget)
	if _, ok := target.Expr.(*ast.CaseExpr); !ok {
		t.Fatalf("target = %T, want *ast.CaseExpr", target.Expr)
	}
	for _, sql := range []string{
		"SELECT CASE nulls WHEN 1 THEN 2 END FROM t",
		"SELECT CASE nulls + 1 WHEN 1 THEN 2 END FROM t",
		"SELECT CASE likec WHEN 'x' THEN 2 END FROM t",
		"SELECT CASE like2 || 'a' WHEN 'xa' THEN 2 END FROM t",
		"SELECT CASE likec = 'x' WHEN TRUE THEN 1 ELSE 0 END FROM t",
		"SELECT CASE likec IN ('x', 'y') WHEN TRUE THEN 1 ELSE 0 END FROM t",
		"SELECT CASE case IS NULL WHEN TRUE THEN 1 ELSE 0 END FROM t",
	} {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			target := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt).TargetList.Items[0].(*ast.ResTarget)
			if _, ok := target.Expr.(*ast.CaseExpr); !ok {
				t.Fatalf("target = %T, want *ast.CaseExpr", target.Expr)
			}
		})
	}
	for _, sql := range []string{
		"SELECT a FROM t WHERE case NOT IN (1, 2)",
		"SELECT a FROM t WHERE case NOT BETWEEN 1 AND 2",
		"SELECT a FROM t WHERE case NOT LIKE 'x%'",
		"SELECT a FROM t ORDER BY case NULLS FIRST",
		"SELECT a FROM t ORDER BY case DESC NULLS LAST",
		"SELECT a FROM t WHERE case LIKEC '1%'",
		"SELECT case a2 FROM t",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseKeywordColumnOuterJoin: a non-reserved keyword followed by the
// legacy outer-join marker (+) is a column, except for the words Oracle still
// reads as their construct (it rejects JSON_VALUE(+) and XMLELEMENT(+)) and
// for pseudo-columns (SYSTIMESTAMP(+) is ORA-30088 even with such a column).
func TestParseKeywordColumnOuterJoin(t *testing.T) {
	for _, word := range []string{"cast", "decode", "case", "interval", "xmlagg"} {
		sql := "SELECT 1 FROM t, u WHERE " + word + "(+) = u.a"
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			where := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt).WhereClause.(*ast.BinaryExpr)
			col, ok := where.Left.(*ast.ColumnRef)
			if !ok || !col.OuterJoin {
				t.Fatalf("left = %s, want an outer-joined column", ast.NodeToString(where.Left))
			}
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
	for _, word := range []string{"json_value", "xmlelement", "treat", "json", "systimestamp"} {
		ParseShouldFail(t, "SELECT 1 FROM t, u WHERE "+word+"(+) = u.a")
	}
	// A qualified reference is a column even for a pseudo-column word.
	ParseAndCheck(t, "SELECT 1 FROM t, u WHERE t.systimestamp(+) = u.a")
}

// TestParseCaseExpressionErrorPosition: deciding between a CASE expression
// and a column named CASE probes the selector without consuming it, so an
// error inside a real CASE expression is still reported where it occurs.
func TestParseCaseExpressionErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		near string
	}{
		{"SELECT CASE x WHEN 1 THEN (1 + ) END FROM t", ")"},
		{"SELECT CASE WHEN a = 1 THEN 1 FROM t", "FROM"},
		{"SELECT CASE likec = 'x' WHEN TRUE THEN (1 + ) END FROM t", ")"},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := Parse(tt.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want a ParseError", err)
			}
			if want := strings.Index(tt.sql, tt.near); pe.Position != want {
				t.Fatalf("error at %d (%q), want %d (%q)", pe.Position, tt.sql[pe.Position:], want, tt.near)
			}
		})
	}
}

// TestParseNestedCaseStaysFast guards parseCaseOrColumn against reparsing:
// a CASE whose selector is another simple CASE, and a column named CASE
// compared with an operand holding another such column, each nested 60
// deep. Parsing either twice per level is exponential (seconds by 20 levels);
// both take milliseconds.
func TestParseNestedCaseStaysFast(t *testing.T) {
	const depth = 60
	selector := "1"
	column := "'1%'"
	for i := 0; i < depth; i++ {
		selector = "CASE " + selector + " WHEN 1 THEN 1 END"
		column = "case LIKEC (" + column + ")"
	}
	for _, sql := range []string{
		"SELECT " + selector + " FROM dual",
		"SELECT a FROM t WHERE " + column,
	} {
		start := time.Now()
		_, _ = Parse(sql)
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("parsing %d nested CASE levels took %s: %.60q...", depth, elapsed, sql)
		}
	}
}
