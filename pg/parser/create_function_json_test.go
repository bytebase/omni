package parser

import "testing"

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
	sqls := []string{
		`CREATE FUNCTION f(data json) RETURNS int AS 'select 1' LANGUAGE sql`,
		`CREATE FUNCTION f(x dec(10,2)) RETURNS int AS 'select 1' LANGUAGE sql`,
		`CREATE FUNCTION f(name national character(10)) RETURNS int AS 'select 1' LANGUAGE sql`,
		`CREATE FUNCTION f(name nchar(10)) RETURNS int AS 'select 1' LANGUAGE sql`,
	}
	for _, sql := range sqls {
		_, err := Parse(sql)
		if err != nil {
			t.Errorf("parse failed for %q: %v", sql, err)
		}
	}
}
