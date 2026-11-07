package catalog

import "testing"

// --- Procedure errors ---

func TestWalkThrough_2_6_CreateProcedureDuplicate(t *testing.T) {
	c := wtSetup(t)
	assertRejectedDDL(t, c, "CREATE PROCEDURE myproc() BEGIN END", 8108)
	results, err := c.Exec("CREATE PROCEDURE myproc() BEGIN END", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, 8108)
}

func TestWalkThrough_2_6_DropProcedureUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE no_such_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, 8108)
}

func TestWalkThrough_2_6_DropProcedureIfExistsUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE IF EXISTS no_such_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, 8108)
}
