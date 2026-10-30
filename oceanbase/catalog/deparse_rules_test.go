package catalog

import (
	"fmt"
	"testing"
)

// TestMySQL_DeparseRules creates views with specific SQL patterns and examines
// SHOW CREATE VIEW output to discover MySQL 8.0's exact deparsing/formatting rules.
// This is a research test — it asserts nothing, only logs results.
func TestMySQL_DeparseRules(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Create base tables used by the views.
	setupSQL := `
		CREATE TABLE t (a INT, b INT, c INT);
		CREATE TABLE t1 (a INT, b INT);
		CREATE TABLE t2 (a INT, b INT);
	`
	if err := ctr.execSQL(setupSQL); err != nil {
		t.Fatalf("setup tables: %v", err)
	}

	type viewTest struct {
		name     string
		createAs string // the SELECT part after "CREATE VIEW vN AS "
	}

	categories := []struct {
		label string
		views []viewTest
	}{
		{
			label: "A. Keyword Casing",
			views: []viewTest{
				{"v1", "SELECT 1"},
				{"v2", "SELECT a FROM t WHERE a > 1 GROUP BY a HAVING COUNT(*) > 1 ORDER BY a LIMIT 10"},
			},
		},
		{
			label: "B. Identifier Quoting",
			views: []viewTest{
				{"v3", "SELECT a, b AS alias1 FROM t"},
				{"v4", "SELECT t.a FROM t"},
				{"v5", "SELECT `select` FROM t"},
			},
		},
		{
			label: "C. Operator Formatting (spacing, case, parens)",
			views: []viewTest{
				{"v6", "SELECT a+b, a-b, a*b, a/b, a DIV b, a MOD b, a%b FROM t"},
				{"v7", "SELECT a=b, a<>b, a!=b, a>b, a<b, a>=b, a<=b, a<=>b FROM t"},
				{"v8", "SELECT a AND b, a OR b, NOT a, a XOR b FROM t"},
				{"v9", "SELECT a|b, a&b, a^b, a<<b, a>>b, ~a FROM t"},
				{"v10", "SELECT a IS NULL, a IS NOT NULL, a IS TRUE, a IS FALSE FROM t"},
				{"v11", "SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t"},
				{"v12", "SELECT a BETWEEN 1 AND 10, a NOT BETWEEN 1 AND 10 FROM t"},
				{"v13", "SELECT a LIKE 'foo%', a NOT LIKE 'bar%', a LIKE 'x' ESCAPE '\\\\' FROM t"},
				{"v14", "SELECT a REGEXP 'pattern', a NOT REGEXP 'pattern' FROM t"},
			},
		},
		{
			label: "D. Function Name Casing",
			views: []viewTest{
				{"v15", "SELECT COUNT(*), SUM(a), AVG(a), MAX(a), MIN(a), COUNT(DISTINCT a) FROM t"},
				{"v16", "SELECT CONCAT(a, b), SUBSTRING(a, 1, 3), TRIM(a), UPPER(a), LOWER(a) FROM t"},
				{"v17", "SELECT NOW(), CURRENT_TIMESTAMP, CURRENT_DATE, CURRENT_TIME, CURRENT_USER FROM t"},
				{"v18", "SELECT IFNULL(a, 0), COALESCE(a, b, 0), NULLIF(a, 0), IF(a > 0, 'yes', 'no') FROM t"},
				{"v19", "SELECT CAST(a AS CHAR), CAST(a AS SIGNED), CONVERT(a, CHAR), CONVERT(a USING utf8mb4) FROM t"},
			},
		},
		{
			label: "E. Expression Structure",
			views: []viewTest{
				{"v20", "SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END FROM t"},
				{"v21", "SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t"},
				{"v22", "SELECT (a + b) * c, a + (b * c) FROM t"},
				{"v23", "SELECT -a, +a, !a FROM t"},
			},
		},
		{
			label: "F. Literals",
			views: []viewTest{
				{"v24", "SELECT 1, 1.5, 'hello', NULL, TRUE, FALSE FROM t"},
				{"v25", "SELECT 0xFF, X'FF', 0b1010, b'1010' FROM t"},
				{"v26", "SELECT _utf8mb4'hello', _latin1'world' FROM t"},
				{"v27", "SELECT DATE '2024-01-01', TIME '12:00:00', TIMESTAMP '2024-01-01 12:00:00' FROM t"},
			},
		},
		{
			label: "G. Subqueries",
			views: []viewTest{
				{"v28", "SELECT (SELECT MAX(a) FROM t) FROM t"},
				{"v29", "SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)"},
				{"v30", "SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t WHERE a > 0)"},
			},
		},
		{
			label: "H. JOIN Formatting",
			views: []viewTest{
				{"v31", "SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a"},
				{"v32", "SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a"},
				{"v33", "SELECT t1.a, t2.b FROM t1 RIGHT JOIN t2 ON t1.a = t2.a"},
				{"v34", "SELECT t1.a, t2.b FROM t1 CROSS JOIN t2"},
				{"v35", "SELECT * FROM t1 NATURAL JOIN t2"},
				{"v36", "SELECT t1.a, t2.b FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a"},
				{"v37", "SELECT t1.a, t2.b FROM t1 JOIN t2 USING (a)"},
				{"v38", "SELECT t1.a, t2.b FROM t1, t2 WHERE t1.a = t2.a"},
			},
		},
		{
			label: "I. UNION/INTERSECT/EXCEPT",
			views: []viewTest{
				{"v39", "SELECT a FROM t UNION SELECT b FROM t"},
				{"v40", "SELECT a FROM t UNION ALL SELECT b FROM t"},
				{"v41", "SELECT a FROM t UNION SELECT b FROM t UNION SELECT c FROM t"},
			},
		},
		{
			label: "J. Window Functions",
			views: []viewTest{
				{"v42", "SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t"},
				{"v43", "SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b) FROM t"},
				{"v44", "SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t"},
			},
		},
		{
			label: "K. Alias Formatting",
			views: []viewTest{
				{"v45", "SELECT a AS col1, b col2 FROM t"},
				{"v46", "SELECT a + b AS sum_col FROM t"},
				{"v47", "SELECT * FROM t AS t1"},
			},
		},
		{
			label: "L. Misc MySQL-specific",
			views: []viewTest{
				{"v48", "SELECT a->>'$.key' FROM t"},
				{"v49", "SELECT INTERVAL 1 DAY + a FROM t"},
			},
		},
		{
			label: "M. Semantic Rewrites (things MySQL changes fundamentally)",
			views: []viewTest{
				// NOT LIKE → not((...like...))
				{"v50", "SELECT a NOT LIKE 'x%' FROM t"},
				// != → <>
				{"v51", "SELECT a != b FROM t"},
				// COUNT(*) → count(0)
				{"v52", "SELECT COUNT(*) FROM t"},
				// SUBSTRING → substr
				{"v53", "SELECT SUBSTRING('abc', 1, 2) FROM t"},
				// CURRENT_TIMESTAMP → now()
				{"v54", "SELECT CURRENT_TIMESTAMP() FROM t"},
				// MOD → %
				{"v55", "SELECT a MOD b FROM t"},
				// +a → a (unary plus dropped)
				{"v56", "SELECT +a FROM t"},
				// ! → (0 = ...)
				{"v57", "SELECT !a FROM t"},
				// AND/OR on INT cols → (0 <> ...) and/or (0 <> ...)
				{"v58", "SELECT a AND b FROM t"},
				// REGEXP → regexp_like()
				{"v59", "SELECT a REGEXP 'x' FROM t"},
				// -> → json_extract()
				{"v60", "SELECT a->'$.key' FROM t"},
			},
		},
		{
			label: "N. Derived Tables & Subquery in FROM",
			views: []viewTest{
				{"v61", "SELECT d.x FROM (SELECT a AS x FROM t) d"},
				{"v62", "SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10"},
			},
		},
		{
			label: "O. GROUP_CONCAT & Special Aggregates",
			views: []viewTest{
				{"v63", "SELECT GROUP_CONCAT(a ORDER BY a SEPARATOR ',') FROM t"},
				{"v64", "SELECT GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';') FROM t"},
			},
		},
		{
			label: "P. Spacing & Comma Rules",
			views: []viewTest{
				// Are there spaces after commas in function args?
				{"v65", "SELECT CONCAT(a, b, c) FROM t"},
				// Spaces in IN list?
				{"v66", "SELECT a IN (1, 2, 3) FROM t"},
				// Space before/after AS?
				{"v67", "SELECT a AS x FROM t"},
			},
		},
		{
			label: "Q. Table Alias Format (AS vs space)",
			views: []viewTest{
				{"v68", "SELECT x.a FROM t AS x"},
				{"v69", "SELECT x.a FROM t x"},
			},
		},
		{
			label: "R. Complex Precedence",
			views: []viewTest{
				{"v70", "SELECT a + b + c FROM t"},
				{"v71", "SELECT a * b + c FROM t"},
				{"v72", "SELECT a + b * c FROM t"},
				{"v73", "SELECT a OR b AND c FROM t"},
				{"v74", "SELECT (a OR b) AND c FROM t"},
				{"v75", "SELECT a > 0 AND b < 10 OR c = 5 FROM t"},
			},
		},
		{
			label: "S. CTE (WITH clause)",
			views: []viewTest{
				{"v76", "WITH cte AS (SELECT a FROM t) SELECT * FROM cte"},
				{"v77", "WITH cte(x) AS (SELECT a FROM t) SELECT x FROM cte"},
			},
		},
		{
			label: "T. Type-Aware Boolean Context (expressions in AND/OR)",
			views: []viewTest{
				// Expression (not column) in boolean context
				{"v78", "SELECT (a + 1) AND b FROM t"},
				{"v79", "SELECT (a > 0) AND (b + 1) FROM t"},
				{"v80", "SELECT (a > 0) AND (b > 0) FROM t"},
				// Function results in boolean context
				{"v81", "SELECT ABS(a) AND b FROM t"},
				{"v82", "SELECT COUNT(*) > 0 FROM t"},
				// CASE in boolean context
				{"v83", "SELECT CASE WHEN a > 0 THEN 1 ELSE 0 END AND b FROM t"},
				// Nested boolean
				{"v84", "SELECT NOT (a + 1) FROM t"},
				{"v85", "SELECT NOT (a > 0) FROM t"},
				// IF result in boolean context
				{"v86", "SELECT IF(a > 0, 1, 0) AND b FROM t"},
				// Subquery in boolean context
				{"v87", "SELECT (SELECT MAX(a) FROM t) AND b FROM t"},
			},
		},
		{
			label: "U. CAST charset inference",
			views: []viewTest{
				{"v88", "SELECT CAST(a AS CHAR(10)) FROM t"},
				{"v89", "SELECT CAST(a AS BINARY) FROM t"},
				{"v90", "SELECT CAST(a AS DECIMAL(10,2)) FROM t"},
				{"v91", "SELECT CAST(a AS UNSIGNED) FROM t"},
				{"v92", "SELECT CAST(a AS DATE) FROM t"},
				{"v93", "SELECT CAST(a AS DATETIME) FROM t"},
				{"v94", "SELECT CAST(a AS JSON) FROM t"},
			},
		},
		{
			label: "V. String vs INT in boolean context",
			views: []viewTest{
				// Need a VARCHAR column to test string behavior
				{"v95_setup", "SELECT 'hello' AND 1 FROM t"},
				{"v96", "SELECT a AND 'hello' FROM t"},
				{"v97", "SELECT CONCAT(a,b) AND 1 FROM t"},
			},
		},
		{
			label: "W. Complex function return types",
			views: []viewTest{
				{"v98", "SELECT IFNULL(a, 0) AND b FROM t"},
				{"v99", "SELECT COALESCE(a, b) AND 1 FROM t"},
				{"v100", "SELECT NULLIF(a, 0) AND b FROM t"},
				{"v101", "SELECT GREATEST(a, b) AND 1 FROM t"},
				{"v102", "SELECT LEAST(a, b) AND 1 FROM t"},
			},
		},
		{
			label: "X. Comparison results NOT wrapped",
			views: []viewTest{
				// These should NOT get (0 <> ...) wrapping because they're already boolean
				{"v103", "SELECT (a = b) AND (a > 0) FROM t"},
				{"v104", "SELECT (a IN (1,2,3)) AND (b BETWEEN 1 AND 10) FROM t"},
				{"v105", "SELECT (a IS NULL) AND (b LIKE 'x%') FROM t"},
				{"v106", "SELECT (a = 1) OR (b = 2) FROM t"},
				{"v107", "SELECT EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0) FROM t"},
			},
		},
		{
			label: "Y. Multi-table column qualification",
			views: []viewTest{
				{"v108", "SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a"},
				{"v109", "SELECT a FROM t WHERE a > 0"},
				{"v110", "SELECT t.a FROM t"},
			},
		},
		{
			label: "Z. Edge cases",
			views: []viewTest{
				// Empty string
				{"v111", "SELECT '' FROM t"},
				// Escaped quotes in strings
				{"v112", "SELECT 'it''s' FROM t"},
				// Backslash in strings
				{"v113", "SELECT 'back\\\\slash' FROM t"},
				// Very long expression
				{"v114", "SELECT a + b + c + a + b + c FROM t"},
				// Nested function calls
				{"v115", "SELECT CONCAT(UPPER(TRIM(a)), LOWER(b)) FROM t"},
				// Multiple aggregates
				{"v116", "SELECT COUNT(*), SUM(a), AVG(b), MAX(c) FROM t GROUP BY a"},
			},
		},
	}

	t.Log("=== MySQL 8.0 Deparse Rules Research ===")
	t.Log("")

	for _, cat := range categories {
		t.Logf("--- %s ---", cat.label)
		for _, vt := range cat.views {
			createSQL := fmt.Sprintf("CREATE VIEW %s AS %s", vt.name, vt.createAs)

			if err := ctr.execSQLDirect(createSQL); err != nil {
				t.Logf("  [%s] CREATE failed: %v", vt.name, err)
				t.Logf("    INPUT:  %s", createSQL)
				t.Logf("")
				continue
			}

			output, err := ctr.showCreateView(vt.name)
			if err != nil {
				t.Logf("  [%s] SHOW CREATE VIEW failed: %v", vt.name, err)
				t.Logf("    INPUT:  %s", createSQL)
				t.Logf("")
				continue
			}

			t.Logf("  [%s]", vt.name)
			t.Logf("    INPUT:  %s", vt.createAs)
			t.Logf("    OUTPUT: %s", output)
			t.Logf("")
		}
		t.Log("")
	}
}
