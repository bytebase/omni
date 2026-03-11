package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCommentOnTable(t *testing.T) {
	result := ParseAndCheck(t, "COMMENT ON TABLE employees IS 'Employee master table'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CommentStmt)
	if !ok {
		t.Fatalf("expected CommentStmt, got %T", raw.Stmt)
	}
	if stmt.ObjectType != ast.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.ObjectType)
	}
	if stmt.Object == nil || stmt.Object.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %v", stmt.Object)
	}
	if stmt.Comment != "Employee master table" {
		t.Errorf("expected 'Employee master table', got %q", stmt.Comment)
	}
}

func TestParseCommentOnTableSchema(t *testing.T) {
	result := ParseAndCheck(t, "COMMENT ON TABLE hr.employees IS 'HR employees'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommentStmt)
	if stmt.Object.Schema != "HR" || stmt.Object.Name != "EMPLOYEES" {
		t.Errorf("expected HR.EMPLOYEES, got %q.%q", stmt.Object.Schema, stmt.Object.Name)
	}
}

func TestParseCommentOnColumn(t *testing.T) {
	result := ParseAndCheck(t, "COMMENT ON COLUMN employees.salary IS 'Monthly salary'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommentStmt)
	if stmt.ObjectType != ast.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE for column comment, got %d", stmt.ObjectType)
	}
	if stmt.Object == nil || stmt.Object.Name != "EMPLOYEES" {
		t.Errorf("expected object name EMPLOYEES, got %v", stmt.Object)
	}
	if stmt.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %q", stmt.Column)
	}
	if stmt.Comment != "Monthly salary" {
		t.Errorf("expected 'Monthly salary', got %q", stmt.Comment)
	}
}

func TestParseCommentOnColumnSchemaQualified(t *testing.T) {
	result := ParseAndCheck(t, "COMMENT ON COLUMN hr.employees.salary IS 'Monthly salary'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommentStmt)
	if stmt.Object == nil || stmt.Object.Schema != "HR" || stmt.Object.Name != "EMPLOYEES" {
		t.Errorf("expected HR.EMPLOYEES, got %v", stmt.Object)
	}
	if stmt.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %q", stmt.Column)
	}
}

func TestParseCommentOnIndex(t *testing.T) {
	result := ParseAndCheck(t, "COMMENT ON INDEX idx_emp IS 'Primary index'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommentStmt)
	if stmt.ObjectType != ast.OBJECT_INDEX {
		t.Errorf("expected OBJECT_INDEX, got %d", stmt.ObjectType)
	}
	if stmt.Object == nil || stmt.Object.Name != "IDX_EMP" {
		t.Errorf("expected IDX_EMP, got %v", stmt.Object)
	}
	if stmt.Comment != "Primary index" {
		t.Errorf("expected 'Primary index', got %q", stmt.Comment)
	}
}

func TestParseCommentOnMaterializedView(t *testing.T) {
	result := ParseAndCheck(t, "COMMENT ON MATERIALIZED VIEW mv_sales IS 'Sales summary'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommentStmt)
	if stmt.ObjectType != ast.OBJECT_MATERIALIZED_VIEW {
		t.Errorf("expected OBJECT_MATERIALIZED_VIEW, got %d", stmt.ObjectType)
	}
}
