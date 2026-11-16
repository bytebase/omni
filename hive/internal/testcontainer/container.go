// Package testcontainer is how hive's tests reach Hive. Open returns a
// connection through a tape (enginetest/tape), which replays the pinned
// server's answers and starts nothing; in a recording (OMNI_RECORD=1) the
// connection is to the pinned server itself. One testcontainer per test binary
// is started by the first test that asks for it and removed by Terminate,
// which the package's TestMain calls.
//
// The pinned server is a HiveServer2 with the image's embedded Derby
// metastore, configured by hive-site.xml beside this file: Hive 4 stores a
// CREATE TABLE as an external table, and the sync's managed-table branch needs
// one, so strict managed tables are off and the warehouse stays on the
// container's filesystem.
//
// gohive v2 hands the driver one thing a tape cannot hold: its error's text
// renders Hive's TStatus as pointer addresses, and the message and error code
// live in fields of a type it does not export. The live connector therefore
// converts an answer that failed into *Error, reading those two fields, and the
// Codec writes them down; so a recording and a replay hand the driver the same
// error. What the driver reports in production still carries gohive's
// rendering — see hive/driver.
package testcontainer

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	gohive "github.com/beltran/gohive/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is hive/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the one Hive the tests ask.
const Default = "default"

// siteFile is the hive-site.xml the pinned server starts with.
var siteFile = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "hive-site.xml")
}()

// user is the login HiveServer2's PLAIN SASL exchange needs. The image runs an
// auth=NONE server that accepts any name.
const user = "hive"

// Error is a statement HiveServer2 answered with a failure: Hive's own error
// code and message, which gohive v2 keeps in fields its error text does not
// show.
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string { return fmt.Sprintf("hive %d: %s", e.Code, e.Message) }

// Open returns t's connection to database, or to no database when it is empty.
// Each call is a stream of its own.
func Open(t testing.TB, database string) *sql.DB {
	t.Helper()
	return tape.Open(t, &tape.Server{Pins: Pins, Live: live(database), Codec: Codec{}}, Default)
}

// live connects to the pinned server, starting it on the first call. It is
// called only in a recording.
func live(database string) func(testing.TB, enginetest.Pin) sqldriver.Connector {
	return func(t testing.TB, pin enginetest.Pin) sqldriver.Connector {
		st := started(t, pin)
		return &connector{live: gohive.NewConnector(gohive.Config{
			Host:     st.host,
			Port:     st.port,
			Auth:     "NONE",
			Username: user,
			Database: database,
			Service:  "hive",
		})}
	}
}

// connector adapts gohive's connector to what a tape records: gohive's
// connection implements the deprecated Execer and Queryer (arguments as
// driver.Value) and nothing else, and its answers' errors are Hive's, spelled
// the way *Error spells them.
type connector struct {
	live sqldriver.Connector
}

func (c *connector) Connect(ctx context.Context) (sqldriver.Conn, error) {
	live, err := c.live.Connect(ctx)
	if err != nil {
		return nil, normalize(err)
	}
	return &conn{Conn: live}, nil
}

func (c *connector) Driver() sqldriver.Driver { return c.live.Driver() }

type conn struct {
	sqldriver.Conn
}

func (c *conn) ExecContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Result, error) {
	if len(args) > 0 {
		return nil, errors.New("hive: a statement with arguments; Hive has no prepared statements")
	}
	ex, ok := c.Conn.(interface {
		ExecContext(context.Context, string, []sqldriver.Value) (sqldriver.Result, error)
	})
	if !ok {
		return nil, errors.New("hive: the client has no ExecContext")
	}
	res, err := ex.ExecContext(ctx, query, nil)
	return res, normalize(err)
}

func (c *conn) QueryContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Rows, error) {
	if len(args) > 0 {
		return nil, errors.New("hive: a statement with arguments; Hive has no prepared statements")
	}
	q, ok := c.Conn.(interface {
		QueryContext(context.Context, string, []sqldriver.Value) (sqldriver.Rows, error)
	})
	if !ok {
		return nil, errors.New("hive: the client has no QueryContext")
	}
	result, err := q.QueryContext(ctx, query, nil)
	if err != nil {
		return nil, normalize(err)
	}
	return &answerRows{Rows: result}, nil
}

// answerRows normalizes the error a read answers: Hive reports an operation
// without a result set from the read rather than from the statement, and gohive
// raises it as the text of a TStatus.
type answerRows struct {
	sqldriver.Rows
}

