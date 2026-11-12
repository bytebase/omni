// Package testcontainer is how oceanbase's tests reach OceanBase. Open and its
// kin give a test a connection through a tape (enginetest/tape), which replays
// the pinned server's answers and starts nothing; in a recording
// (OMNI_RECORD=1) the connection is to the pinned server itself. The packages
// not yet on tapes ask the server directly, through DSN and Database.
//
// The server is one testcontainer per test binary, started by the first test
// that needs it and removed by Terminate, which the package's TestMain calls.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is oceanbase/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Options is how a test's connection differs from the plain one.
type Options struct {
	// Database is selected when the connection opens. Empty selects none.
	Database string
	// MultiStatements lets one request hold several statements.
	MultiStatements bool
	// ParseTime returns DATE and DATETIME values as time.Time.
	ParseTime bool
}

// Open returns t's connection to OceanBase, as root in the test tenant with no
// database selected: a replay of its tape, or the pinned server in a
// recording. Each call is a connection of its own at session defaults, and
// the pool holds that one connection, so USE and user variables hold across
// statements.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return OpenWith(t, Options{})
}

// OpenWith is Open with other connection settings.
func OpenWith(t testing.TB, o Options) *sql.DB {
	t.Helper()
	return tape.Open(t, &tape.Server{Pins: Pins, Live: o.live, Codec: Codec{}}, "default")
}

