package catalog

import (
	"github.com/bytebase/omni/mysql/parser"
	"reflect"
	"strings"
	"testing"
)

// View definitions MySQL 8.0 writes for unaliased select items and for NATURAL
// and USING joins. Each want is the body of MySQL's SHOW CREATE VIEW, which
// TestViewNamesAndJoinColumnsOracle checks against a live server.
//
// An unaliased item is named after its source text, as MySQL names it: case,
// spacing and comments as written (a01, a02, a05, a07), every clause of a
// subquery (a01, a03, a04, a14), the literal forms that name themselves
// (a06), and Name_exp_N past 64 characters (a08, a09). A view column is named
// exactly as its select item (a13).
//
// A NATURAL or USING join over a join qualifies each column it joins on by the
// table its input takes it from (b01-b05, b11), a RIGHT join takes it from its
// right input (b06-b10, b16), a derived table is qualified by its alias
// (b12-b14), and SELECT * lists the joined columns once, first (b01, b17,
// b21, b22). A NATURAL join of inputs that share no column still lists a
// RIGHT join's right input first, and as an outer join it is written with
// on(true) (b23-b25).

var viewNameSchema = []string{
	`CREATE TABLE n1 (a INT, b INT)`,
	`CREATE TABLE j1 (k1 INT, a INT)`,
	`CREATE TABLE j2 (k1 INT, k2 INT, b INT)`,
	`CREATE TABLE j3 (k2 INT, k3 INT, c INT)`,
	`CREATE TABLE j4 (k3 INT, d INT)`,
	`CREATE TABLE j5 (a INT, k1 INT, b INT, k2 INT)`,
	`CREATE TABLE p (k INT, p1 INT)`,
	`CREATE TABLE q (q1 INT, q2 INT)`,
	`CREATE TABLE r (r1 INT, k INT)`,
	`CREATE TABLE s (s1 INT, s2 INT)`,
	`CREATE TABLE u (k INT, u1 INT)`,
}

