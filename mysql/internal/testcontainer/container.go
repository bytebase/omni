// Package testcontainer starts the MySQL server mysql's oracle tests run
// against: one testcontainer per test binary, started by the first test that
// asks and removed by Terminate, which the package's TestMain calls.
//
// A test gets a database of its own on that server, reached through an account
// that holds privileges only on the databases this package creates, so a
// statement a corpus runs (SET GLOBAL, SHUTDOWN, ALTER USER root) cannot change
// the server under the tests that follow. Those statements answer 1227 or 1044,
// which a verdict that looks for 1064 still reads as parsed.
package testcontainer

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/bytebase/omni/enginetest"
)

// Pins is mysql/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

const (
	// The test account and the prefix of its databases. The underscores are
	// escaped in the GRANT, where they would otherwise match any character.
	account        = "omni_test"
	accountSecret  = "omni_test"
	databasePrefix = "omni_t_"
)

// Server is one MySQL server a test binary starts on first use.
type Server struct {
	// Image is the pinned image.
	Image string
	// Args are the server's arguments.
	Args []string

	once      sync.Once
	config    *gomysql.Config
	container testcontainers.Container
	err       error
}

// MySQL80 is the default row of pins.tsv, mysql:8.0.46. The binary log stays
// off so the restricted account can create routines (1419 otherwise).
var MySQL80 = &Server{Image: "mysql:8.0.46", Args: []string{"--skip-log-bin"}}

// CheckPin fails unless the server this process started is the one pin names:
// the image, its repository digest, and the release it reports. A tag can be
// pushed again; a digest cannot.
func (s *Server) CheckPin(ctx context.Context, pin enginetest.Pin) error {
	if pin.Image != s.Image {
		return fmt.Errorf("%s pins %s, but this package starts %s", Pins, pin.Image, s.Image)
	}
	admin, err := s.admin(ctx)
	if err != nil {
		return err
	}
	defer admin.Close()
	info, err := s.container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the MySQL container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the MySQL image: %w", err)
	}
	found := false
	for _, d := range img.RepoDigests {
		found = found || strings.HasSuffix(d, "@"+pin.Digest)
	}
	if !found {
		return fmt.Errorf("the MySQL container runs an image with digests %v, but %s pins %s", img.RepoDigests, Pins, pin.Digest)
	}
	var version string
	if err := admin.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return fmt.Errorf("read the MySQL version: %w", err)
	}
	if "mysql-"+version != pin.Tag {
		return fmt.Errorf("the server reports %q, but %s pins %s", version, Pins, pin.Tag)
	}
	return nil
}

// admin returns a connection that can create databases, starting the server on
// the first call of this process.
func (s *Server) admin(ctx context.Context) (*sql.DB, error) {
	s.once.Do(func() { s.err = s.start(ctx) })
	if s.err != nil {
		return nil, s.err
	}
	db, err := sql.Open("mysql", s.config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func (s *Server) start(ctx context.Context) error {
	c, err := tcmysql.Run(ctx, s.Image,
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
		testcontainers.WithCmdArgs(s.Args...),
	)
	if c != nil {
		s.container = c
	}
	if err != nil {
		return fmt.Errorf("start %s: %w", s.Image, err)
	}
	dsn, err := c.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("connection string: %w", err)
	}
	config, err := gomysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("parse DSN: %w", err)
	}
	config.DBName = ""
	config.MultiStatements = false
	s.config = config

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}
	defer db.Close()
	escaped := strings.ReplaceAll(databasePrefix, "_", `\_`) + "%"
	for _, stmt := range []string{
		fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", account, accountSecret),
		fmt.Sprintf("GRANT ALL ON `%s`.* TO '%s'@'%%'", escaped, account),
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

// Terminate removes the server, if this process started one.
func (s *Server) Terminate() {
	if s.container != nil {
		_ = testcontainers.TerminateContainer(s.container)
	}
}

// Database is a test's own database on a Server.
type Database struct {
	// DB is a pool of connections that start in the database, on the restricted
	// account, pinned to one connection so that USE and session state hold across
	// statements.
	DB   *sql.DB
	Name string
}

var labelChars = regexp.MustCompile(`[^a-z0-9_]`)

// Database creates a database named for label and returns it.
func (s *Server) Database(ctx context.Context, label string) (*Database, error) {
	admin, err := s.admin(ctx)
	if err != nil {
		return nil, err
	}
	defer admin.Close()
	name := databasePrefix + labelChars.ReplaceAllString(strings.ToLower(label), "_")
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+name); err != nil {
		return nil, fmt.Errorf("create database %s: %w", name, err)
	}
	config := s.config.Clone()
	config.User, config.Passwd, config.DBName = account, accountSecret, name
	config.ParseTime, config.MultiStatements = true, true
	// Arguments are inlined, so a statement such as PREPARE, which the server
	// does not take as a prepared statement, can take one.
	config.InterpolateParams = true
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to %s: %w", name, err)
	}
	return &Database{DB: db, Name: name}, nil
}

// Reset empties the database and selects it, as a server of its own would be for
// the next test.
func (d *Database) Reset(ctx context.Context) error {
	for _, stmt := range []string{"DROP DATABASE IF EXISTS " + d.Name, "CREATE DATABASE " + d.Name, "USE " + d.Name} {
		if _, err := d.DB.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

// Close drops the database and closes the pool.
func (d *Database) Close(ctx context.Context) {
	_, _ = d.DB.ExecContext(ctx, "DROP DATABASE IF EXISTS "+d.Name)
	_ = d.DB.Close()
}
