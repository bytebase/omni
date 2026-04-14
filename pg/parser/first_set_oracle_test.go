package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type firstSetOracle struct {
	db  *sql.DB
	ctx context.Context
}

// oracleSetupError remembers a setup failure so every subsequent test
// observes it and can fail or skip uniformly. sync.Once runs exactly
// once, so errors must be captured here rather than returned from Do.
var oracleSetupError error

var (
	firstSetOracleOnce sync.Once
	firstSetOracleInst *firstSetOracle
)

func startFirstSetOracle(t *testing.T) *firstSetOracle {
	t.Helper()
	firstSetOracleOnce.Do(func() {
		ctx := context.Background()
		container, err := tcpg.Run(ctx, "postgres:17-alpine",
			tcpg.WithDatabase("omni_fs"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			oracleSetupError = fmt.Errorf("container start: %w", err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleSetupError = fmt.Errorf("conn string: %w", err)
			return
		}
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleSetupError = fmt.Errorf("db open: %w", err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			oracleSetupError = fmt.Errorf("ping: %w", err)
			return
		}
		firstSetOracleInst = &firstSetOracle{db: db, ctx: ctx}
	})

	if oracleSetupError != nil {
		// CI must fail loudly — otherwise a misconfigured CI silently
		// skips every FIRST-set test and disables the entire guardrail.
		// Local dev without docker is allowed to skip.
		if isCI() {
			t.Fatalf("first-set oracle unavailable in CI: %v", oracleSetupError)
		}
		t.Skipf("first-set oracle unavailable (local dev): %v", oracleSetupError)
	}
	return firstSetOracleInst
}

// isCI reports whether we're running under continuous integration.
// Matches omni's existing conventions and major CI providers.
func isCI() bool {
	// Respect the standard CI env var set by GitHub Actions, CircleCI,
	// GitLab CI, etc. If omni has a project-specific variable, add it here.
	return os.Getenv("CI") == "true" || os.Getenv("CI") == "1"
}

// probeResult classifies a PG error into accept/reject for FIRST-set purposes.
type probeResult int

const (
	probeAccept probeResult = iota // syntax OK, may or may not be semantically valid
	probeReject                    // syntax_error (42601)
)

func classifyProbeErr(err error) probeResult {
	if err == nil {
		return probeAccept
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "42601" {
			return probeReject
		}
	}
	return probeAccept
}

// probe runs one SQL statement inside a savepoint so side effects roll back.
func (o *firstSetOracle) probe(t *testing.T, sqlStr string) probeResult {
	t.Helper()
	tx, err := o.db.BeginTx(o.ctx, nil)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(o.ctx, sqlStr)
	return classifyProbeErr(err)
}
