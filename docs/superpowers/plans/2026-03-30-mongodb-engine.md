# MongoDB Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a hand-written recursive descent parser for MongoDB shell (mongosh) syntax with full AST and byte-offset position tracking, matching the ANTLR grammar at `github.com/bytebase/parser/mongodb`.

**Architecture:** A new `mongo/` package following omni's engine pattern (`pg/`, `mysql/`, `mssql/`, `oracle/`). The parser reads mongosh input character-by-character (lexer) then builds a typed AST via recursive descent. Every node carries `Loc{Start, End}` byte offsets. The parser stops on first error. Test cases are ported from the 309 ANTLR example files.

**Tech Stack:** Pure Go, zero external dependencies. Testing with `testing` package and `testify` (already in go.mod).

**Spec:** `docs/superpowers/specs/2026-03-30-mongodb-engine-design.md`

**ANTLR grammar reference:** The grammar files at `/Users/h3n4l/OpenSource/parser/mongodb/MongoShellParser.g4` and `MongoShellLexer.g4` are the authoritative reference for what syntax to support. The 309 test files at `/Users/h3n4l/OpenSource/parser/mongodb/examples/` provide concrete test cases.

---

## File Structure

```
mongo/
├── parse.go                    # Public API: Parse(input) ([]Statement, error)
├── ast/
│   └── nodes.go                # All AST node types (Node interface, Loc, statements, values)
├── parser/
│   ├── lexer.go                # Tokenizer: character-by-character, produces Token stream
│   ├── token.go                # Token struct, token type constants, keyword map
│   ├── parser.go               # Parser struct, advance/peek/expect, top-level dispatch, ParseError
│   ├── document.go             # Document and key-value parsing: {key: value, ...}
│   ├── expression.go           # Value parsing: literals, arrays, helper functions, regex
│   ├── collection.go           # db.collection.method() + cursor method chains
│   ├── database.go             # db.method() (66+ database methods)
│   ├── shell.go                # show dbs, show collections
│   ├── bulk.go                 # Bulk operations: initializeOrderedBulkOp/UnorderedBulkOp chains
│   ├── connection.go           # Mongo(), connect(), db.getMongo() chains
│   ├── replication.go          # rs.method()
│   ├── sharding.go             # sh.method()
│   ├── encryption.go           # db.getMongo().getKeyVault/getClientEncryption chains
│   ├── plancache.go            # db.collection.getPlanCache() chains
│   ├── stream.go               # sp.method(), sp.x.method()
│   └── native.go               # Top-level function calls: sleep(), load(), etc.
└── parsertest/
    ├── helpers_test.go         # Shared test helpers: parseOne(), mustParse(), assertLoc()
    ├── shell_test.go           # show commands (2 ANTLR examples)
    ├── document_test.go        # Document syntax (document_syntax.js)
    ├── expression_test.go      # Literals, helpers, regex (literals.js, helper_functions.js, regex.js)
    ├── collection_test.go      # Collection methods (48 examples)
    ├── cursor_test.go          # Cursor methods (36 examples)
    ├── database_test.go        # Database methods (61 examples)
    ├── bulk_test.go            # Bulk operations (21 examples)
    ├── connection_test.go      # Connection methods (16 examples)
    ├── replication_test.go     # rs.* methods (14 examples)
    ├── sharding_test.go        # sh.* methods (53 examples)
    ├── encryption_test.go      # Encryption methods (14 examples)
    ├── plancache_test.go       # Plan cache methods (5 examples)
    ├── stream_test.go          # Stream processing (9 examples)
    ├── native_test.go          # Native functions (16 examples)
    ├── comment_test.go         # Comments (comments.js)
    ├── unicode_test.go         # Unicode support (unicode.js)
    ├── error_test.go           # Error cases and positions
    └── position_test.go        # Byte offset verification for every node type
```

---

## Task 1: AST Node Types

**Files:**
- Create: `mongo/ast/nodes.go`

This is the foundation — all other tasks depend on it.

- [ ] **Step 1: Create the AST package with Node interface and Loc**

```go
// mongo/ast/nodes.go
package ast

// Loc represents a source location range (byte offsets).
type Loc struct {
	Start int // inclusive start byte offset
	End   int // exclusive end byte offset
}

// Node is the interface implemented by all MongoDB parse tree nodes.
type Node interface {
	GetLoc() Loc
}
```

- [ ] **Step 2: Add statement node types**

Add these types to `mongo/ast/nodes.go`:

```go
// ShowCommand represents "show dbs", "show databases", or "show collections".
type ShowCommand struct {
	Target string // "databases" or "collections"
	Loc    Loc
}

func (n *ShowCommand) GetLoc() Loc { return n.Loc }

// CursorMethod represents a chained cursor method call like .sort({name: 1}).
type CursorMethod struct {
	Method string // e.g. "sort", "limit", "skip"
	Args   []Node
	Loc    Loc
}

// MethodCall represents a generic chained method call.
type MethodCall struct {
	Method string
	Args   []Node
	Loc    Loc
}

// BulkOperation represents a single operation in a bulk chain.
type BulkOperation struct {
	Method string
	Args   []Node
	Loc    Loc
}

// CollectionStatement represents db.collection.method(...).cursorMethod(...).
type CollectionStatement struct {
	Collection    string         // collection name
	CollectionLoc Loc            // location of the collection name
	AccessMethod  string         // "dot", "bracket", or "getCollection"
	Method        string         // e.g. "find", "insertOne", "aggregate"
	Args          []Node         // arguments to the primary method
	CursorMethods []CursorMethod // chained cursor method calls
	Explain       bool           // true if prefixed with .explain()
	ExplainArgs   []Node         // arguments to explain() if present
	Loc           Loc
}

func (n *CollectionStatement) GetLoc() Loc { return n.Loc }

// DatabaseStatement represents db.method(...) calls.
type DatabaseStatement struct {
	Method string // e.g. "createCollection", "stats", "dropDatabase"
	Args   []Node
	Loc    Loc
}

func (n *DatabaseStatement) GetLoc() Loc { return n.Loc }

// BulkStatement represents db.collection.initializeOrderedBulkOp()... chains.
type BulkStatement struct {
	Collection   string          // collection name
	AccessMethod string          // "dot", "bracket", or "getCollection"
	Ordered      bool            // true for initializeOrderedBulkOp
	Operations   []BulkOperation // chained operations
	Loc          Loc
}

func (n *BulkStatement) GetLoc() Loc { return n.Loc }

// ConnectionStatement represents Mongo(...), connect(...), or db.getMongo() chains.
type ConnectionStatement struct {
	Constructor    string       // "Mongo", "connect", or "getMongo"
	Args           []Node       // constructor arguments
	ChainedMethods []MethodCall // chained method calls
	Loc            Loc
}

func (n *ConnectionStatement) GetLoc() Loc { return n.Loc }

// RsStatement represents rs.method(...).
type RsStatement struct {
	MethodName string
	Args       []Node
	Loc        Loc
}

func (n *RsStatement) GetLoc() Loc { return n.Loc }

// ShStatement represents sh.method(...).
type ShStatement struct {
	MethodName string
	Args       []Node
	Loc        Loc
}

func (n *ShStatement) GetLoc() Loc { return n.Loc }

// EncryptionStatement represents db.getMongo().getKeyVault().method(...)
// or db.getMongo().getClientEncryption().method(...) chains.
type EncryptionStatement struct {
	Target         string       // "keyVault" or "clientEncryption"
	ChainedMethods []MethodCall // chained method calls after getKeyVault/getClientEncryption
	Loc            Loc
}

func (n *EncryptionStatement) GetLoc() Loc { return n.Loc }

// PlanCacheStatement represents db.collection.getPlanCache().method(...) chains.
type PlanCacheStatement struct {
	Collection     string       // collection name
	AccessMethod   string       // "dot", "bracket", or "getCollection"
	ChainedMethods []MethodCall // chained method calls after getPlanCache()
	Loc            Loc
}

func (n *PlanCacheStatement) GetLoc() Loc { return n.Loc }

// SpStatement represents sp.method(...) or sp.resource.method(...).
type SpStatement struct {
	MethodName string // first identifier after "sp."
	SubMethod  string // optional second identifier (sp.x.method)
	Args       []Node
	Loc        Loc
}

func (n *SpStatement) GetLoc() Loc { return n.Loc }

// NativeFunctionCall represents top-level function calls like sleep(1000).
type NativeFunctionCall struct {
	Name string
	Args []Node
	Loc  Loc
}

func (n *NativeFunctionCall) GetLoc() Loc { return n.Loc }
```

- [ ] **Step 3: Add value/expression node types**

Add these types to `mongo/ast/nodes.go`:

```go
// Document represents a BSON document: {key: value, ...}.
type Document struct {
	Pairs []KeyValue
	Loc   Loc
}

func (n *Document) GetLoc() Loc { return n.Loc }

// KeyValue represents a single key-value pair in a document.
type KeyValue struct {
	Key      string // the key (unquoted or quoted, without quotes)
	KeyLoc   Loc    // location of the key
	Value    Node   // the value expression
	Loc      Loc    // spans key: value
}

// Array represents a BSON array: [value, ...].
type Array struct {
	Elements []Node
	Loc      Loc
}

func (n *Array) GetLoc() Loc { return n.Loc }

// StringLiteral represents a quoted string: "hello" or 'hello'.
type StringLiteral struct {
	Value string // unescaped string content
	Loc   Loc
}

func (n *StringLiteral) GetLoc() Loc { return n.Loc }

// NumberLiteral represents a numeric literal: 42, 3.14, 1e10.
type NumberLiteral struct {
	Value   string // raw text representation (preserves precision)
	IsFloat bool   // true if contains '.' or 'e'/'E'
	Loc     Loc
}

func (n *NumberLiteral) GetLoc() Loc { return n.Loc }

// BoolLiteral represents true or false.
type BoolLiteral struct {
	Value bool
	Loc   Loc
}

func (n *BoolLiteral) GetLoc() Loc { return n.Loc }

// NullLiteral represents null.
type NullLiteral struct {
	Loc Loc
}

func (n *NullLiteral) GetLoc() Loc { return n.Loc }

// RegexLiteral represents /pattern/flags.
type RegexLiteral struct {
	Pattern string
	Flags   string
	Loc     Loc
}

func (n *RegexLiteral) GetLoc() Loc { return n.Loc }

// HelperCall represents a BSON helper function call: ObjectId("..."), ISODate(), etc.
type HelperCall struct {
	Name string // e.g. "ObjectId", "ISODate", "UUID", "NumberLong"
	Args []Node
	Loc  Loc
}

func (n *HelperCall) GetLoc() Loc { return n.Loc }

// Identifier represents an unquoted name used as a value (rare, mainly in special contexts).
type Identifier struct {
	Name string
	Loc  Loc
}

func (n *Identifier) GetLoc() Loc { return n.Loc }
```

- [ ] **Step 4: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/ast/`
Expected: SUCCESS, no errors

- [ ] **Step 5: Commit**

```bash
git add mongo/ast/nodes.go
git commit -m "feat(mongo/ast): add MongoDB AST node types

Defines the Node interface, Loc position tracking, statement nodes
(CollectionStatement, DatabaseStatement, BulkStatement, etc.), and
value nodes (Document, Array, StringLiteral, HelperCall, etc.)."
```

---

## Task 2: Token Types and Lexer

**Files:**
- Create: `mongo/parser/token.go`
- Create: `mongo/parser/lexer.go`

- [ ] **Step 1: Define token types and keyword map in token.go**

```go
// mongo/parser/token.go
package parser

// Token represents a lexical token.
type Token struct {
	Type int    // token type constant
	Str  string // token text (for identifiers, strings, numbers)
	Loc  int    // byte offset of token start in input
	End  int    // byte offset of token end (exclusive)
}

// Token type constants.
const (
	tokEOF = 0

	// Literal tokens (600+)
	tokNumber = 600 + iota // numeric literal
	tokString              // string literal (single or double quoted)
	tokRegex               // regex literal /pattern/flags
	tokIdent               // identifier (including $-prefixed)

	// Punctuation
	tokLParen   // (
	tokRParen   // )
	tokLBrace   // {
	tokRBrace   // }
	tokLBracket // [
	tokRBracket // ]
	tokColon    // :
	tokComma    // ,
	tokDot      // .
	tokSemi     // ;
)

