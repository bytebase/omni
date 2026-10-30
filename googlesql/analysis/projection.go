package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// relation-projection resolution
//
// This is the masking-grade core that mirrors the legacy bytebase resolver's
// TableSource model (plugin/parser/bigquery query_span_extractor.go:
// PseudoTable / PhysicalTable / extractTableSourceFromSelect /
// extractTableSourceFromQuerySetOperation). Every FROM relation — a base table,
// a CTE reference, a derived `( query )` subquery, or a join — resolves to an
// ordered output PROJECTION: a list of columns, each carrying its base-column
// lineage and an IsPlainField flag. A `SELECT *` (or `rel.*`) over such a
// relation reproduces that relation's projection (with lineage), NOT a bare
// catalog-blind "*"; a set operation position-merges its two arms' projections.
//
// The one thing omni leaves to the catalog-aware consumer (the bytebase
// extractor) is enumerating a PHYSICAL table's columns: omni has no metadata, so
// a base-table star is carried as a `baseStar` projection element (a single
// element standing for "every column of table T"), which the consumer expands.
// Everything else — CTE/derived projections, column references resolved through
// relations, set-op merges, join concatenation, plain-field propagation — is
// resolved here so the consumer receives concrete per-column lineage.

// projColumn is one resolved output column of a relation's projection.
type projColumn struct {
	name    string      // output column name (lower/written case; consumer upper-cases explicit names)
	sources []ColumnRef // resolved base-column lineage
	plain   bool        // IsPlainField (legacy semantics)

	// baseFieldName marks name as a base-table FIELD passthrough (today: the
	// JOIN ... USING coalesced key, which the legacy resolver named after the
	// left PhysicalTable's field). A consumer reproducing legacy naming renders
	// it in the field's metadata case rather than the written/upper-cased form.
	baseFieldName bool

	// baseStar, when non-nil, means this projection element stands for "every
	// column of this base table" — omni cannot enumerate it (no catalog), so the
	// consumer expands it. name is "*". sources is empty.
	baseStar *ColumnRef

	// baseStarExcept lists base-table column names the consumer must SKIP when
	// expanding baseStar (case-insensitive). Set by a JOIN ... USING coalesce
	// (fix F1): each side's star excludes the key columns, which are projected
	// once as a coalesced concrete column ahead of the stars. Only meaningful
	// when baseStar is non-nil.
	baseStarExcept []string

	// starMerge, when non-nil, marks a set-operation merge position whose other arm
	// is a base-table star of (consumer-known) arity: the consumer expands
	// starMerge.table and unions that table's column at the same ordinal into this
	// position. Used for `concrete UNION ALL (SELECT * FROM base)` (legacy expands
	// the star arm against metadata, then position-merges). See ColumnInfo.
	starMerge *starMerge

	// starGroup, when non-nil, marks this element as a whole `*` / `rel.*` star
	// ITEM carrying its ordered expansion segments plus the EXCEPT/REPLACE
	// modifiers that apply across them, so the modifier application (with the legacy
	// name-collision last-wins dedup) happens once, after expansion, in the
	// metadata-aware consumer. It is produced ONLY for a star that has EXCEPT/REPLACE
	// modifiers; a modifier-free star (even one with base-table segments) expands
	// inline into individual projColumns instead, keeping each segment individually
	// addressable (so an enclosing relation can resolve `rel.col` and a set
	// operation can position-merge). name is "*".
	starGroup *starGroup

	// setOpMerge, when non-nil, marks a single element standing for a WHOLE
	// set-operation merge whose arms cannot be position-merged inline because at
	// least one arm carries an un-enumerable star (a base-table star, an
	// EXCEPT/REPLACE starGroup, or a nested set-op merge). It carries both arms'
	// resolved projections (in wire form); the metadata-aware consumer expands each
	// arm fully (a base-table star against metadata, a starGroup with its modifiers,
	// a nested merge recursively), then position-merges the two expanded lists —
	// output names from the LEFT arm, sources unioned per position, plain=false —
	// reproducing the legacy "fully resolve each arm, then zip"
	// (extractTableSourceFromQuerySetOperation). name is "*". This is the masking-
	// grade fix for star-involving set-op arms whose arity only the consumer knows:
	// a base-star UNION base-star (one merged col per expanded position, NOT two
	// concatenated stars), a nested concrete/star/concrete (the inner star-merge
	// survives the outer merge), and an EXCEPT/REPLACE star arm merged with a
	// concrete arm.
	setOpMerge *setOpMerge
}

// setOpMerge carries a deferred set-operation merge: the two arms' resolved
// projections (wire form), to be expanded and merged by the metadata-aware
// consumer. byName marks a BY NAME / CORRESPONDING operation (fix F3): the
// consumer merges the expanded arms by column NAME (matchColumns restricting/
// ordering the output when present) instead of by ordinal. See
// projColumn.setOpMerge and SetOpMergeInfo.
type setOpMerge struct {
	left         []ColumnInfo
	right        []ColumnInfo
	byName       bool
	matchColumns []string
}

// starGroup is a `*` / `rel.*` star item's resolved expansion: the ordered
// segments to emit and its EXCEPT/REPLACE modifiers (applied across the expanded
// columns by the metadata-aware consumer). See ColumnInfo.StarSegments.
type starGroup struct {
	segments []projColumn
	except   []string
	replace  []StarReplaceItem
}

// starMerge carries a set-operation merge against a base-table-star arm: the base
// table to expand (table) and the output ordinal (index) whose lineage gains that
// table's index-th column. leftStar marks the star arm as the LEFT one (so the
// output NAME is taken from the expanded base column, not from a concrete arm).
type starMerge struct {
	table    ColumnRef
	index    int
	leftStar bool
}

