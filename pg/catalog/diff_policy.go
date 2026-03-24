package catalog

import "sort"

// diffPolicies compares policies between two versions of a relation.
// Identity key is the policy Name.
func diffPolicies(from, to *Catalog, fromRelOID, toRelOID uint32) []PolicyDiffEntry {
	// Build name→*Policy maps.
	fromMap := make(map[string]*Policy)
	for _, p := range from.policiesByRel[fromRelOID] {
		fromMap[p.Name] = p
	}
	toMap := make(map[string]*Policy)
	for _, p := range to.policiesByRel[toRelOID] {
		toMap[p.Name] = p
	}

	var result []PolicyDiffEntry

	// Dropped: in from but not in to.
	for name, fromPol := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, PolicyDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fromPol,
			})
		}
	}

	// Added or modified: in to.
	for name, toPol := range toMap {
		fromPol, ok := fromMap[name]
		if !ok {
			result = append(result, PolicyDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     toPol,
			})
			continue
		}

		// Both exist — compare fields.
		if policiesChanged(fromPol, toPol) {
			result = append(result, PolicyDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fromPol,
				To:     toPol,
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

// policiesChanged returns true if any compared property differs between two policies.
func policiesChanged(a, b *Policy) bool {
	if a.CmdType != b.CmdType {
		return true
	}
	if a.Permissive != b.Permissive {
		return true
	}
	if a.UsingExpr != b.UsingExpr {
		return true
	}
	if a.CheckExpr != b.CheckExpr {
		return true
	}

	// Compare roles — sort copies to avoid mutating originals.
	aRoles := make([]string, len(a.Roles))
	copy(aRoles, a.Roles)
	sort.Strings(aRoles)
	bRoles := make([]string, len(b.Roles))
	copy(bRoles, b.Roles)
	sort.Strings(bRoles)
	if !stringSliceEqual(aRoles, bRoles) {
		return true
	}

	return false
}
