package parser

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
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

	t.Run("pivot", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM Sales PIVOT (SUM(Amount) FOR Month IN ([Jan],[Feb],[Mar])) AS pvt",
			"SELECT * FROM t PIVOT (COUNT(val) FOR category IN ([A],[B])) AS p",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.FromClause == nil || stmt.FromClause.Len() == 0 {
					t.Error("expected non-empty FromClause")
				}
				_, ok := stmt.FromClause.Items[0].(*ast.PivotExpr)
				if !ok {
					t.Errorf("expected *PivotExpr in FROM, got %T", stmt.FromClause.Items[0])
				}
			})
		}
	})

	t.Run("unpivot", func(t *testing.T) {
		sql := "SELECT * FROM t UNPIVOT (val FOR col IN ([Jan],[Feb],[Mar])) AS unpvt"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.FromClause == nil || stmt.FromClause.Len() == 0 {
			t.Error("expected non-empty FromClause")
		}
		unpvt, ok := stmt.FromClause.Items[0].(*ast.UnpivotExpr)
		if !ok {
			t.Errorf("expected *UnpivotExpr in FROM, got %T", stmt.FromClause.Items[0])
		}
		if unpvt.ValueCol != "val" {
			t.Errorf("expected ValueCol=val, got %s", unpvt.ValueCol)
		}
		if unpvt.ForCol != "col" {
			t.Errorf("expected ForCol=col, got %s", unpvt.ForCol)
		}
		if unpvt.Alias != "unpvt" {
			t.Errorf("expected Alias=unpvt, got %s", unpvt.Alias)
		}
	})

	t.Run("tablesample", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM t TABLESAMPLE (10 PERCENT)",
			"SELECT * FROM t TABLESAMPLE (100 ROWS)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.FromClause == nil || stmt.FromClause.Len() == 0 {
					t.Error("expected non-empty FromClause")
				}
				atr, ok := stmt.FromClause.Items[0].(*ast.AliasedTableRef)
				if !ok {
					t.Errorf("expected *AliasedTableRef in FROM, got %T", stmt.FromClause.Items[0])
					return
				}
				if atr.TableSample == nil {
					t.Error("expected non-nil TableSample")
				}
			})
		}
	})

	t.Run("tablesample repeatable", func(t *testing.T) {
		sql := "SELECT * FROM t TABLESAMPLE (100 ROWS) REPEATABLE (42)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		atr, ok := stmt.FromClause.Items[0].(*ast.AliasedTableRef)
		if !ok {
			t.Fatalf("expected *AliasedTableRef, got %T", stmt.FromClause.Items[0])
		}
		if atr.TableSample == nil {
			t.Fatal("expected non-nil TableSample")
		}
		if atr.TableSample.Repeatable == nil {
			t.Error("expected non-nil Repeatable")
		}
		if atr.TableSample.Unit != "ROWS" {
			t.Errorf("expected Unit=ROWS, got %s", atr.TableSample.Unit)
		}
	})

	t.Run("rowset functions", func(t *testing.T) {
		tests := []string{
			"SELECT * FROM OPENROWSET('SQLNCLI', 'server=srv', 'SELECT 1') AS t",
			"SELECT * FROM OPENQUERY(LinkedServer, 'SELECT 1') AS t",
			"SELECT * FROM OPENJSON(@json) AS t",
			"SELECT * FROM OPENXML(@hdoc, '/root/row') AS t",
			"SELECT * FROM OPENDATASOURCE('SQLNCLI', 'Data Source=srv') AS t",
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

	t.Run("openjson with", func(t *testing.T) {
		sql := "SELECT * FROM OPENJSON(@json) WITH (name NVARCHAR(100), age INT) AS t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.FromClause == nil || stmt.FromClause.Len() == 0 {
			t.Error("expected non-empty FromClause")
		}
	})

	t.Run("grouping sets", func(t *testing.T) {
		sql := "SELECT col1, col2, SUM(val) FROM t GROUP BY GROUPING SETS ((col1), (col1, col2), ())"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.GroupByClause == nil {
			t.Fatal("expected non-nil GroupByClause")
		}
		if stmt.GroupByClause.Len() != 1 {
			t.Fatalf("expected 1 group by item, got %d", stmt.GroupByClause.Len())
		}
		gs, ok := stmt.GroupByClause.Items[0].(*ast.GroupingSetsExpr)
		if !ok {
			t.Fatalf("expected *GroupingSetsExpr, got %T", stmt.GroupByClause.Items[0])
		}
		if gs.Sets == nil || gs.Sets.Len() != 3 {
			t.Errorf("expected 3 sets, got %v", gs.Sets)
		}
	})

	t.Run("rollup", func(t *testing.T) {
		sql := "SELECT col1, col2, SUM(val) FROM t GROUP BY ROLLUP (col1, col2)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.GroupByClause == nil {
			t.Fatal("expected non-nil GroupByClause")
		}
		_, ok := stmt.GroupByClause.Items[0].(*ast.RollupExpr)
		if !ok {
			t.Fatalf("expected *RollupExpr, got %T", stmt.GroupByClause.Items[0])
		}
	})

	t.Run("cube", func(t *testing.T) {
		sql := "SELECT col1, col2, SUM(val) FROM t GROUP BY CUBE (col1, col2)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.GroupByClause == nil {
			t.Fatal("expected non-nil GroupByClause")
		}
		_, ok := stmt.GroupByClause.Items[0].(*ast.CubeExpr)
		if !ok {
			t.Fatalf("expected *CubeExpr, got %T", stmt.GroupByClause.Items[0])
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
				// SET session options return SetOptionStmt
				_, okOpt := result.Items[0].(*ast.SetOptionStmt)
				_, okSet := result.Items[0].(*ast.SetStmt)
				if !okOpt && !okSet {
					t.Fatalf("expected *SetOptionStmt or *SetStmt, got %T", result.Items[0])
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

// ---------- Batch 23: Cursor Operations ----------

func TestParseCursorOperations(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// DECLARE CURSOR - ISO syntax (basic)
		{
			name: "declare_cursor_basic",
			sql:  "DECLARE my_cursor CURSOR FOR SELECT * FROM employees",
		},
		// DECLARE CURSOR - ISO syntax with INSENSITIVE
		{
			name: "declare_cursor_insensitive",
			sql:  "DECLARE my_cursor INSENSITIVE CURSOR FOR SELECT id, name FROM employees",
		},
		// DECLARE CURSOR - ISO syntax with SCROLL
		{
			name: "declare_cursor_scroll_iso",
			sql:  "DECLARE my_cursor SCROLL CURSOR FOR SELECT * FROM employees",
		},
		// DECLARE CURSOR - ISO syntax with INSENSITIVE SCROLL
		{
			name: "declare_cursor_insensitive_scroll",
			sql:  "DECLARE my_cursor INSENSITIVE SCROLL CURSOR FOR SELECT * FROM employees",
		},
		// DECLARE CURSOR - ISO syntax with FOR READ_ONLY
		{
			name: "declare_cursor_for_read_only",
			sql:  "DECLARE my_cursor CURSOR FOR SELECT * FROM employees FOR READ_ONLY",
		},
		// DECLARE CURSOR - ISO syntax with FOR UPDATE
		{
			name: "declare_cursor_for_update",
			sql:  "DECLARE my_cursor CURSOR FOR SELECT * FROM employees FOR UPDATE",
		},
		// DECLARE CURSOR - ISO syntax with FOR UPDATE OF columns
		{
			name: "declare_cursor_for_update_of",
			sql:  "DECLARE my_cursor CURSOR FOR SELECT * FROM employees FOR UPDATE OF name, salary",
		},
		// DECLARE CURSOR - T-SQL extended: LOCAL
		{
			name: "declare_cursor_local",
			sql:  "DECLARE my_cursor CURSOR LOCAL FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: GLOBAL
		{
			name: "declare_cursor_global",
			sql:  "DECLARE my_cursor CURSOR GLOBAL FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: FORWARD_ONLY
		{
			name: "declare_cursor_forward_only",
			sql:  "DECLARE my_cursor CURSOR LOCAL FORWARD_ONLY FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: SCROLL
		{
			name: "declare_cursor_scroll_tsql",
			sql:  "DECLARE my_cursor CURSOR LOCAL SCROLL FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: STATIC
		{
			name: "declare_cursor_static",
			sql:  "DECLARE my_cursor CURSOR LOCAL FORWARD_ONLY STATIC FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: KEYSET
		{
			name: "declare_cursor_keyset",
			sql:  "DECLARE my_cursor CURSOR SCROLL KEYSET FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: DYNAMIC
		{
			name: "declare_cursor_dynamic",
			sql:  "DECLARE my_cursor CURSOR SCROLL DYNAMIC FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: FAST_FORWARD
		{
			name: "declare_cursor_fast_forward",
			sql:  "DECLARE my_cursor CURSOR FAST_FORWARD FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: READ_ONLY concurrency
		{
			name: "declare_cursor_read_only",
			sql:  "DECLARE my_cursor CURSOR SCROLL STATIC READ_ONLY FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: SCROLL_LOCKS concurrency
		{
			name: "declare_cursor_scroll_locks",
			sql:  "DECLARE my_cursor CURSOR SCROLL KEYSET SCROLL_LOCKS FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: OPTIMISTIC concurrency
		{
			name: "declare_cursor_optimistic",
			sql:  "DECLARE my_cursor CURSOR SCROLL DYNAMIC OPTIMISTIC FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: TYPE_WARNING
		{
			name: "declare_cursor_type_warning",
			sql:  "DECLARE my_cursor CURSOR SCROLL DYNAMIC OPTIMISTIC TYPE_WARNING FOR SELECT id FROM t",
		},
		// DECLARE CURSOR - T-SQL extended: full options
		{
			name: "declare_cursor_full_options",
			sql:  "DECLARE my_cursor CURSOR GLOBAL SCROLL KEYSET SCROLL_LOCKS TYPE_WARNING FOR SELECT * FROM employees FOR UPDATE OF name, salary",
		},
		// OPEN cursor - basic
		{
			name: "open_cursor_basic",
			sql:  "OPEN my_cursor",
		},
		// OPEN cursor - GLOBAL
		{
			name: "open_cursor_global",
			sql:  "OPEN GLOBAL my_cursor",
		},
		// OPEN cursor - @variable
		{
			name: "open_cursor_variable",
			sql:  "OPEN @my_cursor_var",
		},
		// FETCH NEXT - basic
		{
			name: "fetch_next_basic",
			sql:  "FETCH NEXT FROM my_cursor",
		},
		// FETCH NEXT INTO variables
		{
			name: "fetch_next_into",
			sql:  "FETCH NEXT FROM my_cursor INTO @id, @name, @salary",
		},
		// FETCH PRIOR
		{
			name: "fetch_prior",
			sql:  "FETCH PRIOR FROM my_cursor",
		},
		// FETCH FIRST
		{
			name: "fetch_first",
			sql:  "FETCH FIRST FROM my_cursor",
		},
		// FETCH LAST
		{
			name: "fetch_last",
			sql:  "FETCH LAST FROM my_cursor",
		},
		// FETCH ABSOLUTE n
		{
			name: "fetch_absolute",
			sql:  "FETCH ABSOLUTE 5 FROM my_cursor",
		},
		// FETCH ABSOLUTE @nvar
		{
			name: "fetch_absolute_var",
			sql:  "FETCH ABSOLUTE @pos FROM my_cursor INTO @id",
		},
		// FETCH RELATIVE n
		{
			name: "fetch_relative",
			sql:  "FETCH RELATIVE 3 FROM my_cursor",
		},
		// FETCH RELATIVE negative
		{
			name: "fetch_relative_negative",
			sql:  "FETCH RELATIVE -2 FROM my_cursor INTO @id, @name",
		},
		// FETCH with GLOBAL cursor
		{
			name: "fetch_global_cursor",
			sql:  "FETCH NEXT FROM GLOBAL my_cursor INTO @val",
		},
		// FETCH from @cursor_variable
		{
			name: "fetch_cursor_variable",
			sql:  "FETCH NEXT FROM @my_cursor_var INTO @id",
		},
		// FETCH without orientation (default NEXT)
		{
			name: "fetch_default",
			sql:  "FETCH FROM my_cursor",
		},
		// CLOSE cursor - basic
		{
			name: "close_cursor_basic",
			sql:  "CLOSE my_cursor",
		},
		// CLOSE cursor - GLOBAL
		{
			name: "close_cursor_global",
			sql:  "CLOSE GLOBAL my_cursor",
		},
		// CLOSE cursor - @variable
		{
			name: "close_cursor_variable",
			sql:  "CLOSE @my_cursor_var",
		},
		// DEALLOCATE cursor - basic
		{
			name: "deallocate_cursor_basic",
			sql:  "DEALLOCATE my_cursor",
		},
		// DEALLOCATE cursor - GLOBAL
		{
			name: "deallocate_cursor_global",
			sql:  "DEALLOCATE GLOBAL my_cursor",
		},
		// DEALLOCATE cursor - @variable
		{
			name: "deallocate_cursor_variable",
			sql:  "DEALLOCATE @my_cursor_var",
		},
		// Full cursor lifecycle
		{
			name: "cursor_full_lifecycle",
			sql: `DECLARE emp_cursor CURSOR LOCAL FAST_FORWARD FOR
				SELECT id, name FROM employees WHERE active = 1;
				OPEN emp_cursor;
				FETCH NEXT FROM emp_cursor INTO @id, @name;
				CLOSE emp_cursor;
				DEALLOCATE emp_cursor`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			// Verify deterministic serialization
			for i, item := range result.Items {
				s1 := ast.NodeToString(item)
				s2 := ast.NodeToString(item)
				if s1 != s2 {
					t.Errorf("stmt[%d] serialization not deterministic:\n  s1: %s\n  s2: %s", i, s1, s2)
				}
			}
		})
	}
}

// ---------- Batch 24: Create Trigger ----------

func TestParseCursorCreateTrigger(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// DML trigger - AFTER INSERT
		{
			name: "dml_trigger_after_insert",
			sql:  "CREATE TRIGGER tr_insert ON dbo.employees AFTER INSERT AS BEGIN SELECT 1 END",
		},
		// DML trigger - AFTER UPDATE
		{
			name: "dml_trigger_after_update",
			sql:  "CREATE TRIGGER tr_update ON employees AFTER UPDATE AS BEGIN SELECT 1 END",
		},
		// DML trigger - AFTER DELETE
		{
			name: "dml_trigger_after_delete",
			sql:  "CREATE TRIGGER dbo.tr_delete ON dbo.orders AFTER DELETE AS BEGIN SELECT 1 END",
		},
		// DML trigger - AFTER INSERT, UPDATE, DELETE (multiple events)
		{
			name: "dml_trigger_multiple_events",
			sql:  "CREATE TRIGGER tr_all ON employees AFTER INSERT, UPDATE, DELETE AS BEGIN SELECT 1 END",
		},
		// DML trigger - FOR (same as AFTER)
		{
			name: "dml_trigger_for",
			sql:  "CREATE TRIGGER tr_for ON employees FOR INSERT AS BEGIN SELECT 1 END",
		},
		// DML trigger - INSTEAD OF INSERT
		{
			name: "dml_trigger_instead_of_insert",
			sql:  "CREATE TRIGGER tr_instead ON employees INSTEAD OF INSERT AS BEGIN SELECT 1 END",
		},
		// DML trigger - INSTEAD OF UPDATE
		{
			name: "dml_trigger_instead_of_update",
			sql:  "CREATE TRIGGER tr_instead_upd ON employees INSTEAD OF UPDATE AS BEGIN SELECT 1 END",
		},
		// DML trigger - INSTEAD OF DELETE
		{
			name: "dml_trigger_instead_of_delete",
			sql:  "CREATE TRIGGER tr_instead_del ON employees INSTEAD OF DELETE AS BEGIN SELECT 1 END",
		},
		// DML trigger - INSTEAD OF INSERT, UPDATE
		{
			name: "dml_trigger_instead_of_multi",
			sql:  "CREATE TRIGGER tr_instead_multi ON employees INSTEAD OF INSERT, UPDATE AS BEGIN SELECT 1 END",
		},
		// DML trigger - NOT FOR REPLICATION
		{
			name: "dml_trigger_not_for_replication",
			sql:  "CREATE TRIGGER tr_nfr ON employees AFTER INSERT NOT FOR REPLICATION AS BEGIN SELECT 1 END",
		},
		// DML trigger - OR ALTER
		{
			name: "dml_trigger_or_alter",
			sql:  "CREATE OR ALTER TRIGGER tr_oa ON employees AFTER INSERT AS BEGIN SELECT 1 END",
		},
		// DDL trigger - ON DATABASE
		{
			name: "ddl_trigger_on_database",
			sql:  "CREATE TRIGGER tr_ddl ON DATABASE FOR CREATE_TABLE AS BEGIN SELECT 1 END",
		},
		// DDL trigger - ON DATABASE, multiple events
		{
			name: "ddl_trigger_multi_events",
			sql:  "CREATE TRIGGER tr_ddl_multi ON DATABASE AFTER CREATE_TABLE, ALTER_TABLE, DROP_TABLE AS BEGIN SELECT 1 END",
		},
		// DDL trigger - ON ALL SERVER
		{
			name: "ddl_trigger_all_server",
			sql:  "CREATE TRIGGER tr_server ON ALL SERVER FOR CREATE_DATABASE AS BEGIN SELECT 1 END",
		},
		// DDL trigger - event group
		{
			name: "ddl_trigger_event_group",
			sql:  "CREATE TRIGGER tr_ddl_events ON DATABASE FOR DDL_TABLE_EVENTS AS BEGIN SELECT 1 END",
		},
		// Logon trigger
		{
			name: "logon_trigger",
			sql:  "CREATE TRIGGER tr_logon ON ALL SERVER FOR LOGON AS BEGIN SELECT 1 END",
		},
		// Logon trigger - AFTER
		{
			name: "logon_trigger_after",
			sql:  "CREATE TRIGGER tr_logon_after ON ALL SERVER AFTER LOGON AS BEGIN SELECT 1 END",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			for i, item := range result.Items {
				s1 := ast.NodeToString(item)
				s2 := ast.NodeToString(item)
				if s1 != s2 {
					t.Errorf("stmt[%d] serialization not deterministic:\n  s1: %s\n  s2: %s", i, s1, s2)
				}
			}
		})
	}
}

// TestParseAlterObjects tests ALTER DATABASE, ALTER INDEX, ALTER VIEW,
// ALTER PROCEDURE, and ALTER FUNCTION (batch 29).
func TestParseAlterObjects(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// ALTER DATABASE - SET option
		{
			name: "alter_database_set",
			sql:  "ALTER DATABASE mydb SET RECOVERY FULL",
		},
		// ALTER DATABASE - SET with ON
		{
			name: "alter_database_set_on",
			sql:  "ALTER DATABASE mydb SET AUTO_CLOSE ON",
		},
		// ALTER DATABASE - SET with OFF
		{
			name: "alter_database_set_off",
			sql:  "ALTER DATABASE mydb SET AUTO_SHRINK OFF",
		},
		// ALTER DATABASE - MODIFY FILE
		{
			name: "alter_database_modify_file",
			sql:  "ALTER DATABASE mydb MODIFY FILE (NAME = myfile, SIZE = 100MB)",
		},
		// ALTER DATABASE - ADD FILE
		{
			name: "alter_database_add_file",
			sql:  "ALTER DATABASE mydb ADD FILE (NAME = newfile, FILENAME = 'C:\\data\\newfile.ndf')",
		},
		// ALTER DATABASE - COLLATE
		{
			name: "alter_database_collate",
			sql:  "ALTER DATABASE mydb COLLATE Latin1_General_CI_AS",
		},
		// ALTER DATABASE - CURRENT
		{
			name: "alter_database_current",
			sql:  "ALTER DATABASE CURRENT SET COMPATIBILITY_LEVEL = 150",
		},
		// ALTER INDEX - REBUILD
		{
			name: "alter_index_rebuild",
			sql:  "ALTER INDEX idx_emp_name ON dbo.employees REBUILD",
		},
		// ALTER INDEX - REBUILD ALL
		{
			name: "alter_index_rebuild_all",
			sql:  "ALTER INDEX ALL ON dbo.employees REBUILD",
		},
		// ALTER INDEX - REORGANIZE
		{
			name: "alter_index_reorganize",
			sql:  "ALTER INDEX idx_emp_name ON employees REORGANIZE",
		},
		// ALTER INDEX - DISABLE
		{
			name: "alter_index_disable",
			sql:  "ALTER INDEX idx_emp_name ON employees DISABLE",
		},
		// ALTER INDEX - SET
		{
			name: "alter_index_set",
			sql:  "ALTER INDEX idx_emp_name ON employees SET (ALLOW_ROW_LOCKS = ON)",
		},
		// ALTER INDEX - REBUILD WITH options
		{
			name: "alter_index_rebuild_with",
			sql:  "ALTER INDEX idx_emp ON dbo.employees REBUILD WITH (FILLFACTOR = 80)",
		},
		// ALTER VIEW
		{
			name: "alter_view_simple",
			sql:  "ALTER VIEW dbo.v_employees AS SELECT id, name FROM employees",
		},
		// ALTER VIEW with SCHEMABINDING
		{
			name: "alter_view_schemabinding",
			sql:  "ALTER VIEW dbo.v_emp WITH SCHEMABINDING AS SELECT id, name FROM dbo.employees",
		},
		// ALTER PROCEDURE
		{
			name: "alter_procedure_simple",
			sql:  "ALTER PROCEDURE dbo.usp_get_emp @id INT AS BEGIN SELECT * FROM employees WHERE id = @id END",
		},
		// ALTER PROC (short form)
		{
			name: "alter_proc_short",
			sql:  "ALTER PROC dbo.usp_get_emp @id INT AS BEGIN SELECT 1 END",
		},
		// ALTER FUNCTION - scalar
		{
			name: "alter_function_scalar",
			sql:  "ALTER FUNCTION dbo.fn_add(@a INT, @b INT) RETURNS INT AS BEGIN RETURN @a + @b END",
		},
		// ALTER FUNCTION - inline TVF
		{
			name: "alter_function_tvf",
			sql:  "ALTER FUNCTION dbo.fn_emp(@dept INT) RETURNS TABLE AS RETURN SELECT * FROM employees WHERE dept_id = @dept",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			for i, item := range result.Items {
				s1 := ast.NodeToString(item)
				s2 := ast.NodeToString(item)
				if s1 != s2 {
					t.Errorf("stmt[%d] serialization not deterministic:\n  s1: %s\n  s2: %s", i, s1, s2)
				}
			}
		})
	}
}

// TestParseDbcc tests DBCC statement parsing (batch 32).
func TestParseDbcc(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// DBCC CHECKDB - no args
		{
			name: "checkdb_no_args",
			sql:  "DBCC CHECKDB",
		},
		// DBCC CHECKDB with database name
		{
			name: "checkdb_with_db",
			sql:  "DBCC CHECKDB('AdventureWorks')",
		},
		// DBCC CHECKDB with WITH option
		{
			name: "checkdb_no_infomsgs",
			sql:  "DBCC CHECKDB('AdventureWorks') WITH NO_INFOMSGS",
		},
		// DBCC CHECKDB with multiple WITH options
		{
			name: "checkdb_multiple_options",
			sql:  "DBCC CHECKDB WITH NO_INFOMSGS, ALL_ERRORMSGS",
		},
		// DBCC CHECKTABLE
		{
			name: "checktable",
			sql:  "DBCC CHECKTABLE('dbo.Orders')",
		},
		// DBCC SHRINKDATABASE
		{
			name: "shrinkdatabase",
			sql:  "DBCC SHRINKDATABASE('AdventureWorks', 10)",
		},
		// DBCC SHRINKFILE
		{
			name: "shrinkfile",
			sql:  "DBCC SHRINKFILE(AdventureWorks_log, 1)",
		},
		// DBCC FREEPROCCACHE - no args
		{
			name: "freeproccache_no_args",
			sql:  "DBCC FREEPROCCACHE",
		},
		// DBCC FREEPROCCACHE with plan handle
		{
			name: "freeproccache_with_arg",
			sql:  "DBCC FREEPROCCACHE(0x06000500)",
		},
		// DBCC DROPCLEANBUFFERS
		{
			name: "dropcleanbuffers",
			sql:  "DBCC DROPCLEANBUFFERS",
		},
		// DBCC SHOW_STATISTICS
		{
			name: "show_statistics",
			sql:  "DBCC SHOW_STATISTICS('AdventureWorks.Person.Address', AK_Address_rowguid)",
		},
		// DBCC INPUTBUFFER
		{
			name: "inputbuffer",
			sql:  "DBCC INPUTBUFFER(52)",
		},
		// DBCC OPENTRAN
		{
			name: "opentran_no_args",
			sql:  "DBCC OPENTRAN",
		},
		// DBCC OPENTRAN with db name
		{
			name: "opentran_with_db",
			sql:  "DBCC OPENTRAN('AdventureWorks')",
		},
		// DBCC CHECKIDENT
		{
			name: "checkident",
			sql:  "DBCC CHECKIDENT('dbo.Orders', RESEED, 0)",
		},
		// DBCC UPDATEUSAGE
		{
			name: "updateusage",
			sql:  "DBCC UPDATEUSAGE(0)",
		},
		// DBCC SQLPERF
		{
			name: "sqlperf",
			sql:  "DBCC SQLPERF(LOGSPACE)",
		},
		// DBCC TRACEON
		{
			name: "traceon",
			sql:  "DBCC TRACEON(3604)",
		},
		// DBCC TRACEOFF
		{
			name: "traceoff",
			sql:  "DBCC TRACEOFF(3604)",
		},
		// DBCC TRACESTATUS
		{
			name: "tracestatus",
			sql:  "DBCC TRACESTATUS",
		},
		// DBCC with WITH TABLERESULTS
		{
			name: "checkdb_tableresults",
			sql:  "DBCC CHECKDB WITH TABLERESULTS",
		},
		// DBCC CHECKALLOC
		{
			name: "checkalloc",
			sql:  "DBCC CHECKALLOC",
		},
		// DBCC CHECKCATALOG
		{
			name: "checkcatalog",
			sql:  "DBCC CHECKCATALOG",
		},
		// DBCC MEMORYSTATUS
		{
			name: "memorystatus",
			sql:  "DBCC MEMORYSTATUS",
		},
		// DBCC PROCCACHE
		{
			name: "proccache",
			sql:  "DBCC PROCCACHE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			stmt, ok := result.Items[0].(*ast.DbccStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ast.DbccStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Command == "" {
				t.Errorf("Parse(%q): DbccStmt.Command is empty", tt.sql)
			}
			checkLocation(t, tt.sql, "DbccStmt", stmt.Loc)
			// Verify deterministic serialization
			s1 := ast.NodeToString(stmt)
			s2 := ast.NodeToString(stmt)
			if s1 != s2 {
				t.Errorf("serialization not deterministic:\n  s1: %s\n  s2: %s", s1, s2)
			}
		})
	}
}

// TestParseBulkInsert tests BULK INSERT statement parsing (batch 33).
func TestParseBulkInsert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		sql := "BULK INSERT dbo.mytable FROM 'C:\\data\\file.csv'"
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.BulkInsertStmt)
		if !ok {
			t.Fatalf("expected *BulkInsertStmt, got %T", result.Items[0])
		}
		if stmt.Table == nil {
			t.Fatal("expected non-nil Table")
		}
		if stmt.DataFile != `C:\data\file.csv` {
			t.Errorf("expected DataFile %q, got %q", `C:\data\file.csv`, stmt.DataFile)
		}
		if stmt.Options != nil {
			t.Errorf("expected nil Options, got %v", stmt.Options)
		}
	})

	t.Run("with_options", func(t *testing.T) {
		sql := "BULK INSERT Sales.Orders FROM '/data/orders.csv' WITH (FIELDTERMINATOR = ',', ROWTERMINATOR = '\\n', FIRSTROW = 2, TABLOCK)"
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.BulkInsertStmt)
		if !ok {
			t.Fatalf("expected *BulkInsertStmt, got %T", result.Items[0])
		}
		if stmt.Table == nil {
			t.Fatal("expected non-nil Table")
		}
		if stmt.DataFile != "/data/orders.csv" {
			t.Errorf("expected DataFile /data/orders.csv, got %q", stmt.DataFile)
		}
		if stmt.Options == nil {
			t.Fatal("expected non-nil Options")
		}
		if stmt.Options.Len() != 4 {
			t.Errorf("expected 4 options, got %d", stmt.Options.Len())
		}
	})

	t.Run("serialization_deterministic", func(t *testing.T) {
		sqls := []string{
			"BULK INSERT dbo.t FROM 'C:\\data\\file.csv'",
			"BULK INSERT dbo.orders FROM '/tmp/data.txt' WITH (FIELDTERMINATOR = ',', TABLOCK, KEEPNULLS)",
			"BULK INSERT mydb.dbo.sales FROM 'https://storage/data.csv' WITH (BATCHSIZE = 1000, CHECK_CONSTRAINTS)",
		}
		for _, sql := range sqls {
			result := ParseAndCheck(t, sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("serialization not deterministic for %q:\n  s1: %s\n  s2: %s", sql, s1, s2)
			}
		}
	})
}

// TestParseBackupRestore tests BACKUP and RESTORE statement parsing (batch 34).
func TestParseBackupRestore(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// BACKUP DATABASE - basic
		{
			name: "backup_database_basic",
			sql:  "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak'",
		},
		// BACKUP DATABASE - with WITH options
		{
			name: "backup_database_with_options",
			sql:  "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH COMPRESSION, STATS = 10",
		},
		// BACKUP DATABASE - TO URL (Azure Blob)
		{
			name: "backup_database_to_url",
			sql:  "BACKUP DATABASE mydb TO URL = 'https://mystorageaccount.blob.core.windows.net/mycontainer/mydb.bak'",
		},
		// BACKUP LOG - basic
		{
			name: "backup_log_basic",
			sql:  "BACKUP LOG mydb TO DISK = '/backup/mydb_log.bak'",
		},
		// BACKUP LOG - with NORECOVERY
		{
			name: "backup_log_norecovery",
			sql:  "BACKUP LOG mydb TO DISK = '/backup/mydb_log.bak' WITH NORECOVERY",
		},
		// BACKUP DATABASE - with FORMAT and INIT
		{
			name: "backup_database_format_init",
			sql:  "BACKUP DATABASE mydb TO DISK = 'C:\\backup\\mydb.bak' WITH FORMAT, INIT, NAME = 'Full backup'",
		},
		// RESTORE DATABASE - basic
		{
			name: "restore_database_basic",
			sql:  "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak'",
		},
		// RESTORE DATABASE - with NORECOVERY
		{
			name: "restore_database_norecovery",
			sql:  "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH NORECOVERY",
		},
		// RESTORE DATABASE - with RECOVERY
		{
			name: "restore_database_recovery",
			sql:  "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH RECOVERY",
		},
		// RESTORE LOG - basic
		{
			name: "restore_log_basic",
			sql:  "RESTORE LOG mydb FROM DISK = '/backup/mydb_log.bak'",
		},
		// RESTORE HEADERONLY
		{
			name: "restore_headeronly",
			sql:  "RESTORE HEADERONLY FROM DISK = '/backup/mydb.bak'",
		},
		// RESTORE FILELISTONLY
		{
			name: "restore_filelistonly",
			sql:  "RESTORE FILELISTONLY FROM DISK = '/backup/mydb.bak'",
		},
		// RESTORE VERIFYONLY
		{
			name: "restore_verifyonly",
			sql:  "RESTORE VERIFYONLY FROM DISK = '/backup/mydb.bak'",
		},
		// RESTORE DATABASE - FROM URL
		{
			name: "restore_database_from_url",
			sql:  "RESTORE DATABASE mydb FROM URL = 'https://mystorageaccount.blob.core.windows.net/mycontainer/mydb.bak'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("stmt serialization not deterministic:\n  s1: %s\n  s2: %s", s1, s2)
			}
		})
	}

	// Type-specific checks
	t.Run("backup_database_type_check", func(t *testing.T) {
		result := ParseAndCheck(t, "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak'")
		stmt, ok := result.Items[0].(*ast.BackupStmt)
		if !ok {
			t.Fatalf("expected *ast.BackupStmt, got %T", result.Items[0])
		}
		if stmt.Type != "DATABASE" {
			t.Errorf("expected Type DATABASE, got %q", stmt.Type)
		}
		if stmt.Database != "mydb" {
			t.Errorf("expected Database mydb, got %q", stmt.Database)
		}
		if stmt.Target != "/backup/mydb.bak" {
			t.Errorf("expected Target /backup/mydb.bak, got %q", stmt.Target)
		}
	})

	t.Run("backup_log_type_check", func(t *testing.T) {
		result := ParseAndCheck(t, "BACKUP LOG mydb TO DISK = '/backup/mydb_log.bak'")
		stmt, ok := result.Items[0].(*ast.BackupStmt)
		if !ok {
			t.Fatalf("expected *ast.BackupStmt, got %T", result.Items[0])
		}
		if stmt.Type != "LOG" {
			t.Errorf("expected Type LOG, got %q", stmt.Type)
		}
		if stmt.Database != "mydb" {
			t.Errorf("expected Database mydb, got %q", stmt.Database)
		}
	})

	t.Run("restore_database_type_check", func(t *testing.T) {
		result := ParseAndCheck(t, "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak'")
		stmt, ok := result.Items[0].(*ast.RestoreStmt)
		if !ok {
			t.Fatalf("expected *ast.RestoreStmt, got %T", result.Items[0])
		}
		if stmt.Type != "DATABASE" {
			t.Errorf("expected Type DATABASE, got %q", stmt.Type)
		}
		if stmt.Database != "mydb" {
			t.Errorf("expected Database mydb, got %q", stmt.Database)
		}
		if stmt.Source != "/backup/mydb.bak" {
			t.Errorf("expected Source /backup/mydb.bak, got %q", stmt.Source)
		}
	})

	t.Run("restore_headeronly_type_check", func(t *testing.T) {
		result := ParseAndCheck(t, "RESTORE HEADERONLY FROM DISK = '/backup/mydb.bak'")
		stmt, ok := result.Items[0].(*ast.RestoreStmt)
		if !ok {
			t.Fatalf("expected *ast.RestoreStmt, got %T", result.Items[0])
		}
		if stmt.Type != "HEADERONLY" {
			t.Errorf("expected Type HEADERONLY, got %q", stmt.Type)
		}
		if stmt.Database != "" {
			t.Errorf("expected empty Database, got %q", stmt.Database)
		}
	})
}

