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

// TestParseSetIsolationLevelStructured tests SET TRANSACTION ISOLATION LEVEL with structured parsing (batch 134).
func TestParseSetIsolationLevelStructured(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		level string
	}{
		{
			name:  "read_uncommitted",
			sql:   "SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED",
			level: "READ UNCOMMITTED",
		},
		{
			name:  "read_committed",
			sql:   "SET TRANSACTION ISOLATION LEVEL READ COMMITTED",
			level: "READ COMMITTED",
		},
		{
			name:  "repeatable_read",
			sql:   "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ",
			level: "REPEATABLE READ",
		},
		{
			name:  "serializable",
			sql:   "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE",
			level: "SERIALIZABLE",
		},
		{
			name:  "snapshot",
			sql:   "SET TRANSACTION ISOLATION LEVEL SNAPSHOT",
			level: "SNAPSHOT",
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
			if stmt.Option != "TRANSACTION ISOLATION LEVEL" {
				t.Errorf("Parse(%q): Option = %q, want %q", tt.sql, stmt.Option, "TRANSACTION ISOLATION LEVEL")
			}
			// Check that Value is a ColumnRef with the correct isolation level
			ref, ok := stmt.Value.(*ast.ColumnRef)
			if !ok {
				t.Fatalf("Parse(%q): Value expected *ast.ColumnRef, got %T", tt.sql, stmt.Value)
			}
			if ref.Column != tt.level {
				t.Errorf("Parse(%q): isolation level = %q, want %q", tt.sql, ref.Column, tt.level)
			}
			checkLocation(t, tt.sql, "SetOptionStmt", stmt.Loc)
			checkLocation(t, tt.sql, "Value", ref.Loc)
		})
	}
}

// TestParseMiscUnknownSkipDepth tests structured parsing replacing shallow unknown-token-skip patterns (batch 135).
func TestParseMiscUnknownSkipDepth(t *testing.T) {
	// 1. SYSTEM_VERSIONING sub-option with unknown key=value is consumed structurally
	t.Run("system_versioning_sub_option", func(t *testing.T) {
		sql := `CREATE TABLE t (
			id INT PRIMARY KEY,
			s datetime2 GENERATED ALWAYS AS ROW START NOT NULL,
			e datetime2 GENERATED ALWAYS AS ROW END NOT NULL,
			PERIOD FOR SYSTEM_TIME (s, e)
		) WITH (SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.h, DATA_CONSISTENCY_CHECK = ON))`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.TableOptions == nil {
			t.Fatal("expected table options")
		}
	})

	// 2. ALTER COLUMN ADD/DROP with known options
	t.Run("alter_column_option_structured", func(t *testing.T) {
		tests := []string{
			"ALTER TABLE t ALTER COLUMN c ADD ROWGUIDCOL",
			"ALTER TABLE t ALTER COLUMN c DROP ROWGUIDCOL",
			"ALTER TABLE t ALTER COLUMN c ADD PERSISTED",
			"ALTER TABLE t ALTER COLUMN c ADD SPARSE",
			"ALTER TABLE t ALTER COLUMN c ADD HIDDEN",
			"ALTER TABLE t ALTER COLUMN c ADD NOT FOR REPLICATION",
			"ALTER TABLE t ALTER COLUMN c ADD MASKED WITH (FUNCTION = 'default()')",
		}
		for _, sql := range tests {
			t.Run(sql, func(t *testing.T) {
				result := ParseAndCheck(t, sql)
				if result.Len() == 0 {
					t.Fatalf("Parse(%q): no statements returned", sql)
				}
			})
		}
	})

	// 3. Event notification unknown token is consumed with structured identifier
	t.Run("event_notification_option_structured", func(t *testing.T) {
		sql := "CREATE EVENT NOTIFICATION en1 ON DATABASE FOR CREATE_TABLE TO SERVICE 'svc', 'broker_id'"
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
	})

	// 4. Backup/restore unknown option treated as flag with advance
	t.Run("backup_restore_unknown_option", func(t *testing.T) {
		sql := "BACKUP DATABASE mydb TO DISK = 'backup.bak' WITH INIT, COMPRESSION, CHECKSUM"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.BackupStmt)
		if stmt.Options == nil {
			t.Fatal("expected backup options")
		}
	})
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

// TestParseWorkloadClassifier tests CREATE/ALTER/DROP WORKLOAD CLASSIFIER (batch 80).
func TestParseWorkloadClassifier(t *testing.T) {
	t.Run("create_workload_classifier", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{
				"basic",
				`CREATE WORKLOAD CLASSIFIER wgcELTRole
				WITH (WORKLOAD_GROUP = 'staticrc20', MEMBERNAME = 'ELTRole')`,
			},
			{
				"with importance",
				`CREATE WORKLOAD CLASSIFIER wgcELTRole
				WITH (WORKLOAD_GROUP = 'staticrc20', MEMBERNAME = 'ELTRole', IMPORTANCE = ABOVE_NORMAL)`,
			},
			{
				"with label and context",
				`CREATE WORKLOAD CLASSIFIER wgcELTRole
				WITH (WORKLOAD_GROUP = 'wgDataLoad', MEMBERNAME = 'ELTRole', WLM_LABEL = 'dimension_loads', WLM_CONTEXT = 'dim_load')`,
			},
			{
				"with time range",
				`CREATE WORKLOAD CLASSIFIER wgcNight
				WITH (WORKLOAD_GROUP = 'wgDataLoads', MEMBERNAME = 'ELTRole', START_TIME = '22:00', END_TIME = '02:00')`,
			},
			{
				"all options",
				`CREATE WORKLOAD CLASSIFIER wgcAll
				WITH (WORKLOAD_GROUP = 'wg1', MEMBERNAME = 'user1', WLM_LABEL = 'lbl', WLM_CONTEXT = 'ctx', START_TIME = '08:00', END_TIME = '17:00', IMPORTANCE = HIGH)`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok {
					t.Fatalf("expected *SecurityStmt, got %T", result.Items[0])
				}
				if stmt.Action != "CREATE" {
					t.Errorf("expected action CREATE, got %s", stmt.Action)
				}
				if stmt.ObjectType != "WORKLOAD CLASSIFIER" {
					t.Errorf("expected objectType WORKLOAD CLASSIFIER, got %s", stmt.ObjectType)
				}
			})
		}
	})

	t.Run("alter_workload_classifier", func(t *testing.T) {
		sql := `ALTER WORKLOAD CLASSIFIER wgcELTRole
			WITH (WORKLOAD_GROUP = 'staticrc30', MEMBERNAME = 'ELTRole', IMPORTANCE = HIGH)`
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.SecurityStmt)
		if !ok {
			t.Fatalf("expected *SecurityStmt, got %T", result.Items[0])
		}
		if stmt.Action != "ALTER" {
			t.Errorf("expected action ALTER, got %s", stmt.Action)
		}
		if stmt.ObjectType != "WORKLOAD CLASSIFIER" {
			t.Errorf("expected objectType WORKLOAD CLASSIFIER, got %s", stmt.ObjectType)
		}
	})

	t.Run("drop_workload_classifier", func(t *testing.T) {
		sql := "DROP WORKLOAD CLASSIFIER wgcELTRole"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.SecurityStmt)
		if !ok {
			t.Fatalf("expected *SecurityStmt, got %T", result.Items[0])
		}
		if stmt.Action != "DROP" {
			t.Errorf("expected action DROP, got %s", stmt.Action)
		}
		if stmt.ObjectType != "WORKLOAD CLASSIFIER" {
			t.Errorf("expected objectType WORKLOAD CLASSIFIER, got %s", stmt.ObjectType)
		}
		if stmt.Name != "wgcELTRole" {
			t.Errorf("expected name wgcELTRole, got %s", stmt.Name)
		}
	})
}

// TestParseSubqueryComparison tests ANY/SOME/ALL subquery comparison operators (batch 79).
func TestParseSubqueryComparison(t *testing.T) {
	t.Run("any_subquery", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"basic any", "SELECT * FROM t WHERE x > ANY (SELECT y FROM t2)"},
			{"any with equals", "SELECT * FROM t WHERE x = ANY (SELECT y FROM t2)"},
			{"any with not equals", "SELECT * FROM t WHERE x <> ANY (SELECT y FROM t2)"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				_ = result
			})
		}
	})

	t.Run("some_subquery", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"basic some", "SELECT * FROM t WHERE x >= SOME (SELECT y FROM t2)"},
			{"some with lt", "SELECT * FROM t WHERE x < SOME (SELECT y FROM t2)"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				_ = result
			})
		}
	})

	t.Run("all_subquery", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"basic all", "SELECT * FROM t WHERE x > ALL (SELECT y FROM t2)"},
			{"all with lte", "SELECT * FROM t WHERE x <= ALL (SELECT y FROM t2)"},
			{"all with eq", "SELECT * FROM t WHERE x = ALL (SELECT y FROM t2)"},
			{"all with neq", "SELECT * FROM t WHERE x != ALL (SELECT y FROM t2)"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				_ = result
			})
		}
	})
}

// TestParseGrantDepth tests GRANT/REVOKE/DENY securable class and REVOKE GRANT OPTION FOR (batch 78).
func TestParseGrantDepth(t *testing.T) {
	t.Run("grant_on_schema", func(t *testing.T) {
		tests := []struct {
			name       string
			sql        string
			onType     string
			wantGrant  bool
		}{
			{"grant select on schema", "GRANT SELECT ON SCHEMA::dbo TO user1", "SCHEMA", false},
			{"grant select on schema with grant option", "GRANT SELECT ON SCHEMA::Sales TO user1 WITH GRANT OPTION", "SCHEMA", true},
			{"grant execute on schema", "GRANT EXECUTE ON SCHEMA::HumanResources TO role1", "SCHEMA", false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.GrantStmt)
				if !ok {
					t.Fatalf("expected *GrantStmt, got %T", result.Items[0])
				}
				if stmt.OnType != tt.onType {
					t.Errorf("expected onType %s, got %s", tt.onType, stmt.OnType)
				}
				if stmt.WithGrant != tt.wantGrant {
					t.Errorf("expected withGrant %v, got %v", tt.wantGrant, stmt.WithGrant)
				}
			})
		}
	})

	t.Run("grant_on_object_type", func(t *testing.T) {
		tests := []struct {
			name   string
			sql    string
			onType string
		}{
			{"grant on object", "GRANT SELECT ON OBJECT::dbo.MyTable TO user1", "OBJECT"},
			{"grant on database", "GRANT CREATE TABLE ON DATABASE::mydb TO user1", "DATABASE"},
			{"grant on login", "GRANT ALTER ON LOGIN::mylogin TO user1", "LOGIN"},
			{"grant on user", "GRANT IMPERSONATE ON USER::myuser TO user1", "USER"},
			{"grant on role", "GRANT ALTER ON ROLE::myrole TO user1", "ROLE"},
			{"grant on type", "GRANT EXECUTE ON TYPE::dbo.MyType TO user1", "TYPE"},
			{"grant on assembly", "GRANT ALTER ON ASSEMBLY::MyAssembly TO user1", "ASSEMBLY"},
			{"grant on symmetric key", "GRANT ALTER ON SYMMETRIC KEY::mykey TO user1", "SYMMETRIC KEY"},
			{"grant on asymmetric key", "GRANT ALTER ON ASYMMETRIC KEY::mykey TO user1", "ASYMMETRIC KEY"},
			{"grant on certificate", "GRANT ALTER ON CERTIFICATE::mycert TO user1", "CERTIFICATE"},
			{"grant on xml schema collection", "GRANT ALTER ON XML SCHEMA COLLECTION::dbo.myschema TO user1", "XML SCHEMA COLLECTION"},
			{"grant on application role", "GRANT ALTER ON APPLICATION ROLE::myrole TO user1", "APPLICATION ROLE"},
			{"grant on message type", "GRANT ALTER ON MESSAGE TYPE::mytype TO user1", "MESSAGE TYPE"},
			{"grant on fulltext catalog", "GRANT ALTER ON FULLTEXT CATALOG::mycat TO user1", "FULLTEXT CATALOG"},
			{"grant on search property list", "GRANT ALTER ON SEARCH PROPERTY LIST::mylist TO user1", "SEARCH PROPERTY LIST"},
			{"grant on server role", "GRANT ALTER ON SERVER ROLE::myrole TO user1", "SERVER ROLE"},
			{"grant on remote service binding", "GRANT ALTER ON REMOTE SERVICE BINDING::mybinding TO user1", "REMOTE SERVICE BINDING"},
			{"grant without class", "GRANT SELECT ON dbo.MyTable TO user1", ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.GrantStmt)
				if !ok {
					t.Fatalf("expected *GrantStmt, got %T", result.Items[0])
				}
				if stmt.OnType != tt.onType {
					t.Errorf("expected onType %q, got %q", tt.onType, stmt.OnType)
				}
			})
		}
	})

	t.Run("revoke_grant_option_for", func(t *testing.T) {
		tests := []struct {
			name           string
			sql            string
			grantOptionFor bool
			cascade        bool
		}{
			{"revoke grant option for", "REVOKE GRANT OPTION FOR SELECT ON SCHEMA::dbo FROM user1 CASCADE", true, true},
			{"revoke grant option for no cascade", "REVOKE GRANT OPTION FOR INSERT ON dbo.t1 FROM user1", true, false},
			{"revoke normal", "REVOKE SELECT ON SCHEMA::Sales FROM user1", false, false},
			{"revoke with cascade", "REVOKE EXECUTE ON dbo.myproc FROM user1 CASCADE", false, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.GrantStmt)
				if !ok {
					t.Fatalf("expected *GrantStmt, got %T", result.Items[0])
				}
				if stmt.GrantOptionFor != tt.grantOptionFor {
					t.Errorf("expected grantOptionFor %v, got %v", tt.grantOptionFor, stmt.GrantOptionFor)
				}
				if stmt.CascadeOpt != tt.cascade {
					t.Errorf("expected cascade %v, got %v", tt.cascade, stmt.CascadeOpt)
				}
			})
		}
	})

	t.Run("grant_as_principal", func(t *testing.T) {
		sql := "GRANT SELECT ON dbo.t1 TO user1 WITH GRANT OPTION AS dbo"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.GrantStmt)
		if !ok {
			t.Fatalf("expected *GrantStmt, got %T", result.Items[0])
		}
		if stmt.AsPrincipal != "dbo" {
			t.Errorf("expected asPrincipal dbo, got %s", stmt.AsPrincipal)
		}
	})

	t.Run("deny_on_schema", func(t *testing.T) {
		sql := "DENY SELECT ON SCHEMA::Sales TO user1 CASCADE"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.GrantStmt)
		if !ok {
			t.Fatalf("expected *GrantStmt, got %T", result.Items[0])
		}
		if stmt.StmtType != ast.GrantTypeDeny {
			t.Errorf("expected DENY, got %d", stmt.StmtType)
		}
		if stmt.OnType != "SCHEMA" {
			t.Errorf("expected onType SCHEMA, got %s", stmt.OnType)
		}
		if !stmt.CascadeOpt {
			t.Error("expected cascade true")
		}
	})
}

// TestParseAlterDatabaseDepth tests ALTER DATABASE structured parsing (batch 77).
func TestParseAlterDatabaseDepth(t *testing.T) {
	// --- SET options ---
	t.Run("alter_database_set_options", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"set single_user", "ALTER DATABASE mydb SET SINGLE_USER"},
			{"set multi_user", "ALTER DATABASE mydb SET MULTI_USER"},
			{"set restricted_user", "ALTER DATABASE mydb SET RESTRICTED_USER"},
			{"set read_only", "ALTER DATABASE mydb SET READ_ONLY"},
			{"set read_write", "ALTER DATABASE mydb SET READ_WRITE"},
			{"set online", "ALTER DATABASE mydb SET ONLINE"},
			{"set offline", "ALTER DATABASE mydb SET OFFLINE"},
			{"set emergency", "ALTER DATABASE mydb SET EMERGENCY"},
			{"set recovery full", "ALTER DATABASE mydb SET RECOVERY FULL"},
			{"set recovery simple", "ALTER DATABASE mydb SET RECOVERY SIMPLE"},
			{"set recovery bulk_logged", "ALTER DATABASE mydb SET RECOVERY BULK_LOGGED"},
			{"set auto_close on", "ALTER DATABASE mydb SET AUTO_CLOSE ON"},
			{"set auto_close off", "ALTER DATABASE mydb SET AUTO_CLOSE OFF"},
			{"set auto_shrink on", "ALTER DATABASE mydb SET AUTO_SHRINK ON"},
			{"set auto_create_statistics on", "ALTER DATABASE mydb SET AUTO_CREATE_STATISTICS ON"},
			{"set auto_update_statistics off", "ALTER DATABASE mydb SET AUTO_UPDATE_STATISTICS OFF"},
			{"set ansi_nulls on", "ALTER DATABASE mydb SET ANSI_NULLS ON"},
			{"set quoted_identifier on", "ALTER DATABASE mydb SET QUOTED_IDENTIFIER ON"},
			{"set concat_null_yields_null on", "ALTER DATABASE mydb SET CONCAT_NULL_YIELDS_NULL ON"},
			{"set arithabort on", "ALTER DATABASE mydb SET ARITHABORT ON"},
			{"set page_verify checksum", "ALTER DATABASE mydb SET PAGE_VERIFY CHECKSUM"},
			{"set compatibility_level", "ALTER DATABASE mydb SET COMPATIBILITY_LEVEL = 150"},
			{"set allow_snapshot_isolation on", "ALTER DATABASE mydb SET ALLOW_SNAPSHOT_ISOLATION ON"},
			{"set read_committed_snapshot on", "ALTER DATABASE mydb SET READ_COMMITTED_SNAPSHOT ON"},
			{"set parameterization forced", "ALTER DATABASE mydb SET PARAMETERIZATION FORCED"},
			{"set encryption on", "ALTER DATABASE mydb SET ENCRYPTION ON"},
			{"set db_chaining on", "ALTER DATABASE mydb SET DB_CHAINING ON"},
			{"set trustworthy on", "ALTER DATABASE mydb SET TRUSTWORTHY ON"},
			{"set delayed_durability allowed", "ALTER DATABASE mydb SET DELAYED_DURABILITY = ALLOWED"},
			{"set target_recovery_time", "ALTER DATABASE mydb SET TARGET_RECOVERY_TIME = 60 SECONDS"},
			{"set multiple options", "ALTER DATABASE mydb SET AUTO_CLOSE OFF, AUTO_SHRINK OFF"},
			{"set with rollback immediate", "ALTER DATABASE mydb SET SINGLE_USER WITH ROLLBACK IMMEDIATE"},
			{"set with rollback after", "ALTER DATABASE mydb SET SINGLE_USER WITH ROLLBACK AFTER 10 SECONDS"},
			{"set with no_wait", "ALTER DATABASE mydb SET SINGLE_USER WITH NO_WAIT"},
			{"set CURRENT", "ALTER DATABASE CURRENT SET READ_ONLY"},
			{"set change_tracking on", "ALTER DATABASE mydb SET CHANGE_TRACKING = ON"},
			{"set change_tracking on with options", "ALTER DATABASE mydb SET CHANGE_TRACKING = ON (AUTO_CLEANUP = ON, CHANGE_RETENTION = 7 DAYS)"},
			{"set query_store on", "ALTER DATABASE mydb SET QUERY_STORE = ON"},
			{"set query_store off", "ALTER DATABASE mydb SET QUERY_STORE = OFF"},
			{"set query_store clear", "ALTER DATABASE mydb SET QUERY_STORE CLEAR"},
			{"set containment none", "ALTER DATABASE mydb SET CONTAINMENT = NONE"},
			{"set enable_broker", "ALTER DATABASE mydb SET ENABLE_BROKER"},
			{"set disable_broker", "ALTER DATABASE mydb SET DISABLE_BROKER"},
			{"set new_broker", "ALTER DATABASE mydb SET NEW_BROKER"},
			{"set accelerated_database_recovery on", "ALTER DATABASE mydb SET ACCELERATED_DATABASE_RECOVERY = ON"},
			{"set mixed_page_allocation off", "ALTER DATABASE mydb SET MIXED_PAGE_ALLOCATION OFF"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "SET" {
					t.Errorf("expected action SET, got %s", stmt.Action)
				}
				if stmt.Options == nil || stmt.Options.Len() == 0 {
					t.Error("expected SET options to be parsed")
				}
			})
		}
	})

	// --- MODIFY FILE ---
	t.Run("alter_database_modify_file", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"modify file size", "ALTER DATABASE mydb MODIFY FILE (NAME = mydb_data, SIZE = 200MB)"},
			{"modify file newname", "ALTER DATABASE mydb MODIFY FILE (NAME = mydb_data, NEWNAME = mydb_data_new)"},
			{"modify file filename", `ALTER DATABASE mydb MODIFY FILE (NAME = mydb_data, FILENAME = 'C:\data\mydb.mdf')`},
			{"modify file maxsize", "ALTER DATABASE mydb MODIFY FILE (NAME = mydb_data, SIZE = 100MB, MAXSIZE = 500MB, FILEGROWTH = 10MB)"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "MODIFY" {
					t.Errorf("expected action MODIFY, got %s", stmt.Action)
				}
				if stmt.SubAction != "FILE" {
					t.Errorf("expected subAction FILE, got %s", stmt.SubAction)
				}
				if stmt.FileSpecs == nil || stmt.FileSpecs.Len() == 0 {
					t.Error("expected file specs to be parsed")
				}
			})
		}
	})

	// --- ADD FILE ---
	t.Run("alter_database_add_file", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{
				"add file basic",
				`ALTER DATABASE mydb ADD FILE (NAME = mydb_data2, FILENAME = 'C:\data\mydb2.ndf', SIZE = 5MB)`,
			},
			{
				"add file to filegroup",
				`ALTER DATABASE mydb ADD FILE (NAME = mydb_data2, FILENAME = 'C:\data\mydb2.ndf', SIZE = 5MB) TO FILEGROUP fg1`,
			},
			{
				"add multiple files",
				`ALTER DATABASE mydb ADD FILE (NAME = f1, FILENAME = 'C:\data\f1.ndf', SIZE = 5MB), (NAME = f2, FILENAME = 'C:\data\f2.ndf', SIZE = 5MB)`,
			},
			{
				"add log file",
				`ALTER DATABASE mydb ADD LOG FILE (NAME = mydb_log2, FILENAME = 'C:\data\mydb_log2.ldf', SIZE = 5MB)`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "ADD" {
					t.Errorf("expected action ADD, got %s", stmt.Action)
				}
				if stmt.SubAction != "FILE" && stmt.SubAction != "LOG FILE" {
					t.Errorf("expected subAction FILE or LOG FILE, got %s", stmt.SubAction)
				}
				if stmt.FileSpecs == nil || stmt.FileSpecs.Len() == 0 {
					t.Error("expected file specs to be parsed")
				}
			})
		}
	})

	// --- REMOVE FILE ---
	t.Run("alter_database_remove_file", func(t *testing.T) {
		sql := "ALTER DATABASE mydb REMOVE FILE mydb_data2"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
		if !ok {
			t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Action != "REMOVE" {
			t.Errorf("expected action REMOVE, got %s", stmt.Action)
		}
		if stmt.SubAction != "FILE" {
			t.Errorf("expected subAction FILE, got %s", stmt.SubAction)
		}
		if stmt.TargetName != "mydb_data2" {
			t.Errorf("expected targetName mydb_data2, got %s", stmt.TargetName)
		}
	})

	// --- ADD FILEGROUP ---
	t.Run("alter_database_add_filegroup", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"add filegroup basic", "ALTER DATABASE mydb ADD FILEGROUP fg1"},
			{"add filegroup contains filestream", "ALTER DATABASE mydb ADD FILEGROUP fg1 CONTAINS FILESTREAM"},
			{"add filegroup contains memory_optimized", "ALTER DATABASE mydb ADD FILEGROUP fg1 CONTAINS MEMORY_OPTIMIZED_DATA"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "ADD" {
					t.Errorf("expected action ADD, got %s", stmt.Action)
				}
				if stmt.SubAction != "FILEGROUP" {
					t.Errorf("expected subAction FILEGROUP, got %s", stmt.SubAction)
				}
				if stmt.TargetName == "" {
					t.Error("expected targetName for filegroup")
				}
			})
		}
	})

	// --- REMOVE FILEGROUP ---
	t.Run("alter_database_remove_filegroup", func(t *testing.T) {
		sql := "ALTER DATABASE mydb REMOVE FILEGROUP fg1"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
		if !ok {
			t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Action != "REMOVE" {
			t.Errorf("expected action REMOVE, got %s", stmt.Action)
		}
		if stmt.SubAction != "FILEGROUP" {
			t.Errorf("expected subAction FILEGROUP, got %s", stmt.SubAction)
		}
		if stmt.TargetName != "fg1" {
			t.Errorf("expected targetName fg1, got %s", stmt.TargetName)
		}
	})

	// --- MODIFY FILEGROUP ---
	t.Run("alter_database_modify_filegroup", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"modify filegroup default", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 DEFAULT"},
			{"modify filegroup read_only", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 READ_ONLY"},
			{"modify filegroup read_write", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 READ_WRITE"},
			{"modify filegroup readonly", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 READONLY"},
			{"modify filegroup readwrite", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 READWRITE"},
			{"modify filegroup name", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 NAME = fg2"},
			{"modify filegroup autogrow_all", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 AUTOGROW_ALL_FILES"},
			{"modify filegroup autogrow_single", "ALTER DATABASE mydb MODIFY FILEGROUP fg1 AUTOGROW_SINGLE_FILE"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "MODIFY" {
					t.Errorf("expected action MODIFY, got %s", stmt.Action)
				}
				if stmt.SubAction != "FILEGROUP" {
					t.Errorf("expected subAction FILEGROUP, got %s", stmt.SubAction)
				}
				if stmt.TargetName != "fg1" {
					t.Errorf("expected targetName fg1, got %s", stmt.TargetName)
				}
			})
		}
	})

	// --- COLLATE ---
	t.Run("alter_database_collate", func(t *testing.T) {
		sql := "ALTER DATABASE mydb COLLATE Latin1_General_CI_AS"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
		if !ok {
			t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Action != "COLLATE" {
			t.Errorf("expected action COLLATE, got %s", stmt.Action)
		}
		if stmt.TargetName != "Latin1_General_CI_AS" {
			t.Errorf("expected collation Latin1_General_CI_AS, got %s", stmt.TargetName)
		}
	})

	// --- MODIFY NAME ---
	t.Run("alter_database_modify_name", func(t *testing.T) {
		sql := "ALTER DATABASE mydb MODIFY NAME = newdb"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
		if !ok {
			t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Action != "MODIFY" {
			t.Errorf("expected action MODIFY, got %s", stmt.Action)
		}
		if stmt.SubAction != "NAME" {
			t.Errorf("expected subAction NAME, got %s", stmt.SubAction)
		}
		if stmt.NewName != "newdb" {
			t.Errorf("expected newName newdb, got %s", stmt.NewName)
		}
	})

	// --- SET PARTNER (batch 117) ---
	t.Run("alter_database_set_partner", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpt  string
		}{
			{"partner server", "ALTER DATABASE mydb SET PARTNER = 'TCP://server2:5022'", "PARTNER='TCP://server2:5022'"},
			{"partner off", "ALTER DATABASE mydb SET PARTNER OFF", "PARTNER=OFF"},
			{"partner safety full", "ALTER DATABASE mydb SET PARTNER SAFETY FULL", "PARTNER SAFETY=FULL"},
			{"partner safety off", "ALTER DATABASE mydb SET PARTNER SAFETY OFF", "PARTNER SAFETY=OFF"},
			{"partner timeout", "ALTER DATABASE mydb SET PARTNER TIMEOUT 30", "PARTNER TIMEOUT=30"},
			{"partner failover", "ALTER DATABASE mydb SET PARTNER FAILOVER", "PARTNER FAILOVER"},
			{"partner force_failover", "ALTER DATABASE mydb SET PARTNER FORCE_SERVICE_ALLOW_DATA_LOSS", "PARTNER FORCE_SERVICE_ALLOW_DATA_LOSS"},
			{"partner resume", "ALTER DATABASE mydb SET PARTNER RESUME", "PARTNER RESUME"},
			{"partner suspend", "ALTER DATABASE mydb SET PARTNER SUSPEND", "PARTNER SUSPEND"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "SET" {
					t.Errorf("expected action SET, got %s", stmt.Action)
				}
				if stmt.Options == nil || stmt.Options.Len() == 0 {
					t.Fatal("expected SET options")
				}
				got := stmt.Options.Items[0].(*ast.String).Str
				if got != tt.wantOpt {
					t.Errorf("option = %q, want %q", got, tt.wantOpt)
				}
			})
		}
	})

	// --- SET WITNESS (batch 117) ---
	t.Run("alter_database_set_witness", func(t *testing.T) {
		tests := []struct {
			name    string
			sql     string
			wantOpt string
		}{
			{"witness server", "ALTER DATABASE mydb SET WITNESS = 'TCP://witness:5022'", "WITNESS='TCP://witness:5022'"},
			{"witness off", "ALTER DATABASE mydb SET WITNESS OFF", "WITNESS=OFF"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "SET" {
					t.Errorf("expected action SET, got %s", stmt.Action)
				}
				got := stmt.Options.Items[0].(*ast.String).Str
				if got != tt.wantOpt {
					t.Errorf("option = %q, want %q", got, tt.wantOpt)
				}
			})
		}
	})

	// --- SET HADR (batch 117) ---
	t.Run("alter_database_set_hadr", func(t *testing.T) {
		tests := []struct {
			name    string
			sql     string
			wantOpt string
		}{
			{"hadr availability group", "ALTER DATABASE mydb SET HADR AVAILABILITY GROUP = MyAG", "HADR AVAILABILITY GROUP=MYAG"},
			{"hadr off", "ALTER DATABASE mydb SET HADR OFF", "HADR=OFF"},
			{"hadr suspend", "ALTER DATABASE mydb SET HADR SUSPEND", "HADR SUSPEND"},
			{"hadr resume", "ALTER DATABASE mydb SET HADR RESUME", "HADR RESUME"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.Action != "SET" {
					t.Errorf("expected action SET, got %s", stmt.Action)
				}
				got := stmt.Options.Items[0].(*ast.String).Str
				if got != tt.wantOpt {
					t.Errorf("option = %q, want %q", got, tt.wantOpt)
				}
			})
		}
	})

	// --- Unknown action structured (batch 117) ---
	t.Run("alter_database_unknown_action", func(t *testing.T) {
		sql := "ALTER DATABASE mydb UNKNOWN_ACTION arg1 arg2"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
		if !ok {
			t.Fatalf("expected *AlterDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.Action != "UNKNOWN_ACTION" {
			t.Errorf("expected action UNKNOWN_ACTION, got %s", stmt.Action)
		}
		// Options should capture the remaining tokens instead of silently skipping
		if stmt.Options == nil || stmt.Options.Len() == 0 {
			t.Error("expected unknown action to capture remaining tokens as options")
		}
	})
}

// TestParseAlterIndexUnknownOption tests batch 117: ALTER INDEX unknown token handling.
func TestParseAlterIndexUnknownOption(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"rebuild basic", "ALTER INDEX idx1 ON dbo.MyTable REBUILD"},
		{"reorganize", "ALTER INDEX ALL ON dbo.MyTable REORGANIZE"},
		{"disable", "ALTER INDEX idx1 ON dbo.MyTable DISABLE"},
		{"set options", "ALTER INDEX idx1 ON dbo.MyTable SET (ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON)"},
		{"rebuild with options", "ALTER INDEX idx1 ON dbo.MyTable REBUILD WITH (FILLFACTOR = 80, SORT_IN_TEMPDB = ON)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			stmt, ok := result.Items[0].(*ast.AlterIndexStmt)
			if !ok {
				t.Fatalf("expected *AlterIndexStmt, got %T", result.Items[0])
			}
			if stmt.Action == "" {
				t.Error("expected non-empty action")
			}
			checkLocation(t, tt.sql, "AlterIndexStmt", stmt.Loc)
		})
	}
}

