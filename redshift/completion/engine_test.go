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
				{Name: "sales", Tables: []*metadata.TableMetadata{{Name: "deals", Columns: []*metadata.ColumnMetadata{{Name: "amount", Type: "int"}}}}},
			},
		}}, nil
	}}
	complete := func(request analysis.CompletionRequest, script string, caretOffset int) []analysis.Candidate {
		return Engine.Complete(context.Background(), request, script, 1, caretOffset)
	}

	request := analysis.CompletionRequest{DefaultDatabase: "db", Catalog: catalog}
	if !redshiftOffersCandidate(complete(request, "SELECT * FROM ", 14), "orders", analysis.CandidateTable) {
		t.Fatal("the catalog's table is not offered")
	}
	if redshiftOffersCandidate(complete(request, "SELECT  FROM deals", 7), "amount", analysis.CandidateColumn) {
		t.Fatal("an unqualified name resolves outside the search path")
	}
	request.DefaultSchema = "sales"
	if !redshiftOffersCandidate(complete(request, "SELECT  FROM deals", 7), "amount", analysis.CandidateColumn) {
		t.Fatal("an unqualified name does not resolve in the default schema")
	}
	request.Catalog = analysis.Catalog{}
	if redshiftOffersCandidate(complete(request, "SELECT * FROM ", 14), "orders", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func redshiftOffersCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