var viewNameCases = []struct{ name, sql, want string }{
	{`a01`, `SELECT (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2) FROM n1`,
		"select (with `c` as (select 1 AS `x`) select `c`.`x` from `c` union select 2) AS `(WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)` from `n1`"},
	{`a02`, `select (with c as (select 1 as x) select x from   c union select 2) from n1`,
		"select (with `c` as (select 1 AS `x`) select `c`.`x` from `c` union select 2) AS `(with c as (select 1 as x) select x from   c union select 2)` from `n1`"},
	{`a03`, `SELECT (SELECT 1 UNION SELECT 2) FROM n1`,
		"select (select 1 union select 2) AS `(SELECT 1 UNION SELECT 2)` from `n1`"},
	{`a04`, `SELECT (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1), (SELECT DISTINCT a FROM n1 LIMIT 1) FROM n1`,
		"select (select `n1`.`a` from `n1` group by `n1`.`a` having (`n1`.`a` > 0) order by `n1`.`a` limit 1) AS `(SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1)`,(select distinct `n1`.`a` from `n1` limit 1) AS `(SELECT DISTINCT a FROM n1 LIMIT 1)` from `n1`"},
	{`a05`, `SELECT a+1, a  +  1, count(*), cast(a as char) FROM n1 GROUP BY a`,
		"select (`n1`.`a` + 1) AS `a+1`,(`n1`.`a` + 1) AS `a  +  1`,count(0) AS `count(*)`,cast(`n1`.`a` as char charset utf8mb4) AS `cast(a as char)` from `n1` group by `n1`.`a`"},
	{`a06`, `SELECT (a), ((b)), (1), 0004, - 2, x'41', B'10', true, date '2020-01-01', ('e'), (null) FROM n1`,
		"select `n1`.`a` AS `a`,`n1`.`b` AS `b`,1 AS `1`,4 AS `0004`,-(2) AS `- 2`,0x41 AS `x'41'`,0x02 AS `B'10'`,true AS `true`,DATE'2020-01-01' AS `date '2020-01-01'`,'e' AS `e`,NULL AS `NULL` from `n1`"},
	{`a07`, `SELECT a /* c */ + 1, a/*!50100 +2*/ FROM n1`,
		"select (`n1`.`a` + 1) AS `a /* c */ + 1`,(`n1`.`a` + 2) AS `a +2` from `n1`"},
	{`a08`, `SELECT concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') FROM n1`,
		"select concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') AS `concat('éééééééééééééééééééééééééééééééééééééééééééééééééé')` from `n1`"},
	{`a09`, `SELECT concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') FROM n1`,
		"select concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') AS `Name_exp_1` from `n1`"},
	{`a10`, `SELECT * FROM (SELECT a+1, b FROM n1) d`,
		"select `d`.`a+1` AS `a+1`,`d`.`b` AS `b` from (select (`n1`.`a` + 1) AS `a+1`,`n1`.`b` AS `b` from `n1`) `d`"},
	{`a11`, `SELECT (select   1), (SELECT a FROM n1 LIMIT 1) + 0 FROM n1`,
		"select (select 1) AS `(select   1)`,((select `n1`.`a` from `n1` limit 1) + 0) AS `(SELECT a FROM n1 LIMIT 1) + 0` from `n1`"},
	{`a12`, `SELECT a IN (SELECT a FROM n1), a in (1,2) FROM n1`,
		"select `n1`.`a` in (select `n1`.`a` from `n1`) AS `a IN (SELECT a FROM n1)`,(`n1`.`a` in (1,2)) AS `a in (1,2)` from `n1`"},
	{`a13`, "SELECT a + 1 AS `x + y`, count(*) AS `count(z)` FROM n1",
		"select (`n1`.`a` + 1) AS `x + y`,count(0) AS `count(z)` from `n1`"},
	{`a14`, `SELECT EXISTS (SELECT 1 UNION SELECT 2), a IN (SELECT 1 UNION SELECT 2) FROM n1`,
		"select exists(select 1 union select 2) AS `EXISTS (SELECT 1 UNION SELECT 2)`,`n1`.`a` in (select 1 union select 2) AS `a IN (SELECT 1 UNION SELECT 2)` from `n1`"},
	{`a15`, "SELECT * FROM (SELECT `a`+1 FROM n1) d",
		"select `d`.```a``+1` AS ```a``+1` from (select (`n1`.`a` + 1) AS ```a``+1` from `n1`) `d`"},
	{`a16`, "SELECT * FROM (SELECT a AS `x``y` FROM n1) d",
		"select `d`.`x``y` AS `x``y` from (select `n1`.`a` AS `x``y` from `n1`) `d`"},
	{`b01`, `SELECT * FROM (j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)`,
		"select `j2`.`k2` AS `k2`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c` from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))"},
	{`b02`, `SELECT * FROM j3 JOIN (j1 JOIN j2 USING (k1)) USING (k2)`,
		"select `j3`.`k2` AS `k2`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b` from (`j3` join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) on((`j3`.`k2` = `j2`.`k2`)))"},
	{`b03`, `SELECT j1.a, k2, c FROM (j1 JOIN j2 ON j1.k1 = j2.k1) JOIN j3 USING (k2)`,
		"select `j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j3`.`c` AS `c` from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))"},
	{`b04`, `SELECT * FROM (j1 NATURAL JOIN j2) NATURAL JOIN j3`,
		"select `j2`.`k2` AS `k2`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c` from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))"},
	{`b05`, `SELECT * FROM j3 NATURAL JOIN (j1 NATURAL JOIN j2)`,
		"select `j3`.`k2` AS `k2`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b` from (`j3` join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) on((`j3`.`k2` = `j2`.`k2`)))"},
	{`b06`, `SELECT * FROM j1 RIGHT JOIN j2 USING (k1)`,
		"select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b`,`j1`.`a` AS `a` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))"},
	{`b07`, `SELECT k1 FROM j1 RIGHT JOIN j2 USING (k1)`,
		"select `j2`.`k1` AS `k1` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))"},
	{`b08`, `SELECT * FROM (j1 RIGHT JOIN j2 USING (k1)) JOIN j3 USING (k2)`,
		"select `j2`.`k2` AS `k2`,`j2`.`k1` AS `k1`,`j2`.`b` AS `b`,`j1`.`a` AS `a`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c` from ((`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))"},
	{`b09`, `SELECT * FROM (j1 JOIN j2 USING (k1)) RIGHT JOIN j3 USING (k2)`,
		"select `j3`.`k2` AS `k2`,`j3`.`k3` AS `k3`,`j3`.`c` AS `c`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b` from (`j3` left join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) on((`j3`.`k2` = `j2`.`k2`)))"},
	{`b10`, `SELECT * FROM j1 NATURAL RIGHT JOIN j2`,
		"select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b`,`j1`.`a` AS `a` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))"},
	{`b11`, `SELECT * FROM ((j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)) JOIN j4 USING (k3)`,
		"select `j3`.`k3` AS `k3`,`j2`.`k2` AS `k2`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`b` AS `b`,`j3`.`c` AS `c`,`j4`.`d` AS `d` from (((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`))) join `j4` on((`j3`.`k3` = `j4`.`k3`)))"},
	{`b12`, `SELECT * FROM (SELECT k1, a FROM j1) d JOIN j2 USING (k1)`,
		"select `d`.`k1` AS `k1`,`d`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from ((select `j1`.`k1` AS `k1`,`j1`.`a` AS `a` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))"},
	{`b13`, `SELECT * FROM j2 JOIN (SELECT k2, c FROM j3) d USING (k2)`,
		"select `j2`.`k2` AS `k2`,`j2`.`k1` AS `k1`,`j2`.`b` AS `b`,`d`.`c` AS `c` from (`j2` join (select `j3`.`k2` AS `k2`,`j3`.`c` AS `c` from `j3`) `d` on((`j2`.`k2` = `d`.`k2`)))"},
	{`b14`, `SELECT * FROM (SELECT k1, a FROM j1) d NATURAL JOIN j2`,
		"select `d`.`k1` AS `k1`,`d`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from ((select `j1`.`k1` AS `k1`,`j1`.`a` AS `a` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))"},
	{`b15`, `SELECT * FROM j5 JOIN (SELECT k1, k2 FROM j2) d USING (k2, k1)`,
		"select `j5`.`k1` AS `k1`,`j5`.`k2` AS `k2`,`j5`.`a` AS `a`,`j5`.`b` AS `b` from (`j5` join (select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2` from `j2`) `d` on(((`j5`.`k1` = `d`.`k1`) and (`j5`.`k2` = `d`.`k2`))))"},
	{`b16`, `SELECT * FROM j1 RIGHT JOIN j2 USING (k1) WHERE k1 > 0 ORDER BY k1`,
		"select `j2`.`k1` AS `k1`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b`,`j1`.`a` AS `a` from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`))) where (`j2`.`k1` > 0) order by `j2`.`k1`"},
	{`b17`, `SELECT * FROM u JOIN ((q RIGHT JOIN p ON 1=1) JOIN s ON 1=1) USING (k)`,
		"select `u`.`k` AS `k`,`u`.`u1` AS `u1`,`p`.`p1` AS `p1`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`s`.`s1` AS `s1`,`s`.`s2` AS `s2` from (`u` join ((`p` left join `q` on((1 = 1))) join `s` on((1 = 1))) on((`u`.`k` = `p`.`k`)))"},
	{`b18`, `SELECT * FROM (p JOIN r USING (k)) RIGHT JOIN q ON 1=1`,
		"select `p`.`k` AS `k`,`p`.`p1` AS `p1`,`r`.`r1` AS `r1`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2` from (`q` left join (`p` join `r` on((`p`.`k` = `r`.`k`))) on((1 = 1)))"},
	{`b19`, `SELECT * FROM j1 JOIN j2 USING (K1)`,
		"select `j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`)))"},
	{`b20`, `SELECT * FROM j4, j1 JOIN j2 USING (k1)`,
		"select `j4`.`k3` AS `k3`,`j4`.`d` AS `d`,`j1`.`k1` AS `k1`,`j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from (`j4` join (`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))))"},
	{`b21`, `SELECT * FROM j4 JOIN j1 ON j4.d = j1.a JOIN j2 USING (k1)`,
		"select `j1`.`k1` AS `k1`,`j4`.`k3` AS `k3`,`j4`.`d` AS `d`,`j1`.`a` AS `a`,`j2`.`k2` AS `k2`,`j2`.`b` AS `b` from ((`j4` join `j1` on((`j4`.`d` = `j1`.`a`))) join `j2` on((`j1`.`k1` = `j2`.`k1`)))"},
	{`b22`, `SELECT * FROM u JOIN (q RIGHT JOIN (s RIGHT JOIN p ON 1=1) ON 1=1) USING (k)`,
		"select `u`.`k` AS `k`,`u`.`u1` AS `u1`,`p`.`p1` AS `p1`,`s`.`s1` AS `s1`,`s`.`s2` AS `s2`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2` from (`u` join ((`p` left join `s` on((1 = 1))) left join `q` on((1 = 1))) on((`u`.`k` = `p`.`k`)))"},
	{`b23`, `SELECT * FROM p NATURAL RIGHT JOIN q`,
		"select `q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`k` AS `k`,`p`.`p1` AS `p1` from (`q` left join `p` on(true))"},
	{`b24`, `SELECT * FROM u NATURAL JOIN (p NATURAL RIGHT JOIN q)`,
		"select `u`.`k` AS `k`,`u`.`u1` AS `u1`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`p1` AS `p1` from (`u` join (`q` left join `p` on(true)) on((`u`.`k` = `p`.`k`)))"},
	{`b25`, `SELECT * FROM (q NATURAL LEFT JOIN p) NATURAL JOIN u`,
		"select `p`.`k` AS `k`,`q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`p1` AS `p1`,`u`.`u1` AS `u1` from ((`q` left join `p` on(true)) join `u` on((`p`.`k` = `u`.`k`)))"},
}

