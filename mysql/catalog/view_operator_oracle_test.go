package catalog

import (
	"strings"
	"testing"
)

// TestViewOperatorsMatchMySQL stores each operator shape on a live MySQL 8.0 and
// through this catalog and compares the view text. A view is where a wrong parse
// tree stops being invisible: `NOT a = b` read as `(NOT a) = b` and `a > ALL (q)`
// read as `a > (q)` both store a definition that means something else, and neither
// shows in any accept-or-reject check.
func TestViewOperatorsMatchMySQL(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	if err := ctr.execSQL("CREATE TABLE tt (a INT, b INT, c INT)"); err != nil {
		t.Fatal(err)
	}
	body := func(stmt string) string {
		i := strings.Index(strings.ToLower(stmt), " as select ")
		if i < 0 {
			t.Fatalf("no AS SELECT in %q", stmt)
		}
		return stmt[i+len(" as "):]
	}
	for _, expr := range []string{
		"a || b",
		"a && b",
		"a || b && c",
		"not a = b",
		"not a < b",
		"a = b is null",
		"a is null = b",
		"a between 1 and 2 | 3",
		"a = any (select 1)",
		"a = some (select 1)",
		"a = all (select 1)",
		"a <> any (select 1)",
		"a <> some (select 1)",
		"a <> all (select 1)",
		"a > any (select 1)",
		"a > some (select 1)",
		"a > all (select 1)",
		"a <= all (select 1)",
		"not a = any (select 1)",
		"not a > all (select 1)",
		"not a > any (select 1)",
		"not a = all (select 1)",
		"not a <> any (select 1)",
		"not a <> all (select 1)",
		"not a >= any (select 1)",
		"not a >= all (select 1)",
		"not a < any (select 1)",
		"not a < all (select 1)",
		"not a <= any (select 1)",
		"not a <= all (select 1)",
		"not a in (select 1)",
		"(a > all (select 1)) = 1",
		"(a > any (select 1)) and b",
		"b or (a > all (select 1))",
		"(a = any (select 1)) and b",
		"(a <> all (select 1)) and b",
		"(a > all (select 1)) is null",
		"(a = any (select 1)) is null",
		"1 = (a = any (select 1))",
		"1 = (a <> all (select 1))",
		"(a <> all (select 1)) + 1",
		"(a > all (select 1)) > all (select 1)",
		"(a <> all (select 1)) > all (select 1)",
		"-(a > all (select 1))",
		"abs(a > all (select 1))",
		"coalesce(a = any (select 1), a <> all (select 1))",
		"cast(a > all (select 1) as signed)",
		"if(a > all (select 1), 1, 2)",
		"case a > all (select 1) when 1 then 2 end",
	} {
		t.Run(expr, func(t *testing.T) {
			stmt := "CREATE OR REPLACE VIEW v AS SELECT " + expr + " AS x FROM tt"
			if err := ctr.execSQL(stmt); err != nil {
				t.Fatal(err)
			}
			engine, err := ctr.showCreateView("v")
			if err != nil {
				t.Fatal(err)
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE tt (a INT, b INT, c INT); "+stmt, nil); err != nil {
				t.Fatal(err)
			}
			if got, want := body(c.ShowCreateView("test", "v")), body(engine); got != want {
				t.Errorf("view over %q\n omni %s\nmysql %s", expr, got, want)
			}
		})
	}
}

// TestViewNestedQuantifiedMySQLDiffers records where this catalog keeps what was written
// and MySQL does not, so a server upgrade that changes either answer shows up here.
// MySQL prints a nested quantified comparison with no parentheses. For an arithmetic
// operand on the left that text is a syntax error, so the view is created and then
// answers 1064 whenever it is read. For a right operand it parses, with other grouping:
// the view over `1 = (a > ALL (q))` stores `(1 = a) > all (q)` and returns 0 for a row
// where the same SELECT returns 1.
func TestViewNestedQuantifiedMySQLDiffers(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	for _, stmt := range []string{"CREATE TABLE tt (a INT, b INT, c INT)", "INSERT INTO tt VALUES (2, 0, 0)"} {
		if err := ctr.execSQL(stmt); err != nil {
			t.Fatal(err)
		}
	}

	if err := ctr.execSQL("CREATE OR REPLACE VIEW v AS SELECT (a > all (select 1)) + 1 AS x FROM tt"); err != nil {
		t.Fatal(err)
	}
	if _, err := ctr.showCreateView("v"); err == nil {
		t.Error("MySQL can now read a view over (a > all (q)) + 1; recheck the writer's note")
	}

	const expr = "1 = (a > all (select 1))"
	if err := ctr.execSQL("CREATE OR REPLACE VIEW v AS SELECT " + expr + " AS x FROM tt"); err != nil {
		t.Fatal(err)
	}
	engine, err := ctr.showCreateView("v")
	if err != nil {
		t.Fatal(err)
	}
	if want := "(1 = `tt`.`a`) > all (select 1)"; !strings.Contains(engine, want) {
		t.Errorf("MySQL no longer stores %q as %q: %s", expr, want, engine)
	}

	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE tt (a INT, b INT, c INT); CREATE VIEW v AS SELECT "+expr+" AS x FROM tt", nil); err != nil {
		t.Fatal(err)
	}
	if got, want := c.ShowCreateView("test", "v"), "(1 = (`tt`.`a` > all (select 1)))"; !strings.Contains(got, want) {
		t.Errorf("catalog stored %s, want it to contain %q", got, want)
	}
}
