package parser

import "testing"

// TestTypeBinaryFormsOceanBaseOracle asserts omni agrees with OceanBase CE 4.3.5.6 on
// typeBinaryCases. The verdict comes from PREPARE, which parses the statement without
// running it.
func TestTypeBinaryFormsOceanBaseOracle(t *testing.T) {
	holdFormsToOceanBase(t, typeBinaryCases)
}
