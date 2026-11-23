// Package testcontainer is how trino's tests reach Trino: through a tape
// (enginetest/tape) that replays the pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. That server is
// one testcontainer per test binary, started by the first test that records
// and removed by Terminate, which each package's TestMain calls.
//
// A test asks one thing of Trino: whether its parser accepts a statement
// (Oracle.CheckSyntax). Trino raises errorName SYNTAX_ERROR when its parser
// rejects the input, and for a few checks after the parse that read only the
// text, such as a name with too many parts. Every other outcome, success or a
// semantic error (TABLE_NOT_FOUND, NOT_SUPPORTED, ...), means the grammar
// accepted the statement. A statement may run, against the memory catalog the
// session defaults to; in a recording each top-level test starts with that
// catalog empty.
//
// The oracle's tape is the statement endpoint's: each request is one statement,
// and each answer ok, or the error Trino answered (QueryError).
//
// trino/driver asks more of the server, so OpenDB returns the same tape around
// the client library the driver uses: the tape holds the statements the driver
// sends and the rows the client decoded, and each answer is ok, the client's
// error, or a result set.
package testcontainer

import (
	"bytes"
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	trinoclient "github.com/trinodb/trino-go-client/trino"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is trino/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the one Trino the tests ask.
const Default = "default"

var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}

// syntaxErrorName is Trino's errorName for a syntax error.
const syntaxErrorName = "SYNTAX_ERROR"

// An Oracle answers whether Trino accepts a statement's syntax.
type Oracle struct {
	run func(ctx context.Context, sql string) error
}

// Open returns t's Oracle: a replay of its tape, or the pinned Trino in a
// recording. Each call is a stream of its own.
func Open(t testing.TB) *Oracle {
	t.Helper()
	db := tape.Open(t, server, Default)
	return &Oracle{run: func(ctx context.Context, sql string) error {
		_, err := db.ExecContext(ctx, sql)
		return err
	}}
}

// Judge returns the pinned Trino itself, outside any tape, for the judge of
// compat.tsv, which runs only in a recording.
func Judge(t testing.TB) *Oracle {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return &Oracle{run: newClient(started(t, pin)).run}
}

// Result is the outcome of a syntax check.
type Result struct {
	// Accepted is true unless Trino answered SYNTAX_ERROR.
	Accepted bool
	// ErrorName is Trino's errorName, "" when the statement ran without error.
	ErrorName string
	// ErrorCode is Trino's numeric errorCode (0 when no error).
	ErrorCode int
	// Message is Trino's error message, if any.
	Message string
}

// CheckSyntax sends sql to Trino and reports whether it accepted its syntax.
// An error is no answer from Trino at all.
func (o *Oracle) CheckSyntax(ctx context.Context, sql string) (Result, error) {
	err := o.Run(ctx, sql)
	var qe *QueryError
	switch {
	case err == nil:
		return Result{Accepted: true}, nil
	case errors.As(err, &qe):
		return Result{
			Accepted:  qe.Name != syntaxErrorName,
			ErrorName: qe.Name,
			ErrorCode: qe.Code,
			Message:   qe.Message,
		}, nil
	default:
		return Result{}, err
	}
}

// Run sends sql to Trino. It returns nil when the statement parsed and was
// planned, Trino's answer as a *QueryError when it failed, and any other
// error when Trino gave no answer.
func (o *Oracle) Run(ctx context.Context, sql string) error { return o.run(ctx, sql) }

// QueryError is a Trino server's answer to a statement that failed, as its
// client protocol's QueryError carries it.
type QueryError struct {
	// Name is errorName and Code errorCode: SYNTAX_ERROR and 1 for a
	// rejection, TABLE_NOT_FOUND for a missing table.
	Name string
	Code int
	// Line and Column are errorLocation, from 1; 0 for none.
	Line, Column int
	Message      string
}

func (e *QueryError) Error() string { return fmt.Sprintf("trino %s: %s", e.Name, e.Message) }

// Codec writes a QueryError on a tape as its name, code and location on one
// line, and its message on the lines after.
type Codec struct{}