func (r *answerRows) Next(dest []sqldriver.Value) error { return normalize(r.Rows.Next(dest)) }

func (r *answerRows) ColumnTypeDatabaseTypeName(index int) string {
	if typed, ok := r.Rows.(sqldriver.RowsColumnTypeDatabaseTypeName); ok {
		return typed.ColumnTypeDatabaseTypeName(index)
	}
	return ""
}

// normalize turns gohive's error into one a tape can hold: Hive's message and
// error code as *Error, or, for the answer gohive reports as the text of a
// TStatus, that text with the addresses masked as *statusError. An error that is
// neither is returned as it is, and a recording refuses it.
func normalize(err error) error {
	if err == nil {
		return nil
	}
	if code, message, ok := hiveAnswer(err); ok {
		return &Error{Code: code, Message: message}
	}
	if text, ok := maskStatus(err.Error()); ok {
		return &statusError{text: text}
	}
	return err
}

// pointerRe matches an address gohive prints in place of a TStatus field it
// holds as a pointer.
var pointerRe = regexp.MustCompile(`0x[0-9a-f]+`)

// statusError is a statement HiveServer2 answered with a failure gohive has no
// typed error for: it renders Hive's TStatus, whose fields it keeps as
// pointers, so the addresses are masked and the rest of the text is the answer.
// Hive answers a statement whose operation has no result set this way when its
// rows are read.
type statusError struct {
	text string
}

func (e *statusError) Error() string { return e.text }

// maskStatus returns s with gohive's addresses masked, and reports whether s is
// a TStatus rendering at all.
func maskStatus(s string) (string, bool) {
	if !strings.Contains(s, "TStatus({") {
		return "", false
	}
	return pointerRe.ReplaceAllString(s, "0x"), true
}

// hiveAnswer reads the error code and message gohive keeps of a failed
// statement. They are exported fields of gohive's unexported hiveError, which
// embeds the rendering of Hive's TStatus that its Error() returns.
func hiveAnswer(err error) (code int, message string, ok bool) {
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct || v.Type().String() != "gohive.hiveError" {
		return 0, "", false
	}
	messageField := v.FieldByName("Message")
	codeField := v.FieldByName("ErrorCode")
	if !messageField.IsValid() || messageField.Kind() != reflect.String || !codeField.IsValid() || !codeField.CanInt() {
		return 0, "", false
	}
	return int(codeField.Int()), messageField.String(), true
}

// --- the live server ---

// shared is the one pinned Hive this run started.
var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	host      string
	port      int
	// databases is what the image started with, which ResetToInitial leaves.
	databases map[string]bool
	err       error
}

// started returns the address of the pinned Hive, starting it on the first
// call and failing t unless it runs pin.
func started(t testing.TB, pin enginetest.Pin) *sharedServer {
	t.Helper()
	shared.once.Do(func() {
		shared.pin = pin
		shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start Hive: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started Hive at %s, not %s", shared.pin, pin)
	}
	return &sharedServer{host: shared.host, port: shared.port}
}

type sharedServer struct {
	host string
	port int
}

func start(ctx context.Context, pin enginetest.Pin) (err error) {
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{"10000/tcp"},
			Env:          map[string]string{"SERVICE_NAME": "hiveserver2"},
			Files: []testcontainers.ContainerFile{{
				HostFilePath:      siteFile,
				ContainerFilePath: "/opt/hive/conf/hive-site.xml",
				FileMode:          0o644,
			}},
			WaitingFor: wait.ForLog("Starting HiveServer2").WithStartupTimeout(5 * time.Minute),
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
	mapped, err := c.MappedPort(ctx, "10000/tcp")
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(mapped.Port())
	if err != nil {
		return err
	}
	db := gohive.OpenDB(gohive.Config{
		Host: host, Port: port, Auth: "NONE", Username: user, Database: "default", Service: "hive",
	})
	defer db.Close()
	// The image reports HiveServer2 started before its thrift endpoint
	// answers the first statement, so ask until it does.
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
		return fmt.Errorf("HiveServer2 does not answer: %w", err)
	}
	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return err
	}
	if !strings.HasPrefix(version, pin.Tag[len("rel/release-"):]+" ") {
		return fmt.Errorf("the server reports %q, but pins.tsv pins Hive %s", version, pin.Tag)
	}
	shared.host, shared.port = host, port

	databases, err := queryStrings(ctx, db, "SHOW DATABASES")
	if err != nil {
		return err
	}
	shared.databases = map[string]bool{}
	for _, name := range databases {
		shared.databases[name] = true
	}
	return nil
}

