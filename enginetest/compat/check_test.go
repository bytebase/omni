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

// TestCheckCommitsAGapThatIsNowAnotherOne covers a judge that says where it
// rejects: a row omni got wrong can move to another wrong verdict, which
// -update-compat commits, while a row that agreed cannot.
func TestCheckCommitsAGapThatIsNowAnotherOne(t *testing.T) {
	s := suite(t, "accept\treject@7\tengine\tr/A\tA\nreject@7\treject@7\tengine\tr/B\tB\n",
		map[string]string{"A": "reject@3", "B": "reject@3"})
	rec := run(t, s)
	if !rec.failed("a known gap is now another one") {
		t.Errorf("a gap that moved passed: %q", rec.errors)
	}
	if !rec.failed("a regression") {
		t.Errorf("an agreement that moved passed: %q", rec.errors)
	}

	*updateCompat = true
	defer func() { *updateCompat = false }()
	if rec := run(t, s); !rec.failed("a regression") || rec.failed("a known gap") {
		t.Errorf("-update-compat: %q", rec.errors)
	}
	if got := read(t, s.Path); !strings.Contains(got, "reject@3\treject@7\tengine\tr/A\tA\nreject@7\treject@7\tengine\tr/B\tB\n") {
		t.Errorf("-update-compat did not commit the moved gap alone:\n%s", got)
	}
}

func TestCheckFailsAnUnjudgedRow(t *testing.T) {
	s := suite(t, "?\t?\t-\tissue:BYT-1\tA\n", nil)
	if rec := run(t, s); !rec.failed("not judged yet") {
		t.Errorf("an unjudged row passed: %q", rec.errors)
	}
}

// docsGrammar documents DROP TABLE [IF EXISTS] t and nothing else.
const docsGrammar = `grammar D;
stmt : 'DROP' 'TABLE' ('IF' 'EXISTS')? NAME ;
`

// docsSuite is suite with the docs judge over docsGrammar: no Judge, since
// the recognizer is the judge.
func docsSuite(t *testing.T, rows string, omni map[string]string) Suite {
	t.Helper()
	s := suite(t, rows, omni)
	dir := filepath.Dir(s.Path)
	s.Grammar, s.Tokens = filepath.Join(dir, "docs.g4"), filepath.Join(dir, "tokens.tsv")
	write(t, s.Grammar, docsGrammar)
	write(t, s.Tokens, "NAME\ta\n")
	s.Source, s.Judge = SourceDocs, nil
	return s
}

// The generated rows of docsGrammar, judged by the docs. DROP TABLE IF is
// documented, with IF as the name; DROP TABLE EXISTS a is not.
const docsRows = "" +
	"accept\taccept\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME\tDROP TABLE a\n" +
	"reject\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!drop-1\tTABLE a\n" +
	"accept\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!drop-2\tDROP a\n" +
	"accept\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!drop-3\tDROP TABLE\n" +
	"reject\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!insert\tDROP TABLE a zzz\n" +
	"reject\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!swap\tEXISTS TABLE a\n" +
	"accept\taccept\tdocs\tstmt(IF_EXISTS)?/%empty\tDROP TABLE a\n" +
	"reject\tundocumented\tdocs\tstmt(IF_EXISTS)?/%empty!insert\tDROP TABLE zzz a\n" +
	"accept\taccept\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS\tDROP TABLE IF EXISTS a\n" +
	"accept\tundocumented\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-1\tDROP TABLE EXISTS a\n" +
	"reject\tundocumented\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2\tDROP TABLE IF a\n" +
	"reject\taccept\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!truncate\tDROP TABLE IF\n" +
	"reject\tundocumented\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!insert\tDROP TABLE IF EXISTS zzz a\n" +
	"reject\tundocumented\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!swap\tDROP TABLE TABLE EXISTS a\n"

var docsOmni = map[string]string{"DROP TABLE a": Accept, "DROP TABLE IF EXISTS a": Accept, "DROP TABLE EXISTS a": Accept, "DROP a": Accept, "DROP TABLE": Accept}

