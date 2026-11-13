package parser

import "testing"

// TestStatisticsFormsMatchTiDB holds statistics_forms_test.go's verdicts to
// the pinned server. Every case there was measured against TiDB v8.5.5
// through PREPARE, as the corpus is, so a verdict in that file is a recording
// and this test is what keeps it one: an image upgrade that changes an answer
// fails here rather than leaving a stale expectation behind.
func TestStatisticsFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range statisticsFormCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if accepts != c.accept {
				t.Errorf("TiDB %ss %q, this file records %s",
					verdictName(accepts), c.sql, verdictName(c.accept))
			}
		})
	}
}

// TestStatisticsWordsStayNamesOnTiDB holds the other half of the keyword
// decision to the server: BUCKETS, TOPN, CMSKETCH, SAMPLES, SAMPLERATE,
// STATS, STATISTICS, CARDINALITY, DEPENDENCY, CORRELATION, INCREMENTAL,
// PREDICATE, WIDTH and DEPTH are keywords in TiDB's lexer and ordinary names
// outside this family, so every statement in statisticsWordsAsNames parses on
// v8.5.5 too. A word this parser started treating as a keyword would still
// pass TestStatisticsWordsStayNames if the server had changed its mind; this
// test is what says the server has not.
func TestStatisticsWordsStayNamesOnTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, sql := range statisticsWordsAsNames {
		t.Run(sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
			}
			if !accepts {
				t.Errorf("TiDB rejects %q, which this package records as a statement of ordinary names", sql)
			}
		})
	}
}
