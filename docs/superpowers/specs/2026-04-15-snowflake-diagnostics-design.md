# Snowflake Syntax Diagnostics — Design Spec (T3.4)

Date: 2026-04-15
Author: Claude (agent)
Status: Accepted

---

## Context

T3.4 adds a thin reporting layer on top of the Snowflake parser. Its purpose is to shape
raw parse errors into structured diagnostics suitable for editor integration (LSP/language
server), CI linting, and problem-panel display. The parser already collects all errors via
`ParseBestEffort`; this node just wraps that output in a well-typed, editor-friendly form.

---

## Goals

1. A `Diagnostic` value type with severity, source, message, and a text range.
2. A `Range` / `Position` type pair matching the LSP `Range` concept (line+column+offset).
3. An `Analyze(sql string) []Diagnostic` function as the single public entry point.
4. Tests covering: syntax errors, multi-line input, UTF-8, empty/whitespace input, valid SQL.

---

## Non-Goals

- No ASP walker to detect semantic warnings — that is the advisor's responsibility (T2.6+).
- No configuration, suppression, or rule IDs — this is pure syntax diagnostics only.
- No changes to the parser itself.

---

## Package layout

New package `snowflake/diagnostics/`:

```
snowflake/diagnostics/
  diagnostics.go       — all types + Analyze function
  diagnostics_test.go  — tests
```

A separate package gives a clean API surface (`diagnostics.Analyze`, `diagnostics.Diagnostic`)
that is easier for consumers (e.g. bytebase's language server) to import without pulling in
all parser internals.

---

## Key types

### Severity

```go
type Severity int

const (
    SeverityError   Severity = iota // syntax error — the SQL will not execute
    SeverityWarning                  // reserved for future use
    SeverityInfo                     // reserved for future use
)
```

Lower index = higher severity (matches typical lint tool conventions).

### Position

```go
type Position struct {
    Line   int // 1-based line number
    Column int // 1-based column (byte offset within line)
    Offset int // 0-based byte offset within the source
}
```

`Offset` allows callers to reconstruct the raw position without a line table; `Line`/`Column`
are pre-computed for convenience.

### Range

```go
type Range struct {
    Start Position
    End   Position
}
```

Matches the LSP `Range` shape. `End` is exclusive (same convention as `ast.Loc.End`).

### Diagnostic

```go
type Diagnostic struct {
    Severity Severity
    Range    Range
    Source   string // always "snowflake-parser"
    Message  string
}
```

---

## Analyze algorithm

1. Call `parser.ParseBestEffort(sql)` to obtain `*parser.ParseResult`.
2. Build `parser.NewLineTable(sql)`.
3. For each `parser.ParseError` in `result.Errors`:
   a. Look up `(line, col) = lt.Position(err.Loc.Start)` for the start.
   b. Look up `(line, col) = lt.Position(err.Loc.End)` for the end.
   c. If `Loc.End < 0` (unknown), use `Start` for both start and end.
   d. Emit `Diagnostic{SeverityError, Range{start, end}, "snowflake-parser", err.Msg}`.
4. Return the slice (nil if empty).

No AST walk is needed: the parser already collects all lex and parse errors via `ParseBestEffort`.

---

## Edge cases

- **Empty input**: `ParseBestEffort("")` returns an empty error slice → `Analyze` returns nil.
- **Whitespace-only input**: same as empty — no errors.
- **Valid SQL**: no parse errors → nil result.
- **Unknown `Loc.End`**: use `Start` for end (zero-width underline at the error position).
- **Multibyte UTF-8**: `LineTable.Position` uses byte offsets — column is a byte count, not
  a rune count. This matches the parser's byte-based tokenization.

---

## Testing strategy

1. Single-line syntax error — verify line=1, correct column range, SeverityError.
2. Multi-line SQL, error on line 3 — verify line=3 and correct column.
3. Multiple errors in one input — verify count and individual positions.
4. Valid SQL — zero diagnostics.
5. Empty string — zero diagnostics.
6. Whitespace-only string — zero diagnostics.
7. UTF-8 identifier in error context — verify byte-based column.
8. Source field — always "snowflake-parser".
9. Message is non-empty for all diagnostics.

---

## Alternatives considered

**Place in `snowflake/parser/diagnostics.go`**
Simpler (no new package), but mixes parsing internals with the reporting API.
Separate package is cleaner for consumers and follows the advisor pattern (`snowflake/advisor`).

**Return `([]Diagnostic, error)`**
Unnecessary — all errors are already captured in the diagnostics slice.
A plain slice return keeps the API minimal.

**Walk AST for `Invalid*` nodes**
The Snowflake AST has no `Invalid*` node types at this stage. All error recovery markers
are surfaced as `ParseError` entries. The walk can be added later if needed.
