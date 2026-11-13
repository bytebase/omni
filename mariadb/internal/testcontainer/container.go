// Package testcontainer is how mariadb's tests reach MariaDB. Open and its kin
// give a test a connection through a tape (enginetest/tape), which replays the
// pinned server's answers and starts nothing; in a recording (OMNI_RECORD=1)
// the connection is to the pinned server itself. The compat judge, which runs
// only in a recording, asks the server directly, through Pinned.
//
// The server is the one testdata/pins.tsv names: one testcontainer per test
// binary, started by the first test that needs it and removed by Terminate,
// which the package's TestMain calls.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Image is the pinned MariaDB, as pins.tsv names it, for the packages whose
// tests do not replay tapes yet and start it themselves.
const Image = "mariadb:11.8.8"

// Pins is mariadb/testdata/pins.tsv.
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
	// Session holds system variables set when the connection opens, such as
	// max_statement_time. The tape does not see them being set.
	Session map[string]string
}

// Open returns t's connection to MariaDB, as root with the database test
// selected: a replay of its tape, or the pinned server in a recording. Each
// call is a connection of its own at session defaults, and the pool holds
// that one connection, so USE and user variables hold across statements.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return OpenWith(t, Options{Database: "test"})
}

// OpenWith is Open with other connection settings.
func OpenWith(t testing.TB, o Options) *sql.DB {
	t.Helper()
	return tape.Open(t, &tape.Server{Pins: Pins, Live: o.live, Codec: Codec{}}, "default")
}

// OpenEmpty returns t's connection to a server that holds nothing a test left
// behind, with the database test selected and empty, as a server of its own
// would be. In a recording every database, account and foreign server that
// the image did not start with is dropped first and test is created again, on
// a connection the tape does not see, so what a test sends never depends on
// what ran before it. A request may hold several statements, and dates come
// back as time.Time.
func OpenEmpty(t testing.TB) *sql.DB {
	t.Helper()
	return OpenEmptyWith(t, Options{Database: "test", MultiStatements: true, ParseTime: true})
}

// OpenEmptyWith is OpenEmpty with other connection settings.
func OpenEmptyWith(t testing.TB, o Options) *sql.DB {
	t.Helper()
	Reset(t, empty)
	return OpenWith(t, o)
}

// Reset runs fn on the pinned server, in a recording only: it puts the server
// into the state the requests after it start from, such as dropping what an
// earlier test left, so a replay needs nothing. fn's connection is not
// recorded, and has no database selected.
func Reset(t testing.TB, fn func(ctx context.Context, db *sql.DB) error) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	pin, err := enginetest.LoadPin(Pins, "default")
	if err != nil {
		t.Fatalf("%v", err)
	}
	admin := sql.OpenDB(Options{}.live(t, pin))
	defer admin.Close()
	if err := fn(context.Background(), admin); err != nil {
		t.Fatalf("reset MariaDB: %v", err)
	}
}

