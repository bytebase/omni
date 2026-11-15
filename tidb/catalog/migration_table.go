package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// --- Foreign keys ---

// dropForeignKeys drops the foreign keys of dropped tables, those that go or
// change, and those over a column whose type changes, which TiDB checks
// against the referenced column's type (3780) even when both change.
func (g *generator) dropForeignKeys(diff *SchemaDiff, to *Catalog) {
	retyped := retypedColumns(diff)
	for i := range diff.Tables {
		e := &diff.Tables[i]
		switch e.Action {
		case DiffDrop:
			for _, con := range e.From.Constraints {
				if con.Type == ConForeignKey {
					g.add(OpDropForeignKey, e.From, alterTable(e.From, "DROP FOREIGN KEY "+quoteIdent(con.Name)), "")
				}
			}
		case DiffModify:
			for _, fk := range e.ForeignKeys {
				if fk.Action != DiffAdd {
					g.add(OpDropForeignKey, e.From, alterTable(e.From, "DROP FOREIGN KEY "+quoteIdent(fk.Name)), "")
				}
			}
		default:
		}
	}
	if len(retyped) == 0 {
		return
	}
	changed := changedForeignKeys(diff)
	for _, tbl := range sortedTables(to) {
		for _, con := range tbl.Constraints {
			if con.Type != ConForeignKey || changed[con] || !foreignKeyReadsRetyped(tbl, con, retyped) {
				continue
			}
			g.readdOrder = append(g.readdOrder, con)
			g.add(OpDropForeignKey, tbl, alterTable(tbl, "DROP FOREIGN KEY "+quoteIdent(con.Name)), "")
		}
	}
}

// retypedColumns names the columns of changed tables whose type, charset or
// collation changes, keyed "database.table.column", lower-cased.
func retypedColumns(diff *SchemaDiff) map[string]bool {
	out := make(map[string]bool)
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify {
			continue
		}
		for _, c := range e.Columns {
			if c.Action == DiffModify && columnRetyped(c.From, c.To) {
				out[columnKey(e.To.Database, e.To.Name, c.Name)] = true
			}
		}
	}
	return out
}

func columnRetyped(a, b *Column) bool {
	return !strings.EqualFold(a.ColumnType, b.ColumnType) || !strings.EqualFold(a.Charset, b.Charset) || !strings.EqualFold(a.Collation, b.Collation)
}

func columnKey(db *Database, table, column string) string {
	return toLower(databaseName(db)) + "." + toLower(table) + "." + toLower(column)
}

// changedForeignKeys holds the foreign keys of `to` the diff adds or
// changes, which step 7 adds anyway.
func changedForeignKeys(diff *SchemaDiff) map[*Constraint]bool {
	out := make(map[*Constraint]bool)
	for i := range diff.Tables {
		e := &diff.Tables[i]
		switch e.Action {
		case DiffAdd:
			for _, con := range e.To.Constraints {
				out[con] = true
			}
		case DiffModify:
			for _, fk := range e.ForeignKeys {
				if fk.To != nil {
					out[fk.To] = true
				}
			}
		default:
		}
	}
	return out
}

func foreignKeyReadsRetyped(tbl *Table, con *Constraint, retyped map[string]bool) bool {
	for _, col := range con.Columns {
		if retyped[columnKey(tbl.Database, tbl.Name, col)] {
			return true
		}
	}
	refDB := con.RefDatabase
	if refDB == "" {
		refDB = databaseName(tbl.Database)
	}
	for _, col := range con.RefColumns {
		if retyped[toLower(refDB)+"."+toLower(con.RefTable)+"."+toLower(col)] {
			return true
		}
	}
	return false
}

// addForeignKeys adds the foreign keys of new tables, those that are new or
// changed, and those dropped around a column's type change.
func (g *generator) addForeignKeys(diff *SchemaDiff) {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		switch e.Action {
		case DiffAdd:
			for _, con := range e.To.Constraints {
				if con.Type == ConForeignKey {
					g.add(OpAddForeignKey, e.To, alterTable(e.To, "ADD "+showForeignKey(e.To, con)), "")
				}
			}
		case DiffModify:
			for _, fk := range e.ForeignKeys {
				if fk.To != nil {
					g.add(OpAddForeignKey, e.To, alterTable(e.To, "ADD "+showForeignKey(e.To, fk.To)), "")
				}
			}
		default:
		}
	}
	for _, con := range g.readdOrder {
		g.add(OpAddForeignKey, con.Table, alterTable(con.Table, "ADD "+showForeignKey(con.Table, con)), "")
	}
}

