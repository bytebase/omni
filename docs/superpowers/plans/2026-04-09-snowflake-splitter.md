# Snowflake Statement Splitter (F3) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `Split(input string) []Segment` function in `snowflake/parser/split.go` — a lexer-based streaming state machine that extracts top-level SQL statement byte ranges from a source string, correctly handling Snowflake Scripting `BEGIN..END` blocks (including `DECLARE..BEGIN..END` and nested blocks).

**Architecture:** F3 adds two files (`split.go` and `split_test.go`) to the existing `snowflake/parser` package. `Split` iterates tokens via F2's `Lexer.NextToken()` (no byte-walking), tracks a 3-state state machine (`stateTop`, `stateInDeclare`, `stateInBlock`), and uses a 1-token lookahead buffer to disambiguate TCL `BEGIN` (followed by `TRANSACTION`/`WORK`/`NAME`/`;`/EOF) from Snowflake Scripting `BEGIN`.

**Tech Stack:** Go 1.25, stdlib only. Depends only on F2's existing `Lexer`, `Token`, and keyword constants (`kwBEGIN`, `kwEND`, `kwDECLARE`, `kwTRANSACTION`, `kwWORK`, `kwNAME`).

**Spec:** `docs/superpowers/specs/2026-04-09-snowflake-splitter-design.md` (commit `f1e7d55`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/splitter` on branch `feat/snowflake/splitter`

**Module path:** `github.com/bytebase/omni`. F3 lives in `github.com/bytebase/omni/snowflake/parser` (same package as F2).

**Commit policy:** No commits during implementation. The user reviews the full diff at the end; only after explicit approval do we commit.

---

## File Structure

| File | Responsibility | Approx LOC |
|------|----------------|-----------|
| `snowflake/parser/split.go` | `Segment` struct, `Split` function, state machine, `Empty` method, `isTCLBeginFollower` helper | 180 |
| `snowflake/parser/split_test.go` | Table-driven tests for simple splits, edge cases, BEGIN..END blocks, and the legacy corpus smoke test | 400 |

Total: ~580 LOC across 2 new files. Both live in the existing `snowflake/parser` package — no new subpackage.

---

## Task 1: Scaffold split.go and split_test.go

**Files:**
- Create: `snowflake/parser/split.go`
- Create: `snowflake/parser/split_test.go`

- [ ] **Step 1: Confirm worktree state**

Run: `pwd && git rev-parse --abbrev-ref HEAD`
Expected:
```
/Users/h3n4l/OpenSource/omni/.worktrees/splitter
feat/snowflake/splitter
```

If either is wrong, stop. Do NOT proceed.

- [ ] **Step 2: Write `snowflake/parser/split.go` stub**

```go
package parser
```

- [ ] **Step 3: Write `snowflake/parser/split_test.go` stub**

```go
package parser
```

- [ ] **Step 4: Verify the package still compiles**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 5: Verify existing tests still pass**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

All existing F2 tests (lexer, position, edge cases, error recovery, legacy corpus, keyword completeness) should still pass since we haven't changed any code.

---

## Task 2: Implement Segment struct and Empty method

**Files:**
- Modify: `snowflake/parser/split.go`

- [ ] **Step 1: Replace the stub with the Segment type and Empty method**

Overwrite `snowflake/parser/split.go` with:

```go
package parser

// Segment represents one top-level SQL statement extracted from a source string.
//
// Text is the raw substring of the input from ByteStart (inclusive) to
// ByteEnd (exclusive). It is NOT trimmed — leading/trailing whitespace and
// comments are preserved verbatim. The trailing `;` delimiter (if present)
// is NOT part of Text or ByteEnd; it lives between this segment's ByteEnd
// and the next segment's ByteStart.
type Segment struct {
	Text      string // the raw text of the statement (no trailing semicolon)
	ByteStart int    // inclusive start byte offset in the original source
	ByteEnd   int    // exclusive end byte offset; points AT the trailing ; if present, otherwise at len(input)
}

// Empty reports whether the segment contains no meaningful SQL content —
// that is, only whitespace and comments. It works by re-lexing Text and
// checking whether the first token is tokEOF.
//
// Empty segments are filtered out by Split.
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Type == tokEOF
}
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Verify vet**

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 4: Verify existing tests still pass**

