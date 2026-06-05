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
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateFunctionStmt", sql, n)
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

func TestCreateFunction_OrReplaceTempNoReturns(t *testing.T) {
	// DDL-015: CREATE OR REPLACE TEMP FUNCTION add_n(x INT64, n INT64) AS (x + n)
	cf := cfOf(t, "CREATE OR REPLACE TEMP FUNCTION add_n(x INT64, n INT64) AS (x + n)")
	if !cf.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if cf.Scope != "TEMP" {
		t.Errorf("Scope = %q, want TEMP", cf.Scope)
	}
	if len(cf.Params) != 2 {
		t.Fatalf("Params = %d, want 2", len(cf.Params))
	}
	if cf.Returns != nil {
		t.Errorf("Returns = %v, want nil (no RETURNS)", cf.Returns)
	}
}

func TestCreateFunction_JavaScript(t *testing.T) {
	// DDL-016: LANGUAGE js with a raw-string body and DETERMINISTIC.
	cf := cfOf(t, `CREATE OR REPLACE FUNCTION mydataset.multiplyInputs(x FLOAT64, y FLOAT64) RETURNS FLOAT64 LANGUAGE js AS "return x*y;"`)
	if cf.Language != "js" {
		t.Errorf("Language = %q, want js", cf.Language)
	}
	if !cf.HasBodyString || cf.BodyString != "return x*y;" {
		t.Errorf("BodyString = %q (has=%v), want \"return x*y;\"", cf.BodyString, cf.HasBodyString)
	}
}

func TestCreateFunction_Determinism(t *testing.T) {
	cf := cfOf(t, `CREATE FUNCTION ds.f(x INT64) RETURNS INT64 DETERMINISTIC LANGUAGE js AS "return x"`)
	if cf.Determinism != "DETERMINISTIC" {
		t.Errorf("Determinism = %q, want DETERMINISTIC", cf.Determinism)
	}
	cf2 := cfOf(t, `CREATE FUNCTION ds.f(x INT64) RETURNS INT64 NOT DETERMINISTIC LANGUAGE js AS "return x"`)
	if cf2.Determinism != "NOT DETERMINISTIC" {
		t.Errorf("Determinism = %q, want NOT DETERMINISTIC", cf2.Determinism)
	}
}

func TestCreateFunction_PythonRemote(t *testing.T) {
	// DDL-017: Python UDF with LANGUAGE python WITH CONNECTION and OPTIONS.
	cf := cfOf(t, "CREATE FUNCTION mydataset.py_square(x FLOAT64) RETURNS FLOAT64 LANGUAGE python WITH CONNECTION `my-project.us.my-connection` OPTIONS (entry_point='square') AS \"def square(x):\\n  return x*x\"")
	if cf.Language != "python" {
		t.Errorf("Language = %q, want python", cf.Language)
	}
	if !cf.HasConnection || cf.ConnectionName == "" {
		t.Errorf("Connection: has=%v name=%q, want present", cf.HasConnection, cf.ConnectionName)
	}
	if len(cf.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(cf.Options))
	}
}

func TestCreateFunction_RemoteWithConnection(t *testing.T) {
	// DDL-017: remote function form REMOTE WITH CONNECTION conn.
	cf := cfOf(t, "CREATE FUNCTION ds.f(x INT64) RETURNS INT64 REMOTE WITH CONNECTION `p.us.c` OPTIONS(endpoint='https://x')")
	if !cf.Remote {
		t.Error("Remote = false, want true")
	}
	if !cf.HasConnection {
		t.Error("HasConnection = false, want true")
	}
}

func TestCreateFunction_Aggregate(t *testing.T) {
	// DDL-018: CREATE AGGREGATE FUNCTION … AS (SUM(x * x)).
	cf := cfOf(t, "CREATE AGGREGATE FUNCTION mydataset.sum_of_squares(x FLOAT64) RETURNS FLOAT64 AS (SUM(x * x))")
	if !cf.Aggregate {
		t.Error("Aggregate = false, want true")
	}
	if cf.Body == nil {
		t.Error("Body = nil")
	}
}

func TestCreateFunction_NotAggregateParam(t *testing.T) {
	// A NOT AGGREGATE parameter on an aggregate function.
	cf := cfOf(t, "CREATE AGGREGATE FUNCTION ds.f(x FLOAT64, y FLOAT64 NOT AGGREGATE) AS (SUM(x) + y)")
	if len(cf.Params) != 2 {
		t.Fatalf("Params = %d, want 2", len(cf.Params))
	}
	if !cf.Params[1].NotAggregate {
		t.Error("Params[1].NotAggregate = false, want true")
	}
}

