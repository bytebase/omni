package simulate

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

// recorder stands in for the test Check is handed, so a test can read what
// Check reported.
type recorder struct {
	testing.TB
	errors []string
	fatal  string
}

func (r *recorder) Helper()                         {}
func (r *recorder) Logf(format string, args ...any) {}
func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}
func (r *recorder) Fatalf(format string, args ...any) {
	r.fatal = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

func (r *recorder) failed(substr string) bool {
	for _, e := range append(r.errors, r.fatal) {
		if strings.Contains(e, substr) {
			return true
		}
	}
	return false
}

func (r *recorder) passed() bool { return len(r.errors) == 0 && r.fatal == "" }

// run checks s against a recorder, on a goroutine of its own because Fatalf
// ends the goroutine it runs on.
func run(t *testing.T, s Suite) *recorder {
	t.Helper()
	rec := &recorder{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		check(rec, s)
	}()
	<-done
	return rec
}

const header = "# pin: default img:1 sha256:abc v1\n" + columnsHeader + "\n"

// The engine of these tests: what it answers each statement on fixture fx.
var (
	stored1 = Answer{Forms: []Form{{"TABLE", "t", "CREATE TABLE t (a int)"}}}
	stored2 = Answer{Forms: []Form{{"TABLE", "t", "CREATE TABLE t (a int, b int)"}}}
	refused = Answer{Code: "1060"}
)

// suite writes pins.tsv, a fixture and simulate.tsv into a fresh directory.
// engine maps a statement to the engine's answer and omni to omni's; the
// suite records what the engine was asked in *asked.
func suite(t *testing.T, rows string, engine map[string]Answer, omni map[string]Omni, asked *[]string) Suite {
	t.Helper()
	dir := t.TempDir()
	write(t, filepath.Join(dir, "pins.tsv"), "default\timg:1\tsha256:abc\tv1\n")
	write(t, filepath.Join(dir, "fx.sql"), "CREATE TABLE t (a int)\n")
	// The suite's generated rows are the file's, unless a test says otherwise.
	var generated []Generated
	if rows != "" {
		write(t, filepath.Join(dir, "simulate.tsv"), header+rows)
		f, err := Parse("simulate.tsv", []byte(header+rows))
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range f.Rows {
			if !r.hand() {
				generated = append(generated, Generated{r.Fixture, r.Rule, r.Statement})
			}
		}
	}
	return Suite{
		Rows:     generated,
		Path:     filepath.Join(dir, "simulate.tsv"),
		Pins:     filepath.Join(dir, "pins.tsv"),
		Fixtures: dir,
		Omni:     func(_ Fixture, statement string, _ Answer) Omni { return omni[statement] },
		Engine: func(testing.TB) func(Fixture, string) (Answer, error) {
			return func(fx Fixture, statement string) (Answer, error) {
				if fx.SQL != "CREATE TABLE t (a int)" {
					return Answer{}, fmt.Errorf("fixture %s read as %q", fx.Name, fx.SQL)
				}
				if asked != nil {
					*asked = append(*asked, statement)
				}
				a, ok := engine[statement]
				if !ok {
					return Answer{}, errors.New("connection lost")
				}
				return a, nil
			}
		},
	}
}

func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func updating(t *testing.T) {
	t.Helper()
	if err := flag.Set("update-compat", "true"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = flag.Set("update-compat", "false") })
}

func recording(t *testing.T) {
	t.Helper()
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
}

var ok1, ok2 = stored1.Outcome(nil), stored2.Outcome(nil)

