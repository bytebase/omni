package catalog

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestLoadMetadata(t *testing.T) {
	c := New()
	c.LoadMetadata(&metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "users"}, {Name: "orders"}, {Name: ""}},
		}},
	})
	if got, want := c.Collections(), []string{"orders", "users"}; !slices.Equal(got, want) {
		t.Errorf("collections = %v, want %v", got, want)
	}
}
