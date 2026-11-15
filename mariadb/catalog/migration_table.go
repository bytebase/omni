package catalog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// generateTableDDL produces CREATE TABLE, DROP TABLE, and table-option ALTER TABLE
// operations from the diff. Column changes on a modified table belong to generateColumnDDL.
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

// tableOpRank orders table-level op types: drops, then creates, then option-alters.
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

// tableOptionAlterOp sets each changed option and resets one the target drops to its default,
// since ALTER TABLE keeps an option it does not name.
func tableOptionAlterOp(database, name string, from, to *Table, n *Normalizer) (MigrationOp, bool) {
	var clauses []string
	if n.CanonicalEngine(from) != n.CanonicalEngine(to) {
		clauses = append(clauses, "ENGINE="+tableEngineName(to))
	}
	if toLower(from.Charset) != toLower(to.Charset) || toLower(from.Collation) != toLower(to.Collation) {
		if to.Charset != "" {
			clauses = append(clauses, "DEFAULT CHARSET="+to.Charset)
		}
		if to.Collation != "" {
			clauses = append(clauses, "COLLATE="+to.Collation)
		}
	}
	for _, opt := range []struct {
		name, from, to, reset string
	}{
		{"MIN_ROWS", from.MinRows, to.MinRows, "0"},
		{"MAX_ROWS", from.MaxRows, to.MaxRows, "0"},
		{"AVG_ROW_LENGTH", from.AvgRowLength, to.AvgRowLength, "0"},
		{"PACK_KEYS", from.PackKeys, to.PackKeys, "DEFAULT"},
		{"STATS_PERSISTENT", from.StatsPersistent, to.StatsPersistent, "DEFAULT"},
		{"STATS_AUTO_RECALC", from.StatsAutoRecalc, to.StatsAutoRecalc, "DEFAULT"},
		{"STATS_SAMPLE_PAGES", from.StatsSamplePages, to.StatsSamplePages, "DEFAULT"},
		{"CHECKSUM", from.Checksum, to.Checksum, "0"},
		{"PAGE_CHECKSUM", from.PageChecksum, to.PageChecksum, "DEFAULT"},
		{"DELAY_KEY_WRITE", from.DelayKeyWrite, to.DelayKeyWrite, "0"},
		{"ROW_FORMAT", strings.ToUpper(from.RowFormat), strings.ToUpper(to.RowFormat), "DEFAULT"},
		{"TRANSACTIONAL", from.Transactional, to.Transactional, "DEFAULT"},
		{"KEY_BLOCK_SIZE", strconv.Itoa(from.KeyBlockSize), strconv.Itoa(to.KeyBlockSize), "0"},
	} {
		if strings.EqualFold(opt.from, opt.to) {
			continue
		}
		value := opt.to
		if value == "" || value == "0" && opt.reset == "0" {
			value = opt.reset
		}
		clauses = append(clauses, opt.name+"="+value)
	}
	if from.Comment != to.Comment {
		clauses = append(clauses, "COMMENT="+quoteStringLiteral(to.Comment))
	}
	clauses = append(clauses, engineOptionClauses(from.EngineOptions, to.EngineOptions)...)
	if len(clauses) == 0 {
		return MigrationOp{}, false
	}
	return MigrationOp{
		Type:       OpAlterTable,
		Database:   database,
		ObjectName: name,
		SQL: fmt.Sprintf("ALTER TABLE %s %s",
			tableIdent(to), strings.Join(clauses, " ")),
		Phase:    PhaseMain,
		Priority: priorityTable,
		sortName: tableSortName(database, name),
	}, true
}

// engineOptionClauses sets each engine-defined option of the target that differs and resets
// one the target drops to DEFAULT.
func engineOptionClauses(from, to []EngineOption) []string {
	render := func(o EngineOption) string {
		if o.Quoted {
			return quoteStringLiteral(o.Value)
		}
		return o.Value
	}
	have := make(map[string]string)
	for _, o := range from {
		have[toLower(o.Name)] = render(o)
	}
	var out []string
	want := make(map[string]bool)
	for _, o := range to {
		want[toLower(o.Name)] = true
		if have[toLower(o.Name)] != render(o) {
			out = append(out, "`"+o.Name+"`="+render(o))
		}
	}
	for _, o := range from {
		if !want[toLower(o.Name)] {
			out = append(out, "`"+o.Name+"`=DEFAULT")
		}
	}
	return out
}

