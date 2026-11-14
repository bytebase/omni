package deparse_test

import "testing"

// TestUndropAndDataMetricFunctionDropRoundTrip round-trips the UNDROP forms
// undrop.go reads and DROP FUNCTION with a data metric function's signature.
func TestUndropAndDataMetricFunctionDropRoundTrip(t *testing.T) {
	for _, sql := range []string{
		"UNDROP ACCOUNT my_account",
		"UNDROP AGENT db.s.agent1",
		"UNDROP EXTERNAL VOLUME vol1",
		"UNDROP NOTEBOOK db.s.nb",
		"UNDROP SNAPSHOT example_snapshot",
		"UNDROP SNAPSHOT example_snapshot RENAME TO restored",
		"UNDROP STAGE s.stg",
		"UNDROP STREAMLIT app",
		"UNDROP TYPE db.s.t",
		"UNDROP VARIABLE v",
		"DROP FUNCTION governance.dmfs.count_positive_numbers(TABLE(NUMBER, NUMBER))",
		"DROP FUNCTION IF EXISTS f(TABLE(INT), TABLE(VARCHAR))",
	} {
		assertRoundTrip(t, sql)
	}
}
