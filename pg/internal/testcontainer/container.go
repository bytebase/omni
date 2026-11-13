// Package testcontainer is how pg's tests reach PostgreSQL: through a tape
// (enginetest/tape) that replays the pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. That server is
// one testcontainer per pin role and test binary, started by the first test
// that records on it and removed by Terminate, which the package's TestMain
// calls.
package testcontainer

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is pg/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// The pin roles. Default is the release the parser's grammar is read from,
// and judges compat.tsv; PG18 serves the tests that ask the next major.
const (
	Default = "default"
	PG18    = "pg18"
)

var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}

// Open returns t's connection to PostgreSQL: a replay of its tape, or the
// pinned server in a recording. Each call is a connection of its own, at
// session defaults, to the database of t's top-level test: a recording
// creates that database empty before the test's first connection, so a
// tape never holds what another test left.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return tape.Open(t, server, Default)
}

// OpenRole is Open on the server of another pin row, such as PG18.
func OpenRole(t testing.TB, role string) *sql.DB {
	t.Helper()
	return tape.Open(t, server, role)
}

// Reset runs fn on the default server before a test opens its connection,
// in a recording only: it puts what databases share (roles, tablespaces)
// into the state the test's tape starts from, so a replay needs nothing.
// fn's connection is to the database postgres, and is not recorded.
func Reset(t testing.TB, fn func(ctx context.Context, db *sql.DB) error) {
	t.Helper()
	ResetRole(t, Default, fn)
}

// ResetRole is Reset on the server of another pin row.
func ResetRole(t testing.TB, role string, fn func(ctx context.Context, db *sql.DB) error) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	db := admin(t, role)
	defer db.Close()
	if err := fn(context.Background(), db); err != nil {
		t.Fatalf("reset PostgreSQL: %v", err)
	}
}

// Judge returns a connection to the default server itself, outside any
// tape, for the judge of compat.tsv, which runs only in a recording. A
// statement sent with no arguments reaches the server as written, in one
// simple query.
func Judge(t testing.TB) *sql.DB {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	db := admin(t, Default)
	t.Cleanup(func() { db.Close() })
	return db
}

// admin connects to the database postgres of the server in pin row role,
// starting the server on first use.
func admin(t testing.TB, role string) *sql.DB {
	t.Helper()
	pin, err := enginetest.LoadPin(Pins, role)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return sql.OpenDB(stdlib.GetConnector(*started(t, pin).config("postgres")))
}

// A instance is one pinned server this run started.
type instance struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	base      *pgx.ConnConfig
	err       error

	mu      sync.Mutex
	created map[string]bool // the test databases this run made
}

var shared = struct {
	sync.Mutex
	byRole map[string]*instance
}{byRole: map[string]*instance{}}

// started returns the server of pin, starting it on the first call.
func started(t testing.TB, pin enginetest.Pin) *instance {
	t.Helper()
	shared.Lock()
	in := shared.byRole[pin.Role]
	if in == nil {
		in = &instance{pin: pin, created: map[string]bool{}}
		shared.byRole[pin.Role] = in
	}
	shared.Unlock()
	in.once.Do(func() { in.err = in.start(context.Background()) })
	if in.err != nil {
		t.Fatalf("start PostgreSQL (%s): %v", pin.Role, in.err)
	}
	if in.pin != pin {
		t.Fatalf("this run started PostgreSQL at %s, not %s", in.pin, pin)
	}
	return in
}

// live connects to the database of t's top-level test on the pinned
// server, creating it empty on the test's first connection.
func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	in := started(t, pin)
	name := Database(t)
	in.mu.Lock()
	defer in.mu.Unlock()
	if !in.created[name] {
		db := sql.OpenDB(stdlib.GetConnector(*in.config("postgres")))
		defer db.Close()
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + name + " WITH (FORCE)",
			"CREATE DATABASE " + name,
		} {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
		in.created[name] = true
	}
	return connector{stdlib.GetConnector(*in.config(name))}
}

