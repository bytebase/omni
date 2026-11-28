package analysis

import (
	"context"
	"reflect"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func walkThroughPG(t *testing.T, request analysis.WalkThroughRequest) analysis.WalkThroughResult {
	t.Helper()
	if request.Database == nil {
		request.Database = &analysis.Database{
			Name: "test",
			Metadata: &metadata.DatabaseSchemaMetadata{Name: "test", Schemas: []*metadata.SchemaMetadata{
				{Name: "public", Tables: []*metadata.TableMetadata{{Name: "t1", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "integer", Nullable: true, Position: 1}}}}},
				{Name: "s1"},
				{Name: "alice"},
			}},
		}
	}
	return Engine.WalkThrough(context.Background(), request, nil)
}

func pgTables(database *analysis.Database) map[string][]string {
	tables := map[string][]string{}
	for _, schema := range database.Metadata.GetSchemas() {
		for _, table := range schema.GetTables() {
			tables[schema.GetName()] = append(tables[schema.GetName()], table.GetName())
		}
	}
	return tables
}

func TestEngineWalkThroughMapsSQLStatesByLine(t *testing.T) {
	tests := []struct {
		script string
		code   analysis.WalkThroughCode
	}{
		{"CREATE TABLE t1 (a integer);", analysis.WalkThroughTableExists},
		{"DROP TABLE missing;", analysis.WalkThroughTableNotExists},
		{"ALTER TABLE t1 ADD COLUMN a integer;", analysis.WalkThroughColumnExists},
		{"ALTER TABLE t1 DROP COLUMN nope;", analysis.WalkThroughColumnNotExists},
		{"CREATE SCHEMA s1;", analysis.WalkThroughRelationExists},
		{"CREATE TABLE nosuch.t (a integer);", analysis.WalkThroughSchemaNotExists},
		{"DROP INDEX nope;", analysis.WalkThroughIndexNotExists},
		{"CREATE TABLE t2 (a integer PRIMARY KEY, b integer PRIMARY KEY);", analysis.WalkThroughPrimaryKeyExists},
		{"CREATE SCHEMA pg_mine;", analysis.WalkThroughNameIsKeywordIdentifier},
		{"CREATE TABLE t2 (a nosuchtype);", analysis.WalkThroughIndexNotExists},
		{"ALTER TABLE t1 ALTER COLUMN a TYPE uuid;", analysis.WalkThroughChangeColumnType},
		{"CREATE TABLE t2 (a integer DEFAULT nofunc());", analysis.WalkThroughDDLSimulationFailed},
	}
	for _, tt := range tests {
		got := walkThroughPG(t, analysis.WalkThroughRequest{Script: "SELECT 1;\n" + tt.script})
		if got.Database != nil || got.Finding == nil {
			t.Fatalf("%s: WalkThrough = %+v", tt.script, got)
		}
		f := got.Finding
		if f.Code != tt.code || !f.Error || f.Title != "DDL simulation failed" || f.Content == "" || f.Statement != -1 || f.Start != (analysis.Position{Line: 2}) {
			t.Errorf("%s: finding = %+v, want code %s on line 2", tt.script, f, tt.code)
		}
	}
}

func TestEngineWalkThroughReportsAScriptTheCatalogCannotParse(t *testing.T) {
	got := walkThroughPG(t, analysis.WalkThroughRequest{Script: "CREATE TABLE t2 ("})
	if got.Finding == nil || got.Finding.Code != analysis.WalkThroughDDLSimulationFailed || got.Finding.Title != "DDL simulation failed" || got.Finding.Start != (analysis.Position{}) {
		t.Fatalf("finding = %+v", got.Finding)
	}
}

func TestEngineWalkThroughReturnsTheSchemaOnlyWhenItChanges(t *testing.T) {
	for _, script := range []string{"", "SELECT 1;", "INSERT INTO t1 VALUES (1);"} {
		if got := walkThroughPG(t, analysis.WalkThroughRequest{Script: script}); got.Database != nil || got.Finding != nil {
			t.Fatalf("%q: WalkThrough = %+v", script, got)
		}
	}
	got := walkThroughPG(t, analysis.WalkThroughRequest{Script: "CREATE TABLE t2 (a integer);"})
	if got.Finding != nil || got.Database == nil || got.Database.Name != "test" || got.Database.ObjectCaseSensitive {
		t.Fatalf("WalkThrough = %+v", got)
	}
	if tables := pgTables(got.Database); !reflect.DeepEqual(tables["public"], []string{"t1", "t2"}) {
		t.Fatalf("tables = %v", tables)
	}
}

func TestEngineWalkThroughResolvesBySearchPathAndSessionUser(t *testing.T) {
	tests := []struct {
		name    string
		request analysis.WalkThroughRequest
		schema  string
	}{
		{"configured search path", analysis.WalkThroughRequest{SearchPath: []string{"s1", "public"}}, "s1"},
		{"session user's schema", analysis.WalkThroughRequest{SessionUser: "alice", SearchPath: []string{"$user", "public"}}, "alice"},
		{"set in the script", analysis.WalkThroughRequest{Script: "SET search_path TO s1;\n"}, "s1"},
	}
	for _, tt := range tests {
		tt.request.Script += "CREATE TABLE t9 (a integer);"
		got := walkThroughPG(t, tt.request)
		if got.Finding != nil || got.Database == nil {
			t.Fatalf("%s: WalkThrough = %+v", tt.name, got)
		}
		if tables := pgTables(got.Database); !reflect.DeepEqual(tables[tt.schema], []string{"t9"}) {
			t.Errorf("%s: tables = %v, want t9 in %s", tt.name, tables, tt.schema)
		}
	}
}
