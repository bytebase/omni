package parser

import "testing"

// TestTypeCharsetNameFormsOceanBaseOracle asserts omni agrees with a live OceanBase CE
// 4.3.5.6 on typeCharsetNameCases (container test).
func TestTypeCharsetNameFormsOceanBaseOracle(t *testing.T) {
	holdFormsToOceanBase(t, typeCharsetNameCases)
}
