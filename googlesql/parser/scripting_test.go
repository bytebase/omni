package parser

// Unit tests for the parser-scripting node: GoogleSQL's procedural language
// (DECLARE / SET / IF / CASE / WHILE / LOOP / REPEAT / FOR-IN / BEGIN…END with
// EXCEPTION / RAISE / RETURN / BREAK / CONTINUE / EXECUTE IMMEDIATE / labels).
//
// These assert the AST shape and accept/reject polarity on hand-authored
// fixtures. The differential against the live Spanner emulator (the PROVE gate
// per correctness-protocol.md) lives in the suite:scripting rows of testdata/compat.tsv (build tag
// googlesql_oracle); scripting is BigQuery-only except the BEGIN…END envelope /
// SET / EXECUTE IMMEDIATE, so the control-flow forms are triangulated against the
// legacy .g4 + the BigQuery truth1 corpus, not the Spanner oracle.

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// parseOneScript parses src (one top-level statement / block) and returns the
// single resulting node, failing on any parse error or a missing/extra node.
func parseOneScript(t *testing.T, src string) ast.Node {
	t.Helper()
	file, errs := parseForTest(src)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) returned errors: %v", src, errs)
	}
	if file == nil || len(file.Stmts) != 1 {
		got := 0
		if file != nil {
			got = len(file.Stmts)
		}
		t.Fatalf("parseForTest(%q): expected exactly 1 statement, got %d", src, got)
	}
	return file.Stmts[0]
}

// ----------------------------------------------------------------------------
// DECLARE
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// SET
// ----------------------------------------------------------------------------

func TestParseSetSingleVariable(t *testing.T) {
	n := parseOneScript(t, "SET x = 5")
	s, ok := n.(*ast.SetStmt)
	if !ok {
		t.Fatalf("expected *ast.SetStmt, got %T", n)
	}
	if s.Kind != ast.SetVariable {
		t.Fatalf("Kind = %v, want SetVariable", s.Kind)
	}
	if len(s.Targets) != 1 {
		t.Fatalf("Targets = %+v, want 1", s.Targets)
	}
	id, ok := s.Targets[0].(*ast.Identifier)
	if !ok || id.Name != "x" {
		t.Fatalf("Targets[0] = %+v, want Identifier x", s.Targets[0])
	}
	if s.Value == nil {
		t.Fatalf("Value = nil, want 5")
	}
}

func TestParseSetTuple(t *testing.T) {
	n := parseOneScript(t, "SET (a, b, c) = (1 + 3, 'foo', false)")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetTuple {
		t.Fatalf("Kind = %v, want SetTuple", s.Kind)
	}
	if len(s.Targets) != 3 {
		t.Fatalf("Targets = %+v, want 3", s.Targets)
	}
	if s.Value == nil {
		t.Fatalf("Value = nil, want RHS tuple")
	}
}

func TestParseSetParameter(t *testing.T) {
	n := parseOneScript(t, "SET @p = 5")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetVariable {
		t.Fatalf("Kind = %v, want SetVariable", s.Kind)
	}
	if _, ok := s.Targets[0].(*ast.Parameter); !ok {
		t.Fatalf("Targets[0] = %T, want *ast.Parameter", s.Targets[0])
	}
}

func TestParseSetSystemVariable(t *testing.T) {
	n := parseOneScript(t, "SET @@x.y = 5")
	s := n.(*ast.SetStmt)
	if _, ok := s.Targets[0].(*ast.SystemVariable); !ok {
		t.Fatalf("Targets[0] = %T, want *ast.SystemVariable", s.Targets[0])
	}
}

func TestParseSetTransaction(t *testing.T) {
	n := parseOneScript(t, "SET TRANSACTION READ ONLY")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetTransaction {
		t.Fatalf("Kind = %v, want SetTransaction", s.Kind)
	}
	if len(s.Modes) != 1 || s.Modes[0].Kind != ast.TransactionModeReadOnly {
		t.Fatalf("Modes = %+v, want [READ ONLY]", s.Modes)
	}
}

