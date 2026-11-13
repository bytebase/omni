package parser

import "testing"

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
