package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func mustCreateRoutine(t *testing.T, input string) *ast.CreateRoutineStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateRoutineStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateRoutineStmt", input, node)
	}
	return stmt
}

// mustCreateRoutineDirect parses ONE statement via parseSingle, bypassing F3's
// Split. It is used for the multi-line single-quoted body canaries: those bodies
// contain an embedded ';' that F3's lexer-driven Split mis-segments (the omni
// lexer's scanString cannot span a newline, so a ';' that follows the newline
// inside the body is mistaken for a statement terminator — a Split-layer
// limitation tracked as a flagged divergence). parseSingle receives the whole
// statement text, so it exercises THIS node's body raw-scan in isolation, which
// is the layer this node owns. Forms that survive Split (single-line '…',
// multi-line '…' with no embedded ';', and every $$…$$ body) are still asserted
// end-to-end via mustCreateRoutine / the corpus harness.
func mustCreateRoutineDirect(t *testing.T, input string) *ast.CreateRoutineStmt {
	t.Helper()
	node, errs := parseSingle(input, 0)
	if len(errs) > 0 {
		t.Fatalf("parseSingle %q: %d error(s): %v", input, len(errs), errs)
	}
	stmt, ok := node.(*ast.CreateRoutineStmt)
	if !ok {
		t.Fatalf("parseSingle %q: got %T, want *ast.CreateRoutineStmt", input, node)
	}
	return stmt
}

func mustAlterRoutine(t *testing.T, input string) *ast.AlterRoutineStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterRoutineStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterRoutineStmt", input, node)
	}
	return stmt
}

// assertParseError parses input and asserts at least one error is produced (a
// negative test). It does not require a specific message — only that omni
// rejects the statement.
func assertParseError(t *testing.T, input string) {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) == 0 {
		t.Fatalf("parse %q: expected an error, got none (stmts=%d)", input, len(result.File.Stmts))
	}
}

// ---------------------------------------------------------------------------
// ⭐ CANARY — body handling: single-line '…', MULTI-LINE '…', and $$…$$
// ---------------------------------------------------------------------------

func TestCreateRoutine_OrAlter(t *testing.T) {
	// FUNCTION (T4.5) and PROCEDURE both flow through parseCreateRoutineBody;
	// CREATE OR ALTER threads OrAlter to the shared CreateRoutineStmt node.
	fn := mustCreateRoutine(t, "CREATE OR ALTER FUNCTION multiply(a NUMBER, b NUMBER) RETURNS NUMBER AS 'a * b'")
	if !fn.OrAlter || fn.OrReplace || fn.Kind != ast.RoutineFunction {
		t.Errorf("function: OrAlter=%v OrReplace=%v Kind=%v, want true/false/RoutineFunction", fn.OrAlter, fn.OrReplace, fn.Kind)
	}
	proc := mustCreateRoutine(t, "CREATE OR ALTER PROCEDURE p(a NUMBER) RETURNS NUMBER LANGUAGE SQL AS 'BEGIN RETURN a; END'")
	if !proc.OrAlter || proc.OrReplace || proc.Kind != ast.RoutineProcedure {
		t.Errorf("procedure: OrAlter=%v OrReplace=%v Kind=%v, want true/false/RoutineProcedure", proc.OrAlter, proc.OrReplace, proc.Kind)
	}
	// SECURE pairs with OR ALTER (early-dispatch path in parseCreateStmt).
	secFn := mustCreateRoutine(t, "CREATE OR ALTER SECURE FUNCTION sf(a NUMBER) RETURNS NUMBER AS 'a'")
	if !secFn.OrAlter || !secFn.Secure || secFn.OrReplace {
		t.Errorf("secure function: OrAlter=%v Secure=%v OrReplace=%v, want true/true/false", secFn.OrAlter, secFn.Secure, secFn.OrReplace)
	}
}

func TestRoutineBody_SingleLineSingleQuote(t *testing.T) {
	stmt := mustCreateRoutine(t, "CREATE FUNCTION multiply1 (a number, b number) RETURNS number AS 'a * b'")
	if stmt.Body != "'a * b'" {
		t.Errorf("Body = %q, want %q", stmt.Body, "'a * b'")
	}
	if stmt.BodyDollar {
		t.Errorf("BodyDollar = true, want false")
	}
}

// The defining canary: a MULTI-LINE single-quoted body. The omni lexer cannot
// tokenize this (scanString rejects the embedded newline); the parser must
// raw-scan it from source, verbatim, with the embedded newline and quotes
// preserved, and the statement must parse with ZERO errors. The body also
// contains a ';', which F3's Split mis-segments (a Split-layer limitation; see
// mustCreateRoutineDirect), so this canary is asserted directly against the
// parser via parseSingle — the layer this node owns.
func TestRoutineBody_MultiLineSingleQuote(t *testing.T) {
	input := "CREATE FUNCTION f() RETURNS VARCHAR AS 'var r = 1;\n  return r;'"
	stmt := mustCreateRoutineDirect(t, input)
	want := "'var r = 1;\n  return r;'"
	if stmt.Body != want {
		t.Errorf("Body = %q, want %q", stmt.Body, want)
	}
	if stmt.BodyDollar {
		t.Errorf("BodyDollar = true, want false")
	}
}

