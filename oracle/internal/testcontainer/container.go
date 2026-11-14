// Package testcontainer is how oracle's tests reach Oracle: through a tape
// (enginetest/tape) that replays the pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. The compat
// judge, which runs only in a recording, asks the server directly (Judge).
//
// The server is one testcontainer per test binary, started by the first
// test that records on it and removed by Terminate, which the package's
// TestMain calls.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/sijms/go-ora/v2/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is oracle/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin row of the one server the tests ask.
const Default = "default"

// Password is every account's password: SYSTEM's, which the image takes
// from ORACLE_PASSWORD, and that of each user a test creates, which Open
// logs in with.
const Password = "omni_test"

// service is the pluggable database the tests use.
const service = "FREEPDB1"

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

// Open returns t's connection to Oracle as user, SYSTEM or a user the test
// created with Password: a replay of its tape, or the pinned server in a
// recording. Each call is a session of its own, and the pool holds that one
// session, so session state holds across statements. In a recording each
// top-level test starts on a server that holds nothing another test made:
// before its first session the users and public synonyms the image did not
// start with are dropped, and so is what SYSTEM owns that it did not.
func Open(t testing.TB, user string) *sql.DB {
	t.Helper()
	s := &tape.Server{Pins: Pins, Codec: Codec{}, Live: func(t testing.TB, pin enginetest.Pin) driver.Connector {
		t.Helper()
		return live(t, pin, user)
	}}
	return tape.Open(t, s, Default)
}

var server struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	host      string
	port      string
	err       error

	mu      sync.Mutex
	image   imageState      // what the server held when it started
	emptied map[string]bool // the top-level tests the server was emptied for
}

// imageState is what the pluggable database held when the server started,
// which emptying leaves: its users, its public synonyms, and SYSTEM's
// objects.
type imageState struct {
	users, synonyms, systemObjects map[string]bool
}

// live connects to the pinned server as user, emptying the server first on
// t's top-level test's first session.
func live(t testing.TB, pin enginetest.Pin, user string) driver.Connector {
	t.Helper()
	ctx := context.Background()
	if err := start(ctx, pin); err != nil {
		t.Fatalf("Oracle: %v", err)
	}
	top := enginetest.TopName(t.Name())
	server.mu.Lock()
	defer server.mu.Unlock()
	if !server.emptied[top] {
		if err := empty(ctx); err != nil {
			t.Fatalf("Oracle: empty the server for %s: %v", top, err)
		}
		server.emptied[top] = true
	}
	return connector(user, Password)
}

// connector connects to FREEPDB1 as user.
func connector(user, password string) driver.Connector {
	port, _ := strconv.Atoi(server.port)
	return go_ora.NewConnector(go_ora.BuildUrl(server.host, port, service, user, password, nil))
}

// start starts the pinned server once, checks it is the pinned image and
// release, and reads what it holds.
func start(ctx context.Context, pin enginetest.Pin) error {
	server.once.Do(func() {
		server.pin = pin
		server.emptied = map[string]bool{}
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
			Env:          map[string]string{"ORACLE_PASSWORD": Password},
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
	db := sql.OpenDB(connector("system", Password))
	defer db.Close()
	var version string
	if err := db.QueryRowContext(ctx, `SELECT version_full FROM product_component_version WHERE product LIKE 'Oracle%'`).Scan(&version); err != nil {
		return fmt.Errorf("read the Oracle version: %w", err)
	}
	if version != pin.Tag {
		return fmt.Errorf("the server reports %q, but %s pins %s", version, Pins, pin.Tag)
	}
	if server.image.users, err = names(ctx, db, "SELECT username FROM dba_users"); err != nil {
		return err
	}
	if server.image.synonyms, err = names(ctx, db, "SELECT synonym_name FROM dba_synonyms WHERE owner = 'PUBLIC'"); err != nil {
		return err
	}
	if server.image.systemObjects, err = names(ctx, db, systemObjects); err != nil {
		return err
	}
	return nil
}

// systemObjects lists the objects of SYSTEM's a test may make.
const systemObjects = "SELECT object_type || ' ' || object_name FROM user_objects WHERE object_type IN ('TABLE', 'VIEW', 'SEQUENCE', 'SYNONYM', 'PROCEDURE', 'FUNCTION', 'PACKAGE', 'TYPE')"

func names(ctx context.Context, db *sql.DB, query string) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out[n] = true
	}
	return out, rows.Err()
}

