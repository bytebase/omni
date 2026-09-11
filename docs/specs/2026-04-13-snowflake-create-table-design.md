# T2.2: Snowflake CREATE TABLE Design

**Date:** 2026-04-13
**Node:** T2.2 (DDL: CREATE TABLE)
**Branch:** `feat/snowflake/create-table`
**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/create-table`

## Overview

T2.2 implements CREATE TABLE parsing for the Snowflake engine. This is the most complex DDL statement and the gateway to bytebase's lint rules that depend on table structure (`column_no_null`, `column_require`, `column_maximum_varchar_length`, `table_require_pk`, `table_no_foreign_key`, `naming_table`, `naming_table_no_keyword`, `naming_identifier_case`, `naming_identifier_no_keyword`, `migration_compatibility`).

It partially replaces F4's `unsupported("CREATE")` stub by adding a sub-dispatch after the CREATE keyword. Only CREATE TABLE is implemented; all other CREATE forms remain unsupported. The sub-dispatch pattern is designed so future nodes (T2.3 CREATE VIEW, etc.) can plug in by adding cases.

## Scope

### Table creation forms

| Form | Example | Legacy rule |
|------|---------|-------------|
| Standard | `CREATE TABLE t (col INT, ...)` | `create_table` |
| CTAS | `CREATE TABLE t AS SELECT ...` | `create_table_as_select` |
| CTAS with columns | `CREATE TABLE t (col INT) AS SELECT ...` | `create_table_as_select` |
| LIKE | `CREATE TABLE t LIKE source` | `create_table_like` |
| CLONE | `CREATE TABLE t CLONE source [AT\|BEFORE (...)]` | `create_object_clone` |

### Modifiers

- `OR REPLACE` — recreate if exists
- `TRANSIENT` — transient table (no Fail-safe)
- `TEMPORARY` / `TEMP` — session-scoped temporary table
- `VOLATILE` — synonym for TEMPORARY
- `LOCAL` / `GLOBAL` — optional prefix for TEMPORARY (consumed, no semantic difference in Snowflake)
- `IF NOT EXISTS` — skip if already exists

### Column features

| Feature | AST field | Lint-relevant |
|---------|-----------|---------------|
| Name + data type | `ColumnDef.Name`, `.DataType` | Yes (naming, varchar length) |
| `NOT NULL` / `NULL` | `.NotNull`, `.Nullable` | Yes (column_no_null) |
| `DEFAULT expr` | `.Default` | No |
| `IDENTITY`/`AUTOINCREMENT` | `.Identity` | No |
| `COLLATE 'name'` | `.Collate` | No |
| `WITH MASKING POLICY name` | `.MaskingPolicy` | No |
| Inline constraint (PK/FK/UNIQUE) | `.InlineConstraint` | Yes (PK/FK detection) |
| `COMMENT 'text'` | `.Comment` | No |
| `WITH TAG (...)` | `.Tags` | No |
| Virtual column `AS (expr)` | `.VirtualExpr` | No |

### Table-level constraints

| Constraint | Stored in AST | Lint-relevant |
|------------|---------------|---------------|
| `PRIMARY KEY (cols)` | Yes | Yes (table_require_pk) |
| `FOREIGN KEY (cols) REFERENCES ...` | Yes | Yes (table_no_foreign_key) |
| `UNIQUE (cols)` | Yes | No |
| Constraint properties (ENFORCED, DEFERRABLE, RELY, etc.) | Consumed, not stored | No |
| FK actions (ON DELETE/UPDATE) | Stored in ForeignKeyRef | No |
| FK MATCH (FULL/PARTIAL/SIMPLE) | Stored in ForeignKeyRef | No |

### Table properties

| Property | Stored in AST | Lint-relevant |
|----------|---------------|---------------|
| `CLUSTER BY [LINEAR] (exprs)` | Yes (ClusterBy, Linear) | No |
| `COPY GRANTS` | Yes (CopyGrants) | No |
| `COMMENT = 'text'` | Yes (Comment) | No |
| `WITH TAG (...)` | Yes (Tags) | No |
| `DATA_RETENTION_TIME_IN_DAYS = n` | Consumed, not stored | No |
| `CHANGE_TRACKING = TRUE\|FALSE` | Consumed, not stored | No |
| `DEFAULT_DDL_COLLATION_ = 'name'` | Consumed, not stored | No |
| `WITH ROW ACCESS POLICY ...` | Consumed, not stored | No |
| `STAGE_FILE_FORMAT = (...)` | Consumed, not stored | No |
| `STAGE_COPY_OPTIONS = (...)` | Consumed, not stored | No |

## AST Types

### Node types (3 new, with tags)

```go
// CreateTableStmt represents CREATE [OR REPLACE] [TRANSIENT|TEMPORARY|VOLATILE] TABLE ...
type CreateTableStmt struct {
    OrReplace   bool
    Transient   bool
    Temporary   bool
    Volatile    bool
    IfNotExists bool
    Name        *ObjectName
    Columns     []*ColumnDef
    Constraints []*TableConstraint
    ClusterBy   []Node             // CLUSTER BY expressions; nil if absent
    Linear      bool               // CLUSTER BY LINEAR modifier
    Comment     *string            // COMMENT = 'text'; nil if absent
    CopyGrants  bool
    Tags        []*TagAssignment   // WITH TAG (...); nil if absent
    AsSelect    Node               // CREATE TABLE ... AS SELECT; nil if absent
    Like        *ObjectName        // CREATE TABLE ... LIKE source; nil if absent
    Clone       *CloneSource       // CREATE TABLE ... CLONE source; nil if absent
    Loc         Loc
}

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
    Name             Ident
    DataType         *TypeName          // nil for virtual columns without explicit type
    Default          Node               // DEFAULT expr; nil if absent
    NotNull          bool
    Nullable         bool               // explicit NULL
    Identity         *IdentitySpec      // IDENTITY/AUTOINCREMENT; nil if absent
    Collate          string             // COLLATE 'name'; empty if absent
    MaskingPolicy    *ObjectName        // WITH MASKING POLICY name; nil if absent
    InlineConstraint *InlineConstraint  // inline PK/FK/UNIQUE; nil if absent
    Comment          *string            // COMMENT 'text'; nil if absent
    Tags             []*TagAssignment   // WITH TAG (...); nil if absent
    VirtualExpr      Node               // AS (expr); nil if absent
    Loc              Loc
}

