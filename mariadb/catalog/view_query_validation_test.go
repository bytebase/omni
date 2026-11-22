package catalog

import (
	"errors"
	"strings"
	"testing"
)

func TestViewQueryValidation(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		code int // 0 = expect success
	}{
		// SELECT ... INTO is a parse error (1064) in CREATE/ALTER VIEW.
		{"into_var", "CREATE VIEW _v AS SELECT a INTO @x FROM a", 1064},
		{"into_outfile", "CREATE VIEW _v AS SELECT a INTO OUTFILE 's' FROM a", 1064},

		// [NO]WAIT on locking clause is 1350 before column resolution.
		{"for_update_wait", "CREATE VIEW _v AS SELECT a FOR UPDATE WAIT 1", 1350},
		{"for_update_nowait", "CREATE VIEW _v AS SELECT a FOR UPDATE NOWAIT", 1350},
		{"for_update_plain_no_from", "CREATE VIEW _v AS SELECT a FOR UPDATE", 1054},

		// Parameter marker in DDL is 1064.
		{"param_marker_select", "CREATE VIEW _v AS SELECT ? FROM a", 1064},
		{"param_marker_limit", "CREATE VIEW _v AS SELECT a FROM a LIMIT ?", 1064},

		// Star and table resolution.
		{"bare_star_no_from", "CREATE VIEW _v AS SELECT *", 1096},
		{"bare_star_dual", "CREATE VIEW _v AS SELECT * FROM DUAL", 1096},
		{"qual_star_missing", "CREATE VIEW _v AS SELECT a.*", 1051},
		{"schema_qual_star_missing", "CREATE VIEW _v AS SELECT a.a.*", 1051},
		{"missing_table", "CREATE VIEW _v AS SELECT a FROM nonexistent", 1146},
		{"dual_unknown_col", "CREATE VIEW _v AS SELECT a FROM DUAL", 1054},
		{"dual_literal_ok", "CREATE VIEW _v AS SELECT 1 FROM DUAL", 0},
		{"wrong_schema_3part_col", "CREATE VIEW _v AS SELECT a.a.a FROM a", 1054},
		{"unknown_col_in_like", "CREATE VIEW _v AS SELECT nonexistent LIKE a FROM a", 1054},

		// Duplicate column names and multi-column ROW in target list.
		{"dup_col", "CREATE VIEW _v AS SELECT a, a FROM a", 1060},
		{"row_in_target", "CREATE VIEW _v AS SELECT ROW(a, a) FROM a", 1241},

		// Window frame disallowed on ranking/value window functions (4017).
		{"row_number_frame", "CREATE VIEW _v AS SELECT ROW_NUMBER() OVER (ROWS UNBOUNDED PRECEDING) FROM a", 4017},
		{"rank_named_window_frame", "CREATE VIEW _v AS SELECT RANK() OVER w FROM a WINDOW w AS (ROWS UNBOUNDED PRECEDING)", 4017},
		{"sum_frame_ok", "CREATE VIEW _v AS SELECT SUM(a) OVER (ROWS UNBOUNDED PRECEDING) FROM a", 0},

		// Sequence expressions against regular table (4089).
		{"next_value_for_table", "CREATE VIEW _v AS SELECT NEXT VALUE FOR a FROM a", 4089},
		{"prev_value_for_table", "CREATE VIEW _v AS SELECT PREVIOUS VALUE FOR a FROM a", 4089},
		{"nextval_table", "CREATE VIEW _v AS SELECT NEXTVAL(a) FROM a", 4089},
		{"lastval_table", "CREATE VIEW _v AS SELECT LASTVAL(a) FROM a", 4089},
		{"setval_table", "CREATE VIEW _v AS SELECT SETVAL(a, 1) FROM a", 4089},

		// Non-constant ESCAPE / AGAINST (1210).
		{"like_escape_col", "CREATE VIEW _v AS SELECT a LIKE a ESCAPE a FROM a", 1210},
		{"like_escape_lit_ok", "CREATE VIEW _v AS SELECT a LIKE a ESCAPE '\\\\' FROM a", 0},
		{"match_against_col", "CREATE VIEW _v AS SELECT MATCH(a) AGAINST(a) FROM a", 1210},

		// Schema-qualified special-syntax TRIM/SUBSTRING outside builtin schemas (1128).
		{"qual_substring_from", "CREATE VIEW _v AS SELECT a.SUBSTRING(a FROM a) FROM a", 1128},
		{"qual_trim_leading", "CREATE VIEW _v AS SELECT a.TRIM(LEADING FROM a) FROM a", 1128},

		// Hex literal precision on CURRENT_TIME (1426).
		{"curtime_hex_precision", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x41) FROM a", 1426},
		{"curtime_hex_valid", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x06) FROM a", 0},

		// Qualified DEFAULT(tbl.col) round-trip.
		{"default_col_ok", "CREATE VIEW _v AS SELECT DEFAULT(a) FROM a", 0},
		{"default_qual_col_ok", "CREATE VIEW _v AS SELECT DEFAULT(a.a) FROM a", 0},

		// Enclosing CTE visible inside scalar, EXISTS, IN, quantified, and derived subqueries.
		{"cte_in_scalar_subquery", "CREATE VIEW _v AS WITH c AS (SELECT 1 AS x) SELECT (SELECT x FROM c) AS x", 0},
		{"cte_in_exists_subquery", "CREATE VIEW _v AS WITH c AS (SELECT 1 AS x) SELECT EXISTS(SELECT x FROM c) AS e", 0},
		{"cte_in_in_subquery", "CREATE VIEW _v AS WITH c AS (SELECT 1 AS x) SELECT a IN (SELECT x FROM c) AS m FROM a", 0},
		{"cte_in_derived_subquery", "CREATE VIEW _v AS WITH c AS (SELECT 1 AS x) SELECT * FROM (SELECT x FROM c) d", 0},

		// Huge WEIGHT_STRING LEVEL range must be clamped and complete instantly.
		{"weight_string_huge_level_range", "CREATE VIEW _v AS SELECT WEIGHT_STRING('a' LEVEL 1 - 4294967295) AS w", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := LoadSQL("CREATE DATABASE test; USE test; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT);")
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			results, parseErr := c.Exec(tc.sql, nil)
			var gotCode int
			if parseErr != nil {
				gotCode = 1064
			} else if len(results) > 0 && results[0].Error != nil {
				var ce *Error
				if errors.As(results[0].Error, &ce) {
					gotCode = ce.Code
				} else {
					t.Fatalf("unexpected non-catalog error: %v", results[0].Error)
				}
			}
			if gotCode != tc.code {
				t.Errorf("Exec(%q): got code %d (parseErr=%v, results=%+v), want %d", tc.sql, gotCode, parseErr, results, tc.code)
			}
		})
	}
}

