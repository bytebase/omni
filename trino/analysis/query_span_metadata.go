package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/catalog"
	"github.com/bytebase/omni/trino/parser"
)

// defaultTrinoSchema is the schema assumed for a table reference that omits the
// schema qualifier. The legacy plugin used "public" for the data-access-control
// resource format, and the query-span fixtures encode catalog.public.table; we
// keep that so the resource keys are unchanged across the cutover.
const defaultTrinoSchema = "public"

// parseTrinoSpan is parser.Parse behind a variable so a test can prove a
// statement is parsed once per analysis.
var parseTrinoSpan = parser.Parse

// resolver analyses a single Trino statement and produces a
// base.QuerySpan: the set of physical columns it reads (expanded against
// catalog metadata when available), the predicate columns it filters on, and
// its statement classification.
//
// The extractor delegates table/column lineage to omni's GetQuerySpan,
// then applies bytebase-specific post-processing:
//   - default catalog/schema fill-in for bare table references (Trino's
//     catalog.schema.table → Column Database/Schema/Table),
//   - EXPLAIN unwrap so an EXPLAIN <query> still reports the tables it reads,
//   - expansion of each accessed table to its columns via the metadata getter,
//   - mapping omni's PredicateColumns onto the expanded source columns.
type resolver struct {
	defaultDatabase     string
	defaultSchema       string
	ignoreCaseSensitive bool

	opts ResolveOptions

	// metaCache memoises database metadata lookups within a single span.
	metaCache map[string]*metadata.DatabaseSchemaMetadata
}

// newResolver creates a new Trino query span extractor.
func newResolver(defaultDatabase, defaultSchema string, opts ResolveOptions) *resolver {
	return &resolver{
		defaultDatabase:     defaultDatabase,
		defaultSchema:       defaultSchema,
		opts:                opts,
		ignoreCaseSensitive: opts.IgnoreCaseSensitive,
		metaCache:           make(map[string]*metadata.DatabaseSchemaMetadata),
	}
}

// getQuerySpan extracts the query span for a single Trino statement.
func (q *resolver) getQuerySpan(statement string) (*ResolvedQuerySpan, error) {
	// Split into top-level statements; query-span analysis expects exactly one
	// non-empty statement (matching the legacy behaviour).
	stmts, err := splitSpanStatements(statement)
	if err != nil {
		return nil, err
	}
	nonEmpty := 0
	var single string
	for _, s := range stmts {
		if s.Empty || strings.TrimSpace(s.Text) == "" {
			continue
		}
		nonEmpty++
		single = s.Text
	}
	if nonEmpty == 0 {
		return &ResolvedQuerySpan{
			SourceColumns: ColumnSet{},
			Results:       []ResolvedQuerySpanResult{},
		}, nil
	}
	if nonEmpty > 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", nonEmpty)
	}

	// Resolve the statement type up front so non-SELECT statements short-circuit
	// the way the legacy plugin did. EXPLAIN ANALYZE classifies as Select.
	file, err := parseTrinoSpan(single)
	if err != nil {
		return nil, parser.FirstError(err)
	}
	if file == nil || len(file.Stmts) == 0 {
		return nil, fmt.Errorf("failed to parse Trino statement: %s", single)
	}
	node := file.Stmts[0]
	queryType, _ := spanQueryType(node)
	// Promote a user SELECT/EXPLAIN that references a system schema, matching the
	// legacy containsSystemSchema heuristic.
	switch queryType {
	case Select, Explain:
		if spanContainsSystemSchema(single) {
			queryType = SelectInfoSchema
		}
	default:
	}

	// For statements that don't read columns into a result set (DML, DDL), the
	// legacy plugin returned a basic span: the accessed tables as source
	// columns, no per-column results. CREATE VIEW / EXPLAIN over a query still
	// report the underlying tables.
	readsResults := queryType == Select ||
		queryType == Explain ||
		queryType == SelectInfoSchema

	// Analyze the query the already-parsed statement defines. EXPLAIN and
	// CREATE VIEW/TABLE … AS SELECT carry their query in the AST, so unwrapping
	// it costs no second parse. The catalog scan still reads the whole
	// statement text: it only looks for catalog names the metadata must load.
	query := spanQueryNode(node)
	var rawSpan *QuerySpan
	if query != nil {
		rawSpan, err = getQuerySpanForQuery(query, queryType, newViewState(q.buildSpanCatalog(single)))
		if err != nil {
			return nil, err
		}
	} else {
		rawSpan = &QuerySpan{Type: queryType}
	}

	tables := q.toTableResources(rawSpan)

	if !readsResults {
		// DML/DDL/unknown: report accessed tables as source columns (table-level,
		// no column expansion), no results — mirroring the legacy basic span.
		return &ResolvedQuerySpan{
			Type:          queryType,
			SourceColumns: tables,
			Results:       []ResolvedQuerySpanResult{},
		}, nil
	}

	// Resolve relation aliases (e.g. "u" in "FROM users u") back to physical
	// table names so aliased column references match the expanded source
	// columns.
	aliasMap := q.buildAliasMap(rawSpan)

	// Expand each accessed table to its columns via metadata. If metadata is
	// unavailable (e.g. tests without a getter), fall back to the table-level
	// resource so callers still see what was read. orderedColumns preserves
	// FROM-table-then-metadata column order for positional SELECT * masking.
	//
	// Trino system pseudo-catalogs (system.*) and information_schema are
	// detected and skipped per resolved table inside expandTablesToColumns (not
	// via a query-text substring scan), so a real maskable table is still
	// expanded even when the same statement mentions "system." in a literal.
	orderedColumns, fullSourceColumns, notFound, err := q.expandTablesToColumns(rawSpan)
	if err != nil {
		return nil, err
	}

	results := q.buildResults(rawSpan, fullSourceColumns, orderedColumns, aliasMap)
	predicateColumns := q.mapPredicateColumns(rawSpan, fullSourceColumns, aliasMap)

	span := &ResolvedQuerySpan{
		Type:             queryType,
		SourceColumns:    fullSourceColumns,
		PredicateColumns: predicateColumns,
		Results:          results,
	}
	if notFound != nil {
		span.NotFoundError = notFound
	}
	return span, nil
}

