// Package testcontainer is how tidb's tests reach TiDB: through a tape
// (enginetest/tape) that replays the pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. That server is
// one testcontainer per test binary, started by the first test that records
// and removed by Terminate, which the package's TestMain calls.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
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

// Pins is tidb/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}

// Open returns t's connection to TiDB: a replay of its tape, or the pinned
// server in a recording. Each call is a connection of its own, at session
// defaults with the database test selected.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return tape.Open(t, server, "default")
}

// Reset runs fn on the pinned server before a test opens its connection,
// in a recording only: it puts the server into the state the test's tape
// starts from, so a replay needs nothing. fn's connection is not recorded.
func Reset(t testing.TB, fn func(ctx context.Context, db *sql.DB) error) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	pin, err := enginetest.LoadPin(Pins, "default")
	if err != nil {
		t.Fatalf("%v", err)
	}
	db := sql.OpenDB(live(t, pin))
	defer db.Close()
	if err := fn(context.Background(), db); err != nil {
		t.Fatalf("reset TiDB: %v", err)
	}
}

var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	config    *mysql.Config
	err       error
}

// live connects to the pinned server, starting it on the first call.
func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	shared.once.Do(func() {
		shared.pin = pin
		shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start TiDB: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started TiDB at %s, not %s", shared.pin, pin)
	}
	c, err := mysql.NewConnector(shared.config.Clone())
	if err != nil {
		t.Fatalf("%v", err)
	}
	return c
}

func start(ctx context.Context, pin enginetest.Pin) error {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{"4000/tcp"},
			WaitingFor:   wait.ForListeningPort("4000/tcp").WithStartupTimeout(2 * time.Minute),
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
	port, err := c.MappedPort(ctx, "4000")
	if err != nil {
		return err
	}
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port.Port()
	cfg.DBName = "test"
	cfg.MultiStatements = true
	// A tape holds each statement with its arguments: no prepared statements.
	cfg.InterpolateParams = true
	// A statement that hangs fails its test, not the run.
	cfg.ReadTimeout = 30 * time.Second
	cfg.WriteTimeout = 30 * time.Second
	shared.config = cfg

	db := sql.OpenDB(must(mysql.NewConnector(cfg)))
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
		return fmt.Errorf("TiDB does not answer: %w", err)
	}
	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return err
	}
	if !strings.HasSuffix(version, "-TiDB-"+pin.Tag) {
		return fmt.Errorf("the server reports %q, but pins.tsv pins TiDB %s", version, pin.Tag)
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

// CheckDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func CheckDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the TiDB container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the TiDB image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the TiDB container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
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

func (Codec) DecodeError(text string) error {
	f := strings.SplitN(text, " ", 3)
	if len(f) != 3 {
		return fmt.Errorf("tape: a MySQL error %q", text)
	}
	n, err := strconv.ParseUint(f[0], 10, 16)
	if err != nil {
		return fmt.Errorf("tape: a MySQL error %q", text)
	}
	me := &mysql.MySQLError{Number: uint16(n), Message: f[2]}
	if f[1] != "-" {
		copy(me.SQLState[:], f[1])
	}
	return me
}
