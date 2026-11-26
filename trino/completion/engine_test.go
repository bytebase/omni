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
					{Name: "s1", Tables: []*metadata.TableMetadata{{Name: "t_" + name + "_s1"}}},
					{Name: "s2", Tables: []*metadata.TableMetadata{{Name: "t_" + name + "_s2"}}},
				},
			}}, nil
		},
	}
	request := analysis.CompletionRequest{DefaultDatabase: "db", DefaultSchema: "s1", Catalog: catalog}
	got := Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	if !offersTrinoCandidate(got, "t_db_s1", analysis.CandidateTable) {
		t.Fatal("the default catalog and schema's table is not offered")
	}
	if offersTrinoCandidate(got, "t_db_s2", analysis.CandidateTable) {
		t.Fatal("a table outside the default schema is offered unqualified")
	}
	if !offersTrinoCandidate(got, "other", analysis.CandidateCatalog) {
		t.Fatal("the catalog's databases are not offered as Trino catalogs")
	}
	request.Catalog = analysis.Catalog{}
	got = Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	if offersTrinoCandidate(got, "t_db_s1", analysis.CandidateTable) || offersTrinoCandidate(got, "other", analysis.CandidateCatalog) {
		t.Fatal("a table or catalog is offered without a catalog")
	}
}

func offersTrinoCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
