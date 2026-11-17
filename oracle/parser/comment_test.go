package parser

import (
	"errors"
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

// TestParseCommentOnAcceptedCategories holds the parser to the categories
// the pinned Oracle (testdata/pins.tsv, 23.26.3) reads on COMMENT ON: for
// each, DBMS_SQL.PARSE of "COMMENT ON <category> a IS 's'" answers an error
// about the name and not ORA-32594, so the category was read. A view is
// commented through TABLE, which the dictionary then reports as
// TABLE_TYPE = 'VIEW'.
func TestParseCommentOnAcceptedCategories(t *testing.T) {
	tests := []struct {
		sql    string
		object string
		column string
	}{
		{"COMMENT ON TABLE employees IS 'a table'", "EMPLOYEES", ""},
		{"COMMENT ON TABLE view_commented_as_a_table IS 'a view'", "VIEW_COMMENTED_AS_A_TABLE", ""},
		{"COMMENT ON COLUMN employees.salary IS 'a column'", "EMPLOYEES", "SALARY"},
		{"COMMENT ON MATERIALIZED VIEW mv_sales IS 'a materialized view'", "MV_SALES", ""},
		{"COMMENT ON EDITION e1 IS 'an edition'", "E1", ""},
		{"COMMENT ON INDEXTYPE it IS 'an indextype'", "IT", ""},
		{"COMMENT ON OPERATOR op IS 'an operator'", "OP", ""},
		{"COMMENT ON AUDIT POLICY p IS 'an audit policy'", "P", ""},
		{"COMMENT ON MINING MODEL m IS 'a mining model'", "M", ""},
		{"COMMENT ON PROPERTY GRAPH g IS 'a property graph'", "G", ""},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			stmt, ok := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CommentStmt)
			if !ok {
				t.Fatalf("expected CommentStmt, got %T", result.Items[0].(*ast.RawStmt).Stmt)
			}
			if stmt.Object == nil || stmt.Object.Name != tc.object || stmt.Column != tc.column {
				t.Errorf("expected %q.%q, got %v.%q", tc.object, tc.column, stmt.Object, stmt.Column)
			}
		})
	}
}

// TestParseCommentOnRefusedCategories holds the parser to the categories
// Oracle refuses on COMMENT ON: ORA-32594, "invalid object category", at the
// category word, which is byte 11 of the statement. A grammar that takes one
// of them accepts what the server answers with an error.
func TestParseCommentOnRefusedCategories(t *testing.T) {
	categories := []string{
		"VIEW", "INDEX", "SEQUENCE", "FUNCTION", "PROCEDURE", "PACKAGE", "TRIGGER",
		"TYPE", "SYNONYM", "DATABASE LINK", "DIMENSION", "CLUSTER", "CONSTRAINT",
		"SCHEMA", "LIBRARY", "JAVA SOURCE", "DIRECTORY", "ROLE", "USER", "TABLESPACE", "DOMAIN",
	}
	for _, category := range categories {
		sql := "COMMENT ON " + category + " a IS 's'"
		_, err := Parse(sql)
		if err == nil {
			t.Errorf("Parse(%q): expected a rejection, got success", sql)
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Position != len("COMMENT ON ") {
			t.Errorf("Parse(%q): expected a rejection at %d, got %v", sql, len("COMMENT ON "), err)
		}
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