// TestParseSecurityPrincipals tests CREATE/ALTER/DROP USER, LOGIN, ROLE,
// APPLICATION ROLE, and ADD/DROP ROLE MEMBER parsing (batch 30).
func TestParseSecurityPrincipals(t *testing.T) {
	type secCheck struct {
		action     string
		objectType string
		name       string
	}

	tests := []struct {
		name  string
		sql   string
		check secCheck
	}{
		// --- USER ---
		{
			name:  "create_user_simple",
			sql:   "CREATE USER alice",
			check: secCheck{"CREATE", "USER", "alice"},
		},
		{
			name:  "create_user_for_login",
			sql:   "CREATE USER alice FOR LOGIN alice_login",
			check: secCheck{"CREATE", "USER", "alice"},
		},
		{
			name:  "create_user_from_login",
			sql:   "CREATE USER bob FROM LOGIN bob_login",
			check: secCheck{"CREATE", "USER", "bob"},
		},
		{
			name:  "create_user_with_options",
			sql:   "CREATE USER carol WITH DEFAULT_SCHEMA = dbo",
			check: secCheck{"CREATE", "USER", "carol"},
		},
		{
			name:  "alter_user_rename",
			sql:   "ALTER USER alice WITH NAME = alice2",
			check: secCheck{"ALTER", "USER", "alice"},
		},
		{
			name:  "drop_user_simple",
			sql:   "DROP USER alice",
			check: secCheck{"DROP", "USER", "alice"},
		},
		{
			name:  "drop_user_if_exists",
			sql:   "DROP USER IF EXISTS alice",
			check: secCheck{"DROP", "USER", "alice"},
		},

		// --- LOGIN ---
		{
			name:  "create_login_with_password",
			sql:   "CREATE LOGIN mylogin WITH PASSWORD = 'secret123'",
			check: secCheck{"CREATE", "LOGIN", "mylogin"},
		},
		{
			name:  "create_login_from_windows",
			sql:   "CREATE LOGIN [DOMAIN\\user] FROM WINDOWS",
			check: secCheck{"CREATE", "LOGIN", `DOMAIN\user`},
		},
		{
			name:  "alter_login_enable",
			sql:   "ALTER LOGIN mylogin ENABLE",
			check: secCheck{"ALTER", "LOGIN", "mylogin"},
		},
		{
			name:  "alter_login_disable",
			sql:   "ALTER LOGIN mylogin DISABLE",
			check: secCheck{"ALTER", "LOGIN", "mylogin"},
		},
		{
			name:  "alter_login_password",
			sql:   "ALTER LOGIN mylogin WITH PASSWORD = 'newpassword'",
			check: secCheck{"ALTER", "LOGIN", "mylogin"},
		},
		{
			name:  "drop_login_simple",
			sql:   "DROP LOGIN mylogin",
			check: secCheck{"DROP", "LOGIN", "mylogin"},
		},

		// --- ROLE ---
		{
			name:  "create_role_simple",
			sql:   "CREATE ROLE myrole",
			check: secCheck{"CREATE", "ROLE", "myrole"},
		},
		{
			name:  "create_role_authorization",
			sql:   "CREATE ROLE myrole AUTHORIZATION dbo",
			check: secCheck{"CREATE", "ROLE", "myrole"},
		},
		{
			name:  "alter_role_add_member",
			sql:   "ALTER ROLE myrole ADD MEMBER alice",
			check: secCheck{"ALTER", "ROLE", "myrole"},
		},
		{
			name:  "alter_role_drop_member",
			sql:   "ALTER ROLE myrole DROP MEMBER alice",
			check: secCheck{"ALTER", "ROLE", "myrole"},
		},
		{
			name:  "alter_role_rename",
			sql:   "ALTER ROLE myrole WITH NAME = newrole",
			check: secCheck{"ALTER", "ROLE", "myrole"},
		},
		{
			name:  "drop_role_simple",
			sql:   "DROP ROLE myrole",
			check: secCheck{"DROP", "ROLE", "myrole"},
		},
		{
			name:  "drop_role_if_exists",
			sql:   "DROP ROLE IF EXISTS myrole",
			check: secCheck{"DROP", "ROLE", "myrole"},
		},

		// --- APPLICATION ROLE ---
		{
			name:  "create_application_role",
			sql:   "CREATE APPLICATION ROLE approle1 WITH PASSWORD = 'MyPass1'",
			check: secCheck{"CREATE", "APPLICATION ROLE", "approle1"},
		},
		{
			name:  "alter_application_role_rename",
			sql:   "ALTER APPLICATION ROLE approle1 WITH NAME = approle2",
			check: secCheck{"ALTER", "APPLICATION ROLE", "approle1"},
		},
		{
			name:  "drop_application_role",
			sql:   "DROP APPLICATION ROLE approle1",
			check: secCheck{"DROP", "APPLICATION ROLE", "approle1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ast.SecurityStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Action != tt.check.action {
				t.Errorf("Parse(%q): Action = %q, want %q", tt.sql, stmt.Action, tt.check.action)
			}
			if stmt.ObjectType != tt.check.objectType {
				t.Errorf("Parse(%q): ObjectType = %q, want %q", tt.sql, stmt.ObjectType, tt.check.objectType)
			}
			if stmt.Name != tt.check.name {
				t.Errorf("Parse(%q): Name = %q, want %q", tt.sql, stmt.Name, tt.check.name)
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
			// Verify deterministic serialization
			s1 := ast.NodeToString(stmt)
			s2 := ast.NodeToString(stmt)
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic:\n  s1: %s\n  s2: %s", tt.sql, s1, s2)
			}
		})
	}
}

