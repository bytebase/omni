# Snowflake Statement Splitter (F3) — Design

**DAG node:** F3 — `snowflake/parser` statement splitter
**Migration:** `docs/migration/snowflake/dag.md`
**Branch:** `feat/snowflake/splitter`
**Status:** Approved, ready for plan + implementation.
**Depends on:** F2 (lexer, merged via PR #17) — uses `Lexer.NextToken` as its input source.
**Unblocks:** nothing parser-side (F4 depends on F1+F2, not F3). But F3 is on the critical path for bytebase's `split.go` and `statement_ranges.go` consumers, which can migrate to omni as soon as F3 lands.

## Purpose

F3 extracts top-level SQL statement byte ranges from a source string. Given a SQL input, it produces a list of `Segment{Text, ByteStart, ByteEnd}` values where each segment is one top-level statement. The byte range covers the statement text only — the trailing `;` (if present) is NOT part of the segment's `Text` or `ByteEnd`. This matches `mysql/parser/split.go`'s convention.

Bytebase consumers today:
- `backend/plugin/parser/snowflake/split.go` — returns `[]base.Statement` for the editor's statement-execution feature
- `backend/plugin/parser/snowflake/statement_ranges.go` — returns UTF-16 byte ranges for the editor's statement-highlighting feature

Both consumers currently rely on the legacy ANTLR4 parser's token stream via `base.SplitSQLByLexer(stream, SnowflakeLexerSEMI)`. That base helper splits on every `SEMI` token with **zero depth tracking** — meaning it produces a broken result for Snowflake Scripting `BEGIN..END` blocks. F3 fixes this bug.

## Architecture

F3 is a lexer-based streaming state machine. It reuses F2's `Lexer.NextToken()` rather than doing byte-level scanning (the approach `mysql/parser/split.go` uses). This is a deliberate divergence from the mysql reference — F2's lexer already handles every Snowflake string/comment form correctly (`'...'`, `"..."`, `$$...$$`, `X'...'`, nested `/* */`, `--`, `//`), so re-implementing those byte-walking details would duplicate ~400 LOC of logic. F3 is ~200 LOC instead.

The state machine tracks three states to handle Snowflake Scripting blocks correctly:

- **`stateTop`** — top-level. Semicolons emit statement boundaries.
- **`stateInDeclare`** — inside a `DECLARE` region before the matching `BEGIN` arrives. Semicolons are suppressed (they terminate individual declarations, not the outer statement).
- **`stateInBlock(depth)`** — inside a `BEGIN..END` body. Semicolons at depth > 0 are suppressed. Nested `BEGIN..END` pairs are supported via the depth counter.

TCL `BEGIN` (transaction start) is disambiguated from scripting `BEGIN` by one-token lookahead: if the token after `BEGIN` is `kwTRANSACTION`, `kwWORK`, `kwNAME`, `;`, or EOF, it's a TCL transaction start and no block is entered.

## Public API

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
func (s Segment) Empty() bool

// Split extracts top-level SQL statements from input. The returned slice
// contains one Segment per non-empty statement. Empty statements (lone
// semicolons, comment-only chunks) are filtered out.
//
// Split is infallible — it always returns a result, even for malformed
// input. Lexing errors are suppressed internally; callers who need them
// should use NewLexer(input).Errors() directly.
//
// Split correctly handles:
//   - Single-quoted strings with ' and \\ escapes
//   - Double-quoted identifiers with " escape
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
func Split(input string) []Segment
```

No second `SplitRanges` function. Bytebase consumers that need only byte ranges can call `Split` and extract `{ByteStart, ByteEnd}` from each segment. The convenience wrapper is not worth the API surface.

No error return. The splitter is best-effort; lex errors are available via `Lexer.Errors()` if the caller constructs a lexer explicitly, but `Split` itself is infallible.

## State Machine

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│                        ┌──────────┐                             │
│                        │          │                             │
│                  ┌────►│ stateTop │◄─────────────────┐          │
│                  │     │          │                  │          │
│                  │     └──┬───┬───┘                  │          │
│                  │        │   │                      │          │
│                  │  kwBEGIN   kwDECLARE              │          │
│                  │  (not TCL) │                      │          │
│                  │        │   │                      │          │
│                  │        ▼   ▼                      │          │
│                  │  ┌─────────────┐   kwBEGIN   ┌────┴────────┐ │
│               emit ;│             │─────────────►             │ │
│             (depth  │stateInDeclare│             │stateInBlock │ │
│                =0)  │             │             │  (depth=1)  │ │
│                  │  └─────────────┘             │             │ │
│                  │                              └─┬─────────┬─┘ │
│                  │                                │         │   │
│                  │                           kwBEGIN       kwEND│
│                  │                           depth++       depth│
│                  │                                │         -- │
│                  │                                └─────────┘   │
│                  │                                  (stays)     │
│                  │                                              │
│                  │                      kwEND when depth==0     │
│                  └──────────────────────────────────────────────┘
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Transition rules

**`stateTop`** → initial state. `stmtStart` points at the first byte of the current statement.
- Token `kwDECLARE` → transition to `stateInDeclare`
- Token `kwBEGIN` → one-token lookahead: if the next token is `kwTRANSACTION`, `kwWORK`, `kwNAME`, `int(';')`, or `tokEOF`, stay in `stateTop` (TCL BEGIN); otherwise transition to `stateInBlock` with `depth = 1`
- Token `int(';')` → emit `Segment{Text: trim(input[stmtStart:tok.Loc.End]), ByteStart: stmtStart, ByteEnd: tok.Loc.End}`; set `stmtStart = tok.Loc.End`
- Any other token → stay

**`stateInDeclare`** → inside a DECLARE region, waiting for the matching BEGIN. The semicolons inside DECLARE terminate individual declarations and must not split the outer statement.
- Token `kwBEGIN` → transition to `stateInBlock` with `depth = 1`
- Token `int(';')` → suppressed (no emission)
- Any other token → stay
- Token `tokEOF` → emit remaining input as a best-effort final segment

**`stateInBlock(depth)`** → inside a `BEGIN..END` body. Semicolons are suppressed.
- Token `kwBEGIN` → `depth++` (nested block)
- Token `kwEND` → `depth--`; if `depth == 0`, transition to `stateTop`
- Token `int(';')` → suppressed
- Any other token → stay
- Token `tokEOF` → emit remaining input as a best-effort final segment (unclosed block)

### Final segment emission

After the main loop exits (on `tokEOF`), if `stmtStart < len(input)`, emit a trailing segment from `stmtStart` to `len(input)`. This covers the case where the last statement has no trailing semicolon.

Apply `Segment.Empty()` filter to all emitted segments — skip any segment whose `Text` (after trimming) is whitespace-only.

## Key Snowflake-specific Behaviors

| Input | Expected output | Why |
|-------|-----------------|-----|
| `SELECT 1;` | 1 segment | Trivial split on `;` |
| `SELECT 1; SELECT 2;` | 2 segments | Two top-level statements |
| `SELECT 1` | 1 segment (no trailing `;`) | Trailing statement without `;` |
| `;;;` | 0 segments | All empty, filtered out |
| `SELECT 'a;b';` | 1 segment | `'a;b'` lexes as one `tokString`, semicolon invisible |
| `SELECT $$a;b$$;` | 1 segment | `$$a;b$$` lexes as one `tokString` |
| `SELECT /* ; */ 1;` | 1 segment | Comment dropped by `skipWhitespaceAndComments` |
| `BEGIN; SELECT 1; COMMIT;` | 3 segments | TCL BEGIN (next is `;`) → stateTop → splits normally |
| `BEGIN TRANSACTION; SELECT 1; COMMIT;` | 3 segments | TCL BEGIN (next is `kwTRANSACTION`) → stateTop → splits normally |
| `BEGIN INSERT t VALUES (1); INSERT t VALUES (2); END;` | 1 segment | Scripting BEGIN → stateInBlock → suppresses inner `;` → emits on outer `;` after END |
| `BEGIN BEGIN SELECT 1; END; END;` | 1 segment | Nested blocks handled via depth counter |
| `DECLARE x INTEGER; BEGIN x := 1; END;` | 1 segment | DECLARE → stateInDeclare → suppresses `;` after `INTEGER` → BEGIN → stateInBlock → END → stateTop → emits on final `;` |
| `CREATE PROCEDURE p() AS $$ BEGIN SELECT 1; END; $$;` | 1 segment | Body is one `tokString` (dollar-quoted); splitter never sees inner BEGIN/END/`;` |
| `CREATE TASK t AS BEGIN SELECT 1; END;` | 1 segment | Inline BEGIN..END body; state machine enters stateInBlock on `kwBEGIN`, suppresses inner `;`, emits on outer `;` after END |
| `CREATE TASK t AS DECLARE x FLOAT; BEGIN x := 1; END;` | 1 segment | Inline DECLARE..BEGIN..END body; stateInDeclare → stateInBlock → stateTop |
| `BEGIN SELECT 1;` (unclosed) | 1 segment (best-effort) | stateInBlock stays open until EOF, then emits remaining input |

The `task_scripting.sql` file in the legacy corpus contains the two inline-procedure forms (`CREATE TASK ... AS BEGIN ... END;` and `CREATE TASK ... AS DECLARE ... BEGIN ... END;`) and is the primary regression test for block handling. Expected: **2 segments** (one per CREATE TASK). Under the legacy broken splitter: 5+ fragments.

## Divergence from Legacy Splitter

The legacy bytebase splitter (`base.SplitSQLByLexer(stream, SnowflakeLexerSEMI)`) splits on every `SEMI` token with no depth tracking. This produces incorrect output for any input containing a `BEGIN..END` block — `task_scripting.sql` in particular fragments into 5+ pieces instead of 2.

F3 diverges from this behavior by tracking block depth correctly. This is a **fix**, not a compatibility break:
- The legacy output is broken by any reasonable interpretation (the fragments aren't individually valid SQL).
- Bytebase's editor-side consumers will get better results immediately after migration.
- No downstream behavior depends on the legacy splitter's broken fragmentation.

The divergence is documented in the bytebase migration notes so consumers can validate the new output against their expectations before cutover.

## File Layout

```
snowflake/parser/split.go        # Split function, Segment struct, state machine
snowflake/parser/split_test.go   # Table-driven tests + legacy corpus smoke test
```

Both files added to the existing `snowflake/parser` package (the same package F2's lexer lives in). No new subpackage. No new dependencies beyond `github.com/bytebase/omni/snowflake/ast` (which `snowflake/parser` already imports via F2) — and in fact F3 does not use `ast` directly at all; `Segment` uses plain `int` fields for byte offsets, matching mysql's `Segment` shape.

## Implementation Sketch

```go
package parser

// Segment represents one top-level SQL statement.
type Segment struct {
    Text      string
    ByteStart int
    ByteEnd   int
}

// Empty reports whether the segment contains no meaningful SQL content.
// It re-lexes Text and checks whether the first token is tokEOF.
func (s Segment) Empty() bool {
    return NewLexer(s.Text).NextToken().Type == tokEOF
}

// splitState is the state machine state for Split.
type splitState int

const (
    stateTop splitState = iota
    stateInDeclare
    stateInBlock
)

// Split extracts top-level SQL statements from input.
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
    // scripting. Use a buffered "pending" slot — one-slot ring buffer.
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
            // Semicolons and all other tokens are absorbed — stay in state.
            }

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
            // Semicolons and all other tokens are absorbed.
            }
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