Run: `go test ./snowflake/parser/...`
Expected: clean pass.

---

## Task 3: Implement Split function and state machine

**Files:**
- Modify: `snowflake/parser/split.go`

- [ ] **Step 1: Append the Split function and its helpers to split.go**

Use Edit to append the following code to `snowflake/parser/split.go` (after the `Empty` method, at end of file):

```go

// splitState is the state machine state for Split.
type splitState int

const (
	stateTop splitState = iota
	stateInDeclare
	stateInBlock
)

// Split extracts top-level SQL statements from input.
//
// The returned slice contains one Segment per non-empty statement. Empty
// statements (lone semicolons, comment-only chunks) are filtered out.
//
// Split is infallible — it always returns a result, even for malformed
// input. Lexing errors are suppressed internally; callers who need them
// should use NewLexer(input).Errors() directly.
//
// Split correctly handles:
//   - Single-quoted strings with '' and \ escapes
//   - Double-quoted identifiers with "" escape
//   - $$...$$ dollar strings
//   - X'...' hex literals
//   - Line comments (-- and //) and block comments (/* */ including nested)
//   - Snowflake Scripting BEGIN..END blocks (including nested and DECLARE..BEGIN..END)
//   - Inline procedure bodies (CREATE TASK/PROCEDURE/FUNCTION ... AS BEGIN ... END;)
//
// Split does NOT handle:
//   - IF/FOR/WHILE/REPEAT/CASE at top level (these are only valid inside a BEGIN..END body)
//   - DECLARE CURSOR at top level without a matching BEGIN (unusual; best-effort)
//   - DELIMITER directive (MySQL-specific, not used in Snowflake)
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0

	// We need one-token lookahead after kwBEGIN to disambiguate TCL from
	// scripting. Use a one-slot buffered lookahead.
	var pending *Token
	nextToken := func() Token {
		if pending != nil {
			t := *pending
			pending = nil
			return t
		}
		return l.NextToken()
	}
	peekToken := func() Token {
		if pending == nil {
			t := l.NextToken()
			pending = &t
		}
		return *pending
	}

	// emit creates a Segment covering input[stmtStart:end] — where end is
	// the byte offset BEFORE any trailing delimiter. The caller is
	// responsible for advancing stmtStart past the delimiter.
	emit := func(end int) {
		seg := Segment{
			Text:      input[stmtStart:end],
			ByteStart: stmtStart,
			ByteEnd:   end,
		}
		if !seg.Empty() {
			segments = append(segments, seg)
		}
	}

	for {
		tok := nextToken()
		if tok.Type == tokEOF {
			break
		}

		switch state {
		case stateTop:
			switch tok.Type {
			case kwDECLARE:
				state = stateInDeclare
			case kwBEGIN:
				next := peekToken()
				if isTCLBeginFollower(next) {
					// Stay in stateTop; the buffered peek will be returned
					// by the next nextToken() call and handled normally
					// (including the `;` emission path if applicable).
				} else {
					state = stateInBlock
					depth = 1
				}
			case int(';'):
				// tok.Loc.Start is the position of the `;`, tok.Loc.End is
				// one past the `;`. Segment text stops BEFORE the `;`; the
				// next segment starts AFTER the `;`.
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
			}

		case stateInDeclare:
			switch tok.Type {
			case kwBEGIN:
				state = stateInBlock
				depth = 1
			}
			// Semicolons and all other tokens are absorbed in stateInDeclare.

		case stateInBlock:
			switch tok.Type {
			case kwBEGIN:
				depth++
			case kwEND:
				if depth > 0 {
					depth--
					if depth == 0 {
						state = stateTop
					}
				}
			}
			// Semicolons and all other tokens are absorbed in stateInBlock.
		}
	}

	// Trailing segment — whatever remains after the last `;` (or the whole
	// input if there was no `;`). ByteEnd is len(input) in this case.
	if stmtStart < len(input) {
		emit(len(input))
	}

	return segments
}

// isTCLBeginFollower reports whether tok, appearing immediately after a
// top-level BEGIN keyword, indicates a transaction-start (TCL) rather than
// a Snowflake Scripting block opener.
//
// TCL forms: BEGIN;, BEGIN TRANSACTION, BEGIN WORK, BEGIN NAME <id>, BEGIN EOF
func isTCLBeginFollower(tok Token) bool {
	switch tok.Type {
	case int(';'), tokEOF, kwTRANSACTION, kwWORK, kwNAME:
		return true
	}
	return false
}
```

