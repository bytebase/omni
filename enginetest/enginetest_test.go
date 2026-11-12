package enginetest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pins.tsv")
	text := "# role\timage\tdigest\ttag\ndefault\timg:1\tsha256:abc\tv1\nlower\timg:1\tsha256:abc\tv1\n"
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := LoadPin(path, "default")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.String(), "default img:1 sha256:abc v1"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if _, err := LoadPin(path, "nosuch"); err == nil {
		t.Errorf("a missing role loaded")
	}
	if err := os.WriteFile(path, []byte("default\timg:1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPin(path, "default"); err == nil {
		t.Errorf("a short row loaded")
	}
}

func TestRecording(t *testing.T) {
	t.Setenv("CI", "")
	t.Setenv(RecordEnv, "")
	if Recording(t) {
		t.Errorf("unset records")
	}
	t.Setenv(RecordEnv, "1")
	if !Recording(t) {
		t.Errorf("1 does not record")
	}
}