func TestCheckPassesUndocumentedWhateverOmniSays(t *testing.T) {
	s := docsSuite(t, docsRows, docsOmni)
	// omni now rejects DROP a and accepts DROP TABLE a zzz: neither is
	// documented, so neither fails.
	omni := map[string]string{"DROP TABLE a zzz": Accept}
	for k, v := range docsOmni {
		omni[k] = v
	}
	delete(omni, "DROP a")
	s.Omni = func(sql string) string {
		if v, ok := omni[sql]; ok {
			return v
		}
		return Reject
	}
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("an undocumented row failed: %q %q", rec.errors, rec.fatal)
	}
	*updateCompat = true
	defer func() { *updateCompat = false }()
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("-update-compat: %q %q", rec.errors, rec.fatal)
	}
	got := read(t, s.Path)
	if !strings.Contains(got, "reject\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!drop-2\tDROP a\n") ||
		!strings.Contains(got, "accept\tundocumented\tdocs\tstmt/DROP_TABLE_(IF_EXISTS)?_NAME!insert\tDROP TABLE a zzz\n") {
		t.Errorf("-update-compat did not keep the undocumented rows' omni cells current:\n%s", got)
	}
}

// TestDocsRecordingJudgesByTheRecognizer pins the docs judge: the grammar's
// own statements are accepted, a near-miss the docs derive (DROP TABLE IF,
// with IF as the name) is accepted, and the rest are undocumented. It needs
// no server.
func TestDocsRecordingJudgesByTheRecognizer(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	s := docsSuite(t, "", docsOmni)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); got != "# pin: default img:1 sha256:abc v1\n"+columnsHeader+"\n"+docsRows {
		t.Errorf("the recording wrote:\n%s\nwant:\n%s", got, docsRows)
	}
	t.Setenv(enginetest.RecordEnv, "")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("replay after recording: %q %q", rec.errors, rec.fatal)
	}
}

// TestDocsCellsAreJudgedAgainOnEveryRun pins that a docs cell is the
// recognizer's: a replay fails one it no longer gives, after a grammar fix
// or an edit by hand, and a recording moves it.
func TestDocsCellsAreJudgedAgainOnEveryRun(t *testing.T) {
	edited := strings.Replace(docsRows, "reject\tundocumented\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2", "reject\treject\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2", 1)
	if edited == docsRows {
		t.Fatal("no row edited")
	}
	s := docsSuite(t, edited, docsOmni)
	if rec := run(t, s); !rec.failed("the docs judge now says undocumented, the row records reject") {
		t.Errorf("an edited docs cell passed: %q %q", rec.errors, rec.fatal)
	}
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); !strings.HasSuffix(got, docsRows) {
		t.Errorf("the recording did not move the cell back:\n%s", got)
	}
}

func TestDocsSuiteTakesNoJudge(t *testing.T) {
	s := docsSuite(t, docsRows, docsOmni)
	var asked []string
	s.Judge = judgeOf(nil, &asked)
	if rec := run(t, s); !strings.Contains(rec.fatal, "give the Grammar and no Judge") {
		t.Errorf("a docs suite with a Judge: %q %q", rec.errors, rec.fatal)
	}
}

func TestCheckFailsAnotherPin(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	write(t, s.Pins, "default\timg:2\tsha256:def\tv2\n")
	if rec := run(t, s); !rec.failed("was recorded at pin") {
		t.Errorf("a file recorded at another pin passed: %q", rec.errors)
	}
}

