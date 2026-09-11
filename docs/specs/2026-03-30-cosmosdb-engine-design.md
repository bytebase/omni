# CosmosDB Engine for Omni

## Overview

Add a `cosmosdb/` engine to the Omni SQL toolchain. Phase 1 delivers a pure-Go
hand-written parser for Azure Cosmos DB's NoSQL SQL API, ported from the existing
ANTLR4 grammar at `github.com/bytebase/parser/cosmosdb`. Phase 2 (future) adds
query span / semantic analysis.

CosmosDB's SQL API is SELECT-only (no DDL/DML — mutations go through the SDK),
targeting a document store with JSON documents, containers, and partition keys.
This makes the parser scope significantly smaller than PostgreSQL or MySQL.

### Prior Work

- **`github.com/bytebase/parser/cosmosdb`** — ANTLR4-based parser with complete
  grammar (`CosmosDBLexer.g4`, `CosmosDBParser.g4`) and 37 example SQL files.
  Covers SELECT with TOP/DISTINCT/VALUE, FROM with JOIN/iteration/subqueries,
  WHERE, GROUP BY, HAVING, ORDER BY, OFFSET/LIMIT, all expression types, UDFs,
  array/object construction, parameter references.
- **`github.com/bytebase/bytebase`** — CosmosDB features including query span
  (column-level data lineage analysis).

### Design Decisions

1. **Pure Go, zero runtime dependencies** — matches Omni convention. No ANTLR4
   runtime. Hand-written recursive descent parser.
2. **AST follows Omni conventions** — same Node/ExprNode/TableExpr/StmtNode
   interfaces, Loc struct for position tracking, RawStmt wrapping.
