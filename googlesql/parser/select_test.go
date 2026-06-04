package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-select node (the GoogleSQL query stack: SELECT,
// FROM, joins, set-ops, CTEs, GROUP BY, HAVING, QUALIFY, WINDOW, ORDER/LIMIT,
// UNNEST). The accept/reject behavior is also proven against the live Cloud
// Spanner emulator in select_oracle_test.go (build tag googlesql_oracle); these
// hand-written tests assert the AST STRUCTURE (the structural gate — accept/
// reject alone does not catch wrong nesting/precedence) and cover BigQuery-only
// forms the Spanner emulator feature-rejects (QUALIFY, WITH RECURSIVE,
// CORRESPONDING) which are non-authoritative against that oracle.

// parseQ parses a single query statement and fails the test on any parse error.
// It returns the *ast.QueryStmt (every query_statement produces one).
func parseQ(t *testing.T, sql string) *ast.QueryStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): unexpected errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, len(file.Stmts))
	}
	q, ok := file.Stmts[0].(*ast.QueryStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.QueryStmt", sql, file.Stmts[0])
	}
	return q
}

// selectOf returns the QueryStmt's body as a *ast.SelectStmt, failing if it is
// not a bare SELECT.
func selectOf(t *testing.T, q *ast.QueryStmt) *ast.SelectStmt {
	t.Helper()
	s, ok := q.Body.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("query body is %T, want *ast.SelectStmt", q.Body)
	}
	return s
}

// --- accept-parity: every documented form parses cleanly ---------------------

