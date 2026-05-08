package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseAlterTableAddColumn tests ALTER TABLE ADD column.
func TestParseAlterTableAddColumn(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees ADD (email VARCHAR2(100) NOT NULL)")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	if stmt.Name == nil || stmt.Name.Name != "EMPLOYEES" {
		t.Errorf("expected table name EMPLOYEES, got %v", stmt.Name)
	}
	if stmt.Actions == nil || stmt.Actions.Len() != 1 {
		t.Fatalf("expected 1 action, got %d", stmt.Actions.Len())
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_ADD_COLUMN {
		t.Errorf("expected AT_ADD_COLUMN, got %d", cmd.Action)
	}
	if cmd.ColumnDef == nil {
		t.Fatal("expected non-nil ColumnDef")
	}
	if cmd.ColumnDef.Name != "EMAIL" {
		t.Errorf("expected column name EMAIL, got %q", cmd.ColumnDef.Name)
	}
	if !cmd.ColumnDef.NotNull {
		t.Error("expected NOT NULL on column")
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
}

// TestParseAlterTableAddColumnNoParens tests ALTER TABLE ADD column without parentheses.
func TestParseAlterTableAddColumnNoParens(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE t1 ADD col1 NUMBER")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_ADD_COLUMN {
		t.Errorf("expected AT_ADD_COLUMN, got %d", cmd.Action)
	}
	if cmd.ColumnDef == nil || cmd.ColumnDef.Name != "COL1" {
		t.Errorf("expected column name COL1, got %v", cmd.ColumnDef)
	}
}

// TestParseAlterTableModifyColumn tests ALTER TABLE MODIFY column.
func TestParseAlterTableModifyColumn(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees MODIFY (salary NUMBER(12,2) DEFAULT 0)")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	if stmt.Actions.Len() != 1 {
		t.Fatalf("expected 1 action, got %d", stmt.Actions.Len())
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_MODIFY_COLUMN {
		t.Errorf("expected AT_MODIFY_COLUMN, got %d", cmd.Action)
	}
	if cmd.ColumnDef == nil || cmd.ColumnDef.Name != "SALARY" {
		t.Errorf("expected column name SALARY, got %v", cmd.ColumnDef)
	}
}

// TestParseAlterTableDropColumn tests ALTER TABLE DROP COLUMN.
func TestParseAlterTableDropColumn(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees DROP COLUMN middle_name")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_DROP_COLUMN {
		t.Errorf("expected AT_DROP_COLUMN, got %d", cmd.Action)
	}
	if cmd.ColumnName != "MIDDLE_NAME" {
		t.Errorf("expected column name MIDDLE_NAME, got %q", cmd.ColumnName)
	}
}

// TestParseAlterTableRenameColumn tests ALTER TABLE RENAME COLUMN x TO y.
func TestParseAlterTableRenameColumn(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees RENAME COLUMN old_name TO new_name")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_RENAME_COLUMN {
		t.Errorf("expected AT_RENAME_COLUMN, got %d", cmd.Action)
	}
	if cmd.ColumnName != "OLD_NAME" {
		t.Errorf("expected old name OLD_NAME, got %q", cmd.ColumnName)
	}
	if cmd.NewName != "NEW_NAME" {
		t.Errorf("expected new name NEW_NAME, got %q", cmd.NewName)
	}
}

// TestParseAlterTableAddConstraint tests ALTER TABLE ADD CONSTRAINT.
func TestParseAlterTableAddConstraint(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees ADD CONSTRAINT pk_emp PRIMARY KEY (employee_id)")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_ADD_CONSTRAINT {
		t.Errorf("expected AT_ADD_CONSTRAINT, got %d", cmd.Action)
	}
	if cmd.Constraint == nil {
		t.Fatal("expected non-nil Constraint")
	}
	if cmd.Constraint.Name != "PK_EMP" {
		t.Errorf("expected constraint name PK_EMP, got %q", cmd.Constraint.Name)
	}
	if cmd.Constraint.Type != ast.CONSTRAINT_PRIMARY {
		t.Errorf("expected CONSTRAINT_PRIMARY, got %d", cmd.Constraint.Type)
	}
}

// TestParseAlterTableDropConstraint tests ALTER TABLE DROP CONSTRAINT.
func TestParseAlterTableDropConstraint(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees DROP CONSTRAINT pk_emp")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_DROP_CONSTRAINT {
		t.Errorf("expected AT_DROP_CONSTRAINT, got %d", cmd.Action)
	}
	if cmd.Constraint == nil || cmd.Constraint.Name != "PK_EMP" {
		t.Errorf("expected constraint name PK_EMP, got %v", cmd.Constraint)
	}
}

// TestParseAlterTableRenameTo tests ALTER TABLE RENAME TO new_name.
func TestParseAlterTableRenameTo(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLE employees RENAME TO staff")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_RENAME {
		t.Errorf("expected AT_RENAME, got %d", cmd.Action)
	}
	if cmd.NewName != "STAFF" {
		t.Errorf("expected new name STAFF, got %q", cmd.NewName)
	}
}

