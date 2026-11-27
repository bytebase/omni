package analysis

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineSplit(t *testing.T) {
	script := "SELECT 1;\nCREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; SELECT 2; END;\n;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || segments[1].Text != "\nCREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; SELECT 2; END;" || !segments[2].Empty {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v is not the script's text", s)
		}
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("UPSERT INTO t VALUES (1);"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start.Line != 2 || diagnostics[0].Severity != analysis.SeverityError {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateExplainAndLimit(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT * FROM t AS OF SYSTEM TIME '-10s'"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery("SELECT FROM WHERE"); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("ValidateQuery(SELECT FROM WHERE) err = %v", err)
	}
	for format, want := range map[analysis.ExplainFormat]string{
		analysis.ExplainFormatDefault: "EXPLAIN SELECT 1",
		analysis.ExplainFormatText:    "EXPLAIN SELECT 1",
		analysis.ExplainFormatJSON:    "EXPLAIN (DISTSQL) SELECT 1",
		analysis.ExplainFormatYAML:    "EXPLAIN SELECT 1",
	} {
		if got, err := Engine.Explain("SELECT 1", format); err != nil || got != want {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if _, err := Engine.Explain("EXPLAIN ANALYZE SELECT 1", analysis.ExplainFormatJSON); !errors.Is(err, analysis.ErrExplainAnalyzeExecutes) {
		t.Fatalf("Explain(EXPLAIN ANALYZE) err = %v", err)
	}
	if got := Engine.ResultLimit("SELECT * FROM t AS OF SYSTEM TIME '-10s'", 10, ""); got != "SELECT * FROM t AS OF SYSTEM TIME '-10s' LIMIT 10" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT 1; SELECT 2;", 10, ""); got != "WITH result AS (\nSELECT 1; SELECT 2\n) SELECT * FROM result LIMIT 10;" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}
}

func cockroachSpanCatalog() analysis.Catalog {
	schema := func(name string, tables ...string) *metadata.SchemaMetadata {
		s := &metadata.SchemaMetadata{Name: name}
		for _, table := range tables {
			s.Tables = append(s.Tables, &metadata.TableMetadata{Name: table, Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "INT8"}}})
		}
		return s
	}
	db := &metadata.DatabaseSchemaMetadata{Name: "db", SearchPath: `"$user", s1, public`, Schemas: []*metadata.SchemaMetadata{schema("public", "t", "p"), schema("s1", "t"), schema("s2", "t")}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: db, ObjectCaseSensitive: true}, nil
	}}
}

func TestEngineQuerySpansPutASelectedSchemaBeforePublic(t *testing.T) {
	for _, tc := range []struct {
		schema, statement, want string
	}{
		{"", "SELECT id FROM t", "s1"},
		{"public", "SELECT id FROM t", "public"},
		{"s2", "SELECT id FROM t", "s2"},
		{"s2", "SELECT id FROM p", "public"},
	} {
		spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", DefaultSchema: tc.schema, Catalog: cockroachSpanCatalog()}, []analysis.QuerySpanStatement{{Text: tc.statement}})
		if err != nil {
			t.Fatal(err)
		}
		table := strings.TrimPrefix(tc.statement, "SELECT id FROM ")
		if !spans[0].Results[0].SourceColumns[analysis.ColumnResource{Database: "db", Schema: tc.want, Table: table, Column: "id"}] {
			t.Fatalf("schema %q, %s: spans = %+v", tc.schema, tc.statement, spans[0])
		}
	}
}

func TestEngineQuerySpansNeedASyncedDefaultDatabase(t *testing.T) {
	statements := []analysis.QuerySpanStatement{{Text: "SELECT 1"}}
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "other", Catalog: cockroachSpanCatalog()}, statements)
	if err == nil || err.Error() != `database metadata for database "other" not found (database not synced)` {
		t.Fatalf("err = %v", err)
	}
	if _, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db"}, statements); err == nil {
		t.Fatal("a catalog without GetDatabase must fail")
	}
}

func TestEngineQuerySpansPositionASyntaxErrorInItsStatement(t *testing.T) {
	text := "SELECT 1;\nSELECT * FROM WHERE"
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "db", Catalog: cockroachSpanCatalog()}, []analysis.QuerySpanStatement{{Text: text, Start: analysis.Position{Line: 5, Column: 3}}})
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != analysis.PositionAt(text, strings.Index(text, "WHERE")) || syntaxErr.RawMessage != syntaxErr.Message || !strings.HasPrefix(err.Error(), "failed to get query span from statement: "+text+": ") {
		t.Fatalf("err = %v at %+v", err, syntaxErr.Position)
	}
}