// TestSelect_AcceptForms asserts every documented BigQuery + Spanner query form
// (truth1/{bigquery,spanner}/query.md) parses without error. This is the
// docs-parity accept side; the oracle differential (select_oracle_test.go)
// confirms the Spanner-authoritative subset against the live emulator.
func TestSelect_AcceptForms(t *testing.T) {
	accept := []string{
		// SELECT core.
		"SELECT 1",
		"SELECT name, age FROM dataset.users WHERE age > 18 ORDER BY age DESC LIMIT 10",
		"SELECT * FROM Roster",
		"SELECT name AS n, age FROM dataset.users",
		"SELECT name n FROM t", // implicit alias
		"SELECT g.* FROM groceries AS g",
		"SELECT l.location.* FROM locations l",
		"SELECT * EXCEPT (order_id) FROM orders",
		"SELECT * REPLACE (\"widget\" AS item_name) FROM orders",
		"SELECT * REPLACE (quantity/2 AS quantity) FROM orders",
		"SELECT * EXCEPT (a) REPLACE (b+1 AS b) FROM t",
		"SELECT DISTINCT Name FROM PlayerStats",
		"SELECT ALL name FROM dataset.users",
		"SELECT AS STRUCT 1 AS a, 2 AS b",
		"SELECT AS VALUE STRUCT(1 AS a, 2 AS b)",
		"SELECT AS myproto.Type 1 AS a", // SELECT AS <path>
		// FROM / paths / time-travel.
		"SELECT * FROM dataset.Roster",
		"SELECT * FROM project.dataset.Roster",
		"SELECT * FROM (SELECT \"apple\" AS fruit) AS t",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2017-01-01 10:00:00-07:00'",
		"SELECT a.Title FROM myschema.Albums AS a",
		"SELECT * FROM `my-project`.dataset.table",
		// UNNEST.
		"SELECT * FROM UNNEST([1, 2, 3]) AS num",
		"SELECT num, pos FROM UNNEST([10, 20, 30]) AS num WITH OFFSET AS pos",
		"SELECT SingerId, album FROM Singers, UNNEST(AlbumList) AS album",
		"SELECT * FROM T1 t1, t1.array_column",
		"SELECT * FROM T1 t1, t1.struct_column.array_field",
		// joins.
		"SELECT * FROM a LEFT JOIN b ON a.id = b.id RIGHT JOIN c USING (id)",
		"SELECT Roster.LastName FROM Roster JOIN TeamMascot ON Roster.SchoolID = TeamMascot.SchoolID",
		"SELECT Roster.LastName FROM Roster CROSS JOIN TeamMascot",
		"SELECT * FROM A, B",
		"SELECT s.FirstName FROM Singers s FULL OUTER JOIN Albums a ON s.SingerId = a.SingerId",
		"SELECT s.FirstName FROM Singers s INNER JOIN Albums a ON s.SingerId = a.SingerId",
		"SELECT * FROM Singers s HASH JOIN Albums a ON s.SingerId = a.SingerId",
		"SELECT * FROM (a JOIN b ON a.x = b.x) JOIN c ON a.y = c.y",
		"SELECT * FROM a JOIN b", // INNER JOIN with no ON/USING — oracle-confirmed ACCEPT
		// GROUP BY variants.
		"SELECT SUM(PointsScored) AS total_points, LastName FROM PlayerStats GROUP BY LastName",
		"SELECT SUM(x), LastName, FirstName FROM PlayerStats GROUP BY 2, 3",
		"SELECT a FROM t GROUP BY ALL",
		"SELECT product_type FROM Products GROUP BY GROUPING SETS (product_type, product_name, ())",
		"SELECT product_type FROM Products GROUP BY ROLLUP (product_type, product_name)",
		"SELECT product_type FROM Products GROUP BY CUBE (product_type, product_name)",
		"SELECT SingerId, SUM(Revenue) FROM Sales GROUP BY ROLLUP (SingerId, AlbumId)",
		// HAVING / QUALIFY / WINDOW.
		"SELECT LastName FROM Roster GROUP BY LastName HAVING SUM(PointsScored) > 15",
		"SELECT item FROM Produce QUALIFY RANK() OVER (PARTITION BY category ORDER BY purchases DESC) <= 3",
		"SELECT SingerId, ROW_NUMBER() OVER (PARTITION BY SingerId ORDER BY MarketingBudget DESC) AS rn FROM Albums QUALIFY rn = 1",
		"SELECT item, LAST_VALUE(item) OVER (item_window) AS most_popular FROM Produce WINDOW item_window AS (PARTITION BY category ORDER BY purchases ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING)",
		"SELECT item, LAST_VALUE(item) OVER (d) AS most_popular FROM Produce WINDOW a AS (PARTITION BY category), b AS (a ORDER BY purchases), c AS (b ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING), d AS (c)",
		// ORDER BY / LIMIT.
		"SELECT x, y FROM t ORDER BY x NULLS LAST",
		"SELECT x, y FROM t ORDER BY x DESC NULLS FIRST",
		"SELECT * FROM Singers LIMIT 10 OFFSET 20",
		"SELECT * FROM Singers ORDER BY LastName ASC, FirstName DESC",
		// set-ops.
		"SELECT * FROM Roster UNION ALL SELECT * FROM TeamMascot ORDER BY SchoolID",
		"SELECT * FROM Roster UNION DISTINCT SELECT * FROM TeamMascot",
		"SELECT * FROM A INTERSECT DISTINCT SELECT * FROM B",
		"SELECT * FROM A EXCEPT DISTINCT SELECT * FROM B",
		"(SELECT 1 AS n UNION ALL SELECT 2 AS n) ORDER BY n DESC",
		"(SELECT 1 AS x) UNION ALL (SELECT 2 AS x)",
		// CTEs.
		"WITH subQ1 AS (SELECT SchoolID FROM Roster), subQ2 AS (SELECT OpponentID FROM PlayerStats) SELECT * FROM subQ1 UNION ALL SELECT * FROM subQ2",
		"WITH TopSingers AS (SELECT SingerId FROM Singers ORDER BY LastName LIMIT 10) SELECT * FROM TopSingers WHERE FirstName LIKE 'A%'",
		"WITH Cte1 AS (SELECT 1 AS n), Cte2 AS (SELECT n + 1 AS n FROM Cte1) SELECT * FROM Cte2",
		"WITH Cte1 (n) AS (SELECT 1) SELECT * FROM Cte1", // Spanner explicit columns
		// subqueries.
		"SELECT (SELECT MAX(MarketingBudget) FROM Albums) AS MaxBudget",
		"SELECT * FROM Singers WHERE EXISTS (SELECT 1 FROM Albums WHERE Albums.SingerId = Singers.SingerId)",
		"SELECT * FROM Singers WHERE SingerId IN (SELECT SingerId FROM Albums WHERE Title LIKE 'A%')",
		"SELECT * FROM (SELECT SingerId, COUNT(*) AS albums FROM Albums GROUP BY SingerId) AS sub WHERE albums > 2",
		// Spanner FOR UPDATE.
		"SELECT * FROM Singers WHERE SingerId = 1 FOR UPDATE",
		// TVF.
		"SELECT * FROM my_tvf(arg1, arg2)",
		"SELECT * FROM mydataset.my_tvf('param') AS t",
		// statement-level hint + SELECT-level hint.
		"@{OPTIMIZER_VERSION=2} SELECT * FROM T",
		"SELECT @{FORCE_INDEX=idx} a FROM t",
		// aggregation threshold / differential privacy SELECT WITH.
		"SELECT WITH AGGREGATION_THRESHOLD OPTIONS (threshold=50, privacy_unit_column=id) zip_code, COUNT(*) FROM dataset.users GROUP BY zip_code",
		// trailing comma in select list.
		"SELECT a, b, FROM t",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			file, errs := Parse(sql)
			if len(errs) != 0 {
				t.Fatalf("Parse(%q): want accept, got errors: %v", sql, errs)
			}
			if len(file.Stmts) != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, len(file.Stmts))
			}
		})
	}
}

