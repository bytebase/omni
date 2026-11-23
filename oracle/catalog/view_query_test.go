package catalog

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestViewQueryDerivesColumns(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE dept (deptno NUMBER(2) PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
CREATE VIEW v AS SELECT d.*, UPPER(loc) up FROM dept d;
`)
	if err != nil {
		t.Fatal(err)
	}
	views := c.Metadata().GetSchemas()[0].GetViews()
	if len(views) != 1 {
		t.Fatalf("views = %d, want 1", len(views))
	}
	v := views[0]
	if want := `SELECT d."DEPTNO",d."DNAME",d."LOC", UPPER(loc) up FROM dept d`; v.GetDefinition() != want {
		t.Errorf("definition = %q, want %q", v.GetDefinition(), want)
	}
	want := []struct {
		name, typ string
		nullable  bool
	}{
		{"DEPTNO", "NUMBER(2)", false},
		{"DNAME", "VARCHAR2(14 BYTE)", false},
		{"LOC", "VARCHAR2(13 BYTE)", true},
		{"UP", "VARCHAR2(13 BYTE)", true},
	}
	if len(v.GetColumns()) != len(want) {
		t.Fatalf("columns = %d, want %d", len(v.GetColumns()), len(want))
	}
	for i, w := range want {
		got := v.GetColumns()[i]
		if got.GetName() != w.name || got.GetType() != w.typ || got.GetNullable() != w.nullable {
			t.Errorf("column %d = %s %s nullable=%v, want %s %s nullable=%v", i, got.GetName(), got.GetType(), got.GetNullable(), w.name, w.typ, w.nullable)
		}
	}
}

func TestViewQueryValidation(t *testing.T) {
	const missing = `CREATE VIEW v AS SELECT * FROM nosuch`

	// By default a view over objects the catalog does not hold loads, as a
	// schema dump that names other schemas must.
	c, err := LoadSQL(missing)
	if err != nil {
		t.Fatalf("default config refused the view: %v", err)
	}
	if v := c.Metadata().GetSchemas()[0].GetViews()[0]; v.GetDefinition() != "SELECT * FROM nosuch" || len(v.GetColumns()) != 0 {
		t.Errorf("view = %q with %d columns, want its text and no columns", v.GetDefinition(), len(v.GetColumns()))
	}

	// The pinned server's configuration refuses it as Oracle does.
	c = New()
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(missing, nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if len(results) != 1 || !errors.As(results[0].Error, &refusal) || refusal.Code != "00942" {
		t.Fatalf("results = %+v, want ORA-00942", results)
	}

	// FORCE keeps the view either way.
	results, err = c.Exec(`CREATE FORCE VIEW v AS SELECT * FROM nosuch`, nil)
	if err != nil || results[0].Error != nil {
		t.Fatalf("FORCE view refused: %v %+v", err, results)
	}
}

func TestViewQueryRefusals(t *testing.T) {
	const fixture = `CREATE TABLE bl (b BLOB); CREATE TABLE t (id NUMBER, n NUMBER(5,-1), m NUMBER(*,-2), doc CLOB, "USER" VARCHAR2(20), d DATE, nv NVARCHAR2(10), vc VARCHAR2(10), "CURRENT_DATE" DATE); CREATE TABLE u (id DATE); CREATE TABLE w (id NUMBER);`
	cases := []struct {
		name, sql, code string
	}{
		{"recursive_anchor_reads_itself", `CREATE VIEW v AS WITH r(n) AS (SELECT n FROM r) SELECT n FROM r`, "32040"},
		{"recursive_anchor_union_reads_itself", `CREATE VIEW v AS WITH r(n) AS (SELECT n FROM r UNION ALL SELECT n FROM r) SELECT n FROM r`, "32043"},
		{"compound_unaliased_expression", `CREATE VIEW v AS SELECT 1 FROM dual UNION ALL SELECT 2 FROM dual`, "00998"},
		{"named_window_missing_column", `CREATE VIEW v AS SELECT SUM(id) OVER w s FROM t WINDOW w AS (PARTITION BY nosuch)`, "00904"},
		{"in_subquery_too_wide", `CREATE VIEW v AS SELECT id FROM t WHERE id IN (SELECT id, id FROM t)`, "00913"},
		{"in_row_subquery_too_narrow", `CREATE VIEW v AS SELECT id FROM t WHERE (id, n) IN (SELECT id FROM t)`, "00947"},
		{"qualified_grouped_column", `CREATE VIEW v AS SELECT t.id FROM t GROUP BY id`, ""},
		{"cross_join_lateral_reads_left", `CREATE VIEW v AS SELECT x.id FROM t CROSS JOIN LATERAL (SELECT t.id FROM dual) x`, ""},
		{"case_unknown_function", `CREATE VIEW v AS SELECT CASE WHEN 1=1 THEN TO_MULTI_BYTE('x') ELSE 'x' END x FROM dual`, ""},
		{"coalesce_unknown_function", `CREATE VIEW v AS SELECT COALESCE(TO_MULTI_BYTE('x'), 1) x FROM dual`, ""},
		{"unpivot_columns_by_position", `CREATE VIEW v AS SELECT * FROM (SELECT id n1, TO_CHAR(id) s1, n n2, TO_CHAR(n) s2 FROM t) UNPIVOT ((n, s) FOR src IN ((n1, s1) AS 'A', (n2, s2) AS 'B'))`, ""},
		{"star_not_single_group", `CREATE VIEW v AS SELECT t.*, COUNT(*) c FROM t`, "00937"},
		{"star_not_grouped", `CREATE VIEW v AS SELECT t.*, COUNT(*) c FROM t GROUP BY id`, "00979"},
		{"merged_and_plain_column", `CREATE VIEW v AS SELECT id FROM t JOIN t u USING (id), t w`, "00918"},
		{"fetch_bind", `CREATE VIEW v AS SELECT id FROM t FETCH FIRST :n ROWS ONLY`, "01027"},
		{"offset_bind", `CREATE VIEW v AS SELECT id FROM t OFFSET :n ROWS`, "01027"},
		{"match_recognize_variable_qualifier", `CREATE VIEW v AS SELECT * FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST(a.id) AS first_id PATTERN (a) DEFINE a AS a.id > 0)`, ""},
		{"compound_unknown_function", `CREATE VIEW v AS SELECT TO_MULTI_BYTE('x') x FROM dual UNION ALL SELECT 'x' FROM dual`, ""},
		{"table_constructor_bind", `CREATE VIEW v AS SELECT * FROM TABLE(sys.odcinumberlist(:x))`, "01027"},
		{"window_frame_bind", `CREATE VIEW v AS SELECT SUM(id) OVER (ORDER BY id ROWS :n PRECEDING) s FROM t`, "01027"},
		{"qualified_rowid", `CREATE VIEW v AS SELECT t.ROWID rid FROM t`, ""},
		{"xmltable_bind", `CREATE VIEW v AS SELECT * FROM XMLTABLE('/x' PASSING :x COLUMNS id NUMBER PATH '.')`, "01027"},
		{"order_by_not_grouped", `CREATE VIEW v AS SELECT id, COUNT(*) c FROM t GROUP BY id ORDER BY n`, "00979"},
		{"order_by_alias_grouped", `CREATE VIEW v AS SELECT id x, COUNT(*) c FROM t GROUP BY id ORDER BY x, 2`, ""},
		{"pivot_xml", `CREATE VIEW v AS SELECT n_xml FROM t PIVOT XML (SUM(id) FOR n IN (ANY))`, ""},
		{"with_forward_reference", `CREATE VIEW v AS WITH a AS (SELECT id FROM b), b AS (SELECT id FROM t) SELECT id FROM a`, "00942"},
		{"grouped_subquery_outer_reference", `CREATE VIEW v AS SELECT (SELECT COUNT(*) + u.n FROM t) x FROM t u`, ""},
		{"model_measure_alias", `CREATE VIEW v AS SELECT s FROM t MODEL DIMENSION BY (id) MEASURES (n AS s) RULES (s[1] = 0)`, ""},
		{"model_rule_bind", `CREATE VIEW v AS SELECT n FROM t MODEL DIMENSION BY (id) MEASURES (n) RULES (n[1] = :x)`, "01027"},
		{"having_without_group_by", `CREATE VIEW v AS SELECT id FROM t HAVING id > 1`, "00937"},
		{"distinct_order_by_unselected", `CREATE VIEW v AS SELECT DISTINCT id FROM t ORDER BY n`, "01791"},
		{"distinct_order_by_selected", `CREATE VIEW v AS SELECT DISTINCT id x, n + 1 y FROM t ORDER BY x, n + 1, 1`, ""},
		{"any_subquery_too_wide", `CREATE VIEW v AS SELECT id FROM t WHERE id = ANY (SELECT id, n FROM t)`, "00913"},
		{"all_row_subquery", `CREATE VIEW v AS SELECT id FROM t WHERE (id, n) = ALL (SELECT id, n FROM t)`, ""},
		{"row_equals_subquery", `CREATE VIEW v AS SELECT id FROM t WHERE (id, n) = (SELECT id, n FROM t)`, ""},
		{"xmltable_default_bind", `CREATE VIEW v AS SELECT * FROM XMLTABLE('/x' COLUMNS id NUMBER PATH '.' DEFAULT :x)`, "01027"},
		{"distinct_order_by_built_from_selected", `CREATE VIEW v AS SELECT DISTINCT id FROM t ORDER BY id + 1`, ""},
		{"distinct_order_by_built_from_unselected", `CREATE VIEW v AS SELECT DISTINCT id FROM t ORDER BY id + n`, "01791"},
		{"match_recognize_offset_bind", `CREATE VIEW v AS SELECT * FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST(a.id, :x) AS first_id PATTERN (a) DEFINE a AS id > 0)`, "01027"},
		{"group_by_aggregate_alias", `CREATE VIEW v AS SELECT SUM(n) s FROM t GROUP BY s`, "00934"},
		{"analytic_in_aggregate", `CREATE VIEW v AS SELECT SUM(ROW_NUMBER() OVER (ORDER BY id)) s FROM t`, "30483"},
		{"correlated_ungrouped", `CREATE VIEW v AS SELECT (SELECT t.n FROM dual) x FROM t GROUP BY t.id`, "00979"},
		{"correlated_grouped", `CREATE VIEW v AS SELECT (SELECT t.id FROM dual) x FROM t GROUP BY t.id`, ""},
		{"search_missing_column", `CREATE VIEW v AS WITH r(n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SEARCH DEPTH FIRST BY nosuch SET ord SELECT n FROM r`, "00904"},
		{"undefined_window", `CREATE VIEW v AS SELECT SUM(id) OVER nosuch s FROM t`, "32784"},
		{"qualify_alias", `CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ORDER BY id) rn FROM t QUALIFY rn = 1`, ""},
		{"distinct_to_char_clob", `CREATE VIEW v AS SELECT DISTINCT TO_CHAR(TO_CLOB('x')) x FROM t`, ""},
		{"rollup_composite_column", `CREATE VIEW v AS SELECT id, n, COUNT(*) c FROM t GROUP BY ROLLUP((id, n))`, ""},
		{"qualified_duplicate_column", `CREATE VIEW v AS SELECT x.c FROM (SELECT id c, n c FROM t) x`, "00918"},
		{"date_plus_date", `CREATE VIEW v AS SELECT DATE '2020-01-01' + DATE '2020-01-02' x FROM dual`, "00975"},
		{"recursive_multi_branch_anchor", `CREATE VIEW v AS WITH r(n) AS (SELECT 1 FROM dual UNION ALL SELECT 2 FROM dual UNION ALL SELECT n + 2 FROM r WHERE n < 5) SELECT n FROM r`, ""},
		{"analytic_in_window_order", `CREATE VIEW v AS SELECT SUM(id) OVER (ORDER BY ROW_NUMBER() OVER (ORDER BY id)) s FROM t`, "30483"},
		{"json_table_reads_join_left", `CREATE VIEW v AS SELECT jt.id FROM t CROSS JOIN JSON_TABLE(t.doc, '$' COLUMNS (id NUMBER PATH '$.id')) jt`, ""},
		{"match_recognize_repeated_variable", `CREATE VIEW v AS SELECT * FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST(a.id) AS first_id PATTERN (a b* a) DEFINE a AS id > 0)`, ""},
		{"lob_comparison", `CREATE VIEW v AS SELECT id FROM t WHERE doc = doc`, "00932"},
		{"lob_is_null", `CREATE VIEW v AS SELECT id FROM t WHERE doc IS NULL`, ""},
		{"sum_of_date", `CREATE VIEW v AS SELECT SUM(DATE '2020-01-01') s FROM dual`, "00932"},
		{"keyword_named_column_not_grouped", `CREATE VIEW v AS SELECT "USER", COUNT(*) c FROM t`, "00937"},
		{"using_incompatible_keys", `CREATE VIEW v AS SELECT id FROM t JOIN u USING (id)`, "00932"},
		{"subquery_comparison_types", `CREATE VIEW v AS SELECT id FROM t WHERE id = (SELECT d FROM t)`, "00932"},
		{"in_subquery_types", `CREATE VIEW v AS SELECT id FROM t WHERE id IN (SELECT d FROM t)`, "00932"},
		{"recursive_member_distinct", `CREATE VIEW v AS WITH r(n) AS (SELECT 1 FROM dual UNION ALL SELECT DISTINCT n + 1 FROM r WHERE n < 5) SELECT n FROM r`, "32486"},
		{"between_incompatible", `CREATE VIEW v AS SELECT 1 x FROM dual WHERE 1 BETWEEN DATE '2020-01-01' AND DATE '2020-01-02'`, "00932"},
		{"distinct_lpad_lob", `CREATE VIEW v AS SELECT DISTINCT LPAD(doc, 10, 'x') x FROM t`, "22848"},
		{"parenthesized_order_missing", `CREATE VIEW v AS (SELECT id FROM t) ORDER BY nosuch`, "00904"},
		{"parenthesized_fetch_bind", `CREATE VIEW v AS (SELECT id FROM t) FETCH FIRST :n ROWS ONLY`, "01027"},
		{"in_list_incompatible", `CREATE VIEW v AS SELECT 1 x FROM dual WHERE 1 IN (DATE '2020-01-01')`, "00932"},
		{"in_row_list_width", `CREATE VIEW v AS SELECT id FROM t WHERE (id, n) IN ((1, 2), (3))`, "00947"},
		{"order_by_position_lob", `CREATE VIEW v AS SELECT doc FROM t ORDER BY 1`, "22848"},
		{"group_by_lob_alias", `CREATE VIEW v AS SELECT doc x FROM t GROUP BY x`, "22848"},
		{"simple_case_incompatible", `CREATE VIEW v AS SELECT CASE 1 WHEN DATE '2020-01-01' THEN 1 ELSE 0 END x FROM dual`, "00932"},
		{"coalesce_implicit_number", `CREATE VIEW v AS SELECT COALESCE(1, '2') x FROM dual`, ""},
		{"range_offset_character_key", `CREATE VIEW v AS SELECT SUM(id) OVER (ORDER BY "USER" RANGE 1 PRECEDING) s FROM t`, "00902"},
		{"range_offset_numeric_key", `CREATE VIEW v AS SELECT SUM(id) OVER (ORDER BY id RANGE 1 PRECEDING) s FROM t`, ""},
		{"level_without_connect_by", `CREATE VIEW v AS SELECT LEVEL lvl FROM dual`, "01788"},
		{"level_with_connect_by", `CREATE VIEW v AS SELECT LEVEL lvl FROM dual CONNECT BY LEVEL < 3`, ""},
		{"where_date", `CREATE VIEW v AS SELECT 1 x FROM dual WHERE DATE '2020-01-01'`, "00920"},
		{"group_by_subquery", `CREATE VIEW v AS SELECT COUNT(*) c FROM dual GROUP BY (SELECT 1 FROM dual)`, "22818"},
		{"group_by_subquery_alias", `CREATE VIEW v AS SELECT (SELECT 1 FROM dual) s, COUNT(*) c FROM dual GROUP BY s`, "22818"},
		{"and_with_date_operand", `CREATE VIEW v AS SELECT 1 x FROM dual WHERE 1=1 AND DATE '2020-01-01'`, "00920"},
		{"named_window_not_grouped", `CREATE VIEW v AS SELECT id, COUNT(*) OVER w c FROM t GROUP BY id WINDOW w AS (ORDER BY n)`, "00979"},
		{"as_of_scn_plausible", `CREATE VIEW v AS SELECT id FROM t AS OF SCN 123456`, ""},
		{"as_of_scn_one", `CREATE VIEW v AS SELECT id FROM t AS OF SCN 1`, "08181"},
		{"row_comparison_types", `CREATE VIEW v AS SELECT 1 x FROM dual WHERE (1, DATE '2020-01-01') = (1, 1)`, "00932"},
		{"row_comparison_width", `CREATE VIEW v AS SELECT 1 x FROM dual WHERE (1, 2) = (1, 2, 3)`, "00913"},
		{"date_times_number", `CREATE VIEW v AS SELECT DATE '2020-01-01' * 2 x FROM dual`, "00932"},
		{"qualify_not_grouped", `CREATE VIEW v AS SELECT id, COUNT(*) c FROM t GROUP BY id QUALIFY ROW_NUMBER() OVER (ORDER BY n) = 1`, "00979"},
		{"qualify_grouped_alias", `CREATE VIEW v AS SELECT id, COUNT(*) c FROM t GROUP BY id QUALIFY ROW_NUMBER() OVER (ORDER BY c) = 1`, ""},
		{"number_minus_date", `CREATE VIEW v AS SELECT 1 - DATE '2020-01-01' x FROM dual`, "00932"},
		{"pivot_xml_subquery_missing_table", `CREATE VIEW v AS SELECT n_xml FROM t PIVOT XML (SUM(id) FOR n IN (SELECT id FROM nosuch))`, "00942"},
		{"extract_hour_from_date", `CREATE VIEW v AS SELECT EXTRACT(HOUR FROM DATE '2020-01-01') x FROM dual`, "30076"},
		{"extract_year_from_date", `CREATE VIEW v AS SELECT EXTRACT(YEAR FROM DATE '2020-01-01') x FROM dual`, ""},
		{"json_table_bind_path", `CREATE VIEW v AS SELECT * FROM JSON_TABLE('{}', :p COLUMNS (id NUMBER PATH '$.id')) jt`, "01027"},
		{"quoted_keyword_name", `CREATE VIEW v AS SELECT "CURRENT_DATE" x FROM dual`, "00904"},
		{"compound_order_lob", `CREATE VIEW v AS SELECT doc FROM t UNION ALL SELECT doc FROM t ORDER BY 1`, "22848"},
		{"sum_two_arguments", `CREATE VIEW v AS SELECT SUM(id, id) s FROM t`, "00909"},
		{"upper_two_arguments", `CREATE VIEW v AS SELECT UPPER(id, id) s FROM t`, "00909"},
		{"qualified_star_merged_not_grouped", `CREATE VIEW v AS SELECT t.* FROM t JOIN w USING (id) GROUP BY n`, "00979"},
		{"stddev_interval", `CREATE VIEW v AS SELECT STDDEV(INTERVAL '1' DAY) s FROM dual`, "00932"},
		{"sum_interval", `CREATE VIEW v AS SELECT SUM(INTERVAL '1' DAY) s FROM dual`, ""},
		{"negated_date", `CREATE VIEW v AS SELECT -DATE '2020-01-01' x FROM dual`, "00932"},
		{"interval_leading_precision", `CREATE VIEW v AS SELECT INTERVAL '123' YEAR x FROM dual`, "01873"},
		{"interval_leading_precision_given", `CREATE VIEW v AS SELECT INTERVAL '123' YEAR(3) x FROM dual`, ""},
		{"compound_charset_mismatch", `CREATE VIEW v AS SELECT nv FROM t UNION ALL SELECT vc FROM t`, "12704"},
		{"subset_unknown_member", `CREATE VIEW v AS SELECT first_id FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST(a.id) AS first_id PATTERN (a) SUBSET s = (missing) DEFINE a AS id > 0)`, "62503"},
		{"nvl_date_number", `CREATE VIEW v AS SELECT NVL(DATE '2020-01-01', 1) x FROM dual`, "00932"},
		{"partition_by_lob", `CREATE VIEW v AS SELECT COUNT(*) OVER (PARTITION BY doc) c FROM t`, "22848"},
		{"quoted_pattern_variable", `CREATE VIEW v AS SELECT first_id FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST("a".id) AS first_id PATTERN ("a") DEFINE "a" AS id > 0)`, ""},
		{"range_interval_numeric_key", `CREATE VIEW v AS SELECT SUM(id) OVER (ORDER BY id RANGE INTERVAL '1' DAY PRECEDING) x FROM t`, "00902"},
		{"range_interval_date_key", `CREATE VIEW v AS SELECT SUM(id) OVER (ORDER BY d RANGE INTERVAL '1' DAY PRECEDING) x FROM t`, ""},
		{"table_cast_collection", `CREATE VIEW v AS SELECT column_value cv FROM TABLE(CAST(sys.odcinumberlist(1, 2) AS sys.odcinumberlist))`, ""},
		{"json_object_returning_number", `CREATE VIEW v AS SELECT JSON_OBJECT('a' VALUE 1 RETURNING NUMBER) x FROM dual`, "40449"},
		{"json_value_returning_number", `CREATE VIEW v AS SELECT JSON_VALUE('{"a":1}', '$.a' RETURNING NUMBER) x FROM dual`, ""},
		{"cast_date_to_number", `CREATE VIEW v AS SELECT CAST(DATE '2020-01-01' AS NUMBER) x FROM dual`, "00932"},
		{"subset_two_members", `CREATE VIEW v AS SELECT first_id FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST(a.id) AS first_id PATTERN (a b) SUBSET s = (a, b) DEFINE a AS id > 0, b AS id > 0)`, ""},
		{"qualified_star_only_key_not_grouped", `CREATE VIEW v AS SELECT w.* FROM w JOIN t USING (id) GROUP BY t.n`, "00979"},
		{"scalar_with_over", `CREATE VIEW v AS SELECT UPPER('x') OVER () x FROM dual`, "00923"},
		{"extract_from_number", `CREATE VIEW v AS SELECT EXTRACT(YEAR FROM 1) x FROM dual`, "30076"},
		{"abs_of_date", `CREATE VIEW v AS SELECT ABS(DATE '2020-01-01') x FROM dual`, "00932"},
		{"rowid_with_aggregate", `CREATE VIEW v AS SELECT ROWID rid, COUNT(*) c FROM t`, "00937"},
		{"regexp_substr_national_union", `CREATE VIEW v AS SELECT REGEXP_SUBSTR(nv, N'.') x FROM t UNION ALL SELECT nv FROM t`, ""},
		{"cycle_mark_bind", `CREATE VIEW v AS WITH r(n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 2) CYCLE n SET cyc TO :p DEFAULT 'N' SELECT n, cyc FROM r`, "01027"},
		{"listagg_order_analytic", `CREATE VIEW v AS SELECT LISTAGG(dummy, ',') WITHIN GROUP (ORDER BY ROW_NUMBER() OVER (ORDER BY dummy)) x FROM dual`, "30483"},
		{"distinct_parenthesized_order", `CREATE VIEW v AS SELECT DISTINCT (t.id) x FROM t ORDER BY t.id`, ""},
		{"dual_flashback_bind", `CREATE VIEW v AS SELECT * FROM dual AS OF SCN :p`, "01027"},
		{"aggregate_three_deep", `CREATE VIEW v AS SELECT MAX(SUM(AVG(id))) x FROM t GROUP BY n`, "00935"},
		{"interval_plus_number", `CREATE VIEW v AS SELECT INTERVAL '1' DAY + 1 x FROM dual`, "00932"},
		{"row_number_unordered", `CREATE VIEW v AS SELECT ROW_NUMBER() OVER () x FROM dual`, "30485"},
		{"row_number_named_ordered", `CREATE VIEW v AS SELECT ROW_NUMBER() OVER w x FROM t WINDOW w AS (ORDER BY id)`, ""},
		{"extract_tz_from_region_literal", `CREATE VIEW v AS SELECT EXTRACT(TIMEZONE_HOUR FROM TIMESTAMP '2020-01-01 00:00:00 UTC') x FROM dual`, ""},
		{"interval_times_interval", `CREATE VIEW v AS SELECT INTERVAL '1' DAY * INTERVAL '2' DAY x FROM dual`, "00932"},
		{"number_over_interval", `CREATE VIEW v AS SELECT 2 / INTERVAL '2' DAY x FROM dual`, "00932"},
		{"interval_times_number", `CREATE VIEW v AS SELECT INTERVAL '1' DAY * 2 x FROM dual`, ""},
		{"trim_mode_grouped", `CREATE VIEW v AS SELECT TRIM(BOTH FROM dummy) x, COUNT(*) c FROM dual GROUP BY dummy`, ""},
		{"distinct_order_keyword_column", `CREATE VIEW v AS SELECT DISTINCT id FROM t ORDER BY "CURRENT_DATE"`, "01791"},
		{"recursive_member_reads_twice", `CREATE VIEW v AS WITH r(n) AS (SELECT 1 FROM dual UNION ALL SELECT a.n+1 FROM r a JOIN r b ON a.n=b.n WHERE a.n<2) SELECT n FROM r`, "32042"},
		{"nullif_number_date", `CREATE VIEW v AS SELECT NULLIF(1, DATE '2020-01-01') x FROM dual`, "00932"},
		{"unpivot_bind_label", `CREATE VIEW v AS SELECT * FROM t UNPIVOT (a FOR tag IN (id AS :x))`, "56900"},
		{"rowid_of_distinct_view", `CREATE VIEW v AS SELECT q.ROWID rid FROM (SELECT DISTINCT id FROM t) q`, "01446"},
		{"rowid_of_simple_view", `CREATE VIEW v AS SELECT q.ROWID rid FROM (SELECT id FROM t) q`, ""},
		{"xmlelement_tag_grouped", `CREATE VIEW v AS SELECT XMLELEMENT("tag", COUNT(*)) x FROM dual`, ""},
		{"sum_of_lob", `CREATE VIEW v AS SELECT SUM(doc) x FROM t`, "00932"},
		{"blob_plus_number", `CREATE VIEW v AS SELECT b + 1 x FROM bl`, "00932"},
		{"grouping_without_group_by", `CREATE VIEW v AS SELECT GROUPING(id) g FROM t`, "30481"},
		{"ansi_join_with_plus", `CREATE VIEW v AS SELECT a.id FROM t a LEFT JOIN t b ON a.id=b.id WHERE a.id=b.id(+)`, "25156"},
		{"prior_without_connect_by", `CREATE VIEW v AS SELECT PRIOR id x FROM t`, "01788"},
		{"pivot_lob_grouping_column", `CREATE VIEW v AS SELECT * FROM (SELECT id, n, doc FROM t) PIVOT (SUM(id) FOR n IN (1 AS one))`, "22848"},
		{"fetch_first_date", `CREATE VIEW v AS SELECT dummy FROM dual FETCH FIRST DATE '2020-01-01' ROWS ONLY`, "00932"},
		{"cycle_marks_too_long", `CREATE VIEW v AS WITH r(n) AS (SELECT 1 FROM dual UNION ALL SELECT n+1 FROM r WHERE n<2) CYCLE n SET cyc TO 'YES' DEFAULT 'NO' SELECT n, cyc FROM r`, "32047"},
		{"distinct_analytic_ordered", `CREATE VIEW v AS SELECT SUM(DISTINCT 1) OVER (ORDER BY dummy) x FROM dual`, "30487"},
		{"distinct_analytic_partitioned", `CREATE VIEW v AS SELECT COUNT(DISTINCT dummy) OVER (PARTITION BY dummy) x FROM dual`, ""},
		{"trim_mode_from", `CREATE VIEW v AS SELECT TRIM(BOTH 'x' FROM TO_CHAR(id)) trimmed FROM t`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, strict := range []bool{false, true} {
				c, err := LoadSQL(fixture)
				if err != nil {
					t.Fatal(err)
				}
				if strict {
					c.SetServerConfig(StandardOracle23ServerConfig())
				}
				results, err := c.Exec(tc.sql, nil)
				if err != nil {
					t.Fatal(err)
				}
				var refusal *Error
				got := ""
				if errors.As(results[0].Error, &refusal) {
					got = refusal.Code
				} else if results[0].Error != nil {
					t.Fatalf("strict=%v: %v", strict, results[0].Error)
				}
				want := tc.code
				if !strict {
					want = ""
					if tc.code == "00998" {
						// The catalog refuses an unaliased view expression
						// under any configuration.
						want = tc.code
					}
				}
				if got != want {
					t.Errorf("strict=%v: code %q, want %q", strict, got, want)
				}
			}
		})
	}
}

func TestViewQueryNegativeScales(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (n NUMBER(5,-1), m NUMBER(*,-2));
CREATE VIEW v AS SELECT n, m FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	cols := c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns()
	if len(cols) != 2 || cols[0].GetType() != "NUMBER(5,-1)" || cols[1].GetType() != "NUMBER(*,-2)" {
		t.Fatalf("columns = %v, want NUMBER(5,-1) and NUMBER(*,-2)", cols)
	}
}

func TestViewQueryTrimSourceType(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (name VARCHAR2(100));
CREATE VIEW v AS SELECT TRIM('x' FROM name) a, TRIM(LEADING 'x' FROM name) b FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != "VARCHAR2(100 BYTE)" {
			t.Errorf("%s type = %s, want VARCHAR2(100 BYTE)", col.GetName(), col.GetType())
		}
	}
}

func TestViewQueryPercentileType(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (d DATE, n NUMBER(5,2));
CREATE VIEW v AS SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY d) p, PERCENTILE_DISC(0.5) WITHIN GROUP (ORDER BY d) q, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY n) r FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"DATE", "DATE", "NUMBER"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %s, want %s", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryJSONReturning(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT JSON_VALUE('{"id":1}', '$.id' RETURNING NUMBER) a, JSON_OBJECT('a' VALUE 1 RETURNING CLOB) b, JSON_QUERY('{}', '$' RETURNING VARCHAR2) c, JSON_VALUE('{}', '$.a') d FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"NUMBER", "CLOB", "VARCHAR2(4000 BYTE)", "VARCHAR2(4000 BYTE)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %s, want %s", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryStandardHashLength(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT STANDARD_HASH('x') a, STANDARD_HASH('x', 'SHA256') b, STANDARD_HASH('x', 'MD5') c, STANDARD_HASH('x', 'SHA512') d FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"RAW(20)", "RAW(32)", "RAW(16)", "RAW(64)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %s, want %s", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryCharacterResults(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT 'a' || 'bc' a, 'a' || TO_CHAR(1) b, UPPER(123) u FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"CHAR(3 BYTE)", "VARCHAR2(41 BYTE)", "VARCHAR2(40 BYTE)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %s, want %s", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryDerivedTypesMore(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER NOT NULL, c CLOB, k CHAR(3));
CREATE VIEW v AS SELECT id, COUNT(*) n FROM t GROUP BY ROLLUP (id);
CREATE VIEW w AS SELECT TO_CHAR(c) a, TO_CHAR(k) b, NVL(1, CAST(2 AS BINARY_DOUBLE)) d FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	views := c.Metadata().GetSchemas()[0].GetViews()
	for _, v := range views {
		switch v.GetName() {
		case "V":
			if !v.GetColumns()[0].GetNullable() {
				t.Errorf("V.ID is NOT NULL, want nullable under ROLLUP")
			}
		case "W":
			want := []string{"VARCHAR2(4000 BYTE)", "VARCHAR2(3 BYTE)", "BINARY_DOUBLE"}
			for i, col := range v.GetColumns() {
				if col.GetType() != want[i] {
					t.Errorf("W.%s type = %s, want %s", col.GetName(), col.GetType(), want[i])
				}
			}
		}
	}
}

func TestViewQueryStarKeepsDuplicateNamesApart(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER);
CREATE VIEW v (a, b) AS SELECT * FROM t x CROSS JOIN t y;
CREATE VIEW w AS SELECT N'a' || 'b' n FROM dual;
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		switch v.GetName() {
		case "V":
			if want := `SELECT "X"."ID","Y"."ID" FROM t x CROSS JOIN t y`; v.GetDefinition() != want {
				t.Errorf("V = %q, want %q", v.GetDefinition(), want)
			}
		case "W":
			if got := v.GetColumns()[0].GetType(); got != "NVARCHAR2(2)" {
				t.Errorf("W.N type = %s, want NVARCHAR2(2)", got)
			}
		}
	}
}

func TestViewQueryTypesRound11(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (a NUMBER);
CREATE VIEW v AS SELECT ROUND(CAST(1.25 AS BINARY_DOUBLE), 1) x, ROUND(CAST(1.25 AS BINARY_DOUBLE)) y, LISTAGG(HEXTORAW('AA')) WITHIN GROUP (ORDER BY dummy) z FROM dual;
CREATE VIEW w AS SELECT * FROM t UNPIVOT (val FOR (yr, qtr) IN (a AS (2020, 'Q1')));
`)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		"V": {"NUMBER", "BINARY_DOUBLE", "RAW(2000)"},
		"W": {"NUMBER", "VARCHAR2(2 BYTE)", "NUMBER"},
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		for i, col := range v.GetColumns() {
			if col.GetType() != want[v.GetName()][i] {
				t.Errorf("%s.%s type = %s, want %s", v.GetName(), col.GetName(), col.GetType(), want[v.GetName()][i])
			}
		}
	}
}

