package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Trigger represents a trigger on a relation.
type Trigger struct {
	OID               uint32
	Name              string
	RelOID            uint32
	FuncOID           uint32
	Timing            TriggerTiming
	Events            TriggerEvent
	ForEachRow        bool
	WhenExpr          string
	Columns           []int16  // UPDATE OF columns; nil = all
	Enabled           byte     // 'O'=origin(enabled), 'D'=disabled, 'A'=always, 'R'=replica
	OldTransitionName string   // OLD TABLE AS name (transition tables)
	NewTransitionName string   // NEW TABLE AS name (transition tables)
	IsConstraint      bool     // constraint trigger (pg: tgconstraint != 0)
	Deferrable        bool     // constraint trigger is DEFERRABLE
	Initdeferred      bool     // constraint trigger is INITIALLY DEFERRED
	ConstraintRelOID  uint32   // referenced relation for FK constraint triggers
	Args              []string // trigger function arguments
}

// CreateTriggerStmt creates a new trigger on a relation from a pgparser AST.
//
// pg: src/backend/commands/trigger.c — CreateTrigger
func (c *Catalog) CreateTriggerStmt(stmt *nodes.CreateTrigStmt) error {
	// Map timing bitmask.
	var timing TriggerTiming
	switch {
	case stmt.Timing&nodes.TRIGGER_TYPE_BEFORE != 0:
		timing = TriggerBefore
	case stmt.Timing&nodes.TRIGGER_TYPE_INSTEAD != 0:
		timing = TriggerInsteadOf
	default:
		timing = TriggerAfter
	}

	// Map events bitmask.
	var events TriggerEvent
	if stmt.Events&nodes.TRIGGER_TYPE_INSERT != 0 {
		events |= TriggerEventInsert
	}
	if stmt.Events&nodes.TRIGGER_TYPE_UPDATE != 0 {
		events |= TriggerEventUpdate
	}
	if stmt.Events&nodes.TRIGGER_TYPE_DELETE != 0 {
		events |= TriggerEventDelete
	}
	if stmt.Events&nodes.TRIGGER_TYPE_TRUNCATE != 0 {
		events |= TriggerEventTruncate
	}

	// Function name.
	funcSchema, funcName := qualifiedName(stmt.Funcname)

	// Relation.
	var relSchema, relName string
	if stmt.Relation != nil {
		relSchema = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	trigName := stmt.Trigname
	forEachRow := stmt.Row
	columnNames := stringListItems(stmt.Columns)
	orReplace := stmt.Replace
	isConstraint := stmt.IsConstraint

	// Parse trigger arguments.
	// pg: src/backend/commands/trigger.c — CreateTrigger (trigger args)
	var trigArgs []string
	if stmt.Args != nil {
		for _, item := range stmt.Args.Items {
			trigArgs = append(trigArgs, stringVal(item))
		}
	}

	// Parse constraint trigger referenced relation.
	// pg: src/backend/commands/trigger.c — CreateTrigger (constrrelid)
	var constraintRelOID uint32
	if isConstraint && stmt.Constrrel != nil {
		_, constrRel, constrErr := c.findRelation(stmt.Constrrel.Schemaname, stmt.Constrrel.Relname)
		if constrErr != nil {
			return constrErr
		}
		constraintRelOID = constrRel.OID
	}

	// Find relation.
	_, rel, err := c.findRelation(relSchema, relName)
	if err != nil {
		return err
	}

	// Relkind-specific validations.
	// pg: src/backend/commands/trigger.c — CreateTrigger (relkind switch)
	switch rel.RelKind {
	case 'r', 'm': // tables, matviews
		if timing == TriggerInsteadOf {
			return errInvalidObjectDefinition("INSTEAD OF triggers can only be defined on views")
		}
	case 'p': // partitioned tables
		if timing == TriggerInsteadOf {
			return errInvalidObjectDefinition("INSTEAD OF triggers can only be defined on views")
		}
		// Partitioned tables: ROW triggers with transition tables not supported.
		// pg: src/backend/commands/trigger.c — CreateTrigger (line 257-262)
		if forEachRow && stmt.TransitionRels != nil && len(stmt.TransitionRels.Items) > 0 {
			return &Error{Code: CodeFeatureNotSupported,
				Message: fmt.Sprintf("\"%s\" is a partitioned table\nDetail: ROW triggers with transition tables are not supported on partitioned tables.", rel.Name)}
		}
	case 'v': // views
		// Views can have INSTEAD OF (ROW) triggers and STATEMENT triggers.
		// ROW triggers with BEFORE/AFTER are not allowed on views.
		if timing != TriggerInsteadOf && forEachRow {
			return errInvalidObjectDefinition("ROW triggers with data manipulation statements are not supported for views")
		}
	case 'f': // foreign tables
		if timing == TriggerInsteadOf {
			return errInvalidObjectDefinition("INSTEAD OF triggers can only be defined on views")
		}
	}

	// Disallow ROW-level TRUNCATE triggers.
	// pg: src/backend/commands/trigger.c — CreateTrigger (line 377)
	if forEachRow && events&TriggerEventTruncate != 0 {
		return errInvalidObjectDefinition("TRUNCATE FOR EACH ROW triggers are not supported")
	}

	// TRUNCATE triggers on views and foreign tables are not supported.
	// pg: src/backend/commands/trigger.c — CreateTrigger (line 359-362)
	if events&TriggerEventTruncate != 0 && (rel.RelKind == 'v' || rel.RelKind == 'f') {
		return errInvalidObjectDefinition(fmt.Sprintf(
			"\"%s\" is not a table or materialized view", rel.Name))
	}

	// INSTEAD OF triggers must be row-level, can't have WHEN or columns.
	// pg: src/backend/commands/trigger.c — CreateTrigger (line 383-397)
	if timing == TriggerInsteadOf {
		if !forEachRow {
			return errInvalidObjectDefinition("INSTEAD OF triggers must be FOR EACH ROW")
		}
		if stmt.WhenClause != nil {
			return errInvalidObjectDefinition("INSTEAD OF triggers cannot have WHEN conditions")
		}
		if len(columnNames) > 0 {
			return errInvalidObjectDefinition("INSTEAD OF triggers cannot have column lists")
		}
	}

	// Analyze WHEN clause expression if present.
	// pg: src/backend/commands/trigger.c — CreateTrigger (WHEN clause setup)
	// PG creates two RTEs for OLD (varno=1) and NEW (varno=2).
	var whenExpr string
	var whenAnalyzed AnalyzedExpr
	if stmt.WhenClause != nil {
		analyzed, err := c.AnalyzeTriggerWhenExpr(stmt.WhenClause, rel)
		if err == nil {
			whenAnalyzed = analyzed
			oldRTE := c.buildRelationRTE(rel)
			oldRTE.ERef = "old"
			newRTE := c.buildRelationRTE(rel)
			newRTE.ERef = "new"
			whenExpr = c.DeparseExpr(analyzed, []*RangeTableEntry{oldRTE, newRTE}, true)
		} else {
			return fmt.Errorf("failed to analyze trigger WHEN clause: %w", err)
		}
	}

	// Transition table validations.
	// pg: src/backend/commands/trigger.c — CreateTrigger (line 414-506)
	if stmt.TransitionRels != nil && len(stmt.TransitionRels.Items) > 0 {
		if rel.RelKind == 'f' {
			return errInvalidObjectDefinition("triggers on foreign tables cannot have transition tables")
		}
		if rel.RelKind == 'v' {
			return errInvalidObjectDefinition("triggers on views cannot have transition tables")
		}
		// Inheritance children (non-partition) with ROW triggers + transition tables.
		// pg: src/backend/commands/trigger.c — CreateTrigger (line ~270)
		if rel.PartitionOf == 0 && rel.InhCount > 0 && forEachRow {
			return &Error{Code: CodeFeatureNotSupported,
				Message: "ROW triggers with transition tables are not supported on tables that have inheritance children"}
		}
		if timing != TriggerAfter {
			return errInvalidObjectDefinition("transition table name can only be specified for an AFTER trigger")
		}
		if events&TriggerEventTruncate != 0 {
			return errInvalidObjectDefinition("TRUNCATE triggers with transition tables are not supported")
		}
		// Count events — transition tables require single event.
		eventCount := 0
		if events&TriggerEventInsert != 0 {
			eventCount++
		}
		if events&TriggerEventUpdate != 0 {
			eventCount++
		}
		if events&TriggerEventDelete != 0 {
			eventCount++
		}
		if eventCount > 1 {
			return errInvalidObjectDefinition("transition tables cannot be specified for triggers with more than one event")
		}
		if len(columnNames) > 0 {
			return errInvalidObjectDefinition("transition tables cannot be specified for triggers with column lists")
		}

		// Validate OLD TABLE only for DELETE/UPDATE, NEW TABLE only for INSERT/UPDATE.
		// pg: src/backend/commands/trigger.c — CreateTrigger (line 517-548)
		for _, item := range stmt.TransitionRels.Items {
			tr, ok := item.(*nodes.TriggerTransition)
			if !ok {
				continue
			}
			if !tr.IsNew {
				// OLD TABLE: only for DELETE or UPDATE triggers
				if events&TriggerEventDelete == 0 && events&TriggerEventUpdate == 0 {
					return errInvalidObjectDefinition("OLD TABLE can only be specified for a DELETE or UPDATE trigger")
				}
			} else {
				// NEW TABLE: only for INSERT or UPDATE triggers
				if events&TriggerEventInsert == 0 && events&TriggerEventUpdate == 0 {
					return errInvalidObjectDefinition("NEW TABLE can only be specified for an INSERT or UPDATE trigger")
				}
			}
		}
	}

	// Find trigger function.
	funcBP, wrongRet := c.findTriggerFunc(funcSchema, funcName)
	if funcBP == nil {
		if wrongRet {
			return errInvalidObjectDefinition(fmt.Sprintf(
				"function %q must return type trigger", funcName))
		}
		return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("function %q does not exist", funcName)}
	}

	// Check for duplicate trigger name on this relation.
	// pg: src/backend/commands/trigger.c — CreateTrigger (trigger replacement)
	for _, trig := range c.triggersByRel[rel.OID] {
		if trig.Name == trigName {
			if orReplace {
				// Cannot replace constraint triggers.
				// pg: src/backend/commands/trigger.c — CreateTrigger
				if trig.IsConstraint {
					return errInvalidObjectDefinition(fmt.Sprintf(
						"cannot replace trigger %q because it is a constraint trigger", trigName))
				}
				// Replace: remove old trigger.
				c.removeTrigger(trig)
				break
			}
			return errDuplicateObject("trigger", trigName)
		}
	}

	// Resolve UPDATE OF columns.
	var cols []int16
	if len(columnNames) > 0 {
		cols, err = c.resolveColumnNames(rel, columnNames)
		if err != nil {
			return err
		}
	}

	// Extract transition table names.
	// pg: src/backend/commands/trigger.c — CreateTrigger (transition tables)
	var oldTransName, newTransName string
	if stmt.TransitionRels != nil {
		for _, item := range stmt.TransitionRels.Items {
			if tr, ok := item.(*nodes.TriggerTransition); ok {
				if tr.IsNew {
					newTransName = tr.Name
				} else {
					oldTransName = tr.Name
				}
			}
		}
	}

	// Check that OLD and NEW transition table names are different.
	// pg: src/backend/commands/trigger.c — CreateTrigger (line 549-553)
	if oldTransName != "" && newTransName != "" && oldTransName == newTransName {
		return errInvalidObjectDefinition(fmt.Sprintf("OLD TABLE name and NEW TABLE name cannot be the same"))
	}

	trigOID := c.oidGen.Next()
	trig := &Trigger{
		OID:               trigOID,
		Name:              trigName,
		RelOID:            rel.OID,
		FuncOID:           funcBP.OID,
		Timing:            timing,
		Events:            events,
		ForEachRow:        forEachRow,
		WhenExpr:          whenExpr,
		Columns:           cols,
		Enabled:           'O',
		OldTransitionName: oldTransName,
		NewTransitionName: newTransName,
		IsConstraint:      isConstraint,
		Deferrable:        stmt.Deferrable,
		Initdeferred:      stmt.Initdeferred,
		ConstraintRelOID:  constraintRelOID,
		Args:              trigArgs,
	}

	c.triggers[trigOID] = trig
	c.triggersByRel[rel.OID] = append(c.triggersByRel[rel.OID], trig)

	// pg: src/backend/commands/trigger.c:1074 — trigger auto-drops with table
	c.recordDependency('g', trigOID, 0, 'r', rel.OID, 0, DepAuto)
	// pg: src/backend/commands/trigger.c:1049 — trigger depends on function
	c.recordDependency('g', trigOID, 0, 'f', funcBP.OID, 0, DepNormal)

	// pg: src/backend/commands/trigger.c:1110-1121 — per-column deps for UPDATE OF
	if len(cols) > 0 {
		for _, attnum := range cols {
			c.recordDependency('g', trigOID, 0, 'r', rel.OID, int32(attnum), DepNormal)
		}
	}

	// pg: src/backend/commands/trigger.c:1128-1130 — WHEN clause expression deps
	if whenAnalyzed != nil {
		c.recordDependencyOnSingleRelExpr('g', trigOID, whenAnalyzed, rel.OID,
			DepNormal, DepNormal)
	}

	// pg: src/backend/commands/trigger.c — CreateTrigger (pg_constraint for constraint triggers)
	if isConstraint {
		schema := rel.Schema
		if schema == nil {
			schema = c.schemaByName["public"]
		}
		conOID := c.oidGen.Next()
		con := &Constraint{
			OID:        conOID,
			Name:       trigName,
			Type:       ConstraintTrigger,
			RelOID:     rel.OID,
			Namespace:  schema.OID,
			Deferrable: stmt.Deferrable,
			Deferred:   stmt.Initdeferred,
			Validated:  true,
			ConIsLocal: true,
		}
		c.registerConstraint(rel.OID, con)
		// pg: trigger.c:893-905 — trigger depends on its constraint (internal)
		c.recordDependency('g', trigOID, 0, 'c', conOID, 0, DepInternal)
		// pg: trigger.c:855-867 — constraint depends on relation (auto)
		c.recordDependency('c', conOID, 0, 'r', rel.OID, 0, DepAuto)
	}

	return nil
}

