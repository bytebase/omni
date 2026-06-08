package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// This file is a TEST-ONLY negative/coverage backfill for the PartiQL
// ON CONFLICT and RETURNING grammar (the onConflictClause / conflictTarget /
// conflictAction / returningClause / returningColumn productions of
// PartiQLParser.g4, parsed in partiql/parser/dml.go).
//
// Oracle: every accept/reject verdict below was confirmed against the
// EXECUTABLE generated ANTLR PartiQL parser (github.com/bytebase/parser/partiql,
// `Script()` entry rule) via a throwaway differential probe, cross-checked
// with the grammar text. For these productions omni and the ANTLR oracle
// agree on every case — there are no divergences in this file.
//
// Grammar (PartiQLParser.g4):
//
//	onConflictClause
//	    : ON CONFLICT WHERE expr DO NOTHING                  # OnConflictLegacy
//	    | ON CONFLICT conflictTarget? conflictAction         # OnConflict
//	conflictTarget
//	    : PAREN_LEFT symbolPrimitive (COMMA symbolPrimitive)* PAREN_RIGHT
//	    | ON CONSTRAINT constraintName
//	conflictAction
//	    : DO NOTHING | DO REPLACE EXCLUDED | DO UPDATE EXCLUDED
//	returningClause
//	    : RETURNING returningColumn ( COMMA returningColumn )*
//	returningColumn
//	    : status=(MODIFIED|ALL) age=(OLD|NEW) ASTERISK
//	    | status=(MODIFIED|ALL) age=(OLD|NEW) col=expr
//
// Two oracle-derived nuances worth calling out, because a casual reading of
// the feature list gets them backwards:
//   - `RETURNING MODIFIED *` is REJECTED: every returningColumn REQUIRES an
//     OLD/NEW between the status and the `*`/expr. `MODIFIED *` skips it.
//   - `RETURNING MODIFIED NEW a, b` is REJECTED: each comma-separated column
//     is a full returningColumn, so the bare `b` (no status/mapping) is a
//     syntax error. A well-formed multi-column list repeats the prefix, e.g.
//     `RETURNING MODIFIED NEW a, ALL OLD b`.
// Both are confirmed against the ANTLR oracle (omni matches it), so they live
// among the negative cases below rather than the positives.

// TestParser_OnConflictAccepts covers the conflictTarget / conflictAction
// forms the omni parser implements: column-list target, ON CONSTRAINT target,
// the legacy `WHERE expr DO NOTHING`, and the DO REPLACE/UPDATE EXCLUDED
// actions — both on the legacy `INSERT ... VALUE` form (insertCommandReturning
// / InsertLegacy) and on the RFC-0011 bare-symbol form (# Insert), which also
// permits onConflictClause. All confirmed ACCEPTED by the ANTLR oracle.
func TestParser_OnConflictAccepts(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantStr string // substring NodeToString must contain
	}{
		{
			name:    "target_cols_single_do_nothing",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT (a) DO NOTHING",
			wantStr: "OnConflict{Target:OnConflictTarget{Cols:[VarRef{Name:a}]} Action:DO_NOTHING}",
		},
		{
			name:    "target_cols_multi_do_nothing",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT (a, b) DO NOTHING",
			wantStr: "OnConflictTarget{Cols:[VarRef{Name:a} VarRef{Name:b}]} Action:DO_NOTHING}",
		},
		{
			name:    "target_on_constraint_do_nothing",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT ON CONSTRAINT my_c DO NOTHING",
			wantStr: "OnConflictTarget{ConstraintName:my_c} Action:DO_NOTHING}",
		},
		{
			// legacy onConflictClause: ON CONFLICT WHERE expr DO NOTHING.
			name:    "legacy_where_do_nothing",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT WHERE a > 0 DO NOTHING",
			wantStr: "OnConflict{Action:DO_NOTHING Where:BinaryExpr{Op:> Left:VarRef{Name:a} Right:NumberLit{Val:0}}}",
		},
		{
			name:    "action_do_replace_excluded",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO REPLACE EXCLUDED",
			wantStr: "OnConflict{Action:DO_REPLACE_EXCLUDED}",
		},
		{
			name:    "action_do_update_excluded",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO UPDATE EXCLUDED",
			wantStr: "OnConflict{Action:DO_UPDATE_EXCLUDED}",
		},
		{
			name:    "target_cols_do_replace_excluded",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT (a) DO REPLACE EXCLUDED",
			wantStr: "OnConflictTarget{Cols:[VarRef{Name:a}]} Action:DO_REPLACE_EXCLUDED}",
		},
		{
			name:    "target_on_constraint_do_replace_excluded",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT ON CONSTRAINT c DO REPLACE EXCLUDED",
			wantStr: "OnConflictTarget{ConstraintName:c} Action:DO_REPLACE_EXCLUDED}",
		},
		{
			// RFC-0011 bare-symbol INSERT (no VALUE keyword) also allows
			// onConflictClause (# Insert : ... onConflictClause?).
			name:    "rfc0011_bare_symbol_do_nothing",
			input:   "INSERT INTO foo {'a': 1} ON CONFLICT DO NOTHING",
			wantStr: "InsertStmt{Target:VarRef{Name:foo} Value:TupleLit",
		},
		{
			name:    "rfc0011_bag_value_target_cols",
			input:   "INSERT INTO foo << {'a': 1} >> ON CONFLICT (a) DO NOTHING",
			wantStr: "OnConflictTarget{Cols:[VarRef{Name:a}]} Action:DO_NOTHING}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q) unexpected error: %v", tc.input, err)
			}
			if _, ok := stmt.(*ast.InsertStmt); !ok {
				t.Fatalf("ParseStatement(%q) = %T, want *ast.InsertStmt", tc.input, stmt)
			}
			got := ast.NodeToString(stmt)
			if !strings.Contains(got, tc.wantStr) {
				t.Errorf("NodeToString = %q, want to contain %q", got, tc.wantStr)
			}
		})
	}
}

