package parser

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

func TestMSSQLLocPrecision_SpatialIndexNestedOptions(t *testing.T) {
	sql := `CREATE SPATIAL INDEX SI_t ON dbo.t (geom) USING GEOMETRY_GRID WITH (BOUNDING_BOX = (0, 0, 100, 100), GRIDS = (LEVEL_1 = LOW, LEVEL_2 = MEDIUM, LEVEL_3 = HIGH, LEVEL_4 = HIGH), CELLS_PER_OBJECT = 16)`

	result := ParseAndCheck(t, sql)
	stmt, ok := result.Items[0].(*ast.CreateSpatialIndexStmt)
	if !ok {
		t.Fatalf("stmt type = %T, want *ast.CreateSpatialIndexStmt", result.Items[0])
	}
	if stmt.Options == nil {
		t.Fatal("expected spatial index options")
	}

	var got []string
	for _, item := range stmt.Options.Items {
		opt, ok := item.(*ast.String)
		if !ok {
			t.Fatalf("option type = %T, want *ast.String", item)
		}
		got = append(got, opt.Str)
	}

	want := []string{
		"BOUNDING_BOX=(0, 0, 100, 100)",
		"GRIDS=(LEVEL_1 = LOW, LEVEL_2 = MEDIUM, LEVEL_3 = HIGH, LEVEL_4 = HIGH)",
		"CELLS_PER_OBJECT=16",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("options = %#v, want %#v", got, want)
	}
}

