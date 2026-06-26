package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — foreign keys.
//
// This is the MySQL analog of PG's foreign-key generator. It turns the FK diff
// (TableDiffEntry.ForeignKeys, populated by diff_foreign_key.go) into ALTER TABLE ... ADD/DROP
// FOREIGN KEY DDL. MySQL folds FK changes into ALTER TABLE; a FK cannot be altered in place, so
// a modified FK is a DROP followed by an ADD.
//
// Phase ordering — the crux of coordinating with the index node:
//   - FK ADDs run in PhasePost (priorityForeignKey=99), AFTER the index node's PhaseMain index
//     ADDs (priorityIndex=30) and after table/column DDL. MySQL's ADD CONSTRAINT ... FOREIGN KEY
//     auto-creates a backing index ONLY when no covering index already exists; by deferring the
//     FK add to PhasePost, an index the user declared (added by the index node in PhaseMain) is
//     in place first and is REUSED — no duplicate index, no errno 1061/1553. It also guarantees
//     every referenced table/column exists before the FK is created (the reason PG defers FKs).
//   - FK DROPs run in PhasePost too, but BEFORE the adds (priorityForeignKeyDrop) so a
//     drop-then-add of the same constraint name never collides, and before the leftover backing
//     index drop (errno 1553 — see below).
//
// Backing-index lifecycle on DROP — owned HERE, by contract with the index node:
//
//	When a FK is dropped, MySQL LEAVES the auto-created backing index behind (verified on both
//	live engines: after `ALTER TABLE c DROP FOREIGN KEY fk`, `KEY fk (pid)` remains). The index
//	node EXCLUDES FK-implicit backing indexes from its diff (diff_index.go fkImplicitIndexNames),
//	so it emits no DROP for that leftover — this node owns it. We drop the leftover index iff it
//	was FK-implicit on the `from` side (plain, auto-named after the constraint/first column) AND
//	the user's target has no USER index of that name (so the index node is not keeping it as a
//	real index). The DROP INDEX is emitted AFTER the DROP FOREIGN KEY, because MySQL refuses to
//	drop an index still needed by a FK (errno 1553).
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
			ops = append(ops, foreignKeyOpsForModifiedTable(entry)...)
		case DiffDrop:
			// DROP TABLE removes the FKs (and their backing indexes); nothing to emit.
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

// priorityForeignKeyDrop orders FK DROPs (and the leftover-backing-index DROPs they own) within
// PhasePost BEFORE the FK ADDs (priorityForeignKey=99). A drop-then-add of the same constraint
// name (a MODIFY) must drop first so the name is free; and an FK that moves between tables in the
// same plan must release the old one before the new one is added. It sits just below
// priorityForeignKey so all FK releases precede all FK creations.
const priorityForeignKeyDrop = priorityForeignKey - 1

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
// modified table. A MODIFY is a DROP followed by an ADD (an FK cannot be altered in place); the
// drop in priorityForeignKeyDrop always precedes the add in priorityForeignKey, freeing the
// constraint name for re-creation.
func foreignKeyOpsForModifiedTable(entry *TableDiffEntry) []MigrationOp {
	if len(entry.ForeignKeys) == 0 {
		return nil
	}
	var ops []MigrationOp
	for _, fe := range entry.ForeignKeys {
		switch fe.Action {
		case DiffAdd:
			if fe.To == nil || entry.To == nil {
				continue
			}
			ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), fe.To))
		case DiffDrop:
			if fe.From == nil {
				continue
			}
			ops = append(ops, dropForeignKeyOps(entry, fe.From)...)
		case DiffModify:
			if fe.From == nil || fe.To == nil || entry.To == nil {
				continue
			}
			// DROP the old FK (and its leftover backing index, if FK-implicit and not reused),
			// then ADD the new one. The drop's PhasePost priority precedes the add's.
			ops = append(ops, dropForeignKeyOps(entry, fe.From)...)
			ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), fe.To))
		}
	}
	return ops
}

