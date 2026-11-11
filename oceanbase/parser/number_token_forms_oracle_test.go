package parser

import (
	"context"
	"testing"
)

// TestNumberTokenFormsOceanBaseOracle holds number_token_forms_test.go's verdicts to
// the pinned server. Every verdict there was measured on OceanBase 4.3.5.6 through
// PREPARE, so nothing it names has to exist and nothing runs, and this test is what
// keeps each a recording.
func TestNumberTokenFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range numberTokenStatements() {
		t.Run(c.sql, func(t *testing.T) {
			if got := prepareAcceptsOceanBase(t, ctx, db, c.sql); got != c.accept {
				t.Errorf("OceanBase accepts=%v, this file records %v for %q", got, c.accept, c.sql)
			}
		})
	}
}