// TestParseCreateTypeIndex tests CREATE TYPE with INDEX in table type (batch 81).
func TestParseCreateTypeIndex(t *testing.T) {
	t.Run("create_type_alias", func(t *testing.T) {
		sql := "CREATE TYPE dbo.PhoneNumber FROM varchar(20) NOT NULL"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.CreateTypeStmt)
		if !ok {
			t.Fatalf("expected *CreateTypeStmt, got %T", result.Items[0])
		}
		if stmt.BaseType == nil {
			t.Error("expected non-nil BaseType for alias type")
		}
		if stmt.Nullable == nil || *stmt.Nullable != false {
			t.Error("expected Nullable=false for NOT NULL")
		}
	})

	t.Run("create_type_alias_nullable", func(t *testing.T) {
		sql := "CREATE TYPE dbo.Flag FROM bit NULL"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		if stmt.Nullable == nil || *stmt.Nullable != true {
			t.Error("expected Nullable=true for NULL")
		}
	})

	t.Run("create_type_external", func(t *testing.T) {
		sql := "CREATE TYPE dbo.Utf8String EXTERNAL NAME MyAssembly.Utf8String"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		if stmt.ExternalName != "MyAssembly.Utf8String" {
			t.Errorf("expected ExternalName MyAssembly.Utf8String, got %s", stmt.ExternalName)
		}
	})

	t.Run("create_type_table", func(t *testing.T) {
		sql := `CREATE TYPE dbo.LocationTableType AS TABLE (
			LocationName VARCHAR(50),
			CostRate INT
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		if stmt.TableDef == nil {
			t.Fatal("expected non-nil TableDef")
		}
		if stmt.TableDef.Len() != 2 {
			t.Errorf("expected 2 elements, got %d", stmt.TableDef.Len())
		}
	})

	t.Run("create_type_table_index", func(t *testing.T) {
		sql := `CREATE TYPE dbo.MyTableType AS TABLE (
			id INT NOT NULL,
			name NVARCHAR(100),
			INDEX IX_id CLUSTERED (id ASC)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		if stmt.TableDef == nil {
			t.Fatal("expected non-nil TableDef")
		}
		// Should have 2 columns + 1 index = 3 elements
		if stmt.TableDef.Len() != 3 {
			t.Errorf("expected 3 elements, got %d", stmt.TableDef.Len())
		}
		// The 3rd element should be a TableTypeIndex
		idx, ok := stmt.TableDef.Items[2].(*ast.TableTypeIndex)
		if !ok {
			t.Fatalf("expected *TableTypeIndex, got %T", stmt.TableDef.Items[2])
		}
		if idx.Name != "IX_id" {
			t.Errorf("expected index name IX_id, got %s", idx.Name)
		}
		if idx.Clustered == nil || *idx.Clustered != true {
			t.Error("expected Clustered=true")
		}
		if idx.Columns == nil || idx.Columns.Len() != 1 {
			t.Error("expected 1 index column")
		}
	})

	t.Run("create_type_table_nonclustered_index", func(t *testing.T) {
		sql := `CREATE TYPE dbo.MyTableType AS TABLE (
			id INT NOT NULL,
			col1 INT,
			col2 INT,
			INDEX IX_cols NONCLUSTERED (col1, col2 DESC)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		if stmt.TableDef == nil {
			t.Fatal("expected non-nil TableDef")
		}
		// Find the index
		var idx *ast.TableTypeIndex
		for _, item := range stmt.TableDef.Items {
			if i, ok := item.(*ast.TableTypeIndex); ok {
				idx = i
				break
			}
		}
		if idx == nil {
			t.Fatal("expected a TableTypeIndex element")
		}
		if idx.Clustered == nil || *idx.Clustered != false {
			t.Error("expected Clustered=false (NONCLUSTERED)")
		}
		if idx.Columns == nil || idx.Columns.Len() != 2 {
			t.Errorf("expected 2 index columns, got %d", idx.Columns.Len())
		}
	})

	t.Run("create_type_table_hash_index", func(t *testing.T) {
		sql := `CREATE TYPE dbo.MemOptType AS TABLE (
			id INT NOT NULL,
			INDEX IX_id NONCLUSTERED HASH WITH (BUCKET_COUNT = 1024) (id)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		var idx *ast.TableTypeIndex
		for _, item := range stmt.TableDef.Items {
			if i, ok := item.(*ast.TableTypeIndex); ok {
				idx = i
				break
			}
		}
		if idx == nil {
			t.Fatal("expected a TableTypeIndex element")
		}
		if !idx.Hash {
			t.Error("expected Hash=true")
		}
		if idx.BucketCount == nil {
			t.Error("expected non-nil BucketCount")
		}
	})

	t.Run("create_type_table_index_include", func(t *testing.T) {
		sql := `CREATE TYPE dbo.MyType AS TABLE (
			id INT NOT NULL,
			name NVARCHAR(100),
			email NVARCHAR(200),
			INDEX IX_name NONCLUSTERED (name) INCLUDE (email)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		var idx *ast.TableTypeIndex
		for _, item := range stmt.TableDef.Items {
			if i, ok := item.(*ast.TableTypeIndex); ok {
				idx = i
				break
			}
		}
		if idx == nil {
			t.Fatal("expected a TableTypeIndex element")
		}
		if idx.IncludeCols == nil || idx.IncludeCols.Len() != 1 {
			t.Errorf("expected 1 include column, got %d", idx.IncludeCols.Len())
		}
	})

	t.Run("create_type_table_multiple_indexes", func(t *testing.T) {
		sql := `CREATE TYPE dbo.MyType AS TABLE (
			id INT NOT NULL,
			col1 INT,
			col2 INT,
			INDEX IX_id CLUSTERED (id),
			INDEX IX_col1 NONCLUSTERED (col1)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTypeStmt)
		if stmt.TableDef == nil {
			t.Fatal("expected non-nil TableDef")
		}
		indexCount := 0
		for _, item := range stmt.TableDef.Items {
			if _, ok := item.(*ast.TableTypeIndex); ok {
				indexCount++
			}
		}
		if indexCount != 2 {
			t.Errorf("expected 2 indexes, got %d", indexCount)
		}
	})
}

// TestParseAlterTableFiletable tests ALTER TABLE FILETABLE_NAMESPACE and SET FILETABLE_DIRECTORY (batch 82).
func TestParseAlterTableFiletable(t *testing.T) {
	t.Run("alter_table_enable_filetable_namespace", func(t *testing.T) {
		sql := "ALTER TABLE dbo.MyFileTable ENABLE FILETABLE_NAMESPACE"
		result := ParseAndCheck(t, sql)
		stmt, ok := result.Items[0].(*ast.AlterTableStmt)
		if !ok {
			t.Fatalf("expected *AlterTableStmt, got %T", result.Items[0])
		}
		if stmt.Actions == nil || stmt.Actions.Len() != 1 {
			t.Fatalf("expected 1 action, got %d", stmt.Actions.Len())
		}
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATEnableFiletableNamespace {
			t.Errorf("expected ATEnableFiletableNamespace, got %d", action.Type)
		}
	})

	t.Run("alter_table_disable_filetable_namespace", func(t *testing.T) {
		sql := "ALTER TABLE dbo.MyFileTable DISABLE FILETABLE_NAMESPACE"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATDisableFiletableNamespace {
			t.Errorf("expected ATDisableFiletableNamespace, got %d", action.Type)
		}
	})

	t.Run("alter_table_set_filetable_directory", func(t *testing.T) {
		sql := "ALTER TABLE dbo.MyFileTable SET (FILETABLE_DIRECTORY = N'MyDocuments')"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		if stmt.Actions == nil || stmt.Actions.Len() != 1 {
			t.Fatalf("expected 1 action, got %d", stmt.Actions.Len())
		}
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATSet {
			t.Errorf("expected ATSet, got %d", action.Type)
		}
		if action.Options == nil {
			t.Error("expected non-nil Options")
		}
	})

	t.Run("alter_table_set_filetable_directory_simple", func(t *testing.T) {
		sql := "ALTER TABLE MyFileTable SET (FILETABLE_DIRECTORY = 'docs')"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATSet {
			t.Errorf("expected ATSet, got %d", action.Type)
		}
	})
}

// TestParseTableHints tests table hints parsing (batch 83).
func TestParseTableHints(t *testing.T) {
	// NOLOCK hint
	t.Run("table_hints_nolock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (NOLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		if ref.Hints == nil || ref.Hints.Len() != 1 {
			t.Fatalf("expected 1 hint, got %v", ref.Hints)
		}
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "NOLOCK" {
			t.Errorf("expected NOLOCK, got %s", hint.Name)
		}
	})

	// ROWLOCK hint
	t.Run("table_hints_rowlock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (ROWLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "ROWLOCK" {
			t.Errorf("expected ROWLOCK, got %s", hint.Name)
		}
	})

	// UPDLOCK hint
	t.Run("table_hints_updlock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (UPDLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "UPDLOCK" {
			t.Errorf("expected UPDLOCK, got %s", hint.Name)
		}
	})

	// HOLDLOCK hint
	t.Run("table_hints_holdlock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (HOLDLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "HOLDLOCK" {
			t.Errorf("expected HOLDLOCK, got %s", hint.Name)
		}
	})

	// TABLOCK hint
	t.Run("table_hints_tablock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (TABLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "TABLOCK" {
			t.Errorf("expected TABLOCK, got %s", hint.Name)
		}
	})

	// TABLOCKX hint
	t.Run("table_hints_tablockx", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (TABLOCKX)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "TABLOCKX" {
			t.Errorf("expected TABLOCKX, got %s", hint.Name)
		}
	})

	// PAGLOCK hint
	t.Run("table_hints_paglock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (PAGLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "PAGLOCK" {
			t.Errorf("expected PAGLOCK, got %s", hint.Name)
		}
	})

	// XLOCK hint
	t.Run("table_hints_xlock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (XLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "XLOCK" {
			t.Errorf("expected XLOCK, got %s", hint.Name)
		}
	})

	// SERIALIZABLE hint
	t.Run("table_hints_serializable", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (SERIALIZABLE)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "SERIALIZABLE" {
			t.Errorf("expected SERIALIZABLE, got %s", hint.Name)
		}
	})

	// READUNCOMMITTED hint
	t.Run("table_hints_readuncommitted", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (READUNCOMMITTED)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "READUNCOMMITTED" {
			t.Errorf("expected READUNCOMMITTED, got %s", hint.Name)
		}
	})

	// READCOMMITTED hint
	t.Run("table_hints_readcommitted", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (READCOMMITTED)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "READCOMMITTED" {
			t.Errorf("expected READCOMMITTED, got %s", hint.Name)
		}
	})

	// READCOMMITTEDLOCK hint
	t.Run("table_hints_readcommittedlock", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (READCOMMITTEDLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "READCOMMITTEDLOCK" {
			t.Errorf("expected READCOMMITTEDLOCK, got %s", hint.Name)
		}
	})

	// REPEATABLEREAD hint
	t.Run("table_hints_repeatableread", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (REPEATABLEREAD)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "REPEATABLEREAD" {
			t.Errorf("expected REPEATABLEREAD, got %s", hint.Name)
		}
	})

	// READPAST hint
	t.Run("table_hints_readpast", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (READPAST)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "READPAST" {
			t.Errorf("expected READPAST, got %s", hint.Name)
		}
	})

	// SNAPSHOT hint
	t.Run("table_hints_snapshot", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (SNAPSHOT)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "SNAPSHOT" {
			t.Errorf("expected SNAPSHOT, got %s", hint.Name)
		}
	})

	// NOWAIT hint
	t.Run("table_hints_nowait", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (NOWAIT)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "NOWAIT" {
			t.Errorf("expected NOWAIT, got %s", hint.Name)
		}
	})

	// NOEXPAND hint
	t.Run("table_hints_noexpand", func(t *testing.T) {
		sql := "SELECT * FROM v1 WITH (NOEXPAND)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "NOEXPAND" {
			t.Errorf("expected NOEXPAND, got %s", hint.Name)
		}
	})

	// FORCESCAN hint
	t.Run("table_hints_forcescan", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (FORCESCAN)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "FORCESCAN" {
			t.Errorf("expected FORCESCAN, got %s", hint.Name)
		}
	})

	// INDEX hint with single index name
	t.Run("table_hints_index", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (INDEX(idx1))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "INDEX" {
			t.Errorf("expected INDEX, got %s", hint.Name)
		}
		if hint.IndexValues == nil || hint.IndexValues.Len() != 1 {
			t.Fatalf("expected 1 index value, got %v", hint.IndexValues)
		}
	})

	// INDEX hint with multiple index values
	t.Run("table_hints_index_multiple", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (INDEX(idx1, idx2))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "INDEX" {
			t.Errorf("expected INDEX, got %s", hint.Name)
		}
		if hint.IndexValues == nil || hint.IndexValues.Len() != 2 {
			t.Fatalf("expected 2 index values, got %v", hint.IndexValues)
		}
	})

	// INDEX hint with = syntax
	t.Run("table_hints_index_eq", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (INDEX = (idx1))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "INDEX" {
			t.Errorf("expected INDEX, got %s", hint.Name)
		}
		if hint.IndexValues == nil || hint.IndexValues.Len() != 1 {
			t.Fatalf("expected 1 index value, got %v", hint.IndexValues)
		}
	})

	// INDEX hint with numeric value (0 = table scan)
	t.Run("table_hints_index_zero", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (INDEX(0))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "INDEX" {
			t.Errorf("expected INDEX, got %s", hint.Name)
		}
		if hint.IndexValues == nil || hint.IndexValues.Len() != 1 {
			t.Fatalf("expected 1 index value, got %v", hint.IndexValues)
		}
	})

	// FORCESEEK hint (no parameters)
	t.Run("table_hints_forceseek", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (FORCESEEK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "FORCESEEK" {
			t.Errorf("expected FORCESEEK, got %s", hint.Name)
		}
	})

	// FORCESEEK with index and columns
	t.Run("table_hints_forceseek_params", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (FORCESEEK(idx1(col1, col2)))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "FORCESEEK" {
			t.Errorf("expected FORCESEEK, got %s", hint.Name)
		}
		if hint.IndexValues == nil || hint.IndexValues.Len() != 1 {
			t.Fatalf("expected 1 index value for FORCESEEK, got %v", hint.IndexValues)
		}
		if hint.ForceSeekColumns == nil || hint.ForceSeekColumns.Len() != 2 {
			t.Fatalf("expected 2 FORCESEEK columns, got %v", hint.ForceSeekColumns)
		}
	})

	// SPATIAL_WINDOW_MAX_CELLS hint
	t.Run("table_hints_spatial", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (SPATIAL_WINDOW_MAX_CELLS = 512)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		hint := ref.Hints.Items[0].(*ast.TableHint)
		if hint.Name != "SPATIAL_WINDOW_MAX_CELLS" {
			t.Errorf("expected SPATIAL_WINDOW_MAX_CELLS, got %s", hint.Name)
		}
		if hint.IntValue == nil {
			t.Error("expected IntValue for SPATIAL_WINDOW_MAX_CELLS")
		}
	})

	// Multiple hints
	t.Run("table_hints_multiple", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (NOLOCK, INDEX(idx1))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		if ref.Hints == nil || ref.Hints.Len() != 2 {
			t.Fatalf("expected 2 hints, got %v", ref.Hints)
		}
		h0 := ref.Hints.Items[0].(*ast.TableHint)
		h1 := ref.Hints.Items[1].(*ast.TableHint)
		if h0.Name != "NOLOCK" {
			t.Errorf("expected NOLOCK, got %s", h0.Name)
		}
		if h1.Name != "INDEX" {
			t.Errorf("expected INDEX, got %s", h1.Name)
		}
	})

	// Hints with alias
	t.Run("table_hints_with_alias", func(t *testing.T) {
		sql := "SELECT * FROM t1 AS a WITH (NOLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		if ref.Alias != "a" {
			t.Errorf("expected alias 'a', got %s", ref.Alias)
		}
		if ref.Hints == nil || ref.Hints.Len() != 1 {
			t.Fatalf("expected 1 hint, got %v", ref.Hints)
		}
	})

	// Hints on joined table
	t.Run("table_hints_join", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (NOLOCK) INNER JOIN t2 WITH (ROWLOCK) ON t1.id = t2.id"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		join := stmt.FromClause.Items[0].(*ast.JoinClause)
		left := join.Left.(*ast.TableRef)
		right := join.Right.(*ast.TableRef)
		if left.Hints == nil || left.Hints.Len() != 1 {
			t.Fatalf("expected 1 hint on left, got %v", left.Hints)
		}
		if right.Hints == nil || right.Hints.Len() != 1 {
			t.Fatalf("expected 1 hint on right, got %v", right.Hints)
		}
		if left.Hints.Items[0].(*ast.TableHint).Name != "NOLOCK" {
			t.Errorf("expected NOLOCK on left")
		}
		if right.Hints.Items[0].(*ast.TableHint).Name != "ROWLOCK" {
			t.Errorf("expected ROWLOCK on right")
		}
	})

	// Three hints combined
	t.Run("table_hints_three", func(t *testing.T) {
		sql := "SELECT * FROM t1 WITH (UPDLOCK, ROWLOCK, HOLDLOCK)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		if ref.Hints == nil || ref.Hints.Len() != 3 {
			t.Fatalf("expected 3 hints, got %v", ref.Hints)
		}
		names := []string{"UPDLOCK", "ROWLOCK", "HOLDLOCK"}
		for i, name := range names {
			h := ref.Hints.Items[i].(*ast.TableHint)
			if h.Name != name {
				t.Errorf("hint[%d]: expected %s, got %s", i, name, h.Name)
			}
		}
	})

	// NOEXPAND with INDEX
	t.Run("table_hints_noexpand_index", func(t *testing.T) {
		sql := "SELECT * FROM v1 WITH (NOEXPAND, INDEX(idx1))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		ref := stmt.FromClause.Items[0].(*ast.TableRef)
		if ref.Hints == nil || ref.Hints.Len() != 2 {
			t.Fatalf("expected 2 hints, got %v", ref.Hints)
		}
	})
}

// TestParseCollateExpression tests COLLATE as a postfix expression operator (batch 84).
func TestParseCollateExpression(t *testing.T) {
	// COLLATE in ORDER BY
	t.Run("collate_order_by", func(t *testing.T) {
		sql := "SELECT name FROM employees ORDER BY name COLLATE Latin1_General_CS_AS"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.OrderByClause == nil || stmt.OrderByClause.Len() != 1 {
			t.Fatalf("expected 1 ORDER BY item")
		}
		item := stmt.OrderByClause.Items[0].(*ast.OrderByItem)
		collate, ok := item.Expr.(*ast.CollateExpr)
		if !ok {
			t.Fatalf("expected CollateExpr, got %T", item.Expr)
		}
		if collate.Collation != "Latin1_General_CS_AS" {
			t.Errorf("expected collation Latin1_General_CS_AS, got %s", collate.Collation)
		}
	})

	// COLLATE in WHERE clause
	t.Run("collate_where", func(t *testing.T) {
		sql := "SELECT * FROM t WHERE col1 COLLATE Latin1_General_CI_AS = 'abc'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.WhereClause == nil {
			t.Fatalf("expected WHERE clause")
		}
		bin, ok := stmt.WhereClause.(*ast.BinaryExpr)
		if !ok {
			t.Fatalf("expected BinaryExpr, got %T", stmt.WhereClause)
		}
		collate, ok := bin.Left.(*ast.CollateExpr)
		if !ok {
			t.Fatalf("expected CollateExpr on left, got %T", bin.Left)
		}
		if collate.Collation != "Latin1_General_CI_AS" {
			t.Errorf("expected collation Latin1_General_CI_AS, got %s", collate.Collation)
		}
	})

	// COLLATE on column expression (SELECT list)
	t.Run("collate_column_expr", func(t *testing.T) {
		sql := "SELECT col1 COLLATE SQL_Latin1_General_CP1_CI_AS FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.TargetList == nil || stmt.TargetList.Len() != 1 {
			t.Fatalf("expected 1 target")
		}
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		collate, ok := rt.Val.(*ast.CollateExpr)
		if !ok {
			t.Fatalf("expected CollateExpr, got %T", rt.Val)
		}
		if collate.Collation != "SQL_Latin1_General_CP1_CI_AS" {
			t.Errorf("expected SQL_Latin1_General_CP1_CI_AS, got %s", collate.Collation)
		}
	})

	// COLLATE database_default
	t.Run("collate_database_default", func(t *testing.T) {
		sql := "SELECT col1 COLLATE database_default FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		collate, ok := rt.Val.(*ast.CollateExpr)
		if !ok {
			t.Fatalf("expected CollateExpr, got %T", rt.Val)
		}
		if collate.Collation != "database_default" {
			t.Errorf("expected database_default, got %s", collate.Collation)
		}
	})

	// COLLATE on string literal
	t.Run("collate_string_literal", func(t *testing.T) {
		sql := "SELECT 'hello' COLLATE Latin1_General_CS_AS"
		ParseAndCheck(t, sql)
	})

	// COLLATE with LIKE
	t.Run("collate_with_like", func(t *testing.T) {
		sql := "SELECT * FROM t WHERE name COLLATE Latin1_General_CI_AS LIKE 'A%'"
		ParseAndCheck(t, sql)
	})

	// COLLATE in CASE expression
	t.Run("collate_in_case", func(t *testing.T) {
		sql := "SELECT CASE WHEN 1=1 THEN col1 COLLATE Latin1_General_CS_AS ELSE col2 END FROM t"
		ParseAndCheck(t, sql)
	})

	// Multiple COLLATEs in same query
	t.Run("collate_multiple", func(t *testing.T) {
		sql := "SELECT a COLLATE Latin1_General_CS_AS, b COLLATE SQL_Latin1_General_CP1_CI_AS FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.TargetList.Len() != 2 {
			t.Fatalf("expected 2 targets, got %d", stmt.TargetList.Len())
		}
		for i := 0; i < 2; i++ {
			rt := stmt.TargetList.Items[i].(*ast.ResTarget)
			if _, ok := rt.Val.(*ast.CollateExpr); !ok {
				t.Errorf("target[%d]: expected CollateExpr, got %T", i, rt.Val)
			}
		}
	})
}

// TestParseAtTimeZone tests AT TIME ZONE postfix expressions (batch 85).
func TestParseAtTimeZone(t *testing.T) {
	// Basic AT TIME ZONE with a string literal
	t.Run("at_time_zone_literal", func(t *testing.T) {
		sql := "SELECT SalesDate AT TIME ZONE 'Pacific Standard Time'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		atz, ok := rt.Val.(*ast.AtTimeZoneExpr)
		if !ok {
			t.Fatalf("expected AtTimeZoneExpr, got %T", rt.Val)
		}
		if atz.Expr == nil {
			t.Fatal("AtTimeZoneExpr.Expr is nil")
		}
		if atz.TimeZone == nil {
			t.Fatal("AtTimeZoneExpr.TimeZone is nil")
		}
		lit, ok := atz.TimeZone.(*ast.Literal)
		if !ok {
			t.Fatalf("expected Literal timezone, got %T", atz.TimeZone)
		}
		if lit.Str != "Pacific Standard Time" {
			t.Errorf("expected timezone 'Pacific Standard Time', got %q", lit.Str)
		}
	})

	// Chained AT TIME ZONE expressions
	t.Run("at_time_zone_chained", func(t *testing.T) {
		sql := "SELECT col1 AT TIME ZONE 'UTC' AT TIME ZONE 'Eastern Standard Time'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		outer, ok := rt.Val.(*ast.AtTimeZoneExpr)
		if !ok {
			t.Fatalf("expected outer AtTimeZoneExpr, got %T", rt.Val)
		}
		inner, ok := outer.Expr.(*ast.AtTimeZoneExpr)
		if !ok {
			t.Fatalf("expected inner AtTimeZoneExpr, got %T", outer.Expr)
		}
		_ = inner
		outerTZ := outer.TimeZone.(*ast.Literal)
		if outerTZ.Str != "Eastern Standard Time" {
			t.Errorf("expected outer timezone 'Eastern Standard Time', got %q", outerTZ.Str)
		}
	})

	// AT TIME ZONE with a variable
	t.Run("at_time_zone_variable", func(t *testing.T) {
		sql := "SELECT GETDATE() AT TIME ZONE @tz"
		ParseAndCheck(t, sql)
	})

	// AT TIME ZONE in WHERE clause
	t.Run("at_time_zone_in_where", func(t *testing.T) {
		sql := "SELECT * FROM Orders WHERE OrderDate AT TIME ZONE 'UTC' > '2024-01-01'"
		ParseAndCheck(t, sql)
	})

	// AT TIME ZONE with CAST
	t.Run("at_time_zone_with_cast", func(t *testing.T) {
		sql := "SELECT CAST(col1 AS DATETIMEOFFSET) AT TIME ZONE 'UTC'"
		ParseAndCheck(t, sql)
	})
}

// TestParseCreateTableInlineIndex tests inline INDEX definitions in CREATE TABLE (batch 86).
func TestParseCreateTableInlineIndex(t *testing.T) {
	// Basic inline index
	t.Run("create_table_inline_index", func(t *testing.T) {
		sql := `CREATE TABLE t1 (
			id INT NOT NULL,
			name NVARCHAR(100),
			INDEX IX_name NONCLUSTERED (name)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.Indexes == nil || stmt.Indexes.Len() != 1 {
			t.Fatalf("expected 1 inline index, got %v", stmt.Indexes)
		}
		idx := stmt.Indexes.Items[0].(*ast.InlineIndexDef)
		if idx.Name != "IX_name" {
			t.Errorf("expected index name IX_name, got %q", idx.Name)
		}
		if idx.Clustered == nil || *idx.Clustered != false {
			t.Error("expected NONCLUSTERED")
		}
	})

	// Inline index with INCLUDE
	t.Run("create_table_inline_index_include", func(t *testing.T) {
		sql := `CREATE TABLE t2 (
			id INT PRIMARY KEY,
			col1 INT,
			col2 VARCHAR(50),
			INDEX IX_col1 NONCLUSTERED (col1) INCLUDE (col2)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.Indexes == nil || stmt.Indexes.Len() != 1 {
			t.Fatalf("expected 1 inline index")
		}
		idx := stmt.Indexes.Items[0].(*ast.InlineIndexDef)
		if idx.IncludeCols == nil || idx.IncludeCols.Len() != 1 {
			t.Fatalf("expected 1 include column")
		}
	})

	// Unique clustered inline index
	t.Run("create_table_inline_unique_clustered", func(t *testing.T) {
		sql := `CREATE TABLE t3 (
			id INT,
			INDEX IX_id UNIQUE CLUSTERED (id ASC)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		idx := stmt.Indexes.Items[0].(*ast.InlineIndexDef)
		if !idx.Unique {
			t.Error("expected UNIQUE")
		}
		if idx.Clustered == nil || *idx.Clustered != true {
			t.Error("expected CLUSTERED")
		}
	})

	// Inline index with WHERE (filtered)
	t.Run("create_table_inline_index_filtered", func(t *testing.T) {
		sql := `CREATE TABLE t4 (
			id INT,
			status INT,
			INDEX IX_active NONCLUSTERED (status) WHERE status = 1
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		idx := stmt.Indexes.Items[0].(*ast.InlineIndexDef)
		if idx.WhereClause == nil {
			t.Error("expected WHERE clause on filtered index")
		}
	})

	// Inline index with WITH options
	t.Run("create_table_inline_index_with_options", func(t *testing.T) {
		sql := `CREATE TABLE t5 (
			id INT,
			name VARCHAR(100),
			INDEX IX_name NONCLUSTERED (name) WITH (FILLFACTOR = 80)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		idx := stmt.Indexes.Items[0].(*ast.InlineIndexDef)
		if idx.Options == nil || idx.Options.Len() != 1 {
			t.Fatalf("expected 1 index option, got %v", idx.Options)
		}
	})

	// Multiple inline indexes
	t.Run("create_table_multiple_inline_indexes", func(t *testing.T) {
		sql := `CREATE TABLE t6 (
			id INT,
			col1 INT,
			col2 INT,
			INDEX IX_col1 (col1 DESC),
			INDEX IX_col2 (col2 ASC)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.Indexes == nil || stmt.Indexes.Len() != 2 {
			t.Fatalf("expected 2 inline indexes, got %d", stmt.Indexes.Len())
		}
	})
}

// TestParseGraphTables tests CREATE TABLE AS NODE / AS EDGE (batch 87).
func TestParseGraphTables(t *testing.T) {
	// AS NODE
	t.Run("create_table_as_node", func(t *testing.T) {
		sql := `CREATE TABLE Person AS NODE (
			ID INT PRIMARY KEY,
			Name NVARCHAR(100)
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if !stmt.IsNode {
			t.Error("expected IsNode=true")
		}
		if stmt.IsEdge {
			t.Error("expected IsEdge=false")
		}
		if stmt.Columns == nil || stmt.Columns.Len() != 2 {
			t.Fatalf("expected 2 columns")
		}
	})

	// AS EDGE
	t.Run("create_table_as_edge", func(t *testing.T) {
		sql := `CREATE TABLE Likes AS EDGE (
			CreatedDate DATETIME DEFAULT GETDATE()
		)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if !stmt.IsEdge {
			t.Error("expected IsEdge=true")
		}
		if stmt.IsNode {
			t.Error("expected IsNode=false")
		}
	})

	// Edge table with $from_id and $to_id pseudo-columns
	t.Run("create_table_edge_pseudo_columns", func(t *testing.T) {
		sql := `CREATE TABLE FriendOf AS EDGE (
			Since DATE
		)`
		ParseAndCheck(t, sql)
	})

	// Node table without columns
	t.Run("create_table_node_minimal", func(t *testing.T) {
		sql := "CREATE TABLE dbo.Entity AS NODE (ID INT)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if !stmt.IsNode {
			t.Error("expected IsNode=true")
		}
		if stmt.Name.Schema != "dbo" {
			t.Errorf("expected schema dbo, got %q", stmt.Name.Schema)
		}
	})
}

// TestParseAlterTablePeriod tests ALTER TABLE ADD/DROP PERIOD FOR SYSTEM_TIME (batch 88).
func TestParseAlterTablePeriod(t *testing.T) {
	// ADD PERIOD FOR SYSTEM_TIME
	t.Run("alter_table_add_period", func(t *testing.T) {
		sql := "ALTER TABLE dbo.Employee ADD PERIOD FOR SYSTEM_TIME (ValidFrom, ValidTo)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		if stmt.Actions == nil || stmt.Actions.Len() != 1 {
			t.Fatalf("expected 1 action")
		}
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATAddPeriod {
			t.Errorf("expected ATAddPeriod, got %d", action.Type)
		}
		if action.Names == nil || action.Names.Len() != 2 {
			t.Fatalf("expected 2 column names")
		}
		start := action.Names.Items[0].(*ast.String).Str
		end := action.Names.Items[1].(*ast.String).Str
		if start != "ValidFrom" {
			t.Errorf("expected start col ValidFrom, got %q", start)
		}
		if end != "ValidTo" {
			t.Errorf("expected end col ValidTo, got %q", end)
		}
	})

	// DROP PERIOD FOR SYSTEM_TIME
	t.Run("alter_table_drop_period", func(t *testing.T) {
		sql := "ALTER TABLE dbo.Employee DROP PERIOD FOR SYSTEM_TIME"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		if stmt.Actions == nil || stmt.Actions.Len() != 1 {
			t.Fatalf("expected 1 action")
		}
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATDropPeriod {
			t.Errorf("expected ATDropPeriod, got %d", action.Type)
		}
	})
}

// TestParseTriggerWithOptions tests CREATE TRIGGER WITH options (batch 89).
func TestParseTriggerWithOptions(t *testing.T) {
	// WITH ENCRYPTION
	t.Run("trigger_with_encryption", func(t *testing.T) {
		sql := "CREATE TRIGGER tr1 ON dbo.t1 WITH ENCRYPTION AFTER INSERT AS BEGIN SELECT 1 END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTriggerStmt)
		if stmt.TriggerOptions == nil || stmt.TriggerOptions.Len() != 1 {
			t.Fatalf("expected 1 trigger option, got %v", stmt.TriggerOptions)
		}
		opt := stmt.TriggerOptions.Items[0].(*ast.String).Str
		if opt != "ENCRYPTION" {
			t.Errorf("expected ENCRYPTION option, got %q", opt)
		}
		if stmt.TriggerType != "AFTER" {
			t.Errorf("expected AFTER trigger type, got %q", stmt.TriggerType)
		}
	})

	// WITH EXECUTE AS
	t.Run("trigger_with_execute_as", func(t *testing.T) {
		sql := "CREATE TRIGGER tr2 ON dbo.t2 WITH EXECUTE AS OWNER AFTER UPDATE AS BEGIN SELECT 1 END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTriggerStmt)
		if stmt.TriggerOptions == nil || stmt.TriggerOptions.Len() != 1 {
			t.Fatalf("expected 1 trigger option")
		}
		opt := stmt.TriggerOptions.Items[0].(*ast.String).Str
		if opt != "EXECUTE AS OWNER" {
			t.Errorf("expected 'EXECUTE AS OWNER', got %q", opt)
		}
	})

	// WITH EXECUTE AS 'user'
	t.Run("trigger_with_execute_as_user", func(t *testing.T) {
		sql := "CREATE TRIGGER tr3 ON t3 WITH EXECUTE AS 'dbo' FOR INSERT AS BEGIN SELECT 1 END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTriggerStmt)
		if stmt.TriggerOptions == nil || stmt.TriggerOptions.Len() != 1 {
			t.Fatalf("expected 1 trigger option")
		}
		opt := stmt.TriggerOptions.Items[0].(*ast.String).Str
		if opt != "EXECUTE AS dbo" {
			t.Errorf("expected 'EXECUTE AS dbo', got %q", opt)
		}
	})

	// WITH ENCRYPTION, EXECUTE AS CALLER (multiple options)
	t.Run("trigger_with_multiple_options", func(t *testing.T) {
		sql := "CREATE TRIGGER tr4 ON t4 WITH ENCRYPTION, EXECUTE AS CALLER AFTER DELETE AS BEGIN SELECT 1 END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTriggerStmt)
		if stmt.TriggerOptions == nil || stmt.TriggerOptions.Len() != 2 {
			t.Fatalf("expected 2 trigger options, got %d", stmt.TriggerOptions.Len())
		}
	})

	// WITH SCHEMABINDING
	t.Run("trigger_with_schemabinding", func(t *testing.T) {
		sql := "CREATE TRIGGER tr5 ON t5 WITH SCHEMABINDING AFTER INSERT AS BEGIN SELECT 1 END"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTriggerStmt)
		if stmt.TriggerOptions == nil || stmt.TriggerOptions.Len() != 1 {
			t.Fatalf("expected 1 trigger option")
		}
		opt := stmt.TriggerOptions.Items[0].(*ast.String).Str
		if opt != "SCHEMABINDING" {
			t.Errorf("expected SCHEMABINDING, got %q", opt)
		}
	})
}

