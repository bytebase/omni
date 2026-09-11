# F0: Elasticsearch Bootstrap — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create the `elasticsearch/` package skeleton with public types, YAML golden test harness, and hjson dependency — unblocking F1 and F2.

**Architecture:** Minimal scaffold following mongo's pattern. Public types in top-level `parse.go`, centralized tests in `parsertest/`, YAML goldens lifted verbatim from bytebase. No implementation code — just types and test infrastructure.

**Tech Stack:** Go 1.25, gopkg.in/yaml.v3 (already indirect dep), github.com/hjson/hjson-go/v4 (new direct dep)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap` (branch `feat/elasticsearch/bootstrap`)

---

### File Structure

| Action | Path | Responsibility |
|--------|------|---------------|
| Create | `elasticsearch/parse.go` | Public types: `ParseResult`, `Request`, `SyntaxError`, `Position` |
| Create | `elasticsearch/parser/doc.go` | Package doc placeholder for F1 |
| Create | `elasticsearch/parsertest/helpers_test.go` | Generic `loadYAML[T]` YAML loader + `writeYAML` for record mode |
| Create | `elasticsearch/parsertest/smoke_test.go` | Verify all 8 YAML golden files load without error |
| Create | `elasticsearch/parsertest/testdata/` (8 files) | Verbatim copy of bytebase golden files |
| Modify | `go.mod`, `go.sum` | Add `github.com/hjson/hjson-go/v4` direct dependency |

---

### Task 1: Package skeleton and types

**Files:**
- Create: `elasticsearch/parse.go`
- Create: `elasticsearch/parser/doc.go`

- [ ] **Step 1: Create `elasticsearch/parse.go` with public types**

```go
// Package elasticsearch provides a parser for Kibana Dev Console-style
// Elasticsearch REST request blocks.
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
	StartOffset int      `yaml:"startoffset"`
	EndOffset   int      `yaml:"endoffset"`
}

// SyntaxError represents a parse error at a specific source position.
type SyntaxError struct {
	Position   Position `yaml:"position"`
	Message    string   `yaml:"message"`
	RawMessage string   `yaml:"rawmessage"`
}

// Position represents a location in source text.
type Position struct {
	Line   int `yaml:"line"`
	Column int `yaml:"column"`
}
```

- [ ] **Step 2: Create `elasticsearch/parser/doc.go` placeholder**

```go
// Package parser implements a hand-written recursive-descent parser for
// Kibana Dev Console-style Elasticsearch REST request blocks.
package parser
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap && go build ./elasticsearch/...`
Expected: Clean exit, no errors.

- [ ] **Step 4: Commit skeleton**

```bash
git add elasticsearch/parse.go elasticsearch/parser/doc.go
git commit -m "feat(elasticsearch): add package skeleton with public types

ParseResult, Request, SyntaxError, Position types with YAML tags
matching bytebase golden file format. Empty parser/ package placeholder.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: Add hjson dependency

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: Add hjson-go/v4 to go.mod**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap && go get github.com/hjson/hjson-go/v4`
Expected: `go.mod` updated with new require line, `go.sum` updated.

- [ ] **Step 2: Tidy**

Run: `go mod tidy`
Expected: Clean exit. hjson may drop back to indirect if nothing imports it yet — that's fine, F1 will promote it to direct when it actually imports it.

- [ ] **Step 3: Verify no breakage**

Run: `go build ./...`
Expected: Clean exit, no errors.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "build: add hjson-go/v4 dependency for elasticsearch parser

Required by F1 (parser core) for HJSON tolerance and comment stripping.
Added at bootstrap so the dependency is available when F1 starts.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 3: Copy YAML golden files from bytebase

**Files:**
- Create: `elasticsearch/parsertest/testdata/` (8 YAML files)

Source: `/Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/`

- [ ] **Step 1: Create testdata directory and copy all 8 files**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap
mkdir -p elasticsearch/parsertest/testdata
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/parse-elasticsearch-rest.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/splitter.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/statement-ranges.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/masking_classify_api.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/masking_analyze_body.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/masking_predicate_fields.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/query_span.yaml elasticsearch/parsertest/testdata/
cp /Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/masking_analyze_request.yaml elasticsearch/parsertest/testdata/
```

