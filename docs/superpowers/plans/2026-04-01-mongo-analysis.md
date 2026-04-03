# MongoDB Statement Analysis Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create `mongo/analysis/` package that extracts structural information (operation type, collection, predicate fields, pipeline stages, joins) from parsed MongoDB AST nodes.

**Architecture:** Single package with three source files: `analysis.go` (types + dispatch), `predicate.go` (filter field extraction), `pipeline.go` (aggregate stage analysis). One test file covers all. The package takes `ast.Node` input from `mongo.Parse()` and returns a flat `StatementAnalysis` struct. No external dependencies beyond `mongo/ast`.

**Tech Stack:** Go, stdlib testing, `mongo/ast` package

**Spec:** `docs/superpowers/specs/2026-04-01-mongo-analysis-design.md`

---

### Task 1: Types and Operation Enum

**Files:**
- Create: `mongo/analysis/analysis.go`
- Test: `mongo/analysis/analysis_test.go`

- [ ] **Step 1: Write failing tests for Operation enum methods**

Create `mongo/analysis/analysis_test.go`:

```go
package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/mongo/analysis"
)

func TestOperationIsRead(t *testing.T) {
	tests := []struct {
		op   analysis.Operation
		want bool
	}{
		{analysis.OpFind, true},
		{analysis.OpFindOne, true},
		{analysis.OpAggregate, true},
		{analysis.OpCount, true},
		{analysis.OpDistinct, true},
		{analysis.OpRead, true},
		{analysis.OpWrite, false},
		{analysis.OpAdmin, false},
		{analysis.OpInfo, false},
		{analysis.OpExplain, false},
		{analysis.OpUnknown, false},
	}
	for _, tc := range tests {
		if got := tc.op.IsRead(); got != tc.want {
			t.Errorf("%v.IsRead() = %v, want %v", tc.op, got, tc.want)
		}
	}
}

func TestOperationString(t *testing.T) {
	tests := []struct {
		op   analysis.Operation
		want string
	}{
		{analysis.OpFind, "find"},
		{analysis.OpFindOne, "findOne"},
		{analysis.OpAggregate, "aggregate"},
		{analysis.OpCount, "count"},
		{analysis.OpDistinct, "distinct"},
		{analysis.OpRead, "read"},
		{analysis.OpWrite, "write"},
		{analysis.OpAdmin, "admin"},
		{analysis.OpInfo, "info"},
		{analysis.OpExplain, "explain"},
		{analysis.OpUnknown, "unknown"},
	}
	for _, tc := range tests {
		if got := tc.op.String(); got != tc.want {
			t.Errorf("%v.String() = %q, want %q", tc.op, got, tc.want)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/...`
Expected: compilation error — package does not exist

- [ ] **Step 3: Implement types and Operation enum**

Create `mongo/analysis/analysis.go`:

```go
// Package analysis extracts structural information from parsed MongoDB statements.
package analysis

import "github.com/bytebase/omni/mongo/ast"

// Operation classifies a MongoDB statement by its analytical behavior.
type Operation int

const (
	OpUnknown   Operation = iota
	OpFind                // find
	OpFindOne             // findOne
	OpAggregate           // aggregate
	OpCount               // countDocuments, estimatedDocumentCount, count
	OpDistinct            // distinct
	OpRead                // getIndexes, stats, storageSize, etc.
	OpWrite               // insertOne, updateMany, deleteOne, etc.
	OpAdmin               // createIndex, drop, renameCollection, etc.
	OpInfo                // show dbs, getCollectionNames, serverStatus, etc.
	OpExplain             // .explain() wrapped
)

// IsRead returns true for read operations (find, findOne, aggregate, count, distinct, read).
func (op Operation) IsRead() bool {
	switch op {
	case OpFind, OpFindOne, OpAggregate, OpCount, OpDistinct, OpRead:
		return true
	}
	return false
}

// IsWrite returns true for write operations.
func (op Operation) IsWrite() bool { return op == OpWrite }

// IsAdmin returns true for administrative operations.
func (op Operation) IsAdmin() bool { return op == OpAdmin }

// IsInfo returns true for informational/metadata operations.
func (op Operation) IsInfo() bool { return op == OpInfo }

// String returns a human-readable name for the operation.
func (op Operation) String() string {
	switch op {
	case OpFind:
		return "find"
	case OpFindOne:
		return "findOne"
	case OpAggregate:
		return "aggregate"
	case OpCount:
		return "count"
	case OpDistinct:
		return "distinct"
	case OpRead:
		return "read"
	case OpWrite:
		return "write"
	case OpAdmin:
		return "admin"
	case OpInfo:
		return "info"
	case OpExplain:
		return "explain"
	default:
		return "unknown"
	}
}

// StatementAnalysis is the result of analyzing a parsed MongoDB statement.
type StatementAnalysis struct {
	Operation        Operation
	MethodName       string   // original method name (e.g. "countDocuments" when Operation is OpCount)
	Collection       string   // target collection; empty for db-level/show commands
	PredicateFields  []string // sorted dot-paths from query filter document
	PipelineStages   []string // aggregate only: stage names in order
	ShapePreserving  bool     // true when all pipeline stages preserve document structure
	UnsupportedStage string   // first non-shape-preserving, non-join stage
	Joins            []JoinInfo
}

// JoinInfo records a join extracted from a $lookup or $graphLookup pipeline stage.
type JoinInfo struct {
	Collection string // the "from" field
	AsField    string // the "as" field
}

// Analyze extracts structural information from a single parsed MongoDB AST node.
// Returns nil for unrecognizable nodes.
func Analyze(node ast.Node) *StatementAnalysis {
	// Stub — implemented in subsequent tasks.
	_ = node
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/analysis.go mongo/analysis/analysis_test.go
git commit -m "feat(mongo/analysis): add Operation enum and StatementAnalysis types"
```