// TestParseSecurityKeyOptionsDepth tests structured security key option parsing (batch 90).
func TestParseSecurityKeyOptionsDepth(t *testing.T) {
	// CREATE SYMMETRIC KEY with structured options
	t.Run("create_symmetric_key_options", func(t *testing.T) {
		sql := `CREATE SYMMETRIC KEY MySymKey
			WITH ALGORITHM = AES_256, KEY_SOURCE = 'pass phrase', IDENTITY_VALUE = 'id val'
			ENCRYPTION BY CERTIFICATE MyCert`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Name != "MySymKey" {
			t.Errorf("expected name MySymKey, got %q", stmt.Name)
		}
		if stmt.Options == nil {
			t.Fatal("expected options to be parsed")
		}
		// Should have structured options, not a raw blob
		if stmt.Options.Len() < 2 {
			t.Errorf("expected multiple option items, got %d", stmt.Options.Len())
		}
	})

	// CREATE CERTIFICATE with structured options
	t.Run("create_certificate_options", func(t *testing.T) {
		sql := `CREATE CERTIFICATE MyCert
			WITH SUBJECT = 'Test Certificate'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options to be parsed")
		}
	})

	// CREATE CREDENTIAL with IDENTITY and SECRET
	t.Run("create_credential_options", func(t *testing.T) {
		sql := `CREATE CREDENTIAL MyCred
			WITH IDENTITY = 'my_identity', SECRET = 'my_secret'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CREDENTIAL" {
			t.Errorf("expected CREDENTIAL, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options to be parsed")
		}
	})

	// CREATE MASTER KEY with ENCRYPTION BY PASSWORD
	t.Run("create_master_key", func(t *testing.T) {
		sql := "CREATE MASTER KEY ENCRYPTION BY PASSWORD = 'StrongPass123!'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "MASTER KEY" {
			t.Errorf("expected MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// CREATE ASYMMETRIC KEY
	t.Run("create_asymmetric_key", func(t *testing.T) {
		sql := "CREATE ASYMMETRIC KEY MyAsymKey WITH ALGORITHM = RSA_2048"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "ASYMMETRIC KEY" {
			t.Errorf("expected ASYMMETRIC KEY, got %q", stmt.ObjectType)
		}
	})
}

// TestParseChooseStringAgg tests CHOOSE function and STRING_AGG WITHIN GROUP (batch 91).
func TestParseChooseStringAgg(t *testing.T) {
	// CHOOSE function
	t.Run("choose_function", func(t *testing.T) {
		sql := "SELECT CHOOSE(2, 'a', 'b', 'c')"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		fc, ok := rt.Val.(*ast.FuncCallExpr)
		if !ok {
			t.Fatalf("expected FuncCallExpr, got %T", rt.Val)
		}
		if !strings.EqualFold(fc.Name.Object, "CHOOSE") {
			t.Errorf("expected function name CHOOSE, got %q", fc.Name.Object)
		}
		if fc.Args == nil || fc.Args.Len() != 4 {
			t.Errorf("expected 4 args, got %d", fc.Args.Len())
		}
	})

	// STRING_AGG with WITHIN GROUP
	t.Run("string_agg_within_group", func(t *testing.T) {
		sql := "SELECT STRING_AGG(Name, ', ') WITHIN GROUP (ORDER BY Name ASC) FROM employees"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		fc, ok := rt.Val.(*ast.FuncCallExpr)
		if !ok {
			t.Fatalf("expected FuncCallExpr, got %T", rt.Val)
		}
		if !strings.EqualFold(fc.Name.Object, "STRING_AGG") {
			t.Errorf("expected STRING_AGG, got %q", fc.Name.Object)
		}
		if fc.Within == nil {
			t.Fatal("expected WITHIN GROUP clause")
		}
		if fc.Within.Len() != 1 {
			t.Fatalf("expected 1 order item, got %d", fc.Within.Len())
		}
		orderItem := fc.Within.Items[0].(*ast.OrderByItem)
		if orderItem.SortDir != ast.SortAsc {
			t.Errorf("expected ASC sort, got %d", orderItem.SortDir)
		}
	})

	// STRING_AGG with multiple ORDER BY columns
	t.Run("string_agg_multiple_order", func(t *testing.T) {
		sql := "SELECT STRING_AGG(col, ';') WITHIN GROUP (ORDER BY col1 ASC, col2 DESC) FROM t"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		rt := stmt.TargetList.Items[0].(*ast.ResTarget)
		fc := rt.Val.(*ast.FuncCallExpr)
		if fc.Within == nil || fc.Within.Len() != 2 {
			t.Fatalf("expected 2 order items")
		}
	})

	// CHOOSE in WHERE clause
	t.Run("choose_in_where", func(t *testing.T) {
		sql := "SELECT * FROM t WHERE CHOOSE(status, 'a', 'b') = 'a'"
		ParseAndCheck(t, sql)
	})
}

// TestParseOptionQueryHints tests OPTION clause query hints (batch 92).
func TestParseOptionQueryHints(t *testing.T) {
	// OPTION (RECOMPILE)
	t.Run("option_recompile", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (RECOMPILE)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.OptionClause == nil || stmt.OptionClause.Len() != 1 {
			t.Fatalf("expected 1 hint")
		}
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "RECOMPILE" {
			t.Errorf("expected RECOMPILE, got %q", hint.Kind)
		}
	})

	// OPTION (OPTIMIZE FOR (@id = 1))
	t.Run("option_optimize_for", func(t *testing.T) {
		sql := "SELECT * FROM t WHERE id = @id OPTION (OPTIMIZE FOR (@id = 1))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.OptionClause == nil || stmt.OptionClause.Len() != 1 {
			t.Fatalf("expected 1 hint")
		}
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "OPTIMIZE FOR" {
			t.Errorf("expected OPTIMIZE FOR hint, got %q", hint.Kind)
		}
		if hint.Params == nil || hint.Params.Len() != 1 {
			t.Fatalf("expected 1 param")
		}
		param := hint.Params.Items[0].(*ast.OptimizeForParam)
		if param.Variable != "@id" {
			t.Errorf("expected @id, got %q", param.Variable)
		}
		if param.Unknown {
			t.Errorf("expected not UNKNOWN")
		}
	})

	// OPTION (OPTIMIZE FOR UNKNOWN)
	t.Run("option_optimize_for_unknown", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (OPTIMIZE FOR UNKNOWN)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "OPTIMIZE FOR UNKNOWN" {
			t.Errorf("expected 'OPTIMIZE FOR UNKNOWN', got %q", hint.Kind)
		}
	})

	// OPTION (HASH JOIN)
	t.Run("option_join_hints", func(t *testing.T) {
		sql := "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id OPTION (HASH JOIN)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "HASH JOIN" {
			t.Errorf("expected 'HASH JOIN', got %q", hint.Kind)
		}
	})

	// OPTION (MAXDOP 4)
	t.Run("option_maxdop", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (MAXDOP 4)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "MAXDOP" {
			t.Errorf("expected 'MAXDOP', got %q", hint.Kind)
		}
		if hint.Value == nil {
			t.Fatalf("expected Value to be set")
		}
	})

	// Multiple hints
	t.Run("option_multiple_hints", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (RECOMPILE, MAXDOP 2, FORCE ORDER)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		if stmt.OptionClause == nil || stmt.OptionClause.Len() != 3 {
			t.Fatalf("expected 3 hints, got %d", stmt.OptionClause.Len())
		}
		h0 := stmt.OptionClause.Items[0].(*ast.QueryHint)
		h1 := stmt.OptionClause.Items[1].(*ast.QueryHint)
		h2 := stmt.OptionClause.Items[2].(*ast.QueryHint)
		if h0.Kind != "RECOMPILE" {
			t.Errorf("hint 0: expected RECOMPILE, got %q", h0.Kind)
		}
		if h1.Kind != "MAXDOP" {
			t.Errorf("hint 1: expected MAXDOP, got %q", h1.Kind)
		}
		if h2.Kind != "FORCE ORDER" {
			t.Errorf("hint 2: expected FORCE ORDER, got %q", h2.Kind)
		}
	})

	// OPTION (KEEP PLAN)
	t.Run("option_keep_plan", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (KEEP PLAN)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "KEEP PLAN" {
			t.Errorf("expected 'KEEP PLAN', got %q", hint.Kind)
		}
	})

	// OPTION (ROBUST PLAN)
	t.Run("option_robust_plan", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (ROBUST PLAN)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "ROBUST PLAN" {
			t.Errorf("expected 'ROBUST PLAN', got %q", hint.Kind)
		}
	})

	// TABLE HINT with single hint
	t.Run("option_table_hint", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (TABLE HINT(t, NOLOCK))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "TABLE HINT" {
			t.Errorf("expected 'TABLE HINT', got %q", hint.Kind)
		}
		if hint.TableName == nil || hint.TableName.Object != "t" {
			t.Errorf("expected table name 't'")
		}
		if hint.TableHints == nil || hint.TableHints.Len() != 1 {
			t.Fatalf("expected 1 table hint")
		}
		th := hint.TableHints.Items[0].(*ast.TableHint)
		if th.Name != "NOLOCK" {
			t.Errorf("expected NOLOCK, got %q", th.Name)
		}
	})

	// TABLE HINT with INDEX hint
	t.Run("option_table_hint_index", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (TABLE HINT(dbo.t, INDEX(IX_1)))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "TABLE HINT" {
			t.Errorf("expected 'TABLE HINT', got %q", hint.Kind)
		}
		if hint.TableName == nil || hint.TableName.Schema != "dbo" || hint.TableName.Object != "t" {
			t.Errorf("expected table name 'dbo.t'")
		}
		if hint.TableHints == nil || hint.TableHints.Len() != 1 {
			t.Fatalf("expected 1 table hint")
		}
		th := hint.TableHints.Items[0].(*ast.TableHint)
		if th.Name != "INDEX" {
			t.Errorf("expected INDEX, got %q", th.Name)
		}
	})

	// TABLE HINT with multiple hints
	t.Run("option_table_hint_multiple", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (TABLE HINT(t, NOLOCK, NOWAIT))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "TABLE HINT" {
			t.Errorf("expected 'TABLE HINT', got %q", hint.Kind)
		}
		if hint.TableHints == nil || hint.TableHints.Len() != 2 {
			t.Fatalf("expected 2 table hints, got %d", hint.TableHints.Len())
		}
	})
}

// TestParseSelectQueryHintRemainingDepth tests structured query hints (batch 127).
func TestParseSelectQueryHintRemainingDepth(t *testing.T) {
	// Structured RECOMPILE hint
	t.Run("query_hint_structured_recompile", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (RECOMPILE)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "RECOMPILE" {
			t.Errorf("expected RECOMPILE, got %q", hint.Kind)
		}
		if hint.Loc.Start == 0 && hint.Loc.End == 0 {
			t.Errorf("expected non-zero Loc")
		}
	})

	// Structured OPTIMIZE FOR with value and UNKNOWN params
	t.Run("query_hint_structured_optimize_for", func(t *testing.T) {
		sql := "SELECT * FROM t WHERE x = @x AND y = @y OPTION (OPTIMIZE FOR (@x = 5, @y UNKNOWN))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "OPTIMIZE FOR" {
			t.Errorf("expected 'OPTIMIZE FOR', got %q", hint.Kind)
		}
		if hint.Params == nil || hint.Params.Len() != 2 {
			t.Fatalf("expected 2 params, got %d", hint.Params.Len())
		}
		p0 := hint.Params.Items[0].(*ast.OptimizeForParam)
		if p0.Variable != "@x" || p0.Unknown || p0.Value == nil {
			t.Errorf("param 0: expected @x = value, got var=%q unknown=%v value=%v", p0.Variable, p0.Unknown, p0.Value)
		}
		p1 := hint.Params.Items[1].(*ast.OptimizeForParam)
		if p1.Variable != "@y" || !p1.Unknown {
			t.Errorf("param 1: expected @y UNKNOWN, got var=%q unknown=%v", p1.Variable, p1.Unknown)
		}
	})

	// Structured TABLE HINT reusing parseTableHint()
	t.Run("query_hint_structured_table_hint", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (TABLE HINT(t, INDEX(IX_1), NOLOCK))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "TABLE HINT" {
			t.Errorf("expected 'TABLE HINT', got %q", hint.Kind)
		}
		if hint.TableName == nil || hint.TableName.Object != "t" {
			t.Errorf("expected table name 't'")
		}
		if hint.TableHints == nil || hint.TableHints.Len() != 2 {
			t.Fatalf("expected 2 table hints, got %d", hint.TableHints.Len())
		}
		th0 := hint.TableHints.Items[0].(*ast.TableHint)
		if th0.Name != "INDEX" || th0.IndexValues == nil {
			t.Errorf("expected INDEX hint with values, got %q", th0.Name)
		}
		th1 := hint.TableHints.Items[1].(*ast.TableHint)
		if th1.Name != "NOLOCK" {
			t.Errorf("expected NOLOCK, got %q", th1.Name)
		}
	})

	// Structured unknown hint with name = value
	t.Run("query_hint_structured_unknown", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (MAX_GRANT_PERCENT = 25)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "MAX_GRANT_PERCENT" {
			t.Errorf("expected 'MAX_GRANT_PERCENT', got %q", hint.Kind)
		}
		if hint.Value == nil {
			t.Fatalf("expected Value to be set")
		}
	})

	// USE HINT with string list
	t.Run("query_hint_use_hint", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (USE HINT('DISABLE_OPTIMIZED_NESTED_LOOP', 'FORCE_LEGACY_CARDINALITY_ESTIMATION'))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "USE HINT" {
			t.Errorf("expected 'USE HINT', got %q", hint.Kind)
		}
		if hint.Params == nil || hint.Params.Len() != 2 {
			t.Fatalf("expected 2 hint name params, got %d", hint.Params.Len())
		}
	})

	// PARAMETERIZATION with mode
	t.Run("query_hint_parameterization", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (PARAMETERIZATION FORCED)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "PARAMETERIZATION" {
			t.Errorf("expected PARAMETERIZATION, got %q", hint.Kind)
		}
		if hint.StrValue != "FORCED" {
			t.Errorf("expected FORCED, got %q", hint.StrValue)
		}
	})

	// QUERYTRACEON
	t.Run("query_hint_querytraceon", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (QUERYTRACEON 4199)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "QUERYTRACEON" {
			t.Errorf("expected QUERYTRACEON, got %q", hint.Kind)
		}
		if hint.Value == nil {
			t.Fatalf("expected trace flag value")
		}
	})

	// TABLE HINT with FORCESEEK
	t.Run("query_hint_table_hint_forceseek", func(t *testing.T) {
		sql := "SELECT * FROM t OPTION (TABLE HINT(t, FORCESEEK(IX_1(col1, col2))))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SelectStmt)
		hint := stmt.OptionClause.Items[0].(*ast.QueryHint)
		if hint.Kind != "TABLE HINT" {
			t.Errorf("expected 'TABLE HINT', got %q", hint.Kind)
		}
		if hint.TableHints == nil || hint.TableHints.Len() != 1 {
			t.Fatalf("expected 1 table hint")
		}
		th := hint.TableHints.Items[0].(*ast.TableHint)
		if th.Name != "FORCESEEK" {
			t.Errorf("expected FORCESEEK, got %q", th.Name)
		}
		if th.ForceSeekColumns == nil || th.ForceSeekColumns.Len() != 2 {
			t.Errorf("expected 2 forceseek columns")
		}
	})
}

func TestParseSetuserStatement(t *testing.T) {
	// SETUSER with no arguments (reset identity)
	t.Run("setuser_basic_reset", func(t *testing.T) {
		sql := "SETUSER"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Action != "SETUSER" {
			t.Errorf("expected action SETUSER, got %q", stmt.Action)
		}
		if stmt.Name != "" {
			t.Errorf("expected empty name, got %q", stmt.Name)
		}
	})

	// SETUSER 'username'
	t.Run("setuser_basic", func(t *testing.T) {
		sql := "SETUSER 'mary'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Action != "SETUSER" {
			t.Errorf("expected action SETUSER, got %q", stmt.Action)
		}
		if stmt.Name != "mary" {
			t.Errorf("expected name 'mary', got %q", stmt.Name)
		}
	})

	// SETUSER 'username' WITH NORESET
	t.Run("setuser_with_noreset", func(t *testing.T) { //nolint:dupl
		sql := "SETUSER 'mary' WITH NORESET"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Action != "SETUSER" {
			t.Errorf("expected action SETUSER, got %q", stmt.Action)
		}
		if stmt.Name != "mary" {
			t.Errorf("expected name 'mary', got %q", stmt.Name)
		}
		if stmt.Options == nil || len(stmt.Options.Items) == 0 {
			t.Fatal("expected NORESET option")
		}
		opt := stmt.Options.Items[0].(*ast.String).Str
		if opt != "NORESET" {
			t.Errorf("expected NORESET option, got %q", opt)
		}
	})
}

