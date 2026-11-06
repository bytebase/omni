package parser

import "testing"

// TestAdjacentStringFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// adjacentStringCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it.
func TestAdjacentStringFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range adjacentStringCases {
		t.Run(c.sql, func(t *testing.T) {
			if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
				t.Fatalf("MariaDB accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