// Database is the name of the database Open connects t to: one per
// top-level test, named for it, so that a test reads the same
// current_database() in every recording. PostgreSQL keeps 63 bytes of a
// name, so a long test name ends in a hash of the whole.
func Database(t testing.TB) string {
	top := strings.ToLower(enginetest.TopName(t.Name()))
	name := "omni_" + notNameChars.ReplaceAllString(top, "_")
	if len(name) > 63 {
		sum := sha256.Sum256([]byte(top))
		name = name[:54] + "_" + hex.EncodeToString(sum[:4])
	}
	return name
}

var notNameChars = regexp.MustCompile(`[^a-z0-9_]`)

// statementTimeout fails a statement that hangs, and with it its test and
// not the run.
const statementTimeout = "30s"

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

func (in *instance) start(ctx context.Context) (err error) {
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := tcpg.Run(ctx, in.pin.Image,
		tcpg.WithDatabase("postgres"),
		tcpg.WithUsername("postgres"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(2*time.Minute)),
	)
	if c != nil {
		in.container = c
	}
	if err != nil {
		return err
	}
	if err := CheckDigest(ctx, c, in.pin.Digest); err != nil {
		return err
	}
	dsn, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}
	if in.base, err = pgx.ParseConfig(dsn); err != nil {
		return err
	}
	db := sql.OpenDB(stdlib.GetConnector(*in.config("postgres")))
	defer db.Close()
	var version string
	if err := db.QueryRowContext(ctx, "SHOW server_version").Scan(&version); err != nil {
		return fmt.Errorf("PostgreSQL does not answer: %w", err)
	}
	if want := Version(in.pin.Tag); version != want {
		return fmt.Errorf("the server reports %q, but pins.tsv pins PostgreSQL %s (%s)", version, want, in.pin.Tag)
	}
	return nil
}

// Version is the server_version of a source tag: 17.12 for REL_17_12.
func Version(tag string) string {
	return strings.ReplaceAll(strings.TrimPrefix(tag, "REL_"), "_", ".")
}

// Terminate removes every server this run started.
func Terminate() {
	shared.Lock()
	defer shared.Unlock()
	for _, in := range shared.byRole {
		if in.container != nil {
			_ = testcontainers.TerminateContainer(in.container)
		}
	}
}

// CheckDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func CheckDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the PostgreSQL container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the PostgreSQL image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the PostgreSQL container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

// connector is pgx's, with one change to what its rows report: see rows.
type connector struct{ driver.Connector }

func (c connector) Connect(ctx context.Context) (driver.Conn, error) {
	dc, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return conn{dc}, nil
}

// conn passes the tape what it calls on a connection, and wraps the rows.
type conn struct{ driver.Conn }

func (c conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

func (c conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	r, err := c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return rows{r}, nil
}

func (c conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
}

func (c conn) Ping(ctx context.Context) error { return c.Conn.(driver.Pinger).Ping(ctx) }

func (c conn) ResetSession(ctx context.Context) error {
	return c.Conn.(driver.SessionResetter).ResetSession(ctx)
}

// rows names a column's type as pgx does, except a type pgx does not know:
// pgx names that by its OID, which is the server's next free one when the
// type was made, and so differs from run to run. A tape holds - for it.
type rows struct{ driver.Rows }

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

// Codec renders pgx's *pgconn.PgError as its SQLSTATE and message on one
// line, then a line for each other field the server sent:
//
//	42601 syntax error at or near "FROM"
//	position: 8
//
// A line break inside a field is written \n, and a backslash \\. The
// source file, line and routine the server names are left out: they say
// where in PostgreSQL's code the error was raised, and move with a patch
// release that changes nothing a test reads.
type Codec struct{}

// errorFields are the fields a tape holds after the first line, in order.
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

// plainSeverity is the severity of nearly every error, which a tape leaves
// out.
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
		return fmt.Errorf("tape: a PostgreSQL error %q", text)
	}
	pe := &pgconn.PgError{Severity: plainSeverity, Code: code, Message: unescaper.Replace(message)}
fields:
	for _, line := range lines[1:] {
		name, v, ok := strings.Cut(line, ": ")
		if !ok {
			return fmt.Errorf("tape: a PostgreSQL error field %q", line)
		}
		switch name {
		case "position", "internal position":
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return fmt.Errorf("tape: a PostgreSQL error field %q", line)
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
		return fmt.Errorf("tape: a PostgreSQL error field %q", line)
	}
	pe.SeverityUnlocalized = pe.Severity
	return pe
}
