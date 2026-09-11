# CosmosDB Query Analysis for Omni

## Overview

Add a `cosmosdb/analysis/` package to Omni that extracts field-level information
from parsed CosmosDB SELECT queries: which document properties are projected in
the SELECT list, and which are referenced in the WHERE clause.

This is the Omni-side analysis that Bytebase's query span glue layer will
consume. Omni defines self-contained types and walks its own AST. Bytebase
converts the result into its own `base.QuerySpan` / `base.PathAST` types.

### Prior Work

- **`omni/cosmosdb/parser`** — Phase 1 parser (merged). Produces `*ast.SelectStmt`.
- **`bytebase/bytebase/.../parser/cosmosdb/query_span.go`** — Current ANTLR-based
  implementation (~300 lines). Walks the ANTLR parse tree using listeners to
  extract field paths for data masking.
- **`omni/pg/catalog/`** — PostgreSQL defines its own analyzed `Query` struct;
  Bytebase walks it to build query spans. This is the pattern we follow.

### Design Decisions

1. **Analysis in Omni, glue in Bytebase** — follows the PostgreSQL pattern.
2. **Walk the raw AST** — no semantic analysis layer or catalog needed. The
   current Bytebase code does pure AST walking with alias resolution, and
   that's sufficient for field path extraction.
3. **Return a resolved result struct** — `QueryAnalysis` with projections and
   predicates. Bytebase reads the fields and converts to its own types.
4. **Self-contained types** — Omni defines `FieldPath`, `Selector`,
   `Projection`, `QueryAnalysis`. No dependency on Bytebase's `base` package.

---

## Directory Structure

```
cosmosdb/analysis/
├── analysis.go      # Analyze() entry point, QueryAnalysis/Projection types
├── fieldpath.go     # FieldPath, Selector types
├── extract.go       # extractFieldPaths() recursive walker + helpers
└── analysis_test.go # Table-driven tests
```

---

## Types

File: `cosmosdb/analysis/fieldpath.go`

### Selector

Represents one step in a document property path — either a named property
access or an array index access.

```go
type Selector struct {
    Name       string // property name
    ArrayIndex int    // -1 for item access, >= 0 for array index access
}
```

- `Selector{Name: "address", ArrayIndex: -1}` represents `.address`
- `Selector{Name: "addresses", ArrayIndex: 1}` represents `.addresses[1]`

### FieldPath

A chain of selectors representing a path through a JSON document.

```go
type FieldPath []Selector
```

Example: `c.addresses[1].country` is represented as:

```go
FieldPath{
    {Name: "container", ArrayIndex: -1},  // "c" resolved to container name
    {Name: "addresses", ArrayIndex: 1},
    {Name: "country",   ArrayIndex: -1},
}
```

File: `cosmosdb/analysis/analysis.go`

### Projection

One item in the SELECT list with its output name and source field paths.

```go
type Projection struct {
    Name        string      // output alias; empty if none and not inferrable
    SourcePaths []FieldPath // source field paths this projection references
}
```

A projection can reference multiple source paths when the expression combines
fields: `SELECT c.name ?? c.nickname AS displayName` produces one `Projection`
with two `SourcePaths`.

### QueryAnalysis

The top-level analysis result.

```go
type QueryAnalysis struct {
    Projections []Projection // SELECT list items (empty when SelectStar is true)
    SelectStar  bool         // true if SELECT *
    Predicates  []FieldPath  // field paths referenced in WHERE
}
```

`SelectStar` and `Projections` are mutually exclusive — CosmosDB's grammar
does not allow mixing `*` with explicit properties. When `SelectStar` is true,
Bytebase applies masking policy to all fields rather than checking specific
paths.

---

## Public API

```go
func Analyze(stmt *ast.SelectStmt) *QueryAnalysis
```

Takes a parsed `*ast.SelectStmt` (from `cosmosdb.Parse()`), returns the
analysis result. No error return — if the query parsed successfully, analysis
always succeeds (it is pure AST walking).

---

## Analysis Logic

### Step 1: Build alias resolution map from FROM/JOIN

Walk `stmt.From` and `stmt.Joins` to build a map from alias names to their
resolved source paths.

For a simple query `SELECT c.name FROM container c`:
- `c` -> `FieldPath{{"container", -1}}`

For JOINs `SELECT p.name, t.tag FROM products p JOIN t IN p.tags`:
- `p` -> `FieldPath{{"products", -1}}`
- `t` -> `FieldPath{{"products", -1}, {"tags", -1}}`

For nested JOINs `... JOIN t IN p.tags JOIN s IN t.sizes`:
- `p` -> `FieldPath{{"products", -1}}`
- `t` -> `FieldPath{{"products", -1}, {"tags", -1}}`
- `s` -> `FieldPath{{"products", -1}, {"tags", -1}, {"sizes", -1}}`

The alias map is built by:
1. Walking `stmt.From` for the primary container name and optional alias.
   - `ContainerRef{Name: "products"}` + `AliasedTableExpr{Alias: "p"}` ->
     alias `p` resolves to `products`.
   - If no alias, the container name is used directly.
