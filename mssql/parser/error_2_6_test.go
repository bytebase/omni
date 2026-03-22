package parser

import "testing"

// TestErrorSection2_6 tests soft-fail error detection for truncated subquery,
// CTE, and set operation (UNION/EXCEPT/INTERSECT) scenarios.
func TestErrorSection2_6(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"subquery FROM no content", "SELECT * FROM ("},
		{"CTE AS no query", "WITH cte AS ("},
		{"UNION no right query", "WITH cte AS (SELECT 1) SELECT 1 UNION"},
		{"EXCEPT no right query", "WITH cte AS (SELECT 1) SELECT 1 EXCEPT"},
		{"INTERSECT no right query", "WITH cte AS (SELECT 1) SELECT 1 INTERSECT"},
		{"UNION ALL no right query", "SELECT 1 UNION ALL"},
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