func TestViewDeparseCanonicalization(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "correlated_exists",
			sql:  "CREATE VIEW _v AS SELECT EXISTS(SELECT a) FROM a",
			want: "select exists(select `a`.`a` limit 1) AS `EXISTS(SELECT a)` from `a`",
		},
		{
			name: "correlated_not_in_subquery",
			sql:  "CREATE VIEW _v AS SELECT a NOT IN (SELECT a) FROM a",
			want: "select !(`a`.`a` in (select `a`.`a`)) AS `a NOT IN (SELECT a)` from `a`",
		},
		{
			name: "single_in_and_not_in",
			sql:  "CREATE VIEW _v AS SELECT a IN (a) AS c1, a NOT IN (a) AS c2 FROM a",
			want: "select `a`.`a` = `a`.`a` AS `c1`,`a`.`a` <> `a`.`a` AS `c2` from `a`",
		},
		{
			name: "temporal_and_math_rewrites",
			sql:  "CREATE VIEW _v AS SELECT DATE(a) AS c1, TIME(a) AS c2, TIMESTAMP(a) AS c3, DAY(a) AS c4, MOD(a, b) AS c5 FROM a",
			want: "select cast(`a`.`a` as date) AS `c1`,cast(`a`.`a` as time) AS `c2`,cast(`a`.`a` as datetime) AS `c3`,dayofmonth(`a`.`a`) AS `c4`,`a`.`a` MOD `a`.`b` AS `c5` from `a`",
		},
		{
			name: "interval_rewrites",
			sql:  "CREATE VIEW _v AS SELECT ADDDATE(a, b) AS c1, SUBDATE(a, b) AS c2, DATE_ADD(a, INTERVAL b DAY) AS c3, DATE_SUB(a, INTERVAL b DAY) AS c4, TIMESTAMPADD(DAY, b, a) AS c5, -INTERVAL b DAY + a AS c6 FROM a",
			want: "select `a`.`a` + interval `a`.`b` day AS `c1`,`a`.`a` - interval `a`.`b` day AS `c2`,`a`.`a` + interval `a`.`b` day AS `c3`,`a`.`a` - interval `a`.`b` day AS `c4`,`a`.`a` + interval `a`.`b` day AS `c5`,`a`.`a` - interval `a`.`b` day AS `c6` from `a`",
		},
		{
			name: "min_max_distinct_and_trim",
			sql:  "CREATE VIEW _v AS SELECT MIN(DISTINCT a) AS c1, MAX(DISTINCT a) AS c2, TRIM(a FROM b) AS c3 FROM a",
			want: "select min(`a`.`a`) AS `c1`,max(`a`.`a`) AS `c2`,trim(both `a`.`a` from `a`.`b`) AS `c3` from `a`",
		},
		{
			name: "cast_varchar_and_decimal",
			sql:  "CREATE VIEW _v AS SELECT CAST(a AS VARCHAR(10)) AS c1, CAST(a AS DECIMAL) AS c2 FROM a",
			want: "select cast(`a`.`a` as char(10) charset utf8mb4) AS `c1`,cast(`a`.`a` as decimal(10)) AS `c2` from `a`",
		},
		{
			name: "weight_string_forms",
			sql:  "CREATE VIEW _v AS SELECT WEIGHT_STRING(a) AS c1, WEIGHT_STRING(a AS CHAR(1)) AS c2, WEIGHT_STRING(a AS BINARY(1)) AS c3, WEIGHT_STRING(a LEVEL 1 DESC) AS c4, WEIGHT_STRING(a LEVEL 1 REVERSE) AS c5 FROM a",
			want: "select weight_string(`a`.`a`,0,0,1) AS `c1`,weight_string(`a`.`a`,0,1,65) AS `c2`,weight_string(cast(`a`.`a` as char(1) charset binary),0,1,65) AS `c3`,weight_string(`a`.`a`,0,0,257) AS `c4`,weight_string(`a`.`a`,0,0,65537) AS `c5` from `a`",
		},
		{
			name: "introduced_hex_alias_and_default_col",
			sql:  "CREATE VIEW _v AS SELECT _utf8mb4 0x41, DEFAULT(a) FROM a",
			want: "select _utf8mb4'A' AS `_utf8mb4 0x41`,default(`a`.`a`) AS `DEFAULT(a)` from `a`",
		},
		{
			name: "introduced_string_with_comment_keeps_string_value_alias",
			sql:  "CREATE VIEW _v AS SELECT _utf8mb4 /* comment */ 'hello' FROM a",
			want: "select _utf8mb4'hello' AS `hello` from `a`",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := LoadSQL("CREATE DATABASE test; USE test; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT);")
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			results, err := c.Exec(tc.sql, nil)
			if err != nil || (len(results) > 0 && results[0].Error != nil) {
				t.Fatalf("Exec(%q) failed: parse=%v results=%+v", tc.sql, err, results)
			}
			got := c.ShowCreateView("test", "_v")
			if !strings.Contains(got, " AS "+tc.want) {
				t.Errorf("ShowCreateView:\n got: %s\nwant substring: AS %s", got, tc.want)
			}
		})
	}
}

