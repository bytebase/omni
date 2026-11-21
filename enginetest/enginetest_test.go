package enginetest

import (
	"os"
	"path/filepath"
	"strings"
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
	if got, want := p.ImageRef(), "img:1@sha256:abc"; got != want {
		t.Errorf("ImageRef = %q, want %q", got, want)
	}
	if _, err := LoadPin(path, "nosuch"); err == nil {
		t.Errorf("a missing role loaded")
	}
	if err := os.WriteFile(path, []byte("default\timg:1\tsha256:abc\tv1\t2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err = LoadPin(path, "default")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.String(), "default img:1 sha256:abc v1 judge:2"; got != want {
		t.Errorf("with a judge revision: got %q, want %q", got, want)
	}
	if err := os.WriteFile(path, []byte("default\timg:1\tsha256:abc\tv1\t\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPin(path, "default"); err == nil {
		t.Errorf("an empty judge revision loaded")
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
	if UseLive(t) {
		t.Errorf("unset uses live")
	}
	if Verifying(t) {
		t.Errorf("unset verifies")
	}
	if m, err := RunMode(); err != nil || m != ModeReplay || m.String() != "replay" {
		t.Errorf("unset: RunMode() = (%v, %v), want (replay, nil)", m, err)
	}

	t.Setenv(RecordEnv, "1")
	if !Recording(t) {
		t.Errorf("1 does not record")
	}
	if !UseLive(t) {
		t.Errorf("1 does not use live")
	}
	if Verifying(t) {
		t.Errorf("1 verifies")
	}
	if m, err := RunMode(); err != nil || m != ModeRecord || m.String() != "record" {
		t.Errorf("1: RunMode() = (%v, %v), want (record, nil)", m, err)
	}

	t.Setenv(RecordEnv, "verify")
	if !Recording(t) {
		t.Errorf("verify does not report Recording (live server access)")
	}
	if !UseLive(t) {
		t.Errorf("verify does not report UseLive")
	}
	if !Verifying(t) {
		t.Errorf("verify does not report Verifying")
	}
	if m, err := RunMode(); err != nil || m != ModeVerify || m.String() != "verify" {
		t.Errorf("verify: RunMode() = (%v, %v), want (verify, nil)", m, err)
	}
}

func TestRunModeCIAndUnknownValues(t *testing.T) {
	t.Setenv("CI", "true")

	t.Setenv(RecordEnv, "")
	if m, err := RunMode(); err != nil || m != ModeReplay {
		t.Errorf("unset with CI: got (%v, %v), want (replay, nil)", m, err)
	}

	t.Setenv(RecordEnv, "1")
	if _, err := RunMode(); err == nil || !strings.Contains(err.Error(), "CI") {
		t.Errorf("1 with CI: got %v, want an error naming CI", err)
	}
	if _, err := RecordingMode(); err == nil || !strings.Contains(err.Error(), "CI") {
		t.Errorf("RecordingMode with 1 and CI: got %v, want an error naming CI", err)
	}

	// OMNI_RECORD=verify is allowed in CI so self-hosted verify lanes can check committed answers.
	t.Setenv(RecordEnv, "verify")
	if m, err := RunMode(); err != nil || m != ModeVerify {
		t.Errorf("verify with CI: got (%v, %v), want (verify, nil)", m, err)
	}
	if on, err := RecordingMode(); err != nil || !on {
		t.Errorf("RecordingMode with verify and CI: got (%v, %v), want (true, nil)", on, err)
	}
	if !Verifying(t) || !UseLive(t) {
		t.Errorf("verify with CI: Verifying=%v UseLive=%v, want true true", Verifying(t), UseLive(t))
	}

	t.Setenv("CI", "")
	for _, bad := range []string{"0", "true", "2", "VERIFY"} {
		t.Setenv(RecordEnv, bad)
		_, err := RunMode()
		if err == nil || !strings.Contains(err.Error(), "1") || !strings.Contains(err.Error(), "verify") {
			t.Errorf("OMNI_RECORD=%q: got %v, want an error naming both 1 and verify", bad, err)
		}
	}
}
