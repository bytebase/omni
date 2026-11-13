package parser

import "testing"

// TestDefinerHostFormsOceanBaseOracle asserts omni agrees with OceanBase CE 4.3.5.6 on
// definerHostCases.
func TestDefinerHostFormsOceanBaseOracle(t *testing.T) {
	runFormsOnOceanBase(t, "definer_host_forms", definerHostCases)
}
