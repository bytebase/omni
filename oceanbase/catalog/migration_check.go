package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — CHECK constraints (8.0.16+ only).
//
// This is the MySQL analog of PG's CHECK generator. It turns the check diff (TableDiffEntry.Checks,
// populated by diff_check.go) into ALTER TABLE ... ADD/DROP CHECK DDL. It is wired into
// GenerateMigrationWithNormalizer (migration.go) against the OpAddCheck/OpDropCheck op-types.
//
// VERSION GATE. CHECK is 8.0.16+ only; on 5.7 it is parsed-and-dropped (normalize.go
// CheckSupported). diffChecks already returns nil on 5.7, so the diff carries no check entries and
// this generator would emit nothing anyway — but it gates explicitly too, so the new-table path
// (which reads To.Constraints directly, not the diff) also emits nothing on 5.7. The gate is the
// single owner of the version decision (CheckSupported), never re-derived here.
//
// Two table cases produce check ADDs:
//   - DiffModify: the per-check diff drives ADD / DROP / (DROP+ADD for MODIFY).
//   - DiffAdd: a freshly created table. generate-core's formatCreateTable renders only columns and
//     the inline PRIMARY KEY (migration_table.go) — it does NOT inline CHECK constraints — so every
//     check on a new table is added here via ALTER TABLE ... ADD CONSTRAINT ... CHECK, mirroring how
//     the index node adds a new table's non-PK indexes (addIndexOpsForNewTable).
//
// DiffDrop tables emit nothing: DROP TABLE removes their checks wholesale.
//
// A check cannot be altered in place, so a modified check (changed expression or ENFORCED state) is
// a DROP followed by an ADD. The DROP runs in PhasePre and the ADD in PhaseMain, so a drop-then-add
// of the same name never collides on a name still held by the old definition.
//
// Rendering routes through the same canonical form show.go uses for the stored CHECK line
// (formatCheckClause mirrors showConstraint's ConCheck branch: CONSTRAINT `name` CHECK (expr) with
// a /*!80016 NOT ENFORCED */ marker), so the readback of the emitted DDL canonicalizes equal to
// `to` and apply-correctness holds.
//
// Ordering: check ADDs run in PhaseMain at priorityConstraint (after the table CREATE and column
// ALTERs — so a CHECK that references a freshly added/modified column applies against the final
// column). Check DROPs run in PhasePre (before PhaseMain), so a modified check's drop precedes its
// re-add. sortMigrationOps re-imposes the global phase/priority/name order.
func generateCheckDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	// Version gate: no CHECK in 5.7 stored form, so emit nothing (the new-table path below reads
	// To.Constraints directly, so this guard is what keeps it silent on 5.7).
	if !n.CheckSupported() {
		return nil
	}

	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, addCheckOpsForNewTable(entry)...)
		case DiffModify:
			ops = append(ops, checkOpsForModifiedTable(entry)...)
		case DiffDrop:
			// DROP TABLE removes the checks; nothing to emit.
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			// Drops ahead of adds locally (also enforced by phase) for readability.
			return checkOpRank(ops[i].Type) < checkOpRank(ops[j].Type)
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addCheckOpsForNewTable emits ADD ops for every CHECK constraint on a freshly created table
// (formatCreateTable inlines only columns and the PK, never CHECKs). The checks are ordered by
// lower-cased name so a new table's ADD CHECK sequence is stable across runs.
func addCheckOpsForNewTable(entry *TableDiffEntry) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, con := range orderedChecks(entry.To) {
		ops = append(ops, addCheckOp(entry, table, con))
	}
	return ops
}

// checkOpsForModifiedTable emits ADD / DROP / (DROP+ADD) ops from the per-check diff of a modified
// table. A modified check is a DROP (PhasePre) then an ADD (PhaseMain) — a check has no in-place
// alter, so the drop frees the name for the re-add.
func checkOpsForModifiedTable(entry *TableDiffEntry) []MigrationOp {
	if entry.To == nil || len(entry.Checks) == 0 {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, ce := range entry.Checks {
		switch ce.Action {
		case DiffAdd:
			if ce.To == nil {
				continue
			}
			ops = append(ops, addCheckOp(entry, table, ce.To))
		case DiffDrop:
			if ce.From == nil {
				continue
			}
			ops = append(ops, dropCheckOp(entry, table, ce.From))
		case DiffModify:
			if ce.From == nil || ce.To == nil {
				continue
			}
			ops = append(ops, dropCheckOp(entry, table, ce.From))
			ops = append(ops, addCheckOp(entry, table, ce.To))
		}
	}
	return ops
}

// addCheckOp builds an ALTER TABLE ... ADD CONSTRAINT `name` CHECK (...) op, running in PhaseMain at
// priorityConstraint (after the table CREATE at priorityTable and column ALTERs at priorityColumn,
// so a CHECK referencing a newly added/modified column applies against the final column).
func addCheckOp(entry *TableDiffEntry, table string, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpAddCheck,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s ADD %s", table, formatCheckClause(con)),
		Phase:        PhaseMain,
		Priority:     priorityConstraint,
		sortName:     checkSortName(entry.Database, entry.Name, con.Name),
	}
}

