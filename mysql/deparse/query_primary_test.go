package deparse

import "testing"

// TestDeparseQueryPrimaryLocking: a TABLE or VALUES query primary carries the
// locking clause after it on its SelectStmt, so deparse has to write it. A
// writer that emitted only the primary and its ORDER BY / LIMIT rewrote
// `TABLE t FOR UPDATE` as `table `t“, which is the same query with its
// locking semantics dropped. MySQL 8.0 accepts every input here
// (query_primary_forms_oracle_test.go in mysql/parser holds that column).
func TestDeparseQueryPrimaryLocking(t *testing.T) {
	cases := []struct{ sql, want string }{
		{"TABLE t", "table `t`"},
		{"TABLE t FOR UPDATE", "table `t` for update"},
		{"TABLE t FOR SHARE", "table `t` for share"},
		{"TABLE t LOCK IN SHARE MODE", "table `t` lock in share mode"},
		{"TABLE t FOR UPDATE NOWAIT", "table `t` for update nowait"},
		{"TABLE t FOR UPDATE OF t SKIP LOCKED", "table `t` for update of `t` skip locked"},
		{"TABLE t ORDER BY a LIMIT 1 FOR UPDATE", "table `t` order by `a` limit 1 for update"},
		{"VALUES ROW(1,2)", "values row(1,2)"},
		{"VALUES ROW(1,2) FOR UPDATE", "values row(1,2) for update"},
		{"VALUES ROW(1) ORDER BY 1 LIMIT 1 FOR SHARE", "values row(1) order by 1 limit 1 for share"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := DeparseSelect(parseSelect(t, tc.sql)); got != tc.want {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.sql, got, tc.want)
			}
		})
	}
}
