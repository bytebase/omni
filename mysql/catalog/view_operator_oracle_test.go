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
