package analysis

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/metadata"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT * FROM a;\n;\n  SELECT * FROM b;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT * FROM a" || !segments[1].Empty || segments[2].Text != "SELECT * FROM b" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v is not the script's text", s)
		}
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "SELECT * FROM b" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(";"); err != nil || parsed != nil {
		t.Fatalf("Parse of no statement = %+v, %v", parsed, err)
	}
}

func TestEngineSyntaxErrorShiftsLineAndColumn(t *testing.T) {
	_, err := Engine.Parse("SELECT * FROM t;\n  SELCT * FROM users")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 3}) || syntaxErr.Message != "Syntax error at line 2: "+syntaxErr.RawMessage {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineValidateQueryReturnsTheParsersError(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT * FROM t"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery("SELECT FROM"); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("ValidateQuery(SELECT FROM) err = %v", err)
	}
}

// cassandraSpanCatalog holds keyspace ks, whose table Users is in its second
// schema, and fails on keyspace broken.
func cassandraSpanCatalog() analysis.Catalog {
	ks := &metadata.DatabaseSchemaMetadata{Name: "ks", Schemas: []*metadata.SchemaMetadata{
		{Name: "other"},
		{Name: "", Tables: []*metadata.TableMetadata{{Name: "Users", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "email"}}}}},
	}}
	return analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		switch name {
		case "ks":
			return &analysis.Database{Name: name, Metadata: ks, ObjectCaseSensitive: true}, nil
		case "broken":
			return nil, errors.New("metadata store is down")
		default:
			return nil, nil
		}
	}}
}

func TestEngineQuerySpansExpandAStarFromTheSyncedKeyspace(t *testing.T) {
	statements := []analysis.QuerySpanStatement{{Text: `SELECT * FROM "Users" WHERE id = 1`}, {Text: "SELECT * FROM users"}}
	spans, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "ks", DefaultSchema: "other", Catalog: cassandraSpanCatalog()}, statements)
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 2 || len(spans[0].Results) != 2 || spans[0].Results[1].Name != "email" || !spans[0].PredicateColumns[analysis.ColumnResource{Database: "ks", Table: "Users", Column: "id"}] {
		t.Fatalf("spans = %+v", spans)
	}
	if len(spans[1].Results) != 1 || !spans[1].Results[0].SelectAsterisk {
		t.Fatalf("a star over a table the keyspace does not hold as written = %+v", spans[1].Results)
	}

	spans, err = Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "ks"}, statements[:1])
	if err != nil || len(spans[0].Results) != 1 || !spans[0].Results[0].SelectAsterisk {
		t.Fatalf("without a catalog: spans = %+v, err = %v", spans, err)
	}
}

func TestEngineQuerySpansReturnACatalogError(t *testing.T) {
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "broken", Catalog: cassandraSpanCatalog()}, []analysis.QuerySpanStatement{{Text: "SELECT * FROM t"}})
	if err == nil || err.Error() != `failed to get database metadata for keyspace "broken": metadata store is down` {
		t.Fatalf("err = %v", err)
	}
}

func TestEngineQuerySpansPositionASyntaxErrorInItsStatement(t *testing.T) {
	text := "\n  SELCT * FROM users"
	_, err := Engine.QuerySpans(context.Background(), analysis.QuerySpanRequest{DefaultDatabase: "ks"}, []analysis.QuerySpanStatement{{Text: text, Start: analysis.Position{Line: 5, Column: 4}}})
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 3}) || !strings.HasPrefix(syntaxErr.Message, "Syntax error at line 2: ") {
		t.Fatalf("err = %+v", syntaxErr)
	}
}
