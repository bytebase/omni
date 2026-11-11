package parser

import "testing"

// TestNumberTokenFormsOracle holds number_token_forms_test.go's verdicts to the pinned
// server. Every verdict there was measured on mariadb:11.8.8 through PREPARE, so
// nothing it names has to exist and nothing runs, and this test is what keeps each a
// recording.
func TestNumberTokenFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range numberTokenStatements() {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Errorf("MariaDB 11.8 accepts=%v, this file records %v for %q", got, c.accept, c.sql)
			}
		})
	}
}
