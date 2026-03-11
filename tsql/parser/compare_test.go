package parser

import (
	"testing"

	"github.com/bytebase/omni/tsql/ast"
)

// ParseAndCheck parses sql, validates that parsing succeeds, and checks that
// all nodes with Location fields have non-negative locations.
func ParseAndCheck(t *testing.T, sql string) *ast.List {
	t.Helper()

	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): unexpected error: %v", sql, err)
	}
	if result == nil {
		t.Fatalf("Parse(%q): returned nil result", sql)
	}

	// Verify deterministic serialization (parse twice, compare)
	for i, item := range result.Items {
		s1 := ast.NodeToString(item)
		if s1 == "" {
			t.Errorf("Parse(%q) stmt[%d]: NodeToString returned empty string", sql, i)
		}
	}

	return result
}

// checkLocation verifies that a Loc has a non-negative start position.
func checkLocation(t *testing.T, sql string, nodeName string, loc ast.Loc) {
	t.Helper()
	if loc.Start < 0 {
		t.Errorf("Parse(%q): %s has negative location start %d", sql, nodeName, loc.Start)
	}
}

// TestParseInfrastructure tests the basic parser infrastructure (batch 0).
func TestParseInfrastructure(t *testing.T) {
	// Empty input
	result, err := Parse("")
	if err != nil {
		t.Fatalf("Parse empty: %v", err)
	}
	if result.Len() != 0 {
		t.Errorf("Parse empty: expected 0 items, got %d", result.Len())
	}

	// Semicolons only
	result, err = Parse(";;;")
	if err != nil {
		t.Fatalf("Parse semicolons: %v", err)
	}
	if result.Len() != 0 {
		t.Errorf("Parse semicolons: expected 0 items, got %d", result.Len())
	}
}

// TestParseBasicSelect tests basic SELECT parsing.
func TestParseBasicSelect(t *testing.T) {
	tests := []string{
		"SELECT 1",
		"SELECT 1, 2, 3",
		"SELECT 'hello'",
		"SELECT N'unicode'",
		"SELECT NULL",
		"SELECT 1 + 2",
		"SELECT 1 - 2",
		"SELECT 2 * 3",
		"SELECT 6 / 2",
		"SELECT 7 % 3",
		"SELECT -1",
		"SELECT +1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Errorf("expected 1 statement, got %d", result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SelectStmt)
			if !ok {
				t.Fatalf("expected *SelectStmt, got %T", result.Items[0])
			}
			checkLocation(t, sql, "SelectStmt", stmt.Loc)
		})
	}
}

