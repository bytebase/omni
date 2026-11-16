package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/googlesql/parser"
	"github.com/bytebase/omni/metadata"
)

// querySpanResolver resolves a single GoogleSQL statement against the caller's
// catalog metadata: it delegates the structural lineage to the span walker,
// then applies the catalog-aware post-processing the masking consumer needs —
// the dialect name model onto ColumnResource, the user/system table handling,
// metadata star expansion, and the legacy result naming.
//
// The walker resolves relation projections (CTE/derived/recursive bodies,
// set-operation merges with star markers, SELECT * EXCEPT/REPLACE, JOIN USING
// coalescing, UNNEST element lineage) before handing over; the one piece left to
// this metadata-aware layer is enumerating a PHYSICAL table's columns
// (StarSegment.BaseTable), which the walker cannot do catalog-free.
type querySpanResolver struct {
	opts ResolveOptions

	// metaCache memoises database metadata lookups within a single span.
	metaCache map[string]*metadata.DatabaseSchemaMetadata
}

func newQuerySpanResolver(opts ResolveOptions) *querySpanResolver {
	return &querySpanResolver{
		opts:      opts,
		metaCache: make(map[string]*metadata.DatabaseSchemaMetadata),
	}
}

func (q *querySpanResolver) getQuerySpan(statement string) (*ResolvedQuerySpan, error) {
	// Query-span analysis expects exactly one non-empty statement.
	single, count := q.singleNonEmptyStatement(statement)
	switch {
	case count == 0:
		return &ResolvedQuerySpan{
			SourceColumns: ColumnSet{},
			Results:       []ResolvedQuerySpanResult{},
		}, nil
	case count > 1:
		return nil, fmt.Errorf("expecting only one statement to get query span, but got %d", count)
	}

	rawSpan, err := GetQuerySpan(single, q.opts.Dialect)
	if err != nil {
		return nil, err
	}

	// Reject a statement that reads from both user and system
	// (INFORMATION_SCHEMA / SPANNER_SYS) tables.
	allSystems, mixed := classifyAccess(rawSpan)
	if mixed {
		return nil, ErrMixUserSystemTables
	}
	if allSystems && q.opts.SystemOnlyEmptySpan {
		// A query reading exclusively from system tables early-returns an EMPTY
		// SelectInfoSchema span — no source columns, no results, no metadata
		// access. System metadata is never masked, and resolving these
		// pseudo-tables would only risk a spurious ResourceNotFoundError.
		return &ResolvedQuerySpan{
			Type:          SelectInfoSchema,
			SourceColumns: ColumnSet{},
			Results:       []ResolvedQuerySpanResult{},
		}, nil
	}

	queryType := rawSpan.Type
	if queryType == Unknown && q.opts.SetStatementIsSelect && isSetStatement(single) {
		// A SET statement classifies as Select; the span walker's classifier
		// leaves it Unknown.
		queryType = Select
	}
	// A query reading exclusively from system tables is SelectInfoSchema. The
	// walker's classifier already promotes this, but recompute from the access
	// table system flags too so the result is robust if the two ever disagree
	// (the legacy extractors derived allSystems from the access tables).
	if queryType == Select && allSystems {
		queryType = SelectInfoSchema
	}

	// Top-level SourceColumns is the set of accessed tables at the TABLE level
	// (Column empty) — the legacy spans report accessed tables, not
	// column-expanded source columns. Uses the table name AS WRITTEN in the
	// query, with the dialect defaults filled in.
	accessTables := q.toTableResources(rawSpan)

	// Only a read-only query produces per-column results; DML/DDL/Explain/
	// Unknown keep the accessed tables and no per-column lineage.
	// SelectInfoSchema is a read so it is treated like Select here.
	if queryType != Select && queryType != SelectInfoSchema {
		return &ResolvedQuerySpan{
			Type:          queryType,
			SourceColumns: accessTables,
			Results:       []ResolvedQuerySpanResult{},
		}, nil
	}

	// Resolve relation aliases (e.g. "g" in "FROM galleries g") back to physical
	// table names so aliased column references match the expanded source columns.
	aliasMap := q.buildAliasMap(rawSpan)

	// Expand each accessed table to its columns via metadata: the membership set
	// used for lineage matching. The per-column resources use the
	// METADATA-canonical table name, which is why result lineage carries
	// metadata casing while the table-level SourceColumns above carry the
	// written casing.
	fullSourceColumns, notFound, err := q.expandTablesToColumns(rawSpan)
	if err != nil {
		return nil, err
	}

	results := q.expandColumnInfos(rawSpan.Results, fullSourceColumns, aliasMap)

	span := &ResolvedQuerySpan{
		Type:             queryType,
		SourceColumns:    accessTables,
		PredicateColumns: q.mapPredicateColumns(rawSpan, fullSourceColumns),
		Results:          results,
	}
	if notFound != nil {
		span.NotFoundError = notFound
	}
	return span, nil
}

