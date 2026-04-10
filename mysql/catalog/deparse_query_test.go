package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// parseSelectForDeparse parses a single SELECT and returns the AST node.
func parseSelectForDeparse(t *testing.T, sql string) *nodes.SelectStmt {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", list.Len())
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt, got %T", list.Items[0])
	}
	return sel
}

// setupDeparseTestCatalog creates a catalog with test tables.
func setupDeparseTestCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE testdb; USE testdb;", nil)
	if err != nil {
		t.Fatalf("setup parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("setup exec error: %v", r.Error)
		}
	}

	ddl := `
		CREATE TABLE employees (
			id INT PRIMARY KEY,
			name VARCHAR(100),
			salary DECIMAL(10,2),
			department_id INT
		);
		CREATE TABLE departments (
			id INT PRIMARY KEY,
			name VARCHAR(100)
		);
	`
	results, err = c.Exec(ddl, nil)
	if err != nil {
		t.Fatalf("DDL parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("DDL exec error: %v", r.Error)
		}
	}
	return c
}

// roundTrip parses, analyzes, deparses, re-parses, and re-analyzes.
// Returns both Query IRs and the intermediate SQL.
func roundTrip(t *testing.T, c *Catalog, sql string) (q1 *Query, sql2 string, q2 *Query) {
	t.Helper()

	sel1 := parseSelectForDeparse(t, sql)
	var err error
	q1, err = c.AnalyzeSelectStmt(sel1)
	if err != nil {
		t.Fatalf("analyze q1 error: %v", err)
	}

	sql2 = DeparseQuery(q1)
	t.Logf("Original:  %s", sql)
	t.Logf("Deparsed:  %s", sql2)

	sel2 := parseSelectForDeparse(t, sql2)
	q2, err = c.AnalyzeSelectStmt(sel2)
	if err != nil {
		t.Fatalf("analyze q2 error (deparsed SQL: %s): %v", sql2, err)
	}

	return q1, sql2, q2
}

// compareQueries checks structural equivalence between two Query IRs.
func compareQueries(t *testing.T, q1, q2 *Query) {
	t.Helper()

	// Count non-junk targets.
	nonJunk := func(q *Query) []*TargetEntryQ {
		var out []*TargetEntryQ
		for _, te := range q.TargetList {
			if !te.ResJunk {
				out = append(out, te)
			}
		}
		return out
	}

	tl1 := nonJunk(q1)
	tl2 := nonJunk(q2)
	if len(tl1) != len(tl2) {
		t.Errorf("non-junk target count: q1=%d, q2=%d", len(tl1), len(tl2))
		return
	}

	for i := range tl1 {
		if !strings.EqualFold(tl1[i].ResName, tl2[i].ResName) {
			t.Errorf("target[%d] ResName: q1=%q, q2=%q", i, tl1[i].ResName, tl2[i].ResName)
		}
		// Compare VarExprQ coordinates when both are VarExprQ.
		v1, ok1 := tl1[i].Expr.(*VarExprQ)
		v2, ok2 := tl2[i].Expr.(*VarExprQ)
		if ok1 && ok2 {
			if v1.RangeIdx != v2.RangeIdx {
				t.Errorf("target[%d] VarExprQ.RangeIdx: q1=%d, q2=%d", i, v1.RangeIdx, v2.RangeIdx)
			}
			if v1.AttNum != v2.AttNum {
				t.Errorf("target[%d] VarExprQ.AttNum: q1=%d, q2=%d", i, v1.AttNum, v2.AttNum)
			}
		}
	}

	// Compare RTE count.
	if len(q1.RangeTable) != len(q2.RangeTable) {
		t.Errorf("RangeTable count: q1=%d, q2=%d", len(q1.RangeTable), len(q2.RangeTable))
	}

	// Compare WHERE existence.
	hasWhere1 := q1.JoinTree != nil && q1.JoinTree.Quals != nil
	hasWhere2 := q2.JoinTree != nil && q2.JoinTree.Quals != nil
	if hasWhere1 != hasWhere2 {
		t.Errorf("WHERE presence: q1=%v, q2=%v", hasWhere1, hasWhere2)
	}

	// Compare GROUP BY count.
	if len(q1.GroupClause) != len(q2.GroupClause) {
		t.Errorf("GroupClause count: q1=%d, q2=%d", len(q1.GroupClause), len(q2.GroupClause))
	}

	// Compare HAVING existence.
	if (q1.HavingQual != nil) != (q2.HavingQual != nil) {
		t.Errorf("HAVING presence: q1=%v, q2=%v", q1.HavingQual != nil, q2.HavingQual != nil)
	}

	// Compare ORDER BY count.
	if len(q1.SortClause) != len(q2.SortClause) {
		t.Errorf("SortClause count: q1=%d, q2=%d", len(q1.SortClause), len(q2.SortClause))
	}

	// Compare set op.
	if q1.SetOp != q2.SetOp {
		t.Errorf("SetOp: q1=%d, q2=%d", q1.SetOp, q2.SetOp)
	}
}

