package tape

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

func TestMain(m *testing.M) { os.Exit(Main(m)) }

// fakeT stands in for a test whose failures the harness's tests check.
type fakeT struct {
	testing.TB
	name     string
	errs     []string
	cleanups []func()
}

func (f *fakeT) Name() string      { return f.name }
func (f *fakeT) Helper()           {}
func (f *fakeT) Cleanup(fn func()) { f.cleanups = append(f.cleanups, fn) }
func (f *fakeT) Errorf(format string, args ...any) {
	f.errs = append(f.errs, fmt.Sprintf(format, args...))
}
func (f *fakeT) Fatalf(format string, args ...any) {
	f.Errorf(format, args...)
	runtime.Goexit()
}

// run runs body as a test named name and returns what it logged and the
// failures the harness reported.
func run(t *testing.T, name string, body func(t testing.TB) []string) (log, errs []string) {
	t.Helper()
	f := &fakeT{TB: t, name: name}
	done := make(chan struct{})
	go func() {
		defer close(done)
		log = body(f)
	}()
	<-done
	for i := len(f.cleanups) - 1; i >= 0; i-- {
		f.cleanups[i]()
	}
	return log, f.errs
}

// harness points the package at a temporary tape directory and pins file,
// and forgets what earlier tests opened, as a new process would.
func harness(t *testing.T) (*Server, *fakeEngine) {
	t.Helper()
	old := dir
	dir = t.TempDir()
	t.Cleanup(func() { dir = old })
	resetRegistry()
	pins := filepath.Join(t.TempDir(), "pins.tsv")
	if err := os.WriteFile(pins, []byte("default\timg:1\tsha256:1\t1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	e := &fakeEngine{}
	return &Server{
		Pins:  pins,
		Live:  func(testing.TB, enginetest.Pin) driver.Connector { return fakeConnector{e} },
		Codec: fakeCodec{},
	}, e
}

func resetRegistry() {
	reg.Lock()
	defer reg.Unlock()
	reg.files = map[string]*fileState{}
	reg.counts = map[openKey]int{}
}

func record(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
}

func replay(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "")
	resetRegistry()
}

// session is the body most tests run: every kind of request and answer.
func everything(srv *Server) func(t testing.TB) []string {
	return func(t testing.TB) []string {
		db := Open(t, srv, "default")
		var log []string
		res, err := db.Exec("CREATE TABLE t (a INT)\n-- a second line", 1, "two", nil, []byte("x"))
		if err != nil {
			return append(log, "exec: "+err.Error())
		}
		n, _ := res.RowsAffected()
		log = append(log, fmt.Sprint("affected ", n))

		_, err = db.Exec("BAD")
		var fe *fakeError
		log = append(log, fmt.Sprintf("bad: %v, as fakeError: %v", err, errors.As(err, &fe) && fe.code == 1064))

		rows, err := db.Query("SELECT rows")
		if err != nil {
			return append(log, "query: "+err.Error())
		}
		for {
			types, _ := rows.ColumnTypes()
			for _, c := range types {
				log = append(log, "column "+c.Name()+" "+c.DatabaseTypeName())
			}
			for rows.Next() {
				vals := make([]any, len(types))
				ptrs := make([]any, len(types))
				for i := range vals {
					ptrs[i] = &vals[i]
				}
				if err := rows.Scan(ptrs...); err != nil {
					return append(log, "scan: "+err.Error())
				}
				log = append(log, fmt.Sprintf("row %q", vals))
			}
			if !rows.NextResultSet() {
				break
			}
		}
		rows.Close()

		tx, err := db.Begin()
		if err != nil {
			return append(log, "begin: "+err.Error())
		}
		if _, err := tx.Exec("INSERT INTO t VALUES (1)"); err != nil {
			return append(log, "tx exec: "+err.Error())
		}
		log = append(log, fmt.Sprint("commit: ", tx.Commit()))
		return log
	}
}

func TestRecordThenReplay(t *testing.T) {
	srv, _ := harness(t)
	record(t)
	recorded, errs := run(t, "TestFake", everything(srv))
	if len(errs) > 0 {
		t.Fatalf("recording failed: %v", errs)
	}
	b, err := os.ReadFile(filepath.Join(dir, "TestFake.tape"))
	if err != nil {
		t.Fatal(err)
	}
	tape := string(b)
	for _, want := range []string{
		"-- pin: default img:1 sha256:1 1\n",
		"== TestFake default#1\n> exec\n| CREATE TABLE t (a INT)\n| -- a second line\n: int64 1\n: string\n| two\n: null\n: bytes\n| x\n< ok 4 -\n",
		"< error\n| 1064 near \"BAD\"\n| line 2\n",
		"< row\n: int64 1\n: bytes\n| x\n| y\n",
		"|q \"\\r\"\n",
		"> begin\n< ok\n",
	} {
		if !strings.Contains(tape, want) {
			t.Errorf("the tape lacks %q:\n%s", want, tape)
		}
	}

	replay(t)
	srv.Live = func(t testing.TB, _ enginetest.Pin) driver.Connector {
		t.Fatalf("a replay reached the server")
		return nil
	}
	replayed, errs := run(t, "TestFake", everything(srv))
	if len(errs) > 0 {
		t.Fatalf("replay failed: %v", errs)
	}
	if strings.Join(replayed, "\n") != strings.Join(recorded, "\n") {
		t.Errorf("a replay differs from its recording:\nrecorded:\n%s\nreplayed:\n%s", strings.Join(recorded, "\n"), strings.Join(replayed, "\n"))
	}
	if !strings.Contains(strings.Join(recorded, "\n"), "as fakeError: true") {
		t.Errorf("the codec lost the error's type: %v", recorded)
	}
}

func TestReplayFailsOnAnUnrecordedRequest(t *testing.T) {
	srv, _ := harness(t)
	record(t)
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) > 0 {
		t.Fatal(errs)
	}
	replay(t)
	_, errs := run(t, "TestFake", func(t testing.TB) []string {
		db := Open(t, srv, "default")
		_, err := db.Exec("CREATE TABLE u (a INT)")
		_, err2 := db.Exec("CREATE TABLE t (a INT)") // the stream is broken now
		if err == nil || err2 == nil {
			return []string{"a request the tape doesn't hold succeeded"}
		}
		return nil
	})
	got := strings.Join(errs, "\n")
	if len(errs) != 1 || !strings.Contains(got, `no recorded answer for exec "CREATE TABLE u (a INT)" (request 1 of stream "TestFake default#1"); the tape holds exec "CREATE TABLE t`) ||
		!strings.Contains(got, "OMNI_RECORD=1 go test ./enginetest/tape/ -run '^TestFake$'") {
		t.Errorf("got %q", errs)
	}
}