// singleNonEmptyStatement splits the input and returns its single non-empty
// statement's text together with the non-empty count.
func (q *querySpanResolver) singleNonEmptyStatement(statement string) (string, int) {
	single, count := "", 0
	for _, seg := range parser.Split(statement) {
		if seg.Empty() || strings.TrimSpace(seg.Text) == "" {
			continue
		}
		single = seg.Text
		count++
	}
	return single, count
}

// isSetStatement reports whether the statement parses to a single SET statement.
func isSetStatement(statement string) bool {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return false
	}
	_, ok := file.Stmts[0].(*ast.SetStmt)
	return ok
}

// classifyAccess reports (allSystems, mixed) over the span's access tables,
// using the walker's per-dialect IsSystem flag (BigQuery: INFORMATION_SCHEMA;
// Spanner: INFORMATION_SCHEMA and SPANNER_SYS). It reproduces the legacy
// isMixedQuery: mixed when both a system and a user table are present;
// allSystems when at least one system table is present and no user table is.
func classifyAccess(span *QuerySpan) (allSystems, mixed bool) {
	if span == nil {
		return false, false
	}
	hasSystem, hasUser := false, false
	for _, t := range span.AccessTables {
		if t.IsSystem {
			hasSystem = true
		} else {
			hasUser = true
		}
	}
	if hasSystem && hasUser {
		return false, true
	}
	return hasSystem && !hasUser, false
}

// resourceDatabase resolves the Database field for a resource: the written
// database/dataset with default fill-in (BigQuery), or always the default
// database (Spanner — the legacy resolver silently ignored the db part of a
// 3-part db.schema.table reference).
func (q *querySpanResolver) resourceDatabase(written string) string {
	if q.opts.IgnoreWrittenDatabase || written == "" {
		return q.opts.DefaultDatabase
	}
	return written
}

// toTableResources converts the walker's AccessTables into table-level
// ColumnResource keys (Column empty), applying the dialect name model. The
// table and schema names are kept AS WRITTEN; the project/catalog qualifier is
// dropped (the legacy listeners never recorded it).
func (q *querySpanResolver) toTableResources(span *QuerySpan) ColumnSet {
	out := ColumnSet{}
	if span == nil {
		return out
	}
	// CTE references are not physical tables, but the legacy access-table
	// listeners recorded every FROM table path — CTE references included — so
	// union them into the table-level SourceColumns.
	for _, group := range [][]TableAccess{span.AccessTables, span.CTEReferences} {
		for _, t := range group {
			out[ColumnResource{
				Database: q.resourceDatabase(t.Database),
				Schema:   t.Schema,
				Table:    t.Table,
			}] = true
		}
	}
	return out
}

// expandTablesToColumns expands each accessed table into one resource per column
// from the table's metadata, returning the membership set used for lineage
// matching.
//
// A table whose metadata does not EXIST is kept as a single table-level
// resource so the read is still recorded, and the first ResourceNotFoundError
// encountered is returned (non-fatally) for the caller to attach to the span —
// the masking layer rejects spans carrying NotFoundError. Any OTHER metadata
// error (a store outage, a canceled context) is returned FATALLY: falling
// through would leave result columns with silently-empty lineage and no
// NotFoundError, and the fail-open positional masker would return sensitive
// data unmasked after a transient infrastructure failure.
// Each per-column resource uses the METADATA-canonical table/column names, so
// result lineage carries metadata casing.
func (q *querySpanResolver) expandTablesToColumns(span *QuerySpan) (ColumnSet, *ResourceNotFoundError, error) {
	set := make(ColumnSet)
	var firstNotFound *ResourceNotFoundError
	if span == nil {
		return set, nil, nil
	}
	for _, t := range span.AccessTables {
		db := q.resourceDatabase(t.Database)
		if t.IsSystem {
			// System pseudo-objects are not Bytebase-tracked tables; their
			// metadata is never masked. Resolving them via the metadata getter
			// would yield a ResourceNotFoundError that turns into a hard
			// "failed to mask" rejection of an otherwise successful result, so
			// keep the table-level resource and skip the lookup. (A system-only
			// query already early-returned under SystemOnlyEmptySpan, and a
			// mixed user/system query was rejected.)
			set[ColumnResource{Database: db, Schema: t.Schema, Table: t.Table}] = true
			continue
		}
		canonicalTable, columns, err := q.tableColumns(t.Database, t.Schema, t.Table)
		if err != nil {
			var notFound *ResourceNotFoundError
			if !errors.As(err, &notFound) {
				return nil, nil, err
			}
			if firstNotFound == nil {
				firstNotFound = notFound
			}
			// Keep the table-level resource for the missing object.
			set[ColumnResource{Database: db, Schema: t.Schema, Table: t.Table}] = true
			continue
		}
		if len(columns) == 0 {
			set[ColumnResource{Database: db, Schema: t.Schema, Table: canonicalTable}] = true
			continue
		}
		for _, col := range columns {
			set[ColumnResource{Database: db, Schema: t.Schema, Table: canonicalTable, Column: col}] = true
		}
	}
	return set, firstNotFound, nil
}