Key positioning rules (matching `mysql/parser/split.go`):

- `tok.Loc.Start` for a `;` token is the byte offset OF the `;`.
- `tok.Loc.End` for a `;` token is the byte offset just past the `;`.
- A statement's `ByteEnd` points AT its trailing `;` (so `Text` excludes the `;`).
- The next statement's `ByteStart` is `tok.Loc.End` — one past the `;`, which may include leading whitespace before the next keyword (that's fine; `Text` is raw).

The `pending` slot acts as a 1-token buffer so we can peek ahead after `kwBEGIN` without losing that next token. When the TCL BEGIN path stays in `stateTop`, the buffered next token is returned by the subsequent `nextToken()` call and processed normally — for example, a `;` after `BEGIN;` will emit the segment correctly.

## Testing

Test file: `snowflake/parser/split_test.go`. Run with `go test ./snowflake/parser/...` (scoped — never global).

### Required test categories

1. **Simple splits** — `SELECT 1`, `SELECT 1;`, `SELECT 1; SELECT 2;`, `SELECT 1; SELECT 2; SELECT 3;`
2. **Empty statements** — `;`, `;;;`, `;SELECT 1;`, `SELECT 1;;`, `SELECT 1;;SELECT 2;`
3. **Semicolons in strings** — `SELECT 'a;b';`, `SELECT "a;b";`, `SELECT $$a;b$$;`, `SELECT X'3B';` (`X'3B'` is the byte `;`)
4. **Comments** — `-- comment\nSELECT 1;`, `// comment\nSELECT 1;`, `/* comment */ SELECT 1;`, `SELECT 1 /* a; b */;`, `/* nested /* */ */SELECT 1;`
5. **TCL BEGIN** — `BEGIN; SELECT 1; COMMIT;` (3 segments), `BEGIN TRANSACTION; SELECT 1; COMMIT;` (3 segments), `BEGIN WORK; COMMIT;` (2 segments), `BEGIN NAME tx1; COMMIT;` (2 segments)
6. **Scripting BEGIN..END** — `BEGIN SELECT 1; END;` (1 segment), `BEGIN INSERT t VALUES (1); INSERT t VALUES (2); END;` (1 segment)
7. **Nested BEGIN..END** — `BEGIN BEGIN SELECT 1; END; END;` (1 segment), 3-deep nesting
8. **DECLARE..BEGIN..END** — `DECLARE x INT; BEGIN x := 1; END;` (1 segment)
9. **Inline procedure bodies** — `CREATE TASK t AS BEGIN SELECT 1; END;` (1 segment), `CREATE TASK t AS DECLARE x FLOAT; BEGIN x := 1; END;` (1 segment)
10. **Dollar-quoted bodies** — `CREATE PROCEDURE p() AS $$ BEGIN SELECT 1; END; $$;` (1 segment; internal `;` invisible because the body is one `tokString`)
11. **Unclosed block** — `BEGIN SELECT 1;` (no END) → 1 best-effort segment covering the whole input
12. **Empty input** — `Split("")` returns `nil`
13. **Position accuracy** — for each representative case, assert `ByteStart` and `ByteEnd` are byte-correct against a hand-computed position
14. **Legacy corpus smoke test** — iterate every `.sql` file in `snowflake/parser/testdata/legacy/` (27 files); for each file assert:
    - `Split` returns at least 1 segment
    - Every segment has `ByteStart < ByteEnd`
    - No segment's `Text` is empty
    - The last segment's `ByteEnd` is at or beyond the file's last non-whitespace byte position
    - Specifically assert that `task_scripting.sql` returns exactly **2 segments** (the two CREATE TASK blocks)