---

### Task 2: Collection Statement Analysis (Core Dispatch)

**Files:**
- Modify: `mongo/analysis/analysis.go`
- Modify: `mongo/analysis/analysis_test.go`

- [ ] **Step 1: Write failing tests for collection statement classification**

Append to `mongo/analysis/analysis_test.go`:

```go
import (
	"github.com/bytebase/omni/mongo"
)

// mustAnalyze parses a single statement and runs Analyze on its AST.
func mustAnalyze(t *testing.T, input string) *analysis.StatementAnalysis {
	t.Helper()
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("Parse(%q): %v", input, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse(%q): got %d statements, want 1", input, len(stmts))
	}
	return analysis.Analyze(stmts[0].AST)
}

func TestAnalyzeCollectionFind(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpFind {
		t.Errorf("Operation = %v, want OpFind", sa.Operation)
	}
	if sa.MethodName != "find" {
		t.Errorf("MethodName = %q, want %q", sa.MethodName, "find")
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
}

func TestAnalyzeCollectionFindOne(t *testing.T) {
	sa := mustAnalyze(t, `db.users.findOne()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpFindOne {
		t.Errorf("Operation = %v, want OpFindOne", sa.Operation)
	}
	if sa.MethodName != "findOne" {
		t.Errorf("MethodName = %q, want %q", sa.MethodName, "findOne")
	}
}

func TestAnalyzeCollectionAggregate(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpAggregate {
		t.Errorf("Operation = %v, want OpAggregate", sa.Operation)
	}
}

func TestAnalyzeCollectionWrite(t *testing.T) {
	for _, method := range []string{
		`db.users.insertOne({})`,
		`db.users.insertMany([])`,
		`db.users.updateOne({}, {})`,
		`db.users.updateMany({}, {})`,
		`db.users.deleteOne({})`,
		`db.users.deleteMany({})`,
		`db.users.replaceOne({}, {})`,
	} {
		sa := mustAnalyze(t, method)
		if sa == nil {
			t.Fatalf("Analyze(%q) returned nil", method)
		}
		if sa.Operation != analysis.OpWrite {
			t.Errorf("Analyze(%q).Operation = %v, want OpWrite", method, sa.Operation)
		}
		if sa.Collection != "users" {
			t.Errorf("Analyze(%q).Collection = %q, want %q", method, sa.Collection, "users")
		}
	}
}

func TestAnalyzeCollectionAdmin(t *testing.T) {
	for _, method := range []string{
		`db.users.createIndex({name: 1})`,
		`db.users.dropIndex("name_1")`,
		`db.users.drop()`,
	} {
		sa := mustAnalyze(t, method)
		if sa == nil {
			t.Fatalf("Analyze(%q) returned nil", method)
		}
		if sa.Operation != analysis.OpAdmin {
			t.Errorf("Analyze(%q).Operation = %v, want OpAdmin", method, sa.Operation)
		}
	}
}

func TestAnalyzeCollectionRead(t *testing.T) {
	for _, method := range []string{
		`db.users.getIndexes()`,
		`db.users.stats()`,
		`db.users.storageSize()`,
	} {
		sa := mustAnalyze(t, method)
		if sa == nil {
			t.Fatalf("Analyze(%q) returned nil", method)
		}
		if sa.Operation != analysis.OpRead {
			t.Errorf("Analyze(%q).Operation = %v, want OpRead", method, sa.Operation)
		}
	}
}

func TestAnalyzeCollectionCount(t *testing.T) {
	for _, stmt := range []string{
		`db.users.countDocuments({})`,
		`db.users.estimatedDocumentCount()`,
		`db.users.count()`,
	} {
		sa := mustAnalyze(t, stmt)
		if sa == nil {
			t.Fatalf("Analyze(%q) returned nil", stmt)
		}
		if sa.Operation != analysis.OpCount {
			t.Errorf("Analyze(%q).Operation = %v, want OpCount", stmt, sa.Operation)
		}
	}
}

