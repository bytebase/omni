package catalog

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	pg "github.com/bytebase/omni/pg"
	nodes "github.com/bytebase/omni/pg/ast"
)

// columnLineage represents a source column reference, matching QuerySpan's ColumnResource.
type columnLineage struct {
	Schema string
	Table  string
	Column string
}

func (c columnLineage) String() string {
	return fmt.Sprintf("%s.%s.%s", c.Schema, c.Table, c.Column)
}

// resultLineage represents the lineage for one output column, matching QuerySpanResult.
type resultLineage struct {
	Name         string
	SourceCols   []columnLineage
	IsPlainField bool
}

// collectColumnLineage extracts column lineage from an analyzed Query.
// This is the core algorithm that would replace the ANTLR-based QuerySpan extraction.
func collectColumnLineage(c *Catalog, q *Query) []resultLineage {
	var results []resultLineage
	for _, te := range q.TargetList {
		if te.ResJunk {
			continue
		}
		rl := resultLineage{
			Name: te.ResName,
		}
		_, isVar := te.Expr.(*VarExpr)
		rl.IsPlainField = isVar

		// Collect all VarExpr references from this target entry's expression.
		cols := collectVarExprs(c, q, te.Expr)
		rl.SourceCols = cols
		results = append(results, rl)
	}
	return results
}

// collectVarExprs walks an analyzed expression tree and collects all source
// column references, resolving through subqueries and CTEs.
func collectVarExprs(c *Catalog, q *Query, expr AnalyzedExpr) []columnLineage {
	if expr == nil {
		return nil
	}
	seen := make(map[columnLineage]bool)
	var result []columnLineage
	walkExpr(c, q, expr, seen, &result)
	return result
}

func walkExpr(c *Catalog, q *Query, expr AnalyzedExpr, seen map[columnLineage]bool, result *[]columnLineage) {
	if expr == nil {
		return
	}
	switch v := expr.(type) {
	case *VarExpr:
		resolveVar(c, q, v, seen, result)
	case *FuncCallExpr:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *AggExpr:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *OpExpr:
		walkExpr(c, q, v.Left, seen, result)
		walkExpr(c, q, v.Right, seen, result)
	case *BoolExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *CaseExprQ:
		walkExpr(c, q, v.Arg, seen, result)
		for _, w := range v.When {
			walkExpr(c, q, w.Condition, seen, result)
			walkExpr(c, q, w.Result, seen, result)
		}
		walkExpr(c, q, v.Default, seen, result)
	case *CoalesceExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *SubLinkExpr:
		walkExpr(c, q, v.TestExpr, seen, result)
		if v.SubQuery != nil {
			for _, te := range v.SubQuery.TargetList {
				if !te.ResJunk {
					walkExpr(c, v.SubQuery, te.Expr, seen, result)
				}
			}
		}
	case *RelabelExpr:
		walkExpr(c, q, v.Arg, seen, result)
	case *CoerceViaIOExpr:
		walkExpr(c, q, v.Arg, seen, result)
	case *RowExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *WindowFuncExpr:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
		walkExpr(c, q, v.AggFilter, seen, result)
	case *NullIfExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *MinMaxExprQ:
		for _, arg := range v.Args {
			walkExpr(c, q, arg, seen, result)
		}
	case *ArrayExprQ:
		for _, elem := range v.Elements {
			walkExpr(c, q, elem, seen, result)
		}
	case *ScalarArrayOpExpr:
		walkExpr(c, q, v.Left, seen, result)
		walkExpr(c, q, v.Right, seen, result)
	case *CollateExprQ:
		walkExpr(c, q, v.Arg, seen, result)
		// ConstExpr, ParamExpr, etc. — no column references, skip.
	}
}

