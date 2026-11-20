package catalog

import "sort"

// diffRules compares rewrite rules between two versions of a relation.
func diffRules(fromCat, toCat *Catalog, fromRel, toRel *Relation) []RuleDiffEntry {
	var fromRules, toRules map[string]*Rule
	if fromRel != nil {
		fromRules = fromRel.Rules
	}
	if toRel != nil {
		toRules = toRel.Rules
	}
	if len(fromRules) == 0 && len(toRules) == 0 {
		return nil
	}

	var result []RuleDiffEntry
	for name, fr := range fromRules {
		if fr == nil {
			continue
		}
		if tr := toRules[name]; tr == nil {
			result = append(result, RuleDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fr,
			})
		}
	}
	for name, tr := range toRules {
		if tr == nil {
			continue
		}
		fr := fromRules[name]
		if fr == nil {
			result = append(result, RuleDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     tr,
			})
			continue
		}
		if rulesChanged(fromCat, toCat, fr, tr) {
			result = append(result, RuleDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fr,
				To:     tr,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})
	return result
}

func normalizeRuleMode(m byte) byte {
	if m == 0 {
		return 'O'
	}
	return m
}

func ruleDefinitionEqual(fromCat, toCat *Catalog, a, b *Rule) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Event == b.Event &&
		a.Instead == b.Instead &&
		fromCat.deparseRuleWhere(a) == toCat.deparseRuleWhere(b) &&
		stringSliceEqual(fromCat.deparseRuleActions(a), toCat.deparseRuleActions(b))
}

func rulesChanged(fromCat, toCat *Catalog, a, b *Rule) bool {
	if !ruleDefinitionEqual(fromCat, toCat, a, b) {
		return true
	}
	return normalizeRuleMode(a.Mode) != normalizeRuleMode(b.Mode)
}
