package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// BEGIN
// ---------------------------------------------------------------------------

func TestBegin_Basic(t *testing.T) {
	file, errs := Parse("BEGIN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.BeginStmt)
	if !ok {
		t.Fatalf("expected *ast.BeginStmt, got %T", file.Stmts[0])
	}
	if stmt.Label != "" {
		t.Errorf("Label = %q, want empty", stmt.Label)
	}
	if stmt.Tag() != ast.T_BeginStmt {
		t.Errorf("Tag() = %v, want T_BeginStmt", stmt.Tag())
	}
}

func TestBegin_WithLabel(t *testing.T) {
	file, errs := Parse("BEGIN WITH LABEL my_label")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.BeginStmt)
	if !ok {
		t.Fatalf("expected *ast.BeginStmt, got %T", file.Stmts[0])
	}
	if stmt.Label != "my_label" {
		t.Errorf("Label = %q, want %q", stmt.Label, "my_label")
	}
}

func TestBegin_WithLabelFromLegacyCorpus(t *testing.T) {
	// From doris/parser/testdata/legacy/transaction.sql
	file, errs := Parse("BEGIN WITH LABEL load_1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.BeginStmt)
	if stmt.Label != "load_1" {
		t.Errorf("Label = %q, want %q", stmt.Label, "load_1")
	}
}

// ---------------------------------------------------------------------------
// COMMIT
// ---------------------------------------------------------------------------

func TestCommit_Basic(t *testing.T) {
	file, errs := Parse("COMMIT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CommitStmt)
	if !ok {
		t.Fatalf("expected *ast.CommitStmt, got %T", file.Stmts[0])
	}
	if stmt.Work {
		t.Error("Work should be false")
	}
	if stmt.Chain != "" {
		t.Errorf("Chain = %q, want empty", stmt.Chain)
	}
	if stmt.Release != "" {
		t.Errorf("Release = %q, want empty", stmt.Release)
	}
	if stmt.Tag() != ast.T_CommitStmt {
		t.Errorf("Tag() = %v, want T_CommitStmt", stmt.Tag())
	}
}

func TestCommit_WorkAndChain(t *testing.T) {
	file, errs := Parse("COMMIT WORK AND CHAIN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CommitStmt)
	if !stmt.Work {
		t.Error("Work should be true")
	}
	if stmt.Chain != "AND CHAIN" {
		t.Errorf("Chain = %q, want %q", stmt.Chain, "AND CHAIN")
	}
	if stmt.Release != "" {
		t.Errorf("Release = %q, want empty", stmt.Release)
	}
}

func TestCommit_WorkAndNoChainNoRelease(t *testing.T) {
	file, errs := Parse("COMMIT WORK AND NO CHAIN NO RELEASE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CommitStmt)
	if !stmt.Work {
		t.Error("Work should be true")
	}
	if stmt.Chain != "AND NO CHAIN" {
		t.Errorf("Chain = %q, want %q", stmt.Chain, "AND NO CHAIN")
	}
	if stmt.Release != "NO RELEASE" {
		t.Errorf("Release = %q, want %q", stmt.Release, "NO RELEASE")
	}
}

func TestCommit_Release(t *testing.T) {
	file, errs := Parse("COMMIT RELEASE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CommitStmt)
	if stmt.Release != "RELEASE" {
		t.Errorf("Release = %q, want %q", stmt.Release, "RELEASE")
	}
}

func TestCommit_NoRelease(t *testing.T) {
	file, errs := Parse("COMMIT NO RELEASE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CommitStmt)
	if stmt.Release != "NO RELEASE" {
		t.Errorf("Release = %q, want %q", stmt.Release, "NO RELEASE")
	}
}

// ---------------------------------------------------------------------------
// ROLLBACK
// ---------------------------------------------------------------------------

func TestRollback_Basic(t *testing.T) {
	file, errs := Parse("ROLLBACK")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.RollbackStmt)
	if !ok {
		t.Fatalf("expected *ast.RollbackStmt, got %T", file.Stmts[0])
	}
	if stmt.Work {
		t.Error("Work should be false")
	}
	if stmt.Chain != "" {
		t.Errorf("Chain = %q, want empty", stmt.Chain)
	}
	if stmt.Release != "" {
		t.Errorf("Release = %q, want empty", stmt.Release)
	}
	if stmt.Tag() != ast.T_RollbackStmt {
		t.Errorf("Tag() = %v, want T_RollbackStmt", stmt.Tag())
	}
}

func TestRollback_WorkAndChainRelease(t *testing.T) {
	file, errs := Parse("ROLLBACK WORK AND CHAIN RELEASE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.RollbackStmt)
	if !stmt.Work {
		t.Error("Work should be true")
	}
	if stmt.Chain != "AND CHAIN" {
		t.Errorf("Chain = %q, want %q", stmt.Chain, "AND CHAIN")
	}
	if stmt.Release != "RELEASE" {
		t.Errorf("Release = %q, want %q", stmt.Release, "RELEASE")
	}
}

func TestRollback_WorkAndNoChain(t *testing.T) {
	file, errs := Parse("ROLLBACK WORK AND NO CHAIN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.RollbackStmt)
	if stmt.Chain != "AND NO CHAIN" {
		t.Errorf("Chain = %q, want %q", stmt.Chain, "AND NO CHAIN")
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestBeginStmt_Tag(t *testing.T) {
	file, errs := Parse("BEGIN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_BeginStmt {
		t.Errorf("Tag() = %v, want T_BeginStmt", file.Stmts[0].Tag())
	}
}

func TestCommitStmt_Tag(t *testing.T) {
	file, errs := Parse("COMMIT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CommitStmt {
		t.Errorf("Tag() = %v, want T_CommitStmt", file.Stmts[0].Tag())
	}
}

func TestRollbackStmt_Tag(t *testing.T) {
	file, errs := Parse("ROLLBACK")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_RollbackStmt {
		t.Errorf("Tag() = %v, want T_RollbackStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Loc sanity check
// ---------------------------------------------------------------------------

func TestBegin_Loc(t *testing.T) {
	input := "BEGIN WITH LABEL lbl"
	file, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.BeginStmt)
	loc := ast.NodeLoc(stmt)
	if !loc.IsValid() {
		t.Error("Loc should be valid")
	}
	if loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", loc.Start)
	}
}

// ---------------------------------------------------------------------------
// Multiple TCL statements
// ---------------------------------------------------------------------------

func TestTCL_MultiStatement(t *testing.T) {
	input := "BEGIN; COMMIT; ROLLBACK"
	file, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 3 {
		t.Fatalf("expected 3 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.BeginStmt); !ok {
		t.Errorf("Stmts[0]: expected *ast.BeginStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.CommitStmt); !ok {
		t.Errorf("Stmts[1]: expected *ast.CommitStmt, got %T", file.Stmts[1])
	}
	if _, ok := file.Stmts[2].(*ast.RollbackStmt); !ok {
		t.Errorf("Stmts[2]: expected *ast.RollbackStmt, got %T", file.Stmts[2])
	}
}
