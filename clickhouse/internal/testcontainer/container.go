// Package testcontainer is how clickhouse's driver tests reach ClickHouse:
// through a tape (enginetest/tape) that replays the pinned server's answers,
// or, in a recording (OMNI_RECORD=1), through the pinned server itself. The
// server is one testcontainer, started by the first test that records on it
// and removed by Terminate, which the package's TestMain calls.
//
// ClickHouse's Go client is a database/sql driver, so a test gets an ordinary
// *sql.DB and the tape holds the statements and rows, not a private protocol.
//
// The image disables network access for the default account unless something
// asks it not to, so the container is started with
// CLICKHOUSE_SKIP_USER_SETUP=1: the account the tests use is the image's own
// "default", with no password.
package testcontainer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is clickhouse/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin row of the one server the tests ask.
const Default = "default"

// User is the image's own account, with no password.
const User = "default"

// Port is the native protocol's port, which the client speaks.
const Port = "9000"

var server struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	host      string
	port      string
	err       error

	mu      sync.Mutex
	emptied map[string]bool // the top-level tests the server was emptied for
}

// Open returns t's connection to ClickHouse on the database named for t's
// top-level test: a replay of its tape, or the pinned server in a recording.
// In a recording each top-level test starts on a server that holds nothing
// another test made: before its first connection every database but the
// image's own is dropped and the test's database is created empty.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	s := &tape.Server{Pins: Pins, Codec: Codec{}, Types: Types, Live: func(t testing.TB, pin enginetest.Pin) driver.Connector {
		t.Helper()
		return live(t, pin)
	}}
	return tape.Open(t, s, Default)
}

// Judge returns a live connection to the pinned ClickHouse server, outside any
// tape, for syntax compatibility judging (OMNI_RECORD=1 or OMNI_RECORD=verify).
// Callers must invoke Terminate from TestMain so the container is removed.
func Judge(t testing.TB, pin enginetest.Pin) *sql.DB {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording or verification: a replay reaches no server")
	}
	if err := start(context.Background(), pin); err != nil {
		t.Fatalf("ClickHouse: %v", err)
	}
	db := sql.OpenDB(connector("default"))
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// OpenMissingDatabase returns t's connection to a database the pinned server
// does not have. The client's handshake is refused, which is all ClickHouse
// answers for a database that is not there, so a tape records the refused
// connection. name must be one no test creates; the emptying below drops every
// database the image did not start with and creates none.
func OpenMissingDatabase(t testing.TB, name string) *sql.DB {
	t.Helper()
	s := &tape.Server{Pins: Pins, Codec: Codec{}, Types: Types, Live: func(t testing.TB, pin enginetest.Pin) driver.Connector {
		t.Helper()
		if err := start(context.Background(), pin); err != nil {
			t.Fatalf("ClickHouse: %v", err)
		}
		top := enginetest.TopName(t.Name())
		server.mu.Lock()
		defer server.mu.Unlock()
		if !server.emptied[top] {
			// empty keeps the database it is given and creates it if it is not
			// there; the image's own is in systemDatabases and already there.
			if err := empty(context.Background(), "default"); err != nil {
				t.Fatalf("ClickHouse: empty the server for %s: %v", top, err)
			}
			server.emptied[top] = true
		}
		return connector(name)
	}}
	return tape.Open(t, s, Default)
}

// Database is the database Open selects for t: one per top-level test, named
// for it, so a test reads the same name in every recording.
func Database(t testing.TB) string {
	top := strings.ToLower(enginetest.TopName(t.Name()))
	name := "omni_" + notNameChars.ReplaceAllString(top, "_")
	if len(name) > 63 {
		name = name[:63]
	}
	return name
}

var notNameChars = regexp.MustCompile(`[^a-z0-9_]`)

// live connects to t's database on the pinned server, emptying the server
// first on t's top-level test's first connection.
func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	if err := start(context.Background(), pin); err != nil {
		t.Fatalf("ClickHouse: %v", err)
	}
	database := Database(t)
	top := enginetest.TopName(t.Name())
	server.mu.Lock()
	defer server.mu.Unlock()
	if !server.emptied[top] {
		if err := empty(context.Background(), database); err != nil {
			t.Fatalf("ClickHouse: empty the server for %s: %v", top, err)
		}
		server.emptied[top] = true
	}
	return connector(database)
}

// connector connects to database on the pinned server, with the settings the
// driver opens with.
func connector(database string) driver.Connector {
	return clickhouse.Connector(&clickhouse.Options{
		Addr: []string{server.host + ":" + server.port},
		Auth: clickhouse.Auth{
			Database: database,
			Username: User,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 300,
		},
		DialTimeout: 10 * time.Second,
	})
}

// start starts the pinned server once and checks it is the pinned image.
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

func boot(ctx context.Context, pin enginetest.Pin) (err error) {
	// testcontainers panics when Docker is missing; fail the test instead of
	// the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.ImageRef(),
			ExposedPorts: []string{Port + "/tcp"},
			Env:          map[string]string{"CLICKHOUSE_SKIP_USER_SETUP": "1"},
			// The server's stdout carries the entrypoint's lines only; the
			// server's own "Ready for connections" goes to its log file. The
			// listening port is up a moment before a query can run, so boot
			// waits for an answer below.
			WaitingFor: wait.ForListeningPort("9000/tcp").WithStartupTimeout(2 * time.Minute),
		},
		Started: true,
	})
	if c != nil {
		server.container = c
	}
	if err != nil {
		return fmt.Errorf("start ClickHouse: %w", err)
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		return err
	}
	if server.host, err = c.Host(ctx); err != nil {
		return fmt.Errorf("ClickHouse's host: %w", err)
	}
	mapped, err := c.MappedPort(ctx, Port)
	if err != nil {
		return fmt.Errorf("ClickHouse's port: %w", err)
	}
	server.port = mapped.Port()
	if err := waitForAnswer(ctx, pin); err != nil {
		return err
	}
	return nil
}

