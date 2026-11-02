package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateTableDDL produces CREATE TABLE, DROP TABLE, and table-option ALTER TABLE
// operations from the diff. It is the MySQL analog of PG's generateTableDDL
// (pg/catalog/migration_table.go:21).
//
//   - DiffAdd  → CREATE TABLE (full canonical form: columns + table options).
//   - DiffDrop → DROP TABLE.
//   - DiffModify → table-option ALTER TABLE (ENGINE/CHARSET/COLLATE/COMMENT/ROW_FORMAT)
//     when the table-level options changed. Column changes on a modified table are emitted
//     by generateColumnDDL (one ALTER per column); this generator only handles the
//     table-level option delta so the two never double-emit.
//
// Determinism: drops sorted by name, then creates sorted by name, then option-alters sorted
// by name — and sortMigrationOps re-imposes the global phase/name order on top.
func generateTableDDL(from, to *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:       OpCreateTable,
				Database:   entry.Database,
				ObjectName: entry.Name,
				SQL:        formatCreateTable(entry.To, n),
				Phase:      PhaseMain,
				Priority:   priorityTable,
				sortName:   tableSortName(entry.Database, entry.Name),
			})

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:       OpDropTable,
				Database:   entry.Database,
				ObjectName: entry.Name,
				SQL:        fmt.Sprintf("DROP TABLE %s", tableIdent(entry.From)),
				Phase:      PhasePre,
				Priority:   priorityTable,
				sortName:   tableSortName(entry.Database, entry.Name),
			})

		case DiffModify:
			// Table-level option changes only; column changes are generateColumnDDL's job.
			if entry.From == nil || entry.To == nil {
				continue
			}
			if op, ok := tableOptionAlterOp(entry.Database, entry.Name, entry.From, entry.To, n); ok {
				ops = append(ops, op)
			}
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			// Drops first, then creates, then option-alters — by op-type rank.
			return tableOpRank(ops[i].Type) < tableOpRank(ops[j].Type)
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// tableOpRank orders table-level op types so a deterministic local sort places drops, then
// creates, then option-alters (the global sortMigrationOps re-sorts by phase too).
func tableOpRank(t MigrationOpType) int {
	switch t {
	case OpDropTable:
		return 0
	case OpCreateTable:
		return 1
	default: // OpAlterTable
		return 2
	}
}

// tableOptionAlterOp builds an ALTER TABLE that reconciles the table-level options that
// diff-core's tableOptionsChanged compares (engine, effective charset/collation, comment,
// significant ROW_FORMAT). It emits only the clauses that actually changed, in a stable
// order, so a column-only modification produces NO table-option ALTER (the second return is
// false). Each option is rendered through the same canonical form diff-core compared, so the
// readback of the ALTER canonicalizes equal to `to`.
func tableOptionAlterOp(database, name string, from, to *Table, n *Normalizer) (MigrationOp, bool) {
	var clauses []string

	// ENGINE — render the resolved (canonical) engine when it changed. CanonicalEngine
	// lower-cases and defaults to innodb; we emit the declared name (or the default) so the
	// readback echoes ENGINE=<X>.
	if n.CanonicalEngine(from) != n.CanonicalEngine(to) {
		clauses = append(clauses, "ENGINE="+tableEngineName(to))
	}

	// CHARSET / COLLATE — when the effective (charset, collation) pair changed, emit both.
	// MySQL's ALTER TABLE ... DEFAULT CHARSET=cs COLLATE=coll sets the table default; it does
	// NOT rewrite existing columns' stored charset, but bare columns inherit it and the diff
	// compares the effective pair, so emitting both keeps the readback canonical.
	charsetChanged := foldCharset(from.Charset) != foldCharset(to.Charset)
	collationChanged := n.effectiveTableCollation(from) != n.effectiveTableCollation(to)
	if charsetChanged || collationChanged {
		cs := foldCharset(to.Charset)
		coll := n.effectiveTableCollation(to)
		if cs != "" {
			clauses = append(clauses, "DEFAULT CHARSET="+cs)
		}
		// Emit COLLATE only for a non-default collation. When `coll` is the charset's
		// version default it is redundant — and emitting it risks naming a collation the
		// target server lacks (utf8mb4_0900_ai_ci does not exist on 5.7). Omitting it lets
		// the server apply its own default, which is valid on either version and reads back
		// canonically equal (entry utf8mb4-default-collation).
		if coll != "" && !n.isCharsetDefaultCollation(cs, coll) {
			clauses = append(clauses, "COLLATE="+coll)
		}
	}

	// COMMENT.
	if n.CanonicalComment(from.Comment) != n.CanonicalComment(to.Comment) {
		clauses = append(clauses, "COMMENT="+quoteStringLiteral(to.Comment))
	}

	// ROW_FORMAT — only when normalize-core deems the change significant (rowFormatChanged).
	if rowFormatChanged(from, to, n) {
		clauses = append(clauses, "ROW_FORMAT="+strings.ToUpper(canonicalRowFormatValue(to, n)))
	}

	if len(clauses) == 0 {
		return MigrationOp{}, false
	}

	return MigrationOp{
		Type:       OpAlterTable,
		Database:   database,
		ObjectName: name,
		SQL: fmt.Sprintf("ALTER TABLE %s %s",
			tableIdent(to), strings.Join(clauses, ", ")),
		Phase:    PhaseMain,
		Priority: priorityTable,
		sortName: tableSortName(database, name),
	}, true
}

