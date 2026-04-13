# Elasticsearch Migration Analysis

Migration target: full coverage of the existing hand-written Go parser at `bytebase/bytebase/backend/plugin/parser/elasticsearch` by an equivalent (still hand-written) parser under `omni/elasticsearch`. Bytebase consumption is used to **prioritize** the order of work; it is never used to **scope down** existing parser features.

> **Consumption-only mode.** Unlike other engines in this migration, Elasticsearch has no antlr4 grammar in `bytebase/parser` to map. The existing parser is itself hand-written Go and lives inside `bytebase/bytebase`. This document therefore reflects the **existing parser's surface** rather than a grammar enumeration. The "spec" the existing parser approximates is Kibana's TypeScript console parser (`elastic/kibana/.../requests_utils.ts`); whether to track Kibana drift or freeze on current parser behavior is an open decision noted at the end of this document.

Sources analyzed:

- Existing parser: `/Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch` (8 Go files, 7 test-data YAML golden files, ~3,022 LoC including tests, ~2,156 LoC non-test)
- Downstream consumer: `/Users/h3n4l/OpenSource/bytebase` (3 import sites total — extremely narrow surface)

---

## Parse Surface

### What the parser accepts

Kibana Dev Console-style REST request blocks. Each "statement" is an HTTP method + URL path + optional JSON / HJSON body, optionally followed by additional JSON objects:

```
GET /my-index/_search
{
  "query": { "match": { "field": "value" } }
}

POST /_bulk
{ "index": { "_index": "test" } }
{ "field1": "value1" }
```

A single input may contain multiple such requests. Comments, blank lines, and HJSON conveniences (unquoted keys, trailing commas, triple-quoted literal strings) are tolerated.

### Top-level features

| Feature | Notes |
|---|---|
| Multi-request blocks | Splitter recognizes new requests by line-anchored `^(POST\|HEAD\|GET\|PUT\|DELETE\|PATCH)` regex. Requests separated by blank lines or back-to-back. |
| Multi-document body | After the first body object, additional `{...}` objects on subsequent lines are appended to the same `Request.Data`. Used for `_bulk`, `_msearch`. |
| Comment handling | Line comments (`#`, `//`) and block comments (`/* … */`); detected outside string literals. When present, body is round-tripped through hjson to strip them. |
| Triple-quoted literals | `"""…"""` preserved as-is (no JSON escaping); normalized via `collapseLiteralString()` before downstream consumers see them. |
| HJSON tolerance | Bodies parsed by `github.com/hjson/hjson-go/v4`, then re-marshaled to canonical JSON. |
| Error recovery | On parse error: record byte-offset error, regex-search forward for next method keyword, reset, continue. Yields valid requests + sorted error list even on malformed input. |
| Offset tracking | Every request carries `StartOffset` / `EndOffset` (byte offsets, 0-indexed, end-exclusive). |
| Position reporting | Two parallel systems: byte offset (UTF-8) for `Request`, and LSP `Position` (UTF-16 code units, BMP surrogate-aware) for `GetStatementRanges`. |
| URL parameter handling | `removeTrailingWhitespace()` strips inline comments from URL line while respecting quoted query-string values like `_search?q="hello world"`. |

### What the parser does NOT do

- It does **not** parse the Elasticsearch Query DSL grammar in any deep sense — bodies are treated as opaque JSON/HJSON for splitting and round-tripping. Field-level inspection happens later, in `masking.go`, by walking `map[string]any`.
- It does **not** parse Elasticsearch SQL (`POST _sql`) — the body is opaque JSON; no SQL-string analysis.
- It does **not** parse KQL (Kibana Query Language) — only the request envelope.
- No formal grammar / EBNF exists. The acceptance set is whatever the current Go state machine + Kibana TS reference behavior define.

---

## Parse API Surface

