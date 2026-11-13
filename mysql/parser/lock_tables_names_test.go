package parser

import "testing"

// TestLockTablesNamesOnce: LOCK TABLES names a table once (1066), comparing
// names as written, since whether t and T are one table is the server's
// lower_case_table_names.
func TestLockTablesNamesOnce(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"LOCK TABLES t READ, t READ", false},
		{"LOCK TABLES t AS a READ, u AS a WRITE", false},
		{"LOCK TABLES t READ, t AS a READ", true},
		{"LOCK TABLES t AS a READ, t AS A READ", true},
		{"LOCK TABLES t READ, T READ", true},
	} {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) error = %v, want accept %v", tc.sql, err, tc.accept)
		}
	}
}
