// Package testcontainer is how starrocks's tests reach StarRocks: the
// pinned server, started only in a recording (OMNI_RECORD=1) to judge
// compat.tsv. It is one testcontainer per test binary, started by the first
// test that asks for it and removed by Terminate, which the package's
// TestMain calls.
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

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// database is where a recording's container keeps what the tests do not name,
// so a statement that selects no database still finds one. It is not part of
// the image, and ResetToInitial drops it.
const database = "omni_compat"

// Pins is starrocks/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Judge returns a connection to the pinned server, for a judge of
// compat.tsv, which runs only in a recording. Statements go to the server
// as written, byte for byte, more than one to a request.
func Judge(t testing.TB) *sql.Conn {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	pin, err := enginetest.LoadPin(Pins, "default")
	if err != nil {
		t.Fatalf("%v", err)
	}
	shared.once.Do(func() {
		shared.pin = pin
		shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start StarRocks: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started StarRocks at %s, not %s", shared.pin, pin)
	}
	c, err := mysql.NewConnector(shared.config.Clone())
	if err != nil {
		t.Fatalf("%v", err)
	}
	db := sql.OpenDB(c)
	t.Cleanup(func() { db.Close() })
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("connect to StarRocks: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	config    *mysql.Config
	// databases is what the image started with, which ResetToInitial leaves.
	databases map[string]bool
	err       error
}

func start(ctx context.Context, pin enginetest.Pin) error {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: pin.ImageRef(),
			// The image is built for amd64 only; elsewhere it runs emulated.
			ImagePlatform: "linux/amd64",
			ExposedPorts:  []string{"9030/tcp"},
			WaitingFor:    wait.ForListeningPort("9030/tcp").WithStartupTimeout(5 * time.Minute),
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
	port, err := c.MappedPort(ctx, "9030")
	if err != nil {
		return err
	}
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port.Port()
	// Arguments are inlined, so a statement and its arguments are one request.
	cfg.InterpolateParams = true
	// A statement that hangs fails the recording, not the run.
	cfg.ReadTimeout = 30 * time.Second
	cfg.WriteTimeout = 30 * time.Second

	db := sql.OpenDB(must(mysql.NewConnector(cfg)))
	defer db.Close()
	// The port listens a moment before the front end answers.
	deadline := time.Now().Add(5 * time.Minute)
	for {
		var one int
		err = db.QueryRowContext(ctx, "SELECT 1").Scan(&one)
		if err == nil || time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("StarRocks does not answer: %w", err)
	}
	var version string
	if err := db.QueryRowContext(ctx, "SELECT @@version_comment").Scan(&version); err != nil {
		return err
	}
	// 3.4.10-965be75
	if !strings.HasPrefix(version, pin.Tag+"-") {
		return fmt.Errorf("the server reports %q, but pins.tsv pins StarRocks %s", version, pin.Tag)
	}
	// The front end answers before a back end can store a table.
	if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+database); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, "USE "+database); err != nil {
		return err
	}
	waitDeadline := time.Now().Add(5 * time.Minute)
	for {
		_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS omni_ready (a INT) ENGINE=OLAP DUPLICATE KEY(a) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num" = "1")`)
		if err == nil {
			_, err = db.ExecContext(ctx, "DROP TABLE IF EXISTS omni_ready")
			if err != nil {
				return err
			}
			break
		}
		if time.Now().After(waitDeadline) {
			return fmt.Errorf("StarRocks cannot create a table: %w", err)
		}
		time.Sleep(2 * time.Second)
	}
	// The judge sends a probe and the statement as one request.
	cfg.MultiStatements = true
	shared.config = cfg
	return recordDatabases(ctx, db)
}

// recordDatabases notes the databases the image started with.
func recordDatabases(ctx context.Context, db *sql.DB) error {
	names, err := queryRows(ctx, db, "SHOW DATABASES")
	if err != nil {
		return err
	}
	shared.databases = map[string]bool{}
	for _, name := range names {
		if name == database {
			continue
		}
		shared.databases[strings.ToLower(name)] = true
	}
	return nil
}

func must(c driver.Connector, err error) driver.Connector {
	if err != nil {
		panic(err)
	}
	return c
}

// Terminate removes the server, if this run started one.
func Terminate() {
	if shared.container != nil {
		_ = testcontainers.TerminateContainer(shared.container)
	}
}

