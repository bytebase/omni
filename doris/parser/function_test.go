package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

func TestCreateFunctionAST(t *testing.T) {
	stmt := parseOne(t, "CREATE GLOBAL AGGREGATE FUNCTION IF NOT EXISTS db.f(INT, VARCHAR(*), ...) RETURNS BIGINT INTERMEDIATE STRING PROPERTIES('symbol' = 's')")
	cf, ok := stmt.(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("got %T", stmt)
	}
	if !cf.Global || !cf.Aggregate || !cf.IfNotExists || !cf.Variadic || len(cf.Args) != 2 ||
		cf.Returns.Name != "BIGINT" || cf.Intermediate.Name != "STRING" || len(cf.Properties) != 1 ||
		cf.Name.String() != "db.f" {
		t.Errorf("CreateFunctionStmt = %+v", cf)
	}

	stmt = parseOne(t, "CREATE ALIAS FUNCTION f(INT, INT) WITH PARAMETER(a, b) AS a + b")
	cf = stmt.(*ast.CreateFunctionStmt)
	if !cf.Alias || len(cf.Parameters) != 2 || cf.Body == nil || cf.Returns != nil {
		t.Errorf("alias CreateFunctionStmt = %+v", cf)
	}
}

func TestDropFunctionAST(t *testing.T) {
	stmt := parseOne(t, "DROP SESSION FUNCTION IF EXISTS f(ARRAY<INT>, ...)")
	df, ok := stmt.(*ast.DropFunctionStmt)
	if !ok {
		t.Fatalf("got %T", stmt)
	}
	if df.Global || !df.IfExists || !df.Variadic || len(df.Args) != 1 || df.Args[0].Name != "ARRAY" {
		t.Errorf("DropFunctionStmt = %+v", df)
	}
}

// TestFunctionSignatureGrammar holds the function statements to what
// Doris's two parsers both take (testdata/compat.tsv judged each).
func TestFunctionSignatureGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"DROP FUNCTION a(TINYINT, ...)", true},
		{"DROP FUNCTION a(..., ...)", false},
		{"DROP FUNCTION a(..., TINYINT)", false},
		{"DROP FUNCTION a(, TINYINT)", false},
		{"DROP LOCAL FUNCTION a()", true},
		{"DROP FUNCTION a", false},
		{"CREATE FUNCTION LEFT(INT) RETURNS INT", true},
		{"DROP FUNCTION db.RIGHT(INT)", true},
		{"SHOW CREATE FUNCTION LEFT(INT)", true},
		{"DROP FUNCTION TRIM(INT)", false},
		{"CREATE TABLES FUNCTION a() RETURNS TINYINT", true},
		{"CREATE SESSION FUNCTION a(...) RETURNS VARIANT", true},
		{"CREATE FUNCTION a() RETURNS TINYINT INTERMEDIATE", false},
		{"CREATE FUNCTION a() RETURNS TINYINT PROPERTIES(a = a)", false},
		{"CREATE FUNCTION a()", false},
		{"CREATE ALIAS FUNCTION a() WITH PARAMETER() AS 1", true},
		{"CREATE ALIAS FUNCTION a() WITH PARAMETER(a,) AS 1", false},
		{"CREATE ALIAS FUNCTION a() WITH PARAMETER() AS", false},
		{"SHOW CREATE FUNCTION a()", true},
		{"SHOW CREATE GLOBAL FUNCTION a(INT, ...) FROM db", true},
		{"SHOW CREATE FUNCTION a(INT", false},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}
