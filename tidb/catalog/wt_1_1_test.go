package catalog

import "testing"

// Section 1.1: Exec Result Basics

func TestWalkThrough_1_1_EmptySQL(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("", nil)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results, got %d results", len(results))
	}
}

func TestWalkThrough_1_1_WhitespaceOnlySQL(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("   \n\t  \n  ", nil)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results, got %d results", len(results))
	}
}

func TestWalkThrough_1_1_CommentOnlySQL(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("-- this is a comment\n/* block comment */", nil)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results, got %d results", len(results))
	}
}

func TestWalkThrough_1_1_SingleDDL(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE TABLE t1 (id INT)", nil)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	assertNoError(t, results[0].Error)
}

func TestWalkThrough_1_1_MultipleDDL(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT); CREATE TABLE t3 (id INT)", nil)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		assertNoError(t, r.Error)
		if r.Error != nil {
			t.Errorf("result[%d] unexpected error: %v", i, r.Error)
		}
	}
}

func TestWalkThrough_1_1_ResultIndexMatchesPosition(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT); CREATE TABLE t3 (id INT)", nil)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	for i, r := range results {
		if r.Index != i {
			t.Errorf("result[%d].Index = %d, want %d", i, r.Index, i)
		}
	}
}

func TestWalkThrough_1_1_DMLSkipped(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100))")

	dmlStatements := []string{
		"SELECT * FROM t1",
		"INSERT INTO t1 (id, name) VALUES (1, 'test')",
		"UPDATE t1 SET name = 'updated' WHERE id = 1",
		"DELETE FROM t1 WHERE id = 1",
	}

	for _, sql := range dmlStatements {
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("unexpected parse error for %q: %v", sql, err)
		}
		if len(results) != 1 {
			t.Fatalf("expected 1 result for %q, got %d", sql, len(results))
		}
		if !results[0].Skipped {
			t.Errorf("expected Skipped=true for %q", sql)
		}
	}
}

func TestWalkThrough_1_1_DMLDoesNotModifyState(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")

	db := c.GetDatabase("testdb")
	tableCountBefore := len(db.Tables)

	// Execute DML — should not change anything
	_, err := c.Exec("INSERT INTO t1 (id) VALUES (1); SELECT * FROM t1; UPDATE t1 SET id = 2; DELETE FROM t1", nil)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	tableCountAfter := len(db.Tables)
	if tableCountBefore != tableCountAfter {
		t.Errorf("table count changed from %d to %d after DML", tableCountBefore, tableCountAfter)
	}

	// Verify the original table is still there unchanged
	tbl := db.GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 should still exist after DML")
	}
}

func TestWalkThrough_1_1_UnknownStatementsIgnored(t *testing.T) {
	c := wtSetup(t)

	// FLUSH and ANALYZE are parsed but not handled in processUtility — they return nil
	unsupportedStatements := []string{
		"FLUSH TABLES",
		"ANALYZE TABLE t1",
	}

	// Create a table so ANALYZE TABLE has something to reference
	wtExec(t, c, "CREATE TABLE t1 (id INT)")

	for _, sql := range unsupportedStatements {
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("unexpected parse error for %q: %v", sql, err)
		}
		if len(results) != 1 {
			t.Fatalf("expected 1 result for %q, got %d", sql, len(results))
		}
		if results[0].Error != nil {
			t.Errorf("expected nil error for %q, got: %v", sql, results[0].Error)
		}
	}
}
