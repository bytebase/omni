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
	"io"
	"os"
	"path/filepath"
	"regexp"
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

// A Renamer is the Codec of an engine that makes up names, such as OceanBase's
// for a foreign key written without one: a different name on every run. A
// recording keeps a stand-in for each, in the engine's answers and in the
// requests that repeat the name, and hands the test the stand-in too. So a
// test reads the same text in a recording as in a replay, and a second
// recording writes the same tape.
type Renamer interface {
	Codec
	// MadeUp matches each name the engine makes up.
	MadeUp() *regexp.Regexp
	// StandIn returns what stands for the nth distinct made-up name of a
	// test, counting from 1. It must match MadeUp and be a name the engine
	// would never make up itself.
	StandIn(name string, n int) string
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
	existed  bool      // the file was there when this run first opened it
	orig     []byte    // its bytes then
	old      *tapeFile // nil when there is no committed tape, or it does not parse
	err      error
	poisoned bool // a recording took an answer that can't replay; the committed file stays
	opened   map[string]bool
	pins     map[string]string
	recorded map[string]*stream
	names    *names // the stand-ins this run's recording gave out
}

func file(path string) *fileState {
	f := reg.files[path]
	if f != nil {
		return f
	}
	f = &fileState{path: path, opened: map[string]bool{}, pins: map[string]string{}, recorded: map[string]*stream{}, names: newNames()}
	reg.files[path] = f
	b, err := os.ReadFile(path)
	switch {
	case os.IsNotExist(err):
	case err != nil:
		f.err = err
	default:
		f.existed, f.orig = true, b
		f.old, f.err = parse(path, string(b))
	}
	return f
}

// poison ends the file's recording, because an answer in it can't be
// replayed. A stream that finished earlier may have rewritten the file by
// now, so poison puts the committed bytes back, and the streams that finish
// later leave them alone.
func (f *fileState) poison() error {
	reg.Lock()
	defer reg.Unlock()
	if f.poisoned {
		return nil
	}
	f.poisoned = true
	if f.existed {
		return os.WriteFile(f.path, f.orig, 0o644)
	}
	if err := os.Remove(f.path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
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

	// A recording replaces a tape that does not parse, such as one left with
	// merge conflict markers; anything else wrong with the file stops here.
	if f.err != nil && !(recording && f.existed) {
		t.Fatalf("tape: %v", f.err)
	}
	ses := &session{t: t, name: name, codec: s.Codec, rec: recording, poison: f.poison}
	if r, ok := s.Codec.(Renamer); ok && recording {
		ses.renamer, ses.names = r, f.names
	}
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
		if f.poisoned {
			return // the test failed, and the committed tape stays
		}
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
// exit code. A passing run of every test also accounts for every tape: a
// replay fails on a tape or stream no test opened, which a deleted or
// renamed test leaves behind, and a recording deletes such tapes.
func Main(m *testing.M) int {
	reg.Lock()
	reg.main = true
	reg.Unlock()
	return account(m.Run(), partialRun(testFlag), os.Stderr)
}

// account settles the tapes no test opened once the tests have run. Only a
// passing run of every test can: a test that failed may not have opened its
// streams, and a partial run says nothing about the tapes of the tests it
// left out.
func account(code int, partial bool, stderr io.Writer) int {
	if code != 0 || partial {
		return code
	}
	recording, err := enginetest.RecordingMode()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if problems := leftovers(recording); len(problems) > 0 {
		fmt.Fprintf(stderr, "tape: %s\n", strings.Join(problems, "\ntape: "))
		return 1
	}
	return code
}

func testFlag(name string) string {
	if f := flag.Lookup(name); f != nil {
		return f.Value.String()
	}
	return ""
}

// partialRun reports whether this run leaves tests out: -run and -skip
// select some, and -list and a fuzz worker process run none at all.
func partialRun(flagValue func(name string) string) bool {
	for _, name := range []string{"test.run", "test.skip", "test.list"} {
		if flagValue(name) != "" {
			return true
		}
	}
	return flagValue("test.fuzzworker") == "true"
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
