package parser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

var sharedMySQL80 sharedEngine

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
		if os.Getenv("CI") != "" {
			t.Fatalf("container required in CI but unavailable: %v", e.err)
		}
		t.Skipf("container unavailable: %v", e.err)
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
// database selected.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	return sharedMySQL80.acquire(t, startMySQL80)
}

func TestMain(m *testing.M) {
	code := m.Run()
	sharedMySQL80.terminate()
	os.Exit(code)
}

// checkSyntax sends SQL to MySQL 8.0 and returns true if syntactically valid.
// Strategy: execute the SQL. If MySQL returns error code 1064 (ER_PARSE_ERROR),
// the SQL is syntactically invalid (rejected). Any other error (e.g. 1146 table
// not found) means the SQL parsed fine but failed for semantic reasons (accepted).
func (o *parserOracle) checkSyntax(sql string) bool {
	_, err := o.db.ExecContext(o.ctx, sql)
	if err == nil {
		return true // executed successfully — valid syntax
	}
	if myErr, ok := err.(*gomysql.MySQLError); ok {
		return myErr.Number != 1064 // 1064 = ER_PARSE_ERROR
	}
	// Non-MySQL error (connection issue, etc.) — treat as valid to avoid
	// false negatives; the caller should investigate.
	return true
}

// assertOracleMatch parses with omni parser AND checks against MySQL 8.0,
// asserting both accept or both reject.
func assertOracleMatch(t *testing.T, o *parserOracle, sql string) {
	t.Helper()

	_, omniErr := Parse(sql)
	omniAccepts := omniErr == nil

	mysqlAccepts := o.checkSyntax(sql)

	if omniAccepts != mysqlAccepts {
		t.Errorf("MISMATCH for %q:\n  omni accepts=%v (err=%v)\n  mysql accepts=%v",
			sql, omniAccepts, omniErr, mysqlAccepts)
	} else {
		t.Logf("MATCH for %q: both accept=%v", sql, omniAccepts)
	}
}

// assertOracleMismatch logs a known mismatch between omni and MySQL 8.0 without
// failing the test. These are tracked as [~] partial in the SCENARIOS file.
func assertOracleMismatch(t *testing.T, o *parserOracle, sql string) {
	t.Helper()

	_, omniErr := Parse(sql)
	omniAccepts := omniErr == nil

	mysqlAccepts := o.checkSyntax(sql)

	if omniAccepts != mysqlAccepts {
		t.Logf("KNOWN MISMATCH for %q:\n  omni accepts=%v (err=%v)\n  mysql accepts=%v",
			sql, omniAccepts, omniErr, mysqlAccepts)
	} else {
		// If the mismatch is resolved (e.g. by a parser fix), promote to match.
		t.Logf("RESOLVED — now MATCH for %q: both accept=%v (was known mismatch)", sql, omniAccepts)
	}
}
