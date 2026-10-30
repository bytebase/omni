package catalog

import "testing"

// Section 1.3: Result Metadata — Line, SQL

func TestWalkThrough_1_3_SingleLineMultiStatement(t *testing.T) {
	// Single-line multi-statement: each Result.Line is 1
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT); CREATE TABLE t3 (id INT);"
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	for i, r := range results {
		if r.Line != 1 {
			t.Errorf("result[%d].Line = %d, want 1", i, r.Line)
		}
	}
}

func TestWalkThrough_1_3_MultiLineStatements(t *testing.T) {
	// Multi-line statements: Result.Line matches first line of each statement
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (\n  id INT\n);\nCREATE TABLE t2 (\n  id INT\n);\nCREATE TABLE t3 (\n  id INT\n);"
	// Line 1: CREATE TABLE t1 (
	// Line 2:   id INT
	// Line 3: );
	// Line 4: CREATE TABLE t2 (
	// Line 5:   id INT
	// Line 6: );
	// Line 7: CREATE TABLE t3 (
	// Line 8:   id INT
	// Line 9: );
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	wantLines := []int{1, 4, 7}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}

func TestWalkThrough_1_3_DelimiterMode(t *testing.T) {
	// DELIMITER mode: Result.Line points to correct line in original SQL
	c := wtSetup(t)
	sql := "DELIMITER ;;\nCREATE TABLE t1 (id INT);;\nDELIMITER ;\nCREATE TABLE t2 (id INT);"
	// Line 1: DELIMITER ;;
	// Line 2: CREATE TABLE t1 (id INT);;
	// Line 3: DELIMITER ;
	// Line 4: CREATE TABLE t2 (id INT);
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	wantLines := []int{2, 4}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}

func TestWalkThrough_1_3_BlankLines(t *testing.T) {
	// Statements after blank lines: Line numbers account for blank lines
	c := wtSetup(t)
	sql := "CREATE TABLE t1 (id INT);\n\n\nCREATE TABLE t2 (id INT);\n\nCREATE TABLE t3 (id INT);"
	// Line 1: CREATE TABLE t1 (id INT);
	// Line 2: (blank)
	// Line 3: (blank)
	// Line 4: CREATE TABLE t2 (id INT);
	// Line 5: (blank)
	// Line 6: CREATE TABLE t3 (id INT);
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	wantLines := []int{1, 4, 6}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
	}
}

func TestWalkThrough_1_3_DMLLineNumbers(t *testing.T) {
	// Result.Line for DML (skipped) statements is still correct
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	sql := "SELECT * FROM t1;\nINSERT INTO t1 VALUES (1);\nCREATE TABLE t2 (id INT);"
	// Line 1: SELECT * FROM t1;
	// Line 2: INSERT INTO t1 VALUES (1);
	// Line 3: CREATE TABLE t2 (id INT);
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	// DML statements should have correct line numbers even though skipped
	wantLines := []int{1, 2, 3}
	wantSkipped := []bool{true, true, false}
	for i, r := range results {
		if r.Line != wantLines[i] {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, wantLines[i])
		}
		if r.Skipped != wantSkipped[i] {
			t.Errorf("result[%d].Skipped = %v, want %v", i, r.Skipped, wantSkipped[i])
		}
	}
}