- [ ] **Step 2: Verify compile**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Verify vet**

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 4: Verify existing tests still pass**

Run: `go test ./snowflake/parser/...`
Expected: clean pass. No new tests yet, but the F2 tests should still be green.

- [ ] **Step 5: Smoke-test Split by hand**

Run a quick sanity check via a temporary Go file:

```bash
cat > /tmp/split-smoke.go <<'EOF'
package main

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/parser"
)

func main() {
	cases := []string{
		"SELECT 1;",
		"SELECT 1; SELECT 2;",
		"BEGIN SELECT 1; END;",
		"BEGIN; SELECT 1;",
	}
	for _, c := range cases {
		segs := parser.Split(c)
		fmt.Printf("%q → %d segments\n", c, len(segs))
		for _, s := range segs {
			fmt.Printf("  [%d,%d] %q\n", s.ByteStart, s.ByteEnd, s.Text)
		}
	}
}
EOF
go run /tmp/split-smoke.go
rm /tmp/split-smoke.go
```

Expected output (exact):
```
"SELECT 1;" → 1 segments
  [0,8] "SELECT 1"
"SELECT 1; SELECT 2;" → 2 segments
  [0,8] "SELECT 1"
  [9,18] " SELECT 2"
"BEGIN SELECT 1; END;" → 1 segments
  [0,19] "BEGIN SELECT 1; END"
"BEGIN; SELECT 1;" → 2 segments
  [0,5] "BEGIN"
  [6,15] " SELECT 1"
```

If the output doesn't match exactly, stop and diagnose. The most common issue will be byte-position math (off-by-one).

---

## Task 4: Write basic table-driven tests

**Files:**
- Modify: `snowflake/parser/split_test.go`

- [ ] **Step 1: Replace the stub with the test helper and basic cases**

Overwrite `snowflake/parser/split_test.go` with:

