package catalog

import (
	"errors"
	"testing"

	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

func TestViewQueryAndExpressionValidation(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE db1; CREATE DATABASE db2; USE db1;")
		mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, s VARCHAR(20), b VARBINARY(20), j JSON, ts TIMESTAMP, dt DATETIME(3), g GEOMETRY, pt POINT);")
		mustExec(t, c, "CREATE TABLE db2.t (id INT PRIMARY KEY, b INT);")
		mustExec(t, c, "CREATE FUNCTION db1.my_fn(x INT) RETURNS INT DETERMINISTIC RETURN x;")
		mustExec(t, c, "CREATE FUNCTION db2.other_fn(x INT) RETURNS INT DETERMINISTIC RETURN x;")
		return c
	}

	refusals := []struct {
		name string
		sql  string
		want int
	}{
		// Parameter markers in DDL report 1064.
		{"param_marker_select", "CREATE VIEW v AS SELECT ? FROM t", 1064},
		{"param_marker_limit", "CREATE VIEW v AS SELECT a FROM t LIMIT ?", 1064},
		{"param_marker_ntile", "CREATE VIEW v AS SELECT NTILE(?) OVER () FROM t", 1064},

		// SELECT ... INTO is refused with 1350 ahead of 1054 in position 1.
		{"select_into_var", "CREATE VIEW v AS SELECT missing INTO @x", ErrViewSelectClause},
		{"select_into_outfile", "CREATE VIEW v AS SELECT missing INTO OUTFILE 'f'", ErrViewSelectClause},
		{"select_into_before_for_update", "CREATE VIEW v AS SELECT missing INTO @x FOR UPDATE", ErrViewSelectClause},

		// Unknown functions (1305), resolved ahead of argument column lookup (1054).
		{"unknown_func_no_args", "CREATE VIEW v AS SELECT no_such_fn() FROM t", ErrNoSuchFunction},
		{"unknown_func_schema_qualified", "CREATE VIEW v AS SELECT db2.no_such_fn() FROM t", ErrNoSuchFunction},
		{"unknown_func_missing_arg_col", "CREATE VIEW v AS SELECT no_such_fn(missing_col) FROM t", ErrNoSuchFunction},
		{"unknown_func_missing_schema", "CREATE VIEW v AS SELECT no_such_db.fn() FROM t", ErrNoSuchFunction},

		// Temporal literal validation (1525).
		{"invalid_date_literal", "CREATE VIEW v AS SELECT DATE 's' FROM t", ErrWrongValue},
		{"invalid_date_calendar", "CREATE VIEW v AS SELECT DATE '2020-02-30' FROM t", ErrWrongValue},
		{"invalid_date_zero", "CREATE VIEW v AS SELECT DATE '0000-00-00' FROM t", ErrWrongValue},
		{"invalid_date_with_time", "CREATE VIEW v AS SELECT DATE '2020-01-01 12:00:00' FROM t", ErrWrongValue},
		{"invalid_date_colon_with_time", "CREATE VIEW v AS SELECT DATE '2020:01:01:12:00:00' FROM t", ErrWrongValue},
		{"invalid_time_literal", "CREATE VIEW v AS SELECT TIME 's' FROM t", ErrWrongValue},
		{"invalid_time_overflow", "CREATE VIEW v AS SELECT TIME '900:00:00' FROM t", ErrWrongValue},
		{"invalid_timestamp_literal", "CREATE VIEW v AS SELECT TIMESTAMP 's' FROM t", ErrWrongValue},
		{"invalid_timestamp_date_only", "CREATE VIEW v AS SELECT TIMESTAMP '2020-01-01' FROM t", ErrWrongValue},
		{"invalid_timestamp_colon_date_only", "CREATE VIEW v AS SELECT TIMESTAMP '2020:01:01' FROM t", ErrWrongValue},
		{"invalid_timestamp_dotted_date_only", "CREATE VIEW v AS SELECT TIMESTAMP '2020.01.01' FROM t", ErrWrongValue},
		{"invalid_timestamp_zero", "CREATE VIEW v AS SELECT TIMESTAMP '0000-00-00 00:00:00' FROM t", ErrWrongValue},

		// CAST(... AT TIME ZONE ... AS DATETIME) requires TIMESTAMP or NULL (3998).
		{"cast_at_tz_int", "CREATE VIEW v AS SELECT CAST(a AT TIME ZONE 'UTC' AS DATETIME) FROM t", ErrInvalidCastToTimestampWithTimeZone},
		{"cast_at_tz_interval_int", "CREATE VIEW v AS SELECT CAST(a AT TIME ZONE INTERVAL '+00:00' AS DATETIME) FROM t", ErrInvalidCastToTimestampWithTimeZone},
		{"cast_at_tz_datetime", "CREATE VIEW v AS SELECT CAST(dt AT TIME ZONE 'UTC' AS DATETIME) FROM t", ErrInvalidCastToTimestampWithTimeZone},
		{"cast_at_tz_varchar", "CREATE VIEW v AS SELECT CAST(s AT TIME ZONE 'UTC' AS DATETIME) FROM t", ErrInvalidCastToTimestampWithTimeZone},
		{"cast_at_tz_json_extract", "CREATE VIEW v AS SELECT CAST((j->'$.ts') AT TIME ZONE 'UTC' AS DATETIME) FROM t", ErrInvalidCastToTimestampWithTimeZone},

		// JSON_VALUE argument 1 requires JSON, string, binary, or NULL (3146).
		{"json_value_int_col", "CREATE VIEW v AS SELECT JSON_VALUE(a, '$.x') FROM t", ErrInvalidTypeForJSON},
		{"json_value_int_lit", "CREATE VIEW v AS SELECT JSON_VALUE(123, '$.x') FROM t", ErrInvalidTypeForJSON},
		{"json_value_arith_op", "CREATE VIEW v AS SELECT JSON_VALUE((a + 1), '$.x') FROM t", ErrInvalidTypeForJSON},
		{"json_value_bitand_int", "CREATE VIEW v AS SELECT JSON_VALUE((a & 1), '$.x') FROM t", ErrInvalidTypeForJSON},
		{"json_value_returning_int_col", "CREATE VIEW v AS SELECT JSON_VALUE(a, 's' RETURNING BINARY) FROM t", ErrInvalidTypeForJSON},
		{"json_value_on_empty_int_col", "CREATE VIEW v AS SELECT JSON_VALUE(a, 's' ERROR ON EMPTY) FROM t", ErrInvalidTypeForJSON},

		// Geometry collection constructors require geometric arguments (1367).
		{"linestring_int", "CREATE VIEW v AS SELECT LINESTRING(a) FROM t", ErrIllegalValueForType},
		{"linestring_json_extract", "CREATE VIEW v AS SELECT LINESTRING(j->'$.p') FROM t", ErrIllegalValueForType},
		{"polygon_int", "CREATE VIEW v AS SELECT POLYGON(a, a) FROM t", ErrIllegalValueForType},
		{"multipoint_varchar", "CREATE VIEW v AS SELECT MULTIPOINT(s) FROM t", ErrIllegalValueForType},
		{"multilinestring_null", "CREATE VIEW v AS SELECT MULTILINESTRING(NULL) FROM t", ErrIllegalValueForType},
		{"multipolygon_int", "CREATE VIEW v AS SELECT MULTIPOLYGON(a) FROM t", ErrIllegalValueForType},
		{"geomcollection_int", "CREATE VIEW v AS SELECT GEOMETRYCOLLECTION(a) FROM t", ErrIllegalValueForType},

		// DEFAULT(col) in CREATE VIEW: 1054 missing col, 1364 no default, 1143 top-level select item.
		{"default_missing_col", "CREATE VIEW v AS SELECT DEFAULT(no_col) FROM t", ErrNoSuchColumn},
		{"default_no_default_col", "CREATE VIEW v AS SELECT DEFAULT(id) FROM t", ErrNoDefaultForField},
		{"default_top_level_select", "CREATE VIEW v AS SELECT DEFAULT(a) FROM t", ErrColumnAccessDenied},
		{"default_top_level_aliased", "CREATE VIEW v AS SELECT DEFAULT(a) AS x FROM t", ErrColumnAccessDenied},

		// Undeclared variable (1327) and user variable (1351) in LIMIT / NTILE.
		{"limit_ident", "CREATE VIEW v AS SELECT missing LIMIT n", ErrSPUndeclaredVar},
		{"ntile_ident_before_window", "CREATE VIEW v AS SELECT NTILE(n) OVER no_win FROM t", ErrSPUndeclaredVar},
		{"ntile_user_var_before_window", "CREATE VIEW v AS SELECT NTILE(@x) OVER no_win FROM t", ErrViewSelectVariable},
		{"user_var_in_select", "CREATE VIEW v AS SELECT @x FROM t", ErrViewSelectVariable},

		// Locking clause unresolved table name (3568).
		{"for_update_of_missing", "CREATE VIEW v AS SELECT missing FOR UPDATE OF t", ErrUnresolvedTableLock},
		{"for_update_of_wrong_schema", "CREATE VIEW v AS SELECT a FROM t FOR UPDATE OF db2.t", ErrUnresolvedTableLock},
		{"for_update_of_hidden_by_alias", "CREATE VIEW v AS SELECT u.a FROM t AS u FOR UPDATE OF t", ErrUnresolvedTableLock},

		// Window validation (3579, 1235, 3584, 3585, 3596).
		{"undefined_window", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER w FROM t", ErrWindowNoSuchWindow},
		{"groups_frame", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (GROUPS UNBOUNDED PRECEDING) FROM t", ErrNotSupportedYet},
		{"exclude_frame", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS UNBOUNDED PRECEDING EXCLUDE CURRENT ROW) FROM t", ErrNotSupportedYet},
		{"rows_interval", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS INTERVAL 1 DAY PRECEDING) FROM t", ErrWindowRangeFrameNumericType},
		{"rows_interval_before_end", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN INTERVAL 1 DAY FOLLOWING AND UNBOUNDED PRECEDING) FROM t", ErrWindowRangeFrameNumericType},
		{"frame_start_unbounded_following", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN UNBOUNDED FOLLOWING AND UNBOUNDED PRECEDING) FROM t", ErrWindowFrameStartIllegal},
		{"frame_end_unbounded_preceding", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN 1 FOLLOWING AND UNBOUNDED PRECEDING) FROM t", ErrWindowFrameEndIllegal},

		// VALUES query primary errors (3942, 3943), user var (1351), and unknown column (1054).
		{"values_empty_row", "CREATE VIEW v AS VALUES ROW()", ErrEmptyRowInTVC},
		{"values_default", "CREATE VIEW v AS VALUES ROW(DEFAULT)", ErrDefaultAsValInTVC},
		{"values_user_var", "CREATE VIEW v AS VALUES ROW(@x)", ErrViewSelectVariable},
		{"values_unknown_col", "CREATE VIEW v AS VALUES ROW(no_col)", ErrNoSuchColumn},

		// Top-level ROW() in target list (1241).
		{"row_in_select_list", "CREATE VIEW v AS SELECT ROW(a, id) FROM t", ErrOperandColumns},

		// GROUPING() without WITH ROLLUP (1111).
		{"grouping_no_rollup", "CREATE VIEW v AS SELECT GROUPING(a) FROM t", ErrInvalidGroupFuncUse},
		{"grouping_group_by_no_rollup", "CREATE VIEW v AS SELECT GROUPING(a) FROM t GROUP BY a", ErrInvalidGroupFuncUse},

		// Non-constant ESCAPE / AGAINST (1210).
		{"like_escape_col", "CREATE VIEW v AS SELECT s LIKE 'a%' ESCAPE s FROM t", ErrWrongArguments},
		{"match_against_col", "CREATE VIEW v AS SELECT MATCH(s) AGAINST(s) FROM t", ErrWrongArguments},

		// Duplicate view columns (1060).
		{"dup_derived_col", "CREATE VIEW v AS SELECT a, a FROM t", ErrDupColumn},
		{"dup_explicit_col", "CREATE VIEW v (x, X) AS SELECT id, a FROM t", ErrDupColumn},

		// Missing table (1146) and unknown column (1054).
		{"missing_table_select", "CREATE VIEW v AS SELECT a FROM missing_tbl", ErrNoSuchTable},
		{"missing_table_explicit", "CREATE VIEW v AS TABLE missing_tbl", ErrNoSuchTable},
		{"dual_unknown_col", "CREATE VIEW v AS SELECT a FROM DUAL", ErrNoSuchColumn},
		{"unknown_col_from_table", "CREATE VIEW v AS SELECT no_col FROM t", ErrNoSuchColumn},
		{"wrong_schema_col", "CREATE VIEW v AS SELECT db2.t.a FROM t", ErrNoSuchColumn},
		{"derived_table_unknown_col", "CREATE VIEW v AS SELECT d.no_col FROM (SELECT 1 + 1) AS d", ErrNoSuchColumn},
		{"non_lateral_sibling_table_ref", "CREATE VIEW v AS SELECT d.x FROM t, (SELECT t.a AS x) AS d", ErrNoSuchColumn},
		{"select_alias_same_list_ref", "CREATE VIEW v AS SELECT a + 1 AS x, x + 1 AS y FROM t", ErrNoSuchColumn},
		{"select_alias_where_subquery_ref", "CREATE VIEW v AS SELECT a + 1 AS x FROM t WHERE (SELECT x) = 1", ErrNoSuchColumn},
		{"select_alias_forward_subquery_ref", "CREATE VIEW v AS SELECT (SELECT x + 1) AS y, a + 1 AS x FROM t", ErrNoSuchColumn},
		{"select_alias_qualified_subquery_ref", "CREATE VIEW v AS SELECT a + 1 AS x, (SELECT t.x) AS y FROM t", ErrNoSuchColumn},
		{"reachable_cte_unknown_col", "CREATE VIEW v AS WITH bad AS (SELECT missing_col), used AS (SELECT * FROM bad) SELECT * FROM used", ErrNoSuchColumn},
	}

	for _, tc := range refusals {
		t.Run(tc.name, func(t *testing.T) {
			c := setup(t)
			if got := execErrorCode(t, c, tc.sql); got != tc.want {
				t.Errorf("%s: got error %d, want %d", tc.sql, got, tc.want)
			}
		})
	}

	// Trailing INTO after FROM / LIMIT / locking clauses in CREATE/ALTER VIEW
	// is rejected by the grammar with 1064 (mysql:8.0.46).
	parseRefusals := []struct {
		name string
		sql  string
	}{
		{"select_for_update_into_syntax", "CREATE VIEW v AS SELECT a FOR UPDATE INTO @x"},
		{"select_from_into_syntax", "CREATE VIEW v AS SELECT a FROM t INTO @x"},
		{"select_limit_into_syntax", "CREATE VIEW v AS SELECT a FROM t LIMIT 1 INTO @x"},
		{"alter_view_for_update_into_syntax", "ALTER VIEW v AS SELECT a FOR UPDATE INTO @x"},
	}
	for _, tc := range parseRefusals {
		t.Run("syntax/"+tc.name, func(t *testing.T) {
			c := setup(t)
			_, err := c.Exec(tc.sql, nil)
			var pe *mysqlparser.ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("%s: expected *parser.ParseError (1064), got %v", tc.sql, err)
			}
		})
	}

	// Valid cases that must continue to succeed.
	accepted := []struct {
		name string
		sql  string
	}{
		{"system_var", "CREATE VIEW v AS SELECT @@version FROM t"},
		{"for_update_of_table", "CREATE VIEW v AS SELECT a FROM t FOR UPDATE OF t"},
		{"for_update_of_current_db_table", "CREATE VIEW v AS SELECT a FROM t FOR UPDATE OF db1.t"},
		{"for_update_of_alias", "CREATE VIEW v AS SELECT u.a FROM t AS u FOR UPDATE OF u"},
		{"defined_window", "CREATE VIEW v AS SELECT ROW_NUMBER() OVER w FROM t WINDOW w AS (ORDER BY id)"},
		{"grouping_with_rollup", "CREATE VIEW v AS SELECT a, GROUPING(a) FROM t GROUP BY a WITH ROLLUP"},
		{"like_escape_const", "CREATE VIEW v AS SELECT s LIKE 'a%' ESCAPE '\\\\' FROM t"},
		{"match_against_const", "CREATE VIEW v AS SELECT MATCH(s) AGAINST('foo') FROM t"},
		{"explicit_cols_override_dup", "CREATE VIEW v (c1, c2) AS SELECT a, a FROM t"},
		{"multi_schema_same_table_name", "CREATE VIEW v AS SELECT db1.t.a, db2.t.b FROM db1.t JOIN db2.t ON db1.t.id = db2.t.id"},
		{"cte_chaining", "CREATE VIEW v AS WITH c1 AS (SELECT a FROM t), c2 AS (SELECT a FROM c1) SELECT a FROM c2"},
		{"cte_in_scalar_subquery", "CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT (SELECT x FROM c) AS x"},
		{"cte_in_derived_table", "CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT * FROM (SELECT x FROM c) d"},
		{"cte_in_exists_subquery", "CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT EXISTS(SELECT x FROM c) AS ok"},
		{"cte_in_in_subquery", "CREATE VIEW v AS WITH c AS (SELECT 1 AS x) SELECT a IN (SELECT x FROM c) AS ok FROM t"},
		{"cte_in_paren_query", "CREATE VIEW v AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c)"},
		{"cte_in_values_subquery", "CREATE VIEW v AS WITH c AS (SELECT 1 AS x) VALUES ROW((SELECT x FROM c))"},
		{"unused_cte_missing_column", "CREATE VIEW v AS WITH unused AS (SELECT missing_column) SELECT 1 AS x"},
		{"unused_cte_missing_table", "CREATE VIEW v AS WITH unused AS (SELECT * FROM missing_tbl) SELECT 1 AS x"},
		{"derived_expr_col_name", "CREATE VIEW v AS SELECT d.`1 + 1` FROM (SELECT 1 + 1) AS d"},
		{"derived_func_col_name", "CREATE VIEW v AS SELECT d.`COUNT(*)` FROM (SELECT COUNT(*) FROM t) AS d"},
		{"cte_expr_col_name", "CREATE VIEW v AS WITH c AS (SELECT 1 + 1) SELECT c.`1 + 1` FROM c"},
		{"non_lateral_derived_outer_scope", "CREATE VIEW v AS SELECT (SELECT d.x FROM (SELECT t.a AS x) AS d) AS x FROM t"},
		{"lateral_derived_sibling_scope", "CREATE VIEW v AS SELECT d.x FROM t, LATERAL (SELECT t.a AS x) AS d"},
		{"preceding_select_alias_in_subquery", "CREATE VIEW v AS SELECT a + 1 AS x, (SELECT x + 1) AS y FROM t"},
		{"preceding_select_alias_nested_subquery", "CREATE VIEW v AS SELECT a + 1 AS x, (SELECT (SELECT x + 1)) AS y FROM t"},
		{"preceding_select_alias_in_having_subquery", "CREATE VIEW v AS SELECT a + 1 AS x FROM t GROUP BY a HAVING (SELECT x) > 0"},
		{"recursive_cte_parenthesized", "CREATE VIEW v AS WITH RECURSIVE c(x) AS ((SELECT 1 UNION ALL SELECT x + 1 FROM c WHERE x < 3)) SELECT x FROM c"},
		{"recursive_cte_parenthesized_arms", "CREATE VIEW v AS WITH RECURSIVE c(x) AS ((SELECT 1) UNION ALL (SELECT x + 1 FROM c WHERE x < 3)) SELECT x FROM c"},
		{"recursive_cte_multi_arm", "CREATE VIEW v AS WITH RECURSIVE c(x) AS (SELECT 1 UNION ALL SELECT x + 1 FROM c WHERE x < 3 UNION ALL SELECT x + 2 FROM c WHERE x < 3) SELECT x FROM c"},
		{"qualified_star_using_join", "CREATE VIEW v AS SELECT d.id, d.b FROM (SELECT r.* FROM db1.t AS l JOIN db2.t AS r USING (id)) AS d"},
		{"qualified_star_natural_join", "CREATE VIEW v AS SELECT d.id, d.b FROM (SELECT r.* FROM db1.t AS l NATURAL JOIN db2.t AS r) AS d"},
		{"schema_qualified_star_multi_schema", "CREATE VIEW v AS SELECT d.b FROM (SELECT db2.t.* FROM db1.t JOIN db2.t ON db1.t.id = db2.t.id) AS d"},
		{"dual_literal", "CREATE VIEW v AS SELECT 1 FROM DUAL"},
		// Registered stored functions (current DB and cross-schema) and ST_COLLECT.
		{"stored_func_current_db", "CREATE VIEW v AS SELECT my_fn(a) FROM t"},
		{"stored_func_qualified", "CREATE VIEW v AS SELECT db2.other_fn(a) FROM t"},
		{"st_collect_builtin", "CREATE VIEW v AS SELECT ST_COLLECT(pt) FROM t"},
		// Valid temporal literals.
		{"valid_date_literal", "CREATE VIEW v AS SELECT DATE '2020-01-01' FROM t"},
		{"valid_date_colon_literal", "CREATE VIEW v AS SELECT DATE '2020:01:01' AS d FROM t"},
		{"valid_time_literal", "CREATE VIEW v AS SELECT TIME '-838:59:59' FROM t"},
		{"valid_timestamp_literal", "CREATE VIEW v AS SELECT TIMESTAMP '2020-01-01 12:34:56.123456' FROM t"},
		{"valid_timestamp_colon_date", "CREATE VIEW v AS SELECT TIMESTAMP '2020:01:01 12:34:56' FROM t"},
		{"valid_timestamp_dotted_date", "CREATE VIEW v AS SELECT TIMESTAMP '2020.01.01 12:34:56' AS ts FROM t"},
		// Valid CAST AT TIME ZONE, JSON_VALUE, geometry constructors, and non-top-level DEFAULT(col).
		{"cast_at_tz_timestamp", "CREATE VIEW v AS SELECT CAST(ts AT TIME ZONE 'UTC' AS DATETIME) FROM t"},
		{"cast_at_tz_null", "CREATE VIEW v AS SELECT CAST(NULL AT TIME ZONE 'UTC' AS DATETIME) FROM t"},
		{"json_value_json_col", "CREATE VIEW v AS SELECT JSON_VALUE(j, '$.a') FROM t"},
		{"json_value_bit_literal", "CREATE VIEW v AS SELECT JSON_VALUE(b'00110001', '$') AS x"},
		{"json_value_hex_literal", "CREATE VIEW v AS SELECT JSON_VALUE(0x31, '$') AS x"},
		{"json_value_json_extract", "CREATE VIEW v AS SELECT JSON_VALUE(j->'$.obj', '$.a') AS a FROM t"},
		{"json_value_json_unquote", "CREATE VIEW v AS SELECT JSON_VALUE(j->>'$.obj', '$.a') AS a FROM t"},
		{"json_value_binary_unary", "CREATE VIEW v AS SELECT JSON_VALUE(BINARY s, '$.a') AS a FROM t"},
		{"json_value_binary_introducer_bitor", "CREATE VIEW v AS SELECT JSON_VALUE((_binary'1' | _binary'0'), '$') AS x"},
		{"json_value_binary_introducer_bitnot", "CREATE VIEW v AS SELECT JSON_VALUE(~_binary'1', '$') AS x"},
		{"json_value_binary_introducer_bitshift", "CREATE VIEW v AS SELECT JSON_VALUE((_binary'1' << 1), '$') AS x"},
		{"json_value_cast_char_charset_binary_bitor", "CREATE VIEW v AS SELECT JSON_VALUE((CAST('1' AS CHAR CHARSET binary) | CAST('0' AS CHAR CHARSET binary)), '$') AS x"},
		{"json_value_concat_binary_bitor", "CREATE VIEW v AS SELECT JSON_VALUE((CONCAT(_binary'1', _binary'0') | _binary'0'), '$') AS x"},
		{"json_value_bitand_binary", "CREATE VIEW v AS SELECT JSON_VALUE((b & b), '$.a') AS a FROM t"},
		{"json_value_bitnot_binary", "CREATE VIEW v AS SELECT JSON_VALUE(~b, '$.a') AS a FROM t"},
		{"json_value_bitshift_binary", "CREATE VIEW v AS SELECT JSON_VALUE((b << 1), '$.a') AS a FROM t"},
		{"json_value_varchar_col", "CREATE VIEW v AS SELECT JSON_VALUE(s, '$.a') FROM t"},
		{"json_value_varbinary_col", "CREATE VIEW v AS SELECT JSON_VALUE(b, '$.a') FROM t"},
		{"json_value_null", "CREATE VIEW v AS SELECT JSON_VALUE(NULL, '$.a') FROM t"},
		{"linestring_geom_col", "CREATE VIEW v AS SELECT LINESTRING(g, pt) FROM t"},
		{"linestring_point_expr", "CREATE VIEW v AS SELECT LINESTRING(POINT(a, a), POINT(a, a)) FROM t"},
		{"geomcollection_empty", "CREATE VIEW v AS SELECT GEOMETRYCOLLECTION() FROM t"},
		{"default_in_expr", "CREATE VIEW v AS SELECT DEFAULT(a) + 1 FROM t"},
		{"default_in_where", "CREATE VIEW v AS SELECT a FROM t WHERE DEFAULT(a) IS NULL"},
	}
	for _, tc := range accepted {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			c := setup(t)
			mustExec(t, c, tc.sql)
		})
	}

	t.Run("cross_schema_explicit_table_view", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE VIEW db1.v AS TABLE db2.t")
		v := c.GetDatabase("db1").Views["v"]
		if v == nil {
			t.Fatal("view db1.v not found")
		}
		// A column of another schema's unaliased table is qualified by that
		// schema too: mysql:8.0.46 stores CREATE VIEW a.v AS TABLE o.t as
		// select `o`.`t`.`id` AS `id`,`o`.`t`.`b` AS `b` from `o`.`t`.
		wantDef := "select `db2`.`t`.`id` AS `id`,`db2`.`t`.`b` AS `b` from `db2`.`t`"
		if v.Definition != wantDef {
			t.Errorf("definition = %q, want %q", v.Definition, wantDef)
		}
		if len(v.Columns) != 2 || v.Columns[0] != "id" || v.Columns[1] != "b" {
			t.Errorf("columns = %v, want [id b]", v.Columns)
		}

		mustExec(t, c, "CREATE VIEW db1.v_cte AS WITH t AS (SELECT 9 AS z) TABLE db2.t")
		vCTE := c.GetDatabase("db1").Views["v_cte"]
		if vCTE == nil {
			t.Fatal("view db1.v_cte not found")
		}
		// MySQL drops the unreferenced CTE as well, which the catalog keeps.
		wantCTEDef := "with `t` as (select 9 AS `z`) select `db2`.`t`.`id` AS `id`,`db2`.`t`.`b` AS `b` from `db2`.`t`"
		if vCTE.Definition != wantCTEDef {
			t.Errorf("v_cte definition = %q, want %q", vCTE.Definition, wantCTEDef)
		}
		if len(vCTE.Columns) != 2 || vCTE.Columns[0] != "id" || vCTE.Columns[1] != "b" {
			t.Errorf("v_cte columns = %v, want [id b]", vCTE.Columns)
		}
	})

	t.Run("stored_func_shadowing_builtin_return_type", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE FUNCTION db1.abs() RETURNS VARCHAR(20) DETERMINISTIC RETURN '{\"x\":1}';")
		mustExec(t, c, "CREATE VIEW v AS SELECT JSON_VALUE(db1.abs(), '$.x') AS x")

		mustExec(t, c, "CREATE FUNCTION db1.point() RETURNS INT DETERMINISTIC RETURN 1;")
		if got := execErrorCode(t, c, "CREATE VIEW v_bad_geom AS SELECT LINESTRING(db1.point())"); got != ErrIllegalValueForType {
			t.Errorf("LINESTRING(db1.point()) got %d, want %d", got, ErrIllegalValueForType)
		}
	})

	t.Run("explicit_table_excludes_functional_index_hidden_cols", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE TABLE tf (a INT, INDEX ((a + 1)));")
		mustExec(t, c, "CREATE VIEW vf AS TABLE tf")
		v := c.GetDatabase("db1").Views["vf"]
		if v == nil {
			t.Fatal("view vf not found")
		}
		wantDef := "select `tf`.`a` AS `a` from `tf`"
		if v.Definition != wantDef {
			t.Errorf("definition = %q, want %q", v.Definition, wantDef)
		}
		if len(v.Columns) != 1 || v.Columns[0] != "a" {
			t.Errorf("columns = %v, want [a]", v.Columns)
		}
	})

	t.Run("default_week_format_preserved_in_view", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "SET SESSION default_week_format = 3;")
		mustExec(t, c, "CREATE VIEW vw AS SELECT WEEK('2021-01-01') AS w")
		v := c.GetDatabase("db1").Views["vw"]
		if v == nil {
			t.Fatal("view vw not found")
		}
		wantDef := "select week('2021-01-01',3) AS `w`"
		if v.Definition != wantDef {
			t.Errorf("definition = %q, want %q", v.Definition, wantDef)
		}

		mustExec(t, c, "SET SESSION default_week_format = DEFAULT;")
		mustExec(t, c, "CREATE VIEW vw_def AS SELECT WEEK('2021-01-01') AS w")
		vDef := c.GetDatabase("db1").Views["vw_def"]
		wantDefZero := "select week('2021-01-01',0) AS `w`"
		if vDef.Definition != wantDefZero {
			t.Errorf("default definition = %q, want %q", vDef.Definition, wantDefZero)
		}
	})

	t.Run("explicit_table_excludes_invisible_cols_but_explicit_ref_keeps_them", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE TABLE ti (a INT, secret INT INVISIBLE);")
		mustExec(t, c, "CREATE VIEW vi AS TABLE ti")
		v := c.GetDatabase("db1").Views["vi"]
		if v == nil {
			t.Fatal("view vi not found")
		}
		wantDef := "select `ti`.`a` AS `a` from `ti`"
		if v.Definition != wantDef {
			t.Errorf("TABLE ti definition = %q, want %q", v.Definition, wantDef)
		}
		if len(v.Columns) != 1 || v.Columns[0] != "a" {
			t.Errorf("TABLE ti columns = %v, want [a]", v.Columns)
		}

		mustExec(t, c, "CREATE VIEW vi_named AS SELECT secret FROM ti")
		vNamed := c.GetDatabase("db1").Views["vi_named"]
		wantNamed := "select `ti`.`secret` AS `secret` from `ti`"
		if vNamed.Definition != wantNamed {
			t.Errorf("SELECT secret definition = %q, want %q", vNamed.Definition, wantNamed)
		}

		mustExec(t, c, "CREATE TABLE ti2 (b INT, secret INT INVISIBLE);")
		mustExec(t, c, "CREATE VIEW vi_right_using AS SELECT secret AS x FROM ti RIGHT JOIN ti2 USING (secret)")
		vRightUsing := c.GetDatabase("db1").Views["vi_right_using"]
		wantRightUsing := "select `ti2`.`secret` AS `x` from (`ti2` left join `ti` on((`ti2`.`secret` = `ti`.`secret`)))"
		if vRightUsing.Definition != wantRightUsing {
			t.Errorf("RIGHT JOIN USING invisible definition = %q, want %q", vRightUsing.Definition, wantRightUsing)
		}
	})

	t.Run("dotted_timestamp_literal_time_fsp", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE VIEW v_ts_dot AS SELECT TIME(TIMESTAMP '2020.01.01 12:34:56.1') AS x")
		v := c.GetDatabase("db1").Views["v_ts_dot"]
		wantDef := "select cast(TIMESTAMP'2020.01.01 12:34:56.1' as time(1)) AS `x`"
		if v.Definition != wantDef {
			t.Errorf("definition = %q, want %q", v.Definition, wantDef)
		}
	})

	t.Run("stored_func_cast_as_year_not_rewritten_to_cast", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE FUNCTION cast_as_year(x INT) RETURNS INT DETERMINISTIC RETURN x + 1;")
		mustExec(t, c, "CREATE VIEW vc AS SELECT cast_as_year(2020) AS x")
		v := c.GetDatabase("db1").Views["vc"]
		if v == nil {
			t.Fatal("view vc not found")
		}
		wantDef := "select cast_as_year(2020) AS `x`"
		if v.Definition != wantDef {
			t.Errorf("definition = %q, want %q", v.Definition, wantDef)
		}
	})
}