// canonicalRowFormatValue returns the ROW_FORMAT value to emit when it changed. An
// ignorable (empty/DEFAULT) target side yields "DEFAULT" (reset to engine default); an
// explicit value is emitted verbatim. This mirrors rowFormatChanged's significance rule.
func canonicalRowFormatValue(t *Table, n *Normalizer) string {
	if n.IgnoreRowFormat(t) {
		return "DEFAULT"
	}
	return strings.TrimSpace(t.RowFormat)
}

// formatCreateTable renders a full canonical CREATE TABLE for a target table, with every
// column and the table options in MySQL's stored form for the target version. The output,
// applied and read back, canonicalizes equal to `tbl` because each piece routes through the
// same normalize-core canonical form diff-core compares. It is the MySQL analog of PG's
// FormatCreateTable (pg/catalog/migration_table.go:96).
//
// Scope: columns + table options + inline PRIMARY KEY + (when needed) the AUTO_INCREMENT
// column's supporting key. Other secondary indexes, foreign keys, CHECK constraints, and
// partitioning are breadth concerns and are intentionally NOT rendered here — a table created
// by this node carries its columns, PK, and the one key an AUTO_INCREMENT column must have;
// the index/FK/check/partition nodes emit their own ALTER TABLE ADD ... follow-ups. (The PK is
// rendered inline because a column's NOT NULL canonicalization is PK-aware; the AUTO_INCREMENT
// supporting key is rendered inline because MySQL requires the AUTO_INCREMENT column to be the
// first column of a key AT CREATE TIME — see autoIncSupportingIndex.)
func formatCreateTable(tbl *Table, n *Normalizer) string {
	var b strings.Builder
	b.WriteString("CREATE TABLE ")
	b.WriteString(tableIdent(tbl))
	b.WriteString(" (\n")

	var lines []string
	for _, col := range diffableColumns(tbl) {
		lines = append(lines, "  "+formatColumnDefinition(tbl, col, n))
	}
	if pk := formatInlinePrimaryKey(tbl); pk != "" {
		lines = append(lines, "  "+pk)
	}
	// Inline the AUTO_INCREMENT column's supporting key when its only backing key is a non-PK
	// index (a UNIQUE or secondary key). MySQL rejects CREATE TABLE for an AUTO_INCREMENT column
	// that is not the first column of a key ("there can be only one auto column and it must be
	// defined as a key", errno 1075), and the index node's ADD runs only AFTER the CREATE — so the
	// CREATE itself would fail without the key inline. The index node skips this same index (it
	// consults autoIncSupportingIndex) so the key is not also re-added.
	if idx := autoIncSupportingIndex(tbl); idx != nil {
		lines = append(lines, "  "+formatIndexDefinition(idx, n))
	}
	b.WriteString(strings.Join(lines, ",\n"))
	b.WriteString("\n)")

	if opts := formatCreateTableOptions(tbl, n); opts != "" {
		b.WriteString(" ")
		b.WriteString(opts)
	}

	return b.String()
}

