package catalog

import "testing"

// PLAN REPLAYER reaches the catalog for the reason EXPLAIN does: the
// statement a dump carries may be ALTER TABLE, and whether the server runs it
// depends on ANALYZE. On TiDB v8.5.5
// `PLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN c` adds the
// column, while `PLAN REPLAYER DUMP EXPLAIN ALTER TABLE t ADD COLUMN c`
// writes a dump and leaves the table alone. Both were measured on that
// server, and tidb/parser/plan_replayer_forms_test.go records which forms
// parse.
//
// Exec is an independent consumer of the same parse as review's walk: a
// dispatcher that fell through to its default for a PlanReplayerStmt would
// report success and apply nothing, and the next statement would then fail
// here on a column the server has.

// planReplayerExec runs sql through a fresh catalog and returns its results.
func planReplayerExec(t *testing.T, sql string) []ExecResult {
	t.Helper()
	c := setupWithDB(t)
	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	return results
}

func TestExecAppliesTheDDLAPlanReplayerDumpAnalyzeRuns(t *testing.T) {
	// The column the ANALYZE form adds is there for the statement that
	// follows, so dropping it succeeds.
	results := planReplayerExec(t, "CREATE TABLE t (a int);\nPLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t DROP COLUMN b")
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("statement %d (%s): %v", i, r.SQL, r.Error)
		}
	}
	// The line each result reports is the line the caller wrote, so the
	// dump's own line and not the ALTER's offset inside it.
	if results[1].Line != 2 {
		t.Errorf("PLAN REPLAYER DUMP line = %d, want 2", results[1].Line)
	}
	// The WITH STATS clause is no part of that decision.
	results = planReplayerExec(t, "CREATE TABLE t (a int);\nPLAN REPLAYER DUMP WITH STATS AS OF TIMESTAMP '2024-01-01' EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t DROP COLUMN b")
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("statement %d (%s): %v", i, r.SQL, r.Error)
		}
	}
}

func TestExecLeavesThePlainPlanReplayerDumpAlone(t *testing.T) {
	// A plain dump adds nothing, so dropping the column fails on a catalog
	// that never had it. This is the half that a wrapper applied twice, or
	// applied when the server does not, gets wrong.
	results := planReplayerExec(t, "CREATE TABLE t (a int);\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t DROP COLUMN b")
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	if results[1].Error != nil {
		t.Errorf("PLAN REPLAYER DUMP EXPLAIN ALTER TABLE: %v", results[1].Error)
	}
	if results[2].Error == nil {
		t.Error("DROP COLUMN b succeeded, but a plain dump added no column")
	}
	// ...and a plain dump of a change the catalog would refuse is no error
	// of its own, because it applies nothing.
	for _, r := range planReplayerExec(t, "CREATE TABLE t (a int);\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE t ADD COLUMN a int") {
		if r.Error != nil {
			t.Errorf("%s: %v", r.SQL, r.Error)
		}
	}
}

// TestExecPlanReplayerDDLLandsOnce is the guard against a wrapped DDL applied
// twice: the column the ANALYZE form adds is added once, so adding it again
// fails and no second copy of it exists.
func TestExecPlanReplayerDDLLandsOnce(t *testing.T) {
	results := planReplayerExec(t, "CREATE TABLE t (a int);\nPLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t ADD COLUMN b int")
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	if results[1].Error != nil {
		t.Errorf("PLAN REPLAYER DUMP EXPLAIN ANALYZE: %v", results[1].Error)
	}
	if results[2].Error == nil {
		t.Error("ADD COLUMN b succeeded a second time, so the dump did not add it")
	}
	// Dropping it once succeeds and twice does not, which is the same
	// statement read from the other side.
	results = planReplayerExec(t, "CREATE TABLE t (a int);\nPLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t DROP COLUMN b;\nALTER TABLE t DROP COLUMN b")
	if len(results) != 4 {
		t.Fatalf("results = %d, want 4", len(results))
	}
	if results[2].Error != nil {
		t.Errorf("the first DROP COLUMN b: %v", results[2].Error)
	}
	if results[3].Error == nil {
		t.Error("the second DROP COLUMN b succeeded, so the column was added twice")
	}
}

// TestExecPlanReplayerOfAQueryChangesNothing runs the forms that carry no DDL.
// Each is applied rather than left to the default arm, so none is an error,
// and none touches a schema object.
func TestExecPlanReplayerOfAQueryChangesNothing(t *testing.T) {
	for _, sql := range []string{
		"PLAN REPLAYER DUMP EXPLAIN SELECT 1",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE SELECT 1",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE UPDATE t SET a = 1",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE DELETE FROM t",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE INSERT INTO t VALUES (1)",
		"PLAN REPLAYER DUMP EXPLAIN TABLE t",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE TABLE t",
		"PLAN REPLAYER DUMP EXPLAIN SLOW QUERY",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE SLOW QUERY WHERE a > 10 ORDER BY t LIMIT 10",
		"PLAN REPLAYER DUMP EXPLAIN 'sql.txt'",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE 'sql.txt'",
		"PLAN REPLAYER LOAD '/tmp/a.zip'",
		"PLAN REPLAYER CAPTURE '123' '456'",
		"PLAN REPLAYER CAPTURE REMOVE '123' '456'",
	} {
		for _, r := range planReplayerExec(t, "CREATE TABLE t (a int);\n"+sql) {
			if r.Error != nil {
				t.Errorf("%s: %v", r.SQL, r.Error)
			}
		}
	}
}
