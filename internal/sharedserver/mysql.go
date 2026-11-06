// Package sharedserver hands tests a database of their own on a server that many
// test processes and development sessions share, so that they do not each pay
// for booting one. Booting mysql:8.0 costs about 12 s and 360 MiB, and creating
// a database costs milliseconds.
//
// A test process finds its server in this order:
//
//  1. The DSN in the engine's environment variable, which CI sets. A server that
//     does not answer there fails the test; nothing falls back from a server
//     someone asked for by name.
//  2. The server scripts/dev-db.sh keeps on the engine's conventional local port,
//     if one answers and runs the engine's version.
//  3. A testcontainer of its own, started once per test binary.
//
// A test connects through a restricted account that holds privileges only on
// databases this package names, so a statement a corpus runs cannot reach another
// process through the server: SET GLOBAL, SHUTDOWN, ALTER USER root, and KILL of
// another session answer 1227 or 1044, which a verdict that looks for 1064 still
// reads as parsed. A test that needs a server-wide setting or an account, or whose
// verdict depends on a privilege, does not belong on a shared server.
package sharedserver

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

const (
	// A test's account and database names. The underscores are escaped in the
	// GRANT, where they would otherwise match any character.
	account        = "omni_shared"
	accountSecret  = "omni_shared"
	databasePrefix = "omni_t_"
	// staleAfter is how old a database must be before a later process drops it.
	// A process that gets killed leaves its database behind on a long-lived
	// server; 24 h is well past the longest test timeout.
	staleAfter = 24 * time.Hour
)

// MySQLServer describes one MySQL-protocol server a test can share.
type MySQLServer struct {
	// Image is the pinned image, for a testcontainer of its own.
	Image string
	// Version is what SELECT VERSION() must start with on a server this package
	// finds rather than starts.
	Version string
	// Env is the environment variable naming a running server, as a DSN whose
	// account can create databases and accounts.
	Env string
	// DevDSN is the conventional address of scripts/dev-db.sh's server.
	DevDSN string
	// Args are the server arguments of a testcontainer, CI's included, as
	// dev-db.sh passes them.
	Args []string
	// start boots a testcontainer and returns it with an admin DSN, for a server
	// the mysql module cannot start. Nil means the mysql module.
	start func(ctx context.Context, image string) (testcontainers.Container, string, error)
	// open connects to a DSN. Nil means sql.Open with the mysql driver.
	open func(dsn string) (*sql.DB, error)

	once      sync.Once
	config    *gomysql.Config
	container testcontainers.Container
	err       error
}

// MySQL80 is mysql:8.0. The binary log stays off so the restricted account can
// create routines (1419 otherwise), and the connection limit suits many sessions.
var MySQL80 = &MySQLServer{
	Image:   "mysql:8.0",
	Version: "8.0.",
	Env:     "OMNI_SHARED_MYSQL80_DSN",
	DevDSN:  "root:test@tcp(127.0.0.1:13380)/",
	Args:    []string{"--skip-log-bin", "--max-connections=1000"},
}

// OceanBase is OceanBase CE in MySQL mode, in the build the recorded OceanBase
// verdicts were read on (5.7.25-OceanBase_CE-v4.3.5.6). Its logins name the
// tenant, as root@test, and the restricted account logs in the same way.
var OceanBase = &MySQLServer{
	Image:   OceanBaseImage,
	Version: "5.7.25-OceanBase_CE-v4.3.5.6",
	Env:     "OMNI_SHARED_OCEANBASE_DSN",
	DevDSN:  "root@test:@tcp(127.0.0.1:12881)/",
	start:   StartOceanBase,
	open:    OpenOceanBase,
}

func (s *MySQLServer) openDB(dsn string) (*sql.DB, error) {
	if s.open != nil {
		return s.open(dsn)
	}
	return sql.Open("mysql", dsn)
}

func ping(ctx context.Context, dsn string) (*sql.DB, string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, "", err
	}
	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		db.Close()
		return nil, "", err
	}
	return db, version, nil
}

