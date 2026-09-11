package catalog

import "github.com/bytebase/omni/metadata"

// LoadMetadata adds a Bytebase schema snapshot of one Trino catalog to c, as
// the catalog named meta.Name with its schemas, tables, and views. A snapshot
// names an object as the connector reports it, not as SQL spells it, so the
// names go in as they are: a statement resolves a mixed-case one in quotes, the
// way Trino does.
func (c *Catalog) LoadMetadata(meta *metadata.DatabaseSchemaMetadata) {
	database := c.EnsureCatalog(meta.GetName())
	for _, s := range meta.GetSchemas() {
		schema := database.EnsureSchema(s.GetName())
		for _, t := range s.GetTables() {
			schema.AddTable(t.GetName(), metadataColumns(t.GetColumns())...)
		}
		for _, v := range s.GetViews() {
			schema.AddView(v.GetName(), metadataColumns(v.GetColumns())...).Definition = v.GetDefinition()
		}
	}
}

func metadataColumns(columns []*metadata.ColumnMetadata) []*Column {
	out := make([]*Column, 0, len(columns))
	for _, col := range columns {
		if col.GetName() != "" {
			out = append(out, NewColumn(col.GetName(), col.GetType(), col.GetNullable()))
		}
	}
	return out
}