func TestReplayFailsOnRequestsNeverSent(t *testing.T) {
	srv, _ := harness(t)
	record(t)
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) > 0 {
		t.Fatal(errs)
	}
	replay(t)
	_, errs := run(t, "TestFake", func(t testing.TB) []string {
		db := Open(t, srv, "default")
		_, _ = db.Exec("CREATE TABLE t (a INT)\n-- a second line", 1, "two", nil, []byte("x"))
		return nil
	})
	if len(errs) != 1 || !strings.Contains(errs[0], `holds 5 recorded requests the test never sent, starting with exec "BAD"`) {
		t.Errorf("got %q", errs)
	}
}

func TestReplayNeedsTheTapeAtThePin(t *testing.T) {
	srv, _ := harness(t)
	replay(t)
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) != 1 || !strings.Contains(errs[0], "no tape") {
		t.Errorf("with no tape: %q", errs)
	}

	record(t)
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) > 0 {
		t.Fatal(errs)
	}
	replay(t)
	if _, errs := run(t, "TestOther", func(t testing.TB) []string {
		Open(t, srv, "default")
		return nil
	}); len(errs) != 1 || !strings.Contains(errs[0], "no tape") {
		t.Errorf("another test: %q", errs)
	}
	if _, errs := run(t, "TestFake", func(t testing.TB) []string {
		Open(t, srv, "default")
		Open(t, srv, "default")
		return nil
	}); len(errs) == 0 || !strings.Contains(errs[0], `has no stream "TestFake default#2"`) {
		t.Errorf("a second connection: %q", errs)
	}

	resetRegistry()
	if err := os.WriteFile(srv.Pins, []byte("default\timg:2\tsha256:2\t2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) != 1 || !strings.Contains(errs[0], `was recorded at pin "default img:1 sha256:1 1"`) {
		t.Errorf("a new pin: %q", errs)
	}
}