// EncodeError implements tape.Codec.
func (Codec) EncodeError(err error) (string, bool) {
	var qe *QueryError
	if !errors.As(err, &qe) {
		return "", false
	}
	return fmt.Sprintf("%s %d %d:%d\n%s", qe.Name, qe.Code, qe.Line, qe.Column, qe.Message), true
}

// DecodeError implements tape.Codec.
func (Codec) DecodeError(text string) error {
	head, msg, _ := strings.Cut(text, "\n")
	f := strings.Fields(head)
	if len(f) != 3 {
		return fmt.Errorf("testcontainer: unreadable Trino error on tape: %q", text)
	}
	code, err1 := strconv.Atoi(f[1])
	line, col, ok := strings.Cut(f[2], ":")
	l, err2 := strconv.Atoi(line)
	c, err3 := strconv.Atoi(col)
	if err1 != nil || err2 != nil || err3 != nil || !ok {
		return fmt.Errorf("testcontainer: unreadable Trino error on tape: %q", text)
	}
	return &QueryError{Name: f[0], Code: code, Line: l, Column: c, Message: msg}
}

// --- the live server ---

// memoryLimit caps the container: Trino sizes its JVM to the memory it
// sees, and the box it records on is shared.
const memoryLimit = 4 << 30

// shared is the one pinned Trino this run started.
var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	url       string
	err       error

	mu      sync.Mutex
	emptied map[string]bool // the top-level tests the memory catalog was emptied for
}

// started returns the URL of the pinned Trino, starting it on the first
// call and failing t unless it runs pin.
func started(t testing.TB, pin enginetest.Pin) string {
	t.Helper()
	shared.once.Do(func() {
		shared.pin = pin
		shared.url, shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start Trino: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started Trino at %s, not %s", shared.pin, pin)
	}
	return shared.url
}

// live connects to the pinned Trino. On t's top-level test's first
// connection it empties the memory catalog of what other tests made, so a
// test reads the same answers recorded alone as in a run of its package.
func live(t testing.TB, pin enginetest.Pin) sqldriver.Connector {
	t.Helper()
	c := newClient(started(t, pin))
	emptyOnce(t, c)
	return connector{c}
}

// emptyOnce empties the memory catalog on t's top-level test's first
// connection, so a test reads the same answers recorded alone as in a run of
// its package.
func emptyOnce(t testing.TB, c *client) {
	t.Helper()
	top := enginetest.TopName(t.Name())
	shared.mu.Lock()
	defer shared.mu.Unlock()
	if !shared.emptied[top] {
		if err := c.emptyMemory(context.Background()); err != nil {
			t.Fatalf("empty Trino's memory catalog: %v", err)
		}
		if shared.emptied == nil {
			shared.emptied = map[string]bool{}
		}
		shared.emptied[top] = true
	}
}

// OpenDB returns t's connection to the pinned server in catalog, as
// trino/driver opens one: a *sql.DB over the tape whose live side is the client
// library the driver links. In a recording the client talks to the pinned
// server; in a replay its requests and the rows it decoded come from the tape.
// The pool holds one connection, so a statement's session state holds for the
// next, as on a pinned *sql.Conn.
func OpenDB(t testing.TB, catalog, schema string) *sql.DB {
	t.Helper()
	s := &tape.Server{
		Pins:  Pins,
		Live:  clientLive(catalog, schema),
		Codec: ClientCodec{},
	}
	return tape.Open(t, s, Default)
}

// clientLive connects t's tape to the pinned server with the client library,
// as trino/driver does, in catalog and schema.
func clientLive(catalog, schema string) func(testing.TB, enginetest.Pin) sqldriver.Connector {
	return func(t testing.TB, pin enginetest.Pin) sqldriver.Connector {
		c := newClient(started(t, pin))
		emptyOnce(t, c)
		return clientConnector{url: c.baseURL, catalog: catalog, schema: schema}
	}
}

// clientConnector opens client library connections to the pinned server.
type clientConnector struct {
	url             string
	catalog, schema string
}

func (c clientConnector) Connect(ctx context.Context) (sqldriver.Conn, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	u.User = url.User("omni-driver")
	query := u.Query()
	query.Set("catalog", c.catalog)
	query.Set("schema", c.schema)
	query.Set("source", "omni-driver")
	u.RawQuery = query.Encode()
	conn, err := (&trinoclient.Driver{}).Open(u.String())
	if err != nil {
		return nil, err
	}
	return clientConn{Conn: conn}, nil
}

func (c clientConnector) Driver() sqldriver.Driver { return &trinoclient.Driver{} }

// clientConn adapts the client's connection to what a tape records: the client
// itself offers only Prepare, and database/sql reaches Exec and Query through a
// statement, so the adapter does the same.
type clientConn struct {
	sqldriver.Conn
}

var (
	_ sqldriver.ExecerContext      = clientConn{}
	_ sqldriver.QueryerContext     = clientConn{}
	_ sqldriver.SessionResetter    = clientConn{}
	_ sqldriver.ConnPrepareContext = clientConn{}
)

// PrepareContext prepares one statement. The client offers only this path to
// Exec and Query, which is also how database/sql would reach it.
func (c clientConn) PrepareContext(ctx context.Context, query string) (sqldriver.Stmt, error) {
	preparer, ok := c.Conn.(sqldriver.ConnPrepareContext)
	if !ok {
		return nil, errors.New("testcontainer: the Trino client connection does not prepare")
	}
	return preparer.PrepareContext(ctx, query)
}

func (c clientConn) ExecContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Result, error) {
	stmt, err := c.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	exec, ok := stmt.(sqldriver.StmtExecContext)
	if !ok {
		return nil, errors.New("testcontainer: the Trino client statement does not exec")
	}
	return exec.ExecContext(ctx, args)
}