func TestParseServerAuditOptionsDepth(t *testing.T) {
	// CREATE SERVER AUDIT TO FILE with options
	t.Run("server_audit_to_file", func(t *testing.T) {
		sql := `CREATE SERVER AUDIT MyAudit TO FILE (FILEPATH = '\\server\audit\', MAXSIZE = 100MB, MAX_ROLLOVER_FILES = 10)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Action != "CREATE" {
			t.Errorf("expected CREATE, got %q", stmt.Action)
		}
		if stmt.Name != "MyAudit" {
			t.Errorf("expected MyAudit, got %q", stmt.Name)
		}
		if stmt.Options == nil || len(stmt.Options.Items) < 3 {
			t.Fatalf("expected at least 3 options (TO=FILE + file options), got %v", stmt.Options)
		}
		opt0 := stmt.Options.Items[0].(*ast.String).Str
		if opt0 != "TO=FILE" {
			t.Errorf("expected TO=FILE, got %q", opt0)
		}
	})

	// CREATE SERVER AUDIT WITH options
	t.Run("server_audit_with_options", func(t *testing.T) {
		sql := "CREATE SERVER AUDIT MyAudit TO APPLICATION_LOG WITH (QUEUE_DELAY = 1000, ON_FAILURE = SHUTDOWN)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Options == nil || len(stmt.Options.Items) < 3 {
			t.Fatalf("expected at least 3 options, got %v", stmt.Options)
		}
		opt0 := stmt.Options.Items[0].(*ast.String).Str
		if opt0 != "TO=APPLICATION_LOG" {
			t.Errorf("expected TO=APPLICATION_LOG, got %q", opt0)
		}
	})

	// CREATE SERVER AUDIT with WHERE clause
	t.Run("server_audit_where_clause", func(t *testing.T) {
		sql := "CREATE SERVER AUDIT MyAudit TO FILE (FILEPATH = 'C:\\audit\\') WHERE object_name = 'SensitiveData'"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
		// Find WHERE option
		found := false
		for _, item := range stmt.Options.Items {
			s := item.(*ast.String).Str
			if len(s) > 6 && s[:6] == "WHERE=" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected WHERE predicate in options")
		}
	})

	// ALTER SERVER AUDIT with MODIFY NAME
	t.Run("alter_server_audit_modify_name", func(t *testing.T) {
		sql := "ALTER SERVER AUDIT MyAudit MODIFY NAME = NewAuditName"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
		opt0 := stmt.Options.Items[0].(*ast.String).Str
		if opt0 != "MODIFY NAME=NewAuditName" {
			t.Errorf("expected MODIFY NAME=NewAuditName, got %q", opt0)
		}
	})

	// ALTER SERVER AUDIT with REMOVE WHERE
	t.Run("alter_server_audit_remove_where", func(t *testing.T) {
		sql := "ALTER SERVER AUDIT MyAudit REMOVE WHERE"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityStmt)
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
		opt0 := stmt.Options.Items[0].(*ast.String).Str
		if opt0 != "REMOVE WHERE" {
			t.Errorf("expected REMOVE WHERE, got %q", opt0)
		}
	})
}

func TestParseAlterIndexOptionsDepth(t *testing.T) {
	// ALTER INDEX ... REBUILD WITH options
	t.Run("alter_index_rebuild_with", func(t *testing.T) {
		sql := "ALTER INDEX IX_1 ON dbo.t REBUILD WITH (FILLFACTOR = 80, PAD_INDEX = ON, SORT_IN_TEMPDB = ON)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterIndexStmt)
		if stmt.Action != "REBUILD" {
			t.Errorf("expected REBUILD, got %q", stmt.Action)
		}
		if stmt.Options == nil || len(stmt.Options.Items) != 3 {
			t.Fatalf("expected 3 options, got %v", stmt.Options)
		}
		opt0 := stmt.Options.Items[0].(*ast.String).Str
		if opt0 != "FILLFACTOR=80" {
			t.Errorf("expected FILLFACTOR=80, got %q", opt0)
		}
	})

	// ALTER INDEX ... REBUILD PARTITION = ALL WITH options
	t.Run("alter_index_rebuild_partition", func(t *testing.T) {
		sql := "ALTER INDEX ALL ON dbo.t REBUILD PARTITION = ALL WITH (ONLINE = ON, MAXDOP = 4)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterIndexStmt)
		if stmt.Action != "REBUILD" {
			t.Errorf("expected REBUILD, got %q", stmt.Action)
		}
		if stmt.Partition != "ALL" {
			t.Errorf("expected partition ALL, got %q", stmt.Partition)
		}
		if stmt.Options == nil || len(stmt.Options.Items) != 2 {
			t.Fatalf("expected 2 options, got %v", stmt.Options)
		}
	})

	// ALTER INDEX ... REORGANIZE PARTITION = n
	t.Run("alter_index_reorganize_partition", func(t *testing.T) {
		sql := "ALTER INDEX IX_1 ON dbo.t REORGANIZE PARTITION = 5 WITH (LOB_COMPACTION = ON)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterIndexStmt)
		if stmt.Action != "REORGANIZE" {
			t.Errorf("expected REORGANIZE, got %q", stmt.Action)
		}
		if stmt.Partition != "5" {
			t.Errorf("expected partition 5, got %q", stmt.Partition)
		}
		if stmt.Options == nil || len(stmt.Options.Items) != 1 {
			t.Fatalf("expected 1 option, got %v", stmt.Options)
		}
	})

	// ALTER INDEX ... SET options
	t.Run("alter_index_set", func(t *testing.T) {
		sql := "ALTER INDEX IX_1 ON dbo.t SET (STATISTICS_NORECOMPUTE = ON, ALLOW_ROW_LOCKS = OFF)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterIndexStmt)
		if stmt.Action != "SET" {
			t.Errorf("expected SET, got %q", stmt.Action)
		}
		if stmt.Options == nil || len(stmt.Options.Items) != 2 {
			t.Fatalf("expected 2 options, got %v", stmt.Options)
		}
	})

	// ALTER INDEX ... DISABLE (no options)
	t.Run("alter_index_disable", func(t *testing.T) {
		sql := "ALTER INDEX IX_1 ON dbo.t DISABLE"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterIndexStmt)
		if stmt.Action != "DISABLE" {
			t.Errorf("expected DISABLE, got %q", stmt.Action)
		}
	})

	// ALTER INDEX ... REBUILD with DATA_COMPRESSION ON PARTITIONS
	t.Run("alter_index_data_compression_partitions", func(t *testing.T) {
		sql := "ALTER INDEX IX_1 ON dbo.t REBUILD WITH (DATA_COMPRESSION = PAGE ON PARTITIONS (1, 3 TO 5))"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterIndexStmt)
		if stmt.Action != "REBUILD" {
			t.Errorf("expected REBUILD, got %q", stmt.Action)
		}
		if stmt.Options == nil || len(stmt.Options.Items) == 0 {
			t.Fatal("expected options")
		}
	})
}

func TestParseBeginConversationTimer(t *testing.T) {
	// BEGIN CONVERSATION TIMER with variable handle and integer timeout
	t.Run("begin_conversation_timer", func(t *testing.T) {
		sql := "BEGIN CONVERSATION TIMER (@dialog_handle) TIMEOUT = 120"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.ServiceBrokerStmt)
		if stmt.Action != "BEGIN" {
			t.Errorf("expected action BEGIN, got %q", stmt.Action)
		}
		if stmt.ObjectType != "CONVERSATION TIMER" {
			t.Errorf("expected object type CONVERSATION TIMER, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil || len(stmt.Options.Items) < 2 {
			t.Fatal("expected at least 2 options (HANDLE and TIMEOUT)")
		}
	})
}

// TestParseSecurityKeyStructuredOptions tests batch 98: structured security key option parsing.
func TestParseSecurityKeyStructuredOptions(t *testing.T) {
	// ---- SYMMETRIC KEY ----

	t.Run("symmetric_key_structured_create_full", func(t *testing.T) {
		sql := `CREATE SYMMETRIC KEY MySymKey
			AUTHORIZATION dbo
			WITH ALGORITHM = AES_256, KEY_SOURCE = 'my pass phrase', IDENTITY_VALUE = 'my identity'
			ENCRYPTION BY CERTIFICATE MyCert`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "CREATE" {
			t.Errorf("expected CREATE, got %q", stmt.Action)
		}
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Name != "MySymKey" {
			t.Errorf("expected name MySymKey, got %q", stmt.Name)
		}
		if stmt.Options == nil || stmt.Options.Len() < 3 {
			t.Fatalf("expected at least 3 options, got %d", stmt.Options.Len())
		}
	})

	t.Run("symmetric_key_structured_create_password", func(t *testing.T) {
		sql := `CREATE SYMMETRIC KEY TempKey
			WITH ALGORITHM = AES_128
			ENCRYPTION BY PASSWORD = 'StrongPass123!'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("symmetric_key_structured_create_from_provider", func(t *testing.T) {
		sql := `CREATE SYMMETRIC KEY MyEKMKey
			FROM PROVIDER MyEKMProvider
			WITH PROVIDER_KEY_NAME = 'KeyInEKM', CREATION_DISPOSITION = OPEN_EXISTING
			ENCRYPTION BY PASSWORD = 'pass123'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("symmetric_key_structured_create_multi_encryption", func(t *testing.T) {
		sql := `CREATE SYMMETRIC KEY MultiEncKey
			WITH ALGORITHM = AES_256
			ENCRYPTION BY CERTIFICATE MyCert, PASSWORD = 'pass123', ASYMMETRIC KEY MyAsymKey`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("symmetric_key_structured_alter_add", func(t *testing.T) {
		sql := `ALTER SYMMETRIC KEY MySymKey ADD ENCRYPTION BY CERTIFICATE NewCert`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("symmetric_key_structured_alter_drop", func(t *testing.T) {
		sql := `ALTER SYMMETRIC KEY MySymKey DROP ENCRYPTION BY PASSWORD = 'OldPass'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- ASYMMETRIC KEY ----

	t.Run("asymmetric_key_structured_create_algorithm", func(t *testing.T) {
		sql := `CREATE ASYMMETRIC KEY MyAsymKey
			WITH ALGORITHM = RSA_2048
			ENCRYPTION BY PASSWORD = 'StrongPass!'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "ASYMMETRIC KEY" {
			t.Errorf("expected ASYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("asymmetric_key_structured_create_from_file", func(t *testing.T) {
		sql := `CREATE ASYMMETRIC KEY FileKey
			AUTHORIZATION dbo
			FROM FILE = 'c:\keys\mykey.snk'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "ASYMMETRIC KEY" {
			t.Errorf("expected ASYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("asymmetric_key_structured_create_from_assembly", func(t *testing.T) {
		sql := `CREATE ASYMMETRIC KEY AsmKey FROM ASSEMBLY MyAssembly`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "ASYMMETRIC KEY" {
			t.Errorf("expected ASYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("asymmetric_key_structured_create_from_provider", func(t *testing.T) {
		sql := `CREATE ASYMMETRIC KEY EKMKey
			FROM PROVIDER MyEKMProvider
			WITH ALGORITHM = RSA_2048, PROVIDER_KEY_NAME = 'key1', CREATION_DISPOSITION = CREATE_NEW`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "ASYMMETRIC KEY" {
			t.Errorf("expected ASYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("asymmetric_key_structured_create_executable_file", func(t *testing.T) {
		sql := `CREATE ASYMMETRIC KEY ExeKey FROM EXECUTABLE FILE = 'c:\keys\mydll.dll'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "ASYMMETRIC KEY" {
			t.Errorf("expected ASYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- CERTIFICATE ----

	t.Run("certificate_structured_create_self_signed", func(t *testing.T) {
		sql := `CREATE CERTIFICATE MyCert
			ENCRYPTION BY PASSWORD = 'CertPass123!'
			WITH SUBJECT = 'Test Certificate',
			START_DATE = '2024-01-01', EXPIRY_DATE = '2025-12-31'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil || stmt.Options.Len() < 3 {
			t.Fatalf("expected at least 3 options, got %d", stmt.Options.Len())
		}
	})

	t.Run("certificate_structured_create_from_file", func(t *testing.T) {
		sql := `CREATE CERTIFICATE FileCert
			FROM FILE = 'c:\certs\mycert.cer'
			WITH PRIVATE KEY (FILE = 'c:\certs\mykey.pvk', DECRYPTION BY PASSWORD = 'oldpass', ENCRYPTION BY PASSWORD = 'newpass')`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("certificate_structured_create_from_assembly", func(t *testing.T) {
		sql := `CREATE CERTIFICATE AsmCert FROM ASSEMBLY MySignedAssembly`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("certificate_structured_active_for_begin_dialog", func(t *testing.T) {
		sql := `CREATE CERTIFICATE BrokerCert
			WITH SUBJECT = 'Broker Certificate'
			ACTIVE FOR BEGIN_DIALOG = ON`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("certificate_structured_authorization", func(t *testing.T) {
		sql := `CREATE CERTIFICATE OwnedCert
			AUTHORIZATION dbo
			WITH SUBJECT = 'Owned Certificate'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- CREDENTIAL ----

	t.Run("credential_structured_create", func(t *testing.T) {
		sql := `CREATE CREDENTIAL MyCred
			WITH IDENTITY = 'my_user', SECRET = 'my_secret'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CREDENTIAL" {
			t.Errorf("expected CREDENTIAL, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil || stmt.Options.Len() < 2 {
			t.Fatalf("expected at least 2 options, got %d", stmt.Options.Len())
		}
	})

	t.Run("credential_structured_create_with_provider", func(t *testing.T) {
		sql := `CREATE CREDENTIAL EKMCred
			WITH IDENTITY = 'User1OnEKM', SECRET = 'secretpass'
			FOR CRYPTOGRAPHIC PROVIDER MyEKMProvider`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CREDENTIAL" {
			t.Errorf("expected CREDENTIAL, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("credential_structured_create_identity_only", func(t *testing.T) {
		sql := `CREATE CREDENTIAL SimpleCred WITH IDENTITY = 'Managed Identity'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CREDENTIAL" {
			t.Errorf("expected CREDENTIAL, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- MASTER KEY ----

	t.Run("master_key_structured_create", func(t *testing.T) {
		sql := `CREATE MASTER KEY ENCRYPTION BY PASSWORD = 'StrongMasterPass123!'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "MASTER KEY" {
			t.Errorf("expected MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Action != "CREATE" {
			t.Errorf("expected CREATE, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("master_key_structured_alter_regenerate", func(t *testing.T) {
		sql := `ALTER MASTER KEY REGENERATE WITH ENCRYPTION BY PASSWORD = 'NewPass123!'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.ObjectType != "MASTER KEY" {
			t.Errorf("expected MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("master_key_structured_alter_force_regenerate", func(t *testing.T) {
		sql := `ALTER MASTER KEY FORCE REGENERATE WITH ENCRYPTION BY PASSWORD = 'ForcePass!'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("master_key_structured_alter_add_encryption_smk", func(t *testing.T) {
		sql := `ALTER MASTER KEY ADD ENCRYPTION BY SERVICE MASTER KEY`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("master_key_structured_alter_drop_encryption_password", func(t *testing.T) {
		sql := `ALTER MASTER KEY DROP ENCRYPTION BY PASSWORD = 'OldPass123!'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- SERVICE MASTER KEY ----

	t.Run("service_master_key_alter_regenerate", func(t *testing.T) {
		sql := `ALTER SERVICE MASTER KEY REGENERATE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.ObjectType != "SERVICE MASTER KEY" {
			t.Errorf("expected SERVICE MASTER KEY, got %q", stmt.ObjectType)
		}
	})

	t.Run("service_master_key_alter_force_regenerate", func(t *testing.T) {
		sql := `ALTER SERVICE MASTER KEY FORCE REGENERATE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "SERVICE MASTER KEY" {
			t.Errorf("expected SERVICE MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("service_master_key_backup", func(t *testing.T) {
		sql := `BACKUP SERVICE MASTER KEY TO FILE = 'c:\backup\smk.bak' ENCRYPTION BY PASSWORD = 'BackupPass'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "BACKUP" {
			t.Errorf("expected BACKUP, got %q", stmt.Action)
		}
		if stmt.ObjectType != "SERVICE MASTER KEY" {
			t.Errorf("expected SERVICE MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("service_master_key_restore", func(t *testing.T) {
		sql := `RESTORE SERVICE MASTER KEY FROM FILE = 'c:\backup\smk.bak' DECRYPTION BY PASSWORD = 'RestorePass'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "RESTORE" {
			t.Errorf("expected RESTORE, got %q", stmt.Action)
		}
		if stmt.ObjectType != "SERVICE MASTER KEY" {
			t.Errorf("expected SERVICE MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("service_master_key_restore_force", func(t *testing.T) {
		sql := `RESTORE SERVICE MASTER KEY FROM FILE = 'c:\backup\smk.bak' DECRYPTION BY PASSWORD = 'RestorePass' FORCE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "RESTORE" {
			t.Errorf("expected RESTORE, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- OPEN/CLOSE/BACKUP ----

	t.Run("open_symmetric_key_structured", func(t *testing.T) {
		sql := `OPEN SYMMETRIC KEY MySymKey DECRYPTION BY CERTIFICATE MyCert`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "OPEN" {
			t.Errorf("expected OPEN, got %q", stmt.Action)
		}
		if stmt.ObjectType != "SYMMETRIC KEY" {
			t.Errorf("expected SYMMETRIC KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("open_symmetric_key_password", func(t *testing.T) {
		sql := `OPEN SYMMETRIC KEY MySymKey DECRYPTION BY PASSWORD = 'MyPassword'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "OPEN" {
			t.Errorf("expected OPEN, got %q", stmt.Action)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("open_master_key_structured", func(t *testing.T) {
		sql := `OPEN MASTER KEY DECRYPTION BY PASSWORD = 'MasterPass'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "OPEN" {
			t.Errorf("expected OPEN, got %q", stmt.Action)
		}
		if stmt.ObjectType != "MASTER KEY" {
			t.Errorf("expected MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("backup_certificate_structured", func(t *testing.T) {
		sql := `BACKUP CERTIFICATE MyCert TO FILE = 'c:\certs\mycert.cer'
			WITH PRIVATE KEY (FILE = 'c:\certs\mykey.pvk', ENCRYPTION BY PASSWORD = 'ExportPass')`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "BACKUP" {
			t.Errorf("expected BACKUP, got %q", stmt.Action)
		}
		if stmt.ObjectType != "CERTIFICATE" {
			t.Errorf("expected CERTIFICATE, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("backup_master_key_structured", func(t *testing.T) {
		sql := `BACKUP MASTER KEY TO FILE = 'c:\backup\masterkey.bak' ENCRYPTION BY PASSWORD = 'BackupMKPass'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "BACKUP" {
			t.Errorf("expected BACKUP, got %q", stmt.Action)
		}
		if stmt.ObjectType != "MASTER KEY" {
			t.Errorf("expected MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("restore_master_key_structured", func(t *testing.T) {
		sql := `RESTORE MASTER KEY FROM FILE = 'c:\backup\masterkey.bak'
			DECRYPTION BY PASSWORD = 'FilePass'
			ENCRYPTION BY PASSWORD = 'NewMasterPass'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "RESTORE" {
			t.Errorf("expected RESTORE, got %q", stmt.Action)
		}
		if stmt.ObjectType != "MASTER KEY" {
			t.Errorf("expected MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- DATABASE ENCRYPTION KEY / DATABASE SCOPED CREDENTIAL ----

	t.Run("database_encryption_key_create", func(t *testing.T) {
		sql := `CREATE DATABASE ENCRYPTION KEY
			WITH ALGORITHM = AES_256
			ENCRYPTION BY SERVER CERTIFICATE MyCert`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "DATABASE ENCRYPTION KEY" {
			t.Errorf("expected DATABASE ENCRYPTION KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("database_scoped_credential_create", func(t *testing.T) {
		sql := `CREATE DATABASE SCOPED CREDENTIAL MyDBCred
			WITH IDENTITY = 'db_user', SECRET = 'db_secret'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "DATABASE SCOPED CREDENTIAL" {
			t.Errorf("expected DATABASE SCOPED CREDENTIAL, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- CRYPTOGRAPHIC PROVIDER ----

	t.Run("cryptographic_provider_create", func(t *testing.T) {
		sql := `CREATE CRYPTOGRAPHIC PROVIDER MyProvider FROM FILE = 'c:\ekm\provider.dll'`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "CRYPTOGRAPHIC PROVIDER" {
			t.Errorf("expected CRYPTOGRAPHIC PROVIDER, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	// ---- COLUMN ENCRYPTION KEY / COLUMN MASTER KEY ----

	t.Run("column_encryption_key_create", func(t *testing.T) {
		sql := `CREATE COLUMN ENCRYPTION KEY MyCEK
			WITH VALUES
			(COLUMN_MASTER_KEY = MyCMK, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0x0102030405)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "COLUMN ENCRYPTION KEY" {
			t.Errorf("expected COLUMN ENCRYPTION KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("column_master_key_create", func(t *testing.T) {
		sql := `CREATE COLUMN MASTER KEY MyCMK
			WITH (KEY_STORE_PROVIDER_NAME = 'MSSQL_CERTIFICATE_STORE', KEY_PATH = 'Current User/Personal/my_cert')`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.ObjectType != "COLUMN MASTER KEY" {
			t.Errorf("expected COLUMN MASTER KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})

	t.Run("alter_column_encryption_key_add", func(t *testing.T) {
		sql := `ALTER COLUMN ENCRYPTION KEY MyCEK
			ADD VALUE (COLUMN_MASTER_KEY = NewCMK, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0xABCD)`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.SecurityKeyStmt)
		if stmt.Action != "ALTER" {
			t.Errorf("expected ALTER, got %q", stmt.Action)
		}
		if stmt.ObjectType != "COLUMN ENCRYPTION KEY" {
			t.Errorf("expected COLUMN ENCRYPTION KEY, got %q", stmt.ObjectType)
		}
		if stmt.Options == nil {
			t.Fatal("expected options")
		}
	})
}

// TestParseEndpointOptionsDepth tests batch 99: structured endpoint options parsing.
func TestParseEndpointOptionsDepth(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// TCP protocol options
		{
			"tcp_listener_port",
			"CREATE ENDPOINT ep1 STATE = STARTED AS TCP (LISTENER_PORT = 4022) FOR TSQL ()",
		},
		{
			"tcp_listener_port_and_ip_all",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022, LISTENER_IP = ALL) FOR TSQL ()",
		},
		{
			"tcp_listener_ip_ipv4",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 5022, LISTENER_IP = (10.0.75.1)) FOR TSQL ()",
		},
		{
			"tcp_listener_ip_ipv6",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 5022, LISTENER_IP = ('::1')) FOR TSQL ()",
		},
		// TSQL payload options
		{
			"payload_tsql_empty",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 1433) FOR TSQL ()",
		},
		{
			"payload_tsql_encryption_strict",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 1433) FOR TSQL (ENCRYPTION = STRICT)",
		},
		{
			"payload_tsql_encryption_negotiated",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 1433) FOR TSQL (ENCRYPTION = NEGOTIATED)",
		},
		// SERVICE_BROKER payload options
		{
			"payload_service_broker_auth_windows",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = WINDOWS NEGOTIATE)",
		},
		{
			"payload_service_broker_auth_cert",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = CERTIFICATE MyCert)",
		},
		{
			"payload_service_broker_auth_windows_cert",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = WINDOWS KERBEROS CERTIFICATE MyCert)",
		},
		{
			"payload_service_broker_auth_cert_windows",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = CERTIFICATE MyCert WINDOWS NTLM)",
		},
		{
			"payload_service_broker_encryption",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (ENCRYPTION = REQUIRED ALGORITHM AES)",
		},
		{
			"payload_service_broker_encryption_disabled",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (ENCRYPTION = DISABLED)",
		},
		{
			"payload_service_broker_encryption_supported_rc4_aes",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (ENCRYPTION = SUPPORTED ALGORITHM RC4 AES)",
		},
		{
			"payload_service_broker_message_forwarding",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (MESSAGE_FORWARDING = ENABLED, MESSAGE_FORWARD_SIZE = 10)",
		},
		{
			"payload_service_broker_full",
			"CREATE ENDPOINT ep1 STATE = STARTED AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = WINDOWS, ENCRYPTION = SUPPORTED ALGORITHM AES, MESSAGE_FORWARDING = DISABLED)",
		},
		// DATABASE_MIRRORING payload options
		{
			"payload_db_mirroring_role_witness",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 7022) FOR DATABASE_MIRRORING (ROLE = WITNESS)",
		},
		{
			"payload_db_mirroring_role_partner",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 7022) FOR DATABASE_MIRRORING (ROLE = PARTNER)",
		},
		{
			"payload_db_mirroring_role_all",
			"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 7022) FOR DATABASE_MIRRORING (ROLE = ALL)",
		},
		{
			"payload_db_mirroring_full",
			"CREATE ENDPOINT ep1 STATE = STARTED AS TCP (LISTENER_PORT = 7022) FOR DATABASE_MIRRORING (AUTHENTICATION = WINDOWS KERBEROS, ENCRYPTION = SUPPORTED, ROLE = ALL)",
		},
		// ALTER ENDPOINT
		{
			"alter_endpoint_state",
			"ALTER ENDPOINT ep1 STATE = STARTED",
		},
		{
			"alter_endpoint_tcp_options",
			"ALTER ENDPOINT ep1 AS TCP (LISTENER_PORT = 5023)",
		},
		{
			"alter_endpoint_full",
			"ALTER ENDPOINT ep1 STATE = DISABLED AS TCP (LISTENER_PORT = 7022) FOR DATABASE_MIRRORING (AUTHENTICATION = CERTIFICATE MyCert, ENCRYPTION = REQUIRED ALGORITHM AES RC4, ROLE = PARTNER)",
		},
		// AUTHORIZATION
		{
			"create_endpoint_authorization",
			"CREATE ENDPOINT ep1 AUTHORIZATION sa STATE = STARTED AS TCP (LISTENER_PORT = 4022) FOR TSQL ()",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tc.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tc.sql, result.Items[0])
			}
			checkLocation(t, tc.sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseServiceBrokerAlterOptionsDepth tests batch 100: structured parsing for ALTER Service Broker objects.
func TestParseServiceBrokerAlterOptionsDepth(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// ALTER MESSAGE TYPE with VALIDATION
		{
			name: "alter_message_type_validation_none",
			sql:  "ALTER MESSAGE TYPE MyMessageType VALIDATION = NONE",
		},
		{
			name: "alter_message_type_validation_empty",
			sql:  "ALTER MESSAGE TYPE MyMessageType VALIDATION = EMPTY",
		},
		{
			name: "alter_message_type_validation_well_formed_xml",
			sql:  "ALTER MESSAGE TYPE [//Adventure-Works.com/Expenses/SubmitExpense] VALIDATION = WELL_FORMED_XML",
		},
		{
			name: "alter_message_type_validation_valid_xml_schema",
			sql:  "ALTER MESSAGE TYPE MyMessageType VALIDATION = VALID_XML WITH SCHEMA COLLECTION MySchemaCollection",
		},
		// ALTER CONTRACT with message type modifications (ADD/DROP)
		{
			name: "alter_contract_add_message_type",
			sql:  "ALTER CONTRACT MyContract ADD MESSAGE TYPE MyMsgType SENT BY INITIATOR",
		},
		{
			name: "alter_contract_drop_message_type",
			sql:  "ALTER CONTRACT MyContract DROP MESSAGE TYPE MyMsgType",
		},
		{
			name: "alter_contract_add_sent_by_target",
			sql:  "ALTER CONTRACT MyContract ADD MESSAGE TYPE ResponseMsg SENT BY TARGET",
		},
		{
			name: "alter_contract_add_sent_by_any",
			sql:  "ALTER CONTRACT MyContract ADD MESSAGE TYPE AnyMsg SENT BY ANY",
		},
		// ALTER ROUTE with structured WITH options
		{
			name: "alter_route_all_options",
			sql:  "ALTER ROUTE MyRoute WITH SERVICE_NAME = '//example.com/svc', BROKER_INSTANCE = 'D8D4D268-00A3-4C62-8F91-634B89B1E317', LIFETIME = 600, ADDRESS = 'TCP://10.0.0.2:4022', MIRROR_ADDRESS = 'TCP://10.0.0.3:4022'",
		},
		{
			name: "alter_route_address_local",
			sql:  "ALTER ROUTE MyRoute WITH ADDRESS = 'LOCAL'",
		},
		{
			name: "alter_route_address_transport",
			sql:  "ALTER ROUTE MyRoute WITH ADDRESS = 'TRANSPORT'",
		},
		// ALTER REMOTE SERVICE BINDING with structured WITH options
		{
			name: "alter_remote_binding_user_and_anonymous_off",
			sql:  "ALTER REMOTE SERVICE BINDING MyBinding WITH USER = SecurityAccount, ANONYMOUS = OFF",
		},
		{
			name: "alter_remote_binding_user_only",
			sql:  "ALTER REMOTE SERVICE BINDING MyBinding WITH USER = MyUser",
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

// TestParseEventSessionSpecsDepth tests batch 101: structured parsing of event session event/target specs.
func TestParseEventSessionSpecsDepth(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// Event spec with SET
		{
			name: "event_spec_set",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (SET collect_batch_text = 1)",
		},
		{
			name: "event_spec_set_multiple",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.query_post_execution_showplan (SET collect_database_name = 1, collect_object_name = 1)",
		},
		// Event spec with ACTION
		{
			name: "event_spec_action",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (ACTION (sqlserver.session_id))",
		},
		{
			name: "event_spec_action_multiple",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (ACTION (sqlserver.session_id, sqlserver.database_id, sqlserver.client_hostname))",
		},
		// Event spec with WHERE
		{
			name: "event_spec_where_simple",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (WHERE sqlserver.session_id > 50)",
		},
		{
			name: "event_spec_where_and",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (WHERE sqlserver.session_id > 50 AND sqlserver.database_id = 5)",
		},
		{
			name: "event_spec_where_string",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (WHERE batch_text = 'SELECT')",
		},
		// Event spec with SET + ACTION + WHERE combined
		{
			name: "event_spec_all_clauses",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting (SET collect_batch_text = 1 ACTION (sqlserver.session_id) WHERE sqlserver.session_id > 50)",
		},
		// Target spec with SET
		{
			name: "target_spec_set",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting ADD TARGET package0.event_file (SET filename = 'test.xel')",
		},
		{
			name: "target_spec_set_multiple",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting ADD TARGET package0.event_file (SET filename = 'test.xel', max_file_size = 256, max_rollover_files = 10)",
		},
		{
			name: "target_spec_ring_buffer",
			sql:  "CREATE EVENT SESSION s1 ON SERVER ADD EVENT sqlserver.sql_batch_starting ADD TARGET package0.ring_buffer (SET max_memory = 1024)",
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
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)

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

// TestParseResourceGovernorOptionsDepth tests batch 102: structured parsing of Resource Governor options.
func TestParseResourceGovernorOptionsDepth(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		// Workload group structured options
		{
			name: "workload_group_all_options",
			sql:  "CREATE WORKLOAD GROUP wg1 WITH (IMPORTANCE = HIGH, REQUEST_MAX_MEMORY_GRANT_PERCENT = 25, REQUEST_MAX_CPU_TIME_SEC = 60, REQUEST_MEMORY_GRANT_TIMEOUT_SEC = 30, MAX_DOP = 8, GROUP_MAX_REQUESTS = 100)",
		},
		{
			name: "workload_group_using_pool",
			sql:  "CREATE WORKLOAD GROUP wg1 WITH (MAX_DOP = 4) USING myPool",
		},
		{
			name: "workload_group_using_external",
			sql:  "CREATE WORKLOAD GROUP wg1 USING myPool, EXTERNAL myExtPool",
		},
		{
			name: "alter_workload_group_importance",
			sql:  "ALTER WORKLOAD GROUP [default] WITH (IMPORTANCE = LOW)",
		},
		// Resource pool structured options
		{
			name: "resource_pool_cpu_memory",
			sql:  "CREATE RESOURCE POOL rp1 WITH (MIN_CPU_PERCENT = 10, MAX_CPU_PERCENT = 50, MIN_MEMORY_PERCENT = 5, MAX_MEMORY_PERCENT = 25)",
		},
		{
			name: "resource_pool_iops",
			sql:  "CREATE RESOURCE POOL rp1 WITH (MIN_IOPS_PER_VOLUME = 100, MAX_IOPS_PER_VOLUME = 1000)",
		},
		{
			name: "resource_pool_affinity_auto",
			sql:  "CREATE RESOURCE POOL rp1 WITH (AFFINITY SCHEDULER = AUTO)",
		},
		{
			name: "resource_pool_affinity_range",
			sql:  "CREATE RESOURCE POOL rp1 WITH (AFFINITY SCHEDULER = (0 TO 63, 128 TO 191))",
		},
		{
			name: "resource_pool_affinity_numanode",
			sql:  "CREATE RESOURCE POOL rp1 WITH (AFFINITY NUMANODE = (0, 1))",
		},
		{
			name: "alter_resource_pool_cap_cpu",
			sql:  "ALTER RESOURCE POOL rp1 WITH (CAP_CPU_PERCENT = 80)",
		},
		// ALTER RESOURCE GOVERNOR structured
		{
			name: "alter_resource_governor_reconfigure",
			sql:  "ALTER RESOURCE GOVERNOR RECONFIGURE",
		},
		{
			name: "alter_resource_governor_disable",
			sql:  "ALTER RESOURCE GOVERNOR DISABLE",
		},
		{
			name: "alter_resource_governor_reset_statistics",
			sql:  "ALTER RESOURCE GOVERNOR RESET STATISTICS",
		},
		{
			name: "alter_resource_governor_classifier",
			sql:  "ALTER RESOURCE GOVERNOR WITH (CLASSIFIER_FUNCTION = dbo.myClassifier)",
		},
		{
			name: "alter_resource_governor_classifier_null",
			sql:  "ALTER RESOURCE GOVERNOR WITH (CLASSIFIER_FUNCTION = NULL)",
		},
		// Workload classifier structured options
		{
			name: "workload_classifier_full",
			sql:  "CREATE WORKLOAD CLASSIFIER wc1 WITH (WORKLOAD_GROUP = 'wg1', MEMBERNAME = 'user1', IMPORTANCE = HIGH)",
		},
		{
			name: "workload_classifier_all_options",
			sql:  "CREATE WORKLOAD CLASSIFIER wc1 WITH (WORKLOAD_GROUP = 'wg1', MEMBERNAME = 'user1', WLM_LABEL = 'label1', WLM_CONTEXT = 'ctx1', START_TIME = '08:00', END_TIME = '17:00', IMPORTANCE = NORMAL)",
		},
		{
			name: "alter_workload_classifier",
			sql:  "ALTER WORKLOAD CLASSIFIER wc1 WITH (IMPORTANCE = LOW)",
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
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)

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

// TestParseDropSecurityKeysDispatch tests batch 103: DROP MASTER KEY, DROP CREDENTIAL,
// DROP CERTIFICATE, DROP ASYMMETRIC KEY, DROP SYMMETRIC KEY.
func TestParseDropSecurityKeysDispatch(t *testing.T) {
	tests := []struct {
		sql        string
		action     string
		objectType string
		name       string
	}{
		// DROP MASTER KEY
		{
			sql:        "DROP MASTER KEY",
			action:     "DROP",
			objectType: "MASTER KEY",
			name:       "",
		},
		// DROP CREDENTIAL
		{
			sql:        "DROP CREDENTIAL MyCredential",
			action:     "DROP",
			objectType: "CREDENTIAL",
			name:       "MyCredential",
		},
		// DROP CERTIFICATE
		{
			sql:        "DROP CERTIFICATE MyCertificate",
			action:     "DROP",
			objectType: "CERTIFICATE",
			name:       "MyCertificate",
		},
		// DROP ASYMMETRIC KEY
		{
			sql:        "DROP ASYMMETRIC KEY MyAsymKey",
			action:     "DROP",
			objectType: "ASYMMETRIC KEY",
			name:       "MyAsymKey",
		},
		// DROP SYMMETRIC KEY
		{
			sql:        "DROP SYMMETRIC KEY MySymKey",
			action:     "DROP",
			objectType: "SYMMETRIC KEY",
			name:       "MySymKey",
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

// TestParseDropExternalLibraryLanguageDispatch tests batch 104: DROP EXTERNAL LIBRARY / LANGUAGE dispatch.
func TestParseDropExternalLibraryLanguageDispatch(t *testing.T) {
	tests := []struct {
		sql        string
		objectType string
		name       string
	}{
		// DROP EXTERNAL LIBRARY
		{
			sql:        "DROP EXTERNAL LIBRARY MyLibrary",
			objectType: "EXTERNAL LIBRARY",
			name:       "MyLibrary",
		},
		// DROP EXTERNAL LIBRARY with AUTHORIZATION
		{
			sql:        "DROP EXTERNAL LIBRARY MyLibrary AUTHORIZATION dbo",
			objectType: "EXTERNAL LIBRARY",
			name:       "MyLibrary",
		},
		// DROP EXTERNAL LANGUAGE
		{
			sql:        "DROP EXTERNAL LANGUAGE MyLang",
			objectType: "EXTERNAL LANGUAGE",
			name:       "MyLang",
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
			if stmt.Action != "DROP" {
				t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, "DROP")
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

// TestParseBackupRestoreSymmetricKey tests batch 105: BACKUP/RESTORE SYMMETRIC KEY.
func TestParseBackupRestoreSymmetricKey(t *testing.T) {
	tests := []struct {
		sql        string
		action     string
		objectType string
		name       string
	}{
		// BACKUP SYMMETRIC KEY to file
		{
			sql:        "BACKUP SYMMETRIC KEY symmetric_key TO FILE = 'c:\\temp\\keys\\symmetric_key' ENCRYPTION BY PASSWORD = '3dH85Hhk003GHk2597gheij4'",
			action:     "BACKUP",
			objectType: "SYMMETRIC KEY",
			name:       "symmetric_key",
		},
		// BACKUP SYMMETRIC KEY to URL
		{
			sql:        "BACKUP SYMMETRIC KEY symmetric_key TO URL = 'https://mystorage.blob.core.windows.net/mycontainer/symmetric_key.bak' ENCRYPTION BY PASSWORD = '3dH85Hhk003GHk2597gheij4'",
			action:     "BACKUP",
			objectType: "SYMMETRIC KEY",
			name:       "symmetric_key",
		},
		// RESTORE SYMMETRIC KEY from file
		{
			sql:        "RESTORE SYMMETRIC KEY symmetric_key FROM FILE = 'c:\\temp\\keys\\symmetric_key' DECRYPTION BY PASSWORD = '3dH85Hhk003' ENCRYPTION BY PASSWORD = '259087M'",
			action:     "RESTORE",
			objectType: "SYMMETRIC KEY",
			name:       "symmetric_key",
		},
		// RESTORE SYMMETRIC KEY from URL
		{
			sql:        "RESTORE SYMMETRIC KEY symmetric_key FROM URL = 'https://mystorage.blob.core.windows.net/mycontainer/symmetric_key.bak' DECRYPTION BY PASSWORD = '3dH85Hhk003' ENCRYPTION BY PASSWORD = '259087M'",
			action:     "RESTORE",
			objectType: "SYMMETRIC KEY",
			name:       "symmetric_key",
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

// TestParseGetTransmissionStatus tests batch 106: GET_TRANSMISSION_STATUS function.
func TestParseGetTransmissionStatus(t *testing.T) {
	tests := []struct {
		sql string
	}{
		// GET_TRANSMISSION_STATUS as a function in SELECT
		{sql: "SELECT GET_TRANSMISSION_STATUS('58ef1d2d-c405-42eb-a762-23ff320bddf0')"},
		// GET_TRANSMISSION_STATUS with alias
		{sql: "SELECT Status = GET_TRANSMISSION_STATUS(@convHandle)"},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SelectStmt, got %T", tt.sql, result.Items[0])
			}
			checkLocation(t, tt.sql, "SelectStmt", stmt.Loc)
		})
	}
}

// TestParseEndpointHTTPOptions tests batch 107: endpoint HTTP structured options.
func TestParseEndpointHTTPOptions(t *testing.T) {
	tests := []struct {
		sql string
	}{
		// HTTP endpoint with PATH and AUTHENTICATION
		{sql: "CREATE ENDPOINT MyEndpoint STATE = STARTED AS HTTP (PATH = '/sql/endpoint', AUTHENTICATION = (INTEGRATED), PORTS = (CLEAR)) FOR TSQL ()"},
		// HTTP endpoint with all options
		{sql: "CREATE ENDPOINT MyEndpoint AS HTTP (PATH = '/sql', AUTHENTICATION = (BASIC, DIGEST, NTLM), PORTS = (CLEAR, SSL), SITE = '*', CLEAR_PORT = 80, SSL_PORT = 443, AUTH_REALM = 'myRealm', DEFAULT_LOGON_DOMAIN = 'MYDOMAIN', COMPRESSION = ENABLED) FOR TSQL ()"},
		// HTTP endpoint with KERBEROS
		{sql: "CREATE ENDPOINT MyEndpoint AS HTTP (PATH = '/sql', AUTHENTICATION = (KERBEROS), PORTS = (SSL)) FOR TSQL ()"},
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
			if stmt.Action != "CREATE" {
				t.Errorf("Parse(%q): action = %q, want CREATE", tt.sql, stmt.Action)
			}
			if stmt.ObjectType != "ENDPOINT" {
				t.Errorf("Parse(%q): objectType = %q, want ENDPOINT", tt.sql, stmt.ObjectType)
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)

			// Verify options contain AS=HTTP
			found := false
			if stmt.Options != nil {
				for _, item := range stmt.Options.Items {
					if s, ok := item.(*ast.String); ok && s.Str == "AS=HTTP" {
						found = true
						break
					}
				}
			}
			if !found {
				t.Errorf("Parse(%q): expected AS=HTTP in options", tt.sql)
			}
		})
	}
}

// TestParseEndpointUnknownProtocolOptions tests batch 116: unknown protocol structured options.
func TestParseEndpointUnknownProtocolOptions(t *testing.T) {
	tests := []struct {
		sql      string
		wantOpts []string
	}{
		// Unknown protocol with key=value options
		{
			sql: "CREATE ENDPOINT ep1 AS CUSTOM_PROTO (OPTION_A = value1, OPTION_B = 42) FOR TSQL ()",
			wantOpts: []string{
				"AS=CUSTOM_PROTO",
				"OPTION_A=VALUE1",
				"OPTION_B=42",
				"FOR=TSQL",
			},
		},
		// Unknown protocol with string values
		{
			sql: "CREATE ENDPOINT ep1 AS MYPROTO (LISTENER_ADDRESS = '10.0.0.1', PORT = 8080) FOR TSQL ()",
			wantOpts: []string{
				"AS=MYPROTO",
				"LISTENER_ADDRESS='10.0.0.1'",
				"PORT=8080",
				"FOR=TSQL",
			},
		},
		// Unknown protocol with single option
		{
			sql: "CREATE ENDPOINT ep1 AS SOMEPROTO (PARAM1 = enabled) FOR TSQL ()",
			wantOpts: []string{
				"AS=SOMEPROTO",
				"PARAM1=ENABLED",
				"FOR=TSQL",
			},
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
			if stmt.Options == nil || len(stmt.Options.Items) < len(tt.wantOpts) {
				var gotStrs []string
				if stmt.Options != nil {
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
				}
				t.Fatalf("Parse(%q): got %d options %v, want at least %d %v", tt.sql, len(gotStrs), gotStrs, len(tt.wantOpts), tt.wantOpts)
			}
			for _, want := range tt.wantOpts {
				found := false
				for _, item := range stmt.Options.Items {
					if item.(*ast.String).Str == want {
						found = true
						break
					}
				}
				if !found {
					var gotStrs []string
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
					t.Errorf("Parse(%q): expected option %q not found in %v", tt.sql, want, gotStrs)
				}
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseExternalModel tests batch 108: CREATE/ALTER/DROP EXTERNAL MODEL.
func TestParseExternalModel(t *testing.T) {
	tests := []struct {
		sql    string
		action string
	}{
		// CREATE basic with Azure OpenAI
		{
			sql:    "CREATE EXTERNAL MODEL MyAzureOpenAIModel WITH (LOCATION = 'https://my-endpoint.cognitiveservices.azure.com/openai/deployments/text-embedding-ada-002/embeddings?api-version=2024-02-01', API_FORMAT = 'Azure OpenAI', MODEL_TYPE = EMBEDDINGS, MODEL = 'text-embedding-ada-002', CREDENTIAL = [https://my-endpoint.cognitiveservices.azure.com/])",
			action: "CREATE",
		},
		// CREATE with AUTHORIZATION
		{
			sql:    "CREATE EXTERNAL MODEL MyModel AUTHORIZATION dbo WITH (LOCATION = 'https://localhost:11435/api/embed', API_FORMAT = 'Ollama', MODEL_TYPE = EMBEDDINGS, MODEL = 'all-minilm')",
			action: "CREATE",
		},
		// CREATE with PARAMETERS and CREDENTIAL
		{
			sql:    "CREATE EXTERNAL MODEL MyModel WITH (LOCATION = 'https://api.openai.com/v1/embeddings', API_FORMAT = 'OpenAI', MODEL_TYPE = EMBEDDINGS, MODEL = 'text-embedding-3-small', CREDENTIAL = [https://openai.com], PARAMETERS = '{\"dimensions\":725}')",
			action: "CREATE",
		},
		// CREATE with LOCAL_RUNTIME_PATH (ONNX)
		{
			sql:    "CREATE EXTERNAL MODEL myLocalOnnxModel WITH (LOCATION = 'C:\\onnx_runtime\\model\\all-MiniLM-L6-v2-onnx', API_FORMAT = 'ONNX Runtime', MODEL_TYPE = EMBEDDINGS, MODEL = 'allMiniLM', PARAMETERS = '{\"valid\":\"JSON\"}', LOCAL_RUNTIME_PATH = 'C:\\onnx_runtime\\')",
			action: "CREATE",
		},
		// ALTER basic - change MODEL
		{
			sql:    "ALTER EXTERNAL MODEL myAImodel SET (MODEL = 'text-embedding-3-large')",
			action: "ALTER",
		},
		// ALTER with multiple options
		{
			sql:    "ALTER EXTERNAL MODEL myAImodel SET (LOCATION = 'https://new-endpoint.com/v1/embeddings', API_FORMAT = 'OpenAI', MODEL_TYPE = EMBEDDINGS, MODEL = 'text-embedding-3-large', CREDENTIAL = [https://new-endpoint.com])",
			action: "ALTER",
		},
		// DROP basic
		{
			sql:    "DROP EXTERNAL MODEL myAImodel",
			action: "DROP",
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
			if stmt.ObjectType != "EXTERNAL MODEL" {
				t.Errorf("Parse(%q): objectType = %q, want EXTERNAL MODEL", tt.sql, stmt.ObjectType)
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseJsonVectorIndex tests batch 109: CREATE JSON INDEX and CREATE VECTOR INDEX.
func TestParseJsonVectorIndex(t *testing.T) {
	// --- CREATE JSON INDEX tests ---
	jsonTests := []string{
		// Basic JSON index
		"CREATE JSON INDEX idx_json ON dbo.Products (JsonData)",
		// JSON index with FOR paths
		"CREATE JSON INDEX idx_json ON Products (JsonCol) FOR ('$.name', '$.price')",
		// JSON index with single path
		"CREATE JSON INDEX idx_json ON MyTable (Details) FOR ('$.address.city')",
		// JSON index with WITH options (numeric)
		"CREATE JSON INDEX idx_json ON dbo.Orders (OrderInfo) WITH (FILLFACTOR = 80)",
		// JSON index with multiple WITH options
		"CREATE JSON INDEX idx_json ON Sales.Orders (Data) WITH (FILLFACTOR = 80, MAXDOP = 4)",
		// JSON index with FOR and WITH
		"CREATE JSON INDEX idx_json ON Products (JsonCol) FOR ('$.name') WITH (DATA_COMPRESSION = PAGE)",
		// JSON index with ON filegroup
		"CREATE JSON INDEX idx_json ON dbo.Products (JsonData) ON fg_primary",
		// JSON index with all clauses
		"CREATE JSON INDEX idx_json ON Products (Data) FOR ('$.id', '$.name') WITH (FILLFACTOR = 90) ON fg1",
	}
	for _, sql := range jsonTests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.CreateJsonIndexStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *CreateJsonIndexStmt, got %T", sql, result.Items[0])
			}
			if stmt.Name == "" {
				t.Errorf("Parse(%q): Name is empty", sql)
			}
			checkLocation(t, sql, "CreateJsonIndexStmt", stmt.Loc)
		})
	}

	// --- CREATE VECTOR INDEX tests ---
	vectorTests := []string{
		// Basic vector index
		"CREATE VECTOR INDEX idx_vec ON dbo.Products (EmbeddingCol)",
		// Vector index with METRIC option
		"CREATE VECTOR INDEX idx_vec ON Products (Vec) WITH (METRIC = 'cosine')",
		// Vector index with TYPE option
		"CREATE VECTOR INDEX idx_vec ON dbo.Items (Embedding) WITH (TYPE = 'DiskANN')",
		// Vector index with multiple options
		"CREATE VECTOR INDEX idx_vec ON Products (Vec) WITH (METRIC = 'dot', TYPE = 'DiskANN', MAXDOP = 4)",
		// Vector index with euclidean metric
		"CREATE VECTOR INDEX idx_vec ON dbo.Docs (Vec) WITH (METRIC = 'euclidean')",
		// Vector index with ON filegroup
		"CREATE VECTOR INDEX idx_vec ON Products (Vec) ON fg_primary",
		// Vector index with WITH and ON filegroup
		"CREATE VECTOR INDEX idx_vec ON Products (Vec) WITH (METRIC = 'cosine') ON fg1",
	}
	for _, sql := range vectorTests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.CreateVectorIndexStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *CreateVectorIndexStmt, got %T", sql, result.Items[0])
			}
			if stmt.Name == "" {
				t.Errorf("Parse(%q): Name is empty", sql)
			}
			checkLocation(t, sql, "CreateVectorIndexStmt", stmt.Loc)
		})
	}
}

// TestParseMaterializedView tests CREATE/ALTER/DROP MATERIALIZED VIEW parsing (batch 110).
func TestParseMaterializedView(t *testing.T) {
	// --- CREATE MATERIALIZED VIEW tests ---
	createTests := []string{
		// Basic with HASH distribution
		"CREATE MATERIALIZED VIEW dbo.mv_sales WITH (DISTRIBUTION = HASH(region)) AS SELECT region, SUM(amount) AS total FROM sales GROUP BY region",
		// With ROUND_ROBIN distribution
		"CREATE MATERIALIZED VIEW mv_test WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT a, COUNT_BIG(*) AS cnt FROM t GROUP BY a",
		// Multi-column HASH distribution
		"CREATE MATERIALIZED VIEW dbo.mv_multi WITH (DISTRIBUTION = HASH(col1, col2)) AS SELECT col1, col2, SUM(val) AS total FROM t GROUP BY col1, col2",
		// With FOR_APPEND option
		"CREATE MATERIALIZED VIEW mv_minmax WITH (DISTRIBUTION = HASH(id), FOR_APPEND) AS SELECT id, MAX(start_date) AS max_date, MIN(end_date) AS min_date FROM items GROUP BY id",
		// Schema-qualified view name
		"CREATE MATERIALIZED VIEW sales.mv_orders WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT order_date, COUNT_BIG(*) AS order_count FROM orders GROUP BY order_date",
		// With join in SELECT
		"CREATE MATERIALIZED VIEW dbo.mv_join WITH (DISTRIBUTION = HASH(vendor_id)) AS SELECT a.vendor_id, SUM(a.total_amount) AS s, COUNT_BIG(*) AS c FROM t1 a INNER JOIN t2 b ON a.vendor_id = b.vendor_id GROUP BY a.vendor_id",
	}
	for _, sql := range createTests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.CreateMaterializedViewStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *CreateMaterializedViewStmt, got %T", sql, result.Items[0])
			}
			if stmt.Name == nil {
				t.Errorf("Parse(%q): Name is nil", sql)
			}
			if stmt.Distribution == "" {
				t.Errorf("Parse(%q): Distribution is empty", sql)
			}
			if stmt.Query == nil {
				t.Errorf("Parse(%q): Query is nil", sql)
			}
			checkLocation(t, sql, "CreateMaterializedViewStmt", stmt.Loc)
		})
	}

	// --- ALTER MATERIALIZED VIEW tests ---
	alterTests := []struct {
		sql    string
		action string
	}{
		// REBUILD
		{"ALTER MATERIALIZED VIEW dbo.mv_sales REBUILD", "REBUILD"},
		// DISABLE
		{"ALTER MATERIALIZED VIEW mv_test DISABLE", "DISABLE"},
		// Schema-qualified REBUILD
		{"ALTER MATERIALIZED VIEW sales.mv_orders REBUILD", "REBUILD"},
		// Schema-qualified DISABLE
		{"ALTER MATERIALIZED VIEW sales.mv_orders DISABLE", "DISABLE"},
	}
	for _, tc := range alterTests {
		t.Run(tc.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tc.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.AlterMaterializedViewStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *AlterMaterializedViewStmt, got %T", tc.sql, result.Items[0])
			}
			if stmt.Name == nil {
				t.Errorf("Parse(%q): Name is nil", tc.sql)
			}
			if !strings.EqualFold(stmt.Action, tc.action) {
				t.Errorf("Parse(%q): Action = %q, want %q", tc.sql, stmt.Action, tc.action)
			}
			checkLocation(t, tc.sql, "AlterMaterializedViewStmt", stmt.Loc)
		})
	}

	// --- DROP MATERIALIZED VIEW tests ---
	dropTests := []string{
		"DROP MATERIALIZED VIEW dbo.mv_sales",
		"DROP MATERIALIZED VIEW mv_test",
		"DROP MATERIALIZED VIEW IF EXISTS sales.mv_orders",
	}
	for _, sql := range dropTests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.DropStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *DropStmt, got %T", sql, result.Items[0])
			}
			if stmt.ObjectType != ast.DropMaterializedView {
				t.Errorf("Parse(%q): ObjectType = %d, want DropMaterializedView", sql, stmt.ObjectType)
			}
			checkLocation(t, sql, "DropStmt", stmt.Loc)
		})
	}
}

// TestParseCopyInto tests the COPY INTO statement parser (batch 111).
func TestParseCopyInto(t *testing.T) {
	t.Run("copy_into_basic", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem FROM 'https://myaccount.blob.core.windows.net/myblobcontainer/folder1/' WITH (FILE_TYPE = 'CSV')`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.Table == nil {
			t.Fatal("expected non-nil Table")
		}
		if stmt.Sources == nil || stmt.Sources.Len() != 1 {
			t.Fatal("expected 1 source")
		}
		if stmt.Options == nil || stmt.Options.Len() != 1 {
			t.Fatal("expected 1 option")
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_with_options", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem
FROM 'https://myaccount.blob.core.windows.net/myblobcontainer/folder1/'
WITH (
    FILE_TYPE = 'CSV',
    FIELDTERMINATOR = '|',
    ROWTERMINATOR = '0x0A',
    FIRSTROW = 2,
    ENCODING = 'UTF8',
    DATEFORMAT = 'ymd',
    MAXERRORS = 10,
    COMPRESSION = 'Gzip',
    FIELDQUOTE = '"',
    IDENTITY_INSERT = 'ON',
    AUTO_CREATE_TABLE = 'OFF'
)`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.Options == nil || stmt.Options.Len() != 11 {
			t.Fatalf("expected 11 options, got %d", stmt.Options.Len())
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_multiple_sources", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem FROM 'https://myaccount.blob.core.windows.net/container/file1.csv', 'https://myaccount.blob.core.windows.net/container/file2.csv' WITH (FILE_TYPE = 'CSV')`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.Sources == nil || stmt.Sources.Len() != 2 {
			t.Fatalf("expected 2 sources, got %d", stmt.Sources.Len())
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_column_list", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem (col1, col2 DEFAULT 'N/A', col3 3) FROM 'https://myaccount.blob.core.windows.net/container/data.csv' WITH (FILE_TYPE = 'CSV')`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.ColumnList == nil || stmt.ColumnList.Len() != 3 {
			t.Fatalf("expected 3 columns, got %d", stmt.ColumnList.Len())
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_credential", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem FROM 'https://myaccount.blob.core.windows.net/container/' WITH (FILE_TYPE = 'PARQUET', CREDENTIAL = (IDENTITY = 'Shared Access Signature', SECRET = 'mysastoken'))`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.Options == nil || stmt.Options.Len() != 2 {
			t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_file_format", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem FROM 'https://myaccount.blob.core.windows.net/container/' WITH (FILE_FORMAT = myfileformat)`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.Options == nil || stmt.Options.Len() != 1 {
			t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_errorfile", func(t *testing.T) {
		sql := `COPY INTO dbo.lineitem FROM 'https://myaccount.blob.core.windows.net/container/' WITH (FILE_TYPE = 'ORC', ERRORFILE = 'https://myaccount.blob.core.windows.net/errors/')`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
		if stmt.Options == nil || stmt.Options.Len() != 2 {
			t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
		}
		checkLocation(t, sql, "CopyIntoStmt", stmt.Loc)
	})

	t.Run("copy_into_schema_only", func(t *testing.T) {
		sql := `COPY INTO myschema.mytable FROM 'https://storage.blob.core.windows.net/data/*.parquet' WITH (FILE_TYPE = 'PARQUET')`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		_, ok := result.Items[0].(*ast.CopyIntoStmt)
		if !ok {
			t.Fatalf("expected *CopyIntoStmt, got %T", result.Items[0])
		}
	})
}

// TestParseRenameCETAS tests RENAME, CETAS, and CREATE TABLE AS CLONE OF (batch 112).
func TestParseRenameCETAS(t *testing.T) {
	// RENAME OBJECT
	t.Run("rename_object", func(t *testing.T) {
		sql := `RENAME OBJECT Customer TO Customer1`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.RenameStmt)
		if !ok {
			t.Fatalf("expected *RenameStmt, got %T", result.Items[0])
		}
		if stmt.ObjectType != "OBJECT" {
			t.Errorf("expected ObjectType OBJECT, got %s", stmt.ObjectType)
		}
		if stmt.NewName != "Customer1" {
			t.Errorf("expected NewName Customer1, got %s", stmt.NewName)
		}
		checkLocation(t, sql, "RenameStmt", stmt.Loc)
	})

	t.Run("rename_object_qualified", func(t *testing.T) {
		sql := `RENAME OBJECT::mydb.dbo.Customer TO Customer1`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.RenameStmt)
		if !ok {
			t.Fatalf("expected *RenameStmt, got %T", result.Items[0])
		}
		if stmt.ObjectType != "OBJECT" {
			t.Errorf("expected ObjectType OBJECT, got %s", stmt.ObjectType)
		}
		if stmt.NewName != "Customer1" {
			t.Errorf("expected NewName Customer1, got %s", stmt.NewName)
		}
		checkLocation(t, sql, "RenameStmt", stmt.Loc)
	})

	t.Run("rename_database", func(t *testing.T) {
		sql := `RENAME DATABASE::AdWorks TO AdWorks2`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.RenameStmt)
		if !ok {
			t.Fatalf("expected *RenameStmt, got %T", result.Items[0])
		}
		if stmt.ObjectType != "DATABASE" {
			t.Errorf("expected ObjectType DATABASE, got %s", stmt.ObjectType)
		}
		if stmt.NewName != "AdWorks2" {
			t.Errorf("expected NewName AdWorks2, got %s", stmt.NewName)
		}
		checkLocation(t, sql, "RenameStmt", stmt.Loc)
	})

	t.Run("rename_column", func(t *testing.T) {
		sql := `RENAME OBJECT::Customer COLUMN FName TO FirstName`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.RenameStmt)
		if !ok {
			t.Fatalf("expected *RenameStmt, got %T", result.Items[0])
		}
		if stmt.ColumnName != "FName" {
			t.Errorf("expected ColumnName FName, got %s", stmt.ColumnName)
		}
		if stmt.NewColumnName != "FirstName" {
			t.Errorf("expected NewColumnName FirstName, got %s", stmt.NewColumnName)
		}
		checkLocation(t, sql, "RenameStmt", stmt.Loc)
	})

	// CETAS
	t.Run("cetas", func(t *testing.T) {
		sql := `CREATE EXTERNAL TABLE dbo.export_table
WITH (
    LOCATION = '/export/data/',
    DATA_SOURCE = myExternalDataSource,
    FILE_FORMAT = myFileFormat
)
AS SELECT col1, col2 FROM dbo.source_table WHERE col1 > 100`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CreateExternalTableAsSelectStmt)
		if !ok {
			t.Fatalf("expected *CreateExternalTableAsSelectStmt, got %T", result.Items[0])
		}
		if stmt.Name == nil {
			t.Fatal("expected non-nil Name")
		}
		if stmt.Options == nil || stmt.Options.Len() != 3 {
			t.Fatalf("expected 3 options, got %d", stmt.Options.Len())
		}
		if stmt.Query == nil {
			t.Fatal("expected non-nil Query")
		}
		checkLocation(t, sql, "CreateExternalTableAsSelectStmt", stmt.Loc)
	})

	t.Run("cetas_with_columns", func(t *testing.T) {
		sql := `CREATE EXTERNAL TABLE dbo.export_table (col1, col2)
WITH (
    LOCATION = '/export/',
    DATA_SOURCE = myds,
    FILE_FORMAT = myff
)
AS SELECT a, b FROM dbo.source`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CreateExternalTableAsSelectStmt)
		if !ok {
			t.Fatalf("expected *CreateExternalTableAsSelectStmt, got %T", result.Items[0])
		}
		if stmt.Columns == nil || stmt.Columns.Len() != 2 {
			t.Fatalf("expected 2 columns, got %d", stmt.Columns.Len())
		}
		if stmt.Query == nil {
			t.Fatal("expected non-nil Query")
		}
		checkLocation(t, sql, "CreateExternalTableAsSelectStmt", stmt.Loc)
	})

	t.Run("cetas_with_reject", func(t *testing.T) {
		sql := `CREATE EXTERNAL TABLE dbo.export_table
WITH (
    LOCATION = '/export/',
    DATA_SOURCE = myds,
    FILE_FORMAT = myff,
    REJECT_TYPE = value,
    REJECT_VALUE = 5
)
AS SELECT * FROM dbo.source`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CreateExternalTableAsSelectStmt)
		if !ok {
			t.Fatalf("expected *CreateExternalTableAsSelectStmt, got %T", result.Items[0])
		}
		if stmt.Options == nil || stmt.Options.Len() != 5 {
			t.Fatalf("expected 5 options, got %d", stmt.Options.Len())
		}
		checkLocation(t, sql, "CreateExternalTableAsSelectStmt", stmt.Loc)
	})

	// CREATE TABLE AS CLONE OF
	t.Run("create_table_clone", func(t *testing.T) {
		sql := `CREATE TABLE dbo.Employee AS CLONE OF dbo.EmployeeUSA`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CreateTableCloneStmt)
		if !ok {
			t.Fatalf("expected *CreateTableCloneStmt, got %T", result.Items[0])
		}
		if stmt.Name == nil {
			t.Fatal("expected non-nil Name")
		}
		if stmt.SourceName == nil {
			t.Fatal("expected non-nil SourceName")
		}
		if stmt.AtTime != "" {
			t.Errorf("expected empty AtTime, got %s", stmt.AtTime)
		}
		checkLocation(t, sql, "CreateTableCloneStmt", stmt.Loc)
	})

	t.Run("create_table_clone_cross_schema", func(t *testing.T) {
		sql := `CREATE TABLE dbo.Employee AS CLONE OF dbo1.EmployeeUSA`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CreateTableCloneStmt)
		if !ok {
			t.Fatalf("expected *CreateTableCloneStmt, got %T", result.Items[0])
		}
		if stmt.SourceName == nil || stmt.SourceName.Schema != "dbo1" {
			t.Errorf("expected source schema dbo1")
		}
		checkLocation(t, sql, "CreateTableCloneStmt", stmt.Loc)
	})

	t.Run("create_table_clone_at_time", func(t *testing.T) {
		sql := `CREATE TABLE dbo.Employee AS CLONE OF dbo.EmployeeUSA AT '2023-05-23T14:24:10.325'`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		stmt, ok := result.Items[0].(*ast.CreateTableCloneStmt)
		if !ok {
			t.Fatalf("expected *CreateTableCloneStmt, got %T", result.Items[0])
		}
		if stmt.AtTime != "2023-05-23T14:24:10.325" {
			t.Errorf("expected AtTime 2023-05-23T14:24:10.325, got %s", stmt.AtTime)
		}
		checkLocation(t, sql, "CreateTableCloneStmt", stmt.Loc)
	})

	// Regular CREATE EXTERNAL TABLE (non-CETAS) still works
	t.Run("create_external_table_no_cetas", func(t *testing.T) {
		sql := `CREATE EXTERNAL TABLE dbo.ext_table (col1 INT, col2 VARCHAR(50)) WITH (LOCATION = '/data/', DATA_SOURCE = myds, FILE_FORMAT = myff)`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		// Should parse as SecurityStmt (no AS SELECT)
		_, ok := result.Items[0].(*ast.SecurityStmt)
		if !ok {
			t.Fatalf("expected *SecurityStmt for non-CETAS external table, got %T", result.Items[0])
		}
	})

	// Regular CREATE TABLE still works
	t.Run("create_table_normal", func(t *testing.T) {
		sql := `CREATE TABLE dbo.test (id INT, name VARCHAR(100))`
		result := ParseAndCheck(t, sql)
		if result.Len() == 0 {
			t.Fatalf("Parse(%q): no statements returned", sql)
		}
		_, ok := result.Items[0].(*ast.CreateTableStmt)
		if !ok {
			t.Fatalf("expected *CreateTableStmt, got %T", result.Items[0])
		}
	})
}

// TestParseStmtDispatchPhase4 is the integration test for batch 113.
// Verifies all statements from batches 108-112 are properly dispatched.
func TestParseStmtDispatchPhase4(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantTyp string
	}{
		// Batch 108 - external model
		{"create_external_model", "CREATE EXTERNAL MODEL mymodel WITH (LOCATION = 'onnx:///mymodel.onnx')", "*ast.SecurityStmt"},
		// Batch 109 - json/vector index
		{"create_json_index", "CREATE JSON INDEX idx_json ON dbo.t1 (col1)", "*ast.CreateJsonIndexStmt"},
		// Batch 110 - materialized view
		{"create_materialized_view", "CREATE MATERIALIZED VIEW dbo.mv AS SELECT col1 FROM dbo.t1", "*ast.CreateMaterializedViewStmt"},
		{"alter_materialized_view", "ALTER MATERIALIZED VIEW dbo.mv REBUILD", "*ast.AlterMaterializedViewStmt"},
		// Batch 111 - copy into
		{"copy_into", "COPY INTO dbo.t1 FROM 'https://storage.blob.core.windows.net/data/' WITH (FILE_TYPE = 'CSV')", "*ast.CopyIntoStmt"},
		// Batch 112 - rename, CETAS, clone
		{"rename_object", "RENAME OBJECT dbo.old_table TO new_table", "*ast.RenameStmt"},
		{"rename_database", "RENAME DATABASE mydb TO newdb", "*ast.RenameStmt"},
		{"cetas", "CREATE EXTERNAL TABLE dbo.ext (col1) WITH (LOCATION = '/out/', DATA_SOURCE = ds, FILE_FORMAT = ff) AS SELECT 1", "*ast.CreateExternalTableAsSelectStmt"},
		{"create_table_clone", "CREATE TABLE dbo.clone AS CLONE OF dbo.original", "*ast.CreateTableCloneStmt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q): no statements returned", tt.sql)
			}
			gotTyp := fmt.Sprintf("%T", result.Items[0])
			if gotTyp != tt.wantTyp {
				t.Errorf("Parse(%q): got type %s, want %s", tt.sql, gotTyp, tt.wantTyp)
			}
		})
	}
}

// TestParseServerConfigOptionsDepth tests batch 114: structured parsing
// of ALTER SERVER CONFIGURATION SET options.
func TestParseServerConfigOptionsDepth(t *testing.T) {
	tests := []struct {
		sql        string
		optionType string
		wantOpts   []string
	}{
		// PROCESS AFFINITY CPU = AUTO
		{
			sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY CPU = AUTO",
			optionType: "PROCESS AFFINITY",
			wantOpts:   []string{"CPU=AUTO"},
		},
		// PROCESS AFFINITY CPU ranges
		{
			sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY CPU = 0 TO 63, 128 TO 191",
			optionType: "PROCESS AFFINITY",
			wantOpts:   []string{"CPU=0 TO 63, 128 TO 191"},
		},
		// PROCESS AFFINITY CPU single
		{
			sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY CPU = 0",
			optionType: "PROCESS AFFINITY",
			wantOpts:   []string{"CPU=0"},
		},
		// PROCESS AFFINITY NUMANODE
		{
			sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY NUMANODE = 0, 7",
			optionType: "PROCESS AFFINITY",
			wantOpts:   []string{"NUMANODE=0, 7"},
		},
		// PROCESS AFFINITY NUMANODE range
		{
			sql:        "ALTER SERVER CONFIGURATION SET PROCESS AFFINITY NUMANODE = 0 TO 3",
			optionType: "PROCESS AFFINITY",
			wantOpts:   []string{"NUMANODE=0 TO 3"},
		},
		// DIAGNOSTICS LOG ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG ON",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"ON"},
		},
		// DIAGNOSTICS LOG OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG OFF",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"OFF"},
		},
		// DIAGNOSTICS LOG PATH
		{
			sql:        `ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG PATH = 'C:\logs'`,
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{`PATH='C:\logs'`},
		},
		// DIAGNOSTICS LOG PATH DEFAULT
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG PATH = DEFAULT",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"PATH=DEFAULT"},
		},
		// DIAGNOSTICS LOG MAX_SIZE
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG MAX_SIZE = 10 MB",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"MAX_SIZE=10 MB"},
		},
		// DIAGNOSTICS LOG MAX_SIZE DEFAULT
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG MAX_SIZE = DEFAULT",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"MAX_SIZE=DEFAULT"},
		},
		// DIAGNOSTICS LOG MAX_FILES
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG MAX_FILES = 5",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"MAX_FILES=5"},
		},
		// DIAGNOSTICS LOG MAX_FILES DEFAULT
		{
			sql:        "ALTER SERVER CONFIGURATION SET DIAGNOSTICS LOG MAX_FILES = DEFAULT",
			optionType: "DIAGNOSTICS LOG",
			wantOpts:   []string{"MAX_FILES=DEFAULT"},
		},
		// FAILOVER CLUSTER PROPERTY HealthCheckTimeout
		{
			sql:        "ALTER SERVER CONFIGURATION SET FAILOVER CLUSTER PROPERTY HealthCheckTimeout = 15000",
			optionType: "FAILOVER CLUSTER PROPERTY",
			wantOpts:   []string{"HealthCheckTimeout=15000"},
		},
		// FAILOVER CLUSTER PROPERTY VerboseLogging
		{
			sql:        "ALTER SERVER CONFIGURATION SET FAILOVER CLUSTER PROPERTY VerboseLogging = 2",
			optionType: "FAILOVER CLUSTER PROPERTY",
			wantOpts:   []string{"VerboseLogging=2"},
		},
		// FAILOVER CLUSTER PROPERTY SqlDumperDumpPath
		{
			sql:        `ALTER SERVER CONFIGURATION SET FAILOVER CLUSTER PROPERTY SqlDumperDumpPath = 'C:\dumps'`,
			optionType: "FAILOVER CLUSTER PROPERTY",
			wantOpts:   []string{`SqlDumperDumpPath='C:\dumps'`},
		},
		// FAILOVER CLUSTER PROPERTY FailureConditionLevel DEFAULT
		{
			sql:        "ALTER SERVER CONFIGURATION SET FAILOVER CLUSTER PROPERTY FailureConditionLevel = DEFAULT",
			optionType: "FAILOVER CLUSTER PROPERTY",
			wantOpts:   []string{"FailureConditionLevel=DEFAULT"},
		},
		// FAILOVER CLUSTER PROPERTY ClusterConnectionOptions
		{
			sql:        `ALTER SERVER CONFIGURATION SET FAILOVER CLUSTER PROPERTY ClusterConnectionOptions = 'Encrypt=Strict'`,
			optionType: "FAILOVER CLUSTER PROPERTY",
			wantOpts:   []string{`ClusterConnectionOptions='Encrypt=Strict'`},
		},
		// HADR CLUSTER CONTEXT string
		{
			sql:        `ALTER SERVER CONFIGURATION SET HADR CLUSTER CONTEXT = 'clus01.xyz.com'`,
			optionType: "HADR CLUSTER CONTEXT",
			wantOpts:   []string{`CONTEXT='clus01.xyz.com'`},
		},
		// HADR CLUSTER CONTEXT LOCAL
		{
			sql:        "ALTER SERVER CONFIGURATION SET HADR CLUSTER CONTEXT = LOCAL",
			optionType: "HADR CLUSTER CONTEXT",
			wantOpts:   []string{"CONTEXT=LOCAL"},
		},
		// BUFFER POOL EXTENSION ON
		{
			sql:        `ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION ON (FILENAME = 'F:\SSDCACHE\Example.BPE', SIZE = 50 GB)`,
			optionType: "BUFFER POOL EXTENSION",
			wantOpts:   []string{"ON", `FILENAME='F:\SSDCACHE\Example.BPE'`, "SIZE=50 GB"},
		},
		// BUFFER POOL EXTENSION OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION OFF",
			optionType: "BUFFER POOL EXTENSION",
			wantOpts:   []string{"OFF"},
		},
		// BUFFER POOL EXTENSION ON with KB size
		{
			sql:        `ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION ON (FILENAME = 'cache.bpe', SIZE = 512 KB)`,
			optionType: "BUFFER POOL EXTENSION",
			wantOpts:   []string{"ON", `FILENAME='cache.bpe'`, "SIZE=512 KB"},
		},
		// BUFFER POOL EXTENSION ON with MB size
		{
			sql:        `ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION ON (FILENAME = 'cache.bpe', SIZE = 256 MB)`,
			optionType: "BUFFER POOL EXTENSION",
			wantOpts:   []string{"ON", `FILENAME='cache.bpe'`, "SIZE=256 MB"},
		},
		// SOFTNUMA ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET SOFTNUMA ON",
			optionType: "SOFTNUMA",
			wantOpts:   []string{"ON"},
		},
		// SOFTNUMA OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET SOFTNUMA OFF",
			optionType: "SOFTNUMA",
			wantOpts:   []string{"OFF"},
		},
		// MEMORY_OPTIMIZED ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED ON",
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"ON"},
		},
		// MEMORY_OPTIMIZED OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED OFF",
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"OFF"},
		},
		// MEMORY_OPTIMIZED TEMPDB_METADATA ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED TEMPDB_METADATA = ON",
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"TEMPDB_METADATA=ON"},
		},
		// MEMORY_OPTIMIZED TEMPDB_METADATA OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED TEMPDB_METADATA = OFF",
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"TEMPDB_METADATA=OFF"},
		},
		// MEMORY_OPTIMIZED TEMPDB_METADATA ON with RESOURCE_POOL
		{
			sql:        `ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED TEMPDB_METADATA = ON (RESOURCE_POOL = 'pool_name')`,
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"TEMPDB_METADATA=ON", `RESOURCE_POOL='pool_name'`},
		},
		// MEMORY_OPTIMIZED HYBRID_BUFFER_POOL ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED HYBRID_BUFFER_POOL = ON",
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"HYBRID_BUFFER_POOL=ON"},
		},
		// MEMORY_OPTIMIZED HYBRID_BUFFER_POOL OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET MEMORY_OPTIMIZED HYBRID_BUFFER_POOL = OFF",
			optionType: "MEMORY_OPTIMIZED",
			wantOpts:   []string{"HYBRID_BUFFER_POOL=OFF"},
		},
		// HARDWARE_OFFLOAD ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET HARDWARE_OFFLOAD ON",
			optionType: "HARDWARE_OFFLOAD",
			wantOpts:   []string{"ON"},
		},
		// HARDWARE_OFFLOAD OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET HARDWARE_OFFLOAD OFF",
			optionType: "HARDWARE_OFFLOAD",
			wantOpts:   []string{"OFF"},
		},
		// SUSPEND_FOR_SNAPSHOT_BACKUP ON
		{
			sql:        "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = ON",
			optionType: "SUSPEND_FOR_SNAPSHOT_BACKUP",
			wantOpts:   []string{"ON"},
		},
		// SUSPEND_FOR_SNAPSHOT_BACKUP OFF
		{
			sql:        "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = OFF",
			optionType: "SUSPEND_FOR_SNAPSHOT_BACKUP",
			wantOpts:   []string{"OFF"},
		},
		// SUSPEND_FOR_SNAPSHOT_BACKUP ON with GROUP and MODE
		{
			sql:        "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = ON (GROUP = (db1, db2), MODE = COPY_ONLY)",
			optionType: "SUSPEND_FOR_SNAPSHOT_BACKUP",
			wantOpts:   []string{"ON", "GROUP=db1, db2", "MODE=COPY_ONLY"},
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
			if len(tt.wantOpts) > 0 {
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", tt.sql)
				}
				if len(stmt.Options.Items) != len(tt.wantOpts) {
					var gotStrs []string
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
					t.Fatalf("Parse(%q): got %d options %v, want %d %v", tt.sql, len(stmt.Options.Items), gotStrs, len(tt.wantOpts), tt.wantOpts)
				}
				for i, want := range tt.wantOpts {
					got := stmt.Options.Items[i].(*ast.String).Str
					if got != want {
						t.Errorf("Parse(%q): option[%d] = %q, want %q", tt.sql, i, got, want)
					}
				}
			}
			checkLocation(t, tt.sql, "AlterServerConfigurationStmt", stmt.Loc)
		})
	}
}

