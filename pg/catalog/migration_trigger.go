package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateTriggerDDL produces CREATE TRIGGER and DROP TRIGGER operations from the diff.
func generateTriggerDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, rel := range diff.Relations {
		switch rel.Action {
		case DiffModify:
			targetRel := rel.From
			if targetRel == nil {
				targetRel = rel.To
			}
			ops = append(ops, triggerOpsForRelation(to, rel.SchemaName, rel.Name, targetRel, rel.Triggers)...)
		case DiffAdd:
			// For newly added relations, triggers are not in the diff entry's
			// Triggers slice (the relation itself is new). Generate triggers
			// from the target catalog.
			if rel.To != nil {
				trigs := to.TriggersOf(rel.To.OID)
				for _, trig := range trigs {
					// A new partition clones its parent's triggers itself.
					if trig.ParentTriggerOID == 0 {
						ops = append(ops, buildCreateAndModeTriggerOps(to, rel.SchemaName, rel.Name, trig)...)
					}
				}
			}
		}
	}

	// Sort for determinism: by schema, table, then trigger name.
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		if ops[i].ParentObject != ops[j].ParentObject {
			return ops[i].ParentObject < ops[j].ParentObject
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// triggerOpsForRelation generates trigger DDL for a modified relation.
func triggerOpsForRelation(to *Catalog, schemaName, tableName string, fromRel *Relation, triggers []TriggerDiffEntry) []MigrationOp {
	var ops []MigrationOp

	for _, te := range triggers {
		switch te.Action {
		case DiffAdd:
			if te.To == nil {
				continue
			}
			ops = append(ops, buildCreateAndModeTriggerOps(to, schemaName, tableName, te.To)...)
		case DiffDrop:
			if te.From == nil {
				continue
			}
			ops = append(ops, buildDropTriggerOpForRel(schemaName, tableName, fromRel, te.From.Name, te.From.OID))
		case DiffModify:
			// Modified trigger = DROP old + CREATE new.
			if te.From != nil {
				ops = append(ops, buildDropTriggerOpForRel(schemaName, tableName, fromRel, te.From.Name, te.From.OID))
			}
			if te.To != nil {
				ops = append(ops, buildCreateAndModeTriggerOps(to, schemaName, tableName, te.To)...)
			}
		}
	}

	return ops
}

// buildDropTriggerOp creates a DROP TRIGGER operation.
func buildDropTriggerOp(schemaName, tableName, trigName string, trigOID uint32) MigrationOp {
	return buildDropTriggerOpForRel(schemaName, tableName, nil, trigName, trigOID)
}

func buildDropTriggerOpForRel(schemaName, tableName string, rel *Relation, trigName string, trigOID uint32) MigrationOp {
	qualifiedTable := formatRelationName(schemaName, tableName, rel)
	return MigrationOp{
		Type:          OpDropTrigger,
		SchemaName:    schemaName,
		ObjectName:    trigName,
		ParentObject:  tableName,
		SQL:           fmt.Sprintf("DROP TRIGGER %s ON %s", quoteIdentAlways(trigName), qualifiedTable),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'T',
		ObjOID:        trigOID,
		Priority:      PriorityTrigger,
	}
}

// buildCreateTriggerOp creates a CREATE TRIGGER operation.
func buildCreateTriggerOp(c *Catalog, schemaName, tableName string, trig *Trigger) MigrationOp {
	var rel *Relation
	if c != nil && trig != nil {
		rel = c.GetRelationByOID(trig.RelOID)
	}
	qualifiedTable := formatRelationName(schemaName, tableName, rel)

	var b strings.Builder
	if trig.IsConstraint {
		b.WriteString("CREATE CONSTRAINT TRIGGER ")
	} else {
		b.WriteString("CREATE TRIGGER ")
	}
	b.WriteString(quoteIdentAlways(trig.Name))
	b.WriteString(" ")

	// Timing.
	switch trig.Timing {
	case TriggerBefore:
		b.WriteString("BEFORE")
	case TriggerAfter:
		b.WriteString("AFTER")
	case TriggerInsteadOf:
		b.WriteString("INSTEAD OF")
	}
	b.WriteString(" ")

	// Events.
	var events []string
	if trig.Events&TriggerEventInsert != 0 {
		events = append(events, "INSERT")
	}
	if trig.Events&TriggerEventUpdate != 0 {
		if len(trig.Columns) > 0 {
			// Resolve column attnums to names from the relation.
			rel := c.GetRelation(schemaName, tableName)
			var cols []string
			for _, attnum := range trig.Columns {
				if rel != nil {
					for _, col := range rel.Columns {
						if col.AttNum == attnum {
							cols = append(cols, quoteIdentAlways(col.Name))
							break
						}
					}
				}
			}
			if len(cols) > 0 {
				events = append(events, "UPDATE OF "+strings.Join(cols, ", "))
			} else {
				events = append(events, "UPDATE")
			}
		} else {
			events = append(events, "UPDATE")
		}
	}
	if trig.Events&TriggerEventDelete != 0 {
		events = append(events, "DELETE")
	}
	if trig.Events&TriggerEventTruncate != 0 {
		events = append(events, "TRUNCATE")
	}
	b.WriteString(strings.Join(events, " OR "))

	// ON table.
	b.WriteString(" ON ")
	b.WriteString(qualifiedTable)

	// DEFERRABLE / INITIALLY DEFERRED (constraint triggers only, after ON table).
	if trig.IsConstraint && trig.Deferrable {
		b.WriteString(" DEFERRABLE")
		if trig.Initdeferred {
			b.WriteString(" INITIALLY DEFERRED")
		}
	}

	// REFERENCING transition tables.
	if trig.OldTransitionName != "" || trig.NewTransitionName != "" {
		b.WriteString(" REFERENCING")
		if trig.OldTransitionName != "" {
			b.WriteString(" OLD TABLE AS ")
			b.WriteString(quoteIdentAlways(trig.OldTransitionName))
		}
		if trig.NewTransitionName != "" {
			b.WriteString(" NEW TABLE AS ")
			b.WriteString(quoteIdentAlways(trig.NewTransitionName))
		}
	}

	// FOR EACH ROW / FOR EACH STATEMENT.
	if trig.ForEachRow {
		b.WriteString(" FOR EACH ROW")
	} else {
		b.WriteString(" FOR EACH STATEMENT")
	}

	// WHEN clause.
	if trig.WhenExpr != "" {
		b.WriteString(" WHEN (")
		b.WriteString(trig.WhenExpr)
		b.WriteString(")")
	}

	// EXECUTE FUNCTION.
	funcName := resolveTriggerFuncName(c, trig.FuncOID)
	if funcName == "" {
		funcName = "unknown_function"
	}
	b.WriteString(" EXECUTE FUNCTION ")
	b.WriteString(quoteIdentAlways(funcName))
	b.WriteString("(")
	for i, arg := range trig.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(quoteLiteral(arg))
	}
	b.WriteString(")")

	return MigrationOp{
		Type:          OpCreateTrigger,
		SchemaName:    schemaName,
		ObjectName:    trig.Name,
		ParentObject:  tableName,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'T',
		ObjOID:        trig.OID,
		Priority:      PriorityTrigger,
	}
}

