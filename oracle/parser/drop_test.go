package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseDropTable(t *testing.T) {
	result := ParseAndCheck(t, "DROP TABLE employees")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	drop, ok := raw.Stmt.(*ast.DropStmt)
	if !ok {
		t.Fatalf("expected DropStmt, got %T", raw.Stmt)
	}
	if drop.ObjectType != ast.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", drop.ObjectType)
	}
	if drop.Names == nil || drop.Names.Len() != 1 {
		t.Fatalf("expected 1 name, got %v", drop.Names)
	}
	name := drop.Names.Items[0].(*ast.ObjectName)
	if name.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %q", name.Name)
	}
	if drop.IfExists {
		t.Error("expected IfExists=false")
	}
	if drop.Cascade {
		t.Error("expected Cascade=false")
	}
	if drop.Purge {
		t.Error("expected Purge=false")
	}
}

func TestParseDropTableIfExists(t *testing.T) {
	result := ParseAndCheck(t, "DROP TABLE IF EXISTS employees")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if !drop.IfExists {
		t.Error("expected IfExists=true")
	}
	name := drop.Names.Items[0].(*ast.ObjectName)
	if name.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %q", name.Name)
	}
}

func TestParseDropTableCascadeConstraints(t *testing.T) {
	result := ParseAndCheck(t, "DROP TABLE hr.employees CASCADE CONSTRAINTS")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if !drop.Cascade {
		t.Error("expected Cascade=true")
	}
	name := drop.Names.Items[0].(*ast.ObjectName)
	if name.Schema != "HR" || name.Name != "EMPLOYEES" {
		t.Errorf("expected HR.EMPLOYEES, got %q.%q", name.Schema, name.Name)
	}
}

func TestParseDropTablePurge(t *testing.T) {
	result := ParseAndCheck(t, "DROP TABLE employees PURGE")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if !drop.Purge {
		t.Error("expected Purge=true")
	}
}

func TestParseDropIndex(t *testing.T) {
	result := ParseAndCheck(t, "DROP INDEX idx_emp_name")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_INDEX {
		t.Errorf("expected OBJECT_INDEX, got %d", drop.ObjectType)
	}
	name := drop.Names.Items[0].(*ast.ObjectName)
	if name.Name != "IDX_EMP_NAME" {
		t.Errorf("expected IDX_EMP_NAME, got %q", name.Name)
	}
}

func TestParseDropView(t *testing.T) {
	result := ParseAndCheck(t, "DROP VIEW emp_view")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_VIEW {
		t.Errorf("expected OBJECT_VIEW, got %d", drop.ObjectType)
	}
}

func TestParseDropSequence(t *testing.T) {
	result := ParseAndCheck(t, "DROP SEQUENCE emp_seq")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_SEQUENCE {
		t.Errorf("expected OBJECT_SEQUENCE, got %d", drop.ObjectType)
	}
}

func TestParseDropSynonym(t *testing.T) {
	result := ParseAndCheck(t, "DROP SYNONYM emp_syn")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_SYNONYM {
		t.Errorf("expected OBJECT_SYNONYM, got %d", drop.ObjectType)
	}
}

func TestParseDropPackage(t *testing.T) {
	result := ParseAndCheck(t, "DROP PACKAGE my_pkg")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_PACKAGE {
		t.Errorf("expected OBJECT_PACKAGE, got %d", drop.ObjectType)
	}
}

func TestParseDropPackageBody(t *testing.T) {
	result := ParseAndCheck(t, "DROP PACKAGE BODY my_pkg")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_PACKAGE_BODY {
		t.Errorf("expected OBJECT_PACKAGE_BODY, got %d", drop.ObjectType)
	}
}

func TestParseDropProcedure(t *testing.T) {
	result := ParseAndCheck(t, "DROP PROCEDURE my_proc")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_PROCEDURE {
		t.Errorf("expected OBJECT_PROCEDURE, got %d", drop.ObjectType)
	}
}

func TestParseDropFunction(t *testing.T) {
	result := ParseAndCheck(t, "DROP FUNCTION my_func")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_FUNCTION {
		t.Errorf("expected OBJECT_FUNCTION, got %d", drop.ObjectType)
	}
}

func TestParseDropTrigger(t *testing.T) {
	result := ParseAndCheck(t, "DROP TRIGGER my_trigger")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_TRIGGER {
		t.Errorf("expected OBJECT_TRIGGER, got %d", drop.ObjectType)
	}
}

func TestParseDropType(t *testing.T) {
	result := ParseAndCheck(t, "DROP TYPE my_type")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_TYPE {
		t.Errorf("expected OBJECT_TYPE, got %d", drop.ObjectType)
	}
}

func TestParseDropTypeBody(t *testing.T) {
	result := ParseAndCheck(t, "DROP TYPE BODY my_type")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_TYPE_BODY {
		t.Errorf("expected OBJECT_TYPE_BODY, got %d", drop.ObjectType)
	}
}

func TestParseDropMaterializedView(t *testing.T) {
	result := ParseAndCheck(t, "DROP MATERIALIZED VIEW mv_emp")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_MATERIALIZED_VIEW {
		t.Errorf("expected OBJECT_MATERIALIZED_VIEW, got %d", drop.ObjectType)
	}
}

func TestParseDropDatabaseLink(t *testing.T) {
	result := ParseAndCheck(t, "DROP DATABASE LINK remote_db")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if drop.ObjectType != ast.OBJECT_DATABASE_LINK {
		t.Errorf("expected OBJECT_DATABASE_LINK, got %d", drop.ObjectType)
	}
}

func TestParseDropTableCascadePurge(t *testing.T) {
	result := ParseAndCheck(t, "DROP TABLE employees CASCADE CONSTRAINTS PURGE")
	raw := result.Items[0].(*ast.RawStmt)
	drop := raw.Stmt.(*ast.DropStmt)
	if !drop.Cascade {
		t.Error("expected Cascade=true")
	}
	if !drop.Purge {
		t.Error("expected Purge=true")
	}
}
