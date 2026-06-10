package analysis

import (
	"github.com/bytebase/omni/trino/catalog"
)

// viewState carries the catalog-aware resolution state for one
// GetQuerySpanWithCatalog call, shared across the recursive analysis of view
// definitions. It supplies:
//   - the catalog itself and the CURRENT resolution context for unqualified
//     names (the session catalog/schema at the top level; a view's own
//     catalog/schema while its definition is being analyzed);
//   - a cycle guard and a memo for view projections, so a view referenced
//     several times is analyzed once and a definition cycle in malformed
//     metadata terminates;
//   - a per-statement collector for the base tables discovered inside view
//     definitions, which the caller appends to that statement's AccessTables.
//
// A nil viewState (or one with a nil catalog) disables all catalog-aware
// resolution: behaviour is byte-for-byte the catalog-less GetQuerySpan.
type viewState struct {
	cat *catalog.Catalog

	// curCatalog/curSchema are the resolution context for unqualified names,
	// mirroring catalog.ResolveRelation's session rules with an injected
	// context (the shared catalog is never mutated).
	curCatalog string
	curSchema  string

	// inProgress guards against view-definition cycles.
	inProgress map[viewKey]bool
	// cycleHit records that a cycle was encountered while computing the
	// current projection; projections computed under an open cycle are partial
	// and are discarded rather than memoized.
	cycleHit bool
	// memo caches each view's resolved projection (nil entry = unresolvable).
	memo map[viewKey]*viewProjection

	// tablesOut collects, for the statement currently being analyzed, the
	// qualified base tables read through resolved views. Swapped per
	// getQuerySpanWithViews invocation so nested view analyses collect into
	// their own statement's span.
	tablesOut *[]TableAccess
}

// viewKey identifies a view by its canonical (normalized) location.
type viewKey struct {
	catalog string
	schema  string
	view    string
}

// viewProjection is a view's resolved output: its columns (with base-column
// lineage, names taken positionally from the view's metadata column list when
// present) and the qualified base tables its definition reads.
type viewProjection struct {
	cols   []outCol
	tables []TableAccess
}

// newViewState builds the resolution state for one top-level statement.
// Returns nil when there is no catalog, which disables catalog resolution.
func newViewState(cat *catalog.Catalog) *viewState {
	if cat == nil {
		return nil
	}
	return &viewState{
		cat:        cat,
		curCatalog: cat.CurrentCatalog(),
		curSchema:  cat.CurrentSchema(),
		inProgress: make(map[viewKey]bool),
		memo:       make(map[viewKey]*viewProjection),
	}
}

// resolveParts maps normalized dotted-name parts to a fully-qualified
// (catalog, schema, object) triple using this state's resolution context —
// catalog.ResolveRelation's session rules with the context injected.
func (vs *viewState) resolveParts(parts []string) (catalogName, schemaName, objectName string, ok bool) {
	switch len(parts) {
	case 1:
		if vs.curCatalog == "" || vs.curSchema == "" {
			return "", "", "", false
		}
		return vs.curCatalog, vs.curSchema, parts[0], true
	case 2:
		if vs.curCatalog == "" {
			return "", "", "", false
		}
		return vs.curCatalog, parts[0], parts[1], true
	case 3:
		return parts[0], parts[1], parts[2], true
	default:
		return "", "", "", false
	}
}

// lookupRelation resolves normalized name parts against the catalog. Exactly
// one of view/table is non-nil when found, with the same precedence as
// catalog.ResolveRelation: an existing table shadows a same-named view.
func (vs *viewState) lookupRelation(parts []string) (key viewKey, view *catalog.View, table *catalog.Table, found bool) {
	if vs == nil || vs.cat == nil {
		return viewKey{}, nil, nil, false
	}
	catalogName, schemaName, objectName, ok := vs.resolveParts(parts)
	if !ok {
		return viewKey{}, nil, nil, false
	}
	db := vs.cat.GetCatalog(catalogName)
	if db == nil {
		return viewKey{}, nil, nil, false
	}
	sch := db.GetSchema(schemaName)
	if sch == nil {
		return viewKey{}, nil, nil, false
	}
	key = viewKey{catalog: db.Name, schema: sch.Name, view: objectName}
	if t := sch.GetTable(objectName); t != nil {
		return key, nil, t, true
	}
	if v := sch.GetView(objectName); v != nil {
		return key, v, nil, true
	}
	return viewKey{}, nil, nil, false
}

// tableCols renders a catalog table's columns as resolved output columns, each
// sourced by its own fully-qualified base column. starExpansion uses these to
// expand a star over a catalog-known base table to its exact projection.
func tableCols(key viewKey, t *catalog.Table) []outCol {
	cols := t.Columns()
	out := make([]outCol, 0, len(cols))
	for _, c := range cols {
		out = append(out, outCol{
			name: c.Name,
			sources: []ColumnRef{{
				Catalog: key.catalog,
				Schema:  key.schema,
				Table:   key.view,
				Column:  c.Name,
			}},
		})
	}
	return out
}

