# Elasticsearch Parser Core (F1) -- Design Spec

## Goal

Port the core Elasticsearch REST request parser from `bytebase/bytebase/backend/plugin/parser/elasticsearch/parser.go` (~1262 LoC) to `omni/elasticsearch/parser/parser.go`. This is the foundation node (F1) of the Elasticsearch migration DAG -- every other feature node depends on the parse output.

## Scope

### In scope

1. **State machine parser** -- the hand-written recursive-descent parser that consumes Kibana Dev Console REST request blocks and produces byte-offset-tracked request ranges + syntax errors.
2. **Public entry point** -- `ParseElasticsearchREST(text string) (*ParseResult, error)` that maps internal parser output to the omni-level types defined in `elasticsearch/parse.go`.
3. **Supporting functions** -- all functions from the original parser: `GetAdjustedParsedRequest`, `GetEditorRequest`, `SplitDataIntoJSONObjects`, `ParseLine`, `RemoveTrailingWhitespace`, `CollapseLiteralString`, `IndentData`, `ContainsComments`.
4. **Test suite** -- golden file tests against `parse-elasticsearch-rest.yaml`, inline `TestParse` cases for internal offset tracking, `TestGetEditorRequest` cases, `TestContainsComments` cases.

### Out of scope

- Statement splitting (F3), statement ranges (F4), diagnostics (F5), query type classification (F6), masking analysis (F7), query span (F8) -- these are separate DAG nodes that wrap the parser output.
- Typed body AST -- bodies remain `map[string]any` per the analysis decision.
- Kibana TS upstream alignment -- we freeze on the existing Go parser behavior.

## Architecture

### Package layout

```
elasticsearch/
  parse.go          -- types: ParseResult, Request, SyntaxError, Position (F0, already exists)
  parse_impl.go     -- ParseElasticsearchREST() public entry point (F1, new)
  parser/
    doc.go          -- package doc (F0, already exists)
    parser.go       -- core state machine + all supporting functions (F1, new)
  parsertest/
    helpers_test.go -- loadYAML/writeYAML (F0, already exists)
    smoke_test.go   -- golden file loadability (F0, already exists)
    parser_test.go  -- parser tests (F1, new)
    testdata/
      parse-elasticsearch-rest.yaml -- golden test data (F0, already exists)
```

### Why two packages?

The `parser` sub-package contains the state machine and all internal types (`ParsedRequest`, `SyntaxError` with byte offset, `AdjustedParsedRequest`, `EditorRequest`). The parent `elasticsearch` package owns the public types (`Request`, `SyntaxError` with line/column `Position`) and the public API. This separation keeps the internal parser types from polluting the public namespace and allows future DAG nodes to import `elasticsearch` without pulling in parser internals.

### Type mapping

| Original (bytebase) | Internal (omni parser/) | Public (omni elasticsearch/) |
|---|---|---|
| `parsedRequest` | `parser.ParsedRequest` | -- (internal only) |
| `syntaxError` | `parser.SyntaxError` | `elasticsearch.SyntaxError` |
| `adjustedParsedRequest` | `parser.AdjustedParsedRequest` | -- (internal only) |
| `editorRequest` | `parser.EditorRequest` | -- (internal only) |
| `Request` | -- | `elasticsearch.Request` |
| `ParseResult` | `parser.ParseResult` | `elasticsearch.ParseResult` |

### Error conversion

The original code converts byte-offset errors to 1-based line/column positions using an incremental walk. The omni port preserves this exact algorithm in `parse_impl.go`, mapping `parser.SyntaxError` (byte offset + message) to `elasticsearch.SyntaxError` (Position{Line, Column} + Message + RawMessage).

The `RawMessage` field is always empty (matches original behavior).

## Key behavioral properties preserved

1. **Byte-offset accounting** -- three parallel position systems: `parser.ParsedRequest.StartOffset/EndOffset` (byte offset), `parser.AdjustedParsedRequest.StartLineNumber/EndLineNumber` (0-based line number), and `elasticsearch.SyntaxError.Position` (1-based line/column). All are preserved exactly.

2. **Multi-document JSON splitting** -- `SplitDataIntoJSONObjects` uses brace-depth tracking that respects string boundaries and escape sequences.

3. **Error recovery** -- regex-driven resync to next `^(POST|HEAD|GET|PUT|DELETE|PATCH)` method keyword. Yields valid requests + sorted error list even on malformed input.

4. **Comment handling** -- line (`#`, `//`) and block (`/* */`) comments stripped during whitespace scanning.

5. **HJSON tolerance** -- body round-tripped through `hjson-go/v4` for comment stripping and normalization via `IndentData`.

6. **Triple-quoted literals** -- `"""..."""` preserved as-is, normalized via `CollapseLiteralString`.

## Dependencies

| Dependency | Purpose |
|---|---|
| `github.com/hjson/hjson-go/v4` | HJSON tolerance + comment-strip round-trip |
| `encoding/json` | JSON marshaling for `CollapseLiteralString` and `IndentData` |
| `regexp` | Error recovery resync pattern |
| `unicode/utf8` | Rune-level byte offset tracking |
| `strconv` | Number and hex parsing |
| `strings` | String manipulation throughout |
| `slices` | Error sorting |

No `github.com/pkg/errors` -- uses stdlib `fmt.Errorf` per omni convention.

## Risks

1. **Behavioral drift** -- the port is mechanical, but any subtle difference in error recovery or offset accounting could cause downstream DAG nodes (splitter, statement ranges) to produce wrong results. Mitigation: comprehensive golden tests + inline offset tests.

2. **hjson library version** -- `hjson-go/v4` output formatting must match what the original parser produced. The same library version is used.
