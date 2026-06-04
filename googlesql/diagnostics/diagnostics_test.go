package diagnostics

import (
	"strings"
	"testing"
)

func assertNoDiags(t *testing.T, label string, diags []Diagnostic) {
	t.Helper()
	if len(diags) != 0 {
		t.Errorf("%s: expected zero diagnostics, got %d:", label, len(diags))
		for i, d := range diags {
			t.Errorf("  [%d] %s", i, d.Message)
		}
	}
}

// --- Zero-diagnostic cases (no statement => nothing to stub or reject) ---

func TestAnalyze_EmptyInput(t *testing.T) {
	assertNoDiags(t, "empty", Analyze(""))
}

func TestAnalyze_WhitespaceOnly(t *testing.T) {
	assertNoDiags(t, "spaces", Analyze("   "))
	assertNoDiags(t, "newlines", Analyze("\n\n\n"))
	assertNoDiags(t, "mixed whitespace", Analyze("  \t\n  "))
}

func TestAnalyze_CommentOnly(t *testing.T) {
	assertNoDiags(t, "dash comment", Analyze("-- just a comment"))
	assertNoDiags(t, "block comment", Analyze("/* just a comment */"))
	assertNoDiags(t, "pound comment", Analyze("# pound comment"))
}

func TestAnalyze_LoneSemicolons(t *testing.T) {
	assertNoDiags(t, "semicolons", Analyze(";;;"))
}

// --- Severity / Source / shape ---

func TestSeverity_String(t *testing.T) {
	cases := map[Severity]string{
		SeverityError:   "error",
		SeverityWarning: "warning",
		SeverityInfo:    "info",
		Severity(99):    "unknown",
	}
	for sev, want := range cases {
		if got := sev.String(); got != want {
			t.Errorf("Severity(%d).String() = %q, want %q", int(sev), got, want)
		}
	}
}

func TestAnalyze_StubbedStatementShape(t *testing.T) {
	// A statement whose body is still stubbed (INSERT — the DML node has not
	// landed) produces a "not yet supported" diagnostic. Assert the diagnostic
	// shape. (The query family — SELECT/WITH — is now implemented by
	// parser-select and produces zero diagnostics for valid input; see
	// TestAnalyze_ValidQueryNoDiagnostics.)
	diags := Analyze("INSERT INTO t VALUES (1)")
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d: %+v", len(diags), diags)
	}
	d := diags[0]
	if d.Severity != SeverityError {
		t.Errorf("severity = %v, want SeverityError", d.Severity)
	}
	if d.Source != source {
		t.Errorf("source = %q, want %q", d.Source, source)
	}
	if !strings.Contains(d.Message, "not yet supported") {
		t.Errorf("message = %q, want it to mention 'not yet supported'", d.Message)
	}
	// INSERT starts at byte 0 => line 1, col 1.
	if d.Range.Start.Line != 1 || d.Range.Start.Column != 1 {
		t.Errorf("start = (%d, %d), want (1, 1)", d.Range.Start.Line, d.Range.Start.Column)
	}
}

// TestAnalyze_ValidQueryNoDiagnostics confirms that, with the parser-select node
// landed, a syntactically valid query produces zero diagnostics (no more "not
// yet supported" for the query family). This is the bytebase Diagnose contract:
// a valid BigQuery/Spanner query must not draw a false syntax diagnostic.
func TestAnalyze_ValidQueryNoDiagnostics(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1",
		"SELECT a, b FROM t WHERE a > 1 ORDER BY b LIMIT 10",
		"WITH c AS (SELECT 1 AS n) SELECT n FROM c",
		"SELECT * FROM a JOIN b USING (id)",
		"SELECT * FROM t1 UNION ALL SELECT * FROM t2",
	} {
		if diags := Analyze(sql); len(diags) != 0 {
			t.Errorf("Analyze(%q): got %d diagnostics, want 0: %+v", sql, len(diags), diags)
		}
	}
}

// --- Lex-error mapping with accurate line/column ---

func TestAnalyze_UnterminatedString_Position(t *testing.T) {
	// The unterminated string starts at byte 7 on line 1.
	diags := Analyze("SELECT 'oops")
	var lexDiag *Diagnostic
	for i := range diags {
		if strings.Contains(diags[i].Message, "unterminated string") {
			lexDiag = &diags[i]
		}
	}
	if lexDiag == nil {
		t.Fatalf("no unterminated-string diagnostic in %+v", diags)
	}
	if lexDiag.Range.Start.Offset != 7 {
		t.Errorf("start offset = %d, want 7", lexDiag.Range.Start.Offset)
	}
	if lexDiag.Range.Start.Line != 1 || lexDiag.Range.Start.Column != 8 {
		t.Errorf("start = (%d, %d), want (1, 8)", lexDiag.Range.Start.Line, lexDiag.Range.Start.Column)
	}
}

func TestAnalyze_ErrorOnSecondLine_Position(t *testing.T) {
	// First statement is stubbed (error at line 1); the unterminated string is
	// on line 2. Verify the multi-line offset → (line, col) mapping.
	sql := "SELECT 1;\nSELECT 'oops"
	diags := Analyze(sql)
	var lexDiag *Diagnostic
	for i := range diags {
		if strings.Contains(diags[i].Message, "unterminated string") {
			lexDiag = &diags[i]
		}
	}
	if lexDiag == nil {
		t.Fatalf("no unterminated-string diagnostic in %+v", diags)
	}
	if lexDiag.Range.Start.Line != 2 {
		t.Errorf("start line = %d, want 2", lexDiag.Range.Start.Line)
	}
	// "SELECT '" — the quote is the 8th byte on line 2 => column 8.
	if lexDiag.Range.Start.Column != 8 {
		t.Errorf("start column = %d, want 8", lexDiag.Range.Start.Column)
	}
}

func TestAnalyze_UnknownStatement(t *testing.T) {
	diags := Analyze("FOOBAR BAZ")
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d: %+v", len(diags), diags)
	}
	if !strings.Contains(diags[0].Message, "unknown or unsupported statement") {
		t.Errorf("message = %q, want 'unknown or unsupported statement'", diags[0].Message)
	}
}

func TestAnalyze_RangeIsValid(t *testing.T) {
	// End must be >= Start (never a garbage range) for every diagnostic.
	for _, sql := range []string{"SELECT 'oops", "FOOBAR", "SELECT 1", "/* a"} {
		for _, d := range Analyze(sql) {
			if d.Range.End.Offset < d.Range.Start.Offset {
				t.Errorf("Analyze(%q): End.Offset %d < Start.Offset %d", sql, d.Range.End.Offset, d.Range.Start.Offset)
			}
			if d.Range.Start.Line < 1 || d.Range.Start.Column < 1 {
				t.Errorf("Analyze(%q): start (%d,%d) must be 1-based", sql, d.Range.Start.Line, d.Range.Start.Column)
			}
		}
	}
}
