# Snowflake Query Span Extractor — Design Spec (T3.1)

**Date**: 2026-04-15
**Node**: T3.1
**Package**: `snowflake/analysis`

---

## Goal

Given a parsed Snowflake SELECT/WITH/set-operation AST node, produce a `QuerySpan` that describes:

1. **Result columns** — one entry per SELECT-list position with name + lineage (which source columns it traces back to, or empty for constants/derived).
2. **Source columns** — every source column read by the entire query (flat union from all branches).
3. **CTE transparency** — CTE references resolve to their body's span; CTE sources flow into the outer query's Sources.
4. **Set operators** — positional merge of result columns across branches (UNION/INTERSECT/EXCEPT), or named merge for `UNION ALL BY NAME`.
5. **EXCLUDE** — `SELECT * EXCLUDE (col)` drops the named columns from `*` expansion.
6. **Subqueries in FROM** — resolved to a "virtual table" whose columns come from the subquery's ResultColumns.

---

## Type Definitions

```go
package analysis

// QuerySpan summarises what a query reads and produces.
type QuerySpan struct {
    Results []*ResultColumn
    Sources []*SourceColumn
}

// ResultColumn is one column in the SELECT result set.
type ResultColumn struct {
    Name      string          // output alias or derived name; "*" for unresolved star
    Sources   []*SourceColumn // source columns this result traces back to
    IsDerived bool            // true if computed (function, expression, multi-source)
}

// SourceColumn identifies a column read from a base table or CTE.
type SourceColumn struct {
    Database string // may be empty
    Schema   string // may be empty
    Table    string // table name or CTE alias
    Column   string // column name; "*" for unresolved star
}
```

---

## Architecture: Recursive Extractor

The extractor is a pure recursive function (no global mutable state). The core function signature:

```go
func extractSpan(node ast.Node, scope *queryScope) *QuerySpan
```

`queryScope` carries CTE definitions visible at the current scope level.

### queryScope

```go
type queryScope struct {
    ctes map[string]*QuerySpan // CTE name (uppercase) → its span
}
```

### Extraction Algorithm

#### extractSelectStmt(s *ast.SelectStmt, scope *queryScope) *QuerySpan

1. Build the **FROM scope**: for each item in `s.From`, resolve it to a `tableEntry` (name + its known columns).
2. For each `*SelectTarget` in `s.Targets`:
   - If `Star == true` and `Exclude == nil/empty`: emit one `ResultColumn{Name: "*", Sources: [{Table: qualifier or each_from_table, Column: "*"}]}`
   - If `Star == true` and `Exclude` has entries: emit one `ResultColumn{Name: "*", Sources: filtered}` (drop excluded column names from sources)
   - If `Star == false`: resolve the expression's column refs into `SourceColumn` entries; determine `IsDerived`; pick name from alias or expression
3. Collect all source columns from the FROM scope into `QuerySpan.Sources`.
4. Return the assembled `QuerySpan`.

#### FROM scope building

Each item in `s.From` is either `*ast.TableRef` or `*ast.JoinExpr`.

- `*ast.TableRef` with `Name != nil`: simple table reference. Look up name in scope.ctes first; if found, use the CTE's ResultColumns as the virtual table's columns. Otherwise, treat as a base table with `SourceColumn{Table: name}`.
- `*ast.TableRef` with `Subquery != nil`: recurse into the subquery; the subquery's `ResultColumns` become the virtual table's columns. The alias of the `TableRef` is the virtual table name.
- `*ast.JoinExpr`: recursively build FROM scope for `Left` and `Right`, merge them.

#### Expression column ref extraction (collectRefs)

Walk an expression node with `ast.Inspect`. Collect every `*ast.ColumnRef` found. For each `ColumnRef`:
- 1-part: `{Column: parts[0]}`; resolve table via FROM scope (ambiguous = no table assigned)
- 2-part: `{Table: parts[0], Column: parts[1]}`
- 3-part: `{Schema: parts[0], Table: parts[1], Column: parts[2]}`
- 4-part: `{Database: parts[0], Schema: parts[1], Table: parts[2], Column: parts[3]}`