// Keyword token constants (700+).
// These are recognized during lexing. All keywords also function as identifiers
// (context-sensitive), so the parser can treat them as identifiers where needed.
const (
	kwShow = 700 + iota
	kwDbs
	kwDatabases
	kwCollections
	kwDb
	kwNew
	kwTrue
	kwFalse
	kwNull

	// Collection methods
	kwFind
	kwFindOne
	kwCountDocuments
	kwEstimatedDocumentCount
	kwDistinct
	kwAggregate
	kwGetIndexes
	kwInsertOne
	kwInsertMany
	kwUpdateOne
	kwUpdateMany
	kwDeleteOne
	kwDeleteMany
	kwReplaceOne
	kwFindOneAndUpdate
	kwFindOneAndReplace
	kwFindOneAndDelete
	kwCreateIndex
	kwCreateIndexes
	kwDropIndex
	kwDropIndexes
	kwDrop
	kwRenameCollection
	kwStats
	kwStorageSize
	kwTotalIndexSize
	kwTotalSize
	kwDataSize
	kwIsCapped
	kwValidate
	kwLatencyStats
	kwWatch
	kwBulkWrite
	kwCount
	kwInsert
	kwRemove
	kwUpdate
	kwMapReduce
	kwFindAndModify
	kwExplain
	kwAnalyzeShardKey
	kwConfigureQueryAnalyzer
	kwCompactStructuredEncryptionData
	kwHideIndex
	kwUnhideIndex
	kwReIndex
	kwGetShardDistribution
	kwGetShardVersion
	kwCreateSearchIndex
	kwCreateSearchIndexes
	kwDropSearchIndex
	kwUpdateSearchIndex

	// Cursor methods
	kwSort
	kwLimit
	kwSkip
	kwProjection
	kwProject
	kwBatchSize
	kwClose
	kwCollation
	kwComment
	kwForEach
	kwHasNext
	kwHint
	kwIsClosed
	kwIsExhausted
	kwItcount
	kwMap
	kwMax
	kwMaxAwaitTimeMS
	kwMaxTimeMS
	kwMin
	kwNext
	kwNoCursorTimeout
	kwObjsLeftInBatch
	kwPretty
	kwReadConcern
	kwReadPref
	kwReturnKey
	kwShowRecordId
	kwSize
	kwTailable
	kwToArray
	kwTryNext
	kwAllowDiskUse
	kwAddOption

	// Database methods
	kwGetCollection
	kwGetCollectionNames
	kwGetCollectionInfos
	kwCreateCollection
	kwDropDatabase
	kwHostInfo
	kwListCommands
	kwServerBuildInfo
	kwServerStatus
	kwVersion
	kwRunCommand
	kwAdminCommand
	kwGetName
	kwGetMongo
	kwGetSiblingDB
	kwAuth
	kwChangeUserPassword
	kwCloneDatabase
	kwCommandHelp
	kwCopyDatabase
	kwCreateRole
	kwCreateUser
	kwCreateView
	kwCurrentOp
	kwDropAllRoles
	kwDropAllUsers
	kwDropRole
	kwDropUser
	kwFsyncLock
	kwFsyncUnlock
	kwGetLogComponents
	kwGetProfilingLevel
	kwGetProfilingStatus
	kwGetReplicationInfo
	kwGetRole
	kwGetRoles
	kwGetUser
	kwGetUsers
	kwGrantPrivilegesToRole
	kwGrantRolesToRole
	kwGrantRolesToUser
	kwHello
	kwIsMaster
	kwKillOp
	kwLogout
	kwPrintCollectionStats
	kwPrintReplicationInfo
	kwPrintSecondaryReplicationInfo
	kwPrintShardingStatus
	kwPrintSlaveReplicationInfo
	kwRevokePrivilegesFromRole
	kwRevokeRolesFromRole
	kwRevokeRolesFromUser
	kwRotateCertificates
	kwSetLogLevel
	kwSetProfilingLevel
	kwSetSecondaryOk
	kwSetWriteConcern
	kwShutdownServer
	kwUpdateRole
	kwUpdateUser

	// Bulk operations
	kwInitializeOrderedBulkOp
	kwInitializeUnorderedBulkOp
	kwExecute
	kwGetOperations
	kwToString

	// Connection methods
	kwMongo
	kwConnect
	kwGetDB
	kwGetReadConcern
	kwGetReadPref
	kwGetReadPrefMode
	kwGetReadPrefTagSet
	kwGetWriteConcern
	kwSetReadPref
	kwSetReadConcern
	kwStartSession
	kwGetDBNames

	// Replication
	kwRs

	// Sharding
	kwSh

	// Stream processing
	kwSp

	// Encryption
	kwGetKeyVault
	kwGetClientEncryption

	// Plan cache
	kwGetPlanCache

	// Helper functions
	kwObjectId
	kwISODate
	kwDate
	kwUUID
	kwLong
	kwNumberLong
	kwInt32
	kwNumberInt
	kwDouble
	kwDecimal128
	kwNumberDecimal
	kwTimestamp
	kwRegExp
	kwBinData
	kwBinary
	kwBSONRegExp
	kwHexData
)

// keywords maps case-sensitive keyword strings to token types.
// mongosh keywords are case-sensitive (unlike SQL).
var keywords = map[string]int{
	"show":                           kwShow,
	"dbs":                            kwDbs,
	"databases":                      kwDatabases,
	"collections":                    kwCollections,
	"db":                             kwDb,
	"new":                            kwNew,
	"true":                           kwTrue,
	"false":                          kwFalse,
	"null":                           kwNull,
	"find":                           kwFind,
	"findOne":                        kwFindOne,
	"countDocuments":                 kwCountDocuments,
	"estimatedDocumentCount":         kwEstimatedDocumentCount,
	"distinct":                       kwDistinct,
	"aggregate":                      kwAggregate,
	"getIndexes":                     kwGetIndexes,
	"insertOne":                      kwInsertOne,
	"insertMany":                     kwInsertMany,
	"updateOne":                      kwUpdateOne,
	"updateMany":                     kwUpdateMany,
	"deleteOne":                      kwDeleteOne,
	"deleteMany":                     kwDeleteMany,
	"replaceOne":                     kwReplaceOne,
	"findOneAndUpdate":               kwFindOneAndUpdate,
	"findOneAndReplace":              kwFindOneAndReplace,
	"findOneAndDelete":               kwFindOneAndDelete,
	"createIndex":                    kwCreateIndex,
	"createIndexes":                  kwCreateIndexes,
	"dropIndex":                      kwDropIndex,
	"dropIndexes":                    kwDropIndexes,
	"drop":                           kwDrop,
	"renameCollection":               kwRenameCollection,
	"stats":                          kwStats,
	"storageSize":                    kwStorageSize,
	"totalIndexSize":                 kwTotalIndexSize,
	"totalSize":                      kwTotalSize,
	"dataSize":                       kwDataSize,
	"isCapped":                       kwIsCapped,
	"validate":                       kwValidate,
	"latencyStats":                   kwLatencyStats,
	"watch":                          kwWatch,
	"bulkWrite":                      kwBulkWrite,
	"count":                          kwCount,
	"insert":                         kwInsert,
	"remove":                         kwRemove,
	"update":                         kwUpdate,
	"mapReduce":                      kwMapReduce,
	"findAndModify":                  kwFindAndModify,
	"explain":                        kwExplain,
	"analyzeShardKey":                kwAnalyzeShardKey,
	"configureQueryAnalyzer":         kwConfigureQueryAnalyzer,
	"compactStructuredEncryptionData": kwCompactStructuredEncryptionData,
	"hideIndex":                      kwHideIndex,
	"unhideIndex":                    kwUnhideIndex,
	"reIndex":                        kwReIndex,
	"getShardDistribution":           kwGetShardDistribution,
	"getShardVersion":                kwGetShardVersion,
	"createSearchIndex":              kwCreateSearchIndex,
	"createSearchIndexes":            kwCreateSearchIndexes,
	"dropSearchIndex":                kwDropSearchIndex,
	"updateSearchIndex":              kwUpdateSearchIndex,
	"sort":                           kwSort,
	"limit":                          kwLimit,
	"skip":                           kwSkip,
	"projection":                     kwProjection,
	"project":                        kwProject,
	"batchSize":                      kwBatchSize,
	"close":                          kwClose,
	"collation":                      kwCollation,
	"comment":                        kwComment,
	"forEach":                        kwForEach,
	"hasNext":                        kwHasNext,
	"hint":                           kwHint,
	"isClosed":                       kwIsClosed,
	"isExhausted":                    kwIsExhausted,
	"itcount":                        kwItcount,
	"map":                            kwMap,
	"max":                            kwMax,
	"maxAwaitTimeMS":                 kwMaxAwaitTimeMS,
	"maxTimeMS":                      kwMaxTimeMS,
	"min":                            kwMin,
	"next":                           kwNext,
	"noCursorTimeout":                kwNoCursorTimeout,
	"objsLeftInBatch":                kwObjsLeftInBatch,
	"pretty":                         kwPretty,
	"readConcern":                    kwReadConcern,
	"readPref":                       kwReadPref,
	"returnKey":                      kwReturnKey,
	"showRecordId":                   kwShowRecordId,
	"size":                           kwSize,
	"tailable":                       kwTailable,
	"toArray":                        kwToArray,
	"tryNext":                        kwTryNext,
	"allowDiskUse":                   kwAllowDiskUse,
	"addOption":                      kwAddOption,
	"getCollection":                  kwGetCollection,
	"getCollectionNames":             kwGetCollectionNames,
	"getCollectionInfos":             kwGetCollectionInfos,
	"createCollection":               kwCreateCollection,
	"dropDatabase":                   kwDropDatabase,
	"hostInfo":                       kwHostInfo,
	"listCommands":                   kwListCommands,
	"serverBuildInfo":                kwServerBuildInfo,
	"serverStatus":                   kwServerStatus,
	"version":                        kwVersion,
	"runCommand":                     kwRunCommand,
	"adminCommand":                   kwAdminCommand,
	"getName":                        kwGetName,
	"getMongo":                       kwGetMongo,
	"getSiblingDB":                   kwGetSiblingDB,
	"auth":                           kwAuth,
	"changeUserPassword":             kwChangeUserPassword,
	"cloneDatabase":                  kwCloneDatabase,
	"commandHelp":                    kwCommandHelp,
	"copyDatabase":                   kwCopyDatabase,
	"createRole":                     kwCreateRole,
	"createUser":                     kwCreateUser,
	"createView":                     kwCreateView,
	"currentOp":                      kwCurrentOp,
	"dropAllRoles":                   kwDropAllRoles,
	"dropAllUsers":                   kwDropAllUsers,
	"dropRole":                       kwDropRole,
	"dropUser":                       kwDropUser,
	"fsyncLock":                      kwFsyncLock,
	"fsyncUnlock":                    kwFsyncUnlock,
	"getLogComponents":               kwGetLogComponents,
	"getProfilingLevel":              kwGetProfilingLevel,
	"getProfilingStatus":             kwGetProfilingStatus,
	"getReplicationInfo":             kwGetReplicationInfo,
	"getRole":                        kwGetRole,
	"getRoles":                       kwGetRoles,
	"getUser":                        kwGetUser,
	"getUsers":                       kwGetUsers,
	"grantPrivilegesToRole":          kwGrantPrivilegesToRole,
	"grantRolesToRole":               kwGrantRolesToRole,
	"grantRolesToUser":               kwGrantRolesToUser,
	"hello":                          kwHello,
	"isMaster":                       kwIsMaster,
	"killOp":                         kwKillOp,
	"logout":                         kwLogout,
	"printCollectionStats":           kwPrintCollectionStats,
	"printReplicationInfo":           kwPrintReplicationInfo,
	"printSecondaryReplicationInfo":  kwPrintSecondaryReplicationInfo,
	"printShardingStatus":            kwPrintShardingStatus,
	"printSlaveReplicationInfo":      kwPrintSlaveReplicationInfo,
	"revokePrivilegesFromRole":       kwRevokePrivilegesFromRole,
	"revokeRolesFromRole":            kwRevokeRolesFromRole,
	"revokeRolesFromUser":            kwRevokeRolesFromUser,
	"rotateCertificates":             kwRotateCertificates,
	"setLogLevel":                    kwSetLogLevel,
	"setProfilingLevel":              kwSetProfilingLevel,
	"setSecondaryOk":                 kwSetSecondaryOk,
	"setWriteConcern":                kwSetWriteConcern,
	"shutdownServer":                 kwShutdownServer,
	"updateRole":                     kwUpdateRole,
	"updateUser":                     kwUpdateUser,
	"initializeOrderedBulkOp":        kwInitializeOrderedBulkOp,
	"initializeUnorderedBulkOp":      kwInitializeUnorderedBulkOp,
	"execute":                        kwExecute,
	"getOperations":                  kwGetOperations,
	"toString":                       kwToString,
	"Mongo":                          kwMongo,
	"connect":                        kwConnect,
	"getDB":                          kwGetDB,
	"getReadConcern":                 kwGetReadConcern,
	"getReadPref":                    kwGetReadPref,
	"getReadPrefMode":                kwGetReadPrefMode,
	"getReadPrefTagSet":              kwGetReadPrefTagSet,
	"getWriteConcern":                kwGetWriteConcern,
	"setReadPref":                    kwSetReadPref,
	"setReadConcern":                 kwSetReadConcern,
	"startSession":                   kwStartSession,
	"getDBNames":                     kwGetDBNames,
	"rs":                             kwRs,
	"sh":                             kwSh,
	"sp":                             kwSp,
	"getKeyVault":                    kwGetKeyVault,
	"getClientEncryption":            kwGetClientEncryption,
	"getPlanCache":                   kwGetPlanCache,
	"ObjectId":                       kwObjectId,
	"ISODate":                        kwISODate,
	"Date":                           kwDate,
	"UUID":                           kwUUID,
	"Long":                           kwLong,
	"NumberLong":                     kwNumberLong,
	"Int32":                          kwInt32,
	"NumberInt":                      kwNumberInt,
	"Double":                         kwDouble,
	"Decimal128":                     kwDecimal128,
	"NumberDecimal":                  kwNumberDecimal,
	"Timestamp":                      kwTimestamp,
	"RegExp":                         kwRegExp,
	"BinData":                        kwBinData,
	"Binary":                         kwBinary,
	"BSONRegExp":                     kwBSONRegExp,
	"HexData":                        kwHexData,
}