func TestCreateTableFunction(t *testing.T) {
	// DDL-019: CREATE TABLE FUNCTION … AS query (no RETURNS TABLE).
	cf := cfOf(t, "CREATE OR REPLACE TABLE FUNCTION mydataset.names_by_year(y INT64) AS SELECT year, name FROM t WHERE year = y")
	if !cf.IsTableFunc {
		t.Error("IsTableFunc = false, want true")
	}
	if cf.AsQuery == nil {
		t.Error("AsQuery = nil, want the SELECT body")
	}
	if cf.ReturnsTable {
		t.Error("ReturnsTable = true, want false (no RETURNS TABLE)")
	}
}

func TestCreateTableFunction_ReturnsTable(t *testing.T) {
	// DDL-019: CREATE TABLE FUNCTION … RETURNS TABLE<name STRING, year INT64> AS query.
	cf := cfOf(t, "CREATE OR REPLACE TABLE FUNCTION mydataset.names_by_year(y INT64) RETURNS TABLE<name STRING, year INT64, total INT64> AS SELECT name, year, 1 AS total FROM t WHERE year = y")
	if !cf.IsTableFunc || !cf.ReturnsTable {
		t.Fatalf("IsTableFunc=%v ReturnsTable=%v, want both true", cf.IsTableFunc, cf.ReturnsTable)
	}
	if len(cf.ReturnColumns) != 3 {
		t.Fatalf("ReturnColumns = %d, want 3", len(cf.ReturnColumns))
	}
	if cf.ReturnColumns[0].Name != "name" || cf.ReturnColumns[0].Type.Text != "STRING" {
		t.Errorf("ReturnColumns[0] = %+v, want name STRING", cf.ReturnColumns[0])
	}
}

func TestCreateTableFunction_AnyTypeParam(t *testing.T) {
	// function_parameter with ANY TYPE / a TABLE<…> param.
	cf := cfOf(t, "CREATE TABLE FUNCTION ds.f(t TABLE<x INT64>, p ANY TYPE) AS SELECT 1 AS n")
	if len(cf.Params) != 2 {
		t.Fatalf("Params = %d, want 2", len(cf.Params))
	}
	if cf.Params[0].Type.Text != "TABLE<x INT64>" {
		t.Errorf("Params[0].Type = %q, want TABLE<x INT64>", cf.Params[0].Type.Text)
	}
	if cf.Params[1].Type.Text != "ANY TYPE" {
		t.Errorf("Params[1].Type = %q, want ANY TYPE", cf.Params[1].Type.Text)
	}
}