func TestMSSQLLocPrecision_ExpressionRanges(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{sql: "SELECT c > 0", want: "c > 0"},
		{sql: "SELECT c IN (1, 2)", want: "c IN (1, 2)"},
		{sql: "SELECT c BETWEEN 1 AND 2", want: "c BETWEEN 1 AND 2"},
		{sql: "SELECT c LIKE 'x%'", want: "c LIKE 'x%'"},
		{sql: "SELECT c IS NOT NULL", want: "c IS NOT NULL"},
		{sql: "SELECT a + b * c", want: "a + b * c"},
		{sql: "SELECT (a = 1) AND (b = 2)", want: "(a = 1) AND (b = 2)"},
		{
			sql:  "SELECT name COLLATE Latin1_General_CI_AS LIKE 'A%'",
			want: "name COLLATE Latin1_General_CI_AS LIKE 'A%'",
		},
		{sql: "SELECT dt AT TIME ZONE 'UTC'", want: "dt AT TIME ZONE 'UTC'"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			expr := firstSelectTargetExpr(t, tt.sql)
			if got := locText(t, tt.sql, ast.NodeLoc(expr)); got != tt.want {
				t.Fatalf("loc text = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMSSQLLocPrecision_CheckConstraintExpressionRanges(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE TABLE t (c int, CONSTRAINT ck CHECK (c > 0))",
			want: "c > 0",
		},
		{
			sql:  "CREATE TABLE t (c int CHECK (c IN (1, 2)))",
			want: "c IN (1, 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			var checkExpr ast.ExprNode
			ast.Inspect(ParseAndCheck(t, tt.sql), func(n ast.Node) bool {
				constraint, ok := n.(*ast.ConstraintDef)
				if ok && constraint.Type == ast.ConstraintCheck && checkExpr == nil {
					checkExpr = constraint.Expr
					return false
				}
				return true
			})
			if checkExpr == nil {
				t.Fatal("expected CHECK constraint expression")
			}
			if got := locText(t, tt.sql, ast.NodeLoc(checkExpr)); got != tt.want {
				t.Fatalf("loc text = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMSSQLLocPrecision_FuncCallRanges(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{sql: "SELECT SYSDATETIME()", want: "SYSDATETIME()"},
		{sql: "SELECT COUNT(*)", want: "COUNT(*)"},
		{
			sql:  "SELECT STRING_AGG(name, ',') WITHIN GROUP (ORDER BY name)",
			want: "STRING_AGG(name, ',') WITHIN GROUP (ORDER BY name)",
		},
		{
			sql:  "SELECT SUM(x) OVER (PARTITION BY y)",
			want: "SUM(x) OVER (PARTITION BY y)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			expr := firstSelectTargetExpr(t, tt.sql)
			if got := locText(t, tt.sql, ast.NodeLoc(expr)); got != tt.want {
				t.Fatalf("loc text = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMSSQLLocPrecision_DefaultFuncCallRanges(t *testing.T) {
	sql := "CREATE TABLE t (d datetime2 DEFAULT SYSDATETIME(), e int DEFAULT dbo.f(), c int DEFAULT COUNT(*))"

	var got []string
	ast.Inspect(ParseAndCheck(t, sql), func(n ast.Node) bool {
		column, ok := n.(*ast.ColumnDef)
		if !ok {
			return true
		}
		if fc, ok := column.DefaultExpr.(*ast.FuncCallExpr); ok {
			got = append(got, locText(t, sql, fc.Loc))
		}
		return true
	})

	want := []string{"SYSDATETIME()", "dbo.f()", "COUNT(*)"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("function loc texts = %#v, want %#v", got, want)
	}
}

func TestMSSQLLocPrecision_ColumnRefRanges(t *testing.T) {
	sql := "SELECT c, t.c, dbo.t.c, t.*"

	result := ParseAndCheck(t, sql)
	stmt := result.Items[0].(*ast.SelectStmt)
	var got []string
	for _, item := range stmt.TargetList.Items {
		target := item.(*ast.ResTarget)
		got = append(got, locText(t, sql, ast.NodeLoc(target.Val)))
	}

	want := []string{"c", "t.c", "dbo.t.c", "t.*"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("column loc texts = %#v, want %#v", got, want)
	}
}

func TestMSSQLLocPrecision_ExpressionNodeLocCompleteness(t *testing.T) {
	tests := []string{
		"SELECT (a + b) * dbo.f(c) FROM t WHERE c IS NOT NULL",
		"CREATE TABLE t (c int CHECK (CASE WHEN c BETWEEN 1 AND 2 THEN dbo.f(c) ELSE SYSDATETIME() END IS NOT NULL))",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ast.Inspect(ParseAndCheck(t, sql), func(n ast.Node) bool {
				if !locPrecisionExprNode(n) {
					return true
				}
				loc := ast.NodeLoc(n)
				if !loc.IsValid() {
					t.Errorf("%T has invalid loc: %+v", n, loc)
					return true
				}
				if loc.Start > len(sql) || loc.End > len(sql) || loc.Start > loc.End {
					t.Errorf("%T has out-of-range loc: %+v", n, loc)
				}
				return true
			})
		})
	}
}

func locPrecisionExprNode(n ast.Node) bool {
	switch n.(type) {
	case *ast.AtTimeZoneExpr,
		*ast.BetweenExpr,
		*ast.BinaryExpr,
		*ast.CaseExpr,
		*ast.CaseWhen,
		*ast.ColumnRef,
		*ast.CollateExpr,
		*ast.FuncCallExpr,
		*ast.InExpr,
		*ast.IsExpr,
		*ast.LikeExpr,
		*ast.Literal,
		*ast.ParenExpr,
		*ast.StarExpr,
		*ast.UnaryExpr,
		*ast.VariableRef:
		return true
	default:
		return false
	}
}

func firstSelectTargetExpr(t *testing.T, sql string) ast.ExprNode {
	t.Helper()

	result := ParseAndCheck(t, sql)
	stmt, ok := result.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("stmt type = %T, want *ast.SelectStmt", result.Items[0])
	}
	if stmt.TargetList == nil || len(stmt.TargetList.Items) == 0 {
		t.Fatal("expected SELECT target")
	}
	target, ok := stmt.TargetList.Items[0].(*ast.ResTarget)
	if !ok {
		t.Fatalf("target type = %T, want *ast.ResTarget", stmt.TargetList.Items[0])
	}
	return target.Val
}

func locText(t *testing.T, sql string, loc ast.Loc) string {
	t.Helper()

	if !loc.IsValid() {
		t.Fatalf("invalid loc: %+v", loc)
	}
	if loc.Start > len(sql) || loc.End > len(sql) || loc.Start > loc.End {
		t.Fatalf("loc out of range: %+v for %q", loc, sql)
	}
	return sql[loc.Start:loc.End]
}