// CheckDigest proves the container runs the pinned image. The pull names
// image@digest, so a moved tag cannot select another image; this check reads
// what the daemon actually started, which is where the recorded answers came
// from.
func CheckDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the StarRocks container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the StarRocks image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the StarRocks container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

// Options is how a test's connection differs from the plain one.
type Options struct {
	// Database is selected when the connection opens. Empty selects none.
	Database string
	// MultiStatements lets one request hold several statements.
	MultiStatements bool
}

// Open returns t's connection to StarRocks: a replay of its tape, or the pinned
// server in a recording.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return OpenWith(t, Options{})
}

// OpenWith is Open with other connection settings.
func OpenWith(t testing.TB, o Options) *sql.DB {
	t.Helper()
	return tape.Open(t, &tape.Server{Pins: Pins, Live: o.live, Codec: Codec{}}, "default")
}

// live connects to the pinned server, starting it on the first call. It is
// called only in a recording. Arguments are inlined, because a tape holds no
// driver-level prepared statement.
func (o Options) live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	cfg := config(t, pin)
	cfg.DBName, cfg.MultiStatements = o.Database, o.MultiStatements
	c, err := mysql.NewConnector(cfg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return c
}

// config returns how to connect to the server pin names, starting the server on
// the first call. It fails t unless this run started the pinned server.
func config(t testing.TB, pin enginetest.Pin) *mysql.Config {
	t.Helper()
	shared.once.Do(func() {
		shared.pin = pin
		shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start StarRocks: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started StarRocks at %s, not %s", shared.pin, pin)
	}
	return shared.config.Clone()
}

func quoteName(s string) string { return "`" + strings.ReplaceAll(s, "`", "``") + "`" }

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
		t.Fatalf("reset StarRocks: %v", err)
	}
}

// ResetToInitial drops every database the pinned image did not start with, so
// a test that reads the server's database list sees a fixed set. It runs in a
// recording only.
func ResetToInitial(t testing.TB) {
	t.Helper()
	Reset(t, func(ctx context.Context, db *sql.DB) error {
		databases, err := queryRows(ctx, db, "SHOW DATABASES")
		if err != nil {
			return err
		}
		for _, name := range databases {
			if shared.databases[strings.ToLower(name)] {
				continue
			}
			if _, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+quoteName(name)); err != nil {
				return fmt.Errorf("drop %s: %w", name, err)
			}
		}
		return nil
	})
}

// queryRows reads the first column of every row.
func queryRows(ctx context.Context, db *sql.DB, query string) ([]string, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Codec renders go-sql-driver's *mysql.MySQLError as "<number> <sqlstate>
// <message>", with - for a missing SQLSTATE.
type Codec struct{}

func (Codec) EncodeError(err error) (string, bool) {
	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return "", false
	}
	state := string(me.SQLState[:])
	if me.SQLState == [5]byte{} {
		state = "-"
	}
	return fmt.Sprintf("%d %s %s", me.Number, state, me.Message), true
}

// madeUpRe matches what a new container makes up again: the tablet ID an
// EXPLAIN plan names, and the trace ID in an internal error. The tape holds a
// stand-in for each, so a second recording writes the same bytes.
var madeUpRe = regexp.MustCompile(`tabletList=(?:[0-9]+|omni_tablet_[0-9]+)|(?:\([0-9a-f]{12}\)|omni_trace_[0-9]+)\[INTERNAL_ERROR\]`)

// MadeUp matches made-up identifiers in the server's answers.
func (Codec) MadeUp() *regexp.Regexp { return madeUpRe }

// StandIn returns what stands for the nth distinct made-up identifier of a test.
func (Codec) StandIn(name string, n int) string {
	if strings.HasPrefix(name, "tabletList=") {
		return fmt.Sprintf("tabletList=omni_tablet_%d", n)
	}
	return fmt.Sprintf("omni_trace_%d[INTERNAL_ERROR]", n)
}

func (Codec) DecodeError(text string) error {
	f := strings.SplitN(text, " ", 3)
	if len(f) != 3 {
		return fmt.Errorf("tape: a StarRocks error %q", text)
	}
	n, err := strconv.ParseUint(f[0], 10, 16)
	if err != nil {
		return fmt.Errorf("tape: a StarRocks error %q", text)
	}
	me := &mysql.MySQLError{Number: uint16(n), Message: f[2]}
	if f[1] != "-" {
		copy(me.SQLState[:], f[1])
	}
	return me
}