// OpenDatabase returns t's connection with a database named for label
// selected and empty, as a server of its own would be. In a recording the
// database is dropped and created first, on a connection the tape does not
// see, so what a test sends never depends on what ran before it. A request
// may hold several statements, and dates come back as time.Time.
func OpenDatabase(t testing.TB, label string) *sql.DB {
	t.Helper()
	name := databasePrefix + labelChars.ReplaceAllString(strings.ToLower(label), "_")
	if enginetest.Recording(t) {
		pin, err := enginetest.LoadPin(Pins, "default")
		if err != nil {
			t.Fatalf("%v", err)
		}
		admin := sql.OpenDB(Options{}.live(t, pin))
		defer admin.Close()
		for _, stmt := range []string{"DROP DATABASE IF EXISTS " + name, "CREATE DATABASE " + name} {
			if _, err := admin.ExecContext(context.Background(), stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
	}
	return OpenWith(t, Options{Database: name, MultiStatements: true, ParseTime: true})
}

// liveTimeout bounds one statement on the live server, so one that hangs
// fails its test and not the run.
const liveTimeout = time.Minute

// live connects to the pinned server, starting it on the first call. It is
// called only in a recording. Arguments are inlined, because a tape holds no
// driver-level prepared statement. The retrying connector sits below the
// tape: an answer OceanBase gives while it is between leaders is asked
// again, and never recorded.
func (o Options) live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	ctx := context.Background()
	dsn, err := OceanBase.DSN(ctx)
	if err != nil {
		t.Fatalf("OceanBase: %v", err)
	}
	if err := OceanBase.CheckPin(ctx, pin); err != nil {
		t.Fatalf("%v", err)
	}
	cfg, err := gomysql.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("%v", err)
	}
	cfg.DBName, cfg.MultiStatements, cfg.ParseTime = o.Database, o.MultiStatements, o.ParseTime
	cfg.InterpolateParams = true
	cfg.ReadTimeout, cfg.WriteTimeout = liveTimeout, liveTimeout
	c, err := gomysql.NewConnector(cfg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return retryConnector{c}
}

// Codec renders go-sql-driver's *mysql.MySQLError as "<number> <sqlstate>
// <message>", with - for a missing SQLSTATE.
//
// Every syntax error of OceanBase's parser opens with the same 122
// characters, and a parser package's tapes hold tens of thousands, so a
// tape holds "…" for them: "1064 42000 …near 'x' at line 1".
type Codec struct{}

const (
	syntaxErrorPreamble = "You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use near "
	syntaxErrorShort    = "…near "
)

func (Codec) EncodeError(err error) (string, bool) {
	var me *gomysql.MySQLError
	if !errors.As(err, &me) {
		return "", false
	}
	state := string(me.SQLState[:])
	if me.SQLState == [5]byte{} {
		state = "-"
	}
	message := me.Message
	if rest, ok := strings.CutPrefix(message, syntaxErrorPreamble); ok {
		message = syntaxErrorShort + rest
	}
	return fmt.Sprintf("%d %s %s", me.Number, state, message), true
}

func (Codec) DecodeError(text string) error {
	f := strings.SplitN(text, " ", 3)
	if len(f) != 3 {
		return fmt.Errorf("tape: a MySQL error %q", text)
	}
	n, err := strconv.ParseUint(f[0], 10, 16)
	if err != nil {
		return fmt.Errorf("tape: a MySQL error %q", text)
	}
	me := &gomysql.MySQLError{Number: uint16(n), Message: f[2]}
	if rest, ok := strings.CutPrefix(me.Message, syntaxErrorShort); ok {
		me.Message = syntaxErrorPreamble + rest
	}
	if f[1] != "-" {
		copy(me.SQLState[:], f[1])
	}
	return me
}

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

// oceanBaseStart is how the image's start script starts the observer in SLIM
// mode (fastboot in /root/boot/start.sh): it daemonizes and reads the
// configuration and store under /root/demo.
var oceanBaseStart = []string{"bash", "-c", "cd /root/demo && LD_LIBRARY_PATH=/root/demo/lib bin/observer"}

// OceanBase is the package's OceanBase CE in MySQL mode. Tests log in as root
// in the test tenant, so they may create accounts, use fixed database names and
// change server-wide settings.
var OceanBase = &oceanBase{}

type oceanBase struct {
	once      sync.Once
	dsn       string
	container testcontainers.Container
	err       error

	pinOnce sync.Once
	pinErr  error
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

// Alive reports whether the observer process is running. A statement can
// crash it, and the container outlives it: SLIM does not start it again.
func (o *oceanBase) Alive(ctx context.Context) (bool, error) {
	if o.container == nil {
		return false, errors.New("no OceanBase container: call DSN first")
	}
	code, _, err := o.container.Exec(ctx, []string{"pgrep", "-x", "observer"})
	if err != nil {
		return false, fmt.Errorf("look for the observer: %w", err)
	}
	return code == 0, nil
}

// Restart starts the observer again after a crash, from the store it left,
// and waits until it answers. The container and its mapped port stay, so a
// pool opened on DSN before the crash reconnects by itself.
func (o *oceanBase) Restart(ctx context.Context) error {
	if o.container == nil {
		return errors.New("no OceanBase container: call DSN first")
	}
	code, out, err := o.container.Exec(ctx, oceanBaseStart)
	if err == nil && code != 0 {
		b, _ := io.ReadAll(out)
		err = fmt.Errorf("exit %d: %s", code, b)
	}
	if err != nil {
		return fmt.Errorf("start the observer: %w", err)
	}
	return waitReady(ctx, o.dsn)
}

// CheckPin fails unless the server this process started is the one pin
// names: the image, its digest, and the release it reports. The check runs
// once.
func (o *oceanBase) CheckPin(ctx context.Context, pin enginetest.Pin) error {
	o.pinOnce.Do(func() {
		if pin.Image != OceanBaseImage {
			o.pinErr = fmt.Errorf("%s pins %s, but this package starts %s", Pins, pin.Image, OceanBaseImage)
			return
		}
		if o.pinErr = o.CheckDigest(ctx, pin.Digest); o.pinErr != nil {
			return
		}
		db, err := OpenOceanBase(o.dsn)
		if err != nil {
			o.pinErr = err
			return
		}
		defer db.Close()
		var version string
		if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
			o.pinErr = fmt.Errorf("read the OceanBase version: %w", err)
		} else if version != OceanBaseVersion {
			o.pinErr = fmt.Errorf("the server reports %q, want %q", version, OceanBaseVersion)
		}
	})
	return o.pinErr
}

// CheckDigest fails unless the server this process started runs the image
// with the pinned repository digest. A tag can be pushed again; a digest
// cannot.
func (o *oceanBase) CheckDigest(ctx context.Context, digest string) error {
	if o.container == nil {
		return errors.New("no OceanBase container: call DSN first")
	}
	info, err := o.container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the OceanBase container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the OceanBase image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the OceanBase container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
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

// databasePrefix starts the name of every database Database creates.
const databasePrefix = "omni_t_"

// Database is a test's own database on OceanBase.
type Database struct {
	// DB is a pool of root connections that start in the database, pinned to one
	// connection so that USE and session state hold across statements.
	DB   *sql.DB
	Name string
}

var labelChars = regexp.MustCompile(`[^a-z0-9_]`)

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
