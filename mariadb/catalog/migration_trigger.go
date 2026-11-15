package catalog

import (
	"fmt"
	"slices"
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
			ops = append(ops, buildCreateTriggerOp(te.Database, te.To, nil, triggerPlacement(to, diff, te.To)))
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
				ops = append(ops, buildCreateTriggerOp(te.Database, te.To, te.From, triggerPlacement(to, diff, te.To)))
			}
		}
	}

	orderTriggerCreates(ops, to, diff)
	sort.SliceStable(ops, func(i, j int) bool {
		if ri, rj := triggerOpRank(ops[i].Type), triggerOpRank(ops[j].Type); ri != rj {
			return ri < rj
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// orderTriggerCreates creates triggers in the order the target created them,
// and a trigger after the trigger its FOLLOWS or PRECEDES names, where the plan
// creates that one too: the clause names a trigger that exists.
func orderTriggerCreates(ops []MigrationOp, to *Catalog, diff *SchemaDiff) {
	// A CREATE TRIGGER that names another in FOLLOWS or PRECEDES is created
	// after it: depth is the length of the chain of such references the plan
	// also creates, which orders the creates so each reference exists first.
	creates := make(map[string]int) // lower database.name -> op index
	for i := range ops {
		if ops[i].Type == OpCreateTrigger {
			creates[triggerSortName(ops[i].Database, ops[i].ObjectName)] = i
		}
	}
	var depth func(i int, seen map[int]bool) int
	depth = func(i int, seen map[int]bool) int {
		if seen[i] || ops[i].triggerRef == "" {
			return 0
		}
		seen[i] = true
		ref, ok := creates[triggerSortName(ops[i].Database, ops[i].triggerRef)]
		if !ok {
			return 0
		}
		return depth(ref, seen) + 1
	}
	pos := make(map[string]int)
	if to != nil {
		for _, te := range diff.Triggers {
			if te.To == nil || te.Action == DiffDrop || te.To.Database == nil {
				continue
			}
			if db := to.GetDatabase(te.To.Database.Name); db != nil {
				pos[triggerSortName(te.Database, te.To.Name)] = effectiveTriggerOrder(db, te.To)[toLower(te.To.Name)]
			}
		}
	}
	for i := range ops {
		if ops[i].Type == OpCreateTrigger {
			ops[i].Priority = priorityTrigger + depth(i, map[int]bool{})
			ops[i].sortName = fmt.Sprintf("%s\x00%010d", strings.ToLower(ops[i].Database), pos[ops[i].sortName])
		}
	}
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
func buildCreateTriggerOp(database string, trig, from *Trigger, placement *TriggerOrderInfo) MigrationOp {
	dbName := triggerDatabaseName(trig)
	trigIdent := qualifiedTriggerIdent(dbName, trig.Name)
	onTableIdent := qualifiedTriggerIdent(dbName, trig.Table)

	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TRIGGER %s %s %s ON %s FOR EACH ROW",
		trigIdent, strings.ToUpper(trig.Timing), strings.ToUpper(trig.Event), onTableIdent)
	// The order among a table's triggers of one timing and event is what
	// FOLLOWS or PRECEDES set: they run in it.
	ref := ""
	if o := placement; o != nil && o.TriggerName != "" {
		verb := "PRECEDES"
		if o.Follows {
			verb = "FOLLOWS"
		}
		fmt.Fprintf(&b, " %s %s", verb, migrationQuoteIdent(o.TriggerName))
		ref = o.TriggerName
	}
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
		triggerRef:   ref,
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

// triggerPlacement places a trigger the plan creates, and which names no
// FOLLOWS or PRECEDES of its own, among the triggers of its table, timing and
// event the plan keeps: PRECEDES the first of them the target created after
// it, else FOLLOWS the last created before it. With no kept peer, creation
// order alone places it (orderTriggerCreates).
func triggerPlacement(to *Catalog, diff *SchemaDiff, trig *Trigger) *TriggerOrderInfo {
	if to == nil || trig.Database == nil {
		return nil
	}
	db := to.GetDatabase(trig.Database.Name)
	if db == nil {
		return nil
	}
	created := make(map[string]bool)
	for _, te := range diff.Triggers {
		if te.To != nil && te.Action != DiffDrop && toLower(te.Database) == toLower(trig.Database.Name) {
			created[toLower(te.To.Name)] = true
		}
	}
	// The effective execution order of the group, which FOLLOWS and PRECEDES
	// set, not the raw creation order.
	order := effectiveTriggerOrder(db, trig)
	pos, ok := order[toLower(trig.Name)]
	if !ok {
		return nil
	}
	// The immediate predecessor in effective order, kept or itself created:
	// FOLLOWS it. The creates run in effective order, so it exists first. When
	// trig is first, PRECEDES the next trigger if that one is kept.
	var before, after *Trigger
	for _, peer := range db.Triggers {
		if peer == trig {
			continue
		}
		pp, ok := order[toLower(peer.Name)]
		if !ok {
			continue
		}
		switch {
		case pp < pos && (before == nil || pp > order[toLower(before.Name)]):
			before = peer
		case pp > pos && (after == nil || pp < order[toLower(after.Name)]):
			after = peer
		}
	}
	switch {
	case before != nil:
		return &TriggerOrderInfo{Follows: true, TriggerName: before.Name}
	case after != nil && !created[toLower(after.Name)]:
		return &TriggerOrderInfo{TriggerName: after.Name}
	}
	return nil
}

// effectiveTriggerOrder returns the execution position of each trigger in
// trig's table/timing/event group, starting from creation order and applying
// each trigger's FOLLOWS or PRECEDES, as the server does.
func effectiveTriggerOrder(db *Database, trig *Trigger) map[string]int {
	var group []*Trigger
	for _, t := range db.Triggers {
		if eqFoldStr(t.Table, trig.Table) && eqFoldStr(t.Timing, trig.Timing) && eqFoldStr(t.Event, trig.Event) {
			group = append(group, t)
		}
	}
	slices.SortFunc(group, func(a, b *Trigger) int { return a.created - b.created })
	var seq []*Trigger
	indexOf := func(name string) int {
		for i, t := range seq {
			if eqFoldStr(t.Name, name) {
				return i
			}
		}
		return -1
	}
	for _, t := range group {
		if t.Order == nil || t.Order.TriggerName == "" {
			seq = append(seq, t)
			continue
		}
		ref := indexOf(t.Order.TriggerName)
		if ref < 0 {
			seq = append(seq, t)
			continue
		}
		at := ref + 1
		if t.Order.Follows {
			at = ref + 1
		} else {
			at = ref
		}
		seq = append(seq[:at], append([]*Trigger{t}, seq[at:]...)...)
	}
	out := make(map[string]int, len(seq))
	for i, t := range seq {
		out[toLower(t.Name)] = i
	}
	return out
}
