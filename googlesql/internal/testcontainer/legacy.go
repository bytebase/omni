package testcontainer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// The tests that still ask the emulator live reach it through Host and
// HarnessBinary, until they replay tapes through Open.

// Image is the emulator build the recorded verdicts came from
// (sha256:caf1bd24…). Newer builds drift on a few forms, so it stays pinned.
const Image = "gcr.io/cloud-spanner-emulator/emulator:1.5.54"

var legacy struct {
	once sync.Once
	host string
	err  error
}

// Host returns the emulator address: $SPANNER_EMULATOR_HOST when set, else an
// emulator testcontainer started once per test binary and exported through
// that variable so the harness subprocesses find it too. Docker is assumed,
// so a start failure fails the test.
func Host(t testing.TB) string {
	t.Helper()
	legacy.once.Do(func() {
		if h := os.Getenv("SPANNER_EMULATOR_HOST"); h != "" {
			legacy.host = h
			return
		}
		ctx := context.Background()
		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        Image,
				ExposedPorts: []string{"9010/tcp"},
				WaitingFor:   wait.ForListeningPort("9010/tcp").WithStartupTimeout(2 * time.Minute),
			},
			Started: true,
		})
		if err != nil {
			legacy.err = fmt.Errorf("start Spanner emulator: %w", err)
			return
		}
		host, err := c.Host(ctx)
		if err != nil {
			legacy.err = err
			return
		}
		port, err := c.MappedPort(ctx, "9010/tcp")
		if err != nil {
			legacy.err = err
			return
		}
		legacy.host = fmt.Sprintf("%s:%s", host, port.Port())
		_ = os.Setenv("SPANNER_EMULATOR_HOST", legacy.host)
	})
	if legacy.err != nil {
		t.Fatalf("spanner emulator: %v", legacy.err)
	}
	return legacy.host
}

var legacyHarness struct {
	once    sync.Once
	dir     string // temp dir holding the binary; "" until built
	bin     string
	missing string // harness project path when it is absent (skip)
	err     error
}

// HarnessBinary returns the path of a googlesql/spanner-oracle CLI built once
// per test binary into a private temp dir, so parallel `go test -p=N` packages
// never write (or exec) the same file — building into the module dir raced and
// surfaced as "fork/exec ...: text file busy". The dir lives until
// RemoveHarness, which the owning package's TestMain calls. Skips when the
// harness project is not checked out.
func HarnessBinary(t testing.TB) string {
	t.Helper()
	legacyHarness.once.Do(func() {
		_, thisFile, _, _ := runtime.Caller(0)
		// googlesql/internal/testcontainer/container.go → the project is googlesql/spanner-oracle
		projDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "spanner-oracle"))
		if _, err := os.Stat(projDir); err != nil {
			legacyHarness.missing = projDir
			return
		}
		dir, err := os.MkdirTemp("", "googlesql-spanner-harness-")
		if err != nil {
			legacyHarness.err = fmt.Errorf("temp dir for harness: %w", err)
			return
		}
		legacyHarness.dir = dir
		bin := filepath.Join(dir, "googlesql-spanner")
		build := exec.Command("go", "build", "-o", bin, ".")
		build.Dir = projDir
		if out, err := build.CombinedOutput(); err != nil {
			legacyHarness.err = fmt.Errorf("building harness: %w\n%s", err, out)
			return
		}
		legacyHarness.bin = bin
	})
	if legacyHarness.missing != "" {
		t.Skipf("harness project not found at %s", legacyHarness.missing)
	}
	if legacyHarness.err != nil {
		t.Fatal(legacyHarness.err)
	}
	return legacyHarness.bin
}

// RemoveHarness deletes the temp dir HarnessBinary built into. Call it from
// TestMain after m.Run(); it is a no-op when nothing was built.
func RemoveHarness() {
	if legacyHarness.dir != "" {
		_ = os.RemoveAll(legacyHarness.dir)
	}
}
