// Package testcontainer is how cassandra's tests reach Cassandra: through a
// tape (enginetest/tape) that replays a pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. The server is
// one testcontainer per test binary, started by the first test that records and
// removed by Terminate, which each package's TestMain calls.
//
// Cassandra's client, gocql, is not a database/sql driver, so the tape cannot
// sit at that seam. Client returns a conversation of the driver's own shape
// instead: a request is one statement the driver asks gocql for, with the page
// size and page state it asks with, and an answer is the columns, the rows and
// the state gocql returned for that page. Each cell is written with the Go type
// its CQL column decodes to, because that is what the driver's row conversion
// reads: a CQL time is a time.Duration, a decimal an *inf.Dec, and a list or a
// user-defined type its JSON spelling. container_test.go holds that encoding to
// every kind the driver converts.
//
// A tape records what gocql answered, not the wire traffic: the handshake, the
// prepared statements the sync queries would otherwise use, and gocql's paging
// between pages are not on it. The driver names its keyspace in the query text
// instead of binding it, so a request is stable; the page states the server
// returns are recorded as it sent them, and the corpus is small enough that
// they are the same bytes from one recording to the next.
package testcontainer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/inf.v0"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is cassandra/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Driver is the pin row of the server the driver's tapes were recorded on.
const Driver = "driver"

var server = &tape.Server{Pins: Pins, Codec: codec{}}

// codec is how a Cassandra answer is read for the tape. A query's error is part
// of its answer, so the tape needs no error spelling from the codec.
type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return errors.New(text) }

// Client is the driver's Cassandra client, answered by t's tape. It is what a
// test hands the driver in place of the gocql session; its methods are the
// calls the driver makes, and the recordings hold what gocql answered.
type Client struct {
	tape *tape.Raw
}

// Open returns t's client to the pinned server, in keyspace, which is the
// session's keyspace: an unqualified table name resolves there, as Config.Database
// makes it. A replay holds its tape; a recording holds the pinned server with
// every call written down.
func Open(t testing.TB, keyspace string) *Client {
	t.Helper()
	return &Client{tape: tape.OpenRaw(t, server, Driver, func(request string) (string, error) {
		// A recording reaches the server on its first call, so a test that
		// replays never starts a container.
		return liveDo(t, request, keyspace)
	})}
}

// Query reads one page of one statement.
func (c *Client) Query(ctx context.Context, statement string, pageSize int, pageState []byte) ([]gocql.ColumnInfo, []gocql.RowData, []byte, error) {
	_ = ctx
	a, err := c.do(request{
		Op:        opQuery,
		Statement: statement,
		PageSize:  pageSize,
		PageState: base64.StdEncoding.EncodeToString(pageState),
	})
	if err != nil {
		return nil, nil, nil, err
	}
	columns := make([]gocql.ColumnInfo, 0, len(a.Columns))
	for _, column := range a.Columns {
		columns = append(columns, gocql.ColumnInfo{
			Name:     column.Name,
			TypeInfo: tapeType(column.Type),
		})
	}
	rows := make([]gocql.RowData, 0, len(a.Rows))
	for _, values := range a.Rows {
		row := gocql.RowData{Columns: a.RowColumns, Values: make([]any, 0, len(values))}
		for _, value := range values {
			decoded, err := decodeCell(value)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("testcontainer: unreadable cell: %w", err)
			}
			row.Values = append(row.Values, decoded)
		}
		rows = append(rows, row)
	}
	next, err := base64.StdEncoding.DecodeString(a.PageState)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("testcontainer: unreadable page state: %w", err)
	}
	return columns, rows, next, nil
}

// Exec runs one statement for its effect.
func (c *Client) Exec(ctx context.Context, statement string) error {
	_ = ctx
	_, err := c.do(request{Op: opExec, Statement: statement})
	return err
}

// Close closes the taped conversation.
func (c *Client) Close() {}

// do sends one request and returns its answer.
func (c *Client) do(r request) (answer, error) {
	text, err := r.text()
	if err != nil {
		return answer{}, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return answer{}, err
	}
	var a answer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
		return answer{}, fmt.Errorf("testcontainer: unreadable answer to %s: %w", text, err)
	}
	if a.Error != nil {
		return answer{}, a.Error.err()
	}
	return a, nil
}

// The calls the driver makes of its client, as a tape names them.
const (
	opQuery = "query"
	opExec  = "exec"
)

