package catalog

import "sort"

// diffTriggers compares triggers between two versions of a relation.
// Identity key is the trigger Name. A partition's clone of its parent's
// trigger is left out: PostgreSQL creates and drops it with the parent's.
func diffTriggers(from, to *Catalog, fromRelOID, toRelOID uint32) []TriggerDiffEntry {
	// Build name→*Trigger maps; retain full maps to avoid treating
	// extension-membership (SkipDump) transitions as trigger create/drop.
	fromAll := make(map[string]*Trigger)
	fromMap := make(map[string]*Trigger)
	for _, trig := range from.TriggersOf(fromRelOID) {
		if trig.ParentTriggerOID == 0 {
			fromAll[trig.Name] = trig
			if !from.isExtensionMember(trig.OID) {
				fromMap[trig.Name] = trig
			}
		}
	}
	toAll := make(map[string]*Trigger)
	toMap := make(map[string]*Trigger)
	for _, trig := range to.TriggersOf(toRelOID) {
		if trig.ParentTriggerOID == 0 {
			toAll[trig.Name] = trig
			if !to.isExtensionMember(trig.OID) {
				toMap[trig.Name] = trig
			}
		}
	}

	var result []TriggerDiffEntry

	// Dropped: in from (user-managed) and not present in to at all.
	for name, fromTrig := range fromMap {
		if _, ok := toAll[name]; !ok {
			result = append(result, TriggerDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fromTrig,
			})
		}
	}

	// Added or modified: in to (user-managed).
	for name, toTrig := range toMap {
		fromTrig, ok := fromMap[name]
		if !ok {
			if prevTrig, existed := fromAll[name]; existed {
				if triggersChanged(from, to, prevTrig, toTrig) {
					result = append(result, TriggerDiffEntry{
						Action: DiffModify,
						Name:   name,
						From:   prevTrig,
						To:     toTrig,
					})
				}
				continue
			}
			result = append(result, TriggerDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     toTrig,
			})
			continue
		}

		// Both exist — compare fields.
		if triggersChanged(from, to, fromTrig, toTrig) {
			result = append(result, TriggerDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fromTrig,
				To:     toTrig,
			})
		}
	}

	// Sort by name for determinism.
	sort.Slice(result, func(i, j int) bool {
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// triggersChanged returns true if any compared property differs between two triggers.
func triggersChanged(fromCat, toCat *Catalog, a, b *Trigger) bool {
	if a.Timing != b.Timing {
		return true
	}
	if a.Events != b.Events {
		return true
	}
	if a.ForEachRow != b.ForEachRow {
		return true
	}
	if a.WhenExpr != b.WhenExpr {
		return true
	}
	if a.Enabled != b.Enabled {
		return true
	}
	if a.OldTransitionName != b.OldTransitionName {
		return true
	}
	if a.NewTransitionName != b.NewTransitionName {
		return true
	}

	// Compare function by resolved name, not OID.
	fromFuncName := resolveTriggerFuncName(fromCat, a.FuncOID)
	toFuncName := resolveTriggerFuncName(toCat, b.FuncOID)
	if fromFuncName != toFuncName {
		return true
	}

	// Compare UPDATE OF columns.
	if !sameColumns(fromCat, a.RelOID, a.Columns, toCat, b.RelOID, b.Columns) {
		return true
	}

	// Compare arguments.
	if !stringSliceEqual(a.Args, b.Args) {
		return true
	}

	// Compare constraint trigger properties.
	if a.IsConstraint != b.IsConstraint {
		return true
	}
	if a.Deferrable != b.Deferrable {
		return true
	}
	if a.Initdeferred != b.Initdeferred {
		return true
	}

	return false
}

// resolveTriggerFuncName resolves a trigger's FuncOID to a schema-qualified function name when user-defined.
func resolveTriggerFuncName(c *Catalog, funcOID uint32) string {
	if up := c.userProcs[funcOID]; up != nil && up.Schema != nil {
		return quoteIdentAlways(up.Schema.Name) + "." + quoteIdentAlways(up.Name)
	}
	p := c.procByOID[funcOID]
	if p == nil {
		return ""
	}
	return quoteIdentAlways(p.Name)
}