// TestSelect_BigQueryOnlyForms covers query forms valid in BigQuery that the
// Spanner emulator FEATURE-rejects (parses, then "X is not supported") — so the
// emulator is non-authoritative and they are proven only here, triangulated
// against truth1/bigquery + the legacy ZetaSQL .g4 (divergence ledger #9, #11).
func TestSelect_BigQueryOnlyForms(t *testing.T) {
	for _, sql := range []string{
		// QUALIFY without a preceding window in the SELECT list is BigQuery-valid;
		// divergence #9 (Spanner: "QUALIFY is not supported").
		"SELECT item FROM Produce WHERE Produce.category = 'vegetable' QUALIFY RANK() OVER (PARTITION BY category ORDER BY purchases DESC) <= 3",
		// WITH RECURSIVE: BigQuery-valid; divergence #11 (Spanner: "RECURSIVE is
		// not supported in the WITH clause").
		"WITH RECURSIVE T1 AS ( (SELECT 1 AS n) UNION ALL (SELECT n + 1 AS n FROM T1 WHERE n < 3) ) SELECT n FROM T1",
		"WITH RECURSIVE T1 AS ( (SELECT 1 AS n) UNION ALL (SELECT n + 2 FROM T1 WHERE n < 4)) SELECT * FROM T1 ORDER BY n",
		// Set-op CORRESPONDING / BY NAME-style modifiers (BigQuery): Spanner
		// feature-rejects CORRESPONDING.
		"SELECT 1 AS a UNION ALL STRICT CORRESPONDING SELECT 2 AS a",
		"SELECT 1 AS a FULL OUTER UNION ALL SELECT 2 AS a",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, errs := Parse(sql); len(errs) != 0 {
				t.Errorf("Parse(%q): want accept (BigQuery-valid, union grammar), got: %v", sql, errs)
			}
		})
	}
}

// --- negative tests (required by correctness-protocol) -----------------------

// TestSelect_RejectForms asserts invalid query syntax is rejected (a parse error
// and no statement produced). Without negatives a parser drifts over-permissive.
func TestSelect_RejectForms(t *testing.T) {
	reject := []string{
		"SELECT",                       // empty select list
		"SELECT FROM t",                // empty list before FROM
		"SELECT 1 2 3",                 // two implicit aliases
		"SELECT * FROM",                // FROM with no source
		"SELECT * FROM t WHERE",        // WHERE with no expr
		"SELECT * FROM t GROUP BY",     // GROUP BY with no items
		"SELECT * FROM a JOIN",         // JOIN with no right source
		"SELECT * FROM a CROSS b",      // CROSS without JOIN
		"SELECT * FROM a USING (x)",    // USING with no join
		"FROM t",                       // FROM-first query (oracle rejects)
		"SELECT * FROM t UNION",        // set-op with no right query
		"SELECT * FROM t UNION SELECT 1", // UNION without ALL/DISTINCT
		"WITH c SELECT 1",              // CTE missing AS (query)
		"WITH c AS SELECT 1",           // CTE body not parenthesized
		"WITH c AS (SELECT 1)",         // WITH with no trailing query body
		"SELECT * FROM t ORDER",        // ORDER without BY
		"SELECT * FROM t LIMIT",        // LIMIT with no count
		"SELECT * FROM t GROUP BY ROLLUP", // ROLLUP with no parens
		"SELECT a, FROM",               // trailing comma then nothing
		"(SELECT 1",                    // unbalanced paren
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) {
			file, errs := Parse(sql)
			if len(errs) == 0 {
				t.Errorf("Parse(%q): want reject, got accept (stmts=%d)", sql, len(file.Stmts))
			}
		})
	}
}

// --- structural gate: AST shape ----------------------------------------------

// TestSelect_OrderLimitBindToQuery confirms the ZetaSQL nuance that trailing
// ORDER BY / LIMIT / FOR UPDATE attach to the whole query (QueryStmt), not to the
// inner SELECT or the last set-op arm.
func TestSelect_OrderLimitBindToQuery(t *testing.T) {
	q := parseQ(t, "SELECT 1 AS n UNION ALL SELECT 2 AS n ORDER BY n LIMIT 5 OFFSET 1")
	if _, ok := q.Body.(*ast.SetOperation); !ok {
		t.Fatalf("body = %T, want *ast.SetOperation", q.Body)
	}
	if len(q.OrderBy) != 1 {
		t.Errorf("ORDER BY items = %d, want 1 (on the union)", len(q.OrderBy))
	}
	if q.Limit == nil || q.Offset == nil {
		t.Errorf("LIMIT/OFFSET = %v/%v, want both set on the query", q.Limit, q.Offset)
	}
	// The inner SELECTs must NOT carry the ORDER BY.
	so := q.Body.(*ast.SetOperation)
	if rs, ok := so.Right.(*ast.SelectStmt); ok {
		_ = rs // SelectStmt has no OrderBy field by design — compile-time guarantee
	}
}

