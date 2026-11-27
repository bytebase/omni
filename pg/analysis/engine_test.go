package analysis

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
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

func spanEngineCatalog() analysis.Catalog {
	schema := func(name string) *metadata.SchemaMetadata {
		return &metadata.SchemaMetadata{Name: name, Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}}}
	}
	db := &metadata.DatabaseSchemaMetadata{Name: "db", SearchPath: `"$user", public`, Schemas: []*metadata.SchemaMetadata{schema("public"), schema("s1")}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: db, ObjectCaseSensitive: true}, nil
	}}
}

func TestEngineQuerySpansResolveAlongTheSelectedSchema(t *testing.T) {
	statements := []analysis.QuerySpanStatement{{Text: "SELECT id FROM t;"}, {Text: "\nSELECT id FROM public.t;"}}
	for schema, want := range map[string]string{"": "public", "s1": "s1"} {
		spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", DefaultSchema: schema, Catalog: spanEngineCatalog()}, statements)
		if err != nil {
			t.Fatal(err)
		}
		if len(spans) != 2 || !spans[0].Results[0].SourceColumns[analysis.ColumnResource{Database: "db", Schema: want, Table: "t", Column: "id"}] {
			t.Fatalf("schema %q: spans = %+v", schema, spans)
		}
	}
}

func TestEngineQuerySpansNeedASyncedDefaultDatabase(t *testing.T) {
	statements := []analysis.QuerySpanStatement{{Text: "SELECT 1"}}
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "other", Catalog: spanEngineCatalog()}, statements)
	if err == nil || err.Error() != `database metadata for database "other" not found (database not synced)` {
		t.Fatalf("err = %v", err)
	}
	if _, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db"}, statements); err == nil {
		t.Fatal("a catalog without GetDatabase must fail")
	}
}

func TestEngineQuerySpansPositionASyntaxErrorInItsStatement(t *testing.T) {
	text := "SELECT 1;\nSELECT * FROM"
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", Catalog: spanEngineCatalog()}, []analysis.QuerySpanStatement{{Text: text, Start: analysis.Position{Line: 5, Column: 1}}})
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != analysis.PositionAt(text, len(text)-1) || !strings.HasPrefix(err.Error(), "failed to get query span from statement: "+text+": ") {
		t.Fatalf("err = %v at %+v", err, syntaxErr.Position)
	}
}

func changeEngineRequest(schema string) analysis.ChangeRequest {
	table := func(name, index string) *metadata.TableMetadata {
		return &metadata.TableMetadata{Name: name, Indexes: []*metadata.IndexMetadata{{Name: index}}}
	}
	db := &metadata.DatabaseSchemaMetadata{Name: "db", SearchPath: `"$user", public`, Schemas: []*metadata.SchemaMetadata{
		{Name: "public", Tables: []*metadata.TableMetadata{table("t", "t_idx")}},
		{Name: "s1", Tables: []*metadata.TableMetadata{table("t", "s1_idx")}},
	}}
	return analysis.ChangeRequest{DefaultDatabase: "db", DefaultSchema: schema, Database: &analysis.Database{Name: "db", Metadata: db, ObjectCaseSensitive: true}}
}

func parseForChanges(t *testing.T, script string) []analysis.ParsedStatement {
	t.Helper()
	statements, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	return statements
}

func TestEngineExtractChangesResolveAlongTheSelectedSchema(t *testing.T) {
	for schema, want := range map[string]string{"": "public", "s1": "s1"} {
		changes, err := Engine.ExtractChanges(changeEngineRequest(schema), parseForChanges(t, "UPDATE t SET a = 1;"))
		if err != nil {
			t.Fatal(err)
		}
		if len(changes.Resources) != 1 || changes.Resources[0].Schema != want || changes.DMLCount != 1 {
			t.Fatalf("schema %q: changes = %+v", schema, changes)
		}
	}
}

func TestEngineExtractChangesFindTheTableOfAnIndex(t *testing.T) {
	changes, err := Engine.ExtractChanges(changeEngineRequest(""), parseForChanges(t, "DROP INDEX s1_idx; DROP INDEX s1.s1_idx;"))
	if err != nil {
		t.Fatal(err)
	}
	if len(changes.Resources) != 1 || changes.Resources[0].Schema != "s1" || changes.Resources[0].Table != "t" {
		t.Fatalf("changes = %+v", changes.Resources)
	}
}

func TestEngineExtractChangesNeedSyncedMetadataAndPostgreSQLNodes(t *testing.T) {
	request := changeEngineRequest("")
	request.Database = nil
	if _, err := Engine.ExtractChanges(request, nil); err == nil || err.Error() != `database metadata for database "db" not found (database not synced)` {
		t.Fatalf("err = %v", err)
	}
	for _, node := range []any{nil, "not a node"} {
		_, err := Engine.ExtractChanges(changeEngineRequest(""), []analysis.ParsedStatement{{Node: node}})
		if err == nil || err.Error() != "expected AST for PostgreSQL" {
			t.Fatalf("node %v: err = %v", node, err)
		}
	}
}