3. **Full grammar in one pass** — CosmosDB SQL is small enough (~31 AST node
   types vs PostgreSQL's 210+) to port completely in Phase 1.
4. **Public API matches PostgreSQL** — `Parse(sql) -> []Statement` with byte
   offsets and line:column positions.
5. **Test corpus from ANTLR4 repo** — 37 example SQL files from
   `bytebase/parser/cosmosdb/examples/` used directly.

---

## Directory Structure

```
cosmosdb/
├── parse.go                  # Public API: Parse(sql) -> []Statement
├── parse_test.go             # Public API tests (positions, text slicing)
├── ast/
│   ├── node.go               # Interfaces (Node, ExprNode, TableExpr, StmtNode),
│   │                         # Loc struct, literal nodes
│   ├── parsenodes.go         # All CosmosDB-specific AST node types (~31 types)
│   └── outfuncs.go           # Deterministic AST text dump for golden tests
├── parser/
│   ├── lexer.go              # Hand-written tokenizer
│   ├── parser.go             # Entry point, helpers (advance/peek/match/expect)
│   ├── select.go             # SELECT clause, ORDER BY, GROUP BY, HAVING, OFFSET/LIMIT
│   ├── from.go               # FROM, JOIN, container expressions
│   ├── expr.go               # Pratt parser for scalar expressions
│   ├── function.go           # Built-in function calls + UDF calls
│   ├── parser_test.go        # Tests using 37 example files
│   └── testdata/             # Example SQL files + golden AST output
│       ├── *.sql             # Ported from bytebase/parser/cosmosdb/examples/
│       └── golden/
│           └── *.txt         # Expected AST text output
```

---

## Public API

File: `cosmosdb/parse.go`

```go
package cosmosdb

type Statement struct {
    Text      string       // original SQL text
    AST       ast.Node     // inner statement node (*ast.SelectStmt)
    ByteStart int          // inclusive start byte offset
    ByteEnd   int          // exclusive end byte offset
    Start     Position     // start line:column (1-based)
    End       Position     // end line:column (1-based)
}

type Position struct {
    Line   int  // 1-based line number
    Column int  // 1-based column in bytes
}

func Parse(sql string) ([]Statement, error)
```

CosmosDB SQL is always a single SELECT statement, so `Parse` typically returns
one `Statement`. The slice return type maintains consistency with the Omni
pattern.

---

## AST Node Types

### Core Interfaces

File: `cosmosdb/ast/node.go`

```go
type Node interface { nodeTag() }
type ExprNode interface { Node; exprNode() }
type TableExpr interface { Node; tableExpr() }
type StmtNode interface { Node; stmtNode() }

type Loc struct {
    Start int  // inclusive byte offset, -1 = unknown
    End   int  // exclusive byte offset, -1 = unknown
}
```

### Literal Nodes

Defined in `ast/node.go`, all implement `ExprNode`:

| Node | Fields | Grammar rule |
|------|--------|-------------|
| `StringLit` | Val string, Loc | `string_constant` |
| `NumberLit` | Val string, Loc | `decimal_literal`, `hexadecimal_literal` |
| `BoolLit` | Val bool, Loc | `boolean_constant` |
| `NullLit` | Loc | `null_constant` |
| `UndefinedLit` | Loc | `undefined_constant` |
| `InfinityLit` | Loc | `INFINITY_SYMBOL` |
| `NanLit` | Loc | `NAN_SYMBOL` |
| `List` | Items []Node | Generic list container |

`NumberLit.Val` stores the raw text to preserve the original representation
(integer, float, hex, scientific notation).

### Statement & Clause Nodes

File: `cosmosdb/ast/parsenodes.go`

| Node | Key Fields | Grammar rule |
|------|-----------|-------------|
| `RawStmt` | Stmt StmtNode, StmtLocation int, StmtLen int | Omni convention wrapper |
| `SelectStmt` | Top *int, Distinct bool, Value bool, Star bool, Targets []*TargetEntry, From TableExpr, Joins []*JoinExpr, Where ExprNode, GroupBy []ExprNode, Having ExprNode, OrderBy []*SortExpr, OffsetLimit *OffsetLimitClause, Loc | `select` |
| `TargetEntry` | Expr ExprNode, Alias *string, Loc | `object_property` |
| `SortExpr` | Expr ExprNode, Desc bool, Rank bool, Loc | `sort_expression` |
| `OffsetLimitClause` | Offset ExprNode, Limit ExprNode, Loc | `offset_limit_clause` |
| `JoinExpr` | Source TableExpr, Loc | `join_clause` |

### Table Expression Nodes

Implement `TableExpr`. `DotAccessExpr` and `BracketAccessExpr` implement both
`TableExpr` and `ExprNode` since they appear in both FROM and expression
contexts.

| Node | Key Fields | Grammar rule |
|------|-----------|-------------|
| `ContainerRef` | Name string, Root bool, Loc | `container_name`, `ROOT` |
| `AliasedTableExpr` | Source TableExpr, Alias string, Loc | `from_source` alt1 |
| `ArrayIterationExpr` | Alias string, Source TableExpr, Loc | `from_source` alt2 (`ident IN expr`) |
| `SubqueryExpr` | Select *SelectStmt, Loc | `container_expression` alt5 |

### Expression Nodes

All implement `ExprNode`:

| Node | Key Fields | Grammar rule |
|------|-----------|-------------|
| `ColumnRef` | Name string, Loc | `input_alias` |
| `DotAccessExpr` | Expr ExprNode, Property string, Loc | `scalar_expression` alt11, `container_expression` alt3 |
| `BracketAccessExpr` | Expr ExprNode, Index ExprNode, Loc | `scalar_expression` alt12, `container_expression` alt4 |
| `BinaryExpr` | Op string, Left ExprNode, Right ExprNode, Loc | alts 15-22, 26-28 |
| `UnaryExpr` | Op string, Operand ExprNode, Loc | alts 13-14 |
| `TernaryExpr` | Cond ExprNode, Then ExprNode, Else ExprNode, Loc | alt29 (`? :`) |
| `InExpr` | Expr ExprNode, List []ExprNode, Not bool, Loc | alt23 |
| `BetweenExpr` | Expr ExprNode, Low ExprNode, High ExprNode, Not bool, Loc | alt24 |
| `LikeExpr` | Expr ExprNode, Pattern ExprNode, Escape ExprNode, Not bool, Loc | alt25 |
| `FuncCall` | Name string, Args []ExprNode, Star bool, Loc | `builtin_function_expression` |
| `UDFCall` | Name string, Args []ExprNode, Loc | `udf_scalar_function_expression` |
| `ExistsExpr` | Select *SelectStmt, Loc | alt9 |
| `ArrayExpr` | Select *SelectStmt, Loc | alt10 (`ARRAY(SELECT ...)`) |
| `CreateArrayExpr` | Elements []ExprNode, Loc | `create_array_expression` |
| `CreateObjectExpr` | Fields []*ObjectFieldPair, Loc | `create_object_expression` |
| `ObjectFieldPair` | Key string, Value ExprNode, Loc | `object_field_pair` |
| `ParamRef` | Name string, Loc | `parameter_name` (`@param`) |
| `SubLink` | Select *SelectStmt, Loc | alt8 (`(SELECT ...)` in expr position) |

**Total: ~31 node types.**

---

## Lexer

File: `cosmosdb/parser/lexer.go`

```go
type Token struct {
    Type int
    Str  string
    Loc  int    // byte offset in source
}

type Lexer struct {
    input string
    pos   int    // current read position in input (next byte to consume)
    start int    // start position of the token currently being scanned
    Err   error
}

func NewLexer(input string) *Lexer
func (l *Lexer) NextToken() Token
```

### Token Types

**Literals and identifiers (iota + 1000):**
`tokICONST` (integer), `tokFCONST` (float), `tokHCONST` (hex),
`tokSCONST` (single-quoted string), `tokDCONST` (double-quoted string),
`tokIDENT` (identifier), `tokPARAM` (`@param`).

**Operators and punctuation:**
`tokDOT`, `tokCOMMA`, `tokCOLON`, `tokLPAREN`, `tokRPAREN`, `tokLBRACK`,
`tokRBRACK`, `tokLBRACE`, `tokRBRACE`, `tokSTAR`, `tokPLUS`, `tokMINUS`,
`tokDIV`, `tokMOD`, `tokEQ`, `tokNE` (`!=`), `tokNE2` (`<>`), `tokLT`,
`tokLE`, `tokGT`, `tokGE`, `tokBITAND`, `tokBITOR`, `tokBITXOR`,
`tokBITNOT`, `tokLSHIFT`, `tokRSHIFT`, `tokURSHIFT` (`>>>`),
`tokCONCAT` (`||`), `tokCOALESCE` (`??`), `tokQUESTION` (`?`).

**Keywords (28, case-insensitive):**
`tokSELECT`, `tokFROM`, `tokWHERE`, `tokAND`, `tokOR`, `tokNOT`, `tokIN`,
`tokBETWEEN`, `tokLIKE`, `tokESCAPE`, `tokAS`, `tokJOIN`, `tokTOP`,
`tokDISTINCT`, `tokVALUE`, `tokORDER`, `tokBY`, `tokGROUP`, `tokHAVING`,
`tokOFFSET`, `tokLIMIT`, `tokASC`, `tokDESC`, `tokEXISTS`, `tokTRUE`,
`tokFALSE`, `tokNULL`, `tokUNDEFINED`, `tokUDF`, `tokARRAY`, `tokROOT`,
`tokRANK`.

### Special Cases

- **Case-insensitive keywords**: Identifiers are lowercased before keyword
  lookup. Keywords always emit keyword tokens; the parser accepts them in
  identifier positions.
- **Case-sensitive `Infinity` and `NaN`**: Matched before lowercasing. Emitted
  as distinct `tokINFINITY` and `tokNAN` token types.
- **Line comments**: `--` skips to EOL.
- **String escapes**: `\b`, `\t`, `\n`, `\r`, `\f`, `\"`, `\'`, `\\`, `\/`,
  `\uXXXX`.
- **Multi-character operators**: `??`, `||`, `!=`, `<>`, `<=`, `>=`, `<<`,
  `>>`, `>>>` — resolved by longest-match lookahead.

---

## Parser

### Structure

```go
type Parser struct {
    lexer   *Lexer
    cur     Token    // current token (already consumed)
    prev    Token    // previous token (for error context)
    nextBuf Token    // buffered lookahead token
    hasNext bool     // whether nextBuf is valid
}

func Parse(sql string) (*ast.List, error)
```

### File Responsibilities

| File | Functions |
|------|-----------|
| `parser.go` | `Parse()`, `advance()`, `peek()`, `peekNext()`, `match()`, `expect()`, `parseIdentifier()`, `parsePropertyName()`, top-level `parseSelect()` |
| `select.go` | `parseSelectClause()`, `parseTargetEntry()`, `parseGroupBy()`, `parseHaving()`, `parseOrderBy()`, `parseSortExpr()`, `parseOffsetLimit()` |
| `from.go` | `parseFromClause()`, `parseFromSource()`, `parseContainerExpr()`, `parseJoinClause()` |
| `expr.go` | `parseExpr(precedence)`, prefix parsers (constants, identifiers, unary ops, parenthesized exprs, subqueries, EXISTS, ARRAY, object/array literals, params), infix parsers (binary ops, dot access, bracket access, IN, BETWEEN, LIKE, ternary) |
| `function.go` | `parseFuncCall()`, `parseUDFCall()` |

### Expression Precedence (Pratt Parser)

From lowest to highest binding power:

| Level | Operators | Associativity |
|-------|-----------|---------------|
| 1 | `? :` (ternary) | right |
| 2 | `??` (coalesce) | left |
| 3 | `OR` | left |
| 4 | `AND` | left |
| 5 | `NOT IN`, `NOT BETWEEN`, `NOT LIKE`, `IN`, `BETWEEN`, `LIKE` | non-assoc |
| 6 | `=`, `!=`, `<>`, `<`, `<=`, `>`, `>=` | left |
| 7 | `\|\|` (string concat) | left |
| 8 | `\|` (bitwise OR) | left |
| 9 | `^` (bitwise XOR) | left |
| 10 | `&` (bitwise AND) | left |
| 11 | `<<`, `>>`, `>>>` (shift) | left |
| 12 | `+`, `-` (additive) | left |
| 13 | `*`, `/`, `%` (multiplicative) | left |
| 14 | `NOT`, `~`, unary `+`, unary `-` (prefix) | prefix |
| 15 | `.`, `[]` (property access) | left |

### Keyword-as-Identifier Handling

Two helper functions handle the grammar's context-sensitive identifier rules:

- `parseIdentifier()` — accepts `tokIDENT` plus 20 keywords that are valid as
  identifiers: `IN`, `BETWEEN`, `TOP`, `VALUE`, `ORDER`, `BY`, `GROUP`,
  `OFFSET`, `LIMIT`, `ASC`, `DESC`, `EXISTS`, `LIKE`, `HAVING`, `JOIN`,
  `ESCAPE`, `ARRAY`, `ROOT`, `RANK`.
- `parsePropertyName()` — accepts everything `parseIdentifier()` does, plus 18
  additional keywords: `SELECT`, `FROM`, `WHERE`, `NOT`, `AND`, `OR`, `AS`,
  `TRUE`, `FALSE`, `NULL`, `UNDEFINED`, `UDF`, `DISTINCT`, `ARRAY`, `ROOT`,
  `ESCAPE`, `RANK`.

---

## Test Strategy

### Test Corpus

The 37 example SQL files from `github.com/bytebase/parser/cosmosdb/examples/`
are fetched and placed in `cosmosdb/parser/testdata/`. These cover:

- Basic SELECT, SELECT TOP, SELECT DISTINCT, SELECT VALUE
- FROM with container reference, ROOT, subquery, IN iteration
- JOIN on array properties
- WHERE with all operator types (comparison, logical, BETWEEN, IN, LIKE, ESCAPE)
- GROUP BY, HAVING, ORDER BY with ASC/DESC, OFFSET/LIMIT
- Built-in functions (string, math, type-check, geospatial, aggregate)
- UDF calls
- Array and object construction
- Special constants: Infinity, NaN, undefined, null
- Coalesce (`??`), ternary (`? :`)
- Parameter references (`@param`)
- Keywords as property names and identifiers
- Bracket access and dot access on nested documents
- Line comments

### Test Approach

`parser_test.go` implements:

1. Parse each `.sql` file from `testdata/`
2. Produce deterministic text output via `ast/outfuncs.go`
3. Compare against golden files in `testdata/golden/`
4. Update golden files with a `-update` flag for initial generation

`parse_test.go` tests the public API:

1. Verify `Statement.Text`, `ByteStart`, `ByteEnd`, `Start`, `End` positions
2. Verify round-trip: `statement.Text` can be re-parsed to produce identical AST

No oracle testing against a real CosmosDB instance in Phase 1.

---

## Phase 2 (Future): Semantic Analysis

Not in scope for Phase 1, but the intended direction:

- **Query span / column lineage**: Track output properties back to source
  container document paths through JOINs, subqueries, array iteration, UDFs.
- **Lightweight catalog**: Model containers, partition keys, indexing policies.
  No fixed schema (CosmosDB is schemaless), but container metadata is useful
  for validation.
- **Type inference**: CosmosDB expressions have runtime types (string, number,
  boolean, null, undefined, array, object). Static inference where possible.
- **Built-in function registry**: Validate function names and argument counts
  for CosmosDB's ~100 built-in scalar functions.
