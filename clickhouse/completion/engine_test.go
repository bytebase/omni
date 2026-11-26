package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheCatalog(t *testing.T) {
	catalog := analysis.Catalog{
		GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
			return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
				Name:    name,
				Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "orders"}}}},
			}}, nil
		},
		ListDatabaseNames: func(context.Context) ([]string, error) { return []string{"db", "other"}, nil },
	}
	complete := func(request analysis.CompletionRequest) []analysis.Candidate {
		return Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	}

	request := analysis.CompletionRequest{DefaultDatabase: "db", Catalog: catalog}
	got := complete(request)
	if !clickHouseOffersCandidate(got, "orders", analysis.CandidateTable) || !clickHouseOffersCandidate(got, "other", analysis.CandidateDatabase) {
		t.Fatal("the catalog's table and databases are not offered")
	}
	request.Catalog.ListDatabaseNames = nil
	if clickHouseOffersCandidate(complete(request), "other", analysis.CandidateDatabase) {
		t.Fatal("a database is offered that the catalog cannot list")
	}
	request.Catalog = analysis.Catalog{}
	if clickHouseOffersCandidate(complete(request), "orders", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func clickHouseOffersCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
