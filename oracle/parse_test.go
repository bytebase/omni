package oracle

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseSplitScript(t *testing.T) {
	sql := "SELECT 1 FROM dual;\nBEGIN NULL; END;\n/\nSELECT 2 FROM dual;"

	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(stmts) != 3 {
		t.Fatalf("got %d statements, want 3", len(stmts))
	}

	if stmts[0].Text != "SELECT 1 FROM dual" {
		t.Fatalf("stmt[0].Text = %q", stmts[0].Text)
	}
	if _, ok := stmts[0].AST.(*ast.SelectStmt); !ok {
		t.Fatalf("stmt[0].AST = %T, want SelectStmt", stmts[0].AST)
	}
	if stmts[0].ByteStart != 0 || stmts[0].ByteEnd != len("SELECT 1 FROM dual") {
		t.Fatalf("stmt[0] range = [%d,%d]", stmts[0].ByteStart, stmts[0].ByteEnd)
	}
	if stmts[0].Start != (Position{Line: 1, Column: 1}) {
		t.Fatalf("stmt[0].Start = %+v", stmts[0].Start)
	}

	if stmts[1].Text != "\nBEGIN NULL; END;" {
		t.Fatalf("stmt[1].Text = %q", stmts[1].Text)
	}
	if _, ok := stmts[1].AST.(*ast.PLSQLBlock); !ok {
		t.Fatalf("stmt[1].AST = %T, want PLSQLBlock", stmts[1].AST)
	}
	if ast.NodeLoc(stmts[1].AST).Start != stmts[1].ByteStart+1 {
		t.Fatalf("stmt[1] AST Loc.Start = %d, want %d", ast.NodeLoc(stmts[1].AST).Start, stmts[1].ByteStart+1)
	}

	if stmts[2].Text != "\nSELECT 2 FROM dual" {
		t.Fatalf("stmt[2].Text = %q", stmts[2].Text)
	}
	if stmts[2].Start != (Position{Line: 4, Column: 1}) {
		t.Fatalf("stmt[2].Start = %+v", stmts[2].Start)
	}
}

func TestParseReturnsErrorForSQLPlusCommands(t *testing.T) {
	sql := "SET DEFINE OFF\n" +
		"PROMPT setup\n" +
		"SELECT 1 FROM dual;\n" +
		"SPOOL out.log\n" +
		"SELECT 2 FROM dual;\n" +
		"EXIT SUCCESS\n"

	if _, err := Parse(sql); err == nil {
		t.Fatal("Parse() error = nil, want SQL*Plus command parse error")
	}
}