// A multi-line single-quoted SQL body with NO embedded ';' survives F3's Split
// (one segment), so it is asserted end-to-end through ParseBestEffort — proving
// the body raw-scan works across the full pipeline when Split cooperates. This
// is the create-function example_14 shape.
func TestRoutineBody_MultiLineSingleQuoteSQL(t *testing.T) {
	input := "CREATE OR REPLACE FUNCTION get_x ( id NUMBER )\n" +
		"  RETURNS TABLE (country_code CHAR, country_name VARCHAR)\n" +
		"  AS 'SELECT DISTINCT c.country_code, c.country_name\n" +
		"      FROM user_addresses a, countries c\n" +
		"      WHERE a.user_id = id\n" +
		"      AND c.country_code = a.country_code'"
	stmt := mustCreateRoutine(t, input)
	if !strings.HasPrefix(stmt.Body, "'SELECT DISTINCT") || !strings.HasSuffix(stmt.Body, "a.country_code'") {
		t.Errorf("Body = %q, want the verbatim multi-line single-quoted SELECT", stmt.Body)
	}
	if strings.Count(stmt.Body, "\n") != 3 {
		t.Errorf("Body newline count = %d, want 3 (verbatim preserved)", strings.Count(stmt.Body, "\n"))
	}
	if stmt.ReturnTable == nil {
		t.Errorf("ReturnTable = nil, want the RETURNS TABLE columns")
	}
}

func TestRoutineBody_DollarQuoted(t *testing.T) {
	input := "CREATE OR REPLACE FUNCTION f() RETURNS VARIANT LANGUAGE PYTHON HANDLER = 'udf' AS $$\nimport numpy\ndef udf():\n    return 1\n$$"
	stmt := mustCreateRoutine(t, input)
	if !strings.HasPrefix(stmt.Body, "$$") || !strings.HasSuffix(stmt.Body, "$$") {
		t.Errorf("Body = %q, want a $$-delimited body", stmt.Body)
	}
	if !stmt.BodyDollar {
		t.Errorf("BodyDollar = false, want true")
	}
	if !strings.Contains(stmt.Body, "import numpy") {
		t.Errorf("Body = %q, want to contain the verbatim Python source", stmt.Body)
	}
}

func TestRoutineBody_DollarQuotedEmpty(t *testing.T) {
	// $$ $$ (empty body) — the legacy regression corpus uses this.
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS BOOLEAN AS $$ $$")
	if stmt.Body != "$$ $$" {
		t.Errorf("Body = %q, want %q", stmt.Body, "$$ $$")
	}
	if !stmt.BodyDollar {
		t.Errorf("BodyDollar = false, want true")
	}
}

func TestRoutineBody_SingleQuoteEscapedQuote(t *testing.T) {
	// A doubled single-quote inside the body is an escaped quote and must not
	// terminate the body early.
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS VARCHAR AS 'a ''quoted'' b'")
	if stmt.Body != "'a ''quoted'' b'" {
		t.Errorf("Body = %q, want %q", stmt.Body, "'a ''quoted'' b'")
	}
}

func TestRoutineBody_TrailingSemicolonNotInBody(t *testing.T) {
	// The trailing ';' is the statement delimiter (stripped by Split) and must
	// not be absorbed into the body.
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS FLOAT AS '3.14::FLOAT';")
	if stmt.Body != "'3.14::FLOAT'" {
		t.Errorf("Body = %q, want %q", stmt.Body, "'3.14::FLOAT'")
	}
}

// ---------------------------------------------------------------------------
// PROCEDURE AS <bare scripting block> body — the non-$$, non-quoted form
//
//	CREATE PROCEDURE p() RETURNS ... AS DECLARE ... BEGIN ... END
//
// The body after AS is a Snowflake Scripting block (DECLARE..BEGIN..END or a
// bare BEGIN..END) with NO surrounding quotes or $$. F3's Split already keeps
// the whole block in one segment (the DECLARE/BEGIN..END state machine), and
// parseRoutineBody routes a leading DECLARE/BEGIN body to the scripting parser,
// capturing the block text verbatim into Body. These flow end-to-end through
// mustCreateRoutine (the full Split + parse pipeline).
// ---------------------------------------------------------------------------

func TestRoutineBody_AsDeclareBeginBlock(t *testing.T) {
	input := "CREATE PROCEDURE p() RETURNS INTEGER AS\n" +
		"DECLARE\n  x INTEGER DEFAULT 0;\nBEGIN\n  x := 1;\n  RETURN x;\nEND"
	stmt := mustCreateRoutine(t, input)
	if stmt.Kind != ast.RoutineProcedure {
		t.Errorf("Kind = %v, want RoutineProcedure", stmt.Kind)
	}
	if stmt.BodyDollar {
		t.Errorf("BodyDollar = true, want false")
	}
	if !strings.HasPrefix(stmt.Body, "DECLARE") || !strings.HasSuffix(stmt.Body, "END") {
		t.Errorf("Body = %q, want the verbatim DECLARE..BEGIN..END block", stmt.Body)
	}
	if !strings.Contains(stmt.Body, "x := 1;") {
		t.Errorf("Body = %q, want it to contain the verbatim block statements", stmt.Body)
	}
}

