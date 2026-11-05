package parser

import "testing"

// TestSetListFormsOracle asserts omni agrees with a live MariaDB on setListCases (container test).
// The verdict comes from PREPARE, which parses the statement without running it.
func TestSetListFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range setListCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Fatalf("MariaDB accepts=%v, want %v for %q", got, c.accept, c.sql)
			}
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}
