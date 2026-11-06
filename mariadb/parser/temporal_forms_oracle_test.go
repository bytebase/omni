package parser

import "testing"

// TestMariaDBTemporalFormsMatchEngine holds every row of temporalFormCases to a
// live mariadb:11.8.8 (container test), so the recorded verdicts are the engine's
// and not MySQL's inherited. They differ here: this engine takes a 0x literal and a
// literal up to 18446744073709551615 as a precision, where MySQL takes neither, and
// SYSDATE is a keyword whether or not a space follows it.
func TestMariaDBTemporalFormsMatchEngine(t *testing.T) {
	o := startMariaDB(t)
	if _, err := o.db.ExecContext(o.ctx, temporalFormsSeed); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, c := range temporalFormCases {
		if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
			t.Errorf("mariadb:11.8.8 accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}

// temporalFormsSeed gives the rows columns to read, so that no verdict hangs on an
// unknown name.
const temporalFormsSeed = "CREATE TABLE t (a DATETIME, b INT, c INT, `now` INT, `curdate` INT, `curtime` INT, " +
	"`date_add` INT, `date_sub` INT, `sysdate` INT, `adddate` INT, `subdate` INT)"