// isKeyword returns true if the token type is a keyword (>= 700).
func isKeyword(typ int) bool {
	return typ >= 700
}

// tokenName returns a human-readable name for a token type (for error messages).
func tokenName(typ int) string {
	switch typ {
	case tokEOF:
		return "EOF"
	case tokNumber:
		return "number"
	case tokString:
		return "string"
	case tokRegex:
		return "regex"
	case tokIdent:
		return "identifier"
	case tokLParen:
		return "'('"
	case tokRParen:
		return "')'"
	case tokLBrace:
		return "'{'"
	case tokRBrace:
		return "'}'"
	case tokLBracket:
		return "'['"
	case tokRBracket:
		return "']'"
	case tokColon:
		return "':'"
	case tokComma:
		return "','"
	case tokDot:
		return "'.'"
	case tokSemi:
		return "';'"
	default:
		// For keywords, look up by iterating (only used in error paths)
		for name, t := range keywords {
			if t == typ {
				return "'" + name + "'"
			}
		}
		return "unknown"
	}
}
```

- [ ] **Step 2: Implement the lexer in lexer.go**

```go
// mongo/parser/lexer.go
package parser

import (
	"unicode/utf8"
)

// Lexer tokenizes MongoDB shell input.
type Lexer struct {
	input string
	pos   int // current byte position
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
	}

	start := l.pos
	ch := l.input[l.pos]

	switch ch {
	case '(':
		l.pos++
		return Token{Type: tokLParen, Str: "(", Loc: start, End: l.pos}
	case ')':
		l.pos++
		return Token{Type: tokRParen, Str: ")", Loc: start, End: l.pos}
	case '{':
		l.pos++
		return Token{Type: tokLBrace, Str: "{", Loc: start, End: l.pos}
	case '}':
		l.pos++
		return Token{Type: tokRBrace, Str: "}", Loc: start, End: l.pos}
	case '[':
		l.pos++
		return Token{Type: tokLBracket, Str: "[", Loc: start, End: l.pos}
	case ']':
		l.pos++
		return Token{Type: tokRBracket, Str: "]", Loc: start, End: l.pos}
	case ':':
		l.pos++
		return Token{Type: tokColon, Str: ":", Loc: start, End: l.pos}
	case ',':
		l.pos++
		return Token{Type: tokComma, Str: ",", Loc: start, End: l.pos}
	case '.':
		// Check if this starts a number like .001
		if l.pos+1 < len(l.input) && l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9' {
			return l.scanNumber(start)
		}
		l.pos++
		return Token{Type: tokDot, Str: ".", Loc: start, End: l.pos}
	case ';':
		l.pos++
		return Token{Type: tokSemi, Str: ";", Loc: start, End: l.pos}
	case '"':
		return l.scanDoubleQuotedString(start)
	case '\'':
		return l.scanSingleQuotedString(start)
	case '/':
		// Could be regex literal (line/block comments already consumed)
		return l.scanRegex(start)
	}

	// Numbers (including negative)
	if ch >= '0' && ch <= '9' {
		return l.scanNumber(start)
	}
	if ch == '-' && l.pos+1 < len(l.input) {
		next := l.input[l.pos+1]
		if next >= '0' && next <= '9' || next == '.' {
			return l.scanNumber(start)
		}
	}

	// Identifiers and keywords (including $-prefixed)
	if ch == '$' || ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch >= 0x80 {
		return l.scanIdentifier(start)
	}

	// Unknown character - return as identifier to let parser produce error
	l.pos++
	return Token{Type: tokIdent, Str: l.input[start:l.pos], Loc: start, End: l.pos}
}

// skipWhitespaceAndComments skips spaces, tabs, newlines, and comments.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Whitespace
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.pos++
			continue
		}

		// Line comment: //
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}

		// Block comment: /* ... */
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			l.pos += 2
			for l.pos+1 < len(l.input) {
				if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
					l.pos += 2
					break
				}
				l.pos++
			}
			continue
		}

		break
	}
}

// scanDoubleQuotedString scans a "..." string starting at pos (which points to the opening ").
func (l *Lexer) scanDoubleQuotedString(start int) Token {
	l.pos++ // skip opening "
	var result []byte
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			l.pos++ // skip closing "
			return Token{Type: tokString, Str: string(result), Loc: start, End: l.pos}
		}
		if ch == '\\' {
			l.pos++
			if l.pos < len(l.input) {
				esc := l.input[l.pos]
				switch esc {
				case '"', '\\', '/':
					result = append(result, esc)
				case '\'':
					result = append(result, '\'')
				case 'b':
					result = append(result, '\b')
				case 'f':
					result = append(result, '\f')
				case 'n':
					result = append(result, '\n')
				case 'r':
					result = append(result, '\r')
				case 't':
					result = append(result, '\t')
				case 'u':
					// Unicode escape: \uXXXX
					if l.pos+4 < len(l.input) {
						hex := l.input[l.pos+1 : l.pos+5]
						r := parseHex4(hex)
						if r >= 0 {
							var buf [4]byte
							n := utf8.EncodeRune(buf[:], rune(r))
							result = append(result, buf[:n]...)
							l.pos += 4
						} else {
							result = append(result, '\\', 'u')
						}
					} else {
						result = append(result, '\\', 'u')
					}
				default:
					result = append(result, '\\', esc)
				}
				l.pos++
				continue
			}
		}
		result = append(result, ch)
		l.pos++
	}
	// Unterminated string
	return Token{Type: tokString, Str: string(result), Loc: start, End: l.pos}
}

// scanSingleQuotedString scans a '...' string starting at pos (which points to the opening ').
func (l *Lexer) scanSingleQuotedString(start int) Token {
	l.pos++ // skip opening '
	var result []byte
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			l.pos++ // skip closing '
			return Token{Type: tokString, Str: string(result), Loc: start, End: l.pos}
		}
		if ch == '\\' {
			l.pos++
			if l.pos < len(l.input) {
				esc := l.input[l.pos]
				switch esc {
				case '\'', '\\', '/':
					result = append(result, esc)
				case '"':
					result = append(result, '"')
				case 'b':
					result = append(result, '\b')
				case 'f':
					result = append(result, '\f')
				case 'n':
					result = append(result, '\n')
				case 'r':
					result = append(result, '\r')
				case 't':
					result = append(result, '\t')
				case 'u':
					if l.pos+4 < len(l.input) {
						hex := l.input[l.pos+1 : l.pos+5]
						r := parseHex4(hex)
						if r >= 0 {
							var buf [4]byte
							n := utf8.EncodeRune(buf[:], rune(r))
							result = append(result, buf[:n]...)
							l.pos += 4
						} else {
							result = append(result, '\\', 'u')
						}
					} else {
						result = append(result, '\\', 'u')
					}
				default:
					result = append(result, '\\', esc)
				}
				l.pos++
				continue
			}
		}
		result = append(result, ch)
		l.pos++
	}
	return Token{Type: tokString, Str: string(result), Loc: start, End: l.pos}
}

// scanNumber scans a numeric literal. Handles integers, floats, scientific notation, negative.
// Grammar: -? INT ('.' [0-9]+)? ([eE] [+-]? [0-9]+)?  |  -? '.' [0-9]+ ([eE] [+-]? [0-9]+)?
func (l *Lexer) scanNumber(start int) Token {
	if l.pos < len(l.input) && l.input[l.pos] == '-' {
		l.pos++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		// .NNN form
		l.pos++
		l.scanDigits()
	} else {
		l.scanDigits()
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			l.pos++
			l.scanDigits()
		}
	}
	// Exponent
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		l.scanDigits()
	}
	return Token{Type: tokNumber, Str: l.input[start:l.pos], Loc: start, End: l.pos}
}

func (l *Lexer) scanDigits() {
	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
		l.pos++
	}
}

// scanIdentifier scans an identifier or keyword.
// Identifiers: [$_a-zA-Z][$_a-zA-Z0-9]* (also supports Unicode chars >= 0x80).
func (l *Lexer) scanIdentifier(start int) Token {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '$' || ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			l.pos++
		} else if ch >= 0x80 {
			// Consume full UTF-8 character
			_, size := utf8.DecodeRuneInString(l.input[l.pos:])
			l.pos += size
		} else {
			break
		}
	}

	text := l.input[start:l.pos]

	// Check if it's a keyword
	if kw, ok := keywords[text]; ok {
		return Token{Type: kw, Str: text, Loc: start, End: l.pos}
	}

	return Token{Type: tokIdent, Str: text, Loc: start, End: l.pos}
}

// scanRegex scans a regex literal: /pattern/flags.
// Caller has verified we're at '/' and it's not a comment.
func (l *Lexer) scanRegex(start int) Token {
	l.pos++ // skip opening /
	var pattern []byte
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '/' {
			l.pos++ // skip closing /
			// Scan optional flags
			flagStart := l.pos
			for l.pos < len(l.input) {
				f := l.input[l.pos]
				if f == 'g' || f == 'i' || f == 'm' || f == 's' || f == 'u' || f == 'y' {
					l.pos++
				} else {
					break
				}
			}
			flags := l.input[flagStart:l.pos]
			return Token{Type: tokRegex, Str: string(pattern) + "/" + flags, Loc: start, End: l.pos}
		}
		if ch == '\\' && l.pos+1 < len(l.input) {
			pattern = append(pattern, ch, l.input[l.pos+1])
			l.pos += 2
			continue
		}
		if ch == '\r' || ch == '\n' {
			break // unterminated regex
		}
		pattern = append(pattern, ch)
		l.pos++
	}
	// Unterminated regex - return what we have
	return Token{Type: tokRegex, Str: string(pattern), Loc: start, End: l.pos}
}

// parseHex4 parses a 4-character hex string into an int. Returns -1 on failure.
func parseHex4(s string) int {
	if len(s) != 4 {
		return -1
	}
	var val int
	for _, c := range s {
		val <<= 4
		switch {
		case c >= '0' && c <= '9':
			val |= int(c - '0')
		case c >= 'a' && c <= 'f':
			val |= int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			val |= int(c-'A') + 10
		default:
			return -1
		}
	}
	return val
}
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/parser/`
Expected: SUCCESS (will fail until parser.go exists — create a minimal stub)

Create a minimal `mongo/parser/parser.go` stub so the package compiles:

```go
// mongo/parser/parser.go
package parser
```

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/parser/`
Expected: SUCCESS

- [ ] **Step 4: Commit**

```bash
git add mongo/parser/token.go mongo/parser/lexer.go mongo/parser/parser.go
git commit -m "feat(mongo/parser): add token types and lexer

Implements character-by-character tokenization for mongosh syntax.
Handles strings, numbers, regex literals, identifiers, keywords,
comments, and all punctuation. Case-sensitive keyword matching."
```

---

## Task 3: Parser Core and Public API

**Files:**
- Modify: `mongo/parser/parser.go`
- Create: `mongo/parse.go`

- [ ] **Step 1: Implement Parser struct with advance/peek/expect and error handling**

Replace the stub `mongo/parser/parser.go` with the full parser core:

