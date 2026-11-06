package testcontainer

import (
	"context"
	"errors"
	"os"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	code := m.Run()
	MySQL80.Terminate()
	os.Exit(code)
}

// escapes are statements that would reach past the test's database if the
// account could run them. mysql:8.0 parses every one.
var escapes = []string{
	"SHUTDOWN",
	"RESTART",
	"SET GLOBAL sql_mode = 'ANSI_QUOTES'",
	"SET PERSIST max_connections = 1",
	"ALTER USER root IDENTIFIED BY 'gone'",
	"SET PASSWORD FOR root = 'gone'",
	"CREATE USER intruder",
	"DROP USER root",
	"GRANT ALL ON *.* TO omni_test",
	"DROP DATABASE mysql",
	"CREATE DATABASE outside_the_prefix",
	"FLUSH TABLES WITH READ LOCK",
	"FLUSH PRIVILEGES",
	"LOCK INSTANCE FOR BACKUP",
	"CREATE RESOURCE GROUP rg TYPE = USER",
	"RESET MASTER",
	"PURGE BINARY LOGS BEFORE NOW()",
	"INSTALL PLUGIN p SONAME 'p.so'",
	"KILL 1",
}

// TestRestrictedAccountCannotReachTheServer is the promise that lets a corpus
// run on the server its package shares: whatever it runs stays inside the
// test's database. Each statement parses, so a verdict that looks for 1064
// still reads it as parsed.
func TestRestrictedAccountCannotReachTheServer(t *testing.T) {
	ctx := context.Background()
	own, err := MySQL80.Database(ctx, "own")
	if err != nil {
		t.Fatal(err)
	}
	defer own.Close(ctx)
	other, err := MySQL80.Database(ctx, "other")
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close(ctx)
	if _, err := other.DB.ExecContext(ctx, "CREATE TABLE kept (a INT)"); err != nil {
		t.Fatal(err)
	}

	for _, stmt := range escapes {
		_, err := own.DB.ExecContext(ctx, stmt)
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) {
			t.Errorf("%s: want a server error, got %v", stmt, err)
			continue
		}
		if myErr.Number == 1064 {
			t.Errorf("%s: a parse error says nothing about the privilege: %v", stmt, err)
		}
	}

	var user string
	if err := own.DB.QueryRowContext(ctx, "SELECT CURRENT_USER()").Scan(&user); err != nil {
		t.Fatal(err)
	}
	if user != account+"@%" {
		t.Errorf("connected as %q", user)
	}
	var n int
	if err := other.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM kept").Scan(&n); err != nil {
		t.Errorf("the other database's table is gone: %v", err)
	}
}

// TestResetEmptiesTheDatabase holds Reset to what a server of its own would give.
func TestResetEmptiesTheDatabase(t *testing.T) {
	ctx := context.Background()
	db, err := MySQL80.Database(ctx, "reset")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)
	if _, err := db.DB.ExecContext(ctx, "CREATE TABLE t (a INT)"); err != nil {
		t.Fatal(err)
	}
	if err := db.Reset(ctx); err != nil {
		t.Fatal(err)
	}
	var selected string
	if err := db.DB.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&selected); err != nil {
		t.Fatal(err)
	}
	if selected != db.Name {
		t.Errorf("selected %q, want %q", selected, db.Name)
	}
	var n int
	if err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("%d tables survive Reset", n)
	}
}
