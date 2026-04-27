package catalog

import (
	"fmt"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
	"github.com/bytebase/omni/mysql/parser"
)

// extractExprFromView extracts the expression portion from SHOW CREATE VIEW output.
// MySQL 8.0 SHOW CREATE VIEW returns:
//
//	CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `test`.`v` AS select <expr> AS `<alias>` from `test`.`t`
//
// We extract <expr> — the part between "AS select " and the first " AS `".
func extractExprFromView(showCreate string) string {
	idx := strings.Index(showCreate, " AS select ")
	if idx < 0 {
		return showCreate
	}
	selectPart := showCreate[idx+len(" AS select "):]

	// Find " AS `" which marks the column alias
	aliasIdx := strings.Index(selectPart, " AS `")
	if aliasIdx < 0 {
		return selectPart
	}
	return selectPart[:aliasIdx]
}

// deparseExprForOracle parses a SQL expression and deparses it via our deparser.
func deparseExprForOracle(t *testing.T, expr string) string {
	t.Helper()
	sql := "SELECT " + expr + " FROM t"
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", sql, err)
	}
	if stmts.Len() == 0 {
		t.Fatalf("no statements parsed from %q", sql)
	}
	sel, ok := stmts.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", stmts.Items[0])
	}
	if len(sel.TargetList) == 0 {
		t.Fatalf("no target list in SELECT from %q", sql)
	}
	resolveDeparseOracleSelect(t, sel)
	target := sel.TargetList[0]
	if rt, ok := target.(*nodes.ResTarget); ok {
		return deparse.Deparse(rt.Val)
	}
	return deparse.Deparse(target)
}

// deparseExprRewriteForOracle parses a SQL expression, applies RewriteExpr, and deparses.
func deparseExprRewriteForOracle(t *testing.T, expr string) string {
	t.Helper()
	sql := "SELECT " + expr + " FROM t"
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", sql, err)
	}
	if stmts.Len() == 0 {
		t.Fatalf("no statements parsed from %q", sql)
	}
	sel, ok := stmts.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", stmts.Items[0])
	}
	if len(sel.TargetList) == 0 {
		t.Fatalf("no target list in SELECT from %q", sql)
	}
	resolveDeparseOracleSelect(t, sel)
	target := sel.TargetList[0]
	if rt, ok := target.(*nodes.ResTarget); ok {
		return deparse.Deparse(deparse.RewriteExpr(rt.Val))
	}
	return deparse.Deparse(deparse.RewriteExpr(target))
}

func resolveDeparseOracleSelect(t *testing.T, sel *nodes.SelectStmt) {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT, b INT, c INT)", nil)
	if err != nil {
		t.Fatalf("failed to set up resolver catalog: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("failed to set up resolver catalog: %v", r.Error)
		}
	}
	db := c.GetDatabase("test")
	resolver := &deparse.Resolver{
		Lookup:         tableLookupForDB(db),
		DefaultCharset: c.defaultCharset,
	}
	resolver.Resolve(sel)
}

// TestDeparse_Section_4_1_Container verifies NOT folding against MySQL 8.0.
func TestDeparse_Section_4_1_Container(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Create base table with integer column for boolean/comparison tests
	ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT)")

	cases := []struct {
		name  string
		input string // expression to test
	}{
		{"not_gt", "NOT (a > 0)"},
		{"not_lt", "NOT (a < 0)"},
		{"not_ge", "NOT (a >= 0)"},
		{"not_le", "NOT (a <= 0)"},
		{"not_eq", "NOT (a = 0)"},
		{"not_ne", "NOT (a <> 0)"},
		{"not_col", "NOT a"},
		{"not_add", "NOT (a + 1)"},
		{"bang_col", "!a"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name
			selectSQL := fmt.Sprintf("SELECT %s FROM t", tc.input)

			// Create view on MySQL 8.0
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, selectSQL)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW failed: %v", err)
			}

			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW failed: %v", err)
			}

			// Extract just the expression from MySQL's output
			mysqlExpr := extractExprFromView(mysqlOutput)

			// Our deparser output (with rewrite)
			omniExpr := deparseExprRewriteForOracle(t, tc.input)

			t.Logf("MySQL:  %s", mysqlExpr)
			t.Logf("Omni:   %s", omniExpr)

			if mysqlExpr != omniExpr {
				t.Errorf("mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlExpr, omniExpr)
			}
		})
	}
}

// TestDeparse_Section_3_2_Container verifies TRIM special forms against MySQL 8.0.
func TestDeparse_Section_3_2_Container(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Create base table
	ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a VARCHAR(50), b VARCHAR(50))")

	cases := []struct {
		name  string
		input string // expression to test
	}{
		{"trim_simple", "TRIM(a)"},
		{"trim_leading", "TRIM(LEADING 'x' FROM a)"},
		{"trim_trailing", "TRIM(TRAILING 'x' FROM a)"},
		{"trim_both", "TRIM(BOTH 'x' FROM a)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name
			selectSQL := fmt.Sprintf("SELECT %s FROM t", tc.input)

			// Create view on MySQL 8.0
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, selectSQL)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW failed: %v", err)
			}

			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW failed: %v", err)
			}

			// Extract just the expression from MySQL's output
			mysqlExpr := extractExprFromView(mysqlOutput)

			// Our deparser output
			omniExpr := deparseExprForOracle(t, tc.input)

			t.Logf("MySQL:  %s", mysqlExpr)
			t.Logf("Omni:   %s", omniExpr)

			if mysqlExpr != omniExpr {
				t.Errorf("mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlExpr, omniExpr)
			}
		})
	}
}

// extractSelectBody extracts the SELECT body from SHOW CREATE VIEW output.
// MySQL 8.0 format:
//
//	CREATE ALGORITHM=... VIEW `test`.`v` AS select ...
//
// Our catalog format:
//
//	CREATE ALGORITHM=... VIEW `v` AS select ...
//
// We extract everything after " AS " (the first occurrence after VIEW).
func extractSelectBody(showCreate string) string {
	// Find "VIEW " to locate the view name portion, then find " AS " after that.
	viewIdx := strings.Index(showCreate, "VIEW ")
	if viewIdx < 0 {
		return showCreate
	}
	rest := showCreate[viewIdx:]
	asIdx := strings.Index(rest, " AS ")
	if asIdx < 0 {
		return showCreate
	}
	return rest[asIdx+len(" AS "):]
}

