package catalog

import (
	"slices"
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
// secondary (plain) KEYs, FULLTEXT, and SPATIAL. show.go renders all of these via the index
// path (PRIMARY KEY / UNIQUE KEY / KEY / FULLTEXT KEY / SPATIAL KEY); only FOREIGN KEY and
// CHECK are rendered via the constraint path. So this node owns Table.Indexes in full, and the constraint breadth node
// owns only FK + CHECK — the same split show.go uses (show.go:121).
//
// Identity is the index name, lower-cased (PRIMARY for the primary key). Equality is decided
// by a canonical key (canonicalIndex) that folds every stored-form-significant attribute —
// kind, column parts (name/prefix/expr, DESC via the normalizer), USING type, comment, and
// visibility — so an index whose surface form differs from its SHOW CREATE readback but whose
// stored form is identical produces no phantom diff. Every canonicalization-sensitive part is
// routed through normalize.go (the normalizer's CanonicalIndexColumn handles prefix + DESC;
// CanonicalGeneratedExpr handles a functional index expression).
//
// Every index is the user's own. OceanBase makes no index for a foreign key, where MySQL
// adds one when no index serves it, so there is none here that belongs to a constraint and
// must stay out of the diff (TestOracle_OceanBaseForeignKeyIndex). What a foreign key does
// do to an index is hold the last one over exactly its columns: the foreign key generator
// lets go of it first (foreignKeysHoldingDroppedIndex).
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

// diffableIndexMap indexes a table's indexes by lower-cased name, leaving out the
// generated-invisible-primary-key index: the engine makes it and nobody writes it, as
// diffableColumns leaves out its column.
func diffableIndexMap(t *Table) map[string]*Index {
	if t == nil {
		return map[string]*Index{}
	}
	m := make(map[string]*Index, len(t.Indexes))
	for _, idx := range t.Indexes {
		if idx == nil {
			continue
		}
		m[toLower(idx.Name)] = idx
	}
	return m
}

// indexesChanged reports whether two same-name indexes differ, comparing their canonical keys
// (the MySQL analog of PG's indexesChanged). canonicalIndex folds every stored-form-significant
// attribute into one collision-free key, so this differ never re-implements a per-attribute
// comparison.
func indexesChanged(a, b *Index, n *Normalizer) bool {
	return canonicalIndex(a, n) != canonicalIndex(b, n)
}

// canonicalIndex returns a single stable comparison key for an index, folding the attributes
// that survive into MySQL's stored form. Routed entirely through normalize-core for the
// canonicalization-sensitive parts:
//   - column parts via CanonicalIndexColumn (name/prefix/expr + DESC);
//   - an index-level KEY_BLOCK_SIZE is not part of the key: 8.0 silently absorbs it (never
//     stored on the index, never echoed), so a user form that wrote KEY_BLOCK_SIZE=N must not
//     phantom-diff against the readback that drops it.
//
// Kind (primary/unique/fulltext/spatial), USING type, comment, and visibility are
// compared directly. A fulltext index's parser is compared with
// OceanBase's defaults filled in (canonicalFulltextParser), since the engine prints them
// on every fulltext index. The index NAME is the identity key (handled by the
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
		"vector", strconv.FormatBool(idx.Vector),
		"with", canonicalVectorParams(idx),
		"storing", canonicalStoring(idx),
		"using", canonicalIndexType(idx),
		"comment", n.CanonicalComment(idx.Comment),
		"visible", strconv.FormatBool(idx.Visible),
		"scope", canonicalIndexScope(idx),
		"parser", canonicalFulltextParser(idx),
		"cols", strings.Join(parts, ","),
	)
}

// canonicalStoring returns the STORING columns as a key fragment, folded and sorted: which
// columns an index carries matters and the order they were written in does not.
func canonicalStoring(idx *Index) string {
	names := make([]string, 0, len(idx.Storing))
	for _, name := range idx.Storing {
		names = append(names, strings.ToLower(name))
	}
	slices.Sort(names)
	// Each name carries its length, since a quoted name can hold the separator: STORING (`b,c`)
	// and STORING (b, c) are different columns.
	var b strings.Builder
	for _, name := range names {
		b.WriteString(strconv.Itoa(len(name)))
		b.WriteByte(':')
		b.WriteString(name)
	}
	return b.String()
}

// canonicalVectorParams returns a vector index's parameters as a key fragment, in sorted
// order: the server echoes them in the order they were written, which says nothing about
// the index.
func canonicalVectorParams(idx *Index) string {
	params := slices.Clone(idx.VectorParams)
	slices.Sort(params)
	return strings.Join(params, ",")
}

// canonicalIndexScope returns the index scope as a key fragment, with OceanBase's default
// folded away. 4.3.5 prints LOCAL on every secondary index, partitioned or not, whether or
// not the DDL asked for it, so comparing the raw value makes every table whose DDL omits a
// scope report a spurious index MODIFY against its own readback — a bogus DROP INDEX and
// ADD INDEX in the migration. LOCAL therefore carries no information and folds to "", the
// same treatment IgnoreRowFormat gives an unstated ROW_FORMAT. GLOBAL is only ever present
// because someone wrote it, so it stays and a LOCAL-to-GLOBAL change still diffs.
func canonicalIndexScope(idx *Index) string {
	if strings.EqualFold(strings.TrimSpace(idx.Scope), "LOCAL") {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(idx.Scope))
}

// canonicalIndexType returns the index's USING type as a key fragment, which on OceanBase is
// empty for every index. SHOW CREATE TABLE prints no USING clause, whether the index was
// written USING BTREE, USING HASH or with neither, so a table's readback never says which, and
// a written USING must not make the table differ from it. (information_schema.STATISTICS does
// keep HASH, which changes nothing an index does on this engine.) FULLTEXT, SPATIAL and VECTOR
// are kinds, with flags of their own in the key. TestOracle_OceanBaseIndexOrderAndType holds
// this to the engine.
func canonicalIndexType(*Index) string {
	return ""
}
