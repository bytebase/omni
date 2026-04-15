package diagnostics

import (
	"strings"
	"testing"
)

// -----------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------

// assertNoDiags fails the test if diags is non-nil or non-empty.
func assertNoDiags(t *testing.T, label string, diags []Diagnostic) {
	t.Helper()
	if len(diags) != 0 {
		t.Errorf("%s: expected zero diagnostics, got %d:", label, len(diags))
		for i, d := range diags {
			t.Errorf("  [%d] %s", i, d.Message)
		}
	}
}

// -----------------------------------------------------------------------
// Zero-diagnostic cases
// -----------------------------------------------------------------------

func TestAnalyze_EmptyInput(t *testing.T) {
	assertNoDiags(t, "empty", Analyze(""))
}

func TestAnalyze_WhitespaceOnly(t *testing.T) {
	assertNoDiags(t, "spaces", Analyze("   "))
	assertNoDiags(t, "newlines", Analyze("\n\n\n"))
	assertNoDiags(t, "mixed whitespace", Analyze("  \t\n  "))
}

func TestAnalyze_CommentOnly(t *testing.T) {
	assertNoDiags(t, "line comment", Analyze("-- just a comment"))
	assertNoDiags(t, "block comment", Analyze("/* just a comment */"))
}

func TestAnalyze_ValidSelect(t *testing.T) {
	assertNoDiags(t, "SELECT 1", Analyze("SELECT 1"))
	assertNoDiags(t, "SELECT with alias", Analyze("SELECT 1 AS n"))
	assertNoDiags(t, "SELECT columns", Analyze("SELECT a, b FROM t"))
	assertNoDiags(t, "SELECT with WHERE", Analyze("SELECT id FROM users WHERE active = TRUE"))
}

func TestAnalyze_ValidCTE(t *testing.T) {
	sql := `WITH cte AS (SELECT 1 AS x)
SELECT x FROM cte`
	assertNoDiags(t, "CTE", Analyze(sql))
}

// -----------------------------------------------------------------------
// Single-error cases — line/column accuracy
// -----------------------------------------------------------------------

func TestAnalyze_UnknownStatement(t *testing.T) {
	// "FOOBAR" is not a recognised statement keyword. Error at byte 0.
	diags := Analyze("FOOBAR BAZ")
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d: %+v", len(diags), diags)
	}
	d := diags[0]

	if d.Severity != SeverityError {
		t.Errorf("severity = %v, want SeverityError", d.Severity)
	}
	if d.Source != "snowflake-parser" {
		t.Errorf("source = %q, want %q", d.Source, "snowflake-parser")
	}
	if d.Message == "" {
		t.Errorf("message must not be empty")
	}
	if d.Range.Start.Line != 1 {
		t.Errorf("start line = %d, want 1", d.Range.Start.Line)
	}
	if d.Range.Start.Column != 1 {
		t.Errorf("start column = %d, want 1", d.Range.Start.Column)
	}
	if d.Range.Start.Offset != 0 {
		t.Errorf("start offset = %d, want 0", d.Range.Start.Offset)
	}
}

func TestAnalyze_ErrorOnFirstLine_ColumnAccuracy(t *testing.T) {
	// "@@bogus@@" — bad token at byte 0, col 1 on line 1.
	diags := Analyze("@@bogus@@")
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic, got none")
	}
	d := diags[0]
	if d.Range.Start.Line != 1 {
		t.Errorf("line = %d, want 1", d.Range.Start.Line)
	}
	if d.Range.Start.Column != 1 {
		t.Errorf("column = %d, want 1", d.Range.Start.Column)
	}
}

