package catalog

import "testing"

// A view's column list renames its select items, and whatever names one by the alias it
// had has to follow. An ORDER BY that names a select item is written as its position.
func TestViewColumnListFollowsAliases(t *testing.T) {
	for body, want := range map[string]string{
		"v (p) AS SELECT 1 AS a UNION SELECT 2 ORDER BY 1":              "(select 1 AS `p`) union (select 2 AS `p`) order by 1",
		"v (p) AS SELECT 1 AS a UNION SELECT 2 ORDER BY a":              "(select 1 AS `p`) union (select 2 AS `p`) order by 1",
		"v (p) AS SELECT a AS x FROM t ORDER BY x":                      "select `t`.`a` AS `p` from `t` order by 1",
		"v (p) AS SELECT a AS x FROM t ORDER BY x + 1":                  "select `t`.`a` AS `p` from `t` order by (`p` + 1)",
		"v (p) AS SELECT a + 1 AS x FROM t GROUP BY x HAVING x > 1":     "select (`t`.`a` + 1) AS `p` from `t` group by `p` having (`p` > 1)",
		"v (p) AS (SELECT a AS x FROM t ORDER BY x LIMIT 1) ORDER BY x": "(select `t`.`a` AS `p` from `t` order by 1 limit 1) order by 1",
		"v (b) AS SELECT a AS x FROM t ORDER BY b":                      "select `t`.`a` AS `b` from `t` order by `t`.`b`",
		"v AS SELECT a AS x, b FROM t ORDER BY x DESC, b, a + 1":        "select `t`.`a` AS `x`,`t`.`b` AS `b` from `t` order by 1 desc,2,(`t`.`a` + 1)",
	} {
		c := New()
		if code := omniCode(c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT); CREATE VIEW "+body); code != 0 {
			t.Errorf("%s: error %d", body, code)
			continue
		}
		v := c.GetDatabase("d").Views["v"]
		if v.Definition != want {
			t.Errorf("%s\n got: %s\nwant: %s", body, v.Definition, want)
		}
		// The body is one the catalog reads back to itself.
		again := New()
		if code := omniCode(again, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT); CREATE VIEW v AS "+v.Definition); code != 0 {
			t.Errorf("%s: its body does not load: error %d\n%s", body, code, v.Definition)
		} else if got := again.GetDatabase("d").Views["v"].Definition; got != v.Definition {
			t.Errorf("%s: its body reads back as another\n got: %s\nwant: %s", body, got, v.Definition)
		}
	}
}