func TestRoutineBody_AsBareBeginBlock(t *testing.T) {
	// AS BEGIN ... END (no DECLARE) — BEGIN here is a block opener, not a TCL
	// transaction (we are in a routine body).
	stmt := mustCreateRoutine(t, "CREATE PROCEDURE p() RETURNS INTEGER AS BEGIN RETURN 1; END")
	if stmt.Body != "BEGIN RETURN 1; END" {
		t.Errorf("Body = %q, want %q", stmt.Body, "BEGIN RETURN 1; END")
	}
	if stmt.BodyDollar {
		t.Errorf("BodyDollar = true, want false")
	}
}

func TestRoutineBody_AsBlockNestedEndIf(t *testing.T) {
	// A nested IF..END IF inside the body must NOT close the outer block early:
	// the whole block through the terminating END is the body.
	input := "CREATE PROCEDURE p() RETURNS INTEGER AS\n" +
		"BEGIN\n  IF (1 = 1) THEN\n    RETURN 1;\n  END IF;\n  RETURN 0;\nEND"
	stmt := mustCreateRoutine(t, input)
	if !strings.HasSuffix(stmt.Body, "RETURN 0;\nEND") {
		t.Errorf("Body = %q, want it to end at the terminating END (not END IF)", stmt.Body)
	}
}

func TestRoutineBody_AsBlockTrailingSemicolon(t *testing.T) {
	// The trailing ';' after the block's END is the statement delimiter and must
	// not be part of Body.
	stmt := mustCreateRoutine(t, "CREATE PROCEDURE p() RETURNS INTEGER AS BEGIN RETURN 1; END;")
	if stmt.Body != "BEGIN RETURN 1; END" {
		t.Errorf("Body = %q, want %q (no trailing ';')", stmt.Body, "BEGIN RETURN 1; END")
	}
}

// ---------------------------------------------------------------------------
// CREATE FUNCTION — language × body matrix
//
// Asserted via parseSingle (mustCreateRoutineDirect): the multi-line
// single-quoted handler bodies embed ';' / '{' that F3's Split mis-segments, so
// this matrix isolates the node's own language × body handling. Split-safe forms
// (single-line, no-';' multi-line, $$) are additionally covered end-to-end by the
// dedicated tests and the official corpus harness.
// ---------------------------------------------------------------------------

func TestCreateFunction_LanguageBodyMatrix(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		wantDollar bool
	}{
		{
			"java single-line",
			"CREATE FUNCTION f(x VARCHAR) RETURNS VARCHAR LANGUAGE JAVA HANDLER = 'T.e' AS 'class T {}'",
			false,
		},
		{
			"java multi-line single-quote",
			"CREATE FUNCTION f(x VARCHAR) RETURNS VARCHAR LANGUAGE JAVA HANDLER = 'T.e' AS\n'class T {\n  static String e(String x){ return x; }\n}'",
			false,
		},
		{
			"javascript multi-line single-quote",
			"CREATE FUNCTION f(d DOUBLE) RETURNS DOUBLE LANGUAGE JAVASCRIPT STRICT AS '\nif (D <= 0) { return 1; }\nreturn D;\n'",
			false,
		},
		{
			"python dollar",
			"CREATE FUNCTION f() RETURNS VARIANT LANGUAGE PYTHON RUNTIME_VERSION = '3.10' HANDLER = 'udf' AS $$\ndef udf(): return 1\n$$",
			true,
		},
		{
			"scala dollar runtime float",
			"CREATE FUNCTION f(x VARCHAR) RETURNS VARCHAR LANGUAGE SCALA RUNTIME_VERSION = 2.12 HANDLER='E.e' AS\n$$\nclass E { def e(x:String):String = x }\n$$",
			true,
		},
		{
			"sql single-line",
			"CREATE FUNCTION f(a NUMBER, b NUMBER) RETURNS NUMBER AS 'a * b'",
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := mustCreateRoutineDirect(t, tc.input)
			if stmt.Kind != ast.RoutineFunction {
				t.Errorf("Kind = %v, want RoutineFunction", stmt.Kind)
			}
			if stmt.BodyDollar != tc.wantDollar {
				t.Errorf("BodyDollar = %v, want %v (body=%q)", stmt.BodyDollar, tc.wantDollar, stmt.Body)
			}
			if stmt.Body == "" {
				t.Errorf("Body is empty, want a non-empty verbatim body")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CREATE FUNCTION — structure (args, RETURNS scalar/TABLE, modifiers, no-body)
// ---------------------------------------------------------------------------

func TestCreateFunction_Args(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS INT AS '1'")
		if stmt.Args != nil {
			t.Errorf("Args = %+v, want nil for ()", stmt.Args)
		}
	})
	t.Run("typed args with params", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE FUNCTION f(i NUMERIC(9, 0), s VARCHAR(100)) RETURNS NUMERIC AS '1'")
		if len(stmt.Args) != 2 {
			t.Fatalf("len(Args) = %d, want 2", len(stmt.Args))
		}
		if stmt.Args[0].Name.Normalize() != "I" || stmt.Args[0].Type.Kind != ast.TypeNumber {
			t.Errorf("Args[0] = %+v / %+v, want I NUMERIC", stmt.Args[0].Name, stmt.Args[0].Type)
		}
		if stmt.Args[1].Name.Normalize() != "S" || stmt.Args[1].Type.Kind != ast.TypeVarchar {
			t.Errorf("Args[1] = %+v / %+v, want S VARCHAR", stmt.Args[1].Name, stmt.Args[1].Type)
		}
	})
	t.Run("arg default expr", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE FUNCTION f(x INT DEFAULT 42) RETURNS INT AS '1'")
		if len(stmt.Args) != 1 || stmt.Args[0].Default == nil {
			t.Fatalf("Args = %+v, want one arg with a DEFAULT", stmt.Args)
		}
	})
}