func (c clientConn) QueryContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Rows, error) {
	stmt, err := c.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	queryer, ok := stmt.(sqldriver.StmtQueryContext)
	if !ok {
		stmt.Close()
		return nil, errors.New("testcontainer: the Trino client statement does not query")
	}
	rows, err := queryer.QueryContext(ctx, args)
	if err != nil {
		stmt.Close()
		return nil, err
	}
	return &statementRows{Rows: rows, stmt: stmt}, nil
}

func (c clientConn) ResetSession(ctx context.Context) error {
	if resetter, ok := c.Conn.(sqldriver.SessionResetter); ok {
		return resetter.ResetSession(ctx)
	}
	return nil
}

// statementRows closes the statement a query came from with its rows, and
// forwards the column metadata the client's rows report: a tape holds a column's
// type only when the rows implement RowsColumnTypeDatabaseTypeName.
type statementRows struct {
	sqldriver.Rows
	stmt sqldriver.Stmt
}

var (
	_ sqldriver.RowsColumnTypeScanType         = (*statementRows)(nil)
	_ sqldriver.RowsColumnTypeDatabaseTypeName = (*statementRows)(nil)
	_ sqldriver.RowsColumnTypeLength           = (*statementRows)(nil)
	_ sqldriver.RowsColumnTypePrecisionScale   = (*statementRows)(nil)
	_ sqldriver.RowsColumnTypeNullable         = (*statementRows)(nil)
)

func (r *statementRows) Close() error {
	err := r.Rows.Close()
	if cerr := r.stmt.Close(); err == nil {
		err = cerr
	}
	return err
}

func (r *statementRows) ColumnTypeDatabaseTypeName(index int) string {
	if typed, ok := r.Rows.(sqldriver.RowsColumnTypeDatabaseTypeName); ok {
		return typed.ColumnTypeDatabaseTypeName(index)
	}
	return ""
}

func (r *statementRows) ColumnTypeScanType(index int) reflect.Type {
	if typed, ok := r.Rows.(sqldriver.RowsColumnTypeScanType); ok {
		return typed.ColumnTypeScanType(index)
	}
	return reflect.TypeOf(new(any)).Elem()
}

