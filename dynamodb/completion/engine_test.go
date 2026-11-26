package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheDefaultDatabase(t *testing.T) {
	catalog := analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
			Name:    name,
			Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "Music", Columns: []*metadata.ColumnMetadata{{Name: "Artist"}}}}}},
		}}, nil
	}}
	request := analysis.CompletionRequest{DefaultDatabase: "db", Catalog: catalog}
	if !dynamodbOffers(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "Music", analysis.CandidateTable) {
		t.Fatal("the default database's table is not offered")
	}
	// The caret offset counts runes: the string's two runes are six bytes.
	if !dynamodbOffers(Engine.Complete(context.Background(), request, "-- note\nSELECT '音乐',  FROM Music", 2, 13), "Artist", analysis.CandidateColumn) {
		t.Fatal("the table's column is not offered")
	}
	request.DefaultDatabase = "other"
	if dynamodbOffers(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "Music", analysis.CandidateTable) {
		t.Fatal("a table outside the default database is offered")
	}
	request = analysis.CompletionRequest{DefaultDatabase: "db"}
	if dynamodbOffers(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "Music", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func TestEngineReadOnlyDropsChangeKeywords(t *testing.T) {
	request := analysis.CompletionRequest{}
	if !dynamodbOffers(Engine.Complete(context.Background(), request, "INS", 1, 3), "INSERT", analysis.CandidateKeyword) {
		t.Fatal("INSERT is not offered")
	}
	request.ReadOnly = true
	candidates := Engine.Complete(context.Background(), request, "S", 1, 1)
	if dynamodbOffers(candidates, "SET", analysis.CandidateKeyword) {
		t.Fatal("SET is offered to a read-only editor")
	}
	if !dynamodbOffers(candidates, "SELECT", analysis.CandidateKeyword) {
		t.Fatal("SELECT is not offered to a read-only editor")
	}
}

func dynamodbOffers(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
