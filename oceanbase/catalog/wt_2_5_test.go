package catalog

import "testing"

// Section 2.5: View Errors

func TestWalkThrough_2_5_CreateViewOnExistingName(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	wtExec(t, c, "CREATE VIEW v AS SELECT id FROM t")
	// CREATE VIEW on existing view name without OR REPLACE → error.
	// MySQL treats views and tables in the same namespace: ErrDupTable (1050).
	results, err := c.Exec("CREATE VIEW v AS SELECT id FROM t", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupTable)
}

func TestWalkThrough_2_5_CreateOrReplaceViewOnExisting(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT)")
	wtExec(t, c, "CREATE VIEW v AS SELECT id FROM t")
	// CREATE OR REPLACE VIEW on existing → no error, view is replaced.
	results, err := c.Exec("CREATE OR REPLACE VIEW v AS SELECT id FROM t", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

func TestWalkThrough_2_5_DropViewUnknown(t *testing.T) {
	c := wtSetup(t)
	// DROP VIEW on non-existent view → ErrUnknownTable (1051).
	results, err := c.Exec("DROP VIEW unknown_view", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrUnknownTable)
}

func TestWalkThrough_2_5_DropViewIfExistsUnknown(t *testing.T) {
	c := wtSetup(t)
	// DROP VIEW IF EXISTS on non-existent view → no error.
	results, err := c.Exec("DROP VIEW IF EXISTS unknown_view", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}
