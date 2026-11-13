package catalog

import (
	"strings"
	"testing"
)

// A view's column list renames its select items, and whatever names one by the alias it
// had has to follow. An ORDER BY that names a select item is written as its position, and
// an alias anywhere else in a SELECT stands for the item's expression, so a new name that
// is also a column of the table does not take its place.
func TestViewColumnListFollowsAliases(t *testing.T) {
	for body, want := range map[string]string{
		"v (p) AS SELECT 1 AS a UNION SELECT 2 ORDER BY 1":                                 "(select 1 AS `p`) union (select 2 AS `p`) order by 1",
		"v (p) AS SELECT 1 AS a UNION SELECT 2 ORDER BY a":                                 "(select 1 AS `p`) union (select 2 AS `p`) order by 1",
		"v (p) AS SELECT a AS x FROM t ORDER BY x":                                         "select `t`.`a` AS `p` from `t` order by 1",
		"v (p) AS SELECT a AS x FROM t ORDER BY x + 1":                                     "select `t`.`a` AS `p` from `t` order by (`t`.`a` + 1)",
		"v (p) AS SELECT a + 1 AS x FROM t GROUP BY x HAVING x > 1":                        "select (`t`.`a` + 1) AS `p` from `t` group by (`t`.`a` + 1) having ((`t`.`a` + 1) > 1)",
		"v (p) AS (SELECT a AS x FROM t ORDER BY x LIMIT 1) ORDER BY x":                    "(select `t`.`a` AS `p` from `t` order by 1 limit 1) order by 1",
		"v (b) AS SELECT a AS x FROM t ORDER BY x + 1, b":                                  "select `t`.`a` AS `b` from `t` order by (`t`.`a` + 1),`t`.`b`",
		"v (p) AS SELECT 5 AS x FROM t GROUP BY x ORDER BY x":                              "select 5 AS `p` from `t` group by 1 order by 1",
		"v (b) AS SELECT a AS x FROM t ORDER BY CONVERT(x, SIGNED), EXTRACT(DAY FROM x)":   "select `t`.`a` AS `b` from `t` order by cast(`t`.`a` as signed),extract(day from `t`.`a`)",
		"v (b) AS SELECT a AS x FROM t ORDER BY COALESCE(x, 0) + SUM(x) OVER (ORDER BY x)": "",
		"v (b) AS SELECT a AS x FROM t ORDER BY b":                                         "select `t`.`a` AS `b` from `t` order by `t`.`b`",
		"v AS SELECT a AS x, b FROM t ORDER BY x DESC, b, a + 1":                           "select `t`.`a` AS `x`,`t`.`b` AS `b` from `t` order by 1 desc,2,(`t`.`a` + 1)",
	} {
		c := New()
		if code := omniCode(c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT); CREATE VIEW "+body); code != 0 {
			t.Errorf("%s: error %d", body, code)
			continue
		}
		v := c.GetDatabase("d").Views["v"]
		if want == "" {
			// Only that every alias is gone, and the body reads back to itself.
			if strings.Contains(v.Definition, "`x`") {
				t.Errorf("%s: an alias is left: %s", body, v.Definition)
			}
		} else if v.Definition != want {
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

// A query nested in an ORDER BY is a scope of its own: the column list renames none of
// its aliases, even where one is spelled as an outer one.
func TestViewColumnListLeavesNestedQueries(t *testing.T) {
	c := New()
	body := "v (b) AS SELECT a AS x FROM t ORDER BY (WITH c AS (SELECT 1 AS x HAVING x > 0) SELECT x FROM c)"
	if code := omniCode(c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT); CREATE VIEW "+body); code != 0 {
		t.Fatalf("%s: error %d", body, code)
	}
	def := c.GetDatabase("d").Views["v"].Definition
	if strings.Contains(def, "`b` > 0") || !strings.Contains(def, "`x` > 0") {
		t.Errorf("the nested query's alias was renamed: %s", def)
	}
}
