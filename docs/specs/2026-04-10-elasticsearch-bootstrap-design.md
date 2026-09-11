# F0: Elasticsearch Bootstrap — Design Spec

**DAG node:** F0 (bootstrap)
**Branch:** `feat/elasticsearch/bootstrap`
**Engine:** elasticsearch
**Depends on:** nothing
**Unblocks:** F1 (parser core), F2 (classification)

## Goal

Create the `elasticsearch/` package skeleton, add the hjson dependency, lift the 8 YAML golden files from bytebase, and build the test harness that every subsequent DAG node uses. After this node lands, F1 and F2 can start work immediately.

## Package Layout

```
elasticsearch/
  parse.go                  # Public types: Request, ParseResult, SyntaxError
  parser/                   # Empty directory (F1 populates)
  parsertest/
    helpers_test.go         # Generic YAML loader + shared test assertions
    smoke_test.go           # Verify all 8 YAML files load without error
    testdata/
      parse-elasticsearch-rest.yaml
      splitter.yaml
      statement-ranges.yaml
      masking_classify_api.yaml
      masking_analyze_body.yaml
      masking_predicate_fields.yaml
      query_span.yaml
      masking_analyze_request.yaml
```

`analysis/` is deferred to F6a. No other subdirectories.

## Components

### 1. `elasticsearch/parse.go` — Public Types

Defines the omni-native types that every feature node returns. These are structurally equivalent to the bytebase `Request` / `ParseResult` types but with no dependency on bytebase's `base` package.

```go
package elasticsearch

// ParseResult holds parsed requests and any syntax errors encountered.
type ParseResult struct {
    Requests []*Request     `yaml:"requests"`
    Errors   []*SyntaxError `yaml:"errors,omitempty"`
}

// Request is a single parsed Kibana Dev Console REST request.
type Request struct {
    Method      string   `yaml:"method"`
    URL         string   `yaml:"url"`
    Data        []string `yaml:"data,omitempty"`
    StartOffset int      `yaml:"startoffset"` // Byte offset of method letter (inclusive)
    EndOffset   int      `yaml:"endoffset"`   // Byte offset past request end (exclusive)
}

// SyntaxError represents a parse error at a specific position.
type SyntaxError struct {
    Position   Position `yaml:"position"`
    Message    string   `yaml:"message"`
    RawMessage string   `yaml:"rawmessage"`
}

// Position represents a location in source text.
type Position struct {
    Line   int `yaml:"line"`   // 1-based
    Column int `yaml:"column"` // 1-based
}
```

No function implementations. F1 adds `ParseElasticsearchREST()`, F3 adds `SplitMultiSQL()`, etc.

YAML tags must match the bytebase golden file format exactly. Verified against `parse-elasticsearch-rest.yaml`: requests use `method`/`url`/`data`/`startoffset`/`endoffset`; errors use `position: {line, column}` + `message` + `rawmessage`.

### 2. `elasticsearch/parsertest/helpers_test.go` — Test Harness

Generic YAML loader using Go generics:

```go
package parsertest

import (
    "os"
    "testing"

    "gopkg.in/yaml.v3"
)

// loadYAML loads a YAML file from testdata/ and unmarshals into a slice of T.
func loadYAML[T any](t *testing.T, filename string) []T {
    t.Helper()
    data, err := os.ReadFile("testdata/" + filename)
    if err != nil {
        t.Fatalf("read %s: %v", filename, err)
    }
    var cases []T
    if err := yaml.Unmarshal(data, &cases); err != nil {
        t.Fatalf("unmarshal %s: %v", filename, err)
    }
    return cases
}
```

Each feature node's test file defines its own case struct with YAML tags matching that golden file's schema, then calls `loadYAML[thatStruct](t, "filename.yaml")`.

The `record` pattern (regenerating goldens from actual output) is preserved: each feature test file can set `record := false` and flip it to `true` during development to update goldens.

### 3. `elasticsearch/parsertest/smoke_test.go` — Smoke Test

Schema-agnostic verification that all 8 YAML files load without unmarshal errors:

```go
func TestGoldenFilesLoadable(t *testing.T) {
    files := []string{
        "parse-elasticsearch-rest.yaml",
        "splitter.yaml",
        "statement-ranges.yaml",
        "masking_classify_api.yaml",
        "masking_analyze_body.yaml",
        "masking_predicate_fields.yaml",
        "query_span.yaml",
        "masking_analyze_request.yaml",
    }
    for _, f := range files {
        t.Run(f, func(t *testing.T) {
            loadYAML[map[string]any](t, f)
        })
    }
}
```

This catches YAML syntax errors, file corruption, or missing files. No parser logic exercised — just "can we read and unmarshal the golden files?"

### 4. `elasticsearch/parser/` — Empty Directory

Created with a minimal `doc.go` placeholder:

```go
// Package parser implements a hand-written recursive-descent parser for
// Kibana Dev Console-style Elasticsearch REST request blocks.
package parser
```

F1 populates this with the actual state machine.

### 5. `go.mod` Change

Add `github.com/hjson/hjson-go/v4` as a direct dependency:

```
go get github.com/hjson/hjson-go/v4
```

This dependency is consumed by F1's parser implementation, but adding it at bootstrap ensures the module is available immediately when F1 starts.

## Golden Files

All 8 YAML files are lifted verbatim from `bytebase/backend/plugin/parser/elasticsearch/test-data/`. No modifications. The filenames are preserved exactly (including the hyphen vs underscore inconsistency between `parse-elasticsearch-rest.yaml` and `masking_classify_api.yaml` — matching bytebase is more important than consistency).

Each file maps to a DAG node:

| File | Owning DAG node | Schema (top-level array element shape) |
|------|-----------------|----------------------------------------|
| `parse-elasticsearch-rest.yaml` | F1 | `{description, statement, result: {requests[], errors[]}}` |
| `splitter.yaml` | F3 | `{description, statement, expectedCount, statements[]}` |
| `statement-ranges.yaml` | F4 | `{description, statement, ranges[]}` |
| `masking_classify_api.yaml` | F6a | `{description, method, url, result}` |
| `masking_analyze_body.yaml` | F6a | `{description, method, url, body, result}` |
| `masking_predicate_fields.yaml` | F6b | `{description, body, predicateFields[]}` |
| `query_span.yaml` | F7 | `{description, statement, querySpan}` |
| `masking_analyze_request.yaml` | F7 | `{description, method, url, body, result}` |

## What F0 Does NOT Do

- No parser implementation (F1)
- No `analysis/` directory (F6a)
- No function bodies in `parse.go` — only type definitions
- No `init()` registration (bytebase-side, M1)
- No `catalog/`, `completion/`, `semantic/`, `quality/`, or `deparse/` packages (not needed for this engine)

## Success Criteria

1. `go build ./elasticsearch/...` compiles cleanly
2. `go test ./elasticsearch/parsertest/...` passes — all 8 YAML files load without error
3. `go vet ./elasticsearch/...` clean
4. YAML golden files byte-identical to bytebase originals
5. No new test failures in any existing package (`go test ./mongo/parsertest/...` still green)
