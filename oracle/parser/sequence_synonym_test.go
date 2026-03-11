package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateSequence(t *testing.T) {
	p := newTestParser("SEQUENCE emp_seq")
	stmt := p.parseCreateSequenceStmt(0)
	if stmt == nil {
		t.Fatal("expected CreateSequenceStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "EMP_SEQ" {
		t.Errorf("expected sequence name EMP_SEQ, got %v", stmt.Name)
	}
}

func TestParseCreateSequenceWithOptions(t *testing.T) {
	p := newTestParser("SEQUENCE hr.emp_seq INCREMENT BY 1 START WITH 100 MAXVALUE 999999 MINVALUE 1 CYCLE CACHE 20 ORDER")
	stmt := p.parseCreateSequenceStmt(0)
	if stmt.Name == nil || stmt.Name.Schema != "HR" || stmt.Name.Name != "EMP_SEQ" {
		t.Errorf("expected HR.EMP_SEQ, got %v", stmt.Name)
	}
	if stmt.IncrementBy == nil {
		t.Error("expected non-nil IncrementBy")
	} else {
		n, ok := stmt.IncrementBy.(*ast.NumberLiteral)
		if !ok || n.Ival != 1 {
			t.Errorf("expected IncrementBy=1, got %v", stmt.IncrementBy)
		}
	}
	if stmt.StartWith == nil {
		t.Error("expected non-nil StartWith")
	} else {
		n := stmt.StartWith.(*ast.NumberLiteral)
		if n.Ival != 100 {
			t.Errorf("expected StartWith=100, got %v", n.Ival)
		}
	}
	if stmt.MaxValue == nil {
		t.Error("expected non-nil MaxValue")
	}
	if stmt.MinValue == nil {
		t.Error("expected non-nil MinValue")
	}
	if !stmt.Cycle {
		t.Error("expected Cycle to be true")
	}
	if stmt.Cache == nil {
		t.Error("expected non-nil Cache")
	}
	if !stmt.Order {
		t.Error("expected Order to be true")
	}
}

func TestParseCreateSequenceNoOptions(t *testing.T) {
	p := newTestParser("SEQUENCE s NOMAXVALUE NOMINVALUE NOCYCLE NOCACHE NOORDER")
	stmt := p.parseCreateSequenceStmt(0)
	if !stmt.NoMaxValue {
		t.Error("expected NoMaxValue to be true")
	}
	if !stmt.NoMinValue {
		t.Error("expected NoMinValue to be true")
	}
	if !stmt.NoCycle {
		t.Error("expected NoCycle to be true")
	}
	if !stmt.NoCache {
		t.Error("expected NoCache to be true")
	}
	if !stmt.NoOrder {
		t.Error("expected NoOrder to be true")
	}
}

func TestParseCreateSequenceLoc(t *testing.T) {
	p := newTestParser("SEQUENCE s")
	stmt := p.parseCreateSequenceStmt(0)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

func TestParseCreateSynonym(t *testing.T) {
	p := newTestParser("SYNONYM emp FOR hr.employees")
	stmt := p.parseCreateSynonymStmt(0, false, false)
	if stmt == nil {
		t.Fatal("expected CreateSynonymStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "EMP" {
		t.Errorf("expected synonym name EMP, got %v", stmt.Name)
	}
	if stmt.Target == nil || stmt.Target.Schema != "HR" || stmt.Target.Name != "EMPLOYEES" {
		t.Errorf("expected target HR.EMPLOYEES, got %v", stmt.Target)
	}
}

func TestParseCreatePublicSynonym(t *testing.T) {
	p := newTestParser("SYNONYM emp FOR hr.employees")
	stmt := p.parseCreateSynonymStmt(0, false, true)
	if !stmt.Public {
		t.Error("expected Public to be true")
	}
}

func TestParseCreateOrReplaceSynonym(t *testing.T) {
	p := newTestParser("SYNONYM emp FOR hr.employees")
	stmt := p.parseCreateSynonymStmt(0, true, false)
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

func TestParseCreateSynonymLoc(t *testing.T) {
	p := newTestParser("SYNONYM s FOR t")
	stmt := p.parseCreateSynonymStmt(0, false, false)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

func TestParseCreateDatabaseLink(t *testing.T) {
	p := newTestParser("DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY secret USING 'prod_server'")
	stmt := p.parseCreateDatabaseLinkStmt(0, false)
	if stmt == nil {
		t.Fatal("expected CreateDatabaseLinkStmt, got nil")
	}
	if stmt.Name != "REMOTE_DB" {
		t.Errorf("expected link name REMOTE_DB, got %q", stmt.Name)
	}
	if stmt.ConnectTo != "ADMIN" {
		t.Errorf("expected ConnectTo ADMIN, got %q", stmt.ConnectTo)
	}
	if stmt.Identified != "SECRET" {
		t.Errorf("expected Identified SECRET, got %q", stmt.Identified)
	}
	if stmt.Using != "prod_server" {
		t.Errorf("expected Using prod_server, got %q", stmt.Using)
	}
}

func TestParseCreatePublicDatabaseLink(t *testing.T) {
	p := newTestParser("DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY pass USING 'srv'")
	stmt := p.parseCreateDatabaseLinkStmt(0, true)
	if !stmt.Public {
		t.Error("expected Public to be true")
	}
}

func TestParseCreateDatabaseLinkLoc(t *testing.T) {
	p := newTestParser("DATABASE LINK db CONNECT TO u IDENTIFIED BY p USING 's'")
	stmt := p.parseCreateDatabaseLinkStmt(0, false)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}
