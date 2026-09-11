package catalog

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestLoadMetadata(t *testing.T) {
	c := New()
	c.LoadMetadata(&metadata.DatabaseSchemaMetadata{
		Name: "dynamo",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "Orders"}, {Name: "users"}},
		}},
	})
	for _, name := range []string{"Orders", "users"} {
		if !c.HasTable(name) {
			t.Errorf("table %s not loaded; tables %v", name, c.Tables())
		}
	}
	if got := len(c.Tables()); got != 2 {
		t.Errorf("loaded %d tables, want 2: %v", got, c.Tables())
	}
}