func TestCheckPassesAgreementsAndKnownDifferences(t *testing.T) {
	rows := ok1 + "\t" + ok1 + "\tengine\tfx\tr/A\tA\n" + // both take it and store the same
		"1060\t1060\tengine\tfx\tr/B\tB\n" + // both refuse it alike
		"ok\t" + ok2 + "\tengine\tfx\tr/C\tC\n" + // omni stores another schema
		"ok\t1060\tengine\tfx\tr/D\tD\n" + // only omni takes it
		"1064\t" + ok1 + "\tengine\tfx\tr/E\tE\n" + // only the engine takes it
		"error\t1060\tengine\tfx\tr/F\tF\n" // both refuse, omni with no code
	s := suite(t, rows,
		map[string]Answer{"A": stored1, "B": refused, "C": stored2, "D": refused, "E": stored1, "F": refused},
		map[string]Omni{"A": {Same: true}, "B": {Code: "1060"}, "C": {}, "D": {Same: true}, "E": {Code: "1064"}, "F": {Code: NoCode}}, nil)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("want a pass, got %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckFailsAClosedGapUntilUpdated(t *testing.T) {
	s := suite(t, "ok\t"+ok2+"\tengine\tfx\tr/C\tC\n", map[string]Answer{"C": stored2}, map[string]Omni{"C": {Same: true}}, nil)
	if rec := run(t, s); !rec.failed("a gap closed") {
		t.Errorf("a closed gap passed: %q", rec.errors)
	}
	updating(t)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("-update-compat: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); !strings.Contains(got, ok2+"\t"+ok2+"\tengine\tfx\tr/C\tC\n") {
		t.Errorf("-update-compat did not commit the closed gap:\n%s", got)
	}
}

func TestCheckFailsADifferenceThatMovedUntilUpdated(t *testing.T) {
	// A parser fix: omni no longer answers 1064, and now takes a statement
	// the engine refuses.
	s := suite(t, "1064\t1060\tengine\tfx\tr/B\tB\n", map[string]Answer{"B": refused}, map[string]Omni{"B": {}}, nil)
	rec := run(t, s)
	if !rec.failed("a known difference is now another one") || rec.failed("a regression") {
		t.Errorf("a difference that moved was not reported as one: %q", rec.errors)
	}
	updating(t)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("-update-compat: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); !strings.Contains(got, "ok\t1060\tengine\tfx\tr/B\tB\n") {
		t.Errorf("-update-compat did not commit the moved difference:\n%s", got)
	}
}

func TestCheckFailsARegressionAndNeverWritesIt(t *testing.T) {
	s := suite(t, ok1+"\t"+ok1+"\tengine\tfx\tr/A\tA\n", map[string]Answer{"A": stored1}, map[string]Omni{"A": {Why: "would add b"}}, nil)
	before := read(t, s.Path)
	updating(t)
	rec := run(t, s)
	if !rec.failed("a regression") || !rec.failed("would add b") {
		t.Errorf("a regression passed, or lost its reason: %q", rec.errors)
	}
	if read(t, s.Path) != before {
		t.Errorf("-update-compat wrote a regression")
	}
}

func TestCheckFailsAnUnjudgedRowWithoutAskingTheEngine(t *testing.T) {
	var asked []string
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n?\t?\t-\tfx\tissue:BYT-1\tA\n", map[string]Answer{"B": refused}, nil, &asked)
	if rec := run(t, s); !strings.Contains(rec.fatal, "not judged yet") {
		t.Errorf("an unjudged row passed: %q %q", rec.errors, rec.fatal)
	}
	if len(asked) > 0 {
		t.Errorf("a replay with an unjudged row asked the tape for %q", asked)
	}
}

func TestCheckFailsWhenTheTapeAndTheFileDisagree(t *testing.T) {
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n", map[string]Answer{"B": {Code: "1061"}}, map[string]Omni{"B": {Code: "1060"}}, nil)
	before := read(t, s.Path)
	if rec := run(t, s); !rec.failed("the tape answers 1061, and the row records 1060") {
		t.Errorf("an expected cell that is not the tape's answer passed: %q", rec.errors)
	}
	if read(t, s.Path) != before {
		t.Errorf("a replay wrote the file")
	}
}

func TestCheckFailsAnotherPin(t *testing.T) {
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n", map[string]Answer{"B": refused}, map[string]Omni{"B": {Code: "1060"}}, nil)
	write(t, s.Pins, "default\timg:2\tsha256:def\tv2\n")
	if rec := run(t, s); !strings.Contains(rec.fatal, "was recorded at pin") {
		t.Errorf("a file of another pin passed: %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckHoldsGeneratedRowsToTheSuite(t *testing.T) {
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n1060\t1060\tengine\tfx\tr/Gone\tG\n1060\t1060\tengine\tfx\tissue:BYT-1\tH\n",
		map[string]Answer{"A": stored1, "B": refused, "H": refused},
		map[string]Omni{"A": {Same: true}, "B": {Code: "1060"}, "H": {Code: "1060"}}, nil)
	s.Rows = []Generated{{"fx", "r/A", "A"}, {"fx", "r/B", "B"}}
	if rec := run(t, s); !strings.Contains(rec.fatal, "1 to add, 1 to remove") {
		t.Errorf("a file that is not the suite's rows passed: %q %q", rec.errors, rec.fatal)
	}

	recording(t)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("recording: %q %q", rec.errors, rec.fatal)
	}
	want := header + ok1 + "\t" + ok1 + "\tengine\tfx\tr/A\tA\n1060\t1060\tengine\tfx\tr/B\tB\n1060\t1060\tengine\tfx\tissue:BYT-1\tH\n"
	if got := read(t, s.Path); got != want {
		t.Errorf("the recording wrote:\n%s\nwant the suite's rows in order, then the hand-added row:\n%s", got, want)
	}
}

// A compat.tsv row a person added there reaches the suite as a generated row
// with its hand label. The file must know it as the suite's: read as
// hand-added, it would be missing from every replay and added again by
// every recording.
func TestCheckKnowsAGeneratedRowByTheSuiteNotItsLabel(t *testing.T) {
	rows := "1060\t1060\tengine\tfx\tissue:BYT-1\tB\n"
	s := suite(t, rows, map[string]Answer{"B": refused}, map[string]Omni{"B": {Code: "1060"}}, nil)
	s.Rows = []Generated{{"fx", "issue:BYT-1", "B"}}
	if rec := run(t, s); !rec.passed() {
		t.Errorf("replay: %q %q", rec.errors, rec.fatal)
	}
	recording(t)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); got != header+rows {
		t.Errorf("the recording wrote:\n%s\nwant the one row", got)
	}
}

