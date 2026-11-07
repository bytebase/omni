package catalog

import "testing"

// The function argument forms in a generated column, where the catalog
// serializes the expression rather than keeping its text: a node its type
// switch does not know becomes "(?)", so DDL the parser accepted would reach
// the catalog with its definition gone. POSITION is the node that reaches it
// first, and EXTRACT's compound units arrive with it.
//
// The stored texts are mariadb:11.8.8's, read with SHOW CREATE TABLE; it keeps
// the call under the function's other name, without adding a charset introducer
// to its string literal. An aggregate cannot appear here: the engine answers 1901
// for `g INT AS (COUNT(c1))`.
func TestFunctionArgumentGeneratedColumns(t *testing.T) {
	c := New()
	const ddl = "CREATE DATABASE d; USE d; " +
		"CREATE TABLE tg (s VARCHAR(20), a DATETIME, c1 INT, " +
		"p INT AS (POSITION('x' IN s)) STORED, " +
		"e INT AS (EXTRACT(DAY_HOUR FROM a)) STORED, " +
		"i INT AS (INTERVAL(c1, 1, 2)) STORED)"
	results, err := c.Exec(ddl, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec: %v", r.Error)
		}
	}
	tbl := c.GetDatabase("d").GetTable("tg")
	if tbl == nil {
		t.Fatal("table tg is not in the catalog")
	}
	want := map[string]string{
		"p": "locate('x',`s`)",
		"e": "extract(day_hour from `a`)",
		"i": "interval(`c1`,1,2)",
	}
	for _, col := range tbl.Columns {
		w, ok := want[col.Name]
		if !ok {
			continue
		}
		delete(want, col.Name)
		if col.Generated == nil {
			t.Errorf("column %q: no generated expression", col.Name)
			continue
		}
		if col.Generated.Expr != w {
			t.Errorf("column %q: expression = %q, want %q", col.Name, col.Generated.Expr, w)
		}
	}
	for name := range want {
		t.Errorf("column %q is not in the catalog", name)
	}
}

// A view over POSITION has to be analyzed, not skipped. An expression the
// analyzer does not know fails the whole view body, createView stores a nil
// AnalyzedQuery, and ExpandMergeViews can then neither expand the view nor
// carry its column lineage — all without an error anywhere.
func TestFunctionArgumentViewAnalysis(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; "+
		"CREATE TABLE t (s VARCHAR(20), nn VARCHAR(20) NOT NULL, a DATETIME); "+
		"CREATE VIEW v AS SELECT POSITION('x' IN s) AS p, s AS q FROM t", nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec: %v", r.Error)
		}
	}
	v := c.GetDatabase("d").Views["v"]
	if v == nil {
		t.Fatal("view v is not in the catalog")
	}
	if v.AnalyzedQuery == nil {
		t.Error("the view body was not analyzed, so the view carries no column lineage")
	}
}

// The nullability of a POSITION column is its operands': mariadb:11.8.8 reports
// IS_NULLABLE YES for position('x' in s) over a nullable s, for
// position(s in nn), and for position('x' in NULL), and NO for a NOT NULL
// operand and for two literals. Reading the node as always non-nullable
// records a column the engine calls nullable as NOT NULL.
func TestFunctionArgumentViewColumnNullability(t *testing.T) {
	for _, tc := range []struct {
		expr     string
		nullable bool
	}{
		{"POSITION('x' IN s)", true},
		{"POSITION('x' IN nn)", false},
		{"POSITION(s IN nn)", true},
		{"POSITION(nn IN s)", true},
		{"POSITION('x' IN NULL)", true},
		{"POSITION('x' IN 'abc')", false},
	} {
		c := New()
		results, err := c.Exec("CREATE DATABASE d; USE d; "+
			"CREATE TABLE t (s VARCHAR(20), nn VARCHAR(20) NOT NULL); "+
			"CREATE VIEW v AS SELECT "+tc.expr+" AS p FROM t", nil)
		if err != nil {
			t.Fatalf("Exec(%q): %v", tc.expr, err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("Exec(%q): %v", tc.expr, r.Error)
			}
		}
		v := c.GetDatabase("d").Views["v"]
		if v == nil || len(v.ColumnMetadata) != 1 {
			t.Fatalf("%q: the view has no column metadata", tc.expr)
		}
		if v.ColumnMetadata[0].Nullable != tc.nullable {
			t.Errorf("%q: nullable = %v, want %v", tc.expr, v.ColumnMetadata[0].Nullable, tc.nullable)
		}
	}
}
