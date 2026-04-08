# PartiQL `ast-core` Design Spec

**DAG node:** `ast-core` (node 1 in `docs/migration/partiql/dag.md`)
**Priority:** P0 (foundational — every parser DAG node depends on this)
**Branch:** `feat/partiql/ast-core`
**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-ast-core`
**Linear umbrella:** BYT-9000

## Goal

Create the `partiql/ast` package — the AST type taxonomy for the omni PartiQL parser. This is the **first node** in the PartiQL migration DAG. It defines the type system every subsequent parser node will populate. Per the user decision, scope is **full taxonomy up front** (DAG choice C): the AST covers the entire legacy `bytebase/parser/partiql` ANTLR grammar in one landing, including P1 features like graph patterns, LET/PIVOT/UNPIVOT, window functions, and Ion literals. Subsequent parser DAG nodes (`parser-foundation`, `parser-select`, `parser-dml`, `parser-ddl`, …, `parser-graph`) only fill in parser logic; they should not need to amend `partiql/ast/`.

## Inputs

- `docs/migration/partiql/analysis.md` — full grammar coverage of the legacy parser
- `docs/migration/partiql/dag.md` — migration DAG and node ordering
- `bytebase/parser/partiql/PartiQLParser.g4` — authoritative grammar source (line refs in this spec)
- Reference engines: `cosmosdb/ast/` (most recent NoSQL, sealed sub-interface style), `mongo/ast/` (looser flat style)

## Architecture decisions

### D1. Sealed sub-interface style (rich variant)

Approach 1 from brainstorming. Six sealed sub-interfaces enforce position discipline at every grammar boundary:

```go
// node.go

// Loc is a half-open byte range. -1 in both fields means unknown.
type Loc struct {
    Start int // inclusive
    End   int // exclusive
}

// Node is the root interface for all PartiQL parse-tree nodes.
type Node interface {
    nodeTag()
    GetLoc() Loc
}

// StmtNode marks top-level statement nodes.
type StmtNode interface {
    Node
    stmtNode()
}

// ExprNode marks scalar-position expression nodes.
type ExprNode interface {
    Node
    exprNode()
}

// TableExpr marks FROM-position nodes.
type TableExpr interface {
    Node
    tableExpr()
}

// PathStep marks a single step in a PartiQL path expression.
type PathStep interface {
    Node
    pathStep()
}

// TypeName marks PartiQL type references (used by CAST and DDL).
type TypeName interface {
    Node
    typeName()
}

