package catalog

import (
	"testing"
)

func TestExecChangesCreateTable(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE TABLE t1 (
			id integer PRIMARY KEY,
			name text NOT NULL,
			email text UNIQUE
		);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	if r.Changes == nil {
		t.Fatal("expected Changes, got nil")
	}

	// Should have exactly 1 relation added.
	if len(r.Changes.Relations) != 1 {
		t.Fatalf("expected 1 relation change, got %d", len(r.Changes.Relations))
	}
	rel := r.Changes.Relations[0]
	if rel.Action != DiffAdd {
		t.Errorf("relation action=%d, want DiffAdd(%d)", rel.Action, DiffAdd)
	}
	if rel.Name != "t1" {
		t.Errorf("relation name=%q, want %q", rel.Name, "t1")
	}
	if rel.To == nil {
		t.Fatal("DiffAdd relation should have To set")
	}

	// For DiffAdd, full state is in To. Verify columns exist.
	if len(rel.To.Columns) != 3 {
		t.Errorf("expected 3 columns in To, got %d", len(rel.To.Columns))
	}

	// Constraints and indexes for new relations can be queried from the catalog.
	pkCount, uniqCount := 0, 0
	for _, con := range c.ConstraintsOf(rel.To.OID) {
		switch con.Type {
		case ConstraintPK:
			pkCount++
		case ConstraintUnique:
			uniqCount++
		}
	}
	if pkCount != 1 {
		t.Errorf("expected 1 PK constraint on new table, got %d", pkCount)
	}
	if uniqCount != 1 {
		t.Errorf("expected 1 UNIQUE constraint on new table, got %d", uniqCount)
	}

	// Backing indexes.
	if len(c.IndexesOf(rel.To.OID)) < 2 {
		t.Errorf("expected at least 2 indexes on new table, got %d", len(c.IndexesOf(rel.To.OID)))
	}
}

func TestExecChangesCreateSchema(t *testing.T) {
	c := New()

	results, err := c.Exec("CREATE SCHEMA myschema;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes, got nil")
	}
	if len(r.Changes.Schemas) != 1 {
		t.Fatalf("expected 1 schema change, got %d", len(r.Changes.Schemas))
	}
	if r.Changes.Schemas[0].Action != DiffAdd {
		t.Errorf("schema action=%d, want DiffAdd", r.Changes.Schemas[0].Action)
	}
	if r.Changes.Schemas[0].Name != "myschema" {
		t.Errorf("schema name=%q, want %q", r.Changes.Schemas[0].Name, "myschema")
	}
}