// spanQueryNode returns the query a statement defines, for the statement kinds
// whose reads live in an inner query. EXPLAIN and CREATE … AS SELECT keep the
// query in the parsed AST, so analyzing it parses nothing again.
func spanQueryNode(node ast.Node) *parser.Query {
	switch n := node.(type) {
	case *parser.QueryStmt:
		return n.Query
	case *parser.ExplainStmt:
		if n.Statement != nil {
			return spanQueryNode(n.Statement)
		}
	case *parser.CreateViewStmt:
		return n.Query
	case *parser.CreateTableAsStmt:
		return n.Query
	case *parser.CreateMaterializedViewStmt:
		return n.Query
	}
	return nil
}

// buildSpanCatalog constructs the omni Trino catalog the query-span analysis
// resolves against: views carry their defining query (so lineage through a view
// reaches the underlying base-table columns — masking config only attaches to
// table columns), and tables carry their column lists (so omni expands
// SELECT * to the exact projection). The session context is the extractor's
// default catalog/schema (with the same "public" fallback as
// resolveCatalogSchema, so view resolution and resource keys agree).
//
// Like the completion catalog, every catalog name is registered cheaply but
// full metadata is loaded only for the catalogs the statement needs (the
// default one and any named in the statement). Metadata fetches go through the
// extractor's cache, so the later expandTablesToColumns lookups reuse them.
// Returns nil — omni then behaves exactly catalog-less — when metadata is
// unavailable (e.g. tests without a getter).
func (q *resolver) buildSpanCatalog(statement string) *catalog.Catalog {
	if q.opts.GetDatabaseMetadata == nil || q.opts.ListDatabases == nil {
		return nil
	}
	names, err := q.opts.ListDatabases()
	if err != nil || len(names) == 0 {
		return nil
	}

	defaultDB := catalog.Normalize(q.defaultDatabase)
	cat := catalog.New()
	for _, dbName := range names {
		cat.EnsureCatalog(catalog.Normalize(dbName))
	}

	// Load catalogs transitively: the statement names some catalogs; the view
	// definitions loaded with them may name FURTHER catalogs (a view in the
	// default catalog defined over another catalog's table), which omni needs
	// loaded to resolve lineage through the view. Each round scans the text
	// corpus (statement + every loaded view definition) and loads any
	// newly-referenced catalog; it terminates because each catalog loads at
	// most once.
	corpus := strings.ToLower(statement)
	loaded := make(map[string]bool, len(names))
	for {
		progressed := false
		for _, dbName := range names {
			norm := catalog.Normalize(dbName)
			if loaded[norm] || !catalogNeeded(norm, defaultDB, corpus) {
				continue
			}
			loaded[norm] = true
			progressed = true
			meta, err := q.getDatabaseMetadata(dbName)
			if err != nil || meta == nil {
				continue
			}
			cat.LoadMetadata(meta)
			for _, def := range viewDefinitions(meta) {
				corpus += "\n" + strings.ToLower(def)
			}
		}
		if !progressed {
			break
		}
	}

	if q.defaultDatabase != "" {
		cat.SetCurrentCatalog(defaultDB)
	}
	schema := q.defaultSchema
	if schema == "" {
		schema = defaultTrinoSchema
	}
	cat.SetCurrentSchema(catalog.Normalize(schema))
	return cat
}