```go
// parser.go — primary entry point
func ParseElasticsearchREST(text string) (*ParseResult, error)

type ParseResult struct {
    Requests []*Request
    Errors   []*base.SyntaxError
}

type Request struct {
    Method      string   // GET / POST / PUT / DELETE / HEAD / PATCH
    URL         string   // /my-index/_search?param=value
    Data        []string // body objects (one per JSON document)
    StartOffset int      // byte offset of method letter
    EndOffset   int      // byte offset just past request end
}

// query_type.go — request classification
func ClassifyRequest(method, url string) base.QueryType
// → Select / SelectInfoSchema / DML / DDL / Explain / QueryTypeUnknown

// masking.go — body field analysis for masking & lineage
func AnalyzeRequest(method, url, body string) *RequestAnalysis
// Recognized clauses: match, term, terms, range, exists, prefix, wildcard, regexp,
//                     fuzzy, ids, bool (must/should/must_not/filter), nested,
//                     has_child/has_parent, plus _source/fields/highlight/sort.

// Plus exported supporting types & constants:
//   MaskableAPI, BlockedFeature, RequestAnalysis (type aliases over base.*)
//   APIMaskSearch / APIMaskGetDoc / APIMaskGetSource / APIMaskMGet / APIMaskExplain / APIBlocked
//   BlockedFeatureNames map[BlockedFeature]string

// splitter.go — registered with base, satisfies SplitMultiSQLFunc
func SplitMultiSQL(statement string) ([]base.Statement, error)

// query_span.go — registered with base, satisfies GetQuerySpanFunc
func GetQuerySpan(ctx context.Context, ..., stmt base.Statement, ...) (*base.QuerySpan, error)

// statement_ranges.go — registered with base, satisfies StatementRangeFunc
func GetStatementRanges(_ context.Context, _, statement string) ([]base.Range, error)

// diagnose.go — registered with base, satisfies DiagnoseFunc
func Diagnose(_ context.Context, _, statement string) ([]base.Diagnostic, error)
```

### Plugin registration (init-time)

Each of the four `base`-registered functions self-registers in its file's `init()`:

| File | Register call | Engine key |
|---|---|---|
| `splitter.go:11` | `base.RegisterSplitterFunc` | `storepb.Engine_ELASTICSEARCH` |
| `query_span.go:14` | `base.RegisterGetQuerySpan` | `storepb.Engine_ELASTICSEARCH` |
| `statement_ranges.go:14` | `base.RegisterStatementRangesFunc` | `storepb.Engine_ELASTICSEARCH` |
| `diagnose.go:11` | `base.RegisterDiagnoseFunc` | `storepb.Engine_ELASTICSEARCH` |

The blank import `_ "github.com/bytebase/bytebase/backend/plugin/parser/elasticsearch"` in `backend/server/ultimate.go:34` is what triggers the init chain.

---

## Internal AST / data shape

There is no tree AST. Outputs are:

1. `[]*Request` — flat list of HTTP method + URL + body strings + byte offsets.
2. `map[string]any` — body parsed by hjson, walked dynamically inside `masking.go` for predicate field extraction.
3. `*RequestAnalysis` — masking-relevant metadata: API type, index name, predicate paths, sort fields, source configuration, blocked feature flags.

This is a meaningful divergence from other omni engines (snowflake, mongo, partiql) which all have a typed `ast/` package. **The omni elasticsearch package is unlikely to need an `ast/` directory in the same shape.** Whether to introduce one anyway (e.g., for typed body nodes that bytebase can walk safely) is a design decision for the planning phase.

---

## Per-file inventory

| File | LoC | Role |
|---|---|---|
| `parser.go` | 1262 | Core state-machine parser. Multi-request scanning, comment/whitespace skipping, brace-tracked JSON splitting, HJSON normalization, byte-offset accounting, error recovery via regex. |
| `masking.go` | 462 | Walks the parsed body's `map[string]any` to extract predicate paths and sort fields. Recognizes 13+ ES query clause types. Maintains the URL-pattern lists for `MaskableAPI` classification and the blocked-feature catalog (`aggs`, `suggest`, `script_fields`, `runtime_mappings`, `stored_fields`, `docvalue_fields`). |
| `query_type.go` | 156 | Three URL-pattern lists (`readOnlyPostEndpoints` ×34, `infoSchemaEndpoints` ×3, `dmlPostEndpoints` ×8) + method-based switch → `base.QueryType`. Pure data, no state. |
| `query_span.go` | 78 | Thin wrapper: `ParseElasticsearchREST` + `ClassifyRequest` + `AnalyzeRequest`, then maps predicate dot-paths to `base.PathAST` nodes via `dotPathToPathAST`. |
| `splitter.go` | 80 | Wraps `ParseElasticsearchREST` → `[]base.Statement`. Includes `byteOffsetToPosition()` for 1-based line/column from byte offset. |
| `statement_ranges.go` | 71 | Wraps `ParseElasticsearchREST` → `[]base.Range` (LSP dialect). UTF-16 code-unit aware via `getPositionByByteOffset()`. |
| `diagnose.go` | 27 | Thin wrapper: maps `ParseResult.Errors` → `[]base.Diagnostic` via `base.ConvertSyntaxErrorToDiagnostic`. |

