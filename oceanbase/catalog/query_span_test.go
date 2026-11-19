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
	case *QuantifiedExprQ:
		collectVarExprs(q, e.Arg, out)
	case *OdbcEscapeExprQ:
		collectVarExprs(q, e.Arg, out)
	case *GroupingSetExprQ:
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

// TestLineage_14_5_TemptableView tests that a view written ALGORITHM=TEMPTABLE is read
// through like any other: OceanBase keeps no algorithm
// (TestOracle_OceanBaseViewAttributes), where MySQL's TEMPTABLE view is opaque.
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

	assertLineage(t, lineage, []resultLineage{
		{Name: "name", SourceCols: []columnLineage{{DB: "testdb", Table: "employees", Column: "name"}}},
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

func TestLineageQuantifiedExprAndRemoteViewIsolation(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE orders (id INT, tags INT, secret INT);
		CREATE VIEW v_quant AS SELECT id = ANY(tags) AS hit_any, id = ALL(tags) AS hit_all FROM orders;
		CREATE VIEW v_remote AS SELECT * FROM orders@remote_link;
		CREATE VIEW v_cat AS SELECT * FROM ext_cat.other_db.orders;
	`)

	db := c.GetDatabase("testdb")
	vq := db.Views["v_quant"]
	if vq == nil || vq.AnalyzedQuery == nil {
		t.Fatal("expected v_quant to have AnalyzedQuery")
	}
	deparsedQ := DeparseQuery(vq.AnalyzedQuery)
	if want := "select (`orders`.`id` = any (`orders`.`tags`)) as `hit_any`, (`orders`.`id` = all (`orders`.`tags`)) as `hit_all` from `testdb`.`orders`"; deparsedQ != want {
		t.Errorf("DeparseQuery(v_quant) = %q, want %q", deparsedQ, want)
	}

	sel := parseSelect(t, "SELECT hit_any, hit_all FROM v_quant")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)
	lineage := collectColumnLineage(c, q.ExpandMergeViews(c))
	assertLineage(t, lineage, []resultLineage{
		{Name: "hit_any", SourceCols: []columnLineage{
			{DB: "testdb", Table: "orders", Column: "id"},
			{DB: "testdb", Table: "orders", Column: "tags"},
		}},
		{Name: "hit_all", SourceCols: []columnLineage{
			{DB: "testdb", Table: "orders", Column: "id"},
			{DB: "testdb", Table: "orders", Column: "tags"},
		}},
	})

	// Aggregate inside QuantifiedExpr sets HasAggs.
	aggQ, err := c.AnalyzeSelectStmt(parseSelect(t, "SELECT id = ANY(MAX(tags)) FROM orders"))
	assertNoError(t, err)
	if !aggQ.HasAggs {
		t.Error("expected HasAggs=true for aggregate inside QuantifiedExpr")
	}

	// Remote/catalog views must not expand `*` to local `orders` columns (`id`, `tags`, `secret`).
	for _, viewName := range []string{"v_remote", "v_cat"} {
		v := db.Views[viewName]
		if v == nil {
			t.Fatalf("missing view %s", viewName)
		}
		for _, col := range v.Columns {
			if col == "secret" {
				t.Errorf("view %s leaked local table column %q into view columns %v", viewName, col, v.Columns)
			}
		}
	}
}

func TestFullJoinQualifiedStarsKeepSideColumns(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE l (id INT, lx INT);
		CREATE TABLE r (id INT, rx INT);
		CREATE VIEW v_star AS SELECT * FROM l FULL JOIN r USING (id);
		CREATE VIEW v_lstar AS SELECT l.* FROM l FULL JOIN r USING (id);
		CREATE VIEW v_rstar AS SELECT r.* FROM l FULL JOIN r USING (id);
	`)

	db := c.GetDatabase("testdb")
	if got, want := DeparseQuery(db.Views["v_star"].AnalyzedQuery),
		"select coalesce(`l`.`id`, `r`.`id`) as `id`, `l`.`lx`, `r`.`rx` from `testdb`.`l` full join `testdb`.`r` using (`id`)"; got != want {
		t.Errorf("DeparseQuery(v_star) =\n  %s\nwant\n  %s", got, want)
	}
	if got, want := DeparseQuery(db.Views["v_lstar"].AnalyzedQuery),
		"select `l`.`id`, `l`.`lx` from `testdb`.`l` full join `testdb`.`r` using (`id`)"; got != want {
		t.Errorf("DeparseQuery(v_lstar) =\n  %s\nwant\n  %s", got, want)
	}
	if got, want := DeparseQuery(db.Views["v_rstar"].AnalyzedQuery),
		"select `r`.`id`, `r`.`rx` from `testdb`.`l` full join `testdb`.`r` using (`id`)"; got != want {
		t.Errorf("DeparseQuery(v_rstar) =\n  %s\nwant\n  %s", got, want)
	}

	lQuery, err := c.AnalyzeSelectStmt(parseSelect(t, "SELECT id, lx FROM v_lstar"))
	assertNoError(t, err)
	assertLineage(t, collectColumnLineage(c, lQuery.ExpandMergeViews(c)), []resultLineage{
		{Name: "id", SourceCols: []columnLineage{{DB: "testdb", Table: "l", Column: "id"}}},
		{Name: "lx", SourceCols: []columnLineage{{DB: "testdb", Table: "l", Column: "lx"}}},
	})
}

func TestRemoteCatalogQualifiedColumnResolution(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE t (a INT);
		CREATE VIEW v_multi_cat AS SELECT c1.d.t.a AS a1, c2.d.t.b AS b2 FROM c1.d.t, c2.d.t;
	`)

	vq := c.GetDatabase("testdb").Views["v_multi_cat"].AnalyzedQuery
	if vq == nil {
		t.Fatal("expected v_multi_cat to have AnalyzedQuery")
	}
	v0, ok0 := vq.TargetList[0].Expr.(*VarExprQ)
	v1, ok1 := vq.TargetList[1].Expr.(*VarExprQ)
	if !ok0 || !ok1 || v0.RangeIdx != 0 || v1.RangeIdx != 1 {
		t.Fatalf("expected c1->RangeIdx 0 and c2->RangeIdx 1, got %+v and %+v", vq.TargetList[0].Expr, vq.TargetList[1].Expr)
	}
	if got, want := DeparseQuery(vq), "select `c1`.`d`.`t`.`a` as `a1`, `c2`.`d`.`t`.`b` as `b2` from `c1`.`d`.`t`, `c2`.`d`.`t`"; got != want {
		t.Errorf("DeparseQuery(v_multi_cat) =\n  %s\nwant\n  %s", got, want)
	}

	// Unqualified-catalog `t.a` across two remote `t` relations is ambiguous (1052).
	if _, err := c.AnalyzeSelectStmt(parseSelect(t, "SELECT t.a FROM c1.d.t, c2.d.t")); err == nil {
		t.Fatal("expected ambiguity error for SELECT t.a FROM c1.d.t, c2.d.t")
	} else if catErr, ok := err.(*Error); !ok || catErr.Code != 1052 {
		t.Fatalf("expected error 1052, got %v", err)
	}

	// Four-part `c1.d.t.a` alongside local `t` binds to `c1.d.t` (RangeIdx 1), not local `t` (RangeIdx 0).
	mixedQ, err := c.AnalyzeSelectStmt(parseSelect(t, "SELECT c1.d.t.a FROM t, c1.d.t"))
	assertNoError(t, err)
	if mv, ok := mixedQ.TargetList[0].Expr.(*VarExprQ); !ok || mv.RangeIdx != 1 {
		t.Fatalf("expected c1.d.t.a to bind to RangeIdx 1, got %+v", mixedQ.TargetList[0].Expr)
	}
}
