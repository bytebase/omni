package parser

import "testing"

// TestBindingFormsMatchTiDB holds bindingCases to the engine itself
// (container test). The verdict comes from PREPARE, which parses the
// statement without running it: 1064 is the parser's rejection, and any
// other server error means the text parsed and failed later.
func TestBindingFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, c := range bindingCases {
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

// TestBindingKeywordsStayIdentifiersOnTiDB measures the same name positions
// TestBindingKeywordsStayIdentifiers parses against v8.5.5, so the claim
// that these five words are unreserved is the server's and not this
// parser's (container test).
func TestBindingKeywordsStayIdentifiersOnTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, word := range bindingKeywords {
		for _, form := range bindingKeywordForms {
			sql := fillKeyword(form, word)
			t.Run(sql, func(t *testing.T) {
				accepts, err := tc.prepareSyntax(sql)
				if err != nil {
					t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
				}
				if !accepts {
					t.Errorf("TiDB rejects %q, so %q is not unreserved there", sql, word)
				}
			})
		}
	}
}
