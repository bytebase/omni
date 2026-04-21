package catalog

import "testing"

// --- 4.3 Error Detection in Migrations ---

// TestWalkThrough_4_3_ContinueOnErrorMigration tests that with ContinueOnError,
// a first CREATE succeeds, a duplicate CREATE fails, and a third ALTER on the
// first table succeeds.
func TestWalkThrough_4_3_ContinueOnErrorMigration(t *testing.T) {
	c := wtSetup(t)

	sql := `CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100));
CREATE TABLE t1 (id INT PRIMARY KEY);
ALTER TABLE t1 ADD COLUMN email VARCHAR(255);`

	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// First CREATE succeeds.
	assertNoError(t, results[0].Error)

	// Duplicate CREATE fails.
	assertError(t, results[1].Error, ErrDupTable)

	// ALTER on first table succeeds (ContinueOnError keeps going).
	assertNoError(t, results[2].Error)

	// Verify final state: t1 has id, name, email.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}
	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}
	if tbl.GetColumn("email") == nil {
		t.Error("column email not found after ContinueOnError ALTER")
	}
}

// TestWalkThrough_4_3_ContinueOnErrorCodes tests that ContinueOnError produces
// the correct error code for each failing statement.
func TestWalkThrough_4_3_ContinueOnErrorCodes(t *testing.T) {
	c := wtSetup(t)

	sql := `CREATE TABLE t1 (id INT PRIMARY KEY);
CREATE TABLE t1 (id INT);
ALTER TABLE no_such_table ADD COLUMN x INT;
CREATE TABLE t2 (id INT PRIMARY KEY);`

	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	assertNoError(t, results[0].Error)
	assertError(t, results[1].Error, ErrDupTable)
	assertError(t, results[2].Error, ErrNoSuchTable)
	assertNoError(t, results[3].Error)

	// Both t1 and t2 should exist.
	db := c.GetDatabase("testdb")
	if db.GetTable("t1") == nil {
		t.Error("table t1 not found")
	}
	if db.GetTable("t2") == nil {
		t.Error("table t2 not found")
	}
}

// TestWalkThrough_4_3_FKCycle tests creating a FK cycle (A refs B, B refs A)
// with fk_checks=0, then verifying state after fk_checks=1.
func TestWalkThrough_4_3_FKCycle(t *testing.T) {
	c := wtSetup(t)

	sql := `SET foreign_key_checks = 0;
CREATE TABLE a (
    id INT PRIMARY KEY,
    b_id INT,
    CONSTRAINT fk_a_b FOREIGN KEY (b_id) REFERENCES b (id)
);
CREATE TABLE b (
    id INT PRIMARY KEY,
    a_id INT,
    CONSTRAINT fk_b_a FOREIGN KEY (a_id) REFERENCES a (id)
);
SET foreign_key_checks = 1;`

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec error on stmt %d: %v", r.Index, r.Error)
		}
	}

	db := c.GetDatabase("testdb")

	// Both tables should exist.
	tblA := db.GetTable("a")
	if tblA == nil {
		t.Fatal("table a not found")
	}
	tblB := db.GetTable("b")
	if tblB == nil {
		t.Fatal("table b not found")
	}

	// Verify FK constraints.
	var fkAB *Constraint
	for _, con := range tblA.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_a_b" {
			fkAB = con
			break
		}
	}
	if fkAB == nil {
		t.Fatal("FK fk_a_b not found on table a")
	}
	if fkAB.RefTable != "b" {
		t.Errorf("fk_a_b RefTable: expected 'b', got %q", fkAB.RefTable)
	}

	var fkBA *Constraint
	for _, con := range tblB.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_b_a" {
			fkBA = con
			break
		}
	}
	if fkBA == nil {
		t.Fatal("FK fk_b_a not found on table b")
	}
	if fkBA.RefTable != "a" {
		t.Errorf("fk_b_a RefTable: expected 'a', got %q", fkBA.RefTable)
	}
}

// TestWalkThrough_4_3_DropTableCascadeError tests that dropping a parent table
// referenced by a FK child produces the correct error.
func TestWalkThrough_4_3_DropTableCascadeError(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `CREATE TABLE parent (id INT PRIMARY KEY);
CREATE TABLE child (
    id INT PRIMARY KEY,
    parent_id INT,
    CONSTRAINT fk_child_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
);`)

	// Attempt to drop parent with FK child referencing it.
	results, err := c.Exec("DROP TABLE parent", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrFKCannotDropParent)

	// Parent should still exist.
	if c.GetDatabase("testdb").GetTable("parent") == nil {
		t.Error("parent table should still exist after failed DROP")
	}
}

