package catalog

import "testing"

// --- 5.1 (mapped from starmap 1.1) Column Repositioning Interactions ---

func TestWalkThrough_5_1_AddAfterJustAdded(t *testing.T) {
	// ADD COLUMN x AFTER a, ADD COLUMN y AFTER x — the second command references
	// a column added by the first. TiDB v8.5.5 resolves AFTER against the
	// pre-ALTER schema and rejects it (1054 Unknown column 'x'); the atomic
	// ALTER rolls back, leaving the table unchanged (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	results, _ := c.Exec("ALTER TABLE t ADD COLUMN x INT AFTER a, ADD COLUMN y INT AFTER x", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrNoSuchColumn)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	expected := []string{"a", "b", "c"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns after rollback, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
	}
}

func TestWalkThrough_5_1_AddFirstTwice(t *testing.T) {
	// ADD COLUMN x FIRST, ADD COLUMN y FIRST — both FIRST, y should end up before x
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN x INT FIRST, ADD COLUMN y INT FIRST")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Expected order: y, x, a, b, c
	expected := []string{"y", "x", "a", "b", "c"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
		if tbl.Columns[i].Position != i+1 {
			t.Errorf("column %q: expected position %d, got %d", name, i+1, tbl.Columns[i].Position)
		}
	}
}

func TestWalkThrough_5_1_AddAfterThenDropThat(t *testing.T) {
	// ADD COLUMN x AFTER a, DROP COLUMN a — add after a column that is then dropped in same statement
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN x INT AFTER a, DROP COLUMN a")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Expected order: x, b, c (a was dropped)
	expected := []string{"x", "b", "c"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
		if tbl.Columns[i].Position != i+1 {
			t.Errorf("column %q: expected position %d, got %d", name, i+1, tbl.Columns[i].Position)
		}
	}
}

func TestWalkThrough_5_1_ModifyAfterChain(t *testing.T) {
	// MODIFY COLUMN a INT AFTER c, MODIFY COLUMN b INT AFTER a — chain of AFTER references
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN a INT AFTER c, MODIFY COLUMN b INT AFTER a")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Start: a, b, c
	// After MODIFY a AFTER c: b, c, a
	// After MODIFY b AFTER a: c, a, b
	expected := []string{"c", "a", "b"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
		if tbl.Columns[i].Position != i+1 {
			t.Errorf("column %q: expected position %d, got %d", name, i+1, tbl.Columns[i].Position)
		}
	}
}

func TestWalkThrough_5_1_ModifyFirstThenAddAfter(t *testing.T) {
	// MODIFY COLUMN a INT FIRST, ADD COLUMN x INT AFTER a — FIRST then AFTER the moved column
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	// a is already first, but this tests referencing after MODIFY FIRST
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN c INT FIRST, ADD COLUMN x INT AFTER c")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Start: a, b, c
	// After MODIFY c FIRST: c, a, b
	// After ADD x AFTER c: c, x, a, b
	expected := []string{"c", "x", "a", "b"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
		if tbl.Columns[i].Position != i+1 {
			t.Errorf("column %q: expected position %d, got %d", name, i+1, tbl.Columns[i].Position)
		}
	}
}

func TestWalkThrough_5_1_DropThenReAdd(t *testing.T) {
	// DROP COLUMN a, ADD COLUMN a INT — drop then re-add same column name
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t DROP COLUMN a, ADD COLUMN a INT")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Expected order: b, c, a (dropped from front, re-added at end)
	expected := []string{"b", "c", "a"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
		if tbl.Columns[i].Position != i+1 {
			t.Errorf("column %q: expected position %d, got %d", name, i+1, tbl.Columns[i].Position)
		}
	}
}

func TestWalkThrough_5_1_ChangeThenAddAfterNewName(t *testing.T) {
	// CHANGE COLUMN a a_new INT, ADD COLUMN d INT AFTER a_new — rename then
	// reference the new name. TiDB v8.5.5 resolves AFTER against the pre-ALTER
	// schema, where a_new does not exist, so it is rejected (1054 Unknown column
	// 'a_new') and the ALTER rolls back (measured on a live engine).
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	results, _ := c.Exec("ALTER TABLE t CHANGE COLUMN a a_new INT, ADD COLUMN d INT AFTER a_new", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrNoSuchColumn)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	expected := []string{"a", "b", "c"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns after rollback, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
	}
}

func TestWalkThrough_5_1_ThreeAppends(t *testing.T) {
	// ADD COLUMN x INT, ADD COLUMN y INT, ADD COLUMN z INT — three appends, verify final order
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN x INT, ADD COLUMN y INT, ADD COLUMN z INT")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Expected order: a, b, c, x, y, z
	expected := []string{"a", "b", "c", "x", "y", "z"}
	if len(tbl.Columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(tbl.Columns))
	}
	for i, name := range expected {
		if tbl.Columns[i].Name != name {
			t.Errorf("column %d: expected %q, got %q", i, name, tbl.Columns[i].Name)
		}
		if tbl.Columns[i].Position != i+1 {
			t.Errorf("column %q: expected position %d, got %d", name, i+1, tbl.Columns[i].Position)
		}
	}
}
