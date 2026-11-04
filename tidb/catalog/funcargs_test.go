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
