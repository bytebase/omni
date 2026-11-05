package parser

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// oracleServer is the one postgres:17-alpine the oracles of this package ask.
// Each oracle has a database of its own, because each brings its own fixtures;
// everything they probe runs in a transaction they roll back, or only parses.
var oracleServer struct {
	once  sync.Once
	dsn   *url.URL
	admin *sql.DB
	err   error
}

func startOracleServer(ctx context.Context) (err error) {
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	container, err := tcpg.Run(ctx, "postgres:17-alpine",
		tcpg.WithDatabase("omni_admin"),
		tcpg.WithUsername("postgres"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	if err != nil {
		return fmt.Errorf("container start: %w", err)
	}
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("conn string: %w", err)
	}
	if oracleServer.dsn, err = url.Parse(connStr); err != nil {
		return fmt.Errorf("conn string %q: %w", connStr, err)
	}
	if oracleServer.admin, err = sql.Open("pgx", connStr); err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	if err := oracleServer.admin.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}

// OracleDatabase creates the named database on the shared server and returns a
// pool of connections to it. It drops a database of that name first, so a test
// that calls it again, as under -count=2, starts empty. It is exported for the
// oracles in package parser_test.
func OracleDatabase(ctx context.Context, name string) (*sql.DB, error) {
	oracleServer.once.Do(func() { oracleServer.err = startOracleServer(ctx) })
	if oracleServer.err != nil {
		return nil, oracleServer.err
	}
	for _, stmt := range []string{"DROP DATABASE IF EXISTS " + name + " WITH (FORCE)", "CREATE DATABASE " + name} {
		if _, err := oracleServer.admin.ExecContext(ctx, stmt); err != nil {
			return nil, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	dsn := *oracleServer.dsn
	dsn.Path = "/" + name
	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}
