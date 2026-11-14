// Package testcontainer is how googlesql's tests reach the pinned Spanner
// emulator: through googlesql/spanner-oracle, a harness that sends one
// statement to the emulator and answers a JSON verdict for it
// (googlesql/spanner-oracle/README.md). A test asks through a tape
// (enginetest/tape) that replays the harness's answers, or, in a recording
// (OMNI_RECORD=1), through the harness and the emulator themselves. The
// emulator is one testcontainer per test binary, started by the first test
// that records and removed by Terminate, which each package's TestMain
// calls.
//
// The tape is the harness's: each request is one statement, and each answer
// the harness's JSON for it. The harness is a program of its own module, so
// the tape is the boundary the spec gives such programs: replay builds and
// runs nothing.
package testcontainer

import (
	"bufio"
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is googlesql/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the one emulator the tests ask.
const Default = "default"

var server = &tape.Server{Pins: Pins, Live: live, Codec: noErrors{}}

// A Verdict is the harness's answer for one statement
// (googlesql/spanner-oracle/README.md).
type Verdict struct {
	// Verdict is accept, reject, or error when the harness could not decide.
	Verdict string `json:"verdict"`
	// Kind is how the harness sent the statement: query, dml or ddl.
	Kind string `json:"kind"`
	// Reason is ok, syntax or semantic; for an error, infra or empty.
	Reason string `json:"reason"`
	// Code is the emulator's gRPC status code, and Message its message.
	Code    string `json:"code"`
	Message string `json:"message"`
}

// An Oracle asks the emulator, through the harness, whether it accepts a
// statement.
type Oracle struct {
	ask func(ctx context.Context, sql string) (string, error)
}

// Open returns t's Oracle: a replay of its tape, or, in a recording, a
// harness of its own on the pinned emulator, with a database no other
// stream has touched. Each call is a stream of its own.
func Open(t testing.TB) *Oracle {
	t.Helper()
	db := tape.Open(t, server, Default)
	return &Oracle{ask: func(ctx context.Context, stmt string) (string, error) {
		var answer string
		err := db.QueryRowContext(ctx, stmt).Scan(&answer)
		return answer, err
	}}
}

// Judge returns the pinned emulator itself, outside any tape, for the judge
// of compat.tsv, which runs only in a recording.
func Judge(t testing.TB) *Oracle {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatalf("%v", err)
	}
	h, err := startHarness(started(t, pin))
	if err != nil {
		t.Fatalf("start the Spanner harness: %v", err)
	}
	t.Cleanup(h.close)
	return &Oracle{ask: func(_ context.Context, stmt string) (string, error) { return h.ask(stmt) }}
}

// Verdict asks the emulator about one statement. An error is no answer
// from the harness at all.
func (o *Oracle) Verdict(ctx context.Context, stmt string) (Verdict, error) {
	answer, err := o.ask(ctx, stmt)
	if err != nil {
		return Verdict{}, err
	}
	var v Verdict
	if err := json.Unmarshal([]byte(answer), &v); err != nil {
		return Verdict{}, fmt.Errorf("decode the harness's answer %q: %w", answer, err)
	}
	return v, nil
}

// MustVerdict is Verdict, failing t when the harness gives no answer, or
// answers error, which means it could not decide (an emulator fault, a
// timeout): never an acceptance or a rejection.
func (o *Oracle) MustVerdict(t testing.TB, stmt string) Verdict {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	v, err := o.Verdict(ctx, stmt)
	if err != nil {
		t.Fatalf("ask the emulator about %q: %v", stmt, err)
	}
	if v.Verdict != "accept" && v.Verdict != "reject" {
		t.Fatalf("the emulator gave no verdict for %q: %s %s (%s): %s", stmt, v.Verdict, v.Reason, v.Code, v.Message)
	}
	return v
}

// shared is the one pinned emulator this run started, and the harness
// binary built for it.
var shared struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	host      string
	binDir    string
	bin       string
	err       error
}

// started returns the address of the pinned emulator, with the harness
// built, starting both on the first call and failing t unless the emulator
// runs pin.
func started(t testing.TB, pin enginetest.Pin) string {
	t.Helper()
	shared.once.Do(func() {
		shared.pin = pin
		shared.err = start(context.Background(), pin)
	})
	if shared.err != nil {
		t.Fatalf("start the Spanner emulator: %v", shared.err)
	}
	if shared.pin != pin {
		t.Fatalf("this run started the emulator at %s, not %s", shared.pin, pin)
	}
	return shared.host
}