// --- Tables ---

// createTable creates a table without its foreign keys, which step 7 adds,
// and without the AUTO_INCREMENT and AUTO_RANDOM_BASE counters, which are
// data rather than schema. TiFlash replicas, which CREATE TABLE cannot
// declare, follow it.
func (g *generator) createTable(tbl *Table) {
	cp := *tbl
	cp.AutoIncrement, cp.AutoRandomBase = 0, 0
	cp.Constraints = nil
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey {
			cp.Constraints = append(cp.Constraints, con)
		}
	}
	text := showCreateTable(&cp)
	// Qualify the name, which SHOW CREATE TABLE writes alone.
	text = strings.Replace(text, " "+quoteIdent(tbl.Name)+" (\n", " "+tableIdent(tbl)+" (\n", 1)
	for _, con := range cp.Constraints {
		if con.Type == ConCheck {
			text = strings.Replace(text, showCheck(con), checkDefinition(con), 1)
		}
	}
	g.add(OpCreateTable, tbl, text, "")
	if replica := tiflashReplica(tbl); replica != "" {
		g.add(OpAlterTable, tbl, alterTable(tbl, strings.TrimPrefix(replica, " ")), "")
	}
}

// dropTableParts drops what a changed table loses before anything is added:
// checks, keys (and those an expression key over a changing column needs
// rebuilt), then columns.
func (g *generator) dropTableParts(e *TableDiffEntry) {
	// A retained generated column whose expression stops reading a dropped
	// column is modified before the drop, so the drop no longer fails on the
	// old expression (3108/3837).
	preDrop := preDropModifies(e)
	for _, c := range e.Columns {
		if preDrop[toLower(c.Name)] {
			g.add(OpModifyColumn, e.To, alterTable(e.To, "MODIFY COLUMN "+canonicalColumn(c.To)), modifyColumnWarning(e.From, e.To, c.From, c.To))
		}
	}
	for _, c := range e.Checks {
		if c.Action != DiffAdd {
			g.add(OpDropCheck, e.From, alterTable(e.From, "DROP CHECK "+quoteIdent(c.Name)), "")
		}
	}
	for _, idx := range e.Indexes {
		if idx.Action != DiffAdd {
			g.dropIndex(e.From, idx.From)
		}
	}
	for _, idx := range rebuiltIndexes(e) {
		g.dropIndex(e.From, e.From.indexNamed(idx.Name))
	}
	dropped, _ := droppedColumns(e)
	for _, col := range dropped {
		g.add(OpDropColumn, e.From, alterTable(e.From, "DROP COLUMN "+quoteIdent(col.Name)), "")
	}
	if e.PartitionChanged && e.From.Partitioning != nil && (e.To.Partitioning == nil || addsUniqueKey(e) && !appendsPartitions(e.From.Partitioning, e.To.Partitioning)) {
		// A new unique key must cover the partitioning in force when it is
		// added (8264), so the old partitioning goes first.
		g.add(OpAlterTable, e.From, alterTable(e.From, "REMOVE PARTITIONING"), "")
	}
}

// droppedColumns lists the columns a changed table loses, in the order they
// can go: generated columns first, latest first, since a generated column may
// read the columns before it (3108 otherwise). When the table loses every
// column, the last one is held back as last: TiDB refuses to drop a table's
// last column (1090), so it goes once the new columns are in.
func droppedColumns(e *TableDiffEntry) (now []*Column, last *Column) {
	drop := make(map[string]bool)
	for _, c := range e.Columns {
		if c.Action == DiffDrop {
			drop[toLower(c.Name)] = true
		}
	}
	var generated, plain []*Column
	for _, col := range e.From.Columns {
		switch {
		case !drop[toLower(col.Name)]:
		case col.Generated != nil:
			generated = append([]*Column{col}, generated...)
		default:
			plain = append(plain, col)
		}
	}
	now = append(generated, plain...)
	if len(now) > 0 && len(now) == len(e.From.Columns) {
		return now[:len(now)-1], now[len(now)-1]
	}
	return now, nil
}