A `*ast.StarExpr` inside an expression (e.g. `COUNT(*)`) is treated as a synthetic pseudo-source `{Table: "", Column: "*"}`.

#### IsDerived determination

A `ResultColumn.IsDerived` is `true` when:
- The expression is a `*ast.FuncCallExpr` (function call, including COUNT(*))
- The expression is a `*ast.BinaryExpr`, `*ast.UnaryExpr`, `*ast.CaseExpr`, `*ast.IffExpr`, `*ast.CastExpr`, or any other non-trivial expression
- The expression is a `*ast.Literal` (constants — no sources)
- The expression references more than one source column

A `ResultColumn.IsDerived` is `false` when the expression is a bare `*ast.ColumnRef` (or `*ast.ParenExpr` wrapping a ColumnRef) — direct column pass-through.

#### extractSetOperationStmt(s *ast.SetOperationStmt, scope *queryScope) *QuerySpan

1. Recurse: `leftSpan = extractSpan(s.Left, scope)`, `rightSpan = extractSpan(s.Right, scope)`.
2. If `s.ByName`:
   - Build a map of right ResultColumns by name.
   - Merge: for each left ResultColumn, find matching right by name; merge sources.
   - Append any right ResultColumns not matched by left.
3. Otherwise (positional):
   - Zip `leftSpan.Results` with `rightSpan.Results` positionally.
   - Each merged ResultColumn takes the name from the left branch, sources = union of left+right sources.
   - If branches differ in length, take the longer list (append unmatched tail from longer branch).
4. `Sources` = union of `leftSpan.Sources` + `rightSpan.Sources`.

#### CTE handling

Before extracting the main query body, process `s.With`:
1. Create child `queryScope`.
2. For each `*ast.CTE` in `s.With`:
   a. Recurse: `cteSpan = extractSpan(cte.Query, childScope)`.
   b. If `cte.Columns` is non-empty, rename ResultColumns by position (column aliases).
   c. Store in `childScope.ctes[cte.Name.Normalize()] = cteSpan`.
3. Pass `childScope` into the main body extraction.

---

## Result Column Name Derivation

- If `SelectTarget.Alias` is set → use `Alias.Name`
- Else if `Expr` is a `*ast.ColumnRef` → use last part name
- Else if `Expr` is a `*ast.FuncCallExpr` → use function name (last part of `Name`)
- Else → use empty string (anonymous derived column)

---

## Star Expansion Without Catalog

Since there's no catalog, `SELECT *` produces a single pseudo result column:
```
ResultColumn{Name: "*", Sources: [SourceColumn{Table: <from_alias>, Column: "*"}], IsDerived: false}
```

For `SELECT * EXCLUDE (b, c)`, the `SelectTarget.Exclude` field contains the excluded column names. Since we can't expand `*` without a catalog, we still emit a single `ResultColumn{Name: "*"}` but record the excluded columns in a way that is visible to the caller. The simplest approach: emit one `ResultColumn` per non-excluded known column when the from-source is a CTE (we know its columns), otherwise keep `Name: "*"` with a note.

---

## Deduplication

`QuerySpan.Sources` is the flat union of all source columns referenced anywhere in the query. Deduplication is by value (all 4 fields equal).

---

## Public API

```go
func Extract(node ast.Node) (*QuerySpan, error)
func ExtractSQL(sql string) (*QuerySpan, error)
```

`Extract` accepts `*ast.SelectStmt`, `*ast.SetOperationStmt`, or `*ast.File` (first stmt).
`ExtractSQL` calls `parser.Parse` then `Extract`.

---

## Limitations (by design, not bugs)

- No catalog validation: columns are not checked against real schemas.
- `SELECT *` from a base table produces `Name: "*"` pseudo-column, not expanded columns.
- Recursive CTEs are processed once (no recursion expansion).
- PIVOT/UNPIVOT/LATERAL FLATTEN are not handled (T5.3 scope).
- Subquery expressions in WHERE/HAVING are not traced for result lineage (only FROM subqueries are).
