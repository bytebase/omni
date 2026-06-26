package catalog

import (
	"sort"
	"strconv"
	"strings"
)

// MySQL SDL diff — indexes (the full index/key surface).
//
// This is the MySQL analog of PG's diffIndexes (pg/catalog/diff_index.go). It compares the
// index set of two versions of a table and reports the per-index changes via
// TableDiffEntry.Indexes. The MySQL model keeps the WHOLE key surface in Table.Indexes:
// PRIMARY KEY (Primary), UNIQUE (Unique, MySQL's unique constraint IS a unique index),
// secondary (plain) KEYs, FULLTEXT, and SPATIAL — plus the index that MySQL auto-creates to
// back a foreign key. show.go renders all of these via the index path (PRIMARY KEY / UNIQUE
// KEY / KEY / FULLTEXT KEY / SPATIAL KEY); only FOREIGN KEY and CHECK are rendered via the
// constraint path. So this node owns Table.Indexes in full, and the constraint breadth node
// owns only FK + CHECK — the same split show.go uses (show.go:121).
//
// Identity is the index name, lower-cased (PRIMARY for the primary key). Equality is decided
// by a canonical key (canonicalIndex) that folds every stored-form-significant attribute —
// kind, column parts (name/prefix/expr, version-flagged DESC via the normalizer), USING type,
// comment, visibility, and the 5.7-only index KEY_BLOCK_SIZE — so an index whose surface form
// differs from its SHOW CREATE readback but whose stored form is identical produces no phantom
// diff. Every canonicalization-sensitive part is routed through normalize.go (the normalizer's
// CanonicalIndexColumn handles prefix + version-flagged DESC; CanonicalGeneratedExpr handles a
// functional index expression).
//
// FK-implicit-backing indexes are EXCLUDED (see fkImplicitIndexNames): MySQL auto-creates an
// index to back a foreign key when no user index covers the FK columns, and SHOW CREATE shows
// both the KEY and the CONSTRAINT. The loader synthesizes that same backing index on the user
// side (tablecmds.go ensureFKBackingIndex), so it is present symmetrically — but its lifecycle
// is bound to the foreign key, which the FK breadth node owns. Emitting an ADD/DROP for it here
// would (a) duplicate the FK node's work and (b) risk a duplicate-key-name or "needed in a
// foreign key constraint" (errno 1553) failure. So an index that exists solely to back a FK is
// not a user-managed index change and is dropped from both sides of the diff.
func diffIndexes(from, to *Table, n *Normalizer) []IndexDiffEntry {
	fromMap := diffableIndexMap(from)
	toMap := diffableIndexMap(to)

	var result []IndexDiffEntry

	// Dropped: in from but not in to.
	for name, fromIdx := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, IndexDiffEntry{
				Action: DiffDrop,
				Name:   fromIdx.Name,
				From:   fromIdx,
			})
		}
	}

	// Added or modified: in to.
	for name, toIdx := range toMap {
		fromIdx, ok := fromMap[name]
		if !ok {
			result = append(result, IndexDiffEntry{
				Action: DiffAdd,
				Name:   toIdx.Name,
				To:     toIdx,
			})
			continue
		}
		if indexesChanged(fromIdx, toIdx, n) {
			result = append(result, IndexDiffEntry{
				Action: DiffModify,
				Name:   toIdx.Name,
				From:   fromIdx,
				To:     toIdx,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if a, b := toLower(result[i].Name), toLower(result[j].Name); a != b {
			return a < b
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// diffableIndexMap indexes a table's user-managed indexes by lower-cased name, excluding the
// generated-invisible-primary-key index (engine-synthesized, never user-authored — mirrors
// diffableColumns dropping the GIPK column) and every FK-implicit-backing index. The result is
// the set of indexes whose lifecycle this node owns.
func diffableIndexMap(t *Table) map[string]*Index {
	if t == nil {
		return map[string]*Index{}
	}
	skip := fkImplicitIndexNames(t)
	m := make(map[string]*Index, len(t.Indexes))
	for _, idx := range t.Indexes {
		if idx == nil {
			continue
		}
		if isGeneratedInvisiblePrimaryKeyIndex(idx) {
			continue
		}
		if skip[toLower(idx.Name)] {
			continue
		}
		m[toLower(idx.Name)] = idx
	}
	return m
}

// fkImplicitIndexNames returns the lower-cased names of indexes that exist solely to back a
// foreign key — the ones MySQL auto-creates (and the loader synthesizes via
// ensureFKBackingIndex) when no user index already covers the FK columns. These are owned by
// the FK breadth node, not by this node, so they are excluded from the index diff.
//
// Detection mirrors ensureFKBackingIndex (tablecmds.go) exactly, so it classifies precisely
// the indexes that function produced and never a user-declared index:
//   - For each FK constraint on the table, the backing index it would create has the FK's
//     index columns and a name of either the constraint name (when the FK was named) or the
//     first column's auto-allocated name (when the FK was unnamed, matching allocIndexName).
//     MySQL stores an unnamed FK as `<table>_ibfk_N`; for such a constraint the backing index
//     is named after the first column, so we resolve the name the same way.
//   - The candidate index must (a) carry only the plain attributes a synthesized backing index
//     has — non-unique, non-fulltext, non-spatial, default USING type, no comment, visible, no
//     prefix/expr/DESC parts, no KEY_BLOCK_SIZE — and (b) have column parts exactly equal to
//     the FK's columns. A UNIQUE/named/prefixed index a user wrote on the FK columns therefore
//     stays user-managed (it does not match the plain-backing signature) — and MySQL would have
//     reused it rather than creating an implicit one, so there is no implicit index to skip.
//
// When more than one FK could claim the same index, the first match wins (the index is skipped
// once); the per-FK ownership is the FK node's concern.
func fkImplicitIndexNames(t *Table) map[string]bool {
	skip := make(map[string]bool)
	if t == nil {
		return skip
	}
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConForeignKey {
			continue
		}
		backingName := fkBackingIndexName(t, con)
		if backingName == "" {
			continue
		}
		key := toLower(backingName)
		if skip[key] {
			continue
		}
		idx := findIndexByName(t, backingName)
		if idx == nil {
			continue
		}
		if isPlainBackingIndexFor(idx, con.Columns) {
			skip[key] = true
		}
	}
	return skip
}

// fkBackingIndexName returns the name MySQL/the loader uses for the index backing a foreign
// key: the constraint's IndexName if recorded, else the constraint name when it is a
// user-declared name, else the first FK column (the allocIndexName fallback for an unnamed FK,
// which MySQL records as `<table>_ibfk_N`). It returns "" for a column-less FK (defensive).
func fkBackingIndexName(t *Table, con *Constraint) string {
	if con.IndexName != "" {
		return con.IndexName
	}
	if len(con.Columns) == 0 {
		return ""
	}
	// An unnamed FK is stored as `<table>_ibfk_N`; its backing index is named after the first
	// column (allocIndexName), not the constraint. A user-named FK names its backing index
	// after the constraint. Distinguish by the auto-generated `_ibfk_` shape.
	if con.Name == "" || isAutoFKName(t.Name, con.Name) {
		return con.Columns[0]
	}
	return con.Name
}

// isAutoFKName reports whether name is an auto-generated InnoDB FK constraint name of the form
// `<table>_ibfk_<digits>` (the shape MySQL assigns an unnamed FK; see tablecmds.go
// nextFKGeneratedNumber). Comparison is case-insensitive on the table prefix to match the
// catalog's case-insensitive identifier handling.
func isAutoFKName(tableName, conName string) bool {
	prefix := toLower(tableName) + "_ibfk_"
	low := toLower(conName)
	if !strings.HasPrefix(low, prefix) {
		return false
	}
	digits := low[len(prefix):]
	if digits == "" {
		return false
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// isPlainBackingIndexFor reports whether idx has exactly the shape ensureFKBackingIndex
// produces for an FK on fkCols: a plain (non-unique, non-fulltext, non-spatial), default-type,
// visible, comment-less, block-size-less index whose column parts are the FK columns verbatim
// (no prefix length, no expression, no DESC). Anything richer is a user-managed index that
// merely happens to cover the FK columns, and is diffed normally.
func isPlainBackingIndexFor(idx *Index, fkCols []string) bool {
	if idx.Primary || idx.Unique || idx.Fulltext || idx.Spatial {
		return false
	}
	if idx.IndexType != "" || idx.Comment != "" || idx.KeyBlockSize != 0 || !idx.Visible {
		return false
	}
	if len(idx.Columns) != len(fkCols) {
		return false
	}
	for i, ic := range idx.Columns {
		if ic.Expr != "" || ic.Length != 0 || ic.Descending {
			return false
		}
		if !strings.EqualFold(ic.Name, fkCols[i]) {
			return false
		}
	}
	return true
}

// findIndexByName returns the table's index with the given name (case-insensitive), or nil.
func findIndexByName(t *Table, name string) *Index {
	key := toLower(name)
	for _, idx := range t.Indexes {
		if idx != nil && toLower(idx.Name) == key {
			return idx
		}
	}
	return nil
}

// indexesChanged reports whether two same-name indexes differ, comparing their canonical keys
// (the MySQL analog of PG's indexesChanged). canonicalIndex folds every stored-form-significant
// attribute into one collision-free key, so this differ never re-implements a per-attribute
// comparison.
func indexesChanged(a, b *Index, n *Normalizer) bool {
	return canonicalIndex(a, n) != canonicalIndex(b, n)
}

// canonicalIndex returns a single stable comparison key for an index, folding the attributes
// that survive into MySQL's stored form for the normalizer's version. Routed entirely through
// normalize-core for the canonicalization-sensitive parts:
//   - column parts via CanonicalIndexColumn (name/prefix/expr + version-flagged DESC: 8.0
//     keeps DESC, 5.7 silently stores ascending);
//   - KEY_BLOCK_SIZE is version-flagged the same way DESC is: 5.7 echoes an index-level
//     KEY_BLOCK_SIZE in SHOW CREATE, 8.0 silently absorbs it (never stored on the index, never
//     echoed — verified against both live engines), so it contributes to the key only on 5.7.
//     This is canonicalIndexKeyBlockSize; it keeps an 8.0 user form that wrote KEY_BLOCK_SIZE=N
//     from phantom-diffing against the 8.0 readback that drops it.
//
// Kind (primary/unique/fulltext/spatial), USING type, comment, and visibility are
// version-independent and compared directly. The index NAME is the identity key (handled by the
// caller's map), not part of this content key.
func canonicalIndex(idx *Index, n *Normalizer) string {
	parts := make([]string, 0, len(idx.Columns))
	for _, ic := range idx.Columns {
		parts = append(parts, n.CanonicalIndexColumn(ic))
	}
	return encodeKeyFields(
		"primary", strconv.FormatBool(idx.Primary),
		"unique", strconv.FormatBool(idx.Unique),
		"fulltext", strconv.FormatBool(idx.Fulltext),
		"spatial", strconv.FormatBool(idx.Spatial),
		"using", canonicalIndexType(idx),
		"comment", n.CanonicalComment(idx.Comment),
		"visible", strconv.FormatBool(idx.Visible),
		"kbs", canonicalIndexKeyBlockSize(idx, n),
		"cols", strings.Join(parts, ","),
	)
}

// canonicalIndexType returns the index's USING type, upper-cased, but only when it is a
// user-selectable BTREE/HASH choice. FULLTEXT and SPATIAL set IndexType to "FULLTEXT"/"SPATIAL"
// as a redundant echo of the kind flags (already in the key), and SHOW CREATE never emits a
// USING clause for them, so they are normalized out here to avoid a key that double-counts the
// kind. PRIMARY likewise never carries a meaningful USING in the stored form.
func canonicalIndexType(idx *Index) string {
	if idx.Primary || idx.Fulltext || idx.Spatial {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(idx.IndexType))
}

// canonicalIndexKeyBlockSize returns the index KEY_BLOCK_SIZE as a key fragment, but only on
// 5.7 where SHOW CREATE echoes it. On 8.0 the engine does not store an index-level
// KEY_BLOCK_SIZE (information_schema shows none and SHOW CREATE omits it), so it is dropped from
// the key — otherwise an 8.0 `KEY k (a) KEY_BLOCK_SIZE=4` user form would phantom-diff against
// its readback, which has no block size. FLAG: this version split is local to index diffing; if
// normalize-core later grows a CanonicalIndexKeyBlockSize, route through it.
func canonicalIndexKeyBlockSize(idx *Index, n *Normalizer) string {
	if idx.KeyBlockSize == 0 || n.is80() {
		return ""
	}
	return strconv.Itoa(idx.KeyBlockSize)
}
