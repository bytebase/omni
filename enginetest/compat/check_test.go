package compat

import (
	"errors"
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

// run checks s against a recorder. Check runs on its own goroutine because
// Fatalf ends the goroutine it runs on.
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

const testPin = "default\timg:1\tsha256:abc\tv1"

// suite writes pins.tsv and compat.tsv into a fresh directory. omni maps a
// statement to omni's verdict; a statement it lacks is rejected.
func suite(t *testing.T, rows string, omni map[string]string) Suite {
	t.Helper()
	dir := t.TempDir()
	write(t, filepath.Join(dir, "pins.tsv"), testPin+"\n")
	write(t, filepath.Join(dir, "compat.tsv"), "# pin: default img:1 sha256:abc v1\n"+columnsHeader+"\n"+rows)
	return Suite{
		Path:   filepath.Join(dir, "compat.tsv"),
		Pins:   filepath.Join(dir, "pins.tsv"),
		Source: SourceEngine,
		Omni: func(sql string) string {
			if v, ok := omni[sql]; ok {
				return v
			}
			return Reject
		},
		Judge: func(testing.TB, enginetest.Pin) func(string) (string, error) {
			t.Fatalf("the judge was asked in a replay")
			return nil
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

func TestCheckPassesAgreementsAndKnownGaps(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tr/A\tA\nreject\taccept\tengine\tr/B\tB\naccept\treject\tengine\tr/C\tC\n",
		map[string]string{"A": Accept, "C": Accept})
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("want a pass, got %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckFailsAClosedGapUntilUpdated(t *testing.T) {
	s := suite(t, "reject\taccept\tengine\tr/B\tB\n", map[string]string{"B": Accept})
	if rec := run(t, s); !rec.failed("a gap closed") {
		t.Errorf("a closed gap passed: %q", rec.errors)
	}

	*updateCompat = true
	defer func() { *updateCompat = false }()
	if rec := run(t, s); len(rec.errors) > 0 {
		t.Errorf("-update-compat: %q", rec.errors)
	}
	if got := read(t, s.Path); !strings.Contains(got, "accept\taccept\tengine\tr/B\tB\n") {
		t.Errorf("-update-compat did not commit the closed gap:\n%s", got)
	}
}

func TestCheckFailsARegressionAndNeverWritesIt(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", nil)
	before := read(t, s.Path)
	*updateCompat = true
	defer func() { *updateCompat = false }()
	if rec := run(t, s); !rec.failed("a regression") {
		t.Errorf("a regression passed: %q", rec.errors)
	}
	if read(t, s.Path) != before {
		t.Errorf("-update-compat wrote a regression")
	}
}

func TestCheckFailsAnUnjudgedRow(t *testing.T) {
	s := suite(t, "?\t?\t-\tissue:BYT-1\tA\n", nil)
	if rec := run(t, s); !rec.failed("not judged yet") {
		t.Errorf("an unjudged row passed: %q", rec.errors)
	}
}

func TestCheckPassesUndocumentedWhateverOmniSays(t *testing.T) {
	s := suite(t, "reject\tundocumented\tdocs\tr/A\tA\n", map[string]string{"A": Accept})
	if rec := run(t, s); len(rec.errors) > 0 {
		t.Errorf("an undocumented row failed: %q", rec.errors)
	}
}

func TestCheckFailsAnotherPin(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	write(t, s.Pins, "default\timg:2\tsha256:def\tv2\n")
	if rec := run(t, s); !rec.failed("was recorded at pin") {
		t.Errorf("a file recorded at another pin passed: %q", rec.errors)
	}
}

// judgeOf answers from a map, and counts what it was asked.
func judgeOf(answers map[string]string, asked *[]string) func(testing.TB, enginetest.Pin) func(string) (string, error) {
	return func(testing.TB, enginetest.Pin) func(string) (string, error) {
		return func(sql string) (string, error) {
			*asked = append(*asked, sql)
			if v, ok := answers[sql]; ok {
				return v, nil
			}
			return "", errors.New("connection lost")
		}
	}
}

func TestRecordingJudgesOnlyUnjudgedRows(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n?\t?\t-\tissue:BYT-1\tB\n", map[string]string{"A": Accept, "B": Accept})
	var asked []string
	s.Judge = judgeOf(map[string]string{"B": Reject}, &asked)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("recording: %q %q", rec.errors, rec.fatal)
	}
	if strings.Join(asked, ",") != "B" {
		t.Errorf("the judge was asked %q, want only the unjudged B", asked)
	}
	if got := read(t, s.Path); !strings.Contains(got, "accept\treject\tengine\tissue:BYT-1\tB\n") {
		t.Errorf("the recording did not fill the row:\n%s", got)
	}
}

func TestRecordingAtANewPinJudgesEveryRow(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	write(t, s.Pins, "default\timg:2\tsha256:def\tv2\n")
	var asked []string
	s.Judge = judgeOf(map[string]string{"A": Reject}, &asked)
	// omni still accepts and the new server rejects: a known gap, which passes.
	if rec := run(t, s); len(rec.errors) > 0 {
		t.Errorf("recording at a new pin: %q", rec.errors)
	}
	got := read(t, s.Path)
	if !strings.HasPrefix(got, "# pin: default img:2 sha256:def v2\n") || !strings.Contains(got, "accept\treject\tengine\tr/A\tA\n") {
		t.Errorf("the recording did not rejudge at the new pin:\n%s", got)
	}
}

func TestRecordingKeepsWhatTheJudgeLeftUnanswered(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	s := suite(t, "?\t?\t-\tissue:BYT-1\tA\n?\t?\t-\tissue:BYT-2\tB\n", nil)
	var asked []string
	s.Judge = judgeOf(map[string]string{"B": Reject}, &asked)
	rec := run(t, s)
	if !rec.failed("gave no verdict") {
		t.Errorf("a lost verdict passed: %q", rec.errors)
	}
	got := read(t, s.Path)
	if !strings.Contains(got, "reject\t?\t-\tissue:BYT-1\tA\n") || !strings.Contains(got, "reject\treject\tengine\tissue:BYT-2\tB\n") {
		t.Errorf("want A unjudged and B judged:\n%s", got)
	}
}

func TestRecordingFailsInCI(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "true")
	s := suite(t, "?\t?\t-\tissue:BYT-1\tA\n", nil)
	if rec := run(t, s); !rec.failed("CI") {
		t.Errorf("a recording ran in CI: %q %q", rec.errors, rec.fatal)
	}
}

const tinyGrammar = `
%token	identifier "identifier"
%token	drop "DROP"
%token	table "TABLE"
%start	Stmt
%%
Stmt:
	DropStmt
DropStmt:
	"DROP" "TABLE" identifier
%%
`

// A recording starts an engine's file; a replay without one says how.
func TestRecordingStartsAMissingFile(t *testing.T) {
	s := suite(t, "", nil)
	dir := filepath.Dir(s.Path)
	s.Grammar, s.Tokens = filepath.Join(dir, "g.y"), filepath.Join(dir, "tokens.tsv")
	write(t, s.Grammar, tinyGrammar)
	write(t, s.Tokens, "identifier\tt\n")
	if err := os.Remove(s.Path); err != nil {
		t.Fatal(err)
	}
	s.Path = filepath.Join(dir, "testdata", "compat.tsv") // in a directory not made yet
	if rec := run(t, s); !strings.Contains(rec.fatal, "Record it") {
		t.Errorf("a replay with no file: %q %q", rec.errors, rec.fatal)
	}

	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	var asked []string
	s.Judge = judgeOf(map[string]string{"DROP TABLE t": Accept, "DROP TABLE": Reject, "TABLE t": Reject, "DROP t": Reject, "DROP TABLE t zzz": Reject, "TABLE TABLE t": Reject}, &asked)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	got := read(t, s.Path)
	if !strings.HasPrefix(got, "# pin: default img:1 sha256:abc v1\n"+columnsHeader+"\n") || !strings.Contains(got, "reject\taccept\tengine\tStmt/DropStmt\tDROP TABLE t\n") {
		t.Errorf("the recording wrote:\n%s", got)
	}
}

func TestCheckHoldsGeneratedRowsToTheGrammar(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tDropStmt/DROP_TABLE_identifier\tDROP TABLE t\n"+
		"reject\treject\tengine\tsuite:corpus\tDROP\n", map[string]string{"DROP TABLE t": Accept})
	dir := filepath.Dir(s.Path)
	s.Grammar, s.Tokens = filepath.Join(dir, "g.y"), filepath.Join(dir, "tokens.tsv")
	write(t, s.Grammar, tinyGrammar)
	write(t, s.Tokens, "identifier\tt\n")

	rec := run(t, s)
	if !rec.failed("not what") {
		t.Fatalf("missing generated rows passed: %q", rec.errors)
	}

	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	var asked []string
	s.Judge = judgeOf(map[string]string{"DROP TABLE t": Accept, "DROP TABLE": Reject, "TABLE t": Reject, "DROP t": Reject, "DROP TABLE t zzz": Reject, "TABLE TABLE t": Reject}, &asked)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	f, err := ReadFile(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	var rules []string
	for _, r := range f.Rows {
		rules = append(rules, r.Rule)
	}
	// DropStmt's near-misses are Stmt's again, so only Stmt keeps them.
	want := "Stmt/DropStmt,Stmt/DropStmt!drop-1,Stmt/DropStmt!drop-2,Stmt/DropStmt!drop-3,Stmt/DropStmt!insert,Stmt/DropStmt!swap," +
		"DropStmt/DROP_TABLE_identifier,suite:corpus"
	if got := strings.Join(rules, ","); got != want {
		t.Errorf("rows:\n got %s\nwant %s", got, want)
	}

	t.Setenv(enginetest.RecordEnv, "")
	if rec := run(t, s); len(rec.errors) > 0 {
		t.Errorf("replay after recording: %q", rec.errors)
	}
}

func TestReadFileRejectsMalformedRows(t *testing.T) {
	for _, c := range []struct{ name, rows, want string }{
		{"four columns", "accept\taccept\tengine\tr/A\n", "want 5"},
		{"bad verdict", "yes\taccept\tengine\tr/A\tA\n", "omni"},
		{"bad position", "reject@x\taccept\tengine\tr/A\tA\n", "omni"},
		{"undocumented omni", "undocumented\taccept\tengine\tr/A\tA\n", "judge's answer"},
		{"judged without source", "accept\taccept\t-\tr/A\tA\n", "no source"},
		{"source without verdict", "accept\t?\tengine\tr/A\tA\n", "no expected"},
		{"bad source", "accept\taccept\tserver\tr/A\tA\n", "source"},
		{"space in rule", "accept\taccept\tengine\tr A\tA\n", "space"},
		{"repeated row", "accept\taccept\tengine\tr/A\tA\naccept\taccept\tengine\tr/A\tA\n", "repeats line"},
		{"comment among rows", "accept\taccept\tengine\tr/A\tA\n# late\n", "comment after"},
	} {
		path := filepath.Join(t.TempDir(), "compat.tsv")
		write(t, path, c.rows)
		if _, err := ReadFile(path); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want an error with %q", c.name, err, c.want)
		}
	}
}

func TestReadFileKeepsEveryStatementByte(t *testing.T) {
	path := filepath.Join(t.TempDir(), "compat.tsv")
	stmts := []string{"", "SELECT 1 ", "SELECT 'a\rb'", "SELECT 'a\tb'", "SELECT '\xff'"}
	f := &File{Pin: "default img:1 sha256:abc v1"}
	for i, s := range stmts {
		f.Rows = append(f.Rows, Row{Omni: Reject, Expected: "reject@3", Source: SourceEngine, Rule: fmt.Sprintf("r/%d", i), Statement: s})
	}
	if err := f.Write(path); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i, r := range got.Rows {
		if r.Statement != stmts[i] {
			t.Errorf("row %d: got %q, want %q", i, r.Statement, stmts[i])
		}
	}
}