func TestAnalyzeCollectionDistinct(t *testing.T) {
	sa := mustAnalyze(t, `db.users.distinct("name")`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpDistinct {
		t.Errorf("Operation = %v, want OpDistinct", sa.Operation)
	}
	if sa.MethodName != "distinct" {
		t.Errorf("MethodName = %q, want %q", sa.MethodName, "distinct")
	}
}

func TestAnalyzeCollectionExplain(t *testing.T) {
	sa := mustAnalyze(t, `db.users.explain().find()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpExplain {
		t.Errorf("Operation = %v, want OpExplain", sa.Operation)
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
}

func TestAnalyzeCollectionGetCollection(t *testing.T) {
	sa := mustAnalyze(t, `db.getCollection("users").find()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
	if sa.Operation != analysis.OpFind {
		t.Errorf("Operation = %v, want OpFind", sa.Operation)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -run TestAnalyzeCollection`
Expected: FAIL — `Analyze` returns nil

- [ ] **Step 3: Implement Analyze dispatch and analyzeCollection**

Replace the stub `Analyze` function in `mongo/analysis/analysis.go`:

```go
// Analyze extracts structural information from a single parsed MongoDB AST node.
// Returns nil for unrecognizable nodes.
func Analyze(node ast.Node) *StatementAnalysis {
	switch n := node.(type) {
	case *ast.CollectionStatement:
		return analyzeCollection(n)
	case *ast.DatabaseStatement:
		return analyzeDatabase(n)
	case *ast.ShowCommand:
		return &StatementAnalysis{Operation: OpInfo, MethodName: "show"}
	case *ast.BulkStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: "bulkOp", Collection: n.Collection}
	case *ast.RsStatement:
		return analyzeRs(n)
	case *ast.ShStatement:
		return analyzeSh(n)
	case *ast.EncryptionStatement:
		return analyzeEncryption(n)
	case *ast.PlanCacheStatement:
		return analyzePlanCache(n)
	case *ast.SpStatement:
		return analyzeSp(n)
	case *ast.ConnectionStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Constructor}
	case *ast.NativeFunctionCall:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Name}
	default:
		return nil
	}
}

func analyzeCollection(s *ast.CollectionStatement) *StatementAnalysis {
	a := &StatementAnalysis{
		Collection: s.Collection,
		MethodName: s.Method,
	}

	if s.Explain {
		a.Operation = OpExplain
		return a
	}

	switch s.Method {
	case "find":
		a.Operation = OpFind
		a.PredicateFields = extractPredicateFields(s.Args)
	case "findOne":
		a.Operation = OpFindOne
		a.PredicateFields = extractPredicateFields(s.Args)
	case "aggregate":
		a.Operation = OpAggregate
		analyzePipelineInto(s.Args, a)
	case "countDocuments", "estimatedDocumentCount", "count":
		a.Operation = OpCount
	case "distinct":
		a.Operation = OpDistinct
	case "insertOne", "insertMany", "updateOne", "updateMany",
		"deleteOne", "deleteMany", "replaceOne",
		"findOneAndUpdate", "findOneAndReplace", "findOneAndDelete",
		"bulkWrite", "mapReduce":
		a.Operation = OpWrite
	case "createIndex", "createIndexes", "dropIndex", "dropIndexes",
		"drop", "renameCollection", "reIndex":
		a.Operation = OpAdmin
	case "getIndexes", "stats", "storageSize", "totalIndexSize",
		"totalSize", "dataSize", "validate", "latencyStats",
		"getShardDistribution":
		a.Operation = OpRead
	default:
		a.Operation = OpUnknown
	}

	return a
}

// Stubs for other analyzers — implemented in Task 4.

func analyzeDatabase(s *ast.DatabaseStatement) *StatementAnalysis {
	return &StatementAnalysis{Operation: OpInfo, MethodName: s.Method}
}

func analyzeRs(s *ast.RsStatement) *StatementAnalysis {
	return &StatementAnalysis{Operation: OpWrite, MethodName: s.MethodName}
}

func analyzeSh(s *ast.ShStatement) *StatementAnalysis {
	return &StatementAnalysis{Operation: OpWrite, MethodName: s.MethodName}
}

func analyzeEncryption(s *ast.EncryptionStatement) *StatementAnalysis {
	return &StatementAnalysis{Operation: OpWrite, MethodName: s.Target}
}

func analyzePlanCache(s *ast.PlanCacheStatement) *StatementAnalysis {
	return &StatementAnalysis{Operation: OpWrite, MethodName: "getPlanCache", Collection: s.Collection}
}

func analyzeSp(s *ast.SpStatement) *StatementAnalysis {
	return &StatementAnalysis{Operation: OpWrite, MethodName: s.MethodName}
}
```

Also add stub functions at the bottom of `analysis.go` for `extractPredicateFields` and `analyzePipelineInto` (these will be implemented in Tasks 3 and 4):

```go
// Stubs — implemented in predicate.go and pipeline.go.

func extractPredicateFields(args []ast.Node) []string {
	return nil
}

func analyzePipelineInto(args []ast.Node, a *StatementAnalysis) {
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/analysis.go mongo/analysis/analysis_test.go
git commit -m "feat(mongo/analysis): implement Analyze dispatch and collection method classification"
```

---

### Task 3: Predicate Field Extraction

**Files:**
- Create: `mongo/analysis/predicate.go` (move stubs from analysis.go)
- Modify: `mongo/analysis/analysis.go` (remove stubs)
- Modify: `mongo/analysis/analysis_test.go`

- [ ] **Step 1: Write failing tests for predicate extraction**

Append to `mongo/analysis/analysis_test.go`:

```go
import "slices"

func assertPredicateFields(t *testing.T, sa *analysis.StatementAnalysis, want []string) {
	t.Helper()
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	slices.Sort(want)
	if !slices.Equal(sa.PredicateFields, want) {
		t.Errorf("PredicateFields = %v, want %v", sa.PredicateFields, want)
	}
}

func TestPredicateSimpleField(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({email: "alice@example.com"})`)
	assertPredicateFields(t, sa, []string{"email"})
}

func TestPredicateNestedDocument(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({contact: {phone: "123"}})`)
	assertPredicateFields(t, sa, []string{"contact", "contact.phone"})
}

func TestPredicateDotPathKey(t *testing.T) {
	sa := mustAnalyze(t, `db.users.findOne({"contact.phone": "123"})`)
	assertPredicateFields(t, sa, []string{"contact.phone"})
}

func TestPredicateLogicalOperators(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({
		$or: [
			{email: "a@example.com"},
			{contact: {phone: "123"}}
		],
		$and: [
			{"profile.ssn": "111"},
			{name: "alice"}
		],
		$nor: [{status: "inactive"}]
	})`)
	assertPredicateFields(t, sa, []string{
		"contact", "contact.phone", "email", "name", "profile.ssn", "status",
	})
}

func TestPredicateComparisonOperatorsSkipped(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({age: {$gt: 18, $lt: 65}})`)
	assertPredicateFields(t, sa, []string{"age"})
}