// relation is a resolved FROM source: a base table, a CTE reference, a derived
// subquery, an UNNEST value relation, or a join of these.
type relation struct {
	// name is the relation's reference name for qualified lookups: the alias if
	// present, else the (bare) table name; "" for an unnamed join/derived.
	name string
	// isBase marks a single physical base table (so a bare reference qualified by
	// its name resolves a column to the base table's lineage even though omni
	// cannot enumerate its columns).
	isBase bool
	// baseRef is the base table reference (Database/Schema/Table) when isBase, so a
	// column read off it (`t.col` / a bare `col`) resolves to {…, Column: col}.
	baseRef ColumnRef
	// valueTable marks an UNNEST / correlated-array-path value relation (fix F2):
	// its first column IS the element value, so a field reference through the
	// relation name (`elem.child`) that matches no projection column resolves to
	// the ELEMENT's lineage (a struct field reads the element).
	valueTable bool
	// columns is the ordered resolved projection. For a base table it is a single
	// baseStar element. It is never changed once the relation is built.
	columns []projColumn
	// index answers lookups of columns by name (see projectionColumnSources).
	index *relIndex
}

// resolveSelectProjection resolves a SELECT block to its output projection,
// mirroring the legacy extractTableSourceFromSelect: a bare `*` copies all FROM
// relations' columns; a `rel.*` copies the named relation's columns; an explicit
// item resolves its expression's columns against the relations. fromRels is the
// comma-item relation list (a join collapsed into one combined relation) used to
// expand a bare `*` in FROM order; column / `rel.*` lookups use the SELECT's leaf
// relations (w.leaves), which keep each join side findable by name.
func (w *spanWalker) resolveSelectProjection(stmt *ast.SelectStmt, fromRels []*relation) []projColumn {
	var out []projColumn
	for _, item := range stmt.Items {
		if item == nil {
			continue
		}
		switch {
		case item.Star && item.Expr == nil:
			// Bare `*`: every comma-item relation's columns, in order.
			out = append(out, w.starProjColumns(allColumns(fromRels), item.Modifiers)...)
		case item.Star && item.Expr != nil:
			// `expr.*`: the named relation's columns (or a field-path star).
			out = append(out, w.starProjColumns(w.resolveDotStar(item.Expr), item.Modifiers)...)
		default:
			out = append(out, w.resolveExprColumn(item))
		}
	}
	return out
}

// starProjColumns turns a star item's resolved expansion `segments` into one or
// more projColumns. A modifier-free, fully-concrete star (every segment already
// resolved, no EXCEPT/REPLACE) expands inline into its concrete columns — so a
// set operation can position-merge them and the consumer needs no metadata. A
// star that has modifiers OR contains an un-enumerable base-table segment is kept
// as a single starGroup projColumn, so the consumer expands the base segments and
// applies the modifiers (with the legacy name-collision last-wins dedup) in one
// place.
func (w *spanWalker) starProjColumns(segments []projColumn, mods *ast.StarModifiers) []projColumn {
	except, replace := w.resolveStarModifiers(mods)
	if len(except) == 0 && len(replace) == 0 {
		// Modifier-free: inline the segments (concrete columns stay concrete, a
		// base-table segment stays a baseStar projColumn). Inlining — rather than
		// wrapping in a starGroup — keeps each segment individually addressable so an
		// enclosing relation can resolve `rel.col` / position-merge a base-star arm,
		// and a fully-concrete star can be position-merged by a set operation.
		return segments
	}
	// With EXCEPT/REPLACE the expansion + name-collision dedup must be applied as a
	// unit (and a base segment can only be expanded by the metadata-aware consumer),
	// so keep the star grouped with its modifiers.
	return []projColumn{{
		name: "*",
		starGroup: &starGroup{
			segments: segments,
			except:   except,
			replace:  replace,
		},
	}}
}

// resolveStarModifiers resolves a star item's EXCEPT/REPLACE modifiers into the
// consumer-facing form: the EXCEPT column names, and per-REPLACE the output
// column name plus its replacement expression's resolved source columns. Returns
// (nil, nil) when there are no modifiers.
func (w *spanWalker) resolveStarModifiers(mods *ast.StarModifiers) ([]string, []StarReplaceItem) {
	if mods == nil {
		return nil, nil
	}
	var except []string
	if len(mods.Except) > 0 {
		except = append(except, mods.Except...)
	}
	var replace []StarReplaceItem
	for _, r := range mods.Replace {
		if r == nil || r.Alias == "" {
			continue
		}
		_, sources := w.resolveExprSources(r.Value)
		replace = append(replace, StarReplaceItem{Name: r.Alias, Sources: sources})
	}
	return except, replace
}

// resolveExprColumn resolves a non-star select item to a single output column,
// mirroring the legacy extractSourceColumnSetFromExpr + alias rule. The item is
// never a plain field (an explicit select-list item is IsPlainField=false).
func (w *spanWalker) resolveExprColumn(item *ast.SelectItem) projColumn {
	name, sources := w.resolveExprSources(item.Expr)
	if item.Alias != "" {
		name = item.Alias
	}
	return projColumn{name: name, sources: sources, plain: false}
}

// resolveExprSources resolves an expression's directly-referenced columns to
// their base lineage against the FROM relations, returning a best-effort output
// name and the merged source set (the legacy extractSourceColumnSetFromExpr): each
// maximal dotted column path is resolved by resolvePath; a multi-reference
// expression has no derivable name.
func (w *spanWalker) resolveExprSources(expr ast.Node) (string, []ColumnRef) {
	if expr == nil || w.tablesOnly {
		return "", nil
	}
	paths := w.collectColumnPaths(expr)
	if len(paths) == 0 {
		return "", nil
	}
	// One refUnion for the whole expression: re-unioning the sources gathered
	// so far with each path's made `a1 + a2 + ... + an` cost O(n^2).
	var sources refUnion
	name := ""
	for _, parts := range paths {
		n, refs := w.resolvePath(parts)
		sources.add(refs)
		name = n
	}
	if len(paths) > 1 {
		name = ""
	}
	return name, sources.refs
}

