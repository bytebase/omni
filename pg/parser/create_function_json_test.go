package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestParseCreateFunctionWithColNameKeywordParam(t *testing.T) {
	// Regression cases for the isBuiltinType whitelist drift. These
	// keywords are all in parseSimpleTypename but were missing from
	// create_function.go:isBuiltinType whitelist in the original tree,
	// so parseFuncArg would fail to recognize them as type starts after
	// speculatively consuming the param_name.
	//
	// BOOLEAN is intentionally NOT in this list: BOOLEAN_P is already
	// present in the pre-refactor isBuiltinType whitelist, so it parses
	// correctly today. Do not add it as a regression case.
	cases := []struct {
		sql       string
		paramName string
		typeLast  string // last name component of the parsed TypeName
	}{
		{
			sql:       `CREATE FUNCTION f(data json) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "data",
			typeLast:  "json",
		},
		{
			sql:       `CREATE FUNCTION f(x dec(10,2)) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "x",
			typeLast:  "numeric",
		},
		{
			sql:       `CREATE FUNCTION f(name national character(10)) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "name",
			typeLast:  "bpchar",
		},
		{
			sql:       `CREATE FUNCTION f(name nchar(10)) RETURNS int AS 'select 1' LANGUAGE sql`,
			paramName: "name",
			typeLast:  "bpchar",
		},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			// Walk to FunctionParameter and assert AST shape. A silent
			// regression where parseFuncArg drops the param name into the
			// type would pass an error-only check but fail this assertion.
			cf := firstCreateFunction(t, stmts)
			if cf.Parameters == nil || len(cf.Parameters.Items) != 1 {
				t.Fatalf("expected exactly 1 parameter, got %v", cf.Parameters)
			}
			fp, ok := cf.Parameters.Items[0].(*nodes.FunctionParameter)
			if !ok {
				t.Fatalf("expected FunctionParameter, got %T", cf.Parameters.Items[0])
			}
			if fp.Name != tc.paramName {
				t.Errorf("param name: got %q, want %q", fp.Name, tc.paramName)
			}
			if fp.ArgType == nil || fp.ArgType.Names == nil || len(fp.ArgType.Names.Items) == 0 {
				t.Fatalf("missing ArgType.Names: %+v", fp.ArgType)
			}
			last := fp.ArgType.Names.Items[len(fp.ArgType.Names.Items)-1]
			lastStr, ok := last.(*nodes.String)
			if !ok {
				t.Fatalf("expected last name to be *nodes.String, got %T", last)
			}
			if lastStr.Str != tc.typeLast {
				t.Errorf("type last name: got %q, want %q", lastStr.Str, tc.typeLast)
			}
		})
	}
}

func firstCreateFunction(t *testing.T, stmts *nodes.List) *nodes.CreateFunctionStmt {
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
