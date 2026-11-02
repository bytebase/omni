package catalog

import "testing"

// ---------------------------------------------------------------------------
// Lineage types and walker — test-only reference implementation.
// ---------------------------------------------------------------------------

// columnLineage represents one source column contributing to a target column.
type columnLineage struct {
	DB     string
	Table  string
	Column string
}

// resultLineage represents the lineage of one output column.
type resultLineage struct {
	Name       string
	SourceCols []columnLineage
}

// collectColumnLineage walks the analyzed Query and extracts column-level
// lineage for each non-junk target entry.
func collectColumnLineage(_ *Catalog, q *Query) []resultLineage {
	var results []resultLineage
	for _, te := range q.TargetList {
		if te.ResJunk {
			continue
		}
		var sources []columnLineage
		collectVarExprs(q, te.Expr, &sources)
		results = append(results, resultLineage{
			Name:       te.ResName,
			SourceCols: sources,
		})
	}
	return results
}

// collectVarExprs recursively walks an expression tree to find all VarExprQ
// nodes and resolves them to (db, table, column) tuples via the RangeTable.
func collectVarExprs(q *Query, expr AnalyzedExpr, out *[]columnLineage) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *VarExprQ:
		resolveVar(q, e, out)
	case *OpExprQ:
		collectVarExprs(q, e.Left, out)
		collectVarExprs(q, e.Right, out)
	case *BoolExprQ:
		for _, arg := range e.Args {
			collectVarExprs(q, arg, out)
		}
	case *FuncCallExprQ:
		for _, arg := range e.Args {
			collectVarExprs(q, arg, out)
		}
	case *CaseExprQ:
		collectVarExprs(q, e.TestExpr, out)
		for _, w := range e.Args {
			collectVarExprs(q, w.Cond, out)
			collectVarExprs(q, w.Then, out)
		}
		collectVarExprs(q, e.Default, out)
	case *CoalesceExprQ:
		for _, arg := range e.Args {
			collectVarExprs(q, arg, out)
		}
	case *CastExprQ:
		collectVarExprs(q, e.Arg, out)
	case *NullTestExprQ:
		collectVarExprs(q, e.Arg, out)
	case *InListExprQ:
		collectVarExprs(q, e.Arg, out)
		for _, item := range e.List {
			collectVarExprs(q, item, out)
		}
	case *BetweenExprQ:
		collectVarExprs(q, e.Arg, out)
		collectVarExprs(q, e.Lower, out)
		collectVarExprs(q, e.Upper, out)
	case *SubLinkExprQ:
		if e.Subquery != nil {
			for _, innerTE := range e.Subquery.TargetList {
				if !innerTE.ResJunk {
					collectVarExprs(e.Subquery, innerTE.Expr, out)
				}
			}
		}
	case *ConstExprQ:
		// No column references in constants.
	case *RowExprQ:
		for _, arg := range e.Args {
			collectVarExprs(q, arg, out)
		}
	}
}

// resolveVar resolves a VarExprQ to (db, table, column) by walking the RangeTable.
func resolveVar(q *Query, v *VarExprQ, out *[]columnLineage) {
	if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	colIdx := v.AttNum - 1 // convert to 0-based
	if colIdx < 0 || colIdx >= len(rte.ColNames) {
		return
	}

	switch rte.Kind {
	case RTERelation:
		if rte.Subquery != nil {
			// MERGE view expanded — recurse into the view's analyzed query.
			if colIdx < len(rte.Subquery.TargetList) {
				collectVarExprs(rte.Subquery, rte.Subquery.TargetList[colIdx].Expr, out)
			}
		} else {
			// Physical table or opaque view — terminal.
			*out = append(*out, columnLineage{
				DB:     rte.DBName,
				Table:  rte.TableName,
				Column: rte.ColNames[colIdx],
			})
		}
	case RTESubquery:
		if rte.Subquery != nil && colIdx < len(rte.Subquery.TargetList) {
			collectVarExprs(rte.Subquery, rte.Subquery.TargetList[colIdx].Expr, out)
		}
	case RTECTE:
		if rte.Subquery != nil && colIdx < len(rte.Subquery.TargetList) {
			collectVarExprs(rte.Subquery, rte.Subquery.TargetList[colIdx].Expr, out)
		}
	case RTEJoin:
		*out = append(*out, columnLineage{
			DB:     rte.DBName,
			Table:  rte.TableName,
			Column: rte.ColNames[colIdx],
		})
	}
}

// ---------------------------------------------------------------------------
// Helpers for lineage assertions.
// ---------------------------------------------------------------------------

