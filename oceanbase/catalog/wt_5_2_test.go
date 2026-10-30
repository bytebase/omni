package catalog

import "testing"

// Section 1.2 — Column + Index Interactions
// Tests multi-command ALTER TABLE scenarios where column and index operations interact.

// Scenario: ADD COLUMN x INT, ADD INDEX idx_x (x) — add column then index on it in same ALTER
func TestWalkThrough_5_2_AddColumnThenIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100))")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN x INT, ADD INDEX idx_x (x)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Verify column was added.
	col := tbl.GetColumn("x")
	if col == nil {
		t.Fatal("column x not found")
	}

	// Verify index was created.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_x" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_x not found")
	}
	if found.Unique {
		t.Error("idx_x should not be unique")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "x" {
		t.Errorf("index columns mismatch: %+v", found.Columns)
	}
}

// Scenario: ADD COLUMN x INT, ADD UNIQUE INDEX ux (x) — add column then unique index
func TestWalkThrough_5_2_AddColumnThenUniqueIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100))")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN x INT, ADD UNIQUE INDEX ux (x)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Verify column was added.
	col := tbl.GetColumn("x")
	if col == nil {
		t.Fatal("column x not found")
	}

	// Verify unique index was created.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "ux" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("unique index ux not found")
	}
	if !found.Unique {
		t.Error("ux should be unique")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "x" {
		t.Errorf("unique index columns mismatch: %+v", found.Columns)
	}
}

// Scenario: DROP COLUMN x, DROP INDEX idx_x — drop column and its index simultaneously
func TestWalkThrough_5_2_DropColumnAndIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), x INT, INDEX idx_x (x))")

	// Verify setup: column and index exist.
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.GetColumn("x") == nil {
		t.Fatal("column x should exist before drop")
	}

	wtExec(t, c, "ALTER TABLE t DROP COLUMN x, DROP INDEX idx_x")

	tbl = c.GetDatabase("testdb").GetTable("t")

	// Verify column was dropped.
	if tbl.GetColumn("x") != nil {
		t.Error("column x should have been dropped")
	}

	// Verify index was dropped.
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_x" {
			t.Error("index idx_x should have been dropped")
		}
	}

	// Verify remaining columns are correct.
	if len(tbl.Columns) != 2 {
		t.Errorf("expected 2 columns after drop, got %d", len(tbl.Columns))
	}
}

// Scenario: MODIFY COLUMN x VARCHAR(200), DROP INDEX idx_x, ADD INDEX idx_x (x) — rebuild index after type change
func TestWalkThrough_5_2_ModifyColumnRebuildIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, x VARCHAR(100), INDEX idx_x (x))")

	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN x VARCHAR(200), DROP INDEX idx_x, ADD INDEX idx_x (x)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Verify column type was modified.
	col := tbl.GetColumn("x")
	if col == nil {
		t.Fatal("column x not found")
	}
	if col.ColumnType != "varchar(200)" {
		t.Errorf("expected column type 'varchar(200)', got %q", col.ColumnType)
	}

	// Verify index was recreated.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_x" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_x not found after rebuild")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "x" {
		t.Errorf("index columns mismatch: %+v", found.Columns)
	}
}

// Scenario: CHANGE COLUMN x y INT, ADD INDEX idx_y (y) — rename column then index with new name
func TestWalkThrough_5_2_ChangeColumnThenIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, x INT)")

	wtExec(t, c, "ALTER TABLE t CHANGE COLUMN x y INT, ADD INDEX idx_y (y)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Verify old column is gone and new column exists.
	if tbl.GetColumn("x") != nil {
		t.Error("column x should no longer exist after CHANGE")
	}
	col := tbl.GetColumn("y")
	if col == nil {
		t.Fatal("column y not found after CHANGE")
	}

	// Verify index on new column name.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_y" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_y not found")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "y" {
		t.Errorf("index columns mismatch: %+v", found.Columns)
	}
}

// Scenario: ADD COLUMN x INT, ADD PRIMARY KEY (id, x) — add column then include in new PK
func TestWalkThrough_5_2_AddColumnThenPrimaryKey(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100))")

	wtExec(t, c, "ALTER TABLE t ADD COLUMN x INT NOT NULL, ADD PRIMARY KEY (id, x)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Verify column was added.
	col := tbl.GetColumn("x")
	if col == nil {
		t.Fatal("column x not found")
	}
	if col.Nullable {
		t.Error("column x should be NOT NULL")
	}

	// Verify primary key exists with both columns.
	var pkIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			pkIdx = idx
			break
		}
	}
	if pkIdx == nil {
		t.Fatal("primary key not found")
	}
	if len(pkIdx.Columns) != 2 {
		t.Fatalf("expected 2 PK columns, got %d", len(pkIdx.Columns))
	}
	if pkIdx.Columns[0].Name != "id" {
		t.Errorf("expected first PK column 'id', got %q", pkIdx.Columns[0].Name)
	}
	if pkIdx.Columns[1].Name != "x" {
		t.Errorf("expected second PK column 'x', got %q", pkIdx.Columns[1].Name)
	}
}

// Scenario: DROP INDEX idx_x, ADD INDEX idx_x (x, y) — drop and recreate index with extra column
func TestWalkThrough_5_2_DropAndRecreateIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, x INT, y INT, INDEX idx_x (x))")

	wtExec(t, c, "ALTER TABLE t DROP INDEX idx_x, ADD INDEX idx_x (x, y)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Verify index was recreated with two columns.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_x" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_x not found after recreate")
	}
	if len(found.Columns) != 2 {
		t.Fatalf("expected 2 index columns, got %d", len(found.Columns))
	}
	if found.Columns[0].Name != "x" {
		t.Errorf("expected first index column 'x', got %q", found.Columns[0].Name)
	}
	if found.Columns[1].Name != "y" {
		t.Errorf("expected second index column 'y', got %q", found.Columns[1].Name)
	}
}
