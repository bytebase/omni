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

// TestExecResultLineMethodOpt pins that a method_opt statement reports its own
// line: it is a statement the parser reaches (stmt: method_opt) and the catalog
// runs nothing for, so its result must not fall to the default arm's line 1.
func TestExecResultLineMethodOpt(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test;", nil); err != nil {
		t.Fatalf("Exec error: %v", err)
	}
	results, err := c.Exec("FOR COLUMNS;\nCREATE TABLE t1 (id INT);", nil)
	if err != nil {
		t.Fatalf("Exec error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Line != 1 || results[1].Line != 2 {
		t.Errorf("lines = %d, %d, want 1, 2", results[0].Line, results[1].Line)
	}
}
