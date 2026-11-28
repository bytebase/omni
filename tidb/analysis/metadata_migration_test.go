package analysis

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func tidbSnapshot(db string, column *metadata.ColumnMetadata, views ...*metadata.ViewMetadata) *analysis.Database {
	return &analysis.Database{Name: db, ObjectCaseSensitive: true, Metadata: &metadata.DatabaseSchemaMetadata{
		Name: db,
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{column}}},
			Views:  views,
		}},
	}}
}

func tidbMigration(t *testing.T, from, to *analysis.Database) (string, error) {
	t.Helper()
	return Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{From: from, To: to})
}

func TestEngineMetadataMigrationTiDBAcrossDatabaseNames(t *testing.T) {
	id := &metadata.ColumnMetadata{Name: "id", Type: "int"}
	view := func(db string) *metadata.ViewMetadata {
		return &metadata.ViewMetadata{Name: "v", Definition: "select `" + db + "`.`t`.`id` AS `id` from `" + db + "`.`t`"}
	}
	if sql, err := tidbMigration(t, tidbSnapshot("prod", id, view("prod")), tidbSnapshot("staging", id, view("staging"))); sql != "" || err != nil {
		t.Fatalf("same schema = %q, %v", sql, err)
	}
	sql, err := tidbMigration(t, tidbSnapshot("prod", id), tidbSnapshot("staging", id, view("staging")))
	if err != nil || strings.Contains(sql, "staging") || !strings.Contains(sql, "VIEW `v` (`id`) AS SELECT `t`.`id` AS `id` FROM `t`") {
		t.Fatalf("new view = %q, %v", sql, err)
	}
}

func TestEngineMetadataMigrationTiDBRefusesAChangedUnloadableObject(t *testing.T) {
	id := &metadata.ColumnMetadata{Name: "id", Type: "int"}
	unparsable := &metadata.ColumnMetadata{Name: "id", Type: "int", Default: "((("}
	if _, err := tidbMigration(t, tidbSnapshot("d", id), tidbSnapshot("d", unparsable)); err == nil || !strings.Contains(err.Error(), "table:t") {
		t.Fatalf("changed table error = %v", err)
	}
	if sql, err := tidbMigration(t, tidbSnapshot("d", unparsable), tidbSnapshot("d", unparsable)); sql != "" || err != nil {
		t.Fatalf("same table = %q, %v", sql, err)
	}
}

func TestEngineExtractMetadataTiDB(t *testing.T) {
	ctx := context.Background()
	meta, err := Engine.ExtractMetadata(ctx, "CREATE TABLE t (id INT PRIMARY KEY);\nCREATE VIEW v AS SELECT id FROM t;")
	if err != nil || meta.GetName() != "" || len(meta.GetSchemas()) != 1 || len(meta.GetSchemas()[0].GetViews()) != 1 {
		t.Fatalf("scratch database = %v, %v", meta, err)
	}
	if definition := meta.GetSchemas()[0].GetViews()[0].GetDefinition(); strings.Contains(definition, scratchDatabase) {
		t.Fatalf("scratch view definition = %q", definition)
	}
	meta, err = Engine.ExtractMetadata(ctx, "CREATE DATABASE app;\nUSE app;\nCREATE TABLE t (id INT PRIMARY KEY);\nCREATE VIEW v AS SELECT id FROM t;")
	if err != nil || meta.GetName() != "app" || len(meta.GetSchemas()) != 1 || len(meta.GetSchemas()[0].GetViews()) != 1 {
		t.Fatalf("selected database = %v, %v", meta, err)
	}
	if definition := meta.GetSchemas()[0].GetViews()[0].GetDefinition(); !strings.Contains(definition, "`app`.") {
		t.Fatalf("selected view definition = %q", definition)
	}
}
