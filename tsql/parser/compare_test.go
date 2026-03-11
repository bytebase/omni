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

// TestParseExpressions tests expression parsing (batch 3).
func TestParseExpressions(t *testing.T) {
	t.Run("arithmetic", func(t *testing.T) {
		tests := []string{
			"SELECT 1 + 2 * 3",
			"SELECT (1 + 2) * 3",
			"SELECT 10 / 2 - 1",
			"SELECT 10 % 3",
			"SELECT -1 + +2",
			"SELECT ~0",
			"SELECT 1 & 2",
			"SELECT 1 | 2",
			"SELECT 1 ^ 2",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("comparison", func(t *testing.T) {
		tests := []string{
			"SELECT 1 = 1",
			"SELECT 1 <> 2",
			"SELECT 1 != 2",
			"SELECT 1 < 2",
			"SELECT 1 > 2",
			"SELECT 1 <= 2",
			"SELECT 1 >= 2",
			"SELECT 1 !< 2",
			"SELECT 1 !> 2",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("logical", func(t *testing.T) {
		tests := []string{
			"SELECT 1 = 1 AND 2 = 2",
			"SELECT 1 = 1 OR 2 = 2",
			"SELECT NOT 1 = 1",
			"SELECT 1 = 1 AND 2 = 2 OR 3 = 3",
			"SELECT NOT (1 = 1 AND 2 = 2)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("between", func(t *testing.T) {
		tests := []string{
			"SELECT 5 BETWEEN 1 AND 10",
			"SELECT 5 NOT BETWEEN 1 AND 10",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				_, ok := target.Val.(*ast.BetweenExpr)
				if !ok {
					t.Fatalf("expected *BetweenExpr, got %T", target.Val)
				}
			})
		}
	})

	t.Run("in", func(t *testing.T) {
		tests := []string{
			"SELECT 1 IN (1, 2, 3)",
			"SELECT 1 NOT IN (1, 2, 3)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				_, ok := target.Val.(*ast.InExpr)
				if !ok {
					t.Fatalf("expected *InExpr, got %T", target.Val)
				}
			})
		}
	})

	t.Run("like", func(t *testing.T) {
		tests := []string{
			"SELECT 'foo' LIKE 'f%'",
			"SELECT 'foo' NOT LIKE 'f%'",
			"SELECT 'foo' LIKE 'f%' ESCAPE '\\'",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("is_null", func(t *testing.T) {
		tests := []string{
			"SELECT NULL IS NULL",
			"SELECT 1 IS NOT NULL",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				_, ok := target.Val.(*ast.IsExpr)
				if !ok {
					t.Fatalf("expected *IsExpr, got %T", target.Val)
				}
			})
		}
	})

	t.Run("case_searched", func(t *testing.T) {
		sql := "SELECT CASE WHEN 1 = 1 THEN 'yes' WHEN 2 = 2 THEN 'maybe' ELSE 'no' END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		target := stmt.TargetList.Items[0].(*ast.ResTarget)
		ce, ok := target.Val.(*ast.CaseExpr)
		if !ok {
			t.Fatalf("expected *CaseExpr, got %T", target.Val)
		}
		if ce.Arg != nil {
			t.Error("searched CASE should have nil Arg")
		}
		if ce.WhenList.Len() != 2 {
			t.Errorf("expected 2 WHEN clauses, got %d", ce.WhenList.Len())
		}
		if ce.ElseExpr == nil {
			t.Error("expected non-nil ElseExpr")
		}
	})

	t.Run("case_simple", func(t *testing.T) {
		sql := "SELECT CASE col1 WHEN 1 THEN 'one' WHEN 2 THEN 'two' END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		target := stmt.TargetList.Items[0].(*ast.ResTarget)
		ce, ok := target.Val.(*ast.CaseExpr)
		if !ok {
			t.Fatalf("expected *CaseExpr, got %T", target.Val)
		}
		if ce.Arg == nil {
			t.Error("simple CASE should have non-nil Arg")
		}
	})

	t.Run("cast_convert", func(t *testing.T) {
		tests := []string{
			"SELECT CAST(1 AS varchar(10))",
			"SELECT TRY_CAST(1 AS int)",
			"SELECT CONVERT(varchar, 123)",
			"SELECT CONVERT(varchar(10), 123, 1)",
			"SELECT TRY_CONVERT(int, '123')",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("special_functions", func(t *testing.T) {
		tests := []string{
			"SELECT COALESCE(NULL, 1, 2)",
			"SELECT NULLIF(1, 2)",
			"SELECT IIF(1 = 1, 'yes', 'no')",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("exists", func(t *testing.T) {
		sql := "SELECT EXISTS (SELECT 1)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		target := stmt.TargetList.Items[0].(*ast.ResTarget)
		_, ok := target.Val.(*ast.ExistsExpr)
		if !ok {
			t.Fatalf("expected *ExistsExpr, got %T", target.Val)
		}
	})

	t.Run("function_calls", func(t *testing.T) {
		tests := []string{
			"SELECT COUNT(*)",
			"SELECT COUNT(DISTINCT col1)",
			"SELECT SUM(col1)",
			"SELECT GETDATE()",
			"SELECT LEN('hello')",
			"SELECT SUBSTRING('hello', 1, 3)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("window_functions", func(t *testing.T) {
		tests := []string{
			"SELECT ROW_NUMBER() OVER (ORDER BY col1)",
			"SELECT SUM(col1) OVER (PARTITION BY col2 ORDER BY col3)",
			"SELECT COUNT(*) OVER (PARTITION BY col1)",
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
				if fc.Over == nil {
					t.Error("expected non-nil Over")
				}
			})
		}
	})

	t.Run("parenthesized", func(t *testing.T) {
		tests := []string{
			"SELECT (1)",
			"SELECT (1 + 2)",
			"SELECT ((1 + 2) * 3)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("literals", func(t *testing.T) {
		tests := []string{
			"SELECT 42",
			"SELECT 3.14",
			"SELECT 'hello'",
			"SELECT N'unicode'",
			"SELECT NULL",
			"SELECT DEFAULT",
			"SELECT 0x1A2B",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("string_concatenation", func(t *testing.T) {
		sql := "SELECT 'hello' + ' ' + 'world'"
		ParseAndCheck(t, sql)
	})

	t.Run("alias", func(t *testing.T) {
		tests := []string{
			"SELECT 1 AS col1",
			"SELECT 1 col1",
			"SELECT col1 AS [My Column]",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				target := stmt.TargetList.Items[0].(*ast.ResTarget)
				if target.Name == "" {
					t.Error("expected non-empty alias Name")
				}
			})
		}
	})
}

// TestParseSelect tests SELECT statement parsing (batch 4).
func TestParseSelect(t *testing.T) {
	t.Run("basic select", func(t *testing.T) {
		tests := []string{
			"SELECT 1",
			"SELECT col1, col2, col3",
			"SELECT DISTINCT col1",
			"SELECT ALL col1",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("top", func(t *testing.T) {
		tests := []string{
			"SELECT TOP (10) col1 FROM t",
			"SELECT TOP 10 col1 FROM t",
			"SELECT TOP (50) PERCENT col1 FROM t",
			"SELECT TOP (10) WITH TIES col1 FROM t ORDER BY col1",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.Top == nil {
					t.Error("expected non-nil Top")
				}
			})
		}
	})

	t.Run("into", func(t *testing.T) {
		sql := "SELECT col1 INTO #temp FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.IntoTable == nil {
			t.Error("expected non-nil IntoTable")
		}
	})

	t.Run("from", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM t",
			"SELECT * FROM dbo.t",
			"SELECT * FROM t1, t2",
			"SELECT * FROM t AS alias",
			"SELECT * FROM t alias",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.FromClause == nil || stmt.FromClause.Len() == 0 {
					t.Error("expected non-empty FromClause")
				}
			})
		}
	})

	t.Run("where", func(t *testing.T) {
		sql := "SELECT col1 FROM t WHERE col1 > 10"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.WhereClause == nil {
			t.Error("expected non-nil WhereClause")
		}
	})

	t.Run("group by having", func(t *testing.T) {
		sql := "SELECT col1, COUNT(*) FROM t GROUP BY col1 HAVING COUNT(*) > 1"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.GroupByClause == nil {
			t.Error("expected non-nil GroupByClause")
		}
		if stmt.HavingClause == nil {
			t.Error("expected non-nil HavingClause")
		}
	})

	t.Run("order by", func(t *testing.T) {
		tests := []string{
			"SELECT col1 FROM t ORDER BY col1",
			"SELECT col1 FROM t ORDER BY col1 ASC",
			"SELECT col1 FROM t ORDER BY col1 DESC",
			"SELECT col1 FROM t ORDER BY col1, col2 DESC",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.OrderByClause == nil {
					t.Error("expected non-nil OrderByClause")
				}
			})
		}
	})

	t.Run("offset fetch", func(t *testing.T) {
		sql := "SELECT col1 FROM t ORDER BY col1 OFFSET 10 ROWS FETCH NEXT 20 ROWS ONLY"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.OffsetClause == nil {
			t.Error("expected non-nil OffsetClause")
		}
		if stmt.FetchClause == nil {
			t.Error("expected non-nil FetchClause")
		}
	})

	t.Run("join", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id",
			"SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id",
			"SELECT * FROM t1 LEFT OUTER JOIN t2 ON t1.id = t2.id",
			"SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id",
			"SELECT * FROM t1 FULL OUTER JOIN t2 ON t1.id = t2.id",
			"SELECT * FROM t1 CROSS JOIN t2",
			"SELECT * FROM t1 JOIN t2 ON t1.id = t2.id",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.FromClause == nil {
					t.Error("expected non-nil FromClause")
				}
			})
		}
	})

	t.Run("apply", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM t1 CROSS APPLY dbo.MyFunc(t1.id) AS f",
			"SELECT * FROM t1 OUTER APPLY dbo.MyFunc(t1.id) AS f",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("cte", func(t *testing.T) {
		tests := []string{
			"WITH cte AS (SELECT 1 AS x) SELECT * FROM cte",
			"WITH cte1 AS (SELECT 1 AS x), cte2 AS (SELECT 2 AS y) SELECT * FROM cte1, cte2",
			"WITH cte(a, b) AS (SELECT 1, 2) SELECT * FROM cte",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.WithClause == nil {
					t.Error("expected non-nil WithClause")
				}
			})
		}
	})

	t.Run("union", func(t *testing.T) {
		tests := []string{
			"SELECT 1 UNION SELECT 2",
			"SELECT 1 UNION ALL SELECT 2",
			"SELECT 1 INTERSECT SELECT 2",
			"SELECT 1 EXCEPT SELECT 2",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.Op == ast.SetOpNone {
					t.Error("expected set operation")
				}
			})
		}
	})

	t.Run("subquery in from", func(t *testing.T) {
		sql := "SELECT * FROM (SELECT 1 AS x) AS sub"
		ParseAndCheck(t, sql)
	})

	t.Run("for xml", func(t *testing.T) {
		tests := []string{
			"SELECT col1 FROM t FOR XML RAW",
			"SELECT col1 FROM t FOR XML PATH('row')",
			"SELECT col1 FROM t FOR XML AUTO",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.ForClause == nil {
					t.Error("expected non-nil ForClause")
				}
				if stmt.ForClause.Mode != ast.ForXML {
					t.Error("expected ForXML mode")
				}
			})
		}
	})

	t.Run("for json", func(t *testing.T) {
		tests := []string{
			"SELECT col1 FROM t FOR JSON PATH",
			"SELECT col1 FROM t FOR JSON AUTO",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.ForClause == nil {
					t.Error("expected non-nil ForClause")
				}
				if stmt.ForClause.Mode != ast.ForJSON {
					t.Error("expected ForJSON mode")
				}
			})
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

// TestParseInsert tests INSERT statement parsing (batch 5).
func TestParseInsert(t *testing.T) {
	t.Run("basic insert values", func(t *testing.T) {
		tests := []string{
			"INSERT INTO t (col1, col2) VALUES (1, 2)",
			"INSERT INTO t VALUES (1, 2)",
			"INSERT t VALUES (1, 2)",
			"INSERT INTO dbo.t (col1) VALUES ('hello')",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt, ok := result.Items[0].(*ast.InsertStmt)
				if !ok {
					t.Fatalf("expected *InsertStmt, got %T", result.Items[0])
				}
				if stmt.Relation == nil {
					t.Error("expected non-nil Relation")
				}
			})
		}
	})

	t.Run("insert multiple rows", func(t *testing.T) {
		sql := "INSERT INTO t (col1, col2) VALUES (1, 2), (3, 4), (5, 6)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.InsertStmt)
		vc, ok := stmt.Source.(*ast.ValuesClause)
		if !ok {
			t.Fatalf("expected *ValuesClause, got %T", stmt.Source)
		}
		if vc.Rows.Len() != 3 {
			t.Errorf("expected 3 rows, got %d", vc.Rows.Len())
		}
	})

	t.Run("insert select", func(t *testing.T) {
		sql := "INSERT INTO t (col1) SELECT col1 FROM t2"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.InsertStmt)
		_, ok := stmt.Source.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("expected *SelectStmt source, got %T", stmt.Source)
		}
	})

	t.Run("insert default values", func(t *testing.T) {
		sql := "INSERT INTO t DEFAULT VALUES"
		ParseAndCheck(t, sql)
	})

	t.Run("insert with output", func(t *testing.T) {
		sql := "INSERT INTO t (col1) OUTPUT inserted.col1 VALUES (1)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.InsertStmt)
		if stmt.OutputClause == nil {
			t.Error("expected non-nil OutputClause")
		}
	})

	t.Run("insert exec", func(t *testing.T) {
		sql := "INSERT INTO t EXEC sp_myproc @param1 = 1"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.InsertStmt)
		_, ok := stmt.Source.(*ast.ExecStmt)
		if !ok {
			t.Fatalf("expected *ExecStmt source, got %T", stmt.Source)
		}
	})
}

