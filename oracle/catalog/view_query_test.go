package catalog

import (
	"errors"
	"strings"
	"testing"
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
	const fixture = `CREATE TABLE t (id NUMBER, n NUMBER(5,-1), m NUMBER(*,-2), doc CLOB, "USER" VARCHAR2(20), d DATE); CREATE TABLE u (id DATE);`
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
