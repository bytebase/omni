package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseInsertSimpleValues tests a basic INSERT INTO ... VALUES statement.
func TestParseInsertSimpleValues(t *testing.T) {
	result := ParseAndCheck(t, "INSERT INTO t (c1, c2) VALUES (1, 2)")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	ins, ok := raw.Stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected InsertStmt, got %T", raw.Stmt)
	}
	if ins.InsertType != ast.INSERT_SINGLE {
		t.Errorf("expected INSERT_SINGLE, got %d", ins.InsertType)
	}
	if ins.Table == nil || ins.Table.Name != "T" {
		t.Errorf("expected table T, got %v", ins.Table)
	}
	if ins.Columns == nil || ins.Columns.Len() != 2 {
		t.Fatalf("expected 2 columns, got %v", ins.Columns)
	}
	col0 := ins.Columns.Items[0].(*ast.ColumnRef)
	if col0.Column != "C1" {
		t.Errorf("expected column C1, got %q", col0.Column)
	}
	col1 := ins.Columns.Items[1].(*ast.ColumnRef)
	if col1.Column != "C2" {
		t.Errorf("expected column C2, got %q", col1.Column)
	}
	if ins.Values == nil || ins.Values.Len() != 2 {
		t.Fatalf("expected 2 values, got %v", ins.Values)
	}
}

// TestParseInsertSelect tests INSERT INTO ... SELECT.
func TestParseInsertSelect(t *testing.T) {
	result := ParseAndCheck(t, "INSERT INTO t SELECT a, b FROM s")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	ins, ok := raw.Stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected InsertStmt, got %T", raw.Stmt)
	}
	if ins.InsertType != ast.INSERT_SINGLE {
		t.Errorf("expected INSERT_SINGLE, got %d", ins.InsertType)
	}
	if ins.Table == nil || ins.Table.Name != "T" {
		t.Errorf("expected table T, got %v", ins.Table)
	}
	if ins.Select == nil {
		t.Fatal("expected non-nil Select")
	}
	if ins.Select.TargetList.Len() != 2 {
		t.Errorf("expected 2 select targets, got %d", ins.Select.TargetList.Len())
	}
}

// TestParseInsertAll tests INSERT ALL with multiple INTO clauses.
func TestParseInsertAll(t *testing.T) {
	sql := `INSERT ALL
		INTO t1 (a) VALUES (1)
		INTO t2 (b) VALUES (2)
		SELECT * FROM dual`
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	ins, ok := raw.Stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected InsertStmt, got %T", raw.Stmt)
	}
	if ins.InsertType != ast.INSERT_ALL {
		t.Errorf("expected INSERT_ALL, got %d", ins.InsertType)
	}
	if ins.MultiTable == nil || ins.MultiTable.Len() != 2 {
		t.Fatalf("expected 2 INTO clauses, got %v", ins.MultiTable)
	}
	clause0 := ins.MultiTable.Items[0].(*ast.InsertIntoClause)
	if clause0.Table == nil || clause0.Table.Name != "T1" {
		t.Errorf("expected table T1, got %v", clause0.Table)
	}
	clause1 := ins.MultiTable.Items[1].(*ast.InsertIntoClause)
	if clause1.Table == nil || clause1.Table.Name != "T2" {
		t.Errorf("expected table T2, got %v", clause1.Table)
	}
	if ins.Subquery == nil {
		t.Error("expected non-nil Subquery")
	}
}

// TestParseInsertFirst tests INSERT FIRST with conditional WHEN/THEN/ELSE.
func TestParseInsertFirst(t *testing.T) {
	sql := `INSERT FIRST
		WHEN x > 10 THEN INTO t1 (a) VALUES (1)
		ELSE INTO t2 (b) VALUES (2)
		SELECT * FROM dual`
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	ins, ok := raw.Stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected InsertStmt, got %T", raw.Stmt)
	}
	if ins.InsertType != ast.INSERT_FIRST {
		t.Errorf("expected INSERT_FIRST, got %d", ins.InsertType)
	}
	if ins.MultiTable == nil || ins.MultiTable.Len() != 2 {
		t.Fatalf("expected 2 INTO clauses, got %d", ins.MultiTable.Len())
	}
	// First clause has a WHEN condition
	clause0 := ins.MultiTable.Items[0].(*ast.InsertIntoClause)
	if clause0.When == nil {
		t.Error("expected non-nil When for first clause")
	}
	if clause0.Table == nil || clause0.Table.Name != "T1" {
		t.Errorf("expected table T1, got %v", clause0.Table)
	}
	// Second clause (ELSE) has no When condition
	clause1 := ins.MultiTable.Items[1].(*ast.InsertIntoClause)
	if clause1.When != nil {
		t.Error("expected nil When for ELSE clause")
	}
	if clause1.Table == nil || clause1.Table.Name != "T2" {
		t.Errorf("expected table T2, got %v", clause1.Table)
	}
	if ins.Subquery == nil {
		t.Error("expected non-nil Subquery")
	}
}

// TestParseInsertReturning tests INSERT with RETURNING clause.
func TestParseInsertReturning(t *testing.T) {
	sql := "INSERT INTO t (c1) VALUES (1) RETURNING c1 INTO :b1"
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	ins, ok := raw.Stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected InsertStmt, got %T", raw.Stmt)
	}
	if ins.Returning == nil {
		t.Fatal("expected non-nil Returning")
	}
	if ins.Returning.Len() < 1 {
		t.Error("expected at least 1 returning expression")
	}
}

// TestParseInsertNoColumns tests INSERT without explicit column list.
func TestParseInsertNoColumns(t *testing.T) {
	result := ParseAndCheck(t, "INSERT INTO t VALUES (1, 2, 3)")
	raw := result.Items[0].(*ast.RawStmt)
	ins := raw.Stmt.(*ast.InsertStmt)
	if ins.Columns != nil {
		t.Errorf("expected nil Columns, got %v", ins.Columns)
	}
	if ins.Values == nil || ins.Values.Len() != 3 {
		t.Fatalf("expected 3 values, got %v", ins.Values)
	}
}

// TestParseInsertSchemaQualified tests INSERT with schema-qualified table name.
func TestParseInsertSchemaQualified(t *testing.T) {
	result := ParseAndCheck(t, "INSERT INTO hr.employees (id) VALUES (100)")
	raw := result.Items[0].(*ast.RawStmt)
	ins := raw.Stmt.(*ast.InsertStmt)
	if ins.Table.Schema != "HR" {
		t.Errorf("expected schema HR, got %q", ins.Table.Schema)
	}
	if ins.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table EMPLOYEES, got %q", ins.Table.Name)
	}
}
