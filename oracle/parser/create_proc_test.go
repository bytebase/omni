package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateProcedureSimple(t *testing.T) {
	sql := `CREATE PROCEDURE my_proc (p_id IN NUMBER, p_name IN VARCHAR2) IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateProcedureStmt)
	if !ok {
		t.Fatalf("expected CreateProcedureStmt, got %T", raw.Stmt)
	}
	if stmt.OrReplace {
		t.Error("expected OrReplace to be false")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_PROC" {
		t.Errorf("expected name MY_PROC, got %v", stmt.Name)
	}
	if stmt.Parameters == nil || stmt.Parameters.Len() != 2 {
		t.Fatalf("expected 2 parameters, got %d", stmt.Parameters.Len())
	}

	p0 := stmt.Parameters.Items[0].(*ast.Parameter)
	if p0.Name != "P_ID" {
		t.Errorf("expected param name P_ID, got %q", p0.Name)
	}
	if p0.Mode != "IN" {
		t.Errorf("expected mode IN, got %q", p0.Mode)
	}

	p1 := stmt.Parameters.Items[1].(*ast.Parameter)
	if p1.Name != "P_NAME" {
		t.Errorf("expected param name P_NAME, got %q", p1.Name)
	}

	if stmt.Body == nil {
		t.Error("expected non-nil body")
	}
}

func TestParseCreateOrReplaceProcedure(t *testing.T) {
	sql := `CREATE OR REPLACE PROCEDURE my_proc IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_PROC" {
		t.Errorf("expected name MY_PROC, got %v", stmt.Name)
	}
}

func TestParseCreateWrappedProcedure(t *testing.T) {
	sql := "CREATE OR REPLACE PROCEDURE wrapped_proc WRAPPED\n" +
		"a000000\n" +
		"abcd\n"
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateProcedureStmt)
	if !ok {
		t.Fatalf("expected CreateProcedureStmt, got %T", raw.Stmt)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "WRAPPED_PROC" {
		t.Errorf("expected name WRAPPED_PROC, got %v", stmt.Name)
	}
	if !stmt.Wrapped {
		t.Error("expected Wrapped to be true")
	}
	if stmt.WrappedSource != "WRAPPED\na000000\nabcd" {
		t.Fatalf("expected wrapped source to preserve raw payload, got %q", stmt.WrappedSource)
	}
	if stmt.Body != nil {
		t.Fatalf("expected nil body for wrapped procedure, got %T", stmt.Body)
	}
}

func TestParseCreateWrappedProcedureWithSemicolon(t *testing.T) {
	result := ParseAndCheck(t, "CREATE PROCEDURE wrapped_proc WRAPPED\na000000\nabcd;")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if !stmt.Wrapped {
		t.Fatal("expected Wrapped to be true")
	}
}

func TestParseCreateProcedureAuthID(t *testing.T) {
	sql := `CREATE PROCEDURE my_proc AUTHID CURRENT_USER IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if stmt.AuthID != "CURRENT_USER" {
		t.Fatalf("expected AUTHID CURRENT_USER, got %q", stmt.AuthID)
	}
}

func TestParseCreateFunctionSimple(t *testing.T) {
	sql := `CREATE FUNCTION get_total (p_id IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("expected CreateFunctionStmt, got %T", raw.Stmt)
	}
	if stmt.Name == nil || stmt.Name.Name != "GET_TOTAL" {
		t.Errorf("expected name GET_TOTAL, got %v", stmt.Name)
	}
	if stmt.Parameters == nil || stmt.Parameters.Len() != 1 {
		t.Fatalf("expected 1 parameter, got %d", stmt.Parameters.Len())
	}
	if stmt.ReturnType == nil {
		t.Fatal("expected non-nil ReturnType")
	}
	if stmt.ReturnType.Names.Len() == 0 {
		t.Fatal("expected non-empty ReturnType names")
	}
	nameStr := stmt.ReturnType.Names.Items[0].(*ast.String)
	if nameStr.Str != "NUMBER" {
		t.Errorf("expected return type NUMBER, got %q", nameStr.Str)
	}
	if stmt.Body == nil {
		t.Error("expected non-nil body")
	}
}

func TestParseCreateFunctionDeterministic(t *testing.T) {
	sql := `CREATE FUNCTION calc (x IN NUMBER) RETURN NUMBER DETERMINISTIC IS BEGIN RETURN x * 2; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateFunctionStmt)
	if !stmt.Deterministic {
		t.Error("expected Deterministic to be true")
	}
}

func TestParseCreateFunctionPipelined(t *testing.T) {
	sql := `CREATE FUNCTION pipe_fn RETURN NUMBER PIPELINED IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateFunctionStmt)
	if !stmt.Pipelined {
		t.Error("expected Pipelined to be true")
	}
}

func TestParseCreateFunctionAuthID(t *testing.T) {
	sql := `CREATE FUNCTION get_total RETURN NUMBER AUTHID DEFINER IS BEGIN RETURN 0; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateFunctionStmt)
	if stmt.AuthID != "DEFINER" {
		t.Fatalf("expected AUTHID DEFINER, got %q", stmt.AuthID)
	}
}

func TestParseCreateFunctionMultipleOptions(t *testing.T) {
	sql := `CREATE FUNCTION myfn RETURN NUMBER DETERMINISTIC PARALLEL_ENABLE RESULT_CACHE IS BEGIN RETURN 1; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateFunctionStmt)
	if !stmt.Deterministic {
		t.Error("expected Deterministic to be true")
	}
	if !stmt.Parallel {
		t.Error("expected Parallel to be true")
	}
	if !stmt.ResultCache {
		t.Error("expected ResultCache to be true")
	}
}

