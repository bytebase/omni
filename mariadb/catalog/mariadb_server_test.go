package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmariadb "github.com/testcontainers/testcontainers-go/modules/mariadb"
)

// sharedMariaDB is the one mariadb:11.8.8 server the tests of this package ask.
// startContainer's server is MySQL 8.0, so it is no substitute for a MariaDB
// verdict.
var sharedMariaDB struct {
	once      sync.Once
	container *tcmariadb.MariaDBContainer
	config    *gomysql.Config
	admin     *sql.DB
	err       error
	databases int
}

func startSharedMariaDB(ctx context.Context) error {
	s := &sharedMariaDB
	container, err := tcmariadb.Run(ctx, "mariadb:11.8.8",
		tcmariadb.WithUsername("root"),
		tcmariadb.WithPassword("test"),
	)
	if container != nil {
		s.container = container
	}
	if err != nil {
		return fmt.Errorf("start MariaDB container: %w", err)
	}
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("connection string: %w", err)
	}
	if s.config, err = gomysql.ParseDSN(dsn); err != nil {
		return fmt.Errorf("parse %q: %w", dsn, err)
	}
	if s.admin, err = sql.Open("mysql", dsn); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := s.admin.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}

// mariadbDatabase creates a database for this test on the shared MariaDB and
// returns a pool whose connections start in it. No connection is kept idle, so
// every case that takes one starts a session with no temporary tables and
// default settings. The database is dropped when the test ends.
//
// A statement of a test may also name a database of its own, such as orp or d1,
// and drops it first; tests that ask for a database do not run in parallel.
func mariadbDatabase(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	s := &sharedMariaDB
	s.once.Do(func() { s.err = startSharedMariaDB(ctx) })
	if s.err != nil {
		t.Fatalf("shared MariaDB: %v", s.err)
	}
	s.databases++
	name := fmt.Sprintf("mdb_%d", s.databases)
	if _, err := s.admin.ExecContext(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatalf("create database %s: %v", name, err)
	}
	t.Cleanup(func() {
		if _, err := s.admin.ExecContext(ctx, "DROP DATABASE "+name); err != nil {
			t.Errorf("drop database %s: %v", name, err)
		}
	})
	config := s.config.Clone()
	config.DBName = name
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxIdleConns(0)
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return db
}

func terminateSharedMariaDB() {
	if sharedMariaDB.admin != nil {
		_ = sharedMariaDB.admin.Close()
	}
	if sharedMariaDB.container != nil {
		_ = testcontainers.TerminateContainer(sharedMariaDB.container)
	}
}
