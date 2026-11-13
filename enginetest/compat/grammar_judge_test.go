package compat

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

// grammarSuite is docsSuite judged as a vendor's grammar is: what the
// grammar does not derive is rejected, where docs would leave it
// undocumented.
func grammarSuite(t *testing.T, rows string, omni map[string]string) Suite {
	t.Helper()
	s := docsSuite(t, rows, omni)
	s.Source = SourceGrammar
	return s
}

// grammarRows is docsRows as the grammar judges them.
var grammarRows = strings.NewReplacer("\tundocumented\tdocs\t", "\treject\tgrammar\t", "\tdocs\t", "\tgrammar\t").Replace(docsRows)

func record(t *testing.T) {
	t.Helper()
	t.Setenv(enginetest.RecordEnv, "1")
	t.Setenv("CI", "")
}

// TestGrammarRecordingJudgesByTheRecognizer pins the grammar judge end to
// end: a recording with no server writes accept for what the grammar
// derives and reject for the rest, and the replay holds omni to both.
func TestGrammarRecordingJudgesByTheRecognizer(t *testing.T) {
	record(t)
	s := grammarSuite(t, "", docsOmni)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); got != "# pin: default img:1 sha256:abc v1\n"+columnsHeader+"\n"+grammarRows {
		t.Errorf("the recording wrote:\n%s\nwant:\n%s", got, grammarRows)
	}
	t.Setenv(enginetest.RecordEnv, "")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("replay after recording: %q %q", rec.errors, rec.fatal)
	}
}

// A rejection is a verdict, where undocumented is none: omni moving onto it
// closes a gap, and omni moving off it is a regression.
func TestGrammarRejectionsHoldOmni(t *testing.T) {
	// The rows record that omni accepts DROP a, which the grammar rejects.
	omni := map[string]string{}
	for k, v := range docsOmni {
		omni[k] = v
	}
	delete(omni, "DROP a")
	if rec := run(t, grammarSuite(t, grammarRows, omni)); !rec.failed("a gap closed") {
		t.Errorf("omni now rejects what the grammar rejects, and nothing asked for the row to be updated: %q %q", rec.errors, rec.fatal)
	}
	omni["DROP a"] = Accept
	omni["DROP TABLE a zzz"] = Accept
	if rec := run(t, grammarSuite(t, grammarRows, omni)); !rec.failed("a regression") {
		t.Errorf("omni now accepts what it and the grammar rejected: %q %q", rec.errors, rec.fatal)
	}
}

// TestGrammarCellsAreJudgedAgainOnEveryRun: a grammar cell is the
// recognizer's, so one written by hand fails with no server.
func TestGrammarCellsAreJudgedAgainOnEveryRun(t *testing.T) {
	edited := strings.Replace(grammarRows, "reject\treject\tgrammar\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2", "reject\taccept\tgrammar\tstmt(IF_EXISTS)?/IF_EXISTS!drop-2", 1)
	if edited == grammarRows {
		t.Fatal("no row edited")
	}
	if rec := run(t, grammarSuite(t, edited, docsOmni)); !rec.failed("the grammar judge now says reject, the row records accept") {
		t.Errorf("an edited grammar cell passed: %q %q", rec.errors, rec.fatal)
	}
}

// overrideSuite is grammarSuite with a grammar-overrides.tsv.
func overrideSuite(t *testing.T, rows, overrides string) Suite {
	t.Helper()
	s := grammarSuite(t, rows, docsOmni)
	s.Overrides = filepath.Join(filepath.Dir(s.Path), "grammar-overrides.tsv")
	write(t, s.Overrides, overrides)
	return s
}

// The engine refuses DROP TABLE IF in an action the grammar's rules do not
// show, and takes DROP TABLE EXISTS a in its lexer.
const overrides = "# columns: verdict, source, cite, statement\n" +
	"reject\tgrammar\tgram.y:10 an action refuses IF as a name\tDROP TABLE IF\n" +
	"accept\tdocs\thttps://example.com/drop-table: EXISTS before a name is a name\tDROP TABLE EXISTS a\n"

