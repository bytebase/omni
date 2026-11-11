package parser

import "testing"

// TestNumberTokenFormsOracle holds number_token_forms_test.go's verdicts to the pinned
// server. Every verdict there was measured on mysql:8.0 through PREPARE, so nothing it
// names has to exist and nothing runs, and this test is what keeps each a recording.
func TestNumberTokenFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range numberTokenStatements() {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Errorf("MySQL 8.0 accepts=%v, this file records %v for %q", got, c.accept, c.sql)
			}
		})
	}
}