func TestCheckUsesPinRole(t *testing.T) {
	s := suite(t, "", map[string]string{"A": Accept})
	write(t, s.Pins, testPin+"\nmongosh\tgithub.com/mongodb-js/mongosh@abc\tsha256:def\tv2.0.0\n")
	write(t, s.Path, "# pin: mongosh github.com/mongodb-js/mongosh@abc sha256:def v2.0.0\n"+columnsHeader+"\naccept\taccept\tdocs\tr/A\tA\n")
	s.PinRole = "mongosh"

	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("replay with PinRole=mongosh: %q %q", rec.errors, rec.fatal)
	}

	// Changing the default row does not affect a PinRole=mongosh suite.
	write(t, s.Pins, "default\timg:2\tsha256:999\tv2\nmongosh\tgithub.com/mongodb-js/mongosh@abc\tsha256:def\tv2.0.0\n")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("changing default pin broke PinRole=mongosh: %q %q", rec.errors, rec.fatal)
	}

	// Changing the mongosh row fails replay with a pin mismatch.
	write(t, s.Pins, testPin+"\nmongosh\tgithub.com/mongodb-js/mongosh@xyz\tsha256:999\tv2.1.0\n")
	if rec := run(t, s); !rec.failed("was recorded at pin") {
		t.Errorf("changed PinRole row passed: %q %q", rec.errors, rec.fatal)
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

func TestCheckFailsAnotherJudgeRevision(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	s.JudgeRevision = 1
	if rec := run(t, s); !rec.failed("was recorded at pin") {
		t.Errorf("a file recorded by an earlier judge passed: %q", rec.errors)
	}
}

func TestRecordingAtANewJudgeRevisionJudgesEveryRow(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	s.JudgeRevision = 2
	var asked []string
	s.Judge = judgeOf(map[string]string{"A": Reject}, &asked)
	if rec := run(t, s); len(rec.errors) > 0 {
		t.Errorf("recording at a new judge revision: %q", rec.errors)
	}
	got := read(t, s.Path)
	if !strings.HasPrefix(got, "# pin: default img:1 sha256:abc v1 judge 2\n") || !strings.Contains(got, "accept\treject\tengine\tr/A\tA\n") {
		t.Errorf("the recording did not rejudge at the new revision:\n%s", got)
	}
	// The file now passes.
	t.Setenv(enginetest.RecordEnv, "")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("after the recording: %q %q", rec.errors, rec.fatal)
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

// TestCheckJoinsTheRowsOfMoreGrammars: a second parser grammar read with the
// same lexer adds the rows of its own rules, and a row both generate fails.
func TestCheckJoinsTheRowsOfMoreGrammars(t *testing.T) {
	s := suite(t, "", nil)
	dir := filepath.Dir(s.Path)
	s.Grammar, s.Lexer, s.Tokens = filepath.Join(dir, "P.g4"), filepath.Join(dir, "L.g4"), filepath.Join(dir, "tokens.tsv")
	s.More = []string{filepath.Join(dir, "Q.g4")}
	write(t, s.Grammar, "parser grammar P;\nstmt : DROP TABLE NAME ;\n")
	write(t, s.More[0], "parser grammar Q;\nlegacy : DROP VIEW NAME ;\n")
	write(t, s.Lexer, "lexer grammar L;\nDROP : 'DROP' ;\nTABLE : 'TABLE' ;\nVIEW : 'VIEW' ;\nNAME : [a-z]+ ;\nWS : ' '+ -> channel(2) ;\n")
	write(t, s.Tokens, "NAME\tt\n")

	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
	s.Judge = func(testing.TB, enginetest.Pin) func(string) (string, error) {
		return func(sql string) (string, error) { return Reject, nil }
	}
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	got := read(t, s.Path)
	for _, want := range []string{"\tstmt/DROP_TABLE_NAME\tDROP TABLE t\n", "\tlegacy/DROP_VIEW_NAME\tDROP VIEW t\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the recording holds no row %q:\n%s", want, got)
		}
	}

	write(t, s.More[0], "parser grammar Q;\nstmt : DROP TABLE NAME ;\n")
	if rec := run(t, s); !strings.Contains(rec.fatal, "generate a row twice") {
		t.Errorf("a row two grammars generate: %q %q", rec.errors, rec.fatal)
	}
}

