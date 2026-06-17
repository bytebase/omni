package catalog

import "testing"

func TestCreateDatabase(t *testing.T) {
	c := New()
	_, err := c.Exec("CREATE DATABASE mydb", nil)
	if err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("mydb")
	if db == nil {
		t.Fatal("database not found")
	}
	if db.Name != "mydb" {
		t.Errorf("expected name 'mydb', got %q", db.Name)
	}
}

func TestCreateDatabaseIfNotExists(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE mydb", nil)
	results, _ := c.Exec("CREATE DATABASE IF NOT EXISTS mydb", nil)
	if results[0].Error != nil {
		t.Errorf("IF NOT EXISTS should not error: %v", results[0].Error)
	}
}

func TestCreateDatabaseDuplicate(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE mydb", nil)
	results, _ := c.Exec("CREATE DATABASE mydb", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate database error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrDupDatabase {
		t.Errorf("expected error code %d, got %d", ErrDupDatabase, catErr.Code)
	}
}

func TestDropDatabase(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE mydb", nil)
	_, err := c.Exec("DROP DATABASE mydb", nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("mydb") != nil {
		t.Fatal("database should be dropped")
	}
}

func TestDropDatabaseNotExists(t *testing.T) {
	c := New()
	results, _ := c.Exec("DROP DATABASE noexist", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected error for nonexistent database")
	}
}

func TestDropDatabaseIfExists(t *testing.T) {
	c := New()
	results, _ := c.Exec("DROP DATABASE IF EXISTS noexist", nil)
	if results[0].Error != nil {
		t.Errorf("IF EXISTS should not error: %v", results[0].Error)
	}
}

func TestCreateDatabaseCharset(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE mydb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", nil)
	db := c.GetDatabase("mydb")
	if db == nil {
		t.Fatal("database not found")
	}
	if db.Charset != "utf8mb4" {
		t.Errorf("expected charset utf8mb4, got %q", db.Charset)
	}
	if db.Collation != "utf8mb4_unicode_ci" {
		t.Errorf("expected collation utf8mb4_unicode_ci, got %q", db.Collation)
	}
}

func TestDropDatabaseResetsCurrentDB(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE mydb", nil)
	c.SetCurrentDatabase("mydb")
	c.Exec("DROP DATABASE mydb", nil)
	if c.CurrentDatabase() != "" {
		t.Error("current database should be unset after drop")
	}
}
