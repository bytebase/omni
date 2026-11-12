// Package enginetest holds what every engine's compatibility tests share: the
// pinned servers an engine's answers come from, and the switch between a
// replay, which reads committed answers, and a recording, which asks the
// judge and writes them. It knows no engine.
//
// The design is docs/specs/2026-11-11-engine-compatibility-without-servers-design.md.
package enginetest

import (
	"fmt"
	"os"
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
}

// String renders the pin as committed files name it in their header.
func (p Pin) String() string {
	return strings.Join([]string{p.Role, p.Image, p.Digest, p.Tag}, " ")
}

// LoadPin reads the row for role from a pins.tsv. The file holds one
// tab-separated row per role: role, image, digest, source tag.
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
		if len(f) != 4 {
			return Pin{}, fmt.Errorf("%s:%d: want 4 tab-separated fields (role, image, digest, source tag), got %d", path, i+1, len(f))
		}
		if f[0] == role {
			return Pin{Role: f[0], Image: f[1], Digest: f[2], Tag: f[3]}, nil
		}
	}
	return Pin{}, fmt.Errorf("%s has no %q row", path, role)
}

// RecordEnv is the environment variable that turns a run into a recording.
const RecordEnv = "OMNI_RECORD"

// Recording reports whether this run records answers: OMNI_RECORD=1. It fails
// t on any other non-empty value, and when CI is set, because CI never
// writes answers.
func Recording(t testing.TB) bool {
	t.Helper()
	switch v := os.Getenv(RecordEnv); v {
	case "":
		return false
	case "1":
		if os.Getenv("CI") != "" {
			t.Fatalf("%s=1 with CI set: CI replays committed answers and never records them", RecordEnv)
		}
		return true
	default:
		t.Fatalf("%s=%q: want 1 to record, or unset to replay", RecordEnv, v)
		return false
	}
}
