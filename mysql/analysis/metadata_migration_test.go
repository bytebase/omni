package analysis

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func mysqlSnapshot(name string, columns []*metadata.ColumnMetadata, views ...*metadata.ViewMetadata) *analysis.Database {
	return &analysis.Database{Name: name, Metadata: &metadata.DatabaseSchemaMetadata{
		Name: name,
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: columns}},
			Views:  views,
		}},
	}}
}

func mysqlMigration(t *testing.T, from, to *analysis.Database) (string, error) {
	t.Helper()
	return Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{From: from, To: to})
}

func TestEngineMetadataMigrationMySQLRefusesAChangedUnloadableObject(t *testing.T) {
	id := []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}
	_, err := mysqlMigration(t,
		mysqlSnapshot("d", id, &metadata.ViewMetadata{Name: "v", Definition: "select ((broken from"}),
		mysqlSnapshot("d", id, &metadata.ViewMetadata{Name: "v", Definition: "select ((still broken from"}))
	if err == nil || !strings.Contains(err.Error(), "view:v") {
		t.Fatalf("changed view error = %v", err)
	}
	_, err = mysqlMigration(t, mysqlSnapshot("d", id), mysqlSnapshot("d", []*metadata.ColumnMetadata{{Name: "id", Type: "int", Default: "((("}}))
	if err == nil || !strings.Contains(err.Error(), "table:t") {
		t.Fatalf("incomplete table error = %v", err)
	}
}

func TestEngineMetadataMigrationMySQLAcrossDatabaseNames(t *testing.T) {
	id := []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}
	view := func(db string) *metadata.ViewMetadata {
		return &metadata.ViewMetadata{Name: "v", Definition: "select `" + db + "`.`t`.`id` AS `id` from `" + db + "`.`t`"}
	}
	if sql, err := mysqlMigration(t, mysqlSnapshot("prod", id, view("prod")), mysqlSnapshot("staging", id, view("staging"))); sql != "" || err != nil {
		t.Fatalf("same schema = %q, %v", sql, err)
	}
	sql, err := mysqlMigration(t, mysqlSnapshot("prod", id), mysqlSnapshot("staging", id, view("staging")))
	if err != nil || strings.Contains(sql, "staging") || strings.Contains(sql, "bbcatalog") || !strings.Contains(sql, "VIEW `v` AS select `t`.`id` AS `id` from `t`") {
		t.Fatalf("new view = %q, %v", sql, err)
	}
}

func TestEngineMetadataMigrationMySQLIgnoresStatisticsOfAnIncompleteTable(t *testing.T) {
	snapshot := func(rows int64, extra ...*metadata.ColumnMetadata) *analysis.Database {
		d := mysqlSnapshot("d", append([]*metadata.ColumnMetadata{{Name: "id", Type: "int"}}, extra...))
		d.Metadata.Schemas[0].Tables = append(d.Metadata.Schemas[0].Tables, &metadata.TableMetadata{
			Name:     "unparsed",
			Columns:  []*metadata.ColumnMetadata{{Name: "a", Type: "int", Default: "((("}},
			RowCount: rows, DataSize: rows * 16,
		})
		return d
	}
	sql, err := mysqlMigration(t, snapshot(10), snapshot(20, &metadata.ColumnMetadata{Name: "b", Type: "int", Nullable: true, Default: "NULL"}))
	if err != nil || sql != "ALTER TABLE `t` ADD COLUMN `b` int DEFAULT NULL" {
		t.Fatalf("migration = %q, %v", sql, err)
	}
}

func TestEngineMetadataMigrationMySQLRefusesCaseOnlyIdentities(t *testing.T) {
	snapshot := func(names ...string) *analysis.Database {
		var tables []*metadata.TableMetadata
		for _, name := range names {
			tables = append(tables, &metadata.TableMetadata{Name: name, Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}})
		}
		return &analysis.Database{Metadata: &metadata.DatabaseSchemaMetadata{Name: "d", Schemas: []*metadata.SchemaMetadata{{Tables: tables}}}, ObjectCaseSensitive: true}
	}
	if _, err := mysqlMigration(t, snapshot("Foo"), snapshot("foo")); err == nil || !strings.Contains(err.Error(), "case-only rename") {
		t.Fatalf("case-only rename error = %v", err)
	}
	if _, err := mysqlMigration(t, snapshot("Foo", "foo"), snapshot("Foo")); err == nil || !strings.Contains(err.Error(), "cannot distinguish") {
		t.Fatalf("collision error = %v", err)
	}
}

func TestEngineExtractMetadataMySQL(t *testing.T) {
	ctx := context.Background()
	if meta, err := Engine.ExtractMetadata(ctx, "   \n"); err != nil || len(meta.GetSchemas()) != 0 {
		t.Fatalf("empty script = %v, %v", meta, err)
	}
	meta, err := Engine.ExtractMetadata(ctx, "CREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`));")
	if err != nil || meta.GetName() != "tmp" || len(meta.GetSchemas()) != 1 || len(meta.GetSchemas()[0].GetTables()) != 1 {
		t.Fatalf("extracted = %v, %v", meta, err)
	}
	if _, err := Engine.ExtractMetadata(ctx, "CREATE TABLE ("); err == nil || !strings.HasPrefix(err.Error(), "failed to parse schema DDL: ") {
		t.Fatalf("syntax error = %v", err)
	}
	if _, err := Engine.ExtractMetadata(ctx, "CREATE TABLE t (id int); CREATE TABLE t (id int);"); err == nil || !strings.HasPrefix(err.Error(), "failed to execute schema DDL: ") {
		t.Fatalf("statement error = %v", err)
	}
}
