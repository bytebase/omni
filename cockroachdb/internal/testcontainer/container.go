// Package testcontainer is how CockroachDB's tests reach a server: through a
// tape (enginetest/tape) that replays the pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. That server is
// one single-node testcontainer per pin role and test binary, started by the
// first test that records on it and removed by Terminate, which the package's
// TestMain calls.
package testcontainer

import (
	"bytes"
	"compress/zlib"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is cockroachdb/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin row the tapes replay.
const Default = "default"

var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}

// Open returns t's connection to CockroachDB: a replay of its tape, or the
// pinned server in a recording. Each call is a connection of its own, at
// session defaults, to the database of t's top-level test (Database). In a
// recording each top-level test starts on a database it alone holds: before
// its first connection the database is dropped and created empty.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return tape.Open(t, server, Default)
}

// Database is the name of the database Open connects t to: one per
// top-level test, named for it, so that a test reads the same
// current_database() in every recording.
func Database(t testing.TB) string {
	top := strings.ToLower(enginetest.TopName(t.Name()))
	name := "omni_" + notNameChars.ReplaceAllString(top, "_")
	if len(name) > 63 {
		sum := fmt.Sprintf("%x", top)
		name = name[:54] + "_" + sum[:8]
	}
	return name
}

var notNameChars = regexp.MustCompile(`[^a-z0-9_]`)

// statementTimeout fails a statement that hangs, and with it its test and
// not the run.
const statementTimeout = "30s"

// A instance is one pinned server this run started.
type instance struct {
	once      sync.Once
	termOnce  sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	base      *pgx.ConnConfig
	err       error

	mu      sync.Mutex
	created map[string]bool // the top-level tests whose databases were made
}

// terminate removes this test's server. The first test that starts the
// instance registers it as a cleanup, so a package run holds one server at a
// time instead of one per test.
func (in *instance) terminate() {
	in.termOnce.Do(func() {
		if in.container != nil {
			_ = testcontainers.TerminateContainer(in.container)
		}
	})
}

// One server per top-level test, not one per pin role: a table or view OID
// can collide with an OID in another catalog, and the sync's unqualified
// obj_description reads the other object's comment. Sharing a server would let
// what an earlier test created change those comments in a later test's tape.
var shared = struct {
	sync.Mutex
	byKey map[string]*instance
}{byKey: map[string]*instance{}}

func started(t testing.TB, pin enginetest.Pin) *instance {
	t.Helper()
	key := enginetest.TopName(t.Name()) + "\x00" + pin.Role
	shared.Lock()
	in := shared.byKey[key]
	if in == nil {
		in = &instance{pin: pin, created: map[string]bool{}}
		shared.byKey[key] = in
	}
	shared.Unlock()
	in.once.Do(func() { in.err = in.start(context.Background()) })
	if in.err != nil {
		t.Fatalf("start CockroachDB (%s): %v", pin.Role, in.err)
	}
	if in.pin != pin {
		t.Fatalf("this run started CockroachDB at %s, not %s", in.pin, pin)
	}
	t.Cleanup(in.terminate)
	return in
}

func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	in := started(t, pin)
	top := enginetest.TopName(t.Name())
	name := Database(t)
	in.mu.Lock()
	defer in.mu.Unlock()
	if !in.created[top] {
		db := in.admin()
		defer db.Close()
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quote(name) + " CASCADE",
			"CREATE DATABASE " + quote(name),
		} {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
		in.created[top] = true
	}
	return connector{Connector: stdlib.GetConnector(*in.config(name))}
}

// config returns the connection settings for one database. Arguments are
// inlined (the simple protocol), since a tape holds each statement with its
// arguments and no prepared statement.
func (in *instance) config(database string) *pgx.ConnConfig {
	cfg := in.base.Copy()
	cfg.Database = database
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	cfg.RuntimeParams["statement_timeout"] = statementTimeout
	return cfg
}

func (in *instance) admin() *sql.DB {
	return sql.OpenDB(stdlib.GetConnector(*in.config("defaultdb")))
}