// assertLineage checks that the lineage results match the expected values.
func assertLineage(t *testing.T, got []resultLineage, expected []resultLineage) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("lineage: want %d columns, got %d", len(expected), len(got))
	}
	for i, want := range expected {
		g := got[i]
		if g.Name != want.Name {
			t.Errorf("column[%d].Name: want %q, got %q", i, want.Name, g.Name)
		}
		if len(g.SourceCols) != len(want.SourceCols) {
			t.Errorf("column[%d] %q: want %d sources, got %d: %v",
				i, want.Name, len(want.SourceCols), len(g.SourceCols), g.SourceCols)
			continue
		}
		for j, ws := range want.SourceCols {
			gs := g.SourceCols[j]
			if gs.DB != ws.DB || gs.Table != ws.Table || gs.Column != ws.Column {
				t.Errorf("column[%d] %q source[%d]: want %s.%s.%s, got %s.%s.%s",
					i, want.Name, j,
					ws.DB, ws.Table, ws.Column,
					gs.DB, gs.Table, gs.Column)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Phase 2 lineage tests.
// ---------------------------------------------------------------------------

// TestLineage_14_1_SimpleView tests lineage through a simple view.
func TestLineage_14_1_SimpleView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
		CREATE VIEW emp_names AS SELECT id, name FROM employees;
	`)

	sel := parseSelect(t, "SELECT * FROM emp_names")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	assertLineage(t, lineage, []resultLineage{
		{Name: "id", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "id"}}},
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
	})
}

// TestLineage_14_2_ViewWithExpression tests lineage through a view with an expression.
func TestLineage_14_2_ViewWithExpression(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
		CREATE VIEW emp_salary AS SELECT name, salary * 12 AS annual FROM employees;
	`)

	sel := parseSelect(t, "SELECT * FROM emp_salary")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	assertLineage(t, lineage, []resultLineage{
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
		{Name: "annual", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "salary"}}},
	})
}

// TestLineage_14_3_ViewOfView tests lineage through nested views (view-of-view).
func TestLineage_14_3_ViewOfView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
		CREATE VIEW v1 AS SELECT id, name FROM employees;
		CREATE VIEW v2 AS SELECT name FROM v1;
	`)

	sel := parseSelect(t, "SELECT * FROM v2")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	assertLineage(t, lineage, []resultLineage{
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
	})
}

// TestLineage_14_4_ViewWithJoin tests lineage through a view with a JOIN.
func TestLineage_14_4_ViewWithJoin(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
		CREATE TABLE departments (id INT, name VARCHAR(100), budget DECIMAL(15,2));
		CREATE VIEW emp_dept AS SELECT e.name, d.name AS dept FROM employees e JOIN departments d ON e.department_id = d.id;
	`)

	sel := parseSelect(t, "SELECT * FROM emp_dept")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	assertLineage(t, lineage, []resultLineage{
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
		{Name: "dept", SourceCols: []columnLineage{{DB: "testdb", Table: "departments", Column: "name"}}},
	})
}

func TestLineage_14_4b_ParenthesizedTableReferenceList(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), department_id INT);
		CREATE TABLE departments (id INT, name VARCHAR(100));
	`)

	sel := parseSelect(t, `SELECT e.name AS employee_name, d.name AS department_name FROM (employees e, departments d) WHERE e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	lineage := collectColumnLineage(c, q)
	assertLineage(t, lineage, []resultLineage{
		{Name: "employee_name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
		{Name: "department_name", SourceCols: []columnLineage{{DB: "testdb", Table: "departments", Column: "name"}}},
	})
}

// TestLineage_14_5_TemptableView tests that TEMPTABLE views remain opaque.
func TestLineage_14_5_TemptableView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
		CREATE ALGORITHM=TEMPTABLE VIEW temp_v AS SELECT name FROM employees;
	`)

	sel := parseSelect(t, "SELECT * FROM temp_v")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	// TEMPTABLE view is opaque — lineage stops at the view, not the base table.
	assertLineage(t, lineage, []resultLineage{
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "temp_v", Column: "name"}}},
	})
}

// TestLineage_14_6_CTELineage tests lineage through a CTE.
func TestLineage_14_6_CTELineage(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
	`)

	sel := parseSelect(t, `
		WITH active AS (SELECT id, name FROM employees WHERE is_active = 1)
		SELECT * FROM active
	`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	assertLineage(t, lineage, []resultLineage{
		{Name: "id", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "id"}}},
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
	})
}

// TestLineage_14_7_SubqueryLineage tests lineage through a FROM subquery with aggregate.
func TestLineage_14_7_SubqueryLineage(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
	`)

	sel := parseSelect(t, "SELECT x.total FROM (SELECT COUNT(*) AS total FROM employees) AS x")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	// COUNT(*) has no physical column source — sources should be empty.
	assertLineage(t, lineage, []resultLineage{
		{Name: "total", SourceCols: []columnLineage{}},
	})
}

// TestLineage_14_8_ViewInJoin tests lineage when a view is used in a JOIN.
func TestLineage_14_8_ViewInJoin(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE employees (id INT, name VARCHAR(100), salary DECIMAL(10,2), department_id INT, is_active TINYINT);
		CREATE TABLE departments (id INT, name VARCHAR(100), budget DECIMAL(15,2));
		CREATE VIEW dept_info AS SELECT id, name, budget FROM departments;
	`)

	sel := parseSelect(t, "SELECT e.name, di.budget FROM employees e JOIN dept_info di ON e.department_id = di.id")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	expanded := q.ExpandMergeViews(c)
	lineage := collectColumnLineage(c, expanded)

	assertLineage(t, lineage, []resultLineage{
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
		{Name: "budget", SourceCols: []columnLineage{{DB: "testdb", Table: "departments", Column: "budget"}}},
	})
}
