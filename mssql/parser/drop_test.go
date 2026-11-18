package parser

import "testing"

// TestDropRequiresAName pins that `DROP <object>` names its object: a bare
// keyword and a trailing comma are both syntax errors to SQL Server
// ("Incorrect syntax near the end of the statement").
func TestDropRequiresAName(t *testing.T) {
	for _, sql := range []string{
		"DROP TABLE",
		"DROP VIEW",
		"DROP PROCEDURE",
		"DROP FUNCTION",
		"DROP SCHEMA",
		"DROP DEFAULT",
		"DROP RULE",
		"DROP SYNONYM",
		"DROP SEQUENCE",
		"DROP TABLE a,",
		"DROP TABLE a, , b",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
	for _, sql := range []string{
		"DROP TABLE a",
		"DROP TABLE a, b",
		"DROP TABLE IF EXISTS a",
		"DROP VIEW a.b",
		"DROP DATABASE a",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}