- [ ] **Step 2: Verify byte-identical copies**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap
for f in parse-elasticsearch-rest.yaml splitter.yaml statement-ranges.yaml masking_classify_api.yaml masking_analyze_body.yaml masking_predicate_fields.yaml query_span.yaml masking_analyze_request.yaml; do
  diff "/Users/h3n4l/OpenSource/bytebase/backend/plugin/parser/elasticsearch/test-data/$f" "elasticsearch/parsertest/testdata/$f" || echo "DIFF: $f"
done
```

Expected: No output (all files identical).

- [ ] **Step 3: Verify count**

```bash
ls elasticsearch/parsertest/testdata/*.yaml | wc -l
```

Expected: `8`

---

### Task 4: Test harness and smoke test

**Files:**
- Create: `elasticsearch/parsertest/helpers_test.go`
- Create: `elasticsearch/parsertest/smoke_test.go`

- [ ] **Step 1: Create `helpers_test.go` with YAML loader**

```go
package parsertest

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// loadYAML loads a YAML file from testdata/ and unmarshals into a slice of T.
// Each feature node's test file defines its own case struct with YAML tags
// matching that golden file's schema.
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
	if len(cases) == 0 {
		t.Fatalf("%s: expected at least one test case, got 0", filename)
	}
	return cases
}

// writeYAML writes cases back to a YAML file in testdata/ (record mode).
// Feature test files use this to regenerate goldens during development.
func writeYAML[T any](t *testing.T, filename string, cases []T) {
	t.Helper()
	data, err := yaml.Marshal(cases)
	if err != nil {
		t.Fatalf("marshal %s: %v", filename, err)
	}
	if err := os.WriteFile("testdata/"+filename, data, 0644); err != nil {
		t.Fatalf("write %s: %v", filename, err)
	}
}
```

- [ ] **Step 2: Create `smoke_test.go`**

```go
package parsertest

import "testing"

// TestGoldenFilesLoadable verifies that all 8 YAML golden files from bytebase
// can be read and unmarshaled without error. This is schema-agnostic — it
// deserializes into map[string]any to catch file corruption or YAML syntax
// issues without depending on any parser types.
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

- [ ] **Step 3: Run smoke tests**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap && go test ./elasticsearch/parsertest/... -v`
Expected: 8 subtests pass (one per YAML file).

```
=== RUN   TestGoldenFilesLoadable
=== RUN   TestGoldenFilesLoadable/parse-elasticsearch-rest.yaml
=== RUN   TestGoldenFilesLoadable/splitter.yaml
=== RUN   TestGoldenFilesLoadable/statement-ranges.yaml
=== RUN   TestGoldenFilesLoadable/masking_classify_api.yaml
=== RUN   TestGoldenFilesLoadable/masking_analyze_body.yaml
=== RUN   TestGoldenFilesLoadable/masking_predicate_fields.yaml
=== RUN   TestGoldenFilesLoadable/query_span.yaml
=== RUN   TestGoldenFilesLoadable/masking_analyze_request.yaml
--- PASS: TestGoldenFilesLoadable
PASS
```

- [ ] **Step 4: Commit golden files + test harness**

```bash
git add elasticsearch/parsertest/
git commit -m "test(elasticsearch): add YAML golden test harness and smoke test

Lift 8 YAML golden files verbatim from bytebase elasticsearch parser
test-data/. Add generic loadYAML[T] / writeYAML[T] helpers for the
table-driven test pattern that all subsequent DAG nodes will use.
Smoke test verifies all goldens load without YAML errors.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"
```

---

### Task 5: Final verification

- [ ] **Step 1: Full build**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-elasticsearch-bootstrap && go build ./elasticsearch/...`
Expected: Clean exit.

- [ ] **Step 2: go vet**

Run: `go vet ./elasticsearch/...`
Expected: Clean exit.

- [ ] **Step 3: Elasticsearch tests pass**

Run: `go test ./elasticsearch/...`
Expected: PASS.

- [ ] **Step 4: Existing mongo tests still pass**

Run: `go test ./mongo/parsertest/...`
Expected: PASS (no regressions).

- [ ] **Step 5: Check git log**

Run: `git log --oneline -5`
Expected: 3 clean commits on `feat/elasticsearch/bootstrap` branch:
1. `feat(elasticsearch): add package skeleton with public types`
2. `build: add hjson-go/v4 dependency for elasticsearch parser`
3. `test(elasticsearch): add YAML golden test harness and smoke test`
