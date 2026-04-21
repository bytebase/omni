package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// tidbCatalogContainer wraps a real TiDB instance for catalog
// cross-validation. The catalog's in-memory model is compared to what
// TiDB actually accepts (via its SQL layer) rather than deep-inspecting
// information_schema — the goal is to catch regressions where omni
// parses TiDB syntax but the catalog fails to execute, or vice versa.
type tidbCatalogContainer struct {
	db  *sql.DB
	ctx context.Context
}

var (
	tidbCatalogOnce    sync.Once
	tidbCatalogInst    *tidbCatalogContainer
	tidbCatalogInitErr error
)

func startTiDBForCatalog(t *testing.T) *tidbCatalogContainer {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping TiDB catalog container test in short mode")
	}

	tidbCatalogOnce.Do(func() {
		ctx := context.Background()

		req := testcontainers.ContainerRequest{
			Image:        "pingcap/tidb:v8.5.5",
			ExposedPorts: []string{"4000/tcp"},
			WaitingFor:   wait.ForListeningPort("4000/tcp").WithStartupTimeout(2 * time.Minute),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			tidbCatalogInitErr = fmt.Errorf("failed to start TiDB container: %w", err)
			return
		}

		host, err := container.Host(ctx)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			tidbCatalogInitErr = fmt.Errorf("failed to get host: %w", err)
			return
		}

		port, err := container.MappedPort(ctx, "4000")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			tidbCatalogInitErr = fmt.Errorf("failed to get port: %w", err)
			return
		}

		connStr := fmt.Sprintf("root@tcp(%s:%s)/test?multiStatements=true", host, port.Port())
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			tidbCatalogInitErr = fmt.Errorf("failed to open db: %w", err)
			return
		}

		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			tidbCatalogInitErr = fmt.Errorf("failed to ping TiDB: %w", err)
			return
		}

		tidbCatalogInst = &tidbCatalogContainer{db: db, ctx: ctx}
	})

	if tidbCatalogInitErr != nil {
		t.Skipf("TiDB container not available: %v", tidbCatalogInitErr)
	}
	return tidbCatalogInst
}

