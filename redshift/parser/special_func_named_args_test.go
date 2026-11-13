package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestSpecialFuncNamedArgs: OVERLAY and SUBSTRING also take
// '(' func_arg_list_opt ')', so a user function of that name is called with
// named arguments, in a target list and in FROM.
func TestSpecialFuncNamedArgs(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM overlay(foo => 'abc')",
		"SELECT overlay(foo => 'abc', bar := 1)",
		"SELECT * FROM substring(s => 'abc', n => 2)",
		"SELECT substring(s := 'abc')",
		"SELECT overlay('abc' PLACING 'x' FROM 2)",
		"SELECT substring('abc' FROM 2 FOR 1)",
		"SELECT * FROM overlay('abc', 'x', 2)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{"SELECT overlay(foo => )", "SELECT * FROM substring(s => 'a',)"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

// TestSubstringPlainCall: an empty or delimiter-free SUBSTRING list is
// gram.y's func_arg_list_opt form, an unqualified explicit call.
func TestSubstringPlainCall(t *testing.T) {
	for _, sql := range []string{"SELECT * FROM substring('abc')", "SELECT substring()", "SELECT substring('a', 2)"} {
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		var fc *nodes.FuncCall
		nodes.Inspect(stmts, func(n nodes.Node) bool {
			if f, ok := n.(*nodes.FuncCall); ok && fc == nil {
				fc = f
			}
			return true
		})
		if fc == nil || len(fc.Funcname.Items) != 1 || fc.FuncFormat != int(nodes.COERCE_EXPLICIT_CALL) {
			t.Errorf("Parse(%q) built %+v, want an unqualified explicit substring call", sql, fc)
		}
	}
}
