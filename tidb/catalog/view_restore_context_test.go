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

// TestViewRestoreExprForms: the expression forms below are stored as TiDB
// v8.5.5 stores them (SHOW CREATE VIEW, 2026-11-06): N'...' as _UTF8'...',
// TRIM with a direction and no remstr with the space it removes, CHAR ...
// USING as CHAR_FUNC with the charset as a string, and an (a, b) row as ROW.
// TestDeparse_Section_7_3_ExpressionViews compares the same bodies with the
// server's.
func TestViewRestoreExprForms(t *testing.T) {
	for _, tc := range []struct {
		name, body, want string
		cols             []string
	}{
		{
			"national string",
			"SELECT N'nat', n'x' AS y, N'a' 'b', N'' FROM t",
			"SELECT _UTF8'nat' AS `nat`,_UTF8'x' AS `y`,_UTF8'a' AS `b`,_UTF8'' FROM `d`.`t`",
			[]string{"nat", "y", "b", ""},
		},
		{
			"trim direction",
			"SELECT TRIM(BOTH FROM c), TRIM(LEADING FROM c), TRIM(TRAILING 'x' FROM c) FROM t",
			"SELECT TRIM(BOTH _UTF8MB4' ' FROM `c`) AS `TRIM(BOTH FROM c)`,TRIM(LEADING _UTF8MB4' ' FROM `c`) AS `TRIM(LEADING FROM c)`," +
				"TRIM(TRAILING _UTF8MB4'x' FROM `c`) AS `TRIM(TRAILING 'x' FROM c)` FROM `d`.`t`",
			[]string{"TRIM(BOTH FROM c)", "TRIM(LEADING FROM c)", "TRIM(TRAILING 'x' FROM c)"},
		},
		{
			"char using",
			"SELECT CHAR(65 USING utf8mb4), CHAR(65, 66 USING 'Latin1') AS l, CHAR(65 USING utf8mb3) AS u, CHAR(65) FROM t",
			"SELECT CHAR_FUNC(65, 'utf8mb4') AS `CHAR(65 USING utf8mb4)`,CHAR_FUNC(65, 66, 'latin1') AS `l`," +
				"CHAR_FUNC(65, 'utf8') AS `u`,CHAR_FUNC(65, NULL) AS `CHAR(65)` FROM `d`.`t`",
			[]string{"CHAR(65 USING utf8mb4)", "l", "u", "CHAR(65)"},
		},
		{
			"row shorthand",
			"SELECT (a,b) IN ((1,2)), (a, b) = (1, 2) AS e FROM t",
			"SELECT ROW(`a`,`b`) IN (ROW(1,2)) AS `(a,b) IN ((1,2))`,ROW(`a`,`b`)=ROW(1,2) AS `e` FROM `d`.`t`",
			[]string{"(a,b) IN ((1,2))", "e"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT, c VARCHAR(20)); CREATE VIEW v AS "+tc.body+";")
			v := c.GetDatabase("d").Views["v"]
			if v.Definition != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", v.Definition, tc.want)
			}
			if strings.Join(v.Columns, "|") != strings.Join(tc.cols, "|") {
				t.Errorf("Columns = %q, want %q", v.Columns, tc.cols)
			}
		})
	}
}

