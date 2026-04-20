package catalog

import "testing"

func TestRenameTable(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE TABLE t1 (id INT)", nil)
	_, err := c.Exec("RENAME TABLE t1 TO t2", nil)
	if err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("test")
	if db.GetTable("t1") != nil {
		t.Fatal("old table should not exist")
	}
	if db.GetTable("t2") == nil {
		t.Fatal("new table should exist")
	}
}

func TestRenameTableCrossDatabase(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE db1", nil)
	c.Exec("CREATE DATABASE db2", nil)
	c.SetCurrentDatabase("db1")
	c.Exec("CREATE TABLE t1 (id INT)", nil)
	_, err := c.Exec("RENAME TABLE db1.t1 TO db2.t2", nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("db1").GetTable("t1") != nil {
		t.Fatal("old table should not exist in db1")
	}
	if c.GetDatabase("db2").GetTable("t2") == nil {
		t.Fatal("new table should exist in db2")
	}
}

func TestRenameTableNotFound(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	results, _ := c.Exec("RENAME TABLE noexist TO t2", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected error for missing table")
	}
}
