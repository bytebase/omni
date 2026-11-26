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
				{Name: "PUBLIC", Tables: []*metadata.TableMetadata{{Name: "ORDERS"}}},
				{Name: "SALES", Tables: []*metadata.TableMetadata{{Name: "DEALS"}}},
			},
		}}, nil
	}}
	complete := func(request analysis.CompletionRequest) []analysis.Candidate {
		return Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	}

	request := analysis.CompletionRequest{DefaultDatabase: "db", Catalog: catalog}
	got := complete(request)
	if !snowflakeOffersCandidate(got, "ORDERS", analysis.CandidateTable) || snowflakeOffersCandidate(got, "DEALS", analysis.CandidateTable) {
		t.Fatal("without a default schema, the PUBLIC schema's tables are not the ones offered")
	}
	request.DefaultSchema = "SALES"
	if !snowflakeOffersCandidate(complete(request), "DEALS", analysis.CandidateTable) {
		t.Fatal("the default schema's table is not offered")
	}
	request.Catalog = analysis.Catalog{}
	if snowflakeOffersCandidate(complete(request), "ORDERS", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func snowflakeOffersCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
