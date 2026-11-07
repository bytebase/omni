// Package testcontainer starts the servers oceanbase's oracle tests run
// against: one testcontainer of each per test binary, started by the first test
// that asks and removed by Terminate, which the package's TestMain calls.
//
// MySQL80 answers for the grammar OceanBase inherited from MySQL. A test gets a
// database of its own on it, reached through an account that holds privileges
// only on the databases this package creates, so a statement a corpus runs (SET
// GLOBAL, SHUTDOWN, ALTER USER root) cannot change the server under the tests
// that follow. Those statements answer 1227 or 1044, which a verdict that looks
// for 1064 still reads as parsed. OceanBase answers for OceanBase itself, as
// root.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

// OceanBaseImage is the build the recorded OceanBase verdicts were read on. The
// 4.3.5-lts tag moves with each patch release, so this names the patch.
const OceanBaseImage = "oceanbase/oceanbase-ce:4.3.5.6-106000012026040916"

// OceanBaseVersion is what SELECT VERSION() answers on OceanBaseImage.
const OceanBaseVersion = "5.7.25-OceanBase_CE-v4.3.5.6"

// oceanBaseSocket is the observer's unix socket in a SLIM container.
const oceanBaseSocket = "/root/demo/run/sql.sock"

// oceanBaseSetup runs inside the container once it boots.
//
// SLIM caps the data file at 1 GiB, and the test tenant's flushes then fail with
// OB_SERVER_OUTOF_DISK_SPACE partway through the review walk. External tables
// read files only under secure_file_priv, which the server lets a client set
// only over its unix socket.
var oceanBaseSetup = [][]string{
	{"obclient", "-h127.1", "-P2881", "-uroot@sys", "-e", "ALTER SYSTEM SET datafile_maxsize = '8G'"},
	{"obclient", "-S", oceanBaseSocket, "-uroot@test", "-e", "SET GLOBAL secure_file_priv = '/'"},
}

// OceanBase is the package's OceanBase CE in MySQL mode. Tests log in as root
// in the test tenant, so they may create accounts, use fixed database names and
// change server-wide settings.
var OceanBase = &oceanBase{}

type oceanBase struct {
	once      sync.Once
	dsn       string
	container testcontainers.Container
	err       error
}

// DSN returns root's DSN, naming no database, once the server answers, starting
// it on the first call of this process. Open it with OpenOceanBase.
func (o *oceanBase) DSN(ctx context.Context) (string, error) {
	o.once.Do(func() {
		o.container, o.dsn, o.err = startOceanBase(ctx)
		if o.err == nil {
			o.err = waitReady(ctx, o.dsn)
		}
	})
	return o.dsn, o.err
}

// Terminate removes the server, if this process started one.
func (o *oceanBase) Terminate() {
	if o.container != nil {
		_ = testcontainers.TerminateContainer(o.container)
	}
}

// startOceanBase boots OceanBaseImage in SLIM mode and returns the container
// with a DSN for root in the test tenant, which has no password.
//
// SLIM starts the observer from the store the image ships, already holding the
// test tenant: about 30 s and 3.5 GiB under the oracle tests, where the default
// MINI mode deploys a cluster and creates the tenant, which takes minutes and
// 8 GiB. The deployer refuses fewer than 20000 open files.
func startOceanBase(ctx context.Context) (testcontainers.Container, string, error) {
	c, err := testcontainers.Run(ctx, OceanBaseImage,
		testcontainers.WithEnv(map[string]string{"MODE": "SLIM"}),
		testcontainers.WithExposedPorts("2881/tcp"),
		testcontainers.WithHostConfigModifier(func(hc *container.HostConfig) {
			hc.Ulimits = []*units.Ulimit{{Name: "nofile", Soft: 65536, Hard: 65536}}
		}),
		testcontainers.WithWaitStrategy(wait.ForLog("boot success!").WithStartupTimeout(10*time.Minute)),
	)
	if err != nil {
		return c, "", fmt.Errorf("start %s: %w", OceanBaseImage, err)
	}
	for _, setup := range oceanBaseSetup {
		code, out, err := c.Exec(ctx, setup)
		if err == nil && code != 0 {
			b, _ := io.ReadAll(out)
			err = fmt.Errorf("exit %d: %s", code, b)
		}
		if err != nil {
			return c, "", fmt.Errorf("%v: %w", setup, err)
		}
	}
	host, err := c.Host(ctx)
	if err != nil {
		return c, "", err
	}
	port, err := c.MappedPort(ctx, "2881/tcp")
	if err != nil {
		return c, "", err
	}
	return c, fmt.Sprintf("root@test:@tcp(%s:%s)/", host, port.Port()), nil
}

// waitReady pings dsn until the server answers, for up to three minutes:
// OceanBase takes connections before its tenant can plan.
func waitReady(ctx context.Context, dsn string) error {
	deadline := time.Now().Add(3 * time.Minute)
	for {
		db, err := OpenOceanBase(dsn)
		if err != nil {
			return err
		}
		var version string
		err = db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)
		db.Close()
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("%s does not answer: %w", OceanBaseImage, err)
		}
		time.Sleep(2 * time.Second)
	}
}