func addsUniqueKey(e *TableDiffEntry) bool {
	for _, idx := range e.Indexes {
		if idx.To != nil && (idx.To.Unique || idx.To.Primary) {
			return true
		}
	}
	return false
}

func (g *generator) dropIndex(tbl *Table, idx *Index) {
	if !idx.Primary {
		g.add(OpDropIndex, tbl, alterTable(tbl, "DROP INDEX "+quoteIdent(idx.Name)), "")
		return
	}
	warning := ""
	if primaryKeyClustered(tbl, idx) {
		warning = "TiDB cannot drop a clustered primary key (8200)"
	}
	g.add(OpDropIndex, tbl, alterTable(tbl, "DROP PRIMARY KEY"), warning)
}

// rebuiltIndexes are the keys of `to` the diff leaves as they are but that
// read a column the diff changes through an expression: TiDB refuses to
// change such a column (3837), so the key is dropped and added back.
func rebuiltIndexes(e *TableDiffEntry) []*Index {
	changed := make(map[string]bool)
	for _, idx := range e.Indexes {
		changed[toLower(idx.Name)] = true
	}
	var out []*Index
	for _, idx := range e.To.Indexes {
		if changed[toLower(idx.Name)] || idx.Primary || e.From.indexNamed(idx.Name) == nil {
			continue
		}
		if indexExpressionReads(idx, e.Columns) {
			out = append(out, idx)
		}
	}
	return out
}

func indexExpressionReads(idx *Index, columns []ColumnDiffEntry) bool {
	for _, ic := range idx.Columns {
		if ic.Expr == "" {
			continue
		}
		expr := toLower(ic.Expr)
		for _, c := range columns {
			if c.Action == DiffModify && strings.Contains(expr, toLower(quoteIdent(c.Name))) {
				return true
			}
		}
	}
	return false
}

func (t *Table) indexNamed(name string) *Index {
	for _, idx := range t.Indexes {
		if strings.EqualFold(idx.Name, name) || (idx.Primary && strings.EqualFold(name, "PRIMARY")) {
			return idx
		}
	}
	return nil
}

// alterTable adds what a changed table gains or changes: options, columns,
// keys, checks, then partitioning, which TiDB checks against the keys.
func (g *generator) alterTable(e *TableDiffEntry) {
	from, to := e.From, e.To
	// Changing a table between ordinary and global temporary has no ALTER:
	// the table is recreated, which empties it, so the op carries a warning.
	if from.GlobalTemporary != to.GlobalTemporary || from.OnCommit != to.OnCommit {
		warn := "TiDB cannot change a table to or from GLOBAL TEMPORARY in place; the table is recreated, losing its rows"
		g.add(OpDropTable, from, "DROP TABLE "+tableIdent(from), warn)
		g.createTableWithWarning(to, warn)
		return
	}
	if e.OptionsChanged {
		g.alterTableOptions(from, to)
	}
	g.addColumns(e)
	if _, last := droppedColumns(e); last != nil {
		g.add(OpDropColumn, from, alterTable(from, "DROP COLUMN "+quoteIdent(last.Name)), "")
	}
	preDrop := preDropModifies(e)
	for _, c := range e.Columns {
		if c.Action == DiffModify && !preDrop[toLower(c.Name)] {
			g.add(OpModifyColumn, to, alterTable(to, "MODIFY COLUMN "+canonicalColumn(c.To)), modifyColumnWarning(from, to, c.From, c.To))
		}
	}
	// TTL after the columns it references exist and have their target type.
	g.alterTTL(from, to, func(clause, warning string) { g.add(OpAlterTable, to, alterTable(to, clause), warning) })
	for _, idx := range e.Indexes {
		if idx.Action != DiffDrop {
			g.addIndex(e, idx.To)
		}
	}
	for _, idx := range rebuiltIndexes(e) {
		g.addIndex(e, idx)
	}
	for _, c := range e.Checks {
		if c.Action != DiffDrop {
			g.add(OpAddCheck, to, alterTable(to, "ADD "+checkDefinition(c.To)), "")
		}
	}
	switch {
	case !e.PartitionChanged || to.Partitioning == nil && from.Partitioning != nil:
		// Unchanged, or removed in dropTableParts.
	case from.Partitioning != nil && addsUniqueKey(e) && !appendsPartitions(from.Partitioning, to.Partitioning):
		g.add(OpAlterTable, to, alterTable(to, partitionChange(nil, to.Partitioning)), "")
	default:
		g.add(OpAlterTable, to, alterTable(to, partitionChange(from.Partitioning, to.Partitioning)), "")
	}
}

