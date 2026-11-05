package catalog

import (
	"context"
	"errors"
	"strings"
	"testing"
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
	conn := startDropPeriodMariaDB(t)
	if _, err := conn.ExecContext(ctx, vectorShowCreateDDL); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = conn.ExecContext(ctx, "DROP TABLE IF EXISTS m") })
	var name, got string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE m").Scan(&name, &got); err != nil {
		t.Fatal(err)
	}
	requireInOrder(t, got, vectorShowCreateWant)
}
