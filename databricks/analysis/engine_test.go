package analysis

import (
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineSplitKeepsEmptySegments(t *testing.T) {
	script := "SELECT 1;\n;"
	segments := Engine.Split(script)
	if len(segments) != 2 || segments[0].Text != "SELECT 1;" || !segments[1].Empty {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
		}
	}
}

func TestEngineExplainIgnoresFormat(t *testing.T) {
	for _, format := range []analysis.ExplainFormat{analysis.ExplainFormatDefault, analysis.ExplainFormatJSON} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != "EXPLAIN FORMATTED SELECT 1" {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
}