// TestSelect_SetOpLeftAssoc confirms a uniform-op set-operation chain nests
// left-associatively and carries its operator + ALL/DISTINCT. GoogleSQL requires
// a flat chain to use the SAME operation throughout (see
// TestSelect_MixedSetOpRejected); left-associativity is observed with a repeated
// operator.
func TestSelect_SetOpLeftAssoc(t *testing.T) {
	q := parseQ(t, "SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3")
	so, ok := q.Body.(*ast.SetOperation)
	if !ok {
		t.Fatalf("body = %T, want *ast.SetOperation", q.Body)
	}
	if so.Op != ast.SetOpUnion || !so.All {
		t.Errorf("outer op = %s all=%v, want UNION ALL", so.Op, so.All)
	}
	left, ok := so.Left.(*ast.SetOperation)
	if !ok {
		t.Fatalf("outer.Left = %T, want nested *ast.SetOperation", so.Left)
	}
	if left.Op != ast.SetOpUnion || !left.All {
		t.Errorf("inner op = %s all=%v, want UNION ALL", left.Op, left.All)
	}

	// A mixed-op chain becomes valid once the differing operations are grouped
	// with parentheses: `(a UNION ALL b) INTERSECT DISTINCT c`.
	q = parseQ(t, "(SELECT 1 UNION ALL SELECT 2) INTERSECT DISTINCT SELECT 3")
	so = q.Body.(*ast.SetOperation)
	if so.Op != ast.SetOpIntersect || so.All {
		t.Errorf("grouped chain outer op = %s all=%v, want INTERSECT DISTINCT", so.Op, so.All)
	}
	if lq, ok := so.Left.(*ast.QueryStmt); !ok || !lq.Parens {
		t.Errorf("grouped chain left = %T, want parenthesized *ast.QueryStmt", so.Left)
	}
}

// TestSelect_MixedSetOpRejected confirms GoogleSQL's rule that a flat
// (unparenthesized) set-operation chain must use the same operation throughout;
// mixing operators or quantifiers is a syntax error requiring parentheses
// (oracle-confirmed: "Different set operations cannot be used in the same query
// without using parentheses for grouping").
func TestSelect_MixedSetOpRejected(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1 UNION ALL SELECT 2 INTERSECT DISTINCT SELECT 3",
		"SELECT 1 UNION ALL SELECT 2 UNION DISTINCT SELECT 3",
		"SELECT 1 INTERSECT DISTINCT SELECT 2 EXCEPT DISTINCT SELECT 3",
	} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q): want reject (mixed set operations), got accept", sql)
		}
	}
}

// TestSelect_DotStarModifiers confirms `expr.* EXCEPT(...)` keeps the qualifier
// and the EXCEPT columns.
func TestSelect_DotStarModifiers(t *testing.T) {
	q := parseQ(t, "SELECT t.* EXCEPT (a, b) FROM tbl t")
	sel := selectOf(t, q)
	it := sel.Items[0]
	if !it.Star {
		t.Fatal("item.Star = false, want true (dot-star)")
	}
	pe, ok := it.Expr.(*ast.PathExpr)
	if !ok || len(pe.Parts) != 1 || pe.Parts[0] != "t" {
		t.Errorf("dot-star qualifier = %v, want PathExpr{t}", it.Expr)
	}
	if it.Modifiers == nil || len(it.Modifiers.Except) != 2 {
		t.Errorf("EXCEPT columns = %v, want 2", it.Modifiers)
	}
}

// TestSelect_SubqueryFilled confirms an expression-embedded subquery's Query is
// re-parsed into a real *QueryStmt (the subquery seam with the expressions node).
func TestSelect_SubqueryFilled(t *testing.T) {
	q := parseQ(t, "SELECT (SELECT MAX(x) FROM t2) AS m FROM t1")
	sel := selectOf(t, q)
	sub, ok := sel.Items[0].Expr.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("select item expr = %T, want *ast.SubqueryExpr", sel.Items[0].Expr)
	}
	if sub.Query == nil {
		t.Fatal("SubqueryExpr.Query is nil; parser-select must fill it")
	}
	inner, ok := sub.Query.(*ast.QueryStmt)
	if !ok {
		t.Fatalf("SubqueryExpr.Query = %T, want *ast.QueryStmt", sub.Query)
	}
	if _, ok := inner.Body.(*ast.SelectStmt); !ok {
		t.Errorf("inner query body = %T, want *ast.SelectStmt", inner.Body)
	}

	// EXISTS and ARRAY subqueries are likewise filled.
	q = parseQ(t, "SELECT EXISTS(SELECT 1 FROM t2) FROM t1")
	sel = selectOf(t, q)
	ex := sel.Items[0].Expr.(*ast.ExistsExpr)
	if ex.Query == nil {
		t.Error("ExistsExpr.Query is nil; parser-select must fill it")
	}

	q = parseQ(t, "SELECT ARRAY(SELECT x FROM t2) FROM t1")
	sel = selectOf(t, q)
	arr := sel.Items[0].Expr.(*ast.ArraySubqueryExpr)
	if arr.Query == nil {
		t.Error("ArraySubqueryExpr.Query is nil; parser-select must fill it")
	}
}

