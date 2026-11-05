package parser

import "testing"

// TestSetNamesScopesFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// setNamesScopesCases (container test). The verdict comes from PREPARE, which parses the
// statement without running it.
func TestSetNamesScopesFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range setNamesScopesCases {
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

// TestSetVariableNameSweepOracle tries every keyword in every name position (container test).
func TestSetVariableNameSweepOracle(t *testing.T) {
	o := startMariaDB(t)
	checkSetVariableNameSweep(t, func(sql string) bool {
		return o.prepareAccepts(t, sql)
	})
}
