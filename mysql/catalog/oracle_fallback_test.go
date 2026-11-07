package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// All live catalog oracles share the pinned MySQL 8.4 server. Older
// normalization models are covered by unit fixtures.

type oracleFallbackEngine struct {
	once      sync.Once
	container *testcontainer.DatabaseContainer
	db        *sql.DB
	err       error
}

var oracleFallbacks = map[Version]*oracleFallbackEngine{
	MySQL80: {},
}

// fallbackOracle returns the shared fallback engine for version, starting it
// on first use. Like the other container gates in this package it fails the
// test in CI when the engine cannot start and skips locally.
func fallbackOracle(t *testing.T, version Version, name string) *sql.DB {
	t.Helper()
	eng := oracleFallbacks[version]
	eng.once.Do(func() {
		ctx := context.Background()
		image := testcontainer.Image
		container, err := testcontainer.Start(ctx, "test")
		if err != nil {
			eng.err = fmt.Errorf("start %s container: %w", image, err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "multiStatements=true")
		if err != nil {
			_ = testcontainer.TerminateContainer(container)
			eng.err = err
			return
		}
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			_ = testcontainer.TerminateContainer(container)
			eng.err = err
			return
		}
		if err := db.PingContext(ctx); err != nil {
			_ = db.Close()
			_ = testcontainer.TerminateContainer(container)
			eng.err = fmt.Errorf("ping %s container: %w", image, err)
			return
		}
		eng.container, eng.db = container, db
	})
	if eng.err != nil {
		t.Fatalf("oracle %s required in CI but unavailable: %v", name, eng.err)
	}
	return eng.db
}

func terminateOracleFallbacks() {
	for _, eng := range oracleFallbacks {
		if eng.db != nil {
			_ = eng.db.Close()
		}
		if eng.container != nil {
			_ = testcontainer.TerminateContainer(eng.container)
		}
	}
}
