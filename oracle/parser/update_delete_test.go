package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseUpdateSimple tests UPDATE t SET col = expr WHERE cond.
func TestParseUpdateSimple(t *testing.T) {
	result := ParseAndCheck(t, "UPDATE employees SET salary = 5000 WHERE id = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("expected UpdateStmt, got %T", raw.Stmt)
	}
	if stmt.Table == nil {
		t.Fatal("expected non-nil Table")
	}
	if stmt.SetClauses == nil || stmt.SetClauses.Len() != 1 {
		t.Fatalf("expected 1 SET clause, got %d", stmt.SetClauses.Len())
	}
	sc := stmt.SetClauses.Items[0].(*ast.SetClause)
	if sc.Column == nil || sc.Column.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %v", sc.Column)
	}
	if sc.Value == nil {
		t.Error("expected non-nil Value")
	}
	if stmt.WhereClause == nil {
		t.Error("expected non-nil WhereClause")
	}
}

// TestParseUpdateWithAlias tests UPDATE with table alias.
func TestParseUpdateWithAlias(t *testing.T) {
	result := ParseAndCheck(t, "UPDATE employees e SET e.salary = 5000 WHERE e.id = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.UpdateStmt)
	if stmt.Alias == nil || stmt.Alias.Name != "E" {
		t.Errorf("expected alias E, got %v", stmt.Alias)
	}
	target, ok := stmt.Target.(*ast.TableRef)
	if !ok {
		t.Fatalf("expected table target, got %T", stmt.Target)
	}
	if target.Name == nil || target.Name.Name != "EMPLOYEES" {
		t.Fatalf("expected target table EMPLOYEES, got %#v", target.Name)
	}
	if target.Alias == nil || target.Alias.Name != "E" {
		t.Fatalf("expected target alias E, got %#v", target.Alias)
	}
}

func TestParseUpdateInlineViewTarget(t *testing.T) {
	result := ParseAndCheck(t, "UPDATE (SELECT a FROM t) v SET a = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.UpdateStmt)
	target, ok := stmt.Target.(*ast.SubqueryRef)
	if !ok {
		t.Fatalf("expected subquery target, got %T", stmt.Target)
	}
	if target.Subquery == nil {
		t.Fatal("expected retained target subquery")
	}
	if target.Alias == nil || target.Alias.Name != "V" {
		t.Fatalf("expected target alias V, got %#v", target.Alias)
	}
	if stmt.Table != nil {
		t.Fatalf("expected legacy Table to remain nil for inline view target, got %#v", stmt.Table)
	}
}

func TestParseUpdateInlineViewTargetWithRestriction(t *testing.T) {
	result := ParseAndCheck(t, "UPDATE (SELECT a FROM t WITH CHECK OPTION CONSTRAINT ck_v) v SET a = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.UpdateStmt)
	target, ok := stmt.Target.(*ast.SubqueryRef)
	if !ok {
		t.Fatalf("expected subquery target, got %T", stmt.Target)
	}
	if target.Restriction == nil {
		t.Fatal("expected retained subquery restriction")
	}
	if !target.Restriction.CheckOption {
		t.Fatalf("expected CHECK OPTION restriction, got %#v", target.Restriction)
	}
	if target.Restriction.ConstraintName == nil || target.Restriction.ConstraintName.Name != "CK_V" {
		t.Fatalf("expected restriction constraint CK_V, got %#v", target.Restriction.ConstraintName)
	}
}

// TestParseUpdateMultipleSet tests UPDATE with multiple SET clauses.
func TestParseUpdateMultipleSet(t *testing.T) {
	result := ParseAndCheck(t, "UPDATE employees SET salary = 5000, name = 'John', status = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.UpdateStmt)
	if stmt.SetClauses == nil || stmt.SetClauses.Len() != 3 {
		t.Fatalf("expected 3 SET clauses, got %d", stmt.SetClauses.Len())
	}
}

// TestParseUpdateReturningInto tests UPDATE with RETURNING INTO.
func TestParseUpdateReturningInto(t *testing.T) {
	result := ParseAndCheck(t, "UPDATE employees SET salary = 5000 WHERE id = 1 RETURNING salary INTO :out_sal")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.UpdateStmt)
	if stmt.Returning == nil || stmt.Returning.Len() == 0 {
		t.Fatal("expected non-empty RETURNING list")
	}
}

// TestParseDeleteSimple tests DELETE FROM t WHERE cond.
func TestParseDeleteSimple(t *testing.T) {
	result := ParseAndCheck(t, "DELETE FROM employees WHERE id = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("expected DeleteStmt, got %T", raw.Stmt)
	}
	if stmt.Table == nil {
		t.Fatal("expected non-nil Table")
	}
	if stmt.WhereClause == nil {
		t.Error("expected non-nil WhereClause")
	}
}

// TestParseDeleteWithAlias tests DELETE with table alias.
func TestParseDeleteWithAlias(t *testing.T) {
	result := ParseAndCheck(t, "DELETE FROM employees e WHERE e.id = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.DeleteStmt)
	if stmt.Alias == nil || stmt.Alias.Name != "E" {
		t.Errorf("expected alias E, got %v", stmt.Alias)
	}
	target, ok := stmt.Target.(*ast.TableRef)
	if !ok {
		t.Fatalf("expected table target, got %T", stmt.Target)
	}
	if target.Name == nil || target.Name.Name != "EMPLOYEES" {
		t.Fatalf("expected target table EMPLOYEES, got %#v", target.Name)
	}
	if target.Alias == nil || target.Alias.Name != "E" {
		t.Fatalf("expected target alias E, got %#v", target.Alias)
	}
}

func TestParseDeleteInlineViewTarget(t *testing.T) {
	result := ParseAndCheck(t, "DELETE FROM (SELECT * FROM t) v WHERE id = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.DeleteStmt)
	target, ok := stmt.Target.(*ast.SubqueryRef)
	if !ok {
		t.Fatalf("expected subquery target, got %T", stmt.Target)
	}
	if target.Subquery == nil {
		t.Fatal("expected retained target subquery")
	}
	if target.Alias == nil || target.Alias.Name != "V" {
		t.Fatalf("expected target alias V, got %#v", target.Alias)
	}
	if stmt.Table != nil {
		t.Fatalf("expected legacy Table to remain nil for inline view target, got %#v", stmt.Table)
	}
}

func TestParseDeleteTableCollectionTarget(t *testing.T) {
	result := ParseAndCheck(t, "DELETE FROM TABLE(items) i WHERE i.id = 1")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.DeleteStmt)
	target, ok := stmt.Target.(*ast.TableCollectionExpr)
	if !ok {
		t.Fatalf("expected table collection target, got %T", stmt.Target)
	}
	if target.Expr == nil {
		t.Fatal("expected retained collection expression")
	}
	if target.Alias == nil || target.Alias.Name != "I" {
		t.Fatalf("expected target alias I, got %#v", target.Alias)
	}
}

// TestParseDeleteReturningInto tests DELETE with RETURNING INTO.
func TestParseDeleteReturningInto(t *testing.T) {
	result := ParseAndCheck(t, "DELETE FROM employees WHERE id = 1 RETURNING id INTO :out_id")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.DeleteStmt)
	if stmt.Returning == nil || stmt.Returning.Len() == 0 {
		t.Fatal("expected non-empty RETURNING list")
	}
}
