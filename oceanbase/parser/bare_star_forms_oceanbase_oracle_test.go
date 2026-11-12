package parser

import (
	"context"
	"testing"
)

// TestBareStarFormsOceanBaseOracle asserts omni agrees with a live OceanBase CE 4.3.5.6
// on bareStarCases (container test). The verdict comes from
// PREPARE, which parses the statement without running it.
func TestBareStarFormsOceanBaseOracle(t *testing.T) {
	holdFormsToOceanBase(t, bareStarCases)
}

// holdFormsToOceanBase asserts that a live OceanBase gives each case the verdict it
// records, and that omni gives it too.
func holdFormsToOceanBase(t *testing.T, cases []parserFormCase) {
	t.Helper()
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := prepareAcceptsOceanBase(t, ctx, db, tc.sql); got != tc.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
