package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runParseBestEffortCase is a table-driven helper that asserts the shape
// of ParseBestEffort's result: expected statement count and expected
// number of errors (with optional per-error content matchers).
type parseCase struct {
	name        string
	input       string
	wantStmtCnt int      // expected len(result.File.Stmts)
	wantErrCnt  int      // expected len(result.Errors)
	wantErrMsgs []string // expected substring in each error's Msg (len must match wantErrCnt)
	wantErrLocs []int    // expected Loc.Start per error (len must match wantErrCnt); -1 to skip
}

func runParseBestEffortCase(t *testing.T, c parseCase) {
	t.Helper()
	result := ParseBestEffort(c.input)
	if got := len(result.File.Stmts); got != c.wantStmtCnt {
		t.Errorf("%s: stmt count = %d, want %d", c.name, got, c.wantStmtCnt)
	}
	if got := len(result.Errors); got != c.wantErrCnt {
		t.Errorf("%s: error count = %d, want %d", c.name, got, c.wantErrCnt)
		for i, e := range result.Errors {
			t.Errorf("  [%d] [%d,%d] %s", i, e.Loc.Start, e.Loc.End, e.Msg)
		}
		return
	}
	for i, want := range c.wantErrMsgs {
		if !strings.Contains(result.Errors[i].Msg, want) {
			t.Errorf("%s: error[%d] Msg = %q, want to contain %q",
				c.name, i, result.Errors[i].Msg, want)
		}
	}
	for i, wantLoc := range c.wantErrLocs {
		if wantLoc < 0 {
			continue
		}
		if result.Errors[i].Loc.Start != wantLoc {
			t.Errorf("%s: error[%d] Loc.Start = %d, want %d",
				c.name, i, result.Errors[i].Loc.Start, wantLoc)
		}
	}
}

func TestParse_EmptyInput(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "empty",
		input:       "",
		wantStmtCnt: 0,
		wantErrCnt:  0,
	})
}

func TestParse_WhitespaceOnly(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "whitespace only",
		input:       "   \n\t  ",
		wantStmtCnt: 0,
		wantErrCnt:  0,
	})
}

func TestParse_CommentOnly(t *testing.T) {
	cases := []parseCase{
		{name: "line comment", input: "-- comment\n", wantStmtCnt: 0, wantErrCnt: 0},
		{name: "block comment", input: "/* comment */", wantStmtCnt: 0, wantErrCnt: 0},
		{name: "slash line comment", input: "// comment\n", wantStmtCnt: 0, wantErrCnt: 0},
		{name: "nested block", input: "/* outer /* inner */ still */", wantStmtCnt: 0, wantErrCnt: 0},
	}
	for _, c := range cases {
		runParseBestEffortCase(t, c)
	}
}

func TestParse_SingleUnsupportedSelect(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "single SELECT",
		input:       "SELECT 1;",
		wantStmtCnt: 0,
		wantErrCnt:  1,
		wantErrMsgs: []string{"SELECT statement parsing is not yet supported"},
		wantErrLocs: []int{0},
	})
}

func TestParse_MultiUnsupported(t *testing.T) {
	// "SELECT 1; INSERT INTO t VALUES (1);" has SELECT at byte 0 and
	// INSERT at byte 10 (after "SELECT 1; ").
	runParseBestEffortCase(t, parseCase{
		name:        "SELECT then INSERT",
		input:       "SELECT 1; INSERT INTO t VALUES (1);",
		wantStmtCnt: 0,
		wantErrCnt:  2,
		wantErrMsgs: []string{
			"SELECT statement parsing is not yet supported",
			"INSERT statement parsing is not yet supported",
		},
		wantErrLocs: []int{0, 10},
	})
}

func TestParse_UnknownStatement(t *testing.T) {
	runParseBestEffortCase(t, parseCase{
		name:        "unknown keyword FOO",
		input:       "FOO BAR;",
		wantStmtCnt: 0,
		wantErrCnt:  1,
		wantErrMsgs: []string{"unknown or unsupported statement starting with FOO"},
		wantErrLocs: []int{0},
	})
}

