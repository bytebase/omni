package catalog

import (
	"sort"
	"strings"
)

// tableKey is the identity key for a table: (database, table) lower-cased. MySQL has
// no schemas/namespaces, so the database is the qualifying scope.
type tableKey struct {
	db   string
	name string
}

// diffTables compares base tables between two catalogs (the MySQL analog of PG's
// diffRelations at pg/catalog/diff_table.go:13). It reports tables added (in `to`
// only), dropped (in `from` only), and modified (present in both but differing in
// table-level options or columns). Determinism: results are sorted by (database,
// name, action) so the diff — and the DDL generate-core derives from it — is stable.
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

// buildTableMap indexes every base table in a catalog by (database, name). Temporary
// tables are excluded: they are session-scoped and never part of a persisted schema
// snapshot, so they must not appear in a declarative diff.
func buildTableMap(c *Catalog) map[tableKey]*Table {
	m := make(map[tableKey]*Table)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, tbl := range db.Tables {
			// Guard a nil entry (defensive — the loader never stores nil) and skip
			// temporary tables (session-scoped, never part of a persisted snapshot).
			if tbl == nil || tbl.Temporary {
				continue
			}
			m[tableKey{db: toLower(db.Name), name: name}] = tbl
		}
	}
	return m
}

// compareTable reports whether two same-identity tables differ in table-level options
// or sub-objects, returning a fully populated DiffModify entry when they do. It is the
// MySQL analog of PG's compareRelation. Every option comparison that MySQL rewrites on
// storage is routed through normalize.go. The column sub-diff is the dominant signal;
// the breadth sub-object slices (Indexes/Constraints/ForeignKeys/Checks/PartitionChanged)
// are empty in diff-core but are folded into the "is this table modified?" decision via
// tableSubdiffsChanged, so a breadth node that populates one of them needs no change here.
func compareTable(key tableKey, from, to *Table, n *Normalizer) (TableDiffEntry, bool) {
	entry := TableDiffEntry{
		Action:   DiffModify,
		Database: key.db,
		Name:     key.name,
		From:     from,
		To:       to,
		Columns:  diffColumns(from, to, n),
		// Table-level sub-object diffs, each owned by a breadth node in its own differ (mirroring
		// PG's compareRelation, which calls diffConstraints/diffIndexes/... inline). compareTable
		// and tableSubdiffsChanged fold their results into the "is this table modified?" decision.
		// diffConstraints alone stays a deliberate no-op (PK/UNIQUE are folded into diffIndexes;
		// omni:constraints deferred).
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

// tableOptionsChanged compares the table-level attributes that survive into MySQL's
// stored form, each through normalize-core:
//   - storage engine, resolved to the server default (CanonicalEngine);
//   - the effective (charset, collation) pair, folded to a version-independent identity
//     so utf8/utf8mb3 and the version-divergent utf8mb4 default collation do not
//     phantom-diff (foldCharset + effectiveTableCollation, the same helpers the column
//     resolver uses);
//   - the table COMMENT content (CanonicalComment);
//   - ROW_FORMAT, but ONLY when normalize-core says it is significant: IgnoreRowFormat
//     returns true for an empty/DEFAULT format (environment-dependent, not reconstructible
//     from DDL) and false for an explicit non-default value — so a real DYNAMIC↔COMPRESSED
//     change is detected while a phantom default↔unspecified change is not.
//
// AUTO_INCREMENT=N is a live next-value counter, never schema, so it is unconditionally
// not compared (the constant-true IgnoreTableAutoIncrement documents that choice).
func tableOptionsChanged(from, to *Table, n *Normalizer) bool {
	if n.CanonicalEngine(from) != n.CanonicalEngine(to) {
		return true
	}
	if foldCharset(from.Charset) != foldCharset(to.Charset) {
		return true
	}
	if n.effectiveTableCollation(from) != n.effectiveTableCollation(to) {
		return true
	}
	if n.CanonicalComment(from.Comment) != n.CanonicalComment(to.Comment) {
		return true
	}
	if rowFormatChanged(from, to, n) {
		return true
	}
	return false
}

// rowFormatChanged reports whether the table ROW_FORMAT differs in a way that matters.
// ROW_FORMAT in both the target SDL and the synced SHOW CREATE readback is a DECLARED
// option, not the table's physical row format: the oracle confirms SHOW CREATE echoes
// ROW_FORMAT=X iff the user explicitly declared a non-DEFAULT value (a bare table, and a
// ROW_FORMAT=DEFAULT table, both read back with NO clause; an explicit ROW_FORMAT=DYNAMIC
// — even though DYNAMIC is the 8.0 physical default — IS echoed). So both sides are the
// same kind of declared-option string and are compared directly.
//
// The ignore DECISION is owned by normalize-core's IgnoreRowFormat (empty / DEFAULT →
// ignorable, mapped to "" here; an explicit non-default value is significant). The
// comparison is then:
//   - both ignorable → no diff;
//   - otherwise compare the canonical values (ignorable normalized to "") — so an
//     explicit COMPRESSED vs an unspecified/DEFAULT side IS a real change (the user is
//     adding or removing a recorded create-option), while a bare side and a
//     ROW_FORMAT=DEFAULT side still compare equal (both "").
//
// Normalizing ignorable → "" (NOT to a physical default like DYNAMIC) keeps the result
// independent of innodb_default_row_format, the environment dependency flag
// row-format-default-source warns about: a bare table never acquires a phantom format.
// FLAG: if normalize-core grows a CanonicalRowFormat key, compare those instead — see PR.
func rowFormatChanged(from, to *Table, n *Normalizer) bool {
	if n.IgnoreRowFormat(from) && n.IgnoreRowFormat(to) {
		return false
	}
	return !strings.EqualFold(canonicalRowFormat(from, n), canonicalRowFormat(to, n))
}

// canonicalRowFormat returns the table's declared ROW_FORMAT, or "" when normalize-core
// deems it ignorable (empty / DEFAULT). The value match in rowFormatChanged is a plain
// case-insensitive compare because ROW_FORMAT names (DYNAMIC, COMPRESSED, COMPACT,
// REDUNDANT) have no aliasing or version rules.
func canonicalRowFormat(t *Table, n *Normalizer) string {
	if n.IgnoreRowFormat(t) {
		return ""
	}
	return strings.TrimSpace(t.RowFormat)
}

// tableSubdiffsChanged reports whether any sub-object slice on a TableDiffEntry signals
// a change: the column diff (populated by diff-core) or the breadth-owned
// Indexes/Constraints/ForeignKeys/Checks slices or the PartitionChanged flag (populated
// by later nodes). compareTable calls it so the "is this table modified?" decision lives
// in one place — a breadth node that fills its slice needs no edit to compareTable.
func tableSubdiffsChanged(e *TableDiffEntry) bool {
	return len(e.Columns) > 0 ||
		len(e.Indexes) > 0 ||
		len(e.Constraints) > 0 ||
		len(e.ForeignKeys) > 0 ||
		len(e.Checks) > 0 ||
		e.PartitionChanged
}
