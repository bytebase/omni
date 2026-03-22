package parser

import "testing"

// TestErrorSection2_5 tests soft-fail error detection for truncated SELECT
// clauses: FROM, JOIN, WHERE, GROUP BY, HAVING, ORDER BY, TOP.
func TestErrorSection2_5(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"FROM no table", "SELECT * FROM"},
		{"JOIN no right table", "SELECT * FROM t JOIN"},
		{"LEFT JOIN no right table", "SELECT * FROM t LEFT JOIN"},
		{"CROSS JOIN no right table", "SELECT * FROM t CROSS JOIN"},
		{"JOIN ON no condition", "SELECT * FROM t JOIN t2 ON"},
		{"WHERE no condition", "SELECT * FROM t WHERE"},
		{"GROUP BY no expressions", "SELECT * FROM t GROUP BY"},
		{"HAVING no condition", "SELECT * FROM t HAVING"},
		{"ORDER BY no sort expressions", "SELECT * FROM t ORDER BY"},
		{"TOP no count", "SELECT TOP"},
		{"TOP paren no expression", "SELECT TOP ("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Errorf("Parse(%q): expected error, got nil", tt.sql)
			}
		})
	}
}