// TestParser_ReturningAccepts covers the returningColumn forms the parser
// implements: MODIFIED/ALL × OLD/NEW with `*`, a single non-star expr, and a
// well-formed multi-column list (each column repeats its status+mapping).
// RETURNING is reachable on the legacy `INSERT ... VALUE` form
// (insertCommandReturning) and on DELETE (deleteCommand shares
// returningClause). All confirmed ACCEPTED by the ANTLR oracle.
func TestParser_ReturningAccepts(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantStr string
	}{
		{
			name:    "all_old_star",
			input:   "INSERT INTO foo VALUE {'a': 1} RETURNING ALL OLD *",
			wantStr: "ReturningClause{Items:[ReturningItem{Status:ALL Mapping:OLD Star:true}]}",
		},
		{
			name:    "all_new_star",
			input:   "INSERT INTO foo VALUE {'a': 1} RETURNING ALL NEW *",
			wantStr: "ReturningItem{Status:ALL Mapping:NEW Star:true}",
		},
		{
			name:    "modified_old_star",
			input:   "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED OLD *",
			wantStr: "ReturningItem{Status:MODIFIED Mapping:OLD Star:true}",
		},
		{
			// non-star expr column.
			name:    "modified_new_expr",
			input:   "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED NEW a",
			wantStr: "ReturningItem{Status:MODIFIED Mapping:NEW Expr:VarRef{Name:a}}",
		},
		{
			name:    "modified_new_expr_arith",
			input:   "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED NEW a + 1",
			wantStr: "ReturningItem{Status:MODIFIED Mapping:NEW Expr:BinaryExpr{Op:+",
		},
		{
			// well-formed multi-column list: each column repeats status+mapping.
			name:    "multi_column_each_prefixed",
			input:   "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED NEW a, ALL OLD b",
			wantStr: "Items:[ReturningItem{Status:MODIFIED Mapping:NEW Expr:VarRef{Name:a}} ReturningItem{Status:ALL Mapping:OLD Expr:VarRef{Name:b}}]",
		},
		{
			// RETURNING is shared by deleteCommand; pin both the DeleteStmt
			// shape and the returning item so a regression in either surfaces.
			name:    "delete_returning_star",
			input:   "DELETE FROM foo WHERE a = 1 RETURNING MODIFIED OLD *",
			wantStr: "ReturningItem{Status:MODIFIED Mapping:OLD Star:true}",
		},
		{
			name:    "delete_returning_expr",
			input:   "DELETE FROM foo RETURNING ALL OLD a",
			wantStr: "ReturningItem{Status:ALL Mapping:OLD Expr:VarRef{Name:a}}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q) unexpected error: %v", tc.input, err)
			}
			got := ast.NodeToString(stmt)
			if !strings.Contains(got, "Returning:ReturningClause{") {
				t.Errorf("NodeToString = %q, want a Returning clause", got)
			}
			if !strings.Contains(got, tc.wantStr) {
				t.Errorf("NodeToString = %q, want to contain %q", got, tc.wantStr)
			}
		})
	}
}