func (in *instance) start(ctx context.Context) (err error) {
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        in.pin.ImageRef(),
			ExposedPorts: []string{"26257/tcp"},
			Cmd: []string{
				"start-single-node",
				"--insecure",
				"--listen-addr=:26257",
				"--http-addr=:8080",
				"--store=type=mem,size=25%",
			},
			WaitingFor: wait.ForLog("node startup completed").WithStartupTimeout(2 * time.Minute),
		},
		Started: true,
	})
	if c != nil {
		in.container = c
	}
	if err != nil {
		return err
	}
	if err := CheckDigest(ctx, c, in.pin.Digest); err != nil {
		return err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return err
	}
	port, err := c.MappedPort(ctx, "26257/tcp")
	if err != nil {
		return err
	}
	in.base, err = pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:%s/defaultdb?sslmode=disable", host, port.Port()))
	if err != nil {
		return err
	}
	db := in.admin()
	defer db.Close()
	var version string
	if err := db.QueryRowContext(ctx, "SHOW CLUSTER SETTING version").Scan(&version); err != nil {
		return fmt.Errorf("CockroachDB does not answer: %w", err)
	}
	// SHOW CLUSTER SETTING version answers the release line, 25.2 for
	// v25.2.3; the pin names the image's full tag.
	want := strings.TrimPrefix(in.pin.Tag, "v")
	if parts := strings.Split(want, "."); len(parts) > 2 {
		want = strings.Join(parts[:2], ".")
	}
	if version != want {
		return fmt.Errorf("the server reports %q, but pins.tsv pins CockroachDB %s", version, in.pin.Tag)
	}
	// Automatic statistics move a table's row estimate between recordings;
	// a test that needs an estimate collects statistics explicitly.
	if _, err := db.ExecContext(ctx, "SET CLUSTER SETTING sql.stats.automatic_collection.enabled = false"); err != nil {
		return fmt.Errorf("disable automatic statistics: %w", err)
	}
	return nil
}

// Terminate removes every server this run started.
func Terminate() {
	shared.Lock()
	defer shared.Unlock()
	for _, in := range shared.byKey {
		in.terminate()
	}
}

// CheckDigest proves the container runs the pinned image. The pull names
// image@digest, so a moved tag cannot select another image; this check reads
// what the daemon actually started, which is where the recorded answers came
// from.
func CheckDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the CockroachDB container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the CockroachDB image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the CockroachDB container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

func quote(name string) string { return `"` + strings.ReplaceAll(name, `"`, `""`) + `"` }

// connector is pgx's, with one change to what its rows report: see rows.
type connector struct {
	driver.Connector
}

func (c connector) Connect(ctx context.Context) (driver.Conn, error) {
	dc, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return conn{dc}, nil
}

type conn struct {
	driver.Conn
}

func (c conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

func (c conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	r, err := c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return rows{Rows: r}, nil
}

func (c conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
}

func (c conn) Ping(ctx context.Context) error { return c.Conn.(driver.Pinger).Ping(ctx) }

func (c conn) ResetSession(ctx context.Context) error {
	return c.Conn.(driver.SessionResetter).ResetSession(ctx)
}

// rows names a type pgx does not know by "-" rather than by its OID, which is
// the server's next free one when the type was made and so differs from run to
// run.
type rows struct {
	driver.Rows
}

func (r rows) ColumnTypeDatabaseTypeName(i int) string {
	name := r.Rows.(driver.RowsColumnTypeDatabaseTypeName).ColumnTypeDatabaseTypeName(i)
	if _, err := strconv.ParseUint(name, 10, 32); err == nil {
		return "-"
	}
	return name
}

func (r rows) HasNextResultSet() bool {
	if n, ok := r.Rows.(driver.RowsNextResultSet); ok {
		return n.HasNextResultSet()
	}
	return false
}

func (r rows) NextResultSet() error {
	if n, ok := r.Rows.(driver.RowsNextResultSet); ok {
		return n.NextResultSet()
	}
	return errors.New("no next result set")
}

// planLink is the one made-up value CockroachDB puts in an answer: the page a
// DISTSQL plan links to carries a random flow_id inside its deflated payload.
// A recording replaces that id with a numbered stand-in, so the same plan
// replays byte for byte.
var planLink = regexp.MustCompile(`https?://\S+#[A-Za-z0-9+/_=-]+`)

var flowID = regexp.MustCompile(`"flow_id":"[^"]*"`)

// MadeUp matches the plan link.
func (Codec) MadeUp() *regexp.Regexp { return planLink }

// StandIn replaces the plan's flow_id with a numbered one and re-encodes the
// payload, so the link is the same in every recording.
func (Codec) StandIn(name string, n int) string {
	index := strings.LastIndex(name, "#")
	if index < 0 {
		return name
	}
	payload := name[index+1:]
	normalized := strings.NewReplacer("-", "+", "_", "/").Replace(strings.TrimRight(payload, "="))
	deflated, err := base64.RawStdEncoding.DecodeString(normalized)
	if err != nil {
		return name
	}
	reader, err := zlib.NewReader(bytes.NewReader(deflated))
	if err != nil {
		return name
	}
	plan, err := io.ReadAll(reader)
	_ = reader.Close()
	if err != nil {
		return name
	}
	standIn := fmt.Sprintf(`"flow_id":"00000000-0000-0000-0000-%012d"`, n)
	plan = flowID.ReplaceAllLiteral(plan, []byte(standIn))
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)
	if _, err := writer.Write(plan); err != nil {
		return name
	}
	if err := writer.Close(); err != nil {
		return name
	}
	return name[:index+1] + base64.RawURLEncoding.EncodeToString(compressed.Bytes())
}