// overriddenRows is grammarRows with the two statements as the engine
// judges them.
var overriddenRows = strings.NewReplacer(
	"reject\taccept\tgrammar\tstmt(IF_EXISTS)?/IF_EXISTS!truncate\tDROP TABLE IF\n",
	"reject\treject\tgrammar\tstmt(IF_EXISTS)?/IF_EXISTS!truncate\tDROP TABLE IF\n",
	"accept\treject\tgrammar\tstmt(IF_EXISTS)?/IF_EXISTS!drop-1\tDROP TABLE EXISTS a\n",
	"accept\taccept\tdocs\tstmt(IF_EXISTS)?/IF_EXISTS!drop-1\tDROP TABLE EXISTS a\n",
).Replace(grammarRows)

// TestOverridesCorrectTheRecognizer pins grammar-overrides.tsv: a recording
// writes the override's verdict for its statement, in either direction, a
// replay holds the cell to it, and without the file the cell fails.
func TestOverridesCorrectTheRecognizer(t *testing.T) {
	if overriddenRows == grammarRows {
		t.Fatal("no row overridden")
	}
	record(t)
	s := overrideSuite(t, "", overrides)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); !strings.HasSuffix(got, columnsHeader+"\n"+overriddenRows) {
		t.Errorf("the recording wrote:\n%s\nwant:\n%s", got, overriddenRows)
	}
	t.Setenv(enginetest.RecordEnv, "")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("replay after recording: %q %q", rec.errors, rec.fatal)
	}
	s.Overrides = ""
	if rec := run(t, s); !rec.failed("the grammar judge now says accept, the row records reject (grammar)") ||
		!rec.failed("the grammar judge now says reject, the row records accept (docs)") {
		t.Errorf("a cell only an override gives passed without it: %q %q", rec.errors, rec.fatal)
	}
}

// An override must be a miss of the recognizer's on a row of the file, and
// must say what decides it.
func TestOverridesThatCorrectNothingFail(t *testing.T) {
	for name, tc := range map[string]struct{ overrides, want string }{
		"no row":       {overrides + "reject\tgrammar\tgram.y:1\tDROP TABLE nowhere\n", "no row of the file holds the statement"},
		"no miss":      {overrides + "accept\tgrammar\tgram.y:1\tDROP TABLE a\n", "the recognizer already says accept"},
		"no cite":      {"reject\tgrammar\t \tDROP TABLE IF\n", "no cite"},
		"no verdict":   {"undocumented\tgrammar\tgram.y:1\tDROP TABLE IF\n", `verdict "undocumented"`},
		"docs reject":  {"reject\tdocs\thttps://example.com\tDROP TABLE IF\n", "a docs override only accepts"},
		"engine":       {"reject\tengine\tgram.y:1\tDROP TABLE IF\n", `source "engine"`},
		"twice":        {overrides + "accept\tgrammar\tgram.y:1\tDROP TABLE IF\n", "repeats line 2"},
		"two columns":  {"reject\tDROP TABLE IF\n", "want 4 tab-separated columns"},
		"missing file": {"", "no such file"},
	} {
		t.Run(name, func(t *testing.T) {
			s := overrideSuite(t, overriddenRows, tc.overrides)
			if name == "missing file" {
				s.Overrides += ".absent"
			}
			if rec := run(t, s); !rec.failed(tc.want) {
				t.Errorf("want a failure holding %q, got %q %q", tc.want, rec.errors, rec.fatal)
			}
		})
	}
}

// Only a recognizer is corrected: an engine's answer is the engine's.
func TestOverridesNeedARecognizer(t *testing.T) {
	s := suite(t, "accept\taccept\tengine\tr/A\tA\n", map[string]string{"A": Accept})
	s.Overrides = filepath.Join(filepath.Dir(s.Path), "grammar-overrides.tsv")
	write(t, s.Overrides, overrides)
	if rec := run(t, s); !strings.Contains(rec.fatal, "Overrides correct a recognizer") {
		t.Errorf("an engine suite with overrides: %q %q", rec.errors, rec.fatal)
	}
}