func TestRecordingStartsAMissingFileAndAsksEveryRow(t *testing.T) {
	recording(t)
	var asked []string
	s := suite(t, "", map[string]Answer{"A": stored1, "B": refused}, map[string]Omni{"A": {}, "B": {Code: "1060"}}, &asked)
	s.Rows = []Generated{{"fx", "r/A", "A"}, {"fx", "r/B", "B"}}
	if rec := run(t, s); !rec.passed() {
		t.Errorf("recording: %q %q", rec.errors, rec.fatal)
	}
	if want := header + "ok\t" + ok1 + "\tengine\tfx\tr/A\tA\n1060\t1060\tengine\tfx\tr/B\tB\n"; read(t, s.Path) != want {
		t.Errorf("the recording wrote:\n%s\nwant:\n%s", read(t, s.Path), want)
	}

	// A second recording asks every row again, since its tape is written
	// whole, and changes nothing.
	asked = nil
	before := read(t, s.Path)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("second recording: %q %q", rec.errors, rec.fatal)
	}
	if strings.Join(asked, ",") != "A,B" || read(t, s.Path) != before {
		t.Errorf("the second recording asked %q and left:\n%s", asked, read(t, s.Path))
	}
}

func TestRecordingFailsAnAnswerThatDoesNotReproduce(t *testing.T) {
	recording(t)
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n", map[string]Answer{"B": {Code: "1061"}}, map[string]Omni{"B": {Code: "1060"}}, nil)
	before := read(t, s.Path)
	if rec := run(t, s); !rec.failed("does not reproduce") {
		t.Errorf("an answer that changed at the same pin passed: %q", rec.errors)
	}
	if read(t, s.Path) != before {
		t.Errorf("the recording wrote an answer that does not reproduce")
	}
}

func TestRecordingAtANewPinJudgesEveryRowAgain(t *testing.T) {
	recording(t)
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n", map[string]Answer{"B": {Code: "1061"}}, map[string]Omni{"B": {Code: "1060"}}, nil)
	write(t, s.Pins, "default\timg:2\tsha256:def\tv2\n")
	// omni still answers 1060 and the new server 1061: a known difference.
	if rec := run(t, s); !rec.passed() {
		t.Errorf("recording at a new pin: %q %q", rec.errors, rec.fatal)
	}
	if want := "# pin: default img:2 sha256:def v2\n" + columnsHeader + "\n1060\t1061\tengine\tfx\tr/B\tB\n"; read(t, s.Path) != want {
		t.Errorf("the recording wrote:\n%s\nwant:\n%s", read(t, s.Path), want)
	}
}