// Database creates a database named for label on OceanBase and returns it, with
// a pool of root connections that start in it, pinned to one connection so that
// USE and session state hold across statements.
func (o *oceanBase) Database(ctx context.Context, label string) (*Database, error) {
	dsn, err := o.DSN(ctx)
	if err != nil {
		return nil, err
	}
	name := databasePrefix + labelChars.ReplaceAllString(strings.ToLower(label), "_")
	admin, err := OpenOceanBase(dsn)
	if err != nil {
		return nil, err
	}
	defer admin.Close()
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+name); err != nil {
		return nil, fmt.Errorf("create database %s: %w", name, err)
	}
	config, err := gomysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	config.DBName = name
	config.ParseTime, config.MultiStatements = true, true
	db, err := OpenOceanBase(config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return &Database{DB: db, Name: name}, nil
}

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

// MySQL80 is mysql:8.0. The binary log stays off so the restricted account can
// create routines (1419 otherwise).
var MySQL80 = &Server{Image: "mysql:8.0", Args: []string{"--skip-log-bin"}}

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

// OceanBaseTransient reports an error OceanBase answers while it is between
// leaders or starved for CPU, which the same statement clears once the server
// settles. On a busy box (load average 380 on 20 cores, 2026-11-05) the root
// service is revoked and taken back within a minute, and a log stream changes
// leader.
//
// A verdict must never read one of these as an answer: engineAccepts-style
// helpers take any code but 1064 to mean the statement parsed.
func OceanBaseTransient(err error) bool {
	var me *gomysql.MySQLError
	if errors.As(err, &me) {
		switch me.Number {
		// 4012 is a timeout, 4550 "rootservice is shutdown", 4638 "The
		// RootServer is not the master", 4767 the tenant's log stream between
		// leaders; the rest are the walk's list from 2026-10-28.
		case 4012, 4023, 4038, 4109, 4550, 4638, 4721, 4767, 6005:
			return true
		}
	}
	return false
}

// retryWindow is how long a statement keeps being retried, retryPause the wait
// between tries. Tests shorten them.
var (
	retryWindow = 3 * time.Minute
	retryPause  = 2 * time.Second
)

// OpenOceanBase opens dsn through a connector that runs each statement again
// while OceanBase answers a transient error, for up to three minutes, so no test
// has to wrap its own statements. A statement still failing after that returns
// the last error as a plain error.
func OpenOceanBase(dsn string) (*sql.DB, error) {
	cfg, err := gomysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	connector, err := gomysql.NewConnector(cfg)
	if err != nil {
		return nil, err
	}
	return sql.OpenDB(retryConnector{connector}), nil
}

type retryConnector struct{ driver.Connector }

func (c retryConnector) Connect(ctx context.Context) (driver.Conn, error) {
	var conn driver.Conn
	err := retry(ctx, func() error {
		var err error
		conn, err = c.Connector.Connect(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &retryConn{conn}, nil
}

// retry runs f until it returns something other than a transient error. A
// transient error that outlasts the window comes back as a plain error, not a
// *mysql.MySQLError: a verdict helper reading codes then finds no server answer
// and fails, where 4550 would have read as "parsed".
func retry(ctx context.Context, f func() error) error {
	deadline := time.Now().Add(retryWindow)
	for {
		err := f()
		if !OceanBaseTransient(err) {
			return err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("OceanBase stayed unavailable for %v: %v", retryWindow, err)
		}
		select {
		case <-ctx.Done():
			return err
		case <-time.After(retryPause):
		}
	}
}

// retryConn retries Exec and Query and passes the rest of the driver's
// interfaces through, so database/sql treats it as the mysql connection.
type retryConn struct{ driver.Conn }

func (c *retryConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	var res driver.Result
	err := retry(ctx, func() error {
		var err error
		res, err = c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
		return err
	})
	return res, err
}

func (c *retryConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	var rows driver.Rows
	err := retry(ctx, func() error {
		var err error
		rows, err = c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
		return err
	})
	return rows, err
}

func (c *retryConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.Conn.(driver.ConnPrepareContext).PrepareContext(ctx, query)
}

func (c *retryConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
}

func (c *retryConn) Ping(ctx context.Context) error {
	return retry(ctx, func() error { return c.Conn.(driver.Pinger).Ping(ctx) })
}

func (c *retryConn) ResetSession(ctx context.Context) error {
	return c.Conn.(driver.SessionResetter).ResetSession(ctx)
}

func (c *retryConn) IsValid() bool { return c.Conn.(driver.Validator).IsValid() }

func (c *retryConn) CheckNamedValue(nv *driver.NamedValue) error {
	return c.Conn.(driver.NamedValueChecker).CheckNamedValue(nv)
}