// An alternative the grammar derives and an action always refuses is
// overridden like any row: the recognizer's own complaint, that it does not
// derive the generator's statement, is about another failure.
func TestOverrideRejectsAGeneratedAlternative(t *testing.T) {
	record(t)
	s := overrideSuite(t, "", "reject\tgrammar\tgram.y:30 the action refuses every IF EXISTS\tDROP TABLE IF EXISTS a\n")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	if got := read(t, s.Path); !strings.Contains(got, "accept\treject\tgrammar\tstmt(IF_EXISTS)?/IF_EXISTS\tDROP TABLE IF EXISTS a\n") {
		t.Errorf("the alternative was not overridden:\n%s", got)
	}
}

// TestVendorFilesAreThePinned pins Suite.Vendor: the files a grammar judge
// reads hash to the pin's digest, so neither moves without the other.
func TestVendorFilesAreThePinned(t *testing.T) {
	s := grammarSuite(t, grammarRows, docsOmni)
	s.Vendor = []string{s.Grammar, s.Tokens}
	digest, err := enginetest.FilesDigest(s.Vendor...)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(docsGrammar + "NAME\ta\n"))
	if want := "sha256:" + hex.EncodeToString(sum[:]); digest != want {
		t.Fatalf("FilesDigest = %s, want the SHA-256 of the files' bytes one after the other, %s", digest, want)
	}
	if rec := run(t, s); !strings.Contains(rec.fatal, "a vendor file changes only with its pin") {
		t.Errorf("vendor files that do not hash to the pin's digest passed: %q %q", rec.errors, rec.fatal)
	}
	write(t, s.Pins, "default\timg:1\t"+digest+"\tv1\n")
	rows := "# pin: default img:1 " + digest + " v1\n" + columnsHeader + "\n" + grammarRows
	write(t, s.Path, rows)
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Errorf("vendor files at their pin: %q %q", rec.errors, rec.fatal)
	}
	write(t, s.Grammar, docsGrammar+"// edited\n")
	if rec := run(t, s); !strings.Contains(rec.fatal, "a vendor file changes only with its pin") {
		t.Errorf("an edited vendor file passed: %q %q", rec.errors, rec.fatal)
	}
}

// TestParserAndLexerGrammarsAreReadTogether pins Suite.Lexer: a parser
// grammar published apart from its lexer grammar is read with it, since the
// parser grammar alone does not spell its tokens.
func TestParserAndLexerGrammarsAreReadTogether(t *testing.T) {
	record(t)
	s := grammarSuite(t, "", map[string]string{"DROP a": Accept})
	dir := filepath.Dir(s.Path)
	s.Lexer = filepath.Join(dir, "Tokens.g4")
	write(t, s.Grammar, "parser grammar P;\noptions { tokenVocab = Tokens; }\nstmt : DROP TABLE? NAME ;\n")
	write(t, s.Lexer, "lexer grammar Tokens;\nDROP : 'DROP' ;\nTABLE : 'TABLE' ;\nNAME : [a-z]+ ;\n")
	if rec := run(t, s); len(rec.errors) > 0 || rec.fatal != "" {
		t.Fatalf("recording: %q %q", rec.errors, rec.fatal)
	}
	got := read(t, s.Path)
	for _, want := range []string{
		"accept\taccept\tgrammar\tstmt(TABLE)?/%empty\tDROP a\n",
		// The lexer grammar spells TABLE, and the swap picks among its keywords.
		"reject\taccept\tgrammar\tstmt(TABLE)?/TABLE\tDROP TABLE a\n",
		"reject\treject\tgrammar\tstmt(TABLE)?/TABLE!swap\tDROP DROP a\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the recording lacks %q:\n%s", want, got)
		}
	}

	s.Grammar = filepath.Join(dir, "g.y")
	write(t, s.Grammar, tinyGrammar)
	if rec := run(t, s); !strings.Contains(rec.fatal, "is not an ANTLR grammar") {
		t.Errorf("a bison grammar with a Lexer: %q %q", rec.errors, rec.fatal)
	}
}
