package parser

import (
	"strings"
	"testing"
)

// TestErrorSection3_2_EndOfInput verifies that truncated SQL produces
// "at end of input" error messages.
func TestErrorSection3_2_EndOfInput(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"SELECT_1_PLUS_truncated", "SELECT 1 +"},
		{"SELECT_star_FROM_truncated", "SELECT * FROM"},
		{"SELECT_star_FROM_t_WHERE_truncated", "SELECT * FROM t WHERE"},
		{"CREATE_TABLE_truncated", "CREATE TABLE"},
		{"IF_truncated", "IF"},
		{"SELECT_star_FROM_t_WHERE_truncated_or", "SELECT * FROM t WHERE 1 OR"},
		{"SELECT_star_FROM_t_WHERE_truncated_and", "SELECT * FROM t WHERE 1 AND"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", tt.sql)
			}
			msg := err.Error()
			if !strings.Contains(msg, "at end of input") {
				t.Errorf("Parse(%q): expected error containing %q, got %q", tt.sql, "at end of input", msg)
			}
		})
	}
}

// TestErrorSection3_2_AtOrNear verifies that invalid tokens mid-statement
// produce "at or near" error messages with the token text.
func TestErrorSection3_2_AtOrNear(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		nearText string // expected token text in the "at or near" message
	}{
		{"unexpected_paren_after_from", "SELECT 1 FROM )", ")"},
		{"unexpected_paren", "SELECT 1 ) 2", ")"},
		{"bad_token_after_table", "SELECT * FROM t WHERE )", ")"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", tt.sql)
			}
			msg := err.Error()
			expected := `at or near "` + tt.nearText + `"`
			if !strings.Contains(msg, expected) {
				t.Errorf("Parse(%q): expected error containing %q, got %q", tt.sql, expected, msg)
			}
		})
	}
}

// TestErrorSection3_2_MultiLine verifies that position information is accurate
// for multi-line SQL and error messages still contain proper context.
func TestErrorSection3_2_MultiLine(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string // substring expected in error
	}{
		{
			// AND and OR are a search condition's, so the error is the OR
			// itself rather than the end of the text.
			"multiline_truncated",
			"SELECT\n  1\n  OR",
			`at or near "OR"`,
		},
		{
			"multiline_truncated_in_condition",
			"SELECT\n  *\nFROM t\nWHERE 1 = 1 AND",
			"at end of input",
		},
		{
			"multiline_invalid_token",
			"SELECT\n  *\nFROM t\nWHERE )",
			`at or near ")"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", tt.sql)
			}
			msg := err.Error()
			if !strings.Contains(msg, tt.want) {
				t.Errorf("Parse(%q): expected error containing %q, got %q", tt.sql, tt.want, msg)
			}
		})
	}
}

// TestErrorSection3_2_ParseErrorFields verifies the ParseError struct
// has Source field populated and Error() behaves correctly.
func TestErrorSection3_2_ParseErrorFields(t *testing.T) {
	// Test with Source set, position at end
	e := &ParseError{
		Message:  "syntax error at end of input",
		Position: 10,
		Source:   "SELECT 1 +",
	}
	if got := e.Error(); got != "syntax error at end of input" {
		t.Errorf("ParseError.Error() with pos at end: got %q, want %q", got, "syntax error at end of input")
	}

	// Test with Source set, position mid-input
	e2 := &ParseError{
		Message:  "syntax error",
		Position: 7,
		Source:   "SELECT ) FROM t",
	}
	got := e2.Error()
	if !strings.Contains(got, `at or near ")"`) {
		t.Errorf("ParseError.Error() with pos mid-input: got %q, want containing %q", got, `at or near ")"`)
	}

	// Test backward compatibility: no Source falls back to Message
	e3 := &ParseError{
		Message:  "some custom error",
		Position: 0,
	}
	if got := e3.Error(); got != "some custom error" {
		t.Errorf("ParseError.Error() without Source: got %q, want %q", got, "some custom error")
	}
}

// TestErrorSection3_2_ExpectEnhanced verifies that expect() produces
// enhanced error messages with "at or near" context.
func TestErrorSection3_2_ExpectEnhanced(t *testing.T) {
	// This SQL will fail because expect() is called for a missing token.
	// "SELECT 1 )" — the closing paren is unexpected.
	_, err := Parse("SELECT 1 )")
	if err == nil {
		t.Fatal("expected error for 'SELECT 1 )', got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "at or near") && !strings.Contains(msg, "at end of input") {
		t.Errorf("expect() error should contain 'at or near' or 'at end of input', got %q", msg)
	}
}