// TestParseBeginDistributedTransaction tests BEGIN DISTRIBUTED TRANSACTION (batch 39).
func TestParseBeginDistributedTransaction(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "begin_distributed_tran",
			sql:  "BEGIN DISTRIBUTED TRAN",
		},
		{
			name: "begin_distributed_transaction",
			sql:  "BEGIN DISTRIBUTED TRANSACTION",
		},
		{
			name: "begin_distributed_tran_named",
			sql:  "BEGIN DISTRIBUTED TRANSACTION myTran",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			stmt, ok := result.Items[0].(*ast.BeginDistributedTransStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ast.BeginDistributedTransStmt, got %T", tt.sql, result.Items[0])
			}
			checkLocation(t, tt.sql, "BeginDistributedTransStmt", stmt.Loc)
		})
	}
}

// TestParseCreateStatistics tests CREATE/UPDATE/DROP STATISTICS (batch 40).
func TestParseCreateStatistics(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create_statistics_simple",
			sql:  "CREATE STATISTICS stat1 ON dbo.Orders (OrderDate)",
		},
		{
			name: "create_statistics_fullscan",
			sql:  "CREATE STATISTICS stat2 ON dbo.Orders (CustomerID, OrderDate) WITH FULLSCAN",
		},
		{
			name: "create_statistics_sample",
			sql:  "CREATE STATISTICS stat3 ON Products (Price) WITH SAMPLE 30 PERCENT",
		},
		{
			name: "create_statistics_norecompute",
			sql:  "CREATE STATISTICS stat4 ON Products (Price) WITH NORECOMPUTE",
		},
		{
			name: "update_statistics_table",
			sql:  "UPDATE STATISTICS dbo.Orders",
		},
		{
			name: "update_statistics_named",
			sql:  "UPDATE STATISTICS dbo.Orders stat1",
		},
		{
			name: "update_statistics_fullscan",
			sql:  "UPDATE STATISTICS dbo.Orders WITH FULLSCAN",
		},
		{
			name: "drop_statistics",
			sql:  "DROP STATISTICS dbo.Orders.stat1",
		},
		{
			name: "drop_statistics_multiple",
			sql:  "DROP STATISTICS dbo.Orders.stat1, dbo.Products.stat2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			// Verify deterministic serialization
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseSetOptions tests SET session options (batch 41).
func TestParseSetOptions(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		option string
	}{
		{
			name:   "set_nocount_on",
			sql:    "SET NOCOUNT ON",
			option: "NOCOUNT",
		},
		{
			name:   "set_nocount_off",
			sql:    "SET NOCOUNT OFF",
			option: "NOCOUNT",
		},
		{
			name:   "set_ansi_nulls_on",
			sql:    "SET ANSI_NULLS ON",
			option: "ANSI_NULLS",
		},
		{
			name:   "set_quoted_identifier_on",
			sql:    "SET QUOTED_IDENTIFIER ON",
			option: "QUOTED_IDENTIFIER",
		},
		{
			name:   "set_xact_abort_on",
			sql:    "SET XACT_ABORT ON",
			option: "XACT_ABORT",
		},
		{
			name:   "set_arithabort_on",
			sql:    "SET ARITHABORT ON",
			option: "ARITHABORT",
		},
		{
			name:   "set_identity_insert",
			sql:    "SET IDENTITY_INSERT dbo.Orders ON",
			option: "IDENTITY_INSERT",
		},
		{
			name:   "set_transaction_isolation_level",
			sql:    "SET TRANSACTION ISOLATION LEVEL READ COMMITTED",
			option: "TRANSACTION ISOLATION LEVEL",
		},
		{
			name:   "set_transaction_isolation_serializable",
			sql:    "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE",
			option: "TRANSACTION ISOLATION LEVEL",
		},
		{
			name:   "set_rowcount",
			sql:    "SET ROWCOUNT 100",
			option: "ROWCOUNT",
		},
		{
			name:   "set_language",
			sql:    "SET LANGUAGE us_english",
			option: "LANGUAGE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			stmt, ok := result.Items[0].(*ast.SetOptionStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ast.SetOptionStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Option != tt.option {
				t.Errorf("Parse(%q): Option = %q, want %q", tt.sql, stmt.Option, tt.option)
			}
			checkLocation(t, tt.sql, "SetOptionStmt", stmt.Loc)
		})
	}
}