func (r *statementRows) ColumnTypeLength(index int) (int64, bool) {
	if typed, ok := r.Rows.(sqldriver.RowsColumnTypeLength); ok {
		return typed.ColumnTypeLength(index)
	}
	return 0, false
}

func (r *statementRows) ColumnTypePrecisionScale(index int) (int64, int64, bool) {
	if typed, ok := r.Rows.(sqldriver.RowsColumnTypePrecisionScale); ok {
		return typed.ColumnTypePrecisionScale(index)
	}
	return 0, 0, false
}

func (r *statementRows) ColumnTypeNullable(index int) (bool, bool) {
	if typed, ok := r.Rows.(sqldriver.RowsColumnTypeNullable); ok {
		return typed.ColumnTypeNullable(index)
	}
	return false, false
}

// ClientCodec writes a Trino query error on a tape as the client's fields and
// reads it back into the same type, so a replay's errors.As finds what a
// recording's did.
type ClientCodec struct{}

// tapeClientError is the client's ErrQueryFailed and its ErrTrino reason as a
// tape holds them.
type tapeClientError struct {
	Status   int
	Type     string
	Message  string
	SqlState string
	Name     string
	Code     int
	Line     int
	Column   int
}

// EncodeError implements tape.Codec.
func (ClientCodec) EncodeError(err error) (string, bool) {
	var failed *trinoclient.ErrQueryFailed
	if !errors.As(err, &failed) {
		return "", false
	}
	var reason *trinoclient.ErrTrino
	if !errors.As(err, &reason) {
		return "", false
	}
	b, merr := json.Marshal(tapeClientError{
		Status:   failed.StatusCode,
		Type:     reason.ErrorType,
		Message:  reason.Message,
		SqlState: reason.SqlState,
		Name:     reason.ErrorName,
		Code:     reason.ErrorCode,
		Line:     reason.ErrorLocation.LineNumber,
		Column:   reason.ErrorLocation.ColumnNumber,
	})
	if merr != nil {
		return "", false
	}
	return string(b), true
}

// DecodeError implements tape.Codec.
func (ClientCodec) DecodeError(text string) error {
	var e tapeClientError
	if err := json.Unmarshal([]byte(text), &e); err != nil {
		return fmt.Errorf("testcontainer: unreadable Trino client error on tape: %w", err)
	}
	return &trinoclient.ErrQueryFailed{
		StatusCode: e.Status,
		Reason: &trinoclient.ErrTrino{
			Message:       e.Message,
			SqlState:      e.SqlState,
			ErrorCode:     e.Code,
			ErrorName:     e.Name,
			ErrorType:     e.Type,
			ErrorLocation: trinoclient.ErrorLocation{LineNumber: e.Line, ColumnNumber: e.Column},
		},
	}
}

func start(ctx context.Context, pin enginetest.Pin) (url string, err error) {
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.ImageRef(),
			ExposedPorts: []string{"8080/tcp"},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.Memory = memoryLimit
			},
		},
		Started: true,
	})
	if c != nil {
		shared.container = c
	}
	if err != nil {
		return "", err
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := c.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", err
	}
	url = fmt.Sprintf("http://%s:%s", host, port.Port())
	cl := newClient(url)
	// GenericContainer returns once the process is up; poll until Trino
	// reports it has finished starting.
	deadline := time.Now().Add(3 * time.Minute)
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		version, perr := cl.ping(pingCtx)
		cancel()
		if perr == nil {
			if version != pin.Tag {
				return "", fmt.Errorf("the server reports version %q, but pins.tsv pins Trino %s", version, pin.Tag)
			}
			return url, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("trino did not become ready within timeout: %w", perr)
		}
		time.Sleep(2 * time.Second)
	}
}

// checkDigest proves the container runs the pinned image. The pull names
// image@digest, so a moved tag cannot select another image; this check reads
// what the daemon actually started, which is where the recorded answers came
// from.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Trino container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Trino image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Trino container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