// expandColumnInfos expands a list of ColumnInfo (one query body's resolved
// projection) into ResolvedQuerySpanResult. It is called for the top-level
// results and RECURSIVELY by buildSetOpMergeResults for each arm of a deferred
// set-op merge. Each ColumnInfo is one of: a star item (StarSegments, expanded
// + its EXCEPT/REPLACE applied), a per-position base-star set-op merge
// (StarMerge), a whole deferred set-op merge (SetOpMerge, both arms expanded
// then position-merged), or a concrete output column.
func (q *querySpanResolver) expandColumnInfos(infos []ColumnInfo, fullSourceColumns ColumnSet, aliasMap map[string][]string) []ResolvedQuerySpanResult {
	results := make([]ResolvedQuerySpanResult, 0, len(infos))
	for _, r := range infos {
		switch {
		case r.StarSegments != nil:
			// A `*` / `rel.*` / CTE-or-derived star: the walker resolved its
			// ordered expansion into segments (a base table to enumerate, or an
			// already-resolved concrete column). Expand them, then apply the
			// star's EXCEPT/REPLACE modifiers with the legacy name-collision
			// last-wins dedup.
			results = append(results, q.expandStarSegments(r, fullSourceColumns, aliasMap)...)
		case r.SetOpMerge != nil:
			// A whole deferred set-operation merge whose arms carry
			// un-enumerable stars: expand each arm fully (recursively) and merge
			// them, reproducing the legacy "fully resolve each arm, then zip".
			results = append(results, q.buildSetOpMergeResults(r.SetOpMerge, fullSourceColumns, aliasMap)...)
		case r.StarMerge != nil:
			// A set-operation merge position whose other arm is a base-table
			// star: union the concrete arm's lineage with that table's Index-th
			// column.
			results = append(results, q.buildStarMergeResult(r, fullSourceColumns, aliasMap))
		default:
			sourceColumns := ColumnSet{}
			for _, ref := range r.SourceColumns {
				addMatchingColumns(sourceColumns, fullSourceColumns, ref.Column, ref.Table, aliasMap)
			}
			name := resultName(r)
			// A base-table FIELD passthrough (the JOIN ... USING coalesced key)
			// keeps the field's METADATA case under the spanner naming model —
			// the legacy resolver named it after the left PhysicalTable's field,
			// not the written/upper-cased token.
			if q.opts.CanonicalBaseFieldName && r.BaseFieldName && len(r.SourceColumns) > 0 {
				if canon := q.canonicalColumnCase(r.SourceColumns[0], r.Name); canon != "" {
					name = canon
				}
			}
			results = append(results, ResolvedQuerySpanResult{
				Name:          name,
				SourceColumns: sourceColumns,
				// IsPlainField is carried from the walker: a `SELECT *` / `rel.*`
				// expansion column over a base table is plain; an explicit
				// select-list item, a set-op merge, and a join-left column are
				// not (the walker already applied these rules).
				IsPlainField: r.IsPlain,
			})
		}
	}
	return results
}