// resolvePath resolves one dotted column path (e.g. [col] or [rel, col] or
// [s, field]) to its output name and base lineage. It covers the same cases as the
// legacy getFieldColumnSource, but resolves RELATION-FIRST rather than column-
// first, because omni — unlike the legacy resolver — has no catalog to check
// whether a leading identifier is instead a base-table column:
//   - 1 part [col]: resolve `col` as a column (name = col).
//   - >=2 parts [a, …]: if `a` names a leaf relation, resolve the second part as
//     that relation's column (name = b; trailing struct sub-fields dropped). Else
//     if the dialect-bucketed qualifier names a leaf base table (`proj.ds.t.c`),
//     keep the bucketed ref. Else `a` is a struct/field column root (`s.field`,
//     `s.a.b`) — resolve `a` as a column, dropping the field path (name = a).
//
// A column owned by a CTE / derived relation resolves to that relation's stored
// lineage. A column owned by (or assumed to be on) a base table resolves to a
// written ref the consumer matches against catalog metadata.
func (w *spanWalker) resolvePath(parts []string) (string, []ColumnRef) {
	if len(parts) == 0 || w.tablesOnly {
		return "", nil
	}
	if len(parts) == 1 {
		col := parts[0]
		return col, w.resolveColumn("", col)
	}
	// `rel.col` where the leading identifier names a leaf relation: resolve the
	// SECOND part as that relation's column (any trailing parts are struct
	// sub-fields, which carry no extra base lineage). Relation-first because omni,
	// unlike the legacy resolver, cannot metadata-check whether the leading
	// identifier is instead a column of a base table.
	if w.leaves.find(parts[0]) != nil {
		return parts[1], w.resolveColumn(parts[0], parts[1])
	}
	// A FULLY-qualified column reference whose qualifier names a leaf base relation
	// (e.g. `proj.ds.t.c` over `FROM proj.ds.t`): bucket it by the dialect rule so
	// its Catalog/Database/Schema/Table line up with that table's access.
	ref := columnRefFromParts(parts, w.dialect)
	if w.qualifierNamesBaseRelation(ref) {
		return ref.Column, []ColumnRef{ref}
	}
	// Otherwise the leading identifier is a struct/field column root (`s.field`,
	// `s.a.b`): resolve `parts[0]` as a column (dropping the trailing field path),
	// mirroring the legacy a-as-column attempt for an unqualified field path.
	return parts[0], w.resolveColumn("", parts[0])
}

// qualifierNamesBaseRelation reports whether a bucketed column ref's qualifier
// (its Catalog/Database/Schema/Table) matches a leaf BASE relation's table — i.e.
// the reference is a fully-qualified column of a FROM base table (e.g.
// `proj.ds.t.c`). Used to keep such a reference's dialect-bucketed qualifier
// rather than misreading its leading identifier as a struct-column root.
func (w *spanWalker) qualifierNamesBaseRelation(ref ColumnRef) bool {
	if ref.Table == "" {
		return false
	}
	return w.leaves.hasBase(ref)
}

// resolveColumn resolves a (relationQualifier, column) reference to base lineage
// against the SELECT's leaf relations (w.leaves). relQualifier may be ""
// (unqualified). A CTE/derived relation's column resolves to its stored lineage; a
// base (or base-star-passthrough) relation's column resolves to a written ref the
// consumer matches against metadata; an unresolved reference falls back to a
// written ref (so the consumer can still match it by name).
func (w *spanWalker) resolveColumn(relQualifier, column string) []ColumnRef {
	if refs, ok := w.resolveColumnStrict(relQualifier, column); ok {
		return refs
	}
	// Fallback: a written reference the consumer matches against catalog metadata.
	// Qualified by a base relation → carry the relation's base ref; unqualified or
	// over an unknown relation → a bare column ref.
	if relQualifier != "" {
		if rel := w.leaves.find(relQualifier); rel != nil && rel.isBase {
			ref := rel.baseRef
			ref.Column = column
			return []ColumnRef{ref}
		}
		// Qualifier is not a known relation: emit a written ref (alias.column) so
		// the consumer can still match by column name (table set to the qualifier).
		return []ColumnRef{{Table: relQualifier, Column: column}}
	}
	// Unqualified with exactly ONE base relation in scope: the column can only
	// come from that relation — the strict pass already ruled out every concrete
	// CTE/derived projection — so attribute it. Without the attribution the
	// consumer falls back to matching the bare name across ALL expanded tables,
	// which over-includes whenever another table shares the column name (e.g. a
	// 3-way UNION whose arms' tables have overlapping column names) and diverges
	// from the legacy resolver's exact single-table attribution.
	if rel := w.leaves.soleBase(); rel != nil {
		ref := rel.baseRef
		ref.Column = column
		return []ColumnRef{ref}
	}
	return []ColumnRef{{Column: column}}
}

