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
				ops = append(ops, createAllRelationRuleOps(to, rel.SchemaName, rel.Name, rel.To)...)
				continue
			}
			ops = append(ops, ruleOpsForRelation(from, to, rel.SchemaName, rel.Name, rel.From, rel.To, rel.Rules)...)
		case DiffAdd:
			if rel.To != nil {
				ops = append(ops, createAllRelationRuleOps(to, rel.SchemaName, rel.Name, rel.To)...)
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

func createAllRelationRuleOps(to *Catalog, schemaName, relName string, toRel *Relation) []MigrationOp {
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
		ops = append(ops, buildCreateAndModeRuleOps(to, schemaName, relName, toRel, toRel.Rules[name])...)
	}
	return ops
}

func ruleOpsForRelation(from, to *Catalog, schemaName, relName string, fromRel, toRel *Relation, rules []RuleDiffEntry) []MigrationOp {
	var ops []MigrationOp

	for _, re := range rules {
		switch re.Action {
		case DiffAdd:
			if re.To == nil || toRel == nil {
				continue
			}
			ops = append(ops, buildCreateAndModeRuleOps(to, schemaName, relName, toRel, re.To)...)
		case DiffDrop:
			if re.From == nil {
				continue
			}
			ops = append(ops, buildDropRuleOpForRel(schemaName, relName, fromRel, re.From))
		case DiffModify:
			if re.From == nil || re.To == nil || toRel == nil {
				continue
			}
			if ruleDefinitionEqual(from, to, re.From, re.To) {
				if modeOp, ok := buildRuleModeOp(schemaName, relName, toRel, re.To); ok {
					ops = append(ops, modeOp)
				}
			} else if re.From.Event == re.To.Event {
				ops = append(ops, buildCreateRuleOpWithReplace(to, schemaName, relName, re.To, true))
				if normalizeRuleMode(re.From.Mode) != normalizeRuleMode(re.To.Mode) {
					if modeOp, ok := buildRuleModeOp(schemaName, relName, toRel, re.To); ok {
						ops = append(ops, modeOp)
					}
				}
			} else {
				ops = append(ops, buildDropRuleOpForRel(schemaName, relName, fromRel, re.From))
				ops = append(ops, buildCreateAndModeRuleOps(to, schemaName, relName, toRel, re.To)...)
				if cOp, ok := buildRuleCommentRestoreOp(from, to, schemaName, relName, re.From, re.To); ok {
					ops = append(ops, cOp)
				}
			}
		}
	}
	return ops
}

func buildRuleCommentRestoreOp(from, to *Catalog, schemaName, relName string, fromRule, toRule *Rule) (MigrationOp, bool) {
	if to == nil || toRule == nil {
		return MigrationOp{}, false
	}
	toTxt, ok := to.GetComment('R', toRule.OID, 0)
	if !ok || toTxt == "" {
		return MigrationOp{}, false
	}
	if from != nil && fromRule != nil {
		fromTxt, _ := from.GetComment('R', fromRule.OID, 0)
		if fromTxt != toTxt {
			// Already emitted by generateCommentDDL.
			return MigrationOp{}, false
		}
	}
	desc := schemaName + "\x00" + relName + "\x00" + toRule.Name
	csql := formatCommentSQL(to, 'R', desc, 0, toTxt)
	if csql == "" {
		return MigrationOp{}, false
	}
	return MigrationOp{
		Type:          OpComment,
		ObjectName:    desc,
		SQL:           csql,
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'R',
		Priority:      PriorityMetadata,
	}, true
}

func buildDropRuleOp(schemaName, relName string, fromRule *Rule) MigrationOp {
	return buildDropRuleOpForRel(schemaName, relName, nil, fromRule)
}

