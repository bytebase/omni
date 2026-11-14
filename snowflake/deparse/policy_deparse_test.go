package deparse_test

import "testing"

// DROP of the policies read by their docs pages (parser/policy.go)
// round-trips.
func TestDeparse_DropDocPolicy(t *testing.T) {
	for _, sql := range []string{
		"DROP AGGREGATION POLICY db.s.p",
		"DROP JOIN POLICY p",
		"DROP PROJECTION POLICY p",
		"DROP PRIVACY POLICY p",
		"DROP FEATURE POLICY p",
		"DROP PACKAGES POLICY IF EXISTS p",
		"DROP MAINTENANCE POLICY IF EXISTS p",
		"DROP INGRESS POLICY \"My Policy\"",
	} {
		assertRoundTrip(t, sql)
	}
}
