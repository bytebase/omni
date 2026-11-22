package deparse

import "testing"

func TestViewRewriteAndDeparseCanonicalForms(t *testing.T) {
	lookup := func(name string) *ResolverTable {
		if name == "t" {
			return &ResolverTable{
				Name: "t",
				Columns: []ResolverColumn{
					{Name: "a", Position: 1, DataType: "int"},
					{Name: "b", Position: 2, DataType: "int"},
					{Name: "dt3", Position: 3, DataType: "datetime", Scale: 3},
					{Name: "s", Position: 4, DataType: "varchar"},
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
			"time_timestamp_week_canonical",
			"SELECT TIME(a), TIME(dt3), TIME('12:34:56.123456'), TIMESTAMP(a), TIMESTAMP(dt3), TIMESTAMP('2021-01-01 12:34:56.123456'), TIMESTAMP(a, '01:00:00'), WEEK(a), WEEK(a, 1) FROM t",
			"select cast(`t`.`a` as time) AS `TIME(a)`,cast(`t`.`dt3` as time(3)) AS `TIME(dt3)`,cast('12:34:56.123456' as time(6)) AS `TIME('12:34:56.123456')`,cast(`t`.`a` as datetime) AS `TIMESTAMP(a)`,cast(`t`.`dt3` as datetime(3)) AS `TIMESTAMP(dt3)`,cast('2021-01-01 12:34:56.123456' as datetime(6)) AS `TIMESTAMP('2021-01-01 12:34:56.123456')`,timestamp(`t`.`a`,'01:00:00') AS `TIMESTAMP(a, '01:00:00')`,week(`t`.`a`,0) AS `WEEK(a)`,week(`t`.`a`,1) AS `WEEK(a, 1)` from `t`",
		},
		{
			"time_timestamp_dotted_and_compact_literals",
			"SELECT TIME(TIMESTAMP '2020.01.01 12:34:56.1') AS a, TIME(TIMESTAMP '2020.01.01 12:34:56') AS b, TIME(TIMESTAMP '2020.01.01 12.34.56') AS c, TIME(TIMESTAMP '2020.01.01 12.34.56.12') AS d, TIME(TIMESTAMP '2020.01.01 12:34:56.12+02:00') AS e, TIME(TIMESTAMP '20200101123456.123') AS f, TIME(TIME '12:34:56.12') AS g",
			"select cast(TIMESTAMP'2020.01.01 12:34:56.1' as time(1)) AS `a`,cast(TIMESTAMP'2020.01.01 12:34:56' as time) AS `b`,cast(TIMESTAMP'2020.01.01 12.34.56' as time) AS `c`,cast(TIMESTAMP'2020.01.01 12.34.56.12' as time(2)) AS `d`,cast(TIMESTAMP'2020.01.01 12:34:56.12+02:00' as time(2)) AS `e`,cast(TIMESTAMP'20200101123456.123' as time(3)) AS `f`,cast(TIME'12:34:56.12' as time(2)) AS `g`",
		},
		{
			"position_to_locate",
			"SELECT POSITION(s IN s) FROM t",
			"select locate(`t`.`s`,`t`.`s`) AS `POSITION(s IN s)` from `t`",
		},
		{
			"cast_decimal_defaults",
			"SELECT CAST(a AS DECIMAL), CAST(a AS DECIMAL(5)), CAST(a AS DECIMAL(5,2)) FROM t",
			"select cast(`t`.`a` as decimal(10,0)) AS `CAST(a AS DECIMAL)`,cast(`t`.`a` as decimal(5,0)) AS `CAST(a AS DECIMAL(5))`,cast(`t`.`a` as decimal(5,2)) AS `CAST(a AS DECIMAL(5,2))` from `t`",
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
			"cast_as_year_stored_func_preserved_by_default",
			"SELECT cast_as_year(a) AS x FROM t",
			"select cast_as_year(`t`.`a`) AS `x` from `t`",
		},
		{
			"explicit_table_expands",
			"TABLE t",
			"select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`dt3` AS `dt3`,`t`.`s` AS `s` from `t`",
		},
		{
			"explicit_table_schema_qualified_without_schema_lookup_preserved",
			"TABLE db2.t",
			"table `db2`.`t`",
		},
		{
			"introducer_byte_literals_auto_alias",
			"SELECT _utf8mb4 0x41, _utf8mb4 b'1000001' FROM t",
			"select _utf8mb4'A' AS `_utf8mb4 0x41`,_utf8mb4'A' AS `_utf8mb4 b'1000001'` from `t`",
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

	t.Run("explicit_table_schema_qualified_with_schema_lookup", func(t *testing.T) {
		schemaLookup := func(schema, name string) *ResolverTable {
			if schema == "db2" && name == "t" {
				return &ResolverTable{
					Name: "t",
					Columns: []ResolverColumn{
						{Name: "x", Position: 1, DataType: "int"},
					},
				}
			}
			return nil
		}
		sql := "TABLE db2.t"
		stmt := parseSelect(t, sql)
		r := &Resolver{Lookup: lookup, SchemaLookup: schemaLookup, DefaultCharset: "utf8mb4", Source: sql}
		r.Resolve(stmt)
		RewriteSelectStmt(stmt)
		got := DeparseSelect(stmt)
		want := "select `t`.`x` AS `x` from `db2`.`t`"
		if got != want {
			t.Errorf("got:\n  %s\nwant:\n  %s", got, want)
		}
	})

	t.Run("week_respects_default_week_format", func(t *testing.T) {
		sql := "SELECT WEEK('2021-01-01') AS w"
		stmt := parseSelect(t, sql)
		r := &Resolver{Lookup: lookup, DefaultCharset: "utf8mb4", DefaultWeekFormat: 3, Source: sql}
		r.Resolve(stmt)
		RewriteSelectStmt(stmt)
		if got, want := DeparseSelect(stmt), "select week('2021-01-01',3) AS `w`"; got != want {
			t.Errorf("mode 3: got %q, want %q", got, want)
		}

		stmtUnknown := parseSelect(t, sql)
		rUnknown := &Resolver{Lookup: lookup, DefaultCharset: "utf8mb4", DefaultWeekFormat: -1, Source: sql}
		rUnknown.Resolve(stmtUnknown)
		RewriteSelectStmt(stmtUnknown)
		if got, want := DeparseSelect(stmtUnknown), "select week('2021-01-01') AS `w`"; got != want {
			t.Errorf("unknown mode: got %q, want %q", got, want)
		}
	})

	t.Run("cast_as_year_readback_converges_when_enabled", func(t *testing.T) {
		sql := "SELECT `cast_as_year`(`t`.`a`) AS `CAST(a AS YEAR)` FROM `t`"
		stmt := parseSelect(t, sql)
		r := &Resolver{Lookup: lookup, DefaultCharset: "utf8mb4", NormalizeCastAsYear: true, Source: sql}
		r.Resolve(stmt)
		RewriteSelectStmt(stmt)
		if got, want := DeparseSelect(stmt), "select cast(`t`.`a` as year) AS `CAST(a AS YEAR)` from `t`"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("invisible_columns_using_vs_natural_and_star", func(t *testing.T) {
		invLookup := func(name string) *ResolverTable {
			switch name {
			case "t1":
				return &ResolverTable{
					Name: "t1",
					Columns: []ResolverColumn{
						{Name: "a", Position: 1, DataType: "int"},
						{Name: "k", Position: 2, DataType: "int", Invisible: true},
					},
				}
			case "t2":
				return &ResolverTable{
					Name: "t2",
					Columns: []ResolverColumn{
						{Name: "b", Position: 1, DataType: "int"},
						{Name: "k", Position: 2, DataType: "int", Invisible: true},
					},
				}
			case "t3":
				return &ResolverTable{
					Name: "t3",
					Columns: []ResolverColumn{
						{Name: "c", Position: 1, DataType: "int"},
						{Name: "k", Position: 2, DataType: "int"},
					},
				}
			}
			return nil
		}
		invCases := []struct {
			name string
			sql  string
			want string
		}{
			{
				"right_join_using_invisible_select_k",
				"SELECT k AS x FROM t1 RIGHT JOIN t2 USING (k)",
				"select `t2`.`k` AS `x` from (`t2` left join `t1` on((`t2`.`k` = `t1`.`k`)))",
			},
			{
				"right_join_using_invisible_star",
				"SELECT * FROM t1 RIGHT JOIN t2 USING (k)",
				"select `t2`.`k` AS `k`,`t2`.`b` AS `b`,`t1`.`a` AS `a` from (`t2` left join `t1` on((`t2`.`k` = `t1`.`k`)))",
			},
			{
				"inner_join_using_both_invisible_star",
				"SELECT * FROM t1 JOIN t2 USING (k)",
				"select `t1`.`k` AS `k`,`t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2` on((`t1`.`k` = `t2`.`k`)))",
			},
			{
				"inner_join_using_left_invisible_right_visible_star",
				"SELECT * FROM t1 JOIN t3 USING (k)",
				"select `t1`.`k` AS `k`,`t1`.`a` AS `a`,`t3`.`c` AS `c` from (`t1` join `t3` on((`t1`.`k` = `t3`.`k`)))",
			},
			{
				"natural_join_both_invisible_not_joined",
				"SELECT * FROM t1 NATURAL JOIN t2",
				"select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2`)",
			},
			{
				"natural_join_mixed_invisible_not_joined",
				"SELECT * FROM t1 NATURAL JOIN t3",
				"select `t1`.`a` AS `a`,`t3`.`c` AS `c`,`t3`.`k` AS `k` from (`t1` join `t3`)",
			},
			{
				"nested_join_natural_skips_invisible_binds_visible",
				"SELECT * FROM (t1 JOIN t3 ON t1.a = t3.c) NATURAL JOIN (SELECT 1 AS k) s",
				"select `t3`.`k` AS `k`,`t1`.`a` AS `a`,`t3`.`c` AS `c` from ((`t1` join `t3` on((`t1`.`a` = `t3`.`c`))) join (select 1 AS `k`) `s` on((`t3`.`k` = `s`.`k`)))",
			},
		}
		for _, tc := range invCases {
			t.Run(tc.name, func(t *testing.T) {
				stmt := parseSelect(t, tc.sql)
				r := &Resolver{Lookup: invLookup, DefaultCharset: "utf8mb4", Source: tc.sql}
				r.Resolve(stmt)
				RewriteSelectStmt(stmt)
				if got := DeparseSelect(stmt); got != tc.want {
					t.Errorf("got:\n  %s\nwant:\n  %s", got, tc.want)
				}
			})
		}
	})
}
