package testcontainer

import (
	"bufio"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// ScriptDOM is the pin role of Microsoft's ScriptDOM parser, which
// mssql/scriptdom runs in the pinned .NET SDK image. The row's source tag is
// the ScriptDOM package version mssql/scriptdom/harness.csproj builds
// against.
const ScriptDOM = "scriptdom"

var scriptDOMServer = &tape.Server{Pins: Pins, Live: liveScriptDOM, Codec: noErrors{}}

// OpenScriptDOM returns t's connection to ScriptDOM: a replay of its tape,
// or, in a recording, mssql/scriptdom's harness. A query is one SQL text,
// and its answer one row of one column, shape: the harness's JSON for it
// (mssql/scriptdom/README.md).
func OpenScriptDOM(t testing.TB) *sql.DB {
	t.Helper()
	return tape.Open(t, scriptDOMServer, ScriptDOM)
}

// ScriptDOMShape asks db, an OpenScriptDOM connection, for the harness's
// JSON for sql.
func ScriptDOMShape(ctx context.Context, db *sql.DB, sql string) (string, error) {
	var shape string
	err := db.QueryRowContext(ctx, sql).Scan(&shape)
	return shape, err
}

// harnessDir is mssql/scriptdom.
var harnessDir = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "scriptdom")
}()

// harness is the one ScriptDOM process a recording starts.
var harness struct {
	once sync.Once
	mu   sync.Mutex
	pin  enginetest.Pin
	cmd  *exec.Cmd
	name string
	in   io.WriteCloser
	out  *bufio.Reader
	err  error
}

// packageVersion reads the ScriptDOM version harness.csproj references.
var packageVersion = regexp.MustCompile(`Include="Microsoft\.SqlServer\.TransactSql\.ScriptDom" Version="([^"]+)"`)

// liveScriptDOM starts the harness on first use: it builds and runs inside
// the pinned SDK image, named by its digest, so the only host requirement is
// Docker.
func liveScriptDOM(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	harness.once.Do(func() { harness.err = startHarness(pin) })
	if harness.err != nil {
		t.Fatalf("start the ScriptDOM harness: %v", harness.err)
	}
	if harness.pin != pin {
		t.Fatalf("this run started ScriptDOM at %s, not %s", harness.pin, pin)
	}
	return harnessConnector{}
}

func startHarness(pin enginetest.Pin) error {
	csproj, err := os.ReadFile(filepath.Join(harnessDir, "harness.csproj"))
	if err != nil {
		return err
	}
	m := packageVersion.FindSubmatch(csproj)
	if m == nil || string(m[1]) != pin.Tag {
		return fmt.Errorf("mssql/scriptdom/harness.csproj does not reference ScriptDOM %s, which pins.tsv pins", pin.Tag)
	}
	image, _, _ := strings.Cut(pin.Image, ":")
	// --user keeps bin/, obj/ and the NuGet cache owned by the invoking user;
	// the cache persists across runs, per uid, since runners and developers
	// share a host's /tmp. The SDK's first-run NuGet mutex needs a HOME it
	// can write as that uid.
	cache := filepath.Join(os.TempDir(), fmt.Sprintf("omni-scriptdom-nuget-%d", os.Getuid()))
	if err := os.MkdirAll(cache, 0o755); err != nil {
		return err
	}
	harness.pin = pin
	harness.name = fmt.Sprintf("omni-scriptdom-%d", os.Getpid())
	cmd := exec.Command("docker", "run", "-i", "--rm", "--name", harness.name,
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-v", harnessDir+":/src", "-v", cache+":/nuget", "-w", "/src",
		"-e", "HOME=/tmp", "-e", "DOTNET_CLI_HOME=/tmp", "-e", "NUGET_PACKAGES=/nuget",
		"-e", "DOTNET_CLI_TELEMETRY_OPTOUT=1", "-e", "DOTNET_NOLOGO=1", "-e", "MSSQL_HARNESS_LINE=1",
		image+"@"+pin.Digest, "sh", "-c",
		"dotnet build -c Release -v quiet --nologo 1>&2 && exec dotnet bin/Release/net8.0/mssql-scriptdom-harness.dll")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	harness.cmd, harness.in, harness.out = cmd, in, bufio.NewReader(out)
	return nil
}

// stopHarness ends the harness process and removes its container.
func stopHarness() {
	if harness.cmd == nil {
		return
	}
	harness.in.Close()
	_ = harness.cmd.Wait()
	_ = exec.Command("docker", "rm", "-f", harness.name).Run()
}

// ask sends one SQL text to the harness and returns its JSON line.
func ask(sql string) (string, error) {
	harness.mu.Lock()
	defer harness.mu.Unlock()
	if _, err := fmt.Fprintln(harness.in, base64.StdEncoding.EncodeToString([]byte(sql))); err != nil {
		return "", fmt.Errorf("write to the ScriptDOM harness: %w", err)
	}
	line, err := harness.out.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read from the ScriptDOM harness: %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// harnessConnector is a database/sql connector whose queries the harness
// answers.
type harnessConnector struct{}

func (harnessConnector) Connect(context.Context) (driver.Conn, error) { return harnessConn{}, nil }
func (harnessConnector) Driver() driver.Driver                        { return harnessDriver{} }

type harnessDriver struct{}

func (harnessDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("ScriptDOM: open a connection with testcontainer.OpenScriptDOM")
}

type harnessConn struct{}

func (harnessConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("ScriptDOM: no prepared statements")
}
func (harnessConn) Close() error              { return nil }
func (harnessConn) Begin() (driver.Tx, error) { return nil, errors.New("ScriptDOM: no transactions") }

func (harnessConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, errors.New("ScriptDOM: a query takes no arguments")
	}
	shape, err := ask(query)
	if err != nil {
		return nil, err
	}
	return &shapeRows{shape: shape}, nil
}

// shapeRows is one row of one column, shape.
type shapeRows struct {
	shape string
	done  bool
}

func (*shapeRows) Columns() []string                     { return []string{"shape"} }
func (*shapeRows) Close() error                          { return nil }
func (*shapeRows) ColumnTypeDatabaseTypeName(int) string { return "JSON" }
func (r *shapeRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	dest[0] = r.shape
	return nil
}

// noErrors is the codec of a connection no answer of which is an error:
// the harness reports a parse error inside its JSON.
type noErrors struct{}

func (noErrors) EncodeError(error) (string, bool) { return "", false }
func (noErrors) DecodeError(text string) error {
	return fmt.Errorf("tape: ScriptDOM answered no error, but the tape holds %q", text)
}
