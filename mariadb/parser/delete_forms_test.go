package parser

import "testing"

// TestDeleteSingleTableForms: a single-table DELETE's alias, index hints and
// PARTITION are delete_single_table's and refused before USING, and DELETE
// HISTORY's TIMESTAMP 'string' is a DATETIME.
func TestDeleteSingleTableForms(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"DELETE FROM t USE INDEX(i) WHERE id = 1", true},
		{"DELETE FROM t = a WHERE a.id = 1", true},
		{"DELETE FROM t USE INDEX(i) USING t WHERE t.id = 1", false},
		{"DELETE FROM t a USING t AS a WHERE a.id = 1", false},
		{"DELETE FROM t USING t WHERE t.id = 1", true},
		{"DELETE HISTORY FROM t BEFORE SYSTEM_TIME TIMESTAMP '2020-01-01 12:00:00'", true},
		{"DELETE HISTORY FROM t BEFORE SYSTEM_TIME TIMESTAMP '2020-01-01T12:00:00'", true},
		{"DELETE HISTORY FROM t BEFORE SYSTEM_TIME TIMESTAMP '20200101120000'", true},
		{"DELETE HISTORY FROM t BEFORE SYSTEM_TIME TIMESTAMP '2020-01-01 838:00:00'", false},
		{"DELETE HISTORY FROM t BEFORE SYSTEM_TIME TIMESTAMP 's'", false},
	} {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) error = %v, want accept %v", tc.sql, err, tc.accept)
		}
	}
}
