package analysis

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineSplit(t *testing.T) {
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
}

func TestEngineValidateQueryPositionsASyntaxErrorInTheScript(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SET search_path TO public"); err != nil || !readOnly || returnsData {
		t.Fatalf("ValidateQuery(SET) = %v %v %v", readOnly, returnsData, err)
	}
	_, _, err := Engine.ValidateQuery("SELECT 1;\n  SELECT * FROM;")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 16}) || syntaxErr.Message != `Syntax error at line 2: syntax error at or near ";"` || syntaxErr.RawMessage != `syntax error at or near ";"` {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineDiagnoseAndRanges(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError || diagnostics[0].Message != "syntax error at end of input" {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
	if ranges := Engine.StatementRanges("SELECT 1;\n  SELECT 2;"); len(ranges) != 2 || ranges[1].Start != (analysis.LSPPosition{Line: 1, Character: 2}) {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
	if ranges := Engine.StatementRanges(";"); ranges != nil {
		t.Fatalf("StatementRanges of no statement = %+v", ranges)
	}
}

func TestEngineExplainAndLimit(t *testing.T) {
	for _, format := range []analysis.ExplainFormat{analysis.ExplainFormatDefault, analysis.ExplainFormatJSON, analysis.ExplainFormatXML} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != "EXPLAIN SELECT 1" {
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
}

func TestEngineExtractChangesReadTheScript(t *testing.T) {
	request := analysis.ChangeRequest{DefaultDatabase: "db", Script: "TRUNCATE t; INSERT INTO t VALUES (1), (2); SET search_path TO s1; DELETE FROM t;"}
	changes, err := Engine.ExtractChanges(request, nil)
	if err != nil {
		t.Fatal(err)
	}
	wantResources := []analysis.ChangedResource{
		{Database: "db", Schema: "public", Table: "t"},
		{Database: "db", Schema: "s1", Table: "t"},
		{Database: "db", Schema: "public", Table: "t", AffectsData: true},
	}
	if !slices.Equal(changes.Resources, wantResources) {
		t.Fatalf("resources = %+v", changes.Resources)
	}
	// The script's own extractor prefixes the search path, so SearchPath stays nil.
	wantDML := []analysis.DMLStatement{{Text: "SET LOCAL search_path TO \"s1\";\nDELETE FROM t;"}}
	if !reflect.DeepEqual(changes.DMLStatements, wantDML) || changes.DMLCount != 1 || changes.InsertCount != 2 {
		t.Fatalf("changes = %+v", changes)
	}
	if _, err := Engine.ExtractChanges(analysis.ChangeRequest{Script: "DELETE FROM"}, nil); err == nil {
		t.Fatal("want the script's syntax error")
	}
}
