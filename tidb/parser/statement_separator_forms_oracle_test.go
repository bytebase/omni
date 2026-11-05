package parser

import "testing"

// TestStatementSeparatorFormsMatchTiDB holds separatorFormCases and createTableQueryFormCases
// to the live TiDB v8.5.5 server (container test). The verdict comes from
// PREPARE, which parses the statement without running it.
func TestStatementSeparatorFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, cases := range [][]separatorFormCase{separatorFormCases, createTableQueryFormCases} {
		for _, c := range cases {
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
}
