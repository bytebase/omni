package parser

import "testing"

// TestODBCTableReference: `{ident table_ref}` is an element of a FROM list
// (esc_table_ref), not of a parenthesized nested list (table_ref).
func TestODBCTableReference(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT * FROM { OJ t }", true},
		{"SELECT * FROM t, { OJ u }", true},
		{"SELECT * FROM ({ OJ t })", false},
		{"SELECT * FROM (t, { OJ u })", false},
	} {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) error = %v, want accept %v", tc.sql, err, tc.accept)
		}
	}
}