func TestVerifyPassesOnMatchingFileAndWritesNothing(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "verify")
	t.Setenv("CI", "true")

	s := suite(t, "accept\taccept\tengine\tr/A\tA\nreject\taccept\tengine\tr/B\tB\naccept\treject\tengine\tr/C\tC\n",
		map[string]string{"A": Accept, "C": Accept})
	var asked []string
	s.Judge = judgeOf(map[string]string{"A": Accept, "B": Accept, "C": Reject}, &asked)
	before := read(t, s.Path)

	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("verify: %q %q", rec.errors, rec.fatal)
	}
	if got := strings.Join(asked, ","); got != "A,B,C" {
		t.Errorf("verify asked judge %q, want every row A,B,C", got)
	}
	if read(t, s.Path) != before {
		t.Errorf("verify modified %s", s.Path)
	}

	// Docs suite passes in verify mode too and writes nothing.
	ds := docsSuite(t, docsRows, docsOmni)
	beforeDocs := read(t, ds.Path)
	if rec := run(t, ds); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("docs verify: %q %q", rec.errors, rec.fatal)
	}
	if read(t, ds.Path) != beforeDocs {
		t.Errorf("docs verify modified %s", ds.Path)
	}
}

func TestVerifyFailsOnTamperedExpectedAndWritesNothing(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "verify")
	t.Setenv("CI", "true")

	// Hand-edited expected cell on engine suite.
	s := suite(t, "accept\taccept\tengine\tr/A\tA\nreject\treject\tengine\tr/B\tB\n",
		map[string]string{"A": Accept})
	var asked []string
	s.Judge = judgeOf(map[string]string{"A": Accept, "B": Accept}, &asked)
	before := read(t, s.Path)

	rec := run(t, s)
	if !rec.failed("the judge's answer differs from the committed expected cell") || !rec.failed("the engine judge now says accept, the row records reject") {
		t.Errorf("tampered expected cell passed verify: %q %q", rec.errors, rec.fatal)
	}
	if read(t, s.Path) != before {
		t.Errorf("verify modified %s on mismatch", s.Path)
	}

	// Hand-edited expected cell on docs suite.
	edited := strings.Replace(docsRows, "reject\tundocumented\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2", "reject\treject\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2", 1)
	ds := docsSuite(t, edited, docsOmni)
	beforeDocs := read(t, ds.Path)
	if rec := run(t, ds); !rec.failed("the docs judge now says undocumented, the row records reject") {
		t.Errorf("tampered docs cell passed verify: %q %q", rec.errors, rec.fatal)
	}
	if read(t, ds.Path) != beforeDocs {
		t.Errorf("docs verify modified %s on mismatch", ds.Path)
	}
}

