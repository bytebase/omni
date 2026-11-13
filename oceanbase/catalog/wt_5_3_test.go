package catalog

import "testing"

// Section 1.3: Column + FK Interactions (multi-command ALTER TABLE)

// Scenario 1: ADD COLUMN parent_id INT, ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)
func TestWalkThrough_5_3_AddColumnThenFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL)")

	mustExec(t, c, "ALTER TABLE child ADD COLUMN parent_id INT, ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// Verify column was added.
	col := tbl.GetColumn("parent_id")
	if col == nil {
		t.Fatal("column parent_id not found")
	}

	// Verify FK constraint exists.
	var fkCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk" {
			fkCon = con
			break
		}
	}
	if fkCon == nil {
		t.Fatal("FK constraint 'fk' not found")
	}
	if len(fkCon.Columns) != 1 || fkCon.Columns[0] != "parent_id" {
		t.Errorf("FK columns = %v, want [parent_id]", fkCon.Columns)
	}
	if fkCon.RefTable != "parent" {
		t.Errorf("FK ref table = %q, want 'parent'", fkCon.RefTable)
	}

	// No index is made for the foreign key: OceanBase makes none.
	if len(tbl.Indexes) != 0 {
		t.Errorf("table child declares no index and has %d", len(tbl.Indexes))
	}
}

// Scenario 2: ADD COLUMN parent_id INT, ADD INDEX idx (parent_id), ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)
func TestWalkThrough_5_3_AddColumnExplicitIndexThenFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL)")

	mustExec(t, c, "ALTER TABLE child ADD COLUMN parent_id INT, ADD INDEX idx (parent_id), ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id)")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	// Verify column was added.
	col := tbl.GetColumn("parent_id")
	if col == nil {
		t.Fatal("column parent_id not found")
	}

	// Verify FK constraint exists.
	var fkCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk" {
			fkCon = con
			break
		}
	}
	if fkCon == nil {
		t.Fatal("FK constraint 'fk' not found")
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

	// Verify the declared index is the only one over the column: the foreign key adds none.
	idxCount := 0
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if ic.Name == "parent_id" {
				idxCount++
				break
			}
		}
	}
	if idxCount != 1 {
		t.Errorf("expected exactly 1 index covering parent_id, got %d (indexes: %v)", idxCount, indexNames(tbl))
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

	// Verify no index was left behind.
	if len(tbl.Indexes) != 0 {
		t.Errorf("table child declares no index and has %v", indexNames(tbl))
	}
}

// Scenario 4: DROP FOREIGN KEY fk, DROP INDEX fk. There is no index fk to drop: OceanBase
// makes none for the foreign key and answers 1091 (measured on 4.3.5.6), and the table
// keeps its foreign key.
func TestWalkThrough_5_3_DropFKThenDropIndexNeverMade(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
	wtExec(t, c, "CREATE TABLE child (id INT NOT NULL, parent_id INT, CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id))")

	assertError(t, execError(t, c, "ALTER TABLE child DROP FOREIGN KEY fk, DROP INDEX fk"), 1091)

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}
	kept := false
	for _, con := range tbl.Constraints {
		kept = kept || con.Type == ConForeignKey && con.Name == "fk"
	}
	if !kept {
		t.Error("the refused statement dropped the foreign key")
	}
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

	// The declared index is the only one: the foreign key adds none, whichever comes first.
	if names := indexNames(tbl); len(names) != 1 {
		t.Errorf("indexes %v, want [idx]", names)
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
