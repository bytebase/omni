package completion

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestEngineCompletesFromTheCatalog(t *testing.T) {
	catalog := analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if name != "app" {
			return nil, nil
		}
		return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
			Name:    name,
			Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "orders"}}}},
		}}, nil
	}}
	request := analysis.CompletionRequest{DefaultDatabase: "app", Catalog: catalog}
	if !tidbEngineOffers(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "orders", analysis.CandidateTable) {
		t.Fatal("the default database's table is not offered")
	}
	request.DefaultDatabase = "other"
	if tidbEngineOffers(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "orders", analysis.CandidateTable) {
		t.Fatal("a table of a database other than the default is offered")
	}
	request = analysis.CompletionRequest{DefaultDatabase: "app"}
	if tidbEngineOffers(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "orders", analysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func TestEngineListsDatabases(t *testing.T) {
	catalog := analysis.Catalog{ListDatabaseNames: func(context.Context) ([]string, error) {
		return []string{"sales"}, nil
	}}
	request := analysis.CompletionRequest{Catalog: catalog}
	if !tidbEngineOffers(Engine.Complete(context.Background(), request, "USE ", 1, 4), "sales", analysis.CandidateDatabase) {
		t.Fatal("the catalog's database is not offered")
	}
}

func TestEngineOffersTheResolvedDatabaseName(t *testing.T) {
	catalog := analysis.Catalog{GetDatabase: func(_ context.Context, name string) (*analysis.Database, error) {
		if !strings.EqualFold(name, "sales") {
			return nil, nil
		}
		return &analysis.Database{Name: "Sales", Metadata: &metadata.DatabaseSchemaMetadata{
			Name:    "Sales",
			Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "t"}}}},
		}}, nil
	}}
	script := "SELECT * FROM sales.t JOIN "
	candidates := Engine.Complete(context.Background(), analysis.CompletionRequest{Catalog: catalog}, script, 1, len(script))
	if !tidbEngineOffers(candidates, "Sales", analysis.CandidateDatabase) {
		t.Fatalf("the database is not offered as the catalog names it: %v", candidates)
	}
}

func TestEngineReadOnlyDropsWriteStatements(t *testing.T) {
	request := analysis.CompletionRequest{}
	if !tidbEngineOffers(Engine.Complete(context.Background(), request, "", 1, 0), "INSERT", analysis.CandidateKeyword) {
		t.Fatal("INSERT is not offered")
	}
	request.ReadOnly = true
	candidates := Engine.Complete(context.Background(), request, "", 1, 0)
	if tidbEngineOffers(candidates, "INSERT", analysis.CandidateKeyword) {
		t.Fatal("INSERT is offered to a read-only editor")
	}
	if !tidbEngineOffers(candidates, "SELECT", analysis.CandidateKeyword) {
		t.Fatal("SELECT is not offered to a read-only editor")
	}
}

func tidbEngineOffers(candidates []analysis.Candidate, text string, typ analysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
