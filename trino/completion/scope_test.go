package completion

import "testing"

// TestLexerScope_Aliases pins how the token fallback reads aliases: plain
// identifiers and non-reserved keywords are aliases, a clause keyword with
// an operand (LIMIT 1, FETCH FIRST) is not, and a reserved word never is.
func TestLexerScope_Aliases(t *testing.T) {
	for _, c := range []struct{ sql, table, alias string }{
		{"SELECT x FROM customer c WHERE", "customer", "c"},
		{"SELECT x FROM customer AS c WHERE", "customer", "c"},
		{"SELECT x FROM customer comment WHERE", "customer", "comment"},
		{"SELECT x FROM customer comment", "customer", "comment"},
		{"SELECT x FROM customer WHERE y =", "customer", ""},
		{"SELECT x FROM customer LIMIT 1 WHERE", "customer", ""},
		{"SELECT x FROM customer FETCH FIRST 1 ROWS ONLY WHERE", "customer", ""},
		{"SELECT x FROM customer JOIN orders ON", "customer", ""},
	} {
		span := lexerScope(c.sql)
		if span == nil || len(span.AccessTables) == 0 {
			t.Errorf("%q: no tables recovered", c.sql)
			continue
		}
		got := span.AccessTables[0]
		if got.Table != c.table || got.Alias != c.alias {
			t.Errorf("%q: table/alias = %q/%q, want %q/%q", c.sql, got.Table, got.Alias, c.table, c.alias)
		}
	}
}
