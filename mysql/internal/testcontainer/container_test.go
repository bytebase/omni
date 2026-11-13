package testcontainer

import (
	"context"
	"errors"
	"os"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/tape"
)

func TestMain(m *testing.M) {
	code := tape.Main(m)
	Terminate()
	os.Exit(code)
}

// escapes are statements that would reach past the test's database if the
// account could run them. mysql:8.0.46 parses every one.
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
// run on the server a recording shares: whatever it runs stays inside the
// test's database. Each statement parses, so a verdict that looks for 1064
// still reads it as parsed.
func TestRestrictedAccountCannotReachTheServer(t *testing.T) {
	ctx := context.Background()
	other := OpenDatabase(t, "other")
	if _, err := other.ExecContext(ctx, "CREATE TABLE kept (a INT)"); err != nil {
		t.Fatal(err)
	}
	own := OpenDatabase(t, "own")
	for _, stmt := range escapes {
		_, err := own.ExecContext(ctx, stmt)
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
	if err := own.QueryRowContext(ctx, "SELECT CURRENT_USER()").Scan(&user); err != nil {
		t.Fatal(err)
	}
	if user != account+"@%" {
		t.Errorf("connected as %q", user)
	}
	var n int
	if err := other.QueryRowContext(ctx, "SELECT COUNT(*) FROM kept").Scan(&n); err != nil {
		t.Errorf("the other database's table is gone: %v", err)
	}
}

// TestOpenDatabaseIsEmpty holds OpenDatabase to what a server of its own
// would give: the database selected, and nothing in it.
func TestOpenDatabaseIsEmpty(t *testing.T) {
	ctx := context.Background()
	db := OpenDatabase(t, "empty")
	var selected string
	var n int
	if err := db.QueryRowContext(ctx, "SELECT DATABASE(), (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE())").Scan(&selected, &n); err != nil {
		t.Fatal(err)
	}
	if selected != DatabaseName("empty") || n != 0 {
		t.Errorf("selected %q holding %d tables, want %q holding none", selected, n, DatabaseName("empty"))
	}
}

func TestCodecRoundTrips(t *testing.T) {
	for _, in := range []*gomysql.MySQLError{
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: syntaxErrorPreamble + "'x'\nat line 2"},
		{Number: 1105, Message: "no SQLSTATE"},
	} {
		text, ok := Codec{}.EncodeError(in)
		if !ok {
			t.Fatalf("%v is not encoded", in)
		}
		var out *gomysql.MySQLError
		if err := (Codec{}).DecodeError(text); !errors.As(err, &out) || *out != *in {
			t.Errorf("%q decodes to %#v, want %#v", text, err, in)
		}
	}
	if _, ok := (Codec{}).EncodeError(gomysql.ErrInvalidConn); ok {
		t.Error("a dropped connection was encoded as the engine's answer")
	}
}
