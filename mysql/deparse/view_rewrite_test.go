package deparse

import "testing"

func TestViewRewriteAndDeparseCanonicalForms(t *testing.T) {
	lookup := func(name string) *ResolverTable {
		if name == "t" {
			return &ResolverTable{
				Name: "t",
				Columns: []ResolverColumn{
					{Name: "a", Position: 1},
					{Name: "b", Position: 2},
				},
			}
		}
		return nil
	}

	run := func(t *testing.T, sql string) string {
		t.Helper()
		stmt := parseSelect(t, sql)
		r := &Resolver{Lookup: lookup, DefaultCharset: "utf8mb4", Source: sql}
		r.Resolve(stmt)
		RewriteSelectStmt(stmt)
		return DeparseSelect(stmt)
	}

	cases := []struct {
		name string
		sql  string
		want string
	}{
		{
			"is_not_true_wraps_non_bool",
			"SELECT a IS NOT TRUE FROM t",
			"select ((0 <> `t`.`a`) is not true) AS `a IS NOT TRUE` from `t`",
		},
		{
			"is_not_false_wraps_non_bool",
			"SELECT a IS NOT FALSE FROM t",
			"select ((0 <> `t`.`a`) is not false) AS `a IS NOT FALSE` from `t`",
		},
		{
			"is_unknown_to_is_null",
			"SELECT a IS UNKNOWN, a IS NOT UNKNOWN FROM t",
			"select (`t`.`a` is null) AS `a IS UNKNOWN`,(`t`.`a` is not null) AS `a IS NOT UNKNOWN` from `t`",
		},
		{
			"single_element_in_not_in",
			"SELECT a IN (b), a NOT IN (b) FROM t",
			"select (`t`.`a` = `t`.`b`) AS `a IN (b)`,(`t`.`a` <> `t`.`b`) AS `a NOT IN (b)` from `t`",
		},
		{
			"not_regexp_round_trip",
			"SELECT a NOT REGEXP b FROM t",
			"select (not(regexp_like(`t`.`a`,`t`.`b`))) AS `a NOT REGEXP b` from `t`",
		},
		{
			"not_regexp_like_idempotent",
			"SELECT (not(regexp_like(`t`.`a`,`t`.`b`))) AS `x` FROM `t`",
			"select (not(regexp_like(`t`.`a`,`t`.`b`))) AS `x` from `t`",
		},
		{
			"correlated_subqueries",
			"SELECT EXISTS(SELECT a), a IN (SELECT b), a NOT IN (SELECT b), a = ALL (SELECT b) FROM t",
			"select exists(select `t`.`a`) AS `EXISTS(SELECT a)`,`t`.`a` in (select `t`.`b`) AS `a IN (SELECT b)`,(not(`t`.`a` in (select `t`.`b`))) AS `a NOT IN (SELECT b)`,`t`.`a` = all (select `t`.`b`) AS `a = ALL (SELECT b)` from `t`",
		},
		{
			"binary_unary_to_cast",
			"SELECT BINARY a FROM t",
			"select cast(`t`.`a` as char charset binary) AS `BINARY a` from `t`",
		},
		{
			"date_day_mod_geom_values",
			"SELECT DATE(a), DAY(a), MOD(a, b), GEOMETRYCOLLECTION(), VALUES(a) FROM t",
			"select cast(`t`.`a` as date) AS `DATE(a)`,dayofmonth(`t`.`a`) AS `DAY(a)`,(`t`.`a` % `t`.`b`) AS `MOD(a, b)`,geomcollection() AS `GEOMETRYCOLLECTION()`,NULL AS `VALUES(a)` from `t`",
		},
		{
			"time_timestamp_week_preserved",
			"SELECT TIME('12:34:56.123456'), TIMESTAMP('2021-01-01 12:34:56.123456'), WEEK('2021-01-01') FROM t",
			"select time('12:34:56.123456') AS `TIME('12:34:56.123456')`,timestamp('2021-01-01 12:34:56.123456') AS `TIMESTAMP('2021-01-01 12:34:56.123456')`,week('2021-01-01') AS `WEEK('2021-01-01')` from `t`",
		},
		{
			"schema_qualified_functions_not_rewritten",
			"SELECT db.DAY(a), db.GEOMETRYCOLLECTION(), db.SUBSTRING(a, 1) FROM t",
			"select db.day(`t`.`a`) AS `db.DAY(a)`,db.geometrycollection() AS `db.GEOMETRYCOLLECTION()`,db.substring(`t`.`a`,1) AS `db.SUBSTRING(a, 1)` from `t`",
		},
		{
			"min_max_drop_distinct",
			"SELECT MIN(DISTINCT a), MAX(DISTINCT b) FROM t",
			"select min(`t`.`a`) AS `MIN(DISTINCT a)`,max(`t`.`b`) AS `MAX(DISTINCT b)` from `t`",
		},
		{
			"cast_as_year_stays_executable",
			"SELECT CAST(a AS YEAR) FROM t",
			"select cast(`t`.`a` as year) AS `CAST(a AS YEAR)` from `t`",
		},
		{
			"window_single_bound_expansion",
			"SELECT ROW_NUMBER() OVER (ROWS UNBOUNDED PRECEDING), ROW_NUMBER() OVER (ROWS 1 PRECEDING), ROW_NUMBER() OVER (ROWS CURRENT ROW), ROW_NUMBER() OVER (RANGE UNBOUNDED PRECEDING), ROW_NUMBER() OVER (ORDER BY a RANGE UNBOUNDED PRECEDING) FROM t",
			"select row_number() OVER (ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)  AS `ROW_NUMBER() OVER (ROWS UNBOUNDED PRECEDING)`,row_number() OVER (ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)  AS `ROW_NUMBER() OVER (ROWS 1 PRECEDING)`,row_number() OVER (ROWS BETWEEN CURRENT ROW AND CURRENT ROW)  AS `ROW_NUMBER() OVER (ROWS CURRENT ROW)`,row_number() OVER (RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)  AS `ROW_NUMBER() OVER (RANGE UNBOUNDED PRECEDING)`,row_number() OVER (ORDER BY `t`.`a` RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)  AS `ROW_NUMBER() OVER (ORDER BY a RANGE UNBOUNDED PRECEDING)` from `t`",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := run(t, tc.sql)
			if got != tc.want {
				t.Errorf("got:\n  %s\nwant:\n  %s", got, tc.want)
			}
			// Verify deparsed SQL is idempotent on re-parse + re-resolve + rewrite + deparse.
			stmt2 := parseSelect(t, got)
			r2 := &Resolver{Lookup: lookup, DefaultCharset: "utf8mb4", Source: got}
			r2.Resolve(stmt2)
			RewriteSelectStmt(stmt2)
			if again := DeparseSelect(stmt2); again != got {
				t.Errorf("not idempotent:\n  first:  %s\n  second: %s", got, again)
			}
		})
	}
}