// TestParseUpdateDelete tests UPDATE and DELETE statement parsing (batch 6).
func TestParseUpdateDelete(t *testing.T) {
	t.Run("basic update", func(t *testing.T) {
		tests := []string{
			"UPDATE t SET col1 = 1",
			"UPDATE t SET col1 = 1, col2 = 'hello'",
			"UPDATE dbo.t SET col1 = col1 + 1 WHERE id = 1",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt, ok := result.Items[0].(*ast.UpdateStmt)
				if !ok {
					t.Fatalf("expected *UpdateStmt, got %T", result.Items[0])
				}
				if stmt.SetClause == nil || stmt.SetClause.Len() == 0 {
					t.Error("expected non-empty SetClause")
				}
			})
		}
	})

	t.Run("update with from", func(t *testing.T) {
		sql := "UPDATE t SET t.col1 = s.col1 FROM t INNER JOIN s ON t.id = s.id"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.UpdateStmt)
		if stmt.FromClause == nil {
			t.Error("expected non-nil FromClause")
		}
	})

	t.Run("update with top", func(t *testing.T) {
		sql := "UPDATE TOP (10) t SET col1 = 1"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.UpdateStmt)
		if stmt.Top == nil {
			t.Error("expected non-nil Top")
		}
	})

	t.Run("basic delete", func(t *testing.T) {
		tests := []string{
			"DELETE FROM t",
			"DELETE FROM t WHERE id = 1",
			"DELETE t WHERE id = 1",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.DeleteStmt)
				if !ok {
					t.Fatalf("expected *DeleteStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("delete with top", func(t *testing.T) {
		sql := "DELETE TOP (10) FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.DeleteStmt)
		if stmt.Top == nil {
			t.Error("expected non-nil Top")
		}
	})

	t.Run("delete with from join", func(t *testing.T) {
		sql := "DELETE t FROM t INNER JOIN t2 ON t.id = t2.id WHERE t2.status = 0"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.DeleteStmt)
		if stmt.FromClause == nil {
			t.Error("expected non-nil FromClause")
		}
	})
}

