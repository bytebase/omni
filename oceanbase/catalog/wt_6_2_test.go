package catalog

import (
	"strings"
	"testing"
)

// TestWalkThrough_6_2 tests FK Validation Matrix scenarios from section 2.2
// of the walkthrough. Scenarios already fully covered in wt_2_4_test.go are
// omitted; this file covers new or extended aspects only.

// Scenario 1: DROP TABLE parent when child FK exists, foreign_key_checks=1 — error 3730
// (Already covered by TestWalkThrough_2_4_DropTableReferencedByFK; included here
// under a 6_2 name for completeness of the section proof.)
func TestWalkThrough_6_2_DropParentFKChecksOn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")

	results, err := c.Exec("DROP TABLE parent", &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKCannotDropParent)
}

// Scenario 2: DROP TABLE parent when child FK exists, foreign_key_checks=0 — succeeds,
// child FK becomes orphan (references dropped parent).
func TestWalkThrough_6_2_DropParentFKChecksOff_OrphanState(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")
	wtExec(t, c, "SET foreign_key_checks = 0")
	wtExec(t, c, "DROP TABLE parent")

	// Parent should be gone.
	if c.GetDatabase("testdb").GetTable("parent") != nil {
		t.Fatal("parent table should have been dropped")
	}

	// Child should still exist.
	child := c.GetDatabase("testdb").GetTable("child")
	if child == nil {
		t.Fatal("child table should still exist")
	}

	// Child FK constraint should still reference the dropped parent.
	var fk *Constraint
	for _, con := range child.Constraints {
		if con.Type == ConForeignKey {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("child should still have FK constraint")
	}
	if !strings.EqualFold(fk.RefTable, "parent") {
		t.Errorf("FK should still reference 'parent', got %q", fk.RefTable)
	}
}

// Scenario 3: DROP TABLE child then parent, foreign_key_checks=1 — succeeds.
func TestWalkThrough_6_2_DropChildThenParent(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")

	// Drop child first — removes FK dependency.
	wtExec(t, c, "DROP TABLE child")
	// Now parent can be dropped.
	wtExec(t, c, "DROP TABLE parent")

	if c.GetDatabase("testdb").GetTable("child") != nil {
		t.Error("child should be dropped")
	}
	if c.GetDatabase("testdb").GetTable("parent") != nil {
		t.Error("parent should be dropped")
	}
}

// Scenario 4: DROP COLUMN used in FK on same table, foreign_key_checks=1 — error 1828
// (Already covered by TestWalkThrough_2_4_AlterTableDropColumnUsedInFK; included
// here for section completeness.)
func TestWalkThrough_6_2_DropColumnUsedInFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(id))")

	results, err := c.Exec("ALTER TABLE child DROP COLUMN pid", &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if results[0].Error == nil {
		t.Fatal("expected error when dropping FK column")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != 1828 {
		t.Errorf("expected error code 1828, got %d", catErr.Code)
	}
}

// Scenario 5: CREATE TABLE with FK referencing nonexistent table, foreign_key_checks=0 — succeeds.
// (Already covered by TestWalkThrough_2_4_SetForeignKeyChecksOff; included for completeness.)
func TestWalkThrough_6_2_CreateFKNonexistentTableFKOff(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "SET foreign_key_checks = 0")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES nonexistent(id))")

	child := c.GetDatabase("testdb").GetTable("child")
	if child == nil {
		t.Fatal("child table should exist")
	}
}

// Scenario 6: CREATE TABLE with FK referencing nonexistent column, foreign_key_checks=0 — succeeds.
func TestWalkThrough_6_2_CreateFKNonexistentColumnFKOff(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "SET foreign_key_checks = 0")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, FOREIGN KEY (pid) REFERENCES parent(nonexistent_col))")

	child := c.GetDatabase("testdb").GetTable("child")
	if child == nil {
		t.Fatal("child table should exist")
	}
	// Verify FK constraint is stored.
	var fk *Constraint
	for _, con := range child.Constraints {
		if con.Type == ConForeignKey {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint should be stored")
	}
	if len(fk.RefColumns) == 0 || fk.RefColumns[0] != "nonexistent_col" {
		t.Errorf("FK should reference 'nonexistent_col', got %v", fk.RefColumns)
	}
}

// Scenario 7: ALTER TABLE ADD FK with type mismatch (INT vs VARCHAR), foreign_key_checks=1 — error.
// (Already covered by TestWalkThrough_2_4_AlterTableAddFKIncompatibleColumns; included for completeness.)
func TestWalkThrough_6_2_AddFKTypeMismatch(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT, pid VARCHAR(100))")

	results, err := c.Exec("ALTER TABLE child ADD CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES parent(id)", &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKIncompatibleColumns)
}

// Scenario 8: ALTER TABLE ADD FK where referenced table has no index on referenced columns — error 1822.
// (Already covered by TestWalkThrough_2_4_AlterTableAddFKMissingIndex; included for completeness.)
func TestWalkThrough_6_2_AddFKMissingIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT, val INT)") // val has no index
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT)")

	results, err := c.Exec("ALTER TABLE child ADD CONSTRAINT fk_v FOREIGN KEY (pid) REFERENCES parent(val)", &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKMissingIndex)
}

