package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func testParseAlterTable(input string) (*ast.AlterTableStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.AlterTableStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not an AlterTableStmt"})
	}
	return stmt, result.Errors
}

func requireAlterTable(t *testing.T, input string) *ast.AlterTableStmt {
	t.Helper()
	stmt, errs := testParseAlterTable(input)
	if len(errs) > 0 {
		t.Fatalf("parse errors for %q: %v", input, errs)
	}
	if stmt == nil {
		t.Fatalf("expected AlterTableStmt for %q, got nil", input)
	}
	return stmt
}

func requireAction(t *testing.T, stmt *ast.AlterTableStmt, idx int) *ast.AlterTableAction {
	t.Helper()
	if idx >= len(stmt.Actions) {
		t.Fatalf("expected action[%d], only %d actions", idx, len(stmt.Actions))
	}
	return stmt.Actions[idx]
}

// ---------------------------------------------------------------------------
// RENAME TO (lint-critical: migration_compatibility)
// ---------------------------------------------------------------------------

func TestAlterTable_RenameTo(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t1 RENAME TO t2")
	if stmt.Name.Normalize() != "T1" {
		t.Errorf("name = %q, want T1", stmt.Name.Normalize())
	}
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableRename {
		t.Errorf("kind = %v, want AlterTableRename", a.Kind)
	}
	if a.NewName == nil || a.NewName.Normalize() != "T2" {
		t.Errorf("NewName = %v, want T2", a.NewName)
	}
}

func TestAlterTable_RenameToIfExists(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE IF EXISTS mydb.myschema.t1 RENAME TO t2")
	if !stmt.IfExists {
		t.Error("expected IfExists = true")
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.T1" {
		t.Errorf("name = %q", stmt.Name.Normalize())
	}
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableRename {
		t.Errorf("kind = %v, want AlterTableRename", a.Kind)
	}
}

// ---------------------------------------------------------------------------
// SWAP WITH
// ---------------------------------------------------------------------------

func TestAlterTable_SwapWith(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE orders SWAP WITH orders_backup")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSwapWith {
		t.Errorf("kind = %v, want AlterTableSwapWith", a.Kind)
	}
	if a.NewName == nil || a.NewName.Normalize() != "ORDERS_BACKUP" {
		t.Errorf("NewName = %v", a.NewName)
	}
}

// ---------------------------------------------------------------------------
// ADD COLUMN (lint-critical: column detection)
// ---------------------------------------------------------------------------

func TestAlterTable_AddColumn_Simple(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD COLUMN age INT")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddColumn {
		t.Fatalf("kind = %v, want AlterTableAddColumn", a.Kind)
	}
	if len(a.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(a.Columns))
	}
	col := a.Columns[0]
	if col.Name.Normalize() != "AGE" {
		t.Errorf("col name = %q, want AGE", col.Name.Normalize())
	}
	if col.DataType == nil || col.DataType.Kind != ast.TypeInt {
		t.Errorf("col datatype = %v, want Int", col.DataType)
	}
}

func TestAlterTable_AddColumn_NoColumnKeyword(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD name VARCHAR(100)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddColumn {
		t.Fatalf("kind = %v, want AlterTableAddColumn", a.Kind)
	}
	if len(a.Columns) != 1 || a.Columns[0].Name.Normalize() != "NAME" {
		t.Errorf("unexpected columns: %v", a.Columns)
	}
}

func TestAlterTable_AddColumn_IfNotExists(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD COLUMN IF NOT EXISTS email VARCHAR")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddColumn {
		t.Fatalf("kind = %v, want AlterTableAddColumn", a.Kind)
	}
	if !a.IfNotExists {
		t.Error("expected IfNotExists = true")
	}
	if len(a.Columns) != 1 || a.Columns[0].Name.Normalize() != "EMAIL" {
		t.Errorf("unexpected columns: %v", a.Columns)
	}
}

func TestAlterTable_AddColumn_Multiple(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD COLUMN col1 INT, col2 VARCHAR(50), col3 BOOLEAN")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddColumn {
		t.Fatalf("kind = %v, want AlterTableAddColumn", a.Kind)
	}
	if len(a.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(a.Columns))
	}
	if a.Columns[0].Name.Normalize() != "COL1" {
		t.Errorf("col[0] name = %q", a.Columns[0].Name.Normalize())
	}
	if a.Columns[2].Name.Normalize() != "COL3" {
		t.Errorf("col[2] name = %q", a.Columns[2].Name.Normalize())
	}
}

