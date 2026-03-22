package parser

import "testing"

func TestErrorSection2_2(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"SELECT_1_BETWEEN_truncated", "SELECT 1 BETWEEN"},
		{"SELECT_1_BETWEEN_0_AND_truncated", "SELECT 1 BETWEEN 0 AND"},
		{"SELECT_1_NOT_BETWEEN_truncated", "SELECT 1 NOT BETWEEN"},
		{"SELECT_a_LIKE_truncated", "SELECT 'a' LIKE"},
		{"SELECT_a_NOT_LIKE_truncated", "SELECT 'a' NOT LIKE"},
		{"SELECT_1_IN_LPAREN_truncated", "SELECT 1 IN ("},
		{"SELECT_1_NOT_IN_LPAREN_truncated", "SELECT 1 NOT IN ("},
		{"SELECT_NOT_truncated", "SELECT NOT"},
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
