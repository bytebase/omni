package ddl

import (
	"testing"
)

func TestExecAllSuccess(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE SCHEMA s1;
		CREATE TABLE s1.t1 (id int);
		CREATE INDEX t1_id_idx ON s1.t1 (id);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("stmt %d: unexpected error: %v", i, r.Error)
		}
		if r.Skipped {
			t.Errorf("stmt %d: should not be skipped", i)
		}
		if r.Index != i {
			t.Errorf("stmt %d: Index=%d, want %d", i, r.Index, i)
		}
	}

	// Verify state was applied.
	if c.GetRelation("s1", "t1") == nil {
		t.Error("table s1.t1 should exist after exec")
	}
}

func TestExecStopOnFirstError(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE TABLE t1 (id int);
		CREATE TABLE t1 (id int);
		CREATE TABLE t2 (id int);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results (stopped on error), got %d", len(results))
	}
	if results[0].Error != nil {
		t.Errorf("stmt 0: unexpected error: %v", results[0].Error)
	}
	if results[1].Error == nil {
		t.Error("stmt 1: expected error for duplicate table")
	}
	// t2 should not exist since execution stopped.
	if c.GetRelation("", "t2") != nil {
		t.Error("t2 should not exist — execution should have stopped")
	}
}

func TestExecContinueOnError(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE TABLE t1 (id int);
		CREATE TABLE t1 (id int);
		CREATE TABLE t2 (id int);
	`, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Error != nil {
		t.Errorf("stmt 0: unexpected error: %v", results[0].Error)
	}
	if results[1].Error == nil {
		t.Error("stmt 1: expected error for duplicate table")
	}
	if results[2].Error != nil {
		t.Errorf("stmt 2: unexpected error: %v", results[2].Error)
	}
	// t2 should exist since we continued.
	if c.GetRelation("", "t2") == nil {
		t.Error("t2 should exist — ContinueOnError was true")
	}
}

func TestExecCloneIsolation(t *testing.T) {
	base := New()

	_, err := base.Exec("CREATE TABLE existing (id int);", nil)
	if err != nil {
		t.Fatal(err)
	}

	trial := base.Clone()
	results, err := trial.Exec("CREATE TABLE new_table (x int);", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Error != nil {
		t.Fatalf("exec failed: %v", results)
	}

	// Trial should have both tables.
	if trial.GetRelation("", "existing") == nil {
		t.Error("trial should have 'existing'")
	}
	if trial.GetRelation("", "new_table") == nil {
		t.Error("trial should have 'new_table'")
	}

	// Base should only have the original.
	if base.GetRelation("", "existing") == nil {
		t.Error("base should have 'existing'")
	}
	if base.GetRelation("", "new_table") != nil {
		t.Error("base should NOT have 'new_table'")
	}
}

func TestExecDMLSkipped(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE TABLE t1 (id int);
		INSERT INTO t1 VALUES (1);
		SELECT * FROM t1;
		UPDATE t1 SET id = 2;
		DELETE FROM t1;
		CREATE TABLE t2 (id int);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 6 {
		t.Fatalf("expected 6 results, got %d", len(results))
	}

	// DDL statements should not be skipped.
	if results[0].Skipped {
		t.Error("CREATE TABLE t1 should not be skipped")
	}
	if results[5].Skipped {
		t.Error("CREATE TABLE t2 should not be skipped")
	}

	// DML statements should be skipped.
	for _, i := range []int{1, 2, 3, 4} {
		if !results[i].Skipped {
			t.Errorf("stmt %d (%s) should be skipped", i, results[i].SQL)
		}
		if results[i].Error != nil {
			t.Errorf("stmt %d: skipped should have no error, got: %v", i, results[i].Error)
		}
	}

	// Tables should exist.
	if c.GetRelation("", "t1") == nil {
		t.Error("t1 should exist")
	}
	if c.GetRelation("", "t2") == nil {
		t.Error("t2 should exist")
	}
}

func TestExecNoOpUtility(t *testing.T) {
	c := New()

	// Utility statements that are no-ops (SET, BEGIN/COMMIT) should not be skipped,
	// they are valid utility statements handled by ProcessUtility.
	results, err := c.Exec(`
		SET client_encoding = 'UTF8';
		BEGIN;
		CREATE TABLE t1 (id int);
		COMMIT;
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("stmt %d: unexpected error: %v", i, r.Error)
		}
		if r.Skipped {
			t.Errorf("stmt %d: should not be skipped (utility statements are processed)", i)
		}
	}

	if c.GetRelation("", "t1") == nil {
		t.Error("t1 should exist")
	}
}

