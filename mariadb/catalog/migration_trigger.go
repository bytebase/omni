package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateTriggerDDL turns the trigger diff (SchemaDiff.Triggers) into CREATE TRIGGER and DROP
// TRIGGER ops. There is no ALTER TRIGGER, so a modified trigger is a DROP followed by a CREATE.
//
// CREATE runs in PhaseMain at priorityTrigger, after the table CREATE and column ALTERs, so the
// table the trigger binds to exists. DROP runs in PhasePre, so it precedes the re-create.
//
// Dropping a table drops its triggers, so an explicit DROP TRIGGER for a trigger whose owning
// table is dropped in the same plan is skipped (dropTriggerSuppressed).
//
// Identifiers are fully qualified and backtick-quoted, independent of the current database, and
// the body is emitted verbatim.
func generateTriggerDDL(_, to *Catalog, diff *SchemaDiff, _ *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Triggers {
		te := &diff.Triggers[i]
		switch te.Action {
		case DiffAdd:
			if te.To == nil {
				continue
			}
			// No original session context to preserve.
			ops = append(ops, buildCreateTriggerOp(te.Database, te.To, nil))
		case DiffDrop:
			if te.From == nil || dropTriggerSuppressed(to, te.From) {
				continue
			}
			ops = append(ops, buildDropTriggerOp(te.Database, te.From))
		case DiffModify:
			// The DROP half is suppressed like a plain DiffDrop: a trigger moved to another table
			// while its old table is dropped is already gone with that table.
			if te.From != nil && !dropTriggerSuppressed(to, te.From) {
				ops = append(ops, buildDropTriggerOp(te.Database, te.From))
			}
			if te.To != nil {
				// The recreate runs under the old trigger's session context.
				ops = append(ops, buildCreateTriggerOp(te.Database, te.To, te.From))
			}
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ri, rj := triggerOpRank(ops[i].Type), triggerOpRank(ops[j].Type); ri != rj {
			return ri < rj
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// dropTriggerSuppressed reports whether an explicit DROP TRIGGER for `trig` must be skipped
// because its owning table is absent in `to`, so DROP TABLE already removes the trigger.
func dropTriggerSuppressed(to *Catalog, trig *Trigger) bool {
	if to == nil || trig == nil || trig.Database == nil || trig.Table == "" {
		return false
	}
	db := to.GetDatabase(trig.Database.Name)
	if db == nil {
		return true // the whole database is gone
	}
	return db.GetTable(trig.Table) == nil
}

// buildCreateTriggerOp renders a CREATE TRIGGER op with the trigger name and its ON table
// qualified by the trigger's own database, in original case, and the body verbatim. The CREATE
// is wrapped in a save/restore of the session context it is created under
// (renderWithSessionContext): trig's own where it is known, else that of `from`, the trigger it
// recreates.
func buildCreateTriggerOp(database string, trig, from *Trigger) MigrationOp {
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

	sql := b.String()
	switch {
	case trig.HasSessionContext:
		sql = renderWithSessionContext(sql, triggerSessionContext(trig), false)
	case from != nil && from.HasSessionContext:
		sql = renderWithSessionContext(sql, triggerSessionContext(from), false)
	}

	return MigrationOp{
		Type:         OpCreateTrigger,
		Database:     database,
		ObjectName:   trig.Name,
		ParentObject: trig.Table,
		SQL:          sql,
		Phase:        PhaseMain,
		Priority:     priorityTrigger,
		sortName:     triggerSortName(database, trig.Name),
	}
}

// triggerSessionContext bundles a trigger's captured session context for the recreate framing.
// Triggers have no time_zone.
func triggerSessionContext(t *Trigger) SessionContext {
	return SessionContext{
		SQLMode:             t.SQLMode,
		CharacterSetClient:  t.CharacterSetClient,
		CollationConnection: t.CollationConnection,
	}
}

// buildDropTriggerOp renders a DROP TRIGGER op in PhasePre, qualified by the trigger's own
// database.
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
// the trigger has no database.
func triggerDatabaseName(trig *Trigger) string {
	if trig.Database == nil {
		return ""
	}
	return trig.Database.Name
}

// qualifiedTriggerIdent backtick-quotes `db`.`name`, omitting the qualifier when database is
// empty.
func qualifiedTriggerIdent(database, name string) string {
	if database == "" {
		return migrationQuoteIdent(name)
	}
	return migrationQuoteIdent(database) + "." + migrationQuoteIdent(name)
}

// triggerOpRank orders trigger op types so drops sort ahead of creates.
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