func TestRecordingFailsOnANondeterministicAnswer(t *testing.T) {
	srv, _ := harness(t)
	now := func(t testing.TB) []string {
		db := Open(t, srv, "default")
		var n int
		if err := db.QueryRow("SELECT now").Scan(&n); err != nil {
			return []string{err.Error()}
		}
		return nil
	}
	record(t)
	if _, errs := run(t, "TestFake", now); len(errs) > 0 {
		t.Fatal(errs)
	}
	resetRegistry()
	_, errs := run(t, "TestFake", now)
	if len(errs) != 1 || !strings.Contains(errs[0], `answered request 1 (query "SELECT now") differently than the committed tape`) {
		t.Errorf("got %q", errs)
	}
}

// An engine that makes up names gives different ones on every run. The test
// and the tape see a stand-in for each, in answers and in errors, and the
// engine still gets its own name back wherever the test sends the stand-in.
func TestMadeUpNamesAreRecordedAsStandIns(t *testing.T) {
	srv, e := harness(t)
	srv.Codec = renamingCodec{}
	madeUp := renamingCodec{}.MadeUp()
	body := func(t testing.TB) []string {
		db := Open(t, srv, "default")
		var first, second string
		if err := db.QueryRow("SELECT made").Scan(&first); err != nil {
			return []string{err.Error()}
		}
		if err := db.QueryRow("SELECT made").Scan(&second); err != nil {
			return []string{err.Error()}
		}
		log := []string{first, second}
		_, err := db.Exec("DROP " + madeUp.FindString(first)) // in the request's text
		log = append(log, fmt.Sprint(err))
		_, err = db.Exec("DROP ?", madeUp.FindString(second)) // as an argument
		log = append(log, fmt.Sprint(err))
		_, err = db.Exec("DROP " + madeUp.FindString(first)) // gone by now, and the error names it
		return append(log, fmt.Sprint(err))
	}
	want := []string{"x made_1 y", "x made_2 y", "<nil>", "<nil>", "Error 1091: can't drop made_1"}
	path := filepath.Join(dir, "TestFake.tape")

	record(t)
	log, errs := run(t, "TestFake", body)
	if len(errs) > 0 || fmt.Sprint(log) != fmt.Sprint(want) {
		t.Fatalf("the recording saw %q (errors %q), want %q", log, errs, want)
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(first), "made_9") {
		t.Errorf("the tape holds a name the engine made up:\n%s", first)
	}

	// The engine has moved on, so a second recording meets other names.
	resetRegistry()
	if e.now != 2 {
		t.Fatalf("the engine made %d names, want 2", e.now)
	}
	log, errs = run(t, "TestFake", body)
	if len(errs) > 0 || fmt.Sprint(log) != fmt.Sprint(want) {
		t.Fatalf("the second recording saw %q (errors %q), want %q", log, errs, want)
	}
	if second, _ := os.ReadFile(path); string(second) != string(first) {
		t.Errorf("a second recording wrote another tape:\n%s\nfirst:\n%s", second, first)
	}

	replay(t)
	log, errs = run(t, "TestFake", body)
	if len(errs) > 0 || fmt.Sprint(log) != fmt.Sprint(want) {
		t.Errorf("the replay saw %q (errors %q), want %q", log, errs, want)
	}
}

