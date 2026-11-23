package catalog

import (
	"strings"
	"testing"
)

// A USING alias names the join's merged columns, and only those; a view
// keeps it, and a merged column it names stays unambiguous.
func TestJoinUsingAlias(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE a (id int, x int); CREATE TABLE b (id int, y int); CREATE TABLE c (id int, z int);
CREATE VIEW full_j AS SELECT j.id FROM a FULL JOIN b USING (id) AS j CROSS JOIN c;
CREATE VIEW inner_j AS SELECT j.* FROM a JOIN b USING (id) AS j;`)
	if err != nil {
		t.Fatal(err)
	}
	for view, want := range map[string]string{
		"full_j":  "SELECT j.id FROM ((a FULL JOIN b USING (id) AS j) CROSS JOIN c);",
		"inner_j": "SELECT a.id FROM (a JOIN b USING (id) AS j);",
	} {
		got, err := c.GetViewDefinition("public", view)
		if err != nil {
			t.Fatal(err)
		}
		if got = strings.Join(strings.Fields(got), " "); got != want {
			t.Errorf("%s: %q, want %q", view, got, want)
		}
		if _, err := LoadSQL("CREATE TABLE a (id int, x int); CREATE TABLE b (id int, y int); CREATE TABLE c (id int, z int); CREATE VIEW v AS " + got); err != nil {
			t.Errorf("%s's definition does not load again: %v", view, err)
		}
	}
	if _, err := LoadSQL("CREATE TABLE a (id int, x int); CREATE TABLE b (id int, y int); CREATE VIEW v AS SELECT j.x FROM a JOIN b USING (id) AS j"); err == nil {
		t.Errorf("j.x resolved, but j names only the merged id")
	}
}

// A query's FROM items have system columns; an expression analyzed on its
// own, as an index's or a CHECK constraint's, does not.
func TestSystemColumns(t *testing.T) {
	c, err := LoadSQL("CREATE TABLE a (id int); CREATE VIEW v AS SELECT ctid, (SELECT max(xmin::text) FROM a b) AS m FROM a")
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.GetViewDefinition("public", "v")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.Join(strings.Fields(got), " "), "SELECT ctid, ( SELECT max((b.xmin)::text) AS max FROM a b) AS m FROM a;"; got != want {
		t.Errorf("%q, want %q", got, want)
	}
	for _, sql := range []string{
		"CREATE TABLE a (id int); CREATE VIEW v AS SELECT ctid FROM (SELECT 1 AS id) s",
		"CREATE TABLE a (id int); CREATE VIEW w AS SELECT id FROM a; CREATE VIEW v AS SELECT ctid FROM w",
		"CREATE TABLE a (id int); CREATE INDEX ON a ((ctid IS NOT NULL))",
		"CREATE TABLE a (id int); CREATE INDEX ON a (id) WHERE ctid IS NOT NULL",
		"CREATE TABLE a (id int, CHECK (ctid IS NOT NULL))",
	} {
		if _, err := LoadSQL(sql); err == nil {
			t.Errorf("accepted %s", sql)
		}
	}
}