// buildCreateAndModeTriggerOps creates a trigger and sets the firing mode it
// has in the target: CREATE TRIGGER always starts it in the origin mode.
// PostgreSQL refuses mode changes on views, so only tables need one.
func buildCreateAndModeTriggerOps(c *Catalog, schemaName, tableName string, trig *Trigger) []MigrationOp {
	ops := []MigrationOp{buildCreateTriggerOp(c, schemaName, tableName, trig)}
	var mode string
	switch trig.Enabled {
	case 'D':
		mode = "DISABLE TRIGGER"
	case 'A':
		mode = "ENABLE ALWAYS TRIGGER"
	case 'R':
		mode = "ENABLE REPLICA TRIGGER"
	default:
		return ops
	}
	var rel *Relation
	if c != nil {
		rel = c.GetRelationByOID(trig.RelOID)
	}
	ops = append(ops, MigrationOp{
		Type:          OpAlterTable,
		SchemaName:    schemaName,
		ObjectName:    trig.Name,
		ParentObject:  tableName,
		SQL:           fmt.Sprintf("ALTER TABLE %s %s %s", formatRelationName(schemaName, tableName, rel), mode, quoteIdentAlways(trig.Name)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'g',
		ObjOID:        trig.OID,
		Priority:      PriorityTrigger,
	})
	return ops
}
