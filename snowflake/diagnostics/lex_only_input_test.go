package diagnostics

import "testing"

func TestAnalyze_ReportsUnterminatedCommentAlone(t *testing.T) {
	diags := Analyze("/* unterminated")
	if len(diags) == 0 {
		t.Fatal("Analyze returned no diagnostic for an input that is only an unterminated comment")
	}
	if diags[0].Range.Start.Offset != 0 {
		t.Errorf("Offset = %d, want 0", diags[0].Range.Start.Offset)
	}
}
