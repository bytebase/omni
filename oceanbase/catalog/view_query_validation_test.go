package catalog

import (
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

func TestViewQueryAndExpressionValidation(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE db1; CREATE DATABASE db2; USE db1;")
		mustExec(t, c, "CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));")
		mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), v VECTOR(3));")
		return c
	}

	execCode := func(t *testing.T, c *Catalog, sql string) int {
		t.Helper()
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("unexpected parse error for %q: %v", sql, err)
		}
		for _, r := range results {
			if r.Error != nil {
				var ce *Error
				if errors.As(r.Error, &ce) {
					return ce.Code
				}
				t.Fatalf("unexpected non-catalog error for %q: %v", sql, r.Error)
			}
		}
		return 0
	}

	refusals := []struct {
		name string
		sql  string
		want int
	}{
		// Parameter markers in DDL report 1064.
		{"param_marker_select", "CREATE VIEW _v AS SELECT ?", 1064},
		{"param_marker_limit", "CREATE VIEW _v AS SELECT a FROM a LIMIT ?", 1064},
		{"param_marker_window_frame", "CREATE VIEW _v AS SELECT RANK() OVER(ROWS ? PRECEDING)", 1064},

		// Bare and qualified stars / column refs with missing tables.
		{"bare_star_no_from", "CREATE VIEW _v AS SELECT *", ErrNoTablesUsed},
		{"qualified_star_missing", "CREATE VIEW _v AS SELECT a.*", ErrUnknownTable},
		{"schema_qualified_star_missing", "CREATE VIEW _v AS SELECT a.a.*", ErrUnknownTable},
		{"qualified_col_missing", "CREATE VIEW _v AS SELECT a.a", ErrUnknownTableInFieldList},
		{"schema_qualified_col_missing", "CREATE VIEW _v AS SELECT a.a.a", ErrUnknownTableInFieldList},

		// Duplicate table/alias (1066) vs ambiguous column (1052).
		{"dup_from_table", "CREATE VIEW _v AS SELECT a FROM a, a", 1066},
		{"dup_join_table", "CREATE VIEW _v AS SELECT a FROM a JOIN a ON a", 1066},
		{"dup_join_alias", "CREATE VIEW _v AS SELECT a FROM a JOIN a a", 1066},
		{"ambiguous_col_distinct_alias", "CREATE VIEW _v AS SELECT a FROM a, a zzz", 1052},
		{"ambiguous_col_in_on", "CREATE VIEW _v AS SELECT a FROM a x JOIN a y ON a", 1052},

		// Window validation (5929, 1235, 3596, 1149).
		{"undefined_window_bare", "CREATE VIEW _v AS SELECT RANK() OVER w", ErrOBWindowNotDefined},
		{"undefined_window_paren", "CREATE VIEW _v AS SELECT RANK() OVER (w)", ErrOBWindowNotDefined},
		{"undefined_window_before_col", "CREATE VIEW _v AS SELECT COUNT(missing) OVER w", ErrOBWindowNotDefined},
		{"duplicate_window_name", "CREATE VIEW _v AS SELECT a WINDOW w AS (), w AS ()", ErrNotSupportedYet},
		{"rows_interval_frame", "CREATE VIEW _v AS SELECT RANK() OVER (ROWS INTERVAL 1 DAY PRECEDING)", ErrWindowRangeFrameNumericType},
		{"single_following_bound", "CREATE VIEW _v AS SELECT RANK() OVER (ROWS 1 FOLLOWING)", ErrOBSyntaxError},

		// Unknown charset (1115) and collation (1273).
		{"unknown_collation", "CREATE VIEW _v AS SELECT a COLLATE no_such_coll", ErrUnknownCollation},
		{"unknown_charset_convert", "CREATE VIEW _v AS SELECT CONVERT(a USING no_such_cs)", ErrUnknownCharacterSet},
		{"unknown_charset_char_using", "CREATE VIEW _v AS SELECT CHARACTER(a USING no_such_cs)", ErrUnknownCharacterSet},
		{"unknown_charset_cast", "CREATE VIEW _v AS SELECT CAST(a AS CHARACTER CHARSET no_such_cs)", ErrUnknownCharacterSet},

		// SQL_CALC_FOUND_ROWS (1234).
		{"calc_found_rows", "CREATE VIEW _v AS SELECT SQL_CALC_FOUND_ROWS a FROM a", ErrWrongUsage},

		// VALUES clause validation (3942, 3943, 1054).
		{"values_empty_row", "CREATE VIEW _v AS VALUES ROW()", ErrEmptyRowInTVC},
		{"values_leading_comma", "CREATE VIEW _v AS VALUES ROW(, a)", ErrEmptyRowInTVC},
		{"values_default", "CREATE VIEW _v AS VALUES ROW(DEFAULT)", ErrDefaultAsValInTVC},
		{"values_unknown_col", "CREATE VIEW _v AS VALUES ROW(missing)", ErrNoSuchColumn},

		// Expression analysis for previously unsupported AST nodes (1054).
		{"member_of_missing_col", "CREATE VIEW _v AS SELECT missing MEMBER OF(missing)", ErrNoSuchColumn},
		{"match_against_missing_col", "CREATE VIEW _v AS SELECT MATCH(missing) AGAINST('s')", ErrNoSuchColumn},
		{"convert_missing_col", "CREATE VIEW _v AS SELECT CONVERT(missing, BINARY)", ErrNoSuchColumn},
		{"extract_missing_col", "CREATE VIEW _v AS SELECT EXTRACT(DAY FROM missing)", ErrNoSuchColumn},
		{"default_func_missing_col", "CREATE VIEW _v AS SELECT DEFAULT(missing)", ErrNoSuchColumn},
		{"weight_string_missing_col", "CREATE VIEW _v AS SELECT WEIGHT_STRING(missing AS BINARY(1))", ErrNoSuchColumn},
		{"json_table_missing_col", "CREATE VIEW _v AS SELECT a FROM JSON_TABLE(missing, '$' COLUMNS(a FOR ORDINALITY)) jt", ErrNoSuchColumn},
		{"set_var_missing_col_before_1351", "CREATE VIEW _v AS SELECT @v := missing", ErrNoSuchColumn},
		{"off_is_column_not_bool", "CREATE VIEW _v AS SELECT OFF", ErrNoSuchColumn},

		// JSON_QUERY path validation (5454).
		{"json_query_bad_path", "CREATE VIEW _v AS SELECT JSON_QUERY(a, 'bad')", ErrOBJSONPathSyntax},

		// Internal-only constructs rejected with 1149.
		{"cast_ignore", "CREATE VIEW _v AS SELECT CAST(a AS BINARY IGNORE) FROM a", ErrOBSyntaxError},
		{"internal_check_func", "CREATE VIEW _v AS SELECT CHECK(a) FROM a", ErrOBSyntaxError},
		{"internal_st_asmvt", "CREATE VIEW _v AS SELECT _ST_ASMVT(a) FROM a", ErrOBSyntaxError},

		// Variables in view (1351), remote sequence (5835), UNNEST arity (1582).
		{"user_var_in_view", "CREATE VIEW _v AS SELECT @v", ErrViewSelectVariable},
		{"sys_var_in_view", "CREATE VIEW _v AS SELECT @@autocommit", ErrViewSelectVariable},
		{"var_in_values_subquery_view", "CREATE VIEW _v AS VALUES ROW((SELECT @v))", ErrViewSelectVariable},
		{"remote_sequence_in_view", "CREATE VIEW _v AS SELECT a.a@dblink", ErrOBRemoteDBNotFound},
		{"unnest_col_count_mismatch", "CREATE VIEW _v AS SELECT x FROM UNNEST(a) u(x, y)", ErrWrongParamCount},

		// Unsupported MATCH mode & APPROX without vector distance (1235).
		{"match_query_expansion", "CREATE VIEW _v AS SELECT MATCH(s) AGAINST('x' WITH QUERY EXPANSION) FROM t", ErrNotSupportedYet},
		{"approx_without_vector_sort", "CREATE VIEW _v AS SELECT a FROM a ORDER BY a APPROX LIMIT 1", ErrNotSupportedYet},

		// Unknown unqualified & schema-qualified functions (1305).
		{"unknown_func_no_args", "CREATE VIEW _v AS SELECT a()", ErrNoSuchFunction},
		{"unknown_func_statistics_no_args", "CREATE VIEW _v AS SELECT STATISTICS()", ErrNoSuchFunction},
		{"unknown_func_statistics_lit_arg", "CREATE VIEW _v AS SELECT STATISTICS('s')", ErrNoSuchFunction},
		{"unknown_func_coalesce_aggr", "CREATE VIEW _v AS SELECT COALESCE_AGGR()", ErrNoSuchFunction},
		{"unknown_func_replace_const", "CREATE VIEW _v AS SELECT REPLACE_CONST()", ErrNoSuchFunction},
		{"unknown_func_decrypt", "CREATE VIEW _v AS SELECT DECRYPT()", ErrNoSuchFunction},
		{"unknown_func_qualified", "CREATE VIEW _v AS SELECT a.a()", ErrNoSuchFunction},
		{"unknown_func_qualified_statistics", "CREATE VIEW _v AS SELECT a.STATISTICS()", ErrNoSuchFunction},
		{"unknown_func_qualified_builtin", "CREATE VIEW _v AS SELECT a.CURRENT_USER()", ErrNoSuchFunction},
		{"unknown_func_existing_db_missing_fn", "CREATE VIEW _v AS SELECT db2.no_such_fn()", ErrNoSuchFunction},
		{"unknown_func_col_error_precedes_1305", "CREATE VIEW _v AS SELECT STATISTICS(missing_col)", ErrNoSuchColumn},
		{"unknown_func_aliased_arg_precedes_1305", "CREATE VIEW _v AS SELECT STATISTICS(a zzz)", 1583},

		// CALC_PARTITION_ID argument validation precedes column resolution (1210).
		{"calc_partition_id_2_col", "CREATE VIEW _v AS SELECT CALC_PARTITION_ID(a, a)", ErrWrongArguments},
		{"calc_partition_id_3_col", "CREATE VIEW _v AS SELECT CALC_PARTITION_ID(a, a, a)", ErrWrongArguments},
		{"calc_partition_id_1_arg", "CREATE VIEW _v AS SELECT CALC_PARTITION_ID(1)", ErrWrongArguments},
		{"calc_partition_id_const_first_missing_second", "CREATE VIEW _v AS SELECT CALC_PARTITION_ID(1, missing_col)", ErrNoSuchColumn},

		// Optimizer hint validation precedes column resolution (5063, 1210).
		{"hint_load_batch_size_1_arg", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE(1) */ a", ErrOBUnknownHint},
		{"hint_load_batch_size_2_args", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE(1, 's') */ a", ErrOBUnknownHint},
		{"hint_load_batch_size_trailing_word", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE(1) zzz */ a", ErrOBUnknownHint},
		{"hint_load_batch_size_after_other", "CREATE VIEW _v AS SELECT /*+ PARALLEL(2) LOAD_BATCH_SIZE(1) */ 1", ErrOBUnknownHint},
		{"hint_load_batch_size_subquery", "CREATE VIEW _v AS SELECT (SELECT /*+ LOAD_BATCH_SIZE(1) */ 1)", ErrOBUnknownHint},
		{"hint_direct_bad_load_method", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1, 's') */ a", ErrWrongArguments},
		{"hint_direct_bad_load_method_trailing", "CREATE VIEW _v AS SELECT /*+ DIRECT(FALSE, 0, 'bogus') zzz */ a", ErrWrongArguments},
		// Malformed hint syntax is discarded by OceanBase and falls through to column resolution (1054).
		{"hint_load_batch_size_malformed_empty", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE() */ a", ErrNoSuchColumn},
		{"hint_load_batch_size_malformed_missing_lparen", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE 1) */ a", ErrNoSuchColumn},
		{"hint_load_batch_size_malformed_missing_rparen", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE(1 */ a", ErrNoSuchColumn},
		{"hint_load_batch_size_malformed_trailing_paren", "CREATE VIEW _v AS SELECT /*+ LOAD_BATCH_SIZE(1) (bad) */ a", ErrNoSuchColumn},
		{"hint_direct_malformed_missing_comma", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1 's') */ a", ErrNoSuchColumn},
		{"hint_direct_2_args_valid_falls_through", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1) */ a", ErrNoSuchColumn},
		{"hint_direct_valid_full_falls_through", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1, 'full') */ a", ErrNoSuchColumn},
		{"hint_regular_comment_ignored", "CREATE VIEW _v AS SELECT /* LOAD_BATCH_SIZE(1) */ a", ErrNoSuchColumn},
	}

	for _, tc := range refusals {
		t.Run("refuse/"+tc.name, func(t *testing.T) {
			c := setup(t)
			if got := execCode(t, c, tc.sql); got != tc.want {
				t.Errorf("%s: got error %d, want %d", tc.sql, got, tc.want)
			}
		})
	}

	accepted := []struct {
		name    string
		sql     string
		wantDef string
	}{
		{"defined_window", "CREATE VIEW _v AS SELECT RANK() OVER w FROM a WINDOW w AS (ORDER BY b)", "rank() OVER `w`"},
		{"ranking_drops_frame", "CREATE VIEW _v AS SELECT RANK() OVER (ROWS CURRENT ROW)", "select rank() OVER ( )  AS `RANK() OVER (ROWS CURRENT ROW)`"},
		{"charset_introducer_stripped", "CREATE VIEW _v AS SELECT _utf8mb4 'hello'", "select 'hello' AS `hello`"},
		{"binary_introducer_hex", "CREATE VIEW _v AS SELECT _binary 's'", "select 0x73 AS `s`"},
		{"national_string_alias", "CREATE VIEW _v AS SELECT N's'", "select 's' AS `N's'`"},
		{"geomcollection_rewrite", "CREATE VIEW _v AS SELECT GEOMETRYCOLLECTION()", "select geomcollection() AS `GEOMETRYCOLLECTION()`"},
		{"key_version_rewrite", "CREATE VIEW _v AS SELECT KEY_VERSION()", "select version() AS `KEY_VERSION()`"},
		{"system_user_rewrite", "CREATE VIEW _v AS SELECT SYSTEM_USER()", "select user() AS `SYSTEM_USER()`"},
		{"session_user_rewrite", "CREATE VIEW _v AS SELECT SESSION_USER()", "select user() AS `SESSION_USER()`"},
		{"now_fsp_dropped_in_view", "CREATE VIEW _v AS SELECT NOW(1), SYSDATE(1), CURTIME(1), CURRENT_TIMESTAMP(1)", "select now() AS `NOW(1)`,sysdate() AS `SYSDATE(1)`,curtime() AS `CURTIME(1)`,now() AS `CURRENT_TIMESTAMP(1)`"},
		{"approx_with_vector_sort", "CREATE VIEW _v AS SELECT id FROM t ORDER BY L2_DISTANCE(v, '[1,2,3]') APPROX LIMIT 1", "approx limit 1"},
		{"approx_with_vec_prefix_sort", "CREATE VIEW _v AS SELECT id FROM t ORDER BY VEC_COSINE_DISTANCE(v, '[1,2,3]') APPROX LIMIT 1", "approx limit 1"},
		{"valid_json_query_path", "CREATE VIEW _v AS SELECT JSON_QUERY(s, '$.a') FROM t", "json_query("},
		{"calc_partition_id_const_first", "CREATE VIEW _v AS SELECT CALC_PARTITION_ID(1, a) FROM a", "calc_partition_id("},
		{"hint_direct_2_args_accepted", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1) */ a FROM a", "select `a`.`a` AS `a` from `a`"},
		{"hint_direct_full_accepted", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1, 'full') */ a FROM a", "select `a`.`a` AS `a` from `a`"},
		{"hint_direct_inc_accepted", "CREATE VIEW _v AS SELECT /*+ DIRECT(TRUE, 1, 'inc') */ a FROM a", "select `a`.`a` AS `a` from `a`"},
		{"hint_direct_inc_replace_accepted", "CREATE VIEW _v AS SELECT /*+ DIRECT(FALSE, 0, 'INC_REPLACE') */ a FROM a", "select `a`.`a` AS `a` from `a`"},
		{"hint_regular_comment_accepted", "CREATE VIEW _v AS SELECT /* LOAD_BATCH_SIZE(1) */ a FROM a", "select `a`.`a` AS `a` from `a`"},
	}

	for _, tc := range accepted {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			c := setup(t)
			if got := execCode(t, c, tc.sql); got != 0 {
				t.Fatalf("%s: unexpected error %d", tc.sql, got)
			}
			v := c.GetDatabase("db1").Views["_v"]
			if v == nil {
				t.Fatalf("%s: view _v not created", tc.sql)
			}
			if !strings.Contains(v.Definition, tc.wantDef) {
				t.Errorf("%s: definition %q does not contain %q", tc.sql, v.Definition, tc.wantDef)
			}
		})
	}

	// Variables are allowed in CREATE TABLE ... SELECT (1351 is view-only).
	t.Run("accept/create_table_select_variable", func(t *testing.T) {
		c := setup(t)
		if got := execCode(t, c, "CREATE TABLE t_var AS SELECT @v"); got != 0 {
			t.Fatalf("CREATE TABLE t_var AS SELECT @v: unexpected error %d", got)
		}
	})

	// User-defined stored functions resolve in the current database and via
	// explicit schema qualification, while missing routines fail with 1305.
	t.Run("stored_function_resolution_in_views", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE FUNCTION my_fn(x INT) RETURNS INT DETERMINISTIC RETURN x + 1;")
		mustExec(t, c, "CREATE FUNCTION db2.other_fn() RETURNS INT DETERMINISTIC RETURN 42;")
		mustExec(t, c, "CREATE FUNCTION db2.current_user() RETURNS VARCHAR(32) DETERMINISTIC RETURN 'u';")
		mustExec(t, c, "CREATE FUNCTION db2.calc_partition_id(x INT, y INT) RETURNS INT DETERMINISTIC RETURN x + y;")

		if got := execCode(t, c, "CREATE VIEW v_sf1 AS SELECT my_fn(a) FROM a"); got != 0 {
			t.Fatalf("unqualified stored function: unexpected error %d", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_sf2 AS SELECT db1.my_fn(a) FROM a"); got != 0 {
			t.Fatalf("same-db qualified stored function: unexpected error %d", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_sf3 AS SELECT db2.other_fn()"); got != 0 {
			t.Fatalf("cross-db qualified stored function: unexpected error %d", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_sf4 AS SELECT db2.CURRENT_USER()"); got != 0 {
			t.Fatalf("qualified stored function shadowing builtin: unexpected error %d", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_sf5 AS SELECT db2.CALC_PARTITION_ID(a, a) FROM a"); got != 0 {
			t.Fatalf("qualified stored function shadowing CALC_PARTITION_ID: unexpected error %d", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_sf_miss1 AS SELECT other_fn()"); got != ErrNoSuchFunction {
			t.Fatalf("unqualified call to other db's function: got %d, want %d", got, ErrNoSuchFunction)
		}
		if got := execCode(t, c, "CREATE VIEW v_sf_miss2 AS SELECT db2.my_fn(1)"); got != ErrNoSuchFunction {
			t.Fatalf("wrong schema for stored function: got %d, want %d", got, ErrNoSuchFunction)
		}
	})

	// Cross-database unaliased tables with the same table name do not collide,
	// whereas same-database unaliased tables and matching explicit aliases do (1066).
	t.Run("cross_database_from_uniqueness", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE TABLE db2.t (id INT PRIMARY KEY);")
		if got := execCode(t, c, "CREATE VIEW v_xtab AS SELECT db1.t.id FROM db1.t JOIN db2.t ON db1.t.id = db2.t.id"); got != 0 {
			t.Fatalf("cross-db unaliased join: unexpected error %d", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_same_db AS SELECT db1.t.id FROM db1.t JOIN t ON db1.t.id = t.id"); got != 1066 {
			t.Fatalf("same-db unaliased join: got %d, want 1066", got)
		}
		if got := execCode(t, c, "CREATE VIEW v_same_alias AS SELECT x.id FROM db1.t x JOIN db2.t x ON x.id = x.id"); got != 1066 {
			t.Fatalf("cross-db duplicate alias: got %d, want 1066", got)
		}
	})

	// Hex and bit literals preserve their encoding through AnalyzeSelectStmt + DeparseQuery.
	t.Run("analyzed_hex_and_bit_literals", func(t *testing.T) {
		c := setup(t)
		list, err := parser.Parse("SELECT X'41', 0x42, B'101', 0b110")
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		q, err := c.AnalyzeSelectStmt(list.Items[0].(*nodes.SelectStmt))
		if err != nil {
			t.Fatalf("analyze: %v", err)
		}
		got := DeparseQuery(q)
		for _, want := range []string{"X'41'", "0x42", "b'101'", "0b110"} {
			if !strings.Contains(got, want) {
				t.Errorf("DeparseQuery = %q, want substring %q", got, want)
			}
		}
	})

	// EXTRACT over nullable and non-null columns propagates nullability through view-on-view metadata export.
	t.Run("extract_view_on_view_nullability", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE TABLE dt_tbl (id INT PRIMARY KEY, dt DATETIME NULL, dt_nn DATETIME NOT NULL);")
		mustExec(t, c, "CREATE VIEW v_ext AS SELECT EXTRACT(DAY FROM dt) AS d_null, EXTRACT(DAY FROM dt_nn) AS d_nn FROM dt_tbl;")
		mustExec(t, c, "CREATE VIEW w_ext AS SELECT d_null, d_nn FROM v_ext;")
		md := c.DatabaseMetadata("db1")
		if md == nil {
			t.Fatal("DatabaseMetadata(db1) returned nil")
		}
		var found bool
		for _, sc := range md.Schemas {
			for _, v := range sc.Views {
				if v.Name == "w_ext" {
					found = true
					if len(v.Columns) != 2 {
						t.Fatalf("w_ext columns = %d, want 2", len(v.Columns))
					}
					if !v.Columns[0].Nullable {
						t.Errorf("w_ext.d_null Nullable = false, want true")
					}
					if v.Columns[1].Nullable {
						t.Errorf("w_ext.d_nn Nullable = true, want false")
					}
				}
			}
		}
		if !found {
			t.Fatal("view w_ext not found in DatabaseMetadata")
		}
	})
}
