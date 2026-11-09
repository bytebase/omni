package catalog

import "testing"

// --- 3.7 Routine, Trigger, and Event State ---

// TiDB parses CREATE/DROP PROCEDURE but refuses to run a stored routine (8108),
// so neither enters the catalog through Exec. The routine metadata model the
// catalog still keeps is exercised through the Define* loader path instead
// (see TestDefineProcedure_* in define_test.go).

func TestWalkThrough_3_7_CreateProcedureRejected(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec(`DELIMITER ;;
CREATE PROCEDURE my_proc(IN a INT, OUT b VARCHAR(100))
BEGIN
  SET b = CONCAT('hello', a);
END;;
DELIMITER ;`, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrStoredRoutineUnsupported)

	db := c.GetDatabase("testdb")
	if db != nil && db.Procedures[toLower("my_proc")] != nil {
		t.Error("procedure my_proc must not enter the catalog through Exec on TiDB")
	}
}

func TestWalkThrough_3_7_DropProcedureRejected(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP PROCEDURE my_proc", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrStoredRoutineUnsupported)
}
