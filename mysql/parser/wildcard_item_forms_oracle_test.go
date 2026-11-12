package parser

import "testing"

// TestWildcardItemFormsOracle asserts omni agrees with a live MySQL 8.0 on
// wildcardItemCases, and that the server still refuses every wildcardItemGaps statement
// (container test). The verdict comes from PREPARE, which parses the statement without
// running it.
func TestWildcardItemFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, tc := range wildcardItemCases {
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
	for _, sql := range wildcardItemGaps {
		t.Run(sql, func(t *testing.T) {
			if o.prepareAccepts(t, sql) {
				t.Errorf("MySQL 8.0 accepts %q, which wildcardItemGaps records as refused", sql)
			}
		})
	}
}
