# Snowflake ALTER TABLE Design (T2.3)

## Overview

ALTER TABLE is one of Snowflake's largest DDL statements.  
This spec covers the full action set, splitting actions into:
- **Lint-critical**: ADD COLUMN, DROP COLUMN, RENAME TO, ADD CONSTRAINT, DROP CONSTRAINT, DROP PRIMARY KEY
- **Full-fidelity parse**: all others

## AST Design

### Top-level node

```go
type AlterTableStmt struct {
    IfExists bool
    Name     *ObjectName
    Actions  []*AlterTableAction
    Loc      Loc
}
```

A single ALTER TABLE statement may carry multiple comma-separated actions
(Snowflake allows this for ADD COLUMN, DROP COLUMN, and some others).

### Action discriminated union

```go
type AlterTableActionKind int

const (
    AlterTableRename               AlterTableActionKind = iota
    AlterTableSwapWith
    AlterTableAddColumn
    AlterTableDropColumn
    AlterTableRenameColumn
    AlterTableAlterColumn
    AlterTableAddConstraint
    AlterTableDropConstraint
    AlterTableRenameConstraint
    AlterTableClusterBy
    AlterTableDropClusterKey
    AlterTableRecluster
    AlterTableSuspendRecluster
    AlterTableResumeRecluster
    AlterTableSet
    AlterTableUnset
    AlterTableSetTag
    AlterTableUnsetTag
    AlterTableAddRowAccessPolicy
    AlterTableDropRowAccessPolicy
    AlterTableDropAllRowAccessPolicies
    AlterTableAddSearchOpt
    AlterTableDropSearchOpt
    AlterTableSetMaskingPolicy
    AlterTableUnsetMaskingPolicy
    AlterTableSetColumnTag
    AlterTableUnsetColumnTag
)
```

### Per-action payload (union-style fields)

Only the fields relevant to each Kind are populated. This avoids a large
proliferation of struct types while keeping full fidelity for lint-critical paths.

Key fields:
- `NewName *ObjectName` — RENAME TO, SWAP WITH target
- `Columns []*ColumnDef` — ADD COLUMN definitions (reuse T2.2)
- `DropColumnNames []Ident` — DROP COLUMN list
- `IfExists bool` / `IfNotExists bool` — guards on column ops
- `OldName Ident` + `NewColName Ident` — RENAME COLUMN old TO new
- `ColumnAlter *ColumnAlter` — ALTER/MODIFY COLUMN specification
- `Constraint *TableConstraint` — ADD CONSTRAINT (reuse T2.2)
- `ConstraintName Ident` — DROP CONSTRAINT by name
- `IsPrimaryKey bool` — DROP PRIMARY KEY (not named)
- `DropUnique bool` — DROP UNIQUE (not named)
- `DropForeignKey bool` — DROP FOREIGN KEY
- `Cascade bool`, `Restrict bool` — DROP CONSTRAINT options
- `ClusterBy []Node` — CLUSTER BY expressions
- `Linear bool` — CLUSTER BY LINEAR
- `Props []*TableProp` — SET properties
- `UnsetProps []string` — UNSET property names
- `Tags []*TagAssignment` — SET TAG assignments
- `UnsetTags []*ObjectName` — UNSET TAG names
- `PolicyName *ObjectName` — row access / masking policy
- `PolicyCols []Ident` — row access policy ON (cols)
- `MaskColumn Ident` — SET/UNSET MASKING POLICY column
- `SearchOptOn []Node` — ADD/DROP SEARCH OPTIMIZATION ON targets

### ColumnAlter

Captures ALTER/MODIFY COLUMN sub-actions:

```go
type ColumnAlter struct {
    Column Ident
    Kind   ColumnAlterKind
    // Per-kind fields:
    DataType      *TypeName
    DefaultExpr   Node
    Comment       *string
    MaskingPolicy *ObjectName
    Tags          []*TagAssignment
    UnsetTags     []*ObjectName
}
```

### TableProp

```go
type TableProp struct {
    Name  string
    Value string
}
```

## Grammar coverage

From legacy `alter_table` + `alter_table_alter_column`:

| Form | Kind |
|------|------|
| RENAME TO name | AlterTableRename |
| SWAP WITH name | AlterTableSwapWith |
| ADD [COLUMN] [IF NOT EXISTS] col_def [, ...] | AlterTableAddColumn |
| DROP [COLUMN] [IF EXISTS] col [, ...] | AlterTableDropColumn |
| RENAME COLUMN old TO new | AlterTableRenameColumn |
| ALTER/MODIFY [COLUMN] col [...] | AlterTableAlterColumn |
| ADD [CONSTRAINT name] PK/UK/FK | AlterTableAddConstraint |
| DROP CONSTRAINT name | AlterTableDropConstraint |
| DROP PRIMARY KEY | AlterTableDropConstraint (IsPrimaryKey=true) |
| RENAME CONSTRAINT old TO new | AlterTableRenameConstraint |
| CLUSTER BY [LINEAR] (exprs) | AlterTableClusterBy |
| DROP CLUSTERING KEY | AlterTableDropClusterKey |
| RECLUSTER [...] | AlterTableRecluster |
| SUSPEND RECLUSTER | AlterTableSuspendRecluster |
| RESUME RECLUSTER | AlterTableResumeRecluster |
| SET properties... | AlterTableSet |
| UNSET properties... | AlterTableUnset |
| SET TAG (...) | AlterTableSetTag |
| UNSET TAG (...) | AlterTableUnsetTag |
| ADD ROW ACCESS POLICY | AlterTableAddRowAccessPolicy |
| DROP ROW ACCESS POLICY | AlterTableDropRowAccessPolicy |
| DROP ALL ROW ACCESS POLICIES | AlterTableDropAllRowAccessPolicies |
| ADD SEARCH OPTIMIZATION [ON ...] | AlterTableAddSearchOpt |
| DROP SEARCH OPTIMIZATION [ON ...] | AlterTableDropSearchOpt |
| ALTER/MODIFY COLUMN col SET MASKING POLICY | AlterTableSetMaskingPolicy |
| ALTER/MODIFY COLUMN col UNSET MASKING POLICY | AlterTableUnsetMaskingPolicy |
| ALTER/MODIFY col SET TAG (...) | AlterTableSetColumnTag |
| ALTER/MODIFY col UNSET TAG (...) | AlterTableUnsetColumnTag |

## Walker integration

`walkChildren` must visit:
- `n.Name` (*ObjectName)
- For each action: `NewName`, `Columns` (each ColumnDef), `Constraint`, `PolicyName`
- `ColumnAlter.DataType`, `ColumnAlter.DefaultExpr`, `ColumnAlter.MaskingPolicy`

Since `AlterTableAction` is not a Node itself, walker traverses actions inline
from the `AlterTableStmt` case.

## Lint-critical contract

Downstream advisors will type-assert `*ast.AlterTableStmt` and inspect `Actions`:
- `AlterTableAddColumn` → `Columns[i].Name`, `Columns[i].DataType`
- `AlterTableDropColumn` → `DropColumnNames`
- `AlterTableRename` → `NewName`
- `AlterTableAddConstraint` → `Constraint.Type` (PK/FK/UNIQUE), `Constraint.Columns`
- `AlterTableDropConstraint` → `ConstraintName`, `IsPrimaryKey`
