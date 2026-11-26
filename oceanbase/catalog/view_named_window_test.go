package catalog

import (
	"strings"
	"testing"
)

// TestViewNamedWindowsInlined pins how a view's named windows are stored: each
// written into the OVER that names it, the WINDOW clause dropped. OceanBase
// 4.3.5.6 reads these back (measured 2026-11-26) as `rank() over( ORDER BY a
// ASC )`, `rank() over( PARTITION BY a)` and `sum(a) over( PARTITION BY a ORDER
// BY d ...)`, with no window clause.
func TestViewNamedWindowsInlined(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{"SELECT a FROM a WINDOW w AS (ORDER BY a)", "select `a`.`a` AS `a` from `a`"},
		{"SELECT RANK() OVER w r FROM a WINDOW w AS (ORDER BY a)", "rank() OVER (ORDER BY `a`.`a` )"},
		{"SELECT RANK() OVER (w) r FROM a WINDOW w AS (PARTITION BY a)", "rank() OVER (PARTITION BY `a`.`a` )"},
		{"SELECT SUM(a) OVER (w ORDER BY d) r FROM a WINDOW w AS (PARTITION BY a)", "OVER (PARTITION BY `a`.`a` ORDER BY `a`.`d` )"},
	} {
		c := New()
		if _, err := c.Exec("CREATE DATABASE p; USE p; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT)", nil); err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec("CREATE VIEW _v AS "+tc.body, nil)
		if err != nil || results[0].Error != nil {
			t.Fatalf("%s: %v %v", tc.body, err, results)
		}
		got := c.ShowCreateView("p", "_v")
		if !strings.Contains(got, tc.want) || strings.Contains(strings.ToLower(got), " window ") {
			t.Errorf("%s:\n got %s\nwant %s and no window clause", tc.body, got, tc.want)
		}
	}
}
