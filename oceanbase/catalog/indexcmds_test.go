package catalog

import "testing"

func setupIndexTestTable(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t1 (id INT NOT NULL, name VARCHAR(100), body TEXT)")
	return c
}

func TestCreateIndex(t *testing.T) {
	c := setupIndexTestTable(t)
	mustExec(t, c, "CREATE INDEX idx_name ON t1 (name)")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Indexes) != 1 {
		t.Fatalf("expected 1 index, got %d", len(tbl.Indexes))
	}

	idx := tbl.Indexes[0]
	if idx.Name != "idx_name" {
		t.Errorf("expected index name 'idx_name', got %q", idx.Name)
	}
	if idx.Unique {
		t.Error("expected non-unique index")
	}
	if idx.IndexType != "" {
		t.Errorf("expected empty IndexType (implicit BTREE), got %q", idx.IndexType)
	}
	if len(idx.Columns) != 1 || idx.Columns[0].Name != "name" {
		t.Errorf("expected index on column 'name', got %v", idx.Columns)
	}

	// No constraint should be created for a plain index.
	if len(tbl.Constraints) != 0 {
		t.Errorf("expected 0 constraints, got %d", len(tbl.Constraints))
	}
}

func TestCreateUniqueIndex(t *testing.T) {
	c := setupIndexTestTable(t)
	mustExec(t, c, "CREATE UNIQUE INDEX idx_id ON t1 (id)")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Indexes) != 1 {
		t.Fatalf("expected 1 index, got %d", len(tbl.Indexes))
	}

	idx := tbl.Indexes[0]
	if idx.Name != "idx_id" {
		t.Errorf("expected index name 'idx_id', got %q", idx.Name)
	}
	if !idx.Unique {
		t.Error("expected unique index")
	}

	// Unique index should also create a constraint.
	if len(tbl.Constraints) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(tbl.Constraints))
	}
	con := tbl.Constraints[0]
	if con.Type != ConUniqueKey {
		t.Errorf("expected ConUniqueKey, got %d", con.Type)
	}
	if con.Name != "idx_id" {
		t.Errorf("expected constraint name 'idx_id', got %q", con.Name)
	}
	if con.IndexName != "idx_id" {
		t.Errorf("expected constraint IndexName 'idx_id', got %q", con.IndexName)
	}
}

func TestCreateFulltextIndex(t *testing.T) {
	c := setupIndexTestTable(t)
	mustExec(t, c, "CREATE FULLTEXT INDEX idx_body ON t1 (body)")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Indexes) != 1 {
		t.Fatalf("expected 1 index, got %d", len(tbl.Indexes))
	}

	idx := tbl.Indexes[0]
	if idx.Name != "idx_body" {
		t.Errorf("expected index name 'idx_body', got %q", idx.Name)
	}
	if !idx.Fulltext {
		t.Error("expected fulltext index")
	}
	if idx.IndexType != "FULLTEXT" {
		t.Errorf("expected IndexType 'FULLTEXT', got %q", idx.IndexType)
	}
	if idx.Unique {
		t.Error("fulltext index should not be unique")
	}
}

func TestDropIndex(t *testing.T) {
	c := setupIndexTestTable(t)
	mustExec(t, c, "CREATE INDEX idx_name ON t1 (name)")
	mustExec(t, c, "DROP INDEX idx_name ON t1")

	tbl := c.GetDatabase("test").GetTable("t1")
	if len(tbl.Indexes) != 0 {
		t.Fatalf("expected 0 indexes after drop, got %d", len(tbl.Indexes))
	}
}

func TestDropIndexNotFound(t *testing.T) {
	c := setupIndexTestTable(t)
	results, _ := c.Exec("DROP INDEX nonexistent ON t1", &ExecOptions{ContinueOnError: true})
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("expected error for dropping nonexistent index")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrCantDropKey {
		t.Errorf("expected error code %d, got %d", ErrCantDropKey, catErr.Code)
	}
}

func TestCreateIndexDupKeyName(t *testing.T) {
	c := setupIndexTestTable(t)
	mustExec(t, c, "CREATE INDEX idx_name ON t1 (name)")

	results, _ := c.Exec("CREATE INDEX idx_name ON t1 (id)", &ExecOptions{ContinueOnError: true})
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("expected duplicate key name error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrDupKeyName {
		t.Errorf("expected error code %d, got %d", ErrDupKeyName, catErr.Code)
	}
}