func TestVerifyRejectsUpdateCompatPinMismatchMissingFileGeneratedDiffAndUnjudged(t *testing.T) {
	t.Setenv(enginetest.RecordEnv, "verify")
	t.Setenv("CI", "true")

	// -update-compat cannot be combined with OMNI_RECORD=verify.
	s := suite(t, "reject\taccept\tengine\tr/B\tB\n", map[string]string{"B": Accept})
	var asked []string
	s.Judge = judgeOf(map[string]string{"B": Accept}, &asked)
	before := read(t, s.Path)
	*updateCompat = true
	rec := run(t, s)
	*updateCompat = false
	if !strings.Contains(rec.fatal, "-update-compat cannot be combined with OMNI_RECORD=verify") {
		t.Errorf("-update-compat + verify: %q %q", rec.errors, rec.fatal)
	}
	if len(asked) != 0 || read(t, s.Path) != before {
		t.Errorf("-update-compat + verify asked %q or modified file", asked)
	}

	// Missing compat.tsv fails and creates nothing.
	missing := suite(t, "", nil)
	if err := os.Remove(missing.Path); err != nil {
		t.Fatal(err)
	}
	if rec := run(t, missing); !strings.Contains(rec.fatal, "no ") {
		t.Errorf("missing file: %q %q", rec.errors, rec.fatal)
	}
	if _, err := os.Stat(missing.Path); !os.IsNotExist(err) {
		t.Errorf("verify created missing %s", missing.Path)
	}

	// Pin mismatch fails without asking the judge and writes nothing.
	pinMismatch := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	write(t, pinMismatch.Pins, "default\timg:2\tsha256:def\tv2\n")
	asked = nil
	pinMismatch.Judge = judgeOf(map[string]string{"A": Accept}, &asked)
	beforePin := read(t, pinMismatch.Path)
	if rec := run(t, pinMismatch); !strings.Contains(rec.fatal, "was recorded at pin") {
		t.Errorf("pin mismatch: %q %q", rec.errors, rec.fatal)
	}
	if len(asked) != 0 || read(t, pinMismatch.Path) != beforePin {
		t.Errorf("pin mismatch asked %q or modified file", asked)
	}

	// Unjudged row fails without asking the judge and writes nothing.
	unjudged := suite(t, "?\t?\t-\tissue:BYT-1\tA\n", nil)
	asked = nil
	unjudged.Judge = judgeOf(map[string]string{"A": Reject}, &asked)
	beforeUnjudged := read(t, unjudged.Path)
	if rec := run(t, unjudged); !rec.failed("not judged yet") {
		t.Errorf("unjudged row: %q %q", rec.errors, rec.fatal)
	}
	if len(asked) != 0 || read(t, unjudged.Path) != beforeUnjudged {
		t.Errorf("unjudged row asked %q or modified file", asked)
	}

	// Generated rows mismatch fails and writes nothing.
	gen := suite(t, "accept\taccept\tengine\tDropStmt/DROP_TABLE_identifier\tDROP TABLE t\n", map[string]string{"DROP TABLE t": Accept})
	dir := filepath.Dir(gen.Path)
	gen.Grammar, gen.Tokens = filepath.Join(dir, "g.y"), filepath.Join(dir, "tokens.tsv")
	write(t, gen.Grammar, tinyGrammar)
	write(t, gen.Tokens, "identifier\tt\n")
	asked = nil
	gen.Judge = judgeOf(map[string]string{"DROP TABLE t": Accept}, &asked)
	beforeGen := read(t, gen.Path)
	if rec := run(t, gen); !rec.failed("not what") {
		t.Errorf("generated rows mismatch: %q %q", rec.errors, rec.fatal)
	}
	if read(t, gen.Path) != beforeGen {
		t.Errorf("verify modified %s on generated rows mismatch", gen.Path)
	}

	// Direct (*File).Write callers (such as fuzz-row tests with -update-compat)
	// are also blocked from writing under OMNI_RECORD=verify.
	directPath := filepath.Join(t.TempDir(), "compat.tsv")
	directFile := &File{Pin: "default img:1 sha256:abc v1"}
	if err := directFile.Write(directPath); err == nil || !strings.Contains(err.Error(), "OMNI_RECORD=verify") {
		t.Errorf("File.Write in verify mode: got %v, want error naming OMNI_RECORD=verify", err)
	}
	if _, err := os.Stat(directPath); !os.IsNotExist(err) {
		t.Errorf("File.Write in verify mode created %s", directPath)
	}
}

// TestCheckEverySiteRuleOfAnyGrammar pins a tokens.tsv shared by more than
// one grammar, as Doris's is: an (every site) row names a rule of one of
// them, and a rule none of them has fails.
func TestCheckEverySiteRuleOfAnyGrammar(t *testing.T) {
	s := suite(t, "", nil)
	dir := filepath.Dir(s.Path)
	s.Grammar, s.Tokens = filepath.Join(dir, "one.g4"), filepath.Join(dir, "tokens.tsv")
	s.More = []string{filepath.Join(dir, "two.g4")}
	write(t, s.Grammar, "grammar One;\nstmt : 'DROP' 'TABLE' ('IF' 'EXISTS')? NAME ;\n")
	write(t, s.More[0], "grammar Two;\nshow : 'SHOW' opt_full 'TABLES' ;\nopt_full : | 'FULL' ;\n")
	write(t, s.Tokens, "NAME\ta\nstmt(IF_EXISTS)?\t(every site)\nopt_full\t(every site)\n")
	// The two grammars generate rows the empty file lacks, an error and no
	// fatal.
	if rec := run(t, s); rec.fatal != "" || !rec.failed("to add") {
		t.Errorf("a rule of one grammar: %q, %q", rec.fatal, rec.errors)
	}
	write(t, s.Tokens, "NAME\ta\nopt_tmp\t(every site)\n")
	if rec := run(t, s); !rec.failed("opt_tmp the form (every site), which only a rule takes, and no grammar of the suite has") {
		t.Errorf("a rule of no grammar: %q, %q", rec.fatal, rec.errors)
	}
}
