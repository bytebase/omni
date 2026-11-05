package catalog

import "testing"

// constraintOfIndexScripts: a CHECK or a foreign key may share an index's name without being tied
// to it, and the server keeps it when the index goes.
var constraintOfIndexScripts = map[string][]string{
	"DROP INDEX keeps a CHECK": {
		"CREATE TABLE t (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))", "DROP INDEX i ON t",
	},
	"ALTER DROP INDEX keeps a CHECK": {
		"CREATE TABLE t (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))", "ALTER TABLE t DROP INDEX i",
	},
	"DROP INDEX on a unique key keeps a CHECK written first": {
		"CREATE TABLE t (a INT, CONSTRAINT i CHECK (a > 0), UNIQUE KEY i (a))", "DROP INDEX i ON t",
	},
	"DROP INDEX keeps a foreign key": {
		"CREATE TABLE p (id INT PRIMARY KEY)",
		"CREATE TABLE t (a INT, b INT, KEY j (a), KEY i (b), CONSTRAINT i FOREIGN KEY (a) REFERENCES p (id))",
		"DROP INDEX i ON t",
	},
}

func TestConstraintOfIndexSurvivesOtherNamesOracle(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	for name, script := range constraintOfIndexScripts {
		t.Run(name, func(t *testing.T) {
			if err := ctr.execSQL("DROP TABLE IF EXISTS t; DROP TABLE IF EXISTS p"); err != nil {
				t.Fatal(err)
			}
			c := New()
			mustExec(t, c, "CREATE DATABASE d")
			c.SetCurrentDatabase("d")
			for _, stmt := range script {
				if err := ctr.execSQL(stmt); err != nil {
					t.Fatalf("the server refused %q: %v", stmt, err)
				}
				mustExec(t, c, stmt)
			}
			server, err := ctr.showCreateTable("t")
			if err != nil {
				t.Fatal(err)
			}
			if got := c.ShowCreateTable("d", "t"); got != server {
				t.Errorf("\nomni:   %s\nserver: %s", got, server)
			}
		})
	}
}