// admin returns a connection with the privileges to set up databases, finding or
// starting the server on the first call of this process.
func (s *MySQLServer) admin(ctx context.Context) (*sql.DB, error) {
	s.once.Do(func() { s.err = s.find(ctx) })
	if s.err != nil {
		return nil, s.err
	}
	db, err := s.openDB(s.config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func (s *MySQLServer) find(ctx context.Context) error {
	dsn := os.Getenv(s.Env)
	if dsn != "" {
		db, _, err := ping(ctx, dsn)
		if err != nil {
			return fmt.Errorf("server at $%s does not answer: %w", s.Env, err)
		}
		db.Close()
	} else if found := s.findDev(ctx); found != "" {
		dsn = found
	} else {
		var err error
		if dsn, err = s.startContainer(ctx); err != nil {
			return err
		}
	}
	if err := waitReady(ctx, dsn); err != nil {
		return err
	}
	config, err := gomysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("parse DSN: %w", err)
	}
	config.DBName = ""
	config.MultiStatements = false
	s.config = config
	return s.setUpAccount(ctx)
}

// findDev returns the DSN of the development server if one answers quickly and
// runs the version this server must, and otherwise "".
func (s *MySQLServer) findDev(ctx context.Context) string {
	config, err := gomysql.ParseDSN(s.DevDSN)
	if err != nil {
		return ""
	}
	config.Timeout = time.Second
	dsn := config.FormatDSN()
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	db, version, err := ping(pingCtx, dsn)
	if err != nil {
		return ""
	}
	db.Close()
	if !strings.HasPrefix(version, s.Version) {
		return ""
	}
	return s.DevDSN
}

// waitReady pings dsn until the server answers, for up to three minutes.
// OceanBase takes connections before its tenant can plan, so a first statement
// can fail while the server finishes starting; MySQL answers on the first try.
func waitReady(ctx context.Context, dsn string) error {
	deadline := time.Now().Add(3 * time.Minute)
	for {
		db, _, err := ping(ctx, dsn)
		if err == nil {
			db.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("server does not answer: %w", err)
		}
		time.Sleep(2 * time.Second)
	}
}

func (s *MySQLServer) startContainer(ctx context.Context) (string, error) {
	container, dsn, err := s.run(ctx)
	if container != nil {
		s.container = container
	}
	return dsn, err
}

// run starts a testcontainer of the server and returns it with an admin DSN. On
// failure it still returns whatever container it created, possibly a nil
// pointer, which testcontainers.TerminateContainer accepts.
func (s *MySQLServer) run(ctx context.Context) (testcontainers.Container, string, error) {
	if s.start != nil {
		return s.start(ctx, s.Image)
	}
	container, err := tcmysql.Run(ctx, s.Image,
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
		testcontainers.WithCmdArgs(s.Args...),
	)
	if err != nil {
		return container, "", fmt.Errorf("start %s: %w", s.Image, err)
	}
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return container, "", fmt.Errorf("connection string: %w", err)
	}
	return container, dsn, nil
}

func (s *MySQLServer) setUpAccount(ctx context.Context) error {
	db, err := s.openDB(s.config.FormatDSN())
	if err != nil {
		return err
	}
	defer db.Close()
	escaped := strings.ReplaceAll(databasePrefix, "_", `\_`) + "%"
	for _, stmt := range []string{
		fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", account, accountSecret),
		fmt.Sprintf("GRANT ALL ON `%s`.* TO '%s'@'%%'", escaped, account),
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

// Terminate stops a testcontainer this process started. It leaves a server the
// process found alone.
func (s *MySQLServer) Terminate() {
	if s.container != nil {
		_ = testcontainers.TerminateContainer(s.container)
	}
}

// Database is a test process's own database on a shared server.
type Database struct {
	// DB is a pool of connections that start in the database, on the restricted
	// account, pinned to one connection so that USE and session state hold across
	// statements as on a server of its own.
	DB   *sql.DB
	Name string
}

var labelChars = regexp.MustCompile(`[^a-z0-9_]`)

// Database creates a database named for this process and label, and returns it.
// The label says whose it is when someone reads SHOW DATABASES on the server.
func (s *MySQLServer) Database(ctx context.Context, label string) (*Database, error) {
	admin, err := s.admin(ctx)
	if err != nil {
		return nil, err
	}
	defer admin.Close()
	dropStale(ctx, admin)
	label = labelChars.ReplaceAllString(strings.ToLower(label), "_")
	name := fmt.Sprintf("%s%s_%d_%d", databasePrefix, label, time.Now().Unix(), os.Getpid())
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+name); err != nil {
		return nil, fmt.Errorf("create database %s: %w", name, err)
	}
	config := s.config.Clone()
	config.User, config.Passwd, config.DBName = account, accountSecret, name
	// OceanBase names the tenant in the login, as root@test, and the account
	// lives in the admin's tenant.
	if _, tenant, ok := strings.Cut(s.config.User, "@"); ok {
		config.User += "@" + tenant
	}
	config.ParseTime, config.MultiStatements = true, true
	db, err := s.openDB(config.FormatDSN())
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

// dropStale drops the databases of processes long gone. A failure costs only a
// leftover database, so it is ignored, as is a race with another process.
func dropStale(ctx context.Context, admin *sql.DB) {
	rows, err := admin.QueryContext(ctx, "SHOW DATABASES LIKE '"+strings.ReplaceAll(databasePrefix, "_", `\\_`)+"%'")
	if err != nil {
		return
	}
	var names []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			names = append(names, name)
		}
	}
	rows.Close()
	for _, name := range names {
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}
		started, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
		if err == nil && time.Since(time.Unix(started, 0)) > staleAfter {
			_, _ = admin.ExecContext(ctx, "DROP DATABASE IF EXISTS "+name)
		}
	}
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