// Terminate removes the server this run started, if any. A package's
// TestMain calls it.
func Terminate() {
	if shared.container != nil {
		_ = testcontainers.TerminateContainer(shared.container)
	}
}

// connector is a database/sql connector over the statement endpoint, so
// enginetest/tape can record it: ExecContext sends one statement.
type connector struct{ c *client }

func (c connector) Connect(context.Context) (sqldriver.Conn, error) { return conn(c), nil }
func (c connector) Driver() sqldriver.Driver                        { return nil }

type conn struct{ c *client }

var _ sqldriver.ExecerContext = conn{}

func (conn) Prepare(string) (sqldriver.Stmt, error) {
	return nil, errors.New("testcontainer: Trino statements are not prepared")
}
func (conn) Close() error                 { return nil }
func (conn) Begin() (sqldriver.Tx, error) { return nil, errors.New("testcontainer: no transactions") }

func (c conn) ExecContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Result, error) {
	if len(args) > 0 {
		return nil, errors.New("testcontainer: Trino statements take no arguments")
	}
	if err := c.c.run(ctx, query); err != nil {
		return nil, err
	}
	return noCount{}, nil
}

// noCount is the result of a statement: the endpoint counts nothing the
// tests read.
type noCount struct{}

func (noCount) LastInsertId() (int64, error) { return 0, errors.New("no id") }
func (noCount) RowsAffected() (int64, error) { return 0, errors.New("no count") }

// --- the statement endpoint ---

// client talks to a Trino server over its REST statement protocol, as the
// memory catalog's default schema.
type client struct {
	baseURL string
	http    *http.Client
}

func newClient(baseURL string) *client {
	return &client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 30 * time.Second}}
}

// queryResp models the relevant fields of a /v1/statement response.
type queryResp struct {
	ID         string `json:"id"`
	NextURI    string `json:"nextUri"`
	UpdateType string `json:"updateType"`
	Columns    []struct {
		Name string `json:"name"`
	} `json:"columns"`
	Data  [][]any `json:"data"`
	Stats struct {
		State string `json:"state"`
	} `json:"stats"`
	Error *struct {
		Message       string `json:"message"`
		ErrorCode     int    `json:"errorCode"`
		ErrorName     string `json:"errorName"`
		ErrorLocation *struct {
			LineNumber   int `json:"lineNumber"`
			ColumnNumber int `json:"columnNumber"`
		} `json:"errorLocation"`
	} `json:"error"`
}

// queryError reads the error of a response, or nil.
func (r queryResp) queryError() *QueryError {
	if r.Error == nil {
		return nil
	}
	e := &QueryError{Name: r.Error.ErrorName, Code: r.Error.ErrorCode, Message: r.Error.Message}
	if l := r.Error.ErrorLocation; l != nil {
		e.Line, e.Column = l.LineNumber, l.ColumnNumber
	}
	return e
}

