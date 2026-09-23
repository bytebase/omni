package spannertest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

var harness struct {
	once    sync.Once
	dir     string // temp dir holding the binary; "" until built
	bin     string
	missing string // harness project path when it is absent (skip)
	err     error
}

// HarnessBinary returns the path of a harness/googlesql-spanner CLI built once
// per test binary into a private temp dir, so parallel `go test -p=N` packages
// never write (or exec) the same file — building into the module dir raced and
// surfaced as "fork/exec ...: text file busy". The dir lives until
// RemoveHarness, which the owning package's TestMain calls. Skips when the
// harness project is not checked out.
func HarnessBinary(t testing.TB) string {
	t.Helper()
	harness.once.Do(func() {
		_, thisFile, _, _ := runtime.Caller(0)
		// googlesql/internal/spannertest/harness.go → repo root is ../../..
		repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", ".."))
		projDir := filepath.Join(repoRoot, "harness", "googlesql-spanner")
		if _, err := os.Stat(projDir); err != nil {
			harness.missing = projDir
			return
		}
		dir, err := os.MkdirTemp("", "googlesql-spanner-harness-")
		if err != nil {
			harness.err = fmt.Errorf("temp dir for harness: %w", err)
			return
		}
		harness.dir = dir
		bin := filepath.Join(dir, "googlesql-spanner")
		build := exec.Command("go", "build", "-o", bin, ".")
		build.Dir = projDir
		if out, err := build.CombinedOutput(); err != nil {
			harness.err = fmt.Errorf("building harness: %w\n%s", err, out)
			return
		}
		harness.bin = bin
	})
	if harness.missing != "" {
		t.Skipf("harness project not found at %s", harness.missing)
	}
	if harness.err != nil {
		t.Fatal(harness.err)
	}
	return harness.bin
}

// RemoveHarness deletes the temp dir HarnessBinary built into. Call it from
// TestMain after m.Run(); it is a no-op when nothing was built.
func RemoveHarness() {
	if harness.dir != "" {
		_ = os.RemoveAll(harness.dir)
	}
}
