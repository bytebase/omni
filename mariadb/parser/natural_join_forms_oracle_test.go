package parser

import "testing"

// naturalJoinFormCases are the joins after NATURAL: table_ref NATURAL inner_join
// table_factor, where inner_join is JOIN, INNER JOIN or STRAIGHT_JOIN, and the LEFT and
// RIGHT [OUTER] JOIN arms (MariaDB 11.8 sql/sql_yacc.yy).
var naturalJoinFormCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT * FROM a NATURAL JOIN b", true},
	{"SELECT * FROM a NATURAL INNER JOIN b", true},
	{"SELECT * FROM a NATURAL STRAIGHT_JOIN b", true},
	{"SELECT * FROM a NATURAL LEFT OUTER JOIN b", true},
	{"DELETE a FROM a NATURAL INNER JOIN b", true},
	{"SELECT * FROM a NATURAL CROSS JOIN b", false},
	{"SELECT * FROM a NATURAL INNER b", false},
	{"SELECT * FROM a NATURAL INNER STRAIGHT_JOIN b", false},
	{"SELECT * FROM a NATURAL INNER JOIN b ON a.x = b.x", false},
}

// TestNaturalJoinFormsOracle holds naturalJoinFormCases to a live MariaDB 11.8.8, and
// omni to each verdict.
func TestNaturalJoinFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range naturalJoinFormCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
