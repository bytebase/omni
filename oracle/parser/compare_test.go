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