### Test fixture style

Table-driven with a helper:

```go
func TestSplit_Simple(t *testing.T) {
    cases := []struct {
        name  string
        input string
        want  []Segment  // Text + ByteStart + ByteEnd
    }{
        {"empty", "", nil},
        {"one with semi", "SELECT 1;", []Segment{{"SELECT 1", 0, 8}}},      // ByteEnd=8 is the `;` position
        {"one no semi",   "SELECT 1",  []Segment{{"SELECT 1", 0, 8}}},      // ByteEnd=8 is len(input)
        {"two stmts",     "SELECT 1;SELECT 2;",
            []Segment{{"SELECT 1", 0, 8}, {"SELECT 2", 9, 17}}},            // Second starts at byte 9 (after `;`)
        // ...
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := Split(c.input)
            // assert got equals c.want
        })
    }
}
```

Note on `Text` semantics:
- `Text = input[ByteStart:ByteEnd]` — the raw substring, NOT trimmed
- `ByteEnd` is the position OF the trailing `;` (so Text excludes the `;`), or `len(input)` if there's no trailing `;`
- The next segment's `ByteStart` is the byte AFTER the `;` — meaning leading whitespace may appear at the start of subsequent segments. That's fine; `Empty()` will filter out whitespace-only chunks via re-lexing.