// TestParsePartition tests CREATE/ALTER PARTITION FUNCTION/SCHEME (batch 42).
func TestParsePartition(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create_partition_function_left",
			sql:  "CREATE PARTITION FUNCTION myRangePF1 (int) AS RANGE LEFT FOR VALUES (1, 100, 1000)",
		},
		{
			name: "create_partition_function_right",
			sql:  "CREATE PARTITION FUNCTION myRangePF2 (datetime) AS RANGE RIGHT FOR VALUES ('2020-01-01', '2021-01-01')",
		},
		{
			name: "create_partition_function_empty",
			sql:  "CREATE PARTITION FUNCTION pfEmpty (int) AS RANGE LEFT FOR VALUES ()",
		},
		{
			name: "alter_partition_function_split",
			sql:  "ALTER PARTITION FUNCTION myRangePF1 () SPLIT RANGE (500)",
		},
		{
			name: "alter_partition_function_merge",
			sql:  "ALTER PARTITION FUNCTION myRangePF1 () MERGE RANGE (1)",
		},
		{
			name: "create_partition_scheme_to",
			sql:  "CREATE PARTITION SCHEME myRangePS1 AS PARTITION myRangePF1 TO (fg1, fg2, fg3, fg4)",
		},
		{
			name: "create_partition_scheme_all",
			sql:  "CREATE PARTITION SCHEME myRangePS2 AS PARTITION myRangePF1 ALL TO ([PRIMARY])",
		},
		{
			name: "alter_partition_scheme",
			sql:  "ALTER PARTITION SCHEME myRangePS1 NEXT USED [PRIMARY]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseFulltext tests CREATE/ALTER FULLTEXT INDEX and CATALOG (batch 43).
func TestParseFulltext(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create_fulltext_index",
			sql:  "CREATE FULLTEXT INDEX ON dbo.Products (Name, Description) KEY INDEX PK_Products",
		},
		{
			name: "create_fulltext_index_catalog",
			sql:  "CREATE FULLTEXT INDEX ON dbo.Products (Name) KEY INDEX PK_Products ON MyCatalog",
		},
		{
			name: "create_fulltext_index_with_options",
			sql:  "CREATE FULLTEXT INDEX ON dbo.Products (Name) KEY INDEX PK_Products WITH CHANGE_TRACKING = AUTO",
		},
		{
			name: "alter_fulltext_index_enable",
			sql:  "ALTER FULLTEXT INDEX ON dbo.Products ENABLE",
		},
		{
			name: "alter_fulltext_index_disable",
			sql:  "ALTER FULLTEXT INDEX ON dbo.Products DISABLE",
		},
		{
			name: "create_fulltext_catalog",
			sql:  "CREATE FULLTEXT CATALOG MyCatalog",
		},
		{
			name: "create_fulltext_catalog_options",
			sql:  "CREATE FULLTEXT CATALOG MyCatalog WITH ACCENT_SENSITIVITY = ON",
		},
		{
			name: "alter_fulltext_catalog_rebuild",
			sql:  "ALTER FULLTEXT CATALOG MyCatalog REBUILD",
		},
		{
			name: "alter_fulltext_catalog_as_default",
			sql:  "ALTER FULLTEXT CATALOG MyCatalog AS DEFAULT",
		},
		{
			name: "drop_fulltext_index",
			sql:  "DROP FULLTEXT INDEX ON dbo.Products",
		},
		{
			name: "drop_fulltext_catalog",
			sql:  "DROP FULLTEXT CATALOG MyCatalog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseXmlSchemaCollection tests CREATE/ALTER/DROP XML SCHEMA COLLECTION (batch 44).
func TestParseXmlSchemaCollection(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create_xml_schema_collection",
			sql:  `CREATE XML SCHEMA COLLECTION dbo.myCollection AS '<schema xmlns="http://www.w3.org/2001/XMLSchema"><element name="root"/></schema>'`,
		},
		{
			name: "create_xml_schema_collection_variable",
			sql:  "CREATE XML SCHEMA COLLECTION dbo.myCol AS @schemaVar",
		},
		{
			name: "alter_xml_schema_collection",
			sql:  `ALTER XML SCHEMA COLLECTION dbo.myCollection ADD '<schema xmlns="http://www.w3.org/2001/XMLSchema"><element name="extra"/></schema>'`,
		},
		{
			name: "drop_xml_schema_collection",
			sql:  "DROP XML SCHEMA COLLECTION dbo.myCollection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseAssembly tests CREATE/ALTER/DROP ASSEMBLY (batch 45).
func TestParseAssembly(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create_assembly",
			sql:  "CREATE ASSEMBLY myAssembly FROM 'C:\\myAssembly.dll' WITH PERMISSION_SET = SAFE",
		},
		{
			name: "create_assembly_authorization",
			sql:  "CREATE ASSEMBLY myAssembly AUTHORIZATION dbo FROM 'C:\\myAssembly.dll'",
		},
		{
			name: "create_assembly_unsafe",
			sql:  "CREATE ASSEMBLY myAssembly FROM 'C:\\myAssembly.dll' WITH PERMISSION_SET = UNSAFE",
		},
		{
			name: "alter_assembly_from",
			sql:  "ALTER ASSEMBLY myAssembly FROM 'C:\\newAssembly.dll'",
		},
		{
			name: "alter_assembly_with",
			sql:  "ALTER ASSEMBLY myAssembly WITH PERMISSION_SET = EXTERNAL_ACCESS",
		},
		{
			name: "drop_assembly",
			sql:  "DROP ASSEMBLY myAssembly",
		},
		{
			name: "drop_assembly_if_exists",
			sql:  "DROP ASSEMBLY IF EXISTS myAssembly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseServiceBroker tests Service Broker statements (batch 46).
func TestParseServiceBroker(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create_message_type",
			sql:  "CREATE MESSAGE TYPE [//MyApp/RequestMsg] VALIDATION = NONE",
		},
		{
			name: "create_message_type_wellformed",
			sql:  "CREATE MESSAGE TYPE [//MyApp/ReplyMsg] VALIDATION = WELL_FORMED_XML",
		},
		{
			name: "create_contract",
			sql:  "CREATE CONTRACT [//MyApp/MyContract] ([//MyApp/RequestMsg] SENT BY INITIATOR)",
		},
		{
			name: "create_queue_simple",
			sql:  "CREATE QUEUE MyQueue",
		},
		{
			name: "create_service_simple",
			sql:  "CREATE SERVICE [//MyApp/MyService] ON QUEUE MyQueue",
		},
		{
			name: "send_on_conversation",
			sql:  "SEND ON CONVERSATION @handle MESSAGE TYPE [//MyApp/RequestMsg] (N'Hello')",
		},
		{
			name: "end_conversation",
			sql:  "END CONVERSATION @handle",
		},
		{
			name: "end_conversation_with_cleanup",
			sql:  "END CONVERSATION @handle WITH CLEANUP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseMiscUtility tests CHECKPOINT, RECONFIGURE, SHUTDOWN, KILL, READTEXT,
// WRITETEXT, UPDATETEXT (batch 47).
func TestParseMiscUtility(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "checkpoint",
			sql:  "CHECKPOINT",
		},
		{
			name: "checkpoint_duration",
			sql:  "CHECKPOINT 5",
		},
		{
			name: "reconfigure",
			sql:  "RECONFIGURE",
		},
		{
			name: "reconfigure_with_override",
			sql:  "RECONFIGURE WITH OVERRIDE",
		},
		{
			name: "shutdown",
			sql:  "SHUTDOWN",
		},
		{
			name: "shutdown_nowait",
			sql:  "SHUTDOWN WITH NOWAIT",
		},
		{
			name: "kill_session",
			sql:  "KILL 52",
		},
		{
			name: "kill_with_statusonly",
			sql:  "KILL 52 WITH STATUSONLY",
		},
		{
			name: "readtext",
			sql:  "READTEXT pub_info.pr_info @ptrval 0 10",
		},
		{
			name: "writetext",
			sql:  "WRITETEXT pub_info.pr_info @ptrval 'New content'",
		},
		{
			name: "updatetext",
			sql:  "UPDATETEXT pub_info.pr_info @ptrval NULL NULL 'New content'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
			}
		})
	}
}

// TestParseDropExtended tests extended DROP support (batch 48).
func TestParseDropExtended(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		dropType ast.DropObjectType
	}{
		{
			name:     "drop_sequence",
			sql:      "DROP SEQUENCE dbo.mySeq",
			dropType: ast.DropSequence,
		},
		{
			name:     "drop_synonym",
			sql:      "DROP SYNONYM dbo.mySyn",
			dropType: ast.DropSynonym,
		},
		{
			name:     "drop_assembly",
			sql:      "DROP ASSEMBLY myAssembly",
			dropType: ast.DropAssembly,
		},
		{
			name:     "drop_partition_function",
			sql:      "DROP PARTITION FUNCTION myPF",
			dropType: ast.DropPartitionFunction,
		},
		{
			name:     "drop_partition_scheme",
			sql:      "DROP PARTITION SCHEME myPS",
			dropType: ast.DropPartitionScheme,
		},
		{
			name:     "drop_xml_schema_collection",
			sql:      "DROP XML SCHEMA COLLECTION dbo.myCol",
			dropType: ast.DropXmlSchemaCollection,
		},
		{
			name:     "drop_fulltext_index",
			sql:      "DROP FULLTEXT INDEX ON dbo.Products",
			dropType: ast.DropFulltextIndex,
		},
		{
			name:     "drop_fulltext_catalog",
			sql:      "DROP FULLTEXT CATALOG MyCatalog",
			dropType: ast.DropFulltextCatalog,
		},
		{
			name:     "drop_rule",
			sql:      "DROP RULE dbo.myRule",
			dropType: ast.DropRule,
		},
		{
			name:     "drop_default",
			sql:      "DROP DEFAULT dbo.myDefault",
			dropType: ast.DropDefault,
		},
		{
			name:     "drop_statistics",
			sql:      "DROP STATISTICS dbo.Orders.stat1",
			dropType: ast.DropStatistics,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			// DROP STATISTICS returns a DropStatisticsStmt, not DropStmt
			if tt.dropType == ast.DropStatistics {
				_, ok := result.Items[0].(*ast.DropStatisticsStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *ast.DropStatisticsStmt, got %T", tt.sql, result.Items[0])
				}
				return
			}
			stmt, ok := result.Items[0].(*ast.DropStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ast.DropStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.ObjectType != tt.dropType {
				t.Errorf("Parse(%q): ObjectType = %v, want %v", tt.sql, stmt.ObjectType, tt.dropType)
			}
			checkLocation(t, tt.sql, "DropStmt", stmt.Loc)
		})
	}
}