// Scenario 9: SET foreign_key_checks=0 then CREATE circular FKs then SET foreign_key_checks=1 — both tables valid.
func TestWalkThrough_6_2_CircularFKs(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "SET foreign_key_checks = 0")
	wtExec(t, c, "CREATE TABLE a (id INT PRIMARY KEY, b_id INT, FOREIGN KEY (b_id) REFERENCES b(id))")
	wtExec(t, c, "CREATE TABLE b (id INT PRIMARY KEY, a_id INT, FOREIGN KEY (a_id) REFERENCES a(id))")
	wtExec(t, c, "SET foreign_key_checks = 1")

	tblA := c.GetDatabase("testdb").GetTable("a")
	tblB := c.GetDatabase("testdb").GetTable("b")
	if tblA == nil {
		t.Fatal("table a should exist")
	}
	if tblB == nil {
		t.Fatal("table b should exist")
	}

	// Verify table a has FK referencing b.
	var fkA *Constraint
	for _, con := range tblA.Constraints {
		if con.Type == ConForeignKey {
			fkA = con
			break
		}
	}
	if fkA == nil {
		t.Fatal("table a should have FK constraint")
	}
	if !strings.EqualFold(fkA.RefTable, "b") {
		t.Errorf("table a FK should reference 'b', got %q", fkA.RefTable)
	}

	// Verify table b has FK referencing a.
	var fkB *Constraint
	for _, con := range tblB.Constraints {
		if con.Type == ConForeignKey {
			fkB = con
			break
		}
	}
	if fkB == nil {
		t.Fatal("table b should have FK constraint")
	}
	if !strings.EqualFold(fkB.RefTable, "a") {
		t.Errorf("table b FK should reference 'a', got %q", fkB.RefTable)
	}
}

// Scenario 10: Self-referencing FK (table references itself) — column references own table PK.
func TestWalkThrough_6_2_SelfReferencingFK(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE tree (id INT PRIMARY KEY, parent_id INT, FOREIGN KEY (parent_id) REFERENCES tree(id))")

	tbl := c.GetDatabase("testdb").GetTable("tree")
	if tbl == nil {
		t.Fatal("table tree should exist")
	}

	var fk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("self-referencing FK constraint should exist")
	}
	if !strings.EqualFold(fk.RefTable, "tree") {
		t.Errorf("FK should reference 'tree', got %q", fk.RefTable)
	}
	if len(fk.Columns) != 1 || fk.Columns[0] != "parent_id" {
		t.Errorf("FK columns should be [parent_id], got %v", fk.Columns)
	}
	if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
		t.Errorf("FK ref columns should be [id], got %v", fk.RefColumns)
	}

	// Verify SHOW CREATE TABLE renders the self-ref FK.
	sct := c.ShowCreateTable("testdb", "tree")
	if !strings.Contains(sct, "REFERENCES `tree`") {
		t.Errorf("SHOW CREATE TABLE should contain self-reference, got:\n%s", sct)
	}
}

// Scenario 11: FK column count mismatch (single-column FK referencing composite PK) — error.
// MySQL returns error 1822 (missing index) because the single-column FK cannot match
// the composite PK index (which requires both columns as leading prefix).
func TestWalkThrough_6_2_FKColumnCountMismatch(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id1 INT, id2 INT, PRIMARY KEY (id1, id2))")
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT)")

	// Single-column FK referencing one column of a composite PK.
	// The PK index has (id1, id2), so a reference to just (id1) DOES match
	// the leading prefix. This should SUCCEED in MySQL.
	// But a reference to (id2) alone would NOT match the prefix and would fail.
	// Let's test the failing case: reference id2 (second column of composite PK).
	results, err := c.Exec("ALTER TABLE child ADD CONSTRAINT fk_mismatch FOREIGN KEY (pid) REFERENCES parent(id2)", &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKMissingIndex)
}

// Scenario 12: Cross-database FK: FOREIGN KEY (col) REFERENCES other_db.parent(id) — stored and rendered correctly.
func TestWalkThrough_6_2_CrossDatabaseFK(t *testing.T) {
	c := wtSetup(t)
	// Create the other database and parent table.
	wtExec(t, c, "CREATE DATABASE other_db")
	wtExec(t, c, "CREATE TABLE other_db.parent (id INT PRIMARY KEY)")

	// Create child in testdb with cross-database FK.
	wtExec(t, c, "CREATE TABLE child (id INT, pid INT, CONSTRAINT fk_cross FOREIGN KEY (pid) REFERENCES other_db.parent(id))")

	child := c.GetDatabase("testdb").GetTable("child")
	if child == nil {
		t.Fatal("child table should exist")
	}

	// Verify FK stores the cross-database reference.
	var fk *Constraint
	for _, con := range child.Constraints {
		if con.Type == ConForeignKey {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint should exist")
	}
	if !strings.EqualFold(fk.RefDatabase, "other_db") {
		t.Errorf("FK RefDatabase should be 'other_db', got %q", fk.RefDatabase)
	}
	if !strings.EqualFold(fk.RefTable, "parent") {
		t.Errorf("FK RefTable should be 'parent', got %q", fk.RefTable)
	}

	// Verify SHOW CREATE TABLE renders the cross-database reference.
	sct := c.ShowCreateTable("testdb", "child")
	if !strings.Contains(sct, "`other_db`.`parent`") {
		t.Errorf("SHOW CREATE TABLE should contain cross-database reference `other_db`.`parent`, got:\n%s", sct)
	}
}