func TestAlterTable_AddColumn_WithOptions(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD COLUMN id INT NOT NULL DEFAULT 0 COMMENT 'primary key'")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddColumn {
		t.Fatalf("kind = %v, want AlterTableAddColumn", a.Kind)
	}
	col := a.Columns[0]
	if !col.NotNull {
		t.Error("expected NotNull = true")
	}
	if col.Default == nil {
		t.Error("expected Default != nil")
	}
	if col.Comment == nil || *col.Comment != "primary key" {
		t.Errorf("comment = %v", col.Comment)
	}
}

// ---------------------------------------------------------------------------
// DROP COLUMN (lint-critical: naming convention check)
// ---------------------------------------------------------------------------

func TestAlterTable_DropColumn_Simple(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP COLUMN age")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropColumn {
		t.Fatalf("kind = %v, want AlterTableDropColumn", a.Kind)
	}
	if len(a.DropColumnNames) != 1 || a.DropColumnNames[0].Normalize() != "AGE" {
		t.Errorf("DropColumnNames = %v", a.DropColumnNames)
	}
}

func TestAlterTable_DropColumn_IfExists(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP COLUMN IF EXISTS salary")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropColumn {
		t.Fatalf("kind = %v, want AlterTableDropColumn", a.Kind)
	}
	if !a.IfExists {
		t.Error("expected IfExists = true")
	}
	if len(a.DropColumnNames) != 1 || a.DropColumnNames[0].Normalize() != "SALARY" {
		t.Errorf("DropColumnNames = %v", a.DropColumnNames)
	}
}

func TestAlterTable_DropColumn_Multiple(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP COLUMN col1, col2, col3")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropColumn {
		t.Fatalf("kind = %v, want AlterTableDropColumn", a.Kind)
	}
	if len(a.DropColumnNames) != 3 {
		t.Fatalf("expected 3 drop cols, got %d", len(a.DropColumnNames))
	}
}

func TestAlterTable_DropColumn_NoColumnKeyword(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP salary")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropColumn {
		t.Fatalf("kind = %v, want AlterTableDropColumn", a.Kind)
	}
	if len(a.DropColumnNames) != 1 || a.DropColumnNames[0].Normalize() != "SALARY" {
		t.Errorf("DropColumnNames = %v", a.DropColumnNames)
	}
}

// ---------------------------------------------------------------------------
// RENAME COLUMN
// ---------------------------------------------------------------------------

func TestAlterTable_RenameColumn(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t RENAME COLUMN old_name TO new_name")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableRenameColumn {
		t.Fatalf("kind = %v, want AlterTableRenameColumn", a.Kind)
	}
	if a.OldName.Normalize() != "OLD_NAME" {
		t.Errorf("OldName = %q", a.OldName.Normalize())
	}
	if a.NewColName.Normalize() != "NEW_NAME" {
		t.Errorf("NewColName = %q", a.NewColName.Normalize())
	}
}

// ---------------------------------------------------------------------------
// ADD CONSTRAINT (lint-critical: PK/FK/UK detection)
// ---------------------------------------------------------------------------

func TestAlterTable_AddPrimaryKey(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD PRIMARY KEY (id)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddConstraint {
		t.Fatalf("kind = %v, want AlterTableAddConstraint", a.Kind)
	}
	if a.Constraint == nil {
		t.Fatal("Constraint is nil")
	}
	if a.Constraint.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint type = %v, want PRIMARY KEY", a.Constraint.Type)
	}
	if len(a.Constraint.Columns) != 1 || a.Constraint.Columns[0].Normalize() != "ID" {
		t.Errorf("columns = %v", a.Constraint.Columns)
	}
}

func TestAlterTable_AddPrimaryKeyNamed(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD CONSTRAINT pk_t PRIMARY KEY (id, code)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddConstraint {
		t.Fatalf("kind = %v, want AlterTableAddConstraint", a.Kind)
	}
	if a.Constraint.Name.Normalize() != "PK_T" {
		t.Errorf("constraint name = %q", a.Constraint.Name.Normalize())
	}
	if a.Constraint.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint type = %v", a.Constraint.Type)
	}
	if len(a.Constraint.Columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(a.Constraint.Columns))
	}
}

