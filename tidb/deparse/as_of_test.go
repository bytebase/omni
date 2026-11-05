package deparse

import "testing"

// TestAsOfDeparse checks the clause survives on the way out. A stale read that loses
// it still parses and reads the current data, which no accept-or-reject check can see.
// Each result also reads back to the same text.
func TestAsOfDeparse(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"SELECT 1 FROM t AS OF TIMESTAMP '2020-01-01'", "`t` as of timestamp '2020-01-01'"},
		{"SELECT 1 FROM t AS OF TIMESTAMP now()", "`t` as of timestamp now()"},
		{"SELECT 1 FROM t AS OF TIMESTAMP now() - interval 1 second", "`t` as of timestamp (now() - interval 1 second)"},
		{"SELECT 1 FROM t AS OF TIMESTAMP tidb_bounded_staleness(now(), now())", "`t` as of timestamp tidb_bounded_staleness(now(),now())"},
		{"SELECT 1 FROM t AS a AS OF TIMESTAMP 1", "`t` `a` as of timestamp 1"},
		{"SELECT 1 FROM t a AS OF TIMESTAMP 1 TABLESAMPLE (50)", "`t` `a` as of timestamp 1 tablesample (50)"},
		{"SELECT 1 FROM t AS `of` AS OF TIMESTAMP 1", "`t` `of` as of timestamp 1"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			want := "select 1 AS `1` from " + tc.want
			got := deparseSampleSelect(t, tc.sql)
			if got != want {
				t.Errorf("got  %s\nwant %s", got, want)
			}
			if again := deparseSampleSelect(t, got); again != got {
				t.Errorf("reads back as %s, want %s", again, got)
			}
		})
	}

	t.Run("each side of a join", func(t *testing.T) {
		got := deparseSampleSelect(t, "SELECT 1 AS x FROM t AS OF TIMESTAMP 1, u b AS OF TIMESTAMP 2")
		want := "select 1 AS `x` from (`t` as of timestamp 1 join `u` `b` as of timestamp 2)"
		if got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})
}