// TestParser_OnConflictAndReturningTogether covers the legacy
// insertCommandReturning form where an onConflictClause is followed by a
// well-formed returningClause (`... onConflictClause? returningClause?`). Both
// confirmed ACCEPTED by the ANTLR oracle. (Note the RETURNING must itself be
// well-formed: `... DO NOTHING RETURNING MODIFIED *` is REJECTED because
// `MODIFIED *` lacks OLD/NEW — see TestParser_ReturningRejects.)
func TestParser_OnConflictAndReturningTogether(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantStr string
	}{
		{
			name:    "action_then_returning",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO UPDATE EXCLUDED RETURNING ALL NEW *",
			wantStr: "OnConflict{Action:DO_UPDATE_EXCLUDED} Returning:ReturningClause{Items:[ReturningItem{Status:ALL Mapping:NEW Star:true}]}",
		},
		{
			name:    "target_action_then_returning",
			input:   "INSERT INTO foo VALUE {'a': 1} ON CONFLICT (a) DO NOTHING RETURNING ALL OLD *",
			wantStr: "OnConflict{Target:OnConflictTarget{Cols:[VarRef{Name:a}]} Action:DO_NOTHING} Returning:ReturningClause{Items:[ReturningItem{Status:ALL Mapping:OLD Star:true}]}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q) unexpected error: %v", tc.input, err)
			}
			got := ast.NodeToString(stmt)
			if !strings.Contains(got, tc.wantStr) {
				t.Errorf("NodeToString = %q, want to contain %q", got, tc.wantStr)
			}
		})
	}
}

