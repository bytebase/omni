# MongoDB Auto-Complete Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add auto-complete for MongoDB (mongosh) commands in the web SQL editor, providing context-aware completion candidates based on cursor position and catalog metadata.

**Architecture:** Lexer-based context detection. Tokenize input up to cursor using the existing `mongo/parser` lexer, scan the token sequence right-to-left to classify the completion context (12 context types), then return hardcoded + catalog-sourced candidates filtered by prefix. No parser modifications needed.

**Tech Stack:** Go, existing `mongo/parser` lexer (`NewLexer`/`NextToken`), new `mongo/catalog` and `mongo/completion` packages.

**Design Spec:** `docs/superpowers/specs/2026-04-02-mongo-completion-design.md`

---

## File Structure

| Action | File | Purpose |
|--------|------|---------|
| Modify | `mongo/parser/token.go` | Export minimal token type constants + `IsWord()` method |
| Create | `mongo/catalog/catalog.go` | `Catalog` struct storing collection names |
| Create | `mongo/catalog/catalog_test.go` | Catalog unit tests |
| Create | `mongo/completion/candidates.go` | Hardcoded candidate lists (methods, operators, stages, helpers) |
| Create | `mongo/completion/context.go` | `CompletionContext` type + `detectContext()` from token sequence |
| Create | `mongo/completion/completion.go` | `Complete()` entry point, prefix extraction, filtering |
| Create | `mongo/completion/completion_test.go` | Table-driven tests (context detection, candidates, end-to-end) |

---

### Task 1: Export Token Constants from Parser

**Files:**
- Modify: `mongo/parser/token.go`

This task exposes the minimum token type information needed by the completion package. The lexer (`Lexer`, `Token`, `NewLexer`, `NextToken`) is already exported. Token type constants (`tokEOF`, `tokString`, `tokIdent`) are not.

- [ ] **Step 1: Add exported constants and IsWord method**

Add at the end of `mongo/parser/token.go`, before the `keywords` map:

```go
// Exported token type constants for use by external packages (e.g., completion).
const (
	TokEOF    = tokEOF
	TokString = tokString
	TokIdent  = tokIdent
)

// IsWord returns true if the token is an identifier or keyword (a "word" token).
// This includes collection names, method names, BSON helpers, and JS keywords.
func (t Token) IsWord() bool {
	return t.Type == tokIdent || t.Type >= 700
}
```

- [ ] **Step 2: Verify build**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/...`
Expected: BUILD SUCCESS (no output)

- [ ] **Step 3: Commit**

```bash
git add mongo/parser/token.go
git commit -m "feat(mongo/parser): export token type constants for completion package"
```

---

### Task 2: Create Catalog Package

**Files:**
- Create: `mongo/catalog/catalog.go`
- Create: `mongo/catalog/catalog_test.go`

- [ ] **Step 1: Write the failing test**

Create `mongo/catalog/catalog_test.go`:

```go
package catalog_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mongo/catalog"
)

func TestNewCatalogEmpty(t *testing.T) {
	cat := catalog.New()
	if got := cat.Collections(); len(got) != 0 {
		t.Errorf("new catalog Collections() = %v, want empty", got)
	}
}

func TestAddAndListCollections(t *testing.T) {
	cat := catalog.New()
	cat.AddCollection("users")
	cat.AddCollection("orders")

	got := cat.Collections()
	want := []string{"orders", "users"}
	if !slices.Equal(got, want) {
		t.Errorf("Collections() = %v, want %v", got, want)
	}
}

func TestAddCollectionDedup(t *testing.T) {
	cat := catalog.New()
	cat.AddCollection("users")
	cat.AddCollection("users")
	cat.AddCollection("users")

	got := cat.Collections()
	if len(got) != 1 {
		t.Errorf("Collections() = %v, want 1 entry", got)
	}
}

