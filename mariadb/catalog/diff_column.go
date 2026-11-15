package catalog

import "sort"

// diffColumns compares the columns of two versions of a table. Identity is the
// lower-cased column name; two same-name columns are equal when the Normalizer
// renders them alike (CanonicalColumn). Hidden columns are not compared (see
// diffableColumns).
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

// columnsChanged reports whether two same-name columns differ by CanonicalColumn.
// Position is not compared: a reorder is not a change to the column. The tables
// give CanonicalColumn the context it needs, such as the table's column checks.
func columnsChanged(fromTbl, toTbl *Table, a, b *Column, n *Normalizer) bool {
	return n.CanonicalColumn(fromTbl, a) != n.CanonicalColumn(toTbl, b)
}

// diffableColumns returns the columns of a table that take part in the diff,
// leaving out hidden columns the server maintains itself.
func diffableColumns(t *Table) []*Column {
	cols := make([]*Column, 0, len(t.Columns))
	for _, col := range t.Columns {
		if col.Hidden != ColumnHiddenNone {
			continue
		}
		cols = append(cols, col)
	}
	return cols
}
