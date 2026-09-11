package catalog

import "github.com/bytebase/omni/metadata"

// LoadMetadata adds a Bytebase schema snapshot of one Trino catalog to c, as
// the catalog named meta.Name with its schemas, tables, and views. Names are
// normalized as Trino resolves them.
func (c *Catalog) LoadMetadata(meta *metadata.DatabaseSchemaMetadata) {
	database := c.EnsureCatalog(Normalize(meta.GetName()))
	for _, s := range meta.GetSchemas() {
		schema := database.EnsureSchema(Normalize(s.GetName()))
		for _, t := range s.GetTables() {
			schema.AddTable(Normalize(t.GetName()), metadataColumns(t.GetColumns())...)
		}
		for _, v := range s.GetViews() {
			schema.AddView(Normalize(v.GetName()), metadataColumns(v.GetColumns())...).Definition = v.GetDefinition()
		}
	}
}

func metadataColumns(columns []*metadata.ColumnMetadata) []*Column {
	out := make([]*Column, 0, len(columns))
	for _, col := range columns {
		if col.GetName() != "" {
			out = append(out, NewColumn(Normalize(col.GetName()), col.GetType(), col.GetNullable()))
		}
	}
	return out
}
