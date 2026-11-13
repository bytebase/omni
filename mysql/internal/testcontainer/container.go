// Package testcontainer is how mysql's tests reach MySQL. Open and its kin
// give a test a connection through a tape (enginetest/tape), which replays
// the pinned server's answers and starts nothing; in a recording
// (OMNI_RECORD=1) the connection is to the pinned server itself. The compat
// judge, which runs only in a recording, asks the server directly
// (Server.Judge).
//
// Each row of mysql/testdata/pins.tsv is a server a test may ask for by its
// role. Each is one testcontainer per test binary, started by the first test
// that records on it and removed by Terminate, which the package's TestMain
// calls.
//
// A test that runs what a corpus holds (OpenDatabase) logs in as an account
// that holds privileges only on the databases this package creates, so a
// statement such as SET GLOBAL, SHUTDOWN or ALTER USER root cannot change the
// server under the tests recorded after it. Those statements answer 1227 or
// 1044, which a verdict that looks for 1064 still reads as parsed. A test
// that needs the server whole (OpenRoot) logs in as root.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is mysql/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Roles of pins.tsv, and the server each stands for.
const (
	// Default is mysql:8.0.46 as it starts: lower_case_table_names 0.
	Default = "default"
)

// servers holds the server of each role, started on first use. The binary
// log stays off, so the restricted account can create routines (1419
// otherwise).
var servers = map[string]*Server{
	Default: {role: Default, args: []string{"--skip-log-bin"}},
}

// Options is how a test's connection differs from the plain one.
type Options struct {
	// Role is the pins.tsv row of the server. Empty is Default.
	Role string
	// Database is selected when the connection opens. Empty selects none.
	Database string
	// MultiStatements lets one request hold several statements.
	MultiStatements bool
	// ParseTime returns DATE and DATETIME values as time.Time.
	ParseTime bool
	// restricted logs in as the account that reaches only this package's
	// databases.
	restricted bool
}

func (o Options) role() string {
	if o.Role == "" {
		return Default
	}
	return o.Role
}

// OpenRoot returns t's connection to MySQL as root: a replay of its tape, or
// the pinned server in a recording. Each call is a connection of its own at
// session defaults, and the pool holds that one connection, so USE and
// session state hold across statements.
func OpenRoot(t testing.TB, o Options) *sql.DB {
	t.Helper()
	server := servers[o.role()]
	if server == nil {
		t.Fatalf("testcontainer: no server of role %q", o.role())
	}
	return tape.Open(t, &tape.Server{Pins: Pins, Live: func(t testing.TB, pin enginetest.Pin) driver.Connector {
		return server.connector(t, pin, o)
	}, Codec: Codec{}}, o.role())
}

// OpenEmpty returns t's root connection to a server that holds nothing a test
// left behind, with the database test selected, for a test that names its own
// databases and accounts. In a recording every database and account that is
// not the server's own is dropped first and test is created, on a connection
// the tape does not see. A request may hold several statements, and dates
// come back as time.Time.
func OpenEmpty(t testing.TB) *sql.DB {
	t.Helper()
	return OpenEmptyWith(t, Options{})
}

// OpenEmptyWith is OpenEmpty on another server, or with other settings.
func OpenEmptyWith(t testing.TB, o Options) *sql.DB {
	t.Helper()
	Reset(t, o.role(), empty)
	if o.Database == "" {
		o.Database = "test"
	}
	o.MultiStatements, o.ParseTime = true, true
	return OpenRoot(t, o)
}

// OpenDatabase returns t's connection, on the restricted account, with a
// database named for label selected and empty, as a server of its own would
// be. In a recording the database is dropped and created first, on a
// connection the tape does not see, so what a test sends never depends on
// what ran before it. A request may hold several statements, and dates come
// back as time.Time.
func OpenDatabase(t testing.TB, label string) *sql.DB {
	t.Helper()
	name := DatabaseName(label)
	Reset(t, Default, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{"DROP DATABASE IF EXISTS " + name, "CREATE DATABASE " + name} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("%s: %w", stmt, err)
			}
		}
		return nil
	})
	return OpenRoot(t, Options{Database: name, MultiStatements: true, ParseTime: true, restricted: true})
}

// Reset runs fn as root on the server of role, in a recording only: it puts
// the server into the state the requests after it start from, such as
// dropping what an earlier test left, so a replay needs nothing. fn's
// connection is not recorded.
func Reset(t testing.TB, role string, fn func(ctx context.Context, db *sql.DB) error) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	pin, err := enginetest.LoadPin(Pins, role)
	if err != nil {
		t.Fatalf("%v", err)
	}
	server := servers[role]
	if server == nil {
		t.Fatalf("testcontainer: no server of role %q", role)
	}
	db := sql.OpenDB(server.connector(t, pin, Options{}))
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := fn(context.Background(), db); err != nil {
		t.Fatalf("reset MySQL (%s): %v", role, err)
	}
}

// DatabaseName is the name of the database OpenDatabase selects for label.
func DatabaseName(label string) string {
	return databasePrefix + labelChars.ReplaceAllString(strings.ToLower(label), "_")
}

const (
	// The restricted account and the prefix of its databases. The
	// underscores are escaped in the GRANT, where they would otherwise
	// match any character.
	account        = "omni_test"
	accountSecret  = "omni_test"
	databasePrefix = "omni_t_"
)

var labelChars = regexp.MustCompile(`[^a-z0-9_]`)