// toTableResources converts omni AccessTables into table-level
// Column keys (Column empty), applying default catalog/schema.
func (q *resolver) toTableResources(span *QuerySpan) ColumnSet {
	out := ColumnSet{}
	if span == nil {
		return out
	}
	for _, t := range span.AccessTables {
		db, schema := q.resolveCatalogSchema(t.Catalog, t.Schema)
		out[Column{
			Database: db,
			Schema:   schema,
			Table:    t.Table,
		}] = true
	}
	return out
}

// resolveCatalogSchema fills in the default database (catalog) and schema for a
// table reference that omitted them, mirroring Trino's session-qualified
// resolution. The schema falls back to defaultTrinoSchema ("public") when no
// default schema is configured, matching the legacy resource-key format.
func (q *resolver) resolveCatalogSchema(catalog, schema string) (string, string) {
	db := catalog
	if db == "" {
		db = q.defaultDatabase
	}
	if schema == "" {
		schema = q.defaultSchema
		if schema == "" {
			schema = defaultTrinoSchema
		}
	}
	return db, schema
}

// expandTablesToColumns expands each accessed table into one resource per column
// from the table's metadata. It returns the columns both as an ordered slice
// (AccessTables/FROM order, then metadata column order) and as a set. The
// ordered slice drives positional SELECT * masking, where the order must match
// the executed result's column order; the set is used for membership tests.
//
// A table whose metadata can't be resolved is kept as a single table-level
// resource so the read is still recorded. The first ResourceNotFoundError
// encountered is returned (non-fatally) so the caller can attach it to the span,
// matching the legacy listener behaviour.
func (q *resolver) expandTablesToColumns(span *QuerySpan) ([]Column, ColumnSet, *ResourceNotFoundError, error) {
	ordered := make([]Column, 0)
	set := make(ColumnSet)
	var firstNotFound *ResourceNotFoundError
	// add records a column both in the membership set (deduped, used for lineage
	// matching) and in the ordered slice (NOT deduped, so a self-join such as
	// "FROM users u1 JOIN users u2" contributes each instance's columns and the
	// positional SELECT * masker stays aligned with Trino's 2N output columns).
	add := func(res Column) {
		set[res] = true
		ordered = append(ordered, res)
	}
	if span == nil {
		return ordered, set, nil, nil
	}
	for _, t := range span.AccessTables {
		db, schema := q.resolveCatalogSchema(t.Catalog, t.Schema)
		if isSystemSchemaTable(db, schema) {
			// Trino's "system" catalog (system.runtime/metadata/jdbc.*) and the
			// per-catalog information_schema are pseudo-objects, not
			// Bytebase-tracked databases. Resolving them via the metadata getter
			// would yield a ResourceNotFoundError that sql_service turns into a
			// hard "failed to mask data" rejection of an otherwise-successful
			// result. System metadata is never masked, so keep the table-level
			// resource and skip the lookup. Detection keys on the resolved
			// catalog/schema (not a query-text substring), so a literal
			// containing "system." cannot suppress masking of a real table.
			//
			// Best-effort boundary: a system table contributes one placeholder
			// here but expands to several columns in a real `SELECT *`. For the
			// rare `SELECT * FROM <system> JOIN <real>` shape the positional
			// masker can therefore misalign the trailing real columns. Fully
			// fixing it needs system-table column metadata (not available); such
			// queries are already classified SelectInfoSchema. Tracked as a
			// follow-up rather than blocking the cutover.
			add(Column{Database: db, Schema: schema, Table: t.Table})
			continue
		}
		columns, err := q.tableColumns(db, schema, t.Table)
		if err != nil {
			var notFound *ResourceNotFoundError
			if errors.As(err, &notFound) {
				if firstNotFound == nil {
					firstNotFound = notFound
				}
				// A missing resource is reported on the span; the table-level
				// resource still records the read.
				add(Column{Database: db, Schema: schema, Table: t.Table})
				continue
			}
			// Any other metadata failure is not a missing object; surfacing it
			// prevents a successful span that silently under-reports lineage.
			return nil, nil, nil, err
		}
		if len(columns) == 0 {
			add(Column{Database: db, Schema: schema, Table: t.Table})
			continue
		}
		for _, col := range columns {
			add(Column{Database: db, Schema: schema, Table: t.Table, Column: col})
		}
	}
	return ordered, set, firstNotFound, nil
}

