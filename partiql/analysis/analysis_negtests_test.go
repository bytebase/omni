package analysis

import (
	"strings"
	"testing"
)

// This file backfills negative/coverage cases for the DQL-only validator
// ValidateQuery so that every non-DQL AST kind it switches over has an
// explicit reject test, plus the two remaining DQL accept shapes (EXCEPT
// set-op and EXPLAIN-with-options). The existing analysis_test.go already
// covers INSERT/DELETE/UPDATE (DML), CREATE TABLE / DROP TABLE (DDL),
// EXEC, and the UNION/INTERSECT/EXPLAIN accept paths; this file fills the
// gaps (UPSERT / REPLACE / REMOVE / CREATE INDEX / DROP INDEX, a
// multi-statement reject loop, and the EXCEPT + EXPLAIN(...) accepts).
//
// Oracle: the executable generated ANTLR parser at
// /Users/h3n4l/OpenSource/parser/partiql plus the legacy
// PartiQL{Parser,Lexer}.g4 rules. Each input below was run through both the
// omni parser and the ANTLR Script rule via a throwaway differential probe;
// omni and ANTLR agreed on accept/reject for every case (zero divergences),
// and the AST node type each input parses to was confirmed to land in the
// expected validateStatement switch arm. "accept" for ValidateQuery means a
// nil error; "reject" means a non-nil error whose message names the rejected
// category (DML / DDL / EXEC), matching validateStatement in analysis.go.
//
// Value-expression note: the RFC-0011 UPSERT/REPLACE grammar is
// `UPSERT|REPLACE INTO symbolPrimitive asIdent? value=expr` — there is NO
// `VALUE` keyword on those commands (unlike the INSERT legacy form). So the
// value is a bare expression; `UPSERT INTO t VALUE 1` is rejected by BOTH
// omni and ANTLR ("unexpected token VALUE in expression"), and the correct
// accepted shape is `UPSERT INTO t 1`. The tests use the bare-expression
// form so the statement parses and ValidateQuery reaches its reject arm.

// --- ValidateQuery: remaining non-DQL DML kinds rejected ---

// TestValidateQuery_DMLExtraRejected covers the PartiQL-specific DML kinds
// not exercised by analysis_test.go's DML test: UPSERT, REPLACE, REMOVE.
// Each parses to *ast.UpsertStmt / *ast.ReplaceStmt / *ast.RemoveStmt, which
// validateStatement classifies as DML and rejects.
func TestValidateQuery_DMLExtraRejected(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		// upsertCommand: UPSERT INTO symbolPrimitive value=expr (no VALUE kw).
		{"upsert", "UPSERT INTO t 1"},
		// replaceCommand: REPLACE INTO symbolPrimitive value=expr.
		{"replace", "REPLACE INTO t 1"},
		// removeCommand: REMOVE pathSimple (a bare path).
		{"remove", "REMOVE t"},
		// REMOVE with a path step still classifies as DML.
		{"remove_path_step", "REMOVE t.a"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateQuery(tc.input)
			if err == nil {
				t.Fatalf("ValidateQuery(%q) = nil, want DML error", tc.input)
			}
			if !strings.Contains(err.Error(), "DML") {
				t.Errorf("ValidateQuery(%q) error = %q, want to contain %q", tc.input, err.Error(), "DML")
			}
		})
	}
}

// --- ValidateQuery: remaining non-DQL DDL kinds rejected ---

// TestValidateQuery_DDLExtraRejected covers the index DDL not exercised by
// analysis_test.go's DDL test: CREATE INDEX and DROP INDEX. These parse to
// *ast.CreateIndexStmt / *ast.DropIndexStmt, which validateStatement
// classifies as DDL and rejects.
func TestValidateQuery_DDLExtraRejected(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		// createCommand#CreateIndex: CREATE INDEX ON sym ( pathSimple ,... ).
		{"create_index_single", "CREATE INDEX ON t (a)"},
		{"create_index_multi", "CREATE INDEX ON t (a, b)"},
		// dropCommand#DropIndex: DROP INDEX target ON on.
		{"drop_index", "DROP INDEX i ON t"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateQuery(tc.input)
			if err == nil {
				t.Fatalf("ValidateQuery(%q) = nil, want DDL error", tc.input)
			}
			if !strings.Contains(err.Error(), "DDL") {
				t.Errorf("ValidateQuery(%q) error = %q, want to contain %q", tc.input, err.Error(), "DDL")
			}
		})
	}
}

