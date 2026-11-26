package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/ast"
)

func TestEngineSplitDropsEmptySegments(t *testing.T) {
	script := "SELECT 1;;\n-- c\n;SELECT 2"
	segments := Engine.Split(script)
	if len(segments) != 2 || segments[0].Text != "SELECT 1;" || segments[1].Text != "SELECT 2" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if s.Empty || script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v", s)
		}
	}
}

func TestEngineParse(t *testing.T) {
	parsed, err := Engine.Parse("CREATE TABLE t (a INT);\nSELECT a FROM t;")
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\nSELECT a FROM t;" {
		t.Fatalf("Parse = %+v", parsed)
	}
	if _, ok := parsed[0].Node.(*ast.CreateTableStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(";; -- c\n"); err != nil || parsed == nil || len(parsed) != 0 {
		t.Fatalf("Parse of no statement = %#v, %v", parsed, err)
	}
}

func TestEngineParseErrorShiftsTheColumnOnTheSegmentsFirstLine(t *testing.T) {
	_, err := Engine.Parse("SELECT 1; SELEC 5;")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	want := analysis.Position{Line: 1, Column: 11}
	if syntaxErr.Position != want || syntaxErr.Message != "Syntax error at line 1:11 \n"+syntaxErr.RawMessage || syntaxErr.RawMessage == "" {
		t.Fatalf("err = %+v", syntaxErr)
	}
	_, err = Engine.Parse("SELECT 1;\n  SELEC 5;")
	if !errors.As(err, &syntaxErr) || syntaxErr.Position != (analysis.Position{Line: 2, Column: 3}) {
		t.Fatalf("err = %+v", err)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT FROM t")
	if len(diagnostics) == 0 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateExplainAndLimit(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SET x = 1; SELECT 1"); err != nil || !readOnly || returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	if _, _, err := Engine.ValidateQuery("SELEC 1"); err == nil || errors.As(err, new(*analysis.SyntaxError)) {
		t.Fatalf("ValidateQuery error = %T %v", err, err)
	}
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN USING JSON SELECT 1",
		analysis.ExplainFormatXML:     "EXPLAIN SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if got := Engine.ResultLimit("SELECT * FROM t LIMIT 15;", 10, ""); got != "SELECT * FROM t LIMIT 10\n;" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT * FROM t OFFSET 2;", 10, ""); got != "SELECT * FROM (SELECT * FROM t OFFSET 2) LIMIT 10" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}
	if ranges := Engine.StatementRanges("SELECT 1;\nSELECT 2;"); len(ranges) != 2 || ranges[1].Start.Line != 1 {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
}