func TestViewQueryTypesRound13(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t1 (id NUMBER, secret NUMBER INVISIBLE);
CREATE TABLE t2 (id NUMBER, secret NUMBER INVISIBLE);
CREATE VIEW v AS SELECT TO_CHAR(DATE '2020-01-01', 'MONTH') m, TO_CHAR(DATE '2020-01-01', 'YYYY-MM-DD HH24:MI') d, SIGN(CAST(1 AS BINARY_DOUBLE)) s, TO_CHAR(1, 'FM999') f FROM dual;
CREATE VIEW w AS SELECT * FROM t1 NATURAL JOIN t2;
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		var got []string
		for _, col := range v.GetColumns() {
			got = append(got, col.GetName()+" "+col.GetType())
		}
		want := map[string]string{
			"V": "M ,D VARCHAR2(16 BYTE),S NUMBER,F ",
			"W": "ID NUMBER",
		}[v.GetName()]
		if strings.Join(got, ",") != want {
			t.Errorf("%s columns = %q, want %q", v.GetName(), strings.Join(got, ","), want)
		}
	}
}

func TestViewQueryTypesRound14(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT INTERVAL '1' DAY + DATE '2020-01-01' a, NULLIF(1, CAST(2 AS BINARY_DOUBLE)) b FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"DATE", "BINARY_DOUBLE"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %s, want %s", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryTypesRound15(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER, grp VARCHAR2(10), val DATE);
CREATE VIEW v AS SELECT REPLACE(123, '2', '4') x FROM dual;
CREATE VIEW w AS SELECT * FROM t MATCH_RECOGNIZE (PARTITION BY grp ORDER BY id MEASURES FIRST(a.val) AS first_val ALL ROWS PER MATCH PATTERN (a) DEFINE a AS id > 0);
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		var got []string
		for _, col := range v.GetColumns() {
			got = append(got, col.GetName()+" "+col.GetType())
		}
		want := map[string]string{
			"V": "X VARCHAR2(40 BYTE)",
			"W": "GRP VARCHAR2(10 BYTE),ID NUMBER,FIRST_VAL DATE,VAL DATE",
		}[v.GetName()]
		if strings.Join(got, ",") != want {
			t.Errorf("%s columns = %q, want %q", v.GetName(), strings.Join(got, ","), want)
		}
	}
}

