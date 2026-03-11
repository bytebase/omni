package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCommit(t *testing.T) {
	result := ParseAndCheck(t, "COMMIT")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CommitStmt)
	if !ok {
		t.Fatalf("expected CommitStmt, got %T", raw.Stmt)
	}
	if stmt.Work {
		t.Error("expected Work=false")
	}
}

func TestParseCommitWork(t *testing.T) {
	result := ParseAndCheck(t, "COMMIT WORK")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommitStmt)
	if !stmt.Work {
		t.Error("expected Work=true")
	}
}

func TestParseCommitComment(t *testing.T) {
	result := ParseAndCheck(t, "COMMIT COMMENT 'batch job 123'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommitStmt)
	if stmt.Comment != "batch job 123" {
		t.Errorf("expected comment 'batch job 123', got %q", stmt.Comment)
	}
}

func TestParseCommitForce(t *testing.T) {
	result := ParseAndCheck(t, "COMMIT FORCE 'txn123'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CommitStmt)
	if stmt.Force != "txn123" {
		t.Errorf("expected force 'txn123', got %q", stmt.Force)
	}
}

func TestParseRollback(t *testing.T) {
	result := ParseAndCheck(t, "ROLLBACK")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.RollbackStmt)
	if !ok {
		t.Fatalf("expected RollbackStmt, got %T", raw.Stmt)
	}
	if stmt.Work {
		t.Error("expected Work=false")
	}
	if stmt.ToSavepoint != "" {
		t.Errorf("expected empty ToSavepoint, got %q", stmt.ToSavepoint)
	}
}

func TestParseRollbackWork(t *testing.T) {
	result := ParseAndCheck(t, "ROLLBACK WORK")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.RollbackStmt)
	if !stmt.Work {
		t.Error("expected Work=true")
	}
}

func TestParseRollbackToSavepoint(t *testing.T) {
	result := ParseAndCheck(t, "ROLLBACK TO SAVEPOINT sp1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.RollbackStmt)
	if stmt.ToSavepoint != "SP1" {
		t.Errorf("expected SP1, got %q", stmt.ToSavepoint)
	}
}

func TestParseRollbackTo(t *testing.T) {
	result := ParseAndCheck(t, "ROLLBACK TO sp1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.RollbackStmt)
	if stmt.ToSavepoint != "SP1" {
		t.Errorf("expected SP1, got %q", stmt.ToSavepoint)
	}
}

func TestParseRollbackForce(t *testing.T) {
	result := ParseAndCheck(t, "ROLLBACK FORCE 'txn456'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.RollbackStmt)
	if stmt.Force != "txn456" {
		t.Errorf("expected force 'txn456', got %q", stmt.Force)
	}
}

func TestParseSavepoint(t *testing.T) {
	result := ParseAndCheck(t, "SAVEPOINT sp1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.SavepointStmt)
	if !ok {
		t.Fatalf("expected SavepointStmt, got %T", raw.Stmt)
	}
	if stmt.Name != "SP1" {
		t.Errorf("expected SP1, got %q", stmt.Name)
	}
}

func TestParseSetTransactionReadOnly(t *testing.T) {
	result := ParseAndCheck(t, "SET TRANSACTION READ ONLY")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.SetTransactionStmt)
	if !ok {
		t.Fatalf("expected SetTransactionStmt, got %T", raw.Stmt)
	}
	if !stmt.ReadOnly {
		t.Error("expected ReadOnly=true")
	}
}

func TestParseSetTransactionReadWrite(t *testing.T) {
	result := ParseAndCheck(t, "SET TRANSACTION READ WRITE")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SetTransactionStmt)
	if !stmt.ReadWrite {
		t.Error("expected ReadWrite=true")
	}
}

func TestParseSetTransactionIsolationSerializable(t *testing.T) {
	result := ParseAndCheck(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SetTransactionStmt)
	if stmt.IsolLevel != "SERIALIZABLE" {
		t.Errorf("expected SERIALIZABLE, got %q", stmt.IsolLevel)
	}
}

func TestParseSetTransactionIsolationReadCommitted(t *testing.T) {
	result := ParseAndCheck(t, "SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SetTransactionStmt)
	if stmt.IsolLevel != "READ COMMITTED" {
		t.Errorf("expected READ COMMITTED, got %q", stmt.IsolLevel)
	}
}

func TestParseSetTransactionName(t *testing.T) {
	result := ParseAndCheck(t, "SET TRANSACTION NAME 'my_txn'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SetTransactionStmt)
	if stmt.Name != "my_txn" {
		t.Errorf("expected 'my_txn', got %q", stmt.Name)
	}
}

func TestParseSetTransactionReadOnlyName(t *testing.T) {
	result := ParseAndCheck(t, "SET TRANSACTION READ ONLY NAME 'readonly_txn'")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SetTransactionStmt)
	if !stmt.ReadOnly {
		t.Error("expected ReadOnly=true")
	}
	if stmt.Name != "readonly_txn" {
		t.Errorf("expected 'readonly_txn', got %q", stmt.Name)
	}
}
