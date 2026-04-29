package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateIndex(t *testing.T) {
	p := newTestParser("INDEX idx_emp_name ON employees (last_name)")
	stmt, parseErr1 := p.parseCreateIndexStmt(0)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateIndexStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "IDX_EMP_NAME" {
		t.Errorf("expected index name IDX_EMP_NAME, got %v", stmt.Name)
	}
	if stmt.Table == nil || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table EMPLOYEES, got %v", stmt.Table)
	}
	if stmt.Columns == nil || stmt.Columns.Len() != 1 {
		t.Fatalf("expected 1 column, got %d", stmt.Columns.Len())
	}
	col0 := stmt.Columns.Items[0].(*ast.IndexColumn)
	cr := col0.Expr.(*ast.ColumnRef)
	if cr.Column != "LAST_NAME" {
		t.Errorf("expected column LAST_NAME, got %q", cr.Column)
	}
}

func TestParseCreateUniqueIndex(t *testing.T) {
	p := newTestParser("UNIQUE INDEX idx_emp_id ON hr.employees (employee_id)")
	stmt, parseErr2 := p.parseCreateIndexStmt(0)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if !stmt.Unique {
		t.Error("expected Unique to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "IDX_EMP_ID" {
		t.Errorf("expected index name IDX_EMP_ID, got %v", stmt.Name)
	}
	if stmt.Table == nil || stmt.Table.Schema != "HR" || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table HR.EMPLOYEES, got %v", stmt.Table)
	}
}

func TestParseCreateBitmapIndex(t *testing.T) {
	p := newTestParser("BITMAP INDEX idx_status ON orders (status)")
	stmt, parseErr3 := p.parseCreateIndexStmt(0)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if !stmt.Bitmap {
		t.Error("expected Bitmap to be true")
	}
}

func TestParseCreateIndexMultiColumn(t *testing.T) {
	p := newTestParser("INDEX idx_multi ON t (a ASC, b DESC)")
	stmt, parseErr4 := p.parseCreateIndexStmt(0)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt.Columns == nil || stmt.Columns.Len() != 2 {
		t.Fatalf("expected 2 columns, got %d", stmt.Columns.Len())
	}
	col0 := stmt.Columns.Items[0].(*ast.IndexColumn)
	if col0.Dir != ast.SORTBY_ASC {
		t.Errorf("expected ASC for col0, got %d", col0.Dir)
	}
	col1 := stmt.Columns.Items[1].(*ast.IndexColumn)
	if col1.Dir != ast.SORTBY_DESC {
		t.Errorf("expected DESC for col1, got %d", col1.Dir)
	}
}

func TestParseCreateIndexReverse(t *testing.T) {
	p := newTestParser("INDEX idx_rev ON t (a) REVERSE")
	stmt, parseErr5 := p.parseCreateIndexStmt(0)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if !stmt.Reverse {
		t.Error("expected Reverse to be true")
	}
}

func TestParseCreateIndexTablespace(t *testing.T) {
	p := newTestParser("INDEX idx_ts ON t (a) TABLESPACE users")
	stmt, parseErr6 := p.parseCreateIndexStmt(0)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if stmt.Tablespace != "USERS" {
		t.Errorf("expected tablespace USERS, got %q", stmt.Tablespace)
	}
}

func TestParseCreateIndexLocal(t *testing.T) {
	p := newTestParser("INDEX idx_local ON t (a) LOCAL")
	stmt, parseErr7 := p.parseCreateIndexStmt(0)
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if !stmt.Local {
		t.Error("expected Local to be true")
	}
}

func TestParseCreateIndexGlobal(t *testing.T) {
	p := newTestParser("INDEX idx_global ON t (a) GLOBAL")
	stmt, parseErr8 := p.parseCreateIndexStmt(0)
	if parseErr8 != nil {
		t.Fatalf("parse: %v", parseErr8)
	}
	if !stmt.Global {
		t.Error("expected Global to be true")
	}
}

func TestParseCreateIndexOnline(t *testing.T) {
	p := newTestParser("INDEX idx_online ON t (a) ONLINE")
	stmt, parseErr9 := p.parseCreateIndexStmt(0)
	if parseErr9 != nil {
		t.Fatalf("parse: %v", parseErr9)
	}
	if !stmt.Online {
		t.Error("expected Online to be true")
	}
}

func TestParseCreateIndexLoc(t *testing.T) {
	p := newTestParser("INDEX idx ON t (a)")
	stmt, parseErr10 := p.parseCreateIndexStmt(0)
	if parseErr10 != nil {
		t.Fatalf("parse: %v", parseErr10)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}
