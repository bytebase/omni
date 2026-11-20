package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateRuleDDL produces CREATE RULE, DROP RULE, and ALTER TABLE ... RULE
// firing-mode operations from the schema diff.
func generateRuleDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, rel := range diff.Relations {
		switch rel.Action {
		case DiffModify:
			if rel.From != nil && rel.To != nil && tableRebuilt(from, to, rel.From, rel.To) {
				ops = append(ops, createAllRelationRuleOps(rel.SchemaName, rel.Name, rel.To)...)
				continue
			}
			ops = append(ops, ruleOpsForRelation(rel.SchemaName, rel.Name, rel.From, rel.To, rel.Rules)...)
		case DiffAdd:
			if rel.To != nil {
				ops = append(ops, createAllRelationRuleOps(rel.SchemaName, rel.Name, rel.To)...)
			}
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Phase != ops[j].Phase {
			return ops[i].Phase < ops[j].Phase
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		if ops[i].ParentObject != ops[j].ParentObject {
			return ops[i].ParentObject < ops[j].ParentObject
		}
		if ops[i].ObjectName != ops[j].ObjectName {
			return ops[i].ObjectName < ops[j].ObjectName
		}
		return ops[i].Priority < ops[j].Priority
	})

	return ops
}

func createAllRelationRuleOps(schemaName, relName string, toRel *Relation) []MigrationOp {
	if toRel == nil || len(toRel.Rules) == 0 {
		return nil
	}
	names := make([]string, 0, len(toRel.Rules))
	for name, r := range toRel.Rules {
		if r != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	var ops []MigrationOp
	for _, name := range names {
		ops = append(ops, buildCreateAndModeRuleOps(schemaName, relName, toRel, toRel.Rules[name])...)
	}
	return ops
}

func ruleOpsForRelation(schemaName, relName string, fromRel, toRel *Relation, rules []RuleDiffEntry) []MigrationOp {
	var ops []MigrationOp
	var fromOID uint32
	if fromRel != nil {
		fromOID = fromRel.OID
	}

	for _, re := range rules {
		switch re.Action {
		case DiffAdd:
			if re.To == nil || toRel == nil {
				continue
			}
			ops = append(ops, buildCreateAndModeRuleOps(schemaName, relName, toRel, re.To)...)
		case DiffDrop:
			if re.From == nil {
				continue
			}
			ops = append(ops, buildDropRuleOp(schemaName, relName, re.From.Name, fromOID))
		case DiffModify:
			if re.From == nil || re.To == nil || toRel == nil {
				continue
			}
			if ruleDefinitionEqual(re.From, re.To) {
				if modeOp, ok := buildRuleModeOp(schemaName, relName, toRel, re.To); ok {
					ops = append(ops, modeOp)
				}
			} else {
				ops = append(ops, buildDropRuleOp(schemaName, relName, re.From.Name, fromOID))
				ops = append(ops, buildCreateAndModeRuleOps(schemaName, relName, toRel, re.To)...)
			}
		}
	}
	return ops
}

func buildDropRuleOp(schemaName, relName, ruleName string, fromRelOID uint32) MigrationOp {
	qualifiedRel := migrationQualifiedName(schemaName, relName)
	return MigrationOp{
		Type:          OpDropRule,
		SchemaName:    schemaName,
		ObjectName:    ruleName,
		ParentObject:  relName,
		SQL:           fmt.Sprintf("DROP RULE %s ON %s", quoteIdentAlways(ruleName), qualifiedRel),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'r',
		ObjOID:        fromRelOID,
		Priority:      PriorityRule,
	}
}

func buildCreateAndModeRuleOps(schemaName, relName string, toRel *Relation, r *Rule) []MigrationOp {
	ops := []MigrationOp{buildCreateRuleOp(schemaName, relName, toRel.OID, r)}
	if normalizeRuleMode(r.Mode) != 'O' {
		if modeOp, ok := buildRuleModeOp(schemaName, relName, toRel, r); ok {
			ops = append(ops, modeOp)
		}
	}
	return ops
}

func buildCreateRuleOp(schemaName, relName string, toRelOID uint32, r *Rule) MigrationOp {
	qualifiedRel := migrationQualifiedName(schemaName, relName)
	var b strings.Builder
	b.WriteString("CREATE RULE ")
	b.WriteString(quoteIdentAlways(r.Name))
	b.WriteString(" AS ON ")
	b.WriteString(ruleEventKeyword(r.Event))
	b.WriteString(" TO ")
	b.WriteString(qualifiedRel)
	if r.WhereSQL != "" {
		b.WriteString(" WHERE ")
		b.WriteString(r.WhereSQL)
	}
	b.WriteString(" DO ")
	if r.Instead {
		b.WriteString("INSTEAD ")
	} else {
		b.WriteString("ALSO ")
	}
	switch len(r.ActionsSQL) {
	case 0:
		b.WriteString("NOTHING")
	case 1:
		b.WriteString(r.ActionsSQL[0])
	default:
		b.WriteString("(")
		b.WriteString(strings.Join(r.ActionsSQL, "; "))
		b.WriteString(")")
	}

	return MigrationOp{
		Type:          OpCreateRule,
		SchemaName:    schemaName,
		ObjectName:    r.Name,
		ParentObject:  relName,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        toRelOID,
		Priority:      PriorityRule,
	}
}

func buildRuleModeOp(schemaName, relName string, toRel *Relation, r *Rule) (MigrationOp, bool) {
	if toRel == nil || (toRel.RelKind != 'r' && toRel.RelKind != 'p') {
		return MigrationOp{}, false
	}
	var action string
	switch normalizeRuleMode(r.Mode) {
	case 'D':
		action = "DISABLE RULE"
	case 'A':
		action = "ENABLE ALWAYS RULE"
	case 'R':
		action = "ENABLE REPLICA RULE"
	default:
		action = "ENABLE RULE"
	}
	qualifiedRel := migrationQualifiedName(schemaName, relName)
	return MigrationOp{
		Type:          OpAlterTable,
		SchemaName:    schemaName,
		ObjectName:    r.Name,
		ParentObject:  relName,
		SQL:           fmt.Sprintf("ALTER TABLE %s %s %s", qualifiedRel, action, quoteIdentAlways(r.Name)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        toRel.OID,
		Priority:      PriorityRule + 1,
	}, true
}

// reconcileRulesWithDroppedRelations ensures that when a table or view is
// dropped and recreated (e.g. view column changes or dependent view rebuilds),
// redundant DROP RULE ops on the dropped relation are pruned and all target
// rules on the recreated relation are emitted in PhaseMain.
func reconcileRulesWithDroppedRelations(from, to *Catalog, ops []MigrationOp) []MigrationOp {
	droppedRels := make(map[string]bool)  // "schema\x00name"
	createdRules := make(map[string]bool) // "schema\x00rel\x00rule"
	for _, op := range ops {
		if op.Phase == PhasePre && (op.Type == OpDropTable || op.Type == OpDropView) {
			droppedRels[op.SchemaName+"\x00"+op.ObjectName] = true
		}
		if op.Phase == PhaseMain && op.Type == OpCreateRule {
			createdRules[op.SchemaName+"\x00"+op.ParentObject+"\x00"+op.ObjectName] = true
		}
	}
	if len(droppedRels) == 0 {
		return ops
	}

	filtered := make([]MigrationOp, 0, len(ops))
	for _, op := range ops {
		if op.Phase == PhasePre && op.Type == OpDropRule && droppedRels[op.SchemaName+"\x00"+op.ParentObject] {
			continue
		}
		filtered = append(filtered, op)
	}

	for relKeyStr := range droppedRels {
		schemaName, relName, _ := strings.Cut(relKeyStr, "\x00")
		toRel := to.GetRelation(schemaName, relName)
		if toRel == nil || len(toRel.Rules) == 0 {
			continue
		}
		names := make([]string, 0, len(toRel.Rules))
		for name, r := range toRel.Rules {
			if r != nil && !createdRules[schemaName+"\x00"+relName+"\x00"+name] {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		for _, name := range names {
			filtered = append(filtered, buildCreateAndModeRuleOps(schemaName, relName, toRel, toRel.Rules[name])...)
		}
	}
	return filtered
}
