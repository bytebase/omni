package catalog

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// These hold what the VECTOR type and a vector key leave in the catalog, not that they
// parse. Their expected text is what SHOW CREATE TABLE printed on mariadb:11.8.8 (read on
// 2026-11-05); TestVectorShowCreateOracle holds it to the server.

const vectorShowCreateDDL = `CREATE TABLE m (
  id int(11) NOT NULL,
  v vector(3) NOT NULL,
  w VECTOR(3,2) DEFAULT NULL,
  x VECTOR(1.5) DEFAULT NULL,
  a int(11) DEFAULT NULL,
  b text DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk (a),
  KEY ka (a),
  FULLTEXT KEY ft (b),
  VECTOR KEY k (v) M=6 DISTANCE=cosine COMMENT 'c'
)`

// vectorShowCreateWant is what the DDL's keys and columns read back as, in order.
var vectorShowCreateWant = []string{
	"`v` vector(3) NOT NULL",
	"`w` vector(3) DEFAULT NULL",
	"`x` vector(1) DEFAULT NULL",
	"PRIMARY KEY (`id`)",
	"UNIQUE KEY `uk` (`a`)",
	"KEY `ka` (`a`)",
	"FULLTEXT KEY `ft` (`b`)",
	"VECTOR KEY `k` (`v`) COMMENT 'c' `M`=6 `DISTANCE`=cosine",
}

func requireInOrder(t *testing.T, ddl string, want []string) {
	t.Helper()
	at := 0
	for _, w := range want {
		i := strings.Index(ddl[at:], w)
		if i < 0 {
			t.Errorf("SHOW CREATE TABLE lacks %s after byte %d:\n%s", w, at, ddl)
			return
		}
		at += i + len(w)
	}
}

func TestVectorShowCreate(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; "+vectorShowCreateDDL+";")
	requireInOrder(t, c.ShowCreateTable("d", "m"), vectorShowCreateWant)
	if idx := c.databases["d"].Tables["m"].Indexes[4]; !idx.Vector || idx.IndexType != "VECTOR" || idx.Comment != "c" {
		t.Errorf("the vector key: %+v", idx)
	}
}

// TestVectorKeySpellings: every way to add one leaves the same key, an unnamed one takes
// its column's name, and DEFAULT leaves the option out.
func TestVectorKeySpellings(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR (v))", "VECTOR KEY `v` (`v`)"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR k (v))", "VECTOR KEY `k` (`v`)"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX k (v) DISTANCE='cosine')", "VECTOR KEY `k` (`v`) `DISTANCE`='cosine'"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY k (v) M=default)", "VECTOR KEY `k` (`v`)"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL); CREATE VECTOR INDEX k ON t (v) distance=euclidean", "VECTOR KEY `k` (`v`) `distance`=euclidean"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL); ALTER TABLE t ADD VECTOR k (v) M=8", "VECTOR KEY `k` (`v`) `M`=8"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL); ALTER TABLE t ADD VECTOR (v)", "VECTOR KEY `v` (`v`)"},
	} {
		c := New()
		execAlterClauses(t, c, "CREATE DATABASE d; USE d; "+tc.sql+";")
		ddl := c.ShowCreateTable("d", "t")
		if !strings.Contains(ddl, tc.want+"\n") {
			t.Errorf("%s: SHOW CREATE TABLE lacks %s:\n%s", tc.sql, tc.want, ddl)
		}
	}
}

// TestVectorKeyLike: CREATE TABLE LIKE copies a vector key as one, options and all.
func TestVectorKeyLike(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY k (v) M=6); CREATE TABLE u LIKE t;")
	if ddl := c.ShowCreateTable("d", "u"); !strings.Contains(ddl, "VECTOR KEY `k` (`v`) `M`=6") {
		t.Errorf("LIKE lost the vector key:\n%s", ddl)
	}
}