func TestAlterTable_AddUniqueConstraint(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD UNIQUE (email)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddConstraint {
		t.Fatalf("kind = %v, want AlterTableAddConstraint", a.Kind)
	}
	if a.Constraint.Type != ast.ConstrUnique {
		t.Errorf("constraint type = %v, want UNIQUE", a.Constraint.Type)
	}
}

func TestAlterTable_AddForeignKey(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE orders ADD FOREIGN KEY (customer_id) REFERENCES customers (id)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddConstraint {
		t.Fatalf("kind = %v, want AlterTableAddConstraint", a.Kind)
	}
	if a.Constraint.Type != ast.ConstrForeignKey {
		t.Errorf("constraint type = %v, want FOREIGN KEY", a.Constraint.Type)
	}
	if a.Constraint.References == nil {
		t.Fatal("References is nil")
	}
	if a.Constraint.References.Table.Normalize() != "CUSTOMERS" {
		t.Errorf("references table = %q", a.Constraint.References.Table.Normalize())
	}
}

// ---------------------------------------------------------------------------
// DROP CONSTRAINT (lint-critical)
// ---------------------------------------------------------------------------

func TestAlterTable_DropConstraintNamed(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP CONSTRAINT pk_t")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropConstraint {
		t.Fatalf("kind = %v, want AlterTableDropConstraint", a.Kind)
	}
	if a.ConstraintName.Normalize() != "PK_T" {
		t.Errorf("ConstraintName = %q", a.ConstraintName.Normalize())
	}
}

func TestAlterTable_DropConstraintCascade(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP CONSTRAINT fk_c CASCADE")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropConstraint {
		t.Fatalf("kind = %v, want AlterTableDropConstraint", a.Kind)
	}
	if !a.Cascade {
		t.Error("expected Cascade = true")
	}
}

func TestAlterTable_DropPrimaryKey(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP PRIMARY KEY")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropConstraint {
		t.Fatalf("kind = %v, want AlterTableDropConstraint", a.Kind)
	}
	if !a.IsPrimaryKey {
		t.Error("expected IsPrimaryKey = true")
	}
}

// ---------------------------------------------------------------------------
// RENAME CONSTRAINT
// ---------------------------------------------------------------------------

func TestAlterTable_RenameConstraint(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t RENAME CONSTRAINT old_pk TO new_pk")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableRenameConstraint {
		t.Fatalf("kind = %v, want AlterTableRenameConstraint", a.Kind)
	}
	if a.ConstraintName.Normalize() != "OLD_PK" {
		t.Errorf("ConstraintName = %q", a.ConstraintName.Normalize())
	}
	if a.NewConstraintName.Normalize() != "NEW_PK" {
		t.Errorf("NewConstraintName = %q", a.NewConstraintName.Normalize())
	}
}

// ---------------------------------------------------------------------------
// CLUSTER BY / DROP CLUSTERING KEY
// ---------------------------------------------------------------------------

func TestAlterTable_ClusterBy(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t CLUSTER BY (col1, col2)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableClusterBy {
		t.Fatalf("kind = %v, want AlterTableClusterBy", a.Kind)
	}
	if len(a.ClusterBy) != 2 {
		t.Errorf("ClusterBy exprs = %d, want 2", len(a.ClusterBy))
	}
}

func TestAlterTable_ClusterByLinear(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t CLUSTER BY LINEAR (col1)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableClusterBy {
		t.Fatalf("kind = %v, want AlterTableClusterBy", a.Kind)
	}
	if !a.Linear {
		t.Error("expected Linear = true")
	}
}

func TestAlterTable_DropClusteringKey(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP CLUSTERING KEY")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropClusterKey {
		t.Fatalf("kind = %v, want AlterTableDropClusterKey", a.Kind)
	}
}

func TestAlterTable_SuspendRecluster(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SUSPEND RECLUSTER")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSuspendRecluster {
		t.Fatalf("kind = %v, want AlterTableSuspendRecluster", a.Kind)
	}
}

func TestAlterTable_ResumeRecluster(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t RESUME RECLUSTER")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableResumeRecluster {
		t.Fatalf("kind = %v, want AlterTableResumeRecluster", a.Kind)
	}
}

// ---------------------------------------------------------------------------
// ALTER/MODIFY COLUMN variants
// ---------------------------------------------------------------------------