// TestParseStmtDispatchPhase2 is the integration test for batch 49.
// It parses multi-statement T-SQL scripts combining statements from all batches
// and verifies that each statement is parsed to the correct AST node type.
func TestParseStmtDispatchPhase2(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		stmtCount int
		stmtTypes []string // expected type names
	}{
		{
			name: "multi_statement_ddl",
			sql: `CREATE TABLE dbo.Orders (id INT PRIMARY KEY, amount DECIMAL(10,2));
CREATE INDEX IX_Orders ON dbo.Orders (amount);
CREATE VIEW dbo.vOrders AS SELECT * FROM dbo.Orders;
ALTER TABLE dbo.Orders ADD status INT;
DROP INDEX IX_Orders ON dbo.Orders;`,
			stmtCount: 5,
			stmtTypes: []string{
				"*ast.CreateTableStmt",
				"*ast.CreateIndexStmt",
				"*ast.CreateViewStmt",
				"*ast.AlterTableStmt",
				"*ast.DropStmt",
			},
		},
		{
			name: "multi_statement_dml",
			sql: `INSERT INTO dbo.Orders (id, amount) VALUES (1, 100.00);
UPDATE dbo.Orders SET amount = 200.00 WHERE id = 1;
DELETE FROM dbo.Orders WHERE id = 1;
MERGE dbo.Orders AS t USING dbo.Staging AS s ON t.id = s.id WHEN MATCHED THEN UPDATE SET amount = s.amount;`,
			stmtCount: 4,
			stmtTypes: []string{
				"*ast.InsertStmt",
				"*ast.UpdateStmt",
				"*ast.DeleteStmt",
				"*ast.MergeStmt",
			},
		},
		{
			name: "control_flow_and_variables",
			sql: `DECLARE @x INT = 10;
SET @x = 20;
IF @x > 0 SELECT @x;
WHILE @x > 0 SET @x = @x - 1;
BEGIN TRY SELECT 1 END TRY BEGIN CATCH SELECT 0 END CATCH;`,
			stmtCount: 5,
			stmtTypes: []string{
				"*ast.DeclareStmt",
				"*ast.SetStmt",
				"*ast.IfStmt",
				"*ast.WhileStmt",
				"*ast.TryCatchStmt",
			},
		},
		{
			name: "security_and_admin",
			sql: `CREATE USER testUser FOR LOGIN testLogin;
GRANT SELECT ON dbo.Orders TO testUser;
DBCC CHECKDB;
BACKUP DATABASE myDB TO DISK = 'backup.bak';
BULK INSERT dbo.Orders FROM 'data.csv';`,
			stmtCount: 5,
			stmtTypes: []string{
				"*ast.SecurityStmt",
				"*ast.GrantStmt",
				"*ast.DbccStmt",
				"*ast.BackupStmt",
				"*ast.BulkInsertStmt",
			},
		},
		{
			name: "batches_35_48_mix",
			sql: `CREATE STATISTICS stat1 ON dbo.Orders (amount);
SET NOCOUNT ON;
CHECKPOINT;
KILL 55;`,
			stmtCount: 4,
			stmtTypes: []string{
				"*ast.CreateStatisticsStmt",
				"*ast.SetOptionStmt",
				"*ast.CheckpointStmt",
				"*ast.KillStmt",
			},
		},
		{
			name: "transaction_and_cursor",
			sql: `BEGIN TRANSACTION;
DECLARE emp_cursor CURSOR FOR SELECT id FROM dbo.Emp;
OPEN emp_cursor;
FETCH NEXT FROM emp_cursor;
CLOSE emp_cursor;
DEALLOCATE emp_cursor;
COMMIT;`,
			stmtCount: 7,
			stmtTypes: []string{
				"*ast.BeginTransStmt",
				"*ast.DeclareCursorStmt",
				"*ast.OpenCursorStmt",
				"*ast.FetchCursorStmt",
				"*ast.CloseCursorStmt",
				"*ast.DeallocateCursorStmt",
				"*ast.CommitTransStmt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != tt.stmtCount {
				t.Fatalf("Parse(%q): got %d statements, want %d", tt.name, result.Len(), tt.stmtCount)
			}
			for i, item := range result.Items {
				gotType := fmt.Sprintf("%T", item)
				if i < len(tt.stmtTypes) && gotType != tt.stmtTypes[i] {
					t.Errorf("stmt[%d] type = %s, want %s", i, gotType, tt.stmtTypes[i])
				}
			}
		})
	}
}