// TestSelect_NestedSubqueryFilled confirms a subquery nested inside another
// subquery is also filled (the fill post-pass recurses).
func TestSelect_NestedSubqueryFilled(t *testing.T) {
	q := parseQ(t, "SELECT (SELECT (SELECT 1)) AS m")
	sel := selectOf(t, q)
	outer := sel.Items[0].Expr.(*ast.SubqueryExpr)
	if outer.Query == nil {
		t.Fatal("outer subquery not filled")
	}
	innerSel := outer.Query.(*ast.QueryStmt).Body.(*ast.SelectStmt)
	mid := innerSel.Items[0].Expr.(*ast.SubqueryExpr)
	if mid.Query == nil {
		t.Error("nested subquery not filled (fill pass must recurse)")
	}
}

// TestSelect_CommaJoinPrecedence confirms a top-level comma binds looser than an
// explicit JOIN: `A, B JOIN C` is two FROM items, the second a JoinExpr.
func TestSelect_CommaJoinPrecedence(t *testing.T) {
	q := parseQ(t, "SELECT * FROM A, B JOIN C ON B.x = C.x")
	sel := selectOf(t, q)
	if len(sel.From) != 2 {
		t.Fatalf("FROM items = %d, want 2 (A | B JOIN C)", len(sel.From))
	}
	if _, ok := sel.From[0].(*ast.TableExpr); !ok {
		t.Errorf("FROM[0] = %T, want *ast.TableExpr (A)", sel.From[0])
	}
	je, ok := sel.From[1].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("FROM[1] = %T, want *ast.JoinExpr (B JOIN C)", sel.From[1])
	}
	if je.Type != ast.JoinInner || je.On == nil {
		t.Errorf("join type=%s on=%v, want INNER with ON", je.Type, je.On != nil)
	}
}

// TestSelect_JoinChainLeftDeep confirms a JOIN chain folds left-deep.
func TestSelect_JoinChainLeftDeep(t *testing.T) {
	q := parseQ(t, "SELECT * FROM a JOIN b ON a.x=b.x JOIN c ON b.y=c.y")
	sel := selectOf(t, q)
	if len(sel.From) != 1 {
		t.Fatalf("FROM items = %d, want 1 (single left-deep join tree)", len(sel.From))
	}
	top, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("FROM[0] = %T, want *ast.JoinExpr", sel.From[0])
	}
	if _, ok := top.Left.(*ast.JoinExpr); !ok {
		t.Errorf("top.Left = %T, want nested *ast.JoinExpr (left-deep)", top.Left)
	}
}

// TestSelect_JoinTypes confirms each join_type keyword maps to the right
// JoinType, OUTER is absorbed, and USING/NATURAL/CROSS criteria are correct.
func TestSelect_JoinTypes(t *testing.T) {
	cases := []struct {
		sql      string
		wantType ast.JoinType
		natural  bool
		hasOn    bool
		hasUsing bool
	}{
		{"SELECT * FROM a JOIN b ON a.x=b.x", ast.JoinInner, false, true, false},
		{"SELECT * FROM a INNER JOIN b ON a.x=b.x", ast.JoinInner, false, true, false},
		{"SELECT * FROM a LEFT JOIN b USING (x)", ast.JoinLeft, false, false, true},
		{"SELECT * FROM a LEFT OUTER JOIN b USING (x)", ast.JoinLeft, false, false, true},
		{"SELECT * FROM a RIGHT JOIN b ON a.x=b.x", ast.JoinRight, false, true, false},
		{"SELECT * FROM a FULL OUTER JOIN b ON a.x=b.x", ast.JoinFull, false, true, false},
		{"SELECT * FROM a CROSS JOIN b", ast.JoinCross, false, false, false},
		{"SELECT * FROM a NATURAL JOIN b", ast.JoinInner, true, false, false},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			q := parseQ(t, c.sql)
			sel := selectOf(t, q)
			je, ok := sel.From[0].(*ast.JoinExpr)
			if !ok {
				t.Fatalf("FROM[0] = %T, want *ast.JoinExpr", sel.From[0])
			}
			if je.Type != c.wantType {
				t.Errorf("type = %s, want %s", je.Type, c.wantType)
			}
			if je.Natural != c.natural {
				t.Errorf("natural = %v, want %v", je.Natural, c.natural)
			}
			if (je.On != nil) != c.hasOn {
				t.Errorf("hasOn = %v, want %v", je.On != nil, c.hasOn)
			}
			if (je.Using != nil) != c.hasUsing {
				t.Errorf("hasUsing = %v, want %v", je.Using != nil, c.hasUsing)
			}
		})
	}
}