func TestCollectionsSortOrder(t *testing.T) {
	cat := catalog.New()
	cat.AddCollection("zebra")
	cat.AddCollection("alpha")
	cat.AddCollection("middle")

	got := cat.Collections()
	want := []string{"alpha", "middle", "zebra"}
	if !slices.Equal(got, want) {
		t.Errorf("Collections() = %v, want %v", got, want)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/catalog/`
Expected: FAIL (package does not exist)

- [ ] **Step 3: Write the implementation**

Create `mongo/catalog/catalog.go`:

```go
// Package catalog provides metadata context for MongoDB auto-completion.
package catalog

import "sort"

// Catalog stores collection names for a database context.
type Catalog struct {
	collections map[string]struct{}
}

// New creates an empty Catalog.
func New() *Catalog {
	return &Catalog{collections: make(map[string]struct{})}
}

// AddCollection registers a collection name. Duplicates are ignored.
func (c *Catalog) AddCollection(name string) {
	c.collections[name] = struct{}{}
}

// Collections returns all registered collection names in sorted order.
func (c *Catalog) Collections() []string {
	result := make([]string, 0, len(c.collections))
	for name := range c.collections {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/catalog/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add mongo/catalog/catalog.go mongo/catalog/catalog_test.go
git commit -m "feat(mongo/catalog): add catalog package for completion metadata"
```

---

### Task 3: Create Completion Types and Candidate Lists

**Files:**
- Create: `mongo/completion/completion.go` (types only, stub `Complete`)
- Create: `mongo/completion/candidates.go`

- [ ] **Step 1: Create completion.go with types and stub**

Create `mongo/completion/completion.go`:

```go
// Package completion provides auto-complete for MongoDB shell (mongosh) commands.
package completion

import (
	"strings"

	"github.com/bytebase/omni/mongo/catalog"
	"github.com/bytebase/omni/mongo/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType int

const (
	CandidateKeyword       CandidateType = iota // top-level keywords (db, rs, sh, show, ...)
	CandidateCollection                          // collection name from catalog
	CandidateMethod                              // collection method (find, insertOne, ...)
	CandidateCursorMethod                        // cursor modifier (sort, limit, ...)
	CandidateAggStage                            // aggregation stage ($match, $group, ...)
	CandidateQueryOperator                       // query operator ($gt, $in, ...)
	CandidateBSONHelper                          // BSON constructor (ObjectId, NumberLong, ...)
	CandidateShowTarget                          // show command target (dbs, collections, ...)
	CandidateDbMethod                            // database method (getName, runCommand, ...)
	CandidateRsMethod                            // replica set method (status, conf, ...)
	CandidateShMethod                            // sharding method (addShard, status, ...)
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition/signature
	Comment    string        // optional doc comment
}

// Complete returns completion candidates for the given mongosh input at the cursor offset.
// cat may be nil if no catalog context is available.
func Complete(input string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	if cursorOffset > len(input) {
		cursorOffset = len(input)
	}

	prefix := extractPrefix(input, cursorOffset)
	tokens := tokenize(input, cursorOffset-len(prefix))
	ctx := detectContext(tokens)
	candidates := candidatesForContext(ctx, cat)

	return filterByPrefix(candidates, prefix)
}

// tokenize lexes input up to the given byte offset and returns all tokens.
func tokenize(input string, limit int) []parser.Token {
	if limit > len(input) {
		limit = len(input)
	}
	lex := parser.NewLexer(input[:limit])
	var tokens []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == parser.TokEOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

// extractPrefix returns the partial token the user is typing at cursorOffset.
// Includes $ as a valid prefix character (for $match, $gt, etc.).
func extractPrefix(input string, cursorOffset int) string {
	if cursorOffset > len(input) {
		cursorOffset = len(input)
	}
	i := cursorOffset
	for i > 0 {
		c := input[i-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' || c == '$' {
			i--
		} else {
			break
		}
	}
	return input[i:cursorOffset]
}

// filterByPrefix filters candidates whose Text starts with prefix.
// Matching is case-sensitive (mongosh is case-sensitive).
func filterByPrefix(candidates []Candidate, prefix string) []Candidate {
	if prefix == "" {
		return candidates
	}
	var result []Candidate
	for _, c := range candidates {
		if strings.HasPrefix(c.Text, prefix) {
			result = append(result, c)
		}
	}
	return result
}
```

- [ ] **Step 2: Create candidates.go with all candidate lists**

Create `mongo/completion/candidates.go`:

```go
package completion

import "github.com/bytebase/omni/mongo/catalog"

// candidatesForContext returns the raw candidate list for a given context,
// optionally enriched by the catalog.
func candidatesForContext(ctx completionContext, cat *catalog.Catalog) []Candidate {
	switch ctx {
	case contextTopLevel:
		return topLevelCandidates()
	case contextAfterDbDot:
		return afterDbDotCandidates(cat)
	case contextAfterCollDot:
		return collectionMethodCandidates()
	case contextAfterBracket:
		return bracketCandidates(cat)
	case contextCursorChain:
		return cursorMethodCandidates()
	case contextShowTarget:
		return showTargetCandidates()
	case contextAfterRsDot:
		return rsMethodCandidates()
	case contextAfterShDot:
		return shMethodCandidates()
	case contextAggStage:
		return aggStageCandidates()
	case contextQueryOperator:
		return queryOperatorCandidates()
	case contextInsideArgs:
		return insideArgsCandidates()
	case contextDocumentKey:
		return documentKeyCandidates()
	default:
		return nil
	}
}

func topLevelCandidates() []Candidate {
	keywords := []string{
		"db", "rs", "sh", "sp", "show",
		"sleep", "load", "print", "printjson",
		"quit", "exit", "help", "it", "cls", "version",
	}
	candidates := make([]Candidate, 0, len(keywords)+len(bsonHelpers))
	for _, kw := range keywords {
		candidates = append(candidates, Candidate{Text: kw, Type: CandidateKeyword})
	}
	for _, h := range bsonHelpers {
		candidates = append(candidates, Candidate{Text: h, Type: CandidateBSONHelper})
	}
	return candidates
}

func afterDbDotCandidates(cat *catalog.Catalog) []Candidate {
	var candidates []Candidate
	if cat != nil {
		for _, name := range cat.Collections() {
			candidates = append(candidates, Candidate{Text: name, Type: CandidateCollection})
		}
	}
	for _, m := range dbMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateDbMethod})
	}
	return candidates
}

func collectionMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(collectionMethods))
	for _, m := range collectionMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateMethod})
	}
	return candidates
}

