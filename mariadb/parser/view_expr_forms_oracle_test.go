package parser

import "testing"

// TestViewExprFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// viewExprFormCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it.
func TestViewExprFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range viewExprFormCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