// TestSelect_GroupByVariants confirms each GROUP BY shape maps to the right kind.
func TestSelect_GroupByVariants(t *testing.T) {
	cases := []struct {
		sql      string
		wantKind ast.GroupByKind
		nItems   int
	}{
		{"SELECT a FROM t GROUP BY a, b", ast.GroupByItems, 2},
		{"SELECT a FROM t GROUP BY ALL", ast.GroupByAll, 0},
		{"SELECT a FROM t GROUP BY ROLLUP (a, b)", ast.GroupByItems, 1},
		{"SELECT a FROM t GROUP BY CUBE (a, b)", ast.GroupByItems, 1},
		{"SELECT a FROM t GROUP BY GROUPING SETS (a, b, ())", ast.GroupByItems, 1},
		{"SELECT a FROM t GROUP BY ()", ast.GroupByItems, 1},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			q := parseQ(t, c.sql)
			sel := selectOf(t, q)
			if sel.GroupBy == nil {
				t.Fatal("GroupBy is nil")
			}
			if sel.GroupBy.Kind != c.wantKind {
				t.Errorf("kind = %v, want %v", sel.GroupBy.Kind, c.wantKind)
			}
			if len(sel.GroupBy.Items) != c.nItems {
				t.Errorf("items = %d, want %d", len(sel.GroupBy.Items), c.nItems)
			}
		})
	}

	// ROLLUP item kind.
	q := parseQ(t, "SELECT a FROM t GROUP BY ROLLUP (a, b)")
	gi := selectOf(t, q).GroupBy.Items[0]
	if gi.Kind != ast.GroupingRollup || len(gi.Items) != 2 {
		t.Errorf("ROLLUP item kind=%v items=%d, want GroupingRollup with 2", gi.Kind, len(gi.Items))
	}
}

// TestSelect_SelectAsModifiers confirms AS STRUCT / AS VALUE / AS <path>.
func TestSelect_SelectAsModifiers(t *testing.T) {
	if as := selectOf(t, parseQ(t, "SELECT AS STRUCT 1 AS a")).As; as != ast.SelectAsStruct {
		t.Errorf("AS STRUCT: kind = %v", as)
	}
	if as := selectOf(t, parseQ(t, "SELECT AS VALUE 1")).As; as != ast.SelectAsValue {
		t.Errorf("AS VALUE: kind = %v", as)
	}
	sel := selectOf(t, parseQ(t, "SELECT AS pkg.MyType 1 AS a"))
	if sel.As != ast.SelectAsTypeName || sel.AsTypeName == nil || sel.AsTypeName.String() != "pkg.MyType" {
		t.Errorf("AS <path>: kind=%v name=%v", sel.As, sel.AsTypeName)
	}
}

// TestSelect_DistinctAll confirms the set quantifier flags.
func TestSelect_DistinctAll(t *testing.T) {
	if !selectOf(t, parseQ(t, "SELECT DISTINCT a FROM t")).Distinct {
		t.Error("DISTINCT not set")
	}
	if !selectOf(t, parseQ(t, "SELECT ALL a FROM t")).All {
		t.Error("ALL not set")
	}
}

// TestSelect_CTEShape confirms the WITH clause structure: recursion flag,
// explicit columns, CTE count, and that the body Query is a real *QueryStmt.
func TestSelect_CTEShape(t *testing.T) {
	q := parseQ(t, "WITH a AS (SELECT 1), b (x, y) AS (SELECT 2, 3) SELECT * FROM b")
	if q.With == nil {
		t.Fatal("With is nil")
	}
	if q.With.Recursive {
		t.Error("Recursive = true, want false")
	}
	if len(q.With.CTEs) != 2 {
		t.Fatalf("CTEs = %d, want 2", len(q.With.CTEs))
	}
	if q.With.CTEs[0].Name != "a" || q.With.CTEs[1].Name != "b" {
		t.Errorf("CTE names = %q, %q", q.With.CTEs[0].Name, q.With.CTEs[1].Name)
	}
	if got := q.With.CTEs[1].Columns; len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Errorf("CTE b columns = %v, want [x y]", got)
	}
	if _, ok := q.With.CTEs[0].Query.(*ast.QueryStmt); !ok {
		t.Errorf("CTE query = %T, want *ast.QueryStmt", q.With.CTEs[0].Query)
	}

	// RECURSIVE flag.
	rq := parseQ(t, "WITH RECURSIVE c AS ((SELECT 1 AS n) UNION ALL (SELECT n+1 FROM c WHERE n<3)) SELECT * FROM c")
	if !rq.With.Recursive {
		t.Error("RECURSIVE not captured")
	}
}

// TestSelect_RecursionDepthModifier confirms the WITH DEPTH modifier on a CTE.
func TestSelect_RecursionDepthModifier(t *testing.T) {
	q := parseQ(t, "WITH RECURSIVE c AS ((SELECT 1 AS n) UNION ALL (SELECT n+1 FROM c WHERE n<5)) WITH DEPTH AS d BETWEEN 1 AND 10 SELECT * FROM c")
	cte := q.With.CTEs[0]
	if cte.Depth == nil {
		t.Fatal("recursion depth modifier not parsed")
	}
	if cte.Depth.Alias != "d" {
		t.Errorf("depth alias = %q, want d", cte.Depth.Alias)
	}
	if cte.Depth.Lower == nil || cte.Depth.Upper == nil {
		t.Errorf("BETWEEN bounds = %v/%v, want both set", cte.Depth.Lower, cte.Depth.Upper)
	}
}

