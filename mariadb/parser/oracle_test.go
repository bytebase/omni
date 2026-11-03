package parser

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// sharedMariaDB118 is the MariaDB 11.8 oracle (startMariaDB) the divergence
// inventory and the routine alignment oracle share.
var sharedMariaDB118 sharedEngine

// sharedEngine is an engine started once per test process and handed to every
// test that asks for it, with the `test` database dropped and recreated in
// between. Booting MariaDB costs ~10s; the reset costs milliseconds and gives
// each test the same empty schema a fresh container would.
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

func TestMain(m *testing.M) {
	code := m.Run()
	sharedMariaDB118.terminate()
	os.Exit(code)
}
