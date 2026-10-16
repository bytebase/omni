package catalog

import "testing"

func TestExecResultLine(t *testing.T) {
	sql := "CREATE TABLE t1 (id INT);\nCREATE TABLE t2 (id INT);\nCREATE TABLE t3 (id INT);"

	c := New()
	// Set up a database first.
	c.Exec("CREATE DATABASE test; USE test;", nil) //nolint:errcheck

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	wantLines := []int{1, 2, 3}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}

// TestExecResultLineSplitKeyword pins a statement's line to its first
// keyword when the object keyword follows on the next line.
func TestExecResultLineSplitKeyword(t *testing.T) {
	sql := "CREATE\nTABLE t1 (id INT);\nALTER\nTABLE t1 ADD COLUMN c INT;\nDROP\nTABLE t1;"

	c := New()
	c.Exec("CREATE DATABASE test; USE test;", nil) //nolint:errcheck

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	wantLines := []int{1, 3, 5}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}

func TestExecResultLineWithDelimiter(t *testing.T) {
	sql := "DELIMITER ;;\nCREATE TABLE t1 (id INT);;\nDELIMITER ;\nCREATE TABLE t2 (id INT);"

	c := New()
	c.Exec("CREATE DATABASE test; USE test;", nil) //nolint:errcheck

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	// First CREATE TABLE is on line 2, second is on line 4.
	wantLines := []int{2, 4}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}

// TestExecResultLineSequence pins the line of sequence DDL, which the
// catalog does not model: Exec accepts it without error and reports where
// each statement starts.
func TestExecResultLineSequence(t *testing.T) {
	sql := "CREATE TABLE t1 (id INT);\nCREATE SEQUENCE s1 START WITH 10;\n\nALTER SEQUENCE s1 RESTART;\nDROP SEQUENCE IF EXISTS s1, s2;"

	c := New()
	c.Exec("CREATE DATABASE test; USE test;", nil) //nolint:errcheck

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec error: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("got %d results, want 4", len(results))
	}
	wantLines := []int{1, 2, 4, 5}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("result[%d].Error = %v, want nil", i, r.Error)
		}
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}
