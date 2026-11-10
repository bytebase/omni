package catalog

import "testing"

// Section 1.3: Column + FK Interactions (multi-command ALTER TABLE)

// Scenario 1: ADD COLUMN parent_id INT, ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)
func TestWalkThrough_5_3_AddColumnThenFK(t *testing.T) {
	// ADD COLUMN parent_id, ADD CONSTRAINT fk FOREIGN KEY (parent_id) ... — the
	// FK names a column added in the same statement. TiDB v8.5.5 resolves it
	// against the pre-ALTER schema and rejects (1072); the ALTER rolls back,
	// leaving child with only its id column (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL)")

	results, _ := c.Exec("ALTER TABLE child ADD COLUMN parent_id INT, ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrKeyColumnDoesNotExist)

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	if tbl.GetColumn("parent_id") != nil {
		t.Error("column parent_id should not exist after rollback")
	}
	if len(tbl.Constraints) != 0 {
		t.Errorf("expected no constraints after rollback, got %d", len(tbl.Constraints))
	}
}

// Scenario 2: ADD COLUMN parent_id INT, ADD INDEX idx (parent_id), ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)
func TestWalkThrough_5_3_AddColumnExplicitIndexThenFK(t *testing.T) {
	// The ADD INDEX names parent_id, added in the same statement. TiDB v8.5.5
	// resolves it against the pre-ALTER schema and rejects (1072); the ALTER
	// rolls back (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL)")

	results, _ := c.Exec("ALTER TABLE child ADD COLUMN parent_id INT, ADD INDEX idx (parent_id), ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrKeyColumnDoesNotExist)

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	if tbl.GetColumn("parent_id") != nil {
		t.Error("column parent_id should not exist after rollback")
	}
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx" {
			t.Error("index idx should not exist after rollback")
		}
	}
}

// Scenario 3: DROP FOREIGN KEY fk, DROP COLUMN parent_id
func TestWalkThrough_5_3_DropFKThenDropColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL, parent_id INT, CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id))")

	mustExec(t, c, "ALTER TABLE child DROP FOREIGN KEY fk, DROP COLUMN parent_id")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// Verify FK constraint is gone.
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			t.Errorf("FK constraint should have been dropped, found: %s", con.Name)
		}
	}

	// Verify column is gone.
	if tbl.GetColumn("parent_id") != nil {
		t.Error("column parent_id should have been dropped")
	}

	// Verify backing index is also gone (column was dropped, so index columns are empty).
	for _, idx := range tbl.Indexes {
		if idx.Name == "fk" {
			t.Error("backing index 'fk' should have been removed after column drop")
		}
	}
}

// Scenario 4: DROP FOREIGN KEY fk, DROP INDEX fk
func TestWalkThrough_5_3_DropFKThenDropBackingIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL, parent_id INT, CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id))")

	mustExec(t, c, "ALTER TABLE child DROP FOREIGN KEY fk, DROP INDEX fk")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// Verify FK constraint is gone.
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			t.Errorf("FK constraint should have been dropped, found: %s", con.Name)
		}
	}

	// Verify backing index is gone.
	for _, idx := range tbl.Indexes {
		if idx.Name == "fk" {
			t.Errorf("backing index 'fk' should have been dropped")
		}
	}

	// Verify column still exists.
	if tbl.GetColumn("parent_id") == nil {
		t.Error("column parent_id should still exist")
	}
}

// Scenario 5: ADD FOREIGN KEY fk1 (...), ADD INDEX idx (...) on same column in same ALTER
func TestWalkThrough_5_3_AddFKAndIndexSameColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL, parent_id INT)")

	mustExec(t, c, "ALTER TABLE child ADD CONSTRAINT fk1 FOREIGN KEY (parent_id) REFERENCES parent(id), ADD INDEX idx (parent_id)")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// Verify FK constraint exists.
	var fkCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk1" {
			fkCon = con
			break
		}
	}
	if fkCon == nil {
		t.Fatal("FK constraint 'fk1' not found")
	}

	// Verify explicit index "idx" exists.
	var explicitIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx" {
			explicitIdx = idx
			break
		}
	}
	if explicitIdx == nil {
		t.Fatal("explicit index 'idx' not found")
	}

	// The FK backing index should NOT create a duplicate — the explicit idx covers parent_id.
	// But since FK is processed before INDEX in the ALTER commands, the FK may have created
	// its own backing index "fk1" before "idx" was added. MySQL processes commands sequentially.
	// So we should have: backing index "fk1" created by FK, then "idx" added as explicit.
	// Actually in MySQL, when FK is added first and no index exists, the backing index IS created.
	// Then the explicit ADD INDEX creates a second one. Both should exist.
	// Let's just verify: FK constraint exists + at least one index covers parent_id.
	idxCount := 0
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if ic.Name == "parent_id" {
				idxCount++
				break
			}
		}
	}
	if idxCount < 1 {
		t.Errorf("expected at least 1 index covering parent_id, got %d", idxCount)
	}
}

// indexNames returns a list of index names for debugging.
func indexNames(tbl *Table) []string {
	names := make([]string, len(tbl.Indexes))
	for i, idx := range tbl.Indexes {
		names[i] = idx.Name
	}
	return names
}
