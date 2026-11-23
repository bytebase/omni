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
