# Elasticsearch Migration DAG

Derived from `analysis.md`. Each node is a discrete chunk of work that can be sent to `omni-engine-implementing` (which then delegates to `superpowers:brainstorming` for spec → plan → implement). Nodes that share the same dependency level can run in parallel.

**Architecture decisions** (locked in `analysis.md`):

- Spec target: **freeze on existing Go parser** at `bytebase/bytebase/backend/plugin/parser/elasticsearch/`. Migration is a faithful port. No Kibana TS re-alignment.
- Body representation: **`map[string]any`**, not a typed AST. Failure mode is "unknown clause = opaque pass-through."
- Type ownership: **mongo adapter pattern**. `base.MaskableAPI` / `BlockedFeature` / `RequestAnalysis` stay in bytebase's `base` package; omni defines its own internal `analysis.RequestAnalysis` shape; bytebase's elasticsearch glue files become thin adapters that translate between them.
- hjson dependency (`github.com/hjson/hjson-go/v4`): added to omni's `go.mod`.
- Engine ID dispatch: bytebase-side concern; omni package exposes plain functions, no `Register*Func` init blocks.
- **No `ast/`, `catalog/`, `completion/`, `semantic/`, `quality/`, or `deparse/` packages.** The existing parser doesn't expose any of these concepts and bytebase has no consumer for them. Add later only if a new bytebase feature requires it.

**Package layout** (target):

```
elasticsearch/
  parser/               # Hand-written state-machine REST request parser
    testdata/           # Goldens lifted from bytebase test-data/*.yaml
  parsertest/           # YAML golden runner (table-driven, mongo-style)
  analysis/             # Body walker, AnalyzeRequest, blocked features, predicate paths
  parse.go              # Top-level entry: ParseElasticsearchREST
```

The four wrapper functions (splitter, statement-ranges, diagnose, query-span, classify) live as small files at the top level of the package or in a `wrappers/` subdirectory — exact placement is a brainstorming-time decision.

---

## Nodes

