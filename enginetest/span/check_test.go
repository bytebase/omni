package span

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

// answer is an engine's or omni's span of query Q of these tests.
var answer = Answer{Results: []string{"a", "_col1"}, Lineage: [][]string{{"c.s.t.a"}, {"c.s.t.b", "c.s.t.a"}}, Tables: []string{"c.s.t"}}

// rows renders the three rows of query Q in session c.s on fixture fx.
func rows(omni, engine Answer) string {
	var b strings.Builder
	for _, aspect := range Aspects {
		fmt.Fprintf(&b, "%s\t%s\tengine\tfx\tc.s\t%s\tQ\n", omni.Cell(aspect), engine.Cell(aspect), aspect)
	}
	return b.String()
}

// suite writes pins.tsv, a fixture and span.tsv into a fresh directory.
func suite(t *testing.T, file string, engine, omni Answer, asked *int) Suite {
	t.Helper()
	dir := t.TempDir()
	write(t, filepath.Join(dir, "pins.tsv"), "default\timg:1\tsha256:abc\tv1\n")
	write(t, filepath.Join(dir, "fx.sql"), "CREATE TABLE t (a int, b int);\n")
	if file != "" {
		write(t, filepath.Join(dir, "span.tsv"), header+file)
	}
	return Suite{
		Path:     filepath.Join(dir, "span.tsv"),
		Pins:     filepath.Join(dir, "pins.tsv"),
		Fixtures: dir,
		Queries:  []Query{{Fixture: "fx", Session: "c.s", Statement: "Q"}},
		Omni:     func(Fixture, string, string) Answer { return omni },
		Engine: func(testing.TB) func(Fixture, string, string) (Answer, error) {
			return func(fx Fixture, session, statement string) (Answer, error) {
				if fx.SQL != "CREATE TABLE t (a int, b int);" || session != "c.s" || statement != "Q" {
					return Answer{}, errors.New("asked about something else")
				}
				if asked != nil {
					*asked++
				}
				return engine, nil
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

func TestCells(t *testing.T) {
	for aspect, want := range map[string]string{
		Results: "a, _col1",
		Lineage: "c.s.t.a | c.s.t.a c.s.t.b",
		Tables:  "c.s.t",
	} {
		if got := answer.Cell(aspect); got != want {
			t.Errorf("%s: %q, want %q", aspect, got, want)
		}
	}
	empty := Answer{Results: []string{"", "x"}, Lineage: [][]string{nil, nil}}
	if got := empty.Cell(Results) + " / " + empty.Cell(Lineage) + " / " + empty.Cell(Tables); got != `"", x / - | - / -` {
		t.Errorf("an answer with no names or sources: %q", got)
	}
	if got := (Answer{Error: "COLUMN_NOT_FOUND"}).Cell(Tables); got != "error:COLUMN_NOT_FOUND" {
		t.Errorf("a refusal: %q", got)
	}
}

func TestCheckPassesAgreementsAndKnownDifferences(t *testing.T) {
	other := Answer{Results: []string{"a", "_col1"}, Lineage: [][]string{{"c.s.t.a"}, {}}, Tables: []string{"c.s.t"}}
	if rec := run(t, suite(t, rows(other, answer), answer, other, nil)); !rec.passed() {
		t.Errorf("want a pass, got %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckFailsAClosedGapUntilUpdated(t *testing.T) {
	other := Answer{Results: []string{"a", ""}, Lineage: answer.Lineage, Tables: answer.Tables}
	s := suite(t, rows(other, answer), answer, answer, nil)
	if rec := run(t, s); !rec.failed("a gap closed") {
		t.Errorf("a closed gap passed: %q", rec.errors)
	}
	updating(t)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("-update-compat: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); got != header+rows(answer, answer) {
		t.Errorf("-update-compat did not commit the closed gap:\n%s", got)
	}
}

func TestCheckFailsARegressionAndNeverWritesIt(t *testing.T) {
	s := suite(t, rows(answer, answer), answer, Answer{Error: "omni"}, nil)
	before := read(t, s.Path)
	updating(t)
	if rec := run(t, s); !rec.failed("a regression") {
		t.Errorf("a regression passed: %q", rec.errors)
	}
	if read(t, s.Path) != before {
		t.Errorf("-update-compat wrote a regression")
	}
}

func TestCheckFailsWhenTheTapeAndTheFileDisagree(t *testing.T) {
	s := suite(t, rows(answer, answer), Answer{Error: "TABLE_NOT_FOUND"}, answer, nil)
	if rec := run(t, s); !rec.failed("the tape answers error:TABLE_NOT_FOUND") {
		t.Errorf("an expected cell that is not the tape's answer passed: %q", rec.errors)
	}
}

func TestCheckFailsAnotherPin(t *testing.T) {
	s := suite(t, rows(answer, answer), answer, answer, nil)
	write(t, s.Pins, "default\timg:2\tsha256:def\tv2\n")
	if rec := run(t, s); !strings.Contains(rec.fatal, "was recorded at pin") {
		t.Errorf("a file of another pin passed: %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckHoldsTheFileToTheQueries(t *testing.T) {
	s := suite(t, rows(answer, answer), answer, answer, nil)
	s.Queries = append(s.Queries, Query{Fixture: "fx", Session: "c.s", Statement: "R"})
	if rec := run(t, s); !strings.Contains(rec.fatal, "3 to add, 0 to remove") {
		t.Errorf("a file missing a query's rows passed: %q %q", rec.errors, rec.fatal)
	}
}

func TestRecordingStartsAMissingFileAndAsksEachQueryOnce(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	asked := 0
	s := suite(t, "", answer, answer, &asked)
	if rec := run(t, s); !rec.passed() {
		t.Errorf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); got != header+rows(answer, answer) {
		t.Errorf("the recording wrote:\n%s", got)
	}
	if asked != 1 {
		t.Errorf("the engine was asked %d times about one query", asked)
	}
}

func TestParseRefusesABrokenRow(t *testing.T) {
	for _, row := range []string{
		"a\ta\tengine\tfx\tc.s\tcolumns\tQ\n",                                    // an unknown aspect
		"a\t?\tengine\tfx\tc.s\tresults\tQ\n",                                    // a source with no answer
		"a\ta\t-\tfx\tc.s\tresults\tQ\n",                                         // an answer with no source
		"a\ta\tengine\tfx\tc s\tresults\tQ\n",                                    // a space in the session
		"a\ta\tengine\tfx\tc.s\tresults\n",                                       // six columns
		"\ta\tengine\tfx\tc.s\tresults\tQ\n",                                     // an empty cell
		"a\ta\tengine\tfx\tc.s\tresults\t \n",                                    // an empty statement
		"a\ta\tengine\tfx\tc.s\tresults\tQ\na\ta\tengine\tfx\tc.s\tresults\tQ\n", // a repeat
	} {
		if _, err := Parse("span.tsv", []byte(header+row)); err == nil {
			t.Errorf("Parse took %q", row)
		}
	}
}

func TestGuard(t *testing.T) {
	parse := func(text string) *File {
		f, err := Parse("span.tsv", []byte(text))
		if err != nil {
			t.Fatal(err)
		}
		return f
	}
	agreed := parse(header + "a\ta\tengine\tfx\tc.s\tresults\tQ\n")
	differed := parse(header + "b\ta\tengine\tfx\tc.s\tresults\tQ\n")
	if p := Guard(agreed, differed); len(p) != 1 || !strings.Contains(p[0], "a regression") {
		t.Errorf("a committed regression: %q", p)
	}
	if p := Guard(differed, agreed); len(p) != 0 {
		t.Errorf("a closed gap: %q", p)
	}
	moved := parse(header + "a\tb\tengine\tfx\tc.s\tresults\tQ\n")
	if p := Guard(agreed, moved); len(p) == 0 {
		t.Errorf("an expected cell changed without a pin bump")
	}
	repinned := parse(strings.Replace(header, "v1", "v2", 1) + "b\tb\tengine\tfx\tc.s\tresults\tQ\n")
	if p := Guard(agreed, repinned); len(p) != 0 {
		t.Errorf("a pin bump's new expected cell: %q", p)
	}
	unsaid := parse(header + "c.t.a\t~\tengine\tfx\tc.s\tlineage\tQ\n")
	if p := Guard(unsaid, parse(header+"-\t~\tengine\tfx\tc.s\tlineage\tQ\n")); len(p) != 0 {
		t.Errorf("a move where the engine states nothing: %q", p)
	}
	stated := parse(header + "c.t.a c.t.b\tc.t.a ~\tengine\tfx\tc.s\tlineage\tQ\n")
	if p := Guard(stated, parse(header+"c.t.b\tc.t.a ~\tengine\tfx\tc.s\tlineage\tQ\n")); len(p) != 1 || !strings.Contains(p[0], "a regression") {
		t.Errorf("losing a stated source: %q", p)
	}
}

func TestAgrees(t *testing.T) {
	for _, c := range []struct {
		aspect, omni, expected string
		want                   bool
	}{
		{Results, "a, b", "a, b", true},
		{Results, "a, b", "a, c", false},
		{Lineage, "x.t.a | x.t.b", "x.t.a | x.t.b", true},
		{Lineage, "x.t.a | x.t.b", "x.t.a | ~", true},
		{Lineage, "x.t.a | -", "x.t.a | ~", true},
		{Lineage, "x.t.a | x.t.b x.v.b", "x.t.a | x.v.b ~", true},
		{Lineage, "x.t.a | x.t.b", "x.t.a | x.v.b ~", false},
		{Lineage, "x.t.a | x.t.b", "x.t.a | x.t.c", false},
		{Lineage, "x.t.a", "x.t.a | ~", false},
		{Lineage, "-", "~", true},
		{Lineage, "-", "x.t.a | ~", false},
		{Lineage, "x.t.a", "~", true},
		{Tables, "x.t, x.v", "x.t, ~", true},
		{Tables, "x.v", "x.t, ~", false},
		{Tables, "x.t, x.v", "x.t", false},
		{Tables, "-", "~", true},
		{Tables, "x.t", "error:1146", false},
		{Lineage, "error:omni", "~", false},
	} {
		if got := Agrees(c.aspect, c.omni, c.expected); got != c.want {
			t.Errorf("Agrees(%s, %q, %q) = %v, want %v", c.aspect, c.omni, c.expected, got, c.want)
		}
	}
}

func TestCheckHoldsAnUnstatedSourceAsAgreement(t *testing.T) {
	engine := Answer{Results: []string{"a", "_col1"}, Lineage: [][]string{{"c.s.t.a"}, {Unstated}}, Tables: []string{"c.s.t"}}
	if rec := run(t, suite(t, rows(answer, engine), engine, answer, nil)); !rec.passed() {
		t.Errorf("want a pass, got %q %q", rec.errors, rec.fatal)
	}
	// omni moving inside what the engine leaves unstated still agrees, but
	// is committed, not passed.
	moved := Answer{Results: []string{"a", "_col1"}, Lineage: [][]string{{"c.s.t.a"}, {}}, Tables: []string{"c.s.t"}}
	rec := run(t, suite(t, rows(answer, engine), engine, moved, nil))
	if !rec.failed("states nothing, and still agrees") || rec.failed("a regression") {
		t.Errorf("want the move reported for -update-compat, got %q %q", rec.errors, rec.fatal)
	}
	// Losing a source the engine does state is a regression.
	lost := Answer{Results: []string{"a", "_col1"}, Lineage: [][]string{{}, {"c.s.t.a"}}, Tables: []string{"c.s.t"}}
	if rec := run(t, suite(t, rows(answer, engine), engine, lost, nil)); !rec.failed("a regression") {
		t.Errorf("want a regression, got %q %q", rec.errors, rec.fatal)
	}
}