// viewNameCatalog creates viewNameSchema and every case's view in a fresh
// catalog.
func viewNameCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	script := "CREATE DATABASE test; USE test; " + strings.Join(viewNameSchema, "; ") + ";"
	for _, tc := range viewNameCases {
		script += "\nCREATE VIEW " + tc.name + " AS " + tc.sql + ";"
	}
	results, err := c.Exec(script, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("statement %d: %v", r.Index, r.Error)
		}
	}
	return c
}

func TestViewNamesAndJoinColumns(t *testing.T) {
	c := viewNameCatalog(t)
	// A name can hold a backtick, from a quoted identifier or the source text
	// of an unaliased item, and is written with it doubled, so the definition
	// parses back.
	for _, name := range []string{"a15", "a16"} {
		def := c.GetDatabase("test").Views[name].Definition
		if _, err := parser.Parse(def); err != nil {
			t.Errorf("%s: definition does not parse: %v\n%s", name, err, def)
		}
	}
	for _, tc := range viewNameCases {
		got := extractSelectBody(c.ShowCreateView("test", tc.name))
		if got != tc.want {
			t.Errorf("%s: %s\n got: %s\nwant: %s", tc.name, tc.sql, got, tc.want)
		}
	}
	// A view column is named exactly as its select item: an explicit alias as
	// written, an unaliased item after its text.
	for name, want := range map[string][]string{
		"a05": {"a+1", "a  +  1", "count(*)", "cast(a as char)"},
		"a13": {"x + y", "count(z)"},
		"b06": {"k1", "k2", "b", "a"},
	} {
		if got := c.GetDatabase("test").Views[name].Columns; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: columns %q, want %q", name, got, want)
		}
	}
}

