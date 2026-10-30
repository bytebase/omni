package catalog

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// View columns and column references MariaDB 11.8 writes for unaliased select
// items and for NATURAL and USING joins. The catalog writes view definitions
// in MySQL 8.0's format, which MariaDB's differs from, so the oracle compares
// what does not depend on the format: the view's column names, and the
// qualified column references of its definition in order.
//
// An unaliased item is named after its source text, as MariaDB names it, and
// a NATURAL or USING join over a join qualifies each column it joins on by the
// table its input takes it from.

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

var viewNameCases = []struct{ name, sql string }{
	{`a01`, `SELECT (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2) FROM n1`},
	{`a02`, `select (with c as (select 1 as x) select x from   c union select 2) from n1`},
	{`a03`, `SELECT (SELECT 1 UNION SELECT 2) FROM n1`},
	{`a04`, `SELECT (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1), (SELECT DISTINCT a FROM n1 LIMIT 1) FROM n1`},
	{`a05`, `SELECT a+1, a  +  1, count(*), cast(a as char) FROM n1 GROUP BY a`},
	{`a06`, `SELECT (a), ((b)), (1), 0004, - 2, x'41', B'10', true, date '2020-01-01', ('e'), (null) FROM n1`},
	{`a07`, `SELECT a /* c */ + 1, a/*!50100 +2*/ FROM n1`},
	{`a08`, `SELECT concat('éééééééééééééééééééééééééééééééééééééééééééééééééé') FROM n1`},
	{`a09`, `SELECT concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee') FROM n1`},
	{`a10`, `SELECT * FROM (SELECT a+1, b FROM n1) d`},
	{`a11`, `SELECT (select   1), (SELECT a FROM n1 LIMIT 1) + 0 FROM n1`},
	{`a12`, `SELECT a IN (SELECT a FROM n1), a in (1,2) FROM n1`},
	{`a13`, "SELECT a + 1 AS `x + y`, count(*) AS `count(z)` FROM n1"},
	{`a14`, `SELECT EXISTS (SELECT 1 UNION SELECT 2), a IN (SELECT 1 UNION SELECT 2) FROM n1`},
	{`b01`, `SELECT * FROM (j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)`},
	{`b02`, `SELECT * FROM j3 JOIN (j1 JOIN j2 USING (k1)) USING (k2)`},
	{`b03`, `SELECT j1.a, k2, c FROM (j1 JOIN j2 ON j1.k1 = j2.k1) JOIN j3 USING (k2)`},
	{`b04`, `SELECT * FROM (j1 NATURAL JOIN j2) NATURAL JOIN j3`},
	{`b05`, `SELECT * FROM j3 NATURAL JOIN (j1 NATURAL JOIN j2)`},
	{`b06`, `SELECT * FROM j1 RIGHT JOIN j2 USING (k1)`},
	{`b07`, `SELECT k1 FROM j1 RIGHT JOIN j2 USING (k1)`},
	{`b08`, `SELECT * FROM (j1 RIGHT JOIN j2 USING (k1)) JOIN j3 USING (k2)`},
	{`b09`, `SELECT * FROM (j1 JOIN j2 USING (k1)) RIGHT JOIN j3 USING (k2)`},
	{`b10`, `SELECT * FROM j1 NATURAL RIGHT JOIN j2`},
	{`b11`, `SELECT * FROM ((j1 JOIN j2 USING (k1)) JOIN j3 USING (k2)) JOIN j4 USING (k3)`},
	{`b12`, `SELECT * FROM (SELECT k1, a FROM j1) d JOIN j2 USING (k1)`},
	{`b13`, `SELECT * FROM j2 JOIN (SELECT k2, c FROM j3) d USING (k2)`},
	{`b14`, `SELECT * FROM (SELECT k1, a FROM j1) d NATURAL JOIN j2`},
	{`b15`, `SELECT * FROM j5 JOIN (SELECT k1, k2 FROM j2) d USING (k2, k1)`},
	{`b16`, `SELECT * FROM j1 RIGHT JOIN j2 USING (k1) WHERE k1 > 0 ORDER BY k1`},
	{`b17`, `SELECT * FROM u JOIN ((q RIGHT JOIN p ON 1=1) JOIN s ON 1=1) USING (k)`},
	{`b18`, `SELECT * FROM (p JOIN r USING (k)) RIGHT JOIN q ON 1=1`},
	{`b19`, `SELECT * FROM j1 JOIN j2 USING (K1)`},
	{`b20`, `SELECT * FROM j4, j1 JOIN j2 USING (k1)`},
	{`b21`, `SELECT * FROM j4 JOIN j1 ON j4.d = j1.a JOIN j2 USING (k1)`},
	{`b22`, `SELECT * FROM u JOIN (q RIGHT JOIN (s RIGHT JOIN p ON 1=1) ON 1=1) USING (k)`},
	{`b23`, `SELECT * FROM p NATURAL RIGHT JOIN q`},
	{`b24`, `SELECT * FROM u NATURAL JOIN (p NATURAL RIGHT JOIN q)`},
	{`b25`, `SELECT * FROM (q NATURAL LEFT JOIN p) NATURAL JOIN u`},
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

// TestViewNamesAndJoinColumns pins columns TestViewNamesAndJoinColumnsOracle
// checks against MariaDB.
func TestViewNamesAndJoinColumns(t *testing.T) {
	c := viewNameCatalog(t)
	for name, want := range map[string][]string{
		"a01": {"(WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)"},
		"a02": {"(with c as (select 1 as x) select x from   c union select 2)"},
		"a06": {"a", "b", "1", "0004", "- 2", "x'41'", "B'10'", "TRUE", "date '2020-01-01'", "e", "NULL"},
		"a07": {"a /* c */ + 1", "a +2"},
		"a13": {"x + y", "count(z)"},
		"b06": {"k1", "k2", "b", "a"},
		"b17": {"k", "u1", "q1", "q2", "p1", "s1", "s2"},
	} {
		if got := c.GetDatabase("test").Views[name].Columns; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: columns %q, want %q", name, got, want)
		}
	}
	for name, want := range map[string]string{
		"b01": "from ((`j1` join `j2` on((`j1`.`k1` = `j2`.`k1`))) join `j3` on((`j2`.`k2` = `j3`.`k2`)))",
		"b06": "from (`j2` left join `j1` on((`j1`.`k1` = `j2`.`k1`)))",
		"b12": "from ((select `j1`.`k1` AS `k1`,`j1`.`a` AS `a` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))",
		"b23": "AS select `q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`k` AS `k`,`p`.`p1` AS `p1` from (`q` left join `p` on(1))",
	} {
		if got := c.ShowCreateView("test", name); !strings.HasSuffix(got, want) {
			t.Errorf("%s: %s\nwant it to end with %s", name, got, want)
		}
	}
}