// checkDefinition is a check as DDL writes it. SHOW CREATE TABLE puts the
// expression TiDB stores in two pairs of parentheses, and TiDB stores what
// is inside the one pair CHECK takes, so writing SHOW CREATE TABLE's text
// back would store another pair.
func checkDefinition(con *Constraint) string {
	s := fmt.Sprintf("CONSTRAINT %s CHECK (%s)", quoteIdent(con.Name), stripOuterParens(con.CheckExpr))
	if con.NotEnforced {
		s += " NOT ENFORCED"
	}
	return s
}

// createTableWithWarning creates a table like createTable, but every op it
// emits carries warning.
func (g *generator) createTableWithWarning(tbl *Table, warning string) {
	start := len(g.plan.Ops)
	g.createTable(tbl)
	for i := start; i < len(g.plan.Ops); i++ {
		g.plan.Ops[i].Warning = warning
	}
}

// preDropModifies names the modified columns whose MODIFY must run before the
// table's columns are dropped: a retained generated column whose old
// expression reads a column the diff drops. Its new expression must not read a
// column the diff adds, which would not exist yet; such a column keeps the
// normal order.
func preDropModifies(e *TableDiffEntry) map[string]bool {
	dropped := make(map[string]bool)
	added := make(map[string]bool)
	for _, c := range e.Columns {
		switch c.Action {
		case DiffDrop:
			dropped[toLower(c.Name)] = true
		case DiffAdd:
			added[toLower(c.Name)] = true
		}
	}
	out := make(map[string]bool)
	for _, c := range e.Columns {
		if c.Action != DiffModify || c.From.Generated == nil {
			continue
		}
		if exprReadsAny(c.From.Generated.Expr, dropped) && !exprReadsAny(colGenExpr(c.To), added) {
			out[toLower(c.Name)] = true
		}
	}
	return out
}

func colGenExpr(c *Column) string {
	if c.Generated == nil {
		return ""
	}
	return c.Generated.Expr
}

// exprReadsAny reports whether a generation expression reads any of the
// lower-cased column names, matched as a backquoted identifier.
func exprReadsAny(expr string, names map[string]bool) bool {
	low := toLower(expr)
	for name := range names {
		if strings.Contains(low, "`"+name+"`") {
			return true
		}
	}
	return false
}

func (g *generator) addIndex(e *TableDiffEntry, idx *Index) {
	warning := ""
	if idx.Primary && primaryKeyClustered(e.To, idx) {
		warning = "TiDB cannot add a clustered primary key to an existing table (8200)"
	}
	g.add(OpAddIndex, e.To, alterTable(e.To, "ADD "+showIndex(e.To, idx)), warning)
}

// addColumns adds the new columns in the order of `to`, each after the
// column before it there, which exists by then.
func (g *generator) addColumns(e *TableDiffEntry) {
	added := make(map[string]bool)
	for _, c := range e.Columns {
		if c.Action == DiffAdd {
			added[toLower(c.Name)] = true
		}
	}
	for i, col := range e.To.Columns {
		if !added[toLower(col.Name)] {
			continue
		}
		position := " FIRST"
		if i > 0 {
			position = " AFTER " + quoteIdent(e.To.Columns[i-1].Name)
		}
		warning := ""
		switch {
		case col.AutoIncrement:
			warning = "TiDB cannot add an AUTO_INCREMENT column to an existing table (8200)"
		case col.AutoRandom:
			warning = "TiDB cannot add an AUTO_RANDOM column to an existing table (8216)"
		case col.Generated != nil && col.Generated.Stored:
			warning = "TiDB cannot add a stored generated column to an existing table (3106)"
		}
		g.add(OpAddColumn, e.To, alterTable(e.To, "ADD COLUMN "+canonicalColumn(col)+position), warning)
	}
}

