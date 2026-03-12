package parser

import (
	"fmt"
	"strings"
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

// TestParseEndpoint tests batch 55: CREATE/ALTER/DROP ENDPOINT.
func TestParseEndpoint(t *testing.T) {
	tests := []string{
		// CREATE ENDPOINT basic TCP with SERVICE_BROKER
		"CREATE ENDPOINT MyEndpoint STATE = STARTED AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = WINDOWS)",
		// CREATE ENDPOINT with AUTHORIZATION
		"CREATE ENDPOINT MyEndpoint AUTHORIZATION sa AS TCP (LISTENER_PORT = 4022) FOR DATABASE_MIRRORING (ROLE = PARTNER)",
		// CREATE ENDPOINT stopped
		"CREATE ENDPOINT MyEndpoint STATE = STOPPED AS TCP (LISTENER_PORT = 5022) FOR TSQL ()",
		// ALTER ENDPOINT
		"ALTER ENDPOINT MyEndpoint STATE = STARTED",
		// ALTER ENDPOINT with options
		"ALTER ENDPOINT MyEndpoint AS TCP (LISTENER_PORT = 5023)",
		// DROP ENDPOINT
		"DROP ENDPOINT MyEndpoint",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseServerAudit tests batch 54: SERVER AUDIT, SERVER AUDIT SPECIFICATION, DATABASE AUDIT SPECIFICATION.
func TestParseServerAudit(t *testing.T) {
	tests := []string{
		// CREATE SERVER AUDIT to file
		"CREATE SERVER AUDIT MyAudit TO FILE (FILEPATH = 'C:\\Audits\\', MAXSIZE = 100 MB)",
		// CREATE SERVER AUDIT to application log
		"CREATE SERVER AUDIT MyAudit TO APPLICATION_LOG",
		// CREATE SERVER AUDIT to security log
		"CREATE SERVER AUDIT MyAudit TO SECURITY_LOG",
		// CREATE SERVER AUDIT with options
		"CREATE SERVER AUDIT MyAudit TO FILE (FILEPATH = 'C:\\Audits\\') WITH (QUEUE_DELAY = 1000, ON_FAILURE = CONTINUE)",
		// ALTER SERVER AUDIT
		"ALTER SERVER AUDIT MyAudit WITH (STATE = ON)",
		// ALTER SERVER AUDIT modify name
		"ALTER SERVER AUDIT MyAudit MODIFY NAME = NewAuditName",
		// DROP SERVER AUDIT
		"DROP SERVER AUDIT MyAudit",
		// CREATE SERVER AUDIT SPECIFICATION
		"CREATE SERVER AUDIT SPECIFICATION MySpec FOR SERVER AUDIT MyAudit ADD (FAILED_LOGIN_GROUP)",
		// CREATE SERVER AUDIT SPECIFICATION with state
		"CREATE SERVER AUDIT SPECIFICATION MySpec FOR SERVER AUDIT MyAudit ADD (FAILED_LOGIN_GROUP) WITH (STATE = ON)",
		// ALTER SERVER AUDIT SPECIFICATION
		"ALTER SERVER AUDIT SPECIFICATION MySpec FOR SERVER AUDIT MyAudit ADD (SUCCESSFUL_LOGIN_GROUP) WITH (STATE = ON)",
		// DROP SERVER AUDIT SPECIFICATION
		"DROP SERVER AUDIT SPECIFICATION MySpec",
		// CREATE DATABASE AUDIT SPECIFICATION
		"CREATE DATABASE AUDIT SPECIFICATION MyDbSpec FOR SERVER AUDIT MyAudit ADD (SELECT ON OBJECT::dbo.MyTable BY public)",
		// CREATE DATABASE AUDIT SPECIFICATION with state
		"CREATE DATABASE AUDIT SPECIFICATION MyDbSpec FOR SERVER AUDIT MyAudit ADD (SELECT ON SCHEMA::dbo BY public) WITH (STATE = ON)",
		// ALTER DATABASE AUDIT SPECIFICATION
		"ALTER DATABASE AUDIT SPECIFICATION MyDbSpec FOR SERVER AUDIT MyAudit DROP (SELECT ON OBJECT::dbo.MyTable BY public) WITH (STATE = OFF)",
		// DROP DATABASE AUDIT SPECIFICATION
		"DROP DATABASE AUDIT SPECIFICATION MyDbSpec",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseSecurityContext tests batch 53: EXECUTE AS, REVERT, ALTER AUTHORIZATION.
func TestParseSecurityContext(t *testing.T) {
	tests := []string{
		// EXECUTE AS LOGIN
		"EXECUTE AS LOGIN = 'TestLogin'",
		// EXECUTE AS USER
		"EXECUTE AS USER = 'TestUser'",
		// EXECUTE AS CALLER
		"EXECUTE AS CALLER",
		// EXECUTE AS SELF
		"EXECUTE AS SELF",
		// EXECUTE AS OWNER
		"EXECUTE AS OWNER",
		// EXECUTE AS with NO REVERT
		"EXECUTE AS USER = 'TestUser' WITH NO REVERT",
		// EXECUTE AS with COOKIE
		"EXECUTE AS USER = 'TestUser' WITH COOKIE INTO @cookie",
		// REVERT basic
		"REVERT",
		// REVERT with COOKIE
		"REVERT WITH COOKIE = @cookie",
		// ALTER AUTHORIZATION basic
		"ALTER AUTHORIZATION ON dbo.MyTable TO NewOwner",
		// ALTER AUTHORIZATION with class type
		"ALTER AUTHORIZATION ON OBJECT::dbo.MyTable TO NewOwner",
		// ALTER AUTHORIZATION to SCHEMA OWNER
		"ALTER AUTHORIZATION ON SCHEMA::Sales TO SCHEMA OWNER",
		// ALTER AUTHORIZATION on DATABASE
		"ALTER AUTHORIZATION ON DATABASE::MyDB TO sa",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseServiceBrokerAlterDrop tests batch 52: ALTER/DROP Service Broker objects, BROKER PRIORITY, MOVE CONVERSATION.
func TestParseServiceBrokerAlterDrop(t *testing.T) {
	tests := []string{
		// ALTER QUEUE
		"ALTER QUEUE dbo.ExpenseQueue WITH STATUS = ON",
		"ALTER QUEUE ExpenseQueue WITH ACTIVATION (STATUS = ON, PROCEDURE_NAME = dbo.expense_procedure, MAX_QUEUE_READERS = 5, EXECUTE AS SELF)",
		// ALTER SERVICE
		"ALTER SERVICE MyService ON QUEUE dbo.NewQueue",
		// ALTER ROUTE
		"ALTER ROUTE MyRoute WITH ADDRESS = 'TCP://10.0.0.2:4022'",
		"ALTER ROUTE MyRoute WITH SERVICE_NAME = '//example.com/svc', ADDRESS = 'TCP://10.0.0.2:4022'",
		// ALTER REMOTE SERVICE BINDING
		"ALTER REMOTE SERVICE BINDING MyBinding WITH USER = NewUser",
		"ALTER REMOTE SERVICE BINDING MyBinding WITH USER = NewUser, ANONYMOUS = ON",
		// ALTER MESSAGE TYPE
		"ALTER MESSAGE TYPE MyMessageType",
		// ALTER CONTRACT
		"ALTER CONTRACT MyContract",
		// DROP MESSAGE TYPE
		"DROP MESSAGE TYPE MyMessageType",
		// DROP CONTRACT
		"DROP CONTRACT MyContract",
		// DROP QUEUE
		"DROP QUEUE dbo.ExpenseQueue",
		// DROP SERVICE
		"DROP SERVICE MyService",
		// DROP ROUTE
		"DROP ROUTE MyRoute",
		// DROP REMOTE SERVICE BINDING
		"DROP REMOTE SERVICE BINDING MyBinding",
		// DROP BROKER PRIORITY
		"DROP BROKER PRIORITY MyPriority",
		// CREATE BROKER PRIORITY
		"CREATE BROKER PRIORITY MyPriority FOR CONVERSATION SET (CONTRACT_NAME = MyContract, LOCAL_SERVICE_NAME = MyService, REMOTE_SERVICE_NAME = 'RemoteService', PRIORITY_LEVEL = 5)",
		"CREATE BROKER PRIORITY SimplePriority FOR CONVERSATION",
		"CREATE BROKER PRIORITY AnyPriority FOR CONVERSATION SET (CONTRACT_NAME = ANY, PRIORITY_LEVEL = DEFAULT)",
		// ALTER BROKER PRIORITY
		"ALTER BROKER PRIORITY MyPriority FOR CONVERSATION SET (PRIORITY_LEVEL = 10)",
		// MOVE CONVERSATION
		"MOVE CONVERSATION @dialog_handle TO @new_group_id",
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

// TestParseEventNotification tests batch 56: CREATE/DROP EVENT NOTIFICATION.
func TestParseEventNotification(t *testing.T) {
	tests := []string{
		// CREATE EVENT NOTIFICATION on SERVER
		"CREATE EVENT NOTIFICATION log_ddl1 ON SERVER FOR Object_Created TO SERVICE 'NotifyService', '8140a771-3c4b-4479-8ac0-81008ab17984'",
		// CREATE EVENT NOTIFICATION on DATABASE
		"CREATE EVENT NOTIFICATION Notify_ALTER_T1 ON DATABASE FOR ALTER_TABLE TO SERVICE 'NotifyService', '8140a771-3c4b-4479-8ac0-81008ab17984'",
		// CREATE EVENT NOTIFICATION on QUEUE
		"CREATE EVENT NOTIFICATION NotifyQueue ON QUEUE dbo.ExpenseQueue FOR QUEUE_ACTIVATION TO SERVICE 'NotifyService', '8140a771-3c4b-4479-8ac0-81008ab17984'",
		// CREATE EVENT NOTIFICATION with FAN_IN
		"CREATE EVENT NOTIFICATION log_ddl2 ON SERVER WITH FAN_IN FOR ALTER_TABLE TO SERVICE 'NotifyService', '8140a771-3c4b-4479-8ac0-81008ab17984'",
		// CREATE EVENT NOTIFICATION with multiple events
		"CREATE EVENT NOTIFICATION log_ddl3 ON DATABASE FOR CREATE_TABLE, ALTER_TABLE, DROP_TABLE TO SERVICE 'NotifyService', 'current database'",
		// DROP EVENT NOTIFICATION on SERVER
		"DROP EVENT NOTIFICATION log_ddl1 ON SERVER",
		// DROP EVENT NOTIFICATION on DATABASE
		"DROP EVENT NOTIFICATION Notify_ALTER_T1 ON DATABASE",
		// DROP EVENT NOTIFICATION on QUEUE
		"DROP EVENT NOTIFICATION NotifyQueue ON QUEUE dbo.ExpenseQueue",
		// DROP EVENT NOTIFICATION multiple names
		"DROP EVENT NOTIFICATION log_ddl1, log_ddl2 ON SERVER",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseEventSession tests batch 56: CREATE/ALTER/DROP EVENT SESSION (Extended Events).
func TestParseEventSession(t *testing.T) {
	tests := []string{
		// CREATE EVENT SESSION basic
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting",
		// CREATE EVENT SESSION with multiple events
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.rpc_starting, ADD EVENT sqlserver.sql_batch_starting",
		// CREATE EVENT SESSION with target
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting ADD TARGET package0.event_file (SET filename = N'C:\\xe\\test.xel')",
		// CREATE EVENT SESSION with target and multiple SET options
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting ADD TARGET package0.event_file (SET filename = N'C:\\xe\\test.xel', max_file_size = 256, max_rollover_files = 10)",
		// CREATE EVENT SESSION with WITH options
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting WITH (MAX_MEMORY = 4 MB)",
		// CREATE EVENT SESSION with multiple WITH options
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting WITH (MAX_MEMORY = 4 MB, STARTUP_STATE = ON, TRACK_CAUSALITY = OFF)",
		// CREATE EVENT SESSION on DATABASE
		"CREATE EVENT SESSION test_session ON DATABASE ADD EVENT sqlserver.sql_batch_starting",
		// CREATE EVENT SESSION with event SET and ACTION
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.error_reported (SET collect_message = 1 ACTION (sqlserver.sql_text, sqlserver.tsql_stack) WHERE severity >= 16)",
		// CREATE EVENT SESSION with WHERE predicate
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting (WHERE sqlserver.database_id = 5)",
		// CREATE EVENT SESSION with EVENT_RETENTION_MODE
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting WITH (EVENT_RETENTION_MODE = NO_EVENT_LOSS)",
		// CREATE EVENT SESSION with MAX_DISPATCH_LATENCY
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting WITH (MAX_DISPATCH_LATENCY = 30 SECONDS)",
		// CREATE EVENT SESSION with MEMORY_PARTITION_MODE
		"CREATE EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.sql_batch_starting WITH (MEMORY_PARTITION_MODE = PER_CPU)",
		// ALTER EVENT SESSION - STATE START
		"ALTER EVENT SESSION test_session ON SERVER STATE = START",
		// ALTER EVENT SESSION - STATE STOP
		"ALTER EVENT SESSION test_session ON SERVER STATE = STOP",
		// ALTER EVENT SESSION - ADD EVENT
		"ALTER EVENT SESSION test_session ON SERVER ADD EVENT sqlserver.database_transaction_begin, ADD EVENT sqlserver.database_transaction_end",
		// ALTER EVENT SESSION - DROP EVENT
		"ALTER EVENT SESSION test_session ON SERVER DROP EVENT sqlserver.sql_batch_starting",
		// ALTER EVENT SESSION - ADD TARGET
		"ALTER EVENT SESSION test_session ON SERVER ADD TARGET package0.ring_buffer (SET max_memory = 1024)",
		// ALTER EVENT SESSION - DROP TARGET
		"ALTER EVENT SESSION test_session ON SERVER DROP TARGET package0.event_file",
		// ALTER EVENT SESSION - WITH options
		"ALTER EVENT SESSION test_session ON SERVER WITH (MAX_MEMORY = 8 MB)",
		// ALTER EVENT SESSION on DATABASE
		"ALTER EVENT SESSION test_session ON DATABASE STATE = START",
		// DROP EVENT SESSION
		"DROP EVENT SESSION test_session ON SERVER",
		// DROP EVENT SESSION on DATABASE
		"DROP EVENT SESSION test_session ON DATABASE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseEventSessionIntegration tests batch 56: multi-statement integration.
func TestParseEventSessionIntegration(t *testing.T) {
	sql := `CREATE EVENT SESSION test_session ON SERVER
ADD EVENT sqlserver.rpc_starting,
ADD EVENT sqlserver.sql_batch_starting,
ADD EVENT sqlserver.error_reported
ADD TARGET package0.event_file (SET filename = N'C:\xe\test_session.xel', max_file_size = 256, max_rollover_files = 10)
WITH (MAX_MEMORY = 4 MB);
ALTER EVENT SESSION test_session ON SERVER STATE = START;
ALTER EVENT SESSION test_session ON SERVER STATE = STOP;
DROP EVENT SESSION test_session ON SERVER`
	result := ParseAndCheck(t, sql)
	if result.Len() != 4 {
		for i, item := range result.Items {
			t.Logf("  stmt[%d]: %T -> %s", i, item, ast.NodeToString(item))
		}
		t.Fatalf("Parse: got %d statements, want 4", result.Len())
	}
	for _, item := range result.Items {
		if _, ok := item.(*ast.SecurityStmt); !ok {
			t.Fatalf("expected *SecurityStmt, got %T: %s", item, ast.NodeToString(item))
		}
	}
}

// TestParseExternalDataSource tests batch 57: CREATE/ALTER/DROP EXTERNAL DATA SOURCE.
func TestParseExternalDataSource(t *testing.T) {
	tests := []string{
		// CREATE basic with LOCATION only
		"CREATE EXTERNAL DATA SOURCE MyHadoopCluster WITH (LOCATION = 'hdfs://10.10.10.10:8020')",
		// CREATE with LOCATION and TYPE
		"CREATE EXTERNAL DATA SOURCE MyHadoopCluster WITH (LOCATION = 'hdfs://10.10.10.10:8020', TYPE = HADOOP)",
		// CREATE with LOCATION, TYPE, and RESOURCE_MANAGER_LOCATION
		"CREATE EXTERNAL DATA SOURCE MyHadoopCluster WITH (LOCATION = 'hdfs://10.10.10.10:8020', TYPE = HADOOP, RESOURCE_MANAGER_LOCATION = '10.10.10.10:8032')",
		// CREATE with CREDENTIAL
		"CREATE EXTERNAL DATA SOURCE MyAzureStorage WITH (LOCATION = 'wasbs://container@account.blob.core.windows.net', CREDENTIAL = AzureStorageCredential)",
		// CREATE with TYPE = BLOB_STORAGE
		"CREATE EXTERNAL DATA SOURCE MyAzureBlob WITH (LOCATION = 'https://account.blob.core.windows.net', TYPE = BLOB_STORAGE, CREDENTIAL = AzureStorageCredential)",
		// CREATE with PUSHDOWN
		"CREATE EXTERNAL DATA SOURCE MySqlServer WITH (LOCATION = 'sqlserver://sql.database.windows.net', PUSHDOWN = ON, CREDENTIAL = SqlServerCredentials)",
		// CREATE with CONNECTION_OPTIONS
		"CREATE EXTERNAL DATA SOURCE MyOracleServer WITH (LOCATION = 'oracle://oracle.host:1521', CONNECTION_OPTIONS = 'ImpersonateUser=%CURRENT_USER', CREDENTIAL = OracleProxyCredential, PUSHDOWN = ON)",
		// CREATE SQL Server 2022+ style (no TYPE)
		"CREATE EXTERNAL DATA SOURCE MyABS WITH (LOCATION = 'abs://account.blob.core.windows.net', CREDENTIAL = AzureStorageCredential)",
		// ALTER EXTERNAL DATA SOURCE
		"ALTER EXTERNAL DATA SOURCE hadoop_eds SET LOCATION = 'hdfs://10.10.10.10:8020', RESOURCE_MANAGER_LOCATION = '10.10.10.10:8032'",
		// ALTER with CREDENTIAL
		"ALTER EXTERNAL DATA SOURCE hadoop_eds SET CREDENTIAL = new_hadoop_user",
		// DROP EXTERNAL DATA SOURCE
		"DROP EXTERNAL DATA SOURCE mydatasource",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseExternalTable tests batch 57: CREATE/DROP EXTERNAL TABLE.
func TestParseExternalTable(t *testing.T) {
	tests := []string{
		// CREATE EXTERNAL TABLE basic
		`CREATE EXTERNAL TABLE dbo.ClickStream (
			url VARCHAR(50),
			event_date DATE,
			user_ip VARCHAR(50)
		)
		WITH (
			LOCATION = '/webdata/employee.tbl',
			DATA_SOURCE = MyHadoopCluster,
			FILE_FORMAT = MyFileFormat
		)`,
		// CREATE with reject options
		`CREATE EXTERNAL TABLE dbo.ErrorTable (
			id INT,
			name VARCHAR(100)
		)
		WITH (
			LOCATION = '/data/errors/',
			DATA_SOURCE = MyAzureStorage,
			FILE_FORMAT = CsvFormat,
			REJECT_TYPE = value,
			REJECT_VALUE = 10
		)`,
		// CREATE with percentage reject and sample
		`CREATE EXTERNAL TABLE dbo.DataTable (
			col1 INT,
			col2 NVARCHAR(200)
		)
		WITH (
			LOCATION = '/data/files/',
			DATA_SOURCE = MySource,
			FILE_FORMAT = ParquetFormat,
			REJECT_TYPE = percentage,
			REJECT_VALUE = 5,
			REJECT_SAMPLE_VALUE = 1000
		)`,
		// CREATE with REJECTED_ROW_LOCATION
		`CREATE EXTERNAL TABLE dbo.RejectTest (
			id INT NOT NULL,
			value FLOAT
		)
		WITH (
			LOCATION = '/data/',
			DATA_SOURCE = MySource,
			FILE_FORMAT = CsvFormat,
			REJECT_TYPE = value,
			REJECT_VALUE = 0,
			REJECTED_ROW_LOCATION = '/REJECT_Directory'
		)`,
		// CREATE EXTERNAL TABLE with schema qualification
		`CREATE EXTERNAL TABLE mydb.dbo.ExternalData (
			id BIGINT,
			description NVARCHAR(MAX)
		)
		WITH (
			LOCATION = '/external/',
			DATA_SOURCE = RemoteSource,
			FILE_FORMAT = TextFormat
		)`,
		// DROP EXTERNAL TABLE
		"DROP EXTERNAL TABLE dbo.ClickStream",
		// DROP EXTERNAL TABLE qualified
		"DROP EXTERNAL TABLE mydb.dbo.ExternalData",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseExternalFileFormat tests batch 57: CREATE/DROP EXTERNAL FILE FORMAT.
func TestParseExternalFileFormat(t *testing.T) {
	tests := []string{
		// DELIMITEDTEXT format
		`CREATE EXTERNAL FILE FORMAT textdelimited1 WITH (FORMAT_TYPE = DELIMITEDTEXT, FORMAT_OPTIONS (FIELD_TERMINATOR = '|', DATE_FORMAT = 'MM/dd/yyyy'), DATA_COMPRESSION = 'org.apache.hadoop.io.compress.GzipCodec')`,
		// RCFILE format
		`CREATE EXTERNAL FILE FORMAT rcfile1 WITH (FORMAT_TYPE = RCFILE, SERDE_METHOD = 'org.apache.hadoop.hive.serde2.columnar.LazyBinaryColumnarSerDe', DATA_COMPRESSION = 'org.apache.hadoop.io.compress.DefaultCodec')`,
		// ORC format
		`CREATE EXTERNAL FILE FORMAT orcfile1 WITH (FORMAT_TYPE = ORC, DATA_COMPRESSION = 'org.apache.hadoop.io.compress.SnappyCodec')`,
		// PARQUET format
		`CREATE EXTERNAL FILE FORMAT parquetfile1 WITH (FORMAT_TYPE = PARQUET, DATA_COMPRESSION = 'org.apache.hadoop.io.compress.SnappyCodec')`,
		// PARQUET without compression
		`CREATE EXTERNAL FILE FORMAT parquetfile2 WITH (FORMAT_TYPE = PARQUET)`,
		// DELTA format
		`CREATE EXTERNAL FILE FORMAT DeltaFileFormat WITH (FORMAT_TYPE = DELTA)`,
		// CSV with skip header
		`CREATE EXTERNAL FILE FORMAT skipHeader_CSV WITH (FORMAT_TYPE = DELIMITEDTEXT, FORMAT_OPTIONS (FIELD_TERMINATOR = ',', STRING_DELIMITER = '"', FIRST_ROW = 2, USE_TYPE_DEFAULT = True))`,
		// CSV with ENCODING
		`CREATE EXTERNAL FILE FORMAT utf16format WITH (FORMAT_TYPE = DELIMITEDTEXT, FORMAT_OPTIONS (FIELD_TERMINATOR = ',', ENCODING = 'UTF8'))`,
		// DROP EXTERNAL FILE FORMAT
		"DROP EXTERNAL FILE FORMAT textdelimited1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// ---------- Batch 58: encryption_keys ----------

func TestParseEncryptionKeys(t *testing.T) {
	tests := []struct {
		sql        string
		action     string
		objectType string
		name       string
	}{
		// CREATE COLUMN ENCRYPTION KEY
		{
			sql:        "CREATE COLUMN ENCRYPTION KEY MyCEK WITH VALUES (COLUMN_MASTER_KEY = MyCMK, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0x01AB)",
			action:     "CREATE",
			objectType: "COLUMN ENCRYPTION KEY",
			name:       "MyCEK",
		},
		// CREATE COLUMN ENCRYPTION KEY with two values
		{
			sql:        "CREATE COLUMN ENCRYPTION KEY TwoValueCEK WITH VALUES (COLUMN_MASTER_KEY = CMK1, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0x016E), (COLUMN_MASTER_KEY = CMK2, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0x016E)",
			action:     "CREATE",
			objectType: "COLUMN ENCRYPTION KEY",
			name:       "TwoValueCEK",
		},
		// ALTER COLUMN ENCRYPTION KEY - ADD VALUE
		{
			sql:        "ALTER COLUMN ENCRYPTION KEY MyCEK ADD VALUE (COLUMN_MASTER_KEY = MyCMK2, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0x016E)",
			action:     "ALTER",
			objectType: "COLUMN ENCRYPTION KEY",
			name:       "MyCEK",
		},
		// ALTER COLUMN ENCRYPTION KEY - DROP VALUE
		{
			sql:        "ALTER COLUMN ENCRYPTION KEY MyCEK DROP VALUE (COLUMN_MASTER_KEY = MyCMK)",
			action:     "ALTER",
			objectType: "COLUMN ENCRYPTION KEY",
			name:       "MyCEK",
		},
		// DROP COLUMN ENCRYPTION KEY
		{
			sql:        "DROP COLUMN ENCRYPTION KEY MyCEK",
			action:     "DROP",
			objectType: "COLUMN ENCRYPTION KEY",
			name:       "MyCEK",
		},
		// CREATE COLUMN MASTER KEY
		{
			sql:        "CREATE COLUMN MASTER KEY MyCMK WITH (KEY_STORE_PROVIDER_NAME = N'MSSQL_CERTIFICATE_STORE', KEY_PATH = 'CurrentUser/My/BBF037EC')",
			action:     "CREATE",
			objectType: "COLUMN MASTER KEY",
			name:       "MyCMK",
		},
		// CREATE COLUMN MASTER KEY with ENCLAVE_COMPUTATIONS
		{
			sql:        "CREATE COLUMN MASTER KEY MyCMK WITH (KEY_STORE_PROVIDER_NAME = N'AZURE_KEY_VAULT', KEY_PATH = N'https://myvault.vault.azure.net:443/keys/MyCMK/abc', ENCLAVE_COMPUTATIONS (SIGNATURE = 0xA80F))",
			action:     "CREATE",
			objectType: "COLUMN MASTER KEY",
			name:       "MyCMK",
		},
		// DROP COLUMN MASTER KEY
		{
			sql:        "DROP COLUMN MASTER KEY MyCMK",
			action:     "DROP",
			objectType: "COLUMN MASTER KEY",
			name:       "MyCMK",
		},
		// CREATE DATABASE ENCRYPTION KEY
		{
			sql:        "CREATE DATABASE ENCRYPTION KEY WITH ALGORITHM = AES_256 ENCRYPTION BY SERVER CERTIFICATE MyServerCert",
			action:     "CREATE",
			objectType: "DATABASE ENCRYPTION KEY",
			name:       "",
		},
		// CREATE DATABASE ENCRYPTION KEY with ASYMMETRIC KEY
		{
			sql:        "CREATE DATABASE ENCRYPTION KEY WITH ALGORITHM = AES_128 ENCRYPTION BY SERVER ASYMMETRIC KEY MyAsymKey",
			action:     "CREATE",
			objectType: "DATABASE ENCRYPTION KEY",
			name:       "",
		},
		// ALTER DATABASE ENCRYPTION KEY - REGENERATE
		{
			sql:        "ALTER DATABASE ENCRYPTION KEY REGENERATE WITH ALGORITHM = AES_256",
			action:     "ALTER",
			objectType: "DATABASE ENCRYPTION KEY",
			name:       "",
		},
		// ALTER DATABASE ENCRYPTION KEY - ENCRYPTION BY
		{
			sql:        "ALTER DATABASE ENCRYPTION KEY ENCRYPTION BY SERVER CERTIFICATE MyServerCert",
			action:     "ALTER",
			objectType: "DATABASE ENCRYPTION KEY",
			name:       "",
		},
		// DROP DATABASE ENCRYPTION KEY
		{
			sql:        "DROP DATABASE ENCRYPTION KEY",
			action:     "DROP",
			objectType: "DATABASE ENCRYPTION KEY",
			name:       "",
		},
		// CREATE DATABASE SCOPED CREDENTIAL
		{
			sql:        "CREATE DATABASE SCOPED CREDENTIAL AppCred WITH IDENTITY = 'Mary5', SECRET = 'mypassword'",
			action:     "CREATE",
			objectType: "DATABASE SCOPED CREDENTIAL",
			name:       "AppCred",
		},
		// CREATE DATABASE SCOPED CREDENTIAL with SAS
		{
			sql:        "CREATE DATABASE SCOPED CREDENTIAL MyCredentials WITH IDENTITY = 'SHARED ACCESS SIGNATURE', SECRET = 'sas_token'",
			action:     "CREATE",
			objectType: "DATABASE SCOPED CREDENTIAL",
			name:       "MyCredentials",
		},
		// CREATE DATABASE SCOPED CREDENTIAL with MANAGED IDENTITY
		{
			sql:        "CREATE DATABASE SCOPED CREDENTIAL managed_id WITH IDENTITY = 'Managed Identity'",
			action:     "CREATE",
			objectType: "DATABASE SCOPED CREDENTIAL",
			name:       "managed_id",
		},
		// ALTER DATABASE SCOPED CREDENTIAL
		{
			sql:        "ALTER DATABASE SCOPED CREDENTIAL AppCred WITH IDENTITY = 'Mary6', SECRET = 'newpassword'",
			action:     "ALTER",
			objectType: "DATABASE SCOPED CREDENTIAL",
			name:       "AppCred",
		},
		// DROP DATABASE SCOPED CREDENTIAL
		{
			sql:        "DROP DATABASE SCOPED CREDENTIAL AppCred",
			action:     "DROP",
			objectType: "DATABASE SCOPED CREDENTIAL",
			name:       "AppCred",
		},
		// CREATE CRYPTOGRAPHIC PROVIDER
		{
			sql:        "CREATE CRYPTOGRAPHIC PROVIDER SecurityProvider FROM FILE = 'C:\\SecurityProvider\\SecurityProvider_v1.dll'",
			action:     "CREATE",
			objectType: "CRYPTOGRAPHIC PROVIDER",
			name:       "SecurityProvider",
		},
		// ALTER CRYPTOGRAPHIC PROVIDER - ENABLE
		{
			sql:        "ALTER CRYPTOGRAPHIC PROVIDER SecurityProvider ENABLE",
			action:     "ALTER",
			objectType: "CRYPTOGRAPHIC PROVIDER",
			name:       "SecurityProvider",
		},
		// ALTER CRYPTOGRAPHIC PROVIDER - DISABLE
		{
			sql:        "ALTER CRYPTOGRAPHIC PROVIDER SecurityProvider DISABLE",
			action:     "ALTER",
			objectType: "CRYPTOGRAPHIC PROVIDER",
			name:       "SecurityProvider",
		},
		// ALTER CRYPTOGRAPHIC PROVIDER - FROM FILE
		{
			sql:        "ALTER CRYPTOGRAPHIC PROVIDER SecurityProvider FROM FILE = 'c:\\SecurityProvider\\SecurityProvider_v2.dll' ENABLE",
			action:     "ALTER",
			objectType: "CRYPTOGRAPHIC PROVIDER",
			name:       "SecurityProvider",
		},
		// DROP CRYPTOGRAPHIC PROVIDER
		{
			sql:        "DROP CRYPTOGRAPHIC PROVIDER SecurityProvider",
			action:     "DROP",
			objectType: "CRYPTOGRAPHIC PROVIDER",
			name:       "SecurityProvider",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityKeyStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityKeyStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Action != tt.action {
				t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
			}
			if stmt.ObjectType != tt.objectType {
				t.Errorf("Parse(%q): objectType = %q, want %q", tt.sql, stmt.ObjectType, tt.objectType)
			}
			if stmt.Name != tt.name {
				t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
			}
			checkLocation(t, tt.sql, "SecurityKeyStmt", stmt.Loc)
		})
	}
}

// TestParseResourceGovernor tests batch 59: Resource Governor statements.
func TestParseResourceGovernor(t *testing.T) {
	tests := []struct {
		sql        string
		action     string
		objectType string
		name       string
	}{
		// CREATE WORKLOAD GROUP - basic with no options
		{
			sql:        "CREATE WORKLOAD GROUP newReports",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "newReports",
		},
		// CREATE WORKLOAD GROUP - with multiple options
		{
			sql:        "CREATE WORKLOAD GROUP newReports WITH (REQUEST_MAX_MEMORY_GRANT_PERCENT = 2.5, REQUEST_MAX_CPU_TIME_SEC = 100, MAX_DOP = 4)",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "newReports",
		},
		// CREATE WORKLOAD GROUP - with IMPORTANCE
		{
			sql:        "CREATE WORKLOAD GROUP highPriority WITH (IMPORTANCE = HIGH)",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "highPriority",
		},
		// CREATE WORKLOAD GROUP - with USING
		{
			sql:        "CREATE WORKLOAD GROUP newReports WITH (MAX_DOP = 4) USING myPool",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "newReports",
		},
		// CREATE WORKLOAD GROUP - with USING default
		{
			sql:        "CREATE WORKLOAD GROUP newReports USING [default]",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "newReports",
		},
		// CREATE WORKLOAD GROUP - with USING pool and EXTERNAL pool
		{
			sql:        "CREATE WORKLOAD GROUP newReports USING myPool, EXTERNAL myExtPool",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "newReports",
		},
		// CREATE WORKLOAD GROUP - with all options
		{
			sql:        "CREATE WORKLOAD GROUP fullGroup WITH (IMPORTANCE = MEDIUM, REQUEST_MAX_MEMORY_GRANT_PERCENT = 25, REQUEST_MAX_CPU_TIME_SEC = 60, REQUEST_MEMORY_GRANT_TIMEOUT_SEC = 30, MAX_DOP = 8, GROUP_MAX_REQUESTS = 100)",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "fullGroup",
		},
		// CREATE WORKLOAD GROUP - with GROUP_MAX_TEMPDB_DATA_MB
		{
			sql:        "CREATE WORKLOAD GROUP tempGroup WITH (GROUP_MAX_TEMPDB_DATA_MB = 500, GROUP_MAX_TEMPDB_DATA_PERCENT = 25)",
			action:     "CREATE",
			objectType: "WORKLOAD GROUP",
			name:       "tempGroup",
		},
		// ALTER WORKLOAD GROUP
		{
			sql:        "ALTER WORKLOAD GROUP [default] WITH (IMPORTANCE = LOW)",
			action:     "ALTER",
			objectType: "WORKLOAD GROUP",
			name:       "default",
		},
		// ALTER WORKLOAD GROUP - move to different pool
		{
			sql:        "ALTER WORKLOAD GROUP adHoc USING [default]",
			action:     "ALTER",
			objectType: "WORKLOAD GROUP",
			name:       "adHoc",
		},
		// ALTER WORKLOAD GROUP - with multiple options
		{
			sql:        "ALTER WORKLOAD GROUP myGroup WITH (MAX_DOP = 2, GROUP_MAX_REQUESTS = 50) USING poolA",
			action:     "ALTER",
			objectType: "WORKLOAD GROUP",
			name:       "myGroup",
		},
		// DROP WORKLOAD GROUP
		{
			sql:        "DROP WORKLOAD GROUP adhoc",
			action:     "DROP",
			objectType: "WORKLOAD GROUP",
			name:       "adhoc",
		},
		// CREATE RESOURCE POOL - basic with no options
		{
			sql:        "CREATE RESOURCE POOL bigPool",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "bigPool",
		},
		// CREATE RESOURCE POOL - with CPU options
		{
			sql:        "CREATE RESOURCE POOL adhocPool WITH (MIN_CPU_PERCENT = 10, MAX_CPU_PERCENT = 20, CAP_CPU_PERCENT = 30)",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "adhocPool",
		},
		// CREATE RESOURCE POOL - with memory options
		{
			sql:        "CREATE RESOURCE POOL memPool WITH (MIN_MEMORY_PERCENT = 5, MAX_MEMORY_PERCENT = 15)",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "memPool",
		},
		// CREATE RESOURCE POOL - with IOPS options
		{
			sql:        "CREATE RESOURCE POOL PoolAdmin WITH (MIN_IOPS_PER_VOLUME = 200, MAX_IOPS_PER_VOLUME = 1000)",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "PoolAdmin",
		},
		// CREATE RESOURCE POOL - with AFFINITY SCHEDULER
		{
			sql:        "CREATE RESOURCE POOL cpuPool WITH (AFFINITY SCHEDULER = AUTO)",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "cpuPool",
		},
		// CREATE RESOURCE POOL - with AFFINITY SCHEDULER range
		{
			sql:        "CREATE RESOURCE POOL cpuPool WITH (AFFINITY SCHEDULER = (0 TO 63, 128 TO 191))",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "cpuPool",
		},
		// CREATE RESOURCE POOL - with AFFINITY NUMANODE
		{
			sql:        "CREATE RESOURCE POOL numaPool WITH (AFFINITY NUMANODE = (0, 1))",
			action:     "CREATE",
			objectType: "RESOURCE POOL",
			name:       "numaPool",
		},
		// ALTER RESOURCE POOL
		{
			sql:        "ALTER RESOURCE POOL [default] WITH (MAX_CPU_PERCENT = 25)",
			action:     "ALTER",
			objectType: "RESOURCE POOL",
			name:       "default",
		},
		// ALTER RESOURCE POOL - with multiple options
		{
			sql:        "ALTER RESOURCE POOL adhocPool WITH (MIN_CPU_PERCENT = 10, MAX_CPU_PERCENT = 20, CAP_CPU_PERCENT = 30, MIN_MEMORY_PERCENT = 5, MAX_MEMORY_PERCENT = 15, AFFINITY SCHEDULER = (0 TO 63, 128 TO 191))",
			action:     "ALTER",
			objectType: "RESOURCE POOL",
			name:       "adhocPool",
		},
		// DROP RESOURCE POOL
		{
			sql:        "DROP RESOURCE POOL big_pool",
			action:     "DROP",
			objectType: "RESOURCE POOL",
			name:       "big_pool",
		},
		// CREATE EXTERNAL RESOURCE POOL
		{
			sql:        "CREATE EXTERNAL RESOURCE POOL ep_1 WITH (MAX_CPU_PERCENT = 75, MAX_MEMORY_PERCENT = 30)",
			action:     "CREATE",
			objectType: "EXTERNAL RESOURCE POOL",
			name:       "ep_1",
		},
		// CREATE EXTERNAL RESOURCE POOL - with MAX_PROCESSES
		{
			sql:        "CREATE EXTERNAL RESOURCE POOL ep_2 WITH (MAX_CPU_PERCENT = 50, MAX_MEMORY_PERCENT = 25, MAX_PROCESSES = 4)",
			action:     "CREATE",
			objectType: "EXTERNAL RESOURCE POOL",
			name:       "ep_2",
		},
		// ALTER EXTERNAL RESOURCE POOL
		{
			sql:        "ALTER EXTERNAL RESOURCE POOL ep_1 WITH (MAX_CPU_PERCENT = 50, MAX_MEMORY_PERCENT = 25)",
			action:     "ALTER",
			objectType: "EXTERNAL RESOURCE POOL",
			name:       "ep_1",
		},
		// ALTER EXTERNAL RESOURCE POOL - default
		{
			sql:        `ALTER EXTERNAL RESOURCE POOL [default] WITH (MAX_CPU_PERCENT = 75)`,
			action:     "ALTER",
			objectType: "EXTERNAL RESOURCE POOL",
			name:       "default",
		},
		// DROP EXTERNAL RESOURCE POOL
		{
			sql:        "DROP EXTERNAL RESOURCE POOL ex_pool",
			action:     "DROP",
			objectType: "EXTERNAL RESOURCE POOL",
			name:       "ex_pool",
		},
		// ALTER RESOURCE GOVERNOR - RECONFIGURE
		{
			sql:        "ALTER RESOURCE GOVERNOR RECONFIGURE",
			action:     "ALTER",
			objectType: "RESOURCE GOVERNOR",
			name:       "",
		},
		// ALTER RESOURCE GOVERNOR - DISABLE
		{
			sql:        "ALTER RESOURCE GOVERNOR DISABLE",
			action:     "ALTER",
			objectType: "RESOURCE GOVERNOR",
			name:       "",
		},
		// ALTER RESOURCE GOVERNOR - RESET STATISTICS
		{
			sql:        "ALTER RESOURCE GOVERNOR RESET STATISTICS",
			action:     "ALTER",
			objectType: "RESOURCE GOVERNOR",
			name:       "",
		},
		// ALTER RESOURCE GOVERNOR - WITH CLASSIFIER_FUNCTION
		{
			sql:        "ALTER RESOURCE GOVERNOR WITH (CLASSIFIER_FUNCTION = dbo.rg_classifier)",
			action:     "ALTER",
			objectType: "RESOURCE GOVERNOR",
			name:       "",
		},
		// ALTER RESOURCE GOVERNOR - WITH CLASSIFIER_FUNCTION = NULL
		{
			sql:        "ALTER RESOURCE GOVERNOR WITH (CLASSIFIER_FUNCTION = NULL)",
			action:     "ALTER",
			objectType: "RESOURCE GOVERNOR",
			name:       "",
		},
		// ALTER RESOURCE GOVERNOR - WITH MAX_OUTSTANDING_IO_PER_VOLUME
		{
			sql:        "ALTER RESOURCE GOVERNOR WITH (MAX_OUTSTANDING_IO_PER_VOLUME = 20)",
			action:     "ALTER",
			objectType: "RESOURCE GOVERNOR",
			name:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Action != tt.action {
				t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
			}
			if stmt.ObjectType != tt.objectType {
				t.Errorf("Parse(%q): objectType = %q, want %q", tt.sql, stmt.ObjectType, tt.objectType)
			}
			if stmt.Name != tt.name {
				t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseResourceGovernorIntegration tests batch 59: multi-statement integration.
func TestParseResourceGovernorIntegration(t *testing.T) {
	sql := `CREATE RESOURCE POOL adhocPool WITH (MIN_CPU_PERCENT = 10, MAX_CPU_PERCENT = 20);
CREATE WORKLOAD GROUP adHoc WITH (IMPORTANCE = LOW, MAX_DOP = 4) USING adhocPool;
ALTER RESOURCE GOVERNOR RECONFIGURE;
ALTER WORKLOAD GROUP adHoc WITH (IMPORTANCE = MEDIUM);
ALTER RESOURCE GOVERNOR RECONFIGURE;
DROP WORKLOAD GROUP adHoc;
DROP RESOURCE POOL adhocPool;
ALTER RESOURCE GOVERNOR RECONFIGURE`
	result := ParseAndCheck(t, sql)
	if result.Len() != 8 {
		for i, item := range result.Items {
			t.Logf("  stmt[%d]: %T -> %s", i, item, ast.NodeToString(item))
		}
		t.Fatalf("Parse: got %d statements, want 8", result.Len())
	}
	for _, item := range result.Items {
		if _, ok := item.(*ast.SecurityStmt); !ok {
			t.Fatalf("expected *SecurityStmt, got %T: %s", item, ast.NodeToString(item))
		}
	}
}

// TestParseServerLevelObjects tests batch 60: CREATE/ALTER/DROP SERVER ROLE,
// ALTER SERVER CONFIGURATION, ALTER/BACKUP/RESTORE SERVICE MASTER KEY.
func TestParseServerLevelObjects(t *testing.T) {
	t.Run("server_role", func(t *testing.T) {
		tests := []struct {
			sql        string
			action     string
			objectType string
			name       string
		}{
			// CREATE SERVER ROLE - basic
			{
				sql:        "CREATE SERVER ROLE auditors",
				action:     "CREATE",
				objectType: "SERVER ROLE",
				name:       "auditors",
			},
			// CREATE SERVER ROLE - with AUTHORIZATION
			{
				sql:        "CREATE SERVER ROLE auditors AUTHORIZATION securityadmin",
				action:     "CREATE",
				objectType: "SERVER ROLE",
				name:       "auditors",
			},
			// ALTER SERVER ROLE - ADD MEMBER
			{
				sql:        "ALTER SERVER ROLE sysadmin ADD MEMBER [MyDomain\\TestUser]",
				action:     "ALTER",
				objectType: "SERVER ROLE",
				name:       "sysadmin",
			},
			// ALTER SERVER ROLE - DROP MEMBER
			{
				sql:        "ALTER SERVER ROLE diskadmin DROP MEMBER TestLogin",
				action:     "ALTER",
				objectType: "SERVER ROLE",
				name:       "diskadmin",
			},
			// ALTER SERVER ROLE - WITH NAME
			{
				sql:        "ALTER SERVER ROLE buyers WITH NAME = purchasing",
				action:     "ALTER",
				objectType: "SERVER ROLE",
				name:       "buyers",
			},
			// DROP SERVER ROLE
			{
				sql:        "DROP SERVER ROLE purchasing",
				action:     "DROP",
				objectType: "SERVER ROLE",
				name:       "purchasing",
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				if stmt.ObjectType != tt.objectType {
					t.Errorf("Parse(%q): objectType = %q, want %q", tt.sql, stmt.ObjectType, tt.objectType)
				}
				if stmt.Name != tt.name {
					t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
				}
				checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
			})
		}
	})

	t.Run("alter_server_configuration", func(t *testing.T) {
		tests := []struct {
			sql        string
			optionType string
		}{
			// PROCESS AFFINITY CPU
			{
				sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY CPU = AUTO",
				optionType: "PROCESS AFFINITY",
			},
			// PROCESS AFFINITY CPU range
			{
				sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY CPU = 0 TO 3, 8 TO 11",
				optionType: "PROCESS AFFINITY",
			},
			// PROCESS AFFINITY NUMANODE
			{
				sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY NUMANODE = 0 TO 1",
				optionType: "PROCESS AFFINITY",
			},
			// DIAGNOSTICS LOG ON
			{
				sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG ON",
				optionType: "DIAGNOSTICS LOG",
			},
			// DIAGNOSTICS LOG OFF
			{
				sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG OFF",
				optionType: "DIAGNOSTICS LOG",
			},
			// DIAGNOSTICS LOG PATH
			{
				sql:        `ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG PATH = 'C:\MSSSQL\Logs'`,
				optionType: "DIAGNOSTICS LOG",
			},
			// DIAGNOSTICS LOG MAX_SIZE
			{
				sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG MAX_SIZE = 20 MB",
				optionType: "DIAGNOSTICS LOG",
			},
			// FAILOVER CLUSTER PROPERTY
			{
				sql:        "ALTER SERVER CONFIGURATION SET FAILOVER CLUSTER PROPERTY VerboseLogging = 2",
				optionType: "FAILOVER CLUSTER PROPERTY",
			},
			// HADR CLUSTER CONTEXT
			{
				sql:        `ALTER SERVER CONFIGURATION SET HADR CLUSTER CONTEXT = 'clus01.xyz.com'`,
				optionType: "HADR CLUSTER CONTEXT",
			},
			// HADR CLUSTER CONTEXT LOCAL
			{
				sql:        "ALTER SERVER CONFIGURATION SET HADR CLUSTER CONTEXT = LOCAL",
				optionType: "HADR CLUSTER CONTEXT",
			},
			// BUFFER POOL EXTENSION ON
			{
				sql:        `ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION ON (FILENAME = 'F:\SSDCACHE\Example.BPE', SIZE = 50 GB)`,
				optionType: "BUFFER POOL EXTENSION",
			},
			// BUFFER POOL EXTENSION OFF
			{
				sql:        "ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION OFF",
				optionType: "BUFFER POOL EXTENSION",
			},
			// SOFTNUMA ON
			{
				sql:        "ALTER SERVER CONFIGURATION SET SOFTNUMA ON",
				optionType: "SOFTNUMA",
			},
			// SOFTNUMA OFF
			{
				sql:        "ALTER SERVER CONFIGURATION SET SOFTNUMA OFF",
				optionType: "SOFTNUMA",
			},
			// MEMORY_OPTIMIZED TEMPDB_METADATA ON
			{
				sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED TEMPDB_METADATA = ON",
				optionType: "MEMORY_OPTIMIZED",
			},
			// MEMORY_OPTIMIZED HYBRID_BUFFER_POOL OFF
			{
				sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED HYBRID_BUFFER_POOL = OFF",
				optionType: "MEMORY_OPTIMIZED",
			},
			// HARDWARE_OFFLOAD ON
			{
				sql:        "ALTER SERVER CONFIGURATION SET HARDWARE_OFFLOAD ON",
				optionType: "HARDWARE_OFFLOAD",
			},
			// HARDWARE_OFFLOAD OFF
			{
				sql:        "ALTER SERVER CONFIGURATION SET HARDWARE_OFFLOAD OFF",
				optionType: "HARDWARE_OFFLOAD",
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterServerConfigurationStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *AlterServerConfigurationStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.OptionType != tt.optionType {
					t.Errorf("Parse(%q): optionType = %q, want %q", tt.sql, stmt.OptionType, tt.optionType)
				}
				checkLocation(t, tt.sql, "AlterServerConfigurationStmt", stmt.Loc)
			})
		}
	})

	t.Run("service_master_key", func(t *testing.T) {
		tests := []struct {
			sql        string
			action     string
			objectType string
		}{
			// ALTER SERVICE MASTER KEY REGENERATE
			{
				sql:        "ALTER SERVICE MASTER KEY REGENERATE",
				action:     "ALTER",
				objectType: "SERVICE MASTER KEY",
			},
			// ALTER SERVICE MASTER KEY FORCE REGENERATE
			{
				sql:        "ALTER SERVICE MASTER KEY FORCE REGENERATE",
				action:     "ALTER",
				objectType: "SERVICE MASTER KEY",
			},
			// ALTER SERVICE MASTER KEY WITH OLD_ACCOUNT / OLD_PASSWORD
			{
				sql:        "ALTER SERVICE MASTER KEY WITH OLD_ACCOUNT = 'NT AUTHORITY\\NetworkService', OLD_PASSWORD = 'old_p@ss'",
				action:     "ALTER",
				objectType: "SERVICE MASTER KEY",
			},
			// ALTER SERVICE MASTER KEY WITH NEW_ACCOUNT / NEW_PASSWORD
			{
				sql:        "ALTER SERVICE MASTER KEY WITH NEW_ACCOUNT = 'NT AUTHORITY\\NetworkService', NEW_PASSWORD = 'new_p@ss'",
				action:     "ALTER",
				objectType: "SERVICE MASTER KEY",
			},
			// BACKUP SERVICE MASTER KEY
			{
				sql:        "BACKUP SERVICE MASTER KEY TO FILE = 'c:\\temp\\exportedkey' ENCRYPTION BY PASSWORD = 'sd@34$#%^+)2Dwe23klj'",
				action:     "BACKUP",
				objectType: "SERVICE MASTER KEY",
			},
			// RESTORE SERVICE MASTER KEY
			{
				sql:        "RESTORE SERVICE MASTER KEY FROM FILE = 'c:\\temp\\exportedkey' DECRYPTION BY PASSWORD = 'sd@34$#%^+)2Dwe23klj'",
				action:     "RESTORE",
				objectType: "SERVICE MASTER KEY",
			},
			// RESTORE SERVICE MASTER KEY with FORCE
			{
				sql:        "RESTORE SERVICE MASTER KEY FROM FILE = 'c:\\temp\\exportedkey' DECRYPTION BY PASSWORD = 'sd@34$#%^+)2Dwe23klj' FORCE",
				action:     "RESTORE",
				objectType: "SERVICE MASTER KEY",
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityKeyStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityKeyStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				if stmt.ObjectType != tt.objectType {
					t.Errorf("Parse(%q): objectType = %q, want %q", tt.sql, stmt.ObjectType, tt.objectType)
				}
				checkLocation(t, tt.sql, "SecurityKeyStmt", stmt.Loc)
			})
		}
	})
}

// TestParseFulltextStoplistSearch tests batch 61: FULLTEXT STOPLIST and SEARCH PROPERTY LIST.
func TestParseFulltextStoplistSearch(t *testing.T) {
	t.Run("fulltext_stoplist", func(t *testing.T) {
		tests := []struct {
			sql  string
			name string
		}{
			// CREATE FULLTEXT STOPLIST - empty
			{`CREATE FULLTEXT STOPLIST myStoplist`, "myStoplist"},
			// CREATE FULLTEXT STOPLIST FROM existing
			{`CREATE FULLTEXT STOPLIST myStoplist2 FROM AdventureWorks.otherStoplist`, "myStoplist2"},
			// CREATE FULLTEXT STOPLIST FROM SYSTEM STOPLIST
			{`CREATE FULLTEXT STOPLIST myStoplist3 FROM SYSTEM STOPLIST`, "myStoplist3"},
			// CREATE FULLTEXT STOPLIST with AUTHORIZATION
			{`CREATE FULLTEXT STOPLIST myStoplist AUTHORIZATION dbo`, "myStoplist"},
			// CREATE FULLTEXT STOPLIST FROM source with AUTHORIZATION
			{`CREATE FULLTEXT STOPLIST sl1 FROM db1.sl2 AUTHORIZATION admin1`, "sl1"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateFulltextStoplistStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *CreateFulltextStoplistStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Name != tt.name {
					t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
				}
				checkLocation(t, tt.sql, "CreateFulltextStoplistStmt", stmt.Loc)
			})
		}
	})

	t.Run("alter_fulltext_stoplist", func(t *testing.T) {
		tests := []struct {
			sql    string
			name   string
			action string
		}{
			// ADD stopword
			{`ALTER FULLTEXT STOPLIST CombinedList ADD 'en' LANGUAGE 'Spanish'`, "CombinedList", "ADD"},
			// ADD with N prefix
			{`ALTER FULLTEXT STOPLIST sl1 ADD N'the' LANGUAGE 1033`, "sl1", "ADD"},
			// DROP single stopword
			{`ALTER FULLTEXT STOPLIST sl1 DROP 'en' LANGUAGE 'French'`, "sl1", "DROP"},
			// DROP ALL LANGUAGE
			{`ALTER FULLTEXT STOPLIST sl1 DROP ALL LANGUAGE 'English'`, "sl1", "DROP"},
			// DROP ALL
			{`ALTER FULLTEXT STOPLIST sl1 DROP ALL`, "sl1", "DROP"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterFulltextStoplistStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *AlterFulltextStoplistStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Name != tt.name {
					t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				checkLocation(t, tt.sql, "AlterFulltextStoplistStmt", stmt.Loc)
			})
		}
	})

	t.Run("drop_fulltext_stoplist", func(t *testing.T) {
		sql := `DROP FULLTEXT STOPLIST myStoplist`
		result := ParseAndCheck(t, sql)
		if result.Len() != 1 {
			t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
		}
		stmt, ok := result.Items[0].(*ast.DropFulltextStoplistStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected *DropFulltextStoplistStmt, got %T", sql, result.Items[0])
		}
		if stmt.Name != "myStoplist" {
			t.Errorf("Parse(%q): name = %q, want %q", sql, stmt.Name, "myStoplist")
		}
		checkLocation(t, sql, "DropFulltextStoplistStmt", stmt.Loc)
	})

	t.Run("search_property_list", func(t *testing.T) {
		tests := []struct {
			sql  string
			name string
		}{
			// CREATE SEARCH PROPERTY LIST - empty
			{`CREATE SEARCH PROPERTY LIST DocumentPropertyList`, "DocumentPropertyList"},
			// CREATE SEARCH PROPERTY LIST FROM existing
			{`CREATE SEARCH PROPERTY LIST JobCandidateProperties FROM AdventureWorks2022.DocumentPropertyList`, "JobCandidateProperties"},
			// CREATE SEARCH PROPERTY LIST FROM same db
			{`CREATE SEARCH PROPERTY LIST spl2 FROM spl1`, "spl2"},
			// CREATE SEARCH PROPERTY LIST with AUTHORIZATION
			{`CREATE SEARCH PROPERTY LIST spl1 AUTHORIZATION dbo`, "spl1"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateSearchPropertyListStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *CreateSearchPropertyListStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Name != tt.name {
					t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
				}
				checkLocation(t, tt.sql, "CreateSearchPropertyListStmt", stmt.Loc)
			})
		}
	})

	t.Run("alter_search_property_list", func(t *testing.T) {
		tests := []struct {
			sql    string
			name   string
			action string
		}{
			// ADD property with all options
			{`ALTER SEARCH PROPERTY LIST DocumentPropertyList ADD 'Title' WITH (PROPERTY_SET_GUID = 'F29F85E0-4FF9-1068-AB91-08002B27B3D9', PROPERTY_INT_ID = 2, PROPERTY_DESCRIPTION = 'Title of the item.')`, "DocumentPropertyList", "ADD"},
			// ADD property without description
			{`ALTER SEARCH PROPERTY LIST spl1 ADD 'Author' WITH (PROPERTY_SET_GUID = 'F29F85E0-4FF9-1068-AB91-08002B27B3D9', PROPERTY_INT_ID = 4)`, "spl1", "ADD"},
			// DROP property
			{`ALTER SEARCH PROPERTY LIST DocumentPropertyList DROP 'Comments'`, "DocumentPropertyList", "DROP"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterSearchPropertyListStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *AlterSearchPropertyListStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Name != tt.name {
					t.Errorf("Parse(%q): name = %q, want %q", tt.sql, stmt.Name, tt.name)
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				checkLocation(t, tt.sql, "AlterSearchPropertyListStmt", stmt.Loc)
			})
		}
	})

	t.Run("drop_search_property_list", func(t *testing.T) {
		sql := `DROP SEARCH PROPERTY LIST JobCandidateProperties`
		result := ParseAndCheck(t, sql)
		if result.Len() != 1 {
			t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
		}
		stmt, ok := result.Items[0].(*ast.DropSearchPropertyListStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected *DropSearchPropertyListStmt, got %T", sql, result.Items[0])
		}
		if stmt.Name != "JobCandidateProperties" {
			t.Errorf("Parse(%q): name = %q, want %q", sql, stmt.Name, "JobCandidateProperties")
		}
		checkLocation(t, sql, "DropSearchPropertyListStmt", stmt.Loc)
	})
}

// TestParseSecurityPolicyClassification tests batch 62: SECURITY POLICY, SENSITIVITY CLASSIFICATION, SIGNATURE.
func TestParseSecurityPolicyClassification(t *testing.T) {
	t.Run("security_policy", func(t *testing.T) {
		tests := []struct {
			sql    string
			action string
		}{
			// CREATE with filter predicate
			{`CREATE SECURITY POLICY FederatedSecurityPolicy ADD FILTER PREDICATE rls.fn_securitypredicate(CustomerId) ON dbo.Customer`, "CREATE"},
			// CREATE with multiple predicates and WITH options
			{`CREATE SECURITY POLICY FederatedSecurityPolicy ADD FILTER PREDICATE rls.fn_pred1(CustomerId) ON dbo.Customer, ADD FILTER PREDICATE rls.fn_pred2(VendorId) ON dbo.Vendor WITH (STATE = ON)`, "CREATE"},
			// CREATE with block predicate AFTER INSERT
			{`CREATE SECURITY POLICY rls.SecPol ADD FILTER PREDICATE rls.tenantPred(TenantId) ON dbo.Sales, ADD BLOCK PREDICATE rls.tenantPred(TenantId) ON dbo.Sales AFTER INSERT`, "CREATE"},
			// CREATE with NOT FOR REPLICATION
			{`CREATE SECURITY POLICY pol1 ADD FILTER PREDICATE dbo.fn1(col1) ON dbo.t1 NOT FOR REPLICATION`, "CREATE"},
			// ALTER with STATE = ON
			{`ALTER SECURITY POLICY pol1 WITH (STATE = ON)`, "ALTER"},
			// ALTER ADD predicate
			{`ALTER SECURITY POLICY pol1 ADD FILTER PREDICATE schema_preds.SecPredicate(column1) ON myschema.mytable`, "ALTER"},
			// ALTER DROP predicate
			{`ALTER SECURITY POLICY pol1 DROP FILTER PREDICATE ON myschema.mytable`, "ALTER"},
			// ALTER with ALTER predicate
			{`ALTER SECURITY POLICY pol1 ALTER FILTER PREDICATE schema_preds.SecPredicate2(column1) ON myschema.mytable`, "ALTER"},
			// ALTER with BLOCK predicate BEFORE DELETE
			{`ALTER SECURITY POLICY rls.SecPol ALTER BLOCK PREDICATE rls.tenantPred_v2(TenantId) ON dbo.Sales BEFORE DELETE`, "ALTER"},
			// DROP
			{`DROP SECURITY POLICY secPolicy`, "DROP"},
			// DROP IF EXISTS
			{`DROP SECURITY POLICY IF EXISTS dbo.secPolicy`, "DROP"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityPolicyStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityPolicyStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				checkLocation(t, tt.sql, "SecurityPolicyStmt", stmt.Loc)
			})
		}
	})

	t.Run("sensitivity_classification", func(t *testing.T) {
		tests := []struct {
			sql    string
			action string
		}{
			// ADD with all options
			{`ADD SENSITIVITY CLASSIFICATION TO dbo.sales.price, dbo.sales.discount WITH (LABEL = 'Highly Confidential', INFORMATION_TYPE = 'Financial', RANK = CRITICAL)`, "ADD"},
			// ADD with label only
			{`ADD SENSITIVITY CLASSIFICATION TO dbo.customer.comments WITH (LABEL = 'Confidential', LABEL_ID = '643f7acd-776a-438d-890c-79c3f2a520d6')`, "ADD"},
			// DROP
			{`DROP SENSITIVITY CLASSIFICATION FROM dbo.sales.price`, "DROP"},
			// DROP multiple
			{`DROP SENSITIVITY CLASSIFICATION FROM dbo.sales.price, dbo.sales.discount, SalesLT.Customer.Phone`, "DROP"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SensitivityClassificationStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SensitivityClassificationStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				checkLocation(t, tt.sql, "SensitivityClassificationStmt", stmt.Loc)
			})
		}
	})

	t.Run("signature", func(t *testing.T) {
		tests := []struct {
			sql    string
			action string
		}{
			// ADD SIGNATURE with certificate
			{`ADD SIGNATURE TO HumanResources.uspUpdateEmployeeLogin BY CERTIFICATE HumanResourcesDP`, "ADD"},
			// ADD SIGNATURE with password
			{`ADD SIGNATURE TO sp_signature_demo BY CERTIFICATE cert_signature_demo WITH PASSWORD = 'pGFD4bb925DGvbd2439587y'`, "ADD"},
			// ADD COUNTER SIGNATURE
			{`ADD COUNTER SIGNATURE TO ProcSelectT1 BY CERTIFICATE csSelectT WITH PASSWORD = 'secret'`, "ADD"},
			// DROP SIGNATURE
			{`DROP SIGNATURE FROM sp_signature_demo BY CERTIFICATE cert_signature_demo`, "DROP"},
			// DROP COUNTER SIGNATURE
			{`DROP COUNTER SIGNATURE FROM sp_test BY CERTIFICATE cert1`, "DROP"},
			// ADD SIGNATURE with ASYMMETRIC KEY
			{`ADD SIGNATURE TO dbo.myProc BY ASYMMETRIC KEY myAsymKey`, "ADD"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SignatureStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SignatureStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				checkLocation(t, tt.sql, "SignatureStmt", stmt.Loc)
			})
		}
	})
}

// TestParseSpecializedIndexes tests batch 63: XML INDEX, SELECTIVE XML INDEX, SPATIAL INDEX, AGGREGATE.
func TestParseSpecializedIndexes(t *testing.T) {
	t.Run("xml_index", func(t *testing.T) {
		tests := []struct {
			sql     string
			primary bool
		}{
			// CREATE PRIMARY XML INDEX
			{`CREATE PRIMARY XML INDEX idx_xml ON dbo.Products (CatalogDescription)`, true},
			// CREATE PRIMARY XML INDEX with options
			{`CREATE PRIMARY XML INDEX PXML_ProductModel ON Production.ProductModel (CatalogDescription) WITH (PAD_INDEX = 1)`, true},
			// CREATE secondary XML INDEX FOR VALUE
			{`CREATE XML INDEX idx_xml_value ON dbo.Products (CatalogDescription) USING XML INDEX idx_xml FOR VALUE`, false},
			// CREATE secondary XML INDEX FOR PATH
			{`CREATE XML INDEX idx_xml_path ON dbo.Products (CatalogDescription) USING XML INDEX idx_xml FOR PATH`, false},
			// CREATE secondary XML INDEX FOR PROPERTY
			{`CREATE XML INDEX idx_xml_prop ON dbo.Products (CatalogDescription) USING XML INDEX idx_xml FOR PROPERTY`, false},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateXmlIndexStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *CreateXmlIndexStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Primary != tt.primary {
					t.Errorf("Parse(%q): primary = %v, want %v", tt.sql, stmt.Primary, tt.primary)
				}
				checkLocation(t, tt.sql, "CreateXmlIndexStmt", stmt.Loc)
			})
		}
	})

	t.Run("selective_xml_index", func(t *testing.T) {
		tests := []string{
			`CREATE SELECTIVE XML INDEX sxi_index ON T (xmlcol) FOR (pathab = '/a/b')`,
			`CREATE SELECTIVE XML INDEX sxi_index ON T (xmlcol) FOR (pathab = '/a/b') WITH (PAD_INDEX = 1)`,
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateSelectiveXmlIndexStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *CreateSelectiveXmlIndexStmt, got %T", sql, result.Items[0])
				}
				checkLocation(t, sql, "CreateSelectiveXmlIndexStmt", stmt.Loc)
			})
		}
	})

	t.Run("spatial_index", func(t *testing.T) {
		tests := []string{
			`CREATE SPATIAL INDEX SIndx_SpatialTable_geometry_col1 ON SpatialTable (geometry_col)`,
			`CREATE SPATIAL INDEX SIndx ON SpatialTable (geometry_col) USING GEOMETRY_AUTO_GRID WITH (CELLS_PER_OBJECT = 64)`,
			`CREATE SPATIAL INDEX SIndx ON SpatialTable (geo_col) USING GEOGRAPHY_GRID WITH (GRIDS = (LEVEL_1 = MEDIUM))`,
			`CREATE SPATIAL INDEX SIndx ON dbo.SpatialTable (geometry_col) USING GEOMETRY_GRID WITH (CELLS_PER_OBJECT = 64) ON fg1`,
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateSpatialIndexStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *CreateSpatialIndexStmt, got %T", sql, result.Items[0])
				}
				checkLocation(t, sql, "CreateSpatialIndexStmt", stmt.Loc)
			})
		}
	})

	t.Run("aggregate", func(t *testing.T) {
		tests := []struct {
			sql    string
			isDrop bool
		}{
			// CREATE AGGREGATE
			{`CREATE AGGREGATE dbo.Concatenate (@input nvarchar(4000)) RETURNS nvarchar(4000) EXTERNAL NAME StringUtilities.Concatenate`, false},
			// CREATE AGGREGATE with multiple params
			{`CREATE AGGREGATE dbo.WeightedAvg (@value float, @weight float) RETURNS float EXTERNAL NAME MyAssembly.WeightedAvg`, false},
			// DROP AGGREGATE
			{`DROP AGGREGATE dbo.Concatenate`, true},
			// DROP AGGREGATE IF EXISTS
			{`DROP AGGREGATE IF EXISTS dbo.MyAgg`, true},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				if tt.isDrop {
					stmt, ok := result.Items[0].(*ast.DropAggregateStmt)
					if !ok {
						t.Fatalf("Parse(%q): expected *DropAggregateStmt, got %T", tt.sql, result.Items[0])
					}
					checkLocation(t, tt.sql, "DropAggregateStmt", stmt.Loc)
				} else {
					stmt, ok := result.Items[0].(*ast.CreateAggregateStmt)
					if !ok {
						t.Fatalf("Parse(%q): expected *CreateAggregateStmt, got %T", tt.sql, result.Items[0])
					}
					checkLocation(t, tt.sql, "CreateAggregateStmt", stmt.Loc)
				}
			})
		}
	})
}

// TestParseRestoreExtendedMasterKey tests batch 64: RESTORE extended types, OPEN/CLOSE/BACKUP/RESTORE MASTER KEY.
func TestParseRestoreExtendedMasterKey(t *testing.T) {
	t.Run("restore_verifyonly", func(t *testing.T) {
		tests := []struct {
			sql      string
			restType string
		}{
			{`RESTORE VERIFYONLY FROM DISK = '/backup/test.bak'`, "VERIFYONLY"},
			{`RESTORE LABELONLY FROM DISK = '/backup/test.bak'`, "LABELONLY"},
			{`RESTORE REWINDONLY FROM DISK = '/backup/test.bak'`, "REWINDONLY"},
			{`RESTORE HEADERONLY FROM DISK = '/backup/test.bak'`, "HEADERONLY"},
			{`RESTORE FILELISTONLY FROM DISK = '/backup/test.bak'`, "FILELISTONLY"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.RestoreStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *RestoreStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Type != tt.restType {
					t.Errorf("Parse(%q): type = %q, want %q", tt.sql, stmt.Type, tt.restType)
				}
				checkLocation(t, tt.sql, "RestoreStmt", stmt.Loc)
			})
		}
	})

	t.Run("open_master_key", func(t *testing.T) {
		tests := []string{
			`OPEN MASTER KEY DECRYPTION BY PASSWORD = 'sfj5300osdVdgwdfkli7'`,
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityKeyStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityKeyStmt, got %T", sql, result.Items[0])
				}
				if stmt.Action != "OPEN" || stmt.ObjectType != "MASTER KEY" {
					t.Errorf("Parse(%q): action=%q objectType=%q, want OPEN/MASTER KEY", sql, stmt.Action, stmt.ObjectType)
				}
				checkLocation(t, sql, "SecurityKeyStmt", stmt.Loc)
			})
		}
	})

	t.Run("close_master_key", func(t *testing.T) {
		tests := []string{
			`CLOSE MASTER KEY`,
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityKeyStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityKeyStmt, got %T", sql, result.Items[0])
				}
				if stmt.Action != "CLOSE" || stmt.ObjectType != "MASTER KEY" {
					t.Errorf("Parse(%q): action=%q objectType=%q, want CLOSE/MASTER KEY", sql, stmt.Action, stmt.ObjectType)
				}
				checkLocation(t, sql, "SecurityKeyStmt", stmt.Loc)
			})
		}
	})

	t.Run("backup_master_key", func(t *testing.T) {
		tests := []struct {
			sql    string
			action string
		}{
			{`BACKUP MASTER KEY TO FILE = 'exportedmasterkey.bak' ENCRYPTION BY PASSWORD = 'sd092735kjn'`, "BACKUP"},
			{`RESTORE MASTER KEY FROM FILE = 'exportedmasterkey.bak' DECRYPTION BY PASSWORD = 'sd092735kjn' ENCRYPTION BY PASSWORD = 'newpassword'`, "RESTORE"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityKeyStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityKeyStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action || stmt.ObjectType != "MASTER KEY" {
					t.Errorf("Parse(%q): action=%q objectType=%q, want %s/MASTER KEY", tt.sql, stmt.Action, stmt.ObjectType, tt.action)
				}
				checkLocation(t, tt.sql, "SecurityKeyStmt", stmt.Loc)
			})
		}
	})
}


