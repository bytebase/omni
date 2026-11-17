package driver

import (
	"context"
	"testing"

	"github.com/bytebase/omni/driver"
)

func TestExecute(t *testing.T) {
	d := openTestDriver(t)
	mustExec(t, d.db, `
CREATE SCHEMA app;
CREATE TABLE app.account (id INT8 PRIMARY KEY, v INT8 NOT NULL);`)

	ctx := context.Background()

	t.Run("transaction", func(t *testing.T) {
		affected, err := d.Execute(ctx, `
INSERT INTO app.account VALUES (1, 10);
INSERT INTO app.account VALUES (2, 20);
UPDATE app.account SET v = v + 1 WHERE id = 1;`, driver.ExecuteOptions{Transaction: driver.Transactional})
		if err != nil {
			t.Fatal(err)
		}
		if affected != 3 {
			t.Errorf("affected = %d, want 3", affected)
		}
		var sum int64
		if err := d.db.QueryRowContext(ctx, "SELECT sum(v) FROM app.account;").Scan(&sum); err != nil {
			t.Fatal(err)
		}
		if sum != 31 {
			t.Errorf("sum = %d, want 31", sum)
		}
	})

	t.Run("auto commit", func(t *testing.T) {
		affected, err := d.Execute(ctx, `
INSERT INTO app.account VALUES (3, 30);
INSERT INTO app.account VALUES (4, 40);`, driver.ExecuteOptions{Transaction: driver.AutoCommit})
		if err != nil {
			t.Fatal(err)
		}
		if affected != 2 {
			t.Errorf("affected = %d, want 2", affected)
		}
		var count int64
		if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM app.account WHERE id IN (3, 4);").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})

	t.Run("error rolls back", func(t *testing.T) {
		affected, err := d.Execute(ctx, `
INSERT INTO app.account VALUES (5, 50);
INSERT INTO app.account VALUES (1, 99);`, driver.ExecuteOptions{Transaction: driver.Transactional})
		if err == nil {
			t.Fatal("Execute succeeded, want the duplicate key error")
		}
		if affected != 0 {
			t.Errorf("affected = %d, want 0", affected)
		}
		var count int64
		if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM app.account WHERE id = 5;").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Errorf("id 5 count = %d, want 0 after the rollback", count)
		}
	})

	t.Run("create database", func(t *testing.T) {
		if _, err := d.Execute(ctx, "CREATE DATABASE omni_execute_created;", driver.ExecuteOptions{CreateDatabase: true}); err != nil {
			t.Fatal(err)
		}
		var count int64
		if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM pg_database WHERE datname = 'omni_execute_created';").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("database count = %d, want 1", count)
		}
		// Creating it again is a no-op.
		if _, err := d.Execute(ctx, "CREATE DATABASE omni_execute_created;", driver.ExecuteOptions{CreateDatabase: true}); err != nil {
			t.Fatal(err)
		}
		// A newline is not a statement boundary.
		if _, err := d.Execute(ctx, "CREATE DATABASE\nomni_execute_created_multiline;", driver.ExecuteOptions{CreateDatabase: true}); err != nil {
			t.Fatal(err)
		}
		if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM pg_database WHERE datname = 'omni_execute_created_multiline';").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("multiline database count = %d, want 1", count)
		}
		if _, err := d.Execute(ctx, `create database "team data";`, driver.ExecuteOptions{CreateDatabase: true}); err != nil {
			t.Fatal(err)
		}
		if err := d.db.QueryRowContext(ctx, `SELECT count(*) FROM pg_database WHERE datname = 'team data';`).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("quoted database count = %d, want 1", count)
		}
	})

	t.Run("auto commit keeps script order", func(t *testing.T) {
		if _, err := d.Execute(ctx, "CREATE DATABASE omni_execute_order; DROP DATABASE omni_execute_order;", driver.ExecuteOptions{Transaction: driver.AutoCommit}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("a literal does not change classification", func(t *testing.T) {
		mustExec(t, d.db, "CREATE TABLE app.note (id INT8 PRIMARY KEY, v STRING);")
		_, err := d.Execute(ctx, `
INSERT INTO app.note VALUES (1, 'ok');
INSERT INTO app.note VALUES (1, 'DROP DATABASE');`, driver.ExecuteOptions{Transaction: driver.Transactional})
		if err == nil {
			t.Fatal("Execute succeeded, want the duplicate key error")
		}
		var count int64
		if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM app.note;").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Errorf("count = %d, want 0 after the rollback", count)
		}
	})
}

func TestGetDatabaseInCreateDatabaseStatement(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      string
		wantErr   bool
	}{
		{`CREATE DATABASE "hello" ENCODING "UTF8";`, "hello", false},
		{`CREATE DATABASE "hello";`, "hello", false},
		{`CREATE DATABASE hello;`, "hello", false},
		{`CREATE DATABASE hello ENCODING "UTF8";`, "hello", false},
		{`create database "team data";`, "team data", false},
		{`CREATE DATABASE;`, "", true},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got, err := getDatabaseInCreateDatabaseStatement(tc.statement)
			if tc.wantErr {
				if err == nil {
					t.Fatal("want an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetRoutingIDFromCockroachCloudURL(t *testing.T) {
	for _, tc := range []struct {
		host string
		want string
	}{
		{"routing-id.cockroachlabs.cloud", "routing-id"},
		{"bytebase-cluster-7749.6xw.aws-ap-southeast-1.cockroachlabs.cloud", "bytebase-cluster-7749"},
		{"subdomain.routing-id.cockroachlabs.cloud", "subdomain"},
		{"cockroachlabs.cloud", ""},
		{"example.com", ""},
	} {
		if got := getRoutingIDFromCockroachCloudURL(tc.host); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.host, got, tc.want)
		}
	}
}

// TestGetCurrentDatabaseOwner reads the requested method directly. Bytebase's
// plancheck runner type-asserts it on the concrete driver.
func TestGetCurrentDatabaseOwner(t *testing.T) {
	d := openTestDriver(t)
	owner, err := d.GetCurrentDatabaseOwner(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if owner != "root" {
		t.Errorf("owner = %q, want root", owner)
	}
}
