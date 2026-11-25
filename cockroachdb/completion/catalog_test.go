package completion

import (
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func TestCatalogResolvesThroughSearchPath(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "t"}, {Name: "only_public"}}},
			{Name: "app", Tables: []*metadatapb.TableMetadata{{Name: "t"}, {Name: "Mixed"}}},
		},
	}
	cat := newCatalog(meta, "app")
	require.Equal(t, []string{"app", "public"}, cat.searchPath)
	require.Same(t, cat.schema("app").relations[0], cat.relation("", "t"))
	require.Same(t, cat.schema("public").relations[1], cat.relation("", "only_public"))
	require.Same(t, cat.schema("public").relations[0], cat.relation("public", "t"))
	require.Same(t, cat.schema("app").relations[1], cat.relation("app", "mixed"))
	require.Nil(t, cat.relation("missing", "t"))

	require.Equal(t, []string{"public"}, newCatalog(meta, "").searchPath)
	var nilCat *catalog
	require.Nil(t, nilCat.relation("", "t"))
}

func TestCatalogColumns(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "a"},
					{Name: "rowid", IsInvisible: true},
				},
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{
				{Name: "mv", Definition: "SELECT a AS x, t.a, * FROM t"},
				{Name: "broken", Definition: "SELEC"},
			},
		}},
	}
	cat := newCatalog(meta, "public")
	require.Equal(t, []string{"a"}, cat.relation("", "t").visibleColumns())
	var mv []string
	for _, c := range cat.relation("", "mv").columns {
		mv = append(mv, c.name)
	}
	require.Equal(t, []string{"x", "a", "a"}, mv)
	require.Empty(t, cat.relation("", "broken").columns)
}
