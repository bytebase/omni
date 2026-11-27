package analysis

import (
	"context"
	"errors"
	"slices"
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

func cockroachChangeRequest(schema string) analysis.ChangeRequest {
	table := func(name, index string) *metadata.TableMetadata {
		return &metadata.TableMetadata{Name: name, Indexes: []*metadata.IndexMetadata{{Name: index}}}
	}
	db := &metadata.DatabaseSchemaMetadata{Name: "db", SearchPath: `"$user", app, public`, Schemas: []*metadata.SchemaMetadata{
		{Name: "app", Tables: []*metadata.TableMetadata{table("t", "t_idx")}},
		{Name: "public", Tables: []*metadata.TableMetadata{table("t2", "t2_idx")}, ExternalTables: []*metadata.ExternalTableMetadata{{Name: "ext"}}},
	}}
	return analysis.ChangeRequest{DefaultDatabase: "db", DefaultSchema: schema, Database: &analysis.Database{Name: "db", Metadata: db, ObjectCaseSensitive: true}}
}

func parseCockroachForChanges(t *testing.T, script string) []analysis.ParsedStatement {
	t.Helper()
	parsed, err := ParseScript(script)
	if err != nil {
		t.Fatal(err)
	}
	statements := make([]analysis.ParsedStatement, 0, len(parsed))
	for _, stmt := range parsed {
		statements = append(statements, analysis.ParsedStatement{Node: stmt})
	}
	return statements
}

func cockroachResourceSchemas(changes *analysis.Changes) []string {
	var schemas []string
	for _, resource := range changes.Resources {
		schemas = append(schemas, resource.Schema+"."+resource.Table)
	}
	return schemas
}

func TestEngineExtractChangesResolveAlongTheSelectedSchemaAloneOrTheSyncedSearchPath(t *testing.T) {
	noMetadata := cockroachChangeRequest("")
	noMetadata.Database = nil
	for name, tc := range map[string]struct {
		request analysis.ChangeRequest
		want    string
	}{
		// ext is an external table of public, which the lookup does not count.
		"synced search path":      {cockroachChangeRequest(""), "app.t,public.t2,app.ext"},
		"selected schema alone":   {cockroachChangeRequest("app"), "app.t,app.t2,app.ext"},
		"public without metadata": {noMetadata, "public.t,public.t2,public.ext"},
	} {
		changes, err := Engine.ExtractChanges(tc.request, parseCockroachForChanges(t, "UPDATE t SET a = 1; UPDATE t2 SET a = 1; UPDATE ext SET a = 1;"))
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Join(cockroachResourceSchemas(changes), ","); got != tc.want {
			t.Fatalf("%s: resources = %s", name, got)
		}
	}
	empty := cockroachChangeRequest("")
	empty.Database.Metadata.SearchPath = ""
	changes, err := Engine.ExtractChanges(empty, parseCockroachForChanges(t, "UPDATE t SET a = 1;"))
	if err != nil || strings.Join(cockroachResourceSchemas(changes), ",") != "public.t" {
		t.Fatalf("empty search path: changes = %+v, err = %v", changes, err)
	}
}

func TestEngineExtractChangesFindAnIndexTableOrItsOtherDatabase(t *testing.T) {
	changes, err := Engine.ExtractChanges(cockroachChangeRequest(""), parseCockroachForChanges(t, "DROP INDEX t_idx; DROP INDEX \"T_idx\"; DROP INDEX otherdb.public.x_idx;"))
	if err != nil {
		t.Fatal(err)
	}
	want := []analysis.ChangedResource{{Database: "db", Schema: "app", Table: "t"}, {Database: "otherdb", DatabaseOnly: true}}
	if !slices.Equal(changes.Resources, want) {
		t.Fatalf("resources = %+v", changes.Resources)
	}
}

func TestEngineExtractChangesKeepASearchPathOnlyAfterASet(t *testing.T) {
	changes, err := Engine.ExtractChanges(cockroachChangeRequest(""), parseCockroachForChanges(t, "DELETE FROM t; SET search_path = public; DELETE FROM t2;"))
	if err != nil {
		t.Fatal(err)
	}
	if len(changes.DMLStatements) != 2 || changes.DMLStatements[0].SearchPath != nil || !slices.Equal(changes.DMLStatements[1].SearchPath, []string{"public"}) {
		t.Fatalf("DML statements = %+v", changes.DMLStatements)
	}
}

func TestEngineExtractChangesNeedCockroachDBNodes(t *testing.T) {
	for _, node := range []any{nil, "not a node"} {
		_, err := Engine.ExtractChanges(cockroachChangeRequest(""), []analysis.ParsedStatement{{Node: node}})
		if err == nil || err.Error() != "expected CockroachDB AST" {
			t.Fatalf("node %v: err = %v", node, err)
		}
	}
}