// autoIncSupportingIndex returns the index that formatCreateTable must inline so a table's
// AUTO_INCREMENT column is the first column of a key at CREATE time (MySQL errno 1075), or nil
// when no such inline is needed. It is the single shared predicate for the inline decision: the
// table generator renders the returned index inline, and the index node (addIndexOpsForNewTable)
// skips the SAME index so it is not also added by a follow-up ALTER (which would fail errno 1061,
// duplicate key name).
//
// nil is returned when:
//   - the table has no AUTO_INCREMENT column;
//   - the AUTO_INCREMENT column is already the first column of the inline PRIMARY KEY (the PK
//     satisfies errno 1075, and the PK is rendered by formatInlinePrimaryKey);
//   - no key at all has the AUTO_INCREMENT column first (such a table is invalid in MySQL anyway
//     and the engine would reject it regardless of how the CREATE is rendered).
//
// Selection prefers a USER index (UNIQUE over plain secondary; ties break on lower-cased name) so
// the inline key matches the most common authoring. An FK-IMPLICIT backing index is used only as a
// LAST RESORT — when the AUTO_INCREMENT column's only first-column key is the index MySQL
// auto-creates for a FOREIGN KEY on it. That case is real (a `FOREIGN KEY (id)` on an
// AUTO_INCREMENT `id` with no other key), and without inlining that backing index the CREATE TABLE
// would have an unkeyed AUTO_INCREMENT column (errno 1075): the index node skips FK-implicit
// indexes and the FK node adds the constraint only in a later phase. Inlining the backing index is
// safe — the index node skips it (both because it is FK-implicit and because it is the inlined
// index), and the deferred FK reuses the existing covering index rather than creating a duplicate.
// The generated invisible primary key is never chosen.
func autoIncSupportingIndex(tbl *Table) *Index {
	autoCol := autoIncrementColumnName(tbl)
	if autoCol == "" {
		return nil
	}
	// If the inline PK already starts with the AUTO_INCREMENT column, the PK is the key — no
	// extra inline. (A GIPK-only table has no AUTO_INCREMENT column, handled above.)
	if pk := primaryKeyIndex(tbl); pk != nil && indexFirstColumnIs(pk, autoCol) {
		return nil
	}
	skip := fkImplicitIndexNames(tbl)
	var best, fkBackup *Index
	for _, idx := range tbl.Indexes {
		if idx == nil || idx.Primary || isGeneratedInvisiblePrimaryKeyIndex(idx) {
			continue
		}
		if !indexFirstColumnIs(idx, autoCol) {
			continue
		}
		if skip[toLower(idx.Name)] {
			// FK-implicit backing index: remember it as a fallback only.
			if fkBackup == nil || toLower(idx.Name) < toLower(fkBackup.Name) {
				fkBackup = idx
			}
			continue
		}
		if best == nil || autoIncIndexPreferred(idx, best) {
			best = idx
		}
	}
	if best != nil {
		return best
	}
	return fkBackup
}

// autoIncrementColumnName returns the name of the table's AUTO_INCREMENT column (MySQL allows at
// most one), or "" when there is none. The generated invisible primary key column is never
// AUTO_INCREMENT, so it is naturally excluded.
func autoIncrementColumnName(tbl *Table) string {
	for _, col := range tbl.Columns {
		if col != nil && col.AutoIncrement {
			return col.Name
		}
	}
	return ""
}

// primaryKeyIndex returns the table's PRIMARY KEY index (excluding the generated invisible
// primary key), or nil.
func primaryKeyIndex(tbl *Table) *Index {
	for _, idx := range tbl.Indexes {
		if idx != nil && idx.Primary && !isGeneratedInvisiblePrimaryKeyIndex(idx) {
			return idx
		}
	}
	return nil
}

// indexFirstColumnIs reports whether the index's first key part is the named column (a plain
// column reference, not an expression), case-insensitively. MySQL's errno-1075 rule is about the
// FIRST column of a key, so only the leading part matters.
func indexFirstColumnIs(idx *Index, col string) bool {
	if idx == nil || len(idx.Columns) == 0 {
		return false
	}
	first := idx.Columns[0]
	return first.Expr == "" && strings.EqualFold(first.Name, col)
}

// autoIncIndexPreferred reports whether candidate should be chosen over current as the inlined
// AUTO_INCREMENT supporting key: a UNIQUE key wins over a non-unique one; otherwise the
// lower-cased index name breaks the tie deterministically.
func autoIncIndexPreferred(candidate, current *Index) bool {
	if candidate.Unique != current.Unique {
		return candidate.Unique
	}
	return toLower(candidate.Name) < toLower(current.Name)
}

// formatInlinePrimaryKey renders the inline PRIMARY KEY clause for a CREATE TABLE, or ""
// when the table has no PK. The PK can be expressed via an index (Primary) or a constraint
// (ConPrimaryKey); the loader uses the index form for `PRIMARY KEY (...)` so that is the
// primary source, with the constraint form as a fallback. Prefix lengths and DESC are
// rendered (a PK may carry them); the generated invisible primary key is never emitted (it
// is engine-synthesized).
func formatInlinePrimaryKey(tbl *Table) string {
	for _, idx := range tbl.Indexes {
		if !idx.Primary {
			continue
		}
		if isGeneratedInvisiblePrimaryKeyIndex(idx) {
			return ""
		}
		cols := make([]string, 0, len(idx.Columns))
		for _, ic := range idx.Columns {
			cols = append(cols, migrationIndexColumn(ic))
		}
		return "PRIMARY KEY (" + strings.Join(cols, ",") + ")"
	}
	for _, con := range tbl.Constraints {
		if con.Type != ConPrimaryKey {
			continue
		}
		cols := make([]string, 0, len(con.Columns))
		for _, c := range con.Columns {
			cols = append(cols, migrationQuoteIdent(c))
		}
		return "PRIMARY KEY (" + strings.Join(cols, ",") + ")"
	}
	return ""
}

