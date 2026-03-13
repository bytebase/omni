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
		"ALTER PROCEDURE my_proc COMPILE",
		"ALTER FUNCTION my_func COMPILE",
		"ALTER TRIGGER my_trigger ENABLE",
		"ALTER TRIGGER my_trigger DISABLE",
		"ALTER TYPE my_type COMPILE",
		"ALTER PACKAGE my_pkg COMPILE",
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
		"CREATE SCHEMA AUTHORIZATION hr CREATE TABLE t1 (id NUMBER); CREATE VIEW v1 AS SELECT 1 FROM dual;",
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
		name string
		sql  string
	}{
		{"basic", "CREATE DATABASE mydb"},
		{"with_user_sys_password", "CREATE DATABASE mydb USER SYS IDENTIFIED BY password"},
		{"with_logfile", "CREATE DATABASE mydb LOGFILE GROUP 1 '/u01/log1.log' SIZE 100M"},
		{"with_maxlogfiles", "CREATE DATABASE mydb MAXLOGFILES 16 MAXLOGMEMBERS 3"},
		{"with_maxdatafiles", "CREATE DATABASE mydb MAXDATAFILES 100 MAXINSTANCES 8"},
		{"with_archivelog", "CREATE DATABASE mydb ARCHIVELOG"},
		{"with_noarchivelog", "CREATE DATABASE mydb NOARCHIVELOG"},
		{"with_character_set", "CREATE DATABASE mydb CHARACTER SET AL32UTF8"},
		{"with_national_character_set", "CREATE DATABASE mydb NATIONAL CHARACTER SET AL16UTF16"},
		{"with_datafile", "CREATE DATABASE mydb DATAFILE '/u01/data01.dbf' SIZE 500M AUTOEXTEND ON"},
		{"with_default_tablespace", "CREATE DATABASE mydb DEFAULT TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M"},
		{"with_undo_tablespace", "CREATE DATABASE mydb UNDO TABLESPACE undots DATAFILE '/u01/undo01.dbf' SIZE 200M"},
		{"with_default_temp_tablespace", "CREATE DATABASE mydb DEFAULT TEMPORARY TABLESPACE temp TEMPFILE '/u01/temp01.dbf' SIZE 100M"},
		{"with_force_logging", "CREATE DATABASE mydb FORCE LOGGING"},
		{"complex", "CREATE DATABASE proddb USER SYS IDENTIFIED BY oracle LOGFILE GROUP 1 '/u01/redo01.log' SIZE 50M MAXLOGFILES 16 MAXDATAFILES 1024 ARCHIVELOG CHARACTER SET AL32UTF8 DATAFILE '/u01/system01.dbf' SIZE 700M"},
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
	stmt, ok := raw.Stmt.(*ast.AdminDDLStmt)
	if !ok {
		t.Fatalf("expected *AdminDDLStmt, got %T", raw.Stmt)
	}
	if stmt.ObjectType != ast.OBJECT_DATABASE_LINK {
		t.Errorf("expected OBJECT_DATABASE_LINK, got %d", stmt.ObjectType)
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
