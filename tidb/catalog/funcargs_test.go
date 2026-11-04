package catalog

import "testing"

// A functional key part is where CAST(expr AS type ARRAY) is legal at all —
// TiDB answers 1235 for the cast anywhere else — so the cast has to reach the
// catalog and come back out with its ARRAY. A parse that accepts the DDL and
// then records nothing, or records the cast without the keyword, is a
// different index: the one is multi-valued and the other is not, and the
// accept/reject gate cannot see the difference.
//
// The stored text is the engine's: SHOW CREATE TABLE on TiDB v8.5.5 and MySQL
// 8.0 renders the key part as ((cast(`b` as signed array))). The charset after
// a CHAR cast is omni's own and predates this test.
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
// The stored texts are TiDB v8.5.5's, read with SHOW CREATE TABLE; it keeps
// POSITION where MySQL 8.0 stores locate. An aggregate cannot appear here at
// all: the engine answers 1111 for `g INT AS (COUNT(c1))`.
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
		"p": "position(_utf8mb4'x' in `s`)",
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