func bracketCandidates(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	var candidates []Candidate
	for _, name := range cat.Collections() {
		candidates = append(candidates, Candidate{Text: name, Type: CandidateCollection})
	}
	return candidates
}

func cursorMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(cursorMethods))
	for _, m := range cursorMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateCursorMethod})
	}
	return candidates
}

func showTargetCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(showTargets))
	for _, t := range showTargets {
		candidates = append(candidates, Candidate{Text: t, Type: CandidateShowTarget})
	}
	return candidates
}

func rsMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(rsMethods))
	for _, m := range rsMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateRsMethod})
	}
	return candidates
}

func shMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(shMethods))
	for _, m := range shMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateShMethod})
	}
	return candidates
}

func aggStageCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(aggStages))
	for _, s := range aggStages {
		candidates = append(candidates, Candidate{Text: s, Type: CandidateAggStage})
	}
	return candidates
}

func queryOperatorCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(queryOperators))
	for _, op := range queryOperators {
		candidates = append(candidates, Candidate{Text: op, Type: CandidateQueryOperator})
	}
	return candidates
}

func insideArgsCandidates() []Candidate {
	literals := []string{"true", "false", "null"}
	candidates := make([]Candidate, 0, len(bsonHelpers)+len(literals))
	for _, h := range bsonHelpers {
		candidates = append(candidates, Candidate{Text: h, Type: CandidateBSONHelper})
	}
	for _, l := range literals {
		candidates = append(candidates, Candidate{Text: l, Type: CandidateKeyword})
	}
	return candidates
}

func documentKeyCandidates() []Candidate {
	// Combines query operators + BSON helpers + literals.
	// Field name completion requires catalog field data (future work).
	candidates := queryOperatorCandidates()
	candidates = append(candidates, insideArgsCandidates()...)
	return candidates
}

// --- Hardcoded candidate lists ---

var bsonHelpers = []string{
	"ObjectId", "NumberLong", "NumberInt", "NumberDecimal",
	"Timestamp", "Date", "ISODate", "UUID",
	"MD5", "HexData", "BinData", "Code",
	"DBRef", "MinKey", "MaxKey", "RegExp", "Symbol",
}

var collectionMethods = []string{
	"find", "findOne", "findOneAndDelete", "findOneAndReplace", "findOneAndUpdate",
	"insertOne", "insertMany",
	"updateOne", "updateMany",
	"deleteOne", "deleteMany",
	"replaceOne", "bulkWrite",
	"aggregate",
	"count", "countDocuments", "estimatedDocumentCount",
	"distinct", "mapReduce", "watch",
	"createIndex", "createIndexes",
	"dropIndex", "dropIndexes", "getIndexes", "reIndex",
	"drop", "renameCollection",
	"stats", "dataSize", "storageSize", "totalSize", "totalIndexSize",
	"validate", "explain",
	"getShardDistribution", "latencyStats",
	"getPlanCache",
	"initializeOrderedBulkOp", "initializeUnorderedBulkOp",
}

var cursorMethods = []string{
	"sort", "limit", "skip",
	"toArray", "forEach", "map",
	"hasNext", "next", "itcount", "size",
	"pretty", "hint", "min", "max",
	"readPref", "comment", "batchSize", "close",
	"collation", "noCursorTimeout", "allowPartialResults",
	"returnKey", "showRecordId", "allowDiskUse",
	"maxTimeMS", "readConcern", "writeConcern",
	"tailable", "oplogReplay", "projection",
}

var dbMethods = []string{
	"getName", "getSiblingDB", "getMongo",
	"getCollectionNames", "getCollectionInfos", "getCollection",
	"createCollection", "createView",
	"dropDatabase",
	"adminCommand", "runCommand",
	"getProfilingStatus", "setProfilingLevel",
	"getLogComponents", "setLogLevel",
	"fsyncLock", "fsyncUnlock",
	"currentOp", "killOp",
	"getUser", "getUsers", "createUser", "updateUser",
	"dropUser", "dropAllUsers",
	"grantRolesToUser", "revokeRolesFromUser",
	"getRole", "getRoles", "createRole", "updateRole",
	"dropRole", "dropAllRoles",
	"grantPrivilegesToRole", "revokePrivilegesFromRole",
	"grantRolesToRole", "revokeRolesFromRole",
	"serverStatus", "isMaster", "hello", "hostInfo",
}

var showTargets = []string{
	"dbs", "databases", "collections", "tables",
	"profile", "users", "roles",
	"log", "logs", "startupWarnings",
}

var rsMethods = []string{
	"status", "conf", "config",
	"initiate", "reconfig",
	"add", "addArb",
	"stepDown", "freeze",
	"slaveOk", "secondaryOk",
	"syncFrom",
	"printReplicationInfo", "printSecondaryReplicationInfo",
}

