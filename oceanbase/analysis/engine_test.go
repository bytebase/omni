package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/oceanbase/ast"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT 1;\n-- lead\nSELECT 2;\n;/*! */"
	segments := Engine.Split(script)
	if len(segments) != 4 || segments[0].Text != "SELECT 1;" || segments[1].Text != "\n-- lead\nSELECT 2;" || !segments[2].Empty || segments[3].Text != "/*! */" || segments[3].Empty {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
		}
	}
	if segments := Engine.Split(""); segments != nil {
		t.Fatalf("Split(\"\") = %#v", segments)
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\n-- lead\nSELECT 2;" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(";"); err != nil || parsed != nil {
		t.Fatalf("Parse(;) = %#v, %v", parsed, err)
	}
}

func TestEngineParseErrorShiftsTheFirstLinesColumn(t *testing.T) {
	_, err := Engine.Parse("SELECT 1; SELEC 2;")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 1, Column: 11}) || syntaxErr.Message != "Syntax error at line 1:11: syntax error at or near \"SELEC\"\nrelated text:  SELEC 2;" || syntaxErr.RawMessage != `syntax error at or near "SELEC"` {
		t.Fatalf("err = %+v", syntaxErr)
	}
	_, err = Engine.Parse("SELECT 1;\n  SELEC 2;")
	if !errors.As(err, &syntaxErr) || syntaxErr.Position != (analysis.Position{Line: 2, Column: 3}) {
		t.Fatalf("err = %v", err)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 2, Column: 14}) || diagnostics[0].Severity != analysis.SeverityError || diagnostics[0].Message != "syntax error at end of input" {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateExplainAndLimit(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1; SELECT 2"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	if readOnly, _, err := Engine.ValidateQuery("DELETE FROM t"); err != nil || readOnly {
		t.Fatalf("ValidateQuery(DELETE) = %v %v", readOnly, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery("SELECT 1; SELEC 2"); !errors.As(err, &syntaxErr) || syntaxErr.Position != (analysis.Position{Line: 1, Column: 11}) {
		t.Fatalf("ValidateQuery error = %v", err)
	}

	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN FORMAT=JSON SELECT 1",
		analysis.ExplainFormatXML:     "EXPLAIN SELECT 1",
		analysis.ExplainFormatYAML:    "EXPLAIN SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if _, err := Engine.Explain("ANALYZE TABLE t", analysis.ExplainFormatDefault); !errors.Is(err, analysis.ErrNotExplainable) {
		t.Fatalf("Explain(ANALYZE TABLE) err = %v", err)
	}
	if _, err := Engine.Explain("SELEC 1", analysis.ExplainFormatDefault); !errors.As(err, &syntaxErr) {
		t.Fatalf("Explain(SELEC 1) err = %v", err)
	}

	if got := Engine.ResultLimit("SELECT 1", 10, ""); got != "SELECT 1 LIMIT 10" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT 1; SELECT 2;", 10, ""); got != "SELECT * FROM (SELECT 1; SELECT 2) result LIMIT 10;" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}
}