func TestPredicateNoArgs(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.PredicateFields != nil {
		t.Errorf("PredicateFields = %v, want nil", sa.PredicateFields)
	}
}

func TestPredicateFindOneFields(t *testing.T) {
	sa := mustAnalyze(t, `db.users.findOne({name: "alice", age: 30})`)
	assertPredicateFields(t, sa, []string{"age", "name"})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -run TestPredicate`
Expected: FAIL — `PredicateFields` is nil

- [ ] **Step 3: Implement predicate extraction**

Create `mongo/analysis/predicate.go`:

```go
package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mongo/ast"
)

// extractPredicateFields extracts sorted, deduplicated dot-paths from the
// first argument of a find/findOne call (the filter document).
func extractPredicateFields(args []ast.Node) []string {
	if len(args) == 0 {
		return nil
	}
	doc, ok := args[0].(*ast.Document)
	if !ok {
		return nil
	}
	fields := make(map[string]struct{})
	collectFromDocument(doc, "", fields)
	if len(fields) == 0 {
		return nil
	}
	result := make([]string, 0, len(fields))
	for f := range fields {
		result = append(result, f)
	}
	slices.Sort(result)
	return result
}

// collectFromDocument walks a Document's key-value pairs.
func collectFromDocument(doc *ast.Document, prefix string, fields map[string]struct{}) {
	for _, kv := range doc.Pairs {
		if strings.HasPrefix(kv.Key, "$") {
			if isLogicalOperator(kv.Key) {
				collectFromLogicalOp(kv.Value, prefix, fields)
			}
			continue
		}
		fullPath := joinPath(prefix, kv.Key)
		fields[fullPath] = struct{}{}
		collectFromValue(kv.Value, fullPath, fields)
	}
}

// collectFromValue recurses into nested documents or arrays.
func collectFromValue(node ast.Node, prefix string, fields map[string]struct{}) {
	switch v := node.(type) {
	case *ast.Document:
		collectFromDocument(v, prefix, fields)
	case *ast.Array:
		for _, elem := range v.Elements {
			if doc, ok := elem.(*ast.Document); ok {
				collectFromDocument(doc, prefix, fields)
			}
		}
	}
}

// collectFromLogicalOp handles $and/$or/$nor values.
func collectFromLogicalOp(node ast.Node, prefix string, fields map[string]struct{}) {
	switch v := node.(type) {
	case *ast.Document:
		collectFromDocument(v, prefix, fields)
	case *ast.Array:
		for _, elem := range v.Elements {
			if doc, ok := elem.(*ast.Document); ok {
				collectFromDocument(doc, prefix, fields)
			}
		}
	}
}

func isLogicalOperator(key string) bool {
	switch key {
	case "$and", "$or", "$nor":
		return true
	}
	return false
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
```

Remove the `extractPredicateFields` stub from `analysis.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/predicate.go mongo/analysis/analysis.go
git commit -m "feat(mongo/analysis): implement predicate field extraction from filter documents"
```

---

### Task 4: Pipeline Analysis

**Files:**
- Create: `mongo/analysis/pipeline.go` (move stub from analysis.go)
- Modify: `mongo/analysis/analysis.go` (remove stub)
- Modify: `mongo/analysis/analysis_test.go`

- [ ] **Step 1: Write failing tests for pipeline analysis**

Append to `mongo/analysis/analysis_test.go`:

```go
func TestPipelineShapePreserving(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$match: {name: "alice"}}, {$sort: {name: 1}}, {$limit: 10}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
	if sa.UnsupportedStage != "" {
		t.Errorf("UnsupportedStage = %q, want empty", sa.UnsupportedStage)
	}
	wantStages := []string{"$match", "$sort", "$limit"}
	if !slices.Equal(sa.PipelineStages, wantStages) {
		t.Errorf("PipelineStages = %v, want %v", sa.PipelineStages, wantStages)
	}
	assertPredicateFields(t, sa, []string{"name"})
}

