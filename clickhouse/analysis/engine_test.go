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
	if segments := Engine.Split(""); segments != nil {
		t.Fatalf("Split of nothing = %+v", segments)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT * FROM t WHERE )")
	if len(diagnostics) != 1 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError || diagnostics[0].Message == "" {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateExplainAndLimit(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery("SELECT * FROM t WHERE )"); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("ValidateQuery returns the parser's own error, got %v", err)
	}
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN json = 1, indexes = 1 SELECT 1",
		analysis.ExplainFormatYAML:    "EXPLAIN SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if got := Engine.ResultLimit("SELECT * FROM t;", 10, ""); got != "WITH result AS (SELECT * FROM t) SELECT * FROM result LIMIT 10;" {
		t.Fatalf("ResultLimit = %q", got)
	}
}
