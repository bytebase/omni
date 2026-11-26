package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/trino/parser"
)

func TestEngineSplitAndParseKeepEmptySegments(t *testing.T) {
	script := "SELECT 1;\n;\nSELECT 2"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || !segments[1].Empty || segments[2].Text != "\nSELECT 2" {
		t.Fatalf("Split = %+v", segments)
	}
	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 3 || parsed[1].Node != nil || !parsed[1].Empty {
		t.Fatalf("Parse = %+v", parsed)
	}
	if _, ok := parsed[2].Node.(*parser.QueryStmt); !ok {
		t.Fatalf("node = %T", parsed[2].Node)
	}
	if parsed, err := Engine.Parse(""); err != nil || parsed == nil || len(parsed) != 0 {
		t.Fatalf("Parse of no statement = %#v, %v", parsed, err)
	}
}

func TestEngineSyntaxErrors(t *testing.T) {
	var syntaxErr *analysis.SyntaxError
	_, err := Engine.Parse("SELECT 1; SELEC 2")
	if !errors.As(err, &syntaxErr) || syntaxErr.Position != (analysis.Position{Line: 1, Column: 11}) || syntaxErr.Message != syntaxErr.RawMessage {
		t.Fatalf("Parse err = %+v", err)
	}
	_, _, err = Engine.ValidateQuery("SELECT 1;\n SELEC 2")
	if !errors.As(err, &syntaxErr) || syntaxErr.Position != (analysis.Position{Line: 2, Column: 2}) {
		t.Fatalf("ValidateQuery err = %+v", err)
	}
	_, err = Engine.Explain("SELEC 1", analysis.ExplainFormatDefault)
	if !errors.As(err, &syntaxErr) || syntaxErr.Position != (analysis.Position{Line: 1, Column: 1}) {
		t.Fatalf("Explain err = %+v", err)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT FROM")
	if len(diagnostics) == 0 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateAndExplain(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1; USE c.s"); err != nil || !readOnly || returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN (FORMAT JSON) SELECT 1",
		analysis.ExplainFormatYAML:    "EXPLAIN SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if _, err := Engine.Explain("EXPLAIN ANALYZE SELECT 1", analysis.ExplainFormatDefault); !errors.Is(err, analysis.ErrExplainAnalyzeExecutes) {
		t.Fatalf("Explain(EXPLAIN ANALYZE) err = %v", err)
	}
}
