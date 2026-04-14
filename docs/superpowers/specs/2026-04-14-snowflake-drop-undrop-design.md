# Spec: T2.5 Snowflake DROP / UNDROP DDL

## Scope

Implement DROP and UNDROP statement parsing for Snowflake SQL.
Does NOT include DATABASE or SCHEMA (those belong to T2.1).

## Object types covered

### DROP (with IF EXISTS and CASCADE/RESTRICT where applicable)

| Object type         | IF EXISTS | CASCADE/RESTRICT | Notes                         |
|---------------------|-----------|------------------|-------------------------------|
| TABLE               | yes       | yes              | Core target                   |
| VIEW                | yes       | no               |                               |
| MATERIALIZED VIEW   | yes       | no               |                               |
| DYNAMIC TABLE       | no        | no               | No IF EXISTS in legacy grammar|
| EXTERNAL TABLE      | yes       | yes              |                               |
| STREAM              | yes       | no               |                               |
| TASK                | yes       | no               |                               |
| SEQUENCE            | yes       | yes              |                               |
| STAGE               | yes       | no               |                               |
| FILE FORMAT         | yes       | no               | Two-keyword object type       |
| FUNCTION            | yes       | no               | No arg_types parsing (stub)   |
| PROCEDURE           | yes       | no               | No arg_types parsing (stub)   |
| PIPE                | yes       | no               |                               |
| TAG                 | yes       | no               |                               |
| ROLE                | yes       | no               |                               |
| WAREHOUSE           | yes       | no               |                               |

For object types not in this list, the parser emits a targeted error:
"DROP <X> statement parsing is not yet supported" and skips to the next
statement boundary rather than a generic "DROP not supported" error.

### UNDROP

| Object type   | Notes                     |
|---------------|---------------------------|
| TABLE         | Core target               |
| DYNAMIC TABLE | UNDROP DYNAMIC TABLE name |
| TAG           | UNDROP TAG name           |

DATABASE and SCHEMA UNDROP are excluded (T2.1 territory).

## AST design

### Unified DropStmt

Use a single `DropStmt` type with an `ObjectKind` enum. This is cleaner
than per-object-type structs because DROP forms are structurally near-identical.

```go
type DropObjectKind int

const (
    DropTable DropObjectKind = iota
    DropView
    DropMaterializedView
    DropDynamicTable
    DropExternalTable
    DropStream
    DropTask
    DropSequence
    DropStage
    DropFileFormat
    DropFunction
    DropProcedure
    DropPipe
    DropTag
    DropRole
    DropWarehouse
)

type DropStmt struct {
    Kind     DropObjectKind
    IfExists bool
    Name     *ObjectName
    Cascade  bool // CASCADE (mutually exclusive with Restrict)
    Restrict bool // RESTRICT (mutually exclusive with Cascade)
    Loc      Loc
}
```

### UndropStmt

```go
type UndropObjectKind int

const (
    UndropTable UndropObjectKind = iota
    UndropDynamicTable
    UndropTag
)

type UndropStmt struct {
    Kind UndropObjectKind
    Name *ObjectName
    Loc  Loc
}
```

### Node tags

Add `T_DropStmt` and `T_UndropStmt` to `nodetags.go`.

## Parser dispatch design

```
parseDropStmt()
  advance() // consume DROP
  switch cur:
    TABLE          → parseDropObject(DropTable, ifExistsOK=true, cascadeOK=true)
    VIEW           → parseDropObject(DropView, ifExistsOK=true, cascadeOK=false)
    MATERIALIZED VIEW → parseDropObject(DropMaterializedView, ...)
    DYNAMIC TABLE  → parseDropObject(DropDynamicTable, ifExistsOK=false, ...)
    EXTERNAL TABLE → parseDropObject(DropExternalTable, ifExistsOK=true, cascadeOK=true)
    STREAM         → parseDropObject(DropStream, ifExistsOK=true, ...)
    TASK           → parseDropObject(DropTask, ifExistsOK=true, ...)
    SEQUENCE       → parseDropObject(DropSequence, ifExistsOK=true, cascadeOK=true)
    STAGE          → parseDropObject(DropStage, ifExistsOK=true, ...)
    FILE FORMAT    → parseDropObject(DropFileFormat, ifExistsOK=true, ...)
    FUNCTION       → parseDropObject(DropFunction, ifExistsOK=true, ...)
    PROCEDURE      → parseDropObject(DropProcedure, ifExistsOK=true, ...)
    PIPE           → parseDropObject(DropPipe, ifExistsOK=true, ...)
    TAG            → parseDropObject(DropTag, ifExistsOK=true, ...)
    ROLE           → parseDropObject(DropRole, ifExistsOK=true, ...)
    WAREHOUSE      → parseDropObject(DropWarehouse, ifExistsOK=true, ...)
    default        → targeted unsupported error
```

`parseDropObject(kind, ifExistsOK, cascadeOK)` is a shared helper that:
1. Optionally consumes IF EXISTS (if ifExistsOK and next token matches)
2. Parses the object name
3. Optionally consumes CASCADE or RESTRICT (if cascadeOK)
4. Returns *DropStmt

`parseUndropStmt()` mirrors the same pattern.

## Test coverage

- DROP TABLE [IF EXISTS] name [CASCADE|RESTRICT]
- DROP VIEW, DROP STREAM, DROP TASK, DROP SEQUENCE, DROP STAGE
- DROP FILE FORMAT (two-word object type)
- DROP MATERIALIZED VIEW (two-word object type)
- DROP DYNAMIC TABLE (two-word object type)
- DROP EXTERNAL TABLE (two-word, with CASCADE)
- DROP FUNCTION, DROP PROCEDURE, DROP PIPE, DROP TAG, DROP ROLE, DROP WAREHOUSE
- UNDROP TABLE, UNDROP DYNAMIC TABLE, UNDROP TAG
- Unknown object types produce targeted error (DROP SCHEMA → unsupported)
- Qualified names (schema.table, db.schema.table)
- Loc tracking