## Out of Scope

F3 does NOT include:

| Feature | Where it lives |
|---|---|
| `Parse(sql) (*File, error)` recursive-descent parser | F4 |
| AST construction | F4 |
| Syntax validation | F4 |
| IF/FOR/WHILE/REPEAT/CASE top-level block markers | Handled implicitly because they only exist inside an outer BEGIN..END |
| UTF-16 byte offset conversion | Bytebase consumer does this on its own |
| `DELIMITER` directive | MySQL-specific, not in Snowflake |
| Catalog, completion, semantic, deparse | Separate DAG branches |

## Acceptance Criteria

F3 is complete when:

1. `go build ./snowflake/parser/...` succeeds.
2. `go vet ./snowflake/parser/...` clean.
3. `gofmt -l snowflake/parser/` clean.
4. `go test ./snowflake/parser/...` passes — all F2 tests still pass, plus the new F3 tests.
5. All 27 files in `snowflake/parser/testdata/legacy/` produce at least one non-empty segment from `Split`.
6. `task_scripting.sql` specifically returns exactly **2 segments** (verifying BEGIN..END handling).
7. `go build ./snowflake/...` (isolation check) — F3 builds without depending on any DAG node beyond F1+F2.
8. After merge, `docs/migration/snowflake/dag.md` F3 status is flipped to `done`.

## Files Created

```
snowflake/parser/split.go
snowflake/parser/split_test.go
docs/superpowers/specs/2026-04-09-snowflake-splitter-design.md  (this file)
```

Estimated total: ~550 LOC across 2 code files plus this design document.