func TestExecSetSearchPath(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE SCHEMA myschema;
		SET search_path = myschema, public;
		CREATE TABLE t1 (id int);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("stmt %d: unexpected error: %v", i, r.Error)
		}
	}

	// t1 should be in myschema (first in search path).
	if c.GetRelation("myschema", "t1") == nil {
		t.Error("t1 should exist in myschema")
	}
}

func TestExecLineNumbers(t *testing.T) {
	c := New()

	sql := "CREATE TABLE t1 (id int);\nCREATE TABLE t2 (id int);\nCREATE TABLE t3 (id int);"
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	expectedLines := []int{1, 2, 3}
	for i, r := range results {
		if r.Line != expectedLines[i] {
			t.Errorf("stmt %d: Line=%d, want %d", i, r.Line, expectedLines[i])
		}
	}
}

func TestExecSQLText(t *testing.T) {
	c := New()

	results, err := c.Exec("CREATE TABLE t1 (id int); CREATE TABLE t2 (id int);", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].SQL != "CREATE TABLE t1 (id int)" {
		t.Errorf("stmt 0 SQL: got %q", results[0].SQL)
	}
	if results[1].SQL != "CREATE TABLE t2 (id int)" {
		t.Errorf("stmt 1 SQL: got %q", results[1].SQL)
	}
}

func TestExecWarnings(t *testing.T) {
	c := New()

	results, err := c.Exec(`
		CREATE TABLE t1 (id int);
		CREATE TABLE IF NOT EXISTS t1 (id int);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// First: no error, no warnings.
	if results[0].Error != nil {
		t.Errorf("stmt 0: unexpected error: %v", results[0].Error)
	}

	// Second: no error, but should have a warning.
	if results[1].Error != nil {
		t.Errorf("stmt 1: unexpected error: %v", results[1].Error)
	}
	if len(results[1].Warnings) != 1 {
		t.Fatalf("stmt 1: expected 1 warning, got %d", len(results[1].Warnings))
	}
	if results[1].Warnings[0].Code != CodeWarningSkip {
		t.Errorf("warning code=%q, want %q", results[1].Warnings[0].Code, CodeWarningSkip)
	}
}

func TestExecEmptySQL(t *testing.T) {
	c := New()

	results, err := c.Exec("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty SQL, got %d", len(results))
	}

	results, err = c.Exec("   ", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for whitespace SQL, got %d", len(results))
	}
}

func TestExecParseError(t *testing.T) {
	c := New()

	_, err := c.Exec("CREATE TABLE (;", nil)
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestExecSearchPathNonExistentSchema(t *testing.T) {
	c := New()

	// SET search_path to a non-existent schema should not error (PG behavior).
	results, err := c.Exec(`
		SET search_path = nonexistent, public;
		CREATE TABLE t1 (id int);
	`, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// SET should succeed.
	if results[0].Error != nil {
		t.Errorf("SET should not error: %v", results[0].Error)
	}
	// CREATE TABLE should land in public (first existing schema in path).
	if results[1].Error != nil {
		t.Errorf("CREATE TABLE should succeed: %v", results[1].Error)
	}
	if c.GetRelation("public", "t1") == nil {
		t.Error("t1 should exist in public (nonexistent schema skipped)")
	}
}
