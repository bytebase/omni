package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseInsertStmt parses the input and returns the first statement as
// *ast.InsertStmt. Returns nil if the statement is not an InsertStmt.
func testParseInsertStmt(t *testing.T, input string) *ast.InsertStmt {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected statement, got none")
	}
	stmt, ok := result.File.Stmts[0].(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected *ast.InsertStmt, got %T", result.File.Stmts[0])
	}
	return stmt
}

// testParseInsertMultiStmt parses and returns the first statement as *ast.InsertMultiStmt.
func testParseInsertMultiStmt(t *testing.T, input string) *ast.InsertMultiStmt {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected statement, got none")
	}
	stmt, ok := result.File.Stmts[0].(*ast.InsertMultiStmt)
	if !ok {
		t.Fatalf("expected *ast.InsertMultiStmt, got %T", result.File.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// 1. Simple INSERT INTO ... VALUES
// ---------------------------------------------------------------------------

func TestInsert_BasicValues(t *testing.T) {
	stmt := testParseInsertStmt(t, "INSERT INTO t VALUES (1, 2, 3)")
	if stmt.Overwrite {
		t.Error("Overwrite should be false")
	}
	if stmt.Target.Name.Name != "t" {
		t.Errorf("target = %q, want %q", stmt.Target.Name.Name, "t")
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("columns = %d, want 0", len(stmt.Columns))
	}
	if len(stmt.Values) != 1 {
		t.Fatalf("rows = %d, want 1", len(stmt.Values))
	}
	if len(stmt.Values[0]) != 3 {
		t.Errorf("values[0] = %d, want 3", len(stmt.Values[0]))
	}
	if stmt.Select != nil {
		t.Error("Select should be nil")
	}
}

// ---------------------------------------------------------------------------
// 2. INSERT INTO ... (cols) VALUES (...)
// ---------------------------------------------------------------------------

func TestInsert_WithColumnList(t *testing.T) {
	stmt := testParseInsertStmt(t, "INSERT INTO employees (id, name, dept) VALUES (1, 'Alice', 'Eng')")
	if len(stmt.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "id" {
		t.Errorf("columns[0] = %q, want %q", stmt.Columns[0].Name, "id")
	}
	if stmt.Columns[1].Name != "name" {
		t.Errorf("columns[1] = %q, want %q", stmt.Columns[1].Name, "name")
	}
	if stmt.Columns[2].Name != "dept" {
		t.Errorf("columns[2] = %q, want %q", stmt.Columns[2].Name, "dept")
	}
}

// ---------------------------------------------------------------------------
// 3. INSERT INTO ... SELECT ...
// ---------------------------------------------------------------------------

func TestInsert_Select(t *testing.T) {
	stmt := testParseInsertStmt(t, "INSERT INTO t2 SELECT a, b FROM t1 WHERE a > 0")
	if stmt.Select == nil {
		t.Fatal("Select should not be nil")
	}
	if len(stmt.Values) != 0 {
		t.Errorf("values should be nil, got %d rows", len(stmt.Values))
	}
	sel, ok := stmt.Select.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Select is %T, want *ast.SelectStmt", stmt.Select)
	}
	if len(sel.Targets) != 2 {
		t.Errorf("targets = %d, want 2", len(sel.Targets))
	}
}

// ---------------------------------------------------------------------------
// 4. INSERT OVERWRITE INTO
// ---------------------------------------------------------------------------

func TestInsert_Overwrite(t *testing.T) {
	stmt := testParseInsertStmt(t, "INSERT OVERWRITE INTO t VALUES (42)")
	if !stmt.Overwrite {
		t.Error("Overwrite should be true")
	}
}

// ---------------------------------------------------------------------------
// 5. INSERT INTO ... multi-row VALUES
// ---------------------------------------------------------------------------

func TestInsert_MultiRowValues(t *testing.T) {
	stmt := testParseInsertStmt(t, "INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')")
	if len(stmt.Values) != 3 {
		t.Fatalf("rows = %d, want 3", len(stmt.Values))
	}
	for i, row := range stmt.Values {
		if len(row) != 2 {
			t.Errorf("row[%d] = %d, want 2", i, len(row))
		}
	}
}

// ---------------------------------------------------------------------------
// 6. INSERT qualified table name
// ---------------------------------------------------------------------------

func TestInsert_QualifiedTableName(t *testing.T) {
	stmt := testParseInsertStmt(t, "INSERT INTO mydb.myschema.mytable VALUES (1)")
	if stmt.Target.Name.Name != "mytable" {
		t.Errorf("name = %q, want %q", stmt.Target.Name.Name, "mytable")
	}
	if stmt.Target.Schema.Name != "myschema" {
		t.Errorf("schema = %q, want %q", stmt.Target.Schema.Name, "myschema")
	}
	if stmt.Target.Database.Name != "mydb" {
		t.Errorf("database = %q, want %q", stmt.Target.Database.Name, "mydb")
	}
}

// ---------------------------------------------------------------------------
// 7. INSERT ALL (unconditional)
// ---------------------------------------------------------------------------

func TestInsertAll_Unconditional(t *testing.T) {
	stmt := testParseInsertMultiStmt(t, `
		INSERT ALL
		  INTO t1 VALUES (1, 'a')
		  INTO t2 VALUES (2, 'b')
		SELECT * FROM src
	`)
	if stmt.First {
		t.Error("First should be false for INSERT ALL")
	}
	if stmt.Overwrite {
		t.Error("Overwrite should be false")
	}
	if len(stmt.Branches) != 2 {
		t.Fatalf("branches = %d, want 2", len(stmt.Branches))
	}
	if stmt.Branches[0].Target.Name.Name != "t1" {
		t.Errorf("branch[0].target = %q, want %q", stmt.Branches[0].Target.Name.Name, "t1")
	}
	if stmt.Branches[1].Target.Name.Name != "t2" {
		t.Errorf("branch[1].target = %q, want %q", stmt.Branches[1].Target.Name.Name, "t2")
	}
	if stmt.Select == nil {
		t.Error("Select should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 8. INSERT FIRST (conditional)
// ---------------------------------------------------------------------------

func TestInsertFirst_Conditional(t *testing.T) {
	stmt := testParseInsertMultiStmt(t, `
		INSERT FIRST
		  WHEN amount > 100 THEN INTO big_sales VALUES (id, amount)
		  WHEN amount > 10  THEN INTO mid_sales VALUES (id, amount)
		  ELSE INTO small_sales VALUES (id, amount)
		SELECT id, amount FROM sales
	`)
	if !stmt.First {
		t.Error("First should be true for INSERT FIRST")
	}
	if len(stmt.Branches) != 3 {
		t.Fatalf("branches = %d, want 3", len(stmt.Branches))
	}
	// First two branches should have WHEN conditions
	if stmt.Branches[0].When == nil {
		t.Error("branch[0].When should not be nil")
	}
	if stmt.Branches[1].When == nil {
		t.Error("branch[1].When should not be nil")
	}
	// ELSE branch has no condition
	if stmt.Branches[2].When != nil {
		t.Error("branch[2].When should be nil (ELSE)")
	}
	if stmt.Select == nil {
		t.Error("Select should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 9. INSERT ALL with column list in branch
// ---------------------------------------------------------------------------

func TestInsertAll_WithColumns(t *testing.T) {
	stmt := testParseInsertMultiStmt(t, `
		INSERT ALL
		  INTO t1 (col1, col2) VALUES (a, b)
		SELECT a, b FROM src
	`)
	if len(stmt.Branches) != 1 {
		t.Fatalf("branches = %d, want 1", len(stmt.Branches))
	}
	if len(stmt.Branches[0].Columns) != 2 {
		t.Errorf("columns = %d, want 2", len(stmt.Branches[0].Columns))
	}
}

// ---------------------------------------------------------------------------
// 10. INSERT ... SELECT (WITH CTE)
// ---------------------------------------------------------------------------

func TestInsert_WithCTESelect(t *testing.T) {
	stmt := testParseInsertStmt(t, `
		INSERT INTO t
		WITH cte AS (SELECT 1 AS x)
		SELECT x FROM cte
	`)
	if stmt.Select == nil {
		t.Fatal("Select should not be nil")
	}
	sel, ok := stmt.Select.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Select is %T, want *ast.SelectStmt", stmt.Select)
	}
	if len(sel.With) == 0 {
		t.Error("expected WITH CTEs")
	}
}
