package parser

import "testing"

// TestLoadDataFieldFormsOracle asserts omni agrees with a live mysql:8.0 on
// loadDataFieldCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it.
func TestLoadDataFieldFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range loadDataFieldCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Fatalf("MySQL 8.0 accepts=%v, this file records %v for %q", got, c.accept, c.sql)
			}
			_, err := Parse(c.sql)
			if omniAccepts := err == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL 8.0 accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}