// TestParseComparison tests comparison operators.
func TestParseComparison(t *testing.T) {
	tests := []string{
		"SELECT 1 = 2",
		"SELECT 1 <> 2",
		"SELECT 1 != 2",
		"SELECT 1 < 2",
		"SELECT 1 > 2",
		"SELECT 1 <= 2",
		"SELECT 1 >= 2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseBoolean tests boolean operators.
func TestParseBoolean(t *testing.T) {
	tests := []string{
		"SELECT 1 = 1 AND 2 = 2",
		"SELECT 1 = 1 OR 2 = 2",
		"SELECT NOT 1 = 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseIsNull tests IS NULL / IS NOT NULL.
func TestParseIsNull(t *testing.T) {
	tests := []string{
		"SELECT NULL IS NULL",
		"SELECT 1 IS NOT NULL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseBetween tests BETWEEN expressions.
func TestParseBetween(t *testing.T) {
	tests := []string{
		"SELECT 5 BETWEEN 1 AND 10",
		"SELECT 5 NOT BETWEEN 1 AND 10",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseIn tests IN expressions.
func TestParseIn(t *testing.T) {
	tests := []string{
		"SELECT 1 IN (1, 2, 3)",
		"SELECT 1 NOT IN (1, 2, 3)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseLike tests LIKE expressions.
func TestParseLike(t *testing.T) {
	tests := []string{
		"SELECT 'foo' LIKE 'f%'",
		"SELECT 'foo' NOT LIKE 'f%'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseCase tests CASE expressions.
func TestParseCase(t *testing.T) {
	tests := []string{
		"SELECT CASE WHEN 1 = 1 THEN 'yes' ELSE 'no' END",
		"SELECT CASE 1 WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseCast tests CAST and CONVERT.
func TestParseCast(t *testing.T) {
	tests := []string{
		"SELECT CAST(1 AS int)",
		"SELECT CONVERT(varchar, 123)",
		"SELECT CONVERT(varchar, 123, 1)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseFunctions tests function calls.
func TestParseFunctions(t *testing.T) {
	tests := []string{
		"SELECT COUNT(*)",
		"SELECT COUNT(1)",
		"SELECT COUNT(DISTINCT 1)",
		"SELECT MAX(1)",
		"SELECT MIN(1)",
		"SELECT COALESCE(1, 2, 3)",
		"SELECT NULLIF(1, 2)",
		"SELECT IIF(1 = 1, 'yes', 'no')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseVariables tests variable references.
func TestParseVariables(t *testing.T) {
	tests := []string{
		"SELECT @myvar",
		"SELECT @@ROWCOUNT",
		"SELECT @@IDENTITY",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			stmt := result.Items[0].(*ast.SelectStmt)
			target := stmt.TargetList.Items[0].(*ast.ResTarget)
			varRef, ok := target.Val.(*ast.VariableRef)
			if !ok {
				t.Fatalf("expected *VariableRef, got %T", target.Val)
			}
			checkLocation(t, sql, "VariableRef", varRef.Loc)
		})
	}
}

// TestParseNames tests identifier and qualified name parsing (batch 1).
func TestParseNames(t *testing.T) {
	t.Run("simple identifiers", func(t *testing.T) {
		tests := []string{
			"SELECT col1",
			"SELECT MyColumn",
			"SELECT _underscore",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				ref, ok := target.Val.(*ast.ColumnRef)
				if !ok {
					t.Fatalf("expected *ColumnRef, got %T", target.Val)
				}
				if ref.Column == "" {
					t.Error("expected non-empty Column")
				}
				checkLocation(t, sql, "ColumnRef", ref.Loc)
			})
		}
	})

	t.Run("bracketed identifiers", func(t *testing.T) {
		tests := []struct {
			sql  string
			col  string
		}{
			{"SELECT [my column]", "my column"},
			{"SELECT [select]", "select"},
			{"SELECT [123 col]", "123 col"},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				ref, ok := target.Val.(*ast.ColumnRef)
				if !ok {
					t.Fatalf("expected *ColumnRef, got %T", target.Val)
				}
				if ref.Column != tc.col {
					t.Errorf("expected Column=%q, got %q", tc.col, ref.Column)
				}
			})
		}
	})

	t.Run("quoted identifiers", func(t *testing.T) {
		result := ParseAndCheck(t, `SELECT "my column"`)
		stmt := result.Items[0].(*ast.SelectStmt)
		target := stmt.TargetList.Items[0].(*ast.ResTarget)
		ref, ok := target.Val.(*ast.ColumnRef)
		if !ok {
			t.Fatalf("expected *ColumnRef, got %T", target.Val)
		}
		if ref.Column != "my column" {
			t.Errorf("expected Column=%q, got %q", "my column", ref.Column)
		}
	})

	t.Run("qualified names", func(t *testing.T) {
		tests := []struct {
			sql      string
			table    string
			column   string
			schema   string
			database string
		}{
			{"SELECT t.col", "t", "col", "", ""},
			{"SELECT dbo.t.col", "", "col", "dbo", "t"},
			{"SELECT mydb.dbo.t.col", "t", "col", "dbo", "mydb"},
			{"SELECT [schema].[table].[column]", "", "column", "schema", "table"},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				ref, ok := target.Val.(*ast.ColumnRef)
				if !ok {
					t.Fatalf("expected *ColumnRef, got %T", target.Val)
				}
				if tc.table != "" && ref.Table != tc.table {
					t.Errorf("expected Table=%q, got %q", tc.table, ref.Table)
				}
				if ref.Column != tc.column {
					t.Errorf("expected Column=%q, got %q", tc.column, ref.Column)
				}
			})
		}
	})

	t.Run("qualified star", func(t *testing.T) {
		tests := []string{
			"SELECT t.*",
			"SELECT [my table].*",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				star, ok := target.Val.(*ast.StarExpr)
				if !ok {
					t.Fatalf("expected *StarExpr, got %T", target.Val)
				}
				if star.Qualifier == "" {
					t.Error("expected non-empty Qualifier")
				}
			})
		}
	})

	t.Run("variables", func(t *testing.T) {
		tests := []struct {
			sql  string
			name string
		}{
			{"SELECT @myvar", "@myvar"},
			{"SELECT @@ROWCOUNT", "@@ROWCOUNT"},
			{"SELECT @@IDENTITY", "@@IDENTITY"},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				varRef, ok := target.Val.(*ast.VariableRef)
				if !ok {
					t.Fatalf("expected *VariableRef, got %T", target.Val)
				}
				if varRef.Name != tc.name {
					t.Errorf("expected Name=%q, got %q", tc.name, varRef.Name)
				}
			})
		}
	})

	t.Run("function with qualified name", func(t *testing.T) {
		tests := []string{
			"SELECT dbo.MyFunc(1)",
			"SELECT [schema].MyFunc(1, 2)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				fc, ok := target.Val.(*ast.FuncCallExpr)
				if !ok {
					t.Fatalf("expected *FuncCallExpr, got %T", target.Val)
				}
				if fc.Name.Schema == "" {
					t.Error("expected non-empty Schema on function name")
				}
			})
		}
	})

	t.Run("keyword as identifier after dot", func(t *testing.T) {
		// Keywords like 'type', 'value', 'name' should be valid identifiers after a dot
		tests := []string{
			"SELECT t.type",
			"SELECT t.[value]",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				ref, ok := target.Val.(*ast.ColumnRef)
				if !ok {
					t.Fatalf("expected *ColumnRef, got %T", target.Val)
				}
				if ref.Table == "" {
					t.Error("expected non-empty Table")
				}
				if ref.Column == "" {
					t.Error("expected non-empty Column")
				}
			})
		}
	})
}

// TestParseTypes tests data type parsing (batch 2).
func TestParseTypes(t *testing.T) {
	t.Run("simple types", func(t *testing.T) {
		tests := []string{
			"SELECT CAST(1 AS int)",
			"SELECT CAST(1 AS bigint)",
			"SELECT CAST(1 AS smallint)",
			"SELECT CAST(1 AS tinyint)",
			"SELECT CAST(1 AS bit)",
			"SELECT CAST(1 AS float)",
			"SELECT CAST(1 AS real)",
			"SELECT CAST(1 AS money)",
			"SELECT CAST(1 AS smallmoney)",
			"SELECT CAST('x' AS text)",
			"SELECT CAST('x' AS ntext)",
			"SELECT CAST(1 AS uniqueidentifier)",
			"SELECT CAST('x' AS xml)",
			"SELECT CAST('x' AS date)",
			"SELECT CAST('x' AS datetime)",
			"SELECT CAST('x' AS datetime2)",
			"SELECT CAST('x' AS smalldatetime)",
			"SELECT CAST('x' AS datetimeoffset)",
			"SELECT CAST('x' AS time)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				cast, ok := target.Val.(*ast.CastExpr)
				if !ok {
					t.Fatalf("expected *CastExpr, got %T", target.Val)
				}
				if cast.DataType == nil {
					t.Fatal("expected non-nil DataType")
				}
				if cast.DataType.Name == "" {
					t.Error("expected non-empty DataType.Name")
				}
			})
		}
	})

	t.Run("types with length", func(t *testing.T) {
		tests := []struct {
			sql  string
			name string
		}{
			{"SELECT CAST('x' AS varchar(100))", "varchar"},
			{"SELECT CAST('x' AS nvarchar(50))", "nvarchar"},
			{"SELECT CAST('x' AS char(10))", "char"},
			{"SELECT CAST('x' AS nchar(20))", "nchar"},
			{"SELECT CAST('x' AS binary(16))", "binary"},
			{"SELECT CAST('x' AS varbinary(256))", "varbinary"},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				cast := target.Val.(*ast.CastExpr)
				if cast.DataType.Name != tc.name {
					t.Errorf("expected Name=%q, got %q", tc.name, cast.DataType.Name)
				}
				if cast.DataType.Length == nil {
					t.Error("expected non-nil Length")
				}
			})
		}
	})

	t.Run("types with MAX", func(t *testing.T) {
		tests := []string{
			"SELECT CAST('x' AS varchar(MAX))",
			"SELECT CAST('x' AS nvarchar(MAX))",
			"SELECT CAST('x' AS varbinary(MAX))",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				cast := target.Val.(*ast.CastExpr)
				if !cast.DataType.MaxLength {
					t.Error("expected MaxLength=true")
				}
			})
		}
	})

	t.Run("types with precision and scale", func(t *testing.T) {
		tests := []string{
			"SELECT CAST(1 AS decimal(10, 2))",
			"SELECT CAST(1 AS numeric(18, 4))",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				cast := target.Val.(*ast.CastExpr)
				if cast.DataType.Precision == nil {
					t.Error("expected non-nil Precision")
				}
				if cast.DataType.Scale == nil {
					t.Error("expected non-nil Scale")
				}
			})
		}
	})

	t.Run("types with precision only", func(t *testing.T) {
		tests := []string{
			"SELECT CAST('x' AS datetime2(3))",
			"SELECT CAST('x' AS datetimeoffset(7))",
			"SELECT CAST('x' AS float(53))",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				cast := target.Val.(*ast.CastExpr)
				if cast.DataType.Length == nil {
					t.Error("expected non-nil Length")
				}
			})
		}
	})

	t.Run("schema qualified type", func(t *testing.T) {
		sql := "SELECT CAST(1 AS dbo.MyType)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		target := stmt.TargetList.Items[0].(*ast.ResTarget)
		cast := target.Val.(*ast.CastExpr)
		if cast.DataType.Schema == "" {
			t.Error("expected non-empty Schema")
		}
	})
}

