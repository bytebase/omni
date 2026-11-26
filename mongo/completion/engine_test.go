package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheDefaultDatabase(t *testing.T) {
	catalog := analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "test" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
			Name:    name,
			Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "users"}}}},
		}}, nil
	}}
	request := analysis.CompletionRequest{DefaultDatabase: "test", DefaultSchema: "other", Catalog: catalog}
	if !mongoOffers(Engine.Complete(context.Background(), request, "db.", 1, 3), "users", analysis.CandidateCollection) {
		t.Fatal("the default database's collection is not offered")
	}
	request.Catalog = analysis.Catalog{}
	if mongoOffers(Engine.Complete(context.Background(), request, "db.", 1, 3), "users", analysis.CandidateCollection) {
		t.Fatal("a collection is offered without a catalog")
	}
}

func TestEngineOffersAnUntypedCandidateAsAFunction(t *testing.T) {
	converted := convertCandidates([]Candidate{{Text: "f"}, {Text: "find", Type: CandidateMethod}})
	if converted[0].Type != analysis.CandidateFunction || converted[1].Type != analysis.CandidateMethod {
		t.Fatalf("got types %v and %v", converted[0].Type, converted[1].Type)
	}
}

func mongoOffers(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