```go
// mongo/parser/parser.go
package parser

import (
	"fmt"

	"github.com/bytebase/omni/mongo/ast"
)

// Parser is a recursive descent parser for MongoDB shell syntax.
type Parser struct {
	lexer   *Lexer
	input   string
	cur     Token // current token
	prev    Token // previous token
	nextBuf Token // buffered next token for lookahead
	hasNext bool  // whether nextBuf is valid
}

// Parse parses a mongosh string into a list of AST nodes.
func Parse(input string) ([]ast.Node, error) {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	var stmts []ast.Node
	for p.cur.Type != tokEOF {
		if p.cur.Type == tokSemi {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	return stmts, nil
}

// parseStatement dispatches to the appropriate statement parser based on the current token.
func (p *Parser) parseStatement() (ast.Node, error) {
	switch p.cur.Type {
	case kwShow:
		return p.parseShellCommand()
	case kwDb:
		return p.parseDbStatement()
	case kwRs:
		return p.parseRsStatement()
	case kwSh:
		return p.parseShStatement()
	case kwSp:
		return p.parseSpStatement()
	case kwMongo:
		return p.parseConnectionStatement()
	case kwConnect:
		return p.parseConnectionStatement()
	default:
		// Could be a native function call: identifier(...)
		if p.isIdentLike(p.cur.Type) && p.peekNext().Type == tokLParen {
			return p.parseNativeFunctionCall()
		}
		return nil, p.syntaxErrorAtCur()
	}
}

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// expect consumes the current token if it matches the expected type.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.unexpectedTokenError(tokenType)
}

// match checks if the current token matches any of the given types and consumes it.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// isIdentLike returns true if the token type can be used as an identifier.
// All keywords in mongosh can also serve as identifiers (collection names, field names).
func (p *Parser) isIdentLike(typ int) bool {
	return typ == tokIdent || isKeyword(typ)
}

// identText returns the string for an identifier-like token (identifiers or keywords).
func (p *Parser) identText(tok Token) string {
	return tok.Str
}

// expectIdent consumes the current token if it is identifier-like.
func (p *Parser) expectIdent() (Token, error) {
	if p.isIdentLike(p.cur.Type) {
		return p.advance(), nil
	}
	return Token{}, &ParseError{
		Message:  "expected identifier",
		Position: p.cur.Loc,
	}
}

// parseArguments parses a comma-separated list of values (between parens).
// The opening '(' must already be consumed. Consumes the closing ')'.
func (p *Parser) parseArguments() ([]ast.Node, error) {
	var args []ast.Node
	if p.cur.Type == tokRParen {
		p.advance()
		return args, nil
	}
	for {
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		args = append(args, val)
		if p.cur.Type == tokComma {
			p.advance()
			// Allow trailing comma before )
			if p.cur.Type == tokRParen {
				break
			}
			continue
		}
		break
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}
	return args, nil
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int // byte offset
	Line     int // 1-based
	Column   int // 1-based
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s (line %d, column %d)", e.Message, e.Line, e.Column)
	}
	return e.Message
}

// syntaxErrorAtCur returns a ParseError at the current token.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	line, col := p.lineCol(p.cur.Loc)
	var msg string
	if p.cur.Type == tokEOF {
		msg = "syntax error at end of input"
	} else {
		msg = fmt.Sprintf("syntax error at or near %q", p.cur.Str)
	}
	return &ParseError{
		Message:  msg,
		Position: p.cur.Loc,
		Line:     line,
		Column:   col,
	}
}

// unexpectedTokenError returns an error for an unexpected token when a specific type was expected.
func (p *Parser) unexpectedTokenError(expected int) *ParseError {
	line, col := p.lineCol(p.cur.Loc)
	var msg string
	if p.cur.Type == tokEOF {
		msg = fmt.Sprintf("expected %s, got end of input", tokenName(expected))
	} else {
		msg = fmt.Sprintf("expected %s, got %q", tokenName(expected), p.cur.Str)
	}
	return &ParseError{
		Message:  msg,
		Position: p.cur.Loc,
		Line:     line,
		Column:   col,
	}
}

// lineCol computes 1-based line and column for a byte offset.
func (p *Parser) lineCol(offset int) (int, int) {
	if offset > len(p.input) {
		offset = len(p.input)
	}
	line := 1
	lastNewline := -1
	for i := 0; i < offset; i++ {
		if p.input[i] == '\n' {
			line++
			lastNewline = i
		}
	}
	col := offset - lastNewline
	return line, col
}
```

- [ ] **Step 2: Create the public API in parse.go**

```go
// mongo/parse.go
package mongo

import (
	"github.com/bytebase/omni/mongo/ast"
	"github.com/bytebase/omni/mongo/parser"
)

// Statement is the result of parsing a single MongoDB shell statement.
type Statement struct {
	// Text is the original text of this statement.
	Text string
	// AST is the parsed node.
	AST ast.Node
	// ByteStart is the inclusive start byte offset in the original input.
	ByteStart int
	// ByteEnd is the exclusive end byte offset in the original input.
	ByteEnd int
	// Start is the start position (line:column).
	Start Position
	// End is the end position (line:column).
	End Position
}

// Position represents a location in source text.
type Position struct {
	Line   int // 1-based
	Column int // 1-based
}

// Parse splits and parses a MongoDB shell string into statements.
func Parse(input string) ([]Statement, error) {
	nodes, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}

	lineIndex := buildLineIndex(input)
	var stmts []Statement
	for _, node := range nodes {
		loc := node.GetLoc()
		stmts = append(stmts, Statement{
			Text:      input[loc.Start:loc.End],
			AST:       node,
			ByteStart: loc.Start,
			ByteEnd:   loc.End,
			Start:     offsetToPosition(lineIndex, loc.Start),
			End:       offsetToPosition(lineIndex, loc.End),
		})
	}
	return stmts, nil
}

type lineIndex []int

func buildLineIndex(s string) lineIndex {
	idx := lineIndex{0}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			idx = append(idx, i+1)
		}
	}
	return idx
}

func offsetToPosition(idx lineIndex, offset int) Position {
	lo, hi := 0, len(idx)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if idx[mid] <= offset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return Position{
		Line:   lo + 1,
		Column: offset - idx[lo] + 1,
	}
}
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/...`
Expected: Will fail because `parseValue`, `parseShellCommand`, etc. are not yet defined. Create stubs for all referenced methods in their respective files.

Create these stub files so compilation succeeds:

`mongo/parser/shell.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseShellCommand() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/database.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseDbStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/collection.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseCollectionStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/expression.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseValue() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/document.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseDocument() (*ast.Document, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/replication.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseRsStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/sharding.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseShStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/stream.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseSpStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/connection.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseConnectionStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/bulk.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseBulkStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/encryption.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseEncryptionStatement(stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/plancache.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parsePlanCacheStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

`mongo/parser/native.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseNativeFunctionCall() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
```

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./mongo/...`
Expected: SUCCESS

- [ ] **Step 4: Commit**

```bash
git add mongo/parse.go mongo/parser/parser.go mongo/parser/shell.go mongo/parser/database.go mongo/parser/collection.go mongo/parser/expression.go mongo/parser/document.go mongo/parser/replication.go mongo/parser/sharding.go mongo/parser/stream.go mongo/parser/connection.go mongo/parser/bulk.go mongo/parser/encryption.go mongo/parser/plancache.go mongo/parser/native.go
git commit -m "feat(mongo): parser core and public API

Implements Parser struct with advance/peek/expect, top-level statement
dispatch, ParseError with line/column, and the public mongo.Parse() API.
Statement family parsers are stubs to be implemented in subsequent tasks."
```

---

## Task 4: Expression and Document Parsing

**Files:**
- Modify: `mongo/parser/expression.go`
- Modify: `mongo/parser/document.go`
- Create: `mongo/parsertest/helpers_test.go`
- Create: `mongo/parsertest/expression_test.go`
- Create: `mongo/parsertest/document_test.go`

These are the foundation for all other parsers — every statement contains documents and values.

- [ ] **Step 1: Write tests for value parsing**

```go
// mongo/parsertest/helpers_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

// mustParse parses input and requires exactly one statement with no error.
func mustParse(t *testing.T, input string) mongo.Statement {
	t.Helper()
	stmts, err := mongo.Parse(input)
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	return stmts[0]
}

// mustParseN parses input and requires exactly n statements with no error.
func mustParseN(t *testing.T, input string, n int) []mongo.Statement {
	t.Helper()
	stmts, err := mongo.Parse(input)
	require.NoError(t, err)
	require.Len(t, stmts, n)
	return stmts
}

// mustFail parses input and requires a parse error.
func mustFail(t *testing.T, input string) error {
	t.Helper()
	_, err := mongo.Parse(input)
	require.Error(t, err)
	return err
}

// assertLoc checks that a node's Loc matches expected start and end byte offsets.
func assertLoc(t *testing.T, node ast.Node, start, end int) {
	t.Helper()
	loc := node.GetLoc()
	require.Equal(t, start, loc.Start, "Loc.Start mismatch")
	require.Equal(t, end, loc.End, "Loc.End mismatch")
}
```

```go
// mongo/parsertest/expression_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseLiterals(t *testing.T) {
	// Each test wraps the literal in db.c.find({k: VALUE}) and checks the value node.
	tests := []struct {
		input    string
		nodeType string // expected ast type
	}{
		{`db.c.find({k: "hello"})`, "StringLiteral"},
		{`db.c.find({k: 'world'})`, "StringLiteral"},
		{`db.c.find({k: 25})`, "NumberLiteral"},
		{`db.c.find({k: 3.14})`, "NumberLiteral"},
		{`db.c.find({k: -10})`, "NumberLiteral"},
		{`db.c.find({k: 1.5e10})`, "NumberLiteral"},
		{`db.c.find({k: .001})`, "NumberLiteral"},
		{`db.c.find({k: true})`, "BoolLiteral"},
		{`db.c.find({k: false})`, "BoolLiteral"},
		{`db.c.find({k: null})`, "NullLiteral"},
		{`db.c.find({k: /^test/i})`, "RegexLiteral"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs, ok := stmt.AST.(*ast.CollectionStatement)
			require.True(t, ok)
			require.Len(t, cs.Args, 1)
			doc, ok := cs.Args[0].(*ast.Document)
			require.True(t, ok)
			require.Len(t, doc.Pairs, 1)

			switch tt.nodeType {
			case "StringLiteral":
				_, ok := doc.Pairs[0].Value.(*ast.StringLiteral)
				require.True(t, ok)
			case "NumberLiteral":
				_, ok := doc.Pairs[0].Value.(*ast.NumberLiteral)
				require.True(t, ok)
			case "BoolLiteral":
				_, ok := doc.Pairs[0].Value.(*ast.BoolLiteral)
				require.True(t, ok)
			case "NullLiteral":
				_, ok := doc.Pairs[0].Value.(*ast.NullLiteral)
				require.True(t, ok)
			case "RegexLiteral":
				_, ok := doc.Pairs[0].Value.(*ast.RegexLiteral)
				require.True(t, ok)
			}
		})
	}
}

func TestParseHelperFunctions(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{`db.c.find({_id: ObjectId("507f1f77bcf86cd799439011")})`, "ObjectId"},
		{`db.c.find({_id: ObjectId()})`, "ObjectId"},
		{`db.c.find({d: ISODate("2024-01-15T00:00:00Z")})`, "ISODate"},
		{`db.c.find({d: ISODate()})`, "ISODate"},
		{`db.c.find({d: Date()})`, "Date"},
		{`db.c.find({d: Date("2024-01-15")})`, "Date"},
		{`db.c.find({d: Date(1705276800000)})`, "Date"},
		{`db.c.find({s: UUID("550e8400-e29b-41d4-a716-446655440000")})`, "UUID"},
		{`db.c.find({n: Long(123)})`, "Long"},
		{`db.c.find({n: NumberLong("123456789012345")})`, "NumberLong"},
		{`db.c.find({n: Int32(100)})`, "Int32"},
		{`db.c.find({n: NumberInt(100)})`, "NumberInt"},
		{`db.c.find({v: Double(3.14)})`, "Double"},
		{`db.c.find({a: Decimal128("99.99")})`, "Decimal128"},
		{`db.c.find({a: NumberDecimal("99.99")})`, "NumberDecimal"},
		{`db.c.find({t: Timestamp(1627811580, 1)})`, "Timestamp"},
		{`db.c.find({b: BinData(0, "AQID")})`, "BinData"},
		{`db.c.find({r: BSONRegExp("pattern", "i")})`, "BSONRegExp"},
		{`db.c.find({h: HexData(0, "0123456789abcdef")})`, "HexData"},
		{`db.c.find({r: RegExp("^test", "i")})`, "RegExp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs := stmt.AST.(*ast.CollectionStatement)
			doc := cs.Args[0].(*ast.Document)
			helper, ok := doc.Pairs[0].Value.(*ast.HelperCall)
			require.True(t, ok)
			require.Equal(t, tt.name, helper.Name)
		})
	}
}

func TestParseArrays(t *testing.T) {
	stmt := mustParse(t, `db.c.find({tags: ["a", "b", "c"]})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	arr, ok := doc.Pairs[0].Value.(*ast.Array)
	require.True(t, ok)
	require.Len(t, arr.Elements, 3)
}

