package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseAlterSessionSet tests ALTER SESSION SET with a single parameter.
func TestParseAlterSessionSet(t *testing.T) {
	result := ParseAndCheck(t, "ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY-MM-DD'")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterSessionStmt)
	if !ok {
		t.Fatalf("expected AlterSessionStmt, got %T", raw.Stmt)
	}
	if stmt.SetParams == nil || stmt.SetParams.Len() != 1 {
		t.Fatalf("expected 1 SetParam, got %v", stmt.SetParams)
	}
	sp := stmt.SetParams.Items[0].(*ast.SetParam)
	if sp.Name != "NLS_DATE_FORMAT" {
		t.Errorf("expected param name NLS_DATE_FORMAT, got %q", sp.Name)
	}
	if sp.Value == nil {
		t.Fatal("expected non-nil param value")
	}
	strLit, ok := sp.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral, got %T", sp.Value)
	}
	if strLit.Val != "YYYY-MM-DD" {
		t.Errorf("expected value 'YYYY-MM-DD', got %q", strLit.Val)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
}

// TestParseAlterSessionMultipleParams tests ALTER SESSION SET with multiple params.
func TestParseAlterSessionMultipleParams(t *testing.T) {
	result := ParseAndCheck(t, "ALTER SESSION SET NLS_LANGUAGE = 'AMERICAN' NLS_TERRITORY = 'AMERICA'")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterSessionStmt)
	if !ok {
		t.Fatalf("expected AlterSessionStmt, got %T", raw.Stmt)
	}
	if stmt.SetParams == nil || stmt.SetParams.Len() != 2 {
		t.Fatalf("expected 2 SetParams, got %d", stmt.SetParams.Len())
	}
	sp0 := stmt.SetParams.Items[0].(*ast.SetParam)
	if sp0.Name != "NLS_LANGUAGE" {
		t.Errorf("expected NLS_LANGUAGE, got %q", sp0.Name)
	}
	sp1 := stmt.SetParams.Items[1].(*ast.SetParam)
	if sp1.Name != "NLS_TERRITORY" {
		t.Errorf("expected NLS_TERRITORY, got %q", sp1.Name)
	}
}

// TestParseAlterSystemSet tests ALTER SYSTEM SET with a parameter.
func TestParseAlterSystemSet(t *testing.T) {
	result := ParseAndCheck(t, "ALTER SYSTEM SET open_cursors = 300")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterSystemStmt)
	if !ok {
		t.Fatalf("expected AlterSystemStmt, got %T", raw.Stmt)
	}
	if stmt.SetParams == nil || stmt.SetParams.Len() != 1 {
		t.Fatalf("expected 1 SetParam, got %v", stmt.SetParams)
	}
	sp := stmt.SetParams.Items[0].(*ast.SetParam)
	if sp.Name != "OPEN_CURSORS" {
		t.Errorf("expected param name OPEN_CURSORS, got %q", sp.Name)
	}
	if sp.Value == nil {
		t.Fatal("expected non-nil param value")
	}
}

// TestParseAlterSystemKillSession tests ALTER SYSTEM KILL SESSION.
func TestParseAlterSystemKillSession(t *testing.T) {
	result := ParseAndCheck(t, "ALTER SYSTEM KILL SESSION '12,34'")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterSystemStmt)
	if !ok {
		t.Fatalf("expected AlterSystemStmt, got %T", raw.Stmt)
	}
	if stmt.Kill != "12,34" {
		t.Errorf("expected Kill='12,34', got %q", stmt.Kill)
	}
}

// TestParseAlterIndex tests that ALTER INDEX parses without error.
func TestParseAlterIndex(t *testing.T) {
	result := ParseAndCheck(t, "ALTER INDEX my_idx REBUILD")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	_, ok := raw.Stmt.(*ast.AlterSessionStmt)
	if !ok {
		t.Fatalf("expected AlterSessionStmt (placeholder), got %T", raw.Stmt)
	}
}

// TestParseAlterSequence tests that ALTER SEQUENCE parses without error.
func TestParseAlterSequence(t *testing.T) {
	result := ParseAndCheck(t, "ALTER SEQUENCE my_seq INCREMENT BY 10")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
}

// TestParseAlterView tests that ALTER VIEW parses without error.
func TestParseAlterView(t *testing.T) {
	result := ParseAndCheck(t, "ALTER VIEW my_view COMPILE")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
}

// TestParseAlterLocSet tests that Loc is set on AlterSessionStmt.
func TestParseAlterLocSet(t *testing.T) {
	result := ParseAndCheck(t, "ALTER SESSION SET X = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.AlterSessionStmt)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got End=%d", stmt.Loc.End)
	}
}