// TestParser_OnConflictRejects covers malformed onConflictClause inputs the
// ANTLR oracle REJECTS (syntax error); omni must reject them too. Without
// these the conflictAction parser would be over-permissive.
//
// The asserted substring is omni's message; the contract under test is the
// REJECT verdict (which the ANTLR oracle agrees with), not the exact wording.
func TestParser_OnConflictRejects(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			// DO <garbage>: not NOTHING/REPLACE/UPDATE.
			name:      "do_garbage",
			input:     "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO garbage",
			wantErrIn: "expected NOTHING, REPLACE, or UPDATE after DO",
		},
		{
			// DO REPLACE requires EXCLUDED (doReplace : EXCLUDED).
			name:      "do_replace_missing_excluded",
			input:     "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO REPLACE",
			wantErrIn: "expected EXCLUDED",
		},
		{
			// DO UPDATE requires EXCLUDED (doUpdate : EXCLUDED).
			name:      "do_update_missing_excluded",
			input:     "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO UPDATE",
			wantErrIn: "expected EXCLUDED",
		},
		{
			// The SET / VALUE forms of doReplace are commented out in the
			// grammar (":TODO add the rest"), so `DO REPLACE SET ...` is a
			// syntax error — EXCLUDED is the only legal token after REPLACE.
			name:      "do_replace_set_not_implemented",
			input:     "INSERT INTO foo VALUE {'a': 1} ON CONFLICT DO REPLACE SET a = 1",
			wantErrIn: "expected EXCLUDED",
		},
		{
			// conflictAction is mandatory: `ON CONFLICT` alone (no action) is
			// rejected (the # OnConflict alt is `conflictTarget? conflictAction`).
			name:      "missing_action",
			input:     "INSERT INTO foo VALUE {'a': 1} ON CONFLICT",
			wantErrIn: "expected DO",
		},
		{
			// A conflictTarget without a following conflictAction is rejected.
			name:      "target_without_action",
			input:     "INSERT INTO foo VALUE {'a': 1} ON CONFLICT (a)",
			wantErrIn: "expected DO",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			if _, err := p.ParseStatement(); err == nil {
				t.Fatalf("ParseStatement(%q): expected error, got nil", tc.input)
			} else if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_ReturningRejects covers malformed returningColumn inputs the
// ANTLR oracle REJECTS; omni must reject them too. Each returningColumn is
// `status=(MODIFIED|ALL) age=(OLD|NEW) (ASTERISK | expr)`, so a missing
// status, a missing OLD/NEW, or a bare second column are all syntax errors.
//
// The asserted substring is omni's message; the contract under test is the
// REJECT verdict, confirmed against the ANTLR oracle.
func TestParser_ReturningRejects(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// ----- missing MODIFIED/ALL status -----
		{
			name:      "missing_status_old_star",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING OLD *",
			wantErrIn: "expected MODIFIED or ALL",
		},
		{
			name:      "missing_status_bare_star",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING *",
			wantErrIn: "expected MODIFIED or ALL",
		},
		{
			name:      "missing_status_bare_expr",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING a",
			wantErrIn: "expected MODIFIED or ALL",
		},
		// ----- missing OLD/NEW mapping -----
		{
			// `MODIFIED *` — no OLD/NEW between status and `*`. REJECTED by the
			// oracle, despite "MODIFIED *" reading like a natural shorthand.
			name:      "modified_missing_mapping_star",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED *",
			wantErrIn: "expected OLD or NEW",
		},
		{
			name:      "all_missing_mapping_star",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING ALL *",
			wantErrIn: "expected OLD or NEW",
		},
		{
			name:      "modified_missing_mapping_expr",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED a",
			wantErrIn: "expected OLD or NEW",
		},
		{
			name:      "modified_missing_mapping_eof",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED",
			wantErrIn: "expected OLD or NEW",
		},
		// ----- malformed multi-column list -----
		{
			// `MODIFIED NEW a, b`: the second column `b` has no status/mapping
			// prefix, so it is not a valid returningColumn. REJECTED.
			name:      "multi_column_second_bare",
			input:     "INSERT INTO foo VALUE {'a': 1} RETURNING MODIFIED NEW a, b",
			wantErrIn: "expected MODIFIED or ALL",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			if _, err := p.ParseStatement(); err == nil {
				t.Fatalf("ParseStatement(%q): expected error, got nil", tc.input)
			} else if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_RFC0011RejectsTrailingReturningOrConflict pins a subtle grammar
// boundary: in the top-level `dml` rule, only the legacy `INSERT ... VALUE`
// form (insertCommandReturning, the # DmlInsertReturning alt) carries a
// trailing returningClause. The RFC-0011 bare-symbol INSERT, plus REPLACE and
// UPSERT, are reachable only via `dmlBaseCommand # DmlBase`, which has NO
// returningClause — and only the RFC-0011 INSERT (# Insert) carries an
// onConflictClause; REPLACE/UPSERT carry neither.
//
// Therefore the ANTLR oracle REJECTS every input below (the trailing
// RETURNING / ON CONFLICT cannot attach). omni rejects them at the EOF
// assertion in ParseStatement with "unexpected token <tok> after statement"
// (the RFC-0011 INSERT/REPLACE/UPSERT value parse completes, then the leftover
// keyword is surfaced). Verdict confirmed against the ANTLR oracle.
func TestParser_RFC0011RejectsTrailingReturningOrConflict(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			// RFC-0011 INSERT (no VALUE) + RETURNING: illegal.
			name:      "rfc0011_insert_returning",
			input:     "INSERT INTO foo {'a': 1} RETURNING MODIFIED NEW a",
			wantErrIn: `unexpected token "RETURNING" after statement`,
		},
		{
			// RFC-0011 INSERT may take ON CONFLICT, but NOT a trailing
			// RETURNING after it — the RETURNING is left dangling.
			name:      "rfc0011_insert_conflict_then_returning",
			input:     "INSERT INTO foo {'a': 1} ON CONFLICT DO NOTHING RETURNING ALL NEW *",
			wantErrIn: `unexpected token "RETURNING" after statement`,
		},
		{
			// REPLACE carries neither onConflictClause nor returningClause.
			name:      "replace_returning",
			input:     "REPLACE INTO foo {'a': 1} RETURNING ALL NEW *",
			wantErrIn: `unexpected token "RETURNING" after statement`,
		},
		{
			name:      "replace_on_conflict",
			input:     "REPLACE INTO foo {'a': 1} ON CONFLICT DO NOTHING",
			wantErrIn: `unexpected token "ON" after statement`,
		},
		{
			// UPSERT likewise carries neither.
			name:      "upsert_returning",
			input:     "UPSERT INTO foo {'a': 1} RETURNING ALL NEW *",
			wantErrIn: `unexpected token "RETURNING" after statement`,
		},
		{
			name:      "upsert_on_conflict",
			input:     "UPSERT INTO foo {'a': 1} ON CONFLICT DO NOTHING",
			wantErrIn: `unexpected token "ON" after statement`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			if _, err := p.ParseStatement(); err == nil {
				t.Fatalf("ParseStatement(%q): expected error, got nil", tc.input)
			} else if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
