package parser

import "testing"

// TestJSONQueryFormsOceanBaseOracle asserts omni agrees with OceanBase CE 4.3.5.6 on
// jsonQueryCases. The verdict comes from PREPARE, which parses the statement without
// running it.
func TestJSONQueryFormsOceanBaseOracle(t *testing.T) {
	holdFormsToOceanBase(t, jsonQueryCases)
}