// buildSetOpMergeResults expands a deferred set-operation merge: each arm is
// expanded fully (a base-table star against metadata, an EXCEPT/REPLACE star
// with its modifiers, a nested merge — recursively via expandColumnInfos), then
// the two expanded column lists are merged. This reproduces the legacy
// extractTableSourceFromQuerySetOperation: output NAME comes from the LEFT arm
// (the first-select-name rule), SourceColumns are the union of both arms at
// that position, and the merged column is never a plain field. The walker emits
// this only when a per-position StarMerge cannot express the merge (a
// base-star UNION base-star, or an EXCEPT/REPLACE star arm in a set operation)
// — the arity is known only here, after metadata expansion.
//
// Width handling mirrors masking safety, not the legacy hard error: legacy
// rejects an unequal-width set operation outright (so the query never masks),
// while this merges up to the LONGER arm so no arm's column is dropped (a
// dropped sensitive column would under-attribute). The shorter arm simply
// contributes nothing past its length. Equal-width arms (the only valid set
// operations) merge exactly.
func (q *querySpanResolver) buildSetOpMergeResults(m *SetOpMergeInfo, fullSourceColumns ColumnSet, aliasMap map[string][]string) []ResolvedQuerySpanResult {
	left := q.expandColumnInfos(m.Left, fullSourceColumns, aliasMap)
	right := q.expandColumnInfos(m.Right, fullSourceColumns, aliasMap)
	if m.ByName {
		// BY NAME / CORRESPONDING aligns the arms by output column NAME, not
		// ordinal — an ordinal merge here mis-attributes whenever the arms'
		// column orders differ, and ignoring MatchColumns emits the wrong arity
		// (BY NAME ON (cols) outputs ONLY the listed columns), shifting the
		// positional masker off every real output column.
		return mergeExpandedByName(left, right, m.MatchColumns)
	}
	n := len(left)
	if len(right) > n {
		n = len(right)
	}
	merged := make([]ResolvedQuerySpanResult, 0, n)
	for i := 0; i < n; i++ {
		sourceColumns := ColumnSet{}
		name := ""
		if i < len(left) {
			name = left[i].Name
			for sc := range left[i].SourceColumns {
				sourceColumns[sc] = true
			}
		}
		if i < len(right) {
			if name == "" {
				name = right[i].Name
			}
			for sc := range right[i].SourceColumns {
				sourceColumns[sc] = true
			}
		}
		merged = append(merged, ResolvedQuerySpanResult{
			Name:          name,
			SourceColumns: sourceColumns,
			IsPlainField:  false,
		})
	}
	return merged
}

// mergeExpandedByName merges two expanded set-operation arms by output column
// NAME (case-insensitive — GoogleSQL column names are), mirroring the resolver's
// mergeProjectionsByName for arms that needed metadata expansion first. With
// matchColumns (BY NAME ON (cols)) the output is EXACTLY those columns in list
// order; otherwise the left arm's names in order plus right-only names appended
// (over-inclusion is masking-safe: a trailing extra never shifts earlier
// positions). Each output column's lineage is the union of BOTH arms'
// same-named columns.
func mergeExpandedByName(left, right []ResolvedQuerySpanResult, matchColumns []string) []ResolvedQuerySpanResult {
	namedSources := func(arm []ResolvedQuerySpanResult, name string) ColumnSet {
		out := ColumnSet{}
		for _, r := range arm {
			if strings.EqualFold(r.Name, name) {
				for sc := range r.SourceColumns {
					out[sc] = true
				}
			}
		}
		return out
	}
	mergedFor := func(name string) ResolvedQuerySpanResult {
		sources := namedSources(left, name)
		for sc := range namedSources(right, name) {
			sources[sc] = true
		}
		return ResolvedQuerySpanResult{Name: name, SourceColumns: sources, IsPlainField: false}
	}
	if len(matchColumns) > 0 {
		out := make([]ResolvedQuerySpanResult, 0, len(matchColumns))
		for _, name := range matchColumns {
			out = append(out, mergedFor(name))
		}
		return out
	}
	out := make([]ResolvedQuerySpanResult, 0, len(left)+len(right))
	for _, l := range left {
		out = append(out, mergedFor(l.Name))
	}
	for _, r := range right {
		if !containsResultName(left, r.Name) {
			out = append(out, mergedFor(r.Name))
		}
	}
	return out
}

// containsResultName reports whether any result in list has the given output
// name under case-insensitive comparison.
func containsResultName(list []ResolvedQuerySpanResult, name string) bool {
	for _, r := range list {
		if strings.EqualFold(r.Name, name) {
			return true
		}
	}
	return false
}