| ID | Node | Package | Depends On | Parallel With | Tier | Status |
|----|------|---------|------------|---------------|------|--------|
| **F0** | bootstrap (skeleton + hjson dep + parsertest harness) | `elasticsearch/`, `elasticsearch/parsertest/` | — | F2 | 0 | **done** (PR #23) |
| **F1** | parser core (state machine, single-request happy path, byte offsets, comments, hjson, error recovery, multi-request, multi-document) | `elasticsearch/parser` | F0 | F2 | 1 | **done** (PR #25) |
| **F2** | request classification (URL pattern tables, `ClassifyRequest`) | `elasticsearch/querytype.go` | F0 | F1, F3, F4, F5, F6a | 1 | **done** (PR #24) |
| **F3** | splitter wrapper (`SplitMultiSQL` → `[]base.Statement` equivalent) | `elasticsearch/` | F1 | F4, F5, F6a, F6b | 2 | **done** (PR #27) |
| **F4** | statement ranges wrapper (LSP UTF-16 `Position`, BMP surrogate-aware) | `elasticsearch/` | F1 | F3, F5, F6a, F6b | 2 | **done** (PR #28) |
| **F5** | diagnose wrapper (parse errors → `base.Diagnostic`) | `elasticsearch/` | F1 | F3, F4, F6a, F6b | 2 | **done** (PR #30) |
| **F6a** | body analysis foundation (`MaskableAPI` URL classification, blocked features catalog, `_source`/`fields`/`highlight`/`sort` extraction, `AnalyzeRequest` skeleton) | `elasticsearch/analysis` | F1, F2 | F3, F4, F5 | 2 | **done** (PR #29) |
| **F6b** | predicate walker (recursive walk over `map[string]any` body recognizing the 13+ ES query clause types — match/term/terms/range/exists/prefix/wildcard/regexp/fuzzy/ids/bool/nested/has_child/has_parent — plus the dot-path field extraction) | `elasticsearch/analysis` | F6a | F3, F4, F5 | 3 | **done** (PR #31) |
| **F7** | query span integration (top-level `GetQuerySpan` orchestration, `dotPathToPathAST`, populating `QuerySpan.PredicatePaths` and the omni-internal analysis payload) | `elasticsearch/analysis` (or `elasticsearch/`) | F1, F2, F6b | — | 4 | not started |
| **M1** | bytebase glue adapter rewrite (all 6 glue files at `bytebase/backend/plugin/parser/elasticsearch/` + driver at `bytebase/backend/plugin/db/elasticsearch/elasticsearch.go:28`; mongo adapter pattern; `init()` blocks stay in bytebase) | bytebase repo | F1–F7 | — | — | not started |

---

## Execution Order

### Wave 0 — Bootstrap (1 node, blocking)

```
F0 (skeleton + hjson dep + parsertest harness)
```

Sets up `omni/elasticsearch/`, adds `github.com/hjson/hjson-go/v4` to `go.mod`, ports the YAML golden test runner from `omni/mongo/parsertest/` (or builds an equivalent). All subsequent nodes depend on the harness existing.

### Wave 1 — Two parallel tracks (parser + classification)

```
F0 ──► F1 (parser core)            ┐
   └──► F2 (request classification) ┴── independent of F1, can start day one
```

F1 is the largest single node in the DAG (~1262 LoC of state machine to port + ~6 YAML golden files driving its tests). F2 is small (~156 LoC of pattern tables + a switch) and is the natural "warm-up" node — it touches no parser state, only URL strings.

### Wave 2 — Parser wrappers + body analysis foundation (4 parallel tracks)

```
F1 ──► F3 (splitter)
F1 ──► F4 (statement ranges)
F1 ──► F5 (diagnose)
F1 + F2 ──► F6a (body analysis foundation)
```

All four can run concurrently once F1 lands. F3 / F4 / F5 are thin wrappers (~80 LoC each in the existing parser); F6a is a medium feature node (~150 LoC including `MaskableAPI` URL classification, blocked features, and the simpler body field extractors).

### Wave 3 — Predicate walker (highest single-feature risk)

```
F6a ──► F6b (predicate walker)
```

This is the highest-complexity feature in the entire migration. Recursive walk over dynamic `map[string]any` recognizing 13+ ES query clause types, with compound (`bool.must`, etc.) and nested (`nested`, `has_child`, `has_parent`) variants. Missing a clause type silently breaks lineage in production. Splitting it out from F6a gives it a dedicated PR with focused review.

Can run in parallel with F3, F4, F5 since the wrappers don't touch analysis.

### Wave 4 — Query span integration

```
F1 + F2 + F6b ──► F7 (query span)
```

Once the predicate walker is complete, query span is mostly wiring: call parser, classify, analyze, convert dot paths to `PathAST`, return.

### M1 — Bytebase glue adapter rewrite (BYT-9000)

After F1–F7 are all complete:

1. Rewrite all 6 glue files at `bytebase/backend/plugin/parser/elasticsearch/` as thin adapters calling into `omni/elasticsearch/...`
2. Implement the `omniAnalysisToBaseRequestAnalysis` bridge function (mongo's `omniAnalysisToMasking` is the template)
3. Keep the existing type aliases (`type MaskableAPI = base.MaskableAPI`, etc.) and constant re-exports unchanged so `catalog_masking_elasticsearch_test.go` and `document_masking.go` don't change
4. Keep all four `init()` registration blocks in the bytebase glue (they're what fires from the blank import in `ultimate.go:34`)
5. Verify: existing `bytebase/backend/plugin/parser/elasticsearch/*_test.go` plus the masking integration test pass unchanged

This is bytebase-side work, tracked under BYT-9000.

---

## Critical Path

```
F0 ──► F1 ──► F6a ──► F6b ──► F7 ──► M1
```

**Six sequential nodes deep.** Everything else parallelizes around it:

- F2 (classification) runs concurrently with F1
- F3, F4, F5 (wrappers) run concurrently with F6a/F6b
- The split between F6a and F6b lets the lower-risk body extractors (sort/source/highlight) land first while the predicate walker gets its own focused PR

Peak concurrent work after F1 lands: **4 nodes** (F3 + F4 + F5 + F6a/b).

---

## Test Data Mapping

Each YAML golden file from `bytebase/backend/plugin/parser/elasticsearch/test-data/` lands with a specific node. The file becomes the green-light criterion for that node:

| Node | YAML golden file(s) | What it asserts |
|------|---------------------|-----------------|
| F1 | `parse-elasticsearch-rest.yaml` | Core parser cases — single/multi request, multi-document body, comments, literal strings, error recovery, byte offsets |
| F2 | `masking_classify_api.yaml`* | URL/method → API type. *Note: this YAML actually targets `MaskableAPI` classification (F6a), not `ClassifyRequest`. F2's `query_type.go` is currently tested by inline `_test.go` cases; lift those alongside F2. |
| F3 | `splitter.yaml` | Splitter outputs (text + position ranges per request) |
| F4 | `statement-ranges.yaml` | LSP range outputs with UTF-16 code-unit handling |
| F5 | (inline `diagnose_test.go`) | Currently a 27-LoC stub; lift inline tests alongside F5 |
| F6a | `masking_classify_api.yaml`, `masking_analyze_body.yaml` | `MaskableAPI` URL classification + body-level analysis (sort/source/fields/highlight) |
| F6b | `masking_predicate_fields.yaml` | Recursive query-clause walk → predicate field paths over all 13+ clause types |
| F7 | `query_span.yaml`, `masking_analyze_request.yaml` | Top-level integration: full `AnalyzeRequest` + query span output |

`masking_analyze_request.yaml` is the integration golden for the entire analysis pipeline; it runs as part of F7's acceptance criteria and re-runs in CI for every later node as a regression guard.

---

## Notes

1. **No `ast/` package.** The existing parser produces `[]*Request` (flat HTTP envelope) plus `map[string]any` for the body. There is no tree to walk. This is a deliberate departure from snowflake/mongo/partiql/cosmosdb conventions and is justified in the analysis doc — Elasticsearch Query DSL is open-ended and plugin-extensible, so a typed body AST would force premature commitment.

2. **No `catalog/`, `completion/`, `semantic/`, `quality/`, or `deparse/` packages.** Bytebase has no consumer for any of these on elasticsearch. Add later if a new feature requires it. (Snowflake's DAG carries the same notice.)

3. **Test corpus is small and pre-existing.** Unlike snowflake (which had to scrape docs for an "official" corpus), elasticsearch has 8 YAML golden files in the existing parser's `test-data/` directory. Lift them verbatim — they ARE the spec under the freeze decision.

4. **The 13+ predicate clause types in F6b are the migration's highest single-feature risk.** Walker logic is recursive over `map[string]any` and silently drops unknown clauses. Brainstorming for F6b should explicitly enumerate every clause type the existing walker recognizes, with golden test coverage for each.

5. **F2 classification is structurally simple but semantically dense.** ~45+ URL pattern strings across three lists (`readOnlyPostEndpoints`, `infoSchemaEndpoints`, `dmlPostEndpoints`). Lift the lists verbatim — do not refactor them during migration. This is the "easy" warm-up node.

6. **F1 internal decomposition.** F1 is the largest node and could in principle be split into "single-request happy path", "multi-request + error recovery", "multi-document body", and "comments + hjson normalization". I chose to keep them as one node because the state machine is tightly coupled — you can't realistically test multi-request without single-request, and the error recovery infrastructure threads through every layer. Brainstorming for F1 will break it into commits naturally; the *PR* is one node.

7. **The bytebase glue rewrite (M1) is broader than mongo's was.** Mongo's glue rewrite touched primarily `masking.go`. Elasticsearch's glue touches **all six** files (`parser.go`, `splitter.go`, `query_span.go`, `query_type.go`, `masking.go`, `statement_ranges.go`, `diagnose.go`) plus the driver, because every one of them currently calls `ParseElasticsearchREST` directly. Per-file work is small; volume is higher. M1 may want to be split into 2-3 PRs by file group when it lands.

8. **No `init()` registration blocks in omni.** Engine ID dispatch via `base.RegisterSplitterFunc` etc. lives entirely in bytebase. omni's elasticsearch package exposes plain functions (`Parse`, `Split`, `StatementRanges`, `Diagnose`, `Classify`, `Analyze`, `QuerySpan`); bytebase's glue files keep their existing `init()` blocks and just call into omni from the registered handlers. This means **no omni-side change is needed when a new bytebase consumer wants to use the parser.**

---

## Status Tracking

Update the **Status** column in the node table as work progresses:

- `not started`
- `in progress`
- `done` (PR #N)

When a node moves to `done`, run its associated YAML golden file(s) to confirm green status before marking it complete.
