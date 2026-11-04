package parser

import "testing"

// TestAlterTableClauseFormsOracle asserts omni agrees with a live mysql:8.0 on
// alterTableClauseCases (container test). The verdict comes from PREPARE,
// which parses the statement without running it: 1064 is the parser's
// rejection, and any other server error means the text parsed and failed
// later.
//
// The expectations in alterTableClauseCases are what mysql:8.0 (8.0.46)
// answered on 2026-11-04, so an image upgrade that changes a verdict fails
// here instead of passing silently.
func TestAlterTableClauseFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range alterTableClauseCases {
		t.Run(c.sql, func(t *testing.T) {
			if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
				t.Fatalf("MySQL accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
