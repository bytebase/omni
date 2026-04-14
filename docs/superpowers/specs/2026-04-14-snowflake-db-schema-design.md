# Spec: T2.1 — DATABASE + SCHEMA DDL (Snowflake)

**Date:** 2026-04-14
**Author:** Claude (T2.1 implementation)
**Node:** T2.1 in the Snowflake migration DAG
**Depends on:** T2.2 (CREATE TABLE patterns), T1.4 (SELECT core), F1–F4 (AST/lexer/splitter/parser-entry)

---

## Goal

Implement full parse-tree support for DATABASE and SCHEMA DDL statements in the
Snowflake omni parser. The eight target statements are:

- `CREATE DATABASE`
- `CREATE SCHEMA`
- `ALTER DATABASE`
- `ALTER SCHEMA`
- `DROP DATABASE`
- `DROP SCHEMA`
- `UNDROP DATABASE`
- `UNDROP SCHEMA`

## Source of truth

Legacy grammar: `/parser/snowflake/SnowflakeParser.g4`, rules
`create_database`, `create_schema`, `alter_database`, `alter_schema`,
`drop_database`, `drop_schema`, `undrop_database`, `undrop_schema`.

---

## AST design

### New node types (in `parsenodes.go`)

#### `CreateDatabaseStmt`

```
OrReplace       bool
Transient       bool
IfNotExists     bool
Name            *ObjectName
Clone           *CloneSource         // CLONE ... [AT|BEFORE ...]
DataRetention   *int64               // DATA_RETENTION_TIME_IN_DAYS = n
MaxDataExt      *int64               // MAX_DATA_EXTENSION_TIME_IN_DAYS = n
DefaultDDLColl  *string              // DEFAULT_DDL_COLLATION = 's'
Comment         *string
Tags            []*TagAssignment
Loc             Loc
```

#### `CreateSchemaStmt`

```
OrReplace       bool
Transient       bool
IfNotExists     bool
Name            *ObjectName
Clone           *CloneSource
ManagedAccess   bool                 // WITH MANAGED ACCESS
DataRetention   *int64
MaxDataExt      *int64
DefaultDDLColl  *string
Comment         *string
Tags            []*TagAssignment
Loc             Loc
```

#### `AlterDatabaseStmt`

Represents all ALTER DATABASE variants. Uses a discriminated `Action` field:

```
IfExists    bool
Name        *ObjectName
Action      AlterDatabaseAction  (enum)
// Populated based on Action:
NewName     *ObjectName       // RENAME TO / SWAP WITH
SetProps    *DBSchemaProps    // SET ...
UnsetProps  []string          // UNSET property, ...
Tags        []*TagAssignment  // SET TAG (...)
UnsetTags   []*ObjectName     // UNSET TAG (...)
Loc         Loc
```

`AlterDatabaseAction` enum:
- `AlterDBRename` — RENAME TO
- `AlterDBSwap` — SWAP WITH
- `AlterDBSet` — SET properties
- `AlterDBUnset` — UNSET properties
- `AlterDBSetTag` — SET TAG
- `AlterDBUnsetTag` — UNSET TAG
- `AlterDBEnableReplication` — ENABLE REPLICATION (structural only)
- `AlterDBDisableReplication` — DISABLE REPLICATION (structural only)
- `AlterDBEnableFailover` — ENABLE FAILOVER (structural only)
- `AlterDBDisableFailover` — DISABLE FAILOVER (structural only)
- `AlterDBRefresh` — REFRESH
- `AlterDBPrimary` — PRIMARY

`DBSchemaProps` helper struct (not a Node):
```
DataRetention  *int64
MaxDataExt     *int64
DefaultDDLColl *string
Comment        *string
```

#### `AlterSchemaStmt`

```
IfExists      bool
Name          *ObjectName
Action        AlterSchemaAction  (enum)
NewName       *ObjectName       // RENAME TO / SWAP WITH
SetProps      *DBSchemaProps
UnsetProps    []string
Tags          []*TagAssignment
UnsetTags     []*ObjectName
ManagedAccess *bool             // true=ENABLE, false=DISABLE, nil=not set
Loc           Loc
```

`AlterSchemaAction` enum:
- `AlterSchemaRename`
- `AlterSchemaSwap`
- `AlterSchemaSet`
- `AlterSchemaUnset`
- `AlterSchemaSetTag`
- `AlterSchemaUnsetTag`
- `AlterSchemaEnableManagedAccess`
- `AlterSchemaDisableManagedAccess`

#### `DropDatabaseStmt`

```
IfExists  bool
Name      *ObjectName
Cascade   bool         // CASCADE
Restrict  bool         // RESTRICT
Loc       Loc
```

#### `DropSchemaStmt`

```
IfExists  bool
Name      *ObjectName
Cascade   bool
Restrict  bool
Loc       Loc
```

#### `UndropDatabaseStmt`

```
Name  *ObjectName
Loc   Loc
```

#### `UndropSchemaStmt`

```
Name  *ObjectName
Loc   Loc
```

---

## Parser dispatch design

### `parseCreateStmt` extension (`create_table.go`)