func TestParseRegexConstructor(t *testing.T) {
	stmt := mustParse(t, `db.c.find({name: RegExp("^alice", "i")})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	helper, ok := doc.Pairs[0].Value.(*ast.HelperCall)
	require.True(t, ok)
	require.Equal(t, "RegExp", helper.Name)
	require.Len(t, helper.Args, 2)
}
```

```go
// mongo/parsertest/document_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseDocumentUnquotedKeys(t *testing.T) {
	stmt := mustParse(t, `db.c.find({name: "alice", age: 25})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	require.Len(t, doc.Pairs, 2)
	require.Equal(t, "name", doc.Pairs[0].Key)
	require.Equal(t, "age", doc.Pairs[1].Key)
}

func TestParseDocumentQuotedKeys(t *testing.T) {
	stmt := mustParse(t, `db.c.find({"name": "alice", 'age': 25})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	require.Len(t, doc.Pairs, 2)
	require.Equal(t, "name", doc.Pairs[0].Key)
	require.Equal(t, "age", doc.Pairs[1].Key)
}

func TestParseDocumentDollarKeys(t *testing.T) {
	stmt := mustParse(t, `db.c.find({age: {$gt: 25, $lt: 65}})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	inner, ok := doc.Pairs[0].Value.(*ast.Document)
	require.True(t, ok)
	require.Len(t, inner.Pairs, 2)
	require.Equal(t, "$gt", inner.Pairs[0].Key)
	require.Equal(t, "$lt", inner.Pairs[1].Key)
}

func TestParseDocumentTrailingComma(t *testing.T) {
	stmt := mustParse(t, `db.c.find({name: "alice", age: 25,})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	require.Len(t, doc.Pairs, 2)
}

func TestParseDocumentNested(t *testing.T) {
	stmt := mustParse(t, `db.c.find({profile: {name: "test", settings: {theme: "dark"}}})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	inner, ok := doc.Pairs[0].Value.(*ast.Document)
	require.True(t, ok)
	require.Len(t, inner.Pairs, 2)
	nested, ok := inner.Pairs[1].Value.(*ast.Document)
	require.True(t, ok)
	require.Equal(t, "theme", nested.Pairs[0].Key)
}

func TestParseDocumentEmpty(t *testing.T) {
	stmt := mustParse(t, `db.c.find({})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	require.Len(t, doc.Pairs, 0)
}

func TestParseArrayTrailingComma(t *testing.T) {
	stmt := mustParse(t, `db.c.find({tags: ["a", "b", "c",]})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	arr := doc.Pairs[0].Value.(*ast.Array)
	require.Len(t, arr.Elements, 3)
}
```

- [ ] **Step 2: Implement parseValue in expression.go**

Replace the stub `mongo/parser/expression.go`:

```go
// mongo/parser/expression.go
package parser

import (
	"strings"

	"github.com/bytebase/omni/mongo/ast"
)

// helperNames is the set of recognized BSON helper/constructor function names.
var helperNames = map[int]bool{
	kwObjectId:      true,
	kwISODate:       true,
	kwDate:          true,
	kwUUID:          true,
	kwLong:          true,
	kwNumberLong:    true,
	kwInt32:         true,
	kwNumberInt:     true,
	kwDouble:        true,
	kwDecimal128:    true,
	kwNumberDecimal: true,
	kwTimestamp:     true,
	kwRegExp:        true,
	kwBinData:       true,
	kwBinary:        true,
	kwBSONRegExp:    true,
	kwHexData:       true,
}

// parseValue parses a value: document, array, helper function, regex, or literal.
func (p *Parser) parseValue() (ast.Node, error) {
	switch p.cur.Type {
	case tokLBrace:
		return p.parseDocument()
	case tokLBracket:
		return p.parseArray()
	case tokRegex:
		return p.parseRegexLiteral()
	case tokString:
		return p.parseStringLiteral(), nil
	case tokNumber:
		return p.parseNumberLiteral(), nil
	case kwTrue:
		tok := p.advance()
		return &ast.BoolLiteral{Value: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case kwFalse:
		tok := p.advance()
		return &ast.BoolLiteral{Value: false, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case kwNull:
		tok := p.advance()
		return &ast.NullLiteral{Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case kwNew:
		// "new" keyword — provide helpful error like the ANTLR grammar
		return nil, &ParseError{
			Message:  "'new' keyword is not supported. Use ObjectId(), ISODate(), UUID(), etc. directly without 'new'",
			Position: p.cur.Loc,
			Line:     p.lineColLine(p.cur.Loc),
			Column:   p.lineColCol(p.cur.Loc),
		}
	default:
		// Check for helper function call: ObjectId(...), ISODate(...), etc.
		if helperNames[p.cur.Type] {
			return p.parseHelperCall()
		}
		// Check for Binary.createFromBase64(...) pattern
		if p.cur.Type == kwBinary && p.peekNext().Type == tokDot {
			return p.parseHelperCall()
		}
		return nil, p.syntaxErrorAtCur()
	}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	tok := p.advance()
	return &ast.StringLiteral{Value: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}
}

func (p *Parser) parseNumberLiteral() *ast.NumberLiteral {
	tok := p.advance()
	isFloat := strings.ContainsAny(tok.Str, ".eE")
	return &ast.NumberLiteral{Value: tok.Str, IsFloat: isFloat, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}
}

func (p *Parser) parseRegexLiteral() (ast.Node, error) {
	tok := p.advance()
	// tok.Str is "pattern/flags" (the opening/closing / are stripped by lexer)
	parts := strings.SplitN(tok.Str, "/", 2)
	pattern := parts[0]
	flags := ""
	if len(parts) > 1 {
		flags = parts[1]
	}
	return &ast.RegexLiteral{Pattern: pattern, Flags: flags, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
}

func (p *Parser) parseHelperCall() (ast.Node, error) {
	start := p.cur.Loc
	name := p.cur.Str

	// Handle Binary.createFromBase64(...) special form
	if p.cur.Type == kwBinary && p.peekNext().Type == tokDot {
		p.advance() // consume "Binary"
		p.advance() // consume "."
		sub, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		name = "Binary." + sub.Str
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		return &ast.HelperCall{Name: name, Args: args, Loc: ast.Loc{Start: start, End: p.prev.End}}, nil
	}

	p.advance() // consume helper name
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &ast.HelperCall{Name: name, Args: args, Loc: ast.Loc{Start: start, End: p.prev.End}}, nil
}

func (p *Parser) parseArray() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume [

	var elements []ast.Node
	if p.cur.Type == tokRBracket {
		end := p.advance() // consume ]
		return &ast.Array{Elements: elements, Loc: ast.Loc{Start: start, End: end.End}}, nil
	}

	for {
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		elements = append(elements, val)
		if p.cur.Type == tokComma {
			p.advance()
			// Trailing comma
			if p.cur.Type == tokRBracket {
				break
			}
			continue
		}
		break
	}

	end, err := p.expect(tokRBracket)
	if err != nil {
		return nil, err
	}
	return &ast.Array{Elements: elements, Loc: ast.Loc{Start: start, End: end.End}}, nil
}

// lineColLine returns the 1-based line number for an offset (helper for error creation).
func (p *Parser) lineColLine(offset int) int {
	line, _ := p.lineCol(offset)
	return line
}

// lineColCol returns the 1-based column number for an offset (helper for error creation).
func (p *Parser) lineColCol(offset int) int {
	_, col := p.lineCol(offset)
	return col
}
```

- [ ] **Step 3: Implement parseDocument in document.go**

Replace the stub `mongo/parser/document.go`:

```go
// mongo/parser/document.go
package parser

import (
	"github.com/bytebase/omni/mongo/ast"
)

// parseDocument parses a BSON document: { key: value, ... }
// Supports unquoted keys, quoted keys, $-prefixed keys, and trailing commas.
func (p *Parser) parseDocument() (*ast.Document, error) {
	start := p.cur.Loc
	p.advance() // consume {

	var pairs []ast.KeyValue
	if p.cur.Type == tokRBrace {
		end := p.advance() // consume }
		return &ast.Document{Pairs: pairs, Loc: ast.Loc{Start: start, End: end.End}}, nil
	}

	for {
		kv, err := p.parsePair()
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, kv)
		if p.cur.Type == tokComma {
			p.advance()
			// Trailing comma before }
			if p.cur.Type == tokRBrace {
				break
			}
			continue
		}
		break
	}

	end, err := p.expect(tokRBrace)
	if err != nil {
		return nil, err
	}
	return &ast.Document{Pairs: pairs, Loc: ast.Loc{Start: start, End: end.End}}, nil
}

// parsePair parses a key: value pair.
func (p *Parser) parsePair() (ast.KeyValue, error) {
	pairStart := p.cur.Loc

	// Parse key: quoted string or identifier (including $-prefixed, keywords)
	var key string
	var keyLoc ast.Loc

	if p.cur.Type == tokString {
		tok := p.advance()
		key = tok.Str
		keyLoc = ast.Loc{Start: tok.Loc, End: tok.End}
	} else if p.isIdentLike(p.cur.Type) || p.cur.Type == tokIdent {
		tok := p.advance()
		key = tok.Str
		keyLoc = ast.Loc{Start: tok.Loc, End: tok.End}
	} else {
		return ast.KeyValue{}, &ParseError{
			Message:  "expected key (identifier or string)",
			Position: p.cur.Loc,
			Line:     p.lineColLine(p.cur.Loc),
			Column:   p.lineColCol(p.cur.Loc),
		}
	}

	// Expect colon
	if _, err := p.expect(tokColon); err != nil {
		return ast.KeyValue{}, err
	}

	// Parse value
	val, err := p.parseValue()
	if err != nil {
		return ast.KeyValue{}, err
	}

	valLoc := val.GetLoc()
	return ast.KeyValue{
		Key:    key,
		KeyLoc: keyLoc,
		Value:  val,
		Loc:    ast.Loc{Start: pairStart, End: valLoc.End},
	}, nil
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/parsertest/ -run TestParseLiterals -v`
Expected: FAIL (collection parsing not implemented yet — tests depend on it)

This is expected. We need Task 5 (shell commands) and Task 6 (collection parsing) before these tests pass. Continue to the next task.

- [ ] **Step 5: Commit**

```bash
git add mongo/parser/expression.go mongo/parser/document.go mongo/parsertest/
git commit -m "feat(mongo/parser): expression and document parsing

Implements parseValue (strings, numbers, booleans, null, regex,
arrays, documents, BSON helper calls) and parseDocument (unquoted
keys, quoted keys, dollar keys, trailing commas, nesting).
Adds test infrastructure and tests for expressions and documents."
```

---

## Task 5: Shell Commands (show dbs, show collections)

**Files:**
- Modify: `mongo/parser/shell.go`
- Create: `mongo/parsertest/shell_test.go`

- [ ] **Step 1: Write tests for shell commands**

```go
// mongo/parsertest/shell_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseShowDatabases(t *testing.T) {
	for _, input := range []string{"show dbs", "show databases"} {
		t.Run(input, func(t *testing.T) {
			stmt := mustParse(t, input)
			sc, ok := stmt.AST.(*ast.ShowCommand)
			require.True(t, ok)
			require.Equal(t, "databases", sc.Target)
		})
	}
}

func TestParseShowCollections(t *testing.T) {
	stmt := mustParse(t, "show collections")
	sc, ok := stmt.AST.(*ast.ShowCommand)
	require.True(t, ok)
	require.Equal(t, "collections", sc.Target)
}

func TestParseMultipleShellCommands(t *testing.T) {
	stmts := mustParseN(t, "show dbs; show collections", 2)
	sc1 := stmts[0].AST.(*ast.ShowCommand)
	require.Equal(t, "databases", sc1.Target)
	sc2 := stmts[1].AST.(*ast.ShowCommand)
	require.Equal(t, "collections", sc2.Target)
}
```

- [ ] **Step 2: Implement parseShellCommand**

Replace `mongo/parser/shell.go`:

```go
// mongo/parser/shell.go
package parser

import (
	"github.com/bytebase/omni/mongo/ast"
)

// parseShellCommand parses "show dbs", "show databases", or "show collections".
func (p *Parser) parseShellCommand() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume "show"

	switch p.cur.Type {
	case kwDbs, kwDatabases:
		end := p.advance()
		return &ast.ShowCommand{Target: "databases", Loc: ast.Loc{Start: start, End: end.End}}, nil
	case kwCollections:
		end := p.advance()
		return &ast.ShowCommand{Target: "collections", Loc: ast.Loc{Start: start, End: end.End}}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
```

- [ ] **Step 3: Run tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/parsertest/ -run TestParseShow -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add mongo/parser/shell.go mongo/parsertest/shell_test.go
git commit -m "feat(mongo/parser): show dbs/databases/collections commands"
```

---

## Task 6: Collection Statement Parsing

**Files:**
- Modify: `mongo/parser/database.go` (needs `parseDbStatement` to route to collection parsing)
- Modify: `mongo/parser/collection.go`
- Create: `mongo/parsertest/collection_test.go`

This is the largest and most complex parser. It handles `db.collection.method(...)` and cursor chaining.

- [ ] **Step 1: Write collection parsing tests**

```go
// mongo/parsertest/collection_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseBasicFind(t *testing.T) {
	tests := []struct {
		input      string
		collection string
		method     string
		argCount   int
	}{
		{"db.users.find()", "users", "find", 0},
		{"db.users.find({})", "users", "find", 1},
		{`db.users.find({name: "alice"})`, "users", "find", 1},
		{`db.users.find({age: {$gt: 25}})`, "users", "find", 1},
		{`db.users.findOne({name: "alice"})`, "users", "findOne", 1},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs, ok := stmt.AST.(*ast.CollectionStatement)
			require.True(t, ok)
			require.Equal(t, tt.collection, cs.Collection)
			require.Equal(t, tt.method, cs.Method)
			require.Len(t, cs.Args, tt.argCount)
		})
	}
}

func TestParseFindWithCursorMethods(t *testing.T) {
	stmt := mustParse(t, `db.users.find().sort({age: -1}).limit(10).skip(5)`)
	cs := stmt.AST.(*ast.CollectionStatement)
	require.Equal(t, "find", cs.Method)
	require.Len(t, cs.CursorMethods, 3)
	require.Equal(t, "sort", cs.CursorMethods[0].Method)
	require.Equal(t, "limit", cs.CursorMethods[1].Method)
	require.Equal(t, "skip", cs.CursorMethods[2].Method)
}

func TestParseWriteOperations(t *testing.T) {
	tests := []struct {
		input  string
		method string
	}{
		{`db.users.insertOne({name: "alice"})`, "insertOne"},
		{`db.users.insertMany([{name: "a"}, {name: "b"}])`, "insertMany"},
		{`db.users.updateOne({name: "a"}, {$set: {age: 30}})`, "updateOne"},
		{`db.users.updateMany({}, {$set: {active: true}})`, "updateMany"},
		{`db.users.deleteOne({name: "a"})`, "deleteOne"},
		{`db.users.deleteMany({status: "deleted"})`, "deleteMany"},
		{`db.users.replaceOne({name: "a"}, {name: "b", age: 30})`, "replaceOne"},
		{`db.users.findOneAndUpdate({name: "a"}, {$set: {age: 30}})`, "findOneAndUpdate"},
		{`db.users.findOneAndReplace({name: "a"}, {name: "b"})`, "findOneAndReplace"},
		{`db.users.findOneAndDelete({name: "a"})`, "findOneAndDelete"},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs := stmt.AST.(*ast.CollectionStatement)
			require.Equal(t, tt.method, cs.Method)
		})
	}
}

func TestParseCollectionAccessPatterns(t *testing.T) {
	tests := []struct {
		input        string
		collection   string
		accessMethod string
	}{
		{`db.users.find()`, "users", "dot"},
		{`db["users"].find()`, "users", "bracket"},
		{`db['user-logs'].find()`, "user-logs", "bracket"},
		{`db.getCollection("users").find()`, "users", "getCollection"},
		{`db.getCollection("my.collection").find()`, "my.collection", "getCollection"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs := stmt.AST.(*ast.CollectionStatement)
			require.Equal(t, tt.collection, cs.Collection)
			require.Equal(t, tt.accessMethod, cs.AccessMethod)
		})
	}
}

func TestParseAggregate(t *testing.T) {
	stmt := mustParse(t, `db.orders.aggregate([{$match: {status: "completed"}}, {$group: {_id: "$category", count: {$sum: 1}}}])`)
	cs := stmt.AST.(*ast.CollectionStatement)
	require.Equal(t, "aggregate", cs.Method)
	require.Len(t, cs.Args, 1)
	_, ok := cs.Args[0].(*ast.Array)
	require.True(t, ok)
}

func TestParseIndexOperations(t *testing.T) {
	tests := []struct {
		input  string
		method string
	}{
		{`db.users.createIndex({name: 1})`, "createIndex"},
		{`db.users.dropIndex("name_1")`, "dropIndex"},
		{`db.users.dropIndexes()`, "dropIndexes"},
		{`db.users.getIndexes()`, "getIndexes"},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs := stmt.AST.(*ast.CollectionStatement)
			require.Equal(t, tt.method, cs.Method)
		})
	}
}

func TestParseExplainPrefix(t *testing.T) {
	stmt := mustParse(t, `db.users.explain().find({name: "alice"})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	require.True(t, cs.Explain)
	require.Equal(t, "find", cs.Method)
}
```

- [ ] **Step 2: Implement parseDbStatement in database.go**

This is the main routing function for all `db.` prefixed statements. It handles collection access, database methods, and special forms (bulk, encryption, plan cache).

Replace `mongo/parser/database.go`:

```go
// mongo/parser/database.go
package parser

import (
	"github.com/bytebase/omni/mongo/ast"
)

// dbMethodNames maps database-level method keywords to their string names.
// These are methods called directly on `db` (not on a collection).
var dbMethodNames = map[int]string{
	kwGetCollectionNames:             "getCollectionNames",
	kwGetCollectionInfos:             "getCollectionInfos",
	kwCreateCollection:               "createCollection",
	kwDropDatabase:                   "dropDatabase",
	kwStats:                          "stats",
	kwServerStatus:                   "serverStatus",
	kwServerBuildInfo:                "serverBuildInfo",
	kwVersion:                        "version",
	kwHostInfo:                       "hostInfo",
	kwListCommands:                   "listCommands",
	kwRunCommand:                     "runCommand",
	kwAdminCommand:                   "adminCommand",
	kwGetName:                        "getName",
	kwGetMongo:                       "getMongo",
	kwGetSiblingDB:                   "getSiblingDB",
	kwAggregate:                      "aggregate",
	kwAuth:                           "auth",
	kwChangeUserPassword:             "changeUserPassword",
	kwCloneDatabase:                  "cloneDatabase",
	kwCommandHelp:                    "commandHelp",
	kwCopyDatabase:                   "copyDatabase",
	kwCreateRole:                     "createRole",
	kwCreateUser:                     "createUser",
	kwCreateView:                     "createView",
	kwCurrentOp:                      "currentOp",
	kwDropAllRoles:                   "dropAllRoles",
	kwDropAllUsers:                   "dropAllUsers",
	kwDropRole:                       "dropRole",
	kwDropUser:                       "dropUser",
	kwFsyncLock:                      "fsyncLock",
	kwFsyncUnlock:                    "fsyncUnlock",
	kwGetLogComponents:               "getLogComponents",
	kwGetProfilingLevel:              "getProfilingLevel",
	kwGetProfilingStatus:             "getProfilingStatus",
	kwGetReplicationInfo:             "getReplicationInfo",
	kwGetRole:                        "getRole",
	kwGetRoles:                       "getRoles",
	kwGetUser:                        "getUser",
	kwGetUsers:                       "getUsers",
	kwGrantPrivilegesToRole:          "grantPrivilegesToRole",
	kwGrantRolesToRole:               "grantRolesToRole",
	kwGrantRolesToUser:               "grantRolesToUser",
	kwHello:                          "hello",
	kwIsMaster:                       "isMaster",
	kwKillOp:                         "killOp",
	kwLogout:                         "logout",
	kwPrintCollectionStats:           "printCollectionStats",
	kwPrintReplicationInfo:           "printReplicationInfo",
	kwPrintSecondaryReplicationInfo:  "printSecondaryReplicationInfo",
	kwPrintShardingStatus:            "printShardingStatus",
	kwPrintSlaveReplicationInfo:      "printSlaveReplicationInfo",
	kwRevokePrivilegesFromRole:       "revokePrivilegesFromRole",
	kwRevokeRolesFromRole:            "revokeRolesFromRole",
	kwRevokeRolesFromUser:            "revokeRolesFromUser",
	kwRotateCertificates:             "rotateCertificates",
	kwSetLogLevel:                    "setLogLevel",
	kwSetProfilingLevel:              "setProfilingLevel",
	kwSetSecondaryOk:                 "setSecondaryOk",
	kwSetWriteConcern:                "setWriteConcern",
	kwShutdownServer:                 "shutdownServer",
	kwUpdateRole:                     "updateRole",
	kwUpdateUser:                     "updateUser",
	kwWatch:                          "watch",
}

// parseDbStatement parses all statements starting with "db".
// Routes to: database methods, collection operations, bulk, encryption, plan cache.
func (p *Parser) parseDbStatement() (ast.Node, error) {
	stmtStart := p.cur.Loc
	p.advance() // consume "db"

	// db must be followed by a dot or bracket
	switch p.cur.Type {
	case tokDot:
		p.advance() // consume "."
	case tokLBracket:
		// db["collection"] access
		return p.parseDbBracketAccess(stmtStart)
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// After "db.", check what follows.
	// Could be: database method, getCollection, getMongo (for encryption/connection), or collection name.

	// Check for db.getMongo() — could lead to encryption or connection chain
	if p.cur.Type == kwGetMongo {
		return p.parseDbGetMongo(stmtStart)
	}

	// Check for db.getCollection("name")
	if p.cur.Type == kwGetCollection {
		return p.parseDbGetCollection(stmtStart)
	}

	// Check if the identifier is a known database method
	if name, ok := dbMethodNames[p.cur.Type]; ok {
		// But we need to check if it's followed by "(" — if so it's a db method.
		// If followed by "." it's a collection name that happens to match a keyword.
		if p.peekNext().Type == tokLParen {
			return p.parseDatabaseMethodCall(name, stmtStart)
		}
	}

	// Otherwise it's a collection name (could be any identifier or keyword)
	if !p.isIdentLike(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	collName := p.cur.Str
	collLoc := ast.Loc{Start: p.cur.Loc, End: p.cur.End}
	p.advance() // consume collection name

	// Must be followed by "."
	if _, err := p.expect(tokDot); err != nil {
		return nil, err
	}

	// Check what method follows — special forms first
	switch p.cur.Type {
	case kwInitializeOrderedBulkOp, kwInitializeUnorderedBulkOp:
		return p.parseBulkStatement(collName, collLoc, "dot", stmtStart)
	case kwGetPlanCache:
		return p.parsePlanCacheStatement(collName, collLoc, "dot", stmtStart)
	default:
		return p.parseCollectionStatement(collName, collLoc, "dot", stmtStart)
	}
}

// parseDbBracketAccess handles db["collection"] or db['collection'] access.
func (p *Parser) parseDbBracketAccess(stmtStart int) (ast.Node, error) {
	p.advance() // consume "["
	if p.cur.Type != tokString {
		return nil, &ParseError{
			Message:  "expected string for collection name",
			Position: p.cur.Loc,
			Line:     p.lineColLine(p.cur.Loc),
			Column:   p.lineColCol(p.cur.Loc),
		}
	}
	collName := p.cur.Str
	collLoc := ast.Loc{Start: p.cur.Loc, End: p.cur.End}
	p.advance() // consume string
	if _, err := p.expect(tokRBracket); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokDot); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwInitializeOrderedBulkOp, kwInitializeUnorderedBulkOp:
		return p.parseBulkStatement(collName, collLoc, "bracket", stmtStart)
	case kwGetPlanCache:
		return p.parsePlanCacheStatement(collName, collLoc, "bracket", stmtStart)
	default:
		return p.parseCollectionStatement(collName, collLoc, "bracket", stmtStart)
	}
}

// parseDbGetCollection handles db.getCollection("name").method(...)
func (p *Parser) parseDbGetCollection(stmtStart int) (ast.Node, error) {
	p.advance() // consume "getCollection"
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	if p.cur.Type != tokString {
		return nil, &ParseError{
			Message:  "expected string for collection name",
			Position: p.cur.Loc,
			Line:     p.lineColLine(p.cur.Loc),
			Column:   p.lineColCol(p.cur.Loc),
		}
	}
	collName := p.cur.Str
	collLoc := ast.Loc{Start: p.cur.Loc, End: p.cur.End}
	p.advance() // consume string
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokDot); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwInitializeOrderedBulkOp, kwInitializeUnorderedBulkOp:
		return p.parseBulkStatement(collName, collLoc, "getCollection", stmtStart)
	case kwGetPlanCache:
		return p.parsePlanCacheStatement(collName, collLoc, "getCollection", stmtStart)
	default:
		return p.parseCollectionStatement(collName, collLoc, "getCollection", stmtStart)
	}
}

// parseDbGetMongo handles db.getMongo() which can lead to:
// - encryption: db.getMongo().getKeyVault().xxx() or db.getMongo().getClientEncryption().xxx()
// - connection chain: db.getMongo().getDB() etc.
// - simple db.getMongo() (database method)
func (p *Parser) parseDbGetMongo(stmtStart int) (ast.Node, error) {
	p.advance() // consume "getMongo"
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}

	// If no dot follows, it's just db.getMongo()
	if p.cur.Type != tokDot {
		return &ast.DatabaseStatement{
			Method: "getMongo",
			Args:   nil,
			Loc:    ast.Loc{Start: stmtStart, End: p.prev.End},
		}, nil
	}

	p.advance() // consume "."

	// Check for encryption patterns
	if p.cur.Type == kwGetKeyVault || p.cur.Type == kwGetClientEncryption {
		return p.parseEncryptionStatement(stmtStart)
	}

	// Otherwise it's a connection method chain on getMongo()
	return p.parseGetMongoConnectionChain(stmtStart)
}

// parseGetMongoConnectionChain handles db.getMongo().method() connection chains.
func (p *Parser) parseGetMongoConnectionChain(stmtStart int) (ast.Node, error) {
	var methods []ast.MethodCall
	for {
		if !p.isIdentLike(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		methodStart := p.cur.Loc
		methodName := p.cur.Str
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		methods = append(methods, ast.MethodCall{
			Method: methodName,
			Args:   args,
			Loc:    ast.Loc{Start: methodStart, End: p.prev.End},
		})
		if p.cur.Type != tokDot {
			break
		}
		p.advance() // consume "."
	}
	return &ast.ConnectionStatement{
		Constructor:    "getMongo",
		Args:           nil,
		ChainedMethods: methods,
		Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}

// parseDatabaseMethodCall parses db.methodName(...).
func (p *Parser) parseDatabaseMethodCall(name string, stmtStart int) (ast.Node, error) {
	p.advance() // consume method name
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &ast.DatabaseStatement{
		Method: name,
		Args:   args,
		Loc:    ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
```

- [ ] **Step 3: Implement parseCollectionStatement in collection.go**

Replace `mongo/parser/collection.go`:

```go
// mongo/parser/collection.go
package parser

import (
	"github.com/bytebase/omni/mongo/ast"
)

// parseCollectionStatement parses db.collection.method(...).cursorMethod(...)...
// The collection name, loc, and access method have already been parsed.
// The parser is positioned at the method name token.
func (p *Parser) parseCollectionStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	stmt := &ast.CollectionStatement{
		Collection:    collName,
		CollectionLoc: collLoc,
		AccessMethod:  accessMethod,
	}

	// Check for explain prefix: db.collection.explain().method()
	if p.cur.Type == kwExplain {
		stmt.Explain = true
		explainStart := p.cur.Loc
		p.advance() // consume "explain"
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		stmt.ExplainArgs = args
		_ = explainStart

		// After explain(), expect "." then the actual method
		if _, err := p.expect(tokDot); err != nil {
			return nil, err
		}
	}

	// Parse the primary collection method
	if !p.isIdentLike(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Method = p.cur.Str
	p.advance() // consume method name
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	stmt.Args = args

	// Parse optional cursor method chain: .sort(...).limit(...).skip(...)
	for p.cur.Type == tokDot {
		p.advance() // consume "."
		if !p.isIdentLike(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		cursorStart := p.cur.Loc
		cursorMethod := p.cur.Str
		p.advance() // consume method name
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		cursorArgs, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		stmt.CursorMethods = append(stmt.CursorMethods, ast.CursorMethod{
			Method: cursorMethod,
			Args:   cursorArgs,
			Loc:    ast.Loc{Start: cursorStart, End: p.prev.End},
		})
	}

	stmt.Loc = ast.Loc{Start: stmtStart, End: p.prev.End}
	return stmt, nil
}
```

- [ ] **Step 4: Run collection tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/parsertest/ -run TestParse -v`
Expected: PASS for shell, collection, document, expression tests

- [ ] **Step 5: Commit**

```bash
git add mongo/parser/database.go mongo/parser/collection.go mongo/parsertest/collection_test.go
git commit -m "feat(mongo/parser): collection and database statement parsing

Implements the main db.* routing in parseDbStatement, collection access
patterns (dot, bracket, getCollection), collection method parsing with
cursor method chains, explain prefix, and 66+ database method routing."
```

---

## Task 7: Remaining Statement Families

**Files:**
- Modify: `mongo/parser/replication.go`
- Modify: `mongo/parser/sharding.go`
- Modify: `mongo/parser/stream.go`
- Modify: `mongo/parser/connection.go`
- Modify: `mongo/parser/bulk.go`
- Modify: `mongo/parser/encryption.go`
- Modify: `mongo/parser/plancache.go`
- Modify: `mongo/parser/native.go`
- Create: `mongo/parsertest/replication_test.go`
- Create: `mongo/parsertest/sharding_test.go`
- Create: `mongo/parsertest/stream_test.go`
- Create: `mongo/parsertest/connection_test.go`
- Create: `mongo/parsertest/bulk_test.go`
- Create: `mongo/parsertest/encryption_test.go`
- Create: `mongo/parsertest/plancache_test.go`
- Create: `mongo/parsertest/native_test.go`

These are all structurally similar — each follows the pattern: consume prefix tokens, parse identifier method name, parse arguments.

- [ ] **Step 1: Implement rs.*, sh.*, sp.* parsers**

Replace `mongo/parser/replication.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parseRsStatement parses rs.method(...).
func (p *Parser) parseRsStatement() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume "rs"
	if _, err := p.expect(tokDot); err != nil {
		return nil, err
	}
	methodTok, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &ast.RsStatement{
		MethodName: methodTok.Str,
		Args:       args,
		Loc:        ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
```

Replace `mongo/parser/sharding.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

func (p *Parser) parseShStatement() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume "sh"
	if _, err := p.expect(tokDot); err != nil {
		return nil, err
	}
	methodTok, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &ast.ShStatement{
		MethodName: methodTok.Str,
		Args:       args,
		Loc:        ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
```

Replace `mongo/parser/stream.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parseSpStatement parses sp.method(...) or sp.resource.method(...).
func (p *Parser) parseSpStatement() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume "sp"
	if _, err := p.expect(tokDot); err != nil {
		return nil, err
	}
	firstTok, err := p.expectIdent()
	if err != nil {
		return nil, err
	}

	// Check for sp.resource.method(...) pattern
	if p.cur.Type == tokDot {
		p.advance() // consume "."
		secondTok, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		return &ast.SpStatement{
			MethodName: firstTok.Str,
			SubMethod:  secondTok.Str,
			Args:       args,
			Loc:        ast.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// Simple sp.method(...) pattern
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &ast.SpStatement{
		MethodName: firstTok.Str,
		Args:       args,
		Loc:        ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
```

- [ ] **Step 2: Implement connection, bulk, encryption, plan cache, native parsers**

Replace `mongo/parser/connection.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parseConnectionStatement parses Mongo(...) or connect(...) with optional method chains.
func (p *Parser) parseConnectionStatement() (ast.Node, error) {
	start := p.cur.Loc
	constructor := p.cur.Str
	p.advance() // consume "Mongo" or "connect"
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	var methods []ast.MethodCall
	for p.cur.Type == tokDot {
		p.advance() // consume "."
		if !p.isIdentLike(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		methodStart := p.cur.Loc
		methodName := p.cur.Str
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		methodArgs, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		methods = append(methods, ast.MethodCall{
			Method: methodName,
			Args:   methodArgs,
			Loc:    ast.Loc{Start: methodStart, End: p.prev.End},
		})
	}

	return &ast.ConnectionStatement{
		Constructor:    constructor,
		Args:           args,
		ChainedMethods: methods,
		Loc:            ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
```

Replace `mongo/parser/bulk.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parseBulkStatement parses db.collection.initializeOrderedBulkOp()... chains.
func (p *Parser) parseBulkStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	ordered := p.cur.Type == kwInitializeOrderedBulkOp
	p.advance() // consume init method
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}

	var ops []ast.BulkOperation
	for p.cur.Type == tokDot {
		p.advance() // consume "."
		if !p.isIdentLike(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		opStart := p.cur.Loc
		opMethod := p.cur.Str
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		ops = append(ops, ast.BulkOperation{
			Method: opMethod,
			Args:   args,
			Loc:    ast.Loc{Start: opStart, End: p.prev.End},
		})
	}

	return &ast.BulkStatement{
		Collection:   collName,
		AccessMethod: accessMethod,
		Ordered:      ordered,
		Operations:   ops,
		Loc:          ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
```

Replace `mongo/parser/encryption.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parseEncryptionStatement parses db.getMongo().getKeyVault()... or db.getMongo().getClientEncryption()...
// Parser is positioned at getKeyVault or getClientEncryption token.
func (p *Parser) parseEncryptionStatement(stmtStart int) (ast.Node, error) {
	var target string
	if p.cur.Type == kwGetKeyVault {
		target = "keyVault"
	} else {
		target = "clientEncryption"
	}
	p.advance() // consume getKeyVault/getClientEncryption
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}

	var methods []ast.MethodCall
	for p.cur.Type == tokDot {
		p.advance() // consume "."
		if !p.isIdentLike(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		methodStart := p.cur.Loc
		methodName := p.cur.Str
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		methods = append(methods, ast.MethodCall{
			Method: methodName,
			Args:   args,
			Loc:    ast.Loc{Start: methodStart, End: p.prev.End},
		})
	}

	return &ast.EncryptionStatement{
		Target:         target,
		ChainedMethods: methods,
		Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
```

Replace `mongo/parser/plancache.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parsePlanCacheStatement parses db.collection.getPlanCache().method()... chains.
func (p *Parser) parsePlanCacheStatement(collName string, collLoc ast.Loc, accessMethod string, stmtStart int) (ast.Node, error) {
	p.advance() // consume "getPlanCache"
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}

	var methods []ast.MethodCall
	for p.cur.Type == tokDot {
		p.advance() // consume "."
		if !p.isIdentLike(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		methodStart := p.cur.Loc
		methodName := p.cur.Str
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return nil, err
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		methods = append(methods, ast.MethodCall{
			Method: methodName,
			Args:   args,
			Loc:    ast.Loc{Start: methodStart, End: p.prev.End},
		})
	}

	return &ast.PlanCacheStatement{
		Collection:     collName,
		AccessMethod:   accessMethod,
		ChainedMethods: methods,
		Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
```

Replace `mongo/parser/native.go`:
```go
package parser

import "github.com/bytebase/omni/mongo/ast"

// parseNativeFunctionCall parses top-level function calls: identifier(args...).
func (p *Parser) parseNativeFunctionCall() (ast.Node, error) {
	start := p.cur.Loc
	name := p.cur.Str
	p.advance() // consume function name
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &ast.NativeFunctionCall{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
```

- [ ] **Step 3: Write tests for each family**

Create test files. Each test file follows the same pattern — parse input, assert AST type and key fields. Reference the corresponding ANTLR example files for test cases.

`mongo/parsertest/replication_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseRsStatements(t *testing.T) {
	tests := []struct {
		input  string
		method string
	}{
		{"rs.status()", "status"},
		{"rs.initiate()", "initiate"},
		{`rs.add("mongodb1.example.com:27017")`, "add"},
		{`rs.remove("mongodb1.example.com:27017")`, "remove"},
		{"rs.freeze(120)", "freeze"},
		{"rs.conf()", "conf"},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			rs := stmt.AST.(*ast.RsStatement)
			require.Equal(t, tt.method, rs.MethodName)
		})
	}
}
```

`mongo/parsertest/sharding_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseShStatements(t *testing.T) {
	tests := []struct {
		input  string
		method string
	}{
		{"sh.status()", "status"},
		{`sh.enableSharding("mydb")`, "enableSharding"},
		{`sh.shardCollection("mydb.users", {_id: 1})`, "shardCollection"},
		{`sh.addShard("rs1/mongodb1.example.com:27017")`, "addShard"},
		{"sh.isBalancerRunning()", "isBalancerRunning"},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			sh := stmt.AST.(*ast.ShStatement)
			require.Equal(t, tt.method, sh.MethodName)
		})
	}
}
```

`mongo/parsertest/connection_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseConnectionStatements(t *testing.T) {
	tests := []struct {
		input       string
		constructor string
		chainLen    int
	}{
		{`Mongo("mongodb://localhost:27017")`, "Mongo", 0},
		{`Mongo()`, "Mongo", 0},
		{`connect("mongodb://localhost/mydb")`, "connect", 0},
		{`Mongo("mongodb://localhost").getDB("test")`, "Mongo", 1},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			cs := stmt.AST.(*ast.ConnectionStatement)
			require.Equal(t, tt.constructor, cs.Constructor)
			require.Len(t, cs.ChainedMethods, tt.chainLen)
		})
	}
}
```

`mongo/parsertest/bulk_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseBulkStatements(t *testing.T) {
	stmt := mustParse(t, `db.users.initializeOrderedBulkOp().insert({name: "alice"}).find({status: "inactive"}).update({$set: {status: "archived"}}).execute()`)
	bs := stmt.AST.(*ast.BulkStatement)
	require.Equal(t, "users", bs.Collection)
	require.True(t, bs.Ordered)
	require.Len(t, bs.Operations, 4) // insert, find, update, execute
	require.Equal(t, "insert", bs.Operations[0].Method)
	require.Equal(t, "find", bs.Operations[1].Method)
	require.Equal(t, "update", bs.Operations[2].Method)
	require.Equal(t, "execute", bs.Operations[3].Method)
}

func TestParseBulkUnordered(t *testing.T) {
	stmt := mustParse(t, `db.products.initializeUnorderedBulkOp().insert({sku: "NEW001"}).execute()`)
	bs := stmt.AST.(*ast.BulkStatement)
	require.False(t, bs.Ordered)
}
```

`mongo/parsertest/encryption_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseEncryptionKeyVault(t *testing.T) {
	stmt := mustParse(t, `db.getMongo().getKeyVault().createKey("local")`)
	es := stmt.AST.(*ast.EncryptionStatement)
	require.Equal(t, "keyVault", es.Target)
	require.Len(t, es.ChainedMethods, 1)
	require.Equal(t, "createKey", es.ChainedMethods[0].Method)
}

func TestParseEncryptionClientEncryption(t *testing.T) {
	stmt := mustParse(t, `db.getMongo().getClientEncryption().encrypt(UUID("key"), "value", "AEAD_AES_256_CBC_HMAC_SHA_512-Random")`)
	es := stmt.AST.(*ast.EncryptionStatement)
	require.Equal(t, "clientEncryption", es.Target)
	require.Len(t, es.ChainedMethods, 1)
	require.Equal(t, "encrypt", es.ChainedMethods[0].Method)
}
```

`mongo/parsertest/plancache_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParsePlanCache(t *testing.T) {
	stmt := mustParse(t, `db.users.getPlanCache().clear()`)
	pc := stmt.AST.(*ast.PlanCacheStatement)
	require.Equal(t, "users", pc.Collection)
	require.Len(t, pc.ChainedMethods, 1)
	require.Equal(t, "clear", pc.ChainedMethods[0].Method)
}

func TestParsePlanCacheList(t *testing.T) {
	stmt := mustParse(t, `db.orders.getPlanCache().list()`)
	pc := stmt.AST.(*ast.PlanCacheStatement)
	require.Equal(t, "list", pc.ChainedMethods[0].Method)
}
```

`mongo/parsertest/stream_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseSpStatements(t *testing.T) {
	stmt := mustParse(t, `sp.createStreamProcessor("myProcessor", [{$match: {status: "active"}}])`)
	sp := stmt.AST.(*ast.SpStatement)
	require.Equal(t, "createStreamProcessor", sp.MethodName)
	require.Empty(t, sp.SubMethod)
}

func TestParseSpSubMethod(t *testing.T) {
	stmt := mustParse(t, `sp.myProcessor.start()`)
	sp := stmt.AST.(*ast.SpStatement)
	require.Equal(t, "myProcessor", sp.MethodName)
	require.Equal(t, "start", sp.SubMethod)
}
```

`mongo/parsertest/native_test.go`:
```go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseNativeFunctions(t *testing.T) {
	tests := []struct {
		input string
		name  string
	}{
		{"sleep(1000)", "sleep"},
		{`load("script.js")`, "load"},
		{"pwd()", "pwd"},
		{"quit()", "quit"},
		{`cat("/etc/hosts")`, "cat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := mustParse(t, tt.input)
			nf := stmt.AST.(*ast.NativeFunctionCall)
			require.Equal(t, tt.name, nf.Name)
		})
	}
}
```

- [ ] **Step 4: Run all tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/... -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add mongo/parser/replication.go mongo/parser/sharding.go mongo/parser/stream.go mongo/parser/connection.go mongo/parser/bulk.go mongo/parser/encryption.go mongo/parser/plancache.go mongo/parser/native.go mongo/parsertest/
git commit -m "feat(mongo/parser): all statement family parsers

Implements rs.*, sh.*, sp.*, Mongo()/connect(), bulk operations,
encryption chains, plan cache chains, and native function calls.
Full coverage of all 9 ANTLR grammar statement families."
```

---

## Task 8: Port ANTLR Test Examples

**Files:**
- Create: `mongo/parsertest/antlr_compat_test.go`
- Create: `mongo/parsertest/comment_test.go`
- Create: `mongo/parsertest/unicode_test.go`

Port all 309 ANTLR test examples. Each example file contains multiple statements that should all parse without error.

- [ ] **Step 1: Create the ANTLR compatibility test**

This test reads each ANTLR example file and verifies every statement in it parses without error. The example files live at `/Users/h3n4l/OpenSource/parser/mongodb/examples/`.

```go
// mongo/parsertest/antlr_compat_test.go
package parsertest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/mongo"
	"github.com/stretchr/testify/require"
)

// antlrExamplesDir is the path to the ANTLR example files.
// Set via ANTLR_EXAMPLES_DIR env var, or defaults to the known location.
func antlrExamplesDir() string {
	if dir := os.Getenv("ANTLR_EXAMPLES_DIR"); dir != "" {
		return dir
	}
	// Default: sibling repo location
	return filepath.Join(os.Getenv("HOME"), "OpenSource", "parser", "mongodb", "examples")
}

func TestANTLRExampleFiles(t *testing.T) {
	dir := antlrExamplesDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("ANTLR examples directory not found: %s (set ANTLR_EXAMPLES_DIR)", dir)
	}

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".js" {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			require.NoError(t, err)

			stmts, err := mongo.Parse(string(data))
			require.NoError(t, err, "failed to parse %s", entry.Name())
			require.NotEmpty(t, stmts, "expected at least one statement in %s", entry.Name())

			// Verify all statements have valid position tracking
			for i, stmt := range stmts {
				require.NotNil(t, stmt.AST, "statement %d in %s has nil AST", i, entry.Name())
				require.True(t, stmt.ByteStart >= 0, "statement %d has negative ByteStart", i)
				require.True(t, stmt.ByteEnd > stmt.ByteStart, "statement %d has invalid ByteEnd", i)
				require.True(t, stmt.Start.Line >= 1, "statement %d has invalid start line", i)
				require.True(t, stmt.Start.Column >= 1, "statement %d has invalid start column", i)
			}
		})
	}
}
```

- [ ] **Step 2: Write comment and unicode tests**

```go
// mongo/parsertest/comment_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo"
	"github.com/stretchr/testify/require"
)

func TestParseWithComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		count int
	}{
		{"line comment", "// comment\ndb.users.find()", 1},
		{"inline comment", `db.users.find({name: "alice"}) // inline`, 1},
		{"block comment", "/* comment */\ndb.users.find()", 1},
		{"comment inside document", `db.users.find({/* comment */ name: "test"})`, 1},
		{"multiple statements with comments", "show dbs // list\nshow collections /* list */\ndb.users.find()", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := mongo.Parse(tt.input)
			require.NoError(t, err)
			require.Len(t, stmts, tt.count)
		})
	}
}
```

```go
// mongo/parsertest/unicode_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestParseUnicodeStrings(t *testing.T) {
	tests := []string{
		`db.posts.find({reaction: "👍"})`,
		`db.users.find({name: "नमस्ते"})`,
		`db.users.find({name: "你好世界"})`,
		`db.users.find({greeting: "こんにちは"})`,
		`db.users.find({name: "안녕하세요"})`,
		`db.users.find({name: "Привет мир"})`,
		`db.users.find({name: "José García"})`,
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			stmt := mustParse(t, input)
			_, ok := stmt.AST.(*ast.CollectionStatement)
			require.True(t, ok)
		})
	}
}

