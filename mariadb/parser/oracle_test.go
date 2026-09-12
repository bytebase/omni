package parser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

var (
	sharedMySQL80    sharedEngine // routine alignment oracle (startParserOracle)
	sharedMariaDB118 sharedEngine // MariaDB 11.8 divergence inventory (startMariaDB)
)

// sharedEngine is an engine started once per test process and handed to every
// test that asks for it, with the `test` database dropped and recreated in
// between. Booting MySQL or MariaDB costs ~10s; the reset costs milliseconds
// and gives each test the same empty schema a fresh container would.
type sharedEngine struct {
	once      sync.Once
	container testcontainers.Container
	db        *sql.DB
	err       error
}

// acquire returns the shared engine as a parserOracle, starting it on first
// use via start. It fails the test in CI when the engine cannot start and
// skips locally, so a laptop without Docker still runs the rest of the package.
func (e *sharedEngine) acquire(t *testing.T, start func(ctx context.Context) (testcontainers.Container, *sql.DB, error)) *parserOracle {
	t.Helper()
	ctx := context.Background()
	e.once.Do(func() {
		container, db, err := start(ctx)
		if err != nil {
			e.err = err
			return
		}
		// USE is session state: pin the pool to one connection so every
		// statement runs in the freshly reset database.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		e.container, e.db = container, db
	})
	if e.err != nil {
		t.Fatalf("container required in CI but unavailable: %v", e.err)
	}
	for _, stmt := range []string{"DROP DATABASE IF EXISTS test", "CREATE DATABASE test", "USE test"} {
		if _, err := e.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("reset test database (%s): %v", stmt, err)
		}
	}
	return &parserOracle{db: e.db, ctx: ctx}
}

func (e *sharedEngine) terminate() {
	if e.db != nil {
		_ = e.db.Close()
	}
	if e.container != nil {
		_ = testcontainers.TerminateContainer(e.container)
	}
}

// startMySQL80 boots a mysql:8.0 testcontainer for sharedEngine.acquire.
func startMySQL80(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	container, err := tcmysql.Run(ctx, "mysql:8.0",
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start MySQL container: %w", err)
	}
	connStr, err := container.ConnectionString(ctx, "parseTime=true", "multiStatements=true")
	if err != nil {
		_ = testcontainers.TerminateContainer(container)
		return nil, nil, fmt.Errorf("connection string: %w", err)
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		_ = testcontainers.TerminateContainer(container)
		return nil, nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		_ = testcontainers.TerminateContainer(container)
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}
	return container, db, nil
}

// startParserOracle returns the shared MySQL 8.0 oracle with an empty `test`
// database selected. Used by the routine alignment oracle
// (routine_alignment_test.go). The MariaDB 11.8 divergence inventory uses
// startMariaDB (oracle_corpus_test.go) instead — repointing the routine oracle
// to MariaDB is follow-up (P1+) work, out of PR1 scope.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	return sharedMySQL80.acquire(t, startMySQL80)
}

func TestMain(m *testing.M) {
	code := m.Run()
	sharedMySQL80.terminate()
	sharedMariaDB118.terminate()
	os.Exit(code)
}
