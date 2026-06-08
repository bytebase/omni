package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// TestExceptionSingleWhen tests: EXCEPTION WHEN division_by_zero THEN x := 0;
func TestExceptionSingleWhen(t *testing.T) {
	block := parseOK(t, `BEGIN x := 1; EXCEPTION WHEN division_by_zero THEN x := 0; END`)

	if len(block.Body) != 1 {
		t.Fatalf("expected 1 body statement, got %d", len(block.Body))
	}
	if len(block.Exceptions) != 1 {
		t.Fatalf("expected 1 exception handler, got %d", len(block.Exceptions))
	}
	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Conditions) != 1 || w.Conditions[0] != "division_by_zero" {
		t.Fatalf("expected condition [division_by_zero], got %v", w.Conditions)
	}
	if len(w.Body) != 1 {
		t.Fatalf("expected 1 handler statement, got %d", len(w.Body))
	}
	assign := w.Body[0].(*ast.PLAssign)
	if assign.Target != "x" || assign.Expr != "0" {
		t.Fatalf("expected x := 0, got %s := %s", assign.Target, assign.Expr)
	}
}

// TestExceptionMultipleWhen tests multiple WHEN clauses.
func TestExceptionMultipleWhen(t *testing.T) {
	block := parseOK(t, `
		BEGIN
			x := 1;
		EXCEPTION
			WHEN no_data_found THEN x := -1;
			WHEN too_many_rows THEN x := -2;
		END
	`)

	if len(block.Exceptions) != 2 {
		t.Fatalf("expected 2 exception handlers, got %d", len(block.Exceptions))
	}
	w0 := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w0.Conditions) != 1 || w0.Conditions[0] != "no_data_found" {
		t.Fatalf("handler 0: expected [no_data_found], got %v", w0.Conditions)
	}
	w1 := block.Exceptions[1].(*ast.PLExceptionWhen)
	if len(w1.Conditions) != 1 || w1.Conditions[0] != "too_many_rows" {
		t.Fatalf("handler 1: expected [too_many_rows], got %v", w1.Conditions)
	}
}

// TestExceptionORConditions tests OR-joined conditions.
func TestExceptionORConditions(t *testing.T) {
	block := parseOK(t, `BEGIN NULL; EXCEPTION WHEN division_by_zero OR unique_violation THEN x := 0; END`)

	if len(block.Exceptions) != 1 {
		t.Fatalf("expected 1 exception handler, got %d", len(block.Exceptions))
	}
	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(w.Conditions))
	}
	if w.Conditions[0] != "division_by_zero" {
		t.Fatalf("condition 0: expected division_by_zero, got %s", w.Conditions[0])
	}
	if w.Conditions[1] != "unique_violation" {
		t.Fatalf("condition 1: expected unique_violation, got %s", w.Conditions[1])
	}
}

// TestExceptionSQLSTATE tests SQLSTATE code condition.
func TestExceptionSQLSTATE(t *testing.T) {
	block := parseOK(t, `BEGIN NULL; EXCEPTION WHEN SQLSTATE '22012' THEN x := 0; END`)

	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Conditions) != 1 || w.Conditions[0] != "SQLSTATE '22012'" {
		t.Fatalf("expected [SQLSTATE '22012'], got %v", w.Conditions)
	}
}

// TestExceptionOTHERS tests OTHERS catch-all.
func TestExceptionOTHERS(t *testing.T) {
	block := parseOK(t, `BEGIN NULL; EXCEPTION WHEN OTHERS THEN x := 0; END`)

	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Conditions) != 1 || w.Conditions[0] != "others" {
		t.Fatalf("expected [others], got %v", w.Conditions)
	}
}

