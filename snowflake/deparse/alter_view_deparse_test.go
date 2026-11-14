package deparse_test

import "testing"

// ALTER VIEW and ALTER MATERIALIZED VIEW forms added for the docs' command
// reference round-trip through the deparser.

func TestDeparse_AlterViewDocumentedForms(t *testing.T) {
	assertRoundTrip(t, "ALTER VIEW v SET CHANGE_TRACKING = TRUE")
	assertRoundTrip(t, "ALTER VIEW v SET CHANGE_TRACKING = FALSE")
	assertRoundTrip(t, "ALTER VIEW v SET CONTACT steward = alice, support = bob")
	assertRoundTrip(t, "ALTER VIEW v UNSET CONTACT steward")
	assertRoundTrip(t, "ALTER VIEW v SET TAG t1 = 'x', t2 = 'y'")
	assertRoundTrip(t, "ALTER VIEW v UNSET TAG t1, t2")
	assertRoundTrip(t, "ALTER VIEW v SET DATA_METRIC_SCHEDULE = '5 MINUTE'")
	assertRoundTrip(t, "ALTER VIEW v UNSET DATA_METRIC_SCHEDULE")
	assertRoundTrip(t, "ALTER VIEW v ADD DATA METRIC FUNCTION m ON (a) EXPECTATION e (value = 0)")
	assertRoundTrip(t, "ALTER VIEW v MODIFY DATA METRIC FUNCTION m ON (a) FILTER ()")
	assertRoundTrip(t, "ALTER VIEW v SET JOIN POLICY p FORCE")
	assertRoundTrip(t, "ALTER VIEW v ALTER c SET PROJECTION POLICY p FORCE")
	assertRoundTrip(t, "ALTER VIEW v ALTER c SET MASKING POLICY m FORCE")
	assertRoundTrip(t, "ALTER VIEW v MODIFY c UNSET PROJECTION POLICY, d SET TAG t = 'x'")
	assertRoundTrip(t, "ALTER VIEW v ADD COLUMN IF NOT EXISTS c INT WITH MASKING POLICY m USING (c, d) PROJECTION POLICY p TAG (t = 'x')")
}

func TestDeparse_AlterMaterializedViewDocumentedForms(t *testing.T) {
	assertRoundTrip(t, "ALTER MATERIALIZED VIEW mv SET SECURE CONTACT steward = alice COMMENT = 'c'")
	assertRoundTrip(t, "ALTER MATERIALIZED VIEW mv SET SECURE COMMENT = 'c'")
	assertRoundTrip(t, "ALTER MATERIALIZED VIEW mv SET CONTACT a = b, c = d")
	assertRoundTrip(t, "ALTER MATERIALIZED VIEW mv UNSET CONTACT a")
	assertRoundTrip(t, "ALTER MATERIALIZED VIEW mv SET DATA_METRIC_SCHEDULE = 'TRIGGER_ON_CHANGES'")
	assertRoundTrip(t, "ALTER MATERIALIZED VIEW mv UNSET DATA_METRIC_SCHEDULE")
}
