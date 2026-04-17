package parser

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

type tidbContainer struct {
	db  *sql.DB
	ctx context.Context
}

var (
	tidbOnce    sync.Once
	tidbInst    *tidbContainer
	tidbInitErr error
)

func startTiDB(t *testing.T) *tidbContainer {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping TiDB container test in short mode")
	}

	tidbOnce.Do(func() {
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
			tidbInitErr = fmt.Errorf("failed to start TiDB container: %w", err)
			return
		}

		host, err := container.Host(ctx)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			tidbInitErr = fmt.Errorf("failed to get host: %w", err)
			return
		}

		port, err := container.MappedPort(ctx, "4000")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			tidbInitErr = fmt.Errorf("failed to get port: %w", err)
			return
		}

		connStr := fmt.Sprintf("root@tcp(%s:%s)/test?multiStatements=true", host, port.Port())
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			tidbInitErr = fmt.Errorf("failed to open db: %w", err)
			return
		}

		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			tidbInitErr = fmt.Errorf("failed to ping TiDB: %w", err)
			return
		}

		tidbInst = &tidbContainer{db: db, ctx: ctx}
	})

	if tidbInitErr != nil {
		t.Skipf("TiDB container not available: %v", tidbInitErr)
	}
	return tidbInst
}

// canExecute checks whether TiDB accepts the SQL syntax.
func (tc *tidbContainer) canExecute(sqlStr string) error {
	_, err := tc.db.ExecContext(tc.ctx, sqlStr)
	return err
}

func TestTiDBContainerOracle(t *testing.T) {
	tc := startTiDB(t)

	// Create a test database for isolation.
	tc.db.ExecContext(tc.ctx, "CREATE DATABASE IF NOT EXISTS omni_test")
	tc.db.ExecContext(tc.ctx, "USE omni_test")
	defer tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_test")

	tests := []struct {
		name      string
		sql       string
		wantParse bool   // should our parser accept?
		checkTiDB bool   // also check TiDB execution? (false = parse-only)
		setup     string // optional setup SQL (run before test SQL)
		cleanup   string // optional cleanup SQL (run after test SQL)
	}{
		// AUTO_RANDOM
		{
			name:      "auto_random basic",
			sql:       "CREATE TABLE t_ar1 (id BIGINT AUTO_RANDOM PRIMARY KEY)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_ar1",
		},
		{
			name:      "auto_random shard bits",
			sql:       "CREATE TABLE t_ar2 (id BIGINT AUTO_RANDOM(5) PRIMARY KEY)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_ar2",
		},
		{
			name:      "auto_random shard+range",
			sql:       "CREATE TABLE t_ar3 (id BIGINT AUTO_RANDOM(5, 64) PRIMARY KEY)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_ar3",
		},

		// CLUSTERED / NONCLUSTERED
		{
			name:      "clustered pk",
			sql:       "CREATE TABLE t_cl1 (id INT, PRIMARY KEY (id) CLUSTERED)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_cl1",
		},
		{
			name:      "nonclustered pk",
			sql:       "CREATE TABLE t_cl2 (id INT, PRIMARY KEY (id) NONCLUSTERED)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_cl2",
		},

		// Table options
		{
			name:      "shard_row_id_bits",
			sql:       "CREATE TABLE t_opt1 (id INT) SHARD_ROW_ID_BITS = 4",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_opt1",
		},
		{
			name:      "auto_id_cache",
			sql:       "CREATE TABLE t_opt2 (id INT) AUTO_ID_CACHE = 100",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_opt2",
		},
		{
			name:      "pre_split_regions",
			sql:       "CREATE TABLE t_opt3 (id INT) SHARD_ROW_ID_BITS = 4 PRE_SPLIT_REGIONS = 3",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_opt3",
		},

		// TTL
		{
			name:      "ttl create",
			sql:       "CREATE TABLE t_ttl1 (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_ttl1",
		},
		{
			name:      "remove ttl",
			sql:       "ALTER TABLE t_ttl_rm REMOVE TTL",
			wantParse: true, checkTiDB: true,
			setup:   "CREATE TABLE IF NOT EXISTS t_ttl_rm (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'",
			cleanup: "DROP TABLE IF EXISTS t_ttl_rm",
		},

		// PLACEMENT POLICY (table and database)
		{
			name:      "placement policy table",
			sql:       "CREATE TABLE t_pp1 (id INT) PLACEMENT POLICY = 'default'",
			wantParse: true, checkTiDB: false,
			cleanup: "DROP TABLE IF EXISTS t_pp1",
		},
		{
			name:      "placement policy database",
			sql:       "CREATE DATABASE IF NOT EXISTS d_pp1 PLACEMENT POLICY = 'default'",
			wantParse: true, checkTiDB: false,
			cleanup: "DROP DATABASE IF EXISTS d_pp1",
		},

		// SET TIFLASH REPLICA
		{
			name:      "set tiflash replica",
			sql:       "ALTER TABLE t_tf SET TIFLASH REPLICA 0",
			wantParse: true, checkTiDB: true,
			setup:   "CREATE TABLE IF NOT EXISTS t_tf (id INT PRIMARY KEY)",
			cleanup: "DROP TABLE IF EXISTS t_tf",
		},

		// TiDB comments
		{
			name:      "tidb comment basic",
			sql:       "CREATE TABLE t_cmt1 (id BIGINT /*T! AUTO_RANDOM */ PRIMARY KEY)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_cmt1",
		},
		{
			name:      "tidb comment feature",
			sql:       "CREATE TABLE t_cmt2 (id BIGINT /*T![auto_rand] AUTO_RANDOM */ PRIMARY KEY)",
			wantParse: true, checkTiDB: true,
			cleanup: "DROP TABLE IF EXISTS t_cmt2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setup != "" {
				if _, err := tc.db.ExecContext(tc.ctx, tt.setup); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			// Test our parser.
			_, parseErr := Parse(tt.sql)
			parserAccepts := parseErr == nil

			if parserAccepts != tt.wantParse {
				t.Errorf("parser: got accepts=%v, want %v (err=%v)", parserAccepts, tt.wantParse, parseErr)
			}

			// Test real TiDB.
			if tt.checkTiDB {
				tidbErr := tc.canExecute(tt.sql)
				tidbAccepts := tidbErr == nil

				if parserAccepts != tidbAccepts {
					t.Errorf("MISMATCH: parser accepts=%v, TiDB accepts=%v (tidb err=%v, parse err=%v)",
						parserAccepts, tidbAccepts, tidbErr, parseErr)
				}
			}

			// Cleanup
			if tt.cleanup != "" {
				tc.db.ExecContext(tc.ctx, tt.cleanup)
			}
		})
	}
}
