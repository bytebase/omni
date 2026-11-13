package parser

import "testing"

// TestCollateConvertWeightFormsMatchTiDB holds
// collate_convert_weight_forms_test.go's verdicts to the pinned server. Every
// case there was measured against TiDB v8.5.5 through PREPARE, as the corpus
// is, so an image upgrade that changes an answer fails here rather than
// leaving a stale expectation behind.
func TestCollateConvertWeightFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range collateConvertWeightCases {
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