func TestCreateFunction_ReturnsScalar(t *testing.T) {
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS FLOAT AS '1.0'")
	if stmt.ReturnType == nil || stmt.ReturnType.Kind != ast.TypeFloat {
		t.Fatalf("ReturnType = %+v, want FLOAT", stmt.ReturnType)
	}
	if stmt.ReturnTable != nil {
		t.Errorf("ReturnTable = %+v, want nil for scalar return", stmt.ReturnTable)
	}
}

func TestCreateFunction_ReturnsTable(t *testing.T) {
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS TABLE (x INTEGER, y VARCHAR) AS $$ SELECT 1, 'a' $$")
	if stmt.ReturnTable == nil {
		t.Fatalf("ReturnTable = nil, want columns")
	}
	if len(stmt.ReturnTable) != 2 {
		t.Fatalf("len(ReturnTable) = %d, want 2", len(stmt.ReturnTable))
	}
	if stmt.ReturnTable[0].Name.Normalize() != "X" || stmt.ReturnTable[0].Type.Kind != ast.TypeInt {
		t.Errorf("ReturnTable[0] = %+v, want X INTEGER", stmt.ReturnTable[0])
	}
	if stmt.ReturnType != nil {
		t.Errorf("ReturnType = %+v, want nil for TABLE return", stmt.ReturnType)
	}
}

func TestCreateFunction_ReturnsEmptyTable(t *testing.T) {
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS TABLE () AS $$ SELECT 1 $$")
	if stmt.ReturnTable == nil {
		t.Errorf("ReturnTable = nil, want a non-nil empty slice for RETURNS TABLE ()")
	}
	if len(stmt.ReturnTable) != 0 {
		t.Errorf("len(ReturnTable) = %d, want 0", len(stmt.ReturnTable))
	}
}

func TestCreateFunction_ReturnsNotNull(t *testing.T) {
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS FLOAT NOT NULL AS '1.0'")
	if !stmt.ReturnNotNull {
		t.Errorf("ReturnNotNull = false, want true")
	}
}

func TestCreateFunction_Modifiers(t *testing.T) {
	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE OR REPLACE FUNCTION f() RETURNS INT AS '1'")
		if !stmt.OrReplace {
			t.Errorf("OrReplace = false, want true")
		}
	})
	t.Run("secure", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE SECURE FUNCTION f() RETURNS INT AS '1'")
		if !stmt.Secure {
			t.Errorf("Secure = false, want true")
		}
	})
	t.Run("or replace secure", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE OR REPLACE SECURE FUNCTION f() RETURNS INT AS '1'")
		if !stmt.OrReplace || !stmt.Secure {
			t.Errorf("OrReplace/Secure = %v/%v, want true/true", stmt.OrReplace, stmt.Secure)
		}
	})
	t.Run("temp", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE TEMP FUNCTION f() RETURNS INT AS '1'")
		if !stmt.Temporary {
			t.Errorf("Temporary = false, want true")
		}
	})
	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE FUNCTION IF NOT EXISTS f() RETURNS INT AS '1'")
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
	})
	t.Run("qualified name", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE FUNCTION db.sch.f() RETURNS INT AS '1'")
		if stmt.Name.String() != "db.sch.f" {
			t.Errorf("Name = %q, want db.sch.f", stmt.Name.String())
		}
	})
}