// expandStarSegments expands one star item's resolved segments (in order) into
// per-column results, then applies the star's EXCEPT/REPLACE modifiers. A base
// segment enumerates its physical table's columns via metadata (each a plain
// field per the segment's Plain); a concrete segment is emitted directly with
// its resolved lineage. EXCEPT/REPLACE then act on the expanded column list
// keyed by output name with last-wins on a name collision — the legacy
// starModify semantics (two FROM relations contributing a like-named column
// keep the later one before EXCEPT removes a name).
func (q *querySpanResolver) expandStarSegments(r ColumnInfo, fullSourceColumns ColumnSet, aliasMap map[string][]string) []ResolvedQuerySpanResult {
	var expanded []ResolvedQuerySpanResult
	for _, seg := range r.StarSegments {
		if seg.BaseTable != nil {
			expanded = append(expanded, q.expandBaseTableColumns(*seg.BaseTable, seg.Plain, seg.ExceptColumns)...)
			continue
		}
		// A concrete segment (a CTE/derived/explicit column resolved by the
		// walker).
		sourceColumns := ColumnSet{}
		for _, ref := range seg.Sources {
			addMatchingColumns(sourceColumns, fullSourceColumns, ref.Column, ref.Table, aliasMap)
		}
		// A base-table FIELD passthrough inside a star (the JOIN ... USING
		// coalesced key) keeps the field's METADATA case under the spanner
		// naming model.
		name := seg.Name
		if q.opts.CanonicalBaseFieldName && seg.BaseFieldName && len(seg.Sources) > 0 {
			if canon := q.canonicalColumnCase(seg.Sources[0], name); canon != "" {
				name = canon
			}
		}
		expanded = append(expanded, ResolvedQuerySpanResult{
			Name:          name,
			SourceColumns: sourceColumns,
			IsPlainField:  seg.Plain,
		})
	}
	return applyStarModifiers(expanded, r, fullSourceColumns, aliasMap)
}

// expandBaseTableColumns enumerates a physical base table's columns from
// metadata into per-column plain-field results, each lineage'd to that column.
// A table whose metadata cannot be resolved yields no columns (the access is
// still recorded at the table level by expandTablesToColumns /
// toTableResources, and a ResourceNotFoundError is surfaced there). except
// lists column names to SKIP (case-insensitive): a JOIN ... USING key is
// projected once as a coalesced concrete segment ahead of the side stars, so
// each side's star must not expand it again — re-expanding would shift every
// later position and misalign the positional masker.
func (q *querySpanResolver) expandBaseTableColumns(rel ColumnRef, plain bool, except []string) []ResolvedQuerySpanResult {
	canonicalTable, columns, err := q.tableColumns(rel.Database, rel.Schema, rel.Table)
	if err != nil || len(columns) == 0 {
		return nil
	}
	db := q.resourceDatabase(rel.Database)
	out := make([]ResolvedQuerySpanResult, 0, len(columns))
	for _, col := range columns {
		if containsFold(except, col) {
			continue
		}
		res := ColumnResource{Database: db, Schema: rel.Schema, Table: canonicalTable, Column: col}
		out = append(out, ResolvedQuerySpanResult{
			Name:          col,
			SourceColumns: ColumnSet{res: true},
			IsPlainField:  plain,
		})
	}
	return out
}

// canonicalColumnCase returns the metadata-canonical spelling of column in the
// table named by ref (case-folded match), or "" when the table or column cannot
// be resolved.
func (q *querySpanResolver) canonicalColumnCase(ref ColumnRef, column string) string {
	_, columns, err := q.tableColumns(ref.Database, ref.Schema, ref.Table)
	if err != nil {
		return ""
	}
	for _, col := range columns {
		if strings.EqualFold(col, column) {
			return col
		}
	}
	return ""
}