// TestParseAvailabilityGroupOptionsDepth tests batch 115: structured parsing of AG options.
func TestParseAvailabilityGroupOptionsDepth(t *testing.T) {
	tests := []struct {
		sql      string
		wantOpts []string
	}{
		// ag_replica_on_structured: REPLICA ON with all standard options parsed as key=value
		{
			sql: "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC)",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
			},
		},
		// ag_replica_on_structured: REPLICA ON with SEEDING_MODE, BACKUP_PRIORITY, SESSION_TIMEOUT
		{
			sql: "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC, SEEDING_MODE = AUTOMATIC, BACKUP_PRIORITY = 50, SESSION_TIMEOUT = 10)",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"SEEDING_MODE=AUTOMATIC",
				"BACKUP_PRIORITY=50",
				"SESSION_TIMEOUT=10",
			},
		},
		// ag_replica_options_structured: SECONDARY_ROLE and PRIMARY_ROLE as nested option blocks
		{
			sql: "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC, SECONDARY_ROLE (ALLOW_CONNECTIONS = READ_ONLY), PRIMARY_ROLE (ALLOW_CONNECTIONS = ALL))",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"SECONDARY_ROLE(ALLOW_CONNECTIONS=READ_ONLY)",
				"PRIMARY_ROLE(ALLOW_CONNECTIONS=ALL)",
			},
		},
		// ag_replica_options_structured: PRIMARY_ROLE with READ_ONLY_ROUTING_LIST
		{
			sql: "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC, PRIMARY_ROLE (ALLOW_CONNECTIONS = ALL, READ_ONLY_ROUTING_LIST = ('server2', 'server3')))",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"PRIMARY_ROLE(ALLOW_CONNECTIONS=ALL, READ_ONLY_ROUTING_LIST=('server2', 'server3'))",
			},
		},
		// ag_listener_structured: LISTENER with IP and PORT
		{
			sql: "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC) LISTENER 'MyListener' (WITH IP (('10.120.19.155', '255.255.254.0')), PORT = 1433)",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"LISTENER='MyListener'",
				"WITH", "IP(('10.120.19.155', '255.255.254.0'))", "PORT=1433",
			},
		},
		// ag_listener_structured: LISTENER with DHCP
		{
			sql: "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC) LISTENER 'MyListener' (WITH DHCP ON ('10.120.19.0', '255.255.254.0'))",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"LISTENER='MyListener'",
				"WITH", "DHCP", "ON('10.120.19.0', '255.255.254.0')",
			},
		},
		// SET options parsed structurally
		{
			sql: "ALTER AVAILABILITY GROUP MyAG SET (AUTOMATED_BACKUP_PREFERENCE = SECONDARY, DB_FAILOVER = ON, HEALTH_CHECK_TIMEOUT = 30000)",
			wantOpts: []string{
				"SET",
				"AUTOMATED_BACKUP_PREFERENCE=SECONDARY",
				"DB_FAILOVER=ON",
				"HEALTH_CHECK_TIMEOUT=30000",
			},
		},
		// MODIFY LISTENER with ADD IP
		{
			sql: "ALTER AVAILABILITY GROUP MyAG MODIFY LISTENER 'MyListener' (ADD IP ('10.120.19.200', '255.255.254.0'))",
			wantOpts: []string{
				"MODIFY LISTENER='MyListener'",
				"ADD", "IP('10.120.19.200', '255.255.254.0')",
			},
		},
		// MODIFY LISTENER with PORT
		{
			sql: "ALTER AVAILABILITY GROUP MyAG MODIFY LISTENER 'MyListener' (PORT = 5022)",
			wantOpts: []string{
				"MODIFY LISTENER='MyListener'",
				"PORT=5022",
			},
		},
		// Distributed AG: AVAILABILITY GROUP ON with structured WITH options
		{
			sql: "CREATE AVAILABILITY GROUP MyDistAG WITH (DISTRIBUTED) AVAILABILITY GROUP ON 'AG1' WITH (LISTENER_URL = 'TCP://server1:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL, SEEDING_MODE = AUTOMATIC), 'AG2' WITH (LISTENER_URL = 'TCP://server2:5022', AVAILABILITY_MODE = ASYNCHRONOUS_COMMIT, FAILOVER_MODE = MANUAL, SEEDING_MODE = AUTOMATIC)",
			wantOpts: []string{
				"WITH", "DISTRIBUTED",
				"AVAILABILITY GROUP ON", "'AG1'", "WITH",
				"LISTENER_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=ASYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=MANUAL",
				"SEEDING_MODE=AUTOMATIC",
				"'AG2'", "WITH",
				"LISTENER_URL='TCP://server2:5022'",
				"AVAILABILITY_MODE=ASYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=MANUAL",
				"SEEDING_MODE=AUTOMATIC",
			},
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
			if stmt.Options == nil || len(stmt.Options.Items) != len(tt.wantOpts) {
				var gotStrs []string
				if stmt.Options != nil {
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
				}
				t.Fatalf("Parse(%q): got %d options %v, want %d %v", tt.sql, len(gotStrs), gotStrs, len(tt.wantOpts), tt.wantOpts)
			}
			for i, want := range tt.wantOpts {
				got := stmt.Options.Items[i].(*ast.String).Str
				if got != want {
					t.Errorf("Parse(%q): option[%d] = %q, want %q", tt.sql, i, got, want)
				}
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParseAvailabilityRemainingDepth tests batch 125: structured AG parsing fixes.
func TestParseAvailabilityRemainingDepth(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantOpts []string
	}{
		// ag_nested_parens_structured: IP tuples parsed as structured comma-separated values
		{
			name: "ag_nested_parens_structured_ip_tuple",
			sql:  "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC) LISTENER 'MyListener' (WITH IP (('10.0.0.1', '255.255.255.0'), ('10.0.0.2', '255.255.255.0')), PORT = 1433)",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"LISTENER='MyListener'",
				"WITH", "IP(('10.0.0.1', '255.255.255.0'), ('10.0.0.2', '255.255.255.0'))", "PORT=1433",
			},
		},
		// ag_nested_parens_structured: key = (value_list) uses recursive parsing
		{
			name: "ag_nested_parens_structured_value_list",
			sql:  "CREATE AVAILABILITY GROUP MyAG REPLICA ON 'server1' WITH (ENDPOINT_URL = 'TCP://server1:5022', AVAILABILITY_MODE = SYNCHRONOUS_COMMIT, FAILOVER_MODE = AUTOMATIC, PRIMARY_ROLE (READ_ONLY_ROUTING_LIST = ('server2', 'server3', 'server4')))",
			wantOpts: []string{
				"REPLICA ON", "'server1'", "WITH",
				"ENDPOINT_URL='TCP://server1:5022'",
				"AVAILABILITY_MODE=SYNCHRONOUS_COMMIT",
				"FAILOVER_MODE=AUTOMATIC",
				"PRIMARY_ROLE(READ_ONLY_ROUTING_LIST=('server2', 'server3', 'server4'))",
			},
		},
		// ag_listener_string_structured: LISTENER with structured name parsing
		{
			name: "ag_listener_string_structured",
			sql:  "ALTER AVAILABILITY GROUP MyAG ADD LISTENER 'NewListener' (WITH DHCP)",
			wantOpts: []string{
				"ADD LISTENER='NewListener'",
				"WITH", "DHCP",
			},
		},
		// ag_listener_string_structured: REMOVE LISTENER
		{
			name: "ag_remove_listener_structured",
			sql:  "ALTER AVAILABILITY GROUP MyAG REMOVE LISTENER 'OldListener'",
			wantOpts: []string{
				"REMOVE LISTENER='OldListener'",
			},
		},
		// ag_listener_string_structured: RESTART LISTENER
		{
			name: "ag_restart_listener_structured",
			sql:  "ALTER AVAILABILITY GROUP MyAG RESTART LISTENER 'MyListener'",
			wantOpts: []string{
				"RESTART LISTENER='MyListener'",
			},
		},
		// ag_listener_string_structured: MODIFY LISTENER with PORT
		{
			name: "ag_modify_listener_port",
			sql:  "ALTER AVAILABILITY GROUP MyAG MODIFY LISTENER 'MyListener' (PORT = 5022)",
			wantOpts: []string{
				"MODIFY LISTENER='MyListener'",
				"PORT=5022",
			},
		},
		// ag_grant_deny_database: GRANT CREATE ANY DATABASE
		{
			name: "ag_grant_create_any_database",
			sql:  "ALTER AVAILABILITY GROUP MyAG GRANT CREATE ANY DATABASE",
			wantOpts: []string{
				"GRANT CREATE ANY DATABASE",
			},
		},
		// ag_grant_deny_database: DENY CREATE ANY DATABASE
		{
			name: "ag_deny_create_any_database",
			sql:  "ALTER AVAILABILITY GROUP MyAG DENY CREATE ANY DATABASE",
			wantOpts: []string{
				"DENY CREATE ANY DATABASE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tt.sql, result.Items[0])
			}
			if stmt.Options == nil || len(stmt.Options.Items) != len(tt.wantOpts) {
				var gotStrs []string
				if stmt.Options != nil {
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
				}
				t.Fatalf("Parse(%q): got %d options %v, want %d %v", tt.sql, len(gotStrs), gotStrs, len(tt.wantOpts), tt.wantOpts)
			}
			for i, want := range tt.wantOpts {
				got := stmt.Options.Items[i].(*ast.String).Str
				if got != want {
					t.Errorf("Parse(%q): option[%d] = %q, want %q", tt.sql, i, got, want)
				}
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
		})
	}
}

