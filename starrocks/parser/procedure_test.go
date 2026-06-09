package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseCallProcedureStmt(t *testing.T, sql string) *ast.CallProcedureStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CallProcedureStmt)
	if !ok {
		t.Fatalf("expected *ast.CallProcedureStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseDropProcedureStmt(t *testing.T, sql string) *ast.DropProcedureStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropProcedureStmt)
	if !ok {
		t.Fatalf("expected *ast.DropProcedureStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseCreateProcedureStmt(t *testing.T, sql string) *ast.CreateProcedureStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateProcedureStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateProcedureStmt, got %T", file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// CALL
// ---------------------------------------------------------------------------

func TestCallProcedure_NoArgs(t *testing.T) {
	stmt := parseCallProcedureStmt(t, "CALL my_proc()")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if len(stmt.Args) != 0 {
		t.Errorf("Args: got %d, want 0", len(stmt.Args))
	}
}

func TestCallProcedure_WithArgs(t *testing.T) {
	stmt := parseCallProcedureStmt(t, "CALL my_proc(1, 'a')")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if len(stmt.Args) != 2 {
		t.Fatalf("Args: got %d, want 2", len(stmt.Args))
	}
}

func TestCallProcedure_QualifiedName(t *testing.T) {
	stmt := parseCallProcedureStmt(t, "CALL db.my_proc(1)")
	if stmt.Name.String() != "db.my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "db.my_proc")
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("Args: got %d, want 1", len(stmt.Args))
	}
}

// ---------------------------------------------------------------------------
// DROP PROCEDURE
// ---------------------------------------------------------------------------

func TestDropProcedure_Basic(t *testing.T) {
	stmt := parseDropProcedureStmt(t, "DROP PROCEDURE my_proc")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropProcedure_IfExists(t *testing.T) {
	stmt := parseDropProcedureStmt(t, "DROP PROCEDURE IF EXISTS my_proc")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

// ---------------------------------------------------------------------------
// CREATE PROCEDURE
// ---------------------------------------------------------------------------

func TestCreateProcedure_Basic(t *testing.T) {
	stmt := parseCreateProcedureStmt(t, "CREATE PROCEDURE my_proc() BEGIN SELECT 1; END")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if stmt.OrReplace {
		t.Error("OrReplace should be false")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Parameters) != 0 {
		t.Errorf("Parameters: got %d, want 0", len(stmt.Parameters))
	}
	if stmt.Comment != "" {
		t.Errorf("Comment = %q, want empty", stmt.Comment)
	}
	if !strings.Contains(stmt.Body, "SELECT 1") {
		t.Errorf("Body = %q, want it to contain SELECT 1", stmt.Body)
	}
}

func TestCreateProcedure_WithParams(t *testing.T) {
	stmt := parseCreateProcedureStmt(t, "CREATE PROCEDURE my_proc(IN x INT, OUT y VARCHAR(50)) BEGIN SELECT 1; END")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if len(stmt.Parameters) != 2 {
		t.Fatalf("Parameters: got %d, want 2", len(stmt.Parameters))
	}
	p0 := stmt.Parameters[0]
	if p0.Direction != "IN" {
		t.Errorf("Parameters[0].Direction = %q, want %q", p0.Direction, "IN")
	}
	if p0.Name != "x" {
		t.Errorf("Parameters[0].Name = %q, want %q", p0.Name, "x")
	}
	if p0.Type == nil {
		t.Fatal("Parameters[0].Type is nil")
	}
	p1 := stmt.Parameters[1]
	if p1.Direction != "OUT" {
		t.Errorf("Parameters[1].Direction = %q, want %q", p1.Direction, "OUT")
	}
	if p1.Name != "y" {
		t.Errorf("Parameters[1].Name = %q, want %q", p1.Name, "y")
	}
}

func TestCreateProcedure_OrReplaceWithComment(t *testing.T) {
	stmt := parseCreateProcedureStmt(t, "CREATE OR REPLACE PROCEDURE my_proc(x INT) COMMENT 'my proc' BEGIN SELECT 1; END")
	if !stmt.OrReplace {
		t.Error("OrReplace should be true")
	}
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if len(stmt.Parameters) != 1 {
		t.Fatalf("Parameters: got %d, want 1", len(stmt.Parameters))
	}
	if stmt.Parameters[0].Name != "x" {
		t.Errorf("Parameters[0].Name = %q, want %q", stmt.Parameters[0].Name, "x")
	}
	if stmt.Comment != "my proc" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "my proc")
	}
	if !strings.Contains(stmt.Body, "SELECT 1") {
		t.Errorf("Body = %q, want it to contain SELECT 1", stmt.Body)
	}
}

func TestCreateProcedure_NestedBeginEnd(t *testing.T) {
	// Nested BEGIN...END blocks (without END IF which would confuse the splitter).
	sql := `CREATE PROCEDURE my_proc() BEGIN
  BEGIN
    SELECT 1;
  END;
END`
	stmt := parseCreateProcedureStmt(t, sql)
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
	if !strings.Contains(stmt.Body, "SELECT 1") {
		t.Errorf("Body = %q, want it to contain SELECT 1", stmt.Body)
	}
}

func TestCreateProcedure_IfNotExists(t *testing.T) {
	stmt := parseCreateProcedureStmt(t, "CREATE PROCEDURE IF NOT EXISTS my_proc() BEGIN SELECT 1; END")
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
}

// ---------------------------------------------------------------------------
// Node tags
// ---------------------------------------------------------------------------

func TestProcedureNodeTags(t *testing.T) {
	cp := &ast.CreateProcedureStmt{}
	if cp.Tag() != ast.T_CreateProcedureStmt {
		t.Errorf("CreateProcedureStmt.Tag() = %v, want T_CreateProcedureStmt", cp.Tag())
	}
	call := &ast.CallProcedureStmt{}
	if call.Tag() != ast.T_CallProcedureStmt {
		t.Errorf("CallProcedureStmt.Tag() = %v, want T_CallProcedureStmt", call.Tag())
	}
	dp := &ast.DropProcedureStmt{}
	if dp.Tag() != ast.T_DropProcedureStmt {
		t.Errorf("DropProcedureStmt.Tag() = %v, want T_DropProcedureStmt", dp.Tag())
	}
	pp := &ast.ProcedureParam{}
	if pp.Tag() != ast.T_ProcedureParam {
		t.Errorf("ProcedureParam.Tag() = %v, want T_ProcedureParam", pp.Tag())
	}
}
