package catalog

import "testing"

// Section 1.2 — Column + Index Interactions
// Tests multi-command ALTER TABLE scenarios where column and index operations interact.

// Scenario: ADD COLUMN x INT, ADD INDEX idx_x (x) — add column then index on it in same ALTER
func TestWalkThrough_5_2_AddColumnThenIndex(t *testing.T) {
	// ADD COLUMN x, ADD INDEX idx_x (x) — the index names a column added in the
	// same statement. TiDB v8.5.5 resolves index columns against the pre-ALTER
	// schema and rejects it (1072 column does not exist: x); the ALTER rolls
	// back, leaving the table unchanged (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100))")
	results, _ := c.Exec("ALTER TABLE t ADD COLUMN x INT, ADD INDEX idx_x (x)", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrKeyColumnDoesNotExist)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.GetColumn("x") != nil {
		t.Error("column x should not exist after rollback")
	}
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_x" {
			t.Error("index idx_x should not exist after rollback")
		}
	}
}

// Scenario: ADD COLUMN x INT, ADD UNIQUE INDEX ux (x) — add column then unique index
func TestWalkThrough_5_2_AddColumnThenUniqueIndex(t *testing.T) {
	// Same-statement reference as AddColumnThenIndex: TiDB rejects (1072) and
	// rolls back (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100))")
	results, _ := c.Exec("ALTER TABLE t ADD COLUMN x INT, ADD UNIQUE INDEX ux (x)", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrKeyColumnDoesNotExist)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.GetColumn("x") != nil {
		t.Error("column x should not exist after rollback")
	}
	for _, idx := range tbl.Indexes {
		if idx.Name == "ux" {
			t.Error("unique index ux should not exist after rollback")
		}
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
	// TiDB v8.5.5 resolves the index column against the pre-ALTER schema, where
	// the renamed name y does not exist, so it rejects (1072 column does not
	// exist: y) and rolls back (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, x INT)")
	results, _ := c.Exec("ALTER TABLE t CHANGE COLUMN x y INT, ADD INDEX idx_y (y)", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrKeyColumnDoesNotExist)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.GetColumn("x") == nil {
		t.Error("column x should still exist after rollback")
	}
	if tbl.GetColumn("y") != nil {
		t.Error("column y should not exist after rollback")
	}
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_y" {
			t.Error("index idx_y should not exist after rollback")
		}
	}
}

// Scenario: ADD COLUMN x INT, ADD PRIMARY KEY (id, x) — add column then include in new PK
func TestWalkThrough_5_2_AddColumnThenPrimaryKey(t *testing.T) {
	// The new PK names x, added in the same statement. TiDB v8.5.5 resolves it
	// against the pre-ALTER schema and rejects (1072 column does not exist: x),
	// rolling back (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT NOT NULL, name VARCHAR(100))")
	results, _ := c.Exec("ALTER TABLE t ADD COLUMN x INT NOT NULL, ADD PRIMARY KEY (id, x)", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrKeyColumnDoesNotExist)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.GetColumn("x") != nil {
		t.Error("column x should not exist after rollback")
	}
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			t.Error("primary key should not exist after rollback")
		}
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