// Terminate removes the server, if this run started one.
func Terminate() {
	if shared.container != nil {
		_ = testcontainers.TerminateContainer(shared.container)
	}
}

// CheckDigest fails unless c runs the image with the pinned repository digest.
// A tag can be pushed again; a digest cannot.
func CheckDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Hive container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Hive image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Hive container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

// admin returns a connection to the pinned server outside any tape, which only
// a recording opens.
func admin(t testing.TB) *sql.DB {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer: a replay reaches no server")
	}
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatalf("%v", err)
	}
	st := started(t, pin)
	db := gohive.OpenDB(gohive.Config{
		Host: st.host, Port: st.port, Auth: "NONE", Username: user, Database: "default", Service: "hive",
	})
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// Reset runs fn on the pinned server, in a recording only: it puts the server
// into the state the requests after it start from, such as dropping what an
// earlier test left, so a replay needs nothing. fn's connection is not
// recorded.
func Reset(t testing.TB, fn func(ctx context.Context, db *sql.DB) error) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	db := admin(t)
	if err := fn(context.Background(), db); err != nil {
		t.Fatalf("reset Hive: %v", err)
	}
}

// ResetToInitial drops every database the pinned image did not start with, so
// a test that reads the server's database list sees a fixed set. It runs in a
// recording only.
func ResetToInitial(t testing.TB) {
	t.Helper()
	Reset(t, func(ctx context.Context, db *sql.DB) error {
		databases, err := queryStrings(ctx, db, "SHOW DATABASES")
		if err != nil {
			return err
		}
		for _, name := range databases {
			if shared.databases[name] {
				continue
			}
			if _, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+quoteBacktick(name)+" CASCADE"); err != nil {
				return fmt.Errorf("drop %s: %w", name, err)
			}
		}
		return nil
	})
}

func quoteBacktick(name string) string { return "`" + strings.ReplaceAll(name, "`", "``") + "`" }

// queryStrings reads the first column of every row.
func queryStrings(ctx context.Context, db *sql.DB, query string) ([]string, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Codec writes *Error on a tape as its code on one line and its message on the
// lines after, and hides the two answers that name when a table was made.
type Codec struct{}

// EncodeError implements tape.Codec. An error with Hive's code and message is
// written as its code on the first line; one Hive answered with a TStatus text
// is written with a dash for the code and the text after it.
func (Codec) EncodeError(err error) (string, bool) {
	var status *statusError
	if errors.As(err, &status) {
		return "-\n" + status.text, true
	}
	var e *Error
	if !errors.As(err, &e) {
		return "", false
	}
	return fmt.Sprintf("%d\n%s", e.Code, e.Message), true
}

// DecodeError implements tape.Codec.
func (Codec) DecodeError(text string) error {
	head, message, ok := strings.Cut(text, "\n")
	if !ok {
		return fmt.Errorf("testcontainer: unreadable Hive error on tape: %q", text)
	}
	if head == "-" {
		return &statusError{text: message}
	}
	code, err := strconv.Atoi(head)
	if err != nil {
		return fmt.Errorf("testcontainer: unreadable Hive error on tape: %q", text)
	}
	return &Error{Code: code, Message: message}
}

// madeUpRe matches the two values Hive answers with the moment a table was
// written and nothing else: the epoch in transient_lastDdlTime (and in a
// managed table's last_modified_time), and the creation time DESCRIBE
// FORMATTED spells. A second recording writes no different bytes because both
// get a stand-in. A fixture whose own data is a ten-digit number would be
// narrowed too; none is.
var madeUpRe = regexp.MustCompile(`\b\d{10}\b|\b(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun) (?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) [ 0-9]\d \d\d:\d\d:\d\d [A-Z]{2,5} \d{4}\b`)

// MadeUp matches what a new container makes up again.
func (Codec) MadeUp() *regexp.Regexp { return madeUpRe }

// StandIn returns what stands for a made-up value. Every value of one kind gets
// the same stand-in, so a tape does not depend on how many distinct moments a
// recording's tables happened to be made in.
func (Codec) StandIn(name string, _ int) string {
	if _, err := strconv.Atoi(name); err == nil {
		return "1000000000"
	}
	return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Format("Mon Jan _2 15:04:05 MST 2006")
}
