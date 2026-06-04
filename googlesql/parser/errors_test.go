package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

func TestParseError_Error(t *testing.T) {
	e := &ParseError{Loc: ast.Loc{Start: 3, End: 9}, Msg: "syntax error at or near FOO"}
	if got := e.Error(); got != "syntax error at or near FOO" {
		t.Errorf("Error() = %q, want %q", got, "syntax error at or near FOO")
	}
}

func TestParseError_ImplementsError(t *testing.T) {
	var err error = &ParseError{Msg: "boom"}
	if err.Error() != "boom" {
		t.Errorf("Error() = %q, want %q", err.Error(), "boom")
	}
}

func TestParseError_LocPreserved(t *testing.T) {
	e := &ParseError{Loc: ast.Loc{Start: 10, End: 14}, Msg: "x"}
	if e.Loc != (ast.Loc{Start: 10, End: 14}) {
		t.Errorf("Loc = %+v, want {10, 14}", e.Loc)
	}
}
