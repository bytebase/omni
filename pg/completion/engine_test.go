package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheCatalog(t *testing.T) {
	catalog := analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
			Name:    name,
			Schemas: []*metadata.SchemaMetadata{{Name: "public", Tables: []*metadata.TableMetadata{{Name: "orders"}}}},
		}}, nil
	}}
	request := analysis.CompletionRequest{DefaultDatabase: "db", DefaultSchema: "public", Catalog: catalog}
	if !hasCandidate(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "orders", analysis.CandidateTable) {
		t.Fatal("the catalog's table is not offered")
	}
	request.Catalog = analysis.Catalog{}
	if hasCandidate(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "orders", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func hasCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
