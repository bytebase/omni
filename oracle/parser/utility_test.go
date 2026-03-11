package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseTruncateTable(t *testing.T) {
	result := ParseAndCheck(t, "TRUNCATE TABLE employees")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.TruncateStmt)
	if !ok {
		t.Fatalf("expected TruncateStmt, got %T", raw.Stmt)
	}
	if stmt.Table == nil || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %v", stmt.Table)
	}
	if stmt.Cluster {
		t.Error("expected Cluster=false")
	}
}

func TestParseTruncateTableSchema(t *testing.T) {
	result := ParseAndCheck(t, "TRUNCATE TABLE hr.employees")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.TruncateStmt)
	if stmt.Table.Schema != "HR" || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected HR.EMPLOYEES, got %q.%q", stmt.Table.Schema, stmt.Table.Name)
	}
}

func TestParseTruncateTableCascade(t *testing.T) {
	result := ParseAndCheck(t, "TRUNCATE TABLE employees CASCADE")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.TruncateStmt)
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
}

func TestParseAnalyzeTable(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE TABLE employees COMPUTE STATISTICS")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AnalyzeStmt)
	if !ok {
		t.Fatalf("expected AnalyzeStmt, got %T", raw.Stmt)
	}
	if stmt.ObjectType != ast.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.ObjectType)
	}
	if stmt.Table == nil || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %v", stmt.Table)
	}
	if stmt.Action != "COMPUTE STATISTICS" {
		t.Errorf("expected 'COMPUTE STATISTICS', got %q", stmt.Action)
	}
}

func TestParseAnalyzeIndex(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE INDEX idx_emp VALIDATE STRUCTURE")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.AnalyzeStmt)
	if stmt.ObjectType != ast.OBJECT_INDEX {
		t.Errorf("expected OBJECT_INDEX, got %d", stmt.ObjectType)
	}
	if stmt.Action != "VALIDATE STRUCTURE" {
		t.Errorf("expected 'VALIDATE STRUCTURE', got %q", stmt.Action)
	}
}

func TestParseAnalyzeDeleteStatistics(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE TABLE employees DELETE STATISTICS")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.AnalyzeStmt)
	if stmt.Action != "DELETE STATISTICS" {
		t.Errorf("expected 'DELETE STATISTICS', got %q", stmt.Action)
	}
}

func TestParseExplainPlan(t *testing.T) {
	result := ParseAndCheck(t, "EXPLAIN PLAN FOR SELECT * FROM employees")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.ExplainPlanStmt)
	if !ok {
		t.Fatalf("expected ExplainPlanStmt, got %T", raw.Stmt)
	}
	if stmt.Statement == nil {
		t.Fatal("expected non-nil Statement")
	}
	if _, ok := stmt.Statement.(*ast.SelectStmt); !ok {
		t.Fatalf("expected SelectStmt, got %T", stmt.Statement)
	}
}

func TestParseExplainPlanSetStatementId(t *testing.T) {
	result := ParseAndCheck(t, "EXPLAIN PLAN SET STATEMENT_ID = 'test1' FOR SELECT 1 FROM dual")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.ExplainPlanStmt)
	if stmt.StatementID != "test1" {
		t.Errorf("expected StatementID 'test1', got %q", stmt.StatementID)
	}
}

func TestParseExplainPlanInto(t *testing.T) {
	result := ParseAndCheck(t, "EXPLAIN PLAN INTO plan_table FOR SELECT 1 FROM dual")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.ExplainPlanStmt)
	if stmt.Into == nil || stmt.Into.Name != "PLAN_TABLE" {
		t.Errorf("expected Into=PLAN_TABLE, got %v", stmt.Into)
	}
}

func TestParseFlashbackTableToSCN(t *testing.T) {
	result := ParseAndCheck(t, "FLASHBACK TABLE employees TO SCN 12345")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.FlashbackTableStmt)
	if !ok {
		t.Fatalf("expected FlashbackTableStmt, got %T", raw.Stmt)
	}
	if stmt.Table == nil || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %v", stmt.Table)
	}
	if stmt.ToSCN == nil {
		t.Error("expected non-nil ToSCN")
	}
}

func TestParseFlashbackTableToTimestamp(t *testing.T) {
	result := ParseAndCheck(t, "FLASHBACK TABLE employees TO TIMESTAMP SYSDATE - 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.FlashbackTableStmt)
	if stmt.ToTimestamp == nil {
		t.Error("expected non-nil ToTimestamp")
	}
}

func TestParseFlashbackTableToBeforeDrop(t *testing.T) {
	result := ParseAndCheck(t, "FLASHBACK TABLE employees TO BEFORE DROP")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.FlashbackTableStmt)
	if !stmt.ToBeforeDrop {
		t.Error("expected ToBeforeDrop=true")
	}
}

func TestParseFlashbackTableToBeforeDropRename(t *testing.T) {
	result := ParseAndCheck(t, "FLASHBACK TABLE employees TO BEFORE DROP RENAME TO emp_restored")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.FlashbackTableStmt)
	if !stmt.ToBeforeDrop {
		t.Error("expected ToBeforeDrop=true")
	}
	if stmt.Rename != "EMP_RESTORED" {
		t.Errorf("expected EMP_RESTORED, got %q", stmt.Rename)
	}
}

func TestParsePurgeTable(t *testing.T) {
	result := ParseAndCheck(t, "PURGE TABLE employees")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.PurgeStmt)
	if !ok {
		t.Fatalf("expected PurgeStmt, got %T", raw.Stmt)
	}
	if stmt.ObjectType != ast.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.ObjectType)
	}
	if stmt.Name == nil || stmt.Name.Name != "EMPLOYEES" {
		t.Errorf("expected EMPLOYEES, got %v", stmt.Name)
	}
}

func TestParsePurgeIndex(t *testing.T) {
	result := ParseAndCheck(t, "PURGE INDEX idx_emp")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.PurgeStmt)
	if stmt.ObjectType != ast.OBJECT_INDEX {
		t.Errorf("expected OBJECT_INDEX, got %d", stmt.ObjectType)
	}
}

func TestParsePurgeRecyclebin(t *testing.T) {
	result := ParseAndCheck(t, "PURGE RECYCLEBIN")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.PurgeStmt)
	if stmt.Name == nil || stmt.Name.Name != "RECYCLEBIN" {
		t.Errorf("expected RECYCLEBIN, got %v", stmt.Name)
	}
}

func TestParsePurgeTablespace(t *testing.T) {
	result := ParseAndCheck(t, "PURGE TABLESPACE ts1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.PurgeStmt)
	if stmt.ObjectType != ast.OBJECT_TABLESPACE {
		t.Errorf("expected OBJECT_TABLESPACE, got %d", stmt.ObjectType)
	}
	if stmt.Name == nil || stmt.Name.Name != "TS1" {
		t.Errorf("expected TS1, got %v", stmt.Name)
	}
}