var shMethods = []string{
	"addShard", "addShardTag", "addShardToZone", "addTagRange",
	"disableAutoSplit", "enableAutoSplit",
	"enableSharding", "disableBalancing", "enableBalancing",
	"getBalancerState", "isBalancerRunning",
	"moveChunk",
	"removeRangeFromZone", "removeShard", "removeShardTag", "removeShardFromZone",
	"setBalancerState", "shardCollection",
	"splitAt", "splitFind",
	"startBalancer", "stopBalancer",
	"updateZoneKeyRange",
	"status",
}

var aggStages = []string{
	"$match", "$group", "$project", "$sort", "$limit", "$skip",
	"$unwind", "$lookup", "$graphLookup",
	"$addFields", "$set", "$unset",
	"$out", "$merge",
	"$bucket", "$bucketAuto", "$facet",
	"$replaceRoot", "$replaceWith",
	"$sample", "$count", "$redact",
	"$geoNear", "$setWindowFields", "$fill", "$densify",
	"$unionWith",
	"$collStats", "$indexStats", "$planCacheStats",
	"$search", "$searchMeta", "$changeStream",
}

var queryOperators = []string{
	// Comparison
	"$eq", "$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin",
	// Logical
	"$and", "$or", "$not", "$nor",
	// Element
	"$exists", "$type",
	// Evaluation
	"$regex", "$expr", "$mod", "$text", "$where", "$jsonSchema",
	// Array
	"$all", "$elemMatch", "$size",
	// Geospatial
	"$geoWithin", "$geoIntersects", "$near", "$nearSphere",
}
```

- [ ] **Step 3: Create context.go stub so the package compiles**

Create `mongo/completion/context.go`:

```go
package completion

import "github.com/bytebase/omni/mongo/parser"

// completionContext identifies the kind of completion expected.
type completionContext int

const (
	contextTopLevel      completionContext = iota // start of input or after semicolon
	contextAfterDbDot                             // db.|
	contextAfterCollDot                           // db.users.|
	contextAfterBracket                           // db[|
	contextInsideArgs                             // db.users.find(|
	contextDocumentKey                            // {| or {age: 1, |
	contextQueryOperator                          // {age: {$|
	contextAggStage                               // [{$|
	contextCursorChain                            // db.users.find().|
	contextShowTarget                             // show |
	contextAfterRsDot                             // rs.|
	contextAfterShDot                             // sh.|
)

// detectContext analyzes the token sequence to determine the completion context.
func detectContext(tokens []parser.Token) completionContext {
	// Stub — implemented in Task 4.
	return contextTopLevel
}
```

- [ ] **Step 4: Verify build**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/completion/`
Expected: BUILD SUCCESS

- [ ] **Step 5: Commit**

```bash
git add mongo/completion/completion.go mongo/completion/candidates.go mongo/completion/context.go
git commit -m "feat(mongo/completion): add types, candidate lists, and Complete stub"
```

---

### Task 4: Implement Context Detection

**Files:**
- Modify: `mongo/completion/context.go`
- Create: `mongo/completion/completion_test.go` (context tests)

- [ ] **Step 1: Write context detection tests**

Create `mongo/completion/completion_test.go`:

```go
package completion

import (
	"testing"

	"github.com/bytebase/omni/mongo/parser"
)

// Helper to tokenize and detect context.
func detectContextFromInput(input string) completionContext {
	prefix := extractPrefix(input, len(input))
	tokens := tokenize(input, len(input)-len(prefix))
	return detectContext(tokens)
}

func TestDetectContextTopLevel(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"  "},
		{"db.users.find();\n"},
	}
	for _, tc := range tests {
		got := detectContextFromInput(tc.input)
		if got != contextTopLevel {
			t.Errorf("detectContext(%q) = %v, want contextTopLevel", tc.input, got)
		}
	}
}

func TestDetectContextAfterDbDot(t *testing.T) {
	tests := []string{
		"db.",
		"db.u",
		"db.get",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextAfterDbDot {
			t.Errorf("detectContext(%q) = %v, want contextAfterDbDot", input, got)
		}
	}
}

func TestDetectContextAfterCollDot(t *testing.T) {
	tests := []string{
		"db.users.",
		"db.users.f",
		`db["users"].`,
		`db["users"].f`,
		`db.getCollection("users").`,
		`db.getCollection("users").f`,
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextAfterCollDot {
			t.Errorf("detectContext(%q) = %v, want contextAfterCollDot", input, got)
		}
	}
}

func TestDetectContextAfterBracket(t *testing.T) {
	tests := []string{
		`db[`,
		`db["`,
		`db["us`,
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextAfterBracket {
			t.Errorf("detectContext(%q) = %v, want contextAfterBracket", input, got)
		}
	}
}

func TestDetectContextCursorChain(t *testing.T) {
	tests := []string{
		"db.users.find().",
		"db.users.find().s",
		"db.users.find({}).sort({a:1}).",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextCursorChain {
			t.Errorf("detectContext(%q) = %v, want contextCursorChain", input, got)
		}
	}
}

