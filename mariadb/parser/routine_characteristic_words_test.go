package parser

import "testing"

// TestRoutineCharacteristicWordsWhole holds sp_chistic's data-access words to
// the grammar: CONTAINS SQL, READS SQL DATA and MODIFIES SQL DATA are read
// whole, and SQL starts only SQL SECURITY.
func TestRoutineCharacteristicWordsWhole(t *testing.T) {
	for _, sql := range []string{
		"ALTER PROCEDURE p CONTAINS SQL",
		"ALTER PROCEDURE p READS SQL DATA",
		"ALTER FUNCTION f MODIFIES SQL DATA SQL SECURITY INVOKER",
		"ALTER FUNCTION f LANGUAGE SQL NO SQL",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"ALTER PROCEDURE p CONTAINS",
		"ALTER PROCEDURE p READS SQL",
		"ALTER PROCEDURE p READS DATA",
		"ALTER FUNCTION f MODIFIES SQL",
		"ALTER FUNCTION f SQL",
		"ALTER FUNCTION f SQL COMMENT 'x'",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): want a syntax error", sql)
		}
	}
}
