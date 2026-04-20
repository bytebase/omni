package catalog

import "testing"

// --- Procedure errors ---

func TestWalkThrough_2_6_CreateProcedureDuplicate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE PROCEDURE myproc() BEGIN END")
	results, err := c.Exec("CREATE PROCEDURE myproc() BEGIN END", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupProcedure)
}

func TestWalkThrough_2_6_CreateFunctionDuplicate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE FUNCTION myfunc() RETURNS INT RETURN 1")
	results, err := c.Exec("CREATE FUNCTION myfunc() RETURNS INT RETURN 1", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupFunction)
}

func TestWalkThrough_2_6_DropProcedureUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE no_such_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchProcedure)
}

func TestWalkThrough_2_6_DropFunctionUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP FUNCTION no_such_func", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchFunction)
}

func TestWalkThrough_2_6_DropProcedureIfExistsUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE IF EXISTS no_such_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

// --- Trigger errors ---

func TestWalkThrough_2_6_CreateTriggerDuplicate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t1 (id INT)")
	wtExec(t, c, "CREATE TRIGGER trg1 BEFORE INSERT ON t1 FOR EACH ROW SET @x = 1")
	results, err := c.Exec("CREATE TRIGGER trg1 BEFORE INSERT ON t1 FOR EACH ROW SET @x = 1", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupTrigger)
}

func TestWalkThrough_2_6_DropTriggerUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP TRIGGER no_such_trigger", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchTrigger)
}

func TestWalkThrough_2_6_DropTriggerIfExistsUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP TRIGGER IF EXISTS no_such_trigger", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

// --- Event errors ---

func TestWalkThrough_2_6_CreateEventDuplicate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE EVENT evt1 ON SCHEDULE EVERY 1 DAY DO SELECT 1")
	results, err := c.Exec("CREATE EVENT evt1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupEvent)
}

func TestWalkThrough_2_6_DropEventUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP EVENT no_such_event", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchEvent)
}

func TestWalkThrough_2_6_DropEventIfExistsUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP EVENT IF EXISTS no_such_event", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}