// TestSelect_UnnestSource confirms an UNNEST FROM source with alias + WITH OFFSET.
func TestSelect_UnnestSource(t *testing.T) {
	q := parseQ(t, "SELECT num, pos FROM UNNEST([1,2,3]) AS num WITH OFFSET AS pos")
	sel := selectOf(t, q)
	ue, ok := sel.From[0].(*ast.UnnestExpr)
	if !ok {
		t.Fatalf("FROM[0] = %T, want *ast.UnnestExpr", sel.From[0])
	}
	if ue.Alias != "num" {
		t.Errorf("alias = %q, want num", ue.Alias)
	}
	if !ue.WithOffset || ue.WithOffsetAlias != "pos" {
		t.Errorf("WITH OFFSET = %v alias=%q, want true/pos", ue.WithOffset, ue.WithOffsetAlias)
	}
	if ue.Array == nil {
		t.Error("UNNEST array is nil")
	}
}

// TestSelect_TableSourceShapes confirms path / subquery / TVF FROM sources and
// aliases.
func TestSelect_TableSourceShapes(t *testing.T) {
	// Path with schema/project qualification and alias.
	sel := selectOf(t, parseQ(t, "SELECT * FROM project.dataset.tbl AS x"))
	te := sel.From[0].(*ast.TableExpr)
	if te.Path == nil || len(te.Path.Parts) != 3 || te.Alias != "x" {
		t.Errorf("path source = %v alias=%q, want 3-part path / x", te.Path, te.Alias)
	}
	// Dashed path.
	sel = selectOf(t, parseQ(t, "SELECT * FROM `region-us`.dataset.tbl"))
	_ = sel // backtick-dashed path: just must parse
	sel = selectOf(t, parseQ(t, "SELECT * FROM my-project.ds.tbl"))
	te = sel.From[0].(*ast.TableExpr)
	if te.Path == nil || te.Path.Parts[0] != "my-project" {
		t.Errorf("dashed path first part = %v, want my-project", te.Path)
	}
	// Subquery source.
	sel = selectOf(t, parseQ(t, "SELECT * FROM (SELECT 1 AS n) AS s"))
	te = sel.From[0].(*ast.TableExpr)
	if te.Subquery == nil || te.Alias != "s" {
		t.Errorf("subquery source = %v alias=%q", te.Subquery, te.Alias)
	}
	if _, ok := te.Subquery.(*ast.QueryStmt); !ok {
		t.Errorf("subquery source body = %T, want *ast.QueryStmt", te.Subquery)
	}
	// TVF source.
	sel = selectOf(t, parseQ(t, "SELECT * FROM my_tvf(1, 2) AS t"))
	te = sel.From[0].(*ast.TableExpr)
	if te.Func == nil || te.Alias != "t" {
		t.Errorf("TVF source = %v alias=%q", te.Func, te.Alias)
	}
	// FOR SYSTEM_TIME.
	sel = selectOf(t, parseQ(t, "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'"))
	te = sel.From[0].(*ast.TableExpr)
	if te.SystemTime == nil {
		t.Error("FOR SYSTEM_TIME AS OF not captured")
	}
}

// TestSelect_WindowClause confirms named WINDOW definitions.
func TestSelect_WindowClause(t *testing.T) {
	q := parseQ(t, "SELECT SUM(x) OVER w FROM t WINDOW w AS (PARTITION BY a ORDER BY b), v AS (w)")
	sel := selectOf(t, q)
	if len(sel.Window) != 2 {
		t.Fatalf("WINDOW defs = %d, want 2", len(sel.Window))
	}
	if sel.Window[0].Name != "w" || sel.Window[1].Name != "v" {
		t.Errorf("window names = %q, %q", sel.Window[0].Name, sel.Window[1].Name)
	}
	if sel.Window[0].Spec == nil || len(sel.Window[0].Spec.PartitionBy) != 1 {
		t.Errorf("window w spec partition-by = %v", sel.Window[0].Spec)
	}
	// A base-window-only spec references another window by name.
	if sel.Window[1].Spec == nil || sel.Window[1].Spec.Name != "w" {
		t.Errorf("window v base = %v, want reference to w", sel.Window[1].Spec)
	}
}

// TestSelect_ClauseOrder confirms WHERE→GROUP BY→HAVING→QUALIFY→WINDOW all land
// on the SELECT and ORDER BY/LIMIT land on the QueryStmt.
func TestSelect_ClauseOrder(t *testing.T) {
	q := parseQ(t, "SELECT a, COUNT(*) c FROM t WHERE a > 0 GROUP BY a HAVING c > 1 QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1 WINDOW w AS () ORDER BY a LIMIT 5")
	sel := selectOf(t, q)
	if sel.Where == nil {
		t.Error("WHERE nil")
	}
	if sel.GroupBy == nil {
		t.Error("GROUP BY nil")
	}
	if sel.Having == nil {
		t.Error("HAVING nil")
	}
	if sel.Qualify == nil {
		t.Error("QUALIFY nil")
	}
	if len(sel.Window) != 1 {
		t.Error("WINDOW missing")
	}
	if len(q.OrderBy) != 1 || q.Limit == nil {
		t.Error("ORDER BY / LIMIT must be on the QueryStmt")
	}
}