func TestCreateFunction_VolatilityModifierWords(t *testing.T) {
	// VOLATILE / IMMUTABLE / MEMOIZABLE / STRICT and the CALLED ON NULL INPUT /
	// RETURNS NULL ON NULL INPUT phrases are captured as open-ended options.
	cases := []string{
		"CREATE FUNCTION f(x FLOAT) RETURNS FLOAT VOLATILE MEMOIZABLE AS '1'",
		"CREATE FUNCTION f(x FLOAT) RETURNS FLOAT IMMUTABLE AS '1'",
		"CREATE FUNCTION f(x VARCHAR) RETURNS VARCHAR LANGUAGE JAVA CALLED ON NULL INPUT HANDLER = 'T.e' AS 'x'",
		"CREATE FUNCTION f(x VARCHAR) RETURNS VARCHAR RETURNS NULL ON NULL INPUT AS 'x'",
		"CREATE FUNCTION f(d DOUBLE) RETURNS DOUBLE LANGUAGE JAVASCRIPT STRICT AS 'D'",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			stmt := mustCreateRoutine(t, in)
			if len(stmt.Options) == 0 {
				t.Errorf("Options is empty, want the modifier words captured")
			}
			if stmt.Body == "" {
				t.Errorf("Body is empty, want the body still parsed after the modifiers")
			}
		})
	}
}

func TestCreateFunction_NoBody(t *testing.T) {
	// A function may be defined entirely by its handler (no AS body) — docs
	// example_05 (IMPORTS-only, no AS).
	stmt := mustCreateRoutine(t, "CREATE OR REPLACE FUNCTION f(i INT) RETURNS VARIANT LANGUAGE PYTHON RUNTIME_VERSION = '3.10' HANDLER = 's.snore' IMPORTS = ('@my_stage/sleepy.py')")
	if stmt.Body != "" {
		t.Errorf("Body = %q, want empty for a handler-only function", stmt.Body)
	}
	if stmt.BodyDollar {
		t.Errorf("BodyDollar = true, want false")
	}
}

func TestCreateFunction_OptionsCaptured(t *testing.T) {
	stmt := mustCreateRoutine(t, "CREATE FUNCTION f() RETURNS VARIANT LANGUAGE PYTHON RUNTIME_VERSION = '3.10' PACKAGES = ('numpy','pandas') HANDLER = 'udf' AS $$ x $$")
	names := map[string]bool{}
	for _, o := range stmt.Options {
		names[o.Name] = true
	}
	for _, want := range []string{"LANGUAGE", "RUNTIME_VERSION", "PACKAGES", "HANDLER"} {
		if !names[want] {
			t.Errorf("missing option %q; got %v", want, names)
		}
	}
}

// ---------------------------------------------------------------------------
// CREATE EXTERNAL FUNCTION
// ---------------------------------------------------------------------------

func TestCreateExternalFunction(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE EXTERNAL FUNCTION e() RETURNS INT API_INTEGRATION = i1 AS 'https://host/path'")
		if stmt.Kind != ast.RoutineExternalFunction {
			t.Errorf("Kind = %v, want RoutineExternalFunction", stmt.Kind)
		}
		if stmt.Body != "'https://host/path'" {
			t.Errorf("Body = %q, want the quoted URL", stmt.Body)
		}
	})
	t.Run("secure external", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE OR REPLACE SECURE EXTERNAL FUNCTION e(x INT) RETURNS VARIANT API_INTEGRATION = i1 MAX_BATCH_ROWS = 100 AS 'https://host/path'")
		if stmt.Kind != ast.RoutineExternalFunction || !stmt.Secure || !stmt.OrReplace {
			t.Errorf("Kind/Secure/OrReplace = %v/%v/%v", stmt.Kind, stmt.Secure, stmt.OrReplace)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE PROCEDURE
// ---------------------------------------------------------------------------

func TestCreateProcedure(t *testing.T) {
	t.Run("javascript dollar", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE OR REPLACE PROCEDURE p() RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT AS $$ return 1; $$")
		if stmt.Kind != ast.RoutineProcedure {
			t.Errorf("Kind = %v, want RoutineProcedure", stmt.Kind)
		}
		if !stmt.ReturnNotNull {
			t.Errorf("ReturnNotNull = false, want true")
		}
	})
	t.Run("execute as owner", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE OR REPLACE PROCEDURE p(x FLOAT) RETURNS STRING LANGUAGE JAVASCRIPT STRICT EXECUTE AS OWNER AS $$ return 'x'; $$")
		if stmt.ExecuteAs != ast.ExecuteAsOwner {
			t.Errorf("ExecuteAs = %v, want ExecuteAsOwner", stmt.ExecuteAs)
		}
	})
	t.Run("execute as caller after options", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE OR REPLACE PROCEDURE p(a NUMBER) RETURNS NUMBER LANGUAGE PYTHON HANDLER='main' RUNTIME_VERSION=3.10 PACKAGES = ('snowflake-snowpark-python') EXECUTE AS CALLER AS $$ x $$")
		if stmt.ExecuteAs != ast.ExecuteAsCaller {
			t.Errorf("ExecuteAs = %v, want ExecuteAsCaller", stmt.ExecuteAs)
		}
	})
	t.Run("secure procedure single-quote body", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE SECURE PROCEDURE p() RETURNS STRING LANGUAGE JAVASCRIPT AS ' '")
		if !stmt.Secure {
			t.Errorf("Secure = false, want true")
		}
		if stmt.Body != "' '" {
			t.Errorf("Body = %q, want %q", stmt.Body, "' '")
		}
	})
	t.Run("sql language single-quote", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "create procedure p() returns int language sql as ' '")
		if stmt.Kind != ast.RoutineProcedure || stmt.Body != "' '" {
			t.Errorf("Kind/Body = %v/%q", stmt.Kind, stmt.Body)
		}
	})
	t.Run("secrets quoted-key group", func(t *testing.T) {
		// SECRETS = ('secret_variable_name' = secret_name [, ...]) — a group whose
		// entry KEY is a quoted string (official create-procedure example_06). Each
		// entry binds a string-name to a secret identifier.
		stmt := mustCreateRoutine(t, "CREATE OR ALTER PROCEDURE python_add1(A NUMBER) RETURNS NUMBER LANGUAGE PYTHON HANDLER='main' RUNTIME_VERSION=3.10 EXTERNAL_ACCESS_INTEGRATIONS=(example_integration) secrets=('secret_variable_name'=secret_name) PACKAGES = ('snowflake-snowpark-python') EXECUTE AS CALLER AS $$ x $$")
		var secrets *ast.CopyOption
		for _, o := range stmt.Options {
			if o.Name == "SECRETS" {
				secrets = o
			}
		}
		if secrets == nil {
			t.Fatalf("SECRETS option not captured; got %v", optNames(stmt.Options))
		}
		if len(secrets.Group) != 1 {
			t.Fatalf("SECRETS group len = %d, want 1", len(secrets.Group))
		}
		if got := secrets.Group[0].Name; got != "secret_variable_name" {
			t.Errorf("SECRETS entry key = %q, want %q", got, "secret_variable_name")
		}
	})
	t.Run("secrets multi quoted-key", func(t *testing.T) {
		stmt := mustCreateRoutine(t, "CREATE PROCEDURE p() RETURNS NUMBER LANGUAGE PYTHON HANDLER='main' RUNTIME_VERSION=3.10 SECRETS=('a'=s1,'b'=s2) PACKAGES=('x') AS $$ y $$")
		var secrets *ast.CopyOption
		for _, o := range stmt.Options {
			if o.Name == "SECRETS" {
				secrets = o
			}
		}
		if secrets == nil || len(secrets.Group) != 2 {
			t.Fatalf("SECRETS group = %#v, want 2 entries", secrets)
		}
		if secrets.Group[0].Name != "a" || secrets.Group[1].Name != "b" {
			t.Errorf("SECRETS keys = %q/%q, want a/b", secrets.Group[0].Name, secrets.Group[1].Name)
		}
	})
}

