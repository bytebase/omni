package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// SDL generate — foreign keys.
//
// It turns the FK diff (TableDiffEntry.ForeignKeys, populated by diff_foreign_key.go) into
// ALTER TABLE ... ADD/DROP FOREIGN KEY DDL. A foreign key cannot be altered in place, so a
// modified one is a DROP followed by an ADD.
//
// Phase ordering. FK ADDs and FK DROPs sit in OPPOSITE phases:
//
//   - FK ADDs run in PhasePost (priorityForeignKey=99), after table, column and index DDL, so
//     every referenced table and column exists before the foreign key is created.
//   - FK DROPs run in PhasePre, BEFORE every dependent drop: a referenced table, a
//     referencing column, and an index the foreign key holds. FK constraint drops therefore
//     use priorityForeignKeyConstraintDrop, ahead of the table-drop / index-drop /
//     column-drop priorities the other generators use in PhasePre.
//
// No index belongs to a foreign key. OceanBase makes none, where MySQL adds one when no
// index serves the key, so dropping a foreign key drops nothing else
// (TestOracle_OceanBaseForeignKeyIndex). What OceanBase does do is hold, for a foreign key
// that lives, the last index over exactly its columns: DROP INDEX answers 5320, and DROP
// PRIMARY KEY 4714. A foreign key the plan otherwise leaves alone is therefore dropped first
// and added back when the plan drops that index (foreignKeysHoldingDroppedIndex).
//
// Rendering routes through the same canonical form show.go uses for the stored FK line, so the
// readback of the emitted DDL canonicalizes equal to `to` and apply-correctness holds.
func generateForeignKeyDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, addForeignKeyOpsForNewTable(entry)...)
		case DiffModify:
			ops = append(ops, foreignKeyOpsForModifiedTable(entry, n)...)
		case DiffDrop:
			// The table itself is being dropped, but its FKs must be RELEASED first: a FK on this
			// table that references ANOTHER table being dropped in the same plan would otherwise
			// block that table's DROP (errno 3730 / 1451) — table drops all share priorityTable, so
			// the referenced parent can sort before this child. Emit DROP FOREIGN KEY for every FK
			// here in PhasePre (priorityForeignKeyConstraintDrop), ahead of all table drops.
			ops = append(ops, dropForeignKeyConstraintOpsForDroppedTable(entry)...)
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// dropForeignKeyConstraintOpsForDroppedTable emits a DROP FOREIGN KEY for every FK on a table that
// is itself being DROPPED, in PhasePre ahead of all table drops. This releases the references so a
// referenced table dropped in the same plan is not blocked (errno 3730 / 1451).
func dropForeignKeyConstraintOpsForDroppedTable(entry *TableDiffEntry) []MigrationOp {
	if entry.From == nil {
		return nil
	}
	table := tableIdent(entry.From)
	var ops []MigrationOp
	for _, con := range orderedForeignKeys(entry.From) {
		ops = append(ops, MigrationOp{
			Type:         OpDropForeignKey,
			Database:     entry.Database,
			ObjectName:   entry.Name,
			ParentObject: entry.Name,
			SQL:          fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", table, migrationQuoteIdent(con.Name)),
			Phase:        PhasePre,
			Priority:     priorityForeignKeyConstraintDrop,
			sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name),
		})
	}
	return ops
}

// priorityForeignKeyConstraintDrop orders FK CONSTRAINT drops at the very front of PhasePre —
// ahead of table drops (priorityTable=10), index drops (priorityIndexDrop=15), and column drops
// (priorityColumn-1=19 / priorityColumn=20). A live FK blocks dropping what it needs (a
// referenced table, a referencing column, the last index over its columns), so every FK
// constraint is released before any of those drops run. It is negative so it sorts before
// priorityTable without renumbering the shared constants.
const priorityForeignKeyConstraintDrop = -10

// addForeignKeyOpsForNewTable emits ADD CONSTRAINT ... FOREIGN KEY ops for every FK on a freshly
// created table. formatCreateTable (migration_table.go) renders only columns and the inline
// PRIMARY KEY — it does NOT inline FK constraints — so every FK on a new table is added here in
// PhasePost, after the table (and any referenced tables) exist.
func addForeignKeyOpsForNewTable(entry *TableDiffEntry) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, con := range orderedForeignKeys(entry.To) {
		ops = append(ops, addForeignKeyOp(entry, table, con))
	}
	return ops
}

