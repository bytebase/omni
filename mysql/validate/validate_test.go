package validate

import (
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
	parser "github.com/bytebase/omni/mysql/parser"
)

// requireCode asserts that diags contains at least one Diagnostic with the
// given code. Stable across reorderings.
func requireCode(t *testing.T, diags []Diagnostic, code string) {
	t.Helper()
	for _, d := range diags {
		if d.Code == code {
			return
		}
	}
	t.Fatalf("expected diagnostic code %q, got %v", code, diags)
}

// requireNoCode asserts that diags contains no Diagnostic with the given code.
func requireNoCode(t *testing.T, diags []Diagnostic, code string) {
	t.Helper()
	for _, d := range diags {
		if d.Code == code {
			t.Fatalf("unexpected diagnostic code %q: %v", code, d)
		}
	}
}

func TestValidateEmpty(t *testing.T) {
	diags := Validate(nil, Options{})
	if diags != nil {
		t.Fatalf("expected nil diagnostics, got %v", diags)
	}
}

func TestValidateCleanProcedure(t *testing.T) {
	list, err := parser.Parse(`CREATE PROCEDURE p() BEGIN DECLARE x INT; SET x = 1; END`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	diags := Validate(list, Options{})
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// --- Task 4.1: duplicate DECLARE / label (direct-AST) --------------------

func TestValidateDuplicateVarDirect(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareVarStmt{Names: []string{"x"}},
				&nodes.DeclareVarStmt{Names: []string{"x"}, Loc: nodes.Loc{Start: 42}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "duplicate_variable")
}

func TestValidateDuplicateConditionDirect(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareConditionStmt{Name: "c1"},
				&nodes.DeclareConditionStmt{Name: "C1", Loc: nodes.Loc{Start: 42}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "duplicate_condition")
}

func TestValidateDuplicateCursorDirect(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareCursorStmt{Name: "cur"},
				&nodes.DeclareCursorStmt{Name: "CUR", Loc: nodes.Loc{Start: 42}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "duplicate_cursor")
}

func TestValidateDuplicateLabelDirect(t *testing.T) {
	// Two sibling labeled BEGIN..END blocks inside the same routine scope.
	// Both labels land in the enclosing (routine) scope → collision.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.BeginEndBlock{Label: "lbl"},
				&nodes.BeginEndBlock{Label: "LBL", Loc: nodes.Loc{Start: 42}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "duplicate_label")
}