// isSystemSchemaTable reports whether a resolved table reference targets a Trino
// system pseudo-object: the built-in "system" catalog (system.runtime.*,
// system.metadata.*, system.jdbc.*) or the per-catalog "information_schema".
// These are not Bytebase-tracked databases and their metadata is never masked.
func isSystemSchemaTable(catalog, schema string) bool {
	return strings.EqualFold(catalog, "system") || strings.EqualFold(schema, "information_schema")
}

// buildResults maps omni's per-output-column results onto ResolvedQuerySpanResult.
// Each result's SourceColumns is the subset of the expanded source columns that
// share the result's column name (best-effort; omni does not resolve a select
// item back to its owning relation, so aliasMap translates a relation alias on
// the reference back to its physical table first).
//
// For SELECT *, the expansion uses orderedColumns (FROM-table then metadata
// column order) rather than ranging the source-column map: the masker applies
// per-result maskers positionally against the executed result's column order, so
// a nondeterministic order here could mask the wrong column and leak data.
func (q *resolver) buildResults(span *QuerySpan, fullSourceColumns ColumnSet, orderedColumns []Column, aliasMap map[string][]string) []ResolvedQuerySpanResult {
	if span == nil {
		return []ResolvedQuerySpanResult{}
	}
	results := make([]ResolvedQuerySpanResult, 0, len(span.Results))
	for _, r := range span.Results {
		// Unqualified star "*": expand to every source column, in result order.
		//
		// Best-effort boundary: omni records a star as a single "*" result over
		// the flat AccessTables and does not expose the resolved output relation,
		// so a star over a CTE or derived table that projects a subset/reordering
		// (e.g. WITH c AS (SELECT email FROM users) SELECT * FROM c) is
		// indistinguishable here from SELECT * over the base table and expands to
		// the base table's columns. This matches the prior ANTLR plugin (which
		// produced no columns for such stars) and is tracked as a follow-up; a
		// correct fix needs omni's analysis to expose the star's output columns.
		if r.Star && r.Name == "*" {
			results = appendStarColumns(results, orderedColumns, nil)
			continue
		}
		// Qualified star "<rel>.*": omni names the result "<rel>.*" with a single
		// source ref that is the relation itself. Without this branch it would
		// fall through to the column-name match below, find no column literally
		// named "<rel>.*", and leave SourceColumns empty — the masker would then
		// treat those (possibly sensitive) columns as constants and return them
		// unmasked. Resolve the ref to the specific relation(s) it names (by
		// database/schema/table, so same-named tables in different schemas are
		// distinguished) and expand only those columns, preserving order. The
		// detection inspects the source-ref shape, not just the display name, so a
		// column merely aliased to a name ending in ".*" is not mistaken for a
		// star.
		if r.Star && isQualifiedStar(r) {
			targets := q.qualifiedStarTargets(r.SourceColumns[0], span)
			results = appendStarColumns(results, orderedColumns, targets)
			continue
		}
		sourceColumns := ColumnSet{}
		for _, ref := range r.SourceColumns {
			addMatchingColumns(sourceColumns, fullSourceColumns, ref.Column, ref.Table, aliasMap)
		}
		results = append(results, ResolvedQuerySpanResult{
			Name:          r.Name,
			SourceColumns: sourceColumns,
			// Plainness keys on the MAPPED physical set, not omni's raw ref
			// count: the catalog-aware resolver additively restates a plain
			// column as written + catalog-qualified refs, which dedupe back to
			// the one physical column here. Known drift: an expression over one
			// column repeated (phone || phone) also maps to a single physical
			// column and reads as plain — omni's resolver dedups refs, so the
			// distinction is unrecoverable here; inert for Trino, which is not
			// in EngineSupportQuerySpanPlainField.
			IsPlainField: len(sourceColumns) == 1,
		})
	}
	return results
}

