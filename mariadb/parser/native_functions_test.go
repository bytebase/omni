package parser

import "testing"

// nativeFunctionForms are the argument counts that take a step, and a name that is no valid
// text, beside what TestParseRefusalForms holds.
var nativeFunctionForms = []struct {
	sql    string
	accept bool
}{
	{"SELECT JSON_SET(1,1,1)", true},
	{"SELECT JSON_SET(1,1,1,1)", false},
	{"SELECT JSON_SET(1,1,1,1,1)", true},
	{"SELECT JSON_OBJECT()", true},
	{"SELECT JSON_OBJECT(1)", false},
	{"SELECT JSON_OBJECT(1,2)", true},
	{"SELECT JSON_OBJECT(1,2,3)", false},
	{"SELECT `abs`(1,2)", false},
	{"SELECT d.abs(1,2)", true},
	{"SELECT 1 AS `\xff`", false},
}

func TestNativeFunctionForms(t *testing.T) {
	for _, c := range nativeFunctionForms {
		if _, err := Parse(c.sql); (err == nil) != c.accept {
			t.Errorf("Parse(%q) error = %v, want accept=%v", c.sql, err, c.accept)
		}
	}
}

func TestNativeFunctionFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range nativeFunctionForms {
		if got := o.prepareParses(t, c.sql); got != c.accept {
			t.Errorf("MariaDB parses %q: %v, this file records %v", c.sql, got, c.accept)
		}
	}
}
