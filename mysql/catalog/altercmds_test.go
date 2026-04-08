package catalog

import "testing"

func setupTestTable(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100), age INT)")
	return c
}

func TestAlterTableAddColumn(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD COLUMN email VARCHAR(255) NOT NULL")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(tbl.Columns))
	}

	col := tbl.GetColumn("email")
	if col == nil {
		t.Fatal("column email not found")
	}
	if col.Nullable {
		t.Error("email should not be nullable")
	}
	if col.ColumnType != "varchar(255)" {
		t.Errorf("expected column type 'varchar(255)', got %q", col.ColumnType)
	}
	if col.Position != 4 {
		t.Errorf("expected position 4, got %d", col.Position)
	}

	// Adding duplicate column should fail.
	results, _ := c.Exec("ALTER TABLE t1 ADD COLUMN email VARCHAR(100)", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate column error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrDupColumn {
		t.Errorf("expected error code %d, got %d", ErrDupColumn, catErr.Code)
	}
}

func TestAlterTableAddColumnMulti(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD COLUMN (email VARCHAR(255) NOT NULL, score INT)")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Columns) != 5 {
		t.Fatalf("expected 5 columns, got %d", len(tbl.Columns))
	}

	col := tbl.GetColumn("email")
	if col == nil {
		t.Fatal("column email not found")
	}
	if col.Nullable {
		t.Error("email should not be nullable")
	}
	if col.ColumnType != "varchar(255)" {
		t.Errorf("expected column type 'varchar(255)', got %q", col.ColumnType)
	}

	col2 := tbl.GetColumn("score")
	if col2 == nil {
		t.Fatal("column score not found")
	}
}

func TestAlterTableDropColumn(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 DROP COLUMN age")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(tbl.Columns))
	}

	if tbl.GetColumn("age") != nil {
		t.Error("column age should have been dropped")
	}

	// Check remaining columns have correct positions.
	id := tbl.GetColumn("id")
	if id.Position != 1 {
		t.Errorf("expected id position 1, got %d", id.Position)
	}
	name := tbl.GetColumn("name")
	if name.Position != 2 {
		t.Errorf("expected name position 2, got %d", name.Position)
	}

	// Dropping non-existent column should fail.
	results, _ := c.Exec("ALTER TABLE t1 DROP COLUMN nonexistent", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected no such column error")
	}
}

func TestAlterTableModifyColumn(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 MODIFY COLUMN name VARCHAR(200) NOT NULL")

	tbl := c.GetDatabase("test").GetTable("t1")
	col := tbl.GetColumn("name")
	if col == nil {
		t.Fatal("column name not found")
	}
	if col.ColumnType != "varchar(200)" {
		t.Errorf("expected column type 'varchar(200)', got %q", col.ColumnType)
	}
	if col.Nullable {
		t.Error("name should not be nullable after MODIFY")
	}
	if col.Position != 2 {
		t.Errorf("expected position 2, got %d", col.Position)
	}
}

func TestAlterTableModifyColumnAfterLater(t *testing.T) {
	// Regression: MODIFY COLUMN ... AFTER <col> panicked with slice bounds
	// out of range when the AFTER column was positioned after the modified column.
	c := setupTestTable(t)
	// t1 has: id(1), name(2), age(3)
	// Move id AFTER age → name(1), age(2), id(3)
	mustExec(t, c, "ALTER TABLE t1 MODIFY COLUMN id INT NOT NULL AFTER age")

	tbl := c.GetDatabase("test").GetTable("t1")
	if tbl.Columns[0].Name != "name" {
		t.Errorf("expected column 0 to be 'name', got %q", tbl.Columns[0].Name)
	}
	if tbl.Columns[1].Name != "age" {
		t.Errorf("expected column 1 to be 'age', got %q", tbl.Columns[1].Name)
	}
	if tbl.Columns[2].Name != "id" {
		t.Errorf("expected column 2 to be 'id', got %q", tbl.Columns[2].Name)
	}
}

