package analysis

import (
	"context"
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/snowflake/ast"
)

func TestEngineSplitKeepsEmptySegments(t *testing.T) {
	script := "SELECT 1;;SELECT 2"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || !segments[1].Empty || segments[2].Text != "SELECT 2" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
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
	parsed, err = Engine.Parse("-- c\nSELECT 1;;SELECT 2;")
	if err != nil || len(parsed) != 2 {
		t.Fatalf("Parse skipping an empty segment = %#v, %v", parsed, err)
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

// snowflakeSpanCatalog lists DB1 and GONE; GONE's lookup reports it not found.
func snowflakeSpanCatalog() analysis.Catalog {
	db := &metadata.DatabaseSchemaMetadata{Name: "DB1", Schemas: []*metadata.SchemaMetadata{{
		Name:   "PUBLIC",
		Tables: []*metadata.TableMetadata{{Name: "T1", Columns: []*metadata.ColumnMetadata{{Name: "A"}}}, {Name: "t2", Columns: []*metadata.ColumnMetadata{{Name: "B"}}}},
	}}}
	return analysis.Catalog{
		ListDatabaseNames: func(context.Context) ([]string, error) { return []string{"DB1", "GONE"}, nil },
		GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
			if name == "GONE" {
				return nil, &analysis.ResourceNotFoundError{Database: &name}
			}
			if name != "DB1" {
				return nil, nil
			}
			return &analysis.Database{Name: name, Metadata: db, ObjectCaseSensitive: false}, nil
		},
	}
}

func snowflakeQuerySpan(t *testing.T, statement, schema string) (*analysis.ResolvedSpan, error) {
	t.Helper()
	spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "DB1", DefaultSchema: schema, Catalog: snowflakeSpanCatalog()}, []analysis.QuerySpanStatement{{Text: statement}})
	if err != nil {
		return nil, err
	}
	return spans[0], nil
}

func TestEngineQuerySpansMatchTheCatalogsSpellingExactly(t *testing.T) {
	span, err := snowflakeQuerySpan(t, `SELECT A FROM t1`, "")
	if err != nil || !span.Results[0].SourceColumns[analysis.ColumnResource{Database: "DB1", Schema: "PUBLIC", Table: "T1", Column: "A"}] {
		t.Fatalf("span = %+v, err = %v", span, err)
	}
	// A database that ignores case still needs the folded name spelled as synced.
	if _, err := snowflakeQuerySpan(t, `SELECT B FROM t2`, ""); err == nil || err.Error() != "table DB1.PUBLIC.T2 is not found" {
		t.Fatalf("err = %v", err)
	}
	if _, err := snowflakeQuerySpan(t, `SELECT B FROM "t2"`, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := snowflakeQuerySpan(t, `SELECT A FROM "public".T1`, ""); err == nil || err.Error() != "table DB1.public.T1 is not found" {
		t.Fatalf("err = %v", err)
	}
}

func TestEngineQuerySpansKeepTheAccessedTablesOfANotFoundSpan(t *testing.T) {
	span, err := snowflakeQuerySpan(t, `SELECT * FROM T1 JOIN GONE.PUBLIC.T1 USING (A)`, "")
	if err != nil || span.NotFoundError == nil || len(span.Results) != 0 || !span.SourceColumns[analysis.ColumnResource{Database: "GONE", Schema: "PUBLIC", Table: "T1"}] {
		t.Fatalf("span = %+v, err = %v", span, err)
	}
}

func TestEngineQuerySpansNeedACatalogToResolveATable(t *testing.T) {
	if _, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "DB1"}, []analysis.QuerySpanStatement{{Text: "SELECT 1"}}); err != nil {
		t.Fatalf("a statement without tables needs no catalog: %v", err)
	}
	if _, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "DB1"}, []analysis.QuerySpanStatement{{Text: "SELECT * FROM T1"}}); err == nil {
		t.Fatal("a catalog that cannot look tables up must fail")
	}
}
