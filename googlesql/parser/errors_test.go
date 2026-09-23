package parser

import "testing"

func TestParseError_Error(t *testing.T) {
	e := &ParseError{Position: 3, End: 9, Message: "syntax error at or near FOO"}
	if got := e.Error(); got != "syntax error at or near FOO" {
		t.Errorf("Error() = %q, want %q", got, "syntax error at or near FOO")
	}
}

func TestParseError_ImplementsError(t *testing.T) {
	var err error = &ParseError{Message: "boom"}
	if err.Error() != "boom" {
		t.Errorf("Error() = %q, want %q", err.Error(), "boom")
	}
}

func TestParseError_PositionPreserved(t *testing.T) {
	e := &ParseError{Position: 10, End: 14, Message: "x"}
	if e.Position != 10 || e.End != 14 {
		t.Errorf("Position/End = %d/%d, want 10/14", e.Position, e.End)
	}
}