// addForeignKeyOp builds an ALTER TABLE ... ADD CONSTRAINT `name` FOREIGN KEY (...) REFERENCES
// ... op. It runs in PhasePost at priorityForeignKey so it lands after table/column/index DDL —
// every referenced table exists and any user-declared covering index is already present (so MySQL
// reuses it instead of auto-creating a duplicate backing index).
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

// dropForeignKeyOps builds the ops to drop a foreign key: an ALTER TABLE ... DROP FOREIGN KEY
// `name`, plus — when the FK leaves behind an auto-created backing index that the target no
// longer keeps — an ALTER TABLE ... DROP INDEX `name` for that leftover. Both run in PhasePost at
// priorityForeignKeyDrop; the DROP FOREIGN KEY is emitted first and the DROP INDEX second
// (errno 1553 — MySQL refuses to drop an index a FK still needs).
//
// The leftover index is dropped ONLY when it was FK-implicit on the `from` side (plain,
// auto-named after the constraint / first FK column, reusing the index node's exact detection)
// AND the `to` table has no USER index of that name. This is the mirror of the index node's
// FK-implicit exclusion (diff_index.go diffableIndexMap with userIndexNameSet(to)): the index
// node keeps a user-named index when its FK is dropped, so this node must NOT drop such an index;
// it drops only the engine-synthesized backing index the index node deliberately ignores.
func dropForeignKeyOps(entry *TableDiffEntry, con *Constraint) []MigrationOp {
	if entry.From == nil {
		return nil
	}
	table := tableIdent(entry.From)
	ops := []MigrationOp{{
		Type:         OpDropForeignKey,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", table, migrationQuoteIdent(con.Name)),
		Phase:        PhasePost,
		Priority:     priorityForeignKeyDrop,
		sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name) + ".0",
	}}

	if idxName, ok := leftoverBackingIndexToDrop(entry, con); ok {
		ops = append(ops, MigrationOp{
			Type:         OpDropForeignKey,
			Database:     entry.Database,
			ObjectName:   entry.Name,
			ParentObject: entry.Name,
			SQL:          fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, migrationQuoteIdent(idxName)),
			Phase:        PhasePost,
			Priority:     priorityForeignKeyDrop,
			// ".1" sorts after the DROP FOREIGN KEY ".0" for the same constraint (errno 1553).
			sortName: foreignKeySortName(entry.Database, entry.Name, con.Name) + ".1",
		})
	}
	return ops
}

// leftoverBackingIndexToDrop decides whether dropping `con` (a FK on the `from` table) leaves an
// auto-created backing index this node must drop, and returns its actual (original-case) name.
// It returns the index iff:
//   - the `from` table has an FK-implicit backing index for `con` (the index node's exact
//     detection: fkImplicitIndexNames over `from`), AND
//   - the `to` table has no USER index of that name (so the index node is not keeping it).
//
// When the `to` table keeps a user index of that name, the index node owns it and MySQL leaves it
// intentionally; this node must not drop it. When `to` is absent (the whole table is being
// dropped) this is never reached — DiffDrop tables emit nothing.
func leftoverBackingIndexToDrop(entry *TableDiffEntry, con *Constraint) (string, bool) {
	from := entry.From
	if from == nil {
		return "", false
	}
	implicit := fkImplicitIndexNames(from)
	// Find the backing-index name for THIS constraint: the FK-implicit index whose name matches
	// the auto-derived form for `con` (constraint name / first-column / collision suffix).
	var backing *Index
	for _, idx := range from.Indexes {
		if idx == nil {
			continue
		}
		if !implicit[toLower(idx.Name)] {
			continue
		}
		if isAutoBackingIndexName(idx.Name, con) && isPlainBackingIndexFor(idx, con.Columns) {
			backing = idx
			break
		}
	}
	if backing == nil {
		return "", false
	}
	// Keep it if the target still carries a USER index of that name (index node's domain).
	if userIndexNameSet(entry.To)[toLower(backing.Name)] {
		return "", false
	}
	return backing.Name, true
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
func formatForeignKeyDefinition(con *Constraint) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CONSTRAINT %s FOREIGN KEY (", migrationQuoteIdent(con.Name))
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