// projectionFor returns the view's resolved projection, computing and memoizing
// it on first use. The definition is analyzed recursively (so views over views
// and definitions using CTEs/derived tables resolve) in the view's own
// catalog/schema context. Returns nil — leaving the view opaque, exactly the
// pre-catalog behaviour — for an empty/unanalyzable definition or a definition
// cycle.
func (vs *viewState) projectionFor(key viewKey, v *catalog.View) *viewProjection {
	if p, ok := vs.memo[key]; ok {
		return p
	}
	if vs.inProgress[key] {
		// Definition cycle: record the hit so every projection computed while
		// this cycle was open is discarded rather than memoized partial.
		vs.cycleHit = true
		return nil
	}
	if v.Definition == "" {
		vs.memo[key] = nil
		return nil
	}

	vs.inProgress[key] = true
	savedCatalog, savedSchema := vs.curCatalog, vs.curSchema
	savedCycleHit := vs.cycleHit
	vs.cycleHit = false
	vs.curCatalog, vs.curSchema = key.catalog, key.schema
	span, err := getQuerySpanWithViews(v.Definition, vs)
	vs.curCatalog, vs.curSchema = savedCatalog, savedSchema
	delete(vs.inProgress, key)
	tainted := vs.cycleHit
	vs.cycleHit = savedCycleHit || tainted

	if tainted {
		// This projection was computed with a cyclic dependency unresolved
		// (treated as opaque mid-flight) — a partial result. Do not memoize and
		// report the view unresolvable, so every view in the cycle stays fully
		// opaque (the consumer's metadata expansion applies).
		return nil
	}
	if err != nil || span == nil || len(span.Results) == 0 {
		vs.memo[key] = nil
		return nil
	}

	cols := make([]outCol, 0, len(span.Results))
	for _, r := range span.Results {
		cols = append(cols, outCol{
			name:    r.Name,
			sources: r.SourceColumns,
			// A star the inner analysis could not expand (over an unknown base
			// table, a coalescing join, …) leaves the projection's width
			// unknown; mark it opaque so star expansion through this view
			// bails while named-column lookups still work.
			opaque: r.Name == "*" || isQualifiedStarName(r.Name),
		})
	}
	// The view's metadata column list names the outputs the OUTER query
	// references; apply it positionally over the definition's analyzed columns
	// (an explicit view column list — CREATE VIEW v (ph) AS SELECT phone … —
	// renames them). A COUNT MISMATCH between the metadata column list and the
	// analyzed definition means one of the two is stale; the true output shape
	// is then unknown, so the view is unresolvable — expanding or renaming on
	// either width could misalign the consumer's positional masker.
	if names := v.ColumnNames(); len(names) > 0 {
		if len(names) != len(cols) {
			vs.memo[key] = nil
			return nil
		}
		for i, name := range names {
			if name != "" {
				cols[i].name = name
			}
		}
	}

	// The relations the definition reads — its base tables, and any
	// intermediate views (a query through this view does access them) —
	// qualified with the view's own context when the definition wrote them
	// unqualified, so the consumer resolves them in the right namespace.
	// Tables appended by nested view resolution arrive already qualified.
	tables := make([]TableAccess, 0, len(span.AccessTables))
	for _, t := range span.AccessTables {
		if t.Catalog == "" {
			t.Catalog = key.catalog
		}
		if t.Schema == "" {
			t.Schema = key.schema
		}
		tables = append(tables, t)
	}

	p := &viewProjection{cols: cols, tables: tables}
	vs.memo[key] = p
	return p
}

// collectTables records base tables read through a resolved view for the
// statement currently being analyzed.
func (vs *viewState) collectTables(tables []TableAccess) {
	if vs == nil || vs.tablesOut == nil {
		return
	}
	*vs.tablesOut = append(*vs.tablesOut, tables...)
}

// isQualifiedStarName reports whether a result name has the unexpanded
// qualified-star form "<rel>.*".
func isQualifiedStarName(name string) bool {
	return len(name) > 2 && name[len(name)-2:] == ".*"
}

// appendViewTables merges the base tables collected from resolved views into a
// span's AccessTables, deduplicating on (catalog, schema, table, alias)
// against both the existing entries and each other.
func appendViewTables(span *QuerySpan, tables []TableAccess) {
	if len(tables) == 0 {
		return
	}
	seen := make(map[tableKey]bool, len(span.AccessTables)+len(tables))
	for _, t := range span.AccessTables {
		seen[tableKey{Catalog: t.Catalog, Schema: t.Schema, Table: t.Table, Alias: t.Alias}] = true
	}
	for _, t := range tables {
		k := tableKey{Catalog: t.Catalog, Schema: t.Schema, Table: t.Table, Alias: t.Alias}
		if seen[k] {
			continue
		}
		seen[k] = true
		span.AccessTables = append(span.AccessTables, t)
	}
}