// TestWalkThrough_4_3_AlterMissingTableLine tests that ALTER on a missing table
// produces the correct error at the correct line.
func TestWalkThrough_4_3_AlterMissingTableLine(t *testing.T) {
	c := wtSetup(t)

	sql := "CREATE TABLE t1 (id INT PRIMARY KEY);\n" +
		"CREATE TABLE t2 (id INT PRIMARY KEY);\n" +
		"ALTER TABLE no_such_table ADD COLUMN x INT;\n" +
		"CREATE TABLE t3 (id INT PRIMARY KEY);"

	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	// First two succeed.
	assertNoError(t, results[0].Error)
	assertNoError(t, results[1].Error)

	// Third statement fails with correct error.
	assertError(t, results[2].Error, ErrNoSuchTable)

	// Verify the line number is correct (line 3).
	if results[2].Line != 3 {
		t.Errorf("expected error on line 3, got line %d", results[2].Line)
	}

	// Fourth succeeds (ContinueOnError).
	assertNoError(t, results[3].Error)
}

// TestWalkThrough_4_3_MultipleErrorsAllCodes tests that multiple errors in one
// migration all have correct error codes and correct line numbers.
func TestWalkThrough_4_3_MultipleErrorsAllCodes(t *testing.T) {
	c := wtSetup(t)

	sql := "CREATE TABLE t1 (id INT PRIMARY KEY);\n" +
		"CREATE TABLE t1 (id INT);\n" +
		"ALTER TABLE missing ADD COLUMN x INT;\n" +
		"ALTER TABLE t1 DROP COLUMN no_such;\n" +
		"ALTER TABLE t1 ADD COLUMN val TEXT;"

	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}

	// Statement 0: CREATE TABLE t1 succeeds.
	assertNoError(t, results[0].Error)
	if results[0].Line != 1 {
		t.Errorf("stmt 0: expected line 1, got %d", results[0].Line)
	}

	// Statement 1: duplicate table error.
	assertError(t, results[1].Error, ErrDupTable)
	if results[1].Line != 2 {
		t.Errorf("stmt 1: expected line 2, got %d", results[1].Line)
	}

	// Statement 2: no such table error.
	assertError(t, results[2].Error, ErrNoSuchTable)
	if results[2].Line != 3 {
		t.Errorf("stmt 2: expected line 3, got %d", results[2].Line)
	}

	// Statement 3: DROP COLUMN on nonexistent column returns 1091 (same as DROP INDEX in MySQL 8.0).
	assertError(t, results[3].Error, ErrCantDropKey)
	if results[3].Line != 4 {
		t.Errorf("stmt 3: expected line 4, got %d", results[3].Line)
	}

	// Statement 4: succeeds.
	assertNoError(t, results[4].Error)
	if results[4].Line != 5 {
		t.Errorf("stmt 4: expected line 5, got %d", results[4].Line)
	}

	// Final state: t1 has id, val.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}
	if len(tbl.Columns) != 2 {
		t.Fatalf("expected 2 columns (id, val), got %d", len(tbl.Columns))
	}
	if tbl.GetColumn("val") == nil {
		t.Error("column val not found")
	}
}

// TestWalkThrough_4_3_MixedDMLAndDDL tests that DML statements are skipped,
// DDL statements are executed, and the final state is correct.
func TestWalkThrough_4_3_MixedDMLAndDDL(t *testing.T) {
	c := wtSetup(t)

	sql := `CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100));
INSERT INTO t1 VALUES (1, 'Alice');
ALTER TABLE t1 ADD COLUMN email VARCHAR(255);
SELECT * FROM t1;
UPDATE t1 SET name = 'Bob' WHERE id = 1;
CREATE INDEX idx_name ON t1 (name);
DELETE FROM t1 WHERE id = 1;`

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 7 {
		t.Fatalf("expected 7 results, got %d", len(results))
	}

	// Verify DML statements are skipped.
	if !results[1].Skipped {
		t.Error("INSERT should be skipped")
	}
	if !results[3].Skipped {
		t.Error("SELECT should be skipped")
	}
	if !results[4].Skipped {
		t.Error("UPDATE should be skipped")
	}
	if !results[6].Skipped {
		t.Error("DELETE should be skipped")
	}

	// Verify DDL statements succeeded.
	assertNoError(t, results[0].Error) // CREATE TABLE
	assertNoError(t, results[2].Error) // ALTER TABLE ADD COLUMN
	assertNoError(t, results[5].Error) // CREATE INDEX

	// Verify final state: t1 has id, name, email and index idx_name.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}
	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}
	if tbl.GetColumn("email") == nil {
		t.Error("column email not found")
	}

	// Verify index.
	if findIndex(tbl, "idx_name") == nil {
		t.Error("index idx_name not found")
	}
}