// formatCreateTable renders a CREATE TABLE for a target table with every column, the inline
// PRIMARY KEY, the AUTO_INCREMENT column's supporting key when one is needed, PERIOD FOR
// clauses, and the table options.
//
// Other secondary indexes, foreign keys, CHECK constraints, and partitioning are emitted by
// their own generators as follow-up ALTER TABLE statements. The PK is inline because a
// column's NOT NULL rendering depends on it, and the AUTO_INCREMENT supporting key is inline
// because MariaDB requires that column to be the first column of a key when the table is
// created (see autoIncSupportingIndex).
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
	// A follow-up ADD INDEX runs after the CREATE, which MariaDB rejects for an AUTO_INCREMENT
	// column without a key; the index generator skips this same index.
	if idx := autoIncSupportingIndex(tbl); idx != nil {
		lines = append(lines, "  "+formatIndexDefinition(idx, n))
	}
	// The system-versioning period is declared with its ROW START and ROW END columns, and an
	// application-time period ahead of the keys that are WITHOUT OVERLAPS over it.
	if tbl.AppPeriodName != "" {
		lines = append(lines, "  PERIOD FOR "+backtickIDQ(tbl.AppPeriodName)+" ("+backtickIDQ(tbl.AppPeriodStartCol)+", "+backtickIDQ(tbl.AppPeriodEndCol)+")")
	}
	if tbl.PeriodStartCol != "" {
		lines = append(lines, "  PERIOD FOR SYSTEM_TIME ("+backtickIDQ(tbl.PeriodStartCol)+", "+backtickIDQ(tbl.PeriodEndCol)+")")
	}
	b.WriteString(strings.Join(lines, ",\n"))
	b.WriteString("\n)")

	if opts := formatCreateTableOptions(tbl, n); opts != "" {
		b.WriteString(" ")
		b.WriteString(opts)
	}

	return b.String()
}

// autoIncSupportingIndex returns the index that formatCreateTable must inline so the table's
// AUTO_INCREMENT column is the first column of a key at CREATE time, or nil when no inline is
// needed: no AUTO_INCREMENT column, the inline PRIMARY KEY already starts with it, or no key
// has it first. The index generator (addIndexOpsForNewTable) skips the same index so a
// follow-up ALTER does not add a duplicate.
//
// A user index is preferred (UNIQUE over plain, then lower-cased name). The index created
// implicitly for a FOREIGN KEY is the last resort, because the FK is added in a later phase
// and the CREATE would otherwise leave the column unkeyed.
func autoIncSupportingIndex(tbl *Table) *Index {
	autoCol := autoIncrementColumnName(tbl)
	if autoCol == "" {
		return nil
	}
	if pk := primaryKeyIndex(tbl); pk != nil && indexFirstColumnIs(pk, autoCol) {
		return nil
	}
	skip := fkImplicitIndexNames(tbl)
	var best, fkBackup *Index
	for _, idx := range tbl.Indexes {
		if idx == nil || idx.Primary {
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

// autoIncrementColumnName returns the name of the table's AUTO_INCREMENT column, or "" when
// there is none.
func autoIncrementColumnName(tbl *Table) string {
	for _, col := range tbl.Columns {
		if col != nil && col.AutoIncrement {
			return col.Name
		}
	}
	return ""
}

// primaryKeyIndex returns the table's PRIMARY KEY index, or nil.
func primaryKeyIndex(tbl *Table) *Index {
	for _, idx := range tbl.Indexes {
		if idx != nil && idx.Primary {
			return idx
		}
	}
	return nil
}

// indexFirstColumnIs reports whether the index's first key part is the named column (a plain
// column reference, not an expression), case-insensitively.
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
// when the table has no PK. The index form is preferred and the constraint form is the
// fallback.
func formatInlinePrimaryKey(tbl *Table) string {
	for _, idx := range tbl.Indexes {
		if !idx.Primary {
			continue
		}
		// As SHOW CREATE writes it, so a part WITHOUT OVERLAPS keeps it.
		return showIndex(idx)
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

// migrationIndexColumn renders one index/PK key part: `col`[(len)][ DESC] or (expr).
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

// formatCreateTableOptions is the option list SHOW CREATE TABLE ends the table with, less
// AUTO_INCREMENT=N, and WITH SYSTEM VERSIONING where the table is versioned.
func formatCreateTableOptions(tbl *Table, n *Normalizer) string {
	opts := n.CanonicalTableOptions(tbl)
	if tbl.SystemVersioned {
		opts += " WITH SYSTEM VERSIONING"
	}
	return opts
}

// tableEngineName returns the declared engine, or InnoDB when unspecified.
func tableEngineName(tbl *Table) string {
	e := strings.TrimSpace(tbl.Engine)
	if e == "" {
		return "InnoDB"
	}
	return e
}

// tableIdent returns the backtick-quoted database.table identifier using the table's
// original-case names, not the lower-cased diff keys, because object names are case-significant
// when lower_case_table_names=0. The database qualifier is omitted when the table has no
// Database.
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

// migrationQuoteIdent backtick-quotes an identifier, doubling any embedded backtick.
func migrationQuoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
