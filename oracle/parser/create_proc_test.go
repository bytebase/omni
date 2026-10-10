package parser

import (
	"strings"
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

func TestParseCreateWrappedProcedureWithOperatorPayload(t *testing.T) {
	result := ParseAndCheck(t, "CREATE PROCEDURE wrapped_proc WRAPPED\na000000\nabc=\n")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if !stmt.Wrapped {
		t.Fatal("expected Wrapped to be true")
	}
	if stmt.WrappedSource != "WRAPPED\na000000\nabc=" {
		t.Fatalf("expected wrapped source to preserve opaque payload, got %q", stmt.WrappedSource)
	}
}

func TestParseCreateWrappedProcedureWithCommentBeforeName(t *testing.T) {
	result := ParseAndCheck(t, "CREATE OR REPLACE PROCEDURE\n--++WRAP_VERSION:1\n  wrapped_proc wrapped\na000000\nabc=\n")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateProcedureStmt)
	if stmt.Name == nil || stmt.Name.Name != "WRAPPED_PROC" {
		t.Fatalf("expected name WRAPPED_PROC, got %v", stmt.Name)
	}
	if !stmt.Wrapped {
		t.Fatal("expected Wrapped to be true")
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

func TestParseCallSpec(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want ast.CallSpec
	}{
		{
			name: "java function call spec",
			sql: "CREATE OR REPLACE EDITIONABLE FUNCTION \"S\".\"F\" (p IN BLOB, q IN VARCHAR2) RETURN VARCHAR2\n" +
				"AS LANGUAGE JAVA\n" +
				"NAME 'Signer.check(oracle.sql.BLOB, java.lang.String) return java.lang.String';",
			want: ast.CallSpec{Language: "JAVA", Name: "Signer.check(oracle.sql.BLOB, java.lang.String) return java.lang.String"},
		},
		{
			name: "java procedure call spec after IS",
			sql:  "CREATE PROCEDURE p(a IN VARCHAR2) IS LANGUAGE JAVA NAME 'X.y(java.lang.String)';",
			want: ast.CallSpec{Language: "JAVA", Name: "X.y(java.lang.String)"},
		},
		{
			name: "c call spec with every clause",
			sql: "CREATE FUNCTION f(p IN OUT BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C " +
				"LIBRARY s.lib NAME \"c_f\" AGENT IN (p) WITH CONTEXT " +
				"PARAMETERS (CONTEXT, p BY REFERENCE INT, p INDICATOR SHORT, RETURN UNSIGNED INT);",
			want: ast.CallSpec{Language: "C", Name: "c_f", WithContext: true},
		},
		{
			name: "external call spec",
			sql:  "CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS EXTERNAL NAME c_f LIBRARY lib;",
			want: ast.CallSpec{Language: "C", External: true, Name: "C_F"},
		},
		// The superseded EXTERNAL form with its LANGUAGE and CALLING STANDARD
		// clauses, in the orders Oracle 23ai compiles VALID.
		{
			name: "external language c then name and library",
			sql:  "CREATE OR REPLACE PROCEDURE p AS EXTERNAL LANGUAGE C NAME \"c_p\" LIBRARY lib;",
			want: ast.CallSpec{Language: "C", External: true, Name: "c_p"},
		},
		{
			name: "external calling standard after name",
			sql:  "CREATE OR REPLACE PROCEDURE p AS EXTERNAL LIBRARY lib NAME \"c_p\" CALLING STANDARD C;",
			want: ast.CallSpec{Language: "C", External: true, Name: "c_p", CallingStandard: "C"},
		},
		{
			name: "external calling standard pascal first",
			sql:  "CREATE OR REPLACE PROCEDURE p AS EXTERNAL CALLING STANDARD PASCAL LIBRARY lib;",
			want: ast.CallSpec{Language: "C", External: true, CallingStandard: "PASCAL"},
		},
		{
			name: "external parameters before library",
			sql:  "CREATE OR REPLACE PROCEDURE p(a IN BINARY_INTEGER) AS EXTERNAL PARAMETERS (a INT) LIBRARY lib;",
			want: ast.CallSpec{Language: "C", External: true},
		},
		{
			name: "external with context then name",
			sql:  "CREATE OR REPLACE PROCEDURE p AS EXTERNAL LIBRARY lib WITH CONTEXT NAME \"c_p\";",
			want: ast.CallSpec{Language: "C", External: true, Name: "c_p", WithContext: true},
		},
		{
			name: "language c calling standard before library",
			sql:  "CREATE OR REPLACE PROCEDURE p AS LANGUAGE C CALLING STANDARD C LIBRARY lib;",
			want: ast.CallSpec{Language: "C", CallingStandard: "C"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			var spec *ast.CallSpec
			var body ast.StmtNode
			switch stmt := result.Items[0].(*ast.RawStmt).Stmt.(type) {
			case *ast.CreateFunctionStmt:
				spec, body = stmt.CallSpec, stmt.Body
			case *ast.CreateProcedureStmt:
				spec, body = stmt.CallSpec, stmt.Body
			default:
				t.Fatalf("unexpected statement %T", stmt)
			}
			if body != nil {
				t.Fatalf("Body = %T, want nil for a call spec", body)
			}
			if spec == nil {
				t.Fatal("CallSpec = nil")
			}
			if spec.Language != tt.want.Language || spec.External != tt.want.External ||
				spec.Name != tt.want.Name || spec.WithContext != tt.want.WithContext ||
				spec.CallingStandard != tt.want.CallingStandard {
				t.Fatalf("CallSpec = %+v, want %+v", *spec, tt.want)
			}
			if violations := CheckLocations(t, tt.sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

func TestParseCallSpecCClauses(t *testing.T) {
	sql := "CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C " +
		"NAME \"c_f\" LIBRARY s.lib AGENT IN (p) PARAMETERS (CONTEXT, p INDICATOR STRUCT, RETURN INT);"
	result := ParseAndCheck(t, sql)
	spec := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateFunctionStmt).CallSpec
	if spec.Library == nil || spec.Library.Schema != "S" || spec.Library.Name != "LIB" {
		t.Fatalf("Library = %+v, want S.LIB", spec.Library)
	}
	if got := callSpecWords(spec.AgentIn); got != "P" {
		t.Fatalf("AgentIn = %q, want P", got)
	}
	if got := callSpecWords(spec.Parameters); got != "CONTEXT,P INDICATOR STRUCT,RETURN INT" {
		t.Fatalf("Parameters = %q", got)
	}
}

// TestParseCallSpecExternalParameters covers the external_parameter shapes
// Oracle 23ai parses. Some of them fail a later check for the parameter they
// name (PLS-00235, PLS-00250, PLS-00253); none is a syntax error.
func TestParseCallSpecExternalParameters(t *testing.T) {
	for _, param := range []string{
		"x",
		"x INT",
		"x UNSIGNED INT",
		"x BY REFERENCE",
		"x BY REFERENCE UNSIGNED INT",
		"x BY VALUE INT",
		"x INDICATOR",
		"x INDICATOR STRUCT",
		"x INDICATOR IN INT",
		"x INDICATOR OUT INT",
		"x INDICATOR BY REFERENCE INT",
		"x LENGTH INT",
		"x LENGTH BY VALUE INT",
		"x MAXLEN INT",
		"x CHARSETID UB4",
		"x CHARSETFORM UB1",
		"x DURATION OCIDURATION",
		"x TDO",
		"x NATIVE INT",
		"x NATIVE",
		"x ARRAY",
		"x STRING",
		"x OCISTRING",
		"\"X\" INT",
		"CONTEXT",
		"CONTEXT INT",
		"SELF",
		"SELF TDO",
		"RETURN",
		"RETURN INDICATOR SHORT",
	} {
		sql := "CREATE PROCEDURE p(x IN BINARY_INTEGER) AS LANGUAGE C LIBRARY lib PARAMETERS (" + param + ");"
		t.Run(param, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseCallSpecExternalParameterRejects lists external parameters and
// AGENT IN arguments Oracle 23ai rejects with PLS-00103.
func TestParseCallSpecExternalParameterRejects(t *testing.T) {
	for _, clause := range []string{
		"PARAMETERS (x BY REFERENCE BY REFERENCE INT)",
		"PARAMETERS (x INT INT)",
		"PARAMETERS (x FOO)",
		"PARAMETERS (x UB8)",
		"PARAMETERS (x \"INT\")",
		"PARAMETERS (x UNSIGNED)",
		"PARAMETERS (x LONG RAW)",
		"PARAMETERS (x INDICATOR INDICATOR)",
		"PARAMETERS (x LENGTH LENGTH)",
		"PARAMETERS (x INDICATOR STRUCT STRUCT)",
		"PARAMETERS (x INDICATOR STRUCT BY REFERENCE INT)",
		"PARAMETERS (x INDICATOR TDO)",
		"PARAMETERS (x IN INT)",
		"PARAMETERS (x ARRAY INT)",
		"PARAMETERS (x BY INT)",
		"AGENT IN (x y)",
		"AGENT IN (x, y)",
		"AGENT IN ()",
	} {
		sql := "CREATE PROCEDURE p(x IN BINARY_INTEGER, y IN BINARY_INTEGER) AS LANGUAGE C LIBRARY lib " + clause + ";"
		t.Run(clause, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}

func callSpecWords(list *ast.List) string {
	if list == nil {
		return ""
	}
	var words []string
	for _, item := range list.Items {
		words = append(words, item.(*ast.String).Str)
	}
	return strings.Join(words, ",")
}

// TestParseCallSpecInSubprograms covers call specs where a subprogram is
// defined inside another unit. Oracle 23ai compiles each one VALID, except
// the nested Java call spec, which it rejects with PLS-00999 (implementation
// restriction), a semantic error after parsing.
func TestParseCallSpecInSubprograms(t *testing.T) {
	tests := []string{
		"CREATE OR REPLACE PACKAGE pk AS FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int'; END;",
		"CREATE OR REPLACE PACKAGE BODY pk AS\n" +
			"  FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';\n" +
			"  PROCEDURE p IS BEGIN NULL; END;\n" +
			"  PROCEDURE q AS LANGUAGE C LIBRARY lib;\n" +
			"END pk;",
		"CREATE OR REPLACE TYPE BODY ty AS MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int'; END;",
		"CREATE OR REPLACE PROCEDURE p AS FUNCTION g RETURN NUMBER AS LANGUAGE JAVA NAME 'X.g() return int'; BEGIN NULL; END;",
		"CREATE OR REPLACE PACKAGE pk AS FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }}; PROCEDURE q; END;",
		"CREATE OR REPLACE PACKAGE pk AS FUNCTION f RETURN NUMBER AS MLE MODULE m SIGNATURE 'f()'; PROCEDURE q; END;",
		"CREATE OR REPLACE PACKAGE BODY pk AS\n" +
			"  FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return \"a;b\"; -- x }};\n" +
			"  PROCEDURE z IS BEGIN NULL; END;\n" +
			"END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseCallSpecNamesAreNotCallSpecs checks the names Oracle 23ai compiles
// VALID although they spell a call spec word: parameters named LANGUAGE,
// EXTERNAL, or MLE, a variable so named after the first declaration or in a
// DECLARE block or package, and a quoted "MLE" first declaration.
// TestParseCallSpecRejects has the unquoted first declaration, which Oracle
// reads as a call spec.
func TestParseCallSpecNamesAreNotCallSpecs(t *testing.T) {
	tests := []string{
		"CREATE FUNCTION f(language IN VARCHAR2, external IN NUMBER, mle IN NUMBER) RETURN NUMBER IS BEGIN RETURN 1; END;",
		"CREATE FUNCTION f RETURN NUMBER IS x NUMBER; mle NUMBER; BEGIN mle := 1; RETURN mle; END;",
		"CREATE FUNCTION f RETURN NUMBER IS \"MLE\" NUMBER; BEGIN \"MLE\" := 1; RETURN \"MLE\"; END;",
		"CREATE PROCEDURE p AS BEGIN DECLARE external NUMBER; BEGIN NULL; END; END;",
		"CREATE PACKAGE pk AS mle NUMBER; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			switch stmt := result.Items[0].(*ast.RawStmt).Stmt.(type) {
			case *ast.CreateFunctionStmt:
				if stmt.CallSpec != nil || stmt.Body == nil {
					t.Fatalf("CallSpec = %+v, Body = %v; want a PL/SQL body", stmt.CallSpec, stmt.Body)
				}
			case *ast.CreateProcedureStmt:
				if stmt.CallSpec != nil || stmt.Body == nil {
					t.Fatalf("CallSpec = %+v, Body = %v; want a PL/SQL body", stmt.CallSpec, stmt.Body)
				}
			}
		})
	}
}

// TestParseMLECallSpec covers the MLE call specs, which Oracle 23ai compiles
// past the syntax check: the MODULE form reaches name resolution (PLS-00201
// for an absent module), the inline form ORA-00439 where MLE is not enabled.
// The inline code is kept verbatim with its delimiters; it is JavaScript, so
// its quotes, "--", and ';' are not SQL.
func TestParseMLECallSpec(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		want   ast.CallSpec
		module string
		env    string
	}{
		{
			name:   "module",
			sql:    "CREATE FUNCTION f(a IN NUMBER) RETURN NUMBER AS MLE MODULE m SIGNATURE 'f(number)';",
			want:   ast.CallSpec{MLE: true, Name: "f(number)"},
			module: "M",
		},
		{
			name:   "module with env",
			sql:    "CREATE OR REPLACE FUNCTION f RETURN NUMBER AS MLE MODULE s.m ENV s.e SIGNATURE 'f()';",
			want:   ast.CallSpec{MLE: true, Name: "f()"},
			module: "M",
			env:    "E",
		},
		{
			name: "inline",
			sql:  "CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }};",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Code: "{{ return 1; }}"},
		},
		{
			name: "inline pure",
			sql:  "CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT PURE {{ return 1; }};",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Pure: true, Code: "{{ return 1; }}"},
		},
		{
			name: "inline code the SQL lexer would misread",
			sql:  "CREATE PROCEDURE p AS MLE LANGUAGE JAVASCRIPT {{ let s = \"a;b'\"; -- x }};",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Code: "{{ let s = \"a;b'\"; -- x }}"},
		},
		{
			// A delimiter that is not made of opening brackets closes with
			// itself, and the search for it ignores JavaScript strings.
			name: "identical delimiter",
			sql:  "CREATE PROCEDURE p AS MLE LANGUAGE JAVASCRIPT ## let s = 'it''s'; /* ## ;",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Code: "## let s = 'it''s'; /* ##"},
		},
		{
			name: "quote delimiter",
			sql:  "CREATE PROCEDURE p AS MLE LANGUAGE JAVASCRIPT 'q console.log(1); 'q;",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Code: "'q console.log(1); 'q"},
		},
		{
			// Oracle 23ai takes a reserved word as a delimiter too (BEGIN,
			// SELECT, END, IS, and TABLE all pass the syntax check).
			name: "reserved word delimiter",
			sql:  "CREATE PROCEDURE p AS MLE LANGUAGE JAVASCRIPT BEGIN return; BEGIN;",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Code: "BEGIN return; BEGIN"},
		},
		{
			name: "mirrored bracket delimiter",
			sql:  "CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {< return 1; >};",
			want: ast.CallSpec{MLE: true, Language: "JAVASCRIPT", Code: "{< return 1; >}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			var spec *ast.CallSpec
			switch stmt := result.Items[0].(*ast.RawStmt).Stmt.(type) {
			case *ast.CreateFunctionStmt:
				spec = stmt.CallSpec
			case *ast.CreateProcedureStmt:
				spec = stmt.CallSpec
			default:
				t.Fatalf("unexpected statement %T", stmt)
			}
			if spec == nil {
				t.Fatal("CallSpec = nil")
			}
			if spec.MLE != tt.want.MLE || spec.Language != tt.want.Language || spec.Name != tt.want.Name ||
				spec.Pure != tt.want.Pure || spec.Code != tt.want.Code {
				t.Fatalf("CallSpec = %s, want %+v", ast.NodeToString(spec), tt.want)
			}
			if got := objectNameText(spec.Module); got != tt.module {
				t.Fatalf("Module = %q, want %q", got, tt.module)
			}
			if got := objectNameText(spec.Env); got != tt.env {
				t.Fatalf("Env = %q, want %q", got, tt.env)
			}
			if violations := CheckLocations(t, tt.sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

func objectNameText(n *ast.ObjectName) string {
	if n == nil {
		return ""
	}
	return n.Name
}

// TestParseMLECallSpecRejects lists MLE call specs Oracle 23ai compiles with
// PLS-00103: SIGNATURE is required and takes a string, PURE belongs to the
// inline form only, ENV needs a name, and a ';' must follow the call spec,
// directly after the inline code's closing delimiter. A lone "." is no
// delimiter (PLS-00881), and code without its closing delimiter never ends.
func TestParseMLECallSpecRejects(t *testing.T) {
	tests := []string{
		"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m;",
		"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m SIGNATURE f;",
		"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m PURE SIGNATURE 'f()';",
		"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m ENV SIGNATURE 'f()';",
		"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m SIGNATURE 'f()'",
		"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE SIGNATURE 'f()';",
		"CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }}",
		"CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }} x;",
		"CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT . return 1; .;",
		"CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; ;",
		"CREATE FUNCTION f RETURN NUMBER AS MLE LANGUAGE;",
		"CREATE FUNCTION f RETURN NUMBER AS MLE JAVASCRIPT {{ return 1; }};",
		"CREATE OR REPLACE PACKAGE pk AS FUNCTION f RETURN NUMBER AS MLE MODULE m SIGNATURE 'f()' PROCEDURE q; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}

// TestParseCallSpecRejects lists call specs Oracle 23ai compiles with
// PLS-00103, plus LANGUAGE C without LIBRARY (PLS-00247). Each case past the
// first group ends with ';' so it fails for its own rule, not for a missing
// terminator.
func TestParseCallSpecRejects(t *testing.T) {
	tests := []string{
		// PLS-00103: the call spec's terminating ';' is required.
		"CREATE PROCEDURE p AS LANGUAGE JAVA NAME 'X.y()'",
		"CREATE PROCEDURE p AS EXTERNAL NAME \"c_p\" LIBRARY lib",
		"CREATE OR REPLACE PACKAGE pk AS FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int' PROCEDURE q; END;",
		"CREATE FUNCTION f RETURN NUMBER AS LANGUAGE JAVA;",
		"CREATE FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME;",
		"CREATE FUNCTION f RETURN VARCHAR2 AS LANGUAGE JAVA NAME \"X.y() return java.lang.String\";",
		"CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C NAME \"c_f\";",
		"CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C NAME 'c_f' LIBRARY lib;",
		"CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C LIBRARY lib PARAMETERS ();",
		"CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C LIBRARY lib WITH;",
		"CREATE FUNCTION f(p IN BINARY_INTEGER) RETURN BINARY_INTEGER AS LANGUAGE C LIBRARY lib AGENT (p);",
		// PLS-00103: LANGUAGE takes only C, CALLING STANDARD only C or PASCAL.
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib LANGUAGE PASCAL;",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib CALLING STANDARD;",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib CALLING STANDARD FORTRAN;",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib CALLING \"STANDARD\" C;",
		// PLS-00247: LIBRARY is required.
		"CREATE PROCEDURE p AS EXTERNAL LANGUAGE C;",
		// PLS-00139, -00140, -00142 through -00145, -00171: each C clause at
		// most once.
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib NAME \"a\" NAME \"b\";",
		"CREATE PROCEDURE p AS LANGUAGE C NAME \"c\" LIBRARY lib AGENT IN (x) NAME \"d\";",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib LIBRARY lib;",
		"CREATE PROCEDURE p(a IN BINARY_INTEGER) AS EXTERNAL LIBRARY lib PARAMETERS (a INT) PARAMETERS (a INT);",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib LANGUAGE C LANGUAGE C;",
		"CREATE PROCEDURE p AS LANGUAGE C LIBRARY lib LANGUAGE C;",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib CALLING STANDARD C CALLING STANDARD C;",
		"CREATE PROCEDURE p AS EXTERNAL LIBRARY lib WITH CONTEXT WITH CONTEXT;",
		"CREATE PROCEDURE p(a IN BINARY_INTEGER) AS EXTERNAL LIBRARY lib AGENT IN (a) AGENT IN (a);",
		// PLS-00103: an unquoted LANGUAGE, EXTERNAL, or MLE right after IS|AS
		// starts a call spec in every subprogram, so none of them names the
		// first declaration.
		"CREATE FUNCTION f RETURN NUMBER IS external NUMBER; BEGIN external := 1; RETURN external; END;",
		"CREATE FUNCTION f RETURN NUMBER IS language NUMBER; BEGIN language := 1; RETURN language; END;",
		"CREATE FUNCTION f RETURN NUMBER IS mle NUMBER; BEGIN mle := 1; RETURN mle; END;",
		"CREATE PROCEDURE p AS PROCEDURE q IS external NUMBER; BEGIN NULL; END; BEGIN NULL; END;",
		"CREATE PACKAGE BODY pk AS FUNCTION f RETURN NUMBER IS language NUMBER; BEGIN RETURN 1; END; END;",
		"CREATE TYPE BODY ty AS MEMBER FUNCTION f RETURN NUMBER IS mle NUMBER; BEGIN RETURN 1; END; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}
