package catalog

// Schema represents a PostgreSQL namespace.
type Schema struct {
	OID       uint32
	Name      string
	Owner     string // explicit AUTHORIZATION / ALTER OWNER (for diff/SDL)
	OwnerName string // effective role owner (for DROP ROLE / DROP OWNED / REASSIGN / RENAME)
	Relations map[string]*Relation
	Indexes   map[string]*Index
	Sequences map[string]*Sequence
}

func (s *Schema) effectiveOwner() string {
	if s.OwnerName != "" {
		return s.OwnerName
	}
	return s.Owner
}

// schemaHasObjects returns true if the schema contains any objects.
func (c *Catalog) schemaHasObjects(s *Schema) bool {
	if len(s.Relations) > 0 || len(s.Indexes) > 0 || len(s.Sequences) > 0 {
		return true
	}
	// Check for user-defined types in this namespace.
	for _, bt := range c.typeByOID {
		if bt.Namespace == s.OID && bt.OID >= FirstNormalObjectId && bt.Elem == 0 && bt.Type != 'm' && bt.Type != 'c' {
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

// removeTypesInSchema removes all user-created types in the given schema.
func (c *Catalog) removeTypesInSchema(s *Schema) {
	var types []*BuiltinType
	for _, bt := range c.typeByOID {
		if bt.Namespace != s.OID || bt.OID < FirstNormalObjectId || bt.Elem != 0 || bt.Type == 'm' || bt.Type == 'c' {
			continue
		}
		types = append(types, bt)
	}
	for _, bt := range types {
		if c.typeByOID[bt.OID] != nil {
			c.deleteUserType(bt)
		}
	}
}