// --- ValidateQuery: multi-statement reject loop ---

// TestValidateQuery_MultiStatementRejectsNonDQL asserts the per-statement
// loop in ValidateQuery rejects a script when ANY single statement is
// non-DQL, even though the others are valid DQL. The script parses to
// [*ast.SelectStmt, *ast.UpdateStmt]; the second statement trips the DML
// reject arm.
//
// Note on the input: the task brief names `SELECT *; UPDATE t SET a=1`, but a
// bare `SELECT *` (no FROM) is rejected by BOTH omni and ANTLR
// ("expected FROM") — in PartiQL the project list of a top-level SELECT
// requires a FROM clause to parse. So this uses the parse-valid
// `SELECT * FROM t; UPDATE t SET a=1`, which preserves the intent: a
// multi-statement script with one non-DQL statement must be rejected.
func TestValidateQuery_MultiStatementRejectsNonDQL(t *testing.T) {
	const input = "SELECT * FROM t; UPDATE t SET a=1"
	err := ValidateQuery(input)
	if err == nil {
		t.Fatalf("ValidateQuery(%q) = nil, want DML error from the second statement", input)
	}
	if !strings.Contains(err.Error(), "DML") {
		t.Errorf("ValidateQuery(%q) error = %q, want to contain %q", input, err.Error(), "DML")
	}
}

// TestValidateQuery_MultiStatementAllDQL is the positive companion: a script
// of only DQL statements is accepted (the loop returns nil after visiting
// every item).
func TestValidateQuery_MultiStatementAllDQL(t *testing.T) {
	const input = "SELECT * FROM t; SELECT a FROM s WHERE a = 1"
	if err := ValidateQuery(input); err != nil {
		t.Errorf("ValidateQuery(%q) = %v, want nil", input, err)
	}
}

// --- ValidateQuery: remaining DQL accept shapes ---

// TestValidateQuery_ExceptAccepted covers the EXCEPT set operation, the one
// exprBagOp variant not in analysis_test.go's accept list (which has UNION
// and INTERSECT). EXCEPT parses to *ast.SetOpStmt, which validateStatement
// accepts as DQL.
func TestValidateQuery_ExceptAccepted(t *testing.T) {
	cases := []string{
		"SELECT * FROM t1 EXCEPT SELECT * FROM t2",
		"SELECT * FROM t1 EXCEPT ALL SELECT * FROM t2",
		"SELECT * FROM t1 EXCEPT DISTINCT SELECT * FROM t2",
	}
	for _, input := range cases {
		input := input
		t.Run(input[:min(len(input), 30)], func(t *testing.T) {
			if err := ValidateQuery(input); err != nil {
				t.Errorf("ValidateQuery(%q) = %v, want nil", input, err)
			}
		})
	}
}

// TestValidateQuery_ExplainWithOptionsAccepted covers the EXPLAIN form with a
// parenthesised options list — `EXPLAIN (param value, ...) SELECT ...`, per
// the root rule `(EXPLAIN (PAREN_LEFT explainOption (COMMA explainOption)*
// PAREN_RIGHT)?)? statement` where `explainOption : IDENTIFIER IDENTIFIER`.
// analysis_test.go only covers bare `EXPLAIN <stmt>`. This parses to
// *ast.ExplainStmt (options are parsed and discarded), which is accepted
// regardless of the inner statement.
func TestValidateQuery_ExplainWithOptionsAccepted(t *testing.T) {
	cases := []string{
		"EXPLAIN (a b) SELECT * FROM t",
		"EXPLAIN (a b, c d) SELECT * FROM t",
		// EXPLAIN stays read-only even wrapping DML, with options present.
		"EXPLAIN (a b) UPDATE t SET a = 1",
	}
	for _, input := range cases {
		input := input
		t.Run(input[:min(len(input), 30)], func(t *testing.T) {
			if err := ValidateQuery(input); err != nil {
				t.Errorf("ValidateQuery(%q) = %v, want nil", input, err)
			}
		})
	}
}