// appendStarColumns appends one plain-field result per column in orderedColumns,
// in order (the positional masker maps each result to the executed result column
// at the same index). When targets is non-nil — a qualified star "<rel>.*" — only
// columns whose (database, schema, table) key is in targets are appended; a nil
// targets (unqualified "*") appends every column.
func appendStarColumns(results []ResolvedQuerySpanResult, orderedColumns []Column, targets ColumnSet) []ResolvedQuerySpanResult {
	for _, col := range orderedColumns {
		if targets != nil {
			if !targets[Column{Database: col.Database, Schema: col.Schema, Table: col.Table}] {
				continue
			}
		}
		results = append(results, ResolvedQuerySpanResult{
			Name:          col.Column,
			SourceColumns: ColumnSet{col: true},
			IsPlainField:  true,
		})
	}
	return results
}

// qualifiedStarResult reports whether r is a genuine qualified star "<rel>.*"
// (returning the relation qualifier), as opposed to an ordinary column that
// happens to be aliased to a name ending in ".*". omni represents a qualified
// star as a result whose Name is "<rel>.*" and whose single source ref is the
// relation name itself (Column == <rel>, no catalog/schema/table qualifier); a
// real column aliased to e.g. "u.*" instead carries the underlying column name
// in its source ref, so it is handled as a normal column and not star-expanded
// (which would misalign the positional masker and could leak).
func isQualifiedStar(r ColumnInfo) bool {
	qualifier, ok := qualifiedStarQualifier(r.Name)
	if !ok {
		return false
	}
	// omni represents a qualified star "<...>.<rel>.*" as a single source ref
	// whose Column is the relation name <rel> (any schema/catalog path written
	// before it lands in the ref's Table/Schema fields). A real column aliased to
	// a name ending in ".*" instead carries the underlying column name in Column,
	// which will not equal the qualifier, so it is handled as a normal column and
	// not star-expanded (avoiding a positional-masker misalignment/leak).
	//
	// Matching only on Column == qualifier (not on empty table/schema) is
	// deliberate: it keeps genuine table/schema/catalog-qualified stars like
	// public.users.* working. The sole residual ambiguity is the contrived case
	// of a column literally named like the qualifier and aliased to
	// "<qualifier>.*", which omni renders identically to a real star.
	if len(r.SourceColumns) != 1 {
		return false
	}
	return strings.EqualFold(r.SourceColumns[0].Column, qualifier)
}