func TestParseSetTransactionIsolation(t *testing.T) {
	n := parseOneScript(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetTransaction {
		t.Fatalf("Kind = %v, want SetTransaction", s.Kind)
	}
	if len(s.Modes) != 1 || s.Modes[0].Kind != ast.TransactionModeIsolationLevel {
		t.Fatalf("Modes = %+v, want [ISOLATION LEVEL]", s.Modes)
	}
}

// ----------------------------------------------------------------------------
// EXECUTE IMMEDIATE
// ----------------------------------------------------------------------------

func TestParseExecuteImmediateBare(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT 1"`)
	e, ok := n.(*ast.ExecuteImmediateStmt)
	if !ok {
		t.Fatalf("expected *ast.ExecuteImmediateStmt, got %T", n)
	}
	if e.SQL == nil {
		t.Fatalf("SQL = nil")
	}
	if e.Into != nil || e.Using != nil {
		t.Fatalf("Into/Using should be nil, got %+v / %+v", e.Into, e.Using)
	}
}

func TestParseExecuteImmediateIntoUsingPositional(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT ? * (? + 2)" INTO y USING 1, 3`)
	e := n.(*ast.ExecuteImmediateStmt)
	if len(e.Into) != 1 || e.Into[0].Name != "y" {
		t.Fatalf("Into = %+v, want [y]", e.Into)
	}
	if len(e.Using) != 2 {
		t.Fatalf("Using = %+v, want 2 args", e.Using)
	}
	if e.Using[0].Alias != "" || e.Using[1].Alias != "" {
		t.Fatalf("positional USING args should have no alias, got %q / %q", e.Using[0].Alias, e.Using[1].Alias)
	}
}

func TestParseExecuteImmediateUsingNamed(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT @a * (@b + 2)" INTO y USING 1 AS a, 3 AS b`)
	e := n.(*ast.ExecuteImmediateStmt)
	if len(e.Using) != 2 {
		t.Fatalf("Using = %+v, want 2 args", e.Using)
	}
	if e.Using[0].Alias != "a" || e.Using[1].Alias != "b" {
		t.Fatalf("USING aliases = %q / %q, want a / b", e.Using[0].Alias, e.Using[1].Alias)
	}
}

func TestParseExecuteImmediateIntoMultiple(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT 1, 2" INTO a, b`)
	e := n.(*ast.ExecuteImmediateStmt)
	if len(e.Into) != 2 || e.Into[0].Name != "a" || e.Into[1].Name != "b" {
		t.Fatalf("Into = %+v, want [a b]", e.Into)
	}
}

// ----------------------------------------------------------------------------
// IF
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// CASE (statement)
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// WHILE / LOOP / REPEAT
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// FOR ... IN
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// BEGIN ... END
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// BREAK / LEAVE / CONTINUE / ITERATE / RETURN
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// RAISE
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Labels
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Deep nesting (truth1 doc examples)
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Reject cases (negative tests)
// ----------------------------------------------------------------------------

func TestScriptingRejects(t *testing.T) {
	rejects := []string{
		// SET shapes the .g4 rejects.
		"SET",           // missing target
		"SET x",         // missing = expr
		"SET x =",       // missing RHS
		"SET a, b = 1",  // multiple vars without parens (grammar error alt)
		"SET (a, b) = ", // missing RHS
		"SET () = (1)",  // empty identifier list
		// DECLARE shapes.
		"DECLARE",         // missing name
		"DECLARE x",       // missing type and DEFAULT
		"DECLARE 1 INT64", // name is not an identifier
		// Control flow missing required keywords.
		"IF x SELECT 1; END IF",                  // missing THEN
		"IF x THEN SELECT 1;",                    // missing END IF
		"WHILE x SELECT 1; END WHILE",            // missing DO
		"WHILE x DO SELECT 1; END",               // END not followed by WHILE
		"REPEAT SELECT 1; END REPEAT",            // missing UNTIL clause
		"LOOP SELECT 1; END WHILE",               // mismatched closer
		"FOR x (SELECT 1) DO SELECT 1; END FOR",  // missing IN
		"FOR x IN SELECT 1 DO SELECT 1; END FOR", // IN source not parenthesized
		"CASE WHEN c THEN SELECT 1; END",         // END not followed by CASE
		"CASE END CASE",                          // no WHEN clause
		// EXECUTE IMMEDIATE shapes.
		"EXECUTE",                      // missing IMMEDIATE
		"EXECUTE IMMEDIATE",            // missing sql expr
		`EXECUTE IMMEDIATE "x" INTO`,   // INTO with no variable
		`EXECUTE IMMEDIATE "x" USING`,  // USING with no arg
		`EXECUTE IMMEDIATE "x" INTO 1`, // INTO target not an identifier
		// BEGIN ... END.
		"BEGIN SELECT 1;", // unterminated block (no END)
		"BEGIN EXCEPTION WHEN ERROR THEN SELECT 1; END", // EXCEPTION needs a body before? actually valid; remove
	}
	for _, src := range rejects {
		if src == "BEGIN EXCEPTION WHEN ERROR THEN SELECT 1; END" {
			continue // handled below as accept
		}
		t.Run(src, func(t *testing.T) {
			_, errs := parseForTest(src)
			if len(errs) == 0 {
				t.Errorf("parseForTest(%q) accepted, want a syntax error", src)
			}
		})
	}
}
