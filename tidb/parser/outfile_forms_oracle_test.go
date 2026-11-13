package parser

import "testing"

// TestOutfileFieldsFormsMatchTiDB holds outfile_forms_test.go's verdicts to the
// pinned server. Each case is measured against TiDB v8.5.5 through PREPARE, so
// the verdicts in that file are a recording and this test keeps them one: an
// image upgrade that changes an answer fails here rather than leaving a stale
// expectation behind.
func TestOutfileFieldsFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range outfileFieldsFormCases {
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