// waitForAnswer asks the server a question until it answers with the pinned
// release, which is what "ready" means to the client.
func waitForAnswer(ctx context.Context, pin enginetest.Pin) error {
	deadline := time.Now().Add(2 * time.Minute)
	var last error
	for time.Now().Before(deadline) {
		db := sql.OpenDB(connector("default"))
		var version string
		last = db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)
		_ = db.Close()
		if last == nil {
			if version != pin.Tag {
				return fmt.Errorf("the server reports %q, but %s pins %s", version, Pins, pin.Tag)
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
	return fmt.Errorf("ClickHouse did not answer: %w", last)
}

// empty drops every database but the image's own, creates database empty,
// and drops every account, role, row policy, quota and settings profile but
// the image's own.
func empty(ctx context.Context, database string) error {
	db := sql.OpenDB(connector("default"))
	defer db.Close()
	db.SetMaxOpenConns(1)

	var stmts []string
	databases, err := names(ctx, db, "SELECT name FROM system.databases")
	if err != nil {
		return err
	}
	for _, name := range databases {
		if !systemDatabases[name] && name != database {
			stmts = append(stmts, "DROP DATABASE IF EXISTS `"+escape(name)+"`")
		}
	}
	stmts = append(stmts, "CREATE DATABASE IF NOT EXISTS `"+escape(database)+"`")
	users, err := names(ctx, db, "SELECT name FROM system.users")
	if err != nil {
		return err
	}
	for _, name := range users {
		if name != User {
			stmts = append(stmts, "DROP USER IF EXISTS `"+escape(name)+"`")
		}
	}
	roles, err := names(ctx, db, "SELECT name FROM system.roles")
	if err != nil {
		return err
	}
	for _, name := range roles {
		stmts = append(stmts, "DROP ROLE IF EXISTS `"+escape(name)+"`")
	}
	// A row policy outlives the table and the database it names. The image's
	// own quota and settings profiles live in users.xml, which DROP cannot
	// change; every other one is a test's.
	policies, err := rows(ctx, db, "SELECT short_name, database, table FROM system.row_policies WHERE storage != 'users_xml'")
	if err != nil {
		return err
	}
	for _, p := range policies {
		table := "*"
		if p[2] != "" {
			table = "`" + escape(p[2]) + "`"
		}
		stmts = append(stmts, "DROP ROW POLICY IF EXISTS `"+escape(p[0])+"` ON `"+escape(p[1])+"`."+table)
	}
	for _, kind := range []struct{ table, drop string }{
		{"system.quotas", "DROP QUOTA IF EXISTS "},
		{"system.settings_profiles", "DROP SETTINGS PROFILE IF EXISTS "},
	} {
		entities, err := names(ctx, db, "SELECT name FROM "+kind.table+" WHERE storage != 'users_xml'")
		if err != nil {
			return err
		}
		for _, name := range entities {
			stmts = append(stmts, kind.drop+"`"+escape(name)+"`")
		}
	}
	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

// systemDatabases are the databases the image starts with.
var systemDatabases = map[string]bool{
	"default":            true,
	"system":             true,
	"information_schema": true,
	"INFORMATION_SCHEMA": true,
}

func names(ctx context.Context, db *sql.DB, query string) ([]string, error) {
	found, err := rows(ctx, db, query)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(found))
	for _, row := range found {
		out = append(out, row[0])
	}
	return out, nil
}

// rows answers a query whose columns are all strings.
func rows(ctx context.Context, db *sql.DB, query string) ([][]string, error) {
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer result.Close()
	columns, err := result.Columns()
	if err != nil {
		return nil, err
	}
	var out [][]string
	for result.Next() {
		row := make([]string, len(columns))
		dest := make([]any, len(columns))
		for i := range row {
			dest[i] = &row[i]
		}
		if err := result.Scan(dest...); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, result.Err()
}

func escape(name string) string {
	return strings.ReplaceAll(name, "`", "``")
}

// checkDigest proves the container runs the pinned image. The pull names
// image@digest, so a moved tag cannot select another image; this check reads
// what the daemon actually started, which is where the recorded answers came
// from.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the ClickHouse container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the ClickHouse image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the ClickHouse container runs an image with digests %v, but %s pins %s", img.RepoDigests, Pins, digest)
}

// Terminate removes the server this process started.
func Terminate() {
	if server.container != nil {
		_ = testcontainers.TerminateContainer(server.container)
	}
}

// Codec renders clickhouse-go's *clickhouse.Exception as the fields that make
// up its text, so a replay's error has the recording's type and message.
type Codec struct{}

// exception is one server exception as a tape holds it. StackTrace stays out:
// it names the server's own source paths and says nothing a test reads.
type exception struct {
	Code    int32  `json:"code"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message"`
}

func (Codec) EncodeError(err error) (string, bool) {
	var e *clickhouse.Exception
	if !errors.As(err, &e) {
		return "", false
	}
	b, marshalErr := json.Marshal(exception{Code: e.Code, Name: e.Name, Message: e.Message})
	if marshalErr != nil {
		return "", false
	}
	return string(b), true
}

func (Codec) DecodeError(text string) error {
	var e exception
	if err := json.Unmarshal([]byte(text), &e); err != nil {
		return fmt.Errorf("tape: a ClickHouse error %q: %w", text, err)
	}
	return &clickhouse.Exception{Code: e.Code, Name: e.Name, Message: e.Message}
}
