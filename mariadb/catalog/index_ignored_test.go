package catalog

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"testing"
)

// ignoredIndexCases build table t with MariaDB's index ignorability. ignored is
// INFORMATION_SCHEMA.STATISTICS.IGNORED per index and lines are the index
// lines of SHOW CREATE TABLE t, both as MariaDB 11.8.8 answers them, which
// TestContainerIgnoredIndex holds. VISIBLE is MariaDB's no-op, so it neither
// shows an ignored index again nor appears in the readback, and
// /*!80000 INVISIBLE */ is skipped by its MySQL version.
var ignoredIndexCases = []struct {
	name    string
	stmts   []string
	ignored map[string]bool
	lines   []string
}{
	{
		name: "every kind of index in CREATE TABLE",
		stmts: []string{
			"CREATE TABLE t (a INT, b INT, c TEXT, g GEOMETRY NOT NULL," +
				" KEY i (a) IGNORED, UNIQUE KEY u (b) COMMENT 'u note' IGNORED, KEY n (a, b) NOT IGNORED," +
				" KEY v (b) VISIBLE, KEY iv (a) IGNORED VISIBLE, KEY ni (b) NOT IGNORED IGNORED," +
				" FULLTEXT KEY f (c) IGNORED, SPATIAL KEY s (g) COMMENT 'sp' IGNORED, KEY my (a) /*!80000 INVISIBLE */)",
		},
		ignored: map[string]bool{"i": true, "u": true, "n": false, "v": false, "iv": true, "ni": true, "f": true, "s": true, "my": false},
		lines: []string{
			"KEY `i` (`a`) IGNORED",
			"UNIQUE KEY `u` (`b`) COMMENT 'u note' IGNORED",
			"KEY `n` (`a`,`b`)",
			"KEY `v` (`b`)",
			"KEY `iv` (`a`) IGNORED",
			"KEY `ni` (`b`) IGNORED",
			"FULLTEXT KEY `f` (`c`) IGNORED",
			"SPATIAL KEY `s` (`g`) COMMENT 'sp' IGNORED",
			"KEY `my` (`a`)",
		},
	},
	{
		name: "CREATE INDEX, ALTER TABLE ADD and ALTER INDEX",
		stmts: []string{
			"CREATE TABLE t (a INT, b INT)",
			"CREATE INDEX ci ON t (a) IGNORED",
			"CREATE UNIQUE INDEX cu ON t (b) COMMENT 'c' IGNORED",
			"ALTER TABLE t ADD INDEX ai (b) IGNORED, ADD KEY an (a) NOT IGNORED",
			"ALTER TABLE t ALTER INDEX ci NOT IGNORED",
			"ALTER TABLE t ALTER KEY an IGNORED",
			"ALTER TABLE t ALTER INDEX IF EXISTS nosuch IGNORED",
			"ALTER TABLE t RENAME INDEX ai TO ar",
		},
		ignored: map[string]bool{"ci": false, "cu": true, "ar": true, "an": true},
		lines: []string{
			"UNIQUE KEY `cu` (`b`) COMMENT 'c' IGNORED",
			"KEY `ci` (`a`)",
			"KEY `ar` (`b`) IGNORED",
			"KEY `an` (`a`) IGNORED",
		},
	},
	{
		// IF EXISTS asks about the table the statement found, so an index the
		// same statement adds or renames to the name is left as it is.
		name: "ALTER INDEX IF EXISTS of an index the statement makes",
		stmts: []string{
			"CREATE TABLE t (a INT, b INT, KEY k (a))",
			"ALTER TABLE t ADD INDEX n (b), ALTER INDEX IF EXISTS n IGNORED, ALTER INDEX k IGNORED",
			"ALTER TABLE t RENAME INDEX n TO n2, ALTER INDEX IF EXISTS n2 IGNORED",
		},
		ignored: map[string]bool{"k": true, "n2": false},
		lines:   []string{"KEY `k` (`a`) IGNORED", "KEY `n2` (`b`)"},
	},
	{
		name:    "the last ALTER INDEX of a statement wins",
		stmts:   []string{"CREATE TABLE t (a INT, KEY k (a))", "ALTER TABLE t ALTER INDEX k NOT IGNORED, ALTER INDEX k IGNORED"},
		ignored: map[string]bool{"k": true},
		lines:   []string{"KEY `k` (`a`) IGNORED"},
	},
	{
		name:    "CREATE TABLE LIKE keeps it",
		stmts:   []string{"CREATE TABLE src (a INT, KEY i (a) COMMENT 'c' IGNORED)", "CREATE TABLE t LIKE src"},
		ignored: map[string]bool{"i": true},
		lines:   []string{"KEY `i` (`a`) COMMENT 'c' IGNORED"},
	},
	{
		name:    "a primary key written NOT IGNORED and VISIBLE",
		stmts:   []string{"CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (a) NOT IGNORED VISIBLE, KEY k (b) IGNORED)"},
		ignored: map[string]bool{"PRIMARY": false, "k": true},
		lines:   []string{"PRIMARY KEY (`a`)", "KEY `k` (`b`) IGNORED"},
	},
	{
		// The implicit primary key is the first UNIQUE key of whole NOT NULL
		// columns; an IGNORED one after it, or one with a nullable or prefix
		// part, or a long unique hash, is not promoted.
		name: "unique keys the implicit primary key leaves alone",
		stmts: []string{
			"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, d VARCHAR(10) NOT NULL," +
				" UNIQUE KEY v (b), UNIQUE KEY u (a) IGNORED, UNIQUE KEY n (a, c) IGNORED," +
				" UNIQUE KEY p (d(5)) IGNORED, UNIQUE KEY h (a) USING HASH IGNORED)",
		},
		ignored: map[string]bool{"v": false, "u": true, "n": true, "p": true, "h": true},
		lines: []string{
			"UNIQUE KEY `v` (`b`)",
			"UNIQUE KEY `u` (`a`) IGNORED",
			"UNIQUE KEY `n` (`a`,`c`) IGNORED",
			"UNIQUE KEY `p` (`d`(5)) IGNORED",
			"UNIQUE KEY `h` (`a`) USING HASH IGNORED",
		},
	},
	{
		// A column-level UNIQUE written first is promoted.
		name:    "a column-level unique key written first",
		stmts:   []string{"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL UNIQUE, UNIQUE KEY u (a) IGNORED)"},
		ignored: map[string]bool{"b": false, "u": true},
		lines:   []string{"UNIQUE KEY `b` (`b`)", "UNIQUE KEY `u` (`a`) IGNORED"},
	},
	{
		// MariaDB keeps w ahead of u from the CREATE, when u's column was still
		// nullable, so making it NOT NULL promotes nothing new. The catalog does
		// not track that order and reports nothing either.
		name:    "a key made NOT NULL after the implicit primary key was chosen",
		stmts:   []string{"CREATE TABLE t (a INT, b INT NOT NULL, UNIQUE KEY u (a) IGNORED, UNIQUE KEY w (b))", "ALTER TABLE t MODIFY a INT NOT NULL"},
		ignored: map[string]bool{"u": true, "w": false},
		lines:   []string{"UNIQUE KEY `w` (`b`)", "UNIQUE KEY `u` (`a`) IGNORED"},
	},
	{
		name:    "an IGNORED unique key behind another",
		stmts:   []string{"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY w (a))", "ALTER TABLE t ALTER INDEX w IGNORED", "ALTER TABLE t ADD UNIQUE KEY x (a, b) IGNORED"},
		ignored: map[string]bool{"v": false, "w": true, "x": true},
		lines:   []string{"UNIQUE KEY `v` (`b`)", "UNIQUE KEY `w` (`a`) IGNORED", "UNIQUE KEY `x` (`a`,`b`) IGNORED"},
	},
}

