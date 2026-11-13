package catalog

import "testing"

// TestQueryFirstReadsWhatTheAnalyzerCannotRefuse holds resolveQueryFirst to the forms
// review found it refused: a CTE read from a later CTE or a nested query, DUAL, and
// a derived table's column named by its expression's text. A name the analyzer
// cannot settle is not refused; one it can still is. A parenthesized query's outer
// ORDER BY or LIMIT is checked and leaves no analyzed query that would drop it.
func TestQueryFirstReadsWhatTheAnalyzerCannotRefuse(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x), d AS (SELECT x FROM c) SELECT x FROM d", 0},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT x FROM (SELECT x FROM c) s", 0},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT (SELECT x FROM c) AS y", 0},
		{"CREATE VIEW v AS SELECT 1 AS x FROM DUAL", 0},
		{"CREATE VIEW v AS SELECT d.`b + 1` FROM (SELECT b + 1 FROM a) d", 0},
		{"CREATE VIEW v AS (SELECT b FROM a) ORDER BY b LIMIT 1", 0},
		{"CREATE VIEW v AS (SELECT b FROM a) ORDER BY zz", 1054},
		{"CREATE VIEW v AS SELECT zz FROM (SELECT b FROM a) d", 1054},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT zz FROM c", 1054},
		{"CREATE VIEW v AS SELECT zz FROM DUAL", 1054},
	} {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE a (b INT)")
		if err != nil {
			t.Fatal(err)
		}
		if got := omniCode(c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
		if v := c.GetDatabase("d").Views["v"]; v != nil && v.AnalyzedQuery != nil && tc.sql == "CREATE VIEW v AS (SELECT b FROM a) ORDER BY b LIMIT 1" {
			t.Errorf("%s: kept an analyzed query without its ORDER BY and LIMIT", tc.sql)
		}
	}
}

// TestOuterAliasInSubqueryKeepsNoQuery holds the other half of that bargain: a
// subquery may name the select list of a query around it, which the engine takes
// (TestOracle_OceanBaseQueryFirst), and the view is created — but the alias's
// expression counts range tables from the outer query, so no analyzed query is
// kept rather than one whose references read the wrong table. An alias of the
// query's own select list is unaffected and keeps its query.
func TestOuterAliasInSubqueryKeepsNoQuery(t *testing.T) {
	for _, tc := range []struct {
		sql string
		ir  string
	}{
		{"CREATE VIEW v AS SELECT b AS q FROM a ORDER BY (SELECT q FROM z)", ""},
		{"CREATE VIEW v AS SELECT b AS q FROM a ORDER BY (SELECT q)", ""},
		{"CREATE VIEW v AS SELECT b AS q FROM a GROUP BY q HAVING (SELECT q) > 0", ""},
		{"CREATE VIEW v AS SELECT b AS q FROM a ORDER BY q", "select `a`.`b` as `q` from `d`.`a` order by `a`.`b`"},
	} {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE a (b INT); CREATE TABLE z (c INT)")
		if err != nil {
			t.Fatal(err)
		}
		if got := omniCode(c, tc.sql); got != 0 {
			t.Errorf("%s: refused with %d", tc.sql, got)
			continue
		}
		v := c.GetDatabase("d").Views["v"]
		if v == nil {
			t.Errorf("%s: no view", tc.sql)
			continue
		}
		got := ""
		if v.AnalyzedQuery != nil {
			got = DeparseQuery(v.AnalyzedQuery)
		}
		if got != tc.ir {
			t.Errorf("%s:\n  analyzed query %q\n  want           %q", tc.sql, got, tc.ir)
		}
	}
}
