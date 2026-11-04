package parser

import "testing"

// TestLoadDataFieldFormsOracle asserts omni agrees with this package's
// oracle on loadDataFieldCases (container test). That oracle is mysql:8.0,
// which is what oceanbase/parser has always measured against; the OceanBase
// verdicts themselves were read from oceanbase/oceanbase-ce:4.3.5-lts on
// 2026-11-04 and agree with it on every case here.
func TestLoadDataFieldFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range loadDataFieldCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.checkSyntax(t, c.sql); got != c.accept {
				t.Fatalf("the oracle accepts=%v, this file records %v for %q", got, c.accept, c.sql)
			}
			_, err := Parse(c.sql)
			if omniAccepts := err == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, the oracle accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}
