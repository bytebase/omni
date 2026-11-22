package testcontainer

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

// A Verdict is the googlesql/spanner-oracle harness's answer for one statement
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

// An Oracle asks the pinned Spanner emulator, through googlesql/spanner-oracle,
// whether it accepts a statement.
type Oracle struct {
	ask func(ctx context.Context, sql string) (string, error)
}

// Verdict asks the emulator about one statement. An error is no answer from
// the harness at all.
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

// Judge returns the pinned emulator itself, outside any tape, for the judge of
// spanner/parser/testdata/compat.tsv, which runs only in a recording.
func Judge(t testing.TB) *Oracle {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	bin, err := builtHarness()
	if err != nil {
		t.Fatalf("build the Spanner harness: %v", err)
	}
	host := Host(t)
	h, err := startHarness(bin, host)
	if err != nil {
		t.Fatalf("start the Spanner harness: %v", err)
	}
	t.Cleanup(h.close)
	return &Oracle{ask: func(_ context.Context, stmt string) (string, error) { return h.ask(stmt) }}
}

var harnessState struct {
	mu     sync.Mutex
	binDir string
	bin    string
	err    error
	built  bool
}

// harnessDir is googlesql/spanner-oracle.
var harnessDir = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "googlesql", "spanner-oracle")
}()

func builtHarness() (string, error) {
	harnessState.mu.Lock()
	defer harnessState.mu.Unlock()
	if harnessState.built {
		return harnessState.bin, harnessState.err
	}
	harnessState.built = true
	dir, err := os.MkdirTemp("", "spanner-oracle-harness-")
	if err != nil {
		harnessState.err = err
		return "", err
	}
	harnessState.binDir = dir
	bin := filepath.Join(dir, "spanner-oracle")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = harnessDir
	if out, err := build.CombinedOutput(); err != nil {
		harnessState.err = fmt.Errorf("build %s: %w\n%s", harnessDir, err, out)
		return "", harnessState.err
	}
	harnessState.bin = bin
	return bin, nil
}

func cleanupHarness() {
	harnessState.mu.Lock()
	dir := harnessState.binDir
	harnessState.binDir = ""
	harnessState.bin = ""
	harnessState.err = nil
	harnessState.built = false
	harnessState.mu.Unlock()
	if dir != "" {
		_ = os.RemoveAll(dir)
	}
}

type harness struct {
	mu  sync.Mutex
	cmd *exec.Cmd
	in  io.WriteCloser
	out *bufio.Reader
}

func startHarness(bin, host string) (*harness, error) {
	cmd := exec.Command(bin)
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

func (h *harness) close() {
	_ = h.in.Close()
	_ = h.cmd.Wait()
}
