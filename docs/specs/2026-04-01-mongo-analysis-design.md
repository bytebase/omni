# MongoDB Statement Analysis Package

**Date:** 2026-04-01
**Scope:** New `mongo/analysis/` package in omni; thin integration layer in Bytebase

## Problem

Bytebase's MongoDB masking and query-type classification currently depend on an ANTLR-based parser (`github.com/bytebase/parser/mongodb`). The goal is to migrate to omni's native MongoDB parser (`omni/mongo`). The masking code walks ANTLR parse trees via listener callbacks to extract structural information (collection name, operation type, predicate fields, pipeline stages, joins). Omni's AST already provides this information in a cleaner form, but there is no analysis layer to surface it.

## Decision

Following omni's convention (PG has `pg/catalog/`, CosmosDB has `cosmosdb/analysis/`), MongoDB analysis lives in omni as `mongo/analysis/`. The package provides generic structural analysis. Bytebase keeps only the policy mapping (which operations are maskable, which are blocked).

## Package: `mongo/analysis/`

### Files

| File | Purpose |
|------|---------|
| `analysis.go` | `Analyze()` entry point, `StatementAnalysis` struct, `Operation` enum, dispatch by AST node type |
| `predicate.go` | Predicate field extraction from query filter documents |
| `pipeline.go` | Aggregate pipeline stage classification, `$lookup`/`$graphLookup` join extraction |
| `analysis_test.go` | Table-driven tests |

### Types

```go
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
    OpRead                // getIndexes, stats, storageSize, totalIndexSize, etc.
    OpWrite               // insertOne, insertMany, updateOne, updateMany, deleteOne, deleteMany, replaceOne, etc.
    OpAdmin               // createIndex, createIndexes, dropIndex, dropIndexes, drop, renameCollection, reIndex
    OpInfo                // show dbs, show collections, getCollectionNames, serverStatus, etc.
    OpExplain             // .explain() wrapped
)

func (op Operation) IsRead() bool
func (op Operation) IsWrite() bool
func (op Operation) IsAdmin() bool
func (op Operation) IsInfo() bool
func (op Operation) String() string

// StatementAnalysis is the result of analyzing a parsed MongoDB statement.
type StatementAnalysis struct {
    Operation        Operation
    MethodName       string     // original method name (e.g. "countDocuments" when Operation is OpCount)
    Collection       string     // target collection; empty for db-level/show commands
    PredicateFields  []string   // sorted dot-paths from query filter document
    PipelineStages   []string   // aggregate only: stage names in order (e.g. "$match", "$sort")
    ShapePreserving  bool       // true when all pipeline stages preserve document structure
    UnsupportedStage string     // first stage that is neither shape-preserving nor a join; empty if all supported
    Joins            []JoinInfo // $lookup/$graphLookup join info
}

// JoinInfo records a join extracted from a $lookup or $graphLookup pipeline stage.
type JoinInfo struct {
    Collection string // the "from" field
    AsField    string // the "as" field
}

// Analyze extracts structural information from a single parsed MongoDB AST node.
// Returns nil for unrecognizable nodes.
func Analyze(node ast.Node) *StatementAnalysis
```

### Analysis Dispatch

`Analyze()` type-switches on the AST node:

| AST Node | Handling |
|----------|----------|
| `*ast.CollectionStatement` | Core path: classify by `Method` string, extract predicates for find/findOne, analyze pipeline for aggregate |
| `*ast.DatabaseStatement` | Classify by `Method`: `dropDatabase`/`createCollection` -> OpAdmin, `getCollectionNames`/`serverStatus`/etc. -> OpInfo, `runCommand`/`adminCommand` -> classify by argument keyword heuristic (see below) |
| `*ast.ShowCommand` | OpInfo |
| `*ast.BulkStatement` | OpWrite |
| `*ast.RsStatement` | Classify by `MethodName`: `status`/`conf`/etc. -> OpInfo, others -> OpWrite |
| `*ast.ShStatement` | Classify by `MethodName`: `status`/`getBalancerState`/etc. -> OpInfo, others -> OpWrite |
| `*ast.EncryptionStatement` | Classify by last chained method name: `getKey`/`getKeys`/`decrypt`/`encrypt` -> OpInfo, others -> OpWrite |
| `*ast.PlanCacheStatement` | `list` -> OpInfo, `clear` -> OpWrite |
| `*ast.SpStatement` | Classify by method/sub-method names |
| `*ast.ConnectionStatement` | OpWrite |
| `*ast.NativeFunctionCall` | OpWrite |

### Collection Method Classification

For `*ast.CollectionStatement`, the `Method` field maps to `Operation`:

| Method(s) | Operation |
|-----------|-----------|
| `find` | OpFind |
| `findOne` | OpFindOne |
| `aggregate` | OpAggregate |
| `countDocuments`, `estimatedDocumentCount`, `count` | OpCount |
| `distinct` | OpDistinct |
| `insertOne`, `insertMany`, `updateOne`, `updateMany`, `deleteOne`, `deleteMany`, `replaceOne`, `findOneAndUpdate`, `findOneAndReplace`, `findOneAndDelete`, `bulkWrite`, `mapReduce` | OpWrite |
| `createIndex`, `createIndexes`, `dropIndex`, `dropIndexes`, `drop`, `renameCollection`, `reIndex` | OpAdmin |
| `getIndexes`, `stats`, `storageSize`, `totalIndexSize`, `totalSize`, `dataSize`, `validate`, `latencyStats`, `getShardDistribution` | OpRead |
| `.explain()` wrapper | OpExplain |

