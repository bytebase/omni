package catalog

import "sort"

// Schema is one Trino schema: the third level of the
// catalog -> schema -> table/view hierarchy. A schema owns tables and views,
// each keyed by normalized name. Trino keeps tables and views in the same
// name space (you cannot have a table and a view with the same name in one
// schema); this catalog stores them in separate maps for direct lookup, and
// Relations() merges them into the single de-duplicated relation candidate set
// completion offers after FROM / JOIN.
//
// All name arguments must already be normalized (see the package doc on
// Normalize).
type Schema struct {
	// Name is the normalized schema name.
	Name string
	// tables maps a normalized table name to its Table.
	tables map[string]*Table
	// views maps a normalized view name to its View.
	views map[string]*View
}

// newSchema creates an empty schema with the given (already-normalized) name.
func newSchema(name string) *Schema {
	return &Schema{
		Name:   name,
		tables: make(map[string]*Table),
		views:  make(map[string]*View),
	}
}

// AddTable registers a table with the given (normalized) name and columns,
// returning the created Table. If a table with the same name already exists it
// is replaced (last definition wins, matching how a fresh metadata snapshot is
// loaded).
func (s *Schema) AddTable(name string, columns ...*Column) *Table {
	t := &Table{Name: name, schema: s, colByName: make(map[string]int)}
	for _, col := range columns {
		t.addColumn(col)
	}
	s.tables[name] = t
	return t
}

// GetTable returns the named (normalized) table, or nil if it is not
// registered.
func (s *Schema) GetTable(name string) *Table {
	return s.tables[name]
}

// HasTable reports whether the named (normalized) table is registered.
func (s *Schema) HasTable(name string) bool {
	_, ok := s.tables[name]
	return ok
}

// Tables returns the normalized names of all tables in this schema in sorted
// order.
func (s *Schema) Tables() []string {
	result := make([]string, 0, len(s.tables))
	for name := range s.tables {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// AddView registers a view with the given (normalized) name and columns,
// returning the created View. If a view with the same name already exists it
// is replaced.
func (s *Schema) AddView(name string, columns ...*Column) *View {
	v := &View{Name: name, schema: s}
	for _, col := range columns {
		if col != nil {
			v.columns = append(v.columns, col)
		}
	}
	s.views[name] = v
	return v
}

// GetView returns the named (normalized) view, or nil if it is not registered.
func (s *Schema) GetView(name string) *View {
	return s.views[name]
}

// HasView reports whether the named (normalized) view is registered.
func (s *Schema) HasView(name string) bool {
	_, ok := s.views[name]
	return ok
}

// Views returns the normalized names of all views in this schema in sorted
// order.
func (s *Schema) Views() []string {
	result := make([]string, 0, len(s.views))
	for name := range s.views {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// Relations returns the normalized names of all tables and views in this
// schema in sorted order, de-duplicated (a name present as both a table and a
// view appears once). This is the set of relation candidates completion offers
// after FROM / JOIN.
func (s *Schema) Relations() []string {
	seen := make(map[string]struct{}, len(s.tables)+len(s.views))
	for name := range s.tables {
		seen[name] = struct{}{}
	}
	for name := range s.views {
		seen[name] = struct{}{}
	}
	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