func TestParseCreatePackageSpec(t *testing.T) {
	sql := `CREATE PACKAGE my_pkg IS
  PROCEDURE do_something (p_id NUMBER);
  FUNCTION get_value RETURN NUMBER;
END my_pkg;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreatePackageStmt)
	if !ok {
		t.Fatalf("expected CreatePackageStmt, got %T", raw.Stmt)
	}
	if stmt.IsBody {
		t.Error("expected IsBody to be false for package spec")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_PKG" {
		t.Errorf("expected name MY_PKG, got %v", stmt.Name)
	}
	if stmt.Body == nil || stmt.Body.Len() < 2 {
		t.Fatalf("expected at least 2 declarations in package body, got %d", stmt.Body.Len())
	}
}

func TestParseCreatePackageBody(t *testing.T) {
	sql := `CREATE PACKAGE BODY my_pkg IS
  PROCEDURE do_something (p_id NUMBER) IS
  BEGIN
    NULL;
  END;
  FUNCTION get_value RETURN NUMBER IS
  BEGIN
    RETURN 42;
  END;
END my_pkg;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreatePackageStmt)
	if !ok {
		t.Fatalf("expected CreatePackageStmt, got %T", raw.Stmt)
	}
	if !stmt.IsBody {
		t.Error("expected IsBody to be true for package body")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_PKG" {
		t.Errorf("expected name MY_PKG, got %v", stmt.Name)
	}
	if stmt.Body == nil || stmt.Body.Len() < 2 {
		t.Fatalf("expected at least 2 declarations, got %d", stmt.Body.Len())
	}
}

func TestParseParameterInOutModes(t *testing.T) {
	sql := `CREATE PROCEDURE modes_proc (p_in IN NUMBER, p_out OUT VARCHAR2, p_inout IN OUT NUMBER) IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if stmt.Parameters == nil || stmt.Parameters.Len() != 3 {
		t.Fatalf("expected 3 parameters, got %d", stmt.Parameters.Len())
	}

	p0 := stmt.Parameters.Items[0].(*ast.Parameter)
	if p0.Mode != "IN" {
		t.Errorf("expected mode IN, got %q", p0.Mode)
	}

	p1 := stmt.Parameters.Items[1].(*ast.Parameter)
	if p1.Mode != "OUT" {
		t.Errorf("expected mode OUT, got %q", p1.Mode)
	}

	p2 := stmt.Parameters.Items[2].(*ast.Parameter)
	if p2.Mode != "IN OUT" {
		t.Errorf("expected mode IN OUT, got %q", p2.Mode)
	}
}

func TestParseParameterDefault(t *testing.T) {
	sql := `CREATE PROCEDURE def_proc (p_val NUMBER DEFAULT 10) IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if stmt.Parameters == nil || stmt.Parameters.Len() != 1 {
		t.Fatalf("expected 1 parameter, got %d", stmt.Parameters.Len())
	}

	p0 := stmt.Parameters.Items[0].(*ast.Parameter)
	if p0.Name != "P_VAL" {
		t.Errorf("expected param name P_VAL, got %q", p0.Name)
	}
	if p0.Default == nil {
		t.Error("expected non-nil default expression")
	}
}

func TestParseParameterDefaultAssign(t *testing.T) {
	sql := `CREATE PROCEDURE def_proc (p_val NUMBER := 10) IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	p0 := stmt.Parameters.Items[0].(*ast.Parameter)
	if p0.Default == nil {
		t.Error("expected non-nil default expression")
	}
}

func TestParseCreateOrReplacePackageBody(t *testing.T) {
	sql := `CREATE OR REPLACE PACKAGE BODY my_pkg IS
  PROCEDURE do_it IS
  BEGIN
    NULL;
  END;
END my_pkg;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreatePackageStmt)
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
	if !stmt.IsBody {
		t.Error("expected IsBody to be true")
	}
}

func TestParseCreateProcedureLoc(t *testing.T) {
	sql := `CREATE PROCEDURE my_proc IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got End=%d Start=%d", stmt.Loc.End, stmt.Loc.Start)
	}
}

func TestParseCreateFunctionSchemaQualified(t *testing.T) {
	sql := `CREATE FUNCTION myschema.get_val RETURN NUMBER IS BEGIN RETURN 1; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateFunctionStmt)
	if stmt.Name.Schema != "MYSCHEMA" {
		t.Errorf("expected schema MYSCHEMA, got %q", stmt.Name.Schema)
	}
	if stmt.Name.Name != "GET_VAL" {
		t.Errorf("expected name GET_VAL, got %q", stmt.Name.Name)
	}
}

func TestParseParameterNocopy(t *testing.T) {
	sql := `CREATE PROCEDURE nc_proc (p_data OUT NOCOPY VARCHAR2) IS BEGIN NULL; END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	p0 := stmt.Parameters.Items[0].(*ast.Parameter)
	if p0.Mode != "OUT NOCOPY" {
		t.Errorf("expected mode 'OUT NOCOPY', got %q", p0.Mode)
	}
}