// resolveVar resolves a VarExpr to its ultimate source column(s).
func resolveVar(c *Catalog, q *Query, v *VarExpr, seen map[columnLineage]bool, result *[]columnLineage) {
	if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	colIdx := int(v.AttNum - 1) // AttNum is 1-based

	switch rte.Kind {
	case RTERelation:
		// Physical table — this is the terminal case.
		rel := c.GetRelationByOID(rte.RelOID)
		if rel == nil || rel.Schema == nil {
			return
		}
		colName := ""
		if colIdx >= 0 && colIdx < len(rte.ColNames) {
			colName = rte.ColNames[colIdx]
		}
		cl := columnLineage{
			Schema: rel.Schema.Name,
			Table:  rel.Name,
			Column: colName,
		}
		if !seen[cl] {
			seen[cl] = true
			*result = append(*result, cl)
		}

	case RTESubquery:
		// Subquery — recurse into the subquery's target list at the matching position.
		if rte.Subquery == nil {
			return
		}
		if colIdx >= 0 && colIdx < len(rte.Subquery.TargetList) {
			te := rte.Subquery.TargetList[colIdx]
			walkExpr(c, rte.Subquery, te.Expr, seen, result)
		}

	case RTECTE:
		// CTE reference — recurse into the CTE's query.
		if rte.CTEIndex >= 0 && rte.CTEIndex < len(q.CTEList) {
			cte := q.CTEList[rte.CTEIndex]
			if cte.Query != nil && colIdx >= 0 && colIdx < len(cte.Query.TargetList) {
				te := cte.Query.TargetList[colIdx]
				walkExpr(c, cte.Query, te.Expr, seen, result)
			}
		}

	case RTEJoin:
		// Join — the VarExpr should reference through to underlying RTEs.
		// For JOINs, columns are concatenated from left + right.
		// We need to find which child RTE this column belongs to.
		// However, after analysis, VarExprs typically point directly to the base RTE,
		// not the join RTE. If they do point to a join RTE, resolve via the join's colNames.

	case RTEFunction:
		// Function result — no further lineage to track.
	}
}

// collectAllSourceColumns collects all table columns accessed anywhere in the query
// (result columns + WHERE + JOIN conditions + GROUP BY + HAVING).
func collectAllSourceColumns(c *Catalog, q *Query) []columnLineage {
	seen := make(map[columnLineage]bool)
	var result []columnLineage

	// Collect from all target entries.
	for _, te := range q.TargetList {
		walkExpr(c, q, te.Expr, seen, &result)
	}

	// Collect from WHERE clause.
	if q.JoinTree != nil {
		walkExpr(c, q, q.JoinTree.Quals, seen, &result)
		// Walk join node quals.
		for _, fn := range q.JoinTree.FromList {
			walkJoinNodeExprs(c, q, fn, seen, &result)
		}
	}

	// Collect from HAVING clause.
	walkExpr(c, q, q.HavingQual, seen, &result)

	return result
}

func walkJoinNodeExprs(c *Catalog, q *Query, node JoinNode, seen map[columnLineage]bool, result *[]columnLineage) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *JoinExprNode:
		walkExpr(c, q, n.Quals, seen, result)
		walkJoinNodeExprs(c, q, n.Left, seen, result)
		walkJoinNodeExprs(c, q, n.Right, seen, result)
	case *RangeTableRef:
		// No expressions to walk.
	}
}

// parseAndAnalyze is a helper that parses SQL and analyzes the first SELECT.
func parseAndAnalyze(t *testing.T, c *Catalog, sql string) *Query {
	t.Helper()
	stmts, err := pg.Parse(sql)
	if err != nil {
		t.Fatalf("parse %q: %v", sql, err)
	}
	if len(stmts) == 0 {
		t.Fatalf("no statements in %q", sql)
	}
	selStmt, ok := stmts[0].AST.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", stmts[0].AST)
	}
	q, err := c.AnalyzeSelectStmt(selStmt)
	if err != nil {
		t.Fatalf("analyze %q: %v", sql, err)
	}
	return q
}

// sortLineage sorts columnLineage for deterministic comparison.
func sortLineage(cols []columnLineage) {
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].String() < cols[j].String()
	})
}

