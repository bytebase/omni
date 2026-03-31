# MongoDB Engine Design for Omni

## Overview

Add MongoDB (mongosh) support to omni as a new engine under `mongo/`. This provides a hand-written recursive descent parser with full AST and byte-offset position tracking, covering the complete mongosh syntax. The parser is the first milestone; query analysis and catalog come later.

## Prior Art

Three existing repositories inform this design:

- **bytebase/parser/mongodb** — ANTLR4-based mongosh parser (1097-line grammar, 309 test cases). Serves as the grammar reference.
- **bytebase/gomongo** — Translates mongosh syntax to native Go MongoDB driver calls. Its `Operation` struct validates the statement-centric AST approach.
- **bytebase/bytebase** — MongoDB query span, masking analysis, schema sync. Demonstrates downstream analysis needs.

## Design Decisions

1. **Hand-written recursive descent parser** — matches omni's architecture (pg, mysql, mssql, oracle). Zero external dependencies. The ANTLR grammar is the reference specification, not a runtime dependency.
2. **Statement-centric AST** — the parser desugars method chains into structured statement nodes. `db.users.find({}).sort({}).limit(10)` becomes a single `CollectionStatement` with cursor methods as a slice, not a chain of expression nodes. This aligns with gomongo's `Operation` struct, making future migration straightforward.
3. **AST nodes for all values** — documents, arrays, and literals are AST nodes with `Loc` position tracking, not parsed directly into Go values like `bson.D`. Consumers (gomongo, analysis tools) convert to runtime values when needed.
4. **Full ANTLR grammar parity** — all 9 statement families are parsed. The parser accepts all valid mongosh syntax; support-level filtering is the consumer's responsibility.
5. **Package name: `mongo`** — consistent with omni's naming convention (`pg`, `mysql`, `mssql`, `oracle`).

## Package Layout

```
mongo/
├── ast/              # AST node types
│   ├── nodes.go      # All statement and expression node types
│   └── loc.go        # Loc struct (or reuse from common)
├── parser/           # Recursive descent parser + lexer
│   ├── lexer.go      # Character-by-character tokenization
│   ├── tokens.go     # Token type definitions
│   ├── keywords.go   # mongosh keywords and method names
│   ├── parser.go     # Main dispatch (statement family routing)
│   ├── collection.go # db.collection.method() statements
│   ├── database.go   # db.method() statements
│   ├── document.go   # {key: value} document parsing
│   ├── expression.go # Values, arrays, helpers (ObjectId, ISODate...)
│   ├── shell.go      # show commands
│   ├── bulk.go       # Bulk operations
│   ├── connection.go # Mongo(), connect()
│   ├── replication.go# rs.*
│   ├── sharding.go   # sh.*
│   ├── encryption.go # Key vault, client encryption chains
│   ├── plancache.go  # getPlanCache() chains
│   ├── stream.go     # sp.*
│   └── native.go     # Native function calls
├── parsertest/       # Tests organized by family
│   ├── collection_test.go
│   ├── cursor_test.go
│   ├── database_test.go
│   ├── sharding_test.go
│   ├── bulk_test.go
│   ├── connection_test.go
│   ├── replication_test.go
│   ├── encryption_test.go
│   ├── plancache_test.go
│   ├── stream_test.go
│   ├── shell_test.go
│   ├── native_test.go
│   ├── document_test.go
│   ├── expression_test.go
│   ├── error_test.go
│   └── position_test.go
└── parse.go          # Public API entry point
```

## AST Node Types

### Common Interface

Every node carries source location:

```go
type Loc struct {
    Start int // inclusive byte offset
    End   int // exclusive byte offset
}

type Node interface {
    GetLoc() Loc
}
```

### Statement Nodes

| Node | Syntax | Key Fields |
|------|--------|------------|
| `ShowCommand` | `show dbs`, `show collections` | `Target string` |
| `CollectionStatement` | `db.col.find({...}).sort({...})` | `Collection string`, `Method string`, `Args []Node`, `CursorMethods []CursorMethod` |
| `DatabaseStatement` | `db.createCollection("x")`, `db.stats()`, `db.aggregate(...)`, `db.createUser(...)`, etc. (66+ methods) | `Method string`, `Args []Node` |
| `BulkStatement` | `db.col.initializeOrderedBulkOp()...` | `Collection string`, `Ordered bool`, `Operations []BulkOperation` |
| `ConnectionStatement` | `Mongo("uri")` | `Constructor string`, `Args []Node`, `ChainedMethods []MethodCall` |
| `RsStatement` | `rs.status()` | `MethodName string`, `Args []Node` |
| `ShStatement` | `sh.enableSharding("db")` | `MethodName string`, `Args []Node` |
| `EncryptionStatement` | `db.getMongo().getKeyVault().createKey(...)` | `Target string` (KeyVault/ClientEncryption), `MethodName string`, `Args []Node` |
| `PlanCacheStatement` | `db.col.getPlanCache().clear()` | `Collection string`, `MethodName string`, `Args []Node` |
| `SpStatement` | `sp.createProcessor(...)` | `MethodName string`, `SubMethod string`, `Args []Node` |
| `NativeFunctionCall` | `sleep(1000)` | `Name string`, `Args []Node` |

Supporting types:

```go
type CursorMethod struct {
    Method string
    Args   []Node
    Loc    Loc
}

type BulkOperation struct {
    Method string
    Args   []Node
    Loc    Loc
}

type MethodCall struct {
    Method string
    Args   []Node
    Loc    Loc
}
```

