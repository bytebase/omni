package catalog

import (
	"strings"
	"testing"
)

// viewRestoreCatalog runs a script and fails on any statement's error.
func viewRestoreCatalog(t *testing.T, script string) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec(script, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("statement %d: %v", r.Index, r.Error)
		}
	}
	return c
}

// TestViewBodyUsesSessionDatabase: an unqualified table in a view body is
// the session's, not the view's. On TiDB v8.5.5, USE src; CREATE VIEW
// other.v AS SELECT orders.id FROM orders stores
// SELECT `orders`.`id` AS `id` FROM `src`.`orders`.
func TestViewBodyUsesSessionDatabase(t *testing.T) {
	c := viewRestoreCatalog(t, `CREATE DATABASE src; CREATE DATABASE other; USE src;
CREATE TABLE orders (id INT);
CREATE VIEW other.v AS SELECT orders.id FROM orders;`)
	v := c.GetDatabase("other").Views["v"]
	if v == nil {
		t.Fatal("view other.v not created")
	}
	if want := "SELECT `orders`.`id` AS `id` FROM `src`.`orders`"; v.Definition != want {
		t.Errorf("Definition = %s, want %s", v.Definition, want)
	}
	if got := strings.Join(v.Columns, ","); got != "id" {
		t.Errorf("Columns = %s, want id", got)
	}
}

// TestViewCastIgnoresDisplayWidth: CAST(column AS CHAR) takes the width of
// the column's type, not its display width, which does not bound the value.
// On TiDB v8.5.5 (expression/builtin_cast.go, adjustRetFtForCastString):
// INT(1) ZEROFILL casts as CHAR(10), INT(3) UNSIGNED as CHAR(10), INT(4) as
// CHAR(11), BIGINT(5) as CHAR(20), TINYINT(1) as CHAR(4).
func TestViewCastIgnoresDisplayWidth(t *testing.T) {
	c := viewRestoreCatalog(t, `CREATE DATABASE d; USE d;
CREATE TABLE orders (z INT(1) ZEROFILL, u INT(3) UNSIGNED, s INT(4), bg BIGINT(5), tz TINYINT(1));
CREATE VIEW vc AS SELECT CAST(z AS CHAR), CAST(u AS CHAR), CAST(s AS CHAR), CAST(bg AS CHAR), CAST(tz AS CHAR) FROM orders;`)
	want := "SELECT CAST(`z` AS CHAR(10)) AS `CAST(z AS CHAR)`,CAST(`u` AS CHAR(10)) AS `CAST(u AS CHAR)`," +
		"CAST(`s` AS CHAR(11)) AS `CAST(s AS CHAR)`,CAST(`bg` AS CHAR(20)) AS `CAST(bg AS CHAR)`," +
		"CAST(`tz` AS CHAR(4)) AS `CAST(tz AS CHAR)` FROM `d`.`orders`"
	if got := c.GetDatabase("d").Views["vc"].Definition; got != want {
		t.Errorf("Definition =\n%s\nwant\n%s", got, want)
	}
}

// TestViewUnaryPlusLiteralName: SELECT +1 names its column 1, so a derived
// table over it has a column `1` the outer query can read, and the view keeps
// its analyzed query. TiDB v8.5.5 stores the body as
// SELECT `d`.`1` AS `1` FROM (SELECT +1 AS `1`) AS `d`.
func TestViewUnaryPlusLiteralName(t *testing.T) {
	c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE VIEW vp AS SELECT d.`1` FROM (SELECT +1) d;")
	v := c.GetDatabase("d").Views["vp"]
	if want := "SELECT `d`.`1` AS `1` FROM (SELECT +1 AS `1`) AS `d`"; v.Definition != want {
		t.Errorf("Definition = %s, want %s", v.Definition, want)
	}
	if v.AnalyzedQuery == nil {
		t.Error("AnalyzedQuery is nil: the analyzer did not name SELECT +1's column 1")
	}
}

// TestViewCastWidthFromQualifiedColumn: db.t.c names the table of that
// schema when a.t and b.t are both joined. On TiDB v8.5.5, with
// a.t (s VARCHAR(1)) and b.t (s VARCHAR(100)), the casts are CHAR(100) for
// b.t.s and CHAR(1) for a.t.s.
func TestViewCastWidthFromQualifiedColumn(t *testing.T) {
	c := viewRestoreCatalog(t, `CREATE DATABASE a; CREATE DATABASE b; USE a;
CREATE TABLE t (s VARCHAR(1));
CREATE TABLE b.t (s VARCHAR(100));
CREATE VIEW vj AS SELECT CAST(b.t.s AS CHAR), CAST(a.t.s AS CHAR) FROM a.t JOIN b.t;`)
	want := "SELECT CAST(`b`.`t`.`s` AS CHAR(100)) AS `CAST(b.t.s AS CHAR)`,CAST(`a`.`t`.`s` AS CHAR(1)) AS `CAST(a.t.s AS CHAR)` FROM `a`.`t` JOIN `b`.`t`"
	if got := c.GetDatabase("a").Views["vj"].Definition; got != want {
		t.Errorf("Definition =\n%s\nwant\n%s", got, want)
	}
}

