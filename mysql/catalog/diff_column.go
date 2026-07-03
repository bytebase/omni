package catalog

import "sort"

// diffColumns compares the columns of two versions of a table (the MySQL analog of PG's
// diffColumns at pg/catalog/diff_column.go:7). Identity is the column name (lower-cased)
// within the table; equality is decided by the normalize-core CanonicalColumn aggregate
// key, so a column whose surface form differs from its synced SHOW CREATE readback but
// whose canonical form is identical produces NO phantom diff.
//
// Only DIFFABLE columns are compared — see diffableColumns. MySQL synthesizes columns
// that never appear in user-authored SDL (the my_row_id generated invisible primary key
// and the hidden expression columns that back functional indexes); including them would
// manufacture a phantom column diff on every synced table that has a functional index or
// GIPK. A user-declared INVISIBLE column is genuine schema and IS diffed; only
// engine-synthesized columns are skipped.
func diffColumns(from, to *Table, n *Normalizer) []ColumnDiffEntry {
	fromCols := diffableColumns(from)
	toCols := diffableColumns(to)

	fromMap := make(map[string]*Column, len(fromCols))
	for _, col := range fromCols {
		fromMap[toLower(col.Name)] = col
	}
	toMap := make(map[string]*Column, len(toCols))
	for _, col := range toCols {
		toMap[toLower(col.Name)] = col
	}

	var result []ColumnDiffEntry

	// Dropped: in from but not in to.
	for name, fromCol := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, ColumnDiffEntry{
				Action: DiffDrop,
				Name:   fromCol.Name,
				From:   fromCol,
			})
		}
	}

	// Added or modified: in to.
	for name, toCol := range toMap {
		fromCol, ok := fromMap[name]
		if !ok {
			result = append(result, ColumnDiffEntry{
				Action: DiffAdd,
				Name:   toCol.Name,
				To:     toCol,
			})
			continue
		}
		if columnsChanged(from, to, fromCol, toCol, n) {
			result = append(result, ColumnDiffEntry{
				Action: DiffModify,
				Name:   toCol.Name,
				From:   fromCol,
				To:     toCol,
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

// columnsChanged reports whether two same-name columns differ, comparing the
// normalize-core CanonicalColumn key of each (the MySQL analog of PG's columnsChanged at
// pg/catalog/diff_column.go:66, which compares FormatType). CanonicalColumn folds every
// canonicalization-sensitive aspect into one collision-free key — type (with the
// version's integer-display-width / BOOLEAN / decimal / year rules), the resolved
// (charset, collation) pair, default and ON UPDATE (value-based, EDFT-aware), effective
// nullability (PK-forced + TIMESTAMP magic), generated-column expression and storage,
// AUTO_INCREMENT attribute, comment, invisibility, and SRID — so this differ never
// re-implements any per-aspect comparison.
//
// Column position is intentionally NOT folded into the key: a pure reorder is not a
// column-content change, and MySQL preserves declared order without reordering, so the
// synced readback and a same-order target agree. (Re-positioning support, if generate-
// core ever needs it, belongs in generate-core's ALTER ordering, not in column
// identity.) `from`/`to` are passed for the table context CanonicalColumn needs to
// resolve table-inherited charset/collation, PK membership, and the first-TIMESTAMP rule.
func columnsChanged(fromTbl, toTbl *Table, a, b *Column, n *Normalizer) bool {
	if n.CanonicalColumn(fromTbl, a) == n.CanonicalColumn(toTbl, b) {
		return false
	}
	// Cast-charset wildcard: a charset-less CHAR cast in user SDL must not
	// phantom-diff against the connection charset its readback carries, while
	// an EXPLICIT charset change stays a real diff (both engine-verified; see
	// columnsCastCharsetWildcardEqual).
	return !n.columnsCastCharsetWildcardEqual(fromTbl, toTbl, a, b)
}

// diffableColumns returns the columns of a table that participate in a declarative diff:
// the user-authored columns, excluding every column MySQL synthesizes on its own. Two
// kinds of synthesized column must be dropped, and VisibleColumns() alone is NOT enough:
//   - functional-index expression columns: marked Hidden == ColumnHiddenSystem, so
//     VisibleColumns() already excludes them; and
//   - the generated invisible primary key (my_row_id): marked GeneratedInvisiblePrimaryKey
//     and Invisible, but left Hidden == ColumnHiddenNone, so VisibleColumns() WOULD return
//     it. SHOW CREATE TABLE suppresses it via a separate check (show.go), and a synced
//     dump taken with sql_generate_invisible_primary_key=ON carries it — so comparing it
//     against a target that does not declare it produces a phantom DROP. We exclude it
//     here explicitly.
//
// A user-declared INVISIBLE column (Invisible == true, Hidden == None,
// GeneratedInvisiblePrimaryKey == false) is real schema and is kept.
func diffableColumns(t *Table) []*Column {
	cols := make([]*Column, 0, len(t.Columns))
	for _, col := range t.Columns {
		if col.Hidden != ColumnHiddenNone {
			continue
		}
		if col.GeneratedInvisiblePrimaryKey {
			continue
		}
		cols = append(cols, col)
	}
	return cols
}
