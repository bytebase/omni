// Package testcontainer is how oracle's tests reach Oracle. The compat
// judge, which runs only in a recording (OMNI_RECORD=1), asks the pinned
// server directly (Judge).
//
// The server is one testcontainer per test binary, started by the first
// test that records on it and removed by Terminate, which the package's
// TestMain calls.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
)

// Pins is oracle/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// The image's accounts: SYSTEM, whose password the image takes from
// ORACLE_PASSWORD, and the pluggable database FREEPDB1 the tests use.
const (
	adminPassword = "omni_test"
	service       = "FREEPDB1"
)

// The judge's account, which owns the scratch schema compat.OracleParse
// runs DDL in. It may create the objects a statement row makes, so DDL runs
// and PL/SQL compiles; it holds no system privilege beyond those, so a
// statement on the server or on another user answers ORA-01031, after the
// parse.
const (
	judgeUser   = "OMNI_JUDGE"
	judgeSecret = "omni_judge"
)

// JudgeFixture is the scratch schema's one object, built again after every
// statement that ran: the table t, on which a trigger row's header stands
// (tokens.tsv's TRIGGER_SOURCE), so the trigger's body reaches the PL/SQL
// compiler. It is not named a, the name every other row gives, so the rows
// on a table a find none, and answer at once.
const JudgeFixture = "CREATE TABLE t (a NUMBER)"

var server struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	host      string
	port      string
	err       error
}

// start starts the pinned server once, and checks it is the pinned image
// and release.
func start(ctx context.Context, pin enginetest.Pin) error {
	server.once.Do(func() {
		server.pin = pin
		server.err = boot(ctx, pin)
	})
	if server.err != nil {
		return server.err
	}
	if server.pin != pin {
		return fmt.Errorf("this run started %s, not %s", server.pin, pin)
	}
	return nil
}

func boot(ctx context.Context, pin enginetest.Pin) error {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{"1521/tcp"},
			Env:          map[string]string{"ORACLE_PASSWORD": adminPassword},
			WaitingFor:   wait.ForLog("DATABASE IS READY TO USE!").WithStartupTimeout(10 * time.Minute),
		},
		Started: true,
	})
	if c != nil {
		server.container = c
	}
	if err != nil {
		return fmt.Errorf("start %s: %w", pin.Image, err)
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		return err
	}
	if server.host, err = c.Host(ctx); err != nil {
		return err
	}
	port, err := c.MappedPort(ctx, "1521")
	if err != nil {
		return err
	}
	server.port = port.Port()
	db := sql.OpenDB(connector("system", adminPassword))
	defer db.Close()
	var version string
	if err := db.QueryRowContext(ctx, `SELECT version_full FROM product_component_version WHERE product LIKE 'Oracle%'`).Scan(&version); err != nil {
		return fmt.Errorf("read the Oracle version: %w", err)
	}
	if version != pin.Tag {
		return fmt.Errorf("the server reports %q, but %s pins %s", version, Pins, pin.Tag)
	}
	return nil
}

// connector connects to FREEPDB1 as user.
func connector(user, password string) driver.Connector {
	port, _ := strconv.Atoi(server.port)
	return go_ora.NewConnector(go_ora.BuildUrl(server.host, port, service, user, password, nil))
}

// checkDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Oracle container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Oracle image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Oracle container runs an image with digests %v, but %s pins %s", img.RepoDigests, Pins, digest)
}

// Judge makes the compat judge's account, whose answers go to compat.tsv,
// with a scratch schema holding JudgeFixture, and returns what opens a
// session as that account: the session, and what ends it however it stands
// (Cuttable, and the server process killed). It runs in a recording only,
// and starts the server.
func Judge(t testing.TB, pin enginetest.Pin) func() (*sql.DB, func()) {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer: the judge runs in a recording only")
	}
	ctx := context.Background()
	if err := start(ctx, pin); err != nil {
		t.Fatalf("Oracle: %v", err)
	}
	admin := sql.OpenDB(connector("system", adminPassword))
	t.Cleanup(func() { _ = admin.Close() })
	var users int
	if err := admin.QueryRowContext(ctx, "SELECT COUNT(*) FROM dba_users WHERE username = :1", judgeUser).Scan(&users); err != nil {
		t.Fatalf("Oracle: %v", err)
	}
	var stmts []string
	if users > 0 {
		stmts = append(stmts, "DROP USER "+judgeUser+" CASCADE")
	}
	stmts = append(stmts,
		fmt.Sprintf("CREATE USER %s IDENTIFIED BY %s QUOTA UNLIMITED ON users", judgeUser, judgeSecret),
		"GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE MATERIALIZED VIEW, "+
			"CREATE PROCEDURE, CREATE TYPE, CREATE TRIGGER TO "+judgeUser)
	for _, stmt := range stmts {
		if _, err := admin.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("Oracle: %s: %v", stmt, err)
		}
	}
	port, _ := strconv.Atoi(server.port)
	url := go_ora.BuildUrl(server.host, port, service, judgeUser, judgeSecret, nil)
	open := func() (*sql.DB, func()) {
		db, cut := Cuttable(url)
		var sid, audsid string
		if err := db.QueryRowContext(ctx, "SELECT SYS_CONTEXT('USERENV', 'SID'), SYS_CONTEXT('USERENV', 'SESSIONID') FROM dual").Scan(&sid, &audsid); err != nil {
			t.Fatalf("Oracle: the judge's session: %v", err)
		}
		end := func() {
			cut()
			// A parse the server never finishes keeps its process, and the
			// locks it holds, after its session is killed: kill the process.
			var spid string
			if err := admin.QueryRowContext(ctx, `SELECT p.spid FROM v$session s JOIN v$process p ON p.addr = s.paddr WHERE s.sid = :1 AND s.audsid = :2`, sid, audsid).Scan(&spid); err != nil {
				return // gone already
			}
			if code, _, err := server.container.Exec(ctx, []string{"bash", "-c", "kill -9 " + spid}); err != nil || code != 0 {
				t.Fatalf("Oracle: kill the judge's server process %s: exit %d, %v", spid, code, err)
			}
		}
		return db, end
	}
	db, _ := open()
	defer db.Close()
	if _, err := db.ExecContext(ctx, JudgeFixture); err != nil {
		t.Fatalf("Oracle: %s: %v", JudgeFixture, err)
	}
	return open
}

// Cuttable opens one session to the go-ora URL url, and returns it with
// what cuts its network connection. A statement the server never answers
// leaves go-ora reading for good and the server's session holding its
// locks; cutting the connection ends both, where closing the session would
// wait on the read.
func Cuttable(url string) (*sql.DB, func()) {
	d := &cutDialer{}
	c := go_ora.NewConnector(url).(*go_ora.OracleConnector)
	c.Dialer(d)
	db := sql.OpenDB(c)
	db.SetMaxOpenConns(1)
	return db, d.cut
}

// cutDialer dials go-ora's connections and keeps them, so cut can close
// them.
type cutDialer struct {
	mu    sync.Mutex
	conns []net.Conn
}

func (d *cutDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, err := (&net.Dialer{}).DialContext(ctx, network, address)
	if err == nil {
		d.mu.Lock()
		d.conns = append(d.conns, c)
		d.mu.Unlock()
	}
	return c, err
}

func (d *cutDialer) cut() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.conns {
		_ = c.Close()
	}
	d.conns = nil
}

// Terminate removes the server this process started.
func Terminate() {
	if server.container != nil {
		_ = testcontainers.TerminateContainer(server.container)
	}
}
