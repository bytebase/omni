package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func TestParseWithDelimiter(t *testing.T) {
	sql := "DELIMITER ;;\nCREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\nEND;;\nDELIMITER ;\nSELECT 2;"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("got %d statements, want 2", len(list.Items))
	}
	if _, ok := list.Items[0].(*nodes.CreateFunctionStmt); !ok {
		t.Errorf("stmt[0] type = %T, want *CreateFunctionStmt", list.Items[0])
	}
	if _, ok := list.Items[1].(*nodes.SelectStmt); !ok {
		t.Errorf("stmt[1] type = %T, want *SelectStmt", list.Items[1])
	}
}

func TestParseWithDelimiterLocAbsolute(t *testing.T) {
	sql := "SELECT 1; SELECT 2;"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("got %d statements, want 2", len(list.Items))
	}
	// Second SELECT starts at byte offset 10 in original SQL.
	sel2, ok := list.Items[1].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("stmt[1] type = %T, want *SelectStmt", list.Items[1])
	}
	if sel2.Loc.Start != 10 {
		t.Errorf("stmt[1].Loc.Start = %d, want 10", sel2.Loc.Start)
	}
}

func TestParseErrorPositionAbsolute(t *testing.T) {
	sql := "SELECT 1; INVALID!;"
	_, err := Parse(sql)
	if err == nil {
		t.Fatal("expected error")
	}
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("error type = %T, want *ParseError", err)
	}
	// "INVALID" starts at byte 10. Error position should be >= 10.
	if pe.Position < 10 {
		t.Errorf("ParseError.Position = %d, want >= 10", pe.Position)
	}
	if pe.Line != 1 {
		t.Errorf("ParseError.Line = %d, want 1", pe.Line)
	}
}
