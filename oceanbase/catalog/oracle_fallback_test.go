package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// The TestOracle_* family was written against a hand-provisioned local engine
// (8.0 on :13306). connectOracle still prefers it, because it is the fastest
// thing to iterate against, but when nothing listens there it falls back to
// one testcontainer, shared by every test in the process. That is what lets
// the family run unattended in CI and on a fresh checkout instead of silently
// skipping.

const oracleFallbackImage = "mysql:8.0"

var oracleFallback struct {
	once      sync.Once
	container *tcmysql.MySQLContainer
	db        *sql.DB
	err       error
}

// fallbackOracle returns the shared fallback engine, starting it on first use.
// It fails the test when the engine cannot start.
func fallbackOracle(t testing.TB) *sql.DB {
	t.Helper()
	eng := &oracleFallback
	eng.once.Do(func() {
		ctx := context.Background()
		container, err := tcmysql.Run(ctx, oracleFallbackImage,
			tcmysql.WithDatabase("test"),
			tcmysql.WithUsername("root"),
			tcmysql.WithPassword("test"),
		)
		if err != nil {
			eng.err = fmt.Errorf("start %s container: %w", oracleFallbackImage, err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "multiStatements=true")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			eng.err = err
			return
		}
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			eng.err = err
			return
		}
		if err := db.PingContext(ctx); err != nil {
			_ = db.Close()
			_ = testcontainers.TerminateContainer(container)
			eng.err = fmt.Errorf("ping %s container: %w", oracleFallbackImage, err)
			return
		}
		eng.container, eng.db = container, db
	})
	if eng.err != nil {
		t.Fatalf("oracle %s required in CI but unavailable: %v", oracleFallbackImage, eng.err)
	}
	return eng.db
}

func terminateOracleFallback() {
	if oracleFallback.db != nil {
		_ = oracleFallback.db.Close()
	}
	if oracleFallback.container != nil {
		_ = testcontainers.TerminateContainer(oracleFallback.container)
	}
}
