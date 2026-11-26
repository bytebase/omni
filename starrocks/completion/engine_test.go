package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheCatalog(t *testing.T) {
	catalog := analysis.Catalog{
		ListDatabaseNames: func(context.Context) ([]string, error) { return []string{"db", "otherdb"}, nil },
		GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
			return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
				Name:    name,
				Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "orders_" + name}}}},
			}}, nil
		},
	}
	request := analysis.CompletionRequest{DefaultDatabase: "db", DefaultSchema: "otherdb", Catalog: catalog}
	got := Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	if !offersStarRocksCandidate(got, "orders_db", analysis.CandidateTable) {
		t.Fatal("the default database's table is not offered")
	}
	if offersStarRocksCandidate(got, "orders_otherdb", analysis.CandidateTable) {
		t.Fatal("the default schema is taken for a database")
	}
	if !offersStarRocksCandidate(got, "otherdb", analysis.CandidateDatabase) {
		t.Fatal("the catalog's databases are not offered")
	}
	request.Catalog = analysis.Catalog{}
	got = Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14)
	if offersStarRocksCandidate(got, "orders_db", analysis.CandidateTable) || offersStarRocksCandidate(got, "otherdb", analysis.CandidateDatabase) {
		t.Fatal("a table or database is offered without a catalog")
	}
}

func offersStarRocksCandidate(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