// TestViewRestorePriority: a view stores its SELECT's priority and drops
// SQL_CACHE, as TiDB v8.5.5 stores them (SHOW CREATE VIEW, 2026-11-06). The last
// priority written is the one kept, and SQL_NO_CACHE beats SQL_CACHE in either
// order. TestDeparse_Section_7_3_ExpressionViews compares the same bodies with the
// server's.
func TestViewRestorePriority(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{"SELECT LOW_PRIORITY a FROM t", "SELECT LOW_PRIORITY `a` AS `a` FROM `d`.`t`"},
		{"SELECT HIGH_PRIORITY a FROM t", "SELECT HIGH_PRIORITY `a` AS `a` FROM `d`.`t`"},
		{"SELECT DELAYED a FROM t", "SELECT DELAYED `a` AS `a` FROM `d`.`t`"},
		{"SELECT HIGH_PRIORITY LOW_PRIORITY a FROM t", "SELECT LOW_PRIORITY `a` AS `a` FROM `d`.`t`"},
		{"SELECT SQL_CACHE a FROM t", "SELECT `a` AS `a` FROM `d`.`t`"},
		{"SELECT SQL_CACHE SQL_NO_CACHE a FROM t", "SELECT SQL_NO_CACHE `a` AS `a` FROM `d`.`t`"},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT); CREATE VIEW v AS "+tc.body+";")
			if got := c.GetDatabase("d").Views["v"].Definition; got != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}

// TestViewRestoreODBCEscapes: a view stores an ODBC escape as TiDB v8.5.5 stores
// it (SHOW CREATE VIEW, 2026-11-06). {d}, {t} and {ts} become DATE, TIME and
// TIMESTAMP literals named after the escape's text; any other name, D and TS
// among them, leaves the bare expression, here a string; and {oj ...} leaves its
// join. TiDB names an unaliased {abc 1} after the whole field, {abc 1}; omni names
// it 1, because the escape leaves no node of its own (parseODBCEscape).
func TestViewRestoreODBCEscapes(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{
			"SELECT {d '2020-01-01'}, {t '00:00:00'}, {ts '2020-01-01 00:00:00'}",
			"SELECT DATE '2020-01-01' AS `{d '2020-01-01'}`,TIME '00:00:00' AS `{t '00:00:00'}`," +
				"TIMESTAMP '2020-01-01 00:00:00' AS `{ts '2020-01-01 00:00:00'}`",
		},
		{
			"SELECT {D '2020-01-01'} AS x, {TS '2020-01-01 00:00:00'} AS y",
			"SELECT _UTF8MB4'2020-01-01' AS `x`,_UTF8MB4'2020-01-01 00:00:00' AS `y`",
		},
		{
			"SELECT b FROM t1 WHERE b > {d '2020-01-01'}",
			"SELECT `b` AS `b` FROM `d`.`t1` WHERE `b`>DATE '2020-01-01'",
		},
		{
			"SELECT t1.a FROM {oj t1 LEFT JOIN t2 ON t1.a = t2.a}",
			"SELECT `t1`.`a` AS `a` FROM `d`.`t1` LEFT JOIN `d`.`t2` ON `t1`.`a`=`t2`.`a`",
		},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE TABLE t1 (a INT, b DATE); CREATE TABLE t2 (a INT); CREATE VIEW v AS "+tc.body+";")
			if got := c.GetDatabase("d").Views["v"].Definition; got != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}

// TestViewRestoreAdjacentStrings: a view stores a run of adjacent strings as the
// joined string, named after the first, as TiDB v8.5.5 stores it (SHOW CREATE
// VIEW, 2026-11-06). The first string's leading spaces are dropped from the name,
// and an empty first string leaves the item unaliased, its column named "".
func TestViewRestoreAdjacentStrings(t *testing.T) {
	for _, tc := range []struct {
		body, want string
		cols       []string
	}{
		{"SELECT 'a' 'b'", "SELECT _UTF8MB4'ab' AS `a`", []string{"a"}},
		{
			"SELECT 'a' 'b' 'c', '' 'b', '  x' 'y'",
			"SELECT _UTF8MB4'abc' AS `a`,_UTF8MB4'b',_UTF8MB4'  xy' AS `x`",
			[]string{"a", "", "x"},
		},
		{"SELECT ('a' 'b')", "SELECT (_UTF8MB4'ab') AS `a`", []string{"a"}},
		{"SELECT 'a' 'b' AS c", "SELECT _UTF8MB4'ab' AS `c`", []string{"c"}},
		{
			"SELECT CONCAT('a' 'b', 'c')",
			"SELECT CONCAT(_UTF8MB4'ab', _UTF8MB4'c') AS `CONCAT('a' 'b', 'c')`",
			[]string{"CONCAT('a' 'b', 'c')"},
		},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE VIEW v AS "+tc.body+";")
			v := c.GetDatabase("d").Views["v"]
			if v.Definition != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", v.Definition, tc.want)
			}
			if strings.Join(v.Columns, "|") != strings.Join(tc.cols, "|") {
				t.Errorf("Columns = %q, want %q", v.Columns, tc.cols)
			}
		})
	}
	t.Run("column default", func(t *testing.T) {
		c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE TABLE td (a VARCHAR(10) DEFAULT 'a' 'b');")
		if got := c.ShowCreateTable("d", "td"); !strings.Contains(got, "DEFAULT 'ab'") {
			t.Errorf("SHOW CREATE TABLE = %s, want DEFAULT 'ab'", got)
		}
	})
}