func TestViewConvertPreservesCastIR(t *testing.T) {
	c, err := LoadSQL("CREATE DATABASE test; USE test; CREATE TABLE t (a VARCHAR(10));")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	// CONVERT(a, UNSIGNED) must produce a CastExprQ (not a bare VarExprQ with VARCHAR provenance).
	if results, err := c.Exec("CREATE VIEW v1 AS SELECT CONVERT(a, UNSIGNED) AS c FROM t", nil); err != nil || results[0].Error != nil {
		t.Fatalf("CREATE VIEW v1 failed: parse=%v results=%+v", err, results)
	}
	v1 := c.GetDatabase("test").Views["v1"]
	if v1.AnalyzedQuery == nil || len(v1.AnalyzedQuery.TargetList) != 1 {
		t.Fatalf("v1.AnalyzedQuery = %+v, want 1 target", v1.AnalyzedQuery)
	}
	castQ, ok := v1.AnalyzedQuery.TargetList[0].Expr.(*CastExprQ)
	if !ok || castQ.TargetType == nil || castQ.TargetType.BaseType != BaseTypeBigInt || !castQ.TargetType.Unsigned {
		t.Fatalf("v1 target expr = %#v, want CastExprQ(BIGINT UNSIGNED)", v1.AnalyzedQuery.TargetList[0].Expr)
	}
	if v1.AnalyzedQuery.TargetList[0].ResOrigCol != "" {
		t.Errorf("v1 target ResOrigCol = %q, want empty (conversion must not assign direct column provenance)", v1.AnalyzedQuery.TargetList[0].ResOrigCol)
	}

	// CONVERT(a USING utf8mb4) has no IR node yet, so AnalyzedQuery must stay nil rather than unwrap to VarExprQ.
	if results, err := c.Exec("CREATE VIEW v2 AS SELECT CONVERT(a USING utf8mb4) AS c FROM t", nil); err != nil || results[0].Error != nil {
		t.Fatalf("CREATE VIEW v2 failed: parse=%v results=%+v", err, results)
	}
	v2 := c.GetDatabase("test").Views["v2"]
	if v2.AnalyzedQuery != nil {
		t.Errorf("v2.AnalyzedQuery = %+v, want nil for unsupported CONVERT(... USING ...)", v2.AnalyzedQuery)
	}

	// Hex/bit literals have no distinct ConstExprQ kind yet, so AnalyzedQuery must stay nil
	// rather than serialize them as quoted strings or decimal numbers.
	for _, tc := range []struct{ name, sql string }{
		{"v_hex", "CREATE VIEW v_hex AS SELECT 0x41 AS c"},
		{"v_bit", "CREATE VIEW v_bit AS SELECT b'101' AS c"},
	} {
		if results, err := c.Exec(tc.sql, nil); err != nil || results[0].Error != nil {
			t.Fatalf("Exec(%q) failed: parse=%v results=%+v", tc.sql, err, results)
		}
		if v := c.GetDatabase("test").Views[tc.name]; v.AnalyzedQuery != nil {
			t.Errorf("%s.AnalyzedQuery = %+v, want nil for hex/bit literal", tc.name, v.AnalyzedQuery)
		}
	}
}

