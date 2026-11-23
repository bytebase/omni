package catalog

import (
	"errors"
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
	const fixture = `CREATE TABLE t (id NUMBER, n NUMBER(5,-1), m NUMBER(*,-2));`
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
