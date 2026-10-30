package catalog

import "testing"

// TestWalkThrough_2_4_CreateTableFKUnknownTable tests that CREATE TABLE with a FK
// referencing an unknown table returns ErrFKNoRefTable (1824) when fk_checks=1.
func TestWalkThrough_2_4_CreateTableFKUnknownTable(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKNoRefTable)
}

// TestWalkThrough_2_4_CreateTableFKUnknownColumn tests that CREATE TABLE with a FK
// referencing an unknown column returns an error when fk_checks=1.
// MySQL returns ErrFKMissingIndex (1822) when the referenced column has no matching
// index (which is the case when the column doesn't exist).
func TestWalkThrough_2_4_CreateTableFKUnknownColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	results, err := c.Exec("CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(nonexistent))", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// The referenced column "nonexistent" doesn't exist, so there's no matching index.
	assertError(t, results[0].Error, ErrFKMissingIndex)
}

// TestWalkThrough_2_4_DropTableReferencedByFK tests that DROP TABLE on a table
// referenced by a FK returns ErrFKCannotDropParent (3730) when fk_checks=1.
func TestWalkThrough_2_4_DropTableReferencedByFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")

	results, err := c.Exec("DROP TABLE parent", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKCannotDropParent)
}

// TestWalkThrough_2_4_DropTableReferencedByFKWithChecksOff tests that DROP TABLE
// on a table referenced by a FK succeeds when fk_checks=0.
func TestWalkThrough_2_4_DropTableReferencedByFKWithChecksOff(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")
	wtExec(t, c, "SET foreign_key_checks = 0")

	results, err := c.Exec("DROP TABLE parent", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

// TestWalkThrough_2_4_AlterTableAddFKUnknownTable tests that ALTER TABLE ADD FK
// referencing an unknown table returns ErrFKNoRefTable (1824).
func TestWalkThrough_2_4_AlterTableAddFKUnknownTable(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT)")

	results, err := c.Exec("ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (pid) REFERENCES parent(id)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKNoRefTable)
}

// TestWalkThrough_2_4_AlterTableDropColumnUsedInFK tests that ALTER TABLE DROP COLUMN
// on a column used in a FK constraint returns an error.
func TestWalkThrough_2_4_AlterTableDropColumnUsedInFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")

	results, err := c.Exec("ALTER TABLE child DROP COLUMN pid", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// Error code 1828: Cannot drop column needed in FK constraint.
	if results[0].Error == nil {
		t.Fatal("expected error when dropping column used in FK, got nil")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *catalog.Error, got %T: %v", results[0].Error, results[0].Error)
	}
	if catErr.Code != 1828 {
		t.Errorf("expected error code 1828, got %d: %s", catErr.Code, catErr.Message)
	}
}

// TestWalkThrough_2_4_AlterTableAddFKMissingIndex tests that ALTER TABLE ADD FK
// where the referenced table lacks a matching index returns ErrFKMissingIndex (1822).
func TestWalkThrough_2_4_AlterTableAddFKMissingIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT, val INT)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT)")

	// parent.val has no index, so FK should fail.
	results, err := c.Exec("ALTER TABLE child ADD CONSTRAINT fk_val FOREIGN KEY (pid) REFERENCES parent(val)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKMissingIndex)
}

// TestWalkThrough_2_4_AlterTableAddFKIncompatibleColumns tests that ALTER TABLE ADD FK
// where column types are incompatible returns ErrFKIncompatibleColumns (3780).
func TestWalkThrough_2_4_AlterTableAddFKIncompatibleColumns(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid VARCHAR(100))")

	results, err := c.Exec("ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (pid) REFERENCES parent(id)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKIncompatibleColumns)
}

// TestWalkThrough_2_4_SetForeignKeyChecksOff tests that SET foreign_key_checks=0
// disables FK validation during CREATE TABLE.
func TestWalkThrough_2_4_SetForeignKeyChecksOff(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "SET foreign_key_checks = 0")

	// Should succeed even though "parent" doesn't exist.
	results, err := c.Exec("CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

// TestWalkThrough_2_4_SetForeignKeyChecksOn tests that SET foreign_key_checks=1
// re-enables FK validation.
func TestWalkThrough_2_4_SetForeignKeyChecksOn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "SET foreign_key_checks = 0")
	// Create child with FK to nonexistent parent — should succeed.
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")
	wtExec(t, c, "SET foreign_key_checks = 1")

	// Now creating another table referencing nonexistent parent should fail.
	results, err := c.Exec("CREATE TABLE child2 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKNoRefTable)
}
