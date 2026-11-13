package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateColumnDDL produces ALTER TABLE ... ADD/DROP/MODIFY COLUMN operations for every
// column change within modified tables. A column is re-declared whole with one MODIFY
// COLUMN, so one diff column entry maps to one op:
//
//   - DiffAdd    → ALTER TABLE ... ADD COLUMN <full canonical def>
//   - DiffDrop   → ALTER TABLE ... DROP COLUMN `name`
//   - DiffModify → ALTER TABLE ... MODIFY COLUMN <full canonical def>
//
// MODIFY (not CHANGE) is used because the column name is unchanged — diff-core keys columns
// by name, so a rename is seen as DROP+ADD, never a MODIFY. The full canonical definition
// routes through normalize-core (the same formatColumnDefinition CREATE TABLE uses), so the
// readback of the ALTER canonicalizes equal to `to`.
//
// Four things OceanBase does differently from MySQL shape the ops, each measured on the
// pinned server and held to it by TestOracle_OceanBaseAlterTableClauses:
//
//   - One AUTO_INCREMENT column at a time (1075). When the attribute moves from one column
//     to another, the statement that takes it off runs first (priorityAutoIncrementRelease).
//   - A column of the primary key is not dropped (1235, drop rowkey column) and not made
//     nullable (1171) while the key holds it, and a column the partitioning reads is not
//     dropped (3855). Such a change waits for the statement that replaces the key or the
//     partitioning (columnsThatWait, priorityColumnAfterKeys).
//   - A spatial column becomes NOT NULL only from nullable (1138): ADD COLUMN ... NOT NULL
//     is refused even on an empty table, and so is any MODIFY of a column that is NOT NULL
//     already. It is added, or modified, nullable first and then made NOT NULL
//     (addColumnOps, modifyColumnOps). Its SRID does not change in place (1235), and the
//     MODIFY that asks for it carries a Warning.
//   - A generated column's expression does not change in place (1235): the column is
//     dropped and added again (recreatesColumn).
//
// Determinism: drops are emitted before adds/modifies (so a drop-then-add of the same name
// never collides), each sorted by lower-cased column name; sortMigrationOps re-imposes the
// global order.
func generateColumnDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		// Only MODIFY tables carry column changes (ADD tables ship columns inline in CREATE
		// TABLE; DROP tables vanish whole). diff-core only populates Columns on a DiffModify
		// entry, but guard the action explicitly so a future ADD-with-columns shape cannot
		// leak a duplicate ALTER.
		if entry.Action != DiffModify || len(entry.Columns) == 0 || entry.From == nil || entry.To == nil {
			continue
		}
		// Use the table's original-case identifier (a MODIFY table always has To). Lower-casing
		// the diff identity key would target the wrong object on a case-sensitive server.
		table := tableIdent(entry.To)
		waits := columnsThatWait(entry)
		// losing is the column AUTO_INCREMENT leaves for another one, "" when it does not move.
		losing := ""
		if from, to := toLower(autoIncrementColumnName(entry.From)), toLower(autoIncrementColumnName(entry.To)); to != "" && from != to {
			losing = from
		}

		for _, col := range entry.Columns {
			name := toLower(col.Name)
			switch col.Action {
			case DiffAdd:
				if col.To == nil {
					continue
				}
				ops = append(ops, addColumnOps(entry, table, col.To, n)...)

			case DiffDrop:
				if col.From == nil {
					continue
				}
				op := dropColumnOp(entry, table, col.Name, col.From)
				if waits.drop(name) {
					if name == losing {
						plain := *col.From
						plain.AutoIncrement = false
						ops = append(ops, releaseAutoIncrementOp(entry, table, entry.From, &plain, n))
					}
					op.Phase, op.Priority = PhaseMain, priorityColumnAfterKeys
				}
				ops = append(ops, op)

			case DiffModify:
				if col.From == nil || col.To == nil {
					continue
				}
				// A change MODIFY cannot make is a DROP and an ADD: the generated attribute
				// gained, lost or flipped between VIRTUAL and STORED, and on OceanBase a
				// changed expression.
				if recreatesColumn(entry, col, n) {
					ops = append(ops, dropColumnOp(entry, table, col.Name, col.From))
					ops = append(ops, addColumnOps(entry, table, col.To, n)...)
					continue
				}
				mods := modifyColumnOps(entry, table, col, n)
				switch {
				case waits.null(name) && col.To.Nullable:
					if name == losing {
						plain := *col.From
						plain.AutoIncrement = false
						ops = append(ops, releaseAutoIncrementOp(entry, table, entry.From, &plain, n))
					}
					for i := range mods {
						mods[i].Priority = priorityColumnAfterKeys
					}
				case name == losing:
					for i := range mods {
						mods[i].Priority = priorityAutoIncrementRelease
					}
				}
				ops = append(ops, mods...)
			}
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		// Drops ahead of adds and modifies, as the phases have them. Then lower Priority
		// first (a generated column drops before the plain ones it reads), then name. Ops
		// that tie keep the order they were made in: an ADD and the MODIFY that follows it.
		if ri, rj := columnOpRank(ops[i].Type), columnOpRank(ops[j].Type); ri != rj {
			return ri < rj
		}
		if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// priorityAutoIncrementRelease orders, within PhaseMain, the MODIFY that takes
// AUTO_INCREMENT off a column before every ADD and MODIFY of a column (priorityColumn),
// one of which puts it on another: OceanBase answers 1075, "there can be only one auto
// column", to the second while the first still has it.
const priorityAutoIncrementRelease = priorityColumn - 1

// priorityColumnAfterKeys orders, within PhaseMain, a column change that has to wait for
// the table's key and partition statements: after the index adds and the primary key's
// replacement (priorityIndex), which a repartition that lets a column go precedes
// (priorityRepartitionBeforeKeys).
const priorityColumnAfterKeys = priorityIndex + 1

// waitingColumns names the columns of a modified table that OceanBase lets go only after
// another statement of the plan. See columnsThatWait.
type waitingColumns struct {
	key       map[string]bool // columns of a primary key the plan replaces
	partition map[string]bool // columns the old partitioning reads, when the plan repartitions
	anyDrop   bool            // the old partitioning is an expression: any column may be in it
}

// drop reports whether the column's DROP waits.
func (w waitingColumns) drop(name string) bool {
	return w.anyDrop || w.key[name] || w.partition[name]
}

// null reports whether the column's change to nullable waits.
func (w waitingColumns) null(name string) bool { return w.key[name] }

// columnsThatWait finds the columns whose DROP, or whose change to nullable, must follow
// the plan's key and partition statements, because OceanBase refuses it while the old key
// or the old partitioning still reads the column:
//
//   - A column of the primary key: DROP COLUMN answers 1235 (drop rowkey column is not
//     supported) and a nullable MODIFY 1171. Once DROP PRIMARY KEY, ADD PRIMARY KEY has
//     replaced the key (primaryKeyModifyOp), both are taken. A key the plan drops outright
//     goes in PhasePre, before the column drops, and nothing waits for it.
//   - A column the partitioning reads: DROP COLUMN answers 3855 until the table is
//     partitioned by something else, and OceanBase has no REMOVE PARTITIONING to free it
//     first. The new PARTITION BY therefore runs before the key changes
//     (repartitionsBeforeKeys) and the drop after them.
func columnsThatWait(entry *TableDiffEntry) waitingColumns {
	var w waitingColumns
	for _, ie := range entry.Indexes {
		if ie.Action != DiffModify || ie.From == nil || !ie.From.Primary {
			continue
		}
		w.key = make(map[string]bool, len(ie.From.Columns))
		for _, ic := range ie.From.Columns {
			w.key[toLower(ic.Name)] = true
		}
	}
	if repartitionsBeforeKeys(entry) {
		pi := entry.From.Partitioning
		if pi.Expr != "" || pi.SubExpr != "" {
			// As oldPartitionColumnDropped: the columns of an expression are not read
			// from its text, so every drop waits. Waiting is harmless where it was not needed.
			w.anyDrop = true
		}
		w.partition = map[string]bool{}
		for _, c := range append(append([]string{}, pi.Columns...), pi.SubColumns...) {
			w.partition[toLower(c)] = true
		}
	}
	return w
}

// modifyColumnOp builds ALTER TABLE ... MODIFY COLUMN for col as tbl declares it, in
// PhaseMain at priorityColumn.
func modifyColumnOp(entry *TableDiffEntry, table string, tbl *Table, col *Column, n *Normalizer) MigrationOp {
	return MigrationOp{
		Type:         OpModifyColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s", table, formatColumnDefinition(tbl, col, n)),
		Phase:        PhaseMain,
		Priority:     priorityColumn,
		sortName:     columnSortName(entry.Database, entry.Name, col.Name),
	}
}

// releaseAutoIncrementOp builds the MODIFY that takes AUTO_INCREMENT off a column ahead of
// the statement that changes the column for good, for a column whose own change waits
// (columnsThatWait) while another column takes the attribute. plain is the column as it
// stands in between: as it was, without AUTO_INCREMENT.
func releaseAutoIncrementOp(entry *TableDiffEntry, table string, tbl *Table, plain *Column, n *Normalizer) MigrationOp {
	op := modifyColumnOp(entry, table, tbl, plain, n)
	op.Priority = priorityAutoIncrementRelease
	return op
}

// modifyColumnOps is the MODIFY COLUMN of a changed column, with its Warning when OceanBase
// has no statement for the change (modifyColumnWarning), and in two statements for a
// spatial column that is NOT NULL before and after. OceanBase answers 1138 (Invalid use of
// NULL value) to any MODIFY of such a column, whatever it changes, and takes one that makes
// a nullable column NOT NULL, so the column is modified as it is to be but nullable, and
// then made NOT NULL.
func modifyColumnOps(entry *TableDiffEntry, table string, col ColumnDiffEntry, n *Normalizer) []MigrationOp {
	op := modifyColumnOp(entry, table, entry.To, col.To, n)
	op.Warning = modifyColumnWarning(col, n)
	if op.Warning != "" || !isSpatialColumnType(col.From.DataType) || !isSpatialColumnType(col.To.DataType) ||
		!n.CanonicalNotNull(entry.From, col.From) || !n.CanonicalNotNull(entry.To, col.To) {
		return []MigrationOp{op}
	}
	nullable := *col.To
	nullable.Nullable = true
	return []MigrationOp{modifyColumnOp(entry, table, entry.To, &nullable, n), op}
}

// modifyColumnWarning returns the Warning of a MODIFY COLUMN that OceanBase is known to
// refuse and that the plan has no other statement for, "" for any other.
func modifyColumnWarning(col ColumnDiffEntry, n *Normalizer) string {
	switch {
	case isSpatialColumnType(col.From.DataType) && isSpatialColumnType(col.To.DataType) &&
		(col.From.HasSRID != col.To.HasSRID || col.From.SRID != col.To.SRID):
		return fmt.Sprintf("OceanBase cannot change the SRID of column %s in place (error 1235): drop the column and add it again", migrationQuoteIdent(col.Name))
	case generatedExprChanged(col.From, col.To, n):
		// recreatesColumn found the column read by something a DROP COLUMN would break.
		return fmt.Sprintf("OceanBase cannot change the expression of generated column %s in place (error 1235), and the plan cannot drop and add the column while a key, a constraint, another generated column or the partitioning reads it", migrationQuoteIdent(col.Name))
	}
	return ""
}

// addColumnOps is addColumnOp, in two statements for a spatial column that is NOT NULL.
// OceanBase answers 1138 (Invalid use of NULL value) to ADD COLUMN g GEOMETRY NOT NULL even
// on an empty table, and takes the column nullable and then a MODIFY to NOT NULL, which
// fails when the table has rows, as the ADD does on MySQL. Other types with no implicit
// default (JSON, BLOB, TEXT) are added NOT NULL as written.
func addColumnOps(entry *TableDiffEntry, table string, to *Column, n *Normalizer) []MigrationOp {
	if to.Generated != nil || !isSpatialColumnType(to.DataType) || !n.CanonicalNotNull(entry.To, to) {
		return []MigrationOp{addColumnOp(entry, table, to, n)}
	}
	nullable := *to
	nullable.Nullable = true
	add := addColumnOp(entry, table, &nullable, n)
	notNull := modifyColumnOp(entry, table, entry.To, to, n)
	notNull.Priority = add.Priority
	return []MigrationOp{add, notNull}
}

// isSpatialColumnType reports whether dataType is one of the spatial types.
func isSpatialColumnType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "geometry", "point", "linestring", "polygon", "multipoint", "multilinestring",
		"multipolygon", "geometrycollection", "geomcollection":
		return true
	}
	return false
}

