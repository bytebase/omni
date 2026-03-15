package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// ParseAndCheck is a test helper that parses SQL and returns the AST.
// It fails the test if parsing produces an error.
func ParseAndCheck(t *testing.T, sql string) *ast.List {
	t.Helper()
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): unexpected error: %v", sql, err)
	}
	return result
}

// ParseShouldFail is a test helper that expects parsing to fail.
func ParseShouldFail(t *testing.T, sql string) {
	t.Helper()
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("Parse(%q): expected error, got success", sql)
	}
}

// TestParseEmpty tests parsing of empty input.
func TestParseEmpty(t *testing.T) {
	result := ParseAndCheck(t, "")
	if result.Len() != 0 {
		t.Fatalf("expected 0 statements, got %d", result.Len())
	}
}

// TestParseSemicolon tests parsing of standalone semicolons.
func TestParseSemicolon(t *testing.T) {
	result := ParseAndCheck(t, ";;;")
	if result.Len() != 0 {
		t.Fatalf("expected 0 statements, got %d", result.Len())
	}
}

// TestLexerBasicTokens tests basic lexer functionality.
func TestLexerBasicTokens(t *testing.T) {
	tests := []struct {
		input string
		typ   int
		str   string
	}{
		{"123", tokICONST, "123"},
		{"'hello'", tokSCONST, "hello"},
		{"\"MyCol\"", tokQIDENT, "MyCol"},
		{":param", tokBIND, "param"},
		{":1", tokBIND, "1"},
		{":=", tokASSIGN, ":="},
		{"=>", tokASSOC, "=>"},
		{"..", tokDOTDOT, ".."},
		{"**", tokEXPON, "**"},
		{"||", tokCONCAT, "||"},
		{"<=", tokLESSEQ, "<="},
		{">=", tokGREATEQ, ">="},
		{"!=", tokNOTEQ, "!="},
		{"<>", tokNOTEQ, "<>"},
		{"~=", tokNOTEQ, "~="},
		{"^=", tokNOTEQ, "^="},
		{"<<", tokLABELOPEN, "<<"},
		{">>", tokLABELCLOSE, ">>"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tc.typ {
				t.Errorf("type: got %d, want %d", tok.Type, tc.typ)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
		})
	}
}

// TestLexerStrings tests string literal lexing.
func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input string
		str   string
	}{
		{"'hello'", "hello"},
		{"'it''s'", "it's"},
		{"''", ""},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tokSCONST {
				t.Errorf("type: got %d, want %d", tok.Type, tokSCONST)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
		})
	}
}

// TestLexerQQuote tests q-quote string lexing.
func TestLexerQQuote(t *testing.T) {
	tests := []struct {
		input string
		str   string
	}{
		{"q'[hello world]'", "hello world"},
		{"q'{it's}'", "it's"},
		{"q'(test)'", "test"},
		{"q'<angle>'", "angle"},
		{"q'!bang!'", "bang"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tokSCONST {
				t.Errorf("type: got %d, want %d", tok.Type, tokSCONST)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
		})
	}
}

// TestLexerNCharLiteral tests N'...' national character literal.
func TestLexerNCharLiteral(t *testing.T) {
	lex := NewLexer("N'hello'")
	tok := lex.NextToken()
	if tok.Type != tokNCHARLIT {
		t.Errorf("type: got %d, want %d", tok.Type, tokNCHARLIT)
	}
	if tok.Str != "hello" {
		t.Errorf("str: got %q, want %q", tok.Str, "hello")
	}
}

// TestLexerHints tests hint lexing.
func TestLexerHints(t *testing.T) {
	lex := NewLexer("/*+ FULL(t) */")
	tok := lex.NextToken()
	if tok.Type != tokHINT {
		t.Errorf("type: got %d, want %d", tok.Type, tokHINT)
	}
	if tok.Str != "FULL(t)" {
		t.Errorf("str: got %q, want %q", tok.Str, "FULL(t)")
	}
}

// TestLexerComments tests that comments are skipped.
func TestLexerComments(t *testing.T) {
	lex := NewLexer("-- line comment\n42")
	tok := lex.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Errorf("expected 42, got type=%d ival=%d", tok.Type, tok.Ival)
	}

	lex = NewLexer("/* block comment */ 42")
	tok = lex.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Errorf("expected 42, got type=%d ival=%d", tok.Type, tok.Ival)
	}
}

// TestLexerKeywords tests keyword recognition.
func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input string
		typ   int
	}{
		{"SELECT", kwSELECT},
		{"select", kwSELECT},
		{"Select", kwSELECT},
		{"FROM", kwFROM},
		{"WHERE", kwWHERE},
		{"INSERT", kwINSERT},
		{"UPDATE", kwUPDATE},
		{"DELETE", kwDELETE},
		{"MERGE", kwMERGE},
		{"CREATE", kwCREATE},
		{"ALTER", kwALTER},
		{"DROP", kwDROP},
		{"CONNECT", kwCONNECT},
		{"PRIOR", kwPRIOR},
		{"DECODE", kwDECODE},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tc.typ {
				t.Errorf("keyword %q: got type %d, want %d", tc.input, tok.Type, tc.typ)
			}
		})
	}
}

// TestLexerNumbers tests number literal lexing.
func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input string
		typ   int
		str   string
		ival  int64
	}{
		{"0", tokICONST, "0", 0},
		{"42", tokICONST, "42", 42},
		{"3.14", tokFCONST, "3.14", 0},
		{".5", tokFCONST, ".5", 0},
		{"1e10", tokFCONST, "1e10", 0},
		{"1.5E-3", tokFCONST, "1.5E-3", 0},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tc.typ {
				t.Errorf("type: got %d, want %d", tok.Type, tc.typ)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
			if tc.typ == tokICONST && tok.Ival != tc.ival {
				t.Errorf("ival: got %d, want %d", tok.Ival, tc.ival)
			}
		})
	}
}

// TestNodeToString tests AST serialization.
// TestParseSelectPivot tests PIVOT clause parsing.
func TestParseSelectPivot(t *testing.T) {
	tests := []string{
		// Basic PIVOT with single aggregate
		`SELECT * FROM sales PIVOT (SUM(amount) FOR quarter IN ('Q1', 'Q2', 'Q3', 'Q4'))`,
		// PIVOT with aliased aggregate
		`SELECT * FROM sales PIVOT (SUM(amount) AS total FOR quarter IN ('Q1' AS q1, 'Q2' AS q2))`,
		// PIVOT with multiple aggregates
		`SELECT * FROM sales PIVOT (SUM(amount) AS total, COUNT(*) AS cnt FOR quarter IN ('Q1', 'Q2'))`,
		// PIVOT from subquery
		`SELECT * FROM (SELECT dept_id, quarter, amount FROM sales) PIVOT (SUM(amount) FOR quarter IN ('Q1', 'Q2'))`,
		// PIVOT with multi-column FOR
		`SELECT * FROM sales PIVOT (SUM(amount) FOR (year, quarter) IN ((2023, 'Q1'), (2023, 'Q2')))`,
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
		})
	}
}

// TestParseSelectUnpivot tests UNPIVOT clause parsing.
func TestParseSelectUnpivot(t *testing.T) {
	tests := []string{
		// Basic UNPIVOT
		`SELECT * FROM quarterly_sales UNPIVOT (sales FOR quarter IN (q1, q2, q3, q4))`,
		// UNPIVOT INCLUDE NULLS
		`SELECT * FROM quarterly_sales UNPIVOT INCLUDE NULLS (sales FOR quarter IN (q1, q2, q3, q4))`,
		// UNPIVOT EXCLUDE NULLS
		`SELECT * FROM quarterly_sales UNPIVOT EXCLUDE NULLS (sales FOR quarter IN (q1, q2, q3, q4))`,
		// UNPIVOT with aliases
		`SELECT * FROM quarterly_sales UNPIVOT (sales FOR quarter IN (q1 AS 'Quarter 1', q2 AS 'Quarter 2'))`,
		// UNPIVOT from subquery
		`SELECT * FROM (SELECT id, q1, q2, q3, q4 FROM data) UNPIVOT (val FOR qtr IN (q1, q2, q3, q4))`,
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
		})
	}
}

// TestParseSelectModel tests MODEL clause parsing.
func TestParseSelectModel(t *testing.T) {
	tests := []string{
		// Basic MODEL with DIMENSION BY and MEASURES
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = sales['US', 2023] * 1.1)`,

		// MODEL with PARTITION BY
		`SELECT country, product, year, sales FROM sales_data
		 MODEL
		   PARTITION BY (country)
		   DIMENSION BY (product, year)
		   MEASURES (sales)
		   RULES (sales['Widget', 2024] = sales['Widget', 2023] + 100)`,

		// MODEL with RETURN UPDATED ROWS
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   RETURN UPDATED ROWS
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with RETURN ALL ROWS
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   RETURN ALL ROWS
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with IGNORE NAV
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   IGNORE NAV
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with KEEP NAV
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   KEEP NAV
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with UNIQUE DIMENSION
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   UNIQUE DIMENSION
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with RULES UPDATE
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES UPDATE (sales['US', 2024] = 1000)`,

		// MODEL with RULES UPSERT ALL
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES UPSERT ALL (sales['US', 2024] = 1000)`,

		// MODEL with AUTOMATIC ORDER
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES AUTOMATIC ORDER (sales['US', 2024] = 1000)`,

		// MODEL with SEQUENTIAL ORDER
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES SEQUENTIAL ORDER (sales['US', 2024] = 1000)`,

		// MODEL with ITERATE
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES ITERATE (100) (sales['US', 2024] = sales['US', 2024] + 1)`,

		// MODEL with ITERATE and UNTIL
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES ITERATE (1000) UNTIL (sales['US', 2024] > 10000)
		   (sales['US', 2024] = sales['US', 2024] * 1.1)`,

		// MODEL with multiple rules
		`SELECT country, year, sales, cost FROM sales_data
		 MODEL
		   DIMENSION BY (country, year)
		   MEASURES (sales, cost)
		   RULES (
		     sales['US', 2024] = sales['US', 2023] * 1.1,
		     cost['US', 2024] = cost['US', 2023] * 0.9
		   )`,

		// MODEL with MAIN model name
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   MAIN my_model
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with REFERENCE model
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   REFERENCE ref_model ON (SELECT country, rate FROM exchange_rates)
		     DIMENSION BY (country)
		     MEASURES (rate)
		   MAIN main_model
		   DIMENSION BY (country, year)
		   MEASURES (sales)
		   RULES (sales['US', 2024] = 1000)`,

		// MODEL with aliased MEASURES and DIMENSION BY
		`SELECT country, year, sales FROM sales_data
		 MODEL
		   DIMENSION BY (country AS c, year AS y)
		   MEASURES (amount AS sales)
		   RULES (sales['US', 2024] = 1000)`,
	}
	for _, sql := range tests {
		t.Run(sql[:40], func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
		})
	}
}

// TestParseSelectSample tests SAMPLE clause parsing.
func TestParseSelectSample(t *testing.T) {
	tests := []string{
		// Basic SAMPLE
		`SELECT * FROM employees SAMPLE (10)`,
		// SAMPLE with SEED
		`SELECT * FROM employees SAMPLE (10) SEED (42)`,
		// SAMPLE BLOCK
		`SELECT * FROM employees SAMPLE BLOCK (5)`,
		// SAMPLE BLOCK with SEED
		`SELECT * FROM employees SAMPLE BLOCK (25) SEED (1234)`,
		// SAMPLE in subquery
		`SELECT * FROM (SELECT * FROM employees SAMPLE (50))`,
		// SAMPLE with alias
		`SELECT * FROM employees SAMPLE (10) e`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
		})
	}
}

// TestParseSelectFlashbackQuery tests flashback query clause parsing.
func TestParseSelectFlashbackQuery(t *testing.T) {
	tests := []string{
		// AS OF SCN
		`SELECT * FROM employees AS OF SCN 12345`,
		// AS OF TIMESTAMP with function call
		`SELECT * FROM employees AS OF TIMESTAMP TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD')`,
		// AS OF TIMESTAMP with simple expression
		`SELECT * FROM employees AS OF TIMESTAMP SYSTIMESTAMP`,
		// VERSIONS BETWEEN SCN
		`SELECT * FROM employees VERSIONS BETWEEN SCN 1000 AND 2000`,
		// VERSIONS BETWEEN TIMESTAMP with simple expressions
		`SELECT * FROM employees VERSIONS BETWEEN TIMESTAMP TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD') AND SYSTIMESTAMP`,
		// Flashback with alias
		`SELECT * FROM employees AS OF SCN 12345 e`,
		// Flashback with WHERE
		`SELECT * FROM employees AS OF SCN 12345 WHERE department_id = 10`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
		})
	}
}

