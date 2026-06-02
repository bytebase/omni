package catalog

// Table is a Trino table: the leaf of the
// catalog -> schema -> table -> column hierarchy. It holds an ordered column
// list (ordinal position matters for SELECT * expansion and positional
// references) plus a normalized-name index for O(1) lookup.
//
// All name arguments must already be normalized (see the package doc on
// Normalize).
type Table struct {
	// Name is the normalized table name.
	Name string
	// schema back-references the owning schema (nil only for a detached
	// table constructed directly in a test).
	schema *Schema
	// columns are the table's columns in definition (ordinal) order.
	columns []*Column
	// colByName maps a normalized column name to its index in columns.
	colByName map[string]int
}

// Schema returns the schema that owns this table, or nil if the table was
// constructed detached.
func (t *Table) Schema() *Schema { return t.schema }

// addColumn appends a column, recording its index by its (already-normalized)
// name. A duplicate name keeps the first column's position but updates the
// index to point at the latest definition; callers loading a metadata snapshot
// should not supply duplicates.
func (t *Table) addColumn(col *Column) {
	if col == nil {
		return
	}
	t.colByName[col.Name] = len(t.columns)
	t.columns = append(t.columns, col)
}

// Columns returns the table's columns in ordinal order. The returned slice
// shares backing storage with the table; callers must not mutate it.
func (t *Table) Columns() []*Column { return t.columns }

// GetColumn returns the named (normalized) column, or nil if the table has no
// such column.
func (t *Table) GetColumn(name string) *Column {
	idx, ok := t.colByName[name]
	if !ok {
		return nil
	}
	return t.columns[idx]
}

// HasColumn reports whether the table has a column with the given (normalized)
// name.
func (t *Table) HasColumn(name string) bool {
	_, ok := t.colByName[name]
	return ok
}

// ColumnNames returns the table's column names in ordinal order (the order
// completion offers them and SELECT * expands them). Names are normalized.
func (t *Table) ColumnNames() []string {
	result := make([]string, len(t.columns))
	for i, col := range t.columns {
		result[i] = col.Name
	}
	return result
}

// View is a Trino view (including materialized views). For completion and
// lineage purposes a view exposes a column list just like a table; this
// catalog does not retain the view's defining query.
type View struct {
	// Name is the normalized view name.
	Name string
	// schema back-references the owning schema.
	schema *Schema
	// columns are the view's output columns in order.
	columns []*Column
}

// Schema returns the schema that owns this view, or nil if detached.
func (v *View) Schema() *Schema { return v.schema }

// Columns returns the view's output columns in order. The returned slice
// shares backing storage; callers must not mutate it.
func (v *View) Columns() []*Column { return v.columns }

// ColumnNames returns the view's output column names in order (normalized).
func (v *View) ColumnNames() []string {
	result := make([]string, len(v.columns))
	for i, col := range v.columns {
		result[i] = col.Name
	}
	return result
}

// Column is a single table or view column. It carries the metadata completion
// and lineage need: a normalized name, the Trino type text, and nullability.
type Column struct {
	// Name is the normalized column name.
	Name string
	// Type is the column's Trino type rendered as text (e.g. "varchar(10)",
	// "bigint", "array(integer)"). It is informational for completion and
	// not parsed here.
	Type string
	// Nullable reports whether the column accepts NULL.
	Nullable bool
}

// NewColumn constructs a Column. name must already be normalized (see the
// package doc on Normalize); it is used verbatim as the column's lookup key.
func NewColumn(name, typ string, nullable bool) *Column {
	return &Column{Name: name, Type: typ, Nullable: nullable}
}