func TestAlterTable_AlterColumn_SetDataType(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c SET DATA TYPE VARCHAR(200)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAlterColumn {
		t.Fatalf("kind = %v, want AlterTableAlterColumn", a.Kind)
	}
	if len(a.ColumnAlters) != 1 {
		t.Fatalf("expected 1 ColumnAlter, got %d", len(a.ColumnAlters))
	}
	ca := a.ColumnAlters[0]
	if ca.Column.Normalize() != "C" {
		t.Errorf("column = %q", ca.Column.Normalize())
	}
	if ca.Kind != ast.ColumnAlterSetDataType {
		t.Errorf("alter kind = %v", ca.Kind)
	}
	if ca.DataType == nil || ca.DataType.Kind != ast.TypeVarchar {
		t.Errorf("datatype = %v", ca.DataType)
	}
}

func TestAlterTable_AlterColumn_TypeShorthand(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c TYPE NUMBER(10,2)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAlterColumn {
		t.Fatalf("kind = %v, want AlterTableAlterColumn", a.Kind)
	}
	ca := a.ColumnAlters[0]
	if ca.Kind != ast.ColumnAlterSetDataType {
		t.Errorf("alter kind = %v", ca.Kind)
	}
}

func TestAlterTable_ModifyColumn_DropDefault(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t MODIFY COLUMN c DROP DEFAULT")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAlterColumn {
		t.Fatalf("kind = %v, want AlterTableAlterColumn", a.Kind)
	}
	ca := a.ColumnAlters[0]
	if ca.Kind != ast.ColumnAlterDropDefault {
		t.Errorf("alter kind = %v, want DropDefault", ca.Kind)
	}
}

func TestAlterTable_AlterColumn_SetNotNull(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c SET NOT NULL")
	a := requireAction(t, stmt, 0)
	ca := a.ColumnAlters[0]
	if ca.Kind != ast.ColumnAlterSetNotNull {
		t.Errorf("alter kind = %v, want SetNotNull", ca.Kind)
	}
}

func TestAlterTable_AlterColumn_DropNotNull(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c DROP NOT NULL")
	a := requireAction(t, stmt, 0)
	ca := a.ColumnAlters[0]
	if ca.Kind != ast.ColumnAlterDropNotNull {
		t.Errorf("alter kind = %v, want DropNotNull", ca.Kind)
	}
}

func TestAlterTable_AlterColumn_Comment(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c COMMENT 'my column'")
	a := requireAction(t, stmt, 0)
	ca := a.ColumnAlters[0]
	if ca.Kind != ast.ColumnAlterSetComment {
		t.Errorf("alter kind = %v, want SetComment", ca.Kind)
	}
	if ca.Comment == nil || *ca.Comment != "my column" {
		t.Errorf("comment = %v", ca.Comment)
	}
}

func TestAlterTable_AlterColumn_UnsetComment(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c UNSET COMMENT")
	a := requireAction(t, stmt, 0)
	ca := a.ColumnAlters[0]
	if ca.Kind != ast.ColumnAlterUnsetComment {
		t.Errorf("alter kind = %v, want UnsetComment", ca.Kind)
	}
}

func TestAlterTable_AlterColumn_Parenthesized(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t MODIFY (col1 SET NOT NULL, col2 DROP DEFAULT)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAlterColumn {
		t.Fatalf("kind = %v, want AlterTableAlterColumn", a.Kind)
	}
	if len(a.ColumnAlters) != 2 {
		t.Fatalf("expected 2 ColumnAlters, got %d", len(a.ColumnAlters))
	}
	if a.ColumnAlters[0].Column.Normalize() != "COL1" {
		t.Errorf("col[0] = %q", a.ColumnAlters[0].Column.Normalize())
	}
	if a.ColumnAlters[1].Column.Normalize() != "COL2" {
		t.Errorf("col[1] = %q", a.ColumnAlters[1].Column.Normalize())
	}
}

// ---------------------------------------------------------------------------
// SET MASKING POLICY / UNSET MASKING POLICY
// ---------------------------------------------------------------------------

func TestAlterTable_SetMaskingPolicy(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN ssn SET MASKING POLICY pii_mask")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSetMaskingPolicy {
		t.Fatalf("kind = %v, want AlterTableSetMaskingPolicy", a.Kind)
	}
	if a.MaskColumn.Normalize() != "SSN" {
		t.Errorf("MaskColumn = %q", a.MaskColumn.Normalize())
	}
	if a.MaskingPolicy == nil || a.MaskingPolicy.Normalize() != "PII_MASK" {
		t.Errorf("MaskingPolicy = %v", a.MaskingPolicy)
	}
}