// empty drops what tests created on the server and creates the database test
// again. What the image started with stays: its databases and accounts were
// listed when the server came up.
func empty(ctx context.Context, db *sql.DB) error {
	text := func(s string) string {
		return "'" + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", "''") + "'"
	}
	name := func(s string) string { return "`" + strings.ReplaceAll(s, "`", "``") + "`" }
	var stmts []string
	databases, err := queryRows(ctx, db, "SELECT SCHEMA_NAME, '' FROM information_schema.SCHEMATA ORDER BY 1")
	if err != nil {
		return err
	}
	for _, d := range databases {
		if !shared.databases[strings.ToLower(d[0])] || strings.EqualFold(d[0], "test") {
			stmts = append(stmts, "DROP DATABASE IF EXISTS "+name(d[0]))
		}
	}
	accounts, err := queryRows(ctx, db, "SELECT User, Host FROM mysql.global_priv ORDER BY 1, 2")
	if err != nil {
		return err
	}
	for _, a := range accounts {
		switch {
		case shared.accounts[a]:
		case a[1] == "":
			// A role is an account with no host.
			stmts = append(stmts, "DROP ROLE IF EXISTS "+text(a[0]))
		default:
			stmts = append(stmts, "DROP USER IF EXISTS "+text(a[0])+"@"+text(a[1]))
		}
	}
	servers, err := queryRows(ctx, db, "SELECT Server_name, '' FROM mysql.servers ORDER BY 1")
	if err != nil {
		return err
	}
	for _, s := range servers {
		stmts = append(stmts, "DROP SERVER IF EXISTS "+name(s[0]))
	}
	for _, stmt := range append(stmts, "CREATE DATABASE test") {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

// queryRows reads a result of two text columns.
func queryRows(ctx context.Context, db *sql.DB, query string) ([][2]string, error) {
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

// liveTimeout bounds one statement on the live server, so one that hangs
// fails its test and not the run.
const liveTimeout = time.Minute

// live connects to the pinned server, starting it on the first call. It is
// called only in a recording. Arguments are inlined, because a tape holds no
// driver-level prepared statement.
func (o Options) live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	cfg := config(t, pin)
	cfg.DBName, cfg.MultiStatements, cfg.ParseTime = o.Database, o.MultiStatements, o.ParseTime
	if len(o.Session) > 0 {
		cfg.Params = maps.Clone(o.Session)
	}
	c, err := mysql.NewConnector(cfg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return c
}

// Pinned returns a connection pool to the server pin names, as root with the
// database test selected. Nothing it is asked is recorded: it is for the
// compat judge, which runs only in a recording.
func Pinned(t testing.TB, pin enginetest.Pin) *sql.DB {
	t.Helper()
	db := sql.OpenDB(Options{Database: "test"}.live(t, pin))
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// rootPassword is the root password of the server this package starts.
const rootPassword = "test"

var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	config    *mysql.Config
	// databases and accounts are what the image started with, which empty
	// leaves: database names in lower case, and accounts as user and host.
	databases map[string]bool
	accounts  map[[2]string]bool
	err       error
}

// config returns how to connect to the server pin names, starting the server
// on the first call. It fails t unless the container runs the pinned digest
// and reports the pinned release.
func config(t testing.TB, pin enginetest.Pin) *mysql.Config {
	t.Helper()
	shared.once.Do(func() {
		shared.pin = pin
		shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start MariaDB: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started MariaDB at %s, not %s", shared.pin, pin)
	}
	return shared.config.Clone()
}

func start(ctx context.Context, pin enginetest.Pin) error {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MARIADB_ROOT_PASSWORD": rootPassword,
				"MARIADB_DATABASE":      "test",
			},
			// The entrypoint runs a first server, without the network, to
			// set the database up; the port opens with the second.
			WaitingFor: wait.ForListeningPort("3306/tcp").WithStartupTimeout(2 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		return err
	}
	shared.container = c
	if err := CheckDigest(ctx, c, pin.Digest); err != nil {
		return err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return err
	}
	port, err := c.MappedPort(ctx, "3306")
	if err != nil {
		return err
	}
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = rootPassword
	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port.Port()
	// Arguments are inlined, so a statement and its arguments are one request.
	cfg.InterpolateParams = true
	cfg.ReadTimeout, cfg.WriteTimeout = liveTimeout, liveTimeout
	shared.config = cfg

	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		return err
	}
	db := sql.OpenDB(connector)
	defer db.Close()
	// The port listens a moment before the server accepts connections.
	deadline := time.Now().Add(90 * time.Second)
	for {
		err = db.PingContext(ctx)
		if err == nil || time.Now().After(deadline) {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return fmt.Errorf("MariaDB does not answer: %w", err)
	}
	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return err
	}
	// The source tag is mariadb-<release>, and VERSION() <release>-MariaDB-<build>.
	if release := strings.TrimPrefix(pin.Tag, "mariadb-"); !strings.HasPrefix(version, release+"-MariaDB") {
		return fmt.Errorf("the server reports %q, but pins.tsv pins MariaDB %s", version, pin.Tag)
	}
	databases, err := queryRows(ctx, db, "SELECT SCHEMA_NAME, '' FROM information_schema.SCHEMATA")
	if err != nil {
		return err
	}
	shared.databases = map[string]bool{}
	for _, d := range databases {
		shared.databases[strings.ToLower(d[0])] = true
	}
	accounts, err := queryRows(ctx, db, "SELECT User, Host FROM mysql.global_priv")
	if err != nil {
		return err
	}
	shared.accounts = map[[2]string]bool{}
	for _, a := range accounts {
		shared.accounts[a] = true
	}
	return nil
}

// Terminate removes the server, if this run started one.
func Terminate() {
	if shared.container != nil {
		_ = testcontainers.TerminateContainer(shared.container)
	}
}

// CheckDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func CheckDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the MariaDB container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the MariaDB image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the MariaDB container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

// Codec renders go-sql-driver's *mysql.MySQLError as "<number> <sqlstate>
// <message>", with - for a missing SQLSTATE.
//
// Every syntax error of MariaDB's parser opens with the same 124 characters,
// and a parser package's tapes hold tens of thousands, so a tape holds "…"
// for them: "1064 42000 …near 'x' at line 1".
type Codec struct{}

const (
	syntaxErrorPreamble = "You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near "
	syntaxErrorShort    = "…near "
)

func (Codec) EncodeError(err error) (string, bool) {
	var me *mysql.MySQLError
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
		return fmt.Errorf("tape: a MariaDB error %q", text)
	}
	n, err := strconv.ParseUint(f[0], 10, 16)
	if err != nil {
		return fmt.Errorf("tape: a MariaDB error %q", text)
	}
	me := &mysql.MySQLError{Number: uint16(n), Message: f[2]}
	if rest, ok := strings.CutPrefix(me.Message, syntaxErrorShort); ok {
		me.Message = syntaxErrorPreamble + rest
	}
	if f[1] != "-" {
		copy(me.SQLState[:], f[1])
	}
	return me
}