// foreignKeyOpsForModifiedTable emits ADD / DROP / (DROP+ADD) ops from the per-FK diff of a
// modified table. A MODIFY is a DROP followed by an ADD (an FK cannot be altered in place): the
// DROP runs in PhasePre (priorityForeignKeyConstraintDrop) and the ADD in PhasePost
// (priorityForeignKey), so the old constraint is fully released before the new one is created.
// A foreign key the diff leaves alone is released and added back the same way when the plan
// drops the index it holds.
func foreignKeyOpsForModifiedTable(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	var ops []MigrationOp
	for _, fe := range entry.ForeignKeys {
		switch fe.Action {
		case DiffAdd:
			if fe.To == nil || entry.To == nil {
				continue
			}
			ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), fe.To))
		case DiffDrop:
			if fe.From == nil || entry.From == nil {
				continue
			}
			ops = append(ops, dropForeignKeyOp(entry, fe.From))
		case DiffModify:
			if fe.From == nil || fe.To == nil || entry.From == nil || entry.To == nil {
				continue
			}
			ops = append(ops, dropForeignKeyOp(entry, fe.From))
			ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), fe.To))
		}
	}
	for _, con := range foreignKeysHoldingDroppedIndex(entry, n) {
		ops = append(ops, dropForeignKeyOp(entry, con))
		ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), con))
	}
	return ops
}

// addForeignKeyOp builds an ALTER TABLE ... ADD CONSTRAINT `name` FOREIGN KEY (...) REFERENCES
// ... op. It runs in PhasePost at priorityForeignKey so it lands after table/column/index DDL:
// every referenced table and column exists.
func addForeignKeyOp(entry *TableDiffEntry, table string, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpAddForeignKey,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s ADD %s", table, formatForeignKeyDefinition(con)),
		Phase:        PhasePost,
		Priority:     priorityForeignKey,
		sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name),
	}
}

// dropForeignKeyOp builds ALTER TABLE ... DROP FOREIGN KEY `name`, in PhasePre at
// priorityForeignKeyConstraintDrop: ahead of every drop the foreign key would block.
func dropForeignKeyOp(entry *TableDiffEntry, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpDropForeignKey,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", tableIdent(entry.From), migrationQuoteIdent(con.Name)),
		Phase:        PhasePre,
		Priority:     priorityForeignKeyConstraintDrop,
		sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name),
	}
}

// foreignKeysHoldingDroppedIndex returns the foreign keys of a modified table that the diff
// leaves alone and that hold an index the plan drops. For a foreign key that lives,
// OceanBase refuses to drop the last index over exactly its columns, in any order: DROP
// INDEX answers 5320 ('k': used in a foreign key constraint) and a DROP PRIMARY KEY of its
// own 4714. A wider index, a prefix index and DROP PRIMARY KEY, ADD PRIMARY KEY in one
// statement are not held. So when every such index goes, the foreign key is dropped first
// and added back after: the table then has it with no index, which OceanBase allows.
// TestOracle_OceanBaseForeignKeyIndex holds each of these to the engine.
func foreignKeysHoldingDroppedIndex(entry *TableDiffEntry, n *Normalizer) []*Constraint {
	if entry.From == nil || entry.To == nil {
		return nil
	}
	dropped := map[string]bool{}
	for _, ie := range entry.Indexes {
		if ie.From == nil {
			continue
		}
		switch ie.Action {
		case DiffDrop:
			dropped[toLower(ie.From.Name)] = true
		case DiffModify:
			// A primary key is replaced in one statement, and a change of visibility is an
			// ALTER INDEX: neither drops the index.
			if !ie.From.Primary && (ie.To == nil || !onlyVisibilityDiffers(ie.From, ie.To, n)) {
				dropped[toLower(ie.From.Name)] = true
			}
		}
	}
	if len(dropped) == 0 {
		return nil
	}
	changed := map[string]bool{}
	for _, fe := range entry.ForeignKeys {
		if fe.From != nil {
			changed[toLower(fe.From.Name)] = true
		}
	}
	var held []*Constraint
	for _, con := range orderedForeignKeys(entry.From) {
		if changed[toLower(con.Name)] {
			continue
		}
		holds, keeps := false, false
		for _, idx := range entry.From.Indexes {
			if idx == nil || !indexIsOverColumns(idx, con.Columns) {
				continue
			}
			if dropped[toLower(idx.Name)] {
				holds = true
			} else {
				keeps = true
			}
		}
		if holds && !keeps {
			held = append(held, con)
		}
	}
	return held
}

