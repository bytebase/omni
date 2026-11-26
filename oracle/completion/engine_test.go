package completion

import (
	"context"
	"testing"

	omnianalysis "github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func oracleEngineCatalog() omnianalysis.Catalog {
	return omnianalysis.Catalog{GetDatabase: func(_ context.Context, name string) (*omnianalysis.Database, error) {
		if name != "HR" {
			return nil, nil
		}
		return &omnianalysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
			Name:    name,
			Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "ORDERS"}}}},
		}}, nil
	}}
}

func TestEngineUsesTheDefaultDatabaseAsSchema(t *testing.T) {
	request := omnianalysis.CompletionRequest{DefaultDatabase: "HR", DefaultSchema: "OTHER", Catalog: oracleEngineCatalog()}
	if !offersOracleCandidate(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "ORDERS", omnianalysis.CandidateTable) {
		t.Fatal("the default database's table is not offered")
	}
	request.Catalog = omnianalysis.Catalog{}
	if offersOracleCandidate(Engine.Complete(context.Background(), request, "SELECT * FROM ", 1, 14), "ORDERS", omnianalysis.CandidateTable) {
		t.Fatal("a table is offered without a catalog")
	}
}

func TestEngineReadOnlyOffersNothingInAChange(t *testing.T) {
	request := omnianalysis.CompletionRequest{DefaultDatabase: "HR", Catalog: oracleEngineCatalog()}
	if Engine.Complete(context.Background(), request, "DROP TABLE ", 1, 11) == nil {
		t.Fatal("a change offers nothing outside a read-only editor")
	}
	request.ReadOnly = true
	if got := Engine.Complete(context.Background(), request, "DROP TABLE ", 1, 11); got != nil {
		t.Fatalf("a read-only editor offers %v in a change", got)
	}
}

func TestEngineReportsPLSQLVariablesAsUserVariables(t *testing.T) {
	request := omnianalysis.CompletionRequest{DefaultDatabase: "HR", Catalog: oracleEngineCatalog()}
	script := "DECLARE v_id NUMBER; BEGIN SELECT 1 INTO "
	got := Engine.Complete(context.Background(), request, script, 1, len(script))
	if !offersOracleCandidate(got, "V_ID", omnianalysis.CandidateUserVariable) {
		t.Fatalf("V_ID is not offered as a user variable: %v", got)
	}
	if offersOracleCandidate(got, "V_ID", omnianalysis.CandidateVariable) {
		t.Fatal("V_ID is offered as a system variable")
	}
}

func offersOracleCandidate(candidates []omnianalysis.Candidate, text string, typ omnianalysis.CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}