// addColumnOp builds an ADD COLUMN op. A generated column is added AFTER plain columns
// (priorityColumn+1) so that any base column it references — whether freshly added or being
// MODIFYed in the same PhaseMain batch — is already in place; this is the symmetric partner of
// dropColumnOp's generated-first rule (generated columns drop before, and add after, the plain
// columns they depend on). The global sortMigrationOps honors Priority within PhaseMain.
func addColumnOp(entry *TableDiffEntry, table string, to *Column, n *Normalizer) MigrationOp {
	clause := fmt.Sprintf("ADD COLUMN %s", formatColumnDefinition(entry.To, to, n))
	return MigrationOp{
		Type:         OpAddColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhaseMain,
		Priority:     columnAddPriority(to),
		sortName:     columnSortName(entry.Database, entry.Name, to.Name),
	}
}

// columnAddPriority returns priorityColumn+1 for a generated column (add it after the plain
// columns its expression may reference) and priorityColumn for a plain column.
func columnAddPriority(c *Column) int {
	if c != nil && c.Generated != nil {
		return priorityColumn + 1
	}
	return priorityColumn
}

// dropColumnOp builds a DROP COLUMN op, giving a generated column a lower priority so it is
// dropped before the plain columns its expression may reference (dropping a column a
// generated column reads is error 3108). The global sortMigrationOps also honors Priority
// within PhasePre, so the ordering holds across the whole plan. generateColumnDDL moves the
// drop of a column that has to wait (columnsThatWait) to PhaseMain.
func dropColumnOp(entry *TableDiffEntry, table, colName string, from *Column) MigrationOp {
	clause := fmt.Sprintf("DROP COLUMN %s", migrationQuoteIdent(colName))
	return MigrationOp{
		Type:         OpDropColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhasePre,
		Priority:     columnDropPriority(from),
		sortName:     columnSortName(entry.Database, entry.Name, colName),
	}
}