func start(ctx context.Context, pin enginetest.Pin) (err error) {
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	if err := buildHarness(); err != nil {
		return err
	}
	image, _, _ := strings.Cut(pin.Image, ":")
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        image + "@" + pin.Digest,
			ExposedPorts: []string{"9010/tcp"},
			WaitingFor:   wait.ForListeningPort("9010/tcp").WithStartupTimeout(2 * time.Minute),
		},
		Started: true,
	})
	if c != nil {
		shared.container = c
	}
	if err != nil {
		return err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return err
	}
	port, err := c.MappedPort(ctx, "9010/tcp")
	if err != nil {
		return err
	}
	shared.host = host + ":" + port.Port()
	return nil
}

// harnessDir is googlesql/spanner-oracle.
var harnessDir = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "spanner-oracle")
}()

// buildHarness builds the harness into a temporary directory of this test
// binary's own, so packages recording at once never write or run the same
// file.
func buildHarness() error {
	dir, err := os.MkdirTemp("", "googlesql-spanner-harness-")
	if err != nil {
		return err
	}
	shared.binDir = dir
	bin := filepath.Join(dir, "googlesql-spanner")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = harnessDir
	if out, err := build.CombinedOutput(); err != nil {
		return fmt.Errorf("build %s: %w\n%s", harnessDir, err, out)
	}
	shared.bin = bin
	return nil
}

// Terminate removes the emulator and the harness binary this run started
// and built, if any. A package's TestMain calls it.
func Terminate() {
	if shared.container != nil {
		_ = testcontainers.TerminateContainer(shared.container)
	}
	if shared.binDir != "" {
		_ = os.RemoveAll(shared.binDir)
	}
}

// A harness is one running googlesql/spanner-oracle process, which made a
// database of its own on the emulator when it started.
type harness struct {
	mu  sync.Mutex
	cmd *exec.Cmd
	in  io.WriteCloser
	out *bufio.Reader
}

func startHarness(host string) (*harness, error) {
	cmd := exec.Command(shared.bin)
	cmd.Env = append(os.Environ(), "SPANNER_EMULATOR_HOST="+host, "GOOGLESQL_HARNESS_LINE=1")
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &harness{cmd: cmd, in: in, out: bufio.NewReader(out)}, nil
}

// ask sends one statement to the harness and returns its JSON line.
func (h *harness) ask(stmt string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := fmt.Fprintln(h.in, base64.StdEncoding.EncodeToString([]byte(stmt))); err != nil {
		return "", fmt.Errorf("write to the Spanner harness: %w", err)
	}
	line, err := h.out.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read from the Spanner harness: %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// close ends the harness, which drops its database.
func (h *harness) close() {
	_ = h.in.Close()
	_ = h.cmd.Wait()
}

// live connects to the pinned emulator through a harness of the stream's
// own.
func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	return connector{host: started(t, pin)}
}

// connector is a database/sql connector whose queries the harness answers,
// so enginetest/tape can record it: a query is one statement, and its
// answer one row of one column, the harness's JSON.
type connector struct{ host string }

func (c connector) Connect(context.Context) (driver.Conn, error) {
	h, err := startHarness(c.host)
	if err != nil {
		return nil, err
	}
	return conn{h}, nil
}

func (connector) Driver() driver.Driver { return harnessDriver{} }

type harnessDriver struct{}

func (harnessDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("spanner harness: open a connection with testcontainer.Open")
}

type conn struct{ h *harness }

func (conn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("spanner harness: no prepared statements")
}

func (c conn) Close() error {
	c.h.close()
	return nil
}

func (conn) Begin() (driver.Tx, error) { return nil, errors.New("spanner harness: no transactions") }

func (c conn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, errors.New("spanner harness: a query takes no arguments")
	}
	answer, err := c.h.ask(query)
	if err != nil {
		return nil, err
	}
	return &answerRows{answer: answer}, nil
}

// answerRows is one row of one column, answer.
type answerRows struct {
	answer string
	done   bool
}

func (*answerRows) Columns() []string                     { return []string{"verdict"} }
func (*answerRows) Close() error                          { return nil }
func (*answerRows) ColumnTypeDatabaseTypeName(int) string { return "JSON" }
func (r *answerRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	dest[0] = r.answer
	return nil
}

// noErrors is the codec of a connection no answer of which is an error:
// the harness reports the emulator's errors inside its JSON.
type noErrors struct{}

func (noErrors) EncodeError(error) (string, bool) { return "", false }
func (noErrors) DecodeError(text string) error {
	return fmt.Errorf("tape: the Spanner harness answered no error, but the tape holds %q", text)
}
