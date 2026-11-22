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
