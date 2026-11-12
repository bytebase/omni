package enginetest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// failing records the failure of a call that is expected to fail its test.
type failing struct {
	testing.TB
	msg string
}

func (f *failing) Helper() {}
func (f *failing) Fatalf(format string, args ...any) {
	f.msg = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

// fetchFails runs Fetched expecting it to fail, and returns its message.
func fetchFails(t *testing.T, sources, name string) string {
	t.Helper()
	f := &failing{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		Fetched(f, sources, name)
	}()
	<-done
	if f.msg == "" {
		t.Fatalf("Fetched(%s) did not fail", name)
	}
	return f.msg
}

func TestFetched(t *testing.T) {
	const content = "%token A\n%%\nstmt: A\n%%\n"
	served := content
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/grammar.y" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, served)
	}))
	defer srv.Close()

	cache := t.TempDir()
	oldDir, oldPause := cacheDir, fetchPause
	cacheDir, fetchPause = func() string { return cache }, 0
	defer func() { cacheDir, fetchPause = oldDir, oldPause }()

	sources := filepath.Join(t.TempDir(), "sources.tsv")
	write := func(rows string) {
		t.Helper()
		if err := os.WriteFile(sources, []byte("# name\tsha256\turl\n"+rows), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	sum := checksum([]byte(content))
	write("grammar.y\t" + sum + "\t" + srv.URL + "/grammar.y\n" +
		"gone.y\t" + sum + "\t" + srv.URL + "/gone.y\n" +
		"short.y\tabc\t" + srv.URL + "/grammar.y\n")

	// The first use downloads the file; the second reads the cache.
	path := Fetched(t, sources, "grammar.y")
	if b, err := os.ReadFile(path); err != nil || string(b) != content {
		t.Fatalf("the fetched file holds %q, %v", b, err)
	}
	if again := Fetched(t, sources, "grammar.y"); again != path || requests != 1 {
		t.Errorf("a second use made %d requests and answered %s, want the cached %s", requests, again, path)
	}

	// A cached file that no longer matches is fetched again.
	if err := os.WriteFile(path, []byte("edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(Fetched(t, sources, "grammar.y")); string(b) != content || requests != 2 {
		t.Errorf("an edited cache file was kept: %q after %d requests", b, requests)
	}

	// What cannot be had, or is not what the row pins, fails the test.
	for name, want := range map[string]string{
		"nosuch.y": `has no "nosuch.y" row`,
		"short.y":  "is not a SHA-256",
		"gone.y":   "404 Not Found\nThis package reads it once per machine. To run without the network, put the file with SHA-256 " + sum,
	} {
		if msg := fetchFails(t, sources, name); !strings.Contains(msg, want) {
			t.Errorf("%s: got %q, want %q", name, msg, want)
		}
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	served = "another release"
	if msg := fetchFails(t, sources, "grammar.y"); !strings.Contains(msg, "holds a file with SHA-256 "+checksum([]byte(served))+", but") {
		t.Errorf("a changed file: got %q", msg)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("a file that failed its checksum was cached")
	}
}
