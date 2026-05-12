//go:build oracle

package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type loaderCompatOracle struct {
	db  *sql.DB
	ctx context.Context
}

var (
	loaderCompatOracleOnce     sync.Once
	loaderCompatOracleInst     *loaderCompatOracle
	loaderCompatOracleSetupErr error
)

func startLoaderCompatOracle(t *testing.T) *loaderCompatOracle {
	t.Helper()
	loaderCompatOracleOnce.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				loaderCompatOracleSetupErr = fmt.Errorf("docker provider panic: %v", r)
			}
		}()

		ctx := context.Background()
		container, err := tcpg.Run(ctx, "postgres:17-alpine",
			tcpg.WithDatabase("omni_loader_compat"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			loaderCompatOracleSetupErr = fmt.Errorf("container start: %w", err)
			return
		}

		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			loaderCompatOracleSetupErr = fmt.Errorf("conn string: %w", err)
			return
		}
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			loaderCompatOracleSetupErr = fmt.Errorf("db open: %w", err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			loaderCompatOracleSetupErr = fmt.Errorf("ping: %w", err)
			return
		}

		loaderCompatOracleInst = &loaderCompatOracle{db: db, ctx: ctx}
	})

	if loaderCompatOracleSetupErr != nil {
		t.Skipf("loader compatibility oracle unavailable: %v", loaderCompatOracleSetupErr)
	}
	return loaderCompatOracleInst
}

func (o *loaderCompatOracle) execIsolated(t *testing.T, ddl string) {
	t.Helper()
	if err := o.execIsolatedErr(t, ddl); err != nil {
		t.Fatalf("PG17 rejected SQL that the compatibility corpus expects to be accepted: %v\nSQL:\n%s", err, ddl)
	}
}

func (o *loaderCompatOracle) execIsolatedErr(t *testing.T, ddl string) error {
	t.Helper()
	schema := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_"))
	schema = strings.ReplaceAll(schema, "-", "_")
	if _, err := o.db.ExecContext(o.ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schema)); err != nil {
		return fmt.Errorf("PG cleanup failed: %w", err)
	}
	if _, err := o.db.ExecContext(o.ctx, fmt.Sprintf(`CREATE SCHEMA %q`, schema)); err != nil {
		return fmt.Errorf("PG schema setup failed: %w", err)
	}
	t.Cleanup(func() {
		_, _ = o.db.ExecContext(o.ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schema))
	})

	sql := fmt.Sprintf("SET search_path TO %q, public;\n%s", schema, ddl)
	_, err := o.db.ExecContext(o.ctx, sql)
	return err
}

func TestLoaderCompatPG17OracleAcceptsCorpus(t *testing.T) {
	oracle := startLoaderCompatOracle(t)
	for _, tc := range loaderCompatAcceptCases() {
		t.Run(tc.name, func(t *testing.T) {
			oracle.execIsolated(t, tc.sql)
			if _, err := LoadSQL(tc.sql); err != nil {
				t.Fatalf("LoadSQL rejected PG17-accepted SQL: %v\nSQL:\n%s", err, tc.sql)
			}
		})
	}
}

func TestLoaderCompatPG17OracleRejectsCorpus(t *testing.T) {
	oracle := startLoaderCompatOracle(t)
	for _, tc := range loaderCompatRejectCases() {
		t.Run(tc.name, func(t *testing.T) {
			if oracle.execIsolatedErr(t, tc.sql) == nil {
				t.Fatalf("PG17 accepted SQL that the reject corpus expects to be rejected\nSQL:\n%s", tc.sql)
			}
			if _, err := LoadSQL(tc.sql); err == nil {
				t.Fatalf("LoadSQL accepted SQL that the reject corpus expects to be rejected\nSQL:\n%s", tc.sql)
			}
		})
	}
}