func buildDropRuleOpForRel(schemaName, relName string, fromRel *Relation, fromRule *Rule) MigrationOp {
	qualifiedRel := formatRelationName(schemaName, relName, fromRel)
	return MigrationOp{
		Type:          OpDropRule,
		SchemaName:    schemaName,
		ObjectName:    fromRule.Name,
		ParentObject:  relName,
		SQL:           fmt.Sprintf("DROP RULE %s ON %s", quoteIdentAlways(fromRule.Name), qualifiedRel),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'R',
		ObjOID:        fromRule.OID,
		Priority:      PriorityRule,
	}
}

func buildCreateAndModeRuleOps(to *Catalog, schemaName, relName string, toRel *Relation, r *Rule) []MigrationOp {
	ops := []MigrationOp{buildCreateRuleOp(to, schemaName, relName, r)}
	if normalizeRuleMode(r.Mode) != 'O' {
		if modeOp, ok := buildRuleModeOp(schemaName, relName, toRel, r); ok {
			ops = append(ops, modeOp)
		}
	}
	return ops
}

func buildCreateRuleOp(to *Catalog, schemaName, relName string, r *Rule) MigrationOp {
	return buildCreateRuleOpWithReplace(to, schemaName, relName, r, false)
}

func buildCreateRuleOpWithReplace(to *Catalog, schemaName, relName string, r *Rule, replace bool) MigrationOp {
	var toRel *Relation
	if to != nil && r != nil {
		toRel = to.GetRelationByOID(r.OwnerRelOID)
	}
	qualifiedRel := formatRelationName(schemaName, relName, toRel)
	whereSQL := to.deparseRuleWhere(r)
	actionsSQL := to.deparseRuleActions(r)

	var b strings.Builder
	if replace {
		b.WriteString("CREATE OR REPLACE RULE ")
	} else {
		b.WriteString("CREATE RULE ")
	}
	b.WriteString(quoteIdentAlways(r.Name))
	b.WriteString(" AS ON ")
	b.WriteString(ruleEventKeyword(r.Event))
	b.WriteString(" TO ")
	b.WriteString(qualifiedRel)
	if whereSQL != "" {
		b.WriteString(" WHERE ")
		b.WriteString(whereSQL)
	}
	b.WriteString(" DO ")
	if r.Instead {
		b.WriteString("INSTEAD ")
	} else {
		b.WriteString("ALSO ")
	}
	switch len(actionsSQL) {
	case 0:
		b.WriteString("NOTHING")
	case 1:
		b.WriteString(actionsSQL[0])
	default:
		b.WriteString("(")
		b.WriteString(strings.Join(actionsSQL, "; "))
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
		ObjType:       'R',
		ObjOID:        r.OID,
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
	qualifiedRel := formatRelationName(schemaName, relName, toRel)
	return MigrationOp{
		Type:          OpAlterTable,
		SchemaName:    schemaName,
		ObjectName:    r.Name,
		ParentObject:  relName,
		SQL:           fmt.Sprintf("ALTER TABLE %s %s %s", qualifiedRel, action, quoteIdentAlways(r.Name)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'R',
		ObjOID:        r.OID,
		Priority:      PriorityRule + 1,
	}, true
}

// reconcileRulesWithDroppedRelations ensures that when a table, view, or
// materialized view is dropped and recreated (e.g. view column changes, matview
// definition changes, or dependent view rebuilds), redundant DROP ops for its
// sub-objects are pruned and retained rules, standalone indexes (on views/matviews),
// and triggers (on views) are emitted in PhaseMain.
func reconcileRulesWithDroppedRelations(from, to *Catalog, ops []MigrationOp) []MigrationOp {
	droppedRels := make(map[string]bool)  // "schema\x00name"
	droppedViews := make(map[string]bool) // "schema\x00name" (OpDropView only)
	droppedViewFromOIDs := make(map[uint32]bool)
	createdRules := make(map[string]bool)    // "schema\x00rel\x00rule"
	createdIndexes := make(map[string]bool)  // "schema\x00idx"
	createdTriggers := make(map[string]bool) // "schema\x00rel\x00trig"

	for _, op := range ops {
		if op.Phase == PhasePre && (op.Type == OpDropTable || op.Type == OpDropView) {
			key := op.SchemaName + "\x00" + op.ObjectName
			droppedRels[key] = true
			if op.Type == OpDropView {
				droppedViews[key] = true
				if from != nil {
					if fromRel := from.GetRelation(op.SchemaName, op.ObjectName); fromRel != nil {
						droppedViewFromOIDs[fromRel.OID] = true
					}
				}
			}
		}
		if op.Phase == PhaseMain {
			switch op.Type {
			case OpCreateRule:
				createdRules[op.SchemaName+"\x00"+op.ParentObject+"\x00"+op.ObjectName] = true
			case OpCreateIndex:
				createdIndexes[op.SchemaName+"\x00"+op.ObjectName] = true
			case OpCreateTrigger:
				createdTriggers[op.SchemaName+"\x00"+op.ParentObject+"\x00"+op.ObjectName] = true
			}
		}
	}
	if len(droppedRels) == 0 {
		return ops
	}

	filtered := make([]MigrationOp, 0, len(ops))
	for _, op := range ops {
		if op.Phase == PhasePre {
			if op.Type == OpDropRule && droppedRels[op.SchemaName+"\x00"+op.ParentObject] {
				continue
			}
			if op.Type == OpDropTrigger && droppedViews[op.SchemaName+"\x00"+op.ParentObject] {
				continue
			}
			if op.Type == OpDropIndex && from != nil && op.ObjOID != 0 {
				if idx := from.indexes[op.ObjOID]; idx != nil && droppedViewFromOIDs[idx.RelOID] {
					continue
				}
			}
		}
		filtered = append(filtered, op)
	}

	relKeys := make([]string, 0, len(droppedRels))
	for relKeyStr := range droppedRels {
		relKeys = append(relKeys, relKeyStr)
	}
	sort.Strings(relKeys)

	for _, relKeyStr := range relKeys {
		schemaName, relName, _ := strings.Cut(relKeyStr, "\x00")
		toRel := to.GetRelation(schemaName, relName)
		if toRel == nil {
			continue
		}
		if len(toRel.Rules) > 0 {
			names := make([]string, 0, len(toRel.Rules))
			for name, r := range toRel.Rules {
				if r != nil && !createdRules[schemaName+"\x00"+relName+"\x00"+name] {
					names = append(names, name)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				toRule := toRel.Rules[name]
				filtered = append(filtered, buildCreateAndModeRuleOps(to, schemaName, relName, toRel, toRule)...)
				var fromRule *Rule
				if from != nil {
					if fromRel := from.GetRelation(schemaName, relName); fromRel != nil && fromRel.Rules != nil {
						fromRule = fromRel.Rules[name]
					}
				}
				if cOp, ok := buildRuleCommentRestoreOp(from, to, schemaName, relName, fromRule, toRule); ok {
					filtered = append(filtered, cOp)
				}
			}
		}
		if droppedViews[relKeyStr] {
			for _, idx := range to.IndexesOf(toRel.OID) {
				if idx.ConstraintOID != 0 || idx.ParentIndexOID != 0 {
					continue
				}
				if createdIndexes[schemaName+"\x00"+idx.Name] {
					continue
				}
				sql := buildCreateIndexSQL(idx, toRel, schemaName)
				filtered = append(filtered, MigrationOp{
					Type:          OpCreateIndex,
					SchemaName:    schemaName,
					ObjectName:    idx.Name,
					SQL:           sql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'i',
					ObjOID:        idx.OID,
					Priority:      PriorityIndex,
				})
				filtered = append(filtered, buildIndexColumnRenameOps(schemaName, idx, defaultIndexColumnNames(to, idx), effectiveIndexColumnNames(to, idx))...)
			}
			for _, trig := range to.TriggersOf(toRel.OID) {
				if trig.ParentTriggerOID != 0 {
					continue
				}
				if createdTriggers[schemaName+"\x00"+relName+"\x00"+trig.Name] {
					continue
				}
				filtered = append(filtered, buildCreateTriggerOp(to, schemaName, relName, trig))
			}
		}
	}
	return filtered
}
