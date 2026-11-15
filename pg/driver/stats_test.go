package driver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	pganalysis "github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

func TestCountAffectedRows(t *testing.T) {
	ctx := context.Background()
	for _, role := range []string{"pg12", testcontainer.Default, testcontainer.PG18} {
		t.Run(role, func(t *testing.T) {
			d := openQueryDriver(t, role)
			mustExec(t, d.db, `
				CREATE TABLE s (id int PRIMARY KEY, flag bool NOT NULL);
				INSERT INTO s SELECT g, g <= 10 FROM generate_series(1, 100) g;
				CREATE TABLE big (id int PRIMARY KEY, s_id int NOT NULL, v int);
				INSERT INTO big SELECT g, (g % 100) + 1, 0 FROM generate_series(1, 1000) g;
				CREATE TABLE big_archive (LIKE big);
				ANALYZE s;
				ANALYZE big;
				ANALYZE big_archive;
			`)

			stmts := []string{
				"DELETE FROM big USING s WHERE big.s_id = s.id AND s.flag;",
				"UPDATE big SET v = 1 FROM s WHERE big.s_id = s.id AND s.flag;",
				"DELETE FROM big WHERE s_id IN (SELECT id FROM s WHERE flag);",
				"INSERT INTO big_archive SELECT big.* FROM big JOIN s ON big.s_id = s.id WHERE s.flag;",
			}
			if role != "pg12" {
				stmts = append(stmts, "MERGE INTO big USING s ON big.s_id = s.id AND s.flag WHEN MATCHED THEN UPDATE SET v = 1;")
			}
			for _, stmt := range stmts {
				rows, err := d.CountAffectedRows(ctx, stmt)
				if err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
				if rows < 50 || rows > 200 {
					t.Errorf("%s: estimated rows = %d, want ~100", stmt, rows)
				}
			}

			if _, err := d.CountAffectedRows(ctx, "SELECT * FROM big;"); err == nil {
				t.Error("CountAffectedRows accepted SELECT")
			}

			mustExec(t, d.db, `
				CREATE SCHEMA app;
				CREATE TABLE app.only_in_app (id int);
				INSERT INTO app.only_in_app SELECT generate_series(1, 50);
				ANALYZE app.only_in_app;
			`)
			d.db.SetMaxOpenConns(1)
			rows, err := d.CountAffectedRows(ctx, pganalysis.WithSearchPath("DELETE FROM only_in_app;", []string{"app"}))
			if err != nil || rows != 50 {
				t.Fatalf("WithSearchPath rows = %d, err = %v", rows, err)
			}
			if _, err := d.CountAffectedRows(ctx, "DELETE FROM only_in_app;"); err == nil || !strings.Contains(err.Error(), "does not exist") {
				t.Fatalf("search_path leaked after rollback: err = %v", err)
			}
		})
	}
}

func TestExplainAndRoles(t *testing.T) {
	ctx := context.Background()
	for _, role := range []string{"pg12", testcontainer.Default} {
		t.Run(role, func(t *testing.T) {
			d := openQueryDriver(t, role)
			mustExec(t, d.db, `
				CREATE SCHEMA app;
				CREATE TABLE app.items (id int PRIMARY KEY);
				INSERT INTO app.items SELECT generate_series(1, 25);
				ANALYZE app.items;
				CREATE ROLE explain_owner NOLOGIN;
				GRANT USAGE ON SCHEMA app TO explain_owner;
				GRANT SELECT, INSERT, DELETE ON app.items TO explain_owner;
			`)
			mustExec(t, d.db, fmt.Sprintf("ALTER DATABASE %s OWNER TO explain_owner;", quoteIdent(d.databaseName)))

			owner, err := d.CurrentDatabaseOwner(ctx)
			if err != nil || owner != "explain_owner" {
				t.Fatalf("CurrentDatabaseOwner = %q, %v", owner, err)
			}
			user, err := d.CurrentUser(ctx)
			if err != nil || user != "postgres" {
				t.Fatalf("CurrentUser = %q, %v", user, err)
			}
			sp, err := d.SearchPath(ctx)
			if err != nil || sp == "" {
				t.Fatalf("SearchPath = %q, %v", sp, err)
			}

			td := &Driver{db: d.db, databaseName: d.databaseName, tenantMode: true}
			if err := td.Explain(ctx, []string{`SET LOCAL search_path TO "app"`}, "DELETE FROM items WHERE id = 1"); err != nil {
				t.Fatalf("Explain failed: %v", err)
			}
			if err := td.Explain(ctx, nil, "DELETE FROM missing_table"); err == nil {
				t.Fatal("Explain on missing_table succeeded")
			}

			plan, err := td.ExplainJSON(ctx, []string{`SET LOCAL search_path TO "app"`}, "INSERT INTO items SELECT id + 100 FROM items")
			if err != nil {
				t.Fatalf("ExplainJSON failed: %v", err)
			}
			inserted, err := pganalysis.EstimatedInsertedRowsFromExplainJSON(plan)
			if err != nil || inserted != 25 {
				t.Fatalf("EstimatedInsertedRowsFromExplainJSON = %d, %v", inserted, err)
			}
		})
	}
}
