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
			ops = append(ops, triggerOpsForRelation(to, rel.SchemaName, rel.Name, rel.Triggers)...)
		case DiffAdd:
			// For newly added relations, triggers are not in the diff entry's
			// Triggers slice (the relation itself is new). Generate triggers
			// from the target catalog.
			if rel.To != nil {
				trigs := to.TriggersOf(rel.To.OID)
				for _, trig := range trigs {
					ops = append(ops, buildCreateTriggerOp(to, rel.SchemaName, rel.Name, trig))
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
func triggerOpsForRelation(to *Catalog, schemaName, tableName string, triggers []TriggerDiffEntry) []MigrationOp {
	var ops []MigrationOp

	for _, te := range triggers {
		switch te.Action {
		case DiffAdd:
			if te.To == nil {
				continue
			}
			ops = append(ops, buildCreateTriggerOp(to, schemaName, tableName, te.To))
		case DiffDrop:
			if te.From == nil {
				continue
			}
			ops = append(ops, buildDropTriggerOp(schemaName, tableName, te.From.Name))
		case DiffModify:
			// Modified trigger = DROP old + CREATE new.
			if te.From != nil {
				ops = append(ops, buildDropTriggerOp(schemaName, tableName, te.From.Name))
			}
			if te.To != nil {
				ops = append(ops, buildCreateTriggerOp(to, schemaName, tableName, te.To))
			}
		}
	}

	return ops
}

// buildDropTriggerOp creates a DROP TRIGGER operation.
func buildDropTriggerOp(schemaName, tableName, trigName string) MigrationOp {
	qualifiedTable := migrationQualifiedName(schemaName, tableName)
	return MigrationOp{
		Type:          OpDropTrigger,
		SchemaName:    schemaName,
		ObjectName:    trigName,
		ParentObject:  tableName,
		SQL:           fmt.Sprintf("DROP TRIGGER %s ON %s", quoteIdentAlways(trigName), qualifiedTable),
		Transactional: true,
	}
}

// buildCreateTriggerOp creates a CREATE TRIGGER operation.
func buildCreateTriggerOp(c *Catalog, schemaName, tableName string, trig *Trigger) MigrationOp {
	qualifiedTable := migrationQualifiedName(schemaName, tableName)

	var b strings.Builder
	b.WriteString("CREATE TRIGGER ")
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
		events = append(events, "UPDATE")
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
	b.WriteString("()")

	return MigrationOp{
		Type:          OpCreateTrigger,
		SchemaName:    schemaName,
		ObjectName:    trig.Name,
		ParentObject:  tableName,
		SQL:           b.String(),
		Transactional: true,
	}
}