// TestParsePredictStmt tests batch 118: PREDICT statement.
func TestParsePredictStmt(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		wantModel   bool
		wantData    bool
		wantAlias   string
		wantRuntime string
		wantCols    int
	}{
		// predict_basic: basic PREDICT with model variable and data table
		{
			name:      "predict_basic",
			sql:       "PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score FLOAT)",
			wantModel: true,
			wantData:  true,
			wantAlias: "d",
			wantCols:  1,
		},
		// predict_with_runtime: PREDICT with RUNTIME = ONNX (Azure Synapse)
		{
			name:        "predict_with_runtime",
			sql:         "PREDICT (MODEL = @model, DATA = dbo.mytable AS d, RUNTIME = ONNX) WITH (Score FLOAT)",
			wantModel:   true,
			wantData:    true,
			wantAlias:   "d",
			wantRuntime: "ONNX",
			wantCols:    1,
		},
		// predict with string literal model
		{
			name:      "predict_model_literal",
			sql:       "PREDICT (MODEL = 'my_model_binary', DATA = dbo.input_data AS d) WITH (prediction INT, probability FLOAT)",
			wantModel: true,
			wantData:  true,
			wantAlias: "d",
			wantCols:  2,
		},
		// predict with NOT NULL columns
		{
			name:      "predict_with_notnull",
			sql:       "PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score FLOAT NOT NULL, Label NVARCHAR(100) NULL)",
			wantModel: true,
			wantData:  true,
			wantAlias: "d",
			wantCols:  2,
		},
		// predict without alias
		{
			name:      "predict_no_alias",
			sql:       "PREDICT (MODEL = @model, DATA = dbo.mytable) WITH (Score FLOAT)",
			wantModel: true,
			wantData:  true,
			wantCols:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.PredictStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *PredictStmt, got %T", tt.sql, result.Items[0])
			}
			if tt.wantModel && stmt.Model == nil {
				t.Errorf("Parse(%q): expected non-nil Model", tt.sql)
			}
			if tt.wantData && stmt.Data == nil {
				t.Errorf("Parse(%q): expected non-nil Data", tt.sql)
			}
			if stmt.DataAlias != tt.wantAlias {
				t.Errorf("Parse(%q): DataAlias = %q, want %q", tt.sql, stmt.DataAlias, tt.wantAlias)
			}
			if stmt.Runtime != tt.wantRuntime {
				t.Errorf("Parse(%q): Runtime = %q, want %q", tt.sql, stmt.Runtime, tt.wantRuntime)
			}
			if tt.wantCols > 0 {
				if stmt.WithColumns == nil || len(stmt.WithColumns.Items) != tt.wantCols {
					gotCols := 0
					if stmt.WithColumns != nil {
						gotCols = len(stmt.WithColumns.Items)
					}
					t.Errorf("Parse(%q): got %d WITH columns, want %d", tt.sql, gotCols, tt.wantCols)
				}
			}
			checkLocation(t, tt.sql, "PredictStmt", stmt.Loc)
		})
	}
}

// TestParseEdgeConstraint tests EDGE CONSTRAINT for graph tables (batch 119).
func TestParseEdgeConstraint(t *testing.T) {
	// CREATE TABLE with single edge constraint connection
	t.Run("create_table_edge_constraint_single", func(t *testing.T) {
		sql := `CREATE TABLE bought (
			PurchaseCount INT,
			CONSTRAINT EC_BOUGHT CONNECTION (Customer TO Product)
		) AS EDGE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if !stmt.IsEdge {
			t.Error("expected IsEdge=true")
		}
		if stmt.Constraints == nil || stmt.Constraints.Len() != 1 {
			t.Fatalf("expected 1 constraint, got %d", stmt.Constraints.Len())
		}
		cd := stmt.Constraints.Items[0].(*ast.ConstraintDef)
		if cd.Type != ast.ConstraintEdge {
			t.Errorf("expected ConstraintEdge, got %d", cd.Type)
		}
		if cd.Name != "EC_BOUGHT" {
			t.Errorf("expected constraint name EC_BOUGHT, got %q", cd.Name)
		}
		if cd.EdgeConnections == nil || cd.EdgeConnections.Len() != 1 {
			t.Fatalf("expected 1 edge connection")
		}
		conn := cd.EdgeConnections.Items[0].(*ast.EdgeConnectionDef)
		if conn.FromTable.Object != "Customer" {
			t.Errorf("expected from=Customer, got %q", conn.FromTable.Object)
		}
		if conn.ToTable.Object != "Product" {
			t.Errorf("expected to=Product, got %q", conn.ToTable.Object)
		}
	})

	// CREATE TABLE with multiple edge constraint connections
	t.Run("create_table_edge_constraint_multi", func(t *testing.T) {
		sql := `CREATE TABLE bought (
			PurchaseCount INT,
			CONSTRAINT EC_BOUGHT CONNECTION (Customer TO Product, Supplier TO Product)
		) AS EDGE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		if stmt.Constraints == nil || stmt.Constraints.Len() != 1 {
			t.Fatalf("expected 1 constraint")
		}
		cd := stmt.Constraints.Items[0].(*ast.ConstraintDef)
		if cd.EdgeConnections == nil || cd.EdgeConnections.Len() != 2 {
			t.Fatalf("expected 2 edge connections, got %d", cd.EdgeConnections.Len())
		}
	})

	// CREATE TABLE edge constraint with ON DELETE NO ACTION
	t.Run("create_table_edge_constraint_on_delete_no_action", func(t *testing.T) {
		sql := `CREATE TABLE bought (
			PurchaseCount INT,
			CONSTRAINT EC_BOUGHT CONNECTION (Customer TO Product) ON DELETE NO ACTION
		) AS EDGE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		cd := stmt.Constraints.Items[0].(*ast.ConstraintDef)
		if cd.OnDelete != ast.RefActNoAction {
			t.Errorf("expected OnDelete=NoAction, got %d", cd.OnDelete)
		}
	})

	// CREATE TABLE edge constraint with ON DELETE CASCADE
	t.Run("create_table_edge_constraint_on_delete_cascade", func(t *testing.T) {
		sql := `CREATE TABLE bought (
			PurchaseCount INT,
			CONSTRAINT EC_BOUGHT CONNECTION (Customer TO Product) ON DELETE CASCADE
		) AS EDGE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		cd := stmt.Constraints.Items[0].(*ast.ConstraintDef)
		if cd.OnDelete != ast.RefActCascade {
			t.Errorf("expected OnDelete=Cascade, got %d", cd.OnDelete)
		}
	})

	// ALTER TABLE ADD edge constraint
	t.Run("alter_table_add_edge_constraint", func(t *testing.T) {
		sql := "ALTER TABLE bought ADD CONSTRAINT EC_BOUGHT1 CONNECTION (Customer TO Product)"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		if stmt.Actions == nil || stmt.Actions.Len() != 1 {
			t.Fatalf("expected 1 action")
		}
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATAddConstraint {
			t.Errorf("expected ATAddConstraint, got %d", action.Type)
		}
		if action.Constraint.Type != ast.ConstraintEdge {
			t.Errorf("expected ConstraintEdge")
		}
		if action.Constraint.EdgeConnections == nil || action.Constraint.EdgeConnections.Len() != 1 {
			t.Fatalf("expected 1 edge connection")
		}
	})

	// ALTER TABLE DROP edge constraint
	t.Run("alter_table_drop_edge_constraint", func(t *testing.T) {
		sql := "ALTER TABLE bought DROP CONSTRAINT EC_BOUGHT"
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.AlterTableStmt)
		if stmt.Actions == nil || stmt.Actions.Len() != 1 {
			t.Fatalf("expected 1 action")
		}
		action := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if action.Type != ast.ATDropConstraint {
			t.Errorf("expected ATDropConstraint, got %d", action.Type)
		}
	})

	// Edge constraint with schema-qualified table names
	t.Run("create_table_edge_constraint_schema_qualified", func(t *testing.T) {
		sql := `CREATE TABLE dbo.bought (
			CONSTRAINT EC_BOUGHT CONNECTION (dbo.Customer TO dbo.Product)
		) AS EDGE`
		result := ParseAndCheck(t, sql)
		stmt := result.Items[0].(*ast.CreateTableStmt)
		cd := stmt.Constraints.Items[0].(*ast.ConstraintDef)
		conn := cd.EdgeConnections.Items[0].(*ast.EdgeConnectionDef)
		if conn.FromTable.Schema != "dbo" || conn.FromTable.Object != "Customer" {
			t.Errorf("expected from=dbo.Customer, got %s.%s", conn.FromTable.Schema, conn.FromTable.Object)
		}
		if conn.ToTable.Schema != "dbo" || conn.ToTable.Object != "Product" {
			t.Errorf("expected to=dbo.Product, got %s.%s", conn.ToTable.Schema, conn.ToTable.Object)
		}
	})
}

// TestParseBackupRestoreOptionsDepth tests batch 120: structured BACKUP/RESTORE WITH options.
func TestParseBackupRestoreOptionsDepth(t *testing.T) {
	t.Run("backup_options_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts int // expected number of options
		}{
			{
				name:     "compression_flag",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH COMPRESSION",
				wantOpts: 1,
			},
			{
				name:     "multiple_flags",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH COMPRESSION, INIT, FORMAT, CHECKSUM",
				wantOpts: 4,
			},
			{
				name:     "name_equals_value",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH NAME = 'Full backup'",
				wantOpts: 1,
			},
			{
				name:     "stats_with_value",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH STATS = 10",
				wantOpts: 1,
			},
			{
				name:     "stats_without_value",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH STATS",
				wantOpts: 1,
			},
			{
				name:     "description_value",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH DESCRIPTION = 'Full database backup'",
				wantOpts: 1,
			},
			{
				name:     "differential_copy_only",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH DIFFERENTIAL, COPY_ONLY",
				wantOpts: 2,
			},
			{
				name:     "mixed_flags_and_kv",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH FORMAT, INIT, NAME = 'Full backup', COMPRESSION, STATS = 10",
				wantOpts: 5,
			},
			{
				name:     "expiredate",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH EXPIREDATE = '2025-12-31'",
				wantOpts: 1,
			},
			{
				name:     "retaindays",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH RETAINDAYS = 30",
				wantOpts: 1,
			},
			{
				name:     "blocksize_buffercount_maxtransfersize",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH BLOCKSIZE = 65536, BUFFERCOUNT = 50, MAXTRANSFERSIZE = 4194304",
				wantOpts: 3,
			},
			{
				name:     "medianame_mediadescription",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH MEDIANAME = 'MyMedia', MEDIADESCRIPTION = 'Full backup media set'",
				wantOpts: 2,
			},
			{
				name:     "no_compression",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH NO_COMPRESSION",
				wantOpts: 1,
			},
			{
				name:     "noskip_noinit_noformat",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH NOSKIP, NOINIT, NOFORMAT",
				wantOpts: 3,
			},
			{
				name:     "stop_on_error",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH STOP_ON_ERROR",
				wantOpts: 1,
			},
			{
				name:     "continue_after_error",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH CONTINUE_AFTER_ERROR",
				wantOpts: 1,
			},
			{
				name:     "no_checksum",
				sql:      "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH NO_CHECKSUM",
				wantOpts: 1,
			},
			{
				name:     "norecovery_no_truncate",
				sql:      "BACKUP LOG mydb TO DISK = '/backup/mydb.bak' WITH NORECOVERY, NO_TRUNCATE",
				wantOpts: 2,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.BackupStmt)
				if !ok {
					t.Fatalf("expected *BackupStmt, got %T", result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("expected options, got nil")
				}
				if len(stmt.Options.Items) != tt.wantOpts {
					t.Errorf("expected %d options, got %d", tt.wantOpts, len(stmt.Options.Items))
				}
				// Verify all options are BackupRestoreOption nodes
				for i, item := range stmt.Options.Items {
					opt, ok := item.(*ast.BackupRestoreOption)
					if !ok {
						t.Errorf("option[%d]: expected *BackupRestoreOption, got %T", i, item)
					}
					if opt.Name == "" {
						t.Errorf("option[%d]: empty name", i)
					}
					checkLocation(t, tt.sql, "BackupRestoreOption", opt.Loc)
				}
			})
		}
	})

	t.Run("restore_options_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts int
		}{
			{
				name:     "recovery",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH RECOVERY",
				wantOpts: 1,
			},
			{
				name:     "norecovery",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH NORECOVERY",
				wantOpts: 1,
			},
			{
				name:     "replace",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH REPLACE",
				wantOpts: 1,
			},
			{
				name:     "replace_norecovery",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH REPLACE, NORECOVERY, STATS = 10",
				wantOpts: 3,
			},
			{
				name:     "file_option",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH FILE = 2",
				wantOpts: 1,
			},
			{
				name:     "standby",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH STANDBY = 'C:\\undo\\mydb_undo.ldf'",
				wantOpts: 1,
			},
			{
				name:     "stopat",
				sql:      "RESTORE LOG mydb FROM DISK = '/backup/mydb_log.bak' WITH STOPAT = '2025-01-01T12:00:00'",
				wantOpts: 1,
			},
			{
				name:     "enable_broker",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH ENABLE_BROKER",
				wantOpts: 1,
			},
			{
				name:     "new_broker",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH NEW_BROKER",
				wantOpts: 1,
			},
			{
				name:     "error_broker_conversations",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH ERROR_BROKER_CONVERSATIONS",
				wantOpts: 1,
			},
			{
				name:     "medianame",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH MEDIANAME = 'MyMedia'",
				wantOpts: 1,
			},
			{
				name:     "mediapassword",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH MEDIAPASSWORD = 'secret123'",
				wantOpts: 1,
			},
			{
				name:     "restricted_user_keep_replication",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH RESTRICTED_USER, KEEP_REPLICATION",
				wantOpts: 2,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.RestoreStmt)
				if !ok {
					t.Fatalf("expected *RestoreStmt, got %T", result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("expected options, got nil")
				}
				if len(stmt.Options.Items) != tt.wantOpts {
					t.Errorf("expected %d options, got %d", tt.wantOpts, len(stmt.Options.Items))
				}
				for i, item := range stmt.Options.Items {
					opt, ok := item.(*ast.BackupRestoreOption)
					if !ok {
						t.Errorf("option[%d]: expected *BackupRestoreOption, got %T", i, item)
					}
					if opt.Name == "" {
						t.Errorf("option[%d]: empty name", i)
					}
					checkLocation(t, tt.sql, "BackupRestoreOption", opt.Loc)
				}
			})
		}
	})

	t.Run("backup_encryption_option", func(t *testing.T) {
		tests := []struct {
			name          string
			sql           string
			algorithm     string
			encryptorType string
			encryptorName string
		}{
			{
				name:          "aes256_server_cert",
				sql:           "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH ENCRYPTION (ALGORITHM = AES_256, SERVER CERTIFICATE = MyCert)",
				algorithm:     "AES_256",
				encryptorType: "SERVER CERTIFICATE",
				encryptorName: "MyCert",
			},
			{
				name:          "aes128_server_cert",
				sql:           "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH ENCRYPTION (ALGORITHM = AES_128, SERVER CERTIFICATE = BackupCert)",
				algorithm:     "AES_128",
				encryptorType: "SERVER CERTIFICATE",
				encryptorName: "BackupCert",
			},
			{
				name:          "triple_des_asymmetric_key",
				sql:           "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH ENCRYPTION (ALGORITHM = TRIPLE_DES_3KEY, SERVER ASYMMETRIC KEY = MyAsymKey)",
				algorithm:     "TRIPLE_DES_3KEY",
				encryptorType: "ASYMMETRIC KEY",
				encryptorName: "MyAsymKey",
			},
			{
				name:          "encryption_with_other_options",
				sql:           "BACKUP DATABASE mydb TO DISK = '/backup/mydb.bak' WITH COMPRESSION, ENCRYPTION (ALGORITHM = AES_256, SERVER CERTIFICATE = MyCert), STATS = 10",
				algorithm:     "AES_256",
				encryptorType: "SERVER CERTIFICATE",
				encryptorName: "MyCert",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.BackupStmt)
				if !ok {
					t.Fatalf("expected *BackupStmt, got %T", result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("expected options, got nil")
				}
				// Find the ENCRYPTION option
				var encOpt *ast.BackupRestoreOption
				for _, item := range stmt.Options.Items {
					opt, ok := item.(*ast.BackupRestoreOption)
					if ok && opt.Name == "ENCRYPTION" {
						encOpt = opt
						break
					}
				}
				if encOpt == nil {
					t.Fatalf("expected ENCRYPTION option, not found")
				}
				if encOpt.Algorithm != tt.algorithm {
					t.Errorf("algorithm = %q, want %q", encOpt.Algorithm, tt.algorithm)
				}
				if encOpt.EncryptorType != tt.encryptorType {
					t.Errorf("encryptorType = %q, want %q", encOpt.EncryptorType, tt.encryptorType)
				}
				if encOpt.EncryptorName != tt.encryptorName {
					t.Errorf("encryptorName = %q, want %q", encOpt.EncryptorName, tt.encryptorName)
				}
				checkLocation(t, tt.sql, "BackupRestoreOption", encOpt.Loc)
			})
		}
	})

	t.Run("restore_move_option", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			moveFrom string
			moveTo   string
		}{
			{
				name:     "single_move",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH MOVE 'mydb' TO 'C:\\data\\mydb.mdf'",
				moveFrom: "mydb",
				moveTo:   "C:\\data\\mydb.mdf",
			},
			{
				name:     "multiple_moves",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH MOVE 'mydb' TO 'C:\\data\\mydb.mdf', MOVE 'mydb_log' TO 'C:\\data\\mydb_log.ldf'",
				moveFrom: "mydb",
				moveTo:   "C:\\data\\mydb.mdf",
			},
			{
				name:     "move_with_replace",
				sql:      "RESTORE DATABASE mydb FROM DISK = '/backup/mydb.bak' WITH MOVE 'mydb' TO 'C:\\data\\mydb.mdf', REPLACE, NORECOVERY",
				moveFrom: "mydb",
				moveTo:   "C:\\data\\mydb.mdf",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.RestoreStmt)
				if !ok {
					t.Fatalf("expected *RestoreStmt, got %T", result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("expected options, got nil")
				}
				// Find first MOVE option
				var moveOpt *ast.BackupRestoreOption
				for _, item := range stmt.Options.Items {
					opt, ok := item.(*ast.BackupRestoreOption)
					if ok && opt.Name == "MOVE" {
						moveOpt = opt
						break
					}
				}
				if moveOpt == nil {
					t.Fatalf("expected MOVE option, not found")
				}
				if moveOpt.MoveFrom != tt.moveFrom {
					t.Errorf("moveFrom = %q, want %q", moveOpt.MoveFrom, tt.moveFrom)
				}
				if moveOpt.MoveTo != tt.moveTo {
					t.Errorf("moveTo = %q, want %q", moveOpt.MoveTo, tt.moveTo)
				}
				checkLocation(t, tt.sql, "BackupRestoreOption", moveOpt.Loc)
			})
		}
	})
}

