package parser

import "testing"

func TestParseFromMultipleValuesSubqueriesWithColumnAliases(t *testing.T) {
	tests := []string{
		`SELECT * FROM (VALUES (1)) v(a), (VALUES (2)) w(b)`,
		`SELECT * FROM (VALUES ('week', '7 d')) intervals (str, interval), (VALUES (timestamp '2020-02-29')) ts (ts)`,
		`SELECT interval::interval FROM (VALUES ('7 d')) intervals (interval)`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}
}
