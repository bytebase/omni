package parser

import "testing"

// TestODBCTableReference: `{ OJ table_reference }` is a table reference, not a
// factor, and holds no bare escape.
func TestODBCTableReference(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT * FROM { OJ t LEFT JOIN u ON t.a = u.a }", true},
		{"DELETE t FROM { OJ t }", true},
		{"SELECT * FROM t NATURAL JOIN { OJ u }", false},
		{"SELECT * FROM { OJ { OJ t } }", false},
		{"SELECT * FROM { t }", false},
	} {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) error = %v, want accept %v", tc.sql, err, tc.accept)
		}
	}
}
