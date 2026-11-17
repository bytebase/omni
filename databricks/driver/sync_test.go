package driver

import (
	"testing"

	catalogapi "github.com/databricks/databricks-sdk-go/service/catalog"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

func TestAppendSchemaTablesKeepsTheSyncedKinds(t *testing.T) {
	catalogMap := databricksCatalogMap{
		"cat": databricksSchemaMap{
			"sch": {},
		},
	}
	tables := []catalogapi.TableInfo{
		{
			CatalogName:    "cat",
			SchemaName:     "sch",
			Name:           "z_view",
			TableType:      catalogapi.TableTypeView,
			ViewDefinition: "SELECT a, b FROM cat.sch.t",
			Comment:        "a view",
			Columns: []catalogapi.ColumnInfo{
				{Name: "b", TypeName: catalogapi.ColumnTypeNameInt, Position: 1},
				{Name: "a", TypeName: catalogapi.ColumnTypeNameString, Position: 0},
			},
		},
		{
			CatalogName: "cat",
			SchemaName:  "sch",
			Name:        "m_table",
			TableType:   catalogapi.TableTypeManaged,
			Comment:     "a table",
			Columns: []catalogapi.ColumnInfo{
				{Name: "second", TypeName: catalogapi.ColumnTypeNameLong, Position: 1, Nullable: true, Comment: "second column"},
				{Name: "first", TypeName: catalogapi.ColumnTypeNameString, Position: 0},
			},
		},
		{
			CatalogName: "cat",
			SchemaName:  "sch",
			Name:        "external",
			TableType:   catalogapi.TableTypeExternal,
		},
		{
			CatalogName:    "cat",
			SchemaName:     "sch",
			Name:           "mv",
			TableType:      catalogapi.TableTypeMaterializedView,
			ViewDefinition: "SELECT 1",
			Comment:        "a materialized view",
		},
		{
			CatalogName: "cat",
			SchemaName:  "sch",
			Name:        "stream",
			TableType:   catalogapi.TableTypeStreamingTable,
		},
	}
	if err := appendSchemaTables(catalogMap, tables); err != nil {
		t.Fatalf("appendSchemaTables: %v", err)
	}

	schemas := convertToStorepbSchemas(catalogMap["cat"])
	if len(schemas) != 1 || schemas[0].Name != "sch" {
		t.Fatalf("Schemas = %+v, want one schema sch", schemas)
	}
	schema := schemas[0]
	if len(schema.Tables) != 1 || schema.Tables[0].Name != "m_table" {
		t.Errorf("Tables = %+v, want m_table only", schema.Tables)
	}
	if schema.Tables[0].Comment != "a table" {
		t.Errorf("table comment = %q", schema.Tables[0].Comment)
	}
	if len(schema.Tables[0].Columns) != 2 {
		t.Fatalf("Columns = %+v, want two", schema.Tables[0].Columns)
	}
	if schema.Tables[0].Columns[0].Name != "first" || schema.Tables[0].Columns[1].Name != "second" {
		t.Errorf("Columns are not in position order: %+v", schema.Tables[0].Columns)
	}
	if schema.Tables[0].Columns[1].Type != "LONG" || !schema.Tables[0].Columns[1].Nullable || schema.Tables[0].Columns[1].Comment != "second column" {
		t.Errorf("second column = %+v", schema.Tables[0].Columns[1])
	}

	if len(schema.Views) != 1 || schema.Views[0].Name != "z_view" {
		t.Fatalf("Views = %+v, want z_view", schema.Views)
	}
	if len(schema.Views[0].DependencyColumns) != 2 || schema.Views[0].DependencyColumns[0].Column != "a" {
		t.Errorf("DependencyColumns = %+v, want a and b in order", schema.Views[0].DependencyColumns)
	}
	if len(schema.ExternalTables) != 1 || schema.ExternalTables[0].Name != "external" {
		t.Errorf("ExternalTables = %+v", schema.ExternalTables)
	}
	if len(schema.MaterializedViews) != 1 || schema.MaterializedViews[0].Name != "mv" {
		t.Errorf("MaterializedViews = %+v", schema.MaterializedViews)
	}
}

func TestAppendSchemaTablesNeedsAnInitializedSchema(t *testing.T) {
	err := appendSchemaTables(databricksCatalogMap{}, []catalogapi.TableInfo{{
		CatalogName: "missing",
		SchemaName:  "missing",
		Name:        "t",
		TableType:   catalogapi.TableTypeManaged,
	}})
	if err == nil {
		t.Fatal("appendSchemaTables accepted an uninitialized catalog")
	}
}

func TestVersionFromResult(t *testing.T) {
	for _, tc := range []struct {
		name    string
		result  *driver.Result
		want    string
		wantErr bool
	}{
		{
			name:   "the two-word answer",
			result: &driver.Result{Rows: [][]driver.Value{{{Kind: driver.String, String: "4.2.0 1.0.0"}}}},
			want:   "4.2.0",
		},
		{
			name:    "no rows",
			result:  &driver.Result{},
			wantErr: true,
		},
		{
			name:    "two rows",
			result:  &driver.Result{Rows: [][]driver.Value{{{Kind: driver.String, String: "1"}}, {{Kind: driver.String, String: "2"}}}},
			wantErr: true,
		},
		{
			name:    "not a string",
			result:  &driver.Result{Rows: [][]driver.Value{{{Kind: driver.Int64, Int64: 1}}}},
			wantErr: true,
		},
		{
			name:    "empty",
			result:  &driver.Result{Rows: [][]driver.Value{{{Kind: driver.String, String: "  "}}}},
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := versionFromResult(tc.result)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("versionFromResult accepted %+v", tc.result)
				}
				return
			}
			if err != nil {
				t.Fatalf("versionFromResult: %v", err)
			}
			if got != tc.want {
				t.Errorf("versionFromResult = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestConvertToStorepbSchemasSorts pins the order a sync states, so two syncs
// of one catalog are equal and a tape has one order.
func TestConvertToStorepbSchemasSorts(t *testing.T) {
	schemaMap := databricksSchemaMap{
		"z_schema": {
			{name: "b", typeName: catalogapi.TableTypeManaged, table: &metadatapb.TableMetadata{Name: "b"}},
			{name: "a", typeName: catalogapi.TableTypeManaged, table: &metadatapb.TableMetadata{Name: "a"}},
		},
		"a_schema": {
			{name: "only", typeName: catalogapi.TableTypeManaged, table: &metadatapb.TableMetadata{Name: "only"}},
		},
	}
	schemas := convertToStorepbSchemas(schemaMap)
	if len(schemas) != 2 || schemas[0].Name != "a_schema" || schemas[1].Name != "z_schema" {
		t.Fatalf("Schemas = %+v, want a_schema then z_schema", schemas)
	}
	if len(schemas[1].Tables) != 2 || schemas[1].Tables[0].Name != "a" || schemas[1].Tables[1].Name != "b" {
		t.Errorf("Tables = %+v, want a then b", schemas[1].Tables)
	}
}
