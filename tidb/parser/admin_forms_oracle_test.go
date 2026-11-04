package parser

import "testing"

// TestAdminFormsMatchTiDB holds admin_forms_test.go's verdicts to the pinned
// server. Every case there was measured against TiDB v8.5.5 through PREPARE,
// as the corpus is, so a verdict in that file is a recording and this test is
// what keeps it one: an image upgrade that changes an answer fails here rather
// than leaving a stale expectation behind.
func TestAdminFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, c := range adminFormCases {
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

// TestAdminWordsStayNamesOnTiDB holds the other half of the keyword decision
// to the server: ADMIN's subcommand words are unreserved keywords in TiDB and
// ordinary names outside ADMIN, so every statement in adminWordsAsNames parses
// on v8.5.5 too. A word this parser started treating as a keyword would still
// pass TestAdminWordsStayNames if the server had changed its mind; this test
// is what says the server has not.
func TestAdminWordsStayNamesOnTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, sql := range adminWordsAsNames {
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

// TestAdminIsNotARoutineStatementOnTiDB holds adminRoutineBodyCases to the
// server: each must be a 1064 there, not a statement that parses and fails
// later.
func TestAdminIsNotARoutineStatementOnTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, sql := range adminRoutineBodyCases {
		t.Run(sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
			}
			if accepts {
				t.Errorf("TiDB accepts %q, which this package records as a rejection", sql)
			}
		})
	}
}
