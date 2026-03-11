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

func TestNodeToString(t *testing.T) {
	n := &ast.NumberLiteral{Val: "42", Ival: 42, Loc: ast.Loc{Start: 0, End: -1}}
	s := ast.NodeToString(n)
	if s == "" {
		t.Error("NodeToString returned empty string")
	}
}
