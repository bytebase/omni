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
	// The FK-implicit exclusion is computed against the OTHER side's USER-managed index names
	// (its index names minus its own FK-implicit ones). An FK-implicit index is dropped from a
	// side only when the other side has no USER index of that name. This keeps the exclusion
	// symmetric for an index that is a genuine user index on at least one side (so it is never
	// spuriously added/dropped just because a foreign key appeared/disappeared), while still
	// suppressing a backing index that is FK-implicit on BOTH sides (a FK whose columns changed —
	// owned by the FK node) or that rides with a one-sided FK add/drop. See diffableIndexMap.
	fromMap := diffableIndexMap(from, userIndexNameSet(to))
	toMap := diffableIndexMap(to, userIndexNameSet(from))

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
// diffableColumns dropping the GIPK column) and the FK-implicit-backing indexes the FK node owns.
// The result is the set of indexes whose lifecycle this node owns.
//
// otherSideUserNames is the lower-cased set of the OTHER catalog's USER-managed index names (its
// index names minus its own FK-implicit ones). An FK-implicit index is excluded ONLY when the
// other side has no USER index of that name. This is the firewall against two failure modes:
//   - asymmetric exclusion (errno 1061, duplicate key name): a user `KEY my_idx (pid)` that backs
//     a FK on one side and is a standalone index on the other (FK dropped, index kept) is a USER
//     index on that other side, so it is NOT excluded and is correctly seen as unchanged;
//   - dropping a FK-needed index (errno 1553): a backing index that is FK-implicit on BOTH sides
//     (e.g. a FK whose columns changed, so the backing index columns changed but the name stayed)
//     has no USER index of that name on the other side, so it stays excluded and the FK node owns
//     the change.
//
// A backing index that genuinely appears/disappears WITH a one-sided FK is absent from the other
// side entirely, so it is still excluded. Pass nil to exclude all FK-implicit indexes
// unconditionally (single-table contexts: orderedDiffableIndexes for a new table, no other side).
func diffableIndexMap(t *Table, otherSideUserNames map[string]bool) map[string]*Index {
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
		name := toLower(idx.Name)
		// Exclude an FK-implicit index unless the other side has a USER index of the same name.
		if skip[name] && !otherSideUserNames[name] {
			continue
		}
		m[name] = idx
	}
	return m
}

// userIndexNameSet returns the lower-cased set of a table's USER-managed index names: every index
// name except the generated-invisible-primary-key (never user-authored) and the FK-implicit
// backing indexes (owned by the FK node). It is the cross-side reference for the FK-implicit
// exclusion in diffableIndexMap — "does the other side carry a genuine user index of this name?".
func userIndexNameSet(t *Table) map[string]bool {
	s := make(map[string]bool)
	if t == nil {
		return s
	}
	skip := fkImplicitIndexNames(t)
	for _, idx := range t.Indexes {
		if idx == nil || isGeneratedInvisiblePrimaryKeyIndex(idx) {
			continue
		}
		name := toLower(idx.Name)
		if skip[name] {
			continue
		}
		s[name] = true
	}
	return s
}

// fkImplicitIndexNames returns the lower-cased names of indexes that exist solely to back a
// foreign key — the index MySQL auto-creates (and the loader synthesizes via
// ensureFKBackingIndex) when no user index already covers the FK columns. These are owned by
// the FK breadth node, not by this node, so they are excluded from the index diff. Excluding
// them is symmetric (applied to both sides), so it never breaks idempotence; it prevents the
// index node from emitting an ADD/DROP for an index whose lifecycle the FK node manages — which
// would duplicate the FK node's work (duplicate-key-name on ADD) or fail with errno 1553
// ("needed in a foreign key constraint") on DROP.
//
// Detection combines STRUCTURE and NAME. An index is FK-implicit only when BOTH hold:
//   - it has exactly the shape ensureFKBackingIndex produces — a plain (non-unique/fulltext/
//     spatial), default-type, visible, comment-less, block-size-less index whose column parts are
//     the FK's columns verbatim (no prefix/expression/DESC) — isPlainBackingIndexFor; AND
//   - its name is the AUTO-DERIVED backing name for that FK — isAutoBackingIndexName: the
//     constraint name, the first FK column, or a `<firstFKcol>_<digits>` allocIndexName collision
//     suffix (e.g. `pid_2`). MySQL/the loader names the backing index this way.
//
// The name check is what distinguishes a user's EXPLICITLY-named index from the synthesized one.
// If a user writes `KEY idx_pid (pid)` alongside a `CONSTRAINT fk_c_p FOREIGN KEY (pid)`, the
// index name `idx_pid` is neither the constraint name nor a first-column-derived name, so it is
// NOT FK-implicit — it is a user index this node must emit (else, when the FK is added, the FK
// node would auto-create a differently-named backing index and the user's chosen name `idx_pid`
// would be lost). When the user's index name DOES match the auto-derived form (the constraint
// name, or the column-default name MySQL would pick anyway), excluding it is correct — the FK
// node recreates an identically-named index, and it could not be independently dropped while the
// FK exists (errno 1553). The structural part still keeps a UNIQUE/prefixed/wider index visible.
//
// Each FK claims at most one index (the first match not already claimed by another FK), mirroring
// the one-backing-index-per-FK the engine maintains.
func fkImplicitIndexNames(t *Table) map[string]bool {
	skip := make(map[string]bool)
	if t == nil {
		return skip
	}
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConForeignKey || len(con.Columns) == 0 {
			continue
		}
		for _, idx := range t.Indexes {
			if idx == nil {
				continue
			}
			key := toLower(idx.Name)
			if skip[key] {
				continue
			}
			if isPlainBackingIndexFor(idx, con.Columns) && isAutoBackingIndexName(idx.Name, con) {
				skip[key] = true
				break
			}
		}
	}
	return skip
}

// isAutoBackingIndexName reports whether idxName is the name MySQL/ensureFKBackingIndex would
// assign to the index backing con: the constraint name, the first FK column, or a
// `<firstFKcol>_<digits>` collision suffix that allocIndexName produces when the first-column name
// is already taken. A user index named anything else (e.g. `idx_pid` for an FK on `pid`) is a
// deliberately-named user index and is not treated as the implicit backing index, so this node
// emits it and the FK reuses it — preserving the user's chosen name.
func isAutoBackingIndexName(idxName string, con *Constraint) bool {
	name := toLower(idxName)
	if con.Name != "" && name == toLower(con.Name) {
		return true
	}
	if len(con.Columns) == 0 {
		return false
	}
	first := toLower(con.Columns[0])
	if name == first {
		return true
	}
	// allocIndexName collision form: "<firstcol>_<digits>".
	if rest, ok := strings.CutPrefix(name, first+"_"); ok && rest != "" {
		for _, r := range rest {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	return false
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
