package parser

import "testing"

// TestViewExprFormsMatchEngine holds view_expr_forms_test.go's verdicts to the
// server, so an image upgrade that changes an answer fails here rather than
// leaving a stale expectation behind.
func TestViewExprFormsMatchEngine(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range viewExprFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("TiDB v8.5.5 accepts %q = %v, this file records %v", c.sql, got, c.accept)
			}
		})
	}
}