// A codec may come to spell an error another way. A tape in the old spelling
// still replays, and recording its test again is not a different answer.
func TestACodecsOldSpellingIsTheSameAnswer(t *testing.T) {
	srv, _ := harness(t)
	bad := func(t testing.TB) []string {
		_, err := Open(t, srv, "default").Exec("BAD")
		return []string{fmt.Sprint(err)}
	}
	record(t)
	want, errs := run(t, "TestFake", bad)
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	path := filepath.Join(dir, "TestFake.tape")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	old := strings.Replace(string(b), "| 1064 near", "| error 1064 near", 1)
	if old == string(b) {
		t.Fatalf("the tape holds no error to respell:\n%s", b)
	}
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	resetRegistry()
	if _, errs := run(t, "TestFake", bad); len(errs) > 0 {
		t.Errorf("recording over the old spelling: %q", errs)
	}
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	replay(t)
	if got, errs := run(t, "TestFake", bad); len(errs) > 0 || got[0] != want[0] {
		t.Errorf("replaying the old spelling: %q %q, want %q", got, errs, want)
	}
}

func TestPreparedStatementsAreRefused(t *testing.T) {
	srv, _ := harness(t)
	record(t)
	_, errs := run(t, "TestFake", func(t testing.TB) []string {
		db := Open(t, srv, "default")
		var log []string
		if _, err := db.Prepare("SELECT 1"); !errors.Is(err, errPrepared) {
			log = append(log, fmt.Sprintf("Prepare: %v", err))
		}
		if _, err := db.Exec("SKIP", 1); !errors.Is(err, errPrepared) {
			log = append(log, fmt.Sprintf("a driver that would prepare: %v", err))
		}
		return log
	})
	if len(errs) > 0 {
		t.Error(errs)
	}
}

// An error the engine did not answer, such as a timeout, fails the recording
// and never reaches the tape, where a replay would hand it to the test as the
// engine's verdict.
func TestAnUnrecordableAnswerKeepsTheCommittedTape(t *testing.T) {
	srv, _ := harness(t)
	record(t)
	path := filepath.Join(dir, "TestFake.tape")
	read := func() string {
		t.Helper()
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	healthy := func(t testing.TB) []string {
		_, err := Open(t, srv, "default").Exec("CREATE TABLE t (a INT)")
		return []string{fmt.Sprint(err)}
	}
	// timeout opens a healthy stream beside the one that times out, and the
	// healthy one finishes last.
	timeout := func(t testing.TB) []string {
		_, _ = Open(t, srv, "default").Exec("CREATE TABLE t (a INT)")
		_, err := Open(t, srv, "default").Exec("TIMEOUT")
		return []string{fmt.Sprint(err)}
	}
	refused := func(log, errs []string) {
		t.Helper()
		if len(errs) != 1 || !strings.Contains(errs[0], "i/o timeout is not an answer from the engine") {
			t.Errorf("got %q", errs)
		}
		if len(log) != 1 || log[0] != "i/o timeout" {
			t.Errorf("the test got %q, want the driver's own error", log)
		}
	}

	// With no committed tape, none appears.
	refused(run(t, "TestFake", timeout))
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("a failed recording wrote %s", path)
	}

	// With one, it stays as it was.
	resetRegistry()
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) > 0 {
		t.Fatal(errs)
	}
	committed := read()
	resetRegistry()
	refused(run(t, "TestFake", timeout))
	if got := read(); got != committed {
		t.Errorf("a failed recording rewrote the tape:\n%s", got)
	}

	// A subtest that finished before the failure has rewritten the tape by
	// then. The failure puts the committed tape back, and a subtest after it
	// leaves the tape alone.
	resetRegistry()
	if _, errs := run(t, "TestFake/before", healthy); len(errs) > 0 {
		t.Fatal(errs)
	}
	if read() == committed {
		t.Fatal("the first subtest did not rewrite the tape, so the rest checks nothing")
	}
	refused(run(t, "TestFake/failing", timeout))
	if _, errs := run(t, "TestFake/after", healthy); len(errs) > 0 {
		t.Fatal(errs)
	}
	if got := read(); got != committed {
		t.Errorf("a failed recording left a partial tape:\n%s", got)
	}
}