func optNames(opts []*ast.CopyOption) []string {
	names := make([]string, 0, len(opts))
	for _, o := range opts {
		names = append(names, o.Name)
	}
	return names
}

// ---------------------------------------------------------------------------
// ALTER FUNCTION / PROCEDURE
// ---------------------------------------------------------------------------

func TestAlterFunction(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f() RENAME TO f2")
		if stmt.Procedure {
			t.Errorf("Procedure = true, want false")
		}
		if stmt.Action != ast.AlterRoutineRename || stmt.NewName.String() != "f2" {
			t.Errorf("Action/NewName = %v/%v", stmt.Action, stmt.NewName)
		}
	})
	t.Run("rename with argtypes", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f(NUMBER, VARCHAR) RENAME TO f2")
		if len(stmt.ArgTypes) != 2 {
			t.Fatalf("len(ArgTypes) = %d, want 2", len(stmt.ArgTypes))
		}
		if stmt.ArgTypes[0].Kind != ast.TypeNumber || stmt.ArgTypes[1].Kind != ast.TypeVarchar {
			t.Errorf("ArgTypes = %+v, want NUMBER, VARCHAR", stmt.ArgTypes)
		}
	})
	t.Run("if exists", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION IF EXISTS f() RENAME TO f2")
		if !stmt.IfExists {
			t.Errorf("IfExists = false, want true")
		}
	})
	t.Run("set secure", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f() SET SECURE")
		if stmt.Action != ast.AlterRoutineSet || len(stmt.Options) != 1 || stmt.Options[0].Name != "SECURE" {
			t.Errorf("Action/Options = %v/%+v, want SET [SECURE]", stmt.Action, stmt.Options)
		}
	})
	t.Run("set comment", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f() SET COMMENT = 'hi'")
		if stmt.Action != ast.AlterRoutineSet || stmt.Options[0].Name != "COMMENT" {
			t.Errorf("Options = %+v, want [COMMENT='hi']", stmt.Options)
		}
	})
	t.Run("set log_level", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f(INT) SET LOG_LEVEL = 'DEBUG'")
		if stmt.Options[0].Name != "LOG_LEVEL" {
			t.Errorf("Options = %+v, want LOG_LEVEL", stmt.Options)
		}
	})
	t.Run("unset secure", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f() UNSET SECURE")
		if stmt.Action != ast.AlterRoutineUnset || len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "SECURE" {
			t.Errorf("Action/UnsetProps = %v/%v", stmt.Action, stmt.UnsetProps)
		}
	})
	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f() UNSET COMMENT")
		if stmt.UnsetProps[0] != "COMMENT" {
			t.Errorf("UnsetProps = %v, want [COMMENT]", stmt.UnsetProps)
		}
	})
	t.Run("set api_integration external", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER FUNCTION f(INT) SET API_INTEGRATION = i2")
		if stmt.Options[0].Name != "API_INTEGRATION" {
			t.Errorf("Options = %+v, want API_INTEGRATION", stmt.Options)
		}
	})
}

