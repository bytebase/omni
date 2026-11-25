package diagnostics

import "testing"

// TestAnalyze_ColumnsCountCodePoints: a diagnostic after a multi-byte
// character on the same line reports its column in code points, the unit of
// Bytebase's Position, while Offset stays the byte offset.
func TestAnalyze_ColumnsCountCodePoints(t *testing.T) {
	sql := "SELECT 'é' x y"
	diags := Analyze(sql)
	if len(diags) == 0 {
		t.Fatal("Analyze returned no diagnostic for trailing junk")
	}
	d := diags[0]
	// 'y' is the 14th code point (byte offset 14, since é is two bytes).
	if d.Range.Start.Offset != 14 {
		t.Fatalf("Offset = %d, want 14", d.Range.Start.Offset)
	}
	if d.Range.Start.Line != 1 || d.Range.Start.Column != 14 {
		t.Errorf("start = (%d, %d), want (1, 14) in code points", d.Range.Start.Line, d.Range.Start.Column)
	}
}