// columnDropPriority returns priorityColumn-1 for a generated column (drop it first) and
// priorityColumn for a plain column, so a generated column that references a plain column is
// always dropped before that plain column.
func columnDropPriority(c *Column) int {
	if c != nil && c.Generated != nil {
		return priorityColumn - 1
	}
	return priorityColumn
}

// generatedShapeChanged reports whether a column's GENERATED shape changed in a way MySQL
// cannot express as an in-place MODIFY: gaining or losing the generated attribute, or flipping
// VIRTUAL↔STORED. Such a change must be rendered as DROP + ADD. A change to only the generation
// EXPRESSION (storage mode unchanged) is a normal MODIFY and is not flagged here.
func generatedShapeChanged(from, to *Column) bool {
	if from == nil || to == nil {
		return false
	}
	fromGen, toGen := from.Generated != nil, to.Generated != nil
	if fromGen != toGen {
		return true
	}
	if fromGen && toGen && from.Generated.Stored != to.Generated.Stored {
		return true
	}
	return false
}

// generatedExprChanged reports whether a generated column keeps its storage mode and
// changes its expression.
func generatedExprChanged(from, to *Column, n *Normalizer) bool {
	if from == nil || to == nil || from.Generated == nil || to.Generated == nil || from.Generated.Stored != to.Generated.Stored {
		return false
	}
	return n.CanonicalGeneratedExpr(from.Generated.Expr) != n.CanonicalGeneratedExpr(to.Generated.Expr)
}