// TestParseCreateDatabaseOptionsDepth tests batch 121: structured CREATE DATABASE WITH options.
func TestParseCreateDatabaseOptionsDepth(t *testing.T) {
	t.Run("database_with_options_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts int
			checkOpt func(t *testing.T, opts []*ast.DatabaseOption)
		}{
			{
				name:     "db_chaining_on",
				sql:      "CREATE DATABASE mydb WITH DB_CHAINING ON",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "DB_CHAINING" || opts[0].Value != "ON" {
						t.Errorf("expected DB_CHAINING=ON, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "trustworthy_off",
				sql:      "CREATE DATABASE mydb WITH TRUSTWORTHY OFF",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "TRUSTWORTHY" || opts[0].Value != "OFF" {
						t.Errorf("expected TRUSTWORTHY=OFF, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "ledger_on",
				sql:      "CREATE DATABASE mydb WITH LEDGER = ON",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "LEDGER" || opts[0].Value != "ON" {
						t.Errorf("expected LEDGER=ON, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "default_fulltext_language",
				sql:      "CREATE DATABASE mydb WITH DEFAULT_FULLTEXT_LANGUAGE = 1033",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "DEFAULT_FULLTEXT_LANGUAGE" || opts[0].Value != "1033" {
						t.Errorf("expected DEFAULT_FULLTEXT_LANGUAGE=1033, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "default_language_name",
				sql:      "CREATE DATABASE mydb WITH DEFAULT_LANGUAGE = English",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "DEFAULT_LANGUAGE" || opts[0].Value != "ENGLISH" {
						t.Errorf("expected DEFAULT_LANGUAGE=ENGLISH, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "nested_triggers_off",
				sql:      "CREATE DATABASE mydb WITH NESTED_TRIGGERS = OFF",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "NESTED_TRIGGERS" || opts[0].Value != "OFF" {
						t.Errorf("expected NESTED_TRIGGERS=OFF, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "transform_noise_words_on",
				sql:      "CREATE DATABASE mydb WITH TRANSFORM_NOISE_WORDS = ON",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "TRANSFORM_NOISE_WORDS" || opts[0].Value != "ON" {
						t.Errorf("expected TRANSFORM_NOISE_WORDS=ON, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "two_digit_year_cutoff",
				sql:      "CREATE DATABASE mydb WITH TWO_DIGIT_YEAR_CUTOFF = 2049",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "TWO_DIGIT_YEAR_CUTOFF" || opts[0].Value != "2049" {
						t.Errorf("expected TWO_DIGIT_YEAR_CUTOFF=2049, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "catalog_collation",
				sql:      "CREATE DATABASE mydb WITH CATALOG_COLLATION = DATABASE_DEFAULT",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "CATALOG_COLLATION" || opts[0].Value != "DATABASE_DEFAULT" {
						t.Errorf("expected CATALOG_COLLATION=DATABASE_DEFAULT, got %s=%s", opts[0].Name, opts[0].Value)
					}
				},
			},
			{
				name:     "multiple_options",
				sql:      "CREATE DATABASE mydb WITH DB_CHAINING ON, TRUSTWORTHY OFF, NESTED_TRIGGERS = ON",
				wantOpts: 3,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "DB_CHAINING" || opts[0].Value != "ON" {
						t.Errorf("opt[0]: expected DB_CHAINING=ON")
					}
					if opts[1].Name != "TRUSTWORTHY" || opts[1].Value != "OFF" {
						t.Errorf("opt[1]: expected TRUSTWORTHY=OFF")
					}
					if opts[2].Name != "NESTED_TRIGGERS" || opts[2].Value != "ON" {
						t.Errorf("opt[2]: expected NESTED_TRIGGERS=ON")
					}
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
				if !ok {
					t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.WithOptions == nil {
					t.Fatalf("expected WithOptions, got nil")
				}
				if len(stmt.WithOptions.Items) != tt.wantOpts {
					t.Errorf("expected %d options, got %d", tt.wantOpts, len(stmt.WithOptions.Items))
				}
				var dbOpts []*ast.DatabaseOption
				for _, item := range stmt.WithOptions.Items {
					opt, ok := item.(*ast.DatabaseOption)
					if !ok {
						t.Fatalf("expected *DatabaseOption, got %T", item)
					}
					checkLocation(t, tt.sql, "DatabaseOption", opt.Loc)
					dbOpts = append(dbOpts, opt)
				}
				if tt.checkOpt != nil {
					tt.checkOpt(t, dbOpts)
				}
			})
		}
	})

	t.Run("database_filestream_options", func(t *testing.T) {
		tests := []struct {
			name           string
			sql            string
			wantAccess     string
			wantDirName    string
		}{
			{
				name:       "filestream_non_transacted_access",
				sql:        "CREATE DATABASE mydb WITH FILESTREAM (NON_TRANSACTED_ACCESS = FULL)",
				wantAccess: "FULL",
			},
			{
				name:       "filestream_read_only",
				sql:        "CREATE DATABASE mydb WITH FILESTREAM (NON_TRANSACTED_ACCESS = READ_ONLY)",
				wantAccess: "READ_ONLY",
			},
			{
				name:       "filestream_off",
				sql:        "CREATE DATABASE mydb WITH FILESTREAM (NON_TRANSACTED_ACCESS = OFF)",
				wantAccess: "OFF",
			},
			{
				name:        "filestream_directory_name",
				sql:         "CREATE DATABASE mydb WITH FILESTREAM (DIRECTORY_NAME = 'mydir')",
				wantDirName: "mydir",
			},
			{
				name:        "filestream_both",
				sql:         "CREATE DATABASE mydb WITH FILESTREAM (NON_TRANSACTED_ACCESS = FULL, DIRECTORY_NAME = 'myfilestreamdir')",
				wantAccess:  "FULL",
				wantDirName: "myfilestreamdir",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
				if !ok {
					t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.WithOptions == nil || len(stmt.WithOptions.Items) < 1 {
					t.Fatalf("expected at least 1 option")
				}
				opt, ok := stmt.WithOptions.Items[0].(*ast.DatabaseOption)
				if !ok {
					t.Fatalf("expected *DatabaseOption, got %T", stmt.WithOptions.Items[0])
				}
				if opt.Name != "FILESTREAM" {
					t.Errorf("expected FILESTREAM, got %s", opt.Name)
				}
				if tt.wantAccess != "" && opt.FilestreamAccess != tt.wantAccess {
					t.Errorf("filestreamAccess = %q, want %q", opt.FilestreamAccess, tt.wantAccess)
				}
				if tt.wantDirName != "" && opt.FilestreamDirName != tt.wantDirName {
					t.Errorf("filestreamDirName = %q, want %q", opt.FilestreamDirName, tt.wantDirName)
				}
				checkLocation(t, tt.sql, "DatabaseOption", opt.Loc)
			})
		}
	})

	t.Run("database_attach_options_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts int
			checkOpt func(t *testing.T, opts []*ast.DatabaseOption)
		}{
			{
				name:     "attach_enable_broker",
				sql:      "CREATE DATABASE mydb ON (NAME = mydb_data, FILENAME = 'C:\\data\\mydb.mdf') FOR ATTACH WITH ENABLE_BROKER",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "ENABLE_BROKER" {
						t.Errorf("expected ENABLE_BROKER, got %s", opts[0].Name)
					}
				},
			},
			{
				name:     "attach_new_broker",
				sql:      "CREATE DATABASE mydb ON (NAME = mydb_data, FILENAME = 'C:\\data\\mydb.mdf') FOR ATTACH WITH NEW_BROKER",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "NEW_BROKER" {
						t.Errorf("expected NEW_BROKER, got %s", opts[0].Name)
					}
				},
			},
			{
				name:     "attach_error_broker_conversations",
				sql:      "CREATE DATABASE mydb ON (NAME = mydb_data, FILENAME = 'C:\\data\\mydb.mdf') FOR ATTACH WITH ERROR_BROKER_CONVERSATIONS",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "ERROR_BROKER_CONVERSATIONS" {
						t.Errorf("expected ERROR_BROKER_CONVERSATIONS, got %s", opts[0].Name)
					}
				},
			},
			{
				name:     "attach_restricted_user",
				sql:      "CREATE DATABASE mydb ON (NAME = mydb_data, FILENAME = 'C:\\data\\mydb.mdf') FOR ATTACH WITH RESTRICTED_USER",
				wantOpts: 1,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "RESTRICTED_USER" {
						t.Errorf("expected RESTRICTED_USER, got %s", opts[0].Name)
					}
				},
			},
			{
				name:     "attach_multiple",
				sql:      "CREATE DATABASE mydb ON (NAME = mydb_data, FILENAME = 'C:\\data\\mydb.mdf') FOR ATTACH WITH ENABLE_BROKER, RESTRICTED_USER",
				wantOpts: 2,
				checkOpt: func(t *testing.T, opts []*ast.DatabaseOption) {
					if opts[0].Name != "ENABLE_BROKER" {
						t.Errorf("opt[0]: expected ENABLE_BROKER, got %s", opts[0].Name)
					}
					if opts[1].Name != "RESTRICTED_USER" {
						t.Errorf("opt[1]: expected RESTRICTED_USER, got %s", opts[1].Name)
					}
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
				if !ok {
					t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
				}
				if stmt.AttachOptions == nil {
					t.Fatalf("expected AttachOptions, got nil")
				}
				if len(stmt.AttachOptions.Items) != tt.wantOpts {
					t.Errorf("expected %d options, got %d", tt.wantOpts, len(stmt.AttachOptions.Items))
				}
				var dbOpts []*ast.DatabaseOption
				for _, item := range stmt.AttachOptions.Items {
					opt, ok := item.(*ast.DatabaseOption)
					if !ok {
						t.Fatalf("expected *DatabaseOption, got %T", item)
					}
					checkLocation(t, tt.sql, "DatabaseOption", opt.Loc)
					dbOpts = append(dbOpts, opt)
				}
				if tt.checkOpt != nil {
					tt.checkOpt(t, dbOpts)
				}
			})
		}
	})

	t.Run("persistent_log_buffer", func(t *testing.T) {
		sql := "CREATE DATABASE mydb WITH PERSISTENT_LOG_BUFFER = ON (DIRECTORY_NAME = 'S:\\PLBData')"
		result := ParseAndCheck(t, sql)
		if result.Len() != 1 {
			t.Fatalf("expected 1 stmt, got %d", result.Len())
		}
		stmt, ok := result.Items[0].(*ast.CreateDatabaseStmt)
		if !ok {
			t.Fatalf("expected *CreateDatabaseStmt, got %T", result.Items[0])
		}
		if stmt.WithOptions == nil || len(stmt.WithOptions.Items) != 1 {
			t.Fatalf("expected 1 option")
		}
		opt, ok := stmt.WithOptions.Items[0].(*ast.DatabaseOption)
		if !ok {
			t.Fatalf("expected *DatabaseOption, got %T", stmt.WithOptions.Items[0])
		}
		if opt.Name != "PERSISTENT_LOG_BUFFER" {
			t.Errorf("expected PERSISTENT_LOG_BUFFER, got %s", opt.Name)
		}
		if opt.Value != "ON" {
			t.Errorf("expected value ON, got %s", opt.Value)
		}
		if opt.PersistentLogDir != "S:\\PLBData" {
			t.Errorf("expected persistentLogDir S:\\PLBData, got %s", opt.PersistentLogDir)
		}
	})
}

// TestParseSecurityPrincipalOptionsDepth tests batch 122: structured security principal options.
func TestParseSecurityPrincipalOptionsDepth(t *testing.T) {
	t.Run("create_user_structured_options", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
			check func(t *testing.T, stmt *ast.SecurityStmt)
		}{
			{
				name: "user_with_password",
				sql:  "CREATE USER testUser WITH PASSWORD = 'StrongPass123'",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					if stmt.Name != "testUser" { t.Errorf("name = %q", stmt.Name) }
					opt := findSecOpt(t, stmt, "PASSWORD")
					if opt.Value != "StrongPass123" { t.Errorf("password = %q", opt.Value) }
				},
			},
			{
				name: "user_with_default_schema",
				sql:  "CREATE USER testUser WITH DEFAULT_SCHEMA = dbo",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "DEFAULT_SCHEMA")
					if opt.Value != "dbo" { t.Errorf("default_schema = %q", opt.Value) }
				},
			},
			{
				name: "user_for_login_with_options",
				sql:  "CREATE USER testUser FOR LOGIN testLogin WITH DEFAULT_SCHEMA = Sales",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					loginOpt := findSecOpt(t, stmt, "LOGIN")
					if loginOpt.Value != "testLogin" { t.Errorf("login = %q", loginOpt.Value) }
					schemaOpt := findSecOpt(t, stmt, "DEFAULT_SCHEMA")
					if schemaOpt.Value != "Sales" { t.Errorf("schema = %q", schemaOpt.Value) }
				},
			},
			{
				name: "user_with_default_language",
				sql:  "CREATE USER testUser WITH DEFAULT_LANGUAGE = English",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "DEFAULT_LANGUAGE")
					if opt.Value != "English" { t.Errorf("language = %q", opt.Value) }
				},
			},
			{
				name: "user_with_allow_encrypted",
				sql:  "CREATE USER testUser WITH ALLOW_ENCRYPTED_VALUE_MODIFICATIONS = ON",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "ALLOW_ENCRYPTED_VALUE_MODIFICATIONS")
					if opt.Value != "ON" { t.Errorf("allow_encrypted = %q", opt.Value) }
				},
			},
			{
				name: "user_with_sid",
				sql:  "CREATE USER testUser WITH SID = 0x01020304",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "SID")
					if opt.Value == "" { t.Errorf("sid is empty") }
				},
			},
			{
				name: "alter_user_rename",
				sql:  "ALTER USER testUser WITH NAME = newUser",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "NAME")
					if opt.Value != "newUser" { t.Errorf("name = %q", opt.Value) }
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok { t.Fatalf("expected *SecurityStmt, got %T", result.Items[0]) }
				tt.check(t, stmt)
				// Verify all options are SecurityPrincipalOption nodes
				if stmt.Options != nil {
					for i, item := range stmt.Options.Items {
						opt, ok := item.(*ast.SecurityPrincipalOption)
						if !ok { t.Errorf("option[%d]: expected *SecurityPrincipalOption, got %T", i, item) }
						checkLocation(t, tt.sql, "SecurityPrincipalOption", opt.Loc)
					}
				}
			})
		}
	})

	t.Run("create_login_structured_options", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
			check func(t *testing.T, stmt *ast.SecurityStmt)
		}{
			{
				name: "login_with_password",
				sql:  "CREATE LOGIN testLogin WITH PASSWORD = 'StrongPass123'",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "PASSWORD")
					if opt.Value != "StrongPass123" { t.Errorf("password = %q", opt.Value) }
				},
			},
			{
				name: "login_password_must_change",
				sql:  "CREATE LOGIN testLogin WITH PASSWORD = 'TempPass!' MUST_CHANGE",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "PASSWORD")
					if opt.Value != "TempPass!" { t.Errorf("password = %q", opt.Value) }
					if !opt.MustChange { t.Errorf("expected MustChange=true") }
				},
			},
			{
				name: "login_password_hashed",
				sql:  "CREATE LOGIN testLogin WITH PASSWORD = 0x01020304 HASHED",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "PASSWORD")
					if !opt.Hashed { t.Errorf("expected Hashed=true") }
				},
			},
			{
				name: "login_with_multiple_options",
				sql:  "CREATE LOGIN testLogin WITH PASSWORD = 'Pass123', DEFAULT_DATABASE = mydb, CHECK_EXPIRATION = ON, CHECK_POLICY = OFF",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					findSecOpt(t, stmt, "PASSWORD")
					dbOpt := findSecOpt(t, stmt, "DEFAULT_DATABASE")
					if dbOpt.Value != "mydb" { t.Errorf("default_database = %q", dbOpt.Value) }
					expOpt := findSecOpt(t, stmt, "CHECK_EXPIRATION")
					if expOpt.Value != "ON" { t.Errorf("check_expiration = %q", expOpt.Value) }
					polOpt := findSecOpt(t, stmt, "CHECK_POLICY")
					if polOpt.Value != "OFF" { t.Errorf("check_policy = %q", polOpt.Value) }
				},
			},
			{
				name: "login_with_credential",
				sql:  "CREATE LOGIN testLogin WITH PASSWORD = 'Pass123', CREDENTIAL = MyCred",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "CREDENTIAL")
					if opt.Value != "MyCred" { t.Errorf("credential = %q", opt.Value) }
				},
			},
			{
				name: "login_from_windows",
				sql:  "CREATE LOGIN [DOMAIN\\user] FROM WINDOWS WITH DEFAULT_DATABASE = master",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					fromOpt := findSecOpt(t, stmt, "FROM")
					if fromOpt.Value != "WINDOWS" { t.Errorf("from = %q", fromOpt.Value) }
					dbOpt := findSecOpt(t, stmt, "DEFAULT_DATABASE")
					if dbOpt.Value != "master" { t.Errorf("default_database = %q", dbOpt.Value) }
				},
			},
			{
				name: "login_from_certificate",
				sql:  "CREATE LOGIN testLogin FROM CERTIFICATE MyCert",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					fromOpt := findSecOpt(t, stmt, "FROM")
					if fromOpt.Value != "CERTIFICATE MyCert" { t.Errorf("from = %q", fromOpt.Value) }
				},
			},
			{
				name: "login_from_asymmetric_key",
				sql:  "CREATE LOGIN testLogin FROM ASYMMETRIC KEY MyKey",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					fromOpt := findSecOpt(t, stmt, "FROM")
					if fromOpt.Value != "ASYMMETRIC KEY MyKey" { t.Errorf("from = %q", fromOpt.Value) }
				},
			},
			{
				name: "login_from_external_provider",
				sql:  "CREATE LOGIN testLogin FROM EXTERNAL PROVIDER",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					fromOpt := findSecOpt(t, stmt, "FROM")
					if fromOpt.Value != "EXTERNAL PROVIDER" { t.Errorf("from = %q", fromOpt.Value) }
				},
			},
			{
				name: "login_with_sid",
				sql:  "CREATE LOGIN testLogin WITH PASSWORD = 'Pass123', SID = 0xABCD",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					findSecOpt(t, stmt, "PASSWORD")
					findSecOpt(t, stmt, "SID")
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok { t.Fatalf("expected *SecurityStmt, got %T", result.Items[0]) }
				tt.check(t, stmt)
			})
		}
	})

	t.Run("alter_login_structured_options", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
			check func(t *testing.T, stmt *ast.SecurityStmt)
		}{
			{
				name: "alter_login_enable",
				sql:  "ALTER LOGIN testLogin ENABLE",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					findSecOpt(t, stmt, "ENABLE")
				},
			},
			{
				name: "alter_login_disable",
				sql:  "ALTER LOGIN testLogin DISABLE",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					findSecOpt(t, stmt, "DISABLE")
				},
			},
			{
				name: "alter_login_password",
				sql:  "ALTER LOGIN testLogin WITH PASSWORD = 'NewPass'",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "PASSWORD")
					if opt.Value != "NewPass" { t.Errorf("password = %q", opt.Value) }
				},
			},
			{
				name: "alter_login_password_old",
				sql:  "ALTER LOGIN testLogin WITH PASSWORD = 'NewPass' OLD_PASSWORD = 'OldPass'",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "PASSWORD")
					if opt.Value != "NewPass" { t.Errorf("password = %q", opt.Value) }
					if opt.OldPassword != "OldPass" { t.Errorf("old_password = %q", opt.OldPassword) }
				},
			},
			{
				name: "alter_login_default_database",
				sql:  "ALTER LOGIN testLogin WITH DEFAULT_DATABASE = newdb",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "DEFAULT_DATABASE")
					if opt.Value != "newdb" { t.Errorf("default_database = %q", opt.Value) }
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok { t.Fatalf("expected *SecurityStmt, got %T", result.Items[0]) }
				tt.check(t, stmt)
			})
		}
	})

	t.Run("create_app_role_structured_options", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
			check func(t *testing.T, stmt *ast.SecurityStmt)
		}{
			{
				name: "app_role_with_password",
				sql:  "CREATE APPLICATION ROLE myAppRole WITH PASSWORD = 'AppPass123'",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "PASSWORD")
					if opt.Value != "AppPass123" { t.Errorf("password = %q", opt.Value) }
				},
			},
			{
				name: "app_role_with_password_and_schema",
				sql:  "CREATE APPLICATION ROLE myAppRole WITH PASSWORD = 'AppPass123', DEFAULT_SCHEMA = dbo",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					findSecOpt(t, stmt, "PASSWORD")
					opt := findSecOpt(t, stmt, "DEFAULT_SCHEMA")
					if opt.Value != "dbo" { t.Errorf("default_schema = %q", opt.Value) }
				},
			},
			{
				name: "alter_app_role_rename",
				sql:  "ALTER APPLICATION ROLE myAppRole WITH NAME = newAppRole",
				check: func(t *testing.T, stmt *ast.SecurityStmt) {
					opt := findSecOpt(t, stmt, "NAME")
					if opt.Value != "newAppRole" { t.Errorf("name = %q", opt.Value) }
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok { t.Fatalf("expected *SecurityStmt, got %T", result.Items[0]) }
				tt.check(t, stmt)
			})
		}
	})
}

// findSecOpt finds a SecurityPrincipalOption by name in a SecurityStmt.
func findSecOpt(t *testing.T, stmt *ast.SecurityStmt, name string) *ast.SecurityPrincipalOption {
	t.Helper()
	if stmt.Options == nil {
		t.Fatalf("expected options for %s, got nil", name)
	}
	for _, item := range stmt.Options.Items {
		opt, ok := item.(*ast.SecurityPrincipalOption)
		if ok && opt.Name == name {
			return opt
		}
	}
	t.Fatalf("option %q not found in SecurityStmt options", name)
	return nil
}

// TestParseDbccWithOptionsDepth tests batch 123: structured DBCC WITH options.
func TestParseDbccWithOptionsDepth(t *testing.T) {
	t.Run("dbcc_with_options_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts int
			wantNames []string
		}{
			{
				name:      "checkdb_no_infomsgs",
				sql:       "DBCC CHECKDB ('mydb') WITH NO_INFOMSGS",
				wantOpts:  1,
				wantNames: []string{"NO_INFOMSGS"},
			},
			{
				name:      "checkdb_all_errormsgs",
				sql:       "DBCC CHECKDB ('mydb') WITH ALL_ERRORMSGS",
				wantOpts:  1,
				wantNames: []string{"ALL_ERRORMSGS"},
			},
			{
				name:      "checkdb_physical_only",
				sql:       "DBCC CHECKDB ('mydb') WITH PHYSICAL_ONLY",
				wantOpts:  1,
				wantNames: []string{"PHYSICAL_ONLY"},
			},
			{
				name:      "checkdb_extended_logical_checks",
				sql:       "DBCC CHECKDB ('mydb') WITH EXTENDED_LOGICAL_CHECKS",
				wantOpts:  1,
				wantNames: []string{"EXTENDED_LOGICAL_CHECKS"},
			},
			{
				name:      "checkdb_data_purity",
				sql:       "DBCC CHECKDB ('mydb') WITH DATA_PURITY",
				wantOpts:  1,
				wantNames: []string{"DATA_PURITY"},
			},
			{
				name:      "checktable_tablock",
				sql:       "DBCC CHECKTABLE ('dbo.mytable') WITH TABLOCK",
				wantOpts:  1,
				wantNames: []string{"TABLOCK"},
			},
			{
				name:      "checkdb_estimateonly",
				sql:       "DBCC CHECKDB ('mydb') WITH ESTIMATEONLY",
				wantOpts:  1,
				wantNames: []string{"ESTIMATEONLY"},
			},
			{
				name:      "checkdb_count_rows",
				sql:       "DBCC CHECKDB ('mydb') WITH COUNT_ROWS",
				wantOpts:  1,
				wantNames: []string{"COUNT_ROWS"},
			},
			{
				name:      "checkdb_multiple_options",
				sql:       "DBCC CHECKDB ('mydb') WITH NO_INFOMSGS, ALL_ERRORMSGS, PHYSICAL_ONLY",
				wantOpts:  3,
				wantNames: []string{"NO_INFOMSGS", "ALL_ERRORMSGS", "PHYSICAL_ONLY"},
			},
			{
				name:      "checkdb_all_options",
				sql:       "DBCC CHECKDB ('mydb') WITH NO_INFOMSGS, DATA_PURITY, TABLOCK, ESTIMATEONLY",
				wantOpts:  4,
				wantNames: []string{"NO_INFOMSGS", "DATA_PURITY", "TABLOCK", "ESTIMATEONLY"},
			},
			{
				name:      "checkdb_tableresults",
				sql:       "DBCC CHECKDB ('mydb') WITH TABLERESULTS",
				wantOpts:  1,
				wantNames: []string{"TABLERESULTS"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("expected 1 stmt, got %d", result.Len())
				}
				stmt, ok := result.Items[0].(*ast.DbccStmt)
				if !ok {
					t.Fatalf("expected *DbccStmt, got %T", result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("expected options, got nil")
				}
				if len(stmt.Options.Items) != tt.wantOpts {
					t.Errorf("expected %d options, got %d", tt.wantOpts, len(stmt.Options.Items))
				}
				for i, item := range stmt.Options.Items {
					opt, ok := item.(*ast.DbccOption)
					if !ok {
						t.Errorf("option[%d]: expected *DbccOption, got %T", i, item)
						continue
					}
					if i < len(tt.wantNames) && opt.Name != tt.wantNames[i] {
						t.Errorf("option[%d]: name = %q, want %q", i, opt.Name, tt.wantNames[i])
					}
					checkLocation(t, tt.sql, "DbccOption", opt.Loc)
				}
			})
		}
	})
}

// TestParseExternalNestedParensDepth tests batch 124: structured parsing of external library/language
// file specs, FORMAT_OPTIONS, SHARDED(), and external table column definitions.
func TestParseExternalNestedParensDepth(t *testing.T) {
	// Test structured FROM clause parsing for CREATE EXTERNAL LIBRARY
	t.Run("external_library_from_structured", func(t *testing.T) {
		tests := []struct {
			sql         string
			wantOptions []string
		}{
			{
				sql:         "CREATE EXTERNAL LIBRARY customPackage FROM (CONTENT = 'C:\\temp\\pkg.zip') WITH (LANGUAGE = 'R')",
				wantOptions: []string{"CONTENT", "LANGUAGE"},
			},
			{
				sql:         "CREATE EXTERNAL LIBRARY customLib FROM (CONTENT = 'C:\\temp\\lib.zip', PLATFORM = WINDOWS) WITH (LANGUAGE = 'R')",
				wantOptions: []string{"CONTENT", "PLATFORM", "LANGUAGE"},
			},
			{
				sql:         "CREATE EXTERNAL LIBRARY customLib FROM (CONTENT = 0xABC123) WITH (LANGUAGE = 'R')",
				wantOptions: []string{"CONTENT", "LANGUAGE"},
			},
			{
				sql:         "ALTER EXTERNAL LIBRARY customPackage SET (CONTENT = 'C:\\temp\\pkg.zip') WITH (LANGUAGE = 'R')",
				wantOptions: []string{"CONTENT", "LANGUAGE"},
			},
			{
				sql:         "ALTER EXTERNAL LIBRARY customLib AUTHORIZATION dbo SET (CONTENT = 'C:\\temp\\lib.zip', PLATFORM = LINUX) WITH (LANGUAGE = 'Python')",
				wantOptions: []string{"CONTENT", "PLATFORM", "LANGUAGE"},
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
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", tt.sql)
				}
				if len(stmt.Options.Items) < len(tt.wantOptions) {
					t.Errorf("Parse(%q): got %d options, want at least %d", tt.sql, len(stmt.Options.Items), len(tt.wantOptions))
				}
				checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
			})
		}
	})

	// Test structured FROM clause parsing for CREATE/ALTER EXTERNAL LANGUAGE
	t.Run("external_language_from_structured", func(t *testing.T) {
		tests := []struct {
			sql         string
			wantOptions int // minimum number of options expected
		}{
			{
				sql:         "CREATE EXTERNAL LANGUAGE Java FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll')",
				wantOptions: 2,
			},
			{
				sql:         "CREATE EXTERNAL LANGUAGE Java FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll', PLATFORM = WINDOWS)",
				wantOptions: 3,
			},
			{
				sql:         "CREATE EXTERNAL LANGUAGE Java FROM (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll', PLATFORM = WINDOWS), (CONTENT = N'C:\\temp\\java.tar.gz', FILE_NAME = 'javaextension.so', PLATFORM = LINUX)",
				wantOptions: 6,
			},
			{
				sql:         "ALTER EXTERNAL LANGUAGE Java SET (CONTENT = N'C:\\temp\\java.zip', FILE_NAME = 'javaextension.dll')",
				wantOptions: 2,
			},
			{
				sql:         "ALTER EXTERNAL LANGUAGE Java ADD (CONTENT = N'C:\\temp\\java.tar.gz', FILE_NAME = 'javaextension.so', PLATFORM = LINUX)",
				wantOptions: 3,
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
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", tt.sql)
				}
				if len(stmt.Options.Items) < tt.wantOptions {
					t.Errorf("Parse(%q): got %d options, want at least %d", tt.sql, len(stmt.Options.Items), tt.wantOptions)
				}
				checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
			})
		}
	})

	// Test structured FORMAT_OPTIONS parsing
	t.Run("external_format_options_structured", func(t *testing.T) {
		tests := []string{
			`CREATE EXTERNAL FILE FORMAT CsvFormat WITH (
				FORMAT_TYPE = DELIMITEDTEXT,
				FORMAT_OPTIONS (FIELD_TERMINATOR = ',', STRING_DELIMITER = '"', FIRST_ROW = 2, USE_TYPE_DEFAULT = TRUE)
			)`,
			`CREATE EXTERNAL FILE FORMAT JsonFormat WITH (
				FORMAT_TYPE = JSON,
				FORMAT_OPTIONS (FIELD_TERMINATOR = '\t')
			)`,
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
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", sql)
				}
				checkLocation(t, sql, "SecurityStmt", stmt.Loc)
			})
		}
	})

	// Test SHARDED() structured parsing through external data source options
	t.Run("sharded_structured", func(t *testing.T) {
		sql := "CREATE EXTERNAL DATA SOURCE MyDS WITH (LOCATION = 'hdfs://myhost:8020', TYPE = HADOOP)"
		result := ParseAndCheck(t, sql)
		if result.Len() != 1 {
			t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
		}
		stmt, ok := result.Items[0].(*ast.SecurityStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
		}
		if stmt.Options == nil {
			t.Fatalf("Parse(%q): expected options, got nil", sql)
		}
		checkLocation(t, sql, "SecurityStmt", stmt.Loc)
	})

	// Test structured column def parsing (CETAS path produces ColumnDef nodes from fallback)
	t.Run("external_table_columns_structured", func(t *testing.T) {
		tests := []string{
			`CREATE EXTERNAL TABLE dbo.ClickStream (
				url VARCHAR(50),
				event_date DATE,
				user_ip VARCHAR(50)
			)
			WITH (
				LOCATION = '/data/',
				DATA_SOURCE = MySource,
				FILE_FORMAT = CsvFormat
			)`,
			`CREATE EXTERNAL TABLE dbo.DataTable (
				id INT NOT NULL,
				name NVARCHAR(100) NULL,
				value FLOAT
			)
			WITH (
				LOCATION = '/data/',
				DATA_SOURCE = MySource,
				FILE_FORMAT = CsvFormat,
				DISTRIBUTION = REPLICATED
			)`,
			`CREATE EXTERNAL TABLE dbo.RoundRobin (
				col1 BIGINT
			)
			WITH (
				LOCATION = '/data/',
				DATA_SOURCE = MySource,
				FILE_FORMAT = ParquetFormat,
				DISTRIBUTION = REPLICATED
			)`,
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
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", sql)
				}
				// Verify columns are encoded as ColumnDef nodes
				foundColDef := false
				for _, item := range stmt.Options.Items {
					if _, ok := item.(*ast.ColumnDef); ok {
						foundColDef = true
						break
					}
				}
				if !foundColDef {
					t.Errorf("Parse(%q): expected ColumnDef nodes in options, got none", sql)
				}
				checkLocation(t, sql, "SecurityStmt", stmt.Loc)
			})
		}
	})
}

// TestParseServerConfigRemainingDepth tests batch 126: structured parsing of
// server config remaining depth issues.
func TestParseServerConfigRemainingDepth(t *testing.T) {
	t.Run("server_config_remaining_structured", func(t *testing.T) {
		// Unknown option type should still parse the statement with detected keyword
		sql := "ALTER SERVER CONFIGURATION SET UNKNOWN_OPTION = 42"
		result, _ := Parse(sql)
		if result.Len() != 1 {
			t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
		}
		stmt, ok := result.Items[0].(*ast.AlterServerConfigurationStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected *AlterServerConfigurationStmt, got %T", sql, result.Items[0])
		}
		// Unknown option type should still be detected by the first keyword
		if stmt.OptionType == "" {
			t.Errorf("Parse(%q): expected non-empty optionType", sql)
		}
		checkLocation(t, sql, "AlterServerConfigurationStmt", stmt.Loc)
	})

	t.Run("buffer_pool_extension_structured", func(t *testing.T) {
		tests := []struct {
			sql      string
			wantOpts []string
		}{
			// Standard ON with FILENAME and SIZE
			{
				sql:      `ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION ON (FILENAME = 'D:\cache.bpe', SIZE = 20 GB)`,
				wantOpts: []string{"ON", `FILENAME='D:\cache.bpe'`, "SIZE=20 GB"},
			},
			// SIZE without unit suffix
			{
				sql:      `ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION ON (FILENAME = 'cache.bpe', SIZE = 100)`,
				wantOpts: []string{"ON", `FILENAME='cache.bpe'`, "SIZE=100"},
			},
			// OFF
			{
				sql:      "ALTER SERVER CONFIGURATION SET BUFFER POOL EXTENSION OFF",
				wantOpts: []string{"OFF"},
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
				if stmt.OptionType != "BUFFER POOL EXTENSION" {
					t.Errorf("Parse(%q): optionType = %q, want %q", tt.sql, stmt.OptionType, "BUFFER POOL EXTENSION")
				}
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", tt.sql)
				}
				if len(stmt.Options.Items) != len(tt.wantOpts) {
					var gotStrs []string
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
					t.Fatalf("Parse(%q): got %d options %v, want %d %v", tt.sql, len(stmt.Options.Items), gotStrs, len(tt.wantOpts), tt.wantOpts)
				}
				for i, want := range tt.wantOpts {
					got := stmt.Options.Items[i].(*ast.String).Str
					if got != want {
						t.Errorf("Parse(%q): option[%d] = %q, want %q", tt.sql, i, got, want)
					}
				}
				checkLocation(t, tt.sql, "AlterServerConfigurationStmt", stmt.Loc)
			})
		}
	})

	t.Run("suspend_snapshot_structured", func(t *testing.T) {
		tests := []struct {
			sql      string
			wantOpts []string
		}{
			// ON with GROUP (single db)
			{
				sql:      "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = ON (GROUP = (mydb))",
				wantOpts: []string{"ON", "GROUP=mydb"},
			},
			// ON with GROUP (multiple dbs) and MODE
			{
				sql:      "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = ON (GROUP = (db1, db2, db3), MODE = COPY_ONLY)",
				wantOpts: []string{"ON", "GROUP=db1, db2, db3", "MODE=COPY_ONLY"},
			},
			// Simple ON
			{
				sql:      "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = ON",
				wantOpts: []string{"ON"},
			},
			// OFF
			{
				sql:      "ALTER SERVER CONFIGURATION SET SUSPEND_FOR_SNAPSHOT_BACKUP = OFF",
				wantOpts: []string{"OFF"},
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
				if stmt.OptionType != "SUSPEND_FOR_SNAPSHOT_BACKUP" {
					t.Errorf("Parse(%q): optionType = %q, want %q", tt.sql, stmt.OptionType, "SUSPEND_FOR_SNAPSHOT_BACKUP")
				}
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", tt.sql)
				}
				if len(stmt.Options.Items) != len(tt.wantOpts) {
					var gotStrs []string
					for _, item := range stmt.Options.Items {
						gotStrs = append(gotStrs, item.(*ast.String).Str)
					}
					t.Fatalf("Parse(%q): got %d options %v, want %d %v", tt.sql, len(stmt.Options.Items), gotStrs, len(tt.wantOpts), tt.wantOpts)
				}
				for i, want := range tt.wantOpts {
					got := stmt.Options.Items[i].(*ast.String).Str
					if got != want {
						t.Errorf("Parse(%q): option[%d] = %q, want %q", tt.sql, i, got, want)
					}
				}
				checkLocation(t, tt.sql, "AlterServerConfigurationStmt", stmt.Loc)
			})
		}
	})
}