var qualifiedColumnRef = regexp.MustCompile("`[^`]+`\\.`[^`]+`")

// TestViewNamesAndJoinColumnsOracle creates every case's view on MariaDB 11.8
// and compares its columns in information_schema, and the qualified column
// references of its SHOW CREATE VIEW, with the catalog's.
func TestViewNamesAndJoinColumnsOracle(t *testing.T) {
	db := startOrReplaceMariaDB(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, stmt := range viewNameSchema {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	c := viewNameCatalog(t)
	for _, tc := range viewNameCases {
		if _, err := conn.ExecContext(ctx, "CREATE VIEW "+tc.name+" AS "+tc.sql); err != nil {
			t.Errorf("%s: MariaDB rejects %s: %v", tc.name, tc.sql, err)
			continue
		}
		var name, show, cs, coll string
		if err := conn.QueryRowContext(ctx, "SHOW CREATE VIEW "+tc.name).Scan(&name, &show, &cs, &coll); err != nil {
			t.Fatal(err)
		}
		rows, err := conn.QueryContext(ctx, "SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", tc.name)
		if err != nil {
			t.Fatal(err)
		}
		var cols []string
		for rows.Next() {
			var col string
			if err := rows.Scan(&col); err != nil {
				t.Fatal(err)
			}
			cols = append(cols, col)
		}
		rows.Close()
		if got := c.GetDatabase("test").Views[tc.name].Columns; !reflect.DeepEqual(got, cols) {
			t.Errorf("%s: catalog columns %q, MariaDB columns %q", tc.name, got, cols)
		}
		mine := c.ShowCreateView("test", tc.name)
		if got, want := qualifiedColumnRef.FindAllString(mine, -1), qualifiedColumnRef.FindAllString(show, -1); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: catalog references %q, MariaDB references %q\ncatalog: %s\nMariaDB: %s", tc.name, got, want, mine, show)
		}
	}
}
