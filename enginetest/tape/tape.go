// Package tape records what a hand-written test sends to an engine and the
// engine's answers, and replays them so the test runs with no server.
//
// A test calls its engine's testcontainer.Open(t), which calls Open here
// and gets an ordinary *sql.DB. In a replay, the default, that connection
// answers from testdata/tapes/<top-level test>.tape. With OMNI_RECORD=1 it
// talks to the pinned server and writes the tape. With OMNI_RECORD=verify it
// talks to the pinned server, fails on any difference from the committed tape,
// and writes or deletes nothing on disk.
//
// An engine whose client library is not database/sql calls OpenRaw instead and
// writes the exchange itself: a request and its answer are opaque text, which
// the engine's testcontainer spells so that one request always answers the same
// text. Replay holds a test to
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
	// called only in a recording (OMNI_RECORD=1) or verification
	// (OMNI_RECORD=verify). The connector must inline arguments, since a tape
	// refuses driver-level prepared statements.
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
// its tape, or, in a recording (OMNI_RECORD=1), the pinned server with every
// exchange written down, or, in a verification (OMNI_RECORD=verify), the pinned
// server compared against the committed tape with nothing written to disk. Each
// call is one stream, named by t, role and how many connections t opened before;
// the pool holds one connection, so session state (USE, SET SESSION) holds as
// it would on one live connection. A connection the engine refuses is recorded
// as a connect exchange, and a replay answers the pool's use of it with that
// error. The package's TestMain must run tests through Main.
func Open(t testing.TB, s *Server, role string) *sql.DB {
	t.Helper()
	ses, f, pin, mode := begin(t, s, role)

	c := &connector{s: ses}
	if mode != enginetest.ModeReplay {
		c.live = s.Live(t, pin)
	}

	db := sql.OpenDB(c)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = db.Close()
		end(t, f, ses, role, pin, mode)
	})
	return db
}

// begin opens t's stream in pin row role: a replay of its tape, a recording
// stream that will be written down, or a verification stream compared against
// the committed tape without writing. Each stream is named by t, role and how
// many streams t opened before.
func begin(t testing.TB, s *Server, role string) (*session, *fileState, enginetest.Pin, enginetest.Mode) {
	t.Helper()
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("tape: %v", err)
	}
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
	// merge conflict markers; replay and verify stop on any parse or read error.
	if f.err != nil && !(mode == enginetest.ModeRecord && f.existed) {
		t.Fatalf("tape: %v", f.err)
	}
	ses := &session{t: t, name: name, codec: s.Codec, mode: mode}
	if mode == enginetest.ModeRecord {
		ses.poison = f.poison
	}
	if r, ok := s.Codec.(Renamer); ok && mode != enginetest.ModeReplay {
		ses.renamer, ses.names = r, f.names
	}
	var committed *stream
	if f.old != nil && f.old.pins[role] == pin.String() {
		committed = f.old.stream(name)
	}
	if committed != nil {
		ses.want = committed.exchanges
	}

	if mode != enginetest.ModeRecord {
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
	return ses, f, pin, mode
}

// end accounts for one stream once its test is done: a replay fails on
// answers no request asked for, a verification checks every live request and
// answer against the committed stream without writing, and a recording replaces
// the file it opened.
func end(t testing.TB, f *fileState, ses *session, role string, pin enginetest.Pin, mode enginetest.Mode) {
	t.Helper()
	switch mode {
	case enginetest.ModeReplay:
		ses.finish()
		return
	case enginetest.ModeVerify:
		ses.verify()
		return
	}
	ses.reproduces()
	reg.Lock()
	defer reg.Unlock()
	if f.poisoned {
		return // the test failed, and the committed tape stays
	}
	f.pins[role] = pin.String()
	f.recorded[ses.name] = &stream{name: ses.name, exchanges: ses.got}
	if err := f.write(); err != nil {
		t.Errorf("tape: %v", err)
	}
}

// OpenRaw returns t's stream in pin row role. live, when the run records or
// verifies, answers a request from the pinned server; a replay never calls it
// and may pass nil.
func OpenRaw(t testing.TB, s *Server, role string, live func(req string) (string, error)) *Raw {
	t.Helper()
	ses, f, pin, mode := begin(t, s, role)
	r := &Raw{ses: ses}
	if mode != enginetest.ModeReplay {
		if live == nil {
			t.Fatalf("tape: OpenRaw needs a live function in %s mode", mode)
		}
		r.live = live
	}
	t.Cleanup(func() { end(t, f, ses, role, pin, mode) })
	return r
}

// Raw is one stream of opaque exchanges, for an engine whose client library is
// not database/sql: a request is the text the driver asked its client, and the
// answer is the text the client returned. Neither is interpreted here; each
// engine's testcontainer writes them so that the same request always answers
// the same text.
type Raw struct {
	ses  *session
	live func(req string) (string, error)
}

// Do sends one request and returns its answer. In a recording it asks the
// pinned server and writes the answer down; in a replay it returns the
// recorded answer, failing the test when the tape holds another request next.
// A returned error means the answer could not be produced at all, which a
// tape cannot hold: the request or the server is not the engine's answer.
//
// A Server whose Codec is a Renamer has the values the engine made up replaced
// in the answer and in the requests that repeat them, as a database/sql
// stream does.
func (r *Raw) Do(text string) (string, error) {
	req := request{kind: "raw", text: text}
	if r.live == nil {
		a, err := r.ses.replay(req)
		if err != nil {
			return "", err
		}
		return a.text, nil
	}
	ans, err := r.live(r.ses.show(text))
	if err != nil {
		return "", err
	}
	ans = r.ses.hide(ans)
	r.ses.record(req, answer{kind: "raw", text: ans})
	return ans, nil
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

// Live skips t outside a recording (OMNI_RECORD=1), including in replay and
// verification (OMNI_RECORD=verify). It is for the few tests with nothing to
// replay and no committed tape to verify: a test whose input is the engine's
// own data that the repo may not hold, as Oracle's PL/SQL dictionary test reads
// the source Oracle ships with. In a recording it runs against the live server
// and records nothing. It is the only environment gate (§4.5 of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md).
func Live(t testing.TB) {
	t.Helper()
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatalf("tape: %v", err)
	}
	if mode != enginetest.ModeRecord {
		t.Skipf("tape: %s has nothing to replay and runs only against the pinned server:\n  %s", t.Name(), enginetest.RecordCommand(t))
	}
}

// Main runs the package's tests, as TestMain's m.Run does, and returns the
// exit code. A passing run of every test also accounts for every tape: a
// replay or verification fails on a tape or stream no test opened, which a
// deleted or renamed test leaves behind, and only a recording deletes such
// tapes.
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
	mode, err := enginetest.RunMode()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if problems := leftovers(mode); len(problems) > 0 {
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
// recording deletes the tapes. In replay and verification it deletes nothing.
func leftovers(mode enginetest.Mode) []string {
	reg.Lock()
	defer reg.Unlock()
	paths, _ := filepath.Glob(filepath.Join(dir, "*.tape"))
	var problems []string
	for _, p := range paths {
		f, ok := reg.files[p]
		if !ok {
			if mode == enginetest.ModeRecord {
				if err := os.Remove(p); err != nil {
					problems = append(problems, err.Error())
				}
				continue
			}
			problems = append(problems, fmt.Sprintf("%s belongs to no test that ran; delete it if its test is gone", p))
			continue
		}
		if mode == enginetest.ModeRecord || f.old == nil {
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