func TestDetectContextShowTarget(t *testing.T) {
	tests := []string{
		"show ",
		"show d",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextShowTarget {
			t.Errorf("detectContext(%q) = %v, want contextShowTarget", input, got)
		}
	}
}

func TestDetectContextAfterRsDot(t *testing.T) {
	tests := []string{"rs.", "rs.s"}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextAfterRsDot {
			t.Errorf("detectContext(%q) = %v, want contextAfterRsDot", input, got)
		}
	}
}

func TestDetectContextAfterShDot(t *testing.T) {
	tests := []string{"sh.", "sh.a"}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextAfterShDot {
			t.Errorf("detectContext(%q) = %v, want contextAfterShDot", input, got)
		}
	}
}

func TestDetectContextAggStage(t *testing.T) {
	tests := []string{
		"db.users.aggregate([{$",
		"db.users.aggregate([{$m",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextAggStage {
			t.Errorf("detectContext(%q) = %v, want contextAggStage", input, got)
		}
	}
}

func TestDetectContextQueryOperator(t *testing.T) {
	tests := []string{
		"db.users.find({age: {$",
		"db.users.find({age: {$g",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextQueryOperator {
			t.Errorf("detectContext(%q) = %v, want contextQueryOperator", input, got)
		}
	}
}

func TestDetectContextInsideArgs(t *testing.T) {
	tests := []string{
		"db.users.find(",
		"db.users.insertOne(",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextInsideArgs {
			t.Errorf("detectContext(%q) = %v, want contextInsideArgs", input, got)
		}
	}
}

