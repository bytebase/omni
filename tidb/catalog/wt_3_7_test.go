package catalog

import "testing"

// --- 3.7 Routine, Trigger, and Event State ---

func TestWalkThrough_3_7_CreateProcedure(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `DELIMITER ;;
CREATE PROCEDURE my_proc(IN a INT, OUT b VARCHAR(100))
BEGIN
  SET b = CONCAT('hello', a);
END;;
DELIMITER ;`)

	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	proc := db.Procedures[toLower("my_proc")]
	if proc == nil {
		t.Fatal("procedure my_proc not found in database.Procedures")
	}
	if proc.Name != "my_proc" {
		t.Errorf("expected name my_proc, got %s", proc.Name)
	}
	if !proc.IsProcedure {
		t.Error("expected IsProcedure=true")
	}
	if len(proc.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(proc.Params))
	}
	// Check param a
	if proc.Params[0].Direction != "IN" {
		t.Errorf("param 0 direction: expected IN, got %s", proc.Params[0].Direction)
	}
	if proc.Params[0].Name != "a" {
		t.Errorf("param 0 name: expected a, got %s", proc.Params[0].Name)
	}
	if proc.Params[0].TypeName != "INT" {
		t.Errorf("param 0 type: expected INT, got %s", proc.Params[0].TypeName)
	}
	// Check param b
	if proc.Params[1].Direction != "OUT" {
		t.Errorf("param 1 direction: expected OUT, got %s", proc.Params[1].Direction)
	}
	if proc.Params[1].Name != "b" {
		t.Errorf("param 1 name: expected b, got %s", proc.Params[1].Name)
	}
	if proc.Params[1].TypeName != "VARCHAR(100)" {
		t.Errorf("param 1 type: expected VARCHAR(100), got %s", proc.Params[1].TypeName)
	}
}

func TestWalkThrough_3_7_DropProcedure(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `DELIMITER ;;
CREATE PROCEDURE my_proc()
BEGIN
  SELECT 1;
END;;
DELIMITER ;`)

	// Verify it exists first.
	db := c.GetDatabase("testdb")
	if db.Procedures[toLower("my_proc")] == nil {
		t.Fatal("procedure should exist before drop")
	}

	wtExec(t, c, "DROP PROCEDURE my_proc")

	if db.Procedures[toLower("my_proc")] != nil {
		t.Error("procedure my_proc should be removed after DROP PROCEDURE")
	}
}