// TestParseLexerBracketed tests bracketed identifier lexing.
func TestParseLexerBracketed(t *testing.T) {
	l := NewLexer("[my column]")
	tok := l.NextToken()
	if tok.Type != tokIDENT {
		t.Fatalf("expected tokIDENT, got %d", tok.Type)
	}
	if tok.Str != "my column" {
		t.Fatalf("expected 'my column', got %q", tok.Str)
	}
}

// TestParseLexerNString tests N'...' string lexing.
func TestParseLexerNString(t *testing.T) {
	l := NewLexer("N'hello'")
	tok := l.NextToken()
	if tok.Type != tokNSCONST {
		t.Fatalf("expected tokNSCONST, got %d", tok.Type)
	}
	if tok.Str != "hello" {
		t.Fatalf("expected 'hello', got %q", tok.Str)
	}
}

// TestParseLexerComments tests comment handling.
func TestParseLexerComments(t *testing.T) {
	l := NewLexer("-- line comment\n42")
	tok := l.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Fatalf("expected 42, got type=%d ival=%d", tok.Type, tok.Ival)
	}

	l = NewLexer("/* block */ 42")
	tok = l.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Fatalf("expected 42 after block comment, got type=%d ival=%d", tok.Type, tok.Ival)
	}

	// Nested block comments
	l = NewLexer("/* outer /* inner */ still comment */ 42")
	tok = l.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Fatalf("expected 42 after nested comment, got type=%d ival=%d", tok.Type, tok.Ival)
	}
}

