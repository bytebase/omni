package parser

import "testing"

// TestNumberTokenFormsMatchTiDB holds number_token_forms_test.go's verdicts to the
// pinned server. Every verdict there was measured on TiDB v8.5.5 through PREPARE, as
// the corpus is, and this test is what keeps each a recording.
func TestNumberTokenFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range numberTokenStatements() {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if accepts != c.accept {
				t.Errorf("TiDB accepts=%v, this file records %v for %q", accepts, c.accept, c.sql)
			}
		})
	}
}
