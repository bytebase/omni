package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseCallProcedureStmt(t *testing.T, sql string) *ast.CallProcedureStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CallProcedureStmt)
	if !ok {
		t.Fatalf("expected *ast.CallProcedureStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseDropProcedureStmt(t *testing.T, sql string) *ast.DropProcedureStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropProcedureStmt)
	if !ok {
		t.Fatalf("expected *ast.DropProcedureStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseCreateProcedureStmt(t *testing.T, sql string) *ast.CreateProcedureStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
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

// DROP (PROCEDURE | PROC) (DorisParser.g4:44): PROC is the rule's other
// spelling (compat.tsv's statement(PROCEDURE|PROC)/PROC row).
func TestDropProcedure_ProcSpelling(t *testing.T) {
	stmt := parseDropProcedureStmt(t, "DROP PROC my_proc")
	if stmt.Name.String() != "my_proc" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_proc")
	}
}

// ---------------------------------------------------------------------------
// ALTER PROCEDURE
// ---------------------------------------------------------------------------

func parseAlterProcedureStmt(t *testing.T, sql string) *ast.AlterProcedureStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterProcedureStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterProcedureStmt, got %T", file.Stmts[0])
	}
	return stmt
}

// ALTER PROCEDURE reads everything after the name as free text
// (DorisParser.g4:43), which is what the engine takes: it parses ALTER
// PROCEDURE a() and ALTER PROCEDURE a(a zzz) (compat.tsv's statement(.)* rows)
// and refuses a name without parentheses.
func TestAlterProcedure(t *testing.T) {
	tests := []struct {
		sql     string
		name    string
		rawText string
	}{
		{"ALTER PROCEDURE a()", "a", ""},
		{"ALTER PROCEDURE a(zzz)", "a", "zzz"},
		{"ALTER PROCEDURE a(a zzz)", "a", "a zzz"},
		{"ALTER PROCEDURE a() zzz", "a", "zzz"},
		{"ALTER PROC db.a()", "db.a", ""},
	}
	for _, tc := range tests {
		stmt := parseAlterProcedureStmt(t, tc.sql)
		if stmt.Name.String() != tc.name {
			t.Errorf("Parse(%q): Name = %q, want %q", tc.sql, stmt.Name.String(), tc.name)
		}
		if stmt.RawText != tc.rawText {
			t.Errorf("Parse(%q): RawText = %q, want %q", tc.sql, stmt.RawText, tc.rawText)
		}
	}
	for _, sql := range []string{"ALTER PROCEDURE a", "ALTER PROCEDURE a(", "ALTER PROCEDURE ()"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
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

// The grammar reads the parameter list and the text after it freely
// (DorisParser.g4:43 LEFT_PAREN .*? RIGHT_PAREN .*?), so Doris takes a
// procedure with no body (compat.tsv's statement(OR_REPLACE)?/%empty row) and
// one whose parentheses hold text that is no parameter list. Both keep that
// text in RawText.
func TestCreateProcedure_FreeText(t *testing.T) {
	tests := []struct {
		sql     string
		rawText string
	}{
		{"CREATE PROCEDURE a()", ""},
		{"CREATE PROCEDURE a(zzz)", "zzz"},
		{"CREATE PROCEDURE a(a zzz)", "a zzz"},
		{"CREATE PROCEDURE a() zzz", "zzz"},
	}
	for _, tc := range tests {
		stmt := parseCreateProcedureStmt(t, tc.sql)
		if stmt.RawText != tc.rawText {
			t.Errorf("Parse(%q): RawText = %q, want %q", tc.sql, stmt.RawText, tc.rawText)
		}
		if stmt.Body != "" || len(stmt.Parameters) != 0 {
			t.Errorf("Parse(%q): Body = %q Parameters = %d, want neither", tc.sql, stmt.Body, len(stmt.Parameters))
		}
	}
}

// REPLACE PROCEDURE is the grammar's third verb (DorisParser.g4:43), and
// compat.tsv's statement(ALTER|CREATE_(OR_REPLACE)?|REPLACE)/REPLACE row says
// the engine takes it.
func TestCreateProcedure_Replace(t *testing.T) {
	stmt := parseCreateProcedureStmt(t, "REPLACE PROCEDURE a()")
	if !stmt.Replace {
		t.Error("Replace should be true")
	}
	if stmt.OrReplace {
		t.Error("OrReplace should be false")
	}
	if stmt.Name.String() != "a" {
		t.Errorf("Name = %q, want a", stmt.Name.String())
	}
}

// The name and its parentheses are not free: the engine refuses
// CREATE PROCEDURE a (compat.tsv's SHOW_CREATE_PROCEDURE...!drop-1 row).
func TestCreateProcedure_NeedsParentheses(t *testing.T) {
	for _, sql := range []string{"CREATE PROCEDURE a", "CREATE PROCEDURE a("} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
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
	ap := &ast.AlterProcedureStmt{}
	if ap.Tag() != ast.T_AlterProcedureStmt {
		t.Errorf("AlterProcedureStmt.Tag() = %v, want T_AlterProcedureStmt", ap.Tag())
	}
	pp := &ast.ProcedureParam{}
	if pp.Tag() != ast.T_ProcedureParam {
		t.Errorf("ProcedureParam.Tag() = %v, want T_ProcedureParam", pp.Tag())
	}
}
