# Snowflake Syntax Diagnostics — Implementation Plan (T3.4)

Date: 2026-04-15
Status: In progress

---

## Steps

### Step 1 — Core types + Analyze (diagnostics.go)

File: `snowflake/diagnostics/diagnostics.go`

- Define `Severity` (int, Error/Warning/Info constants + String())
- Define `Position` struct (Line, Column, Offset)
- Define `Range` struct (Start, End Position)
- Define `Diagnostic` struct (Severity, Range, Source, Message)
- Implement `Analyze(sql string) []Diagnostic`:
  - Call `parser.ParseBestEffort(sql)`
  - Build `parser.NewLineTable(sql)`
  - Convert each `ParseError` → `Diagnostic`
  - Return nil for zero-error inputs

### Step 2 — Tests (diagnostics_test.go)

File: `snowflake/diagnostics/diagnostics_test.go`

- Test: empty input → nil
- Test: whitespace-only input → nil
- Test: valid SELECT → nil
- Test: single-line syntax error → line=1, correct column, SeverityError
- Test: multi-line SQL, error on line 3 → line=3
- Test: multiple errors → correct count
- Test: Source == "snowflake-parser"
- Test: Message non-empty
- Test: UTF-8 input (multibyte chars before error) → byte-based column

### Step 3 — Run tests + fmt

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/snowflake-diagnostics
go test ./snowflake/... -count=1
gofmt -w snowflake/diagnostics/
```

### Step 4 — Commit + push + PR

Commit 1: docs — spec + plan
Commit 2: feat — types + Analyze implementation
Commit 3: test — diagnostics_test.go

---

## Dependency tree

- diagnostics.go depends on `snowflake/parser` (ParseBestEffort, NewLineTable, ParseError)
- diagnostics_test.go depends on diagnostics.go only
