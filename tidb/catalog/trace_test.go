package catalog

import "testing"

// TRACE reaches the catalog because TRACE runs the statement it traces, and
// what it does there depends on TRACE PLAN: a plain TRACE and a
// TRACE FORMAT = '...' run the statement, so a SET lands, while TRACE PLAN
// and TRACE PLAN TARGET = '...' trace the optimizer and leave the catalog as
// it was. Both verdicts were measured against TiDB v8.5.5: on that server
// `TRACE SET foreign_key_checks = 0` leaves the session variable at 0 and
// `TRACE PLAN SET foreign_key_checks = 0` leaves it at 1.
//
// SET is the one traceable statement the catalog models. Upstream's
// TraceableStmt names no DDL at all -- `TRACE ALTER TABLE t ADD COLUMN c` is
// a 1064, which tidb/parser/trace_forms_test.go records -- so unlike EXPLAIN
// there is no traced DDL to apply, and foreign_key_checks is what makes the
// arm observable.
//
// Exec is an independent consumer of the same parse as review's walk: a
// dispatcher that fell through to its default for a TraceStmt would report
// success and apply nothing, and the FOREIGN KEY the next statement needs
// would then be refused.

// traceExec runs sql through a fresh catalog and returns its results.
func traceExec(t *testing.T, sql string) []ExecResult {
	t.Helper()
	c := setupWithDB(t)
	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	return results
}

// traceSetChecks runs the statement through a fresh catalog and reports the
// foreign-key-check state it left behind, failing on any error.
func traceSetChecks(t *testing.T, sql string) bool {
	t.Helper()
	c := setupWithDB(t)
	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	for i, r := range results {
		if r.Error != nil {
			t.Fatalf("statement %d (%s): %v", i, r.SQL, r.Error)
		}
	}
	return c.ForeignKeyChecks()
}

func TestExecAppliesTheSetATraceRuns(t *testing.T) {
	// The catalog starts with the checks on, so a TRACE that did not apply
	// its SET would leave this true.
	for _, sql := range []string{
		"TRACE SET foreign_key_checks = 0",
		"TRACE FORMAT = 'row' SET foreign_key_checks = 0",
		"TRACE FORMAT = 'json' SET foreign_key_checks = 0",
	} {
		if traceSetChecks(t, sql) {
			t.Errorf("%s left foreign_key_checks on, but TRACE runs the statement it traces", sql)
		}
	}
	// And back on again, so the arm is reading the value rather than
	// clearing a flag.
	if !traceSetChecks(t, "SET foreign_key_checks = 0;\nTRACE SET foreign_key_checks = 1") {
		t.Error("TRACE SET foreign_key_checks = 1 left the checks off")
	}
}

func TestExecLeavesTheTracePlanAlone(t *testing.T) {
	// TRACE PLAN traces the optimizer and runs nothing, so the checks stay
	// on.
	for _, sql := range []string{
		"TRACE PLAN SET foreign_key_checks = 0",
		"TRACE PLAN TARGET = 'estimation' SET foreign_key_checks = 0",
	} {
		if !traceSetChecks(t, sql) {
			t.Errorf("%s turned foreign_key_checks off, but TRACE PLAN runs nothing", sql)
		}
	}
}

// TestExecTraceSetReachesTheNextStatement is the effect the flag has rather
// than the flag itself: with the checks off the catalog takes a FOREIGN KEY
// to a table that is not there, and with them on it does not. A TRACE that
// applied nothing would fail the first half; a TRACE PLAN that applied its
// SET would fail the second.
func TestExecTraceSetReachesTheNextStatement(t *testing.T) {
	const fk = "CREATE TABLE child (a int, FOREIGN KEY (a) REFERENCES nosuch (id))"

	results := traceExec(t, "TRACE SET foreign_key_checks = 0;\n"+fk)
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("statement %d (%s): %v", i, r.SQL, r.Error)
		}
	}

	results = traceExec(t, "TRACE PLAN SET foreign_key_checks = 0;\n"+fk)
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	if results[0].Error != nil {
		t.Errorf("TRACE PLAN SET: %v", results[0].Error)
	}
	if results[1].Error == nil {
		t.Error("the FOREIGN KEY was taken, but TRACE PLAN never ran the SET that would allow it")
	}
}

// TestExecTraceReportsTheCallersLine pins which line a traced statement's
// result carries: the TRACE's own, not the offset of what it carries.
func TestExecTraceReportsTheCallersLine(t *testing.T) {
	results := traceExec(t, "CREATE TABLE t (a int);\nTRACE SET foreign_key_checks = 0;\nSELECT 1")
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	if results[1].Line != 2 {
		t.Errorf("TRACE line = %d, want 2", results[1].Line)
	}
}

// TestExecTraceOfARowChangeIsNoCatalogChange runs the traceable statements
// the catalog does not model. Each is applied rather than left to the default
// arm, so none is an error, and none touches a schema object.
func TestExecTraceOfARowChangeIsNoCatalogChange(t *testing.T) {
	for _, sql := range []string{
		"TRACE SELECT 1",
		"TRACE SELECT a FROM t",
		"TRACE INSERT INTO t VALUES (1)",
		"TRACE UPDATE t SET a = 1",
		"TRACE DELETE FROM t",
		"TRACE REPLACE INTO t VALUES (1)",
		"TRACE BEGIN",
		"TRACE COMMIT",
		"TRACE ROLLBACK",
		"TRACE SAVEPOINT sp1",
		"TRACE RELEASE SAVEPOINT sp1",
		"TRACE ANALYZE TABLE t",
		"TRACE LOAD DATA INFILE '/tmp/a.csv' INTO TABLE t",
		"TRACE PLAN SELECT 1",
		"TRACE FORMAT = 'json' SELECT 1",
		"TRACE TABLE t",
		"TRACE VALUES ROW(1)",
	} {
		for _, r := range traceExec(t, "CREATE TABLE t (a int);\n"+sql) {
			if r.Error != nil {
				t.Errorf("%s: %v", r.SQL, r.Error)
			}
		}
	}
}