### Test data (golden YAML, will become omni golden inputs)

| File | Covers |
|---|---|
| `parse-elasticsearch-rest.yaml` | Core parser cases — single/multi request, multi-document body, comments, literal strings, error recovery, offset tracking. |
| `splitter.yaml` | Statement splitter outputs (text + position ranges). |
| `statement-ranges.yaml` | LSP range outputs with UTF-16 code-unit handling. |
| `query_span.yaml` | Query span / lineage extraction. |
| `masking_analyze_request.yaml` | Top-level `AnalyzeRequest` cases. |
| `masking_classify_api.yaml` | URL/method → `MaskableAPI` classification. |
| `masking_predicate_fields.yaml` | Recursive query-clause walk → predicate field paths. |
| `masking_analyze_body.yaml` | Body-level analysis (sort, fields, _source). |

These golden files are the highest-fidelity spec we have and should drive the test plan for every feature node.

---

## Bytebase Consumption

Only **3 import sites** exist in the entire bytebase codebase. This is by far the narrowest consumption surface of any omni-migrated engine so far.

| File:line | Import style | Calls / uses |
|---|---|---|
| `backend/plugin/db/elasticsearch/elasticsearch.go:28` | named (`parser`) | `parser.ParseElasticsearchREST(statement)` at line 395 — driver parses REST blocks before HTTP execution |
| `backend/api/v1/catalog_masking_elasticsearch_test.go:14` | named (`esparser`) | Type aliases only: `MaskableAPI`, `BlockedFeature`, `RequestAnalysis` |
| `backend/server/ultimate.go:34` | blank (`_`) | Triggers the four `init()` registrations into `base` |

### Feature Dependency Map

| Bytebase feature | Parser API used (direct or via base dispatch) | Data extracted |
|---|---|---|
| **REST execution** (driver) | `ParseElasticsearchREST` (direct call) | `Request{Method, URL, Data}` — passed straight to HTTP client |
| **Statement splitting** (SQL editor preprocessing) | `SplitMultiSQL` (via `base` dispatch) | `[]base.Statement` with position ranges per request |
| **Query type detection** (cost & routing) | `ClassifyRequest` via `GetQuerySpan` | `base.QueryType` enum |
| **Lineage / query span** (masking pre-check, audit) | `GetQuerySpan` → `ParseElasticsearchREST` + `AnalyzeRequest` | `QuerySpan.PredicatePaths`, `QuerySpan.ElasticsearchAnalysis` (API type, index, blocked features, source config) |
| **Result masking** (post-execution) | Reads `ElasticsearchAnalysis` from query span; does NOT call parser directly | Drives `maskElasticsearchHitsColumn`, `maskElasticsearchMGetSource`, etc. |
| **IDE diagnostics** (syntax errors) | `Diagnose` (via `base` dispatch) | `[]base.Diagnostic` with line/column |
| **IDE statement ranges** (highlighting / nav) | `GetStatementRanges` (via `base` dispatch) | `[]base.Range` (LSP, UTF-16) |

### Active vs dead code

**Every exported function is consumed.** No P1-due-to-non-use exists. The priority gradient must come from dependency order and risk, not from a consumption gap.

### Bytebase-side test coverage

