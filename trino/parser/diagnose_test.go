package parser

import (
	"strings"
	"testing"
)

// TestDiagnose_CleanInputNoStubNoise verifies the wiring shape: Diagnose runs
// Parse and surfaces ParseErrors as Diagnostics. The CREATE/DROP/ALTER dispatch
// recognizes the leading keyword but routes an UNRECOGNIZED object keyword to
// the `unsupported` stub (e.g. `CREATE INDEX` / `ALTER INDEX`, which Trino has
// no statement form for and also rejects). That stub path is what this test
// exercises. (Every real Trino DDL/DML/query/admin form is now implemented by
// the parser-* nodes, so a valid statement no longer hits the stub.)
func TestDiagnose_StubbedStatementProducesDiagnostic(t *testing.T) {
	diags := Diagnose("ALTER INDEX i RENAME TO j")
	if len(diags) != 1 {
		t.Fatalf("got %d diagnostics, want 1: %+v", len(diags), diags)
	}
	if !strings.Contains(diags[0].Msg, "not yet supported") {
		t.Errorf("Msg = %q, want 'not yet supported'", diags[0].Msg)
	}
}

// TestDiagnose_LexError verifies a lexer error (unterminated string) surfaces
// as a diagnostic with a position.
func TestDiagnose_LexError(t *testing.T) {
	diags := Diagnose("SELECT 'oops")
	if len(diags) == 0 {
		t.Fatal("want at least one diagnostic for unterminated string, got none")
	}
	found := false
	for _, d := range diags {
		if strings.Contains(d.Msg, "unterminated string") {
			found = true
			if !d.Loc.IsValid() {
				t.Errorf("diagnostic %q has invalid Loc %+v", d.Msg, d.Loc)
			}
		}
	}
	if !found {
		t.Errorf("diagnostics = %+v, want one mentioning 'unterminated string'", diags)
	}
}

// TestDiagnose_EmptyInput verifies clean empty input yields no diagnostics.
func TestDiagnose_EmptyInput(t *testing.T) {
	if diags := Diagnose("   \n-- c\n"); len(diags) != 0 {
		t.Errorf("got %+v, want no diagnostics", diags)
	}
}

// TestDiagnose_MultipleStatements verifies diagnostics are collected across all
// segments. A valid SELECT, a valid INSERT, a valid ALTER TABLE, and a valid
// COMMENT (all now implemented by parser-select / parser-dml / parser-ddl)
// yield none, while the two unrecognized-object stubs (CREATE INDEX, ALTER
// INDEX — no Trino statement form) each yield one.
func TestDiagnose_MultipleStatements(t *testing.T) {
	diags := Diagnose("SELECT 1; INSERT INTO t VALUES (1); CREATE INDEX i ON t (a); ALTER INDEX i RENAME TO j")
	if len(diags) != 2 {
		t.Fatalf("got %d diagnostics, want 2 (CREATE INDEX + ALTER INDEX stubs; SELECT and INSERT parse): %+v", len(diags), diags)
	}
}

// TestDiagnose_UnknownStatement verifies an unrecognized leading token yields
// an "unknown or unsupported statement" diagnostic (distinct from the
// "not yet supported" stub message).
func TestDiagnose_UnknownStatement(t *testing.T) {
	diags := Diagnose("FROBNICATE foo")
	if len(diags) != 1 {
		t.Fatalf("got %d diagnostics, want 1: %+v", len(diags), diags)
	}
	if !strings.Contains(diags[0].Msg, "unknown or unsupported statement") {
		t.Errorf("Msg = %q, want 'unknown or unsupported statement'", diags[0].Msg)
	}
}