func TestParse_LexErrorPropagated(t *testing.T) {
	// "'unterminated" triggers an unterminated-string LexError. It should
	// appear in result.Errors as a ParseError (lex errors are promoted).
	result := ParseBestEffort("'unterminated")
	if len(result.Errors) == 0 {
		t.Fatalf("expected at least one error, got none")
	}
	// Check that at least one error mentions "unterminated".
	found := false
	for _, e := range result.Errors {
		if strings.Contains(e.Msg, "unterminated") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a 'unterminated' error, got: %+v", result.Errors)
	}
}

func TestParse_BeginEndBlockOneError(t *testing.T) {
	// F3's Split treats "BEGIN SELECT 1; END;" as ONE segment, so parseSingle
	// sees the whole block as one statement starting with BEGIN. The BEGIN
	// stub emits ONE ParseError for the whole block.
	runParseBestEffortCase(t, parseCase{
		name:        "BEGIN..END block",
		input:       "BEGIN SELECT 1; END;",
		wantStmtCnt: 0,
		wantErrCnt:  1,
		wantErrMsgs: []string{"BEGIN statement parsing is not yet supported"},
		wantErrLocs: []int{0},
	})
}

func TestParse_StrictVsBestEffort(t *testing.T) {
	// Parse returns the first error; ParseBestEffort returns all errors.
	input := "SELECT 1; INSERT INTO t VALUES (1);"

	file, err := Parse(input)
	if err == nil {
		t.Errorf("Parse: expected error, got nil")
	} else {
		pe, ok := err.(*ParseError)
		if !ok {
			t.Errorf("Parse: expected *ParseError, got %T", err)
		} else if !strings.Contains(pe.Msg, "SELECT") {
			t.Errorf("Parse: first error Msg = %q, want to contain SELECT", pe.Msg)
		}
	}
	if file == nil {
		t.Errorf("Parse: file is nil, want non-nil")
	}

	result := ParseBestEffort(input)
	if len(result.Errors) != 2 {
		t.Errorf("ParseBestEffort: got %d errors, want 2", len(result.Errors))
	}
}

func TestParse_StrictNoErrors(t *testing.T) {
	// Empty input has no errors; Parse should return (file, nil).
	file, err := Parse("")
	if err != nil {
		t.Errorf("Parse(\"\") error = %v, want nil", err)
	}
	if file == nil {
		t.Errorf("Parse(\"\") file = nil, want non-nil")
	}
}

func TestParse_AbsoluteSegmentPositions(t *testing.T) {
	// Given "SELECT 1; SELECT 2;", the second SELECT error's Loc.Start
	// should be 10 (the absolute byte position), NOT 0 (which would be
	// relative to the second segment).
	result := ParseBestEffort("SELECT 1; SELECT 2;")
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %+v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Loc.Start != 0 {
		t.Errorf("first error Loc.Start = %d, want 0", result.Errors[0].Loc.Start)
	}
	if result.Errors[1].Loc.Start != 10 {
		t.Errorf("second error Loc.Start = %d, want 10", result.Errors[1].Loc.Start)
	}
}

func TestParse_LegacyCorpus(t *testing.T) {
	// Smoke test: run ParseBestEffort against every .sql file in the
	// legacy corpus. Assert no panic; log the ParseError count per file
	// (most files will have many "X not supported" errors — that's fine
	// for stubs-only F4).
	corpusDir := "testdata/legacy"
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("failed to read corpus dir %s: %v", corpusDir, err)
	}

	fileCount := 0
	totalErrors := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		fileCount++

		t.Run(entry.Name(), func(t *testing.T) {
			path := filepath.Join(corpusDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			// Must not panic.
			result := ParseBestEffort(string(data))
			totalErrors += len(result.Errors)
			t.Logf("%s: %d stmts, %d errors", entry.Name(), len(result.File.Stmts), len(result.Errors))
		})
	}

	if fileCount == 0 {
		t.Errorf("found 0 .sql files in %s — corpus missing?", corpusDir)
	}
	t.Logf("legacy corpus: %d files, %d total parse errors", fileCount, totalErrors)
}
