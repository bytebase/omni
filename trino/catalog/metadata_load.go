package catalog

import "github.com/bytebase/omni/metadata"

// LoadMetadata adds a Bytebase schema snapshot of one Trino catalog to c, as
// the catalog named meta.Name with its schemas, tables, and views. Names go in
// folded, as Trino folds an unquoted identifier, so a statement resolves them
// as it spells them: a table the snapshot records as Employees resolves, and
// completes, as employees. Bytebase's Trino completer is built on that, and
// stores no case of its own.
func (c *Catalog) LoadMetadata(meta *metadata.DatabaseSchemaMetadata) {
	database := c.EnsureCatalog(Normalize(meta.GetName()))
	for _, s := range meta.GetSchemas() {
		schema := database.EnsureSchema(Normalize(s.GetName()))
		for _, t := range s.GetTables() {
			if t.GetName() != "" {
				schema.AddTable(Normalize(t.GetName()), metadataColumns(t.GetColumns())...)
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() != "" {
				schema.AddView(Normalize(v.GetName()), metadataColumns(v.GetColumns())...).Definition = v.GetDefinition()
			}
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
