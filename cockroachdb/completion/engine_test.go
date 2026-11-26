package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheCatalog(t *testing.T) {
	catalog := analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "db" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
			Name: name,
			Schemas: []*metadata.SchemaMetadata{
				{Name: "public", Tables: []*metadata.TableMetadata{{Name: "orders"}}},
				{Name: "sales", Tables: []*metadata.TableMetadata{{Name: "deals"}}},
			},
		}}, nil
	}}
	complete := func(request analysis.CompletionRequest) []analysis.Candidate {
		return Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	}

	request := analysis.CompletionRequest{DefaultDatabase: "db", Catalog: catalog}
	got := complete(request)
	if !cockroachDBOffersCandidate(got, "orders", analysis.CandidateTable) || cockroachDBOffersCandidate(got, "deals", analysis.CandidateTable) {
		t.Fatal("without a default schema, the public schema's tables are not the ones offered")
	}
	request.DefaultSchema = "sales"
	if !cockroachDBOffersCandidate(complete(request), "deals", analysis.CandidateTable) {
		t.Fatal("the default schema's table is not offered")
	}
	request.Catalog = analysis.Catalog{}
	if cockroachDBOffersCandidate(complete(request), "orders", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func cockroachDBOffersCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