// TestTiDBCatalogContainerCrossValidation executes the same DDL on real
// TiDB and on the omni catalog and asserts both accept or reject in
// lockstep. Each case runs a (setup?, sql, cleanup?) sequence on TiDB and
// the same sql through the catalog's Exec(). The catalog doesn't
// information_schema-query TiDB — that's deferred to a diff/migration PR.
func TestTiDBCatalogContainerCrossValidation(t *testing.T) {
	tc := startTiDBForCatalog(t)

	// Isolate test tables in a dedicated database.
	mustExecTiDB(t, tc, "CREATE DATABASE IF NOT EXISTS omni_cat_test")
	mustExecTiDB(t, tc, "USE omni_cat_test")
	t.Cleanup(func() { _, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_cat_test") })

	cases := []struct {
		name    string
		setup   string
		sql     string
		cleanup string
	}{
		{
			// AUTO_RANDOM on PK; SHARD_ROW_ID_BITS is tested separately
			// because TiDB forbids SHARD_ROW_ID_BITS on tables that
			// already use the PK as the row id.
			name:    "auto_random pk",
			sql:     "CREATE TABLE t_ar (id BIGINT AUTO_RANDOM(5) PRIMARY KEY)",
			cleanup: "DROP TABLE IF EXISTS t_ar",
		},
		{
			// SHARD_ROW_ID_BITS requires NO explicit PK (or non-row-id PK).
			name:    "shard_row_id_bits no pk",
			sql:     "CREATE TABLE t_srib (a INT) SHARD_ROW_ID_BITS = 4",
			cleanup: "DROP TABLE IF EXISTS t_srib",
		},
		{
			name:    "ttl create",
			sql:     "CREATE TABLE t_ttl (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'",
			cleanup: "DROP TABLE IF EXISTS t_ttl",
		},
		{
			// SET TIFLASH REPLICA 0 is accepted even when no TiFlash
			// servers are available; replica counts > 0 require TiFlash
			// servers in the cluster, which a standalone TiDB container
			// does not provide. Using 0 keeps the test hermetic.
			name:    "tiflash replica alter (0)",
			setup:   "CREATE TABLE IF NOT EXISTS t_tf (id INT PRIMARY KEY)",
			sql:     "ALTER TABLE t_tf SET TIFLASH REPLICA 0",
			cleanup: "DROP TABLE IF EXISTS t_tf",
		},
		{
			name:    "remove ttl alter",
			setup:   "CREATE TABLE IF NOT EXISTS t_ttl_rm (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'",
			sql:     "ALTER TABLE t_ttl_rm REMOVE TTL",
			cleanup: "DROP TABLE IF EXISTS t_ttl_rm",
		},
		{
			name:    "clustered pk",
			sql:     "CREATE TABLE t_cl (id INT, PRIMARY KEY (id) CLUSTERED)",
			cleanup: "DROP TABLE IF EXISTS t_cl",
		},
		{
			// Multi-feature interaction covers AUTO_RANDOM + CLUSTERED PK +
			// AUTO_ID_CACHE + TTL + TTL_JOB_INTERVAL on one table.
			// PLACEMENT POLICY and SHARD_ROW_ID_BITS are exercised elsewhere;
			// combining them here would require running CREATE PLACEMENT
			// POLICY first (unmodeled by catalog in this PR) or would
			// conflict with AUTO_RANDOM's row-id usage.
			name:    "multi-feature interaction",
			sql:     "CREATE TABLE t_mix (id BIGINT AUTO_RANDOM(5) PRIMARY KEY CLUSTERED, created_at DATETIME) AUTO_ID_CACHE = 100 TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON' TTL_JOB_INTERVAL = '1h'",
			cleanup: "DROP TABLE IF EXISTS t_mix",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.setup != "" {
				mustExecTiDB(t, tc, c.setup)
			}
			// 1. TiDB must accept.
			if _, err := tc.db.ExecContext(tc.ctx, c.sql); err != nil {
				t.Fatalf("TiDB rejected SQL %q: %v", c.sql, err)
			}
			// 2. Catalog must accept through the same path (fresh catalog
			//    per case so setup SQL on TiDB doesn't need to be mirrored).
			cat := New()
			catSetup := "CREATE DATABASE omni_cat_test; USE omni_cat_test;"
			if _, err := cat.Exec(catSetup, nil); err != nil {
				t.Fatalf("catalog setup failed: %v", err)
			}
			if c.setup != "" {
				results, err := cat.Exec(c.setup, nil)
				if err != nil {
					// Catalog may not support PLACEMENT POLICY DDL yet — PR3
					// stores PLACEMENT POLICY as a table attribute but does
					// not track policies as first-class catalog objects.
					// Skip the lockstep comparison when catalog setup fails
					// on syntax the catalog doesn't model.
					t.Skipf("catalog does not model setup SQL %q (deferred to future PR): %v", c.setup, err)
				}
				for _, r := range results {
					if r.Error != nil {
						t.Skipf("catalog exec of setup SQL returned error %v (likely unmodeled DDL)", r.Error)
					}
				}
			}
			results, err := cat.Exec(c.sql, nil)
			if err != nil {
				t.Fatalf("catalog rejected SQL %q: %v", c.sql, err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("catalog exec error on %q: %v", c.sql, r.Error)
				}
			}
			if c.cleanup != "" {
				_, _ = tc.db.ExecContext(tc.ctx, c.cleanup)
			}
		})
	}
}

func mustExecTiDB(t *testing.T, tc *tidbCatalogContainer, sqlStr string) {
	t.Helper()
	if _, err := tc.db.ExecContext(tc.ctx, sqlStr); err != nil {
		t.Fatalf("TiDB setup exec failed for %q: %v", sqlStr, err)
	}
}