```go
package parser

import (
	"reflect"
	"testing"
)

// runSplitCases is a table-driven helper that asserts Split(input) returns
// exactly the expected Segment list.
func runSplitCases(t *testing.T, cases []struct {
	name  string
	input string
	want  []Segment
}) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Split(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Split(%q) mismatch\n got: %+v\nwant: %+v", c.input, got, c.want)
			}
		})
	}
}

func TestSplit_Empty(t *testing.T) {
	if got := Split(""); got != nil {
		t.Errorf("Split(\"\") = %+v, want nil", got)
	}
}

func TestSplit_Simple(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"one with semi", "SELECT 1;", []Segment{{"SELECT 1", 0, 8}}},
		{"one no semi", "SELECT 1", []Segment{{"SELECT 1", 0, 8}}},
		{"two stmts", "SELECT 1;SELECT 2;", []Segment{
			{"SELECT 1", 0, 8},
			{"SELECT 2", 9, 17},
		}},
		{"three stmts", "A;B;C;", []Segment{
			{"A", 0, 1},
			{"B", 2, 3},
			{"C", 4, 5},
		}},
		{"trailing whitespace", "SELECT 1;  ", []Segment{{"SELECT 1", 0, 8}}},
		{"leading whitespace", "  SELECT 1;", []Segment{{"  SELECT 1", 0, 10}}},
	}
	runSplitCases(t, cases)
}

func TestSplit_EmptyStatements(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"lone semi", ";", nil},
		{"triple semi", ";;;", nil},
		{"leading semi", ";SELECT 1;", []Segment{{"SELECT 1", 1, 9}}},
		{"trailing semi", "SELECT 1;;", []Segment{{"SELECT 1", 0, 8}}},
		{"embedded semi", "SELECT 1;;SELECT 2;", []Segment{
			{"SELECT 1", 0, 8},
			{"SELECT 2", 10, 18},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_Positions(t *testing.T) {
	// Verify byte offsets are accurate for a representative multi-statement input.
	input := "SELECT 1; SELECT 2; SELECT 3"
	//        0123456789...
	//        s1=[0,8] ; s2=[9,18] ; s3=[19,28]
	got := Split(input)
	want := []Segment{
		{"SELECT 1", 0, 8},
		{" SELECT 2", 9, 18},
		{" SELECT 3", 19, 28},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split mismatch\n got: %+v\nwant: %+v", got, want)
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `go test ./snowflake/parser/...`
Expected: clean pass. All F2 tests plus 4 new F3 tests (`TestSplit_Empty`, `TestSplit_Simple`, `TestSplit_EmptyStatements`, `TestSplit_Positions`).

- [ ] **Step 3: Run with verbose output to verify each subtest**

Run: `go test -v ./snowflake/parser/... -run TestSplit 2>&1 | tail -30`
Expected: each subtest reported as PASS:
```
=== RUN   TestSplit_Empty
--- PASS: TestSplit_Empty
=== RUN   TestSplit_Simple
=== RUN   TestSplit_Simple/one_with_semi
--- PASS: TestSplit_Simple/one_with_semi
...
```

All subtests pass. If any fail, diagnose the exact byte-offset mismatch and fix either the test or the production code.

---

## Task 5: Write edge case tests

**Files:**
- Modify: `snowflake/parser/split_test.go`

- [ ] **Step 1: Append edge case tests to split_test.go**

Use Edit to append the following code to the end of `snowflake/parser/split_test.go`:

```go

func TestSplit_SemisInStrings(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"single-quoted", "SELECT 'a;b';", []Segment{{"SELECT 'a;b'", 0, 12}}},
		{"doubled-quote in string", "SELECT 'a''b;c';", []Segment{{"SELECT 'a''b;c'", 0, 15}}},
		{"dollar-quoted", "SELECT $$a;b$$;", []Segment{{"SELECT $$a;b$$", 0, 14}}},
		{"hex literal", "SELECT X'3B';", []Segment{{"SELECT X'3B'", 0, 12}}},
		{"quoted ident with semi", `SELECT "a;b";`, []Segment{{`SELECT "a;b"`, 0, 12}}},
	}
	runSplitCases(t, cases)
}

func TestSplit_Comments(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"dash line before", "-- comment\nSELECT 1;", []Segment{
			{"-- comment\nSELECT 1", 0, 19},
		}},
		{"slash line before", "// comment\nSELECT 1;", []Segment{
			{"// comment\nSELECT 1", 0, 19},
		}},
		{"block before", "/* comment */ SELECT 1;", []Segment{
			{"/* comment */ SELECT 1", 0, 22},
		}},
		{"inline block", "SELECT /* a; b */ 1;", []Segment{
			{"SELECT /* a; b */ 1", 0, 19},
		}},
		{"nested block", "/* outer /* inner */ still */ SELECT 1;", []Segment{
			{"/* outer /* inner */ still */ SELECT 1", 0, 38},
		}},
		{"comment-only filtered", "-- just a comment\n", nil},
		{"block comment only filtered", "/* only */", nil},
	}
	runSplitCases(t, cases)
}

