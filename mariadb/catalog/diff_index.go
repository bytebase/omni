package catalog

import (
	"sort"
	"strings"
)

// diffIndexes compares the keys of two versions of a table: the primary key,
// unique keys, and plain, FULLTEXT and SPATIAL keys. Identity is the lower-cased
// key name (PRIMARY for the primary key); two same-name keys are equal when the
// Normalizer renders them alike (CanonicalIndex).
//
// An index that exists only to back a foreign key is left out of both sides
// (fkImplicitIndexNames): its lifecycle follows the foreign key, and adding or
// dropping it separately fails with a duplicate key name or with errno 1553.
func diffIndexes(from, to *Table, n *Normalizer) []IndexDiffEntry {
	// An FK-implicit index is left out of a side unless the other side has a user
	// index of that name, so an index that is user-managed on one side is not reported
	// as added or dropped when a foreign key appears or disappears.
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

// diffableIndexMap indexes a table's indexes by lower-cased name, leaving out the
// FK-implicit ones unless the other catalog has a user index of the same name
// (otherSideUserNames). That keeps an index that backs a key on one side and
// stands alone on the other from being reported as added or dropped (errno 1061 on
// apply), and keeps a backing index whose key changed out of this differ (errno
// 1553). A nil set leaves out every FK-implicit index.
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
		name := toLower(idx.Name)
		if skip[name] && !otherSideUserNames[name] {
			continue
		}
		m[name] = idx
	}
	return m
}

// userIndexNameSet returns the lower-cased names of a table's indexes that are
// not FK-implicit.
func userIndexNameSet(t *Table) map[string]bool {
	s := make(map[string]bool)
	if t == nil {
		return s
	}
	skip := fkImplicitIndexNames(t)
	for _, idx := range t.Indexes {
		if idx == nil {
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

// fkImplicitIndexNames returns the lower-cased names of the indexes that exist
// only to back a foreign key. An index is FK-implicit when it has the shape
// ensureFKBackingIndex creates (isPlainBackingIndexFor) and a name that function
// gives (isAutoBackingIndexName), so an index the user named explicitly on the
// same columns stays a user index. Each key claims at most one index.
func fkImplicitIndexNames(t *Table) map[string]bool {
	skip := make(map[string]bool)
	for _, idx := range fkImplicitBackingIndexes(t) {
		skip[toLower(idx.Name)] = true
	}
	return skip
}

// fkImplicitBackingIndexes maps each foreign key of t to the FK-implicit index it claims (see
// fkImplicitIndexNames). A foreign key that claims none is absent.
func fkImplicitBackingIndexes(t *Table) map[*Constraint]*Index {
	claims := make(map[*Constraint]*Index)
	if t == nil {
		return claims
	}
	claimed := make(map[string]bool)
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConForeignKey || len(con.Columns) == 0 {
			continue
		}
		for _, idx := range t.Indexes {
			if idx == nil {
				continue
			}
			key := toLower(idx.Name)
			if claimed[key] {
				continue
			}
			if isPlainBackingIndexFor(idx, con.Columns) && isAutoBackingIndexName(idx.Name, con) {
				claimed[key] = true
				claims[con] = idx
				break
			}
		}
	}
	return claims
}

// isAutoBackingIndexName reports whether idxName is a name ensureFKBackingIndex
// gives the index that backs con: the constraint name, the first key column, or
// `<firstcol>_<digits>` from allocIndexName when that name is taken. An index with
// any other name was named by the user and is not the backing index.
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

// indexesChanged reports whether two same-name indexes differ by CanonicalIndex.
func indexesChanged(a, b *Index, n *Normalizer) bool {
	return canonicalIndex(a, n) != canonicalIndex(b, n)
}

// canonicalIndex is the comparison key of an index: the text SHOW CREATE TABLE
// gives it.
func canonicalIndex(idx *Index, n *Normalizer) string {
	return n.CanonicalIndex(idx)
}
