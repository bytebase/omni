package catalog

import "testing"

// TestExecSetStatement checks a wrapped statement reaches the catalog, and each option lasts for
// that statement alone, as mariadb:11.8.8 has it.
func TestExecSetStatement(t *testing.T) {
	exec := func(t *testing.T, c *Catalog, sql string) error {
		t.Helper()
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range results {
			if r.Error != nil {
				return r.Error
			}
		}
		return nil
	}
	setup := "CREATE DATABASE d; USE d; "

	t.Run("a wrapped DDL reaches the catalog", func(t *testing.T) {
		c := New()
		if err := exec(t, c, setup+"SET STATEMENT sql_mode = '' FOR CREATE TABLE x (a INT);"); err != nil {
			t.Fatal(err)
		}
		if c.GetDatabase("d").GetTable("x") == nil {
			t.Error("table x is missing")
		}
		if err := exec(t, c, "SET STATEMENT sql_mode = '' FOR DROP TABLE x;"); err != nil {
			t.Fatal(err)
		}
		if c.GetDatabase("d").GetTable("x") != nil {
			t.Error("table x is still there")
		}
	})

	const fk = "CREATE TABLE y (a INT, FOREIGN KEY (a) REFERENCES nosuch (id))"
	t.Run("foreign_key_checks lasts for the statement", func(t *testing.T) {
		c := New()
		if err := exec(t, c, setup+"SET STATEMENT foreign_key_checks = 0 FOR "+fk+";"); err != nil {
			t.Fatalf("the option did not apply: %v", err)
		}
		if !c.foreignKeyChecks || !c.session.ForeignKeyChecks {
			t.Error("foreign_key_checks stayed off")
		}
		if err := exec(t, c, "CREATE TABLE z (a INT, FOREIGN KEY (a) REFERENCES nosuch (id));"); err == nil {
			t.Error("the next statement ran with the checks off")
		}
	})

	t.Run("the last option wins", func(t *testing.T) {
		c := New()
		if err := exec(t, c, setup+"SET STATEMENT foreign_key_checks = 0, foreign_key_checks = 1 FOR "+fk+";"); err == nil {
			t.Error("the checks were off")
		}
	})

	t.Run("a variable the options name goes back", func(t *testing.T) {
		c := New()
		if err := exec(t, c, setup+"SET STATEMENT foreign_key_checks = 0 FOR SET foreign_key_checks = 0;"); err != nil {
			t.Fatal(err)
		}
		if !c.foreignKeyChecks {
			t.Error("foreign_key_checks stayed off")
		}
		if err := exec(t, c, "SET sql_mode = 'ANSI'; SET STATEMENT sql_mode = '' FOR SET sql_mode = 'TRADITIONAL';"); err != nil {
			t.Fatal(err)
		}
		if c.session.SQLMode != "ANSI" {
			t.Errorf("sql_mode = %q, want ANSI", c.session.SQLMode)
		}
	})

	t.Run("a variable the options do not name stays set", func(t *testing.T) {
		c := New()
		if err := exec(t, c, setup+"SET STATEMENT max_statement_time = 1 FOR SET foreign_key_checks = 0;"); err != nil {
			t.Fatal(err)
		}
		if c.foreignKeyChecks || c.session.ForeignKeyChecks {
			t.Error("foreign_key_checks went back")
		}
	})

	t.Run("wrappers nest", func(t *testing.T) {
		c := New()
		if err := exec(t, c, setup+"SET STATEMENT foreign_key_checks = 0 FOR SET STATEMENT sql_mode = '' FOR "+fk+";"); err != nil {
			t.Fatal(err)
		}
		if !c.foreignKeyChecks {
			t.Error("foreign_key_checks stayed off")
		}
	})

	t.Run("a query or a compound statement changes nothing", func(t *testing.T) {
		c := New()
		results, err := c.Exec(setup+"SET STATEMENT max_statement_time = 1 FOR SELECT 1; SET STATEMENT max_statement_time = 1 FOR IF 1 THEN SELECT 1; END IF;", nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("statement %d: %v", r.Index, r.Error)
			}
		}
		if got := results[len(results)-2]; !got.Skipped {
			t.Error("the wrapped SELECT was not skipped as a query")
		}
	})
}