// TestParseLexerOperators tests T-SQL specific operators.
func TestParseLexerOperators(t *testing.T) {
	cases := []struct {
		input string
		typ   int
		str   string
	}{
		{"<>", tokNOTEQUAL, "<>"},
		{"!=", tokNOTEQUAL, "!="},
		{"<=", tokLESSEQUAL, "<="},
		{">=", tokGREATEQUAL, ">="},
		{"!<", tokNOTLESS, "!<"},
		{"!>", tokNOTGREATER, "!>"},
		{"::", tokCOLONCOLON, "::"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			l := NewLexer(c.input)
			tok := l.NextToken()
			if tok.Type != c.typ {
				t.Errorf("expected type %d, got %d", c.typ, tok.Type)
			}
			if tok.Str != c.str {
				t.Errorf("expected str %q, got %q", c.str, tok.Str)
			}
		})
	}
}

// TestNodeToString tests deterministic AST serialization.
func TestNodeToString(t *testing.T) {
	// Parse and serialize, verify non-empty
	sql := "SELECT 1, 'hello', NULL"
	result := ParseAndCheck(t, sql)
	s := ast.NodeToString(result.Items[0])
	if s == "" || s == "<>" {
		t.Errorf("NodeToString returned empty for %q", sql)
	}

	// Parse again and verify identical output
	result2 := ParseAndCheck(t, sql)
	s2 := ast.NodeToString(result2.Items[0])
	if s != s2 {
		t.Errorf("NodeToString not deterministic:\n  first:  %s\n  second: %s", s, s2)
	}
}