// TableConstraint represents a table-level constraint (out-of-line).
type TableConstraint struct {
    Type       ConstraintType  // ConstrPrimaryKey/ConstrForeignKey/ConstrUnique
    Name       Ident           // CONSTRAINT name; zero if unnamed
    Columns    []Ident         // constrained column names
    References *ForeignKeyRef  // FK only; nil otherwise
    Comment    *string         // inline COMMENT 'text'; nil if absent
    Loc        Loc
}
```

### Enum

```go
type ConstraintType int

const (
    ConstrPrimaryKey ConstraintType = iota
    ConstrForeignKey
    ConstrUnique
)
```

### Helper structs (not Nodes, no tags)

```go
// InlineConstraint represents a column-level constraint.
type InlineConstraint struct {
    Type       ConstraintType
    Name       Ident           // CONSTRAINT name; zero if unnamed
    References *ForeignKeyRef  // for FK; nil otherwise
    Loc        Loc
}

// ForeignKeyRef holds REFERENCES clause details.
type ForeignKeyRef struct {
    Table    *ObjectName
    Columns  []Ident
    OnDelete ReferenceAction
    OnUpdate ReferenceAction
    Match    string  // "FULL"/"PARTIAL"/"SIMPLE"; empty if absent
}

// ReferenceAction enumerates FK referential actions.
type ReferenceAction int

const (
    RefActNone       ReferenceAction = iota // not specified
    RefActCascade                           // CASCADE
    RefActSetNull                           // SET NULL
    RefActSetDefault                        // SET DEFAULT
    RefActRestrict                          // RESTRICT
    RefActNoAction                          // NO ACTION
)

// IdentitySpec holds IDENTITY/AUTOINCREMENT configuration.
type IdentitySpec struct {
    Start     *int64  // START WITH value; nil if default
    Increment *int64  // INCREMENT BY value; nil if default
    Order     *bool   // true=ORDER, false=NOORDER, nil=unspecified
}

// TagAssignment is a single TAG name = 'value' pair.
type TagAssignment struct {
    Name  *ObjectName
    Value string
}

// CloneSource holds CLONE source with optional time travel.
type CloneSource struct {
    Source   *ObjectName
    AtBefore string  // "AT" or "BEFORE"; empty if no time travel
    Kind     string  // "TIMESTAMP"/"OFFSET"/"STATEMENT"
    Value    string  // the time travel value
}
```

### Node tags (3 new)

```
T_CreateTableStmt
T_ColumnDef
T_TableConstraint
```

## Parser Structure

### Dispatch change (parser.go)

```go
case kwCREATE:
    return p.parseCreateStmt()  // replaces p.unsupported("CREATE")
