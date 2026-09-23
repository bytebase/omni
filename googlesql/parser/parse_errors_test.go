package parser

import (
	"errors"
	"testing"
)

// TestParse_ReturnsParseErrors pins the error contract every omni engine
// shares: Parse returns nil on success, otherwise an error that errors.As
// resolves to the first *ParseError with a byte Position, while AllErrors
// still yields every error a multi-statement script produced.
func TestParse_ReturnsParseErrors(t *testing.T) {
	if _, err := Parse("SELECT 1"); err != nil {
		t.Fatalf("Parse(valid) = %v, want nil", err)
	}

	file, err := Parse("SELECT 1 ))); SELECT 2; SELECT 3 )))")
	if err == nil {
		t.Fatal("Parse(two bad segments) returned nil error")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("errors.As(*ParseError) failed on %T", err)
	}
	if pe.Position != len("SELECT 1 ") {
		t.Errorf("first error Position = %d, want %d (the first ')')", pe.Position, len("SELECT 1 "))
	}
	if pe.Message == "" {
		t.Error("first error has an empty Message")
	}
	if got := AllErrors(err); len(got) < 2 {
		t.Errorf("AllErrors = %d errors, want at least 2 (one per bad segment): %v", len(got), got)
	}
	if file == nil || len(file.Stmts) != 1 {
		t.Errorf("File keeps the statements that parsed: got %v", file)
	}
	if AllErrors(nil) != nil || AllErrors(errors.New("other")) != nil {
		t.Error("AllErrors must be nil for nil and for non-parse errors")
	}
}