func TestDetectContextDocumentKey(t *testing.T) {
	tests := []string{
		"db.users.find({",
		"db.users.insertOne({name: 1, ",
	}
	for _, input := range tests {
		got := detectContextFromInput(input)
		if got != contextDocumentKey {
			t.Errorf("detectContext(%q) = %v, want contextDocumentKey", input, got)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/completion/ -run TestDetectContext -v`
Expected: Most tests FAIL (stub returns contextTopLevel for everything)

- [ ] **Step 3: Implement detectContext**

Replace the full contents of `mongo/completion/context.go`:

```go
package completion

import "github.com/bytebase/omni/mongo/parser"

// completionContext identifies the kind of completion expected.
type completionContext int

const (
	contextTopLevel      completionContext = iota // start of input or after semicolon
	contextAfterDbDot                             // db.|
	contextAfterCollDot                           // db.users.|
	contextAfterBracket                           // db[|
	contextInsideArgs                             // db.users.find(|
	contextDocumentKey                            // {| or {age: 1, |
	contextQueryOperator                          // {age: {$|
	contextAggStage                               // [{$|
	contextCursorChain                            // db.users.find().|
	contextShowTarget                             // show |
	contextAfterRsDot                             // rs.|
	contextAfterShDot                             // sh.|
)

// detectContext analyzes the token sequence to determine the completion context.
// It scans right-to-left, tracking brace/bracket depth, and matches trailing patterns.
func detectContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	if n == 0 {
		return contextTopLevel
	}

	last := tokens[n-1]

	// After semicolon → top-level.
	if last.Str == ";" {
		return contextTopLevel
	}

	// show <target>
	if last.Str == "show" {
		return contextShowTarget
	}

	// Check for $-prefixed token at end: could be agg stage or query operator.
	if last.IsWord() && len(last.Str) > 0 && last.Str[0] == '$' {
		ctx := classifyDollarContext(tokens)
		if ctx != contextTopLevel {
			return ctx
		}
	}

	// Trailing dot patterns: db., rs., sh., db.coll., ...().
	if last.Str == "." {
		return classifyDotContext(tokens)
	}

	// Trailing bracket: db[
	if last.Str == "[" && n >= 2 && tokens[n-2].Str == "db" {
		return contextAfterBracket
	}

	// Trailing open brace: { → could be document key, inside args, or agg stage.
	if last.Str == "{" {
		return classifyBraceContext(tokens)
	}

	// Trailing open paren: method( → inside args.
	if last.Str == "(" {
		return contextInsideArgs
	}

	// Trailing comma or colon inside braces → document key.
	if last.Str == "," || last.Str == ":" {
		if insideBrace(tokens) {
			return contextDocumentKey
		}
	}

	// If none of the above matched, the user may be typing a prefix
	// after one of the patterns above. Strip the trailing word token
	// and re-check.
	if last.IsWord() && n >= 2 {
		return detectContext(tokens[:n-1])
	}

	return contextTopLevel
}

// classifyDotContext determines context when the last token is ".".
func classifyDotContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	// Need at least 2 tokens: <something> "."
	if n < 2 {
		return contextTopLevel
	}

	before := tokens[n-2]

	// db.
	if before.Str == "db" {
		return contextAfterDbDot
	}
	// rs.
	if before.Str == "rs" {
		return contextAfterRsDot
	}
	// sh.
	if before.Str == "sh" {
		return contextAfterShDot
	}

	// (...). → cursor chain (method call result)
	if before.Str == ")" {
		// Walk back to find the matching "(" to verify it's a method call.
		openIdx := findMatchingOpen(tokens, n-2, '(', ')')
		if openIdx >= 0 {
			return contextCursorChain
		}
	}

	// db["coll"]. or db.getCollection("coll"). → after coll dot
	if before.Str == "]" && n >= 4 {
		// db [ "string" ]  .
		bracketOpen := findMatchingOpen(tokens, n-2, '[', ']')
		if bracketOpen >= 1 && tokens[bracketOpen-1].Str == "db" {
			return contextAfterCollDot
		}
	}

	// db.ident. → after coll dot (ident is the collection name)
	if before.IsWord() && n >= 4 && tokens[n-3].Str == "." && tokens[n-4].Str == "db" {
		return contextAfterCollDot
	}

	// db.getCollection("coll"). → after coll dot
	// Pattern: db . getCollection ( "string" ) .
	// before is ")" at n-2
	// Already handled by the ")" case above — but we should return contextAfterCollDot
	// if the chain starts with db.getCollection(...).
	// The ")" case returns contextCursorChain. We need to check further.
	// Let's handle this: if the ) case finds db.getCollection pattern, return contextAfterCollDot.
	// Actually, let's check: does classifyDotContext get called when last="." and before=")"?
	// Yes. We returned contextCursorChain above. But db.getCollection("x"). should be contextAfterCollDot.
	// We need to check what's before the matching "(".

	// Re-check: if before is ")" and it's a db.getCollection() or db.ident.method() pattern
	if before.Str == ")" {
		openIdx := findMatchingOpen(tokens, n-2, '(', ')')
		if openIdx >= 0 {
			// Check if this is db.getCollection("x").
			if openIdx >= 3 &&
				tokens[openIdx-1].Str == "getCollection" &&
				tokens[openIdx-2].Str == "." &&
				tokens[openIdx-3].Str == "db" {
				return contextAfterCollDot
			}
			// Check if this is db.coll.method().
			// Pattern: db . coll . method ( ... ) .
			if openIdx >= 5 &&
				tokens[openIdx-1].IsWord() &&
				tokens[openIdx-2].Str == "." &&
				tokens[openIdx-3].IsWord() &&
				tokens[openIdx-4].Str == "." &&
				tokens[openIdx-5].Str == "db" {
				return contextCursorChain
			}
			return contextCursorChain
		}
	}

	return contextTopLevel
}

// classifyDollarContext determines if a trailing $-prefixed token is an agg stage or query operator.
func classifyDollarContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	// Need at least: { $ or [ { $
	if n < 2 {
		return contextTopLevel
	}

	before := tokens[n-2]

	// {$ → could be agg stage or query operator depending on what's before {.
	if before.Str == "{" {
		// Look before the { to see if it follows [ (agg stage) or : (query operator).
		if n >= 3 {
			beforeBrace := tokens[n-3]
			if beforeBrace.Str == "[" {
				return contextAggStage
			}
			if beforeBrace.Str == ":" {
				return contextQueryOperator
			}
			if beforeBrace.Str == "," {
				// Could be inside an array of stages: [{$match:{}}, {$
				// Walk back to find [ before the matching context.
				if isInsideArray(tokens[:n-2]) {
					return contextAggStage
				}
			}
		}
		// Bare {$ with nothing before → agg stage is unlikely, try query operator.
		return contextQueryOperator
	}

	return contextTopLevel
}

// classifyBraceContext determines context when the last token is "{".
func classifyBraceContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	if n < 2 {
		return contextDocumentKey
	}

	before := tokens[n-2]

	// "(" "{" → inside args with document (e.g., find({)
	if before.Str == "(" {
		return contextDocumentKey
	}

	// ":" "{" → nested document (could be query operator context)
	if before.Str == ":" {
		return contextDocumentKey
	}

	// "," "{" → array element or next arg
	if before.Str == "," {
		return contextDocumentKey
	}

	// "[" "{" → array of documents (e.g., aggregate pipeline)
	if before.Str == "[" {
		return contextDocumentKey
	}

	return contextDocumentKey
}

// findMatchingOpen walks tokens backward from pos to find the matching open delimiter.
// open/close are the delimiter pair (e.g., '(', ')').
// Returns the index of the matching open delimiter, or -1 if not found.
func findMatchingOpen(tokens []parser.Token, pos int, open, close byte) int {
	if pos < 0 || pos >= len(tokens) {
		return -1
	}
	depth := 0
	for i := pos; i >= 0; i-- {
		if len(tokens[i].Str) == 1 {
			ch := tokens[i].Str[0]
			if ch == close {
				depth++
			} else if ch == open {
				depth--
				if depth == 0 {
					return i
				}
			}
		}
	}
	return -1
}

// insideBrace returns true if the token stream has an unclosed "{".
func insideBrace(tokens []parser.Token) bool {
	depth := 0
	for _, tok := range tokens {
		if tok.Str == "{" {
			depth++
		} else if tok.Str == "}" {
			depth--
		}
	}
	return depth > 0
}