func TestAnalyze_UnsupportedStatementInMiddle(t *testing.T) {
	// "SELECT 1; @@bogus@@" — @@ starts at byte 10 (after "SELECT 1; ").
	// Line is still 1, column 11.
	sql := "SELECT 1; @@bogus@@"
	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic, got none")
	}
	d := diags[0]
	if d.Range.Start.Line != 1 {
		t.Errorf("line = %d, want 1", d.Range.Start.Line)
	}
	// "SELECT 1; " is 10 bytes, so the bad token starts at column 11.
	if d.Range.Start.Column != 11 {
		t.Errorf("column = %d, want 11 (bad token starts after 'SELECT 1; ')", d.Range.Start.Column)
	}
	if d.Range.Start.Offset != 10 {
		t.Errorf("offset = %d, want 10", d.Range.Start.Offset)
	}
}

// -----------------------------------------------------------------------
// Multi-line inputs
// -----------------------------------------------------------------------

func TestAnalyze_MultiLine_ErrorOnLine3(t *testing.T) {
	// Error should be reported on line 3 (the bad token).
	sql := "SELECT 1;\nSELECT 2;\n@@bogus@@"
	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic")
	}
	d := diags[0]
	if d.Range.Start.Line != 3 {
		t.Errorf("line = %d, want 3", d.Range.Start.Line)
	}
}

func TestAnalyze_MultiLine_ErrorOnLine2(t *testing.T) {
	// "SELECT 1\n@@bad" — @@ introduces an error on line 2.
	// We use an unterminated string on line 2.
	sql := "SELECT 1;\n'unterminated"
	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic")
	}
	// The lex error for the unterminated string should be on line 2.
	found := false
	for _, d := range diags {
		if d.Range.Start.Line == 2 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a diagnostic on line 2; got: %+v", diags)
	}
}

func TestAnalyze_CRLF_LineEndings(t *testing.T) {
	// CRLF line endings — error should still report line 3.
	sql := "SELECT 1;\r\nSELECT 2;\r\n@@bogus@@"
	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic")
	}
	d := diags[0]
	if d.Range.Start.Line != 3 {
		t.Errorf("line = %d, want 3 (CRLF endings)", d.Range.Start.Line)
	}
}

// -----------------------------------------------------------------------
// Multiple errors
// -----------------------------------------------------------------------

func TestAnalyze_MultipleErrors(t *testing.T) {
	// Two malformed statements.
	sql := "@@bogus1@@;\n@@bogus2@@"
	diags := Analyze(sql)
	if len(diags) < 2 {
		t.Fatalf("expected at least 2 diagnostics, got %d", len(diags))
	}
}

func TestAnalyze_LexError_UnterminatedString(t *testing.T) {
	// The lexer emits a lex error for the unterminated string; the parser also
	// records a parse error because the resulting tokInvalid looks like an
	// unknown statement start. We just need at least one diagnostic, and at
	// least one should mention "unterminated".
	diags := Analyze("'unterminated string")
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic")
	}
	for _, d := range diags {
		if d.Severity != SeverityError {
			t.Errorf("severity = %v, want SeverityError", d.Severity)
		}
	}
	found := false
	for _, d := range diags {
		if strings.Contains(d.Message, "unterminated") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a diagnostic containing 'unterminated'; got: %+v", diags)
	}
}

// -----------------------------------------------------------------------
// Source field and message invariants
// -----------------------------------------------------------------------

func TestAnalyze_SourceField(t *testing.T) {
	diags := Analyze("INVALID_STATEMENT")
	if len(diags) == 0 {
		t.Fatal("expected diagnostics")
	}
	for i, d := range diags {
		if d.Source != "snowflake-parser" {
			t.Errorf("diags[%d].Source = %q, want \"snowflake-parser\"", i, d.Source)
		}
	}
}

func TestAnalyze_MessageNonEmpty(t *testing.T) {
	diags := Analyze("INVALID_STATEMENT")
	if len(diags) == 0 {
		t.Fatal("expected diagnostics")
	}
	for i, d := range diags {
		if d.Message == "" {
			t.Errorf("diags[%d].Message is empty", i)
		}
	}
}