func TestViewSchemaQualifiedStarAcrossDatabasesAndSequenceFuncs(t *testing.T) {
	c, err := LoadSQL(`
		CREATE DATABASE db1;
		CREATE DATABASE db2;
		CREATE TABLE db1.t (x INT);
		CREATE TABLE db2.t (y INT);
		USE db1;
		CREATE SEQUENCE s;
	`)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	// db2.t.* must find db2.t even when db1.t precedes it in FROM.
	if results, err := c.Exec("CREATE VIEW v_star AS SELECT db2.t.* FROM db1.t, db2.t", nil); err != nil || results[0].Error != nil {
		t.Fatalf("CREATE VIEW v_star failed: parse=%v results=%+v", err, results)
	}
	// Sequence functions referencing an existing sequence must not fail with 1054.
	for _, sql := range []string{
		"CREATE VIEW v_seq1 AS SELECT NEXTVAL(s) AS n FROM t",
		"CREATE VIEW v_seq2 AS SELECT LASTVAL(s) AS n FROM t",
		"CREATE VIEW v_seq3 AS SELECT SETVAL(s, 1) AS n FROM t",
	} {
		if results, err := c.Exec(sql, nil); err != nil || results[0].Error != nil {
			t.Fatalf("Exec(%q) failed: parse=%v results=%+v", sql, err, results)
		}
	}

	// Wrapper CTEs on a parenthesized query must be preserved on the returned AnalyzedQuery.
	if results, err := c.Exec("CREATE VIEW v_paren_cte AS WITH c AS (SELECT 1 AS x) (SELECT x FROM c)", nil); err != nil || results[0].Error != nil {
		t.Fatalf("CREATE VIEW v_paren_cte failed: parse=%v results=%+v", err, results)
	}
	vParen := c.GetDatabase("db1").Views["v_paren_cte"]
	if vParen.AnalyzedQuery == nil || len(vParen.AnalyzedQuery.CTEList) != 1 || vParen.AnalyzedQuery.CTEList[0].Name != "c" {
		t.Fatalf("v_paren_cte.AnalyzedQuery.CTEList = %+v, want [c]", vParen.AnalyzedQuery)
	}
	if dep := DeparseQuery(vParen.AnalyzedQuery); !strings.Contains(strings.ToUpper(dep), "WITH") {
		t.Errorf("DeparseQuery(v_paren_cte) = %q, want WITH clause preserved", dep)
	}
}