2. Walking each `JoinExpr` whose source is an `ArrayIterationExpr`:
   - `ArrayIterationExpr{Alias: "t", Source: <DotAccessExpr p.tags>}` ->
     resolve `p` through the alias map, append `tags`, register alias `t`.

### Step 2: Extract projections

If `stmt.Star` is true, set `SelectStar: true` and skip projection extraction.

Otherwise, for each `TargetEntry` in `stmt.Targets`:
1. Walk `TargetEntry.Expr` with `extractFieldPaths()` to collect source paths.
2. Determine the output name:
   - Use `TargetEntry.Alias` if present.
   - Otherwise use the last property name in the path (e.g., `c.name` -> `name`).
   - Otherwise empty string (for expressions without alias).

### Step 3: Extract predicates

If `stmt.Where` is not nil, walk the WHERE expression with
`extractFieldPaths()` to collect all referenced field paths.

### The core walker: `extractFieldPaths`

```go
func extractFieldPaths(expr ast.ExprNode, aliases aliasMap) []FieldPath
```

Recursive switch on AST node type:

| AST Node | Action |
|----------|--------|
| `ColumnRef` | Start new path; resolve name through alias map |
| `DotAccessExpr` | Recurse into Expr, append property to each path |
| `BracketAccessExpr` | Recurse into Expr, append selector (item or array index) |
| `BinaryExpr` | Collect paths from Left and Right |
| `UnaryExpr` | Recurse into Operand |
| `TernaryExpr` | Collect from Cond, Then, Else |
| `FuncCall` | Collect from all Args |
| `UDFCall` | Collect from all Args |
| `InExpr` | Collect from Expr and all List items |
| `BetweenExpr` | Collect from Expr, Low, High |
| `LikeExpr` | Collect from Expr and Pattern |
| `CreateArrayExpr` | Collect from all Elements |
| `CreateObjectExpr` | Collect from all Field values |
| Literals (`StringLit`, `NumberLit`, `BoolLit`, `NullLit`, `UndefinedLit`, `InfinityLit`, `NanLit`) | Return nil |
| `ParamRef` | Return nil |
| `SubLink`, `ExistsExpr`, `ArrayExpr` | No traversal into subqueries (matches current Bytebase behavior) |

---

## Test Strategy

Table-driven tests in `analysis_test.go`. Each test case: SQL string ->
expected `QueryAnalysis`.

### Test cases

**SELECT * detection:**
- `SELECT * FROM c` -> `SelectStar: true`, no projections

**Simple projections:**
- `SELECT c.name FROM c` -> projection `name` with path `[c, name]`
- `SELECT c.name AS n FROM c` -> projection `n` with path `[c, name]`

**Expression projections (multiple source paths):**
- `SELECT c.name ?? c.nickname FROM c` -> one projection with two paths

**Dot access chains:**
- `SELECT c.address.city FROM c` -> path `[c, address, city]`

**Bracket access:**
- `SELECT c.addresses[1].country FROM c` -> path with array selector
- `SELECT c["field"] FROM c` -> path with item selector

**Function calls:**
- `SELECT CONCAT(c.first, c.last) FROM c` -> two source paths
- `SELECT udf.myFunc(c.name) FROM c` -> one source path

**JOINs:**
- `SELECT p.name, t.tag FROM products p JOIN t IN p.tags` -> two projections,
  `t.tag` resolves through alias to `products.tags.tag`
- Nested JOINs: `... JOIN t IN p.tags JOIN s IN t.sizes`

**WHERE predicates:**
- `SELECT * FROM c WHERE c.country = "US"` -> predicate path `[c, country]`
- `SELECT * FROM c WHERE c.age > 18 AND c.active = true` -> two predicate paths
- `SELECT * FROM c WHERE c.population BETWEEN 100000 AND 5000000` -> one path
- `SELECT * FROM c WHERE c.country IN ("US", "UK")` -> one path

**Alias resolution:**
- `SELECT t.name FROM container t` -> path `[container, name]`

**SELECT VALUE:**
- `SELECT VALUE c.name FROM c` -> projection with path `[c, name]`

**Edge cases:**
- `SELECT 1` (no FROM) -> no paths
- `SELECT VALUE COUNT(1) FROM c` -> no field paths (literal argument)

---

## Bytebase Glue Layer (out of scope for Omni)

For context, Bytebase's adapter would look roughly like:

```go
func GetQuerySpan(ctx context.Context, _ base.GetQuerySpanContext, stmt base.Statement, ...) (*base.QuerySpan, error) {
    stmts, err := cosmosdb.Parse(stmt.Text)
    // ...
    sel := stmts[0].AST.(*ast.SelectStmt)
    qa := analysis.Analyze(sel)

    // Convert qa.Projections -> []base.QuerySpanResult
    // Convert qa.Predicates -> map[string]*base.PathAST
    // Convert analysis.FieldPath -> base.PathAST linked list
}
```

This is Bytebase's responsibility and not part of this spec.
