package parser

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// srContainer wraps a StarRocks allin1 container for parse-conformance testing.
// StarRocks speaks the MySQL wire protocol on port 9030 (not 3306), so we
// connect with go-sql-driver/mysql via a raw GenericContainer (there is no
// testcontainers StarRocks module). arm64 note (#65095): the FE can come up
// without a healthy BE on arm64, so we force linux/amd64 (qemu emulation).
type srContainer struct {
	db  *sql.DB
	ctx context.Context
}

var (
	srOnce    sync.Once
	srInst    *srContainer
	srInitErr error
)

// startStarRocks starts (once) the StarRocks 3.4 container and returns a ready
// srContainer. Skipped in -short mode — omni CI runs `go test -short ./...`, so
// container tests must be Short-gated (mirrors tidb/parser/tidb_container_test.go).
func startStarRocks(t *testing.T) *srContainer {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping StarRocks container test in short mode")
	}

	srOnce.Do(func() {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:         "starrocks/allin1-ubuntu:3.4.10",
			ImagePlatform: "linux/amd64", // arm64 #65095 workaround
			ExposedPorts:  []string{"9030/tcp"},
			WaitingFor:    wait.ForListeningPort("9030/tcp").WithStartupTimeout(4 * time.Minute),
		}
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			srInitErr = fmt.Errorf("start StarRocks container: %w", err)
			return
		}
		host, err := container.Host(ctx)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			srInitErr = err
			return
		}
		port, err := container.MappedPort(ctx, "9030")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			srInitErr = err
			return
		}
		dsn := fmt.Sprintf("root@tcp(%s:%s)/?multiStatements=true&timeout=10s", host, port.Port())
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			srInitErr = err
			return
		}

		// Active readiness: the FE accepts connections before it can serve
		// queries. Poll until SELECT 1 succeeds. FE readiness suffices for a
		// parse oracle (syntax errors are raised at parse time, before BE).
		deadline := time.Now().Add(4 * time.Minute)
		for {
			ok := func() bool {
				cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				var one int
				return db.QueryRowContext(cctx, "SELECT 1").Scan(&one) == nil && one == 1
			}()
			if ok {
				break
			}
			if time.Now().After(deadline) {
				_ = db.Close()
				_ = testcontainers.TerminateContainer(container)
				srInitErr = fmt.Errorf("readiness: SELECT 1 never succeeded")
				return
			}
			time.Sleep(3 * time.Second)
		}

		// Select a default database so analysis-stage probes don't all fail with
		// "No database selected" (which would obscure whether they parsed).
		if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS srtest"); err != nil {
			_ = db.Close()
			_ = testcontainers.TerminateContainer(container)
			srInitErr = fmt.Errorf("create srtest db: %w", err)
			return
		}
		_ = db.Close()
		db, err = sql.Open("mysql", fmt.Sprintf("root@tcp(%s:%s)/srtest?multiStatements=true&timeout=10s", host, port.Port()))
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			srInitErr = err
			return
		}
		// Container intentionally leaked for singleton reuse; process exit reaps it.
		srInst = &srContainer{db: db, ctx: ctx}
	})

	if srInitErr != nil {
		t.Fatalf("StarRocks container: %v", srInitErr)
	}
	return srInst
}

// syntaxErrMarker is the prefix StarRocks' parser emits for a genuine PARSE
// failure. CRITICAL: StarRocks codes BOTH parse AND semantic ("analyzing")
// errors as MySQL 1064, so 1064 alone does not mean rejected — the reliable
// discriminator is this message. (See reference_starrocks_oracle_1064_gotcha.)
const syntaxErrMarker = "Getting syntax error"

// accepts reports whether StarRocks accepted sql at the PARSE layer: executed
// OK, or failed with any non-syntax error (semantic/runtime) — both imply it
// parsed. Only a "Getting syntax error" message means rejected.
func (c *srContainer) accepts(sql string) bool {
	ctx, cancel := context.WithTimeout(c.ctx, 8*time.Second)
	defer cancel()
	_, err := c.db.ExecContext(ctx, sql)
	if err == nil {
		return true
	}
	if myErr, ok := err.(*gomysql.MySQLError); ok {
		return !strings.Contains(myErr.Message, syntaxErrMarker)
	}
	return true // non-driver error (conn/timeout) — treat as accepted
}

// assertParity asserts the omni parser and StarRocks agree on sql (both accept
// or both reject) — the accept+reject conformance discipline.
func (c *srContainer) assertParity(t *testing.T, sql string, wantAccept bool) {
	t.Helper()
	_, errs := Parse(sql)
	omni := len(errs) == 0
	sr := c.accepts(sql)
	if omni != sr {
		t.Errorf("MISMATCH %q: omni=%v starrocks=%v (omniErrs=%v)", sql, omni, sr, errs)
	}
	if omni != wantAccept {
		t.Errorf("omni %q: got accept=%v, want %v (errs=%v)", sql, omni, wantAccept, errs)
	}
}

// TestStarRocksConnectivity is the harness smoke test: the container is reachable.
func TestStarRocksConnectivity(t *testing.T) {
	c := startStarRocks(t)
	var one int
	if err := c.db.QueryRowContext(c.ctx, "SELECT 1").Scan(&one); err != nil || one != 1 {
		t.Fatalf("SELECT 1 = %d, err=%v", one, err)
	}
}