// qualifiedStarTargets resolves the relation(s) named by a qualified-star source
// ref into a set of (database, schema, table) keys, so the star expands to
// exactly that relation's columns — distinguishing same-named tables in
// different schemas or catalogs. omni encodes the written qualifier in the ref
// with the parts shifted right: Column is the relation/alias, Table is the schema
// part, and Schema is the catalog part.
func (q *resolver) qualifiedStarTargets(ref ColumnRef, span *QuerySpan) ColumnSet {
	targets := make(ColumnSet)
	// Alias-qualified star (u.*): collect every base table that carries the
	// alias. A reused/shadowed alias yields several, which over-includes (the
	// extra columns are ignored positionally) rather than risking a leak.
	aliasMatched := false
	for _, t := range span.AccessTables {
		if t.Alias != "" && strings.EqualFold(t.Alias, ref.Column) {
			db, schema := q.resolveCatalogSchema(t.Catalog, t.Schema)
			targets[Column{Database: db, Schema: schema, Table: t.Table}] = true
			aliasMatched = true
		}
	}
	if aliasMatched {
		return targets
	}
	// Table-name-qualified star (users.* / public.users.* / cat.sch.users.*):
	// ref.Column is the table, ref.Table the schema part, ref.Schema the catalog.
	// Only target a base table actually present in AccessTables, so a qualifier
	// naming a CTE/derived relation (which omni omits from AccessTables) does not
	// accidentally target an unrelated physical table of the same name. A
	// qualified star over a CTE/derived relation thus yields no columns here —
	// the same best-effort boundary as an unqualified star over such a relation
	// (omni does not expose the resolved output projection; tracked as a
	// follow-up).
	wantDB, wantSchema := q.resolveCatalogSchema(ref.Schema, ref.Table)
	for _, t := range span.AccessTables {
		db, schema := q.resolveCatalogSchema(t.Catalog, t.Schema)
		if strings.EqualFold(db, wantDB) && strings.EqualFold(schema, wantSchema) && strings.EqualFold(t.Table, ref.Column) {
			targets[Column{Database: db, Schema: schema, Table: t.Table}] = true
		}
	}
	return targets
}

// qualifiedStarQualifier returns the relation qualifier of a qualified-star
// result name like "u.*" or "catalog.schema.t.*" (the rightmost component before
// the trailing ".*"), and whether name is a qualified star. Column matching keys
// on the table/alias, so only the last component is returned.
func qualifiedStarQualifier(name string) (string, bool) {
	if !strings.HasSuffix(name, ".*") {
		return "", false
	}
	qualifier := strings.TrimSuffix(name, ".*")
	if i := strings.LastIndex(qualifier, "."); i >= 0 {
		qualifier = qualifier[i+1:]
	}
	if qualifier == "" {
		return "", false
	}
	return qualifier, true
}

// mapPredicateColumns maps omni's predicate column refs onto the expanded source
// columns (matching on column name, and table when the ref carries one; a
// relation alias on the ref is resolved via aliasMap).
func (*resolver) mapPredicateColumns(span *QuerySpan, fullSourceColumns ColumnSet, aliasMap map[string][]string) ColumnSet {
	out := make(ColumnSet)
	if span == nil {
		return out
	}
	for _, ref := range span.PredicateColumns {
		addMatchingColumns(out, fullSourceColumns, ref.Column, ref.Table, aliasMap)
	}
	return out
}

// addMatchingColumns adds every column in fullSourceColumns whose name matches
// refColumn (case-insensitively, and whose table matches refTable when refTable
// is non-empty) to dst. refTable may be a relation alias; see tableMatches.
func addMatchingColumns(dst, fullSourceColumns ColumnSet, refColumn, refTable string, aliasMap map[string][]string) {
	for sc := range fullSourceColumns {
		if !strings.EqualFold(sc.Column, refColumn) {
			continue
		}
		if refTable == "" || tableMatches(sc.Table, refTable, aliasMap) {
			dst[sc] = true
		}
	}
}

// tableMatches reports whether the physical table scTable is named by refTable,
// either directly (its own name) or as a relation alias for it. Alias resolution
// is additive: the written refTable and every physical table the alias maps to
// are all accepted. omni's AccessTables is a flat, scope-less list, so an alias
// reused or shadowed across subqueries yields several candidates; accepting all
// of them over-includes (conservatively masks more) rather than risk an
// under-match that would leave a sensitive column unmasked.
func tableMatches(scTable, refTable string, aliasMap map[string][]string) bool {
	if strings.EqualFold(scTable, refTable) {
		return true
	}
	for _, phys := range aliasMap[strings.ToLower(refTable)] {
		if strings.EqualFold(scTable, phys) {
			return true
		}
	}
	return false
}

