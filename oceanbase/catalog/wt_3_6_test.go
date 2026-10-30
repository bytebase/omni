package catalog

import (
	"strings"
	"testing"
)

// TestWalkThrough_3_6_CreateViewExists verifies that CREATE VIEW adds the view
// to database.Views.
func TestWalkThrough_3_6_CreateViewExists(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100))")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id, name FROM t1")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found in database.Views")
	}
	if v.Name != "v1" {
		t.Errorf("expected view name 'v1', got %q", v.Name)
	}
}

// TestWalkThrough_3_6_CreateViewDefinition verifies that the Definition field
// stores the deparsed SQL (not the raw input).
func TestWalkThrough_3_6_CreateViewDefinition(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100))")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id, name FROM t1")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found")
	}
	// Definition should be the deparsed SELECT, not empty.
	if v.Definition == "" {
		t.Fatal("view definition should not be empty")
	}
	// The deparsed SQL should reference the columns and table.
	def := strings.ToLower(v.Definition)
	if !strings.Contains(def, "select") {
		t.Errorf("definition should contain SELECT, got %q", v.Definition)
	}
	if !strings.Contains(def, "t1") {
		t.Errorf("definition should reference t1, got %q", v.Definition)
	}
}

// TestWalkThrough_3_6_CreateViewAttributes verifies Algorithm, Definer, and
// SqlSecurity are preserved.
func TestWalkThrough_3_6_CreateViewAttributes(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	wtExec(t, c, "CREATE ALGORITHM=MERGE DEFINER=`admin`@`localhost` SQL SECURITY INVOKER VIEW v1 AS SELECT id FROM t1")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found")
	}
	if !strings.EqualFold(v.Algorithm, "MERGE") {
		t.Errorf("expected Algorithm 'MERGE', got %q", v.Algorithm)
	}
	// The parser stores the definer as-is (without backtick-quoting).
	// Backtick formatting is applied only in showCreateView via formatDefiner.
	if !strings.Contains(v.Definer, "admin") || !strings.Contains(v.Definer, "localhost") {
		t.Errorf("expected Definer to contain 'admin' and 'localhost', got %q", v.Definer)
	}
	if !strings.EqualFold(v.SqlSecurity, "INVOKER") {
		t.Errorf("expected SqlSecurity 'INVOKER', got %q", v.SqlSecurity)
	}
}

// TestWalkThrough_3_6_CreateViewDefaultAttributes verifies defaults when
// Algorithm/Definer/SqlSecurity are not specified.
func TestWalkThrough_3_6_CreateViewDefaultAttributes(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id FROM t1")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found")
	}
	// Definer defaults to `root`@`%` per viewcmds.go.
	if v.Definer != "`root`@`%`" {
		t.Errorf("expected default Definer '`root`@`%%`', got %q", v.Definer)
	}
}

// TestWalkThrough_3_6_CreateViewColumns verifies that the column list is
// derived from the SELECT target list.
func TestWalkThrough_3_6_CreateViewColumns(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100), age INT)")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id, name FROM t1")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found")
	}
	if len(v.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d: %v", len(v.Columns), v.Columns)
	}
	if v.Columns[0] != "id" {
		t.Errorf("expected first column 'id', got %q", v.Columns[0])
	}
	if v.Columns[1] != "name" {
		t.Errorf("expected second column 'name', got %q", v.Columns[1])
	}
}

// TestWalkThrough_3_6_CreateOrReplaceView verifies that CREATE OR REPLACE VIEW
// updates an existing view.
func TestWalkThrough_3_6_CreateOrReplaceView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100), age INT)")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id FROM t1")

	// Verify initial state.
	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found after CREATE")
	}
	if len(v.Columns) != 1 {
		t.Fatalf("expected 1 column initially, got %d", len(v.Columns))
	}

	// Replace with a different SELECT.
	wtExec(t, c, "CREATE OR REPLACE VIEW v1 AS SELECT id, name, age FROM t1")

	v = db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found after CREATE OR REPLACE")
	}
	if len(v.Columns) != 3 {
		t.Fatalf("expected 3 columns after replace, got %d: %v", len(v.Columns), v.Columns)
	}
}

// TestWalkThrough_3_6_AlterView verifies that ALTER VIEW updates definition
// and attributes.
func TestWalkThrough_3_6_AlterView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100))")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id FROM t1")

	// ALTER VIEW changes definition and attributes.
	wtExec(t, c, "ALTER VIEW v1 AS SELECT id, name FROM t1")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found after ALTER")
	}
	if len(v.Columns) != 2 {
		t.Fatalf("expected 2 columns after ALTER, got %d", len(v.Columns))
	}
	// Definition should be updated.
	def := strings.ToLower(v.Definition)
	if !strings.Contains(def, "name") {
		t.Errorf("definition after ALTER should reference 'name', got %q", v.Definition)
	}
}

// TestWalkThrough_3_6_DropView verifies that DROP VIEW removes the view from
// database.Views.
func TestWalkThrough_3_6_DropView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	wtExec(t, c, "CREATE VIEW v1 AS SELECT id FROM t1")

	db := c.GetDatabase("testdb")
	if db.Views[toLower("v1")] == nil {
		t.Fatal("view v1 should exist before DROP")
	}

	wtExec(t, c, "DROP VIEW v1")

	if db.Views[toLower("v1")] != nil {
		t.Fatal("view v1 should not exist after DROP")
	}
}

// TestWalkThrough_3_6_ViewReferencingTableColumns verifies that a view
// referencing table columns resolves correctly and the column names appear
// in the view's Columns list.
func TestWalkThrough_3_6_ViewReferencingTableColumns(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		salary DECIMAL(10,2)
	)`)
	wtExec(t, c, "CREATE VIEW emp_names AS SELECT id, first_name, last_name FROM employees")

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("emp_names")]
	if v == nil {
		t.Fatal("view emp_names not found")
	}
	expectedCols := []string{"id", "first_name", "last_name"}
	if len(v.Columns) != len(expectedCols) {
		t.Fatalf("expected %d columns, got %d: %v", len(expectedCols), len(v.Columns), v.Columns)
	}
	for i, want := range expectedCols {
		if v.Columns[i] != want {
			t.Errorf("column %d: expected %q, got %q", i, want, v.Columns[i])
		}
	}

	// Verify the definition references the employees table.
	def := strings.ToLower(v.Definition)
	if !strings.Contains(def, "employees") {
		t.Errorf("definition should reference employees table, got %q", v.Definition)
	}
}