// ignoredIndexErrorCases end in a statement MariaDB 11.8.8 answers with code
// (0 for none), held by TestContainerIgnoredIndexErrors. 4174 refuses an
// IGNORED primary key, declared or implicit, and any ALTER INDEX of the primary
// key; 1176 is an ALTER INDEX naming no index of the table as the statement
// found it.
var ignoredIndexErrorCases = []struct {
	name  string
	setup []string
	stmt  string
	code  int
}{
	{name: "an ignored primary key", stmt: "CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) IGNORED)", code: 4174},
	{name: "an ignored named primary key", stmt: "CREATE TABLE t (a INT NOT NULL, CONSTRAINT pk PRIMARY KEY (a) IGNORED)", code: 4174},
	{name: "an ignored implicit primary key", stmt: "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) IGNORED)", code: 4174},
	{name: "the first unique key ignored", stmt: "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) IGNORED, UNIQUE KEY v (b))", code: 4174},
	{name: "the first NOT NULL unique key ignored", stmt: "CREATE TABLE t (a INT, b INT NOT NULL, UNIQUE KEY u (a), UNIQUE KEY w (b) IGNORED)", code: 4174},
	{name: "AUTO_INCREMENT makes the column NOT NULL", stmt: "CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) IGNORED)", code: 4174},
	{name: "a descending part", stmt: "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a DESC) IGNORED)", code: 4174},
	{name: "the second unique key ignored", stmt: "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) IGNORED)"},
	{name: "a nullable part", stmt: "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) IGNORED)"},
	{name: "a prefix part", stmt: "CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(5)) IGNORED)"},
	{name: "a long unique hash", stmt: "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) USING HASH IGNORED)"},
	{name: "a declared primary key besides", stmt: "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) IGNORED, b INT NOT NULL, PRIMARY KEY (b))"},

	{name: "ADD PRIMARY KEY IGNORED", setup: []string{"CREATE TABLE t (a INT)"}, stmt: "ALTER TABLE t ADD PRIMARY KEY (a) IGNORED", code: 4174},
	{name: "ALTER INDEX PRIMARY NOT IGNORED", setup: []string{"CREATE TABLE t (a INT NOT NULL PRIMARY KEY)"}, stmt: "ALTER TABLE t ALTER INDEX `PRIMARY` NOT IGNORED", code: 4174},
	{name: "ALTER INDEX IF EXISTS PRIMARY", setup: []string{"CREATE TABLE t (a INT NOT NULL PRIMARY KEY)"}, stmt: "ALTER TABLE t ALTER INDEX IF EXISTS `PRIMARY` IGNORED", code: 4174},
	{name: "ALTER INDEX of the implicit primary key", setup: []string{"CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a))"}, stmt: "ALTER TABLE t ALTER INDEX u IGNORED", code: 4174},
	{name: "ALTER INDEX NOT IGNORED of the implicit primary key", setup: []string{"CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a))"}, stmt: "ALTER TABLE t ALTER INDEX u NOT IGNORED", code: 4174},
	{name: "ALTER INDEX PRIMARY without one", setup: []string{"CREATE TABLE t (a INT)"}, stmt: "ALTER TABLE t ALTER INDEX `PRIMARY` IGNORED", code: 1176},
	{name: "ALTER INDEX of no index", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t ALTER INDEX nosuch IGNORED", code: 1176},
	{name: "ALTER INDEX IF EXISTS of no index", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t ALTER INDEX IF EXISTS nosuch IGNORED"},
	{name: "ALTER INDEX then DROP INDEX", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t ALTER INDEX k IGNORED, DROP INDEX k", code: 1176},
	{name: "DROP INDEX then ALTER INDEX", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t DROP INDEX k, ALTER INDEX k IGNORED", code: 1176},
	{name: "ALTER INDEX then RENAME INDEX", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t ALTER INDEX k IGNORED, RENAME INDEX k TO k2", code: 1176},
	{name: "RENAME INDEX then ALTER INDEX", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t RENAME INDEX k TO k2, ALTER INDEX k2 IGNORED", code: 1176},
	{name: "ADD INDEX then ALTER INDEX", setup: []string{"CREATE TABLE t (a INT)"}, stmt: "ALTER TABLE t ADD INDEX n (a), ALTER INDEX n IGNORED", code: 1176},
	{name: "ALTER INDEX IF EXISTS then DROP INDEX", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t ALTER INDEX IF EXISTS k IGNORED, DROP INDEX k", code: 1176},
	{name: "ALTER INDEX IF EXISTS then RENAME INDEX", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t ALTER INDEX IF EXISTS k IGNORED, RENAME INDEX k TO k2", code: 1176},
	{name: "an index dropped and added again", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t DROP INDEX k, ADD INDEX k (a), ALTER INDEX k IGNORED", code: 1176},
	{name: "an index dropped and added again, IF EXISTS", setup: []string{"CREATE TABLE t (a INT, KEY k (a))"}, stmt: "ALTER TABLE t DROP INDEX k, ADD INDEX k (a), ALTER INDEX IF EXISTS k IGNORED", code: 1176},

	{name: "DROP PRIMARY KEY promotes an ignored key", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a), UNIQUE KEY u (b) IGNORED)"}, stmt: "ALTER TABLE t DROP PRIMARY KEY", code: 4174},
	{name: "MODIFY NOT NULL promotes an ignored key", setup: []string{"CREATE TABLE t (a INT, UNIQUE KEY u (a) IGNORED)"}, stmt: "ALTER TABLE t MODIFY a INT NOT NULL", code: 4174},
	{name: "ALTER TABLE DROP INDEX promotes an ignored key", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY w (a) IGNORED)"}, stmt: "ALTER TABLE t DROP INDEX v", code: 4174},
	{name: "DROP INDEX promotes an ignored key", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY w (a) IGNORED)"}, stmt: "DROP INDEX v ON t", code: 4174},
	{name: "CREATE UNIQUE INDEX IGNORED as the implicit primary key", setup: []string{"CREATE TABLE t (a INT NOT NULL)"}, stmt: "CREATE UNIQUE INDEX u ON t (a) IGNORED", code: 4174},
	{name: "ADD UNIQUE IGNORED as the implicit primary key", setup: []string{"CREATE TABLE t (a INT NOT NULL)"}, stmt: "ALTER TABLE t ADD UNIQUE KEY u (a) IGNORED", code: 4174},
	{name: "ADD UNIQUE IGNORED behind another", setup: []string{"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b))"}, stmt: "ALTER TABLE t ADD UNIQUE KEY u (a) IGNORED"},
}

// indexLines returns the index lines of a SHOW CREATE TABLE text, sorted.
func indexLines(show string) []string {
	var out []string
	for _, line := range strings.Split(show, "\n") {
		line = strings.TrimSuffix(strings.TrimSpace(line), ",")
		for _, p := range []string{"PRIMARY KEY ", "UNIQUE KEY ", "KEY ", "FULLTEXT KEY ", "SPATIAL KEY "} {
			if strings.HasPrefix(line, p) {
				out = append(out, line)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

func sortedCopy(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

func catalogIgnored(c *Catalog, table string) map[string]bool {
	got := map[string]bool{}
	for _, idx := range c.GetDatabase("testdb").GetTable(table).Indexes {
		got[idx.Name] = !idx.Visible
	}
	return got
}

func equalIgnored(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || w != v {
			return false
		}
	}
	return true
}

func equalLines(a, b []string) bool {
	return strings.Join(a, "\n") == strings.Join(b, "\n")
}

// TestIgnoredIndexState pins the catalog's ignorability and SHOW CREATE TABLE
// index lines after each case.
func TestIgnoredIndexState(t *testing.T) {
	for _, tc := range ignoredIndexCases {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			for _, stmt := range tc.stmts {
				mustExec(t, c, stmt)
			}
			if got := catalogIgnored(c, "t"); !equalIgnored(got, tc.ignored) {
				t.Errorf("ignored %v, want %v", got, tc.ignored)
			}
			if got, want := indexLines(c.ShowCreateTable("testdb", "t")), sortedCopy(tc.lines); !equalLines(got, want) {
				t.Errorf("index lines\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
		})
	}
}

// TestIgnoredIndexErrors pins the catalog's answer to each case, and that a
// refused statement leaves the table as it was.
func TestIgnoredIndexErrors(t *testing.T) {
	for _, tc := range ignoredIndexErrorCases {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			for _, stmt := range tc.setup {
				mustExec(t, c, stmt)
			}
			before := ""
			if len(tc.setup) > 0 {
				before = c.ShowCreateTable("testdb", "t")
			}
			code, msg := catalogErr(t, c, tc.stmt)
			if code != tc.code {
				t.Fatalf("catalog answers %d %q, want %d", code, msg, tc.code)
			}
			if code == 0 {
				return
			}
			if len(tc.setup) == 0 {
				if c.GetDatabase("testdb").GetTable("t") != nil {
					t.Errorf("a refused CREATE TABLE left table t")
				}
				return
			}
			if after := c.ShowCreateTable("testdb", "t"); after != before {
				t.Errorf("a refused statement changed the table:\n%s\nwas\n%s", after, before)
			}
		})
	}
}

// TestIgnoredIndexAddPrimaryKeyLeavesColumns pins that a refused ADD PRIMARY
// KEY ... IGNORED does not make its columns NOT NULL: MariaDB 11.8.8 reads back
// `a` int(11) DEFAULT NULL (TestContainerIgnoredIndexErrors).
func TestIgnoredIndexAddPrimaryKeyLeavesColumns(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE t (a INT)")
	if code, _ := catalogErr(t, c, "ALTER TABLE t ADD PRIMARY KEY (a) IGNORED"); code != 4174 {
		t.Fatalf("code %d, want 4174", code)
	}
	if !c.GetDatabase("testdb").GetTable("t").GetColumn("a").Nullable {
		t.Errorf("a refused ADD PRIMARY KEY made a NOT NULL")
	}
}

// resetIgnoredIndexTables drops the tables the cases build.
func resetIgnoredIndexTables(t *testing.T, ctx context.Context, conn *sql.Conn) {
	t.Helper()
	if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS t, src"); err != nil {
		t.Fatalf("reset: %v", err)
	}
}

func showCreateOn(t *testing.T, ctx context.Context, conn *sql.Conn, table string) string {
	t.Helper()
	var name, show string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE "+table).Scan(&name, &show); err != nil {
		t.Fatalf("SHOW CREATE TABLE %s: %v", table, err)
	}
	return show
}

// TestContainerIgnoredIndex runs each case on MariaDB 11.8.8 and holds the
// pinned ignorability and index lines to what the server answers, then loads
// the server's SHOW CREATE TABLE into a fresh catalog, which must read the same
// ignorability back (container test).
func TestContainerIgnoredIndex(t *testing.T) {
	db := mariadbDatabase(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("conn: %v", err)
	}
	defer conn.Close()
	for _, tc := range ignoredIndexCases {
		t.Run(tc.name, func(t *testing.T) {
			resetIgnoredIndexTables(t, ctx, conn)
			for _, stmt := range tc.stmts {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("MariaDB exec %q: %v", stmt, err)
				}
			}
			rows, err := conn.QueryContext(ctx, "SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'")
			if err != nil {
				t.Fatalf("STATISTICS: %v", err)
			}
			engine := map[string]bool{}
			for rows.Next() {
				var name, ignored string
				if err := rows.Scan(&name, &ignored); err != nil {
					t.Fatalf("scan: %v", err)
				}
				engine[name] = ignored == "YES"
			}
			if err := rows.Close(); err != nil {
				t.Fatalf("rows: %v", err)
			}
			if !equalIgnored(engine, tc.ignored) {
				t.Errorf("MariaDB ignored %v, case pins %v", engine, tc.ignored)
			}
			show := showCreateOn(t, ctx, conn, "t")
			if got, want := indexLines(show), sortedCopy(tc.lines); !equalLines(got, want) {
				t.Errorf("MariaDB index lines\n%s\ncase pins\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
			readback := setupWithDB(t)
			mustExec(t, readback, show)
			if got := catalogIgnored(readback, "t"); !equalIgnored(got, tc.ignored) {
				t.Errorf("MariaDB's readback loads as ignored %v, want %v:\n%s", got, tc.ignored, show)
			}
		})
	}
}

// TestContainerIgnoredIndexErrors asks MariaDB 11.8.8 for each case's code,
// and that a refused statement leaves the table as SHOW CREATE TABLE had it
// (container test).
func TestContainerIgnoredIndexErrors(t *testing.T) {
	db := mariadbDatabase(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("conn: %v", err)
	}
	defer conn.Close()
	for _, tc := range ignoredIndexErrorCases {
		t.Run(tc.name, func(t *testing.T) {
			resetIgnoredIndexTables(t, ctx, conn)
			for _, stmt := range tc.setup {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("MariaDB exec %q: %v", stmt, err)
				}
			}
			before := ""
			if len(tc.setup) > 0 {
				before = showCreateOn(t, ctx, conn, "t")
			}
			_, err := conn.ExecContext(ctx, tc.stmt)
			code, msg := mariadbErr(t, tc.stmt, err)
			if code != tc.code {
				t.Fatalf("MariaDB answers %d %q, case pins %d", code, msg, tc.code)
			}
			if code != 0 && before != "" {
				if after := showCreateOn(t, ctx, conn, "t"); after != before {
					t.Errorf("MariaDB changed the table on a refused statement:\n%s\nwas\n%s", after, before)
				}
			}
		})
	}
}