// recreatesColumn reports whether a modified column is dropped and added again: when its
// generated shape changed, which no engine alters in place, and when only its expression
// did, which OceanBase refuses as well (1235, Modify generated column definition not
// supported) where MySQL takes a MODIFY. The values of a generated column are computed, so
// nothing is lost. An index over the column goes with the DROP and is added again by the
// index generator (recreatedColumns).
//
// An expression change is left a MODIFY, with a Warning, when something else reads the
// column that a DROP COLUMN would be refused for or would take along: the primary key, a
// foreign key, a CHECK, another generated column or the partitioning.
func recreatesColumn(entry *TableDiffEntry, col ColumnDiffEntry, n *Normalizer) bool {
	if generatedShapeChanged(col.From, col.To) {
		return true
	}
	return generatedExprChanged(col.From, col.To, n) && !columnIsReadElsewhere(entry.From, col.Name) && !columnIsReadElsewhere(entry.To, col.Name)
}

// recreatedColumns returns the lower-cased names of the columns of a modified table that
// the plan drops and adds again (recreatesColumn).
func recreatedColumns(entry *TableDiffEntry, n *Normalizer) map[string]bool {
	var set map[string]bool
	if entry.From == nil || entry.To == nil {
		return set
	}
	for _, col := range entry.Columns {
		if col.Action != DiffModify || col.From == nil || col.To == nil || !recreatesColumn(entry, col, n) {
			continue
		}
		if set == nil {
			set = map[string]bool{}
		}
		set[toLower(col.Name)] = true
	}
	return set
}