### Value/Expression Nodes

| Node | Syntax |
|------|--------|
| `Document` | `{name: "alice", age: 25}` — contains `[]KeyValue` |
| `KeyValue` | `name: "alice"` — `Key string`, `Value Node` |
| `Array` | `[1, 2, 3]` — contains `[]Node` |
| `StringLiteral` | `"hello"` or `'hello'` — `Value string` |
| `NumberLiteral` | `42`, `3.14`, `1e10` — `Value string`, `IsFloat bool` |
| `BoolLiteral` | `true`, `false` — `Value bool` |
| `NullLiteral` | `null` |
| `RegexLiteral` | `/pattern/flags` — `Pattern string`, `Flags string` |
| `HelperCall` | `ObjectId("...")`, `ISODate()` — `Name string`, `Args []Node` |
| `Identifier` | Unquoted names, `$gt` operators — `Name string` |

## Lexer Design

Character-by-character tokenization following `pg/parser/lexer.go` patterns.

Token categories:

- **Keywords**: `show`, `db`, `rs`, `sh`, `sp`, `true`, `false`, `null`
- **Method name keywords**: `find`, `findOne`, `insertOne`, `aggregate`, `sort`, `limit`, etc. Context-sensitive — also valid as identifiers (collection names, field names).
- **Punctuation**: `(`, `)`, `{`, `}`, `[`, `]`, `.`, `,`, `:`, `;`
- **Literals**: Strings (single/double quoted with escape sequences), numbers (int/float/scientific notation), regex (`/pattern/flags`)
- **Dollar identifiers**: `$gt`, `$match`, `$lookup` — the `$` is part of the token
- **Identifiers**: Collection names, field names, generic method names
- **Comments**: `//` line comments, `/* */` block comments — consumed, not emitted
- **Whitespace**: Consumed, not emitted

## Parser Design

### Top-Level Dispatch

The parser reads the first token(s) to route to a statement family parser:

| First tokens | Parser function |
|---|---|
| `show` | `parseShellCommand()` |
| `db` `.` identifier `(` | `parseDatabaseStatement()` |
| `db` `.` identifier `.` | `parseCollectionOrSpecial()` |
| `rs` `.` | `parseRsStatement()` |
| `sh` `.` | `parseShStatement()` |
| `sp` `.` | `parseSpStatement()` |
| `Mongo` or `connect` | `parseConnectionStatement()` |
| identifier `(` | `parseNativeFunctionCall()` |

### Collection/Special Disambiguation

`parseCollectionOrSpecial()` handles the ambiguity after `db.identifier.`:
- If the next method is `initializeOrderedBulkOp`/`initializeUnorderedBulkOp` → `parseBulkStatement()`
- If the next method is `getPlanCache` → `parsePlanCacheStatement()`
- If `db.getMongo()` → peek for `.getKeyVault()`/`.getClientEncryption()` → `parseEncryptionStatement()`
- Otherwise → `parseCollectionStatement()` (the common case)

### Document Parsing Quirks

MongoDB documents differ from JSON:
- Unquoted keys: `{name: "alice"}`
- Dollar-prefixed keys: `{$gt: 25}`
- Trailing commas: `{a: 1, b: 2,}`
- Helper values: `{_id: ObjectId("..."), date: ISODate()}`
- Regex values: `{name: /^alice/i}`

### Method Chain Desugaring

`db.users.find({age: {$gt: 25}}).sort({name: 1}).limit(10)` is parsed as:

1. Identify `db` → collection access → collection name `"users"`
2. Parse primary method: `find` with document argument
3. See `.sort(...)` → append to `CursorMethods` slice
4. See `.limit(...)` → append to `CursorMethods` slice
5. Return single `CollectionStatement` node

### Error Handling

Stop on first error, consistent with all other omni engines. Return `nil` AST and a `ParseError` with position info. No error recovery or accumulation. The `ParseError` includes line, column, and byte position:

```go
type ParseError struct {
    Message  string
    Position int  // byte offset
    Line     int  // 1-based
    Column   int  // 1-based
}
```

## Public API

```go
package mongo

type Statement struct {
    Text      string
    AST       ast.Node
    ByteStart int
    ByteEnd   int
    Start     Position
    End       Position
}

type Position struct {
    Line   int // 1-based
    Column int // 1-based
}

func Parse(input string) ([]Statement, error)
```

Identical shape to `pg.Parse()`.

## Test Strategy

1. **Port all 309 ANTLR test examples** from `parser/mongodb/examples/`. Each test verifies: (a) parses without error, (b) produces expected AST structure, (c) position tracking is correct.
2. **Tests organized by family** in `mongo/parsertest/` — collection (48), cursor (36), database (61), sharding (53), bulk (21), connection (16), replication (14), encryption (14), stream (9), plan cache (5), native (16), shell (2), plus document/expression/error/position tests.
3. **Error case tests** — malformed syntax, missing parens, bad document structure, `new` keyword rejection.
4. **Position tracking tests** — verify `Loc{Start, End}` byte offsets for every node type, especially inside nested documents and method chains.

## Future Work (Not In Scope)

- **Query analysis**: Field-level lineage through aggregation pipelines (query span from Bytebase)
- **Catalog**: In-memory tracking of collections, indexes, JSON schema validators, views
- **gomongo migration**: Replace ANTLR parser + translator with omni's AST as the `Operation` source
- **Completion**: mongosh code completion using the parser
- **Migration/diffing**: Index and validator schema diffing