// Codec renders pgx's *pgconn.PgError as its SQLSTATE and message on one
// line, then a line for each other field the server sent:
//
//	42601 syntax error at or near "FROM"
//	position: 8
//
// A line break inside a field is written \n, and a backslash \\. The source
// file, line and routine the server names are left out: they say where in
// CockroachDB's code the error was raised and move with a patch release.
type Codec struct{}

var errorFields = []struct {
	name string
	get  func(*pgconn.PgError) *string
}{
	{"severity", func(e *pgconn.PgError) *string { return &e.Severity }},
	{"detail", func(e *pgconn.PgError) *string { return &e.Detail }},
	{"hint", func(e *pgconn.PgError) *string { return &e.Hint }},
	{"internal query", func(e *pgconn.PgError) *string { return &e.InternalQuery }},
	{"where", func(e *pgconn.PgError) *string { return &e.Where }},
	{"schema", func(e *pgconn.PgError) *string { return &e.SchemaName }},
	{"table", func(e *pgconn.PgError) *string { return &e.TableName }},
	{"column", func(e *pgconn.PgError) *string { return &e.ColumnName }},
	{"data type", func(e *pgconn.PgError) *string { return &e.DataTypeName }},
	{"constraint", func(e *pgconn.PgError) *string { return &e.ConstraintName }},
}

const plainSeverity = "ERROR"

var (
	escaper   = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "\r", `\r`)
	unescaper = strings.NewReplacer(`\\`, `\`, `\n`, "\n", `\r`, "\r")
)

func (Codec) EncodeError(err error) (string, bool) {
	var pe *pgconn.PgError
	if !errors.As(err, &pe) {
		return "", false
	}
	var b strings.Builder
	b.WriteString(pe.Code + " " + escaper.Replace(pe.Message))
	if pe.Position != 0 {
		fmt.Fprintf(&b, "\nposition: %d", pe.Position)
	}
	if pe.InternalPosition != 0 {
		fmt.Fprintf(&b, "\ninternal position: %d", pe.InternalPosition)
	}
	for _, f := range errorFields {
		v := *f.get(pe)
		if v == "" || f.name == "severity" && v == plainSeverity {
			continue
		}
		b.WriteString("\n" + f.name + ": " + escaper.Replace(v))
	}
	return b.String(), true
}

func (Codec) DecodeError(text string) error {
	lines := strings.Split(text, "\n")
	code, message, ok := strings.Cut(lines[0], " ")
	if !ok || len(code) != 5 {
		return fmt.Errorf("tape: a CockroachDB error %q", text)
	}
	pe := &pgconn.PgError{Severity: plainSeverity, Code: code, Message: unescaper.Replace(message)}
fields:
	for _, line := range lines[1:] {
		name, v, ok := strings.Cut(line, ": ")
		if !ok {
			return fmt.Errorf("tape: a CockroachDB error field %q", line)
		}
		switch name {
		case "position", "internal position":
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return fmt.Errorf("tape: a CockroachDB error field %q", line)
			}
			if name == "position" {
				pe.Position = int32(n)
			} else {
				pe.InternalPosition = int32(n)
			}
			continue
		}
		for _, f := range errorFields {
			if f.name == name {
				*f.get(pe) = unescaper.Replace(v)
				continue fields
			}
		}
		return fmt.Errorf("tape: a CockroachDB error field %q", line)
	}
	pe.SeverityUnlocalized = pe.Severity
	return pe
}
