# T1.8 — Snowflake Statement-Classification Helper Design

## Scope

T1.8 adds a small utility package `snowflake/analysis` that classifies any parsed
Snowflake AST node into one of six statement categories: SELECT, DML, DDL, SHOW,
DESCRIBE, or Other. An additional convenience function classifies raw SQL strings
directly.

This is a P0 utility used by bytebase to route statements to the right handler
(e.g. query-span extraction vs. schema migration vs. read-only check).

## Package

New package: `snowflake/analysis`

Files:
- `classify.go` — Category enum, Classify(node), ClassifySQL(sql)
- `classify_test.go` — test coverage

## Category enum

```go
type Category int

const (
    CategoryUnknown  Category = iota
    CategorySelect            // SELECT, WITH+SELECT, UNION/INTERSECT/EXCEPT
    CategoryDML               // INSERT, UPDATE, DELETE, MERGE
    CategoryDDL               // CREATE/ALTER/DROP/UNDROP/TRUNCATE/COMMENT ON
    CategoryShow              // SHOW (forward-compatible placeholder)
    CategoryDescribe          // DESCRIBE/DESC (forward-compatible placeholder)
    CategoryOther             // USE, SET, EXPLAIN, CALL, etc.
)
```

## Classify(node ast.Node) Category

Uses a type switch on the concrete `ast.Node` type.

### CategorySelect
- `*ast.SelectStmt`
- `*ast.SetOperationStmt`

### CategoryDML
- `*ast.InsertStmt`
- `*ast.InsertMultiStmt`
- `*ast.UpdateStmt`
- `*ast.DeleteStmt`
- `*ast.MergeStmt`

### CategoryDDL
- `*ast.CreateTableStmt`
- `*ast.AlterTableStmt`
- `*ast.CreateDatabaseStmt`
- `*ast.AlterDatabaseStmt`
- `*ast.DropDatabaseStmt`
- `*ast.UndropDatabaseStmt`
- `*ast.CreateSchemaStmt`
- `*ast.AlterSchemaStmt`
- `*ast.DropSchemaStmt`
- `*ast.UndropSchemaStmt`
- `*ast.CreateViewStmt`
- `*ast.CreateMaterializedViewStmt`
- `*ast.AlterViewStmt`
- `*ast.AlterMaterializedViewStmt`
- `*ast.DropStmt`
- `*ast.UndropStmt`

### CategoryShow / CategoryDescribe
Not yet implemented in the parser; placeholders ensure forward-compatible routing
when those nodes are added later. The type switch default handles them as
CategoryUnknown until dedicated node types exist.

### Default (CategoryUnknown)
- `nil` input
- unrecognized node types (e.g. expression nodes passed by mistake)
- future node types not yet classified

## ClassifySQL(sql string) Category

Parses the first statement from `sql` using `parser.ParseBestEffort`, then calls
`Classify` on the first node. Returns `CategoryUnknown` if:
- the SQL is empty
- parsing yields no nodes (e.g. only whitespace or comments)
- the first node is nil

## String() method

```
CategoryUnknown  → "Unknown"
CategorySelect   → "SELECT"
CategoryDML      → "DML"
CategoryDDL      → "DDL"
CategoryShow     → "SHOW"
CategoryDescribe → "DESCRIBE"
CategoryOther    → "Other"
```

## Design notes

- No new AST nodes are added; this is purely a read-only classification.
- The package imports only `snowflake/ast` and `snowflake/parser` — no circular deps.
- ClassifySQL intentionally ignores parse errors to be lenient; callers that need
  strict validation should use the diagnostics package.
- CategoryOther is intentionally NOT the same as CategoryUnknown: Unknown means
  "we don't know how to classify this node", while Other means "we recognize this
  as a statement but it is not SELECT/DML/DDL/SHOW/DESCRIBE" (e.g. USE, SET, CALL).
  Since the parser does not yet emit nodes for those statements, both currently
  fall through to CategoryUnknown in practice; the distinction becomes meaningful
  once those nodes are added.