// empty drops what tests made since the server started: users with all
// they own, public synonyms, and SYSTEM's own objects.
func empty(ctx context.Context) error {
	db := sql.OpenDB(connector("system", Password))
	defer db.Close()
	db.SetMaxOpenConns(1)
	var stmts []string
	users, err := names(ctx, db, "SELECT username FROM dba_users")
	if err != nil {
		return err
	}
	for _, u := range sorted(users) {
		if !server.image.users[u] {
			stmts = append(stmts, `DROP USER "`+u+`" CASCADE`)
		}
	}
	synonyms, err := names(ctx, db, "SELECT synonym_name FROM dba_synonyms WHERE owner = 'PUBLIC'")
	if err != nil {
		return err
	}
	for _, s := range sorted(synonyms) {
		if !server.image.synonyms[s] {
			stmts = append(stmts, `DROP PUBLIC SYNONYM "`+s+`"`)
		}
	}
	objects, err := names(ctx, db, systemObjects)
	if err != nil {
		return err
	}
	for _, o := range sorted(objects) {
		if !server.image.systemObjects[o] {
			kind, name, _ := strings.Cut(o, " ")
			stmt := "DROP " + kind + ` "` + name + `"`
			if kind == "TABLE" {
				stmt += " CASCADE CONSTRAINTS PURGE"
			}
			stmts = append(stmts, stmt)
		}
	}
	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

func sorted(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
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
	admin := sql.OpenDB(connector("system", Password))
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

// Codec renders go-ora's *network.OracleError as "<code> <message>": the
// ORA- number, and the error's text without the position go-ora appends.
// It is also the Renamer of the names Oracle makes up, which differ on
// every run: a constraint's or an index's SYS_C, a materialized view's
// constraint's SYS_C_SNAP$_, a LOB segment's and its
// index's SYS_LOB and SYS_IL, an identity column's sequence ISEQ$$_.
type Codec struct{}

func (Codec) EncodeError(err error) (string, bool) {
	var oe *network.OracleError
	if !errors.As(err, &oe) {
		return "", false
	}
	return strconv.Itoa(oe.ErrCode) + " " + oe.ErrMsg, true
}

func (Codec) DecodeError(text string) error {
	code, msg, ok := strings.Cut(text, " ")
	n, err := strconv.Atoi(code)
	if !ok || err != nil {
		return fmt.Errorf("tape: an Oracle error %q", text)
	}
	oe := network.NewOracleError(n)
	oe.ErrMsg = msg
	return oe
}

var madeUp = regexp.MustCompile(`SYS_C\d{6,}|SYS_C_SNAP\$_\d+|SYS_(?:LOB|IL)\d{10}C\d{5}\$\$|ISEQ\$\$_\d+`)

// MadeUp matches a name Oracle made up.
func (Codec) MadeUp() *regexp.Regexp { return madeUp }

// StandIn is the nth made-up name's stand-in, of its own kind, numbered
// from 9 million, which Oracle reaches on no server a test starts.
func (Codec) StandIn(name string, n int) string {
	switch {
	case strings.HasPrefix(name, "SYS_C_SNAP$_"):
		return fmt.Sprintf("SYS_C_SNAP$_9%06d", n)
	case strings.HasPrefix(name, "SYS_C"):
		return fmt.Sprintf("SYS_C9%06d", n)
	case strings.HasPrefix(name, "ISEQ$$_"):
		return fmt.Sprintf("ISEQ$$_9%06d", n)
	}
	return fmt.Sprintf("%s9%09dC00001$$", strings.TrimRight(name[:7], "0123456789"), n)
}

// OpenLive returns a session as user to the pinned server outside any tape,
// for a test with nothing to replay: its input is Oracle's own data, which
// a tape may not hold. It skips t outside a recording (tape.Live).
func OpenLive(t testing.TB, user string) *sql.DB {
	t.Helper()
	tape.Live(t)
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := start(context.Background(), pin); err != nil {
		t.Fatalf("Oracle: %v", err)
	}
	db := sql.OpenDB(connector(user, Password))
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}
