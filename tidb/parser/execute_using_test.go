package parser

import "testing"

// TestExecuteUsingUserVariables holds EXECUTE ... USING to UserVariableList.
func TestExecuteUsingUserVariables(t *testing.T) {
	for _, sql := range []string{
		"EXECUTE s",
		"EXECUTE s USING @a",
		"EXECUTE s USING @a, @`b`",
		"EXECUTE s USING @''",
		"EXECUTE s USING @``",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"EXECUTE s USING a",
		"EXECUTE s USING @a, b",
		"EXECUTE s USING @a + 1",
		"EXECUTE s USING 1",
		"EXECUTE s USING @@sql_mode",
		"EXECUTE s USING `@a`",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): want a syntax error", sql)
		}
	}
}