// resolveColumnStrict resolves a (relationQualifier, column) reference ONLY when
// it can be tied to a concrete relation projection. It returns nil when no leaf
// relation resolves it, so the caller can apply the legacy fallback.
//
//   - Qualified `rel.col`: the named leaf relation's column. A base or base-star-
//     passthrough relation resolves `col` to its base-table lineage (the written
//     column name keyed to the base table); a CTE/derived relation matches `col`
//     by name in its concrete projection.
//   - Unqualified `col`: only a CTE/derived relation with a CONCRETE column named
//     `col` resolves here. A base table's columns are not enumerable, so an
//     unqualified column over base tables is left to the consumer to match by name
//     (omni cannot know which base table — or which arm of a join — owns it).
func (w *spanWalker) resolveColumnStrict(relQualifier, column string) ([]ColumnRef, bool) {
	if relQualifier != "" {
		rel := w.leaves.find(relQualifier)
		if rel == nil {
			return nil, false
		}
		if rel.isBase {
			ref := rel.baseRef
			ref.Column = column
			return []ColumnRef{ref}, true
		}
		if refs, ok := projectionColumnSources(rel, column, true /* allowBaseStar */); ok {
			return refs, true
		}
		// A VALUE relation (UNNEST element): a qualified reference that names no
		// projection column is a struct FIELD of the element (`elem.child`), so it
		// reads the element's lineage (fix F2).
		if rel.valueTable && len(rel.columns) > 0 {
			return rel.columns[0].sources, true
		}
		return nil, false
	}
	return w.leaves.derivedColumn(column)
}

// resolveDotStar resolves a `rel.*` (or qualified/field-path star) to the named
// relation's projection columns, mirroring the legacy extractWildFromExpr: the
// leading path names a relation whose columns are reproduced. A multi-part
// qualifier (`schema.table.*` / `db.schema.table.*`) matches the relation by its
// TRAILING part with the prefix verified against the relation's base reference —
// the legacy resolver ERRORED on these (fail-closed), so resolving them is an
// omni improvement, but an UNRESOLVABLE wild path must FAIL CLOSED too (the
// structural rule): silently yielding no columns would give the masker zero
// result maskers and return every output column unmasked.
func (w *spanWalker) resolveDotStar(expr ast.Node) []projColumn {
	parts := exprToParts(expr)
	if len(parts) == 0 {
		return nil
	}
	// `rel.*`: a SINGLE-part qualifier naming a leaf relation reproduces that
	// relation's projection. The single-part guard is load-bearing: a MULTI-part
	// path whose head names a relation is a STRUCT-FIELD star through it
	// (`d.s.*` — the star expands the struct field's sub-columns, which omni
	// cannot enumerate metadata-free); returning the relation's whole projection
	// would misalign the positional masker (gate-round finding: the first
	// result's masker lands on the struct's first sub-column → leak). Those fall
	// through to the fail-closed branch, matching the legacy resolver's error.
	if len(parts) == 1 {
		if rel := w.leaves.find(parts[0]); rel != nil {
			if rel.valueTable {
				// An UNNEST/value relation's `elem.*` expands the ELEMENT's struct
				// sub-fields — N engine output columns omni cannot enumerate;
				// returning the relation's single projection column would shift
				// every later position (re-verify finding). Fail closed.
				w.failClosed(fmt.Errorf("cannot enumerate %s.* over a value table's element (fail closed)", parts[0]))
				return nil
			}
			return rel.columns
		}
	}
	// `schema.table.*` / `dataset.table.*`: the head does NOT name a relation,
	// the trailing part names a base relation, and the written prefix matches the
	// relation's schema OR database/dataset qualifier (non-empty — an unqualified
	// FROM accepts no written prefix; the engine would reject the mismatched
	// range variable anyway, and the legacy resolver errored on every
	// schema-qualified star).
	if len(parts) >= 2 && w.leaves.find(parts[0]) == nil {
		last := parts[len(parts)-1]
		if rel := w.leaves.find(last); rel != nil && rel.isBase {
			prefix := parts[len(parts)-2]
			schemaMatch := rel.baseRef.Schema != "" && strings.EqualFold(rel.baseRef.Schema, prefix)
			databaseMatch := rel.baseRef.Database != "" && strings.EqualFold(rel.baseRef.Database, prefix)
			if schemaMatch || databaseMatch {
				return rel.columns
			}
		}
	}
	// Unresolvable wild path (a struct-path star or a qualifier naming nothing in
	// FROM): fail closed rather than silently produce zero output columns.
	w.failClosed(fmt.Errorf("cannot resolve %s.* to a FROM relation's columns (fail closed)", strings.Join(parts, ".")))
	return nil
}

