package catalog

import (
	"testing"
)

// TestWalkThrough_4_2_AddColumnThenIndex tests that after adding a column,
// an index can be created on it and references the new column.
func TestWalkThrough_4_2_AddColumnThenIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT PRIMARY KEY)")
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN email VARCHAR(255)")
	wtExec(t, c, "CREATE INDEX idx_email ON t1 (email)")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}

	// Verify column exists.
	col := tbl.GetColumn("email")
	if col == nil {
		t.Fatal("column email not found")
	}

	// Verify index exists and references new column.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_email" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_email not found")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "email" {
		t.Fatalf("expected index column [email], got %v", found.Columns)
	}
}

// TestWalkThrough_4_2_RenameColumnIndexRef tests that after renaming a column,
// existing indexes still reference the correct (new) column name.
func TestWalkThrough_4_2_RenameColumnIndexRef(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100))")
	wtExec(t, c, "CREATE INDEX idx_name ON t1 (name)")
	wtExec(t, c, "ALTER TABLE t1 RENAME COLUMN name TO full_name")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}

	// Column should be renamed.
	if tbl.GetColumn("name") != nil {
		t.Error("old column name 'name' should not exist")
	}
	if tbl.GetColumn("full_name") == nil {
		t.Fatal("renamed column 'full_name' not found")
	}

	// Index should reference the new name.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_name not found")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "full_name" {
		t.Fatalf("expected index column [full_name], got %v", indexColNames(found))
	}
}

// TestWalkThrough_4_2_AddColumnThenFK tests that after adding a column,
// a foreign key using that column is correctly created.
func TestWalkThrough_4_2_AddColumnThenFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT PRIMARY KEY)")
	wtExec(t, c, "ALTER TABLE child ADD COLUMN parent_id INT")
	wtExec(t, c, "ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// Find FK constraint.
	var fk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_parent" {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint fk_parent not found")
	}
	if len(fk.Columns) != 1 || fk.Columns[0] != "parent_id" {
		t.Fatalf("expected FK columns [parent_id], got %v", fk.Columns)
	}
	if fk.RefTable != "parent" {
		t.Fatalf("expected RefTable=parent, got %s", fk.RefTable)
	}
	if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
		t.Fatalf("expected FK RefColumns [id], got %v", fk.RefColumns)
	}
}

// TestWalkThrough_4_2_DropIndexRecreate tests dropping an index and re-creating
// it with different columns: the old is gone, the new is present.
func TestWalkThrough_4_2_DropIndexRecreate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT PRIMARY KEY, a INT, b INT, c INT)")
	wtExec(t, c, "CREATE INDEX idx_ab ON t1 (a, b)")

	// Verify initial index.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if findIndex(tbl, "idx_ab") == nil {
		t.Fatal("initial idx_ab not found")
	}

	// Drop and re-create with different columns.
	wtExec(t, c, "DROP INDEX idx_ab ON t1")
	if findIndex(tbl, "idx_ab") != nil {
		t.Fatal("idx_ab should be gone after DROP")
	}

	wtExec(t, c, "CREATE INDEX idx_ab ON t1 (b, c)")
	idx := findIndex(tbl, "idx_ab")
	if idx == nil {
		t.Fatal("re-created idx_ab not found")
	}
	names := indexColNames(idx)
	if len(names) != 2 || names[0] != "b" || names[1] != "c" {
		t.Fatalf("expected index columns [b, c], got %v", names)
	}
}

// TestWalkThrough_4_2_AlterTableMultipleCommands tests that multiple ALTER TABLE
// sub-commands applied in sequence produce the correct cumulative effect.
func TestWalkThrough_4_2_AlterTableMultipleCommands(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT PRIMARY KEY, a INT, b INT)")

	// Multiple ALTERs in sequence (separate statements, not multi-command).
	wtExec(t, c, "ALTER TABLE t1 ADD COLUMN c VARCHAR(50)")
	wtExec(t, c, "ALTER TABLE t1 DROP COLUMN b")
	wtExec(t, c, "ALTER TABLE t1 MODIFY COLUMN a BIGINT")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}

	// Should have: id, a, c
	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}
	if tbl.Columns[0].Name != "id" {
		t.Errorf("col 0: expected id, got %s", tbl.Columns[0].Name)
	}
	if tbl.Columns[1].Name != "a" {
		t.Errorf("col 1: expected a, got %s", tbl.Columns[1].Name)
	}
	if tbl.Columns[2].Name != "c" {
		t.Errorf("col 2: expected c, got %s", tbl.Columns[2].Name)
	}

	// a should be BIGINT now.
	colA := tbl.GetColumn("a")
	if colA.DataType != "bigint" {
		t.Errorf("expected column a type=bigint, got %s", colA.DataType)
	}

	// b should be gone.
	if tbl.GetColumn("b") != nil {
		t.Error("column b should have been dropped")
	}
}

