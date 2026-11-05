package parser

import "testing"

// TestSetStatementFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// setStatementCases (container test). The verdict comes from PREPARE, which parses the statement
// without running it. A 1064 does not depend on the table these forms name, so the test holds no
// schema: a statement PREPARE can no longer resolve still answers something other than 1064.
func TestSetStatementFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range setStatementCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, c.accept, c.sql)
			}
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}