// migrationIndexColumn renders one index/PK key part: `col`[(len)][ DESC] or (expr). DESC is
// rendered unconditionally (it is harmless on 5.7, which stores ascending — the diff already
// drops 5.7 DESC via CanonicalIndexColumn, so a round-tripped PK key still compares equal).
func migrationIndexColumn(ic *IndexColumn) string {
	var b strings.Builder
	if ic.Expr != "" {
		b.WriteString("(")
		b.WriteString(ic.Expr)
		b.WriteString(")")
	} else {
		b.WriteString(migrationQuoteIdent(ic.Name))
		if ic.Length > 0 {
			fmt.Fprintf(&b, "(%d)", ic.Length)
		}
	}
	if ic.Descending {
		b.WriteString(" DESC")
	}
	return b.String()
}

// formatCreateTableOptions renders the trailing table options (ENGINE/CHARSET/COLLATE/
// COMMENT/ROW_FORMAT) in canonical form for a CREATE TABLE. AUTO_INCREMENT=N is never
// emitted (ignore-in-diff: it is a live counter, not schema). ENGINE is always emitted
// (MySQL always echoes it); charset/collation are emitted when the table declares a charset
// so the readback's effective pair matches.
func formatCreateTableOptions(tbl *Table, n *Normalizer) string {
	var parts []string

	parts = append(parts, "ENGINE="+tableEngineName(tbl))

	if cs := foldCharset(tbl.Charset); cs != "" {
		parts = append(parts, "DEFAULT CHARSET="+cs)
		// Emit COLLATE only for a non-default collation (see tableOptionAlterOp): a bare
		// CHARSET table renders no COLLATE, so the server applies its version default and
		// the DDL is valid on both 5.7 and 8.0 (never naming a missing utf8mb4_0900_ai_ci).
		if coll := n.effectiveTableCollation(tbl); coll != "" && !n.isCharsetDefaultCollation(cs, coll) {
			parts = append(parts, "COLLATE="+coll)
		}
	}

	if !n.IgnoreRowFormat(tbl) {
		parts = append(parts, "ROW_FORMAT="+strings.ToUpper(strings.TrimSpace(tbl.RowFormat)))
	}

	if c := n.CanonicalComment(tbl.Comment); c != "" {
		parts = append(parts, "COMMENT="+quoteStringLiteral(c))
	}

	return strings.Join(parts, " ")
}

// tableEngineName returns the engine name to emit for a table — the declared engine, or the
// server default (InnoDB) when unspecified, so the readback's ENGINE clause matches the
// canonical engine diff-core compares.
func tableEngineName(tbl *Table) string {
	e := strings.TrimSpace(tbl.Engine)
	if e == "" {
		return "InnoDB"
	}
	return e
}

// tableIdent returns the backtick-quoted database.table identifier for migration DDL,
// using the table's ORIGINAL-CASE names (Table.Name and Table.Database.Name), not the
// lower-cased diff identity keys (TableDiffEntry.Database/Name are folded to lower case for
// matching). On a case-sensitive server (lower_case_table_names=0, the Linux default) the
// stored object name is case-significant, so the DDL must reproduce the declared casing or it
// targets the wrong/non-existent object. The database qualifier is omitted when the table has
// no Database (the synced single-database release path may load with no qualifying database).
func tableIdent(tbl *Table) string {
	if tbl.Database == nil || tbl.Database.Name == "" {
		return migrationQuoteIdent(tbl.Name)
	}
	return migrationQuoteIdent(tbl.Database.Name) + "." + migrationQuoteIdent(tbl.Name)
}

// tableSortName is the stable secondary sort key for a table-level op: lower-cased
// database.table.
func tableSortName(database, table string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table)
}

// migrationQuoteIdent backtick-quotes a MySQL identifier, doubling any embedded backtick.
// Migration DDL always quotes identifiers to avoid reserved-word collisions.
func migrationQuoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// quoteStringLiteral single-quotes a string for a DDL string-literal context (COMMENT,
// string default), escaping backslashes and single quotes the way MySQL's stored form does
// (reusing show.go's escapeComment, which doubles ' and escapes \).
func quoteStringLiteral(s string) string {
	return "'" + escapeComment(s) + "'"
}
