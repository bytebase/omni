# T2.4 — Snowflake VIEW + MATERIALIZED VIEW DDL Design

## Scope

T2.4 adds CREATE VIEW, CREATE MATERIALIZED VIEW, ALTER VIEW, and ALTER MATERIALIZED VIEW
to the Snowflake parser. DROP VIEW and DROP MATERIALIZED VIEW are already handled by T2.5
(DropStmt) and must not be touched here.

## Source of truth

Legacy grammar: `/Users/h3n4l/OpenSource/parser/snowflake/SnowflakeParser.g4`
- `create_view` (line 2850)
- `create_materialized_view` (line 1883)
- `alter_view` (line 1438)
- `alter_materialized_view` (line 906)

Official Snowflake docs syntax matches the legacy grammar; the two sources agree.

## AST nodes

### ViewColumn

Not a Node (value struct embedded in CreateViewStmt / CreateMaterializedViewStmt).

```go
type ViewColumn struct {
    Name          Ident
    MaskingPolicy *ObjectName // WITH MASKING POLICY p; nil if absent
    MaskingUsing  []Ident     // USING (col, ...); nil if absent
    Tags          []*TagAssignment
    Comment       *string
    Loc           Loc
}
```

### CreateViewStmt (Node, tag T_CreateViewStmt)

```go
type CreateViewStmt struct {
    OrReplace   bool
    Secure      bool
    Recursive   bool
    IfNotExists bool
    Name        *ObjectName
    Columns     []*ViewColumn    // optional col list from ( column_list_with_comment )
    ViewCols    []*ViewColumn    // view_col* — column-level policy/tag bindings
    CopyGrants  bool
    Comment     *string
    Tags        []*TagAssignment // WITH TAG (...)
    RowPolicy   *RowAccessPolicy // WITH ROW ACCESS POLICY ...; nil if absent
    Query       Node             // the AS query
    Loc         Loc
}
```

### RowAccessPolicy (helper struct, not a Node)

```go
type RowAccessPolicy struct {
    PolicyName *ObjectName
    Columns    []Ident
}
```

### CreateMaterializedViewStmt (Node, tag T_CreateMaterializedViewStmt)

Same fields as CreateViewStmt minus Recursive, plus:

```go
    ClusterBy  []Node // CLUSTER BY (exprs); nil if absent
    Linear     bool   // CLUSTER BY LINEAR modifier
```

### AlterViewAction enum

```go
type AlterViewAction int

const (
    AlterViewRename                    AlterViewAction = iota
    AlterViewSetComment
    AlterViewUnsetComment
    AlterViewSetSecure
    AlterViewUnsetSecure
    AlterViewSetTag
    AlterViewUnsetTag
    AlterViewAddRowAccessPolicy
    AlterViewDropRowAccessPolicy
    AlterViewDropAllRowAccessPolicies
    AlterViewColumnSetMaskingPolicy
    AlterViewColumnUnsetMaskingPolicy
    AlterViewColumnSetTag
    AlterViewColumnUnsetTag
)
```

### AlterViewStmt (Node, tag T_AlterViewStmt)

```go
type AlterViewStmt struct {
    IfExists   bool
    Name       *ObjectName
    Action     AlterViewAction
    NewName    *ObjectName      // RENAME TO
    Comment    *string          // SET COMMENT
    Secure     bool             // SET SECURE (true) or UNSET SECURE (false)
    Tags       []*TagAssignment // SET TAG
    UnsetTags  []*ObjectName    // UNSET TAG
    PolicyName *ObjectName      // ADD/DROP ROW ACCESS POLICY
    PolicyCols []Ident          // ON (col, ...)
    Column     Ident            // ALTER COLUMN col
    MaskingPolicy *ObjectName   // SET MASKING POLICY p
    MaskingUsing  []Ident       // USING (col, ...)
    Loc        Loc
}
```

### AlterMaterializedViewAction enum

```go
type AlterMaterializedViewAction int

const (
    AlterMVRename           AlterMaterializedViewAction = iota
    AlterMVClusterBy
    AlterMVDropClusteringKey
    AlterMVSuspend
    AlterMVResume
    AlterMVSuspendRecluster
    AlterMVResumeRecluster
    AlterMVSetSecure
    AlterMVUnsetSecure
    AlterMVSetComment
    AlterMVUnsetComment
)
```

