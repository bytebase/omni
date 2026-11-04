package parser

import "testing"

// TestTiDBFunctionArgumentOracle holds every row of funcArgCases to a live
// TiDB v8.5.5 (container test), so the recorded verdicts are the engine's and
// not a reading of the manual. The reject rows matter most: a quantifier set
// or an argument count taken from MySQL's grammar is not evidence of TiDB's,
// and the two differ inside these parentheses — TiDB takes DISTINCT on the
// statistical aggregates and ALL on GROUP_CONCAT where MySQL rejects both, and
// rejects COUNT(ALL *) where MySQL parses it.
func TestTiDBFunctionArgumentOracle(t *testing.T) {
	tc := startTiDB(t)

	tc.db.ExecContext(tc.ctx, "CREATE DATABASE IF NOT EXISTS omni_funcargs")
	tc.db.ExecContext(tc.ctx, "USE omni_funcargs")
	defer tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_funcargs")

	for _, ddl := range []string{
		"CREATE TABLE IF NOT EXISTS t (c1 INT, c2 INT, b JSON, a DATETIME, profit INT)",
		"CREATE TABLE IF NOT EXISTS sales (profit INT)",
	} {
		if err := tc.canExecute(ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, c := range funcArgCases {
		accepts := !tidbRejectedSyntax(t, tc.canExecute(c.sql))
		if accepts != c.accept {
			t.Errorf("TiDB accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
