package analysis

import (
	"errors"
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

func TestEngineValidateAndExplain(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1; SELECT 2"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery("SELECT * FROM"); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("ValidateQuery returns the parser's own error, got %v", err)
	}
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN FORMATTED SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
}
