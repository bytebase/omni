// Package catalog provides metadata context for PartiQL auto-completion
// and semantic analysis.
//
// PartiQL targets AWS DynamoDB (schema-on-read), so the catalog is
// minimal: table names are the primary metadata. Future extensions may
// add known attribute paths, type constraints, or index metadata.
package catalog

import "sort"

// Catalog stores table names for a PartiQL database context.
type Catalog struct {
	tables map[string]struct{}
}

// New creates an empty Catalog.
func New() *Catalog {
	return &Catalog{tables: make(map[string]struct{})}
}

// AddTable registers a table name. Duplicates are ignored.
func (c *Catalog) AddTable(name string) {
	c.tables[name] = struct{}{}
}

// Tables returns all registered table names in sorted order.
func (c *Catalog) Tables() []string {
	result := make([]string, 0, len(c.tables))
	for name := range c.tables {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// HasTable reports whether the named table is registered.
func (c *Catalog) HasTable(name string) bool {
	_, ok := c.tables[name]
	return ok
}
