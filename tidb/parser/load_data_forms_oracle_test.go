package parser

import "testing"

// TestLoadDataFormsMatchTiDB holds load_data_forms_test.go's verdicts to the
// pinned server. Every case there was measured against TiDB v8.5.5 through
// PREPARE, as the corpus is, so a verdict in that file is a recording and
// this test is what keeps it one: an image upgrade that changes an answer
// fails here rather than leaving a stale expectation behind.
func TestLoadDataFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, c := range loadDataFormCases {
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