// TestParseExternalLibrary tests batch 65: CREATE/ALTER/DROP EXTERNAL LIBRARY.
func TestParseExternalLibrary(t *testing.T) {
	tests := []string{
		// CREATE basic with file path
		"CREATE EXTERNAL LIBRARY customPackage FROM (CONTENT = 'C:\\Program Files\\Microsoft SQL Server\\MSSQL14.MSSQLSERVER\\customPackage.zip') WITH (LANGUAGE = 'R')",
		// CREATE with AUTHORIZATION
		"CREATE EXTERNAL LIBRARY customLib AUTHORIZATION dbo FROM (CONTENT = 'C:\\temp\\lib.zip') WITH (LANGUAGE = 'Python')",
		// CREATE with PLATFORM
		"CREATE EXTERNAL LIBRARY customLib FROM (CONTENT = 'C:\\temp\\lib.zip', PLATFORM = WINDOWS) WITH (LANGUAGE = 'R')",
		// CREATE with varbinary literal
		"CREATE EXTERNAL LIBRARY customLib FROM (CONTENT = 0xABC123) WITH (LANGUAGE = 'R')",
		// CREATE for Java language
		"CREATE EXTERNAL LIBRARY customJar FROM (CONTENT = 'C:\\temp\\customJar.jar') WITH (LANGUAGE = 'Java')",
		// ALTER with SET
		"ALTER EXTERNAL LIBRARY customPackage SET (CONTENT = 'C:\\temp\\customPackage.zip') WITH (LANGUAGE = 'R')",
		// ALTER with AUTHORIZATION
		"ALTER EXTERNAL LIBRARY customLib AUTHORIZATION dbo SET (CONTENT = 'C:\\temp\\lib.zip') WITH (LANGUAGE = 'Python')",
		// DROP
		"DROP EXTERNAL LIBRARY customPackage",
		// DROP with AUTHORIZATION
		"DROP EXTERNAL LIBRARY customPackage AUTHORIZATION dbo",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseExternalLanguage tests batch 65: CREATE/ALTER/DROP EXTERNAL LANGUAGE.
func TestParseExternalLanguage(t *testing.T) {
	tests := []string{
		// CREATE basic
		"CREATE EXTERNAL LANGUAGE Java FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll')",
		// CREATE with AUTHORIZATION
		"CREATE EXTERNAL LANGUAGE Java AUTHORIZATION dbo FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll')",
		// CREATE with PLATFORM
		"CREATE EXTERNAL LANGUAGE Java FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll', PLATFORM = WINDOWS)",
		// CREATE with dual platform specs
		"CREATE EXTERNAL LANGUAGE Java FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll', PLATFORM = WINDOWS), (CONTENT = N'C:\\temp\\java.tar.gz', FILE_NAME = 'javaextension.so', PLATFORM = LINUX)",
		// ALTER with SET
		"ALTER EXTERNAL LANGUAGE Java SET (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll')",
		// ALTER with ADD
		"ALTER EXTERNAL LANGUAGE Java ADD (CONTENT = N'C:\\temp\\java.tar.gz', FILE_NAME = 'javaextension.so', PLATFORM = LINUX)",
		// ALTER with REMOVE PLATFORM
		"ALTER EXTERNAL LANGUAGE Java REMOVE PLATFORM LINUX",
		// DROP
		"DROP EXTERNAL LANGUAGE Java",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseCreateAvailabilityGroup tests batch 66: CREATE AVAILABILITY GROUP.
func TestParseCreateAvailabilityGroup(t *testing.T) {
	tests := []string{
		// Basic CREATE AVAILABILITY GROUP with minimal options
		"CREATE AVAILABILITY GROUP MyAG WITH (CLUSTER_TYPE = WSFC) FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
		// CREATE with multiple databases
		"CREATE AVAILABILITY GROUP MyAG FOR DATABASE db1, db2, db3 REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL)",
		// CREATE with multiple WITH options
		"CREATE AVAILABILITY GROUP MyAG WITH (AUTOMATED_BACKUP_PREFERENCE = SECONDARY, FAILURE_CONDITION_LEVEL = 3, HEALTH_CHECK_TIMEOUT = 30000, DB_FAILOVER = ON) FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
		// CREATE with LISTENER
		"CREATE AVAILABILITY GROUP MyAG FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC) LISTENER 'MyAGListener' (WITH IP ((N'10.120.19.155', N'255.255.254.0')), PORT = 1433)",
		// CREATE with CLUSTER_TYPE = NONE (Linux)
		"CREATE AVAILABILITY GROUP MyAG WITH (CLUSTER_TYPE = NONE) FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = EXTERNAL)",
		// CREATE with DTC_SUPPORT
		"CREATE AVAILABILITY GROUP MyAG WITH (DTC_SUPPORT = PER_DB) FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
		// CREATE with REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT
		"CREATE AVAILABILITY GROUP MyAG WITH (REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT = 1) FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
		// CREATE with SEEDING_MODE, BACKUP_PRIORITY in replica
		"CREATE AVAILABILITY GROUP MyAG FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC, SEEDING_MODE = AUTOMATIC, BACKUP_PRIORITY = 50)",
		// CREATE with SECONDARY_ROLE and PRIMARY_ROLE
		"CREATE AVAILABILITY GROUP MyAG FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC, SECONDARY_ROLE (ALLOW_CONNECTIONS = READ_ONLY), PRIMARY_ROLE (ALLOW_CONNECTIONS = ALL))",
		// CREATE DISTRIBUTED
		"CREATE AVAILABILITY GROUP MyDistAG WITH (DISTRIBUTED) AVAILABILITY GROUP ON 'AG1' WITH (LISTENER_URL = 'TCP://server1:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL, SEEDING_MODE = AUTOMATIC), 'AG2' WITH (LISTENER_URL = 'TCP://server2:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL, SEEDING_MODE = AUTOMATIC)",
		// CREATE with CONFIGURATION_ONLY availability mode
		"CREATE AVAILABILITY GROUP MyAG FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = CONFIGURATION_ONLY, FAILOVER_MODE = MANUAL)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Action != "CREATE" {
				t.Errorf("Parse(%q): Action = %q, want CREATE", sql, stmt.Action)
			}
			if stmt.ObjectType != "AVAILABILITY GROUP" {
				t.Errorf("Parse(%q): ObjectType = %q, want AVAILABILITY GROUP", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseAlterAvailabilityGroup tests batch 66: ALTER AVAILABILITY GROUP.
func TestParseAlterAvailabilityGroup(t *testing.T) {
	tests := []string{
		// SET option
		"ALTER AVAILABILITY GROUP MyAG SET (AUTOMATED_BACKUP_PREFERENCE = SECONDARY)",
		// ADD DATABASE
		"ALTER AVAILABILITY GROUP MyAG ADD DATABASE MyNewDB",
		// REMOVE DATABASE
		"ALTER AVAILABILITY GROUP MyAG REMOVE DATABASE MyOldDB",
		// ADD REPLICA ON
		"ALTER AVAILABILITY GROUP MyAG ADD REPLICA ON 'server2' WITH (ENDPOINT_URL = 'TCP://server2:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
		// MODIFY REPLICA ON
		"ALTER AVAILABILITY GROUP MyAG MODIFY REPLICA ON 'server2' WITH (AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT)",
		// REMOVE REPLICA ON
		"ALTER AVAILABILITY GROUP MyAG REMOVE REPLICA ON 'server2'",
		// JOIN
		"ALTER AVAILABILITY GROUP MyAG JOIN",
		// JOIN AVAILABILITY GROUP ON (distributed)
		"ALTER AVAILABILITY GROUP MyDistAG JOIN AVAILABILITY GROUP ON 'AG1' WITH (LISTENER_URL = 'TCP://server1:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL, SEEDING_MODE = AUTOMATIC), 'AG2' WITH (LISTENER_URL = 'TCP://server2:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL, SEEDING_MODE = AUTOMATIC)",
		// GRANT CREATE ANY DATABASE
		"ALTER AVAILABILITY GROUP MyAG GRANT CREATE ANY DATABASE",
		// DENY CREATE ANY DATABASE
		"ALTER AVAILABILITY GROUP MyAG DENY CREATE ANY DATABASE",
		// FAILOVER
		"ALTER AVAILABILITY GROUP MyAG FAILOVER",
		// FORCE_FAILOVER_ALLOW_DATA_LOSS
		"ALTER AVAILABILITY GROUP MyAG FORCE_FAILOVER_ALLOW_DATA_LOSS",
		// ADD LISTENER
		"ALTER AVAILABILITY GROUP MyAG ADD LISTENER 'MyAGListener' (WITH IP ((N'10.120.19.155', N'255.255.254.0')), PORT = 1433)",
		// MODIFY LISTENER
		"ALTER AVAILABILITY GROUP MyAG MODIFY LISTENER 'MyAGListener' (ADD IP (N'10.120.19.200', N'255.255.254.0'))",
		// RESTART LISTENER
		"ALTER AVAILABILITY GROUP MyAG RESTART LISTENER 'MyAGListener'",
		// REMOVE LISTENER
		"ALTER AVAILABILITY GROUP MyAG REMOVE LISTENER 'MyAGListener'",
		// OFFLINE
		"ALTER AVAILABILITY GROUP MyAG OFFLINE",
		// SET multiple options
		"ALTER AVAILABILITY GROUP MyAG SET (DB_FAILOVER = ON, FAILURE_CONDITION_LEVEL = 4, HEALTH_CHECK_TIMEOUT = 60000)",
		// MODIFY AVAILABILITY GROUP ON (distributed)
		"ALTER AVAILABILITY GROUP MyDistAG MODIFY AVAILABILITY GROUP ON 'AG1' WITH (LISTENER_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Action != "ALTER" {
				t.Errorf("Parse(%q): Action = %q, want ALTER", sql, stmt.Action)
			}
			if stmt.ObjectType != "AVAILABILITY GROUP" {
				t.Errorf("Parse(%q): ObjectType = %q, want AVAILABILITY GROUP", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseDropAvailabilityGroup tests batch 66: DROP AVAILABILITY GROUP.
func TestParseDropAvailabilityGroup(t *testing.T) {
	tests := []string{
		"DROP AVAILABILITY GROUP MyAG",
		"DROP AVAILABILITY GROUP [MyAG]",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Action != "DROP" {
				t.Errorf("Parse(%q): Action = %q, want DROP", sql, stmt.Action)
			}
			if stmt.ObjectType != "AVAILABILITY GROUP" {
				t.Errorf("Parse(%q): ObjectType = %q, want AVAILABILITY GROUP", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseCreateDefault tests batch 67: CREATE DEFAULT.
func TestParseCreateDefault(t *testing.T) {
	tests := []string{
		// Basic CREATE DEFAULT with string constant
		"CREATE DEFAULT phonedflt AS 'unknown'",
		// CREATE DEFAULT with numeric constant
		"CREATE DEFAULT zero_default AS 0",
		// CREATE DEFAULT with schema-qualified name
		"CREATE DEFAULT dbo.datedflt AS GETDATE()",
		// CREATE DEFAULT with N-string
		"CREATE DEFAULT ndflt AS N'unknown'",
		// CREATE DEFAULT with negative number
		"CREATE DEFAULT neg_default AS -1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Action != "CREATE" {
				t.Errorf("Parse(%q): Action = %q, want CREATE", sql, stmt.Action)
			}
			if stmt.ObjectType != "DEFAULT" {
				t.Errorf("Parse(%q): ObjectType = %q, want DEFAULT", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseCreateRule tests batch 67: CREATE RULE.
func TestParseCreateRule(t *testing.T) {
	tests := []string{
		// Range rule
		"CREATE RULE range_rule AS @range >= 1000 AND @range < 20000",
		// List rule
		"CREATE RULE list_rule AS @list IN ('1389', '0736', '0877')",
		// Pattern rule
		"CREATE RULE pattern_rule AS @value LIKE '__-%[0-9]'",
		// Schema-qualified name
		"CREATE RULE dbo.positive_rule AS @val > 0",
		// Simple comparison
		"CREATE RULE min_rule AS @val >= 0",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Action != "CREATE" {
				t.Errorf("Parse(%q): Action = %q, want CREATE", sql, stmt.Action)
			}
			if stmt.ObjectType != "RULE" {
				t.Errorf("Parse(%q): ObjectType = %q, want RULE", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseAlterDatabaseScopedConfig tests batch 67: ALTER DATABASE SCOPED CONFIGURATION.
func TestParseAlterDatabaseScopedConfig(t *testing.T) {
	tests := []string{
		// SET MAXDOP
		"ALTER DATABASE SCOPED CONFIGURATION SET MAXDOP = 4",
		// SET ON/OFF option
		"ALTER DATABASE SCOPED CONFIGURATION SET LEGACY_CARDINALITY_ESTIMATION = ON",
		// SET with PRIMARY value
		"ALTER DATABASE SCOPED CONFIGURATION SET MAXDOP = PRIMARY",
		// SET PARAMETER_SNIFFING
		"ALTER DATABASE SCOPED CONFIGURATION SET PARAMETER_SNIFFING = OFF",
		// SET QUERY_OPTIMIZER_HOTFIXES
		"ALTER DATABASE SCOPED CONFIGURATION SET QUERY_OPTIMIZER_HOTFIXES = ON",
		// FOR SECONDARY SET
		"ALTER DATABASE SCOPED CONFIGURATION FOR SECONDARY SET MAXDOP = PRIMARY",
		// CLEAR PROCEDURE_CACHE
		"ALTER DATABASE SCOPED CONFIGURATION CLEAR PROCEDURE_CACHE",
		// SET IDENTITY_CACHE
		"ALTER DATABASE SCOPED CONFIGURATION SET IDENTITY_CACHE = OFF",
		// SET OPTIMIZE_FOR_AD_HOC_WORKLOADS
		"ALTER DATABASE SCOPED CONFIGURATION SET OPTIMIZE_FOR_AD_HOC_WORKLOADS = ON",
		// SET ELEVATE_ONLINE
		"ALTER DATABASE SCOPED CONFIGURATION SET ELEVATE_ONLINE = WHEN_SUPPORTED",
		// SET ELEVATE_RESUMABLE
		"ALTER DATABASE SCOPED CONFIGURATION SET ELEVATE_RESUMABLE = FAIL_UNSUPPORTED",
		// SET BATCH_MODE_ON_ROWSTORE
		"ALTER DATABASE SCOPED CONFIGURATION SET BATCH_MODE_ON_ROWSTORE = ON",
		// SET PAUSED_RESUMABLE_INDEX_ABORT_DURATION_MINUTES
		"ALTER DATABASE SCOPED CONFIGURATION SET PAUSED_RESUMABLE_INDEX_ABORT_DURATION_MINUTES = 60",
		// SET VERBOSE_TRUNCATION_WARNINGS
		"ALTER DATABASE SCOPED CONFIGURATION SET VERBOSE_TRUNCATION_WARNINGS = ON",
		// SET LAST_QUERY_PLAN_STATS
		"ALTER DATABASE SCOPED CONFIGURATION SET LAST_QUERY_PLAN_STATS = ON",
		// SET LIGHTWEIGHT_QUERY_PROFILING
		"ALTER DATABASE SCOPED CONFIGURATION SET LIGHTWEIGHT_QUERY_PROFILING = OFF",
		// CLEAR PROCEDURE_CACHE with plan handle
		"ALTER DATABASE SCOPED CONFIGURATION CLEAR PROCEDURE_CACHE 0x060006001ECA270E",
		// SET LEDGER_DIGEST_STORAGE_ENDPOINT OFF
		"ALTER DATABASE SCOPED CONFIGURATION SET LEDGER_DIGEST_STORAGE_ENDPOINT = OFF",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Action != "ALTER" {
				t.Errorf("Parse(%q): Action = %q, want ALTER", sql, stmt.Action)
			}
			if stmt.ObjectType != "DATABASE SCOPED CONFIGURATION" {
				t.Errorf("Parse(%q): ObjectType = %q, want DATABASE SCOPED CONFIGURATION", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseIntegrationPhase3 tests batch 68: integration test for all phase 3 statement types.
// Verifies that the main dispatcher correctly routes to parsers from batches 50-67.
func TestParseIntegrationPhase3(t *testing.T) {
	tests := []string{
		// Batch 50: Service Broker - CREATE ROUTE, CREATE REMOTE SERVICE BINDING
		"CREATE ROUTE ExpenseRoute WITH SERVICE_NAME = 'ExpenseService', ADDRESS = 'TCP://expense.example.com:4022'",
		"CREATE REMOTE SERVICE BINDING MyBinding TO SERVICE 'TargetService' WITH USER = MyUser",
		// Batch 51: ENABLE/DISABLE TRIGGER, ALTER TRIGGER
		"ENABLE TRIGGER trg1 ON dbo.Orders",
		"DISABLE TRIGGER ALL ON DATABASE",
		// Batch 52: Service Broker ALTER/DROP, MOVE CONVERSATION, BROKER PRIORITY
		"CREATE BROKER PRIORITY BrokerPriority1 FOR CONVERSATION SET (CONTRACT_NAME = MyContract, PRIORITY_LEVEL = 5)",
		// Batch 53: Security context
		"EXECUTE AS USER = 'testuser'",
		"REVERT",
		"ALTER AUTHORIZATION ON OBJECT::dbo.MyTable TO NewOwner",
		// Batch 54: Server audit
		"CREATE SERVER AUDIT MyAudit TO FILE (FILEPATH = 'C:\\AuditLogs')",
		// Batch 55: Endpoint
		"CREATE ENDPOINT MyEndpoint STATE = STARTED AS TCP (LISTENER_PORT = 4022) FOR TSQL()",
		// Batch 56: Event notification/session
		"CREATE EVENT SESSION MySession ON SERVER ADD EVENT sqlserver.sql_statement_completed ADD TARGET package0.event_file(SET filename = 'audit.xel')",
		// Batch 57: External objects
		"CREATE EXTERNAL DATA SOURCE MySource WITH (LOCATION = 'hdfs://10.10.10.10:8020')",
		"DROP EXTERNAL DATA SOURCE MySource",
		"DROP EXTERNAL LIBRARY MyLib",
		"DROP EXTERNAL LANGUAGE MyLang",
		"DROP EXTERNAL RESOURCE POOL MyPool",
		// Batch 58: Encryption keys
		"CREATE COLUMN MASTER KEY MyCMK WITH (KEY_STORE_PROVIDER_NAME = 'MSSQL_CERTIFICATE_STORE')",
		// Batch 59: Resource governor
		"CREATE WORKLOAD GROUP TestGroup USING [default]",
		"ALTER RESOURCE GOVERNOR RECONFIGURE",
		// Batch 60: Server-level objects
		"CREATE SERVER ROLE MyServerRole",
		// Batch 61: Fulltext extensions
		"CREATE FULLTEXT STOPLIST MyStoplist",
		// Batch 62: Security policy, classification, signature
		"CREATE SECURITY POLICY dbo.SecurityFilter ADD FILTER PREDICATE dbo.fn_securitypredicate(TenantId) ON dbo.Sales",
		// Batch 63: Specialized indexes, aggregate
		"CREATE SPATIAL INDEX SIndx_SpatialTable ON dbo.SpatialTable (geometry_col)",
		// Batch 64: Extended restore, master key ops
		"OPEN MASTER KEY DECRYPTION BY PASSWORD = 'password123'",
		"CLOSE MASTER KEY",
		// Batch 65: External library/language
		"CREATE EXTERNAL LIBRARY MyLib FROM (CONTENT = 0x504B) WITH (LANGUAGE = 'R')",
		// Batch 66: Availability group
		"CREATE AVAILABILITY GROUP MyAG FOR DATABASE MyDB REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
		"ALTER AVAILABILITY GROUP MyAG FAILOVER",
		"DROP AVAILABILITY GROUP MyAG",
		// Batch 67: Deprecated misc
		"CREATE DEFAULT phonedflt AS 'unknown'",
		"CREATE RULE range_rule AS @range >= 1000 AND @range < 20000",
		"ALTER DATABASE SCOPED CONFIGURATION SET MAXDOP = 4",
		// BEGIN DIALOG (service broker, newly wired)
		"BEGIN DIALOG @dialog FROM SERVICE 'InitiatorService' TO SERVICE 'TargetService' ON CONTRACT 'MyContract'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() < 1 {
				t.Fatalf("Parse(%q): got %d statements, want >= 1", sql, result.Len())
			}
		})
	}
}

// TestParseCreateTableDepth tests batch 69: temporal table support, advanced column options, table storage.
func TestParseCreateTableDepth(t *testing.T) {
	t.Run("temporal table with system versioning", func(t *testing.T) {
		sql := `CREATE TABLE dbo.Employee (
			EmployeeID int PRIMARY KEY,
			Name nvarchar(100),
			Salary decimal(10,2),
			ValidFrom datetime2 GENERATED ALWAYS AS ROW START HIDDEN NOT NULL,
			ValidTo datetime2 GENERATED ALWAYS AS ROW END HIDDEN NOT NULL,
			PERIOD FOR SYSTEM_TIME (ValidFrom, ValidTo)
		)
		WITH (SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.EmployeeHistory))`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.PeriodStartCol == "" {
			t.Error("expected PeriodStartCol to be set")
		}
		if stmt.TableOptions == nil {
			t.Error("expected TableOptions to be set")
		}
	})

	t.Run("sparse column", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			optional_data varchar(100) SPARSE NULL
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col := stmt.Columns.Items[1].(*ast.ColumnDef)
		if !col.Sparse {
			t.Error("expected Sparse to be true")
		}
	})

	t.Run("masked column", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			email varchar(100) MASKED WITH (FUNCTION = 'email()') NULL,
			ssn varchar(11) MASKED WITH (FUNCTION = 'partial(0,"XXX-XX-",4)') NULL
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col := stmt.Columns.Items[1].(*ast.ColumnDef)
		if col.MaskFunction != "email()" {
			t.Errorf("expected MaskFunction 'email()', got %q", col.MaskFunction)
		}
	})

	t.Run("encrypted column", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			ssn varchar(11) ENCRYPTED WITH (
				COLUMN_ENCRYPTION_KEY = MyCEK,
				ENCRYPTION_TYPE = DETERMINISTIC,
				ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'
			)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col := stmt.Columns.Items[1].(*ast.ColumnDef)
		if col.EncryptedWith == nil {
			t.Error("expected EncryptedWith to be set")
		} else {
			if col.EncryptedWith.EncryptionType != "DETERMINISTIC" {
				t.Errorf("expected DETERMINISTIC, got %q", col.EncryptedWith.EncryptionType)
			}
		}
	})

	t.Run("memory optimized table", func(t *testing.T) {
		sql := `CREATE TABLE dbo.MemTable (
			Id int NOT NULL PRIMARY KEY NONCLUSTERED,
			Name nvarchar(100) NOT NULL
		) WITH (MEMORY_OPTIMIZED = ON, DURABILITY = SCHEMA_AND_DATA)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.TableOptions == nil || stmt.TableOptions.Len() != 2 {
			t.Errorf("expected 2 table options, got %v", stmt.TableOptions)
		}
	})

	t.Run("filestream column", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id uniqueidentifier ROWGUIDCOL NOT NULL,
			doc varbinary(max) FILESTREAM NULL
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col0 := stmt.Columns.Items[0].(*ast.ColumnDef)
		if !col0.Rowguidcol {
			t.Error("expected Rowguidcol to be true")
		}
		col1 := stmt.Columns.Items[1].(*ast.ColumnDef)
		if !col1.Filestream {
			t.Error("expected Filestream to be true")
		}
	})

	t.Run("on filegroup", func(t *testing.T) {
		sql := `CREATE TABLE t (id int NOT NULL) ON MyFilegroup`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.OnFilegroup != "MyFilegroup" {
			t.Errorf("expected OnFilegroup 'MyFilegroup', got %q", stmt.OnFilegroup)
		}
	})

	t.Run("textimage_on", func(t *testing.T) {
		sql := `CREATE TABLE t (id int NOT NULL, doc text) ON PRIMARY TEXTIMAGE_ON LobGroup`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if !strings.EqualFold(stmt.OnFilegroup, "PRIMARY") {
			t.Errorf("expected OnFilegroup 'PRIMARY', got %q", stmt.OnFilegroup)
		}
		if !strings.EqualFold(stmt.TextImageOn, "LobGroup") {
			t.Errorf("expected TextImageOn 'LobGroup', got %q", stmt.TextImageOn)
		}
	})

	t.Run("generated always as row start end", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			start_time datetime2 GENERATED ALWAYS AS ROW START NOT NULL,
			end_time datetime2 GENERATED ALWAYS AS ROW END NOT NULL,
			PERIOD FOR SYSTEM_TIME (start_time, end_time)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col1 := stmt.Columns.Items[1].(*ast.ColumnDef)
		if col1.GeneratedAlways == nil {
			t.Fatal("expected GeneratedAlways to be set")
		}
		if col1.GeneratedAlways.Kind != "ROW" || col1.GeneratedAlways.StartEnd != "START" {
			t.Errorf("expected ROW START, got %s %s", col1.GeneratedAlways.Kind, col1.GeneratedAlways.StartEnd)
		}
		if stmt.PeriodStartCol != "start_time" || stmt.PeriodEndCol != "end_time" {
			t.Errorf("expected period(start_time, end_time), got (%s, %s)", stmt.PeriodStartCol, stmt.PeriodEndCol)
		}
	})

	t.Run("hidden column", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			ts datetime2 GENERATED ALWAYS AS ROW START HIDDEN NOT NULL,
			te datetime2 GENERATED ALWAYS AS ROW END HIDDEN NOT NULL,
			PERIOD FOR SYSTEM_TIME (ts, te)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col1 := stmt.Columns.Items[1].(*ast.ColumnDef)
		if !col1.Hidden {
			t.Error("expected Hidden to be true")
		}
	})

	t.Run("not for replication", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int IDENTITY(1,1) NOT FOR REPLICATION NOT NULL
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		col := stmt.Columns.Items[0].(*ast.ColumnDef)
		if !col.NotForReplication {
			t.Error("expected NotForReplication to be true")
		}
	})

	t.Run("system versioning with options", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id int NOT NULL,
			s datetime2 GENERATED ALWAYS AS ROW START NOT NULL,
			e datetime2 GENERATED ALWAYS AS ROW END NOT NULL,
			PERIOD FOR SYSTEM_TIME (s, e)
		) WITH (SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.tHistory, DATA_CONSISTENCY_CHECK = ON))`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.TableOptions == nil {
			t.Fatal("expected TableOptions")
		}
		opt := stmt.TableOptions.Items[0].(*ast.TableOption)
		if opt.HistoryTable != "dbo.tHistory" {
			t.Errorf("expected dbo.tHistory, got %q", opt.HistoryTable)
		}
		if opt.DataConsistencyCheck != "ON" {
			t.Errorf("expected ON, got %q", opt.DataConsistencyCheck)
		}
	})

	t.Run("data compression", func(t *testing.T) {
		sql := `CREATE TABLE t (id int NOT NULL) WITH (DATA_COMPRESSION = PAGE)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.TableOptions == nil {
			t.Fatal("expected TableOptions")
		}
		opt := stmt.TableOptions.Items[0].(*ast.TableOption)
		if opt.Name != "DATA_COMPRESSION" || opt.Value != "PAGE" {
			t.Errorf("expected DATA_COMPRESSION=PAGE, got %s=%s", opt.Name, opt.Value)
		}
	})
}

// TestParseAlterTableDepth tests ALTER TABLE extended features (batch 70).
func TestParseAlterTableDepth(t *testing.T) {
	// Multiple comma-separated ADD actions
	t.Run("alter_table_multiple_actions", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t ADD col1 int, col2 varchar(50)",
			"ALTER TABLE t ADD col1 int NULL, CONSTRAINT PK_t PRIMARY KEY (id)",
			"ALTER TABLE t ADD col1 int, col2 nvarchar(100) NOT NULL, col3 datetime DEFAULT GETDATE()",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				if stmt.Actions == nil || stmt.Actions.Len() < 2 {
					t.Fatalf("expected multiple actions, got %v", stmt.Actions)
				}
			})
		}
	})

	// SWITCH PARTITION
	t.Run("alter_table_switch_partition", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t SWITCH TO t2",
			"ALTER TABLE t SWITCH PARTITION 1 TO t2",
			"ALTER TABLE t SWITCH PARTITION 1 TO t2 PARTITION 2",
			"ALTER TABLE dbo.source SWITCH PARTITION 5 TO dbo.target PARTITION 5",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != ast.ATSwitchPartition {
					t.Errorf("expected ATSwitchPartition, got %d", action.Type)
				}
				if action.TargetName == nil {
					t.Fatal("expected target table name")
				}
			})
		}
	})

	// CHECK/NOCHECK CONSTRAINT
	t.Run("alter_table_check_nocheck", func(t *testing.T) {
		tests := []struct {
			sql     string
			actType ast.AlterTableActionType
		}{
			{"ALTER TABLE t CHECK CONSTRAINT ALL", ast.ATCheckConstraint},
			{"ALTER TABLE t NOCHECK CONSTRAINT ALL", ast.ATNocheckConstraint},
			{"ALTER TABLE t CHECK CONSTRAINT FK_t_col", ast.ATCheckConstraint},
			{"ALTER TABLE t NOCHECK CONSTRAINT FK_t_col", ast.ATNocheckConstraint},
			{"ALTER TABLE t WITH CHECK CHECK CONSTRAINT ALL", ast.ATCheckConstraint},
			{"ALTER TABLE t WITH NOCHECK NOCHECK CONSTRAINT FK_t_col", ast.ATNocheckConstraint},
			{"ALTER TABLE t CHECK CONSTRAINT ck1, ck2", ast.ATCheckConstraint},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != tc.actType {
					t.Errorf("expected action type %d, got %d", tc.actType, action.Type)
				}
			})
		}
	})

	// ENABLE/DISABLE TRIGGER (at ALTER TABLE level)
	t.Run("alter_table_enable_disable_trigger", func(t *testing.T) {
		tests := []struct {
			sql     string
			actType ast.AlterTableActionType
		}{
			{"ALTER TABLE t ENABLE TRIGGER ALL", ast.ATEnableTrigger},
			{"ALTER TABLE t DISABLE TRIGGER ALL", ast.ATDisableTrigger},
			{"ALTER TABLE t ENABLE TRIGGER trg1", ast.ATEnableTrigger},
			{"ALTER TABLE t DISABLE TRIGGER trg1, trg2", ast.ATDisableTrigger},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != tc.actType {
					t.Errorf("expected action type %d, got %d", tc.actType, action.Type)
				}
			})
		}
	})

	// ENABLE/DISABLE CHANGE_TRACKING
	t.Run("alter_table_change_tracking", func(t *testing.T) {
		tests := []struct {
			sql     string
			actType ast.AlterTableActionType
		}{
			{"ALTER TABLE t ENABLE CHANGE_TRACKING", ast.ATEnableChangeTracking},
			{"ALTER TABLE t DISABLE CHANGE_TRACKING", ast.ATDisableChangeTracking},
			{"ALTER TABLE t ENABLE CHANGE_TRACKING WITH (TRACK_COLUMNS_UPDATED = ON)", ast.ATEnableChangeTracking},
		}
		for _, tc := range tests {
			t.Run(tc.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tc.sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != tc.actType {
					t.Errorf("expected action type %d, got %d", tc.actType, action.Type)
				}
			})
		}
	})

	// REBUILD
	t.Run("alter_table_rebuild", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t REBUILD",
			"ALTER TABLE t REBUILD PARTITION = ALL",
			"ALTER TABLE t REBUILD PARTITION = 1",
			"ALTER TABLE t REBUILD WITH (DATA_COMPRESSION = PAGE)",
			"ALTER TABLE t REBUILD PARTITION = ALL WITH (DATA_COMPRESSION = ROW)",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != ast.ATRebuild {
					t.Errorf("expected ATRebuild, got %d", action.Type)
				}
			})
		}
	})

	// SET options
	t.Run("alter_table_set", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t SET (LOCK_ESCALATION = TABLE)",
			"ALTER TABLE t SET (LOCK_ESCALATION = AUTO)",
			"ALTER TABLE t SET (LOCK_ESCALATION = DISABLE)",
			"ALTER TABLE t SET (FILESTREAM_ON = myFilegroup)",
			"ALTER TABLE t SET (SYSTEM_VERSIONING = ON)",
			"ALTER TABLE t SET (SYSTEM_VERSIONING = OFF)",
			"ALTER TABLE t SET (SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.history))",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != ast.ATSet {
					t.Errorf("expected ATSet, got %d", action.Type)
				}
			})
		}
	})

	// ALTER COLUMN with COLLATE
	t.Run("alter_column_collate", func(t *testing.T) {
		sql := "ALTER TABLE t ALTER COLUMN col1 varchar(100) COLLATE Latin1_General_CI_AS NOT NULL"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATAlterColumn {
			t.Errorf("expected ATAlterColumn, got %d", action.Type)
		}
		if action.Collation != "Latin1_General_CI_AS" {
			t.Errorf("expected collation Latin1_General_CI_AS, got %s", action.Collation)
		}
	})

	// ALTER COLUMN ADD/DROP ROWGUIDCOL, PERSISTED, SPARSE, HIDDEN, MASKED
	t.Run("alter_column_add_drop", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t ALTER COLUMN col1 ADD ROWGUIDCOL",
			"ALTER TABLE t ALTER COLUMN col1 DROP ROWGUIDCOL",
			"ALTER TABLE t ALTER COLUMN col1 ADD PERSISTED",
			"ALTER TABLE t ALTER COLUMN col1 DROP PERSISTED",
			"ALTER TABLE t ALTER COLUMN col1 ADD SPARSE",
			"ALTER TABLE t ALTER COLUMN col1 DROP SPARSE",
			"ALTER TABLE t ALTER COLUMN col1 ADD HIDDEN",
			"ALTER TABLE t ALTER COLUMN col1 DROP HIDDEN",
			"ALTER TABLE t ALTER COLUMN col1 ADD MASKED WITH (FUNCTION = 'default()')",
			"ALTER TABLE t ALTER COLUMN col1 DROP MASKED",
			"ALTER TABLE t ALTER COLUMN col1 ADD NOT FOR REPLICATION",
			"ALTER TABLE t ALTER COLUMN col1 DROP NOT FOR REPLICATION",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.AlterTableStmt)
				action := stmt.Actions.Items[0].(*ast.AlterTableAction)
				if action.Type != ast.ATAlterColumnAddDrop {
					t.Errorf("expected ATAlterColumnAddDrop, got %d", action.Type)
				}
			})
		}
	})

	// DROP with IF EXISTS
	t.Run("alter_table_drop_if_exists", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t DROP CONSTRAINT IF EXISTS ck1",
			"ALTER TABLE t DROP COLUMN IF EXISTS col1",
			"ALTER TABLE t DROP COLUMN IF EXISTS col1, col2",
			"ALTER TABLE t DROP CONSTRAINT IF EXISTS ck1, ck2",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				ParseAndCheck(t, sql)
			})
		}
	})

	// Multiple DROP items
	t.Run("alter_table_drop_multiple", func(t *testing.T) {
		sql := "ALTER TABLE t DROP COLUMN col1, col2, col3"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		// Multiple drops should produce multiple actions
		if stmt.Actions == nil || stmt.Actions.Len() < 2 {
			t.Fatalf("expected multiple drop actions, got %v", stmt.Actions)
		}
	})
}

// TestParseRoutineOptionsDepth tests WITH options for proc/func/view (batch 71).
func TestParseRoutineOptionsDepth(t *testing.T) {
	// CREATE PROCEDURE WITH options
	t.Run("create_proc_with_options", func(t *testing.T) {
		tests := []string{
			"CREATE PROCEDURE dbo.p1 WITH RECOMPILE AS SELECT 1",
			"CREATE PROCEDURE dbo.p1 WITH ENCRYPTION AS SELECT 1",
			"CREATE PROCEDURE dbo.p1 WITH EXECUTE AS OWNER AS SELECT 1",
			"CREATE PROCEDURE dbo.p1 WITH RECOMPILE, ENCRYPTION AS SELECT 1",
			"CREATE PROCEDURE dbo.p1 WITH EXECUTE AS 'dbo' AS SELECT 1",
			"CREATE OR ALTER PROCEDURE dbo.p1 WITH RECOMPILE AS SELECT 1",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.CreateProcedureStmt)
				if stmt.Options == nil || stmt.Options.Len() == 0 {
					t.Fatal("expected WITH options")
				}
			})
		}
	})

	// CREATE FUNCTION WITH options
	t.Run("create_func_with_options", func(t *testing.T) {
		tests := []string{
			"CREATE FUNCTION dbo.f1() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN 1 END",
			"CREATE FUNCTION dbo.f1() RETURNS int WITH RETURNS NULL ON NULL INPUT AS BEGIN RETURN 1 END",
			"CREATE FUNCTION dbo.f1() RETURNS int WITH CALLED ON NULL INPUT AS BEGIN RETURN 1 END",
			"CREATE FUNCTION dbo.f1() RETURNS int WITH EXECUTE AS CALLER AS BEGIN RETURN 1 END",
			"CREATE FUNCTION dbo.f1() RETURNS int WITH SCHEMABINDING, RETURNS NULL ON NULL INPUT AS BEGIN RETURN 1 END",
			"CREATE FUNCTION dbo.f1() RETURNS int WITH ENCRYPTION AS BEGIN RETURN 1 END",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.CreateFunctionStmt)
				if stmt.Options == nil || stmt.Options.Len() == 0 {
					t.Fatal("expected WITH options")
				}
			})
		}
	})

	// CREATE VIEW WITH options
	t.Run("create_view_with_options", func(t *testing.T) {
		tests := []string{
			"CREATE VIEW dbo.v1 WITH SCHEMABINDING AS SELECT 1 AS a",
			"CREATE VIEW dbo.v1 WITH VIEW_METADATA AS SELECT 1 AS a",
			"CREATE VIEW dbo.v1 WITH ENCRYPTION AS SELECT 1 AS a",
			"CREATE VIEW dbo.v1 WITH SCHEMABINDING, VIEW_METADATA AS SELECT 1 AS a",
			"CREATE VIEW dbo.v1 WITH SCHEMABINDING, ENCRYPTION, VIEW_METADATA AS SELECT 1 AS a",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				stmt := result.Items[0].(*ast.CreateViewStmt)
				if stmt.Options == nil || stmt.Options.Len() == 0 {
					t.Fatal("expected WITH options")
				}
			})
		}
	})

	// Existing tests should still pass (no WITH)
	t.Run("proc_no_options", func(t *testing.T) {
		sql := "CREATE PROCEDURE dbo.p1 AS SELECT 1"
		result := ParseAndCheck(t, sql)
		_ = result.Items[0].(*ast.CreateProcedureStmt)
	})

	t.Run("func_no_options", func(t *testing.T) {
		sql := "CREATE FUNCTION dbo.f1() RETURNS int AS BEGIN RETURN 1 END"
		result := ParseAndCheck(t, sql)
		_ = result.Items[0].(*ast.CreateFunctionStmt)
	})

	t.Run("view_no_options", func(t *testing.T) {
		sql := "CREATE VIEW dbo.v1 AS SELECT 1 AS a"
		result := ParseAndCheck(t, sql)
		_ = result.Items[0].(*ast.CreateViewStmt)
	})
}

// TestParseServiceBrokerDepth tests batch 72: structural parsing of Service Broker statements.
func TestParseServiceBrokerDepth(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// BEGIN DIALOG CONVERSATION - full form
		{
			name: "begin_dialog_full",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target' ON CONTRACT [//MyApp/Contract]",
		},
		{
			name: "begin_dialog_with_broker_guid",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target', 'a326e034-d4cf-4e8b-8d98-4d7e1926c904' ON CONTRACT [//MyApp/Contract]",
		},
		{
			name: "begin_dialog_with_current_database",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target', 'CURRENT DATABASE' ON CONTRACT [//MyApp/Contract]",
		},
		{
			name: "begin_dialog_with_lifetime",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target' ON CONTRACT [//MyApp/Contract] WITH LIFETIME = 60",
		},
		{
			name: "begin_dialog_with_encryption",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target' ON CONTRACT [//MyApp/Contract] WITH ENCRYPTION = OFF",
		},
		{
			name: "begin_dialog_with_related_conversation",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target' ON CONTRACT [//MyApp/Contract] WITH RELATED_CONVERSATION = @existing_handle",
		},
		{
			name: "begin_dialog_with_related_group",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target' ON CONTRACT [//MyApp/Contract] WITH RELATED_CONVERSATION_GROUP = @group_id",
		},
		{
			name: "begin_dialog_with_multiple_options",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target' ON CONTRACT [//MyApp/Contract] WITH RELATED_CONVERSATION = @existing_handle, LIFETIME = 600, ENCRYPTION = ON",
		},
		{
			name: "begin_dialog_no_contract",
			sql:  "BEGIN DIALOG CONVERSATION @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target'",
		},
		{
			name: "begin_dialog_minimal",
			sql:  "BEGIN DIALOG @dialog_handle FROM SERVICE [//MyApp/Initiator] TO SERVICE '//MyApp/Target'",
		},
		// RECEIVE - column list and INTO
		{
			name: "receive_star",
			sql:  "RECEIVE * FROM ExpenseQueue",
		},
		{
			name: "receive_columns",
			sql:  "RECEIVE conversation_handle, message_type_name, message_body FROM ExpenseQueue",
		},
		{
			name: "receive_top",
			sql:  "RECEIVE TOP (1) * FROM ExpenseQueue",
		},
		{
			name: "receive_into",
			sql:  "RECEIVE TOP (1) conversation_handle, message_body FROM ExpenseQueue INTO @tableVar",
		},
		{
			name: "receive_where_handle",
			sql:  "RECEIVE * FROM ExpenseQueue WHERE conversation_handle = @handle",
		},
		{
			name: "receive_where_group",
			sql:  "RECEIVE * FROM ExpenseQueue WHERE conversation_group_id = @group_id",
		},
		{
			name: "receive_with_alias",
			sql:  "RECEIVE message_type_name AS MsgType, message_body AS Body FROM ExpenseQueue",
		},
		// CREATE CONTRACT - message type definitions
		{
			name: "contract_single_initiator",
			sql:  "CREATE CONTRACT [//MyApp/Contract] ([//MyApp/Request] SENT BY INITIATOR)",
		},
		{
			name: "contract_single_target",
			sql:  "CREATE CONTRACT [//MyApp/Contract] ([//MyApp/Reply] SENT BY TARGET)",
		},
		{
			name: "contract_single_any",
			sql:  "CREATE CONTRACT [//MyApp/Contract] ([//MyApp/Msg] SENT BY ANY)",
		},
		{
			name: "contract_multiple_types",
			sql:  "CREATE CONTRACT [//MyApp/Contract] ([//MyApp/Request] SENT BY INITIATOR, [//MyApp/Reply] SENT BY TARGET, [//MyApp/Status] SENT BY ANY)",
		},
		{
			name: "contract_with_default",
			sql:  "CREATE CONTRACT [//MyApp/Contract] ([DEFAULT] SENT BY ANY)",
		},
		{
			name: "contract_with_authorization",
			sql:  "CREATE CONTRACT [//MyApp/Contract] AUTHORIZATION dbo ([//MyApp/Request] SENT BY INITIATOR)",
		},
		// CREATE/ALTER QUEUE - ACTIVATION options
		{
			name: "queue_with_status",
			sql:  "CREATE QUEUE ExpenseQueue WITH STATUS = ON",
		},
		{
			name: "queue_with_retention",
			sql:  "CREATE QUEUE ExpenseQueue WITH RETENTION = ON",
		},
		{
			name: "queue_with_activation",
			sql:  "CREATE QUEUE ExpenseQueue WITH ACTIVATION (STATUS = ON, PROCEDURE_NAME = dbo.expense_proc, MAX_QUEUE_READERS = 5, EXECUTE AS SELF)",
		},
		{
			name: "queue_activation_execute_as_user",
			sql:  "CREATE QUEUE ExpenseQueue WITH ACTIVATION (PROCEDURE_NAME = expense_proc, MAX_QUEUE_READERS = 3, EXECUTE AS 'ExpenseUser')",
		},
		{
			name: "queue_activation_execute_as_owner",
			sql:  "CREATE QUEUE ExpenseQueue WITH ACTIVATION (PROCEDURE_NAME = expense_proc, MAX_QUEUE_READERS = 1, EXECUTE AS OWNER)",
		},
		{
			name: "queue_with_poison_message",
			sql:  "CREATE QUEUE ExpenseQueue WITH POISON_MESSAGE_HANDLING (STATUS = OFF)",
		},
		{
			name: "queue_full_options",
			sql:  "CREATE QUEUE ExpenseQueue WITH STATUS = ON, RETENTION = OFF, ACTIVATION (STATUS = ON, PROCEDURE_NAME = dbo.expense_proc, MAX_QUEUE_READERS = 10, EXECUTE AS SELF), POISON_MESSAGE_HANDLING (STATUS = ON)",
		},
		{
			name: "queue_on_filegroup",
			sql:  "CREATE QUEUE ExpenseQueue ON ExpenseFileGroup",
		},
		{
			name: "alter_queue_activation",
			sql:  "ALTER QUEUE ExpenseQueue WITH ACTIVATION (STATUS = ON, PROCEDURE_NAME = dbo.expense_proc, MAX_QUEUE_READERS = 5, EXECUTE AS SELF)",
		},
		{
			name: "alter_queue_poison",
			sql:  "ALTER QUEUE ExpenseQueue WITH POISON_MESSAGE_HANDLING (STATUS = OFF)",
		},
		// END CONVERSATION - WITH options
		{
			name: "end_conversation_simple",
			sql:  "END CONVERSATION @handle",
		},
		{
			name: "end_conversation_with_error",
			sql:  "END CONVERSATION @handle WITH ERROR = 100 DESCRIPTION = 'An error occurred'",
		},
		{
			name: "end_conversation_with_cleanup",
			sql:  "END CONVERSATION @handle WITH CLEANUP",
		},
		{
			name: "end_conversation_error_variable",
			sql:  "END CONVERSATION @handle WITH ERROR = @errorCode DESCRIPTION = @errorDesc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			stmt, ok := result.Items[0].(*ast.ServiceBrokerStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *ast.ServiceBrokerStmt, got %T", tt.sql, result.Items[0])
			}
			checkLocation(t, tt.sql, "ServiceBrokerStmt", stmt.Loc)

			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic:\n  s1=%s\n  s2=%s", tt.sql, s1, s2)
			}
			if s1 == "" {
				t.Errorf("Parse(%q): empty serialization", tt.sql)
			}
		})
	}
}

// TestParseWindowFrame tests ROWS/RANGE/GROUPS window frame specification in OVER clause (batch 73).
func TestParseWindowFrame(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// ROWS frame
		{"rows_unbounded_preceding", "SELECT SUM(x) OVER (ORDER BY id ROWS UNBOUNDED PRECEDING)"},
		{"rows_n_preceding", "SELECT SUM(x) OVER (ORDER BY id ROWS 3 PRECEDING)"},
		{"rows_current_row", "SELECT SUM(x) OVER (ORDER BY id ROWS CURRENT ROW)"},
		{"rows_between_unbounded_preceding_and_current_row", "SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)"},
		{"rows_between_n_preceding_and_n_following", "SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING)"},
		{"rows_between_current_row_and_unbounded_following", "SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)"},
		{"rows_between_unbounded_preceding_and_unbounded_following", "SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)"},
		{"rows_between_n_following_and_n_following", "SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN 2 FOLLOWING AND 10 FOLLOWING)"},
		{"rows_between_n_preceding_and_current_row", "SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN 5 PRECEDING AND CURRENT ROW)"},
		// RANGE frame
		{"range_unbounded_preceding", "SELECT SUM(x) OVER (ORDER BY id RANGE UNBOUNDED PRECEDING)"},
		{"range_current_row", "SELECT SUM(x) OVER (ORDER BY id RANGE CURRENT ROW)"},
		{"range_between_unbounded_preceding_and_current_row", "SELECT SUM(x) OVER (ORDER BY id RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)"},
		{"range_between_current_row_and_unbounded_following", "SELECT SUM(x) OVER (ORDER BY id RANGE BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)"},
		{"range_between_unbounded_preceding_and_unbounded_following", "SELECT SUM(x) OVER (ORDER BY id RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)"},
		// GROUPS frame (SQL Server 2022+)
		{"groups_unbounded_preceding", "SELECT SUM(x) OVER (ORDER BY id GROUPS UNBOUNDED PRECEDING)"},
		{"groups_between_n_preceding_and_n_following", "SELECT SUM(x) OVER (ORDER BY id GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING)"},
		{"groups_between_current_row_and_unbounded_following", "SELECT SUM(x) OVER (ORDER BY id GROUPS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)"},
		// Combined with PARTITION BY
		{"partition_and_rows_frame", "SELECT SUM(x) OVER (PARTITION BY dept ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING)"},
		{"partition_and_range_frame", "SELECT AVG(salary) OVER (PARTITION BY dept ORDER BY hire_date RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): expected 1 statement, got %d", tt.sql, result.Len())
			}

			s1 := ast.NodeToString(result.Items[0])
			s2 := ast.NodeToString(result.Items[0])
			if s1 != s2 {
				t.Errorf("Parse(%q): serialization not deterministic:\n  s1=%s\n  s2=%s", tt.sql, s1, s2)
			}
			if s1 == "" {
				t.Errorf("Parse(%q): empty serialization", tt.sql)
			}

			// Verify the serialization contains WINDOWFRAME
			if !strings.Contains(s1, "WINDOWFRAME") {
				t.Errorf("Parse(%q): expected WINDOWFRAME in serialization, got: %s", tt.sql, s1)
			}
		})
	}
}

func TestParseForXmlJsonOptions(t *testing.T) {
	t.Run("for_xml_type_root", func(t *testing.T) {
		tests := []struct {
			sql  string
			mode ast.ForMode
			sub  string
			chk  func(t *testing.T, fc *ast.ForClause)
		}{
			{
				sql:  "SELECT col1 FROM t FOR XML RAW, TYPE",
				mode: ast.ForXML, sub: "RAW",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Type {
						t.Error("expected Type=true")
					}
				},
			},
			{
				sql:  "SELECT col1 FROM t FOR XML RAW('Element'), TYPE, ROOT('MyRoot')",
				mode: ast.ForXML, sub: "RAW",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if fc.ElementName != "Element" {
						t.Errorf("expected ElementName='Element', got %q", fc.ElementName)
					}
					if !fc.Type {
						t.Error("expected Type=true")
					}
					if !fc.Root {
						t.Error("expected Root=true")
					}
					if fc.RootName != "MyRoot" {
						t.Errorf("expected RootName='MyRoot', got %q", fc.RootName)
					}
				},
			},
			{
				sql:  "SELECT col1 FROM t FOR XML AUTO, ROOT",
				mode: ast.ForXML, sub: "AUTO",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Root {
						t.Error("expected Root=true")
					}
					if fc.RootName != "" {
						t.Errorf("expected empty RootName, got %q", fc.RootName)
					}
				},
			},
			{
				sql:  "SELECT col1 FROM t FOR XML PATH('row'), ROOT('data'), TYPE",
				mode: ast.ForXML, sub: "PATH",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if fc.ElementName != "row" {
						t.Errorf("expected ElementName='row', got %q", fc.ElementName)
					}
					if !fc.Root {
						t.Error("expected Root=true")
					}
					if fc.RootName != "data" {
						t.Errorf("expected RootName='data', got %q", fc.RootName)
					}
					if !fc.Type {
						t.Error("expected Type=true")
					}
				},
			},
			{
				sql:  "SELECT col1 FROM t FOR XML AUTO, BINARY BASE64",
				mode: ast.ForXML, sub: "AUTO",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.BinaryBase64 {
						t.Error("expected BinaryBase64=true")
					}
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.ForClause == nil {
					t.Fatal("expected non-nil ForClause")
				}
				fc := stmt.ForClause
				if fc.Mode != tt.mode {
					t.Errorf("expected mode %d, got %d", tt.mode, fc.Mode)
				}
				if fc.SubMode != tt.sub {
					t.Errorf("expected subMode %q, got %q", tt.sub, fc.SubMode)
				}
				tt.chk(t, fc)
			})
		}
	})

	t.Run("for_xml_elements", func(t *testing.T) {
		tests := []struct {
			sql  string
			chk  func(t *testing.T, fc *ast.ForClause)
		}{
			{
				sql: "SELECT col1 FROM t FOR XML RAW, ELEMENTS",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Elements {
						t.Error("expected Elements=true")
					}
					if fc.ElementsMode != "" {
						t.Errorf("expected empty ElementsMode, got %q", fc.ElementsMode)
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML AUTO, ELEMENTS XSINIL",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Elements {
						t.Error("expected Elements=true")
					}
					if fc.ElementsMode != "XSINIL" {
						t.Errorf("expected ElementsMode='XSINIL', got %q", fc.ElementsMode)
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML PATH('row'), ELEMENTS ABSENT",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Elements {
						t.Error("expected Elements=true")
					}
					if fc.ElementsMode != "ABSENT" {
						t.Errorf("expected ElementsMode='ABSENT', got %q", fc.ElementsMode)
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML RAW, XMLDATA",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.XmlData {
						t.Error("expected XmlData=true")
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML AUTO, XMLSCHEMA",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.XmlSchema {
						t.Error("expected XmlSchema=true")
					}
					if fc.XmlSchemaURI != "" {
						t.Errorf("expected empty XmlSchemaURI, got %q", fc.XmlSchemaURI)
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML AUTO, XMLSCHEMA('http://example.com/ns')",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.XmlSchema {
						t.Error("expected XmlSchema=true")
					}
					if fc.XmlSchemaURI != "http://example.com/ns" {
						t.Errorf("expected XmlSchemaURI='http://example.com/ns', got %q", fc.XmlSchemaURI)
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML EXPLICIT, XMLDATA",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if fc.SubMode != "EXPLICIT" {
						t.Errorf("expected SubMode='EXPLICIT', got %q", fc.SubMode)
					}
					if !fc.XmlData {
						t.Error("expected XmlData=true")
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR XML AUTO, TYPE, ROOT('root'), ELEMENTS XSINIL",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Type {
						t.Error("expected Type=true")
					}
					if !fc.Root {
						t.Error("expected Root=true")
					}
					if fc.RootName != "root" {
						t.Errorf("expected RootName='root', got %q", fc.RootName)
					}
					if !fc.Elements {
						t.Error("expected Elements=true")
					}
					if fc.ElementsMode != "XSINIL" {
						t.Errorf("expected ElementsMode='XSINIL', got %q", fc.ElementsMode)
					}
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.ForClause == nil {
					t.Fatal("expected non-nil ForClause")
				}
				tt.chk(t, stmt.ForClause)
			})
		}
	})

	t.Run("for_json_root", func(t *testing.T) {
		tests := []struct {
			sql  string
			chk  func(t *testing.T, fc *ast.ForClause)
		}{
			{
				sql: "SELECT col1 FROM t FOR JSON PATH, ROOT('data')",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if fc.SubMode != "PATH" {
						t.Errorf("expected SubMode='PATH', got %q", fc.SubMode)
					}
					if !fc.Root {
						t.Error("expected Root=true")
					}
					if fc.RootName != "data" {
						t.Errorf("expected RootName='data', got %q", fc.RootName)
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR JSON AUTO, ROOT",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if fc.SubMode != "AUTO" {
						t.Errorf("expected SubMode='AUTO', got %q", fc.SubMode)
					}
					if !fc.Root {
						t.Error("expected Root=true")
					}
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.ForClause == nil {
					t.Fatal("expected non-nil ForClause")
				}
				if stmt.ForClause.Mode != ast.ForJSON {
					t.Error("expected ForJSON mode")
				}
				tt.chk(t, stmt.ForClause)
			})
		}
	})

	t.Run("for_json_options", func(t *testing.T) {
		tests := []struct {
			sql  string
			chk  func(t *testing.T, fc *ast.ForClause)
		}{
			{
				sql: "SELECT col1 FROM t FOR JSON PATH, INCLUDE_NULL_VALUES",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.IncludeNullValues {
						t.Error("expected IncludeNullValues=true")
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR JSON AUTO, WITHOUT_ARRAY_WRAPPER",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.WithoutArrayWrapper {
						t.Error("expected WithoutArrayWrapper=true")
					}
				},
			},
			{
				sql: "SELECT col1 FROM t FOR JSON PATH, ROOT('result'), INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER",
				chk: func(t *testing.T, fc *ast.ForClause) {
					if !fc.Root {
						t.Error("expected Root=true")
					}
					if fc.RootName != "result" {
						t.Errorf("expected RootName='result', got %q", fc.RootName)
					}
					if !fc.IncludeNullValues {
						t.Error("expected IncludeNullValues=true")
					}
					if !fc.WithoutArrayWrapper {
						t.Error("expected WithoutArrayWrapper=true")
					}
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SelectStmt)
				if stmt.ForClause == nil {
					t.Fatal("expected non-nil ForClause")
				}
				if stmt.ForClause.Mode != ast.ForJSON {
					t.Error("expected ForJSON mode")
				}
				tt.chk(t, stmt.ForClause)
			})
		}
	})
}

// TestParseCompoundAssignment tests compound assignment operators (batch 75).
func TestParseCompoundAssignment(t *testing.T) {
	t.Run("update_compound_assignment", func(t *testing.T) {
		tests := []struct {
			sql string
			op  string
		}{
			{"UPDATE t SET col += 1", "+="},
			{"UPDATE t SET col -= 1", "-="},
			{"UPDATE t SET col *= 2", "*="},
			{"UPDATE t SET col /= 2", "/="},
			{"UPDATE t SET col %= 3", "%="},
			{"UPDATE t SET col &= 0xFF", "&="},
			{"UPDATE t SET col |= 0x01", "|="},
			{"UPDATE t SET col ^= 0x0F", "^="},
			// Simple = should have empty operator (default)
			{"UPDATE t SET col = 1", ""},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.UpdateStmt)
				if stmt.SetClause == nil || stmt.SetClause.Len() == 0 {
					t.Fatal("expected non-empty SetClause")
				}
				se := stmt.SetClause.Items[0].(*ast.SetExpr)
				if se.Operator != tt.op {
					t.Errorf("expected operator %q, got %q", tt.op, se.Operator)
				}
			})
		}
	})

	t.Run("update_compound_with_variable", func(t *testing.T) {
		tests := []struct {
			sql string
			op  string
		}{
			{"UPDATE t SET @v += col", "+="},
			{"UPDATE t SET @v -= col", "-="},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.UpdateStmt)
				se := stmt.SetClause.Items[0].(*ast.SetExpr)
				if se.Variable == "" {
					t.Error("expected variable in SetExpr")
				}
				if se.Operator != tt.op {
					t.Errorf("expected operator %q, got %q", tt.op, se.Operator)
				}
			})
		}
	})

	t.Run("update_compound_mixed", func(t *testing.T) {
		sql := "UPDATE t SET col1 += 1, col2 = 2, col3 *= 3"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.UpdateStmt)
		if stmt.SetClause.Len() != 3 {
			t.Fatalf("expected 3 set exprs, got %d", stmt.SetClause.Len())
		}
		ops := []string{"+=", "", "*="}
		for i, expected := range ops {
			se := stmt.SetClause.Items[i].(*ast.SetExpr)
			if se.Operator != expected {
				t.Errorf("item[%d]: expected operator %q, got %q", i, expected, se.Operator)
			}
		}
	})

	t.Run("set_var_compound_assignment", func(t *testing.T) {
		tests := []struct {
			sql string
			op  string
		}{
			{"SET @x += 1", "+="},
			{"SET @x -= 1", "-="},
			{"SET @x *= 2", "*="},
			{"SET @x /= 2", "/="},
			{"SET @x %= 3", "%="},
			{"SET @x &= 0xFF", "&="},
			{"SET @x |= 0x01", "|="},
			{"SET @x ^= 0x0F", "^="},
			// Simple = should have empty operator
			{"SET @x = 1", ""},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SetStmt)
				if stmt.Operator != tt.op {
					t.Errorf("expected operator %q, got %q", tt.op, stmt.Operator)
				}
			})
		}
	})

	t.Run("set_var_compound_with_expr", func(t *testing.T) {
		// Compound assignment with a complex expression
		sql := "SET @total += @price * @quantity"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SetStmt)
		if stmt.Operator != "+=" {
			t.Errorf("expected operator '+=', got %q", stmt.Operator)
		}
		if stmt.Value == nil {
			t.Error("expected non-nil Value")
		}
	})

	t.Run("serialization_roundtrip", func(t *testing.T) {
		// Ensure compound assignment operators appear in serialized output
		sqls := []string{
			"UPDATE t SET col += 1",
			"SET @x -= 5",
		}
		for _, sql := range sqls {
			result := ParseAndCheck(t, sql)
			s := ast.NodeToString(result.Items[0])
			if s == "" {
				t.Errorf("NodeToString returned empty for %q", sql)
			}
		}
	})
}

// TestParseCreateDatabaseDepth tests CREATE DATABASE full options (batch 76).
func TestParseCreateDatabaseDepth(t *testing.T) {
	t.Run("create_database_filespec", func(t *testing.T) {
		// ON PRIMARY with file specs
		tests := []struct {
			name string
			sql  string
		}{
			{
				"single filespec",
				`CREATE DATABASE mydb
				ON PRIMARY
				( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf', SIZE = 10MB, MAXSIZE = 100MB, FILEGROWTH = 5MB )`,
			},
			{
				"multiple filespecs",
				`CREATE DATABASE mydb
				ON PRIMARY
				( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf', SIZE = 10MB ),
				( NAME = mydb_data2, FILENAME = 'C:\data\mydb2.ndf', SIZE = 5MB, MAXSIZE = UNLIMITED, FILEGROWTH = 10% )`,
			},
			{
				"filespec with GB size",
				`CREATE DATABASE mydb
				ON PRIMARY
				( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf', SIZE = 1GB, MAXSIZE = 50GB, FILEGROWTH = 100MB )`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
				if !ok {
					t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Name != "mydb" {
					t.Errorf("expected name mydb, got %s", stmt.Name)
				}
				if stmt.OnPrimary == nil {
					t.Fatalf("expected OnPrimary, got nil")
				}
			})
		}
	})

	t.Run("create_database_log_on", func(t *testing.T) {
		sql := `CREATE DATABASE mydb
			ON PRIMARY
			( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf', SIZE = 10MB )
			LOG ON
			( NAME = mydb_log, FILENAME = 'C:\data\mydb.ldf', SIZE = 5MB, MAXSIZE = 25MB, FILEGROWTH = 5MB )`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.LogOn == nil {
			t.Fatalf("expected LogOn, got nil")
		}
		if stmt.LogOn.Len() != 1 {
			t.Errorf("expected 1 log file, got %d", stmt.LogOn.Len())
		}
	})

	t.Run("create_database_collate", func(t *testing.T) {
		sql := `CREATE DATABASE mydb COLLATE Latin1_General_CI_AS`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Collation != "Latin1_General_CI_AS" {
			t.Errorf("expected collation Latin1_General_CI_AS, got %s", stmt.Collation)
		}
	})

	t.Run("create_database_for_attach", func(t *testing.T) {
		tests := []struct {
			name              string
			sql               string
			forAttach         bool
			forAttachRebuild  bool
			hasAttachOptions  bool
		}{
			{
				"basic attach",
				`CREATE DATABASE mydb
				ON ( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf' )
				FOR ATTACH`,
				true, false, false,
			},
			{
				"attach with broker options",
				`CREATE DATABASE mydb
				ON ( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf' )
				FOR ATTACH WITH ENABLE_BROKER`,
				true, false, true,
			},
			{
				"attach rebuild log",
				`CREATE DATABASE mydb
				ON ( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf' )
				FOR ATTACH_REBUILD_LOG`,
				false, true, false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
				if !ok {
					t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.ForAttach != tt.forAttach {
					t.Errorf("ForAttach: expected %v, got %v", tt.forAttach, stmt.ForAttach)
				}
				if stmt.ForAttachRebuildLog != tt.forAttachRebuild {
					t.Errorf("ForAttachRebuildLog: expected %v, got %v", tt.forAttachRebuild, stmt.ForAttachRebuildLog)
				}
				if tt.hasAttachOptions && stmt.AttachOptions == nil {
					t.Error("expected AttachOptions, got nil")
				}
			})
		}
	})

	t.Run("create_database_snapshot", func(t *testing.T) {
		sql := `CREATE DATABASE mydb_snapshot
			ON ( NAME = mydb_data, FILENAME = 'C:\snapshots\mydb.ss' )
			AS SNAPSHOT OF mydb`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.SnapshotOf != "mydb" {
			t.Errorf("expected SnapshotOf mydb, got %s", stmt.SnapshotOf)
		}
	})

	t.Run("create_database_containment", func(t *testing.T) {
		sql := `CREATE DATABASE mydb CONTAINMENT = PARTIAL`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Containment != "PARTIAL" {
			t.Errorf("expected containment PARTIAL, got %s", stmt.Containment)
		}
	})

	t.Run("create_database_with_options", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{
				"with ledger",
				`CREATE DATABASE mydb WITH LEDGER = ON`,
			},
			{
				"with nested triggers and db chaining",
				`CREATE DATABASE mydb WITH NESTED_TRIGGERS = ON, DB_CHAINING OFF`,
			},
			{
				"with default language",
				`CREATE DATABASE mydb WITH DEFAULT_LANGUAGE = us_english`,
			},
			{
				"with filestream option",
				`CREATE DATABASE mydb WITH FILESTREAM ( NON_TRANSACTED_ACCESS = FULL, DIRECTORY_NAME = 'mydir' )`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
				if !ok {
					t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.WithOptions == nil {
					t.Fatal("expected WithOptions, got nil")
				}
			})
		}
	})

	t.Run("create_database_filegroup", func(t *testing.T) {
		sql := `CREATE DATABASE mydb
			ON PRIMARY
			( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf', SIZE = 10MB )
			FILEGROUP fg1
			( NAME = mydb_fg1, FILENAME = 'C:\data\mydb_fg1.ndf', SIZE = 5MB )
			LOG ON
			( NAME = mydb_log, FILENAME = 'C:\data\mydb.ldf', SIZE = 5MB )`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.OnPrimary == nil {
			t.Fatal("expected OnPrimary, got nil")
		}
		if stmt.Filegroups == nil {
			t.Fatal("expected Filegroups, got nil")
		}
		if stmt.Filegroups.Len() != 1 {
			t.Errorf("expected 1 filegroup, got %d", stmt.Filegroups.Len())
		}
		if stmt.LogOn == nil {
			t.Fatal("expected LogOn, got nil")
		}
	})

	t.Run("create_database_filegroup_contains_filestream", func(t *testing.T) {
		sql := `CREATE DATABASE mydb
			ON PRIMARY
			( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf' )
			FILEGROUP fsg CONTAINS FILESTREAM DEFAULT
			( NAME = mydb_fs, FILENAME = 'C:\data\mydb_fs' )`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Filegroups == nil {
			t.Fatal("expected Filegroups, got nil")
		}
		fg := stmt.Filegroups.Items[0].(*ast.DatabaseFilegroup)
		if !fg.ContainsFilestream {
			t.Error("expected ContainsFilestream true")
		}
		if !fg.IsDefault {
			t.Error("expected IsDefault true")
		}
	})

	t.Run("create_database_full_example", func(t *testing.T) {
		sql := `CREATE DATABASE mydb
			CONTAINMENT = NONE
			ON PRIMARY
			( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf', SIZE = 100MB, MAXSIZE = 1GB, FILEGROWTH = 10MB )
			LOG ON
			( NAME = mydb_log, FILENAME = 'C:\data\mydb.ldf', SIZE = 50MB, MAXSIZE = 500MB, FILEGROWTH = 10% )
			COLLATE SQL_Latin1_General_CP1_CI_AS
			WITH LEDGER = OFF`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Containment != "NONE" {
			t.Errorf("expected NONE, got %s", stmt.Containment)
		}
		if stmt.OnPrimary == nil {
			t.Fatal("expected OnPrimary")
		}
		if stmt.LogOn == nil {
			t.Fatal("expected LogOn")
		}
		if stmt.Collation != "SQL_Latin1_General_CP1_CI_AS" {
			t.Errorf("expected SQL_Latin1_General_CP1_CI_AS, got %s", stmt.Collation)
		}
		if stmt.WithOptions == nil {
			t.Fatal("expected WithOptions")
		}
	})

	t.Run("create_database_attach_with_multiple_options", func(t *testing.T) {
		sql := `CREATE DATABASE mydb
			ON ( NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf' )
			FOR ATTACH WITH ENABLE_BROKER, RESTRICTED_USER`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if !stmt.ForAttach {
			t.Error("expected ForAttach true")
		}
		if stmt.AttachOptions == nil {
			t.Fatal("expected AttachOptions")
		}
		if stmt.AttachOptions.Len() != 2 {
			t.Errorf("expected 2 attach options, got %d", stmt.AttachOptions.Len())
		}
	})
}
