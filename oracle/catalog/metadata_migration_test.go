package catalog

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

func oracleMigration(t *testing.T, from, to *metadatapb.DatabaseSchemaMetadata) string {
	t.Helper()
	database := func(meta *metadatapb.DatabaseSchemaMetadata) *analysis.Database {
		if meta == nil {
			return nil
		}
		return &analysis.Database{Metadata: meta}
	}
	sql, err := Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{From: database(from), To: database(to)})
	if err != nil {
		t.Fatalf("MetadataMigration: %v", err)
	}
	return sql
}

func oracleOneTable() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{{
			Name:    "T",
			Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER"}},
		}},
	}}}
}

func TestEngineMetadataMigrationOracle(t *testing.T) {
	source := oracleOneTable()
	target := proto.CloneOf(source)
	target.Schemas[0].Tables[0].Columns = append(target.Schemas[0].Tables[0].Columns, &metadatapb.ColumnMetadata{Name: "NOTE", Type: "VARCHAR2(50 BYTE)", Nullable: true})

	if sql := oracleMigration(t, source, source); sql != "" {
		t.Fatalf("same schema = %q", sql)
	}
	if sql := oracleMigration(t, source, target); !strings.Contains(sql, `ALTER TABLE "APP"."T" ADD "NOTE" VARCHAR2(50 BYTE);`) || !strings.HasSuffix(sql, "\n") {
		t.Fatalf("added column = %q", sql)
	}
	if sql := oracleMigration(t, nil, source); !strings.Contains(sql, `CREATE TABLE "APP"."T"`) || !strings.HasSuffix(sql, "\n") {
		t.Fatalf("from nil = %q", sql)
	}
}

func TestEngineMetadataMigrationOracleLeavesOutIdentityAndVisibility(t *testing.T) {
	source := oracleOneTable()
	target := proto.CloneOf(source)
	table := target.Schemas[0].Tables[0]
	table.Columns = append(table.Columns, &metadatapb.ColumnMetadata{
		Name: "SEQ", Type: "NUMBER", IsIdentity: true,
		IdentityGeneration: metadatapb.ColumnMetadata_ALWAYS, IdentitySeed: 1, IdentityIncrement: 1,
	})
	table.Indexes = []*metadatapb.IndexMetadata{{Name: "IDX_T_ID", Expressions: []string{"ID"}}}

	sql := oracleMigration(t, source, target)
	if sql == "" || strings.Contains(sql, "IDENTITY") || strings.Contains(sql, "GENERATED") || strings.Contains(sql, "INVISIBLE") {
		t.Fatalf("migration = %q", sql)
	}
	if !table.Columns[1].IsIdentity || table.Indexes[0].Visible {
		t.Fatalf("the request's metadata changed: %v", table)
	}
}

func TestEngineExtractMetadataOracle(t *testing.T) {
	meta, err := Engine.ExtractMetadata(context.Background(), "CREATE TABLE t (id NUMBER);")
	want, wantErr := GetDatabaseMetadata("CREATE TABLE t (id NUMBER);")
	if err != nil || wantErr != nil || !proto.Equal(meta, want) || len(meta.GetSchemas()) == 0 {
		t.Fatalf("extracted = %v, %v", meta, err)
	}
}