// TestVectorKeyRefusals: what 11.8.8 refuses when the statement runs.
func TestVectorKeyRefusals(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"CREATE TABLE t2 (v VECTOR(3), VECTOR KEY (v));", 1252},
		{"CREATE TABLE t2 (v INT NOT NULL, VECTOR KEY (v));", 1210},
		{"CREATE TABLE t2 (v VECTOR(3) NOT NULL, w VECTOR(3) NOT NULL, VECTOR KEY (v), VECTOR KEY (w));", 1235},
		{"CREATE VECTOR INDEX i2 ON u (a);", 1210},
		{"CREATE VECTOR INDEX i2 ON u (n);", 1252},
		{"CREATE VECTOR INDEX i2 ON t (v);", 1235},
		{"ALTER TABLE t ADD VECTOR i2 (v);", 1235},
		{"ALTER TABLE t ADD VECTOR i (v);", 1061},
	} {
		c := New()
		execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT NOT NULL, n VECTOR(3), v VECTOR(3) NOT NULL, VECTOR KEY i (v)); CREATE TABLE u (a INT NOT NULL, n VECTOR(3));")
		err := execAlterClausesExpectError(t, c, tc.sql)
		var cerr *Error
		if !errors.As(err, &cerr) || cerr.Code != tc.code {
			t.Errorf("%s: error %v, want %d", tc.sql, err, tc.code)
		}
	}

	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v));")
	execAlterClauses(t, c, "ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v);")
	if n := len(c.databases["d"].Tables["t"].Indexes); n != 1 {
		t.Errorf("IF NOT EXISTS left %d indexes", n)
	}
	execAlterClauses(t, c, "DROP INDEX i ON t;")
	if n := len(c.databases["d"].Tables["t"].Indexes); n != 0 {
		t.Errorf("DROP INDEX left %d indexes", n)
	}
}

// TestVectorShowCreateOracle holds the expected text to the server.
func TestVectorShowCreateOracle(t *testing.T) {
	ctx := context.Background()
	conn := mariadbDatabase(t)
	if _, err := conn.ExecContext(ctx, vectorShowCreateDDL); err != nil {
		t.Fatal(err)
	}
	var name, got string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE m").Scan(&name, &got); err != nil {
		t.Fatal(err)
	}
	requireInOrder(t, got, vectorShowCreateWant)
}