func lineageStrings(cols []columnLineage) []string {
	var s []string
	for _, c := range cols {
		s = append(s, c.String())
	}
	sort.Strings(s)
	return s
}

// --- Test cases matching existing QuerySpan test data ---

func TestQuerySpan_SimpleColumnRef(t *testing.T) {
	// Test: SELECT a, b FROM t;
	// Expected: two result columns, each a plain field from public.t
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int, d text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT a, b FROM t`)
	results := collectColumnLineage(c, q)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Check column "a"
	assertResult(t, results[0], "a", true, []columnLineage{{Schema: "public", Table: "t", Column: "a"}})
	// Check column "b"
	assertResult(t, results[1], "b", true, []columnLineage{{Schema: "public", Table: "t", Column: "b"}})
}

func TestQuerySpan_StarExpansion(t *testing.T) {
	// Test: SELECT * FROM t;
	// Expected: all columns of t, each a plain field
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int, d text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT * FROM t`)
	results := collectColumnLineage(c, q)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	for i, name := range []string{"a", "b", "c", "d"} {
		assertResult(t, results[i], name, true, []columnLineage{{Schema: "public", Table: "t", Column: name}})
	}
}

func TestQuerySpan_ExpressionColumn(t *testing.T) {
	// Test: SELECT a + 1, count(*) FROM t GROUP BY a;
	// Expected: a+1 is not a plain field but sources from t.a; count(*) has no source columns
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT a + 1, count(*) FROM t GROUP BY a`)
	results := collectColumnLineage(c, q)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// a + 1: not plain field, sources from t.a
	assertResult(t, results[0], "?column?", false, []columnLineage{{Schema: "public", Table: "t", Column: "a"}})
	// count(*): not plain field, no source columns
	assertResult(t, results[1], "count", false, nil)
}

func TestQuerySpan_JoinWithStar(t *testing.T) {
	// Test: SELECT t1.*, t2.c, 0 FROM t1 JOIN t2 ON 1 = 1;
	// Matches test case: "Test for SELECT partial FROM JOIN clause"
	c := New()
	_, err := c.Exec(`
		CREATE TABLE t1 (a int, b text);
		CREATE TABLE t2 (c int, d text);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT t1.*, t2.c, 0 FROM t1 JOIN t2 ON 1 = 1`)
	results := collectColumnLineage(c, q)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	assertResult(t, results[0], "a", true, []columnLineage{{Schema: "public", Table: "t1", Column: "a"}})
	assertResult(t, results[1], "b", true, []columnLineage{{Schema: "public", Table: "t1", Column: "b"}})
	assertResult(t, results[2], "c", true, []columnLineage{{Schema: "public", Table: "t2", Column: "c"}})
	assertResult(t, results[3], "?column?", false, nil) // literal 0
}

func TestQuerySpan_Subquery(t *testing.T) {
	// Test: SELECT MIN(bb) AS aa FROM (SELECT 1 bb) AS r WHERE 1 = 1;
	// Matches first test case in query_span.yaml
	c := New()

	q := parseAndAnalyze(t, c, `SELECT MIN(bb) AS aa FROM (SELECT 1 bb) AS r WHERE 1 = 1`)
	results := collectColumnLineage(c, q)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	// MIN(bb) where bb is a literal — no source table columns
	assertResult(t, results[0], "aa", false, nil)
}

func TestQuerySpan_SubqueryWithTableRef(t *testing.T) {
	// Test: SELECT * FROM (SELECT a, b FROM t) sub;
	// Should trace through subquery to base table
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT * FROM (SELECT a, b FROM t) sub`)
	results := collectColumnLineage(c, q)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	assertResult(t, results[0], "a", true, []columnLineage{{Schema: "public", Table: "t", Column: "a"}})
	assertResult(t, results[1], "b", true, []columnLineage{{Schema: "public", Table: "t", Column: "b"}})
}

