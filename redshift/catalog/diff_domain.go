package catalog

import "sort"

// diffDomains compares domain types between two catalogs and returns diff entries.
// Identity key is schemaName + "." + typeName, resolved via typeByOID.
func diffDomains(from, to *Catalog) []DomainDiffEntry {
	type domainInfo struct {
		schemaName string
		typeName   string
		domain     *DomainType
		catalog    *Catalog
	}

	// Build name-based maps from both catalogs.
	fromMap := make(map[string]*domainInfo)
	for typeOID, dt := range from.domainTypes {
		bt := from.typeByOID[typeOID]
		if bt == nil {
			continue
		}
		schema := from.schemas[bt.Namespace]
		if schema == nil {
			continue
		}
		key := schema.Name + "." + bt.TypeName
		fromMap[key] = &domainInfo{
			schemaName: schema.Name,
			typeName:   bt.TypeName,
			domain:     dt,
			catalog:    from,
		}
	}

	toMap := make(map[string]*domainInfo)
	for typeOID, dt := range to.domainTypes {
		bt := to.typeByOID[typeOID]
		if bt == nil {
			continue
		}
		schema := to.schemas[bt.Namespace]
		if schema == nil {
			continue
		}
		key := schema.Name + "." + bt.TypeName
		toMap[key] = &domainInfo{
			schemaName: schema.Name,
			typeName:   bt.TypeName,
			domain:     dt,
			catalog:    to,
		}
	}

	var entries []DomainDiffEntry

	// Dropped: in from but not in to.
	for key, fi := range fromMap {
		if _, ok := toMap[key]; !ok {
			entries = append(entries, DomainDiffEntry{
				Action:     DiffDrop,
				SchemaName: fi.schemaName,
				Name:       fi.typeName,
				From:       fi.domain,
			})
		}
	}

	// Added: in to but not in from.
	for key, ti := range toMap {
		if _, ok := fromMap[key]; !ok {
			entries = append(entries, DomainDiffEntry{
				Action:     DiffAdd,
				SchemaName: ti.schemaName,
				Name:       ti.typeName,
				To:         ti.domain,
			})
		}
	}

	// Modified: in both — compare fields.
	for key, fi := range fromMap {
		ti, ok := toMap[key]
		if !ok {
			continue
		}
		if domainsChanged(fi.domain, fi.catalog, ti.domain, ti.catalog) {
			entries = append(entries, DomainDiffEntry{
				Action:     DiffModify,
				SchemaName: fi.schemaName,
				Name:       fi.typeName,
				From:       fi.domain,
				To:         ti.domain,
			})
		}
	}

	// Sort for determinism: by schema name, then type name.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].SchemaName != entries[j].SchemaName {
			return entries[i].SchemaName < entries[j].SchemaName
		}
		if entries[i].Name != entries[j].Name {
			return entries[i].Name < entries[j].Name
		}
		return entries[i].Action < entries[j].Action
	})

	return entries
}

// domainsChanged returns true if any compared property differs between two domains.
func domainsChanged(a *DomainType, aCat *Catalog, b *DomainType, bCat *Catalog) bool {
	// Compare base type using FormatType (never raw OIDs).
	if aCat.FormatType(a.BaseTypeOID, a.BaseTypMod) != bCat.FormatType(b.BaseTypeOID, b.BaseTypMod) {
		return true
	}
	if a.NotNull != b.NotNull {
		return true
	}
	if a.Default != b.Default {
		return true
	}
	if domainConstraintsChanged(a.Constraints, b.Constraints) {
		return true
	}
	return false
}

// domainConstraintsChanged returns true if domain constraints differ.
// Constraints are compared by name; for matching names, CheckExpr is compared.
func domainConstraintsChanged(a, b []*DomainConstraint) bool {
	if len(a) != len(b) {
		return true
	}
	aMap := make(map[string]string, len(a))
	for _, c := range a {
		aMap[c.Name] = c.CheckExpr
	}
	for _, c := range b {
		expr, ok := aMap[c.Name]
		if !ok {
			return true // constraint in b not in a
		}
		if expr != c.CheckExpr {
			return true // expression changed
		}
	}
	return false
}
