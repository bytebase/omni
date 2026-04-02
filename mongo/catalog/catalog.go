// Package catalog provides metadata context for MongoDB auto-completion.
package catalog

import "sort"

// Catalog stores collection names for a database context.
type Catalog struct {
	collections map[string]struct{}
}

// New creates an empty Catalog.
func New() *Catalog {
	return &Catalog{collections: make(map[string]struct{})}
}

// AddCollection registers a collection name. Duplicates are ignored.
func (c *Catalog) AddCollection(name string) {
	c.collections[name] = struct{}{}
}

// Collections returns all registered collection names in sorted order.
func (c *Catalog) Collections() []string {
	result := make([]string, 0, len(c.collections))
	for name := range c.collections {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
