package catalog

import (
	"testing"
)

// --- Section 7.2 (starmap): Table State Consistency (6 scenarios) ---
// These tests verify internal catalog consistency after ALTER TABLE operations.

// checkPositionsSequential verifies that all column positions in the table are
// 1-based, sequential, and gap-free.
func checkPositionsSequential(t *testing.T, tbl *Table) {
	t.Helper()
	for i, col := range tbl.Columns {
		expected := i + 1
		if col.Position != expected {
			t.Errorf("column %q: expected position %d, got %d", col.Name, expected, col.Position)
		}
	}
}

// checkColByNameConsistent verifies that the colByName map is consistent with
// the actual Columns slice: every column is findable via GetColumn and the
// returned column matches the one in the slice.
func checkColByNameConsistent(t *testing.T, tbl *Table) {
	t.Helper()
	for i, col := range tbl.Columns {
		got := tbl.GetColumn(col.Name)
		if got == nil {
			t.Errorf("GetColumn(%q) returned nil, but column exists at index %d", col.Name, i)
			continue
		}
		if got != col {
			t.Errorf("GetColumn(%q) returned different pointer than Columns[%d]", col.Name, i)
		}
	}
}

// Scenario 1: After ADD COLUMN, all column positions are sequential (1-based, no gaps)
func TestWalkThrough_11_2_AddColumnPositions(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100), age INT)")

	// Add column at end.
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN email VARCHAR(255)")
	tbl := c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	if len(tbl.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(tbl.Columns))
	}

	// Add column FIRST.
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN flag TINYINT FIRST")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	if len(tbl.Columns) != 5 {
		t.Fatalf("expected 5 columns, got %d", len(tbl.Columns))
	}
	if tbl.Columns[0].Name != "flag" {
		t.Errorf("expected first column 'flag', got %q", tbl.Columns[0].Name)
	}

	// Add column AFTER a specific column.
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN score INT AFTER name")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	if len(tbl.Columns) != 6 {
		t.Fatalf("expected 6 columns, got %d", len(tbl.Columns))
	}
}

// Scenario 2: After DROP COLUMN, remaining column positions are resequenced
func TestWalkThrough_11_2_DropColumnResequence(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (a INT, b INT, c INT, d INT, e INT)")

	// Drop from middle.
	wtExec(t, c, "ALTER TABLE t1 DROP COLUMN c")
	tbl := c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	if len(tbl.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(tbl.Columns))
	}
	// Expected order: a, b, d, e
	expectedNames := []string{"a", "b", "d", "e"}
	for i, name := range expectedNames {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
	}

	// Drop from first.
	wtExec(t, c, "ALTER TABLE t1 DROP COLUMN a")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}

	// Drop from last.
	wtExec(t, c, "ALTER TABLE t1 DROP COLUMN e")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	if len(tbl.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(tbl.Columns))
	}
}

// Scenario 3: After MODIFY COLUMN FIRST, positions reflect new order
func TestWalkThrough_11_2_ModifyColumnFirst(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (a INT, b INT, c INT)")

	// Move last column to first.
	wtExec(t, c, "ALTER TABLE t1 MODIFY COLUMN c INT FIRST")
	tbl := c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	expectedNames := []string{"c", "a", "b"}
	for i, name := range expectedNames {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
	}

	// Move middle column to first.
	wtExec(t, c, "ALTER TABLE t1 MODIFY COLUMN a INT FIRST")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkPositionsSequential(t, tbl)

	expectedNames = []string{"a", "c", "b"}
	for i, name := range expectedNames {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
	}
}

