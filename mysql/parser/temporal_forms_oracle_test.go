package parser

import "testing"

// TestTemporalFormsMatchEngine holds every row of temporalFormCases to a live
// mysql:8.0 (container test), so the recorded verdicts are the engine's and not a
// reading of the grammar. PREPARE parses the statement and stops, so no row runs.
func TestTemporalFormsMatchEngine(t *testing.T) {
	o := startParserOracle(t)
	if _, err := o.db.ExecContext(o.ctx, temporalFormsSeed); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, c := range temporalFormCases {
		if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
			t.Errorf("mysql:8.0 accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}

// temporalFormsSeed gives the rows columns to read, so that no verdict hangs on an
// unknown name.
const temporalFormsSeed = "CREATE TABLE t (a DATETIME, b INT, c INT, `now` INT, `curdate` INT, `curtime` INT, " +
	"`date_add` INT, `date_sub` INT, `sysdate` INT, `adddate` INT, `subdate` INT)"
