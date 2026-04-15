package parser

import "testing"

func TestDiagnoseCleanInput(t *testing.T) {
	// Empty input should produce no diagnostics (no statements to parse).
	diags := Diagnose("")
	if len(diags) != 0 {
		t.Errorf("expected 0 diagnostics for empty input, got %d", len(diags))
	}
}

func TestDiagnoseLexErrors(t *testing.T) {
	// Unterminated string literal is a genuine lex error.
	diags := Diagnose("SELECT 'unterminated")
	found := false
	for _, d := range diags {
		if d.Msg == "unterminated string literal" {
			found = true
		}
	}
	if !found {
		t.Error("expected unterminated string diagnostic")
	}
}

func TestDiagnoseUnterminatedComment(t *testing.T) {
	diags := Diagnose("SELECT /* unterminated")
	found := false
	for _, d := range diags {
		if d.Msg == "unterminated block comment" {
			found = true
		}
	}
	if !found {
		t.Error("expected unterminated comment diagnostic")
	}
}

func TestDiagnoseMultiStatement(t *testing.T) {
	// INSERT and TRUNCATE are still unsupported; each produces a diagnostic.
	diags := Diagnose("INSERT INTO t; TRUNCATE TABLE t")
	if len(diags) < 2 {
		t.Errorf("expected at least 2 diagnostics, got %d", len(diags))
	}
}

func TestDiagnosticPositions(t *testing.T) {
	// The unterminated string starts at byte 7 (after "SELECT ").
	diags := Diagnose("SELECT 'bad")
	if len(diags) == 0 {
		t.Fatal("expected diagnostics")
	}
	for _, d := range diags {
		if d.Msg == "unterminated string literal" {
			if d.Loc.Start != 7 {
				t.Errorf("Loc.Start = %d, want 7", d.Loc.Start)
			}
			return
		}
	}
	t.Error("unterminated string literal diagnostic not found")
}

func TestDiagnoseUnterminatedBacktickIdentifier(t *testing.T) {
	// Unterminated backtick-quoted identifier is also a genuine lex error.
	diags := Diagnose("SELECT `col")
	found := false
	for _, d := range diags {
		if d.Msg == "unterminated backtick-quoted identifier" {
			found = true
		}
	}
	if !found {
		t.Error("expected unterminated backtick-quoted identifier diagnostic")
	}
}

func TestDiagnoseReturnsDiagnosticSlice(t *testing.T) {
	// Diagnose on a supported-but-stubbed keyword returns a non-nil slice.
	diags := Diagnose("SELECT 1")
	if diags == nil {
		t.Error("expected non-nil diagnostics slice")
	}
}
