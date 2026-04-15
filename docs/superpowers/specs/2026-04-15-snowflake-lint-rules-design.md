# Spec: Snowflake Lint Rules (T2.7)

## Context

T2.7 adds 14 production-quality lint rules to the Snowflake advisor framework (T2.6). All AST dependencies are available: CREATE TABLE (T2.2), ALTER TABLE (T2.3), DROP/UNDROP (T2.5), advisor framework (T2.6), and SELECT core (T1.4).

## Rules

### Column rules (3)

**column_no_null** — Every column in a CREATE TABLE must have a NOT NULL constraint (either `col TYPE NOT NULL` or no explicit `NULL`). Virtual columns (WITH AS expr) are exempt because they cannot have NOT NULL. Fires on ColumnDef where `NotNull == false && Nullable == false` is NOT sufficient — we want explicit NOT NULL. Design decision: fire when `NotNull == false` (i.e., NOT NULL absent) AND `VirtualExpr == nil`.

**column_require** — A configured list of column names must appear in every CREATE TABLE. Config = slice of required column names (case-insensitive match via Snowflake normalize). Default = empty (no-op).

**column_maximum_varchar_length** — VARCHAR/STRING/TEXT/CHAR columns in CREATE TABLE whose declared length > configured max. TypeKind TypeVarchar and TypeChar with `Params[0] > max`. For types with no Params (unbounded), skip. Default max = 1024.

### Table rules (2)

**table_require_pk** — Every CREATE TABLE (non-TEMPORARY, non-TRANSIENT, non-VOLATILE, and not AS SELECT / LIKE / CLONE) must define a PRIMARY KEY. Check: any ColumnDef with inline `ConstrPrimaryKey` OR any TableConstraint with `ConstrPrimaryKey`. Except for derived tables (AsSelect/Like/Clone non-nil) and temp tables.

**table_no_foreign_key** — No FOREIGN KEY constraints. Checks: ColumnDef.InlineConstraint.Type == ConstrForeignKey OR TableConstraint.Type == ConstrForeignKey.

### Naming rules (4)

**naming_table** — Table name (last part of ObjectName) in CREATE TABLE must match configured regex. Default: `^[a-z][a-z0-9_]{0,63}$`. Match against unquoted name as-is (preserving source case, since lint runs before resolution).

**naming_table_no_keyword** — Table name (last part of ObjectName) in CREATE TABLE must not be a reserved keyword. Uses `parser.IsReservedKeyword`. Quoted names are exempt (they can always be used as identifiers).

**naming_identifier_case** — All identifiers (column names, table names, aliases) must match a configured case convention. Default: LOWER. Conventions: LOWER (all lowercase), UPPER (all uppercase), CAMEL (camelCase). Applied to unquoted identifiers only. Visits: ColumnDef.Name, CreateTableStmt.Name.Name, TableRef.Alias, SelectTarget.Alias, UpdateStmt.Target.Name, DeleteStmt.Target.Name.

**naming_identifier_no_keyword** — No unquoted identifier may be a reserved keyword. Applies to column names, table aliases, CTE names. Visits ColumnDef.Name (in CreateTableStmt), TableRef.Alias, SelectTarget.Alias, CTE.Name.

### Query rules (3)

**select_no_select_all** — SELECT * (bare star) not allowed. Identical trigger to the example NoSelectStarRule but with rule ID `snowflake.select.no-select-all` and Severity ERROR. Note: the example rule has a different ID and Severity, so these are distinct.

**where_require_select** — SELECT without a WHERE clause. Only fires for SelectStmt with `From != nil` (i.e., not `SELECT 1`). CTE body SELECTs also fire.

**where_require_update_delete** — UPDATE or DELETE without a WHERE clause.

### Migration rules (2)

**migration_compatibility** — Forbids destructive DDL: DROP TABLE, DROP COLUMN (in ALTER TABLE), ALTER TABLE DROP constraint. Also fires on DropStmt with Kind==DropTable; AlterTableAction with Kind==AlterTableDropColumn.

**table_drop_naming_convention** — When DROP TABLE is used, the table name must match a configured regex (allows "archive" naming patterns). Default: `^_deleted$|^deleted_` i.e. name starts with `deleted_` OR equals `_deleted`. Match against the unquoted last-part name.

## AST Field Mapping

| Rule | Trigger Node | Key Fields |
|------|-------------|------------|
| column_no_null | *CreateTableStmt | .Columns[i].NotNull, .VirtualExpr |
| column_require | *CreateTableStmt | .Columns[i].Name |
| column_maximum_varchar_length | *CreateTableStmt | .Columns[i].DataType.Kind, .Params[0] |
| table_require_pk | *CreateTableStmt | .Columns[i].InlineConstraint, .Constraints[i].Type |
| table_no_foreign_key | *CreateTableStmt | .Columns[i].InlineConstraint.Type, .Constraints[i].Type |
| naming_table | *CreateTableStmt | .Name.Name |
| naming_table_no_keyword | *CreateTableStmt | .Name.Name.Quoted, IsReservedKeyword |
| naming_identifier_case | *CreateTableStmt, *SelectStmt | .Columns[i].Name, .Name.Name, .Targets[i].Alias |
| naming_identifier_no_keyword | *CreateTableStmt, *SelectStmt | .Columns[i].Name, .Targets[i].Alias, .With[i].Name |
| select_no_select_all | *SelectStmt | .Targets[i].Star, StarExpr.Qualifier |
| where_require_select | *SelectStmt | .Where, .From |
| where_require_update_delete | *UpdateStmt, *DeleteStmt | .Where |
| migration_compatibility | *DropStmt, *AlterTableStmt | .Kind, .Actions[i].Kind |
| table_drop_naming_convention | *DropStmt | .Kind==DropTable, .Name.Name |

## Design Decisions

1. **naming_identifier_case**: Applied to unquoted identifiers only. Quoted identifiers are intentionally case-sensitive and must not be folded.
2. **table_require_pk**: Temporary/transient/volatile tables are exempt (they're often used as scratch space).
3. **column_no_null**: Virtual columns (computed AS expr) are always exempt — Snowflake does not allow NOT NULL on virtual columns.
4. **select_no_select_all vs example NoSelectStarRule**: The production rule has ID `snowflake.select.no-select-all` (matching the bytebase naming convention) and Severity ERROR (stricter).
5. **where_require_select**: `SELECT 1` (no FROM clause) is exempt — table-less selects are always safe.
6. **migration_compatibility**: Only fires for structural drops: DROP TABLE, DROP COLUMN. Does NOT fire for DROP VIEW, DROP INDEX, etc.
