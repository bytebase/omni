package catalog

import "sort"

// Schema represents a PostgreSQL namespace.
type Schema struct {
	OID        uint32
	Name       string
	Owner      string // explicit AUTHORIZATION / ALTER OWNER (for diff/SDL)
	OwnerName  string // effective role owner (for DROP ROLE / DROP OWNED / REASSIGN / RENAME)
	Relations  map[string]*Relation
	Indexes    map[string]*Index
	Sequences  map[string]*Sequence
	Collations map[string]uint32
}

func (s *Schema) effectiveOwner() string {
	if s.OwnerName != "" {
		return s.OwnerName
	}
	return s.Owner
}

// schemaHasObjects returns true if the schema contains any objects.
func (c *Catalog) schemaHasObjects(s *Schema) bool {
	if len(s.Relations) > 0 || len(s.Indexes) > 0 || len(s.Sequences) > 0 || len(s.Collations) > 0 {
		return true
	}
	// Check for user-defined types in this namespace.
	for _, bt := range c.typeByOID {
		if bt.Namespace == s.OID && bt.OID >= FirstNormalObjectId && bt.Elem == 0 && bt.Type != 'c' {
			return true
		}
	}
	// Check for user functions in this schema.
	for _, up := range c.userProcs {
		if up.Schema.OID == s.OID {
			return true
		}
	}
	for _, ext := range c.extensions {
		if ext.SchemaOID == s.OID {
			return true
		}
	}
	for _, nsOID := range c.operatorSchemas {
		if nsOID == s.OID {
			return true
		}
	}
	for _, opc := range c.opClasses {
		if opc.Schema == s.OID {
			return true
		}
	}
	for _, fam := range c.opFamilies {
		if fam.Schema == s.OID {
			return true
		}
	}
	return false
}

// removeFunctionsInSchema removes all user functions in the given schema,
// cascading to any remaining dependents.
func (c *Catalog) removeFunctionsInSchema(s *Schema) {
	var oids []uint32
	for oid, up := range c.userProcs {
		if up.Schema.OID == s.OID {
			oids = append(oids, oid)
		}
	}
	sort.Slice(oids, func(i, j int) bool { return oids[i] < oids[j] })
	for _, oid := range oids {
		up := c.userProcs[oid]
		if up == nil {
			continue
		}
		c.dropDependents('f', oid)
		if c.userProcs[oid] != nil {
			c.removeFunction(oid, up.Name)
		}
	}
}

// owningRangeType finds the range BuiltinType whose companion multirange is multirangeOID.
func (c *Catalog) owningRangeType(multirangeOID uint32) *BuiltinType {
	for rangeOID, rt := range c.rangeTypes {
		if rt.MultirangeOID == multirangeOID {
			return c.typeByOID[rangeOID]
		}
	}
	return nil
}

// removeTypesInSchema removes all user-created types in the given schema,
// cascading to any remaining dependents (including range types whose companion
// multirange still resides in s).
func (c *Catalog) removeTypesInSchema(s *Schema) {
	var types []*BuiltinType
	seen := map[uint32]bool{}
	var orphanMultiranges []*BuiltinType
	for _, bt := range c.typeByOID {
		if bt.Namespace != s.OID || bt.OID < FirstNormalObjectId || bt.Elem != 0 || bt.Type == 'c' {
			continue
		}
		if bt.Type == 'm' {
			if rng := c.owningRangeType(bt.OID); rng != nil {
				if !seen[rng.OID] {
					seen[rng.OID] = true
					types = append(types, rng)
				}
			} else {
				orphanMultiranges = append(orphanMultiranges, bt)
			}
			continue
		}
		if !seen[bt.OID] {
			seen[bt.OID] = true
			types = append(types, bt)
		}
	}
	sort.Slice(types, func(i, j int) bool { return types[i].OID < types[j].OID })
	for _, bt := range types {
		if c.typeByOID[bt.OID] == nil {
			continue
		}
		for _, depRelOID := range c.typeRelationDependents(bt) {
			if depRel := c.relationByOID[depRelOID]; depRel != nil {
				c.cascadeDropRelation(depRel)
			}
		}
		c.dropTypeDependents(bt)
		if c.typeByOID[bt.OID] != nil {
			c.deleteUserType(bt)
		}
	}
	sort.Slice(orphanMultiranges, func(i, j int) bool { return orphanMultiranges[i].OID < orphanMultiranges[j].OID })
	for _, mr := range orphanMultiranges {
		if c.typeByOID[mr.OID] != nil {
			c.dropTypeDependents(mr)
			c.deleteUserType(mr)
		}
	}
}