func TestARecordingReplacesATapeThatDoesNotParse(t *testing.T) {
	srv, _ := harness(t)
	path := filepath.Join(dir, "TestFake.tape")
	if err := os.WriteFile(path, []byte("<<<<<<< HEAD\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	replay(t)
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) != 1 || !strings.Contains(errs[0], `TestFake.tape:1: unexpected line "<<<<<<< HEAD"`) {
		t.Errorf("a replay: %q", errs)
	}
	record(t)
	resetRegistry()
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) > 0 {
		t.Fatalf("a recording: %q", errs)
	}
	replay(t)
	if _, errs := run(t, "TestFake", everything(srv)); len(errs) > 0 {
		t.Errorf("the replay of it: %q", errs)
	}
}

// Only a passing run of every test settles the tapes no test opened. A run
// that failed, selected tests, or ran none (-list, a fuzz worker) proves
// nothing about them, and must neither fail on them nor delete them.
func TestOnlyAPassingRunOfEveryTestAccountsForTapes(t *testing.T) {
	for _, c := range []struct {
		flags map[string]string
		want  bool
	}{
		{nil, false},
		{map[string]string{"test.count": "2", "test.fuzzworker": "false"}, false},
		{map[string]string{"test.run": "^TestX$"}, true},
		{map[string]string{"test.skip": "TestX"}, true},
		{map[string]string{"test.list": "."}, true},
		{map[string]string{"test.fuzzworker": "true"}, true},
	} {
		if got := partialRun(func(name string) string { return c.flags[name] }); got != c.want {
			t.Errorf("partialRun with %v = %v, want %v", c.flags, got, c.want)
		}
	}

	harness(t)
	stray := filepath.Join(dir, "TestGone.tape")
	if err := os.WriteFile(stray, []byte("-- pin: default img:1 sha256:1 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	for _, mode := range []func(*testing.T){record, replay} {
		mode(t)
		if code := account(1, false, &out); code != 1 {
			t.Errorf("a failed run exits %d, want 1", code)
		}
		if code := account(0, true, &out); code != 0 {
			t.Errorf("a partial run exits %d, want 0", code)
		}
		if _, err := os.Stat(stray); err != nil || out.Len() > 0 {
			t.Fatalf("a failed or partial run settled the tapes: %v, %q", err, out.String())
		}
	}

	replay(t)
	if code := account(0, false, &out); code != 1 || !strings.Contains(out.String(), "TestGone.tape belongs to no test that ran") {
		t.Errorf("a passing replay of every test exits %d with %q", code, out.String())
	}
	record(t)
	if code := account(0, false, &out); code != 0 {
		t.Errorf("a passing recording of every test exits %d", code)
	}
	if _, err := os.Stat(stray); !os.IsNotExist(err) {
		t.Errorf("a passing recording of every test kept %s", stray)
	}
}

func TestLeftovers(t *testing.T) {
	srv, _ := harness(t)
	record(t)
	if _, errs := run(t, "TestFake", func(t testing.TB) []string {
		Open(t, srv, "default")
		Open(t, srv, "default")
		return nil
	}); len(errs) > 0 {
		t.Fatal(errs)
	}
	stray := filepath.Join(dir, "TestGone.tape")
	if err := os.WriteFile(stray, []byte("-- pin: default img:1 sha256:1 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	replay(t)
	if _, errs := run(t, "TestFake", func(t testing.TB) []string {
		Open(t, srv, "default")
		return nil
	}); len(errs) > 0 {
		t.Fatal(errs)
	}
	got := strings.Join(leftovers(false), "\n")
	for _, want := range []string{"TestGone.tape belongs to no test that ran", `no test opened stream "TestFake default#2"`} {
		if !strings.Contains(got, want) {
			t.Errorf("leftovers lack %q: %s", want, got)
		}
	}
	if got := leftovers(true); len(got) > 0 {
		t.Errorf("a recording's leftovers: %q", got)
	}
	if _, err := os.Stat(stray); !os.IsNotExist(err) {
		t.Errorf("a recording kept %s", stray)
	}
}

// A request's text or argument over 64 KB is kept as a digest: the tape stays
// small, and a replay still tells one big statement from another.
func TestBigRequestsAreKeptAsDigests(t *testing.T) {
	srv, _ := harness(t)
	big := "SELECT 0" + strings.Repeat(" OR 1", 20000)
	send := func(stmt string) func(t testing.TB) []string {
		return func(t testing.TB) []string {
			db := Open(t, srv, "default")
			_, err1 := db.Exec(stmt)
			_, err2 := db.Exec("SET @v = ?", stmt)
			return []string{fmt.Sprint(err1, err2)}
		}
	}
	record(t)
	if _, errs := run(t, "TestFake", send(big)); len(errs) > 0 {
		t.Fatal(errs)
	}
	b, err := os.ReadFile(filepath.Join(dir, "TestFake.tape"))
	if err != nil {
		t.Fatal(err)
	}
	if len(b) > 1000 || strings.Count(string(b), "|# sha256:") != 2 || !strings.Contains(string(b), fmt.Sprintf(" %d\n", len(big))) {
		t.Errorf("the tape is %d bytes:\n%.600s", len(b), b)
	}

	replay(t)
	if log, errs := run(t, "TestFake", send(big)); len(errs) > 0 || log[0] != "<nil> <nil>" {
		t.Errorf("replay: %q %q", log, errs)
	}
	resetRegistry()
	other := strings.Replace(big, "OR 1", "OR 2", 1)
	if _, errs := run(t, "TestFake", send(other)); len(errs) == 0 || !strings.Contains(errs[0], "no recorded answer for exec <sha256:") {
		t.Errorf("another statement of the same size replayed: %q", errs)
	}
}

// A request that repeats the text of the one before it writes ^ for it, and
// reads back whole. A tape written before the mark existed reads the same.
func TestRepeatedRequestTextIsWrittenOnce(t *testing.T) {
	prepare := func(arg string) exchange {
		return exchange{request{kind: "exec", text: "PREPARE s FROM ?", args: []value{{kind: "string", text: arg}}}, answer{kind: "ok", affected: "0", lastID: "0"}}
	}
	f := &tapeFile{pins: map[string]string{}, streams: []*stream{
		{name: "a", exchanges: []exchange{prepare("SELECT 1"), prepare("SELECT 2"), {request{kind: "query", text: "PREPARE s FROM ?"}, answer{kind: "badconn"}}, prepare("SELECT 3")}},
		{name: "b", exchanges: []exchange{prepare("SELECT 4")}},
	}}
	text := f.encode()
	want := "\n== a\n> exec\n| PREPARE s FROM ?\n: string\n| SELECT 1\n< ok 0 0\n> exec ^\n: string\n| SELECT 2\n< ok 0 0\n" +
		"> query\n| PREPARE s FROM ?\n< badconn\n> exec\n| PREPARE s FROM ?\n: string\n| SELECT 3\n< ok 0 0\n" +
		"\n== b\n> exec\n| PREPARE s FROM ?\n: string\n| SELECT 4\n< ok 0 0\n"
	if text != want {
		t.Fatalf("encoded:\n%s\nwant:\n%s", text, want)
	}
	back, err := parse("x.tape", text)
	if err != nil {
		t.Fatal(err)
	}
	if got := back.streams[0].exchanges[1].req.text; got != "PREPARE s FROM ?" {
		t.Errorf("the second request reads back with text %q", got)
	}
	whole := strings.Replace(text, "> exec ^\n", "> exec\n| PREPARE s FROM ?\n", 1)
	old, err := parse("old.tape", whole)
	if err != nil || old.encode() != text {
		t.Errorf("a tape with every text written out reads differently: %v", err)
	}
}

// Two requests often differ only in an argument, so the failure names it.
func TestReplayNamesTheArgumentThatDiffers(t *testing.T) {
	srv, _ := harness(t)
	prepare := func(stmt string) func(t testing.TB) []string {
		return func(t testing.TB) []string {
			_, _ = Open(t, srv, "default").Exec("PREPARE s FROM ?", stmt)
			return nil
		}
	}
	record(t)
	if _, errs := run(t, "TestFake", prepare("SELECT 1")); len(errs) > 0 {
		t.Fatal(errs)
	}
	replay(t)
	_, errs := run(t, "TestFake", prepare("SELECT 2"))
	if len(errs) != 1 || !strings.Contains(errs[0], `no recorded answer for exec "PREPARE s FROM ?" string "SELECT 2"`) || !strings.Contains(errs[0], `the tape holds exec "PREPARE s FROM ?" string "SELECT 1" there`) {
		t.Errorf("got %q", errs)
	}
}

func TestParseRejects(t *testing.T) {
	for _, c := range []struct{ name, tape, want string }{
		{"a request before a stream", "> exec\n| x\n< ok\n", "before any stream"},
		{"no answer", "== s\n> exec\n| x\n", "want an answer"},
		{"no text", "== s\n> exec\n< ok\n", "want text"},
		{"an unknown request", "== s\n> fetch\n< ok\n", "unknown request"},
		{"a bad int", "== s\n> query\n| x\n< result\n< column \"a\" \"INT\"\n< row\n: int64 x\n", "a int64 value"},
		{"a short row", "== s\n> query\n| x\n< result\n< column \"a\" \"INT\"\n< row\n", "want a value"},
		{"a stream twice", "== s\n== s\n", "appears twice"},
		{"a stray line", "== s\nhello\n", "unexpected line"},
		{"a ditto with nothing before it", "== s\n> exec ^\n< ok 0 0\n", "no exec before it"},
		{"a ditto of another kind", "== s\n> query\n| x\n< badconn\n> exec ^\n< ok 0 0\n", "no exec before it"},
		{"a digest in an answer", "== s\n> query\n| x\n< result\n< column \"a\" \"TEXT\"\n< row\n: bytes\n|# sha256:ab 70000\n", "answers are kept whole"},
	} {
		if _, err := parse("x.tape", c.tape); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want %q", c.name, err, c.want)
		}
	}
}

func TestFormatRoundTrips(t *testing.T) {
	f := &tapeFile{pins: map[string]string{"default": "default a b c", "lower": "lower d e f"}, streams: []*stream{{
		name: "TestX/sub name default#1",
		exchanges: []exchange{
			{request{kind: "query", text: "", args: []value{{kind: "string", text: ""}, {kind: "float64", text: "1.5"}, {kind: "bool", text: "true"}, {kind: "time", text: "2026-11-12T00:00:00Z"}}},
				answer{kind: "rows", sets: []resultSet{{cols: []column{{"a \"b\"", "UNSIGNED BIGINT"}}, rows: [][]value{{{kind: "bytes", text: "\n"}}, {{kind: "bytes", text: "\x00\xff"}}, {{kind: "string", text: "| q"}}}}, {}}}},
			{request{kind: "exec", text: "|q x"}, answer{kind: "badconn"}},
			{request{kind: "exec", sum: "sha256:ab 70000", args: []value{{kind: "bytes", sum: "sha256:cd 80000"}}}, answer{kind: "ok", affected: "1", lastID: "-"}},
			{request{kind: "rollback"}, answer{kind: "error", err: "\ttab"}},
		},
	}}}
	text := f.encode()
	back, err := parse("x.tape", text)
	if err != nil {
		t.Fatalf("%v\n%s", err, text)
	}
	if again := back.encode(); again != text {
		t.Errorf("a round trip changed the tape:\n%s\nbecame:\n%s", text, again)
	}
}

// The fake engine.

type fakeEngine struct {
	now  int
	made map[string]bool // the names it made up and has not dropped
}

type fakeConnector struct{ e *fakeEngine }

func (c fakeConnector) Connect(context.Context) (driver.Conn, error) { return &fakeConn{c.e}, nil }
func (c fakeConnector) Driver() driver.Driver                        { return nil }

type fakeConn struct{ e *fakeEngine }

func (c *fakeConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("prepared") }
func (c *fakeConn) Close() error                        { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)           { return fakeTx{}, nil }

func (c *fakeConn) ExecContext(_ context.Context, q string, args []driver.NamedValue) (driver.Result, error) {
	switch q {
	case "BAD":
		return nil, &fakeError{1064, "near \"BAD\"\nline 2"}
	case "SKIP":
		return nil, driver.ErrSkip
	case "TIMEOUT":
		return nil, errors.New("i/o timeout")
	}
	if name, ok := strings.CutPrefix(q, "DROP "); ok {
		if name == "?" {
			name = args[0].Value.(string)
		}
		if !c.e.made[name] {
			return nil, &fakeError{1091, "can't drop " + name}
		}
		delete(c.e.made, name)
		return driver.RowsAffected(1), nil
	}
	return driver.RowsAffected(len(args)), nil
}

func (c *fakeConn) QueryContext(_ context.Context, q string, _ []driver.NamedValue) (driver.Rows, error) {
	switch q {
	case "SELECT rows":
		return &fakeRows{sets: []fakeSet{
			{[]string{"a", "b"}, []string{"INT", "TEXT"}, [][]driver.Value{{int64(1), []byte("x\ny")}, {nil, []byte("\r")}}},
			{[]string{"c"}, []string{"BOOL"}, [][]driver.Value{{true}}},
		}}, nil
	case "SELECT now":
		c.e.now++
		return &fakeRows{sets: []fakeSet{{[]string{"now"}, []string{"INT"}, [][]driver.Value{{int64(c.e.now)}}}}}, nil
	case "SELECT made":
		c.e.now++
		name := fmt.Sprintf("made_%d", 9000+c.e.now)
		if c.e.made == nil {
			c.e.made = map[string]bool{}
		}
		c.e.made[name] = true
		return &fakeRows{sets: []fakeSet{{[]string{"name"}, []string{"TEXT"}, [][]driver.Value{{[]byte("x " + name + " y")}}}}}, nil
	}
	return nil, &fakeError{1146, "no such table"}
}

type fakeTx struct{}

func (fakeTx) Commit() error   { return nil }
func (fakeTx) Rollback() error { return nil }

type fakeSet struct {
	cols, types []string
	rows        [][]driver.Value
}

type fakeRows struct {
	sets     []fakeSet
	set, row int
}

func (r *fakeRows) Columns() []string                       { return r.sets[r.set].cols }
func (r *fakeRows) ColumnTypeDatabaseTypeName(i int) string { return r.sets[r.set].types[i] }
func (r *fakeRows) Close() error                            { return nil }
func (r *fakeRows) HasNextResultSet() bool                  { return r.set+1 < len(r.sets) }
func (r *fakeRows) NextResultSet() error {
	if !r.HasNextResultSet() {
		return io.EOF
	}
	r.set, r.row = r.set+1, 0
	return nil
}

func (r *fakeRows) Next(dest []driver.Value) error {
	s := r.sets[r.set]
	if r.row >= len(s.rows) {
		return io.EOF
	}
	copy(dest, s.rows[r.row])
	r.row++
	return nil
}

type fakeError struct {
	code int
	msg  string
}

func (e *fakeError) Error() string { return fmt.Sprintf("Error %d: %s", e.code, e.msg) }

type fakeCodec struct{}

func (fakeCodec) EncodeError(err error) (string, bool) {
	var fe *fakeError
	if !errors.As(err, &fe) {
		return "", false
	}
	return fmt.Sprintf("%d %s", fe.code, fe.msg), true
}

func (fakeCodec) DecodeError(text string) error {
	// An older spelling the codec still reads.
	text = strings.TrimPrefix(text, "error ")
	code, msg, _ := strings.Cut(text, " ")
	n, err := strconv.Atoi(code)
	if err != nil {
		return errors.New(text)
	}
	return &fakeError{n, msg}
}

// renamingCodec is fakeCodec for an engine that makes up names: made_ and a
// number over 9000. Its stand-ins count from 1.
type renamingCodec struct{ fakeCodec }

var fakeMadeUp = regexp.MustCompile(`made_\d+`)

func (renamingCodec) MadeUp() *regexp.Regexp { return fakeMadeUp }

func (renamingCodec) StandIn(_ string, n int) string { return fmt.Sprintf("made_%d", n) }
