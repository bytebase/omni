package analysis

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func mariadbMigration(t *testing.T, from, to *metadata.DatabaseSchemaMetadata) (string, error) {
	t.Helper()
	return Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{
		From: &analysis.Database{Metadata: from, ObjectCaseSensitive: true},
		To:   &analysis.Database{Metadata: to, ObjectCaseSensitive: true},
	})
}

func TestEngineMetadataMigrationMariaDBRefusesAChangedUnloadableObject(t *testing.T) {
	withRoutine := func(body string) *metadata.DatabaseSchemaMetadata {
		return &metadata.DatabaseSchemaMetadata{Name: "d", Schemas: []*metadata.SchemaMetadata{{
			Procedures: []*metadata.ProcedureMetadata{{Name: "p", Definition: body}},
		}}}
	}
	if _, err := mariadbMigration(t, withRoutine("BEGIN\n  SELECT 1;\nEND"), withRoutine("BEGIN\n  SELECT 2;\nEND")); err == nil || !strings.Contains(err.Error(), "procedure:p") {
		t.Fatalf("changed routine error = %v", err)
	}
	if sql, err := mariadbMigration(t, withRoutine("BEGIN\n  SELECT 1;\nEND"), withRoutine("BEGIN\n  SELECT 1;\nEND")); sql != "" || err != nil {
		t.Fatalf("same routine = %q, %v", sql, err)
	}
}

func TestEngineMetadataMigrationMariaDBTableCharsets(t *testing.T) {
	source := &metadata.DatabaseSchemaMetadata{Name: "d", Schemas: []*metadata.SchemaMetadata{{}}}
	target := func(charset, collation string) *metadata.DatabaseSchemaMetadata {
		return &metadata.DatabaseSchemaMetadata{Name: "d", Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name: "t", Engine: "InnoDB", Charset: charset, Collation: collation,
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}},
			}},
		}}}
	}
	sql, err := mariadbMigration(t, source, target("latin1", "latin1_swedish_ci"))
	if err != nil || !strings.Contains(sql, "DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci") {
		t.Fatalf("explicit charset = %q, %v", sql, err)
	}
	sql, err = mariadbMigration(t, source, target("", ""))
	if err != nil || !strings.Contains(sql, "CREATE TABLE `t`") || strings.Contains(sql, "CHARSET=") || strings.Contains(sql, "COLLATE=") {
		t.Fatalf("implicit charset = %q, %v", sql, err)
	}
}

func TestEngineExtractMetadataMariaDB(t *testing.T) {
	meta, err := Engine.ExtractMetadata(context.Background(), "CREATE TABLE t (id int) WITH SYSTEM VERSIONING;")
	if err != nil || meta.GetName() != "tmp" || len(meta.GetSchemas()) != 1 || len(meta.GetSchemas()[0].GetTables()) != 1 {
		t.Fatalf("extracted = %v, %v", meta, err)
	}
}