// modifyColumnWarning says why TiDB refuses a column change, or "".
func modifyColumnWarning(fromTbl, toTbl *Table, a, b *Column) string {
	switch {
	case a.AutoRandom != b.AutoRandom || b.AutoRandomShardBits < a.AutoRandomShardBits:
		return "TiDB cannot add or remove AUTO_RANDOM, or lower its shard bits (8216)"
	case !a.AutoIncrement && b.AutoIncrement:
		return "TiDB cannot make an existing column AUTO_INCREMENT (8200)"
	case a.AutoIncrement && !b.AutoIncrement:
		return "TiDB removes AUTO_INCREMENT only with tidb_allow_remove_auto_inc on (8200)"
	case (a.Generated != nil && a.Generated.Stored) || (b.Generated != nil && b.Generated.Stored):
		return "TiDB cannot change a stored generated column (3106)"
	case a.Charset != "" && b.Charset != "" && !strings.EqualFold(a.Charset, b.Charset) &&
		!(strings.EqualFold(a.Charset, "utf8") && strings.EqualFold(b.Charset, "utf8mb4")):
		return "TiDB changes a column's charset only from utf8 to utf8mb4 (8200)"
	case columnRetyped(a, b) && !widens(a, b) && inClusteredPrimaryKey(fromTbl, a.Name):
		return "TiDB changes a clustered primary key column's type only where no row is rewritten (8200)"
	case columnRetyped(a, b) && readByGeneratedColumn(fromTbl, a.Name):
		return "TiDB cannot change the type of a column a generated column reads (3106)"
	}
	return ""
}

func readByGeneratedColumn(tbl *Table, column string) bool {
	ref := toLower(quoteIdent(column))
	for _, col := range tbl.Columns {
		if col.Generated != nil && strings.Contains(toLower(col.Generated.Expr), ref) {
			return true
		}
	}
	return false
}

// integerWidth orders the integer types by width.
var integerWidth = map[string]int{"tinyint": 1, "smallint": 2, "mediumint": 3, "int": 4, "bigint": 8}

// widens reports a type change TiDB makes without rewriting rows: to a wider
// integer of the same signedness, or to a longer VARCHAR of the same
// collation.
func widens(a, b *Column) bool {
	unsigned := func(c *Column) bool { return strings.Contains(toLower(c.ColumnType), "unsigned") }
	wa, aInt := integerWidth[toLower(a.DataType)]
	wb, bInt := integerWidth[toLower(b.DataType)]
	if aInt && bInt {
		return wb >= wa && unsigned(a) == unsigned(b)
	}
	if strings.EqualFold(a.DataType, "varchar") && strings.EqualFold(b.DataType, "varchar") &&
		strings.EqualFold(a.Charset, b.Charset) && strings.EqualFold(a.Collation, b.Collation) {
		return typeLength(b.ColumnType) >= typeLength(a.ColumnType)
	}
	return false
}

// typeLength is the length in a column type such as varchar(20).
func typeLength(columnType string) int {
	n := 0
	if open := strings.IndexByte(columnType, '('); open >= 0 {
		for _, r := range columnType[open+1:] {
			if r < '0' || r > '9' {
				break
			}
			n = n*10 + int(r-'0')
		}
	}
	return n
}

func inClusteredPrimaryKey(tbl *Table, column string) bool {
	for _, idx := range tbl.Indexes {
		if !idx.Primary || !primaryKeyClustered(tbl, idx) {
			continue
		}
		for _, ic := range idx.Columns {
			if strings.EqualFold(ic.Name, column) {
				return true
			}
		}
	}
	return false
}

