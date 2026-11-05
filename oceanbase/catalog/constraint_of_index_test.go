package catalog

import (
	"context"
	"strings"
	"testing"
	"time"
)

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

func TestOracle_ConstraintOfIndexSurvivesOtherNames(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for name, script := range constraintOfIndexScripts {
		t.Run(name, func(t *testing.T) {
			for _, stmt := range []string{"DROP DATABASE IF EXISTS zz_coi", "CREATE DATABASE zz_coi", "USE zz_coi"} {
				if _, err := o.db.ExecContext(ctx, stmt); err != nil {
					t.Fatal(err)
				}
			}
			c := New()
			mustExec(t, c, "CREATE DATABASE d")
			c.SetCurrentDatabase("d")
			for _, stmt := range script {
				if _, err := o.db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("the server refused %q: %v", stmt, err)
				}
				mustExec(t, c, stmt)
			}
			var table, server string
			if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE t").Scan(&table, &server); err != nil {
				t.Fatal(err)
			}
			got := c.GetDatabase("d").GetTable("t")
			for _, want := range []string{"i` CHECK", "i` FOREIGN KEY"} {
				if strings.Contains(server, want) != strings.Contains(c.ShowCreateTable("d", "t"), want) {
					t.Errorf("%q differs:\nomni:   %s\nserver: %s", want, c.ShowCreateTable("d", "t"), server)
				}
			}
			if got.Indexes == nil {
				t.Errorf("no index left")
			}
		})
	}
}