// isInsideArray returns true if the token stream has an unclosed "[".
func isInsideArray(tokens []parser.Token) bool {
	depth := 0
	for _, tok := range tokens {
		if tok.Str == "[" {
			depth++
		} else if tok.Str == "]" {
			depth--
		}
	}
	return depth > 0
}
```

- [ ] **Step 4: Run context detection tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/completion/ -run TestDetectContext -v`
Expected: ALL PASS

- [ ] **Step 5: Fix any failing tests**

If any tests fail, adjust the `detectContext` implementation to match the expected patterns. The most common issues will be around the order of pattern checks and edge cases in delimiter matching.

- [ ] **Step 6: Commit**

```bash
git add mongo/completion/context.go mongo/completion/completion_test.go
git commit -m "feat(mongo/completion): implement context detection from token sequence"
```

---

### Task 5: End-to-End Complete() Tests

**Files:**
- Modify: `mongo/completion/completion_test.go`

- [ ] **Step 1: Add end-to-end tests to completion_test.go**

Append to `mongo/completion/completion_test.go`:

```go
func TestExtractPrefix(t *testing.T) {
	tests := []struct {
		input  string
		cursor int
		want   string
	}{
		{"db.", 3, ""},
		{"db.us", 5, "us"},
		{"db.users.find().s", 17, "s"},
		{"db.users.find({age: {$g", 23, "$g"},
		{"db.users.aggregate([{$m", 24, "$m"},
		{"", 0, ""},
		{"show d", 6, "d"},
	}
	for _, tc := range tests {
		got := extractPrefix(tc.input, tc.cursor)
		if got != tc.want {
			t.Errorf("extractPrefix(%q, %d) = %q, want %q", tc.input, tc.cursor, got, tc.want)
		}
	}
}

func TestFilterByPrefix(t *testing.T) {
	candidates := []Candidate{
		{Text: "find"},
		{Text: "findOne"},
		{Text: "forEach"},
		{Text: "aggregate"},
	}

	// Case-sensitive: "f" matches "find", "findOne", "forEach"
	got := filterByPrefix(candidates, "f")
	if len(got) != 3 {
		t.Errorf("filterByPrefix(..., 'f') returned %d candidates, want 3", len(got))
	}

	// Case-sensitive: "F" matches nothing
	got = filterByPrefix(candidates, "F")
	if len(got) != 0 {
		t.Errorf("filterByPrefix(..., 'F') returned %d candidates, want 0", len(got))
	}

	// Empty prefix returns all
	got = filterByPrefix(candidates, "")
	if len(got) != 4 {
		t.Errorf("filterByPrefix(..., '') returned %d candidates, want 4", len(got))
	}
}

func TestCompleteEndToEnd(t *testing.T) {
	cat := newTestCatalog("users", "orders")

	tests := []struct {
		name      string
		input     string
		cursor    int
		wantTexts []string // subset of expected candidate texts
		wantType  CandidateType
	}{
		{
			name:      "db. with catalog",
			input:     "db.",
			cursor:    3,
			wantTexts: []string{"users", "orders"},
			wantType:  CandidateCollection,
		},
		{
			name:      "db. includes db methods",
			input:     "db.",
			cursor:    3,
			wantTexts: []string{"getName", "runCommand"},
			wantType:  CandidateDbMethod,
		},
		{
			name:      "db.users.f prefix filter",
			input:     "db.users.f",
			cursor:    10,
			wantTexts: []string{"find", "findOne", "findOneAndDelete", "findOneAndReplace", "findOneAndUpdate"},
			wantType:  CandidateMethod,
		},
		{
			name:      "cursor chain",
			input:     "db.users.find().s",
			cursor:    17,
			wantTexts: []string{"sort", "skip", "size", "showRecordId"},
			wantType:  CandidateCursorMethod,
		},
		{
			name:      "agg stage",
			input:     "db.users.aggregate([{$m",
			cursor:    23,
			wantTexts: []string{"$match", "$merge"},
			wantType:  CandidateAggStage,
		},
		{
			name:      "query operator",
			input:     "db.users.find({age: {$g",
			cursor:    23,
			wantTexts: []string{"$gt", "$gte", "$geoWithin", "$geoIntersects"},
			wantType:  CandidateQueryOperator,
		},
		{
			name:      "bracket access",
			input:     `db["sys`,
			cursor:    7,
			wantTexts: []string{},
			wantType:  CandidateCollection,
		},
		{
			name:      "show target",
			input:     "show d",
			cursor:    6,
			wantTexts: []string{"dbs", "databases"},
			wantType:  CandidateShowTarget,
		},
		{
			name:      "rs methods",
			input:     "rs.",
			cursor:    3,
			wantTexts: []string{"status", "conf", "initiate"},
			wantType:  CandidateRsMethod,
		},
		{
			name:      "sh methods",
			input:     "sh.",
			cursor:    3,
			wantTexts: []string{"addShard", "enableSharding", "status"},
			wantType:  CandidateShMethod,
		},
		{
			name:      "top level",
			input:     "",
			cursor:    0,
			wantTexts: []string{"db", "rs", "sh", "show"},
			wantType:  CandidateKeyword,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Complete(tc.input, tc.cursor, cat)
			for _, want := range tc.wantTexts {
				found := false
				for _, c := range got {
					if c.Text == want {
						found = true
						if c.Type != tc.wantType {
							t.Errorf("Complete(%q, %d): candidate %q type = %v, want %v",
								tc.input, tc.cursor, want, c.Type, tc.wantType)
						}
						break
					}
				}
				if !found {
					t.Errorf("Complete(%q, %d): missing expected candidate %q (got %v)",
						tc.input, tc.cursor, want, candidateTexts(got))
				}
			}
		})
	}
}