// systemDatabases are the databases a MySQL server starts with.
var systemDatabases = map[string]bool{"information_schema": true, "mysql": true, "performance_schema": true, "sys": true}

// systemAccounts are the accounts a MySQL server starts with, and the
// restricted account.
var systemAccounts = map[string]bool{"root": true, "mysql.infoschema": true, "mysql.session": true, "mysql.sys": true, account: true}

// empty drops what tests created on the server and creates the database test,
// and puts back the global settings a test may change.
func empty(ctx context.Context, db *sql.DB) error {
	var stmts []string
	databases, err := queryPairs(ctx, db, "SELECT SCHEMA_NAME, '' FROM information_schema.SCHEMATA")
	if err != nil {
		return err
	}
	for _, d := range databases {
		if !systemDatabases[strings.ToLower(d[0])] {
			stmts = append(stmts, "DROP DATABASE `"+strings.ReplaceAll(d[0], "`", "``")+"`")
		}
	}
	accounts, err := queryPairs(ctx, db, "SELECT User, Host FROM mysql.user")
	if err != nil {
		return err
	}
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }
	for _, a := range accounts {
		if !systemAccounts[a[0]] {
			stmts = append(stmts, "DROP USER "+quote(a[0])+"@"+quote(a[1]))
		}
	}
	stmts = append(stmts,
		"SET GLOBAL foreign_key_checks = 1",
		"SET GLOBAL show_gipk_in_create_table_and_information_schema = ON",
		"SET GLOBAL sql_generate_invisible_primary_key = OFF",
		"CREATE DATABASE test")
	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

func queryPairs(ctx context.Context, db *sql.DB, query string) ([][2]string, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var out [][2]string
	for rows.Next() {
		var r [2]string
		if err := rows.Scan(&r[0], &r[1]); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// Server is one pinned MySQL a test binary starts on first use.
type Server struct {
	role string
	args []string

	once      sync.Once
	pin       enginetest.Pin
	config    *gomysql.Config // root, no database
	container testcontainers.Container
	err       error
}

// liveTimeout bounds one statement on the live server, so one that hangs fails
// its test and not the run.
const liveTimeout = time.Minute

// connector connects to the pinned server, starting it on the first call. It
// is called only in a recording. Arguments are inlined, because a tape holds
// no driver-level prepared statement.
func (s *Server) connector(t testing.TB, pin enginetest.Pin, o Options) driver.Connector {
	t.Helper()
	if err := s.start(context.Background(), pin); err != nil {
		t.Fatalf("MySQL (%s): %v", s.role, err)
	}
	cfg := s.config.Clone()
	if o.restricted {
		cfg.User, cfg.Passwd = account, accountSecret
	}
	cfg.DBName, cfg.MultiStatements, cfg.ParseTime = o.Database, o.MultiStatements, o.ParseTime
	cfg.InterpolateParams = true
	cfg.ReadTimeout, cfg.WriteTimeout = liveTimeout, liveTimeout
	c, err := gomysql.NewConnector(cfg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return c
}

// start starts the server once, checks it is the pinned one, and creates the
// restricted account.
func (s *Server) start(ctx context.Context, pin enginetest.Pin) error {
	s.once.Do(func() {
		s.pin = pin
		s.err = s.boot(ctx, pin)
	})
	if s.err != nil {
		return s.err
	}
	if s.pin != pin {
		return fmt.Errorf("this run started %s, not %s", s.pin, pin)
	}
	return nil
}

func (s *Server) boot(ctx context.Context, pin enginetest.Pin) error {
	c, err := tcmysql.Run(ctx, pin.Image,
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
		testcontainers.WithCmdArgs(s.args...),
	)
	if c != nil {
		s.container = c
	}
	if err != nil {
		return fmt.Errorf("start %s: %w", pin.Image, err)
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		return err
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
	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return fmt.Errorf("read the MySQL version: %w", err)
	}
	if "mysql-"+version != pin.Tag {
		return fmt.Errorf("the server reports %q, but %s pins %s", version, Pins, pin.Tag)
	}
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

// checkDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
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
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the MySQL container runs an image with digests %v, but %s pins %s", img.RepoDigests, Pins, digest)
}

// Judge returns a connection the tape does not see, for the compat judge,
// whose answers go to compat.tsv: on the restricted account, with a database
// named for label selected and empty. It runs in a recording only.
func Judge(t testing.TB, pin enginetest.Pin, label string) *sql.DB {
	t.Helper()
	server := servers[pin.Role]
	if server == nil {
		t.Fatalf("testcontainer: no server of role %q", pin.Role)
	}
	name := DatabaseName(label)
	Reset(t, pin.Role, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{"DROP DATABASE IF EXISTS " + name, "CREATE DATABASE " + name} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("%s: %w", stmt, err)
			}
		}
		return nil
	})
	db := sql.OpenDB(server.connector(t, pin, Options{Database: name, restricted: true}))
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// Terminate removes every server this process started.
func Terminate() {
	for _, s := range servers {
		if s.container != nil {
			_ = testcontainers.TerminateContainer(s.container)
		}
	}
}

// Codec renders go-sql-driver's *mysql.MySQLError as "<number> <sqlstate>
// <message>", with - for a missing SQLSTATE.
//
// Every syntax error opens with the same 134 characters, and a parser
// package's tapes hold thousands, so a tape holds "…" for them:
// "1064 42000 …near 'x' at line 1".
type Codec struct{}

const (
	syntaxErrorPreamble = "You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near "
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
