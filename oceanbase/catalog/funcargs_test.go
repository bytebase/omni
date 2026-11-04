package catalog

import "testing"

// A functional key part is where CAST(expr AS type ARRAY) is legal at all —
// the engine answers 1235 for the cast anywhere else — so the cast has to
// reach the catalog and come back out with its ARRAY. A parse that accepts the DDL and
// then records nothing, or records the cast without the keyword, is a
// different index: the one is multi-valued and the other is not, and the
// accept/reject gate cannot see the difference.
//
// The stored text follows mysql:8.0, whose SHOW CREATE TABLE renders the key
// part as ((cast(`b` as signed array))). The charset after a CHAR cast is
// omni's own and predates this test.
func TestFunctionArgumentKeyParts(t *testing.T) {
	c := New()
	const ddl = "CREATE DATABASE d; USE d; " +
		"CREATE TABLE tm (b JSON, " +
		"KEY mv ((CAST(b AS SIGNED ARRAY))), " +
		"KEY mvc ((CAST(b AS CHAR(10) ARRAY))), " +
		"KEY plain ((CAST(b AS SIGNED))), " +
		"KEY ex ((EXTRACT(DAY_HOUR FROM b))), " +
		"KEY pos ((POSITION('x' IN b))))"
	results, err := c.Exec(ddl, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec: %v", r.Error)
		}
	}
	tbl := c.GetDatabase("d").GetTable("tm")
	if tbl == nil {
		t.Fatal("table tm is not in the catalog")
	}
	want := map[string]string{
		"mv":    "cast(`b` as signed array)",
		"mvc":   "cast(`b` as char(10) charset utf8mb4 array)",
		"plain": "cast(`b` as signed)",
		"ex":    "extract(day_hour from `b`)",
		"pos":   "position('x' in `b`)",
	}
	for _, idx := range tbl.Indexes {
		w, ok := want[idx.Name]
		if !ok {
			t.Errorf("unexpected index %q", idx.Name)
			continue
		}
		delete(want, idx.Name)
		if len(idx.Columns) != 1 {
			t.Errorf("index %q: %d key parts, want 1", idx.Name, len(idx.Columns))
			continue
		}
		if got := idx.Columns[0].Expr; got != w {
			t.Errorf("index %q: key part = %q, want %q", idx.Name, got, w)
		}
	}
	for name := range want {
		t.Errorf("index %q is not in the catalog", name)
	}
}

// The same forms in a generated column, where the catalog serializes the
// expression rather than keeping its text: a node its type switch does not
// know becomes "(?)", so DDL the parser accepted would reach the catalog with
// its definition gone. POSITION is the node that reaches it first.
//
// The stored texts follow mysql:8.0, which keeps the call under the function's
// other name. An aggregate cannot appear here at all: the engine answers 1111
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
// The verdicts follow mysql:8.0, whose functional-index rules this
// validator holds the engine to: it answers 3758 for a disallowed function in
// either operand and creates the index over a column.
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
