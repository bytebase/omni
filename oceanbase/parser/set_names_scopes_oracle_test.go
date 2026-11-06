package parser

import (
	"context"
	"testing"
)

// TestSetNamesScopesFormsOceanBaseOracle asserts omni agrees with a live OceanBase on
// setNamesScopesCases. startParserOracle boots MySQL in this package, so a real verdict needs
// the server itself: oceanBaseOracle runs on the package's own OceanBase (see
// testcontainer.OceanBase):
//
//	go test ./oceanbase/parser/ -run 'SetNamesScopes|SetVariableName' -count=1
func TestSetNamesScopesFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range setNamesScopesCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := prepareAcceptsOceanBase(t, ctx, db, c.sql); got != c.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", got, c.accept, c.sql)
			}
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestSetVariableNameSweepOceanBaseOracle tries every keyword in every name position.
func TestSetVariableNameSweepOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	checkSetVariableNameSweep(t, func(sql string) bool {
		return prepareAcceptsOceanBase(t, ctx, db, sql)
	})
}
