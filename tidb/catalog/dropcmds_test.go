package catalog

import "testing"

func TestDropTable(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE TABLE t1 (id INT)", nil)
	_, err := c.Exec("DROP TABLE t1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("test").GetTable("t1") != nil {
		t.Fatal("table should be dropped")
	}
}

func TestDropTableIfExists(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	results, _ := c.Exec("DROP TABLE IF EXISTS noexist", nil)
	if results[0].Error != nil {
		t.Errorf("IF EXISTS should not error: %v", results[0].Error)
	}
}

func TestDropTableNotExists(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	results, _ := c.Exec("DROP TABLE noexist", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected error")
	}
}

func TestDropMultipleTables(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE TABLE t1 (id INT)", nil)
	c.Exec("CREATE TABLE t2 (id INT)", nil)
	c.Exec("DROP TABLE t1, t2", nil)
	db := c.GetDatabase("test")
	if db.GetTable("t1") != nil || db.GetTable("t2") != nil {
		t.Fatal("both tables should be dropped")
	}
}

func TestTruncateTable(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE TABLE t1 (id INT AUTO_INCREMENT PRIMARY KEY)", nil)
	results, _ := c.Exec("TRUNCATE TABLE t1", nil)
	if results[0].Error != nil {
		t.Fatalf("truncate failed: %v", results[0].Error)
	}
}

func TestTruncateTableNotExists(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	results, _ := c.Exec("TRUNCATE TABLE noexist", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected error")
	}
}