// TestParseAlterTableMultipleActions tests ALTER TABLE with multiple actions.
func TestParseAlterTableMultipleActions(t *testing.T) {
	sql := "ALTER TABLE t1 ADD (col1 NUMBER) MODIFY (col2 VARCHAR2(50))"
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	if stmt.Actions.Len() != 2 {
		t.Fatalf("expected 2 actions, got %d", stmt.Actions.Len())
	}
	cmd1 := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd1.Action != ast.AT_ADD_COLUMN {
		t.Errorf("expected AT_ADD_COLUMN for action 1, got %d", cmd1.Action)
	}
	cmd2 := stmt.Actions.Items[1].(*ast.AlterTableCmd)
	if cmd2.Action != ast.AT_MODIFY_COLUMN {
		t.Errorf("expected AT_MODIFY_COLUMN for action 2, got %d", cmd2.Action)
	}
}

func TestP2AlterTablePartitionOptionsAndActionBoundaries(t *testing.T) {
	sql := "ALTER TABLE sales ADD PARTITION p1 VALUES LESS THAN (10) TABLESPACE ts1 MOVE PARTITION p2 TABLESPACE ts2 UPDATE INDEXES"
	stmt := parseAlterTableForP2(t, sql)
	if stmt.Actions.Len() != 2 {
		t.Fatalf("expected 2 actions, got %d", stmt.Actions.Len())
	}

	add := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if add.Action != ast.AT_ADD_PARTITION || add.Subtype != "PARTITION" || add.ColumnName != "P1" {
		t.Fatalf("unexpected ADD PARTITION command: action=%d subtype=%q name=%q", add.Action, add.Subtype, add.ColumnName)
	}
	assertAlterTableOption(t, add, "VALUES", "LESS THAN ( 10 )")
	assertAlterTableOption(t, add, "TABLESPACE", "TS1")
	if option := findAlterTableOption(add, "MOVE"); option != nil {
		t.Fatalf("ADD PARTITION options swallowed next action: %#v", option)
	}

	move := stmt.Actions.Items[1].(*ast.AlterTableCmd)
	if move.Action != ast.AT_MOVE || move.Subtype != "PARTITION" || move.ColumnName != "P2" {
		t.Fatalf("unexpected MOVE PARTITION command: action=%d subtype=%q name=%q", move.Action, move.Subtype, move.ColumnName)
	}
	assertAlterTableOption(t, move, "TABLESPACE", "TS2")
	assertAlterTableOption(t, move, "UPDATE", "INDEXES")
}

func TestP2AlterTableSupplementalLogOptions(t *testing.T) {
	stmt := parseAlterTableForP2(t, "ALTER TABLE sales ADD SUPPLEMENTAL LOG DATA (PRIMARY KEY) COLUMNS")
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_ALTER_PROPERTY || cmd.Subtype != "ADD SUPPLEMENTAL LOG" {
		t.Fatalf("unexpected supplemental log command: action=%d subtype=%q", cmd.Action, cmd.Subtype)
	}
	assertAlterTableOption(t, cmd, "DATA", "( PRIMARY KEY )")
	assertAlterTableOption(t, cmd, "COLUMNS", "")
}

func TestP2AlterTableOptionLoc(t *testing.T) {
	stmt := parseAlterTableForP2(t, "ALTER TABLE sales MOVE PARTITION p1 TABLESPACE ts1 UPDATE INDEXES")
	cmd := stmt.Actions.Items[0].(*ast.AlterTableCmd)
	opt := findAlterTableOption(cmd, "TABLESPACE")
	if opt == nil {
		t.Fatalf("expected TABLESPACE option in %#v", cmd.Options)
	}
	if opt.Loc.Start <= cmd.Loc.Start || opt.Loc.End > cmd.Loc.End || opt.Loc.Start >= opt.Loc.End {
		t.Fatalf("option Loc=%+v is not inside command Loc=%+v", opt.Loc, cmd.Loc)
	}
}

func TestP2AlterTableMovePartitionRequiresName(t *testing.T) {
	ParseShouldFail(t, "ALTER TABLE sales MOVE PARTITION TABLESPACE ts1")
}

func parseAlterTableForP2(t *testing.T, sql string) *ast.AlterTableStmt {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	if stmt.Actions == nil || stmt.Actions.Len() == 0 {
		t.Fatalf("expected ALTER TABLE actions, got %#v", stmt.Actions)
	}
	return stmt
}

func assertAlterTableOption(t *testing.T, cmd *ast.AlterTableCmd, key, value string) {
	t.Helper()
	opt := findAlterTableOption(cmd, key)
	if opt == nil {
		t.Fatalf("expected option %q in %#v", key, cmd.Options)
	}
	if opt.Value != value {
		t.Fatalf("option %q value=%q, want %q", key, opt.Value, value)
	}
	if opt.Loc.Start < cmd.Loc.Start || opt.Loc.End > cmd.Loc.End || opt.Loc.Start >= opt.Loc.End {
		t.Fatalf("option %q Loc=%+v is not inside command Loc=%+v", key, opt.Loc, cmd.Loc)
	}
}

func findAlterTableOption(cmd *ast.AlterTableCmd, key string) *ast.DDLOption {
	if cmd.Options == nil {
		return nil
	}
	for _, item := range cmd.Options.Items {
		opt, ok := item.(*ast.DDLOption)
		if ok && opt.Key == key {
			return opt
		}
	}
	return nil
}
