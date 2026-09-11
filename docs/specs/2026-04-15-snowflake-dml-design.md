# Snowflake DML (T5.1) — Design

**DAG node:** T5.1 — DML: INSERT/UPDATE/DELETE/MERGE
**Branch:** `feat/snowflake/dml`
**Depends on:** T1.4 (SELECT core), T1.6 (CTEs), expressions
**Unblocks:** T3.2 (deparse)

## Purpose

T5.1 adds the four primary DML statement parsers for Snowflake SQL. These replace
the `unsupported("INSERT")`, `unsupported("UPDATE")`, `unsupported("DELETE")`, and
`unsupported("MERGE")` stubs in `parser.go`.

## Scope

Derived from the legacy `SnowflakeParser.g4` grammar:
- `insert_statement` — single-row INSERT with VALUES or SELECT
- `insert_multi_table_statement` — INSERT ALL / INSERT FIRST (multi-table)
- `update_statement` — UPDATE with optional FROM (Snowflake extension)
- `delete_statement` — DELETE with optional USING
- `merge_statement` — MERGE INTO ... USING ... ON ... WHEN clauses

## AST Types

### InsertStmt (single-row INSERT)

```go
type InsertStmt struct {
    Overwrite bool
    Target    *ObjectName
    Columns   []Ident    // optional; nil if not specified
    Values    [][]Node   // VALUES rows; nil if SELECT used
    Select    Node       // SELECT body; nil if VALUES used
    Loc       Loc
}
```

### InsertMultiStmt (INSERT ALL / INSERT FIRST)

```go
type InsertMultiStmt struct {
    Overwrite bool
    First     bool                 // true = FIRST, false = ALL
    Branches  []*InsertMultiBranch // INTO targets
    Select    Node                 // driving SELECT
    Loc       Loc
}

type InsertMultiBranch struct {
    When    Node    // nil for unconditional; non-nil for WHEN cond THEN
    Target  *ObjectName
    Columns []Ident
    Values  []Node  // single-row values; nil if VALUES omitted
    Loc     Loc
}
```

### UpdateStmt

```go
type UpdateStmt struct {
    Target *ObjectName
    Alias  Ident
    Sets   []*UpdateSet
    From   []Node   // FROM items (Snowflake extension)
    Where  Node
    Loc    Loc
}

type UpdateSet struct {
    Column Ident
    Value  Node
    Loc    Loc
}
```

### DeleteStmt

```go
type DeleteStmt struct {
    Target *ObjectName
    Alias  Ident
    Using  []Node   // USING items
    Where  Node
    Loc    Loc
}
```

### MergeStmt

```go
type MergeStmt struct {
    Target      *ObjectName
    TargetAlias Ident
    Source      Node   // table ref or subquery
    SourceAlias Ident
    On          Node
    Whens       []*MergeWhen
    Loc         Loc
}

type MergeWhen struct {
    Matched       bool        // true = WHEN MATCHED
    BySource      bool        // WHEN NOT MATCHED BY SOURCE
    ByTarget      bool        // WHEN NOT MATCHED [BY TARGET]
    AndCond       Node        // optional AND condition
    Action        MergeAction
    Sets          []*UpdateSet // for UPDATE SET
    InsertCols    []Ident      // for INSERT (cols)
    InsertVals    []Node       // for INSERT VALUES (exprs)
    InsertDefault bool         // for INSERT VALUES DEFAULT
    Loc           Loc
}

type MergeAction int
const (
    MergeActionUpdate MergeAction = iota
    MergeActionDelete
    MergeActionInsert
)
```

## Parser strategy

- `parseInsertStmt` dispatches on ALL/FIRST for multi-table vs single-row
- Reuse `parseFromClause`, `parseFromItem`, `parseTableRef`, `parseQueryExpr`, `parseExprList`
- `parseOptionalAlias` for aliases throughout
- MERGE source uses `parsePrimarySource` (handles table refs and subqueries)
- WHEN NOT MATCHED BY SOURCE/TARGET: `SOURCE` is `kwSOURCE`, `TARGET` is `tokIdent`

## Walker

All new Node types get cases in `walk_generated.go` (regenerated via `go generate ./snowflake/ast/...`).

## Test plan

- INSERT single: VALUES, SELECT, OVERWRITE, column list, multi-value
- INSERT ALL/FIRST: unconditional and conditional
- UPDATE: basic SET, FROM, WHERE, alias
- DELETE: basic, USING, WHERE, alias
- MERGE: WHEN MATCHED UPDATE/DELETE, WHEN NOT MATCHED INSERT, BY SOURCE/TARGET, multiple WHEN