func TestAlterTable_UnsetMaskingPolicy(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN ssn UNSET MASKING POLICY")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableUnsetMaskingPolicy {
		t.Fatalf("kind = %v, want AlterTableUnsetMaskingPolicy", a.Kind)
	}
	if a.MaskColumn.Normalize() != "SSN" {
		t.Errorf("MaskColumn = %q", a.MaskColumn.Normalize())
	}
}

func TestAlterTable_ModifyColumn_SetMaskingPolicyWithForce(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t MODIFY COLUMN email SET MASKING POLICY email_mask USING (email, scope) FORCE")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSetMaskingPolicy {
		t.Fatalf("kind = %v, want AlterTableSetMaskingPolicy", a.Kind)
	}
}

// ---------------------------------------------------------------------------
// SET / UNSET table properties
// ---------------------------------------------------------------------------

func TestAlterTable_SetComment(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET COMMENT = 'my table'")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSet {
		t.Fatalf("kind = %v, want AlterTableSet", a.Kind)
	}
	if len(a.Props) != 1 {
		t.Fatalf("expected 1 prop, got %d", len(a.Props))
	}
	if a.Props[0].Name != "COMMENT" {
		t.Errorf("prop name = %q", a.Props[0].Name)
	}
	if a.Props[0].Value != "my table" {
		t.Errorf("prop value = %q", a.Props[0].Value)
	}
}

func TestAlterTable_SetDataRetention(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET DATA_RETENTION_TIME_IN_DAYS = 7")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSet {
		t.Fatalf("kind = %v, want AlterTableSet", a.Kind)
	}
	if len(a.Props) < 1 || a.Props[0].Name != "DATA_RETENTION_TIME_IN_DAYS" {
		t.Errorf("props = %v", a.Props)
	}
}

func TestAlterTable_SetChangeTracking(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET CHANGE_TRACKING = TRUE")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSet {
		t.Fatalf("kind = %v, want AlterTableSet", a.Kind)
	}
}

func TestAlterTable_UnsetProps(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t UNSET DATA_RETENTION_TIME_IN_DAYS, COMMENT")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableUnset {
		t.Fatalf("kind = %v, want AlterTableUnset", a.Kind)
	}
	if len(a.UnsetProps) != 2 {
		t.Fatalf("expected 2 unset props, got %d", len(a.UnsetProps))
	}
}

// ---------------------------------------------------------------------------
// SET TAG / UNSET TAG
// ---------------------------------------------------------------------------

func TestAlterTable_SetTag(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t SET TAG (my_tag = 'v1', other_tag = 'v2')")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSetTag {
		t.Fatalf("kind = %v, want AlterTableSetTag", a.Kind)
	}
	if len(a.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(a.Tags))
	}
	if a.Tags[0].Value != "v1" {
		t.Errorf("tag[0] value = %q", a.Tags[0].Value)
	}
}

func TestAlterTable_UnsetTag(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t UNSET TAG (my_tag, other_tag)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableUnsetTag {
		t.Fatalf("kind = %v, want AlterTableUnsetTag", a.Kind)
	}
	if len(a.UnsetTags) != 2 {
		t.Fatalf("expected 2 unset tags, got %d", len(a.UnsetTags))
	}
}

// ---------------------------------------------------------------------------
// ROW ACCESS POLICY
// ---------------------------------------------------------------------------

func TestAlterTable_AddRowAccessPolicy(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD ROW ACCESS POLICY rap1 ON (col1, col2)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddRowAccessPolicy {
		t.Fatalf("kind = %v, want AlterTableAddRowAccessPolicy", a.Kind)
	}
	if a.PolicyName == nil || a.PolicyName.Normalize() != "RAP1" {
		t.Errorf("PolicyName = %v", a.PolicyName)
	}
	if len(a.PolicyCols) != 2 {
		t.Errorf("PolicyCols = %v", a.PolicyCols)
	}
}