// priorityCheckDrop orders a CHECK DROP within PhasePre BEFORE any column DROP. MySQL CASCADE-DROPS
// a CHECK constraint when a column the constraint references is dropped (verified against live 8.0:
// `ALTER TABLE ... DROP COLUMN a` silently removes a `CHECK (a > 0)`), after which an explicit
// `DROP CHECK` for that constraint fails with errno 3821 ("Check constraint is not found in the
// table"). Because the migration plan is applied one statement at a time, the check DROP must
// precede the column DROP so the explicit DROP CHECK runs while the constraint still exists; the
// column DROP then proceeds with no constraint left to cascade. Column drops run at priorityColumn-1
// (generated) / priorityColumn (plain) within PhasePre, and table drops at priorityTable=10; this
// sits below the column-drop range and above the table-drop (matching the index node's
// priorityIndexDrop), so the order is table-drop (10) < check-drop (15) < generated-col-drop (19) <
// plain-col-drop (20). Dropping a check whose columns are NOT being dropped is harmless to sequence
// first (it is independent), so always ordering check drops ahead of column drops is safe and simple
// — and it also keeps a modified check's drop ahead of its PhaseMain re-add.
const priorityCheckDrop = priorityColumn - 5

// dropCheckOp builds an ALTER TABLE ... DROP CHECK `name` op (the 8.0.16+ syntax for removing a
// check constraint), running in PhasePre so all check drops precede any add (a modified check's
// drop frees its name before the re-add in PhaseMain) AND precede column drops (so a check referencing
// a dropped column is removed explicitly before MySQL would cascade-drop it — see priorityCheckDrop).
func dropCheckOp(entry *TableDiffEntry, table string, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpDropCheck,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP CHECK %s", table, migrationQuoteIdent(con.Name)),
		Phase:        PhasePre,
		Priority:     priorityCheckDrop,
		sortName:     checkSortName(entry.Database, entry.Name, con.Name),
	}
}

// formatCheckClause renders the CHECK clause that follows ADD in an ALTER TABLE, in MySQL's
// canonical stored form: CONSTRAINT `name` CHECK (expr) with an optional /*!80016 NOT ENFORCED */
// marker. It mirrors show.go's showConstraint ConCheck branch (the stored CHECK line) so the
// readback of the ADD canonicalizes equal to `to`. The expression is emitted as stored on the
// constraint (CheckExpr); the diff's canonical comparison — not the emitted text — is what absorbs
// surface-form differences, and the engine re-canonicalizes the emitted expression on apply.
//
// The /*!80016 ... */ version guard makes NOT ENFORCED a no-op below 8.0.16, matching how show.go
// renders it, so the same DDL is valid across the 8.0 line.
func formatCheckClause(con *Constraint) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CONSTRAINT %s CHECK (%s)", migrationQuoteIdent(con.Name), con.CheckExpr)
	if con.NotEnforced {
		b.WriteString(" /*!80016 NOT ENFORCED */")
	}
	return b.String()
}

// orderedChecks returns a table's CHECK constraints in a deterministic order (by lower-cased name),
// so a new table's ADD CHECK sequence is stable across runs.
func orderedChecks(t *Table) []*Constraint {
	var out []*Constraint
	if t == nil {
		return out
	}
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConCheck || con.Name == "" {
			continue
		}
		out = append(out, con)
	}
	sort.Slice(out, func(i, j int) bool {
		return toLower(out[i].Name) < toLower(out[j].Name)
	})
	return out
}

// checkOpRank orders check op types so drops sort ahead of adds locally (phase ordering also
// enforces this globally).
func checkOpRank(t MigrationOpType) int {
	if t == OpDropCheck {
		return 0
	}
	return 1
}

// checkSortName is the stable secondary sort key for a check op: lower-cased
// database.table.constraint.
func checkSortName(database, table, check string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(check)
}
