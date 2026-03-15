package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateTypeObject(t *testing.T) {
	p := newTestParser("TYPE address_t AS OBJECT (street VARCHAR2(100), city VARCHAR2(50), zip NUMBER(5))")
	stmt := p.parseCreateTypeStmt(0, false, false, false, false)
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
	stmt := p.parseCreateTypeStmt(0, true, false, false, false)
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

func TestParseCreateTypeTableOf(t *testing.T) {
	p := newTestParser("TYPE num_table AS TABLE OF NUMBER")
	stmt := p.parseCreateTypeStmt(0, false, false, false, false)
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
	stmt := p.parseCreateTypeStmt(0, false, false, false, false)
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
	stmt := p.parseCreateTypeStmt(0, false, false, false, false)
	if !stmt.IsBody {
		t.Error("expected IsBody to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_TYPE" {
		t.Errorf("expected type name MY_TYPE, got %v", stmt.Name)
	}
}

func TestParseCreateTypeSchemaQualified(t *testing.T) {
	p := newTestParser("TYPE hr.address_t AS OBJECT (id NUMBER)")
	stmt := p.parseCreateTypeStmt(0, false, false, false, false)
	if stmt.Name == nil || stmt.Name.Schema != "HR" || stmt.Name.Name != "ADDRESS_T" {
		t.Errorf("expected HR.ADDRESS_T, got %v", stmt.Name)
	}
}

func TestParseCreateTypeLoc(t *testing.T) {
	p := newTestParser("TYPE t AS OBJECT (id NUMBER)")
	stmt := p.parseCreateTypeStmt(0, false, false, false, false)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}