func TestSplit_TCLBegin(t *testing.T) {
	// BEGIN followed by ;, TRANSACTION, WORK, NAME, or EOF is a TCL
	// statement, not a Snowflake Scripting block.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"bare begin semi", "BEGIN; SELECT 1; COMMIT;", []Segment{
			{"BEGIN", 0, 5},
			{" SELECT 1", 6, 15},
			{" COMMIT", 16, 23},
		}},
		{"begin transaction", "BEGIN TRANSACTION; SELECT 1; COMMIT;", []Segment{
			{"BEGIN TRANSACTION", 0, 17},
			{" SELECT 1", 18, 27},
			{" COMMIT", 28, 35},
		}},
		{"begin work", "BEGIN WORK; COMMIT;", []Segment{
			{"BEGIN WORK", 0, 10},
			{" COMMIT", 11, 18},
		}},
		{"begin name", "BEGIN NAME tx1; COMMIT;", []Segment{
			{"BEGIN NAME tx1", 0, 14},
			{" COMMIT", 15, 22},
		}},
		{"begin at eof", "BEGIN", []Segment{
			{"BEGIN", 0, 5},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_ScriptingBegin(t *testing.T) {
	// BEGIN followed by non-TCL tokens (INSERT, SELECT, DECLARE, etc.) is
	// a Snowflake Scripting block — the whole BEGIN..END is one statement.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"simple block", "BEGIN SELECT 1; END;", []Segment{
			{"BEGIN SELECT 1; END", 0, 19},
		}},
		{"two stmts in block", "BEGIN INSERT INTO t VALUES (1); INSERT INTO t VALUES (2); END;", []Segment{
			{"BEGIN INSERT INTO t VALUES (1); INSERT INTO t VALUES (2); END", 0, 61},
		}},
		{"nested block", "BEGIN BEGIN SELECT 1; END; END;", []Segment{
			{"BEGIN BEGIN SELECT 1; END; END", 0, 30},
		}},
		{"block then top-level", "BEGIN SELECT 1; END; SELECT 2;", []Segment{
			{"BEGIN SELECT 1; END", 0, 19},
			{" SELECT 2", 20, 29},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_DeclareBegin(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"standalone declare begin", "DECLARE x INTEGER; BEGIN x := 1; END;", []Segment{
			{"DECLARE x INTEGER; BEGIN x := 1; END", 0, 36},
		}},
		{"multiple decls", "DECLARE x INT; y FLOAT; BEGIN x := 1; END;", []Segment{
			{"DECLARE x INT; y FLOAT; BEGIN x := 1; END", 0, 41},
		}},
		{"declare then top-level", "DECLARE x INT; BEGIN SELECT 1; END; SELECT 2;", []Segment{
			{"DECLARE x INT; BEGIN SELECT 1; END", 0, 34},
			{" SELECT 2", 35, 44},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_InlineProcedureBody(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"create task begin end", "CREATE TASK t AS BEGIN SELECT 1; END;", []Segment{
			{"CREATE TASK t AS BEGIN SELECT 1; END", 0, 36},
		}},
		{"create task declare begin end", "CREATE TASK t AS DECLARE x FLOAT; BEGIN x := 1; END;", []Segment{
			{"CREATE TASK t AS DECLARE x FLOAT; BEGIN x := 1; END", 0, 51},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_DollarQuotedBody(t *testing.T) {
	// Dollar-quoted strings are opaque to the splitter — the entire $$..$$
	// block lexes as one tokString, so internal ; BEGIN END are invisible.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"proc with dollar body", "CREATE PROCEDURE p() AS $$ BEGIN SELECT 1; END; $$;", []Segment{
			{"CREATE PROCEDURE p() AS $$ BEGIN SELECT 1; END; $$", 0, 50},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_UnclosedBlock(t *testing.T) {
	// BEGIN without a matching END should produce a best-effort single
	// segment covering the whole input.
	got := Split("BEGIN SELECT 1;")
	want := []Segment{{"BEGIN SELECT 1;", 0, 15}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split(unclosed BEGIN) = %+v, want %+v", got, want)
	}
}

func TestSplit_TrailingWhitespaceAfterSemi(t *testing.T) {
	// Whitespace-only content after the last `;` should NOT produce a
	// trailing segment (filtered by Empty()).
	got := Split("SELECT 1;   \n\t")
	want := []Segment{{"SELECT 1", 0, 8}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split mismatch\n got: %+v\nwant: %+v", got, want)
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `go test ./snowflake/parser/...`
Expected: clean pass.

- [ ] **Step 3: Verbose run to spot any failures**

Run: `go test -v ./snowflake/parser/... -run TestSplit 2>&1 | grep -E "^(=== RUN|--- (PASS|FAIL))" | tail -60`
Expected: every subtest PASS. No FAILs.

If any tests fail:
- If the failure is byte-offset arithmetic (e.g. expected ByteEnd=36 got 37), the test expectation is probably wrong — recompute by hand and fix the test
- If the failure is a missing segment or extra segment, the state machine logic is wrong — re-read `split.go` Task 3 step 1 and fix
- If the failure is in the Text field, the slicing bounds are wrong — check `input[stmtStart:end]`

---

## Task 6: Write legacy corpus smoke test

**Files:**
- Modify: `snowflake/parser/split_test.go`

- [ ] **Step 1: Append the legacy corpus test**

Use Edit to append the following code to the end of `snowflake/parser/split_test.go`:

```go

func TestSplit_LegacyCorpus(t *testing.T) {
	corpusDir := "testdata/legacy"
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("failed to read corpus dir %s: %v", corpusDir, err)
	}

	sqlCount := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		sqlCount++

		t.Run(entry.Name(), func(t *testing.T) {
			path := filepath.Join(corpusDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			segments := Split(string(data))

			// Every corpus file must produce at least one segment.
			if len(segments) == 0 {
				t.Errorf("%s: Split returned 0 segments", entry.Name())
				return
			}

			// Every segment must have valid byte range and non-empty text.
			for i, seg := range segments {
				if seg.ByteStart >= seg.ByteEnd {
					t.Errorf("%s: segment[%d] ByteStart=%d >= ByteEnd=%d",
						entry.Name(), i, seg.ByteStart, seg.ByteEnd)
				}
				if seg.Empty() {
					t.Errorf("%s: segment[%d] Text is empty or whitespace-only: %q",
						entry.Name(), i, seg.Text)
				}
				// Text must equal the raw slice.
				if seg.Text != string(data)[seg.ByteStart:seg.ByteEnd] {
					t.Errorf("%s: segment[%d] Text != input[%d:%d]",
						entry.Name(), i, seg.ByteStart, seg.ByteEnd)
				}
			}
		})
	}

	if sqlCount == 0 {
		t.Errorf("found 0 .sql files in %s — corpus missing?", corpusDir)
	}
	if sqlCount != 27 {
		t.Logf("note: found %d .sql files (expected 27 from the legacy lift)", sqlCount)
	}
}

func TestSplit_TaskScriptingHasTwoSegments(t *testing.T) {
	// task_scripting.sql contains two CREATE TASK statements, each with
	// an inline BEGIN..END body. The legacy (dumb) splitter would produce
	// 5+ fragments. F3's state machine must produce exactly 2 segments.
	data, err := os.ReadFile("testdata/legacy/task_scripting.sql")
	if err != nil {
		t.Fatalf("read task_scripting.sql: %v", err)
	}
	segments := Split(string(data))
	if len(segments) != 2 {
		t.Errorf("Split(task_scripting.sql) returned %d segments, want 2:", len(segments))
		for i, seg := range segments {
			t.Errorf("  [%d] [%d,%d] %q", i, seg.ByteStart, seg.ByteEnd, seg.Text)
		}
	}
}
```

- [ ] **Step 2: Add the required imports**

The new tests use `os`, `path/filepath`, and `strings`. Use Edit to update the import block at the top of `split_test.go`:

```go
package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)
```

(Replace the existing import block `import (\n\t"reflect"\n\t"testing"\n)` with the expanded version above.)

- [ ] **Step 3: Run the tests**

Run: `go test ./snowflake/parser/...`
Expected: clean pass. The legacy corpus test adds 27 sub-tests (one per `.sql` file) plus the `task_scripting` assertion.

- [ ] **Step 4: Verbose verification**

Run: `go test -v ./snowflake/parser/... -run "TestSplit_LegacyCorpus|TestSplit_TaskScriptingHasTwoSegments" 2>&1 | tail -40`
Expected: every `.sql` sub-test PASS, and `TestSplit_TaskScriptingHasTwoSegments` PASS.

If `TestSplit_TaskScriptingHasTwoSegments` fails, the state machine is not handling `CREATE TASK ... AS BEGIN ... END;` correctly — likely the `kwBEGIN` lookahead isn't transitioning to `stateInBlock`, or the `kwDECLARE` case is mishandling the `DECLARE x FLOAT; BEGIN` sequence in the second task.

If any corpus file fails, the splitter is producing either too many or zero segments. Read the failing file and trace through the state machine by hand to find the bug.

---

## Task 7: Final acceptance criteria sweep

This task runs every acceptance criterion from the spec back-to-back to confirm F3 is complete.

- [ ] **Step 1: Build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 2: Vet**

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Gofmt**

Run: `gofmt -l snowflake/parser/`
Expected: no output (no files need reformatting).

If any file is listed, run `gofmt -w snowflake/parser/` to apply and re-run the check.

- [ ] **Step 4: Test**

Run: `go test ./snowflake/parser/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

- [ ] **Step 5: Verbose test run to confirm test counts**

Run: `go test -v ./snowflake/parser/... 2>&1 | tail -60`
Expected: all tests PASS. F2 tests (lexer, position, edge, recovery, legacy corpus, keyword completeness) still green. F3 tests (TestSplit_*) all green.

- [ ] **Step 6: Package isolation**

Run: `go build ./snowflake/...`
Expected: no output, exit 0. Confirms `snowflake/parser` builds without depending on any package beyond F1's `snowflake/ast`.

- [ ] **Step 7: List all files for review**

Run: `git status snowflake/parser docs/superpowers`
Expected output:
```
On branch feat/snowflake/splitter
Changes not staged for commit:  (none)
Untracked files:
	snowflake/parser/split.go
	snowflake/parser/split_test.go
```

Plus the already-committed spec doc and plan doc at `docs/superpowers/specs/2026-04-09-snowflake-splitter-design.md` and `docs/superpowers/plans/2026-04-09-snowflake-splitter.md` (which will be tracked).

- [ ] **Step 8: STOP and present for review**

Do NOT commit. Output a summary:

> F3 (snowflake/parser splitter) implementation complete.
>
> All 7 tasks done. Acceptance criteria green:
> - `go build ./snowflake/parser/...` ✓
> - `go vet ./snowflake/parser/...` ✓
> - `gofmt -l snowflake/parser/` clean ✓
> - `go test ./snowflake/parser/...` ✓
> - 27 legacy corpus files produce ≥1 non-empty segments ✓
> - task_scripting.sql returns exactly 2 segments ✓
> - `go build ./snowflake/...` (isolation) ✓
>
> Files created:
> - snowflake/parser/split.go
> - snowflake/parser/split_test.go
> - docs/superpowers/plans/2026-04-09-snowflake-splitter.md (this plan)
>
> Please review the diff. After your approval I will commit and run finishing-a-development-branch.

---

## Spec Coverage Checklist

| Spec section | Covered by |
|--------------|-----------|
| Purpose | Implicitly by having the Split function |
| Architecture (lexer-based, 3-state machine) | Task 3 |
| Public API (Segment, Split, Empty) | Tasks 2 + 3 |
| State machine transitions (stateTop/stateInDeclare/stateInBlock) | Task 3 |
| TCL vs scripting BEGIN disambiguation | Task 3 (isTCLBeginFollower) + Task 5 tests |
| Final segment emission (no trailing `;`) | Task 3 (`if stmtStart < len(input)` branch) |
| Empty segment filtering | Task 2 (Empty method) + Task 3 (`if !seg.Empty()`) |
| Key Snowflake behaviors (strings, comments, blocks) | Task 5 tests |
| Divergence from legacy (task_scripting.sql → 2 segments) | Task 6 explicit assertion |
| File layout (split.go + split_test.go) | Task 1 scaffold |
| 14 test categories | Tasks 4-6 |
| 27-file legacy corpus smoke test | Task 6 |
| Acceptance criteria 1-8 | Task 7 |
