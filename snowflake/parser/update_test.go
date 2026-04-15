package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseUpdateStmt parses and returns the first statement as *ast.UpdateStmt.
func testParseUpdateStmt(t *testing.T, input string) *ast.UpdateStmt {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected statement, got none")
	}
	stmt, ok := result.File.Stmts[0].(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("expected *ast.UpdateStmt, got %T", result.File.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// 1. Basic UPDATE ... SET
// ---------------------------------------------------------------------------

func TestUpdate_BasicSet(t *testing.T) {
	stmt := testParseUpdateStmt(t, "UPDATE employees SET salary = 100000")
	if stmt.Target.Name.Name != "employees" {
		t.Errorf("target = %q, want %q", stmt.Target.Name.Name, "employees")
	}
	if len(stmt.Sets) != 1 {
		t.Fatalf("sets = %d, want 1", len(stmt.Sets))
	}
	if stmt.Sets[0].Column.Name.Name != "salary" {
		t.Errorf("column = %q, want %q", stmt.Sets[0].Column.Name.Name, "salary")
	}
	if stmt.Where != nil {
		t.Error("Where should be nil")
	}
	if len(stmt.From) != 0 {
		t.Errorf("from = %d, want 0", len(stmt.From))
	}
}

// ---------------------------------------------------------------------------
// 2. UPDATE with multiple SET assignments
// ---------------------------------------------------------------------------

func TestUpdate_MultipleSet(t *testing.T) {
	stmt := testParseUpdateStmt(t, "UPDATE t SET a = 1, b = 2, c = 3")
	if len(stmt.Sets) != 3 {
		t.Fatalf("sets = %d, want 3", len(stmt.Sets))
	}
	cols := []string{"a", "b", "c"}
	for i, want := range cols {
		if stmt.Sets[i].Column.Name.Name != want {
			t.Errorf("sets[%d].column = %q, want %q", i, stmt.Sets[i].Column.Name.Name, want)
		}
	}
}

// ---------------------------------------------------------------------------
// 3. UPDATE with WHERE clause
// ---------------------------------------------------------------------------

func TestUpdate_WithWhere(t *testing.T) {
	stmt := testParseUpdateStmt(t, "UPDATE t SET x = 5 WHERE id = 10")
	if stmt.Where == nil {
		t.Fatal("Where should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 4. UPDATE with table alias
// ---------------------------------------------------------------------------

func TestUpdate_WithAlias(t *testing.T) {
	stmt := testParseUpdateStmt(t, "UPDATE employees e SET e.salary = 90000 WHERE e.id = 1")
	if stmt.Alias.Name != "e" {
		t.Errorf("alias = %q, want %q", stmt.Alias.Name, "e")
	}
}

// ---------------------------------------------------------------------------
// 5. UPDATE with FROM clause (Snowflake extension)
// ---------------------------------------------------------------------------

func TestUpdate_WithFrom(t *testing.T) {
	stmt := testParseUpdateStmt(t, `
		UPDATE target_tbl
		SET target_tbl.value = src.value
		FROM source_tbl src
		WHERE target_tbl.id = src.id
	`)
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	if stmt.Where == nil {
		t.Error("Where should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 6. UPDATE with qualified table name
// ---------------------------------------------------------------------------

func TestUpdate_QualifiedName(t *testing.T) {
	stmt := testParseUpdateStmt(t, "UPDATE db.schema.t SET a = 1")
	if stmt.Target.Name.Name != "t" {
		t.Errorf("name = %q, want %q", stmt.Target.Name.Name, "t")
	}
	if stmt.Target.Schema.Name != "schema" {
		t.Errorf("schema = %q, want %q", stmt.Target.Schema.Name, "schema")
	}
	if stmt.Target.Database.Name != "db" {
		t.Errorf("database = %q, want %q", stmt.Target.Database.Name, "db")
	}
}

// ---------------------------------------------------------------------------
// 7. UPDATE with expression in SET value
// ---------------------------------------------------------------------------

func TestUpdate_ExpressionValue(t *testing.T) {
	stmt := testParseUpdateStmt(t, "UPDATE t SET price = price * 1.1 WHERE category = 'A'")
	if len(stmt.Sets) != 1 {
		t.Fatalf("sets = %d, want 1", len(stmt.Sets))
	}
	// The value should be a BinaryExpr (price * 1.1)
	if _, ok := stmt.Sets[0].Value.(*ast.BinaryExpr); !ok {
		t.Errorf("expected *ast.BinaryExpr, got %T", stmt.Sets[0].Value)
	}
}
