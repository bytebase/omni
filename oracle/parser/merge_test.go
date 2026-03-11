package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseMergeMatchedUpdate tests MERGE with WHEN MATCHED THEN UPDATE SET.
func TestParseMergeMatchedUpdate(t *testing.T) {
	sql := `MERGE INTO t USING src ON (t.id = src.id) WHEN MATCHED THEN UPDATE SET t.name = src.name, t.val = src.val`
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	m, ok := raw.Stmt.(*ast.MergeStmt)
	if !ok {
		t.Fatalf("expected MergeStmt, got %T", raw.Stmt)
	}
	if m.Target.Name != "T" {
		t.Errorf("expected target T, got %q", m.Target.Name)
	}
	if m.Clauses.Len() != 1 {
		t.Fatalf("expected 1 clause, got %d", m.Clauses.Len())
	}
	mc := m.Clauses.Items[0].(*ast.MergeClause)
	if !mc.Matched {
		t.Errorf("expected Matched=true")
	}
	if mc.UpdateSet == nil || mc.UpdateSet.Len() != 2 {
		t.Errorf("expected 2 SET clauses, got %v", mc.UpdateSet)
	}
}

// TestParseMergeNotMatchedInsert tests MERGE with WHEN NOT MATCHED THEN INSERT.
func TestParseMergeNotMatchedInsert(t *testing.T) {
	sql := `MERGE INTO t USING src ON (t.id = src.id) WHEN NOT MATCHED THEN INSERT (id, name) VALUES (src.id, src.name)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	m := raw.Stmt.(*ast.MergeStmt)
	if m.Clauses.Len() != 1 {
		t.Fatalf("expected 1 clause, got %d", m.Clauses.Len())
	}
	mc := m.Clauses.Items[0].(*ast.MergeClause)
	if mc.Matched {
		t.Errorf("expected Matched=false")
	}
	if mc.InsertCols == nil || mc.InsertCols.Len() != 2 {
		t.Errorf("expected 2 insert cols, got %v", mc.InsertCols)
	}
	if mc.InsertVals == nil || mc.InsertVals.Len() != 2 {
		t.Errorf("expected 2 insert vals, got %v", mc.InsertVals)
	}
}

// TestParseMergeBothClauses tests MERGE with both WHEN MATCHED and WHEN NOT MATCHED.
func TestParseMergeBothClauses(t *testing.T) {
	sql := `MERGE INTO target t
		USING source s ON (t.id = s.id)
		WHEN MATCHED THEN UPDATE SET t.col1 = s.col1
		WHEN NOT MATCHED THEN INSERT (id, col1) VALUES (s.id, s.col1)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	m := raw.Stmt.(*ast.MergeStmt)
	if m.Target.Name != "TARGET" {
		t.Errorf("expected target TARGET, got %q", m.Target.Name)
	}
	if m.TargetAlias == nil || m.TargetAlias.Name != "T" {
		t.Errorf("expected target alias T, got %v", m.TargetAlias)
	}
	if m.SourceAlias == nil || m.SourceAlias.Name != "S" {
		t.Errorf("expected source alias S, got %v", m.SourceAlias)
	}
	if m.Clauses.Len() != 2 {
		t.Fatalf("expected 2 clauses, got %d", m.Clauses.Len())
	}
	mc0 := m.Clauses.Items[0].(*ast.MergeClause)
	if !mc0.Matched {
		t.Errorf("clause 0: expected Matched=true")
	}
	if mc0.UpdateSet == nil || mc0.UpdateSet.Len() != 1 {
		t.Errorf("clause 0: expected 1 SET clause")
	}
	mc1 := m.Clauses.Items[1].(*ast.MergeClause)
	if mc1.Matched {
		t.Errorf("clause 1: expected Matched=false")
	}
	if mc1.InsertCols == nil || mc1.InsertCols.Len() != 2 {
		t.Errorf("clause 1: expected 2 insert cols")
	}
}

// TestParseMergeMatchedDelete tests MERGE with WHEN MATCHED THEN DELETE.
func TestParseMergeMatchedDelete(t *testing.T) {
	sql := `MERGE INTO t USING src ON (t.id = src.id) WHEN MATCHED THEN DELETE`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	m := raw.Stmt.(*ast.MergeStmt)
	if m.Clauses.Len() != 1 {
		t.Fatalf("expected 1 clause, got %d", m.Clauses.Len())
	}
	mc := m.Clauses.Items[0].(*ast.MergeClause)
	if !mc.Matched {
		t.Errorf("expected Matched=true")
	}
	if !mc.IsDelete {
		t.Errorf("expected IsDelete=true")
	}
}

// TestParseMergeSourceAlias tests MERGE with source table alias.
func TestParseMergeSourceAlias(t *testing.T) {
	sql := `MERGE INTO employees e USING new_employees n ON (e.emp_id = n.emp_id) WHEN MATCHED THEN UPDATE SET e.name = n.name`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	m := raw.Stmt.(*ast.MergeStmt)
	if m.Target.Name != "EMPLOYEES" {
		t.Errorf("expected target EMPLOYEES, got %q", m.Target.Name)
	}
	if m.TargetAlias == nil || m.TargetAlias.Name != "E" {
		t.Errorf("expected target alias E, got %v", m.TargetAlias)
	}
	src, ok := m.Source.(*ast.TableRef)
	if !ok {
		t.Fatalf("expected TableRef source, got %T", m.Source)
	}
	if src.Name.Name != "NEW_EMPLOYEES" {
		t.Errorf("expected source NEW_EMPLOYEES, got %q", src.Name.Name)
	}
	if m.SourceAlias == nil || m.SourceAlias.Name != "N" {
		t.Errorf("expected source alias N, got %v", m.SourceAlias)
	}
	if m.On == nil {
		t.Errorf("expected ON condition")
	}
}