func TestCreateTableFunction_NoBodyAcceptedPerGrammar(t *testing.T) {
	// DEFEND (review finding): the legacy create_table_function_statement makes the
	// AS body OPTIONAL (`opt_as_query_or_string?`). BigQuery docs require AS query,
	// but the migration target is legacy-grammar parity and there is no Spanner
	// oracle for TVFs — so we FOLLOW THE GRAMMAR and accept a body-less TVF (the
	// "needs a body" rule is a ZetaSQL semantic check, not a parse rule).
	cf := cfOf(t, "CREATE TABLE FUNCTION ds.f(y INT64)")
	if !cf.IsTableFunc {
		t.Error("IsTableFunc = false, want true")
	}
	if cf.AsQuery != nil || cf.HasBodyString {
		t.Error("expected no body for a body-less TVF")
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

func cpOf(t *testing.T, sql string) *ast.CreateProcedureStmt {
	t.Helper()
	n := parseDDL(t, sql)
	cp, ok := n.(*ast.CreateProcedureStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateProcedureStmt", sql, n)
	}
	return cp
}

func TestCreateProcedure_BeginEnd(t *testing.T) {
	// DDL-020: CREATE OR REPLACE PROCEDURE … (IN tbl STRING, OUT result INT64) BEGIN … END;
	cp := cpOf(t, "CREATE OR REPLACE PROCEDURE mydataset.SelectFromTable(IN tbl STRING, OUT result INT64) BEGIN SELECT COUNT(*) FROM t; END")
	if cp.Name.String() != "mydataset.SelectFromTable" {
		t.Errorf("Name = %q", cp.Name.String())
	}
	if len(cp.Params) != 2 {
		t.Fatalf("Params = %d, want 2", len(cp.Params))
	}
	if cp.Params[0].Mode != ast.ParamModeIn || cp.Params[0].Name != "tbl" {
		t.Errorf("Params[0] = mode %v name %q, want IN tbl", cp.Params[0].Mode, cp.Params[0].Name)
	}
	if cp.Params[1].Mode != ast.ParamModeOut || cp.Params[1].Name != "result" {
		t.Errorf("Params[1] = mode %v name %q, want OUT result", cp.Params[1].Mode, cp.Params[1].Name)
	}
	if cp.BodyText == "" {
		t.Error("BodyText = empty, want the BEGIN…END block text")
	}
}

func TestCreateProcedure_NoParams(t *testing.T) {
	cp := cpOf(t, "CREATE PROCEDURE ds.p() BEGIN SELECT 1; END")
	if len(cp.Params) != 0 {
		t.Errorf("Params = %d, want 0", len(cp.Params))
	}
}

func TestCreateProcedure_BodyWithIfFunction(t *testing.T) {
	// Regression (review finding): the IF / CASE / WHILE / FOR keywords appear as
	// FUNCTION calls inside a body — `IF(cond, a, b)`, `SELECT CASE WHEN … END` —
	// and must NOT be counted as procedural block openers (else the body looks
	// unterminated). The '(' lookahead excludes the IF() function call.
	cp := cpOf(t, "CREATE PROCEDURE ds.p() BEGIN SELECT IF(TRUE, 1, 0); SELECT IFNULL(x, 0) FROM t; END")
	if cp.BodyText == "" {
		t.Fatal("BodyText empty; IF(...) function miscounted as a block opener")
	}
	// A CASE expression inside the body is balanced by its own END and must not
	// leak the outer block.
	cp2 := cpOf(t, "CREATE PROCEDURE ds.p() BEGIN SELECT CASE WHEN x THEN 1 ELSE 2 END FROM t; END")
	if cp2.BodyText == "" {
		t.Fatal("BodyText empty; CASE expression mishandled")
	}
}

func TestCreateProcedure_NestedBeginEnd(t *testing.T) {
	// A nested BEGIN…END and an IF…END IF must not close the outer block early.
	cp := cpOf(t, "CREATE PROCEDURE ds.p() BEGIN IF TRUE THEN SELECT 1; END IF; BEGIN SELECT 2; END; END")
	if cp.BodyText == "" {
		t.Fatal("BodyText empty")
	}
	// The captured block must include the final END (the body must extend past the
	// inner blocks).
	if got := cp.BodyText; got[len(got)-3:] != "END" {
		t.Errorf("BodyText does not end at the outer END: %q", got)
	}
}

func TestCreateProcedure_LanguageBody(t *testing.T) {
	cp := cpOf(t, "CREATE PROCEDURE ds.p(x INT64) OPTIONS(strict_mode=true) LANGUAGE python AS \"pass\"")
	if cp.Language != "python" {
		t.Errorf("Language = %q, want python", cp.Language)
	}
	if !cp.HasBodyString || cp.BodyString != "pass" {
		t.Errorf("BodyString = %q (has=%v), want pass", cp.BodyString, cp.HasBodyString)
	}
	if len(cp.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(cp.Options))
	}
}

func TestCreateProcedure_ExternalSecurity(t *testing.T) {
	cp := cpOf(t, "CREATE PROCEDURE ds.p() EXTERNAL SECURITY INVOKER BEGIN SELECT 1; END")
	if cp.ExternalSecurity != "INVOKER" {
		t.Errorf("ExternalSecurity = %q, want INVOKER", cp.ExternalSecurity)
	}
}

func TestCreateProcedure_INOUTParamMode(t *testing.T) {
	cp := cpOf(t, "CREATE PROCEDURE ds.p(INOUT x INT64) BEGIN SELECT x; END")
	if cp.Params[0].Mode != ast.ParamModeInout {
		t.Errorf("Params[0].Mode = %v, want INOUT", cp.Params[0].Mode)
	}
}

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
		t.Fatalf("Parse(%q): statement is %T, want *ast.BQDropStmt", sql, n)
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

func TestDropTableFunction(t *testing.T) {
	// DDL-050.
	d := bqDropOf(t, "DROP TABLE FUNCTION IF EXISTS mydataset.my_tvf")
	if d.Object != ast.BQDropTableFunction {
		t.Errorf("Object = %v, want TABLE FUNCTION", d.Object)
	}
}

func TestDropProcedure(t *testing.T) {
	// DDL-051.
	d := bqDropOf(t, "DROP PROCEDURE IF EXISTS mydataset.my_procedure")
	if d.Object != ast.BQDropProcedure {
		t.Errorf("Object = %v, want PROCEDURE", d.Object)
	}
}

func TestDropFunction_WithOverloadParams(t *testing.T) {
	// opt_function_parameters disambiguating an overload — consumed structurally.
	d := bqDropOf(t, "DROP FUNCTION ds.f(INT64, STRING)")
	if d.Object != ast.BQDropFunction || d.Name.String() != "ds.f" {
		t.Errorf("got Object=%v Name=%q", d.Object, d.Name.String())
	}
}