// -----------------------------------------------------------------------
// Severity type
// -----------------------------------------------------------------------

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		s    Severity
		want string
	}{
		{SeverityError, "error"},
		{SeverityWarning, "warning"},
		{SeverityInfo, "info"},
		{Severity(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Severity(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}

// -----------------------------------------------------------------------
// Range sanity: Start.Offset <= End.Offset
// -----------------------------------------------------------------------

func TestAnalyze_RangeSanity(t *testing.T) {
	inputs := []string{
		"FOOBAR",
		"INSERT INTO t VALUES (1)",
		"'unterminated",
	}
	for _, sql := range inputs {
		diags := Analyze(sql)
		for i, d := range diags {
			if d.Range.Start.Offset > d.Range.End.Offset {
				t.Errorf("sql=%q diag[%d]: Start.Offset(%d) > End.Offset(%d)",
					sql, i, d.Range.Start.Offset, d.Range.End.Offset)
			}
			if d.Range.Start.Line > d.Range.End.Line {
				t.Errorf("sql=%q diag[%d]: Start.Line(%d) > End.Line(%d)",
					sql, i, d.Range.Start.Line, d.Range.End.Line)
			}
		}
	}
}

// -----------------------------------------------------------------------
// UTF-8 / multibyte characters
// -----------------------------------------------------------------------

func TestAnalyze_UTF8_ByteBasedColumn(t *testing.T) {
	// Place an unsupported statement on a line that starts with multi-byte
	// UTF-8 characters so we can verify the column is byte-based.
	//
	// Line 2 is: "αβ INSERT INTO t VALUES (1)"
	//   α = 2 bytes (0xCE 0xB1)
	//   β = 2 bytes (0xCE 0xB2)
	//   space = 1 byte
	//   INSERT starts at byte offset 5 within that line → column 6
	//
	// The Snowflake lexer does not support non-ASCII identifiers; α and β will
	// each cause an "invalid byte" lex error. We are primarily testing that the
	// LineTable produces the correct byte-based column for the error on line 2.
	sql := "SELECT 1;\nαβ INSERT INTO t VALUES (1)"

	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics")
	}

	// All diagnostics are on line 2 (the second line).
	for _, d := range diags {
		if d.Range.Start.Line != 2 {
			t.Errorf("expected all diagnostics on line 2; got line %d (msg=%q)",
				d.Range.Start.Line, d.Message)
		}
	}

	// The first diagnostic is the "invalid byte" for α at column 1 (offset 10
	// in the full input, which is the first byte of line 2).
	d0 := diags[0]
	if d0.Range.Start.Column != 1 {
		t.Errorf("diags[0] column = %d, want 1 (first byte of line 2)", d0.Range.Start.Column)
	}
	// Offset 10 = length of "SELECT 1;\n" (10 bytes).
	if d0.Range.Start.Offset != 10 {
		t.Errorf("diags[0] offset = %d, want 10", d0.Range.Start.Offset)
	}
}

// -----------------------------------------------------------------------
// Position consistency: Offset matches Line+Column
// -----------------------------------------------------------------------

func TestAnalyze_PositionConsistency(t *testing.T) {
	// For a multi-line input, verify that the Offset of the error on line N
	// is consistent with the line/column values.
	sql := "SELECT 1;\nSELECT 2;\n@@bogus@@"
	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("expected diagnostics")
	}
	d := diags[0]
	// Line 3, column 1 → offset should equal the byte position of "@@bogus@@" in sql.
	wantOffset := strings.Index(sql, "@@bogus@@")
	if d.Range.Start.Offset != wantOffset {
		t.Errorf("Start.Offset = %d, want %d", d.Range.Start.Offset, wantOffset)
	}
	if d.Range.Start.Line != 3 {
		t.Errorf("Start.Line = %d, want 3", d.Range.Start.Line)
	}
	if d.Range.Start.Column != 1 {
		t.Errorf("Start.Column = %d, want 1", d.Range.Start.Column)
	}
}
