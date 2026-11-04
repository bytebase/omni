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
