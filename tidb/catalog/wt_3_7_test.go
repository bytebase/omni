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

func TestWalkThrough_3_7_CreateFunction(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `DELIMITER ;;
CREATE FUNCTION my_func(a INT, b INT) RETURNS INT
DETERMINISTIC
RETURN a + b;;
DELIMITER ;`)

	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	fn := db.Functions[toLower("my_func")]
	if fn == nil {
		t.Fatal("function my_func not found in database.Functions")
	}
	if fn.Name != "my_func" {
		t.Errorf("expected name my_func, got %s", fn.Name)
	}
	if fn.IsProcedure {
		t.Error("expected IsProcedure=false for function")
	}
	// Return type should contain "int"
	if fn.Returns == "" {
		t.Error("expected non-empty Returns for function")
	}
	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.Params[0].Name != "a" {
		t.Errorf("param 0 name: expected a, got %s", fn.Params[0].Name)
	}
	if fn.Params[1].Name != "b" {
		t.Errorf("param 1 name: expected b, got %s", fn.Params[1].Name)
	}
}

func TestWalkThrough_3_7_AlterProcedure(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `DELIMITER ;;
CREATE PROCEDURE my_proc()
BEGIN
  SELECT 1;
END;;
DELIMITER ;`)

	wtExec(t, c, "ALTER PROCEDURE my_proc COMMENT 'updated comment'")

	db := c.GetDatabase("testdb")
	proc := db.Procedures[toLower("my_proc")]
	if proc == nil {
		t.Fatal("procedure my_proc not found")
	}
	if proc.Characteristics["COMMENT"] != "updated comment" {
		t.Errorf("expected COMMENT 'updated comment', got %q", proc.Characteristics["COMMENT"])
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

func TestWalkThrough_3_7_CreateTrigger(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT, val INT)")
	wtExec(t, c, `DELIMITER ;;
CREATE TRIGGER trg_before_insert BEFORE INSERT ON t1 FOR EACH ROW
BEGIN
  SET NEW.val = NEW.val + 1;
END;;
DELIMITER ;`)

	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	trg := db.Triggers[toLower("trg_before_insert")]
	if trg == nil {
		t.Fatal("trigger trg_before_insert not found in database.Triggers")
	}
	if trg.Name != "trg_before_insert" {
		t.Errorf("expected name trg_before_insert, got %s", trg.Name)
	}
	if trg.Timing != "BEFORE" {
		t.Errorf("expected timing BEFORE, got %s", trg.Timing)
	}
	if trg.Event != "INSERT" {
		t.Errorf("expected event INSERT, got %s", trg.Event)
	}
	if trg.Table != "t1" {
		t.Errorf("expected table t1, got %s", trg.Table)
	}
}

func TestWalkThrough_3_7_DropTrigger(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	wtExec(t, c, `DELIMITER ;;
CREATE TRIGGER trg1 AFTER DELETE ON t1 FOR EACH ROW
BEGIN
  SELECT 1;
END;;
DELIMITER ;`)

	db := c.GetDatabase("testdb")
	if db.Triggers[toLower("trg1")] == nil {
		t.Fatal("trigger should exist before drop")
	}

	wtExec(t, c, "DROP TRIGGER trg1")

	if db.Triggers[toLower("trg1")] != nil {
		t.Error("trigger trg1 should be removed after DROP TRIGGER")
	}
}

func TestWalkThrough_3_7_CreateEvent(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE EVENT my_event
ON SCHEDULE EVERY 1 HOUR
ON COMPLETION PRESERVE
ENABLE
COMMENT 'hourly cleanup'
DO DELETE FROM t1 WHERE created < NOW() - INTERVAL 1 DAY`)

	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	ev := db.Events[toLower("my_event")]
	if ev == nil {
		t.Fatal("event my_event not found in database.Events")
	}
	if ev.Name != "my_event" {
		t.Errorf("expected name my_event, got %s", ev.Name)
	}
	if ev.Schedule == "" {
		t.Error("expected non-empty schedule")
	}
	// Enable should be ENABLE or default
	if ev.Enable != "" && ev.Enable != "ENABLE" {
		t.Errorf("expected Enable ENABLE or empty, got %s", ev.Enable)
	}
	if ev.OnCompletion != "PRESERVE" {
		t.Errorf("expected OnCompletion PRESERVE, got %s", ev.OnCompletion)
	}
}

func TestWalkThrough_3_7_AlterEvent(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE EVENT my_event ON SCHEDULE EVERY 1 HOUR DO SELECT 1")

	wtExec(t, c, "ALTER EVENT my_event ON SCHEDULE EVERY 2 HOUR DISABLE")

	db := c.GetDatabase("testdb")
	ev := db.Events[toLower("my_event")]
	if ev == nil {
		t.Fatal("event my_event not found")
	}
	if ev.Enable != "DISABLE" {
		t.Errorf("expected Enable DISABLE after alter, got %s", ev.Enable)
	}
	// Schedule should be updated
	if ev.Schedule == "" {
		t.Error("expected non-empty schedule after alter")
	}
}

func TestWalkThrough_3_7_AlterEventRename(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE EVENT old_event ON SCHEDULE EVERY 1 HOUR DO SELECT 1")

	wtExec(t, c, "ALTER EVENT old_event RENAME TO new_event")

	db := c.GetDatabase("testdb")
	if db.Events[toLower("old_event")] != nil {
		t.Error("old_event should not exist after rename")
	}
	ev := db.Events[toLower("new_event")]
	if ev == nil {
		t.Fatal("new_event should exist after rename")
	}
	if ev.Name != "new_event" {
		t.Errorf("expected name new_event, got %s", ev.Name)
	}
}

func TestWalkThrough_3_7_DropEvent(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE EVENT my_event ON SCHEDULE EVERY 1 HOUR DO SELECT 1")

	db := c.GetDatabase("testdb")
	if db.Events[toLower("my_event")] == nil {
		t.Fatal("event should exist before drop")
	}

	wtExec(t, c, "DROP EVENT my_event")

	if db.Events[toLower("my_event")] != nil {
		t.Error("event my_event should be removed after DROP EVENT")
	}
}