// mergeProjections merges two set-operation arms' projections position-wise,
// reproducing the legacy extractTableSourceFromQuerySetOperation (fully resolve
// each arm, then zip: output names from the LEFT arm, sources unioned per
// position). Because omni — unlike the catalog-aware legacy resolver — cannot
// expand a base-table star before the merge, an arm carrying an un-enumerable
// star is handled by deferring the whole merge to the metadata-aware consumer.
//
//   - Both arms fully concrete (no star markers): a plain position-wise union.
//     Arity mismatch ⇒ left wins unchanged (the legacy length-mismatch guard;
//     legacy errors on an unequal-width set op, which the masking layer rejects —
//     omni's left-only best effort never DROPS a left column, so it does not
//     under-attribute the surviving result).
//   - Concrete LEFT + single base-star RIGHT (the common `concrete UNION ALL
//     SELECT * FROM base`): keep the left columns (names from the left) and mark
//     each with a starMerge so the consumer expands the right base table and
//     unions its i-th column into position i. plain=false.
//   - Single base-star LEFT + concrete RIGHT: symmetric, names from the expanded
//     LEFT base table (legacy first-select-name) via starMerge{leftStar:true}.
//   - Both single base-star OVER THE SAME TABLE: idempotent — the merge of a base
//     table's star with itself is that same star (a plain passthrough). This is
//     the recursive-CTE base-star shape (`(SELECT * FROM seed) UNION ALL SELECT *
//     FROM r`, where the recursive self-reference resolves back to seed); legacy
//     publishes seed's columns once (plain), not duplicated.
//   - Any other star-involving combination (base-star × base-star over DIFFERENT
//     tables, an EXCEPT/REPLACE starGroup arm, a nested set-op merge whose arm
//     carries a starMerge, …): defer the whole merge to the consumer via a single
//     setOpMerge element carrying both arms' wire-form projections. The consumer
//     expands each arm fully and position-merges — the only place the star arities
//     are known. This preserves every arm's lineage (no concatenation, no dropped
//     star marker), closing the set-op masking leaks.
func mergeProjections(left, right []projColumn) []projColumn {
	if len(right) == 0 {
		return left
	}
	if len(left) == 0 {
		return right
	}
	lStar := singleBaseStar(left)
	rStar := singleBaseStar(right)
	lConcrete := allConcrete(left)
	rConcrete := allConcrete(right)

	switch {
	case lConcrete && rConcrete:
		// A position-wise union of the sources, names from the left (see
		// setOpMerger.merge).
		m := setOpMerger{cols: left}
		m.merge(right)
		return m.result()
	case lConcrete && rStar != nil:
		// Concrete LEFT, single base-star RIGHT: per-position starMerge.
		merged := make([]projColumn, len(left))
		for i := range left {
			merged[i] = projColumn{
				name:      left[i].name,
				sources:   left[i].sources,
				plain:     false,
				starMerge: &starMerge{table: *rStar, index: i},
			}
		}
		return merged
	case lStar != nil && rConcrete:
		// Single base-star LEFT, concrete RIGHT: per-position starMerge (names from
		// the expanded left base table).
		merged := make([]projColumn, len(right))
		for i := range right {
			merged[i] = projColumn{
				name:      "", // named from the expanded left base table by the consumer
				sources:   right[i].sources,
				plain:     false,
				starMerge: &starMerge{table: *lStar, index: i, leftStar: true},
			}
		}
		return merged
	case lStar != nil && rStar != nil && sameBaseTable(*lStar, *rStar):
		// Both single base-star over the SAME table: idempotent passthrough (the
		// recursive-CTE base-star shape). Keep one base-star, plain.
		return []projColumn{{name: "*", plain: true, baseStar: lStar}}
	default:
		// Any other star-involving combination: defer the whole merge to the
		// consumer, which expands each arm's stars and position-merges.
		return []projColumn{{
			name: "*",
			setOpMerge: &setOpMerge{
				left:  projColumnsToColumnInfos(left),
				right: projColumnsToColumnInfos(right),
			},
		}}
	}
}

// mergeProjectionsByName merges a BY NAME / CORRESPONDING set operation's arms
// by case-insensitive column NAME (fix F3), reproducing the engine's column
// matching: output order is the LEFT arm's columns (each unioning every
// same-named right column's lineage), followed by any right-only names (the
// FULL BY NAME shape; for the plain/LEFT variants a trailing extra never
// shifts an earlier position, so including it is over-attribution-safe). A
// non-empty matchColumns (`ON (cols)` / `BY (cols)`) restricts the output to
// exactly those columns, in list order. Merged columns are never plain.
//
// When either arm carries a non-concrete marker (a base-table star, a star
// group, or a nested merge) the name partition needs the arm's expanded column
// NAMES, which only the metadata-aware consumer knows — the whole merge defers
// via a single setOpMerge element with byName set (NEVER a silent ordinal
// merge, which is the verified mis-attribution this fix closes).
func mergeProjectionsByName(left, right []projColumn, matchColumns []string) []projColumn {
	if len(right) == 0 {
		return left
	}
	if len(left) == 0 {
		return right
	}
	if !allConcrete(left) || !allConcrete(right) {
		return []projColumn{{
			name: "*",
			setOpMerge: &setOpMerge{
				left:         projColumnsToColumnInfos(left),
				right:        projColumnsToColumnInfos(right),
				byName:       true,
				matchColumns: matchColumns,
			},
		}}
	}
	if len(matchColumns) > 0 {
		l, r := newNamedSources(left), newNamedSources(right)
		out := make([]projColumn, 0, len(matchColumns))
		for _, m := range matchColumns {
			k := foldKey(m)
			out = append(out, projColumn{
				name:    m,
				sources: unionColumnRefs(l.get(k), r.get(k)),
				plain:   false,
			})
		}
		return out
	}
	// The left columns, each unioned with the same-named right columns, then the
	// right-only columns (see setOpMerger.mergeByName).
	m := setOpMerger{cols: left}
	m.mergeByName(right, nil)
	return m.result()
}

// setOpMerger merges the arms of a set-operation chain one at a time, each
// merge giving what mergeProjections or mergeProjectionsByName gives for the
// projection so far and the next arm. It builds the concrete unions itself
// (position-wise, or by name without a column list) and hands every other
// merge to those functions. While every merge so far was a concrete union, the
// sources of each output column are held in a refUnion and extended in place.
// Building a new union of both arms' full source lists at every arm copied the
// sources gathered so far again at each one, so a chain of n arms reading n
// different tables cost O(n^2).
type setOpMerger struct {
	// cols is the projection so far. While unions is non-nil, the merges were
	// concrete unions: each column has only its name set, and unions[i] holds
	// column i's sources.
	cols   []projColumn
	unions []refUnion
	// byName maps foldKey of each name in cols to the columns with that name,
	// in order. A BY NAME merge builds it; it is dropped with unions.
	byName map[string][]int
	// concrete caches allConcrete(cols) while checked is set.
	concrete, checked bool
}

// colsConcrete reports whether every column so far is concrete. A BY NAME
// chain can widen the projection at every arm, so checking it all again at
// each arm cost O(n^2).
func (m *setOpMerger) colsConcrete() bool {
	if m.unions != nil {
		return true // accumulated columns are concrete by construction
	}
	if !m.checked {
		m.concrete, m.checked = allConcrete(m.cols), true
	}
	return m.concrete
}

