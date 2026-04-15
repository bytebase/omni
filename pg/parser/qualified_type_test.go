package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParseQualifiedTypeInFuncTypePositions covers the regression
// surface for the parseFuncType partial-snapshot bug. Before the fix,
// `parseFuncType`'s speculative %TYPE branch saved cur/prev/nextBuf/
// hasNext/lexer.Err but NOT lexer.pos/start/state, so any qualified
// type at any parseFuncType call site corrupted the lexer position on
// rollback and produced a downstream syntax error.
//
// parseFuncType is called from FOUR sites; this test covers all of them:
//   - create_function.go:104  (function return type)
//   - create_function.go:291  (function parameter type)
//   - create_function.go:359  (RETURNS TABLE column type)
//   - define.go:329           (parseDefArg, e.g. CREATE OPERATOR LEFTARG/RIGHTARG)
//
// 3-component qualified names (e.g. db.schema.mytype) are deliberately
// NOT tested here. They fail in non-func_type contexts too (CREATE TABLE,
// CAST, ALTER TABLE), so they are a separate parseGenericType limitation
// that is out of scope for this fix.
func TestParseQualifiedTypeInFuncTypePositions(t *testing.T) {
	cases := []struct {
		name      string
		sql       string
		paramType []string // expected ArgType.Names for the first parameter (or RETURNS table column)
	}{
		// parseFuncArg argType — function parameter
		{
			name:      "function param qualified pg_catalog.int4",
			sql:       `CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramType: []string{"pg_catalog", "int4"},
		},
		{
			name:      "function param custom schema",
			sql:       `CREATE FUNCTION f(x schema.mytype) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramType: []string{"schema", "mytype"},
		},
		// sanity: unqualified parameter still works
		{
			name:      "function param unqualified int (sanity)",
			sql:       `CREATE FUNCTION f(x int) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramType: []string{"pg_catalog", "int4"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			cf := firstCreateFunctionFromStmts(t, stmts)
			if cf.Parameters == nil || len(cf.Parameters.Items) == 0 {
				t.Fatalf("expected at least 1 parameter, got %v", cf.Parameters)
			}
			fp, ok := cf.Parameters.Items[0].(*nodes.FunctionParameter)
			if !ok {
				t.Fatalf("expected FunctionParameter, got %T", cf.Parameters.Items[0])
			}
			assertTypeNameNames(t, fp.ArgType, tc.paramType)
		})
	}
}

// TestParseFunctionReturnsQualifiedType covers the function return type
// position (create_function.go:104). Before the fix, RETURNS pg_catalog.int4
// would corrupt the lexer state and fail.
func TestParseFunctionReturnsQualifiedType(t *testing.T) {
	cases := []struct {
		name       string
		sql        string
		returnType []string
		setof      bool
	}{
		{
			name:       "RETURNS pg_catalog.int4",
			sql:        `CREATE FUNCTION f() RETURNS pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
			returnType: []string{"pg_catalog", "int4"},
		},
		{
			name:       "RETURNS schema.mytype",
			sql:        `CREATE FUNCTION f() RETURNS schema.mytype AS 'select 1' LANGUAGE sql`,
			returnType: []string{"schema", "mytype"},
		},
		// sanity: SETOF qualified — bypasses the speculative branch via the SETOF prefix
		{
			name:       "RETURNS SETOF pg_catalog.int4 (sanity)",
			sql:        `CREATE FUNCTION f() RETURNS SETOF pg_catalog.int4 AS 'select 1' LANGUAGE sql`,
			returnType: []string{"pg_catalog", "int4"},
			setof:      true,
		},
		// sanity: unqualified return still works
		{
			name:       "RETURNS int (sanity)",
			sql:        `CREATE FUNCTION f() RETURNS int AS 'select 1' LANGUAGE sql`,
			returnType: []string{"pg_catalog", "int4"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			cf := firstCreateFunctionFromStmts(t, stmts)
			if cf.ReturnType == nil {
				t.Fatalf("expected non-nil ReturnType")
			}
			assertTypeNameNames(t, cf.ReturnType, tc.returnType)
			if cf.ReturnType.Setof != tc.setof {
				t.Errorf("Setof: got %v, want %v", cf.ReturnType.Setof, tc.setof)
			}
		})
	}
}