// indexIsOverColumns reports whether idx is an index over exactly cols, in any order: whole
// columns, no expression and no prefix.
func indexIsOverColumns(idx *Index, cols []string) bool {
	if idx.Fulltext || idx.Spatial || idx.Vector || len(idx.Columns) != len(cols) {
		return false
	}
	want := make(map[string]bool, len(cols))
	for _, c := range cols {
		want[toLower(c)] = true
	}
	for _, ic := range idx.Columns {
		if ic.Expr != "" || ic.Length != 0 || !want[toLower(ic.Name)] {
			return false
		}
		delete(want, toLower(ic.Name))
	}
	return len(want) == 0
}

// orderedForeignKeys returns a table's FOREIGN KEY constraints in a deterministic order (by
// lower-cased constraint name) so a new table's ADD FOREIGN KEY sequence is stable across runs.
func orderedForeignKeys(t *Table) []*Constraint {
	if t == nil {
		return nil
	}
	var fks []*Constraint
	for _, con := range t.Constraints {
		if con != nil && con.Type == ConForeignKey {
			fks = append(fks, con)
		}
	}
	sort.Slice(fks, func(i, j int) bool {
		return toLower(fks[i].Name) < toLower(fks[j].Name)
	})
	return fks
}

// formatForeignKeyDefinition renders the constraint clause that follows ADD in an ALTER TABLE, in
// MySQL's canonical stored form:
//
//	CONSTRAINT `name` FOREIGN KEY (`c1`, `c2`) REFERENCES `db`.`tbl` (`r1`, `r2`)
//	    [ON DELETE <action>] [ON UPDATE <action>]
//
// It mirrors show.go's showConstraint (the stored FK line) so the readback of the ADD
// canonicalizes equal to `to`. The ON DELETE / ON UPDATE clauses are emitted ONLY for a
// non-default action (CASCADE / SET NULL / SET DEFAULT): RESTRICT and NO ACTION (and an absent
// clause) are the engine's default and are omitted, so the emitted DDL round-trips to the same
// stored form the canonical key folds to (see diff_foreign_key.go canonicalFKAction). Emitting an
// explicit RESTRICT/NO ACTION would also round-trip on equality, but omitting matches show.go and
// keeps the surface minimal.
//
// A name nobody wrote is left out, so OceanBase makes its own (formatCheckClause).
func formatForeignKeyDefinition(con *Constraint) string {
	var b strings.Builder
	if !madeUpForeignKeyName(con) {
		fmt.Fprintf(&b, "CONSTRAINT %s ", migrationQuoteIdent(con.Name))
	}
	b.WriteString("FOREIGN KEY (")
	b.WriteString(quoteIdentList(con.Columns))
	b.WriteString(") REFERENCES ")
	if con.RefDatabase != "" {
		b.WriteString(migrationQuoteIdent(con.RefDatabase) + "." + migrationQuoteIdent(con.RefTable))
	} else {
		b.WriteString(migrationQuoteIdent(con.RefTable))
	}
	b.WriteString(" (")
	b.WriteString(quoteIdentList(con.RefColumns))
	b.WriteString(")")

	if con.OnDelete != "" && !isFKDefaultAction(con.OnDelete) {
		fmt.Fprintf(&b, " ON DELETE %s", strings.ToUpper(con.OnDelete))
	}
	if con.OnUpdate != "" && !isFKDefaultAction(con.OnUpdate) {
		fmt.Fprintf(&b, " ON UPDATE %s", strings.ToUpper(con.OnUpdate))
	}
	return b.String()
}

// isFKDefaultAction reports whether a referential action is a MySQL default that is omitted from
// the rendered FK clause. RESTRICT, NO ACTION, and an absent clause are all the "default" (no
// referential action) and are omitted — matching the canonical-equality collapse in
// canonicalFKAction, so the emitted DDL reads back to a form that diffs empty against `to` on both
// versions (8.0 omits NO ACTION but would echo RESTRICT; 5.7 omits RESTRICT but would echo NO
// ACTION — omitting BOTH from what we emit sidesteps the per-version echo entirely).
func isFKDefaultAction(action string) bool {
	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "", "RESTRICT", "NO ACTION":
		return true
	default:
		return false
	}
}

// quoteIdentList backtick-quotes each identifier and joins with ", " (the spacing show.go's FK
// line uses), preserving order.
func quoteIdentList(names []string) string {
	if len(names) == 0 {
		return ""
	}
	quoted := make([]string, len(names))
	for i, nm := range names {
		quoted[i] = migrationQuoteIdent(nm)
	}
	return strings.Join(quoted, ", ")
}

// foreignKeySortName is the stable secondary sort key for a FK op: lower-cased
// database.table.constraint.
func foreignKeySortName(database, table, constraint string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(constraint)
}
