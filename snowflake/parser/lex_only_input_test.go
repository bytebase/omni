package parser

import "testing"

// TestParse_ReportsLexErrorInDroppedSegment: an input that lexes to nothing
// (an unterminated block comment) yields no segment, so the per-segment
// promotion never runs; strict Parse must still report the lex error.
func TestParse_ReportsLexErrorInDroppedSegment(t *testing.T) {
	for _, sql := range []string{"/* unterminated", "SELECT 1; /* unterminated"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil, want the unterminated-comment error", sql)
		}
	}
	if r := ParseBestEffort("/* unterminated"); len(r.Errors) != 0 {
		t.Errorf("ParseBestEffort keeps its per-segment reporting: errors=%v", r.Errors)
	}
}
