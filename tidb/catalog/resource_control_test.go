package catalog

import "testing"

// What the resource control statements leave behind in the catalog.
//
// Nothing, by design: a resource group, a calibration and a query-watch entry
// are none of the catalog's objects, so processUtility does not act on them.
// What they must still carry is their line, because Exec reports one
// ExecResult per statement and a caller shows the error against it. stmtLocStart
// reads Loc.Start through a type switch and answers 0 for a statement it does
// not list, so a statement missing from that switch reports line 1 wherever it
// stands, and no corpus row can see it: the corpus records whether a statement
// parses and nothing about what the catalog did with it.
func TestResourceControlExecLines(t *testing.T) {
	stmts := []string{
		"CREATE RESOURCE GROUP rg RU_PER_SEC = 1000",
		"ALTER RESOURCE GROUP rg QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION KILL)",
		"DROP RESOURCE GROUP IF EXISTS rg",
		"SET RESOURCE GROUP DEFAULT",
		"CALIBRATE RESOURCE WORKLOAD TPCC",
		"CALIBRATE RESOURCE START_TIME now() DURATION '20m'",
		"QUERY WATCH ADD RESOURCE GROUP rg SQL TEXT SIMILAR TO 'select 1'",
		"QUERY WATCH REMOVE 1",
	}

	// One statement per line, so each result's line is its index plus one.
	sql := ""
	for _, s := range stmts {
		sql += s + ";\n"
	}

	c := New()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != len(stmts) {
		t.Fatalf("got %d results, want %d", len(results), len(stmts))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("%s: error %v, want none", stmts[i], r.Error)
		}
		if r.Skipped {
			t.Errorf("%s: skipped, want executed as a utility statement", stmts[i])
		}
		if r.Line != i+1 {
			t.Errorf("%s: line %d, want %d", stmts[i], r.Line, i+1)
		}
	}

	// None of them is a catalog object, so the catalog is still empty.
	if got := len(c.Databases()); got != 0 {
		t.Errorf("the catalog holds %d databases, want none", got)
	}
}
