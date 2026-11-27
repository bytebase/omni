package analysis

import (
	"context"
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineSplitKeepsSemicolonsAndEmptySegments(t *testing.T) {
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

func TestEngineRangesAndDiagnose(t *testing.T) {
	if ranges := Engine.StatementRanges(""); ranges != nil {
		t.Fatalf("StatementRanges of no statement = %+v", ranges)
	}
	if ranges := Engine.StatementRanges("SELECT 1;\nSELECT 2"); len(ranges) != 2 || ranges[1].Start.Line != 1 {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineSyntaxErrorIsPositionedInTheScript(t *testing.T) {
	_, _, err := Engine.ValidateQuery("SELECT 1;\nSELECT * FROM")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position.Line != 2 || syntaxErr.Message == "" || syntaxErr.Message != syntaxErr.RawMessage {
		t.Fatalf("err = %+v", syntaxErr)
	}
	if _, err := Engine.Explain("SELECT 1; SELECT 2", analysis.ExplainFormatDefault); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("Explain of two statements err = %v", err)
	}
}

func TestEngineExplainIgnoresFormatAndLimit(t *testing.T) {
	for _, format := range []analysis.ExplainFormat{analysis.ExplainFormatDefault, analysis.ExplainFormatJSON} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != "EXPLAIN SELECT 1" {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if got, err := Engine.Explain("TRACE ALL SELECT 1", analysis.ExplainFormatJSON); err != nil || got != "EXPLAIN SELECT 1" {
		t.Fatalf("Explain(TRACE) = %q, %v", got, err)
	}
	if got := Engine.ResultLimit("SELECT * FROM t", 10, ""); got != "SELECT * FROM t LIMIT 10" {
		t.Fatalf("ResultLimit = %q", got)
	}
}

func TestEngineQuerySpansReadTablesInTheDefaultDatabaseWithoutTheCatalog(t *testing.T) {
	unread := analysis.Catalog{GetDatabase: func(context.Context, string) (*analysis.Database, error) {
		return nil, errors.New("the catalog is read")
	}}
	statements := []analysis.QuerySpanStatement{{Text: "SELECT id FROM t;"}, {Text: "\nSELECT * FROM other.t1"}}
	spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", DefaultSchema: "s1", Catalog: unread}, statements)
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 2 || !spans[0].SourceColumns[analysis.ColumnResource{Database: "db", Table: "t"}] || !spans[1].SourceColumns[analysis.ColumnResource{Database: "other", Table: "t1"}] {
		t.Fatalf("spans = %+v", spans)
	}
	if spans[0].Type != analysis.Select || spans[0].Results == nil || len(spans[0].Results) != 0 || spans[0].PredicateColumns == nil || len(spans[0].PredicateColumns) != 0 {
		t.Fatalf("span = %+v, want a table-level span", spans[0])
	}
}

func TestEngineQuerySpansFoldSystemDatabasesOnlyWhenIgnoringCase(t *testing.T) {
	for ignore, want := range map[bool]analysis.QueryType{false: analysis.Select, true: analysis.SelectInfoSchema} {
		spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", IgnoreCaseSensitive: ignore}, []analysis.QuerySpanStatement{{Text: "SELECT * FROM SYS.grants_to_roles"}})
		if err != nil || spans[0].Type != want {
			t.Fatalf("ignore case %v: spans = %+v, err = %v", ignore, spans, err)
		}
	}
}

func TestEngineQuerySpansShiftASyntaxErrorToItsStatementStart(t *testing.T) {
	text := "SELECT a FROM t xx yy"
	start := analysis.Position{Line: 4, Column: 3}
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db"}, []analysis.QuerySpanStatement{{Text: text, Start: start}})
	_, _, unshifted := Engine.ValidateQuery(text)
	var syntaxErr, want *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) || !errors.As(unshifted, &want) {
		t.Fatalf("err = %v, ValidateQuery err = %v", err, unshifted)
	}
	if syntaxErr.Position != want.Position.Shift(start) || syntaxErr.Message != want.Message || syntaxErr.RawMessage != want.RawMessage {
		t.Fatalf("err = %+v, want %+v shifted to %+v", syntaxErr, want, start)
	}
}
