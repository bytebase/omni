package deparse

import "testing"

// TestAutoAliasWithoutSource: with no source text, DeparseSelect rebuilds an
// unaliased subquery's name from the AST. It used to keep only the select
// list, FROM and WHERE of a plain query, so a WITH, a set operation or a
// trailing clause vanished from the name; and writing a set operation cleared
// its WITH from the AST, which dropped it from the definition too.
func TestAutoAliasWithoutSource(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{
			"SELECT (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)",
			"select (with `c` as (select 1 AS `x`) select `x` from `c` union select 2) AS `(WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)`",
		},
		{
			"SELECT (SELECT a FROM t GROUP BY a HAVING a > 0 LIMIT 1), (SELECT a FROM t ORDER BY a DESC LIMIT 2, 1)",
			"select (select `a` from `t` group by `a` having (`a` > 0) limit 1) AS `(SELECT a FROM t GROUP BY a HAVING a > 0 LIMIT 1)`,(select `a` from `t` order by `a` desc limit 2,1) AS `(SELECT a FROM t ORDER BY a DESC LIMIT 2, 1)`",
		},
		{
			"SELECT EXISTS (SELECT DISTINCT a FROM t UNION ALL SELECT 1)",
			"select exists(select distinct `a` from `t` union all select 1) AS `EXISTS(SELECT DISTINCT a FROM t UNION ALL SELECT 1)`",
		},
		{
			"SELECT (VALUES ROW(1) ORDER BY 1 LIMIT 1)",
			"select (values row(1) order by 1 limit 1) AS `(VALUES ROW(1) ORDER BY 1 LIMIT 1)`",
		},
	} {
		stmt := parseSelectOrFatal(t, tc.sql, tc.sql)
		if got := DeparseSelect(stmt); got != tc.want {
			t.Errorf("DeparseSelect(%s)\n got: %s\nwant: %s", tc.sql, got, tc.want)
		}
		if got := DeparseSelect(stmt); got != tc.want {
			t.Errorf("DeparseSelect(%s) again\n got: %s\nwant: %s", tc.sql, got, tc.want)
		}
	}
}
