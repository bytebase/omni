package sharedserver

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// OceanBaseImage is the build the recorded OceanBase verdicts were read on. The
// 4.3.5-lts tag moves with each patch release, so this names the patch.
const OceanBaseImage = "oceanbase/oceanbase-ce:4.3.5.6-106000012026040916"

// OceanBaseSocket is the observer's unix socket in a SLIM container.
const OceanBaseSocket = "/root/demo/run/sql.sock"

// oceanBaseSetup runs inside a SLIM container once it boots; CI and dev-db.sh run
// the same two commands.
//
// SLIM caps the data file at 1 GiB. With parser, catalog and review run in a row
// on one server, the cap makes the test tenant's flushes fail with
// OB_SERVER_OUTOF_DISK_SPACE, and on a busy box its root service shut down partway
// through the review walk; without it the file grew to 4.2 GiB and all three
// passed (2026-11-05). External tables read files only under secure_file_priv,
// which the server lets a client set only over its unix socket.
var oceanBaseSetup = [][]string{
	{"obclient", "-h127.1", "-P2881", "-uroot@sys", "-e", "ALTER SYSTEM SET datafile_maxsize = '8G'"},
	{"obclient", "-S", OceanBaseSocket, "-uroot@test", "-e", "SET GLOBAL secure_file_priv = '/'"},
}

// StartOceanBase boots image in SLIM mode and returns the container with a DSN
// for root in the test tenant, which has no password.
//
// SLIM starts the observer from the store the image ships, already holding the
// test tenant: about 30 s and 3.5 GiB under the oracle tests, where the default
// MINI mode deploys a cluster and creates the tenant, which takes minutes and
// 8 GiB. The deployer refuses fewer than 20000 open files.
func StartOceanBase(ctx context.Context, image string) (testcontainers.Container, string, error) {
	c, err := testcontainers.Run(ctx, image,
		testcontainers.WithEnv(map[string]string{"MODE": "SLIM"}),
		testcontainers.WithExposedPorts("2881/tcp"),
		testcontainers.WithHostConfigModifier(func(hc *container.HostConfig) {
			hc.Ulimits = []*units.Ulimit{{Name: "nofile", Soft: 65536, Hard: 65536}}
		}),
		testcontainers.WithWaitStrategy(wait.ForLog("boot success!").WithStartupTimeout(10*time.Minute)),
	)
	if err != nil {
		return c, "", fmt.Errorf("start %s: %w", image, err)
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

// PrivateServer is a server one test process has to itself, for tests that need
// what a shared server withholds: root, fixed database names, server-wide
// settings. It is the server Env names, when set, and otherwise a testcontainer
// started on first use and kept for the rest of the process.
//
// A server Env names belongs to one process at a time. CI runs the packages that
// use it one after another, as it does the trino packages.
type PrivateServer struct {
	Env   string
	Image string
	start func(ctx context.Context, image string) (testcontainers.Container, string, error)

	once      sync.Once
	dsn       string
	container testcontainers.Container
	err       error
}

// PrivateOceanBase is OceanBase for the tests that run as root: the oracles in
// oceanbase/parser, oceanbase/catalog and oceanbase/review.
var PrivateOceanBase = &PrivateServer{
	Env:   "OCEANBASE_WALK_DSN",
	Image: OceanBaseImage,
	start: StartOceanBase,
}

// DSN returns the server's admin DSN, naming no database, once it answers.
func (p *PrivateServer) DSN(ctx context.Context) (string, error) {
	p.once.Do(func() {
		if p.dsn = os.Getenv(p.Env); p.dsn == "" {
			p.container, p.dsn, p.err = p.start(ctx, p.Image)
			if p.err != nil {
				return
			}
		}
		if err := waitReady(ctx, p.dsn); err != nil {
			p.err = fmt.Errorf("%s: %w", p.Image, err)
		}
	})
	return p.dsn, p.err
}

// Terminate stops a testcontainer this process started. It leaves a server Env
// names alone.
func (p *PrivateServer) Terminate() {
	if p.container != nil {
		_ = testcontainers.TerminateContainer(p.container)
	}
}
