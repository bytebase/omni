// Package enginetest holds what every engine's compatibility tests share: the
// pinned servers an engine's answers come from, and the switch among a replay
// (reading committed answers), a recording (asking the judge and writing them),
// and a verification (asking the judge again without writing). It knows no
// engine.
//
// The design is docs/specs/2026-11-11-engine-compatibility-without-servers-design.md.
package enginetest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A Pin is one row of an engine's testdata/pins.tsv: the exact server
// release its answers come from.
type Pin struct {
	Role   string // "default" judges the rows; other roles serve tests that need another server
	Image  string // the image reference tests start, e.g. pingcap/tidb:v8.5.5
	Digest string // the image's repository digest, sha256:...
	Tag    string // the source tag the release and its grammar were built from
	// Judge is the revision of how the engine's tests read the server's
	// answers, "" for the first. The reading is part of what an answer
	// means: a harness that comes to read an error as a rejection where it
	// read an acceptance answers differently at the same server, so a new
	// reading is a pin bump, like a new image, and its recording may move
	// expected cells.
	Judge string
}

// ImageRef is the reference a testcontainer starts: the image and the
// pinned digest. A recording pulls this reference rather than the tag, so a
// registry moving the tag cannot select another image. A testcontainer still
// checks the running container's digest (checkDigest) because the reference
// says what Docker was asked for, not what it ran.
func (p Pin) ImageRef() string { return p.Image + "@" + p.Digest }

// String renders the pin as committed files name it in their header.
func (p Pin) String() string {
	s := strings.Join([]string{p.Role, p.Image, p.Digest, p.Tag}, " ")
	if p.Judge != "" {
		s += " judge:" + p.Judge
	}
	return s
}

// LoadPin reads the row for role from a pins.tsv. The file holds one
// tab-separated row per role: role, image, digest, source tag, and, once
// the engine's reading of the server has changed, the judge revision.
func LoadPin(path, role string) (Pin, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Pin{}, err
	}
	for i, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 4 && len(f) != 5 {
			return Pin{}, fmt.Errorf("%s:%d: want 4 or 5 tab-separated fields (role, image, digest, source tag, judge revision), got %d", path, i+1, len(f))
		}
		if f[0] == role {
			p := Pin{Role: f[0], Image: f[1], Digest: f[2], Tag: f[3]}
			if len(f) == 5 {
				if f[4] == "" {
					return Pin{}, fmt.Errorf("%s:%d: an empty judge revision; leave the field out for the first", path, i+1)
				}
				p.Judge = f[4]
			}
			return p, nil
		}
	}
	return Pin{}, fmt.Errorf("%s has no %q row", path, role)
}

// FilesDigest returns the digest a pin holds for vendor files committed
// as published: sha256: and the SHA-256 of their bytes, one file after
// another, which is what sha256sum prints for cat of the same files.
func FilesDigest(paths ...string) (string, error) {
	h := sha256.New()
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return "", err
		}
		h.Write(b)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

// RecordEnv is the environment variable that selects a live run mode: "1" to
// record answers, "verify" to verify committed answers against the pinned
// server without writing, or unset to replay committed answers.
const RecordEnv = "OMNI_RECORD"

// Mode is how a test run handles committed engine answers.
type Mode int

const (
	// ModeReplay reads committed answers and reaches no server (OMNI_RECORD unset).
	ModeReplay Mode = iota
	// ModeRecord asks the pinned server or judge and writes committed answers (OMNI_RECORD=1).
	ModeRecord
	// ModeVerify asks the pinned server or judge again, fails on any difference
	// from committed answers, and writes or deletes nothing on disk (OMNI_RECORD=verify).
	ModeVerify
)

// String names m for diagnostics.
func (m Mode) String() string {
	switch m {
	case ModeReplay:
		return "replay"
	case ModeRecord:
		return "record"
	case ModeVerify:
		return "verify"
	default:
		return fmt.Sprintf("Mode(%d)", int(m))
	}
}

// RunMode returns the run mode selected by OMNI_RECORD: ModeReplay when unset,
// ModeRecord for "1", and ModeVerify for "verify". It fails on any other
// non-empty value, and on OMNI_RECORD=1 when CI is set because CI never writes
// answers. OMNI_RECORD=verify is allowed when CI is set so a self-hosted verify
// lane can check committed answers against pinned servers without writing.
func RunMode() (Mode, error) {
	switch v := os.Getenv(RecordEnv); v {
	case "":
		return ModeReplay, nil
	case "1":
		if os.Getenv("CI") != "" {
			return ModeReplay, fmt.Errorf("%s=1 with CI set: CI replays committed answers and never records them", RecordEnv)
		}
		return ModeRecord, nil
	case "verify":
		return ModeVerify, nil
	default:
		return ModeReplay, fmt.Errorf("%s=%q: want 1 to record, verify to verify against the pinned server, or unset to replay", RecordEnv, v)
	}
}

// Recording reports whether this run talks to the pinned server or judge rather
// than replaying committed answers: OMNI_RECORD=1 (record) or OMNI_RECORD=verify
// (verify). Callers that prepare or query the live server (such as testcontainer
// Reset and Judge helpers) use this so both recording and verification run
// against the pinned server. Callers that write or delete committed files must
// check RunMode or Verifying so OMNI_RECORD=verify stays strictly read-only.
// It fails t on an unknown OMNI_RECORD value, and on OMNI_RECORD=1 when CI is set.
func Recording(t testing.TB) bool {
	t.Helper()
	on, err := RecordingMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return on
}

// RecordingMode is Recording for code that has no test, such as TestMain: it
// returns true for both OMNI_RECORD=1 and OMNI_RECORD=verify. Callers that
// write or delete files must call RunMode instead.
func RecordingMode() (bool, error) {
	m, err := RunMode()
	if err != nil {
		return false, err
	}
	return m != ModeReplay, nil
}

// UseLive reports whether this run talks to the pinned server or judge
// (OMNI_RECORD=1 or OMNI_RECORD=verify).
func UseLive(t testing.TB) bool {
	t.Helper()
	return Recording(t)
}

// Verifying reports whether this run is in read-only verification mode
// (OMNI_RECORD=verify). It fails t on an unknown OMNI_RECORD value, and on
// OMNI_RECORD=1 when CI is set.
func Verifying(t testing.TB) bool {
	t.Helper()
	m, err := RunMode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return m == ModeVerify
}

// RunCommand is the go test command that reruns t's top-level test from the
// repo root, for messages that tell the reader what to run.
func RunCommand(t testing.TB) string {
	return fmt.Sprintf("go test %s -run '^%s$'", packageArg(), TopName(t.Name()))
}

// RecordCommand is RunCommand as a recording.
func RecordCommand(t testing.TB) string { return RecordEnv + "=1 " + RunCommand(t) }

// TopName is the top-level test of a test name: the part before any subtest.
func TopName(name string) string {
	top, _, _ := strings.Cut(name, "/")
	return top
}

// packageArg names the package under test as a go test argument from the
// repo root, such as ./tidb/parser/. Tests run in their package directory.
func packageArg() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			rel, err := filepath.Rel(dir, wd)
			if err != nil {
				return "."
			}
			return "./" + filepath.ToSlash(rel) + "/"
		}
		if filepath.Dir(dir) == dir {
			return "."
		}
	}
}