Add `kwDATABASE` and `kwSCHEMA` cases to the existing switch. Both DATABASE and
SCHEMA accept TRANSIENT; SCHEMA does NOT accept TEMPORARY/VOLATILE.

The `transient` boolean already computed by `parseCreateStmt` is forwarded to
the new sub-parsers.

### New top-level dispatchers in `parser.go`

Replace `unsupported("ALTER")`, `unsupported("DROP")`, `unsupported("UNDROP")`
with sub-dispatchers that peek at the object-type keyword:

```
parseAlterStmt():
  consume ALTER
  switch cur.Type:
    case kwDATABASE: return parseAlterDatabaseStmt(...)
    case kwSCHEMA:   return parseAlterSchemaStmt(...)
    default:         return unsupported("ALTER")

parseDropStmt():
  consume DROP
  switch cur.Type:
    case kwDATABASE: return parseDropDatabaseStmt(...)
    case kwSCHEMA:   return parseDropSchemaStmt(...)
    default:         return unsupported("DROP")

parseUndropStmt():
  consume UNDROP
  switch cur.Type:
    case kwDATABASE: return parseUndropDatabaseStmt()
    case kwSCHEMA:   return parseUndropSchemaStmt()
    default:         return unsupported("UNDROP")
```

---

## Property parsing helpers

### `parseDBSchemaProps()` — shared by CREATE and ALTER SET

Parses zero or more optional properties from the set:
- `DATA_RETENTION_TIME_IN_DAYS = n`
- `MAX_DATA_EXTENSION_TIME_IN_DAYS = n`
- `DEFAULT_DDL_COLLATION = 'str'` (also `DEFAULT_DDL_COLLATION_`)
- `COMMENT = 'str'`

Returns a `*DBSchemaProps`.

### `parsePropertyNameList()` — for UNSET

Returns a `[]string` of uppercased property names.

### `parseUnsetTagList()` — for UNSET TAG

Parses `TAG ( name, name, ... )`, returns `[]*ObjectName`.

---

## Walker entries

The genwalker will add cases for:
- `*CreateDatabaseStmt` — walk Name, Clone.Source (via *ObjectName)
- `*CreateSchemaStmt` — walk Name, Clone.Source
- `*AlterDatabaseStmt` — walk Name, NewName, Tags[].Name, UnsetTags
- `*AlterSchemaStmt` — walk Name, NewName, Tags[].Name, UnsetTags
- `*DropDatabaseStmt` — walk Name
- `*DropSchemaStmt` — walk Name
- `*UndropDatabaseStmt` — walk Name
- `*UndropSchemaStmt` — walk Name

Walker generation: `go run ./snowflake/ast/cmd/genwalker` from the worktree root.

---

## NodeTag additions (in `nodetags.go`)

Eight new tags:
`T_CreateDatabaseStmt`, `T_CreateSchemaStmt`,
`T_AlterDatabaseStmt`, `T_AlterSchemaStmt`,
`T_DropDatabaseStmt`, `T_DropSchemaStmt`,
`T_UndropDatabaseStmt`, `T_UndropSchemaStmt`

---

## Test coverage targets

File: `snowflake/parser/database_schema_test.go`

| Group | Cases |
|-------|-------|
| `TestCreateDatabase_*` | Basic, OrReplace, Transient, IfNotExists, WithClone, WithProps, WithTags |
| `TestCreateSchema_*` | Basic, OrReplace, Transient, IfNotExists, WithClone, ManagedAccess, WithProps, WithTags |
| `TestAlterDatabase_*` | RenameToNew, SwapWith, SetProps, SetComment, UnsetProp, SetTag, UnsetTag, EnableReplication, DisableReplication, EnableFailover, Refresh, Primary |
| `TestAlterSchema_*` | Rename, Swap, SetProps, Unset, SetTag, UnsetTag, EnableManagedAccess, DisableManagedAccess |
| `TestDropDatabase_*` | Basic, IfExists, Cascade, Restrict |
| `TestDropSchema_*` | Basic, IfExists, Cascade, Restrict |
| `TestUndropDatabase_*` | Basic |
| `TestUndropSchema_*` | Basic |

Target: ~45 test cases.

---

## Files changed

| File | Change |
|------|--------|
| `snowflake/ast/parsenodes.go` | Add 8 stmt nodes + helpers |
| `snowflake/ast/nodetags.go` | Add 8 NodeTag constants + String() cases |
| `snowflake/ast/walk_generated.go` | Regenerated by genwalker |
| `snowflake/parser/create_table.go` | Add DATABASE/SCHEMA cases to CREATE sub-dispatch |
| `snowflake/parser/parser.go` | Replace ALTER/DROP/UNDROP stubs with dispatchers |
| `snowflake/parser/database_schema.go` | New: all 8 statement parsers + helpers |
| `snowflake/parser/database_schema_test.go` | New: ~45 tests |
| `docs/superpowers/specs/2026-04-14-snowflake-db-schema-design.md` | This file |
| `docs/superpowers/plans/2026-04-14-snowflake-db-schema.md` | Plan |