func TestParseAnalyticFunctions(t *testing.T) {
	tests := []string{
		"SELECT COUNT(*) OVER () FROM t",
		"SELECT SUM(salary) OVER (PARTITION BY dept_id) FROM employees",
		"SELECT ROW_NUMBER() OVER (ORDER BY hire_date) FROM employees",
		"SELECT RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) FROM employees",
		"SELECT AVG(salary) OVER (ORDER BY hire_date ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM employees",
		"SELECT SUM(amount) OVER (ORDER BY txn_date RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM txns",
		"SELECT SUM(salary) OVER (ORDER BY hire_date ROWS UNBOUNDED PRECEDING) FROM employees",
		"SELECT AVG(salary) OVER (ORDER BY dept_id GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM employees",
		"SELECT MIN(salary) KEEP (DENSE_RANK FIRST ORDER BY hire_date) FROM employees",
		"SELECT MAX(salary) KEEP (DENSE_RANK LAST ORDER BY hire_date DESC) FROM employees",
		"SELECT dept_id, MIN(salary) KEEP (DENSE_RANK FIRST ORDER BY hire_date) FROM employees GROUP BY dept_id",
		"SELECT MIN(salary) KEEP (DENSE_RANK FIRST ORDER BY hire_date) OVER (PARTITION BY dept_id) FROM employees",
		"SELECT LISTAGG(name, ', ') WITHIN GROUP (ORDER BY name) FROM employees",
		"SELECT dept_id, LISTAGG(name, ', ') WITHIN GROUP (ORDER BY hire_date) FROM employees GROUP BY dept_id",
		"SELECT ROW_NUMBER() OVER (ORDER BY salary), RANK() OVER (ORDER BY salary), DENSE_RANK() OVER (ORDER BY salary) FROM employees",
		"SELECT SUM(salary) OVER (ORDER BY hire_date ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING) FROM employees",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseGroupingExtensions(t *testing.T) {
	tests := []string{
		"SELECT dept_id, SUM(salary) FROM employees GROUP BY ROLLUP(dept_id)",
		"SELECT dept_id, job_id, SUM(salary) FROM employees GROUP BY ROLLUP(dept_id, job_id)",
		"SELECT dept_id, SUM(salary) FROM employees GROUP BY CUBE(dept_id, job_id)",
		"SELECT dept_id, SUM(salary) FROM employees GROUP BY GROUPING SETS(dept_id, job_id)",
		"SELECT dept_id, SUM(salary) FROM employees GROUP BY GROUPING SETS(ROLLUP(dept_id), CUBE(job_id))",
		"SELECT dept_id, SUM(salary) FROM employees GROUP BY dept_id, ROLLUP(job_id)",
		"SELECT dept_id, SUM(salary) FROM employees GROUP BY CUBE(dept_id, job_id, mgr_id)",
		"SELECT dept_id, GROUPING(dept_id), SUM(salary) FROM employees GROUP BY ROLLUP(dept_id)",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseSelectFullReview tests all BNF branches from batch 88 SELECT review.
func TestParseSelectFullReview(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// WINDOW clause
		{"window_basic", `SELECT emp_id, salary, SUM(salary) OVER w FROM employees WINDOW w AS (PARTITION BY dept_id ORDER BY salary)`},
		{"window_multiple", `SELECT emp_id, SUM(salary) OVER w1, AVG(salary) OVER w2 FROM employees WINDOW w1 AS (PARTITION BY dept_id), w2 AS (ORDER BY hire_date)`},
		{"window_with_frame", `SELECT emp_id, SUM(salary) OVER w FROM employees WINDOW w AS (ORDER BY hire_date ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)`},
		{"window_ref_existing", `SELECT emp_id, SUM(salary) OVER w2 FROM employees WINDOW w1 AS (PARTITION BY dept_id), w2 AS (w1 ORDER BY salary)`},

		// QUALIFY clause
		{"qualify_basic", `SELECT emp_id, salary, ROW_NUMBER() OVER (ORDER BY salary DESC) AS rn FROM employees QUALIFY rn <= 10`},
		{"qualify_with_func", `SELECT * FROM employees QUALIFY ROW_NUMBER() OVER (PARTITION BY dept_id ORDER BY salary DESC) = 1`},

		// ORDER SIBLINGS BY
		{"order_siblings_by", `SELECT emp_id, mgr_id, last_name FROM employees START WITH mgr_id IS NULL CONNECT BY PRIOR emp_id = mgr_id ORDER SIBLINGS BY last_name`},
		{"order_siblings_by_multi", `SELECT emp_id, last_name FROM employees CONNECT BY PRIOR emp_id = mgr_id START WITH mgr_id IS NULL ORDER SIBLINGS BY last_name ASC, emp_id DESC`},

		// WITH clause SEARCH
		{"with_search_breadth", `WITH RECURSIVE emp_tree (emp_id, mgr_id, lvl) AS (SELECT emp_id, mgr_id, 1 FROM employees WHERE mgr_id IS NULL UNION ALL SELECT e.emp_id, e.mgr_id, t.lvl + 1 FROM employees e JOIN emp_tree t ON e.mgr_id = t.emp_id) SEARCH BREADTH FIRST BY emp_id SET order_col SELECT * FROM emp_tree ORDER BY order_col`},
		{"with_search_depth", `WITH RECURSIVE emp_tree (emp_id, mgr_id) AS (SELECT emp_id, mgr_id FROM employees WHERE mgr_id IS NULL UNION ALL SELECT e.emp_id, e.mgr_id FROM employees e JOIN emp_tree t ON e.mgr_id = t.emp_id) SEARCH DEPTH FIRST BY emp_id SET order_col SELECT * FROM emp_tree`},

		// WITH clause CYCLE
		{"with_cycle", `WITH RECURSIVE emp_tree (emp_id, mgr_id) AS (SELECT emp_id, mgr_id FROM employees WHERE mgr_id IS NULL UNION ALL SELECT e.emp_id, e.mgr_id FROM employees e JOIN emp_tree t ON e.mgr_id = t.emp_id) CYCLE emp_id SET is_cycle TO 'Y' DEFAULT 'N' SELECT * FROM emp_tree`},

		// Partition extension clause
		{"partition_ext_name", `SELECT * FROM sales PARTITION (sales_q1_2024)`},
		{"partition_ext_for", `SELECT * FROM sales PARTITION FOR (TO_DATE('2024-01-15', 'YYYY-MM-DD'))`},
		{"subpartition_ext", `SELECT * FROM sales SUBPARTITION (sales_q1_east)`},
		{"subpartition_ext_for", `SELECT * FROM sales SUBPARTITION FOR ('EAST', TO_DATE('2024-01-15', 'YYYY-MM-DD'))`},

		// Dblink
		{"dblink_basic", `SELECT * FROM employees@remote_db`},
		{"dblink_with_schema", `SELECT * FROM hr.employees@remote_db WHERE dept_id = 10`},

		// TABLE() collection expression
		{"table_collection", `SELECT * FROM TABLE(my_function(1, 2))`},
		{"table_collection_alias", `SELECT * FROM TABLE(get_employees(10)) t`},

		// PIVOT XML
		{"pivot_xml", `SELECT * FROM sales PIVOT XML (SUM(amount) FOR quarter IN (SELECT DISTINCT quarter FROM all_quarters))`},

		// CROSS APPLY / OUTER APPLY
		{"cross_apply", `SELECT e.emp_id, d.dept_name FROM employees e CROSS APPLY (SELECT dept_name FROM departments WHERE dept_id = e.dept_id) d`},
		{"outer_apply", `SELECT e.emp_id, d.dept_name FROM employees e OUTER APPLY (SELECT dept_name FROM departments WHERE dept_id = e.dept_id) d`},

		// VERSIONS PERIOD FOR / AS OF PERIOD FOR
		{"versions_period_for", `SELECT * FROM employees VERSIONS PERIOD FOR valid_time BETWEEN TIMESTAMP '2024-01-01 00:00:00' AND TIMESTAMP '2024-06-30 23:59:59'`},
		{"as_of_period_for", `SELECT * FROM employees AS OF PERIOD FOR valid_time TIMESTAMP '2024-01-01 00:00:00'`},

		// CONTAINERS / SHARDS
		{"containers", `SELECT * FROM CONTAINERS(hr.employees)`},
		{"shards", `SELECT * FROM SHARDS(hr.orders)`},

		// Basic set operations
		{"union", `SELECT a FROM t1 UNION SELECT b FROM t2`},
		{"union_all", `SELECT a FROM t1 UNION ALL SELECT b FROM t2`},
		{"intersect", `SELECT a FROM t1 INTERSECT SELECT b FROM t2`},
		{"minus", `SELECT a FROM t1 MINUS SELECT b FROM t2`},

		// FETCH FIRST
		{"offset_fetch", `SELECT * FROM employees ORDER BY salary DESC OFFSET 10 ROWS FETCH FIRST 5 ROWS ONLY`},
		{"fetch_percent", `SELECT * FROM employees ORDER BY salary DESC FETCH FIRST 10 PERCENT ROWS WITH TIES`},
		{"fetch_next", `SELECT * FROM employees FETCH NEXT 20 ROWS ONLY`},

		// FOR UPDATE variants
		{"for_update_basic", `SELECT * FROM employees FOR UPDATE`},
		{"for_update_of", `SELECT * FROM employees FOR UPDATE OF salary, commission`},
		{"for_update_nowait", `SELECT * FROM employees FOR UPDATE NOWAIT`},
		{"for_update_wait", `SELECT * FROM employees FOR UPDATE WAIT 5`},
		{"for_update_skip_locked", `SELECT * FROM employees FOR UPDATE SKIP LOCKED`},

		// WITH clause basic
		{"with_basic", `WITH dept_avg AS (SELECT dept_id, AVG(salary) avg_sal FROM employees GROUP BY dept_id) SELECT e.emp_id FROM employees e JOIN dept_avg d ON e.dept_id = d.dept_id`},
		{"with_columns", `WITH dept_info(did, avg_sal) AS (SELECT dept_id, AVG(salary) FROM employees GROUP BY dept_id) SELECT * FROM dept_info`},

		// Hierarchical query
		{"connect_by_basic", `SELECT emp_id, mgr_id, LEVEL FROM employees CONNECT BY PRIOR emp_id = mgr_id START WITH mgr_id IS NULL`},
		{"start_with_first", `SELECT emp_id FROM employees START WITH emp_id = 100 CONNECT BY PRIOR emp_id = mgr_id`},
		{"connect_by_nocycle", `SELECT emp_id FROM employees CONNECT BY NOCYCLE PRIOR emp_id = mgr_id`},

		// JOIN variants
		{"inner_join", `SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id`},
		{"left_outer_join", `SELECT * FROM t1 LEFT OUTER JOIN t2 ON t1.id = t2.id`},
		{"right_join", `SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id`},
		{"full_outer_join", `SELECT * FROM t1 FULL OUTER JOIN t2 ON t1.id = t2.id`},
		{"cross_join", `SELECT * FROM t1 CROSS JOIN t2`},
		{"natural_join", `SELECT * FROM t1 NATURAL JOIN t2`},
		{"natural_left_join", `SELECT * FROM t1 NATURAL LEFT JOIN t2`},
		{"join_using", `SELECT * FROM t1 JOIN t2 USING (id)`},

		// Subquery in FROM
		{"subquery_from", `SELECT * FROM (SELECT emp_id, salary FROM employees WHERE dept_id = 10) sub WHERE sub.salary > 5000`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() < 1 {
				t.Fatalf("expected at least 1 statement, got %d", result.Len())
			}
		})
	}
}

func TestParseLateralView(t *testing.T) {
	tests := []string{
		"SELECT * FROM employees e, LATERAL (SELECT * FROM departments d WHERE d.dept_id = e.dept_id) ld",
		"SELECT * FROM t1, LATERAL (SELECT * FROM t2 WHERE t2.id = t1.id)",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseXmlTable(t *testing.T) {
	tests := []string{
		"SELECT x.* FROM xml_data, XMLTABLE('/root/row' PASSING xml_col COLUMNS id NUMBER PATH 'id', name VARCHAR2(100) PATH 'name') x",
		"SELECT x.* FROM xml_data, XMLTABLE('/root/row' PASSING xml_col COLUMNS seq_no FOR ORDINALITY, val VARCHAR2(100) PATH 'val') x",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseJsonTable(t *testing.T) {
	tests := []string{
		"SELECT jt.* FROM json_data, JSON_TABLE(json_col, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id', name VARCHAR2(100) PATH '$.name')) jt",
		"SELECT jt.* FROM json_data, JSON_TABLE(json_col, '$.rows[*]' COLUMNS (seq_no FOR ORDINALITY, val VARCHAR2(100) PATH '$.val')) jt",
		"SELECT jt.* FROM json_data, JSON_TABLE(json_col, '$' COLUMNS (NESTED PATH '$.items[*]' COLUMNS (item_id NUMBER PATH '$.id'))) jt",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseOraclePlusJoin(t *testing.T) {
	tests := []string{
		"SELECT * FROM t1, t2 WHERE t1.id = t2.id(+)",
		"SELECT * FROM t1, t2 WHERE t1.id(+) = t2.id",
		"SELECT * FROM emp, dept WHERE emp.dept_id = dept.id(+) AND dept.name(+) = 'Sales'",
		"SELECT * FROM t1, t2 WHERE t1.col(+) IS NULL",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCursorExpr(t *testing.T) {
	tests := []string{
		"SELECT CURSOR(SELECT * FROM employees) FROM dual",
		"SELECT CURSOR(SELECT dept_id, name FROM departments WHERE active = 1) FROM dual",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseMultisetOps(t *testing.T) {
	tests := []string{
		"SELECT col1 MULTISET UNION col2 FROM t",
		"SELECT col1 MULTISET UNION ALL col2 FROM t",
		"SELECT col1 MULTISET UNION DISTINCT col2 FROM t",
		"SELECT col1 MULTISET INTERSECT col2 FROM t",
		"SELECT col1 MULTISET EXCEPT col2 FROM t",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseTreatExpr(t *testing.T) {
	tests := []string{
		"SELECT TREAT(val AS NUMBER) FROM t",
		"SELECT TREAT(obj AS employee_type) FROM objects",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseXmlExpressions(t *testing.T) {
	tests := []string{
		// XMLELEMENT
		"SELECT XMLELEMENT(NAME \"employee\", XMLELEMENT(NAME \"name\", e.name)) FROM employees e",
		"SELECT XMLELEMENT(NAME \"emp\", e.name, e.salary) FROM employees e",
		// XMLFOREST
		"SELECT XMLFOREST(e.name, e.salary) FROM employees e",
		// XMLAGG
		"SELECT XMLAGG(XMLELEMENT(NAME \"name\", e.name) ORDER BY e.name) FROM employees e",
		// XMLROOT
		"SELECT XMLROOT(xml_col, VERSION '1.0') FROM xml_data",
		// XMLPARSE
		"SELECT XMLPARSE(CONTENT '<a>test</a>') FROM dual",
		// XMLSERIALIZE
		"SELECT XMLSERIALIZE(CONTENT xml_col AS VARCHAR2(4000)) FROM xml_data",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestNodeToString(t *testing.T) {
	n := &ast.NumberLiteral{Val: "42", Ival: 42, Loc: ast.Loc{Start: 0, End: -1}}
	s := ast.NodeToString(n)
	if s == "" {
		t.Error("NodeToString returned empty string")
	}
}

func TestParseJsonExpressions(t *testing.T) {
	tests := []string{
		// JSON_OBJECT
		"SELECT JSON_OBJECT('name', e.name, 'salary', e.salary) FROM employees e",
		// JSON_ARRAY
		"SELECT JSON_ARRAY(1, 2, 3) FROM dual",
		// JSON_VALUE
		"SELECT JSON_VALUE(json_col, '$.name') FROM json_data",
		// JSON_VALUE with RETURNING
		"SELECT JSON_VALUE(json_col, '$.id' RETURNING NUMBER) FROM json_data",
		// JSON_QUERY
		"SELECT JSON_QUERY(json_col, '$.address') FROM json_data",
		// JSON_EXISTS
		"SELECT * FROM json_data WHERE JSON_EXISTS(json_col, '$.name')",
		// JSON_MERGEPATCH
		"SELECT JSON_MERGEPATCH(json_col, '{}') FROM json_data",
		// IS JSON condition
		"SELECT * FROM t WHERE col IS JSON",
		// IS NOT JSON condition
		"SELECT * FROM t WHERE col IS NOT JSON",
	}

	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseLockTable(t *testing.T) {
	tests := []string{
		"LOCK TABLE employees IN ROW SHARE MODE",
		"LOCK TABLE employees IN EXCLUSIVE MODE NOWAIT",
		"LOCK TABLE hr.employees IN SHARE MODE WAIT 10",
		"LOCK TABLE t IN ROW EXCLUSIVE MODE",
		"LOCK TABLE t IN SHARE ROW EXCLUSIVE MODE",
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCallStmt(t *testing.T) {
	tests := []string{
		"CALL my_proc()",
		"CALL my_proc(1, 'hello')",
		"CALL hr.my_func(42) INTO :result",
		"CALL dbms_output.put_line('test')",
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRenameStmt(t *testing.T) {
	tests := []string{
		"RENAME old_table TO new_table",
		"RENAME emp TO employees",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSetRole(t *testing.T) {
	tests := []string{
		"SET ROLE dba_role",
		"SET ROLE role1, role2",
		"SET ROLE ALL",
		"SET ROLE ALL EXCEPT restricted_role",
		"SET ROLE NONE",
		"SET ROLE admin_role IDENTIFIED BY secret",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSetConstraints(t *testing.T) {
	tests := []string{
		"SET CONSTRAINTS ALL IMMEDIATE",
		"SET CONSTRAINTS ALL DEFERRED",
		"SET CONSTRAINT emp_fk IMMEDIATE",
		"SET CONSTRAINTS emp_fk, dept_fk DEFERRED",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAudit(t *testing.T) {
	tests := []string{
		"AUDIT SELECT ON hr.employees",
		"AUDIT INSERT, UPDATE, DELETE ON hr.employees BY ACCESS",
		"AUDIT SELECT ON hr.employees WHENEVER SUCCESSFUL",
		"AUDIT SELECT ON hr.employees WHENEVER NOT SUCCESSFUL",
		"AUDIT ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseNoaudit(t *testing.T) {
	tests := []string{
		"NOAUDIT SELECT ON hr.employees",
		"NOAUDIT ALL",
		"NOAUDIT INSERT ON hr.employees WHENEVER SUCCESSFUL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAssociateStatistics(t *testing.T) {
	tests := []string{
		"ASSOCIATE STATISTICS WITH FUNCTIONS my_func USING my_stats_type",
		"ASSOCIATE STATISTICS WITH COLUMNS employees.salary USING my_stats",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDisassociateStatistics(t *testing.T) {
	tests := []string{
		"DISASSOCIATE STATISTICS FROM FUNCTIONS my_func",
		"DISASSOCIATE STATISTICS FROM COLUMNS employees.salary FORCE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateUser(t *testing.T) {
	tests := []string{
		"CREATE USER scott IDENTIFIED BY tiger",
		"CREATE USER app_user IDENTIFIED BY password123",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterUser(t *testing.T) {
	tests := []string{
		"ALTER USER scott IDENTIFIED BY newpass",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDropUser(t *testing.T) {
	tests := []string{
		"DROP USER scott",
		"DROP USER scott CASCADE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateRole(t *testing.T) {
	tests := []string{
		"CREATE ROLE admin_role",
		"CREATE ROLE secure_role IDENTIFIED BY secret",
		"CREATE ROLE open_role NOT IDENTIFIED",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDropRole(t *testing.T) {
	tests := []string{
		"DROP ROLE admin_role",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateProfile(t *testing.T) {
	tests := []string{
		"CREATE PROFILE app_profile LIMIT SESSIONS_PER_USER 5",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTablespace(t *testing.T) {
	tests := []string{
		"CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M",
		"DROP TABLESPACE temp_ts",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateDirectory(t *testing.T) {
	tests := []string{
		"CREATE DIRECTORY data_dir AS '/u01/data'",
		"DROP DIRECTORY data_dir",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateContext(t *testing.T) {
	tests := []string{
		"CREATE CONTEXT app_ctx USING ctx_pkg",
		"DROP CONTEXT app_ctx",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateCluster(t *testing.T) {
	tests := []string{
		"CREATE CLUSTER emp_dept (dept_id NUMBER)",
		"DROP CLUSTER emp_dept",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateDimension(t *testing.T) {
	tests := []string{
		"CREATE DIMENSION time_dim LEVEL day IS t.day",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateFlashbackArchive(t *testing.T) {
	tests := []string{
		"CREATE FLASHBACK ARCHIVE default_archive TABLESPACE ts1 RETENTION 1 YEAR",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateJava(t *testing.T) {
	tests := []string{
		"CREATE JAVA SOURCE NAMED my_java AS 'public class Foo {}'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateLibrary(t *testing.T) {
	tests := []string{
		"CREATE LIBRARY my_lib AS '/usr/lib/mylib.so'",
		"DROP LIBRARY my_lib",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateMviewLog(t *testing.T) {
	tests := []string{
		"CREATE MATERIALIZED VIEW LOG ON employees",
		"DROP MATERIALIZED VIEW LOG ON employees",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterProcedure(t *testing.T) {
	tests := []string{
		// ALTER PROCEDURE
		"ALTER PROCEDURE my_proc COMPILE",
		"ALTER PROCEDURE hr.my_proc COMPILE DEBUG",
		"ALTER PROCEDURE my_proc COMPILE REUSE SETTINGS",
		"ALTER PROCEDURE my_proc COMPILE DEBUG REUSE SETTINGS",
		"ALTER PROCEDURE IF EXISTS my_proc COMPILE",
		"ALTER PROCEDURE my_proc EDITIONABLE",
		"ALTER PROCEDURE my_proc NONEDITIONABLE",
		"ALTER PROCEDURE my_proc COMPILE PLSQL_OPTIMIZE_LEVEL = 2 REUSE SETTINGS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterFunction(t *testing.T) {
	tests := []string{
		"ALTER FUNCTION my_func COMPILE",
		"ALTER FUNCTION hr.my_func COMPILE DEBUG",
		"ALTER FUNCTION my_func COMPILE REUSE SETTINGS",
		"ALTER FUNCTION IF EXISTS my_func COMPILE DEBUG REUSE SETTINGS",
		"ALTER FUNCTION my_func EDITIONABLE",
		"ALTER FUNCTION my_func NONEDITIONABLE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterPackage(t *testing.T) {
	tests := []string{
		"ALTER PACKAGE my_pkg COMPILE",
		"ALTER PACKAGE hr.my_pkg COMPILE BODY",
		"ALTER PACKAGE my_pkg COMPILE SPECIFICATION",
		"ALTER PACKAGE my_pkg COMPILE PACKAGE",
		"ALTER PACKAGE my_pkg COMPILE BODY DEBUG",
		"ALTER PACKAGE my_pkg COMPILE BODY REUSE SETTINGS",
		"ALTER PACKAGE my_pkg COMPILE BODY DEBUG REUSE SETTINGS",
		"ALTER PACKAGE my_pkg EDITIONABLE",
		"ALTER PACKAGE my_pkg NONEDITIONABLE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterTrigger(t *testing.T) {
	tests := []string{
		"ALTER TRIGGER my_trigger ENABLE",
		"ALTER TRIGGER my_trigger DISABLE",
		"ALTER TRIGGER hr.my_trigger ENABLE",
		"ALTER TRIGGER my_trigger COMPILE",
		"ALTER TRIGGER my_trigger COMPILE DEBUG",
		"ALTER TRIGGER my_trigger COMPILE DEBUG REUSE SETTINGS",
		"ALTER TRIGGER my_trigger RENAME TO new_trigger_name",
		"ALTER TRIGGER IF EXISTS my_trigger ENABLE",
		"ALTER TRIGGER my_trigger EDITIONABLE",
		"ALTER TRIGGER my_trigger NONEDITIONABLE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterProcedureLegacy(t *testing.T) {
	// Legacy tests that should still pass
	tests := []string{
		"ALTER TYPE my_type COMPILE",
		"ALTER MATERIALIZED VIEW mv1 COMPILE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLExit(t *testing.T) {
	tests := []string{
		"BEGIN EXIT; END;",
		"BEGIN EXIT WHEN x > 10; END;",
		"BEGIN EXIT outer_loop; END;",
		"BEGIN EXIT outer_loop WHEN done; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLContinue(t *testing.T) {
	tests := []string{
		"BEGIN CONTINUE; END;",
		"BEGIN CONTINUE WHEN x < 5; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLForall(t *testing.T) {
	tests := []string{
		"BEGIN FORALL i IN 1..tab.COUNT INSERT INTO t VALUES (tab(i)); END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLPipeRow(t *testing.T) {
	tests := []string{
		"BEGIN PIPE ROW (out_rec); END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLPragma(t *testing.T) {
	tests := []string{
		"DECLARE PRAGMA AUTONOMOUS_TRANSACTION; BEGIN NULL; END;",
		"DECLARE PRAGMA EXCEPTION_INIT(e_custom, -20001); BEGIN NULL; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLCaseStmt(t *testing.T) {
	tests := []string{
		"BEGIN CASE x WHEN 1 THEN y := 10; WHEN 2 THEN y := 20; ELSE y := 0; END CASE; END;",
		"BEGIN CASE WHEN x > 0 THEN y := 1; ELSE y := 0; END CASE; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLTypeDecl(t *testing.T) {
	tests := []string{
		"DECLARE TYPE t_tab IS TABLE OF VARCHAR2(100); BEGIN NULL; END;",
		"DECLARE TYPE t_arr IS VARRAY(10) OF NUMBER; BEGIN NULL; END;",
		"DECLARE TYPE t_cur IS REF CURSOR; BEGIN NULL; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePLSQLCall(t *testing.T) {
	tests := []string{
		"BEGIN dbms_output.put_line('hello'); END;",
		"BEGIN my_proc(1, 2); END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 51: CREATE TABLE partitioning
// ---------------------------------------------------------------------------

func TestParseCreateTablePartition(t *testing.T) {
	tests := []string{
		"CREATE TABLE sales (id NUMBER, sale_date DATE) PARTITION BY RANGE (sale_date) (PARTITION p1 VALUES LESS THAN (1000), PARTITION p2 VALUES LESS THAN (MAXVALUE));",
		"CREATE TABLE regions (id NUMBER, region VARCHAR2(50)) PARTITION BY LIST (region) (PARTITION p_east VALUES ('East'), PARTITION p_west VALUES ('West'));",
		"CREATE TABLE logs (id NUMBER) PARTITION BY HASH (id) (PARTITION p1, PARTITION p2, PARTITION p3);",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 52: ALTER INDEX/VIEW/SEQUENCE full
// ---------------------------------------------------------------------------

func TestParseAlterIndexFull(t *testing.T) {
	tests := []string{
		"ALTER INDEX idx_emp REBUILD;",
		"ALTER INDEX hr.idx_emp RENAME TO idx_emp_new;",
		"ALTER INDEX idx_emp MONITORING USAGE;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterViewFull(t *testing.T) {
	tests := []string{
		"ALTER VIEW emp_view COMPILE;",
		"ALTER VIEW hr.emp_view ADD CONSTRAINT pk_view PRIMARY KEY (id) DISABLE NOVALIDATE;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterSequenceFull(t *testing.T) {
	tests := []string{
		"ALTER SEQUENCE emp_seq INCREMENT BY 5 MAXVALUE 10000;",
		"ALTER SEQUENCE hr.emp_seq NOCACHE NOCYCLE;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 53: CREATE SCHEMA
// ---------------------------------------------------------------------------

func TestParseCreateSchema(t *testing.T) {
	tests := []string{
		"CREATE SCHEMA AUTHORIZATION hr;",
		"CREATE SCHEMA AUTHORIZATION hr CREATE TABLE t1 (id NUMBER);",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 85: CREATE SCHEMA (proper)
// ---------------------------------------------------------------------------

func TestParseCreateSchemaProper(t *testing.T) {
	tests := []string{
		// Basic: authorization only
		"CREATE SCHEMA AUTHORIZATION hr;",
		// With a single CREATE TABLE
		"CREATE SCHEMA AUTHORIZATION hr CREATE TABLE t1 (id NUMBER);",
		// With CREATE TABLE and CREATE VIEW
		"CREATE SCHEMA AUTHORIZATION hr CREATE TABLE t1 (id NUMBER) CREATE VIEW v1 AS SELECT id FROM t1;",
		// With CREATE TABLE, CREATE VIEW, and GRANT
		"CREATE SCHEMA AUTHORIZATION hr CREATE TABLE t1 (id NUMBER, name VARCHAR2(100)) CREATE VIEW v1 AS SELECT id, name FROM t1 GRANT SELECT ON t1 TO scott;",
		// With multiple tables and grants
		"CREATE SCHEMA AUTHORIZATION sales CREATE TABLE orders (order_id NUMBER, customer_id NUMBER) CREATE TABLE items (item_id NUMBER, order_id NUMBER) GRANT SELECT ON orders TO hr GRANT SELECT ON items TO hr;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 54: FLASHBACK DATABASE
// ---------------------------------------------------------------------------

func TestParseFlashbackDatabase(t *testing.T) {
	tests := []string{
		"FLASHBACK DATABASE TO SCN 12345;",
		"FLASHBACK DATABASE TO TIMESTAMP SYSTIMESTAMP - 1;",
		"FLASHBACK DATABASE TO RESTORE POINT before_upgrade;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 55: EXISTS subquery + scalar subquery in parens
// ---------------------------------------------------------------------------

func TestParseExistsSubquery(t *testing.T) {
	tests := []string{
		"SELECT 1 FROM dual WHERE EXISTS (SELECT 1 FROM emp);",
		"SELECT 1 FROM dual WHERE NOT EXISTS (SELECT id FROM dept WHERE dept.id = 1);",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseScalarSubquery(t *testing.T) {
	tests := []string{
		"SELECT (SELECT MAX(salary) FROM emp) AS max_sal FROM dual;",
		"SELECT e.name, (SELECT d.name FROM dept d WHERE d.id = e.dept_id) FROM emp e;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 56: Compound trigger + DDL trigger
// ---------------------------------------------------------------------------

func TestParseCompoundTrigger(t *testing.T) {
	tests := []string{
		`CREATE OR REPLACE TRIGGER trg_audit
  FOR INSERT OR UPDATE ON employees
  COMPOUND TRIGGER
  BEFORE STATEMENT IS BEGIN NULL; END BEFORE STATEMENT;
  AFTER EACH ROW IS BEGIN NULL; END AFTER EACH ROW;
  END trg_audit;`,
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDDLTrigger(t *testing.T) {
	tests := []string{
		"CREATE OR REPLACE TRIGGER trg_ddl AFTER CREATE ON DATABASE BEGIN NULL; END;",
		"CREATE TRIGGER trg_logon AFTER LOGON ON DATABASE BEGIN NULL; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTypeBodyMethods(t *testing.T) {
	tests := []string{
		// MEMBER PROCEDURE
		`CREATE OR REPLACE TYPE BODY employee_type AS
  MEMBER PROCEDURE set_name(p_name VARCHAR2) IS
  BEGIN
    SELF.name := p_name;
  END set_name;
END;`,
		// MEMBER FUNCTION
		`CREATE TYPE BODY my_type AS
  MEMBER FUNCTION get_name RETURN VARCHAR2 IS
  BEGIN
    RETURN SELF.name;
  END get_name;
END;`,
		// STATIC FUNCTION
		`CREATE TYPE BODY my_type AS
  STATIC FUNCTION create_default RETURN my_type IS
  BEGIN
    RETURN my_type('default');
  END create_default;
END;`,
		// STATIC PROCEDURE
		`CREATE TYPE BODY my_type AS
  STATIC PROCEDURE log_count IS
  BEGIN
    NULL;
  END log_count;
END;`,
		// MAP MEMBER FUNCTION
		`CREATE TYPE BODY rational_type AS
  MAP MEMBER FUNCTION to_real RETURN NUMBER IS
  BEGIN
    RETURN num / den;
  END to_real;
END;`,
		// ORDER MEMBER FUNCTION
		`CREATE TYPE BODY my_type AS
  ORDER MEMBER FUNCTION compare(other my_type) RETURN INTEGER IS
  BEGIN
    IF SELF.val < other.val THEN RETURN -1;
    ELSIF SELF.val > other.val THEN RETURN 1;
    ELSE RETURN 0;
    END IF;
  END compare;
END;`,
		// Multiple members
		`CREATE OR REPLACE TYPE BODY person_type AS
  MEMBER FUNCTION get_name RETURN VARCHAR2 IS
  BEGIN
    RETURN first_name;
  END get_name;
  MEMBER PROCEDURE set_name(p_name VARCHAR2) IS
  BEGIN
    first_name := p_name;
  END set_name;
  STATIC FUNCTION count_all RETURN NUMBER IS
    v_count NUMBER;
  BEGIN
    SELECT COUNT(*) INTO v_count FROM persons;
    RETURN v_count;
  END count_all;
  MAP MEMBER FUNCTION to_string RETURN VARCHAR2 IS
  BEGIN
    RETURN first_name;
  END to_string;
END;`,
		// Member with DECLARE section (local variables)
		`CREATE TYPE BODY my_type AS
  MEMBER FUNCTION compute RETURN NUMBER IS
    v_result NUMBER;
  BEGIN
    v_result := val * 2;
    RETURN v_result;
  END compute;
END;`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw, ok := result.Items[0].(*ast.RawStmt)
			if !ok {
				t.Fatalf("expected *RawStmt, got %T", result.Items[0])
			}
			stmt, ok := raw.Stmt.(*ast.CreateTypeStmt)
			if !ok {
				t.Fatalf("expected *CreateTypeStmt, got %T", raw.Stmt)
			}
			if !stmt.IsBody {
				t.Fatal("expected IsBody to be true")
			}
			if stmt.Body == nil || len(stmt.Body.Items) == 0 {
				t.Fatal("expected non-empty Body")
			}
		})
	}
}

func TestParseCreateTypeBodyConstructor(t *testing.T) {
	tests := []string{
		// Basic constructor
		`CREATE TYPE BODY my_type AS
  CONSTRUCTOR FUNCTION my_type(p_val NUMBER) RETURN SELF AS RESULT IS
  BEGIN
    SELF.val := p_val;
    RETURN;
  END;
END;`,
		// Constructor with SELF parameter
		`CREATE TYPE BODY my_type AS
  CONSTRUCTOR FUNCTION my_type(SELF IN OUT NOCOPY my_type, p_name VARCHAR2) RETURN SELF AS RESULT IS
  BEGIN
    SELF.name := p_name;
    RETURN;
  END;
END;`,
		// Constructor plus member methods
		`CREATE OR REPLACE TYPE BODY address_type AS
  CONSTRUCTOR FUNCTION address_type(p_street VARCHAR2, p_city VARCHAR2) RETURN SELF AS RESULT IS
  BEGIN
    SELF.street := p_street;
    SELF.city := p_city;
    RETURN;
  END;
  MEMBER FUNCTION get_full_address RETURN VARCHAR2 IS
  BEGIN
    RETURN street || ', ' || city;
  END get_full_address;
END;`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw, ok := result.Items[0].(*ast.RawStmt)
			if !ok {
				t.Fatalf("expected *RawStmt, got %T", result.Items[0])
			}
			stmt, ok := raw.Stmt.(*ast.CreateTypeStmt)
			if !ok {
				t.Fatalf("expected *CreateTypeStmt, got %T", raw.Stmt)
			}
			if !stmt.IsBody {
				t.Fatal("expected IsBody to be true")
			}
			if stmt.Body == nil || len(stmt.Body.Items) == 0 {
				t.Fatal("expected non-empty Body")
			}
		})
	}
}

// TestParseCreateUserFull tests full CREATE USER parsing with structured options.
func TestParseCreateUserFull(t *testing.T) {
	tests := []string{
		"CREATE USER scott IDENTIFIED BY tiger",
		"CREATE USER ext_user IDENTIFIED EXTERNALLY",
		"CREATE USER ext_user IDENTIFIED EXTERNALLY AS 'CN=scott,O=myco'",
		"CREATE USER global_user IDENTIFIED GLOBALLY",
		"CREATE USER global_user IDENTIFIED GLOBALLY AS 'CN=scott,O=myco'",
		"CREATE USER schema_only NO AUTHENTICATION",
		"CREATE USER IF NOT EXISTS scott IDENTIFIED BY tiger",
		"CREATE USER scott IDENTIFIED BY tiger DEFAULT TABLESPACE users",
		"CREATE USER scott IDENTIFIED BY tiger TEMPORARY TABLESPACE temp",
		"CREATE USER scott IDENTIFIED BY tiger LOCAL TEMPORARY TABLESPACE temp",
		"CREATE USER scott IDENTIFIED BY tiger QUOTA 100M ON users",
		"CREATE USER scott IDENTIFIED BY tiger QUOTA UNLIMITED ON users",
		"CREATE USER scott IDENTIFIED BY tiger QUOTA 100M ON users QUOTA UNLIMITED ON temp",
		"CREATE USER scott IDENTIFIED BY tiger PROFILE app_profile",
		"CREATE USER scott IDENTIFIED BY tiger PASSWORD EXPIRE",
		"CREATE USER scott IDENTIFIED BY tiger ACCOUNT LOCK",
		"CREATE USER scott IDENTIFIED BY tiger ACCOUNT UNLOCK",
		"CREATE USER scott IDENTIFIED BY tiger ENABLE EDITIONS",
		"CREATE USER scott IDENTIFIED BY tiger DEFAULT COLLATION USING_NLS_COMP",
		"CREATE USER c##scott IDENTIFIED BY tiger CONTAINER = ALL",
		"CREATE USER scott IDENTIFIED BY tiger CONTAINER = CURRENT",
		"CREATE USER scott IDENTIFIED BY tiger DEFAULT TABLESPACE users TEMPORARY TABLESPACE temp QUOTA 100M ON users PROFILE default PASSWORD EXPIRE ACCOUNT LOCK",
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			_, ok := raw.Stmt.(*ast.CreateUserStmt)
			if !ok {
				t.Fatalf("expected *CreateUserStmt, got %T", raw.Stmt)
			}
		})
	}
}

// TestParseAlterUserFull tests full ALTER USER parsing with structured options.
func TestParseAlterUserFull(t *testing.T) {
	tests := []string{
		"ALTER USER scott IDENTIFIED BY lion REPLACE tiger",
		"ALTER USER scott IDENTIFIED BY newpass",
		"ALTER USER scott IDENTIFIED EXTERNALLY",
		"ALTER USER scott IDENTIFIED GLOBALLY AS 'CN=scott,O=myco'",
		"ALTER USER scott NO AUTHENTICATION",
		"ALTER USER IF EXISTS scott IDENTIFIED BY tiger",
		"ALTER USER scott DEFAULT TABLESPACE users",
		"ALTER USER scott TEMPORARY TABLESPACE temp",
		"ALTER USER scott LOCAL TEMPORARY TABLESPACE temp",
		"ALTER USER scott QUOTA 50M ON users",
		"ALTER USER scott QUOTA UNLIMITED ON users",
		"ALTER USER scott PROFILE app_profile",
		"ALTER USER scott DEFAULT ROLE connect, resource",
		"ALTER USER scott DEFAULT ROLE ALL",
		"ALTER USER scott DEFAULT ROLE ALL EXCEPT dba",
		"ALTER USER scott DEFAULT ROLE NONE",
		"ALTER USER scott PASSWORD EXPIRE",
		"ALTER USER scott ACCOUNT LOCK",
		"ALTER USER scott ACCOUNT UNLOCK",
		"ALTER USER scott ENABLE EDITIONS",
		"ALTER USER c##scott CONTAINER = ALL",
		"ALTER USER scott IDENTIFIED BY newpass DEFAULT TABLESPACE users QUOTA 100M ON users PROFILE default PASSWORD EXPIRE ACCOUNT UNLOCK",
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			_, ok := raw.Stmt.(*ast.AlterUserStmt)
			if !ok {
				t.Fatalf("expected *AlterUserStmt, got %T", raw.Stmt)
			}
		})
	}
}

// TestParseCreateProfileFull tests full CREATE PROFILE parsing with structured options.
func TestParseCreateProfileFull(t *testing.T) {
	tests := []string{
		"CREATE PROFILE app_profile LIMIT SESSIONS_PER_USER 10",
		"CREATE PROFILE app_profile LIMIT CPU_PER_SESSION 10000",
		"CREATE PROFILE app_profile LIMIT CPU_PER_CALL 3000",
		"CREATE PROFILE app_profile LIMIT CONNECT_TIME 60",
		"CREATE PROFILE app_profile LIMIT IDLE_TIME 30",
		"CREATE PROFILE app_profile LIMIT LOGICAL_READS_PER_SESSION UNLIMITED",
		"CREATE PROFILE app_profile LIMIT LOGICAL_READS_PER_CALL DEFAULT",
		"CREATE PROFILE app_profile LIMIT PRIVATE_SGA 100K",
		"CREATE PROFILE app_profile LIMIT COMPOSITE_LIMIT 5000000",
		"CREATE PROFILE app_profile LIMIT FAILED_LOGIN_ATTEMPTS 5",
		"CREATE PROFILE app_profile LIMIT PASSWORD_LIFE_TIME 90",
		"CREATE PROFILE app_profile LIMIT PASSWORD_REUSE_TIME 365",
		"CREATE PROFILE app_profile LIMIT PASSWORD_REUSE_MAX 10",
		"CREATE PROFILE app_profile LIMIT PASSWORD_LOCK_TIME 1",
		"CREATE PROFILE app_profile LIMIT PASSWORD_GRACE_TIME 7",
		"CREATE PROFILE app_profile LIMIT INACTIVE_ACCOUNT_TIME 90",
		"CREATE PROFILE app_profile LIMIT PASSWORD_VERIFY_FUNCTION verify_func",
		"CREATE PROFILE app_profile LIMIT PASSWORD_VERIFY_FUNCTION NULL",
		"CREATE PROFILE app_profile LIMIT PASSWORD_ROLLOVER_TIME 1",
		"CREATE MANDATORY PROFILE mandatory_prof LIMIT SESSIONS_PER_USER 5",
		"CREATE PROFILE app_profile LIMIT SESSIONS_PER_USER 10 CPU_PER_SESSION UNLIMITED FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LIFE_TIME 60",
		"CREATE PROFILE app_profile LIMIT SESSIONS_PER_USER 10 CONTAINER = ALL",
		"CREATE PROFILE app_profile LIMIT SESSIONS_PER_USER 10 LIMIT FAILED_LOGIN_ATTEMPTS 5",
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			_, ok := raw.Stmt.(*ast.CreateProfileStmt)
			if !ok {
				t.Fatalf("expected *CreateProfileStmt, got %T", raw.Stmt)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 60: admin_ddl_option_parsing
// ---------------------------------------------------------------------------

func TestParseCreateTablespaceFull(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// Basic tablespace with DATAFILE and SIZE
		{"datafile_size", "CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M"},
		// Multiple datafiles
		{"multi_datafile", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 500M, '/u01/data02.dbf' SIZE 500M"},
		// AUTOEXTEND ON
		{"autoextend_on", "CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M AUTOEXTEND ON NEXT 10M MAXSIZE 1G"},
		// AUTOEXTEND OFF
		{"autoextend_off", "CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M AUTOEXTEND OFF"},
		// AUTOEXTEND UNLIMITED
		{"autoextend_unlimited", "CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M AUTOEXTEND ON MAXSIZE UNLIMITED"},
		// LOGGING/NOLOGGING
		{"logging", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M LOGGING"},
		{"nologging", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M NOLOGGING"},
		// ONLINE/OFFLINE
		{"online", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M ONLINE"},
		{"offline", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M OFFLINE"},
		// EXTENT MANAGEMENT LOCAL
		{"extent_local_auto", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M EXTENT MANAGEMENT LOCAL AUTOALLOCATE"},
		{"extent_local_uniform", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M EXTENT MANAGEMENT LOCAL UNIFORM SIZE 1M"},
		// SEGMENT SPACE MANAGEMENT
		{"segment_auto", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M SEGMENT SPACE MANAGEMENT AUTO"},
		{"segment_manual", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M SEGMENT SPACE MANAGEMENT MANUAL"},
		// BIGFILE / SMALLFILE
		{"bigfile", "CREATE BIGFILE TABLESPACE big_ts DATAFILE '/u01/big.dbf' SIZE 10G"},
		{"smallfile", "CREATE SMALLFILE TABLESPACE small_ts DATAFILE '/u01/small.dbf' SIZE 100M"},
		// TEMPORARY tablespace
		{"temporary", "CREATE TEMPORARY TABLESPACE temp_ts TEMPFILE '/u01/temp01.dbf' SIZE 500M"},
		// UNDO tablespace
		{"undo", "CREATE UNDO TABLESPACE undo_ts DATAFILE '/u01/undo01.dbf' SIZE 200M"},
		// BLOCKSIZE
		{"blocksize", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M BLOCKSIZE 8K"},
		// Retention
		{"retention_guarantee", "CREATE UNDO TABLESPACE undo_ts DATAFILE '/u01/undo01.dbf' SIZE 200M RETENTION GUARANTEE"},
		{"retention_noguarantee", "CREATE UNDO TABLESPACE undo_ts DATAFILE '/u01/undo01.dbf' SIZE 200M RETENTION NOGUARANTEE"},
		// DEFAULT COMPRESS
		{"compress", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M DEFAULT COMPRESS"},
		{"nocompress", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 100M DEFAULT NOCOMPRESS"},
		// REUSE
		{"datafile_reuse", "CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M REUSE"},
		// Combined options
		{"combined", "CREATE TABLESPACE data01 DATAFILE '/u01/data01.dbf' SIZE 500M AUTOEXTEND ON NEXT 50M MAXSIZE 2G LOGGING EXTENT MANAGEMENT LOCAL AUTOALLOCATE SEGMENT SPACE MANAGEMENT AUTO"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			_, ok := raw.Stmt.(*ast.CreateTablespaceStmt)
			if !ok {
				t.Fatalf("expected *CreateTablespaceStmt, got %T", raw.Stmt)
			}
		})
	}
}

func TestParseCreateClusterFull(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// Basic indexed cluster
		{"basic", "CREATE CLUSTER personnel (department NUMBER(4))"},
		// With SIZE
		{"with_size", "CREATE CLUSTER personnel (department NUMBER(4)) SIZE 512"},
		// With TABLESPACE
		{"with_tablespace", "CREATE CLUSTER emp_dept (deptno NUMBER(3)) SIZE 600 TABLESPACE users"},
		// Hash cluster with HASHKEYS
		{"hash_basic", "CREATE CLUSTER language (cust_language VARCHAR2(3)) SIZE 512 HASHKEYS 10"},
		// Hash cluster with HASH IS expression
		{"hash_expr", "CREATE CLUSTER address (postal_code NUMBER, country_id CHAR(2)) HASHKEYS 20 HASH IS MOD(postal_code + country_id, 101)"},
		// SINGLE TABLE hash cluster
		{"single_table", "CREATE CLUSTER cust_orders (customer_id NUMBER(6)) SIZE 512 SINGLE TABLE HASHKEYS 100"},
		// INDEX clause explicit
		{"index_explicit", "CREATE CLUSTER my_cluster (id NUMBER) INDEX"},
		// CACHE / NOCACHE
		{"cache", "CREATE CLUSTER my_cluster (id NUMBER) CACHE"},
		{"nocache", "CREATE CLUSTER my_cluster (id NUMBER) NOCACHE"},
		// Physical attributes
		{"pctfree", "CREATE CLUSTER my_cluster (id NUMBER) PCTFREE 20"},
		// SORT column
		{"sort_column", "CREATE CLUSTER sorted_cl (id NUMBER, ts DATE SORT) HASHKEYS 100"},
		// Multiple columns
		{"multi_col", "CREATE CLUSTER mc (a NUMBER, b VARCHAR2(10))"},
		// STORAGE clause (parsed but not deeply)
		{"storage", "CREATE CLUSTER personnel (department NUMBER(4)) SIZE 512 STORAGE (INITIAL 100K NEXT 50K)"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			_, ok := raw.Stmt.(*ast.CreateClusterStmt)
			if !ok {
				t.Fatalf("expected *CreateClusterStmt, got %T", raw.Stmt)
			}
		})
	}
}

func TestParseCreateDimensionFull(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// Basic dimension with single level
		{"single_level", "CREATE DIMENSION time_dim LEVEL day IS (t.day_id)"},
		// Multiple levels
		{"multi_level", `CREATE DIMENSION time_dim
			LEVEL day IS (t.day_id)
			LEVEL month IS (t.month_id)
			LEVEL year IS (t.year_id)`},
		// Hierarchy
		{"hierarchy", `CREATE DIMENSION time_dim
			LEVEL day IS (t.day_id)
			LEVEL month IS (t.month_id)
			LEVEL year IS (t.year_id)
			HIERARCHY time_rollup (
				day CHILD OF month CHILD OF year
			)`},
		// Attribute
		{"attribute", `CREATE DIMENSION time_dim
			LEVEL day IS (t.day_id)
			LEVEL month IS (t.month_id)
			ATTRIBUTE day DETERMINES (t.day_name)`},
		// Extended attribute
		{"extended_attr", `CREATE DIMENSION time_dim
			LEVEL day IS (t.day_id)
			ATTRIBUTE day_info LEVEL day DETERMINES (t.day_name, t.day_desc)`},
		// JOIN KEY
		{"join_key", `CREATE DIMENSION customers_dim
			LEVEL customer IS (customers.cust_id)
			LEVEL city IS (customers.cust_city)
			LEVEL country IS (countries.country_id)
			HIERARCHY geog_rollup (
				customer CHILD OF city CHILD OF country
				JOIN KEY (customers.country_id) REFERENCES country
			)`},
		// SKIP WHEN NULL
		{"skip_when_null", `CREATE DIMENSION customers_dim
			LEVEL customer IS (customers.cust_id)
			LEVEL status IS (customers.cust_marital_status) SKIP WHEN NULL
			LEVEL city IS (customers.cust_city)`},
		// Full example from Oracle docs
		{"full_example", `CREATE DIMENSION customers_dim
			LEVEL customer IS (customers.cust_id)
			LEVEL city IS (customers.cust_city)
			LEVEL state IS (customers.cust_state_province)
			LEVEL country IS (countries.country_id)
			LEVEL subregion IS (countries.country_subregion)
			LEVEL region IS (countries.country_region)
			HIERARCHY geog_rollup (
				customer CHILD OF
				city CHILD OF
				state CHILD OF
				country CHILD OF
				subregion CHILD OF
				region
				JOIN KEY (customers.country_id) REFERENCES country
			)
			ATTRIBUTE customer DETERMINES
			(cust_first_name, cust_last_name, cust_gender,
			 cust_marital_status, cust_year_of_birth,
			 cust_income_level, cust_credit_limit)
			ATTRIBUTE country DETERMINES (countries.country_name)`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			_, ok := raw.Stmt.(*ast.CreateDimensionStmt)
			if !ok {
				t.Fatalf("expected *CreateDimensionStmt, got %T", raw.Stmt)
			}
		})
	}
}

// TestParseCreateDatabase tests CREATE DATABASE statements.
func TestParseCreateDatabase(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		minOpts  int // minimum expected options count
		checkName string // expected database name (empty = no check)
	}{
		{"basic", "CREATE DATABASE mydb", 0, "MYDB"},
		{"no_name", "CREATE DATABASE USER SYS IDENTIFIED BY pass1", 1, ""},
		{"with_user_sys_password", "CREATE DATABASE mydb USER SYS IDENTIFIED BY password", 1, "MYDB"},
		{"with_user_system_password", "CREATE DATABASE mydb USER SYSTEM IDENTIFIED BY syspass", 1, "MYDB"},
		{"with_both_users", "CREATE DATABASE mydb USER SYS IDENTIFIED BY pass1 USER SYSTEM IDENTIFIED BY pass2", 2, "MYDB"},
		{"controlfile_reuse", "CREATE DATABASE mydb CONTROLFILE REUSE", 1, "MYDB"},
		{"with_logfile", "CREATE DATABASE mydb LOGFILE GROUP 1 '/u01/log1.log' SIZE 100M", 1, "MYDB"},
		{"with_logfile_multiple_groups", "CREATE DATABASE mydb LOGFILE GROUP 1 '/u01/log1.log' SIZE 100M, GROUP 2 '/u01/log2.log' SIZE 100M", 1, "MYDB"},
		{"with_maxlogfiles", "CREATE DATABASE mydb MAXLOGFILES 16 MAXLOGMEMBERS 3", 2, "MYDB"},
		{"with_maxloghistory", "CREATE DATABASE mydb MAXLOGHISTORY 100", 1, "MYDB"},
		{"with_maxdatafiles", "CREATE DATABASE mydb MAXDATAFILES 100 MAXINSTANCES 8", 2, "MYDB"},
		{"with_archivelog", "CREATE DATABASE mydb ARCHIVELOG", 1, "MYDB"},
		{"with_noarchivelog", "CREATE DATABASE mydb NOARCHIVELOG", 1, "MYDB"},
		{"with_force_logging", "CREATE DATABASE mydb FORCE LOGGING", 1, "MYDB"},
		{"with_set_standby_nologging", "CREATE DATABASE mydb SET STANDBY NOLOGGING FOR DATA AVAILABILITY", 1, "MYDB"},
		{"with_character_set", "CREATE DATABASE mydb CHARACTER SET AL32UTF8", 1, "MYDB"},
		{"with_national_character_set", "CREATE DATABASE mydb NATIONAL CHARACTER SET AL16UTF16", 1, "MYDB"},
		{"with_set_default_bigfile", "CREATE DATABASE mydb SET DEFAULT BIGFILE TABLESPACE", 1, "MYDB"},
		{"with_set_default_smallfile", "CREATE DATABASE mydb SET DEFAULT SMALLFILE TABLESPACE", 1, "MYDB"},
		{"with_default_tablespace", "CREATE DATABASE mydb DEFAULT TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M", 1, "MYDB"},
		{"with_bigfile_default_tablespace", "CREATE DATABASE mydb BIGFILE DEFAULT TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M", 1, "MYDB"},
		{"with_undo_tablespace", "CREATE DATABASE mydb UNDO TABLESPACE undots DATAFILE '/u01/undo01.dbf' SIZE 200M", 1, "MYDB"},
		{"with_default_temp_tablespace", "CREATE DATABASE mydb DEFAULT TEMPORARY TABLESPACE temp TEMPFILE '/u01/temp01.dbf' SIZE 100M", 1, "MYDB"},
		{"with_datafile", "CREATE DATABASE mydb DATAFILE '/u01/data01.dbf' SIZE 500M AUTOEXTEND ON", 1, "MYDB"},
		{"with_set_time_zone", "CREATE DATABASE mydb SET TIME_ZONE = '+00:00'", 1, "MYDB"},
		{"with_set_time_zone_region", "CREATE DATABASE mydb SET TIME_ZONE = 'US/Eastern'", 1, "MYDB"},
		{"with_enable_pluggable_database", "CREATE DATABASE mydb ENABLE PLUGGABLE DATABASE", 1, "MYDB"},
		{"complex", "CREATE DATABASE proddb USER SYS IDENTIFIED BY oracle USER SYSTEM IDENTIFIED BY manager CONTROLFILE REUSE LOGFILE GROUP 1 '/u01/redo01.log' SIZE 50M MAXLOGFILES 16 MAXLOGMEMBERS 3 MAXDATAFILES 1024 MAXINSTANCES 8 ARCHIVELOG FORCE LOGGING CHARACTER SET AL32UTF8 NATIONAL CHARACTER SET AL16UTF16 DATAFILE '/u01/system01.dbf' SIZE 700M DEFAULT TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 500M DEFAULT TEMPORARY TABLESPACE temp TEMPFILE '/u01/temp01.dbf' SIZE 100M UNDO TABLESPACE undots DATAFILE '/u01/undo01.dbf' SIZE 200M", 15, "PRODDB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "CREATE" {
				t.Errorf("expected action CREATE, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DATABASE {
				t.Errorf("expected object type OBJECT_DATABASE, got %d", stmt.ObjectType)
			}
			if tt.checkName != "" && (stmt.Name == nil || stmt.Name.Name != tt.checkName) {
				name := ""
				if stmt.Name != nil {
					name = stmt.Name.Name
				}
				t.Errorf("expected name %q, got %q", tt.checkName, name)
			}
			optCount := 0
			if stmt.Options != nil {
				optCount = len(stmt.Options.Items)
			}
			if optCount < tt.minOpts {
				t.Errorf("expected at least %d options, got %d", tt.minOpts, optCount)
			}
		})
	}
}

// TestParseAlterDatabase tests ALTER DATABASE statements.
func TestParseAlterDatabase(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"mount", "ALTER DATABASE MOUNT"},
		{"open", "ALTER DATABASE OPEN"},
		{"open_resetlogs", "ALTER DATABASE OPEN RESETLOGS"},
		{"archivelog", "ALTER DATABASE ARCHIVELOG"},
		{"noarchivelog", "ALTER DATABASE NOARCHIVELOG"},
		{"force_logging", "ALTER DATABASE FORCE LOGGING"},
		{"no_force_logging", "ALTER DATABASE NO FORCE LOGGING"},
		{"rename_file", "ALTER DATABASE RENAME FILE '/u01/old.dbf' TO '/u01/new.dbf'"},
		{"backup_controlfile", "ALTER DATABASE BACKUP CONTROLFILE TO '/backup/control01.bkp'"},
		{"backup_controlfile_trace", "ALTER DATABASE BACKUP CONTROLFILE TO TRACE"},
		{"add_logfile", "ALTER DATABASE ADD LOGFILE GROUP 3 '/u01/redo03.log' SIZE 50M"},
		{"drop_logfile", "ALTER DATABASE DROP LOGFILE GROUP 2"},
		{"add_datafile", "ALTER DATABASE ADD DATAFILE '/u01/data02.dbf' SIZE 500M"},
		{"recover", "ALTER DATABASE RECOVER AUTOMATIC"},
		{"set_default_tablespace", "ALTER DATABASE SET DEFAULT TABLESPACE users"},
		{"flashback_on", "ALTER DATABASE FLASHBACK ON"},
		{"flashback_off", "ALTER DATABASE FLASHBACK OFF"},
		{"named", "ALTER DATABASE mydb MOUNT"},
		{"enable_block_change_tracking", "ALTER DATABASE ENABLE BLOCK CHANGE TRACKING"},
		{"standby", "ALTER DATABASE ACTIVATE STANDBY DATABASE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Errorf("expected action ALTER, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DATABASE {
				t.Errorf("expected object type OBJECT_DATABASE, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseCreateControlfile tests CREATE CONTROLFILE statements.
func TestParseCreateControlfile(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic_noresetlogs", "CREATE CONTROLFILE REUSE DATABASE mydb NORESETLOGS NOARCHIVELOG"},
		{"resetlogs", "CREATE CONTROLFILE DATABASE mydb RESETLOGS ARCHIVELOG"},
		{"set_database", "CREATE CONTROLFILE SET DATABASE newdb RESETLOGS"},
		{"with_logfile", "CREATE CONTROLFILE REUSE DATABASE mydb NORESETLOGS LOGFILE GROUP 1 '/u01/log1.log' SIZE 500K, GROUP 2 '/u01/log2.log' SIZE 500K"},
		{"with_datafile", "CREATE CONTROLFILE DATABASE mydb NORESETLOGS DATAFILE '/u01/data1.dbf', '/u01/data2.dbf'"},
		{"with_maxlogfiles", "CREATE CONTROLFILE DATABASE mydb NORESETLOGS MAXLOGFILES 32 MAXLOGMEMBERS 2 MAXDATAFILES 100"},
		{"with_maxinstances", "CREATE CONTROLFILE DATABASE mydb NORESETLOGS MAXINSTANCES 8 MAXLOGHISTORY 449"},
		{"with_character_set", "CREATE CONTROLFILE DATABASE mydb NORESETLOGS CHARACTER SET WE8DEC"},
		{"with_force_logging", "CREATE CONTROLFILE DATABASE mydb NORESETLOGS FORCE LOGGING"},
		{"full_example", "CREATE CONTROLFILE REUSE DATABASE \"demo\" NORESETLOGS NOARCHIVELOG MAXLOGFILES 32 MAXLOGMEMBERS 2 MAXDATAFILES 32 MAXINSTANCES 1 MAXLOGHISTORY 449 LOGFILE GROUP 1 '/path/log1.f' SIZE 500K, GROUP 2 '/path/log2.f' SIZE 500K DATAFILE '/path/db1.f', '/path/db2.dbf' CHARACTER SET WE8DEC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "CREATE" {
				t.Errorf("expected action CREATE, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_CONTROLFILE {
				t.Errorf("expected object type OBJECT_CONTROLFILE, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseAlterDatabaseDictionary tests ALTER DATABASE DICTIONARY statements.
func TestParseAlterDatabaseDictionary(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"encrypt_credentials", "ALTER DATABASE DICTIONARY ENCRYPT CREDENTIALS"},
		{"rekey_credentials", "ALTER DATABASE DICTIONARY REKEY CREDENTIALS"},
		{"delete_credentials_key", "ALTER DATABASE DICTIONARY DELETE CREDENTIALS KEY"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Errorf("expected action ALTER, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DATABASE_DICTIONARY {
				t.Errorf("expected object type OBJECT_DATABASE_DICTIONARY, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseCreateDatabaseLinkStillWorks verifies CREATE DATABASE LINK still works after dispatch changes.
func TestParseCreateDatabaseLinkStillWorks(t *testing.T) {
	result := ParseAndCheck(t, "CREATE DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY pass USING 'srv'")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	_, ok := raw.Stmt.(*ast.CreateDatabaseLinkStmt)
	if !ok {
		t.Fatalf("expected *CreateDatabaseLinkStmt, got %T", raw.Stmt)
	}
}

// TestParseAlterDatabaseLinkStillWorks verifies ALTER DATABASE LINK still works after dispatch changes.
func TestParseAlterDatabaseLinkStillWorks(t *testing.T) {
	result := ParseAndCheck(t, "ALTER DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY pass")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterDatabaseLinkStmt)
	if !ok {
		t.Fatalf("expected *AlterDatabaseLinkStmt, got %T", raw.Stmt)
	}
	if stmt.Name == nil {
		t.Fatal("expected non-nil Name")
	}
}

// TestParseCreateDiskgroup tests CREATE DISKGROUP statements.
func TestParseCreateDiskgroup(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "CREATE DISKGROUP dg1 DISK '/dev/sda1' NAME disk1"},
		{"normal_redundancy", "CREATE DISKGROUP dg1 NORMAL REDUNDANCY DISK '/dev/sda1', '/dev/sdb1'"},
		{"high_redundancy", "CREATE DISKGROUP dg1 HIGH REDUNDANCY DISK '/dev/sda1', '/dev/sdb1', '/dev/sdc1'"},
		{"external_redundancy", "CREATE DISKGROUP dg1 EXTERNAL REDUNDANCY DISK '/dev/sda1'"},
		{"with_failgroup", "CREATE DISKGROUP dg1 NORMAL REDUNDANCY FAILGROUP fg1 DISK '/dev/sda1' FAILGROUP fg2 DISK '/dev/sdb1'"},
		{"with_attribute", "CREATE DISKGROUP dg1 DISK '/dev/sda1' ATTRIBUTE 'compatible.asm' = '19.0'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "CREATE" {
				t.Errorf("expected action CREATE, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DISKGROUP {
				t.Errorf("expected OBJECT_DISKGROUP, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseAlterDiskgroup tests ALTER DISKGROUP statements.
func TestParseAlterDiskgroup(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"add_disk", "ALTER DISKGROUP dg1 ADD DISK '/dev/sdc1' NAME disk3"},
		{"drop_disk", "ALTER DISKGROUP dg1 DROP DISK disk2"},
		{"resize_disk", "ALTER DISKGROUP dg1 RESIZE DISK disk1 SIZE 500M"},
		{"rebalance", "ALTER DISKGROUP dg1 REBALANCE POWER 5"},
		{"mount", "ALTER DISKGROUP dg1 MOUNT"},
		{"dismount", "ALTER DISKGROUP dg1 DISMOUNT"},
		{"check_all", "ALTER DISKGROUP dg1 CHECK ALL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Errorf("expected action ALTER, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DISKGROUP {
				t.Errorf("expected OBJECT_DISKGROUP, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseDropDiskgroup tests DROP DISKGROUP statements.
func TestParseDropDiskgroup(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "DROP DISKGROUP dg1"},
		{"force", "DROP DISKGROUP dg1 FORCE"},
		{"including_contents", "DROP DISKGROUP dg1 INCLUDING CONTENTS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "DROP" {
				t.Errorf("expected action DROP, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DISKGROUP {
				t.Errorf("expected OBJECT_DISKGROUP, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseCreatePDB tests CREATE PLUGGABLE DATABASE statements.
func TestParseCreatePDB(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER admin1 IDENTIFIED BY pass"},
		{"from_seed", "CREATE PLUGGABLE DATABASE pdb1 FROM SEED ADMIN USER pdbadm IDENTIFIED BY pass"},
		{"clone", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1"},
		{"relocate", "CREATE PLUGGABLE DATABASE pdb3 FROM pdb1 RELOCATE"},
		{"using_xml", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml'"},
		{"as_clone", "CREATE PLUGGABLE DATABASE pdb2 AS CLONE USING '/tmp/pdb1.xml'"},
		{"file_name_convert", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER admin1 IDENTIFIED BY pass FILE_NAME_CONVERT = ('/u01/pdb_seed/', '/u01/pdb1/')"},
		{"storage", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER admin1 IDENTIFIED BY pass STORAGE UNLIMITED"},
		{"create_file_dest", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER admin1 IDENTIFIED BY pass CREATE_FILE_DEST = '/u01/oradata/pdb1'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "CREATE" {
				t.Errorf("expected action CREATE, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_PLUGGABLE_DATABASE {
				t.Errorf("expected OBJECT_PLUGGABLE_DATABASE, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseAlterPDB tests ALTER PLUGGABLE DATABASE statements.
func TestParseAlterPDB(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"open", "ALTER PLUGGABLE DATABASE pdb1 OPEN"},
		{"close", "ALTER PLUGGABLE DATABASE pdb1 CLOSE"},
		{"open_all", "ALTER PLUGGABLE DATABASE ALL OPEN"},
		{"close_immediate", "ALTER PLUGGABLE DATABASE pdb1 CLOSE IMMEDIATE"},
		{"save_state", "ALTER PLUGGABLE DATABASE pdb1 SAVE STATE"},
		{"discard_state", "ALTER PLUGGABLE DATABASE pdb1 DISCARD STATE"},
		{"unplug", "ALTER PLUGGABLE DATABASE pdb1 UNPLUG INTO '/tmp/pdb1.xml'"},
		{"default_tablespace", "ALTER PLUGGABLE DATABASE pdb1 DEFAULT TABLESPACE users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Errorf("expected action ALTER, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_PLUGGABLE_DATABASE {
				t.Errorf("expected OBJECT_PLUGGABLE_DATABASE, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseDropPDB tests DROP PLUGGABLE DATABASE statements.
func TestParseDropPDB(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "DROP PLUGGABLE DATABASE pdb1"},
		{"including_datafiles", "DROP PLUGGABLE DATABASE pdb1 INCLUDING DATAFILES"},
		{"keep_datafiles", "DROP PLUGGABLE DATABASE pdb1 KEEP DATAFILES"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "DROP" {
				t.Errorf("expected action DROP, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_PLUGGABLE_DATABASE {
				t.Errorf("expected OBJECT_PLUGGABLE_DATABASE, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParsePluggableDatabaseBatch95 tests CREATE/ALTER/DROP PLUGGABLE DATABASE
// statements with full BNF coverage.
func TestParsePluggableDatabaseBatch95(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// === CREATE PLUGGABLE DATABASE ===
		// create_pdb_from_seed
		{"create_seed_basic", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass"},
		{"create_seed_as_seed", "CREATE PLUGGABLE DATABASE pdb1 AS SEED ADMIN USER pdbadm IDENTIFIED BY pass"},
		{"create_seed_as_app_container", "CREATE PLUGGABLE DATABASE pdb1 AS APPLICATION CONTAINER ADMIN USER pdbadm IDENTIFIED BY pass"},
		{"create_seed_roles", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass ROLES = (dba, connect)"},
		{"create_seed_parallel", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass PARALLEL 4"},
		{"create_seed_default_tablespace", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass DEFAULT TABLESPACE users"},
		{"create_seed_file_name_convert", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass FILE_NAME_CONVERT = ('/u01/seed/', '/u01/pdb1/')"},
		{"create_seed_file_name_convert_none", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass FILE_NAME_CONVERT = NONE"},
		{"create_seed_service_name_convert", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass SERVICE_NAME_CONVERT = ('svc1', 'svc2')"},
		{"create_seed_storage_unlimited", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass STORAGE UNLIMITED"},
		{"create_seed_storage_maxsize", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass STORAGE (MAXSIZE 2G)"},
		{"create_seed_storage_multi", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass STORAGE (MAXSIZE 2G MAX_AUDIT_SIZE 100M MAX_DIAG_SIZE 50M)"},
		{"create_seed_path_prefix", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass PATH_PREFIX = '/u01/pdb1/'"},
		{"create_seed_path_prefix_none", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass PATH_PREFIX = NONE"},
		{"create_seed_tempfile_reuse", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass TEMPFILE REUSE"},
		{"create_seed_user_tablespaces", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass USER_TABLESPACES = ALL"},
		{"create_seed_user_tablespaces_except", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass USER_TABLESPACES = ALL EXCEPT (system, sysaux)"},
		{"create_seed_user_tablespaces_list", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass USER_TABLESPACES = (users, data)"},
		{"create_seed_user_tablespaces_none", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass USER_TABLESPACES = NONE"},
		{"create_seed_standbys_all", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass STANDBYS = ALL"},
		{"create_seed_standbys_none", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass STANDBYS = NONE"},
		{"create_seed_standbys_list", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass STANDBYS = (cdb1, cdb2)"},
		{"create_seed_logging", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass LOGGING"},
		{"create_seed_nologging", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass NOLOGGING"},
		{"create_seed_create_file_dest", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass CREATE_FILE_DEST = '/u01/oradata/pdb1'"},
		{"create_seed_create_file_dest_none", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass CREATE_FILE_DEST = NONE"},
		{"create_seed_snapshot", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass SNAPSHOT = MANUAL"},
		{"create_seed_snapshot_every", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass SNAPSHOT = EVERY 4 HOURS"},
		{"create_seed_snapshot_none", "CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadm IDENTIFIED BY pass SNAPSHOT = NONE"},
		// create_pdb_clone
		{"create_clone_basic", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1"},
		{"create_clone_dblink", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1@remote_cdb"},
		{"create_clone_proxy", "CREATE PLUGGABLE DATABASE pdb2 AS PROXY FROM pdb1@remote_cdb"},
		{"create_clone_proxy_host_port", "CREATE PLUGGABLE DATABASE pdb2 AS PROXY FROM pdb1@remote_cdb HOST 'myhost' PORT 1521"},
		{"create_clone_parallel", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 PARALLEL 4"},
		{"create_clone_file_name_convert", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 FILE_NAME_CONVERT = ('/u01/pdb1/', '/u01/pdb2/')"},
		{"create_clone_snapshot_copy", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 SNAPSHOT COPY"},
		{"create_clone_snapshot_copy_no_data", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 SNAPSHOT COPY NO DATA"},
		{"create_clone_no_data", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 NO DATA"},
		{"create_clone_relocate", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 RELOCATE"},
		{"create_clone_relocate_keep", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 RELOCATE KEEP SOURCE"},
		{"create_clone_keystore", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 KEYSTORE IDENTIFIED BY pass1"},
		{"create_clone_keystore_decrypt", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 KEYSTORE IDENTIFIED BY pass1 DECRYPT USING 'secret'"},
		{"create_clone_refresh_manual", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 REFRESH MODE MANUAL"},
		{"create_clone_refresh_every", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 REFRESH MODE EVERY 30 MINUTES"},
		{"create_clone_refresh_none", "CREATE PLUGGABLE DATABASE pdb2 FROM pdb1 REFRESH MODE NONE"},
		// create_pdb_from_xml
		{"create_xml_basic", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml'"},
		{"create_xml_copy", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' COPY"},
		{"create_xml_move", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' MOVE"},
		{"create_xml_nocopy", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' NOCOPY"},
		{"create_xml_source_file_name_convert", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' SOURCE_FILE_NAME_CONVERT = ('/old/', '/new/')"},
		{"create_xml_source_file_directory", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' SOURCE_FILE_DIRECTORY = '/tmp/files'"},
		{"create_xml_as_clone", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' AS CLONE"},
		{"create_xml_decrypt", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' DECRYPT USING 'secret'"},
		{"create_xml_keystore", "CREATE PLUGGABLE DATABASE pdb1 USING '/tmp/pdb1.xml' KEYSTORE IDENTIFIED BY pass1"},
		// === ALTER PLUGGABLE DATABASE ===
		// pdb_open / pdb_close / pdb_save_or_discard_state
		{"alter_open_read_write", "ALTER PLUGGABLE DATABASE pdb1 OPEN READ WRITE"},
		{"alter_open_read_only", "ALTER PLUGGABLE DATABASE pdb1 OPEN READ ONLY"},
		{"alter_open_restricted", "ALTER PLUGGABLE DATABASE pdb1 OPEN RESTRICTED"},
		{"alter_open_force", "ALTER PLUGGABLE DATABASE pdb1 OPEN FORCE"},
		{"alter_open_resetlogs", "ALTER PLUGGABLE DATABASE pdb1 OPEN RESETLOGS"},
		{"alter_open_upgrade", "ALTER PLUGGABLE DATABASE pdb1 OPEN UPGRADE"},
		{"alter_open_hybrid_read_only", "ALTER PLUGGABLE DATABASE pdb1 OPEN HYBRID READ ONLY"},
		{"alter_close_abort", "ALTER PLUGGABLE DATABASE pdb1 CLOSE ABORT"},
		{"alter_close_relocate", "ALTER PLUGGABLE DATABASE pdb1 CLOSE RELOCATE"},
		{"alter_close_norelocate", "ALTER PLUGGABLE DATABASE pdb1 CLOSE NORELOCATE"},
		{"alter_close_relocate_to", "ALTER PLUGGABLE DATABASE pdb1 CLOSE RELOCATE TO 'inst1'"},
		{"alter_save_state", "ALTER PLUGGABLE DATABASE pdb1 SAVE STATE"},
		{"alter_discard_state", "ALTER PLUGGABLE DATABASE pdb1 DISCARD STATE"},
		{"alter_save_state_instances", "ALTER PLUGGABLE DATABASE pdb1 SAVE STATE INSTANCES = (inst1, inst2)"},
		// pdb_change_state_from_root
		{"alter_all_open", "ALTER PLUGGABLE DATABASE ALL OPEN"},
		{"alter_all_except_open", "ALTER PLUGGABLE DATABASE ALL EXCEPT (pdb1, pdb2) OPEN"},
		{"alter_all_close", "ALTER PLUGGABLE DATABASE ALL CLOSE IMMEDIATE"},
		{"alter_open_instances", "ALTER PLUGGABLE DATABASE pdb1 OPEN INSTANCES = (inst1)"},
		{"alter_open_services", "ALTER PLUGGABLE DATABASE pdb1 OPEN SERVICES = (svc1, svc2)"},
		// pdb_unplug_clause
		{"alter_unplug", "ALTER PLUGGABLE DATABASE pdb1 UNPLUG INTO '/tmp/pdb1.xml'"},
		{"alter_unplug_encrypt", "ALTER PLUGGABLE DATABASE pdb1 UNPLUG INTO '/tmp/pdb1.xml' ENCRYPT USING 'secret'"},
		// pdb_settings_clauses
		{"alter_default_edition", "ALTER PLUGGABLE DATABASE DEFAULT EDITION ora_edition"},
		{"alter_set_default_tablespace_bigfile", "ALTER PLUGGABLE DATABASE SET DEFAULT TABLESPACE BIGFILE"},
		{"alter_set_default_tablespace_smallfile", "ALTER PLUGGABLE DATABASE SET DEFAULT TABLESPACE SMALLFILE"},
		{"alter_default_tablespace", "ALTER PLUGGABLE DATABASE pdb1 DEFAULT TABLESPACE users"},
		{"alter_default_temp_tablespace", "ALTER PLUGGABLE DATABASE pdb1 DEFAULT TEMPORARY TABLESPACE temp"},
		{"alter_rename_global_name", "ALTER PLUGGABLE DATABASE pdb1 RENAME GLOBAL_NAME TO pdb1.example.com"},
		{"alter_set_time_zone", "ALTER PLUGGABLE DATABASE pdb1 SET TIME_ZONE = '+05:30'"},
		{"alter_storage", "ALTER PLUGGABLE DATABASE pdb1 STORAGE (MAXSIZE 10G)"},
		{"alter_storage_unlimited", "ALTER PLUGGABLE DATABASE pdb1 STORAGE UNLIMITED"},
		{"alter_logging", "ALTER PLUGGABLE DATABASE pdb1 LOGGING"},
		{"alter_nologging", "ALTER PLUGGABLE DATABASE pdb1 NOLOGGING"},
		{"alter_enable_force_logging", "ALTER PLUGGABLE DATABASE pdb1 ENABLE FORCE LOGGING"},
		{"alter_disable_force_logging", "ALTER PLUGGABLE DATABASE pdb1 DISABLE FORCE LOGGING"},
		{"alter_refresh_mode_manual", "ALTER PLUGGABLE DATABASE pdb1 REFRESH MODE MANUAL"},
		{"alter_refresh_mode_every", "ALTER PLUGGABLE DATABASE pdb1 REFRESH MODE EVERY 60 MINUTES"},
		{"alter_refresh_mode_none", "ALTER PLUGGABLE DATABASE pdb1 REFRESH MODE NONE"},
		{"alter_refresh_switchover", "ALTER PLUGGABLE DATABASE pdb1 REFRESH SWITCHOVER TO PRIMARY remote_link"},
		{"alter_containers_default_target", "ALTER PLUGGABLE DATABASE CONTAINERS DEFAULT TARGET pdb1"},
		{"alter_containers_host", "ALTER PLUGGABLE DATABASE CONTAINERS HOST = 'myhost'"},
		{"alter_containers_port", "ALTER PLUGGABLE DATABASE CONTAINERS PORT = 1521"},
		{"alter_priority", "ALTER PLUGGABLE DATABASE pdb1 PRIORITY 5"},
		{"alter_priority_none", "ALTER PLUGGABLE DATABASE pdb1 PRIORITY NONE"},
		// pdb_datafile_clause
		{"alter_datafile_online", "ALTER PLUGGABLE DATABASE pdb1 DATAFILE '/u01/pdb1/data01.dbf' ONLINE"},
		{"alter_datafile_offline", "ALTER PLUGGABLE DATABASE pdb1 DATAFILE ALL OFFLINE"},
		{"alter_datafile_number", "ALTER PLUGGABLE DATABASE pdb1 DATAFILE 5 ONLINE"},
		// pdb_recovery_clauses
		{"alter_recover", "ALTER PLUGGABLE DATABASE pdb1 RECOVER"},
		{"alter_recover_until_cancel", "ALTER PLUGGABLE DATABASE pdb1 RECOVER UNTIL CANCEL"},
		{"alter_recover_until_time", "ALTER PLUGGABLE DATABASE pdb1 RECOVER UNTIL TIME '2024-01-01'"},
		{"alter_recover_until_change", "ALTER PLUGGABLE DATABASE pdb1 RECOVER UNTIL CHANGE 12345"},
		{"alter_recover_standby", "ALTER PLUGGABLE DATABASE pdb1 RECOVER STANDBY UNTIL CANCEL"},
		{"alter_recover_backup_controlfile", "ALTER PLUGGABLE DATABASE pdb1 RECOVER USING BACKUP CONTROLFILE"},
		{"alter_backup_begin", "ALTER PLUGGABLE DATABASE pdb1 BACKUP BEGIN"},
		{"alter_backup_end", "ALTER PLUGGABLE DATABASE pdb1 BACKUP END"},
		{"alter_enable_recovery", "ALTER PLUGGABLE DATABASE pdb1 ENABLE RECOVERY"},
		{"alter_disable_recovery", "ALTER PLUGGABLE DATABASE pdb1 DISABLE RECOVERY"},
		// application_clauses
		{"alter_app_begin_install", "ALTER PLUGGABLE DATABASE APPLICATION myapp BEGIN INSTALL VERSION '1.0'"},
		{"alter_app_begin_install_comment", "ALTER PLUGGABLE DATABASE APPLICATION myapp BEGIN INSTALL VERSION '1.0' COMMENT 'initial install'"},
		{"alter_app_begin_patch", "ALTER PLUGGABLE DATABASE APPLICATION myapp BEGIN PATCH 100"},
		{"alter_app_begin_patch_min_ver", "ALTER PLUGGABLE DATABASE APPLICATION myapp BEGIN PATCH 100 MINIMUM VERSION '1.0'"},
		{"alter_app_begin_upgrade", "ALTER PLUGGABLE DATABASE APPLICATION myapp BEGIN UPGRADE FROM '1.0' TO '2.0'"},
		{"alter_app_begin_uninstall", "ALTER PLUGGABLE DATABASE APPLICATION myapp BEGIN UNINSTALL"},
		{"alter_app_end_install", "ALTER PLUGGABLE DATABASE APPLICATION myapp END INSTALL"},
		{"alter_app_end_install_version", "ALTER PLUGGABLE DATABASE APPLICATION myapp END INSTALL VERSION '2.0'"},
		{"alter_app_end_patch", "ALTER PLUGGABLE DATABASE APPLICATION myapp END PATCH"},
		{"alter_app_end_upgrade", "ALTER PLUGGABLE DATABASE APPLICATION myapp END UPGRADE"},
		{"alter_app_end_uninstall", "ALTER PLUGGABLE DATABASE APPLICATION myapp END UNINSTALL"},
		{"alter_app_set_patch", "ALTER PLUGGABLE DATABASE APPLICATION myapp SET PATCH 200"},
		{"alter_app_set_version", "ALTER PLUGGABLE DATABASE APPLICATION myapp SET VERSION '2.0'"},
		{"alter_app_set_compat", "ALTER PLUGGABLE DATABASE APPLICATION myapp SET COMPATIBILITY VERSION '1.0'"},
		{"alter_app_set_compat_current", "ALTER PLUGGABLE DATABASE APPLICATION myapp SET COMPATIBILITY VERSION CURRENT"},
		{"alter_app_sync", "ALTER PLUGGABLE DATABASE APPLICATION myapp SYNC"},
		{"alter_app_sync_to_version", "ALTER PLUGGABLE DATABASE APPLICATION myapp SYNC TO '2.0'"},
		{"alter_app_sync_to_patch", "ALTER PLUGGABLE DATABASE APPLICATION myapp SYNC TO PATCH 100"},
		{"alter_app_multi_sync", "ALTER PLUGGABLE DATABASE APPLICATION app1, app2 SYNC"},
		{"alter_app_all_sync", "ALTER PLUGGABLE DATABASE APPLICATION ALL SYNC"},
		{"alter_app_all_except_sync", "ALTER PLUGGABLE DATABASE APPLICATION ALL EXCEPT (app1) SYNC"},
		// snapshot_clauses
		{"alter_snapshot_none", "ALTER PLUGGABLE DATABASE pdb1 SNAPSHOT NONE"},
		{"alter_snapshot_manual", "ALTER PLUGGABLE DATABASE pdb1 SNAPSHOT MANUAL"},
		{"alter_snapshot_every", "ALTER PLUGGABLE DATABASE pdb1 SNAPSHOT EVERY 4 HOURS"},
		{"alter_materialize", "ALTER PLUGGABLE DATABASE pdb1 MATERIALIZE"},
		{"alter_create_snapshot", "ALTER PLUGGABLE DATABASE pdb1 CREATE SNAPSHOT snap1"},
		{"alter_drop_snapshot", "ALTER PLUGGABLE DATABASE pdb1 DROP SNAPSHOT snap1"},
		{"alter_set_max_snapshots", "ALTER PLUGGABLE DATABASE pdb1 SET MAX_PDB_SNAPSHOTS = 10"},
		{"alter_prepare_mirror", "ALTER PLUGGABLE DATABASE pdb1 PREPARE MIRROR COPY mc1"},
		{"alter_prepare_mirror_redundancy", "ALTER PLUGGABLE DATABASE pdb1 PREPARE MIRROR COPY mc1 WITH NORMAL REDUNDANCY"},
		{"alter_prepare_mirror_for_db", "ALTER PLUGGABLE DATABASE pdb1 PREPARE MIRROR COPY mc1 FOR DATABASE db1"},
		{"alter_drop_mirror_copy", "ALTER PLUGGABLE DATABASE pdb1 DROP MIRROR COPY mc1"},
		{"alter_enable_lost_write", "ALTER PLUGGABLE DATABASE pdb1 ENABLE LOST WRITE PROTECTION"},
		{"alter_disable_lost_write", "ALTER PLUGGABLE DATABASE pdb1 DISABLE LOST WRITE PROTECTION"},
		{"alter_enable_backup", "ALTER PLUGGABLE DATABASE pdb1 ENABLE BACKUP"},
		{"alter_disable_backup", "ALTER PLUGGABLE DATABASE pdb1 DISABLE BACKUP"},
		// === DROP PLUGGABLE DATABASE ===
		{"drop_basic", "DROP PLUGGABLE DATABASE pdb1"},
		{"drop_if_exists", "DROP PLUGGABLE DATABASE IF EXISTS pdb1"},
		{"drop_force", "DROP PLUGGABLE DATABASE pdb1 FORCE"},
		{"drop_including_datafiles", "DROP PLUGGABLE DATABASE pdb1 INCLUDING DATAFILES"},
		{"drop_keep_datafiles", "DROP PLUGGABLE DATABASE pdb1 KEEP DATAFILES"},
		{"drop_force_including", "DROP PLUGGABLE DATABASE pdb1 FORCE INCLUDING DATAFILES"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.ObjectType != ast.OBJECT_PLUGGABLE_DATABASE {
				t.Errorf("expected OBJECT_PLUGGABLE_DATABASE, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseAdministerKeyMgmt tests ADMINISTER KEY MANAGEMENT statements.
func TestParseAdministerKeyMgmt(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"create_keystore", "ADMINISTER KEY MANAGEMENT CREATE KEYSTORE '/u01/keystore' IDENTIFIED BY password1"},
		{"create_auto_login", "ADMINISTER KEY MANAGEMENT CREATE AUTO_LOGIN KEYSTORE FROM KEYSTORE '/u01/keystore' IDENTIFIED BY password1"},
		{"create_local_auto_login", "ADMINISTER KEY MANAGEMENT CREATE LOCAL AUTO_LOGIN KEYSTORE FROM KEYSTORE '/u01/keystore' IDENTIFIED BY password1"},
		{"open_keystore", "ADMINISTER KEY MANAGEMENT SET KEYSTORE OPEN IDENTIFIED BY password1"},
		{"close_keystore", "ADMINISTER KEY MANAGEMENT SET KEYSTORE CLOSE IDENTIFIED BY password1"},
		{"set_key", "ADMINISTER KEY MANAGEMENT SET KEY IDENTIFIED BY password1 WITH BACKUP"},
		{"set_key_tag", "ADMINISTER KEY MANAGEMENT SET KEY USING TAG 'quarterly_key' IDENTIFIED BY password1 WITH BACKUP USING 'Q1 key rotation'"},
		{"create_key", "ADMINISTER KEY MANAGEMENT CREATE KEY IDENTIFIED BY password1 WITH BACKUP"},
		{"use_key", "ADMINISTER KEY MANAGEMENT USE KEY 'key_id_123' IDENTIFIED BY password1 WITH BACKUP"},
		{"backup_keystore", "ADMINISTER KEY MANAGEMENT BACKUP KEYSTORE IDENTIFIED BY password1 TO '/backup/'"},
		{"alter_password", "ADMINISTER KEY MANAGEMENT ALTER KEYSTORE PASSWORD IDENTIFIED BY old_pass SET new_pass WITH BACKUP"},
		{"export_keys", "ADMINISTER KEY MANAGEMENT EXPORT KEYS WITH SECRET 'my_secret' TO '/tmp/export.p12' IDENTIFIED BY password1"},
		{"import_keys", "ADMINISTER KEY MANAGEMENT IMPORT KEYS WITH SECRET 'my_secret' FROM '/tmp/export.p12' IDENTIFIED BY password1 WITH BACKUP"},
		{"merge_keystore", "ADMINISTER KEY MANAGEMENT MERGE KEYSTORE '/u01/ks1' AND '/u01/ks2' IDENTIFIED BY password1 INTO NEW KEYSTORE '/u01/merged' IDENTIFIED BY password2"},
		{"add_secret", "ADMINISTER KEY MANAGEMENT ADD SECRET 'secret1' FOR CLIENT 'client1' IDENTIFIED BY password1"},
		{"set_tag", "ADMINISTER KEY MANAGEMENT SET TAG 'mytag' FOR 'key123' IDENTIFIED BY password1 WITH BACKUP"},
		{"external_store", "ADMINISTER KEY MANAGEMENT SET KEY IDENTIFIED BY EXTERNAL STORE WITH BACKUP"},
		{"container_all", "ADMINISTER KEY MANAGEMENT SET KEY IDENTIFIED BY password1 WITH BACKUP CONTAINER = ALL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ADMINISTER" {
				t.Errorf("expected action ADMINISTER, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_KEY_MANAGEMENT {
				t.Errorf("expected OBJECT_KEY_MANAGEMENT, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseCreateAuditPolicy tests CREATE AUDIT POLICY statements.
func TestParseCreateAuditPolicy(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "CREATE AUDIT POLICY my_policy ACTIONS SELECT ON hr.employees"},
		{"all_actions", "CREATE AUDIT POLICY my_policy ACTIONS ALL"},
		{"multiple_actions", "CREATE AUDIT POLICY my_policy ACTIONS INSERT ON hr.employees, DELETE ON hr.employees"},
		{"privileges", "CREATE AUDIT POLICY priv_policy PRIVILEGES CREATE TABLE, DROP ANY TABLE"},
		{"roles", "CREATE AUDIT POLICY role_policy ROLES dba, resource"},
		{"when_condition", "CREATE AUDIT POLICY cond_policy ACTIONS SELECT ON hr.employees WHEN 'SYS_CONTEXT(''USERENV'', ''SESSION_USER'') = ''HR''' EVALUATE PER SESSION"},
		{"system_actions", "CREATE AUDIT POLICY sys_policy ACTIONS CREATE TABLE, ALTER TABLE, DROP TABLE"},
		{"container_current", "CREATE AUDIT POLICY my_policy ACTIONS SELECT CONTAINER = CURRENT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "CREATE" {
				t.Errorf("expected action CREATE, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_AUDIT_POLICY {
				t.Errorf("expected OBJECT_AUDIT_POLICY, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseAlterAuditPolicy tests ALTER AUDIT POLICY statements.
func TestParseAlterAuditPolicy(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"add_actions", "ALTER AUDIT POLICY my_policy ADD ACTIONS DELETE ON hr.employees"},
		{"drop_actions", "ALTER AUDIT POLICY my_policy DROP ACTIONS SELECT ON hr.employees"},
		{"add_privileges", "ALTER AUDIT POLICY my_policy ADD PRIVILEGES CREATE ANY TABLE"},
		{"add_roles", "ALTER AUDIT POLICY my_policy ADD ROLES connect"},
		{"condition", "ALTER AUDIT POLICY my_policy CONDITION 'SYS_CONTEXT(''USERENV'',''IP_ADDRESS'') = ''10.0.0.1''' EVALUATE PER SESSION"},
		{"drop_condition", "ALTER AUDIT POLICY my_policy CONDITION DROP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Errorf("expected action ALTER, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_AUDIT_POLICY {
				t.Errorf("expected OBJECT_AUDIT_POLICY, got %d", stmt.ObjectType)
			}
		})
	}
}

// TestParseDropAuditPolicy tests DROP AUDIT POLICY statements.
func TestParseDropAuditPolicy(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "DROP AUDIT POLICY my_policy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw := result.Items[0].(*ast.RawStmt)
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "DROP" {
				t.Errorf("expected action DROP, got %q", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_AUDIT_POLICY {
				t.Errorf("expected OBJECT_AUDIT_POLICY, got %d", stmt.ObjectType)
			}
		})
	}
}

// ---------- Batch 66-79 Tests ----------

// adminDDLTest is a helper for testing AdminDDLStmt parsing.
func adminDDLTest(t *testing.T, sql string, wantAction string, wantType ast.ObjectType) {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
	if !ok {
		t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
	}
	if stmt.Action != wantAction {
		t.Errorf("expected action %q, got %q", wantAction, stmt.Action)
	}
	if stmt.ObjectType != wantType {
		t.Errorf("expected object type %d, got %d", wantType, stmt.ObjectType)
	}
}

// TestBatch66_AnalyticViewHierarchy tests CREATE/ALTER/DROP ANALYTIC VIEW, ATTRIBUTE DIMENSION, HIERARCHY.
func TestBatch66_AnalyticViewHierarchy(t *testing.T) {
	// CREATE/ALTER ANALYTIC VIEW now use dedicated AST types (upgraded in batch 93)
	t.Run("create_analytic_view", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE ANALYTIC VIEW sales_av USING sales_fact DIMENSION BY (time_dim) MEASURES (amount)")
		raw := result.Items[0].(*ast.RawStmt)
		_, ok := raw.Stmt.(*ast.CreateAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *CreateAnalyticViewStmt, got %T", raw.Stmt)
		}
	})
	t.Run("alter_analytic_view", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER ANALYTIC VIEW sales_av RENAME TO sales_av2")
		raw := result.Items[0].(*ast.RawStmt)
		_, ok := raw.Stmt.(*ast.AlterAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *AlterAnalyticViewStmt, got %T", raw.Stmt)
		}
	})

	// DROP and ATTRIBUTE DIMENSION / HIERARCHY still use AdminDDLStmt
	adminTests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"drop_analytic_view", "DROP ANALYTIC VIEW sales_av", "DROP", ast.OBJECT_ANALYTIC_VIEW},
		{"create_attribute_dimension", "CREATE ATTRIBUTE DIMENSION time_attr_dim USING time_dim", "CREATE", ast.OBJECT_ATTRIBUTE_DIMENSION},
		{"alter_attribute_dimension", "ALTER ATTRIBUTE DIMENSION time_attr_dim RENAME TO time_attr_dim2", "ALTER", ast.OBJECT_ATTRIBUTE_DIMENSION},
		{"drop_attribute_dimension", "DROP ATTRIBUTE DIMENSION time_attr_dim", "DROP", ast.OBJECT_ATTRIBUTE_DIMENSION},
		{"create_hierarchy", "CREATE HIERARCHY time_hier USING time_attr_dim", "CREATE", ast.OBJECT_HIERARCHY},
		{"alter_hierarchy", "ALTER HIERARCHY time_hier RENAME TO time_hier2", "ALTER", ast.OBJECT_HIERARCHY},
		{"drop_hierarchy", "DROP HIERARCHY time_hier", "DROP", ast.OBJECT_HIERARCHY},
	}
	for _, tt := range adminTests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch67_Domain tests CREATE/ALTER/DROP DOMAIN.
func TestBatch67_Domain(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_domain", "CREATE DOMAIN email_domain AS VARCHAR2(255) NOT NULL", "CREATE", ast.OBJECT_DOMAIN},
		{"alter_domain", "ALTER DOMAIN email_domain ADD CONSTRAINT chk_email CHECK (VALUE LIKE '%@%')", "ALTER", ast.OBJECT_DOMAIN},
		{"drop_domain", "DROP DOMAIN email_domain", "DROP", ast.OBJECT_DOMAIN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch68_IndextypeOperator tests CREATE/ALTER/DROP INDEXTYPE and OPERATOR.
// Updated in batch 92 to use dedicated AST types.
func TestBatch68_IndextypeOperator(t *testing.T) {
	t.Run("create_indextype", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_indextype FOR my_operator(NUMBER) USING my_type")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateIndextypeStmt)
		if !ok {
			t.Fatalf("expected *CreateIndextypeStmt, got %T", raw.Stmt)
		}
		if stmt.Name.Name != "MY_INDEXTYPE" {
			t.Errorf("expected MY_INDEXTYPE, got %q", stmt.Name.Name)
		}
	})
	t.Run("alter_indextype", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_indextype ADD my_operator2(VARCHAR2)")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterIndextypeStmt)
		if !ok {
			t.Fatalf("expected *AlterIndextypeStmt, got %T", raw.Stmt)
		}
		if stmt.Name.Name != "MY_INDEXTYPE" {
			t.Errorf("expected MY_INDEXTYPE, got %q", stmt.Name.Name)
		}
	})
	t.Run("drop_indextype", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEXTYPE my_indextype")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.DropStmt)
		if !ok {
			t.Fatalf("expected *DropStmt, got %T", raw.Stmt)
		}
		if stmt.ObjectType != ast.OBJECT_INDEXTYPE {
			t.Errorf("expected OBJECT_INDEXTYPE, got %d", stmt.ObjectType)
		}
	})
	t.Run("create_operator", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING my_eq_func")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateOperatorStmt)
		if !ok {
			t.Fatalf("expected *CreateOperatorStmt, got %T", raw.Stmt)
		}
		if stmt.Name.Name != "MY_EQ" {
			t.Errorf("expected MY_EQ, got %q", stmt.Name.Name)
		}
	})
	t.Run("alter_operator", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER OPERATOR my_eq ADD BINDING (VARCHAR2, VARCHAR2) RETURN NUMBER USING my_eq_str")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterOperatorStmt)
		if !ok {
			t.Fatalf("expected *AlterOperatorStmt, got %T", raw.Stmt)
		}
		if stmt.Name.Name != "MY_EQ" {
			t.Errorf("expected MY_EQ, got %q", stmt.Name.Name)
		}
	})
	t.Run("drop_operator", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP OPERATOR my_eq")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.DropStmt)
		if !ok {
			t.Fatalf("expected *DropStmt, got %T", raw.Stmt)
		}
		if stmt.ObjectType != ast.OBJECT_OPERATOR {
			t.Errorf("expected OBJECT_OPERATOR, got %d", stmt.ObjectType)
		}
	})
}

// TestBatch69_LockdownProfileOutline tests CREATE/ALTER/DROP LOCKDOWN PROFILE and OUTLINE.
func TestBatch69_LockdownProfileOutline(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_lockdown_profile", "CREATE LOCKDOWN PROFILE my_lockdown_prof", "CREATE", ast.OBJECT_LOCKDOWN_PROFILE},
		{"alter_lockdown_profile", "ALTER LOCKDOWN PROFILE my_lockdown_prof DISABLE STATEMENT = ('ALTER SYSTEM')", "ALTER", ast.OBJECT_LOCKDOWN_PROFILE},
		{"drop_lockdown_profile", "DROP LOCKDOWN PROFILE my_lockdown_prof", "DROP", ast.OBJECT_LOCKDOWN_PROFILE},
		{"create_outline", "CREATE OUTLINE my_outline FOR CATEGORY my_cat ON SELECT * FROM t", "CREATE", ast.OBJECT_OUTLINE},
		{"alter_outline", "ALTER OUTLINE my_outline RENAME TO my_outline2", "ALTER", ast.OBJECT_OUTLINE},
		{"drop_outline", "DROP OUTLINE my_outline", "DROP", ast.OBJECT_OUTLINE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch70_MaterializedZonemapInmemoryJoinGroup tests CREATE/ALTER/DROP MATERIALIZED ZONEMAP and INMEMORY JOIN GROUP.
func TestBatch70_MaterializedZonemapInmemoryJoinGroup(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_materialized_zonemap", "CREATE MATERIALIZED ZONEMAP sales_zmap ON sales (region_id, product_id)", "CREATE", ast.OBJECT_MATERIALIZED_ZONEMAP},
		{"alter_materialized_zonemap", "ALTER MATERIALIZED ZONEMAP sales_zmap ENABLE PRUNING", "ALTER", ast.OBJECT_MATERIALIZED_ZONEMAP},
		{"drop_materialized_zonemap", "DROP MATERIALIZED ZONEMAP sales_zmap", "DROP", ast.OBJECT_MATERIALIZED_ZONEMAP},
		{"create_inmemory_join_group", "CREATE INMEMORY JOIN GROUP my_jg (t1(id), t2(t1_id))", "CREATE", ast.OBJECT_INMEMORY_JOIN_GROUP},
		{"alter_inmemory_join_group", "ALTER INMEMORY JOIN GROUP my_jg ADD (t3(t1_id))", "ALTER", ast.OBJECT_INMEMORY_JOIN_GROUP},
		{"drop_inmemory_join_group", "DROP INMEMORY JOIN GROUP my_jg", "DROP", ast.OBJECT_INMEMORY_JOIN_GROUP},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch71_JsonDualityView tests CREATE/ALTER/DROP JSON RELATIONAL DUALITY VIEW.
// Updated in batch 93 to use dedicated AST types for CREATE/ALTER.
func TestBatch71_JsonDualityView(t *testing.T) {
	t.Run("create_json_duality_view", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE JSON RELATIONAL DUALITY VIEW emp_dv AS SELECT employee_id FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		_, ok := raw.Stmt.(*ast.CreateJsonDualityViewStmt)
		if !ok {
			t.Fatalf("expected *CreateJsonDualityViewStmt, got %T", raw.Stmt)
		}
	})
	t.Run("alter_json_duality_view", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER JSON RELATIONAL DUALITY VIEW emp_dv ENABLE LOGICAL REPLICATION")
		raw := result.Items[0].(*ast.RawStmt)
		_, ok := raw.Stmt.(*ast.AlterJsonDualityViewStmt)
		if !ok {
			t.Fatalf("expected *AlterJsonDualityViewStmt, got %T", raw.Stmt)
		}
	})
	t.Run("drop_json_duality_view", func(t *testing.T) {
		adminDDLTest(t, "DROP JSON RELATIONAL DUALITY VIEW emp_dv", "DROP", ast.OBJECT_JSON_DUALITY_VIEW)
	})
}

// TestBatch72_AlterMiscRound3 tests ALTER FLASHBACK ARCHIVE, ALTER RESOURCE COST, ALTER ROLLBACK SEGMENT.
func TestBatch72_AlterMiscRound3(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"alter_flashback_archive", "ALTER FLASHBACK ARCHIVE fla1 SET DEFAULT", "ALTER", ast.OBJECT_FLASHBACK_ARCHIVE},
		{"alter_flashback_archive_purge", "ALTER FLASHBACK ARCHIVE fla1 PURGE BEFORE TIMESTAMP SYSTIMESTAMP - 30", "ALTER", ast.OBJECT_FLASHBACK_ARCHIVE},
		{"alter_resource_cost", "ALTER RESOURCE COST CPU_PER_SESSION 100", "ALTER", ast.OBJECT_RESOURCE_COST},
		{"alter_rollback_segment", "ALTER ROLLBACK SEGMENT rbs1 ONLINE", "ALTER", ast.OBJECT_ROLLBACK_SEGMENT},
		{"alter_rollback_segment_offline", "ALTER ROLLBACK SEGMENT rbs1 OFFLINE", "ALTER", ast.OBJECT_ROLLBACK_SEGMENT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch73_RollbackSegmentEdition tests CREATE/DROP ROLLBACK SEGMENT and EDITION.
func TestBatch73_RollbackSegmentEdition(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_rollback_segment", "CREATE ROLLBACK SEGMENT rbs1 TABLESPACE undotbs", "CREATE", ast.OBJECT_ROLLBACK_SEGMENT},
		{"drop_rollback_segment", "DROP ROLLBACK SEGMENT rbs1", "DROP", ast.OBJECT_ROLLBACK_SEGMENT},
		{"create_edition", "CREATE EDITION e2 AS CHILD OF ora$base", "CREATE", ast.OBJECT_EDITION},
		{"drop_edition", "DROP EDITION e2", "DROP", ast.OBJECT_EDITION},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch74_TablespaceSet tests CREATE/ALTER/DROP TABLESPACE SET.
func TestBatch74_TablespaceSet(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_tablespace_set", "CREATE TABLESPACE SET ts_set1 IN SHARDSPACE shardspace1 USING TEMPLATE (DATAFILE SIZE 100M)", "CREATE", ast.OBJECT_TABLESPACE_SET},
		{"alter_tablespace_set", "ALTER TABLESPACE SET ts_set1 ADD DATAFILE SIZE 200M", "ALTER", ast.OBJECT_TABLESPACE_SET},
		{"drop_tablespace_set", "DROP TABLESPACE SET ts_set1", "DROP", ast.OBJECT_TABLESPACE_SET},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch75_MleEnvModule tests CREATE/DROP MLE ENV and MLE MODULE.
func TestBatch75_MleEnvModule(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_mle_env", "CREATE MLE ENV my_env IMPORTS ('module1')", "CREATE", ast.OBJECT_MLE_ENV},
		{"drop_mle_env", "DROP MLE ENV my_env", "DROP", ast.OBJECT_MLE_ENV},
		{"create_mle_module", "CREATE MLE MODULE my_module LANGUAGE JAVASCRIPT AS 'export function hello() { return 1; }'", "CREATE", ast.OBJECT_MLE_MODULE},
		{"drop_mle_module", "DROP MLE MODULE my_module", "DROP", ast.OBJECT_MLE_MODULE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch76_PfileSpfile tests CREATE PFILE and SPFILE.
func TestBatch76_PfileSpfile(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_pfile", "CREATE PFILE FROM SPFILE", "CREATE", ast.OBJECT_PFILE},
		{"create_pfile_named", "CREATE PFILE '/tmp/initTEST.ora' FROM SPFILE '/u01/app/oracle/dbs/spfileTEST.ora'", "CREATE", ast.OBJECT_PFILE},
		{"create_spfile", "CREATE SPFILE FROM PFILE", "CREATE", ast.OBJECT_SPFILE},
		{"create_spfile_named", "CREATE SPFILE '/u01/app/oracle/dbs/spfileTEST.ora' FROM PFILE '/tmp/initTEST.ora'", "CREATE", ast.OBJECT_SPFILE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch77_PropertyGraphVectorIndex tests CREATE PROPERTY GRAPH and VECTOR INDEX.
func TestBatch77_PropertyGraphVectorIndex(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_property_graph", "CREATE PROPERTY GRAPH my_graph VERTEX TABLES (persons) EDGE TABLES (friendships)", "CREATE", ast.OBJECT_PROPERTY_GRAPH},
		{"create_vector_index", "CREATE VECTOR INDEX vec_idx ON docs (embedding) ORGANIZATION NEIGHBOR PARTITIONS", "CREATE", ast.OBJECT_VECTOR_INDEX},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch78_RestorePointMisc tests CREATE/DROP RESTORE POINT, LOGICAL PARTITION TRACKING, PMEM FILESTORE.
func TestBatch78_RestorePointMisc(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		action string
		obj    ast.ObjectType
	}{
		{"create_restore_point", "CREATE RESTORE POINT before_upgrade", "CREATE", ast.OBJECT_RESTORE_POINT},
		{"create_restore_point_guaranteed", "CREATE RESTORE POINT before_upgrade GUARANTEE FLASHBACK DATABASE", "CREATE", ast.OBJECT_RESTORE_POINT},
		{"drop_restore_point", "DROP RESTORE POINT before_upgrade", "DROP", ast.OBJECT_RESTORE_POINT},
		{"create_logical_partition_tracking", "CREATE LOGICAL PARTITION TRACKING my_tracking ON my_table", "CREATE", ast.OBJECT_LOGICAL_PARTITION_TRACKING},
		{"create_pmem_filestore", "CREATE PMEM FILESTORE my_pmem MOUNTPOINT '/pmem0' SIZE 100G", "CREATE", ast.OBJECT_PMEM_FILESTORE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminDDLTest(t, tt.sql, tt.action, tt.obj)
		})
	}
}

// TestBatch79_TruncateClusterDropTypeBody tests TRUNCATE CLUSTER and DROP TYPE BODY.
func TestBatch79_TruncateClusterDropTypeBody(t *testing.T) {
	// TRUNCATE CLUSTER
	t.Run("truncate_cluster", func(t *testing.T) {
		result := ParseAndCheck(t, "TRUNCATE CLUSTER my_cluster")
		if result.Len() != 1 {
			t.Fatalf("expected 1 statement, got %d", result.Len())
		}
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.TruncateStmt)
		if !ok {
			t.Fatalf("expected *TruncateStmt, got %T", raw.Stmt)
		}
		if !stmt.Cluster {
			t.Errorf("expected Cluster=true")
		}
	})

	// TRUNCATE CLUSTER with options
	t.Run("truncate_cluster_storage", func(t *testing.T) {
		result := ParseAndCheck(t, "TRUNCATE CLUSTER my_cluster DROP STORAGE")
		if result.Len() != 1 {
			t.Fatalf("expected 1 statement, got %d", result.Len())
		}
		raw := result.Items[0].(*ast.RawStmt)
		_, ok := raw.Stmt.(*ast.TruncateStmt)
		if !ok {
			t.Fatalf("expected *TruncateStmt, got %T", raw.Stmt)
		}
	})

	// DROP TYPE BODY
	t.Run("drop_type_body", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP TYPE BODY my_type")
		if result.Len() != 1 {
			t.Fatalf("expected 1 statement, got %d", result.Len())
		}
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.DropStmt)
		if !ok {
			t.Fatalf("expected *DropStmt, got %T", raw.Stmt)
		}
		if stmt.ObjectType != ast.OBJECT_TYPE_BODY {
			t.Errorf("expected OBJECT_TYPE_BODY, got %d", stmt.ObjectType)
		}
	})
}

// TestParseAlterIndexProper tests ALTER INDEX with proper parsing.
func TestParseAlterIndexProper(t *testing.T) {
	// REBUILD
	t.Run("alter_index_rebuild", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX my_idx REBUILD")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterIndexStmt)
		if !ok {
			t.Fatalf("expected *AlterIndexStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "REBUILD" {
			t.Errorf("expected REBUILD, got %q", stmt.Action)
		}
	})

	// REBUILD ONLINE TABLESPACE
	t.Run("alter_index_rebuild_online_tablespace", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX hr.idx1 REBUILD ONLINE TABLESPACE users_ts")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterIndexStmt)
		if !ok {
			t.Fatalf("expected *AlterIndexStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "REBUILD" {
			t.Errorf("expected REBUILD, got %q", stmt.Action)
		}
		if !stmt.Online {
			t.Error("expected Online=true")
		}
		if stmt.Tablespace != "USERS_TS" {
			t.Errorf("expected USERS_TS, got %q", stmt.Tablespace)
		}
	})

	// REBUILD PARTITION
	t.Run("alter_index_rebuild_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 REBUILD PARTITION p1 TABLESPACE ts1 ONLINE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterIndexStmt)
		if !ok {
			t.Fatalf("expected *AlterIndexStmt, got %T", raw.Stmt)
		}
		if stmt.Partition != "P1" {
			t.Errorf("expected P1, got %q", stmt.Partition)
		}
	})

	// REBUILD REVERSE
	t.Run("alter_index_rebuild_reverse", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 REBUILD REVERSE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if !stmt.Reverse {
			t.Error("expected Reverse=true")
		}
	})

	// RENAME TO
	t.Run("alter_index_rename", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX old_idx RENAME TO new_idx")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterIndexStmt)
		if !ok {
			t.Fatalf("expected *AlterIndexStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "RENAME" {
			t.Errorf("expected RENAME, got %q", stmt.Action)
		}
		if stmt.NewName != "NEW_IDX" {
			t.Errorf("expected NEW_IDX, got %q", stmt.NewName)
		}
	})

	// COALESCE
	t.Run("alter_index_coalesce", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 COALESCE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "COALESCE" {
			t.Errorf("expected COALESCE, got %q", stmt.Action)
		}
	})

	// COALESCE CLEANUP ONLY
	t.Run("alter_index_coalesce_cleanup_only", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 COALESCE CLEANUP ONLY")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if !stmt.Cleanup {
			t.Error("expected Cleanup=true")
		}
		if !stmt.CleanupOnly {
			t.Error("expected CleanupOnly=true")
		}
	})

	// MONITORING USAGE
	t.Run("alter_index_monitoring_usage", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MONITORING USAGE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "MONITORING_USAGE" {
			t.Errorf("expected MONITORING_USAGE, got %q", stmt.Action)
		}
	})

	// NOMONITORING USAGE
	t.Run("alter_index_nomonitoring_usage", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 NOMONITORING USAGE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "NOMONITORING_USAGE" {
			t.Errorf("expected NOMONITORING_USAGE, got %q", stmt.Action)
		}
	})

	// UNUSABLE
	t.Run("alter_index_unusable", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 UNUSABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "UNUSABLE" {
			t.Errorf("expected UNUSABLE, got %q", stmt.Action)
		}
	})

	// UNUSABLE ONLINE
	t.Run("alter_index_unusable_online", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 UNUSABLE ONLINE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "UNUSABLE" {
			t.Errorf("expected UNUSABLE, got %q", stmt.Action)
		}
		if !stmt.Online {
			t.Error("expected Online=true")
		}
	})

	// VISIBLE
	t.Run("alter_index_visible", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 VISIBLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "VISIBLE" {
			t.Errorf("expected VISIBLE, got %q", stmt.Action)
		}
	})

	// INVISIBLE
	t.Run("alter_index_invisible", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 INVISIBLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "INVISIBLE" {
			t.Errorf("expected INVISIBLE, got %q", stmt.Action)
		}
	})

	// ENABLE
	t.Run("alter_index_enable", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 ENABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "ENABLE" {
			t.Errorf("expected ENABLE, got %q", stmt.Action)
		}
	})

	// DISABLE
	t.Run("alter_index_disable", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 DISABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "DISABLE" {
			t.Errorf("expected DISABLE, got %q", stmt.Action)
		}
	})

	// COMPILE
	t.Run("alter_index_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 COMPILE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "COMPILE" {
			t.Errorf("expected COMPILE, got %q", stmt.Action)
		}
	})

	// SHRINK SPACE
	t.Run("alter_index_shrink_space", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 SHRINK SPACE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "SHRINK_SPACE" {
			t.Errorf("expected SHRINK_SPACE, got %q", stmt.Action)
		}
	})

	// SHRINK SPACE COMPACT CASCADE
	t.Run("alter_index_shrink_space_compact_cascade", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 SHRINK SPACE COMPACT CASCADE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if !stmt.Compact {
			t.Error("expected Compact=true")
		}
		if !stmt.Cascade {
			t.Error("expected Cascade=true")
		}
	})

	// PARALLEL
	t.Run("alter_index_parallel", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 PARALLEL 4")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "PARALLEL" {
			t.Errorf("expected PARALLEL, got %q", stmt.Action)
		}
		if stmt.Parallel != "4" {
			t.Errorf("expected 4, got %q", stmt.Parallel)
		}
	})

	// NOPARALLEL
	t.Run("alter_index_noparallel", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 NOPARALLEL")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "NOPARALLEL" {
			t.Errorf("expected NOPARALLEL, got %q", stmt.Action)
		}
	})

	// LOGGING
	t.Run("alter_index_logging", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 LOGGING")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "LOGGING" {
			t.Errorf("expected LOGGING, got %q", stmt.Action)
		}
	})

	// NOLOGGING
	t.Run("alter_index_nologging", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 NOLOGGING")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "NOLOGGING" {
			t.Errorf("expected NOLOGGING, got %q", stmt.Action)
		}
	})

	// UPDATE BLOCK REFERENCES
	t.Run("alter_index_update_block_references", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 UPDATE BLOCK REFERENCES")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "UPDATE_BLOCK_REFERENCES" {
			t.Errorf("expected UPDATE_BLOCK_REFERENCES, got %q", stmt.Action)
		}
	})

	// IF EXISTS
	t.Run("alter_index_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX IF EXISTS idx1 REBUILD")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterIndexStmt)
		if !stmt.IfExists {
			t.Error("expected IfExists=true")
		}
	})
}

// TestParseAlterViewProper tests ALTER VIEW with proper parsing.
func TestParseAlterViewProper(t *testing.T) {
	// COMPILE
	t.Run("alter_view_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW hr.emp_view COMPILE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterViewStmt)
		if !ok {
			t.Fatalf("expected *AlterViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "COMPILE" {
			t.Errorf("expected COMPILE, got %q", stmt.Action)
		}
	})

	// ADD CONSTRAINT
	t.Run("alter_view_add_constraint", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view ADD CONSTRAINT pk_emp PRIMARY KEY (emp_id) DISABLE NOVALIDATE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterViewStmt)
		if !ok {
			t.Fatalf("expected *AlterViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "ADD_CONSTRAINT" {
			t.Errorf("expected ADD_CONSTRAINT, got %q", stmt.Action)
		}
		if stmt.Constraint == nil {
			t.Fatal("expected non-nil Constraint")
		}
	})

	// MODIFY CONSTRAINT ... RELY
	t.Run("alter_view_modify_constraint_rely", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view MODIFY CONSTRAINT pk_emp RELY")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "MODIFY_CONSTRAINT" {
			t.Errorf("expected MODIFY_CONSTRAINT, got %q", stmt.Action)
		}
		if stmt.ConstraintName != "PK_EMP" {
			t.Errorf("expected PK_EMP, got %q", stmt.ConstraintName)
		}
		if !stmt.Rely {
			t.Error("expected Rely=true")
		}
	})

	// MODIFY CONSTRAINT ... NORELY
	t.Run("alter_view_modify_constraint_norely", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view MODIFY CONSTRAINT pk_emp NORELY")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if !stmt.NoRely {
			t.Error("expected NoRely=true")
		}
	})

	// DROP CONSTRAINT
	t.Run("alter_view_drop_constraint", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view DROP CONSTRAINT pk_emp")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "DROP_CONSTRAINT" {
			t.Errorf("expected DROP_CONSTRAINT, got %q", stmt.Action)
		}
		if stmt.ConstraintName != "PK_EMP" {
			t.Errorf("expected PK_EMP, got %q", stmt.ConstraintName)
		}
	})

	// READ ONLY
	t.Run("alter_view_read_only", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view READ ONLY")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "READ_ONLY" {
			t.Errorf("expected READ_ONLY, got %q", stmt.Action)
		}
	})

	// READ WRITE
	t.Run("alter_view_read_write", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view READ WRITE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "READ_WRITE" {
			t.Errorf("expected READ_WRITE, got %q", stmt.Action)
		}
	})

	// EDITIONABLE
	t.Run("alter_view_editionable", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view EDITIONABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "EDITIONABLE" {
			t.Errorf("expected EDITIONABLE, got %q", stmt.Action)
		}
	})

	// NONEDITIONABLE
	t.Run("alter_view_noneditionable", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_view NONEDITIONABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "NONEDITIONABLE" {
			t.Errorf("expected NONEDITIONABLE, got %q", stmt.Action)
		}
	})

	// IF EXISTS
	t.Run("alter_view_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW IF EXISTS emp_view COMPILE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if !stmt.IfExists {
			t.Error("expected IfExists=true")
		}
	})
}

// TestParseAlterSequenceProper tests ALTER SEQUENCE with proper parsing.
func TestParseAlterSequenceProper(t *testing.T) {
	// Basic options
	t.Run("alter_sequence_increment_by", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 INCREMENT BY 5")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterSequenceStmt)
		if !ok {
			t.Fatalf("expected *AlterSequenceStmt, got %T", raw.Stmt)
		}
		if stmt.IncrementBy == nil {
			t.Error("expected non-nil IncrementBy")
		}
	})

	// MAXVALUE / NOMAXVALUE
	t.Run("alter_sequence_maxvalue", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 MAXVALUE 1000")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if stmt.MaxValue == nil {
			t.Error("expected non-nil MaxValue")
		}
	})

	t.Run("alter_sequence_nomaxvalue", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOMAXVALUE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoMaxValue {
			t.Error("expected NoMaxValue=true")
		}
	})

	// MINVALUE / NOMINVALUE
	t.Run("alter_sequence_minvalue", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 MINVALUE 1")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if stmt.MinValue == nil {
			t.Error("expected non-nil MinValue")
		}
	})

	t.Run("alter_sequence_nominvalue", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOMINVALUE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoMinValue {
			t.Error("expected NoMinValue=true")
		}
	})

	// CYCLE / NOCYCLE
	t.Run("alter_sequence_cycle", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 CYCLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Cycle {
			t.Error("expected Cycle=true")
		}
	})

	t.Run("alter_sequence_nocycle", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOCYCLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoCycle {
			t.Error("expected NoCycle=true")
		}
	})

	// CACHE / NOCACHE
	t.Run("alter_sequence_cache", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 CACHE 20")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if stmt.Cache == nil {
			t.Error("expected non-nil Cache")
		}
	})

	t.Run("alter_sequence_nocache", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOCACHE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoCache {
			t.Error("expected NoCache=true")
		}
	})

	// ORDER / NOORDER
	t.Run("alter_sequence_order", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 ORDER")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Order {
			t.Error("expected Order=true")
		}
	})

	t.Run("alter_sequence_noorder", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOORDER")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoOrder {
			t.Error("expected NoOrder=true")
		}
	})

	// KEEP / NOKEEP
	t.Run("alter_sequence_keep", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 KEEP")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Keep {
			t.Error("expected Keep=true")
		}
	})

	t.Run("alter_sequence_nokeep", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOKEEP")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoKeep {
			t.Error("expected NoKeep=true")
		}
	})

	// SCALE / NOSCALE
	t.Run("alter_sequence_scale_extend", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 SCALE EXTEND")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Scale {
			t.Error("expected Scale=true")
		}
		if !stmt.ScaleExtend {
			t.Error("expected ScaleExtend=true")
		}
	})

	t.Run("alter_sequence_noscale", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOSCALE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoScale {
			t.Error("expected NoScale=true")
		}
	})

	// SHARD / NOSHARD
	t.Run("alter_sequence_shard_extend", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 SHARD EXTEND")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Shard {
			t.Error("expected Shard=true")
		}
		if !stmt.ShardExtend {
			t.Error("expected ShardExtend=true")
		}
	})

	t.Run("alter_sequence_noshard", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 NOSHARD")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.NoShard {
			t.Error("expected NoShard=true")
		}
	})

	// RESTART
	t.Run("alter_sequence_restart", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 RESTART")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Restart {
			t.Error("expected Restart=true")
		}
	})

	t.Run("alter_sequence_restart_with", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 RESTART WITH 100")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Restart {
			t.Error("expected Restart=true")
		}
		if stmt.RestartWith == nil {
			t.Error("expected non-nil RestartWith")
		}
	})

	// GLOBAL / SESSION
	t.Run("alter_sequence_global", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 GLOBAL")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Global {
			t.Error("expected Global=true")
		}
	})

	t.Run("alter_sequence_session", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 SESSION")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.Session {
			t.Error("expected Session=true")
		}
	})

	// Multiple options
	t.Run("alter_sequence_multiple_options", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE seq1 INCREMENT BY 10 MAXVALUE 9999 CACHE 50 CYCLE ORDER")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if stmt.IncrementBy == nil {
			t.Error("expected non-nil IncrementBy")
		}
		if stmt.MaxValue == nil {
			t.Error("expected non-nil MaxValue")
		}
		if stmt.Cache == nil {
			t.Error("expected non-nil Cache")
		}
		if !stmt.Cycle {
			t.Error("expected Cycle=true")
		}
		if !stmt.Order {
			t.Error("expected Order=true")
		}
	})

	// IF EXISTS
	t.Run("alter_sequence_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SEQUENCE IF EXISTS seq1 INCREMENT BY 1")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSequenceStmt)
		if !stmt.IfExists {
			t.Error("expected IfExists=true")
		}
	})
}

func TestParseAlterType(t *testing.T) {
	tests := []string{
		// COMPILE variations
		"ALTER TYPE my_type COMPILE",
		"ALTER TYPE my_schema.my_type COMPILE",
		"ALTER TYPE my_type COMPILE SPECIFICATION",
		"ALTER TYPE my_type COMPILE BODY",
		"ALTER TYPE my_type COMPILE DEBUG",
		"ALTER TYPE my_type COMPILE SPECIFICATION REUSE SETTINGS",
		"ALTER TYPE my_type COMPILE BODY DEBUG PLSQL_OPTIMIZE_LEVEL = 2 REUSE SETTINGS",

		// ADD ATTRIBUTE
		"ALTER TYPE person_t ADD ATTRIBUTE (email VARCHAR2(100)) CASCADE",
		"ALTER TYPE person_t ADD ATTRIBUTE (phone NUMBER) INVALIDATE",
		"ALTER TYPE person_t ADD ATTRIBUTE (addr VARCHAR2(200), zip NUMBER) CASCADE NOT INCLUDING TABLE DATA",

		// DROP ATTRIBUTE
		"ALTER TYPE person_t DROP ATTRIBUTE (email) CASCADE",
		"ALTER TYPE person_t DROP ATTRIBUTE (email, phone) INVALIDATE",

		// MODIFY ATTRIBUTE
		"ALTER TYPE person_t MODIFY ATTRIBUTE (email VARCHAR2(200)) CASCADE",
		"ALTER TYPE person_t MODIFY ATTRIBUTE (email VARCHAR2(200), phone NUMBER(20)) CASCADE INCLUDING TABLE DATA",

		// ADD method
		"ALTER TYPE data_typ1 ADD MEMBER FUNCTION get_name RETURN VARCHAR2 CASCADE",
		"ALTER TYPE data_typ1 ADD STATIC PROCEDURE init(p1 NUMBER) CASCADE",
		"ALTER TYPE data_typ1 ADD MAP MEMBER FUNCTION cmp RETURN NUMBER CASCADE",
		"ALTER TYPE data_typ1 ADD ORDER MEMBER FUNCTION cmp(other data_typ1) RETURN NUMBER CASCADE",
		"ALTER TYPE data_typ1 ADD CONSTRUCTOR FUNCTION data_typ1(x NUMBER) RETURN SELF AS RESULT CASCADE",

		// DROP method
		"ALTER TYPE data_typ1 DROP MEMBER FUNCTION get_name RETURN VARCHAR2 CASCADE",
		"ALTER TYPE data_typ1 DROP STATIC PROCEDURE init(p1 NUMBER) INVALIDATE",
		"ALTER TYPE data_typ1 DROP MAP MEMBER FUNCTION cmp RETURN NUMBER CASCADE",

		// NOT INSTANTIABLE / NOT FINAL / FINAL / INSTANTIABLE
		"ALTER TYPE my_type NOT INSTANTIABLE CASCADE",
		"ALTER TYPE my_type NOT FINAL CASCADE",
		"ALTER TYPE my_type FINAL CASCADE",
		"ALTER TYPE my_type INSTANTIABLE CASCADE",

		// MODIFY LIMIT (varray)
		"ALTER TYPE phone_list_t MODIFY LIMIT 20 CASCADE",

		// MODIFY ELEMENT TYPE (collection)
		"ALTER TYPE phone_list_t MODIFY ELEMENT TYPE VARCHAR2(64) CASCADE",

		// EDITIONABLE / NONEDITIONABLE
		"ALTER TYPE my_type EDITIONABLE",
		"ALTER TYPE my_type NONEDITIONABLE",

		// IF EXISTS
		"ALTER TYPE IF EXISTS my_type COMPILE",

		// RESET
		"ALTER TYPE my_type RESET",

		// CASCADE with INCLUDING/NOT INCLUDING TABLE DATA
		"ALTER TYPE person_t ADD ATTRIBUTE (x NUMBER) CASCADE INCLUDING TABLE DATA",
		"ALTER TYPE person_t ADD ATTRIBUTE (x NUMBER) CASCADE NOT INCLUDING TABLE DATA",
		"ALTER TYPE person_t ADD ATTRIBUTE (x NUMBER) CASCADE CONVERT TO SUBSTITUTABLE",

		// FORCE (exceptions clause)
		"ALTER TYPE person_t ADD ATTRIBUTE (x NUMBER) CASCADE FORCE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify AST node type
	t.Run("alter_type_ast_check", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER TYPE my_schema.my_type ADD ATTRIBUTE (email VARCHAR2(100)) CASCADE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterTypeStmt)
		if stmt.Name == nil {
			t.Fatal("expected non-nil Name")
		}
		if stmt.Action != "ADD_ATTRIBUTE" {
			t.Errorf("expected Action=ADD_ATTRIBUTE, got %q", stmt.Action)
		}
	})

	t.Run("alter_type_compile_check", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER TYPE my_type COMPILE SPECIFICATION DEBUG REUSE SETTINGS")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterTypeStmt)
		if stmt.Action != "COMPILE" {
			t.Errorf("expected Action=COMPILE, got %q", stmt.Action)
		}
		if stmt.CompileTarget != "SPECIFICATION" {
			t.Errorf("expected CompileTarget=SPECIFICATION, got %q", stmt.CompileTarget)
		}
		if !stmt.Debug {
			t.Error("expected Debug=true")
		}
		if !stmt.ReuseSettings {
			t.Error("expected ReuseSettings=true")
		}
	})
}

func TestParseAlterMaterializedView(t *testing.T) {
	tests := []string{
		// COMPILE
		"ALTER MATERIALIZED VIEW mv1 COMPILE",
		// CONSIDER FRESH
		"ALTER MATERIALIZED VIEW hr.mv_emp CONSIDER FRESH",
		// REFRESH COMPLETE
		"ALTER MATERIALIZED VIEW mv1 REFRESH COMPLETE",
		// REFRESH FAST
		"ALTER MATERIALIZED VIEW mv1 REFRESH FAST",
		// REFRESH FORCE
		"ALTER MATERIALIZED VIEW mv1 REFRESH FORCE",
		// REFRESH ON DEMAND
		"ALTER MATERIALIZED VIEW mv1 REFRESH ON DEMAND",
		// REFRESH ON COMMIT
		"ALTER MATERIALIZED VIEW mv1 REFRESH ON COMMIT",
		// REFRESH FAST ON COMMIT
		"ALTER MATERIALIZED VIEW mv1 REFRESH FAST ON COMMIT",
		// REFRESH COMPLETE ON DEMAND
		"ALTER MATERIALIZED VIEW mv1 REFRESH COMPLETE ON DEMAND",
		// REFRESH with START WITH / NEXT
		"ALTER MATERIALIZED VIEW mv1 REFRESH START WITH SYSDATE NEXT SYSDATE + 1",
		// REFRESH WITH PRIMARY KEY
		"ALTER MATERIALIZED VIEW mv1 REFRESH WITH PRIMARY KEY",
		// REFRESH USING ROLLBACK SEGMENT
		"ALTER MATERIALIZED VIEW mv1 REFRESH USING ROLLBACK SEGMENT rbs1",
		// REFRESH USING ENFORCED CONSTRAINTS
		"ALTER MATERIALIZED VIEW mv1 REFRESH USING ENFORCED CONSTRAINTS",
		// REFRESH USING TRUSTED CONSTRAINTS
		"ALTER MATERIALIZED VIEW mv1 REFRESH USING TRUSTED CONSTRAINTS",
		// REFRESH combined
		"ALTER MATERIALIZED VIEW mv1 REFRESH FORCE ON DEMAND START WITH SYSDATE NEXT SYSDATE + 7",
		// ENABLE QUERY REWRITE
		"ALTER MATERIALIZED VIEW mv1 ENABLE QUERY REWRITE",
		// DISABLE QUERY REWRITE
		"ALTER MATERIALIZED VIEW mv1 DISABLE QUERY REWRITE",
		// ENABLE ON QUERY COMPUTATION
		"ALTER MATERIALIZED VIEW mv1 REFRESH ENABLE ON QUERY COMPUTATION",
		// DISABLE ON QUERY COMPUTATION
		"ALTER MATERIALIZED VIEW mv1 REFRESH DISABLE ON QUERY COMPUTATION",
		// SHRINK SPACE
		"ALTER MATERIALIZED VIEW mv1 SHRINK SPACE",
		// SHRINK SPACE COMPACT
		"ALTER MATERIALIZED VIEW mv1 SHRINK SPACE COMPACT",
		// SHRINK SPACE CASCADE
		"ALTER MATERIALIZED VIEW mv1 SHRINK SPACE CASCADE",
		// CACHE / NOCACHE
		"ALTER MATERIALIZED VIEW mv1 CACHE",
		"ALTER MATERIALIZED VIEW mv1 NOCACHE",
		// PARALLEL / NOPARALLEL
		"ALTER MATERIALIZED VIEW mv1 PARALLEL",
		"ALTER MATERIALIZED VIEW mv1 PARALLEL 4",
		"ALTER MATERIALIZED VIEW mv1 NOPARALLEL",
		// LOGGING / NOLOGGING
		"ALTER MATERIALIZED VIEW mv1 LOGGING",
		"ALTER MATERIALIZED VIEW mv1 NOLOGGING",
		// ENABLE CONCURRENT REFRESH
		"ALTER MATERIALIZED VIEW mv1 ENABLE CONCURRENT REFRESH",
		// DISABLE CONCURRENT REFRESH
		"ALTER MATERIALIZED VIEW mv1 DISABLE CONCURRENT REFRESH",
		// IF EXISTS
		"ALTER MATERIALIZED VIEW IF EXISTS mv1 COMPILE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseAlterDatabaseLink tests parsing of ALTER DATABASE LINK statements.
func TestParseAlterDatabaseLink(t *testing.T) {
	tests := []string{
		// Basic ALTER DATABASE LINK
		"ALTER DATABASE LINK private_link CONNECT TO hr IDENTIFIED BY hr_password",
		// PUBLIC
		"ALTER PUBLIC DATABASE LINK pub_link CONNECT TO scott IDENTIFIED BY tiger",
		// SHARED
		"ALTER SHARED DATABASE LINK shared_link CONNECT TO hr IDENTIFIED BY hr_pass",
		// SHARED PUBLIC with AUTHENTICATED BY
		"ALTER SHARED PUBLIC DATABASE LINK shared_pub_link CONNECT TO scott IDENTIFIED BY tiger AUTHENTICATED BY hr IDENTIFIED BY hr_pass",
		// Schema-qualified link name
		"ALTER DATABASE LINK my_schema.my_link CONNECT TO admin IDENTIFIED BY secret",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify AST node type
	t.Run("alter_dblink_ast_check", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER SHARED PUBLIC DATABASE LINK mylink CONNECT TO scott IDENTIFIED BY tiger AUTHENTICATED BY hr IDENTIFIED BY hr_pass")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterDatabaseLinkStmt)
		if stmt.Name == nil {
			t.Fatal("expected non-nil Name")
		}
		if !stmt.Shared {
			t.Fatal("expected Shared=true")
		}
		if !stmt.Public {
			t.Fatal("expected Public=true")
		}
		if stmt.ConnectUser != "SCOTT" {
			t.Fatalf("expected ConnectUser=SCOTT, got %s", stmt.ConnectUser)
		}
		if stmt.ConnectPassword != "TIGER" {
			t.Fatalf("expected ConnectPassword=TIGER, got %s", stmt.ConnectPassword)
		}
		if stmt.AuthenticatedUser != "HR" {
			t.Fatalf("expected AuthenticatedUser=HR, got %s", stmt.AuthenticatedUser)
		}
		if stmt.AuthenticatedPass != "HR_PASS" {
			t.Fatalf("expected AuthenticatedPass=HR_PASS, got %s", stmt.AuthenticatedPass)
		}
	})
}

// TestParseAlterSynonym tests parsing of ALTER SYNONYM statements.
func TestParseAlterSynonym(t *testing.T) {
	tests := []string{
		// COMPILE
		"ALTER SYNONYM my_syn COMPILE",
		// PUBLIC COMPILE
		"ALTER PUBLIC SYNONYM emp_table COMPILE",
		// EDITIONABLE
		"ALTER SYNONYM my_syn EDITIONABLE",
		// NONEDITIONABLE
		"ALTER SYNONYM my_syn NONEDITIONABLE",
		// IF EXISTS
		"ALTER SYNONYM IF EXISTS my_syn COMPILE",
		// PUBLIC IF EXISTS
		"ALTER PUBLIC SYNONYM IF EXISTS my_syn COMPILE",
		// Schema-qualified
		"ALTER SYNONYM my_schema.my_syn COMPILE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify AST node type
	t.Run("alter_synonym_ast_check", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER PUBLIC SYNONYM IF EXISTS my_syn EDITIONABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterSynonymStmt)
		if stmt.Name == nil {
			t.Fatal("expected non-nil Name")
		}
		if !stmt.Public {
			t.Fatal("expected Public=true")
		}
		if !stmt.IfExists {
			t.Fatal("expected IfExists=true")
		}
		if stmt.Action != "EDITIONABLE" {
			t.Fatalf("expected Action=EDITIONABLE, got %s", stmt.Action)
		}
	})
}

// TestParseAlterDatabaseProper tests ALTER DATABASE clause parsing (batch 87).
func TestParseAlterDatabaseProper(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// startup_clauses
		{"mount", "ALTER DATABASE MOUNT"},
		{"mount_standby", "ALTER DATABASE MOUNT STANDBY DATABASE"},
		{"mount_clone", "ALTER DATABASE MOUNT CLONE DATABASE"},
		{"open_default", "ALTER DATABASE OPEN"},
		{"open_read_write", "ALTER DATABASE OPEN READ WRITE"},
		{"open_read_only", "ALTER DATABASE OPEN READ ONLY"},
		{"open_resetlogs", "ALTER DATABASE OPEN RESETLOGS"},
		{"open_noresetlogs", "ALTER DATABASE OPEN NORESETLOGS"},
		{"open_upgrade", "ALTER DATABASE OPEN READ WRITE RESETLOGS UPGRADE"},
		{"open_downgrade", "ALTER DATABASE OPEN DOWNGRADE"},

		// recovery_clauses
		{"recover_database", "ALTER DATABASE RECOVER DATABASE"},
		{"recover_automatic", "ALTER DATABASE RECOVER AUTOMATIC DATABASE"},
		{"recover_managed_standby", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE"},
		{"recover_managed_cancel", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE CANCEL"},
		{"recover_managed_disconnect", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE DISCONNECT FROM SESSION"},
		{"begin_backup", "ALTER DATABASE BEGIN BACKUP"},
		{"end_backup", "ALTER DATABASE END BACKUP"},
		{"recover_datafile", "ALTER DATABASE RECOVER DATAFILE '/u01/data01.dbf'"},
		{"recover_tablespace", "ALTER DATABASE RECOVER TABLESPACE users"},

		// database_file_clauses
		{"rename_file", "ALTER DATABASE RENAME FILE '/old/file.dbf' TO '/new/file.dbf'"},
		{"rename_file_multi", "ALTER DATABASE RENAME FILE '/old1.dbf', '/old2.dbf' TO '/new1.dbf', '/new2.dbf'"},
		{"create_datafile", "ALTER DATABASE CREATE DATAFILE '/u01/data01.dbf'"},
		{"create_datafile_as", "ALTER DATABASE CREATE DATAFILE '/u01/data01.dbf' AS '/u01/data02.dbf'"},
		{"datafile_online", "ALTER DATABASE DATAFILE '/u01/data01.dbf' ONLINE"},
		{"datafile_offline", "ALTER DATABASE DATAFILE '/u01/data01.dbf' OFFLINE"},
		{"datafile_resize", "ALTER DATABASE DATAFILE '/u01/data01.dbf' RESIZE 100M"},
		{"datafile_autoextend", "ALTER DATABASE DATAFILE '/u01/data01.dbf' AUTOEXTEND ON NEXT 10M MAXSIZE 1G"},
		{"datafile_end_backup", "ALTER DATABASE DATAFILE '/u01/data01.dbf' END BACKUP"},
		{"tempfile_resize", "ALTER DATABASE TEMPFILE '/u01/temp01.dbf' RESIZE 200M"},
		{"tempfile_autoextend", "ALTER DATABASE TEMPFILE '/u01/temp01.dbf' AUTOEXTEND ON NEXT 50M"},
		{"tempfile_online", "ALTER DATABASE TEMPFILE '/u01/temp01.dbf' ONLINE"},
		{"tempfile_offline", "ALTER DATABASE TEMPFILE '/u01/temp01.dbf' OFFLINE"},
		{"tempfile_drop", "ALTER DATABASE TEMPFILE '/u01/temp01.dbf' DROP INCLUDING DATAFILES"},
		{"move_datafile", "ALTER DATABASE MOVE DATAFILE '/u01/old.dbf' TO '/u01/new.dbf'"},

		// logfile_clauses
		{"add_logfile", "ALTER DATABASE ADD LOGFILE GROUP 3 '/u01/redo03.log' SIZE 100M"},
		{"add_logfile_member", "ALTER DATABASE ADD LOGFILE MEMBER '/u01/redo03b.log' TO GROUP 3"},
		{"drop_logfile", "ALTER DATABASE DROP LOGFILE GROUP 3"},
		{"drop_logfile_member", "ALTER DATABASE DROP LOGFILE MEMBER '/u01/redo03b.log'"},
		{"add_standby_logfile", "ALTER DATABASE ADD STANDBY LOGFILE GROUP 10 '/u01/standby10.log' SIZE 100M"},
		{"drop_standby_logfile", "ALTER DATABASE DROP STANDBY LOGFILE GROUP 10"},
		{"clear_logfile", "ALTER DATABASE CLEAR LOGFILE GROUP 3"},
		{"clear_unarchived_logfile", "ALTER DATABASE CLEAR UNARCHIVED LOGFILE GROUP 3"},
		{"switch_logfile", "ALTER DATABASE SWITCH ALL LOGFILE"},

		// controlfile_clauses
		{"backup_controlfile_to", "ALTER DATABASE BACKUP CONTROLFILE TO '/u01/backup.ctl'"},
		{"backup_controlfile_reuse", "ALTER DATABASE BACKUP CONTROLFILE TO '/u01/backup.ctl' REUSE"},
		{"backup_controlfile_trace", "ALTER DATABASE BACKUP CONTROLFILE TO TRACE"},
		{"backup_controlfile_trace_as", "ALTER DATABASE BACKUP CONTROLFILE TO TRACE AS '/u01/trace.sql'"},

		// standby_database_clauses
		{"activate_standby", "ALTER DATABASE ACTIVATE STANDBY DATABASE"},
		{"activate_physical_standby", "ALTER DATABASE ACTIVATE PHYSICAL STANDBY DATABASE"},
		{"set_standby_maximize", "ALTER DATABASE SET STANDBY DATABASE TO MAXIMIZE PROTECTION"},
		{"register_logfile", "ALTER DATABASE REGISTER LOGFILE '/u01/arch01.log'"},
		{"convert_to_standby", "ALTER DATABASE CONVERT TO PHYSICAL STANDBY"},

		// default_settings_clauses
		{"set_default_bigfile", "ALTER DATABASE SET DEFAULT BIGFILE TABLESPACE"},
		{"default_tablespace", "ALTER DATABASE DEFAULT TABLESPACE users"},
		{"default_temp_tablespace", "ALTER DATABASE DEFAULT TEMPORARY TABLESPACE temp"},
		{"rename_global_name", "ALTER DATABASE RENAME GLOBAL_NAME TO mydb.world"},
		{"enable_block_tracking", "ALTER DATABASE ENABLE BLOCK CHANGE TRACKING"},
		{"disable_block_tracking", "ALTER DATABASE DISABLE BLOCK CHANGE TRACKING"},
		{"flashback_on", "ALTER DATABASE FLASHBACK ON"},
		{"flashback_off", "ALTER DATABASE FLASHBACK OFF"},
		{"set_time_zone", "ALTER DATABASE SET TIME_ZONE = '+08:00'"},
		{"default_edition", "ALTER DATABASE DEFAULT EDITION myedition"},

		// security_clause
		{"guard_all", "ALTER DATABASE GUARD ALL"},
		{"guard_standby", "ALTER DATABASE GUARD STANDBY"},
		{"guard_none", "ALTER DATABASE GUARD NONE"},

		// instance_clauses
		{"enable_instance", "ALTER DATABASE ENABLE INSTANCE 'inst1'"},
		{"disable_instance", "ALTER DATABASE DISABLE INSTANCE 'inst1'"},

		// with database name
		{"named_db_mount", "ALTER DATABASE mydb MOUNT"},
		{"named_db_open", "ALTER DATABASE mydb OPEN READ WRITE"},

		// supplemental logging
		{"add_supplemental_log", "ALTER DATABASE ADD SUPPLEMENTAL LOG DATA"},
		{"drop_supplemental_log", "ALTER DATABASE DROP SUPPLEMENTAL LOG DATA"},

		// lost write protection / suspend resume
		{"suspend", "ALTER DATABASE SUSPEND"},
		{"resume", "ALTER DATABASE RESUME"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw, ok := result.Items[0].(*ast.RawStmt)
			if !ok {
				t.Fatalf("expected *RawStmt, got %T", result.Items[0])
			}
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Fatalf("expected action ALTER, got %s", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DATABASE {
				t.Fatalf("expected OBJECT_DATABASE, got %d", stmt.ObjectType)
			}
			// Verify that options were actually parsed (not empty skip)
			if stmt.Options == nil || len(stmt.Options.Items) == 0 {
				t.Fatal("expected parsed options, got nil/empty")
			}
		})
	}
}

// TestParseCreateControlfileProper tests CREATE CONTROLFILE parsing (batch 87).
func TestParseCreateControlfileProper(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"basic", "CREATE CONTROLFILE SET DATABASE mydb LOGFILE GROUP 1 '/u01/redo01.log' SIZE 100M RESETLOGS"},
		{"reuse", "CREATE CONTROLFILE REUSE DATABASE mydb LOGFILE GROUP 1 '/u01/redo01.log' SIZE 50M NORESETLOGS"},
		{"with_datafile", "CREATE CONTROLFILE SET DATABASE mydb LOGFILE GROUP 1 '/u01/redo01.log' SIZE 100M RESETLOGS DATAFILE '/u01/data01.dbf'"},
		{"full", "CREATE CONTROLFILE REUSE SET DATABASE mydb LOGFILE GROUP 1 '/u01/redo01.log' SIZE 100M RESETLOGS DATAFILE '/u01/data01.dbf' MAXLOGFILES 16 MAXLOGMEMBERS 3 MAXDATAFILES 100 MAXINSTANCES 8 ARCHIVELOG"},
		{"noresetlogs_noarchive", "CREATE CONTROLFILE DATABASE mydb LOGFILE GROUP 1 '/u01/redo01.log' SIZE 50M NORESETLOGS NOARCHIVELOG"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw, ok := result.Items[0].(*ast.RawStmt)
			if !ok {
				t.Fatalf("expected *RawStmt, got %T", result.Items[0])
			}
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "CREATE" {
				t.Fatalf("expected action CREATE, got %s", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_CONTROLFILE {
				t.Fatalf("expected OBJECT_CONTROLFILE, got %d", stmt.ObjectType)
			}
			if stmt.Options == nil || len(stmt.Options.Items) == 0 {
				t.Fatal("expected parsed options, got nil/empty")
			}
		})
	}
}

// TestParseAlterDatabaseDictionaryProper tests ALTER DATABASE DICTIONARY parsing (batch 87).
func TestParseAlterDatabaseDictionaryProper(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		key  string
	}{
		{"encrypt", "ALTER DATABASE DICTIONARY ENCRYPT CREDENTIALS", "ENCRYPT_CREDENTIALS"},
		{"rekey", "ALTER DATABASE DICTIONARY REKEY CREDENTIALS", "REKEY_CREDENTIALS"},
		{"delete", "ALTER DATABASE DICTIONARY DELETE CREDENTIALS KEY", "DELETE_CREDENTIALS_KEY"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
			raw, ok := result.Items[0].(*ast.RawStmt)
			if !ok {
				t.Fatalf("expected *RawStmt, got %T", result.Items[0])
			}
			stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
			if !ok {
				t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
			}
			if stmt.Action != "ALTER" {
				t.Fatalf("expected action ALTER, got %s", stmt.Action)
			}
			if stmt.ObjectType != ast.OBJECT_DATABASE_DICTIONARY {
				t.Fatalf("expected OBJECT_DATABASE_DICTIONARY, got %d", stmt.ObjectType)
			}
			if stmt.Options == nil || len(stmt.Options.Items) == 0 {
				t.Fatal("expected parsed options, got nil/empty")
			}
			opt := stmt.Options.Items[0].(*ast.DDLOption)
			if opt.Key != tc.key {
				t.Fatalf("expected key %q, got %q", tc.key, opt.Key)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 94: database_controlfile — review+fix
// ---------------------------------------------------------------------------

// TestParseDatabaseControlfileBatch94 tests all BNF branches for batch 94.
func TestParseDatabaseControlfileBatch94(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// --- FLASHBACK DATABASE: STANDBY/PLUGGABLE, BEFORE variants, database name ---
		{"flashback_standby", "FLASHBACK STANDBY DATABASE TO SCN 12345"},
		{"flashback_pluggable", "FLASHBACK PLUGGABLE DATABASE pdb1 TO TIMESTAMP SYSTIMESTAMP - 1"},
		{"flashback_before_scn", "FLASHBACK DATABASE TO BEFORE SCN 99999"},
		{"flashback_before_timestamp", "FLASHBACK DATABASE TO BEFORE TIMESTAMP SYSTIMESTAMP - 1"},
		{"flashback_before_resetlogs", "FLASHBACK DATABASE TO BEFORE RESETLOGS"},
		{"flashback_named_db", "FLASHBACK DATABASE mydb TO RESTORE POINT rp1"},
		{"flashback_standby_named", "FLASHBACK STANDBY DATABASE mydb TO SCN 555"},

		// --- DROP DATABASE ---
		{"drop_database", "DROP DATABASE"},

		// --- CREATE DATABASE: missing BNF branches ---
		{"create_db_sysaux_datafile", "CREATE DATABASE mydb SYSAUX DATAFILE '/u01/sysaux01.dbf' SIZE 500M"},
		{"create_db_default_local_temp", "CREATE DATABASE mydb DEFAULT LOCAL TEMPORARY TABLESPACE localtemp TEMPFILE '/u01/localtemp.dbf' SIZE 100M"},
		{"create_db_default_local_temp_for_all", "CREATE DATABASE mydb DEFAULT LOCAL TEMPORARY TABLESPACE localtemp FOR ALL TEMPFILE '/u01/localtemp.dbf' SIZE 100M"},
		{"create_db_extent_management", "CREATE DATABASE mydb EXTENT MANAGEMENT LOCAL"},
		{"create_db_extent_management_uniform", "CREATE DATABASE mydb EXTENT MANAGEMENT LOCAL UNIFORM SIZE 1M"},
		{"create_db_using_mirror_copy", "CREATE DATABASE mydb USING MIRROR COPY mirror1"},
		{"create_db_enable_pdb_seed", "CREATE DATABASE mydb ENABLE PLUGGABLE DATABASE SEED FILE_NAME_CONVERT = ('/u01/pdb/', '/u01/seed/')"},
		{"create_db_enable_pdb_seed_system", "CREATE DATABASE mydb ENABLE PLUGGABLE DATABASE SEED SYSTEM DATAFILES SIZE 200M AUTOEXTEND ON NEXT 10M MAXSIZE UNLIMITED"},
		{"create_db_enable_pdb_seed_sysaux", "CREATE DATABASE mydb ENABLE PLUGGABLE DATABASE SEED SYSAUX DATAFILES SIZE 200M"},
		{"create_db_enable_pdb_seed_undo", "CREATE DATABASE mydb ENABLE PLUGGABLE DATABASE SEED LOCAL UNDO ON"},

		// --- ALTER DATABASE: missing BNF branches ---
		// controlfile_clauses: CREATE STANDBY CONTROLFILE
		{"alter_db_create_standby_ctl", "ALTER DATABASE CREATE PHYSICAL STANDBY CONTROLFILE AS '/u01/standby.ctl'"},
		{"alter_db_create_logical_standby_ctl", "ALTER DATABASE CREATE LOGICAL STANDBY CONTROLFILE AS '/u01/standby.ctl' REUSE"},
		{"alter_db_create_far_sync_ctl", "ALTER DATABASE CREATE FAR SYNC INSTANCE CONTROLFILE AS '/u01/farsync.ctl' REUSE"},

		// commit_switchover_clause
		{"alter_db_prepare_switchover_logical", "ALTER DATABASE PREPARE TO SWITCHOVER TO LOGICAL STANDBY"},
		{"alter_db_prepare_switchover_primary", "ALTER DATABASE PREPARE TO SWITCHOVER TO PRIMARY DATABASE"},
		{"alter_db_commit_switchover_standby", "ALTER DATABASE COMMIT TO SWITCHOVER TO STANDBY"},
		{"alter_db_commit_switchover_logical", "ALTER DATABASE COMMIT TO SWITCHOVER TO LOGICAL STANDBY"},
		{"alter_db_commit_switchover_primary", "ALTER DATABASE COMMIT TO SWITCHOVER TO PRIMARY DATABASE"},

		// start/stop standby apply
		{"alter_db_start_standby_apply", "ALTER DATABASE START LOGICAL STANDBY APPLY"},
		{"alter_db_start_standby_apply_immediate", "ALTER DATABASE START LOGICAL STANDBY APPLY IMMEDIATE"},
		{"alter_db_stop_standby_apply", "ALTER DATABASE STOP LOGICAL STANDBY APPLY"},

		// switchover/failover clause
		{"alter_db_switchover_to", "ALTER DATABASE SWITCHOVER TO PRIMARY VERIFY FORCE mydb"},
		{"alter_db_failover_to", "ALTER DATABASE FAILOVER TO PHYSICAL STANDBY FORCE mydb"},

		// instance_clauses: ENABLE/DISABLE RESTRICTED SESSION
		{"alter_db_enable_restricted", "ALTER DATABASE ENABLE RESTRICTED SESSION"},
		{"alter_db_disable_restricted", "ALTER DATABASE DISABLE RESTRICTED SESSION"},

		// cdb_fleet_clauses
		{"alter_db_set_lead_cdb", "ALTER DATABASE SET LEAD_CDB = mydb"},
		{"alter_db_set_lead_cdb_none", "ALTER DATABASE SET LEAD_CDB = NONE"},
		{"alter_db_set_lead_cdb_uri", "ALTER DATABASE SET LEAD_CDB_URI = 'jdbc:oracle:thin:@//host:1521/mydb'"},
		{"alter_db_set_property", "ALTER DATABASE SET PROPERTY myprop = myval"},

		// replay_upgrade_clause
		{"alter_db_start_replay", "ALTER DATABASE START REPLAY"},
		{"alter_db_stop_replay", "ALTER DATABASE STOP REPLAY"},

		// ENCRYPT/DECRYPT for datafile
		{"alter_db_datafile_encrypt", "ALTER DATABASE DATAFILE '/u01/data01.dbf' ENCRYPT"},
		{"alter_db_datafile_decrypt", "ALTER DATABASE DATAFILE '/u01/data01.dbf' DECRYPT"},

		// NO FORCE LOGGING (already has basic test, ensure it works)
		{"alter_db_no_force_logging", "ALTER DATABASE NO FORCE LOGGING"},

		// ARCHIVELOG / MANUAL / NOARCHIVELOG in logfile_clauses
		{"alter_db_archivelog", "ALTER DATABASE ARCHIVELOG"},
		{"alter_db_noarchivelog", "ALTER DATABASE NOARCHIVELOG"},

		// SET STANDBY NOLOGGING FOR ...
		{"alter_db_set_standby_nologging_data", "ALTER DATABASE SET STANDBY NOLOGGING FOR DATA AVAILABILITY"},
		{"alter_db_set_standby_nologging_load", "ALTER DATABASE SET STANDBY NOLOGGING FOR LOAD PERFORMANCE"},

		// supplemental_db_logging sub-clauses
		{"alter_db_add_supp_log_pk", "ALTER DATABASE ADD SUPPLEMENTAL LOG DATA (PRIMARY KEY COLUMNS)"},
		{"alter_db_add_supp_log_unique", "ALTER DATABASE ADD SUPPLEMENTAL LOG DATA (UNIQUE COLUMNS)"},
		{"alter_db_add_supp_log_fk", "ALTER DATABASE ADD SUPPLEMENTAL LOG DATA (FOREIGN KEY COLUMNS)"},
		{"alter_db_add_supp_log_all", "ALTER DATABASE ADD SUPPLEMENTAL LOG DATA (ALL COLUMNS)"},
		{"alter_db_drop_supp_log_pk", "ALTER DATABASE DROP SUPPLEMENTAL LOG DATA (PRIMARY KEY COLUMNS)"},
		{"alter_db_add_supp_log_subset", "ALTER DATABASE ADD SUPPLEMENTAL LOG DATA SUBSET DATABASE REPLICATION"},
		{"alter_db_drop_supp_log_subset", "ALTER DATABASE DROP SUPPLEMENTAL LOG DATA SUBSET DATABASE REPLICATION"},

		// SWITCH LOGFILE TO BLOCK SIZE
		{"alter_db_switch_logfile_blocksize", "ALTER DATABASE SWITCH LOGFILE TO BLOCK SIZE 4096"},

		// CLEAR LOGFILE with UNRECOVERABLE DATAFILE
		{"alter_db_clear_logfile_unrecoverable", "ALTER DATABASE CLEAR LOGFILE GROUP 3 UNRECOVERABLE DATAFILE"},

		// ADD LOGFILE with INSTANCE/THREAD
		{"alter_db_add_logfile_instance", "ALTER DATABASE ADD LOGFILE INSTANCE 'inst1' GROUP 5 '/u01/redo05.log' SIZE 100M"},
		{"alter_db_add_logfile_thread", "ALTER DATABASE ADD LOGFILE THREAD 2 GROUP 5 '/u01/redo05.log' SIZE 100M"},

		// DEFAULT LOCAL TEMPORARY TABLESPACE
		{"alter_db_default_local_temp", "ALTER DATABASE DEFAULT LOCAL TEMPORARY TABLESPACE localtemp"},

		// SET UNDO TABLESPACE
		{"alter_db_set_undo_tablespace", "ALTER DATABASE SET UNDO TABLESPACE = undots2"},

		// ENABLE/DISABLE LOST WRITE PROTECTION
		{"alter_db_enable_lost_write", "ALTER DATABASE ENABLE LOST WRITE PROTECTION"},
		{"alter_db_disable_lost_write", "ALTER DATABASE DISABLE LOST WRITE PROTECTION"},

		// PREPARE MIRROR COPY / DROP MIRROR COPY
		{"alter_db_prepare_mirror", "ALTER DATABASE PREPARE MIRROR COPY mirrorname"},
		{"alter_db_drop_mirror", "ALTER DATABASE DROP MIRROR COPY"},

		// ENABLE BLOCK CHANGE TRACKING USING FILE '+diskgroup'
		{"alter_db_enable_bct_asm", "ALTER DATABASE ENABLE BLOCK CHANGE TRACKING USING FILE '+DATA/bct.f'"},

		// CREATE CONTROLFILE: SET STANDBY NOLOGGING
		{"create_ctlfile_standby_nologging", "CREATE CONTROLFILE SET DATABASE mydb LOGFILE GROUP 1 '/u01/redo01.log' SIZE 100M RESETLOGS SET STANDBY NOLOGGING FOR DATA AVAILABILITY"},

		// recover with UNTIL CONSISTENT, LOGFILE, TEST, ALLOW CORRUPTION
		{"alter_db_recover_until_time", "ALTER DATABASE RECOVER DATABASE UNTIL TIME '2024-01-15:12:00:00'"},
		{"alter_db_recover_until_consistent", "ALTER DATABASE RECOVER DATABASE UNTIL CONSISTENT"},
		{"alter_db_recover_standby_until", "ALTER DATABASE RECOVER STANDBY DATABASE UNTIL CANCEL"},
		{"alter_db_recover_using_backup_ctl", "ALTER DATABASE RECOVER DATABASE USING BACKUP CONTROLFILE"},

		// managed standby with USING CURRENT LOGFILE, UNTIL CHANGE, etc.
		{"alter_db_managed_standby_using_current", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE USING CURRENT LOGFILE"},
		{"alter_db_managed_standby_until_change", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE UNTIL CHANGE 123456"},
		{"alter_db_managed_standby_nodelay", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE NODELAY"},
		{"alter_db_managed_standby_finish_force", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE FINISH FORCE"},

		// REGISTER LOGFILE from managed standby context
		{"alter_db_register_logfile_or_replace", "ALTER DATABASE REGISTER LOGFILE '/u01/arch01.log' OR REPLACE"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ParseAndCheck(t, tc.sql)
		})
	}
}

// TestParseInfrastructureSharedHelpers tests shared parser helpers
// that are used across all statement parsers (batch 87).
func TestParseInfrastructureSharedHelpers(t *testing.T) {
	// Test schema-qualified names in various statements
	tests := []struct {
		name string
		sql  string
	}{
		// Schema-qualified object names
		{"schema_qualified_select", `SELECT * FROM hr.employees`},
		{"schema_qualified_insert", `INSERT INTO hr.employees (id) VALUES (1)`},
		{"dblink_reference", `SELECT * FROM employees@remote_db`},
		{"schema_dblink_reference", `SELECT * FROM hr.employees@remote_db`},
		// Column references: simple, table-qualified, schema-qualified
		{"column_ref_simple", `SELECT name FROM employees`},
		{"column_ref_table_qualified", `SELECT e.name FROM employees e`},
		{"column_ref_schema_table_qualified", `SELECT hr.employees.name FROM hr.employees`},
		{"column_ref_star", `SELECT e.* FROM employees e`},
		// Bind variables
		{"bind_variable_named", `SELECT * FROM employees WHERE id = :emp_id`},
		{"bind_variable_positional", `SELECT * FROM employees WHERE id = :1`},
		// Pseudo-columns
		{"pseudo_rownum", `SELECT * FROM employees WHERE ROWNUM <= 10`},
		{"pseudo_rowid", `SELECT ROWID, name FROM employees`},
		{"pseudo_level", `SELECT LEVEL, name FROM employees CONNECT BY PRIOR id = manager_id`},
		{"pseudo_sysdate", `SELECT SYSDATE FROM DUAL`},
		{"pseudo_systimestamp", `SELECT SYSTIMESTAMP FROM DUAL`},
		// Expression basics
		{"expr_concatenation", `SELECT first_name || ' ' || last_name FROM employees`},
		{"expr_arithmetic", `SELECT salary * 1.1 + bonus FROM employees`},
		{"expr_exponent", `SELECT 2 ** 10 FROM DUAL`},
		{"expr_between", `SELECT * FROM employees WHERE salary BETWEEN 1000 AND 5000`},
		{"expr_in_list", `SELECT * FROM employees WHERE dept_id IN (10, 20, 30)`},
		{"expr_like", `SELECT * FROM employees WHERE name LIKE 'A%'`},
		{"expr_is_null", `SELECT * FROM employees WHERE manager_id IS NULL`},
		{"expr_is_not_null", `SELECT * FROM employees WHERE manager_id IS NOT NULL`},
		{"expr_case", `SELECT CASE WHEN salary > 5000 THEN 'high' ELSE 'low' END FROM employees`},
		// Multi-statement parsing
		{"multi_statement", `SELECT 1 FROM DUAL; SELECT 2 FROM DUAL`},
		// Keywords as identifiers
		{"keyword_as_column", `SELECT "TYPE", "NAME" FROM my_table`},
		// Type parsing in DDL context
		{"type_varchar2", `CREATE TABLE t (c VARCHAR2(100))`},
		{"type_number_precision", `CREATE TABLE t (c NUMBER(10, 2))`},
		{"type_timestamp_tz", `CREATE TABLE t (c TIMESTAMP WITH TIME ZONE)`},
		{"type_interval_ym", `CREATE TABLE t (c INTERVAL YEAR TO MONTH)`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() == 0 {
				t.Fatal("expected at least 1 statement")
			}
		})
	}
}

// TestParseInfrastructureLocTracking tests that AST nodes have Loc set.
func TestParseInfrastructureLocTracking(t *testing.T) {
	result := ParseAndCheck(t, `SELECT * FROM employees WHERE id = 1`)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	if raw.StmtLocation != 0 {
		t.Errorf("expected StmtLocation=0, got %d", raw.StmtLocation)
	}
	if raw.StmtLen <= 0 {
		t.Errorf("expected StmtLen>0, got %d", raw.StmtLen)
	}
}

// TestParseInfrastructureErrorRecovery tests that parse errors are reported properly.
func TestParseInfrastructureErrorRecovery(t *testing.T) {
	// Completely invalid SQL should produce an error
	_, err := Parse("!!! not valid sql !!!")
	if err == nil {
		t.Fatal("expected parse error for invalid SQL")
	}
}

// TestParseDMLReview tests INSERT/UPDATE/DELETE/MERGE statements against BNF.
func TestParseDMLReview(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// --- INSERT ---
		{
			name: "insert_basic_values",
			sql:  `INSERT INTO employees (employee_id, first_name, last_name) VALUES (1, 'John', 'Doe')`,
		},
		{
			name: "insert_subquery",
			sql:  `INSERT INTO emp_archive SELECT * FROM employees WHERE hire_date < TO_DATE('2020-01-01', 'YYYY-MM-DD')`,
		},
		{
			name: "insert_partition",
			sql:  `INSERT INTO sales PARTITION (sales_q1_2024) (prod_id, amount) VALUES (100, 500.00)`,
		},
		{
			name: "insert_partition_for",
			sql:  `INSERT INTO sales PARTITION FOR (TO_DATE('2024-03-15', 'YYYY-MM-DD')) (prod_id, amount) VALUES (100, 500.00)`,
		},
		{
			name: "insert_subpartition",
			sql:  `INSERT INTO sales SUBPARTITION (sales_q1_west) (prod_id, amount) VALUES (101, 250.00)`,
		},
		{
			name: "insert_dblink",
			sql:  `INSERT INTO employees@remote_db (employee_id, first_name) VALUES (999, 'Remote')`,
		},
		{
			name: "insert_with_alias",
			sql:  `INSERT INTO employees e (employee_id) VALUES (1)`,
		},
		{
			name: "insert_returning",
			sql:  `INSERT INTO employees (employee_id, first_name) VALUES (1, 'John') RETURNING employee_id INTO :id`,
		},
		{
			name: "insert_return_keyword",
			sql:  `INSERT INTO employees (employee_id) VALUES (1) RETURN employee_id INTO :id`,
		},
		{
			name: "insert_error_logging",
			sql:  `INSERT INTO employees (employee_id) VALUES (1) LOG ERRORS INTO err_employees ('batch1') REJECT LIMIT 100`,
		},
		{
			name: "insert_error_logging_unlimited",
			sql:  `INSERT INTO employees (employee_id) VALUES (1) LOG ERRORS REJECT LIMIT UNLIMITED`,
		},
		{
			name: "insert_all_unconditional",
			sql: `INSERT ALL
				INTO sal_history (empid, hire_date, salary) VALUES (empno, hiredate, sal)
				INTO mgr_history (empid, manager, salary) VALUES (empno, mgr, sal)
				SELECT empno, hiredate, sal, mgr FROM emp WHERE deptno = 10`,
		},
		{
			name: "insert_first_conditional",
			sql: `INSERT FIRST
				WHEN sal > 10000 THEN INTO high_sal (empid, salary) VALUES (empno, sal)
				WHEN sal > 5000 THEN INTO mid_sal (empid, salary) VALUES (empno, sal)
				ELSE INTO low_sal (empid, salary) VALUES (empno, sal)
				SELECT empno, sal FROM emp`,
		},
		{
			name: "insert_all_conditional",
			sql: `INSERT ALL
				WHEN deptno = 10 THEN INTO dept10 (empid) VALUES (empno)
				WHEN deptno = 20 THEN INTO dept20 (empid) VALUES (empno)
				SELECT empno, deptno FROM emp`,
		},
		{
			name: "insert_hints",
			sql:  `INSERT /*+ APPEND */ INTO employees (employee_id) VALUES (1)`,
		},
		{
			name: "insert_set_clause",
			sql:  `INSERT INTO employees SET employee_id = 1, first_name = 'John', last_name = 'Doe'`,
		},
		{
			name: "insert_by_name",
			sql:  `INSERT INTO target_table BY NAME SELECT * FROM source_table`,
		},

		// --- UPDATE ---
		{
			name: "update_basic",
			sql:  `UPDATE employees SET salary = 50000 WHERE employee_id = 100`,
		},
		{
			name: "update_multi_column",
			sql:  `UPDATE employees SET salary = 50000, commission_pct = 0.1 WHERE department_id = 80`,
		},
		{
			name: "update_subquery_value",
			sql:  `UPDATE employees SET (salary, commission_pct) = (SELECT avg_sal, avg_comm FROM dept_avg WHERE dept_id = 80) WHERE department_id = 80`,
		},
		{
			name: "update_partition",
			sql:  `UPDATE employees PARTITION (p_2024) SET salary = salary * 1.1`,
		},
		{
			name: "update_dblink",
			sql:  `UPDATE employees@remote_db SET salary = 50000 WHERE employee_id = 100`,
		},
		{
			name: "update_returning",
			sql:  `UPDATE employees SET salary = 50000 WHERE employee_id = 100 RETURNING salary INTO :new_sal`,
		},
		{
			name: "update_return_keyword",
			sql:  `UPDATE employees SET salary = 50000 WHERE employee_id = 100 RETURN salary INTO :new_sal`,
		},
		{
			name: "update_error_logging",
			sql:  `UPDATE employees SET salary = 50000 WHERE department_id = 80 LOG ERRORS INTO err_employees ('update_batch') REJECT LIMIT 50`,
		},
		{
			name: "update_set_default",
			sql:  `UPDATE employees SET commission_pct = DEFAULT WHERE department_id = 10`,
		},
		{
			name: "update_with_alias",
			sql:  `UPDATE employees e SET e.salary = 50000 WHERE e.employee_id = 100`,
		},
		{
			name: "update_with_hints",
			sql:  `UPDATE /*+ PARALLEL(employees, 4) */ employees SET salary = salary * 1.1`,
		},

		// --- DELETE ---
		{
			name: "delete_basic",
			sql:  `DELETE FROM employees WHERE employee_id = 100`,
		},
		{
			name: "delete_without_from",
			sql:  `DELETE employees WHERE employee_id = 100`,
		},
		{
			name: "delete_partition",
			sql:  `DELETE FROM sales PARTITION (sales_q1_2024) WHERE amount < 10`,
		},
		{
			name: "delete_subpartition",
			sql:  `DELETE FROM sales SUBPARTITION (sales_q1_west)`,
		},
		{
			name: "delete_dblink",
			sql:  `DELETE FROM employees@remote_db WHERE employee_id = 999`,
		},
		{
			name: "delete_returning",
			sql:  `DELETE FROM employees WHERE employee_id = 100 RETURNING first_name, last_name INTO :fname, :lname`,
		},
		{
			name: "delete_error_logging",
			sql:  `DELETE FROM employees WHERE department_id = 80 LOG ERRORS INTO err_employees REJECT LIMIT UNLIMITED`,
		},
		{
			name: "delete_with_alias",
			sql:  `DELETE FROM employees e WHERE e.department_id = 10`,
		},
		{
			name: "delete_hints",
			sql:  `DELETE /*+ PARALLEL */ FROM employees WHERE hire_date < SYSDATE - 3650`,
		},

		// --- MERGE ---
		{
			name: "merge_basic",
			sql: `MERGE INTO target t USING source s ON (t.id = s.id)
				WHEN MATCHED THEN UPDATE SET t.val = s.val
				WHEN NOT MATCHED THEN INSERT (id, val) VALUES (s.id, s.val)`,
		},
		{
			name: "merge_update_where",
			sql: `MERGE INTO target t USING source s ON (t.id = s.id)
				WHEN MATCHED THEN UPDATE SET t.val = s.val WHERE t.status = 'ACTIVE'`,
		},
		{
			name: "merge_update_delete_where",
			sql: `MERGE INTO target t USING source s ON (t.id = s.id)
				WHEN MATCHED THEN UPDATE SET t.val = s.val DELETE WHERE t.status = 'DELETED'`,
		},
		{
			name: "merge_insert_where",
			sql: `MERGE INTO target t USING source s ON (t.id = s.id)
				WHEN NOT MATCHED THEN INSERT (id, val) VALUES (s.id, s.val) WHERE s.active = 1`,
		},
		{
			name: "merge_error_logging",
			sql: `MERGE INTO target t USING source s ON (t.id = s.id)
				WHEN MATCHED THEN UPDATE SET t.val = s.val
				WHEN NOT MATCHED THEN INSERT (id, val) VALUES (s.id, s.val)
				LOG ERRORS INTO merge_errors ('batch1') REJECT LIMIT 200`,
		},
		{
			name: "merge_using_subquery",
			sql: `MERGE INTO employees e
				USING (SELECT id, new_sal FROM salary_changes) sc ON (e.employee_id = sc.id)
				WHEN MATCHED THEN UPDATE SET e.salary = sc.new_sal`,
		},
		{
			name: "merge_with_hints",
			sql: `MERGE /*+ PARALLEL(4) */ INTO target t USING source s ON (t.id = s.id)
				WHEN MATCHED THEN UPDATE SET t.val = s.val`,
		},
		{
			name: "merge_set_default",
			sql: `MERGE INTO target t USING source s ON (t.id = s.id)
				WHEN MATCHED THEN UPDATE SET t.val = DEFAULT`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("expected at least 1 statement, got 0")
			}
			_ = ast.NodeToString(result.Items[0])
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 90: CREATE TABLE full (774-line BNF review)
// ---------------------------------------------------------------------------

func TestParseCreateTableFull(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// Basic CREATE TABLE
		{
			name: "basic_table",
			sql:  "CREATE TABLE employees (id NUMBER, name VARCHAR2(100))",
		},
		// OR REPLACE
		{
			name: "or_replace",
			sql:  "CREATE OR REPLACE TABLE employees (id NUMBER)",
		},
		// GLOBAL TEMPORARY
		{
			name: "global_temporary",
			sql:  "CREATE GLOBAL TEMPORARY TABLE temp_data (id NUMBER) ON COMMIT PRESERVE ROWS",
		},
		// PRIVATE TEMPORARY
		{
			name: "private_temporary",
			sql:  "CREATE PRIVATE TEMPORARY TABLE ora$ptt_temp (id NUMBER) ON COMMIT DELETE ROWS",
		},
		// SHARDED
		{
			name: "sharded_table",
			sql:  "CREATE SHARDED TABLE customers (id NUMBER, name VARCHAR2(200))",
		},
		// DUPLICATED
		{
			name: "duplicated_table",
			sql:  "CREATE DUPLICATED TABLE ref_data (code VARCHAR2(10), descr VARCHAR2(100))",
		},
		// IF NOT EXISTS
		{
			name: "if_not_exists",
			sql:  "CREATE TABLE IF NOT EXISTS audit_log (id NUMBER, msg CLOB)",
		},
		// SHARING clause
		{
			name: "sharing_metadata",
			sql:  "CREATE TABLE hr.employees SHARING = METADATA (id NUMBER, name VARCHAR2(100))",
		},
		{
			name: "sharing_data",
			sql:  "CREATE TABLE hr.employees SHARING = DATA (id NUMBER)",
		},
		{
			name: "sharing_extended_data",
			sql:  "CREATE TABLE hr.employees SHARING = EXTENDED DATA (id NUMBER)",
		},
		{
			name: "sharing_none",
			sql:  "CREATE TABLE hr.employees SHARING = NONE (id NUMBER)",
		},
		// CTAS
		{
			name: "ctas",
			sql:  "CREATE TABLE emp_copy AS SELECT * FROM employees",
		},
		// Column with SORT
		{
			name: "column_sort",
			sql:  "CREATE TABLE t1 (id NUMBER SORT)",
		},
		// Column with VISIBLE
		{
			name: "column_visible",
			sql:  "CREATE TABLE t1 (id NUMBER VISIBLE, secret VARCHAR2(100) INVISIBLE)",
		},
		// Column with DEFAULT ON NULL
		{
			name: "default_on_null",
			sql:  "CREATE TABLE t1 (id NUMBER, status VARCHAR2(10) DEFAULT ON NULL 'ACTIVE')",
		},
		// Column with DEFAULT ON NULL FOR INSERT ONLY
		{
			name: "default_on_null_insert_only",
			sql:  "CREATE TABLE t1 (id NUMBER, status VARCHAR2(10) DEFAULT ON NULL FOR INSERT ONLY 'ACTIVE')",
		},
		// Column with DOMAIN
		{
			name: "column_domain",
			sql:  "CREATE TABLE t1 (email DOMAIN hr.email_domain, name VARCHAR2(100))",
		},
		// Identity column: GENERATED ALWAYS
		{
			name: "identity_always",
			sql:  "CREATE TABLE t1 (id NUMBER GENERATED ALWAYS AS IDENTITY, name VARCHAR2(100))",
		},
		// Identity column: GENERATED BY DEFAULT
		{
			name: "identity_by_default",
			sql:  "CREATE TABLE t1 (id NUMBER GENERATED BY DEFAULT AS IDENTITY, name VARCHAR2(100))",
		},
		// Identity column: GENERATED BY DEFAULT ON NULL
		{
			name: "identity_by_default_on_null",
			sql:  "CREATE TABLE t1 (id NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY, name VARCHAR2(100))",
		},
		// Identity column with options
		{
			name: "identity_with_options",
			sql:  "CREATE TABLE t1 (id NUMBER GENERATED ALWAYS AS IDENTITY (START WITH 100 INCREMENT BY 10 MAXVALUE 999999 NOCYCLE CACHE 20))",
		},
		// Virtual column
		{
			name: "virtual_column",
			sql:  "CREATE TABLE t1 (price NUMBER, qty NUMBER, total NUMBER GENERATED ALWAYS AS (price * qty) VIRTUAL)",
		},
		// Column with ENCRYPT
		{
			name: "column_encrypt",
			sql:  "CREATE TABLE t1 (id NUMBER, ssn VARCHAR2(11) ENCRYPT)",
		},
		// Column with ENCRYPT spec
		{
			name: "column_encrypt_spec",
			sql:  "CREATE TABLE t1 (id NUMBER, ssn VARCHAR2(11) ENCRYPT USING 'AES256')",
		},
		// Column with COLLATE
		{
			name: "column_collate",
			sql:  "CREATE TABLE t1 (id NUMBER, name VARCHAR2(100) COLLATE BINARY_CI)",
		},
		// Multiple column constraints
		{
			name: "multiple_constraints",
			sql:  "CREATE TABLE t1 (id NUMBER PRIMARY KEY, email VARCHAR2(200) NOT NULL UNIQUE, dept_id NUMBER REFERENCES departments(id))",
		},
		// Table-level constraints
		{
			name: "table_constraints",
			sql: `CREATE TABLE t1 (id NUMBER, name VARCHAR2(100), dept_id NUMBER,
				CONSTRAINT pk_t1 PRIMARY KEY (id),
				CONSTRAINT uq_name UNIQUE (name),
				CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments(id) ON DELETE CASCADE,
				CHECK (id > 0))`,
		},
		// TABLESPACE
		{
			name: "tablespace",
			sql:  "CREATE TABLE t1 (id NUMBER) TABLESPACE users",
		},
		// ORGANIZATION HEAP
		{
			name: "org_heap",
			sql:  "CREATE TABLE t1 (id NUMBER) ORGANIZATION HEAP",
		},
		// ORGANIZATION INDEX
		{
			name: "org_index",
			sql:  "CREATE TABLE t1 (id NUMBER PRIMARY KEY) ORGANIZATION INDEX",
		},
		// SEGMENT CREATION IMMEDIATE
		{
			name: "segment_creation_immediate",
			sql:  "CREATE TABLE t1 (id NUMBER) SEGMENT CREATION IMMEDIATE",
		},
		// SEGMENT CREATION DEFERRED
		{
			name: "segment_creation_deferred",
			sql:  "CREATE TABLE t1 (id NUMBER) SEGMENT CREATION DEFERRED",
		},
		// LOGGING / NOLOGGING
		{
			name: "logging",
			sql:  "CREATE TABLE t1 (id NUMBER) LOGGING",
		},
		{
			name: "nologging",
			sql:  "CREATE TABLE t1 (id NUMBER) NOLOGGING",
		},
		// CACHE / NOCACHE
		{
			name: "cache",
			sql:  "CREATE TABLE t1 (id NUMBER) CACHE",
		},
		{
			name: "nocache",
			sql:  "CREATE TABLE t1 (id NUMBER) NOCACHE",
		},
		// PARALLEL / NOPARALLEL
		{
			name: "parallel",
			sql:  "CREATE TABLE t1 (id NUMBER) PARALLEL",
		},
		{
			name: "noparallel",
			sql:  "CREATE TABLE t1 (id NUMBER) NOPARALLEL",
		},
		// COMPRESS / NOCOMPRESS
		{
			name: "compress",
			sql:  "CREATE TABLE t1 (id NUMBER) COMPRESS",
		},
		{
			name: "nocompress",
			sql:  "CREATE TABLE t1 (id NUMBER) NOCOMPRESS",
		},
		// ROW STORE COMPRESS
		{
			name: "row_store_compress_advanced",
			sql:  "CREATE TABLE t1 (id NUMBER) ROW STORE COMPRESS ADVANCED",
		},
		// READ ONLY / READ WRITE
		{
			name: "read_only",
			sql:  "CREATE TABLE t1 (id NUMBER) READ ONLY",
		},
		{
			name: "read_write",
			sql:  "CREATE TABLE t1 (id NUMBER) READ WRITE",
		},
		// INDEXING ON / OFF
		{
			name: "indexing_on",
			sql:  "CREATE TABLE t1 (id NUMBER) INDEXING ON",
		},
		{
			name: "indexing_off",
			sql:  "CREATE TABLE t1 (id NUMBER) INDEXING OFF",
		},
		// RESULT_CACHE
		{
			name: "result_cache_default",
			sql:  "CREATE TABLE t1 (id NUMBER) RESULT_CACHE (MODE DEFAULT)",
		},
		{
			name: "result_cache_force",
			sql:  "CREATE TABLE t1 (id NUMBER) RESULT_CACHE (MODE FORCE)",
		},
		// ROWDEPENDENCIES / NOROWDEPENDENCIES
		{
			name: "rowdependencies",
			sql:  "CREATE TABLE t1 (id NUMBER) ROWDEPENDENCIES",
		},
		{
			name: "norowdependencies",
			sql:  "CREATE TABLE t1 (id NUMBER) NOROWDEPENDENCIES",
		},
		// ROW MOVEMENT
		{
			name: "enable_row_movement",
			sql:  "CREATE TABLE t1 (id NUMBER) ENABLE ROW MOVEMENT",
		},
		{
			name: "disable_row_movement",
			sql:  "CREATE TABLE t1 (id NUMBER) DISABLE ROW MOVEMENT",
		},
		// FLASHBACK ARCHIVE
		{
			name: "flashback_archive",
			sql:  "CREATE TABLE t1 (id NUMBER) FLASHBACK ARCHIVE fba1",
		},
		{
			name: "no_flashback_archive",
			sql:  "CREATE TABLE t1 (id NUMBER) NO FLASHBACK ARCHIVE",
		},
		// ENABLE / DISABLE constraint
		{
			name: "enable_primary_key",
			sql:  "CREATE TABLE t1 (id NUMBER, CONSTRAINT pk_t1 PRIMARY KEY (id)) ENABLE PRIMARY KEY",
		},
		{
			name: "disable_constraint",
			sql:  "CREATE TABLE t1 (id NUMBER, CONSTRAINT ck_t1 CHECK (id > 0)) DISABLE CONSTRAINT ck_t1",
		},
		// PARTITION BY RANGE
		{
			name: "partition_by_range",
			sql: `CREATE TABLE sales (id NUMBER, sale_date DATE, amount NUMBER)
				PARTITION BY RANGE (sale_date)
				(PARTITION p2023 VALUES LESS THAN (TO_DATE('2024-01-01','YYYY-MM-DD')),
				 PARTITION p2024 VALUES LESS THAN (MAXVALUE))`,
		},
		// PARTITION BY LIST
		{
			name: "partition_by_list",
			sql: `CREATE TABLE customers (id NUMBER, region VARCHAR2(20))
				PARTITION BY LIST (region)
				(PARTITION p_east VALUES ('East', 'Northeast'),
				 PARTITION p_west VALUES ('West'))`,
		},
		// PARTITION BY HASH
		{
			name: "partition_by_hash",
			sql: `CREATE TABLE orders (id NUMBER, customer_id NUMBER)
				PARTITION BY HASH (customer_id)
				(PARTITION p1, PARTITION p2, PARTITION p3, PARTITION p4)`,
		},
		// DEFAULT COLLATION
		{
			name: "default_collation",
			sql:  "CREATE TABLE t1 (id NUMBER, name VARCHAR2(100)) DEFAULT COLLATION USING_NLS_COMP",
		},
		// Immutable table: NO DROP
		{
			name: "immutable_no_drop",
			sql:  "CREATE TABLE t1 (id NUMBER) NO DROP",
		},
		{
			name: "immutable_no_drop_until",
			sql:  "CREATE TABLE t1 (id NUMBER) NO DROP UNTIL 30 DAYS IDLE",
		},
		// Immutable table: NO DELETE
		{
			name: "immutable_no_delete_locked",
			sql:  "CREATE TABLE t1 (id NUMBER) NO DELETE LOCKED",
		},
		{
			name: "immutable_no_delete_until",
			sql:  "CREATE TABLE t1 (id NUMBER) NO DELETE UNTIL 365 DAYS AFTER INSERT",
		},
		// Immutable table: NO DROP + NO DELETE combined
		{
			name: "immutable_no_drop_no_delete",
			sql:  "CREATE TABLE t1 (id NUMBER) NO DROP UNTIL 16 DAYS IDLE NO DELETE UNTIL 16 DAYS AFTER INSERT",
		},
		// Blockchain: HASHING USING / VERSION
		{
			name: "blockchain_hash_version",
			sql:  "CREATE TABLE t1 (id NUMBER) NO DROP UNTIL 16 DAYS IDLE NO DELETE UNTIL 16 DAYS AFTER INSERT HASHING USING 'SHA2_512' VERSION 'v2'",
		},
		// MEMOPTIMIZE FOR READ
		{
			name: "memoptimize_read",
			sql:  "CREATE TABLE t1 (id NUMBER) MEMOPTIMIZE FOR READ",
		},
		// MEMOPTIMIZE FOR WRITE
		{
			name: "memoptimize_write",
			sql:  "CREATE TABLE t1 (id NUMBER) MEMOPTIMIZE FOR WRITE",
		},
		// MEMOPTIMIZE FOR READ + WRITE
		{
			name: "memoptimize_read_write",
			sql:  "CREATE TABLE t1 (id NUMBER) MEMOPTIMIZE FOR READ MEMOPTIMIZE FOR WRITE",
		},
		// PARENT
		{
			name: "parent_table",
			sql:  "CREATE TABLE child_tab (id NUMBER) PARENT hr.parent_tab",
		},
		// PCTFREE / PCTUSED / INITRANS
		{
			name: "physical_attrs",
			sql:  "CREATE TABLE t1 (id NUMBER) PCTFREE 20 PCTUSED 40 INITRANS 2 TABLESPACE users",
		},
		// Combined options
		{
			name: "combined_options",
			sql: `CREATE TABLE hr.employees (
				id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
				name VARCHAR2(200) NOT NULL,
				email VARCHAR2(200) UNIQUE,
				salary NUMBER DEFAULT 0 CHECK (salary >= 0),
				dept_id NUMBER REFERENCES departments(id) ON DELETE SET NULL
			) TABLESPACE hr_data
			  SEGMENT CREATION IMMEDIATE
			  LOGGING
			  CACHE
			  PARALLEL
			  COMPRESS
			  INDEXING ON
			  ENABLE ROW MOVEMENT`,
		},
		// Deferrable constraint
		{
			name: "deferrable_constraint",
			sql:  "CREATE TABLE t1 (id NUMBER CONSTRAINT pk_t1 PRIMARY KEY DEFERRABLE INITIALLY DEFERRED, name VARCHAR2(100) NOT NULL)",
		},
		// Check constraint
		{
			name: "check_constraint_inline",
			sql:  "CREATE TABLE t1 (id NUMBER, age NUMBER CHECK (age >= 0 AND age <= 150))",
		},
		// FOREIGN KEY with ON DELETE
		{
			name: "foreign_key_cascade",
			sql:  "CREATE TABLE t1 (id NUMBER, parent_id NUMBER REFERENCES t2(id) ON DELETE CASCADE)",
		},
		{
			name: "foreign_key_set_null",
			sql:  "CREATE TABLE t1 (id NUMBER, parent_id NUMBER REFERENCES t2(id) ON DELETE SET NULL)",
		},
		// Named column constraints
		{
			name: "named_column_constraints",
			sql:  "CREATE TABLE t1 (id NUMBER CONSTRAINT nn_id NOT NULL CONSTRAINT pk_id PRIMARY KEY)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("expected at least 1 statement, got 0")
			}
			_ = ast.NodeToString(result.Items[0])
		})
	}
}

func TestParseAlterTableFull(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// --- Column clauses ---
		// add_column_clause
		{"add_column_single", `ALTER TABLE employees ADD (email VARCHAR2(100))`},
		{"add_column_multiple", `ALTER TABLE employees ADD (first_name VARCHAR2(50), last_name VARCHAR2(50) NOT NULL)`},
		{"add_column_no_parens", `ALTER TABLE employees ADD email VARCHAR2(100)`},
		{"add_column_default", `ALTER TABLE employees ADD (hire_date DATE DEFAULT SYSDATE)`},
		{"add_column_identity", `ALTER TABLE employees ADD (id NUMBER GENERATED ALWAYS AS IDENTITY)`},
		{"add_column_identity_options", `ALTER TABLE employees ADD (id NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY (START WITH 1 INCREMENT BY 1 MAXVALUE 99999 NOCYCLE CACHE 20))`},
		{"add_column_encrypt", `ALTER TABLE employees ADD (ssn VARCHAR2(11) ENCRYPT)`},
		{"add_column_invisible", `ALTER TABLE employees ADD (internal_code NUMBER INVISIBLE)`},
		{"add_column_virtual", `ALTER TABLE employees ADD (full_name VARCHAR2(100) GENERATED ALWAYS AS (first_name || ' ' || last_name) VIRTUAL)`},

		// modify_column_clauses
		{"modify_column_type", `ALTER TABLE employees MODIFY (salary NUMBER(12,2))`},
		{"modify_column_default", `ALTER TABLE employees MODIFY (status VARCHAR2(10) DEFAULT 'ACTIVE')`},
		{"modify_column_null", `ALTER TABLE employees MODIFY (email VARCHAR2(200) NOT NULL)`},
		{"modify_column_multiple", `ALTER TABLE employees MODIFY (first_name VARCHAR2(100), last_name VARCHAR2(100))`},
		{"modify_column_visibility", `ALTER TABLE employees MODIFY email VISIBLE`},
		{"modify_column_invisible", `ALTER TABLE employees MODIFY internal_code INVISIBLE`},
		{"modify_column_drop_identity", `ALTER TABLE employees MODIFY (id DROP IDENTITY)`},

		// drop_column_clause
		{"drop_column", `ALTER TABLE employees DROP COLUMN email`},
		{"drop_column_cascade", `ALTER TABLE employees DROP COLUMN dept_id CASCADE CONSTRAINTS`},
		{"drop_columns_parens", `ALTER TABLE employees DROP (col1, col2, col3)`},
		{"drop_column_checkpoint", `ALTER TABLE employees DROP COLUMN old_data CASCADE CONSTRAINTS CHECKPOINT 1000`},
		{"drop_unused_columns", `ALTER TABLE employees DROP UNUSED COLUMNS`},
		{"drop_columns_continue", `ALTER TABLE employees DROP COLUMNS CONTINUE`},
		{"drop_unused_checkpoint", `ALTER TABLE employees DROP UNUSED COLUMNS CHECKPOINT 500`},

		// SET UNUSED
		{"set_unused_column", `ALTER TABLE employees SET UNUSED COLUMN old_data`},
		{"set_unused_columns", `ALTER TABLE employees SET UNUSED (col1, col2)`},
		{"set_unused_cascade", `ALTER TABLE employees SET UNUSED COLUMN old_data CASCADE CONSTRAINTS`},
		{"set_unused_online", `ALTER TABLE employees SET UNUSED (col1) ONLINE`},

		// rename_column_clause
		{"rename_column", `ALTER TABLE employees RENAME COLUMN old_name TO new_name`},

		// PERIOD clauses
		{"add_period", `ALTER TABLE employees ADD (PERIOD FOR valid_time (start_date, end_date))`},
		{"add_period_no_cols", `ALTER TABLE employees ADD (PERIOD FOR valid_time)`},
		{"drop_period", `ALTER TABLE employees DROP (PERIOD FOR valid_time)`},

		// MODIFY NESTED TABLE / LOB / VARRAY
		{"modify_nested_table", `ALTER TABLE departments MODIFY NESTED TABLE employees_nt RETURN AS LOCATOR`},
		{"modify_nested_table_value", `ALTER TABLE departments MODIFY NESTED TABLE items RETURN AS VALUE`},
		{"modify_lob", `ALTER TABLE documents MODIFY LOB (content) (CACHE)`},
		{"modify_varray", `ALTER TABLE employees MODIFY VARRAY phone_numbers (CACHE)`},

		// --- Constraint clauses ---
		{"add_constraint_pk", `ALTER TABLE employees ADD CONSTRAINT pk_emp PRIMARY KEY (employee_id)`},
		{"add_constraint_fk", `ALTER TABLE employees ADD CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments (dept_id)`},
		{"add_constraint_unique", `ALTER TABLE employees ADD CONSTRAINT uq_email UNIQUE (email)`},
		{"add_constraint_check", `ALTER TABLE employees ADD CONSTRAINT ck_salary CHECK (salary > 0)`},
		{"modify_constraint", `ALTER TABLE employees MODIFY CONSTRAINT pk_emp ENABLE VALIDATE`},
		{"modify_primary_key", `ALTER TABLE employees MODIFY PRIMARY KEY DISABLE`},
		{"modify_unique", `ALTER TABLE employees MODIFY UNIQUE (email) NOVALIDATE`},
		{"drop_constraint", `ALTER TABLE employees DROP CONSTRAINT fk_dept CASCADE`},
		{"drop_primary_key", `ALTER TABLE employees DROP PRIMARY KEY CASCADE`},
		{"drop_unique", `ALTER TABLE employees DROP UNIQUE (email)`},
		{"drop_pk_keep_index", `ALTER TABLE employees DROP PRIMARY KEY KEEP INDEX`},
		{"drop_constraint_online", `ALTER TABLE employees DROP CONSTRAINT ck_sal ONLINE`},
		{"rename_constraint", `ALTER TABLE employees RENAME CONSTRAINT old_ck TO new_ck`},

		// --- Enable/Disable clauses ---
		{"enable_constraint", `ALTER TABLE employees ENABLE CONSTRAINT ck_salary`},
		{"disable_constraint", `ALTER TABLE employees DISABLE CONSTRAINT fk_dept`},
		{"enable_pk", `ALTER TABLE employees ENABLE PRIMARY KEY`},
		{"disable_unique", `ALTER TABLE employees DISABLE UNIQUE (email)`},
		{"enable_validate", `ALTER TABLE employees ENABLE VALIDATE CONSTRAINT ck_salary`},
		{"enable_novalidate", `ALTER TABLE employees ENABLE NOVALIDATE PRIMARY KEY`},
		{"enable_table_lock", `ALTER TABLE employees ENABLE TABLE LOCK`},
		{"disable_table_lock", `ALTER TABLE employees DISABLE TABLE LOCK`},
		{"enable_all_triggers", `ALTER TABLE employees ENABLE ALL TRIGGERS`},
		{"disable_all_triggers", `ALTER TABLE employees DISABLE ALL TRIGGERS`},
		{"enable_row_movement", `ALTER TABLE employees ENABLE ROW MOVEMENT`},
		{"disable_row_movement", `ALTER TABLE employees DISABLE ROW MOVEMENT`},

		// --- Alter table properties ---
		{"logging", `ALTER TABLE employees LOGGING`},
		{"nologging", `ALTER TABLE employees NOLOGGING`},
		{"compress", `ALTER TABLE employees COMPRESS`},
		{"compress_basic", `ALTER TABLE employees COMPRESS BASIC`},
		{"nocompress", `ALTER TABLE employees NOCOMPRESS`},
		{"row_store_compress", `ALTER TABLE employees ROW STORE COMPRESS ADVANCED`},
		{"parallel", `ALTER TABLE employees PARALLEL 4`},
		{"noparallel", `ALTER TABLE employees NOPARALLEL`},
		{"cache", `ALTER TABLE employees CACHE`},
		{"pctfree", `ALTER TABLE employees PCTFREE 20`},
		{"initrans", `ALTER TABLE employees INITRANS 4`},
		{"storage", `ALTER TABLE employees STORAGE (INITIAL 64K NEXT 64K)`},
		{"read_only", `ALTER TABLE employees READ ONLY`},
		{"read_write", `ALTER TABLE employees READ WRITE`},
		{"result_cache", `ALTER TABLE employees RESULT_CACHE (MODE FORCE)`},
		{"rename_to", `ALTER TABLE employees RENAME TO staff`},

		// --- Flashback archive ---
		{"flashback_archive", `ALTER TABLE employees FLASHBACK ARCHIVE fba1`},
		{"flashback_archive_default", `ALTER TABLE employees FLASHBACK ARCHIVE`},
		{"no_flashback_archive", `ALTER TABLE employees NO FLASHBACK ARCHIVE`},

		// --- SHRINK SPACE ---
		{"shrink_space", `ALTER TABLE employees SHRINK SPACE`},
		{"shrink_space_cascade", `ALTER TABLE employees SHRINK SPACE CASCADE`},

		// --- INMEMORY ---
		{"inmemory", `ALTER TABLE employees INMEMORY`},
		{"no_inmemory", `ALTER TABLE employees NO INMEMORY`},

		// --- MEMOPTIMIZE ---
		{"memoptimize_read", `ALTER TABLE employees MEMOPTIMIZE FOR READ`},
		{"no_memoptimize_read", `ALTER TABLE employees NO MEMOPTIMIZE FOR READ`},
		{"memoptimize_write", `ALTER TABLE employees MEMOPTIMIZE FOR WRITE`},

		// --- ALLOCATE / DEALLOCATE ---
		{"allocate_extent", `ALTER TABLE employees ALLOCATE EXTENT`},
		{"allocate_extent_size", `ALTER TABLE employees ALLOCATE EXTENT (SIZE 100K)`},
		{"deallocate_unused", `ALTER TABLE employees DEALLOCATE UNUSED`},
		{"deallocate_unused_keep", `ALTER TABLE employees DEALLOCATE UNUSED KEEP 100K`},

		// --- REFRESH ---
		{"refresh", `ALTER TABLE employees REFRESH`},
		{"no_refresh", `ALTER TABLE employees NO REFRESH`},

		// --- Immutable / Blockchain table clauses ---
		{"no_drop", `ALTER TABLE ledger NO DROP`},
		{"no_delete", `ALTER TABLE ledger NO DELETE`},

		// --- Partitioning operations ---
		{"add_partition", `ALTER TABLE sales ADD PARTITION sales_q5 VALUES LESS THAN (TO_DATE('2025-04-01', 'YYYY-MM-DD'))`},
		{"drop_partition", `ALTER TABLE sales DROP PARTITION sales_q1`},
		{"drop_partition_for", `ALTER TABLE sales DROP PARTITION FOR (TO_DATE('2024-01-15', 'YYYY-MM-DD'))`},
		{"drop_subpartition", `ALTER TABLE sales DROP SUBPARTITION sales_q1_east`},
		{"truncate_partition", `ALTER TABLE sales TRUNCATE PARTITION sales_q1`},
		{"truncate_partition_storage", `ALTER TABLE sales TRUNCATE PARTITION sales_q1 DROP STORAGE`},
		{"truncate_subpartition", `ALTER TABLE sales TRUNCATE SUBPARTITION sales_q1_east`},
		{"truncate_partition_cascade", `ALTER TABLE sales TRUNCATE PARTITION old_data REUSE STORAGE CASCADE`},
		{"split_partition", `ALTER TABLE sales SPLIT PARTITION sales_q1 AT (TO_DATE('2024-02-01', 'YYYY-MM-DD')) INTO (PARTITION sales_jan, PARTITION sales_feb_mar)`},
		{"merge_partitions", `ALTER TABLE sales MERGE PARTITIONS sales_q1, sales_q2 INTO PARTITION sales_h1`},
		{"exchange_partition", `ALTER TABLE sales EXCHANGE PARTITION sales_q1 WITH TABLE sales_staging`},
		{"exchange_partition_opts", `ALTER TABLE sales EXCHANGE PARTITION sales_q1 WITH TABLE staging INCLUDING INDEXES WITHOUT VALIDATION`},
		{"coalesce_partition", `ALTER TABLE hash_tab COALESCE PARTITION`},
		{"coalesce_subpartition", `ALTER TABLE hash_tab COALESCE SUBPARTITION`},
		{"modify_partition", `ALTER TABLE sales MODIFY PARTITION sales_q1 PCTFREE 30`},
		{"modify_subpartition", `ALTER TABLE sales MODIFY SUBPARTITION sales_q1_east TABLESPACE ts2`},
		{"modify_default_attrs", `ALTER TABLE sales MODIFY DEFAULT ATTRIBUTES TABLESPACE ts_new`},
		{"rename_partition", `ALTER TABLE sales RENAME PARTITION old_part TO new_part`},
		{"rename_subpartition", `ALTER TABLE sales RENAME SUBPARTITION old_sub TO new_sub`},

		// --- SET operations ---
		{"set_interval", `ALTER TABLE sales SET INTERVAL (NUMTOYMINTERVAL(1, 'MONTH'))`},
		{"set_partitioning_auto", `ALTER TABLE sales SET PARTITIONING AUTOMATIC`},
		{"set_partitioning_manual", `ALTER TABLE sales SET PARTITIONING MANUAL`},

		// --- MOVE ---
		{"move_table", `ALTER TABLE employees MOVE`},
		{"move_table_online", `ALTER TABLE employees MOVE ONLINE`},
		{"move_table_tablespace", `ALTER TABLE employees MOVE TABLESPACE ts_new`},
		{"move_partition", `ALTER TABLE sales MOVE PARTITION sales_q1 TABLESPACE ts_archive`},
		{"move_subpartition", `ALTER TABLE sales MOVE SUBPARTITION sales_q1_east TABLESPACE ts_archive`},

		// --- Supplemental log ---
		{"add_supplemental_log", `ALTER TABLE employees ADD SUPPLEMENTAL LOG DATA (ALL) COLUMNS`},
		{"add_supplemental_log_pk", `ALTER TABLE employees ADD SUPPLEMENTAL LOG DATA (PRIMARY KEY) COLUMNS`},
		{"drop_supplemental_log", `ALTER TABLE employees DROP SUPPLEMENTAL LOG DATA (ALL) COLUMNS`},

		// --- Multiple actions ---
		{"multi_action", `ALTER TABLE employees ADD (bonus NUMBER) MODIFY (salary NUMBER(12,2))`},
		{"multi_enable_disable", `ALTER TABLE employees ENABLE CONSTRAINT ck_sal ENABLE ALL TRIGGERS`},

		// --- MODIFY OPAQUE TYPE ---
		{"modify_opaque_type", `ALTER TABLE xml_tab MODIFY OPAQUE TYPE xmlcol STORE AS OBJECT RELATIONAL (xmldata)`},

		// --- Schema-qualified ---
		{"schema_qualified", `ALTER TABLE hr.employees ADD (middle_name VARCHAR2(50))`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("expected at least 1 statement, got 0")
			}
			_ = ast.NodeToString(result.Items[0])
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 92: index_indextype_operator
// ---------------------------------------------------------------------------

func TestParseIndexIndextypeOperator(t *testing.T) {
	// ===== CREATE INDEX =====
	t.Run("create_index_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_emp ON employees (last_name)")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateIndexStmt)
		if !ok { t.Fatalf("expected *CreateIndexStmt, got %T", raw.Stmt) }
		if stmt.Name.Name != "IDX_EMP" { t.Errorf("expected IDX_EMP, got %q", stmt.Name.Name) }
		if stmt.Table.Name != "EMPLOYEES" { t.Errorf("expected EMPLOYEES, got %q", stmt.Table.Name) }
	})
	t.Run("create_unique_index", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE UNIQUE INDEX idx_emp_id ON hr.employees (employee_id)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Unique { t.Error("expected Unique=true") }
	})
	t.Run("create_bitmap_index", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE BITMAP INDEX idx_status ON orders (status)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Bitmap { t.Error("expected Bitmap=true") }
	})
	t.Run("create_multivalue_index", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MULTIVALUE INDEX idx_tags ON products (tags)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Multivalue { t.Error("expected Multivalue=true") }
	})
	t.Run("create_index_if_not_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX IF NOT EXISTS idx_emp ON employees (last_name)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.IfNotExists { t.Error("expected IfNotExists=true") }
	})
	t.Run("create_index_reverse", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_rev ON t (a) REVERSE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Reverse { t.Error("expected Reverse=true") }
	})
	t.Run("create_index_nosort", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_ns ON t (a) NOSORT")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.NoSort { t.Error("expected NoSort=true") }
	})
	t.Run("create_index_sort", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_s ON t (a) SORT")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Sort { t.Error("expected Sort=true") }
	})
	t.Run("create_index_visible", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_v ON t (a) VISIBLE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Visible { t.Error("expected Visible=true") }
	})
	t.Run("create_index_invisible", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_i ON t (a) INVISIBLE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Invisible { t.Error("expected Invisible=true") }
	})
	t.Run("create_index_logging", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_l ON t (a) LOGGING")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Logging { t.Error("expected Logging=true") }
	})
	t.Run("create_index_nologging", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_nl ON t (a) NOLOGGING")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.NoLogging { t.Error("expected NoLogging=true") }
	})
	t.Run("create_index_compress_advanced", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_ca ON t (a) COMPRESS ADVANCED HIGH")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Compress != "ADVANCED HIGH" { t.Errorf("expected 'ADVANCED HIGH', got %q", stmt.Compress) }
	})
	t.Run("create_index_domain", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_ctx ON docs (content) INDEXTYPE IS CTXSYS.CONTEXT")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.IndexType == nil { t.Fatal("expected non-nil IndexType") }
	})
	t.Run("create_index_parameters", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_ctx ON docs (content) INDEXTYPE IS CTXSYS.CONTEXT PARAMETERS ('WORDLIST my_wordlist')")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Parameters != "WORDLIST my_wordlist" { t.Errorf("expected parameters, got %q", stmt.Parameters) }
	})
	t.Run("create_index_cluster", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_clust ON CLUSTER emp_dept_cluster")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Cluster == nil { t.Fatal("expected non-nil Cluster") }
		if stmt.Cluster.Name != "EMP_DEPT_CLUSTER" { t.Errorf("expected EMP_DEPT_CLUSTER, got %q", stmt.Cluster.Name) }
	})
	t.Run("create_index_pctfree", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_pf ON t (a) PCTFREE 20")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.PctFree != "20" { t.Errorf("expected 20, got %q", stmt.PctFree) }
	})
	t.Run("create_index_local", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_local ON t (a) LOCAL")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Local { t.Error("expected Local=true") }
	})
	t.Run("create_index_global", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_global ON t (a) GLOBAL")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Global { t.Error("expected Global=true") }
	})
	t.Run("create_index_deferred_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_di ON t (a) DEFERRED INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Invalidation != "DEFERRED" { t.Errorf("expected DEFERRED, got %q", stmt.Invalidation) }
	})
	t.Run("create_index_indexing_full", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_if ON t (a) INDEXING FULL")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.IndexingFull { t.Error("expected IndexingFull=true") }
	})
	t.Run("create_index_multi_options", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE UNIQUE INDEX idx_m ON t (a, b DESC) TABLESPACE ts1 ONLINE PARALLEL 4 COMPRESS 2 NOLOGGING VISIBLE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if !stmt.Unique { t.Error("expected Unique=true") }
		if stmt.Tablespace != "TS1" { t.Errorf("expected TS1, got %q", stmt.Tablespace) }
		if !stmt.Online { t.Error("expected Online=true") }
		if stmt.Parallel != "4" { t.Errorf("expected 4, got %q", stmt.Parallel) }
		if stmt.Compress != "2" { t.Errorf("expected 2, got %q", stmt.Compress) }
		if !stmt.NoLogging { t.Error("expected NoLogging=true") }
		if !stmt.Visible { t.Error("expected Visible=true") }
	})
	t.Run("create_index_nocompress", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_nc ON t (a) NOCOMPRESS")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Compress != "NOCOMPRESS" { t.Errorf("expected NOCOMPRESS, got %q", stmt.Compress) }
	})
	t.Run("create_index_noparallel", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_np ON t (a) NOPARALLEL")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Parallel != "NOPARALLEL" { t.Errorf("expected NOPARALLEL, got %q", stmt.Parallel) }
	})
	t.Run("create_index_immediate_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEX idx_ii ON t (a) IMMEDIATE INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndexStmt)
		if stmt.Invalidation != "IMMEDIATE" { t.Errorf("expected IMMEDIATE, got %q", stmt.Invalidation) }
	})

	// ===== ALTER INDEX =====
	t.Run("alter_index_deallocate_unused", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 DEALLOCATE UNUSED")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "DEALLOCATE_UNUSED" { t.Errorf("expected DEALLOCATE_UNUSED, got %q", stmt.Action) }
	})
	t.Run("alter_index_deallocate_unused_keep", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 DEALLOCATE UNUSED KEEP 100M")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.DeallocateKeep != "100M" { t.Errorf("expected 100M, got %q", stmt.DeallocateKeep) }
	})
	t.Run("alter_index_allocate_extent", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 ALLOCATE EXTENT")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "ALLOCATE_EXTENT" { t.Errorf("expected ALLOCATE_EXTENT, got %q", stmt.Action) }
	})
	t.Run("alter_index_parameters", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 PARAMETERS ('my_params')")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Parameters != "my_params" { t.Errorf("expected my_params, got %q", stmt.Parameters) }
	})
	t.Run("alter_index_deferred_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 DEFERRED INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Invalidation != "DEFERRED" { t.Errorf("expected DEFERRED, got %q", stmt.Invalidation) }
	})
	t.Run("alter_index_immediate_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 IMMEDIATE INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Invalidation != "IMMEDIATE" { t.Errorf("expected IMMEDIATE, got %q", stmt.Invalidation) }
	})
	t.Run("alter_index_modify_default_attrs", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MODIFY DEFAULT ATTRIBUTES TABLESPACE ts_new")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "MODIFY_DEFAULT_ATTRIBUTES" { t.Errorf("expected MODIFY_DEFAULT_ATTRIBUTES, got %q", stmt.Action) }
		if stmt.Tablespace != "TS_NEW" { t.Errorf("expected TS_NEW, got %q", stmt.Tablespace) }
	})
	t.Run("alter_index_modify_default_attrs_for_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MODIFY DEFAULT ATTRIBUTES FOR PARTITION p1 LOGGING")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.ModifyDefaultFor != "P1" { t.Errorf("expected P1, got %q", stmt.ModifyDefaultFor) }
		if !stmt.Logging { t.Error("expected Logging=true") }
	})
	t.Run("alter_index_add_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 ADD PARTITION p_new TABLESPACE ts1")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "ADD_PARTITION" { t.Errorf("expected ADD_PARTITION, got %q", stmt.Action) }
		if stmt.AddPartitionName != "P_NEW" { t.Errorf("expected P_NEW, got %q", stmt.AddPartitionName) }
	})
	t.Run("alter_index_modify_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MODIFY PARTITION p1 UNUSABLE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "MODIFY_PARTITION" { t.Errorf("expected MODIFY_PARTITION, got %q", stmt.Action) }
		if stmt.ModifyPartAction != "UNUSABLE" { t.Errorf("expected UNUSABLE, got %q", stmt.ModifyPartAction) }
	})
	t.Run("alter_index_rename_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 RENAME PARTITION old_p TO new_p")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "RENAME_PARTITION" { t.Errorf("expected RENAME_PARTITION, got %q", stmt.Action) }
		if stmt.NewName != "NEW_P" { t.Errorf("expected NEW_P, got %q", stmt.NewName) }
	})
	t.Run("alter_index_rename_subpartition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 RENAME SUBPARTITION old_sp TO new_sp")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "RENAME_SUBPARTITION" { t.Errorf("expected RENAME_SUBPARTITION, got %q", stmt.Action) }
	})
	t.Run("alter_index_drop_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 DROP PARTITION p1")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "DROP_PARTITION" { t.Errorf("expected DROP_PARTITION, got %q", stmt.Action) }
		if stmt.Partition != "P1" { t.Errorf("expected P1, got %q", stmt.Partition) }
	})
	t.Run("alter_index_split_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 SPLIT PARTITION p_old AT (100)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "SPLIT_PARTITION" { t.Errorf("expected SPLIT_PARTITION, got %q", stmt.Action) }
		if stmt.SplitPartition != "P_OLD" { t.Errorf("expected P_OLD, got %q", stmt.SplitPartition) }
	})
	t.Run("alter_index_coalesce_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 COALESCE PARTITION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "COALESCE_PARTITION" { t.Errorf("expected COALESCE_PARTITION, got %q", stmt.Action) }
	})
	t.Run("alter_index_modify_subpartition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MODIFY SUBPARTITION sp1 UNUSABLE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Action != "MODIFY_SUBPARTITION" { t.Errorf("expected MODIFY_SUBPARTITION, got %q", stmt.Action) }
	})
	t.Run("alter_index_pctfree", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 PCTFREE 30")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.PctFree != "30" { t.Errorf("expected 30, got %q", stmt.PctFree) }
	})
	t.Run("alter_index_rebuild_compress_advanced", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 REBUILD COMPRESS ADVANCED HIGH")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Compress != "ADVANCED HIGH" { t.Errorf("expected 'ADVANCED HIGH', got %q", stmt.Compress) }
	})
	t.Run("alter_index_rebuild_parameters", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 REBUILD PARAMETERS ('rebuild_params')")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Parameters != "rebuild_params" { t.Errorf("expected rebuild_params, got %q", stmt.Parameters) }
	})
	t.Run("alter_index_rebuild_deferred_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 REBUILD DEFERRED INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Invalidation != "DEFERRED" { t.Errorf("expected DEFERRED, got %q", stmt.Invalidation) }
	})
	t.Run("alter_index_modify_partition_parameters", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MODIFY PARTITION p1 PARAMETERS ('some_params')")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Parameters != "some_params" { t.Errorf("expected some_params, got %q", stmt.Parameters) }
	})
	t.Run("alter_index_modify_partition_coalesce", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 MODIFY PARTITION p1 COALESCE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.ModifyPartAction != "COALESCE" { t.Errorf("expected COALESCE, got %q", stmt.ModifyPartAction) }
	})
	t.Run("alter_index_rebuild_subpartition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEX idx1 REBUILD SUBPARTITION sp1 ONLINE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndexStmt)
		if stmt.Subpartition != "SP1" { t.Errorf("expected SP1, got %q", stmt.Subpartition) }
		if !stmt.Online { t.Error("expected Online=true") }
	})

	// ===== DROP INDEX =====
	t.Run("drop_index_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEX hr.idx_emp")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if stmt.ObjectType != ast.OBJECT_INDEX { t.Errorf("expected OBJECT_INDEX, got %d", stmt.ObjectType) }
	})
	t.Run("drop_index_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEX IF EXISTS idx_emp")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.IfExists { t.Error("expected IfExists=true") }
	})
	t.Run("drop_index_online", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEX idx_emp ONLINE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.Online { t.Error("expected Online=true") }
	})
	t.Run("drop_index_force", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEX idx_emp FORCE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.Force { t.Error("expected Force=true") }
	})
	t.Run("drop_index_deferred_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEX idx_emp ONLINE FORCE DEFERRED INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.Online { t.Error("expected Online=true") }
		if !stmt.Force { t.Error("expected Force=true") }
		if stmt.Invalidation != "DEFERRED" { t.Errorf("expected DEFERRED, got %q", stmt.Invalidation) }
	})
	t.Run("drop_index_immediate_invalidation", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEX idx_emp IMMEDIATE INVALIDATION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if stmt.Invalidation != "IMMEDIATE" { t.Errorf("expected IMMEDIATE, got %q", stmt.Invalidation) }
	})

	// ===== CREATE INDEXTYPE =====
	t.Run("create_indextype_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_itype FOR my_op(NUMBER) USING my_impl_type")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if stmt.Name.Name != "MY_ITYPE" { t.Errorf("expected MY_ITYPE, got %q", stmt.Name.Name) }
		if len(stmt.Operators) != 1 { t.Fatalf("expected 1 operator, got %d", len(stmt.Operators)) }
		if stmt.UsingType == nil { t.Fatal("expected non-nil UsingType") }
	})
	t.Run("create_or_replace_indextype", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OR REPLACE INDEXTYPE my_itype FOR my_op(NUMBER, VARCHAR2) USING my_type")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if !stmt.OrReplace { t.Error("expected OrReplace=true") }
		if len(stmt.Operators[0].ParamTypes) != 2 { t.Errorf("expected 2 param types, got %d", len(stmt.Operators[0].ParamTypes)) }
	})
	t.Run("create_indextype_multi_operators", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_itype FOR op1(NUMBER), op2(VARCHAR2) USING my_type")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if len(stmt.Operators) != 2 { t.Fatalf("expected 2 operators, got %d", len(stmt.Operators)) }
	})
	t.Run("create_indextype_local_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_itype FOR my_op(NUMBER) USING my_type WITH LOCAL PARTITION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if !stmt.WithLocal { t.Error("expected WithLocal=true") }
	})
	t.Run("create_indextype_local_range_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_itype FOR my_op(NUMBER) USING my_type WITH LOCAL RANGE PARTITION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if !stmt.WithRange { t.Error("expected WithRange=true") }
	})
	t.Run("create_indextype_storage_table", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_itype FOR my_op(NUMBER) USING my_type WITH SYSTEM MANAGED STORAGE TABLES")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if stmt.StorageTable != "SYSTEM" { t.Errorf("expected SYSTEM, got %q", stmt.StorageTable) }
	})
	t.Run("create_indextype_if_not_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE IF NOT EXISTS my_itype FOR my_op(NUMBER) USING my_type")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if !stmt.IfNotExists { t.Error("expected IfNotExists=true") }
	})
	t.Run("create_indextype_sharing", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE INDEXTYPE my_itype SHARING = METADATA FOR my_op(NUMBER) USING my_type")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateIndextypeStmt)
		if stmt.Sharing != "METADATA" { t.Errorf("expected METADATA, got %q", stmt.Sharing) }
	})

	// ===== ALTER INDEXTYPE =====
	t.Run("alter_indextype_add", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_itype ADD my_op2(VARCHAR2)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if stmt.Action != "ADD_DROP" { t.Errorf("expected ADD_DROP, got %q", stmt.Action) }
		if !stmt.Modifications[0].Add { t.Error("expected Add=true") }
	})
	t.Run("alter_indextype_drop", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_itype DROP my_op2(VARCHAR2)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if stmt.Modifications[0].Add { t.Error("expected Add=false") }
	})
	t.Run("alter_indextype_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_itype COMPILE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if stmt.Action != "COMPILE" { t.Errorf("expected COMPILE, got %q", stmt.Action) }
	})
	t.Run("alter_indextype_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE IF EXISTS my_itype COMPILE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if !stmt.IfExists { t.Error("expected IfExists=true") }
	})
	t.Run("alter_indextype_using_type", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_itype ADD my_op2(NUMBER) USING my_new_type")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if stmt.UsingType == nil { t.Fatal("expected non-nil UsingType") }
	})
	t.Run("alter_indextype_local_partition", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_itype ADD my_op2(NUMBER) WITH LOCAL PARTITION")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if !stmt.WithLocal { t.Error("expected WithLocal=true") }
	})
	t.Run("alter_indextype_storage_table", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER INDEXTYPE my_itype ADD my_op2(NUMBER) WITH SYSTEM MANAGED STORAGE TABLES")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterIndextypeStmt)
		if stmt.StorageTable != "SYSTEM" { t.Errorf("expected SYSTEM, got %q", stmt.StorageTable) }
	})

	// ===== DROP INDEXTYPE =====
	t.Run("drop_indextype_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEXTYPE my_itype")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if stmt.ObjectType != ast.OBJECT_INDEXTYPE { t.Errorf("expected OBJECT_INDEXTYPE, got %d", stmt.ObjectType) }
	})
	t.Run("drop_indextype_force", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEXTYPE my_itype FORCE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.Force { t.Error("expected Force=true") }
	})
	t.Run("drop_indextype_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP INDEXTYPE IF EXISTS my_itype")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.IfExists { t.Error("expected IfExists=true") }
	})

	// ===== CREATE OPERATOR =====
	t.Run("create_operator_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING my_eq_func")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateOperatorStmt)
		if stmt.Name.Name != "MY_EQ" { t.Errorf("expected MY_EQ, got %q", stmt.Name.Name) }
		if len(stmt.Bindings) != 1 { t.Fatalf("expected 1 binding, got %d", len(stmt.Bindings)) }
		if stmt.Bindings[0].ReturnType != "NUMBER" { t.Errorf("expected NUMBER, got %q", stmt.Bindings[0].ReturnType) }
	})
	t.Run("create_or_replace_operator", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OR REPLACE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING my_eq_func")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateOperatorStmt)
		if !stmt.OrReplace { t.Error("expected OrReplace=true") }
	})
	t.Run("create_operator_schema_func", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING hr.my_eq_func")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateOperatorStmt)
		if stmt.Bindings[0].UsingFunc == nil { t.Fatal("expected non-nil UsingFunc") }
	})
	t.Run("create_operator_ancillary", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OPERATOR my_anc BINDING (NUMBER) RETURN NUMBER ANCILLARY TO my_primary(NUMBER) USING my_anc_func")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateOperatorStmt)
		if stmt.Bindings[0].AncillaryTo == nil { t.Fatal("expected non-nil AncillaryTo") }
	})
	t.Run("create_operator_context", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OPERATOR my_op BINDING (NUMBER) RETURN NUMBER WITH INDEX CONTEXT, SCAN CONTEXT my_impl_type COMPUTE ANCILLARY DATA USING my_func")
		b := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateOperatorStmt).Bindings[0]
		if !b.WithIndexCtx { t.Error("expected WithIndexCtx=true") }
		if b.ScanCtxType != "MY_IMPL_TYPE" { t.Errorf("expected MY_IMPL_TYPE, got %q", b.ScanCtxType) }
		if !b.ComputeAnc { t.Error("expected ComputeAnc=true") }
	})
	t.Run("create_operator_sharing", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING my_eq_func SHARING = METADATA")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateOperatorStmt)
		if stmt.Sharing != "METADATA" { t.Errorf("expected METADATA, got %q", stmt.Sharing) }
	})

	// ===== ALTER OPERATOR =====
	t.Run("alter_operator_add_binding", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER OPERATOR my_eq ADD BINDING (VARCHAR2, VARCHAR2) RETURN NUMBER USING my_eq_str")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterOperatorStmt)
		if stmt.Action != "ADD_BINDING" { t.Errorf("expected ADD_BINDING, got %q", stmt.Action) }
		if stmt.Binding == nil { t.Fatal("expected non-nil Binding") }
	})
	t.Run("alter_operator_drop_binding", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER OPERATOR my_eq DROP BINDING (VARCHAR2, VARCHAR2)")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterOperatorStmt)
		if stmt.Action != "DROP_BINDING" { t.Errorf("expected DROP_BINDING, got %q", stmt.Action) }
		if len(stmt.DropTypes) != 2 { t.Errorf("expected 2 types, got %d", len(stmt.DropTypes)) }
	})
	t.Run("alter_operator_drop_binding_force", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER OPERATOR my_eq DROP BINDING (NUMBER) FORCE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterOperatorStmt)
		if !stmt.DropForce { t.Error("expected DropForce=true") }
	})
	t.Run("alter_operator_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER OPERATOR my_eq COMPILE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterOperatorStmt)
		if stmt.Action != "COMPILE" { t.Errorf("expected COMPILE, got %q", stmt.Action) }
	})
	t.Run("alter_operator_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER OPERATOR IF EXISTS my_eq COMPILE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterOperatorStmt)
		if !stmt.IfExists { t.Error("expected IfExists=true") }
	})

	// ===== DROP OPERATOR =====
	t.Run("drop_operator_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP OPERATOR my_eq")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if stmt.ObjectType != ast.OBJECT_OPERATOR { t.Errorf("expected OBJECT_OPERATOR, got %d", stmt.ObjectType) }
	})
	t.Run("drop_operator_force", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP OPERATOR my_eq FORCE")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.Force { t.Error("expected Force=true") }
	})
	t.Run("drop_operator_if_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "DROP OPERATOR IF EXISTS my_eq")
		stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DropStmt)
		if !stmt.IfExists { t.Error("expected IfExists=true") }
	})

	// ===== Serialization =====
	t.Run("serialize_all", func(t *testing.T) {
		sqls := []string{
			"CREATE UNIQUE BITMAP INDEX IF NOT EXISTS hr.idx1 ON hr.emp (last_name DESC, first_name) TABLESPACE ts1 ONLINE PARALLEL 4 COMPRESS 2 NOLOGGING VISIBLE REVERSE",
			"ALTER INDEX IF EXISTS hr.idx1 REBUILD PARTITION p1 TABLESPACE ts1 ONLINE REVERSE PARALLEL 4 COMPRESS 2 LOGGING",
			"CREATE OR REPLACE INDEXTYPE my_itype FOR my_op(NUMBER, VARCHAR2) USING my_type WITH LOCAL RANGE PARTITION",
			"CREATE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING hr.my_eq_func",
			"ALTER OPERATOR my_eq ADD BINDING (VARCHAR2, VARCHAR2) RETURN NUMBER USING my_eq_str",
			"ALTER INDEXTYPE my_itype ADD my_op2(NUMBER)",
		}
		for _, sql := range sqls {
			result := ParseAndCheck(t, sql)
			s := ast.NodeToString(result.Items[0])
			if s == "" { t.Errorf("expected non-empty serialization for %q", sql) }
		}
	})
}