func TestAlterProcedure(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER PROCEDURE pr() RENAME TO p2")
		if !stmt.Procedure {
			t.Errorf("Procedure = false, want true")
		}
		if stmt.Action != ast.AlterRoutineRename || stmt.NewName.String() != "p2" {
			t.Errorf("Action/NewName = %v/%v", stmt.Action, stmt.NewName)
		}
	})
	t.Run("execute as caller", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER PROCEDURE p(VARCHAR) EXECUTE AS CALLER")
		if stmt.Action != ast.AlterRoutineExecuteAs || stmt.ExecuteAs != ast.ExecuteAsCaller {
			t.Errorf("Action/ExecuteAs = %v/%v", stmt.Action, stmt.ExecuteAs)
		}
	})
	t.Run("execute as owner", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER PROCEDURE p() EXECUTE AS OWNER")
		if stmt.ExecuteAs != ast.ExecuteAsOwner {
			t.Errorf("ExecuteAs = %v, want ExecuteAsOwner", stmt.ExecuteAs)
		}
	})
	t.Run("set comment", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER PROCEDURE p(INT, INT) SET COMMENT = 'x'")
		if stmt.Action != ast.AlterRoutineSet || stmt.Options[0].Name != "COMMENT" {
			t.Errorf("Options = %+v", stmt.Options)
		}
	})
	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterRoutine(t, "ALTER PROCEDURE IF EXISTS p() UNSET COMMENT")
		if !stmt.IfExists || stmt.UnsetProps[0] != "COMMENT" {
			t.Errorf("IfExists/UnsetProps = %v/%v", stmt.IfExists, stmt.UnsetProps)
		}
	})
}

// ---------------------------------------------------------------------------
// Negative tests — omni must reject malformed routine DDL
// ---------------------------------------------------------------------------

func TestRoutine_Negatives(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"function missing parens", "CREATE FUNCTION f RETURNS INT AS '1'"},
		{"function missing returns", "CREATE FUNCTION f() AS '1'"},
		{"function returns missing type", "CREATE FUNCTION f() RETURNS AS '1'"},
		{"function as with no body", "CREATE FUNCTION f() RETURNS INT AS"},
		{"function unterminated single-quote body", "CREATE FUNCTION f() RETURNS INT AS 'abc"},
		{"function unterminated dollar body", "CREATE FUNCTION f() RETURNS INT AS $$abc"},
		{"function arg missing type", "CREATE FUNCTION f(x) RETURNS INT AS '1'"},
		{"alter function missing signature parens", "ALTER FUNCTION f RENAME TO g"},
		{"alter function bad action", "ALTER FUNCTION f() FROBNICATE"},
		{"alter procedure set nothing", "ALTER PROCEDURE p() SET"},
		{"external function not function", "CREATE EXTERNAL TABLE e() RETURNS INT AS 'x'"},
		{"procedure execute as bad", "CREATE PROCEDURE p() RETURNS INT LANGUAGE SQL EXECUTE AS NOBODY AS '1'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertParseError(t, tc.input)
		})
	}
}

// ---------------------------------------------------------------------------
// Loc spans
// ---------------------------------------------------------------------------

func TestRoutine_LocSpansBody(t *testing.T) {
	input := "CREATE FUNCTION f() RETURNS INT AS '1 + 1'"
	stmt := mustCreateRoutine(t, input)
	// Loc.End must reach the closing quote of the body (the last byte of input).
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d (end of body)", stmt.Loc.End, len(input))
	}
	if input[stmt.Loc.Start:stmt.Loc.End] != input {
		t.Errorf("Loc span = %q, want the full statement", input[stmt.Loc.Start:stmt.Loc.End])
	}
}

