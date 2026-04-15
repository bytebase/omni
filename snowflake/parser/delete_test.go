package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseDeleteStmt parses and returns the first statement as *ast.DeleteStmt.
func testParseDeleteStmt(t *testing.T, input string) *ast.DeleteStmt {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected statement, got none")
	}
	stmt, ok := result.File.Stmts[0].(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("expected *ast.DeleteStmt, got %T", result.File.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// 1. Basic DELETE FROM
// ---------------------------------------------------------------------------

func TestDelete_Basic(t *testing.T) {
	stmt := testParseDeleteStmt(t, "DELETE FROM employees")
	if stmt.Target.Name.Name != "employees" {
		t.Errorf("target = %q, want %q", stmt.Target.Name.Name, "employees")
	}
	if stmt.Where != nil {
		t.Error("Where should be nil")
	}
	if len(stmt.Using) != 0 {
		t.Errorf("using = %d, want 0", len(stmt.Using))
	}
}

// ---------------------------------------------------------------------------
// 2. DELETE FROM ... WHERE
// ---------------------------------------------------------------------------

func TestDelete_WithWhere(t *testing.T) {
	stmt := testParseDeleteStmt(t, "DELETE FROM t WHERE id = 42")
	if stmt.Where == nil {
		t.Fatal("Where should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 3. DELETE with alias
// ---------------------------------------------------------------------------

func TestDelete_WithAlias(t *testing.T) {
	stmt := testParseDeleteStmt(t, "DELETE FROM employees e WHERE e.dept = 'HR'")
	if stmt.Alias.Name != "e" {
		t.Errorf("alias = %q, want %q", stmt.Alias.Name, "e")
	}
	if stmt.Where == nil {
		t.Error("Where should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 4. DELETE ... USING (Snowflake extension)
// ---------------------------------------------------------------------------

func TestDelete_WithUsing(t *testing.T) {
	stmt := testParseDeleteStmt(t, `
		DELETE FROM orders o
		USING customers c
		WHERE o.customer_id = c.id AND c.active = false
	`)
	if len(stmt.Using) != 1 {
		t.Fatalf("using = %d, want 1", len(stmt.Using))
	}
	if stmt.Where == nil {
		t.Error("Where should not be nil")
	}
	// Check alias on the USING source
	ref, ok := stmt.Using[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("using[0] is %T, want *ast.TableRef", stmt.Using[0])
	}
	if ref.Alias.Name != "c" {
		t.Errorf("using[0].alias = %q, want %q", ref.Alias.Name, "c")
	}
}

// ---------------------------------------------------------------------------
// 5. DELETE USING multiple tables
// ---------------------------------------------------------------------------

func TestDelete_UsingMultiple(t *testing.T) {
	stmt := testParseDeleteStmt(t, `
		DELETE FROM t
		USING a, b
		WHERE t.x = a.x AND t.y = b.y
	`)
	if len(stmt.Using) != 2 {
		t.Fatalf("using = %d, want 2", len(stmt.Using))
	}
}

// ---------------------------------------------------------------------------
// 6. DELETE qualified table
// ---------------------------------------------------------------------------

func TestDelete_QualifiedName(t *testing.T) {
	stmt := testParseDeleteStmt(t, "DELETE FROM mydb.myschema.t WHERE id = 1")
	if stmt.Target.Name.Name != "t" {
		t.Errorf("name = %q, want %q", stmt.Target.Name.Name, "t")
	}
	if stmt.Target.Schema.Name != "myschema" {
		t.Errorf("schema = %q, want %q", stmt.Target.Schema.Name, "myschema")
	}
}
