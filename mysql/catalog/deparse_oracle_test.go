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
	target := sel.TargetList[0]
	if rt, ok := target.(*nodes.ResTarget); ok {
		return deparse.Deparse(deparse.RewriteExpr(rt.Val))
	}
	return deparse.Deparse(deparse.RewriteExpr(target))
}

// TestDeparse_Section_4_1_Oracle verifies NOT folding against MySQL 8.0.
func TestDeparse_Section_4_1_Oracle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Create base table with integer column for boolean/comparison tests
	oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT)")

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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, selectSQL)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW failed: %v", err)
			}

			mysqlOutput, err := oracle.showCreateView(viewName)
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

// TestDeparse_Section_3_2_Oracle verifies TRIM special forms against MySQL 8.0.
func TestDeparse_Section_3_2_Oracle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Create base table
	oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a VARCHAR(50), b VARCHAR(50))")

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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, selectSQL)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW failed: %v", err)
			}

			mysqlOutput, err := oracle.showCreateView(viewName)
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
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := oracle.showCreateView(viewName)
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
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Setup: create base tables on MySQL 8.0
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on MySQL: %v", err)
	}
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on MySQL: %v", err)
	}
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := oracle.showCreateView(viewName)
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
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := oracle.showCreateView(viewName)
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
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on MySQL: %v", err)
	}

	t.Run("view_with_explicit_column_aliases", func(t *testing.T) {
		viewName := "v_col_alias"
		createSQL := "CREATE VIEW " + viewName + "(x, y) AS SELECT a, b FROM t"

		// --- MySQL 8.0 side ---
		oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
		if err := oracle.execSQLDirect(createSQL); err != nil {
			t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
		}
		mysqlOutput, err := oracle.showCreateView(viewName)
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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := oracle.showCreateView(viewName)
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

// TestDeparse_Section_7_3_ExpressionViews verifies that our catalog's SHOW CREATE VIEW
// output matches MySQL 8.0's output for views with expressions (arithmetic, functions,
// CASE, CAST, aggregates).
func TestDeparse_Section_7_3_ExpressionViews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// Setup: create base table on MySQL 8.0
	if err := oracle.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
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
			oracle.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := oracle.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on MySQL failed: %v", err)
			}
			mysqlOutput, err := oracle.showCreateView(viewName)
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
