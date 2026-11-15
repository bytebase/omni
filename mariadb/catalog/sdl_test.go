package catalog

import (
	"strings"
	"testing"
)

func TestLoadSDLScopesNestedCTEs(t *testing.T) {
	// The b inside a's derived table is its CTE, so a does not depend on view b.
	c, err := LoadSDL(dbDDL + "CREATE VIEW a AS SELECT * FROM (WITH b AS (SELECT 1 AS x) SELECT x FROM b) q; CREATE VIEW b AS SELECT * FROM a;")
	if err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("app").Views["b"] == nil {
		t.Error("view b not loaded")
	}
	// A non-recursive CTE's body reads the table its name shadows.
	if _, err := LoadSDL(dbDDL + "CREATE VIEW v AS WITH t AS (SELECT * FROM t) SELECT * FROM t; CREATE TABLE t (a INT);"); err != nil {
		t.Errorf("a CTE body reading the table it shadows: %v", err)
	}
}

func TestLoadSDLRefusesParameterMarkers(t *testing.T) {
	for _, stmt := range []string{
		"CREATE TABLE t (a INT DEFAULT ?)",
		"CREATE TABLE t (a INT, CHECK (a > ?))",
		"CREATE TABLE t (a INT, b INT AS (a + ?))",
	} {
		if _, err := LoadSDL(dbDDL + stmt + ";"); err == nil || !strings.Contains(err.Error(), "parameter marker") {
			t.Errorf("%s: err = %v, want the parameter marker refused", stmt, err)
		}
	}
}

func TestLoadSDLEndsInTheLastSelectedDatabase(t *testing.T) {
	c, err := LoadSDL("CREATE DATABASE a; CREATE DATABASE b; USE a; CREATE TABLE t (x INT); USE b;")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec("CREATE TABLE u (x INT);", nil); err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("b").Tables["u"] == nil {
		t.Error("u not created in b, the database the SDL last selects")
	}
}

func TestLoadSDLKeepsOperandCTEsInTheirOperand(t *testing.T) {
	c, err := LoadSDL(dbDDL + "CREATE VIEW a AS SELECT * FROM b UNION ALL (WITH b AS (SELECT 1 AS x) SELECT * FROM b); CREATE VIEW b AS SELECT 2 AS x;")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.ShowCreateView("app", "a"); !strings.Contains(got, "AS select `b`.`x` AS `x` from `b` union all") {
		t.Errorf("view a defined before b:\n%s", got)
	}
}

func TestLoadSDLCopiesSeparateIndexesWithLike(t *testing.T) {
	c, err := LoadSDL(dbDDL + "CREATE TABLE dst LIKE src; CREATE TABLE src (a INT); CREATE UNIQUE INDEX k ON src (a);")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.ShowCreateTable("app", "dst"); !strings.Contains(got, "UNIQUE KEY `k` (`a`)") {
		t.Errorf("dst lost src's index:\n%s", got)
	}
}