// applyStarModifiers applies a star item's EXCEPT/REPLACE modifiers to its
// expanded column list, reproducing the legacy starModify: results are keyed by
// output name (case-insensitive), last-wins on a collision; an EXCEPT name
// drops its column entirely; a REPLACE (expr AS name) re-points the named
// column's lineage to the replacement expression's sources and clears its
// plain-field flag. Original order is preserved. When the star has no modifiers
// the expanded list is returned unchanged.
func applyStarModifiers(expanded []ResolvedQuerySpanResult, r ColumnInfo, fullSourceColumns ColumnSet, aliasMap map[string][]string) []ResolvedQuerySpanResult {
	if len(r.StarExcept) == 0 && len(r.StarReplace) == 0 {
		return expanded
	}
	type item struct {
		id    int
		field ResolvedQuerySpanResult
	}
	// Key by output name; a name collision keeps the LAST occurrence (legacy
	// fieldItemMap overwrite), but the surviving entry keeps its FIRST ordinal
	// so the output order is stable.
	order := map[string]int{}
	byName := map[string]item{}
	next := 0
	for _, f := range expanded {
		key := strings.ToLower(f.Name)
		if _, ok := order[key]; !ok {
			order[key] = next
			next++
		}
		byName[key] = item{id: order[key], field: f}
	}
	for _, name := range r.StarExcept {
		delete(byName, strings.ToLower(name))
	}
	for _, rep := range r.StarReplace {
		key := strings.ToLower(rep.Name)
		if _, ok := byName[key]; !ok {
			continue
		}
		set := ColumnSet{}
		for _, ref := range rep.Sources {
			addMatchingColumns(set, fullSourceColumns, ref.Column, ref.Table, aliasMap)
		}
		byName[key] = item{
			id:    byName[key].id,
			field: ResolvedQuerySpanResult{Name: rep.Name, SourceColumns: set, IsPlainField: false},
		}
	}
	items := make([]item, 0, len(byName))
	for _, it := range byName {
		items = append(items, it)
	}
	slices.SortFunc(items, func(a, b item) int { return a.id - b.id })
	out := make([]ResolvedQuerySpanResult, 0, len(items))
	for _, it := range items {
		out = append(out, it.field)
	}
	return out
}

// buildStarMergeResult builds a set-operation merge position whose other arm is
// a base-table star: the concrete arm's lineage (r.SourceColumns) is unioned
// with the star table's Index-th column (expanded from metadata). When the star
// arm is the LEFT one (StarMerge.LeftStar) the output name comes from that
// column (the legacy first-select-name rule); otherwise the concrete arm's name
// is kept. The merged column is never a plain field.
//
// Naming: a concrete-arm output name is uppercased (the legacy expression-name
// rendering, both engines). A LEFT-star-derived name is uppercased under the
// BigQuery model but keeps the star column's metadata case under the spanner
// model (the legacy spanner first-arm rule passed the PhysicalTable field name
// through unchanged).
func (q *querySpanResolver) buildStarMergeResult(r ColumnInfo, fullSourceColumns ColumnSet, aliasMap map[string][]string) ResolvedQuerySpanResult {
	sourceColumns := ColumnSet{}
	for _, ref := range r.SourceColumns {
		addMatchingColumns(sourceColumns, fullSourceColumns, ref.Column, ref.Table, aliasMap)
	}
	name := strings.ToUpper(r.Name)
	// A StarMerge always carries a bare single-base-table star (a
	// USING-coalesced or otherwise multi-segment star arm defers to SetOpMerge
	// instead), so there are no except columns to apply here.
	starCols := q.expandBaseTableColumns(r.StarMerge.Table, false, nil)
	if idx := r.StarMerge.Index; idx >= 0 && idx < len(starCols) {
		for sc := range starCols[idx].SourceColumns {
			sourceColumns[sc] = true
		}
		if r.StarMerge.LeftStar {
			name = starCols[idx].Name
			if q.opts.UppercaseStarMergeName {
				name = strings.ToUpper(name)
			}
		}
	}
	return ResolvedQuerySpanResult{
		Name:          name,
		SourceColumns: sourceColumns,
		IsPlainField:  false,
	}
}

// resultName renders an explicit select-item's output column name to match the
// legacy extractors (extractTableSourceFromSelect): the name is UPPER-CASED,
// and for an expression that has no written name of its own the walker leaves
// Name empty — in that case the legacy code used the name of the single bare
// column the expression's DFS surfaced (e.g. `ID+1` → "ID"), which is the
// walker's first source column. A name-less expression with no column reference
// (a literal, or an opaque scalar subquery) stays "".
func resultName(r ColumnInfo) string {
	name := r.Name
	if name == "" && len(r.SourceColumns) > 0 {
		name = r.SourceColumns[0].Column
	}
	return strings.ToUpper(name)
}