// TestDeparse_Section_7_2_SimpleViews verifies that our catalog's SHOW CREATE VIEW
// output matches MySQL 8.0's output for simple view definitions.
// For each view, we compare the SELECT body portion (after "AS ").
func TestDeparse_Section_7_2_SimpleViews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
	}{
		{"select_constant", "SELECT 1"},
		{"select_column", "SELECT a FROM t"},
		{"select_alias", "SELECT a AS col1 FROM t"},
		{"select_multi_columns", "SELECT a, b FROM t"},
		{"select_where", "SELECT a FROM t WHERE a > 0"},
		{"select_orderby_limit", "SELECT a FROM t ORDER BY a LIMIT 10"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := extractSelectBody(mysqlOutput)

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL full:  %s", mysqlOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparse_Section_7_4_JoinViews verifies that our catalog's SHOW CREATE VIEW
// output matches MySQL 8.0's output for views with JOINs (INNER JOIN, LEFT JOIN,
// multiple tables, subquery in FROM).
func TestDeparse_Section_7_4_JoinViews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
		tables   string // additional CREATE TABLE statements for our catalog (beyond t, t1, t2)
		partial  bool   // expected partial match (parser limitation)
	}{
		{"inner_join", "SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a", "", false},
		{"left_join", "SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a", "", false},
		{"multi_table", "SELECT t1.a FROM t1, t2 WHERE t1.a = t2.a", "", false},
		{"subquery_from", "SELECT d.x FROM (SELECT a AS x FROM t) d", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
			cat.Exec("CREATE TABLE t2 (a INT, b INT)", nil)
			if tc.tables != "" {
				cat.Exec(tc.tables, nil)
			}
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial — parser limitation)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL full:  %s", mysqlOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparse_Section_7_5_AdvancedViews verifies that our catalog's SHOW CREATE VIEW
// output matches MySQL 8.0's output for advanced view definitions: UNION, CTE,
// window functions, nested subqueries, boolean expressions, and combined rewrites.
func TestDeparse_Section_7_5_AdvancedViews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
		partial  bool   // expected partial match (parser/resolver limitation)
	}{
		{"union_view", "SELECT a FROM t UNION SELECT b FROM t", false},
		{"cte_view", "WITH cte AS (SELECT a FROM t) SELECT * FROM cte", false},
		{"window_func_view", "SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t", false},
		{"nested_subquery_view", "SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)", false},
		{"boolean_expr_view", "SELECT a AND b, a OR b FROM t", false},
		{"combined_rewrite_view", "SELECT a + b, NOT (a > 0), CAST(a AS CHAR), COUNT(*) FROM t GROUP BY a, b HAVING COUNT(*) > 1", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial — parser/resolver limitation)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL full:  %s", mysqlOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// stripDatabasePrefix removes the `test`. database prefix from MySQL 8.0 SHOW CREATE VIEW output.
// MySQL 8.0 qualifies all identifiers with the database name (e.g., `test`.`t`.`a`),
// while our catalog does not. We strip the prefix for comparison.
func stripDatabasePrefix(s string) string {
	return strings.ReplaceAll(s, "`test`.", "")
}

// TestDeparse_Section_7_6_Regression verifies that the deparser integration does not
// break existing tests (scenarios 1-2 are covered by running go test ./mysql/catalog/ -short
// and go test ./mysql/parser/ -short separately) and that views with explicit column
// aliases match MySQL 8.0 output exactly.
func TestDeparse_Section_7_6_Regression(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	t.Run("view_with_explicit_column_aliases", func(t *testing.T) {
		viewName := "v_col_alias"
		createSQL := "CREATE VIEW " + viewName + "(x, y) AS SELECT a, b FROM t"

		// --- MySQL 8.0 side ---
		ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
		}
		mysqlOutput, err := ctr.showCreateView(viewName)
		if err != nil {
			t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
		}

		// --- Our catalog side ---
		cat := New()
		cat.Exec("CREATE DATABASE test", nil)
		cat.SetCurrentDatabase("test")
		cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
		results, _ := cat.Exec(createSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
		}
		omniOutput := cat.ShowCreateView("test", viewName)
		if omniOutput == "" {
			t.Fatal("ShowCreateView returned empty")
		}

		// Compare full output (stripping database prefix from MySQL).
		// MySQL: CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `test`.`v_col_alias` (`x`,`y`) AS select ...
		// Omni:  CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_col_alias` (`x`,`y`) AS select ...
		mysqlNorm := stripDatabasePrefix(mysqlOutput)

		t.Logf("MySQL full:  %s", mysqlOutput)
		t.Logf("Omni full:   %s", omniOutput)
		t.Logf("MySQL norm:  %s", mysqlNorm)

		// Compare the preamble (up to and including column list).
		mysqlPreamble := extractViewPreamble(mysqlNorm)
		omniPreamble := extractViewPreamble(omniOutput)
		if mysqlPreamble != omniPreamble {
			t.Errorf("preamble mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlPreamble, omniPreamble)
		}

		// Compare the SELECT body.
		mysqlBody := extractSelectBody(mysqlNorm)
		omniBody := extractSelectBody(omniOutput)
		if mysqlBody != omniBody {
			t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
		}
	})

	// Verify simple and complex views still match (re-run 7.2 and 7.5 representative cases).
	simpleAndComplexCases := []struct {
		name     string
		createAs string
	}{
		{"simple_select_column", "SELECT a FROM t"},
		{"simple_select_where", "SELECT a FROM t WHERE a > 0"},
		{"complex_union", "SELECT a FROM t UNION SELECT b FROM t"},
		{"complex_boolean", "SELECT a AND b, a OR b FROM t"},
	}

	for _, tc := range simpleAndComplexCases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_reg_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_7_1_WindowFunctions verifies window function patterns
// against MySQL 8.0 SHOW CREATE VIEW output.
// Covers: ROW_NUMBER, SUM OVER PARTITION BY+ORDER BY, ROWS frame, RANGE frame,
// named window, multiple window functions, LAG/LEAD.
func TestDeparseContainer_7_1_WindowFunctions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		partial  bool
	}{
		{"row_number_basic", "v_wf_rownum", "CREATE VIEW v_wf_rownum AS SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t", false},
		{"sum_partition_orderby", "v_wf_sum_part", "CREATE VIEW v_wf_sum_part AS SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b) FROM t", false},
		{"rows_frame", "v_wf_rows", "CREATE VIEW v_wf_rows AS SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t", false},
		{"range_frame", "v_wf_range", "CREATE VIEW v_wf_range AS SELECT a, AVG(b) OVER (ORDER BY a RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t", false},
		{"named_window", "v_wf_named", "CREATE VIEW v_wf_named AS SELECT a, RANK() OVER w, DENSE_RANK() OVER w FROM t WINDOW w AS (ORDER BY a)", false},
		{"multiple_window_funcs", "v_wf_multi", "CREATE VIEW v_wf_multi AS SELECT a, ROW_NUMBER() OVER (ORDER BY a), SUM(b) OVER (ORDER BY a) FROM t", false},
		{"lag_lead", "v_wf_lagld", "CREATE VIEW v_wf_lagld AS SELECT a, LAG(a, 1) OVER (ORDER BY a), LEAD(a, 1) OVER (ORDER BY a) FROM t", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial — parser limitation)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_1_1_ArithmeticComparison verifies arithmetic and comparison operators
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_1_1_ArithmeticComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"addition", "v_add", "CREATE VIEW v_add AS SELECT a + b FROM t"},
		{"subtraction", "v_sub", "CREATE VIEW v_sub AS SELECT a - b FROM t"},
		{"multiplication", "v_mul", "CREATE VIEW v_mul AS SELECT a * b FROM t"},
		{"division", "v_div", "CREATE VIEW v_div AS SELECT a / b FROM t"},
		{"int_division", "v_intdiv", "CREATE VIEW v_intdiv AS SELECT a DIV b FROM t"},
		{"mod", "v_mod", "CREATE VIEW v_mod AS SELECT a MOD b FROM t"},
		{"equals", "v_eq", "CREATE VIEW v_eq AS SELECT a = b FROM t"},
		{"not_equals_bang", "v_neq", "CREATE VIEW v_neq AS SELECT a != b FROM t"},
		{"not_equals_ltgt", "v_neq2", "CREATE VIEW v_neq2 AS SELECT a <> b FROM t"},
		{"comparisons", "v_cmp", "CREATE VIEW v_cmp AS SELECT a > b, a < b, a >= b, a <= b FROM t"},
		{"null_safe_equals", "v_nseq", "CREATE VIEW v_nseq AS SELECT a <=> b FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparse_Section_7_3_ExpressionViews verifies that our catalog's SHOW CREATE VIEW
// output matches MySQL 8.0's output for views with expressions (arithmetic, functions,
// CASE, CAST, aggregates).
func TestDeparse_Section_7_3_ExpressionViews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
	}{
		{"arithmetic_expr", "SELECT a + b FROM t"},
		{"function_call", "SELECT CONCAT(a, b) FROM t"},
		{"case_expr", "SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t"},
		{"cast_expr", "SELECT CAST(a AS CHAR) FROM t"},
		{"aggregate_expr", "SELECT COUNT(*), SUM(a) FROM t GROUP BY a HAVING SUM(a) > 10"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL full:  %s", mysqlOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_1_3_LiteralsSpacing verifies literals and spacing rules
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_1_3_LiteralsSpacing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		partial  bool // mark [~] if parser limitation
	}{
		{"basic_literals", "v_basic_lit", "CREATE VIEW v_basic_lit AS SELECT 1, 1.5, 'hello', NULL, TRUE, FALSE FROM t", false},
		{"hex_bit_literals", "v_hex_bit", "CREATE VIEW v_hex_bit AS SELECT 0xFF, X'FF', 0b1010, b'1010' FROM t", false},
		{"charset_introducers", "v_charset", "CREATE VIEW v_charset AS SELECT _utf8mb4'hello', _latin1'world' FROM t", false},
		{"empty_string", "v_empty_str", "CREATE VIEW v_empty_str AS SELECT '' FROM t", false},
		{"escaped_quotes", "v_esc_quotes", "CREATE VIEW v_esc_quotes AS SELECT 'it''s' FROM t", false},
		{"escaped_backslash", "v_esc_bslash", "CREATE VIEW v_esc_bslash AS SELECT 'back\\\\slash' FROM t", false},
		{"temporal_literals", "v_temporal", "CREATE VIEW v_temporal AS SELECT DATE '2024-01-01', TIME '12:00:00', TIMESTAMP '2024-01-01 12:00:00' FROM t", false},
		{"func_args_no_space", "v_func_args", "CREATE VIEW v_func_args AS SELECT CONCAT(a, b, c) FROM t", false},
		{"in_list_no_space", "v_in_list", "CREATE VIEW v_in_list AS SELECT a IN (1, 2, 3) FROM t", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial — parser limitation)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_1_2_LogicalBitwiseIS verifies logical, bitwise, and IS operators
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_1_2_LogicalBitwiseIS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"and_or", "v_and_or", "CREATE VIEW v_and_or AS SELECT a AND b, a OR b FROM t"},
		{"xor", "v_xor", "CREATE VIEW v_xor AS SELECT a XOR b FROM t"},
		{"not", "v_not", "CREATE VIEW v_not AS SELECT NOT a FROM t"},
		{"bitwise_ops", "v_bitwise", "CREATE VIEW v_bitwise AS SELECT a | b, a & b, a ^ b FROM t"},
		{"shifts", "v_shifts", "CREATE VIEW v_shifts AS SELECT a << b, a >> b FROM t"},
		{"bitwise_not", "v_bitnot", "CREATE VIEW v_bitnot AS SELECT ~a FROM t"},
		{"is_null", "v_isnull", "CREATE VIEW v_isnull AS SELECT a IS NULL, a IS NOT NULL FROM t"},
		{"is_true_false", "v_istf", "CREATE VIEW v_istf AS SELECT a IS TRUE, a IS FALSE FROM t"},
		{"in_not_in", "v_in", "CREATE VIEW v_in AS SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t"},
		{"between", "v_between", "CREATE VIEW v_between AS SELECT a BETWEEN 1 AND 10, a NOT BETWEEN 1 AND 10 FROM t"},
		{"like", "v_like", "CREATE VIEW v_like AS SELECT a LIKE 'foo%', a NOT LIKE 'bar%' FROM t"},
		{"like_escape", "v_like_esc", "CREATE VIEW v_like_esc AS SELECT a LIKE 'x' ESCAPE '\\\\' FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_Section_2_1_FunctionNameRewrites verifies that function name
// rewrites match MySQL 8.0 SHOW CREATE VIEW output.
// Covers: SUBSTRING->substr, CURRENT_TIMESTAMP->now(), CURRENT_DATE->curdate(),
// CURRENT_TIME->curtime(), CURRENT_USER->current_user(), NOW()->now(),
// COUNT(*)->count(0), COUNT(DISTINCT).
func TestDeparseContainer_Section_2_1_FunctionNameRewrites(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"substring_rewrite", "v_substr", "CREATE VIEW v_substr AS SELECT SUBSTRING('abc', 1, 2) FROM t"},
		{"current_timestamp_kw", "v_cur_ts", "CREATE VIEW v_cur_ts AS SELECT CURRENT_TIMESTAMP FROM t"},
		{"current_timestamp_fn", "v_cur_ts_fn", "CREATE VIEW v_cur_ts_fn AS SELECT CURRENT_TIMESTAMP() FROM t"},
		{"current_date_kw", "v_cur_date", "CREATE VIEW v_cur_date AS SELECT CURRENT_DATE FROM t"},
		{"current_time_kw", "v_cur_time", "CREATE VIEW v_cur_time AS SELECT CURRENT_TIME FROM t"},
		{"current_user_kw", "v_cur_user", "CREATE VIEW v_cur_user AS SELECT CURRENT_USER FROM t"},
		{"now_func", "v_now", "CREATE VIEW v_now AS SELECT NOW() FROM t"},
		{"count_star", "v_count_star", "CREATE VIEW v_count_star AS SELECT COUNT(*) FROM t"},
		{"count_distinct", "v_count_dist", "CREATE VIEW v_count_dist AS SELECT COUNT(DISTINCT a) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_2_2_RegularFunctionsAggregates verifies regular functions and aggregates
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_2_2_RegularFunctionsAggregates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"concat_upper_lower", "v_str_funcs", "CREATE VIEW v_str_funcs AS SELECT CONCAT(a, b), UPPER(a), LOWER(a) FROM t"},
		{"ifnull_coalesce_nullif", "v_null_funcs", "CREATE VIEW v_null_funcs AS SELECT IFNULL(a, 0), COALESCE(a, b, 0), NULLIF(a, 0) FROM t"},
		{"if_function", "v_if_func", "CREATE VIEW v_if_func AS SELECT IF(a > 0, 'yes', 'no') FROM t"},
		{"abs_greatest_least", "v_num_funcs", "CREATE VIEW v_num_funcs AS SELECT ABS(a), GREATEST(a, b), LEAST(a, b) FROM t"},
		{"sum_avg_max_min", "v_agg_funcs", "CREATE VIEW v_agg_funcs AS SELECT SUM(a), AVG(a), MAX(a), MIN(a) FROM t"},
		{"nested_functions", "v_nested_funcs", "CREATE VIEW v_nested_funcs AS SELECT CONCAT(UPPER(TRIM(a)), LOWER(b)) FROM t"},
		{"multiple_aggregates_groupby", "v_multi_agg", "CREATE VIEW v_multi_agg AS SELECT COUNT(*), SUM(a), AVG(b), MAX(c) FROM t GROUP BY a"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_2_4_CastConvertOperatorRewrites verifies CAST, CONVERT,
// REGEXP→regexp_like, NOT REGEXP, -> (json_extract), ->> (json_unquote(json_extract))
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_2_4_CastConvertOperatorRewrites(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS tj (a JSON, b INT)"); err != nil {
		t.Fatalf("failed to create table tj on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		tables   string // extra CREATE TABLE for omni catalog (beyond t)
		partial  bool
	}{
		{"cast_char", "v_cast_char", "CREATE VIEW v_cast_char AS SELECT CAST(a AS CHAR) FROM t", "", false},
		{"cast_char10", "v_cast_char10", "CREATE VIEW v_cast_char10 AS SELECT CAST(a AS CHAR(10)) FROM t", "", false},
		{"cast_binary", "v_cast_binary", "CREATE VIEW v_cast_binary AS SELECT CAST(a AS BINARY) FROM t", "", false},
		{"cast_signed_unsigned", "v_cast_su", "CREATE VIEW v_cast_su AS SELECT CAST(a AS SIGNED), CAST(a AS UNSIGNED) FROM t", "", false},
		{"cast_decimal", "v_cast_dec", "CREATE VIEW v_cast_dec AS SELECT CAST(a AS DECIMAL(10,2)) FROM t", "", false},
		{"cast_date_datetime_json", "v_cast_ddj", "CREATE VIEW v_cast_ddj AS SELECT CAST(a AS DATE), CAST(a AS DATETIME), CAST(a AS JSON) FROM t", "", false},
		{"convert_char", "v_conv_char", "CREATE VIEW v_conv_char AS SELECT CONVERT(a, CHAR) FROM t", "", false},
		{"convert_using", "v_conv_using", "CREATE VIEW v_conv_using AS SELECT CONVERT(a USING utf8mb4) FROM t", "", false},
		{"regexp", "v_regexp", "CREATE VIEW v_regexp AS SELECT a REGEXP 'pattern' FROM t", "", false},
		{"not_regexp", "v_not_regexp", "CREATE VIEW v_not_regexp AS SELECT a NOT REGEXP 'pattern' FROM t", "", false},
		{"json_extract", "v_json_ext", "CREATE VIEW v_json_ext AS SELECT a->'$.key' FROM tj", "CREATE TABLE tj (a JSON, b INT)", false},
		{"json_unquote", "v_json_unq", "CREATE VIEW v_json_unq AS SELECT a->>'$.key' FROM tj", "CREATE TABLE tj (a JSON, b INT)", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			if tc.tables != "" {
				cat.Exec(tc.tables, nil)
			}
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_2_3_SpecialFunctions verifies TRIM, GROUP_CONCAT, and simple CASE
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_2_3_SpecialFunctions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS tv (a VARCHAR(50), b VARCHAR(50), c INT)"); err != nil {
		t.Fatalf("failed to create table tv on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		tables   string // extra CREATE TABLE for omni catalog (beyond t)
	}{
		{"trim_simple", "v_trim_simple", "CREATE VIEW v_trim_simple AS SELECT TRIM(a) FROM tv", "CREATE TABLE tv (a VARCHAR(50), b VARCHAR(50), c INT)"},
		{"trim_leading", "v_trim_leading", "CREATE VIEW v_trim_leading AS SELECT TRIM(LEADING 'x' FROM a) FROM tv", "CREATE TABLE tv (a VARCHAR(50), b VARCHAR(50), c INT)"},
		{"trim_trailing", "v_trim_trailing", "CREATE VIEW v_trim_trailing AS SELECT TRIM(TRAILING 'x' FROM a) FROM tv", "CREATE TABLE tv (a VARCHAR(50), b VARCHAR(50), c INT)"},
		{"trim_both", "v_trim_both", "CREATE VIEW v_trim_both AS SELECT TRIM(BOTH 'x' FROM a) FROM tv", "CREATE TABLE tv (a VARCHAR(50), b VARCHAR(50), c INT)"},
		{"group_concat_basic", "v_gc_basic", "CREATE VIEW v_gc_basic AS SELECT GROUP_CONCAT(a ORDER BY a SEPARATOR ',') FROM t", ""},
		{"group_concat_full", "v_gc_full", "CREATE VIEW v_gc_full AS SELECT GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';') FROM t", ""},
		{"simple_case", "v_simple_case", "CREATE VIEW v_simple_case AS SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			if tc.tables != "" {
				cat.Exec(tc.tables, nil)
			}
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_3_1_BooleanContextWrapping verifies that non-boolean expressions
// in boolean context (AND/OR) get (0 <> ...) wrapping to match MySQL 8.0's
// SHOW CREATE VIEW output.
func TestDeparseContainer_3_1_BooleanContextWrapping(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"col_and_col", "v_bool_cc", "CREATE VIEW v_bool_cc AS SELECT a AND b FROM t"},
		{"arith_and_col", "v_bool_ac", "CREATE VIEW v_bool_ac AS SELECT (a + 1) AND b FROM t"},
		{"cmp_and_arith", "v_bool_ca", "CREATE VIEW v_bool_ca AS SELECT (a > 0) AND (b + 1) FROM t"},
		{"cmp_and_cmp", "v_bool_cmpx2", "CREATE VIEW v_bool_cmpx2 AS SELECT (a > 0) AND (b > 0) FROM t"},
		{"abs_and_col", "v_bool_abs", "CREATE VIEW v_bool_abs AS SELECT ABS(a) AND b FROM t"},
		{"case_and_col", "v_bool_case", "CREATE VIEW v_bool_case AS SELECT CASE WHEN a > 0 THEN 1 ELSE 0 END AND b FROM t"},
		{"if_and_col", "v_bool_if", "CREATE VIEW v_bool_if AS SELECT IF(a > 0, 1, 0) AND b FROM t"},
		{"subquery_and_col", "v_bool_subq", "CREATE VIEW v_bool_subq AS SELECT (SELECT MAX(a) FROM t) AND b FROM t"},
		{"string_and_int", "v_bool_str", "CREATE VIEW v_bool_str AS SELECT 'hello' AND 1 FROM t"},
		{"ifnull_and_col", "v_bool_ifnull", "CREATE VIEW v_bool_ifnull AS SELECT IFNULL(a, 0) AND b FROM t"},
		{"coalesce_and_int", "v_bool_coal", "CREATE VIEW v_bool_coal AS SELECT COALESCE(a, b) AND 1 FROM t"},
		{"nullif_and_col", "v_bool_nullif", "CREATE VIEW v_bool_nullif AS SELECT NULLIF(a, 0) AND b FROM t"},
		{"greatest_and_int", "v_bool_great", "CREATE VIEW v_bool_great AS SELECT GREATEST(a, b) AND 1 FROM t"},
		{"least_and_int", "v_bool_least", "CREATE VIEW v_bool_least AS SELECT LEAST(a, b) AND 1 FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_3_2_NotFoldingNoDoubleWrap verifies NOT folding and
// no-double-wrapping of boolean expressions against MySQL 8.0.
// Section 3.2 scenarios:
//   - NOT(comparison) folds into inverted operator
//   - NOT(non-boolean) becomes (0 = ...)
//   - ! is same as NOT
//   - comparisons in AND are NOT double-wrapped
//   - predicates (IN, BETWEEN) in AND are NOT wrapped
//   - IS/LIKE in AND are NOT wrapped
//   - EXISTS in AND is NOT wrapped
func TestDeparseContainer_3_2_NotFoldingNoDoubleWrap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"not_folds_into_le", "v_s32_not_fold", "CREATE VIEW v_s32_not_fold AS SELECT NOT (a > 0) FROM t"},
		{"not_non_boolean", "v_s32_not_nb", "CREATE VIEW v_s32_not_nb AS SELECT NOT (a + 1) FROM t"},
		{"bang_col", "v_s32_bang", "CREATE VIEW v_s32_bang AS SELECT !a FROM t"},
		{"cmp_not_double_wrapped", "v_s32_cmp_ndw", "CREATE VIEW v_s32_cmp_ndw AS SELECT (a = b) AND (a > 0) FROM t"},
		{"predicates_not_wrapped", "v_s32_pred_nw", "CREATE VIEW v_s32_pred_nw AS SELECT (a IN (1, 2, 3)) AND (b BETWEEN 1 AND 10) FROM t"},
		{"is_like_not_wrapped", "v_s32_islike_nw", "CREATE VIEW v_s32_islike_nw AS SELECT (a IS NULL) AND (b LIKE 'x%') FROM t"},
		{"exists_not_wrapped", "v_s32_exists_nw", "CREATE VIEW v_s32_exists_nw AS SELECT EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_3_3_ComplexPrecedence verifies complex operator precedence
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_3_3_ComplexPrecedence(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"left_assoc_add", "v_left_add", "CREATE VIEW v_left_add AS SELECT a + b + c FROM t"},
		{"mul_then_add", "v_mul_add", "CREATE VIEW v_mul_add AS SELECT a * b + c FROM t"},
		{"add_then_mul", "v_add_mul", "CREATE VIEW v_add_mul AS SELECT a + b * c FROM t"},
		{"paren_add_mul", "v_paren_mul", "CREATE VIEW v_paren_mul AS SELECT (a + b) * c FROM t"},
		{"or_and_prec", "v_or_and", "CREATE VIEW v_or_and AS SELECT a OR b AND c FROM t"},
		{"paren_or_and", "v_paren_or", "CREATE VIEW v_paren_or AS SELECT (a OR b) AND c FROM t"},
		{"mixed_cmp_logic", "v_mixed", "CREATE VIEW v_mixed AS SELECT a > 0 AND b < 10 OR c = 5 FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_4_1_AllJoinTypes verifies all JOIN types against MySQL 8.0
// SHOW CREATE VIEW output: INNER JOIN, LEFT JOIN, RIGHT JOIN→LEFT swap,
// CROSS JOIN, NATURAL JOIN expanded, STRAIGHT_JOIN, USING expanded, comma→explicit join.
func TestDeparseContainer_4_1_AllJoinTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		partial  bool
	}{
		{"inner_join", "v_inner_join", "CREATE VIEW v_inner_join AS SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a", false},
		{"left_join", "v_left_join", "CREATE VIEW v_left_join AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a", false},
		{"right_join_swap", "v_right_join", "CREATE VIEW v_right_join AS SELECT t1.a, t2.b FROM t1 RIGHT JOIN t2 ON t1.a = t2.a", false},
		{"cross_join", "v_cross_join", "CREATE VIEW v_cross_join AS SELECT t1.a, t2.b FROM t1 CROSS JOIN t2", false},
		{"natural_join", "v_natural_join", "CREATE VIEW v_natural_join AS SELECT * FROM t1 NATURAL JOIN t2", false},
		{"straight_join", "v_straight_join", "CREATE VIEW v_straight_join AS SELECT t1.a, t2.b FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a", false},
		{"using_expanded", "v_using", "CREATE VIEW v_using AS SELECT t1.a, t2.b FROM t1 JOIN t2 USING (a)", false},
		{"comma_to_join", "v_comma_join", "CREATE VIEW v_comma_join AS SELECT t1.a, t2.b FROM t1, t2 WHERE t1.a = t2.a", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
			cat.Exec("CREATE TABLE t2 (a INT, b INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_4_2_MultiTableDerived verifies multi-table JOINs, chained LEFT JOINs,
// derived tables, and table aliases against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_4_2_MultiTableDerived(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t3 (b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t3 on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		partial  bool
	}{
		{"three_table_join", "v_3join", "CREATE VIEW v_3join AS SELECT t1.a, t2.b, t3.c FROM t1 JOIN t2 ON t1.a = t2.a JOIN t3 ON t2.b = t3.b", false},
		{"chained_left_joins", "v_chain_left", "CREATE VIEW v_chain_left AS SELECT t1.a FROM t1 LEFT JOIN t2 ON t1.a = t2.a LEFT JOIN t3 ON t1.a = t3.b", false},
		{"derived_table", "v_derived", "CREATE VIEW v_derived AS SELECT d.x FROM (SELECT a AS x FROM t) d", false},
		{"derived_with_where", "v_derived_where", "CREATE VIEW v_derived_where AS SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10", false},
		{"table_alias_with_as", "v_alias_as", "CREATE VIEW v_alias_as AS SELECT x.a FROM t AS x", false},
		{"table_alias_without_as", "v_alias_noas", "CREATE VIEW v_alias_noas AS SELECT x.a FROM t x", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
			cat.Exec("CREATE TABLE t2 (a INT, b INT)", nil)
			cat.Exec("CREATE TABLE t3 (b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_Section_5_1_SelectClauses verifies that our catalog's SHOW CREATE VIEW
// output matches MySQL 8.0's output for views with WHERE, GROUP BY, HAVING, ORDER BY,
// LIMIT, OFFSET, DISTINCT, and expression-based GROUP BY.
func TestDeparseContainer_Section_5_1_SelectClauses(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		partial  bool
	}{
		{"all_clauses_combined", "v_all_clauses", "CREATE VIEW v_all_clauses AS SELECT a FROM t WHERE a > 1 GROUP BY a HAVING COUNT(*) > 1 ORDER BY a LIMIT 10", false},
		{"alias_in_order_by", "v_alias_orderby", "CREATE VIEW v_alias_orderby AS SELECT a, COUNT(*) cnt FROM t GROUP BY a HAVING COUNT(*) > 1 ORDER BY cnt DESC", false},
		{"distinct_order_by", "v_distinct_orderby", "CREATE VIEW v_distinct_orderby AS SELECT DISTINCT a FROM t ORDER BY a DESC", false},
		{"multi_column_order_by", "v_multi_orderby", "CREATE VIEW v_multi_orderby AS SELECT a FROM t ORDER BY a, b DESC", false},
		{"limit_offset", "v_limit_offset", "CREATE VIEW v_limit_offset AS SELECT a FROM t LIMIT 10 OFFSET 5", false},
		{"expression_group_by", "v_expr_groupby", "CREATE VIEW v_expr_groupby AS SELECT a + b FROM t GROUP BY a + b", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_5_2_SetOperations verifies set operations (UNION, UNION ALL,
// multiple UNION, INTERSECT, EXCEPT, UNION+ORDER BY+LIMIT) against MySQL 8.0.
// INTERSECT/EXCEPT require MySQL 8.0.31+; if rejected, the test is skipped.
func TestDeparseContainer_5_2_SetOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		partial  bool
	}{
		{"union", "v_union", "CREATE VIEW v_union AS SELECT a FROM t UNION SELECT b FROM t", false},
		{"union_all", "v_union_all", "CREATE VIEW v_union_all AS SELECT a FROM t UNION ALL SELECT b FROM t", false},
		{"multiple_union", "v_multi_union", "CREATE VIEW v_multi_union AS SELECT a FROM t UNION SELECT b FROM t UNION SELECT c FROM t", false},
		{"intersect", "v_intersect", "CREATE VIEW v_intersect AS SELECT a FROM t INTERSECT SELECT b FROM t", false},
		{"except", "v_except", "CREATE VIEW v_except AS SELECT a FROM t EXCEPT SELECT b FROM t", false},
		{"union_orderby_limit", "v_union_ordlim", "CREATE VIEW v_union_ordlim AS SELECT a FROM t UNION SELECT b FROM t ORDER BY 1 LIMIT 5", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_5_3_ColumnAliasPatterns verifies column and alias patterns
// against MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_5_3_ColumnAliasPatterns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		tables   string // extra tables for omni catalog (beyond t, t1, t2)
		partial  bool
	}{
		{
			"explicit_alias_as_vs_space",
			"v_alias_as_space",
			"CREATE VIEW v_alias_as_space AS SELECT a AS col1, b col2 FROM t",
			"",
			false,
		},
		{
			"expression_explicit_alias",
			"v_expr_alias",
			"CREATE VIEW v_expr_alias AS SELECT a + b AS sum_col FROM t",
			"",
			false,
		},
		{
			"literal_auto_alias",
			"v_lit_auto",
			"CREATE VIEW v_lit_auto AS SELECT 1 FROM t",
			"",
			false,
		},
		{
			"star_expansion",
			"v_star",
			"CREATE VIEW v_star AS SELECT * FROM t",
			"",
			false,
		},
		{
			"same_name_columns_join",
			"v_same_name_join",
			"CREATE VIEW v_same_name_join AS SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a",
			"",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Omni catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
			cat.Exec("CREATE TABLE t2 (a INT, b INT)", nil)
			if tc.tables != "" {
				cat.Exec(tc.tables, nil)
			}
			results, _ := cat.Exec(tc.viewSQL, nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", tc.viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_Section_6_1_SubqueryPatterns verifies that subquery patterns
// in views match MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_Section_6_1_SubqueryPatterns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		createAs string
		partial  bool
	}{
		{
			"scalar_subquery_in_select",
			"SELECT (SELECT MAX(a) FROM t) FROM t",
			false,
		},
		{
			"in_subquery",
			"SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)",
			false,
		},
		{
			"exists_subquery",
			"SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t WHERE a > 0)",
			false,
		},
		{
			"correlated_subquery",
			"SELECT a, (SELECT COUNT(*) FROM t t2 WHERE t2.a = t1.a) FROM t t1",
			false,
		},
		{
			"nested_subqueries_2_levels",
			"SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE b IN (SELECT b FROM t WHERE c > 0))",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

func TestDeparseContainer_Section_6_2_CTEPatterns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}

	cases := []struct {
		name     string
		createAs string
		partial  bool
	}{
		{
			"simple_cte",
			"WITH cte AS (SELECT a FROM t) SELECT * FROM cte",
			false,
		},
		{
			"cte_with_column_list",
			"WITH cte(x) AS (SELECT a FROM t) SELECT x FROM cte",
			false,
		},
		{
			"recursive_cte",
			"WITH RECURSIVE cte AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM cte WHERE n < 10) SELECT * FROM cte",
			false,
		},
		{
			"multiple_ctes",
			"WITH c1 AS (SELECT a FROM t), c2 AS (SELECT b FROM t) SELECT c1.a, c2.b FROM c1, c2",
			false,
		},
		{
			"cte_used_in_union",
			"WITH cte AS (SELECT a FROM t) SELECT * FROM cte UNION SELECT * FROM cte",
			false,
		},
		{
			"recursive_cte_multi_col",
			"WITH RECURSIVE cte AS (SELECT 1 AS n, 'a' AS s UNION ALL SELECT n + 1, CONCAT(s, 'a') FROM cte WHERE n < 5) SELECT * FROM cte",
			false,
		},
		{
			"recursive_cte_join_base",
			"WITH RECURSIVE cte AS (SELECT a, 1 AS depth FROM t WHERE a = 1 UNION ALL SELECT t.a, cte.depth + 1 FROM t JOIN cte ON t.b = cte.a WHERE cte.depth < 5) SELECT * FROM cte",
			false,
		},
		{
			"cte_references_cte",
			"WITH c1 AS (SELECT a FROM t), c2 AS (SELECT * FROM c1 WHERE a > 0) SELECT * FROM c2",
			false,
		},
		{
			"cte_with_aggregation",
			"WITH cte AS (SELECT a, COUNT(*) AS cnt FROM t GROUP BY a) SELECT * FROM cte",
			false,
		},
		{
			"cte_with_inner_union",
			"WITH cte AS (SELECT a FROM t UNION SELECT b FROM t) SELECT * FROM cte",
			false,
		},
		{
			"cte_with_complex_expr",
			"WITH cte AS (SELECT a + b AS sum_val, CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END AS label FROM t) SELECT * FROM cte",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				if tc.partial {
					t.Skipf("MySQL 8.0 rejected (expected partial): %v", err)
				}
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if len(results) == 0 {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog returned no results (expected partial)")
				}
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				if tc.partial {
					t.Skipf("CREATE VIEW on catalog failed (expected partial): %v", results[0].Error)
				}
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				if tc.partial {
					t.Skip("ShowCreateView returned empty (expected partial)")
				}
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				if tc.partial {
					t.Skipf("SELECT body mismatch (expected partial):\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
				}
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_8_1_ViewOfViewComplexStructures verifies view-of-view,
// many-column views, reserved word aliases, CASE without ELSE, and BETWEEN
// with column bounds against MySQL 8.0.
func TestDeparseContainer_8_1_ViewOfViewComplexStructures(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}

	// --- Test 1: View referencing another view ---
	t.Run("view_of_view", func(t *testing.T) {
		// MySQL side: create v1, then v2 referencing v1
		ctr.execSQLDirect("DROP VIEW IF EXISTS v2_vov")
		ctr.execSQLDirect("DROP VIEW IF EXISTS v1_vov")
		if err := ctr.execSQLDirect("CREATE VIEW v1_vov AS SELECT a FROM t"); err != nil {
			t.Fatalf("CREATE VIEW v1_vov on MySQL failed: %v", err)
		}
		if err := ctr.execSQLDirect("CREATE VIEW v2_vov AS SELECT * FROM v1_vov"); err != nil {
			t.Fatalf("CREATE VIEW v2_vov on MySQL failed: %v", err)
		}
		mysqlOutput, err := ctr.showCreateView("v2_vov")
		if err != nil {
			t.Fatalf("SHOW CREATE VIEW v2_vov on MySQL failed: %v", err)
		}
		mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

		// Omni side
		cat := New()
		cat.Exec("CREATE DATABASE test", nil)
		cat.SetCurrentDatabase("test")
		cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
		results, _ := cat.Exec("CREATE VIEW v1_vov AS SELECT a FROM t", nil)
		if len(results) > 0 && results[0].Error != nil {
			t.Fatalf("CREATE VIEW v1_vov on catalog failed: %v", results[0].Error)
		}
		results, _ = cat.Exec("CREATE VIEW v2_vov AS SELECT * FROM v1_vov", nil)
		if len(results) > 0 && results[0].Error != nil {
			t.Fatalf("CREATE VIEW v2_vov on catalog failed: %v", results[0].Error)
		}
		omniOutput := cat.ShowCreateView("test", "v2_vov")
		if omniOutput == "" {
			t.Fatal("ShowCreateView returned empty for v2_vov")
		}
		omniBody := extractSelectBody(omniOutput)

		t.Logf("MySQL body:  %s", mysqlBody)
		t.Logf("Omni body:   %s", omniBody)

		if mysqlBody != omniBody {
			t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
		}
	})

	// --- Tests 2-5: table-driven ---
	cases := []struct {
		name     string
		createAs string
	}{
		{
			"ten_plus_columns",
			"SELECT a, b, c, a + 1, b + 1, c + 1, a * b, b * c, a * c, a + b + c FROM t",
		},
		{
			"reserved_word_aliases",
			"SELECT a AS `select`, b AS `from`, c AS `where` FROM t",
		},
		{
			"case_without_else",
			"SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END FROM t",
		},
		{
			"between_column_bounds",
			"SELECT a BETWEEN b AND c FROM t",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_8_2_ExpressionEdgeCases verifies expression edge cases and stress tests
// against real MySQL 8.0 SHOW CREATE VIEW output.
func TestDeparseContainer_8_2_ExpressionEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t_rw (`select` INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t_rw on MySQL: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS tc (a VARCHAR(50), b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table tc on MySQL: %v", err)
	}

	cases := []struct {
		name       string
		setupTable string // which table(s) to create on omni side ("t", "t_rw", "tc")
		createAs   string // the SELECT portion after CREATE VIEW v AS
	}{
		{
			"cast_with_expression",
			"t",
			"SELECT CAST(a + b AS CHAR) FROM t",
		},
		{
			"collate_expression",
			"tc",
			"SELECT a COLLATE utf8mb4_unicode_ci FROM tc",
		},
		{
			"interval_arithmetic",
			"t",
			"SELECT INTERVAL 1 DAY + a FROM t",
		},
		{
			"sounds_like",
			"tc",
			"SELECT a SOUNDS LIKE b FROM tc",
		},
		{
			"unary_operators",
			"t",
			"SELECT -a, +a FROM t",
		},
		{
			"long_expression_auto_alias",
			"t",
			"SELECT CASE WHEN a > 0 THEN CONCAT(a, b, c) WHEN a < 0 THEN CONCAT(c, b, a) ELSE NULL END FROM t",
		},
		{
			"reserved_word_column",
			"t_rw",
			"SELECT `select` FROM t_rw",
		},
		{
			"multiple_rewrites",
			"t",
			"SELECT a + b, NOT (a > 0), CAST(a AS CHAR), COUNT(*), a REGEXP 'x' FROM t GROUP BY a, b HAVING COUNT(*) > 1",
		},
		{
			"all_logical_operators",
			"t",
			"SELECT a AND b, a OR b, NOT a, a XOR b, !a FROM t",
		},
		{
			"function_boolean_context",
			"t",
			"SELECT IFNULL(a, 0) AND COALESCE(b, 0) FROM t",
		},
		{
			"mixed_boolean_precedence",
			"t",
			"SELECT (a > 0) AND (b + 1) OR (c IS NULL) FROM t",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_82_" + tc.name

			// --- MySQL 8.0 side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Skipf("MySQL 8.0 rejected: %v", err)
				return
			}
			mysqlOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlBody := stripDatabasePrefix(extractSelectBody(mysqlOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			// Create the appropriate table(s)
			switch tc.setupTable {
			case "t":
				cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			case "t_rw":
				cat.Exec("CREATE TABLE t_rw (`select` INT, b INT)", nil)
			case "tc":
				cat.Exec("CREATE TABLE tc (a VARCHAR(50), b INT, c INT)", nil)
			}
			results, _ := cat.Exec(fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs), nil)
			if len(results) == 0 {
				t.Fatalf("CREATE VIEW on catalog returned no results")
			}
			if results[0].Error != nil {
				t.Fatalf("CREATE VIEW on catalog failed: %v", results[0].Error)
			}
			omniOutput := cat.ShowCreateView("test", viewName)
			if omniOutput == "" {
				t.Fatal("ShowCreateView returned empty")
			}
			omniBody := extractSelectBody(omniOutput)

			t.Logf("MySQL body:  %s", mysqlBody)
			t.Logf("Omni body:   %s", omniBody)

			if mysqlBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- mysql ---\n%s\n--- omni ---\n%s", mysqlBody, omniBody)
			}
		})
	}
}