// TestSelect_ParenthesizedQueryPrimary confirms a parenthesized query is a
// query_primary that can be a set-op operand and carries Parens=true.
func TestSelect_ParenthesizedQueryPrimary(t *testing.T) {
	q := parseQ(t, "(SELECT 1 AS x) UNION ALL (SELECT 2 AS x)")
	so, ok := q.Body.(*ast.SetOperation)
	if !ok {
		t.Fatalf("body = %T, want *ast.SetOperation", q.Body)
	}
	lq, ok := so.Left.(*ast.QueryStmt)
	if !ok || !lq.Parens {
		t.Errorf("left operand = %T (parens=%v), want parenthesized *ast.QueryStmt", so.Left, ok && lq.Parens)
	}
}

// TestSelect_WalkVisitsClauses confirms the generated walker descends into the
// query tree (a Visitor sees the SELECT, FROM table, WHERE expr, etc.). The
// query-span extractor relies on this traversal.
func TestSelect_WalkVisitsClauses(t *testing.T) {
	q := parseQ(t, "SELECT a, b FROM tbl WHERE a > 1 GROUP BY a ORDER BY b")
	var sawSelect, sawTable, sawJoinOrPath bool
	ast.Inspect(q, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SelectStmt:
			sawSelect = true
		case *ast.TableExpr:
			sawTable = true
			if v.Path != nil {
				sawJoinOrPath = true
			}
		}
		return true
	})
	if !sawSelect || !sawTable || !sawJoinOrPath {
		t.Errorf("walk coverage: select=%v table=%v path=%v", sawSelect, sawTable, sawJoinOrPath)
	}
}

// TestSelect_LocSpans confirms nodes carry plausible source spans (Loc.Start <
// Loc.End and within the input) — query-span / Diagnose need accurate positions.
func TestSelect_LocSpans(t *testing.T) {
	sql := "SELECT a FROM t WHERE a > 1"
	q := parseQ(t, sql)
	if !q.Loc.IsValid() || q.Loc.Start != 0 || q.Loc.End != len(sql) {
		t.Errorf("QueryStmt Loc = %+v, want [0,%d)", q.Loc, len(sql))
	}
	sel := selectOf(t, q)
	if !sel.Loc.IsValid() || sel.Loc.Start != 0 {
		t.Errorf("SelectStmt Loc = %+v, want valid starting at 0", sel.Loc)
	}
	// The first select item "a" spans bytes 7..8.
	it := sel.Items[0]
	if it.Loc.Start != 7 {
		t.Errorf("first item Loc.Start = %d, want 7", it.Loc.Start)
	}
}

// TestSelect_DiagnoseNoFalsePositive confirms a valid query draws zero parse
// diagnostics (the bytebase Diagnose contract for the query family).
func TestSelect_DiagnoseNoFalsePositive(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1",
		"WITH c AS (SELECT 1 AS n) SELECT n FROM c",
		"SELECT * FROM a JOIN b USING (id) WHERE a.x > 0 ORDER BY a.y",
	} {
		_, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("Parse(%q): want 0 diagnostics, got %v", sql, errs)
		}
	}
}

// TestSelect_TPCHCorpus is the completeness gate's real-world coverage probe:
// it parses every ZetaSQL TPC-H query (22 complex analytic queries with nested
// subqueries, multi-table joins, comma cross joins, CTEs, aggregates, GROUP BY,
// HAVING, ORDER BY, correlated subqueries, and scalar subqueries) and asserts a
// clean parse — broad coverage beyond the hand-written cases. If the corpus is
// not checked out (CI), the test skips. The corpus is the canonical ZetaSQL
// reference workload the legacy ANTLR grammar was validated against.
func TestSelect_TPCHCorpus(t *testing.T) {
	dir := "/Users/h3n4l/OpenSource/parser/googlesql/examples/zetasql/examples/tpch"
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil || len(files) == 0 {
		t.Skipf("TPC-H corpus not available at %s", dir)
	}
	for _, f := range files {
		f := f
		t.Run(filepath.Base(f), func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("read %s: %v", f, err)
			}
			file, errs := Parse(string(data))
			if len(errs) != 0 {
				t.Errorf("%s: parse errors: %v", filepath.Base(f), errs)
			}
			if len(file.Stmts) == 0 {
				t.Errorf("%s: produced no statements", filepath.Base(f))
			}
		})
	}
}

// TestSelect_ErrorPositions confirms reject diagnostics point at a sensible
// location (non-negative offset within the statement) so Diagnose can underline.
func TestSelect_ErrorPositions(t *testing.T) {
	_, errs := Parse("SELECT * FROM a JOIN") // JOIN with no right source
	if len(errs) == 0 {
		t.Fatal("want a parse error for JOIN with no right source")
	}
	if errs[0].Loc.Start < 0 {
		t.Errorf("error Loc.Start = %d, want >= 0", errs[0].Loc.Start)
	}
	if !strings.Contains(strings.ToLower(errs[0].Msg), "syntax error") {
		t.Errorf("error message = %q, want a syntax error", errs[0].Msg)
	}
}
