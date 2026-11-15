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
