package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — triggers (CREATE TRIGGER / DROP TRIGGER).
//
// Turns the trigger diff (SchemaDiff.Triggers, populated by diff_trigger.go) into DDL. MySQL has
// NO ALTER TRIGGER, so every modification is a DROP followed by a CREATE — diff_trigger.go already
// emits DiffModify for a changed trigger, and this generator renders it as the drop+create pair.
//
// Triggers are database-level objects, so unlike the index/partition generators this one reads
// SchemaDiff.Triggers directly (not the per-table entries). It is wired into
// GenerateMigrationWithNormalizer (migration.go) via the OpCreateTrigger / OpDropTrigger op-types
// at priorityTrigger.
//
// Ordering (sortMigrationOps imposes the global order):
//   - CREATE TRIGGER runs in PhaseMain at priorityTrigger (70). priorityTrigger > priorityTable
//     (10) and > priorityColumn (20), so a trigger is created AFTER its table's CREATE TABLE and
//     any column ALTERs in the same plan — the table the trigger binds to already exists.
//   - DROP TRIGGER runs in PhasePre, so it precedes a same-name re-create (a DiffModify) and any
//     table re-create, and is emitted BEFORE the PhaseMain CREATE half of a modify.
//
// Table-drop cascade: MySQL auto-drops a table's triggers when the table is dropped (verified on
// live 5.7.25 + 8.0.32). When a trigger's owning table is itself being dropped in this plan
// (absent in `to`), emitting an explicit DROP TRIGGER would fail — the table drop already removed
// it. dropTriggerSuppressed detects that case and skips the redundant op; the DROP TABLE
// (PhasePre, priorityTable) cascades to the trigger.
//
// Rendering uses the fully-qualified form CREATE TRIGGER db.trig ... ON db.table (and
// DROP TRIGGER db.trig), every identifier backtick-quoted, which is current-database-independent
// and matches what the engine accepts (verified on live engines). The body is emitted verbatim —
// MySQL stores it verbatim, so the readback of the emitted DDL canonicalizes equal to `to` and
// apply-correctness holds by construction.
//
// implemented by omni:triggers breadth node
func generateTriggerDDL(_, to *Catalog, diff *SchemaDiff, _ *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Triggers {
		te := &diff.Triggers[i]
		switch te.Action {
		case DiffAdd:
			if te.To == nil {
				continue
			}
			ops = append(ops, buildCreateTriggerOp(te.Database, te.To))
		case DiffDrop:
			if te.From == nil || dropTriggerSuppressed(to, te.From) {
				continue
			}
			ops = append(ops, buildDropTriggerOp(te.Database, te.From))
		case DiffModify:
			// No ALTER TRIGGER in MySQL: DROP old then CREATE new. The drop (PhasePre) always
			// precedes the create (PhaseMain), freeing the name for re-creation. The DROP half is
			// gated by the SAME table-drop-cascade suppression as a plain DiffDrop: a trigger whose
			// owning table is dropped in this plan (e.g. a trigger relocated to a new table while its
			// OLD table is dropped — identity is (database, name), so a table move is a MODIFY) has
			// already been cascade-removed by the DROP TABLE, so an explicit DROP TRIGGER on the old
			// table would fail (errno 1360). The CREATE half (on te.To's surviving table) still runs.
			if te.From != nil && !dropTriggerSuppressed(to, te.From) {
				ops = append(ops, buildDropTriggerOp(te.Database, te.From))
			}
			if te.To != nil {
				ops = append(ops, buildCreateTriggerOp(te.Database, te.To))
			}
		}
	}

	// Determinism: stable by (op-rank, sortName) so drops sort ahead of creates locally (phase
	// ordering also enforces this globally) and the emitted sequence is byte-stable across runs.
	sort.SliceStable(ops, func(i, j int) bool {
		if ri, rj := triggerOpRank(ops[i].Type), triggerOpRank(ops[j].Type); ri != rj {
			return ri < rj
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// dropTriggerSuppressed reports whether an explicit DROP TRIGGER for `trig` must be skipped
// because the trigger's owning table is being dropped in the same plan (absent in `to`), in
// which case DROP TABLE already cascades to the trigger and a standalone DROP TRIGGER would
// fail. The table is looked up in `to` under the trigger's own database; if the database or
// table is gone, the drop is cascaded and the op is suppressed.
func dropTriggerSuppressed(to *Catalog, trig *Trigger) bool {
	if to == nil || trig == nil || trig.Database == nil || trig.Table == "" {
		return false
	}
	db := to.GetDatabase(trig.Database.Name)
	if db == nil {
		return true // whole database gone → table (and its triggers) cascaded
	}
	return db.GetTable(trig.Table) == nil
}

// buildCreateTriggerOp renders a CREATE TRIGGER op in MySQL's fully-qualified canonical form. The
// trigger name and its ON table are both qualified by the trigger's OWN database (a trigger always
// lives in the same database as its table). Timing/Event are emitted in their stored upper-case
// form; the body is appended verbatim.
//
// Identifiers use the ORIGINAL-CASE names (trig.Database.Name, trig.Name, trig.Table), not the
// lower-cased diff key `database` — on a case-sensitive server (lower_case_table_names=0, the Linux
// default) the stored object name is case-significant, so the DDL must reproduce the declared
// casing, the same rationale as tableIdent in migration_table.go. `database` (the lower-cased diff
// key) is retained only for the Database/sortName ordering metadata.
func buildCreateTriggerOp(database string, trig *Trigger) MigrationOp {
	dbName := triggerDatabaseName(trig)
	trigIdent := qualifiedTriggerIdent(dbName, trig.Name)
	onTableIdent := qualifiedTriggerIdent(dbName, trig.Table)

	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TRIGGER %s %s %s ON %s FOR EACH ROW",
		trigIdent, strings.ToUpper(trig.Timing), strings.ToUpper(trig.Event), onTableIdent)
	if body := strings.TrimSpace(trig.Body); body != "" {
		b.WriteString(" ")
		b.WriteString(body)
	}

	return MigrationOp{
		Type:         OpCreateTrigger,
		Database:     database,
		ObjectName:   trig.Name,
		ParentObject: trig.Table,
		SQL:          b.String(),
		Phase:        PhaseMain,
		Priority:     priorityTrigger,
		sortName:     triggerSortName(database, trig.Name),
	}
}

// buildDropTriggerOp renders a DROP TRIGGER op (PhasePre, so it precedes a same-name re-create
// and the PhaseMain CREATE half of a modify). The trigger name is qualified by its own database
// (original case, like buildCreateTriggerOp).
func buildDropTriggerOp(database string, trig *Trigger) MigrationOp {
	return MigrationOp{
		Type:         OpDropTrigger,
		Database:     database,
		ObjectName:   trig.Name,
		ParentObject: trig.Table,
		SQL:          fmt.Sprintf("DROP TRIGGER %s", qualifiedTriggerIdent(triggerDatabaseName(trig), trig.Name)),
		Phase:        PhasePre,
		Priority:     priorityTrigger,
		sortName:     triggerSortName(database, trig.Name),
	}
}

// triggerDatabaseName returns the trigger's owning database name in its original case, or "" when
// the trigger was loaded without a qualifying database (the synced single-database release path).
func triggerDatabaseName(trig *Trigger) string {
	if trig.Database == nil {
		return ""
	}
	return trig.Database.Name
}

// qualifiedTriggerIdent backtick-quotes `db`.`name`, omitting the database qualifier when the
// database scope is empty (the synced single-database release path may load with no qualifying
// database, exactly like tableIdent in migration_table.go).
func qualifiedTriggerIdent(database, name string) string {
	if database == "" {
		return migrationQuoteIdent(name)
	}
	return migrationQuoteIdent(database) + "." + migrationQuoteIdent(name)
}

// triggerOpRank orders trigger op types so drops sort ahead of creates locally (phase ordering
// also enforces this globally).
func triggerOpRank(t MigrationOpType) int {
	if t == OpDropTrigger {
		return 0
	}
	return 1
}

// triggerSortName is the stable secondary sort key for a trigger op: lower-cased database.name.
func triggerSortName(database, name string) string {
	return strings.ToLower(database) + "." + strings.ToLower(name)
}
