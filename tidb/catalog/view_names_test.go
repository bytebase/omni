package catalog

import (
	"github.com/bytebase/omni/tidb/parser"
	"reflect"
	"strings"
	"testing"
)

// View columns TiDB v8.5 derives for unaliased select items and for NATURAL
// and USING joins. TiDB writes its view definitions in its own format and
// keeps USING and NATURAL as written, so the oracle compares the view's
// column names.
//
// An unaliased item is named after its source text, as TiDB names it, and
// SELECT * over a NATURAL or USING join lists the columns it joins on once,
// first.
//
// Not covered: TiDB rejects a WITH inside a scalar subquery of a view and a
// scalar subquery that may return more than one row (a01-a03 of the mysql
// cases), and it lists a name longer than 64 bytes as name_exp_N in lower
// case where the catalog writes Name_exp_N (a08, a09), a known gap.

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
	{`a04`, `SELECT (SELECT a FROM n1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1), (SELECT DISTINCT a FROM n1 LIMIT 1) FROM n1`},
	{`a05`, `SELECT a+1, a  +  1, count(*), cast(a as char) FROM n1 GROUP BY a`},
	{`a06`, `SELECT (a), ((b)), (1), 0004, - 2, x'41', B'10', true, date '2020-01-01', ('e'), (null) FROM n1`},
	{`a07`, `SELECT a /* c */ + 1, a/*!50100 +2*/ FROM n1`},
	{`a10`, `SELECT * FROM (SELECT a+1, b FROM n1) d`},
	{`a11`, `SELECT (select   1), (SELECT a FROM n1 LIMIT 1) + 0 FROM n1`},
	{`a12`, `SELECT a IN (SELECT a FROM n1), a in (1,2) FROM n1`},
	{`a13`, "SELECT a + 1 AS `x + y`, count(*) AS `count(z)` FROM n1"},
	{`a14`, `SELECT EXISTS (SELECT 1 UNION SELECT 2), a IN (SELECT 1 UNION SELECT 2) FROM n1`},
	{`a15`, "SELECT * FROM (SELECT `a`+1 FROM n1) d"},
	{`a16`, "SELECT * FROM (SELECT a AS `x``y` FROM n1) d"},
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
// checks against TiDB.
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
	for name, want := range map[string][]string{
		"a06": {"(a)", "((b))", "1", "0004", "- 2", "x'41'", "B'10'", "TRUE", "date '2020-01-01'", "e", "NULL"},
		"a07": {"a /* c  + 1", "a +2"},
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
		"b06": "from (`j2` left join `j1` on((`j2`.`k1` = `j1`.`k1`)))",
		"b12": "from ((select `j1`.`k1` AS `k1`,`j1`.`a` AS `a` from `j1`) `d` join `j2` on((`d`.`k1` = `j2`.`k1`)))",
		"b23": "AS select `q`.`q1` AS `q1`,`q`.`q2` AS `q2`,`p`.`k` AS `k`,`p`.`p1` AS `p1` from (`q` left join `p` on(true))",
	} {
		if got := c.ShowCreateView("test", name); !strings.HasSuffix(got, want) {
			t.Errorf("%s: %s\nwant it to end with %s", name, got, want)
		}
	}
}

// TestViewNamesAndJoinColumnsOracle creates every case's view on TiDB and
// compares its columns in information_schema with the catalog's.
func TestViewNamesAndJoinColumnsOracle(t *testing.T) {
	tc := startTiDBForCatalog(t)
	mustExecTiDB(t, tc, "DROP DATABASE IF EXISTS omni_view_names")
	mustExecTiDB(t, tc, "CREATE DATABASE omni_view_names")
	mustExecTiDB(t, tc, "USE omni_view_names")
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_view_names") })
	for _, stmt := range viewNameSchema {
		mustExecTiDB(t, tc, stmt)
	}
	c := viewNameCatalog(t)
	for _, v := range viewNameCases {
		if _, err := tc.db.ExecContext(tc.ctx, "CREATE VIEW "+v.name+" AS "+v.sql); err != nil {
			t.Errorf("%s: TiDB rejects %s: %v", v.name, v.sql, err)
			continue
		}
		rows, err := tc.db.QueryContext(tc.ctx, "SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'omni_view_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", v.name)
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
		if got := c.GetDatabase("test").Views[v.name].Columns; !reflect.DeepEqual(got, cols) {
			t.Errorf("%s: catalog columns %q, TiDB columns %q", v.name, got, cols)
		}
	}
}
