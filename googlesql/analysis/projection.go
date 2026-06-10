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
	// baseStar element.
	columns []projColumn
}

// resolveSelectProjection resolves a SELECT block to its output projection,
// mirroring the legacy extractTableSourceFromSelect: a bare `*` copies all FROM
// relations' columns; a `rel.*` copies the named relation's columns; an explicit
// item resolves its expression's columns against the relations. fromRels is the
// comma-item relation list (a join collapsed into one combined relation) used to
// expand a bare `*` in FROM order; column / `rel.*` lookups use the SELECT's leaf
// relations (w.leafRels), which keep each join side findable by name.
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
	if expr == nil {
		return "", nil
	}
	paths := w.collectColumnPaths(expr)
	if len(paths) == 0 {
		return "", nil
	}
	var sources []ColumnRef
	name := ""
	for _, parts := range paths {
		n, refs := w.resolvePath(parts)
		sources = unionColumnRefs(sources, refs)
		name = n
	}
	if len(paths) > 1 {
		name = ""
	}
	return name, sources
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
	if len(parts) == 0 {
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
	if findRelation(w.leafRels, parts[0]) != nil {
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
	for _, rel := range w.leafRels {
		if !rel.isBase {
			continue
		}
		b := rel.baseRef
		if strings.EqualFold(b.Table, ref.Table) &&
			strings.EqualFold(b.Database, ref.Database) &&
			strings.EqualFold(b.Schema, ref.Schema) &&
			strings.EqualFold(b.Catalog, ref.Catalog) {
			return true
		}
	}
	return false
}

// resolveColumn resolves a (relationQualifier, column) reference to base lineage
// against the SELECT's leaf relations (w.leafRels). relQualifier may be ""
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
		if rel := findRelation(w.leafRels, relQualifier); rel != nil && rel.isBase {
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
	if rel := soleBaseRelation(w.leafRels); rel != nil {
		ref := rel.baseRef
		ref.Column = column
		return []ColumnRef{ref}
	}
	return []ColumnRef{{Column: column}}
}

// soleBaseRelation returns the single base-table relation in rels when exactly
// one exists, else nil. With several base relations a bare column is ambiguous
// (omni cannot enumerate base-table columns metadata-free), so the caller keeps
// the bare ref and the consumer matches additively by name.
func soleBaseRelation(rels []*relation) *relation {
	var sole *relation
	for _, rel := range rels {
		if !rel.isBase {
			continue
		}
		if sole != nil {
			return nil
		}
		sole = rel
	}
	return sole
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
		rel := findRelation(w.leafRels, relQualifier)
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
	for _, rel := range w.leafRels {
		if rel.isBase {
			continue
		}
		if refs, ok := projectionColumnSources(rel, column, false /* allowBaseStar */); ok {
			return refs, true
		}
	}
	return nil, false
}

// projectionColumnSources returns the base lineage of the column named `column`
// in a relation's projection, or nil if the relation has no such named column.
//
// A concrete projection column matches by name. A baseStar element (a whole base
// table whose columns omni cannot enumerate) matches ANY column name — but ONLY
// when allowBaseStar is set: a baseStar match is sound for a QUALIFIED reference
// (`rel.col` names the relation, so `col` must be one of its columns) but NOT for
// an UNQUALIFIED reference (omni cannot know which of several base relations — or
// join arms — owns a bare `col`; that is left to the metadata-aware consumer to
// match by name).
func projectionColumnSources(rel *relation, column string, allowBaseStar bool) ([]ColumnRef, bool) {
	for _, c := range rel.columns {
		if c.baseStar != nil {
			if !allowBaseStar {
				continue
			}
			// A USING-coalesced star excludes its key columns — those are owned by
			// the coalesced concrete column (earlier in the projection), not by this
			// star segment.
			if nameInListFold(c.baseStarExcept, column) {
				continue
			}
			ref := *c.baseStar
			ref.Column = column
			return []ColumnRef{ref}, true
		}
		if strings.EqualFold(c.name, column) {
			// Found — even when its lineage is empty (e.g. a CTE column projected from
			// a literal like `1 AS n`); return found=true so the caller does NOT apply
			// the unresolved-reference fallback (which would wrongly fabricate a bare
			// column ref). Empty-but-resolved lineage is the legacy result.
			return c.sources, true
		}
	}
	return nil, false
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
	// `rel.*`: the first part names a leaf relation (legacy tries the leading
	// identifier as a relation for a `a.*`).
	if rel := findRelation(w.leafRels, parts[0]); rel != nil {
		return rel.columns
	}
	// `schema.table.*` (or `db.schema.table.*`): the trailing part names a base
	// relation and the written prefix must agree with its schema qualifier.
	if len(parts) >= 2 {
		last := parts[len(parts)-1]
		if rel := findRelation(w.leafRels, last); rel != nil && rel.isBase {
			schemaPart := parts[len(parts)-2]
			if rel.baseRef.Schema == "" || strings.EqualFold(rel.baseRef.Schema, schemaPart) {
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
		if len(left) != len(right) {
			return left
		}
		merged := make([]projColumn, len(left))
		for i := range left {
			merged[i] = projColumn{
				name:    left[i].name,
				sources: unionColumnRefs(left[i].sources, right[i].sources),
				plain:   false,
			}
		}
		return merged
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
		out := make([]projColumn, 0, len(matchColumns))
		for _, m := range matchColumns {
			out = append(out, projColumn{
				name:    m,
				sources: unionColumnRefs(namedSources(left, m), namedSources(right, m)),
				plain:   false,
			})
		}
		return out
	}
	out := make([]projColumn, 0, len(left)+len(right))
	for _, l := range left {
		out = append(out, projColumn{
			name:    l.name,
			sources: unionColumnRefs(l.sources, namedSources(right, l.name)),
			plain:   false,
		})
	}
	for _, r := range right {
		if !projHasName(left, r.name) {
			out = append(out, projColumn{name: r.name, sources: r.sources, plain: false})
		}
	}
	return out
}

// namedSources returns the union of the sources of every concrete projection
// column whose name matches `name` case-insensitively.
func namedSources(cols []projColumn, name string) []ColumnRef {
	var out []ColumnRef
	for _, c := range cols {
		if strings.EqualFold(c.name, name) {
			out = unionColumnRefs(out, c.sources)
		}
	}
	return out
}

// projHasName reports whether any projection column is named `name`
// (case-insensitive).
func projHasName(cols []projColumn, name string) bool {
	for _, c := range cols {
		if strings.EqualFold(c.name, name) {
			return true
		}
	}
	return false
}

// coalesceUsingJoin builds the projection of `left JOIN right USING (keys)`
// (fix F1): each key ONCE, first, in USING order — its lineage the union of
// BOTH sides' key columns — then the left side's non-key columns (rewrapped
// non-plain, the legacy join-left rule), then the right side's (plainness
// kept). It returns ok=false when a side cannot be name-partitioned without
// metadata (a deferred set-op marker in its projection); the caller fails
// closed.
func coalesceUsingJoin(left, right *relation, keys []string) (*relation, bool) {
	leftRest, ok := stripUsingKeys(left.columns, keys)
	if !ok {
		return nil, false
	}
	rightRest, ok := stripUsingKeys(right.columns, keys)
	if !ok {
		return nil, false
	}
	cols := make([]projColumn, 0, len(keys)+len(leftRest)+len(rightRest))
	for _, key := range keys {
		cols = append(cols, projColumn{
			name:          key,
			sources:       unionColumnRefs(usingKeyLineage(left.columns, key), usingKeyLineage(right.columns, key)),
			plain:         false,
			baseFieldName: true,
		})
	}
	for _, c := range leftRest {
		cols = append(cols, rewrapNonPlain(c))
	}
	cols = append(cols, rightRest...)
	return &relation{columns: cols}, true
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
		case nameMatchesAnyFold(c.name, keys):
			// Dropped: this side's key column is folded into the coalesced key.
		default:
			out = append(out, c)
		}
	}
	return out, true
}

// usingKeyLineage returns the lineage a join side contributes to the USING key
// `key`: every concrete column named like it, every base-table star that may
// own it (the star's table with Column=key — when several stars could own the
// key, all contribute: over-attribution is safe, under is not), and a star
// group's segments likewise (its EXCEPT removing the key means the side does
// not project it; a REPLACE re-points it to the replacement's sources).
func usingKeyLineage(cols []projColumn, key string) []ColumnRef {
	var out []ColumnRef
	for _, c := range cols {
		switch {
		case c.starMerge != nil || c.setOpMerge != nil:
			// Unreachable behind stripUsingKeys' ok=false; contribute nothing.
		case c.starGroup != nil:
			if nameInListFold(c.starGroup.except, key) {
				continue
			}
			if srcs, ok := replaceSourcesFor(c.starGroup.replace, key); ok {
				out = unionColumnRefs(out, srcs)
				continue
			}
			out = unionColumnRefs(out, usingKeyLineage(c.starGroup.segments, key))
		case c.baseStar != nil:
			if nameInListFold(c.baseStarExcept, key) {
				continue
			}
			ref := *c.baseStar
			ref.Column = key
			out = unionColumnRefs(out, []ColumnRef{ref})
		case strings.EqualFold(c.name, key):
			out = unionColumnRefs(out, c.sources)
		}
	}
	return out
}

// replaceSourcesFor returns the replacement sources of the REPLACE item named
// `name` (case-insensitive), if any.
func replaceSourcesFor(replace []StarReplaceItem, name string) ([]ColumnRef, bool) {
	for _, r := range replace {
		if strings.EqualFold(r.Name, name) {
			return r.Sources, true
		}
	}
	return nil, false
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

// nameInListFold reports whether list contains name under case-insensitive
// comparison.
func nameInListFold(list []string, name string) bool {
	for _, e := range list {
		if strings.EqualFold(e, name) {
			return true
		}
	}
	return false
}

// nameMatchesAnyFold reports whether name matches any of names
// (case-insensitive).
func nameMatchesAnyFold(name string, names []string) bool {
	for _, n := range names {
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}

// appendNamesFold returns a COPY of list with each name appended unless already
// present (case-insensitive). The copy keeps the caller's slice unaliased (a
// join side's projection is shared with its leaf relation).
func appendNamesFold(list, names []string) []string {
	out := append([]string{}, list...)
	for _, n := range names {
		if !nameInListFold(out, n) {
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

// findRelation returns the FROM relation whose reference name matches `name`
// (case-insensitive, GoogleSQL identifier rule), or nil.
func findRelation(rels []*relation, name string) *relation {
	for _, rel := range rels {
		if rel.name != "" && strings.EqualFold(rel.name, name) {
			return rel
		}
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
