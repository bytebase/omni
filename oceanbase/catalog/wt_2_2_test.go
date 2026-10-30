package catalog

import "testing"

// Section 2.2: Table and Column Errors

func TestWalkThrough_2_2_CreateTableDuplicate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("CREATE TABLE t (id INT)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupTable)
}

func TestWalkThrough_2_2_CreateTableIfNotExistsOnExisting(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("CREATE TABLE IF NOT EXISTS t (id INT)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

func TestWalkThrough_2_2_CreateTableDuplicateColumn(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE TABLE t (id INT, id INT)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupColumn)
}

func TestWalkThrough_2_2_CreateTableMultiplePrimaryKeys(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50) PRIMARY KEY)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrMultiplePriKey)
}

func TestWalkThrough_2_2_DropTableUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP TABLE unknown_tbl", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// DROP TABLE on non-existent table returns ErrUnknownTable (1051).
	assertError(t, results[0].Error, ErrUnknownTable)
}

func TestWalkThrough_2_2_DropTableIfExistsUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP TABLE IF EXISTS unknown_tbl", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

func TestWalkThrough_2_2_AlterTableUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("ALTER TABLE unknown_tbl ADD COLUMN x INT", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchTable)
}

func TestWalkThrough_2_2_AlterTableAddColumnDuplicate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("ALTER TABLE t ADD COLUMN id INT", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupColumn)
}

func TestWalkThrough_2_2_AlterTableDropColumnUnknown(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("ALTER TABLE t DROP COLUMN unknown_col", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// MySQL 8.0 returns error 1091 (ErrCantDropKey) for DROP COLUMN on nonexistent column.
	// The catalog matches this behavior. The scenario listed ErrNoSuchColumn (1054)
	// but the actual MySQL behavior is 1091.
	assertError(t, results[0].Error, ErrCantDropKey)
}

func TestWalkThrough_2_2_AlterTableModifyColumnUnknown(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("ALTER TABLE t MODIFY COLUMN unknown_col VARCHAR(50)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchColumn)
}

func TestWalkThrough_2_2_AlterTableChangeColumnUnknownSource(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("ALTER TABLE t CHANGE COLUMN unknown_col new_col INT", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchColumn)
}

func TestWalkThrough_2_2_AlterTableChangeColumnDuplicateTarget(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(50))")
	results, err := c.Exec("ALTER TABLE t CHANGE COLUMN id name INT", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupColumn)
}

func TestWalkThrough_2_2_AlterTableAddPKWhenPKExists(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50))")
	results, err := c.Exec("ALTER TABLE t ADD PRIMARY KEY (name)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrMultiplePriKey)
}

func TestWalkThrough_2_2_AlterTableRenameColumnUnknown(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	results, err := c.Exec("ALTER TABLE t RENAME COLUMN unknown_col TO new_col", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchColumn)
}

func TestWalkThrough_2_2_RenameTableUnknownSource(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("RENAME TABLE unknown_tbl TO new_tbl", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchTable)
}

func TestWalkThrough_2_2_RenameTableToExistingTarget(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	wtExec(t, c, "CREATE TABLE t2 (id INT)")
	results, err := c.Exec("RENAME TABLE t1 TO t2", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupTable)
}

func TestWalkThrough_2_2_DropTableMultiPartialFailure(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	// DROP TABLE t1, t2 — t1 exists, t2 does not.
	// t1 should be dropped, then error on t2.
	results, err := c.Exec("DROP TABLE t1, t2", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrUnknownTable)

	// Verify t1 was actually dropped before the error on t2.
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("database testdb not found")
	}
	if db.GetTable("t1") != nil {
		t.Error("t1 should have been dropped before error on t2")
	}
}