// TestBatch93_ViewsAllTypes tests CREATE/ALTER VIEW, MATERIALIZED VIEW, ANALYTIC VIEW, and JSON DUALITY VIEW.
func TestBatch93_ViewsAllTypes(t *testing.T) {
	// ===== CREATE VIEW =====
	t.Run("create_view_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateViewStmt)
		if !ok {
			t.Fatalf("expected *CreateViewStmt, got %T", raw.Stmt)
		}
		if stmt.Name == nil || stmt.Name.Name != "EMP_V" {
			t.Errorf("expected name=EMP_V, got %v", stmt.Name)
		}
		if stmt.Query == nil {
			t.Error("expected query to be set")
		}
	})

	t.Run("create_view_with_columns", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v (eid, ename) AS SELECT employee_id, last_name FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.Columns == nil || len(stmt.Columns.Items) != 2 {
			t.Errorf("expected 2 columns, got %v", stmt.Columns)
		}
	})

	t.Run("create_view_or_replace", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OR REPLACE VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.OrReplace {
			t.Error("expected OrReplace=true")
		}
	})

	t.Run("create_view_if_not_exists", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW IF NOT EXISTS emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.IfNotExists {
			t.Error("expected IfNotExists=true")
		}
	})

	t.Run("create_view_force", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE FORCE VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.Force {
			t.Error("expected Force=true")
		}
	})

	t.Run("create_view_noforce", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE NO FORCE VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.NoForce {
			t.Error("expected NoForce=true")
		}
	})

	t.Run("create_view_editioning", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE EDITIONING VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.Editioning {
			t.Error("expected Editioning=true")
		}
	})

	t.Run("create_view_editionable", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE EDITIONABLE VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.Editionable != "EDITIONABLE" {
			t.Errorf("expected Editionable=EDITIONABLE, got %q", stmt.Editionable)
		}
	})

	t.Run("create_view_noneditionable", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE NONEDITIONABLE VIEW emp_v AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.Editionable != "NONEDITIONABLE" {
			t.Errorf("expected Editionable=NONEDITIONABLE, got %q", stmt.Editionable)
		}
	})

	t.Run("create_view_sharing", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v SHARING = METADATA AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.Sharing != "METADATA" {
			t.Errorf("expected Sharing=METADATA, got %q", stmt.Sharing)
		}
	})

	t.Run("create_view_bequeath", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v BEQUEATH DEFINER AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.Bequeath != "DEFINER" {
			t.Errorf("expected Bequeath=DEFINER, got %q", stmt.Bequeath)
		}
	})

	t.Run("create_view_with_check_option_constraint", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v AS SELECT * FROM employees WITH CHECK OPTION CONSTRAINT emp_ck")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.WithCheckOpt {
			t.Error("expected WithCheckOpt=true")
		}
		if stmt.ConstraintName != "EMP_CK" {
			t.Errorf("expected ConstraintName=EMP_CK, got %q", stmt.ConstraintName)
		}
	})

	t.Run("create_view_with_read_only_constraint", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v AS SELECT * FROM employees WITH READ ONLY CONSTRAINT emp_ro")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.WithReadOnly {
			t.Error("expected WithReadOnly=true")
		}
		if stmt.ConstraintName != "EMP_RO" {
			t.Errorf("expected ConstraintName=EMP_RO, got %q", stmt.ConstraintName)
		}
	})

	t.Run("create_view_container_map", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE VIEW emp_v CONTAINER_MAP AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.ContainerMap {
			t.Error("expected ContainerMap=true")
		}
	})

	// ===== ALTER VIEW =====
	t.Run("alter_view_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v COMPILE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterViewStmt)
		if !ok {
			t.Fatalf("expected *AlterViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "COMPILE" {
			t.Errorf("expected Action=COMPILE, got %q", stmt.Action)
		}
	})

	t.Run("alter_view_add_constraint", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v ADD CONSTRAINT emp_pk PRIMARY KEY (employee_id) DISABLE NOVALIDATE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "ADD_CONSTRAINT" {
			t.Errorf("expected ADD_CONSTRAINT, got %q", stmt.Action)
		}
	})

	t.Run("alter_view_modify_constraint_rely", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v MODIFY CONSTRAINT emp_pk RELY")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "MODIFY_CONSTRAINT" {
			t.Errorf("expected MODIFY_CONSTRAINT, got %q", stmt.Action)
		}
		if !stmt.Rely {
			t.Error("expected Rely=true")
		}
	})

	t.Run("alter_view_drop_constraint", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v DROP CONSTRAINT emp_pk")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "DROP_CONSTRAINT" {
			t.Errorf("expected DROP_CONSTRAINT, got %q", stmt.Action)
		}
	})

	t.Run("alter_view_read_only", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v READ ONLY")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "READ_ONLY" {
			t.Errorf("expected READ_ONLY, got %q", stmt.Action)
		}
	})

	t.Run("alter_view_read_write", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v READ WRITE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "READ_WRITE" {
			t.Errorf("expected READ_WRITE, got %q", stmt.Action)
		}
	})

	t.Run("alter_view_editionable", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v EDITIONABLE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "EDITIONABLE" {
			t.Errorf("expected EDITIONABLE, got %q", stmt.Action)
		}
	})

	t.Run("alter_view_annotations", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER VIEW emp_v ANNOTATIONS (ADD note 'some note')")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterViewStmt)
		if stmt.Action != "ANNOTATIONS" {
			t.Errorf("expected ANNOTATIONS, got %q", stmt.Action)
		}
	})

	// ===== CREATE MATERIALIZED VIEW =====
	t.Run("create_mview_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateViewStmt)
		if !ok {
			t.Fatalf("expected *CreateViewStmt, got %T", raw.Stmt)
		}
		if !stmt.Materialized {
			t.Error("expected Materialized=true")
		}
	})

	t.Run("create_mview_with_refresh", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv REFRESH FAST ON COMMIT AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.RefreshMethod != "FAST" {
			t.Errorf("expected RefreshMethod=FAST, got %q", stmt.RefreshMethod)
		}
		if stmt.RefreshMode != "ON COMMIT" {
			t.Errorf("expected RefreshMode=ON COMMIT, got %q", stmt.RefreshMode)
		}
	})

	t.Run("create_mview_with_build", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv BUILD DEFERRED REFRESH COMPLETE ON DEMAND AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.BuildMode != "DEFERRED" {
			t.Errorf("expected BuildMode=DEFERRED, got %q", stmt.BuildMode)
		}
	})

	t.Run("create_mview_on_prebuilt", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv ON PREBUILT TABLE REFRESH FAST AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.OnPrebuilt {
			t.Error("expected OnPrebuilt=true")
		}
	})

	t.Run("create_mview_never_refresh", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv NEVER REFRESH AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.NeverRefresh {
			t.Error("expected NeverRefresh=true")
		}
	})

	t.Run("create_mview_cache_parallel", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv CACHE PARALLEL 4 AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if stmt.CacheMode != "CACHE" {
			t.Errorf("expected CacheMode=CACHE, got %q", stmt.CacheMode)
		}
		if stmt.ParallelMode != "PARALLEL" {
			t.Errorf("expected ParallelMode=PARALLEL, got %q", stmt.ParallelMode)
		}
	})

	t.Run("create_mview_enable_query_rewrite", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv ENABLE QUERY REWRITE AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.EnableQuery {
			t.Error("expected EnableQuery=true")
		}
	})

	t.Run("create_mview_with_pk", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW emp_mv REFRESH FAST WITH PRIMARY KEY AS SELECT * FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.CreateViewStmt)
		if !stmt.WithPK {
			t.Error("expected WithPK=true")
		}
	})

	// ===== ALTER MATERIALIZED VIEW =====
	t.Run("alter_mview_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW emp_mv COMPILE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterMaterializedViewStmt)
		if !ok {
			t.Fatalf("expected *AlterMaterializedViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "COMPILE" {
			t.Errorf("expected COMPILE, got %q", stmt.Action)
		}
	})

	t.Run("alter_mview_refresh_fast", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW emp_mv REFRESH FAST ON COMMIT")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterMaterializedViewStmt)
		if stmt.Action != "REFRESH" {
			t.Errorf("expected REFRESH, got %q", stmt.Action)
		}
		if stmt.RefreshMethod != "FAST" {
			t.Errorf("expected RefreshMethod=FAST, got %q", stmt.RefreshMethod)
		}
	})

	t.Run("alter_mview_cache", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW emp_mv CACHE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterMaterializedViewStmt)
		if stmt.Action != "CACHE" {
			t.Errorf("expected CACHE, got %q", stmt.Action)
		}
	})

	t.Run("alter_mview_parallel", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW emp_mv PARALLEL 4")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterMaterializedViewStmt)
		if stmt.Action != "PARALLEL" {
			t.Errorf("expected PARALLEL, got %q", stmt.Action)
		}
	})

	t.Run("alter_mview_shrink_space", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW emp_mv SHRINK SPACE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterMaterializedViewStmt)
		if stmt.Action != "SHRINK" {
			t.Errorf("expected SHRINK, got %q", stmt.Action)
		}
	})

	t.Run("alter_mview_consider_fresh", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW emp_mv CONSIDER FRESH")
		raw := result.Items[0].(*ast.RawStmt)
		stmt := raw.Stmt.(*ast.AlterMaterializedViewStmt)
		if stmt.Action != "CONSIDER_FRESH" {
			t.Errorf("expected CONSIDER_FRESH, got %q", stmt.Action)
		}
	})

	// ===== CREATE MATERIALIZED VIEW LOG =====
	t.Run("create_mview_log_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW LOG ON employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateMviewLogStmt)
		if !ok {
			t.Fatalf("expected *CreateMviewLogStmt, got %T", raw.Stmt)
		}
		if stmt.OnTable == nil || stmt.OnTable.Name != "EMPLOYEES" {
			t.Errorf("expected OnTable=EMPLOYEES, got %v", stmt.OnTable)
		}
	})

	t.Run("create_mview_log_with_options", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW LOG ON employees TABLESPACE ts1 NOLOGGING")
		raw := result.Items[0].(*ast.RawStmt)
		_, ok := raw.Stmt.(*ast.CreateMviewLogStmt)
		if !ok {
			t.Fatalf("expected *CreateMviewLogStmt, got %T", raw.Stmt)
		}
	})

	t.Run("create_mview_log_with_clauses", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW LOG ON employees WITH PRIMARY KEY, ROWID (salary) INCLUDING NEW VALUES")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateMviewLogStmt)
		if !ok {
			t.Fatalf("expected *CreateMviewLogStmt, got %T", raw.Stmt)
		}
		if !stmt.WithPK {
			t.Error("expected WithPK=true")
		}
		if !stmt.WithRowID {
			t.Error("expected WithRowID=true")
		}
		if !stmt.Including {
			t.Error("expected Including=true")
		}
	})

	t.Run("create_mview_log_with_purge", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE MATERIALIZED VIEW LOG ON employees WITH PRIMARY KEY PURGE IMMEDIATE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateMviewLogStmt)
		if !ok {
			t.Fatalf("expected *CreateMviewLogStmt, got %T", raw.Stmt)
		}
		if stmt.PurgeMode != "IMMEDIATE_SYNC" {
			t.Errorf("expected PurgeMode=IMMEDIATE_SYNC, got %q", stmt.PurgeMode)
		}
	})

	// ===== ALTER MATERIALIZED VIEW LOG =====
	t.Run("alter_mview_log_add_columns", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW LOG ON employees ADD (salary, commission_pct)")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterMviewLogStmt)
		if !ok {
			t.Fatalf("expected *AlterMviewLogStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "ADD" {
			t.Errorf("expected ADD, got %q", stmt.Action)
		}
	})

	t.Run("alter_mview_log_shrink_space", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW LOG ON employees SHRINK SPACE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterMviewLogStmt)
		if !ok {
			t.Fatalf("expected *AlterMviewLogStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "SHRINK" {
			t.Errorf("expected SHRINK, got %q", stmt.Action)
		}
	})

	t.Run("alter_mview_log_purge", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER MATERIALIZED VIEW LOG ON employees PURGE IMMEDIATE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterMviewLogStmt)
		if !ok {
			t.Fatalf("expected *AlterMviewLogStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "PURGE" {
			t.Errorf("expected PURGE, got %q", stmt.Action)
		}
	})

	// ===== CREATE ANALYTIC VIEW =====
	t.Run("create_analytic_view_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE ANALYTIC VIEW sales_av USING sales_fact DIMENSION BY (time_attr HIERARCHY time_hier) MEASURES (sales)")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *CreateAnalyticViewStmt, got %T", raw.Stmt)
		}
		if stmt.Name == nil || stmt.Name.Name != "SALES_AV" {
			t.Errorf("expected name=SALES_AV, got %v", stmt.Name)
		}
	})

	t.Run("create_analytic_view_with_options", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE OR REPLACE ANALYTIC VIEW sales_av USING sales_fact DIMENSION BY (time_attr HIERARCHY time_hier) MEASURES (sales FACT sales)")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *CreateAnalyticViewStmt, got %T", raw.Stmt)
		}
		if !stmt.OrReplace {
			t.Error("expected OrReplace=true")
		}
	})

	// ===== ALTER ANALYTIC VIEW =====
	t.Run("alter_analytic_view_rename", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER ANALYTIC VIEW sales_av RENAME TO sales_av2")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *AlterAnalyticViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "RENAME" {
			t.Errorf("expected RENAME, got %q", stmt.Action)
		}
		if stmt.NewName == nil || stmt.NewName.Name != "SALES_AV2" {
			t.Errorf("expected NewName=SALES_AV2, got %v", stmt.NewName)
		}
	})

	t.Run("alter_analytic_view_compile", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER ANALYTIC VIEW sales_av COMPILE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *AlterAnalyticViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "COMPILE" {
			t.Errorf("expected COMPILE, got %q", stmt.Action)
		}
	})

	t.Run("alter_analytic_view_add_cache", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER ANALYTIC VIEW sales_av ADD CACHE STORE AS COLUMNSTORE")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterAnalyticViewStmt)
		if !ok {
			t.Fatalf("expected *AlterAnalyticViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "ADD_CACHE" {
			t.Errorf("expected ADD_CACHE, got %q", stmt.Action)
		}
	})

	// ===== CREATE JSON RELATIONAL DUALITY VIEW =====
	t.Run("create_json_duality_view_basic", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE JSON RELATIONAL DUALITY VIEW emp_dv AS SELECT employee_id FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateJsonDualityViewStmt)
		if !ok {
			t.Fatalf("expected *CreateJsonDualityViewStmt, got %T", raw.Stmt)
		}
		if stmt.Name == nil || stmt.Name.Name != "EMP_DV" {
			t.Errorf("expected name=EMP_DV, got %v", stmt.Name)
		}
	})

	t.Run("create_json_duality_view_with_replication", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE JSON RELATIONAL DUALITY VIEW emp_dv ENABLE LOGICAL REPLICATION AS SELECT employee_id FROM employees")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.CreateJsonDualityViewStmt)
		if !ok {
			t.Fatalf("expected *CreateJsonDualityViewStmt, got %T", raw.Stmt)
		}
		if !stmt.EnableLogicalReplication {
			t.Error("expected EnableLogicalReplication=true")
		}
	})

	// ===== ALTER JSON RELATIONAL DUALITY VIEW =====
	t.Run("alter_json_duality_view", func(t *testing.T) {
		result := ParseAndCheck(t, "ALTER JSON RELATIONAL DUALITY VIEW emp_dv ENABLE LOGICAL REPLICATION")
		raw := result.Items[0].(*ast.RawStmt)
		stmt, ok := raw.Stmt.(*ast.AlterJsonDualityViewStmt)
		if !ok {
			t.Fatalf("expected *AlterJsonDualityViewStmt, got %T", raw.Stmt)
		}
		if stmt.Action != "ENABLE_LOGICAL_REPLICATION" {
			t.Errorf("expected ENABLE_LOGICAL_REPLICATION, got %q", stmt.Action)
		}
	})

	// ===== DROP =====
	t.Run("drop_mview_log", func(t *testing.T) {
		ParseAndCheck(t, "DROP MATERIALIZED VIEW LOG ON employees")
	})

	t.Run("drop_analytic_view", func(t *testing.T) {
		ParseAndCheck(t, "DROP ANALYTIC VIEW sales_av")
	})

	// ===== Serialization =====
	t.Run("serialize_views", func(t *testing.T) {
		sqls := []string{
			"CREATE VIEW emp_v AS SELECT * FROM employees",
			"CREATE OR REPLACE VIEW emp_v AS SELECT * FROM employees WITH CHECK OPTION",
			"CREATE FORCE VIEW emp_v AS SELECT e.employee_id eid, e.last_name ename FROM employees e",
			"ALTER VIEW emp_v COMPILE",
			"ALTER VIEW emp_v READ ONLY",
			"ALTER VIEW emp_v ADD CONSTRAINT emp_pk PRIMARY KEY (employee_id) DISABLE NOVALIDATE",
			"CREATE MATERIALIZED VIEW emp_mv REFRESH FAST ON COMMIT ENABLE QUERY REWRITE AS SELECT * FROM employees",
			"CREATE MATERIALIZED VIEW emp_mv BUILD DEFERRED AS SELECT * FROM employees",
			"ALTER MATERIALIZED VIEW emp_mv REFRESH COMPLETE ON DEMAND",
			"ALTER MATERIALIZED VIEW emp_mv CONSIDER FRESH",
			"ALTER MATERIALIZED VIEW emp_mv SHRINK SPACE COMPACT",
			"CREATE MATERIALIZED VIEW LOG ON employees WITH PRIMARY KEY INCLUDING NEW VALUES",
			"ALTER MATERIALIZED VIEW LOG ON employees ADD (salary)",
			"CREATE ANALYTIC VIEW sales_av USING sales_fact DIMENSION BY (time_attr HIERARCHY time_hier) MEASURES (sales)",
			"ALTER ANALYTIC VIEW sales_av COMPILE",
			"ALTER ANALYTIC VIEW sales_av RENAME TO sales_av_new",
			"CREATE JSON RELATIONAL DUALITY VIEW emp_dv AS SELECT employee_id FROM employees",
			"ALTER JSON RELATIONAL DUALITY VIEW emp_dv DISABLE LOGICAL REPLICATION",
		}
		for _, sql := range sqls {
			result := ParseAndCheck(t, sql)
			s := ast.NodeToString(result.Items[0])
			if s == "" {
				t.Errorf("expected non-empty serialization for %q", sql)
			}
		}
	})
}
