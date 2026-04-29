package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseIdentifier tests parsing of simple unquoted identifiers.
func TestParseIdentifier(t *testing.T) {
	p := newTestParser("my_table")
	name, parseErr1 := p.parseIdentifier()
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if name != "MY_TABLE" {
		t.Errorf("expected MY_TABLE, got %q", name)
	}
}

// TestParseQuotedIdentifier tests parsing of "double-quoted" identifiers.
func TestParseQuotedIdentifier(t *testing.T) {
	p := newTestParser(`"MyColumn"`)
	name, parseErr2 := p.parseIdentifier()
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if name != "MyColumn" {
		t.Errorf("expected MyColumn, got %q", name)
	}
}

// TestParseObjectNameSimple tests parsing a simple (unqualified) object name.
func TestParseObjectNameSimple(t *testing.T) {
	p := newTestParser("employees")
	obj, parseErr3 := p.parseObjectName()
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if obj.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %q", obj.Name)
	}
	if obj.Schema != "" {
		t.Errorf("expected empty schema, got %q", obj.Schema)
	}
	if obj.DBLink != "" {
		t.Errorf("expected empty dblink, got %q", obj.DBLink)
	}
}

// TestParseObjectNameSchemaQualified tests parsing schema.object.
func TestParseObjectNameSchemaQualified(t *testing.T) {
	p := newTestParser("hr.employees")
	obj, parseErr4 := p.parseObjectName()
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if obj.Schema != "HR" {
		t.Errorf("expected schema HR, got %q", obj.Schema)
	}
	if obj.Name != "EMPLOYEES" {
		t.Errorf("expected name EMPLOYEES, got %q", obj.Name)
	}
}

// TestParseObjectNameWithDBLink tests parsing schema.object@dblink.
func TestParseObjectNameWithDBLink(t *testing.T) {
	p := newTestParser("hr.employees@remote_db")
	obj, parseErr5 := p.parseObjectName()
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if obj.Schema != "HR" {
		t.Errorf("expected schema HR, got %q", obj.Schema)
	}
	if obj.Name != "EMPLOYEES" {
		t.Errorf("expected name EMPLOYEES, got %q", obj.Name)
	}
	if obj.DBLink != "REMOTE_DB" {
		t.Errorf("expected dblink REMOTE_DB, got %q", obj.DBLink)
	}
}

// TestParseObjectNameDBLinkOnly tests parsing object@dblink without schema.
func TestParseObjectNameDBLinkOnly(t *testing.T) {
	p := newTestParser("employees@remote_db")
	obj, parseErr6 := p.parseObjectName()
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if obj.Schema != "" {
		t.Errorf("expected empty schema, got %q", obj.Schema)
	}
	if obj.Name != "EMPLOYEES" {
		t.Errorf("expected name EMPLOYEES, got %q", obj.Name)
	}
	if obj.DBLink != "REMOTE_DB" {
		t.Errorf("expected dblink REMOTE_DB, got %q", obj.DBLink)
	}
}

// TestParseObjectNameQuoted tests parsing "Schema"."Object".
func TestParseObjectNameQuoted(t *testing.T) {
	p := newTestParser(`"MySchema"."MyTable"`)
	obj, parseErr7 := p.parseObjectName()
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if obj.Schema != "MySchema" {
		t.Errorf("expected schema MySchema, got %q", obj.Schema)
	}
	if obj.Name != "MyTable" {
		t.Errorf("expected name MyTable, got %q", obj.Name)
	}
}

// TestParseColumnRefSimple tests parsing a simple column reference.
func TestParseColumnRefSimple(t *testing.T) {
	p := newTestParser("salary")
	col, parseErr8 := p.parseColumnRef()
	if parseErr8 != nil {
		t.Fatalf("parse: %v", parseErr8)
	}
	if col.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %q", col.Column)
	}
	if col.Table != "" {
		t.Errorf("expected empty table, got %q", col.Table)
	}
}

