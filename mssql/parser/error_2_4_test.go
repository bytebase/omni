package parser

import "testing"

func TestErrorSection2_4(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"COLLATE_no_collation", "SELECT 'a' COLLATE"},
		{"AT_TIME_ZONE_no_zone", "SELECT GETDATE() AT TIME ZONE"},
		{"unary_minus_no_operand", "SELECT -"},
		{"unary_plus_no_operand", "SELECT +"},
		{"bitwise_not_no_operand", "SELECT ~"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.sql)
			}
		})
	}
}