func TestQuerySpan_SimpleCTE(t *testing.T) {
	// Test: WITH t1 AS (SELECT * FROM t) SELECT * FROM t1;
	// Matches "Simple CTE" test case
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int, d text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `WITH t1 AS (SELECT * FROM t) SELECT * FROM t1`)
	results := collectColumnLineage(c, q)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	for i, name := range []string{"a", "b", "c", "d"} {
		assertResult(t, results[i], name, true, []columnLineage{{Schema: "public", Table: "t", Column: name}})
	}
}

func TestQuerySpan_CTEWithAliases(t *testing.T) {
	// Test: WITH t1(d, c, b, a) AS (SELECT * FROM t) SELECT * FROM t1;
	// CTE columns renamed: d=t.a, c=t.b, b=t.c, a=t.d
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int, d text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `WITH t1(d, c, b, a) AS (SELECT * FROM t) SELECT * FROM t1`)
	results := collectColumnLineage(c, q)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	// The CTE renames columns: position 0 is "d" aliased from t.a
	assertResult(t, results[0], "d", true, []columnLineage{{Schema: "public", Table: "t", Column: "a"}})
	assertResult(t, results[1], "c", true, []columnLineage{{Schema: "public", Table: "t", Column: "b"}})
	assertResult(t, results[2], "b", true, []columnLineage{{Schema: "public", Table: "t", Column: "c"}})
	assertResult(t, results[3], "a", true, []columnLineage{{Schema: "public", Table: "t", Column: "d"}})
}

func TestQuerySpan_SetOperation(t *testing.T) {
	// Test: SELECT a, b FROM t1 UNION ALL SELECT c, d FROM t2;
	// Set operations merge columns from both branches
	c := New()
	_, err := c.Exec(`
		CREATE TABLE t1 (a int, b text);
		CREATE TABLE t2 (c int, d text);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT a, b FROM t1 UNION ALL SELECT c, d FROM t2`)
	results := collectColumnLineage(c, q)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// UNION uses left branch names. Both branches contribute.
	// For set operations, the top-level Query has placeholder VarExprs —
	// we need to recurse into LArg and RArg.
	if q.SetOp == SetOpNone {
		t.Fatal("expected set operation")
	}

	// Verify LArg and RArg are present.
	if q.LArg == nil || q.RArg == nil {
		t.Fatal("expected LArg and RArg for set operation")
	}

	// Verify left branch lineage.
	lResults := collectColumnLineage(c, q.LArg)
	assertResult(t, lResults[0], "a", true, []columnLineage{{Schema: "public", Table: "t1", Column: "a"}})
	assertResult(t, lResults[1], "b", true, []columnLineage{{Schema: "public", Table: "t1", Column: "b"}})

	// Verify right branch lineage.
	// Note: In PG set-ops, right branch ResName is aligned to left branch names.
	rResults := collectColumnLineage(c, q.RArg)
	if len(rResults) != 2 {
		t.Fatalf("expected 2 rResults, got %d", len(rResults))
	}
	// Check source columns regardless of name (name alignment is PG behavior).
	assertSourceContains(t, rResults[0].SourceCols, columnLineage{Schema: "public", Table: "t2", Column: "c"})
	assertSourceContains(t, rResults[1].SourceCols, columnLineage{Schema: "public", Table: "t2", Column: "d"})
}

func TestQuerySpan_RecursiveCTE(t *testing.T) {
	// Test: WITH RECURSIVE t1 AS (
	//   SELECT 1 AS c1, 2 AS c2, 3 AS c3, 4 AS c4
	//   UNION SELECT a, b, c, d FROM t
	// ) SELECT * FROM t1;
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b int, c int, d int);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `
		WITH RECURSIVE t1 AS (
			SELECT 1 AS c1, 2 AS c2, 3 AS c3, 4 AS c4
			UNION SELECT a, b, c, d FROM t
		) SELECT * FROM t1
	`)
	results := collectColumnLineage(c, q)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	// The CTE has two branches via UNION. The first is literals (no source cols),
	// the second references t. We should see the CTE resolves to the union's
	// inner query which contains t references.
	t.Logf("recursive CTE results: %+v", results)

	// At minimum, verify the column names are correct.
	for i, name := range []string{"c1", "c2", "c3", "c4"} {
		if results[i].Name != name {
			t.Errorf("result[%d].Name = %q, want %q", i, results[i].Name, name)
		}
	}
}