// TestParseColumnRefQualified tests parsing table.column.
func TestParseColumnRefQualified(t *testing.T) {
	p := newTestParser("e.salary")
	col, parseErr9 := p.parseColumnRef()
	if parseErr9 != nil {
		t.Fatalf("parse: %v", parseErr9)
	}
	if col.Table != "E" {
		t.Errorf("expected table E, got %q", col.Table)
	}
	if col.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %q", col.Column)
	}
}

// TestParseColumnRefSchemaQualified tests parsing schema.table.column.
func TestParseColumnRefSchemaQualified(t *testing.T) {
	p := newTestParser("hr.employees.salary")
	col, parseErr10 := p.parseColumnRef()
	if parseErr10 != nil {
		t.Fatalf("parse: %v", parseErr10)
	}
	if col.Schema != "HR" {
		t.Errorf("expected schema HR, got %q", col.Schema)
	}
	if col.Table != "EMPLOYEES" {
		t.Errorf("expected table EMPLOYEES, got %q", col.Table)
	}
	if col.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %q", col.Column)
	}
}

// TestParseColumnRefStar tests parsing table.*.
func TestParseColumnRefStar(t *testing.T) {
	p := newTestParser("e.*")
	col, parseErr11 := p.parseColumnRef()
	if parseErr11 != nil {
		t.Fatalf("parse: %v", parseErr11)
	}
	if col.Table != "E" {
		t.Errorf("expected table E, got %q", col.Table)
	}
	if col.Column != "*" {
		t.Errorf("expected column *, got %q", col.Column)
	}
}

// TestParseBindVariable tests parsing bind variables.
func TestParseBindVariable(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{":param1", "param1"},
		{":1", "1"},
		{":emp_id", "emp_id"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			bv, parseErr1 := p.parseBindVariable()
			if parseErr1 != nil {
				t.Fatalf("parse: %v", parseErr1)
			}
			if bv == nil {
				t.Fatal("expected non-nil BindVariable")
			}
			if bv.Name != tc.name {
				t.Errorf("expected name %q, got %q", tc.name, bv.Name)
			}
		})
	}
}

// TestParsePseudoColumn tests parsing pseudo-columns.
func TestParsePseudoColumn(t *testing.T) {
	tests := []struct {
		input string
		ptype ast.PseudoColumnType
	}{
		{"ROWID", ast.PSEUDO_ROWID},
		{"ROWNUM", ast.PSEUDO_ROWNUM},
		{"LEVEL", ast.PSEUDO_LEVEL},
		{"SYSDATE", ast.PSEUDO_SYSDATE},
		{"SYSTIMESTAMP", ast.PSEUDO_SYSTIMESTAMP},
		{"USER", ast.PSEUDO_USER},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			pc, parseErr2 := p.parsePseudoColumn()
			if parseErr2 != nil {
				t.Fatalf("parse: %v", parseErr2)
			}
			if pc == nil {
				t.Fatal("expected non-nil PseudoColumn")
			}
			if pc.Type != tc.ptype {
				t.Errorf("expected type %d, got %d", tc.ptype, pc.Type)
			}
		})
	}
}

// TestParseObjectNameLoc tests that locations are recorded.
func TestParseObjectNameLoc(t *testing.T) {
	p := newTestParser("hr.employees")
	obj, parseErr12 := p.parseObjectName()
	if parseErr12 != nil {
		t.Fatalf("parse: %v", parseErr12)
	}
	if obj.Loc.Start != 0 {
		t.Errorf("expected Start=0, got %d", obj.Loc.Start)
	}
	if obj.Loc.End <= obj.Loc.Start {
		t.Errorf("expected End > Start, got End=%d Start=%d", obj.Loc.End, obj.Loc.Start)
	}
}

// newTestParser creates a parser for testing helper functions.
func newTestParser(sql string) *Parser {
	p := &Parser{
		lexer: NewLexer(sql),
	}
	p.advance()
	return p
}