// TestDeparseQuery_16_1_SimpleRoundTrip tests a simple SELECT with WHERE.
func TestDeparseQuery_16_1_SimpleRoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, _, q2 := roundTrip(t, c, "SELECT name, salary FROM employees WHERE salary > 50000")
	compareQueries(t, q1, q2)
}

// TestDeparseQuery_16_2_JoinRoundTrip tests a JOIN round-trip.
func TestDeparseQuery_16_2_JoinRoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, _, q2 := roundTrip(t, c, "SELECT e.name, d.name FROM employees e JOIN departments d ON e.department_id = d.id")
	compareQueries(t, q1, q2)
}

// TestDeparseQuery_16_3_AggregateRoundTrip tests GROUP BY + HAVING round-trip.
func TestDeparseQuery_16_3_AggregateRoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, _, q2 := roundTrip(t, c, "SELECT department_id, COUNT(*) FROM employees GROUP BY department_id HAVING COUNT(*) > 5")
	compareQueries(t, q1, q2)
}

// TestDeparseQuery_16_4_CTERoundTrip tests CTE round-trip.
func TestDeparseQuery_16_4_CTERoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, _, q2 := roundTrip(t, c, "WITH cte AS (SELECT id, name FROM employees) SELECT id, name FROM cte")
	compareQueries(t, q1, q2)
}

// TestDeparseQuery_16_5_SetOpRoundTrip tests UNION ALL round-trip.
func TestDeparseQuery_16_5_SetOpRoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, _, q2 := roundTrip(t, c, "SELECT name FROM employees UNION ALL SELECT name FROM departments")
	compareQueries(t, q1, q2)
}

// TestDeparseQuery_BareLiteral tests SELECT 1 (no FROM).
func TestDeparseQuery_BareLiteral(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, sql2, _ := roundTrip(t, c, "SELECT 1")
	if q1 == nil {
		t.Fatal("q1 is nil")
	}
	// Should not contain "from".
	if strings.Contains(strings.ToLower(sql2), "from") {
		t.Errorf("deparsed bare literal should not have FROM: %s", sql2)
	}
}

// TestDeparseQuery_DeparseOutput tests that DeparseQuery produces valid SQL text.
func TestDeparseQuery_DeparseOutput(t *testing.T) {
	c := setupDeparseTestCatalog(t)

	tests := []struct {
		name string
		sql  string
	}{
		{"simple_select", "SELECT id, name FROM employees"},
		{"where_clause", "SELECT name FROM employees WHERE salary > 50000"},
		{"order_by", "SELECT name, salary FROM employees ORDER BY salary DESC"},
		{"limit", "SELECT name FROM employees LIMIT 10"},
		{"distinct", "SELECT DISTINCT department_id FROM employees"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := parseSelectForDeparse(t, tt.sql)
			q, err := c.AnalyzeSelectStmt(sel)
			if err != nil {
				t.Fatalf("analyze error: %v", err)
			}
			result := DeparseQuery(q)
			if result == "" {
				t.Fatal("DeparseQuery returned empty string")
			}
			t.Logf("SQL: %s => %s", tt.sql, result)

			// Verify the deparsed SQL can be parsed back.
			_, err = parser.Parse(result)
			if err != nil {
				t.Fatalf("deparsed SQL failed to parse: %v\nSQL: %s", err, result)
			}
		})
	}
}
