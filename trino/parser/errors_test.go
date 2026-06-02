package parser

import (
	"testing"

	"github.com/bytebase/omni/trino/ast"
)

// TestParseError_Error verifies ParseError implements error by returning its
// message verbatim, leaving line/column conversion to a caller-side LineTable.
func TestParseError_Error(t *testing.T) {
	pe := &ParseError{Loc: ast.Loc{Start: 7, End: 13}, Msg: "syntax error at or near FROM"}
	if got := pe.Error(); got != "syntax error at or near FROM" {
		t.Fatalf("ParseError.Error() = %q, want %q", got, "syntax error at or near FROM")
	}
}

// TestParseError_IsErrorInterface verifies *ParseError satisfies the error
// interface so Parse can return it polymorphically.
func TestParseError_IsErrorInterface(t *testing.T) {
	var err error = &ParseError{Msg: "boom"}
	if err.Error() != "boom" {
		t.Fatalf("error.Error() = %q, want %q", err.Error(), "boom")
	}
}
