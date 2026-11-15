package catalog

import (
	"fmt"
	"slices"
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
	return tableOptionAlterOpExceptEngine(database, name, from, to, n, false)
}

// tableOptionAlterOpExceptEngine is tableOptionAlterOp, omitting the ENGINE
// clause when skipEngine is set: foldEngineIntoAutoKey moves the engine change
// into the key statement and regenerates the rest.
func tableOptionAlterOpExceptEngine(database, name string, from, to *Table, n *Normalizer, skipEngine bool) (MigrationOp, bool) {
	var clauses []string
	if !skipEngine && n.CanonicalEngine(from) != n.CanonicalEngine(to) {
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
		Priority: tableOptionPriority(from, to),
		sortName: tableSortName(database, name),
	}, true
}

// tableOptionPriority runs a table's option ALTER ahead of its column changes,
// unless it changes the engine while a column gives up AUTO_INCREMENT: the
// target engine may refuse the column as the source holds it, as InnoDB
// refuses an AUTO_INCREMENT column that does not lead a key (1075). An engine
// change that instead needs the key the plan builds is folded into that
// statement (foldEngineIntoAutoKey).
func tableOptionPriority(from, to *Table) int {
	if !eqFoldStr(from.Engine, to.Engine) {
		for _, col := range from.Columns {
			if target := to.GetColumn(col.Name); col.AutoIncrement && (target == nil || !target.AutoIncrement) {
				return priorityColumn + 4
			}
		}
	}
	return priorityTable
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
	// column without a key; the index generator skips this same index. When that key is IGNORED,
	// inline the implicit primary key too, else the IGNORED key is promoted to it (4174).
	if pk := inlineImplicitPrimaryKey(tbl); pk != nil {
		lines = append(lines, "  "+formatIndexDefinition(pk, n))
	}
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
	if pk := primaryKeyIndex(tbl); pk != nil && indexBacksAutoColumn(pk, autoCol, tbl.Engine) {
		return nil
	}
	skip := fkImplicitIndexNames(tbl)
	var best, fkBackup *Index
	for _, idx := range tbl.Indexes {
		if idx == nil || idx.Primary {
			continue
		}
		// MyISAM and Aria take a key the column is any plain part of.
		if !indexBacksAutoColumn(idx, autoCol, tbl.Engine) {
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

// autoIncIndexPreferred reports whether candidate should be chosen over current, a key the
// target places before it, as the inlined AUTO_INCREMENT supporting key: a UNIQUE key wins over
// a non-unique one, and otherwise the target's order holds, since the first UNIQUE key over NOT
// NULL columns is the implicit primary key, which an IGNORED key cannot be (4174).
func autoIncIndexPreferred(candidate, current *Index) bool {
	return candidate.Unique && !current.Unique
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

// restatePageChecksum sets PAGE_CHECKSUM again after the last ALTER of a table
// that keeps it: MariaDB drops it in an ALTER that rebuilds an InnoDB or MyISAM
// table without naming it (see alterTableKeepingOptions).
func restatePageChecksum(ops []MigrationOp, diff *SchemaDiff) []MigrationOp {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action == DiffDrop || e.To == nil || e.To.PageChecksum == "" || eqFoldStr(e.To.Engine, "Aria") {
			continue
		}
		table := tableIdent(e.To)
		last := -1
		for j, op := range ops {
			if op.Database == e.Database && toLower(op.ObjectName) == toLower(e.Name) && strings.HasPrefix(op.SQL, "ALTER TABLE "+table+" ") {
				last = j
			}
		}
		if last < 0 || strings.Contains(ops[last].SQL, "PAGE_CHECKSUM=") {
			continue
		}
		restate := ops[last]
		restate.Type, restate.SQL, restate.Warning = OpAlterTable, fmt.Sprintf("ALTER TABLE %s PAGE_CHECKSUM=%s", table, e.To.PageChecksum), ""
		restate.alterClause, restate.alterParts, restate.addsIndex, restate.dropsIndex, restate.column = "", nil, nil, "", ""
		ops = slices.Insert(ops, last+1, restate)
	}
	return ops
}

// foldEngineIntoAutoKey folds an ENGINE change into the single ALTER that
// defines the AUTO_INCREMENT column's key, when the engines differ on whether
// that column may be a non-leading key part (MyISAM allows it, InnoDB does
// not). The whole ALTER is then checked once against the target engine with the
// final key, which neither an engine-first nor an engine-last order achieves in
// both directions. Measured on mariadb:11.8.8 and mariadb (MyISAM): 1075.
func foldEngineIntoAutoKey(ops []MigrationOp, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify || e.From == nil || e.To == nil || eqFoldStr(e.From.Engine, e.To.Engine) {
			continue
		}
		autoCol := autoIncrementColumnName(e.To)
		if autoCol == "" {
			continue
		}
		// Only when the AI column's non-leading status is allowed by one
		// engine and not the other, and it is non-leading somewhere.
		if engineAllowsNonFirstAutoKey(e.From.Engine) == engineAllowsNonFirstAutoKey(e.To.Engine) {
			continue
		}
		if leadsEveryKey(e.From, autoCol) && leadsEveryKey(e.To, autoCol) {
			continue
		}
		key := -1
		for j := range ops {
			if ops[j].Type == OpAddIndex && opOfTable(ops[j], e) && ops[j].addsIndex != nil && indexHasPlainColumn(ops[j].addsIndex, autoCol) {
				key = j
			}
		}
		if key < 0 {
			continue
		}
		engine := "ENGINE=" + tableEngineName(e.To)
		ops[key].alterParts = append(partsOf(ops[key]), alterPart{clause: engine})
		ops[key].SQL = "ALTER TABLE " + tableIdent(e.To) + " " + joinAlterParts(ops[key].alterParts)
		// Regenerate the standalone option op without the ENGINE clause, which
		// now rides the key statement; drop it when nothing else remains.
		replacement, remains := tableOptionAlterOpExceptEngine(e.Database, e.Name, e.From, e.To, n, true)
		var kept []MigrationOp
		for _, op := range ops {
			// Only the standalone table-option op, not a PARTITION BY or REMOVE
			// PARTITIONING ALTER of the same table.
			if op.Type == OpAlterTable && op.Database == e.Database && toLower(op.ObjectName) == toLower(e.Name) && op.alterParts == nil && !strings.Contains(op.SQL, "PARTITION") {
				if !remains {
					continue
				}
				op = replacement
			}
			kept = append(kept, op)
		}
		ops = kept
	}
	return ops
}

// leadsEveryKey reports whether col leads every key it is a part of.
func leadsEveryKey(tbl *Table, col string) bool {
	for _, idx := range tbl.Indexes {
		if indexHasPlainColumn(idx, col) && !indexFirstColumnIs(idx, col) {
			return false
		}
	}
	return true
}

// inlineImplicitPrimaryKey returns the unique key formatCreateTable must inline
// as the implicit primary key, when the AUTO_INCREMENT supporting key it inlines
// is IGNORED and a non-ignored unique key over NOT NULL columns precedes it: on
// its own the IGNORED key would be promoted to the implicit primary key (4174).
func inlineImplicitPrimaryKey(tbl *Table) *Index {
	ai := autoIncSupportingIndex(tbl)
	if ai == nil || ai.Visible || formatInlinePrimaryKey(tbl) != "" {
		return nil
	}
	pk := tablePrimaryKey(tbl, true)
	if pk == nil || pk == ai || pk.Primary || !pk.Visible {
		return nil
	}
	return pk
}
