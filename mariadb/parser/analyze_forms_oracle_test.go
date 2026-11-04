package parser

import "testing"

// TestAnalyzeFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// analyzeFormCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it.
func TestAnalyzeFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range analyzeFormCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}

// TestMariaDBRefusesMySQLHistograms records the false accept this package
// still carries, so it is a measured and named gap rather than an untested
// one. MariaDB 11.8.8 answers 1064 for both arms and omni takes both; the
// test fails if the server ever starts taking them, or if omni stops.
//
// Removing the arms moves mariadb/review's recorded expectations for "a
// histogram of MySQL's, which MariaDB refuses", so it belongs to its own
// change (container test).
func TestMariaDBRefusesMySQLHistograms(t *testing.T) {
	o := startMariaDB(t)
	for _, sql := range []string{
		"ANALYZE TABLE t UPDATE HISTOGRAM ON b",
		"ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 BUCKETS",
		"ANALYZE TABLE t UPDATE HISTOGRAM ON b USING DATA '{}'",
		"ANALYZE TABLE t DROP HISTOGRAM ON b",
	} {
		t.Run(sql, func(t *testing.T) {
			if o.prepareAccepts(t, sql) {
				t.Errorf("MariaDB 11.8.8 now accepts %q, which this package records as a rejection", sql)
			}
			if _, err := Parse(sql); err != nil {
				t.Errorf("omni now rejects %q; this package records it as the false accept left for its own change: %v", sql, err)
			}
		})
	}
}
