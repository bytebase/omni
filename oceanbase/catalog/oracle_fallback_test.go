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

// The TestOracle_* family was written against a hand-provisioned pair of local
// engines (8.0 on :13306, 5.7 on :13307). connectOracle still prefers those,
// because they are the fastest thing to iterate against, but when nothing
// listens there it falls back to one testcontainer per version, shared by
// every test in the process. That is what lets the family run unattended in
// CI and on a fresh checkout instead of silently skipping.

type oracleFallbackEngine struct {
	once      sync.Once
	container *tcmysql.MySQLContainer
	db        *sql.DB
	err       error
}

var oracleFallbacks = map[Version]*oracleFallbackEngine{
	MySQL57: {},
	MySQL80: {},
}

func oracleFallbackImage(version Version) (string, []testcontainers.ContainerCustomizer) {
	opts := []testcontainers.ContainerCustomizer{
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
	}
	if version == MySQL57 {
		// mysql:5.7 only ships linux/amd64; arm64 hosts run it under emulation.
		return "mysql:5.7", append(opts, testcontainers.WithImagePlatform("linux/amd64"))
	}
	return "mysql:8.0", opts
}

// fallbackOracle returns the shared fallback engine for version, starting it
// on first use. Like the other container gates in this package it fails the
// test in CI when the engine cannot start and skips locally.
func fallbackOracle(t *testing.T, version Version, name string) *sql.DB {
	t.Helper()
	eng := oracleFallbacks[version]
	eng.once.Do(func() {
		ctx := context.Background()
		image, opts := oracleFallbackImage(version)
		container, err := tcmysql.Run(ctx, image, opts...)
		if err != nil {
			eng.err = fmt.Errorf("start %s container: %w", image, err)
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
			_ = testcontainers.TerminateContainer(eng.container)
		}
	}
}
