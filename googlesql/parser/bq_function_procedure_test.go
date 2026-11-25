package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl-bigquery node: CREATE [AGGREGATE] FUNCTION /
// TABLE FUNCTION / PROCEDURE and their DROP forms. These are BigQuery-only at the
// GoogleSQL union level (the Spanner emulator rejects most of them), so the
// accept/reject verdicts here are triangulated against the legacy
// GoogleSQLParser.g4 + the BigQuery truth1 corpus (DDL-015..020/049/050/051), not
// the Spanner differential. Structural-gate assertions on the produced AST.

func cfOf(t *testing.T, sql string) *ast.CreateFunctionStmt {
	t.Helper()
	n := parseDDL(t, sql)
	cf, ok := n.(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): statement is %T, want *ast.CreateFunctionStmt", sql, n)
	}
	return cf
}

func TestCreateFunction_SQLBody(t *testing.T) {
	// DDL-015: CREATE FUNCTION mydataset.square(x FLOAT64) RETURNS FLOAT64 AS (x * x)
	cf := cfOf(t, "CREATE FUNCTION mydataset.square(x FLOAT64) RETURNS FLOAT64 AS (x * x)")
	if cf.Name.String() != "mydataset.square" {
		t.Errorf("Name = %q, want mydataset.square", cf.Name.String())
	}
	if len(cf.Params) != 1 || cf.Params[0].Name != "x" || cf.Params[0].Type.Text != "FLOAT64" {
		t.Fatalf("Params = %+v, want [x FLOAT64]", cf.Params)
	}
	if cf.Returns == nil || cf.Returns.Text != "FLOAT64" {
		t.Errorf("Returns = %v, want FLOAT64", cf.Returns)
	}
	if cf.Body == nil {
		t.Error("Body = nil, want the (x * x) expression")
	}
	if cf.Aggregate || cf.IsTableFunc {
		t.Errorf("Aggregate=%v IsTableFunc=%v, want both false", cf.Aggregate, cf.IsTableFunc)
	}
}

func TestCreateFunction_Rejects(t *testing.T) {
	cases := []string{
		"CREATE FUNCTION ds.f",                        // missing parameter list
		"CREATE FUNCTION ds.f(x INT64) RETURNS",       // RETURNS with no type
		"CREATE FUNCTION ds.f(x INT64) AS (SELECT 1)", // sql_function_body error: query body must be a scalar subquery
		"CREATE FUNCTION ds.f(x INT64) AS",            // AS with no body
		"CREATE AGGREGATE ds.f(x INT64) AS (x)",       // AGGREGATE without FUNCTION
		"CREATE FUNCTION (x INT64) AS (x)",            // missing name
		"CREATE TABLE FUNCTION ds.f AS SELECT 1",      // TABLE FUNCTION missing parameter list
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- CREATE PROCEDURE ---

func TestCreateProcedure_Rejects(t *testing.T) {
	cases := []string{
		"CREATE PROCEDURE ds.p",                        // missing parameter list
		"CREATE PROCEDURE ds.p()",                      // missing body
		"CREATE PROCEDURE ds.p() BEGIN SELECT 1;",      // unterminated BEGIN…END (no END)
		"CREATE PROCEDURE ds.p(x) BEGIN SELECT 1; END", // procedure parameter without a type
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- DROP FUNCTION / TABLE FUNCTION / PROCEDURE ---

func bqDropOf(t *testing.T, sql string) *ast.BQDropStmt {
	t.Helper()
	n := parseDDL(t, sql)
	d, ok := n.(*ast.BQDropStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): statement is %T, want *ast.BQDropStmt", sql, n)
	}
	return d
}

func TestDropFunction(t *testing.T) {
	// DDL-049.
	d := bqDropOf(t, "DROP FUNCTION IF EXISTS mydataset.my_function")
	if d.Object != ast.BQDropFunction {
		t.Errorf("Object = %v, want FUNCTION", d.Object)
	}
	if !d.IfExists {
		t.Error("IfExists = false, want true")
	}
	if d.Name.String() != "mydataset.my_function" {
		t.Errorf("Name = %q", d.Name.String())
	}
}

func TestDropFunction_WithOverloadParams(t *testing.T) {
	// opt_function_parameters disambiguating an overload — consumed structurally.
	d := bqDropOf(t, "DROP FUNCTION ds.f(INT64, STRING)")
	if d.Object != ast.BQDropFunction || d.Name.String() != "ds.f" {
		t.Errorf("got Object=%v Name=%q", d.Object, d.Name.String())
	}
}
