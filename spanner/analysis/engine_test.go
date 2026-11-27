package analysis

import (
	"context"
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/spanner/ast"
	"github.com/bytebase/omni/spanner/parser"
)

func TestEngineSplitEndsAtTheFollowingSemicolon(t *testing.T) {
	script := "SELECT 1;\n;\nSELECT 2"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || !segments[1].Empty || segments[2].Text != "\nSELECT 2" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v", s)
		}
	}
}

func TestEngineParseKeepsAStatementItRejects(t *testing.T) {
	parsed, err := Engine.Parse("SELECT 1;\n;\nSELEC 2;")
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\nSELEC 2;" || parsed[1].Node != nil {
		t.Fatalf("Parse = %+v", parsed)
	}
	if _, ok := parsed[0].Node.(*ast.QueryStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(" ; "); err != nil || parsed != nil {
		t.Fatalf("Parse of no statement = %#v, %v", parsed, err)
	}
}

func TestEngineDiagnoseDropsNotYetSupportedStubs(t *testing.T) {
	if diagnostics := Engine.Diagnose("DEFINE MACRO m 1"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of a stub = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT FROM")
	if len(diagnostics) == 0 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateAndLimit(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1; SELECT 2"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	if readOnly, _, err := Engine.ValidateQuery("DELETE FROM t WHERE true"); err != nil || readOnly {
		t.Fatalf("ValidateQuery(DELETE) = %v %v", readOnly, err)
	}
	if got := Engine.ResultLimit("SELECT 1 LIMIT 50", 10, ""); got != "SELECT 1 LIMIT 10" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT 1 FOR UPDATE", 10, ""); got != "SELECT 1 FOR UPDATE" {
		t.Fatalf("ResultLimit(FOR UPDATE) = %q", got)
	}
}

func spannerSpanCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "ds", Schemas: []*metadata.SchemaMetadata{
		{Name: "s1", Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id"}}}}},
	}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "ds" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: db, ObjectCaseSensitive: true}, nil
	}}
}

func TestEngineQuerySpansResolveSchemasInTheDefaultDatabase(t *testing.T) {
	spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "ds", Catalog: spannerSpanCatalog()}, []analysis.QuerySpanStatement{
		{Text: "SELECT id FROM other.s1.t WHERE id = 1"}, {Text: "SELECT * FROM INFORMATION_SCHEMA.TABLES"}, {Text: "SET x = 1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !spans[0].Results[0].SourceColumns[analysis.ColumnResource{Database: "ds", Schema: "s1", Table: "t", Column: "id"}] || spans[0].PredicateColumns != nil {
		t.Fatalf("span = %+v", spans[0])
	}
	if spans[1].Type != analysis.SelectInfoSchema || len(spans[1].SourceColumns) != 0 {
		t.Fatalf("system span = %+v", spans[1])
	}
	if spans[2].Type != analysis.Select {
		t.Fatalf("SET type = %v", spans[2].Type)
	}
}

func TestEngineQuerySpansWithoutAGetterReportTheDatabaseNotFound(t *testing.T) {
	spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "ds"}, []analysis.QuerySpanStatement{{Text: "SELECT id FROM t"}})
	if err != nil {
		t.Fatal(err)
	}
	if notFound := spans[0].NotFoundError; notFound == nil || *notFound.Database != "ds" {
		t.Fatalf("NotFoundError = %+v", notFound)
	}
}

func TestEngineQuerySpansReturnTheParserErrorAsIs(t *testing.T) {
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "ds", Catalog: spannerSpanCatalog()}, []analysis.QuerySpanStatement{{Text: "SELECT * FROM", Start: analysis.Position{Line: 5, Column: 1}}})
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) || err != error(parseErr) {
		t.Fatalf("err = %#v", err)
	}
}
