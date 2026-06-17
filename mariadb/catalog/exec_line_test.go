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