func TestExecChangesDropTableCascade(t *testing.T) {
	c := New()

	// Set up two tables with a FK relationship.
	_, err := c.Exec(`
		CREATE TABLE parent (id integer PRIMARY KEY);
		CREATE TABLE child (
			id integer PRIMARY KEY,
			parent_id integer REFERENCES parent(id)
		);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	// DROP parent CASCADE should drop the FK on child too.
	results, err := c.Exec("DROP TABLE parent CASCADE;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	if r.Changes == nil {
		t.Fatal("expected Changes, got nil")
	}

	// Should have at least 2 relation changes: parent dropped + child modified (FK removed).
	parentDropped := false
	childModified := false
	for _, rel := range r.Changes.Relations {
		if rel.Name == "parent" && rel.Action == DiffDrop {
			parentDropped = true
		}
		if rel.Name == "child" && rel.Action == DiffModify {
			childModified = true
			// The FK constraint should be dropped.
			fkDropped := false
			for _, con := range rel.Constraints {
				if con.Action == DiffDrop && con.From != nil && con.From.Type == ConstraintFK {
					fkDropped = true
				}
			}
			if !fkDropped {
				t.Error("child's FK constraint should be dropped in cascade")
			}
		}
	}
	if !parentDropped {
		t.Error("parent table should be in Changes as DiffDrop")
	}
	if !childModified {
		t.Error("child table should be in Changes as DiffModify (FK removed)")
	}
}

func TestExecChangesAlterTable(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE t1 (id integer);", nil)
	if err != nil {
		t.Fatal(err)
	}

	results, err := c.Exec("ALTER TABLE t1 ADD COLUMN name text NOT NULL DEFAULT 'unknown';", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes, got nil")
	}
	if len(r.Changes.Relations) != 1 {
		t.Fatalf("expected 1 relation change, got %d", len(r.Changes.Relations))
	}

	rel := r.Changes.Relations[0]
	if rel.Action != DiffModify {
		t.Errorf("relation action=%d, want DiffModify(%d)", rel.Action, DiffModify)
	}

	// Should have 1 column added.
	colAdded := false
	for _, col := range rel.Columns {
		if col.Action == DiffAdd && col.Name == "name" {
			colAdded = true
		}
	}
	if !colAdded {
		t.Error("column 'name' should be added")
	}
}

func TestExecChangesMultiStatement(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE SCHEMA s1;
		CREATE TABLE s1.t1 (id integer);
		CREATE INDEX t1_idx ON s1.t1 (id);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// Statement 0: CREATE SCHEMA
	if results[0].Changes == nil {
		t.Fatal("stmt 0: expected Changes")
	}
	if len(results[0].Changes.Schemas) != 1 {
		t.Errorf("stmt 0: expected 1 schema change, got %d", len(results[0].Changes.Schemas))
	}

	// Statement 1: CREATE TABLE
	if results[1].Changes == nil {
		t.Fatal("stmt 1: expected Changes")
	}
	if len(results[1].Changes.Relations) != 1 {
		t.Errorf("stmt 1: expected 1 relation change, got %d", len(results[1].Changes.Relations))
	}

	// Statement 2: CREATE INDEX
	if results[2].Changes == nil {
		t.Fatal("stmt 2: expected Changes")
	}
	// Index addition shows up as a relation modification.
	if len(results[2].Changes.Relations) != 1 {
		t.Errorf("stmt 2: expected 1 relation change (index added), got %d", len(results[2].Changes.Relations))
	}
	if results[2].Changes.Relations[0].Action != DiffModify {
		t.Errorf("stmt 2: relation action=%d, want DiffModify", results[2].Changes.Relations[0].Action)
	}
	idxAdded := false
	for _, idx := range results[2].Changes.Relations[0].Indexes {
		if idx.Action == DiffAdd {
			idxAdded = true
		}
	}
	if !idxAdded {
		t.Error("stmt 2: expected an index DiffAdd entry")
	}
}

func TestExecChangesErrorNoChanges(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE t1 (id integer);", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Duplicate table should error — no Changes.
	results, err := c.Exec("CREATE TABLE t1 (id integer);", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error == nil {
		t.Fatal("expected error for duplicate table")
	}
	if results[0].Changes != nil {
		t.Error("Changes should be nil when Error is set")
	}
}

func TestExecChangesSkippedNoChanges(t *testing.T) {
	c := New()

	results, err := c.Exec("SELECT 1;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Changes != nil {
		t.Error("Changes should be nil for skipped (DML) statements")
	}
}

func TestExecChangesNoOpUtility(t *testing.T) {
	c := New()

	// SET, BEGIN, COMMIT are no-ops — should have nil Changes (empty diff).
	results, err := c.Exec("SET client_encoding = 'UTF8';", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Changes != nil {
		t.Error("Changes should be nil for no-op utility statements")
	}
}

func TestExecChangesCreateView(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE t1 (id integer, name text);", nil)
	if err != nil {
		t.Fatal(err)
	}

	results, err := c.Exec("CREATE VIEW v1 AS SELECT id, name FROM t1;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes for CREATE VIEW")
	}
	// View should appear as a relation add.
	viewAdded := false
	for _, rel := range r.Changes.Relations {
		if rel.Action == DiffAdd && rel.Name == "v1" {
			viewAdded = true
		}
	}
	if !viewAdded {
		t.Error("view v1 should appear as DiffAdd in Changes.Relations")
	}
}

func TestExecChangesCreateSequence(t *testing.T) {
	c := New()

	results, err := c.Exec("CREATE SEQUENCE my_seq START 1 INCREMENT 1;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes for CREATE SEQUENCE")
	}
	if len(r.Changes.Sequences) != 1 {
		t.Fatalf("expected 1 sequence change, got %d", len(r.Changes.Sequences))
	}
	if r.Changes.Sequences[0].Action != DiffAdd {
		t.Errorf("sequence action=%d, want DiffAdd", r.Changes.Sequences[0].Action)
	}
	if r.Changes.Sequences[0].Name != "my_seq" {
		t.Errorf("sequence name=%q, want %q", r.Changes.Sequences[0].Name, "my_seq")
	}
}

func TestExecChangesRenameTable(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE old_name (id integer);", nil)
	if err != nil {
		t.Fatal(err)
	}

	results, err := c.Exec("ALTER TABLE old_name RENAME TO new_name;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes for RENAME TABLE")
	}
	// Rename shows up as drop old + add new in diff.
	if len(r.Changes.Relations) < 1 {
		t.Fatal("expected at least 1 relation change for rename")
	}
}

func TestExecChangesIfNotExists(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE t1 (id integer);", nil)
	if err != nil {
		t.Fatal(err)
	}

	// IF NOT EXISTS on existing table — no structural change.
	results, err := c.Exec("CREATE TABLE IF NOT EXISTS t1 (id integer);", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error != nil {
		t.Fatalf("IF NOT EXISTS should not error: %v", results[0].Error)
	}
	// No structural change → Changes should be nil.
	if results[0].Changes != nil {
		t.Error("IF NOT EXISTS on existing table should produce nil Changes")
	}
}

func TestExecChangesDropColumn(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE t1 (id integer, name text, email text);", nil)
	if err != nil {
		t.Fatal(err)
	}

	results, err := c.Exec("ALTER TABLE t1 DROP COLUMN email;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes for DROP COLUMN")
	}
	if len(r.Changes.Relations) != 1 {
		t.Fatalf("expected 1 relation change, got %d", len(r.Changes.Relations))
	}

	colDropped := false
	for _, col := range r.Changes.Relations[0].Columns {
		if col.Action == DiffDrop && col.Name == "email" {
			colDropped = true
		}
	}
	if !colDropped {
		t.Error("column 'email' should be DiffDrop in Changes")
	}
}

func TestExecChangesCreateFunction(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE FUNCTION add_one(x integer) RETURNS integer
		LANGUAGE sql AS 'SELECT x + 1';
	`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Changes == nil {
		t.Fatal("expected Changes for CREATE FUNCTION")
	}
	if len(r.Changes.Functions) != 1 {
		t.Fatalf("expected 1 function change, got %d", len(r.Changes.Functions))
	}
	if r.Changes.Functions[0].Action != DiffAdd {
		t.Errorf("function action=%d, want DiffAdd", r.Changes.Functions[0].Action)
	}
}
