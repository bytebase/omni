package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// The view tests below hold the catalog's SHOW CREATE VIEW body to
// OceanBase's, for the definitions the two print alike once OceanBase's
// `test`. qualifier is taken off. Where they print a definition differently
// (boolean contexts, function-name case, joins, set operations, window
// functions, CTEs), the catalog prints MySQL 8.0's text, and mysql/catalog's
// copy of this file holds those cases to MySQL.

// extractSelectBody extracts the SELECT body from SHOW CREATE VIEW output.
// OceanBase format:
//
//	CREATE VIEW `v` AS select ...
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
// output matches OceanBase's output for simple view definitions.
// For each view, we compare the SELECT body portion (after "AS ").
func TestDeparse_Section_7_2_SimpleViews(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
	}{
		{"select_constant", "SELECT 1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := extractSelectBody(engineOutput)

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

			t.Logf("OceanBase full:  %s", engineOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparse_Section_7_4_JoinViews verifies that our catalog's SHOW CREATE VIEW
// output matches OceanBase's output for views with JOINs (INNER JOIN, LEFT JOIN,
// multiple tables, subquery in FROM).
func TestDeparse_Section_7_4_JoinViews(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
		tables   string // additional CREATE TABLE statements for our catalog (beyond t, t1, t2)
	}{
		{"subquery_from", "SELECT d.x FROM (SELECT a AS x FROM t) d", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase full:  %s", engineOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// stripDatabasePrefix removes the `test`. database prefix from OceanBase SHOW CREATE VIEW output.
// OceanBase qualifies all identifiers with the database name (e.g., `test`.`t`.`a`),
// while our catalog does not. We strip the prefix for comparison.
func stripDatabasePrefix(s string) string {
	return strings.ReplaceAll(s, "`test`.", "")
}

// TestDeparse_Section_7_6_Regression verifies that the deparser integration does not
// break existing tests (scenarios 1-2 are covered by running go test ./oceanbase/catalog/
// and go test ./oceanbase/parser/ separately) and that views with explicit column
// aliases match OceanBase output exactly.
func TestDeparse_Section_7_6_Regression(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	// Verify simple and complex views still match (re-run 7.2 and 7.5 representative cases).
	simpleAndComplexCases := []struct {
		name     string
		createAs string
	}{
		{"simple_select_column", "SELECT a FROM t"},
		{"simple_select_where", "SELECT a FROM t WHERE a > 0"},
	}

	for _, tc := range simpleAndComplexCases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_reg_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_1_1_ArithmeticComparison verifies arithmetic and comparison operators
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_1_1_ArithmeticComparison(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
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
		{"mod", "v_mod", "CREATE VIEW v_mod AS SELECT a MOD b FROM t"},
		{"equals", "v_eq", "CREATE VIEW v_eq AS SELECT a = b FROM t"},
		{"not_equals_bang", "v_neq", "CREATE VIEW v_neq AS SELECT a != b FROM t"},
		{"not_equals_ltgt", "v_neq2", "CREATE VIEW v_neq2 AS SELECT a <> b FROM t"},
		{"comparisons", "v_cmp", "CREATE VIEW v_cmp AS SELECT a > b, a < b, a >= b, a <= b FROM t"},
		{"null_safe_equals", "v_nseq", "CREATE VIEW v_nseq AS SELECT a <=> b FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparse_Section_7_3_ExpressionViews verifies that our catalog's SHOW CREATE VIEW
// output matches OceanBase's output for views with expressions (arithmetic, functions,
// CASE, CAST, aggregates).
func TestDeparse_Section_7_3_ExpressionViews(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		createAs string // the SELECT portion after CREATE VIEW v AS
	}{
		{"arithmetic_expr", "SELECT a + b FROM t"},
		{"case_expr", "SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t"},
		{"cast_expr", "SELECT CAST(a AS CHAR) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase full:  %s", engineOutput)
			t.Logf("Omni full:   %s", omniOutput)
			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_1_3_LiteralsSpacing verifies literals and spacing rules
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_1_3_LiteralsSpacing(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"escaped_quotes", "v_esc_quotes", "CREATE VIEW v_esc_quotes AS SELECT 'it''s' FROM t"},
		{"escaped_backslash", "v_esc_bslash", "CREATE VIEW v_esc_bslash AS SELECT 'back\\\\slash' FROM t"},
		{"in_list_no_space", "v_in_list", "CREATE VIEW v_in_list AS SELECT a IN (1, 2, 3) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_1_2_LogicalBitwiseIS verifies logical, bitwise, and IS operators
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_1_2_LogicalBitwiseIS(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"bitwise_ops", "v_bitwise", "CREATE VIEW v_bitwise AS SELECT a | b, a & b, a ^ b FROM t"},
		{"shifts", "v_shifts", "CREATE VIEW v_shifts AS SELECT a << b, a >> b FROM t"},
		{"bitwise_not", "v_bitnot", "CREATE VIEW v_bitnot AS SELECT ~a FROM t"},
		{"is_null", "v_isnull", "CREATE VIEW v_isnull AS SELECT a IS NULL, a IS NOT NULL FROM t"},
		{"in_not_in", "v_in", "CREATE VIEW v_in AS SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t"},
		{"like_escape", "v_like_esc", "CREATE VIEW v_like_esc AS SELECT a LIKE 'x' ESCAPE '\\\\' FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_Section_2_1_FunctionNameRewrites verifies that function name
// rewrites match OceanBase SHOW CREATE VIEW output.
// Covers: SUBSTRING->substr, CURRENT_TIMESTAMP->now(), CURRENT_DATE->curdate(),
// CURRENT_TIME->curtime(), CURRENT_USER->current_user(), NOW()->now(),
// COUNT(*)->count(0), COUNT(DISTINCT).
func TestDeparseContainer_Section_2_1_FunctionNameRewrites(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
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
		{"count_distinct", "v_count_dist", "CREATE VIEW v_count_dist AS SELECT COUNT(DISTINCT a) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_2_2_RegularFunctionsAggregates verifies regular functions and aggregates
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_2_2_RegularFunctionsAggregates(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"sum_avg_max_min", "v_agg_funcs", "CREATE VIEW v_agg_funcs AS SELECT SUM(a), AVG(a), MAX(a), MIN(a) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_2_4_CastConvertOperatorRewrites verifies CAST, CONVERT,
// REGEXP→regexp_like, NOT REGEXP, -> (json_extract), ->> (json_unquote(json_extract))
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_2_4_CastConvertOperatorRewrites(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS tj (a JSON, b INT)"); err != nil {
		t.Fatalf("failed to create table tj on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		tables   string // extra CREATE TABLE for omni catalog (beyond t)
	}{
		{"cast_char", "v_cast_char", "CREATE VIEW v_cast_char AS SELECT CAST(a AS CHAR) FROM t", ""},
		{"cast_char10", "v_cast_char10", "CREATE VIEW v_cast_char10 AS SELECT CAST(a AS CHAR(10)) FROM t", ""},
		{"cast_binary", "v_cast_binary", "CREATE VIEW v_cast_binary AS SELECT CAST(a AS BINARY) FROM t", ""},
		{"cast_signed_unsigned", "v_cast_su", "CREATE VIEW v_cast_su AS SELECT CAST(a AS SIGNED), CAST(a AS UNSIGNED) FROM t", ""},
		{"convert_char", "v_conv_char", "CREATE VIEW v_conv_char AS SELECT CONVERT(a, CHAR) FROM t", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_2_3_SpecialFunctions verifies TRIM, GROUP_CONCAT, and simple CASE
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_2_3_SpecialFunctions(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS tv (a VARCHAR(50), b VARCHAR(50), c INT)"); err != nil {
		t.Fatalf("failed to create table tv on OceanBase: %v", err)
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
		{"simple_case", "v_simple_case", "CREATE VIEW v_simple_case AS SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_3_1_BooleanContextWrapping verifies that non-boolean expressions
// in boolean context (AND/OR) get (0 <> ...) wrapping to match OceanBase's
// SHOW CREATE VIEW output.
func TestDeparseContainer_3_1_BooleanContextWrapping(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"cmp_and_cmp", "v_bool_cmpx2", "CREATE VIEW v_bool_cmpx2 AS SELECT (a > 0) AND (b > 0) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_3_2_NotFoldingNoDoubleWrap verifies NOT folding and
// no-double-wrapping of boolean expressions against OceanBase.
// Section 3.2 scenarios:
//   - NOT(comparison) folds into inverted operator
//   - NOT(non-boolean) becomes (0 = ...)
//   - ! is same as NOT
//   - comparisons in AND are NOT double-wrapped
//   - predicates (IN, BETWEEN) in AND are NOT wrapped
//   - IS/LIKE in AND are NOT wrapped
//   - EXISTS in AND is NOT wrapped
func TestDeparseContainer_3_2_NotFoldingNoDoubleWrap(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"cmp_not_double_wrapped", "v_s32_cmp_ndw", "CREATE VIEW v_s32_cmp_ndw AS SELECT (a = b) AND (a > 0) FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_3_3_ComplexPrecedence verifies complex operator precedence
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_3_3_ComplexPrecedence(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table on OceanBase: %v", err)
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
		{"mixed_cmp_logic", "v_mixed", "CREATE VIEW v_mixed AS SELECT a > 0 AND b < 10 OR c = 5 FROM t"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_4_2_MultiTableDerived verifies multi-table JOINs, chained LEFT JOINs,
// derived tables, and table aliases against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_4_2_MultiTableDerived(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t3 (b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t3 on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"derived_table", "v_derived", "CREATE VIEW v_derived AS SELECT d.x FROM (SELECT a AS x FROM t) d"},
		{"derived_with_where", "v_derived_where", "CREATE VIEW v_derived_where AS SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10"},
		{"table_alias_with_as", "v_alias_as", "CREATE VIEW v_alias_as AS SELECT x.a FROM t AS x"},
		{"table_alias_without_as", "v_alias_noas", "CREATE VIEW v_alias_noas AS SELECT x.a FROM t x"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_Section_5_1_SelectClauses verifies that our catalog's SHOW CREATE VIEW
// output matches OceanBase's output for views with WHERE, GROUP BY, HAVING, ORDER BY,
// LIMIT, OFFSET, DISTINCT, and expression-based GROUP BY.
func TestDeparseContainer_Section_5_1_SelectClauses(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
	}{
		{"limit_offset", "v_limit_offset", "CREATE VIEW v_limit_offset AS SELECT a FROM t LIMIT 10 OFFSET 5"},
		{"expression_group_by", "v_expr_groupby", "CREATE VIEW v_expr_groupby AS SELECT a + b FROM t GROUP BY a + b"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_5_3_ColumnAliasPatterns verifies column and alias patterns
// against OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_5_3_ColumnAliasPatterns(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t2 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t2 on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		viewName string
		viewSQL  string
		tables   string // extra tables for omni catalog (beyond t, t1, t2)
	}{
		{
			"explicit_alias_as_vs_space",
			"v_alias_as_space",
			"CREATE VIEW v_alias_as_space AS SELECT a AS col1, b col2 FROM t",
			"",
		},
		{
			"expression_explicit_alias",
			"v_expr_alias",
			"CREATE VIEW v_expr_alias AS SELECT a + b AS sum_col FROM t",
			"",
		},
		{
			"literal_auto_alias",
			"v_lit_auto",
			"CREATE VIEW v_lit_auto AS SELECT 1 FROM t",
			"",
		},
		{
			"star_expansion",
			"v_star",
			"CREATE VIEW v_star AS SELECT * FROM t",
			"",
		},
	}

	// Two columns of one name: OceanBase refuses the view, as MySQL does
	// (ER_DUP_FIELDNAME, 1060).
	// omni gap: the catalog takes it.
	t.Run("same_name_columns_join", func(t *testing.T) {
		const view = "CREATE VIEW v_same_name_join AS SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a"
		oracleRefuses(t, ctr, view, 1060)
		cat := New()
		cat.Exec("CREATE DATABASE test", nil)
		cat.SetCurrentDatabase("test")
		cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
		cat.Exec("CREATE TABLE t2 (a INT, b INT)", nil)
		assertIntEq(t, "omni error code", omniCode(cat, view), 0)
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + tc.viewName)
			if err := ctr.execSQLDirect(tc.viewSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(tc.viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_Section_6_1_SubqueryPatterns verifies that subquery patterns
// in views match OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_Section_6_1_SubqueryPatterns(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t1 (a INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t1 on OceanBase: %v", err)
	}

	cases := []struct {
		name     string
		createAs string
	}{
		{
			"scalar_subquery_in_select",
			"SELECT (SELECT MAX(a) FROM t) FROM t",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

			// --- Our catalog side ---
			cat := New()
			cat.Exec("CREATE DATABASE test", nil)
			cat.SetCurrentDatabase("test")
			cat.Exec("CREATE TABLE t (a INT, b INT, c INT)", nil)
			cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_8_1_ViewOfViewComplexStructures verifies view-of-view,
// many-column views, reserved word aliases, CASE without ELSE, and BETWEEN
// with column bounds against OceanBase.
func TestDeparseContainer_8_1_ViewOfViewComplexStructures(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base table on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}

	// --- Test 1: View referencing another view ---
	t.Run("view_of_view", func(t *testing.T) {
		// OceanBase side: create v1, then v2 referencing v1
		ctr.execSQLDirect("DROP VIEW IF EXISTS v2_vov")
		ctr.execSQLDirect("DROP VIEW IF EXISTS v1_vov")
		if err := ctr.execSQLDirect("CREATE VIEW v1_vov AS SELECT a FROM t"); err != nil {
			t.Fatalf("CREATE VIEW v1_vov on OceanBase failed: %v", err)
		}
		if err := ctr.execSQLDirect("CREATE VIEW v2_vov AS SELECT * FROM v1_vov"); err != nil {
			t.Fatalf("CREATE VIEW v2_vov on OceanBase failed: %v", err)
		}
		engineOutput, err := ctr.showCreateView("v2_vov")
		if err != nil {
			t.Fatalf("SHOW CREATE VIEW v2_vov on OceanBase failed: %v", err)
		}
		engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

		t.Logf("OceanBase body:  %s", engineBody)
		t.Logf("Omni body:   %s", omniBody)

		if engineBody != omniBody {
			t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
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
			"between_column_bounds",
			"SELECT a BETWEEN b AND c FROM t",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}

// TestDeparseContainer_8_2_ExpressionEdgeCases verifies expression edge cases and stress tests
// against real OceanBase SHOW CREATE VIEW output.
func TestDeparseContainer_8_2_ExpressionEdgeCases(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create base tables on OceanBase
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table t on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS t_rw (`select` INT, b INT)"); err != nil {
		t.Fatalf("failed to create table t_rw on OceanBase: %v", err)
	}
	if err := ctr.execSQLDirect("CREATE TABLE IF NOT EXISTS tc (a VARCHAR(50), b INT, c INT)"); err != nil {
		t.Fatalf("failed to create table tc on OceanBase: %v", err)
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
			"reserved_word_column",
			"t_rw",
			"SELECT `select` FROM t_rw",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viewName := "v_82_" + tc.name

			// --- OceanBase side ---
			ctr.execSQLDirect("DROP VIEW IF EXISTS " + viewName)
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", viewName, tc.createAs)
			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Fatalf("CREATE VIEW on OceanBase failed: %v", err)
			}
			engineOutput, err := ctr.showCreateView(viewName)
			if err != nil {
				t.Fatalf("SHOW CREATE VIEW on OceanBase failed: %v", err)
			}
			engineBody := stripDatabasePrefix(extractSelectBody(engineOutput))

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

			t.Logf("OceanBase body:  %s", engineBody)
			t.Logf("Omni body:   %s", omniBody)

			if engineBody != omniBody {
				t.Errorf("SELECT body mismatch:\n--- oceanbase ---\n%s\n--- omni ---\n%s", engineBody, omniBody)
			}
		})
	}
}
