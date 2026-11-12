package enginetest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A source is one row of an engine's testdata/sources.tsv: a vendor file
// the engine's tests read and the repo does not hold.
type source struct {
	name   string // the file's name
	sha256 string // the checksum of its content, in hex
	url    string // where the pinned release publishes it
}

// cacheDir is where fetched files are kept: per user, so every checkout and
// worktree shares one download.
var cacheDir = func() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "omni-enginetest")
}

// fetchAttempts and fetchPause bound a download. Tests shorten the pause.
var (
	fetchAttempts = 3
	fetchPause    = 2 * time.Second
)

// Fetched returns the local path of a vendor file the repo does not hold
// because its license is not the repo's, such as a grammar under the GPL or
// Mulan PubL. sources is the engine's testdata/sources.tsv, one
// tab-separated row per file: name, SHA-256, URL. The file is downloaded on
// first use into the user's cache directory and its checksum is checked on
// every use, so the content is as pinned as a committed file would be.
//
// A file that cannot be had fails t: these tests need the network once per
// machine, as a container test needs its image.
func Fetched(t testing.TB, sources, name string) string {
	t.Helper()
	src, err := loadSource(sources, name)
	if err != nil {
		t.Fatalf("%v", err)
	}
	path := filepath.Join(cacheDir(), src.sha256[:16]+"-"+filepath.Base(src.name))
	if b, err := os.ReadFile(path); err == nil && checksum(b) == src.sha256 {
		return path
	}
	b, err := download(src.url)
	if err != nil {
		t.Fatalf("fetch %s (%s): %v\nThis package reads it once per machine. To run without the network, put the file with SHA-256 %s at %s",
			src.name, src.url, err, src.sha256, path)
	}
	if got := checksum(b); got != src.sha256 {
		t.Fatalf("%s holds a file with SHA-256 %s, but %s pins %s", src.url, got, sources, src.sha256)
	}
	if err := writeAtomic(path, b); err != nil {
		t.Fatalf("cache %s: %v", src.name, err)
	}
	return path
}

func loadSource(path, name string) (source, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return source{}, err
	}
	for i, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 3 {
			return source{}, fmt.Errorf("%s:%d: want 3 tab-separated fields (name, SHA-256, URL), got %d", path, i+1, len(f))
		}
		if f[0] != name {
			continue
		}
		if sum, err := hex.DecodeString(f[1]); err != nil || len(sum) != sha256.Size {
			return source{}, fmt.Errorf("%s:%d: %q is not a SHA-256 in hex", path, i+1, f[1])
		}
		return source{name: f[0], sha256: f[1], url: f[2]}, nil
	}
	return source{}, fmt.Errorf("%s has no %q row", path, name)
}

func checksum(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func download(url string) ([]byte, error) {
	client := &http.Client{Timeout: time.Minute}
	var last error
	for i := 0; i < fetchAttempts; i++ {
		if i > 0 {
			time.Sleep(fetchPause)
		}
		b, err := get(client, url)
		if err == nil {
			return b, nil
		}
		last = err
	}
	return nil, last
}

func get(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the server answered %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// writeAtomic writes the file under a temporary name and renames it into
// place, so a test of another package never reads half of it.
func writeAtomic(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(b); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