func TestAlterTable_DropRowAccessPolicy(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP ROW ACCESS POLICY rap1")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropRowAccessPolicy {
		t.Fatalf("kind = %v, want AlterTableDropRowAccessPolicy", a.Kind)
	}
	if a.PolicyName == nil || a.PolicyName.Normalize() != "RAP1" {
		t.Errorf("PolicyName = %v", a.PolicyName)
	}
}

func TestAlterTable_DropAllRowAccessPolicies(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP ALL ROW ACCESS POLICIES")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropAllRowAccessPolicies {
		t.Fatalf("kind = %v, want AlterTableDropAllRowAccessPolicies", a.Kind)
	}
}

// ---------------------------------------------------------------------------
// SEARCH OPTIMIZATION
// ---------------------------------------------------------------------------

func TestAlterTable_AddSearchOptimization(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddSearchOpt {
		t.Fatalf("kind = %v, want AlterTableAddSearchOpt", a.Kind)
	}
	if len(a.SearchOptOn) != 0 {
		t.Errorf("expected no ON targets, got %v", a.SearchOptOn)
	}
}

func TestAlterTable_DropSearchOptimization(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP SEARCH OPTIMIZATION")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropSearchOpt {
		t.Fatalf("kind = %v, want AlterTableDropSearchOpt", a.Kind)
	}
}

func TestAlterTable_AddSearchOptimizationOnEquality(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(col1)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddSearchOpt {
		t.Fatalf("kind = %v, want AlterTableAddSearchOpt", a.Kind)
	}
	if len(a.SearchOptOn) != 1 {
		t.Errorf("expected 1 ON target, got %v", a.SearchOptOn)
	}
}

// ---------------------------------------------------------------------------
// Column SET/UNSET TAG
// ---------------------------------------------------------------------------

func TestAlterTable_SetColumnTag(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c SET TAG (sensitivity = 'pii')")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableSetColumnTag {
		t.Fatalf("kind = %v, want AlterTableSetColumnTag", a.Kind)
	}
	if a.TagColumn.Normalize() != "C" {
		t.Errorf("TagColumn = %q", a.TagColumn.Normalize())
	}
	if len(a.Tags) != 1 || a.Tags[0].Value != "pii" {
		t.Errorf("tags = %v", a.Tags)
	}
}

func TestAlterTable_UnsetColumnTag(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ALTER COLUMN c UNSET TAG (sensitivity)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableUnsetColumnTag {
		t.Fatalf("kind = %v, want AlterTableUnsetColumnTag", a.Kind)
	}
	if a.TagColumn.Normalize() != "C" {
		t.Errorf("TagColumn = %q", a.TagColumn.Normalize())
	}
}

// ---------------------------------------------------------------------------
// Recluster
// ---------------------------------------------------------------------------

func TestAlterTable_Recluster(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t RECLUSTER MAX_SIZE = 1000 WHERE id > 0")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableRecluster {
		t.Fatalf("kind = %v, want AlterTableRecluster", a.Kind)
	}
	if a.ReclusterMaxSize == nil || *a.ReclusterMaxSize != 1000 {
		t.Errorf("ReclusterMaxSize = %v", a.ReclusterMaxSize)
	}
	if a.ReclusterWhere == nil {
		t.Error("expected ReclusterWhere != nil")
	}
}

// ---------------------------------------------------------------------------
// Multi-action comma-separated (top-level)
// ---------------------------------------------------------------------------

func TestAlterTable_MultiAction_DropUnique(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t DROP UNIQUE (col1)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropConstraint {
		t.Fatalf("kind = %v, want AlterTableDropConstraint", a.Kind)
	}
	if !a.DropUnique {
		t.Error("expected DropUnique = true")
	}
}

// ---------------------------------------------------------------------------
// Location tracking
// ---------------------------------------------------------------------------

func TestAlterTable_LocTracking(t *testing.T) {
	input := "ALTER TABLE t ADD COLUMN c INT"
	stmt, errs := testParseAlterTable(input)
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	if !stmt.Loc.IsValid() {
		t.Error("stmt.Loc is not valid")
	}
	// Loc.Start is at the TABLE token (ALTER is consumed by the dispatcher),
	// consistent with how AlterDatabaseStmt and AlterSchemaStmt set their loc.
	if stmt.Loc.Start <= 0 {
		t.Errorf("stmt.Loc.Start = %d, expected > 0 (TABLE keyword position)", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("stmt.Loc.End (%d) <= Start (%d)", stmt.Loc.End, stmt.Loc.Start)
	}
}