```

### New file: create_table.go

```
parseCreateStmt()              -- CREATE dispatch: OR REPLACE, table type, sub-switch TABLE/default
parseCreateTableStmt(...)      -- IF NOT EXISTS, name, branch: LIKE/CLONE/AS/(columns)/columns+AS
parseColumnDeclItems()         -- ( columnDef|tableConstraint , ... )
parseColumnDef()               -- name type [options...]
parseColumnOptions(col)        -- loop: NOT NULL/NULL/DEFAULT/IDENTITY/COLLATE/MASKING/CONSTRAINT/COMMENT/TAG/AS
parseInlineConstraint()        -- [CONSTRAINT name] PK/UNIQUE/FK REFERENCES
parseOutOfLineConstraint()     -- [CONSTRAINT name] PK(cols)/UNIQUE(cols)/FK(cols) REFERENCES
parseForeignKeyRef()           -- REFERENCES table(cols) [MATCH] [ON DELETE/UPDATE]
parseConstraintProperties()    -- ENFORCED/DEFERRABLE/RELY etc. (consumed, discarded)
parseIdentitySpec()            -- (start, incr) / START WITH / INCREMENT BY / ORDER/NOORDER
parseTagAssignments()          -- ( name = 'val' , ... )
parseCloneSource()             -- CLONE name [AT|BEFORE (TIMESTAMP => / OFFSET => / STATEMENT =>)]
parseTableProperties(stmt)     -- CLUSTER BY, COMMENT, COPY GRANTS, tags, retention, etc.
```

### Disambiguation logic

After consuming `CREATE [OR REPLACE] [table_type] TABLE [IF NOT EXISTS] name`:

1. If `LIKE` -> parse LIKE source, return
2. If `CLONE` -> parse clone source with optional time travel, return
3. If `AS` -> parse CTAS query, return
4. If `(` -> parse column/constraint list
5. After column list: parse table properties, then optionally `AS SELECT` (CTAS with columns)

For `(` after table name, the column list is always a column/constraint list (not a subquery) because CTAS without `AS` keyword is not valid Snowflake syntax.

### Column vs. constraint disambiguation inside parentheses

At each item in the comma-separated list:
- `CONSTRAINT` / `PRIMARY` / `UNIQUE` / `FOREIGN` -> table-level constraint
- Otherwise -> column definition (starts with identifier)

## Walker

`genwalker` auto-generates cases for the 3 new Node types:

- **CreateTableStmt**: walks `Name` (*ObjectName), `Columns` ([]*ColumnDef), `Constraints` ([]*TableConstraint), `ClusterBy` ([]Node), `AsSelect` (Node), `Like` (*ObjectName)
- **ColumnDef**: walks `DataType` (*TypeName), `Default` (Node), `MaskingPolicy` (*ObjectName), `VirtualExpr` (Node)
- **TableConstraint**: no Node children to walk

Helper structs (InlineConstraint, ForeignKeyRef, etc.) are not Nodes and are not walked. Lint rules that need FK details find the TableConstraint/ColumnDef and inspect fields directly.

## Test Strategy

### Unit tests (create_table_test.go)

**Table creation forms:**
- Standard CREATE TABLE with columns
- CTAS (CREATE TABLE ... AS SELECT)
- CTAS with columns (CREATE TABLE (col INT) AS SELECT)
- LIKE (CREATE TABLE ... LIKE source)
- CLONE (CREATE TABLE ... CLONE source)
- CLONE with AT/BEFORE time travel

**Modifiers:**
- OR REPLACE
- TRANSIENT, TEMPORARY, VOLATILE
- IF NOT EXISTS
- Combinations (OR REPLACE + TRANSIENT + IF NOT EXISTS)

**Column features:**
- Basic column (name + type)
- NOT NULL / NULL
- DEFAULT expr
- IDENTITY/AUTOINCREMENT with (start, incr)
- IDENTITY with START WITH / INCREMENT BY
- COLLATE
- WITH MASKING POLICY
- Inline PRIMARY KEY, UNIQUE, FOREIGN KEY REFERENCES
- COMMENT
- WITH TAG
- Virtual column AS (expr)

**Table constraints:**
- PRIMARY KEY (col1, col2)
- UNIQUE (col)
- FOREIGN KEY (col) REFERENCES other(col) with ON DELETE/UPDATE
- Named constraints (CONSTRAINT name ...)
- Mixed columns and constraints

**Table properties:**
- CLUSTER BY (exprs)
- CLUSTER BY LINEAR
- COPY GRANTS
- COMMENT = 'text'
- WITH TAG
- DATA_RETENTION_TIME_IN_DAYS (consumed without error)
- CHANGE_TRACKING (consumed without error)

**Legacy corpus:**
- Any CREATE TABLE statements in corpus A/B should parse without errors

### Test pattern

Each test case: SQL input -> parse -> verify AST fields. Follow the same `[]struct{name, sql, check}` pattern used in select_test.go and expr_test.go.

## Estimated Size

~1,200-1,500 LOC:
- `snowflake/ast/parsenodes.go`: +200 LOC (new types)
- `snowflake/ast/nodetags.go`: +15 LOC (3 tags + String cases)
- `snowflake/parser/create_table.go`: ~800 LOC (new file)
- `snowflake/parser/parser.go`: ~5 LOC (dispatch change)
- `snowflake/parser/create_table_test.go`: ~400 LOC (new file)
- `snowflake/ast/walk_generated.go`: regenerated

## Dependencies

- Reuses `parseDataType` from T1.2 (datatypes.go)
- Reuses `parseExpr` from T1.3 (expr.go)
- Reuses `parseQueryExpr` from T1.4/T1.7 (select.go) for CTAS
- Reuses `parseObjectName`, `parseIdent` from T1.1 (identifiers.go)
- All Tier 1 nodes are merged to main