// merge merges the next arm position-wise, as mergeProjections does.
func (m *setOpMerger) merge(right []projColumn) {
	if len(right) == 0 {
		return
	}
	if len(m.cols) > 0 && m.colsConcrete() && allConcrete(right) {
		if len(m.cols) != len(right) {
			return // arity mismatch: the left projection wins unchanged
		}
		m.accumulate()
		for i := range right {
			m.unions[i].add(right[i].sources)
		}
		return
	}
	m.cols, m.checked = mergeProjections(m.result(), right), false
}

// mergeByName merges the next arm by column name, as mergeProjectionsByName
// does.
func (m *setOpMerger) mergeByName(right []projColumn, matchColumns []string) {
	if len(right) == 0 {
		return
	}
	if len(m.cols) > 0 && len(matchColumns) == 0 && m.colsConcrete() && allConcrete(right) {
		m.accumulate()
		if m.byName == nil {
			m.byName = make(map[string][]int, len(m.cols))
			for i, c := range m.cols {
				k := foldKey(c.name)
				m.byName[k] = append(m.byName[k], i)
			}
		}
		// Each column so far gains the sources of the right columns with its
		// name.
		r := newNamedSources(right)
		for k := range r.byName {
			if cols := m.byName[k]; len(cols) > 0 {
				srcs := r.get(k)
				for _, i := range cols {
					m.unions[i].add(srcs)
				}
			}
		}
		// Then each right column whose name none of them has is appended.
		n := len(m.cols)
		for _, c := range right {
			if _, ok := m.byName[foldKey(c.name)]; !ok {
				m.cols = append(m.cols, projColumn{name: c.name})
				m.unions = append(m.unions, refUnion{refs: c.sources})
			}
		}
		for i := n; i < len(m.cols); i++ {
			k := foldKey(m.cols[i].name)
			m.byName[k] = append(m.byName[k], i)
		}
		return
	}
	m.cols, m.checked = mergeProjectionsByName(m.result(), right, matchColumns), false
}

// accumulate turns the projection so far, all concrete, into the accumulated
// form: a merged column keeps only its name (it is never plain), and its
// sources move to unions.
func (m *setOpMerger) accumulate() {
	if m.unions != nil {
		return
	}
	cols := make([]projColumn, len(m.cols))
	m.unions = make([]refUnion, len(m.cols))
	for i, c := range m.cols {
		cols[i] = projColumn{name: c.name}
		m.unions[i] = refUnion{refs: c.sources}
	}
	m.cols = cols
}

// result returns the projection so far, with the accumulated sources written
// back into its columns.
func (m *setOpMerger) result() []projColumn {
	if m.unions != nil {
		for i := range m.cols {
			m.cols[i].sources = m.unions[i].refs
		}
		m.unions = nil
		m.byName = nil
	}
	return m.cols
}

// namedSources gives, for a folded name, the union of the sources of every
// column of cols whose name matches it case-insensitively, in column order.
type namedSources struct {
	cols   []projColumn
	byName map[string][]int       // foldKey(name) → the columns with that name
	unions map[string][]ColumnRef // computed on first use
}

func newNamedSources(cols []projColumn) *namedSources {
	x := &namedSources{cols: cols, byName: make(map[string][]int, len(cols))}
	for i, c := range cols {
		k := foldKey(c.name)
		x.byName[k] = append(x.byName[k], i)
	}
	return x
}

// get returns the union for the folded name key.
func (x *namedSources) get(key string) []ColumnRef {
	if srcs, ok := x.unions[key]; ok {
		return srcs
	}
	var u refUnion
	for _, i := range x.byName[key] {
		u.add(x.cols[i].sources)
	}
	if x.unions == nil {
		x.unions = make(map[string][]ColumnRef)
	}
	x.unions[key] = u.refs
	return u.refs
}

// coalesceUsingColumns builds the projection of `left JOIN right USING (keys)`
// (fix F1): each key ONCE, first, in USING order — its lineage the union of
// BOTH sides' key columns — then the left side's non-key columns (rewrapped
// non-plain, the legacy join-left rule), then the right side's (plainness
// kept). leftEnd is the length of the keys and left columns together. It
// returns ok=false when a side cannot be name-partitioned without metadata (a
// deferred set-op marker in its projection); the caller fails closed. The
// result is a new slice; neither input is modified.
func coalesceUsingColumns(left, right []projColumn, keys []string) (cols []projColumn, leftEnd int, ok bool) {
	leftRest, ok := stripUsingKeys(left, keys)
	if !ok {
		return nil, 0, false
	}
	rightRest, ok := stripUsingKeys(right, keys)
	if !ok {
		return nil, 0, false
	}
	cols = make([]projColumn, 0, len(keys)+len(leftRest)+len(rightRest))
	cols = append(cols, usingKeyColumns(left, right, keys)...)
	for _, c := range leftRest {
		cols = append(cols, rewrapNonPlain(c))
	}
	leftEnd = len(cols)
	cols = append(cols, rightRest...)
	return cols, leftEnd, true
}

// usingKeyColumns returns the coalesced key columns of a USING join, one per
// key in USING order, each reading the key's lineage from both sides.
func usingKeyColumns(left, right []projColumn, keys []string) []projColumn {
	l, r := usingKeyLineages(left, keys), usingKeyLineages(right, keys)
	cols := make([]projColumn, 0, len(keys))
	for i, key := range keys {
		cols = append(cols, projColumn{
			name:          key,
			sources:       unionColumnRefs(l[i].refs, r[i].refs),
			plain:         false,
			baseFieldName: true,
		})
	}
	return cols
}

// sameNameSetFold reports whether a and b have the same length and hold the
// same names under case-insensitive comparison.
func sameNameSetFold(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	setA, setB := newFoldSet(a), newFoldSet(b)
	for _, n := range a {
		if !setB.has(n) {
			return false
		}
	}
	for _, n := range b {
		if !setA.has(n) {
			return false
		}
	}
	return true
}

