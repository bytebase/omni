package parser

import "testing"

// TestMariaDBFunctionArgumentOracle holds every row of funcArgCases to a live
// mariadb:11.8.8 (container test), so the recorded verdicts are the engine's
// and not MySQL's inherited. They differ inside these parentheses: this engine
// takes DISTINCT on JSON_ARRAYAGG and an argument list with it, takes no ALL
// there or on JSON_OBJECTAGG, and has no CAST ... ARRAY at all.
func TestMariaDBFunctionArgumentOracle(t *testing.T) {
	o := startMariaDB(t)

	for _, ddl := range []string{
		"CREATE TABLE t (c1 INT, c2 INT, b JSON, a DATETIME, profit INT)",
		"CREATE TABLE sales (profit INT)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, c := range funcArgCases {
		if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
			t.Errorf("mariadb:11.8.8 accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
