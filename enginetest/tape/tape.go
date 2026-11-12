// Package tape records what a hand-written test sends to an engine and the
// engine's answers, and replays them so the test runs with no server.
//
// A test calls its engine's testcontainer.Open(t), which calls Open here
// and gets an ordinary *sql.DB. In a replay, the default, that connection
// answers from testdata/tapes/<top-level test>.tape. With OMNI_RECORD=1 it
// talks to the pinned server and writes the tape. Replay holds a test to
// its tape in both directions: a request the tape does not hold next fails,
// and so do recorded requests the test never sends. So what a test sends
// must depend only on its inputs, never on omni's results.
//
// The design is §4 and Appendix C of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md.
package tape

import (
	"database/sql"
	"database/sql/driver"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

// A Codec turns the driver's errors into tape text and back.
type Codec interface {
	// EncodeError renders an error the engine answered, and false for any
	// other error (a timeout, a driver fault), which a recording refuses.
	EncodeError(err error) (string, bool)
	// DecodeError rebuilds the driver's own error type from its text, so a
	// test's errors.As still finds it in a replay.
	DecodeError(text string) error
}

// A Server is how an engine's tests reach one pinned server.
type Server struct {
	// Pins is the engine's pins.tsv. An absolute path: tests of every
	// package of the engine use it.
	Pins string
	// Live connects to the pinned server, starting it on first use. It is
	// called only in a recording. The connector must inline arguments, since
	// a tape refuses driver-level prepared statements.
	Live func(t testing.TB, pin enginetest.Pin) driver.Connector
	// Codec renders the driver's errors.
	Codec Codec
}

// dir holds a package's tapes, relative to the package directory a test
// runs in.
var dir = filepath.Join("testdata", "tapes")

var reg = struct {
	sync.Mutex
	main   bool
	files  map[string]*fileState
	counts map[openKey]int
}{files: map[string]*fileState{}, counts: map[openKey]int{}}

type openKey struct {
	t    testing.TB
	role string
}

// fileState is one tape file during a run: what is committed, and what a
// recording has written so far.
type fileState struct {
	path     string
	old      *tapeFile // nil when there is no committed tape
	err      error
	opened   map[string]bool
	pins     map[string]string
	recorded map[string]*stream
}

func file(path string) *fileState {
	f := reg.files[path]
	if f != nil {
		return f
	}
	f = &fileState{path: path, opened: map[string]bool{}, pins: map[string]string{}, recorded: map[string]*stream{}}
	reg.files[path] = f
	b, err := os.ReadFile(path)
	switch {
	case os.IsNotExist(err):
	case err != nil:
		f.err = err
	default:
		f.old, f.err = parse(path, string(b))
	}
	return f
}

// Open returns t's connection to the server in pin row role: a replay of
// its tape, or, in a recording, the pinned server with every exchange
// written down. Each call is one stream, named by t, role and how many
// connections t opened before; the pool holds one connection, so session
// state (USE, SET SESSION) holds as it would on one live connection. The
// package's TestMain must run tests through Main.
func Open(t testing.TB, s *Server, role string) *sql.DB {
	t.Helper()
	recording := enginetest.Recording(t)
	pin, err := enginetest.LoadPin(s.Pins, role)
	if err != nil {
		t.Fatalf("tape: %v", err)
	}

	reg.Lock()
	if !reg.main {
		reg.Unlock()
		t.Fatalf("tape: tape.Open needs tape.Main: call os.Exit(tape.Main(m)) from the package's TestMain")
	}
	f := file(filepath.Join(dir, enginetest.TopName(t.Name())+".tape"))
	key := openKey{t, role}
	reg.counts[key]++
	name := fmt.Sprintf("%s %s#%d", t.Name(), role, reg.counts[key])
	f.opened[name] = true
	reg.Unlock()

	if f.err != nil {
		t.Fatalf("tape: %v", f.err)
	}
	ses := &session{t: t, name: name, codec: s.Codec, rec: recording}
	var committed *stream
	if f.old != nil && f.old.pins[role] == pin.String() {
		committed = f.old.stream(name)
	}
	if committed != nil {
		ses.want = committed.exchanges
	}

	c := &connector{s: ses}
	if recording {
		c.live = s.Live(t, pin)
	} else {
		switch {
		case f.old == nil:
			t.Fatalf("tape: no tape %s. Record it:\n  %s", f.path, recordCommand(t))
		case f.old.pins[role] != pin.String():
			t.Fatalf("tape: %s was recorded at pin %q, but %s pins %q. Record it:\n  %s",
				f.path, f.old.pins[role], s.Pins, pin, recordCommand(t))
		case committed == nil:
			t.Fatalf("tape: %s has no stream %q. Record it:\n  %s", f.path, name, recordCommand(t))
		}
	}

	db := sql.OpenDB(c)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = db.Close()
		if !recording {
			ses.finish()
			return
		}
		ses.reproduces()
		reg.Lock()
		defer reg.Unlock()
		f.pins[role] = pin.String()
		f.recorded[name] = &stream{name: name, exchanges: ses.got}
		if err := f.write(); err != nil {
			t.Errorf("tape: %v", err)
		}
	})
	return db
}

// write saves the streams this run recorded, in name order. A recording
// replaces the whole file, so record whole top-level tests.
func (f *fileState) write() error {
	out := &tapeFile{pins: f.pins}
	names := make([]string, 0, len(f.recorded))
	for n := range f.recorded {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		out.streams = append(out.streams, f.recorded[n])
	}
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(f.path, []byte(out.encode()), 0o644)
}

// Main runs the package's tests, as TestMain's m.Run does, and returns the
// exit code. A run with no -run or -skip also accounts for every tape:
// a replay fails on a tape or stream no test opened, which a deleted or
// renamed test leaves behind, and a recording deletes such tapes.
func Main(m *testing.M) int {
	reg.Lock()
	reg.main = true
	reg.Unlock()
	code := m.Run()
	if filtered() {
		return code
	}
	recording, err := enginetest.RecordingMode()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if code != 0 && !recording {
		return code // a test that failed may not have opened its streams
	}
	problems := leftovers(recording)
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "tape: %s\n", strings.Join(problems, "\ntape: "))
		return 1
	}
	return code
}

func filtered() bool {
	for _, name := range []string{"test.run", "test.skip"} {
		if f := flag.Lookup(name); f != nil && f.Value.String() != "" {
			return true
		}
	}
	return false
}

// leftovers reports the tapes and streams no test opened, and in a
// recording deletes the tapes.
func leftovers(recording bool) []string {
	reg.Lock()
	defer reg.Unlock()
	paths, _ := filepath.Glob(filepath.Join(dir, "*.tape"))
	var problems []string
	for _, p := range paths {
		f, ok := reg.files[p]
		if !ok {
			if recording {
				if err := os.Remove(p); err != nil {
					problems = append(problems, err.Error())
				}
				continue
			}
			problems = append(problems, fmt.Sprintf("%s belongs to no test that ran; delete it if its test is gone", p))
			continue
		}
		if recording || f.old == nil {
			continue
		}
		for _, s := range f.old.streams {
			if !f.opened[s.name] {
				problems = append(problems, fmt.Sprintf("%s: no test opened stream %q; record its test again", p, s.name))
			}
		}
	}
	return problems
}

func recordCommand(t testing.TB) string { return enginetest.RecordCommand(t) }
