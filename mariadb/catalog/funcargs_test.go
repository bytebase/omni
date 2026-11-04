package catalog

import "testing"

// The function argument forms in a generated column, where the catalog
// serializes the expression rather than keeping its text: a node its type
// switch does not know becomes "(?)", so DDL the parser accepted would reach
// the catalog with its definition gone. POSITION is the node that reaches it
// first, and EXTRACT's compound units arrive with it.
//
// The stored texts are mariadb:11.8.8's, read with SHOW CREATE TABLE; it keeps
// the call under the function's other name. The charset introducer on a string
// literal is omni's own — this engine stores concat('x',`s`) too — and predates
// this test. An aggregate cannot appear here at all: the engine answers 1901
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
		"p": "locate(_utf8mb4'x',`s`)",
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

// A function the engine disallows in a functional index is still disallowed
// inside POSITION's operands: the validator walks the expression, and a node
// its type switch does not know returns false without looking at either
// operand, so KEY ((position('x' in rand()))) would be recorded as an index
// the engine refuses to create.
//
// This engine has no functional-index syntax of its own —
// mariadb:11.8.8 answers 1064 for KEY ((...)) — so what this pins is the
// inherited validator: a disallowed function cannot hide inside POSITION.
func TestFunctionArgumentFunctionalIndexValidation(t *testing.T) {
	for _, tc := range []struct {
		ddl    string
		reject bool
	}{
		{"CREATE TABLE ti (s VARCHAR(20), KEY k ((POSITION('x' IN s))))", false},
		{"CREATE TABLE ti (s VARCHAR(20), KEY k ((POSITION('x' IN RAND()))))", true},
		{"CREATE TABLE ti (s VARCHAR(20), KEY k ((POSITION(NOW() IN s))))", true},
		{"CREATE TABLE ti (s VARCHAR(20), KEY k ((RAND())))", true},
	} {
		c := New()
		results, err := c.Exec("CREATE DATABASE d; USE d; "+tc.ddl, nil)
		if err != nil {
			t.Fatalf("Exec(%q): %v", tc.ddl, err)
		}
		var failed error
		for _, r := range results {
			if r.Error != nil {
				failed = r.Error
			}
		}
		if (failed != nil) != tc.reject {
			t.Errorf("%q: error = %v, want rejected = %v", tc.ddl, failed, tc.reject)
			continue
		}
		if tc.reject {
			continue
		}
		if tbl := c.GetDatabase("d").GetTable("ti"); tbl == nil || len(tbl.Indexes) != 1 {
			t.Errorf("%q: the index is not in the catalog", tc.ddl)
		}
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