func TestRoutine_LocMultiLineBody(t *testing.T) {
	input := "CREATE FUNCTION f() RETURNS VARCHAR AS 'line1\nline2'"
	stmt := mustCreateRoutine(t, input)
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d (end of multi-line body)", stmt.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// Official docs corpus — every CREATE FUNCTION / PROCEDURE example in the
// create-function and create-procedure corpora must parse with zero errors. The
// official docs are the authoritative oracle (truth1); the legacy .g4 corpus is
// the regression baseline (truth2).
//
// Each corpus file is exactly ONE statement, so it is driven through parseSingle
// on the whole file text (trailing ';' stripped) rather than ParseBestEffort.
// This is deliberate: the multi-line single-quoted handler bodies (example_01 /
// _03, etc.) embed a ';' that F3's lexer-driven Split mis-segments — a Split-layer
// limitation (the omni lexer's scanString cannot span a newline) tracked as a
// flagged divergence, NOT this node's concern. parseSingle exercises THIS node's
// body raw-scan against the real documented forms, including those multi-line
// single-quoted bodies, which all parse. The only skips are CREATE OR ALTER …
// (create-function example_15/16, create-procedure example_05/06), owned by the
// separate parser-or-alter DAG node, and the SELECT FROM TABLE(...) line
// (create-function example_12), owned by the select-core node.
// ---------------------------------------------------------------------------

var routineCorpusDirs = []string{
	"testdata/official/create-function",
	"testdata/official/create-procedure",
}

func TestRoutine_OfficialCorpus(t *testing.T) {
	for _, dir := range routineCorpusDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertRoutineFileParses(t, string(data))
			})
		}
	}
}

func TestRoutine_LegacyCorpus(t *testing.T) {
	// The legacy create_function.sql / create_procedure.sql carry the regression
	// cases; each line is a single, ';'-terminated statement with a Split-safe
	// $$ or short single-quoted body, so they are driven statement-by-statement
	// via Split + parseSingle.
	for _, path := range []string{"testdata/legacy/create_function.sql", "testdata/legacy/create_procedure.sql"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		t.Run(path, func(t *testing.T) {
			assertRoutineStatementsParse(t, string(data))
		})
	}
}

// assertRoutineFileParses treats a single-statement corpus file as ONE
// statement: it strips a trailing ';' and parses the whole text via parseSingle
// (bypassing F3's Split, which mis-segments multi-line single-quoted bodies). It
// asserts the statement parses with no errors and to a routine node. CREATE OR
// ALTER FUNCTION/PROCEDURE is a routine and is asserted; a non-routine statement
// owned by another DAG node (e.g. a SELECT) is skipped via routineStmtKind.
func assertRoutineFileParses(t *testing.T, sql string) {
	t.Helper()
	text := strings.TrimRight(strings.TrimSpace(sql), ";")
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	upper := strings.ToUpper(text)

	kind, want := routineStmtKind(upper)
	if kind == "" {
		return // context statement owned by another DAG node (e.g. SELECT)
	}

	node, errs := parseSingle(text, 0)
	if len(errs) > 0 {
		t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
		return
	}
	if !want(node) {
		t.Errorf("statement %q (%s) parsed to unexpected type %T", text, kind, node)
	}
}

// assertRoutineStatementsParse parses sql via Split + parseSingle and asserts
// that every CREATE FUNCTION / EXTERNAL FUNCTION / PROCEDURE and ALTER FUNCTION /
// PROCEDURE statement parses with no errors and to the expected AST type.
// Statements owned by other DAG nodes (a SELECT, a CREATE OR ALTER) are skipped.
// Used for the legacy corpus, whose bodies are Split-safe.
func assertRoutineStatementsParse(t *testing.T, sql string) {
	t.Helper()
	for _, seg := range Split(sql) {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)

		kind, want := routineStmtKind(upper)
		if kind == "" {
			continue // context statement owned by another DAG node (e.g. SELECT)
		}

		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if len(errs) > 0 {
			t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
			continue
		}
		if !want(node) {
			t.Errorf("statement %q (%s) parsed to unexpected type %T", text, kind, node)
		}
	}
}

// routineStmtKind classifies an uppercased statement by its leading keywords,
// returning a kind label and a predicate checking the parsed node type. Returns
// ("", nil) for statements this node does not own.
func routineStmtKind(upper string) (string, func(ast.Node) bool) {
	isCreate := routineHasCreatePrefix(upper)
	switch {
	case isCreate:
		return "CREATE-ROUTINE", func(n ast.Node) bool { _, ok := n.(*ast.CreateRoutineStmt); return ok }
	case strings.HasPrefix(upper, "ALTER FUNCTION "), strings.HasPrefix(upper, "ALTER PROCEDURE "):
		return "ALTER-ROUTINE", func(n ast.Node) bool { _, ok := n.(*ast.AlterRoutineStmt); return ok }
	}
	return "", nil
}

// routineHasCreatePrefix reports whether an uppercased statement is a CREATE
// FUNCTION / EXTERNAL FUNCTION / PROCEDURE, allowing for the OR REPLACE / SECURE
// / TEMP / TEMPORARY modifiers between CREATE and the object keyword.
func routineHasCreatePrefix(upper string) bool {
	if !strings.HasPrefix(upper, "CREATE ") {
		return false
	}
	rest := strings.TrimPrefix(upper, "CREATE ")
	for _, mod := range []string{"OR REPLACE ", "OR ALTER ", "SECURE ", "TEMPORARY ", "TEMP ", "EXTERNAL "} {
		rest = strings.TrimPrefix(rest, mod)
	}
	return strings.HasPrefix(rest, "FUNCTION ") || strings.HasPrefix(rest, "FUNCTION(") ||
		strings.HasPrefix(rest, "PROCEDURE ") || strings.HasPrefix(rest, "PROCEDURE(")
}