// stripUsingKeys returns a join side's projection with the USING key columns
// removed: a concrete column named like a key is dropped (it is owned by the
// coalesced key column); a base-table star carries the keys in its except list
// (the consumer's expansion skips them); a star group adds them to its EXCEPT
// modifier. ok=false when an element cannot be name-partitioned without
// metadata (a starMerge / setOpMerge marker, whose column names only the
// consumer knows). All returned elements are copies — the side relations stay
// leaf-resolvable with their own un-coalesced projections.
func stripUsingKeys(cols []projColumn, keys []string) ([]projColumn, bool) {
	keySet := newFoldSet(keys)
	out := make([]projColumn, 0, len(cols))
	for _, c := range cols {
		switch {
		case c.starMerge != nil || c.setOpMerge != nil:
			return nil, false
		case c.starGroup != nil:
			for _, s := range c.starGroup.segments {
				if s.starMerge != nil || s.setOpMerge != nil || s.starGroup != nil {
					return nil, false
				}
			}
			c.starGroup = &starGroup{
				segments: append([]projColumn{}, c.starGroup.segments...),
				except:   appendNamesFold(c.starGroup.except, keys),
				replace:  c.starGroup.replace,
			}
			out = append(out, c)
		case c.baseStar != nil:
			c.baseStarExcept = appendNamesFold(c.baseStarExcept, keys)
			out = append(out, c)
		case keySet.has(c.name):
			// Dropped: this side's key column is folded into the coalesced key.
		default:
			out = append(out, c)
		}
	}
	return out, true
}

// usingKeyLineages returns, for each USING key, the lineage a join side
// contributes to it: every concrete column named like it, every base-table
// star that may own it (the star's table with Column=key — when several stars
// could own the key, all contribute: over-attribution is safe, under is not),
// and a star group's segments likewise (its EXCEPT removing the key means the
// side does not project it; a REPLACE re-points it to the replacement's
// sources). Each key's lineage is in column order. The keys are resolved
// together in one pass over the side, which a wide USING list over a wide side
// needs: a pass per key costs the product of the two widths.
func usingKeyLineages(cols []projColumn, keys []string) []refUnion {
	folded := make([]string, len(keys))
	all := make([]int, len(keys))
	for i, key := range keys {
		folded[i] = foldKey(key)
		all[i] = i
	}
	out := make([]refUnion, len(keys))
	addUsingKeyLineages(cols, keys, folded, all, out)
	return out
}

// addUsingKeyLineages adds to out[i] the lineage cols contribute to keys[i],
// for each i in which; folded[i] is foldKey(keys[i]).
func addUsingKeyLineages(cols []projColumn, keys, folded []string, which []int, out []refUnion) {
	if len(which) == 0 {
		return
	}
	var byName map[string][]int // a folded key → its entries of which
	for _, c := range cols {
		switch {
		case c.starMerge != nil || c.setOpMerge != nil:
			// Unreachable behind stripUsingKeys' ok=false; contribute nothing.
		case c.starGroup != nil:
			except := newFoldSet(c.starGroup.except)
			replace := make(map[string][]ColumnRef, len(c.starGroup.replace))
			for _, r := range c.starGroup.replace {
				if k := foldKey(r.Name); !hasRefs(replace, k) {
					replace[k] = r.Sources
				}
			}
			var rest []int
			for _, i := range which {
				if _, ok := except[folded[i]]; ok {
					continue
				}
				if srcs, ok := replace[folded[i]]; ok {
					out[i].add(srcs)
					continue
				}
				rest = append(rest, i)
			}
			addUsingKeyLineages(c.starGroup.segments, keys, folded, rest, out)
		case c.baseStar != nil:
			except := newFoldSet(c.baseStarExcept)
			for _, i := range which {
				if _, ok := except[folded[i]]; ok {
					continue
				}
				ref := *c.baseStar
				ref.Column = keys[i]
				out[i].add([]ColumnRef{ref})
			}
		default:
			if byName == nil {
				byName = make(map[string][]int, len(which))
				for _, i := range which {
					byName[folded[i]] = append(byName[folded[i]], i)
				}
			}
			for _, i := range byName[foldKey(c.name)] {
				out[i].add(c.sources)
			}
		}
	}
}

func hasRefs(m map[string][]ColumnRef, k string) bool {
	_, ok := m[k]
	return ok
}

// rewrapNonPlain marks a join-left projection column non-plain (the legacy
// joinTable rewraps the anchor side without IsPlainField), including a star
// element's per-segment plainness (the wire carries plainness on the star /
// its segments).
func rewrapNonPlain(c projColumn) projColumn {
	c.plain = false
	if c.starGroup != nil {
		segs := make([]projColumn, len(c.starGroup.segments))
		for i, s := range c.starGroup.segments {
			s.plain = false
			segs[i] = s
		}
		c.starGroup = &starGroup{segments: segs, except: c.starGroup.except, replace: c.starGroup.replace}
	}
	return c
}

// appendNamesFold returns a COPY of list with each name appended unless already
// present (case-insensitive). The copy keeps the caller's slice unaliased (a
// join side's projection is shared with its leaf relation).
func appendNamesFold(list, names []string) []string {
	out := append([]string{}, list...)
	seen := newFoldSet(list)
	for _, n := range names {
		k := foldKey(n)
		if _, ok := seen[k]; !ok {
			if seen == nil {
				seen = make(foldSet, len(names))
			}
			seen[k] = struct{}{}
			out = append(out, n)
		}
	}
	return out
}

// allConcrete reports whether every element of a projection is a concrete output
// column (no baseStar / starGroup / starMerge / setOpMerge marker) — i.e. the
// projection can be position-merged inline without metadata.
func allConcrete(cols []projColumn) bool {
	for _, c := range cols {
		if c.baseStar != nil || c.starGroup != nil || c.starMerge != nil || c.setOpMerge != nil {
			return false
		}
	}
	return true
}

