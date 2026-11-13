// Package testcontainer pins the MariaDB server the tests of this engine run
// against, and starts it. Each package starts its own container on first use,
// so the pin is what keeps them on one version: a test that records a verdict
// records it from the server in testdata/pins.tsv, and an upgrade there means
// rerunning every oracle test and reading the drift.
package testcontainer

import (
	"context"
	"database/sql"
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

// Image is the pinned MariaDB, as pins.tsv names it. mariadb/parser,
// mariadb/catalog, mariadb/review and mariadb/validate all start it.
const Image = "mariadb:11.8.8"

// Pins is mariadb/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// rootPassword is the root password of the server Pinned starts.
const rootPassword = "test"

var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	config    *mysql.Config
	err       error
}

// Pinned returns a connection pool to the server pin names, as root with the
// database test selected, starting the server on the first call. It fails t
// unless the container runs the pinned digest and reports the pinned release.
// Nothing it is asked is recorded: it is for the compat judge, which runs
// only in a recording.
func Pinned(t testing.TB, pin enginetest.Pin) *sql.DB {
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
	c, err := mysql.NewConnector(shared.config.Clone())
	if err != nil {
		t.Fatalf("%v", err)
	}
	db := sql.OpenDB(c)
	t.Cleanup(func() { _ = db.Close() })
	return db
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
	cfg.DBName = "test"
	// Arguments are inlined, so a statement and its arguments are one request.
	cfg.InterpolateParams = true
	// A statement that hangs fails its test, not the run.
	cfg.ReadTimeout = 30 * time.Second
	cfg.WriteTimeout = 30 * time.Second
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