// ping verifies the server is reachable and finished starting, returning
// its version.
func (c *client) ping(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/info", nil)
	if err != nil {
		return "", err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("trino /v1/info: status %d", resp.StatusCode)
	}
	var in struct {
		Starting    bool `json:"starting"`
		NodeVersion struct {
			Version string `json:"version"`
		} `json:"nodeVersion"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&in); err != nil {
		return "", err
	}
	if in.Starting {
		return in.NodeVersion.Version, fmt.Errorf("trino is still starting")
	}
	return in.NodeVersion.Version, nil
}

// run sends sql and follows it until it fails or is planned. A query that
// reaches RUNNING is cancelled there, so large results are not streamed; a
// statement that changes something (an updateType) runs to the end, so the
// next statement finds its change made.
func (c *client) run(ctx context.Context, sql string) error {
	_, err := c.follow(ctx, sql, false)
	return err
}

// query runs sql to the end and returns its rows.
func (c *client) query(ctx context.Context, sql string) ([][]any, error) {
	return c.follow(ctx, sql, true)
}

func (c *client) follow(ctx context.Context, sql string, all bool) ([][]any, error) {
	resp, err := c.post(ctx, sql)
	if err != nil {
		return nil, err
	}
	var rows [][]any
	for i := 0; i < 10000; i++ {
		if qe := resp.queryError(); qe != nil {
			return nil, qe
		}
		rows = append(rows, resp.Data...)
		if resp.NextURI == "" {
			return rows, nil
		}
		if !all && resp.UpdateType == "" && (resp.Stats.State == "RUNNING" || resp.Stats.State == "FINISHED") {
			// Parsing, analysis and planning have all succeeded.
			c.cancel(resp.NextURI)
			return nil, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
		if resp, err = c.get(ctx, resp.NextURI); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("trino: statement %q did not settle", truncate(sql, 60))
}

// emptyMemory drops every table, view and schema of the memory catalog but
// its default schema, which the image starts with empty.
func (c *client) emptyMemory(ctx context.Context) error {
	rows, err := c.query(ctx, "SELECT table_schema, table_name, table_type FROM memory.information_schema.tables "+
		"WHERE table_schema <> 'information_schema' ORDER BY 1, 2")
	if err != nil {
		return err
	}
	var stmts []string
	for _, r := range rows {
		kind := "TABLE"
		if r[2] == "VIEW" {
			kind = "VIEW"
		}
		stmts = append(stmts, fmt.Sprintf("DROP %s memory.%s.%s", kind, quote(r[0]), quote(r[1])))
	}
	schemas, err := c.query(ctx, "SELECT schema_name FROM memory.information_schema.schemata "+
		"WHERE schema_name NOT IN ('information_schema', 'default') ORDER BY 1")
	if err != nil {
		return err
	}
	for _, r := range schemas {
		stmts = append(stmts, "DROP SCHEMA memory."+quote(r[0]))
	}
	for _, s := range stmts {
		if _, err := c.query(ctx, s); err != nil {
			return fmt.Errorf("%s: %w", s, err)
		}
	}
	return nil
}

// quote writes a name as a quoted identifier.
func quote(name any) string {
	return `"` + strings.ReplaceAll(fmt.Sprint(name), `"`, `""`) + `"`
}

func (c *client) post(ctx context.Context, sql string) (queryResp, error) {
	return c.postIn(ctx, sql, memoryDefault)
}

// postIn sends sql to run in session in.
func (c *client) postIn(ctx context.Context, sql string, in session) (queryResp, error) {
	// Trino can answer 503 SERVER_STARTING_UP briefly; retry a few times.
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/statement", bytes.NewBufferString(sql))
		if err != nil {
			return queryResp{}, err
		}
		req.Header.Set("X-Trino-User", "omni-oracle")
		req.Header.Set("X-Trino-Catalog", in.catalog)
		req.Header.Set("X-Trino-Schema", in.schema)
		req.Header.Set("Content-Type", "text/plain")
		resp, err := c.http.Do(req)
		if err != nil {
			return queryResp{}, err
		}
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			lastErr = fmt.Errorf("trino /v1/statement: 503")
			select {
			case <-ctx.Done():
				return queryResp{}, ctx.Err()
			case <-time.After(time.Second):
			}
			continue
		}
		return decode(resp)
	}
	return queryResp{}, lastErr
}

func (c *client) get(ctx context.Context, uri string) (queryResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return queryResp{}, err
	}
	req.Header.Set("X-Trino-User", "omni-oracle")
	resp, err := c.http.Do(req)
	if err != nil {
		return queryResp{}, err
	}
	return decode(resp)
}

// cancel best-effort terminates a running query so result data is not streamed.
func (c *client) cancel(uri string) {
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return
	}
	req.Header.Set("X-Trino-User", "omni-oracle")
	if resp, err := c.http.Do(req); err == nil {
		resp.Body.Close()
	}
}

func decode(resp *http.Response) (queryResp, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return queryResp{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return queryResp{}, fmt.Errorf("trino: status %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var qr queryResp
	if err := json.Unmarshal(body, &qr); err != nil {
		return queryResp{}, fmt.Errorf("trino: decode response: %w", err)
	}
	return qr, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