// vectorStateCases run on the catalog and on the server: what the last statement answers, and
// the keys the table has after it, which on a refusal are the keys it had before.
var vectorStateCases = []struct {
	name  string
	setup string
	stmt  string
	code  int
	keys  []string
}{
	// The last setting of an engine option wins, DEFAULT clearing it.
	{name: "default clears", stmt: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 M=DEFAULT)", keys: []string{"VECTOR KEY `i` (`v`)\n"}},
	{name: "last wins", stmt: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 M=7)", keys: []string{"VECTOR KEY `i` (`v`) `M`=7\n"}},
	{name: "last moves", stmt: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=cosine M=7)", keys: []string{"VECTOR KEY `i` (`v`) `DISTANCE`=cosine `M`=7\n"}},
	{name: "any case", stmt: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) m=6 M=DEFAULT DISTANCE=cosine)", keys: []string{"VECTOR KEY `i` (`v`) `DISTANCE`=cosine\n"}},

	// CREATE OR REPLACE.
	{
		name:  "replaces a vector index",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, a INT, VECTOR INDEX i (v) M=6, KEY k (a))",
		stmt:  "CREATE OR REPLACE VECTOR INDEX i ON t (v) M=8",
		keys:  []string{"KEY `k` (`a`)", "VECTOR KEY `i` (`v`) `M`=8\n"},
	},
	{
		name:  "a second name is still a second vector index",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)",
		stmt:  "CREATE OR REPLACE VECTOR INDEX j ON t (v) M=9",
		code:  1235,
		keys:  []string{"VECTOR KEY `i` (`v`) `M`=6\n"},
	},
	{
		name:  "plain over vector",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, a INT, VECTOR INDEX i (v))",
		stmt:  "CREATE OR REPLACE INDEX i ON t (a)",
		keys:  []string{"KEY `i` (`a`)\n"},
	},
	{
		name:  "vector over plain",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, a INT, KEY k (a))",
		stmt:  "CREATE OR REPLACE VECTOR INDEX k ON t (v)",
		keys:  []string{"VECTOR KEY `k` (`v`)\n"},
	},
	{
		name:  "a plain replacement goes last among the plain keys",
		setup: "CREATE TABLE t (a INT, b INT, KEY a (a), KEY b (b))",
		stmt:  "CREATE OR REPLACE INDEX a ON t (b)",
		keys:  []string{"KEY `b` (`b`)", "KEY `a` (`b`)\n"},
	},
	{
		name:  "a unique key replaced by a plain one",
		setup: "CREATE TABLE t (a INT, UNIQUE KEY u (a))",
		stmt:  "CREATE OR REPLACE INDEX u ON t (a)",
		keys:  []string{"KEY `u` (`a`)\n"},
	},
	{
		name:  "a failed replacement keeps the old key: missing column",
		setup: "CREATE TABLE t (a INT, KEY k (a) COMMENT 'c')",
		stmt:  "CREATE OR REPLACE INDEX k ON t (nosuch)",
		code:  1072,
		keys:  []string{"KEY `k` (`a`) COMMENT 'c'\n"},
	},
	{
		name:  "a failed replacement keeps the old key: no vector column",
		setup: "CREATE TABLE t (a INT NOT NULL, KEY k (a) COMMENT 'c')",
		stmt:  "CREATE OR REPLACE VECTOR INDEX k ON t (a)",
		code:  1210,
		keys:  []string{"KEY `k` (`a`) COMMENT 'c'\n"},
	},

	// A vector key is held to the table the ALTER leaves.
	{
		name:  "a column added after it",
		setup: "CREATE TABLE t (id INT)",
		stmt:  "ALTER TABLE t ADD VECTOR i (v), ADD COLUMN v VECTOR(3) NOT NULL",
		keys:  []string{"VECTOR KEY `i` (`v`)\n"},
	},
	{
		name:  "made NOT NULL after it",
		setup: "CREATE TABLE t (id INT, v VECTOR(3))",
		stmt:  "ALTER TABLE t ADD VECTOR i (v), MODIFY v VECTOR(3) NOT NULL",
		keys:  []string{"VECTOR KEY `i` (`v`)\n"},
	},
	{
		name:  "dropped after it",
		setup: "CREATE TABLE t (id INT, v VECTOR(3) NOT NULL)",
		stmt:  "ALTER TABLE t ADD VECTOR i (v), DROP COLUMN v",
		code:  1072,
	},
	{
		name:  "renamed after it",
		setup: "CREATE TABLE t (id INT, w VECTOR(3) NOT NULL)",
		stmt:  "ALTER TABLE t ADD VECTOR i (w), CHANGE w z VECTOR(3) NOT NULL",
		code:  1072,
	},
	{
		name:  "added after it, the wrong type",
		setup: "CREATE TABLE t (id INT)",
		stmt:  "ALTER TABLE t ADD VECTOR i (w), ADD COLUMN w INT",
		code:  1210,
	},
	{
		name:  "added after it, nullable",
		setup: "CREATE TABLE t (id INT)",
		stmt:  "ALTER TABLE t ADD VECTOR i (w), ADD COLUMN w VECTOR(3) NULL",
		code:  1252,
	},
	{
		name:  "the old vector index dropped after it",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v))",
		stmt:  "ALTER TABLE t ADD VECTOR j (v), DROP INDEX i",
		keys:  []string{"VECTOR KEY `j` (`v`)\n"},
	},
	{
		name:  "two in one statement",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, w VECTOR(3) NOT NULL)",
		stmt:  "ALTER TABLE t ADD VECTOR k (v), ADD VECTOR l (w)",
		code:  1235,
	},
	{
		name:  "a plain key added after it",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, a INT)",
		stmt:  "ALTER TABLE t ADD VECTOR i (v), ADD KEY z (a)",
		keys:  []string{"KEY `z` (`a`)", "VECTOR KEY `i` (`v`)\n"},
	},

	// IF NOT EXISTS, twice and beside a plain ADD.
	{
		name:  "twice, on a key the table has",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)",
		stmt:  "ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD VECTOR IF NOT EXISTS i (v)",
		keys:  []string{"VECTOR KEY `i` (`v`) `M`=6\n"},
	},
	{
		name:  "twice, on a key it lacks",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL)",
		stmt:  "ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD VECTOR IF NOT EXISTS i (v)",
		keys:  []string{"VECTOR KEY `i` (`v`)\n"},
	},
	{
		name:  "a plain vector key, then IF NOT EXISTS",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL)",
		stmt:  "ALTER TABLE t ADD VECTOR i (v), ADD VECTOR IF NOT EXISTS i (v)",
		keys:  []string{"VECTOR KEY `i` (`v`)\n"},
	},
	{
		name:  "IF NOT EXISTS, then a plain vector key",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL)",
		stmt:  "ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD VECTOR i (v)",
		code:  1061,
	},
	{
		name:  "IF NOT EXISTS, then another kind of key",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, a INT)",
		stmt:  "ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD INDEX i (a)",
		code:  1061,
	},
	{
		name:  "another kind of key, then IF NOT EXISTS",
		setup: "CREATE TABLE t (v VECTOR(3) NOT NULL, a INT)",
		stmt:  "ALTER TABLE t ADD INDEX i (a), ADD VECTOR IF NOT EXISTS i (v)",
		code:  1061,
	},
	{
		name:  "IF NOT EXISTS on a column added after it",
		setup: "CREATE TABLE t (id INT)",
		stmt:  "ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD COLUMN v VECTOR(3) NOT NULL",
		keys:  []string{"VECTOR KEY `i` (`v`)\n"},
	},
}

