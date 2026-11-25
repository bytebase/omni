package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/trino/parser"
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

	// A parse error's position is a byte offset into the whole script.
	bad := "CREATE TABLE t (a integer);\nCREATE TABLE ("
	results, _ = newMemory().Exec(bad, nil)
	var pe *parser.ParseError
	if len(results) != 2 || !errors.As(results[1].Error, &pe) || pe.Position < strings.Index(bad, "\n") {
		t.Fatalf("got %+v, want a parse error positioned in the second statement", results)
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
		"CREATE TABLE bin (v varbinary DEFAULT X'0a ff')",
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
	if col := s.GetTable("bin").GetColumn("v"); col == nil || col.Default != "X'0AFF'" {
		t.Errorf("bin.v = %+v, want default X'0AFF'", col)
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

// TestExecQueryReferences pins how a query's column references resolve:
// qualified by schema and table, through parentheses, and a qualified star
// with aliases.
func TestExecQueryReferences(t *testing.T) {
	c := New()
	s := c.EnsureCatalog("memory").EnsureSchema("s")
	c.SetCurrentCatalog("memory")
	c.SetCurrentSchema("s")
	for _, sql := range []string{
		"CREATE TABLE t (id integer, a varchar)",
		"CREATE TABLE r (id bigint)",
		"CREATE TABLE q1 AS SELECT s.t.id, (a) FROM s.t",
		"CREATE TABLE q2 AS SELECT s.t.* FROM s.t",
		"CREATE TABLE q3 AS SELECT t.* AS (x, y) FROM t",
		"CREATE TABLE q4 AS SELECT id FROM t JOIN r USING (id)",
	} {
		if results, err := c.Exec(sql, nil); err != nil || len(results) != 1 || results[0].Error != nil {
			t.Fatalf("%s: got %+v, %v", sql, results, err)
		}
	}
	for table, want := range map[string]string{"q1": "id,a", "q2": "id,a", "q3": "x,y", "q4": "id"} {
		if got := strings.Join(s.GetTable(table).ColumnNames(), ","); got != want {
			t.Errorf("%s's columns = %s, want %s", table, got, want)
		}
	}
}

// TestExecParenJoinAndViewContext pins an alias on a parenthesized join,
// and that a view in another schema than the session's keeps no
// definition, which would be read in the wrong schema.
func TestExecParenJoinAndViewContext(t *testing.T) {
	c := New()
	db := c.EnsureCatalog("memory")
	s := db.EnsureSchema("s")
	other := db.EnsureSchema("other")
	c.SetCurrentCatalog("memory")
	c.SetCurrentSchema("s")
	for _, sql := range []string{
		"CREATE TABLE l (a integer)",
		"CREATE TABLE r (b bigint)",
		"CREATE TABLE joined AS SELECT j.a, j.b FROM (l CROSS JOIN r) AS j",
		"CREATE VIEW other.v AS SELECT a FROM l",
		"CREATE TABLE ul (a varchar, id integer)",
		"CREATE TABLE ur (id integer, b bigint)",
		"CREATE TABLE u AS SELECT * FROM ul JOIN ur USING (id)",
		"CREATE VIEW sv AS SELECT a FROM l",
		"ALTER VIEW sv RENAME TO other.sv",
		"CREATE TABLE x (x bigint)",
		"CREATE TABLE cu (one, two, three, four) AS SELECT * FROM x, ul JOIN ur USING (id)",
		"CREATE TABLE rowt (p ROW(\"a b\" integer))",
		"CREATE VIEW rowv AS SELECT p.\"a b\" FROM rowt",
		"CREATE TABLE uz (c varchar, id integer)",
		"CREATE VIEW nested AS SELECT * FROM ul JOIN (ur JOIN uz USING (id)) USING (id)",
		"CREATE VIEW nestedcols AS SELECT id, a, b, c FROM ul JOIN (ur JOIN uz USING (id)) USING (id)",
		"CREATE TABLE quoted (p ROW(\"a,b\" integer, \"row(x)\" bigint))",
		// Shapes Exec takes without typing them, as Trino takes them.
		"CREATE TABLE pj AS SELECT l.a, r.b FROM (l CROSS JOIN r)",
		"CREATE TABLE nj AS SELECT * FROM ul NATURAL JOIN ur",
		"CREATE TABLE rec AS WITH RECURSIVE nums(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM nums WHERE n < 3) SELECT n FROM nums",
	} {
		if results, err := c.Exec(sql, nil); err != nil || len(results) != 1 || results[0].Error != nil {
			t.Fatalf("%s: got %+v, %v", sql, results, err)
		}
	}
	if got := strings.Join(s.GetTable("joined").ColumnNames(), ","); got != "a,b" {
		t.Errorf("joined's columns = %s, want a,b", got)
	}
	if v := other.GetView("v"); v == nil || v.Definition != "" || strings.Join(v.ColumnNames(), ",") != "a" {
		t.Errorf("other.v = %+v, want column a and no definition", v)
	}
	if v := other.GetView("sv"); v == nil || v.Definition != "" {
		t.Errorf("other.sv = %+v, want no definition after a move", v)
	}
	if got := strings.Join(s.GetTable("u").ColumnNames(), ","); got != "id,a,b" {
		t.Errorf("u's columns = %s, want id,a,b", got)
	}
	var types []string
	for _, col := range s.GetTable("cu").Columns() {
		types = append(types, col.Type)
	}
	if got := strings.Join(types, ","); got != "bigint,integer,varchar,bigint" {
		t.Errorf("cu's types = %s, want x's bigint before the join's", got)
	}
	if cols := s.GetView("rowv").Columns(); len(cols) != 1 || cols[0].Type != "integer" {
		t.Errorf("rowv's columns = %+v, want an integer", cols)
	}
	if got := strings.Join(s.GetView("nested").ColumnNames(), ","); got != "id,a,b,c" {
		t.Errorf("nested's columns = %s, want id,a,b,c", got)
	}
	var nestedTypes []string
	for _, col := range s.GetView("nestedcols").Columns() {
		nestedTypes = append(nestedTypes, col.Type)
	}
	if got := strings.Join(nestedTypes, ","); got != "integer,varchar,bigint,varchar" {
		t.Errorf("nestedcols's types = %s", got)
	}
	if got, want := s.GetTable("quoted").ShowCreate("memory"), "CREATE TABLE memory.s.quoted (\n   p ROW(\"a,b\" integer, \"row(x)\" bigint)\n)"; got != want {
		t.Errorf("quoted = %q, want %q", got, want)
	}
	if got := strings.Join(s.GetTable("pj").ColumnNames(), ","); got != "a,b" {
		t.Errorf("pj's columns = %s, want a,b", got)
	}
}