// TestViewSchemaQualifiedTable: a table written with its schema is that
// schema's, not the session's table of the same name. On TiDB v8.5.5, with
// a.t (s VARCHAR(1), x INT) and b.t (s VARCHAR(100), y INT, z INT),
// USE a; CREATE VIEW vq AS SELECT CAST(s AS CHAR), b.* FROM b.t AS b stores
// the width and the columns of b.t.
func TestViewSchemaQualifiedTable(t *testing.T) {
	c := viewRestoreCatalog(t, `CREATE DATABASE a; CREATE DATABASE b; USE a;
CREATE TABLE t (s VARCHAR(1), x INT);
CREATE TABLE b.t (s VARCHAR(100), y INT, z INT);
CREATE VIEW vq AS SELECT CAST(s AS CHAR), b.* FROM b.t AS b;`)
	v := c.GetDatabase("a").Views["vq"]
	want := "SELECT CAST(`s` AS CHAR(100)) AS `CAST(s AS CHAR)`,`b`.`b`.`s` AS `s`,`b`.`b`.`y` AS `y`,`b`.`b`.`z` AS `z` FROM `b`.`t` AS `b`"
	if v.Definition != want {
		t.Errorf("Definition =\n%s\nwant\n%s", v.Definition, want)
	}
	if got := strings.Join(v.Columns, ","); got != "CAST(s AS CHAR),s,y,z" {
		t.Errorf("Columns = %s", got)
	}
}

// TestViewRestoreReviewForms pins three bodies against what TiDB v8.5.5
// stores for them (SHOW CREATE VIEW), and checks each stored body reads back
// through CREATE VIEW to the same text.
func TestViewRestoreReviewForms(t *testing.T) {
	for _, tc := range []struct {
		name, setup, body, want string
		cols                    []string
	}{
		{
			// A column a RIGHT JOIN ... USING joins on is the right
			// input's, and so is its width.
			"right join using", "CREATE TABLE t1 (s VARCHAR(1)); CREATE TABLE t2 (s VARCHAR(100))",
			"SELECT CAST(s AS CHAR) FROM t1 RIGHT JOIN t2 USING (s)",
			"SELECT CAST(`s` AS CHAR(100)) AS `CAST(s AS CHAR)` FROM `d`.`t1` RIGHT JOIN `d`.`t2` USING (`s`)",
			[]string{"CAST(s AS CHAR)"},
		},
		{
			// An empty string literal names its column "", in a derived
			// table and in a CTE alike.
			"empty name in a derived table", "",
			"SELECT * FROM (SELECT '') d",
			"SELECT `d`.`` FROM (SELECT _UTF8MB4'') AS `d`",
			[]string{""},
		},
		{
			"empty name in a CTE", "",
			"WITH c AS (SELECT '', 1) SELECT * FROM c",
			"WITH `c` AS (SELECT _UTF8MB4'',1 AS `1`) SELECT `c`.``,`c`.`1` AS `1` FROM `c`",
			[]string{"", "1"},
		},
		{
			// CONVERT ... USING names the charset as a lower-case string.
			"convert using", "",
			"SELECT CONVERT('x' USING UTF8MB4) AS a, CONVERT('x' USING 'Latin1') AS b",
			"SELECT CONVERT(_UTF8MB4'x' USING 'utf8mb4') AS `a`,CONVERT(_UTF8MB4'x' USING 'latin1') AS `b`",
			[]string{"a", "b"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setup := "CREATE DATABASE d; USE d;"
			if tc.setup != "" {
				setup += " " + tc.setup + ";"
			}
			c := viewRestoreCatalog(t, setup+" CREATE VIEW v AS "+tc.body+";")
			v := c.GetDatabase("d").Views["v"]
			if v.Definition != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", v.Definition, tc.want)
			}
			if strings.Join(v.Columns, "|") != strings.Join(tc.cols, "|") {
				t.Errorf("Columns = %q, want %q", v.Columns, tc.cols)
			}
			again := viewRestoreCatalog(t, setup+" CREATE VIEW v AS "+v.Definition+";")
			if got := again.GetDatabase("d").Views["v"].Definition; got != v.Definition {
				t.Errorf("reads back as\n%s\nwas\n%s", got, v.Definition)
			}
		})
	}
}
