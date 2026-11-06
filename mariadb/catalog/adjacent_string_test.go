package catalog

import "testing"

// TestAnalyzeAdjacentStringNames: MariaDB names a run of adjacent strings after
// the whole value, introducer or not (mariadb:11.8.8, 2026-11-06: SELECT 'a' 'b'
// is column ab, SELECT ” 'b' is column b, SELECT N'a' 'b' is column ab).
func TestAnalyzeAdjacentStringNames(t *testing.T) {
	c := New()
	for _, tc := range []struct{ sql, want string }{
		{"SELECT 'a' 'b'", "ab"},
		{"SELECT '' 'b'", "b"},
		{"SELECT 'a' 'b' 'c'", "abc"},
		{"SELECT N'a' 'b'", "ab"},
		{"SELECT _utf8mb4'a' 'b'", "ab"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			q, err := c.AnalyzeSelectStmt(parseSelect(t, tc.sql))
			if err != nil {
				t.Fatalf("AnalyzeSelectStmt(%q): %v", tc.sql, err)
			}
			if got := q.TargetList[0].ResName; got != tc.want {
				t.Errorf("%q column = %q, want %q", tc.sql, got, tc.want)
			}
		})
	}
}