// PatternNode marks GPML graph-match pattern nodes.
type PatternNode interface {
    Node
    patternNode()
}
```

**Rationale:** PartiQL's grammar has more distinct categories than cosmosdb's. Path navigation (PartiQL-unique), graph patterns (GPML), and type names are first-class enough to deserve their own marker interfaces. Compile-time discrimination prevents an entire class of parser bugs (e.g., handing a `MatchExpr` to a slot that expects `BinaryExpr`).

### D2. Multi-interface satisfaction for unified grammar productions

Three node types satisfy more than one interface, mirroring grammar productions that legitimately appear in multiple positions:

| Node | Interfaces | Why |
|------|-----------|-----|
| `PathExpr` | `ExprNode` + `TableExpr` | `FROM a.b[0]` is a legal table source; same path also valid in scalar position |
| `VarRef` | `ExprNode` + `TableExpr` | `FROM tbl` resolves through the same identifier rule as scalar `tbl` |
| `SubLink` | `ExprNode` + `TableExpr` | `(SELECT …)` is valid in both positions |

These are documented inline in `tableexprs.go` with comment links back to `exprs.go`.

### D3. Byte-offset `Loc` only

Each node embeds a `Loc` field. `Loc` carries `Start` and `End` byte offsets only — no line/column. Line/column conversion happens at the public `partiql/parse.go` boundary (mongo-style), not in the AST. Synthetic nodes (created post-parse) use `Loc{-1, -1}`.

### D4. Files split by category

```
partiql/ast/
├── node.go        # interfaces, Loc, List
├── literals.go    # all literal nodes (string, number, bool, null, missing, date/time/timestamp, ion)
├── exprs.go       # expression nodes (operators, predicates, function calls, CASE, CAST,
│                  #   subqueries, paths, variables, collection literals)
├── stmts.go       # statement nodes + DML/SELECT clause helpers
├── tableexprs.go  # FROM-position nodes (table refs, joins, unpivot)
├── types.go       # type name nodes
├── patterns.go    # graph pattern nodes (GPML / MATCH clauses)
├── outfuncs.go    # NodeToString(Node) — deterministic debug dump
└── ast_test.go    # unit tests for tag dispatch + Loc handling
```

Splitting `parsenodes.go` into 6 thematic files (`literals.go`, `exprs.go`, `stmts.go`, `tableexprs.go`, `types.go`, `patterns.go`) plus the foundation `node.go` and the utility `outfuncs.go` — vs cosmosdb's single `parsenodes.go`. The split is justified because PartiQL has ~3× the node count and the grammar categories are very distinct. Each file targets ≤ ~400 lines.

### D5. Single `TypeRef` for all type names

Type names use one flat node:

```go
type TypeRef struct {
    Name         string  // canonical uppercase: "INT", "DECIMAL", "TIME", "BAG", ...
    Args         []int   // optional precision/scale/length: DECIMAL(10,2) -> [10,2]
    WithTimeZone bool    // for TIME / TIMESTAMP
    Loc          Loc
}
```

Discrimination by `Name` rather than a Go enum. This matches the PartiQL grammar's treatment (a name token with optional modifiers) and avoids a 25-arm enum that would mostly be one-to-one with the name string.

### D6. `FuncCall` covers regular calls, aggregates, and windowed calls

Window functions, aggregates, and ordinary function calls use one node:

```go
type FuncCall struct {
    Name       string
    Args       []ExprNode
    Quantifier QuantifierKind  // NONE/DISTINCT/ALL — for aggregates
    Star       bool            // true for COUNT(*)
    Over       *WindowSpec     // non-nil for window calls (LAG/LEAD)
    Loc        Loc
}
```

This avoids three separate node types for what's structurally the same shape. Whether a particular `FuncCall` is "really" an aggregate or a window function is determined by parser context, not by node identity.

`EXTRACT`, `TRIM`, `SUBSTRING`, `CAST`, `CASE`, `COALESCE`, and `NULLIF` get their own dedicated nodes because the PartiQL grammar gives them special syntactic forms (keywords inside parens or non-comma argument syntax). `DATE_ADD` and `DATE_DIFF` use plain `FuncCall` because their syntax is just a comma-separated arg list.

### D7. No `Walker` / `Visitor`, no smart constructors

Out of scope. Go type switches are sufficient for the parser, the future analyzer, and the future completion engine. Adding a walker now would be speculative. Smart constructors (`NewBinaryExpr(...)`) add maintenance burden without benefit when struct literals are clear.

## Type taxonomy

Around **73 node types** across the 6 thematic files, plus `node.go` and `outfuncs.go`. Tagged ★ marks PartiQL-unique features.

### `literals.go` (9 types, all `ExprNode`)

| Node | Fields | Grammar |
|------|--------|---------|
| `StringLit` | `Val string` | `'…'` and `"…"` strings |
| `NumberLit` | `Val string` (raw text preserved) | integer / decimal / scientific |
| `BoolLit` | `Val bool` | `TRUE` / `FALSE` |
| `NullLit` | — | `NULL` |
| `MissingLit` ★ | — | `MISSING` |
| `DateLit` | `Val string` | `DATE 'YYYY-MM-DD'` |
| `TimeLit` | `Val string`, `Precision *int`, `WithTimeZone bool` | `TIME [(p)] [WITH TIME ZONE] '…'` |
| `TimestampLit` | `Val string`, `Precision *int`, `WithTimeZone bool` | `TIMESTAMP [(p)] [WITH TIME ZONE] '…'` |
| `IonLit` ★ | `Text string` (verbatim backtick contents) | `` ` … ` `` |

### `stmts.go` — top-level statements (`StmtNode`)

| Node | Key fields |
|------|-----------|
| `SelectStmt` | `Quantifier`, `Star`, `Value`, `Pivot *PivotProjection`, `Targets []*TargetEntry`, `From TableExpr`, `Let []*LetBinding`, `Where`, `GroupBy *GroupByClause`, `Having`, `OrderBy []*OrderByItem`, `Limit`, `Offset` |
| `SetOpStmt` | `Op SetOpKind` (UNION/INTERSECT/EXCEPT), `Quantifier QuantifierKind`, `Outer bool`, `Left StmtNode`, `Right StmtNode` |
| `ExplainStmt` | `Inner StmtNode` |
| `InsertStmt` | `Target TableExpr`, `AsAlias *string`, `Value ExprNode`, `OnConflict *OnConflict`, `Returning *ReturningClause` |
| `UpdateStmt` | `Source TableExpr`, `Sets []*SetAssignment`, `Where`, `Returning` |
| `DeleteStmt` | `Source TableExpr`, `Where`, `Returning` |
| `UpsertStmt` | (same shape as `InsertStmt`) |
| `ReplaceStmt` | (same shape as `InsertStmt`) |
| `RemoveStmt` ★ | `Path *PathExpr` |
| `CreateTableStmt` | `Name *VarRef` |
| `CreateIndexStmt` | `Table *VarRef`, `Paths []*PathExpr` |
| `DropTableStmt` | `Name *VarRef` |
| `DropIndexStmt` | `Index *VarRef`, `Table *VarRef` |
| `ExecStmt` | `Name string`, `Args []ExprNode` |

### `stmts.go` — SELECT clause helpers (regular `Node`)

| Node | Fields |
|------|--------|
| `TargetEntry` | `Expr`, `Alias *string` |
| `PivotProjection` ★ | `Value`, `At` |
| `LetBinding` ★ | `Expr`, `Alias string` |
| `GroupByClause` | `Partial bool`, `Items []*GroupByItem`, `GroupAs *string` |
| `GroupByItem` | `Expr`, `Alias *string` |
| `OrderByItem` | `Expr`, `Desc bool`, `NullsFirst bool`, `NullsExplicit bool` |

### `stmts.go` — DML helpers (regular `Node` / enums)

| Node | Fields |
|------|--------|
| `SetAssignment` | `Target *PathExpr`, `Value ExprNode` |
| `OnConflict` | `Target *OnConflictTarget`, `Action OnConflictAction`, `Where ExprNode` |
| `OnConflictTarget` | `Cols []*VarRef` _or_ `ConstraintName string` |
| `OnConflictAction` (enum) | `DoNothing`, `DoReplaceExcluded`, `DoUpdateExcluded` (matches the legacy ANTLR stub scope; spec gap noted) |
| `ReturningClause` | `Items []*ReturningItem` |
| `ReturningItem` | `Status` (`MODIFIED`/`ALL`), `Mapping` (`OLD`/`NEW`), `Star bool`, `Expr ExprNode` |

**Enums declared in `stmts.go`:** `SetOpKind`, `QuantifierKind`, `OnConflictAction`, `ReturningStatus`, `ReturningMapping`.

### `exprs.go` — operators & predicates (`ExprNode`)

| Node | Fields | Covers |
|------|--------|--------|
| `BinaryExpr` | `Op BinOp`, `Left`, `Right` | `OR`, `AND`, `\|\|`, `+`, `-`, `*`, `/`, `%`, `=`, `<>`, `<`, `>`, `<=`, `>=` |
| `UnaryExpr` | `Op UnOp`, `Operand` | `NOT`, unary `+`/`-` |
| `InExpr` | `Expr`, `List []ExprNode` _or_ `Subquery StmtNode`, `Not bool` | `expr [NOT] IN …` |
| `BetweenExpr` | `Expr`, `Low`, `High`, `Not bool` | `expr [NOT] BETWEEN low AND high` |
| `LikeExpr` | `Expr`, `Pattern`, `Escape ExprNode`, `Not bool` | `expr [NOT] LIKE pat [ESCAPE c]` |
| `IsExpr` | `Expr`, `Type IsType` (`NULL`/`MISSING`/`TRUE`/`FALSE`), `Not bool` | `IS [NOT] (NULL\|MISSING\|TRUE\|FALSE)` |

### `exprs.go` — special-form expressions (`ExprNode`)

| Node | Fields | Notes |
|------|--------|-------|
| `FuncCall` | `Name string`, `Args []ExprNode`, `Quantifier QuantifierKind`, `Star bool`, `Over *WindowSpec` | Generic function call; covers aggregates and window functions |
| `CaseExpr` | `Operand ExprNode` (nil for searched), `Whens []*CaseWhen`, `Else ExprNode` | Both `CASE WHEN` and `CASE expr WHEN` |
| `CaseWhen` | `When`, `Then` | Regular `Node` helper |
| `CastExpr` | `Kind CastKind` (`CAST`/`CAN_CAST`/`CAN_LOSSLESS_CAST`), `Expr`, `AsType TypeName` | |
| `ExtractExpr` | `Field string`, `From ExprNode` | `EXTRACT(field FROM expr)` |
| `TrimExpr` | `Spec` (`LEADING`/`TRAILING`/`BOTH`/none), `Sub ExprNode`, `From ExprNode` | `TRIM([spec] [sub] FROM target)` |
| `SubstringExpr` | `Expr`, `From`, `For ExprNode` | `SUBSTRING(expr FROM start [FOR len])` and the comma form |
| `CoalesceExpr` | `Args []ExprNode` | |
| `NullIfExpr` | `Left`, `Right` | |

### `exprs.go` — paths, variables, parameters, subqueries

| Node | Interfaces | Fields |
|------|-----------|--------|
| `PathExpr` ★ | `ExprNode` + `TableExpr` | `Root ExprNode`, `Steps []PathStep` |
| `VarRef` | `ExprNode` + `TableExpr` | `Name string`, `AtPrefixed bool`, `CaseSensitive bool` |
| `ParamRef` | `ExprNode` | — (`?`) |
| `SubLink` | `ExprNode` + `TableExpr` | `Stmt StmtNode` |

### `exprs.go` — collection literals (`ExprNode`)

| Node | Fields | Grammar |
|------|--------|---------|
| `ListLit` | `Items []ExprNode` | `[…]` |
| `BagLit` ★ | `Items []ExprNode` | `<<…>>` |
| `TupleLit` ★ | `Pairs []*TuplePair` | `{k:v, …}` |
| `TuplePair` | `Key ExprNode`, `Value ExprNode` | Regular `Node` helper |

### `exprs.go` — window spec (regular `Node`)

| Node | Fields |
|------|--------|
| `WindowSpec` | `PartitionBy []ExprNode`, `OrderBy []*OrderByItem` |

### `exprs.go` — path steps (`PathStep`)

| Node | Fields | Grammar |
|------|--------|---------|
| `DotStep` | `Field string`, `CaseSensitive bool` | `.field` / `."Field"` |
| `AllFieldsStep` | — | `.*` |
| `IndexStep` | `Index ExprNode` | `[expr]` |
| `WildcardStep` | — | `[*]` |

**Enums declared in `exprs.go`:** `BinOp`, `UnOp`, `IsType`, `CastKind`, `TrimSpec`.

### `tableexprs.go` (`TableExpr`)

| Node | Fields |
|------|--------|
| `TableRef` | `Name string`, `Schema string`, `CaseSensitive bool` |
| `AliasedSource` ★ | `Source TableExpr`, `As *string`, `At *string`, `By *string` |
| `JoinExpr` | `Kind JoinKind`, `Left TableExpr`, `Right TableExpr`, `On ExprNode` |
| `UnpivotExpr` ★ | `Source ExprNode`, `As *string`, `At *string`, `By *string` |

`PathExpr`, `VarRef`, and `SubLink` from `exprs.go` also implement `TableExpr`. Documented inline in `tableexprs.go`.

**Enum declared in `tableexprs.go`:** `JoinKind` (`CROSS`, `INNER`, `LEFT`, `RIGHT`, `FULL`, `OUTER`).

### `types.go` (`TypeName`)

| Node | Fields | Covers |
|------|--------|--------|
| `TypeRef` | `Name string`, `Args []int`, `WithTimeZone bool` | All PartiQL types: `INT`, `INTEGER`, `BIGINT`, `SMALLINT`, `BOOL`, `BOOLEAN`, `NULL`, `MISSING`, `STRING`, `SYMBOL`, `BLOB`, `CLOB`, `ANY`, `DATE`, `TIME`, `TIMESTAMP`, `REAL`, `DOUBLE PRECISION`, `DECIMAL(p,s)`, `NUMERIC(p,s)`, `FLOAT(n)`, `VARCHAR(n)`, `CHAR(n)`, `CHARACTER(n)`, `STRUCT`, `TUPLE`, `LIST`, `BAG`, `SEXP` |

### `patterns.go` ★ (graph patterns)

Graph Pattern Matching scaffolding. The marker interface (`PatternNode`), the node names, and the multi-interface relationship to `ExprNode` are stable in this spec. Precise field shapes inside `NodePattern`/`EdgePattern` may be refined when `parser-graph` (DAG node 16) is implemented by reading `PartiQLParser.g4` lines 314–382.

| Node | Implements | Notes |
|------|-----------|-------|
| `MatchExpr` | `ExprNode` | Top-level `MATCH (graph_expr, pattern, …)` |
| `GraphPattern` | `PatternNode` | Container for one pattern: optional selector + restrictor + variable + parts |
| `NodePattern` | `PatternNode` | `(var:Label WHERE …)` |
| `EdgePattern` | `PatternNode` | `-[var:Label]->`, `<-[]-`, `~[]~`, etc. |
| `PatternQuantifier` | regular `Node` | `+`, `*`, `{m,n}` |
| `PatternSelector` | regular `Node` | `ANY`, `ALL SHORTEST`, `SHORTEST k` |

**Enums declared in `patterns.go`:** `EdgeDirection` (LEFT, RIGHT, UNDIRECTED, ...), `PatternRestrictor` (TRAIL, ACYCLIC, SIMPLE), `SelectorKind` (ANY, ALL_SHORTEST, SHORTEST_K).

### `outfuncs.go`

`NodeToString(Node) string` — deterministic textual dump of any AST node, written as a giant type switch over every node type. Pattern follows `cosmosdb/ast/outfuncs.go`. Invariants:

- Stable field ordering (matches struct declaration order)
- Parens around nested expressions for unambiguous reading
- `Loc` is **not** dumped (positions are tested separately by `TestGetLoc`)
- Returns non-empty for every node, including zero values

Used by `ast_test.go` for snapshot-style assertions, by future parser golden tests, and for REPL debugging.

## Coverage cross-check

Every grammar area in `analysis.md`'s "Full Coverage Target" maps to ≥1 node type:

| analysis.md grammar area | Node type(s) |
|--------------------------|--------------|
| DDL (CREATE/DROP TABLE/INDEX) | `CreateTableStmt`, `CreateIndexStmt`, `DropTableStmt`, `DropIndexStmt` |
| DML — INSERT (legacy + RFC 0011) | `InsertStmt` |
| DML — UPDATE (incl. `FROM … SET …`) | `UpdateStmt` |
| DML — DELETE (incl. `FROM … DELETE …`) | `DeleteStmt` |
| DML — UPSERT / REPLACE / REMOVE | `UpsertStmt`, `ReplaceStmt`, `RemoveStmt` |
| ON CONFLICT (incl. legacy `EXCLUDED` stub) | `OnConflict`, `OnConflictTarget`, `OnConflictAction` enum |
| RETURNING | `ReturningClause`, `ReturningItem` |
| SELECT variants (`*`, items, `VALUE`, `PIVOT`) | `SelectStmt` (Star/Value/Pivot/Targets fields) |
| LET clause | `LetBinding` (`SelectStmt.Let` field) |
| FROM with `AS`/`AT`/`BY` | `AliasedSource` |
| UNPIVOT | `UnpivotExpr` |
| WHERE / GROUP BY / GROUP PARTIAL / GROUP AS / HAVING | `SelectStmt.Where`, `GroupByClause`, `SelectStmt.Having` |
| ORDER BY (ASC/DESC, NULLS FIRST/LAST) | `OrderByItem` |
| LIMIT / OFFSET | `SelectStmt.Limit`, `SelectStmt.Offset` |
| Joins (CROSS/INNER/LEFT/RIGHT/FULL/OUTER) | `JoinExpr`, `JoinKind` enum |
| Set ops (UNION/INTERSECT/EXCEPT, OUTER, DISTINCT/ALL) | `SetOpStmt` |
| Window functions (LAG/LEAD with OVER) | `FuncCall.Over`, `WindowSpec` |
| Aggregates (COUNT/SUM/AVG/MIN/MAX, DISTINCT/ALL, COUNT(*)) | `FuncCall.Quantifier`, `FuncCall.Star` |
| Built-ins (CAST/EXTRACT/TRIM/SUBSTRING/COALESCE/NULLIF/CASE) | `CastExpr`, `ExtractExpr`, `TrimExpr`, `SubstringExpr`, `CoalesceExpr`, `NullIfExpr`, `CaseExpr` |
| `CAST`/`CAN_CAST`/`CAN_LOSSLESS_CAST` | `CastExpr.Kind` |
| Predicates (IN/BETWEEN/LIKE/IS NULL/MISSING/TRUE/FALSE) | `InExpr`, `BetweenExpr`, `LikeExpr`, `IsExpr` |
| Operator hierarchy (OR/AND/NOT/+/-/*/etc.) | `BinaryExpr`, `UnaryExpr` |
| Path navigation (`.field`, `[expr]`, `[*]`, `.*`, chained) | `PathExpr`, `DotStep`, `AllFieldsStep`, `IndexStep`, `WildcardStep` |
| Variables (`@id`, `id`, `?`) | `VarRef`, `ParamRef` |
| Collection literals (list/bag/tuple) | `ListLit`, `BagLit`, `TupleLit`, `TuplePair` |
| Subqueries | `SubLink` |
| Type system (full PartiQL type list) | `TypeRef` |
| Date/time literals (DATE/TIME/TIMESTAMP) | `DateLit`, `TimeLit`, `TimestampLit` |
| Ion backtick literals | `IonLit` |
| Graph pattern matching (MATCH, nodes/edges/quantifiers/selectors/restrictors) | `MatchExpr`, `GraphPattern`, `NodePattern`, `EdgePattern`, `PatternQuantifier`, `PatternSelector` |
| EXEC | `ExecStmt` |
| EXPLAIN | `ExplainStmt` |

During implementation, do a second-pass verification by reading `bytebase/parser/partiql/PartiQLParser.g4` line by line and checking each production maps to a node type. Any missed productions update this spec and the file.

## Test plan

All tests live in a single `partiql/ast/ast_test.go` and use only the standard `testing` package.

### 1. Compile-time interface assertions

For every node type, declare `var _ <Interface> = (*Type)(nil)` so the file fails to compile if any node's interface set drifts:

```go
var _ ExprNode  = (*BinaryExpr)(nil)
var _ ExprNode  = (*PathExpr)(nil)
var _ TableExpr = (*PathExpr)(nil)  // dual-purpose
var _ ExprNode  = (*VarRef)(nil)
var _ TableExpr = (*VarRef)(nil)
// ... one block per category
```

This costs nothing at runtime and catches the multi-interface contracts.

### 2. `TestGetLoc` — `Loc` round-trip

Table-driven, one row per node type (~73 rows). Construct an instance with `Loc{Start: 10, End: 20}`, call `GetLoc()`, assert `{10, 20}`.

### 3. `TestNodeToString` — `outfuncs.go` smoke (golden assertions)

~30 representative cases, each pairing a hand-built node with an expected string literal. Covers one node per category, all multi-interface nodes, and at least one deeply nested case (e.g., a `SelectStmt` whose `From` is a `PathExpr` with multiple steps).

### 4. `TestNodeToString_AllNodesCovered` — reflection safety net

Walks every exported type in the package, checks it implements `Node`, and calls `NodeToString` on a zero value. Asserts no panic and a non-empty result. Catches the case where a new node type is added but never wired into `outfuncs.go`.

### Test scope

- Run with `go test ./partiql/ast/...` (no global runs, per the implementing skill's scoping rule)
- Zero external dependencies — standard library only
- Target ~100% line coverage on `outfuncs.go`

## Test corpus (project-wide forward reference)

**This section describes the test strategy for the PartiQL migration as a whole, not just `ast-core`.** `ast-core` itself only tests its type system (the four test groups above) — there is nothing parser-related to validate yet because no parser exists. The project-wide corpus policy is recorded here so it is visible from this spec and need not be re-derived for every later parser DAG node.

Per the test strategy decision (made during brainstorming for `ast-core` and applicable to all subsequent parser DAG nodes), parser nodes build their regression test corpora from **two authoritative sources**:

### 1. Legacy ANTLR test fixtures (parity check)

Every input from the legacy `bytebase/parser` and bytebase wrapper test corpora is re-run against the new omni parser to verify behavioral parity. This catches any case where the new parser accepts/rejects a statement differently from the ANTLR baseline:

- `bytebase/parser/partiql/parser_test.go` — basic ANTLR parse-tree assertions and the `simple.sql` smoke fixture
- `bytebase/backend/plugin/parser/partiql/test-data/test_split.yaml` — statement-splitting fixtures (drives `parser-foundation`'s splitter)
- `bytebase/backend/plugin/parser/partiql/test-data/test_completion.yaml` — auto-complete fixtures (drives the `completion` DAG node)
- `bytebase/backend/plugin/parser/partiql/split_test.go` and `completion_test.go` — additional Go-level test cases beyond the YAML fixtures

### 2. AWS DynamoDB PartiQL reference examples (coverage corpus)

Every PartiQL example published under `docs.aws.amazon.com/amazondynamodb/latest/developerguide/ql-reference.html` and its sub-pages, recursively. **Already crawled and assembled** during this brainstorming session by an explicit sub-agent task: 63 examples across 19 pages, covering SELECT (with key/non-key/IN/BETWEEN/OR/ORDER BY/document paths/indexes), INSERT (legacy + RFC 0011 with tuple value literals), UPDATE (multi-SET, `list_append`, REMOVE, `set_add`, RETURNING ALL OLD/NEW), DELETE (with RETURNING), and all six built-in functions (`size`, `exists`, `attribute_type`, `begins_with`, `contains`, `missing`). Bag literals (`<<…>>`), nested document paths (`a.b.c[0]`), parameterized statements (`?`), and the `IS MISSING` predicate are all represented.

Output paths (currently in the main worktree, to be brought into the parser worktrees as needed):
- `partiql/parser/testdata/aws-corpus/index.json` — manifest with id, label, source URL, sql for each example
- `partiql/parser/testdata/aws-corpus/<id>.partiql` — one file per example (63 files)

Two entries (`select-001`, `insert-002`) are syntax skeletons with backticks/brackets and are not valid PartiQL — they are flagged in the index for filtering when running through a parser.

### Wiring rule for parser DAG nodes

Each parser DAG node (`parser-foundation` and onwards) is responsible for adding golden tests that load both corpora as inputs. A grammar feature is considered "covered by omni" only when:

- the legacy fixture inputs that exercise it parse without error and the AST matches the expected snapshot, **and**
- the AWS reference examples that exercise it parse without error and the AST matches the expected snapshot.

Missed grammar features are caught the first time an example fails — there is no need to maintain a parallel hand-written checklist.

### Why this lives in the `ast-core` spec

`ast-core` has nothing to validate against these corpora because no parser exists yet. But the AST shape decided here directly determines what golden snapshots are possible — so the corpus policy is recorded alongside the type design rather than buried in a parser-foundation spec written weeks later. Subsequent parser node specs reference back to this section instead of re-deriving the policy.

## Acceptance criteria

The `ast-core` DAG node is **done** when:

1. All ~73 node types defined across `node.go`, `literals.go`, `stmts.go`, `exprs.go`, `tableexprs.go`, `types.go`, `patterns.go`, and `outfuncs.go`
2. Every node has a doc comment naming the grammar rule(s) from `analysis.md` that it represents
3. Every node implements `nodeTag()` and `GetLoc() Loc`
4. Every node implements at least one sub-interface (`StmtNode`, `ExprNode`, `TableExpr`, `PathStep`, `TypeName`, `PatternNode`) — except small clause helpers (`TargetEntry`, `CaseWhen`, `LetBinding`, `GroupByClause`, `GroupByItem`, `OrderByItem`, `SetAssignment`, `OnConflict`, `OnConflictTarget`, `ReturningClause`, `ReturningItem`, `WindowSpec`, `TuplePair`, `PivotProjection`, `PatternQuantifier`, `PatternSelector`) which are bare `Node`
5. `outfuncs.NodeToString` covers every type — verified by `TestNodeToString_AllNodesCovered`
6. `go test ./partiql/ast/...` passes
7. `go vet ./partiql/ast/...` clean
8. `gofmt` clean
9. Coverage cross-check: every grammar area in `analysis.md`'s "Full Coverage Target" maps to ≥1 node type, with the cross-check section above re-verified during implementation against `PartiQLParser.g4`

## Non-goals

- **No parser logic** — that's `parser-foundation` through `parser-graph` (DAG nodes 4–7, 12–18). `ast-core` only defines types; no recursive-descent code, no tokens, no error recovery
- **No deparse / SQL printing** — separate (out-of-scope) package per `dag.md` "Out of Scope"
- **No semantic analysis or type checking**
- **No `Walker` / `Visitor` interface** — Go type switches are sufficient; cosmosdb doesn't ship one and we have no consumer that needs one. Add later only if/when a real consumer requires it.
- **No smart constructors / builder helpers** — parsers and tests construct nodes directly with struct literals. Less indirection, less to maintain.
- **No round-trip tests against either test corpus** (legacy ANTLR fixtures or the AWS PartiQL reference examples) — that is gated by `parser-foundation`. See "Test corpus (project-wide forward reference)" above for the full strategy.
- **Graph-pattern AST shape may be refined** when `parser-graph` (DAG node 16) is implemented. The marker interface (`PatternNode`), node names, and multi-interface relationship to `ExprNode` are stable; only field shapes inside `NodePattern`/`EdgePattern` are at risk.
- **No `Loc` line/column** — only byte offsets. Conversion happens at the public `partiql/parse.go` boundary.

## Risks & mitigations

| Risk | Mitigation |
|------|-----------|
| Missing a grammar rule from `analysis.md` | Acceptance criterion 9 forces a manual line-by-line re-read of `PartiQLParser.g4` during implementation; the cross-check table in this spec is the working list |
| Graph pattern AST shape needs adjustment when `parser-graph` lands | Marker interface is stable; node names are stable; only the field shapes inside `NodePattern`/`EdgePattern` are at risk. Documented as a known refinement point. |
| Multi-interface nodes (`PathExpr`, `VarRef`, `SubLink`) confuse parser callers | Compile-time `var _` assertions in `ast_test.go` make the dual-interface explicit; doc comments on each multi-interface node spell out the grammar reason |
| `outfuncs.go` becomes stale when node types are added later | `TestNodeToString_AllNodesCovered` reflection test catches missing arms |
| Field bloat on large nodes (`SelectStmt` has 13 fields) | Acceptable — mirrors PostgreSQL `SelectStmt` and cosmosdb `SelectStmt`. Subdividing would create artificial structure. |

## References

- `docs/migration/partiql/analysis.md` — full legacy grammar coverage
- `docs/migration/partiql/dag.md` — migration node ordering
- `bytebase/parser/partiql/PartiQLParser.g4` — authoritative grammar source
- `cosmosdb/ast/node.go`, `cosmosdb/ast/parsenodes.go` — sealed-interface reference
- `mongo/ast/nodes.go` — looser flat reference (not adopted, but compared)