func TestCompleteBracketWithCatalog(t *testing.T) {
	cat := newTestCatalog("system.profile", "users", "system.views")
	got := Complete(`db["sys`, 7, cat)

	texts := candidateTexts(got)
	for _, want := range []string{"system.profile", "system.views"} {
		found := false
		for _, text := range texts {
			if text == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf(`Complete("db[\"sys", 7): missing %q in %v`, want, texts)
		}
	}
}

func TestCompleteNilCatalog(t *testing.T) {
	// Should work without panic when catalog is nil.
	got := Complete("db.", 3, nil)
	// Should still have db methods, just no collections.
	hasMethod := false
	for _, c := range got {
		if c.Type == CandidateDbMethod {
			hasMethod = true
			break
		}
	}
	if !hasMethod {
		t.Error("Complete with nil catalog should still return db methods")
	}
	for _, c := range got {
		if c.Type == CandidateCollection {
			t.Error("Complete with nil catalog should not return collections")
		}
	}
}

func TestCompleteGetCollectionDot(t *testing.T) {
	cat := newTestCatalog("users")
	got := Complete(`db.getCollection("users").f`, 27, cat)
	texts := candidateTexts(got)
	for _, want := range []string{"find", "findOne"} {
		found := false
		for _, text := range texts {
			if text == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf(`Complete(db.getCollection("users").f): missing %q in %v`, want, texts)
		}
	}
}

func TestCompleteCursorOvershoot(t *testing.T) {
	// Cursor past end of input should not panic.
	got := Complete("db.", 100, nil)
	if got == nil {
		t.Error("Complete with oversized cursor should not return nil")
	}
}

// --- Helpers ---

func newTestCatalog(names ...string) *catalog.Catalog {
	cat := catalog.New()
	for _, name := range names {
		cat.AddCollection(name)
	}
	return cat
}

func candidateTexts(candidates []Candidate) []string {
	texts := make([]string, len(candidates))
	for i, c := range candidates {
		texts[i] = c.Text
	}
	return texts
}
```

Add the missing import at the top of the test file:

```go
package completion

import (
	"testing"

	"github.com/bytebase/omni/mongo/catalog"
	"github.com/bytebase/omni/mongo/parser"
)
```

Note: the `parser` import may become unused after removing or adjusting earlier test helpers. Keep only imports that are used.

- [ ] **Step 2: Run all tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/completion/ -v`
Expected: ALL PASS

- [ ] **Step 3: Fix any failing tests**

Iterate on context detection or candidate generation if any end-to-end tests fail. Common issues:
- Bracket access prefix extraction (unterminated string before `[` may consume the `"` character)
- The `db.getCollection("x").` pattern requires the full token walk-back
- Comma inside documents for contextDocumentKey

- [ ] **Step 4: Commit**

```bash
git add mongo/completion/completion_test.go
git commit -m "test(mongo/completion): add end-to-end completion tests"
```

---

### Task 6: Run Full Test Suite and Clean Up

**Files:**
- All files from previous tasks

- [ ] **Step 1: Run all mongo package tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/...`
Expected: ALL PASS (parser tests, analysis tests, catalog tests, completion tests)

- [ ] **Step 2: Run vet and build**

Run: `cd /Users/h3n4l/OpenSource/omni && go vet ./mongo/... && go build ./mongo/...`
Expected: No errors or warnings

- [ ] **Step 3: Verify no regressions in existing parser tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/parsertest/`
Expected: ALL PASS

- [ ] **Step 4: Fix any issues found**

If any tests fail or vet reports issues, fix them before proceeding.

- [ ] **Step 5: Commit any fixes**

```bash
git add -A mongo/
git commit -m "fix(mongo/completion): address test and lint issues"
```

Only create this commit if changes were needed. Skip if all passed clean.

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1 | Export token constants | `mongo/parser/token.go` |
| 2 | Create catalog package | `mongo/catalog/catalog.go`, `catalog_test.go` |
| 3 | Create completion types + candidates | `mongo/completion/completion.go`, `candidates.go`, `context.go` (stub) |
| 4 | Implement context detection + tests | `mongo/completion/context.go`, `completion_test.go` |
| 5 | End-to-end Complete() tests | `mongo/completion/completion_test.go` |
| 6 | Full test suite + cleanup | All files |
