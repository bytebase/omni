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

// --- Task 4.2: HANDLER condition list (direct-AST) -----------------------

func TestValidateHandlerDuplicateCond(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareHandlerStmt{
					Action: "CONTINUE",
					Conditions: []nodes.HandlerCondValue{
						{Kind: nodes.HandlerCondSQLState, Value: "23000"},
						{Kind: nodes.HandlerCondSQLState, Value: "23000"},
					},
					Stmt: &nodes.BeginEndBlock{},
					Loc:  nodes.Loc{Start: 10},
				},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "duplicate_handler_condition")
}

func TestValidateHandlerUndeclaredCond(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareHandlerStmt{
					Action: "EXIT",
					Conditions: []nodes.HandlerCondValue{
						{Kind: nodes.HandlerCondName, Value: "no_such_cond"},
					},
					Stmt: &nodes.BeginEndBlock{},
					Loc:  nodes.Loc{Start: 10},
				},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_condition")
}

func TestValidateHandlerDeclaredCondOK(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareConditionStmt{Name: "cond_ok"},
				&nodes.DeclareHandlerStmt{
					Action: "EXIT",
					Conditions: []nodes.HandlerCondValue{
						{Kind: nodes.HandlerCondName, Value: "COND_OK"},
					},
					Stmt: &nodes.BeginEndBlock{},
				},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_condition")
	requireNoCode(t, diags, "duplicate_handler_condition")
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

// --- Task 4.3: LEAVE / ITERATE label resolution --------------------------

func TestValidateLeaveUndeclaredLabel(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.LeaveStmt{Label: "nope", Loc: nodes.Loc{Start: 10}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_label")
}

func TestValidateLeaveBlockLabelOK(t *testing.T) {
	// LEAVE can target a block label.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{
				Label: "outer",
				Stmts: []nodes.Node{
					&nodes.LeaveStmt{Label: "outer"},
				},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_label")
}

func TestValidateIterateUndeclaredLabel(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.IterateStmt{Label: "nope", Loc: nodes.Loc{Start: 10}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_loop_label")
}

func TestValidateIterateBlockLabelRejected(t *testing.T) {
	// ITERATE inside a loop-less labeled block: block label is NOT a loop
	// label, so ITERATE must be rejected even though the label exists.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{
				Label: "blk",
				Stmts: []nodes.Node{
					&nodes.IterateStmt{Label: "blk"},
				},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_loop_label")
}

// --- Task 4.4: OPEN/FETCH/CLOSE cursor reference -------------------------

func TestValidateOpenUndeclaredCursor(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.OpenCursorStmt{Name: "no_cur", Loc: nodes.Loc{Start: 10}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_cursor")
}

func TestValidateFetchUndeclaredCursor(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.FetchCursorStmt{Name: "no_cur", Loc: nodes.Loc{Start: 10}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_cursor")
}

func TestValidateCloseUndeclaredCursor(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.CloseCursorStmt{Name: "no_cur", Loc: nodes.Loc{Start: 10}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_cursor")
}

func TestValidateCursorDeclaredOK(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareCursorStmt{Name: "cur"},
				&nodes.OpenCursorStmt{Name: "CUR"},
				&nodes.FetchCursorStmt{Name: "cur"},
				&nodes.CloseCursorStmt{Name: "Cur"},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_cursor")
}

// --- Task 4.5: RETURN only inside function body --------------------------

func TestValidateReturnOutsideFunction(t *testing.T) {
	// Procedure body: RETURN is not allowed.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.ReturnStmt{Loc: nodes.Loc{Start: 10}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "return_outside_function")
}

func TestValidateReturnInsideFunction(t *testing.T) {
	// Function body: RETURN is allowed.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: false, // function
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.ReturnStmt{},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "return_outside_function")
}

// --- Task 4.6: Function body must contain RETURN -------------------------

func TestValidateFunctionMissingReturn(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: false, // function
			Body: &nodes.BeginEndBlock{
				Loc: nodes.Loc{Start: 7},
				Stmts: []nodes.Node{
					&nodes.DeclareVarStmt{Names: []string{"x"}},
				},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "function_missing_return")
}

func TestValidateFunctionHasNestedReturn(t *testing.T) {
	// RETURN inside IF/THEN counts toward the HAS_RETURN flag.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: false,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.IfStmt{
					ThenList: []nodes.Node{&nodes.ReturnStmt{}},
				},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "function_missing_return")
}

// --- Task 4.7: SET target resolution (Bug 2 groundwork) ------------------

func TestValidateSetUndeclaredVar(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.Assignment{
					Column: &nodes.ColumnRef{
						Column: "totally_not_a_real_thing",
						Loc:    nodes.Loc{Start: 10},
					},
				},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_variable")
}

func TestValidateSetDeclaredVar(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareVarStmt{Names: []string{"x"}},
				&nodes.Assignment{Column: &nodes.ColumnRef{Column: "X"}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_variable")
}

func TestValidateSetUserVarSkipped(t *testing.T) {
	// @user variables carry the @-prefix in Column; validator must ignore.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.Assignment{Column: &nodes.ColumnRef{Column: "@foo"}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_variable")
}

func TestValidateSetQualifiedSkipped(t *testing.T) {
	// Table-qualified column → treated as a system-variable reference.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.Assignment{Column: &nodes.ColumnRef{Table: "global", Column: "sql_mode"}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_variable")
}

func TestValidateSetOutsideRoutineSkipped(t *testing.T) {
	// A bare Assignment at the top level — no scope — must not emit.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.Assignment{Column: &nodes.ColumnRef{Column: "whatever"}},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_variable")
}

func TestValidateProcedureNoReturnOK(t *testing.T) {
	// Procedures are not required to contain RETURN.
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{
				Stmts: []nodes.Node{&nodes.DeclareVarStmt{Names: []string{"x"}}},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "function_missing_return")
}

func TestValidateIterateLoopLabelOK(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.WhileStmt{
					Label: "lp",
					Stmts: []nodes.Node{
						&nodes.IterateStmt{Label: "lp"},
						&nodes.LeaveStmt{Label: "lp"},
					},
				},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "undeclared_label")
	requireNoCode(t, diags, "undeclared_loop_label")
}

// --- Phase 4 reviewer follow-ups: walker invariants (regression cover) ---

// TestValidateLeaveAcrossHandlerBarrier asserts that LEAVE inside a handler
// body cannot target a label declared in the enclosing routine scope.
// scope.go:lookupLabel stops walking at scopeHandlerBody; this fires
// "undeclared_label".
func TestValidateLeaveAcrossHandlerBarrier(t *testing.T) {
	// outer: BEGIN (label "lbl")
	//   DECLARE EXIT HANDLER FOR SQLEXCEPTION
	//     BEGIN
	//       LEAVE lbl;   -- illegal: handler body cannot see outer label
	//     END;
	// END
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{
				Label: "lbl",
				Stmts: []nodes.Node{
					&nodes.DeclareHandlerStmt{
						Action:     "EXIT",
						Conditions: []nodes.HandlerCondValue{{Kind: nodes.HandlerCondSQLException}},
						Stmt: &nodes.BeginEndBlock{
							Stmts: []nodes.Node{
								&nodes.LeaveStmt{Label: "lbl"},
							},
						},
					},
				},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_label")
}

// TestValidateNestedVarShadowAllowed asserts MySQL's rule that re-declaring
// a variable name in a nested BEGIN/END scope is legal — the inner DECLARE
// shadows the outer; no "duplicate_variable" diagnostic.
func TestValidateNestedVarShadowAllowed(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{
				Stmts: []nodes.Node{
					&nodes.DeclareVarStmt{Names: []string{"x"}},
					&nodes.BeginEndBlock{
						Stmts: []nodes.Node{
							&nodes.DeclareVarStmt{Names: []string{"x"}},
						},
					},
				},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "duplicate_variable")
}

// TestValidateReturnInHandlerInFunction asserts that RETURN inside a handler
// body inside a function body is legal. isFunction must propagate through
// scopeHandlerBody, so no "return_outside_function" fires.
func TestValidateReturnInHandlerInFunction(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: false,
			Returns:     &nodes.DataType{}, // shape-only; validator doesn't care
			Body: &nodes.BeginEndBlock{
				Stmts: []nodes.Node{
					&nodes.DeclareHandlerStmt{
						Action:     "EXIT",
						Conditions: []nodes.HandlerCondValue{{Kind: nodes.HandlerCondSQLException}},
						Stmt:       &nodes.ReturnStmt{Expr: &nodes.IntLit{Value: 0}},
					},
				},
			},
		},
	}}
	diags := Validate(list, Options{})
	requireNoCode(t, diags, "return_outside_function")
}