// TestVectorState holds the catalog to vectorStateCases.
func TestVectorState(t *testing.T) {
	for _, tc := range vectorStateCases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			execAlterClauses(t, c, "CREATE DATABASE d; USE d;")
			if tc.setup != "" {
				execAlterClauses(t, c, tc.setup+";")
			}
			before := c.ShowCreateTable("d", "t")
			if tc.code != 0 {
				var cerr *Error
				if err := execAlterClausesExpectError(t, c, tc.stmt+";"); !errors.As(err, &cerr) || cerr.Code != tc.code {
					t.Fatalf("error %v, want %d", err, tc.code)
				}
				if after := c.ShowCreateTable("d", "t"); after != before {
					t.Errorf("the refused statement changed the table:\n%s\nwas:\n%s", after, before)
				}
			} else {
				execAlterClauses(t, c, tc.stmt+";")
			}
			if len(tc.keys) > 0 {
				requireInOrder(t, c.ShowCreateTable("d", "t"), tc.keys)
			}
		})
	}
}

// TestVectorStateOracle holds vectorStateCases to the server.
func TestVectorStateOracle(t *testing.T) {
	ctx := context.Background()
	conn := mariadbDatabase(t)
	for _, tc := range vectorStateCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS t"); err != nil {
				t.Fatal(err)
			}
			if tc.setup != "" {
				if _, err := conn.ExecContext(ctx, tc.setup); err != nil {
					t.Fatalf("setup: %v", err)
				}
			}
			var code int
			if _, err := conn.ExecContext(ctx, tc.stmt); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatal(err)
				}
				code = int(myErr.Number)
			}
			if code != tc.code {
				t.Fatalf("the server answered %d, the record says %d", code, tc.code)
			}
			if len(tc.keys) > 0 {
				var name, got string
				if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE t").Scan(&name, &got); err != nil {
					t.Fatal(err)
				}
				requireInOrder(t, got, tc.keys)
			}
		})
	}
}

// TestFailedCreateIndexLeavesNoHiddenColumns: a statement that fails leaves the table as it was, its
// columns and its indexes. MariaDB 11.8.8 has no functional key part, so `((a+1))` is a 1064 the
// parser refuses, and the statement here fails in the catalog on a column the table lacks.
func TestFailedCreateIndexLeavesNoHiddenColumns(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (a INT)", nil); err != nil {
		t.Fatal(err)
	}
	before := len(c.GetDatabase("d").GetTable("t").Columns)
	results, err := c.Exec("CREATE INDEX i ON t (a, nosuch)", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Error == nil {
		t.Fatalf("results = %+v, want the statement's 1072", results)
	}
	tbl := c.GetDatabase("d").GetTable("t")
	if got := len(tbl.Columns); got != before {
		t.Errorf("the failed statement left %d columns, want %d", got, before)
	}
	if len(tbl.Indexes) != 0 {
		t.Errorf("the failed statement left %d indexes", len(tbl.Indexes))
	}
}
