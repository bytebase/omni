package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/pg/ast"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT 1;\n-- lead\nSELECT 2;\n;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || segments[1].Text != "\n-- lead\nSELECT 2;" || !segments[2].Empty {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
		}
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
	inserts, err := Engine.Parse("INSERT INTO t VALUES (1)")
	if err != nil {
		t.Fatal(err)
	}
	if types := Engine.StatementTypes(inserts[0]); len(types) != 1 || types[0] != analysis.Insert {
		t.Fatalf("StatementTypes = %v", types)
	}
	if types := Engine.StatementTypes(analysis.ParsedStatement{}); types != nil {
		t.Fatalf("StatementTypes of no node = %v", types)
	}
}

func TestEngineParseErrorIsPositionedOnTheScriptsLines(t *testing.T) {
	_, err := Engine.Parse("SELECT 1;\nSELECT * FROM;")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position.Line != 2 || syntaxErr.Message != `Syntax error at line 2: syntax error at or near ";"` || syntaxErr.RawMessage != `syntax error at or near ";"` {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError {
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
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN (FORMAT JSON) SELECT 1",
		analysis.ExplainFormatXML:     "EXPLAIN (FORMAT XML) SELECT 1",
		analysis.ExplainFormatYAML:    "EXPLAIN (FORMAT YAML) SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if _, err := Engine.Explain("EXPLAIN ANALYZE SELECT 1", analysis.ExplainFormatDefault); !errors.Is(err, analysis.ErrExplainAnalyzeExecutes) {
		t.Fatalf("Explain(EXPLAIN ANALYZE) err = %v", err)
	}
	if got := Engine.ResultLimit("SELECT 1", 10, ""); got != "SELECT 1 LIMIT 10" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT 1; SELECT 2;", 10, ""); got != "WITH result AS (\nSELECT 1; SELECT 2\n) SELECT * FROM result LIMIT 10;" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}
	if ranges := Engine.StatementRanges("SELECT 1;\nSELECT 2"); len(ranges) != 2 || ranges[1].Start.Line != 1 {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
}