// buildAliasMap maps each relation alias (lower-cased) to the physical table
// names it stands for, so aliased column references resolve back to base tables
// during column matching. Because omni's AccessTables is flat (carries no SQL
// scope), an alias reused in different scopes maps to several tables; all are
// kept so matching can over-include rather than under-match. An alias equal to
// its own table name carries no information and is skipped.
func (*resolver) buildAliasMap(span *QuerySpan) map[string][]string {
	if span == nil {
		return nil
	}
	out := make(map[string][]string)
	for _, t := range span.AccessTables {
		if t.Alias == "" || strings.EqualFold(t.Alias, t.Table) {
			continue
		}
		key := strings.ToLower(t.Alias)
		if !containsFold(out[key], t.Table) {
			out[key] = append(out[key], t.Table)
		}
	}
	return out
}

// containsFold reports whether list contains s under case-insensitive
// comparison.
func containsFold(list []string, s string) bool {
	for _, e := range list {
		if strings.EqualFold(e, s) {
			return true
		}
	}
	return false
}

// tableColumns returns the column names of the given table or view, honouring
// the case-sensitivity flag for object lookup.
func (q *resolver) tableColumns(db, schema, table string) ([]string, error) {
	dbMeta, err := q.getDatabaseMetadata(db)
	if err != nil {
		return nil, err
	}

	schemaMeta := findSpanSchema(dbMeta, schema)
	if schemaMeta == nil && q.ignoreCaseSensitive {
		for _, name := range listSpanSchemaNames(dbMeta) {
			if strings.EqualFold(name, schema) {
				schemaMeta = findSpanSchema(dbMeta, name)
				break
			}
		}
	}
	if schemaMeta == nil {
		return nil, &ResourceNotFoundError{
			Database: &db,
			Schema:   &schema,
		}
	}

	// Table first.
	tableMeta := findSpanTable(schemaMeta, table)
	if tableMeta == nil && q.ignoreCaseSensitive {
		for _, name := range listSpanTableNames(schemaMeta) {
			if strings.EqualFold(name, table) {
				tableMeta = findSpanTable(schemaMeta, name)
				break
			}
		}
	}
	if tableMeta != nil {
		var cols []string
		for _, col := range tableMeta.GetColumns() {
			cols = append(cols, col.Name)
		}
		return cols, nil
	}

	// Then view.
	viewMeta := findSpanView(schemaMeta, table)
	if viewMeta == nil && q.ignoreCaseSensitive {
		for _, name := range listSpanViewNames(schemaMeta) {
			if strings.EqualFold(name, table) {
				viewMeta = findSpanView(schemaMeta, name)
				break
			}
		}
	}
	if viewMeta != nil {
		var cols []string
		for _, col := range viewMeta.GetColumns() {
			cols = append(cols, col.Name)
		}
		return cols, nil
	}

	return nil, &ResourceNotFoundError{
		Database: &db,
		Schema:   &schema,
		Table:    &table,
	}
}

// getDatabaseMetadata fetches (and caches) metadata for the given database.
func (q *resolver) getDatabaseMetadata(database string) (*metadata.DatabaseSchemaMetadata, error) {
	if database == "" {
		database = q.defaultDatabase
	}
	if meta, ok := q.metaCache[database]; ok {
		return meta, nil
	}
	if q.opts.GetDatabaseMetadata == nil {
		return nil, &ResourceNotFoundError{Database: &database}
	}
	meta, err := q.opts.GetDatabaseMetadata(database)
	if err != nil {
		var notFound *ResourceNotFoundError
		if errors.As(err, &notFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get database metadata for database: %s: %w", database, err)
	}
	if meta == nil {
		empty := &metadata.DatabaseSchemaMetadata{}
		q.metaCache[database] = empty
		return empty, nil
	}
	q.metaCache[database] = meta
	return meta, nil
}

// nodeText returns the source substring covered by node n within statement, or
// "" when the node's location is unusable.
// nodeSpan returns the source byte range of an omni statement node. Most omni
// node types expose Span() ast.Loc; *parser.QueryStmt carries its range on a
// plain Loc field. (ast.NodeLoc only knows ast-package nodes — File/Identifier/
// QualifiedName — not the parser-package statement nodes, so it cannot be used
// to unwrap an EXPLAIN's inner statement.)
func nodeSpan(n ast.Node) (ast.Loc, bool) {
	if qs, ok := n.(*parser.QueryStmt); ok {
		return qs.Loc, true
	}
	if sp, ok := n.(interface{ Span() ast.Loc }); ok {
		return sp.Span(), true
	}
	return ast.Loc{}, false
}
