package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheCatalog(t *testing.T) {
	catalog := analysis.Catalog{
		ListDatabaseNames: func(context.Context) ([]string, error) { return []string{"db", "other"}, nil },
		GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
			return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
				Name: name,
				Schemas: []*metadata.SchemaMetadata{
					{Name: "dbo", Tables: []*metadata.TableMetadata{{Name: "in_dbo"}}},
					{Name: "s1", Tables: []*metadata.TableMetadata{{Name: "in_s1"}}},
				},
			}}, nil
		},
	}
	request := analysis.CompletionRequest{DefaultDatabase: "db", DefaultSchema: "s1", Catalog: catalog}
	if !offersMSSQLCandidate(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "in_s1", analysis.CandidateTable) {
		t.Fatal("the default schema's table is not offered")
	}
	if !offersMSSQLCandidate(Engine.Complete(context.Background(), request, "USE ", 1, 4), "other", analysis.CandidateDatabase) {
		t.Fatal("the catalog's database names are not offered")
	}
	request.Catalog = analysis.Catalog{}
	if offersMSSQLCandidate(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "in_s1", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func offersMSSQLCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