func TestRecordingStopsWhenTheEngineGivesNoAnswer(t *testing.T) {
	recording(t)
	s := suite(t, "?\t?\t-\tfx\tissue:BYT-1\tA\n", nil, nil, nil)
	before := read(t, s.Path)
	if rec := run(t, s); !strings.Contains(rec.fatal, "the engine gave no answer: connection lost") {
		t.Errorf("no answer was read as one: %q %q", rec.errors, rec.fatal)
	}
	if read(t, s.Path) != before {
		t.Errorf("the recording wrote a row the engine did not answer")
	}
}

func TestRecordingFailsInCI(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "true")
	s := suite(t, "1060\t1060\tengine\tfx\tr/B\tB\n", map[string]Answer{"B": refused}, map[string]Omni{"B": {Code: "1060"}}, nil)
	if rec := run(t, s); !strings.Contains(rec.fatal, "CI") {
		t.Errorf("a recording ran in CI: %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckFailsARowWhoseFixtureIsMissing(t *testing.T) {
	s := suite(t, "1060\t1060\tengine\tnosuch\tr/B\tB\n", map[string]Answer{"B": refused}, nil, nil)
	if rec := run(t, s); !strings.Contains(rec.fatal, "names fixture nosuch") {
		t.Errorf("a row on a missing fixture passed: %q %q", rec.errors, rec.fatal)
	}
}

func TestOutcomeIsOfTheRowsOwnText(t *testing.T) {
	madeUp := regexp.MustCompile(`c_\d+`)
	form := func(names ...string) Answer {
		return Answer{Forms: []Form{{"TABLE", "t", "CREATE TABLE t (CONSTRAINT " + names[0] + " CHECK (a), CONSTRAINT " + names[1] + " CHECK (b), KEY (" + names[2] + "))"}}}
	}
	first := form("c_1", "c_2", "c_1").Outcome(madeUp)
	if !isStored(first) {
		t.Fatalf("outcome %q is not ok:<hash>", first)
	}
	if later := form("c_7", "c_9", "c_7").Outcome(madeUp); later != first {
		t.Errorf("the hash moved with the stand-ins a tape gave the made-up names: %s, %s", first, later)
	}
	if swapped := form("c_1", "c_2", "c_2").Outcome(madeUp); swapped == first {
		t.Errorf("the hash does not tell one made-up name from another")
	}
	if form("c_1", "c_2", "c_1").Outcome(nil) == form("c_7", "c_9", "c_7").Outcome(nil) {
		t.Errorf("with no made-up names, two texts hash alike")
	}
	if got := (Answer{Code: "1146", Forms: stored1.Forms}).Outcome(nil); got != "1146" {
		t.Errorf("a refusal's outcome is %q, want its code", got)
	}
	if (Answer{}).Outcome(nil) == stored1.Outcome(nil) || !isStored((Answer{}).Outcome(nil)) {
		t.Errorf("an empty database is not a stored schema of its own")
	}
}

func TestParseRejectsMalformedRows(t *testing.T) {
	for _, c := range []struct{ rows, want string }{
		{"1060\t1060\tengine\tfx\tr/B\n", "want 6 tab-separated columns"},
		{"accept\t1060\tengine\tfx\tr/B\tB\n", "omni:"},
		{"1060\tok\tengine\tfx\tr/B\tB\n", "expected:"},
		{"1060\terror\tengine\tfx\tr/B\tB\n", "expected:"},
		{"1060\tok:123\tengine\tfx\tr/B\tB\n", "expected:"},
		{"01060\t1060\tengine\tfx\tr/B\tB\n", "omni:"},
		{"1060\t1060\tgrammar\tfx\tr/B\tB\n", "source"},
		{"1060\t?\tengine\tfx\tr/B\tB\n", "no expected outcome"},
		{"1060\t1060\t-\tfx\tr/B\tB\n", "with no source"},
		{"1060\t1060\tengine\t\tr/B\tB\n", "empty fixture"},
		{"1060\t1060\tengine\tf x\tr/B\tB\n", "holds a space"},
		{"1060\t1060\tengine\tfx\tr/B\tB\n1060\t1060\tengine\tfx\tr/B\tB\n", "repeats line"},
		{"1060\t1060\tengine\tfx\tr/B\tB\n# late\n", "a comment after the first row"},
	} {
		if _, err := Parse("simulate.tsv", []byte(header+c.rows)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: got %v, want an error holding %q", c.rows, err, c.want)
		}
	}
	// One statement on two fixtures is two rows, and a statement keeps its tabs.
	f, err := Parse("simulate.tsv", []byte(header+"1060\t1060\tengine\tfx\tr/B\tB\n1060\t1060\tengine\tfy\tr/B\tB\tx \n"))
	if err != nil || len(f.Rows) != 2 || f.Rows[1].Statement != "B\tx " {
		t.Errorf("got %v %+v", err, f)
	}
}

func TestFromCompat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "compat.tsv")
	write(t, path, "# pin: p\n# columns: omni, expected, source, rule, statement\n"+
		"accept\taccept\tengine\tr/A\tALTER 1\n"+
		"reject\taccept\tengine\tr/A2\tALTER 1\n"+ // the same statement by another rule
		"accept\treject\tengine\tr/B\tALTER 2\n"+ // the engine rejects it
		"reject\taccept\tengine\tr/C!drop-1\tALTER 3\n"+ // a near-miss the engine accepts
		"accept\taccept\tengine\tr/D\tSELECT 1\n") // not picked
	got, err := FromCompat(path, []string{"f1", "f2"}, func(s string) bool { return strings.HasPrefix(s, "ALTER") })
	if err != nil {
		t.Fatal(err)
	}
	want := []Generated{{"f1", "r/A", "ALTER 1"}, {"f1", "r/C!drop-1", "ALTER 3"}, {"f2", "r/A", "ALTER 1"}, {"f2", "r/C!drop-1", "ALTER 3"}}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGuard(t *testing.T) {
	parse := func(rows string) *File {
		t.Helper()
		f, err := Parse("simulate.tsv", []byte(rows))
		if err != nil {
			t.Fatal(err)
		}
		return f
	}
	const pin1, pin2 = "# pin: default e:1 sha256:1 1\n", "# pin: default e:2 sha256:2 2\n"
	for _, c := range []struct {
		name          string
		before, after string
		want          string // a substring of the problem, or "" for none
	}{
		{"a closed gap", "ok\t" + ok1 + "\tengine\tfx\tr/A\tA\n", ok1 + "\t" + ok1 + "\tengine\tfx\tr/A\tA\n", ""},
		{"a committed regression", "1060\t1060\tengine\tfx\tr/A\tA\n", "ok\t1060\tengine\tfx\tr/A\tA\n", "a regression"},
		{"a difference that moved", "1064\t1060\tengine\tfx\tr/A\tA\n", "ok\t1060\tengine\tfx\tr/A\tA\n", ""},
		{"a recording fills a row", "?\t?\t-\tfx\tissue:X\tA\n", "ok\t1060\tengine\tfx\tissue:X\tA\n", ""},
		{"an edited expected", "ok\t1060\tengine\tfx\tr/A\tA\n", "ok\t" + ok1 + "\tengine\tfx\tr/A\tA\n", "without a pin bump"},
		{"a pin bump", pin1 + "ok\t1060\tengine\tfx\tr/A\tA\n", pin2 + "ok\t1061\tengine\tfx\tr/A\tA\n", ""},
		{"the same row on another fixture", "ok\t1060\tengine\tfx\tr/A\tA\n", "ok\t1060\tengine\tfx\tr/A\tA\n1061\t1061\tengine\tfy\tr/A\tA\n", ""},
	} {
		got := strings.Join(Guard(parse(c.before), parse(c.after)), "\n")
		if c.want == "" && got != "" || !strings.Contains(got, c.want) {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
	if !PinChanged(parse(pin1), parse(pin2)) || PinChanged(parse(pin1), parse(pin1)) {
		t.Errorf("PinChanged does not read the pin line")
	}
}
