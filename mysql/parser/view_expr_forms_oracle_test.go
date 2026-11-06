package parser

import "testing"

// TestViewExprFormsOracle asserts omni agrees with a live MySQL 8.0 on
// viewExprFormCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it.
func TestViewExprFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, tc := range viewExprFormCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MySQL 8.0 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL 8.0 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