### AlterMaterializedViewStmt (Node, tag T_AlterMaterializedViewStmt)

```go
type AlterMaterializedViewStmt struct {
    Name      *ObjectName
    Action    AlterMaterializedViewAction
    NewName   *ObjectName      // RENAME TO
    ClusterBy []Node           // CLUSTER BY (exprs)
    Linear    bool
    Comment   *string
    Secure    bool
    Loc       Loc
}
```

Note: The legacy grammar for ALTER MATERIALIZED VIEW does NOT include IF EXISTS
(the grammar uses `id_` not `if_exists? object_name`), consistent with Snowflake docs
for this command.

## Parser dispatch

### parseCreateStmt (create_table.go)

Add cases for `kwSECURE`, `kwRECURSIVE`, `kwVIEW`, `kwMATERIALIZED` after the existing
temporary/transient checks. `SECURE` and `RECURSIVE` are view-only modifiers that appear
before `VIEW`.

```
CREATE [OR REPLACE] [SECURE] [RECURSIVE] VIEW ...
CREATE [OR REPLACE] [SECURE] MATERIALIZED VIEW ...
```

Decision tree:
1. After consuming OR REPLACE, check for SECURE (consume it if present).
2. Check for RECURSIVE (consume it if present) → must be VIEW.
3. Check for MATERIALIZED (consume it if present) → must be VIEW next.
4. Check for VIEW → parseCreateViewStmt.
5. Existing TABLE/DATABASE/SCHEMA cases.

### parseAlterStmt (database_schema.go)

Add cases for `kwVIEW` and `kwMATERIALIZED` in the switch.

## Column list parsing

The legacy grammar supports two column-list forms:

1. `( column_list_with_comment )` — column names with optional COMMENT per column.
   Parsed as ViewColumn with Name and Comment only.

2. `view_col*` — zero or more `column_name WITH MASKING POLICY p [WITH TAG (...)]`
   entries outside the parens. Parsed as ViewColumn with MaskingPolicy and/or Tags.

The parser checks for `(` to decide which form to parse. After the optional `( ... )`,
it may encounter `view_col*` entries (column name followed by WITH MASKING POLICY or
WITH TAG without a comma separator — each is a distinct col entry).

## Row access policy

`[WITH] ROW ACCESS POLICY policy_name ON (col1 [, col2 ...])` is stored in RowAccessPolicy.

## Walker

walkChildren must be extended (manually, since the generator runs at build time):
- CreateViewStmt: Walk Name, walk Query, walk each ViewCol's MaskingPolicy/Tags.
- CreateMaterializedViewStmt: same + ClusterBy exprs.
- AlterViewStmt: Walk Name, NewName, PolicyName, MaskingPolicy.
- AlterMaterializedViewStmt: Walk Name, NewName, ClusterBy exprs.

Due to the complexity of nested non-Node structs (ViewColumn, RowAccessPolicy),
the walker in walk_generated.go will walk the top-level *ObjectName and Node children
only; ViewColumn and RowAccessPolicy are treated as opaque value structs.

## Test coverage

Tests in `snowflake/parser/view_test.go`:
- CREATE VIEW basic, OR REPLACE, SECURE, RECURSIVE, IF NOT EXISTS
- CREATE VIEW with column list + comments
- CREATE VIEW with COPY GRANTS, COMMENT, WITH TAG, WITH ROW ACCESS POLICY
- CREATE VIEW AS SELECT / WITH CTE query
- CREATE MATERIALIZED VIEW basic + CLUSTER BY
- ALTER VIEW RENAME TO
- ALTER VIEW SET/UNSET COMMENT
- ALTER VIEW SET/UNSET SECURE
- ALTER VIEW SET/UNSET TAG
- ALTER VIEW ADD/DROP ROW ACCESS POLICY
- ALTER VIEW ALTER COLUMN SET/UNSET MASKING POLICY
- ALTER VIEW ALTER COLUMN SET/UNSET TAG
- ALTER MATERIALIZED VIEW RENAME, CLUSTER BY, DROP CLUSTERING KEY
- ALTER MATERIALIZED VIEW SUSPEND, RESUME, SUSPEND RECLUSTER, RESUME RECLUSTER
- ALTER MATERIALIZED VIEW SET/UNSET SECURE, SET/UNSET COMMENT
