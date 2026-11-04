package parser

import "testing"

// TestFunctionArgumentOracle holds every row of funcArgCases to a live
// mysql:8.0 (container test), so the recorded verdicts are the engine's and
// not a reading of the manual. The reject rows matter most: which quantifier
// an aggregate takes is not uniform and not shared with the forks of this
// grammar — this engine rejects STD(DISTINCT c1) and GROUP_CONCAT(ALL c1),
// which TiDB parses, and parses COUNT(ALL *), which TiDB rejects.
func TestFunctionArgumentOracle(t *testing.T) {
	o := startParserOracle(t)

	for _, ddl := range []string{
		"CREATE TABLE t (c1 INT, c2 INT, b JSON, a DATETIME, profit INT)",
		"CREATE TABLE sales (profit INT)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, c := range funcArgCases {
		if accepts := o.checkSyntax(t, c.sql); accepts != c.accept {
			t.Errorf("mysql:8.0 accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
