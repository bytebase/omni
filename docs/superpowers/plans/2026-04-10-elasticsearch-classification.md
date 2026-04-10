# F2: Elasticsearch Request Classification — Implementation Plan

**Goal:** Port `ClassifyRequest` and URL pattern tables from bytebase's `query_type.go` into
`omni/elasticsearch/querytype.go`, plus table-driven tests in `elasticsearch/parsertest/querytype_test.go`.

**Architecture:** Standalone file with no bytebase imports. Omni-native `QueryType` enum defined
alongside `ClassifyRequest`. URL pattern lists ported verbatim.

**Tech Stack:** Go stdlib only (`strings`).

**Branch:** `feat/elasticsearch/classification`

---

### Task 1: Write `elasticsearch/querytype.go`

**Files:**
- Create: `elasticsearch/querytype.go`

- [ ] **Step 1: Create the file with QueryType enum and URL pattern lists**

Port verbatim:
- `QueryType` type and constants (Select, SelectInfoSchema, DML, DDL, Explain, Unknown)
- `readOnlyPostEndpoints` slice (22 strings)
- `infoSchemaEndpoints` slice (3 strings)
- `dmlPostEndpoints` slice (7 strings)
- `ClassifyRequest(method, url string) QueryType`
- `classifyPostRequest(url string) QueryType`
- `isInfoSchemaURL(url string) bool`
- `isReadOnlyPostURL(url string) bool`
- `isDMLPostURL(url string) bool`
- `isDocumentURL(url string) bool`
- `isDocumentWriteURL(url string) bool`

Replace all `base.QueryType` references with the omni-native `QueryType`.
Replace all `base.Select`, `base.DML`, etc. with `QueryTypeSelect`, `QueryTypeDML`, etc.

- [ ] **Step 2: Build check**

Run: `go build ./elasticsearch/...`
Expected: Clean exit.

- [ ] **Step 3: Commit**

```
feat(elasticsearch): add F2 request classification

Port ClassifyRequest and three URL pattern tables (readOnlyPostEndpoints,
infoSchemaEndpoints, dmlPostEndpoints) verbatim from bytebase's
query_type.go. Define omni-native QueryType enum with no bytebase imports.
```

---

### Task 2: Write `elasticsearch/parsertest/querytype_test.go`

**Files:**
- Create: `elasticsearch/parsertest/querytype_test.go`

- [ ] **Step 1: Write table-driven tests covering all decision branches**

Test cases must cover:
- HEAD → Select
- GET non-info-schema → Select
- GET with `_cat/`, `_cluster/`, `_nodes/` → SelectInfoSchema
- DELETE with `_doc/` → DML
- DELETE without `_doc/` → DDL
- PUT with `_doc/`, `_doc`, `_create/`, `_bulk` → DML
- PUT index-level → DDL
- PATCH → DML
- POST with `_explain` → Explain
- POST with info schema URL → SelectInfoSchema
- POST with read-only POST URL → Select (e.g., `_search`, `_count`)
- POST with DML URL → DML (e.g., `_bulk`, `_doc`)
- POST default → DDL
- Unknown method → Unknown
- Case-insensitive method matching (e.g., `get` → Select)

- [ ] **Step 2: Run tests**

Run: `go test ./elasticsearch/parsertest/... -v -run TestClassifyRequest`
Expected: All subtests pass.

- [ ] **Step 3: Run full test suite**

Run: `go test ./elasticsearch/... && go test ./elasticsearch/parsertest/...`
Expected: All pass.

- [ ] **Step 4: Commit**

```
test(elasticsearch): add F2 classification tests

Table-driven tests for ClassifyRequest covering all HTTP methods and
URL-pattern branches. Ported from bytebase's query_type_test.go style.
```

---

### Task 3: Write spec and plan documents

- [ ] Create `docs/superpowers/specs/2026-04-10-elasticsearch-classification-design.md`
- [ ] Create `docs/superpowers/plans/2026-04-10-elasticsearch-classification.md`
- [ ] Commit

---

### Task 4: Final verification

- [ ] `go build ./elasticsearch/...` — clean
- [ ] `go vet ./elasticsearch/...` — clean
- [ ] `go test ./elasticsearch/...` — all pass
- [ ] `go test ./elasticsearch/parsertest/...` — all pass (including smoke test)
- [ ] Push branch and create PR
