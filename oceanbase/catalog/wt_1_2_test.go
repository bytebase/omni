package catalog

import (
	"testing"
)

// --- Default mode: stop at first error ---

func TestWalkThrough_1_2_DefaultStopsAtFirstError(t *testing.T) {
	c := wtSetup(t)
	// Three statements: first succeeds, second fails (dup table), third should not run.
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// Should have exactly 2 results: success + error (stops there).
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	assertNoError(t, results[0].Error)
	assertError(t, results[1].Error, ErrDupTable)
}

func TestWalkThrough_1_2_DefaultStatementsAfterErrorNotExecuted(t *testing.T) {
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// Only 2 results should be returned — third statement not reached.
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// Verify: no result with Index==2.
	for _, r := range results {
		if r.Index == 2 {
			t.Error("third statement should not have been executed")
		}
	}
}

func TestWalkThrough_1_2_DefaultCatalogReflectsPreError(t *testing.T) {
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"
	_, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	// t1 should exist (created before error).
	if db.GetTable("t1") == nil {
		t.Error("t1 should exist — it was created before the error")
	}
	// t2 should NOT exist (after error).
	if db.GetTable("t2") != nil {
		t.Error("t2 should not exist — execution stopped at error")
	}
}

// --- ContinueOnError mode ---

func TestWalkThrough_1_2_ContinueOnErrorAllAttempted(t *testing.T) {
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"
	opts := &ExecOptions{ContinueOnError: true}
	results, err := c.Exec(sql, opts)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// All 3 statements should be attempted.
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestWalkThrough_1_2_ContinueOnErrorSuccessAfterFailure(t *testing.T) {
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"
	opts := &ExecOptions{ContinueOnError: true}
	_, err := c.Exec(sql, opts)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	// t2 should exist even though t1 dup failed in the middle.
	if db.GetTable("t2") == nil {
		t.Error("t2 should exist — ContinueOnError should continue past failures")
	}
}

func TestWalkThrough_1_2_ContinueOnErrorPerStatementErrors(t *testing.T) {
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT);"
	opts := &ExecOptions{ContinueOnError: true}
	results, err := c.Exec(sql, opts)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// First: success
	assertNoError(t, results[0].Error)
	// Second: dup table error
	assertError(t, results[1].Error, ErrDupTable)
	// Third: success
	assertNoError(t, results[2].Error)
}

func TestWalkThrough_1_2_ContinueOnErrorMultipleErrors(t *testing.T) {
	c := wtSetup(t)
	// Two errors: dup table and alter on unknown table.
	sql := `CREATE TABLE t1 (id INT);
CREATE TABLE t1 (id INT);
ALTER TABLE nosuch ADD COLUMN x INT;
CREATE TABLE t2 (id INT);`
	opts := &ExecOptions{ContinueOnError: true}
	results, err := c.Exec(sql, opts)
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
}

// --- Parse error ---

func TestWalkThrough_1_2_ParseErrorReturnsTopLevelError(t *testing.T) {
	c := wtSetup(t)
	// Intentionally bad SQL that the parser cannot parse.
	results, err := c.Exec("CREATE TABLE ???", nil)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if results != nil {
		t.Errorf("expected nil results on parse error, got %d results", len(results))
	}
}

// --- DELIMITER-containing SQL ---

func TestWalkThrough_1_2_DelimiterSplitting(t *testing.T) {
	c := wtSetup(t)
	sql := `DELIMITER ;;
CREATE TABLE t1 (id INT);;
CREATE TABLE t2 (id INT);;
DELIMITER ;`
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	// Should have 2 results (one per CREATE TABLE).
	if len(results) != 2 {
		t.Fatalf("expected 2 results for DELIMITER SQL, got %d", len(results))
	}
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	if db.GetTable("t1") == nil {
		t.Error("t1 should exist after DELIMITER SQL")
	}
	if db.GetTable("t2") == nil {
		t.Error("t2 should exist after DELIMITER SQL")
	}
}