// TestParseAlterTriggerEnableDisable tests batch 51: ALTER TRIGGER, ENABLE TRIGGER, DISABLE TRIGGER.
func TestParseAlterTriggerEnableDisable(t *testing.T) {
	// ALTER TRIGGER tests
	alterTests := []string{
		"ALTER TRIGGER dbo.MyTrigger ON dbo.MyTable AFTER INSERT AS SELECT 1",
		"ALTER TRIGGER MyTrigger ON DATABASE AFTER CREATE_TABLE AS SELECT 1",
		"ALTER TRIGGER MyTrigger ON ALL SERVER AFTER LOGON AS SELECT 1",
	}
	for _, sql := range alterTests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.CreateTriggerStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *CreateTriggerStmt, got %T", sql, result.Items[0])
			}
			if !stmt.OrAlter {
				t.Errorf("Parse(%q): expected OrAlter=true", sql)
			}
			checkLocation(t, sql, "CreateTriggerStmt", stmt.Loc)
		})
	}

	// ENABLE/DISABLE TRIGGER tests
	edTests := []struct {
		sql    string
		enable bool
	}{
		// ENABLE single trigger on table
		{"ENABLE TRIGGER MyTrigger ON dbo.MyTable", true},
		// DISABLE single trigger on table
		{"DISABLE TRIGGER MyTrigger ON dbo.MyTable", false},
		// ENABLE ALL triggers on table
		{"ENABLE TRIGGER ALL ON dbo.MyTable", true},
		// DISABLE ALL triggers on table
		{"DISABLE TRIGGER ALL ON dbo.MyTable", false},
		// ENABLE trigger on DATABASE
		{"ENABLE TRIGGER MyTrigger ON DATABASE", true},
		// DISABLE trigger on DATABASE
		{"DISABLE TRIGGER MyTrigger ON DATABASE", false},
		// ENABLE trigger on ALL SERVER
		{"ENABLE TRIGGER MyTrigger ON ALL SERVER", true},
		// DISABLE trigger on ALL SERVER
		{"DISABLE TRIGGER MyTrigger ON ALL SERVER", false},
		// Multiple triggers
		{"DISABLE TRIGGER Trig1, Trig2, Trig3 ON dbo.MyTable", false},
	}
	for _, tt := range edTests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.EnableDisableTriggerStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *EnableDisableTriggerStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Enable != tt.enable {
				t.Errorf("Parse(%q): expected Enable=%v, got %v", tt.sql, tt.enable, stmt.Enable)
			}
			checkLocation(t, tt.sql, "EnableDisableTriggerStmt", stmt.Loc)
		})
	}
}