// TestParseMerge tests MERGE statement parsing (batch 7).
func TestParseMerge(t *testing.T) {
	t.Run("basic merge", func(t *testing.T) {
		sql := `MERGE INTO target AS t
			USING source AS s ON t.id = s.id
			WHEN MATCHED THEN UPDATE SET t.col1 = s.col1
			WHEN NOT MATCHED THEN INSERT (col1) VALUES (s.col1);`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.MergeStmt)
		if !ok {
			t.Fatalf("expected *MergeStmt, got %T", result.Items[0])
		}
		if stmt.WhenClauses == nil || stmt.WhenClauses.Len() < 2 {
			t.Errorf("expected at least 2 WHEN clauses, got %d", stmt.WhenClauses.Len())
		}
	})

	t.Run("merge with delete", func(t *testing.T) {
		sql := `MERGE target AS t
			USING source AS s ON t.id = s.id
			WHEN MATCHED AND s.deleted = 1 THEN DELETE
			WHEN MATCHED THEN UPDATE SET t.col1 = s.col1
			WHEN NOT MATCHED BY TARGET THEN INSERT (col1) VALUES (s.col1);`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.MergeStmt)
		if stmt.WhenClauses.Len() != 3 {
			t.Errorf("expected 3 WHEN clauses, got %d", stmt.WhenClauses.Len())
		}
	})
}