func TestPipelineUnsupportedGroup(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$group: {_id: "$status"}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.ShapePreserving {
		t.Error("ShapePreserving = true, want false")
	}
	if sa.UnsupportedStage != "$group" {
		t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, "$group")
	}
}

func TestPipelineUnsupportedProject(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$match: {name: "alice"}}, {$project: {name: 1}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.UnsupportedStage != "$project" {
		t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, "$project")
	}
}

func TestPipelineEmpty(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
	if sa.UnsupportedStage != "" {
		t.Errorf("UnsupportedStage = %q, want empty", sa.UnsupportedStage)
	}
}

func TestPipelineNoArgs(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpAggregate {
		t.Errorf("Operation = %v, want OpAggregate", sa.Operation)
	}
}

func TestPipelineLookupSimple(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "userId", as: "orders"}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if len(sa.Joins) != 1 {
		t.Fatalf("Joins: got %d, want 1", len(sa.Joins))
	}
	if sa.Joins[0].Collection != "orders" {
		t.Errorf("Joins[0].Collection = %q, want %q", sa.Joins[0].Collection, "orders")
	}
	if sa.Joins[0].AsField != "orders" {
		t.Errorf("Joins[0].AsField = %q, want %q", sa.Joins[0].AsField, "orders")
	}
	if sa.UnsupportedStage != "" {
		t.Errorf("UnsupportedStage = %q, want empty", sa.UnsupportedStage)
	}
}

func TestPipelineLookupPipelineForm(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$match: {status: "active"}}], as: "orders"}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.UnsupportedStage != "$lookup" {
		t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, "$lookup")
	}
}

func TestPipelineGraphLookup(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$graphLookup: {from: "employees", startWith: "$reportsTo", connectFromField: "reportsTo", connectToField: "name", as: "hierarchy"}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if len(sa.Joins) != 1 {
		t.Fatalf("Joins: got %d, want 1", len(sa.Joins))
	}
	if sa.Joins[0].Collection != "employees" {
		t.Errorf("Joins[0].Collection = %q, want %q", sa.Joins[0].Collection, "employees")
	}
	if sa.Joins[0].AsField != "hierarchy" {
		t.Errorf("Joins[0].AsField = %q, want %q", sa.Joins[0].AsField, "hierarchy")
	}
}

func TestPipelineUnwind(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$unwind: "$tags"}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
}

func TestPipelineMatchWithUnwind(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$match: {status: "active"}}, {$unwind: "$tags"}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	assertPredicateFields(t, sa, []string{"status"})
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
}

