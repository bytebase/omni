package ddl

// Schema represents a PostgreSQL namespace.
type Schema struct {
	OID       uint32
	Name      string
	Owner     string // role name of the schema owner (empty = current user)
	Relations map[string]*Relation
	Indexes   map[string]*Index
	Sequences map[string]*Sequence
}

// schemaHasObjects returns true if the schema contains any objects.
func (c *Catalog) schemaHasObjects(s *Schema) bool {
	if len(s.Relations) > 0 || len(s.Indexes) > 0 || len(s.Sequences) > 0 {
		return true
	}
	// Check for user-defined types (enum/domain/range) in this namespace.
	for _, bt := range c.typeByOID {
		if bt.Namespace == s.OID && (bt.Type == 'e' || bt.Type == 'd' || bt.Type == 'r') {
			return true
		}
	}
	// Check for user functions in this schema.
	for _, up := range c.userProcs {
		if up.Schema.OID == s.OID {
			return true
		}
	}
	return false
}

// removeFunctionsInSchema removes all user functions in the given schema.
func (c *Catalog) removeFunctionsInSchema(s *Schema) {
	for oid, up := range c.userProcs {
		if up.Schema.OID == s.OID {
			c.removeFunction(oid, up.Name)
		}
	}
}

// removeTypesInSchema removes all user-created types (enum/domain/range) in the given schema.
func (c *Catalog) removeTypesInSchema(s *Schema) {
	for oid, bt := range c.typeByOID {
		if bt.Namespace != s.OID {
			continue
		}
		switch bt.Type {
		case 'e':
			delete(c.enumTypes, oid)
		case 'd':
			delete(c.domainTypes, oid)
		case 'r':
			delete(c.rangeTypes, oid)
		default:
			continue
		}
		// Remove array type.
		if bt.Array != 0 {
			if at := c.typeByOID[bt.Array]; at != nil {
				delete(c.typeByName, typeKey{ns: at.Namespace, name: at.TypeName})
			}
			delete(c.typeByOID, bt.Array)
		}
		delete(c.typeByName, typeKey{ns: bt.Namespace, name: bt.TypeName})
		delete(c.typeByOID, oid)
		c.removeDepsOf('t', oid)
		c.removeDepsOn('t', oid)
	}
}