// columnIsReadElsewhere reports whether tbl's primary key, a foreign key, a CHECK, another
// generated column or the partitioning reads the column. Expressions are searched for the
// name as text, which can only say yes too often.
func columnIsReadElsewhere(tbl *Table, column string) bool {
	name := toLower(column)
	inText := func(expr string) bool { return strings.Contains(toLower(expr), name) }
	for _, idx := range tbl.Indexes {
		if idx == nil || !idx.Primary {
			continue
		}
		for _, ic := range idx.Columns {
			if toLower(ic.Name) == name {
				return true
			}
		}
	}
	for _, con := range tbl.Constraints {
		if con == nil {
			continue
		}
		switch con.Type {
		case ConForeignKey:
			for _, c := range con.Columns {
				if toLower(c) == name {
					return true
				}
			}
		case ConCheck:
			if inText(con.CheckExpr) {
				return true
			}
		}
	}
	for _, c := range tbl.Columns {
		if c != nil && c.Generated != nil && toLower(c.Name) != name && inText(c.Generated.Expr) {
			return true
		}
	}
	if pi := tbl.Partitioning; pi != nil {
		if inText(pi.Expr) || inText(pi.SubExpr) {
			return true
		}
		for _, c := range append(append(append([]string{}, pi.Columns...), pi.SubColumns...), pi.keyColumns...) {
			if toLower(c) == name {
				return true
			}
		}
	}
	return false
}

// columnOpRank orders column op types so drops sort ahead of adds/modifies locally.
func columnOpRank(t MigrationOpType) int {
	if t == OpDropColumn {
		return 0
	}
	return 1
}

// columnSortName is the stable secondary sort key for a column op: lower-cased
// database.table.column.
func columnSortName(database, table, column string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(column)
}

