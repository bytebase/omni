# F2: Elasticsearch Request Classification — Design Spec

**DAG node:** F2 (request classification)
**Branch:** `feat/elasticsearch/classification`
**Engine:** elasticsearch
**Depends on:** F0 (bootstrap)
**Unblocks:** F6a (body analysis foundation)
**Parallel with:** F1 (parser core), F3, F4, F5

## Goal

Port `ClassifyRequest(method, url string)` from `bytebase/backend/plugin/parser/elasticsearch/query_type.go`
into `omni/elasticsearch/querytype.go`. This function maps an HTTP method + URL path to a `QueryType`
value (Select / SelectInfoSchema / DML / DDL / Explain / Unknown). It is pure data + switch logic —
no parser state, no body analysis, no external dependencies beyond `strings`.

## QueryType

Define an omni-native `QueryType` type in `elasticsearch/querytype.go`. Equivalent to `base.QueryType`
in bytebase but with no dependency on bytebase's `base` package.

```go
// QueryType classifies an Elasticsearch REST request for cost and routing decisions.
type QueryType int

const (
    QueryTypeUnknown       QueryType = iota
    QueryTypeSelect                  // Read-only data query
    QueryTypeSelectInfoSchema        // Read-only cluster/node metadata
    QueryTypeDML                     // Document write (index, update, delete)
    QueryTypeDDL                     // Index/alias/mapping management
    QueryTypeExplain                 // Query explanation
)
```

Matching bytebase's integer values is not required — bytebase's glue file will translate between the
two type systems at M1 time. What matters is that the five semantic categories are faithfully
represented.

## URL Pattern Lists

Three lists are ported verbatim from bytebase without refactoring, reordering, or improvement:

- `readOnlyPostEndpoints` — POST endpoints that are read-only (search, count, etc.)
- `infoSchemaEndpoints` — cluster/node metadata endpoints
- `dmlPostEndpoints` — document write endpoints via POST

The constraint "verbatim" means: same strings, same order, same comments, same variable names.

## ClassifyRequest

The decision tree mirrors bytebase exactly:

| Method | Condition | QueryType |
|--------|-----------|-----------|
| HEAD | — | Select |
| GET | URL contains info schema pattern | SelectInfoSchema |
| GET | otherwise | Select |
| DELETE | URL contains `_doc/` | DML |
| DELETE | otherwise | DDL |
| PUT | URL contains `_doc/`, `_doc`, `_create/`, or `_bulk` | DML |
| PUT | otherwise | DDL |
| PATCH | — | DML |
| POST | URL contains `_explain/` or ends with `_explain` | Explain |
| POST | URL contains info schema pattern | SelectInfoSchema |
| POST | URL contains read-only POST pattern | Select |
| POST | URL contains DML POST pattern | DML |
| POST | otherwise | DDL |
| other | — | Unknown |

Method matching is case-insensitive (normalized to upper). URL matching is case-insensitive
(normalized to lower). Matching is substring (`strings.Contains`) or suffix (`strings.HasSuffix`)
as in the original.

## Test Strategy

No YAML golden file exists for `ClassifyRequest` — bytebase tests it with inline Go table-driven
cases in `query_type_test.go`. Since that file does not exist in the bytebase repo (the test file
was not present), port tests that exercise the full decision tree:

- Each HTTP method
- InfoSchema vs non-InfoSchema GET
- Document vs index-level DELETE and PUT
- All POST branches (explain, info schema, read-only, DML, default DDL)
- Unknown method

Tests live in `elasticsearch/parsertest/querytype_test.go` using Go table-driven style.

## Files

| Action | Path |
|--------|------|
| Create | `elasticsearch/querytype.go` |
| Create | `elasticsearch/parsertest/querytype_test.go` |

## Success Criteria

1. `go build ./elasticsearch/...` compiles cleanly
2. `go test ./elasticsearch/...` passes
3. `go test ./elasticsearch/parsertest/...` passes (including existing smoke test)
4. `go vet ./elasticsearch/...` clean
