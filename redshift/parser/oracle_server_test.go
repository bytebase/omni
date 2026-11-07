package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/redshift/internal/testcontainer"
)

// oracleServer is the pinned PostgreSQL 18 compatibility baseline the package's oracles ask.
// Each oracle has a database of its own, because each brings its own fixtures;
// everything they probe runs in a transaction they roll back, or only parses.
var oracleServer struct {
	once  sync.Once
	dsn   *url.URL
	admin *sql.DB
	err   error
}

// startOracleServer retries creating the container because a busy CI runner can
// transiently blow a creation deadline (image pull, dockerd contention) or lose
// a fresh container to the Ryuk reaper race. Genuine unavailability still fails:
// every attempt's error is preserved.
func startOracleServer(ctx context.Context) error {
	if err := dockerPreflight(2 * time.Second); err != nil {
		return fmt.Errorf("docker preflight: %w", err)
	}
	const attempts = 3
	var errs []error
	for attempt := 1; attempt <= attempts; attempt++ {
		if attempt > 1 {
			time.Sleep(2 * time.Second)
		}
		err := startOracleServerAttempt(ctx)
		if err == nil {
			return nil
		}
		errs = append(errs, fmt.Errorf("attempt %d/%d: %w", attempt, attempts, err))
	}
	return errors.Join(errs...)
}

// startOracleServerAttempt makes one full attempt, terminating what it created
// when it fails so that retries do not accumulate containers on the runner.
func startOracleServerAttempt(ctx context.Context) (err error) {
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	// Generous per-attempt timeout: on a clean CI runner the image pull alone
	// can exceed a short one.
	startCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	container, err := testcontainer.Start(startCtx, "omni_admin")
	defer func() {
		if err != nil {
			// Run can return a partially-created container alongside its
			// error, and TerminateContainer takes a nil one.
			_ = testcontainer.TerminateContainer(container)
		}
	}()
	if err != nil {
		return fmt.Errorf("container start: %w", err)
	}
	connStr, err := container.ConnectionString(startCtx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("conn string: %w", err)
	}
	if oracleServer.dsn, err = url.Parse(connStr); err != nil {
		return fmt.Errorf("conn string %q: %w", connStr, err)
	}
	if oracleServer.admin, err = sql.Open("pgx", connStr); err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	if err := oracleServer.admin.PingContext(startCtx); err != nil {
		oracleServer.admin.Close()
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