// sameBaseTable reports whether two base-table refs name the same physical table
// (case-insensitive on every qualifier component, the GoogleSQL identifier rule).
func sameBaseTable(a, b ColumnRef) bool {
	return strings.EqualFold(a.Catalog, b.Catalog) &&
		strings.EqualFold(a.Database, b.Database) &&
		strings.EqualFold(a.Schema, b.Schema) &&
		strings.EqualFold(a.Table, b.Table)
}

// allColumns concatenates every relation's projection columns in FROM order
// (the legacy `fromFields`).
func allColumns(rels []*relation) []projColumn {
	var out []projColumn
	for _, rel := range rels {
		out = append(out, rel.columns...)
	}
	return out
}

// singleBaseStar returns the base-table ref when the projection is exactly a
// single baseStar element (an un-enumerable base table star), else nil. A star
// carrying a USING except list does not qualify — the starMerge / idempotent
// merge paths have no way to carry the exclusions, so such a projection takes
// the deferred-merge path instead (which preserves them in wire form).
func singleBaseStar(cols []projColumn) *ColumnRef {
	if len(cols) == 1 && cols[0].baseStar != nil && cols[0].starMerge == nil && len(cols[0].baseStarExcept) == 0 {
		return cols[0].baseStar
	}
	return nil
}

// projColumnsToColumnInfos converts a resolved top-level projection into the
// ColumnInfo wire form the metadata-aware consumer (the bytebase extractor)
// reads. A concrete column maps to a plain ColumnInfo; a base-table star (or a
// star item with modifiers / base segments) maps to a StarSegments item the
// consumer expands; a set-operation star-merge position carries StarMerge.
func projColumnsToColumnInfos(cols []projColumn) []ColumnInfo {
	if cols == nil {
		return nil
	}
	out := make([]ColumnInfo, 0, len(cols))
	for _, c := range cols {
		switch {
		case c.setOpMerge != nil:
			out = append(out, ColumnInfo{
				Name: "*",
				SetOpMerge: &SetOpMergeInfo{
					Left:         c.setOpMerge.left,
					Right:        c.setOpMerge.right,
					ByName:       c.setOpMerge.byName,
					MatchColumns: c.setOpMerge.matchColumns,
				},
			})
		case c.starGroup != nil:
			out = append(out, ColumnInfo{
				Name:         "*",
				StarSegments: segmentsToWire(c.starGroup.segments),
				StarExcept:   c.starGroup.except,
				StarReplace:  c.starGroup.replace,
			})
		case c.baseStar != nil && c.starMerge == nil:
			// A standalone base-table star (a bare `SELECT *` over a single base
			// table, or a `rel.*` over a base table) — one segment to expand,
			// skipping any USING-coalesced key columns.
			ref := *c.baseStar
			out = append(out, ColumnInfo{
				Name:         "*",
				StarSegments: []StarSegment{{BaseTable: &ref, ExceptColumns: c.baseStarExcept, Plain: c.plain}},
			})
		case c.starMerge != nil:
			info := ColumnInfo{
				Name:          c.name,
				SourceColumns: c.sources,
				IsPlain:       c.plain,
				StarMerge: &StarMergeInfo{
					Table:    c.starMerge.table,
					Index:    c.starMerge.index,
					LeftStar: c.starMerge.leftStar,
				},
			}
			out = append(out, info)
		default:
			out = append(out, ColumnInfo{
				Name:          c.name,
				SourceColumns: c.sources,
				IsPlain:       c.plain,
				BaseFieldName: c.baseFieldName,
			})
		}
	}
	return out
}

// segmentsToWire converts a star item's resolved expansion segments into the
// StarSegment wire form: a baseStar element becomes a BaseTable segment (the
// consumer enumerates the table's columns), a concrete element becomes a resolved
// (Name, Sources, Plain) segment the consumer emits directly.
func segmentsToWire(segs []projColumn) []StarSegment {
	out := make([]StarSegment, 0, len(segs))
	for _, s := range segs {
		if s.baseStar != nil {
			ref := *s.baseStar
			out = append(out, StarSegment{BaseTable: &ref, ExceptColumns: s.baseStarExcept, Plain: s.plain})
			continue
		}
		out = append(out, StarSegment{Name: s.name, Sources: s.sources, Plain: s.plain, BaseFieldName: s.baseFieldName})
	}
	return out
}

// collectColumnPaths returns the raw dotted-name component lists of every column
// reference directly mentioned in expr (in source order, excluding nested
// subqueries), mirroring the legacy getPossibleColumnResources. The flat parts —
// not a dialect-bucketed ColumnRef — are what relation-aware resolution needs (a
// 2-part `rel.col` must be tried as relation.column, not dataset.column).
func (w *spanWalker) collectColumnPaths(expr ast.Node) [][]string {
	if expr == nil {
		return nil
	}
	var paths [][]string
	ew := &exprWalk{
		w:         w,
		followSub: false,
		onParts:   func(parts []string) { paths = append(paths, parts) },
	}
	ew.walk(expr)
	return paths
}

// exprToParts renders a dotted-name expression (Identifier / PathExpr /
// FieldAccess chain) into its component parts, or nil if it is not a plain name
// chain.
func exprToParts(expr ast.Node) []string {
	switch e := expr.(type) {
	case *ast.Identifier:
		if e.Name == "" {
			return nil
		}
		return []string{e.Name}
	case *ast.PathExpr:
		if len(e.Parts) == 0 {
			return nil
		}
		return append([]string{}, e.Parts...)
	case *ast.FieldAccess:
		return flattenFieldAccess(e)
	case *ast.ParenExpr:
		return exprToParts(e.Expr)
	}
	return nil
}