- `catalog_masking_elasticsearch_test.go` (317 LoC) — exercises masking output (consumes `RequestAnalysis`) but does not exercise parser entry points end-to-end.
- `db/elasticsearch/elasticsearch_test.go` — driver auth/HTTP only, no parser flow.
- **Gap:** no bytebase-side end-to-end test of `parse → split → query span → mask`. The parser's own `*_test.go` + YAML goldens are the integration surface.

---

## Dependencies

| External | Purpose | Reuse in omni? |
|---|---|---|
| `github.com/hjson/hjson-go/v4` | HJSON tolerance + comment-strip round-trip | Yes — port verbatim. omni already pulls it for nothing else; new module dep. |
| `github.com/pkg/errors` | Error wrapping | omni convention may differ — check existing engines. |
| `github.com/bytebase/lsp-protocol` | LSP `Position` for `GetStatementRanges` | omni already uses for other engines. |
| stdlib `encoding/json`, `regexp`, `strings`, `strconv`, `slices`, `unicode/utf8` | — | — |

| Internal (bytebase) | omni equivalent |
|---|---|
| `backend/plugin/parser/base` (Splitter / QuerySpan / Diagnose / StatementRanges register hooks; types `Statement`, `QuerySpan`, `SyntaxError`, `QueryType`, `Range`, `Diagnostic`, `PathAST`, `ConvertSyntaxErrorToDiagnostic`) | omni's `base` package — confirm all needed symbols exist when implementing each tier |
| `backend/generated-go/store` (`storepb.Engine_ELASTICSEARCH`, `storepb.Position`) | Out of scope for omni. Engine ID dispatch lives entirely on the bytebase side; omni exposes functions, bytebase decides how to register them. |
| `base.MaskableAPI`, `base.BlockedFeature`, `base.RequestAnalysis` (elasticsearch-specific types in bytebase's `base`) | Stay in bytebase's `base` package. omni produces its own internal analysis shape; the bytebase glue file converts via an adapter (mongo pattern — see "Bytebase Glue Rewrite" below). |

---

## Complexity Hot Spots

Ordered by risk, highest first:

1. **Byte-offset / UTF-16 accounting.** Three parallel position systems (byte offset → 1-based line/col → LSP UTF-16 code units) with surrogate-pair handling for non-BMP characters. Off-by-one bugs here are silent until they hit production input. → Drives the test plan for splitter, statement_ranges, and parser error positions.
2. **Multi-document JSON splitting** (`splitDataIntoJSONObjects`, parser.go:345-393). Brace-depth tracking that must respect string boundaries and escape sequences. Edge cases in `_bulk` / `_msearch` payloads.
3. **Masking predicate extraction** (`masking.go:388-461`). Recursive walk over dynamic `map[string]any` recognizing 13+ ES query clauses. Missing a clause type silently breaks lineage.
4. **Error recovery state machine** (parser.go:474-545). Regex-driven resync to next method keyword while preserving precise error positions and not losing valid trailing requests.
5. **Comment / HJSON normalization** (parser.go:263-296). Detection-then-round-trip via hjson library; behavioral fidelity depends on hjson library's output stability.
6. **Triple-quoted literal preservation** (parser.go:863-871). Must round-trip identically through `collapseLiteralString` and back.

Lower risk:

- `ClassifyRequest` (pure data + switch).
- Thin wrappers (`Diagnose`, `query_span.go`).
- LSP `Range` construction (modulo UTF-16 hot spot above).

---

## Priority Ranking & Full Coverage Target

Because every public function is consumed, the typical P0/P1 split collapses. Instead, the migration is sequenced by **dependency order** and **risk**:

| Tier | Functions | Rationale | Driving consumer |
|---|---|---|---|
| **Foundation** | `ParseElasticsearchREST` + `Request` / `ParseResult` types + offset machinery | Everything else wraps this; cannot be deferred. Highest implementation risk (state machine + offsets). | All. |
| **Core wrappers** | `SplitMultiSQL`, `GetStatementRanges`, `Diagnose` | Pure transformations of `ParseResult`. Low logic risk; main hazard is UTF-16 in `GetStatementRanges`. | IDE preprocessing + bytebase splitter dispatch. |
| **Classification** | `ClassifyRequest` + URL pattern tables | Independent of body parsing; can be implemented in parallel with foundation. | Cost / routing / query span. |
| **Masking analysis** | `AnalyzeRequest`, `RequestAnalysis`, blocked-feature catalog, predicate walker, `MaskableAPI` constants | Built on top of the parsed body and `ClassifyRequest`. Highest single-feature complexity (clause walker). | Masking pre-check & post-exec column masking. |
| **Query span integration** | `GetQuerySpan`, `dotPathToPathAST` | Trivial wiring once classification + masking exist. | Lineage. |
| **Plugin registration** | Four `init()` blocks → `base.Register*Func` | Last step before downstream import swap. | Server bootstrap. |
| **Driver swap** | Update `backend/plugin/db/elasticsearch/elasticsearch.go:28` import + `backend/server/ultimate.go:34` blank import | Out of scope for omni; lives in bytebase repo. Tracked under BYT-9000. | — |

### Full coverage target (everything that must exist in omni)

| Existing parser feature | Tier |
|---|---|
| Multi-request REST block parsing | Foundation |
| Multi-document body inside one request | Foundation |
| Line + block comment handling (`#`, `//`, `/* */`) | Foundation |
| Triple-quoted literal string preservation | Foundation |
| HJSON tolerance + comment-strip round-trip | Foundation |
| Error recovery via method-keyword regex resync | Foundation |
| Byte-offset (`StartOffset`/`EndOffset`) tracking | Foundation |
| 1-based line/column reporting in syntax errors | Foundation |
| LSP UTF-16 code-unit `Position` reporting | Core wrapper (`GetStatementRanges`) |
| Statement splitting → `base.Statement` array | Core wrapper (`SplitMultiSQL`) |
| Diagnostic conversion → `base.Diagnostic` array | Core wrapper (`Diagnose`) |
| URL pattern → `base.QueryType` (Select/SelectInfoSchema/DML/DDL/Explain) | Classification |
| `MaskableAPI` classification (Search/GetDoc/GetSource/MGet/Explain/Blocked) | Masking |
| Blocked feature detection (aggs, suggest, script_fields, runtime_mappings, stored_fields, docvalue_fields) | Masking |
| `_source` / `fields` / `highlight` / `sort` extraction | Masking |
| Predicate field walk over 13+ query clause types (match, term, terms, range, exists, prefix, wildcard, regexp, fuzzy, ids, bool[must/should/must_not/filter], nested, has_child/has_parent) | Masking |
| Dot-path → `base.PathAST` chain (`dotPathToPathAST`) | Query span |
| `ElasticsearchAnalysis` payload on `QuerySpan` | Query span |
| All 7 test-data YAML golden files passing | Cross-cutting |

---

## Resolved Decisions

All planning-blocking questions have been answered before entering the planning phase:

1. **Spec target — freeze on the existing Go parser.** Treat `bytebase/bytebase/backend/plugin/parser/elasticsearch/` plus its 7 YAML golden files as the literal specification. Do not re-align with current Kibana TS upstream during this migration; that can be a separate follow-up. Migration is a faithful port.
2. **Where do `MaskableAPI` / `BlockedFeature` / `RequestAnalysis` live? — Stay in bytebase's `base` package; omni defines its own internal analysis shape; bytebase's elasticsearch glue file becomes a thin adapter.** This is the same pattern mongo uses today (see `bytebase/backend/plugin/parser/mongodb/masking.go` calling `omniAnalysisToMasking`). Properties of this choice:
   - `bytebase/backend/plugin/parser/base/document_analysis.go` keeps owning the public type contract — no rename, no move.
   - omni's `elasticsearch` package designs its own internal `analysis` types freely, named however is cleanest for omni.
   - `bytebase/backend/plugin/parser/elasticsearch/masking.go` becomes a thin facade: calls into omni for parsing + analysis, converts the result to `*base.RequestAnalysis`, keeps the existing short-name type aliases (`type MaskableAPI = base.MaskableAPI`, etc.) and constant re-exports unchanged.
   - **`backend/api/v1/catalog_masking_elasticsearch_test.go` does not change** — same import line, same type names. The migration is invisible to bytebase consumers.
3. **Body representation — stay with `map[string]any`.** Do not introduce a typed body AST. Rationale: matches the freeze-on-existing decision (#1), Elasticsearch Query DSL is genuinely open-ended (hundreds of clause types, plugin-extensible) so a typed AST forces premature commitment, and the failure mode of `map[string]any` is safer (unknown clauses pass through opaquely, matching existing behavior). The omni `elasticsearch` package likely has no `catalog/`, `completion/`, or `semantic/` directories, so the typed-AST convention's main customers don't exist for this engine.
4. **hjson dependency — accepted.** Add `github.com/hjson/hjson-go/v4` to omni's `go.mod`. Port the comment-strip round-trip logic verbatim from the existing parser.
5. **Engine ID — out of scope.** omni exposes functions; bytebase decides how to dispatch them by engine ID on its side. The omni `elasticsearch` package's public API does not need to know `storepb.Engine_ELASTICSEARCH` exists.

---

## Bytebase Glue Rewrite

Once omni's elasticsearch package is complete, the existing bytebase glue files in `bytebase/backend/plugin/parser/elasticsearch/` must be rewritten as thin adapters that call omni and produce the same `base.*` types they produce today. This step is bytebase-side work, tracked under BYT-9000, and is the very last step of the migration.

**Files to rewrite as adapters:**

| File | Current behavior | Adapter behavior |
|---|---|---|
| `parser.go` | Owns the entire ~1262 LoC state machine + `ParseElasticsearchREST` entry point | Becomes ~30 LoC: re-export `ParseElasticsearchREST` as a thin call into `omni/elasticsearch.Parse` (or equivalent), translating omni's result types into the existing `Request` / `ParseResult` shape that the driver depends on. |
| `splitter.go` | Calls `ParseElasticsearchREST` → builds `[]base.Statement` | Calls `omni/elasticsearch` splitter; returns `[]base.Statement`. Registration block stays. |
| `query_span.go` | Calls `ParseElasticsearchREST` + `AnalyzeRequest` + `dotPathToPathAST` → `base.QuerySpan` | Calls omni's query-span function; converts omni's analysis to `base.RequestAnalysis` via adapter; populates `base.QuerySpan.ElasticsearchAnalysis` and `PredicatePaths`. Registration block stays. |
| `query_type.go` | Owns 45+ URL pattern lists + `ClassifyRequest` | Becomes a thin call into `omni/elasticsearch.ClassifyRequest`. The pattern lists move into omni. |
| `masking.go` | Owns `AnalyzeRequest`, predicate walker, blocked-feature catalog | Becomes the adapter file: calls `omni/elasticsearch` analysis, converts result to `base.RequestAnalysis`, keeps the type aliases and constant re-exports. This is the file that pays the mongo-pattern conversion cost. |
| `statement_ranges.go` | Calls `ParseElasticsearchREST` → builds `[]base.Range` (LSP UTF-16) | Calls omni's statement-ranges function. Registration block stays. |
| `diagnose.go` | Calls `ParseElasticsearchREST` → `[]base.Diagnostic` | Calls omni's diagnose function. Registration block stays. |
| Driver: `backend/plugin/db/elasticsearch/elasticsearch.go:28` | `parser.ParseElasticsearchREST(statement)` | Either keep going through the bytebase glue (preferred — single import surface) or import omni directly. Recommendation: keep going through the bytebase glue so the driver doesn't grow a new dependency. |
| Server bootstrap: `backend/server/ultimate.go:34` | Blank import to fire bytebase glue's `init()` registrations | Unchanged. Bytebase glue's `init()` blocks stay; they're what calls into omni. |

**Notable difference from mongo's glue rewrite:** mongo's bridge work was concentrated in `masking.go` (one file). Elasticsearch's bridge touches **all six** glue files plus the driver, because every one of them currently calls `ParseElasticsearchREST` directly. Per-file the work is small, but it's spread out — plan accordingly.

**No bytebase consumers outside `backend/plugin/parser/elasticsearch/` and the driver need to change.** The masking test (`catalog_masking_elasticsearch_test.go`) and the api/v1 masking layer (`document_masking.go`) interact only with `base.*` types and the existing short-name aliases — both of which are preserved by the adapter pattern.
