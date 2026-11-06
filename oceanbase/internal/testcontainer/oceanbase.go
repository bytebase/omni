package testcontainer

import (
	"context"
	"fmt"
	"io"
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
