package catalog

import "testing"

// EXPLAIN reaches the catalog because TiDB explains ALTER TABLE, and what it
// does there depends on ANALYZE: EXPLAIN ANALYZE runs the statement it
// carries, so the change lands, while a plain EXPLAIN reports a plan and
// leaves the catalog as it was. Both verdicts are measured against TiDB
// v8.5.5 in tidb/parser/explain_forms_test.go.
//
// Exec is an independent consumer of the same parse as review's walk: a
// dispatcher that fell through to its default for an ExplainStmt would report
// success and apply nothing, and the next statement would then fail here on a
// column the server has.

// explainExec runs sql through a fresh catalog and returns the error of the
// statement at index i, or nil.
func explainExec(t *testing.T, sql string) []ExecResult {
	t.Helper()
	c := setupWithDB(t)
	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	return results
}

func TestExecAppliesTheDDLExplainAnalyzeRuns(t *testing.T) {
	// The column EXPLAIN ANALYZE adds is there for the statement that
	// follows, so dropping it succeeds.
	results := explainExec(t, "CREATE TABLE t (a int);\nEXPLAIN ANALYZE ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t DROP COLUMN b")
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("statement %d (%s): %v", i, r.SQL, r.Error)
		}
	}
	// The line each result reports is the line the caller wrote, so the
	// EXPLAIN's own line and not the ALTER's offset inside it.
	if results[1].Line != 2 {
		t.Errorf("EXPLAIN ANALYZE line = %d, want 2", results[1].Line)
	}
}

func TestExecLeavesThePlainExplainAlone(t *testing.T) {
	// A plain EXPLAIN adds nothing, so dropping the column fails on a
	// catalog that never had it.
	results := explainExec(t, "CREATE TABLE t (a int);\nEXPLAIN ALTER TABLE t ADD COLUMN b int;\nALTER TABLE t DROP COLUMN b")
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	if results[1].Error != nil {
		t.Errorf("EXPLAIN ALTER TABLE: %v", results[1].Error)
	}
	if results[2].Error == nil {
		t.Error("DROP COLUMN b succeeded, but a plain EXPLAIN added no column")
	}
	// ...and a plain EXPLAIN of a change the catalog would refuse is no
	// error of its own, because it applies nothing.
	for _, r := range explainExec(t, "CREATE TABLE t (a int);\nEXPLAIN ALTER TABLE t ADD COLUMN a int") {
		if r.Error != nil {
			t.Errorf("%s: %v", r.SQL, r.Error)
		}
	}
}

func TestExecExplainOfAQueryChangesNothing(t *testing.T) {
	for _, sql := range []string{
		"EXPLAIN SELECT 1",
		"EXPLAIN ANALYZE SELECT 1",
		"EXPLAIN ANALYZE UPDATE t SET a = 1",
		"EXPLAIN t",
		"EXPLAIN FOR CONNECTION 1",
	} {
		for _, r := range explainExec(t, "CREATE TABLE t (a int);\n"+sql) {
			if r.Error != nil {
				t.Errorf("%s: %v", r.SQL, r.Error)
			}
		}
	}
}
