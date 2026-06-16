package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// CTE-prefixed DELETE (PR3): StarRocks allows a WITH clause before DELETE
// (but NOT before UPDATE).

func TestCTEDelete(t *testing.T) {
	file, errs := Parse("WITH c AS (SELECT id FROM s) DELETE FROM t WHERE id IN (SELECT id FROM c)")
	if len(errs) > 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	del, ok := file.Stmts[0].(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.DeleteStmt", file.Stmts[0])
	}
	if del.With == nil || len(del.With.CTEs) != 1 || del.With.CTEs[0].Name != "c" {
		t.Errorf("With = %+v, want one CTE named c", del.With)
	}
	if del.Target == nil || del.Target.Parts[len(del.Target.Parts)-1] != "t" {
		t.Errorf("target = %+v, want t", del.Target)
	}
}

func TestCTEDeleteMultiple(t *testing.T) {
	file, errs := Parse("WITH a AS (SELECT 1), b AS (SELECT 2) DELETE FROM t WHERE k IN (SELECT x FROM a)")
	if len(errs) > 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	del := file.Stmts[0].(*ast.DeleteStmt)
	if del.With == nil || len(del.With.CTEs) != 2 {
		t.Errorf("want 2 CTEs, got %+v", del.With)
	}
}

// Regression: CTE-prefixed SELECT must still parse as a SelectStmt with its WITH.
func TestCTESelectStillWorks(t *testing.T) {
	file, errs := Parse("WITH c AS (SELECT a FROM real_t) SELECT a FROM c")
	if len(errs) > 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	sel, ok := file.Stmts[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.SelectStmt", file.Stmts[0])
	}
	if sel.With == nil {
		t.Error("With = nil, want the CTE clause")
	}
}

// CTE-prefixed UPDATE is also accepted by StarRocks (the grammar's
// updateStatement takes withClause?). The earlier "rejected" assumption came
// from a probe that used UPDATE … FROM c, which StarRocks rejects for an
// unrelated reason (it has no UPDATE … FROM form).
func TestCTEUpdate(t *testing.T) {
	file, errs := Parse("WITH c AS (SELECT id FROM s) UPDATE t SET v = 1 WHERE id IN (SELECT id FROM c)")
	if len(errs) > 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	upd, ok := file.Stmts[0].(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.UpdateStmt", file.Stmts[0])
	}
	if upd.With == nil || len(upd.With.CTEs) != 1 || upd.With.CTEs[0].Name != "c" {
		t.Errorf("With = %+v, want one CTE named c", upd.With)
	}
}

// WITH … INSERT is NOT allowed by StarRocks (insertStatement has no withClause)
// and must stay rejected.
func TestCTEInsertRejected(t *testing.T) {
	_, errs := Parse("WITH c AS (SELECT 1) INSERT INTO t SELECT x FROM c")
	if len(errs) == 0 {
		t.Fatal("expected a parse error for CTE-prefixed INSERT, got none")
	}
}
