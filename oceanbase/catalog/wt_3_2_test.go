package catalog

import "testing"

// TestWalkThrough_3_2_AlterDatabaseCharset verifies that ALTER DATABASE
// changes the database charset and derives the default collation.
func TestWalkThrough_3_2_AlterDatabaseCharset(t *testing.T) {
	c := wtSetup(t)

	// Default charset is utf8mb4
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	if db.Charset != "utf8mb4" {
		t.Fatalf("expected default charset utf8mb4, got %s", db.Charset)
	}

	wtExec(t, c, "ALTER DATABASE testdb CHARACTER SET latin1")

	db = c.GetDatabase("testdb")
	if db.Charset != "latin1" {
		t.Errorf("expected charset latin1, got %s", db.Charset)
	}
	// When charset is changed without explicit collation, default collation is derived.
	if db.Collation != "latin1_swedish_ci" {
		t.Errorf("expected collation latin1_swedish_ci, got %s", db.Collation)
	}
}

// TestWalkThrough_3_2_AlterDatabaseCollation verifies that ALTER DATABASE
// changes the database collation.
func TestWalkThrough_3_2_AlterDatabaseCollation(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, "ALTER DATABASE testdb COLLATE utf8mb4_unicode_ci")

	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	if db.Collation != "utf8mb4_unicode_ci" {
		t.Errorf("expected collation utf8mb4_unicode_ci, got %s", db.Collation)
	}
	// Charset should remain unchanged when only collation is set.
	if db.Charset != "utf8mb4" {
		t.Errorf("expected charset utf8mb4, got %s", db.Charset)
	}
}

// TestWalkThrough_3_2_RenameTable verifies that RENAME TABLE removes the old
// name and adds the new name with the same columns and indexes.
func TestWalkThrough_3_2_RenameTable(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, "CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100))")
	wtExec(t, c, "CREATE INDEX idx_name ON t1 (name)")

	// Capture column and index info before rename.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("t1 not found before rename")
	}
	origColCount := len(tbl.Columns)
	origIdxCount := len(tbl.Indexes)

	wtExec(t, c, "RENAME TABLE t1 TO t2")

	// Old name must be gone.
	if c.GetDatabase("testdb").GetTable("t1") != nil {
		t.Error("t1 should not exist after rename")
	}

	// New name must be present.
	tbl2 := c.GetDatabase("testdb").GetTable("t2")
	if tbl2 == nil {
		t.Fatal("t2 not found after rename")
	}

	// Same columns.
	if len(tbl2.Columns) != origColCount {
		t.Errorf("expected %d columns, got %d", origColCount, len(tbl2.Columns))
	}
	if tbl2.GetColumn("id") == nil {
		t.Error("column 'id' missing after rename")
	}
	if tbl2.GetColumn("name") == nil {
		t.Error("column 'name' missing after rename")
	}

	// Same indexes.
	if len(tbl2.Indexes) != origIdxCount {
		t.Errorf("expected %d indexes, got %d", origIdxCount, len(tbl2.Indexes))
	}
}

// TestWalkThrough_3_2_RenameTableCrossDatabase verifies that RENAME TABLE
// moves a table from one database to another.
func TestWalkThrough_3_2_RenameTableCrossDatabase(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, "CREATE DATABASE otherdb")
	wtExec(t, c, "CREATE TABLE t1 (id INT PRIMARY KEY, val TEXT)")

	wtExec(t, c, "RENAME TABLE testdb.t1 TO otherdb.t1")

	// Old location must be empty.
	if c.GetDatabase("testdb").GetTable("t1") != nil {
		t.Error("t1 should not exist in testdb after cross-db rename")
	}

	// New location must have the table.
	tbl := c.GetDatabase("otherdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("t1 not found in otherdb after cross-db rename")
	}
	if tbl.GetColumn("id") == nil {
		t.Error("column 'id' missing after cross-db rename")
	}
	if tbl.GetColumn("val") == nil {
		t.Error("column 'val' missing after cross-db rename")
	}
}

// TestWalkThrough_3_2_TruncateTable verifies that TRUNCATE TABLE keeps the
// table but resets AUTO_INCREMENT to 0.
func TestWalkThrough_3_2_TruncateTable(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, "CREATE TABLE t1 (id INT AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=100")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("t1 not found")
	}
	if tbl.AutoIncrement != 100 {
		t.Fatalf("expected AUTO_INCREMENT=100 before truncate, got %d", tbl.AutoIncrement)
	}

	wtExec(t, c, "TRUNCATE TABLE t1")

	tbl = c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("t1 should still exist after truncate")
	}
	if tbl.AutoIncrement != 0 {
		t.Errorf("expected AUTO_INCREMENT=0 after truncate, got %d", tbl.AutoIncrement)
	}
	// Columns should still be present.
	if tbl.GetColumn("id") == nil {
		t.Error("column 'id' missing after truncate")
	}
}