// TestParseFunctionReturnsTypePctType covers the success path of the
// speculative branch (the %TYPE form) — which must keep working after
// the rollback fix.
func TestParseFunctionReturnsTypePctType(t *testing.T) {
	sql := `CREATE FUNCTION f() RETURNS schema.tab.col%TYPE AS 'select 1' LANGUAGE sql`
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	cf := firstCreateFunctionFromStmts(t, stmts)
	if cf.ReturnType == nil {
		t.Fatalf("expected non-nil ReturnType")
	}
	if !cf.ReturnType.PctType {
		t.Errorf("expected PctType=true on %%TYPE form, got false")
	}
	assertTypeNameNames(t, cf.ReturnType, []string{"schema", "tab", "col"})
}

// TestParseReturnsTableQualifiedColumnType covers parseTableFuncColumn
// (create_function.go:359) — the RETURNS TABLE (col type) path that
// also goes through parseFuncType.
func TestParseReturnsTableQualifiedColumnType(t *testing.T) {
	sql := `CREATE FUNCTION f() RETURNS TABLE (x pg_catalog.int4, y schema.mytype) AS 'select 1' LANGUAGE sql`
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	cf := firstCreateFunctionFromStmts(t, stmts)
	if cf.Parameters == nil || len(cf.Parameters.Items) != 2 {
		t.Fatalf("expected 2 RETURNS TABLE columns, got %v", cf.Parameters)
	}
	col1, ok := cf.Parameters.Items[0].(*nodes.FunctionParameter)
	if !ok {
		t.Fatalf("expected FunctionParameter for column 1, got %T", cf.Parameters.Items[0])
	}
	if col1.Name != "x" {
		t.Errorf("column 1 name: got %q, want %q", col1.Name, "x")
	}
	assertTypeNameNames(t, col1.ArgType, []string{"pg_catalog", "int4"})

	col2, ok := cf.Parameters.Items[1].(*nodes.FunctionParameter)
	if !ok {
		t.Fatalf("expected FunctionParameter for column 2, got %T", cf.Parameters.Items[1])
	}
	if col2.Name != "y" {
		t.Errorf("column 2 name: got %q, want %q", col2.Name, "y")
	}
	assertTypeNameNames(t, col2.ArgType, []string{"schema", "mytype"})
}

// TestParseCreateOperatorQualifiedArg covers parseDefArg (define.go:329)
// — the CREATE OPERATOR LEFTARG/RIGHTARG path that also goes through
// parseFuncType.
func TestParseCreateOperatorQualifiedArg(t *testing.T) {
	cases := []string{
		`CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = pg_catalog.int4, RIGHTARG = int4)`,
		`CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = int4, RIGHTARG = pg_catalog.int4)`,
		`CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = pg_catalog.int4, RIGHTARG = pg_catalog.int4)`,
		// sanity: unqualified still works
		`CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = int4, RIGHTARG = int4)`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := Parse(sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// firstCreateFunctionFromStmts extracts the first CreateFunctionStmt from
// the parser's RawStmt-wrapped output. Same shape as the helper in
// create_function_json_test.go but kept separate to avoid coupling the
// two test files.
func firstCreateFunctionFromStmts(t *testing.T, stmts *nodes.List) *nodes.CreateFunctionStmt {
	t.Helper()
	if stmts == nil || len(stmts.Items) == 0 {
		t.Fatalf("expected at least one statement, got nil/empty")
	}
	raw, ok := stmts.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt wrapper, got %T", stmts.Items[0])
	}
	cf, ok := raw.Stmt.(*nodes.CreateFunctionStmt)
	if !ok {
		t.Fatalf("expected CreateFunctionStmt inside RawStmt, got %T", raw.Stmt)
	}
	return cf
}

// assertTypeNameNames walks ArgType.Names and asserts it matches the
// expected list of name components. Each component must be a *nodes.String.
func assertTypeNameNames(t *testing.T, tn *nodes.TypeName, want []string) {
	t.Helper()
	if tn == nil {
		t.Fatalf("expected non-nil TypeName")
	}
	if tn.Names == nil {
		t.Fatalf("expected non-nil TypeName.Names")
	}
	if len(tn.Names.Items) != len(want) {
		var got []string
		for _, item := range tn.Names.Items {
			if s, ok := item.(*nodes.String); ok {
				got = append(got, s.Str)
			} else {
				got = append(got, "<?>")
			}
		}
		t.Fatalf("Names length: got %d %v, want %d %v", len(tn.Names.Items), got, len(want), want)
	}
	for i, item := range tn.Names.Items {
		s, ok := item.(*nodes.String)
		if !ok {
			t.Errorf("Names[%d] type: got %T, want *nodes.String", i, item)
			continue
		}
		if s.Str != want[i] {
			t.Errorf("Names[%d]: got %q, want %q", i, s.Str, want[i])
		}
	}
}
