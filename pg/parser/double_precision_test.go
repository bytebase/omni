package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParseCreateFunctionDoublePrecisionAllArgClasses covers the regression
// surface for the parseFuncArg pushBack bug. Before the fix, pushBack
// rewrote the rolled-back token's Type to IDENT unconditionally on
// rollback, so DOUBLE_P (the only UnreservedKeyword in
// simpleTypenameLeadTokens) lost its type and broke every parseFuncArg
// alternative ending with `double precision`.
//
// The peek-then-commit rewrite eliminates the speculative consume
// entirely; this test locks down all 6 visible SQL forms identified by
// the codex review of the plan doc.
func TestParseCreateFunctionDoublePrecisionAllArgClasses(t *testing.T) {
	// Note: omni's parser uses FUNC_PARAM_IN as the implicit default mode
	// when no arg_class is specified (see create_function.go parseFuncArg
	// `mode := nodes.FUNC_PARAM_IN`). PG itself uses FUNC_PARAM_DEFAULT
	// for the same case, but omni's normalization is FUNC_PARAM_IN.
	cases := []struct {
		name      string
		sql       string
		paramName string                      // "" for paramless
		mode      nodes.FunctionParameterMode // expected param mode
		typeNames []string                    // expected ArgType.Names
	}{
		// case 5: bare func_type
		{
			name:      "case 5 paramless double precision",
			sql:       `CREATE FUNCTION f(double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "float8"},
		},
		// case 4: arg_class func_type (4 variants)
		{
			name:      "case 4 IN double precision",
			sql:       `CREATE FUNCTION f(IN double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "float8"},
		},
		{
			name:      "case 4 OUT double precision",
			sql:       `CREATE FUNCTION f(OUT double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_OUT,
			typeNames: []string{"pg_catalog", "float8"},
		},
		{
			name:      "case 4 INOUT double precision",
			sql:       `CREATE FUNCTION f(INOUT double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_INOUT,
			typeNames: []string{"pg_catalog", "float8"},
		},
		{
			name:      "case 4 VARIADIC double precision",
			sql:       `CREATE FUNCTION f(VARIADIC double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_VARIADIC,
			typeNames: []string{"pg_catalog", "float8"},
		},
		// case 5 with default value
		{
			name:      "case 5 double precision DEFAULT 1.0",
			sql:       `CREATE FUNCTION f(double precision DEFAULT 1.0) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "float8"},
		},
		// sanity: case 3 (param_name func_type) with explicit name (must keep working)
		{
			name:      "case 3 sanity: explicit param name + double precision",
			sql:       `CREATE FUNCTION f(p double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "p",
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "float8"},
		},
		// sanity: case 1 (arg_class param_name func_type) with explicit name
		{
			name:      "case 1 sanity: IN p double precision",
			sql:       `CREATE FUNCTION f(IN p double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "p",
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "float8"},
		},
		// sanity: case 2 (param_name arg_class func_type) — Oracle-style
		{
			name:      "case 2 sanity: p IN double precision",
			sql:       `CREATE FUNCTION f(p IN double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "p",
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "float8"},
		},
		// sanity: paramless single-token type (must keep working)
		{
			name:      "case 5 sanity: paramless int",
			sql:       `CREATE FUNCTION f(int) RETURNS int AS 'select 1' LANGUAGE sql`,
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "int4"},
		},
		// cross-bug interaction: case 3 with qualified type
		// (covers parseFuncArg + parseFuncType together)
		{
			name:      "cross sanity: p pg_catalog.int4",
			sql:       `CREATE FUNCTION f(p pg_catalog.int4) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "p",
			mode:      nodes.FUNC_PARAM_IN,
			typeNames: []string{"pg_catalog", "int4"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			cf := firstCreateFunctionFromStmts(t, stmts)
			if cf.Parameters == nil || len(cf.Parameters.Items) != 1 {
				t.Fatalf("expected exactly 1 parameter, got %v", cf.Parameters)
			}
			fp, ok := cf.Parameters.Items[0].(*nodes.FunctionParameter)
			if !ok {
				t.Fatalf("expected FunctionParameter, got %T", cf.Parameters.Items[0])
			}
			if fp.Name != tc.paramName {
				t.Errorf("Name: got %q, want %q", fp.Name, tc.paramName)
			}
			if fp.Mode != tc.mode {
				t.Errorf("Mode: got %v, want %v", fp.Mode, tc.mode)
			}
			assertTypeNameNames(t, fp.ArgType, tc.typeNames)
		})
	}
}