func TestAlterTableChangeColumn(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 CHANGE COLUMN name full_name VARCHAR(200)")

	tbl := c.GetDatabase("test").GetTable("t1")

	if tbl.GetColumn("name") != nil {
		t.Error("old column 'name' should no longer exist")
	}

	col := tbl.GetColumn("full_name")
	if col == nil {
		t.Fatal("column full_name not found")
	}
	if col.ColumnType != "varchar(200)" {
		t.Errorf("expected column type 'varchar(200)', got %q", col.ColumnType)
	}
	if col.Position != 2 {
		t.Errorf("expected position 2, got %d", col.Position)
	}

	// Changing to a name that already exists should fail.
	results, _ := c.Exec("ALTER TABLE t1 CHANGE COLUMN age full_name INT", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate column error")
	}
}

func TestAlterTableAddIndex(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD INDEX idx_name (name)")

	tbl := c.GetDatabase("test").GetTable("t1")

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
	if found.Unique {
		t.Error("idx_name should not be unique")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "name" {
		t.Error("idx_name should have column 'name'")
	}
}

func TestAlterTableDropIndex(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD INDEX idx_name (name)")
	mustExec(t, c, "ALTER TABLE t1 DROP INDEX idx_name")

	tbl := c.GetDatabase("test").GetTable("t1")
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			t.Fatal("index idx_name should have been dropped")
		}
	}
}

func TestAlterTableAddPrimaryKey(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD PRIMARY KEY (id)")

	tbl := c.GetDatabase("test").GetTable("t1")

	var pkIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			pkIdx = idx
			break
		}
	}
	if pkIdx == nil {
		t.Fatal("primary key index not found")
	}
	if pkIdx.Name != "PRIMARY" {
		t.Errorf("expected PK index name 'PRIMARY', got %q", pkIdx.Name)
	}
	if !pkIdx.Unique {
		t.Error("PK should be unique")
	}

	// id column should be NOT NULL.
	id := tbl.GetColumn("id")
	if id.Nullable {
		t.Error("id should not be nullable after adding PK")
	}

	// Check PK constraint.
	var pkCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			pkCon = con
			break
		}
	}
	if pkCon == nil {
		t.Fatal("PK constraint not found")
	}

	// Adding another PK should fail.
	results, _ := c.Exec("ALTER TABLE t1 ADD PRIMARY KEY (name)", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected multiple primary key error")
	}
}

func TestAlterTableRenameColumn(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 RENAME COLUMN name TO full_name")

	tbl := c.GetDatabase("test").GetTable("t1")

	if tbl.GetColumn("name") != nil {
		t.Error("old column 'name' should no longer exist")
	}
	col := tbl.GetColumn("full_name")
	if col == nil {
		t.Fatal("column full_name not found")
	}
	if col.Position != 2 {
		t.Errorf("expected position 2, got %d", col.Position)
	}

	// Renaming to existing name should fail.
	results, _ := c.Exec("ALTER TABLE t1 RENAME COLUMN full_name TO id", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate column error")
	}
}

func TestAlterTableAddColumnFirst(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD COLUMN email VARCHAR(255) FIRST")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(tbl.Columns))
	}

	email := tbl.GetColumn("email")
	if email == nil {
		t.Fatal("column email not found")
	}
	if email.Position != 1 {
		t.Errorf("expected email at position 1, got %d", email.Position)
	}

	id := tbl.GetColumn("id")
	if id.Position != 2 {
		t.Errorf("expected id at position 2, got %d", id.Position)
	}

	name := tbl.GetColumn("name")
	if name.Position != 3 {
		t.Errorf("expected name at position 3, got %d", name.Position)
	}

	age := tbl.GetColumn("age")
	if age.Position != 4 {
		t.Errorf("expected age at position 4, got %d", age.Position)
	}
}

func TestAlterTableAddColumnAfter(t *testing.T) {
	c := setupTestTable(t)
	mustExec(t, c, "ALTER TABLE t1 ADD COLUMN email VARCHAR(255) AFTER id")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(tbl.Columns))
	}

	id := tbl.GetColumn("id")
	if id.Position != 1 {
		t.Errorf("expected id at position 1, got %d", id.Position)
	}

	email := tbl.GetColumn("email")
	if email == nil {
		t.Fatal("column email not found")
	}
	if email.Position != 2 {
		t.Errorf("expected email at position 2, got %d", email.Position)
	}

	name := tbl.GetColumn("name")
	if name.Position != 3 {
		t.Errorf("expected name at position 3, got %d", name.Position)
	}

	age := tbl.GetColumn("age")
	if age.Position != 4 {
		t.Errorf("expected age at position 4, got %d", age.Position)
	}
}
