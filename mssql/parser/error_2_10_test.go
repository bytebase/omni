package parser

import (
	"testing"
)

func TestErrorSection2_10(t *testing.T) {
	// Error cases: truncated control flow / declaration statements
	errorTests := []struct {
		name string
		sql  string
	}{
		{"IF_no_condition", "IF"},
		{"WHILE_no_condition", "WHILE"},
		{"DECLARE_no_type", "DECLARE @v"},
		{"SET_var_eq_no_value", "SET @v ="},
		{"EXEC_no_proc", "EXEC"},
		{"PRINT_no_expr", "PRINT"},
	}
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.sql)
			}
		})
	}

	// Valid cases (should NOT error)
	validTests := []struct {
		name string
		sql  string
	}{
		{"RETURN_no_value", "RETURN"},
		{"BEGIN_TRANSACTION", "BEGIN TRANSACTION"},
		{"THROW_rethrow", "THROW"},
	}
	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.sql, err)
			}
		})
	}

	// IF 1=1 with no body — this should error because there's no statement body
	t.Run("IF_condition_no_body", func(t *testing.T) {
		_, err := Parse("IF 1=1")
		if err == nil {
			t.Errorf("Parse(%q) expected error, got nil", "IF 1=1")
		}
	})
}