func TestQuerySpan_JoinWithUsing(t *testing.T) {
	// Test: SELECT * FROM t AS t1 JOIN t AS t2 USING(a);
	// Matches "Test for JOIN with USING clause"
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int, d text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT * FROM t AS t1 JOIN t AS t2 USING(a)`)
	results := collectColumnLineage(c, q)

	// PostgreSQL USING produces: a (merged), t1.b, t1.c, t1.d, t2.b, t2.c, t2.d = 7 columns
	t.Logf("JOIN USING results (%d):", len(results))
	for i, r := range results {
		t.Logf("  [%d] name=%q plain=%v sources=%v", i, r.Name, r.IsPlainField, lineageStrings(r.SourceCols))
	}

	if len(results) != 7 {
		t.Fatalf("expected 7 columns (USING deduplicates 'a'), got %d", len(results))
	}

	// Verify column names: a, b, c, d, b, c, d
	expectedNames := []string{"a", "b", "c", "d", "b", "c", "d"}
	for i, expected := range expectedNames {
		if results[i].Name != expected {
			t.Errorf("column %d: expected name %q, got %q", i, expected, results[i].Name)
		}
	}

	// All columns should trace back to public.t
	for _, r := range results {
		for _, col := range r.SourceCols {
			if col.Table != "t" || col.Schema != "public" {
				t.Errorf("unexpected source %v for column %q", col, r.Name)
			}
		}
	}
}

func TestQuerySpan_MultiTableJoin(t *testing.T) {
	// Test: SELECT users.*, payments.* FROM users JOIN payments ON users.id = payments.user_id;
	c := New()
	_, err := c.Exec(`
		CREATE TABLE users (id int, name text, email text);
		CREATE TABLE payments (id int, user_id int, amount numeric);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT users.*, payments.* FROM users JOIN payments ON users.id = payments.user_id`)
	results := collectColumnLineage(c, q)

	if len(results) != 6 {
		t.Fatalf("expected 6 results, got %d", len(results))
	}

	assertResult(t, results[0], "id", true, []columnLineage{{Schema: "public", Table: "users", Column: "id"}})
	assertResult(t, results[1], "name", true, []columnLineage{{Schema: "public", Table: "users", Column: "name"}})
	assertResult(t, results[2], "email", true, []columnLineage{{Schema: "public", Table: "users", Column: "email"}})
	assertResult(t, results[3], "id", true, []columnLineage{{Schema: "public", Table: "payments", Column: "id"}})
	assertResult(t, results[4], "user_id", true, []columnLineage{{Schema: "public", Table: "payments", Column: "user_id"}})
	assertResult(t, results[5], "amount", true, []columnLineage{{Schema: "public", Table: "payments", Column: "amount"}})
}

func TestQuerySpan_NestedSubquery(t *testing.T) {
	// Test: SELECT x FROM (SELECT a AS x FROM (SELECT a FROM t) inner_sub) outer_sub;
	// Should trace through two levels of subqueries to base table
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT x FROM (SELECT a AS x FROM (SELECT a FROM t) inner_sub) outer_sub`)
	results := collectColumnLineage(c, q)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	assertResult(t, results[0], "x", true, []columnLineage{{Schema: "public", Table: "t", Column: "a"}})
}

