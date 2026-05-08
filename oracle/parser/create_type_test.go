package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateTypeObject(t *testing.T) {
	p := newTestParser("TYPE address_t AS OBJECT (street VARCHAR2(100), city VARCHAR2(50), zip NUMBER(5))")
	stmt, parseErr1 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateTypeStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "ADDRESS_T" {
		t.Errorf("expected type name ADDRESS_T, got %v", stmt.Name)
	}
	if stmt.Attributes == nil || stmt.Attributes.Len() != 3 {
		t.Fatalf("expected 3 attributes, got %d", stmt.Attributes.Len())
	}
	attr0 := stmt.Attributes.Items[0].(*ast.ColumnDef)
	if attr0.Name != "STREET" {
		t.Errorf("expected attribute name STREET, got %q", attr0.Name)
	}
	if attr0.TypeName == nil || attr0.TypeName.Names.Len() == 0 {
		t.Error("expected non-nil TypeName for first attribute")
	}
}

func TestParseCreateTypeOrReplace(t *testing.T) {
	p := newTestParser("TYPE my_type AS OBJECT (id NUMBER)")
	stmt, parseErr2 := p.parseCreateTypeStmt(0, true, false, false, false)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

func TestParseCreateTypeTableOf(t *testing.T) {
	p := newTestParser("TYPE num_table AS TABLE OF NUMBER")
	stmt, parseErr3 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if stmt.AsTable == nil {
		t.Fatal("expected non-nil AsTable")
	}
	if stmt.AsTable.Names == nil || stmt.AsTable.Names.Len() == 0 {
		t.Error("expected AsTable type names")
	}
	nameStr := stmt.AsTable.Names.Items[0].(*ast.String)
	if nameStr.Str != "NUMBER" {
		t.Errorf("expected NUMBER, got %q", nameStr.Str)
	}
}

func TestParseCreateTypeVarray(t *testing.T) {
	p := newTestParser("TYPE phone_list AS VARRAY (10) OF VARCHAR2(20)")
	stmt, parseErr4 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt.AsVarray == nil {
		t.Fatal("expected non-nil AsVarray")
	}
	if stmt.VarraySize == nil {
		t.Fatal("expected non-nil VarraySize")
	}
	sizeNum := stmt.VarraySize.(*ast.NumberLiteral)
	if sizeNum.Ival != 10 {
		t.Errorf("expected varray size 10, got %d", sizeNum.Ival)
	}
}

func TestParseCreateTypeBody(t *testing.T) {
	p := newTestParser("TYPE BODY my_type IS BEGIN NULL; END my_type")
	stmt, parseErr5 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if !stmt.IsBody {
		t.Error("expected IsBody to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_TYPE" {
		t.Errorf("expected type name MY_TYPE, got %v", stmt.Name)
	}
}

func TestParseCreateTypeSchemaQualified(t *testing.T) {
	p := newTestParser("TYPE hr.address_t AS OBJECT (id NUMBER)")
	stmt, parseErr6 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if stmt.Name == nil || stmt.Name.Schema != "HR" || stmt.Name.Name != "ADDRESS_T" {
		t.Errorf("expected HR.ADDRESS_T, got %v", stmt.Name)
	}
}

func TestParseCreateTypeLoc(t *testing.T) {
	p := newTestParser("TYPE t AS OBJECT (id NUMBER)")
	stmt, parseErr7 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

func TestP2CreateTypeBodyMemberStructure(t *testing.T) {
	result := ParseAndCheck(t, `CREATE TYPE BODY person_type AS
  MEMBER FUNCTION get_name RETURN VARCHAR2 IS
  BEGIN
    RETURN first_name;
  END get_name;
  STATIC PROCEDURE set_count(p_count NUMBER) IS
  BEGIN
    NULL;
  END set_count;
END;`)
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateTypeStmt)
	if !ok {
		t.Fatalf("expected CreateTypeStmt, got %T", raw.Stmt)
	}
	if !stmt.IsBody {
		t.Fatal("expected type body")
	}
	if stmt.Body == nil || stmt.Body.Len() != 2 {
		t.Fatalf("expected 2 body members, got %#v", stmt.Body)
	}

	memberFn := stmt.Body.Items[0].(*ast.TypeBodyMember)
	if memberFn.Kind != ast.TYPE_BODY_MEMBER {
		t.Fatalf("member 1 kind=%d, want MEMBER", memberFn.Kind)
	}
	fn, ok := memberFn.Subprog.(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("member 1 subprogram=%T, want function", memberFn.Subprog)
	}
	fnBlock, ok := fn.Body.(*ast.PLSQLBlock)
	if !ok {
		t.Fatalf("function body=%T, want PLSQLBlock", fn.Body)
	}
	if fn.Name == nil || fn.Name.Name != "GET_NAME" || fnBlock.Statements == nil || fnBlock.Statements.Len() != 1 {
		t.Fatalf("function body was not structured: %#v", fn)
	}
	if _, ok := fnBlock.Statements.Items[0].(*ast.PLSQLReturn); !ok {
		t.Fatalf("function statement=%T, want PLSQLReturn", fnBlock.Statements.Items[0])
	}

	memberProc := stmt.Body.Items[1].(*ast.TypeBodyMember)
	if memberProc.Kind != ast.TYPE_BODY_STATIC {
		t.Fatalf("member 2 kind=%d, want STATIC", memberProc.Kind)
	}
	proc, ok := memberProc.Subprog.(*ast.CreateProcedureStmt)
	if !ok {
		t.Fatalf("member 2 subprogram=%T, want procedure", memberProc.Subprog)
	}
	procBlock, ok := proc.Body.(*ast.PLSQLBlock)
	if !ok {
		t.Fatalf("procedure body=%T, want PLSQLBlock", proc.Body)
	}
	if proc.Name == nil || proc.Name.Name != "SET_COUNT" || procBlock.Statements == nil || procBlock.Statements.Len() != 1 {
		t.Fatalf("procedure body was not structured: %#v", proc)
	}
	if _, ok := procBlock.Statements.Items[0].(*ast.PLSQLNull); !ok {
		t.Fatalf("procedure statement=%T, want PLSQLNull", procBlock.Statements.Items[0])
	}
}

func TestP2CreateTypeBodyRequiresPLSQLBlock(t *testing.T) {
	ParseShouldFail(t, "CREATE TYPE BODY person_type AS MEMBER PROCEDURE p IS END p; END;")
}

func TestP2CreateTypeBodyLoc(t *testing.T) {
	stmt := ParseAndCheck(t, `CREATE TYPE BODY person_type AS
  MEMBER PROCEDURE p IS
  BEGIN
    NULL;
  END p;
END;`).Items[0].(*ast.RawStmt).Stmt.(*ast.CreateTypeStmt)
	if stmt.Loc.Start != 0 || stmt.Loc.End <= stmt.Loc.Start {
		t.Fatalf("invalid type body Loc=%+v", stmt.Loc)
	}
	member := stmt.Body.Items[0].(*ast.TypeBodyMember)
	if member.Loc.Start <= stmt.Loc.Start || member.Loc.End > stmt.Loc.End {
		t.Fatalf("member Loc=%+v not inside stmt Loc=%+v", member.Loc, stmt.Loc)
	}
}
