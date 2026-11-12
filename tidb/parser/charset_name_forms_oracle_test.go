package parser

import "testing"

// TestCharsetNameFormsTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// charsetNameCases (container test).
func TestCharsetNameFormsTiDBOracle(t *testing.T) {
	holdFormsToTiDB(t, charsetNameCases)
}
