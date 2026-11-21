package catalog

import (
	"testing"
)

func TestViewQueryAndExpressionValidation(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE db1; CREATE DATABASE db2; USE db1;")
		mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT, s VARCHAR(20));")
		mustExec(t, c, "CREATE TABLE db2.t (id INT PRIMARY KEY, b INT);")
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

		// SELECT ... INTO is refused with 1350 ahead of 1054.
		{"select_into_var", "CREATE VIEW v AS SELECT missing INTO @x", ErrViewSelectClause},
		{"select_into_outfile", "CREATE VIEW v AS SELECT missing INTO OUTFILE 'f'", ErrViewSelectClause},

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
		{"dual_literal", "CREATE VIEW v AS SELECT 1 FROM DUAL"},
	}
	for _, tc := range accepted {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			c := setup(t)
			mustExec(t, c, tc.sql)
		})
	}
}