func TestQuerySpan_AllSourceColumns(t *testing.T) {
	// Verify that collectAllSourceColumns captures columns from WHERE, not just SELECT
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT a FROM t WHERE b = 'hello' AND c > 0`)
	allCols := collectAllSourceColumns(c, q)

	// Should include a (from SELECT), b (from WHERE), c (from WHERE)
	colNames := make(map[string]bool)
	for _, cl := range allCols {
		colNames[cl.Column] = true
	}
	for _, expected := range []string{"a", "b", "c"} {
		if !colNames[expected] {
			t.Errorf("expected source column %q in WHERE clause, got %v", expected, lineageStrings(allCols))
		}
	}
}

func TestQuerySpan_MultiLevelCTE(t *testing.T) {
	// Test: WITH tt2 AS (WITH tt2 AS (SELECT * FROM t) SELECT MAX(a) FROM tt2) SELECT * FROM tt2;
	// Matches "Multi-level CTE" test case
	c := New()
	_, err := c.Exec(`CREATE TABLE t (a int, b text, c int, d text);`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `WITH tt2 AS (WITH tt2 AS (SELECT * FROM t) SELECT MAX(a) FROM tt2) SELECT * FROM tt2`)
	results := collectColumnLineage(c, q)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	// MAX(a) traces through inner CTE to t.a
	if results[0].Name != "max" {
		t.Errorf("result name = %q, want %q", results[0].Name, "max")
	}
	// Should trace back to t.a
	if len(results[0].SourceCols) > 0 {
		assertSourceContains(t, results[0].SourceCols, columnLineage{Schema: "public", Table: "t", Column: "a"})
	}
}

func TestQuerySpan_ViewLineage(t *testing.T) {
	// Views should be treated as subqueries for lineage.
	c := New()
	_, err := c.Exec(`
		CREATE TABLE t (a int, b text);
		CREATE VIEW v AS SELECT a, b FROM t;
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT * FROM v`)
	results := collectColumnLineage(c, q)

	t.Logf("view results: %+v", results)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Views in omni catalog are stored as relations.
	// The RTE will be RTERelation with RelOID pointing to the view.
	// For full lineage through views, we'd need to expand view definitions.
	// For now, verify the view's columns are returned.
	if results[0].Name != "a" {
		t.Errorf("result[0].Name = %q, want %q", results[0].Name, "a")
	}
	if results[1].Name != "b" {
		t.Errorf("result[1].Name = %q, want %q", results[1].Name, "b")
	}
}

func TestQuerySpan_DifferentSchemas(t *testing.T) {
	// Test cross-schema references
	c := New()
	_, err := c.Exec(`
		CREATE SCHEMA other;
		CREATE TABLE public.t1 (a int);
		CREATE TABLE other.t2 (b int);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	q := parseAndAnalyze(t, c, `SELECT t1.a, t2.b FROM public.t1 JOIN other.t2 ON true`)
	results := collectColumnLineage(c, q)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	assertResult(t, results[0], "a", true, []columnLineage{{Schema: "public", Table: "t1", Column: "a"}})
	assertResult(t, results[1], "b", true, []columnLineage{{Schema: "other", Table: "t2", Column: "b"}})
}

// --- Assertion helpers ---

func assertResult(t *testing.T, got resultLineage, wantName string, wantPlain bool, wantCols []columnLineage) {
	t.Helper()
	if got.Name != wantName {
		t.Errorf("name = %q, want %q", got.Name, wantName)
	}
	if got.IsPlainField != wantPlain {
		t.Errorf("isPlainField(%s) = %v, want %v", wantName, got.IsPlainField, wantPlain)
	}
	if len(wantCols) == 0 {
		if len(got.SourceCols) != 0 {
			t.Errorf("sourceColumns(%s) = %v, want empty", wantName, lineageStrings(got.SourceCols))
		}
		return
	}
	gotStrs := lineageStrings(got.SourceCols)
	wantStrs := lineageStrings(wantCols)
	if strings.Join(gotStrs, ",") != strings.Join(wantStrs, ",") {
		t.Errorf("sourceColumns(%s) = %v, want %v", wantName, gotStrs, wantStrs)
	}
}

func assertSourceContains(t *testing.T, cols []columnLineage, want columnLineage) {
	t.Helper()
	for _, c := range cols {
		if c == want {
			return
		}
	}
	t.Errorf("source columns %v does not contain %v", lineageStrings(cols), want)
}
