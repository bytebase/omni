package catalog

import "testing"

// TestWalkThrough_2_3_CreateIndexDupName tests that CREATE INDEX with a duplicate
// index name returns ErrDupKeyName (1061).
func TestWalkThrough_2_3_CreateIndexDupName(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100))")
	wtExec(t, c, "CREATE INDEX idx_name ON t (name)")

	results, err := c.Exec("CREATE INDEX idx_name ON t (id)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupKeyName)
}

// TestWalkThrough_2_3_CreateIndexIfNotExists tests that CREATE INDEX IF NOT EXISTS
// is a syntax error: MySQL's create_index_stmt has no IF NOT EXISTS
// (mysql:8.0.46 and mysql:8.4.13 answer 1064).
func TestWalkThrough_2_3_CreateIndexIfNotExists(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100))")
	wtExec(t, c, "CREATE INDEX idx_name ON t (name)")

	if _, err := c.Exec("CREATE INDEX IF NOT EXISTS idx_name ON t (id)", nil); err == nil {
		t.Fatal("expected a syntax error")
	}
}

// TestWalkThrough_2_3_DropIndexUnknown tests that DROP INDEX on a nonexistent
// index returns ErrCantDropKey (1091).
func TestWalkThrough_2_3_DropIndexUnknown(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")

	results, err := c.Exec("DROP INDEX no_such_idx ON t", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrCantDropKey)
}

// TestWalkThrough_2_3_AlterTableDropIndexUnknown tests that ALTER TABLE DROP INDEX
// on a nonexistent index returns ErrCantDropKey (1091).
func TestWalkThrough_2_3_AlterTableDropIndexUnknown(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")

	results, err := c.Exec("ALTER TABLE t DROP INDEX no_such_idx", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrCantDropKey)
}

// TestWalkThrough_2_3_AlterTableAddUniqueIndexDupName tests that ALTER TABLE ADD
// UNIQUE INDEX with a duplicate index name returns ErrDupKeyName (1061).
func TestWalkThrough_2_3_AlterTableAddUniqueIndexDupName(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100))")
	wtExec(t, c, "CREATE INDEX my_idx ON t (name)")

	results, err := c.Exec("ALTER TABLE t ADD UNIQUE INDEX my_idx (id)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupKeyName)
}
