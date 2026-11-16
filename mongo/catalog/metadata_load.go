package catalog

import "github.com/bytebase/omni/metadata"

// LoadMetadata adds the collections and views of a Bytebase schema snapshot to
// c. Sync reports each collection as a table and each view as a view; both are
// query targets, so completion offers both.
func (c *Catalog) LoadMetadata(meta *metadata.DatabaseSchemaMetadata) {
	for _, s := range meta.GetSchemas() {
		for _, t := range s.GetTables() {
			if t.GetName() != "" {
				c.AddCollection(t.GetName())
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() != "" {
				c.AddCollection(v.GetName())
			}
		}
	}
}
