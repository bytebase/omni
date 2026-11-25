package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestExecResults pins what Exec reports beside the schema it stores: one
// result per statement with its line, a stop at the first refusal unless
// ContinueOnError, and Trino's errorName on a refusal.
func TestExecResults(t *testing.T) {
	newMemory := func() *Catalog {
		c := New()
		c.EnsureCatalog("memory").EnsureSchema("s")
		c.SetCurrentCatalog("memory")
		c.SetCurrentSchema("s")
		return c
	}
	sql := "CREATE TABLE t (a integer);\nDROP TABLE nope;\nCREATE TABLE \"U\" (\"B\" bigint)"
	results, err := newMemory().Exec(sql, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].Error != nil || results[1].Line != 2 {
		t.Fatalf("got %+v, want two results, the second on line 2", results)
	}
	var e *Error
	if !errors.As(results[1].Error, &e) || e.Name != "TABLE_NOT_FOUND" {
		t.Fatalf("got %v, want TABLE_NOT_FOUND", results[1].Error)
	}

	c := newMemory()
	results, err = c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil || len(results) != 3 || results[2].Error != nil {
		t.Fatalf("got %+v, %v, want three results", results, err)
	}
	// Trino stores a quoted name in lower case too.
	u := c.GetCatalog("memory").GetSchema("s").GetTable("u")
	if u == nil || u.GetColumn("b") == nil {
		t.Fatalf("want table u with column b")
	}
	if got, want := u.ShowCreate("memory"), "CREATE TABLE memory.s.u (\n   b bigint\n)"; got != want {
		t.Errorf("ShowCreate = %q, want %q", got, want)
	}

	// A statement that does not parse fails alone: the one before it runs,
	// and with ContinueOnError so does the one after.
	c = newMemory()
	results, err = c.Exec("CREATE TABLE first (x integer); CREATE TABLE (; CREATE TABLE third (x integer)", &ExecOptions{ContinueOnError: true})
	if err != nil || len(results) != 3 || results[0].Error != nil || results[1].Error == nil || results[2].Error != nil {
		t.Fatalf("got %+v, %v, want the second of three to fail", results, err)
	}
	if s := c.GetCatalog("memory").GetSchema("s"); !s.HasTable("first") || !s.HasTable("third") {
		t.Errorf("want tables first and third")
	}
}

// TestExecQueries pins what CREATE TABLE AS and CREATE VIEW store from
// their query, and the decimal bounds Trino holds a column type to.
func TestExecQueries(t *testing.T) {
	c := New()
	s := c.EnsureCatalog("memory").EnsureSchema("s")
	c.SetCurrentCatalog("memory")
	c.SetCurrentSchema("s")
	sql := "CREATE TABLE t (id integer, d decimal(38,38));\n" +
		"CREATE TABLE u AS SELECT id, d FROM t UNION ALL SELECT id, d FROM t;\n" +
		"CREATE VIEW v AS SELECT id FROM t"
	results, err := c.Exec(sql, nil)
	if err != nil || len(results) != 3 || results[2].Error != nil {
		t.Fatalf("got %+v, %v", results, err)
	}
	if got, want := s.GetTable("u").ShowCreate("memory"), "CREATE TABLE memory.s.u (\n   id integer,\n   d decimal(38, 38)\n)"; got != want {
		t.Errorf("u = %q, want %q", got, want)
	}
	if got := s.GetView("v").Definition; got != "SELECT id FROM t" {
		t.Errorf("v's definition = %q", got)
	}
	for _, bad := range []string{"decimal(39,0)", "decimal(0)", "decimal(5,6)"} {
		results, _ := c.Exec("CREATE TABLE bad (a "+bad+")", nil)
		var e *Error
		if len(results) != 1 || !errors.As(results[0].Error, &e) || e.Name != "INVALID_FUNCTION_ARGUMENT" {
			t.Errorf("%s: got %+v, want INVALID_FUNCTION_ARGUMENT", bad, results)
		}
	}
}

// TestExecQueryShapes pins the query shapes CREATE TABLE AS and CREATE VIEW
// type beyond a plain SELECT, and the types and defaults Exec takes.
func TestExecQueryShapes(t *testing.T) {
	c := New()
	s := c.EnsureCatalog("memory").EnsureSchema("s")
	c.SetCurrentCatalog("memory")
	c.SetCurrentSchema("s")
	for _, sql := range []string{
		"CREATE TABLE l (id integer, a varchar)",
		"CREATE TABLE r (id integer, b varchar)",
		"CREATE TABLE j AS SELECT * FROM l JOIN r USING (id)",
		"CREATE TABLE vals (id) AS VALUES (1), (2)",
		"CREATE TABLE rows (payload ROW(x integer))",
		"CREATE VIEW rv AS SELECT payload.x FROM rows",
		"CREATE TABLE ts (a timestamp DEFAULT TIMESTAMP '2020-01-01 00:00:00', b interval day to second)",
	} {
		if results, err := c.Exec(sql, nil); err != nil || len(results) != 1 || results[0].Error != nil {
			t.Fatalf("%s: got %+v, %v", sql, results, err)
		}
	}
	if got := s.GetTable("j").ColumnNames(); strings.Join(got, ",") != "id,a,b" {
		t.Errorf("j's columns = %v, want id, a, b", got)
	}
	if col := s.GetTable("vals").GetColumn("id"); col == nil || col.Type != "integer" {
		t.Errorf("vals.id = %+v, want an integer", col)
	}
	if cols := s.GetView("rv").Columns(); len(cols) != 1 || cols[0].Name != "x" || cols[0].Type != "integer" {
		t.Errorf("rv's columns = %+v, want x integer", cols)
	}
	// A single-field interval is no Trino type, and is refused, not a panic.
	results, err := c.Exec("CREATE TABLE i (a INTERVAL DAY)", nil)
	var e *Error
	if err != nil || len(results) != 1 || !errors.As(results[0].Error, &e) || e.Name != "TYPE_NOT_FOUND" {
		t.Errorf("INTERVAL DAY: got %+v, %v, want TYPE_NOT_FOUND", results, err)
	}
}