// mapPredicateColumns maps the walker's row-condition column refs (WHERE,
// JOIN ... ON / USING, HAVING) onto the expanded source columns — matching on
// column name, and on the table, database and schema the walk resolved the ref
// to. The resolved resources therefore carry the dialect's name model: the
// default database fill-in, the Spanner schema bucketing, and the
// METADATA-canonical table and column names.
//
// A qualified ref whose relation the expansion cannot attribute — a recursive
// CTE's self-reference, a relation only an outer query body knows — falls back
// to matching the column name alone, because under-attributing a row condition
// leaves a sensitive predicate unchecked while over-attributing only makes the
// check stricter. A map-independent ref that matches no expanded column (a
// referenced column the metadata does not hold, a pseudo-column) is dropped: a
// resource the metadata does not have cannot carry a policy.
func (q *querySpanResolver) mapPredicateColumns(span *QuerySpan, fullSourceColumns ColumnSet) ColumnSet {
	out := ColumnSet{}
	if span == nil {
		return out
	}
	for _, ref := range span.filterColumns {
		if ref.Column == "" {
			continue
		}
		if q.addMatchingPredicateColumns(out, fullSourceColumns, ref, true) {
			continue
		}
		// A qualified ref the expansion cannot attribute (a recursive CTE's
		// self-reference, a relation only an outer body knows) falls back to the
		// column name alone.
		if ref.Table != "" {
			q.addMatchingPredicateColumns(out, fullSourceColumns, ColumnRef{Column: ref.Column}, false)
		}
	}
	return out
}

// addMatchingPredicateColumns adds every expanded column a filter ref names and
// reports whether it found one. A ref the walk resolved to a relation is matched
// on the table — narrowed by the database and schema it carries, so the same
// table name in two datasets or schemas of the statement does not over-match —
// and a star ref (a subquery projecting `*` it could not enumerate) matches every
// column of its relation. With matchTable false only the column name is compared.
func (q *querySpanResolver) addMatchingPredicateColumns(dst, fullSourceColumns ColumnSet, ref ColumnRef, matchTable bool) bool {
	star := ref.Column == filterStarColumn
	found := false
	for sc := range fullSourceColumns {
		if sc.Column == "" {
			continue
		}
		if matchTable && ref.Table != "" {
			if !strings.EqualFold(sc.Table, ref.Table) {
				continue
			}
			if ref.Database != "" && !strings.EqualFold(sc.Database, q.resourceDatabase(ref.Database)) {
				continue
			}
			if ref.Schema != "" && !strings.EqualFold(sc.Schema, ref.Schema) {
				continue
			}
		}
		if !star && !strings.EqualFold(sc.Column, ref.Column) {
			continue
		}
		dst[sc] = true
		found = true
	}
	return found
}

// addMatchingColumns adds every column in fullSourceColumns whose name matches
// refColumn (case-insensitively — GoogleSQL column names are case-insensitive —
// and whose table matches refTable when refTable is non-empty) to dst. refTable
// may be a relation alias; see tableMatches.
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
// are all accepted. The walker's AccessTables is a flat, scope-less list, so an
// alias reused or shadowed across subqueries yields several candidates;
// accepting all of them over-includes (conservatively masks more) rather than
// risk an under-match that would leave a sensitive column unmasked.
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

// buildAliasMap maps each relation alias (lower-cased) to the metadata-canonical
// physical table names it stands for, so aliased column references resolve back
// to the expanded (metadata-cased) source columns during column matching.
// Because the walker's AccessTables is flat (carries no SQL scope), an alias
// reused in different scopes maps to several tables; all are kept so matching
// can over-include rather than under-match. An alias equal to its own table name
// carries no information and is skipped.
func (q *querySpanResolver) buildAliasMap(span *QuerySpan) map[string][]string {
	if span == nil {
		return nil
	}
	out := make(map[string][]string)
	for _, t := range span.AccessTables {
		if t.Alias == "" || strings.EqualFold(t.Alias, t.Table) {
			continue
		}
		table := t.Table
		if !t.IsSystem {
			if canonical, _, err := q.tableColumns(t.Database, t.Schema, t.Table); err == nil && canonical != "" {
				table = canonical
			}
		}
		key := strings.ToLower(t.Alias)
		if !containsFold(out[key], table) {
			out[key] = append(out[key], table)
		}
	}
	return out
}

// containsFold reports whether list contains s under case-insensitive comparison.
func containsFold(list []string, s string) bool {
	return slices.ContainsFunc(list, func(e string) bool { return strings.EqualFold(e, s) })
}

