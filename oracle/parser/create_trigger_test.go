package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateTriggerBeforeInsert(t *testing.T) {
	p := newTestParser("TRIGGER my_trigger BEFORE INSERT ON employees BEGIN NULL; END;")
	stmt, parseErr1 := p.parseCreateTriggerStmt(0, false, false, false, false)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateTriggerStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_TRIGGER" {
		t.Errorf("expected trigger name MY_TRIGGER, got %v", stmt.Name)
	}
	if stmt.Timing != ast.TRIGGER_BEFORE {
		t.Errorf("expected TRIGGER_BEFORE, got %d", stmt.Timing)
	}
	if stmt.Events == nil || stmt.Events.Len() != 1 {
		t.Fatalf("expected 1 event, got %d", stmt.Events.Len())
	}
	ev := stmt.Events.Items[0].(*ast.Integer)
	if ast.TriggerEvent(ev.Ival) != ast.TRIGGER_INSERT {
		t.Errorf("expected TRIGGER_INSERT, got %d", ev.Ival)
	}
	if stmt.Table == nil || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table EMPLOYEES, got %v", stmt.Table)
	}
	if stmt.ForEachRow {
		t.Error("expected ForEachRow to be false")
	}
	if stmt.Body == nil {
		t.Error("expected non-nil Body")
	}
}

func TestParseCreateTriggerAfterUpdateForEachRow(t *testing.T) {
	p := newTestParser("TRIGGER audit_trigger AFTER UPDATE ON employees FOR EACH ROW BEGIN NULL; END;")
	stmt, parseErr2 := p.parseCreateTriggerStmt(0, false, false, false, false)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if stmt == nil {
		t.Fatal("expected CreateTriggerStmt, got nil")
	}
	if stmt.Timing != ast.TRIGGER_AFTER {
		t.Errorf("expected TRIGGER_AFTER, got %d", stmt.Timing)
	}
	if stmt.Events == nil || stmt.Events.Len() != 1 {
		t.Fatalf("expected 1 event, got %d", stmt.Events.Len())
	}
	ev := stmt.Events.Items[0].(*ast.Integer)
	if ast.TriggerEvent(ev.Ival) != ast.TRIGGER_UPDATE {
		t.Errorf("expected TRIGGER_UPDATE, got %d", ev.Ival)
	}
	if !stmt.ForEachRow {
		t.Error("expected ForEachRow to be true")
	}
	if stmt.Body == nil {
		t.Error("expected non-nil Body")
	}
}

func TestParseCreateTriggerInsteadOf(t *testing.T) {
	p := newTestParser("TRIGGER instead_trigger INSTEAD OF INSERT ON my_view BEGIN NULL; END;")
	stmt, parseErr3 := p.parseCreateTriggerStmt(0, false, false, false, false)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if stmt == nil {
		t.Fatal("expected CreateTriggerStmt, got nil")
	}
	if stmt.Timing != ast.TRIGGER_INSTEAD_OF {
		t.Errorf("expected TRIGGER_INSTEAD_OF, got %d", stmt.Timing)
	}
	if stmt.Table == nil || stmt.Table.Name != "MY_VIEW" {
		t.Errorf("expected table MY_VIEW, got %v", stmt.Table)
	}
}

func TestParseCreateTriggerMultipleEvents(t *testing.T) {
	p := newTestParser("TRIGGER multi_trigger BEFORE INSERT OR UPDATE OR DELETE ON employees BEGIN NULL; END;")
	stmt, parseErr4 := p.parseCreateTriggerStmt(0, false, false, false, false)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt == nil {
		t.Fatal("expected CreateTriggerStmt, got nil")
	}
	if stmt.Events == nil || stmt.Events.Len() != 3 {
		t.Fatalf("expected 3 events, got %d", stmt.Events.Len())
	}
	ev0 := stmt.Events.Items[0].(*ast.Integer)
	ev1 := stmt.Events.Items[1].(*ast.Integer)
	ev2 := stmt.Events.Items[2].(*ast.Integer)
	if ast.TriggerEvent(ev0.Ival) != ast.TRIGGER_INSERT {
		t.Errorf("expected TRIGGER_INSERT, got %d", ev0.Ival)
	}
	if ast.TriggerEvent(ev1.Ival) != ast.TRIGGER_UPDATE {
		t.Errorf("expected TRIGGER_UPDATE, got %d", ev1.Ival)
	}
	if ast.TriggerEvent(ev2.Ival) != ast.TRIGGER_DELETE {
		t.Errorf("expected TRIGGER_DELETE, got %d", ev2.Ival)
	}
}

func TestParseCreateTriggerWithWhen(t *testing.T) {
	p := newTestParser("TRIGGER when_trigger BEFORE INSERT ON employees FOR EACH ROW WHEN (NEW.salary > 1000) BEGIN NULL; END;")
	stmt, parseErr5 := p.parseCreateTriggerStmt(0, false, false, false, false)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if stmt == nil {
		t.Fatal("expected CreateTriggerStmt, got nil")
	}
	if stmt.When == nil {
		t.Error("expected non-nil When clause")
	}
	if !stmt.ForEachRow {
		t.Error("expected ForEachRow to be true")
	}
}

func TestParseCreateOrReplaceTrigger(t *testing.T) {
	p := newTestParser("TRIGGER hr.my_trigger BEFORE DELETE ON hr.employees BEGIN NULL; END;")
	stmt, parseErr6 := p.parseCreateTriggerStmt(0, true, false, false, false)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if stmt == nil {
		t.Fatal("expected CreateTriggerStmt, got nil")
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
	if stmt.Name == nil || stmt.Name.Schema != "HR" || stmt.Name.Name != "MY_TRIGGER" {
		t.Errorf("expected trigger name HR.MY_TRIGGER, got %v", stmt.Name)
	}
	if stmt.Table == nil || stmt.Table.Schema != "HR" || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table HR.EMPLOYEES, got %v", stmt.Table)
	}
}

func TestParseCreateTriggerFullSQL(t *testing.T) {
	sql := `CREATE OR REPLACE TRIGGER audit_emp
  AFTER INSERT OR UPDATE ON employees
  FOR EACH ROW
  BEGIN
    NULL;
  END;`
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateTriggerStmt)
	if !ok {
		t.Fatalf("expected CreateTriggerStmt, got %T", raw.Stmt)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "AUDIT_EMP" {
		t.Errorf("expected trigger name AUDIT_EMP, got %v", stmt.Name)
	}
	if stmt.Timing != ast.TRIGGER_AFTER {
		t.Errorf("expected TRIGGER_AFTER, got %d", stmt.Timing)
	}
	if stmt.Events == nil || stmt.Events.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", stmt.Events.Len())
	}
	if !stmt.ForEachRow {
		t.Error("expected ForEachRow to be true")
	}
}
