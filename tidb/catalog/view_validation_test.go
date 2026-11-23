package catalog

import (
	"errors"
	"strings"
	"testing"
)

func TestViewValidation(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		for _, stmt := range []string{
			"CREATE DATABASE db",
			"USE db",
			"CREATE TABLE t (a INT PRIMARY KEY, b INT, s VARCHAR(10))",
			"CREATE SEQUENCE seq",
		} {
			res, err := c.Exec(stmt, nil)
			if err != nil {
				t.Fatalf("setup %q parse err: %v", stmt, err)
			}
			for _, r := range res {
				if r.Error != nil {
					t.Fatalf("setup %q exec err: %v", stmt, r.Error)
				}
			}
		}
		return c
	}

	assertCode := func(t *testing.T, c *Catalog, sql string, wantCode int) {
		t.Helper()
		res, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("Exec(%q) returned parse err: %v, want code %d", sql, err, wantCode)
		}
		if len(res) != 1 || res[0].Error == nil {
			t.Fatalf("Exec(%q) succeeded, want error code %d", sql, wantCode)
		}
		var catErr *Error
		if !errors.As(res[0].Error, &catErr) {
			t.Fatalf("Exec(%q) err = %T (%v), want *Error with code %d", sql, res[0].Error, res[0].Error, wantCode)
		}
		if catErr.Code != wantCode {
			t.Errorf("Exec(%q) code = %d (%s), want %d", sql, catErr.Code, catErr.Message, wantCode)
		}
	}

	assertOK := func(t *testing.T, c *Catalog, sql string) {
		t.Helper()
		res, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("Exec(%q) parse err: %v", sql, err)
		}
		for _, r := range res {
			if r.Error != nil {
				t.Fatalf("Exec(%q) unexpected err: %v", sql, r.Error)
			}
		}
	}

	t.Run("select_into_clause_1350", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT * INTO OUTFILE 'x'", ErrViewSelectClause)
		assertCode(t, c, "CREATE VIEW v AS SELECT a FROM t INTO OUTFILE 'x'", ErrViewSelectClause)
		assertCode(t, c, "CREATE VIEW v AS SELECT a FROM t INTO OUTFILE 'x' FIELDS TERMINATED BY ','", ErrViewSelectClause)
	})

	t.Run("variables_1351", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT @x FROM t", ErrViewSelectVariable)
		assertCode(t, c, "CREATE VIEW v AS SELECT @@sql_mode FROM t", ErrViewSelectVariable)
		assertCode(t, c, "CREATE VIEW v AS SELECT @x := a FROM t", ErrViewSelectVariable)
		assertCode(t, c, "CREATE VIEW v AS SELECT a >= @x := a FROM t", ErrViewSelectVariable)
	})

	t.Run("param_marker_1149", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT ? FROM t", ErrSyntax)
		assertCode(t, c, "CREATE VIEW v AS SELECT a FROM t LIMIT ?", ErrSyntax)
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS ? PRECEDING) FROM t", ErrSyntax)
	})

	t.Run("duplicate_cte_1066", func(t *testing.T) {
		c := setup(t)
		// Duplicate CTE fires before star expansion (1051).
		assertCode(t, c, "CREATE VIEW v AS WITH c AS (SELECT *), c AS (SELECT *) SELECT *", ErrNonuniqTable)
		assertCode(t, c, "CREATE VIEW v AS WITH c AS (SELECT 1), C AS (SELECT 2) SELECT * FROM c", ErrNonuniqTable)
		assertOK(t, c, "CREATE VIEW v AS WITH c1 AS (SELECT 1 AS x), c2 AS (SELECT x FROM c1) SELECT x FROM c2")
	})

	t.Run("row_target_1241", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW(a, b) FROM t", ErrOperandColumns)
		assertCode(t, c, "CREATE VIEW v AS SELECT (ROW(a, b)) FROM t", ErrOperandColumns)
		// ROW in a comparison predicate is valid.
		assertOK(t, c, "CREATE VIEW v AS SELECT ROW(a, b) = ROW(1, 2) FROM t")
	})

	t.Run("sequence_target_1347", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT NEXT VALUE FOR t FROM t", ErrWrongObject)
		assertCode(t, c, "CREATE VIEW v AS SELECT NEXTVAL(t) FROM t", ErrWrongObject)
		assertCode(t, c, "CREATE VIEW v AS SELECT LASTVAL(t) FROM t", ErrWrongObject)
		assertCode(t, c, "CREATE VIEW v AS SELECT SETVAL(t, 1) FROM t", ErrWrongObject)
		// Real sequence succeeds and does not fail unknown-column checking.
		assertOK(t, c, "CREATE VIEW v AS SELECT NEXT VALUE FOR seq FROM t")
		assertOK(t, c, "CREATE VIEW v2 AS SELECT NEXTVAL(seq), LASTVAL(seq), SETVAL(seq, 10) FROM t")
	})

	t.Run("schema_qualified_function_1305", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT db.foo() FROM t", ErrNoSuchFunction)
	})

	t.Run("builtin_arity_1582", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT ASCII() FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT COALESCE() FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT IF() FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT IF(a, b) FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT LEFT(s) FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT REPLACE(s, 'a') FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT INSERT() FROM t", ErrWrongParamCount)
		assertCode(t, c, "CREATE VIEW v AS SELECT PASSWORD() FROM t", ErrWrongParamCount)
		// Valid arities succeed.
		assertOK(t, c, "CREATE VIEW v AS SELECT ASCII(s), COALESCE(a, b), IF(a > 0, a, b), LEFT(s, 2), LOG(a), LOG(2, a), FORMAT(a, 2), FORMAT(a, 2, 'en_US'), WEEK(s), WEEK(s, 1) FROM t")
	})

	t.Run("window_validation", func(t *testing.T) {
		c := setup(t)
		// 1210 fires before 3579 for NTILE / NTH_VALUE with column arg.
		assertCode(t, c, "CREATE VIEW v AS SELECT NTILE(a) OVER missing_w FROM t", ErrWrongArguments)
		assertCode(t, c, "CREATE VIEW v AS SELECT NTH_VALUE(a, b) OVER missing_w FROM t", ErrWrongArguments)
		assertCode(t, c, "CREATE VIEW v AS SELECT NTILE(0) OVER () FROM t", ErrWrongArguments)
		assertCode(t, c, "CREATE VIEW v AS SELECT NTILE(-1) OVER () FROM t", ErrWrongArguments)
		assertCode(t, c, "CREATE VIEW v AS SELECT NTILE(-(-(-2))) OVER () FROM t", ErrWrongArguments)
		assertCode(t, c, "CREATE VIEW v AS SELECT NTH_VALUE(a, -1) OVER () FROM t", ErrWrongArguments)
		// 3579 undefined window.
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER missing_w FROM t", ErrWindowNoSuchWindow)
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (missing_w) FROM t", ErrWindowNoSuchWindow)
		// 3591 duplicate window name.
		assertCode(t, c, "CREATE VIEW v AS SELECT * FROM t WINDOW w AS (), W AS ()", ErrWindowDuplicateName)
		// 3584 frame start UNBOUNDED FOLLOWING (fires before 3585).
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN UNBOUNDED FOLLOWING AND UNBOUNDED PRECEDING) FROM t", ErrWindowFrameStartIllegal)
		// 3585 frame end UNBOUNDED PRECEDING (fires before 3596).
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED PRECEDING) FROM t", ErrWindowFrameEndIllegal)
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN INTERVAL a MICROSECOND FOLLOWING AND UNBOUNDED PRECEDING) FROM t", ErrWindowFrameEndIllegal)
		// 3596 ROWS with INTERVAL frame bound.
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ROWS INTERVAL a MICROSECOND PRECEDING) FROM t", ErrWindowRowsIntervalUse)
		// Valid window forms succeed, including double-negated positive constant args.
		assertOK(t, c, "CREATE VIEW v1 AS SELECT ROW_NUMBER() OVER w, NTILE(4) OVER w, NTH_VALUE(a, 2) OVER w FROM t WINDOW w AS (ORDER BY a ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)")
		assertOK(t, c, "CREATE VIEW v2 AS SELECT ROW_NUMBER() OVER (ORDER BY a RANGE INTERVAL 1 DAY PRECEDING) FROM t")
		assertOK(t, c, "CREATE VIEW v3 AS SELECT NTILE(-(-2)) OVER () AS bucket, NTH_VALUE(a, -(-2)) OVER () AS second_val FROM t")
		assertOK(t, c, "CREATE VIEW v4 AS SELECT NTILE(-(1 - 3)) OVER () AS bucket FROM t")
	})

	t.Run("unknown_columns_1054", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT a", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT a AS x", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT missing FROM t", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT ASCII(missing) FROM t", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT TRIM(BOTH missing FROM s) FROM t", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT other_db.t.a FROM t", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT CASE missing WHEN a THEN b END FROM t", ErrNoSuchColumn)
		// Valid column refs, correlated subqueries, CTEs, derived tables, and aliases succeed.
		assertOK(t, c, "CREATE VIEW v1 AS SELECT db.t.a, t.b, s FROM t WHERE a > 0")
		assertOK(t, c, "CREATE VIEW v2 AS SELECT a, (SELECT b FROM t t2 WHERE t2.a = t.a) AS sub FROM t")
		assertOK(t, c, "CREATE VIEW v3 AS SELECT a + 1 AS total FROM t ORDER BY total")
		assertOK(t, c, "CREATE VIEW v4 AS WITH cte AS (SELECT a AS x FROM t) SELECT x FROM cte")
		assertOK(t, c, "CREATE VIEW v5 AS SELECT d.x FROM (SELECT a AS x FROM t) d")

		// TiDB hidden _tidb_rowid handle: available on tables without a
		// clustered PK (no PK or NONCLUSTERED PK), unavailable on clustered PK
		// tables and views that do not project it.
		assertOK(t, c, "CREATE TABLE t_no_pk (x INT)")
		assertOK(t, c, "CREATE TABLE t_nonclustered (x INT PRIMARY KEY NONCLUSTERED)")
		assertOK(t, c, "CREATE VIEW v_rowid_no_pk AS SELECT _tidb_rowid AS rid FROM t_no_pk")
		assertOK(t, c, "CREATE VIEW v_rowid_nonclustered AS SELECT t_nonclustered._tidb_rowid AS rid FROM t_nonclustered")
		assertCode(t, c, "CREATE VIEW v_rowid_clustered AS SELECT _tidb_rowid AS rid FROM t", ErrNoSuchColumn)
		assertOK(t, c, "CREATE VIEW v_plain AS SELECT x FROM t_no_pk")
		assertCode(t, c, "CREATE VIEW v_rowid_from_view AS SELECT _tidb_rowid AS rid FROM v_plain", ErrNoSuchColumn)
	})

	t.Run("optional_braces_functions_restore", func(t *testing.T) {
		c := setup(t)
		assertOK(t, c, "CREATE VIEW v_role AS SELECT CURRENT_ROLE FROM t")
		if got := c.ShowCreateView("db", "v_role"); !strings.Contains(got, "SELECT CURRENT_ROLE() AS `CURRENT_ROLE` FROM `db`.`t`") {
			t.Errorf("ShowCreateView(v_role) = %q", got)
		}
		assertOK(t, c, "CREATE VIEW v_tso AS SELECT TIDB_CURRENT_TSO FROM t")
		if got := c.ShowCreateView("db", "v_tso"); !strings.Contains(got, "SELECT TIDB_CURRENT_TSO() AS `TIDB_CURRENT_TSO` FROM `db`.`t`") {
			t.Errorf("ShowCreateView(v_tso) = %q", got)
		}
		// Backtick-quoted name stays a column ref, not a niladic function call.
		assertOK(t, c, "CREATE TABLE t_quoted (`CURRENT_ROLE` INT, `TIDB_CURRENT_TSO` INT)")
		assertOK(t, c, "CREATE VIEW v_quoted AS SELECT `CURRENT_ROLE`, `TIDB_CURRENT_TSO` FROM t_quoted")
		if got := c.ShowCreateView("db", "v_quoted"); !strings.Contains(got, "SELECT `CURRENT_ROLE` AS `CURRENT_ROLE`,`TIDB_CURRENT_TSO` AS `TIDB_CURRENT_TSO` FROM `db`.`t_quoted`") {
			t.Errorf("ShowCreateView(v_quoted) = %q", got)
		}
	})

	t.Run("values_primary_1051", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS VALUES ROW()", ErrUnknownTable)
		assertCode(t, c, "CREATE VIEW v AS VALUES ROW(1)", ErrUnknownTable)
		assertCode(t, c, "CREATE VIEW v AS VALUES ROW(a, b)", ErrUnknownTable)
		assertCode(t, c, "CREATE VIEW v AS VALUES ROW(1), ROW(2)", ErrUnknownTable)
	})

	t.Run("noop_functions_1235", func(t *testing.T) {
		c := setup(t)
		// Fires before star expansion (1051).
		assertCode(t, c, "CREATE VIEW v AS SELECT * GROUP BY a ASC", ErrNotSupportedYet)
		assertCode(t, c, "CREATE VIEW v AS SELECT * GROUP BY a DESC", ErrNotSupportedYet)
		assertCode(t, c, "CREATE VIEW v AS SELECT SQL_CALC_FOUND_ROWS *", ErrNotSupportedYet)
		// Plain GROUP BY without ASC/DESC succeeds.
		assertOK(t, c, "CREATE VIEW v_gb AS SELECT a FROM t GROUP BY a")
	})

	t.Run("match_against_1105", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT MATCH(s) AGAINST('foo') FROM t", ErrUnknownError)
		assertCode(t, c, "CREATE VIEW v AS SELECT MATCH(s) AGAINST('foo' IN BOOLEAN MODE) FROM t", ErrUnknownError)
		assertCode(t, c, "CREATE VIEW v AS SELECT MATCH(s) AGAINST('foo' IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION) FROM t", ErrUnknownError)
	})

	t.Run("cast_array_and_groups_frame_1235", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT CAST(a AS BINARY ARRAY) FROM t", ErrNotSupportedYet)
		assertCode(t, c, "CREATE VIEW v AS SELECT ROW_NUMBER() OVER (GROUPS UNBOUNDED PRECEDING) FROM t", ErrNotSupportedYet)
		assertOK(t, c, "CREATE VIEW v_cast AS SELECT CAST(a AS BINARY) FROM t")
	})

	t.Run("unknown_unqualified_function_1305_and_approx_percentile_1105", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT unknown_fn() FROM t", ErrNoSuchFunction)
		assertCode(t, c, "CREATE VIEW v AS SELECT unknown_fn(a) FROM t", ErrNoSuchFunction)
		assertCode(t, c, "CREATE VIEW v AS SELECT POINT() FROM t", ErrNoSuchFunction)
		assertCode(t, c, "CREATE VIEW v AS SELECT APPROX_PERCENTILE(a) FROM t", ErrUnknownError)
		assertOK(t, c, "CREATE VIEW v_ok AS SELECT APPROX_PERCENTILE(a, 50), ABS(a), TIDB_VERSION(), TIDB_ENCODE_SQL_DIGEST('select 1'), TIDB_DECODE_BASE64_KEY(''), TIDB_DECODE_BINARY_PLAN(''), TIDB_ENCODE_RECORD_KEY('db', 't', 1), TIDB_ENCODE_INDEX_KEY('db', 't', 'idx', 1) FROM t")
	})

	t.Run("temporal_literals_1292", func(t *testing.T) {
		c := setup(t)
		assertCode(t, c, "CREATE VIEW v AS SELECT DATE 'binary' FROM t", ErrTruncatedWrongValue)
		assertCode(t, c, "CREATE VIEW v AS SELECT TIME 'binary' FROM t", ErrTruncatedWrongValue)
		assertCode(t, c, "CREATE VIEW v AS SELECT TIMESTAMP 'binary' FROM t", ErrTruncatedWrongValue)
		assertCode(t, c, "CREATE VIEW v AS SELECT DATE '2024-02-30' FROM t", ErrTruncatedWrongValue)
		assertCode(t, c, "CREATE VIEW v AS SELECT TIME '900:00:00' FROM t", ErrTruncatedWrongValue)
		assertOK(t, c, "CREATE VIEW v_tmp AS SELECT DATE '2024-02-29', TIME '-12:34:56.789', TIMESTAMP '2024-01-01 12:00:00', TIMESTAMP '2024-01-01' FROM t")
	})

	t.Run("json_operator_types_3146", func(t *testing.T) {
		c := setup(t)
		assertOK(t, c, "CREATE TABLE tj (id INT PRIMARY KEY, j JSON, s VARCHAR(100))")
		assertCode(t, c, "CREATE VIEW v AS SELECT id -> '$.a' FROM tj", ErrInvalidTypeForJSON)
		assertCode(t, c, "CREATE VIEW v AS SELECT id ->> '$.a' FROM tj", ErrInvalidTypeForJSON)
		assertCode(t, c, "CREATE VIEW v AS SELECT id MEMBER OF(id) FROM tj", ErrInvalidTypeForJSON)
		// Unknown column still fails with 1054 before 3146.
		assertCode(t, c, "CREATE VIEW v AS SELECT missing -> '$.a' FROM tj", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT missing MEMBER OF(id) FROM tj", ErrNoSuchColumn)
		assertCode(t, c, "CREATE VIEW v AS SELECT id MEMBER OF(missing) FROM tj", ErrNoSuchColumn)
		// JSON and string columns succeed.
		assertOK(t, c, "CREATE VIEW v_json AS SELECT j -> '$.a' AS c1, s ->> '$.a' AS c2, id MEMBER OF(j) AS c3, id MEMBER OF('[1,2]') AS c4 FROM tj")
		// An opaque inner scope (derived table or CTE) may shadow an outer
		// non-JSON column with a JSON/string column; do not attribute the outer
		// column's type across the opaque boundary.
		assertOK(t, c, "CREATE TABLE t_outer (j INT)")
		assertOK(t, c, "CREATE VIEW v_shadow_derived AS SELECT (SELECT j -> '$.a' FROM (SELECT '{\"a\":1}' AS j) d) AS x FROM t_outer")
		assertOK(t, c, "CREATE VIEW v_shadow_cte AS WITH cte AS (SELECT '{\"a\":1}' AS j) SELECT j -> '$.a' AS x FROM cte")
		// A correlated reference through a non-opaque inner scope still checks
		// the outer column's type.
		assertCode(t, c, "CREATE VIEW v_corr_int AS SELECT (SELECT t_outer.j -> '$.a' FROM tj) AS x FROM t_outer", ErrInvalidTypeForJSON)
	})
}