// removeTrigger removes a trigger from all catalog maps.
func (c *Catalog) removeTrigger(trig *Trigger) {
	delete(c.triggers, trig.OID)
	c.removeComments('g', trig.OID)
	c.removeDepsOf('g', trig.OID)

	// Remove from triggersByRel.
	list := c.triggersByRel[trig.RelOID]
	for i, t := range list {
		if t.OID == trig.OID {
			c.triggersByRel[trig.RelOID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

// findTriggerFunc finds a function suitable for a trigger.
// Returns (proc, false) on success, (nil, true) if found but wrong return type,
// (nil, false) if not found at all.
func (c *Catalog) findTriggerFunc(schemaName, funcName string) (*BuiltinProc, bool) {
	// Triggers call functions with 0 declared args that return TRIGGER.
	foundAny := false
	for _, p := range c.procByName[funcName] {
		if p.NArgs != 0 {
			continue
		}
		if schemaName != "" {
			if schemaName == "pg_catalog" {
				// Built-in functions live in pg_catalog; allow them through.
				// Only skip if this is a user-defined proc in a different schema.
				if up := c.userProcs[p.OID]; up != nil {
					continue // user proc, not pg_catalog
				}
			} else {
				up := c.userProcs[p.OID]
				if up == nil || c.schemaByName[schemaName] == nil || up.Schema.OID != c.schemaByName[schemaName].OID {
					continue
				}
			}
		}
		foundAny = true
		if p.RetType == TRIGGEROID {
			return p, false
		}
	}
	return nil, foundAny
}

// TriggersOf returns all triggers on the given relation.
func (c *Catalog) TriggersOf(relOID uint32) []*Trigger {
	return c.triggersByRel[relOID]
}
