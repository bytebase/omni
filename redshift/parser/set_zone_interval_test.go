package parser

import (
	"errors"
	"testing"
)

// TestSetZoneInterval holds zone_value's interval forms and SET TRANSACTION
// SNAPSHOT to gram.y.
func TestSetZoneInterval(t *testing.T) {
	for _, sql := range []string{
		"SET TIME ZONE INTERVAL '-08:00' HOUR TO MINUTE",
		"SET TIME ZONE INTERVAL '1' HOUR",
		"SET TIME ZONE INTERVAL '1'",
		"SET TIME ZONE INTERVAL(3) '1'",
		"SET TRANSACTION SNAPSHOT '00000003-0000001B-1'",
		"SET TRANSACTION READ ONLY",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for sql, pos := range map[string]int{
		"SET TIME ZONE INTERVAL '1' YEAR": 27,
		"SET TIME ZONE INTERVAL":          22,
		"SET TIME ZONE INTERVAL(3)":       25,
		"SET TRANSACTION":                 15,
		"SET TRANSACTION SNAPSHOT":        24,
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Position != pos {
			t.Errorf("Parse(%q) = %v, want an error at %d", sql, err, pos)
		}
	}
}
