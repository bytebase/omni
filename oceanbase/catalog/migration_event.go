package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — scheduled events (CREATE / ALTER / DROP EVENT).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so event diffs (SchemaDiff.Events)
// become DDL. Mirrors generateTableDDL's shape. The three cases, from EventDiffEntry.Action:
//
//   - DiffAdd    → CREATE EVENT … (PhaseMain, priorityEvent)
//   - DiffModify → ALTER EVENT  … (PhaseMain, priorityEvent) — preferred over DROP+CREATE
//   - DiffDrop   → DROP EVENT   … (PhasePre,  priorityEvent) — drops run before creates
//
// ALTER vs DROP+CREATE: MySQL's ALTER EVENT can change every mutable aspect an SDL diff can
// observe — the schedule (ON SCHEDULE …), ON COMPLETION, status (ENABLE/DISABLE/DISABLE ON SLAVE),
// COMMENT, and body (DO …). The event name is its identity (a rename is a drop+add of distinct
// keys, handled by the Add/Drop arms), so a MODIFY never needs to change the name. ALTER is
// therefore always applicable for a same-name change and is preferred: it preserves the event
// object (and its bookkeeping) instead of dropping and recreating it. Both forms are apply-correct
// — verified on live 5.7.25 + 8.0.32 — but ALTER is the smaller, safer change.
//
// Only OpCreateEvent / OpDropEvent op-types exist (declared in migration.go); a MODIFY is rendered
// as an OpCreateEvent-typed op carrying ALTER EVENT SQL (the op-type is a coarse drop-vs-forward
// classifier for ordering/advice; the SQL string is what executes). This keeps the op-type set as
// declared without an OpAlterEvent.
//
// DEFINER is never emitted: it is ignore-in-diff (always resolves to the executing user) and
// omitting it lets the apply readback match the synced form (“ DEFINER=`root`@`%` “). The
// schedule is emitted verbatim from the target — applying `EVERY 1 HOUR` lets the engine inject
// its own STARTS, whose readback canonicalizes equal to the target (canonicalizeSchedule strips
// STARTS), so apply-correctness holds.
//
// implemented by omni:events breadth node
func generateEventDDL(_, _ *Catalog, diff *SchemaDiff, _ *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Events {
		e := &diff.Events[i]
		switch e.Action {
		case DiffAdd:
			// A brand-new event emits a BARE CREATE (→ server-default session): no original
			// context to preserve.
			ops = append(ops, MigrationOp{
				Type:       OpCreateEvent,
				Database:   e.Database,
				ObjectName: e.Name,
				SQL:        renderCreateEvent(e.To),
				Phase:      PhaseMain,
				Priority:   priorityEvent,
				sortName:   eventSortName(e.Database, e.Name),
			})
		case DiffModify:
			if eventModifyNeedsRecreate(e.From, e.To) {
				// ALTER EVENT cannot express this change on every supported version (the one known
				// case: 5.7's ALTER EVENT … ON SCHEDULE does NOT clear an existing ENDS — verified
				// on live 5.7.25, where the old ENDS survives; 8.0.32 clears it). Fall back to
				// DROP + CREATE, which is always correct: a PhasePre DROP removes the old event
				// before the PhaseMain CREATE re-establishes the target with no ENDS. The recreate
				// runs under the OLD event's session context (e.From), including its time_zone.
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
			// ALTER EVENT … DO <body> RE-STAMPS the event's sql_mode / character_set_client /
			// collation_connection from the executing session (verified on live 8.0.32 — a bare
			// ALTER EVENT DO under a different session flips ROUTINES/EVENTS.SQL_MODE), so the
			// ALTER must run under the OLD event's context too, not just the DROP+CREATE path.
			ops = append(ops, MigrationOp{
				Type:       OpCreateEvent, // forward (PhaseMain) op; SQL is ALTER EVENT (see file doc)
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

// eventModifyNeedsRecreate reports whether a same-name event MODIFY must be applied as DROP+CREATE
// rather than ALTER EVENT, because ALTER cannot express the change on every supported version.
//
// The one known case is ENDS removal: 5.7's ALTER EVENT … ON SCHEDULE does not clear a previously
// set ENDS (the old ENDS survives the re-specification — verified on live 5.7.25), so an event that
// drops its ENDS bound can only reach the target via DROP+CREATE there. 8.0 clears ENDS on ALTER,
// but DROP+CREATE is correct on 8.0 too, so the fallback is applied version-independently for a
// single, simple code path. Every other mutable aspect (schedule interval/unit, EVERY↔AT, status,
// ON COMPLETION, comment, body) is reachable via ALTER on both versions (verified), so this is the
// only trigger.
func eventModifyNeedsRecreate(from, to *Event) bool {
	if from == nil || to == nil {
		return false
	}
	return scheduleHasEnds(from.Schedule) && !scheduleHasEnds(to.Schedule)
}

// scheduleHasEnds reports whether a raw schedule string carries a top-level ENDS clause. It reuses
// the same quote-aware clause scan the canonical key uses, so a timestamp literal that happens to
// contain the word ENDS is not mistaken for a clause.
func scheduleHasEnds(raw string) bool {
	_, ends := peelScheduleClause(collapseSpaces(raw), "ENDS")
	return ends != ""
}

// renderCreateEvent renders a CREATE EVENT statement for the target event. The clauses follow
// MySQL's SHOW CREATE EVENT order (sans DEFINER): name, ON SCHEDULE, ON COMPLETION, status,
// COMMENT, DO body. Defaults (ON COMPLETION NOT PRESERVE, ENABLE) are emitted explicitly so the
// statement is self-describing and matches the readback form.
func renderCreateEvent(e *Event) string {
	var b strings.Builder
	b.WriteString("CREATE EVENT ")
	b.WriteString(quoteIdent(e.Name))
	writeEventBodyClauses(&b, e)
	return b.String()
}

// renderCreateEventWithContext renders the target event's CREATE and, when the OLD event
// (`from`) carries synced session context, wraps it in a save/restore of that context —
// sql_mode, character_set_client, collation_connection, AND the session time_zone — so an
// event re-created via DROP+CREATE keeps the schedule/behavior it had under its original
// session (a non-UTC time_zone in particular changes how a literal schedule is
// interpreted). Without synced context (`from` nil / unset) it is a bare CREATE.
func renderCreateEventWithContext(to, from *Event) string {
	return frameEventWithContext(renderCreateEvent(to), from)
}

// renderAlterEventWithContext renders the target event's ALTER and, when the OLD event
// (`from`) carries synced session context, wraps it in the same save/restore. This is
// REQUIRED, not cosmetic: ALTER EVENT … DO <body> re-stamps sql_mode/charset/collation from
// the executing session (verified on live 8.0.32), so a bare ALTER under the deploy session
// would rewrite the event's stored modes. time_zone is not re-stamped by ALTER, but framing
// it is harmless and keeps the CREATE and ALTER paths uniform.
func renderAlterEventWithContext(to, from *Event) string {
	return frameEventWithContext(renderAlterEvent(to), from)
}

// frameEventWithContext wraps an event CREATE/ALTER statement in a save/restore of the OLD
// event's session context when that context is present; otherwise returns the bare statement.
func frameEventWithContext(stmtSQL string, from *Event) string {
	if from != nil && from.HasSessionContext {
		return renderWithSessionContext(stmtSQL, eventSessionContext(from), true)
	}
	return stmtSQL
}

// eventSessionContext bundles an event's captured session context for the framing, including
// the event-only time_zone axis.
func eventSessionContext(e *Event) SessionContext {
	return SessionContext{
		SQLMode:             e.SQLMode,
		CharacterSetClient:  e.CharacterSetClient,
		CollationConnection: e.CollationConnection,
		TimeZone:            e.TimeZone,
	}
}

// renderAlterEvent renders an ALTER EVENT that re-establishes every mutable aspect of the target
// event in one statement. ALTER EVENT accepts the same clause set as CREATE (minus IF NOT EXISTS),
// and re-specifying all of them makes the event match the target regardless of its prior state.
func renderAlterEvent(e *Event) string {
	var b strings.Builder
	b.WriteString("ALTER EVENT ")
	b.WriteString(quoteIdent(e.Name))
	writeEventBodyClauses(&b, e)
	return b.String()
}

// renderDropEvent renders DROP EVENT IF EXISTS for the source event. IF EXISTS makes the drop
// tolerant of an event already gone (idempotent apply).
func renderDropEvent(e *Event) string {
	return "DROP EVENT IF EXISTS " + quoteIdent(e.Name)
}

// writeEventBodyClauses appends the shared CREATE/ALTER EVENT clause sequence (everything after
// the event name): ON SCHEDULE, ON COMPLETION, status, COMMENT, DO body. Every clause is emitted
// UNCONDITIONALLY (including the explicit defaults NOT PRESERVE / ENABLE and an explicit
// COMMENT ”) so that for the ALTER path a clause omitted from the target actually RESETS the
// event to that target value. MySQL's ALTER EVENT leaves any omitted clause unchanged, so to
// CLEAR a comment the statement must say COMMENT ” — verified on live 5.7.25 + 8.0.32: an ALTER
// without a COMMENT clause keeps the old comment, while COMMENT ” removes it (and a CREATE/ALTER
// with COMMENT ” reads back with no COMMENT clause, canonicalizing equal to an empty comment). The
// same reasoning is why status and ON COMPLETION are always rendered. The schedule is the one
// clause that cannot be empty for a valid event, so it is emitted whenever present.
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

	// Always emit COMMENT (COMMENT '' when empty) so an ALTER that drops a comment actually clears
	// it; an empty comment reads back as no clause on both versions, so this stays apply-correct.
	fmt.Fprintf(b, " COMMENT '%s'", escapeComment(e.Comment))

	if body := strings.TrimSpace(e.Body); body != "" {
		b.WriteString(" DO ")
		b.WriteString(body)
	}
}

// quoteIdent wraps an identifier in backticks, escaping any embedded backtick by doubling it
// (MySQL's identifier-quoting rule). Event names are simple identifiers in practice, but this
// keeps the rendered DDL correct for any legal name.
func quoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
