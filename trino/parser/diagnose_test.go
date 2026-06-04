package parser

import (
	"strings"
	"testing"
)

// TestDiagnose_CleanInputNoStubNoise verifies the wiring shape: Diagnose runs
// Parse and surfaces ParseErrors as Diagnostics. While a statement body is still
// stubbed (no DAG node has implemented it), even valid SQL of that form yields a
// "not yet supported" diagnostic — those vanish as later nodes implement real
// parsing. INSERT (parser-dml's job) is still stubbed; SELECT is now implemented
// by parser-select, so this test uses a still-stubbed form.
func TestDiagnose_StubbedStatementProducesDiagnostic(t *testing.T) {
	diags := Diagnose("INSERT INTO t VALUES (1)")
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
// segments. A valid SELECT (now implemented by parser-select) yields none, while
// the two still-stubbed DML statements (INSERT, DELETE) each yield one.
func TestDiagnose_MultipleStatements(t *testing.T) {
	diags := Diagnose("SELECT 1; INSERT INTO t VALUES (1); DELETE FROM t")
	if len(diags) != 2 {
		t.Fatalf("got %d diagnostics, want 2 (INSERT + DELETE stubs; SELECT now parses): %+v", len(diags), diags)
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
