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
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
)

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
	err       error
}

func start(ctx context.Context, pin enginetest.Pin) error {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: pin.Image,
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
	// The judge sends a probe and the statement as one request.
	cfg.MultiStatements = true
	shared.config = cfg
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
