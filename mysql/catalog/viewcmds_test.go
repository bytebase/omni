package catalog

import "testing"

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