// Scenario 4: After RENAME COLUMN, index column references updated
func TestWalkThrough_11_2_RenameColumnIndexRefs(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100), age INT, INDEX idx_name (name), INDEX idx_combo (name, age))")

	// Rename the column used in indexes.
	wtExec(t, c, "ALTER TABLE t1 RENAME COLUMN name TO full_name")
	tbl := c.GetDatabase("testdb").GetTable("t1")

	// Old name should not be findable.
	if tbl.GetColumn("name") != nil {
		t.Error("old column 'name' should no longer exist")
	}
	// New name should be findable.
	if tbl.GetColumn("full_name") == nil {
		t.Fatal("column 'full_name' not found after rename")
	}

	// Verify index column references reflect the new name.
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if ic.Name == "name" {
				t.Errorf("index %q still references old column name 'name'", idx.Name)
			}
		}
	}

	// Specifically check idx_name.
	var idxName *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			idxName = idx
			break
		}
	}
	if idxName == nil {
		t.Fatal("index idx_name not found")
	}
	if len(idxName.Columns) != 1 || idxName.Columns[0].Name != "full_name" {
		t.Errorf("idx_name: expected column 'full_name', got %v", idxName.Columns)
	}

	// Check composite index idx_combo.
	var idxCombo *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_combo" {
			idxCombo = idx
			break
		}
	}
	if idxCombo == nil {
		t.Fatal("index idx_combo not found")
	}
	if len(idxCombo.Columns) != 2 {
		t.Fatalf("idx_combo: expected 2 columns, got %d", len(idxCombo.Columns))
	}
	if idxCombo.Columns[0].Name != "full_name" {
		t.Errorf("idx_combo column 0: expected 'full_name', got %q", idxCombo.Columns[0].Name)
	}
	if idxCombo.Columns[1].Name != "age" {
		t.Errorf("idx_combo column 1: expected 'age', got %q", idxCombo.Columns[1].Name)
	}

	// Positions should still be consistent.
	checkPositionsSequential(t, tbl)
	checkColByNameConsistent(t, tbl)
}

// Scenario 5: After DROP INDEX, remaining indexes unaffected
func TestWalkThrough_11_2_DropIndexRemainingUnaffected(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100), age INT, INDEX idx_name (name), INDEX idx_age (age), INDEX idx_combo (name, age))")

	// Drop the middle index.
	wtExec(t, c, "ALTER TABLE t1 DROP INDEX idx_age")
	tbl := c.GetDatabase("testdb").GetTable("t1")

	// idx_age should be gone.
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_age" {
			t.Fatal("index idx_age should have been dropped")
		}
	}

	// idx_name should still exist and be intact.
	var idxName *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			idxName = idx
			break
		}
	}
	if idxName == nil {
		t.Fatal("index idx_name should still exist")
	}
	if len(idxName.Columns) != 1 || idxName.Columns[0].Name != "name" {
		t.Errorf("idx_name columns changed unexpectedly")
	}

	// idx_combo should still exist and be intact.
	var idxCombo *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_combo" {
			idxCombo = idx
			break
		}
	}
	if idxCombo == nil {
		t.Fatal("index idx_combo should still exist")
	}
	if len(idxCombo.Columns) != 2 {
		t.Fatalf("idx_combo: expected 2 columns, got %d", len(idxCombo.Columns))
	}
	if idxCombo.Columns[0].Name != "name" || idxCombo.Columns[1].Name != "age" {
		t.Errorf("idx_combo columns changed unexpectedly")
	}

	// Columns should be unaffected.
	checkPositionsSequential(t, tbl)
	if len(tbl.Columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(tbl.Columns))
	}
}

// Scenario 6: colByName index is consistent after every ALTER TABLE operation
func TestWalkThrough_11_2_ColByNameConsistency(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100), age INT)")
	tbl := c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)

	// ADD COLUMN
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN email VARCHAR(255)")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)

	// ADD COLUMN FIRST
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN flag TINYINT FIRST")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)

	// DROP COLUMN
	wtExec(t, c, "ALTER TABLE t1 DROP COLUMN age")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)
	// Dropped column must not be findable.
	if tbl.GetColumn("age") != nil {
		t.Error("dropped column 'age' should not be findable via GetColumn")
	}

	// MODIFY COLUMN (change type, move FIRST)
	wtExec(t, c, "ALTER TABLE t1 MODIFY COLUMN email VARCHAR(500) FIRST")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)

	// CHANGE COLUMN (rename)
	wtExec(t, c, "ALTER TABLE t1 CHANGE COLUMN name full_name VARCHAR(200)")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)
	if tbl.GetColumn("name") != nil {
		t.Error("old column 'name' should not be findable after CHANGE")
	}
	if tbl.GetColumn("full_name") == nil {
		t.Error("new column 'full_name' should be findable after CHANGE")
	}

	// RENAME COLUMN
	wtExec(t, c, "ALTER TABLE t1 RENAME COLUMN full_name TO display_name")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)
	if tbl.GetColumn("full_name") != nil {
		t.Error("old column 'full_name' should not be findable after RENAME")
	}
	if tbl.GetColumn("display_name") == nil {
		t.Error("new column 'display_name' should be findable after RENAME")
	}

	// ADD COLUMN AFTER
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN score INT AFTER email")
	tbl = c.GetDatabase("testdb").GetTable("t1")
	checkColByNameConsistent(t, tbl)
	checkPositionsSequential(t, tbl)
}