func TestViewQueryTypesRound16(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT INTERVAL '400' DAY(3) a, INTERVAL '1 2:3:4.5' DAY TO SECOND(3) b, INTERVAL '5' YEAR(4) c, RAWTOHEX(123) d, RAWTOHEX(HEXTORAW('AB')) e FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"INTERVAL DAY(3) TO SECOND(6)", "INTERVAL DAY(2) TO SECOND(3)", "INTERVAL YEAR(4) TO MONTH", "", "VARCHAR2(2 BYTE)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryUnresolvedViews(t *testing.T) {
	// A view the catalog cannot derive leaves a view over it unexpanded.
	c, err := LoadSQL(`
CREATE VIEW a AS SELECT * FROM ext.t;
CREATE VIEW b AS SELECT * FROM a;
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		if v.GetName() == "B" && v.GetDefinition() != "SELECT * FROM a" {
			t.Errorf("B = %q, want its text unexpanded", v.GetDefinition())
		}
	}

	// IF NOT EXISTS keeps an existing view without reading the new query.
	c, err = LoadSQL(`CREATE TABLE t (id NUMBER); CREATE VIEW v AS SELECT id FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(`CREATE VIEW IF NOT EXISTS v AS SELECT id FROM missing_table`, nil)
	if err != nil || results[0].Error != nil {
		t.Fatalf("IF NOT EXISTS on an existing view: %v %+v", err, results)
	}
}

func TestViewQueryTypesRound18(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (nd NCLOB, j JSON);
CREATE VIEW v AS SELECT MOD(1, CAST(2 AS BINARY_DOUBLE)) a, CONCAT(nd, 'x') b, JSON_MERGEPATCH(j, '{"x":1}') c FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"BINARY_DOUBLE", "NCLOB", "JSON"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryReplaceLength(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (s VARCHAR2(10)); CREATE VIEW v AS SELECT REPLACE(s, 'a', 'b') a, REPLACE(s, 'a', 'abcdefghij') b, REPLACE(s, 'a') c, REPLACE(123, '2', '4') d FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"VARCHAR2(10 BYTE)", "", "VARCHAR2(10 BYTE)", "VARCHAR2(40 BYTE)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryTypesRound20(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (a NVARCHAR2(10), b NVARCHAR2(20));
CREATE VIEW v AS SELECT CASE WHEN 1=1 THEN a ELSE b END x, SIN(CAST(1 AS BINARY_FLOAT)) y, COALESCE(1, '2') z FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"NVARCHAR2(20)", "BINARY_DOUBLE", "NUMBER"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryTrimDatetime(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT TRIM(DATE '2020-01-01') a, LTRIM(DATE '2020-01-01') b FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	for _, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != "VARCHAR2(75 BYTE)" {
			t.Errorf("%s type = %q, want VARCHAR2(75 BYTE)", col.GetName(), col.GetType())
		}
	}
}

func TestViewQueryNationalResults(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (n NCHAR(10)); CREATE VIEW v AS SELECT SUBSTR(n, 2) a, TO_NCHAR(DATE '2020-01-01', 'YYYY') b, TO_NCHAR(SYSDATE) c FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"NVARCHAR2(10)", "NVARCHAR2(4)", ""}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryTypesRound23(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (a VARCHAR2(10 CHAR), b VARCHAR2(20 CHAR));
CREATE VIEW v AS SELECT DECODE(1, 1, NULL, 2) x, NVL(a, b) y, LOG(10, CAST(100 AS BINARY_DOUBLE)) z FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"VARCHAR2(40 BYTE)", "VARCHAR2(20 CHAR)", "BINARY_DOUBLE"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryPeriodsAndNames(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (id NUMBER, PERIOD FOR valid_time); CREATE TABLE u ("CAFÉ" NUMBER);`)
	if err != nil {
		t.Fatal(err)
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(`CREATE VIEW v AS SELECT id FROM t AS OF PERIOD FOR valid_time TIMESTAMP '2020-01-01 00:00:00';
CREATE VIEW w AS SELECT * FROM (SELECT café+1 FROM u);`, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("%s: %v", r.SQL, r.Error)
		}
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		if v.GetName() == "W" && v.GetColumns()[0].GetName() != "CAFÉ+1" {
			t.Errorf("W column = %q, want %q", v.GetColumns()[0].GetName(), "CAFÉ+1")
		}
	}
}

func TestViewQueryTypesRound25(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (n NCHAR(10)); CREATE VIEW v AS SELECT TRIM(n) a, LTRIM(n) b, FLOOR(DATE '2020-01-01') c, CEIL(SYSTIMESTAMP) d FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"NVARCHAR2(10)", "NVARCHAR2(10)", "DATE", "DATE"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}

func TestViewQueryQualifiedStarUsingKey(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER, n NUMBER);
CREATE TABLE u (id NUMBER);
CREATE VIEW v AS SELECT t.* FROM t JOIN u USING (id);
`)
	if err != nil {
		t.Fatal(err)
	}
	v := c.Metadata().GetSchemas()[0].GetViews()[0]
	if want := `SELECT "ID",t."N" FROM t JOIN u USING (id)`; v.GetDefinition() != want {
		t.Fatalf("definition = %q, want %q", v.GetDefinition(), want)
	}
	// The stored text loads again.
	if _, err := LoadSQL("CREATE TABLE t (id NUMBER, n NUMBER); CREATE TABLE u (id NUMBER); CREATE VIEW v AS " + v.GetDefinition()); err != nil {
		t.Fatal(err)
	}
}

func TestViewQueryTypesRound27(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT ATAN2(CAST(1 AS BINARY_FLOAT), 2) a, SUBSTR(DATE '2020-01-01', 2) b FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"BINARY_DOUBLE", "VARCHAR2(75 BYTE)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
	c, err = LoadSQL(`CREATE TABLE p (id NUMBER) PARTITION BY RANGE (id) (PARTITION p1 VALUES LESS THAN (10));`)
	if err != nil {
		t.Fatal(err)
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	for sql, code := range map[string]string{
		`CREATE VIEW v AS SELECT id FROM p PARTITION (no_such)`: "02149",
		`CREATE VIEW w AS SELECT id FROM p PARTITION (p1)`:      "",
	} {
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatal(err)
		}
		var refusal *Error
		got := ""
		if errors.As(results[0].Error, &refusal) {
			got = refusal.Code
		}
		if got != code {
			t.Errorf("%s: code %q, want %q", sql, got, code)
		}
	}
}

func TestViewQueryStarQualifiesInvisibleClash(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (a NUMBER, id NUMBER INVISIBLE);
CREATE TABLE u (id NUMBER);
CREATE VIEW v AS SELECT * FROM t CROSS JOIN u;
`)
	if err != nil {
		t.Fatal(err)
	}
	v := c.Metadata().GetSchemas()[0].GetViews()[0]
	if want := `SELECT "A","U"."ID" FROM t CROSS JOIN u`; v.GetDefinition() != want {
		t.Errorf("definition = %q, want %q", v.GetDefinition(), want)
	}
}

func TestViewQuerySubstrbBytes(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (c VARCHAR2(100 BYTE)); CREATE VIEW v AS SELECT SUBSTRB(c, 1, 2) x FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns()[0].GetType(); got != "VARCHAR2(2 BYTE)" {
		t.Errorf("type = %q, want VARCHAR2(2 BYTE)", got)
	}
}

func TestViewQueryNestedComparisonsStayLinear(t *testing.T) {
	q := "SELECT 1 FROM dual"
	for i := 0; i < 40; i++ {
		q = "SELECT 1 FROM dual WHERE (" + q + ") = (SELECT 1 FROM dual)"
	}
	done := make(chan error, 1)
	go func() {
		_, err := LoadSQL("CREATE VIEW v AS SELECT 1 x FROM dual WHERE (" + q + ") = (SELECT 1 FROM dual)")
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("nested subquery comparisons took more than 10s")
	}
}

func TestViewQueryNestedWithStaysLinear(t *testing.T) {
	q := "SELECT 1 x FROM dual"
	for i := 0; i < 40; i++ {
		q = "WITH c AS (SELECT (" + q + ") x FROM dual) SELECT x FROM c"
	}
	done := make(chan error, 1)
	go func() {
		_, err := LoadSQL("CREATE VIEW v AS " + q)
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("nested WITH clauses took more than 10s")
	}
}

func TestViewQueryRound33(t *testing.T) {
	// Nested recursive WITH anchors stay linear.
	q := "SELECT 1 x FROM dual"
	for i := 0; i < 30; i++ {
		q = fmt.Sprintf("WITH r%d(x) AS (SELECT (%s) FROM dual UNION ALL SELECT x + 1 FROM r%d WHERE x < 0) SELECT x FROM r%d", i, q, i, i)
	}
	done := make(chan error, 1)
	go func() {
		_, err := LoadSQL("CREATE VIEW v AS " + q)
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("nested recursive WITH anchors took more than 10s")
	}

	// MAX_STRING_SIZE=EXTENDED keeps a concatenation's length.
	c, err := LoadSQL(`CREATE TABLE t (c VARCHAR2(3000 BYTE), id NUMBER, PERIOD FOR valid_time);`)
	if err != nil {
		t.Fatal(err)
	}
	cfg := StandardOracle23ServerConfig()
	cfg.MaxStringSize = "EXTENDED"
	c.SetServerConfig(cfg)
	results, err := c.Exec(`CREATE VIEW v AS SELECT c || c x FROM t;
CREATE VIEW w AS SELECT id FROM t AS OF PERIOD FOR valid_time :p;`, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Error != nil {
		t.Fatal(results[0].Error)
	}
	var refusal *Error
	if !errors.As(results[1].Error, &refusal) || refusal.Code != "01027" {
		t.Errorf("PERIOD FOR bind: %v, want ORA-01027", results[1].Error)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		if v.GetName() == "V" && v.GetColumns()[0].GetType() != "VARCHAR2(6000 BYTE)" {
			t.Errorf("V.X type = %q, want VARCHAR2(6000 BYTE)", v.GetColumns()[0].GetType())
		}
	}
}

func TestViewQueryLongExpressionStaysFast(t *testing.T) {
	expr := "1" + strings.Repeat("+1", 4000)
	done := make(chan error, 1)
	go func() {
		c, err := LoadSQL("CREATE TABLE t (id NUMBER);")
		if err == nil {
			_, err = c.Exec("CREATE VIEW v AS SELECT "+expr+" x, COUNT(*) c FROM t GROUP BY id", nil)
		}
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("a 4,000-term grouped expression took more than 10s")
	}
}

func TestViewQueryPartitionForBind(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE p (id NUMBER) PARTITION BY RANGE (id) (PARTITION p1 VALUES LESS THAN (10));`)
	if err != nil {
		t.Fatal(err)
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(`CREATE VIEW v AS SELECT id FROM p PARTITION FOR (:p)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if !errors.As(results[0].Error, &refusal) || refusal.Code != "01027" {
		t.Errorf("PARTITION FOR bind: %v, want ORA-01027", results[0].Error)
	}
}

func TestViewQueryRound36(t *testing.T) {
	// Nested TABLE() collections stay linear.
	q := "SELECT 1 FROM dual"
	for i := 0; i < 30; i++ {
		q = "SELECT column_value FROM TABLE(sys.odcinumberlist((" + q + ")))"
	}
	done := make(chan error, 1)
	go func() {
		_, err := LoadSQL("CREATE VIEW v AS SELECT (" + q + ") x, COUNT(*) c FROM dual")
		done <- err
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("nested TABLE() collections took more than 10s")
	}

	// t.* whose USING key would be ambiguous bare keeps its text.
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER, n NUMBER);
CREATE TABLE u (id NUMBER);
CREATE TABLE w (id NUMBER);
CREATE VIEW v AS SELECT t.* FROM t JOIN u USING (id) CROSS JOIN w;
`)
	if err != nil {
		t.Fatal(err)
	}
	v := c.Metadata().GetSchemas()[0].GetViews()[0]
	if want := "SELECT t.* FROM t JOIN u USING (id) CROSS JOIN w"; v.GetDefinition() != want {
		t.Errorf("definition = %q, want %q", v.GetDefinition(), want)
	}
}

func TestViewQueryStarOverUnnamedDuplicates(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v (a, b) AS SELECT * FROM (SELECT 1 id FROM dual), (SELECT 2 id FROM dual);`)
	if err != nil {
		t.Fatal(err)
	}
	v := c.Metadata().GetSchemas()[0].GetViews()[0]
	if want := "SELECT * FROM (SELECT 1 id FROM dual), (SELECT 2 id FROM dual)"; v.GetDefinition() != want {
		t.Errorf("definition = %q, want %q", v.GetDefinition(), want)
	}
}

func TestViewQueryStarOverDuplicatesInOneItem(t *testing.T) {
	c, err := LoadSQL(`
CREATE VIEW v (a, b) AS SELECT q.* FROM (SELECT 1 id, 2 id FROM dual) q;
CREATE VIEW w (a, b) AS SELECT * FROM (SELECT 1 id, 2 id FROM dual) q;
`)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"V": "SELECT q.* FROM (SELECT 1 id, 2 id FROM dual) q",
		"W": "SELECT * FROM (SELECT 1 id, 2 id FROM dual) q",
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		if v.GetDefinition() != want[v.GetName()] {
			t.Errorf("%s = %q, want %q", v.GetName(), v.GetDefinition(), want[v.GetName()])
		}
	}
}

func TestViewQueryStarOverRepeatedQualifier(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER);
CREATE TABLE u (n NUMBER);
CREATE VIEW v (a, b, c) AS SELECT * FROM t CROSS JOIN u CROSS JOIN t;
`)
	if err != nil {
		t.Fatal(err)
	}
	v := c.Metadata().GetSchemas()[0].GetViews()[0]
	if want := "SELECT * FROM t CROSS JOIN u CROSS JOIN t"; v.GetDefinition() != want {
		t.Errorf("definition = %q, want %q", v.GetDefinition(), want)
	}
}

func TestViewQueryRound41(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t ("*" NUMBER, id NUMBER, name VARCHAR2(10));
CREATE VIEW q AS SELECT t."*" a FROM t;
CREATE VIEW m AS SELECT x FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES FIRST(a.name) || 'x' AS x PATTERN (a) DEFINE a AS id > 0);
CREATE VIEW v AS SELECT 1 x FROM dual;
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		switch v.GetName() {
		case "Q":
			if len(v.GetColumns()) != 1 || v.GetDefinition() != `SELECT t."*" a FROM t` {
				t.Errorf("Q = %q with %d columns", v.GetDefinition(), len(v.GetColumns()))
			}
		case "M":
			if got := v.GetColumns()[0].GetType(); got != "VARCHAR2(11 BYTE)" {
				t.Errorf("M.X type = %q, want VARCHAR2(11 BYTE)", got)
			}
		}
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(`CREATE OR REPLACE VIEW v AS SELECT x FROM v`, nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if !errors.As(results[0].Error, &refusal) || refusal.Code != "01731" {
		t.Errorf("self-reading view: %v, want ORA-01731", results[0].Error)
	}
}

func TestViewQueryCircularThroughSynonym(t *testing.T) {
	c, err := LoadSQL(`
CREATE VIEW v AS SELECT 1 x FROM dual;
CREATE SYNONYM s FOR v;
CREATE VIEW w AS SELECT x FROM s;
`)
	if err != nil {
		t.Fatal(err)
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(`CREATE OR REPLACE VIEW v AS SELECT x FROM w`, nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if !errors.As(results[0].Error, &refusal) || refusal.Code != "01731" {
		t.Errorf("cycle through a synonym: %v, want ORA-01731", results[0].Error)
	}
}

func TestViewQuerySysDualPastLocalDual(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE dual (x NUMBER);`)
	if err != nil {
		t.Fatal(err)
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(`CREATE VIEW v AS SELECT dummy FROM sys.dual`, nil)
	if err != nil || results[0].Error != nil {
		t.Fatalf("SYS.DUAL past a local DUAL: %v %+v", err, results)
	}
}

func TestViewQueryRound45(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (id NUMBER);
CREATE SYNONYM dual FOR t;
CREATE VIEW v AS SELECT * FROM dual;
CREATE VIEW w AS SELECT COLLECT(id) x FROM t;
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		switch v.GetName() {
		case "V":
			if v.GetColumns()[0].GetName() != "ID" {
				t.Errorf("V reads %q, want the synonym's ID", v.GetColumns()[0].GetName())
			}
		case "W":
			if v.GetColumns()[0].GetType() != "" {
				t.Errorf("W.X type = %q, want unknown", v.GetColumns()[0].GetType())
			}
		}
	}
}

func TestViewQueryReplaceNchar(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (n NCHAR(10)); CREATE VIEW v AS SELECT REPLACE(n, N'a', N'b') x FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns()[0].GetType(); got != "NVARCHAR2(10)" {
		t.Errorf("type = %q, want NVARCHAR2(10)", got)
	}
}

func TestViewQueryUnknownThroughConcat(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (s VARCHAR2(10)); CREATE VIEW v AS SELECT REPLACE(s, 'a', 'bb') || 'x' x FROM t;`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns()[0].GetType(); got != "" {
		t.Errorf("type = %q, want unknown", got)
	}
}

func TestViewQueryTypesRound48(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE t (a VARCHAR2(1 CHAR), b VARCHAR2(5 BYTE));
CREATE VIEW v AS SELECT a x FROM t UNION ALL SELECT b FROM t;
CREATE VIEW w AS SELECT SYS_CONTEXT('USERENV', 'AUTHENTICATION_DATA', 4000) x FROM dual;
`)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"V": "VARCHAR2(5 CHAR)", "W": "VARCHAR2(4000 BYTE)"}
	for _, v := range c.Metadata().GetSchemas()[0].GetViews() {
		if got := v.GetColumns()[0].GetType(); got != want[v.GetName()] {
			t.Errorf("%s.X type = %q, want %q", v.GetName(), got, want[v.GetName()])
		}
	}
}

func TestViewQueryDecodePrecedence(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT DECODE(1d, 1, 1, 0) x FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns()[0].GetType(); got != "BINARY_DOUBLE" {
		t.Errorf("type = %q, want BINARY_DOUBLE", got)
	}
}

func TestViewQueryUserenvTypes(t *testing.T) {
	c, err := LoadSQL(`CREATE VIEW v AS SELECT USERENV('SESSIONID') a, USERENV('LANG') b FROM dual;`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"NUMBER", "VARCHAR2(52 BYTE)"}
	for i, col := range c.Metadata().GetSchemas()[0].GetViews()[0].GetColumns() {
		if col.GetType() != want[i] {
			t.Errorf("%s type = %q, want %q", col.GetName(), col.GetType(), want[i])
		}
	}
}