func TestParseBracketAccessUnicode(t *testing.T) {
	stmt := mustParse(t, `db["文章"].find({})`)
	cs := stmt.AST.(*ast.CollectionStatement)
	require.Equal(t, "文章", cs.Collection)
	require.Equal(t, "bracket", cs.AccessMethod)
}
```

- [ ] **Step 3: Run all tests including ANTLR compatibility**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/... -v -count=1`
Expected: ALL PASS (or note specific ANTLR example failures to fix)

- [ ] **Step 4: Fix any ANTLR example failures**

If any ANTLR example files fail to parse, fix the parser to handle those cases. Common issues to watch for:
- Keywords used as collection names (e.g., `db.update.find()`)
- The `aggregate-js-functions.js` example (may contain JS function syntax not in our scope)
- Edge cases in document/regex parsing

Iterate until all 309 examples pass.

- [ ] **Step 5: Commit**

```bash
git add mongo/parsertest/antlr_compat_test.go mongo/parsertest/comment_test.go mongo/parsertest/unicode_test.go
git commit -m "test(mongo/parser): port all 309 ANTLR example tests

Adds ANTLR compatibility test that reads all example .js files from
the reference parser and verifies they parse successfully with correct
position tracking. Also adds comment and unicode test coverage."
```

---

## Task 9: Error Cases and Position Tracking Tests