### `runCommand`/`adminCommand` Classification

For `db.runCommand({...})` and `db.adminCommand({...})`, classification uses keyword substring matching on the serialized argument text (matching existing Bytebase behavior):

- Contains `find`, `aggregate`, `count`, `distinct` -> OpRead
- Contains `serverStatus`, `listCollections`, `listIndexes`, `listDatabases`, `collStats`, `dbStats`, `hostInfo`, `buildInfo`, `connectionStatus` -> OpInfo
- Contains `create`, `drop`, `createIndexes`, `dropIndexes`, `renameCollection`, `collMod` -> OpAdmin
- Default -> OpWrite

### Predicate Field Extraction (`predicate.go`)

Applies to OpFind and OpFindOne. Extracts dot-delimited field paths from the first argument (filter document):

- Regular keys become field paths: `{email: "x"}` -> `["email"]`
- Nested documents extend the path: `{contact: {phone: "x"}}` -> `["contact", "contact.phone"]`
- `$and`, `$or`, `$nor` are logical operators: recurse into their array/document values with the current prefix
- Other `$`-prefixed keys (`$gt`, `$in`, etc.) are comparison operators: skip, not field paths
- Arrays in values: recurse into document elements
- Result is sorted and deduplicated

### Pipeline Analysis (`pipeline.go`)

Applies to OpAggregate. Walks the first argument (pipeline array):

**Shape-preserving stages** (document structure is retained):
`$match`, `$sort`, `$limit`, `$skip`, `$sample`, `$addFields`, `$set`, `$unset`, `$geoNear`, `$setWindowFields`, `$fill`, `$redact`, `$unwind`

**Join stages:**
- `$lookup` (simple form): extract `from` and `as` fields -> `JoinInfo`. Pipeline-form `$lookup` (has `pipeline` key) sets `UnsupportedStage`.
- `$graphLookup`: extract `from` and `as` fields -> `JoinInfo`.

**Any other stage** (e.g. `$group`, `$project`, `$facet`): sets `UnsupportedStage`, `ShapePreserving = false`, returns early.

`$match` stages also contribute predicate fields via the same extraction logic as find filters.

## Bytebase Integration

### `masking.go` (rewrite)

Becomes a ~60-line mapper:

1. Parse with `mongo.Parse(statement)`
2. Call `analysis.Analyze(stmts[0].AST)`
3. Map `analysis.Operation` to `base.MongoDBMaskableAPI`:
   - `OpFind` -> `MaskableAPIFind`
   - `OpFindOne` -> `MaskableAPIFindOne`
   - `OpAggregate` without unsupported stage -> `MaskableAPIAggregate`
   - `OpAggregate` with unsupported stage -> `MaskableAPIUnsupportedRead`
   - `OpCount`, `OpDistinct`, `OpRead` -> `MaskableAPIUnsupportedRead`
   - All others -> `nil` (not relevant to masking)
4. Copy `PredicateFields`, `Joins` -> `JoinedCollections`

### `query_type.go` (rewrite)

Becomes a ~30-line mapper:

1. Parse with `mongo.Parse(statement)`
2. Call `analysis.Analyze(stmts[0].AST)`
3. Map `analysis.Operation` to `base.QueryType`:
   - `OpFind`, `OpFindOne`, `OpCount`, `OpDistinct` -> `Select`
   - `OpAggregate` -> `Select` (unless pipeline contains `$out`/`$merge` -> `DML`)
   - `OpWrite` -> `DML`
   - `OpAdmin` -> `DDL`
   - `OpInfo`, `OpRead` -> `SelectInfoSchema`
   - `OpExplain` -> `Explain`

### `query_span.go` (unchanged)

Already parser-agnostic. Calls `AnalyzeMaskingStatement` and `GetQueryType`, builds `PathAST` from predicate fields. No changes needed.

### Deleted code

- All ANTLR imports (`antlr4-go/antlr/v4`, `github.com/bytebase/parser/mongodb`)
- `parseMongoShellRaw()` and ANTLR error listener
- Listener structs (`maskingAnalysisListener`, `queryTypeListener`)
- Value extraction helpers (`extractDocumentValue`, `extractArrayValue`, `extractPairKey`, `extractStringLiteralValue`, `extractCollectionName`, `unquoteMongoString`)
- `base.ANTLRAST` usage in MongoDB path

### Unchanged

- `base.MongoDBAnalysis` struct
- `base.MongoDBMaskableAPI` enum
- `base.MongoDBJoinedCollection` struct
- `query_span.go` (`GetQuerySpan`, `dotPathToPathAST`)
- `document_masking.go` (consumer)
- All downstream test expectations (same inputs/outputs)

## Testing

### In omni (`mongo/analysis/analysis_test.go`)

Table-driven tests covering:
- Operation classification for all AST node types
- Collection access patterns (dot, getCollection)
- Predicate field extraction (nested docs, logical operators, `$`-operators)
- Pipeline analysis (shape-preserving, unsupported, joins)
- `$lookup` simple vs pipeline form
- `$graphLookup` join extraction
- Edge cases (empty args, no pipeline, nil nodes)

### In Bytebase

- `masking_test.go`: thin tests verifying the `Operation` -> `MaskableAPI` mapping
- `query_type_test.go`: thin tests verifying `Operation` -> `QueryType` mapping
- Existing e2e masking tests (`document_masking.go`) unchanged — they validate the full pipeline end-to-end