// request is one call as a tape holds it.
type request struct {
	Op        string `json:"op"`
	Statement string `json:"statement"`
	PageSize  int    `json:"pageSize,omitempty"`
	PageState string `json:"pageState,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one call's answer as a tape holds it: the columns and rows of one
// page, the state that continues after it, or the error gocql returned.
type answer struct {
	Columns []column `json:"columns,omitempty"`
	// RowColumns are the names gocql gives the values of one row, which differ
	// from Columns only where a tuple column expands into one value per
	// element ("name[0]", "name[1]", ...).
	RowColumns []string     `json:"rowColumns,omitempty"`
	Rows       [][]cell     `json:"rows,omitempty"`
	PageState  string       `json:"pageState,omitempty"`
	Error      *clientError `json:"error,omitempty"`
}

// column is one column as a tape holds it: its name and the CQL type gocql
// described it with.
type column struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// clientError is an error the server answered, as a tape holds it. Cassandra's
// errors are its request errors, whose code and message are what a caller
// reads.
type clientError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// requestError is a gocql.RequestError a replay rebuilds. gocql returns its own
// unexported frame type, which a tape cannot reconstruct; a caller that reads
// the code and message finds them.
type requestError struct {
	code    int
	message string
}

func (e requestError) Code() int       { return e.code }
func (e requestError) Message() string { return e.message }
func (e requestError) Error() string   { return e.message }

func (e *clientError) err() error {
	return requestError{code: e.Code, message: e.Message}
}

// encodeError renders err as an answer's error half.
func encodeError(err error) *clientError {
	var requestErr gocql.RequestError
	if errors.As(err, &requestErr) {
		return &clientError{Code: requestErr.Code(), Message: requestErr.Message()}
	}
	return nil
}

// The Go types a cell is written as. A CQL column decodes to exactly one of
// them, and the driver's row conversion keys on it.
const (
	kindString       = "string"
	kindInt          = "int"
	kindInt64        = "int64"
	kindInt32        = "int32"
	kindInt16        = "int16"
	kindInt8         = "int8"
	kindBool         = "bool"
	kindFloat32      = "float32"
	kindFloat64      = "float64"
	kindBytes        = "bytes"
	kindTimestamp    = "timestamp"
	kindDuration     = "duration"
	kindLongDuration = "gocql.Duration"
	kindUUID         = "uuid"
	kindDecimal      = "inf.Dec"
	kindVarint       = "big.Int"
	kindJSON         = "json"
)

// cell is one value as a tape holds it: the Go type it decoded to, whether the
// column was NULL, and the value's spelling. A collection or a user-defined
// type is its JSON text, which is what the driver shows.
type cell struct {
	Kind   string `json:"kind"`
	Null   bool   `json:"null,omitempty"`
	Text   string `json:"text,omitempty"`
	Months int32  `json:"months,omitempty"`
	Days   int32  `json:"days,omitempty"`
	Nanos  int64  `json:"nanos,omitempty"`
}

// encodeCell reads one value gocql handed back.
func encodeCell(v any) (cell, error) {
	switch v := v.(type) {
	case *string:
		if v == nil {
			return cell{Kind: kindString, Null: true}, nil
		}
		return cell{Kind: kindString, Text: *v}, nil
	case *int:
		if v == nil {
			return cell{Kind: kindInt, Null: true}, nil
		}
		return cell{Kind: kindInt, Text: strconv.Itoa(*v)}, nil
	case *int64:
		if v == nil {
			return cell{Kind: kindInt64, Null: true}, nil
		}
		return cell{Kind: kindInt64, Text: strconv.FormatInt(*v, 10)}, nil
	case *int32:
		if v == nil {
			return cell{Kind: kindInt32, Null: true}, nil
		}
		return cell{Kind: kindInt32, Text: strconv.FormatInt(int64(*v), 10)}, nil
	case *int16:
		if v == nil {
			return cell{Kind: kindInt16, Null: true}, nil
		}
		return cell{Kind: kindInt16, Text: strconv.FormatInt(int64(*v), 10)}, nil
	case *int8:
		if v == nil {
			return cell{Kind: kindInt8, Null: true}, nil
		}
		return cell{Kind: kindInt8, Text: strconv.FormatInt(int64(*v), 10)}, nil
	case *bool:
		if v == nil {
			return cell{Kind: kindBool, Null: true}, nil
		}
		return cell{Kind: kindBool, Text: strconv.FormatBool(*v)}, nil
	case *float32:
		if v == nil {
			return cell{Kind: kindFloat32, Null: true}, nil
		}
		return cell{Kind: kindFloat32, Text: strconv.FormatFloat(float64(*v), 'g', -1, 32)}, nil
	case *float64:
		if v == nil {
			return cell{Kind: kindFloat64, Null: true}, nil
		}
		return cell{Kind: kindFloat64, Text: strconv.FormatFloat(*v, 'g', -1, 64)}, nil
	case *[]byte:
		if v == nil {
			return cell{Kind: kindBytes, Null: true}, nil
		}
		return cell{Kind: kindBytes, Text: base64.StdEncoding.EncodeToString(*v)}, nil
	case *time.Time:
		if v == nil {
			return cell{Kind: kindTimestamp, Null: true}, nil
		}
		return cell{Kind: kindTimestamp, Text: v.Format(time.RFC3339Nano)}, nil
	case *time.Duration:
		if v == nil {
			return cell{Kind: kindDuration, Null: true}, nil
		}
		return cell{Kind: kindDuration, Text: strconv.FormatInt(int64(*v), 10)}, nil
	case *gocql.Duration:
		if v == nil {
			return cell{Kind: kindLongDuration, Null: true}, nil
		}
		return cell{Kind: kindLongDuration, Months: v.Months, Days: v.Days, Nanos: v.Nanoseconds}, nil
	case *gocql.UUID:
		if v == nil {
			return cell{Kind: kindUUID, Null: true}, nil
		}
		return cell{Kind: kindUUID, Text: v.String()}, nil
	case **inf.Dec:
		if v == nil || *v == nil {
			return cell{Kind: kindDecimal, Null: true}, nil
		}
		return cell{Kind: kindDecimal, Text: (*v).String()}, nil
	case **big.Int:
		if v == nil || *v == nil {
			return cell{Kind: kindVarint, Null: true}, nil
		}
		return cell{Kind: kindVarint, Text: (*v).String()}, nil
	default:
		encoded, err := json.Marshal(v)
		if err != nil {
			return cell{}, fmt.Errorf("testcontainer: %T is not a value a tape can hold: %w", v, err)
		}
		return cell{Kind: kindJSON, Text: string(encoded)}, nil
	}
}

// decodeCell rebuilds the pointer gocql hands the driver for a cell.
func decodeCell(c cell) (any, error) {
	switch c.Kind {
	case kindString:
		if c.Null {
			return (*string)(nil), nil
		}
		v := c.Text
		return &v, nil
	case kindInt:
		if c.Null {
			return (*int)(nil), nil
		}
		v, err := strconv.Atoi(c.Text)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindInt64:
		if c.Null {
			return (*int64)(nil), nil
		}
		v, err := strconv.ParseInt(c.Text, 10, 64)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindInt32:
		if c.Null {
			return (*int32)(nil), nil
		}
		v, err := strconv.ParseInt(c.Text, 10, 32)
		if err != nil {
			return nil, err
		}
		n := int32(v)
		return &n, nil
	case kindInt16:
		if c.Null {
			return (*int16)(nil), nil
		}
		v, err := strconv.ParseInt(c.Text, 10, 16)
		if err != nil {
			return nil, err
		}
		n := int16(v)
		return &n, nil
	case kindInt8:
		if c.Null {
			return (*int8)(nil), nil
		}
		v, err := strconv.ParseInt(c.Text, 10, 8)
		if err != nil {
			return nil, err
		}
		n := int8(v)
		return &n, nil
	case kindBool:
		if c.Null {
			return (*bool)(nil), nil
		}
		v, err := strconv.ParseBool(c.Text)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindFloat32:
		if c.Null {
			return (*float32)(nil), nil
		}
		v, err := strconv.ParseFloat(c.Text, 32)
		if err != nil {
			return nil, err
		}
		n := float32(v)
		return &n, nil
	case kindFloat64:
		if c.Null {
			return (*float64)(nil), nil
		}
		v, err := strconv.ParseFloat(c.Text, 64)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindBytes:
		if c.Null {
			return (*[]byte)(nil), nil
		}
		v, err := base64.StdEncoding.DecodeString(c.Text)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindTimestamp:
		if c.Null {
			return (*time.Time)(nil), nil
		}
		v, err := time.Parse(time.RFC3339Nano, c.Text)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindDuration:
		if c.Null {
			return (*time.Duration)(nil), nil
		}
		v, err := strconv.ParseInt(c.Text, 10, 64)
		if err != nil {
			return nil, err
		}
		d := time.Duration(v)
		return &d, nil
	case kindLongDuration:
		if c.Null {
			return (*gocql.Duration)(nil), nil
		}
		v := gocql.Duration{Months: c.Months, Days: c.Days, Nanoseconds: c.Nanos}
		return &v, nil
	case kindUUID:
		if c.Null {
			return (*gocql.UUID)(nil), nil
		}
		v, err := gocql.ParseUUID(c.Text)
		if err != nil {
			return nil, err
		}
		return &v, nil
	case kindDecimal:
		if c.Null {
			return (**inf.Dec)(nil), nil
		}
		v, ok := new(inf.Dec).SetString(c.Text)
		if !ok {
			return nil, fmt.Errorf("testcontainer: %q is not a decimal", c.Text)
		}
		return &v, nil
	case kindVarint:
		if c.Null {
			return (**big.Int)(nil), nil
		}
		v, ok := new(big.Int).SetString(c.Text, 10)
		if !ok {
			return nil, fmt.Errorf("testcontainer: %q is not an integer", c.Text)
		}
		return &v, nil
	case kindJSON:
		decoder := json.NewDecoder(strings.NewReader(c.Text))
		decoder.UseNumber()
		var v any
		if err := decoder.Decode(&v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("testcontainer: unknown cell kind %q", c.Kind)
	}
}

// tapeType is a column's CQL type as a tape holds it. The driver reads only the
// type itself, so the wire marshalling a live TypeInfo also offers stays
// unimplemented here.
type tapeType int

func (t tapeType) Type() gocql.Type { return gocql.Type(t) }

func (tapeType) Zero() any { return nil }

func (tapeType) Marshal(any) ([]byte, error) {
	return nil, errors.New("testcontainer: a tape holds no encoded CQL value")
}

func (tapeType) Unmarshal([]byte, any) error {
	return errors.New("testcontainer: a tape holds no encoded CQL value")
}

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	decoder := json.NewDecoder(strings.NewReader(string(b)))
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// Reset empties the pinned server of every user keyspace and runs statements,
// so a tape holds one test's corpus whatever the test binary ran before. It
// does nothing in a replay, whose tape already holds the answers the corpus
// produced when it was recorded.
func Reset(t testing.TB, statements ...string) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	session, err := liveSession("")
	if err != nil {
		t.Fatalf("start Cassandra: %v", err)
	}
	ctx := context.Background()
	iter := session.Query("SELECT keyspace_name FROM system_schema.keyspaces").Iter()
	var name string
	var names []string
	for iter.Scan(&name) {
		names = append(names, name)
	}
	if err := iter.Close(); err != nil {
		t.Fatalf("list keyspaces: %v", err)
	}
	for _, name := range names {
		if strings.HasPrefix(name, "system") {
			continue
		}
		if err := session.Query("DROP KEYSPACE " + name).ExecContext(ctx); err != nil {
			t.Fatalf("drop keyspace %s: %v", name, err)
		}
	}
	for _, statement := range statements {
		if err := session.Query(statement).ExecContext(ctx); err != nil {
			t.Fatalf("%s:\n%v", statement, err)
		}
	}
}

// --- the live server ---

const (
	// startupTimeout is how long Cassandra's first boot may take: a JVM
	// starting, migrating its system keyspaces and listening for CQL.
	startupTimeout = 5 * time.Minute
	connectTimeout = 30 * time.Second
)

// sessions holds one session per keyspace a recording asked for: a session
// created with a keyspace resolves unqualified names there, which is what the
// driver's session does.
var sessions = struct {
	sync.Mutex
	byKeyspace map[string]*gocql.Session
	started    bool
	startErr   error
	container  testcontainers.Container
	address    string
	tag        string
}{byKeyspace: map[string]*gocql.Session{}}

// liveSession returns the pinned server's session in keyspace, starting the
// container on the first call. Only a recording calls it.
func liveSession(keyspace string) (*gocql.Session, error) {
	sessions.Lock()
	defer sessions.Unlock()
	if session, ok := sessions.byKeyspace[keyspace]; ok {
		return session, nil
	}
	if !sessions.started {
		sessions.started = true
		sessions.startErr = startContainer()
	}
	if sessions.startErr != nil {
		return nil, sessions.startErr
	}
	session, err := connect(keyspace)
	if err != nil {
		return nil, err
	}
	sessions.byKeyspace[keyspace] = session
	return session, nil
}

// startContainer starts the pinned server and holds its address for the
// sessions a recording opens.
func startContainer() (err error) {
	ctx := context.Background()
	pin, pinErr := enginetest.LoadPin(Pins, Driver)
	if pinErr != nil {
		return pinErr
	}
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	container, err := testcontainers.Run(ctx, pin.Image,
		testcontainers.WithEnv(map[string]string{
			// A recording needs one node's answers, not a cluster's, and the
			// heap the machine can spare.
			"CASSANDRA_CLUSTER_NAME": "omni",
			"MAX_HEAP_SIZE":          "1G",
			"HEAP_NEWSIZE":           "200M",
		}),
		testcontainers.WithExposedPorts("9042/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Starting listening for CQL clients").WithStartupTimeout(startupTimeout)),
	)
	if err != nil {
		return err
	}
	if err := checkDigest(ctx, container, pin.Digest); err != nil {
		return err
	}
	port, err := container.MappedPort(ctx, "9042/tcp")
	if err != nil {
		return err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return err
	}
	sessions.container = container
	sessions.address = host + ":" + port.Port()
	sessions.tag = pin.Tag
	return nil
}

// connect opens a session on the pinned server in keyspace.
func connect(keyspace string) (*gocql.Session, error) {
	ctx := context.Background()
	cluster := gocql.NewCluster(sessions.address)
	cluster.ConnectTimeout = connectTimeout
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("Cassandra does not answer: %w", err)
	}
	var version string
	if err := session.Query("SELECT release_version FROM system.local").ScanContext(ctx, &version); err != nil {
		return nil, fmt.Errorf("Cassandra does not answer: %w", err)
	}
	if version != sessions.tag {
		return nil, fmt.Errorf("the server reports %q, but pins.tsv pins Cassandra %s", version, sessions.tag)
	}
	return session, nil
}

func checkDigest(ctx context.Context, container testcontainers.Container, digest string) error {
	info, err := container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Cassandra container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Cassandra image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Cassandra container runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// Terminate closes the sessions this run started and removes the server, so a
// recording leaves no container behind even when the reaper is disabled.
func Terminate() {
	sessions.Lock()
	defer sessions.Unlock()
	for _, session := range sessions.byKeyspace {
		session.Close()
	}
	if sessions.container != nil {
		_ = testcontainers.TerminateContainer(sessions.container)
	}
}

// --- the live client's answers ---

// liveDo answers one request from the pinned server, in the encoding a tape
// holds.
func liveDo(t testing.TB, text, keyspace string) (string, error) {
	t.Helper()
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	session, err := liveSession(keyspace)
	if err != nil {
		t.Fatalf("start Cassandra: %v", err)
	}
	ctx := context.Background()
	var a answer
	switch req.Op {
	case opExec:
		if err := session.Query(req.Statement).ExecContext(ctx); err != nil {
			if e := encodeError(err); e != nil {
				a.Error = e
				break
			}
			return "", err
		}
	case opQuery:
		pageState, err := base64.StdEncoding.DecodeString(req.PageState)
		if err != nil {
			return "", fmt.Errorf("testcontainer: unreadable page state in %q: %w", text, err)
		}
		columns, rows, next, err := liveQuery(ctx, session, req.Statement, req.PageSize, pageState)
		if err != nil {
			if e := encodeError(err); e != nil {
				a.Error = e
				break
			}
			return "", err
		}
		for _, info := range columns {
			a.Columns = append(a.Columns, column{Name: info.Name, Type: int(info.TypeInfo.Type())})
		}
		if len(rows) > 0 {
			a.RowColumns = rows[0].Columns
		}
		for _, row := range rows {
			values := make([]cell, 0, len(row.Values))
			for _, value := range row.Values {
				encoded, err := encodeCell(value)
				if err != nil {
					return "", err
				}
				values = append(values, encoded)
			}
			a.Rows = append(a.Rows, values)
		}
		a.PageState = base64.StdEncoding.EncodeToString(next)
	default:
		return "", fmt.Errorf("testcontainer: unknown operation %q", req.Op)
	}
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// liveQuery is the page read the driver makes of gocql, against the pinned
// server: one page, its columns, and the state that continues after it.
func liveQuery(ctx context.Context, session *gocql.Session, statement string, pageSize int, pageState []byte) ([]gocql.ColumnInfo, []gocql.RowData, []byte, error) {
	iter := session.Query(statement).PageSize(pageSize).PageState(pageState).IterContext(ctx)
	defer func() { _ = iter.Close() }()
	nextPageState := iter.PageState()
	var rows []gocql.RowData
	for {
		rowData, err := iter.RowData()
		if err != nil {
			return nil, nil, nil, err
		}
		if !iter.Scan(rowData.Values...) {
			break
		}
		rows = append(rows, rowData)
	}
	if err := iter.Close(); err != nil {
		return nil, nil, nil, err
	}
	return iter.Columns(), rows, nextPageState, nil
}