func TestPipelineMatchLogicalOps(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$match: {$or: [{age: {$gt: 18}}, {name: "alice"}]}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	assertPredicateFields(t, sa, []string{"age", "name"})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -run TestPipeline`
Expected: FAIL — pipeline fields not populated

- [ ] **Step 3: Implement pipeline analysis**

Create `mongo/analysis/pipeline.go`:

```go
package analysis

import "github.com/bytebase/omni/mongo/ast"

// shapePreservingStages are pipeline stages whose output retains the
// original document structure.
var shapePreservingStages = map[string]bool{
	"$match":           true,
	"$sort":            true,
	"$limit":           true,
	"$skip":            true,
	"$sample":          true,
	"$addFields":       true,
	"$set":             true,
	"$unset":           true,
	"$geoNear":         true,
	"$setWindowFields": true,
	"$fill":            true,
	"$redact":          true,
	"$unwind":          true,
}

// analyzePipelineInto extracts pipeline stage info into a StatementAnalysis.
func analyzePipelineInto(args []ast.Node, a *StatementAnalysis) {
	if len(args) == 0 {
		a.ShapePreserving = true
		return
	}
	arr, ok := args[0].(*ast.Array)
	if !ok {
		a.UnsupportedStage = "unknown"
		return
	}

	fields := make(map[string]struct{})
	a.ShapePreserving = true

	for _, elem := range arr.Elements {
		stageDoc, ok := elem.(*ast.Document)
		if !ok || len(stageDoc.Pairs) == 0 {
			a.UnsupportedStage = "unknown"
			a.ShapePreserving = false
			return
		}

		stageName := stageDoc.Pairs[0].Key
		a.PipelineStages = append(a.PipelineStages, stageName)

		switch {
		case shapePreservingStages[stageName]:
			if stageName == "$match" {
				if doc, ok := stageDoc.Pairs[0].Value.(*ast.Document); ok {
					collectFromDocument(doc, "", fields)
				}
			}
		case stageName == "$lookup":
			join, unsupported := extractLookup(stageDoc.Pairs[0].Value)
			if unsupported {
				a.UnsupportedStage = "$lookup"
				a.ShapePreserving = false
				return
			}
			if join != nil {
				a.Joins = append(a.Joins, *join)
			}
		case stageName == "$graphLookup":
			join := extractGraphLookup(stageDoc.Pairs[0].Value)
			if join != nil {
				a.Joins = append(a.Joins, *join)
			}
		default:
			a.UnsupportedStage = stageName
			a.ShapePreserving = false
			return
		}
	}

	if len(fields) > 0 {
		a.PredicateFields = sortedKeys(fields)
	}
}

// extractLookup parses a $lookup stage value.
// Returns (nil, true) for pipeline-form $lookup.
func extractLookup(node ast.Node) (*JoinInfo, bool) {
	doc, ok := node.(*ast.Document)
	if !ok {
		return nil, false
	}
	var from, as string
	for _, kv := range doc.Pairs {
		switch kv.Key {
		case "pipeline":
			return nil, true
		case "from":
			if sl, ok := kv.Value.(*ast.StringLiteral); ok {
				from = sl.Value
			}
		case "as":
			if sl, ok := kv.Value.(*ast.StringLiteral); ok {
				as = sl.Value
			}
		}
	}
	if from == "" || as == "" {
		return nil, false
	}
	return &JoinInfo{Collection: from, AsField: as}, false
}

// extractGraphLookup parses a $graphLookup stage value.
func extractGraphLookup(node ast.Node) *JoinInfo {
	doc, ok := node.(*ast.Document)
	if !ok {
		return nil
	}
	var from, as string
	for _, kv := range doc.Pairs {
		switch kv.Key {
		case "from":
			if sl, ok := kv.Value.(*ast.StringLiteral); ok {
				from = sl.Value
			}
		case "as":
			if sl, ok := kv.Value.(*ast.StringLiteral); ok {
				as = sl.Value
			}
		}
	}
	if from == "" || as == "" {
		return nil
	}
	return &JoinInfo{Collection: from, AsField: as}
}

// sortedKeys returns sorted keys from a string set.
func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
```

Add the `slices` import to `pipeline.go`:

```go
import (
	"slices"

	"github.com/bytebase/omni/mongo/ast"
)
```

Remove the `analyzePipelineInto` stub from `analysis.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/pipeline.go mongo/analysis/analysis.go
git commit -m "feat(mongo/analysis): implement aggregate pipeline analysis with join extraction"
```

---

### Task 5: Non-Collection Statement Analyzers

**Files:**
- Modify: `mongo/analysis/analysis.go`
- Modify: `mongo/analysis/analysis_test.go`

- [ ] **Step 1: Write failing tests for non-collection statements**

Append to `mongo/analysis/analysis_test.go`:

```go
func TestAnalyzeShowDbs(t *testing.T) {
	sa := mustAnalyze(t, `show dbs`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeShowCollections(t *testing.T) {
	sa := mustAnalyze(t, `show collections`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeDatabaseDropDatabase(t *testing.T) {
	sa := mustAnalyze(t, `db.dropDatabase()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpAdmin {
		t.Errorf("Operation = %v, want OpAdmin", sa.Operation)
	}
}

func TestAnalyzeDatabaseCreateCollection(t *testing.T) {
	sa := mustAnalyze(t, `db.createCollection("test")`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpAdmin {
		t.Errorf("Operation = %v, want OpAdmin", sa.Operation)
	}
}

func TestAnalyzeDatabaseGetCollectionNames(t *testing.T) {
	sa := mustAnalyze(t, `db.getCollectionNames()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeDatabaseServerStatus(t *testing.T) {
	sa := mustAnalyze(t, `db.serverStatus()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeRsStatus(t *testing.T) {
	sa := mustAnalyze(t, `rs.status()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeRsInitiate(t *testing.T) {
	sa := mustAnalyze(t, `rs.initiate()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("Operation = %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeShStatus(t *testing.T) {
	sa := mustAnalyze(t, `sh.status()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeShAddShard(t *testing.T) {
	sa := mustAnalyze(t, `sh.addShard("shard1/host:port")`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("Operation = %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeBulkStatement(t *testing.T) {
	sa := mustAnalyze(t, `db.users.initializeOrderedBulkOp().insert({name: "alice"}).execute()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("Operation = %v, want OpWrite", sa.Operation)
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
}

func TestAnalyzeNativeFunction(t *testing.T) {
	sa := mustAnalyze(t, `sleep(1000)`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("Operation = %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeNil(t *testing.T) {
	sa := analysis.Analyze(nil)
	if sa != nil {
		t.Errorf("Analyze(nil) = %v, want nil", sa)
	}
}
```

- [ ] **Step 2: Run tests to verify some fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -run "TestAnalyze(Show|Database|Rs|Sh|Bulk|Native|Nil)"`
Expected: FAIL — database/rs/sh stubs return wrong operations

- [ ] **Step 3: Implement full non-collection analyzers**

Replace the stub functions in `mongo/analysis/analysis.go`:

```go
func analyzeDatabase(s *ast.DatabaseStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.Method}
	switch s.Method {
	case "dropDatabase", "createCollection", "createView":
		a.Operation = OpAdmin
	case "getCollectionNames", "getCollectionInfos", "serverStatus",
		"serverBuildInfo", "version", "hostInfo", "getName",
		"listCommands", "stats":
		a.Operation = OpInfo
	case "runCommand", "adminCommand":
		a.Operation = classifyCommandArgs(s.Args)
	case "getSiblingDB", "getMongo":
		a.Operation = OpWrite
	default:
		a.Operation = OpInfo
	}
	return a
}

func classifyCommandArgs(args []ast.Node) Operation {
	if len(args) == 0 {
		return OpWrite
	}
	doc, ok := args[0].(*ast.Document)
	if !ok {
		return OpWrite
	}
	// Check first key as command name.
	if len(doc.Pairs) > 0 {
		cmd := doc.Pairs[0].Key
		switch cmd {
		case "find", "aggregate", "count", "distinct":
			return OpRead
		case "serverStatus", "listCollections", "listIndexes", "listDatabases",
			"collStats", "dbStats", "hostInfo", "buildInfo", "connectionStatus":
			return OpInfo
		case "create", "drop", "createIndexes", "dropIndexes",
			"renameCollection", "collMod":
			return OpAdmin
		}
	}
	return OpWrite
}

func analyzeRs(s *ast.RsStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.MethodName}
	switch s.MethodName {
	case "status", "conf", "config", "printReplicationInfo", "printSecondaryReplicationInfo":
		a.Operation = OpInfo
	default:
		a.Operation = OpWrite
	}
	return a
}

func analyzeSh(s *ast.ShStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.MethodName}
	switch s.MethodName {
	case "status", "getBalancerState", "isBalancerRunning":
		a.Operation = OpInfo
	default:
		a.Operation = OpWrite
	}
	return a
}

func analyzeEncryption(s *ast.EncryptionStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.Target}
	if len(s.ChainedMethods) > 0 {
		last := s.ChainedMethods[len(s.ChainedMethods)-1].Name
		switch last {
		case "getKey", "getKeyByAltName", "getKeys", "decrypt", "encrypt", "encryptExpression":
			a.Operation = OpInfo
			a.MethodName = last
			return a
		}
	}
	a.Operation = OpWrite
	return a
}

func analyzePlanCache(s *ast.PlanCacheStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: "getPlanCache", Collection: s.Collection}
	if len(s.ChainedMethods) > 0 {
		last := s.ChainedMethods[len(s.ChainedMethods)-1].Name
		switch last {
		case "list", "help":
			a.Operation = OpInfo
			a.MethodName = last
			return a
		}
	}
	a.Operation = OpWrite
	return a
}

func analyzeSp(s *ast.SpStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.MethodName}
	if s.SubMethod != "" {
		switch s.SubMethod {
		case "stats", "sample":
			a.Operation = OpInfo
			return a
		}
	} else {
		switch s.MethodName {
		case "listConnections", "listStreamProcessors":
			a.Operation = OpInfo
			return a
		}
	}
	a.Operation = OpWrite
	return a
}
```

Also handle the nil case in `Analyze`:

```go
func Analyze(node ast.Node) *StatementAnalysis {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	// ... rest unchanged
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/analysis.go mongo/analysis/analysis_test.go
git commit -m "feat(mongo/analysis): implement all non-collection statement analyzers"
```

---

### Task 6: runCommand Classification and Edge Cases

**Files:**
- Modify: `mongo/analysis/analysis_test.go`

- [ ] **Step 1: Write tests for runCommand and edge cases**

Append to `mongo/analysis/analysis_test.go`:

```go
func TestAnalyzeRunCommandFind(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({find: "users", filter: {status: "active"}})`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpRead {
		t.Errorf("Operation = %v, want OpRead", sa.Operation)
	}
}

func TestAnalyzeRunCommandCreate(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({create: "newCollection"})`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpAdmin {
		t.Errorf("Operation = %v, want OpAdmin", sa.Operation)
	}
}

func TestAnalyzeRunCommandInsert(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({insert: "users", documents: [{name: "alice"}]})`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("Operation = %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeRunCommandServerStatus(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({serverStatus: 1})`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeAdminCommand(t *testing.T) {
	sa := mustAnalyze(t, `db.adminCommand({listDatabases: 1})`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("Operation = %v, want OpInfo", sa.Operation)
	}
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -run TestAnalyze -v`
Expected: PASS (these should already work from Task 5 implementation)

- [ ] **Step 3: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/analysis_test.go
git commit -m "test(mongo/analysis): add runCommand and adminCommand classification tests"
```

---

### Task 7: Full Integration Test (Masking Parity)

**Files:**
- Modify: `mongo/analysis/analysis_test.go`

This task adds a comprehensive table-driven test that mirrors Bytebase's `masking_test.go` expectations, verifying that the analysis package produces the correct data for every masking scenario.

- [ ] **Step 1: Write integration test**

Append to `mongo/analysis/analysis_test.go`:

```go
// TestAnalyzeMaskingParity verifies that analysis output matches
// the data needed by Bytebase's masking flow for every known scenario.
func TestAnalyzeMaskingParity(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantOp           analysis.Operation
		wantMethod       string
		wantCollection   string
		wantPredicates   []string
		wantShapePres    bool
		wantUnsupported  string
		wantJoins        []analysis.JoinInfo
		wantNil          bool
	}{
		{
			name:           "find with predicate",
			input:          `db.users.find({email: "alice@example.com"})`,
			wantOp:         analysis.OpFind,
			wantMethod:     "find",
			wantCollection: "users",
			wantPredicates: []string{"email"},
		},
		{
			name:           "findOne with dot path key",
			input:          `db.users.findOne({"contact.phone": "123"})`,
			wantOp:         analysis.OpFindOne,
			wantMethod:     "findOne",
			wantCollection: "users",
			wantPredicates: []string{"contact.phone"},
		},
		{
			name:           "getCollection access",
			input:          `db.getCollection("users").find({email: "x"})`,
			wantOp:         analysis.OpFind,
			wantMethod:     "find",
			wantCollection: "users",
			wantPredicates: []string{"email"},
		},
		{
			name:           "aggregate shape-preserving",
			input:          `db.users.aggregate([{$match: {name: "alice"}}])`,
			wantOp:         analysis.OpAggregate,
			wantMethod:     "aggregate",
			wantCollection: "users",
			wantPredicates: []string{"name"},
			wantShapePres:  true,
		},
		{
			name:           "aggregate $group unsupported",
			input:          `db.users.aggregate([{$group: {_id: "$status"}}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantUnsupported: "$group",
		},
		{
			name:           "aggregate $lookup simple",
			input:          `db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "userId", as: "orders"}}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
			wantJoins:      []analysis.JoinInfo{{Collection: "orders", AsField: "orders"}},
		},
		{
			name:            "aggregate $lookup pipeline form",
			input:           `db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$match: {status: "active"}}], as: "orders"}}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$lookup",
		},
		{
			name:           "aggregate $graphLookup",
			input:          `db.users.aggregate([{$graphLookup: {from: "employees", startWith: "$reportsTo", connectFromField: "reportsTo", connectToField: "name", as: "hierarchy"}}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
			wantJoins:      []analysis.JoinInfo{{Collection: "employees", AsField: "hierarchy"}},
		},
		{
			name:           "countDocuments",
			input:          `db.users.countDocuments({})`,
			wantOp:         analysis.OpCount,
			wantMethod:     "countDocuments",
			wantCollection: "users",
		},
		{
			name:           "distinct",
			input:          `db.users.distinct("name")`,
			wantOp:         analysis.OpDistinct,
			wantMethod:     "distinct",
			wantCollection: "users",
		},
		{
			name:    "write method",
			input:   `db.users.insertOne({name: "alice"})`,
			wantOp:  analysis.OpWrite,
			wantCollection: "users",
		},
		{
			name:  "show dbs",
			input: `show dbs`,
			wantOp: analysis.OpInfo,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sa := mustAnalyze(t, tc.input)
			if tc.wantNil {
				if sa != nil {
					t.Fatalf("got %+v, want nil", sa)
				}
				return
			}
			if sa == nil {
				t.Fatal("Analyze returned nil")
			}
			if sa.Operation != tc.wantOp {
				t.Errorf("Operation = %v, want %v", sa.Operation, tc.wantOp)
			}
			if tc.wantMethod != "" && sa.MethodName != tc.wantMethod {
				t.Errorf("MethodName = %q, want %q", sa.MethodName, tc.wantMethod)
			}
			if tc.wantCollection != "" && sa.Collection != tc.wantCollection {
				t.Errorf("Collection = %q, want %q", sa.Collection, tc.wantCollection)
			}
			if tc.wantPredicates != nil {
				slices.Sort(tc.wantPredicates)
				if !slices.Equal(sa.PredicateFields, tc.wantPredicates) {
					t.Errorf("PredicateFields = %v, want %v", sa.PredicateFields, tc.wantPredicates)
				}
			}
			if sa.ShapePreserving != tc.wantShapePres {
				t.Errorf("ShapePreserving = %v, want %v", sa.ShapePreserving, tc.wantShapePres)
			}
			if sa.UnsupportedStage != tc.wantUnsupported {
				t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, tc.wantUnsupported)
			}
			if tc.wantJoins != nil {
				if !slices.Equal(sa.Joins, tc.wantJoins) {
					t.Errorf("Joins = %v, want %v", sa.Joins, tc.wantJoins)
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run the full test suite**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/analysis/... -v`
Expected: PASS

- [ ] **Step 3: Also run existing mongo parser tests to ensure no regressions**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/... -v`
Expected: PASS — all existing tests still pass

- [ ] **Step 4: Commit**

```bash
cd /Users/h3n4l/OpenSource/omni
git add mongo/analysis/analysis_test.go
git commit -m "test(mongo/analysis): add masking parity integration tests"
```
