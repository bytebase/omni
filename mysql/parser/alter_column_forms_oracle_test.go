package parser

import "testing"

// TestAlterColumnFormsOracle asserts omni agrees with a live mysql:8.0 on
// alterColumnFormCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it: 1064 is the parser's rejection, and
// any other server error means the text parsed and failed later.
func TestAlterColumnFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range alterColumnFormCases {
		t.Run(c.sql, func(t *testing.T) {
			if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
				t.Fatalf("MySQL accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
