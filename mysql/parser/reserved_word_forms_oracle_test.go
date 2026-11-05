package parser

import "testing"

// TestReservedWordFormsOracle asserts omni agrees with a live MySQL 8.0 on reservedWordCases
// (container test). The verdict comes from PREPARE, which parses the statement without running it.
func TestReservedWordFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range reservedWordCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Fatalf("MySQL 8.0 accepts=%v, want %v for %q", got, c.accept, c.sql)
			}
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL 8.0 accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestReservedWordSweepOracle holds every registered keyword to a live MySQL 8.0 (container test).
func TestReservedWordSweepOracle(t *testing.T) {
	o := startParserOracle(t)
	checkReservedWordSweep(t, func(sql string) bool {
		return o.prepareAccepts(t, sql)
	})
}