// TestWalkThrough_4_2_RenameTableThenView tests renaming a table then creating
// a view on the new name: the view should resolve.
func TestWalkThrough_4_2_RenameTableThenView(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE old_t (id INT PRIMARY KEY, val TEXT)")
	wtExec(t, c, "RENAME TABLE old_t TO new_t")

	db := c.GetDatabase("testdb")
	if db.GetTable("old_t") != nil {
		t.Error("old_t should no longer exist")
	}
	if db.GetTable("new_t") == nil {
		t.Fatal("new_t should exist after rename")
	}

	wtExec(t, c, "CREATE VIEW v1 AS SELECT id, val FROM new_t")

	v := db.Views[toLower("v1")]
	if v == nil {
		t.Fatal("view v1 not found")
	}
	// The view should exist in the database and have a definition.
	if v.Definition == "" {
		t.Error("view v1 should have a non-empty definition")
	}
}

// TestWalkThrough_4_2_ChangeColumnTypeGeneratedColumn tests that after changing
// a column type, a dependent generated column is still recorded.
func TestWalkThrough_4_2_ChangeColumnTypeGeneratedColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		a INT,
		b INT GENERATED ALWAYS AS (a * 2) STORED
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t1")
	colB := tbl.GetColumn("b")
	if colB == nil {
		t.Fatal("column b not found")
	}
	if colB.Generated == nil {
		t.Fatal("column b should be generated")
	}
	if !colB.Generated.Stored {
		t.Error("column b should be STORED")
	}

	// Change column a type to BIGINT.
	wtExec(t, c, "ALTER TABLE t1 MODIFY COLUMN a BIGINT")

	// Verify generated column b is still recorded.
	colB = tbl.GetColumn("b")
	if colB == nil {
		t.Fatal("column b not found after modify")
	}
	if colB.Generated == nil {
		t.Fatal("column b should still be generated after modifying a")
	}
	if !colB.Generated.Stored {
		t.Error("column b should still be STORED")
	}
}

// TestWalkThrough_4_2_ConvertCharset tests CONVERT TO CHARACTER SET updates
// all string columns.
func TestWalkThrough_4_2_ConvertCharset(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100),
		bio TEXT,
		age INT,
		tag ENUM('a','b')
	)`)

	// Convert to latin1.
	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET latin1")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}

	// Table-level charset/collation should be updated.
	if tbl.Charset != "latin1" {
		t.Errorf("expected table charset=latin1, got %s", tbl.Charset)
	}

	// String columns should have latin1.
	for _, colName := range []string{"name", "bio", "tag"} {
		col := tbl.GetColumn(colName)
		if col == nil {
			t.Fatalf("column %s not found", colName)
		}
		if col.Charset != "latin1" {
			t.Errorf("column %s: expected charset=latin1, got %s", colName, col.Charset)
		}
	}

	// Non-string columns should not have charset changed.
	colID := tbl.GetColumn("id")
	if colID.Charset != "" {
		t.Errorf("INT column id should not have charset, got %s", colID.Charset)
	}
	colAge := tbl.GetColumn("age")
	if colAge.Charset != "" {
		t.Errorf("INT column age should not have charset, got %s", colAge.Charset)
	}
}

// --- helpers ---

func findIndex(tbl *Table, name string) *Index {
	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == toLower(name) {
			return idx
		}
	}
	return nil
}

func indexColNames(idx *Index) []string {
	names := make([]string, len(idx.Columns))
	for i, ic := range idx.Columns {
		names[i] = ic.Name
	}
	return names
}
