package parser

import "testing"

func TestErrorSection2_1(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"SELECT_1_OR_truncated", "SELECT 1 OR"},
		{"SELECT_1_AND_truncated", "SELECT 1 AND"},
		{"SELECT_1_LT_truncated", "SELECT 1 <"},
		{"SELECT_1_GT_truncated", "SELECT 1 >"},
		{"SELECT_1_EQ_truncated", "SELECT 1 ="},
		{"SELECT_1_LE_truncated", "SELECT 1 <="},
		{"SELECT_1_GE_truncated", "SELECT 1 >="},
		{"SELECT_1_NE_truncated", "SELECT 1 <>"},
		{"SELECT_1_PLUS_truncated", "SELECT 1 +"},
		{"SELECT_1_MINUS_truncated", "SELECT 1 -"},
		{"SELECT_1_MUL_truncated", "SELECT 1 *"},
		{"SELECT_1_DIV_truncated", "SELECT 1 /"},
		{"SELECT_1_MOD_truncated", "SELECT 1 %"},
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
