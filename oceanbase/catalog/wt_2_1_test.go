package catalog

import "testing"

// Section 2.1: Database Errors

func TestWalkThrough_2_1_CreateDatabaseDuplicate(t *testing.T) {
	c := wtSetup(t)
	// "testdb" already exists from wtSetup
	results, err := c.Exec("CREATE DATABASE testdb", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrDupDatabase)
}

func TestWalkThrough_2_1_CreateDatabaseIfNotExistsOnExisting(t *testing.T) {
	c := wtSetup(t)
	// "testdb" already exists from wtSetup
	results, err := c.Exec("CREATE DATABASE IF NOT EXISTS testdb", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

func TestWalkThrough_2_1_DropDatabaseUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP DATABASE nope", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrUnknownDatabase)
}

func TestWalkThrough_2_1_DropDatabaseIfExistsOnUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("DROP DATABASE IF EXISTS nope", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertNoError(t, results[0].Error)
}

func TestWalkThrough_2_1_UseUnknownDatabase(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("USE nope", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrUnknownDatabase)
}

func TestWalkThrough_2_1_CreateTableWithoutUse(t *testing.T) {
	c := New() // no database selected
	results, err := c.Exec("CREATE TABLE t (id INT)", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoDatabaseSelected)
}

func TestWalkThrough_2_1_AlterTableWithoutUse(t *testing.T) {
	c := New() // no database selected
	results, err := c.Exec("ALTER TABLE t ADD COLUMN x INT", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoDatabaseSelected)
}

func TestWalkThrough_2_1_AlterDatabaseUnknown(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("ALTER DATABASE nope CHARACTER SET utf8", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrUnknownDatabase)
}

func TestWalkThrough_2_1_TruncateUnknownTable(t *testing.T) {
	c := wtSetup(t)
	results, err := c.Exec("TRUNCATE TABLE nope", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertError(t, results[0].Error, ErrNoSuchTable)
}
