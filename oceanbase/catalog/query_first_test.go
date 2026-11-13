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
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) (SELECT zz FROM c)", 1054},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c) ORDER BY zz", 1054},
		{"CREATE VIEW v AS SELECT zz FROM DUAL", 1054},
		{"CREATE VIEW v AS (SELECT b+1 FROM a) ORDER BY `b+1`", 0},
		{"CREATE VIEW v AS (SELECT count(*) FROM a) ORDER BY `count(*)`", 0},
		{"CREATE VIEW v AS SELECT d.`b+1` FROM (SELECT * FROM (SELECT b+1 FROM a) s) d", 0},
		{"CREATE VIEW v AS WITH c AS (SELECT b+1 FROM a), e AS (SELECT * FROM c) SELECT `b+1` FROM e", 0},
		{"CREATE VIEW v AS SELECT b+1 FROM a UNION ALL SELECT b+2 FROM a ORDER BY `b+1`", 0},
		{"CREATE VIEW v AS SELECT b FROM a UNION ALL SELECT b FROM a ORDER BY zz", 1054},
		{"CREATE VIEW v AS SELECT d.`0004` FROM (SELECT 0004) d", 0},
		{"CREATE VIEW v AS SELECT d.`(1)` FROM (SELECT (1)) d", 0},
		{"CREATE VIEW v AS SELECT d.`zz` FROM (SELECT 4) d", 1054},
		{"CREATE VIEW v AS WITH RECURSIVE r AS (SELECT 1+1 UNION ALL SELECT 3 FROM r WHERE 0) SELECT `1+1` FROM r", 0},
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

// TestAnalyzedQueryKeptOnlyWhenItIsFaithful holds the other half of that bargain:
// a query the engine takes (TestOracle_OceanBaseQueryFirst) creates its view, and
// an analyzed query is kept only where it says what the statement said. A subquery
// naming an outer select alias keeps none, since the alias's expression counts
// range tables from the outer query; so does a parenthesized query, whose own WITH,
// ORDER BY and LIMIT a Query has no layer for. An alias of the query's own select
// list, and a plain WITH, are unaffected and keep theirs.
func TestAnalyzedQueryKeptOnlyWhenItIsFaithful(t *testing.T) {
	for _, tc := range []struct {
		sql string
		ir  string
	}{
		{"CREATE VIEW v AS SELECT b AS q FROM a ORDER BY (SELECT q FROM z)", ""},
		{"CREATE VIEW v AS SELECT b AS q FROM a ORDER BY (SELECT q)", ""},
		{"CREATE VIEW v AS SELECT b AS q FROM a GROUP BY q HAVING (SELECT q) > 0", ""},
		{"CREATE VIEW v AS SELECT b AS q FROM a ORDER BY q", "select `a`.`b` as `q` from `d`.`a` order by `a`.`b`"},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c) ORDER BY x", ""},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c)", ""},
		{"CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT x FROM c", "with `c`(`x`) as (select 1 as `x`) select `c`.`x` from `c`"},
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