// TestViewRestoreFetchFirst: a view stores FETCH FIRST as the LIMIT it is, as
// TiDB v8.5.5 stores it (SHOW CREATE VIEW, 2026-11-06), LIMIT 1 when the count is
// left out.
func TestViewRestoreFetchFirst(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{"SELECT a FROM t FETCH FIRST 5 ROWS ONLY", "SELECT `a` AS `a` FROM `d`.`t` LIMIT 5"},
		{"SELECT a FROM t ORDER BY a FETCH NEXT ROW ONLY", "SELECT `a` AS `a` FROM `d`.`t` ORDER BY `a` LIMIT 1"},
		{
			"SELECT a FROM t UNION SELECT a FROM t FETCH FIRST 2 ROWS ONLY",
			"SELECT `a` AS `a` FROM `d`.`t` UNION SELECT `a` AS `a` FROM `d`.`t` LIMIT 2",
		},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT); CREATE VIEW v AS "+tc.body+";")
			if got := c.GetDatabase("d").Views["v"].Definition; got != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}

// TestViewRestoreWindowNullTreatment: RESPECT NULLS and FROM FIRST are the
// defaults and a view stores neither, as TiDB v8.5.5 stores it (SHOW CREATE VIEW,
// 2026-11-07). TiDB parses IGNORE NULLS and FROM LAST but refuses them at view
// creation (1235), so only the default forms are stored on both sides.
func TestViewRestoreWindowNullTreatment(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{
			"SELECT FIRST_VALUE(val) RESPECT NULLS OVER (PARTITION BY p) AS x FROM t",
			"SELECT FIRST_VALUE(`val`) OVER (PARTITION BY `p`) AS `x` FROM `d`.`t`",
		},
		{
			"SELECT NTH_VALUE(val, 2) FROM FIRST OVER (ORDER BY p) AS x FROM t",
			"SELECT NTH_VALUE(`val`, 2) OVER (ORDER BY `p`) AS `x` FROM `d`.`t`",
		},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE TABLE t (val INT, p INT); CREATE VIEW v AS "+tc.body+";")
			if got := c.GetDatabase("d").Views["v"].Definition; got != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}

// TestViewRestoreIntroducerLiteral: a view stores a charset-introduced hex or bit
// literal as the plain literal, dropping the introducer, as TiDB v8.5.5 stores it
// (SHOW CREATE VIEW, 2026-11-07).
func TestViewRestoreIntroducerLiteral(t *testing.T) {
	for _, tc := range []struct{ body, want string }{
		{"SELECT _utf8 0x61 AS x", "SELECT x'61' AS `x`"},
		{"SELECT _latin1 B'101' AS y", "SELECT b'101' AS `y`"},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := viewRestoreCatalog(t, "CREATE DATABASE d; USE d; CREATE VIEW v AS "+tc.body+";")
			if got := c.GetDatabase("d").Views["v"].Definition; got != tc.want {
				t.Errorf("Definition =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}