// TestExceptionMultipleStatements tests multiple statements in handler body.
func TestExceptionMultipleStatements(t *testing.T) {
	block := parseOK(t, `
		BEGIN
			x := 1;
		EXCEPTION
			WHEN division_by_zero THEN
				x := 0;
				y := 'error';
				z := x + y;
		END
	`)

	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Body) != 3 {
		t.Fatalf("expected 3 handler statements, got %d", len(w.Body))
	}
}

// TestExceptionGetStackedDiag tests GET STACKED DIAGNOSTICS inside handler.
func TestExceptionGetStackedDiag(t *testing.T) {
	block := parseOK(t, `
		BEGIN
			x := 1;
		EXCEPTION
			WHEN OTHERS THEN
				GET STACKED DIAGNOSTICS msg = MESSAGE_TEXT;
		END
	`)

	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Body) != 1 {
		t.Fatalf("expected 1 handler statement, got %d", len(w.Body))
	}
	diag := w.Body[0].(*ast.PLGetDiag)
	if !diag.IsStacked {
		t.Fatal("expected GET STACKED DIAGNOSTICS")
	}
}

// TestExceptionRaiseInsideHandler tests RAISE inside handler (re-raise).
func TestExceptionRaiseInsideHandler(t *testing.T) {
	// Note: RAISE is not yet implemented, so we test with a simple assignment
	// that simulates what would happen. Once RAISE is implemented, this can
	// be updated. For now, we verify the structure works with existing stmts.
	block := parseOK(t, `
		BEGIN
			x := 1;
		EXCEPTION
			WHEN OTHERS THEN
				x := -1;
		END
	`)

	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Body) != 1 {
		t.Fatalf("expected 1 handler statement, got %d", len(w.Body))
	}
}

// TestExceptionNestedBlock tests nested block with own exception handler.
func TestExceptionNestedBlock(t *testing.T) {
	block := parseOK(t, `
		BEGIN
			BEGIN
				x := 1;
			EXCEPTION
				WHEN division_by_zero THEN x := 0;
			END;
		EXCEPTION
			WHEN OTHERS THEN x := -1;
		END
	`)

	// Outer block has 1 body stmt (the nested block) and 1 exception handler
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 outer body statement, got %d", len(block.Body))
	}
	if len(block.Exceptions) != 1 {
		t.Fatalf("expected 1 outer exception handler, got %d", len(block.Exceptions))
	}

	// Inner block
	inner := block.Body[0].(*ast.PLBlock)
	if len(inner.Exceptions) != 1 {
		t.Fatalf("expected 1 inner exception handler, got %d", len(inner.Exceptions))
	}
	innerW := inner.Exceptions[0].(*ast.PLExceptionWhen)
	if innerW.Conditions[0] != "division_by_zero" {
		t.Fatalf("inner handler: expected division_by_zero, got %s", innerW.Conditions[0])
	}
}

// TestExceptionEmptyBody tests exception block with no statements before EXCEPTION.
func TestExceptionEmptyBody(t *testing.T) {
	block := parseOK(t, `BEGIN EXCEPTION WHEN OTHERS THEN x := 0; END`)

	if len(block.Body) != 0 {
		t.Fatalf("expected 0 body statements, got %d", len(block.Body))
	}
	if len(block.Exceptions) != 1 {
		t.Fatalf("expected 1 exception handler, got %d", len(block.Exceptions))
	}
}

// TestExceptionSQLSTATEAndNamedOR tests SQLSTATE and named conditions combined via OR.
func TestExceptionSQLSTATEAndNamedOR(t *testing.T) {
	block := parseOK(t, `BEGIN NULL; EXCEPTION WHEN SQLSTATE '23505' OR unique_violation THEN x := 0; END`)

	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(w.Conditions))
	}
	if w.Conditions[0] != "SQLSTATE '23505'" {
		t.Fatalf("condition 0: expected SQLSTATE '23505', got %s", w.Conditions[0])
	}
	if w.Conditions[1] != "unique_violation" {
		t.Fatalf("condition 1: expected unique_violation, got %s", w.Conditions[1])
	}
}
