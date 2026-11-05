package review

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// postgresServer is the one server of an image the oracle tests of this package
// ask. A version decides some cases, so the image names the server; everything
// else about a test's state lives in its database.
type postgresServer struct {
	once  sync.Once
	dsn   *url.URL
	admin *sql.DB
	err   error
}

var postgresServers = struct {
	sync.Mutex
	byImage map[string]*postgresServer
}{byImage: map[string]*postgresServer{}}

func (s *postgresServer) start(ctx context.Context, image string) (err error) {
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	container, err := tcpg.Run(ctx, image,
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
	if s.dsn, err = url.Parse(connStr); err != nil {
		return fmt.Errorf("conn string %q: %w", connStr, err)
	}
	if s.admin, err = sql.Open("pgx", connStr); err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	return s.admin.PingContext(ctx)
}

// startPostgres returns a fresh database on the PostgreSQL the oracle tests of
// this package ask.
func startPostgres(t *testing.T) *sql.DB {
	t.Helper()
	return startPostgresImage(t, "postgres:17-alpine")
}

// startPostgresImage returns a fresh database on the server of the image, for
// the cases a version decides. The database is named omni_review, which cases
// name in their text; the tests that ask for it do not run in parallel.
func startPostgresImage(t *testing.T, image string) *sql.DB {
	t.Helper()
	return startPostgresDatabase(t, image, "omni_review")
}

// startPostgresDatabase returns the named database on the server of the image,
// created for this test. Roles belong to the server, so they outlive the
// database and are visible to every database of the server.
func startPostgresDatabase(t *testing.T, image, name string) *sql.DB {
	t.Helper()
	ctx := context.Background()
	postgresServers.Lock()
	s := postgresServers.byImage[image]
	if s == nil {
		s = &postgresServer{}
		postgresServers.byImage[image] = s
	}
	postgresServers.Unlock()
	s.once.Do(func() { s.err = s.start(ctx, image) })
	if s.err != nil {
		t.Fatalf("postgres oracle unavailable: %v", s.err)
	}
	drop := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", name)
	for _, stmt := range []string{drop, "CREATE DATABASE " + name} {
		if _, err := s.admin.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		if _, err := s.admin.ExecContext(ctx, drop); err != nil {
			t.Errorf("%s: %v", drop, err)
		}
	})
	dsn := *s.dsn
	dsn.Path = "/" + name
	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return db
}