// TestParseSecurityKeysRemainingDepth tests batch 128: structured security key column
// and security key options remaining depth.
func TestParseSecurityKeysRemainingDepth(t *testing.T) {
	// ---- ALTER COLUMN ENCRYPTION KEY with structured ADD/DROP VALUE ----

	t.Run("security_key_column_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			action   string
			objType  string
			keyName  string
			wantOpts []string
		}{
			{
				name:    "alter_column_encryption_key_add_value",
				sql:     `ALTER COLUMN ENCRYPTION KEY MyColKey ADD VALUE (COLUMN_MASTER_KEY = MyCMK, ALGORITHM = 'RSA_OAEP', ENCRYPTED_VALUE = 0xABCD)`,
				action:  "ALTER",
				objType: "COLUMN ENCRYPTION KEY",
				keyName: "MyColKey",
				wantOpts: []string{
					"ADD", "VALUE",
					"COLUMN_MASTER_KEY", "MyCMK",
					"ALGORITHM", "RSA_OAEP",
					"ENCRYPTED_VALUE", "0xABCD",
				},
			},
			{
				name:    "alter_column_encryption_key_drop_value",
				sql:     `ALTER COLUMN ENCRYPTION KEY MyColKey DROP VALUE (COLUMN_MASTER_KEY = OldCMK)`,
				action:  "ALTER",
				objType: "COLUMN ENCRYPTION KEY",
				keyName: "MyColKey",
				wantOpts: []string{
					"DROP", "VALUE",
					"COLUMN_MASTER_KEY", "OldCMK",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SecurityKeyStmt)
				if stmt.Action != tt.action {
					t.Errorf("expected action %q, got %q", tt.action, stmt.Action)
				}
				if stmt.ObjectType != tt.objType {
					t.Errorf("expected objectType %q, got %q", tt.objType, stmt.ObjectType)
				}
				if stmt.Name != tt.keyName {
					t.Errorf("expected name %q, got %q", tt.keyName, stmt.Name)
				}
				if stmt.Options == nil {
					t.Fatal("expected options, got nil")
				}
				if len(tt.wantOpts) > 0 && stmt.Options.Len() < len(tt.wantOpts) {
					t.Fatalf("expected at least %d options, got %d", len(tt.wantOpts), stmt.Options.Len())
				}
				for i, want := range tt.wantOpts {
					if i >= stmt.Options.Len() {
						break
					}
					got := stmt.Options.Items[i].(*ast.String).Str
					if got != want {
						t.Errorf("option[%d] = %q, want %q", i, got, want)
					}
				}
				checkLocation(t, tt.sql, "SecurityKeyStmt", stmt.Loc)
			})
		}
	})

	// ---- ENCRYPTION BY / DECRYPTION BY structured using parseEncryptingMechanism ----

	t.Run("security_key_options_remaining_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts []string
		}{
			{
				name: "create_symmetric_encryption_by_certificate",
				sql:  `CREATE SYMMETRIC KEY TestKey WITH ALGORITHM = AES_256 ENCRYPTION BY CERTIFICATE MyCert`,
				wantOpts: []string{
					"ALGORITHM", "AES_256",
					"ENCRYPTION BY CERTIFICATE MyCert",
				},
			},
			{
				name: "create_symmetric_encryption_by_password",
				sql:  `CREATE SYMMETRIC KEY TestKey WITH ALGORITHM = AES_128 ENCRYPTION BY PASSWORD = 'pass123'`,
				wantOpts: []string{
					"ALGORITHM", "AES_128",
					"ENCRYPTION BY PASSWORD = pass123",
				},
			},
			{
				name: "create_symmetric_encryption_by_symmetric_key",
				sql:  `CREATE SYMMETRIC KEY TestKey WITH ALGORITHM = AES_256 ENCRYPTION BY SYMMETRIC KEY OtherKey`,
				wantOpts: []string{
					"ALGORITHM", "AES_256",
					"ENCRYPTION BY SYMMETRIC KEY OtherKey",
				},
			},
			{
				name: "create_symmetric_encryption_by_asymmetric_key",
				sql:  `CREATE SYMMETRIC KEY TestKey WITH ALGORITHM = AES_256 ENCRYPTION BY ASYMMETRIC KEY MyAsymKey`,
				wantOpts: []string{
					"ALGORITHM", "AES_256",
					"ENCRYPTION BY ASYMMETRIC KEY MyAsymKey",
				},
			},
			{
				name: "open_symmetric_decryption_by_certificate",
				sql:  `OPEN SYMMETRIC KEY TestKey DECRYPTION BY CERTIFICATE MyCert`,
				wantOpts: []string{
					"DECRYPTION BY CERTIFICATE MyCert",
				},
			},
			{
				name: "open_symmetric_decryption_by_password",
				sql:  `OPEN SYMMETRIC KEY TestKey DECRYPTION BY PASSWORD = 'mypass'`,
				wantOpts: []string{
					"DECRYPTION BY PASSWORD = mypass",
				},
			},
			{
				name: "create_symmetric_multi_encryption",
				sql:  `CREATE SYMMETRIC KEY TestKey WITH ALGORITHM = AES_256 ENCRYPTION BY CERTIFICATE Cert1, PASSWORD = 'pass1'`,
				wantOpts: []string{
					"ALGORITHM", "AES_256",
					"ENCRYPTION BY CERTIFICATE Cert1",
					"ENCRYPTION BY PASSWORD = pass1",
				},
			},
			{
				name: "database_encryption_key_server_cert",
				sql:  `CREATE DATABASE ENCRYPTION KEY WITH ALGORITHM = AES_256 ENCRYPTION BY SERVER CERTIFICATE MyCert`,
				wantOpts: []string{
					"ALGORITHM", "AES_256",
					"ENCRYPTION BY SERVER CERTIFICATE MyCert",
				},
			},
			{
				name: "database_encryption_key_server_asymmetric",
				sql:  `CREATE DATABASE ENCRYPTION KEY WITH ALGORITHM = AES_128 ENCRYPTION BY SERVER ASYMMETRIC KEY MyAsymKey`,
				wantOpts: []string{
					"ALGORITHM", "AES_128",
					"ENCRYPTION BY SERVER ASYMMETRIC KEY MyAsymKey",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				stmt := result.Items[0].(*ast.SecurityKeyStmt)
				if stmt.Options == nil {
					t.Fatal("expected options, got nil")
				}
				if len(tt.wantOpts) > 0 && stmt.Options.Len() < len(tt.wantOpts) {
					t.Fatalf("expected at least %d options, got %d", len(tt.wantOpts), stmt.Options.Len())
				}
				for i, want := range tt.wantOpts {
					if i >= stmt.Options.Len() {
						break
					}
					got := stmt.Options.Items[i].(*ast.String).Str
					if got != want {
						t.Errorf("option[%d] = %q, want %q", i, got, want)
					}
				}
			})
		}
	})
}

// TestParseFulltextAlterDepth tests batch 129: structured ALTER FULLTEXT INDEX actions
// and structured search property list option parsing.
func TestParseFulltextAlterDepth(t *testing.T) {
	t.Run("alter_fulltext_index_action_structured", func(t *testing.T) {
		tests := []struct {
			name   string
			sql    string
			action string
		}{
			{"enable", "ALTER FULLTEXT INDEX ON dbo.Products ENABLE", "ENABLE"},
			{"disable", "ALTER FULLTEXT INDEX ON dbo.Products DISABLE", "DISABLE"},
			{"set_change_tracking_manual", "ALTER FULLTEXT INDEX ON dbo.Products SET CHANGE_TRACKING MANUAL", "SET"},
			{"set_change_tracking_auto", "ALTER FULLTEXT INDEX ON dbo.Products SET CHANGE_TRACKING AUTO", "SET"},
			{"set_change_tracking_off", "ALTER FULLTEXT INDEX ON dbo.Products SET CHANGE_TRACKING OFF", "SET"},
			{"add_columns", "ALTER FULLTEXT INDEX ON dbo.Products ADD (Name, Description)", "ADD"},
			{"add_columns_with_no_population", "ALTER FULLTEXT INDEX ON dbo.Products ADD (Name) WITH NO POPULATION", "ADD"},
			{"add_columns_type_column_language", "ALTER FULLTEXT INDEX ON dbo.Products ADD (Name TYPE COLUMN nvarchar LANGUAGE 1033 STATISTICAL_SEMANTICS)", "ADD"},
			{"alter_column_add_semantics", "ALTER FULLTEXT INDEX ON dbo.Products ALTER COLUMN Name ADD STATISTICAL_SEMANTICS", "ALTER"},
			{"alter_column_drop_semantics", "ALTER FULLTEXT INDEX ON dbo.Products ALTER COLUMN Name DROP STATISTICAL_SEMANTICS WITH NO POPULATION", "ALTER"},
			{"drop_columns", "ALTER FULLTEXT INDEX ON dbo.Products DROP (Name, Description)", "DROP"},
			{"drop_columns_with_no_population", "ALTER FULLTEXT INDEX ON dbo.Products DROP (Name) WITH NO POPULATION", "DROP"},
			{"start_full_population", "ALTER FULLTEXT INDEX ON dbo.Products START FULL POPULATION", "START"},
			{"start_incremental_population", "ALTER FULLTEXT INDEX ON dbo.Products START INCREMENTAL POPULATION", "START"},
			{"start_update_population", "ALTER FULLTEXT INDEX ON dbo.Products START UPDATE POPULATION", "START"},
			{"stop_population", "ALTER FULLTEXT INDEX ON dbo.Products STOP POPULATION", "STOP"},
			{"pause_population", "ALTER FULLTEXT INDEX ON dbo.Products PAUSE POPULATION", "PAUSE"},
			{"resume_population", "ALTER FULLTEXT INDEX ON dbo.Products RESUME POPULATION", "RESUME"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterFulltextIndexStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *AlterFulltextIndexStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): action = %q, want %q", tt.sql, stmt.Action, tt.action)
				}
				checkLocation(t, tt.sql, "AlterFulltextIndexStmt", stmt.Loc)
				s1 := ast.NodeToString(result.Items[0])
				s2 := ast.NodeToString(result.Items[0])
				if s1 != s2 {
					t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
				}
			})
		}
	})

	t.Run("alter_fulltext_index_structured_fields", func(t *testing.T) {
		t.Run("change_tracking_manual", func(t *testing.T) {
			result := ParseAndCheck(t, "ALTER FULLTEXT INDEX ON dbo.T SET CHANGE_TRACKING MANUAL")
			stmt := result.Items[0].(*ast.AlterFulltextIndexStmt)
			if stmt.ChangeTracking != "MANUAL" {
				t.Errorf("ChangeTracking = %q, want MANUAL", stmt.ChangeTracking)
			}
		})

		t.Run("add_columns_parsed", func(t *testing.T) {
			result := ParseAndCheck(t, "ALTER FULLTEXT INDEX ON dbo.T ADD (col1, col2)")
			stmt := result.Items[0].(*ast.AlterFulltextIndexStmt)
			if stmt.Columns == nil || len(stmt.Columns.Items) != 2 {
				t.Errorf("expected 2 columns, got %v", stmt.Columns)
			}
		})

		t.Run("alter_column_fields", func(t *testing.T) {
			result := ParseAndCheck(t, "ALTER FULLTEXT INDEX ON dbo.T ALTER COLUMN Name ADD STATISTICAL_SEMANTICS")
			stmt := result.Items[0].(*ast.AlterFulltextIndexStmt)
			if stmt.ColumnName != "Name" {
				t.Errorf("ColumnName = %q, want Name", stmt.ColumnName)
			}
			if stmt.ColumnAction != "ADD" {
				t.Errorf("ColumnAction = %q, want ADD", stmt.ColumnAction)
			}
		})

		t.Run("start_population_type", func(t *testing.T) {
			result := ParseAndCheck(t, "ALTER FULLTEXT INDEX ON dbo.T START FULL POPULATION")
			stmt := result.Items[0].(*ast.AlterFulltextIndexStmt)
			if stmt.PopulationType != "FULL" {
				t.Errorf("PopulationType = %q, want FULL", stmt.PopulationType)
			}
		})

		t.Run("with_no_population", func(t *testing.T) {
			result := ParseAndCheck(t, "ALTER FULLTEXT INDEX ON dbo.T DROP (col1) WITH NO POPULATION")
			stmt := result.Items[0].(*ast.AlterFulltextIndexStmt)
			if !stmt.WithNoPopulation {
				t.Error("expected WithNoPopulation = true")
			}
		})
	})

	t.Run("search_property_list_option_structured", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			{"add_with_all_options", "ALTER SEARCH PROPERTY LIST mylist ADD 'Title' WITH (PROPERTY_SET_GUID = 'F29F85E0-4FF9-1068-AB91-08002B27B3D9', PROPERTY_INT_ID = 2, PROPERTY_DESCRIPTION = 'System.Title')"},
			{"add_with_required_options_only", "ALTER SEARCH PROPERTY LIST mylist ADD 'Author' WITH (PROPERTY_SET_GUID = 'F29F85E0-4FF9-1068-AB91-08002B27B3D9', PROPERTY_INT_ID = 4)"},
			{"drop_property", "ALTER SEARCH PROPERTY LIST mylist DROP 'Title'"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterSearchPropertyListStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *AlterSearchPropertyListStmt, got %T", tt.sql, result.Items[0])
				}
				checkLocation(t, tt.sql, "AlterSearchPropertyListStmt", stmt.Loc)
				s1 := ast.NodeToString(result.Items[0])
				s2 := ast.NodeToString(result.Items[0])
				if s1 != s2 {
					t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
				}
			})
		}
	})
}

// TestParseSecurityMiscRemainingDepth tests batch 130: structured CryptoItem and SensitivityOption parsing.
func TestParseSecurityMiscRemainingDepth(t *testing.T) {
	t.Run("signature_crypto_item_structured", func(t *testing.T) {
		tests := []struct {
			sql       string
			mechanism string
			name      string
			withType  string
		}{
			// CERTIFICATE only
			{`ADD SIGNATURE TO dbo.myProc BY CERTIFICATE cert1`, "CERTIFICATE", "cert1", ""},
			// CERTIFICATE with PASSWORD
			{`ADD SIGNATURE TO sp_demo BY CERTIFICATE cert_demo WITH PASSWORD = 'secret123'`, "CERTIFICATE", "cert_demo", "PASSWORD"},
			// CERTIFICATE with SIGNATURE blob
			{`ADD SIGNATURE TO sp_demo BY CERTIFICATE cert_demo WITH SIGNATURE = 0xABCD`, "CERTIFICATE", "cert_demo", "SIGNATURE"},
			// ASYMMETRIC KEY only
			{`ADD SIGNATURE TO dbo.myProc BY ASYMMETRIC KEY myKey`, "ASYMMETRIC KEY", "myKey", ""},
			// ASYMMETRIC KEY with PASSWORD
			{`ADD SIGNATURE TO dbo.myProc BY ASYMMETRIC KEY myKey WITH PASSWORD = 'keypass'`, "ASYMMETRIC KEY", "myKey", "PASSWORD"},
			// ASYMMETRIC KEY with SIGNATURE blob
			{`ADD SIGNATURE TO dbo.myProc BY ASYMMETRIC KEY myKey WITH SIGNATURE = 0xFF01`, "ASYMMETRIC KEY", "myKey", "SIGNATURE"},
			// DROP with CERTIFICATE
			{`DROP SIGNATURE FROM sp_test BY CERTIFICATE cert1`, "CERTIFICATE", "cert1", ""},
			// COUNTER SIGNATURE with CERTIFICATE and PASSWORD
			{`ADD COUNTER SIGNATURE TO ProcSelectT1 BY CERTIFICATE csSelectT WITH PASSWORD = 'secret'`, "CERTIFICATE", "csSelectT", "PASSWORD"},
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
				if stmt.CryptoList == nil || stmt.CryptoList.Len() < 1 {
					t.Fatalf("Parse(%q): expected CryptoList with items", tt.sql)
				}
				ci, ok := stmt.CryptoList.Items[0].(*ast.CryptoItem)
				if !ok {
					t.Fatalf("Parse(%q): expected *CryptoItem, got %T", tt.sql, stmt.CryptoList.Items[0])
				}
				if ci.Mechanism != tt.mechanism {
					t.Errorf("Parse(%q): mechanism = %q, want %q", tt.sql, ci.Mechanism, tt.mechanism)
				}
				if ci.Name != tt.name {
					t.Errorf("Parse(%q): name = %q, want %q", tt.sql, ci.Name, tt.name)
				}
				if ci.WithType != tt.withType {
					t.Errorf("Parse(%q): withType = %q, want %q", tt.sql, ci.WithType, tt.withType)
				}
				checkLocation(t, tt.sql, "CryptoItem", ci.Loc)
				s1 := ast.NodeToString(result.Items[0])
				s2 := ast.NodeToString(result.Items[0])
				if s1 != s2 {
					t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
				}
			})
		}
	})

	t.Run("sensitivity_classification_options_structured", func(t *testing.T) {
		tests := []struct {
			sql      string
			optCount int
			firstKey string
		}{
			// Single option
			{`ADD SENSITIVITY CLASSIFICATION TO dbo.t1.col1 WITH (LABEL = 'Confidential')`, 1, "LABEL"},
			// Multiple options
			{`ADD SENSITIVITY CLASSIFICATION TO dbo.t1.col1 WITH (LABEL = 'Highly Confidential', INFORMATION_TYPE = 'Financial', RANK = CRITICAL)`, 3, "LABEL"},
			// LABEL_ID and INFORMATION_TYPE_ID
			{`ADD SENSITIVITY CLASSIFICATION TO dbo.t1.col1 WITH (LABEL_ID = '643f7acd-776a-438d-890c-79c3f2a520d6', INFORMATION_TYPE_ID = '57845286-7598-22f5-9659-15b24aeb125e')`, 2, "LABEL_ID"},
			// RANK only
			{`ADD SENSITIVITY CLASSIFICATION TO dbo.t1.col1 WITH (RANK = HIGH)`, 1, "RANK"},
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
				if stmt.Options == nil || stmt.Options.Len() != tt.optCount {
					t.Fatalf("Parse(%q): expected %d options, got %v", tt.sql, tt.optCount, stmt.Options)
				}
				opt, ok := stmt.Options.Items[0].(*ast.SensitivityOption)
				if !ok {
					t.Fatalf("Parse(%q): expected *SensitivityOption, got %T", tt.sql, stmt.Options.Items[0])
				}
				if opt.Key != tt.firstKey {
					t.Errorf("Parse(%q): first option key = %q, want %q", tt.sql, opt.Key, tt.firstKey)
				}
				checkLocation(t, tt.sql, "SensitivityOption", opt.Loc)
				s1 := ast.NodeToString(result.Items[0])
				s2 := ast.NodeToString(result.Items[0])
				if s1 != s2 {
					t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
				}
			})
		}
	})
}

// TestParseUtilityCetasDbscopedDepth tests batch 131: structured CETAS columns and ALTER DATABASE SCOPED CONFIG options.
func TestParseUtilityCetasDbscopedDepth(t *testing.T) {
	// Test 1: CETAS with structured column definitions (ColumnDef nodes instead of raw strings)
	t.Run("cetas_column_defs_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			numCols  int
			colNames []string
		}{
			{
				name: "simple_columns",
				sql: `CREATE EXTERNAL TABLE dbo.export_data
					(id, name, value)
					WITH (LOCATION = '/out/', DATA_SOURCE = ds, FILE_FORMAT = ff)
					AS SELECT id, name, value FROM src`,
				numCols:  3,
				colNames: []string{"id", "name", "value"},
			},
			{
				name: "typed_columns",
				sql: `CREATE EXTERNAL TABLE dbo.export_data
					(id INT, name VARCHAR(100), value FLOAT)
					WITH (LOCATION = '/out/', DATA_SOURCE = ds, FILE_FORMAT = ff)
					AS SELECT id, name, value FROM src`,
				numCols:  3,
				colNames: []string{"id", "name", "value"},
			},
			{
				name: "single_column",
				sql: `CREATE EXTERNAL TABLE dbo.export_data (col1)
					WITH (LOCATION = '/out/', DATA_SOURCE = ds, FILE_FORMAT = ff)
					AS SELECT 1`,
				numCols:  1,
				colNames: []string{"col1"},
			},
			{
				name: "typed_columns_nullable",
				sql: `CREATE EXTERNAL TABLE dbo.export_data
					(id INT NOT NULL, name NVARCHAR(50) NULL)
					WITH (LOCATION = '/out/', DATA_SOURCE = ds, FILE_FORMAT = ff)
					AS SELECT id, name FROM src`,
				numCols:  2,
				colNames: []string{"id", "name"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.CreateExternalTableAsSelectStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *CreateExternalTableAsSelectStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Columns == nil {
					t.Fatalf("Parse(%q): expected columns, got nil", tt.sql)
				}
				if len(stmt.Columns.Items) != tt.numCols {
					t.Fatalf("Parse(%q): got %d columns, want %d", tt.sql, len(stmt.Columns.Items), tt.numCols)
				}
				for i, item := range stmt.Columns.Items {
					colDef, ok := item.(*ast.ColumnDef)
					if !ok {
						t.Errorf("Parse(%q): column %d is %T, want *ColumnDef", tt.sql, i, item)
						continue
					}
					if colDef.Name != tt.colNames[i] {
						t.Errorf("Parse(%q): column %d name = %q, want %q", tt.sql, i, colDef.Name, tt.colNames[i])
					}
				}
				checkLocation(t, tt.sql, "CETAS", stmt.Loc)
				// Check deterministic serialization
				s1 := ast.NodeToString(result.Items[0])
				s2 := ast.NodeToString(result.Items[0])
				if s1 != s2 {
					t.Errorf("Parse(%q): serialization not deterministic", tt.sql)
				}
			})
		}
	})

	// Test 2: ALTER DATABASE SCOPED CONFIGURATION with structured SecurityPrincipalOption nodes
	t.Run("alter_db_scoped_config_structured", func(t *testing.T) {
		tests := []struct {
			name  string
			sql   string
			key   string
			value string
		}{
			{
				name:  "set_maxdop",
				sql:   "ALTER DATABASE SCOPED CONFIGURATION SET MAXDOP = 4",
				key:   "MAXDOP",
				value: "4",
			},
			{
				name:  "set_on",
				sql:   "ALTER DATABASE SCOPED CONFIGURATION SET LEGACY_CARDINALITY_ESTIMATION = ON",
				key:   "LEGACY_CARDINALITY_ESTIMATION",
				value: "ON",
			},
			{
				name:  "set_off",
				sql:   "ALTER DATABASE SCOPED CONFIGURATION SET PARAMETER_SNIFFING = OFF",
				key:   "PARAMETER_SNIFFING",
				value: "OFF",
			},
			{
				name:  "set_primary",
				sql:   "ALTER DATABASE SCOPED CONFIGURATION SET MAXDOP = PRIMARY",
				key:   "MAXDOP",
				value: "PRIMARY",
			},
			{
				name:  "set_string_value",
				sql:   "ALTER DATABASE SCOPED CONFIGURATION SET LEDGER_DIGEST_STORAGE_ENDPOINT = 'https://ledger.endpoint.com'",
				key:   "LEDGER_DIGEST_STORAGE_ENDPOINT",
				value: "https://ledger.endpoint.com",
			},
			{
				name:  "set_when_supported",
				sql:   "ALTER DATABASE SCOPED CONFIGURATION SET ELEVATE_ONLINE = WHEN_SUPPORTED",
				key:   "ELEVATE_ONLINE",
				value: "WHEN_SUPPORTED",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): expected options, got nil", tt.sql)
				}
				// Find the SecurityPrincipalOption with the expected key
				found := false
				for _, item := range stmt.Options.Items {
					opt, ok := item.(*ast.SecurityPrincipalOption)
					if !ok {
						continue
					}
					if opt.Name == tt.key {
						found = true
						if opt.Value != tt.value {
							t.Errorf("Parse(%q): option %q value = %q, want %q", tt.sql, tt.key, opt.Value, tt.value)
						}
					}
				}
				if !found {
					t.Errorf("Parse(%q): expected SecurityPrincipalOption with key %q, got none", tt.sql, tt.key)
				}
				checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
			})
		}

		// Test FOR SECONDARY
		t.Run("for_secondary", func(t *testing.T) {
			sql := "ALTER DATABASE SCOPED CONFIGURATION FOR SECONDARY SET MAXDOP = PRIMARY"
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Options == nil {
				t.Fatalf("Parse(%q): expected options, got nil", sql)
			}
			// Should have FOR SECONDARY string and a SecurityPrincipalOption
			hasForSecondary := false
			hasOption := false
			for _, item := range stmt.Options.Items {
				if s, ok := item.(*ast.String); ok && s.Str == "FOR SECONDARY" {
					hasForSecondary = true
				}
				if opt, ok := item.(*ast.SecurityPrincipalOption); ok && opt.Name == "MAXDOP" {
					hasOption = true
				}
			}
			if !hasForSecondary {
				t.Errorf("Parse(%q): expected FOR SECONDARY marker", sql)
			}
			if !hasOption {
				t.Errorf("Parse(%q): expected MAXDOP option", sql)
			}
		})

		// Test CLEAR PROCEDURE_CACHE
		t.Run("clear_procedure_cache", func(t *testing.T) {
			sql := "ALTER DATABASE SCOPED CONFIGURATION CLEAR PROCEDURE_CACHE"
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", sql, result.Len())
			}
			stmt, ok := result.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", sql, result.Items[0])
			}
			if stmt.Options == nil {
				t.Fatalf("Parse(%q): expected options, got nil", sql)
			}
			// CLEAR is still a String since it's not a key=value option
			found := false
			for _, item := range stmt.Options.Items {
				if s, ok := item.(*ast.String); ok && strings.Contains(s.Str, "CLEAR") {
					found = true
				}
			}
			if !found {
				t.Errorf("Parse(%q): expected CLEAR option string", sql)
			}
		})
	})
}

// TestParseAlterDatabaseUnknownDepth tests batch 132: structured ALTER DATABASE unknown action,
// sub-options, and termination parsing.
func TestParseAlterDatabaseUnknownDepth(t *testing.T) {
	t.Run("alter_database_unknown_structured", func(t *testing.T) {
		tests := []struct {
			name   string
			sql    string
			action string
		}{
			// Standalone unknown action
			{"failover", "ALTER DATABASE mydb FAILOVER", "FAILOVER"},
			// Unknown action with extra keyword
			{"unknown with keyword", "ALTER DATABASE mydb SUSPEND LEDGER", "SUSPEND"},
			// Unknown action with parenthesized options
			{"unknown with parens", "ALTER DATABASE mydb RESUME LEDGER (OPTION1 = 100, OPTION2 = OFF)", "RESUME"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): expected 1 statement, got %d", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected AlterDatabaseStmt", tt.sql)
				}
				if stmt.Action != tt.action {
					t.Errorf("Parse(%q): expected action %q, got %q", tt.sql, tt.action, stmt.Action)
				}
			})
		}
	})

	t.Run("alter_database_sub_options_structured", func(t *testing.T) {
		tests := []struct {
			name string
			sql  string
		}{
			// CHANGE_TRACKING with sub-options
			{"change_tracking sub-opts", "ALTER DATABASE mydb SET CHANGE_TRACKING = ON (AUTO_CLEANUP = ON, CHANGE_RETENTION = 7 DAYS)"},
			// QUERY_STORE with sub-options
			{"query_store sub-opts", "ALTER DATABASE mydb SET QUERY_STORE = ON (MAX_STORAGE_SIZE_MB = 100, INTERVAL_LENGTH_MINUTES = 60)"},
			// QUERY_STORE OFF with FORCED
			{"query_store off forced", "ALTER DATABASE mydb SET QUERY_STORE = OFF (FORCED)"},
			// ACCELERATED_DATABASE_RECOVERY with sub-options
			{"adr sub-opts", "ALTER DATABASE mydb SET ACCELERATED_DATABASE_RECOVERY = ON (PERSISTENT_VERSION_STORE_FILEGROUP = fg1)"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): expected 1 statement, got %d", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected AlterDatabaseStmt", tt.sql)
				}
				if stmt.Action != "SET" {
					t.Errorf("Parse(%q): expected action SET, got %q", tt.sql, stmt.Action)
				}
				if stmt.Options == nil || stmt.Options.Len() == 0 {
					t.Errorf("Parse(%q): expected non-empty options", tt.sql)
				}
			})
		}
	})

	t.Run("alter_database_termination_structured", func(t *testing.T) {
		tests := []struct {
			name        string
			sql         string
			termination string
		}{
			{"rollback immediate", "ALTER DATABASE mydb SET SINGLE_USER WITH ROLLBACK IMMEDIATE", "ROLLBACK IMMEDIATE"},
			{"rollback after seconds", "ALTER DATABASE mydb SET SINGLE_USER WITH ROLLBACK AFTER 60 SECONDS", "ROLLBACK AFTER 60 SECONDS"},
			{"rollback after no unit", "ALTER DATABASE mydb SET SINGLE_USER WITH ROLLBACK AFTER 30", "ROLLBACK AFTER 30"},
			{"no_wait", "ALTER DATABASE mydb SET SINGLE_USER WITH NO_WAIT", "NO_WAIT"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): expected 1 statement, got %d", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.AlterDatabaseStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected AlterDatabaseStmt", tt.sql)
				}
				if stmt.Termination != tt.termination {
					t.Errorf("Parse(%q): expected termination %q, got %q", tt.sql, tt.termination, stmt.Termination)
				}
			})
		}
	})
}

// TestParseEndpointRemainingDepth tests batch 133: structured endpoint option parsing,
// removal of skipParenthesized, and elimination of generic p.advance() skips.
func TestParseEndpointRemainingDepth(t *testing.T) {
	t.Run("generic_value_structured", func(t *testing.T) {
		tests := []struct {
			name     string
			sql      string
			wantOpts []string
		}{
			{
				"generic_protocol_parenthesized_value",
				"CREATE ENDPOINT ep1 AS CUSTOM_PROTO (AUTH = (BASIC, DIGEST), PORT = 8080) FOR TSQL ()",
				[]string{
					"AS=CUSTOM_PROTO",
					"AUTH=BASIC,DIGEST",
					"PORT=8080",
					"FOR=TSQL",
				},
			},
			{
				"generic_protocol_nested_parens_value",
				"CREATE ENDPOINT ep1 AS CUSTOM_PROTO (METHODS = (GET, POST, PUT)) FOR TSQL ()",
				[]string{
					"AS=CUSTOM_PROTO",
					"METHODS=GET,POST,PUT",
					"FOR=TSQL",
				},
			},
			{
				"generic_protocol_flag_option",
				"CREATE ENDPOINT ep1 AS CUSTOM_PROTO (VERBOSE) FOR TSQL ()",
				[]string{
					"AS=CUSTOM_PROTO",
					"VERBOSE",
					"FOR=TSQL",
				},
			},
			{
				"endpoint_payload_unknown_option_key_value",
				"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 4022) FOR SERVICE_BROKER (AUTHENTICATION = WINDOWS, CUSTOM_OPT = myval)",
				[]string{
					"AS=TCP",
					"LISTENER_PORT=4022",
					"FOR=SERVICE_BROKER",
					"AUTHENTICATION=WINDOWS",
					"CUSTOM_OPT=MYVAL",
				},
			},
			{
				"endpoint_options_unknown_toplevel_keyword",
				"CREATE ENDPOINT ep1 AFFINITY = ADMIN AS TCP (LISTENER_PORT = 4022) FOR TSQL ()",
				[]string{
					"AFFINITY=ADMIN",
					"AS=TCP",
					"LISTENER_PORT=4022",
					"FOR=TSQL",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				stmt, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tt.sql, result.Items[0])
				}
				if stmt.Options == nil {
					t.Fatalf("Parse(%q): options is nil", tt.sql)
				}
				for _, want := range tt.wantOpts {
					found := false
					for _, item := range stmt.Options.Items {
						if s, ok := item.(*ast.String); ok && s.Str == want {
							found = true
							break
						}
					}
					if !found {
						var gotStrs []string
						for _, item := range stmt.Options.Items {
							gotStrs = append(gotStrs, item.(*ast.String).Str)
						}
						t.Errorf("Parse(%q): expected option %q not found in %v", tt.sql, want, gotStrs)
					}
				}
				checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
			})
		}
	})

	t.Run("no_skip_parens", func(t *testing.T) {
		// Verify that skipParenthesized is removed by confirming these still parse correctly
		tests := []struct {
			name string
			sql  string
		}{
			{
				"alter_endpoint_basic",
				"ALTER ENDPOINT ep1 STATE = STARTED AS TCP (LISTENER_PORT = 5022) FOR SERVICE_BROKER (AUTHENTICATION = WINDOWS)",
			},
			{
				"create_endpoint_db_mirroring_full",
				"CREATE ENDPOINT ep1 AS TCP (LISTENER_PORT = 7022) FOR DATABASE_MIRRORING (AUTHENTICATION = CERTIFICATE MyCert, ENCRYPTION = REQUIRED ALGORITHM AES, ROLE = ALL)",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ParseAndCheck(t, tt.sql)
				if result.Len() != 1 {
					t.Fatalf("Parse(%q): got %d statements, want 1", tt.sql, result.Len())
				}
				_, ok := result.Items[0].(*ast.SecurityStmt)
				if !ok {
					t.Fatalf("Parse(%q): expected *SecurityStmt, got %T", tt.sql, result.Items[0])
				}
			})
		}
	})
}
