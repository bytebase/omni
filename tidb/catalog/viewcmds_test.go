package catalog

import (
	"strings"
	"testing"
)

func TestCreateView(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	_, err := c.Exec("CREATE VIEW v1 AS SELECT 1", nil)
	if err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("test")
	if db.Views[toLower("v1")] == nil {
		t.Fatal("view should exist")
	}
}

func TestCreateViewOrReplace(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE VIEW v1 AS SELECT 1", nil)
	results, _ := c.Exec("CREATE OR REPLACE VIEW v1 AS SELECT 2", nil)
	if results[0].Error != nil {
		t.Fatalf("OR REPLACE should not error: %v", results[0].Error)
	}
	if c.GetDatabase("test").Views[toLower("v1")] == nil {
		t.Fatal("view should still exist after replace")
	}
}

func TestCreateViewDuplicate(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE VIEW v1 AS SELECT 1", nil)
	results, _ := c.Exec("CREATE VIEW v1 AS SELECT 2", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestDropView(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE VIEW v1 AS SELECT 1", nil)
	_, err := c.Exec("DROP VIEW v1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("test").Views[toLower("v1")] != nil {
		t.Fatal("view should be dropped")
	}
}

func TestDropViewIfExists(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	results, _ := c.Exec("DROP VIEW IF EXISTS noexist", nil)
	if results[0].Error != nil {
		t.Errorf("IF EXISTS should not error: %v", results[0].Error)
	}
}

// TestSection_7_1_ViewCreationPipeline verifies that createView() stores the
// body as TiDB v8.5.5 does (restored by its parser after the planner names
// and expands the select fields) instead of the raw SelectText. The expected
// strings are the server's SHOW CREATE VIEW output for the same statements.
func TestSection_7_1_ViewCreationPipeline(t *testing.T) {
	t.Run("createView_uses_deparser", func(t *testing.T) {
		// Create a view with a schema-aware SELECT; the stored Definition
		// should be the deparsed output, not the raw input.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t (a INT, b INT)", nil)

		results, _ := c.Exec("CREATE VIEW v1 AS SELECT a, b FROM t", nil)
		if results[0].Error != nil {
			t.Fatalf("CREATE VIEW error: %v", results[0].Error)
		}

		db := c.GetDatabase("test")
		v := db.Views[toLower("v1")]
		if v == nil {
			t.Fatal("view v1 should exist")
		}

		// TiDB names every field and qualifies the table by its schema,
		// but leaves the column references as written.
		if !strings.Contains(v.Definition, "FROM `test`.`t`") {
			t.Errorf("Definition should qualify the table, got: %s", v.Definition)
		}
		if !strings.Contains(v.Definition, "AS `a`") {
			t.Errorf("Definition should contain alias AS `a`, got: %s", v.Definition)
		}

		// The raw input "SELECT a, b FROM t" should NOT be stored verbatim.
		if v.Definition == "SELECT a, b FROM t" {
			t.Error("Definition should be deparsed, not raw input")
		}

		expected := "SELECT `a` AS `a`,`b` AS `b` FROM `test`.`t`"
		if v.Definition != expected {
			t.Errorf("Definition mismatch:\n  got:  %q\n  want: %q", v.Definition, expected)
		}
	})

	t.Run("definition_contains_deparsed_sql", func(t *testing.T) {
		// Verify that View.Definition has the * expanded, as TiDB stores
		// it, and the rest restored as written.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t (a INT)", nil)

		c.Exec("CREATE VIEW v1 AS SELECT * FROM t WHERE a > 0", nil)
		db := c.GetDatabase("test")
		v := db.Views[toLower("v1")]

		// * is expanded to schema-qualified columns
		expected := "SELECT `test`.`t`.`a` AS `a` FROM `test`.`t` WHERE `a`>0"
		if v.Definition != expected {
			t.Errorf("Definition mismatch:\n  got:  %q\n  want: %q", v.Definition, expected)
		}
	})

	t.Run("preamble_format", func(t *testing.T) {
		// Verify SHOW CREATE VIEW preamble matches TiDB's format, which always
		// lists the view's columns:
		// CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`) AS ...
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t (a INT)", nil)
		c.Exec("CREATE VIEW v1 AS SELECT a FROM t", nil)

		ddl := c.ShowCreateView("test", "v1")
		expectedPrefix := "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` (`a`) AS "
		if !strings.HasPrefix(ddl, expectedPrefix) {
			t.Errorf("SHOW CREATE VIEW preamble mismatch:\n  got:  %q\n  want prefix: %q", ddl, expectedPrefix)
		}

		// Full output check
		expectedFull := expectedPrefix + "SELECT `a` AS `a` FROM `test`.`t`"
		if ddl != expectedFull {
			t.Errorf("SHOW CREATE VIEW full mismatch:\n  got:  %q\n  want: %q", ddl, expectedFull)
		}
	})
}
