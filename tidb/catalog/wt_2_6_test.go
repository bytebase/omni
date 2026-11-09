package catalog

import "testing"

// --- Stored routines ---

// TiDB parses CREATE PROCEDURE and DROP PROCEDURE but refuses to run a stored
// routine, rejecting both at execution (8108). The catalog keeps a routine
// model for the Define* loader path (see define_test.go) but matches TiDB at
// Exec. The dup/no-such-routine paths MySQL would reach are therefore
// unreachable through Exec on TiDB.

func TestWalkThrough_2_6_CreateProcedureRejected(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("CREATE PROCEDURE myproc() BEGIN END", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrStoredRoutineUnsupported)
}

func TestWalkThrough_2_6_DropProcedureRejected(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE no_such_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrStoredRoutineUnsupported)
}

func TestWalkThrough_2_6_DropProcedureIfExistsRejected(t *testing.T) {
	// TiDB rejects the statement type before it checks existence, so even the
	// IF EXISTS form of an absent routine answers 8108.
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE IF EXISTS no_such_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrStoredRoutineUnsupported)
}
