package catalog

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

// The snapshot's names go in folded, so every lookup here spells them the way
// a statement does.
func TestLoadMetadata(t *testing.T) {
	c := New()
	c.LoadMetadata(&metadata.DatabaseSchemaMetadata{
		Name: "Hive",
		Schemas: []*metadata.SchemaMetadata{{
			Name: "Sales",
			Tables: []*metadata.TableMetadata{{
				Name:    "Orders",
				Columns: []*metadata.ColumnMetadata{{Name: "ID", Type: "bigint"}, {Name: "note", Type: "varchar", Nullable: true}},
			}, {Name: ""}},
			Views: []*metadata.ViewMetadata{{
				Name:       "Recent",
				Definition: "SELECT id FROM orders",
				Columns:    []*metadata.ColumnMetadata{{Name: "id", Type: "bigint"}},
			}, {Name: ""}},
		}},
	})

	database := c.GetCatalog(Normalize("Hive"))
	if database == nil {
		t.Fatalf("catalog not loaded; catalogs %v", c.Catalogs())
	}
	schema := database.GetSchema(Normalize("Sales"))
	if schema == nil {
		t.Fatalf("schema not loaded; schemas %v", database.Schemas())
	}
	table := schema.GetTable(Normalize("Orders"))
	if table == nil || len(table.Columns()) != 2 {
		t.Fatalf("table = %+v, want Orders with two columns", table)
	}
	if id := table.GetColumn(Normalize("ID")); id == nil || id.Type != "bigint" || id.Nullable {
		t.Errorf("ID = %+v, want a non-null bigint", id)
	}
	if note := table.GetColumn(Normalize("note")); note == nil || !note.Nullable {
		t.Errorf("note = %+v, want a nullable column", note)
	}
	view := schema.GetView(Normalize("Recent"))
	if view == nil || view.Definition != "SELECT id FROM orders" || len(view.Columns()) != 1 {
		t.Errorf("view = %+v, want Recent with its definition and one column", view)
	}
	// An unnamed relation would complete as a blank candidate.
	if len(schema.Tables()) != 1 || len(schema.Views()) != 1 {
		t.Errorf("tables %v, views %v; want only the named ones", schema.Tables(), schema.Views())
	}
}
