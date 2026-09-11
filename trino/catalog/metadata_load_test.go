package catalog

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestLoadMetadata(t *testing.T) {
	c := New()
	c.LoadMetadata(&metadata.DatabaseSchemaMetadata{
		Name: "Hive",
		Schemas: []*metadata.SchemaMetadata{{
			Name: "Sales",
			Tables: []*metadata.TableMetadata{{
				Name:    "Orders",
				Columns: []*metadata.ColumnMetadata{{Name: "ID", Type: "bigint"}, {Name: "note", Type: "varchar", Nullable: true}},
			}},
			Views: []*metadata.ViewMetadata{{
				Name:       "Recent",
				Definition: "SELECT id FROM orders",
				Columns:    []*metadata.ColumnMetadata{{Name: "id", Type: "bigint"}},
			}},
		}},
	})

	database := c.GetCatalog("Hive")
	if database == nil {
		t.Fatalf("catalog not loaded; catalogs %v", c.Catalogs())
	}
	schema := database.GetSchema("Sales")
	if schema == nil {
		t.Fatalf("schema not loaded; schemas %v", database.Schemas())
	}
	table := schema.GetTable("Orders")
	if table == nil || len(table.Columns()) != 2 {
		t.Fatalf("table = %+v, want Orders with two columns", table)
	}
	if id := table.GetColumn("ID"); id == nil || id.Type != "bigint" || id.Nullable {
		t.Errorf("ID = %+v, want a non-null bigint", id)
	}
	if note := table.GetColumn("note"); note == nil || !note.Nullable {
		t.Errorf("note = %+v, want a nullable column", note)
	}
	// A statement spelling the name unquoted resolves it folded, so it must not
	// find the mixed-case one.
	if schema.GetTable(Normalize("Orders")) != nil {
		t.Error("the folded name resolves the table the connector reports as Orders")
	}
	view := schema.GetView("Recent")
	if view == nil || view.Definition != "SELECT id FROM orders" || len(view.Columns()) != 1 {
		t.Errorf("view = %+v, want Recent with its definition and one column", view)
	}
}