// formatColumnDefinition renders a single column in MySQL's canonical stored form for the
// target version: `name` <type> [CHARACTER SET .. COLLATE ..] [GENERATED ALWAYS AS (..)
// VIRTUAL|STORED] [NULL|NOT NULL] [DEFAULT ..] [AUTO_INCREMENT] [ON UPDATE ..] [COMMENT ..]
// [INVISIBLE] [/*!80003 SRID n */]. It is shared by CREATE TABLE, ADD COLUMN, and MODIFY
// COLUMN, so all three express a column identically and round-trip through the same
// normalize-core canonical form diff-core compares — no ad-hoc surface rebuilding that would
// re-introduce diff noise.
//
// The clause ORDER follows MySQL's own SHOW CREATE rendering (show.go), which is also the
// order MySQL accepts on input, so a rendered definition is both valid DDL and canonical.
func formatColumnDefinition(tbl *Table, c *Column, n *Normalizer) string {
	var b strings.Builder
	b.WriteString(migrationQuoteIdent(c.Name))
	b.WriteString(" ")
	b.WriteString(n.CanonicalColumnType(c))

	// CHARACTER SET / COLLATE — emit the resolved pair when it differs from the table's
	// effective pair, so the column's stored charset/collation matches `to` after readback. A
	// bare column (same effective pair as the table) emits nothing and inherits, exactly as
	// MySQL stores it; an explicit divergent charset/collation is rendered in full. Resolving
	// through ResolveColumnCharsetCollation (the same key diff-core compares) keeps the
	// version-flagged echo asymmetry out of the rendering.
	if cs, coll := n.ResolveColumnCharsetCollation(tbl, c); cs != "" {
		tcs := foldCharset(tbl.Charset)
		tcoll := n.effectiveTableCollation(tbl)
		if cs != tcs {
			b.WriteString(" CHARACTER SET ")
			b.WriteString(cs)
			// The column declares a DIFFERENT charset than the table, so it cannot inherit the
			// table COLLATE; a bare `CHARACTER SET cs` falls back to cs's SERVER default. Emit
			// COLLATE only when the target collation is NOT that charset's default — so a
			// column whose collation IS the default renders no COLLATE, while a non-default one
			// is explicit.
			if coll != "" && !n.isCharsetDefaultCollation(cs, coll) {
				b.WriteString(" COLLATE ")
				b.WriteString(coll)
			}
		} else if coll != "" && coll != tcoll {
			// Same charset as the table but a collation that differs from what the column would
			// INHERIT from the table (tcoll) → emit COLLATE. The comparison is against the
			// table's effective collation, NOT the charset default: when the table itself
			// carries a non-default collation, a column at the charset default still differs
			// from the table and MUST be rendered explicitly, or it would wrongly inherit the
			// table's collation on apply (non-convergent). tcoll is already the version-correct
			// effective table collation, so this stays version-robust.
			b.WriteString(" COLLATE ")
			b.WriteString(coll)
		}
	}

	// Generated columns: GENERATED ALWAYS AS (expr) VIRTUAL|STORED. A generated column cannot
	// carry a DEFAULT/AUTO_INCREMENT/ON UPDATE, so this branch renders nullability + comment +
	// invisibility and returns. The expression is rendered verbatim (its stored form); the
	// diff compares it via CanonicalGeneratedExpr, so a round-trip matches.
	if c.Generated != nil {
		mode := "VIRTUAL"
		if c.Generated.Stored {
			mode = "STORED"
		}
		b.WriteString(" GENERATED ALWAYS AS (")
		b.WriteString(c.Generated.Expr)
		b.WriteString(") ")
		b.WriteString(mode)
		if n.CanonicalNotNull(tbl, c) {
			b.WriteString(" NOT NULL")
		}
		writeColumnComment(&b, c)
		if c.Invisible {
			b.WriteString(" /*!80023 INVISIBLE */")
		}
		return b.String()
	}

	// Nullability. Render NOT NULL when the canonical (post-rewrite) state is NOT NULL. For a
	// nullable column we render NULL only for TIMESTAMP (MySQL's SHOW CREATE shows an explicit
	// NULL for a nullable TIMESTAMP, and rendering it keeps the EDFT-off magic from forcing
	// NOT NULL on apply); other nullable columns omit the keyword (the default).
	if n.CanonicalNotNull(tbl, c) {
		b.WriteString(" NOT NULL")
	} else if isTimestampType(c.DataType) {
		b.WriteString(" NULL")
	}

	// SRID (spatial) — MySQL places it after NOT NULL, before DEFAULT. Emit on PRESENCE so
	// an explicit `SRID 0` round-trips (SHOW CREATE echoes `/*!80003 SRID 0 */`); a no-SRID
	// column leaves HasSRID false and emits nothing.
	if c.HasSRID {
		fmt.Fprintf(&b, " /*!80003 SRID %d */", c.SRID)
	}

	// DEFAULT. Render the column's declared default in canonical surface form (reusing
	// show.go's formatDefault, which normalizes CURRENT_TIMESTAMP and quoting). A bare first
	// TIMESTAMP under EDFT=OFF has no declared default but the engine injects
	// DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP on apply; CanonicalTimestampDefaults
	// predicts the same for `to`, so omitting it here still round-trips. We render only what is
	// declared to avoid emitting an injected default that a re-apply would also inject.
	if c.Default != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(formatDefault(*c.Default, c))
	}

	// AUTO_INCREMENT.
	if c.AutoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	}

	// ON UPDATE.
	if c.OnUpdate != "" {
		b.WriteString(" ON UPDATE ")
		b.WriteString(formatOnUpdate(c.OnUpdate, c.OnUpdateKind))
	}

	writeColumnComment(&b, c)

	if c.Invisible {
		b.WriteString(" /*!80023 INVISIBLE */")
	}

	return b.String()
}

// writeColumnComment appends a COMMENT clause when the column has one, escaped the way
// MySQL's stored form does.
func writeColumnComment(b *strings.Builder, c *Column) {
	if c.Comment != "" {
		b.WriteString(" COMMENT ")
		b.WriteString(quoteStringLiteral(c.Comment))
	}
}
