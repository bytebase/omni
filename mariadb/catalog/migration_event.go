package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateEventDDL turns SchemaDiff.Events into CREATE EVENT, ALTER EVENT, and DROP EVENT ops.
// A modified event is an ALTER EVENT, rendered as an OpCreateEvent-typed op because no
// OpAlterEvent exists, unless eventModifyNeedsRecreate calls for DROP + CREATE.
//
// DEFINER is never emitted; the diff ignores it. The schedule is emitted verbatim from the
// target.
func generateEventDDL(_, _ *Catalog, diff *SchemaDiff, _ *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Events {
		e := &diff.Events[i]
		switch e.Action {
		case DiffAdd:
			// No original session context to preserve.
			ops = append(ops, MigrationOp{
				Type:       OpCreateEvent,
				Database:   e.Database,
				ObjectName: e.Name,
				SQL:        renderCreateEventWithContext(e.To, nil),
				Phase:      PhaseMain,
				Priority:   priorityEvent,
				sortName:   eventSortName(e.Database, e.Name),
			})
		case DiffModify:
			if eventModifyNeedsRecreate(e.From, e.To) {
				// The recreate runs under the old event's session context, including time_zone.
				ops = append(ops,
					MigrationOp{
						Type:       OpDropEvent,
						Database:   e.Database,
						ObjectName: e.Name,
						SQL:        renderDropEvent(e.From),
						Phase:      PhasePre,
						Priority:   priorityEvent,
						sortName:   eventSortName(e.Database, e.Name),
					},
					MigrationOp{
						Type:       OpCreateEvent,
						Database:   e.Database,
						ObjectName: e.Name,
						SQL:        renderCreateEventWithContext(e.To, e.From),
						Phase:      PhaseMain,
						Priority:   priorityEvent,
						sortName:   eventSortName(e.Database, e.Name),
					},
				)
				continue
			}
			// The ALTER also runs under the old event's session context, so the stored modes do
			// not come from the deploy session.
			ops = append(ops, MigrationOp{
				Type:       OpCreateEvent, // SQL is ALTER EVENT
				Database:   e.Database,
				ObjectName: e.Name,
				SQL:        renderAlterEventWithContext(e.To, e.From),
				Phase:      PhaseMain,
				Priority:   priorityEvent,
				sortName:   eventSortName(e.Database, e.Name),
			})
		case DiffDrop:
			ops = append(ops, MigrationOp{
				Type:       OpDropEvent,
				Database:   e.Database,
				ObjectName: e.Name,
				SQL:        renderDropEvent(e.From),
				Phase:      PhasePre,
				Priority:   priorityEvent,
				sortName:   eventSortName(e.Database, e.Name),
			})
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// eventSortName is the stable secondary sort key for an event op: database.name lower-cased.
func eventSortName(database, name string) string {
	return toLower(database) + "." + toLower(name)
}

// eventModifyNeedsRecreate reports whether a same-name event MODIFY is applied as DROP + CREATE
// rather than ALTER EVENT. That is the case when the target drops the ENDS bound, which ALTER
// EVENT ... ON SCHEDULE is not relied on to clear.
func eventModifyNeedsRecreate(from, to *Event) bool {
	if from == nil || to == nil {
		return false
	}
	return scheduleHasEnds(from.Schedule) && !scheduleHasEnds(to.Schedule)
}

// scheduleHasEnds reports whether a raw schedule string carries a top-level ENDS clause, using
// the same quote-aware scan as the canonical key.
func scheduleHasEnds(raw string) bool {
	_, ends := peelScheduleClause(collapseSpaces(raw), "ENDS")
	return ends != ""
}

// renderCreateEvent renders a CREATE EVENT statement in SHOW CREATE EVENT clause order, without
// DEFINER. Defaults (ON COMPLETION NOT PRESERVE, ENABLE) are emitted explicitly.
func renderCreateEvent(e *Event) string {
	var b strings.Builder
	b.WriteString("CREATE EVENT ")
	b.WriteString(eventIdent(e))
	writeEventBodyClauses(&b, e)
	return b.String()
}

// renderCreateEventWithContext renders the target event's CREATE, framed by frameEventWithContext.
func renderCreateEventWithContext(to, from *Event) string {
	return frameEventWithContext(renderCreateEvent(to), to, from)
}

// renderAlterEventWithContext renders the target event's ALTER, framed by frameEventWithContext.
func renderAlterEventWithContext(to, from *Event) string {
	return frameEventWithContext(renderAlterEvent(to), to, from)
}

// frameEventWithContext wraps an event statement in a save/restore of the session context it
// runs under: the target event's own where it is known, else the old event's.
func frameEventWithContext(stmtSQL string, to, from *Event) string {
	switch {
	case to != nil && to.HasSessionContext:
		return renderWithSessionContext(stmtSQL, eventSessionContext(to), true)
	case from != nil && from.HasSessionContext:
		return renderWithSessionContext(stmtSQL, eventSessionContext(from), true)
	}
	return stmtSQL
}

// eventSessionContext bundles an event's captured session context, including time_zone.
func eventSessionContext(e *Event) SessionContext {
	return SessionContext{
		SQLMode:             e.SQLMode,
		CharacterSetClient:  e.CharacterSetClient,
		CollationConnection: e.CollationConnection,
		TimeZone:            e.TimeZone,
	}
}

// renderAlterEvent renders an ALTER EVENT that re-states every mutable aspect of the target
// event, so the event matches the target regardless of its prior state.
func renderAlterEvent(e *Event) string {
	var b strings.Builder
	b.WriteString("ALTER EVENT ")
	b.WriteString(eventIdent(e))
	writeEventBodyClauses(&b, e)
	return b.String()
}

// renderDropEvent renders DROP EVENT IF EXISTS for the source event.
func renderDropEvent(e *Event) string {
	return "DROP EVENT IF EXISTS " + eventIdent(e)
}

// eventIdent names an event with its database, so the plan does not depend on
// the database the connection has selected.
func eventIdent(e *Event) string {
	if e.Database != nil && e.Database.Name != "" {
		return quoteIdent(e.Database.Name) + "." + quoteIdent(e.Name)
	}
	return quoteIdent(e.Name)
}

// writeEventBodyClauses appends the CREATE/ALTER EVENT clauses after the event name: ON SCHEDULE,
// ON COMPLETION, status, COMMENT, DO body. Every clause but the schedule is emitted
// unconditionally, including defaults and an empty COMMENT, because ALTER EVENT keeps a clause
// it does not name.
func writeEventBodyClauses(b *strings.Builder, e *Event) {
	if sched := strings.TrimSpace(e.Schedule); sched != "" {
		b.WriteString(" ON SCHEDULE ")
		b.WriteString(sched)
	}

	if eqFoldStr(strings.TrimSpace(e.OnCompletion), "PRESERVE") {
		b.WriteString(" ON COMPLETION PRESERVE")
	} else {
		b.WriteString(" ON COMPLETION NOT PRESERVE")
	}

	switch canonicalEventStatus(e.Enable) {
	case "DISABLE":
		b.WriteString(" DISABLE")
	case "DISABLE ON SLAVE":
		b.WriteString(" DISABLE ON SLAVE")
	default:
		b.WriteString(" ENABLE")
	}

	fmt.Fprintf(b, " COMMENT '%s'", escapeComment(e.Comment))

	if body := strings.TrimSpace(e.Body); body != "" {
		b.WriteString(" DO ")
		b.WriteString(body)
	}
}