// TestViewNamesAndJoinColumnsOracle creates every case's view on MySQL 8.0
// and checks its SHOW CREATE VIEW body against the pinned want, and the view's
// columns in information_schema against the catalog's.
func TestViewNamesAndJoinColumnsOracle(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	for _, stmt := range viewNameSchema {
		if err := ctr.execSQLDirect(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	c := viewNameCatalog(t)
	for _, tc := range viewNameCases {
		if err := ctr.execSQLDirect("CREATE VIEW " + tc.name + " AS " + tc.sql); err != nil {
			t.Errorf("%s: MySQL rejects %s: %v", tc.name, tc.sql, err)
			continue
		}
		show, err := ctr.showCreateView(tc.name)
		if err != nil {
			t.Fatal(err)
		}
		if got := extractSelectBody(show); got != tc.want {
			t.Errorf("%s: MySQL writes\n%s\nwant\n%s", tc.name, got, tc.want)
		}
		cols, err := queryStrings(ctr, "SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = '"+tc.name+"' ORDER BY ORDINAL_POSITION")
		if err != nil {
			t.Fatal(err)
		}
		if got := c.GetDatabase("test").Views[tc.name].Columns; !reflect.DeepEqual(got, cols) {
			t.Errorf("%s: catalog columns %q, MySQL columns %q", tc.name, got, cols)
		}
	}
}
