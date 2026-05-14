package catalog

import (
	"fmt"
	"sort"
	"strings"
)

func generateEventTriggerDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.EventTriggers {
		switch entry.Action {
		case DiffAdd:
			if entry.To != nil {
				ops = append(ops, buildCreateEventTriggerOps(to, entry.To)...)
			}
		case DiffDrop:
			if entry.From != nil {
				ops = append(ops, buildDropEventTriggerOp(entry.From))
			}
		case DiffModify:
			if structuralEventTriggerChange(from, to, entry.From, entry.To) {
				if entry.From != nil {
					ops = append(ops, buildDropEventTriggerOp(entry.From))
				}
				if entry.To != nil {
					ops = append(ops, buildCreateEventTriggerOps(to, entry.To)...)
				}
			} else if entry.To != nil {
				ops = append(ops, buildAlterEventTriggerEnableOp(entry.To))
			}
		}
	}

	sort.Slice(ops, func(i, j int) bool {
		if ops[i].ObjectName != ops[j].ObjectName {
			return ops[i].ObjectName < ops[j].ObjectName
		}
		return eventTriggerOpRank(ops[i].Type) < eventTriggerOpRank(ops[j].Type)
	})
	return ops
}

func eventTriggerOpRank(t MigrationOpType) int {
	switch t {
	case OpDropEventTrigger:
		return 0
	case OpCreateEventTrigger:
		return 1
	case OpAlterEventTrigger:
		return 2
	default:
		return 3
	}
}

func buildDropEventTriggerOp(evt *EventTrigger) MigrationOp {
	return MigrationOp{
		Type:          OpDropEventTrigger,
		ObjectName:    evt.Name,
		SQL:           fmt.Sprintf("DROP EVENT TRIGGER %s", quoteIdentAlways(evt.Name)),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'E',
		ObjOID:        evt.OID,
		Priority:      PriorityEventTrigger,
	}
}

func buildCreateEventTriggerOps(c *Catalog, evt *EventTrigger) []MigrationOp {
	var ops []MigrationOp
	var b strings.Builder
	b.WriteString("CREATE EVENT TRIGGER ")
	b.WriteString(quoteIdentAlways(evt.Name))
	b.WriteString(" ON ")
	b.WriteString(quoteIdentAlways(evt.EventName))
	if len(evt.Tags) > 0 {
		var tags []string
		for _, tag := range evt.Tags {
			tags = append(tags, quoteLiteral(tag))
		}
		b.WriteString("\n         WHEN TAG IN (")
		b.WriteString(strings.Join(tags, ", "))
		b.WriteString(")")
	}
	b.WriteString("\n   EXECUTE FUNCTION ")
	b.WriteString(resolveEventTriggerFuncSQLName(c, evt.FuncOID))
	b.WriteString("()")

	ops = append(ops, MigrationOp{
		Type:          OpCreateEventTrigger,
		ObjectName:    evt.Name,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'E',
		ObjOID:        evt.OID,
		Priority:      PriorityEventTrigger,
	})
	if evt.Enabled != 'O' {
		ops = append(ops, buildAlterEventTriggerEnableOp(evt))
	}
	return ops
}

func buildAlterEventTriggerEnableOp(evt *EventTrigger) MigrationOp {
	var action string
	switch evt.Enabled {
	case 'D':
		action = "DISABLE"
	case 'A':
		action = "ENABLE ALWAYS"
	case 'R':
		action = "ENABLE REPLICA"
	default:
		action = "ENABLE"
	}
	return MigrationOp{
		Type:          OpAlterEventTrigger,
		ObjectName:    evt.Name,
		SQL:           fmt.Sprintf("ALTER EVENT TRIGGER %s %s", quoteIdentAlways(evt.Name), action),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'E',
		ObjOID:        evt.OID,
		Priority:      PriorityEventTrigger,
	}
}