// namedColumns is the table-or-view subset of the metadata model the column
// lookup needs: both expose a canonical name and an ordered column-name list.
type namedColumns struct {
	canonical string
	columns   []string
}

// tableColumns returns the metadata-canonical table name and the column names of
// the given table or view, per the dialect metadata model:
//   - spanner (SchemaScopedMetadata): the schema named by the reference
//     (case-SENSITIVE; "" is the default schema) within the default database;
//   - bigquery: the written dataset resolved as the DATABASE (default fill-in)
//     with its single unnamed schema.
//
// Table names fall back to a case-folded match in both models. Tables shadow
// views on a name collision (the legacy lookup tried tables first).
func (q *querySpanResolver) tableColumns(writtenDatabase, schema, table string) (string, []string, error) {
	lookupDatabase := q.opts.DefaultDatabase
	lookupSchema := schema
	if !q.opts.SchemaScopedMetadata {
		lookupDatabase = q.resourceDatabase(writtenDatabase)
		lookupSchema = ""
	}
	dbMetadata, err := q.getDatabaseMetadata(lookupDatabase)
	if err != nil {
		return "", nil, err
	}

	schemaMeta := findSchema(dbMetadata, lookupSchema)
	if schemaMeta == nil {
		return "", nil, &ResourceNotFoundError{
			Database: &lookupDatabase,
			Schema:   &lookupSchema,
		}
	}

	for _, lookup := range []func(string) *namedColumns{
		func(name string) *namedColumns {
			meta := findTable(schemaMeta, name)
			if meta == nil {
				return nil
			}
			out := &namedColumns{canonical: meta.GetName()}
			for _, col := range meta.GetColumns() {
				out.columns = append(out.columns, col.GetName())
			}
			return out
		},
		func(name string) *namedColumns {
			meta := findView(schemaMeta, name)
			if meta == nil {
				return nil
			}
			out := &namedColumns{canonical: meta.GetName()}
			for _, col := range meta.GetColumns() {
				out.columns = append(out.columns, col.GetName())
			}
			return out
		},
	} {
		found := lookup(table)
		if found == nil {
			// Case-folded fallback over the object names (the legacy behavior
			// for table identifiers, which resolve case-insensitively).
			for _, name := range append(listTableNames(schemaMeta), listViewNames(schemaMeta)...) {
				if strings.EqualFold(name, table) {
					if found = lookup(name); found != nil {
						break
					}
				}
			}
		}
		if found != nil {
			return found.canonical, found.columns, nil
		}
	}

	return "", nil, &ResourceNotFoundError{
		Database: &lookupDatabase,
		Schema:   &lookupSchema,
		Table:    &table,
	}
}

// getDatabaseMetadata fetches (and caches) metadata for the given database.
func (q *querySpanResolver) getDatabaseMetadata(database string) (*metadata.DatabaseSchemaMetadata, error) {
	if database == "" {
		database = q.opts.DefaultDatabase
	}
	if meta, ok := q.metaCache[database]; ok {
		return meta, nil
	}
	if q.opts.GetDatabaseMetadata == nil {
		return nil, &ResourceNotFoundError{Database: &database}
	}
	meta, err := q.opts.GetDatabaseMetadata(database)
	switch {
	case err != nil:
		var notFound *ResourceNotFoundError
		if errors.As(err, &notFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get database metadata for database: %s: %w", database, err)
	case meta == nil:
		empty := &metadata.DatabaseSchemaMetadata{}
		q.metaCache[database] = empty
		return empty, nil
	default:
		q.metaCache[database] = meta
		return meta, nil
	}
}

func findSchema(meta *metadata.DatabaseSchemaMetadata, name string) *metadata.SchemaMetadata {
	for _, schema := range meta.GetSchemas() {
		if schema.GetName() == name {
			return schema
		}
	}
	return nil
}

func findTable(schema *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	for _, table := range schema.GetTables() {
		if table.GetName() == name {
			return table
		}
	}
	return nil
}

func findView(schema *metadata.SchemaMetadata, name string) *metadata.ViewMetadata {
	for _, view := range schema.GetViews() {
		if view.GetName() == name {
			return view
		}
	}
	return nil
}

func listTableNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.GetTables()))
	for _, table := range schema.GetTables() {
		names = append(names, table.GetName())
	}
	slices.Sort(names)
	return names
}

func listViewNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.GetViews()))
	for _, view := range schema.GetViews() {
		names = append(names, view.GetName())
	}
	slices.Sort(names)
	return names
}