**Files:**
- Create: `mongo/parsertest/error_test.go`
- Create: `mongo/parsertest/position_test.go`

- [ ] **Step 1: Write error case tests**

```go
// mongo/parsertest/error_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/parser"
	"github.com/stretchr/testify/require"
)

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty method", "db.users."},
		{"missing parens", "db.users.find"},
		{"unclosed brace", `db.users.find({name: "alice")`},
		{"unclosed bracket", `db.users.find([1, 2, 3)`},
		{"unknown top-level", "foobar"},
		{"invalid show target", "show tables"},
		{"new keyword", `db.users.find({_id: new ObjectId()})`},
		{"missing collection", "db..find()"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mustFail(t, tt.input)
			pe, ok := err.(*parser.ParseError)
			require.True(t, ok, "expected ParseError, got %T", err)
			require.True(t, pe.Line > 0, "expected line > 0")
			require.True(t, pe.Column > 0, "expected column > 0")
		})
	}
}

func TestNewKeywordError(t *testing.T) {
	err := mustFail(t, `db.users.find({_id: new ObjectId("abc")})`)
	pe := err.(*parser.ParseError)
	require.Contains(t, pe.Message, "'new' keyword is not supported")
}
```

- [ ] **Step 2: Write position tracking tests**

```go
// mongo/parsertest/position_test.go
package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/ast"
	"github.com/stretchr/testify/require"
)

func TestPositionTracking(t *testing.T) {
	input := `db.users.find({name: "alice"})`
	stmt := mustParse(t, input)
	require.Equal(t, 0, stmt.ByteStart)
	require.Equal(t, len(input), stmt.ByteEnd)
	require.Equal(t, mongo.Position{Line: 1, Column: 1}, stmt.Start)
}

func TestPositionMultipleStatements(t *testing.T) {
	input := "show dbs\nshow collections"
	stmts := mustParseN(t, input, 2)
	require.Equal(t, 0, stmts[0].ByteStart)
	require.Equal(t, 8, stmts[0].ByteEnd) // "show dbs"
	require.Equal(t, 1, stmts[0].Start.Line)
	require.Equal(t, 9, stmts[1].ByteStart) // "show collections" starts after \n
	require.Equal(t, 2, stmts[1].Start.Line)
}

func TestPositionDocumentNodes(t *testing.T) {
	input := `db.c.find({k: "v"})`
	stmt := mustParse(t, input)
	cs := stmt.AST.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	// Document starts at { (byte 10)
	require.Equal(t, 10, doc.Loc.Start)
	// Document ends after } (byte 19)
	require.Equal(t, 19, doc.Loc.End)
}
```

- [ ] **Step 3: Run all tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/... -v -count=1`
Expected: ALL PASS

- [ ] **Step 4: Commit**

```bash
git add mongo/parsertest/error_test.go mongo/parsertest/position_test.go
git commit -m "test(mongo/parser): error cases and position tracking tests

Verifies ParseError includes line/column for all error cases,
'new' keyword gives helpful error message, and byte offsets are
correct for statements, documents, and multi-statement input."
```

---

## Task 10: Final Verification and Cleanup

**Files:**
- All `mongo/` files

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./mongo/... -v -count=1`
Expected: ALL PASS

- [ ] **Step 2: Run go vet**

Run: `cd /Users/h3n4l/OpenSource/omni && go vet ./mongo/...`
Expected: No issues

- [ ] **Step 3: Verify the entire repo still builds and passes**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./... && go test ./... -count=1`
Expected: ALL PASS (existing tests unaffected)

- [ ] **Step 4: Final commit if any cleanup was needed**

```bash
git add -A mongo/
git commit -m "chore(mongo): final cleanup and verification"
```