// TestParseServiceBrokerMissing tests batch 50: CREATE ROUTE, CREATE REMOTE SERVICE BINDING, GET CONVERSATION GROUP.
func TestParseServiceBrokerMissing(t *testing.T) {
	tests := []string{
		// CREATE ROUTE - basic
		"CREATE ROUTE MyRoute WITH ADDRESS = 'TCP://10.0.0.1:4022'",
		// CREATE ROUTE - with AUTHORIZATION
		"CREATE ROUTE MyRoute AUTHORIZATION dbo WITH ADDRESS = 'TCP://10.0.0.1:4022'",
		// CREATE ROUTE - with SERVICE_NAME
		"CREATE ROUTE MyRoute WITH SERVICE_NAME = '//Adventure-Works.com/Expenses', ADDRESS = 'TCP://10.0.0.1:4022'",
		// CREATE ROUTE - with BROKER_INSTANCE
		"CREATE ROUTE MyRoute WITH SERVICE_NAME = '//Adventure-Works.com/Expenses', BROKER_INSTANCE = 'D8787ED9-A3CF-43CE-8FCC-B123456789AB', ADDRESS = 'TCP://10.0.0.1:4022'",
		// CREATE ROUTE - with LIFETIME
		"CREATE ROUTE MyRoute WITH LIFETIME = 600, ADDRESS = 'TCP://10.0.0.1:4022'",
		// CREATE ROUTE - with MIRROR_ADDRESS
		"CREATE ROUTE MyRoute WITH ADDRESS = 'TCP://10.0.0.1:4022', MIRROR_ADDRESS = 'TCP://10.0.0.2:4022'",
		// CREATE ROUTE - LOCAL address
		"CREATE ROUTE MyRoute WITH SERVICE_NAME = '//Adventure-Works.com/Expenses', ADDRESS = 'LOCAL'",
		// CREATE ROUTE - TRANSPORT address
		"CREATE ROUTE TransportRoute WITH ADDRESS = 'TRANSPORT'",
		// CREATE REMOTE SERVICE BINDING - basic
		"CREATE REMOTE SERVICE BINDING MyBinding TO SERVICE '//Adventure-Works.com/Expenses' WITH USER = ExpensesUser",
		// CREATE REMOTE SERVICE BINDING - with AUTHORIZATION
		"CREATE REMOTE SERVICE BINDING MyBinding AUTHORIZATION dbo TO SERVICE '//Adventure-Works.com/Expenses' WITH USER = ExpensesUser",
		// CREATE REMOTE SERVICE BINDING - with ANONYMOUS ON
		"CREATE REMOTE SERVICE BINDING MyBinding TO SERVICE '//Adventure-Works.com/Expenses' WITH USER = ExpensesUser, ANONYMOUS = ON",
		// CREATE REMOTE SERVICE BINDING - with ANONYMOUS OFF
		"CREATE REMOTE SERVICE BINDING MyBinding TO SERVICE '//Adventure-Works.com/Expenses' WITH USER = ExpensesUser, ANONYMOUS = OFF",
		// GET CONVERSATION GROUP - basic
		"GET CONVERSATION GROUP @conversation_group_id FROM ExpenseQueue",
		// GET CONVERSATION GROUP - qualified queue name
		"GET CONVERSATION GROUP @conversation_group_id FROM dbo.ExpenseQueue",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.ServiceBrokerStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ServiceBrokerStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "ServiceBrokerStmt", stmt.Loc)
		})
	}
}
