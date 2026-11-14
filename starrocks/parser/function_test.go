package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

func TestCreateFunctionAST(t *testing.T) {
	stmt := parseOne(t, "CREATE OR REPLACE GLOBAL AGGREGATE FUNCTION IF NOT EXISTS db.f(INT, ...) RETURNS BIGINT PROPERTIES('symbol' = 's')")
	cf, ok := stmt.(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("got %T", stmt)
	}
	if !cf.OrReplace || !cf.Global || !cf.Aggregate || !cf.IfNotExists || !cf.Variadic ||
		len(cf.Args) != 1 || cf.Returns.Name != "BIGINT" || len(cf.Properties) != 1 {
		t.Errorf("CreateFunctionStmt = %+v", cf)
	}
	stmt = parseOne(t, "CREATE TABLE FUNCTION f() RETURNS INT symbol = 's' type = 'StarrocksJar'")
	cf = stmt.(*ast.CreateFunctionStmt)
	if !cf.Table || len(cf.Properties) != 2 || cf.Properties[1].Key != "type" {
		t.Errorf("table CreateFunctionStmt = %+v", cf)
	}
}

func TestDropFunctionAST(t *testing.T) {
	stmt := parseOne(t, "DROP GLOBAL FUNCTION IF EXISTS f(, INT, ...)")
	df, ok := stmt.(*ast.DropFunctionStmt)
	if !ok {
		t.Fatalf("got %T", stmt)
	}
	if !df.Global || !df.IfExists || !df.Variadic || len(df.Args) != 1 {
		t.Errorf("DropFunctionStmt = %+v", df)
	}
}

// TestFunctionSignatureGrammar holds the function statements to
// StarRocks.g4 (testdata/compat.tsv judged each).
func TestFunctionSignatureGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"DROP FUNCTION a(, BOOLEAN)", true},
		{"DROP FUNCTION a(,...)", true},
		{"DROP FUNCTION a(...)", false},
		{"DROP FUNCTION a(,)", false},
		{"DROP FUNCTION a(zzz)", false},
		{"DROP FUNCTION a(DECIMAL64(10, 2), NUMERIC, NUMBER(3), DECIMAL128, BINARY)", true},
		{"CREATE TABLE t (a NUMERIC(10, 2), b DECIMAL32(5))", true},
		{"SELECT CAST(1 AS DECIMAL64(10, 2))", true},
		{"CREATE TABLE FUNCTION a() RETURNS BOOLEAN", true},
		{"CREATE OR REPLACE FUNCTION a() RETURNS BOOLEAN", true},
		{"CREATE FUNCTION a() RETURNS BOOLEAN AS $$a$$", false},
		{"CREATE FUNCTION a() RETURNS BOOLEAN a = 'a' zzz", false},
		{"CREATE FUNCTION a() RETURNS", false},
		{"CREATE ALIAS FUNCTION a() WITH PARAMETER() AS 1", false},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}