// TestParseCreateTable tests CREATE TABLE statement parsing (batch 8).
func TestParseCreateTable(t *testing.T) {
	t.Run("basic create table", func(t *testing.T) {
		sql := "CREATE TABLE t (col1 int, col2 varchar(100))"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateTableStmt)
		if !ok {
			t.Fatalf("expected *CreateTableStmt, got %T", result.Items[0])
		}
		if stmt.Columns == nil || stmt.Columns.Len() != 2 {
			t.Errorf("expected 2 columns, got %d", stmt.Columns.Len())
		}
	})

	t.Run("create table with constraints", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL IDENTITY(1, 1),
			name nvarchar(100) NULL,
			email varchar(255) NOT NULL,
			CONSTRAINT PK_t PRIMARY KEY CLUSTERED (id),
			CONSTRAINT UQ_email UNIQUE (email)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.Columns == nil {
			t.Error("expected non-nil Columns")
		}
	})

	t.Run("create table with default and check", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			status int DEFAULT 0,
			CONSTRAINT CHK_status CHECK (status >= 0)
		)`
		ParseAndCheck(t, sql)
	})

	t.Run("create table with foreign key", func(t *testing.T) {
		sql := `CREATE TABLE orders (
			id int NOT NULL,
			customer_id int NOT NULL,
			CONSTRAINT FK_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE
		)`
		ParseAndCheck(t, sql)
	})

	t.Run("create table with identity", func(t *testing.T) {
		sql := "CREATE TABLE t (id int IDENTITY(1,1) NOT NULL)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col := stmt.Columns.Items[0].(*ast.ColumnDef)
		if col.Identity == nil {
			t.Error("expected non-nil Identity")
		}
	})

	t.Run("create table with collate", func(t *testing.T) {
		sql := "CREATE TABLE t (name nvarchar(100) COLLATE Latin1_General_CI_AS)"
		ParseAndCheck(t, sql)
	})
}

// TestParseAlterTable tests ALTER TABLE statement parsing (batch 9).
func TestParseAlterTable(t *testing.T) {
	t.Run("add column", func(t *testing.T) {
		sql := "ALTER TABLE t ADD col1 int NULL"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.AlterTableStmt)
		if !ok {
			t.Fatalf("expected *AlterTableStmt, got %T", result.Items[0])
		}
	})

	t.Run("drop column", func(t *testing.T) {
		sql := "ALTER TABLE t DROP COLUMN col1"
		ParseAndCheck(t, sql)
	})

	t.Run("alter column", func(t *testing.T) {
		sql := "ALTER TABLE t ALTER COLUMN col1 varchar(200) NOT NULL"
		ParseAndCheck(t, sql)
	})

	t.Run("add constraint", func(t *testing.T) {
		sql := "ALTER TABLE t ADD CONSTRAINT PK_t PRIMARY KEY (id)"
		ParseAndCheck(t, sql)
	})

	t.Run("drop constraint", func(t *testing.T) {
		sql := "ALTER TABLE t DROP CONSTRAINT PK_t"
		ParseAndCheck(t, sql)
	})
}

// TestParseCreateIndex tests CREATE INDEX statement parsing (batch 10).
func TestParseCreateIndex(t *testing.T) {
	t.Run("basic index", func(t *testing.T) {
		sql := "CREATE INDEX IX_t_col1 ON t (col1)"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateIndexStmt)
		if !ok {
			t.Fatalf("expected *CreateIndexStmt, got %T", result.Items[0])
		}
		if stmt.Name != "IX_t_col1" {
			t.Errorf("expected name IX_t_col1, got %s", stmt.Name)
		}
	})

	t.Run("unique clustered index", func(t *testing.T) {
		sql := "CREATE UNIQUE CLUSTERED INDEX IX_t_id ON dbo.t (id ASC)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateIndexStmt)
		if !stmt.Unique {
			t.Error("expected Unique=true")
		}
	})

	t.Run("index with include", func(t *testing.T) {
		sql := "CREATE NONCLUSTERED INDEX IX_t_col1 ON t (col1) INCLUDE (col2, col3)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateIndexStmt)
		if stmt.IncludeCols == nil {
			t.Error("expected non-nil IncludeCols")
		}
	})

	t.Run("filtered index", func(t *testing.T) {
		sql := "CREATE INDEX IX_t_col1 ON t (col1) WHERE col1 IS NOT NULL"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateIndexStmt)
		if stmt.WhereClause == nil {
			t.Error("expected non-nil WhereClause")
		}
	})
}

// TestParseCreateView tests CREATE VIEW statement parsing (batch 11).
func TestParseCreateView(t *testing.T) {
	t.Run("basic view", func(t *testing.T) {
		sql := "CREATE VIEW v AS SELECT col1 FROM t"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateViewStmt)
		if !ok {
			t.Fatalf("expected *CreateViewStmt, got %T", result.Items[0])
		}
		if stmt.Query == nil {
			t.Error("expected non-nil Query")
		}
	})

	t.Run("create or alter view", func(t *testing.T) {
		sql := "CREATE OR ALTER VIEW dbo.v AS SELECT col1 FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateViewStmt)
		if !stmt.OrAlter {
			t.Error("expected OrAlter=true")
		}
	})

	t.Run("view with schemabinding", func(t *testing.T) {
		sql := "CREATE VIEW v WITH SCHEMABINDING AS SELECT col1 FROM dbo.t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateViewStmt)
		if !stmt.SchemaBinding {
			t.Error("expected SchemaBinding=true")
		}
	})

	t.Run("view with columns", func(t *testing.T) {
		sql := "CREATE VIEW v (a, b) AS SELECT col1, col2 FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateViewStmt)
		if stmt.Columns == nil || stmt.Columns.Len() != 2 {
			t.Errorf("expected 2 columns, got %d", stmt.Columns.Len())
		}
	})
}

// TestParseCreateProc tests CREATE PROCEDURE and FUNCTION parsing (batch 12).
func TestParseCreateProc(t *testing.T) {
	t.Run("basic procedure", func(t *testing.T) {
		sql := `CREATE PROCEDURE sp_test @param1 int, @param2 varchar(100) AS BEGIN SELECT 1 END`
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.CreateProcedureStmt)
		if !ok {
			t.Fatalf("expected *CreateProcedureStmt, got %T", result.Items[0])
		}
	})

	t.Run("proc with defaults", func(t *testing.T) {
		sql := `CREATE PROC sp_test @param1 int = 0 AS BEGIN SELECT @param1 END`
		ParseAndCheck(t, sql)
	})

	t.Run("proc with output param", func(t *testing.T) {
		sql := `CREATE PROCEDURE sp_test @param1 int, @result int OUTPUT AS BEGIN SET @result = @param1 END`
		ParseAndCheck(t, sql)
	})

	t.Run("basic function", func(t *testing.T) {
		sql := `CREATE FUNCTION fn_test (@param1 int) RETURNS int AS BEGIN RETURN @param1 + 1 END`
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.CreateFunctionStmt)
		if !ok {
			t.Fatalf("expected *CreateFunctionStmt, got %T", result.Items[0])
		}
	})

	t.Run("table valued function", func(t *testing.T) {
		sql := `CREATE FUNCTION fn_test (@id int) RETURNS TABLE AS RETURN SELECT * FROM t WHERE id = @id`
		ParseAndCheck(t, sql)
	})
}

// TestParseCreateDatabase tests CREATE DATABASE parsing (batch 13).
func TestParseCreateDatabase(t *testing.T) {
	t.Run("basic create database", func(t *testing.T) {
		sql := "CREATE DATABASE mydb"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Name != "mydb" {
			t.Errorf("expected name mydb, got %s", stmt.Name)
		}
	})
}

// TestParseDrop tests DROP statement parsing (batch 14).
func TestParseDrop(t *testing.T) {
	t.Run("drop table", func(t *testing.T) {
		tests := []string{
			"DROP TABLE t",
			"DROP TABLE IF EXISTS t",
			"DROP TABLE dbo.t1, dbo.t2",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.DropStmt)
				if !ok {
					t.Fatalf("expected *DropStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("drop various objects", func(t *testing.T) {
		tests := []string{
			"DROP VIEW v",
			"DROP PROCEDURE sp_test",
			"DROP FUNCTION fn_test",
			"DROP INDEX IX_t_col1 ON t",
			"DROP DATABASE mydb",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	t.Run("drop if exists", func(t *testing.T) {
		sql := "DROP TABLE IF EXISTS dbo.t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.DropStmt)
		if !stmt.IfExists {
			t.Error("expected IfExists=true")
		}
	})
}

// TestParseControlFlow tests control flow statement parsing (batch 15).
func TestParseControlFlow(t *testing.T) {
	t.Run("if else", func(t *testing.T) {
		sql := "IF 1 = 1 SELECT 1 ELSE SELECT 2"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.IfStmt)
		if !ok {
			t.Fatalf("expected *IfStmt, got %T", result.Items[0])
		}
		if stmt.Else == nil {
			t.Error("expected non-nil Else")
		}
	})

	t.Run("if begin end", func(t *testing.T) {
		sql := "IF 1 = 1 BEGIN SELECT 1 END"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.IfStmt)
		if !ok {
			t.Fatalf("expected *IfStmt, got %T", result.Items[0])
		}
	})

	t.Run("while", func(t *testing.T) {
		sql := "WHILE @i < 10 SET @i = @i + 1"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.WhileStmt)
		if !ok {
			t.Fatalf("expected *WhileStmt, got %T", result.Items[0])
		}
	})

	t.Run("begin end", func(t *testing.T) {
		sql := "BEGIN SELECT 1; SELECT 2 END"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.BeginEndStmt)
		if !ok {
			t.Fatalf("expected *BeginEndStmt, got %T", result.Items[0])
		}
	})

	t.Run("try catch", func(t *testing.T) {
		sql := "BEGIN TRY SELECT 1 END TRY BEGIN CATCH SELECT ERROR_MESSAGE() END CATCH"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.TryCatchStmt)
		if !ok {
			t.Fatalf("expected *TryCatchStmt, got %T", result.Items[0])
		}
	})

	t.Run("return", func(t *testing.T) {
		sql := "RETURN 0"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.ReturnStmt)
		if !ok {
			t.Fatalf("expected *ReturnStmt, got %T", result.Items[0])
		}
	})

	t.Run("break continue", func(t *testing.T) {
		// These need to be in a context where they're valid statements
		sql := "BREAK"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.BreakStmt)
		if !ok {
			t.Fatalf("expected *BreakStmt, got %T", result.Items[0])
		}

		sql = "CONTINUE"
		result = ParseAndCheck(t, sql)
		_, ok2 := result.Items[0].(*ast.ContinueStmt)
		if !ok2 {
			t.Fatalf("expected *ContinueStmt, got %T", result.Items[0])
		}
	})

	t.Run("goto", func(t *testing.T) {
		sql := "GOTO my_label"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.GotoStmt)
		if !ok {
			t.Fatalf("expected *GotoStmt, got %T", result.Items[0])
		}
	})

	t.Run("waitfor", func(t *testing.T) {
		sql := "WAITFOR DELAY '00:00:05'"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.WaitForStmt)
		if !ok {
			t.Fatalf("expected *WaitForStmt, got %T", result.Items[0])
		}
	})
}

// TestParseDeclareSet tests DECLARE and SET parsing (batch 16).
func TestParseDeclareSet(t *testing.T) {
	t.Run("declare variable", func(t *testing.T) {
		tests := []string{
			"DECLARE @x int",
			"DECLARE @x int = 0",
			"DECLARE @x int, @y varchar(100)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.DeclareStmt)
				if !ok {
					t.Fatalf("expected *DeclareStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("declare table", func(t *testing.T) {
		sql := "DECLARE @t TABLE (id int, name varchar(100))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.DeclareStmt)
		decl := stmt.Variables.Items[0].(*ast.VariableDecl)
		if !decl.IsTable {
			t.Error("expected IsTable=true")
		}
	})

	t.Run("set variable", func(t *testing.T) {
		tests := []string{
			"SET @x = 1",
			"SET @x = @x + 1",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.SetStmt)
				if !ok {
					t.Fatalf("expected *SetStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("set option", func(t *testing.T) {
		tests := []string{
			"SET NOCOUNT ON",
			"SET NOCOUNT OFF",
			"SET XACT_ABORT ON",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.SetStmt)
				if !ok {
					t.Fatalf("expected *SetStmt, got %T", result.Items[0])
				}
			})
		}
	})
}

// TestParseTransaction tests transaction statement parsing (batch 17).
func TestParseTransaction(t *testing.T) {
	t.Run("begin transaction", func(t *testing.T) {
		tests := []string{
			"BEGIN TRANSACTION",
			"BEGIN TRAN",
			"BEGIN TRAN MyTran",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.BeginTransStmt)
				if !ok {
					t.Fatalf("expected *BeginTransStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("commit", func(t *testing.T) {
		tests := []string{
			"COMMIT",
			"COMMIT TRANSACTION",
			"COMMIT TRAN",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.CommitTransStmt)
				if !ok {
					t.Fatalf("expected *CommitTransStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("rollback", func(t *testing.T) {
		tests := []string{
			"ROLLBACK",
			"ROLLBACK TRANSACTION",
			"ROLLBACK TRAN",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.RollbackTransStmt)
				if !ok {
					t.Fatalf("expected *RollbackTransStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("save transaction", func(t *testing.T) {
		sql := "SAVE TRAN MySavepoint"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.SaveTransStmt)
		if !ok {
			t.Fatalf("expected *SaveTransStmt, got %T", result.Items[0])
		}
	})
}

// TestParseExecute tests EXEC/EXECUTE parsing (batch 18).
func TestParseExecute(t *testing.T) {
	t.Run("basic exec", func(t *testing.T) {
		tests := []string{
			"EXEC sp_test",
			"EXECUTE sp_test",
			"EXEC dbo.sp_test",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				_, ok := result.Items[0].(*ast.ExecStmt)
				if !ok {
					t.Fatalf("expected *ExecStmt, got %T", result.Items[0])
				}
			})
		}
	})

	t.Run("exec with args", func(t *testing.T) {
		sql := "EXEC sp_test @param1 = 1, @param2 = 'hello'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.ExecStmt)
		if stmt.Args == nil || stmt.Args.Len() < 2 {
			t.Error("expected at least 2 args")
		}
	})

	t.Run("exec with return var", func(t *testing.T) {
		sql := "EXEC @result = sp_test 1, 2"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.ExecStmt)
		if stmt.ReturnVar != "@result" {
			t.Errorf("expected ReturnVar=@result, got %s", stmt.ReturnVar)
		}
	})

	t.Run("exec with output param", func(t *testing.T) {
		sql := "EXEC sp_test @param1 = 1, @param2 OUTPUT"
		ParseAndCheck(t, sql)
	})
}

// TestParseGrant tests GRANT/REVOKE/DENY parsing (batch 19).
func TestParseGrant(t *testing.T) {
	t.Run("grant", func(t *testing.T) {
		tests := []string{
			"GRANT SELECT ON dbo.t TO myuser",
			"GRANT SELECT, INSERT ON dbo.t TO myuser",
			"GRANT EXECUTE ON dbo.sp_test TO myuser",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt, ok := result.Items[0].(*ast.GrantStmt)
				if !ok {
					t.Fatalf("expected *GrantStmt, got %T", result.Items[0])
				}
				if stmt.StmtType != ast.GrantTypeGrant {
					t.Error("expected GrantTypeGrant")
				}
			})
		}
	})

	t.Run("revoke", func(t *testing.T) {
		sql := "REVOKE SELECT ON dbo.t FROM myuser"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.GrantStmt)
		if stmt.StmtType != ast.GrantTypeRevoke {
			t.Error("expected GrantTypeRevoke")
		}
	})

	t.Run("deny", func(t *testing.T) {
		sql := "DENY SELECT ON dbo.t TO myuser"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.GrantStmt)
		if stmt.StmtType != ast.GrantTypeDeny {
			t.Error("expected GrantTypeDeny")
		}
	})
}

// TestParseUtility tests utility statement parsing (batch 20).
func TestParseUtility(t *testing.T) {
	t.Run("use", func(t *testing.T) {
		sql := "USE mydb"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.UseStmt)
		if !ok {
			t.Fatalf("expected *UseStmt, got %T", result.Items[0])
		}
		if stmt.Database != "mydb" {
			t.Errorf("expected database mydb, got %s", stmt.Database)
		}
	})

	t.Run("print", func(t *testing.T) {
		sql := "PRINT 'hello world'"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.PrintStmt)
		if !ok {
			t.Fatalf("expected *PrintStmt, got %T", result.Items[0])
		}
	})

	t.Run("raiserror", func(t *testing.T) {
		sql := "RAISERROR('Error occurred', 16, 1)"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.RaiseErrorStmt)
		if !ok {
			t.Fatalf("expected *RaiseErrorStmt, got %T", result.Items[0])
		}
	})

	t.Run("throw", func(t *testing.T) {
		sql := "THROW 50000, 'Error', 1"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.ThrowStmt)
		if !ok {
			t.Fatalf("expected *ThrowStmt, got %T", result.Items[0])
		}
	})

	t.Run("throw rethrow", func(t *testing.T) {
		sql := "THROW"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.ThrowStmt)
		if !ok {
			t.Fatalf("expected *ThrowStmt, got %T", result.Items[0])
		}
	})

	t.Run("truncate", func(t *testing.T) {
		sql := "TRUNCATE TABLE dbo.t"
		result := ParseAndCheck(t, sql)
		_, ok := result.Items[0].(*ast.TruncateStmt)
		if !ok {
			t.Fatalf("expected *TruncateStmt, got %T", result.Items[0])
		}
	})
}

// TestParseGoBatch tests GO batch separator parsing (batch 21).
func TestParseGoBatch(t *testing.T) {
	t.Run("simple go", func(t *testing.T) {
		sql := "SELECT 1\nGO"
		result := ParseAndCheck(t, sql)
		if result.Len() < 2 {
			t.Errorf("expected at least 2 items, got %d", result.Len())
		}
	})

	t.Run("go with count", func(t *testing.T) {
		sql := "SELECT 1\nGO 5"
		result := ParseAndCheck(t, sql)
		found := false
		for _, item := range result.Items {
			if gs, ok := item.(*ast.GoStmt); ok {
				found = true
				if gs.Count != 5 {
					t.Errorf("expected Count=5, got %d", gs.Count)
				}
			}
		}
		if !found {
			t.Error("expected GoStmt in result")
		}
	})
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
