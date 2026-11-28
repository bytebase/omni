package analysis

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func oceanBaseMigration(t *testing.T, from, to *metadata.SchemaMetadata) (string, error) {
	t.Helper()
	database := func(schema *metadata.SchemaMetadata) *analysis.Database {
		return &analysis.Database{Metadata: &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{schema}}, ObjectCaseSensitive: true}
	}
	return Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{From: database(from), To: database(to)})
}

func TestEngineMetadataMigrationOceanBase(t *testing.T) {
	id := &metadata.ColumnMetadata{Name: "id", Type: "int"}
	extra := &metadata.ColumnMetadata{Name: "name", Type: "varchar(64)", Nullable: true, Default: "NULL"}
	sql, err := oceanBaseMigration(t,
		&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{id}}}},
		&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{id, extra}}}})
	if err != nil || sql != "ALTER TABLE `t` ADD COLUMN `name` varchar(64) DEFAULT NULL" {
		t.Fatalf("migration = %q, %v", sql, err)
	}
}

func TestEngineMetadataMigrationOceanBaseRefusesEventWithoutDefinition(t *testing.T) {
	withEvent := &metadata.SchemaMetadata{Events: []*metadata.EventMetadata{{Name: "e"}}}
	if _, err := oceanBaseMigration(t, withEvent, &metadata.SchemaMetadata{}); err == nil || !strings.Contains(err.Error(), "event:e") {
		t.Fatalf("dropped event error = %v", err)
	}
	if _, err := oceanBaseMigration(t, &metadata.SchemaMetadata{}, withEvent); err == nil || !strings.Contains(err.Error(), "event:e") {
		t.Fatalf("added event error = %v", err)
	}
}

func TestEngineExtractMetadataOceanBase(t *testing.T) {
	meta, err := Engine.ExtractMetadata(context.Background(), "CREATE TABLE `t` (`a` int NOT NULL, PRIMARY KEY (`a`));")
	if err != nil || meta.GetName() != "tmp" || len(meta.GetSchemas()) != 1 || len(meta.GetSchemas()[0].GetTables()) != 1 {
		t.Fatalf("extracted = %v, %v", meta, err)
	}
}
