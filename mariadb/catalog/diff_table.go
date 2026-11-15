package catalog

import (
	"sort"
	"strings"
)

// tableKey is the identity key of a table: (database, table), lower-cased.
type tableKey struct {
	db   string
	name string
}

// diffTables compares the base tables of two catalogs and reports those added,
// dropped and modified, sorted by (database, name, action).
func diffTables(from, to *Catalog, n *Normalizer) []TableDiffEntry {
	fromMap := buildTableMap(from)
	toMap := buildTableMap(to)

	var result []TableDiffEntry

	// Dropped: in from but not in to.
	for key, fromTbl := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, TableDiffEntry{
				Action:   DiffDrop,
				Database: key.db,
				Name:     key.name,
				From:     fromTbl,
			})
		}
	}

	// Added or modified: in to.
	for key, toTbl := range toMap {
		fromTbl, ok := fromMap[key]
		if !ok {
			result = append(result, TableDiffEntry{
				Action:   DiffAdd,
				Database: key.db,
				Name:     key.name,
				To:       toTbl,
			})
			continue
		}
		if entry, changed := compareTable(key, fromTbl, toTbl, n); changed {
			result = append(result, entry)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Database != result[j].Database {
			return result[i].Database < result[j].Database
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// buildTableMap indexes every base table of a catalog by (database, name).
// Temporary tables are left out: they are session-scoped and not part of a schema
// snapshot.
func buildTableMap(c *Catalog) map[tableKey]*Table {
	m := make(map[tableKey]*Table)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, tbl := range db.Tables {
			if tbl == nil || tbl.Temporary {
				continue
			}
			m[tableKey{db: toLower(db.Name), name: name}] = tbl
		}
	}
	return m
}

// compareTable reports whether two same-identity tables differ in their options
// or in a sub-object, returning the DiffModify entry when they do.
func compareTable(key tableKey, from, to *Table, n *Normalizer) (TableDiffEntry, bool) {
	entry := TableDiffEntry{
		Action:   DiffModify,
		Database: key.db,
		Name:     key.name,
		From:     from,
		To:       to,
		Columns:  diffColumns(from, to, n),
		// Sub-object diffs.
		Indexes:          diffIndexes(from, to, n),
		Constraints:      diffConstraints(from, to, n),
		ForeignKeys:      diffForeignKeys(from, to, n),
		Checks:           diffChecks(from, to, n),
		PartitionChanged: diffPartitions(from, to, n),
	}

	if !tableOptionsChanged(from, to, n) && !tableSubdiffsChanged(&entry) {
		return TableDiffEntry{}, false
	}
	return entry, true
}

// tableOptionsChanged compares the table options as SHOW CREATE TABLE gives
// them (less AUTO_INCREMENT, which is data, not schema), then system versioning
// and the system and application-time periods.
func tableOptionsChanged(from, to *Table, n *Normalizer) bool {
	return n.CanonicalTableOptions(from) != n.CanonicalTableOptions(to) ||
		from.SystemVersioned != to.SystemVersioned ||
		!strings.EqualFold(from.PeriodStartCol, to.PeriodStartCol) ||
		!strings.EqualFold(from.PeriodEndCol, to.PeriodEndCol) ||
		!strings.EqualFold(from.AppPeriodName, to.AppPeriodName) ||
		!strings.EqualFold(from.AppPeriodStartCol, to.AppPeriodStartCol) ||
		!strings.EqualFold(from.AppPeriodEndCol, to.AppPeriodEndCol)
}

// tableSubdiffsChanged reports whether any sub-object diff of the entry is
// non-empty.
func tableSubdiffsChanged(e *TableDiffEntry) bool {
	return len(e.Columns) > 0 ||
		len(e.Indexes) > 0 ||
		len(e.Constraints) > 0 ||
		len(e.ForeignKeys) > 0 ||
		len(e.Checks) > 0 ||
		e.PartitionChanged
}
