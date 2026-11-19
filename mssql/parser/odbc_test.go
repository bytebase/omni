package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestOdbcEscapes pins the ODBC escapes of a scalar expression, as the pinned
// SQL Server answered them under SET PARSEONLY ON: the rows of testdata/
// compat.tsv under odbcFunctionCall, odbcLiteral, expressionPrimary/
// odbcFunctionCall, likePredicate and suite:lexer.
func TestOdbcEscapes(t *testing.T) {
	for _, sql := range []string{
		"SELECT { FN CURRENT_TIME }",
		"SELECT { FN CURRENT_TIMESTAMP }",
		"SELECT { FN CURRENT_TIME() }",
		"SELECT { FN CURRENT_TIME(1) }",
		"SELECT { FN CURRENT_TIME(zzz) }",
		"SELECT { FN DATABASE() }",
		"SELECT { FN USER() }",
		"SELECT { FN CURRENT_DATE() }",
		"SELECT { FN TRUNCATE(1, 1) }",
		"SELECT { FN LEFT(1, 2) }",
		"SELECT { FN RIGHT(1, 2) }",
		"SELECT { FN INSERT(1, 2, 3, 4) }",
		"SELECT { FN CONVERT(1, SQL_VARCHAR) }",
		"SELECT { FN EXTRACT(year FROM 1) }",
		"SELECT { FN TIMESTAMPADD(SQL_TSI_DAY, 1, a) }",
		"SELECT { FN TIMESTAMPDIFF(SQL_TSI_DAY, a, b) }",
		"SELECT {fn NOW()}",
		"SELECT {fn UCASE(a)}",
		"SELECT { FN CURRENT_TIME } zzz",
		"SELECT {d '2020-01-01'}",
		"SELECT {t '12:00:00'}",
		"SELECT {ts '2020-01-01 12:00:00'}",
		"SELECT {D '2020-01-01'}",
		"IF 1 LIKE 1 { ESCAPE 1 } SELECT 1",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// An escape names a function SQL Server maps to a T-SQL one.
		"SELECT { FN FOO() }",
		"SELECT { FN COOKIE(1, a) }",
		"SELECT { FN NOW(1) }",
		// Each named form keeps its argument count.
		"SELECT { FN CURRENT_DATE(1) }",
		"SELECT { FN CURRENT_TIME(1, 2) }",
		"SELECT { FN DATABASE(1) }",
		"SELECT { FN TRUNCATE(1) }",
		"SELECT { FN TRUNCATE(1, 2, 3) }",
		"SELECT { FN LEFT(1) }",
		"SELECT { FN INSERT(1, 2) }",
		"SELECT { FN CURRENT_TIME(1) zzz }",
		// CONVERT takes an ODBC type name, EXTRACT an ODBC part, and
		// TIMESTAMPADD an SQL_TSI interval.
		"SELECT { FN CONVERT(1, a) }",
		"SELECT { FN CONVERT(1, INT) }",
		"SELECT { FN CONVERT(1, varchar) }",
		"SELECT { FN CONVERT(1, SQL_VARCHAR, 3) }",
		"SELECT { FN EXTRACT(a FROM 1) }",
		"SELECT { FN EXTRACT(quarter FROM 1) }",
		"SELECT { FN TIMESTAMPADD(day, 1, a) }",
		// The literal escape takes d, t and ts.
		"SELECT {x 'a'}",
		"SELECT {dd '2020-01-01'}",
		// An argument is required where the form has one.
		"SELECT {fn CONVERT(, SQL_VARCHAR)}",
		"SELECT {fn EXTRACT(year FROM )}",
		// FN and the closing brace are required.
		"SELECT { CURRENT_TIME }",
		"SELECT { FN CURRENT_TIME",
		"SELECT FN CURRENT_TIME }",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want SQL Server's rejection", sql)
		}
	}
}

// TestOdbcEscapeAST pins the tree and the locations of the escapes.
func TestOdbcEscapeAST(t *testing.T) {
	const sql = "SELECT { FN CONVERT(1, SQL_VARCHAR) }"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	item := list.Items[0].(*ast.SelectStmt).TargetList.Items[0].(*ast.ResTarget)
	call, ok := item.Val.(*ast.OdbcFunctionCallExpr)
	if !ok {
		t.Fatalf("target = %T, want *OdbcFunctionCallExpr", item.Val)
	}
	if call.Name != "CONVERT" || call.Convert != "SQL_VARCHAR" || call.Args == nil || call.Args.Len() != 1 {
		t.Errorf("call = %+v, want CONVERT to SQL_VARCHAR with one argument", call)
	}
	if got := sql[call.Loc.Start:call.Loc.End]; got != "{ FN CONVERT(1, SQL_VARCHAR) }" {
		t.Errorf("Loc covers %q, want the whole escape", got)
	}

	list, err = Parse("SELECT {d '2020-01-01'}")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	item = list.Items[0].(*ast.SelectStmt).TargetList.Items[0].(*ast.ResTarget)
	lit, ok := item.Val.(*ast.OdbcLiteralExpr)
	if !ok {
		t.Fatalf("target = %T, want *OdbcLiteralExpr", item.Val)
	}
	if lit.Kind != "d" || lit.Value.Str != "2020-01-01" {
		t.Errorf("literal = %+v, want {d '2020-01-01'}", lit)
	}

	list, err = Parse("SELECT a FROM t WHERE a LIKE 'x' { ESCAPE '!' }")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	like, ok := list.Items[0].(*ast.SelectStmt).WhereClause.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("WHERE = %T, want *LikeExpr", list.Items[0].(*ast.SelectStmt).WhereClause)
	}
	if !like.EscapeBraced || like.Escape == nil {
		t.Errorf("LikeExpr = %+v, want a braced escape", like)
	}
}
