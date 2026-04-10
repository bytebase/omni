# Elasticsearch Parser Core (F1) -- Implementation Plan

## Overview

Port the core Elasticsearch REST request parser from bytebase to omni. This is a faithful, mechanical port -- no refactoring, no simplification, no behavioral changes.

## Steps

### Step 1: Create `elasticsearch/parser/parser.go`

Port the entire state machine from `bytebase/bytebase/backend/plugin/parser/elasticsearch/parser.go`. This includes:

- `state` struct (renamed from `parser` to avoid package/type name collision)
- `newState()` constructor
- `Parse()` public entry point that creates a state and runs `parse()`
- All state machine methods: `parse()`, `multiRequest()`, `request()`, `method()`, `url()`, `object()`, `value()`, `word()`, `number()`, `array()`, `string()`
- Navigation methods: `next()`, `nextEmptyInput()`, `nNextEmptyInput()`, `nextOneOf()`, `reset()`, `peek()`
- Whitespace/comment methods: `white()`, `strictWhite()`, `newLine()`, `comment()`
- Request tracking: `addRequestStart()`, `addRequestEnd()`, `updateRequestEnd()`
- String utilities: `nextUpTo()`, `includes()`
- Supporting functions exported for use by `parse_impl.go`:
  - `GetAdjustedParsedRequest()` -- byte offset to line number mapping
  - `GetEditorRequest()` -- extract method/URL/data from text
  - `SplitDataIntoJSONObjects()` -- multi-document body splitting
  - `ParseLine()` -- extract method and URL from a request line
  - `RemoveTrailingWhitespace()` -- strip trailing comments from URL
  - `CollapseLiteralString()` -- normalize triple-quoted literals
  - `IndentData()` -- HJSON round-trip for comment stripping
  - `ContainsComments()` -- detect comments outside string literals

Adaptations from original:
- Replace `github.com/pkg/errors` with stdlib `fmt.Errorf`
- Remove `storepb` and `base` imports (bytebase-specific)
- Export types and functions that need to be accessed from `parse_impl.go` and tests
- Rename `parser` struct to `state` to avoid name collision with package name

### Step 2: Create `elasticsearch/parse_impl.go`

Add the public `ParseElasticsearchREST(text string) (*ParseResult, error)` function that:

1. Calls `parser.Parse(text)` to get raw parse result
2. Iterates over `parser.ParsedRequest` results
3. For each request, computes `AdjustedParsedRequest` via `GetAdjustedParsedRequest()`
4. Extracts `EditorRequest` via `GetEditorRequest()`
5. Applies comment stripping (`ContainsComments` + `IndentData`) and literal string normalization (`CollapseLiteralString`) to body data
6. Maps to `elasticsearch.Request` with original byte offsets
7. Converts `parser.SyntaxError` (byte offset) to `elasticsearch.SyntaxError` (line/column Position) using the incremental walk algorithm from the original
8. Returns `*ParseResult` with requests and errors

### Step 3: Create `elasticsearch/parsertest/parser_test.go`

Four test functions:

1. `TestParseElasticsearchREST` -- load `parse-elasticsearch-rest.yaml` via `loadYAML`, iterate cases, compare with `require.Equal`. Includes `record` flag for golden regeneration.
2. `TestParse` -- inline test cases for internal `ParsedRequest` offsets (ported from bytebase's `TestParse`).
3. `TestGetEditorRequest` -- inline test cases for editor request extraction (ported from bytebase's `TestGetEditorRequest`).
4. `TestContainsComments` -- inline test cases for comment detection (ported from bytebase's `TestContainsComments`).

### Step 4: Write design spec and implementation plan

Create:
- `docs/superpowers/specs/2026-04-10-elasticsearch-parser-core-design.md`
- `docs/superpowers/plans/2026-04-10-elasticsearch-parser-core.md`

### Step 5: Verify

- `go build ./elasticsearch/...` compiles cleanly
- `go test ./elasticsearch/parsertest/` -- all tests pass
- `go build ./...` -- full repo still compiles

## Verification checklist

- [ ] `parser/parser.go` compiles with no bytebase dependencies
- [ ] All 5 golden test cases in `parse-elasticsearch-rest.yaml` pass
- [ ] All 6 `TestParse` inline cases pass (offset tracking)
- [ ] All 7 `TestGetEditorRequest` inline cases pass
- [ ] All 7 `TestContainsComments` inline cases pass
- [ ] Full repo `go build ./...` succeeds
- [ ] No behavioral differences from original parser
