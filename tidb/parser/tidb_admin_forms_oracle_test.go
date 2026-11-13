package parser

import (
	"errors"
	"maps"
	"slices"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestTiDBAdminFormsMatchTiDB holds tidb_admin_forms_test.go's verdicts to the
// pinned server, so a verdict in that file stays a recording: an image upgrade
// that changes an answer fails here rather than leaving a stale expectation behind.
func TestTiDBAdminFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range tidbAdminCases {
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

// TestTiDBAdminRefusalsHappenWhileParsing holds tidbAdminRefusedWhileParsing to the
// server. A code alone cannot say where it was raised, 1105 least of all, so each
// statement is prepared twice: alone, and with a syntax error after it. The same code
// both times means the server stopped before it read the tail. That is the reading
// `RECOMMEND INDEX APPLY 1` fails, which answers 1105 alone, from the executor, and
// 1064 with the tail.
func TestTiDBAdminRefusalsHappenWhileParsing(t *testing.T) {
	tc := startTiDB(t)
	prepareCode := func(t *testing.T, sql string) uint16 {
		t.Helper()
		conn, err := tc.db.Conn(tc.ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if _, err := conn.ExecContext(tc.ctx, "SET @omni_admin = ?", sql); err != nil {
			t.Fatal(err)
		}
		_, err = conn.ExecContext(tc.ctx, "PREPARE omni_admin_stmt FROM @omni_admin")
		if err == nil {
			_, _ = conn.ExecContext(tc.ctx, "DEALLOCATE PREPARE omni_admin_stmt")
			return 0
		}
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) {
			t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
		}
		return myErr.Number
	}
	const tail = "\n; SELECT FROM"
	for _, c := range tidbAdminRefusedWhileParsing {
		t.Run(c.sql, func(t *testing.T) {
			if got := prepareCode(t, c.sql); got != c.code {
				t.Errorf("TiDB answers %d for %q, this file records %d", got, c.sql, c.code)
			}
			if got := prepareCode(t, c.sql+tail); got != c.code {
				t.Errorf("with a syntax error after it TiDB answers %d for %q: the refusal did not come from the parser", got, c.sql)
			}
		})
	}
	for _, sql := range []string{"RECOMMEND INDEX APPLY 1", "RECOMMEND INDEX IGNORE 1"} {
		t.Run(sql, func(t *testing.T) {
			if got := prepareCode(t, sql+tail); got != 1064 {
				t.Errorf("with a syntax error after it TiDB answers %d for %q, so 1105 alone is not the parser's", got, sql)
			}
		})
	}
}

// TestTiDBAdminWordsStayNamesOnTiDB holds the other half of the bare-word decision
// to the server: every statement in tidbAdminWordsAsNames parses on v8.5.5 too. A word
// this parser started treating as a keyword would still pass TestTiDBAdminWordsStayNames
// if the server had changed its mind; this test is what says the server has not.
func TestTiDBAdminWordsStayNamesOnTiDB(t *testing.T) {
	tc := startTiDB(t)
	for _, sql := range tidbAdminWordsAsNames {
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

// tidbReservedWordsTheLexerLacks are words TiDB reserves that omni's lexer has no
// keyword for, so parseIdent reads them as names and `CREATE TABLE ilike (a int)` is
// an over-accept. omni reads CURRENT_ROLE() and TIDB_CURRENT_TSO() as function calls
// through that same reading, so this cannot be fixed in parseIdent alone. The name
// positions of TestKillFunctionNameSweepMatchesTiDB leave them out.
var tidbReservedWordsTheLexerLacks = []string{"current_role", "ilike", "tidb_current_tso"}

// TestKillFunctionNameSweepMatchesTiDB walks every word either lexer knows through the
// positions of this family that depend on what the word is:
//
//   - KILL word() and KILL word (1), which hold tidbLexerWords and
//     tidbBuiltinOnlyWords to the server. The function form takes only a word TiDB's
//     lexer has no token for, which omni's own keyword table cannot say, and the
//     spaced form is what separates a word only btFuncTokenMap holds.
//   - SET CONFIG tikv word = 1 and SET CONFIG tikv a-word = 1, the names a config item
//     is made of, which take TiDB's Identifier: an unreserved keyword is in and a
//     reserved one out. Every other name position in the family reads the same way.
//
// A table copied from upstream is what a release changes without anyone noticing, so
// this sweeps the words rather than trusting the copy.
func TestKillFunctionNameSweepMatchesTiDB(t *testing.T) {
	tc := startTiDB(t)

	words := map[string]bool{"foo": true, "_foo": true}
	maps.Copy(words, tidbLexerWords)
	maps.Copy(words, tidbBuiltinOnlyWords)
	for w := range keywords {
		words[w] = true
	}
	for cs := range tidbCharsets {
		words["_"+cs] = true
	}

	check := func(t *testing.T, sql string) {
		t.Helper()
		accepts, err := tc.prepareSyntax(sql)
		if err != nil {
			t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
		}
		_, perr := Parse(sql)
		if omni := perr == nil; omni != accepts {
			t.Errorf("omni %ss %q, TiDB %ss it (err=%v)", verdictName(omni), sql, verdictName(accepts), perr)
		}
	}
	for _, w := range slices.Sorted(maps.Keys(words)) {
		// Not a word: the table holds a few spellings the scanner takes apart.
		if strings.ContainsAny(w, " .") {
			continue
		}
		check(t, "KILL "+w+"()")
		check(t, "KILL "+w+" (1)")
		if !slices.Contains(tidbReservedWordsTheLexerLacks, w) {
			check(t, "SET CONFIG tikv "+w+" = 1")
			check(t, "SET CONFIG tikv a-"+w+" = 1")
		}
	}
}