// alterTableOptions changes each option that differs, one statement each.
func (g *generator) alterTableOptions(from, to *Table) {
	set := func(clause, warning string) { g.add(OpAlterTable, to, alterTable(to, clause), warning) }
	fromCharset, fromCollation := tableCharsetCollation(from)
	toCharset, toCollation := tableCharsetCollation(to)
	if !strings.EqualFold(fromCharset, toCharset) || !strings.EqualFold(fromCollation, toCollation) {
		warning := ""
		if !strings.EqualFold(fromCharset, toCharset) && !(strings.EqualFold(fromCharset, "utf8") && strings.EqualFold(toCharset, "utf8mb4")) {
			warning = "TiDB changes a table's charset only from utf8 to utf8mb4 (8200)"
		}
		set("DEFAULT CHARSET="+toCharset+" COLLATE="+toCollation, warning)
	}
	if from.Comment != to.Comment {
		set("COMMENT='"+outputFormat(to.Comment)+"'", "")
	}
	if from.ShardRowIDBits != to.ShardRowIDBits {
		set(fmt.Sprintf("SHARD_ROW_ID_BITS=%d", to.ShardRowIDBits), "")
	}
	if from.PreSplitRegions != to.PreSplitRegions {
		set(fmt.Sprintf("PRE_SPLIT_REGIONS=%d", to.PreSplitRegions), "TiDB cannot change PRE_SPLIT_REGIONS of an existing table (8200)")
	}
	if from.AutoIDCache != to.AutoIDCache {
		set(fmt.Sprintf("AUTO_ID_CACHE=%d", to.AutoIDCache), "")
	}
	if from.Compression != to.Compression {
		set("COMPRESSION='"+outputFormat(to.Compression)+"'", "")
	}
	if !strings.EqualFold(from.PlacementPolicy, to.PlacementPolicy) {
		policy := "DEFAULT"
		if to.PlacementPolicy != "" {
			policy = quoteIdent(to.PlacementPolicy)
		}
		set("PLACEMENT POLICY="+policy, "")
	}
	if tiflashReplica(from) != tiflashReplica(to) {
		clause := tiflashReplica(to)
		if clause == "" {
			clause = " SET TIFLASH REPLICA 0"
		}
		set(strings.TrimPrefix(clause, " "), "")
	}
}

func (g *generator) alterTTL(from, to *Table, set func(clause, warning string)) {
	if to.TTLColumn == "" {
		if from.TTLColumn != "" {
			set("REMOVE TTL", "")
		}
		return
	}
	if from.TTLColumn == "" || !strings.EqualFold(from.TTLColumn, to.TTLColumn) || ttlInterval(from.TTLInterval) != ttlInterval(to.TTLInterval) {
		set("TTL="+quoteIdent(to.TTLColumn)+" + "+ttlInterval(to.TTLInterval), "")
	}
	if from.TTLColumn == "" && !to.TTLEnable || from.TTLColumn != "" && from.TTLEnable != to.TTLEnable {
		enable := "ON"
		if !to.TTLEnable {
			enable = "OFF"
		}
		set("TTL_ENABLE='"+enable+"'", "")
	}
	if ttlJobInterval(from) != ttlJobInterval(to) {
		set("TTL_JOB_INTERVAL='"+ttlJobInterval(to)+"'", "")
	}
}

func ttlJobInterval(t *Table) string {
	if t.TTLColumn == "" || t.TTLJobInterval != "" {
		return t.TTLJobInterval
	}
	return "1h"
}

// partitionChange is the clause that turns one partitioning into another:
// ADD PARTITION when RANGE or LIST partitions are only appended, REMOVE
// PARTITIONING when none is left, and the whole new clause otherwise, which
// TiDB applies by reorganizing the rows.
func partitionChange(from, to *PartitionInfo) string {
	if to == nil {
		return "REMOVE PARTITIONING"
	}
	if from != nil && appendsPartitions(from, to) {
		defs := make([]string, 0, len(to.Partitions)-len(from.Partitions))
		for _, pd := range to.Partitions[len(from.Partitions):] {
			defs = append(defs, showPartitionDefinition(to, pd))
		}
		return "ADD PARTITION (" + strings.Join(defs, ", ") + ")"
	}
	return strings.TrimPrefix(showPartitioning(to), "\n")
}

func appendsPartitions(from, to *PartitionInfo) bool {
	if !strings.HasPrefix(to.Type, "RANGE") && !strings.HasPrefix(to.Type, "LIST") {
		return false
	}
	if from.IntervalText != "" || to.IntervalText != "" || len(to.Partitions) <= len(from.Partitions) {
		return false
	}
	head := *to
	head.Partitions = to.Partitions[:len(from.Partitions)]
	return showPartitioning(from) == showPartitioning(&head)
}

// sortedTables lists a catalog's base tables by database and name.
func sortedTables(c *Catalog) []*Table {
	m := tablesOf(c)
	keys := make([]objectKey, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].db != keys[j].db {
			return keys[i].db < keys[j].db
		}
		return keys[i].name < keys[j].name
	})
	out := make([]*Table, len(keys))
	for i, k := range keys {
		out[i] = m[k]
	}
	return out
}
